package timer

import (
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"
)

//定义一个超时函数
func myFunc(v ...interface{}) {
	fmt.Printf("No.%d function calld. delay %d second(s)\n", v[0].(int), v[1].(int))
}

func TestTimer(t *testing.T) {
	exitChan := make(chan os.Signal, 1)
	for i := 1; i < 5; i++ {
		go func(i int) {
			timer := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second)
			timer.Run()
		}(i)
	}

	signal.Notify(exitChan, os.Interrupt)
	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	<- exitChan
}

func TestOneTimer(t *testing.T) {
	exitChan := make(chan os.Signal, 1)
	go func(i int) {
		timer := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second)
		timer.Run()
	}(9)

	signal.Notify(exitChan, os.Interrupt)
	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	<- exitChan
}
