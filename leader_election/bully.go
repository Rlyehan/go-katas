package main

import (
	"fmt"
	"time"
	"sync"
)

type Node struct {
	id int
	up bool
	leader bool
	mu sync.Mutex
}

type Message struct {
	source *Node
	kind string
}

func startElection(node *Node, nodes []*Node, messageChannels []chan Message) {
	if node.up == false {
		return
	}
	fmt.Println("Starting election:", node.id)
	highestID := node.id
	for _, peer := range nodes {
		if peer.id > node.id && peer.up {
			go startElection(peer, nodes, messageChannels)
			messageChannels[peer.id] <- Message{source: node, kind: "election"}
			message := <-messageChannels[peer.id]
			if message.kind == "answer" && message.source.id > highestID {
				highestID = message.source.id
			}
		}
	}
	if highestID == node.id {
		node.mu.Lock()
		node.leader = true
		node.mu.Unlock()
		fmt.Println("Election finished, new leader:", node.id)
		for _, peer := range nodes {
			if peer.id < node.id && peer.up {
				messageChannels[peer.id] <-Message{source: node, kind: "victory"}
				peer.mu.Lock()
				peer.leader = false
				peer.mu.Unlock()
			}
		}
	} else {
		node.mu.Lock()
		node.leader = false
		node.mu.Unlock()
	}
}

func main() {
	nodes := []*Node{}
	for i := 0; i < 5; i++ {
		nodes = append(nodes, &Node{id: i, up: true, leader: false})
	}
	nodes[4].leader = true

	messageChannels := make([]chan Message, len(nodes))
	for i := range messageChannels {
		messageChannels[i] = make(chan Message)
	}

	go startElection(nodes[0], nodes, messageChannels)
	time.Sleep(time.Second * 5)

	nodes[4].mu.Lock()
	nodes[4].up = false
	nodes[4].mu.Unlock()

	fmt.Println("Node 4 down")

	go startElection(nodes[0], nodes, messageChannels)

	time.Sleep(time.Second * 5)
}