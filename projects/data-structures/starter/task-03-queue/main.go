package main

// Queue is a First-In, First-Out (FIFO) collection backed by a slice.
type Queue struct {
	data []int
}

// Enqueue adds val to the back of the queue.
func (q *Queue) Enqueue(val int) {
	// TODO: implement
}

// Dequeue removes and returns the front element.
// Returns (0, false) if the queue is empty.
func (q *Queue) Dequeue() (int, bool) {
	// TODO: implement
	return 0, false
}

// Peek returns the front element without removing it.
// Returns (0, false) if the queue is empty.
func (q *Queue) Peek() (int, bool) {
	// TODO: implement
	return 0, false
}

// IsEmpty returns true if the queue has no elements.
func (q *Queue) IsEmpty() bool {
	// TODO: implement
	return true
}

// Size returns the number of elements in the queue.
func (q *Queue) Size() int {
	// TODO: implement
	return 0
}
