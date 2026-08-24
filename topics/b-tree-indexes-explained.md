---
title: B-Tree Indexes and Why They Make Queries Fast
category: Software & Programming
tags: [Databases, Performance]
duration: 8 min read
relatedCourses: [databases-and-async-go, senior-engineer-interview]
relatedProjects: [key-value-store]
relatedTopics: [connection-pooling-explained, acid-vs-base-consistency-models]
---

## TL;DR

- Without an index, finding a row means scanning the whole table — O(n). A B-tree index turns that into O(log n) by keeping keys sorted in a shallow, wide tree structure.
- "B" doesn't stand for "binary" — B-trees are wide (hundreds of keys per node, not two), specifically so the tree stays shallow even with millions of rows, minimizing the number of disk reads needed.
- An index speeds up reads that use it, but every index also slows down writes (each insert/update/delete has to update the index too) and costs storage — indexes are a tradeoff, not a free win.
- A composite index's column *order* determines which queries it can actually help — this is the single most common indexing mistake.

## Why a Full Table Scan Is Slow

Without an index, `SELECT * FROM users WHERE email = 'x@example.com'` means the database reads every row, checks if its `email` matches, and moves on. For a table with 10 million rows, that's 10 million comparisons for a query that should return at most one row. This is O(n) — the cost grows linearly with table size, no matter how selective the query actually is.

## The B-Tree Structure

A B-tree keeps keys sorted, and organizes them into a shallow, wide tree where each node holds many keys (not just one or two) and many child pointers:

```
                    [ 50 | 150 ]
                   /      |      \
          [10|30]     [70|110]    [200|300]
         /   |   \    /   |   \   /   |   \
       ...  ...  ... ...  ... ... ...  ... ...
```

To find a key, start at the root and compare against the keys in that node to decide which child pointer to follow — repeat at each level until reaching a leaf. Because each node holds hundreds of keys (the actual number, called the tree's "fanout," is tuned to match the database's disk page size), even a table with tens of millions of rows only needs a tree 3-4 levels deep. That's the entire point: 3-4 comparisons (or disk reads, in practice) to find any row, versus millions in a full scan.

This is why it's a "B"-tree, not a binary tree: a binary tree (2 children per node) would need `log2(10,000,000)` ≈ 24 levels to hold the same data — a B-tree with a fanout of, say, 200 needs only `log200(10,000,000)` ≈ 3 levels. Fewer levels means fewer disk reads to traverse from root to leaf, and disk reads (even on an SSD) are the dominant cost in a database query — this is the whole reason the tree is optimized to be wide and shallow rather than the classic balanced-binary-tree shape you'd use for an in-memory structure.

## What "Using the Index" Actually Means

```sql
CREATE INDEX idx_users_email ON users(email);

SELECT * FROM users WHERE email = 'x@example.com';
-- query planner sees a usable index on `email`, walks the B-tree
-- (root -> ... -> leaf), finds the matching entry, which points
-- to the exact row location -> one targeted row fetch instead of a full scan
```

The index itself doesn't store the whole row — it stores the indexed column's value plus a pointer to where the actual row lives (in PostgreSQL, a "TID"; in MySQL/InnoDB, the primary key, since InnoDB's tables themselves are structured as a B-tree keyed by primary key). The query still does one more lookup to fetch the full row after finding it in the index, unless the index alone contains every column the query needs (a "covering index," a further optimization).

## The Write Cost — Why You Don't Index Everything

Every `INSERT`, `UPDATE` (of an indexed column), or `DELETE` has to update every index on that table, not just the underlying rows. A table with five indexes means every write does up to six pieces of work (the row itself, plus five index updates) instead of one. This is the direct tradeoff: indexes make specific reads fast at the cost of making all writes on that table slower, plus consuming additional disk space per index. This is exactly why "just add an index" isn't a free optimization — it's a deliberate cost/benefit call for each specific query pattern.

## Composite Indexes and Column Order

An index on multiple columns, `CREATE INDEX idx ON orders(customer_id, created_at)`, is sorted first by `customer_id`, then by `created_at` *within* each `customer_id`. This ordering has a direct, non-obvious consequence: the index can efficiently serve queries that filter on `customer_id` alone, or on `customer_id` AND `created_at` together — but it **cannot** efficiently serve a query that filters on `created_at` alone, because the index isn't sorted by `created_at` globally, only within each `customer_id` group.

```sql
-- uses the index efficiently:
WHERE customer_id = 42
WHERE customer_id = 42 AND created_at > '2026-01-01'

-- does NOT use this index efficiently (falls back to a full scan):
WHERE created_at > '2026-01-01'
```

This is the single most common real-world indexing bug: adding a composite index in the "wrong" column order for the queries that actually run, or adding one index per column separately when what the query patterns actually need is one composite index in the right order.

## Common Pitfalls

- **Indexing every column "just in case"** — each additional index adds write overhead and storage cost for every insert/update, whether or not any query actually benefits from it. Index based on real, known query patterns, not speculatively.
- **Wrong composite index column order** — as above; the leading column(s) of a composite index must match what your actual queries filter on most selectively and most often.
- **Not checking whether the query planner actually uses the index** — an index existing doesn't guarantee the database chooses to use it (e.g., for a query returning most of the table, a full scan can genuinely be faster). Use `EXPLAIN`/`EXPLAIN ANALYZE` to confirm, rather than assuming.
- **Low-selectivity indexes** — an index on a boolean column (`is_active`) with a 90/10 split of values often isn't very useful, since the database still has to fetch most of the table either way; indexes pay off most on columns with high cardinality (many distinct values) used in selective filters.
- **Forgetting indexes matter for `JOIN`s and `ORDER BY`, not just `WHERE`** — a join condition or sort column without a matching index can force a full scan or an expensive in-memory sort, even if the query has no `WHERE` clause at all.
