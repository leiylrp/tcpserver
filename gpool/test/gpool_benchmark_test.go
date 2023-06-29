

package test

import (
	"fmt"
	"runtime"
	"sync"
	"tcpserver/gpool/pool"
	"testing"
	"time"
)

const (
	RunTimes           = 1000000
	BenchParam         = 10
	BenchAntsSize      = 200000
	DefaultExpiredTime = 10 * time.Second
)

const (
	Param    = 100
	AntsSize = 1000
	TestSize = 10000
	n        = 100000
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	// GiB // 1073741824
	// TiB // 1099511627776             (超过了int32的范围)
	// PiB // 1125899906842624
	// EiB // 1152921504606846976
	// ZiB // 1180591620717411303424    (超过了int64的范围)
	// YiB // 1208925819614629174706176
)

var curMem uint64

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func demoTest() {
	fmt.Println("I am testing...")
}

func demoTest1()  {
	//time.Sleep(2*time.Second)
	fmt.Println("T am testing1...")
}

func demoPoolFunc(args interface{}) {
	n := args.(int)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func longRunningFunc() {
	for {
	runtime.Gosched()
}
}

func longRunningPoolFunc(arg interface{}) {
	if ch, ok := arg.(chan struct{}); ok {
		<-ch
		return
	}
	for {
		runtime.Gosched()
	}
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSemaphore(b *testing.B) {
	var wg sync.WaitGroup
	sema := make(chan struct{}, BenchAntsSize)

	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc()
				<-sema
				wg.Done()
			}()
		}
	wg.Wait()
	}
}


func BenchmarkGoroutinesThroughput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			go demoFunc()
		}
	}
}

func BenchmarkSemaphoreThroughput(b *testing.B) {
	sema := make(chan struct{}, BenchAntsSize)
	for i := 0; i < b.N; i++ {
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc()
				<-sema
			}()
		}
	}
}

// TestNoPool 		memory usage:15 MB--- PASS: TestNoPool (0.19s)
func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc()
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestAntsPoolGetWorkerFromCache(t *testing.T) {
	p, _ := pool.NewPool(TestSize)
	defer p.Release()

	for i := 0; i < AntsSize; i++ {
		_ = p.Submit(demoFunc)
	}
	t.Logf("pool, running workers number:%d", p.Running())
	time.Sleep(2 * pool.DefaultCleanIntervalTime)
	_ = p.Submit(demoFunc)
	t.Logf("pool, running workers number:%d", p.Running())
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestSubmit(t *testing.T) {
	defer pool.Release()
	demofunc := func() {
		fmt.Println("I am demofunc...")
	}
	pool.Submit(func() {
		demofunc()
	})
}

// 13MB 0.17s
func TestAntsPool(t *testing.T) {
	defer pool.Release()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		_ = pool.Submit(func() {
			demoTest()
			wg.Done()
		})
	}
	wg.Wait()

	t.Logf("pool, capacity:%d", pool.Cap())
	t.Logf("pool, running workers number:%d", pool.Running())
	t.Logf("pool, free workers number:%d", pool.Free())

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestPool(t *testing.T) {
	defer pool.Release()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_ = pool.Submit(func() {
			go demoTest()
			go demoTest1()
			wg.Done()
		})
	}()
	time.Sleep(time.Second)
	go func() {
		_ = pool.Submit(func() {
			demoTest1()
			wg.Done()
		})
	}()
	wg.Wait()

	t.Logf("pool, capacity:%d", pool.Cap())
	t.Logf("pool, running workers number:%d", pool.Running())
	t.Logf("pool, free workers number:%d", pool.Free())

	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

//func TestPurge(t *testing.T) {
//	p, err := pool.NewPool(10)
//	assert.NoErrorf(t, err, "create TimingPool failed: %v", err)
//	defer p.Release()
//	_ = p.Submit(demoFunc)
//	time.Sleep(3 * pool.DefaultCleanIntervalTime)
//	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
//	p1, err := NewPoolWithFunc(10, demoPoolFunc)
//	assert.NoErrorf(t, err, "create TimingPoolWithFunc failed: %v", err)
//	defer p1.Release()
//	_ = p1.Invoke(1)
//	time.Sleep(3 * DefaultCleanIntervalTime)
//	assert.EqualValues(t, 0, p.Running(), "all p should be purged")
//}
