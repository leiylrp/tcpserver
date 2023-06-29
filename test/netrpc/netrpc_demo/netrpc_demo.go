package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Server struct {
}

type Req struct {
	Num1, Num2 int
}

type Res struct {
	Sum   int
	Multi int
}

func (s *Server) Add(req Req, res *Res) error {
	res.Sum = req.Num1 + req.Num2
	fmt.Println("rpc process...")
	return nil
}

func main() {
	_ = rpc.Register(new(Server))
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	if err := http.Serve(l, nil); err != nil {
		fmt.Println("server Serve error: ", err)
		os.Exit(1)
	}
	fmt.Println(l)
}
