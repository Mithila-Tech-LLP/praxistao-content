# Chapter 32: Building with ClickHouse in Go

A complete website analytics system in Go + ClickHouse that tracks page views, unique users, and funnel conversions in real-time.

## Table of Contents

1. Async Inserts — The Right Way to Write
2. Query Patterns for Analytics
3. Schema Migration in ClickHouse
4. Mini Project: Website Analytics System
5. Exercises

---

## 1. Async Inserts — The Right Way to Write

ClickHouse is optimized for bulk inserts, not single-row inserts. Inserting one row at a time creates millions of tiny parts on disk, degrading performance.

**The right pattern:** Buffer events in memory, flush periodically as batches.

```go
package analytics

import (
    "context"
    "log"
    "sync"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
    "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type PageView struct {
    Timestamp  time.Time
    UserID     string
    Page       string
    Referrer   string
    DurationMS uint32
    Country    string
}

type Writer struct {
    conn    driver.Conn
    buffer  []PageView
    mu      sync.Mutex
    flushCh chan struct{}
}

func NewWriter(conn driver.Conn, flushInterval time.Duration) *Writer {
    w := &Writer{
        conn:    conn,
        flushCh: make(chan struct{}, 1),
    }
    go w.autoFlush(flushInterval)
    return w
}

func (w *Writer) Write(pv PageView) {
    w.mu.Lock()
    w.buffer = append(w.buffer, pv)
    shouldFlush := len(w.buffer) >= 10000 // flush at 10k rows
    w.mu.Unlock()

    if shouldFlush {
        select {
        case w.flushCh <- struct{}{}:
        default:
        }
    }
}

func (w *Writer) autoFlush(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for {
        select {
        case <-ticker.C:
            w.Flush(context.Background())
        case <-w.flushCh:
            w.Flush(context.Background())
        }
    }
}

func (w *Writer) Flush(ctx context.Context) {
    w.mu.Lock()
    if len(w.buffer) == 0 {
        w.mu.Unlock()
        return
    }
    batch := make([]PageView, len(w.buffer))
    copy(batch, w.buffer)
    w.buffer = w.buffer[:0]
    w.mu.Unlock()

    if err := w.insertBatch(ctx, batch); err != nil {
        log.Printf("clickhouse flush error: %v", err)
    } else {
        log.Printf("flushed %d events to ClickHouse", len(batch))
    }
}

func (w *Writer) insertBatch(ctx context.Context, pvs []PageView) error {
    stmt, err := w.conn.PrepareBatch(ctx, "INSERT INTO page_views")
    if err != nil {
        return err
    }

    for _, pv := range pvs {
        if err := stmt.Append(
            pv.Timestamp,
            pv.UserID,
            pv.Page,
            pv.Referrer,
            pv.DurationMS,
            pv.Country,
        ); err != nil {
            return err
        }
    }
    return stmt.Send()
}
```

---

## 2. Query Patterns for Analytics

```go
type Analytics struct{ conn driver.Conn }

type DayStat struct {
    Date        time.Time
    Views       uint64
    UniqueUsers uint64
    AvgDuration float64
}

func (a *Analytics) DailyStats(ctx context.Context, days int) ([]DayStat, error) {
    rows, err := a.conn.Query(ctx, `
        SELECT
            toDate(timestamp) AS day,
            count() AS views,
            countDistinct(user_id) AS unique_users,
            avg(duration_ms) AS avg_duration
        FROM page_views
        WHERE timestamp >= today() - ?
        GROUP BY day
        ORDER BY day
    `, days)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []DayStat
    for rows.Next() {
        var s DayStat
        if err := rows.Scan(&s.Date, &s.Views, &s.UniqueUsers, &s.AvgDuration); err != nil {
            return nil, err
        }
        stats = append(stats, s)
    }
    return stats, rows.Err()
}

type PageStat struct {
    Page        string
    Views       uint64
    UniqueUsers uint64
    AvgDuration float64
}

func (a *Analytics) TopPages(ctx context.Context, days int, limit int) ([]PageStat, error) {
    rows, err := a.conn.Query(ctx, `
        SELECT
            page,
            count() AS views,
            countDistinct(user_id) AS unique_users,
            avg(duration_ms) AS avg_duration
        FROM page_views
        WHERE timestamp >= today() - ?
        GROUP BY page
        ORDER BY views DESC
        LIMIT ?
    `, days, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []PageStat
    for rows.Next() {
        var s PageStat
        rows.Scan(&s.Page, &s.Views, &s.UniqueUsers, &s.AvgDuration)
        stats = append(stats, s)
    }
    return stats, rows.Err()
}

type HourStat struct {
    Hour  time.Time
    Views uint64
}

func (a *Analytics) HourlyTrend(ctx context.Context, hours int) ([]HourStat, error) {
    rows, err := a.conn.Query(ctx, `
        SELECT
            toStartOfHour(timestamp) AS hour,
            count() AS views
        FROM page_views
        WHERE timestamp >= now() - INTERVAL ? HOUR
        GROUP BY hour
        ORDER BY hour
    `, hours)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []HourStat
    for rows.Next() {
        var s HourStat
        rows.Scan(&s.Hour, &s.Views)
        stats = append(stats, s)
    }
    return stats, rows.Err()
}
```

