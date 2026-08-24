package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

// Block is one entry in the chain. Nonce is filled in by mining.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// serializeForMining builds the bytes proof-of-work hashes: everything
// that identifies the block's content, plus the candidate nonce.
// (Provided for you -- this is the same idea as Task 01's Serialize,
// extended with the nonce.)
func serializeForMining(b *Block, nonce int64) []byte {
	var buf bytes.Buffer
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(b.Timestamp))
	buf.Write(ts)
	buf.Write(b.Data)
	buf.Write(b.PrevBlockHash)
	n := make([]byte, 8)
	binary.BigEndian.PutUint64(n, uint64(nonce))
	buf.Write(n)
	return buf.Bytes()
}

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProofOfWork creates a ProofOfWork for b with the given difficulty
// (number of leading zero bits the resulting hash must have).
func NewProofOfWork(b *Block, difficulty int) *ProofOfWork {
	panic("TODO: implement NewProofOfWork -- build Target = 1 << (256 - difficulty)")
}

// Run searches for a nonce that satisfies the target, starting from 0,
// and returns the winning nonce and its resulting hash.
func (pow *ProofOfWork) Run() (int64, []byte) {
	panic("TODO: implement Run")
}

// Validate reports whether pow.Block's current Nonce and Hash actually
// satisfy pow.Target.
func (pow *ProofOfWork) Validate() bool {
	panic("TODO: implement Validate")
}
