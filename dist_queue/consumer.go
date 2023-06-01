package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type TakeArgs struct{}

type TakeReply struct {
	Item string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	for {
		var reply TakeReply
		err = client.Call("QueueService.Take", TakeArgs{}, &reply)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Received:", reply.Item)
		}

		time.Sleep(2 * time.Second)
	}
}
