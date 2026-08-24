package main

import "testing"

func TestLRUCache_BasicPutGet(t *testing.T) {
	c := NewLRUCache(3)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(3, 30)

	cases := []struct {
		key  int
		want int
	}{
		{1, 10},
		{2, 20},
		{3, 30},
	}
	for _, tc := range cases {
		got, ok := c.Get(tc.key)
		if !ok {
			t.Errorf("Get(%d) returned false, want true", tc.key)
		}
		if got != tc.want {
			t.Errorf("Get(%d) = %d, want %d", tc.key, got, tc.want)
		}
	}
}

func TestLRUCache_GetMissing(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 100)

	_, ok := c.Get(99)
	if ok {
		t.Error("Get on missing key should return false")
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	// capacity 2: put 1, put 2, put 3 → key 1 should be evicted (LRU)
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(3, 30) // evicts key 1

	_, ok := c.Get(1)
	if ok {
		t.Error("Key 1 should have been evicted after adding key 3 to a full cache")
	}

	got2, ok := c.Get(2)
	if !ok || got2 != 20 {
		t.Errorf("Get(2) = (%d, %v), want (20, true)", got2, ok)
	}

	got3, ok := c.Get(3)
	if !ok || got3 != 30 {
		t.Errorf("Get(3) = (%d, %v), want (30, true)", got3, ok)
	}
}

func TestLRUCache_ReAccessPreventsEviction(t *testing.T) {
	// capacity 2: put 1, put 2, get 1 (makes 2 the LRU), put 3 → key 2 should be evicted
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Get(1) // re-access key 1 — now key 2 is LRU
	c.Put(3, 30) // evicts key 2

	_, ok := c.Get(2)
	if ok {
		t.Error("Key 2 should have been evicted (it was LRU when key 3 was inserted)")
	}

	got1, ok := c.Get(1)
	if !ok || got1 != 10 {
		t.Errorf("Get(1) = (%d, %v), want (10, true) — key 1 was recently accessed", got1, ok)
	}

	got3, ok := c.Get(3)
	if !ok || got3 != 30 {
		t.Errorf("Get(3) = (%d, %v), want (30, true)", got3, ok)
	}
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(1, 99) // update value

	got, ok := c.Get(1)
	if !ok || got != 99 {
		t.Errorf("Get(1) after update = (%d, %v), want (99, true)", got, ok)
	}
}

func TestLRUCache_UpdateDoesNotGrowCapacity(t *testing.T) {
	// updating an existing key must not count as a new entry
	c := NewLRUCache(2)
	c.Put(1, 10)
	c.Put(2, 20)
	c.Put(1, 11) // update key 1, not a new entry
	c.Put(3, 30) // now at capacity=2 with keys 1 and 3; key 2 should be evicted

	_, ok := c.Get(2)
	if ok {
		t.Error("Key 2 should be evicted — updating key 1 made key 2 the LRU")
	}

	got1, ok := c.Get(1)
	if !ok || got1 != 11 {
		t.Errorf("Get(1) = (%d, %v), want (11, true)", got1, ok)
	}
}

func TestLRUCache_CapacityOne(t *testing.T) {
	c := NewLRUCache(1)
	c.Put(1, 10)
	c.Put(2, 20) // evicts key 1

	_, ok := c.Get(1)
	if ok {
		t.Error("Key 1 should be evicted in capacity-1 cache when key 2 is added")
	}

	got, ok := c.Get(2)
	if !ok || got != 20 {
		t.Errorf("Get(2) = (%d, %v), want (20, true)", got, ok)
	}
}
