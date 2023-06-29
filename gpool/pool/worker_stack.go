package pool

import (
	"time"
)

type workerStack struct {
	items	[]*goWorker

	expiry	[]*goWorker
}

func newWorkerStack(size int) *workerStack {
	return &workerStack{
		items: make([]*goWorker, size),
	}
}

func (wq *workerStack) len() int {
	// fmt.Println("workerStack len:", len(wq.items))
	return len(wq.items)
}

func (wq *workerStack) isEmpty() bool {
	return len(wq.items) == 0
}

func (wq *workerStack) insert(worker *goWorker) error {
	wq.items = append(wq.items, worker)
	return nil
}

func (wq *workerStack) detach() *goWorker {
	l := wq.len()
	// fmt.Println("worker_stack len: ", l)
	if l == 0 {
		return nil
	}

	w := wq.items[l-1]
	wq.items[l-1] = nil // avoid memory leaks
	wq.items = wq.items[:l-1]

	return w
}

// retrieveExpiry 获取到期的worker
func (wq *workerStack) retrieveExpiry(duration time.Duration) []*goWorker {
	//log.Println("清理过期worker")
	n := wq.len()
	if n == 0 {
		return nil
	}

	expiryTime := time.Now().Add(-duration)

	// 通过二分查找找到最后一个过期的worker
	index := wq.binarySearch(0, n-1, expiryTime)

	// 清除上一批worker
	wq.expiry = wq.expiry[:0]
	if index != -1 {
		wq.expiry = append(wq.expiry, wq.items[:index+1]...)
		m := copy(wq.items, wq.items[index+1:])
		for i := m; i < n; i++ {
			wq.items[i] = nil
		}
		wq.items = wq.items[:m]
	}
	return wq.expiry
}

func (wq *workerStack) binarySearch(l, r int, expiryTime time.Time) int {
	var mid int
	for l <= r {
		mid = l + (r - l) >> 1
		if expiryTime.Before(wq.items[mid].recycleTime) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return r
}

func (wq *workerStack) reset() {
	for i := 0; i < wq.len(); i++ {
		wq.items[i].task <- nil
		wq.items[i] = nil
	}
	wq.items = wq.items[:0]
}

