package main

type BSTNode struct {
	Val   int
	Left  *BSTNode
	Right *BSTNode
}

type BST struct {
	Root *BSTNode
}

func (t *BST) Insert(val int) {
	// TODO: implement
}

func (t *BST) Search(val int) bool {
	// TODO: implement
	return false
}

func (t *BST) InOrder() []int {
	// TODO: implement — return all values in ascending order
	return []int{}
}

func (t *BST) Min() (int, bool) {
	// TODO: implement — return minimum value, false if tree is empty
	return 0, false
}

func (t *BST) Max() (int, bool) {
	// TODO: implement — return maximum value, false if tree is empty
	return 0, false
}
