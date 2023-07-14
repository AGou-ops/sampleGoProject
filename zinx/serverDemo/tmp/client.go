package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	defer conn.Close()
	if err != nil {
		log.Println("cannot dial tcp server: ", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("input something here: ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Read from stdin failed: ", err)
		}
		conn.Write([]byte(msg))

		// writer := bufio.NewWriter(os.Stdout)
		// if err != nil {
		// 	log.Println("Write to stdout error: ", err)
		// }
		// buf := make([]byte, 1024)
		// conn.Read(buf)
		// writer.WriteString(string(buf))
		// 读取服务器的响应
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Printf("Server response: %s", response)
	}
}
