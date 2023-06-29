package cond_demo

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCondDemo(t *testing.T) {
	q := NewQueue()
	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			items := q.GetMany(n)
			fmt.Printf("%2d: %2d\n", n, items)
			wg.Done()
		}(n)
	}
	time.Sleep(2*time.Second)
	arr := make([]Item, 100)
	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}
	q.PutArray(arr)
	wg.Wait()
}


type S struct {
	f1 int
	f2 *s
}

type s struct {
	name string
}

// TestNocopy
func TestNocopy(t *testing.T) {
	sOld := S{
		f1: 2,
		f2: &s{name: "tom"},
	}
	sNew := sOld
	sNew.f1 = 1
	sNew.f2.name = "marry"
	fmt.Println(sOld.f1, sOld.f2)
}

// TestRuntimeChecking Builder{addr *Builder  buf []byte}  b.addr = b   a.addr = b ==> panic
func TestRuntimeChecking(t *testing.T) {
	var b strings.Builder
	b.Write([]byte("a"))
	a := b
	a.Write([]byte("b")) // 调用write 时，需要进行nocopy检查，由于b中的addr是指针类型，所以复制之后a.addr 指向的还是b，所以panic
	fmt.Println("a", a)
}





func TestCondNocopy(t *testing.T) {

}