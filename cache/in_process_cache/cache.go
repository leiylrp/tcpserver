package in_process_cache

import (
	"fmt"
	"runtime"
)

type Cache interface {
	// Get 从fifo缓存中获取元素
	Get(key string) interface{}

	// Set 从fifo缓存中新增元素
	Set(key string, val interface{})

	// Del 删除元素
	Del(key string)

	// DelOldest 删除链表头元素
	DelOldest()

	// Len 缓存的大小
	Len() int
}

// Value *list.Element 中Value的长度
type Value interface {
	Len()	int
}

// CalcLen 计算占用内存
func CalcLen(val interface{}) int {
	var n int
	switch v := val.(type) {
	case Value:
		n = v.Len()

	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 8 + len(v)
		}

	case bool, uint8, int8:
		n = 1

	case int16, uint16:
		n = 2

	case int32, uint32, float32:
		n = 4

	case int64, uint64, float64:
		n = 8

	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 8
		} else {
			n = 4
		}

	case complex64:
		n = 8

	case complex128:
		n = 16

	default:
		panic(fmt.Sprintf("%T is not implement cache.Value", val))
	}
	return n
}