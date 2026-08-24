# Chapter 51: PostgreSQL and sqlx

The in-memory stores from previous chapters can't survive a restart. Real services need a database. PostgreSQL is the gold standard for relational data: ACID transactions, complex queries, and decades of battle-testing. This chapter connects Go to PostgreSQL using `database/sql` and `sqlx`.

## Table of Contents

1. [database/sql Architecture](#1-databasesql-architecture)
2. [Connecting with sqlx](#2-connecting-with-sqlx)
3. [Schema and Migrations](#3-schema-and-migrations)
4. [CRUD with sqlx](#4-crud-with-sqlx)
5. [Transactions](#5-transactions)
6. [Connection Pool Tuning](#6-connection-pool-tuning)
7. [Query Patterns and Pitfalls](#7-query-patterns-and-pitfalls)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. database/sql Architecture

```
Application
    ↓
database/sql (standard library)
    ↓
Driver (lib/pq or pgx)
    ↓
PostgreSQL TCP connection
```

`database/sql` manages a **connection pool** — you rarely work with individual connections. The pool:
- Keeps `MaxOpenConns` connections alive
- Reuses idle connections (saves TCP handshake + auth overhead)
- Opens new connections on demand up to the max
- Validates idle connections with `Ping` before returning them

```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // Side-effect import registers the "postgres" driver
)

db, err := sql.Open("postgres", "postgres://user:pass@host/db?sslmode=require")
```

---

## 2. Connecting with sqlx

`sqlx` extends `database/sql` with named queries, struct scanning, and in-query slices.

```bash
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
# or pgx (faster, more features):
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/stdlib
```

```go
package db

import (
    "context"
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type Config struct {
    URL             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func Connect(cfg Config) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", cfg.URL)
    if err != nil {
        return nil, fmt.Errorf("connect postgres: %w", err)
    }

    db.SetMaxOpenConns(cfg.MaxOpenConns)
    db.SetMaxIdleConns(cfg.MaxIdleConns)
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

    // Verify connectivity:
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return db, nil
}

// Sensible defaults for most services:
var DefaultConfig = Config{
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 30 * time.Minute,
    ConnMaxIdleTime: 5 * time.Minute,
}
```

---

## 3. Schema and Migrations

Never modify production database schema manually. Use migrations: versioned, ordered SQL files.

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir migrations -seq create_users
```

```sql
-- migrations/000001_create_users.up.sql
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

-- migrations/000002_create_notes.up.sql
CREATE TABLE notes (
    id         BIGSERIAL PRIMARY KEY,
    owner_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(500) NOT NULL,
    content    TEXT NOT NULL DEFAULT '',
    status     VARCHAR(20) NOT NULL DEFAULT 'draft'
                   CHECK (status IN ('draft', 'published', 'archived')),
    version    INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notes_owner_id ON notes(owner_id);
CREATE INDEX idx_notes_updated_at ON notes(updated_at DESC);

-- Full-text search index:
ALTER TABLE notes ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(content, '')), 'B')
    ) STORED;
CREATE INDEX idx_notes_fts ON notes USING GIN(search_vector);
```

```sql
-- migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS users;

-- migrations/000002_create_notes.down.sql
DROP TABLE IF EXISTS notes;
```

```go
// Run migrations at startup:
import "github.com/golang-migrate/migrate/v4"
import _ "github.com/golang-migrate/migrate/v4/database/postgres"
import _ "github.com/golang-migrate/migrate/v4/source/file"

func RunMigrations(dbURL, migrationsDir string) error {
    m, err := migrate.New("file://"+migrationsDir, dbURL)
    if err != nil { return fmt.Errorf("new migrate: %w", err) }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("run migrations: %w", err)
    }
    return nil
}
```

---

## 4. CRUD with sqlx

```go
package store

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")
var ErrEmailTaken = errors.New("email already taken")

type User struct {
    ID           int64     `db:"id"`
    Name         string    `db:"name"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
}

type UserStore struct {
    db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
    return &UserStore{db: db}
}

// Create inserts a new user and returns it with the DB-assigned ID.
func (s *UserStore) Create(ctx context.Context, name, email, passwordHash string) (*User, error) {
    const q = `
        INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, password_hash, created_at, updated_at`

    var u User
    err := s.db.QueryRowxContext(ctx, q, name, email, passwordHash).StructScan(&u)
    if err != nil {
        // Detect duplicate email (Postgres error code 23505):
        var pqErr *pq.Error
        if errors.As(err, &pqErr) && pqErr.Code == "23505" {
            return nil, ErrEmailTaken
        }
        return nil, fmt.Errorf("insert user: %w", err)
    }
    return &u, nil
}

// GetByEmail fetches a user by email for authentication.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
    const q = `SELECT * FROM users WHERE email = $1`
    var u User
    err := s.db.GetContext(ctx, &u, q, email)
    if errors.Is(err, sql.ErrNoRows) { return nil, ErrNotFound }
    if err != nil { return nil, fmt.Errorf("get user by email: %w", err) }
    return &u, nil
}

// GetByID fetches a user by ID.
func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
    const q = `SELECT * FROM users WHERE id = $1`
    var u User
    err := s.db.GetContext(ctx, &u, q, id)
    if errors.Is(err, sql.ErrNoRows) { return nil, ErrNotFound }
    if err != nil { return nil, fmt.Errorf("get user: %w", err) }
    return &u, nil
}

// Note store:
type Note struct {
    ID        int64     `db:"id"`
    OwnerID   int64     `db:"owner_id"`
    Title     string    `db:"title"`
    Content   string    `db:"content"`
    Status    string    `db:"status"`
    Version   int       `db:"version"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type NoteStore struct {
    db *sqlx.DB
}

func (s *NoteStore) Create(ctx context.Context, ownerID int64, title, content string) (*Note, error) {
    const q = `
        INSERT INTO notes (owner_id, title, content)
        VALUES ($1, $2, $3)
        RETURNING *`
    var n Note
    if err := s.db.QueryRowxContext(ctx, q, ownerID, title, content).StructScan(&n); err != nil {
        return nil, fmt.Errorf("create note: %w", err)
    }
    return &n, nil
}

func (s *NoteStore) GetByID(ctx context.Context, id int64) (*Note, error) {
    const q = `SELECT * FROM notes WHERE id = $1`
    var n Note
    if err := s.db.GetContext(ctx, &n, q, id); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, ErrNotFound }
        return nil, fmt.Errorf("get note: %w", err)
    }
    return &n, nil
}

// Update uses optimistic locking via the version column:
func (s *NoteStore) Update(ctx context.Context, id, ownerID int64, title, content string, version int) (*Note, error) {
    const q = `
        UPDATE notes
        SET title = COALESCE(NULLIF($1, ''), title),
            content = COALESCE(NULLIF($2, ''), content),
            version = version + 1,
            updated_at = NOW()
        WHERE id = $3
          AND owner_id = $4
          AND version = $5
        RETURNING *`

    var n Note
    err := s.db.QueryRowxContext(ctx, q, title, content, id, ownerID, version).StructScan(&n)
    if errors.Is(err, sql.ErrNoRows) {
        // Either not found, not owner, or version conflict — check which:
        var exists bool
        s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1)`, id).Scan(&exists)
        if !exists { return nil, ErrNotFound }
        return nil, ErrVersionConflict
    }
    if err != nil { return nil, fmt.Errorf("update note: %w", err) }
    return &n, nil
}

// List notes with pagination:
func (s *NoteStore) ListByOwner(ctx context.Context, ownerID int64, offset, limit int) ([]*Note, int, error) {
    // Get total count:
    var total int
    if err := s.db.QueryRowContext(ctx,
        `SELECT COUNT(*) FROM notes WHERE owner_id = $1`, ownerID,
    ).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("count notes: %w", err)
    }

    // Get page:
    const q = `
        SELECT * FROM notes
        WHERE owner_id = $1
        ORDER BY updated_at DESC
        LIMIT $2 OFFSET $3`

    var notes []*Note
    if err := s.db.SelectContext(ctx, &notes, q, ownerID, limit, offset); err != nil {
        return nil, 0, fmt.Errorf("list notes: %w", err)
    }
    return notes, total, nil
}

// Full-text search:
func (s *NoteStore) Search(ctx context.Context, ownerID int64, query string) ([]*Note, error) {
    const q = `
        SELECT *, ts_rank(search_vector, plainto_tsquery('english', $2)) AS rank
        FROM notes
        WHERE owner_id = $1
          AND search_vector @@ plainto_tsquery('english', $2)
        ORDER BY rank DESC
        LIMIT 20`

    var notes []*Note
    if err := s.db.SelectContext(ctx, &notes, q, ownerID, query); err != nil {
        return nil, fmt.Errorf("search notes: %w", err)
    }
    return notes, nil
}
```

---

## 5. Transactions

```go
// Atomic operation: create a note and a tag in one transaction.
func (s *NoteStore) CreateWithTags(ctx context.Context, ownerID int64, title, content string, tagNames []string) (*Note, error) {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil { return nil, fmt.Errorf("begin tx: %w", err) }
    defer tx.Rollback()  // No-op if already committed

    // Create note:
    var n Note
    if err := tx.QueryRowxContext(ctx,
        `INSERT INTO notes (owner_id, title, content) VALUES ($1, $2, $3) RETURNING *`,
        ownerID, title, content,
    ).StructScan(&n); err != nil {
        return nil, fmt.Errorf("create note in tx: %w", err)
    }

    // Create tags and associations:
    for _, name := range tagNames {
        var tagID int64
        err := tx.QueryRowContext(ctx,
            `INSERT INTO tags (name, owner_id) VALUES ($1, $2)
             ON CONFLICT (name, owner_id) DO UPDATE SET name = EXCLUDED.name
             RETURNING id`,
            name, ownerID,
        ).Scan(&tagID)
        if err != nil { return nil, fmt.Errorf("upsert tag %q: %w", name, err) }

        if _, err := tx.ExecContext(ctx,
            `INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
            n.ID, tagID,
        ); err != nil {
            return nil, fmt.Errorf("link tag: %w", err)
        }
    }

    if err := tx.Commit(); err != nil { return nil, fmt.Errorf("commit: %w", err) }
    return &n, nil
}

// Helper for running any function inside a transaction with automatic rollback:
func WithTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()
    if err := fn(tx); err != nil { return err }
    return tx.Commit()
}
```

---

## 6. Connection Pool Tuning

```go
// Rules of thumb:
// MaxOpenConns: CPU cores * 4 for OLTP, or test with benchmarks
// MaxIdleConns: MaxOpenConns / 5 (don't hold too many idle connections — each uses memory)
// ConnMaxLifetime: 30min (force reconnect to pick up DB config changes / load balancer rotation)
// ConnMaxIdleTime: 5min (close idle connections sooner to reduce DB-side load)

