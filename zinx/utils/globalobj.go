package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AGou-ops/zinx/ziface"
)

type globalObj struct {
	TcpServer        ziface.IServer
	Host             string
	TcpPort          int
	Name             string
	Version          string
	MaxConn          int
	MaxPacketSize    uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

var GlobalObject *globalObj

func (g *globalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		log.Println("ReadFile failed: ", err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		log.Println("unmarshal failed: ", err)
	}
}

func init() {
	GlobalObject = &globalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "Zinx server app",
		Version:          "v0.9",
		MaxConn:          2,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	// GlobalObject.Reload()
}
