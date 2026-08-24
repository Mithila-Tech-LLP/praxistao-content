package main

import "time"

// TokenBucket is a deterministic, time-injected token-bucket rate limiter.
type TokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// NewTokenBucket returns a bucket that starts full (tokens = capacity) at
// time now.
func NewTokenBucket(capacity, refillRate float64, now time.Time) *TokenBucket {
	// TODO: return a *TokenBucket with tokens initialized to capacity and
	// lastRefill set to now. Never call time.Now() anywhere in this type.
	return &TokenBucket{capacity: capacity, refillRate: refillRate, lastRefill: now}
}

// Allow refills the bucket based on elapsed time since the last call, then
// attempts to consume one token.
func (b *TokenBucket) Allow(now time.Time) bool {
	// TODO: compute elapsed := now.Sub(b.lastRefill).
	// TODO: add elapsed.Seconds() * b.refillRate tokens, capped at b.capacity.
	// TODO: update b.lastRefill = now.
	// TODO: if b.tokens >= 1, subtract 1 and return true; otherwise return false.
	return false
}

func main() {}
