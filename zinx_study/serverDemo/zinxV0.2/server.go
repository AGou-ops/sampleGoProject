package main

import "github.com/AGou-ops/zinx/znet"

func main() {
	s := znet.NewServer("zinx v0.2")
	s.Server()
}
