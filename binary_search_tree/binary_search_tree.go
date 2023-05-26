package main

import (
	"fmt"
)

type Node struct {
	Key int
	Left *Node
	Right *Node
}

type Queue []*Node

func (q *Queue) Enqueue(n *Node) {
	*q = append(*q, n)
}

func (q *Queue) Dequeue() (*Node, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	node := (*q)[0]
	*q = (*q)[1:]
	return node, true
}

func (n *Node) Insert(key int) *Node {
	if n==nil {
		return &Node{Key: key}
	}
	if key < n.Key {
		n.Left = n.Left.Insert(key)
	}
	if key > n.Key {
		n.Right = n.Right.Insert(key)
	}
	return n
}

func (n *Node) exists(key int) bool {
	if n==nil {
		return false
	}
	if key == n.Key {
		return true
	}
	if key < n.Key {
		return n.Left.exists(key)
	} else {
		return n.Right.exists(key)
	}
}

func (n *Node) InOrderTraversal() {
	if n != nil {
		n.Left.InOrderTraversal()
		fmt.Println(n.Key)
		n.Right.InOrderTraversal()
	}
}

func (n *Node) PreOrderTraversal() {
	if n != nil {
		fmt.Println(n.Key)
		n.Left.PreOrderTraversal()
		n.Right.PreOrderTraversal()
	}
}

func (n *Node) PostOrderTraversal() {
	if n != nil {
		n.Left.PostOrderTraversal()
		n.Right.PostOrderTraversal()
		fmt.Println(n.Key)
	}
}

func (n *Node) BreadthFirstTraversal() {
	var q Queue
	q.Enqueue(n)
	for len(q) > 0 {
		node, _ := q.Dequeue()
		fmt.Println(node.Key)
		if node.Left != nil {
			q.Enqueue(node.Left)
		}
		if node.Right != nil {
			q.Enqueue(node.Right)
		}
	}
}

func (n *Node) isEqualTo(m *Node) bool {
	if n == nil && m == nil {
		return true
	}
	if n == nil || m == nil {
		return false
	}
	return n.Key == m.Key && n.Left.isEqualTo(m.Left) && n.Right.isEqualTo(m.Right)
}

func main() {
	root := &Node{Key: 10}
	root = root.Insert(5)
	root = root.Insert(15)
	root = root.Insert(3)
	root = root.Insert(7)
	root = root.Insert(13)
	root = root.Insert(17)

	fmt.Println("In-Order Traversal:")
	root.InOrderTraversal()

	fmt.Println("Pre-Order Traversal:")
	root.PreOrderTraversal()

	fmt.Println("Post-Order Traversal:")
	root.PostOrderTraversal()

	fmt.Println("Breadth-First Traversal:")
	root.BreadthFirstTraversal()

	fmt.Println("Does 7 exist in the tree?")
	fmt.Println(root.exists(7))

	fmt.Println("Does 20 exist in the tree?")
	fmt.Println(root.exists(20))

	root2 := &Node{Key: 10}
	root2 = root2.Insert(5)
	root2 = root2.Insert(15)
	root2 = root2.Insert(3)
	root2 = root2.Insert(7)
	root2 = root2.Insert(13)
	root2 = root2.Insert(17)

	fmt.Println("Does the two trees equal?")
	fmt.Println(root.isEqualTo(root2))
}