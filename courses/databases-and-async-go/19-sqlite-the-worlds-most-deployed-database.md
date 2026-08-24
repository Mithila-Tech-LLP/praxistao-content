# Chapter 19: SQLite — The World's Most Deployed Database

SQLite is in your phone, your browser, your operating system. It powers WhatsApp messages, Firefox bookmarks, Apple Photos, and millions of mobile apps. It has no server, no daemon, no configuration — it's just a file. And it's remarkably capable.

## Table of Contents

1. What SQLite Is (and Is Not)
2. SQLite Architecture — The File Format
3. When to Use SQLite
4. Setting Up SQLite in Go
5. SQLite-Specific Features
6. WAL Mode — Enabling Concurrent Reads
7. SQLite for Testing
8. Mini Project: Local CLI Password Manager
9. Exercises

---

## 1. What SQLite Is (and Is Not)

**What SQLite is:**
- A complete SQL database engine in a single C library (~750 KB)
- The database is a single file (`.db`)
- No server, no network, no configuration
- ACID-compliant, supports transactions, foreign keys, triggers
- The most widely deployed database engine in the world

**What SQLite is NOT:**
- A replacement for PostgreSQL or MySQL in multi-user server applications
- Suitable for high-concurrency write workloads (one writer at a time)
- Suitable for databases > a few GB (works, but not optimized)

**Where SQLite shines:**
- Mobile apps (iOS, Android)
- Embedded systems and IoT devices
- Desktop applications
- Development and testing environments
- Read-heavy applications with occasional writes
- Replacing config files and CSVs with a proper queryable store

---

## 2. SQLite Architecture — The File Format

SQLite stores everything in a single binary file. The file is divided into **pages** (default 4096 bytes). The first page is the file header (100 bytes describing the format version, page size, etc.).

```
SQLite file (.db)
├── Page 1: File header + root of database catalog (sqlite_schema table)
├── Page 2: B-tree root for "users" table
├── Page 3: B-tree interior node for "users" table
├── Page 4: B-tree leaf node (actual rows)
├── Page 5: B-tree root for "users_email_idx" index
└── ...
```

The file format is stable and documented — a database created in 2005 still opens today. SQLite files are **portable**: you can copy a SQLite file between macOS, Linux, and Windows and it works.

```bash
# Inspect a SQLite database
sqlite3 myapp.db

sqlite> .tables         -- list tables
sqlite> .schema users   -- show CREATE TABLE for users
sqlite> PRAGMA page_size; -- show page size
sqlite> PRAGMA integrity_check; -- verify database integrity
```

---

## 3. When to Use SQLite

| Use Case | SQLite? |
|----------|---------|
| Mobile app local storage | Yes |
| Unit test database | Yes (in-memory) |
| Local CLI tools | Yes |
| Config/state files | Yes |
| Single-user desktop app | Yes |
| Read-heavy web app (< 10 writes/sec) | Yes |
| Multi-user web app | No — use PostgreSQL |
| > 10 concurrent writers | No |
| Multi-server deployment | No |
| Very large datasets (> 10 GB) | Maybe, test first |

SQLite's official documentation says it's appropriate for databases up to terabytes in size, but performance degrades vs PostgreSQL for multi-user workloads.

---

## 4. Setting Up SQLite in Go

There are two Go drivers for SQLite:

**mattn/go-sqlite3** — uses CGo (compiles SQLite C code), fastest, requires a C compiler.
```bash
go get github.com/mattn/go-sqlite3
```

**modernc.org/sqlite** — pure Go, no CGo, easier cross-compilation, slightly slower.
```bash
go get modernc.org/sqlite
```

For most projects, use `modernc.org/sqlite` for simplicity.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite" // register the sqlite driver
)

