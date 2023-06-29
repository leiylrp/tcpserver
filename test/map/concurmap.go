package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var lock sync.Mutex
	m := make(map[int]int)
	go func() {
		//开一个协程写map
		for i := 0; i < 1000; i++ {
			lock.Lock()
			m[i] = i
			fmt.Println("======write map", i)
			lock.Unlock()
		}
	}()

	go func() {
		//开一个协程读map
		for i := 0; i < 1000; i++ {
			lock.Lock()
			fmt.Println("======read map",m[i])
			lock.Unlock()
		}
	}()

	//time.Sleep(time.Second * 20)
	time.Sleep(time.Second * 20)
}
