package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"tcpserver/gpool/internal"
	"time"
)

type Pool struct {

	// 池的容量
	capacity		int32

	// 当前正在运行的goroutine的数量
	running			int32

	// 指数退避实现的自旋锁 worker queue  状态位+CAS+进队列/出队列
	lock 			sync.Locker

	// 通知池自行关闭
	state 			int32

	// 广播 等待获取一个idle worker
	cond 			*sync.Cond

	// 工作队列 存储可用worker的slice
	workers 		workerArray

	// 缓冲池 快速获取可用worker 减少了内存分配，减轻GC压力 垃圾回收，加锁自旋都很消耗CPU
	workerCache		sync.Pool

	waiting			int32

	heartbeatDone	int32

	stopHeartbeat	context.CancelFunc

	options 		*Options
}

// NewPool 初始化一个指数退避实现的自旋锁，初始化缓冲池，初始化定时器定期清理过期的worker
func NewPool(size int, options ...Option) (*Pool, error){
	log.Println("new pool ...")
	opts := loadOptions(options...)

	if size <= 0 {
		size = -1
	}

	// 清理间隔
	if expiry := opts.ExpiryDuration; expiry < 0 {
		return nil, ErrInvalidPoolExpiry
	} else if expiry == 0 {
		opts.ExpiryDuration = DefaultCleanIntervalTime
	}

	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}

	p := &Pool{
		capacity: int32(size),
		lock:     internal.NewSpinLock(),
		options:  opts,
	}
	// New是一个接口， p.workerCache实现了New方法
	p.workerCache.New = func() interface{} {
		return &goWorker{
			Parent: p,
			task: make(chan func(), workerChanCap),
		}
	}
	if p.options.PreAlloc {
		if size == -1 {
			return nil, ErrInvalidPreAllocSize
		}
		p.workers = newWorkerArray(loopQueueType, size)
	} else {
		p.workers = newWorkerArray(stackType, 0)
	}

	// 基于自旋锁的cond  sync.Cond 经常用在多个 goroutine 等待，一个 goroutine 通知（事件发生）的场景。如果是一个通知，一个等待，使用互斥锁或 channel 就能搞定了。
	p.cond = sync.NewCond(p.lock)

	// Start a goroutine to clean up expired workers periodically.
	// 控制链路|安全传递数据 定期清理过期的worker
	var ctx context.Context
	ctx, p.stopHeartbeat = context.WithCancel(context.Background())
	go p.purgePeriodically(ctx)

	return p, nil
}

// purgePeriodically 定期清理过期的worker
// purgePeriodically clears expired workers periodically which runs in an individual goroutine, as a scavenger.
func (p *Pool) purgePeriodically(ctx context.Context){
	log.Println("goroutine 3s间隔定期清理过期worker")
	// Ticker是一个周期触发定时的计时器，它会按照一个时间间隔往channel发送系统当前时间，而channel的接收者可以以固定的时间间隔从channel中读取事件。
	heartbeat := time.NewTicker(3*time.Second)
	defer func() {
		// 停止定时器
		heartbeat.Stop()
		atomic.StoreInt32(&p.heartbeatDone, 1)
	}()

	for {
		select {
		// 从定时器获取数据
		case <-heartbeat.C:
		case <-ctx.Done():
			return
		}

		if p.IsClosed() {
			break
		}

		p.lock.Lock()
		expiredWorkers := p.workers.retrieveExpiry(p.options.ExpiryDuration)
		p.lock.Unlock()

		// Notify obsolete workers to stop.
		// This notification must be outside the p.lock, since w.task
		// may be blocking and may consume a lot of time if many workers
		// are located on non-local CPUs.
		for i := range expiredWorkers {
			expiredWorkers[i].task <- nil
			expiredWorkers[i] = nil
		}

		// There might be a situation where all workers have been cleaned up(no worker is running),
		// or another case where the pool capacity has been Tuned up,
		// while some invokers still get stuck in "p.cond.Wait()",
		// then it ought to wake all those invokers.
		if p.Running() == 0 || (p.Waiting() > 0 && p.Free() > 0) {
			p.cond.Broadcast()
		}
	}
}

