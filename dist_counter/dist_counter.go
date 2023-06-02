package main

import (
	"fmt"
	"sync"
)

type Node struct {
	counter int
	leader bool
	mutex sync.RWMutex
}

func NewNode() *Node {
	return &Node {
		counter: 0,
		leader: false,
	}
}

func (n *Node) Increment() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.counter++
}

func (n *Node) Decrement() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.counter--
}

func (n *Node) GetCounter() int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.counter
}

func main() {
	nodes := []*Node{}
	for i := 0; i < 5; i++ {
		nodes = append(nodes, NewNode())
	}

	nodes[0].leader = true
	fmt.Println("Node 0 selected as leader")

	nodes[0].Increment()
	nodes[0].Increment()
	nodes[0].Decrement()

	for i, node := range nodes[1:] {
		node.counter = nodes[0].GetCounter()
		fmt.Printf("Synced Node %d with leader's counter value", i+1)
	}

	for i, node := range nodes {
		fmt.Printf("Node %d counter: %d\n", i, node.GetCounter())
	}

	nodes[0].leader = false
	nodes[1].leader = true
	fmt.Println("Node 1 selected as leader")

	nodes[1].Increment()

	for i, node := range nodes {
		if !node.leader {
			node.counter = nodes[1].GetCounter()
			fmt.Printf("Synced Node %d with new leader's counter value", i)
		}
	}

	for i, node := range nodes {
		fmt.Printf("Node %d counter: %d\n", i, node.GetCounter())
	}
}