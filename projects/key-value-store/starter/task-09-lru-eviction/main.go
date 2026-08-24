package main

type lruEntry struct {
	key, val   string
	prev, next *lruEntry
}

type Store struct {
	cap        int
	data       map[string]*lruEntry
	head, tail *lruEntry
}

func NewStoreWithCapacity(maxKeys int) *Store {
	// TODO: init sentinels, set cap
	return &Store{}
}

func (s *Store) Get(key string) (string, bool) {
	// TODO: return value, move to front (most recently used)
	return "", false
}

func (s *Store) Set(key, value string) {
	// TODO: insert/update; if at capacity, evict LRU (tail side) first
}
