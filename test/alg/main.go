package main

import "fmt"

func main() {
	a := []int{1,2}
	arr := test(a)
	fmt.Println(arr)
}

func test(arr []int) []int{
	lenArr := len(arr)
	if lenArr == 0 {
		return nil
	}
	help := make([]int,lenArr)
	fmt.Println("=======", len(help), cap(help))
	help[0] = 1
	fmt.Println("---->", help)
	help = append(help, 1)
	return help
}
