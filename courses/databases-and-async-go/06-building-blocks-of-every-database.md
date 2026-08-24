# Chapter 06: The Building Blocks of Every Database

Before we touch PostgreSQL, MongoDB, or any real database, we need to understand the four concepts that every database is built on. These aren't PostgreSQL-specific or MySQL-specific — they are universal. Every database you use in this course (and in your career) implements some version of what you'll learn here.

## Table of Contents

1. Transactions — Doing Things Atomically
2. ACID — The Four Promises Every Database Makes
3. Concurrency — What Happens When Everyone Writes at Once
4. The CAP Theorem — The Fundamental Trade-off
5. Quick Reference Card
6. Exercises

---

## 1. Transactions — Doing Things Atomically

Imagine you are transferring $100 from your bank account to a friend's account. Two things must happen:

1. Subtract $100 from your account
2. Add $100 to your friend's account

What if the power goes out after step 1 but before step 2? Your $100 is gone, but your friend never received it. The money disappeared.

A **transaction** is a group of operations that either all succeed or all fail — together, as one unit. You cannot have step 1 succeed and step 2 fail.

```sql
BEGIN;  -- start the transaction

UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

COMMIT; -- apply both changes, or ROLLBACK to undo everything
```

If anything goes wrong between `BEGIN` and `COMMIT`, you call `ROLLBACK` and the database undoes all changes as if nothing happened. The bank transfer either completes fully or not at all.

In Go with `database/sql`:

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq" // PostgreSQL driver
)

func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    // Begin a transaction
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    // If anything goes wrong, roll back automatically
    defer tx.Rollback()

    // Debit from sender
    _, err = tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
    if err != nil {
        return err // Rollback() will be called by defer
    }

    // Credit to receiver
    _, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
    if err != nil {
        return err
    }

    // Everything succeeded — commit
    return tx.Commit()
}
```

The `defer tx.Rollback()` is a safety net. If `Commit()` is called successfully, `Rollback()` is a no-op. If anything fails before `Commit()`, the rollback fires automatically.

---

## 2. ACID — The Four Promises Every Database Makes

ACID stands for Atomicity, Consistency, Isolation, and Durability. These are the four guarantees that SQL databases (and many NoSQL databases) provide.

### Atomicity — All or Nothing

We already covered this with transactions. Atomicity means a transaction either completes fully or has no effect at all. There is no "partial success."

Real-world analogy: A vending machine either gives you the snack AND takes your money, or it returns your money and gives you nothing. It never takes your money without giving the snack.

### Consistency — The Database Is Always Valid

Consistency means the database is always in a valid state — all rules and constraints are satisfied, before and after every transaction.

If you have a rule "every order must have a customer", the database will reject any INSERT that violates this. You cannot end up with orphaned data.

```sql
-- This would fail if customer_id 999 doesn't exist
INSERT INTO orders (customer_id, total) VALUES (999, 50.00);
-- ERROR: foreign key constraint violation
```

The database enforces your rules, so your application doesn't have to.

### Isolation — Transactions Don't Step on Each Other

Isolation means concurrent transactions don't see each other's intermediate (uncommitted) state. Each transaction runs as if it is the only one running.

Without isolation, you could read data that another transaction is in the middle of modifying — data that might be rolled back moments later.

```
Transaction A: reads Alice's balance ($500)
Transaction B: starts transferring $500 OUT of Alice's account
Transaction A: reads Alice's balance again... is it $500 or $0?
```

Isolation levels control how strictly this is enforced. We'll go deep on isolation in Chapter 11.

### Durability — Committed Data Survives Crashes

Durability means once a transaction is committed, it stays committed — even if the server crashes immediately after.

The database achieves this using the **Write-Ahead Log (WAL)**. Before any page is changed, the change is first written to the WAL on disk. If the server crashes, it replays the WAL on restart to recover committed transactions.

```
Without durability:
1. User inserts a row
2. Database confirms "inserted!"
3. Power goes out
4. Database restarts — the row is GONE

With durability:
1. User inserts a row
2. Database writes the change to the WAL (disk)
3. Database confirms "inserted!"
4. Power goes out
5. Database replays WAL on restart — the row is THERE
```

---

## 3. Concurrency — What Happens When Everyone Writes at Once

A production database might have 10,000 users querying and updating simultaneously. Without careful coordination, they'd corrupt each other's data.

### The Lost Update Problem

```
Time → →
User A: reads stock count = 10
User B: reads stock count = 10
User A: sets stock count = 9  (bought 1 item)
User B: sets stock count = 9  (bought 1 item)
Result: stock count = 9, but 2 items were sold!
```

Both users read the same value, both decremented by 1, and one update was lost.

### Locks

The simplest solution: locks. Before modifying data, grab an exclusive lock on it. Others must wait.

```sql
-- SELECT FOR UPDATE locks the row, preventing others from modifying it
BEGIN;
SELECT stock FROM products WHERE id = 5 FOR UPDATE;
-- At this point, no other transaction can UPDATE this row
UPDATE products SET stock = stock - 1 WHERE id = 5;
COMMIT;
```

Locks work but slow things down — everyone waits for the lock holder to finish.

### MVCC — Reads Never Block Writes

Modern databases use **Multi-Version Concurrency Control (MVCC)**. Instead of blocking reads when writing, the database keeps multiple versions of each row.

```
Row: user_id=1, name="Alice"

