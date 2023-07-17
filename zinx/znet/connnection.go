package znet

import (
	"errors"
	"io"
	"log"
	"net"

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
	// Router ziface.IRouter
	MsgHander ziface.IMsgHandler
}

// 初始化连接模块的方法
func NewConnection(
	conn *net.TCPConn,
	ConnID uint32,
	msgHander ziface.IMsgHandler,
) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    ConnID,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		MsgHander: msgHander,
	}
}

func (c *Connection) StartReader() {
	log.Printf("Reader Goroutine is running")
	defer log.Println("connID", c.ConnID, "is Stopped")
	defer c.Stop()

	for {
		// buf := make([]byte, utils.GlobalObject.MaxConn)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	log.Println("Recv buf error: ", err)
		// 	continue
		// }
		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	log.Println("ConnID", c.ConnID, " handler error: ", err)
		// 	break
		// }

		// 创建一个拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的msg HEAD
		headData := make([]byte, dp.GetHeadlen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			log.Println("read message header error: ", err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			log.Println("Unpack error: ", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				log.Println("read msg data error: ", err)
				break
			}
		}

		// 拆包，得到msgId 和 msgdataLen 放到一个msg消息中

		// 根据dataLen，再次读取data，放在msg.data字段中

		req := Request{
			conn: c,
			msg:  msg,
		}
		// 从路由中，找到注册绑定的Conn对应的router调用
		// go func(request ziface.IRequest) {
		// c.Router.PreHandle(request)
		// c.Router.Handle(request)
		// c.Router.PostHandle(request)
		// }(&req)
		go c.MsgHander.DoMsgHandler(&req)

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

// 提供一个SendMsg，先封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	// 将data进行封包，格式msgdataLen, msgID, data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		return errors.New("pack message error" + err.Error())
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		log.Printf("Write msg id %d, error: %s", msgId, err)
		return errors.New("conn Write error")
	}
	return nil
}
