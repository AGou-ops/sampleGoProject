package main

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/AGou-ops/zinx/ziface"
	"github.com/AGou-ops/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	log.Println("call PreHandle function")
	_, err := request.GetConnection().
		GetTCPConnection().
		Write([]byte("before ping\n"))
	if err != nil {
		log.Println("Err call PreHandle function ", err)
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	log.Println("call Handle function")
	_, err := request.GetConnection().
		GetTCPConnection().
		Write([]byte("main...ping..\n"))
	if err != nil {
		log.Println("Err call Handle function ", err)
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	log.Println("call PostHandle function")
	_, err := request.GetConnection().
		GetTCPConnection().
		Write([]byte("After ping\n"))
	if err != nil {
		log.Println("Err call Posthandle function ", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Server()
}

