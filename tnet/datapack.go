package tnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"tcpserver/iface"
)

type DataPack struct {}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// Packing |ID|Len|Data|
func (dp *DataPack) Packing(message iface.IMessage) ([]byte, error) {
	// 创建一个缓冲区， 小端序将数据以id+len+data的写入到buffer缓冲区
	dataBuffer := bytes.NewBuffer([]byte{})
	// 注意网络传输的数据是以小端序的
	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

func (dp *DataPack) Unpacking(binaryMessage []byte) (iface.IMessage, error) {
	reader := bytes.NewReader(binaryMessage)

	msg := &Message{}

	// interface 类型要记得加地址 否则报 binary.Read: invalid type uint32
	if err := binary.Read(reader, binary.LittleEndian, &msg.MsgID); err != nil {
		fmt.Println("========1 unpack err=====", err)
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.MsgLen); err != nil {
		fmt.Println("========2 unpack err=====", err)
		return nil, err
	}

	return	msg, nil
}

func (dp *DataPack) GetHeaderLen() uint32 {
	return 8
}


