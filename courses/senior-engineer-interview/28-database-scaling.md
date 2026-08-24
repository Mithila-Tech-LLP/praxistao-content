# Chapter 28: Database Scaling — Replication, Sharding & Connection Pooling

Scaling databases is one of the most common system design topics. You need to know read replicas, sharding strategies, connection pooling, and partitioning deeply — not just the names.

## Table of Contents

1. [Vertical vs Horizontal Scaling](#1-vertical-vs-horizontal-scaling)
2. [Read Replicas](#2-read-replicas)
3. [Database Sharding](#3-database-sharding)
4. [Partitioning](#4-partitioning)
5. [Connection Pooling](#5-connection-pooling)
6. [Caching Architecture](#6-caching-architecture)
7. [CQRS Pattern](#7-cqrs-pattern)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Vertical vs Horizontal Scaling

**Vertical scaling (scale up):** bigger machine. More CPU, RAM, faster disk. Simple — no code changes. But there's a ceiling: the biggest AWS RDS instance is still one machine.

**Horizontal scaling (scale out):** more machines. Distribute load. Requires architectural changes — connections, consistency, data distribution.

For databases, the typical scaling journey:
1. Start: single PostgreSQL instance
2. Add read replicas (read scale)
3. Add caching (read scale + latency)
4. Partition tables (write scale within one server)
5. Shard across multiple databases (write scale across servers)
6. Each step adds operational complexity. Don't jump steps.

---

## 2. Read Replicas

Most applications are read-heavy (80-90% reads). Read replicas are copies of the primary database that serve read traffic.

```
                    ┌─────────────┐
    Writes ────────▶│   Primary   │─── WAL stream ───▶ Replica 1
                    │  (leader)   │                    Replica 2
                    └─────────────┘                    Replica 3
    Reads ──────────────────────────────▶ (load balanced across replicas)
```

### Replication Lag

Replicas are typically asynchronous — the primary doesn't wait for replicas to confirm before responding to the client. This means replicas are slightly behind.

```
Problem: read-your-own-write inconsistency
  User updates profile picture → write goes to primary
  User immediately loads profile page → reads from replica (old picture!)
  
Solutions:
  1. "Read-after-write" routing: if user just wrote, route next read to primary
  2. Sticky sessions: always route a user's reads to the same replica
  3. Wait-for-replica-catch-up: add short delay before redirecting read
  4. Accept eventual consistency: most social media reads are "good enough"

-- Check replication lag in PostgreSQL:
SELECT 
    client_addr,
    write_lag,
    flush_lag,
    replay_lag
FROM pg_stat_replication;
```

### Go Code: Primary/Replica Routing

```go
type DB struct {
    primary  *sql.DB
    replicas []*sql.DB
    counter  atomic.Int64
}

func (db *DB) Write() *sql.DB {
    return db.primary
}

func (db *DB) Read() *sql.DB {
    // Round-robin across replicas
    idx := db.counter.Add(1) % int64(len(db.replicas))
    return db.replicas[idx]
}

// Usage:
func (r *UserRepo) UpdateProfile(ctx context.Context, u *User) error {
    _, err := r.db.Write().ExecContext(ctx, "UPDATE users SET ...")
    return err
}

func (r *UserRepo) GetUser(ctx context.Context, id int) (*User, error) {
    row := r.db.Read().QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
    // ...
}
```

---

## 3. Database Sharding

Sharding splits data across multiple database instances. Each shard holds a subset of rows. Enables horizontal write scaling when a single primary can't keep up.

### Sharding Strategies

**Range-based sharding:**
```
Shard 1: user_id 1 – 1,000,000
Shard 2: user_id 1,000,001 – 2,000,000
Shard 3: user_id 2,000,001 – 3,000,000

Pros: range queries efficient (all 2023 orders on one shard)
Cons: hotspot problem (new users all go to shard 3)
```

**Hash-based sharding:**
```
shard_index = hash(user_id) % num_shards

Pros: even distribution, no hotspots
Cons: range queries require scanning all shards
     Resharding when adding new shards requires rebalancing data
```

**Directory-based sharding:**
```
A lookup service maps entity_id → shard_id
Pros: flexible, can handle uneven data sizes, easy resharding
Cons: lookup service becomes a bottleneck/SPOF
```

### The Problems With Sharding

```
1. Cross-shard queries: "find all orders > $100" must go to all shards
2. Cross-shard transactions: no native ACID across shards
3. Resharding: adding a shard requires migrating data
4. Schema changes: must be applied to all shards
5. Joins: impossible across shards (data must be co-located)

Solution for cross-shard joins: denormalize (duplicate data) or use a different storage layer
```

### Go Code: Shard Router

```go
type ShardRouter struct {
    shards []*sql.DB
}

func (r *ShardRouter) ShardFor(userID int64) *sql.DB {
    idx := userID % int64(len(r.shards))
    return r.shards[idx]
}

// Always pass entity ID to route correctly:
func (r *UserRepo) GetOrder(ctx context.Context, userID, orderID int64) (*Order, error) {
    shard := r.router.ShardFor(userID)
    var o Order
    shard.QueryRowContext(ctx, "SELECT * FROM orders WHERE id = $1 AND user_id = $2",
        orderID, userID).Scan(&o.ID, &o.Total)
    return &o, nil
}
```

---

## 4. Partitioning

Partitioning splits a large table into smaller physical pieces within the same database. Unlike sharding, partitioning is transparent to queries — it's a storage optimization.

```sql
-- Range partitioning by date:
CREATE TABLE orders (
    id BIGINT,
    user_id BIGINT,
    created_at TIMESTAMP,
    total DECIMAL
) PARTITION BY RANGE (created_at);

CREATE TABLE orders_2023 PARTITION OF orders
    FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE orders_2024 PARTITION OF orders
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- Queries with WHERE created_at = '2024-06-01' only scan the 2024 partition
-- PostgreSQL prunes irrelevant partitions automatically

-- Hash partitioning (even distribution):
CREATE TABLE users (id BIGINT, name TEXT) PARTITION BY HASH (id);
CREATE TABLE users_p0 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE users_p1 PARTITION OF users FOR VALUES WITH (MODULUS 4, REMAINDER 1);
-- etc.

-- Benefits:
-- Partition pruning: queries touch only relevant partitions
-- Easier archival: DROP TABLE orders_2020 (instant vs DELETE with millions of rows)
-- Parallel queries: multiple workers can scan different partitions simultaneously
```

---

## 5. Connection Pooling

Already covered in Chapter 26 (PgBouncer). Key additions for scaling:

```
Per-service pools:
  Service A → PgBouncer pool (20 connections)
  Service B → PgBouncer pool (20 connections)
  
  Both connect to the same PostgreSQL instance.
  PostgreSQL sees 40 connections max.
  Service A can use 100 goroutines that time-share the 20 connections.

Connection lifecycle:
  app opens connection → PgBouncer assigns a real PG connection
  app runs query
  app closes connection (transaction mode) → PgBouncer recycles PG connection
  PG connection is now available for the next request

Setting max_open_conns in Go:
  db.SetMaxOpenConns(25)       // max connections to PgBouncer (or direct PG)
  db.SetMaxIdleConns(10)       // keep this many warm
  db.SetConnMaxLifetime(5 * time.Minute)  // recycle to prevent stale connections
```

---

## 6. Caching Architecture

Adding a cache in front of the database is often more impactful than database scaling:

```
Read request:
  1. Check L1 cache (in-process, e.g. Go map with TTL) — microseconds
  2. Check L2 cache (Redis) — sub-millisecond
  3. Read from primary/replica database — milliseconds

Cache hit ratios:
  Even 95% cache hit ratio means 20x fewer database reads
  99% hit ratio means 100x fewer database reads
  
Cache invalidation strategies:
  TTL-based: simple but stale data during TTL window
  Write-through: update cache when writing to DB (consistent but slower writes)
  Event-based: publish change events, subscribers invalidate their cache keys
```

---

## 7. CQRS Pattern

Command Query Responsibility Segregation: separate the write model from the read model.

```
Write side:                           Read side:
  Commands → Command Handler            Queries → Read Model (optimized projections)
  Updates primary DB                    Reads from denormalized read DB / cache

Example:
  Write: normalize order data into orders, order_items, products tables
  Read: maintain a pre-computed "order_summary" table updated by events

Why:
  - Write side: normalized schema, ACID transactions, correctness
  - Read side: denormalized for fast reads, optimized for specific query shapes
  - Each side can scale independently
```

---

## 8. Interview Questions & Model Answers

**Q: How would you scale a database that's hitting write capacity limits?**

"First I'd confirm the bottleneck is writes, not reads. If it's writes, I'd look at: (1) batch writes instead of per-row inserts, (2) asynchronous write paths with a queue (write to Kafka, flush to DB in batches), (3) table partitioning to spread I/O across files, (4) offloading to a write-optimized database like Cassandra for specific high-write data types. If all else fails, horizontal sharding — but that significantly increases operational complexity and should be a last resort."

**Q: What is the difference between sharding and partitioning?**

"Partitioning is within a single database instance — the table is split into multiple physical storage segments, but it's transparent to queries and managed by the database engine. PostgreSQL's declarative partitioning is a good example. Sharding is across multiple database instances — data is split and routed at the application layer. Sharding gives you more raw capacity but requires application-level routing logic, makes cross-shard queries expensive, and complicates transactions. Always partition before sharding."

---

## Summary

- Scale path: single DB → read replicas → caching → partitioning → sharding. Each step adds complexity.
- **Read replicas:** async replication. Watch for lag. Route analytics/reporting to replicas.
- **Sharding:** data partitioned across multiple DB instances. Hash-based avoids hotspots; range-based enables efficient range queries.
- **Partitioning:** within one DB. PostgreSQL range/hash partitioning enables partition pruning and easy archival.
- **Connection pooling:** PgBouncer multiplexes thousands of app connections onto few real DB connections.
- **Caching:** Redis in front of DB reduces read load dramatically. TTL + cache-aside is the standard pattern.
- **CQRS:** separate write and read models for independent scaling and optimization.
