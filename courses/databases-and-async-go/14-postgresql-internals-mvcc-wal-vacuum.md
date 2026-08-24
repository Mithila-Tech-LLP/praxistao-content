# Chapter 14: PostgreSQL Internals — MVCC, WAL, and VACUUM

Most developers use PostgreSQL as a black box. This chapter opens the box. Understanding how PostgreSQL works internally makes you a dramatically better database user — you'll know why certain operations are slow, why VACUUM exists, and how your data survives crashes.

## Table of Contents

1. How PostgreSQL Stores Data — Pages and Tuples
2. MVCC — How PostgreSQL Handles Concurrent Access
3. The Write-Ahead Log — Crash Safety
4. VACUUM — Cleaning Up Dead Tuples
5. Checkpoints — Writing Dirty Pages to Disk
6. Monitoring Internals in Real Time
7. Connection Pooling with pgxpool
8. Exercises

---

## 1. How PostgreSQL Stores Data — Pages and Tuples

PostgreSQL stores table data in **heap files**. Each heap file is divided into **pages** of exactly 8 KB. Each page holds multiple **tuples** (rows).

```
Heap File for "users" table
┌──────────────────────────────────┐
│ Page 0 (8 KB)                    │
│  ┌────────────────────────────┐  │
│  │ Page Header (24 bytes)     │  │
│  │  - LSN (last WAL position) │  │
│  │  - checksum                │  │
│  │  - free space info         │  │
│  ├────────────────────────────┤  │
│  │ Item Pointers              │  │ ← array of (offset, length) for each tuple
│  ├────────────────────────────┤  │
│  │   (free space)             │  │
│  ├────────────────────────────┤  │
│  │ Tuple 3 (row)              │  │
│  │ Tuple 2 (row)              │  │
│  │ Tuple 1 (row)              │  │ ← tuples grow from end toward start
│  └────────────────────────────┘  │
└──────────────────────────────────┘
│ Page 1 (8 KB)                    │
│  ...                             │
```

Each tuple (row) has a **tuple header** containing:
- `xmin`: the transaction ID that created this tuple
- `xmax`: the transaction ID that deleted this tuple (or 0 if alive)
- `ctid`: physical location (page number, item index) — used by indexes

These fields are the foundation of MVCC.

---

## 2. MVCC — How PostgreSQL Handles Concurrent Access

MVCC (Multi-Version Concurrency Control) means PostgreSQL stores **multiple versions** of each row. When a row is updated, PostgreSQL doesn't overwrite the old data — it inserts a new version.

### How an UPDATE Works

```sql
-- Initial state: user 1, name = "Alice"
-- Row in page: xmin=100, xmax=0, name="Alice"

-- Transaction 200 runs: UPDATE users SET name = 'Alicia' WHERE id = 1
-- Result:
-- Old row: xmin=100, xmax=200, name="Alice"   ← marked as deleted by txn 200
-- New row: xmin=200, xmax=0,   name="Alicia"  ← created by txn 200
```

### Visibility Rules

When a transaction reads a row, it applies visibility rules:
- A row is visible if `xmin` committed **before** the current transaction started
- A row is invisible if `xmax` committed **before** the current transaction started (it was deleted)

This means **readers never block writers and writers never block readers**. Each transaction sees a consistent snapshot of the database.

```
Time ──────────────────────────────────────────────►
Txn 100: INSERT name="Alice" ... COMMIT
Txn 200: UPDATE name="Alicia"... (in progress)
Txn 300: SELECT * FROM users
         → sees name="Alice" (txn 200 not yet committed)
Txn 200: COMMIT
Txn 400: SELECT * FROM users
         → sees name="Alicia" (txn 200 now committed)
```

### The Dead Tuple Problem

Because UPDATE creates new rows instead of overwriting, the old rows ("dead tuples") accumulate. They're not visible to any transaction, but they still take up disk space. This is why VACUUM exists.

---

## 3. The Write-Ahead Log — Crash Safety

