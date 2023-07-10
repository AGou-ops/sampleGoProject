package ziface

// 路由的抽象接口
type IRouter interface {
	// 处理connection之前的hook方法
	PreHandle(request IRequest)

	// 处理connection业务的主要hook方法
	Handle(request IRequest)

	// 处理connection业务之后的主要hook方法
	PostHandle(request IRequest)
}
