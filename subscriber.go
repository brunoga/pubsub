package pubsub

import (
	"fmt"
	"time"
)

// Subscriber consume messages published by Publishers.
type Subscriber[T any] struct {
	ch chan T
}

// NewSubscriber creates a new Subscriber with a buffer of size.
func NewSubscriber[T any](size int) *Subscriber[T] {
	return &Subscriber[T]{
		ch: make(chan T, size),
	}
}

// Get returns the next message published by a Publisher that is directed to
// this Subscriber. If no message is available, Get will block until one is
// available or the timeout is reached.
func (s *Subscriber[T]) Get(timeout time.Duration) (T, error) {
	var zero T
	select {
	case v := <-s.ch:
		return v, nil
	case <-time.After(timeout):
		return zero, fmt.Errorf("timeout")
	}
}

// put is used by Publishers to send messages to this Subscriber.
func (s *Subscriber[T]) put(v T) error {
	select {
	case s.ch <- v:
		return nil
	default:
		return fmt.Errorf("not ready")
	}
}
