package pubsub

import (
	"testing"
	"time"
)

func TestNewPublisher(t *testing.T) {
	p := NewPublisher[int](1)
	if p == nil {
		t.Fatal("publisher is nil")
	}
}

func TestPublisher_Subscribe(t *testing.T) {
	p := NewPublisher[int](1)
	s := NewSubscriber[int](1)

	p.Subscribe(s)

	if len(p.subscribers) != 1 {
		t.Fatalf("expected 1, got %v", len(p.subscribers))
	}
	if p.subscribers[s] != s {
		t.Fatalf("expected %v, got %v", s, p.subscribers[s])
	}
}

func TestPublisher_Unsubscribe(t *testing.T) {
	p := NewPublisher[int](1)
	s := NewSubscriber[int](1)

	p.Subscribe(s)

	p.Unsubscribe(s)

	if len(p.subscribers) != 0 {
		t.Fatalf("expected 0, got %v", len(p.subscribers))
	}
	if p.subscribers[s] != nil {
		t.Fatalf("expected nil, got %v", p.subscribers[s])
	}
}

func TestPublisher_Publish(t *testing.T) {
	p := NewPublisher[int](1)
	s := NewSubscriber[int](1)

	p.Subscribe(s)

	err := p.Publish(1)
	if err != nil {
		t.Fatal(err)
	}

	v, err := s.Get(1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}
}

func TestPublisher_Publish_NotReady(t *testing.T) {
	p := NewPublisher[int](1)
	s := NewSubscriber[int](1)

	p.Subscribe(s)

	err := p.Publish(1)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Publish(2)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPublisher_Publish_MultipleSubscribers(t *testing.T) {
	p := NewPublisher[int](1)
	s1 := NewSubscriber[int](1)
	s2 := NewSubscriber[int](1)

	p.Subscribe(s1)
	p.Subscribe(s2)

	err := p.Publish(1)
	if err != nil {
		t.Fatal(err)
	}

	v1, err := s1.Get(1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if v1 != 1 {
		t.Fatalf("expected 1, got %v", v1)
	}

	v2, err := s2.Get(1 * time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if v2 != 1 {
		t.Fatalf("expected 1, got %v", v2)
	}
}

func TestPublisher_Publish_MultipleSubscribers_NotReady(t *testing.T) {
	p := NewPublisher[int](1)
	s1 := NewSubscriber[int](1)
	s2 := NewSubscriber[int](1)

	p.Subscribe(s1)
	p.Subscribe(s2)

	err := p.Publish(1)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Publish(2)
	if err == nil {
		t.Fatal("expected error")
	}
}
