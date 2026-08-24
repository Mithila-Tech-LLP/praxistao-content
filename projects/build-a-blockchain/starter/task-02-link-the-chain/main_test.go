package main

import "testing"

func TestNewBlockchainStartsValid(t *testing.T) {
	bc := NewBlockchain()
	if len(bc.Blocks) != 1 {
		t.Fatalf("expected exactly 1 genesis block, got %d", len(bc.Blocks))
	}
	if !bc.IsValid() {
		t.Fatal("expected a fresh chain to be valid")
	}
}

func TestAddBlockLinksCorrectly(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("first")
	bc.AddBlock("second")
	bc.AddBlock("third")

	if len(bc.Blocks) != 4 {
		t.Fatalf("expected 4 blocks (genesis + 3), got %d", len(bc.Blocks))
	}
	if !bc.IsValid() {
		t.Fatal("expected a chain built only with AddBlock to be valid")
	}
}

func TestTamperedDataInvalidatesChain(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("first")
	bc.AddBlock("second")

	if !bc.IsValid() {
		t.Fatal("expected chain to be valid before tampering")
	}

	// Corrupt an old block's data WITHOUT recomputing its hash.
	bc.Blocks[1].Data = []byte("tampered!")

	if bc.IsValid() {
		t.Fatal("expected IsValid to detect tampering with an old block's data")
	}
}

func TestTamperedPrevHashInvalidatesChain(t *testing.T) {
	bc := NewBlockchain()
	bc.AddBlock("first")
	bc.AddBlock("second")

	bc.Blocks[2].PrevBlockHash = []byte("not the real previous hash")

	if bc.IsValid() {
		t.Fatal("expected IsValid to detect a broken PrevBlockHash link")
	}
}

func TestGenesisHasZeroPrevHash(t *testing.T) {
	g := NewGenesisBlock()
	for _, bb := range g.PrevBlockHash {
		if bb != 0 {
			t.Fatal("expected genesis block's PrevBlockHash to be all zero bytes")
		}
	}
}