The WAL (Write-Ahead Log) is the mechanism that makes PostgreSQL crash-safe.

**The principle:** Before modifying any data page, write the intended change to the WAL on disk. If the server crashes, replay the WAL to recover.

```
Without WAL:                        With WAL:
1. Modify page in buffer pool       1. Write change to WAL (disk)
2. Flush page to disk               2. Confirm to client: "committed!"
   [CRASH here → data lost]         3. Modify page in buffer pool
                                    4. [CRASH here]
                                    5. On restart: replay WAL → data recovered
```

The WAL is sequential (append-only), so WAL writes are very fast. Data page writes can happen lazily in the background.

### WAL Files on Disk

```bash
ls $PGDATA/pg_wal/
# 000000010000000000000001
# 000000010000000000000002
# ...
```

Each WAL file is 16 MB. PostgreSQL archives old WAL files for replication and point-in-time recovery (PITR).

### Important WAL Configuration

```sql
-- Show WAL-related settings
SHOW wal_level;         -- minimal, replica, or logical
SHOW max_wal_size;      -- maximum WAL size before checkpoint
SHOW wal_buffers;       -- size of WAL buffer in RAM
SHOW synchronous_commit; -- controls WAL durability vs performance
```

`synchronous_commit = off` gives a performance boost (no fsync per transaction) at the cost of potentially losing the last few milliseconds of transactions on a crash. Useful for non-critical writes (logs, analytics).

---

## 4. VACUUM — Cleaning Up Dead Tuples

Because UPDATE creates new row versions, tables accumulate dead tuples over time. VACUUM reclaims this space.

```sql
-- Manual vacuum (reclaims space, doesn't return to OS)
VACUUM users;

-- Full vacuum (rewrites the table, returns space to OS, locks the table!)
VACUUM FULL users;

-- Analyze (updates statistics for the query planner)
VACUUM ANALYZE users;

-- See dead tuple counts
SELECT relname, n_live_tup, n_dead_tup, last_vacuum, last_autovacuum
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC;
```

### AUTOVACUUM

PostgreSQL runs AUTOVACUUM automatically in the background. It triggers when dead tuples exceed a threshold:

```
threshold = autovacuum_vacuum_threshold + autovacuum_vacuum_scale_factor × table_rows
          = 50 + 0.2 × 1,000,000
          = 200,050 dead tuples before autovacuum fires
```

For write-heavy tables, tune autovacuum to run more aggressively:

```sql
-- Per-table autovacuum tuning
ALTER TABLE user_events SET (
    autovacuum_vacuum_scale_factor = 0.01,  -- vacuum at 1% dead tuples
    autovacuum_vacuum_threshold = 100
);
```

### Transaction ID Wraparound

PostgreSQL uses 32-bit transaction IDs (XIDs). At ~2 billion transactions, the XID wraps around and rows look like they're from the future — which breaks visibility rules catastrophically.

VACUUM prevents this by advancing the `relfrozenxid` — marking old rows as "frozen" (visible to all transactions). PostgreSQL will **force** a VACUUM before XID wraparound by entering emergency mode. Keep AUTOVACUUM running and monitor `pg_database.datfrozenxid`.

---

## 5. Checkpoints — Writing Dirty Pages to Disk

PostgreSQL keeps recently modified pages in the shared buffer pool (RAM). A **checkpoint** writes all dirty (modified) pages to disk.

```
Shared Buffer Pool (RAM)
  [modified page A] → checkpoint → disk
  [modified page B] → checkpoint → disk
  [clean page C]    (no write needed)
```

After a checkpoint, PostgreSQL can truncate WAL files older than the checkpoint (the WAL records are no longer needed for recovery — the pages are safely on disk).

Checkpoints happen:
- Every `checkpoint_timeout` seconds (default: 5 minutes)
- When WAL reaches `max_wal_size` (default: 1 GB)

You can force a checkpoint:
```sql
CHECKPOINT;  -- use carefully in production; it's expensive
```

