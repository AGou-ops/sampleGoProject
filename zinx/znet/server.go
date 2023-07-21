package znet

import (
	"fmt"
	"log"
	"net"

	"github.com/AGou-ops/zinx/utils"
	"github.com/AGou-ops/zinx/ziface"
)

// IServer 的实现接口
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	// Router    ziface.IRouter
	MsgHandler ziface.IMsgHandler
	ConnMgr    ziface.IConnManager
}

// 启动服务器
func (s *Server) Start() {
	log.Printf("Server started at: %s:%d \n", s.IP, s.Port)
	log.Printf("%+v", utils.GlobalObject)
	go func() {
		// 开启工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个TCP的Address
		addr, err := net.ResolveTCPAddr(
			s.IPVersion,
			fmt.Sprintf("%s:%d", s.IP, s.Port),
		)
		if err != nil {
			log.Println(err)
			return
		}
		// 监听服务器的地址
		listenrer, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("start server success, ", s.Name)

		var cid uint32 = 0

		for {
			// 服务器端接受数据
			conn, err := listenrer.AcceptTCP()
			if err != nil {
				log.Println("AcceptTCP ERR: ", err)
				continue
			}

			// 如果超过最大链接的数量，那么则关闭新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				log.Println(
					"TOO Many Connections. max connection: ",
					utils.GlobalObject.MaxConn,
				)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	log.Println("zinx server name: ", s.Name, "CLOSED!")
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Server() {
	s.Start()
	// TODO: something else.
	// 阻塞状态
	select {}
}

// AddRouter implements ziface.IServer.
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	log.Println("AddRouter successfully")
}

// 获取当前server的连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
