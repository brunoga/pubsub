package pubsub

import (
	"fmt"
	"time"
)

type Subscriber[T any] struct {
	ch chan T
}

func NewSubscriber[T any](size int) *Subscriber[T] {
	return &Subscriber[T]{
		ch: make(chan T, size),
	}
}

func (s *Subscriber[T]) Get(timeout time.Duration) (T, error) {
	var zero T
	select {
	case v := <-s.ch:
		return v, nil
	case <-time.After(timeout):
		return zero, fmt.Errorf("timeout")
	}
}

func (s *Subscriber[T]) put(v T) error {
	select {
	case s.ch <- v:
		return nil
	default:
		return fmt.Errorf("not ready")
	}
}
