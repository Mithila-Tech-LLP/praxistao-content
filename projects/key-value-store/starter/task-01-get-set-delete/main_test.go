package main

import (
	"fmt"
	"sync"
	"testing"
)

// --- basic CRUD ---

func TestSetAndGet(t *testing.T) {
	s := NewStore()
	s.Set("name", "alice")

	v, ok := s.Get("name")
	if !ok {
		t.Fatal("expected key 'name' to exist")
	}
	if v != "alice" {
		t.Fatalf("got %q, want %q", v, "alice")
	}
}

func TestGetMissing(t *testing.T) {
	s := NewStore()
	v, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
	if v != "" {
		t.Fatalf("expected empty string, got %q", v)
	}
}

func TestOverwrite(t *testing.T) {
	s := NewStore()
	s.Set("k", "first")
	s.Set("k", "second")

	v, ok := s.Get("k")
	if !ok || v != "second" {
		t.Fatalf("expected 'second', got %q (ok=%v)", v, ok)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()
	s.Set("x", "1")

	deleted := s.Delete("x")
	if !deleted {
		t.Fatal("expected Delete to return true for existing key")
	}

	_, ok := s.Get("x")
	if ok {
		t.Fatal("expected key to be gone after Delete")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewStore()
	deleted := s.Delete("nope")
	if deleted {
		t.Fatal("expected Delete to return false for missing key")
	}
}

func TestExists(t *testing.T) {
	s := NewStore()
	if s.Exists("k") {
		t.Fatal("expected Exists=false before Set")
	}
	s.Set("k", "v")
	if !s.Exists("k") {
		t.Fatal("expected Exists=true after Set")
	}
	s.Delete("k")
	if s.Exists("k") {
		t.Fatal("expected Exists=false after Delete")
	}
}

func TestMultipleKeys(t *testing.T) {
	s := NewStore()
	pairs := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	for k, v := range pairs {
		s.Set(k, v)
	}
	for k, want := range pairs {
		got, ok := s.Get(k)
		if !ok || got != want {
			t.Errorf("key %q: got %q (ok=%v), want %q", k, got, ok, want)
		}
	}
}

// --- concurrency ---

func TestConcurrentWrites(t *testing.T) {
	t.Parallel()
	s := NewStore()
	const n = 1000
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			s.Set(key, fmt.Sprintf("val-%d", i))
		}()
	}
	wg.Wait()
	// just checking no panic / data race; count keys
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("key-%d", i)
		if !s.Exists(key) {
			t.Errorf("key %q missing after concurrent write", key)
		}
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	t.Parallel()
	s := NewStore()
	s.Set("shared", "0")

	var wg sync.WaitGroup
	const n = 500
	wg.Add(n * 2)

	// writers
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			s.Set("shared", fmt.Sprintf("%d", i))
		}()
	}
	// readers — must not panic or return inconsistent state
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Get("shared")
		}()
	}
	wg.Wait()
}

func TestConcurrentDelete(t *testing.T) {
	t.Parallel()
	s := NewStore()
	const n = 200
	for i := 0; i < n; i++ {
		s.Set(fmt.Sprintf("k%d", i), "v")
	}
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			s.Delete(fmt.Sprintf("k%d", i))
		}()
	}
	wg.Wait()
}
