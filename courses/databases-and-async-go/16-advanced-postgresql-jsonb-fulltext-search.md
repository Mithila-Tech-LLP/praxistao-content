# Chapter 16: Advanced PostgreSQL — JSONB, Full-Text Search, and More

PostgreSQL is secretly multiple databases in one. This chapter unlocks the features that most developers don't know exist: document storage with JSONB, built-in full-text search, window functions, and table partitioning. These can replace entire separate systems.

## Table of Contents

1. JSONB — A Document Database Inside PostgreSQL
2. Full-Text Search — Built In, No Elasticsearch Needed
3. Arrays — Storing Lists in a Column
4. Window Functions Deep Dive
5. Table Partitioning — Handling Huge Tables
6. Mini Project: Full-Text Article Search in Go
7. Exercises

---

## 1. JSONB — A Document Database Inside PostgreSQL

`JSONB` stores JSON as a binary format. Unlike `JSON` (stored as-is, slow to query), JSONB is parsed on insert and stored efficiently — making queries fast.

### Basic JSONB Operations

```sql
CREATE TABLE products (
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'
);

INSERT INTO products (name, metadata) VALUES
    ('Laptop', '{"brand":"Apple","specs":{"ram":16,"storage":512},"tags":["electronics","computers"]}'),
    ('Phone',  '{"brand":"Samsung","specs":{"ram":8,"storage":256},"tags":["electronics","mobile"]}');

-- Access a top-level key (returns JSON)
SELECT metadata->'brand' FROM products;           -- "Apple", "Samsung"

-- Access a top-level key (returns text, no quotes)
SELECT metadata->>'brand' FROM products;          -- Apple, Samsung

-- Access nested key
SELECT metadata->'specs'->>'ram' FROM products;  -- 16, 8

-- Check if key exists
SELECT * FROM products WHERE metadata ? 'brand';

-- Check if object contains key-value pair
SELECT * FROM products WHERE metadata @> '{"brand":"Apple"}';

-- Check if array element exists
SELECT * FROM products WHERE metadata->'tags' ? 'electronics';
```

### Modifying JSONB

```sql
-- Add or update a key
UPDATE products
SET metadata = metadata || '{"in_stock": true}'
WHERE id = 1;

-- Remove a key
UPDATE products
SET metadata = metadata - 'in_stock'
WHERE id = 1;

-- Update a nested value
UPDATE products
SET metadata = jsonb_set(metadata, '{specs,ram}', '32')
WHERE id = 1;
```

### Indexing JSONB

```sql
-- GIN index for any JSONB query (@>, ?, ?|, ?&)
CREATE INDEX idx_products_metadata ON products USING gin(metadata);

-- Expression index for a specific key (faster for single-key queries)
CREATE INDEX idx_products_brand ON products ((metadata->>'brand'));
```

With a GIN index, `WHERE metadata @> '{"brand":"Apple"}'` is very fast even with millions of rows.

### JSONB in Go

```go
import (
    "encoding/json"
    "github.com/jackc/pgx/v5"
)

type ProductMetadata struct {
    Brand  string   `json:"brand"`
    Tags   []string `json:"tags"`
    InStock bool    `json:"in_stock"`
}

type Product struct {
    ID       int64
    Name     string
    Metadata ProductMetadata
}

func GetProduct(ctx context.Context, id int64) (*Product, error) {
    var p Product
    var metaJSON []byte

    err := pool.QueryRow(ctx,
        "SELECT id, name, metadata FROM products WHERE id = $1", id,
    ).Scan(&p.ID, &p.Name, &metaJSON)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(metaJSON, &p.Metadata); err != nil {
        return nil, err
    }
    return &p, nil
}

func UpdateMetadata(ctx context.Context, id int64, meta ProductMetadata) error {
    metaJSON, err := json.Marshal(meta)
    if err != nil {
        return err
    }
    _, err = pool.Exec(ctx,
        "UPDATE products SET metadata = $1 WHERE id = $2",
        metaJSON, id,
    )
    return err
}
```

---

## 2. Full-Text Search — Built In, No Elasticsearch Needed

For most applications, PostgreSQL's built-in full-text search is sufficient — no need for a separate Elasticsearch cluster.

### How It Works

PostgreSQL converts text into a `tsvector` (text search vector) — a sorted list of **lexemes** (normalized word forms) with position information. Queries use `tsquery` — boolean expressions of lexemes.

