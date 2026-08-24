package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	palindromes := []string{"racecar", "abba", "a", "", "madam", "level"}
	for _, s := range palindromes {
		if !IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = false, want true", s)
		}
	}
	notPalindromes := []string{"hello", "world", "abcd", "Racecar"}
	for _, s := range notPalindromes {
		if IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = true, want false", s)
		}
	}
}
