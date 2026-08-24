# Chapter 21: MySQL — The Web's Favorite Database

MySQL powers WordPress, Wikipedia, YouTube, Facebook (before they built their own systems), and millions of other websites. It's the "M" in LAMP stack and the most popular database for web applications in the world. Understanding MySQL is essential for any web developer.

## Table of Contents

1. MySQL History and the Oracle Question
2. MySQL vs PostgreSQL — Honest Trade-offs
3. InnoDB — MySQL's Default Storage Engine
4. MySQL Architecture
5. MySQL-Specific SQL Features
6. Replication and High Availability
7. MySQL 8.0 Modern Features
8. Building with MySQL in Go
9. Exercises

---

## 1. MySQL History and the Oracle Question

MySQL was created in 1995 by Michael Widenius (who named it after his daughter "My"). In 2010, Oracle acquired Sun Microsystems (which had bought MySQL), causing community concern about Oracle's stewardship.

This led to the **MariaDB** fork — a drop-in replacement for MySQL created by MySQL's original developers. MariaDB is MySQL-compatible but adds features independently.

For most purposes, MySQL and MariaDB are interchangeable at the SQL level. This chapter covers MySQL 8.0.

---

## 2. MySQL vs PostgreSQL — Honest Trade-offs

| Aspect | MySQL | PostgreSQL |
|--------|-------|------------|
| Maturity | 1995 | 1986 |
| SQL standards | Less strict | Very strict |
| JSON support | JSON type (8.0) | JSONB (much better) |
| Full-text search | Built-in (limited) | Built-in (better) |
| Replication | Binary log (BinLog) | WAL streaming |
| Max table size | 256 TB | Unlimited |
| Concurrency | Good | Excellent (MVCC) |
| Extensions | Limited | Rich ecosystem |
| Community | Huge | Large and growing |

**Choose MySQL when:**
- Existing team experience
- LAMP/LEMP stack applications
- WordPress, Drupal, many PHP frameworks

**Choose PostgreSQL when:**
- New projects with no legacy constraints
- Complex queries, JSONB, geospatial
- Strict SQL standards compliance

---

## 3. InnoDB — MySQL's Default Storage Engine

InnoDB is MySQL's default storage engine since MySQL 5.5. It provides:

- **ACID transactions** with row-level locking
- **Foreign key support**
- **MVCC** (multi-version concurrency control)
- **Buffer pool** for caching data pages in memory
- **Clustered primary key index** — data is physically ordered by PK

### Clustered Primary Key

Unlike PostgreSQL (heap storage), InnoDB **physically sorts the table by primary key**. This means:

```sql
-- Sequential PK inserts: fast (appends to end)
INSERT INTO users (id, name) VALUES (1, 'A'), (2, 'B'), (3, 'C');

-- Random/UUID PKs: slow (must insert in the middle, causes page splits)
-- For high-insert tables, avoid UUIDs as PKs in MySQL!
INSERT INTO users (id, name) VALUES (UUID(), 'A');  -- bad for performance
```

**Best practice for MySQL:** Use `BIGINT AUTO_INCREMENT` as the primary key. If you need UUIDs for external exposure, add a separate `uuid` column with a unique index.

---

## 4. MySQL Architecture

```
Client (your Go app)
    │
    │ MySQL protocol (port 3306)
    ▼
Connection handler (one thread per connection)
    │
    ├── Query cache (deprecated in 8.0)
    ├── Parser & Optimizer
    │
    ▼
InnoDB Storage Engine
    │
    ├── Buffer Pool (caches pages in RAM)
    ├── Log Buffer → redo log (InnoDB's WAL)
    │
    ▼
Disk
    ├── Tablespace files (.ibd per table in 8.0)
    └── Undo logs (for MVCC)
```

**Key difference from PostgreSQL:** MySQL uses a **redo log** (InnoDB) in addition to the **binary log** (MySQL server layer). The binary log is used for replication; the redo log is for crash recovery.

---

## 5. MySQL-Specific SQL Features

### AUTO_INCREMENT

```sql
CREATE TABLE users (
    id    BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name  VARCHAR(100) NOT NULL
);

-- MySQL returns the last inserted ID
INSERT INTO users (email, name) VALUES ('alice@example.com', 'Alice');
SELECT LAST_INSERT_ID();  -- returns the new id
```

### String Types

```sql
-- VARCHAR: variable length (up to specified max)
name VARCHAR(100)      -- good default

-- TEXT: no size limit (stored separately from row for large values)
content TEXT           -- for long content

-- CHAR: fixed length (pads with spaces)
country_code CHAR(2)   -- exactly 2 characters, faster for fixed-length lookups
```

### DATETIME vs TIMESTAMP

```sql
-- DATETIME: stores literally what you give it (no timezone)
-- Range: 1000-01-01 to 9999-12-31
created_at DATETIME DEFAULT CURRENT_TIMESTAMP

-- TIMESTAMP: stores UTC, auto-converts on read
-- Range: 1970-01-01 to 2038-01-19 (32-bit Unix timestamp!)
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```

**Warning:** `TIMESTAMP` has a 2038 problem. For new tables, use `DATETIME` and handle timezone in your application, or use MySQL 8+ with `TIMESTAMP(6)` that avoids this.

### ENUM Type

```sql
-- MySQL-specific ENUM type
status ENUM('pending', 'active', 'suspended') NOT NULL DEFAULT 'pending'
```

ENUM stores values as integers (1, 2, 3) but displays as strings. Avoid if values might change — adding new ENUM values requires an ALTER TABLE.

### Stored Procedures and Triggers

```sql
-- Trigger: automatically update updated_at
CREATE TRIGGER users_updated
BEFORE UPDATE ON users
FOR EACH ROW SET NEW.updated_at = NOW();
```