---

## 3. Schema Migration in ClickHouse

```go
var schema = []string{
    `CREATE TABLE IF NOT EXISTS page_views (
        timestamp   DateTime64(3) CODEC(Delta, ZSTD),
        user_id     String CODEC(ZSTD),
        page        LowCardinality(String),
        referrer    String CODEC(ZSTD),
        duration_ms UInt32 CODEC(Delta, ZSTD),
        country     LowCardinality(String) DEFAULT ''
    ) ENGINE = MergeTree()
    PARTITION BY toYYYYMM(timestamp)
    ORDER BY (toDate(timestamp), page, user_id)
    TTL toDate(timestamp) + INTERVAL 90 DAY`,

    `CREATE TABLE IF NOT EXISTS daily_stats (
        date         Date,
        page         LowCardinality(String),
        views        UInt64,
        unique_users UInt64,
        avg_duration Float64
    ) ENGINE = SummingMergeTree((views, unique_users))
    ORDER BY (date, page)`,
}

func Migrate(conn driver.Conn) error {
    ctx := context.Background()
    for _, stmt := range schema {
        if err := conn.Exec(ctx, stmt); err != nil {
            return fmt.Errorf("migrate: %w", err)
        }
    }
    return nil
}
```

`LowCardinality(String)` is a ClickHouse optimization for columns with < 10,000 distinct values (like country names, page paths). It uses dictionary encoding and is 3-5x faster for GROUP BY.

`CODEC(Delta, ZSTD)` compresses timestamps and sequential numbers much better than default compression.

---

## 4. Mini Project: Website Analytics System

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "strings"
    "time"

    "github.com/ClickHouse/clickhouse-go/v2"
)

var (
    writer    *Writer
    analytics *Analytics
)

func main() {
    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{"localhost:9000"},
    })
    if err != nil {
        log.Fatal(err)
    }

    if err := Migrate(conn); err != nil {
        log.Fatal("migrate:", err)
    }

    writer = NewWriter(conn, 2*time.Second)
    analytics = &Analytics{conn: conn}

    // Simulate traffic
    go simulateTraffic()

    mux := http.NewServeMux()
    mux.HandleFunc("POST /track", handleTrack)
    mux.HandleFunc("GET /analytics/daily", handleDailyStats)
    mux.HandleFunc("GET /analytics/pages", handleTopPages)
    mux.HandleFunc("GET /analytics/hourly", handleHourly)

    log.Println("Analytics API on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleTrack(w http.ResponseWriter, r *http.Request) {
    var pv PageView
    if err := json.NewDecoder(r.Body).Decode(&pv); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    pv.Timestamp = time.Now()
    writer.Write(pv)
    w.WriteHeader(204)
}

func handleDailyStats(w http.ResponseWriter, r *http.Request) {
    stats, err := analytics.DailyStats(r.Context(), 30)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(stats)
}

func handleTopPages(w http.ResponseWriter, r *http.Request) {
    stats, err := analytics.TopPages(r.Context(), 30, 20)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(stats)
}

func handleHourly(w http.ResponseWriter, r *http.Request) {
    stats, err := analytics.HourlyTrend(r.Context(), 24)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(stats)
}

func simulateTraffic() {
    pages := []string{"/home", "/products", "/blog", "/about", "/contact"}
    countries := []string{"US", "UK", "DE", "FR", "JP", "IN"}

    ticker := time.NewTicker(10 * time.Millisecond)
    for range ticker.C {
        writer.Write(PageView{
            Timestamp:  time.Now(),
            UserID:     fmt.Sprintf("user_%d", rand.Intn(10000)),
            Page:       pages[rand.Intn(len(pages))],
            Referrer:   "https://google.com",
            DurationMS: uint32(500 + rand.Intn(5000)),
            Country:    countries[rand.Intn(len(countries))],
        })
    }
}
```

Test:
```bash
# Track a page view
curl -X POST localhost:8080/track \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user123","page":"/home","duration_ms":1200,"country":"US"}'

# Get analytics
curl localhost:8080/analytics/daily
curl localhost:8080/analytics/pages
curl localhost:8080/analytics/hourly
```

---

## Summary

- Always batch-insert into ClickHouse. Use an in-memory buffer that flushes periodically (2 seconds) or at a size threshold (10,000 rows).
- `LowCardinality(String)` for columns with few distinct values (country, status, page) — 3-5x faster GROUP BY.
- `CODEC(Delta, ZSTD)` for timestamps and sequences — better compression.
- `SummingMergeTree` automatically aggregates during background merges — perfect for pre-computed daily/hourly stats.

### Exercises

**Easy:** Add a `GET /analytics/countries` endpoint that returns the top 10 countries by page views.

**Medium:** Add a "funnel" query: what % of users who viewed `/products` also viewed `/cart`, and then `/checkout`? Use ClickHouse's `windowFunnel` or `sequenceCount` functions.

**Hard:** Add a materialized view `hourly_page_views` using `SummingMergeTree` that pre-aggregates hourly stats. Query from both the raw table and the materialized view and compare performance on 100 million rows.
