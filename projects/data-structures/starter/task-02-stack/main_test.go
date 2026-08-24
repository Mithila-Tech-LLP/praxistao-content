package main

import "testing"

func TestStack_IsEmptyInitially(t *testing.T) {
	s := &Stack{}
	if !s.IsEmpty() {
		t.Error("new stack should be empty")
	}
}

func TestStack_SizeZero(t *testing.T) {
	s := &Stack{}
	if s.Size() != 0 {
		t.Errorf("new stack size: got %d, want 0", s.Size())
	}
}

func TestStack_PushAndSize(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if s.Size() != 3 {
		t.Errorf("size after 3 pushes: got %d, want 3", s.Size())
	}
	if s.IsEmpty() {
		t.Error("stack should not be empty after pushes")
	}
}

func TestStack_PopOrder(t *testing.T) {
	s := &Stack{}
	s.Push(10)
	s.Push(20)
	s.Push(30)

	val, ok := s.Pop()
	if !ok || val != 30 {
		t.Errorf("Pop 1: got (%d, %v), want (30, true)", val, ok)
	}
	val, ok = s.Pop()
	if !ok || val != 20 {
		t.Errorf("Pop 2: got (%d, %v), want (20, true)", val, ok)
	}
	val, ok = s.Pop()
	if !ok || val != 10 {
		t.Errorf("Pop 3: got (%d, %v), want (10, true)", val, ok)
	}
}

func TestStack_PopEmpty(t *testing.T) {
	s := &Stack{}
	val, ok := s.Pop()
	if ok {
		t.Error("Pop on empty: ok should be false")
	}
	if val != 0 {
		t.Errorf("Pop on empty: val should be 0, got %d", val)
	}
}

func TestStack_Peek(t *testing.T) {
	s := &Stack{}
	s.Push(5)
	s.Push(10)

	val, ok := s.Peek()
	if !ok || val != 10 {
		t.Errorf("Peek: got (%d, %v), want (10, true)", val, ok)
	}
	// Peek should not remove the element
	if s.Size() != 2 {
		t.Errorf("Size after Peek: got %d, want 2", s.Size())
	}
}

func TestStack_PeekEmpty(t *testing.T) {
	s := &Stack{}
	val, ok := s.Peek()
	if ok {
		t.Error("Peek on empty: ok should be false")
	}
	if val != 0 {
		t.Errorf("Peek on empty: val should be 0, got %d", val)
	}
}

func TestStack_SizeAfterPop(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Push(2)
	s.Pop()
	if s.Size() != 1 {
		t.Errorf("Size after pop: got %d, want 1", s.Size())
	}
}

func TestStack_IsEmptyAfterPopAll(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Pop()
	if !s.IsEmpty() {
		t.Error("stack should be empty after popping all elements")
	}
}

func TestStack_MultiplePopsBeyondEmpty(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Pop()
	_, ok1 := s.Pop()
	_, ok2 := s.Pop()
	if ok1 || ok2 {
		t.Error("popping beyond empty should return ok=false")
	}
}

func TestStack_SingleElement(t *testing.T) {
	s := &Stack{}
	s.Push(42)

	top, ok := s.Peek()
	if !ok || top != 42 {
		t.Errorf("Peek single: got (%d, %v), want (42, true)", top, ok)
	}

	val, ok := s.Pop()
	if !ok || val != 42 {
		t.Errorf("Pop single: got (%d, %v), want (42, true)", val, ok)
	}

	if !s.IsEmpty() {
		t.Error("should be empty after popping single element")
	}
}

func TestStack_PushNegatives(t *testing.T) {
	s := &Stack{}
	s.Push(-5)
	s.Push(-1)
	val, ok := s.Pop()
	if !ok || val != -1 {
		t.Errorf("Pop negative: got (%d, %v), want (-1, true)", val, ok)
	}
}
