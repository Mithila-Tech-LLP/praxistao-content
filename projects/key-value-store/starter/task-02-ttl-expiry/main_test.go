package main

import (
	"sync"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s := NewStore()
	t.Cleanup(s.Close)
	return s
}

// --- basic TTL ---

func TestSetExpires(t *testing.T) {
	s := newTestStore(t)
	s.SetWithTTL("k", "v", 50*time.Millisecond)

	v, ok := s.Get("k")
	if !ok || v != "v" {
		t.Fatalf("expected key to exist immediately after SetWithTTL")
	}

	time.Sleep(120 * time.Millisecond)

	_, ok = s.Get("k")
	if ok {
		t.Fatal("expected key to be expired after TTL elapsed")
	}
}

func TestSetClearsTTL(t *testing.T) {
	s := newTestStore(t)
	s.SetWithTTL("k", "v", 50*time.Millisecond)
	// overwrite without TTL — key should live forever
	s.Set("k", "persistent")

	time.Sleep(120 * time.Millisecond)

	v, ok := s.Get("k")
	if !ok || v != "persistent" {
		t.Fatalf("expected key to persist after Set cleared TTL, got %q ok=%v", v, ok)
	}
}

func TestTTLReturnValue(t *testing.T) {
	s := newTestStore(t)
	s.SetWithTTL("k", "v", 1*time.Second)

	remaining, ok := s.TTL("k")
	if !ok {
		t.Fatal("expected TTL to return ok=true")
	}
	if remaining <= 0 || remaining > 1*time.Second {
		t.Fatalf("unexpected remaining TTL: %v", remaining)
	}
}

func TestTTLMissingKey(t *testing.T) {
	s := newTestStore(t)
	_, ok := s.TTL("ghost")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestTTLNoExpiry(t *testing.T) {
	s := newTestStore(t)
	s.Set("k", "v")
	_, ok := s.TTL("k")
	if ok {
		t.Fatal("expected ok=false for key without TTL")
	}
}

func TestDeleteExpiredKey(t *testing.T) {
	s := newTestStore(t)
	s.SetWithTTL("k", "v", 30*time.Millisecond)
	time.Sleep(60 * time.Millisecond)
	deleted := s.Delete("k")
	if deleted {
		t.Fatal("expected Delete to return false for already-expired key")
	}
}

// --- multiple keys with different TTLs ---

func TestMultipleTTLs(t *testing.T) {
	s := newTestStore(t)
	s.SetWithTTL("short", "s", 40*time.Millisecond)
	s.SetWithTTL("long", "l", 500*time.Millisecond)
	s.Set("immortal", "i")

	time.Sleep(80 * time.Millisecond)

	if _, ok := s.Get("short"); ok {
		t.Error("'short' should have expired")
	}
	if _, ok := s.Get("long"); !ok {
		t.Error("'long' should still be alive")
	}
	if _, ok := s.Get("immortal"); !ok {
		t.Error("'immortal' should still be alive")
	}
}

// --- concurrency ---

func TestConcurrentTTLSets(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	const n = 200
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.SetWithTTL("racing", "v", 50*time.Millisecond)
		}()
	}
	wg.Wait()
}

func TestConcurrentGetAfterExpiry(t *testing.T) {
	t.Parallel()
	s := newTestStore(t)
	s.SetWithTTL("k", "v", 30*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Get("k") // must not panic
		}()
	}
	wg.Wait()
}
