package znet

import "github.com/AGou-ops/zinx/ziface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRouter) {
}

func (br *BaseRouter) Handle(request ziface.IRouter) {
}

func (br *BaseRouter) Posthandle(request ziface.IRouter) {
}
