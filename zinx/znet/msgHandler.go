package znet

import (
	"log"

	"github.com/AGou-ops/zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
