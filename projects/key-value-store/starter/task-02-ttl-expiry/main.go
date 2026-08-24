package main

import (
	"sync"
	"time"
)

// Store is an in-memory key-value store with optional TTL expiry.
type Store struct {
	mu     sync.RWMutex
	data   map[string]string
	expiry map[string]time.Time
	done   chan struct{}
}

// NewStore creates a Store with background expiry cleanup every 100 ms.
func NewStore() *Store {
	panic("TODO: implement NewStore — init maps, start background cleanup goroutine")
}

// Close stops the background cleanup goroutine.
func (s *Store) Close() {
	panic("TODO: implement Close — signal the background goroutine to stop")
}

// Set stores value under key and removes any existing TTL on that key.
func (s *Store) Set(key, value string) {
	panic("TODO: implement Set")
}

// SetWithTTL stores value under key and schedules expiry after ttl.
func (s *Store) SetWithTTL(key, value string, ttl time.Duration) {
	panic("TODO: implement SetWithTTL")
}

// Get returns the value and true if key exists and has not expired.
// An expired key is treated identically to a missing key.
func (s *Store) Get(key string) (string, bool) {
	panic("TODO: implement Get — check expiry before returning")
}

// TTL returns the remaining lifetime of a key.
// Returns (remaining, true) if the key has an active TTL.
// Returns (0, false) if the key does not exist, has no TTL, or has expired.
func (s *Store) TTL(key string) (time.Duration, bool) {
	panic("TODO: implement TTL")
}

// Delete removes key. Returns true if the key existed and had not expired.
func (s *Store) Delete(key string) bool {
	panic("TODO: implement Delete")
}
