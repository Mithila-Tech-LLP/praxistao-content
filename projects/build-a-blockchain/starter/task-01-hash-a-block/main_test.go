package main

import (
	"bytes"
	"testing"
)

func TestHashDeterministic(t *testing.T) {
	h1 := Hash([]byte("hello"))
	h2 := Hash([]byte("hello"))
	if !bytes.Equal(h1, h2) {
		t.Fatal("expected the same input to hash to the same output")
	}
}

func TestHashAvalanche(t *testing.T) {
	h1 := Hash([]byte("hello"))
	h2 := Hash([]byte("hellp"))
	if bytes.Equal(h1, h2) {
		t.Fatal("expected different inputs to produce different hashes")
	}
	// Sanity: hashes should be 32 bytes (SHA-256).
	if len(h1) != 32 {
		t.Fatalf("expected a 32-byte hash, got %d bytes", len(h1))
	}
}

func TestNewBlockSetsHash(t *testing.T) {
	b := NewBlock([]byte("genesis"), []byte{})
	if len(b.Hash) == 0 {
		t.Fatal("expected NewBlock to set a non-empty Hash")
	}
	if !bytes.Equal(b.Hash, b.ComputeHash()) {
		t.Fatal("expected b.Hash to match a freshly recomputed hash")
	}
}

func TestComputeHashExcludesHashField(t *testing.T) {
	b := NewBlock([]byte("data"), []byte{})
	original := b.ComputeHash()

	// Deliberately corrupt the stored Hash field only.
	b.Hash = []byte("not the real hash")

	// ComputeHash must still return the SAME value as before -- proving
	// it never serialized b.Hash into its own computation.
	if !bytes.Equal(b.ComputeHash(), original) {
		t.Fatal("expected ComputeHash to be unaffected by the current value of b.Hash")
	}
}

func TestDifferentDataDifferentHash(t *testing.T) {
	b1 := NewBlock([]byte("data one"), []byte{})
	b2 := NewBlock([]byte("data two"), []byte{})
	if bytes.Equal(b1.Hash, b2.Hash) {
		t.Fatal("expected blocks with different data to have different hashes")
	}
}

func TestSamePrevHashLinkage(t *testing.T) {
	genesis := NewBlock([]byte("genesis"), []byte{})
	next := NewBlock([]byte("second block"), genesis.Hash)
	if !bytes.Equal(next.PrevBlockHash, genesis.Hash) {
		t.Fatal("expected next.PrevBlockHash to equal genesis.Hash")
	}
}