```sql
-- Convert text to tsvector
SELECT to_tsvector('english', 'The quick brown fox jumps over the lazy dog');
-- 'brown':3 'dog':9 'fox':4 'jump':5 'lazi':8 'quick':2
-- Stop words ("the", "over") are removed; words are normalized

-- Create a tsquery
SELECT to_tsquery('english', 'fox & jump');
-- 'fox' & 'jump'

-- Check if document matches query
SELECT to_tsvector('english', 'The fox jumps') @@ to_tsquery('english', 'fox & jump');
-- true
```

### Setting Up Full-Text Search

```sql
CREATE TABLE articles (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    search_vec TSVECTOR GENERATED ALWAYS AS (
        setweight(to_tsvector('english', title), 'A') ||
        setweight(to_tsvector('english', content), 'B')
    ) STORED,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- GIN index on the precomputed vector
CREATE INDEX idx_articles_search ON articles USING gin(search_vec);
```

`GENERATED ALWAYS AS ... STORED` means PostgreSQL automatically updates `search_vec` on every INSERT and UPDATE. `setweight` gives title matches a higher relevance than content matches.

### Searching

```sql
-- Find articles matching a query, ranked by relevance
SELECT
    id,
    title,
    ts_rank(search_vec, query) AS rank,
    ts_headline('english', content, query, 'MaxWords=20') AS snippet
FROM articles,
     to_tsquery('english', 'database & performance') AS query
WHERE search_vec @@ query
ORDER BY rank DESC
LIMIT 10;
```

`ts_headline` generates highlighted snippets — the matching words surrounded by context.

### Full-Text Search in Go

```go
func SearchArticles(ctx context.Context, q string) ([]Article, error) {
    // Convert user query to tsquery safely
    // plainto_tsquery handles phrases, no special syntax needed
    rows, err := pool.Query(ctx, `
        SELECT
            id, title,
            ts_rank(search_vec, plainto_tsquery('english', $1)) AS rank,
            ts_headline('english', content, plainto_tsquery('english', $1),
                'MaxWords=30, MinWords=10') AS snippet
        FROM articles
        WHERE search_vec @@ plainto_tsquery('english', $1)
        ORDER BY rank DESC
        LIMIT 20
    `, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var articles []Article
    for rows.Next() {
        var a Article
        if err := rows.Scan(&a.ID, &a.Title, &a.Rank, &a.Snippet); err != nil {
            return nil, err
        }
        articles = append(articles, a)
    }
    return articles, rows.Err()
}
```

### Trigram Search for Fuzzy Matching

For "did you mean?" style fuzzy search, use the `pg_trgm` extension:

```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_products_name_trgm ON products USING gin(name gin_trgm_ops);

-- Find similar names (handles typos)
SELECT name, similarity(name, 'labtop') AS sim
FROM products
WHERE name % 'labtop'  -- % operator uses trigram similarity
ORDER BY sim DESC
LIMIT 5;
-- Returns: "laptop", "laptops", "laptop bag", ...
```

---

## 3. Arrays — Storing Lists in a Column

```sql
CREATE TABLE posts (
    id   BIGSERIAL PRIMARY KEY,
    tags TEXT[]
);

INSERT INTO posts (tags) VALUES (ARRAY['go', 'databases', 'tutorial']);

-- Find posts containing a specific tag
SELECT * FROM posts WHERE 'go' = ANY(tags);

-- Find posts containing all of these tags
SELECT * FROM posts WHERE tags @> ARRAY['go', 'databases'];

-- Find posts containing any of these tags
SELECT * FROM posts WHERE tags && ARRAY['go', 'python'];

-- Expand array to rows
SELECT id, UNNEST(tags) AS tag FROM posts;

-- Aggregate: collect all tags with counts
SELECT tag, COUNT(*) AS usage
FROM posts, UNNEST(tags) AS tag
GROUP BY tag
ORDER BY usage DESC;
```

Arrays in Go:

```go
// pgx handles arrays natively
var tags []string
err := pool.QueryRow(ctx, "SELECT tags FROM posts WHERE id = $1", id).Scan(&tags)

// Insert an array
_, err = pool.Exec(ctx,
    "UPDATE posts SET tags = $1 WHERE id = $2",
    []string{"go", "database"}, id,
)
```

---

## 4. Window Functions Deep Dive

Window functions compute values across a "window" of related rows without collapsing them into a single row (unlike GROUP BY).

```sql
-- Rank posts by likes within each user
SELECT
    user_id,
    title,
    likes,
    RANK() OVER (PARTITION BY user_id ORDER BY likes DESC) AS rank_in_user,
    ROW_NUMBER() OVER (ORDER BY likes DESC) AS global_rank
FROM posts;

-- Running total
SELECT
    date,
    revenue,
    SUM(revenue) OVER (ORDER BY date) AS running_total
FROM daily_sales;

-- Compare to previous row (LAG) or next row (LEAD)
SELECT
    date,
    revenue,
    LAG(revenue) OVER (ORDER BY date) AS prev_day,
    revenue - LAG(revenue) OVER (ORDER BY date) AS change
FROM daily_sales;

-- Moving average (7-day)
SELECT
    date,
    revenue,
    AVG(revenue) OVER (ORDER BY date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS ma7
FROM daily_sales;
```

