package znet

import (
	"log"

	"github.com/AGou-ops/zinx/utils"
	"github.com/AGou-ops/zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make(
			[]chan ziface.IRequest,
			utils.GlobalObject.WorkerPoolSize,
		),
	}
}

// DoMsgHandler 调度对应router的处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		log.Println("api msgID: ", request.GetMsgID(), "is not found")
	}
	// handler.PreHandle(request)
	handler.Handle(request)
	// handler.PostHandle(request)
}

// AddRouter 添加路由
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		log.Println("msgID already exists: ", msgID)
		return
	}
	mh.Apis[msgID] = router
	log.Println("Router successfully added: ", msgID)
}

// 启动一个worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolsize开启worker，每一个workder使用一个Goroutine来进行承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 启动一个workder
		// 第0个worker就用第0个channel
		mh.TaskQueue[i] = make(
			chan ziface.IRequest,
			utils.GlobalObject.MaxWorkerTaskLen,
		)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) startOneWorker(
	workerID int,
	TaskQueue chan ziface.IRequest,
) {
	log.Println("WorkerID: ", workerID, " is started!")
	// 不断阻塞等待对应的消息队列的消息
	for {
		select {
		case request := <-TaskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息发送给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	log.Println("Add ConnID: ", request.GetConnection().GetConnID(), "request MsgID: ", request.GetMsgID(), "to workerID: ", workerID)

	// 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
