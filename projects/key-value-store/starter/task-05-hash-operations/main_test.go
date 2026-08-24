package main

import "testing"

func TestHSet_HGet(t *testing.T) {
	s := NewStore()
	s.HSet("user:1", "name", "Alice")
	s.HSet("user:1", "age", "30")

	v, ok := s.HGet("user:1", "name")
	if !ok || v != "Alice" {
		t.Fatalf("expected Alice, got %q ok=%v", v, ok)
	}

	v, ok = s.HGet("user:1", "age")
	if !ok || v != "30" {
		t.Fatalf("expected 30, got %q ok=%v", v, ok)
	}
}

func TestHGet_MissingField(t *testing.T) {
	s := NewStore()
	s.HSet("user:1", "name", "Alice")

	_, ok := s.HGet("user:1", "missing")
	if ok {
		t.Fatal("expected ok=false for missing field")
	}
}

func TestHGet_MissingKey(t *testing.T) {
	s := NewStore()
	_, ok := s.HGet("nokey", "field")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestHGetAll(t *testing.T) {
	s := NewStore()
	s.HSet("h", "f1", "v1")
	s.HSet("h", "f2", "v2")

	all := s.HGetAll("h")
	if len(all) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(all))
	}
	if all["f1"] != "v1" || all["f2"] != "v2" {
		t.Fatalf("unexpected values: %v", all)
	}
}

func TestHGetAll_MissingKey(t *testing.T) {
	s := NewStore()
	empty := s.HGetAll("missing")
	if len(empty) != 0 {
		t.Fatalf("expected empty map for missing key, got %v", empty)
	}
}

func TestHDel(t *testing.T) {
	s := NewStore()
	s.HSet("h", "f1", "v1")
	s.HSet("h", "f2", "v2")
	s.HSet("h", "f3", "v3")

	count := s.HDel("h", "f1", "f2")
	if count != 2 {
		t.Fatalf("expected count=2, got %d", count)
	}

	all := s.HGetAll("h")
	if len(all) != 1 || all["f3"] != "v3" {
		t.Fatalf("unexpected remaining fields: %v", all)
	}
}

func TestHDel_MissingField(t *testing.T) {
	s := NewStore()
	s.HSet("h", "f1", "v1")

	count := s.HDel("h", "missing")
	if count != 0 {
		t.Fatalf("expected count=0 for missing field, got %d", count)
	}
}

func TestHExists(t *testing.T) {
	s := NewStore()
	s.HSet("h", "field", "val")

	if !s.HExists("h", "field") {
		t.Fatal("expected true for existing field")
	}
	if s.HExists("h", "other") {
		t.Fatal("expected false for missing field")
	}
	if s.HExists("missing", "field") {
		t.Fatal("expected false for missing key")
	}
}
