package main

import (
	"fmt"
	"testing"
)

func TestBloomFilter_AddAndMightContain(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	items := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
	for _, item := range items {
		bf.Add(item)
	}

	for _, item := range items {
		if !bf.MightContain(item) {
			t.Errorf("MightContain(%q) should return true after Add — bloom filter must not have false negatives", item)
		}
	}
}

func TestBloomFilter_FalseNegativeFree(t *testing.T) {
	// A bloom filter MUST NOT produce false negatives: if we added X, MightContain(X) must be true.
	bf := NewBloomFilter(1000, 3)
	words := []string{"go", "golang", "gopher", "goroutine", "garbage", "collector"}
	for _, w := range words {
		bf.Add(w)
	}
	for _, w := range words {
		if !bf.MightContain(w) {
			t.Errorf("False negative: MightContain(%q) = false, but it was added", w)
		}
	}
}

func TestBloomFilter_FalsePositiveRateReasonable(t *testing.T) {
	// With size=1000 and 3 hashes, false positive rate should stay well below 50%
	// for items that were never added. We test 100 unique non-added items.
	bf := NewBloomFilter(1000, 3)
	for i := 0; i < 20; i++ {
		bf.Add(fmt.Sprintf("added-item-%d", i))
	}

	falsePositives := 0
	total := 100
	for i := 0; i < total; i++ {
		// Use a prefix that does not overlap with "added-item-*"
		if bf.MightContain(fmt.Sprintf("not-in-filter-%d", i)) {
			falsePositives++
		}
	}

	// False positive rate must be below 50% for a well-sized filter
	if falsePositives > total/2 {
		t.Errorf("False positive rate too high: %d/%d — filter may be broken", falsePositives, total)
	}
}

func TestBloomFilter_EmptyFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	// Nothing added — MightContain should return false for any item
	// (A correct empty filter has all bits 0, so all hashes miss.)
	if bf.MightContain("anything") {
		t.Error("MightContain on empty filter should return false")
	}
}

func TestBloomFilter_SingleItem(t *testing.T) {
	bf := NewBloomFilter(1000, 3)
	bf.Add("only-item")

	if !bf.MightContain("only-item") {
		t.Error("MightContain should return true for the only added item")
	}
}

func TestBloomFilter_DifferentSizesAndHashes(t *testing.T) {
	configs := []struct {
		size      int
		numHashes int
	}{
		{100, 1},
		{500, 5},
		{2000, 4},
	}

	for _, cfg := range configs {
		bf := NewBloomFilter(cfg.size, cfg.numHashes)
		bf.Add("hello")
		bf.Add("world")

		if !bf.MightContain("hello") {
			t.Errorf("size=%d hashes=%d: MightContain(hello) should be true", cfg.size, cfg.numHashes)
		}
		if !bf.MightContain("world") {
			t.Errorf("size=%d hashes=%d: MightContain(world) should be true", cfg.size, cfg.numHashes)
		}
	}
}
