package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = false
var globalNum = 1

func main() {
	var wg sync.WaitGroup

	m := sync.Mutex{}
	c := sync.NewCond(&m)
	go func() {
		wg.Add(1)
		// this go routine wait for changes to the sharedRsc
		c.L.Lock()
		for sharedRsc == false {
			fmt.Println("goroutine1 wait")
			c.Wait()
		}
		fmt.Println("goroutine1", sharedRsc)
		time.Sleep(2*time.Second)
		fmt.Printf("goroutine1 modify globalNum = %d\n", globalNum+1)
		c.L.Unlock()
		wg.Done()
	}()

	go func() {
		// this go routine wait for changes to the sharedRsc
		wg.Add(1)
		c.L.Lock()
		for sharedRsc == false {
			fmt.Println("goroutine2 wait")
			c.Wait()
		}
		fmt.Println("goroutine2", sharedRsc)
		time.Sleep(2*time.Second)
		fmt.Printf("goroutine2 modify globalNum = %d\n", globalNum+1)
		c.L.Unlock()
		wg.Done()
	}()

	go add(c, 3, &wg)

	// this one writes changes to sharedRsc
	time.Sleep(2 * time.Second)
	c.L.Lock()
	fmt.Println("main goroutine ready")
	sharedRsc = true
	fmt.Println("main goroutine globalNum--1--", globalNum)
	c.Broadcast()
	fmt.Println("main goroutine broadcast")
	fmt.Println("main goroutine globalNum--2--", globalNum)
	c.L.Unlock()
	wg.Wait()
	fmt.Printf("main globalNum = %d", globalNum)
}

func add(c *sync.Cond, n int, wg *sync.WaitGroup) {
	// this go routine wait for changes to the sharedRsc
	wg.Add(1)
	c.L.Lock()
	for sharedRsc == false {
		fmt.Println("goroutine3 wait")
		c.Wait()
	}
	fmt.Println("goroutine3", sharedRsc)
	time.Sleep(2*time.Second)
	fmt.Printf("goroutine3 modify globalNum = %d\n", globalNum+n)
	c.L.Unlock()
	wg.Done()
}



