package main

import (
	"sync/atomic"
	"time"
	"unsafe"
	"fmt"
	"sync"
)

type Node struct {
	value interface{}
	next  *Node
}

type LockFreeStack struct {
	head 	*Node
	toFree	[]*Node
}

func NewLockFreeStack() *LockFreeStack {
	s := &LockFreeStack{
		toFree: make([]*Node, 0),
	}
	go s.cleanup()
	return s
}

func (s *LockFreeStack) cleanup() {
	for {
		time.Sleep(time.Second)
		toFree := s.toFree
		s.toFree = nil
		for _, node := range toFree {
			if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&node.next))) != nil {
				s.toFree = append(s.toFree, node)
				continue
			}
		}
	}
}

func (s *LockFreeStack) Push(v interface{}) {
	newNode := &Node{value: v}
	for {
		oldHead := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.head)))
		newNode.next = (*Node)(oldHead)
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&s.head)), oldHead, unsafe.Pointer(newNode)) {
			break
	}
	time.Sleep(time.Microsecond)
	}
	fmt.Printf("Pushed: %v\n", v)
}

func (s *LockFreeStack) Pop() (v interface{}, ok bool) {
	for {
		oldHead := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.head)))
		if oldHead == nil {
			return nil, false
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&s.head)), oldHead, unsafe.Pointer((*Node)(oldHead).next)) {
			v := (*Node)(oldHead).value
			s.toFree = append(s.toFree, (*Node)(oldHead))
			fmt.Printf("Popped: %v\n", v)
			return v, true
		}
		time.Sleep(time.Microsecond)
	}
}

func main() {
	// usage
	stack := NewLockFreeStack()
	n := 100

	var wg sync.WaitGroup
	wg.Add(2 * n)

	for i := 0; i < n; i++ {
		go func(i int) {
			stack.Push(i)
			wg.Done()
		}(i)
	}

	for i := 0; i < n; i++ {
		go func() {
			_, ok := stack.Pop()
			if ok {
				wg.Done()
			} else {
				fmt.Println("Error: popped from an empty stack")
			}
		}()
	}

	wg.Wait()
}
