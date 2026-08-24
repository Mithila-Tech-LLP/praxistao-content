# Chapter 12: Schema Design and Normalization

A database schema is like the blueprints of a building. Good blueprints make construction straightforward. Bad blueprints cause expensive problems later. This chapter teaches you how to design schemas that grow with your application without falling apart.

## Table of Contents

1. What Is a Schema?
2. The Three Relationships
3. Normalization — Organizing Data Cleanly
4. Choosing the Right Data Types
5. Timestamps and the Timezone Trap
6. Migrations — Changing a Live Database Safely
7. Major Project: Social Media App Schema
8. Exercises

---

## 1. What Is a Schema?

A schema is the structure of your database: which tables exist, what columns they have, and what relationships connect them. Before writing a single INSERT, you design the schema.

```sql
-- A schema defines structure, not data
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    username   TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE posts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

Good schema design means:
- Data is stored in one place, not duplicated
- Relationships are explicit (foreign keys)
- Constraints prevent invalid data

---

## 2. The Three Relationships

Almost all data relationships are one of three types:

### One-to-Many (Most Common)

One user has many posts. One post belongs to one user.

```sql
-- The "many" side holds the foreign key
CREATE TABLE posts (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),  -- points to ONE user
    content TEXT
);
```

The "many" table (posts) stores the foreign key. One user row ↔ many post rows.

### Many-to-Many

A post can have many tags. A tag can apply to many posts.

```sql
-- You need a junction table
CREATE TABLE tags (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE post_tags (
    post_id BIGINT REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INT    REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)  -- composite primary key prevents duplicates
);
```

The junction table (`post_tags`) holds both foreign keys. To find all tags for post 5:
```sql
SELECT tags.name
FROM tags
JOIN post_tags ON tags.id = post_tags.tag_id
WHERE post_tags.post_id = 5;
```

### One-to-One (Rare)

One user has exactly one profile. Used to split a table that would be too wide.

```sql
CREATE TABLE user_profiles (
    user_id    BIGINT PRIMARY KEY REFERENCES users(id),
    bio        TEXT,
    avatar_url TEXT
);
```

---

## 3. Normalization — Organizing Data Cleanly

Normalization is the process of organizing your schema to eliminate redundancy. Each piece of data should exist in exactly one place.

### Without Normalization (Bad)

```sql
-- Storing city in every user row
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    name        TEXT,
    city_name   TEXT,
    city_country TEXT,
    city_pop    INT   -- duplicated for every user in that city!
);
```

Problems:
- "New York" is spelled differently by different rows ("New York" vs "new york")
- Updating city population requires updating thousands of rows
- Deleting all users from a city also deletes the city information

### With Normalization (Good)

```sql
CREATE TABLE cities (
    id       SERIAL PRIMARY KEY,
    name     TEXT NOT NULL,
    country  TEXT NOT NULL,
    pop      INT
);

CREATE TABLE users (
    id      SERIAL PRIMARY KEY,
    name    TEXT,
    city_id INT REFERENCES cities(id)  -- just a pointer
);
```

Now "New York" data lives in one row. Updates go to one place. This is **Third Normal Form (3NF)**.

### When to Denormalize

Normalization is great for writes. But heavily normalized data can require many JOINs for reads, which is slow for analytics.

For **analytics and reporting**, it's often better to denormalize: precompute totals, store redundant data, accept some duplication in exchange for query speed.

```sql
-- Denormalized analytics table (for speed, not for transactional writes)
CREATE TABLE daily_stats (
    date         DATE PRIMARY KEY,
    new_users    INT,
    total_posts  INT,
    revenue      DECIMAL(12,2)
);
```

The rule: **normalize your transactional data, denormalize your analytical data**.

---

## 4. Choosing the Right Data Types

| Use case | Best type | Avoid |
|----------|-----------|-------|
| User IDs, row IDs | `BIGSERIAL` or `UUID` | `INT` (can run out at 2 billion) |
| Prices, money | `NUMERIC(12,2)` | `FLOAT` (floating point errors!) |
| Short text (<255 chars) | `VARCHAR(255)` or `TEXT` | `CHAR(255)` (wastes space) |
| Long text | `TEXT` | `VARCHAR(10000)` (arbitrary limit) |
| True/false | `BOOLEAN` | `INT` with 0/1 |
| Timestamps | `TIMESTAMPTZ` | `TIMESTAMP` (loses timezone info) |
| Small numbers (0-255) | `SMALLINT` | `BIGINT` (wastes 6 bytes per row) |
| JSON documents | `JSONB` | `TEXT` (can't query JSON efficiently) |
| Network addresses | `INET` | `TEXT` (can't do IP range queries) |

### Never Use FLOAT for Money

```go
// BUG: floating point imprecision
price := 9.99
tax := 0.08
total := price * (1 + tax) // = 10.790000000000001 (not 10.79!)
```

Use `NUMERIC(precision, scale)` in PostgreSQL and `decimal.Decimal` in Go:

```go
import "github.com/shopspring/decimal"

price, _ := decimal.NewFromString("9.99")
tax, _ := decimal.NewFromString("0.08")
total := price.Mul(decimal.NewFromString("1").Add(tax))
fmt.Println(total) // 10.7892 (exact)
```

### UUID vs BIGSERIAL

```sql
-- BIGSERIAL: auto-incrementing integer (fast, sequential, leaks row count)
id BIGSERIAL PRIMARY KEY

