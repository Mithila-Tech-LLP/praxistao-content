# Chapter 135: Database at Scale — Read Replicas, pgBouncer, and Connection Pooling

A single PostgreSQL instance handles most workloads well. But once you have dozens of pods, millions of rows, and thousands of concurrent users, you hit three distinct walls: connection limits, read throughput, and table scan performance. This chapter covers the standard tools for breaking through each wall.

---

## The Scaling Problem: Connection Limits

PostgreSQL's default `max_connections` is 100. That sounds like a lot until you factor in how Go services actually deploy.

A typical service has a `database/sql` pool configured like this:

```go
db, err := sql.Open("pgx", dsn)
if err != nil {
    log.Fatal(err)
}
db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

With 50 pods running, you have:

```
50 pods × 10 connections = 500 connections
```

PostgreSQL refuses anything above `max_connections`. New connection attempts see:

```
pq: sorry, too many clients already
```

You can raise `max_connections` in `postgresql.conf`, but there is a hard cost: each PostgreSQL connection is a separate OS process. Each process consumes roughly 5-10 MB of RAM. At 500 connections:

```
500 connections × 7 MB average = 3.5 GB RAM just for connection overhead
```

That memory is not available for shared buffers, query plans, or actual data caching. Raising `max_connections` to 1000 to handle 100 pods just moves the problem and burns more RAM.

There are two approaches to solving this: connection pooling (reduce the number of connections that reach PostgreSQL) and read replicas (distribute query load across multiple instances).

---

## Connection Pooling with pgBouncer

pgBouncer is a lightweight connection pooler that sits between your application and PostgreSQL. Applications connect to pgBouncer; pgBouncer maintains a small pool of real connections to PostgreSQL.

```
[Pod 1 pool: 10 conns ]
[Pod 2 pool: 10 conns ]
[Pod 3 pool: 10 conns ]     pgBouncer            PostgreSQL
        ...              (1000 clients in  -->  (20 server conns)
[Pod 99 pool: 10 conns ]    20 server conns)
[Pod 100 pool: 10 conns]
```

The application code does not change. The connection string points to pgBouncer's host and port instead of PostgreSQL directly.

### Pooling Modes

pgBouncer has three pooling modes. The mode determines when a server connection is returned to the pool.

**Session mode**: a client holds a server connection for the entire session. This behaves like a direct connection. Useful for compatibility but offers minimal connection savings.

**Transaction mode**: the server connection is returned to the pool after each transaction commits or rolls back. This is the right mode for Go services. A pod with 10 idle goroutines holds zero server connections; it only borrows one when a transaction is in flight.

**Statement mode**: the server connection is returned after every single SQL statement. You cannot use multi-statement transactions in this mode. Rarely useful in practice.

Use transaction mode unless you have a specific reason not to.

### pgBouncer Configuration

pgBouncer is configured with an INI file:

```ini
[databases]
mydb = host=127.0.0.1 port=5432 dbname=mydb

[pgbouncer]
listen_addr = 0.0.0.0
listen_port = 6432
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt

; Pooling mode
pool_mode = transaction

; Maximum client connections pgBouncer accepts
max_client_conn = 1000

; How many server connections per user+database pair
default_pool_size = 20

; Hard limit on server connections per pool
server_pool_size = 10

; Log connections
log_connections = 0
log_disconnections = 0

; How long to wait for a server connection before returning error
query_wait_timeout = 120

; Keep server connections open when idle
server_idle_timeout = 600
```

The `userlist.txt` stores credentials:

```
"appuser" "md5hashofpassword"
```

Your application connects to pgBouncer on port 6432 exactly as it would connect to PostgreSQL on port 5432:

```go
dsn := "postgres://appuser:password@pgbouncer-host:6432/mydb"
db, err := sql.Open("pgx", dsn)
```

### Why Not Just Use database/sql Pools?

`database/sql` pools connections per process. With 100 pods each holding 10 connections, you have 1000 connections to PostgreSQL. pgBouncer centralizes this:

```
Without pgBouncer:
  100 pods × 10 connections = 1000 PostgreSQL connections

