# Chapter 66: ShopMetrics — Integrating VaultDB + StreamFlow into a Production Service

Part 1 built VaultDB. Part 2 built StreamFlow. Part 3 combines them into **ShopMetrics** — a real-time analytics service for an e-commerce platform.

ShopMetrics receives order events via StreamFlow, stores them in VaultDB, and exposes a query API. It demonstrates exactly how a production service integrates a custom database with a custom message broker.

## Table of Contents

1. What ShopMetrics Does
2. System Design
3. The VaultDB Schema
4. Receiving Events from StreamFlow
5. The Query API
6. Graceful Shutdown
7. Exercises

---

## 1. What ShopMetrics Does

An e-commerce company needs to answer:
- "How many orders were placed today?"
- "What's the total revenue in the last hour?"
- "Which user has spent the most this month?"
- "What's the average order value by region?"

These queries run against a stream of `order.placed` events flowing through StreamFlow.

**ShopMetrics:**
1. Subscribes to the "orders" topic in StreamFlow
2. Writes each order to VaultDB
3. Exposes HTTP endpoints to query the data

---

## 2. System Design

```
┌──────────────────────────────────────────────────────────────────┐
│                        ShopMetrics Service                        │
│                                                                    │
│   ┌─────────────────┐        ┌────────────────────────────────┐  │
│   │  StreamFlow     │        │        HTTP Query API          │  │
│   │  Consumer       │        │  GET /metrics/today            │  │
│   │  (orders topic) │──────► │  GET /metrics/revenue?hours=1  │  │
│   │                 │        │  GET /metrics/top-users        │  │
│   └────────┬────────┘        └──────────────┬─────────────────┘  │
│            │                                │                      │
│            ▼                                ▼                      │
│   ┌────────────────────────────────────────────────────────────┐  │
│   │                       VaultDB                              │  │
│   │   Table: orders                                            │  │
│   │   (order_id, user_id, amount, region, created_at)         │  │
│   └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 3. Project Structure

```
shopmetrics/
├── main.go
├── db/
│   └── store.go      ← VaultDB wrapper
├── ingest/
│   └── consumer.go   ← StreamFlow consumer
├── api/
│   └── server.go     ← HTTP query API
└── go.mod
```

---

## 4. VaultDB Store

```go
// db/store.go
package db

import (
    "fmt"
    "time"

    vaultdb "github.com/yourname/vaultdb"
)

type OrderRecord struct {
    OrderID   string
    UserID    string
    Amount    float64
    Region    string
    CreatedAt time.Time
}

type DailyStats struct {
    Date        string
    OrderCount  int
    TotalRevenue float64
    AvgOrderValue float64
}

type TopUser struct {
    UserID    string
    TotalSpend float64
    OrderCount int
}

type Store struct {
    db *vaultdb.DB
}

func Open(dir string) (*Store, error) {
    db, err := vaultdb.Open(dir)
    if err != nil {
        return nil, fmt.Errorf("vaultdb open: %w", err)
    }

    s := &Store{db: db}
    if err := s.createSchema(); err != nil {
        return nil, err
    }
    return s, nil
}

