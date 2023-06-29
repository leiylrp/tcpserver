package waitgroup_demo

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 同步多个goroutine之间的工作

func WaitGroupTest() {
	wg := sync.WaitGroup{}
	var result int64 = 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(delta int) {
			defer wg.Done()
			atomic.AddInt64(&result, int64(delta))

		}(i)
	}
	wg.Wait()
	fmt.Println(result)
}
