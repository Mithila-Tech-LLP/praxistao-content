package main

import (
	"sort"
	"sync"
	"testing"
)

// --- SAdd ---

func TestSAddBasic(t *testing.T) {
	s := NewStore()
	n := s.SAdd("tags", "go", "redis", "cache")
	if n != 3 {
		t.Fatalf("expected 3 added, got %d", n)
	}
}

func TestSAddDuplicates(t *testing.T) {
	s := NewStore()
	s.SAdd("tags", "go", "redis")
	n := s.SAdd("tags", "go", "python") // "go" is duplicate
	if n != 1 {
		t.Fatalf("expected 1 newly added, got %d", n)
	}
}

func TestSAddIdempotent(t *testing.T) {
	s := NewStore()
	s.SAdd("s", "a", "a", "a")
	if s.SCard("s") != 1 {
		t.Fatalf("expected cardinality 1, got %d", s.SCard("s"))
	}
}

// --- SMembers ---

func TestSMembers(t *testing.T) {
	s := NewStore()
	s.SAdd("s", "c", "a", "b")
	members := s.SMembers("s")
	sort.Strings(members)
	want := []string{"a", "b", "c"}
	for i, w := range want {
		if members[i] != w {
			t.Fatalf("SMembers[%d]: got %q, want %q", i, members[i], w)
		}
	}
}

func TestSMembersEmpty(t *testing.T) {
	s := NewStore()
	got := s.SMembers("ghost")
	if got == nil {
		t.Fatal("SMembers should return non-nil slice for missing key")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

// --- SIsMember ---

func TestSIsMember(t *testing.T) {
	s := NewStore()
	s.SAdd("s", "alice", "bob")

	if !s.SIsMember("s", "alice") {
		t.Fatal("expected alice to be member")
	}
	if s.SIsMember("s", "carol") {
		t.Fatal("carol should not be member")
	}
	if s.SIsMember("ghost", "alice") {
		t.Fatal("missing key should return false")
	}
}

// --- SRem ---

func TestSRem(t *testing.T) {
	s := NewStore()
	s.SAdd("s", "a", "b", "c")
	n := s.SRem("s", "a", "c", "x") // "x" not present
	if n != 2 {
		t.Fatalf("expected 2 removed, got %d", n)
	}
	if s.SIsMember("s", "a") {
		t.Fatal("'a' should have been removed")
	}
	if !s.SIsMember("s", "b") {
		t.Fatal("'b' should still be present")
	}
}

func TestSRemMissingKey(t *testing.T) {
	s := NewStore()
	n := s.SRem("ghost", "x")
	if n != 0 {
		t.Fatalf("expected 0 for missing key, got %d", n)
	}
}

// --- SCard ---

func TestSCard(t *testing.T) {
	s := NewStore()
	if s.SCard("s") != 0 {
		t.Fatal("expected 0 for missing key")
	}
	s.SAdd("s", "a", "b", "c")
	if s.SCard("s") != 3 {
		t.Fatalf("expected 3, got %d", s.SCard("s"))
	}
	s.SRem("s", "a")
	if s.SCard("s") != 2 {
		t.Fatalf("expected 2 after remove, got %d", s.SCard("s"))
	}
}

// --- concurrency ---

func TestConcurrentSAdd(t *testing.T) {
	t.Parallel()
	s := NewStore()
	const n = 500
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			s.SAdd("shared", string(rune('a'+i%26)))
		}()
	}
	wg.Wait()
	// cardinality should be at most 26 (alphabet) and at least 1
	c := s.SCard("shared")
	if c < 1 || c > 26 {
		t.Fatalf("unexpected cardinality %d", c)
	}
}

func TestConcurrentSAddSRem(t *testing.T) {
	t.Parallel()
	s := NewStore()
	s.SAdd("s", "a", "b", "c", "d", "e")

	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() { defer wg.Done(); s.SAdd("s", "z") }()
		go func() { defer wg.Done(); s.SRem("s", "z") }()
	}
	wg.Wait()
}
