package main

import "sync"

type Store struct {
	mu      sync.RWMutex
	strings map[string]string
	hashes  map[string]map[string]string
}

func NewStore() *Store {
	return &Store{
		strings: make(map[string]string),
		hashes:  make(map[string]map[string]string),
	}
}

func (s *Store) HSet(key, field, value string) {
	// TODO: set field in hash at key
}

func (s *Store) HGet(key, field string) (string, bool) {
	// TODO: get field from hash at key
	return "", false
}

func (s *Store) HGetAll(key string) map[string]string {
	// TODO: return copy of all fields for key
	return map[string]string{}
}

func (s *Store) HDel(key string, fields ...string) int {
	// TODO: delete fields, return count deleted
	return 0
}

func (s *Store) HExists(key, field string) bool {
	// TODO
	return false
}
