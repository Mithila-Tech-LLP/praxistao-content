# Chapter 10: Indexes — Making Queries Fast

Imagine a book with 1,000 pages and no index. To find every mention of "transactions", you'd read every single page. With an index at the back, you jump straight to the right pages. Database indexes do exactly the same thing — they let you find rows without reading every row.

## Table of Contents

1. The Problem: Full Table Scans
2. What an Index Is (and How It Works)
3. B-Tree Indexes — The Default
4. Creating and Using Indexes in SQL
5. Types of Indexes
6. Reading Query Plans with EXPLAIN
7. The Index Trade-off
8. Exercises

---

## 1. The Problem: Full Table Scans

Suppose you have a `users` table with 10 million rows:

```sql
SELECT * FROM users WHERE email = 'alice@example.com';
```

Without an index, PostgreSQL reads every row, checks if `email` matches, and returns matches. For 10 million rows, that's 10 million comparisons — potentially seconds of work.

This is called a **full table scan** (or sequential scan). It's fine for small tables but devastating for large ones.

---

## 2. What an Index Is (and How It Works)

An index is a separate data structure that maps index key values to row locations (page number + row offset). It's maintained by the database automatically as you insert, update, and delete rows.

Think of it like a library card catalog (or today, a library search terminal). Instead of walking every aisle looking for a book about "databases", you look it up in the catalog: "databases → Aisle 7, Shelf 3". Direct navigation.

```
Without index:                     With index on email:
Full scan of users table           Index lookup:
Row 1: email=bob@...      ✗        email='alice@...' → page 847, row 3
Row 2: email=carol@...    ✗        Read only page 847
Row 3: email=alice@...    ✓        Return row
...10 million rows...
```

---

## 3. B-Tree Indexes — The Default

The most common index type is the **B-Tree** (Balanced Tree). Every `CREATE INDEX` in PostgreSQL and MySQL creates a B-Tree by default.

A B-Tree is a sorted tree structure. Each internal node holds sorted keys and pointers to child nodes. Leaf nodes hold the actual index entries (key + row pointer).

```
                    [50]
                   /    \
               [25]      [75]
              /   \      /   \
         [10,20] [30,40] [60,70] [80,90]
          |  |    |  |    |  |    |  |
         rows   rows   rows   rows
```

To find email = 'alice@...':
1. Start at the root, compare with sorted keys, go left or right
2. Repeat down each level
3. At a leaf node, follow the pointer to the actual row

A B-Tree with 1 billion rows has a height of about 30. So any lookup takes at most 30 comparisons — regardless of table size. That's O(log n) instead of O(n).

---

## 4. Creating and Using Indexes in SQL

```sql
-- Create an index on the email column
CREATE INDEX idx_users_email ON users(email);

-- Now this query uses the index automatically:
SELECT * FROM users WHERE email = 'alice@example.com';

-- Create a unique index (also enforces uniqueness)
CREATE UNIQUE INDEX idx_users_email_unique ON users(email);

-- Create a composite index (on multiple columns)
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at);

-- Drop an index
DROP INDEX idx_users_email;
```

