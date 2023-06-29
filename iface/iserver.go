package iface

type IServer interface {

	Start()

	Serve()

	Stop()

	GetConnManager() IConnManager

	AddRouter(msgID uint32, router IRouter)

	SetOnConnStartDo(func(connection IConnection))

	SetOnConnStopDO(func(connection IConnection))

	CallOnConnStartDo(connection IConnection)

	CallOnConnStopDo(connection IConnection)
}