-- UUID: random 128-bit identifier (slower inserts, no leaks, globally unique)
id UUID DEFAULT gen_random_uuid() PRIMARY KEY
```

Use UUIDs when:
- You need globally unique IDs across multiple databases
- You don't want to reveal how many rows you have
- You generate IDs before inserting (e.g., in application code)

Use BIGSERIAL when:
- Maximum insert speed matters
- The ID is internal (not exposed in URLs)

---

## 5. Timestamps and the Timezone Trap

**Always use `TIMESTAMPTZ` (timestamp with time zone), never `TIMESTAMP`.**

```sql
-- BAD: stores time without timezone info
created_at TIMESTAMP DEFAULT NOW()

-- GOOD: stores time in UTC, displays in any timezone
created_at TIMESTAMPTZ DEFAULT NOW()
```

`TIMESTAMPTZ` stores everything in UTC internally. When you read it, PostgreSQL converts to your session's timezone. This means:
- No bugs when your server changes timezone (daylight saving, deployment moves)
- Correct comparisons across timezones

In Go:

```go
// When reading from PostgreSQL, always use time.Time
type User struct {
    ID        int
    CreatedAt time.Time // maps to TIMESTAMPTZ
}

// When inserting, time.Time is automatically handled
db.Exec(
    "INSERT INTO users (email, created_at) VALUES ($1, $2)",
    "alice@example.com",
    time.Now().UTC(), // always store UTC
)
```

---

## 6. Migrations — Changing a Live Database Safely

Your schema will change over time. Adding columns, renaming tables, dropping old fields. You can't just edit `CREATE TABLE` — the table already exists with data.

**Migrations** are SQL files that describe schema changes in order.

```
migrations/
  0001_create_users.sql
  0002_add_phone_to_users.sql
  0003_create_posts.sql
  0004_add_post_index.sql
```

Each migration runs once, in order. A `schema_migrations` table tracks which have run.

### Using golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

```go
package main

import (
    "log"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(databaseURL string) {
    m, err := migrate.New("file://migrations", databaseURL)
    if err != nil {
        log.Fatal("creating migrator:", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("running migrations:", err)
    }
    log.Println("migrations applied")
}
```

### Safe Migration Practices

**Adding a column:** Safe. Existing rows get a NULL or DEFAULT.
```sql
ALTER TABLE users ADD COLUMN phone TEXT;
```

**Dropping a column:** Dangerous on live systems! Mark it as unused first, then drop in a later migration.
```sql
-- Migration 1: Stop writing to the column
-- Migration 2 (weeks later): Drop it
ALTER TABLE users DROP COLUMN old_phone;
```

**Renaming a column:** Dangerous. Instead: add new column, backfill, drop old column.

**Adding a NOT NULL column without a default:** Dangerous on large tables (locks the table). Always add with a DEFAULT first, then remove the default later.

```sql
-- Safe way to add NOT NULL column to large table:
-- Step 1: add with default
ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'active' NOT NULL;
-- Step 2 (later): remove the default if not needed
ALTER TABLE users ALTER COLUMN status DROP DEFAULT;
```

---

## 7. Major Project: Social Media App Schema

Design a complete schema for a minimal social media app with users, posts, follows, likes, and comments.

```sql
-- Users
CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL CHECK (length(username) BETWEEN 3 AND 30),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio           TEXT,
    avatar_url    TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ  -- soft delete: NULL = active
);

-- Posts
CREATE TABLE posts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL CHECK (length(content) BETWEEN 1 AND 280),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Follows (many-to-many: users follow users)
CREATE TABLE follows (
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id),
    CHECK (follower_id != followee_id)  -- can't follow yourself
);

-- Likes (many-to-many: users like posts)
CREATE TABLE likes (
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)  -- one like per user per post
);

-- Comments
CREATE TABLE comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL CHECK (length(content) BETWEEN 1 AND 500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_posts_user_id ON posts(user_id, created_at DESC);
CREATE INDEX idx_follows_followee ON follows(followee_id);
CREATE INDEX idx_likes_post ON likes(post_id);
CREATE INDEX idx_comments_post ON comments(post_id, created_at ASC);
-- Don't index soft-deleted rows
CREATE INDEX idx_posts_active ON posts(user_id, created_at DESC) WHERE deleted_at IS NULL;
```

Query examples:

```sql
-- Get a user's timeline (posts from people they follow)
SELECT posts.*, users.username
FROM posts
JOIN follows ON follows.followee_id = posts.user_id
JOIN users ON users.id = posts.user_id
WHERE follows.follower_id = $1
  AND posts.deleted_at IS NULL
ORDER BY posts.created_at DESC
LIMIT 20;

-- Get like count for each post
SELECT post_id, COUNT(*) as like_count
FROM likes
GROUP BY post_id;
```

---

## Summary

- **Schema design** defines tables, columns, types, and relationships before data is stored.
- Three relationship types: **one-to-many** (FK on "many" side), **many-to-many** (junction table), **one-to-one** (FK + PRIMARY KEY).
- **Normalization** stores each fact once. Eliminates update anomalies and data inconsistency.
- Use `NUMERIC` for money, `TIMESTAMPTZ` for timestamps, `BIGSERIAL` or `UUID` for IDs.
- **Migrations** track schema changes in version-ordered SQL files. Tools like `golang-migrate` automate applying them.

### Exercises

**Easy:** Design a schema for a library (books, authors, borrowers, loans). Identify the relationship type (one-to-many or many-to-many) for each pair.

**Medium:** Write a migration to add a `profile_picture_url` column to an existing `users` table, then another migration to backfill all existing NULL values with a default avatar URL.

**Hard:** Design a schema for an e-commerce system with products, variants (size/color), orders, order items, shipping addresses, and payment methods. Account for products being discontinued and orders needing to remember the price they paid even if the product price changes later.
