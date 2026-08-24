package main

import (
	"bytes"
	"testing"
)

func TestRunFindsValidNonce(t *testing.T) {
	b := &Block{Timestamp: 1, Data: []byte("hello"), PrevBlockHash: make([]byte, 32)}
	pow := NewProofOfWork(b, 12) // small difficulty so the test is fast

	nonce, hash := pow.Run()

	b.Nonce = nonce
	b.Hash = hash

	if !pow.Validate() {
		t.Fatal("expected the nonce/hash Run() found to satisfy Validate()")
	}
}

func TestValidateRejectsWrongNonce(t *testing.T) {
	b := &Block{Timestamp: 1, Data: []byte("hello"), PrevBlockHash: make([]byte, 32)}
	pow := NewProofOfWork(b, 12)

	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash

	if !pow.Validate() {
		t.Fatal("expected the real solution to validate")
	}

	// Corrupt the nonce without recomputing the hash.
	b.Nonce = nonce + 1
	if pow.Validate() {
		t.Fatal("expected Validate to reject a nonce that doesn't match the stored hash")
	}
}

func TestValidateRejectsWrongHash(t *testing.T) {
	b := &Block{Timestamp: 1, Data: []byte("hello"), PrevBlockHash: make([]byte, 32)}
	pow := NewProofOfWork(b, 12)

	nonce, _ := pow.Run()
	b.Nonce = nonce
	b.Hash = []byte("not the real hash, but exactly 32 bytes long!!")

	if pow.Validate() {
		t.Fatal("expected Validate to reject a stored hash that doesn't match recomputation")
	}
}

func TestHigherDifficultyStricterTarget(t *testing.T) {
	low := NewProofOfWork(&Block{}, 8)
	high := NewProofOfWork(&Block{}, 16)

	// A higher difficulty means a SMALLER target (fewer valid hashes).
	if high.Target.Cmp(low.Target) != -1 {
		t.Fatal("expected higher difficulty to produce a smaller target")
	}
}

func TestDifferentDataDifferentSolution(t *testing.T) {
	b1 := &Block{Timestamp: 1, Data: []byte("data one"), PrevBlockHash: make([]byte, 32)}
	b2 := &Block{Timestamp: 1, Data: []byte("data two"), PrevBlockHash: make([]byte, 32)}

	_, hash1 := NewProofOfWork(b1, 12).Run()
	_, hash2 := NewProofOfWork(b2, 12).Run()

	if bytes.Equal(hash1, hash2) {
		t.Fatal("expected different block data to mine to different hashes")
	}
}
