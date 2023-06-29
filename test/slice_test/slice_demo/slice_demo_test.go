package slice_demo

import (
	"fmt"
	"testing"
)

func TestSliceDemo(t *testing.T) {
	sliceDemo()
}

func TestMyAppend(t *testing.T) {
	arr := make([]int, 3, 6)
	for i := 0; i < 3; i++ {
		arr[i] = i
	}
	fmt.Println("arr:", arr, "len:", len(arr), "cap:", cap(arr))
	myAppend(arr)
	fmt.Println(arr)
}

func TestMap(t *testing.T) {
	m := make(map[int]int, 10)
	for i := 0; i < 10; i++ {
		m[i] = i
	}
	fmt.Println(m, len(m))
	//for key := range m {
	//	delete(m, key)
	//}
	delete(m, 1)
	fmt.Println(m, len(m))
}