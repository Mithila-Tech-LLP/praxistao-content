# Chapter 30: Building with Cassandra in Go

Complete Go patterns for Cassandra: session management, prepared statements, batches, pagination, and an IoT pipeline project.

## Table of Contents

1. Production Session Setup
2. CRUD with Prepared Statements
3. Pagination with Cassandra
4. Batch Operations
5. Error Handling
6. Mini Project: IoT Sensor Pipeline
7. Exercises

---

## 1. Production Session Setup

```go
package cassandra

import (
    "fmt"
    "time"

    "github.com/gocql/gocql"
)

type Config struct {
    Hosts       []string
    Keyspace    string
    Username    string
    Password    string
    Consistency gocql.Consistency
}

func NewSession(cfg Config) (*gocql.Session, error) {
    cluster := gocql.NewCluster(cfg.Hosts...)
    cluster.Keyspace = cfg.Keyspace
    cluster.Consistency = cfg.Consistency
    cluster.ProtoVersion = 4
    cluster.ConnectTimeout = 10 * time.Second
    cluster.Timeout = 5 * time.Second

    if cfg.Username != "" {
        cluster.Authenticator = gocql.PasswordAuthenticator{
            Username: cfg.Username,
            Password: cfg.Password,
        }
    }

    // Connection pool: 2 connections per host
    cluster.NumConns = 2

    // Retry policy: retry twice on timeouts/unavailable
    cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
        Min:        time.Millisecond * 100,
        Max:        time.Second * 5,
        NumRetries: 3,
    }

    session, err := cluster.CreateSession()
    if err != nil {
        return nil, fmt.Errorf("create cassandra session: %w", err)
    }
    return session, nil
}
```

---

## 2. CRUD with Prepared Statements

Prepared statements parse the CQL once and cache the plan. Use them for any query that runs more than once.

```go
type SensorRepository struct {
    session *gocql.Session
    insert  *gocql.Query
    select_ *gocql.Query
    delete_ *gocql.Query
}

func NewSensorRepository(session *gocql.Session) (*SensorRepository, error) {
    r := &SensorRepository{session: session}
    var err error

    r.insert, err = session.Prepare(`
        INSERT INTO sensor_readings (sensor_id, recorded_at, temperature, humidity)
        VALUES (?, ?, ?, ?)
        USING TTL 2592000
    `) // 30 days TTL
    if err != nil {
        return nil, err
    }

    r.select_, err = session.Prepare(`
        SELECT sensor_id, recorded_at, temperature, humidity
        FROM sensor_readings
        WHERE sensor_id = ?
          AND recorded_at >= ?
          AND recorded_at <  ?
        ORDER BY recorded_at DESC
        LIMIT ?
    `)
    if err != nil {
        return nil, err
    }

    return r, nil
}

type Reading struct {
    SensorID   gocql.UUID
    RecordedAt time.Time
    Temp       float32
    Humidity   float32
}

func (r *SensorRepository) Insert(sensorID gocql.UUID, ts time.Time, temp, humidity float32) error {
    return r.insert.Bind(sensorID, ts, temp, humidity).Exec()
}

func (r *SensorRepository) GetReadings(sensorID gocql.UUID, from, to time.Time, limit int) ([]Reading, error) {
    iter := r.select_.Bind(sensorID, from, to, limit).Iter()
    defer iter.Close()

    var readings []Reading
    var rd Reading
    for iter.Scan(&rd.SensorID, &rd.RecordedAt, &rd.Temp, &rd.Humidity) {
        readings = append(readings, rd)
    }
    return readings, iter.Close()
}
```

---

## 3. Pagination with Cassandra

Cassandra uses **paging state** for pagination, not OFFSET (which would require scanning all previous rows).

```go
func (r *SensorRepository) GetPage(
    sensorID gocql.UUID,
    from, to time.Time,
    pageSize int,
    pagingState []byte, // nil for first page
) ([]Reading, []byte, error) {

    query := r.session.Query(`
        SELECT sensor_id, recorded_at, temperature
        FROM sensor_readings
        WHERE sensor_id = ?
          AND recorded_at >= ?
          AND recorded_at < ?
        ORDER BY recorded_at DESC
    `, sensorID, from, to).PageSize(pageSize)

    if pagingState != nil {
        query = query.PageState(pagingState)
    }

    iter := query.Iter()
    defer iter.Close()

    var readings []Reading
    var rd Reading
    for iter.Scan(&rd.SensorID, &rd.RecordedAt, &rd.Temp) {
        readings = append(readings, rd)
    }

    nextPageState := iter.PageState() // nil if no more pages
    return readings, nextPageState, iter.Close()
}
```

Client usage (HTTP handler example):

```go
func handleListReadings(w http.ResponseWriter, r *http.Request) {
    // Decode paging state from base64 query param
    var pagingState []byte
    if ps := r.URL.Query().Get("cursor"); ps != "" {
        pagingState, _ = base64.StdEncoding.DecodeString(ps)
    }

    readings, nextState, err := repo.GetPage(sensorID, from, to, 50, pagingState)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    resp := map[string]interface{}{
        "data": readings,
    }
    if nextState != nil {
        resp["next_cursor"] = base64.StdEncoding.EncodeToString(nextState)
    }
    json.NewEncoder(w).Encode(resp)
}
```

---

## 4. Batch Operations

Cassandra has three batch types:
- **Logged batch:** atomic across all statements (slower, avoid for performance)
- **Unlogged batch:** no atomicity guarantee, but fast. Use only for same-partition writes.
- **Counter batch:** for updating counter columns

