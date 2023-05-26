package main

import (
	"fmt"
	"sync"
)

type Message struct {
	Body string
}

type Subscriber interface {
	Notify(Message)
}

type Broker struct {
	sync.RWMutex
	subscibers map[Subscriber]struct{}
}

func NewBroker() *Broker {
	return &Broker{subscibers: make(map[Subscriber]struct{})}
}

func (b *Broker) Subscribe(s Subscriber) {
	b.Lock()
	b.subscibers[s] = struct{}{}
	b.Unlock()
}

func (b *Broker) Unsubscribe(s Subscriber) {
	b.Lock()
	delete(b.subscibers, s)
	b.Unlock()
}

func (b *Broker) Publish(m Message) {
	b.RLock()
	defer b.RUnlock()
	for s := range b.subscibers {
		s.Notify(m)
	}
}

type TestSubsriber struct {
	ID string
}

func (t *TestSubsriber) Notify(m Message) {
	fmt.Printf("Subscriber %s received message: %s\n", t.ID, m.Body)
}

func main() {
	broker := NewBroker()

	s1 := &TestSubsriber{ID: "1"}
	s2 := &TestSubsriber{ID: "2"}

	broker.Subscribe(s1)
	broker.Subscribe(s2)

	broker.Publish(Message{Body: "Hello World!"})

	broker.Unsubscribe(s2)

	broker.Publish(Message{Body: "Hello Again!"})
}