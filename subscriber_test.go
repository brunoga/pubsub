package pubsub

import "testing"

func TestNewSubscriber(t *testing.T) {
	s := NewSubscriber[int](1)
	if s == nil {
		t.Fatal("subscriber is nil")
	}
}

func TestSubscriber_Get_Timeout(t *testing.T) {
	s := NewSubscriber[int](1)

	_, err := s.Get(0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSubscriber_Get(t *testing.T) {
	s := NewSubscriber[int](1)

	err := s.put(1)
	if err != nil {
		t.Fatal(err)
	}

	v, err := s.Get(0)
	if err != nil {
		t.Fatal(err)
	}

	if v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}
}
