package main

import "testing"

func TestTx_ExecAppliesAll(t *testing.T) {
	s := NewStore()
	tx := s.Begin()
	tx.Set("a", "1")
	tx.Set("b", "2")
	tx.Set("c", "3")

	if err := tx.Exec(); err != nil {
		t.Fatalf("Exec error: %v", err)
	}

	for key, want := range map[string]string{"a": "1", "b": "2", "c": "3"} {
		v, ok := s.Get(key)
		if !ok || v != want {
			t.Fatalf("key %q: expected %q got %q ok=%v", key, want, v, ok)
		}
	}
}

func TestTx_Discard(t *testing.T) {
	s := NewStore()
	s.Set("x", "original")
	tx := s.Begin()
	tx.Set("x", "modified")
	tx.Discard()

	v, _ := s.Get("x")
	if v != "original" {
		t.Fatalf("expected original after discard, got %q", v)
	}
}

func TestTx_Delete(t *testing.T) {
	s := NewStore()
	s.Set("k", "v")
	tx := s.Begin()
	tx.Delete("k")
	tx.Exec()

	_, ok := s.Get("k")
	if ok {
		t.Fatal("expected key deleted after tx exec")
	}
}

func TestTx_ExecOnDiscardedTx(t *testing.T) {
	s := NewStore()
	tx := s.Begin()
	tx.Set("a", "1")
	tx.Discard()

	// Exec after Discard should be a no-op (queue was cleared)
	if err := tx.Exec(); err != nil {
		t.Fatalf("Exec after discard should not error, got: %v", err)
	}

	_, ok := s.Get("a")
	if ok {
		t.Fatal("key should not exist: queue was discarded before exec")
	}
}

func TestTx_SetAfterExecIsNoOp(t *testing.T) {
	s := NewStore()
	tx := s.Begin()
	tx.Set("a", "1")
	tx.Exec()

	// Exec clears the queue. A Set queued after Exec
	// has no effect on the store unless Exec is called again.
	tx.Set("a", "changed")
	// NOT calling Exec — the store should retain the first exec's value
	v, ok := s.Get("a")
	if !ok || v != "1" {
		t.Fatalf("expected a=1 (no re-exec after second Set), got %q ok=%v", v, ok)
	}
}
