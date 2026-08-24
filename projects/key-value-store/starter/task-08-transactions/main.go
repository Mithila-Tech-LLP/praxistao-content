package main

import "sync"

type Store struct {
	mu   sync.Mutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Get(key string) (string, bool) { v, ok := s.data[key]; return v, ok }
func (s *Store) Set(key, value string)         { s.data[key] = value }

type cmd struct{ op, key, val string }

type Tx struct {
	store *Store
	queue []cmd
}

func (s *Store) Begin() *Tx {
	return &Tx{store: s}
}

func (t *Tx) Set(key, value string) {
	// TODO: queue a set command
}

func (t *Tx) Delete(key string) {
	// TODO: queue a delete command
}

func (t *Tx) Exec() error {
	// TODO: execute all queued commands atomically (hold store lock)
	return nil
}

func (t *Tx) Discard() {
	// TODO: clear the queue
}
