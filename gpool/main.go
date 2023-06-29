package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"tcpserver/gpool/pool"
	"time"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main1() {

	defer pool.Release()
	runTimes := 1000
	//runtime.GOMAXPROCS()
	//debug.SetMaxThreads()
	// Use the common pool.
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = pool.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", pool.Running())
	fmt.Printf("finish all tasks.\n")

	// Use the pool with a method,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	//p, _ := pool.NewPoolWithFunc(10, func(i interface{}) {
	//	myFunc(i)
	//	wg.Done()
	//})
	//defer p.Release()
	//// Submit tasks one by one.
	//for i := 0; i < runTimes; i++ {
	//	wg.Add(1)
	//	_ = p.Invoke(int32(i))
	//}
	//wg.Wait()
	//fmt.Printf("running goroutines: %d\n", p.Running())
	//fmt.Printf("finish all tasks, result is %d\n", sum)
	//if sum != 499500 {
	//	panic("the final result is wrong!!!")
	//}
}


func main() {
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