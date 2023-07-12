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
	Router    ziface.IRouter
}

// 启动服务器
func (s *Server) Start() {
	log.Printf("Server started at: %s:%d \n", s.IP, s.Port)
	log.Printf("%+v", utils.GlobalObject)
	go func() {
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

			// go func() {
			// 	for {
			// 		buf := make([]byte, 512)
			// 		cnt, err := conn.Read(buf)
			// 		if err != nil {
			// 			log.Println(err)
			// 			continue
			// 		}
			// 		fmt.Println("recv client buf: ", string(buf), cnt)
			// 		if _, err := conn.Write(buf[:cnt]); err != nil {
			// 			log.Println("Write Back error: ", err)
			// 			continue
			// 		}
			//
			// 	}
			// }()
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO: 停止服务
}

// 运行服务器
func (s *Server) Server() {
	s.Start()
	// TODO: something else.
	// 阻塞状态
	select {}
}

// AddRouter implements ziface.IServer.
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	log.Println("AddRouter successfully")
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