With pgBouncer:
  100 pods × 10 connections to pgBouncer = 1000 pgBouncer clients
  pgBouncer maintains 20 server connections to PostgreSQL
```

The 100 pods see no difference. PostgreSQL sees 20 connections instead of 1000.

### Caveat: Session-Level Features

Transaction pooling breaks features that require a persistent server session. The server connection changes between transactions, so anything set with `SET` or registered with `LISTEN` vanishes.

Specifically, **prepared statements** are problematic. When your driver sends `PREPARE`, that prepared statement lives on a specific server connection. The next transaction might land on a different server connection that does not have it.

Solutions:

1. Disable prepared statements in the driver. With `pgx`:

```go
config, err := pgx.ParseConfig(dsn)
config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
db = stdlib.OpenDB(*config)
```

2. Use pgBouncer 1.21+ which has `prepared_statement_cache_size` and handles prepared statement proxying transparently.

3. Set `server_reset_query = DEALLOCATE ALL` in pgbouncer.ini so prepared statements are cleaned up when a server connection is returned to the pool.

Advisory locks (`pg_advisory_lock`), `LISTEN`/`NOTIFY`, and `SET LOCAL` for session variables all require either session mode or careful redesign.

---

## Read Replicas

PostgreSQL streaming replication sends the Write-Ahead Log (WAL) from the primary to one or more replica servers. Replicas apply the WAL and stay synchronized. Replicas are read-only; all writes must go to the primary.

```
                WAL stream
Primary  ─────────────────────>  Replica 1
         ─────────────────────>  Replica 2

Writes --> Primary
Reads  --> Replica 1 or Replica 2
```

Reads typically represent 80-90% of database load. Routing them to replicas lets the primary focus on writes and frees up IOPS and CPU.

### Two-Pool Pattern in Go

The cleanest Go pattern is two separate connection pools:

```go
package db

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
)

type Pools struct {
    Write *sql.DB
    Read  *sql.DB
}

func NewPools(writeDSN, readDSN string) (*Pools, error) {
    write, err := openPool(writeDSN)
    if err != nil {
        return nil, fmt.Errorf("write pool: %w", err)
    }

    read, err := openPool(readDSN)
    if err != nil {
        return nil, fmt.Errorf("read pool: %w", err)
    }

    return &Pools{Write: write, Read: read}, nil
}

func openPool(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(2 * time.Minute)
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
```

Initialize with separate DSNs pointing at pgBouncer instances in front of the primary and replica:

```go
pools, err := db.NewPools(
    "postgres://app:pass@pgbouncer-primary:6432/mydb",
    "postgres://app:pass@pgbouncer-replica:6432/mydb",
)
```

### Smart Routing in the Repository Layer

The routing decision lives in the repository, not in business logic. Reads use `dbRead`; writes use `dbWrite`.

```go
package order

import (
    "context"
    "database/sql"
    "time"
)

type Order struct {
    ID        int64
    UserID    int64
    Total     float64
    CreatedAt time.Time
}

type Repository struct {
    dbWrite *sql.DB
    dbRead  *sql.DB
}

func NewRepository(dbWrite, dbRead *sql.DB) *Repository {
    return &Repository{dbWrite: dbWrite, dbRead: dbRead}
}

// GetOrder reads from replica — safe for read replica.
func (r *Repository) GetOrder(ctx context.Context, id int64) (*Order, error) {
    const q = `
        SELECT id, user_id, total, created_at
        FROM orders
        WHERE id = $1
    `
    // safe for read replica
    row := r.dbRead.QueryRowContext(ctx, q, id)
    var o Order
    if err := row.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt); err != nil {
        return nil, err
    }
    return &o, nil
}

// ListOrders reads from replica — safe for read replica.
func (r *Repository) ListOrders(ctx context.Context, userID int64) ([]*Order, error) {
    const q = `
        SELECT id, user_id, total, created_at
        FROM orders
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 100
    `
    // safe for read replica
    rows, err := r.dbRead.QueryContext(ctx, q, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []*Order
    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt); err != nil {
            return nil, err
        }
        orders = append(orders, &o)
    }
    return orders, rows.Err()
}