func (s *Store) createSchema() error {
    _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS orders (
        order_id   TEXT,
        user_id    TEXT,
        amount     FLOAT,
        region     TEXT,
        created_at TEXT
    )`)
    return err
}

func (s *Store) InsertOrder(o OrderRecord) error {
    _, err := s.db.Exec(
        `INSERT INTO orders (order_id, user_id, amount, region, created_at) VALUES (?, ?, ?, ?, ?)`,
        o.OrderID, o.UserID, o.Amount, o.Region, o.CreatedAt.Format(time.RFC3339),
    )
    return err
}

func (s *Store) DailyStats(date string) (*DailyStats, error) {
    rows, err := s.db.Query(
        `SELECT COUNT(*), SUM(amount) FROM orders WHERE created_at LIKE ?`,
        date+"%",
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    stats := &DailyStats{Date: date}
    if rows.Next() {
        rows.Scan(&stats.OrderCount, &stats.TotalRevenue)
        if stats.OrderCount > 0 {
            stats.AvgOrderValue = stats.TotalRevenue / float64(stats.OrderCount)
        }
    }
    return stats, nil
}

func (s *Store) RevenueLastN(hours int) (float64, error) {
    cutoff := time.Now().Add(-time.Duration(hours) * time.Hour).Format(time.RFC3339)
    rows, err := s.db.Query(
        `SELECT SUM(amount) FROM orders WHERE created_at >= ?`, cutoff)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    var total float64
    if rows.Next() {
        rows.Scan(&total)
    }
    return total, nil
}

func (s *Store) TopUsers(limit int) ([]TopUser, error) {
    rows, err := s.db.Query(
        `SELECT user_id, SUM(amount), COUNT(*) FROM orders GROUP BY user_id ORDER BY SUM(amount) DESC LIMIT ?`,
        limit,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []TopUser
    for rows.Next() {
        var u TopUser
        rows.Scan(&u.UserID, &u.TotalSpend, &u.OrderCount)
        users = append(users, u)
    }
    return users, nil
}

func (s *Store) Close() error {
    return s.db.Close()
}
```

---

## 5. StreamFlow Consumer (Ingestion)

```go
// ingest/consumer.go
package ingest

import (
    "context"
    "encoding/json"
    "log"
    "time"

    sfclient "github.com/yourname/streamflow/client"
    "github.com/yourname/shopmetrics/db"
)

type OrderEvent struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    Region    string    `json:"region"`
    OccurredAt time.Time `json:"occurred_at"`
}

type Consumer struct {
    store    *db.Store
    consumer *sfclient.Consumer
}

func New(brokerAddr string, store *db.Store) (*Consumer, error) {
    c, err := sfclient.NewConsumer(sfclient.ConsumerConfig{
        Addr:       brokerAddr,
        GroupID:    "shopmetrics",
        Topic:      "orders",
        Partition:  0,
        AutoCommit: false,
        MaxBytes:   2 * 1024 * 1024,
    })
    if err != nil {
        return nil, err
    }
    return &Consumer{store: store, consumer: c}, nil
}

func (c *Consumer) Run(ctx context.Context) {
    log.Println("[ingest] starting order consumer")

    for {
        records, err := c.consumer.Poll(ctx)
        if err != nil {
            if ctx.Err() != nil {
                log.Println("[ingest] stopping")
                return
            }
            log.Printf("[ingest] poll error: %v", err)
            time.Sleep(time.Second)
            continue
        }

        var lastOffset int64 = -1
        processedCount := 0

        for _, r := range records {
            var event OrderEvent
            if err := json.Unmarshal(r.Value, &event); err != nil {
                log.Printf("[ingest] bad event at offset %d: %v", r.Offset, err)
                lastOffset = r.Offset
                continue
            }

            order := db.OrderRecord{
                OrderID:   event.OrderID,
                UserID:    event.UserID,
                Amount:    event.Amount,
                Region:    event.Region,
                CreatedAt: event.OccurredAt,
            }

            if err := c.store.InsertOrder(order); err != nil {
                log.Printf("[ingest] insert error: %v", err)
                // Don't advance offset — retry this batch
                break
            }

            lastOffset = r.Offset
            processedCount++
        }

        if lastOffset >= 0 {
            c.consumer.CommitOffset(lastOffset + 1)
            log.Printf("[ingest] processed %d orders, committed offset %d",
                processedCount, lastOffset+1)
        }
    }
}

func (c *Consumer) Close() error {
    return c.consumer.Close()
}
```

---

## 6. HTTP Query API

```go
// api/server.go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "github.com/yourname/shopmetrics/db"
)

type Server struct {
    store *db.Store
}

func New(store *db.Store) *Server {
    return &Server{store: store}
}

func (s *Server) Routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /metrics/today", s.handleToday)
    mux.HandleFunc("GET /metrics/revenue", s.handleRevenue)
    mux.HandleFunc("GET /metrics/top-users", s.handleTopUsers)
    mux.HandleFunc("GET /health", s.handleHealth)
    return mux
}

func (s *Server) handleToday(w http.ResponseWriter, r *http.Request) {
    date := time.Now().Format("2006-01-02")
    stats, err := s.store.DailyStats(date)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleRevenue(w http.ResponseWriter, r *http.Request) {
    hoursStr := r.URL.Query().Get("hours")
    hours, _ := strconv.Atoi(hoursStr)
    if hours <= 0 {
        hours = 1
    }

    revenue, err := s.store.RevenueLastN(hours)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "hours":   hours,
        "revenue": revenue,
    })
}

func (s *Server) handleTopUsers(w http.ResponseWriter, r *http.Request) {
    limitStr := r.URL.Query().Get("limit")
    limit, _ := strconv.Atoi(limitStr)
    if limit <= 0 || limit > 100 {
        limit = 10
    }

    users, err := s.store.TopUsers(limit)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(users)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

---

## 7. Final main.go

```go
// main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourname/shopmetrics/api"
    "github.com/yourname/shopmetrics/db"
    "github.com/yourname/shopmetrics/ingest"
)

