package timer

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/*
  tips:
	一个网络服务程序时需要管理大量客户端连接的，
	其中每个客户端连接都需要管理它的 timeout 时间。
	通常连接的超时管理一般设置为30~60秒不等，并不需要太精确的时间控制。
	另外由于服务端管理着多达数万到数十万不等的连接数，
	因此我们没法为每个连接使用一个Timer，那样太消耗资源不现实。
	用时间轮的方式来管理和维护大量的timer调度，会解决上面的问题。
*/

/*
		map的底层实现原理
		为什么遍历map是无序的？
		如何实现有序遍历map？
		为什么Go map是非线程安全的？
		线程安全的map如何实现?
		Go sync.map 和原生 map 谁的性能好，为什么？
		为什么 Go map 的负载因子是 6.5？
		map扩容策略是什么?

 */

type TimeWheel struct {
	// 时间轮的名称
	name string

	// 刻度的间隔
	interval int64

	// 每个时间轮上的刻度数
	scales int

	// 当前时间指针的指向
	curIndex int

	// 每个刻度所能存放timer的最大容量
	maxCap int

	// 当前时间轮上的所有timer
	timerQueue map[int]map[uint32]*Timer

	// 下一层时间轮
	nextTimeWheel *TimeWheel

	sync.RWMutex
}

func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	fmt.Println("---------", name, interval, scales, maxCap)
	tw := &TimeWheel{
		name: name,
		interval: interval,
		scales: scales,
		maxCap: maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
	}
	for i := 0; i < scales; i++ {
		tw.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}
	return tw
}

/*
	将一个timer定时器加入到分层时间轮中
	tID: 每个定时器timer的唯一标识
	t: 当前被加入时间轮的定时器
	forceNext: 是否强制的将定时器添加到下一层时间轮
	我们采用的算法是：
	如果当前timer的超时时间间隔 大于一个刻度，那么进行hash计算 找到对应的刻度上添加
	如果当前的timer的超时时间间隔 小于一个刻度 :
					如果没有下一轮时间轮
*/

// addTimer 逐级往下
func (tw *TimeWheel) addTimer(tID uint32, timer *Timer, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errstr := fmt.Sprintf("addTimer function err : %s", err)
			return errors.New(errstr)
		}
		return nil
	}()
	delayInterval := timer.unixTs - UnixMilli()

	// 如果当前超时间间隔大于一个刻度的时间间隔
	if delayInterval >= tw.interval {
		// 得到需要跨越的刻度
		dn := delayInterval / tw.interval
		// 因为是环形队列，所以要取余
		tw.timerQueue[(tw.curIndex + int(dn)) % tw.scales][tID] = timer
		return nil
	}
	if delayInterval < tw.interval && tw.nextTimeWheel == nil {
		if forceNext {
			//如果设置为强制移至下一个刻度，那么将定时器移至下一个刻度
			//这种情况，主要是时间轮自动轮转的情况
			//因为这是底层时间轮，该定时器在转动的时候，如果没有被调度者取走的话，该定时器将不会再被发现
			//因为时间轮刻度已经过去，如果不强制把该定时器Timer移至下时刻，就永远不会被取走并触发调用
			//所以这里强制将timer移至下个刻度的集合中，等待调用者在下次轮转之前取走该定时器
			tw.timerQueue[(tw.curIndex+1) % tw.scales][tID] = timer
		} else {
			//如果手动添加定时器，那么直接将timer添加到对应底层时间轮的当前刻度集合中
			tw.timerQueue[tw.curIndex][tID] = timer
		}
		return nil
	}
	if delayInterval < tw.interval {
		return tw.nextTimeWheel.AddTimer(tID, timer)
	}
	return nil
}

//AddTimer 添加一个timer到一个时间轮中(非时间轮自转情况)  lock
func (tw *TimeWheel) AddTimer(tID uint32, timer *Timer) error {
	tw.Lock()
	defer tw.Unlock()
	return tw.addTimer(tID, timer, false)
}

// DeleteTimeWheel 删除一个定时器任务
func (tw *TimeWheel) DeleteTimeWheel(tID uint32) {
	tw.Lock()
	defer tw.Unlock()
	for i := 0; i < tw.scales; i++ {
		if _, ok := tw.timerQueue[i][tID]; ok {
			delete(tw.timerQueue[i], tID)
		}
	}
}

func (tw *TimeWheel) AddTimeWheel(next *TimeWheel) {
	tw.nextTimeWheel = next
	fmt.Println("add next TimeWheel success! cur TimeWheel ", tw, "next TimeWheel", next)
}

/*
	启动时间轮
*/
func (tw *TimeWheel) run() {
	for {
		// sleep
		time.Sleep(time.Duration(tw.interval)*time.Millisecond)
		tw.Lock()
		curTimer := tw.timerQueue[tw.curIndex]
		tw.timerQueue[tw.curIndex] = make(map[uint32]*Timer, tw.maxCap)
		for tID, timer := range curTimer {
			//这里属于时间轮自动转动，forceNext设置为true
			_ = tw.addTimer(tID, timer, true)
		}

		nextTimer := tw.timerQueue[(tw.curIndex+1) % tw.scales]
		tw.timerQueue[(tw.curIndex+1) % tw.scales] = make(map[uint32]*Timer, tw.maxCap)
		for ntID, ntimer := range nextTimer {
			_ = tw.addTimer(ntID, ntimer, true)
		}

		tw.curIndex = (tw.curIndex + 1) % tw.scales
		tw.Unlock()
	}
}

func (tw *TimeWheel) Run() {
	go tw.run()
	fmt.Printf("timerwheel name = %s is running...", tw.name)
}

//GetTimerWithIn 获取定时器在一段时间间隔内的Timer
func (tw *TimeWheel) GetTimerWithIn(duration time.Duration) map[uint32]*Timer {

	leaftw := tw
	for leaftw.nextTimeWheel != nil {
		leaftw = leaftw.nextTimeWheel
	}
	leaftw.Lock()
	defer leaftw.Unlock()
	//返回的Timer集合
	timerList := make(map[uint32]*Timer)

	now := UnixMilli()

	for tID, timer := range leaftw.timerQueue[leaftw.curIndex] {
		if timer.unixTs - now < int64(duration/1e6) {
			timerList[tID] = timer
			delete(leaftw.timerQueue[leaftw.curIndex], tID)
		}
	}
	return timerList
}