// CreateOrder writes to primary.
func (r *Repository) CreateOrder(ctx context.Context, userID int64, total float64) (*Order, error) {
    const q = `
        INSERT INTO orders (user_id, total, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id, user_id, total, created_at
    `
    row := r.dbWrite.QueryRowContext(ctx, q, userID, total)
    var o Order
    if err := row.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt); err != nil {
        return nil, err
    }
    return &o, nil
}

// UpdateOrder writes to primary.
func (r *Repository) UpdateOrder(ctx context.Context, id int64, total float64) error {
    const q = `UPDATE orders SET total = $1 WHERE id = $2`
    _, err := r.dbWrite.ExecContext(ctx, q, total, id)
    return err
}

// DeleteOrder writes to primary.
func (r *Repository) DeleteOrder(ctx context.Context, id int64) error {
    const q = `DELETE FROM orders WHERE id = $1`
    _, err := r.dbWrite.ExecContext(ctx, q, id)
    return err
}
```

### Replication Lag

Streaming replication is asynchronous by default. The replica applies WAL after the primary commits. This lag is typically 10-100ms on a healthy setup but can grow under write load or network issues.

Check replication lag on the primary:

```sql
SELECT
    application_name,
    state,
    replay_lag,
    write_lag,
    flush_lag
FROM pg_stat_replication;
```

The lag creates a consistency problem. A user creates an order, the handler returns the new order ID, and the client immediately calls GET /orders/{id}. If that GET routes to the replica and the replica has not yet received that WAL segment, it returns `not found` or stale data.

**Solutions:**

**Option A — Primary read window.** After a write, store a timestamp in a per-user or per-session store (Redis, a cookie, a signed token). For the next N seconds, any read for that user goes to the primary.

**Option B — Write timestamp in context.** Inject the write timestamp into the request context. The repository checks it before choosing which pool to query. See the Hard exercise at the end of this chapter for a full implementation.

**Option C — Primary for the immediate response.** The handler that processes the write reads from the primary to construct its response. Subsequent reads can use the replica because the user already has the data they just wrote.

---

## PgCat: A Modern Pooler

PgCat is a PostgreSQL pooler written in Rust, originally built at Instacart. It handles connection pooling like pgBouncer but adds features that eliminate the need for the two-pool pattern in application code.

PgCat includes a query parser that inspects each SQL statement. `SELECT` statements are automatically routed to replicas; everything else goes to the primary. The application uses a single connection string pointing at PgCat, and routing is transparent.

Other features: health checks on replicas (automatically stops routing to lagging replicas), load balancing across multiple replicas, and a Prometheus metrics endpoint.

PgCat is worth evaluating for new deployments. For existing services already using the two-pool pattern, the migration cost may not justify the switch immediately. pgBouncer remains more battle-tested and simpler to operate.

---

## Partitioning for Scale

When a table grows past roughly 50-100 million rows, queries slow down even with proper indexes. The index itself becomes large; vacuum takes longer; autovacuum struggles to keep up. Partitioning splits the table into smaller physical segments while keeping the logical table interface unchanged.

### Range Partitioning by Date

```sql
CREATE TABLE orders (
    id         BIGSERIAL,
    user_id    BIGINT NOT NULL,
    total      NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);
```

Create partitions manually:

```sql
CREATE TABLE orders_2024_01
    PARTITION OF orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE orders_2024_02
    PARTITION OF orders
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

Or use `pg_partman` to automate partition creation and retention:

```sql
CREATE EXTENSION pg_partman;

SELECT pg_partman.create_parent(
    p_parent_table   => 'public.orders',
    p_control        => 'created_at',
    p_interval       => 'monthly',
    p_premake        => 3,     -- create 3 future partitions in advance
    p_start_partition => '2024-01-01'
);

-- Configure retention: drop partitions older than 12 months
UPDATE pg_partman.part_config
SET retention = '12 months',
    retention_keep_table = false
WHERE parent_table = 'public.orders';
```

Run the partition maintenance job (typically called from a cron or pg_cron):

```sql
SELECT pg_partman.run_maintenance('public.orders');
```

### Partition Pruning

When a query includes a filter on `created_at`, the query planner scans only the relevant partitions:

