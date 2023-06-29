package pkg

import (
	"tcpserver/iface"
)

type GlobalOBJ struct {

	TCPSERVER				iface.IServer

	ServerName				string

	HOST 					string

	PORT					int

	ServerVersions			string

	AllowMaxConn			int

	AllowMaxDataPack		uint32

	AllowMaxTaskQueue		uint32

	WorkerPoolSize			uint32
}

var GlobalObject *GlobalOBJ

//func (g *GlobalOBJ) ReloadUserSettingInfo() {
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("Recovered in ReloadUserSettingInfo", r)
//		}
//	}()
//	jsonData, err := ioutil.ReadFile("tcpserver/tcpserver.json")
//	if err != nil {
//		fmt.Println("-------err", err)
//		panic("读取json数据错误")
//	}
//	err = json.Unmarshal(jsonData, &GlobalObject)
//	if err != nil {
//		panic("unmarshal json err")
//	}
//}

func init() {
	GlobalObject = &GlobalOBJ{
		ServerName:             	"TCPSERVER",
		HOST:                   	"127.0.0.1",
		PORT:                   	 8999,
		ServerVersions:         	 "V1.0",
		WorkerPoolSize: 	 		 10,
		AllowMaxConn:           	 12000,
		AllowMaxDataPack:       	 4096,
		AllowMaxTaskQueue:   		 1024,
	}

	//GlobalObject.ReloadUserSettingInfo()
}
