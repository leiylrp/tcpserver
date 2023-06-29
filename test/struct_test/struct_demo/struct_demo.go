package struct_demo

import (
	"fmt"
)

type s struct {
	name string
	age  []int
}

func structDemo() {
	s1 := s{
		name: "tom",
		age:  []int{11, 12},
	}
	s2 := s{
		name: "tom",
		age:  []int{22, 33},
	}
	//if s1 == s2 {
	//	fmt.Println("===========")
	//	return
	//}
	fmt.Println("-------------", s1, s2)
}



func s12() {
	st1 := struct {
		name string
		age int
	}{
		name: "tom",
		age: 18,
	}

	st2 := struct {
		name string
		age int
	}{
		name: "marry",
		age: 19,
	}
	if st1 == st2 {
		fmt.Println("xxxxxxxxxxx")
	}
}
