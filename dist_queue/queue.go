package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type QueueService struct {
	queue []string
	mutex *sync.Mutex
}

type PutArgs struct {
	Item string
}

type TakeArgs struct {
}

type TakeReply struct {
	Item string
}

func (t *QueueService) Put(args *PutArgs, reply *bool) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.queue = append(t.queue, args.Item)
	*reply = true
	return nil
}

func (t *QueueService) Take(args *TakeArgs, reply *TakeReply) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if len(t.queue) == 0 {
		return errors.New("Queue is empty")
	}

	*reply = TakeReply{Item: t.queue[0]}
	t.queue = t.queue[1:]
	return nil
}

func main() {
	queueService := &QueueService{
		queue: []string{},
		mutex: &sync.Mutex{},
	}
	rpc.Register(queueService)
	rpc.HandleHTTP()

	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	go http.Serve(listener, nil)

	select {}
}