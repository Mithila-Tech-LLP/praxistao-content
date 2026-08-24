# Chapter 11: ACID Transactions — Never Lose Your Data

You just learned what ACID means in theory. Now let's make it real with code. This chapter shows you exactly how to use transactions in Go with PostgreSQL, handle the tricky edge cases, and build systems that never lose data.

## Table of Contents

1. Transactions in Go — The Complete Pattern
2. The Bank Transfer in Full Go Code
3. Isolation Levels Deep Dive
4. SELECT FOR UPDATE — Locking Specific Rows
5. Savepoints — Partial Rollbacks
6. Deadlock Handling in Go
7. Mini Project: Inventory Management with Atomic Updates
8. Exercises

---

## 1. Transactions in Go — The Complete Pattern

The `database/sql` package provides three methods for transactions:

```go
db.Begin()         // start a transaction → returns *sql.Tx
tx.Exec(...)       // run a statement within the transaction
tx.QueryRow(...)   // run a query within the transaction
tx.Commit()        // apply all changes
tx.Rollback()      // undo all changes
```

The canonical transaction pattern in Go:

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func withTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    // If fn panics or returns an error, roll back
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p) // re-panic after rollback
        }
    }()

    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}
```

This helper wraps any function in a transaction. If the function returns an error, the transaction rolls back. If it succeeds, it commits. This pattern keeps your business logic clean.

### Using the helper:

```go
err := withTransaction(db, func(tx *sql.Tx) error {
    // Do work inside the transaction
    _, err := tx.Exec("UPDATE products SET stock = stock - 1 WHERE id = $1", productID)
    if err != nil {
        return err // triggers rollback
    }

    _, err = tx.Exec("INSERT INTO orders (product_id, user_id) VALUES ($1, $2)", productID, userID)
    if err != nil {
        return err // triggers rollback
    }

    return nil // triggers commit
})

if err != nil {
    log.Printf("transaction failed: %v", err)
}
```

---

## 2. The Bank Transfer in Full Go Code

Let's implement the classic bank transfer completely:

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type Account struct {
    ID      int
    Balance float64
}

func transfer(db *sql.DB, fromID, toID int, amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }

    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin: %w", err)
    }
    defer tx.Rollback() // no-op if Commit succeeds

    // Lock both accounts in a consistent order to prevent deadlock.
    // Always lock the lower ID first.
    first, second := fromID, toID
    if toID < fromID {
        first, second = toID, fromID
    }

    // SELECT FOR UPDATE locks the rows so no other transaction
    // can modify them until we COMMIT or ROLLBACK.
    var fromBalance, toBalance float64

    row := tx.QueryRow(
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", first)
    if err := row.Scan(&fromBalance); err != nil {
        return fmt.Errorf("lock account %d: %w", first, err)
    }

    row = tx.QueryRow(
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", second)
    if err := row.Scan(&toBalance); err != nil {
        return fmt.Errorf("lock account %d: %w", second, err)
    }

    // Re-assign to correct variables based on original fromID/toID
    if fromID == first {
        // fromBalance is already correct
    } else {
        fromBalance, toBalance = toBalance, fromBalance
    }

    // Check sufficient funds
    if fromBalance < amount {
        return ErrInsufficientFunds
    }

    // Debit
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID)
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }

    // Credit
    _, err = tx.Exec(
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID)
    if err != nil {
        return fmt.Errorf("credit: %w", err)
    }

    // Record the transaction history
    _, err = tx.Exec(
        "INSERT INTO transfer_log (from_id, to_id, amount) VALUES ($1, $2, $3)",
        fromID, toID, amount)
    if err != nil {
        return fmt.Errorf("log: %w", err)
    }

    return tx.Commit()
}

func main() {
    db, err := sql.Open("postgres", "postgres://localhost/bank?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = transfer(db, 1, 2, 100.00)
    if errors.Is(err, ErrInsufficientFunds) {
        fmt.Println("Transfer failed: not enough funds")
    } else if err != nil {
        fmt.Println("Transfer failed:", err)
    } else {
        fmt.Println("Transfer successful!")
    }
}
```

Notice the important detail: we lock accounts in **consistent order** (lower ID first). This prevents deadlocks when two concurrent transfers lock the same accounts in opposite orders.

---

## 3. Isolation Levels Deep Dive

