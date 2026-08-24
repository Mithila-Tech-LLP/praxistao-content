package main

import (
	"reflect"
	"testing"
)

func TestTrie_Search_AfterInsert(t *testing.T) {
	tr := NewTrie()
	words := []string{"apple", "app", "application", "apply", "apt", "banana"}
	for _, w := range words {
		tr.Insert(w)
	}

	for _, w := range words {
		if !tr.Search(w) {
			t.Errorf("Search(%q) should return true after inserting it", w)
		}
	}
}

func TestTrie_Search_NotInserted(t *testing.T) {
	tr := NewTrie()
	tr.Insert("apple")
	tr.Insert("app")

	notInserted := []string{"ap", "appl", "apples", "application", "banana"}
	for _, w := range notInserted {
		if tr.Search(w) {
			t.Errorf("Search(%q) should return false — not inserted", w)
		}
	}
}

func TestTrie_Search_EmptyTrie(t *testing.T) {
	tr := NewTrie()
	if tr.Search("anything") {
		t.Error("Search on empty trie should return false")
	}
}

func TestTrie_StartsWith(t *testing.T) {
	tr := NewTrie()
	for _, w := range []string{"apple", "app", "application", "apt", "banana"} {
		tr.Insert(w)
	}

	cases := []struct {
		prefix string
		want   bool
	}{
		{"app", true},
		{"appl", true},
		{"ap", true},
		{"a", true},
		{"ban", true},
		{"b", true},
		{"c", false},
		{"apz", false},
		{"applications", false},
	}
	for _, tc := range cases {
		got := tr.StartsWith(tc.prefix)
		if got != tc.want {
			t.Errorf("StartsWith(%q) = %v, want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestTrie_StartsWith_EmptyPrefix(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")
	// empty prefix matches everything
	if !tr.StartsWith("") {
		t.Error("StartsWith(\"\") should return true when trie has words")
	}
}

func TestTrie_WordsWithPrefix(t *testing.T) {
	tr := NewTrie()
	for _, w := range []string{"apple", "app", "application", "apply", "apt", "banana", "band"} {
		tr.Insert(w)
	}

	cases := []struct {
		prefix string
		want   []string
	}{
		{"app", []string{"app", "apple", "application", "apply"}},
		{"apt", []string{"apt"}},
		{"ban", []string{"banana", "band"}},
		{"xyz", []string{}},
	}

	for _, tc := range cases {
		got := tr.WordsWithPrefix(tc.prefix)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("WordsWithPrefix(%q) = %v, want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestTrie_WordsWithPrefix_EmptyPrefix(t *testing.T) {
	tr := NewTrie()
	words := []string{"cat", "car", "card", "care", "dog"}
	for _, w := range words {
		tr.Insert(w)
	}

	got := tr.WordsWithPrefix("")
	if len(got) != len(words) {
		t.Errorf("WordsWithPrefix(\"\") should return all %d words, got %v", len(words), got)
	}
}

func TestTrie_WordsWithPrefix_Sorted(t *testing.T) {
	tr := NewTrie()
	for _, w := range []string{"zoo", "zebra", "zero", "zeal"} {
		tr.Insert(w)
	}

	got := tr.WordsWithPrefix("ze")
	want := []string{"zeal", "zebra", "zero"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordsWithPrefix(\"ze\") = %v, want %v (must be sorted)", got, want)
	}
}

func TestTrie_InsertDuplicate(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")
	tr.Insert("hello") // insert same word twice

	got := tr.WordsWithPrefix("hello")
	if len(got) != 1 {
		t.Errorf("WordsWithPrefix after duplicate insert should return 1 word, got %v", got)
	}
}
