package tnet

type Message struct {
	// uint32 4字节
	MsgID		uint32
	MsgLen		uint32
	MsgData		[]byte
}

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		MsgID: 		msgID,
		MsgLen: 	uint32(len(data)),
		MsgData: 	data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

func (m *Message) GetMsgData() []byte {
	return m.MsgData
}

func (m *Message) SetMsgID(id uint32) {
	m.MsgID = id
}

func (m *Message) SetMsgLen(lens uint32) {
	m.MsgLen = lens
}

func (m *Message) SetMsgData(data []byte) {
	m.MsgData = data
}

