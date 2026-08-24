# Chapter 24: Index Deep Dive — B-Trees, Types, Strategies & Anti-Patterns

Indexes are the single most impactful performance lever in relational databases. Senior engineers must understand not just how to create them, but why they work, when they help, and when they hurt.

## Table of Contents

1. [How B-Tree Indexes Work](#1-how-b-tree-indexes-work)
2. [Index Types in PostgreSQL](#2-index-types-in-postgresql)
3. [Composite Indexes & Column Order](#3-composite-indexes--column-order)
4. [Covering Indexes](#4-covering-indexes)
5. [Partial Indexes](#5-partial-indexes)
6. [When Indexes DON'T Help](#6-when-indexes-dont-help)
7. [The N+1 Query Problem](#7-the-n1-problem)
8. [EXPLAIN ANALYZE](#8-explain-analyze)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. How B-Tree Indexes Work

A B-tree (Balanced Tree) index stores a sorted copy of the indexed column(s) plus pointers to the actual table rows.

```
Table: users(id, name, email, age)
Index on: email

B-tree index structure:
       [m@x.com]
       /         \
[a@x.com]       [z@x.com]
  /   \           /   \
[a@]  [j@]     [s@]   [z@]
 ↓     ↓         ↓      ↓
row   row       row    row
ptr   ptr       ptr    ptr

Finding "j@example.com":
1. Compare with root: j < m → go left
2. Compare with left: j > a → go right
3. Found! Follow pointer to heap row
Steps: O(log n) — for 1M rows: ~20 comparisons
```

**Without an index:** PostgreSQL reads every row in the table sequentially until it finds matches. For 1M rows with 10 matching: reads 1M rows. With index: reads ~20 nodes to find the start, follows pointers to 10 rows.

### Index Storage Cost

Every index is an additional data structure on disk. Writing to a table updates all its indexes. A table with 10 indexes has 10x the write overhead. Indexes trade write performance for read performance.

---

## 2. Index Types in PostgreSQL

```sql
-- B-TREE (default): for comparison operators (=, <, >, <=, >=, BETWEEN, LIKE prefix)
CREATE INDEX idx_email ON users(email);
-- Use for: exact match, range queries, sorting

-- HASH: for equality only (=), not ranges
CREATE INDEX idx_hash_email ON users USING HASH (email);
-- Faster for equality on large tables, but doesn't support range queries

-- GIN (Generalized Inverted Index): for array, JSONB, full-text search
CREATE INDEX idx_tags ON posts USING GIN (tags);    -- array
CREATE INDEX idx_data ON events USING GIN (data);   -- JSONB
CREATE INDEX idx_fts ON articles USING GIN (to_tsvector('english', content)); -- full-text

-- GIST: for geometric types, range types, full-text
CREATE INDEX idx_location ON venues USING GIST (location);  -- geographic point

-- BRIN (Block Range Index): for very large tables with naturally ordered data (e.g., timestamps)
-- Very small index, great for append-only tables
CREATE INDEX idx_created ON events USING BRIN (created_at);
```

---

## 3. Composite Indexes & Column Order

A composite index on (a, b, c) can serve queries filtering on:
- Just `a`
- `a` and `b`
- `a`, `b`, and `c`
- But NOT just `b` or just `c` (the leftmost prefix rule)

```sql
-- Index on (country, city, zip)
CREATE INDEX idx_location ON users(country, city, zip);

-- USES the index:
WHERE country = 'US'                                  -- prefix
WHERE country = 'US' AND city = 'NYC'                -- prefix
WHERE country = 'US' AND city = 'NYC' AND zip = '10001' -- full key

-- CANNOT use the index:
WHERE city = 'NYC'                                    -- no country
WHERE zip = '10001'                                   -- no country/city

-- Column order rule: put high-cardinality columns first IF used in equality,
-- put range/sort columns last
-- Query: WHERE country = 'US' AND age > 30 ORDER BY age
-- Best index: (country, age) — equality first, then range/sort
```

---

## 4. Covering Indexes

An index that contains all columns needed by a query allows an "index-only scan" — the database never needs to touch the heap (table rows).

```sql
-- Query: SELECT name, email FROM users WHERE age > 25 ORDER BY age
-- Without covering index: index scan + heap fetch for each row
-- With covering index:
CREATE INDEX idx_covering ON users(age) INCLUDE (name, email);
-- The INCLUDE columns are stored in the index leaf pages
-- Now the query can be served entirely from the index

-- Check if index-only scan happened:
EXPLAIN SELECT name, email FROM users WHERE age > 25;
-- Look for "Index Only Scan" vs "Index Scan"
```

---

## 5. Partial Indexes

A partial index only indexes a subset of rows where a condition is true. Smaller, faster to build, faster to scan.

```sql
-- Only index active users (most queries are for active users)
CREATE INDEX idx_active_users ON users(email) WHERE active = true;

-- Only index unprocessed orders
CREATE INDEX idx_pending_orders ON orders(created_at) WHERE status = 'pending';

-- Result: index is much smaller than a full-table index
-- Queries with the WHERE clause in the predicate use this index
-- Queries for all users (including inactive) do NOT use this index
```

---

## 6. When Indexes DON'T Help

**Case 1: Low cardinality columns**
```sql
-- Index on a boolean column (only 2 values): rarely useful
-- The planner might decide a full table scan is cheaper if 40% of rows are 'active'
CREATE INDEX idx_active ON users(active); -- often useless

-- Exception: partial index!
CREATE INDEX idx_inactive ON users(id) WHERE active = false;
-- If only 1% of users are inactive, this is very selective and useful
```

**Case 2: Functions or type casts on the indexed column**
```sql
-- DOES NOT use index on email:
WHERE UPPER(email) = 'ALICE@EXAMPLE.COM'
WHERE created_at::date = '2024-01-01'

-- FIX: create an expression index
CREATE INDEX idx_upper_email ON users (UPPER(email));
CREATE INDEX idx_date ON events ((created_at::date));

-- Then the query must match the expression exactly
```

**Case 3: LIKE with leading wildcard**
```sql
-- CANNOT use B-tree index:
WHERE name LIKE '%alice%'   -- leading wildcard
WHERE name LIKE '%alice'    -- leading wildcard

-- CAN use B-tree index:
WHERE name LIKE 'alice%'    -- no leading wildcard

-- For wildcard search: use GIN with pg_trgm extension
CREATE EXTENSION pg_trgm;
CREATE INDEX idx_name_trgm ON users USING GIN (name gin_trgm_ops);
-- Now LIKE '%alice%' can use this index
```

**Case 4: Table is too small**
For tables under ~1000 rows, PostgreSQL often prefers a sequential scan because the overhead of following index pointers exceeds the benefit.

---

## 7. The N+1 Problem

The N+1 problem is one of the most common performance bugs in production services. It happens when you execute N additional queries for each of N results.

```go
// N+1 in Go with a SQL driver:

// 1 query: get all users (N users returned)
rows, _ := db.QueryContext(ctx, "SELECT id, name FROM users")
var users []User
for rows.Next() {
    var u User
    rows.Scan(&u.ID, &u.Name)
    users = append(users, u)
}

// N queries: one per user to get their orders
for i := range users {
    // This runs N times — for 1000 users, 1001 total queries!
    orders, _ := db.QueryContext(ctx, "SELECT * FROM orders WHERE user_id = $1", users[i].ID)
    // ...
}

// FIX: JOIN or IN clause
users_with_orders, _ := db.QueryContext(ctx, `
    SELECT u.id, u.name, o.id as order_id, o.total
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
`)

// OR for ORM-style: batch load
userIDs := extractIDs(users)
orders, _ := db.QueryContext(ctx, `
    SELECT * FROM orders WHERE user_id = ANY($1)
`, pq.Array(userIDs))
```

---

## 8. EXPLAIN ANALYZE

```sql
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'alice@example.com';

-- Output:
-- Index Scan using idx_email on users  (cost=0.42..8.44 rows=1 width=48) (actual time=0.015..0.016 rows=1 loops=1)
--   Index Cond: (email = 'alice@example.com')
-- Planning Time: 0.065 ms
-- Execution Time: 0.032 ms

-- Key things to look for:
-- "Seq Scan" on large table = missing index or low-cardinality column
-- "Nested Loop" with large row counts = N+1 or missing index on join column  
-- High "actual rows" >> "rows estimate" = stale statistics (run ANALYZE)
-- "Bitmap Heap Scan" = moderate number of matching rows (between index and seq scan)
```

---

## 9. Interview Questions & Model Answers

**Q: How does a B-tree index work and what operations does it support?**

"A B-tree index keeps a sorted copy of the indexed column values with pointers to the heap rows. Looking up a value takes O(log n) comparisons by traversing the tree from root to leaf. Because values are sorted, B-trees support equality (=), range queries (<, >, BETWEEN), prefix LIKE queries, ORDER BY, and MIN/MAX without a full table scan. They're the default index type in PostgreSQL for a reason — they handle the vast majority of query patterns."

**Q: Explain the leftmost prefix rule for composite indexes.**

"A composite index on (a, b, c) can only be used if the query filters on the leftmost columns. Filtering on just `a` uses the index. Filtering on `a` and `b` uses the index. Filtering on just `b` or just `c` cannot use it — the index is sorted by `a` first, so without an `a` filter, you'd have to scan the whole index anyway. The practical implication: order composite index columns by how they'll be used in queries — equality columns first, then range or sort columns."

**Q: What is a covering index and when would you use one?**

"A covering index includes all the columns a query needs — so the database can answer the query entirely from the index without fetching heap rows. This eliminates a second I/O operation (the 'heap fetch') per matched row. Use it when a query runs very frequently and always accesses the same set of columns. In PostgreSQL, use INCLUDE to add extra columns to the index leaf pages without including them in the sorted key — keeps the index smaller while still covering the query."

---

## Summary

- **B-tree indexes:** sorted structure enabling O(log n) lookups. Support equality, range, sort, prefix LIKE.
- **Composite index:** leftmost prefix rule. Put equality columns before range/sort columns.
- **Covering index:** all needed columns in the index → index-only scan → no heap fetch.
- **Partial index:** only index rows matching a condition → smaller, faster for selective queries.
- **When indexes fail:** low cardinality, leading wildcard LIKE, functions on indexed column, small tables.
- **N+1 problem:** N queries for N results. Fix with JOINs, IN clauses, or batch loading.
- `EXPLAIN ANALYZE` is your weapon: look for Seq Scan on large tables and Nested Loop with high row counts.
