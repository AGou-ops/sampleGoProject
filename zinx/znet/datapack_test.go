package znet

import (
	"io"
	"log"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:9899")
	defer listener.Close()
	if err != nil {
		log.Println("err occurred: ", err)
		return
	}

	go func() {
		for {
			conn, connErr := listener.Accept()
			if connErr != nil {
				log.Println("err occurred: ", err)
				return
			}
			go func(conn net.Conn) {
				// 处理客户端的请求
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadlen())
					_, readErr := io.ReadFull(conn, headData)
					if readErr != nil {
						log.Println("read head err: ", err)
						break
					}
					msgHead, unpackErr := dp.Unpack(headData)
					if unpackErr != nil {
						log.Println("server unpack err: ", unpackErr)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						if _, readErr := io.ReadFull(conn, msg.Data); readErr != nil {
							log.Println("server unpack err: ", err)
							return
						}
						log.Println(
							"Recv MsgID: ",
							msg.Id,
							", dataLen: ",
							msg.DataLen,
							", data: ",
							string(msg.Data),
						)
					}
				}
			}(conn)
		}
	}()
	conn, err := net.Dial("tcp", "127.0.0.1:9899")
	if err != nil {
		log.Println("err dial server", err)
		return
	}
	dp := NewDataPack()
	msg1 := &Message{
		Id:      0,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		log.Println("client pack msg1 error: ", err)
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		log.Println("client pack msg2 error: ", err)
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