func main() {
    brokerAddr := envOr("BROKER_ADDR", "localhost:9999")
    dataDir := envOr("DATA_DIR", "./data")
    httpAddr := envOr("HTTP_ADDR", ":8080")

    // Open VaultDB
    store, err := db.Open(dataDir)
    if err != nil {
        log.Fatalf("store: %v", err)
    }
    defer store.Close()

    // Start StreamFlow consumer
    consumer, err := ingest.New(brokerAddr, store)
    if err != nil {
        log.Fatalf("consumer: %v", err)
    }
    defer consumer.Close()

    ctx, cancel := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer cancel()

    go consumer.Run(ctx)

    // Start HTTP server
    srv := &http.Server{
        Addr:    httpAddr,
        Handler: api.New(store).Routes(),
    }

    go func() {
        log.Printf("ShopMetrics API on %s", httpAddr)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("http: %v", err)
        }
    }()

    <-ctx.Done()

    // Graceful HTTP shutdown
    shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutCancel()
    srv.Shutdown(shutCtx)

    log.Println("ShopMetrics shutdown complete")
}

func envOr(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

---

## 8. Running the Full Stack

```bash
# Terminal 1: start StreamFlow broker
cd streamflow && go run main.go

# Terminal 2: start ShopMetrics
cd shopmetrics && go run main.go

# Terminal 3: produce some order events
cat > produce_orders.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "time"

    sfclient "github.com/yourname/streamflow/client"
)

type OrderEvent struct {
    OrderID    string    `json:"order_id"`
    UserID     string    `json:"user_id"`
    Amount     float64   `json:"amount"`
    Region     string    `json:"region"`
    OccurredAt time.Time `json:"occurred_at"`
}

func main() {
    admin, _ := sfclient.NewAdminClient("localhost:9999")
    admin.CreateTopic("orders", 1)
    admin.Close()

    p, _ := sfclient.NewProducer(sfclient.ProducerConfig{Addr: "localhost:9999"})
    defer p.Close()

    regions := []string{"us-east", "us-west", "eu-west", "ap-south"}
    for i := 0; i < 200; i++ {
        order := OrderEvent{
            OrderID:    fmt.Sprintf("ord-%04d", i),
            UserID:     fmt.Sprintf("user-%03d", rand.Intn(50)),
            Amount:     float64(rand.Intn(500)) + rand.Float64(),
            Region:     regions[rand.Intn(len(regions))],
            OccurredAt: time.Now(),
        }
        data, _ := json.Marshal(order)
        p.SendSync(sfclient.ProducerRecord{Topic: "orders", Value: data})
        time.Sleep(10 * time.Millisecond)
    }
    fmt.Println("Produced 200 orders")
}
EOF
go run produce_orders.go

# Query the API
curl http://localhost:8080/metrics/today
curl http://localhost:8080/metrics/revenue?hours=1
curl "http://localhost:8080/metrics/top-users?limit=5"
```

---

## Summary

You've integrated two systems you built from scratch into a production service:
- **VaultDB** stores orders with full SQL query support
- **StreamFlow** delivers order events reliably with consumer group semantics
- **ShopMetrics** ties them together with an HTTP API

Every major tech company's analytics systems follow this same pattern: events flow through a message broker, land in a database, and get queried by an API. You've built the full stack yourself.

### Exercises

**Easy:** Add a `GET /metrics/by-region` endpoint that returns total revenue grouped by region.

**Medium:** Add an in-memory cache to the query API: results for `today` are cached for 5 seconds. After 5 seconds the cache expires and the next request refreshes it. This is how real analytics APIs avoid hammering the database.

**Hard — Final Major Project:** Build `shopmetrics-cli`, a terminal dashboard that polls `GET /metrics/today` and `GET /metrics/top-users` every 3 seconds and displays a live-updating table using the `github.com/charmbracelet/bubbletea` TUI library. The dashboard should show: orders today, revenue today, top 5 users, and a rolling message count from the last 10 poll intervals.
