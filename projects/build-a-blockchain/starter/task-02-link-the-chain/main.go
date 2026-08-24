package main

import (
	"crypto/sha256"
)

// Block is one entry in the chain.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// Hash returns the SHA-256 hash of data.
func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func (b *Block) Serialize() []byte {
	panic("TODO: implement Serialize (copy from Task 01)")
}

func (b *Block) ComputeHash() []byte {
	panic("TODO: implement ComputeHash (copy from Task 01)")
}

func NewBlock(data []byte, prevBlockHash []byte) *Block {
	panic("TODO: implement NewBlock (copy from Task 01)")
}

// Blockchain links blocks together.
type Blockchain struct {
	Blocks []*Block
}

// NewGenesisBlock creates the first block in a chain, with PrevBlockHash
// set to 32 zero bytes.
func NewGenesisBlock() *Block {
	panic("TODO: implement NewGenesisBlock")
}

// NewBlockchain creates a new chain containing only a genesis block.
func NewBlockchain() *Blockchain {
	panic("TODO: implement NewBlockchain")
}

// AddBlock creates a new block containing data, links it to the current
// last block, and appends it to the chain.
func (bc *Blockchain) AddBlock(data string) {
	panic("TODO: implement AddBlock")
}

// IsValid walks the entire chain and returns false if any block's
// stored hash or link to its predecessor doesn't check out.
func (bc *Blockchain) IsValid() bool {
	panic("TODO: implement IsValid")
}