func main() {
    // Open (or create) a SQLite database file
    db, err := sql.Open("sqlite", "myapp.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // SQLite is a single file — verify it works
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    // Create a table
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS notes (
        id         INTEGER PRIMARY KEY AUTOINCREMENT,
        title      TEXT NOT NULL,
        content    TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
    if err != nil {
        log.Fatal(err)
    }

    // Insert a row
    result, err := db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)",
        "First Note", "Hello, SQLite!")
    if err != nil {
        log.Fatal(err)
    }
    id, _ := result.LastInsertId()
    fmt.Println("Inserted note ID:", id)

    // Query rows
    rows, err := db.Query("SELECT id, title, content FROM notes")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var title, content string
        rows.Scan(&id, &title, &content)
        fmt.Printf("[%d] %s: %s\n", id, title, content)
    }
}
```

Note: SQLite uses `?` for placeholders, not `$1` like PostgreSQL.

---

## 5. SQLite-Specific Features

### PRAGMAs — Configuring SQLite

PRAGMAs are SQLite-specific commands for configuration:

```go
func configureSQLite(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA journal_mode = WAL",       // enable WAL mode (crucial!)
        "PRAGMA synchronous = NORMAL",     // balance between safety and speed
        "PRAGMA foreign_keys = ON",        // enforce foreign keys (OFF by default!)
        "PRAGMA cache_size = -64000",      // 64 MB page cache
        "PRAGMA temp_store = MEMORY",      // store temp tables in RAM
        "PRAGMA busy_timeout = 5000",      // wait up to 5s when DB is locked
    }
    for _, p := range pragmas {
        if _, err := db.Exec(p); err != nil {
            return fmt.Errorf("pragma %q: %w", p, err)
        }
    }
    return nil
}
```

**Always set `PRAGMA foreign_keys = ON`** — SQLite doesn't enforce foreign keys by default!

### AUTOINCREMENT vs INTEGER PRIMARY KEY

```sql
-- INTEGER PRIMARY KEY alone: rowid alias, reuses deleted IDs
id INTEGER PRIMARY KEY

-- AUTOINCREMENT: strictly increasing, never reuses IDs (slower)
id INTEGER PRIMARY KEY AUTOINCREMENT
```

For most cases, `INTEGER PRIMARY KEY` (without AUTOINCREMENT) is fine and faster.

### Full-Text Search with FTS5

```sql
-- Create a virtual FTS5 table
CREATE VIRTUAL TABLE articles_fts USING fts5(
    title,
    content,
    content='articles',
    content_rowid='id'
);

-- Populate the index
INSERT INTO articles_fts(rowid, title, content)
SELECT id, title, content FROM articles;

-- Search
SELECT rowid, highlight(articles_fts, 1, '<b>', '</b>') as snippet
FROM articles_fts
WHERE articles_fts MATCH 'database performance'
ORDER BY rank;
```

---

## 6. WAL Mode — Enabling Concurrent Reads

By default, SQLite uses rollback journal mode: a write locks the entire database file. With WAL (Write-Ahead Log) mode, readers and writers don't block each other.

```go
db.Exec("PRAGMA journal_mode = WAL")
```

In WAL mode:
- Multiple readers can read simultaneously
- One writer at a time (writes don't block reads)
- WAL file grows as writes happen; periodically checkpointed back to the main file

For web applications using SQLite (like Litestream or Turso), WAL mode is essential.

### Connection Pool for SQLite

SQLite allows only **one writer at a time**. If multiple goroutines try to write simultaneously, they serialize. Configure properly:

```go
db, _ := sql.Open("sqlite", "myapp.db")

// Allow multiple readers, but serialized writers
db.SetMaxOpenConns(1)  // for write-heavy: serialize all access
// OR
db.SetMaxOpenConns(10) // for read-heavy with WAL mode

// Keep connections alive (opening SQLite is cheap, but still)
db.SetMaxIdleConns(5)
```

---

## 7. SQLite for Testing

SQLite's biggest superpower for Go developers: **in-memory databases** for testing.

```go
package mypackage_test

import (
    "database/sql"
    "testing"
    _ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
    // In-memory database: created fresh, destroyed when closed
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }

    // Apply your schema migrations
    _, err = db.Exec(`
        CREATE TABLE users (
            id    INTEGER PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            name  TEXT NOT NULL
        );
        CREATE TABLE orders (
            id      INTEGER PRIMARY KEY,
            user_id INTEGER REFERENCES users(id),
            total   REAL NOT NULL
        );
    `)
    if err != nil {
        t.Fatal(err)
    }

    // Register cleanup
    t.Cleanup(func() { db.Close() })
    return db
}

