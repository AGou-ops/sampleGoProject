package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	defer listener.Close()
	if err != nil {
		log.Println("cannnot start a tcp server: ", err)
		return
	}
	log.Println("TCP server started at: ", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		log.Println("New connection: ", conn.RemoteAddr().String())
		if err != nil {
			log.Println("cannnot accept client req: ", err)
			return
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// buf := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed connection: ", conn.RemoteAddr().String())
				break
			}
			log.Println("read error: ", err)
			break
		}

		// 检查消息是否为exit
		if msg == "exit\n" {
			log.Println("Client has sent exit message. Closing connection...")
			conn.Close()
			break
		}
		log.Println("rev msg from client: ", strings.TrimRight(msg, "\n"))

		if _, err = conn.Write([]byte(msg)); err != nil {
			log.Println("write back to client error: ", err)
			break
		}
	}
}
