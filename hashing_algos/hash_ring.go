package main

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
)

const VirtualNodes = 10

type HashRing []uint32

func (h HashRing) Len() int {
	return len(h)
}

func (h HashRing) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h HashRing) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type Node struct {
	Id       string
	VirtualNodes HashRing
}

type ConsistentHash struct {
    Nodes map[uint32]*Node
    HashRing HashRing
}

func NewNode(id string) *Node {
	node := &Node{
		Id: id,
	}

	for i := 0; i < VirtualNodes; i++ {
		node.VirtualNodes = append(node.VirtualNodes, hashVal(id+strconv.Itoa(i)))
	}

	sort.Sort(node.VirtualNodes)

	return node
}

func NewConsistentHash() *ConsistentHash {
    return &ConsistentHash{
        Nodes: make(map[uint32]*Node),
    }
}

func (ch *ConsistentHash) AddNode(node *Node) {
    for _, vNode := range node.VirtualNodes {
        ch.Nodes[vNode] = node
        ch.HashRing = append(ch.HashRing, vNode)
    }
    sort.Sort(ch.HashRing)
}

func (ch *ConsistentHash) RemoveNode(node *Node) {
    for _, vNode := range node.VirtualNodes {
        delete(ch.Nodes, vNode)
        index := sort.Search(len(ch.HashRing), func(i int) bool {
            return ch.HashRing[i] >= vNode
        })
        ch.HashRing = append(ch.HashRing[:index], ch.HashRing[index+1:]...)
    }
}

func (ch *ConsistentHash) GetNode(key string) *Node {
	hash := hashVal(key)
	idx := sort.Search(len(ch.HashRing), func(i int) bool { return ch.HashRing[i] >= hash })
	if idx == len(ch.HashRing) {
		idx = 0
	}
	return ch.Nodes[ch.HashRing[idx]]
}

func hashVal(key string) uint32 {
	h := sha256.New()
	h.Write([]byte(key))
	hash := h.Sum(nil)

	return uint32(hash[0])<<24 + uint32(hash[1])<<16 + uint32(hash[2])<<8 + uint32(hash[3])
}

func main() {
	nodes := []*Node{
		NewNode("node1"),
		NewNode("node2"),
		NewNode("node3"),
	}

    ch := NewConsistentHash()

	for _, node := range nodes {
		ch.AddNode(node)
	}

	fmt.Println(ch.GetNode("alpha").Id)  
	fmt.Println(ch.GetNode("beta").Id) 
	fmt.Println(ch.GetNode("gamma").Id) 

    node4 := NewNode("node4")
	ch.AddNode(node4)
	fmt.Println(ch.GetNode("someKey").Id)

    ch.RemoveNode(nodes[0])
	fmt.Println(ch.GetNode("someKey").Id)
}
