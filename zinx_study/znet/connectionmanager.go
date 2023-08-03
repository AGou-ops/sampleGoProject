package znet

import (
	"errors"
	"log"
	"sync"

	"github.com/AGou-ops/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		// connLock:    sync.RWMutex{},
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，添加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到Connmanager中
	connMgr.connections[conn.GetConnID()] = conn
	log.Println(
		"Connection, connID",
		conn.GetConnID(),
		" add to ConnManager successfully: conn num = ",
		connMgr.Len(),
	)
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map，添加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	log.Println(
		"ConnManager del conn, connID",
		conn.GetConnID(),
		" successfully: conn num = ",
		connMgr.Len(),
	)
}

// 根据connID查找连接
func (connMgr *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	log.Println(
		"Clear All connections successfully, conn num = ",
		connMgr.Len(),
	)
}
