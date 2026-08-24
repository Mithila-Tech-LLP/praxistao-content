package main

import "testing"

func TestEvenOrOdd(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "even"},
		{1, "odd"},
		{2, "even"},
		{3, "odd"},
		{4, "even"},
		{-4, "even"},
		{-7, "odd"},
		{100, "even"},
		{101, "odd"},
	}
	for _, tt := range tests {
		got := EvenOrOdd(tt.n)
		if got != tt.want {
			t.Errorf("EvenOrOdd(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
