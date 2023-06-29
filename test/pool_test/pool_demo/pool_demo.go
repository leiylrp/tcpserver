package pool_demo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Person struct {
	Name string
	Age  int
}

func initPool() *sync.Pool {
	cache := &sync.Pool{}
	cache.New = func() interface{} {
		fmt.Println("hello, nice to meet you!")
		return &Person{}
	}
	return cache
}

func PoolDemo1() {
	poolCache := initPool()
	person := poolCache.Get().(*Person)
	fmt.Println("首次从sync.Pool中获取person：", person)

	person.Name = "jack"
	person.Age = 18
	poolCache.Put(person)

	p := &Person{
		Name: "tom",
		Age:  20,
	}
	poolCache.Put(p)
	fmt.Println("设置的对象Name: ", person.Name)
	fmt.Println("设置的对象Age: ", person.Age)

	fmt.Println("Pool 中第一个对象，调用Get方法获取：", poolCache.Get().(*Person))
	fmt.Println("Pool 中第二个对象，调用Get方法获取：", poolCache.Get().(*Person))
	fmt.Println("Pool 中第三个对象，调用Get方法获取: ", poolCache.Get().(*Person))
}

type cache struct {
	pool *sync.Pool
}

// NewCache 加锁自旋和垃圾回收CPU占有率高  主要还是内存，CPU是附带的
func NewCache() *cache {
	poolCache := &sync.Pool{}
	poolCache.New = func() interface{} {
		fmt.Println("hello, nice to meet you!")
		return []byte{}
	}
	return &cache{
		pool: poolCache,
	}
}

type poolCache struct {
	pool   sync.Pool
	cnt    int32
	maxCnt int32
}

func (p *poolCache) Get() interface{} {
	return p.pool.Get()
}

func (p *poolCache) Put(val interface{}) {
	if unsafe.Sizeof(val) > 1024 {
		return
	}
	cnt := atomic.AddInt32(&p.cnt, 1)
	if cnt >= p.maxCnt {
		return
	}
	p.pool.Put(val)
}
