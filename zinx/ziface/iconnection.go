package ziface

import "net"

type IConnection interface {
	// 启动连接
	Start()

	// 停止连接
	Stop()

	// 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接ID
	GetConnID() uint32

	// 获取客户端的TCP状态
	RemoteAddr() net.Addr

	// 发送数据
	SendMsg(msgId uint32, data []byte) error
}

//
// // 处理连接业务的方法
// type HandleFunc func(*net.TCPConn, []byte, int) error