PostgreSQL supports four isolation levels. Each prevents different types of anomalies.

### READ COMMITTED (Default in PostgreSQL)

Each statement in a transaction sees the latest committed data at the time that statement runs. Two reads in the same transaction can see different data if another transaction commits between them.

```sql
-- Session A:
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- Returns 500

-- Session B (commits between A's two reads):
UPDATE accounts SET balance = 600 WHERE id = 1;
COMMIT;

-- Session A (same transaction, different read):
SELECT balance FROM accounts WHERE id = 1;  -- Returns 600 (changed!)
COMMIT;
```

This is called a **non-repeatable read** — reading the same row twice gives different results.

### REPEATABLE READ

The transaction sees a consistent snapshot of the database from when it started. Reads are stable within the transaction.

```sql
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
BEGIN;
SELECT balance FROM accounts WHERE id = 1;  -- Returns 500
-- (even if Session B updates and commits here)
SELECT balance FROM accounts WHERE id = 1;  -- Still returns 500
COMMIT;
```

In Go:

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})
```

### SERIALIZABLE

The strongest isolation level. Transactions are executed as if they ran one after another (serially), even though they actually run concurrently. PostgreSQL detects any read-write conflicts and aborts transactions that could produce non-serializable results.

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
```

Use SERIALIZABLE when your logic absolutely requires no concurrent interference (financial calculations, etc.). Be prepared to retry serialization failures.

| Isolation Level   | Dirty Read | Non-Repeatable Read | Phantom Read |
|------------------|------------|---------------------|--------------|
| READ UNCOMMITTED | Possible   | Possible            | Possible     |
| READ COMMITTED   | Prevented  | Possible            | Possible     |
| REPEATABLE READ  | Prevented  | Prevented           | Prevented*   |
| SERIALIZABLE     | Prevented  | Prevented           | Prevented    |

*PostgreSQL's REPEATABLE READ prevents phantoms too (stronger than SQL standard requires).

---

## 4. SELECT FOR UPDATE — Locking Specific Rows

`SELECT FOR UPDATE` acquires an exclusive lock on the selected rows. Other transactions trying to `SELECT FOR UPDATE` or `UPDATE` the same rows will wait until you commit or rollback.

```sql
-- Lock the row while checking inventory
BEGIN;
SELECT quantity FROM inventory WHERE product_id = 5 FOR UPDATE;
-- No one else can update this row until we finish
UPDATE inventory SET quantity = quantity - 1 WHERE product_id = 5;
COMMIT;
```

Variants:
- `FOR UPDATE` — exclusive lock, others wait
- `FOR SHARE` — shared lock, others can also acquire FOR SHARE but not FOR UPDATE
- `FOR UPDATE SKIP LOCKED` — skip rows already locked (great for job queues)
- `FOR UPDATE NOWAIT` — fail immediately if can't get the lock

```go
// Job queue pattern: grab the next unclaimed job without waiting
row := tx.QueryRow(`
    SELECT id, payload
    FROM jobs
    WHERE status = 'pending'
    ORDER BY created_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
`)
```

`SKIP LOCKED` is a powerful pattern for work queues — multiple workers each grab different jobs without interfering with each other.

---

## 5. Savepoints — Partial Rollbacks

Sometimes you want to roll back only part of a transaction, not all of it. **Savepoints** let you mark checkpoints within a transaction.

```go
func processOrders(db *sql.DB, orders []Order) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    successCount := 0
    for _, order := range orders {
        // Create a savepoint before processing each order
        _, err := tx.Exec("SAVEPOINT sp_order")
        if err != nil {
            return err
        }

        err = processOneOrder(tx, order)
        if err != nil {
            // Roll back to the savepoint (undo just this order)
            tx.Exec("ROLLBACK TO SAVEPOINT sp_order")
            log.Printf("skipping order %d: %v", order.ID, err)
            // Continue with the next order
        } else {
            // Release the savepoint (we don't need to roll back to it)
            tx.Exec("RELEASE SAVEPOINT sp_order")
            successCount++
        }
    }

    log.Printf("processed %d/%d orders", successCount, len(orders))
    return tx.Commit()
}
```

Savepoints are useful when processing a batch where you want to skip bad items instead of rolling back everything.

---

## 6. Deadlock Handling in Go

