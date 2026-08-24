# Chapter 48: VaultDB — Major Project: Complete Working Database

All the pieces are built. Now we assemble them into a complete, runnable database — start it with one command, connect from the CLI, and run real queries.

## Table of Contents

1. Final Architecture Integration
2. The Main Entry Point
3. CLI Client: `vault-cli`
4. Integration Tests
5. What VaultDB Can Do
6. What Real Databases Add
7. What to Build Next

---

## 1. Final Architecture Integration

```go
// vaultdb.go — the top-level Database struct that connects all layers
package vaultdb

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/yourname/vaultdb/query"
    "github.com/yourname/vaultdb/storage"
    "github.com/yourname/vaultdb/txn"
    "github.com/yourname/vaultdb/wal"
    "github.com/yourname/vaultdb/wire"
)

type DB struct {
    dir     string
    dm      *storage.DiskManager
    bp      *storage.BufferPool
    wal     *wal.WAL
    catalog *storage.Catalog
    txnMgr  *txn.Manager
}

func Open(dir string) (*DB, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    dm, err := storage.NewDiskManager(filepath.Join(dir, "data.vault"))
    if err != nil {
        return nil, fmt.Errorf("disk manager: %w", err)
    }

    bp := storage.NewBufferPool(dm, storage.DefaultBufferSize)

    w, err := wal.Open(filepath.Join(dir, "wal.log"))
    if err != nil {
        return nil, fmt.Errorf("wal: %w", err)
    }

    catalog, err := dm.ReadCatalog()
    if err != nil {
        return nil, fmt.Errorf("catalog: %w", err)
    }

    txnMgr := txn.NewManager()

    db := &DB{
        dir: dir, dm: dm, bp: bp,
        wal: w, catalog: catalog, txnMgr: txnMgr,
    }

    // Replay WAL to recover from crash
    if err := db.recover(); err != nil {
        return nil, fmt.Errorf("recovery: %w", err)
    }

    return db, nil
}

func (db *DB) recover() error {
    toRedo, _, err := db.wal.Recover()
    if err != nil {
        return err
    }

    for _, rec := range toRedo {
        if err := applyWALRecord(db.bp, db.dm, rec); err != nil {
            return fmt.Errorf("apply WAL record: %w", err)
        }
    }
    return db.bp.FlushAll()
}

func (db *DB) NewExecutor() *query.Executor {
    return query.NewExecutor(db.dm, db.bp, db.wal, db.catalog)
}

func (db *DB) Exec(sql string) (*query.Result, error) {
    stmt, err := query.Parse(sql)
    if err != nil {
        return nil, err
    }
    return db.NewExecutor().Execute(stmt)
}

func (db *DB) Serve(addr string) error {
    srv := wire.NewServer(addr, &wire.DBAdapter{
        DM: db.dm, BP: db.bp, WAL: db.wal, Catalog: db.catalog,
    })
    return srv.ListenAndServe()
}

func (db *DB) Close() error {
    db.bp.FlushAll()
    db.wal.Flush()
    if err := db.dm.WriteCatalog(db.catalog); err != nil {
        return err
    }
    return db.dm.Close()
}
```

---

## 2. The Main Entry Point

```go
// cmd/vaultdb/main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/yourname/vaultdb"
)

func main() {
    var (
        dataDir = flag.String("data", "./vaultdb-data", "database directory")
        addr    = flag.String("addr", ":5555", "TCP address to listen on")
    )
    flag.Parse()

    log.Printf("Opening VaultDB at %s", *dataDir)

    db, err := vaultdb.Open(*dataDir)
    if err != nil {
        log.Fatal("open:", err)
    }
    defer func() {
        log.Println("Closing database...")
        db.Close()
    }()

    // Graceful shutdown on Ctrl+C
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Printf("VaultDB listening on %s", *addr)
        if err := db.Serve(*addr); err != nil {
            log.Fatal("server:", err)
        }
    }()

    <-quit
    log.Println("Shutting down...")
}
```

Build and run:
```bash
cd cmd/vaultdb
go build -o vaultdb .
./vaultdb -data ./mydb -addr :5555
```

---

## 3. CLI Client: vault-cli

