# Chapter 17: PostgreSQL Performance Tuning and Scaling

Your schema is designed, your indexes are in place. Now you need to make PostgreSQL go fast — and keep it fast as your data grows. This chapter is a practical performance guide used by engineers at companies running PostgreSQL at massive scale.

## Table of Contents

1. Reading EXPLAIN ANALYZE Like a Pro
2. Configuration Tuning — The Most Impactful Settings
3. Index Strategies for Production
4. Replication — Read Replicas and High Availability
5. Connection Pooling at Scale
6. Benchmarking with pgbench
7. Exercises

---

## 1. Reading EXPLAIN ANALYZE Like a Pro

`EXPLAIN ANALYZE` is your most powerful tool. Here's how to read it:

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON o.user_id = u.id
WHERE u.created_at > '2024-01-01'
GROUP BY u.id, u.name
ORDER BY order_count DESC
LIMIT 10;
```

Example output:
```
Limit  (cost=1234.56..1234.58 rows=10 width=32) (actual time=12.345..12.348 rows=10 loops=1)
  ->  Sort  (cost=1234.56..1256.78 rows=8888 width=32) (actual time=12.343..12.344 rows=10 loops=1)
        Sort Key: (count(o.id)) DESC
        Sort Method: top-N heapsort  Memory: 25kB
        ->  HashAggregate  (cost=987.65..1098.76 rows=8888 width=32) (actual time=10.123..11.234 rows=7500 loops=1)
              Buckets: 16384  Batches: 1  Memory Usage: 1281kB
              ->  Hash Left Join  (cost=234.56..876.54 rows=22222 width=16) (actual time=2.345..8.765 rows=22000 loops=1)
                    Hash Cond: (o.user_id = u.id)
                    ->  Seq Scan on orders o  (cost=0.00..456.78 rows=22222 width=8) (actual time=0.012..3.456 rows=22000 loops=1)
                    ->  Hash  (cost=198.76..198.76 rows=2864 width=16) (actual time=2.123..2.123 rows=2864 loops=1)
                          Buckets: 4096  Batches: 1  Memory Usage: 160kB
                          ->  Index Scan using idx_users_created_at on users u  (cost=0.43..198.76 rows=2864 width=16) (actual time=0.034..1.876 rows=2864 loops=1)
                                Index Cond: (created_at > '2024-01-01 00:00:00+00'::timestamptz)
Buffers: shared hit=345 read=123
Planning Time: 0.456 ms
Execution Time: 12.567 ms
```

What to look for:

**`actual time=X..Y rows=Z`**: X is startup time, Y is total time. Z is rows produced. If `rows` estimate is very different from actual, statistics are stale — run `ANALYZE`.

**`Seq Scan`**: Reading every row. Fine for small tables, bad for large. Ask: is there a missing index?

**`Index Scan`**: Using a B-Tree index. Usually good.

**`Index Only Scan`**: Serving the query entirely from the index. Best case.

**`Hash Join`**: Building a hash table from the smaller side, probing with the larger. Good for large datasets.

**`Nested Loop`**: For each row in outer, look up one row in inner (usually using an index). Best when outer is small.

**`Buffers: shared hit=X read=Y`**: `hit` = served from RAM (buffer pool). `read` = fetched from disk. Aim for high hit ratio.

### Finding Slow Queries Automatically

```sql
-- Requires pg_stat_statements extension
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Top 10 slowest queries by total execution time
SELECT
    round(mean_exec_time::numeric, 2) AS avg_ms,
    calls,
    round(total_exec_time::numeric, 2) AS total_ms,
    left(query, 100) AS query
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;
```

---

## 2. Configuration Tuning — The Most Impactful Settings

PostgreSQL's defaults are very conservative. For a dedicated database server, tune these settings in `postgresql.conf`:

### Memory

```ini
# shared_buffers: PostgreSQL's buffer pool. Set to 25% of total RAM.
shared_buffers = 4GB     # for a 16GB server

# effective_cache_size: Hint for the planner. Set to 75% of total RAM.
effective_cache_size = 12GB

# work_mem: Memory per sort/hash operation. Can be used by many operations at once.
# Too high = OOM. Too low = disk sorts (slow). 
# Formula: total_ram / (max_connections * max_parallel_workers_per_gather)
work_mem = 64MB

# maintenance_work_mem: Memory for VACUUM, CREATE INDEX, etc.
maintenance_work_mem = 1GB
```

### Write Performance

```ini
# checkpoint_completion_target: Spread checkpoint writes over this fraction of the checkpoint interval.
checkpoint_completion_target = 0.9

# wal_buffers: Size of WAL buffer. Set to shared_buffers / 32, max 16MB.
wal_buffers = 16MB

# max_wal_size: Trigger checkpoint when WAL exceeds this. Increase to reduce checkpoint frequency.
max_wal_size = 4GB

# synchronous_commit: Set to 'off' for non-critical workloads for 2-3x write speedup.
# Risk: last ~1ms of transactions might be lost on crash (not corruption).
synchronous_commit = on   # keep 'on' for financial data
```

### Connections

```ini
# Never exceed this in application code — use pgxpool or PgBouncer
max_connections = 100
```

Apply changes:
```bash
# After editing postgresql.conf:
docker exec pg pg_ctl reload  # for most settings (no restart needed)
# Some settings require restart (shared_buffers, max_connections)
```

---

## 3. Index Strategies for Production

### The Most Missed Index: Foreign Keys

```sql
-- Always index foreign key columns
CREATE INDEX idx_orders_user_id ON orders(user_id);
-- Without this, every JOIN from users to orders is a full scan of orders!
```

### Covering Indexes

If your query only needs certain columns, include them in the index:

```sql
-- Query: SELECT email, name FROM users WHERE status = 'active'
-- Normal index:
CREATE INDEX idx_users_status ON users(status);
-- → Index Scan + heap fetch (two I/Os)

