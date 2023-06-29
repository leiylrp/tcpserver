package main

import "fmt"

func main1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("========", err)
		}
		fmt.Println("<--------->")
	}()
	//fmt.Println("result...")
	//var a, b, x int
	//a = 3
	//x = a/b
	//fmt.Println("result...", x)

	for i := 0; i < 100; i++ {
		fmt.Println("========i=====", i)
		if i == 50 {
			panic("已达上限")
		}
	}
}

func main() {
	main1()
	fmt.Println("----panic recover-----")
}
