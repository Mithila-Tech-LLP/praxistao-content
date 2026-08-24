package main

import (
	"testing"
	"time"
)

func TestTokenBucket_BasicConsumptionAndRefill(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	b := NewTokenBucket(2, 1, t0)

	if !b.Allow(t0) {
		t.Error("Allow(t0) #1: expected true (bucket starts full)")
	}
	if !b.Allow(t0) {
		t.Error("Allow(t0) #2: expected true (one token left)")
	}
	if b.Allow(t0) {
		t.Error("Allow(t0) #3: expected false (no tokens, no time elapsed)")
	}
	if !b.Allow(t0.Add(2 * time.Second)) {
		t.Error("Allow(t0+2s): expected true (2 tokens refilled)")
	}
}
