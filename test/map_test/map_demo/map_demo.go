package map_demo

import "fmt"

type Student struct {
	name string
	age  int
}

var list1 map[string]*Student
var list2 map[string]Student

func mapDemo1() {
	stu1 := Student{name: "tom", age: 18}
	stu2 := Student{name: "marry", age: 19}
	list1 = make(map[string]*Student)
	list1["stu1"] = &stu1
	list1["stu2"] = &stu2

	fmt.Println(list1["stu1"])
	list1["stu1"].name = "jack"

	s1 := Student{name: "tom", age: 18}
	list2 = make(map[string]Student)
	list2["s1"] = s1

	tmpS1 := list2["s1"]
	tmpS1.name = "alice"
	list2["s1"] = tmpS1
}

type student struct {
	name string
	age  int
}

func stuDemo() {
	m := make(map[string]*student)
	stus := []student{
		{name: "tom", age: 18},
		{name: "marry", age: 17},
		{name: "jack", age: 16},
	}

	for _, stu := range stus {
		m[stu.name] = &stu
	}
	for n, stu := range m {
		fmt.Println(n, "--->", stu)
	}
}
