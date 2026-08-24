package main

type MinHeap struct {
	data []int
}

func (h *MinHeap) Push(val int) {
	// TODO: append and sift up
}

func (h *MinHeap) Pop() (int, bool) {
	// TODO: swap root with last, shrink, sift down; return (val, true) or (0, false) if empty
	return 0, false
}

func (h *MinHeap) Peek() (int, bool) {
	// TODO: return root without removing; (0, false) if empty
	return 0, false
}

func (h *MinHeap) Size() int {
	return len(h.data)
}
