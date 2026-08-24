# Chapter 20: Building with SQLite in Go

SQLite's killer feature for Go developers is simplicity: no Docker, no server, no connection string secrets. Drop a file, open it, query it. This chapter covers advanced SQLite patterns and builds a complete offline-first task manager.

## Table of Contents

1. Connection String Options and Pragmas
2. Migrations in SQLite
3. SQLite-Specific SQL
4. Performance Tuning SQLite
5. Backing Up a Live SQLite Database
6. Project: Offline-First Task Manager CLI
7. Exercises

---

## 1. Connection String Options and Pragmas

```go
package main

import (
    "database/sql"
    "fmt"
    _ "modernc.org/sqlite"
)

func openDB(path string) (*sql.DB, error) {
    // SQLite connection string options
    // _journal_mode=WAL: write-ahead log mode
    // _foreign_keys=on:  enforce FK constraints
    // _busy_timeout=5000: wait 5s if DB is locked
    // _synchronous=NORMAL: balance safety and speed
    dsn := fmt.Sprintf("%s?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=5000&_synchronous=NORMAL", path)
    
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }

    // SQLite can handle multiple readers with WAL, but only one writer
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)

    return db, db.Ping()
}
```

You can set pragmas either in the DSN or with explicit `PRAGMA` SQL statements:

```go
func setPragmas(db *sql.DB) {
    // After opening the connection
    db.Exec("PRAGMA cache_size = -32000")  // 32 MB cache
    db.Exec("PRAGMA temp_store = MEMORY")   // temp tables in RAM
    db.Exec("PRAGMA mmap_size = 268435456") // 256 MB memory-mapped I/O
}
```

---

## 2. Migrations in SQLite

SQLite has a built-in user_version pragma that you can use to track schema version:

