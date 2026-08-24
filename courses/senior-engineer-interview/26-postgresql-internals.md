# Chapter 26: PostgreSQL Internals — WAL, VACUUM & Query Planning

Understanding PostgreSQL internals helps you design more reliable systems, diagnose production issues faster, and answer the deep database questions that Google/Stripe/Uber interviewers ask senior candidates.

## Table of Contents

1. [Write-Ahead Log (WAL)](#1-write-ahead-log-wal)
2. [Checkpoints & Crash Recovery](#2-checkpoints--crash-recovery)
3. [VACUUM & Dead Tuple Management](#3-vacuum--dead-tuple-management)
4. [Query Planning — How PostgreSQL Chooses Plans](#4-query-planning--how-postgresql-chooses-plans)
5. [EXPLAIN ANALYZE Deep Dive](#5-explain-analyze-deep-dive)
6. [Connection Pooling with PgBouncer](#6-connection-pooling-with-pgbouncer)
7. [PostgreSQL Configuration Knobs](#7-postgresql-configuration-knobs)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Write-Ahead Log (WAL)

WAL is PostgreSQL's durability mechanism. Before modifying data pages in memory, PostgreSQL first writes a record of the change to the WAL (on disk). This ensures that even if the server crashes before flushing dirty data pages, the changes can be replayed from the WAL on restart.

```
Without WAL:
  1. Update page in shared_buffers (in memory)
  2. Crash! Page not flushed yet.
  3. Restart: page is lost. Data loss.

With WAL:
  1. Write WAL record to WAL file (sequential write, fast)
  2. Update page in shared_buffers (in memory)
  3. Commit: WAL is fsync'd to disk → transaction confirmed
  4. Crash! 
  5. Restart: replay WAL records → all committed transactions restored
```

### WAL and Replication

WAL is also the basis for streaming replication:
- Primary writes WAL records to WAL files
- Standby reads WAL records via `pg_receivewal` or streaming
- Standby replays WAL records to stay in sync with primary

```sql
-- Check WAL status:
SELECT * FROM pg_stat_replication;  -- active standby connections
SELECT pg_current_wal_lsn();        -- current WAL position on primary
SELECT pg_last_wal_replay_lsn();    -- last replayed position on standby

-- WAL configuration in postgresql.conf:
-- wal_level = replica        (or logical for logical replication)
-- max_wal_senders = 10       (max standbys)
-- wal_keep_size = 1GB        (retain WAL for standbys that are behind)
```

---

## 2. Checkpoints & Crash Recovery

PostgreSQL periodically writes dirty pages from shared_buffers to data files. This is called a **checkpoint**.

```
Timeline:
  ←———checkpoint———→ ←———WAL records———→ ←crash→ ←replay from checkpoint→

After crash:
  1. Start from last completed checkpoint (data files are consistent up to here)
  2. Replay WAL records from that checkpoint forward
  3. All committed transactions are recovered
  4. Aborted transactions' changes are never written (they don't commit)
```

```sql
-- Force a checkpoint:
CHECKPOINT;

-- Checkpoint configuration:
-- checkpoint_timeout = 5min     (max time between checkpoints)
-- max_wal_size = 1GB            (WAL size triggers checkpoint)
-- checkpoint_completion_target = 0.9  (spread I/O over 90% of checkpoint interval)
```

---

## 3. VACUUM & Dead Tuple Management

Every UPDATE in PostgreSQL creates a new row version (MVCC) and marks the old row as dead. Dead tuples waste disk space and slow down queries. VACUUM reclaims this space.

```sql
-- Dead tuple accumulation:
-- UPDATE users SET name = 'Bob' WHERE id = 1
-- Before: [xmin=5, xmax=0]  {id=1, name='Alice'}
-- After:  [xmin=5, xmax=100] {id=1, name='Alice'}  ← dead tuple
--         [xmin=100, xmax=0]  {id=1, name='Bob'}    ← live tuple

-- VACUUM removes dead tuples (but doesn't shrink the file):
VACUUM users;

-- VACUUM FULL compacts the table (requires exclusive lock — avoid on large tables!):
VACUUM FULL users;  -- dangerous! locks the table

-- ANALYZE updates statistics for the query planner:
ANALYZE users;

-- VACUUM ANALYZE: both together (safest routine maintenance):
VACUUM ANALYZE users;

-- Monitor dead tuples:
SELECT 
    relname,
    n_live_tup,
    n_dead_tup,
    ROUND(n_dead_tup::numeric / NULLIF(n_live_tup + n_dead_tup, 0) * 100, 2) AS dead_pct,
    last_vacuum,
    last_autovacuum
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;
```

### Table Bloat

If VACUUM can't keep up (e.g., long-running transactions hold back the oldest snapshot), tables bloat with dead tuples:

```sql
-- Find bloated tables:
SELECT 
    schemaname, tablename,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname || '.' || tablename)) AS table_size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname || '.' || tablename) DESC;

-- Check for long-running transactions that prevent vacuum:
SELECT pid, age(clock_timestamp(), query_start) AS elapsed, query
FROM pg_stat_activity
WHERE state != 'idle' AND query_start < NOW() - INTERVAL '5 minutes'
ORDER BY elapsed DESC;
```

---

## 4. Query Planning — How PostgreSQL Chooses Plans

PostgreSQL's query planner picks the cheapest execution strategy based on cost estimates derived from table statistics.

### Statistics & the Planner

```sql
-- PostgreSQL collects statistics automatically (via autovacuum):
-- pg_stats: per-column statistics
SELECT tablename, attname, n_distinct, correlation, most_common_vals
FROM pg_stats
WHERE tablename = 'users';

-- n_distinct: estimated unique values (negative = fraction of total rows)
-- correlation: correlation between column order and physical row order (1.0 = sequential access)
-- most_common_vals: most common values (used for cardinality estimates)

-- Statistics target (how much data to sample):
ALTER TABLE orders ALTER COLUMN status SET STATISTICS 500;
-- Default is 100. Increase for columns with complex distributions.
```

### Join Strategies

PostgreSQL picks one of three join algorithms:
- **Hash Join:** build hash table from smaller relation, probe with larger. Good for large tables without indexes.
- **Nested Loop Join:** for each row in outer table, scan inner table. Fast when inner side has an index and few rows match.
- **Merge Join:** sort both sides on join key, merge. Good when both sides are pre-sorted or have indexes.

```sql
-- See which join is being used:
EXPLAIN SELECT e.name, d.name
FROM employees e JOIN departments d ON e.department_id = d.id;

-- Example output:
-- Hash Join  (cost=3.27..7.64 rows=10 width=64)
--   Hash Cond: (e.department_id = d.id)
--   ->  Seq Scan on employees  (cost=0.00..1.10 rows=10 width=36)
--   ->  Hash  (cost=2.15..2.15 rows=15 width=40)
--         ->  Seq Scan on departments  (cost=0.00..2.15 rows=15 width=40)
```

---

## 5. EXPLAIN ANALYZE Deep Dive

`EXPLAIN ANALYZE` runs the actual query and shows what happened (not just what was estimated):

```sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT o.id, o.total, u.name
FROM orders o
JOIN users u ON o.user_id = u.id
WHERE o.status = 'pending'
ORDER BY o.created_at DESC
LIMIT 20;

-- Example output with annotations:
-- Sort  (cost=45.23..45.25 rows=8) (actual time=0.52..0.53 rows=8 loops=1)
--   Sort Key: o.created_at DESC
--   Sort Method: quicksort  Memory: 25kB          ← fits in work_mem, no disk sort
--   Buffers: shared hit=12                         ← all from cache, no disk read
--   ->  Hash Join  (cost=2.18..45.12 rows=8) (actual time=0.18..0.49 rows=8 loops=1)
--         Hash Cond: (o.user_id = u.id)
--         Buffers: shared hit=12
--         ->  Seq Scan on orders  (cost=0.00..42.50 rows=8) (actual time=0.05..0.35 rows=8 loops=1)
--               Filter: (status = 'pending')
--               Rows Removed by Filter: 1490      ← 1490 rows scanned but removed — add index!
--               Buffers: shared hit=10
--         ->  Hash  (cost=1.15..1.15 rows=15) (actual time=0.10..0.10 rows=15 loops=1)
--               Buckets: 1024  Batches: 1  Memory Usage: 9kB
--               ->  Seq Scan on users  (cost=0.00..1.15 rows=15 width=40) (actual time=0.04..0.07 rows=15 loops=1)
-- Planning Time: 0.122 ms
-- Execution Time: 0.581 ms

-- Key things to look for:
-- "Rows Removed by Filter: 1490" on a Seq Scan → you need an index on that column
-- "actual rows=1490 estimate rows=1" → planner underestimated, stats are stale (run ANALYZE)
-- "Batches: 4" on Hash → hash exceeded work_mem, spilled to disk
-- "loops=N" → this node ran N times (Nested Loop!), multiply the actual times
```

---

## 6. Connection Pooling with PgBouncer

PostgreSQL creates a new OS process per connection. This is expensive. At 10,000 concurrent connections, PostgreSQL struggles — each process uses ~10MB RAM and context-switching overhead is huge.

**PgBouncer** sits between your app and PostgreSQL, maintaining a small pool of real database connections and multiplexing many app connections onto them.

```
App servers (100 processes × 10 connections = 1000 app connections)
      ↓
PgBouncer (pool of 20 real database connections)
      ↓
PostgreSQL (20 connections — manageable)
```

```ini
# pgbouncer.ini
[databases]
mydb = host=127.0.0.1 port=5432 dbname=mydb

[pgbouncer]
pool_mode = transaction      # release connection after each transaction (best for most apps)
# pool_mode = session        # hold connection for entire session (needed for session-level features)
# pool_mode = statement      # release after each statement (very limited use)

max_client_conn = 1000       # max app connections
default_pool_size = 20       # real PG connections per user/db pair
```

---

## 7. PostgreSQL Configuration Knobs

```ini
# postgresql.conf — key settings for production:

# Memory
shared_buffers = 25% of RAM         # PostgreSQL's main page cache
effective_cache_size = 75% of RAM    # hint to planner about OS cache
work_mem = 64MB                      # per-sort/hash-join memory (×parallel workers!)
maintenance_work_mem = 512MB         # for VACUUM, CREATE INDEX

# Write performance
wal_buffers = 64MB                   # WAL buffer (16MB default is often too small)
checkpoint_completion_target = 0.9   # spread checkpoint I/O
random_page_cost = 1.1              # for SSD; default 4.0 is for spinning disk!

# Connections
max_connections = 100                # keep low; use PgBouncer for more
```

---

## 8. Interview Questions & Model Answers

**Q: What is WAL and why is it critical for durability?**

"WAL (Write-Ahead Log) ensures durability by writing a record of every change to a sequential log on disk before modifying data pages in memory. A transaction is considered committed only after its WAL record is flushed to disk. If the server crashes, PostgreSQL replays WAL records from the last checkpoint on startup — no committed data is lost. WAL also enables streaming replication: standbys read the primary's WAL stream and replay it to maintain a consistent copy."

**Q: Why does a long-running transaction cause table bloat in PostgreSQL?**

"PostgreSQL's MVCC keeps old row versions visible to any transaction that might still need them. The oldest active transaction determines the 'horizon' — VACUUM cannot remove dead tuples that are newer than the oldest transaction's snapshot. If a transaction runs for hours, VACUUM is blocked from reclaiming dead tuples for the entire duration. The dead tuples accumulate in the table file (bloat). The fix is to avoid long-running transactions, monitor `pg_stat_activity` for old connections, and set `idle_in_transaction_session_timeout`."

---

## Summary

- **WAL:** sequential write of every change before modifying data pages. Enables crash recovery and replication.
- **Checkpoint:** periodic flush of dirty pages to disk. Recovery starts from last checkpoint and replays WAL.
- **VACUUM:** removes dead tuples created by MVCC. Autovacuum handles it automatically; long transactions can block it.
- **VACUUM FULL:** compacts the file but requires exclusive lock — avoid on large production tables.
- **Query planner:** uses pg_stats to estimate row counts and pick join strategies (Hash/Nested Loop/Merge).
- **EXPLAIN ANALYZE:** see actual vs estimated rows. Watch for high "Rows Removed by Filter" and planner mis-estimates.
- **PgBouncer:** connection pooler. Transaction mode is the default — releases PG connection after each transaction.
- **random_page_cost = 1.1** for SSDs — one of the highest-impact config changes.
