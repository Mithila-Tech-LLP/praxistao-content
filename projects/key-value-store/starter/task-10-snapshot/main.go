package main

import (
	"encoding/json"
	"io"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) { s.mu.Lock(); s.data[key] = value; s.mu.Unlock() }
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Save(w io.Writer) error {
	// TODO: serialize s.data as JSON to w
	return nil
}

func (s *Store) Load(r io.Reader) error {
	// TODO: decode JSON from r into s.data, replacing current data
	return nil
}

var _ = json.Marshal
