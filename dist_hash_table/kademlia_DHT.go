package main

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

type NodeID [20]byte

type Node struct {
	NodeID NodeID
	Data   map[string]string
}

func NewNode() *Node {
	rand.Seed(time.Now().UnixNano())
	node := &Node{Data: make(map[string]string)}
	rand.Read(node.NodeID[:]) // Randomly generate NodeID.
	return node
}

func (node *Node) CloserThan(other *Node, keyHash []byte) bool {
	return closer(node.NodeID[:], other.NodeID[:], keyHash)
}

func closer(a, b, key []byte) bool {
	distA := xor(a, key)
	distB := xor(b, key)
	return big.NewInt(0).SetBytes(distA).Cmp(big.NewInt(0).SetBytes(distB)) < 0
}

func xor(a, b []byte) []byte {
	l := len(a)
	if len(b) < l {
		l = len(b)
	}
	res := make([]byte, l)
	for i := 0; i < l; i++ {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func main() {
	nodes := []*Node{
		NewNode(),
		NewNode(),
		NewNode(),
		NewNode(),
	}

	h := sha1.New()
	h.Write([]byte("Hello"))
	keyHash := h.Sum(nil)

	var closestNode *Node
	for _, node := range nodes {
		if closestNode == nil || node.CloserThan(closestNode, keyHash) {
			closestNode = node
		}
	}

	closestNode.Data["Hello"] = "World"
	for _, node := range nodes {
		fmt.Printf("NodeID: %x Data: %v\n", node.NodeID, node.Data)
	}
}
