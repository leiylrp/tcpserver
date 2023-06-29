package pool

import (
	"time"
)

type loopQueue struct {
	items	[]*goWorker
	expiry	[]*goWorker
	head 	int
	tail 	int
	size    int
	isFull	bool
}

func newWorkerLoopQueue(size int) *loopQueue {
	return &loopQueue{
		items: 	make([]*goWorker, size),
		size: 	size,
	}
}

func (wq *loopQueue) len() int {
	if wq.isEmpty() {
		return 0
	}
	if wq.head == wq.tail {
		if wq.isFull {
			return wq.size
		}
		return 0
	}
	if wq.tail > wq.head {
		return wq.tail - wq.head
	}

	return wq.size - wq.head + wq.tail
	/*
			0		1		2		3		4		5		6
			1		1		1		1		1
			head									tail

			1		1		1		0		0		1		1
									tail			head
	 */
}

func (wq *loopQueue) isEmpty() bool {
	return wq.head == wq.tail && !wq.isFull
}

// insert 环形队列插入一个worker
func (wq *loopQueue) insert(worker *goWorker) error {

	if wq.size == 0 {
		return errQueueIsReleased
	}
	if wq.isFull {
		return errQueueIsFull
	}

	wq.items[wq.tail] = worker
	wq.tail++

	if wq.tail == wq.size {
		wq.tail = 0
	}

	if wq.tail == wq.head {
		wq.isFull = true
	}
	return nil
}

// detach 环形队列删除一个worker
func (wq *loopQueue) detach() *goWorker {
	if wq.isEmpty() {
		return nil
	}
	worker := wq.items[wq.head]
	wq.items[wq.head] = nil
	wq.head++
	if wq.head == wq.size {
		wq.head = 0
	}
	wq.isFull = false
	return worker
}

// retrieveExpiry 环形队列清理过期的worker
func (wq *loopQueue) retrieveExpiry(duration time.Duration) []*goWorker {
	expiryTime := time.Now().Add(-duration)

	index := wq.binarySearch(expiryTime)

	if index == -1 {
		return nil
	}

	wq.expiry = wq.expiry[:0]

	if wq.head <= index {
		wq.expiry = append(wq.expiry, wq.items[wq.head:index+1]...)
		for i := wq.head; i <= index; i++ {
			wq.items[i] = nil
		}
	} else {
		wq.expiry = append(wq.expiry, wq.items[wq.head:]...)
		wq.expiry = append(wq.expiry, wq.items[0:index+1]...)
		var i int
		for i = 0; i < index+1; i++ {
			wq.items[i] = nil
		}
		for i = wq.head; i < wq.size; i++ {
			wq.items[i] = nil
		}
	}
	head := (index + 1) % wq.size
	wq.head = head

	if len(wq.expiry) > 0 {
		wq.isFull = false
	}
	return wq.expiry
}
/*
	-------------------->wq.head < index
	nil		nil		1		2		3		4		5		nil
	0		1		2		3		4		5		6		7
					head	index							tail

	-------------------->wq.head > index
	2		3		4		5		nil		nil		nil		1
	0		1		2		3		4		5		6		7
	index							tail					head
---------------------------------------------------------------------------------------------------------

	2		3		4		5		nil		nil		nil		1
	0		1		2		3		4		5		6		7
									tail					head
															baseHead	itLen = 8
			baseMid = (mid + baseHead + itLen) % iLen

--------------------------------------------------------------------------
	1		2		3		4		nil		nil		nil		0
							r								l
			mid = l + (r - l) >> 1
	len = 8

*/
func (wq *loopQueue) binarySearch(expiryTime time.Time) int {
	if wq.isEmpty() || expiryTime.Before(wq.items[wq.head].recycleTime) {
		return -1
	}
	var itsLen, mid, baseMid, baseHead int
	itsLen = len(wq.items)
	r := (itsLen - wq.head + wq.tail - 1) % itsLen
	baseHead = wq.head
	l := 0

	for l <= r {
		mid = l + (r - l) >> 1
		baseMid = (mid + baseHead + itsLen) % itsLen
		if expiryTime.Before(wq.items[baseMid].recycleTime) {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return (r + baseHead + itsLen) % itsLen
}

func (wq *loopQueue) reset() {
	if wq.isEmpty() {
		return
	}
	Releasing:
		if worker := wq.detach(); worker != nil {
			worker.task <- nil
			goto Releasing
		}
		wq.items = wq.items[:0]
		wq.size = 0
		wq.head = 0
		wq.tail = 0
}


