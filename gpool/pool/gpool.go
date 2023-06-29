package pool

import (
	"errors"
	"log"
	"math"
	"os"
	"runtime"
	"time"
)

const (
	// DefaultPoolSize is the default capacity for a default goroutine pool.
	DefaultPoolSize = math.MaxInt32

	// DefaultCleanIntervalTime is the interval time to clean up goroutines. 清理goroutine的间隔时间
	DefaultCleanIntervalTime = time.Second
)

const (
	// OPENED represents that the pool is opened.
	OPENED = iota

	// CLOSED represents that the pool is closed.
	CLOSED
)

var (
	//--------------------------Error types for the Ants API------------------------------

	// ErrLackPoolFunc will be returned when invokers don't provide function for pool.
	ErrLackPoolFunc = errors.New("must provide function for pool")

	// ErrInvalidPoolExpiry will be returned when setting a negative number as the periodic duration to purge goroutines.
	ErrInvalidPoolExpiry = errors.New("invalid expiry for pool")

	// ErrPoolClosed will be returned when submitting task to a closed pool.
	ErrPoolClosed = errors.New("this pool has been closed")

	// ErrPoolOverload will be returned when the pool is full and no workers available.
	ErrPoolOverload = errors.New("too many goroutines blocked on submit or Nonblocking is set")

	// ErrInvalidPreAllocSize will be returned when trying to set up a negative capacity under PreAlloc mode.
	ErrInvalidPreAllocSize = errors.New("can not set up a negative capacity under PreAlloc mode")

	// ErrTimeout will be returned after the operations timed out.
	ErrTimeout = errors.New("operation timed out")
)

var (

	// 判断一个worker的channel是否应该是一个buffered channel来获取最佳性能
	/*
		如果 GOMAXPROCS=1，则使用阻塞 workerChan。这会立即将 Serve 切换到 WorkerFunc，从而提高性能（至少在 go1.5 下）
	 */
	workerChanCap = func() int {
		// runtime.GOMAXPROCS(逻辑CPU数量)
		// <1：不修改任何数值。
		// =1：单核心执行。
		// >1：多核并发执行。
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}
		return 1
	}()

	// 返回的是*Logger, *Logger实现了Logger的Printf，所以*Logger也就是Logger
	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))

	defaultPool, _ = NewPool(DefaultPoolSize)

)

// Logger is used for logging formatted messages
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...interface{})
}

// Submit submits a task to pool.
func Submit(task func()) error {
	return defaultPool.Submit(task)
}

// Running returns the number of the currently running goroutines.
func Running() int {
	return defaultPool.Running()
}

// Cap returns the capacity of this default pool.
func Cap() int {
	return defaultPool.Cap()
}

// Free returns the available goroutines to work.
func Free() int {
	return defaultPool.Free()
}

// Release Closes the default pool.
func Release() {
	defaultPool.Release()
}

// Reboot reboots the default pool.
func Reboot() {
	defaultPool.Reboot()
}
