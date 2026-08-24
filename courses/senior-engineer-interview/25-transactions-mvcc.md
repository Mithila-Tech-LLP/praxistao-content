# Chapter 25: Transactions & MVCC — Isolation Levels, Locking & Deadlocks

Transactions are what make databases trustworthy. Understanding MVCC, isolation levels, and locking deeply is essential for senior backend interviews — especially for database-intensive systems at top companies.

## Table of Contents

1. [ACID Revisited](#1-acid-revisited)
2. [Transaction Isolation Levels](#2-transaction-isolation-levels)
3. [MVCC — How PostgreSQL Avoids Read Locks](#3-mvcc--how-postgresql-avoids-read-locks)
4. [Row-Level Locking](#4-row-level-locking)
5. [Deadlocks in Databases](#5-deadlocks-in-databases)
6. [Optimistic vs Pessimistic Locking](#6-optimistic-vs-pessimistic-locking)
7. [Distributed Transactions](#7-distributed-transactions)
8. [Go Code Patterns for Transactions](#8-go-patterns-for-transactions)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. ACID Revisited

| Property | What It Means | What Breaks Without It |
|---|---|---|
| **Atomicity** | All changes succeed or all are rolled back | Half-completed transfers: balance debited but not credited |
| **Consistency** | Constraints are never violated | Foreign key violation, negative balance |
| **Isolation** | Concurrent transactions don't interfere | Dirty reads, lost updates, phantom reads |
| **Durability** | Committed data survives crashes | Data loss on server restart |

---

## 2. Transaction Isolation Levels

Isolation levels control what anomalies can occur when transactions run concurrently.

### The Four Anomalies

```
DIRTY READ: reading uncommitted changes from another transaction
  T1 writes X=10 (not committed)
  T2 reads X=10 (dirty read)
  T1 rolls back: T2 saw a value that never existed!

NON-REPEATABLE READ: reading the same row twice gives different results
  T1 reads X=5
  T2 updates X=10 and commits
  T1 reads X=10 — different from first read!

PHANTOM READ: a query returns different rows when run twice
  T1: SELECT * FROM orders WHERE total > 100 → 5 rows
  T2: INSERT INTO orders(total=200) and commits
  T1: SELECT * WHERE total > 100 → 6 rows (phantom appeared!)

LOST UPDATE: two transactions overwrite each other's writes
  T1 reads X=5, plans to set X=6
  T2 reads X=5, sets X=7, commits
  T1 sets X=6, commits — T2's update is LOST
```

### Isolation Levels and What They Prevent

| Isolation Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| READ UNCOMMITTED | Possible | Possible | Possible |
| READ COMMITTED | Prevented | Possible | Possible |
| REPEATABLE READ | Prevented | Prevented | Possible* |
| SERIALIZABLE | Prevented | Prevented | Prevented |

*PostgreSQL's REPEATABLE READ also prevents phantom reads due to MVCC.

```sql
-- Set isolation level for a transaction
BEGIN;
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
-- ... your queries ...
COMMIT;

-- Or per connection:
SET default_transaction_isolation = 'repeatable read';
```

**Default:** PostgreSQL defaults to READ COMMITTED — each statement sees data committed before that statement began.

---

## 3. MVCC — How PostgreSQL Avoids Read Locks

MVCC (Multi-Version Concurrency Control) is the key to PostgreSQL's high concurrency: **readers never block writers, and writers never block readers**.

### How It Works

Every row in PostgreSQL has two system columns:
- `xmin`: transaction ID that inserted this row
- `xmax`: transaction ID that deleted this row (0 if still alive)

```
Row in heap:   [xmin=100] [xmax=0] {name: "Alice", salary: 50000}

T1 (txid=150) updates salary:
  New row:     [xmin=150] [xmax=0]   {name: "Alice", salary: 60000}
  Old row:     [xmin=100] [xmax=150] {name: "Alice", salary: 50000}  ← marked as deleted by T1

T2 (txid=140) reads BEFORE T1 commits:
  T2's snapshot: sees transactions committed before 140
  T1 has txid=150 > 140, so T2 sees the OLD row (salary: 50000)
  
T2 reads AFTER T1 commits:
  T2's new statement (READ COMMITTED mode) has a new snapshot
  T2 sees the NEW row (salary: 60000)
```

MVCC means read queries don't need to acquire locks on rows — they just read the correct version for their snapshot. This is why `SELECT` never blocks in PostgreSQL.

### VACUUM and Dead Tuples

Old row versions ("dead tuples") accumulate and must be cleaned up. VACUUM reclaims this space:

```sql
-- Manual vacuum:
VACUUM ANALYZE users;

-- Check dead tuple accumulation:
SELECT relname, n_dead_tup, n_live_tup FROM pg_stat_user_tables ORDER BY n_dead_tup DESC;

-- Autovacuum runs automatically based on thresholds:
-- autovacuum_vacuum_threshold = 50 (rows)
-- autovacuum_vacuum_scale_factor = 0.2 (20% of table)
-- Vacuum triggers when: dead_tup > threshold + scale_factor * live_tup
```

---

## 4. Row-Level Locking

PostgreSQL provides row-level locks for explicit synchronization:

```sql
-- SELECT FOR UPDATE: exclusive lock on selected rows
-- Blocks other SELECT FOR UPDATE/SHARE, and updates to these rows
BEGIN;
SELECT * FROM seats WHERE id = 5 FOR UPDATE;  -- lock this seat
-- Check if seat is available
UPDATE seats SET status = 'booked' WHERE id = 5;
COMMIT;

-- SELECT FOR SHARE: shared lock — allows other FOR SHARE but blocks FOR UPDATE
-- Use when: you need to read a row and prevent it from being deleted, but don't need to update it
SELECT * FROM prices WHERE product_id = 1 FOR SHARE;

-- SKIP LOCKED: skip rows that are locked by other transactions
-- Use for task queues — pick available work without waiting
SELECT * FROM jobs WHERE status = 'pending' LIMIT 1 FOR UPDATE SKIP LOCKED;
-- This pattern avoids contention in concurrent queue consumers
```

### Advisory Locks

For application-level locking:

```sql
-- Acquire an advisory lock on an arbitrary integer key
SELECT pg_advisory_lock(12345);  -- blocks until acquired
-- ... do work requiring distributed lock ...
SELECT pg_advisory_unlock(12345);

-- Non-blocking try:
SELECT pg_try_advisory_lock(12345);  -- returns true if acquired, false if already locked
```

---

## 5. Deadlocks in Databases

Database deadlocks happen the same way as in application code: two transactions each hold a lock and wait for the other's.

```sql
-- Deadlock example:
-- T1: UPDATE accounts SET balance -= 100 WHERE id = 1;  (locks row 1)
-- T2: UPDATE accounts SET balance -= 100 WHERE id = 2;  (locks row 2)
-- T1: UPDATE accounts SET balance += 100 WHERE id = 2;  (waits for T2 to release row 2)
-- T2: UPDATE accounts SET balance += 100 WHERE id = 1;  (waits for T1 to release row 1)
-- DEADLOCK: PostgreSQL detects and kills one transaction

-- PostgreSQL error:
-- ERROR: deadlock detected
-- DETAIL: Process 1234 waits for ShareLock on transaction 5678; blocked by process 5678.
--         Process 5678 waits for ShareLock on transaction 1234; blocked by process 1234.

-- FIX: always lock rows in the same order
-- Both transactions should update row 1 first, then row 2
BEGIN;
UPDATE accounts SET balance -= 100 WHERE id = LEAST(from_id, to_id);
UPDATE accounts SET balance += 100 WHERE id = GREATEST(from_id, to_id);
COMMIT;
```

---

## 6. Optimistic vs Pessimistic Locking

**Pessimistic locking:** Acquire a lock before reading, hold it until done. Use when conflicts are frequent.

```sql
-- Pessimistic: lock the row before reading to ensure nobody else modifies it
BEGIN;
SELECT * FROM inventory WHERE product_id = 1 FOR UPDATE;
-- Nobody else can update this row until we COMMIT
UPDATE inventory SET quantity -= 1 WHERE product_id = 1;
COMMIT;
```

**Optimistic locking:** Don't lock. Read, compute, then verify nothing changed before writing. If something changed, retry. Use when conflicts are rare.

```sql
-- Optimistic: use a version column
-- Read:
SELECT id, quantity, version FROM inventory WHERE product_id = 1;
-- Returns: id=1, quantity=10, version=5

-- Update: include version in WHERE — fails if version changed
UPDATE inventory 
SET quantity = quantity - 1, version = version + 1
WHERE product_id = 1 AND version = 5;
-- Returns: 0 rows updated → conflict detected → retry

-- In Go:
func decrementInventory(ctx context.Context, db *sql.DB, productID int) error {
    for attempt := 0; attempt < 3; attempt++ {
        var qty, version int
        db.QueryRowContext(ctx, "SELECT quantity, version FROM inventory WHERE product_id = $1", productID).
            Scan(&qty, &version)
        
        result, _ := db.ExecContext(ctx,
            "UPDATE inventory SET quantity = $1, version = $2 WHERE product_id = $3 AND version = $4",
            qty-1, version+1, productID, version)
        
        affected, _ := result.RowsAffected()
        if affected > 0 { return nil } // success
        // Conflict: someone else updated — retry
    }
    return errors.New("conflict after 3 retries")
}
```

---

## 7. Distributed Transactions

Distributed transactions span multiple databases or services. This is fundamentally harder than single-database transactions.

**Two-Phase Commit (2PC):**
1. **Prepare phase:** Coordinator asks all participants "can you commit?"
2. **Commit phase:** If all say yes, coordinator sends "commit" to all

```
Coordinator
  ↓ PREPARE
Participant A: "yes, I can commit"
Participant B: "yes, I can commit"
  ↓ COMMIT
Participant A: committed
Participant B: committed
```

**The problem with 2PC:** If the coordinator fails after some participants committed but before others, you're stuck in an inconsistent state. This is why distributed transactions are rarely used in practice.

**Modern approach:** Sagas — a sequence of local transactions with compensating actions for rollback.

---

## 8. Go Patterns for Transactions

```go
// Standard transaction pattern in Go
func transferMoney(ctx context.Context, db *sql.DB, from, to string, amount float64) error {
    tx, err := db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable, // or sql.LevelReadCommitted
    })
    if err != nil { return err }
    
    // ALWAYS defer rollback — it's a no-op if transaction was already committed
    defer tx.Rollback()
    
    // Debit
    _, err = tx.ExecContext(ctx, "UPDATE accounts SET balance -= $1 WHERE id = $2", amount, from)
    if err != nil { return fmt.Errorf("debit: %w", err) }
    
    // Credit
    _, err = tx.ExecContext(ctx, "UPDATE accounts SET balance += $1 WHERE id = $2", amount, to)
    if err != nil { return fmt.Errorf("credit: %w", err) }
    
    return tx.Commit()
}

// Transaction helper to reduce boilerplate
func withTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback() // no-op if committed
    
    if err := fn(tx); err != nil { return err }
    return tx.Commit()
}

// Usage:
err := withTransaction(ctx, db, func(tx *sql.Tx) error {
    _, err := tx.ExecContext(ctx, "UPDATE ...")
    if err != nil { return err }
    _, err = tx.ExecContext(ctx, "INSERT ...")
    return err
})
```

---

## 9. Interview Questions & Model Answers

**Q: What is MVCC and how does it improve database concurrency?**

"MVCC (Multi-Version Concurrency Control) keeps multiple versions of each row simultaneously. When a row is updated, PostgreSQL creates a new version of the row rather than overwriting the old one. Readers see the version of the row that was current at the start of their transaction (or statement, in READ COMMITTED mode). This means readers never need to acquire locks on rows — they just read the appropriate version. Writers only conflict with other writers. The result: reads and writes can proceed concurrently, and a long-running read transaction doesn't block writes."

**Q: What is the difference between SELECT FOR UPDATE and optimistic locking?**

"SELECT FOR UPDATE is pessimistic — it acquires an exclusive lock on rows immediately and holds it until the transaction commits. Other transactions that try to update those rows must wait. This prevents lost updates but can cause contention. Optimistic locking doesn't acquire a lock — it reads the current version, does its computation, and then updates with a WHERE clause that includes the expected version. If the row was changed by someone else (version mismatch), 0 rows are updated and the application retries. Pessimistic locking is safer for high-contention scenarios; optimistic is better for read-heavy workloads where conflicts are rare."

---

## Summary

- ACID: Atomicity (all-or-nothing), Consistency (constraints hold), Isolation (concurrency anomalies prevented), Durability (crashes don't lose committed data).
- Four anomalies: dirty read, non-repeatable read, phantom read, lost update.
- PostgreSQL default: READ COMMITTED. Serializable prevents all anomalies.
- **MVCC:** multiple row versions — readers see the version current at their snapshot, writers create new versions. No read locks needed.
- **VACUUM:** reclaims space from dead row versions created by MVCC.
- **SELECT FOR UPDATE:** exclusive lock for pessimistic concurrency. SKIP LOCKED for queue patterns.
- **Optimistic locking:** version column, retry on conflict. Good for read-heavy workloads.
- In Go: always `defer tx.Rollback()`, it's a no-op if already committed.
