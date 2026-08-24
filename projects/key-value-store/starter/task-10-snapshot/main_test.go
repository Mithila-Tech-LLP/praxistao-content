package main

import (
	"bytes"
	"testing"
)

func TestSaveLoad_ThreeKeys(t *testing.T) {
	s := NewStore()
	s.Set("key1", "value1")
	s.Set("key2", "value2")
	s.Set("key3", "value3")

	var buf bytes.Buffer
	if err := s.Save(&buf); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	s2 := NewStore()
	if err := s2.Load(&buf); err != nil {
		t.Fatalf("Load error: %v", err)
	}

	for key, want := range map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	} {
		v, ok := s2.Get(key)
		if !ok || v != want {
			t.Fatalf("key %q: expected %q got %q ok=%v", key, want, v, ok)
		}
	}
}

func TestLoad_ReplacesExistingData(t *testing.T) {
	s := NewStore()
	s.Set("a", "1")

	var buf bytes.Buffer
	s.Save(&buf)

	s2 := NewStore()
	s2.Set("b", "2") // pre-existing data in target store
	s2.Load(&buf)

	_, ok := s2.Get("b")
	if ok {
		t.Fatal("Load should replace existing data; 'b' should not exist after load")
	}
	v, ok := s2.Get("a")
	if !ok || v != "1" {
		t.Fatalf("expected a=1 after load, got %q ok=%v", v, ok)
	}
}

func TestSaveLoad_EmptyStore(t *testing.T) {
	s := NewStore()

	var buf bytes.Buffer
	if err := s.Save(&buf); err != nil {
		t.Fatalf("Save of empty store error: %v", err)
	}

	s2 := NewStore()
	s2.Set("x", "y")
	if err := s2.Load(&buf); err != nil {
		t.Fatalf("Load error: %v", err)
	}
	_, ok := s2.Get("x")
	if ok {
		t.Fatal("loading empty snapshot should clear existing keys")
	}
}
