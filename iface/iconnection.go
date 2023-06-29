package iface

import "net"

type IConnection interface {
	Start()

	Stop()

	GetConnOfSocket()	*net.TCPConn

	GetConnID()			uint32

	Send(msgID uint32, data []byte)	error

	RemoteAddr()	net.Addr

	SetProperty(key string, val interface{})

	GetProperty(key string) (interface{}, error)

	RemoveProperty(key string)
}
