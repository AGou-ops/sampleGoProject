package main

import (
	"log"

	"github.com/AGou-ops/zinx/ziface"
	"github.com/AGou-ops/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	log.Println("call Handle function")
	// _, err := request.GetConnection().
	// 	GetTCPConnection().
	// 	Write([]byte("main...ping..\n"))
	// if err != nil {
	// 	log.Println("Err call Handle function ", err)
	// }

	// 先读取客户端的数据
	log.Printf(
		"recv from client: [id: %d]: %s ",
		request.GetMsgID(),
		string(request.GetData()),
	)
	if err := request.GetConnection().SendMsg(0, []byte("ping...")); err != nil {
		log.Println("send error: ", err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (pr *HelloRouter) Handle(request ziface.IRequest) {
	log.Println("call Hello_Handle function")
	// _, err := request.GetConnection().
	// 	GetTCPConnection().
	// 	Write([]byte("main...ping..\n"))
	// if err != nil {
	// 	log.Println("Err call Handle function ", err)
	// }

	// 先读取客户端的数据
	log.Printf(
		"recv from client: [id: %d]: %s ",
		request.GetMsgID(),
		string(request.GetData()),
	)
	if err := request.GetConnection().SendMsg(1, []byte("hello...")); err != nil {
		log.Println("send error: ", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Server()
}
