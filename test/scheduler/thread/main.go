package main

import (
	"runtime/debug"
)

func main() {
	//var x int
	//runtime.GOMAXPROCS(0)
	debug.SetMaxThreads(6)
	//fmt.Println(threads, P)
	//for i := 0; i < threads-5; i++ {
	//	go func() {
	//		for {
	//			x++
	//		}
	//	}()
	//}
	//time.Sleep(time.Second)

}
