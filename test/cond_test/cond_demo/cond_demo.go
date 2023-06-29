package cond_demo

import (
	"fmt"
	"sync"
	"time"
)

var shareRecourse = false

func CondTest() {
	var wg sync.WaitGroup
	wg.Add(2)
	l := sync.Mutex{}
	c := sync.NewCond(&l)
	go func() {
		c.L.Lock()
		for shareRecourse == false {
			fmt.Println("goroutine1 wait...")
			c.Wait()
		}
		fmt.Println("goroutine1:", shareRecourse)
		c.L.Unlock()
		wg.Done()
	}()

	go func() {
		c.L.Lock()
		for shareRecourse == false {
			fmt.Println("goroutine2 wait...")
			c.Wait()
		}
		fmt.Println("goroutine2:", shareRecourse)
		c.L.Unlock()
		wg.Done()
	}()

	time.Sleep(2*time.Second)
	c.L.Lock()
	fmt.Println("main goroutine ready")
	shareRecourse = true
	c.Signal()
	c.Signal()
	//c.Broadcast()
	c.L.Unlock()
	wg.Wait()
}
