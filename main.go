package main

import (
	"fmt"
	"tcpserver/iface"
	"tcpserver/tnet"
)

type PingRouter struct {
	tnet.BaseRouter
}

func (pr *PingRouter) Handler(request iface.IRequest) {
	fmt.Println("Call PingRouter Handler.....")
	err := request.GetConnection().Send(200, []byte("Ping Client--------------->"))
	if err != nil {
		fmt.Println(err)
	}
}

func (pr *PingRouter) PostHandler(request iface.IRequest)  {
	fmt.Println(".........Call PingRouter PostHandler.........")
}

type HelloRouter struct {
	tnet.BaseRouter
}

func (hr *HelloRouter) PreHandler(request iface.IRequest) {
	fmt.Println("Call HelloRouter PreHandler.....")
	err := request.GetConnection().Send(400, []byte("--------------->Hello world"))
	if err != nil {
		fmt.Println(err)
	}
}

// DoConnectionBegin 连接开始之前执行任务
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called...")
	if err := conn.Send(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println("连接之前失败", err)
	}
	fmt.Println("============>Set conn property...")
	conn.SetProperty("Name", "shine")
	conn.SetProperty("Date", "2022-08-21")
	name, _ := conn.GetProperty("Name")
	date, _ := conn.GetProperty("Date")
	fmt.Println("name-------date", name, date)
}

func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("====> DoConnectionLost is Called...")

	if name, err := conn.GetProperty("Name"); err != nil {
		fmt.Println("Name=", name)
	}

	if data, err := conn.GetProperty("Data"); err != nil {
		fmt.Println("Name=", data)
	}
}

func main() {
	server := tnet.NewServer("TCPSERVER")
	server.SetOnConnStartDo(DoConnectionBegin)
	server.SetOnConnStopDO(DoConnectionLost)
	server.AddRouter(1, &PingRouter{})
	server.AddRouter(2, &HelloRouter{})

	server.Serve()
}

func init() {
	fmt.Println("init successfully")
}
