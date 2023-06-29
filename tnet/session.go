package tnet

import (
	"fmt"
	"strconv"
	"tcpserver/iface"
	"tcpserver/pkg"
)

type Session struct {
	RouterApis		map[uint32]iface.IRouter

	TaskQueue		[]chan	iface.IRequest

	WorkerPoolSize	uint32

	//lock 			sync.Mutex
}

func NewMsgHandler() *Session {
	return &Session{
		RouterApis:     make(map[uint32]iface.IRouter),
		TaskQueue:      make([]chan iface.IRequest, pkg.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: pkg.GlobalObject.WorkerPoolSize,
	}
}

func (mh *Session) AddRouter(msgID uint32, router iface.IRouter) {
	if _, ok := mh.RouterApis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.RouterApis[msgID] = router
	fmt.Println("Session   Add api MsgID = ", msgID, "success !")
}

func (mh *Session) DoMsgHandler(request iface.IRequest) {
	reqHandler, ok := mh.RouterApis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! NEED Register!")
		return
	}
	fmt.Println("msgID,  connection, data: ", request.GetMsgID(), string(request.GetData()))
	reqHandler.PreHandler(request)
	reqHandler.Handler(request)
	reqHandler.PostHandler(request)
}

//func (mh *Session) StartWorkerGPool() {
//
//	defer pool.Release()
//	runTimes := 5
//	var wg sync.WaitGroup
//
//	syncCalculateSum := func() {
//		fmt.Println("----------------> Hello World <-----------------")
//		wg.Done()
//	}
//
//	for i := 0; i < runTimes; i++ {
//		wg.Add(1)
//		_ = pool.Submit(syncCalculateSum)
//	}
//	wg.Wait()
//	fmt.Printf("running goroutines: %d\n", pool.Running())
//	fmt.Printf("finish all tasks.\n")
//}

func (mh *Session) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 别晕 一个worker对应一个queue，也就是WorkerPoolSize 与 queue的数量是一致的
		// 一个chan queue 里的大小是AllowMaxTaskQueue
		mh.TaskQueue[i] = make(chan iface.IRequest, pkg.GlobalObject.AllowMaxTaskQueue)
		go mh.startOneWorker(uint32(i), mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 发送数据到任务队列
func (mh *Session) SendMsgToTaskQueue(request iface.IRequest) {
	workID := request.GetConnection().GetConnID() % pkg.GlobalObject.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "request MsgID = ", request.GetMsgID(), "to WorkerID = ", workID, "request data = ", string(request.GetData()))
	mh.TaskQueue[workID] <- request
}

// startOneWorker 每个任务一个goroutine
func (mh *Session) startOneWorker(workID uint32, taskChan chan iface.IRequest) {
	fmt.Println("WorkID = ", workID, "is start...")
	for  {
		select {
		case request := <- taskChan:
			mh.DoMsgHandler(request)
		}
	}
}


