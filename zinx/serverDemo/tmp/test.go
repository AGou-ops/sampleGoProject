package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const headerSize = 4 // 头部长度的字节数

type Cat struct {
	test string
}

func main() {
	// 启动服务器
	go startServer()

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	// 发送消息
	message := "Hello, Server!"
	sendMessage(conn, message)

	// 读取服务器响应
	response, err := readMessage(conn)
	if err != nil {
		fmt.Println("读取消息失败:", err)
		return
	}
	fmt.Println("服务器响应:", response)
}

func startServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动服务器失败:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务器已启动，等待连接...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("客户端 %s 已连接\n", conn.RemoteAddr().String())

	defer conn.Close()

	// 读取消息
	message, err := readMessage(conn)
	if err != nil {
		fmt.Println("读取消息失败:", err)
		return
	}
	fmt.Println("收到消息:", message)

	// 发送响应
	response := "Hello, Client!"
	sendMessage(conn, response)
}

func sendMessage(conn net.Conn, message string) error {
	// 计算消息长度
	messageLength := len(message)

	// 将消息长度写入头部
	header := make([]byte, headerSize)
	binary.BigEndian.PutUint32(header, uint32(messageLength))
	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("写入消息头部失败: %v", err)
	}

	// 写入消息体
	if _, err := conn.Write([]byte(message)); err != nil {
		return fmt.Errorf("写入消息体失败: %v", err)
	}

	return nil
}

func readMessage(conn net.Conn) (string, error) {
	// 读取消息头部
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(conn, header); err != nil {
		return "", fmt.Errorf("读取消息头部失败: %v", err)
	}

	// 解析消息长度
	messageLength := binary.BigEndian.Uint32(header)

	// 读取消息体
	message := make([]byte, messageLength)
	if _, err := io.ReadFull(conn, message); err != nil {
		return "", fmt.Errorf("读取消息体失败: %v", err)
	}

	return string(message), nil
}
