# Chapter 31: ClickHouse — Analytics at Ludicrous Speed

ClickHouse is to analytics what Redis is to caching — absurdly fast. It can scan billions of rows per second on a single server. This chapter explains how it works and when to use it.

## Table of Contents

1. What ClickHouse Is and Why It's Blazing Fast
2. Columnar vs Row Storage
3. MergeTree — ClickHouse's Core Engine
4. ClickHouse SQL Dialect
5. Table Engines
6. Real-Time Ingestion
7. Building with ClickHouse in Go
8. Exercises

---

## 1. What ClickHouse Is and Why It's Blazing Fast

ClickHouse is an open-source column-oriented OLAP (Online Analytical Processing) database built by Yandex. It's designed for analytical queries on huge datasets.

**Benchmark comparison on a 100 million row dataset:**

| Query | PostgreSQL | ClickHouse |
|-------|-----------|------------|
| `COUNT(*) WHERE date > X` | 3.5 seconds | 0.01 seconds |
| `SUM(revenue) GROUP BY country` | 8.2 seconds | 0.05 seconds |
| `SELECT * WHERE user_id = X` | 0.001 seconds | 0.8 seconds |

ClickHouse is 100-1000x faster for analytical queries. But slower for point lookups (fetching one row). This is the column store trade-off.

**Use ClickHouse for:**
- Web analytics (user behavior, funnels, cohorts)
- Application performance monitoring
- Log aggregation and search
- Ad tech (impressions, clicks, conversions)
- Business intelligence and reporting

**Don't use ClickHouse for:**
- Transactional systems (no UPDATE/DELETE by primary key)
- Point lookups of single rows
- Systems requiring row-level ACID transactions

---

## 2. Columnar vs Row Storage

### Row Storage (PostgreSQL, MySQL)

Each row is stored together. Great for retrieving complete records.

```
Row 1: [id=1, name="Alice", country="US", revenue=100]
Row 2: [id=2, name="Bob",   country="UK", revenue=250]
Row 3: [id=3, name="Carol", country="US", revenue=75]
```

To compute `SUM(revenue) WHERE country='US'`:
- Read all 3 rows from disk (even though we only need country and revenue columns)
- Filter by country
- Sum revenue

### Columnar Storage (ClickHouse)

Each column is stored separately. Great for aggregations over many rows.

```
id column:      [1, 2, 3, ...]
name column:    ["Alice", "Bob", "Carol", ...]
country column: ["US", "UK", "US", ...]
revenue column: [100, 250, 75, ...]
```

To compute `SUM(revenue) WHERE country='US'`:
- Read only the `country` column (much less data)
- Find matching rows
- Read only those rows from the `revenue` column
- Sum — never touched `id` or `name` at all

Columnar storage also compresses much better: a column of the same type with similar values compresses 5-10x better than mixed row data.

---

## 3. MergeTree — ClickHouse's Core Engine

MergeTree is ClickHouse's primary table engine. It writes data in sorted "parts" on disk and periodically merges them (like LSM-tree in Cassandra/RocksDB).

```sql
CREATE TABLE page_views (
    date        Date,
    user_id     UInt64,
    page        String,
    session_id  String,
    duration_ms UInt32
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(date)      -- partition by month
ORDER BY (date, user_id)         -- primary sort key (determines sparse index)
SETTINGS index_granularity = 8192; -- one index entry per 8192 rows
```

**PARTITION BY:** ClickHouse organizes data into partitions. Queries on a specific month only scan that month's data.

**ORDER BY:** Data within a partition is sorted by this key. This creates a sparse index — ClickHouse knows which 8192-row granules to read for a given value range.

**When you query `WHERE date = '2024-01-15' AND user_id = 42`:** ClickHouse reads only the relevant partition and only the granules that could contain that user_id range. 99% of data is skipped.

---

## 4. ClickHouse SQL Dialect

ClickHouse speaks SQL with extensions:

```sql
-- Aggregate functions
SELECT
    toStartOfDay(timestamp) AS day,
    count() AS events,
    countDistinct(user_id) AS unique_users,
    avg(duration_ms) AS avg_duration,
    quantile(0.95)(duration_ms) AS p95_duration
FROM page_views
WHERE date >= '2024-01-01' AND date < '2024-02-01'
GROUP BY day
ORDER BY day;

-- Array functions
SELECT arrayJoin(tags) AS tag, count() AS cnt
FROM posts
GROUP BY tag
ORDER BY cnt DESC
LIMIT 20;

-- URL parsing
SELECT
    domain(url) AS site,
    count() AS hits
FROM access_logs
GROUP BY site
ORDER BY hits DESC
LIMIT 10;

-- Date/time functions
SELECT
    toStartOfHour(created_at) AS hour,
    count() AS signups
FROM users
WHERE created_at >= now() - INTERVAL 7 DAY
GROUP BY hour
ORDER BY hour;

-- Window functions
SELECT
    user_id,
    revenue,
    sum(revenue) OVER (PARTITION BY country ORDER BY date) AS cumulative
FROM orders;
```

### Materialized Views

Materialized views in ClickHouse update automatically as data is inserted:

```sql
-- Source table
CREATE TABLE events (
    timestamp DateTime,
    event_type String,
    user_id UInt64
) ENGINE = MergeTree() ORDER BY timestamp;

-- Materialized view: hourly event counts
CREATE MATERIALIZED VIEW event_counts_hourly
ENGINE = SummingMergeTree()
ORDER BY (hour, event_type)
POPULATE AS
SELECT
    toStartOfHour(timestamp) AS hour,
    event_type,
    count() AS cnt
FROM events
GROUP BY hour, event_type;

-- INSERT into events → automatically updates event_counts_hourly
-- Query the view (pre-aggregated, instant):
SELECT hour, event_type, sum(cnt) FROM event_counts_hourly
WHERE hour >= now() - INTERVAL 24 HOUR
GROUP BY hour, event_type;
```

