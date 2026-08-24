# Chapter 13: PostgreSQL — The Reliable Giant

PostgreSQL is used by Instagram, Shopify, Reddit, Twitch, and millions of other applications. It's been under active development since 1986 and is widely considered the most feature-rich and reliable open-source database in existence. This chapter gives you the foundation to use it confidently.

## Table of Contents

1. What Makes PostgreSQL Special
2. PostgreSQL Architecture
3. Installation with Docker
4. Connecting and Basic Operations
5. PostgreSQL Data Types Beyond the Basics
6. Extensions — PostgreSQL's Superpower
7. Roles and Permissions
8. Connecting from Go with pgx
9. Exercises

---

## 1. What Makes PostgreSQL Special

PostgreSQL stands out from other databases in several ways:

**Standards compliance.** PostgreSQL follows the SQL standard very closely. If you learn SQL on PostgreSQL, it mostly transfers to other databases.

**ACID by default.** Every operation is transactional. You never worry about data being half-written.

**Extensibility.** You can add new data types, operators, index types, and functions. This gave us PostGIS (geospatial), pgvector (vector search), and TimescaleDB (time-series) — all just PostgreSQL extensions.

**JSONB.** PostgreSQL handles JSON documents natively with full indexing. It's essentially a document database built into a relational database.

**Window functions, CTEs, lateral joins.** The most powerful SQL dialect outside of enterprise databases.

**Free and open source.** No licensing fees, no vendor lock-in.

---

## 2. PostgreSQL Architecture

Understanding the architecture helps you understand performance.

```
Client (your Go app)
    │
    │ TCP connection (port 5432)
    ▼
Postmaster process
    │ forks a new backend process per connection
    ▼
Backend process (one per client connection)
    │
    ├── Parser & Planner  → parses SQL, chooses query plan
    ├── Executor          → executes the plan
    │
    ▼
Shared Memory
    ├── Shared Buffer Pool  → caches data pages in RAM (shared_buffers)
    ├── WAL Buffer          → buffers WAL records before flush to disk
    └── Lock Table          → tracks all held locks
    │
    ▼
Disk
    ├── Data files ($PGDATA/base/)  → tables, indexes stored as pages
    ├── WAL files ($PGDATA/pg_wal/) → write-ahead log
    └── pg_log/                      → server logs
```

**Key implication:** PostgreSQL uses **one OS process per connection**. Each connection uses ~5-10 MB of memory. With 1000 connections, that's 5-10 GB just for connection overhead. This is why connection pooling (PgBouncer, pgxpool) is essential for production.

---

## 3. Installation with Docker

The fastest way to run PostgreSQL locally:

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_USER=dev \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  postgres:16

# Connect with psql
docker exec -it postgres psql -U dev -d myapp
```

Or with Docker Compose (better for projects):

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

```bash
docker compose up -d
```

---

## 4. Connecting and Basic Operations

```bash
# Connect with psql
psql postgresql://dev:secret@localhost:5432/myapp

# Useful psql commands
\l          -- list databases
\dt         -- list tables
\d users    -- describe the users table
\timing     -- show query execution time
\q          -- quit
```

Essential SQL for getting oriented:

```sql
-- What version is running?
SELECT version();

-- What databases exist?
SELECT datname FROM pg_database;

-- What tables are in this database?
SELECT tablename FROM pg_tables WHERE schemaname = 'public';

-- How big is each table?
SELECT relname, pg_size_pretty(pg_total_relation_size(relid))
FROM pg_stat_user_tables
ORDER BY pg_total_relation_size(relid) DESC;

-- What queries are currently running?
SELECT pid, query, state, duration
FROM pg_stat_activity
WHERE state != 'idle';
```

---

## 5. PostgreSQL Data Types Beyond the Basics

PostgreSQL has a much richer type system than most databases.

### JSONB — Binary JSON

```sql
CREATE TABLE events (
    id      BIGSERIAL PRIMARY KEY,
    payload JSONB NOT NULL
);

INSERT INTO events (payload) VALUES
    ('{"type": "click", "user_id": 42, "button": "buy"}'),
    ('{"type": "view",  "user_id": 7,  "page": "/home"}');

-- Query inside JSON
SELECT payload->>'type' AS event_type,
       payload->>'user_id' AS user_id
FROM events
WHERE payload->>'type' = 'click';

-- Check if JSON contains a key-value pair
SELECT * FROM events WHERE payload @> '{"type": "click"}';

-- Index JSONB for fast queries
CREATE INDEX idx_events_payload ON events USING gin(payload);
```

### Arrays

```sql
CREATE TABLE articles (
    id   SERIAL PRIMARY KEY,
    tags TEXT[]
);

INSERT INTO articles (tags) VALUES (ARRAY['go', 'database', 'tutorial']);

-- Query articles with a specific tag
SELECT * FROM articles WHERE 'go' = ANY(tags);

