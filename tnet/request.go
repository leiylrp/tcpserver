package tnet

import "tcpserver/iface"

type Request struct {

	connection 	iface.IConnection

	message 	iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.connection
}

func (r *Request) GetMsgID() uint32 {
	return r.message.GetMsgID()
}

func (r *Request) GetData() []byte {
	return r.message.GetMsgData()
}
