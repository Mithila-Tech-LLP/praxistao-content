package main

// Node is a single element in the linked list.
type Node struct {
	Val  int
	Next *Node
}

// LinkedList holds a reference to the head node.
type LinkedList struct {
	Head *Node
}

// Append adds val at the tail of the list.
func (l *LinkedList) Append(val int) {
	// TODO: implement
}

// Prepend adds val at the head of the list.
func (l *LinkedList) Prepend(val int) {
	// TODO: implement
}

// Delete removes the first occurrence of val.
// Does nothing if val is not found or the list is empty.
func (l *LinkedList) Delete(val int) {
	// TODO: implement
}

// Contains returns true if val exists in the list.
func (l *LinkedList) Contains(val int) bool {
	// TODO: implement
	return false
}

// ToSlice returns all values in order from head to tail.
func (l *LinkedList) ToSlice() []int {
	// TODO: implement
	return []int{}
}
