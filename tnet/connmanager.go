package tnet

import (
	"errors"
	"fmt"
	"sync"
	"tcpserver/iface"
)

type ConnManager struct{

	connections 	map[uint32]iface.IConnection

	connLock		sync.RWMutex
}

func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

func (cm *ConnManager) AddConn(connection iface.IConnection) {
	// map 不是线程安全的
	fmt.Println("添加连接到连接管理器前需要先上锁")
	cm.connLock.Lock()
	fmt.Println("连接管理器添加连接ID", connection.GetConnID(), "上锁成功")
	defer func() {
		cm.connLock.Unlock()
		fmt.Println("连接管理器添加连接ID", connection.GetConnID(), "解锁成功")
	}()
	connID := connection.GetConnID()
	cm.connections[connID] = connection
	fmt.Println("connID", connection.GetConnID(), "add to ConnManager successfully:conn num = ", cm.ConnLen())
}

func (cm *ConnManager) GetConn(connID uint32) (iface.IConnection, error) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection is not exits")
}

func (cm *ConnManager) RemoveConn(conn iface.IConnection)  {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.ConnLen())
}

// ClearConn 清除该server上的所有conn
func (cm *ConnManager) ClearConn() {
	fmt.Println("ClearConn 获取锁...")
	cm.connLock.Lock()
	fmt.Println("ClearConn 上锁成功...")
	defer func() {
		cm.connLock.Unlock()
		fmt.Println("ClearConn 解锁成功!")
	}()
	fmt.Println("for ClearConn...")
	for connID, conn := range cm.connections {
		fmt.Println("服务器清除所有连接")
		conn.Stop()

		delete(cm.connections, connID)
	}
	fmt.Println("===============================================================")
	fmt.Println("Clear All connections success! conn num = ", cm.ConnLen())
}

// ConnLen 得到当前连接总数  Add里方法调用了len方法，加锁的话会死锁
func (cm *ConnManager) ConnLen() int {
	return len(cm.connections)
}

func (cm *ConnManager) RemoveOneConn(connID uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	if conn, ok := cm.connections[connID]; ok {
		conn.Stop()
		delete(cm.connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
	}
	fmt.Println("connID", connID, "remove from ConnManager successfully:conn num = ", cm.ConnLen())
}
