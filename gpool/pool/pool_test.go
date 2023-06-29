package pool

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

// 池中是如何定期清理过期的worker？
/*
	在池中我是通过cancelFunc实现的定时清扫功能
	pool中有一个stopHeartbeat 通过一个定时器设计定时功能,NewTicker，内部,

sendTime接收一个管道作为参数，其主要任务是向管道中写入当前时间。
	将runtimeTimer交给系统线程管理
创建Ticker时生成的管道含有一个缓冲区（make(chan Time, 1)），但是Ticker触发的事件确是周期性的，如果管道中的数据没有被取走，那么sendTime()也不会阻塞，而是直接退出，带来的后果是本次事件会丢失。
	timer 小根堆
	ticker
 */

func TestTicker(t *testing.T) {
	ch := make(chan os.Signal, 1)
	tick := time.NewTicker(2*time.Second)
	var ctx context.Context
	var cancelfunc context.CancelFunc
	ctx, cancelfunc = context.WithCancel(context.Background())
	go ticker(ctx, tick)
	signal.Notify(ch, os.Interrupt)
	<- ch
	cancelfunc()
}

func ticker(ctx context.Context, tick *time.Ticker) {
	for {
		select {
		case <- ctx.Done():
			fmt.Println("over...")
			return
		case <- tick.C:
		}
		fmt.Println("定期")
	}
}

func TestFunc(t *testing.T) {
	testfunc := func() {
		fmt.Println("i am test func")
	}
	testfunc()
}

func TestRetry(t *testing.T) {
	l := sync.Mutex{}
	c := sync.NewCond(&l)
	testfunc := func() {
		fmt.Println("i am test func")
	}
	l.Lock()
	if testfunc == nil {
		fmt.Println("nil")
		} else {
			fmt.Println("normal")
			l.Unlock()
		}
	retry:
		fmt.Println("retry")
		c.Wait()
		if testfunc != nil{
			goto retry
		}
		l.Unlock()
	return
}

// 如何将任务提交到worker？如何获取一个worker执行任务
// 紧急任务你有考虑过吗？优先级的问题。
/*
	多个队列（提高复杂度，多个队列多把锁）
	优先级队列  priority queue 实现是完全可以接受的,调度算法本质上都大同小异,自己做降级不就ok了
	多个 pool 就有很多管理难度，我就需要监控多个 pool，参数里也要调整多个 pool 的大小，而且高优队列可能处于长期闲置状态，在平时浪费资源。
*/