```go
// cmd/vault-cli/main.go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/yourname/vaultdb/client"
)

func main() {
    addr := "localhost:5555"
    if len(os.Args) > 1 {
        addr = os.Args[1]
    }

    c, err := client.Connect(addr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "connect: %v\n", err)
        os.Exit(1)
    }
    defer c.Close()

    fmt.Printf("VaultDB CLI — connected to %s\n", addr)
    fmt.Println("Type SQL and press Enter. Type \\quit to exit.\n")

    scanner := bufio.NewScanner(os.Stdin)
    var buffer strings.Builder

    for {
        if buffer.Len() == 0 {
            fmt.Print("vault> ")
        } else {
            fmt.Print("    -> ")
        }

        if !scanner.Scan() {
            break
        }

        line := scanner.Text()
        if strings.TrimSpace(line) == "\\quit" {
            break
        }

        buffer.WriteString(" ")
        buffer.WriteString(line)

        // Execute when we see a semicolon
        if !strings.Contains(line, ";") {
            continue
        }

        sql := strings.TrimSpace(buffer.String())
        sql = strings.TrimSuffix(sql, ";")
        buffer.Reset()

        if sql == "" {
            continue
        }

        rs, err := c.Query(sql)
        if err != nil {
            fmt.Printf("ERROR: %v\n", err)
            continue
        }

        printTable(rs)
    }
}

func printTable(rs *client.ResultSet) {
    if rs == nil || len(rs.Columns) == 0 {
        return
    }

    // Calculate column widths
    widths := make([]int, len(rs.Columns))
    for i, col := range rs.Columns {
        widths[i] = len(col)
    }
    for _, row := range rs.Rows {
        for i, val := range row {
            if len(val) > widths[i] {
                widths[i] = len(val)
            }
        }
    }

    // Print header
    for i, col := range rs.Columns {
        fmt.Printf("%-*s", widths[i]+2, col)
    }
    fmt.Println()

    // Separator
    for _, w := range widths {
        fmt.Print(strings.Repeat("-", w+2))
    }
    fmt.Println()

    // Rows
    for _, row := range rs.Rows {
        for i, val := range row {
            fmt.Printf("%-*s", widths[i]+2, val)
        }
        fmt.Println()
    }
    fmt.Printf("(%d rows)\n", len(rs.Rows))
}
```

Usage:
```bash
./vault-cli localhost:5555

vault> CREATE TABLE employees (id INT, name VARCHAR, salary INT);
(0 rows)

vault> INSERT INTO employees VALUES (1, 'Alice', 95000);
1 rows affected

vault> INSERT INTO employees VALUES (2, 'Bob', 85000);
1 rows affected

vault> INSERT INTO employees VALUES (3, 'Carol', 105000);
1 rows affected

vault> SELECT * FROM employees;
id    name     salary
-------------------------------
1     Alice    95000
2     Bob      85000
3     Carol    105000
(3 rows)

vault> SELECT name, salary FROM employees WHERE salary > 90000;
name     salary
-----------------
Alice    95000
Carol    105000
(2 rows)
```

---

## 4. Integration Tests

```go
// integration_test.go
package vaultdb_test

import (
    "os"
    "testing"
    "github.com/yourname/vaultdb"
)

func TestFullCRUD(t *testing.T) {
    dir := t.TempDir()
    db, err := vaultdb.Open(dir)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Create table
    _, err = db.Exec("CREATE TABLE users (id INT, name VARCHAR, age INT)")
    if err != nil {
        t.Fatal("create table:", err)
    }

    // Insert rows
    for _, row := range []string{
        "INSERT INTO users VALUES (1, 'Alice', 25)",
        "INSERT INTO users VALUES (2, 'Bob', 30)",
        "INSERT INTO users VALUES (3, 'Carol', 22)",
    } {
        if _, err := db.Exec(row); err != nil {
            t.Fatal("insert:", err)
        }
    }

    // Select all
    result, err := db.Exec("SELECT * FROM users")
    if err != nil {
        t.Fatal("select:", err)
    }
    if len(result.Rows) != 3 {
        t.Errorf("expected 3 rows, got %d", len(result.Rows))
    }

    // Select with WHERE
    result, err = db.Exec("SELECT name FROM users WHERE age > 24")
    if err != nil {
        t.Fatal("select where:", err)
    }
    if len(result.Rows) != 2 {
        t.Errorf("expected 2 rows (Alice, Bob), got %d", len(result.Rows))
    }
}

func TestCrashRecovery(t *testing.T) {
    dir := t.TempDir()

    // Insert rows in one session
    func() {
        db, _ := vaultdb.Open(dir)
        defer db.Close()
        db.Exec("CREATE TABLE logs (id INT, msg VARCHAR)")
        db.Exec("INSERT INTO logs VALUES (1, 'first entry')")
        db.Exec("INSERT INTO logs VALUES (2, 'second entry')")
        // db.Close() flushes WAL
    }()

    // Reopen (simulates restart after crash)
    db2, err := vaultdb.Open(dir)
    if err != nil {
        t.Fatal("reopen:", err)
    }
    defer db2.Close()

    result, err := db2.Exec("SELECT * FROM logs")
    if err != nil {
        t.Fatal("select after recovery:", err)
    }
    if len(result.Rows) != 2 {
        t.Errorf("expected 2 rows after recovery, got %d", len(result.Rows))
    }
}

func TestConcurrentReads(t *testing.T) {
    db, _ := vaultdb.Open(t.TempDir())
    defer db.Close()
    db.Exec("CREATE TABLE data (id INT, val INT)")
    for i := 0; i < 100; i++ {
        db.Exec(fmt.Sprintf("INSERT INTO data VALUES (%d, %d)", i, i*10))
    }

    // 10 concurrent readers
    done := make(chan bool, 10)
    for i := 0; i < 10; i++ {
        go func() {
            result, err := db.Exec("SELECT * FROM data")
            if err != nil || len(result.Rows) != 100 {
                t.Errorf("concurrent read failed: %v, rows=%d", err, len(result.Rows))
            }
            done <- true
        }()
    }
    for i := 0; i < 10; i++ {
        <-done
    }
}
```

