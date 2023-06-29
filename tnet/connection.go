package tnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"tcpserver/gpool/pool"
	"tcpserver/iface"
)

type Connection struct {
	Parent			iface.IServer

	ConnSocket		*net.TCPConn

	ConnID			uint32

	isClose			bool

	ExitChan		chan bool

	msgChan			chan []byte

	MsgHandler		iface.ISession

	property 		map[string]interface{}

	propertyLock	sync.RWMutex
}


func NewConnection(parent iface.IServer, connSocket *net.TCPConn, connID uint32, msgHandler iface.ISession) *Connection {
	c := &Connection{
		Parent: 		parent,
		ConnSocket: 	connSocket,
		ConnID: 		connID,
		isClose: 		false,
		ExitChan:      	make(chan bool, 1),
		msgChan:     	make(chan []byte),
		MsgHandler: 	msgHandler,
		property: 		make(map[string]interface{}),
	}
	fmt.Println("NewConnection:", c.ConnID)
	c.Parent.GetConnManager().AddConn(c)
	return c
}

// StartReader 开始读取连接的消息，需要对tcp包进行处理
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer func() {
		fmt.Println("读关闭")
		c.Stop()
	}()
	for {

		dataPack := NewDataPack()
		headLen := dataPack.GetHeaderLen()

		// 只有8字节大小，第一次只读取八字节
		headData := make([]byte, headLen)
		connSocket := c.GetConnOfSocket()
		_, err := io.ReadFull(connSocket, headData)
		if err != nil {
			return
		}

		msg, err := dataPack.Unpacking(headData)
		if err != nil {
			return
		}
		var msgData []byte
		if msg.GetMsgLen() > 0 {
			msgData = make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(connSocket, msgData); err != nil {
				return
			}
		}
		msg.SetMsgData(msgData)
		// 将消息和连接封装成一个request，加入到协程池
		req := &Request{
			connection: c,
			message: msg,
		}
		//fmt.Println("==========", &req)
		//if pkg.GlobalObject.WorkerPoolSize > 0 {
		//	c.MsgHandler.SendMsgToTaskQueue(req)
		//} else {
		//	go c.MsgHandler.DoMsgHandler(req)
		//}
		taskFunc := func() {
			c.MsgHandler.DoMsgHandler(req)
		}
		err = pool.Submit(taskFunc)
		if err != nil {
			panic(err)
		}
	}

}

func (c *Connection) StartWriter() {
	for {
		select {
		case data := <- c.msgChan:
			if _, err := c.GetConnOfSocket().Write(data); err != nil {
				fmt.Println("写数据失败，err", err)
				return
			}
		case <- c.ExitChan:
			return
		}
	}
}


/*
					 ________ go c.StartReader() --->拆包  SendMagToTasksQueue or DoMsgHandler
					 |
go conn.Start()------|----c.Parent.CallOnConnStartDo(c)
					 |
					 ________ go c.StartWriter() --->for { select{ case } }
 */


func (c *Connection) Start()  {

	go c.StartReader()

	go c.StartWriter()

	c.Parent.CallOnConnStartDo(c)
}

func (c *Connection) Stop() {

	c.Parent.CallOnConnStopDo(c)

	if c.isClose == true {
		return
	}

	fmt.Println("--------------------------------------------------->")
	c.Parent.GetConnManager().RemoveConn(c)
	// c.Parent.GetConnManager().ClearConn()
	c.isClose = true

	c.ExitChan <- true

	err := c.ConnSocket.Close()
	if err != nil {
		panic(err)
	}
	close(c.msgChan)
	close(c.ExitChan)
	fmt.Println("------------->")
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// Send 开始写消息，需要对tcp包进行处理
func (c *Connection) Send(msgID uint32, data []byte) error {
	if c.isClose {
		return errors.New("conn is close")
	}
	dataPack := NewDataPack()

	// 头部占8字节，
	dataMsg, err := dataPack.Packing(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("数据封包err", err)
		return err
	}

	c.msgChan <- dataMsg
	return nil
}

func (c *Connection) GetConnOfSocket() *net.TCPConn {
	return c.ConnSocket
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = val
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if val, ok := c.property[key]; ok {
		return val, nil
	}
	return nil, errors.New("nil")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}


