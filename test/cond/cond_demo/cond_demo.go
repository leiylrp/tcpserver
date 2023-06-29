package cond_demo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)


// Item sync.Cond testing
type Item = int

type Queue struct {
	items     []Item
	itemAdded sync.Cond
}

func NewQueue() *Queue {
	q := new(Queue)
	q.itemAdded.L = &sync.Mutex{}
	return q
}

func (q *Queue) PutArray(items []Item) {
	q.itemAdded.L.Lock()
	q.itemAdded.L.Unlock()
	q.items = append(q.items, items...)
	q.itemAdded.Signal()
}

func (q *Queue) Put(item Item) {
	q.itemAdded.L.Lock()
	defer q.itemAdded.L.Unlock()
	q.items = append(q.items, item)
	q.itemAdded.Signal()
}

func (q *Queue) GetMany(n int) []Item {
	q.itemAdded.L.Lock()
	defer q.itemAdded.L.Unlock()

	for len(q.items) < n {
		q.itemAdded.Wait()
	}

	items := q.items[:n:n]
	q.items = q.items[n:]
	return items
}

type cond struct {
	checker copyChecker
}

type copyChecker uintptr

func (c *copyChecker) check() {
	fmt.Printf("Before: c: %12v, *c: %12v, uintptr(*c): %12v, uintptr(unsafe.Pointer(c)): %12v\n", c, *c, uintptr(*c), uintptr(unsafe.Pointer(c)))
	swapped := atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c)))
	fmt.Printf("After : c: %12v, *c: %12v, uintptr(*c): %12v, uintptr(unsafe.Pointer(c)): %12v, swapped: %12v\n", c, *c, uintptr(*c), uintptr(unsafe.Pointer(c)), swapped)
}
