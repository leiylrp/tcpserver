package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	haichen()
}

func haichen() {
	ch := make(chan int)
	go func() {
		ch <- 10
		ch <- 20
	}()
	fmt.Println(<-ch)
}

func unInitializationChan() {
	var ch chan int
	//ch <- 1
	//<- ch
	close(ch)
	fmt.Println("close未初始化channel：panic")
}

func unbufferedChan() {
	//ch := make(chan int)
	////fmt.Println("无缓冲channel读：", <-ch)
	//ch <- 1
	//fmt.Println("无缓冲channel读：", <-ch)



	//c := make(chan int)
	//data, ok := <- c
	//if ok {
	//	fmt.Println("data：", data)
	//}
	//c <- 10

	//c := make(chan int)
	//go func() {
	//	data, ok := <- c
	//	if ok {
	//		fmt.Println("data:", data)
	//	}
	//	fmt.Println("=====", data)
	//}()
	//c <- 10

	//fmt.Println("fatal error: all goroutines are asleep - deadlock!")
}

func channelDemo() {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	var m int
	for  {
		select {
		case m = <- ch:
			if m == 1 {
				fmt.Println(m)
			} else {
				ch <- m
			}
		}
		if m == 1 {
			break
		}
	}
	l := len(ch)
	for i := 0; i < l; i++ {
		fmt.Println(<- ch)
	}
	fmt.Println("=================")
}

func ChanConcurrency() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan os.Signal, 1)
	ch1 <- 1
	ch2 <- 1
	go func() {
		fmt.Println("ch1")
		ch1 <- 2
		fmt.Println("ch1 block...")
		ch2 <- 2
	}()
	signal.Notify(ch3, os.Interrupt)
	<- ch3
	fmt.Println("over")
}

