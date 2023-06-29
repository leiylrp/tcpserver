package tnet

import (
	"fmt"
	"net"
	"tcpserver/iface"
	"tcpserver/pkg"
)

type Server struct {
	Name          string
	IPVersion     string
	IP            string
	Port          int
	ConnManager   iface.IConnManager
	Router        iface.IRouter
	Session       iface.ISession
	DoOnConnStart func(connection iface.IConnection)
	DoOnConnStop  func(connection iface.IConnection)

}

func NewServer(name string) *Server {
	return &Server{
		Name:        pkg.GlobalObject.ServerName,
		IPVersion:   "tcp4",
		IP:          pkg.GlobalObject.HOST,
		Port:        pkg.GlobalObject.PORT,
		ConnManager: NewConnManager(),
		Session:     NewMsgHandler(),
	}
}

func (s *Server) Start()  {
	go func() {
		// TODO:开启协程池
		//s.Session.StartWorkerPool()
		// s.Session.StartWorkerGPool()

		tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8888")
		if err != nil {
			fmt.Println("server socket and bind err", err)
		}

		listener, err := net.ListenTCP("tcp4", tcpAddr)
		if err != nil {
			fmt.Println("server listen err", err)
			return
		}
		fmt.Println("start tcp server success, ", s.Name, "success, Listenning...")
		var connID uint32
		for {
			// 获得连接的文件描述符socket fd
			tcpConn, err := listener.AcceptTCP()
			connID++
			fmt.Println("connID:", connID)
			if err != nil {
				fmt.Println("server accept arr", err)
			}
			if s.ConnManager.ConnLen() > pkg.GlobalObject.AllowMaxConn {
				fmt.Println("====> too many Connections MaxConn = ", pkg.GlobalObject.AllowMaxConn)
				tcpConn.Close()
				continue
			}
			fmt.Println("当前连接数小于server中允许的最大连接数")

			// 获得当前连接的对象，进行处理
			conn := NewConnection(s, tcpConn, connID, s.Session)
			fmt.Println("开启goroutine处理连接内的业务消息...")
			go conn.Start()
		}
	}()
}
/*
main.go(server.Serve)---------------------------------------select{}
						|
						-----go s.Start()----()这里开一个协程是为了执行完start的主逻辑之后能回到主协程执行其它逻辑
											|
											------go func()---- for()
															  |	go conn.Start()------
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()
															  |	go conn.Start()

 */

// Serve 控制服务器的连接与关闭等主要功能
func (s *Server) Serve() {
	fmt.Println("server start to run")
	go s.Start()

	select { }
}

func (s *Server) Stop()  {
	fmt.Println("关闭连接ing")
	s.GetConnManager().ClearConn()
}

func (s *Server) GetConnManager() iface.IConnManager {
	return s.ConnManager
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.Session.AddRouter(msgID, router)
}

func (s *Server) SetOnConnStartDo(hookFunc func(connection iface.IConnection)) {
	s.DoOnConnStart = hookFunc
}

func (s *Server) SetOnConnStopDO(hookFunc func(connection iface.IConnection)) {
	s.DoOnConnStop = hookFunc
}

func (s *Server) CallOnConnStartDo(connection iface.IConnection) {
	if s.DoOnConnStart != nil {
		fmt.Println("call on conn start")
		s.DoOnConnStart(connection)
	}
}

func (s *Server) CallOnConnStopDo(connection iface.IConnection) {
	if s.DoOnConnStop != nil {
		fmt.Println("call on conn stop")
		s.DoOnConnStop(connection)
	}
}





