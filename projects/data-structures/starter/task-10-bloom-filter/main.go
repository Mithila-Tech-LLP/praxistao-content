package main

type BloomFilter struct {
	bits      []bool
	size      int
	numHashes int
}

func NewBloomFilter(size, numHashes int) *BloomFilter {
	return &BloomFilter{bits: make([]bool, size), size: size, numHashes: numHashes}
}

func (b *BloomFilter) Add(item string) {
	// TODO: for each hash seed 0..numHashes-1, compute hash(item+seed) % size, set bit
}

func (b *BloomFilter) MightContain(item string) bool {
	// TODO: return true only if ALL hash positions are set
	return false
}
