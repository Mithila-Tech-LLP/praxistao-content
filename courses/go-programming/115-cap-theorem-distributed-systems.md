# Chapter 97: CAP Theorem, Consistency, and Distributed Systems

When you have multiple database nodes, something has to give under network partition. CAP theorem formalizes what. Understanding consistency models tells you what guarantees you actually get from databases like PostgreSQL, MongoDB, Cassandra, and Redis.

## Table of Contents

1. [CAP Theorem](#1-cap-theorem)
2. [Consistency Models](#2-consistency-models)
3. [Eventual Consistency in Practice](#3-eventual-consistency-in-practice)
4. [Distributed Locks Revisited](#4-distributed-locks-revisited)
5. [Sharding and Partitioning](#5-sharding-and-partitioning)
6. [PACELC — Beyond CAP](#6-pacelc--beyond-cap)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. CAP Theorem

A distributed system can guarantee at most 2 of 3:

```
C - Consistency:    Every read sees the most recent write (or an error)
A - Availability:   Every request gets a response (not an error)
P - Partition Tolerance: System works even when nodes can't communicate
```

Partition tolerance is non-negotiable in practice (networks drop packets). So the real choice is **CP vs AP** during a partition.

| Database | Choice | Behavior during partition |
|----------|--------|--------------------------|
| PostgreSQL (primary) | CP | Reads from replica may be stale; writes go to primary only |
| Cassandra | AP | Continues accepting writes/reads; may return stale data |
| MongoDB (replica set) | CP | Primary election; brief unavailability |
| Redis Cluster | AP (configurable) | Can return stale data from remaining nodes |
| Zookeeper/etcd | CP | Rejects operations until quorum restored |

### What CP means in practice

```go
// PostgreSQL: when you read from a replica, you may get old data
// This is fine for many reads but not for cases where consistency matters

// Safe: always read from primary for consistency-critical paths
db.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1", id)

// Stale: reading from a replica (may lag by 10ms-minutes)
replicaDB.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id=$1", id)

// Better: use read-your-writes consistency
// After a write, read back from primary (or include a min_lsn hint)
```

---

## 2. Consistency Models

From strongest to weakest:

```
Linearizability (strongest):
  Every operation appears to take effect instantaneously at some point
  between its invocation and completion.
  Example: distributed compare-and-swap (etcd)

Sequential Consistency:
  All operations appear to occur in some sequential order consistent
  with each process's program order.
  Example: single-master database with synchronous replication

Causal Consistency:
  Causally related operations appear in the same order on all nodes.
  Concurrent operations may differ.
  Example: MongoDB causal sessions

Monotonic Reads:
  Once you read a value, subsequent reads return the same or newer value.
  (no "going back in time")

Read-Your-Writes:
  After writing a value, you always see your own write.
  Other readers may still see the old value.
  Example: sticky sessions to the same DB node

Eventual Consistency (weakest):
  If no new updates are made, all replicas eventually converge to the same value.
  Example: DNS, Cassandra with eventual consistency
```

---

## 3. Eventual Consistency in Practice

```go
// Conflict resolution: last-write-wins (LWW) using timestamps
type VersionedValue struct {
    Value     string
    Timestamp int64 // unix nanoseconds
    NodeID    string
}

func merge(a, b VersionedValue) VersionedValue {
    if a.Timestamp > b.Timestamp { return a }
    if b.Timestamp > a.Timestamp { return b }
    // Same timestamp: use node ID as tiebreaker (deterministic but arbitrary)
    if a.NodeID > b.NodeID { return a }
    return b
}

// Vector clocks: track causality across nodes
type VectorClock map[string]int

func (vc VectorClock) Increment(nodeID string) VectorClock {
    result := make(VectorClock, len(vc))
    for k, v := range vc { result[k] = v }
    result[nodeID]++
    return result
}

func (vc VectorClock) HappensBefore(other VectorClock) bool {
    for nodeID, ts := range vc {
        if ts > other[nodeID] { return false }
    }
    return true
}

func (vc VectorClock) Concurrent(other VectorClock) bool {
    return !vc.HappensBefore(other) && !other.HappensBefore(vc)
}

// CRDTs: data structures that merge conflict-free
// Counter CRDT (G-Counter): each node has its own counter; total = sum of all
type GCounter struct {
    counts map[string]int // nodeID → increment count
}

func (g *GCounter) Increment(nodeID string) { g.counts[nodeID]++ }
func (g *GCounter) Value() int {
    total := 0
    for _, v := range g.counts { total += v }
    return total
}
func (g *GCounter) Merge(other GCounter) {
    for nodeID, v := range other.counts {
        if v > g.counts[nodeID] { g.counts[nodeID] = v }
    }
}
```

### Read-Repair pattern

```go
// Read from multiple replicas; if they disagree, repair the stale one
func readWithRepair(ctx context.Context, replicas []DB, key string) (string, error) {
    type result struct{ value string; version int64; db DB }
    
    results := make([]result, 0, len(replicas))
    for _, db := range replicas {
        v, ver, err := db.GetWithVersion(ctx, key)
        if err == nil { results = append(results, result{v, ver, db}) }
    }
    if len(results) == 0 { return "", errors.New("all replicas unavailable") }
    
    // Find the most recent value
    best := results[0]
    for _, r := range results[1:] {
        if r.version > best.version { best = r }
    }
    
    // Repair stale replicas
    for _, r := range results {
        if r.version < best.version {
            go r.db.SetWithVersion(ctx, key, best.value, best.version)
        }
    }
    
    return best.value, nil
}
```

---

## 4. Distributed Locks Revisited

Distributed locks are hard. Redis `SETNX` works for many cases but has edge cases:

```go
// The safe way: use RedLock algorithm for safety, or etcd for strong consistency
import clientv3 "go.etcd.io/etcd/client/v3"
import "go.etcd.io/etcd/client/v3/concurrency"

func acquireEtcdLock(ctx context.Context, client *clientv3.Client, key string, ttl int) (*concurrency.Mutex, error) {
    session, err := concurrency.NewSession(client, concurrency.WithTTL(ttl))
    if err != nil { return nil, err }
    
    mu := concurrency.NewMutex(session, "/locks/"+key)
    if err := mu.Lock(ctx); err != nil {
        session.Close()
        return nil, err
    }
    return mu, nil
}

func withLock(ctx context.Context, client *clientv3.Client, key string, fn func() error) error {
    mu, err := acquireEtcdLock(ctx, client, key, 30)
    if err != nil { return fmt.Errorf("acquire lock: %w", err) }
    defer mu.Unlock(ctx)
    return fn()
}
```

---

## 5. Sharding and Partitioning

Sharding splits data across nodes for horizontal scaling.

```go
// Hash-based sharding: consistent hashing
type ConsistentHashRing struct {
    replicas int
    ring     map[int]string  // hash → node name
    nodes    []int           // sorted hash values
    mu       sync.RWMutex
}

func NewConsistentHashRing(replicas int) *ConsistentHashRing {
    return &ConsistentHashRing{replicas: replicas, ring: make(map[int]string)}
}

func (r *ConsistentHashRing) AddNode(node string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    for i := range r.replicas {
        hash := int(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s-%d", node, i))))
        r.ring[hash] = node
        r.nodes = append(r.nodes, hash)
    }
    sort.Ints(r.nodes)
}

func (r *ConsistentHashRing) GetNode(key string) string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    if len(r.nodes) == 0 { return "" }
    
    hash := int(crc32.ChecksumIEEE([]byte(key)))
    
    // Find first node with hash >= key hash (clockwise)
    idx := sort.SearchInts(r.nodes, hash)
    if idx >= len(r.nodes) { idx = 0 } // wrap around
    
    return r.ring[r.nodes[idx]]
}

// Range-based sharding: user IDs 0-999 → shard 0, 1000-1999 → shard 1, etc.
func getUserShard(userID int64, numShards int) int {
    return int(userID) % numShards
}

// Directory-based sharding: look up which shard holds a given key
type ShardDirectory struct {
    mu     sync.RWMutex
    lookup map[string]int // key prefix → shard ID
}

func (d *ShardDirectory) GetShard(key string) int {
    d.mu.RLock()
    defer d.mu.RUnlock()
    for prefix, shard := range d.lookup {
        if strings.HasPrefix(key, prefix) { return shard }
    }
    return 0
}
```

### Resharding

```go
// When adding a new shard: move a portion of data from existing shards
// Classic pattern: dual-write then migrate

func migrateShard(ctx context.Context, src, dst DB, keyRange [2]string) error {
    cursor := src.Scan(ctx, keyRange[0], keyRange[1])
    var batch []KV
    
    for cursor.Next() {
        batch = append(batch, cursor.Value())
        if len(batch) >= 100 {
            if err := dst.BatchSet(ctx, batch); err != nil { return err }
            batch = batch[:0]
        }
    }
    if len(batch) > 0 { return dst.BatchSet(ctx, batch) }
    return cursor.Err()
}
```

---

## 6. PACELC — Beyond CAP

CAP only describes behavior during partitions (P). PACELC extends this: even when there is no partition (E), you trade off latency (L) vs consistency (C).

```
PACELC:
  During Partition:   Choose between A (availability) and C (consistency)
  Else (no partition): Choose between L (latency) and C (consistency)

PostgreSQL:     PC/EC  — CP during partition, consistent reads (higher latency)
Cassandra:      PA/EL  — AP during partition, low latency (eventual consistency) 
Dynamo/DynamoDB: PA/EL — same as Cassandra
HBase:          PC/EC  — strong consistency, higher latency
Zookeeper:      PC/EC  — CP, used for coordination
```

---

## Summary

- **CAP**: in real distributed systems, choose between CP (prefer consistency) or AP (prefer availability) during network partitions
- **Consistency models**: linearizability → sequential → causal → read-your-writes → eventual; choose weakest model that satisfies your requirements
- **Eventual consistency**: requires conflict resolution (LWW, vector clocks, CRDTs)
- **Distributed locks**: Redis SETNX for most cases; etcd for strong consistency guarantees
- **Sharding**: consistent hashing for load distribution; be careful about cross-shard queries and resharding complexity

## Exercises

### Easy
1. Implement a `GSet` (Grow-only Set) CRDT: you can add elements but never remove them. The merge operation is set union. Show that concurrent adds on two nodes always converge after merging.
2. Build a simple consistent hash ring with 3 nodes and 100 virtual nodes. Distribute 10,000 keys and verify the distribution is roughly equal across nodes.
3. Implement `ReadRepair` for a simple in-memory replicated key-value store with 3 nodes. Introduce artificial version drift and verify that reads heal the stale replicas.

### Medium
4. Implement an **optimistic concurrency control** system: records have a version number. `Update(id, data, expectedVersion)` succeeds only if the current version matches. Return `ErrConflict` otherwise. This prevents lost updates without locking.
5. Build a **two-phase commit** protocol (simplified): a coordinator node asks all participant nodes to `Prepare(txn)`. If all say yes, it sends `Commit(txn)`. If any say no, it sends `Abort(txn)`. Simulate a participant failure and verify the coordinator handles it.
6. Implement a **read-your-writes** guarantee for a PostgreSQL read replica setup: after a write to the primary, store the `pg_current_wal_lsn()` in a cookie. All subsequent reads from this client include `SET LOCAL synchronous_commit = local; SET LOCAL min_recovery_apply_delay = 0;` and wait until the replica catches up to that LSN.

### Hard
7. Build a **distributed counter** using optimistic concurrency: N nodes each maintain a local count. Periodically, nodes gossip their counts to each other and merge using the G-Counter CRDT. Show that the global count is eventually consistent and never lost even if nodes fail.
8. Implement **Raft leader election** (simplified): 3 nodes, each with a term counter and a `Voted` flag. Implement `RequestVote` and `AppendEntries` RPCs. Show that a single leader is elected and that a failed leader causes a new election. This is how etcd, CockroachDB, and TiKV work internally.
