# Chapter 76: Database Migrations with golang-migrate

Every application that talks to a database has a schema — tables, columns, indexes, constraints. That schema changes over time: you add a column, drop an old one, rename a table. Without a migration system, those changes live only in someone's head or a one-off SQL script run on the server. With migrations, every change is a versioned file that lives in git alongside the code that needs it.

## Table of Contents

1. [Why Migrations?](#1-why-migrations)
2. [Migration File Naming Convention](#2-migration-file-naming-convention)
3. [Writing Safe Migrations](#3-writing-safe-migrations)
4. [golang-migrate: CLI and Library](#4-golang-migrate-cli-and-library)
5. [Embedding Migrations in Your Binary](#5-embedding-migrations-in-your-binary)
6. [Running Migrations at Startup](#6-running-migrations-automatically-at-startup)
7. [Best Practices](#7-best-practices)
8. [Squashing Migrations](#8-squashing-migrations)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Migrations?

Imagine two developers, Alice and Bob, both working on the same Go service. Alice adds a `role` column to the `users` table on her laptop. Bob pulls her code, runs it, and gets a panic because the `role` column does not exist in his local database.

Migrations solve this by treating schema changes exactly like code changes: written in a file, committed to git, reviewed in pull requests, and applied consistently everywhere — laptop, staging, production.

```
Without migrations                With migrations
─────────────────                 ───────────────
"Oh, you need to run this SQL     git pull
 manually on your machine"         make migrate-up
                                   ✓ schema is now at version 7
```

The other half of migrations is rollbacks. Every change should be reversible. If version 7 breaks production, you run one command and you are back to version 6.

---

## 2. Migration File Naming Convention

golang-migrate uses a sequential numbering scheme:

```
migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql
  000002_add_posts_table.up.sql
  000002_add_posts_table.down.sql
  000003_add_user_role.up.sql
  000003_add_user_role.down.sql
  000004_add_posts_index.up.sql
  000004_add_posts_index.down.sql
```

Rules:
- The **number prefix** is zero-padded to 6 digits. Always increments — never reuses a number.
- The **name part** is descriptive and uses underscores. Read like a commit message: what does this migration do?
- `.up.sql` contains the **forward** change (apply).
- `.down.sql` contains the **reverse** change (rollback).

```
000001_create_users.up.sql   → CREATE TABLE users (...)
000001_create_users.down.sql → DROP TABLE users
```

The tool tracks which migrations have been applied in a special table it creates automatically (`schema_migrations` by default).

---

## 3. Writing Safe Migrations

### The up migration

Use idempotent SQL so the migration can safely be re-run (or so errors are obvious):

```sql
-- 000001_create_users.up.sql
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    email         TEXT         NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    name          TEXT         NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
```

Adding a column to an existing table:

```sql
-- 000003_add_user_role.up.sql
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'member';

-- Add a check constraint to limit valid values
ALTER TABLE users
    ADD CONSTRAINT IF NOT EXISTS chk_users_role
    CHECK (role IN ('member', 'admin', 'moderator'));
```

Creating a table for posts:

```sql
-- 000002_add_posts_table.up.sql
CREATE TABLE IF NOT EXISTS posts (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      TEXT         NOT NULL,
    slug       TEXT         NOT NULL UNIQUE,
    body       TEXT         NOT NULL DEFAULT '',
    published  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id  ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_posts_slug     ON posts (slug);
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts (published) WHERE published = TRUE;
```

### The down migration

The down migration must cleanly undo what the up migration did:

```sql
-- 000001_create_users.down.sql
DROP TABLE IF EXISTS users;

-- 000002_add_posts_table.down.sql
DROP TABLE IF EXISTS posts;

-- 000003_add_user_role.down.sql
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role;
ALTER TABLE users DROP COLUMN  IF EXISTS role;
```

### What to avoid

```
BAD                                  GOOD
──────────────────────────           ─────────────────────────────────────
DROP COLUMN email                    Deprecation period: rename to
                                     email_deprecated, deploy code that
                                     no longer reads it, then DROP in a
                                     later migration

RENAME COLUMN name TO full_name      Same: add full_name, copy data,
                                     update code, drop name later

DROP TABLE users                     Only after all references removed
                                     and data migrated elsewhere

ALTER TABLE ADD COLUMN NOT NULL      Add as nullable first, backfill,
  without a DEFAULT                  then add the NOT NULL constraint
```

### Large table migrations

Adding an index on a large table locks writes. Use `CONCURRENTLY`:

```sql
-- 000004_add_posts_index.up.sql
-- NOTE: CREATE INDEX CONCURRENTLY cannot run inside a transaction.
-- golang-migrate wraps each migration in a transaction by default.
-- Add this comment to disable the transaction wrapper:

-- migrate: no-transaction

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_title
    ON posts USING gin(to_tsvector('english', title));
```

```sql
-- 000004_add_posts_index.down.sql
-- migrate: no-transaction

DROP INDEX CONCURRENTLY IF EXISTS idx_posts_title;
```

---

## 4. golang-migrate: CLI and Library

### Install the CLI

```bash
# macOS
brew install golang-migrate

# Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify
migrate -version
```

### Create a new migration file

```bash
migrate create -ext sql -dir ./migrations -seq add_user_avatar
# Creates:
#   migrations/000005_add_user_avatar.up.sql
#   migrations/000005_add_user_avatar.down.sql
```

### Run migrations

```bash
# Apply all pending migrations
migrate -path ./migrations -database "postgres://user:pass@localhost/mydb?sslmode=disable" up

# Roll back the last applied migration
migrate -path ./migrations -database "postgres://user:pass@localhost/mydb?sslmode=disable" down 1

# Roll back all applied migrations
migrate -path ./migrations -database "postgres://user:pass@localhost/mydb?sslmode=disable" down

# Check the current version
migrate -path ./migrations -database "postgres://user:pass@localhost/mydb?sslmode=disable" version

# Jump to a specific version (use with care)
migrate -path ./migrations -database "postgres://user:pass@localhost/mydb?sslmode=disable" goto 3
```

### Use the library directly in Go

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dsn string) error {
    m, err := migrate.New(
        "file://migrations",                  // source: local directory
        dsn,                                   // target: postgres connection string
    )
    if err != nil {
        return fmt.Errorf("create migrator: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        return fmt.Errorf("run migrations: %w", err)
    }
    return nil
}
```

`migrate.ErrNoChange` is returned when all migrations are already applied — that is not an error; the database is up to date.

---

## 5. Embedding Migrations in Your Binary

Shipping the migration SQL files alongside the binary as separate files is fragile — the files can go missing. Embed them directly into the binary using Go's `embed` package:

```go
// migrations/embed.go
package migrations

import "embed"

// FS holds all migration files, compiled into the binary.
//go:embed *.sql
var FS embed.FS
```

Then reference the embedded filesystem instead of the local directory:

```go
import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/source/iofs"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"

    "myapp/migrations"
)

func runMigrations(dsn string) error {
    // iofs.New wraps an embed.FS as a migrate source
    src, err := iofs.New(migrations.FS, ".")
    if err != nil {
        return fmt.Errorf("create migration source: %w", err)
    }

    m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
    if err != nil {
        return fmt.Errorf("create migrator: %w", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        return fmt.Errorf("run migrations: %w", err)
    }
    return nil
}
```

Your directory structure:

```
myapp/
  migrations/
    embed.go           ← the //go:embed declaration
    000001_create_users.up.sql
    000001_create_users.down.sql
    000002_add_posts_table.up.sql
    000002_add_posts_table.down.sql
  main.go
```

Now `go build` compiles the SQL files into the binary. No separate deployment of SQL files needed.

---

## 6. Running Migrations Automatically at Startup

Many services run pending migrations automatically when they boot. This keeps deployment simple: push code → restart service → schema is up to date.

```go
// internal/database/database.go
package database

import (
    "context"
    "errors"
    "fmt"
    "log/slog"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/source/iofs"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "myapp/migrations"
)

type DB struct {
    *sqlx.DB
}

func Open(dsn string) (*DB, error) {
    db, err := sqlx.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("open db: %w", err)
    }

    if err := db.PingContext(context.Background()); err != nil {
        return nil, fmt.Errorf("ping db: %w", err)
    }

    return &DB{db}, nil
}

func (db *DB) MigrateUp() error {
    src, err := iofs.New(migrations.FS, ".")
    if err != nil {
        return fmt.Errorf("migration source: %w", err)
    }

    m, err := migrate.NewWithSourceInstance("iofs", src, db.DriverName()+"://"+db.Stats().MaxOpenConnections)
    if err != nil {
        return fmt.Errorf("migrator: %w", err)
    }
    defer m.Close()

    version, dirty, _ := m.Version()
    slog.Info("running migrations", "current_version", version, "dirty", dirty)

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            slog.Info("migrations: no changes, schema is up to date")
            return nil
        }
        return fmt.Errorf("migrate up: %w", err)
    }

    version, _, _ = m.Version()
    slog.Info("migrations complete", "version", version)
    return nil
}

// main.go
func main() {
    db, err := database.Open(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }

    if err := db.MigrateUp(); err != nil {
        log.Fatalf("migrations failed: %v", err)
    }

    // Start the server
    srv := server.New(db)
    srv.ListenAndServe(":8080")
}
```

### When auto-migration at startup is a bad idea

It is fine for small teams and simple schemas. For large production databases:
- Multiple instances starting at the same time may race to apply the same migration.
- Long-running migrations (adding an index on 100M rows) block startup for minutes.

The fix for races is that golang-migrate uses an advisory lock — only one instance applies migrations, the others wait. The fix for long migrations is to run them as a separate deployment step before restarting the service.

---

## 7. Best Practices

**Never edit a committed migration.** Once a migration file is merged to main, treat it as immutable. Other developers and all environments have already applied it. If you got something wrong, write a new migration that fixes it.

```
WRONG: edit 000003_add_user_role.up.sql after it has been applied
RIGHT: create 000006_fix_user_role_constraint.up.sql
```

**Always write the down migration.** It is tempting to leave `down.sql` empty when you are moving fast. You will regret it the first time a bad deploy needs a quick rollback.

**Test the rollback.** Add this to your CI pipeline:

```bash
# Apply all migrations
migrate -path ./migrations -database "$TEST_DSN" up

# Roll back one
migrate -path ./migrations -database "$TEST_DSN" down 1

# Apply it again — should work cleanly
migrate -path ./migrations -database "$TEST_DSN" up
```

**Use transactions in migrations.** golang-migrate wraps each migration in a transaction by default. If a migration fails halfway, the partial changes are rolled back. Only opt out (`-- migrate: no-transaction`) when you must (e.g., `CREATE INDEX CONCURRENTLY`).

**Back up before big migrations.** Before running a migration that alters a large table in production, take a snapshot. It is faster to restore a snapshot than to write a recovery migration under pressure.

---

## 8. Squashing Migrations

After months of development, you might have 200 migration files. Bootstrapping a new environment replays all 200 in order, which can take minutes. Squashing condenses them into one "baseline" migration:

```
Before squash                 After squash
─────────────                 ────────────
000001_create_users           000000_baseline.up.sql   ← full schema snapshot
000002_add_posts              000000_baseline.down.sql ← DROP everything
000003_add_role
...
000199_fix_index
```

The workflow:

1. Dump the current schema: `pg_dump --schema-only mydb > migrations/000000_baseline.up.sql`
2. Add `-- migrate: no-transaction` if needed.
3. Write the down migration (`DROP TABLE ...` for every table in reverse dependency order).
4. Delete migrations 000001 through 000199.
5. In production (where the schema already exists), force the version to 0: `migrate ... force 0`. The baseline migration is skipped because the schema already exists.
6. New environments apply only the single baseline file.

A simpler alternative: keep old migrations but add a `SKIP_OLD_MIGRATIONS=true` environment variable that marks them as applied without running them.

---

## Summary

- Migrations are versioned SQL files committed alongside your code — `.up.sql` applies, `.down.sql` reverses
- Naming: `000001_describe_what_it_does.up.sql` — sequential, zero-padded, descriptive
- Write safe SQL: `CREATE TABLE IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`, never edit a committed file
- golang-migrate CLI: `migrate up/down/version/goto`; library: `migrate.New(source, dsn).Up()`
- Embed migrations with `//go:embed *.sql` and `iofs.New` so the binary is self-contained
- Auto-migrate at startup is fine for small teams; for large production databases, run as a separate deploy step
- Always write the down migration and test rollback in CI
- Squash 200+ migrations into a baseline to keep new environment setup fast

## Exercises

### Easy
1. Create a `migrations/` directory with three migration pairs: `000001_create_users`, `000002_create_posts`, `000003_add_published_column`. Run them against a local PostgreSQL instance with the `migrate` CLI. Roll back one step, then re-apply.
2. Add `//go:embed *.sql` to a `migrations/embed.go` file and write a `RunMigrations(dsn string) error` function using `iofs.New`. Call it from `main()` and verify it prints "no changes" on a second run.
3. Write a bash script that runs `migrate up`, then `migrate down 1`, then `migrate up` again. Run it in CI against a test database. Assert the exit code is 0 at each step.

### Medium
4. Add a `version` command to your service: `GET /admin/migrations/version` returns the current migration version and whether the database is dirty. Use `m.Version()` from the library.
5. Write migration 000004 that adds a `posts.view_count BIGINT DEFAULT 0` column. Write migration 000005 that creates a partial index on `posts WHERE published = TRUE`. Verify that `migrate down 2` removes both, and `migrate up` reapplies them.
6. Implement a `MigrateWithTimeout(dsn string, timeout time.Duration) error` function that wraps migration in a context with deadline. If the migration takes longer than `timeout`, it returns an error. Test with a migration that intentionally sleeps (use `pg_sleep`).

### Hard
7. Build a **migration smoke test harness**: for each migration from 1 to N, apply migrations 1..i, run a Go function that validates the schema at that version (check that expected tables and columns exist using `information_schema`), then roll back to 0. This catches regressions where a down migration doesn't clean up properly.
8. Implement the **expand-contract pattern** for a zero-downtime column rename: migration A adds `full_name` alongside the existing `name`, migration B (deployed after all app instances use `full_name`) drops `name`. Write both migrations, write the Go repository code that reads from `full_name` if present and falls back to `name`, and write a test that verifies the service works correctly at both schema versions.