```go
package main

import (
    "database/sql"
    "fmt"
)

type migration struct {
    version int
    sql     string
}

var migrations = []migration{
    {1, `CREATE TABLE tasks (
        id         INTEGER PRIMARY KEY,
        title      TEXT NOT NULL,
        done       INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`},
    {2, `ALTER TABLE tasks ADD COLUMN priority INTEGER NOT NULL DEFAULT 0`},
    {3, `CREATE TABLE projects (
        id   INTEGER PRIMARY KEY,
        name TEXT NOT NULL UNIQUE
    );
    ALTER TABLE tasks ADD COLUMN project_id INTEGER REFERENCES projects(id);`},
}

func migrateDB(db *sql.DB) error {
    var version int
    db.QueryRow("PRAGMA user_version").Scan(&version)
    fmt.Println("Current schema version:", version)

    for _, m := range migrations {
        if m.version <= version {
            continue // already applied
        }
        fmt.Printf("Applying migration v%d...\n", m.version)

        tx, err := db.Begin()
        if err != nil {
            return err
        }

        if _, err := tx.Exec(m.sql); err != nil {
            tx.Rollback()
            return fmt.Errorf("migration v%d: %w", m.version, err)
        }

        // Update version INSIDE the same transaction
        tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", m.version))
        if err := tx.Commit(); err != nil {
            return err
        }
        fmt.Printf("Migration v%d applied\n", m.version)
    }
    return nil
}
```

### SQLite ALTER TABLE Limitations

SQLite only supports a small subset of ALTER TABLE:
- `ALTER TABLE t ADD COLUMN c type` ✓
- `ALTER TABLE t RENAME TO new_name` ✓
- `ALTER TABLE t RENAME COLUMN c TO new_c` ✓ (SQLite 3.25+)
- `ALTER TABLE t DROP COLUMN c` ✓ (SQLite 3.35+)
- Changing column type, adding constraints to existing columns: NOT SUPPORTED

To change a column type, you must recreate the table:

```sql
-- SQLite's way to change column type
BEGIN;
CREATE TABLE tasks_new (
    id    INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    done  BOOLEAN NOT NULL DEFAULT FALSE  -- changed type
);
INSERT INTO tasks_new SELECT * FROM tasks;
DROP TABLE tasks;
ALTER TABLE tasks_new RENAME TO tasks;
COMMIT;
```

---

## 3. SQLite-Specific SQL

### Upsert with INSERT OR REPLACE / ON CONFLICT

```sql
-- Upsert: replace if primary key conflicts
INSERT OR REPLACE INTO tasks (id, title, done) VALUES (1, 'Updated title', 0);

-- Upsert with ON CONFLICT (SQLite 3.24+)
INSERT INTO tasks (id, title) VALUES (1, 'New title')
ON CONFLICT(id) DO UPDATE SET title = EXCLUDED.title;
```

### WITH RECURSIVE — Tree Queries

SQLite supports CTEs and recursive queries:

```sql
-- Find all subtasks recursively
WITH RECURSIVE subtasks AS (
    SELECT id, title, parent_id, 0 AS depth
    FROM tasks WHERE id = 5                -- root task
    UNION ALL
    SELECT t.id, t.title, t.parent_id, s.depth + 1
    FROM tasks t
    JOIN subtasks s ON t.parent_id = s.id  -- children
)
SELECT * FROM subtasks ORDER BY depth, id;
```

### JSON Functions

SQLite has built-in JSON functions:

```sql
-- Store and query JSON
CREATE TABLE config (key TEXT PRIMARY KEY, value TEXT);
INSERT INTO config VALUES ('settings', '{"theme":"dark","lang":"en"}');

SELECT json_extract(value, '$.theme') FROM config WHERE key = 'settings';
-- dark

SELECT json_set(value, '$.theme', 'light') FROM config WHERE key = 'settings';
-- {"theme":"light","lang":"en"}
```

---

## 4. Performance Tuning SQLite

### Batch Inserts — Always Use Transactions

```go
func insertManyTasks(db *sql.DB, tasks []Task) error {
    // Without transaction: 1000 inserts = 1000 disk syncs = VERY SLOW
    // With transaction: 1000 inserts = 1 disk sync = FAST

    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare("INSERT INTO tasks (title, done) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, t := range tasks {
        if _, err := stmt.Exec(t.Title, t.Done); err != nil {
            return err
        }
    }
    return tx.Commit()
}
```

Benchmarks: inserting 100,000 rows:
- Without transaction: ~100 seconds
- With transaction: ~0.1 seconds (1000x faster!)

### Prepared Statements

Prepare a statement once, execute many times:

```go
stmt, err := db.Prepare("SELECT * FROM tasks WHERE done = ? LIMIT ?")
defer stmt.Close()

// Reuse stmt many times without re-parsing
stmt.QueryRow(false, 10)
stmt.QueryRow(true, 5)
```

### Index Your Queries

```sql
-- Common query: active tasks by project, ordered by creation
CREATE INDEX idx_tasks_project_done ON tasks(project_id, done, created_at DESC);
```

---

## 5. Backing Up a Live SQLite Database

Unlike PostgreSQL, you can copy a SQLite file while it's in use — with WAL mode and the backup API:

```go
import "github.com/mattn/go-sqlite3"

func backupDatabase(srcPath, dstPath string) error {
    srcDB, err := sql.Open("sqlite3", srcPath)
    if err != nil {
        return err
    }
    defer srcDB.Close()

    dstDB, err := sql.Open("sqlite3", dstPath)
    if err != nil {
        return err
    }
    defer dstDB.Close()

    srcConn, err := srcDB.Conn(context.Background())
    if err != nil {
        return err
    }
    defer srcConn.Close()

    return srcConn.Raw(func(srcRaw interface{}) error {
        srcSqlite := srcRaw.(*sqlite3.SQLiteConn)
        dstConn, _ := dstDB.Conn(context.Background())
        defer dstConn.Close()

        return dstConn.Raw(func(dstRaw interface{}) error {
            dstSqlite := dstRaw.(*sqlite3.SQLiteConn)
            backup, err := dstSqlite.Backup("main", srcSqlite, "main")
            if err != nil {
                return err
            }
            _, err = backup.Step(-1) // copy all pages
            backup.Finish()
            return err
        })
    })
}
```

Or just use the SQLite `VACUUM INTO` command (simpler):

```sql
VACUUM INTO '/path/to/backup.db';
```

---

## 6. Project: Offline-First Task Manager CLI

A complete command-line task manager:

```go
package main

import (
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "time"

    _ "modernc.org/sqlite"
)

type Task struct {
    ID        int
    Title     string
    Done      bool
    Priority  int
    CreatedAt time.Time
}

func main() {
    db := mustOpenDB("tasks.db")
    defer db.Close()
    mustMigrate(db)

    if len(os.Args) < 2 {
        printUsage()
        return
    }

    switch os.Args[1] {
    case "add":
        if len(os.Args) < 3 {
            fmt.Println("Usage: task add <title>")
            return
        }
        addTask(db, os.Args[2])

    case "list":
        listTasks(db, false)

    case "done":
        if len(os.Args) < 3 {
            fmt.Println("Usage: task done <id>")
            return
        }
        id, _ := strconv.Atoi(os.Args[2])
        markDone(db, id)

    case "delete":
        if len(os.Args) < 3 {
            fmt.Println("Usage: task delete <id>")
            return
        }
        id, _ := strconv.Atoi(os.Args[2])
        deleteTask(db, id)

    case "all":
        listTasks(db, true)

    default:
        printUsage()
    }
}

func mustOpenDB(path string) *sql.DB {
    db, err := sql.Open("sqlite",
        path+"?_journal_mode=WAL&_foreign_keys=on")
    if err != nil {
        fmt.Fprintf(os.Stderr, "open db: %v\n", err)
        os.Exit(1)
    }
    return db
}

func mustMigrate(db *sql.DB) {
    db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id         INTEGER PRIMARY KEY,
        title      TEXT NOT NULL,
        done       INTEGER NOT NULL DEFAULT 0,
        priority   INTEGER NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
}

func addTask(db *sql.DB, title string) {
    result, err := db.Exec("INSERT INTO tasks (title) VALUES (?)", title)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    id, _ := result.LastInsertId()
    fmt.Printf("Added task #%d: %s\n", id, title)
}

func listTasks(db *sql.DB, showAll bool) {
    query := "SELECT id, title, done, created_at FROM tasks WHERE done = 0 ORDER BY priority DESC, created_at"
    if showAll {
        query = "SELECT id, title, done, created_at FROM tasks ORDER BY done, priority DESC, created_at"
    }

    rows, err := db.Query(query)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer rows.Close()

    fmt.Println("─────────────────────────────")
    count := 0
    for rows.Next() {
        var id int
        var title string
        var done int
        var createdAt time.Time
        rows.Scan(&id, &title, &done, &createdAt)

        status := "[ ]"
        if done == 1 {
            status = "[x]"
        }
        fmt.Printf("%s #%3d  %s  (%s)\n", status, id, title, createdAt.Format("Jan 2"))
        count++
    }
    if count == 0 {
        fmt.Println("No tasks!")
    }
    fmt.Println("─────────────────────────────")
}

func markDone(db *sql.DB, id int) {
    result, err := db.Exec("UPDATE tasks SET done = 1 WHERE id = ?", id)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    if n, _ := result.RowsAffected(); n == 0 {
        fmt.Printf("Task #%d not found\n", id)
        return
    }
    fmt.Printf("Task #%d marked done!\n", id)
}

func deleteTask(db *sql.DB, id int) {
    db.Exec("DELETE FROM tasks WHERE id = ?", id)
    fmt.Printf("Task #%d deleted\n", id)
}

func printUsage() {
    fmt.Println(`Usage:
  task add <title>    Add a new task
  task list           List pending tasks
  task all            List all tasks (including done)
  task done <id>      Mark task as done
  task delete <id>    Delete a task`)
}
```

```bash
go run main.go add "Buy groceries"
go run main.go add "Write blog post"
go run main.go list
go run main.go done 1
```

---

## Summary

- Set pragmas in the DSN or with explicit `PRAGMA` statements after opening. Always enable WAL and foreign keys.
- Use `PRAGMA user_version` for simple schema migrations in embedded SQLite applications.
- SQLite ALTER TABLE is limited — use create-new-table-and-insert for structural changes.
- Batch inserts MUST use transactions — without them, each INSERT is a separate disk sync.
- `VACUUM INTO` backs up a live database instantly.

### Exercises

**Easy:** Add a `priority` command to the task manager: `task priority <id> <1-3>`. Display high-priority tasks with a `!` marker.

**Medium:** Add project support: `task project create <name>`, `task add <title> --project <name>`. Group tasks by project in the list view.

**Hard:** Implement offline sync: tasks have a `synced_at` column. Build a function that POSTs all unsynced tasks to a mock HTTP server and marks them synced on success.
