package iface

type ISession interface {

	AddRouter(msgID uint32, router IRouter)

	DoMsgHandler(request IRequest)

	SendMsgToTaskQueue(request IRequest)

	StartWorkerPool()

	//StartWorkerGPool()
}

// Provider 抽象出来的Provider接口，用于封装CSession底层的存储细节
type Provider interface {
	SessionInit(sid string) (ICSession, error)

	SessionDestroy(sid string)

	SessionRead(sid string) (ICSession, error)

	SessionUpdate(sid string)

	SessionGC(maxLiftTime int64)
}

// ICSession 存储客户端Session
type ICSession interface {

	Set(key, val interface{})

	Get(key interface{}) (val interface{})

	Delete(key interface{})

	CSessionID() string
}