package main

import "sync"

// Store holds both string keys and list keys.
type Store struct {
	mu    sync.RWMutex
	data  map[string]string
	lists map[string][]string
}

// NewStore creates and returns an empty Store.
func NewStore() *Store {
	panic("TODO: implement NewStore")
}

// --- string operations (from Task 01) ---

func (s *Store) Set(key, value string) {
	panic("TODO: implement Set")
}

func (s *Store) Get(key string) (string, bool) {
	panic("TODO: implement Get")
}

func (s *Store) Delete(key string) bool {
	panic("TODO: implement Delete")
}

// --- list operations ---

// LPush prepends vals to the list at key. Each value is prepended individually,
// so LPush("k","a","b","c") results in ["c","b","a"].
// Returns the new length of the list.
func (s *Store) LPush(key string, vals ...string) int {
	panic("TODO: implement LPush")
}

// RPush appends vals to the list at key in order.
// Returns the new length of the list.
func (s *Store) RPush(key string, vals ...string) int {
	panic("TODO: implement RPush")
}

// LPop removes and returns the first element. Returns ("", false) if the list is empty.
func (s *Store) LPop(key string) (string, bool) {
	panic("TODO: implement LPop")
}

// RPop removes and returns the last element. Returns ("", false) if the list is empty.
func (s *Store) RPop(key string) (string, bool) {
	panic("TODO: implement RPop")
}

// LRange returns elements from index start to stop (inclusive).
// Negative indices count from the end: -1 is the last element.
// Out-of-bounds indices are clamped, not an error.
func (s *Store) LRange(key string, start, stop int) []string {
	panic("TODO: implement LRange")
}

// LLen returns the number of elements in the list, or 0 if the key does not exist.
func (s *Store) LLen(key string) int {
	panic("TODO: implement LLen")
}
