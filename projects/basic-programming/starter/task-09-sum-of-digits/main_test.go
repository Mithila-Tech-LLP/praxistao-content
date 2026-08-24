package main

import "testing"

func TestSumOfDigits(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0},
		{1, 1},
		{9, 9},
		{10, 1},
		{123, 6},
		{456, 15},
		{999, 27},
		{-123, 6},
		{-456, 15},
	}
	for _, tt := range tests {
		got := SumOfDigits(tt.n)
		if got != tt.want {
			t.Errorf("SumOfDigits(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
