package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)

	AddRouter(msgID uint32, router IRouter)

	StartWorkerPool()

	// 将消息发送给消息任务队列进行处理
	SendMsgToTaskQueue(request IRequest)
}
