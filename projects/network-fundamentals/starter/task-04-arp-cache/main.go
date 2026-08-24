package main

// ARPCache is a bounded IP -> MAC cache with LRU eviction.
type ARPCache struct {
	capacity int
	// TODO: add fields to track entries and recency order, e.g. a
	// container/list.List for recency and a map[string]*list.Element for
	// O(1) lookup.
}

// NewARPCache returns an empty cache bounded to capacity entries.
func NewARPCache(capacity int) *ARPCache {
	// TODO: initialize and return an *ARPCache with the given capacity.
	return &ARPCache{capacity: capacity}
}

// Set stores or updates the MAC for ip, evicting the least-recently-used
// entry if the cache is full.
func (c *ARPCache) Set(ip, mac string) {
	// TODO: if ip already exists, update its MAC and mark it most-recently-used.
	// TODO: otherwise, if the cache is at capacity, evict the
	// least-recently-used entry before inserting the new one.
	// TODO: insert the new entry as most-recently-used.
}

// Lookup returns the MAC for ip, refreshing it as most-recently-used.
func (c *ARPCache) Lookup(ip string) (mac string, ok bool) {
	// TODO: look up ip; if found, mark it most-recently-used and return its MAC.
	// TODO: return ok=false if not present.
	return "", false
}

func main() {}