-- Covering index:
CREATE INDEX idx_users_status_covering ON users(status) INCLUDE (email, name);
-- → Index Only Scan (one I/O) — no heap fetch needed!
```

### Conditional (Partial) Indexes

```sql
-- Only 10% of orders are 'pending'. Index only those:
CREATE INDEX idx_orders_pending ON orders(created_at) WHERE status = 'pending';

-- Dramatically smaller index, faster for the common case
SELECT * FROM orders WHERE status = 'pending' ORDER BY created_at;
```

### Index Maintenance

```sql
-- Find unused indexes (expensive to maintain, never used)
SELECT indexrelname, idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;

-- Rebuild bloated indexes
REINDEX TABLE CONCURRENTLY users; -- rebuilds without locking!
```

---

## 4. Replication — Read Replicas and High Availability

### Streaming Replication

PostgreSQL replication works by shipping WAL records from the primary to one or more replicas. Replicas apply the WAL and stay in sync.

```
Primary ──WAL stream──► Replica 1 (hot standby)
         ──WAL stream──► Replica 2 (hot standby)
```

With Docker Compose for local testing:

```yaml
# The replica can serve SELECT queries, offloading the primary
services:
  postgres-primary:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: secret
      POSTGRES_USER: replicator
    command: >
      postgres
      -c wal_level=replica
      -c max_wal_senders=3
      -c hot_standby=on

  postgres-replica:
    image: postgres:16
    depends_on:
      - postgres-primary
    # Replica connects to primary and streams WAL
```

### Read/Write Splitting in Go

```go
package main

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    primaryDB *pgxpool.Pool // all writes go here
    replicaDB *pgxpool.Pool // reads can go here
)

func init() {
    primaryDB, _ = pgxpool.New(context.Background(), "postgres://primary:5432/myapp")
    replicaDB, _ = pgxpool.New(context.Background(), "postgres://replica:5432/myapp")
}

func GetUser(ctx context.Context, id int64) (*User, error) {
    // Reads go to replica
    var u User
    err := replicaDB.QueryRow(ctx,
        "SELECT id, email, name FROM users WHERE id = $1", id,
    ).Scan(&u.ID, &u.Email, &u.Name)
    return &u, err
}

func CreateUser(ctx context.Context, email, name string) (int64, error) {
    // Writes MUST go to primary
    var id int64
    err := primaryDB.QueryRow(ctx,
        "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id",
        email, name,
    ).Scan(&id)
    return id, err
}
```

**Warning:** Replication lag means your replica may be slightly behind. Never read from a replica data you just wrote — you might not see it yet. For read-your-writes consistency, route that specific user's reads to the primary for a short window after a write.

---

## 5. Connection Pooling at Scale

For > 50 connections, PgBouncer (a standalone pooler) is more efficient than pgxpool alone:

```
Your App (100 goroutines)
    ↓ 100 connections to PgBouncer
PgBouncer (transaction pooling)
    ↓ 20 connections to PostgreSQL
PostgreSQL (max_connections = 100)
```

PgBouncer in transaction pooling mode reuses connections — a connection is returned to the pool after each transaction. 1000 concurrent users can share 20 database connections.

```ini
# pgbouncer.ini
[databases]
myapp = host=localhost port=5432 dbname=myapp

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
server_idle_timeout = 600
```

---

## 6. Benchmarking with pgbench

`pgbench` is a built-in PostgreSQL benchmarking tool:

```bash
# Initialize with a standard test schema
pgbench -h localhost -U dev -i -s 50 myapp
# -s 50 = scale factor (50 * 100,000 = 5 million rows)

# Run a standard TPC-B benchmark
pgbench -h localhost -U dev -T 60 -c 20 -j 4 myapp
# -T 60 = run for 60 seconds
# -c 20 = 20 concurrent clients
# -j 4  = 4 worker threads
```

Output:
```
latency average = 3.456 ms
latency stddev = 1.234 ms
tps = 5789.123456 (without initial connection time)
```

Custom benchmark scripts (test your specific queries):

```sql
-- custom_benchmark.sql
\set user_id random(1, 100000)
SELECT id, email FROM users WHERE id = :user_id;
```

```bash
pgbench -h localhost -U dev -T 30 -c 50 -f custom_benchmark.sql myapp
```

---

## Summary

- `EXPLAIN (ANALYZE, BUFFERS)` is your primary tool. Look for Seq Scans on large tables, estimate vs actual row discrepancies, and high `read` vs `hit` buffer counts.
- Key config settings: `shared_buffers` (25% RAM), `work_mem` (careful!), `max_wal_size` (reduce checkpoint frequency).
- Always index foreign keys. Use covering indexes and partial indexes for the most frequently-run queries.
- Streaming replication creates hot-standby replicas that can serve reads. Route writes to primary, reads to replica in Go.
- PgBouncer multiplexes thousands of app connections into a small number of PostgreSQL connections.

### Exercises

**Easy:** Run `EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@test.com'` before and after creating an index on email. Compare execution times.

**Medium:** Use `pg_stat_statements` to find the top 5 slowest queries in your application. Add indexes where missing.

**Hard:** Set up two PostgreSQL instances with streaming replication using Docker Compose. Write a Go program that writes to the primary and reads from the replica with a lag detection check.
