package main

import "sync"

// Store is an in-memory key-value store safe for concurrent use.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewStore creates and returns an empty Store.
func NewStore() *Store {
	panic("TODO: implement NewStore")
}

// Set stores value under key, overwriting any previous value.
func (s *Store) Set(key, value string) {
	panic("TODO: implement Set")
}

// Get returns the value and true if key exists, or ("", false) if not.
func (s *Store) Get(key string) (string, bool) {
	panic("TODO: implement Get")
}

// Delete removes key from the store. Returns true if the key existed.
func (s *Store) Delete(key string) bool {
	panic("TODO: implement Delete")
}

// Exists returns true if key is present in the store.
func (s *Store) Exists(key string) bool {
	panic("TODO: implement Exists")
}
