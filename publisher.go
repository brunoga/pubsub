package pubsub

import (
	"fmt"
	"sync"
)

// Publisher publishes messages to be consumed by Subscribers.
type Publisher[T any] struct {
	ch chan T

	m           sync.Mutex
	subscribers map[*Subscriber[T]]*Subscriber[T]
}

// NewPublisher creates a new Publisher with a buffer of size.
func NewPublisher[T any](size int) *Publisher[T] {
	p := &Publisher[T]{
		ch:          make(chan T, size),
		subscribers: make(map[*Subscriber[T]]*Subscriber[T]),
	}

	go p.loop()

	return p
}

// Subscribe adds a new Subscriber to this Publisher.
func (p *Publisher[T]) Subscribe(s *Subscriber[T]) {
	p.m.Lock()
	defer p.m.Unlock()

	p.subscribers[s] = s
}

// Unsubscribe removes a Subscriber from this Publisher.
func (p *Publisher[T]) Unsubscribe(s *Subscriber[T]) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, s)
}

// Publish sends a message to all Subscribers of this Publisher. If the internal
// channle buffer is full, publish will return an error and callers should
// retry (hopefully with a backoff).
func (p *Publisher[T]) Publish(v T) error {
	select {
	case p.ch <- v:
		return nil
	default:
		return fmt.Errorf("not ready")
	}
}

func (p *Publisher[T]) loop() {
	// Reads messages form the internal channel and sends them to all
	// subscribers.
	for v := range p.ch {
		// Locking here is fine because we are usually not going to have a lot
		// of subscribers and put will always return immediatelly. There is an
		// isseu with a Subscriber not being ready to process the message
		// immediatelly and, in this case, the message would be lost.
		p.m.Lock()
		for _, s := range p.subscribers {
			err := s.put(v)
			if err != nil {
				fmt.Printf("publisher loop error: %v\n", err)
			}
		}
		p.m.Unlock()
	}
}
