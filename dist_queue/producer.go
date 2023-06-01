package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type PutArgs struct {
	Item string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}

	count := 1
	for {
		item := fmt.Sprintf("Message %d", count)
		var reply bool
		err = client.Call("QueueService.Put", PutArgs{Item: item}, &reply)
		if err != nil {
			log.Fatal("Error:", err)
		}
		fmt.Printf("Produced: %s\n", item)

		time.Sleep(1 * time.Second)
		count++
	}
}
