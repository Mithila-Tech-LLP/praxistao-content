# Chapter 78: Transactions, Locking, and Concurrency Control

Two users buy the last ticket at the same time. A payment is debited but the order never commits. A report reads half-updated data. All of these are concurrency bugs in databases. PostgreSQL gives you a full suite of tools to prevent them — transactions, isolation levels, row locks, optimistic versioning, and advisory locks. This chapter covers each one with real Go code.

## Table of Contents

1. [ACID Properties](#1-acid-properties)
2. [Transaction Isolation Levels](#2-transaction-isolation-levels)
3. [Concurrency Anomalies](#3-concurrency-anomalies)
4. [Go Transaction Pattern](#4-go-transaction-pattern)
5. [Pessimistic Locking](#5-pessimistic-locking)
6. [Optimistic Locking](#6-optimistic-locking)
7. [Advisory Locks](#7-advisory-locks)
8. [Deadlock Detection and Avoidance](#8-deadlock-detection-and-avoidance)
9. [MVCC — How PostgreSQL Does It](#9-mvcc--how-postgresql-does-it)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. ACID Properties

Every database transaction must satisfy four guarantees:

```
ACID

A — Atomicity    All operations in a transaction succeed together,
                 or none of them do. No partial commits.
                 Analogy: a bank transfer debits $100 AND credits $100.
                 If the credit fails, the debit is also reversed.

C — Consistency  A transaction takes the database from one valid state
                 to another. Constraints (NOT NULL, UNIQUE, FK) are
                 checked at commit time.
                 Analogy: you can't overdraw an account below zero
                 if a CHECK(balance >= 0) constraint exists.

I — Isolation    Concurrent transactions behave as if they ran serially.
                 One transaction's in-progress changes are not visible
                 to others (at the right isolation level).
                 Analogy: two cashiers ringing up purchases at the
                 same time don't interfere with each other's drawers.

D — Durability   Once committed, a transaction survives crashes.
                 PostgreSQL writes to the WAL (Write-Ahead Log) before
                 acknowledging the commit.
                 Analogy: once the receipt prints, the sale is recorded
                 even if the power goes out one second later.
```

---

## 2. Transaction Isolation Levels

PostgreSQL supports four isolation levels. The default is **Read Committed**.

```
Level                  | What you get
───────────────────────┼─────────────────────────────────────────────────
Read Uncommitted       | Identical to Read Committed in PostgreSQL.
                       | PG never exposes uncommitted data.
───────────────────────┼─────────────────────────────────────────────────
Read Committed         | Each statement sees a fresh snapshot of
(PG default)           | committed data as of that statement's start.
                       | Two SELECT statements in the same transaction
                       | can return different rows.
───────────────────────┼─────────────────────────────────────────────────
Repeatable Read        | All statements in the transaction see the
                       | same snapshot (taken at transaction start).
                       | Rows read at the beginning are still there
                       | at the end. Phantom reads are also prevented
                       | in PostgreSQL (unlike the SQL standard).
───────────────────────┼─────────────────────────────────────────────────
Serializable           | Strongest guarantee. Transactions behave as if
                       | they executed one at a time in some serial order.
                       | PostgreSQL uses SSI (Serializable Snapshot
                       | Isolation) — low overhead, detects conflicts,
                       | may abort and retry on conflict.
```

Setting isolation level in Go:

```go
import (
    "context"
    "database/sql"
)

// Default (Read Committed)
tx, err := db.BeginTx(ctx, nil)

// Repeatable Read
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})

// Serializable
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})

// Read-only transaction (PostgreSQL can optimize these)
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
    ReadOnly:  true,
})
```

---

## 3. Concurrency Anomalies

Each isolation level protects against a different set of anomalies:

```
Anomaly              | Read Committed | Repeatable Read | Serializable
─────────────────────┼────────────────┼─────────────────┼─────────────
Dirty Read           | Protected      | Protected        | Protected
Non-Repeatable Read  | Possible       | Protected        | Protected
Phantom Read         | Possible       | Protected (PG)   | Protected
Serialization Anomaly| Possible       | Possible         | Protected
```

### Dirty read

Transaction A reads a row that Transaction B has modified but not yet committed. If B rolls back, A has read data that never existed. PostgreSQL never allows this at any level.

### Non-repeatable read

```
Transaction A                    Transaction B
────────────                     ────────────
SELECT balance FROM accounts
WHERE id = 1;  → 1000

                                 UPDATE accounts SET balance = 0
                                 WHERE id = 1;
                                 COMMIT;

SELECT balance FROM accounts
WHERE id = 1;  → 0   ← different result!
```

With **Repeatable Read**, Transaction A's second SELECT still returns 1000.

### Phantom read

```
Transaction A                    Transaction B
────────────                     ────────────
SELECT COUNT(*) FROM orders
WHERE status = 'pending';  → 5

                                 INSERT INTO orders (...) VALUES (...);
                                 -- status = 'pending'
                                 COMMIT;

SELECT COUNT(*) FROM orders
WHERE status = 'pending';  → 6  ← new row appeared (phantom)
```

PostgreSQL's Repeatable Read also prevents phantoms (the SQL standard says it doesn't have to).

### Serialization anomaly

Two transactions each read data and write based on what they read. Both reads happen before either write, so each transaction "sees" the old state. The combined result could never happen if they ran serially.

```
Transaction A: "if no doctor is on call, I'll take the shift"
Transaction B: "if no doctor is on call, I'll take the shift"

Both read: 0 doctors on call
Both write: I am on call

Result: 2 doctors on call — neither would have volunteered if they
        had seen the other sign up first.
```

**Serializable** prevents this. PostgreSQL detects the conflict and aborts one transaction with error `40001 (serialization_failure)` — the caller must retry.

---

## 4. Go Transaction Pattern

The canonical Go pattern for safe transactions:

```go
func Transfer(ctx context.Context, db *sql.DB, fromID, toID int64, amount float64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    // Deferred rollback is a safety net.
    // If Commit() succeeds, Rollback() is a no-op.
    // If anything returns early with an error, Rollback() cleans up.
    defer tx.Rollback()

    // Debit
    var balance float64
    err = tx.QueryRowContext(ctx,
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE",
        fromID,
    ).Scan(&balance)
    if err != nil {
        return fmt.Errorf("get sender: %w", err)
    }
    if balance < amount {
        return fmt.Errorf("insufficient funds: have %.2f, need %.2f", balance, amount)
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID,
    )
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }

    // Credit
    _, err = tx.ExecContext(ctx,
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID,
    )
    if err != nil {
        return fmt.Errorf("credit: %w", err)
    }

    // Only reaches here if both operations succeeded
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    return nil
}
```

Key points:
- `defer tx.Rollback()` immediately after `BeginTx` — always, no exceptions.
- `tx.Commit()` at the very end, only after all operations succeed.
- All queries inside the transaction use `tx`, not `db`. Using `db` would run outside the transaction.
- Return errors wrapped with `fmt.Errorf(...%w...)` so callers can inspect them.

### Helper to reduce boilerplate

```go
// WithTx runs fn inside a transaction. Commits if fn returns nil, rolls back otherwise.
func WithTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn func(*sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, opts)
    if err != nil {
        return fmt.Errorf("begin: %w", err)
    }
    defer tx.Rollback()

    if err := fn(tx); err != nil {
        return err
    }
    return tx.Commit()
}

// Usage
err := WithTx(ctx, db, nil, func(tx *sql.Tx) error {
    if _, err := tx.ExecContext(ctx, "UPDATE ..."); err != nil {
        return err
    }
    if _, err := tx.ExecContext(ctx, "INSERT ..."); err != nil {
        return err
    }
    return nil
})
```

---

## 5. Pessimistic Locking

Pessimistic locking assumes conflicts will happen and prevents them upfront by locking rows before reading them.

### SELECT ... FOR UPDATE

Acquires an exclusive row-level lock. Other transactions trying to lock the same row will block until this transaction commits or rolls back.

```go
func BookSeat(ctx context.Context, db *sql.DB, seatID, userID int64) error {
    return WithTx(ctx, db, nil, func(tx *sql.Tx) error {
        // Lock the seat row — no other transaction can lock it concurrently
        var status string
        err := tx.QueryRowContext(ctx,
            "SELECT status FROM seats WHERE id = $1 FOR UPDATE",
            seatID,
        ).Scan(&status)
        if err != nil {
            return fmt.Errorf("get seat: %w", err)
        }

        if status != "available" {
            return fmt.Errorf("seat %d is already %s", seatID, status)
        }

        _, err = tx.ExecContext(ctx,
            "UPDATE seats SET status = 'booked', user_id = $1 WHERE id = $2",
            userID, seatID,
        )
        return err
    })
}
```

### SELECT ... FOR UPDATE SKIP LOCKED

Used for job queues: each worker picks a job and skips rows already locked by other workers. No worker waits; no job is processed twice.

```go
func ClaimNextJob(ctx context.Context, tx *sql.Tx) (*Job, error) {
    var job Job
    err := tx.QueryRowContext(ctx, `
        SELECT id, payload, attempts
        FROM jobs
        WHERE status = 'pending'
          AND run_at <= NOW()
        ORDER BY run_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `).Scan(&job.ID, &job.Payload, &job.Attempts)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil // no jobs available right now
    }
    if err != nil {
        return nil, fmt.Errorf("claim job: %w", err)
    }

    // Mark as running
    _, err = tx.ExecContext(ctx,
        "UPDATE jobs SET status = 'running', started_at = NOW() WHERE id = $1",
        job.ID,
    )
    return &job, err
}

// Worker loop
func runWorker(ctx context.Context, db *sql.DB) {
    for {
        err := WithTx(ctx, db, nil, func(tx *sql.Tx) error {
            job, err := ClaimNextJob(ctx, tx)
            if err != nil { return err }
            if job == nil { return nil } // nothing to do

            if err := processJob(ctx, job); err != nil {
                // Mark failed — transaction still commits (we want to record the failure)
                _, _ = tx.ExecContext(ctx,
                    "UPDATE jobs SET status='failed', error=$1 WHERE id=$2",
                    err.Error(), job.ID,
                )
                return nil
            }

            _, err = tx.ExecContext(ctx,
                "UPDATE jobs SET status='done', finished_at=NOW() WHERE id=$1",
                job.ID,
            )
            return err
        })
        if err != nil { log.Printf("worker error: %v", err) }

        select {
        case <-ctx.Done(): return
        case <-time.After(1 * time.Second):
        }
    }
}
```

### SELECT ... FOR SHARE

Acquires a shared lock. Multiple transactions can hold a shared lock simultaneously, but none can acquire an exclusive lock (`FOR UPDATE`) while shared locks are held. Use it when you need to read and guarantee the row won't be deleted or exclusively locked by another transaction.

```go
// Read user, prevent concurrent deletion, but allow other readers
var user User
err := tx.QueryRowContext(ctx,
    "SELECT id, email, name FROM users WHERE id = $1 FOR SHARE",
    userID,
).Scan(&user.ID, &user.Email, &user.Name)
```

---

## 6. Optimistic Locking

Optimistic locking assumes conflicts are rare and does not lock rows upfront. Instead, it detects conflicts at write time using a `version` column.

```sql
-- Migration: add version column
ALTER TABLE products ADD COLUMN IF NOT EXISTS version BIGINT NOT NULL DEFAULT 1;
```

```go
type Product struct {
    ID      int64
    Name    string
    Price   float64
    Version int64  // read this, pass it back on update
}

func UpdateProduct(ctx context.Context, db *sql.DB, p *Product) error {
    result, err := db.ExecContext(ctx, `
        UPDATE products
        SET name = $1, price = $2, version = version + 1, updated_at = NOW()
        WHERE id = $3
          AND version = $4   -- only update if no one else changed it
    `, p.Name, p.Price, p.ID, p.Version)
    if err != nil {
        return fmt.Errorf("update product: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("rows affected: %w", err)
    }
    if rows == 0 {
        // Either the product doesn't exist OR someone else updated it first
        return ErrConflict
    }

    p.Version++ // update caller's copy
    return nil
}

var ErrConflict = errors.New("conflict: record was modified by another transaction")

// Caller retries on conflict
func UpdateProductWithRetry(ctx context.Context, db *sql.DB, id int64, applyFn func(*Product)) error {
    const maxRetries = 3
    for attempt := range maxRetries {
        // Fresh read each attempt — get the latest version
        product, err := GetProduct(ctx, db, id)
        if err != nil { return err }

        applyFn(product) // let the caller set the new values

        err = UpdateProduct(ctx, db, product)
        if errors.Is(err, ErrConflict) {
            if attempt == maxRetries-1 {
                return fmt.Errorf("update failed after %d attempts: %w", maxRetries, err)
            }
            time.Sleep(time.Duration(attempt+1) * 50 * time.Millisecond)
            continue
        }
        return err
    }
    return nil
}
```

### Pessimistic vs Optimistic: when to use each

```
Pessimistic (FOR UPDATE)          Optimistic (version column)
────────────────────────          ───────────────────────────
High conflict rate                Low conflict rate
Short critical section            Long think time between read and write
Simple: no retry logic            Requires retry on conflict
Blocks other transactions         Never blocks — reads are always free
Good for: job queues, seats,      Good for: profile edits, config
  inventory with high contention    updates, low-traffic records
```

---

## 7. Advisory Locks

Advisory locks are application-level locks stored in PostgreSQL's memory — not tied to any row or table. They are useful for things like ensuring only one instance of a cron job runs at a time.

```go
// pg_try_advisory_lock(key) — non-blocking, returns true if lock acquired
func TryAdvisoryLock(ctx context.Context, db *sql.DB, key int64) (bool, error) {
    var acquired bool
    err := db.QueryRowContext(ctx,
        "SELECT pg_try_advisory_lock($1)", key,
    ).Scan(&acquired)
    return acquired, err
}

// pg_advisory_unlock(key) — releases the lock
func AdvisoryUnlock(ctx context.Context, db *sql.DB, key int64) error {
    _, err := db.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", key)
    return err
}

// Usage: ensure only one instance of a maintenance job runs
const maintenanceJobKey = 12345 // any int64 your team agrees on

func RunMaintenanceJob(ctx context.Context, db *sql.DB) error {
    acquired, err := TryAdvisoryLock(ctx, db, maintenanceJobKey)
    if err != nil {
        return fmt.Errorf("advisory lock: %w", err)
    }
    if !acquired {
        log.Println("maintenance job already running on another instance, skipping")
        return nil
    }
    defer AdvisoryUnlock(ctx, db, maintenanceJobKey)

    return doMaintenance(ctx, db)
}
```

### Session vs transaction advisory locks

```go
// Session-scoped: lock held until explicitly released or connection closes
db.ExecContext(ctx, "SELECT pg_advisory_lock($1)", key)           // blocks until acquired
db.ExecContext(ctx, "SELECT pg_try_advisory_lock($1)", key)       // returns immediately
db.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", key)         // release

// Transaction-scoped: released automatically when transaction ends
tx.ExecContext(ctx, "SELECT pg_try_advisory_xact_lock($1)", key)  // released at COMMIT/ROLLBACK
```

Transaction-scoped locks are safer — you cannot forget to release them. Session-scoped locks are necessary when the lock must outlive a single transaction (e.g., a long-running background process).

---

## 8. Deadlock Detection and Avoidance

A deadlock occurs when two transactions each hold a lock the other needs:

```
Transaction A                    Transaction B
────────────                     ────────────
LOCK row users WHERE id=1        LOCK row users WHERE id=2
(waiting for id=2) ──────────────(waiting for id=1)
```

Both are stuck forever. PostgreSQL detects this within a second and aborts one transaction with error code `40P01 (deadlock_detected)`.

```go
func Transfer(ctx context.Context, db *sql.DB, fromID, toID int64, amount float64) error {
    // BAD: if two concurrent transfers go in opposite directions,
    // they can deadlock:
    //   Transfer(A→B) locks A, then tries to lock B
    //   Transfer(B→A) locks B, then tries to lock A
    //   → deadlock

    // GOOD: always lock in the same order (lower ID first)
    // This ensures both transactions acquire locks in the same sequence,
    // so one always finishes first and the other proceeds without conflict.
    first, second := fromID, toID
    if fromID > toID {
        first, second = toID, fromID
    }

    return WithTx(ctx, db, nil, func(tx *sql.Tx) error {
        // Lock both rows in a consistent order
        _, err := tx.ExecContext(ctx, `
            SELECT id FROM accounts
            WHERE id IN ($1, $2)
            ORDER BY id
            FOR UPDATE
        `, first, second)
        if err != nil { return err }

        // Now update safely
        if _, err := tx.ExecContext(ctx,
            "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
            amount, fromID,
        ); err != nil { return err }

        _, err = tx.ExecContext(ctx,
            "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
            amount, toID,
        )
        return err
    })
}
```

### Handling deadlock errors in Go

```go
import "github.com/lib/pq"

func isDeadlock(err error) bool {
    var pqErr *pq.Error
    return errors.As(err, &pqErr) && pqErr.Code == "40P01"
}

func isSerializationFailure(err error) bool {
    var pqErr *pq.Error
    return errors.As(err, &pqErr) && pqErr.Code == "40001"
}

func WithRetryOnConflict(ctx context.Context, fn func() error) error {
    const maxAttempts = 5
    for attempt := range maxAttempts {
        err := fn()
        if err == nil { return nil }
        if isDeadlock(err) || isSerializationFailure(err) {
            if attempt == maxAttempts-1 { return err }
            // Exponential backoff with jitter
            delay := time.Duration(attempt+1)*50*time.Millisecond +
                time.Duration(rand.Intn(50))*time.Millisecond
            select {
            case <-time.After(delay):
            case <-ctx.Done(): return ctx.Err()
            }
            continue
        }
        return err // non-retryable error
    }
    return nil
}
```

---

## 9. MVCC — How PostgreSQL Does It

MVCC (Multi-Version Concurrency Control) is how PostgreSQL lets readers and writers run concurrently without blocking each other. Understanding it explains why the isolation levels work the way they do.

```
Without MVCC (naive locking):         With MVCC (PostgreSQL):
──────────────────────────────        ──────────────────────────────
Reader waits for Writer                Reader sees old version of row
Writer waits for Reader                Writer creates new version of row
High contention = low throughput       Old versions cleaned up by VACUUM
```

Every row in PostgreSQL has two hidden system columns:
- `xmin`: the transaction ID that created this row version
- `xmax`: the transaction ID that deleted this row version (0 if alive)

```
accounts table:
id  | balance | xmin | xmax
────┼─────────┼──────┼─────
1   |  1000   |  50  |   0   ← visible to everyone: created by tx 50, not deleted
1   |   900   |  75  |   0   ← also exists: created by tx 75 (UPDATE = delete old + insert new)
```

When Transaction A (started before tx 75) reads `accounts WHERE id=1`, it sees `xmin=50` (balance=1000) because `xmin=75` is newer than its snapshot. Transaction B (started after tx 75) sees `xmin=75` (balance=900).

The old version (`xmin=50, xmax=75`) is kept around until no active transaction can see it anymore, then `VACUUM` removes it.

```
Key consequences:
  • Reads never block writes
  • Writes never block reads
  • An UPDATE is internally: mark old version as deleted (set xmax) +
    insert new version (set xmin)
  • VACUUM removes "dead" row versions (xmax is set and no transaction
    can see them anymore)
  • Table bloat occurs if VACUUM doesn't run or rows are updated heavily
```

---

## Summary

- **ACID**: Atomicity (all-or-nothing), Consistency (constraints hold), Isolation (concurrent = serial), Durability (committed = permanent)
- **Isolation levels**: Read Committed (default, two SELECTs may differ), Repeatable Read (same snapshot throughout), Serializable (no anomalies, may abort)
- **Go pattern**: `BeginTx` → `defer tx.Rollback()` → do work → `tx.Commit()`
- **Pessimistic locking**: `FOR UPDATE` (exclusive), `FOR UPDATE SKIP LOCKED` (job queues), `FOR SHARE` (shared)
- **Optimistic locking**: `version` column; `UPDATE WHERE version=$old_version`; check `RowsAffected() == 0` means conflict
- **Advisory locks**: application-level locks via `pg_try_advisory_lock(key)`; ideal for cron jobs and migration gates
- **Deadlocks**: always acquire locks in the same order; catch `40P01` and retry
- **MVCC**: readers see a consistent snapshot of old row versions; writers create new versions; no reader/writer blocking

## Exercises

### Easy
1. Write a `Transfer(ctx, db, fromID, toID, amount)` function using `BeginTx`, `defer tx.Rollback()`, and explicit `tx.Commit()`. Test it: transfer $100 from account 1 to account 2, verify balances before and after.
2. Implement a `GetProductForUpdate(ctx, tx, id)` function that uses `SELECT ... FOR UPDATE` and returns an error if the product is not found. Call it inside a transaction, modify the price, then update.
3. Add a `version BIGINT DEFAULT 1` column to a `products` table. Write `UpdateProductOptimistic(ctx, db, product)` that only updates if `version` matches. Return a custom `ErrConflict` when `RowsAffected() == 0`.

### Medium
4. Build a **job queue** backed by a PostgreSQL table with columns `id, payload, status, run_at, started_at, finished_at`. Use `SELECT ... FOR UPDATE SKIP LOCKED` to claim jobs. Run 5 concurrent worker goroutines and verify no job is processed twice (use a mutex-protected map in tests to track job IDs).
5. Simulate a **serialization failure** by running two goroutines that both read a count, then write based on the count (the on-call doctor scenario from Section 3). Use `Serializable` isolation. One goroutine should receive a `40001` error. Add retry logic using `WithRetryOnConflict`.
6. Write an `AdvisoryLock` helper that acquires a session-level advisory lock and returns a `release func()`. Use it in a cron job function. Simulate two concurrent goroutines calling the cron job and verify only one actually runs (the other returns early).

### Hard
7. Implement a **full optimistic-locking retry loop** for a high-conflict scenario: 50 goroutines all try to increment the same counter (a `counter` table with `value` and `version` columns). No `FOR UPDATE`, no serializable isolation — pure optimistic locking with retry. Measure: how many retries occur total? At what concurrency does the retry rate become unacceptable?
8. Build a **transaction outbox pattern**: inside a single transaction, insert a `domain_events` row alongside the main business data. A separate goroutine polls `domain_events WHERE published = FALSE FOR UPDATE SKIP LOCKED`, publishes each event to a message queue, then marks it `published = TRUE`. This guarantees at-least-once delivery without distributed transactions. Implement both sides and write an integration test that kills the publisher mid-flight and verifies unpublished events are retried on restart.
