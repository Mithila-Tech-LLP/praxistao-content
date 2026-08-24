package main

import "testing"

func TestMaxOfArray(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{1}, 1},
		{[]int{3, 1, 4, 1, 5, 9, 2, 6}, 9},
		{[]int{-3, -1, -4}, -1},
		{[]int{0, 0, 0}, 0},
		{[]int{100, 99, 98}, 100},
		{[]int{1, 2, 3, 4, 5}, 5},
	}
	for _, tt := range tests {
		got := MaxOfArray(tt.nums)
		if got != tt.want {
			t.Errorf("MaxOfArray(%v) = %d, want %d", tt.nums, got, tt.want)
		}
	}
}
