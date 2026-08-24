package main

import "testing"

func TestCountCharacters(t *testing.T) {
	t.Run("hello", func(t *testing.T) {
		got := CountCharacters("hello")
		want := map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}
		for ch, wc := range want {
			if got[ch] != wc {
				t.Errorf("CountCharacters(%q)[%q] = %d, want %d", "hello", string(ch), got[ch], wc)
			}
		}
	})
	t.Run("empty", func(t *testing.T) {
		got := CountCharacters("")
		if len(got) != 0 {
			t.Errorf("CountCharacters(%q) should be empty, got %v", "", got)
		}
	})
	t.Run("single", func(t *testing.T) {
		got := CountCharacters("aaa")
		if got['a'] != 3 {
			t.Errorf("CountCharacters(%q)['a'] = %d, want 3", "aaa", got['a'])
		}
	})
}