func TestCreateUser(t *testing.T) {
    db := newTestDB(t)

    // Test your actual function
    id, err := CreateUser(db, "alice@example.com", "Alice")
    if err != nil {
        t.Fatal("create user:", err)
    }
    if id <= 0 {
        t.Fatal("expected positive id, got", id)
    }
}
```

In-memory databases start empty, run your schema, and are destroyed after each test. No cleanup scripts, no test database, no test isolation issues. Tests run in milliseconds.

---

## 8. Mini Project: Local CLI Password Manager

A command-line password manager that stores encrypted entries in SQLite.

```go
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "database/sql"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "os"

    _ "modernc.org/sqlite"
)

const dbPath = "passwords.db"

func openDB() *sql.DB {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    db.Exec("PRAGMA journal_mode = WAL")
    db.Exec("PRAGMA foreign_keys = ON")
    db.Exec(`CREATE TABLE IF NOT EXISTS passwords (
        id       INTEGER PRIMARY KEY,
        site     TEXT NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL  -- stored encrypted
    )`)
    return db
}

func encrypt(key []byte, plaintext string) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, encoded string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return "", err
    }
    block, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(block)
    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    return string(plaintext), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: passwords [add|get|list] [args...]")
        os.Exit(1)
    }

    // In production, derive this from a master password using PBKDF2
    key := make([]byte, 32) // 256-bit key
    copy(key, []byte("your-32-byte-master-key-here!!"))

    db := openDB()
    defer db.Close()

    switch os.Args[1] {
    case "add":
        if len(os.Args) < 5 {
            fmt.Println("Usage: passwords add <site> <username> <password>")
            os.Exit(1)
        }
        site, user, pass := os.Args[2], os.Args[3], os.Args[4]
        encrypted, err := encrypt(key, pass)
        if err != nil {
            log.Fatal(err)
        }
        db.Exec("INSERT INTO passwords (site, username, password) VALUES (?, ?, ?)",
            site, user, encrypted)
        fmt.Println("Saved!")

    case "get":
        if len(os.Args) < 3 {
            fmt.Println("Usage: passwords get <site>")
            os.Exit(1)
        }
        rows, _ := db.Query(
            "SELECT site, username, password FROM passwords WHERE site LIKE ?",
            "%"+os.Args[2]+"%")
        defer rows.Close()
        for rows.Next() {
            var site, user, enc string
            rows.Scan(&site, &user, &enc)
            pass, _ := decrypt(key, enc)
            fmt.Printf("Site: %s  User: %s  Pass: %s\n", site, user, pass)
        }

    case "list":
        rows, _ := db.Query("SELECT site, username FROM passwords ORDER BY site")
        defer rows.Close()
        for rows.Next() {
            var site, user string
            rows.Scan(&site, &user)
            fmt.Printf("%-30s %s\n", site, user)
        }
    }
}
```

---

## Summary

- SQLite is a file-based, serverless SQL database ideal for embedded, mobile, and local applications.
- Always enable `PRAGMA journal_mode = WAL` for concurrent reads and `PRAGMA foreign_keys = ON`.
- Use `modernc.org/sqlite` (pure Go) or `mattn/go-sqlite3` (CGo, faster).
- In-memory SQLite (`:memory:`) is perfect for unit tests — fast, isolated, no cleanup.
- SQLite handles one writer at a time; for multi-writer scenarios, use PostgreSQL.

### Exercises

**Easy:** Build a Go CLI tool that stores a shopping list in SQLite. Commands: `add <item>`, `list`, `done <id>`.

**Medium:** Write a test suite for a data access layer using in-memory SQLite. Run schema migrations in a `TestMain` function shared across all tests.

**Hard:** Implement the password manager with a real master password: use PBKDF2 to derive the encryption key from the master password, and store a salt in a separate config table.
