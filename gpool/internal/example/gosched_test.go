package example

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

// 内核占用了6个系统线程
func TestMain1(t *testing.T) {
	nv := 10
	ov := debug.SetMaxThreads(nv)
	fmt.Println(fmt.Sprintf("Change max threads %d=>%d", ov, nv))

	var wg sync.WaitGroup
	c := make(chan bool, 0)
	for i := 0; i < 10; i++ {
		fmt.Println(fmt.Sprintf("Start goroutine #%v", i))

		wg.Add(1)
		go func() {
			c <- true
			defer wg.Done()
			runtime.LockOSThread()
			time.Sleep(10 * time.Second)
			fmt.Println("Goroutine quit")
		}()

		<- c
		fmt.Println(fmt.Sprintf("Start goroutine #%v ok", i))
	}

	fmt.Println("Wait for all goroutines about 10s...")
	wg.Wait()

	fmt.Println("All goroutines done")
}

func a() {
	fmt.Println("testing!")
}


func TestGoschedTest(t *testing.T) {
	e := make(chan os.Signal, 1)
	//runtime.GOMAXPROCS(4)
	fmt.Println("P numbers", runtime.GOMAXPROCS(0))
	debug.SetMaxThreads(6)
	go GoschedTest("hello")
	go sum()
	go sum1()
	go sum2()
	go sum3()
	go sum4()
	go sum()
	go sum1()
	go sum2()
	go sum3()
	go sum4()
	signal.Notify(e, os.Interrupt)
	<- e
}

func GoschedTest(s string) {
	for i := 0; i < 10; i++ {
		fmt.Println("hh", i, s)
		runtime.Gosched()
	}
}

func sum() {
	fmt.Println("sum")
}
func sum1() {
	fmt.Println("sum1")
}
func sum2() {
	runtime.Gosched()
	fmt.Println("sum2")
}
func sum3() {
	fmt.Println("sum3")
}
func sum4() {
	fmt.Println("sum4")
}

func say(s string) {
	for i := 0; i < 2; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}

func TestSay(t *testing.T) {
	runtime.GOMAXPROCS(1)
	go say("world")
	say("hello")
}

func TestOutput(t *testing.T) {
	//runtime.GOMAXPROCS(1)
	go output("goroutine 2")
	output("goroutine 1")
}

func output(s string){
	for i:=0;i<3;i++{
		//runtime.Gosched()
		fmt.Println(s)
	}
}
