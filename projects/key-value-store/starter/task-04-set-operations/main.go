package main

import "sync"

// Store holds string keys and set keys.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
	sets map[string]map[string]struct{}
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

// --- set operations ---

// SAdd adds members to the set at key. Returns count of newly added members.
// Duplicate members are silently ignored.
func (s *Store) SAdd(key string, members ...string) int {
	panic("TODO: implement SAdd")
}

// SMembers returns all members of the set in arbitrary order.
// Returns an empty (non-nil) slice if key does not exist.
func (s *Store) SMembers(key string) []string {
	panic("TODO: implement SMembers")
}

// SIsMember returns true if member belongs to the set.
func (s *Store) SIsMember(key, member string) bool {
	panic("TODO: implement SIsMember")
}

// SRem removes members from the set. Returns count of actually removed members.
func (s *Store) SRem(key string, members ...string) int {
	panic("TODO: implement SRem")
}

// SCard returns the number of members in the set (cardinality).
// Returns 0 if key does not exist.
func (s *Store) SCard(key string) int {
	panic("TODO: implement SCard")
}