```sql
-- Only scans orders_2024_03, not the full table
SELECT * FROM orders
WHERE created_at >= '2024-03-01' AND created_at < '2024-04-01';
```

Each partition has its own index, smaller than a single global index on the full table. Vacuum runs per partition and completes faster. Dropping old data is instant (`DROP TABLE orders_2024_01`) rather than a slow `DELETE`.

For partitioning to work in Go with `database/sql`, no changes are required. The table name stays `orders`; PostgreSQL routes rows to the correct partition transparently.

---

## VACUUM and Table Bloat

PostgreSQL uses MVCC (Multi-Version Concurrency Control). When a row is updated or deleted, the old version is not immediately removed. It is marked dead and left in place until VACUUM collects it. These dead tuples are called bloat.

### Finding Bloated Tables

```sql
SELECT
    relname,
    n_live_tup,
    n_dead_tup,
    ROUND(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) AS dead_pct,
    last_autovacuum,
    last_autoanalyze
FROM pg_stat_user_tables
ORDER BY n_dead_tup DESC
LIMIT 20;
```

### Autovacuum Thresholds

The default autovacuum trigger is `autovacuum_vacuum_scale_factor = 0.2`. Autovacuum runs when dead tuples exceed 20% of the table. For a 100-million-row table:

```
100,000,000 × 0.2 = 20,000,000 dead rows before autovacuum fires
```

That is a large accumulation of bloat. On high-write tables, lower the threshold:

```sql
ALTER TABLE orders
    SET (autovacuum_vacuum_scale_factor = 0.01);
```

Now autovacuum fires at 1% dead tuples — 1 million rows for that same table. For very large tables, use an absolute threshold instead:

```sql
ALTER TABLE orders
    SET (
        autovacuum_vacuum_scale_factor = 0.0,
        autovacuum_vacuum_threshold = 10000
    );
```

This fires vacuum when there are 10,000 dead tuples, regardless of table size.

### Measuring Actual Bloat

`pg_stat_user_tables` estimates dead tuples. For a precise measurement, use `pgstattuple`:

```sql
CREATE EXTENSION pgstattuple;

SELECT
    table_len,
    tuple_count,
    dead_tuple_count,
    dead_tuple_percent,
    free_space,
    free_percent
FROM pgstattuple('orders');
```

If `dead_tuple_percent` is above 10-15%, run `VACUUM ANALYZE orders` manually or investigate why autovacuum is not keeping up (it may be blocked by long-running transactions).

---

## When to Scale Beyond Single PostgreSQL

Read replicas, pgBouncer, and partitioning can take a single PostgreSQL instance very far. A well-tuned single primary with two read replicas handles tens of thousands of queries per second. Exhaust these options before considering distributed databases.

### Citus: Horizontal Sharding

Citus is a PostgreSQL extension that distributes table rows across worker nodes. You choose a distribution column (e.g., `tenant_id`), and Citus shards rows by hashing it. Each worker holds a subset of the data.

```sql
-- On Citus coordinator
SELECT create_distributed_table('orders', 'tenant_id');
```

Queries that include the distribution key are fast — they hit only the relevant shard. Cross-shard queries (aggregations across all tenants, joins without the distribution key) are slow because the coordinator must gather results from all workers.

Citus fits multi-tenant SaaS well: each tenant's data lives on one shard, and most queries are scoped to a single tenant.

### Aurora PostgreSQL

Amazon Aurora PostgreSQL is a managed engine compatible with PostgreSQL. Storage auto-scales up to 128 TB. It supports up to 15 read replicas with replica lag typically under 10ms. Aurora Global Database replicates across AWS regions with a recovery point objective under 1 second.

Aurora removes the operational burden of managing replication. If you are already on AWS and want read replicas without running your own replication setup, Aurora is the practical default.

### Managed Alternatives

**Neon**: serverless PostgreSQL with branching and scale-to-zero. Useful for development environments and low-traffic services that benefit from branching for testing migrations.

**Supabase**: PostgreSQL-based backend-as-a-service with built-in auth, real-time subscriptions via `LISTEN`/`NOTIFY`, and a REST API generated from the schema.

