package fifo

import (
	"container/list"
	"tcpserver/cache/in_process_cache"
)

type fifo struct {
	// maxBytes group cache中使用最大的entry数
	maxBytes int

	// onEvicted 从缓存中移除entry时调用回调函数, 默认nil
	onEvicted func(key string, val interface{})

	// usedBytes 已使用的字节数，只包括value, key不算
	usedBytes int

	ll *list.List
	cache map[string]*list.Element
}

func NewFIFO(maxBytes int, onEvicted func(key string, val interface{})) in_process_cache.Cache {
	return &fifo{
		maxBytes: maxBytes,
		onEvicted: onEvicted,
		ll: list.New(),
		cache: make(map[string]*list.Element),
	}
}

type entry struct {
	key 	string
	value 	interface{}
}

// Len 计算entry中value的长度大小
func (e *entry)	Len() int {
	return in_process_cache.CalcLen(e.value)
}

// Get 获取缓存中的value
func (f *fifo) Get(key string) interface{} {
	if ele, ok := f.cache[key]; ok {
		return ele.Value.(*entry).value
	}
	return nil
}

// Set 缓存中新增value
func (f *fifo) Set(key string, val interface{}) {
	if ele, ok := f.cache[key]; ok {
		f.ll.MoveToBack(ele)
		en := ele.Value.(*entry)
		f.usedBytes = f.usedBytes - in_process_cache.CalcLen(en.value) + in_process_cache.CalcLen(val)
		en.value = val
		return
	}

	en := entry{key: key, value: val}
	ele := f.ll.PushBack(en)
	f.cache[key] = ele
	f.usedBytes = f.usedBytes + en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Del(key string) {
	if ele, ok := f.cache[key]; ok {
		f.removeElement(ele)
	}
}

func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	f.ll.Remove(e)

	en := e.Value.(*entry)

	f.usedBytes -= en.Len()

	delete(f.cache, en.key)

	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

