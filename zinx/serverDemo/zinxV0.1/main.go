package main

import "github.com/AGou-ops/zinx/znet"

func main() {
	s := znet.NewServer("zinx v0.1")
	s.Server()
}
