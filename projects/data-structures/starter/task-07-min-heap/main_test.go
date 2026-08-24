package main

import "testing"

func TestMinHeap_PopEmpty(t *testing.T) {
	h := &MinHeap{}
	_, ok := h.Pop()
	if ok {
		t.Error("Pop on empty heap should return false")
	}
}

func TestMinHeap_PeekEmpty(t *testing.T) {
	h := &MinHeap{}
	_, ok := h.Peek()
	if ok {
		t.Error("Peek on empty heap should return false")
	}
}

func TestMinHeap_SingleElement(t *testing.T) {
	h := &MinHeap{}
	h.Push(7)

	if h.Size() != 1 {
		t.Errorf("Size after one push = %d, want 1", h.Size())
	}

	peek, ok := h.Peek()
	if !ok || peek != 7 {
		t.Errorf("Peek = (%d, %v), want (7, true)", peek, ok)
	}

	val, ok := h.Pop()
	if !ok || val != 7 {
		t.Errorf("Pop = (%d, %v), want (7, true)", val, ok)
	}

	if h.Size() != 0 {
		t.Errorf("Size after popping last element = %d, want 0", h.Size())
	}
}

func TestMinHeap_PopReturnsAscendingOrder(t *testing.T) {
	h := &MinHeap{}
	input := []int{9, 3, 7, 1, 5, 8, 2, 6, 4}
	for _, v := range input {
		h.Push(v)
	}

	if h.Size() != len(input) {
		t.Errorf("Size = %d, want %d", h.Size(), len(input))
	}

	prev := -1
	for i := 1; i <= len(input); i++ {
		val, ok := h.Pop()
		if !ok {
			t.Fatalf("Pop returned false on iteration %d", i)
		}
		if val <= prev {
			t.Errorf("Pop iteration %d: got %d which is not > previous %d (not ascending)", i, val, prev)
		}
		prev = val
	}

	if h.Size() != 0 {
		t.Errorf("Size after all pops = %d, want 0", h.Size())
	}
}

func TestMinHeap_PeekDoesNotRemove(t *testing.T) {
	h := &MinHeap{}
	h.Push(3)
	h.Push(1)
	h.Push(2)

	for i := 0; i < 3; i++ {
		val, ok := h.Peek()
		if !ok || val != 1 {
			t.Errorf("Peek iteration %d = (%d, %v), want (1, true)", i, val, ok)
		}
	}

	if h.Size() != 3 {
		t.Errorf("Size after 3 Peeks = %d, want 3 (Peek should not remove)", h.Size())
	}
}

func TestMinHeap_PeekEqualsFirstPop(t *testing.T) {
	h := &MinHeap{}
	for _, v := range []int{10, 4, 6, 1, 8} {
		h.Push(v)
	}

	peek, _ := h.Peek()
	pop, _ := h.Pop()
	if peek != pop {
		t.Errorf("Peek (%d) should equal first Pop (%d)", peek, pop)
	}
}

func TestMinHeap_SizeTracking(t *testing.T) {
	h := &MinHeap{}
	for i := 1; i <= 5; i++ {
		h.Push(i)
		if h.Size() != i {
			t.Errorf("Size after %d pushes = %d, want %d", i, h.Size(), i)
		}
	}
	for i := 4; i >= 0; i-- {
		h.Pop()
		if h.Size() != i {
			t.Errorf("Size after pop = %d, want %d", h.Size(), i)
		}
	}
}

func TestMinHeap_DuplicateValues(t *testing.T) {
	h := &MinHeap{}
	for _, v := range []int{3, 1, 2, 1, 3, 2} {
		h.Push(v)
	}

	expected := []int{1, 1, 2, 2, 3, 3}
	for _, want := range expected {
		got, ok := h.Pop()
		if !ok || got != want {
			t.Errorf("Pop = (%d, %v), want (%d, true)", got, ok, want)
		}
	}
}
