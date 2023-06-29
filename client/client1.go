package main

import (
	"fmt"
	"io"
	"net"
	"tcpserver/tnet"
	"time"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)
	//1. 直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8888")

	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	for {
		dataPack := tnet.NewDataPack()
		binaryWriteData, err := dataPack.Packing(tnet.NewMessage(2, []byte("hello my world... Client Test Message...")))

		if err != nil {
			fmt.Println("Pack error", err)
			return
		}
		if _, err = conn.Write(binaryWriteData); err != nil {
			fmt.Println("client write error", err)
			return
		}
		binaryReadHeadData := make([]byte, dataPack.GetHeaderLen())
		if _, err = io.ReadFull(conn, binaryReadHeadData); err != nil {
			fmt.Println("client read head error", err)
			break
		}

		msgHeadData, err :=  dataPack.Unpacking(binaryReadHeadData)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		if msgHeadData.GetMsgLen() > 0 {
			msg := msgHeadData.(*tnet.Message)
			msg.MsgData = make([]byte, msgHeadData.GetMsgLen())
			if _, err = io.ReadFull(conn, msg.MsgData); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
			fmt.Println("Recv Server Msg: ID = ", msg.GetMsgID(), "len = ", msg.MsgLen, "data = ", string(msg.MsgData))
		}
		//cpu阻塞   不断循环判断，进程会把cpu跑死
		time.Sleep(1 * time.Second)
	}
}
