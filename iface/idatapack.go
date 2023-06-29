package iface

type IDataPack interface {

	GetHeaderLen() uint32

	Packing(message IMessage) ([]byte ,error)

	Unpacking(binaryMessage []byte) (IMessage, error)
}
