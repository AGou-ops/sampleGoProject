package znet

import (
	"log"
	"net"

	"github.com/AGou-ops/zinx/utils"
	"github.com/AGou-ops/zinx/ziface"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 连接是否关闭
	isClosed bool

	// 告知当前连接已经退出的channel
	ExitChan chan bool

	// 当前连接的处理router
	Router ziface.IRouter
}

// 初始化连接模块的方法
func NewConnection(
	conn *net.TCPConn,
	ConnID uint32,
	router ziface.IRouter,
) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   ConnID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
}

func (c *Connection) StartReader() {
	log.Printf("Reader Goroutine is running")
	defer log.Println("connID", c.ConnID, "is Stopped")
	defer c.Stop()

	for {
		buf := make([]byte, utils.GlobalObject.MaxConn)
		_, err := c.Conn.Read(buf)
		if err != nil {
			log.Println("Recv buf error: ", err)
			continue
		}
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	log.Println("ConnID", c.ConnID, " handler error: ", err)
		// 	break
		// }

		req := Request{
			conn: c,
			data: []byte{},
		}
		// 从路由中，找到注册绑定的Conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// 启动连接
func (c *Connection) Start() {
	log.Println("Connection started: ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	// TODO: 启动从当前连接写数据的业务
}

// 停止连接
func (c *Connection) Stop() {
	log.Println("Connection stopped: ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取客户端的TCP状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) Send(data []byte) error {
	return nil
}
