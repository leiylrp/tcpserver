package internal

import (
	"runtime"
	"sync"
	"sync/atomic"
)

/*
	任何协程不安全的都可以用读写锁+装饰器模式实现  自旋 等待 获取锁 锁竞争
	Mutex state sema sema+1挂起，sema-1唤醒
 */
// spinLock 自旋锁 实现了Lock和Unlock，所以spinLock也是一个Lock  0没上锁 1上锁 为了尽快拿到锁
// 进队列再被唤醒是很慢的，这里由CPU换时间
type spinlock uint32

const (
	maxBackOff = 16
)

type lockState uint32

const (
	unlock	lockState = iota		// 无锁
	lock							// 有锁
)

// Lock 使用了指数退避算法进行上锁
func (sl *spinlock) Lock() {
	var backoff = 1
	// 通过CAS原子操作一次性操作成功，如果不符合条件，则指数自旋
	for !atomic.CompareAndSwapUint32((*uint32)(sl), uint32(unlock), uint32(lock)) {
		for i := 0; i < backoff; i++ {
			// 调度g
			runtime.Gosched()
		}
		if backoff < maxBackOff {
			backoff = backoff << 1
		}
	}
}

func (sl *spinlock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), uint32(unlock))
}

func NewSpinLock() sync.Locker {
	return new(spinlock)
}

/*
	为什么要设计成这样一个指数退避的自旋锁？
 */
