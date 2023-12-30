package pubsub

import (
	"fmt"
	"sync"
)

type Publisher[T any] struct {
	ch chan T

	m           sync.Mutex
	subscribers map[*Subscriber[T]]*Subscriber[T]
}

func NewPublisher[T any](size int) *Publisher[T] {
	p := &Publisher[T]{
		ch:          make(chan T, size),
		subscribers: make(map[*Subscriber[T]]*Subscriber[T]),
	}

	go p.loop()

	return p
}

func (p *Publisher[T]) Subscribe(s *Subscriber[T]) {
	p.m.Lock()
	defer p.m.Unlock()

	p.subscribers[s] = s
}

func (p *Publisher[T]) Unsubscribe(s *Subscriber[T]) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, s)
}

func (p *Publisher[T]) Publish(v T) error {
	select {
	case p.ch <- v:
		return nil
	default:
		return fmt.Errorf("not ready")
	}
}

func (p *Publisher[T]) loop() {
	for v := range p.ch {
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