Transaction A updates name to "Alicia"
↓
Database stores TWO versions:
Version 1 (old): name="Alice",  visible to transactions started before the update
Version 2 (new): name="Alicia", visible to transactions started after the update

Other transactions reading while A is running see the OLD version.
After A commits, new transactions see the NEW version.
```

MVCC means **reads never block writes and writes never block reads**. This is why PostgreSQL can handle thousands of concurrent transactions without constant lock contention.

### Deadlocks

When two transactions each hold a lock and each wait for the other's lock, you get a **deadlock**.

```
Transaction A: locks Row 1, waits for Row 2
Transaction B: locks Row 2, waits for Row 1
→ Neither can proceed. Deadlock!
```

Databases detect deadlocks automatically and kill one transaction (returning an error), letting the other proceed. Your application should retry on deadlock errors.

---

## 4. The CAP Theorem — The Fundamental Trade-off

The CAP theorem says: in a distributed system, you can only guarantee two of these three properties:

- **C — Consistency**: Every read sees the most recent write (or an error)
- **A — Availability**: Every request gets a response (not an error)
- **P — Partition Tolerance**: The system works even when the network splits

Since network partitions are inevitable in distributed systems (cables get unplugged, servers crash), you always have P. The real choice is between C and A.

**CP systems** (choose consistency): If two parts of the cluster can't communicate, they refuse to serve requests rather than risk returning stale data. Example: PostgreSQL (single node is consistent), ZooKeeper, etcd.

**AP systems** (choose availability): If partitioned, keep serving requests even if data might be stale. Accept that different nodes might have different views. Example: Cassandra, DynamoDB, CouchDB.

```
                  CAP Triangle
                  
                 Consistency
                     /\
                    /  \
                   /    \
                  / CA   \
                 /--------\
        CP      /          \   AP
               /            \
         Availability ---- Partition
                            Tolerance
         
         CA = single node databases (can't scale across network)
         CP = ZooKeeper, etcd, HBase
         AP = Cassandra, DynamoDB, Riak
```

**For this course:** Most SQL databases (PostgreSQL, MySQL) are essentially CP — they sacrifice availability for consistency. Redis and Cassandra sacrifice consistency for availability.

---

## 5. Quick Reference Card

| Concept | One-Line Explanation | Go Code Pattern |
|---------|---------------------|-----------------|
| Transaction | Group of operations that all succeed or all fail | `tx, _ := db.Begin()` ... `tx.Commit()` |
| Atomicity | All-or-nothing execution | `defer tx.Rollback()` as safety net |
| Consistency | Database rules are always satisfied | Foreign keys, NOT NULL, CHECK constraints |
| Isolation | Transactions don't see each other's work-in-progress | `SET TRANSACTION ISOLATION LEVEL` |
| Durability | Committed data survives crashes | WAL (write-ahead log) on disk |
| Lock | Exclusive access to a row during a transaction | `SELECT ... FOR UPDATE` |
| MVCC | Multiple versions of rows for concurrent access | Built into PostgreSQL, MySQL InnoDB |
| CAP Theorem | C + A + P: pick two (P is always required) | Architecture choice, not code |

---

## Summary

- A **transaction** groups multiple operations into one atomic unit — all succeed or all roll back.
- **ACID** = Atomicity (all-or-nothing), Consistency (rules always hold), Isolation (transactions don't interfere), Durability (committed data survives crashes).
- **Concurrency problems** include lost updates and dirty reads — solved by locks and MVCC.
- **MVCC** lets reads and writes happen simultaneously without blocking each other.
- The **CAP theorem** says distributed databases must choose between consistency and availability.

### Quick Check

1. What happens to a transaction if the server crashes between BEGIN and COMMIT?
2. What is the "lost update" problem and how do locks solve it?
3. In the CAP theorem, why is P (partition tolerance) always required in distributed systems?

### Exercises

**Easy:** Write a Go function that wraps two SQL UPDATE statements in a transaction. Use `defer tx.Rollback()` as a safety net.

**Medium:** Simulate a lost update without transactions: write a Go program where two goroutines read and then write back an integer value from an in-memory counter, causing one update to be lost.

**Hard:** Look up PostgreSQL's four isolation levels (READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE). Write a short description of what anomalies each level prevents.