---

## 5. Table Engines

| Engine | Use Case |
|--------|----------|
| `MergeTree` | Default. Sorted storage with sparse index. |
| `ReplicatedMergeTree` | MergeTree with automatic replication (production). |
| `SummingMergeTree` | Automatically sums numeric columns during merges. Great for pre-aggregation. |
| `AggregatingMergeTree` | Stores partial aggregation states. For complex materialized views. |
| `CollapsingMergeTree` | Handles updates by storing old + new versions, collapsing on merge. |
| `Kafka` | Read directly from Kafka topics (no consumer code needed). |
| `Memory` | In-memory table, lost on restart. Good for temporary data. |
| `Log` | Append-only, no index. For tiny tables. |

---

## 6. Real-Time Ingestion

### Via HTTP (simplest)

```bash
# Insert via HTTP API
echo "2024-01-15 10:00:00,alice,/home,1200" | \
  curl -s "http://localhost:8123/?query=INSERT+INTO+page_views+FORMAT+CSV" --data-binary @-
```

### Via Kafka Engine

```sql
-- ClickHouse reads from Kafka automatically
CREATE TABLE kafka_events (
    timestamp DateTime,
    user_id   UInt64,
    event     String
) ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'user-events',
    kafka_group_name = 'clickhouse',
    kafka_format = 'JSONEachRow';

-- Materialized view that moves Kafka data to a real table
CREATE MATERIALIZED VIEW events_mv TO events AS
SELECT timestamp, user_id, event FROM kafka_events;
```

---

## 7. Building with ClickHouse in Go

```bash
go get github.com/ClickHouse/clickhouse-go/v2
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func main() {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{"localhost:9000"},
        Auth: clickhouse.Auth{
            Database: "default",
            Username: "default",
            Password: "",
        },
        MaxOpenConns:     10,
        MaxIdleConns:     5,
        ConnMaxLifetime:  time.Hour,
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    if err := conn.Ping(ctx); err != nil {
        log.Fatal("ping:", err)
    }
    fmt.Println("Connected to ClickHouse!")

    // Create table
    conn.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS page_views (
            date        Date DEFAULT toDate(timestamp),
            timestamp   DateTime,
            user_id     UInt64,
            page        String,
            duration_ms UInt32
        ) ENGINE = MergeTree()
        PARTITION BY toYYYYMM(date)
        ORDER BY (date, user_id)
    `)

    // Bulk insert (critical for performance: batch, don't insert one-by-one)
    batchInsert(ctx, conn)

    // Analytical query
    analyticsQuery(ctx, conn)
}

func batchInsert(ctx context.Context, conn driver.Conn) {
    batch, err := conn.PrepareBatch(ctx, "INSERT INTO page_views")
    if err != nil {
        log.Fatal(err)
    }

    now := time.Now()
    pages := []string{"/home", "/products", "/about", "/blog"}

    for i := 0; i < 100000; i++ {
        ts := now.Add(-time.Duration(i) * time.Second)
        if err := batch.Append(
            ts.Format("2006-01-02"),     // date
            ts,                          // timestamp
            uint64(i%1000),             // user_id
            pages[i%len(pages)],        // page
            uint32(100+i%2000),         // duration_ms
        ); err != nil {
            log.Fatal(err)
        }
    }

    if err := batch.Send(); err != nil {
        log.Fatal("send batch:", err)
    }
    fmt.Println("Inserted 100,000 rows")
}

func analyticsQuery(ctx context.Context, conn driver.Conn) {
    rows, err := conn.Query(ctx, `
        SELECT
            page,
            count() AS views,
            countDistinct(user_id) AS unique_users,
            avg(duration_ms) AS avg_duration
        FROM page_views
        WHERE date >= today() - 7
        GROUP BY page
        ORDER BY views DESC
        LIMIT 10
    `)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    fmt.Printf("\n%-20s %10s %12s %12s\n", "Page", "Views", "Unique Users", "Avg Duration")
    fmt.Println(strings.Repeat("-", 58))

    for rows.Next() {
        var page string
        var views, unique uint64
        var avgDur float64
        if err := rows.Scan(&page, &views, &unique, &avgDur); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%-20s %10d %12d %12.0fms\n", page, views, unique, avgDur)
    }
}
```

---

## Summary

- ClickHouse is a columnar database optimized for analytical queries — 100-1000x faster than row stores for aggregations.
- MergeTree engine partitions by date/time and sorts by a primary key for efficient range scans.
- Materialized views update automatically on insert — use them to pre-aggregate data for common queries.
- Always use **batch inserts** — ClickHouse is optimized for large bulk inserts, not one-row-at-a-time.
- Not suitable for transactional workloads, point lookups, or frequent row-level updates.

### Exercises

**Easy:** Create a `page_views` table and insert 1,000 rows with random data. Run a query that counts page views per page and measures execution time.

**Medium:** Create a materialized view that pre-aggregates hourly page view counts by page. Compare query time for hourly stats: reading from raw table vs materialized view on 10 million rows.

**Hard:** Set up ClickHouse with a Kafka engine table that reads from a Kafka topic. Write a Go producer that publishes 100,000 events/second to Kafka and verify ClickHouse ingests them in near-real-time.
