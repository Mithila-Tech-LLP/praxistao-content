package main

import (
	"sort"
	"testing"
)

func TestHashMap_SetAndGet(t *testing.T) {
	m := &HashMap{}
	m.Set("foo", 1)
	m.Set("bar", 2)
	m.Set("baz", 3)

	cases := []struct {
		key  string
		want int
	}{
		{"foo", 1},
		{"bar", 2},
		{"baz", 3},
	}
	for _, tc := range cases {
		got, ok := m.Get(tc.key)
		if !ok {
			t.Errorf("Get(%q) returned false, want true", tc.key)
		}
		if got != tc.want {
			t.Errorf("Get(%q) = %d, want %d", tc.key, got, tc.want)
		}
	}
}

func TestHashMap_GetMissingKey(t *testing.T) {
	m := &HashMap{}
	m.Set("hello", 99)

	_, ok := m.Get("world")
	if ok {
		t.Error("Get on missing key should return false")
	}
}

func TestHashMap_Overwrite(t *testing.T) {
	m := &HashMap{}
	m.Set("key", 10)
	m.Set("key", 20) // overwrite

	got, ok := m.Get("key")
	if !ok {
		t.Error("Get after overwrite returned false")
	}
	if got != 20 {
		t.Errorf("Get after overwrite = %d, want 20", got)
	}

	// Keys() must not contain duplicate
	keys := m.Keys()
	count := 0
	for _, k := range keys {
		if k == "key" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Keys() should contain 'key' exactly once after overwrite, found %d times", count)
	}
}

func TestHashMap_Delete(t *testing.T) {
	m := &HashMap{}
	m.Set("a", 1)
	m.Set("b", 2)
	m.Delete("a")

	_, ok := m.Get("a")
	if ok {
		t.Error("Get after Delete should return false")
	}

	got, ok := m.Get("b")
	if !ok || got != 2 {
		t.Errorf("Get(b) after deleting a = (%d, %v), want (2, true)", got, ok)
	}
}

func TestHashMap_DeleteNonExistent(t *testing.T) {
	m := &HashMap{}
	// deleting a key that was never set should not panic
	m.Delete("ghost")
}

func TestHashMap_Keys(t *testing.T) {
	m := &HashMap{}
	inserted := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for i, k := range inserted {
		m.Set(k, i)
	}

	keys := m.Keys()
	sort.Strings(keys)
	sort.Strings(inserted)

	if len(keys) != len(inserted) {
		t.Errorf("Keys() length = %d, want %d", len(keys), len(inserted))
	}
	for i := range inserted {
		if keys[i] != inserted[i] {
			t.Errorf("Keys()[%d] = %q, want %q", i, keys[i], inserted[i])
		}
	}
}

func TestHashMap_KeysAfterDelete(t *testing.T) {
	m := &HashMap{}
	m.Set("x", 1)
	m.Set("y", 2)
	m.Set("z", 3)
	m.Delete("y")

	keys := m.Keys()
	for _, k := range keys {
		if k == "y" {
			t.Error("Keys() should not contain deleted key 'y'")
		}
	}
	if len(keys) != 2 {
		t.Errorf("Keys() length after delete = %d, want 2", len(keys))
	}
}

func TestHashMap_ZeroValue(t *testing.T) {
	m := &HashMap{}
	m.Set("zero", 0)

	got, ok := m.Get("zero")
	if !ok {
		t.Error("Get for key with zero value should return true")
	}
	if got != 0 {
		t.Errorf("Get for zero value = %d, want 0", got)
	}
}
