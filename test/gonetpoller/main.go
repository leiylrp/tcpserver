package gonetpoller

import "net"

func GoNetPoller() {
	listener, err := net.Listen("tcp", ":8888")
	conn, err := listener.Accept()
	conn.Read()
}
