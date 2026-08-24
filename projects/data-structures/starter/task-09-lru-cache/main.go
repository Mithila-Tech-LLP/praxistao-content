package main

type lruNode struct {
	key, val   int
	prev, next *lruNode
}

type LRUCache struct {
	cap        int
	cache      map[int]*lruNode
	head, tail *lruNode // sentinel nodes
}

func NewLRUCache(capacity int) *LRUCache {
	// TODO: init head/tail sentinels, cache map
	return &LRUCache{}
}

func (c *LRUCache) Get(key int) (int, bool) {
	// TODO: return value if exists (and move to front), else (0, false)
	return 0, false
}

func (c *LRUCache) Put(key, value int) {
	// TODO: insert/update; if at capacity, evict LRU before inserting
}
