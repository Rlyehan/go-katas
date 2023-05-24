package main

import (
	"fmt"
	"sync"
	"time"
)

type RingBuffer struct {
	mu 			*sync.Mutex
	notEmpty 	*sync.Cond
	notFull 	*sync.Cond
	buf 		[]interface{}
	indexes 	[2]int
	size 		int
	count 		int
}

func NewRingBuffer(size int) *RingBuffer {
	rb := &RingBuffer{
		buf: make([]interface{}, size),
		size:  size,
	}
	rb.mu = &sync.Mutex{}
	rb.notEmpty = sync.NewCond(rb.mu)
	rb.notFull = sync.NewCond(rb.mu)

	return rb
}

func (rb *RingBuffer) Enqueue(val interface{}) {
	rb.mu.Lock()

	for rb.count == rb.size {
		fmt.Println("Buffer full, waiting")
		rb.notFull.Wait()
	}

	rb.buf[rb.indexes[1]] = val
	rb.indexes[1] = (rb.indexes[1] + 1) % rb.size
	rb.count++

	rb.notEmpty.Signal()
	rb.mu.Unlock()
}

func (rb *RingBuffer) Dequeue() interface{} {
	rb.mu.Lock()

	for rb.count == 0 {
		fmt.Println("Buffer empty, waiting")
		rb.notEmpty.Wait()
	}

	val := rb.buf[rb.indexes[0]]
	rb.buf[rb.indexes[0]] = nil
	rb.indexes[0] = (rb.indexes[0] + 1) % rb.size
	rb.count--

	rb.notFull.Signal()
	rb.mu.Unlock()

	return val
}

func main() {
	rb := NewRingBuffer(10)

	go func() {
		for i := 0; i < 20; i++ {
			rb.Enqueue(i)
			fmt.Println("Enqueued", i)
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			val := rb.Dequeue()
			fmt.Println("Dequeued", val)
		}
	}()

	time.Sleep(time.Second * 2)
}