**PlanetScale**: MySQL-based (not PostgreSQL), mentioned here only because it appears in the same conversations. Not relevant for PostgreSQL workloads.

### Rule of Thumb

Follow this order before reaching for distributed systems:

```
1. Vertical scaling (bigger instance: more CPU, more RAM, faster NVMe)
2. pgBouncer (reduce connection overhead)
3. Read replicas (distribute read load)
4. Partitioning (manage large tables)
5. Caching (Redis/Memcached for hot reads)
6. Horizontal sharding (Citus, distributed PostgreSQL)
```

Each step adds operational complexity. Most production services never need step 6.

---

## Summary

A single PostgreSQL instance has three main scaling constraints: connection limits, read throughput, and large table performance.

pgBouncer addresses connection limits by multiplexing thousands of application connections into a small pool of real server connections. Transaction mode gives the best efficiency for Go services. The application code does not change — only the connection string target changes.

Read replicas address read throughput by offloading SELECT queries to replicas that apply WAL from the primary. In Go, the two-pool pattern (`dbWrite` and `dbRead`) routes queries at the repository layer. Replication lag (10-100ms typical) requires careful handling for read-your-writes consistency.

PgCat is a modern alternative to pgBouncer that adds transparent replica routing, removing the need for the two-pool pattern.

Table partitioning addresses large table performance. Range partitioning by date, automated with `pg_partman`, enables partition pruning and keeps index and vacuum costs manageable.

MVCC bloat accumulates on high-write tables. Lowering `autovacuum_vacuum_scale_factor` on large tables prevents bloat from building up between vacuum runs.

When single-instance PostgreSQL is genuinely exhausted, Citus provides horizontal sharding (best for multi-tenant workloads) and Aurora PostgreSQL provides managed read replicas and storage scaling.

---

## Exercises

### Easy

Write a complete `pgbouncer.ini` configuration file for a local development setup. Requirements: connect to a PostgreSQL database named `appdb` running on `localhost:5432`; listen on port `6432`; use transaction pooling mode; allow a maximum of 500 client connections; use a default pool size of 10 server connections per database/user pair; set query wait timeout to 30 seconds.

### Medium

Refactor the following single-pool repository into a dual-pool repository that routes reads to the replica and writes to the primary:

```go
type OrderRepo struct {
    db *sql.DB
}

func (r *OrderRepo) GetOrder(ctx context.Context, id int64) (*Order, error)
func (r *OrderRepo) ListOrders(ctx context.Context, userID int64) ([]*Order, error)
func (r *OrderRepo) CreateOrder(ctx context.Context, userID int64, total float64) (*Order, error)
func (r *OrderRepo) UpdateOrder(ctx context.Context, id int64, total float64) error
func (r *OrderRepo) DeleteOrder(ctx context.Context, id int64) error
```

Replace the single `*sql.DB` field with `dbWrite` and `dbRead`. `GetOrder` and `ListOrders` should use `dbRead`. `CreateOrder`, `UpdateOrder`, and `DeleteOrder` should use `dbWrite`. Write a constructor that accepts two DSNs and returns an initialized `*OrderRepo` or an error.

### Hard

Implement read-your-writes consistency to avoid the post-write stale read problem with replicas.

Requirements:

1. Define a context key type and a helper that injects a `time.Time` write timestamp into a `context.Context`.

2. Write an HTTP middleware that reads a write timestamp from a header (e.g., `X-Write-Time`) and injects it into the request context. On responses from write endpoints, set this header to the current time.

3. In `GetOrder`, check the context for a write timestamp. If the timestamp exists and is less than 2 seconds ago, read from `dbWrite` (primary). If the timestamp is absent or older than 2 seconds, read from `dbRead` (replica).

4. Write a test that:
   - Calls `CreateOrder` and captures the returned timestamp.
   - Immediately calls `GetOrder` with a context carrying that timestamp — verify it hits `dbWrite`.
   - Creates a context with a timestamp 3 seconds in the past — verify `GetOrder` hits `dbRead`.

Use `*sql.DB` mock or an interface-based approach to verify which pool is queried in the test without needing a real database connection.
