package main

import "sync"

type ZMember struct {
	Member string
	Score  float64
}

type Store struct {
	mu   sync.RWMutex
	sets map[string][]ZMember // kept sorted by score ascending
}

func NewStore() *Store {
	return &Store{sets: make(map[string][]ZMember)}
}

func (s *Store) ZAdd(key string, score float64, member string) {
	// TODO: insert or update member; keep slice sorted by score ascending
}

func (s *Store) ZScore(key, member string) (float64, bool) {
	// TODO
	return 0, false
}

func (s *Store) ZRange(key string, start, stop int) []string {
	// TODO: return member names at indices start..stop inclusive (-1 means last)
	return []string{}
}

func (s *Store) ZRangeWithScores(key string, start, stop int) []ZMember {
	// TODO
	return []ZMember{}
}

func (s *Store) ZRank(key, member string) (int, bool) {
	// TODO: return 0-based rank (position) of member
	return 0, false
}