When a deadlock is detected, PostgreSQL kills one of the transactions and returns an error code `40P01`. Your Go code should detect this and retry.

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "time"
    "github.com/lib/pq"
)

func isDeadlock(err error) bool {
    var pqErr *pq.Error
    if errors.As(err, &pqErr) {
        return pqErr.Code == "40P01" // deadlock_detected
    }
    return false
}

func withRetry(db *sql.DB, maxRetries int, fn func(*sql.Tx) error) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := withTransaction(db, fn)
        if err == nil {
            return nil
        }
        if isDeadlock(err) {
            wait := time.Duration(attempt*attempt*10) * time.Millisecond
            fmt.Printf("deadlock on attempt %d, retrying in %v\n", attempt+1, wait)
            time.Sleep(wait)
            continue
        }
        return err // non-deadlock error, don't retry
    }
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

The key insight: deadlocks are expected and normal in concurrent systems. Your code must handle them gracefully by retrying.

---

## 7. Mini Project: Inventory Management with Atomic Updates

Build a simple inventory system that prevents overselling using transactions.

```go
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

var ErrOutOfStock = errors.New("product out of stock")

func setup(db *sql.DB) {
    db.Exec(`CREATE TABLE IF NOT EXISTS products (
        id      SERIAL PRIMARY KEY,
        name    TEXT NOT NULL,
        stock   INT NOT NULL CHECK (stock >= 0)
    )`)
    db.Exec(`CREATE TABLE IF NOT EXISTS orders (
        id         SERIAL PRIMARY KEY,
        product_id INT REFERENCES products(id),
        quantity   INT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
    )`)
    // Seed some products
    db.Exec(`INSERT INTO products (name, stock) VALUES ('Widget', 10) ON CONFLICT DO NOTHING`)
}

func purchaseProduct(db *sql.DB, productID, userID, qty int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Lock the product row and check stock
    var stock int
    err = tx.QueryRow(
        "SELECT stock FROM products WHERE id = $1 FOR UPDATE",
        productID,
    ).Scan(&stock)
    if err == sql.ErrNoRows {
        return fmt.Errorf("product %d not found", productID)
    }
    if err != nil {
        return err
    }

    if stock < qty {
        return ErrOutOfStock
    }

    // Deduct stock (the CHECK constraint prevents going below 0 as a safety net)
    _, err = tx.Exec(
        "UPDATE products SET stock = stock - $1 WHERE id = $2",
        qty, productID)
    if err != nil {
        return err
    }

    // Record the order
    _, err = tx.Exec(
        "INSERT INTO orders (product_id, quantity) VALUES ($1, $2)",
        productID, qty)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func main() {
    db, err := sql.Open("postgres", "postgres://localhost/inventory?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    setup(db)

    // Try to buy 3 widgets
    err = purchaseProduct(db, 1, 42, 3)
    if errors.Is(err, ErrOutOfStock) {
        fmt.Println("Sorry, out of stock!")
    } else if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Purchase successful!")
    }
}
```

The `CHECK (stock >= 0)` constraint is a second safety net — even if there's a bug in the application logic, the database will refuse to let stock go negative.

---

## Summary

- The canonical Go transaction pattern: `Begin()` → `defer Rollback()` → work → `Commit()`.
- `defer tx.Rollback()` is a no-op after a successful `Commit()` — safe to always include.
- **READ COMMITTED** (default): sees latest committed data per-statement. Can have non-repeatable reads.
- **REPEATABLE READ**: consistent snapshot for the whole transaction. Prevents non-repeatable reads.
- **SERIALIZABLE**: strongest isolation, prevents all anomalies, may require retries.
- `SELECT FOR UPDATE` locks rows; `SKIP LOCKED` enables efficient job queues.
- Deadlocks are normal — detect error code `40P01` and retry with backoff.

### Exercises

**Easy:** Write a Go function that runs three SQL statements in a single transaction. Force a rollback by returning an error from the middle statement and verify neither the first nor third statements' effects are visible.

**Medium:** Build a simple job queue using `SELECT FOR UPDATE SKIP LOCKED`. Write two goroutines that race to claim jobs from the queue without ever claiming the same job.

**Hard:** Implement a serializable transfer where you first read the sum of all account balances, then make a transfer, and assert the total sum is unchanged. Use `SERIALIZABLE` isolation and handle serialization failure retries.
