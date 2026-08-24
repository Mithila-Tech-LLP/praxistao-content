package main

const numBuckets = 16

type entry struct {
	key string
	val int
}

type HashMap struct {
	buckets [numBuckets][]entry
}

func (m *HashMap) Set(key string, val int) {
	// TODO: implement — hash key, find/replace in bucket or append
}

func (m *HashMap) Get(key string) (int, bool) {
	// TODO: implement
	return 0, false
}

func (m *HashMap) Delete(key string) {
	// TODO: implement
}

func (m *HashMap) Keys() []string {
	// TODO: implement — return all keys in any order
	return []string{}
}
