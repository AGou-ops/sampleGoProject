package main

import (
	"log"
	"net/rpc"
)

type Result struct {
	Num, Ans int
}

func main() {
	client, _ := rpc.DialHTTP("tcp", "localhost:1234")
	var result Result
	// 同步调用
	if err := client.Call("Cal.Square", 14, &result); err != nil {
		log.Fatal("Failed to call Cal.Square", err)
	}

	log.Printf("%d^2 = %d", result.Num, result.Ans)

	// 异步调用
	// 因为 client.Go 是异步调用，因此第一次打印 result，result 没有被赋值。
	// 而通过调用 <-asyncCall.Done，阻塞当前程序直到 RPC 调用结束，因此第二次打印 result 时，能够看到正确的赋值。
	asyncCall := client.Go("Cal.Square", 12, &result, nil)
	log.Printf("%d^2 = %d", result.Num, result.Ans)

	<-asyncCall.Done
	log.Printf("%d^2 = %d", result.Num, result.Ans)
}