---

## 6. Replication and High Availability

MySQL replication uses the **binary log** (binlog), which records every change as either:
- **Statement-based:** the SQL statement itself (smaller, but can be non-deterministic)
- **Row-based:** the actual before/after row values (larger, but deterministic)
- **Mixed:** MySQL chooses based on the statement

```
Primary (writes)
    │
    │ binlog stream
    ▼
Replica 1 (read-only)
Replica 2 (read-only)
```

Setting up:
```sql
-- On primary
SET GLOBAL binlog_format = 'ROW';
CREATE USER 'replicator'@'%' IDENTIFIED BY 'secret';
GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';
SHOW MASTER STATUS;  -- note File and Position

-- On replica
CHANGE MASTER TO
    MASTER_HOST='primary_host',
    MASTER_USER='replicator',
    MASTER_PASSWORD='secret',
    MASTER_LOG_FILE='mysql-bin.000001',
    MASTER_LOG_POS=123;
START SLAVE;
SHOW SLAVE STATUS\G  -- check Seconds_Behind_Master
```

### MySQL Group Replication (InnoDB Cluster)

For multi-primary or automatic failover, use MySQL InnoDB Cluster with MySQL Group Replication. This is more complex but provides automatic failover without manual intervention.

---

## 7. MySQL 8.0 Modern Features

MySQL 8.0 (2018+) added many PostgreSQL-like features:

### CTEs and Window Functions

```sql
-- CTEs (WITH clause)
WITH monthly_revenue AS (
    SELECT DATE_FORMAT(created_at, '%Y-%m') AS month,
           SUM(total) AS revenue
    FROM orders
    GROUP BY month
)
SELECT * FROM monthly_revenue ORDER BY month;

-- Window functions
SELECT
    user_id,
    SUM(amount) OVER (PARTITION BY user_id ORDER BY created_at) AS running_total
FROM payments;
```

### JSON Type

```sql
CREATE TABLE products (
    id       INT AUTO_INCREMENT PRIMARY KEY,
    metadata JSON
);

INSERT INTO products (metadata)
VALUES ('{"brand":"Apple","price":999}');

-- Extract JSON values
SELECT metadata->>'$.brand' FROM products;

-- Index a JSON path
CREATE INDEX idx_products_brand ON products ((metadata->>'$.brand'));
```

### Descending Indexes

```sql
-- MySQL 8.0 supports DESC indexes natively
CREATE INDEX idx_orders_created_desc ON orders(created_at DESC);
-- Efficient for ORDER BY created_at DESC LIMIT 10
```

---

## 8. Building with MySQL in Go

```bash
go get github.com/go-sql-driver/mysql
```

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID        int
    Email     string
    Name      string
    CreatedAt time.Time
}

func main() {
    // DSN format: user:password@tcp(host:port)/dbname?params
    dsn := "dev:secret@tcp(localhost:3306)/myapp?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // MySQL connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := db.Ping(); err != nil {
        log.Fatal("cannot connect:", err)
    }
    fmt.Println("Connected to MySQL!")

    // Insert with LAST_INSERT_ID
    result, err := db.Exec(
        "INSERT INTO users (email, name) VALUES (?, ?)",
        "alice@example.com", "Alice",
    )
    if err != nil {
        log.Fatal(err)
    }
    id, _ := result.LastInsertId()
    fmt.Println("Inserted ID:", id)

    // Query
    rows, err := db.Query("SELECT id, email, name, created_at FROM users LIMIT 10")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%+v\n", u)
    }
}
```

Important DSN parameters:
- `parseTime=true` — parse `DATETIME`/`TIMESTAMP` as `time.Time`
- `charset=utf8mb4` — use full Unicode (emoji support). Never use `utf8` in MySQL (it's broken — only 3 bytes)
- `loc=UTC` — interpret times as UTC

### Handling MySQL-Specific Errors in Go

```go
import "github.com/go-sql-driver/mysql"

func isMySQLError(err error, code uint16) bool {
    if mysqlErr, ok := err.(*mysql.MySQLError); ok {
        return mysqlErr.Number == code
    }
    return false
}

// MySQL error codes
// 1062: duplicate entry (unique violation)
// 1451: foreign key constraint fails (parent row referenced)
// 1452: foreign key constraint fails (child row doesn't exist)
// 1205: lock wait timeout exceeded
// 1213: deadlock found

func createUser(db *sql.DB, email, name string) (int64, error) {
    result, err := db.Exec("INSERT INTO users (email, name) VALUES (?, ?)", email, name)
    if err != nil {
        if isMySQLError(err, 1062) {
            return 0, fmt.Errorf("email %q already exists", email)
        }
        return 0, err
    }
    return result.LastInsertId()
}
```

---

## Summary

- MySQL's InnoDB engine provides ACID transactions, row-level locking, and MVCC.
- The clustered primary key physically orders rows — sequential integer PKs are much faster than random UUIDs.
- Use `utf8mb4` charset always (the `utf8` charset in MySQL is actually broken 3-byte UTF-8).
- `TIMESTAMP` columns have a 2038 problem — use `DATETIME` for new tables.
- MySQL 8.0 added CTEs, window functions, and better JSON support — always use 8.0+.
- The Go driver is `go-sql-driver/mysql`; always set `parseTime=true` in the DSN.

### Exercises

**Easy:** Set up MySQL with Docker. Create a `products` table and insert/query rows from Go using the `go-sql-driver/mysql`.

**Medium:** Write a Go function that handles MySQL's duplicate key error (1062) specifically, returning a friendly error message.

**Hard:** Configure MySQL binary log replication between two Docker containers. Write a Go program that writes to the primary and reads from the replica, measuring replication lag.
