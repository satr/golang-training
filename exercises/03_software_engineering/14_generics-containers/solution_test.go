package genericscontainers

import (
	"testing"
)

func TestStackPushPop(t *testing.T) {
	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	v, err := s.Pop()
	if err != nil || v != 2 {
		t.Errorf("want 2, got %v %v", v, err)
	}
	v, err = s.Pop()
	if err != nil || v != 1 {
		t.Errorf("want 1, got %v %v", v, err)
	}
}

func TestStackEmpty(t *testing.T) {
	s := NewStack[int]()
	if !s.IsEmpty() {
		t.Error("expected empty")
	}
	_, err := s.Pop()
	if err == nil {
		t.Error("expected error on pop empty")
	}
}

func TestStackPeek(t *testing.T) {
	s := NewStack[string]()
	s.Push("hello")
	v, err := s.Peek()
	if err != nil || v != "hello" {
		t.Errorf("want hello, got %v %v", v, err)
	}
	if s.Size() != 1 {
		t.Error("peek should not change size")
	}
}

func TestMapSetGet(t *testing.T) {
	m := NewMap[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	v, ok := m.Get("one")
	if !ok || v != 1 {
		t.Errorf("want 1, got %v %v", v, ok)
	}
	if m.Len() != 2 {
		t.Errorf("len = %d, want 2", m.Len())
	}
}

func TestMapMissing(t *testing.T) {
	m := NewMap[string, int]()
	_, ok := m.Get("nope")
	if ok {
		t.Error("expected missing key")
	}
}
