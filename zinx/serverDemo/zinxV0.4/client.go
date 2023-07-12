package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("client started")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		log.Println("err occurred: ", err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello world."))
		if err != nil {
			log.Println(err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Println("err occurred: ", err)
			return
		}
		fmt.Println("server back: ", string(buf), "\t", cnt)

		time.Sleep(time.Second)
	}
}
