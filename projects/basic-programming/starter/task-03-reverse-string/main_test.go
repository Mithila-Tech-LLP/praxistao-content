package main

import "testing"

func TestReverseString(t *testing.T) {
	tests := []struct{ in, want string }{
		{"hello", "olleh"},
		{"Go", "oG"},
		{"", ""},
		{"a", "a"},
		{"racecar", "racecar"},
		{"abcdef", "fedcba"},
	}
	for _, tt := range tests {
		got := ReverseString(tt.in)
		if got != tt.want {
			t.Errorf("ReverseString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
