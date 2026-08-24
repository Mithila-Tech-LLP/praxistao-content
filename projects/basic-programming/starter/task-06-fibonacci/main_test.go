package main

import "testing"

func TestFibonacci(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 3},
		{5, 5}, {6, 8}, {7, 13}, {8, 21}, {9, 34}, {10, 55},
	}
	for _, tt := range tests {
		got := Fibonacci(tt.n)
		if got != tt.want {
			t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
