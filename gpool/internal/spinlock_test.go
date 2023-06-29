package internal

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type originSpinLock uint32

// NewOriginSpinLock 普通的锁
func NewOriginSpinLock() sync.Locker {
	return new(originSpinLock)
}

func (sl *originSpinLock) Lock() {

	if !atomic.CompareAndSwapUint32((*uint32)(sl), uint32(unlock), uint32(lock)) {
		runtime.Gosched()
	}

}

func (sl *originSpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

/*
	正常模式和饥饿模式 所以说mutex是非公平锁。这也是为了性能考虑。尽量减少runtime的调度
 */
// BenchmarkMutex-8   	24541030	        47.09 ns/op
func BenchmarkMutex(b *testing.B) {
	m := sync.Mutex{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			// nolint:static check
			m.Unlock()
		}
	})
}



// BenchmarkSpinLock-8   	55439286	        23.32 ns/op  执行了5千万次，每次花费23.32ns
func BenchmarkSpinLock(b *testing.B) {
	spin := NewOriginSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}

func BenchmarkBackOffSpinLock(b *testing.B) {
	spin := NewSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}