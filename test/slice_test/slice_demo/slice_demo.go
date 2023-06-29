package slice_demo

import "fmt"

func sliceDemo() {
	slice := []int{0,1,2,3,4,5,6,7,8,9}
	// --------------------
	s1 := slice[2:5]
	s2 := slice[2:6:7]
	// --------------------

	s2 = append(s2, 100)
	// ---------------------

	s2 = append(s2, 200)
	// --------------------
	s1[2] = 20
	// --------------------
}



func myAppend(s []int) {
	s[1] = 100
	fmt.Println(s)

	s = append(s, 100)
	fmt.Println("arr: ", s, "len:", len(s), "cap:", cap(s))

}

func myAppendPtr(s *[]int) {
	*s = append(*s, 100)
}