// 返回一个可用的worker   如何获取worker
// retrieveWorker returns an available worker to run the tasks. 返回一个可用的worker
func (p *Pool) retrieveWorker() (w *goWorker) {
	spawnWorker := func() {
		w = p.workerCache.Get().(*goWorker)
		fmt.Println("<<<==================从pool中取出worker")
		w.run()
	}
	fmt.Println("等待锁...")
	p.lock.Lock()
	fmt.Println("持有锁")
	w = p.workers.detach()
	if w != nil { // first try to fetch the worker from the queue
		fmt.Println("<<<======================从 worker_array 中获取 可执行的worker对象")
		p.lock.Unlock()
	} else if capacity := p.Cap(); capacity == -1 || capacity > p.Running() {
		// if the worker queue is empty and we don't run out of the pool capacity,
		// then just spawn a new worker goroutine.
		p.lock.Unlock()
		spawnWorker()
	} else { // otherwise, we'll have to keep them blocked and wait for at least one worker to be put back into pool.
		if p.options.Nonblocking {
			p.lock.Unlock()
			return
		}
		fmt.Println("------------retry before")
		retry:
			if p.options.MaxBlockingTasks != 0 && p.Waiting() >= p.options.MaxBlockingTasks {
				p.lock.Unlock()
				return
			}
			fmt.Println("---------------------------------")
			p.addWaiting(1)
			p.cond.Wait() // block and wait for an available worker
			p.addWaiting(-1)

			if p.IsClosed() {
				p.lock.Unlock()
				return
			}

			var nw int
			// 获取正在运行的 worker 的数量 如果没有正在运行的worker，从缓冲池中获取
			if nw = p.Running(); nw == 0 { // awakened by the scavenger
				p.lock.Unlock()
				spawnWorker()
				return
			}
			if w = p.workers.detach(); w == nil {
				if nw < p.Cap() {
					p.lock.Unlock()
					spawnWorker()
					return
				}
				goto retry
			}
			p.lock.Unlock()
		}
	return
}

// 返回worker
// revertWorker puts a worker back into free pool, recycling the goroutines.
func (p *Pool) revertWorker(worker *goWorker) bool {
	log.Println("将执行完的worker放回worker_array")
	if capacity := p.Cap(); (capacity > 0 && p.Running() > capacity) || p.IsClosed() {
		p.cond.Broadcast()
		return false
	}
	worker.recycleTime = time.Now()
	p.lock.Lock()

	// To avoid memory leaks, add a double check in the lock scope.
	// Issue: https://github.com/panjf2000/ants/issues/113
	if p.IsClosed() {
		p.lock.Unlock()
		return false
	}

	err := p.workers.insert(worker)
	if err != nil {
		p.lock.Unlock()
		return false
	}
	time.Sleep(time.Second)
	fmt.Println("<<---------worker_array nums:", p.workers.len())
	// Notify the invoker stuck in 'retrieveWorker()' of there is an available worker in the worker queue.
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

func (p *Pool) addRunning(delta int) {
	atomic.AddInt32(&p.running, int32(delta))
}

func (p *Pool) addWaiting(delta int) {
	atomic.AddInt32(&p.waiting, int32(delta))
}

func (p *Pool) Submit(task func()) error {
	if p.IsClosed() {
		return ErrPoolClosed
	}
	var w *goWorker
	// 如果有可用worker
	fmt.Println("--------send request to worker---------")
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	fmt.Println("--------worker channel receive task--------------")
	w.task <- task
	return nil
}

// Running 返回正在运行的worker数量
func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) Free() int {
	c := p.Cap()
	if c < 0 {
		return -1
	}
	return c - p.Running()
}

func (p *Pool) Waiting() int {
	return int(atomic.LoadInt32(&p.waiting))
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

func (p *Pool) Tune(size int) {
	capacity := p.Cap()
	if capacity == -1 || size <= 0 || size == capacity || p.options.PreAlloc {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
	if size > capacity {
		if size-capacity == 1 {
			p.cond.Signal()
			return
		}
		p.cond.Broadcast()
	}
}

func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.state) == CLOSED
}

func (p *Pool) Release() {
	if !atomic.CompareAndSwapInt32(&p.state, OPENED, CLOSED) {
		return
	}
	p.lock.Lock()
	p.workers.reset()
	p.lock.Unlock()
	// There might be some callers waiting in retrieveWorker(), so we need to wake them up to prevent
	// those callers blocking infinitely.
	p.cond.Broadcast()
}

func (p *Pool) ReleaseTimeout(timeout time.Duration) error {
	if p.IsClosed() || p.stopHeartbeat == nil {
		return ErrPoolClosed
	}

	p.stopHeartbeat()
	p.stopHeartbeat = nil
	p.Release()

	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		if p.Running() == 0 && atomic.LoadInt32(&p.heartbeatDone) == 1 {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ErrTimeout
}

func (p *Pool) Reboot() {
	if atomic.CompareAndSwapInt32(&p.state, CLOSED, OPENED) {
		atomic.StoreInt32(&p.heartbeatDone, 0)
		var ctx context.Context
		ctx, p.stopHeartbeat = context.WithCancel(context.Background())
		go p.purgePeriodically(ctx)
	}
}




