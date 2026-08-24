package main

import "testing"

func TestLRU_BasicEviction(t *testing.T) {
	s := NewStoreWithCapacity(2)
	s.Set("A", "1")
	s.Set("B", "2")
	s.Set("C", "3") // capacity full: A (LRU) should be evicted

	_, ok := s.Get("A")
	if ok {
		t.Fatal("A should have been evicted (was LRU)")
	}
	v, ok := s.Get("B")
	if !ok || v != "2" {
		t.Fatalf("B should still exist, got %q ok=%v", v, ok)
	}
	v, ok = s.Get("C")
	if !ok || v != "3" {
		t.Fatalf("C should exist, got %q ok=%v", v, ok)
	}
}

func TestLRU_GetRefreshesRecency(t *testing.T) {
	s := NewStoreWithCapacity(2)
	s.Set("A", "1")
	s.Set("B", "2")
	s.Get("A") // A is now most recently used; B becomes LRU
	s.Set("D", "4") // B should be evicted, not A

	_, ok := s.Get("B")
	if ok {
		t.Fatal("B should have been evicted (was LRU after Get(A))")
	}
	_, ok = s.Get("A")
	if !ok {
		t.Fatal("A should still exist (was refreshed by Get)")
	}
	_, ok = s.Get("D")
	if !ok {
		t.Fatal("D should exist")
	}
}

func TestLRU_UpdateExistingKeyDoesNotEvict(t *testing.T) {
	s := NewStoreWithCapacity(2)
	s.Set("A", "1")
	s.Set("B", "2")
	s.Set("A", "updated") // update existing key — should NOT evict B

	v, ok := s.Get("A")
	if !ok || v != "updated" {
		t.Fatalf("A should be updated, got %q ok=%v", v, ok)
	}
	_, ok = s.Get("B")
	if !ok {
		t.Fatal("B should not be evicted when updating an existing key")
	}
}

func TestLRU_CapacityOne(t *testing.T) {
	s := NewStoreWithCapacity(1)
	s.Set("X", "x")
	s.Set("Y", "y") // X evicted

	_, ok := s.Get("X")
	if ok {
		t.Fatal("X should have been evicted")
	}
	v, ok := s.Get("Y")
	if !ok || v != "y" {
		t.Fatalf("Y should exist, got %q ok=%v", v, ok)
	}
}

func TestLRU_EvictionOrderAfterMultipleGets(t *testing.T) {
	s := NewStoreWithCapacity(3)
	s.Set("A", "1")
	s.Set("B", "2")
	s.Set("C", "3")
	s.Get("A") // order now: C(LRU) B A(MRU)
	s.Get("B") // order now: C(LRU) A B(MRU)
	s.Set("D", "4") // C should be evicted

	_, ok := s.Get("C")
	if ok {
		t.Fatal("C should have been evicted (was LRU)")
	}
	for _, k := range []string{"A", "B", "D"} {
		if _, ok := s.Get(k); !ok {
			t.Fatalf("%s should still exist", k)
		}
	}
}
