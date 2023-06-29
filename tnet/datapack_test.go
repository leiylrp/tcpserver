package tnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 创建socket TCP  守候socket
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建一个 go 负责从客户端处理业务
	go func() {
		// 从客户端读取数据， 拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
				return
			}
			/*
				    conn, err := listenner.Accept()
				    go fun(conn net.Conn) {

					}(conn)
					这里是一个匿名函数，同时返回的也是一个闭包  闭包 = 函数 + 外层变量的引用 内部函数引用外部函数变量
					// 闭包：一个函数和与其相关的引用环境组合的一个整体（实体）
			*/
			go func(conn net.Conn) {
				// 处理客户端的请求

				// ---->拆包过程<------

				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeaderLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}

					// headData 是读到的二进制流的长度
					msgHead, err := dp.Unpacking(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}

					// 如果msgHead 是有数据的，需要进行第二次读取
					if msgHead.GetMsgLen() > 0 {
						// 第二次从conn读，根据head中的datalen读取data内容
						// 类型断言：将接口转会具体的数据类型
						msg := msgHead.(*Message)
						msg.MsgData = make([]byte, msg.GetMsgLen())

						// 根据dataLen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.MsgData)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						// 完整的一个消息已经读取完毕
						fmt.Println("----> Recv MsgID:", msg.MsgID, ", datalen = ", msg.MsgLen, ", data = ", string(msg.MsgData))
					}

				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	// 创建一个封包对象dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	// 封装的一个msg1包
	msg1 := &Message{
		MsgID: 1,
		MsgLen:   5,
		MsgData:      []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendDatal, err := dp.Packing(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	//封装第二个msg2包
	msg2 := &Message{
		MsgID: 2,
		MsgLen:   8,
		MsgData:      []byte{'m', 'y', 'w', 'o', 'r', 'l', 'd', '!'},
	}
	sendDatal2, err := dp.Packing(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}

	// 将两个包粘在一起
	sendDatal = append(sendDatal, sendDatal2...)

	// 一次性发送给服务端
	conn.Write(sendDatal)

	// 客户端阻塞
	select {}

}

