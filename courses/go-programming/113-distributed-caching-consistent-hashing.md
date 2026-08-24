# Chapter 113: Distributed Caching and Consistent Hashing

Chapter 80 taught you how to cache with a single Redis instance. That works until it doesn't: one instance runs out of memory, becomes a single point of failure, and turns into a bandwidth bottleneck when every service instance hammers the same node. The fix is to spread the cache across many nodes — and the moment you do that, a new question appears: *which node holds which key?* Answer it naively and adding one cache node invalidates almost your entire cache at once. Answer it with **consistent hashing** and adding a node moves only the keys it needs to own. This chapter covers distributed cache topologies, the caching patterns you need beyond cache-aside, invalidation, hot keys, and a from-scratch Go implementation of a consistent hash ring with virtual nodes.

## Table of Contents

1. [Cache Topologies — Local, Distributed, Layered](#1-cache-topologies--local-distributed-layered)
2. [Caching Patterns — Cache-Aside, Write-Through, Write-Behind](#2-caching-patterns--cache-aside-write-through-write-behind)
3. [Invalidation Strategies](#3-invalidation-strategies)
4. [Hot Keys](#4-hot-keys)
5. [Why Modulo Hashing Breaks](#5-why-modulo-hashing-breaks)
6. [Consistent Hashing from Scratch](#6-consistent-hashing-from-scratch)
7. [Virtual Nodes](#7-virtual-nodes)
8. [A Sharded Cache Client in Go](#8-a-sharded-cache-client-in-go)
9. [Consistent Hashing in the Wild](#9-consistent-hashing-in-the-wild)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Cache Topologies — Local, Distributed, Layered

There are three places a cache can live, and mature systems usually use more than one.

```
Local (in-process)              Distributed                     Layered (near cache)
┌──────────────┐            ┌──────────────┐                 ┌──────────────┐
│ App instance │            │ App instance │                 │ App instance │
│ ┌──────────┐ │            └──────┬───────┘                 │ ┌──────────┐ │
│ │ map/LRU  │ │                   │ network hop             │ │ L1 local │ │
│ └──────────┘ │            ┌──────▼───────┐                 │ └────┬─────┘ │
└──────────────┘            │ Redis node(s)│                 └──────┼───────┘
                            └──────────────┘                        │ on L1 miss
 nanosecond reads            shared by all                   ┌──────▼───────┐
 not shared, per-copy        instances, survives             │ L2 Redis     │
 duplication                 app restarts                    └──────────────┘
```

| Topology | Latency | Shared across instances? | Hard part |
|----------|---------|--------------------------|-----------|
| Local (in-process map, LRU from Ch 42) | ~100ns | No — each instance has its own copy | Invalidation across instances |
| Distributed (Redis/Memcached cluster) | ~0.5–2ms | Yes | Sharding keys across nodes |
| Layered (local L1 + distributed L2) | best of both | L2 yes, L1 no | Keeping L1 copies from going stale |

The rule of thumb:

- **Local cache** for small, hot, rarely-changing data: feature flags, config, country lists. Tolerate staleness with a short TTL (5–60s).
- **Distributed cache** for everything users see: sessions, rendered product pages, query results. One source of cached truth for all app instances.
- **Layered** when a distributed cache alone can't absorb the read rate — the local L1 shields the L2 from repeated reads of the same hot keys (more on hot keys in Section 4).

---

## 2. Caching Patterns — Cache-Aside, Write-Through, Write-Behind

Chapter 80 implemented cache-aside and write-through against a single Redis instance. Quick recap, then the one we haven't built yet.

| Pattern | Read path | Write path | Risk |
|---------|-----------|------------|------|
| **Cache-aside** | app checks cache → on miss, reads DB, fills cache | app writes DB, deletes/updates cache | brief staleness between DB write and invalidation |
| **Write-through** | same as cache-aside | app writes DB **and** cache synchronously | write latency = DB + cache |
| **Write-behind** (write-back) | reads hit cache | app writes **cache only**; a background worker flushes to DB in batches | data loss if cache dies before flush |

### Write-behind in Go

Write-behind trades durability for write throughput. It shines for high-frequency, low-value writes: view counters, "last seen" timestamps, telemetry. Never use it for money.

```go
package writebehind

import (
	"context"
	"log/slog"
	"time"
)

// Write is one pending database write.
type Write struct {
	Key   string
	Value int64
}

// Flusher persists a batch of writes to the database.
type Flusher interface {
	Flush(ctx context.Context, batch []Write) error
}

// Buffer accepts writes instantly and flushes them in batches.
type Buffer struct {
	in      chan Write
	flusher Flusher
	size    int           // flush when this many writes are pending
	every   time.Duration // ... or when this much time has passed
}

func NewBuffer(f Flusher, size int, every time.Duration) *Buffer {
	return &Buffer{
		in:      make(chan Write, 10_000),
		flusher: f,
		size:    size,
		every:   every,
	}
}

// Add records a write. It returns immediately — this is the whole point.
func (b *Buffer) Add(w Write) {
	select {
	case b.in <- w:
	default:
		// Buffer full: shed load rather than block the request path.
		slog.Warn("write-behind buffer full, dropping write", "key", w.Key)
	}
}

// Run flushes batches until ctx is cancelled. Call it in its own goroutine.
func (b *Buffer) Run(ctx context.Context) {
	ticker := time.NewTicker(b.every)
	defer ticker.Stop()

	batch := make([]Write, 0, b.size)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		// Coalesce: keep only the latest write per key.
		latest := make(map[string]int64, len(batch))
		for _, w := range batch {
			latest[w.Key] = w.Value
		}
		coalesced := make([]Write, 0, len(latest))
		for k, v := range latest {
			coalesced = append(coalesced, Write{Key: k, Value: v})
		}

		if err := b.flusher.Flush(ctx, coalesced); err != nil {
			slog.Error("write-behind flush failed", "err", err, "count", len(coalesced))
			return // keep batch? here we drop; a real system would retry
		}
		batch = batch[:0]
	}

	for {
		select {
		case w := <-b.in:
			batch = append(batch, w)
			if len(batch) >= b.size {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush() // final drain on shutdown (see Ch 134 for graceful shutdown)
			return
		}
	}
}
```

Note the **coalescing** step: if a counter was incremented 500 times in one second, the database sees one `UPDATE`, not 500. That's where the throughput win comes from — and also why a crash loses up to one batch of data.

---

## 3. Invalidation Strategies

> "There are only two hard things in Computer Science: cache invalidation and naming things." — Phil Karlton

Four strategies, from bluntest to sharpest:

### 3.1 TTL — let it expire

Every cache entry gets a lifetime. Simple, self-healing, and your safety net even when you also invalidate explicitly. The trade-off is a staleness window equal to the TTL.

```go
rdb.Set(ctx, "product:42", data, 5*time.Minute) // stale for at most 5 minutes
```

### 3.2 Explicit delete on write

The cache-aside classic: after writing the database, delete the cache key. Prefer **delete** over **update** — deleting is idempotent and immune to the race where two concurrent writers update the cache in the wrong order.

```go
func (s *ProductService) UpdatePrice(ctx context.Context, id int64, price int64) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE products SET price_cents = $1 WHERE id = $2`, price, id); err != nil {
		return err
	}
	// Delete, don't update: next reader re-fills from the source of truth.
	return s.rdb.Del(ctx, fmt.Sprintf("product:%d", id)).Err()
}
```

### 3.3 Versioned keys — invalidate by never invalidating

Instead of deleting entries, change the key. Keep a version counter per entity; readers include it in the cache key. Bumping the version orphans all old entries (they age out via TTL).

```go
// Read path: key includes the current version.
func (s *ProductService) cacheKey(ctx context.Context, id int64) (string, error) {
	ver, err := s.rdb.Get(ctx, fmt.Sprintf("product:%d:ver", id)).Int64()
	if err == redis.Nil {
		ver = 1
	} else if err != nil {
		return "", err
	}
	return fmt.Sprintf("product:%d:v%d", id, ver), nil
}

// Write path: one INCR invalidates every cached representation of the product —
// the JSON blob, the rendered HTML, the search snippet — in a single command.
func (s *ProductService) bumpVersion(ctx context.Context, id int64) error {
	return s.rdb.Incr(ctx, fmt.Sprintf("product:%d:ver", id)).Err()
}
```

This costs one extra `GET` per read (batch it with a pipeline) but makes "invalidate everything derived from this entity" a one-liner.

### 3.4 Broadcast invalidation — for local L1 caches

The hard case: layered caches. When instance A updates a product, instances B and C still hold the old value in their in-process L1. Redis Pub/Sub (Ch 80, Section 5) is the standard fix — every instance subscribes to an invalidation channel:

```go
// Every app instance runs this subscriber for its local L1 cache.
func (c *LayeredCache) listenForInvalidations(ctx context.Context) {
	sub := c.rdb.Subscribe(ctx, "cache:invalidate")
	defer sub.Close()

	for {
		select {
		case msg, ok := <-sub.Channel():
			if !ok {
				return
			}
			c.l1.Delete(msg.Payload) // drop the key from the local LRU
		case <-ctx.Done():
			return
		}
	}
}

// Writers publish after updating the database and deleting from L2.
func (c *LayeredCache) Invalidate(ctx context.Context, key string) error {
	if err := c.rdb.Del(ctx, key).Err(); err != nil { // L2
		return err
	}
	return c.rdb.Publish(ctx, "cache:invalidate", key).Err() // everyone's L1
}
```

Pub/Sub is fire-and-forget — an instance that was disconnected misses the message. That's why L1 entries **always** get a short TTL too: broadcast handles the common case fast, TTL guarantees an upper bound on staleness.

---

## 4. Hot Keys

Sharding spreads *keys* evenly. It does nothing about *traffic* per key. When a celebrity posts, or your homepage pins one product, a single key can absorb 100,000 reads/sec — and whichever cache node owns it melts while the others idle.

### Detecting hot keys

Count accesses cheaply with a sampled counter (count 1 in every 100 reads to keep overhead low), or use `redis-cli --hotkeys` against an LFU-configured Redis.

### Mitigation 1: request coalescing with singleflight

If the hot key ever misses (expiry, eviction), thousands of goroutines stampede to the database. Chapter 80 fixed this with a Redis `SETNX` fill-lock; **within one process**, `golang.org/x/sync/singleflight` is simpler — concurrent callers for the same key share one execution:

```go
import "golang.org/x/sync/singleflight"

type Cache struct {
	rdb   *redis.Client
	db    ProductDB
	group singleflight.Group
}

func (c *Cache) GetProduct(ctx context.Context, id int64) (*Product, error) {
	key := fmt.Sprintf("product:%d", id)

	if data, err := c.rdb.Get(ctx, key).Bytes(); err == nil {
		var p Product
		if json.Unmarshal(data, &p) == nil {
			return &p, nil
		}
	}

	// 10,000 concurrent misses for the same key → exactly one DB query.
	v, err, _ := c.group.Do(key, func() (any, error) {
		p, err := c.db.GetProduct(ctx, id)
		if err != nil {
			return nil, err
		}
		if data, err := json.Marshal(p); err == nil {
			c.rdb.Set(ctx, key, data, 5*time.Minute)
		}
		return p, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*Product), nil
}
```

### Mitigation 2: local L1 for the hottest keys

A 1-second local TTL on a key read 100,000 times/sec turns 100,000 network reads into ~1 per instance per second. Staleness bound: 1 second. Almost always acceptable for read-heavy hot data.

### Mitigation 3: key replication (splitting)

Store the hot key N times — `product:42#0` … `product:42#7` — and have each reader pick a random replica. The load spreads across up to N cache nodes. The cost: writers must update all N copies, so reserve this for keys you've *measured* to be hot.

```go
func hotKey(base string, replicas int) string {
	return fmt.Sprintf("%s#%d", base, rand.IntN(replicas)) // math/rand/v2
}
```

---

## 5. Why Modulo Hashing Breaks

Say you shard keys across 3 cache nodes the obvious way:

```go
node := nodes[crc32.ChecksumIEEE([]byte(key)) % uint32(len(nodes))]
```

This works — until the node count changes. Watch what happens to 10 keys when you go from 3 nodes to 4:

```
key    hash%3 (before)   hash%4 (after)   moved?
k1        0                 1              YES
k2        1                 1              no
k3        2                 3              YES
k4        0                 0              no
k5        1                 2              YES
k6        2                 1              YES
k7        0                 3              YES
k8        1                 0              YES
k9        2                 2              no
k10       0                 1              YES
```

Roughly **3 out of 4 keys** (in general, `1 - 1/N` of them) land on a different node. Every moved key is now a cache miss. Adding one node to relieve pressure just invalidated ~75% of your cache — the database gets hit by the full miss storm at the exact moment you were scaling because load was already high. This is how "we added a cache node" becomes an outage.

What we want instead: adding an Nth node moves only ~`1/N` of keys — the ones the new node should own — and nothing else.

---

## 6. Consistent Hashing from Scratch

The idea, in three steps:

1. Imagine the whole hash space (0 to 2³²−1 for a 32-bit hash) bent into a **circle**: value 2³²−1 wraps around to 0.
2. **Hash each node's name** onto the circle. Each node "sits" at its hash position.
3. To place a key: hash it onto the same circle, then **walk clockwise** until you meet the first node. That node owns the key.

```
                    0 / 2³²
                      │
             k4 ●     ▼      ● k1
                 ┌─────────┐
        node C ■ │         │ ■ node A
                 │  hash   │
                 │  ring   │      ● k2
             k3 ●│         │
                 └─────────┘
                      ■ node B
                      ▲
                      ● (k1, k2 walk clockwise → node B owns them)

■ = node position (hash of node name)     ● = key position (hash of key)
Each key belongs to the first node found moving clockwise from it.
```

Now the magic:

- **Add node D** between A and B → D takes over only the arc between A and D. Keys on that arc move from B to D. **Every other key stays put.**
- **Remove node B** → B's arc merges into the next node clockwise. Only B's keys move.

On average, adding a node to an N-node ring relocates `1/(N+1)` of keys — the theoretical minimum.

### The implementation

"Walk clockwise" translates to: keep node hashes in a sorted slice, and binary-search for the first hash ≥ the key's hash (wrapping to index 0 past the end). Go's `sort.Search` does exactly this.

```go
// Package ring implements a consistent hash ring with virtual nodes.
package ring

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Ring struct {
	mu       sync.RWMutex
	replicas int               // virtual nodes per real node (Section 7)
	hashes   []uint32          // sorted positions on the ring
	owner    map[uint32]string // position → real node name
	nodes    map[string]struct{}
}

// New creates a ring. replicas is the number of virtual nodes per real node;
// use 1 to see the problem virtual nodes solve, 100–200 in practice.
func New(replicas int) *Ring {
	return &Ring{
		replicas: replicas,
		owner:    make(map[uint32]string),
		nodes:    make(map[string]struct{}),
	}
}

func (r *Ring) hash(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// Add places a node (and its virtual replicas) on the ring.
func (r *Ring) Add(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodes[node]; ok {
		return
	}
	r.nodes[node] = struct{}{}

	for i := range r.replicas { // Go 1.22: range over int
		h := r.hash(node + "#" + strconv.Itoa(i))
		r.owner[h] = node
		r.hashes = append(r.hashes, h)
	}
	sort.Slice(r.hashes, func(i, j int) bool { return r.hashes[i] < r.hashes[j] })
}

// Remove takes a node and all its virtual replicas off the ring.
func (r *Ring) Remove(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.nodes[node]; !ok {
		return
	}
	delete(r.nodes, node)

	for i := range r.replicas {
		h := r.hash(node + "#" + strconv.Itoa(i))
		delete(r.owner, h)
	}
	kept := r.hashes[:0]
	for _, h := range r.hashes {
		if _, ok := r.owner[h]; ok {
			kept = append(kept, h)
		}
	}
	r.hashes = kept
}

// Get returns the node that owns key, walking clockwise from the key's hash.
func (r *Ring) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.hashes) == 0 {
		return "", false
	}

	h := r.hash(key)
	// First position on the ring >= h ...
	idx := sort.Search(len(r.hashes), func(i int) bool { return r.hashes[i] >= h })
	if idx == len(r.hashes) {
		idx = 0 // ... wrapping around past the top of the ring
	}
	return r.owner[r.hashes[idx]], true
}
```

(One honest caveat: two virtual node names could CRC32-collide, in which case the later `Add` silently overwrites the position's owner. With ~hundreds of virtual nodes the odds are tiny and the consequence is a slightly uneven split; production libraries use 64-bit hashes to make it negligible.)

### Verifying minimal movement

```go
package main

import (
	"fmt"
	"strconv"

	"example.com/cache/ring"
)

func main() {
	const keys = 100_000

	before := ring.New(100)
	for _, n := range []string{"cache-a", "cache-b", "cache-c"} {
		before.Add(n)
	}

	after := ring.New(100)
	for _, n := range []string{"cache-a", "cache-b", "cache-c", "cache-d"} {
		after.Add(n)
	}

	moved := 0
	for i := range keys {
		k := "key-" + strconv.Itoa(i)
		b, _ := before.Get(k)
		a, _ := after.Get(k)
		if a != b {
			moved++
		}
	}
	fmt.Printf("moved %d of %d keys (%.1f%%)\n", moved, keys, 100*float64(moved)/keys)
}
```

```
moved 24923 of 100000 keys (24.9%)
```

Adding a 4th node moved ~25% of keys — exactly the `1/4` the new node must own. Modulo hashing would have moved ~75%.

---

## 7. Virtual Nodes

Try the ring with `New(1)` — one position per real node — and count key ownership:

```
cache-a: 61,483 keys    ← 6x more load than cache-c!
cache-b: 28,006 keys
cache-c: 10,511 keys
```

With only 3 random positions on the circle, the arcs between them are wildly unequal — that's just what a few random points on a circle look like. Whichever node owns the longest arc gets the most keys. Worse, when a node dies, its **entire** arc dumps onto the *single* next node clockwise, potentially cascading it into overload too.

**Virtual nodes** fix both problems. Each real node appears at many positions — `cache-a#0`, `cache-a#1`, …, `cache-a#99` — so each real node owns many *small* arcs scattered around the ring:

```
1 virtual node each:              100 virtual nodes each:

   ┌───AAAAAAAAAAAAA──┐              ┌─ab─c─a─cc─b─a─b─┐
   │                  │              │                  │
   C                  A              c                  a
   │                  │              │                  │
   └──CCC───BBBBBBB───┘              └─ca─b─a─c─b─ab─c──┘

  arcs are huge and uneven         arcs are small and interleaved
  node death → 1 neighbor          node death → load spreads over
  absorbs everything               ALL remaining nodes
```

Re-running the distribution test with `New(100)`:

```
cache-a: 33,842 keys
cache-b: 32,910 keys
cache-c: 33,248 keys     ← within ~2% of even
```

And when a node is removed, its 100 little arcs are absorbed by whichever node happens to be next after each one — statistically, *all* remaining nodes share the load.

How many virtual nodes? The standard deviation of load shrinks roughly with `1/√(replicas)`. 100–200 per node gives ~5–10% imbalance and is the common default (Amazon's Dynamo paper used 100–200; groupcache defaults to 50). More replicas cost memory and slightly slower `Add`/`Remove` — `Get` stays `O(log(nodes × replicas))`.

Virtual nodes have a bonus feature: **weighting**. Give a machine with double the RAM double the virtual nodes, and it will own roughly double the keys.

---

## 8. A Sharded Cache Client in Go

Now put the ring in front of several independent Redis instances (this is exactly how Memcached clients work — the servers don't even know about each other; all the intelligence lives in the client):

```go
package shardedcache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"example.com/cache/ring"
)

var ErrNoNodes = errors.New("shardedcache: no cache nodes available")

type Cache struct {
	ring    *ring.Ring
	clients map[string]*redis.Client
}

// New connects one Redis client per node. addrs maps node name → host:port.
func New(addrs map[string]string) *Cache {
	c := &Cache{
		ring:    ring.New(150),
		clients: make(map[string]*redis.Client, len(addrs)),
	}
	for name, addr := range addrs {
		c.ring.Add(name)
		c.clients[name] = redis.NewClient(&redis.Options{Addr: addr})
	}
	return c
}

func (c *Cache) clientFor(key string) (*redis.Client, error) {
	node, ok := c.ring.Get(key)
	if !ok {
		return nil, ErrNoNodes
	}
	return c.clients[node], nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	rdb, err := c.clientFor(key)
	if err != nil {
		return nil, err
	}
	return rdb.Get(ctx, key).Bytes()
}

func (c *Cache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	rdb, err := c.clientFor(key)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, val, ttl).Err()
}

// RemoveNode drops a dead node from the ring. Keys it owned now route to the
// next node clockwise — they'll miss once and re-fill from the database.
func (c *Cache) RemoveNode(name string) {
	c.ring.Remove(name)
	if rdb, ok := c.clients[name]; ok {
		rdb.Close()
		delete(c.clients, name)
	}
}
```

Usage:

```go
cache := shardedcache.New(map[string]string{
	"cache-a": "10.0.1.10:6379",
	"cache-b": "10.0.1.11:6379",
	"cache-c": "10.0.1.12:6379",
})

cache.Set(ctx, "product:42", data, 5*time.Minute) // → lands on exactly one node
val, err := cache.Get(ctx, "product:42")          // → same node, guaranteed
```

Two operational notes:

- **Failure handling**: because a cache is allowed to lose data, node failure is cheap — remove it from the ring, absorb the misses, done. This is much easier than the database sharding in Chapter 116, where a lost node means lost data and the ring must be paired with replication.
- **Where does the node list come from?** Hardcoded config to start; service discovery (Chapter 110) once nodes come and go dynamically. Every client must build the ring from the *same* node list and the *same* hash function, or they'll disagree about ownership and your hit rate quietly craters.

---

## 9. Consistent Hashing in the Wild

You will rarely maintain a hand-rolled ring in production, but you'll meet the idea everywhere:

| System | How it uses the idea |
|--------|----------------------|
| **Memcached clients** | Classic client-side ring, exactly like Section 8 (ketama algorithm) |
| **Redis Cluster** | A variation: 16,384 fixed **hash slots** (`CRC16(key) % 16384`), slots assigned to nodes; moving a slot moves its keys. Same goal, coarser granularity |
| **groupcache** (Go, by Brad Fitzpatrick) | Peers form a ring; each key has one owner peer that fills it via singleflight — hot-key protection built in |
| **DynamoDB / Cassandra** | Data partitioning across storage nodes with virtual nodes ("vnodes") |
| **Kafka consumers, load balancers** | Sticky assignment of partitions/clients to workers with minimal reshuffling |

And you've already seen it pay off in this course's neighborhood: Chapter 116 reuses this exact ring to route *database* shards, where minimal key movement matters even more because moving a key means physically migrating rows.

---

## Summary

- **Topologies**: local caches are fastest but per-instance; distributed caches are shared truth; layered (L1 + L2) shields the distributed tier from hot-read traffic.
- **Patterns**: cache-aside and write-through (Ch 80) cover most needs; **write-behind** batches and coalesces writes for huge throughput at the cost of possible data loss — counters yes, money no.
- **Invalidation**: always set a TTL as the safety net; prefer *delete* over *update* on writes; versioned keys make multi-representation invalidation atomic; Pub/Sub broadcasts invalidations to local L1 caches.
- **Hot keys**: sharding balances keys, not traffic. Coalesce misses with `singleflight`, absorb reads with a short-TTL local cache, split measured hot keys across replicas.
- **Modulo hashing** remaps `1 − 1/N` of keys when the node count changes — a self-inflicted cache flush. **Consistent hashing** remaps only `1/N`: nodes sit on a hash ring, each key belongs to the first node clockwise, and `sort.Search` finds it in `O(log n)`.
- **Virtual nodes** (100–200 per real node) fix the uneven-arc problem, spread a dead node's load across all survivors, and enable capacity weighting.
- Redis Cluster, Memcached clients, DynamoDB, Cassandra, and groupcache are all variations on this one idea.

## Exercises

### Easy
1. Implement the modulo-hashing experiment from Section 5 in code: assign 10,000 keys to 3 nodes with `hash % 3`, then to 4 nodes with `hash % 4`, and print the percentage of keys that moved. Repeat going from 9 to 10 nodes — does the damage get better or worse as N grows?
2. Using the `Ring` from Section 6, add nodes `a`, `b`, `c` with 150 virtual nodes each and distribute 100,000 keys. Print each node's key count and the max/min ratio. Re-run with 1, 10, and 1,000 virtual nodes and tabulate how the imbalance shrinks.
3. Write a table test for `Ring.Get` that verifies: (a) the same key always returns the same node, (b) `Get` on an empty ring returns `ok == false`, and (c) after `Remove("b")`, no key maps to `b`.

### Medium
4. Extend `Ring` with weights: `AddWeighted(node string, weight int)` gives a node `weight × replicas` virtual nodes. Add one node with weight 2 and two with weight 1, distribute 100,000 keys, and verify the heavy node owns ~50%.
5. Build a `LayeredCache` combining an in-process LRU (reuse your Chapter 42 implementation) as L1 with Redis as L2, including the Pub/Sub broadcast invalidation from Section 3.4. Write a test with two `LayeredCache` instances sharing one Redis: update through instance A and assert instance B's L1 entry is gone within 100ms.
6. Implement `GetReplicated(key string, n int) []string` on the ring: return the first `n` **distinct real nodes** walking clockwise (skip virtual nodes of already-chosen reals). This is the primitive replicated caches and Dynamo-style stores use to place N copies of each key.

### Hard
7. Build a **hot-key detector and auto-splitter**: wrap the sharded cache client with a sampled per-key counter (count every 100th access). When a key exceeds a threshold rate, transparently switch it to split mode — writes fan out to `key#0..key#7`, reads pick a random replica — and demote it when traffic drops. Prove with a load test (1 goroutine per 1,000 req/s) that the hottest node's load drops after the split kicks in.
8. Implement **bounded-load consistent hashing** (Google's "Consistent Hashing with Bounded Loads"): track how many keys each node currently owns; on `Get`, if the clockwise owner is already above `⌈avg × 1.25⌉` keys, continue clockwise to the next node with spare capacity. Compare the max/avg load ratio against the plain ring under a skewed (Zipfian) key distribution.