PostgreSQL automatically uses an index when it determines it will be faster than a full scan. For small tables, it may choose a full scan even with an index (reading 10 rows doesn't need an index).

### Composite Indexes

A composite index on `(user_id, created_at)` can answer:

```sql
-- Uses the index (leftmost prefix)
SELECT * FROM orders WHERE user_id = 5;

-- Uses the index (full composite)
SELECT * FROM orders WHERE user_id = 5 AND created_at > '2024-01-01';

-- Does NOT use the composite index (not leading with user_id)
SELECT * FROM orders WHERE created_at > '2024-01-01';
```

The rule: a composite index can be used if your WHERE clause starts with the leftmost columns of the index.

---

## 5. Types of Indexes

### Hash Indexes

Hash indexes map a key directly to its location using a hash function. Extremely fast for equality lookups (`=`), but useless for range queries (`<`, `>`).

```sql
CREATE INDEX idx_users_email_hash ON users USING hash(email);

-- Fast:
SELECT * FROM users WHERE email = 'alice@example.com'; ✓

-- Cannot use hash index:
SELECT * FROM users WHERE email > 'a'; ✗
```

### Partial Indexes

Index only rows that match a condition. Smaller and faster:

```sql
-- Only index active users
CREATE INDEX idx_active_users ON users(email) WHERE active = true;

-- This query uses the partial index:
SELECT * FROM users WHERE email = 'alice@example.com' AND active = true;
```

Great for filtering out NULL values, soft-deleted rows, or inactive records.

### Covering Indexes (INCLUDE)

Include extra columns in the index so the database doesn't need to fetch the actual row:

```sql
-- Index email, but also include name and id for "index-only scans"
CREATE INDEX idx_users_email_covering ON users(email) INCLUDE (id, name);

-- This query can be answered entirely from the index:
SELECT id, name FROM users WHERE email = 'alice@example.com';
-- No need to go to the actual table!
```

### GIN and GiST Indexes

For special data types:

```sql
-- GIN index for full-text search
CREATE INDEX idx_articles_fts ON articles USING gin(to_tsvector('english', content));

-- GIN index for JSONB queries
CREATE INDEX idx_users_metadata ON users USING gin(metadata);

-- GiST index for geographic data
CREATE INDEX idx_locations_coords ON locations USING gist(coordinates);
```

---

## 6. Reading Query Plans with EXPLAIN

`EXPLAIN` shows you how PostgreSQL plans to execute a query. `EXPLAIN ANALYZE` actually runs the query and shows real timings.

```sql
-- Show the plan without running
EXPLAIN SELECT * FROM users WHERE email = 'alice@example.com';

-- Show the plan AND run it (get real timing)
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'alice@example.com';
```

**Without an index:**
```
Seq Scan on users  (cost=0.00..18334.00 rows=1 width=64)
  Filter: (email = 'alice@example.com')
Planning time: 0.1 ms
Execution time: 1234.5 ms  ← 1.2 seconds!
```

**With an index:**
```
Index Scan using idx_users_email on users  (cost=0.43..8.45 rows=1 width=64)
  Index Cond: (email = 'alice@example.com')
Planning time: 0.2 ms
Execution time: 0.1 ms  ← 0.1 milliseconds!
```

Key things to look for:
- `Seq Scan` = full table scan (often bad for large tables)
- `Index Scan` = using an index (good)
- `Index Only Scan` = answer from index alone, no table fetch (great)
- `cost=X..Y` = estimated cost (Y is the total; lower is better)
- `rows=N` = estimated number of rows returned

In Go, you can run EXPLAIN and print the result:

```go
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func explainQuery(db *sql.DB, query string, args ...interface{}) {
    rows, err := db.Query("EXPLAIN ANALYZE "+query, args...)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var line string
        rows.Scan(&line)
        fmt.Println(line)
    }
}

func main() {
    db, _ := sql.Open("postgres", "postgres://localhost/mydb?sslmode=disable")
    defer db.Close()

    explainQuery(db, "SELECT * FROM users WHERE email = $1", "alice@example.com")
}
```

---

## 7. The Index Trade-off

Indexes are not free. Every index:

- **Speeds up reads** (SELECT, WHERE, JOIN, ORDER BY)
- **Slows down writes** (INSERT, UPDATE, DELETE must update every index)
- **Uses disk space** (an index can be as large as the table itself)

Rules of thumb:
- Always index foreign key columns (speeds up JOINs)
- Always index columns used in frequent WHERE clauses
- Don't index columns with very few distinct values (e.g., a `status` column with only 'active'/'inactive' — on a 50/50 split, a full scan is often faster)
- Don't add indexes you don't need — they slow down every write

```sql
-- Good candidates for indexes:
CREATE INDEX ON users(email);          -- frequent lookup column
CREATE INDEX ON orders(user_id);       -- foreign key
CREATE INDEX ON sessions(expires_at); -- used in WHERE expires_at < NOW()

-- Poor candidates:
CREATE INDEX ON users(gender);         -- only 2-3 distinct values
CREATE INDEX ON logs(message);         -- long text, rarely searched exactly
```

---

## Summary

- Without indexes, queries do a **full table scan** — reading every row. O(n).
- **B-Tree indexes** keep keys sorted, enabling O(log n) lookups.
- **Hash indexes** are fastest for equality, useless for ranges.
- **Partial indexes** only index rows matching a condition (smaller and faster).
- **Covering indexes** include extra columns to allow index-only scans.
- `EXPLAIN ANALYZE` shows how PostgreSQL executes a query and where time is spent.
- Indexes speed up reads but slow down writes — add them where needed, not everywhere.

### Quick Check

1. What type of scan does PostgreSQL use when there is no index?
2. Can a composite index on `(user_id, created_at)` be used for `WHERE created_at > '2024-01-01'`?
3. What is the main trade-off of adding many indexes?

### Exercises

**Easy:** Create a table with 100,000 rows (you can use `generate_series`). Run a SELECT with EXPLAIN ANALYZE before and after adding an index. Compare the execution times.

**Medium:** Create a partial index that only indexes users where `deleted_at IS NULL`. Write a query that will use this index and verify with EXPLAIN.

**Hard:** Read a PostgreSQL EXPLAIN ANALYZE output for a multi-table JOIN query. Identify which join algorithm (Hash Join, Nested Loop, Merge Join) was used and whether indexes are being used on both sides of the join.
