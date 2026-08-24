package main

// Block is one entry in the chain.
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// Hash returns the SHA-256 hash of data.
func Hash(data []byte) []byte {
	panic("TODO: implement Hash using crypto/sha256")
}

// Serialize returns a deterministic byte representation of the block,
// suitable for hashing. It does NOT include the block's own Hash field.
func (b *Block) Serialize() []byte {
	panic("TODO: implement Serialize")
}

// ComputeHash computes and returns this block's hash from its other fields.
func (b *Block) ComputeHash() []byte {
	panic("TODO: implement ComputeHash")
}

// NewBlock creates a block with the given data and previous hash, and
// sets its Hash field by calling ComputeHash.
func NewBlock(data []byte, prevBlockHash []byte) *Block {
	panic("TODO: implement NewBlock")
}
