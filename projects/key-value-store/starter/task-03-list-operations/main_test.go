package main

import (
	"reflect"
	"sync"
	"testing"
)

// --- LPush / RPush ---

func TestRPushOrder(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b", "c")
	got := s.LRange("q", 0, -1)
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("RPush: got %v, want %v", got, want)
	}
}

func TestLPushOrder(t *testing.T) {
	s := NewStore()
	s.LPush("q", "a", "b", "c")
	got := s.LRange("q", 0, -1)
	// each element is prepended: after "a" → ["a"], after "b" → ["b","a"], after "c" → ["c","b","a"]
	want := []string{"c", "b", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LPush: got %v, want %v", got, want)
	}
}

func TestPushReturnLength(t *testing.T) {
	s := NewStore()
	n := s.RPush("q", "x")
	if n != 1 {
		t.Fatalf("expected length 1, got %d", n)
	}
	n = s.RPush("q", "y", "z")
	if n != 3 {
		t.Fatalf("expected length 3, got %d", n)
	}
}

// --- LPop / RPop ---

func TestLPop(t *testing.T) {
	s := NewStore()
	s.RPush("q", "first", "second", "third")

	v, ok := s.LPop("q")
	if !ok || v != "first" {
		t.Fatalf("LPop: got %q ok=%v, want 'first' ok=true", v, ok)
	}
	if s.LLen("q") != 2 {
		t.Fatalf("expected length 2 after LPop, got %d", s.LLen("q"))
	}
}

func TestRPop(t *testing.T) {
	s := NewStore()
	s.RPush("q", "first", "second", "third")

	v, ok := s.RPop("q")
	if !ok || v != "third" {
		t.Fatalf("RPop: got %q ok=%v, want 'third' ok=true", v, ok)
	}
}

func TestPopEmpty(t *testing.T) {
	s := NewStore()
	_, ok := s.LPop("empty")
	if ok {
		t.Fatal("LPop on missing key should return ok=false")
	}
	_, ok = s.RPop("empty")
	if ok {
		t.Fatal("RPop on missing key should return ok=false")
	}
}

func TestPopUntilEmpty(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b")
	s.LPop("q")
	s.LPop("q")
	_, ok := s.LPop("q")
	if ok {
		t.Fatal("LPop on emptied list should return ok=false")
	}
}

// --- LRange ---

func TestLRangeAll(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b", "c", "d")
	got := s.LRange("q", 0, -1)
	want := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LRange(0,-1): got %v, want %v", got, want)
	}
}

func TestLRangeSubset(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b", "c", "d")
	got := s.LRange("q", 1, 2)
	want := []string{"b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LRange(1,2): got %v, want %v", got, want)
	}
}

func TestLRangeNegativeIndices(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b", "c", "d")
	got := s.LRange("q", -2, -1)
	want := []string{"c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LRange(-2,-1): got %v, want %v", got, want)
	}
}

func TestLRangeOutOfBounds(t *testing.T) {
	s := NewStore()
	s.RPush("q", "a", "b")
	got := s.LRange("q", 0, 100)
	want := []string{"a", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LRange(0,100): got %v, want %v", got, want)
	}
}

func TestLRangeMissing(t *testing.T) {
	s := NewStore()
	got := s.LRange("ghost", 0, -1)
	if len(got) != 0 {
		t.Fatalf("expected empty slice for missing key, got %v", got)
	}
}

// --- LLen ---

func TestLLen(t *testing.T) {
	s := NewStore()
	if s.LLen("q") != 0 {
		t.Fatal("expected LLen 0 for missing key")
	}
	s.RPush("q", "x", "y", "z")
	if s.LLen("q") != 3 {
		t.Fatalf("expected LLen 3, got %d", s.LLen("q"))
	}
}

// --- concurrency ---

func TestConcurrentRPush(t *testing.T) {
	t.Parallel()
	s := NewStore()
	const n = 500
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.RPush("shared", "item")
		}()
	}
	wg.Wait()
	if s.LLen("shared") != n {
		t.Fatalf("expected %d items, got %d", n, s.LLen("shared"))
	}
}

func TestConcurrentLPopRPush(t *testing.T) {
	t.Parallel()
	s := NewStore()
	for i := 0; i < 100; i++ {
		s.RPush("q", "v")
	}
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() { defer wg.Done(); s.LPop("q") }()
		go func() { defer wg.Done(); s.RPush("q", "new") }()
	}
	wg.Wait()
}
