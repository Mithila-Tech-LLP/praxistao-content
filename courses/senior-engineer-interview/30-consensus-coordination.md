# Chapter 30: Consensus & Coordination — Raft, etcd & Distributed Locks

Distributed coordination is what makes leader election, service discovery, and distributed locking work. Senior engineers at infrastructure-heavy companies (Uber, Google, Stripe) are expected to understand consensus deeply.

## Table of Contents

1. [Why Consensus is Hard](#1-why-consensus-is-hard)
2. [The Raft Consensus Algorithm](#2-the-raft-consensus-algorithm)
3. [etcd — Distributed Key-Value Store](#3-etcd--distributed-key-value-store)
4. [Distributed Locks](#4-distributed-locks)
5. [Leader Election in Go](#5-leader-election-in-go)
6. [Service Discovery](#6-service-discovery)
7. [Interview Questions & Model Answers](#7-interview-questions--model-answers)
8. [Summary](#summary)

---

## 1. Why Consensus is Hard

In a distributed system, you need multiple nodes to agree on a single value (the "leader", the "current term", the "committed log index"). This is hard because:

```
Problems:
  - Nodes can crash at any time
  - Messages can be delayed, reordered, or lost (but not fabricated — Byzantine faults are different)
  - There's no global clock — you can't tell if a node is slow or dead
  
FLP Impossibility Theorem (1985):
  In a purely asynchronous system with even ONE possible failure, 
  no deterministic consensus algorithm can always terminate.
  
Practical solution: add timeouts (make the system partially synchronous)
  "If a node doesn't respond within 150ms, assume it's dead and elect a new leader"
```

---

## 2. The Raft Consensus Algorithm

Raft was designed to be "understandable" (unlike Paxos). It's used in etcd, CockroachDB, TiKV, and many production systems.

### Core Concepts

```
Node roles:
  Leader:    receives all writes, replicates to followers
  Follower:  receives log entries from leader, votes in elections
  Candidate: requesting votes to become leader

Terms:
  Raft divides time into terms (monotonically increasing integers)
  Each term starts with an election
  A term can have 0 or 1 leaders
```

### Leader Election

```
Normal operation:
  Leader sends heartbeats to all followers every 50ms
  Followers reset their "election timeout" (150-300ms) on each heartbeat

Leader fails:
  Followers stop receiving heartbeats
  First follower whose timer fires becomes a Candidate
  Candidate increments term, votes for itself, sends RequestVote to all
  
Voting:
  A node grants a vote if:
  1. It hasn't voted in this term yet
  2. The candidate's log is at least as up-to-date as its own
  
  Candidate wins if it gets majority (N/2 + 1) votes
  New leader immediately sends heartbeats to suppress other elections

Split vote:
  If two candidates get equal votes, wait for next timeout, start new election
  Randomized timeouts prevent this from repeating forever
```

### Log Replication

```
Client request → Leader

Leader:
  1. Appends entry to its own log (not committed yet)
  2. Sends AppendEntries RPC to all followers
  3. When majority confirm receipt: entry is committed
  4. Leader applies entry to state machine, responds to client
  5. Next heartbeat tells followers to commit the entry

Safety guarantee:
  An entry is committed only when a MAJORITY has it in their log
  Therefore: any new leader must have ALL committed entries (voting ensures this)
  
  If 3 nodes: 2 must confirm → a new leader must have been in that majority
  Even if the third node (minority) had stale log, the leader always has everything committed
```

### Why Raft is Safe

```
The "Log Matching" invariant:
  If two log entries have the same index and term, the logs are identical up to that point.
  This is enforced by leaders including the previous entry's (index, term) in AppendEntries.

The "Leader Completeness" property:
  If a log entry is committed in a term, it appears in the logs of ALL future leaders.
  Ensured by: votes are only given to candidates whose log is at least as up-to-date.
```

---

## 3. etcd — Distributed Key-Value Store

etcd is a distributed key-value store built on Raft. It's the backbone of Kubernetes (stores cluster state) and many distributed systems.

```go
import clientv3 "go.etcd.io/etcd/client/v3"

// Connect to etcd cluster
client, _ := clientv3.New(clientv3.Config{
    Endpoints:   []string{"localhost:2379", "localhost:2380", "localhost:2381"},
    DialTimeout: 5 * time.Second,
})
defer client.Close()

ctx := context.Background()

// Put (linearizable write — goes through Raft):
_, err := client.Put(ctx, "/config/feature_flag", "enabled")

// Get (linearizable read by default):
resp, err := client.Get(ctx, "/config/feature_flag")
for _, kv := range resp.Kvs {
    fmt.Printf("Key: %s, Value: %s\n", kv.Key, kv.Value)
}

// Watch — receive events when keys change:
watchChan := client.Watch(ctx, "/config/", clientv3.WithPrefix())
for event := range watchChan {
    for _, ev := range event.Events {
        fmt.Printf("Event: %s %s %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
    }
}

// Transactions (compare-and-swap):
_, err = client.Txn(ctx).
    If(clientv3.Compare(clientv3.Value("/leader"), "=", "node-1")).
    Then(clientv3.OpPut("/leader", "node-2")).
    Else(clientv3.OpGet("/leader")).
    Commit()
```

---

## 4. Distributed Locks

A distributed lock ensures that at most one node executes a critical section at a time across a distributed system.

### Lock with etcd (Using Leases)

```go
import (
    clientv3 "go.etcd.io/etcd/client/v3"
    "go.etcd.io/etcd/client/v3/concurrency"
)

func doWithDistributedLock(client *clientv3.Client, lockKey string, fn func() error) error {
    // Create a session with a 30-second TTL
    // If the node holding the lock crashes, the lease expires and the lock is released
    session, err := concurrency.NewSession(client, concurrency.WithTTL(30))
    if err != nil { return err }
    defer session.Close()
    
    // Try to acquire the lock
    mutex := concurrency.NewMutex(session, lockKey)
    if err := mutex.Lock(context.Background()); err != nil {
        return fmt.Errorf("failed to acquire lock: %w", err)
    }
    defer mutex.Unlock(context.Background())
    
    // Execute the critical section
    return fn()
}

// Usage:
err := doWithDistributedLock(client, "/locks/inventory-update", func() error {
    // Only one node executes this at a time
    return updateInventory()
})
```

### Lock with Redis (Redlock)

```go
// Simplified Redis distributed lock using SET NX PX
func acquireRedisLock(rdb *redis.Client, lockKey, lockValue string, ttl time.Duration) (bool, error) {
    result, err := rdb.SetNX(context.Background(), lockKey, lockValue, ttl).Result()
    return result, err
}

func releaseRedisLock(rdb *redis.Client, lockKey, lockValue string) error {
    // Lua script ensures atomic check-and-delete
    // Only delete if we own the lock (value matches)
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `)
    return script.Run(context.Background(), rdb, []string{lockKey}, lockValue).Err()
}

func withRedisLock(rdb *redis.Client, lockKey string, fn func() error) error {
    lockValue := uuid.New().String() // unique per lock holder
    
    acquired, err := acquireRedisLock(rdb, lockKey, lockValue, 30*time.Second)
    if err != nil { return err }
    if !acquired { return errors.New("lock not acquired") }
    
    defer releaseRedisLock(rdb, lockKey, lockValue)
    return fn()
}
```

### Distributed Lock Considerations

```
Lease/TTL expiry:
  The lock holder must renew the lease if the operation takes longer than the TTL
  If the lock holder is paused (GC pause, overload), the lock might expire mid-operation
  
Fencing tokens:
  Each lock acquisition gets an incrementing token (like a term in Raft)
  All operations must include the token; the resource rejects old tokens
  Prevents "zombie" processes from acting after their lock has expired

Split-brain:
  Redis single-node: if Redis goes down, no locking at all
  Redlock (multi-node): requires majority of Redis nodes — more complex
  etcd: Raft consensus ensures safety even with minority failures
  
Recommendation: use etcd for production distributed locks; Redis is simpler but less safe
```

---

## 5. Leader Election in Go

```go
// Leader election using etcd concurrency package
type LeaderElector struct {
    client   *clientv3.Client
    session  *concurrency.Session
    election *concurrency.Election
    nodeID   string
}

func NewLeaderElector(client *clientv3.Client, electionKey, nodeID string) (*LeaderElector, error) {
    session, err := concurrency.NewSession(client, concurrency.WithTTL(15))
    if err != nil { return nil, err }
    
    return &LeaderElector{
        client:   client,
        session:  session,
        election: concurrency.NewElection(session, electionKey),
        nodeID:   nodeID,
    }, nil
}

func (le *LeaderElector) Run(ctx context.Context, onLeader func(ctx context.Context)) error {
    defer le.session.Close()
    
    for {
        // Campaign blocks until this node becomes leader
        if err := le.election.Campaign(ctx, le.nodeID); err != nil {
            if ctx.Err() != nil { return nil } // context cancelled
            return err
        }
        
        fmt.Printf("Node %s is now the leader\n", le.nodeID)
        
        // Create a sub-context; cancel it to resign leadership
        leaderCtx, cancel := context.WithCancel(ctx)
        
        // Run the leader-specific work
        go onLeader(leaderCtx)
        
        // Watch for leadership loss (session expiry, resign)
        select {
        case <-le.session.Done():
            fmt.Println("Session expired, lost leadership")
            cancel()
        case <-ctx.Done():
            le.election.Resign(context.Background())
            cancel()
            return nil
        }
    }
}
```

---

## 6. Service Discovery

etcd and Consul are commonly used for service discovery — nodes register their addresses, clients look them up.

```go
// Service registration:
func registerService(client *clientv3.Client, serviceName, addr string) error {
    session, _ := concurrency.NewSession(client, concurrency.WithTTL(10))
    key := fmt.Sprintf("/services/%s/%s", serviceName, addr)
    _, err := client.Put(context.Background(), key, addr,
        clientv3.WithLease(session.Lease()))
    return err
    // Key is automatically deleted when session expires (node goes down)
}

// Service discovery:
func discoverServices(client *clientv3.Client, serviceName string) ([]string, error) {
    prefix := fmt.Sprintf("/services/%s/", serviceName)
    resp, err := client.Get(context.Background(), prefix, clientv3.WithPrefix())
    if err != nil { return nil, err }
    
    var addrs []string
    for _, kv := range resp.Kvs {
        addrs = append(addrs, string(kv.Value))
    }
    return addrs, nil
}
```

---

## 7. Interview Questions & Model Answers

**Q: Explain how Raft achieves consensus.**

"Raft divides time into terms. In each term, one leader is elected by getting votes from a majority of nodes. The election rule ensures a candidate can only win if its log is at least as up-to-date as any voting node's log, which guarantees the new leader has all previously committed entries. The leader then replicates log entries by sending AppendEntries RPCs; an entry is committed when a majority confirms receipt. This means committed entries are safe even if the leader fails — any new leader must have been in the majority that acknowledged the entry."

**Q: What is the difference between a distributed lock with Redis vs etcd?**

"A Redis-based distributed lock is simpler: SET key NX PX TTL tries to atomically set a key if it doesn't exist. The main risk is that Redis is single-node by default — if Redis goes down, the lock is gone. Redlock uses multiple Redis nodes for safety but is controversial (Martin Kleppmann wrote a famous critique about fencing token issues). etcd-based locks use Raft consensus: the lock operation itself goes through the Raft log, so it's safe as long as a majority of etcd nodes are available. etcd also provides fencing tokens (revision numbers) natively. For production systems with strong safety requirements, etcd is the better choice."

---

## Summary

- **Consensus** is needed for: leader election, distributed locks, consistent configuration.
- **Raft:** leader election via majority vote. Log replication to majority = committed. New leader always has all committed entries.
- **etcd:** Raft-based KV store. Used in Kubernetes. Provides: KV with linearizable reads, watch, transactions, leases.
- **Distributed locks:** use leases/TTLs so the lock releases if the holder crashes. Use fencing tokens to protect against zombies.
- etcd > Redis for distributed locks where correctness matters.
- **Service discovery:** nodes register with a lease; clients watch for changes. Lease expiry means automatic deregistration on failure.
