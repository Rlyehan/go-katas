package main

import (
	"fmt"
)

type RingBuffer struct {
	buf		[]interface{}
	indexes [2]int
	size	int
	count	int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buf: make([]interface{}, size),
		size: size,
	}
}

func (rb *RingBuffer) Enqueue(val interface{}) {
	if rb.count == rb.size {
		fmt.Println("Buffer is full")
		return
	}

	rb.buf[rb.indexes[1]] = val
	rb.indexes[1] = (rb.indexes[1] + 1) % rb.size
	rb.count++
}

func (rb *RingBuffer) Dequeue() interface{} {
	if rb.count == 0 {
		fmt.Println("Buffer is empty")
		return nil
	}

	val := rb.buf[rb.indexes[0]]
	rb.buf[rb.indexes[0]] = nil
	rb.indexes[0] = (rb.indexes[0] + 1) % rb.size
	rb.count--

	return val
}

func main() {
	rb := NewRingBuffer(10)

	for i := 0; i < 12; i++ {
		rb.Enqueue(i)
	}

	for i := 0; i < 12; i++ {
		fmt.Println(rb.Dequeue())
	}
}