# Chapter 79: Advanced PostgreSQL — JSONB, Full-Text Search, and Indexes

PostgreSQL is not just a relational database. Its JSONB column type handles semi-structured data. Its built-in full-text search rivals dedicated search engines for moderate scale. Its index types — B-tree, GIN, GiST, BRIN — each solve a different problem. This chapter covers all three.

## Table of Contents

1. [JSONB — Semi-Structured Data](#1-jsonb--semi-structured-data)
2. [Full-Text Search](#2-full-text-search)
3. [Index Types](#3-index-types)
4. [Advanced Queries in Go](#4-advanced-queries-in-go)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. JSONB — Semi-Structured Data

JSONB stores JSON as a binary parse tree, not as raw text. This means:
- Keys are sorted and deduplicated
- Operators and indexes work on individual paths
- Writes are slightly slower (parsing), reads are faster (no re-parse)

```sql
-- Schema: a products table with flexible attributes per category
CREATE TABLE products (
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    price      NUMERIC(10,2) NOT NULL,
    category   TEXT NOT NULL,
    attributes JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Insert a product with category-specific attributes
INSERT INTO products (name, price, category, attributes) VALUES
    ('Sony WH-1000XM5', 349.99, 'electronics', '{
        "brand": "Sony",
        "battery_hours": 30,
        "noise_cancelling": true,
        "connectivity": ["bluetooth", "aux"]
    }');

INSERT INTO products (name, price, category, attributes) VALUES
    ('Levi''s 501', 59.99, 'clothing', '{
        "brand": "Levi''s",
        "sizes": ["28x30", "30x32", "32x34"],
        "material": "100% cotton",
        "fit": "original"
    }');
```

### JSONB Operators

```sql
-- → returns JSONB value (keeps as JSON)
SELECT attributes->'brand'            FROM products;  -- "Sony"

-- ->> returns TEXT value (extracts as string)
SELECT attributes->>'brand'           FROM products;  -- Sony

-- #> navigates a path (returns JSONB)
SELECT attributes#>'{connectivity,0}' FROM products;  -- "bluetooth"

-- #>> navigates a path (returns TEXT)
SELECT attributes#>>'{connectivity,0}' FROM products; -- bluetooth

-- @> contains (checks if left contains right)
SELECT * FROM products
WHERE attributes @> '{"noise_cancelling": true}';

-- ? key exists
SELECT * FROM products
WHERE attributes ? 'battery_hours';

-- ?| any key exists
SELECT * FROM products
WHERE attributes ?| ARRAY['battery_hours', 'material'];

-- ?& all keys exist
SELECT * FROM products
WHERE attributes ?& ARRAY['brand', 'fit'];
```

### Modifying JSONB

```sql
-- jsonb_set: update a path
UPDATE products
SET attributes = jsonb_set(attributes, '{battery_hours}', '35')
WHERE id = 1;

-- || merge (right wins on conflict)
UPDATE products
SET attributes = attributes || '{"wireless": true}'
WHERE category = 'electronics';

-- #- delete a key at path
UPDATE products
SET attributes = attributes #- '{connectivity,1}'
WHERE id = 1;
```

### JSONB Indexes

```sql
-- GIN index on all JSONB paths — fast @>, ?, ?|, ?& operators
CREATE INDEX idx_products_attributes ON products USING GIN (attributes);

-- GIN index on a specific path — smaller, faster if you only query one key
CREATE INDEX idx_products_brand ON products USING GIN ((attributes->'brand'));

-- BTREE index on extracted scalar value — fast =, <, > on a specific path
CREATE INDEX idx_products_battery ON products ((attributes->>'battery_hours'));
```

---

## 2. Full-Text Search

PostgreSQL full-text search converts text into lexemes (stemmed, normalized tokens) and builds a GIN index for fast lookup.

```sql
-- tsvector: the searchable representation of a document
SELECT to_tsvector('english', 'The quick brown fox jumped over lazy dogs');
-- 'brown':3 'dog':8 'fox':4 'jump':5 'lazi':7 'quick':2

-- tsquery: a search expression
SELECT to_tsquery('english', 'fox & dog');     -- fox AND dog
SELECT to_tsquery('english', 'fox | cat');     -- fox OR cat
SELECT plainto_tsquery('english', 'quick fox'); -- plain text → AND
SELECT websearch_to_tsquery('english', '"quick fox" -lazy'); -- Google-style

-- @@ match operator
SELECT to_tsvector('english', 'The quick brown fox') @@ to_tsquery('english', 'fox');  -- true
```

### Generated Columns for FTS

```sql
-- Add full-text search to an existing table
ALTER TABLE products ADD COLUMN search_vector TSVECTOR
    GENERATED ALWAYS AS (
        to_tsvector('english',
            coalesce(name, '') || ' ' ||
            coalesce(category, '') || ' ' ||
            coalesce(attributes->>'brand', '')
        )
    ) STORED;

-- GIN index on the generated column
CREATE INDEX idx_products_search ON products USING GIN (search_vector);

-- Search with ranking
SELECT
    id,
    name,
    price,
    ts_rank(search_vector, query) AS rank,
    ts_headline('english', name, query) AS headline
FROM products, websearch_to_tsquery('english', 'noise cancelling sony') AS query
WHERE search_vector @@ query
ORDER BY rank DESC
LIMIT 20;
```

### FTS in Go with sqlx

```go
type Product struct {
    ID       int64   `db:"id"`
    Name     string  `db:"name"`
    Price    float64 `db:"price"`
    Category string  `db:"category"`
    Rank     float64 `db:"rank"`
    Headline string  `db:"headline"`
}

func (r *ProductRepository) Search(ctx context.Context, query string, limit int) ([]*Product, error) {
    sql := `
        SELECT
            id, name, price, category,
            ts_rank(search_vector, to_tsquery('english', $1)) AS rank,
            ts_headline('english', name || ' ' || category, to_tsquery('english', $1),
                'MaxWords=15, MinWords=5, ShortWord=3, HighlightAll=false') AS headline
        FROM products
        WHERE search_vector @@ websearch_to_tsquery('english', $1)
        ORDER BY rank DESC
        LIMIT $2`
    
    var products []*Product
    if err := r.db.SelectContext(ctx, &products, sql, query, limit); err != nil {
        return nil, fmt.Errorf("search: %w", err)
    }
    return products, nil
}
```

---

## 3. Index Types

### B-tree (default)
Best for: equality, range, prefix on sortable types (integers, text, timestamps).

```sql
-- Default: CREATE INDEX creates a B-tree
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_orders_created ON orders (created_at DESC);

-- Composite B-tree: column order matters
-- This covers: WHERE status = ? ORDER BY created_at
-- But NOT: WHERE created_at > ? (leading column not used)
CREATE INDEX idx_orders_status_created ON orders (status, created_at DESC);

-- Partial index: only index the rows you actually query
CREATE INDEX idx_orders_pending ON orders (created_at)
WHERE status = 'pending';
-- Much smaller and faster than indexing all orders
```

### GIN — Generalized Inverted Index
Best for: arrays, JSONB, full-text search, tsvector. Each element maps to a list of row IDs.

```sql
-- Array containment
CREATE INDEX idx_tags ON posts USING GIN (tags);
SELECT * FROM posts WHERE tags @> '{go, performance}';

-- JSONB
CREATE INDEX idx_metadata ON events USING GIN (metadata);

-- Full-text
CREATE INDEX idx_search ON articles USING GIN (search_vector);
```

### GiST — Generalized Search Tree
Best for: geometric types, ranges, IP addresses (inet/cidr), PostGIS geometries.

```sql
-- Range queries (int4range, tstzrange, etc.)
CREATE INDEX idx_bookings_period ON bookings USING GIST (period);
-- Find overlapping bookings
SELECT * FROM bookings
WHERE period && tstzrange(NOW(), NOW() + INTERVAL '1 day');

-- IP address range
CREATE INDEX idx_sessions_ip ON sessions USING GIST (inet(ip_address) inet_ops);
```

### BRIN — Block Range Index
Best for: very large tables where column value correlates with physical storage order (logs, events, IoT time series). 10-100× smaller than B-tree. Slightly less precise (block-level, not row-level).

```sql
-- Perfect for append-only time-series: new rows always have newer timestamps
CREATE INDEX idx_events_ts ON events USING BRIN (created_at);
-- Index is ~8 KB even for a 1 TB table
```

### Index Strategies

```sql
-- EXPLAIN ANALYZE to see if an index is being used
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM orders WHERE user_id = 42 AND status = 'pending';

-- Check index sizes
SELECT
    schemaname || '.' || tablename AS table,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- Check index usage (is it being used at all?)
SELECT
    indexrelname AS index,
    idx_scan AS scans,
    idx_tup_read AS tuples_read
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;
-- Unused indexes (idx_scan = 0) just waste write performance and disk space — drop them
```

---

## 4. Advanced Queries in Go

### JSONB with pgx/sqlx

```go
import "encoding/json"

type Attributes map[string]interface{}

// Scan JSONB into a Go map
type Product struct {
    ID         int64      `db:"id"`
    Name       string     `db:"name"`
    Attributes Attributes `db:"attributes"`
}

// sqlx can scan JSONB into []byte; unmarshal manually
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
    var p struct {
        ID         int64  `db:"id"`
        Name       string `db:"name"`
        Attributes []byte `db:"attributes"`
    }
    err := r.db.QueryRowxContext(ctx,
        "SELECT id, name, attributes FROM products WHERE id = $1", id,
    ).StructScan(&p)
    if err != nil { return nil, err }
    
    var attrs Attributes
    if err := json.Unmarshal(p.Attributes, &attrs); err != nil {
        return nil, fmt.Errorf("unmarshal attributes: %w", err)
    }
    return &Product{ID: p.ID, Name: p.Name, Attributes: attrs}, nil
}

// Store JSONB from a Go value
func (r *ProductRepository) UpdateAttributes(ctx context.Context, id int64, attrs Attributes) error {
    data, err := json.Marshal(attrs)
    if err != nil { return fmt.Errorf("marshal: %w", err) }
    
    _, err = r.db.ExecContext(ctx,
        "UPDATE products SET attributes = $1 WHERE id = $2", data, id)
    return err
}
```

### Using pgx directly (recommended for JSONB)

```go
import "github.com/jackc/pgx/v5"

// pgx v5 natively handles JSONB with []byte or json.RawMessage
func (r *ProductRepository) GetAttributes(ctx context.Context, id int64) (json.RawMessage, error) {
    var attrs json.RawMessage
    err := r.db.QueryRow(ctx,
        "SELECT attributes FROM products WHERE id = $1", id,
    ).Scan(&attrs)
    return attrs, err
}

// Query JSONB with operators using pgx
func (r *ProductRepository) FindByAttribute(ctx context.Context, key, value string) ([]*Product, error) {
    query := fmt.Sprintf(`{"` + key + `": %q}`, value) // e.g. {"brand": "Sony"}
    rows, err := r.db.Query(ctx,
        "SELECT id, name, price FROM products WHERE attributes @> $1",
        query,
    )
    if err != nil { return nil, err }
    defer rows.Close()
    
    var products []*Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil { return nil, err }
        products = append(products, &p)
    }
    return products, rows.Err()
}
```

---

## Summary

**JSONB**:
- `->` returns JSONB; `->>` returns text
- `@>` for containment; `?` for key existence
- GIN index on the whole column or a specific path
- `jsonb_set` and `||` for updates; `#-` for deletes

**Full-text search**:
- `tsvector` = searchable representation; `tsquery` = search expression
- Use `GENERATED ALWAYS AS ... STORED` for search vectors to keep them updated automatically
- GIN index on the tsvector column
- `ts_rank` for relevance scoring; `ts_headline` for result snippets

**Index types**:
- B-tree: everything else (equality, range, sort)
- GIN: arrays, JSONB, tsvector
- GiST: ranges, geometric, IP
- BRIN: very large append-only tables

## Exercises

### Easy
1. Add a JSONB `metadata` column to a `users` table. Insert a user with `{"plan": "pro", "trial_ends": "2026-08-01"}`. Write a query that finds all users whose plan is "pro".
2. Create a `posts` table with a tsvector generated column. Insert 10 posts. Run a full-text search for "distributed systems" and order by rank.
3. Use `EXPLAIN ANALYZE` on a query with and without an index. Document the difference in total cost and execution time.

### Medium
4. Build a `ProductSearch` function in Go that accepts a query string and optional filters (category, min/max price). Combine full-text search with regular WHERE clauses. Use `websearch_to_tsquery` so users can type `"noise cancelling" -wired`.
5. Implement **faceted search**: given a search query, return not just the results but also aggregate counts by category. `SELECT category, COUNT(*) FROM products WHERE search_vector @@ $1 GROUP BY category`.
6. Write a benchmark comparing JSONB containment queries with and without a GIN index on a table of 100k rows. How many rows per second can you read with each?

### Hard
7. Design a **multi-tenant attribute schema** using JSONB: each tenant defines their own required fields. Store the schema as JSONB in a `tenant_schemas` table. Write a PostgreSQL function that validates a product's attributes against its tenant's schema before inserting.
8. Implement **real-time search with change tracking**: when a product is updated, automatically update its `search_vector`. Use a PostgreSQL trigger that calls `to_tsvector` and stores the result. Then implement a Go webhook that the trigger fires via `pg_notify`, which your server listens to via `LISTEN/NOTIFY`.
