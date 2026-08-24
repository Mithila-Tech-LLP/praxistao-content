---
title: Optimistic vs Pessimistic Locking
category: Software & Programming
tags: [Databases, Concurrency]
duration: 7 min read
relatedCourses: [databases-and-async-go, go-programming]
relatedProjects: [key-value-store]
relatedTopics: [acid-vs-base-consistency-models, b-tree-indexes-explained]
---

## TL;DR

- **Pessimistic locking**: acquire a lock before touching the data, assuming a conflict is likely — other transactions wait or fail immediately.
- **Optimistic locking**: don't lock anything up front; assume conflicts are rare, and check at write time whether the data changed since you read it. If it did, reject and let the caller retry.
- The right choice depends on **contention**, not personal preference: frequent conflicts on the same row favor pessimistic; rare conflicts favor optimistic (since most optimistic attempts just succeed on the first try).
- Optimistic locking is usually implemented with a `version` column; pessimistic locking is usually `SELECT ... FOR UPDATE`.

## Pessimistic Locking

Grab an exclusive lock on a row before reading it for an update, so no other transaction can read-for-update (or in some isolation levels, even read) the same row until you're done:

```sql
BEGIN;
SELECT * FROM accounts WHERE id = 42 FOR UPDATE; -- locks this row
-- ... application logic decides the new balance ...
UPDATE accounts SET balance = balance - 100 WHERE id = 42;
COMMIT; -- lock released
```

Any other transaction trying to `SELECT ... FOR UPDATE` the same row blocks until this one commits or rolls back. This guarantees no other transaction can interleave a conflicting change — but it does so by making that other transaction wait, which is exactly the cost: under high contention on the same rows, transactions queue up behind each other, and a slow transaction holding the lock stalls everyone behind it.

## Optimistic Locking

Don't lock anything. Read the row along with a version marker (a `version` integer column, or the row's `updated_at` timestamp). When writing, check that the version is still what you read — if someone else changed the row in between, the write fails, and the caller decides whether to retry.

```sql
-- read
SELECT id, balance, version FROM accounts WHERE id = 42;
-- suppose we got: balance=500, version=7

-- write, checking the version hasn't moved
UPDATE accounts
SET balance = 400, version = version + 1
WHERE id = 42 AND version = 7;
-- if another transaction already updated this row (version is now 8),
-- this UPDATE matches zero rows — the application must detect that and retry
```

```go
result, err := db.Exec(
    "UPDATE accounts SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3",
    newBalance, accountID, expectedVersion,
)
rows, _ := result.RowsAffected()
if rows == 0 {
    // someone else updated this row first — reload and retry, or surface a conflict to the caller
    return ErrConflict
}
```

No lock is ever held while the application decides what to write — the only "check" happens atomically at the moment of the write itself, via the `WHERE version = $3` clause. This means other transactions never wait on this one; they just might occasionally have to retry if they collide.

## Choosing Based on Contention

The entire decision comes down to one question: **how often do multiple transactions actually try to modify the same row at the same time?**

- **High contention** (many transactions racing to update the same hot row — a popular product's inventory count, a single counter everyone increments): pessimistic locking is usually better. Under high contention, optimistic locking's retries themselves become frequent, and every failed attempt wasted the work of reading, computing, and attempting the write for nothing — that waste compounds under heavy contention. Pessimistic locking's queuing is more predictable and doesn't waste application-level work.
- **Low contention** (most rows are touched by at most one transaction at a time — most typical CRUD operations on individual user records): optimistic locking usually wins. Almost every write succeeds on the first attempt with zero lock-holding overhead; the rare conflict just costs one retry, which is cheaper overall than every transaction pessimistically locking rows that, in practice, almost never actually collide.

## The Retry Loop, Done Properly

Optimistic locking isn't complete without deciding what happens on a version mismatch. The typical pattern is a bounded retry loop:

```go
func updateBalanceOptimistic(db *sql.DB, accountID int, delta int) error {
    for attempt := 0; attempt < 3; attempt++ {
        var balance, version int
        db.QueryRow("SELECT balance, version FROM accounts WHERE id = $1", accountID).
            Scan(&balance, &version)

        newBalance := balance + delta
        result, _ := db.Exec(
            "UPDATE accounts SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3",
            newBalance, accountID, version,
        )
        if rows, _ := result.RowsAffected(); rows == 1 {
            return nil // success
        }
        // version mismatch — someone else updated it; loop and retry with fresh data
    }
    return errors.New("too many conflicting updates, giving up")
}
```

An unbounded retry loop, or one with no backoff, risks the same kind of retry-storm behavior seen in distributed systems generally under sustained contention — a bounded attempt count with a clear failure path back to the caller is the safer default.

## Common Pitfalls

- **Using optimistic locking under genuinely high contention** — this doesn't fail loudly, it just quietly wastes work: most attempts collide, retry, collide again, and overall throughput on that hot row ends up worse than a pessimistic lock would have given, while looking superficially "lock-free."
- **Holding a pessimistic lock across a slow operation** — `SELECT ... FOR UPDATE` followed by a slow external API call before the `UPDATE`/`COMMIT` holds the lock for the entire duration of that slow call, blocking every other transaction wanting the same row. Keep the locked section as short as possible — do slow work *before* acquiring the lock, not while holding it.
- **Forgetting to increment the version on every successful write** — if the version column isn't actually bumped, optimistic locking silently stops detecting conflicts at all; two concurrent writes can both succeed against the same stale version.
- **No retry logic for optimistic locking failures** — a version mismatch is an expected, routine outcome under any real contention, not an exceptional error; the caller needs an explicit retry (or a clear conflict response) rather than just surfacing a raw failure.