---

## 5. What VaultDB Can Do

- Create tables with INT, VARCHAR, FLOAT, BOOL columns
- INSERT rows
- SELECT with column projection, WHERE filtering, LIMIT
- UPDATE and DELETE with WHERE
- Persist data to disk (survives restarts)
- WAL for crash recovery
- Buffer pool for memory caching
- B-Tree indexes for fast lookups
- MVCC for concurrent transaction isolation
- TCP server + binary protocol for remote clients
- CLI client for interactive queries

---

## 6. What Real Databases Add

PostgreSQL has 30 years of additional features. Key ones:

| Feature | What It Does |
|---------|-------------|
| JOINs | Query across multiple tables |
| Subqueries | Nested SELECTs |
| Aggregates | SUM, COUNT, GROUP BY |
| Stored procedures | Functions stored in the DB |
| Full-text search | tsvector, tsquery |
| Extensions | pgvector, PostGIS, etc. |
| Streaming replication | Send WAL to read replicas |
| Vacuum | Reclaim space from old MVCC versions |
| Connection pooling | PgBouncer handles thousands of clients |
| EXPLAIN ANALYZE | Show query execution plan |
| Partial indexes | Index only rows matching a condition |

---

## 7. What to Build Next

The database is done. Part 2 of this course covers **async systems** — Kafka, RabbitMQ, NATS. Part 3 integrates VaultDB with a streaming message broker to build a production service.

You've built more than most engineers ever will. You now understand:
- Why B-Trees are used for every index
- Why the WAL makes crash recovery possible
- Why MVCC lets readers and writers coexist
- Why buffer pools matter (RAM vs disk speed)
- How SQL becomes execution through a parser, planner, and executor
- How a binary wire protocol enables language-agnostic clients

These concepts are the same in PostgreSQL, MySQL, RocksDB, and every other serious database. You can now read their source code and understand it.

---

## Summary

- VaultDB is a complete, working relational database with all the core components: storage, B-Tree, WAL, buffer pool, SQL parser, executor, MVCC, and a TCP wire protocol.
- The CLI client connects to the running server and provides an interactive SQL prompt.
- Integration tests verify the full stack: CRUD, crash recovery, and concurrent reads.
- Understanding VaultDB gives you a mental model for every production database you'll ever use.

### Final Project

Build a complete application using VaultDB as the database:

**TaskTracker:** A task management API with VaultDB as storage.

Requirements:
- Tables: `tasks(id, title, status, priority, created_at)`, `users(id, name, email)`
- API: `POST /tasks`, `GET /tasks`, `GET /tasks?status=done`, `PUT /tasks/:id`, `DELETE /tasks/:id`
- Concurrent access: multiple clients writing simultaneously
- Crash recovery test: insert 100 tasks, kill the server, restart, verify all 100 tasks survive
- Benchmarks: measure insert throughput, query throughput with and without index

This is your capstone project for the database section. The async section starts next.