In Go, window function results are scanned like any other column:

```go
rows, err := pool.Query(ctx, `
    SELECT user_id, post_id, likes,
           RANK() OVER (PARTITION BY user_id ORDER BY likes DESC) AS rank
    FROM posts
`)
for rows.Next() {
    var userID, postID, likes, rank int
    rows.Scan(&userID, &postID, &likes, &rank)
    fmt.Printf("User %d: post %d has rank %d with %d likes\n", userID, postID, rank, likes)
}
```

---

## 5. Table Partitioning — Handling Huge Tables

When a table grows to billions of rows, queries slow down even with indexes. Partitioning splits the table into smaller physical tables (partitions) while presenting a single logical table.

```sql
-- Partition by month (range partitioning)
CREATE TABLE events (
    id         BIGSERIAL,
    user_id    BIGINT NOT NULL,
    event_type TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL
) PARTITION BY RANGE (occurred_at);

-- Create monthly partitions
CREATE TABLE events_2024_01 PARTITION OF events
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE events_2024_02 PARTITION OF events
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- PostgreSQL automatically routes inserts to the right partition
INSERT INTO events (user_id, event_type, occurred_at)
VALUES (1, 'click', '2024-01-15 10:00:00+00');
-- Goes into events_2024_01

-- Queries that filter on occurred_at only scan relevant partitions
EXPLAIN SELECT * FROM events WHERE occurred_at >= '2024-01-01' AND occurred_at < '2024-02-01';
-- "Seq Scan on events_2024_01" → skips all other partitions!
```

**Drop old data instantly** by dropping a partition (O(1), no row-by-row delete):

```sql
DROP TABLE events_2024_01; -- drops all January 2024 data instantly
```

---

## 6. Mini Project: Full-Text Article Search in Go

Build a Go HTTP server with full-text search:

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    "github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type Article struct {
    ID      int64   `json:"id"`
    Title   string  `json:"title"`
    Rank    float64 `json:"rank"`
    Snippet string  `json:"snippet"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    if q == "" {
        http.Error(w, "q parameter required", 400)
        return
    }

    rows, err := db.Query(r.Context(), `
        SELECT id, title,
               ts_rank(search_vec, plainto_tsquery('english', $1)) AS rank,
               ts_headline('english', content, plainto_tsquery('english', $1),
                   'MaxWords=25, StartSel=<b>, StopSel=</b>') AS snippet
        FROM articles
        WHERE search_vec @@ plainto_tsquery('english', $1)
        ORDER BY rank DESC
        LIMIT 10
    `, q)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()

    var results []Article
    for rows.Next() {
        var a Article
        if err := rows.Scan(&a.ID, &a.Title, &a.Rank, &a.Snippet); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        results = append(results, a)
    }

    if results == nil {
        results = []Article{}
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func main() {
    var err error
    db, err = pgxpool.New(context.Background(), "postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("GET /search", searchHandler)
    log.Println("Search API on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Test:
```bash
curl "http://localhost:8080/search?q=database+performance"
```

---

## Summary

- **JSONB** gives you document-database flexibility inside PostgreSQL with full indexing. Use `@>`, `?`, `->>` operators.
- **Full-text search** with `tsvector`/`tsquery` handles most search needs without Elasticsearch. Use `GENERATED ALWAYS AS STORED` columns for auto-updating search vectors.
- **pg_trgm** provides fuzzy matching for autocomplete and typo-tolerance.
- **Arrays** let you store lists of values with efficient `ANY` and `@>` queries.
- **Window functions** (RANK, LAG, SUM OVER) compute running totals, rankings, and comparisons without subqueries.
- **Partitioning** splits huge tables by range, list, or hash. Queries only scan relevant partitions.

### Exercises

**Easy:** Create an `articles` table with a GENERATED tsvector column and GIN index. Insert 10 articles, then search for keywords and verify the GIN index is used.

**Medium:** Build a product search endpoint that combines full-text search (for relevance) with exact tag filtering (using `tags && ARRAY[...]`). Return results ranked by a combination of recency and relevance.

**Hard:** Implement table partitioning on a high-volume `logs` table partitioned by month. Write a Go function that automatically creates next month's partition if it doesn't exist (idempotent).
