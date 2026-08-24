package main

import "testing"

func TestQueue_IsEmptyInitially(t *testing.T) {
	q := &Queue{}
	if !q.IsEmpty() {
		t.Error("new queue should be empty")
	}
}

func TestQueue_SizeZero(t *testing.T) {
	q := &Queue{}
	if q.Size() != 0 {
		t.Errorf("new queue size: got %d, want 0", q.Size())
	}
}

func TestQueue_EnqueueAndSize(t *testing.T) {
	q := &Queue{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Size() != 3 {
		t.Errorf("size after 3 enqueues: got %d, want 3", q.Size())
	}
	if q.IsEmpty() {
		t.Error("queue should not be empty after enqueues")
	}
}

func TestQueue_DequeueOrder(t *testing.T) {
	q := &Queue{}
	q.Enqueue(10)
	q.Enqueue(20)
	q.Enqueue(30)

	val, ok := q.Dequeue()
	if !ok || val != 10 {
		t.Errorf("Dequeue 1: got (%d, %v), want (10, true)", val, ok)
	}
	val, ok = q.Dequeue()
	if !ok || val != 20 {
		t.Errorf("Dequeue 2: got (%d, %v), want (20, true)", val, ok)
	}
	val, ok = q.Dequeue()
	if !ok || val != 30 {
		t.Errorf("Dequeue 3: got (%d, %v), want (30, true)", val, ok)
	}
}

func TestQueue_DequeueEmpty(t *testing.T) {
	q := &Queue{}
	val, ok := q.Dequeue()
	if ok {
		t.Error("Dequeue on empty: ok should be false")
	}
	if val != 0 {
		t.Errorf("Dequeue on empty: val should be 0, got %d", val)
	}
}

func TestQueue_Peek(t *testing.T) {
	q := &Queue{}
	q.Enqueue(5)
	q.Enqueue(10)

	val, ok := q.Peek()
	if !ok || val != 5 {
		t.Errorf("Peek: got (%d, %v), want (5, true)", val, ok)
	}
	// Peek should not remove element
	if q.Size() != 2 {
		t.Errorf("Size after Peek: got %d, want 2", q.Size())
	}
}

func TestQueue_PeekEmpty(t *testing.T) {
	q := &Queue{}
	val, ok := q.Peek()
	if ok {
		t.Error("Peek on empty: ok should be false")
	}
	if val != 0 {
		t.Errorf("Peek on empty: val should be 0, got %d", val)
	}
}

func TestQueue_IsEmptyAfterDequeueAll(t *testing.T) {
	q := &Queue{}
	q.Enqueue(1)
	q.Dequeue()
	if !q.IsEmpty() {
		t.Error("queue should be empty after dequeuing all elements")
	}
}

func TestQueue_SizeDecrementsOnDequeue(t *testing.T) {
	q := &Queue{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Dequeue()
	if q.Size() != 1 {
		t.Errorf("Size after dequeue: got %d, want 1", q.Size())
	}
}

func TestQueue_MultipleDequeuesBeyondEmpty(t *testing.T) {
	q := &Queue{}
	q.Enqueue(1)
	q.Dequeue()
	_, ok1 := q.Dequeue()
	_, ok2 := q.Dequeue()
	if ok1 || ok2 {
		t.Error("dequeuing beyond empty should return ok=false")
	}
}

func TestQueue_InterleavedEnqueueDequeue(t *testing.T) {
	q := &Queue{}
	q.Enqueue(1)
	q.Enqueue(2)

	val, _ := q.Dequeue() // 1
	if val != 1 {
		t.Errorf("interleaved: got %d, want 1", val)
	}

	q.Enqueue(3)
	val, _ = q.Dequeue() // 2
	if val != 2 {
		t.Errorf("interleaved: got %d, want 2", val)
	}
	val, _ = q.Dequeue() // 3
	if val != 3 {
		t.Errorf("interleaved: got %d, want 3", val)
	}
}

func TestQueue_SingleElement(t *testing.T) {
	q := &Queue{}
	q.Enqueue(99)

	front, ok := q.Peek()
	if !ok || front != 99 {
		t.Errorf("Peek single: got (%d, %v), want (99, true)", front, ok)
	}

	val, ok := q.Dequeue()
	if !ok || val != 99 {
		t.Errorf("Dequeue single: got (%d, %v), want (99, true)", val, ok)
	}
	if !q.IsEmpty() {
		t.Error("should be empty after dequeuing single element")
	}
}
