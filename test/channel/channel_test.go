package main

import (
	"fmt"
	"testing"
)

func TestChanInit(t *testing.T) {
	c := make(chan int)
	ch := new(chan int)
	c <- 18
	fmt.Println(c, ch)
}

func TestChanDemo(t *testing.T) {
	channelDemo()
}

func TestChanConcurrency(t *testing.T) {
	ChanConcurrency()
}
