# Chapter 83: Time-Series Databases — InfluxDB and TimescaleDB

A time-series database (TSDB) is optimized for data that arrives in order of time: metrics, sensor readings, financial ticks, IoT data. The core assumption — new data is always newer than existing data — enables aggressive compression and specialized query optimizations that general-purpose databases cannot match.

## Table of Contents

1. [When to Use a TSDB](#1-when-to-use-a-tsdb)
2. [InfluxDB — Line Protocol and Flux](#2-influxdb--line-protocol-and-flux)
3. [TimescaleDB — PostgreSQL Extension](#3-timescaledb--postgresql-extension)
4. [Patterns in Go](#4-patterns-in-go)
5. [Comparison and Choice Guide](#5-comparison-and-choice-guide)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. When to Use a TSDB

**Use a TSDB when:**
- Data is append-only and timestamped
- You query by time ranges and aggregate (avg, sum, rate, percentile)
- Data volume is high (millions of points per minute)
- You need automatic data downsampling/retention policies

**Stick with PostgreSQL when:**
- You have < 10M rows/month
- You need JOINs with other tables
- Strong ACID consistency matters
- You don't want to operate another database

---

## 2. InfluxDB — Line Protocol and Flux

InfluxDB is a purpose-built TSDB. Data is organized into **buckets**, written using **line protocol** (a compact text format), and queried with **Flux** (a functional query language).

### Data model

```
measurement,tag1=val1,tag2=val2 field1=v1,field2=v2 timestamp_ns

cpu,host=server01,region=us-east usage_idle=80.5,usage_user=15.2 1688000000000000000
```

- **Measurement**: like a table name (`cpu`, `http_requests`, `temperature`)
- **Tags**: indexed metadata, low cardinality (`host`, `region`, `status_code`)
- **Fields**: actual data values, not indexed (`value`, `duration_ms`, `error_count`)
- **Timestamp**: nanoseconds since Unix epoch

### Writing to InfluxDB from Go

```go
import (
    influxdb2 "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
    "github.com/influxdata/influxdb-client-go/v2/api/write"
)

type MetricsWriter struct {
    client   influxdb2.Client
    writeAPI api.WriteAPI
    org      string
    bucket   string
}

func NewMetricsWriter(url, token, org, bucket string) *MetricsWriter {
    client := influxdb2.NewClientWithOptions(url, token,
        influxdb2.DefaultOptions().
            SetBatchSize(5000).         // buffer 5000 points before flush
            SetFlushInterval(1000),     // flush every 1 second
    )
    return &MetricsWriter{
        client:   client,
        writeAPI: client.WriteAPI(org, bucket),
        org:      org,
        bucket:   bucket,
    }
}

// Record an HTTP request metric
func (w *MetricsWriter) RecordHTTP(method, path, status string, durationMs float64) {
    point := write.NewPoint("http_requests",
        map[string]string{
            "method": method,
            "path":   path,
            "status": status,
        },
        map[string]interface{}{
            "duration_ms": durationMs,
            "count":       1,
        },
        time.Now(),
    )
    w.writeAPI.WritePoint(point) // non-blocking, batched
}

// Record a system metric
func (w *MetricsWriter) RecordCPU(host string, usagePercent float64) {
    w.writeAPI.WriteRecord(fmt.Sprintf(
        "cpu,host=%s usage=%f %d",
        host, usagePercent, time.Now().UnixNano(),
    ))
}

func (w *MetricsWriter) Flush() {
    w.writeAPI.Flush()
}

func (w *MetricsWriter) Close() {
    w.client.Close()
}

// Handle write errors asynchronously
func (w *MetricsWriter) watchErrors(ctx context.Context) {
    errorsCh := w.writeAPI.Errors()
    for {
        select {
        case err := <-errorsCh:
            log.Printf("influx write error: %v", err)
        case <-ctx.Done():
            return
        }
    }
}
```

### Querying with Flux

```go
type QueryResult struct {
    Time   time.Time
    Value  float64
    Labels map[string]string
}

func (w *MetricsWriter) QueryP99Latency(ctx context.Context, hours int) ([]QueryResult, error) {
    queryAPI := w.client.QueryAPI(w.org)
    
    flux := fmt.Sprintf(`
        from(bucket: "%s")
          |> range(start: -%dh)
          |> filter(fn: (r) => r._measurement == "http_requests")
          |> filter(fn: (r) => r._field == "duration_ms")
          |> aggregateWindow(every: 5m, fn: (tables=<-, column) =>
               tables |> quantile(q: 0.99, column: column))
          |> yield(name: "p99")
    `, w.bucket, hours)
    
    result, err := queryAPI.Query(ctx, flux)
    if err != nil { return nil, fmt.Errorf("query: %w", err) }
    defer result.Close()
    
    var rows []QueryResult
    for result.Next() {
        record := result.Record()
        rows = append(rows, QueryResult{
            Time:  record.Time(),
            Value: record.Value().(float64),
            Labels: map[string]string{
                "path":   record.Field(),
                "status": fmt.Sprintf("%v", record.ValueByKey("status")),
            },
        })
    }
    return rows, result.Err()
}
```

---

## 3. TimescaleDB — PostgreSQL Extension

TimescaleDB extends PostgreSQL with time-series capabilities. You get all of PostgreSQL (SQL, JOINs, ACID, indexes) plus automatic time-based partitioning (hypertables), compression, and continuous aggregates.

```sql
-- Install extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create a hypertable (partitioned by time automatically)
CREATE TABLE metrics (
    time        TIMESTAMPTZ NOT NULL,
    host        TEXT NOT NULL,
    service     TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    value       DOUBLE PRECISION NOT NULL
);

-- Convert to hypertable: TimescaleDB partitions by 1-day chunks
SELECT create_hypertable('metrics', 'time', chunk_time_interval => INTERVAL '1 day');

-- Create indexes on commonly filtered columns
CREATE INDEX idx_metrics_host_time ON metrics (host, time DESC);
CREATE INDEX idx_metrics_service_name ON metrics (service, metric_name, time DESC);

-- Compression: compress chunks older than 7 days (often 10-20x smaller)
ALTER TABLE metrics SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'host,service',
    timescaledb.compress_orderby = 'time DESC'
);
SELECT add_compression_policy('metrics', INTERVAL '7 days');

-- Retention: auto-drop chunks older than 90 days
SELECT add_retention_policy('metrics', INTERVAL '90 days');
```

### Time-series queries

```sql
-- Recent values for a specific host
SELECT time, metric_name, value
FROM metrics
WHERE host = 'server01'
  AND time > NOW() - INTERVAL '1 hour'
ORDER BY time DESC;

-- time_bucket: equivalent to date_trunc but with arbitrary intervals
SELECT
    time_bucket('5 minutes', time) AS bucket,
    service,
    AVG(value) AS avg_value,
    MAX(value) AS max_value,
    percentile_disc(0.99) WITHIN GROUP (ORDER BY value) AS p99
FROM metrics
WHERE metric_name = 'cpu_usage'
  AND time > NOW() - INTERVAL '24 hours'
GROUP BY bucket, service
ORDER BY bucket;

-- Rate of change (using lag)
SELECT
    time,
    service,
    value,
    value - LAG(value) OVER (PARTITION BY service ORDER BY time) AS delta
FROM metrics
WHERE metric_name = 'request_count'
  AND time > NOW() - INTERVAL '1 hour';
```

### Continuous aggregates (materialized views refreshed automatically)

```sql
-- Pre-compute 1-minute averages automatically
CREATE MATERIALIZED VIEW metrics_1min
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 minute', time) AS bucket,
    host,
    service,
    metric_name,
    AVG(value)  AS avg_value,
    MAX(value)  AS max_value,
    MIN(value)  AS min_value,
    COUNT(*)    AS sample_count
FROM metrics
GROUP BY bucket, host, service, metric_name
WITH NO DATA;

-- Auto-refresh as new data arrives
SELECT add_continuous_aggregate_policy('metrics_1min',
    start_offset  => INTERVAL '10 minutes',
    end_offset    => INTERVAL '1 minute',
    schedule_interval => INTERVAL '1 minute'
);
```

### Using TimescaleDB from Go

```go
// TimescaleDB is just PostgreSQL — use the same pgx/sqlx drivers
type TimescaleStore struct {
    db *sqlx.DB
}

func (s *TimescaleStore) InsertBatch(ctx context.Context, points []MetricPoint) error {
    // Use COPY for maximum insert throughput
    tx, _ := s.db.BeginTx(ctx, nil)
    stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("metrics", "time", "host", "service", "metric_name", "value"))
    
    for _, p := range points {
        stmt.ExecContext(ctx, p.Time, p.Host, p.Service, p.Name, p.Value)
    }
    stmt.ExecContext(ctx) // flush
    stmt.Close()
    return tx.Commit()
}

type BucketedMetric struct {
    Bucket  time.Time `db:"bucket"`
    Service string    `db:"service"`
    Avg     float64   `db:"avg_value"`
    P99     float64   `db:"p99"`
}

func (s *TimescaleStore) QueryBucketed(ctx context.Context, metricName, service string, hours int) ([]BucketedMetric, error) {
    query := `
        SELECT
            time_bucket('5 minutes', time) AS bucket,
            service,
            AVG(value) AS avg_value,
            percentile_disc(0.99) WITHIN GROUP (ORDER BY value) AS p99
        FROM metrics
        WHERE metric_name = $1
          AND service = $2
          AND time > NOW() - ($3 || ' hours')::interval
        GROUP BY bucket, service
        ORDER BY bucket`
    
    var results []BucketedMetric
    err := s.db.SelectContext(ctx, &results, query, metricName, service, hours)
    return results, err
}
```

---

## 4. Patterns in Go

### Metrics middleware for HTTP handlers

```go
type MetricsMiddleware struct {
    writer *MetricsWriter
    host   string
}

func (m *MetricsMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        rec := &statusRecorder{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(rec, r)
        
        m.writer.RecordHTTP(
            r.Method,
            r.URL.Path,
            strconv.Itoa(rec.statusCode),
            float64(time.Since(start).Milliseconds()),
        )
    })
}

type statusRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}
```

---

## 5. Comparison and Choice Guide

| | InfluxDB | TimescaleDB |
|---|---|---|
| Base | Purpose-built TSDB | PostgreSQL extension |
| Query language | Flux (functional) | SQL + time_bucket |
| Existing PostgreSQL integration | No | Yes — same connection, JOINs work |
| Compression | Built-in | Built-in |
| Horizontal scaling | InfluxDB Clustered | Distributed via Timescale Cloud |
| Setup complexity | Moderate | Low if you already run PostgreSQL |
| Free tier | Limited | Yes (TimescaleDB OSS) |

**Choose TimescaleDB if**: you already run PostgreSQL and want minimal ops overhead. You can JOIN metrics to your `users` table directly in SQL.

**Choose InfluxDB if**: your team is comfortable with Flux, you need Telegraf/Kapacitor integration, or you're building a standalone metrics pipeline separate from your main DB.

---

## Summary

- TSDBs assume append-only, time-ordered data — this assumption drives massive compression gains
- **InfluxDB**: line protocol for writes, Flux for queries, native compression
- **TimescaleDB**: PostgreSQL + hypertables + `time_bucket()` + continuous aggregates
- Use `COPY` (PostgreSQL) or batched writes (InfluxDB async API) for high-throughput ingestion
- **Continuous aggregates** (TimescaleDB) or **tasks** (InfluxDB) for pre-computing rollups

## Exercises

### Easy
1. Set up TimescaleDB locally with Docker. Create a `cpu_usage` hypertable with fields: `time`, `host`, `value`. Insert 10,000 rows and run a `time_bucket` query to get 1-minute averages.
2. Write a Go struct `CPUMetric` and a function `BulkInsert(ctx, []CPUMetric)` that uses PostgreSQL `COPY` to insert 10,000 rows as fast as possible.
3. Create a continuous aggregate on the `cpu_usage` table for 1-minute buckets. Query both the raw table and the aggregate and compare query times.

### Medium
4. Build an HTTP metrics recorder: middleware that records `method`, `path`, `status`, `duration_ms` to InfluxDB. Query the last 30 minutes of P99 latency by path. Display results in a table.
5. Implement a **data retention tester**: insert one year of synthetic CPU metrics into TimescaleDB. Apply a 90-day retention policy, then verify that older chunks are automatically dropped. Measure query time before and after enabling compression.
6. Compare InfluxDB vs TimescaleDB insert throughput: write a benchmark that inserts 1 million metric points to each. Compare: total time, storage size, and query time for a 1-hour window query.

### Hard
7. Build a **real-time anomaly detector**: use TimescaleDB to compute rolling averages with `time_bucket_gapfill`. Detect when a metric exceeds 2× the rolling average of the previous hour. Fire a webhook when an anomaly is detected. Run this as a background goroutine using `LISTEN/NOTIFY`.
8. Implement **automatic downsampling**: store raw data at 1-second resolution. After 1 day, downsample to 1-minute averages and delete the raw data. After 1 week, downsample to 1-hour averages. Use TimescaleDB continuous aggregates and retention policies for the chain.