Monitor checkpoint frequency:
```sql
SELECT checkpoints_timed, checkpoints_req, buffers_checkpoint
FROM pg_stat_bgwriter;
```

If `checkpoints_req` (forced by WAL size) is high, increase `max_wal_size`.

---

## 6. Monitoring Internals in Real Time

PostgreSQL exposes rich internal state through system views.

### Active Queries

```sql
SELECT pid, usename, state, query, now() - query_start AS duration
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC;

-- Kill a slow query
SELECT pg_cancel_backend(pid);   -- sends SIGINT (gentle)
SELECT pg_terminate_backend(pid); -- sends SIGTERM (forceful)
```

### Lock Monitoring

```sql
-- See all locks currently held
SELECT
    locktype, relation::regclass, mode, granted, pid,
    left(query, 50) AS query
FROM pg_locks
JOIN pg_stat_activity USING (pid)
WHERE relation IS NOT NULL;

-- Find blocking queries
SELECT
    blocked.pid, blocked.query AS blocked_query,
    blocking.pid AS blocking_pid, blocking.query AS blocking_query
FROM pg_stat_activity AS blocked
JOIN pg_stat_activity AS blocking
    ON blocking.pid = ANY(pg_blocking_pids(blocked.pid))
WHERE blocked.wait_event_type = 'Lock';
```

### Table Statistics

```sql
-- Cache hit ratio (want > 99%)
SELECT
    sum(heap_blks_read) AS heap_read,
    sum(heap_blks_hit)  AS heap_hit,
    round(sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read) + 0.001) * 100, 2) AS ratio
FROM pg_statio_user_tables;

-- Index usage
SELECT
    indexrelname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

---

## 7. Connection Pooling with pgxpool

As mentioned in Chapter 13, PostgreSQL creates one process per connection. Each process uses ~5-10 MB. With 1000 concurrent connections, that's 5-10 GB just for overhead. The solution: **connection pooling**.

pgxpool (bundled with pgx) maintains a pool of persistent connections and multiplexes many Go goroutines across them.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    ctx := context.Background()

    config, err := pgxpool.ParseConfig("postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal(err)
    }

    // Pool configuration
    config.MaxConns = 25              // maximum connections
    config.MinConns = 5               // always keep 5 alive
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute
    config.HealthCheckPeriod = time.Minute

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Pool stats
    stats := pool.Stat()
    fmt.Printf("Pool: total=%d, idle=%d, acquired=%d\n",
        stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns())
}
```

For even larger scale, consider **PgBouncer** — a standalone connection pooler that sits between your app and PostgreSQL, supporting thousands of app connections with only a few dozen actual PostgreSQL connections.

---

## Summary

- PostgreSQL stores data as 8 KB **pages** in heap files. Each row (tuple) has `xmin`/`xmax` transaction IDs.
- **MVCC** stores multiple row versions. Readers see a consistent snapshot without blocking writers.
- The **WAL** writes changes to disk before modifying pages. Replay on restart achieves crash safety.
- **VACUUM** cleans dead tuples left behind by MVCC. AUTOVACUUM does this automatically.
- **Checkpoints** flush modified pages from RAM to disk, allowing old WAL to be truncated.
- Monitor with `pg_stat_activity` (active queries), `pg_locks` (lock contention), and `pg_stat_user_tables` (vacuum stats).

### Exercises

**Easy:** Run `SELECT xmin, xmax, * FROM users LIMIT 5` in psql. What do xmin and xmax represent for an active row?

**Medium:** Create a table, insert 100,000 rows, then delete 50,000 of them. Run `VACUUM ANALYZE` and then check `pg_stat_user_tables` to see `n_dead_tup` before and after.

**Hard:** Read the PostgreSQL documentation on `pg_stat_bgwriter`. Set up a pgbench workload and monitor checkpoint frequency. Tune `max_wal_size` and `checkpoint_completion_target` to reduce checkpoint I/O spikes.
