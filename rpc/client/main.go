package main

import (
	"fmt"
	"net/rpc"
)

type Req struct {
	Num1, Num2 int
}

type Res struct {
	Sum int
}

func main() {
	req := Req{Num1: 1, Num2: 2}
	var res Res
	client, err := rpc.DialHTTP("tcp", ":8888")
	if err != nil {
		fmt.Println("DialHTTP error:", err)
		return
	}
	call := client.Go("Server.Add", req, &res, nil)
	<- call.Done
	fmt.Println("res:", res.Sum)
}
