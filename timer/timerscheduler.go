package timer

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// 时间轮调度器  ---> timer delayFunc timeWheel
const (
	// MaxChanBuffer 最大缓冲出发函数队列大小
	MaxChanBuffer = 2048

	// MaxTimeDelay 默认最大误差时间
	MaxTimeDelay = 100
)

type TimerScheduler struct {

	// 当前调度器的最高时间轮
	tw *TimeWheel

	// ID累加器
	IDGen	uint32

	// 已经触发定时器的channel
	triggerChan	chan *DelayFunc

	// 读写锁
	sync.RWMutex
}

// NewTimerScheduler 时间轮调度器
func NewTimerScheduler() *TimerScheduler {

	secondTw := NewTimeWheel(SecondName, SecondInterval, SecondScales, TimersMaxCap)

	minuteTw := NewTimeWheel(MinuteName, MinuteInterval, MinuteScales, TimersMaxCap)

	hourTw := NewTimeWheel(HourName, HourInterval, HourScales, TimersMaxCap)

	hourTw.AddTimeWheel(minuteTw)
	minuteTw.AddTimeWheel(secondTw)

	secondTw.run()
	minuteTw.run()
	hourTw.run()

	return &TimerScheduler{
		tw: hourTw,
		triggerChan: make(chan *DelayFunc, MaxChanBuffer),
	}
}

func (ts *TimerScheduler) CreateTimerAt(df *DelayFunc, unixNano int64) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()

	ts.IDGen++
	return ts.IDGen, ts.tw.AddTimer(ts.IDGen, NewTimerAt(df, unixNano))
}

func (ts *TimerScheduler) CreateTimerAfter(df *DelayFunc, duration time.Duration) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()

	ts.IDGen++
	return ts.IDGen, ts.tw.AddTimer(ts.IDGen, NewTimerAfter(df, duration))
}

func (ts *TimerScheduler) CancelTimer(tID uint32) {
	ts.Lock()
	defer ts.Unlock()

	tw := ts.tw
	for tw != nil {
		tw.DeleteTimeWheel(tID)
		tw = tw.nextTimeWheel
	}
}

func (ts *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return ts.triggerChan
}

func (ts *TimerScheduler) Start() {
	go func() {
		for {
			now := UnixMilli()
			timers := ts.tw.GetTimerWithIn(MaxTimeDelay * time.Millisecond)
			for _, timer := range timers {
				if math.Abs(float64(now-timer.unixTs)) > MaxTimeDelay {
					fmt.Println("want call at ", timer.unixTs, "; real call at", now, "; delay ", now-timer.unixTs)
					return
				}
				ts.triggerChan <- timer.delayFunc
			}
			time.Sleep(MaxTimeDelay / 2 * time.Millisecond)
		}
	}()
}

//func NewAutoExecTimerScheduler() *TimerScheduler {
//	autoExecScheduler := NewTimerScheduler()
//	autoExecScheduler.Start()
//
//	go func() {
//		delayFunc := autoExecScheduler.triggerChan
//		for df := range delayFunc {
//			go df.Call()
//		}
//	}()
//	return autoExecScheduler
//}


