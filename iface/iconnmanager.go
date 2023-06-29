package iface

type IConnManager interface {
	AddConn(connection IConnection)

	GetConn(connID uint32)	(IConnection, error)

	RemoveConn(conn IConnection)

	ClearConn()

	ConnLen() int
}
