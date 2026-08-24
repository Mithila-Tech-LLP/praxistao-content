package main

// Stack is a Last-In, First-Out (LIFO) collection backed by a slice.
type Stack struct {
	data []int
}

// Push adds val to the top of the stack.
func (s *Stack) Push(val int) {
	// TODO: implement
}

// Pop removes and returns the top element.
// Returns (0, false) if the stack is empty.
func (s *Stack) Pop() (int, bool) {
	// TODO: implement
	return 0, false
}

// Peek returns the top element without removing it.
// Returns (0, false) if the stack is empty.
func (s *Stack) Peek() (int, bool) {
	// TODO: implement
	return 0, false
}

// IsEmpty returns true if the stack has no elements.
func (s *Stack) IsEmpty() bool {
	// TODO: implement
	return true
}

// Size returns the number of elements in the stack.
func (s *Stack) Size() int {
	// TODO: implement
	return 0
}
