# Chapter 116: Data Partitioning, Sharding, and Multi-Tenancy

A single PostgreSQL instance on the fastest available hardware will eventually hit its ceiling. Before that happens — usually long before — you need a strategy for partitioning data: splitting it so that each query touches less of it. This chapter covers vertical partitioning (splitting columns), horizontal partitioning (splitting rows, a.k.a. sharding), and the multi-tenancy problem: how to serve many customers from one database without letting their data bleed together.

## Table of Contents

1. [Vertical Partitioning](#1-vertical-partitioning)
2. [Horizontal Partitioning — Sharding](#2-horizontal-partitioning--sharding)
3. [PostgreSQL Table Partitioning](#3-postgresql-table-partitioning)
4. [Shard Key Selection](#4-shard-key-selection)
5. [Consistent Hashing for Application-Level Sharding](#5-consistent-hashing-for-application-level-sharding)
6. [Cross-Shard Queries and Denormalization](#6-cross-shard-queries-and-denormalization)
7. [Resharding: Adding a Shard](#7-resharding-adding-a-shard)
8. [Multi-Tenancy Models](#8-multi-tenancy-models)
9. [Row-Level Security in PostgreSQL](#9-row-level-security-in-postgresql)
10. [Go Middleware for Tenant Context](#10-go-middleware-for-tenant-context)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Vertical Partitioning

Vertical partitioning splits a wide table into a hot table (columns touched by most queries) and a cold table (large or rarely accessed columns). The classic offender is a `products` table that stores both the listing metadata and the raw HTML blob used only when rendering a product detail page.

```sql
-- Before: one wide table — every listing query loads blobs
CREATE TABLE products (
    id               BIGSERIAL PRIMARY KEY,
    name             TEXT        NOT NULL,
    price            NUMERIC(12,2) NOT NULL,
    description_html TEXT,
    raw_html_blob    TEXT,
    image_blob       BYTEA
);

-- After: split into hot + cold
CREATE TABLE products (
    id    BIGSERIAL PRIMARY KEY,
    name  TEXT          NOT NULL,
    price NUMERIC(12,2) NOT NULL
);

CREATE TABLE product_content (
    product_id       BIGINT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    description_html TEXT,
    raw_html_blob    TEXT,
    image_blob       BYTEA
);
```

Product listing queries never load blobs. The `product_content` row is fetched only when rendering the detail page, keeping the hot path narrow and cache-friendly.

A further step is to move binary blobs entirely out of the database into object storage (S3, GCS) and store only the URL:

```sql
ALTER TABLE product_content
    DROP COLUMN image_blob,
    ADD COLUMN  image_url TEXT;
```

```go
// UploadProductImage stores the blob in S3 and records the URL in the DB.
func (s *ProductService) UploadProductImage(ctx context.Context, productID int64, data []byte) error {
    key := fmt.Sprintf("products/%d/image.jpg", productID)
    url, err := s.s3.PutObject(ctx, key, data)
    if err != nil {
        return fmt.Errorf("s3 upload: %w", err)
    }

    _, err = s.db.ExecContext(ctx,
        `UPDATE product_content SET image_url = $1 WHERE product_id = $2`,
        url, productID,
    )
    return err
}
```

---

## 2. Horizontal Partitioning — Sharding

Horizontal partitioning keeps all columns but splits rows across multiple database instances, each called a **shard**. A routing function maps a shard key (usually a user ID or tenant ID) to the shard that owns that row.

```
               hash(user_id) % 3

User 101  →  shard 2   [DB-2: rows where user_id % 3 = 2]
User 202  →  shard 0   [DB-0: rows where user_id % 3 = 0]
User 303  →  shard 1   [DB-1: rows where user_id % 3 = 1]
```

Each shard is an independent PostgreSQL instance — its own host, storage, and connection pool. Together they hold 100 % of the data; individually each holds roughly 1/N.

```go
// ShardedDB routes queries to the correct shard by key.
type ShardedDB struct {
    shards []*sql.DB
}

func NewShardedDB(dsns []string) (*ShardedDB, error) {
    dbs := make([]*sql.DB, len(dsns))
    for i, dsn := range dsns {
        db, err := sql.Open("pgx", dsn)
        if err != nil {
            return nil, fmt.Errorf("shard %d: %w", i, err)
        }
        dbs[i] = db
    }
    return &ShardedDB{shards: dbs}, nil
}

// shardFor returns the DB connection responsible for the given key.
func (s *ShardedDB) shardFor(key int64) *sql.DB {
    idx := key % int64(len(s.shards))
    if idx < 0 {
        idx = -idx
    }
    return s.shards[idx]
}

func (s *ShardedDB) QueryRowContext(ctx context.Context, shardKey int64, query string, args ...any) *sql.Row {
    return s.shardFor(shardKey).QueryRowContext(ctx, query, args...)
}

func (s *ShardedDB) ExecContext(ctx context.Context, shardKey int64, query string, args ...any) (sql.Result, error) {
    return s.shardFor(shardKey).ExecContext(ctx, query, args...)
}
```

---

## 3. PostgreSQL Table Partitioning

PostgreSQL's built-in table partitioning is single-instance: one PostgreSQL server, but the rows are stored in separate physical child tables called **partitions**. The planner routes queries and `INSERT`s automatically.

**Range partitioning** is ideal for time-series data — old partitions can be dropped in O(1) (no `DELETE` needed):

```sql
-- Range partitioning by date (time-series logs)
CREATE TABLE events (
    id          BIGSERIAL,
    user_id     BIGINT NOT NULL,
    event_type  TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    payload     JSONB
) PARTITION BY RANGE (created_at);

CREATE TABLE events_2025_01 PARTITION OF events
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE events_2025_02 PARTITION OF events
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
```

**Hash partitioning** spreads rows evenly when there is no natural range boundary:

```sql
-- Hash partitioning for even distribution by user_id
CREATE TABLE orders (
    id       BIGSERIAL,
    user_id  BIGINT NOT NULL,
    total    NUMERIC(12,2),
    status   TEXT
) PARTITION BY HASH (user_id);

CREATE TABLE orders_p0 PARTITION OF orders FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE orders_p1 PARTITION OF orders FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE orders_p2 PARTITION OF orders FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE orders_p3 PARTITION OF orders FOR VALUES WITH (MODULUS 4, REMAINDER 3);
```

Key properties:
- A `WHERE created_at BETWEEN '2025-01-01' AND '2025-01-31'` query only scans `events_2025_01` (**partition pruning** — the other partitions are skipped entirely).
- Creating an index on the parent table automatically creates matching indexes on all child partitions.
- Dropping a partition (`DROP TABLE events_2024_01`) is instant and reclaims storage without a table scan.

---

## 4. Shard Key Selection

Choosing the wrong shard key is the most common sharding mistake. The key cannot be changed after data is loaded without a full resharding operation.

| Candidate key | Cardinality | Distribution | Query alignment | Verdict |
|---------------|-------------|--------------|-----------------|---------|
| `user_id` (UUID/int) | Very high | Even | Most queries filter by user | Excellent |
| `tenant_id` (SaaS) | High | Varies | Most queries are tenant-scoped | Excellent |
| `created_at` | Very high | Uneven — writes cluster at "now" | Range scans only | Bad (hotspot) |
| `status` (enum) | Very low | Skewed | Rarely the primary filter | Bad |
| `country_code` | Low (~200) | Skewed | Rarely primary filter | Bad |
| `boolean` flag | 2 | Extremely skewed | Rarely primary filter | Terrible |

The hotspot problem with time-based keys is severe:

```
Time-based shard key — all writes go to one "hot" shard:

  t=now ──────────────────────────► time
                                     │
                    shard 0 (cold)   │   ← only old reads
                    shard 1 (cold)   │   ← only old reads
                    shard 2 (HOT)    │   ← 100% of writes land here
                                     ▼
  Adding a new time partition does not redistribute existing load.
  Write throughput is bounded by the single hot shard.
```

For SaaS products, `tenant_id` is usually the right choice: most application queries are already scoped to one tenant, so the routing is free — the filter that was already in your `WHERE` clause becomes the shard selector.

---

## 5. Consistent Hashing for Application-Level Sharding

Simple modulo routing (`key % N`) breaks when you add a shard: every key potentially maps to a different shard, forcing a full data migration. Consistent hashing was covered in Chapter 97; here we apply it to DB routing.

```go
// ConsistentHashRing maps string keys to named nodes with minimal remapping
// when nodes are added or removed. (Implementation omitted — see Ch. 97.)
type ConsistentHashRing struct{ /* ... */ }

func (r *ConsistentHashRing) GetNode(key string) string { /* ... */ return "" }
func (r *ConsistentHashRing) AddNode(node string)        { /* ... */ }
func (r *ConsistentHashRing) RemoveNode(node string)     { /* ... */ }

// DBRouter uses a consistent hash ring to select the shard for a given key.
type DBRouter struct {
    ring   *ConsistentHashRing
    shards map[string]*sql.DB // node name → DB connection
}

func NewDBRouter(nodes map[string]string) (*DBRouter, error) {
    ring := &ConsistentHashRing{}
    shards := make(map[string]*sql.DB, len(nodes))
    for name, dsn := range nodes {
        ring.AddNode(name)
        db, err := sql.Open("pgx", dsn)
        if err != nil {
            return nil, fmt.Errorf("open shard %s: %w", name, err)
        }
        shards[name] = db
    }
    return &DBRouter{ring: ring, shards: shards}, nil
}

func (r *DBRouter) DBFor(key string) *sql.DB {
    node := r.ring.GetNode(key)
    return r.shards[node]
}

func (r *DBRouter) QueryRow(ctx context.Context, shardKey string, query string, args ...any) *sql.Row {
    db := r.DBFor(shardKey)
    return db.QueryRowContext(ctx, query, args...)
}

func (r *DBRouter) Exec(ctx context.Context, shardKey string, query string, args ...any) (sql.Result, error) {
    db := r.DBFor(shardKey)
    return db.ExecContext(ctx, query, args...)
}
```

When a new shard is added, only ~1/N keys need to move (the keys that were hashed to the arc of the ring now claimed by the new node). All other keys continue pointing to the same shard unchanged.

---

## 6. Cross-Shard Queries and Denormalization

The hardest constraint in a sharded system is that **joins across shards are not possible at the database level**. A query like:

```sql
SELECT o.id, u.email FROM orders o JOIN users u ON o.user_id = u.id
```

breaks when `orders` and `users` may reside on different shards.

The rule: **avoid cross-shard joins — denormalize what you need**. Embed frequently read fields from the joined table directly into the fact table:

```sql
-- Denormalized: user fields embedded in orders
CREATE TABLE orders (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    user_email TEXT   NOT NULL,  -- duplicated from users table
    user_name  TEXT   NOT NULL,  -- duplicated from users table
    total      NUMERIC(12,2),
    status     TEXT
);
```

The data may go stale if the user changes their email, but for most order history use cases that is an acceptable tradeoff. What you lose in normalisation you gain in query simplicity.

For aggregates that span all shards (total revenue, active user counts), fan the query out to every shard and reduce the results in Go:

```go
// FanOutQuery runs the same query on every shard and collects rows.
func (r *DBRouter) FanOutQuery(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
    type shardResult struct {
        rows []map[string]any
        err  error
    }

    results := make(chan shardResult, len(r.shards))

    for name, db := range r.shards {
        go func(name string, db *sql.DB) {
            rows, err := db.QueryContext(ctx, query, args...)
            if err != nil {
                results <- shardResult{err: fmt.Errorf("shard %s: %w", name, err)}
                return
            }
            defer rows.Close()

            cols, _ := rows.Columns()
            var out []map[string]any
            for rows.Next() {
                vals := make([]any, len(cols))
                ptrs := make([]any, len(cols))
                for i := range vals { ptrs[i] = &vals[i] }
                rows.Scan(ptrs...)
                row := make(map[string]any, len(cols))
                for i, c := range cols { row[c] = vals[i] }
                out = append(out, row)
            }
            results <- shardResult{rows: out}
        }(name, db)
    }

    var all []map[string]any
    for range r.shards {
        res := <-results
        if res.err != nil { return nil, res.err }
        all = append(all, res.rows...)
    }
    return all, nil
}
```

For complex analytics, maintain a dedicated analytics database (a read replica that aggregates from all shards) rather than fanning out on every request.

---

## 7. Resharding: Adding a Shard

When you add a new shard, ~1/N of existing keys must migrate. The naive approach (stop writes, copy, resume) requires downtime. The **double-write pattern** achieves zero downtime:

```
Step 1  Start writing to both old shard and new shard for the affected key range.
Step 2  Backfill: copy historical data from old shard to new shard.
Step 3  Verify: read from new shard, compare checksums with old shard.
Step 4  Cut reads over to new shard.
Step 5  Stop writing to old shard.
Step 6  Delete migrated data from old shard.
```

```go
// DoubleWriteDB wraps two DB connections during a shard migration window.
type DoubleWriteDB struct {
    primary   *sql.DB // old shard — always the source of truth
    secondary *sql.DB // new shard — being populated
    active    bool    // flip to true when migration starts
}

func (d *DoubleWriteDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
    result, err := d.primary.ExecContext(ctx, query, args...)
    if err != nil {
        return result, err
    }
    if d.active {
        // Best-effort write to secondary; don't fail the request if it fails.
        go d.secondary.ExecContext(context.Background(), query, args...)
    }
    return result, nil
}

// QueryRowContext reads from secondary once verified, primary otherwise.
func (d *DoubleWriteDB) QueryRowContext(ctx context.Context, readFromSecondary bool, query string, args ...any) *sql.Row {
    if d.active && readFromSecondary {
        return d.secondary.QueryRowContext(ctx, query, args...)
    }
    return d.primary.QueryRowContext(ctx, query, args...)
}
```

The secondary write is fire-and-forget (`go ...`). Any gaps are filled by the backfill job in Step 2. The primary remains authoritative until Step 4 flips the read flag.

---

## 8. Multi-Tenancy Models

Three standard models exist, each with a different isolation/cost tradeoff:

| Model | Isolation | Cost | Complexity | Best for |
|-------|-----------|------|------------|----------|
| One DB per tenant | Strongest | Highest (N DB instances) | High (N connection pools) | Enterprise, compliance requirements |
| One schema per tenant | Good | Medium (one DB, N schemas) | Medium | Mid-market SaaS |
| Shared tables + `tenant_id` | Weakest | Lowest | Low (add `WHERE tenant_id=?`) | High-volume SMB SaaS |

**Schema-per-tenant** uses PostgreSQL's `search_path` to route queries transparently:

```sql
-- Provision a new tenant
CREATE SCHEMA tenant_abc123;

CREATE TABLE tenant_abc123.orders (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    total      NUMERIC(12,2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- At query time — set search_path to the tenant's schema
SET search_path = tenant_abc123;
SELECT * FROM orders;  -- hits tenant_abc123.orders
```

**Shared tables with `tenant_id`** is the simplest model. The composite index puts `tenant_id` first so all per-tenant queries hit a tight index range:

```sql
CREATE TABLE orders (
    id         BIGSERIAL PRIMARY KEY,
    tenant_id  UUID NOT NULL,
    user_id    BIGINT NOT NULL,
    total      NUMERIC(12,2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Composite index: tenant_id first so queries within a tenant are fast
CREATE INDEX idx_orders_tenant_user ON orders(tenant_id, user_id);
```

Every query must carry a `WHERE tenant_id = $1` clause. Forgetting it returns all tenants' data — which is why Row-Level Security exists.

---

## 9. Row-Level Security in PostgreSQL

Row-Level Security (RLS) enforces tenant isolation at the database engine level, independent of application code. Even if a bug causes missing `WHERE tenant_id=?` clauses, the policy blocks cross-tenant reads.

```sql
-- Enable RLS on the table
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- Policy: every SELECT/UPDATE/DELETE automatically filters to the current tenant
CREATE POLICY tenant_isolation ON orders
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

-- The app connects as 'app_user', not as a superuser.
-- Superusers bypass RLS; ordinary roles do not.
GRANT SELECT, INSERT, UPDATE, DELETE ON orders TO app_user;
```

`current_setting('app.tenant_id')` reads a PostgreSQL session variable. The application sets it at the start of each transaction:

```sql
SET LOCAL app.tenant_id = '550e8400-e29b-41d4-a716-446655440000';
```

`SET LOCAL` is transaction-scoped — it is automatically cleared when the transaction commits or rolls back, so there is no risk of the value leaking between requests when connection pooling is in use.

---

## 10. Go Middleware for Tenant Context

The full tenant isolation pipeline in Go: HTTP middleware extracts the tenant ID, stores it in `context.Context`, and a DB helper sets the PostgreSQL session variable before any queries run.

```go
import (
    "context"
    "database/sql"
    "fmt"
    "net/http"

    "github.com/google/uuid"
)

type contextKey string

const tenantKey contextKey = "tenantID"

// TenantMiddleware extracts the tenant ID from the X-Tenant-ID header
// (or from JWT claims in production) and injects it into the request context.
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        raw := r.Header.Get("X-Tenant-ID")
        if raw == "" {
            http.Error(w, "missing X-Tenant-ID", http.StatusUnauthorized)
            return
        }

        tenantID, err := uuid.Parse(raw)
        if err != nil {
            http.Error(w, "invalid tenant ID", http.StatusBadRequest)
            return
        }

        ctx := context.WithValue(r.Context(), tenantKey, tenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// TenantFromContext retrieves the tenant ID stored by TenantMiddleware.
func TenantFromContext(ctx context.Context) (uuid.UUID, bool) {
    id, ok := ctx.Value(tenantKey).(uuid.UUID)
    return id, ok
}

// WithTenantTx opens a transaction, sets the tenant session variable via
// SET LOCAL, then calls fn. SET LOCAL is undone automatically on COMMIT/ROLLBACK.
func WithTenantTx(ctx context.Context, db *sql.DB, tenantID uuid.UUID, fn func(tx *sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    if _, err = tx.ExecContext(ctx, "SET LOCAL app.tenant_id = $1", tenantID); err != nil {
        tx.Rollback()
        return fmt.Errorf("set tenant context: %w", err)
    }

    if err = fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
}

// Example handler that ties middleware + WithTenantTx together.
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
    tenantID, ok := TenantFromContext(r.Context())
    if !ok {
        http.Error(w, "no tenant in context", http.StatusUnauthorized)
        return
    }

    var orders []Order
    err := WithTenantTx(r.Context(), h.db, tenantID, func(tx *sql.Tx) error {
        // RLS policy enforces tenant_id automatically — no WHERE clause needed,
        // but adding it is still good practice for index usage.
        rows, err := tx.QueryContext(r.Context(),
            `SELECT id, user_id, total, created_at FROM orders ORDER BY created_at DESC LIMIT 50`)
        if err != nil {
            return err
        }
        defer rows.Close()
        for rows.Next() {
            var o Order
            rows.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt)
            orders = append(orders, o)
        }
        return rows.Err()
    })

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSON(w, orders)
}
```

---

## Summary

- **Vertical partitioning** splits a wide table into hot (frequently accessed) and cold (rarely accessed) columns; move binary blobs to S3 and store only the URL in the DB
- **Horizontal partitioning (sharding)** splits rows across multiple database instances using a routing function on a shard key; each shard holds an independent subset of rows
- **PostgreSQL table partitioning** is single-instance: range partitioning suits time-series data (drop old partitions instantly); hash partitioning distributes rows evenly when there is no natural range
- **Shard key selection** is critical: choose high-cardinality, evenly distributed keys that align with query patterns; time-based keys create write hotspots; for SaaS, `tenant_id` is usually ideal
- **Consistent hashing** minimises data movement when shards are added or removed: only ~1/N keys remap rather than all of them
- **Cross-shard joins are impossible** at the DB level; denormalize required fields into the fact table and fan aggregation queries out to all shards in Go
- **Resharding** without downtime uses the double-write pattern: write to both old and new shard, backfill, verify, cut over reads, then stop writing to the old shard
- **Multi-tenancy models**: DB-per-tenant (strongest isolation), schema-per-tenant (middle ground), shared tables + `tenant_id` (lowest cost); choose based on isolation requirements and scale
- **Row-Level Security** enforces tenant isolation at the PostgreSQL engine level via `CREATE POLICY` and `current_setting('app.tenant_id')`; `SET LOCAL` is transaction-scoped and safe with connection pools
- **Go middleware** injects the tenant ID into `context.Context` via HTTP middleware and uses `WithTenantTx` to bind the PostgreSQL session variable before any queries execute

---

## Exercises

### Easy
1. Write the DDL for a `logs` table using range partitioning by `created_at`, with partitions for January, February, and March 2025. Add the appropriate index on the parent table.
2. Implement the `ShardedDB` struct with a `shardFor(key int64) *sql.DB` method and `QueryRowContext`/`ExecContext` wrappers. Write a unit test that verifies calls with the same key always route to the same shard and calls with keys that hash to different remainders route to different shards.
3. Create a PostgreSQL RLS policy for a `projects` table with a `tenant_id UUID` column. Enable RLS, create the policy using `current_setting('app.tenant_id')::uuid`, and write the `GRANT` statement for `app_user`. Verify by attempting a query without setting the session variable.

### Medium
4. Implement a `DBRouter` backed by a consistent hash ring. Populate it with three named shards (`shard-a`, `shard-b`, `shard-c`). Write a test that adds a fourth shard (`shard-d`) and counts how many of 1,000 test keys change their assigned shard — it must be approximately 250 (1/4), not all 1,000.
5. Implement `DoubleWriteDB` with an `active` flag, primary and secondary `*sql.DB` fields, and `ExecContext`. Write a test that enables the double-write mode and confirms that after 100 writes, both primary and secondary databases contain the same row count. Ensure secondary write failures do not propagate errors to the caller.
6. Implement a schema-per-tenant connection router: given a `tenantID string`, it sets `search_path = <tenantID>` on a connection and returns a `*sql.Tx`. Wrap it in a `WithSchemaTx(ctx, db, tenantID, fn)` function. Write a test that creates two schemas with separate `orders` tables, inserts rows into each, and confirms that queries via `WithSchemaTx` only see rows for the correct tenant.

### Hard
7. Build a zero-downtime resharding tool for `ShardedDB`. It should: (a) accept an old shard list and a new shard list; (b) identify which keys from a provided key set need to migrate (keys whose target shard differs between old and new routing); (c) copy rows for those keys from the old shard to the new shard using `INSERT … ON CONFLICT DO NOTHING`; (d) verify row counts match; (e) print a migration report. Test it by seeding 10,000 rows across three shards, adding a fourth shard, and confirming no data loss after migration.
8. Implement a multi-tenancy benchmark that compares all three models (DB-per-tenant, schema-per-tenant, shared tables + RLS) using `testing.B`. Provision three tenants per model in a test PostgreSQL instance. Each benchmark iteration should: insert one row and then query the 50 most recent rows for a tenant. Report ns/op and allocs/op for each model and explain in comments which model wins at different tenant counts and why.