// Expose pool metrics (Chapter 101 covers Prometheus):
func PoolStats(db *sqlx.DB) sql.DBStats {
    return db.Stats()
    // Stats.OpenConnections — currently open
    // Stats.InUse          — currently in use by queries
    // Stats.Idle           — waiting to be used
    // Stats.WaitCount      — queries that waited for a connection
    // Stats.WaitDuration   — total time queries waited
}

// Alert if WaitDuration is increasing — means pool is too small or queries are too slow.
```

---

## 7. Query Patterns and Pitfalls

```go
// Pitfall 1: N+1 queries — loading associations one by one:
// BAD:
for _, note := range notes {
    note.Tags, _ = store.GetTagsForNote(ctx, note.ID)  // N queries!
}

// GOOD: fetch all in one query with JOIN or IN:
const q = `
    SELECT n.*, t.id as tag_id, t.name as tag_name
    FROM notes n
    LEFT JOIN note_tags nt ON nt.note_id = n.id
    LEFT JOIN tags t ON t.id = nt.tag_id
    WHERE n.owner_id = $1`
// Then group by note ID in Go

// Pitfall 2: sqlx IN query — use sqlx.In for dynamic slices:
ids := []int64{1, 2, 3}
query, args, err := sqlx.In(`SELECT * FROM notes WHERE id IN (?)`, ids)
if err != nil { ... }
query = db.Rebind(query)  // Convert ? to $1, $2, $3 for postgres
db.SelectContext(ctx, &notes, query, args...)