```go
func (r *SensorRepository) BulkInsert(readings []Reading) error {
    // Unlogged batch: fastest for bulk inserts to different partitions
    batch := r.session.NewBatch(gocql.UnloggedBatch)

    for _, rd := range readings {
        batch.Query(
            "INSERT INTO sensor_readings (sensor_id, recorded_at, temperature, humidity) VALUES (?, ?, ?, ?)",
            rd.SensorID, rd.RecordedAt, rd.Temp, rd.Humidity,
        )
    }

    return r.session.ExecuteBatch(batch)
}

// Concurrent bulk insert: split into batches of 100
func (r *SensorRepository) BulkInsertConcurrent(readings []Reading) error {
    const batchSize = 100
    var wg sync.WaitGroup
    errs := make(chan error, len(readings)/batchSize+1)

    for i := 0; i < len(readings); i += batchSize {
        end := i + batchSize
        if end > len(readings) {
            end = len(readings)
        }
        chunk := readings[i:end]

        wg.Add(1)
        go func(chunk []Reading) {
            defer wg.Done()
            if err := r.BulkInsert(chunk); err != nil {
                errs <- err
            }
        }(chunk)
    }

    wg.Wait()
    close(errs)

    for err := range errs {
        if err != nil {
            return err
        }
    }
    return nil
}
```

---

## 5. Error Handling

```go
import "github.com/gocql/gocql"

func isTimeout(err error) bool {
    if err == nil {
        return false
    }
    _, ok := err.(*gocql.RequestErrWriteTimeout)
    if ok {
        return true
    }
    _, ok = err.(*gocql.RequestErrReadTimeout)
    return ok
}

func isUnavailable(err error) bool {
    _, ok := err.(*gocql.RequestErrUnavailable)
    return ok
}

func handleCassandraError(err error) error {
    if err == nil {
        return nil
    }
    if err == gocql.ErrNotFound {
        return nil // not an error, just no results
    }
    if isTimeout(err) {
        return fmt.Errorf("cassandra timeout — try again: %w", err)
    }
    if isUnavailable(err) {
        return fmt.Errorf("cassandra unavailable — not enough replicas: %w", err)
    }
    return fmt.Errorf("cassandra: %w", err)
}
```

---

## 6. Mini Project: IoT Sensor Pipeline

A complete pipeline: simulate sensors writing data → Cassandra stores it → API reads it.

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/gocql/gocql"
)

var session *gocql.Session

func main() {
    cluster := gocql.NewCluster("localhost")
    cluster.Keyspace = "iot"
    cluster.Consistency = gocql.Quorum

    var err error
    session, err = cluster.CreateSession()
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    // Create schema
    session.Query(`CREATE TABLE IF NOT EXISTS sensor_readings (
        sensor_id   UUID,
        recorded_at TIMESTAMP,
        temperature FLOAT,
        humidity    FLOAT,
        PRIMARY KEY (sensor_id, recorded_at)
    ) WITH CLUSTERING ORDER BY (recorded_at DESC)
    AND default_time_to_live = 2592000`).Exec()

    // Simulate 5 sensors writing every 100ms
    for i := 0; i < 5; i++ {
        sensorID := gocql.MustRandomUUID()
        go simulateSensor(sensorID, i)
    }

    // API
    http.HandleFunc("GET /sensors/{id}/readings", handleReadings)

    log.Println("IoT API on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func simulateSensor(sensorID gocql.UUID, index int) {
    ticker := time.NewTicker(100 * time.Millisecond)
    for range ticker.C {
        temp := 20.0 + float32(index*2) + float32(rand.Intn(10))/10
        humidity := 50.0 + float32(rand.Intn(20))

        if err := session.Query(
            "INSERT INTO sensor_readings (sensor_id, recorded_at, temperature, humidity) VALUES (?, ?, ?, ?)",
            sensorID, time.Now(), temp, humidity,
        ).Exec(); err != nil {
            log.Printf("sensor %s write error: %v", sensorID, err)
        }
    }
}

func handleReadings(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    sensorID, err := gocql.ParseUUID(idStr)
    if err != nil {
        http.Error(w, "invalid sensor ID", 400)
        return
    }

    // Last hour of readings
    from := time.Now().Add(-time.Hour)
    to := time.Now()

    iter := session.Query(`
        SELECT recorded_at, temperature, humidity
        FROM sensor_readings
        WHERE sensor_id = ?
          AND recorded_at >= ?
          AND recorded_at <= ?
        LIMIT 1000
    `, sensorID, from, to).Iter()

    type point struct {
        Time      time.Time `json:"time"`
        Temp      float32   `json:"temperature"`
        Humidity  float32   `json:"humidity"`
    }

    var points []point
    var ts time.Time
    var temp, humidity float32
    for iter.Scan(&ts, &temp, &humidity) {
        points = append(points, point{Time: ts, Temp: temp, Humidity: humidity})
    }
    if err := iter.Close(); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "sensor_id": sensorID,
        "count":     len(points),
        "readings":  points,
    })
}
```

---

## Summary

- Use **prepared statements** for all repeated queries in Cassandra — they parse once, execute many times.
- **Paging state** (not OFFSET) for pagination — Cassandra can't skip rows efficiently.
- **Unlogged batch** for bulk inserts to different partitions (fastest). Logged batch for atomic multi-statement operations (slower).
- Check for `gocql.ErrNotFound`, `RequestErrWriteTimeout`, `RequestErrUnavailable` — handle each differently.

### Exercises

**Easy:** Write a Go function that inserts 1,000 sensor readings using both individual inserts (with transaction) and unlogged batch. Compare the time taken.

**Medium:** Implement cursor-based pagination for the sensor readings API. Encode the paging state as base64 and pass it as a `cursor` query parameter.

**Hard:** Add a "latest reading per sensor" table: `sensor_latest (sensor_id, recorded_at, temperature)` updated on every insert using a Cassandra lightweight transaction (`IF` condition). This ensures only newer readings update the latest entry.
