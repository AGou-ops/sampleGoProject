package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/AGou-ops/zinx/znet"
)

func main() {
	log.Println("client1 started")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		log.Println("err occurred: ", err)
		return
	}

	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(
			znet.NewMsgPackage(1, []byte("zinx0.6 client1 test msg")),
		)
		if err != nil {
			log.Println("Pack msg error: ", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			log.Println("write to conn error: ", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadlen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			log.Println("read head error ", err)
			break
		}
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			log.Println("client unpack msg head error: ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				log.Println("read msg data error: ", err)
				break
			}
			log.Printf(
				"--> Recv serverID: %d, Data: %s, dataLen: %d",
				msg.Id,
				msg.Data,
				msg.DataLen,
			)
		}

		time.Sleep(time.Second)
	}
}
