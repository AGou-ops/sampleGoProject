package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AGou-ops/zinx/ziface"
)

type globalObj struct {
	TcpServer     ziface.IServer
	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxConn       int
	MaxPacketSize int
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
		TcpServer:     nil,
		Host:          "0.0.0.0",
		TcpPort:       8999,
		Name:          "Zinx server app",
		Version:       "v0.4",
		MaxConn:       999999,
		MaxPacketSize: 512,
	}
	GlobalObject.Reload()
}