// Pitfall 3: Scanning into wrong types — always check db tags match column names:
type Note struct {
    ID      int64  `db:"id"`         // Must match column name exactly
    OwnerID int64  `db:"owner_id"`   // Snake_case matches Postgres convention
}

// Pitfall 4: Forgetting to close rows:
rows, _ := db.QueryxContext(ctx, query, args...)
defer rows.Close()  // Always!
for rows.Next() {
    var n Note
    rows.StructScan(&n)
    notes = append(notes, n)
}
if err := rows.Err(); err != nil { ... }  // Check for iteration errors

// Pitfall 5: NULL columns need pointer types or sql.NullXxx:
type User struct {
    AvatarURL *string        `db:"avatar_url"`      // nil if NULL
    DeletedAt sql.NullTime   `db:"deleted_at"`       // .Valid == false if NULL
}
```

---

## Summary

- `database/sql` is the standard interface; `sqlx` adds `StructScan`, `GetContext`, `SelectContext`, named queries
- Use `RETURNING *` to get the inserted/updated row in one round trip
- Detect PostgreSQL unique constraint violations by checking `pq.Error.Code == "23505"`
- Transactions: always `defer tx.Rollback()` + explicit `tx.Commit()` at the end
- Optimistic locking: include `WHERE version = $N` in UPDATE; zero rows affected = conflict
- Built-in full-text search via `tsvector` + `plainto_tsquery` + GIN index
- Pool: 25 max open, 5 max idle, 30min lifetime for most OLTP workloads
- N+1: fetch related data with JOIN or `IN` — never with per-row queries

---

## Exercises

### Easy
1. Write `store.UserStore.Update(ctx, id, name string) (*User, error)` that updates the user's name and `updated_at`. Use `RETURNING *`. Write a test using `dockertest` or a real local PostgreSQL to verify the update.
2. Implement soft deletes: add `deleted_at TIMESTAMPTZ` column to `notes`. `Delete()` sets `deleted_at = NOW()` instead of deleting the row. All query methods filter `WHERE deleted_at IS NULL`. Add `Restore(ctx, id)` to undelete.
3. Write a query that returns the 10 most recently updated notes for a user, along with a count of how many tags each note has. Use a `LEFT JOIN` + `GROUP BY` + `COUNT`. Map the result to a struct with a `TagCount int` field.

### Medium
4. Implement **pagination with cursor** instead of LIMIT/OFFSET. The cursor is `(updated_at, id)` from the last item. Next page: `WHERE (updated_at, id) < ($cursor_time, $cursor_id) ORDER BY updated_at DESC, id DESC LIMIT $n`. This is stable even when items are added/removed between pages (unlike OFFSET).
5. Write a `BulkCreate(ctx, notes []Note) ([]Note, error)` that inserts multiple notes in a single `INSERT INTO notes (...) VALUES ($1,$2,$3), ($4,$5,$6), ...`. Use `sqlx.Named` or build the query manually. Benchmark it against N individual inserts for N=100, 1000.
6. Build a **repository pattern**: define `NoteRepository` interface, implement `PostgresNoteRepository`, and `MockNoteRepository`. Tests use the mock; integration tests use Postgres. Show that the handler tests don't import `sqlx` at all — only the repository interface.

### Hard
7. Implement **row-level locking** for the collaborative editing lock feature: use `SELECT ... FOR UPDATE SKIP LOCKED` to take an advisory lock on a note row. If another connection holds the lock, return "note is being edited" immediately. The lock is released when the transaction ends (30s timeout). Test with two goroutines racing to lock the same note.
8. Write a **database migration test**: test that each migration is: (1) idempotent (running `up` twice doesn't fail or corrupt), (2) reversible (running `down` leaves the DB in the same state as before `up`). Test all migrations from 1 to N, verifying the schema after each step. Use a Dockerized PostgreSQL via `testcontainers-go`.