-- Aggregate tags across all articles
SELECT UNNEST(tags) as tag, COUNT(*) as count
FROM articles
GROUP BY tag
ORDER BY count DESC;
```

### UUID

```sql
-- Enable the uuid extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sessions (
    id         UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Range Types

```sql
-- Store a time range for a booking
CREATE TABLE bookings (
    id       SERIAL PRIMARY KEY,
    room_id  INT,
    during   TSTZRANGE NOT NULL,
    EXCLUDE USING gist(room_id WITH =, during WITH &&)  -- prevent overlapping bookings!
);

INSERT INTO bookings (room_id, during)
VALUES (101, '[2024-03-01 10:00, 2024-03-01 12:00)');

-- Find all bookings that overlap with a given range
SELECT * FROM bookings
WHERE during && '[2024-03-01 09:00, 2024-03-01 11:00)';
```

---

## 6. Extensions — PostgreSQL's Superpower

Extensions add functionality to PostgreSQL without modifying the core.

```sql
-- Full-text search helpers
CREATE EXTENSION IF NOT EXISTS pg_trgm;      -- trigram similarity
CREATE EXTENSION IF NOT EXISTS unaccent;     -- remove accents from searches

-- Monitoring
CREATE EXTENSION IF NOT EXISTS pg_stat_statements; -- track all query stats

-- Vector search
CREATE EXTENSION IF NOT EXISTS vector;       -- pgvector for AI embeddings

-- Check installed extensions
SELECT name, default_version, installed_version
FROM pg_available_extensions
WHERE installed_version IS NOT NULL;
```

---

## 7. Roles and Permissions

PostgreSQL uses a role system for access control.

```sql
-- Create a read-only role for analytics
CREATE ROLE analytics_reader;
GRANT CONNECT ON DATABASE myapp TO analytics_reader;
GRANT USAGE ON SCHEMA public TO analytics_reader;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO analytics_reader;

-- Create a role for the application
CREATE ROLE app_user LOGIN PASSWORD 'secret';
GRANT CONNECT ON DATABASE myapp TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- Revoke dangerous permissions
REVOKE CREATE ON SCHEMA public FROM PUBLIC;  -- prevent creating tables
```

**Best practice:** Your application should connect as a role with only the permissions it needs — no more. Never connect as `postgres` (the superuser) from application code.

---

## 8. Connecting from Go with pgx

`pgx` is the best Go driver for PostgreSQL. It's faster than `database/sql` and supports PostgreSQL-specific features.

```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool  # connection pooling
```

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
    ID    int
    Email string
    Name  string
}

func main() {
    ctx := context.Background()

    // Create a connection pool (reuses connections, not one per query)
    pool, err := pgxpool.New(ctx, "postgres://dev:secret@localhost:5432/myapp")
    if err != nil {
        log.Fatal("connect:", err)
    }
    defer pool.Close()

    // Ping to verify connection
    if err := pool.Ping(ctx); err != nil {
        log.Fatal("ping:", err)
    }
    fmt.Println("Connected to PostgreSQL!")

    // Insert a user
    var newID int
    err = pool.QueryRow(ctx,
        "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id",
        "alice@example.com", "Alice",
    ).Scan(&newID)
    if err != nil {
        log.Fatal("insert:", err)
    }
    fmt.Println("Created user with ID:", newID)

    // Query users
    rows, err := pool.Query(ctx, "SELECT id, email, name FROM users LIMIT 10")
    if err != nil {
        log.Fatal("query:", err)
    }
    defer rows.Close()

    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Email, &u.Name); err != nil {
            log.Fatal("scan:", err)
        }
        fmt.Printf("User: %+v\n", u)
    }

    if err := rows.Err(); err != nil {
        log.Fatal("rows error:", err)
    }
}
```

### Connection Pool Configuration

```go
config, _ := pgxpool.ParseConfig("postgres://dev:secret@localhost:5432/myapp")
config.MaxConns = 20          // maximum connections in the pool
config.MinConns = 5           // always keep 5 connections warm
config.MaxConnLifetime = time.Hour  // recycle connections after 1 hour
config.MaxConnIdleTime = 30 * time.Minute

pool, err := pgxpool.NewWithConfig(ctx, config)
```

---

## Summary

- PostgreSQL is ACID-compliant, extensible, and the most feature-rich open-source database.
- Its process-per-connection model requires connection pooling (pgxpool) in production.
- JSONB, arrays, UUID, and range types give you much more power than basic `INT/TEXT`.
- Extensions (pg_trgm, pgvector, pg_stat_statements) add powerful capabilities without leaving PostgreSQL.
- Always connect as a least-privilege role, never as the superuser.
- Use `pgxpool` in Go for efficient connection management.

### Exercises

**Easy:** Start PostgreSQL with Docker. Connect with psql. Create a `products` table with `id`, `name`, `price`, and `created_at`. Insert 5 rows and SELECT them all.

**Medium:** Add a JSONB column `metadata` to your products table. Insert products with different metadata shapes. Query for products where the metadata contains a specific key-value pair using the `@>` operator.

**Hard:** Create a PostgreSQL role with SELECT-only access. Connect from Go using that role and verify you can SELECT but not INSERT. What error do you get when trying to INSERT?
