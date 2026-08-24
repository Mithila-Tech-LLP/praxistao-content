# Chapter 69: NebulaDB — Introduction to Building a Vector Database from Scratch

You've built VaultDB — a traditional SQL database with B-Trees, WAL, and MVCC. Now we build **NebulaDB**: a Qdrant-inspired vector database from scratch in Go. By the end of this project, you'll understand exactly what Qdrant does internally.

## Table of Contents

1. What We're Building
2. NebulaDB Architecture Overview
3. Project Setup and Structure
4. Core Data Types
5. The Collection Manager
6. Configuration and Startup
7. Exercises

---

## 1. What We're Building

NebulaDB is a simplified but real vector database with:

| Feature | NebulaDB | Qdrant |
|---------|---------|--------|
| Collections (named vector groups) | ✅ | ✅ |
| Points (id + vector + payload) | ✅ | ✅ |
| HNSW index (built from scratch) | ✅ | ✅ |
| Payload filtering | ✅ | ✅ |
| Payload indexes (keyword, numeric) | ✅ | ✅ |
| WAL for durability | ✅ | ✅ |
| Snapshots | ✅ | ✅ |
| REST HTTP API | ✅ | ✅ |
| Quantization | ❌ | ✅ |
| Distributed mode | ❌ | ✅ |
| Named vectors | ❌ | ✅ |

The missing features are excluded to keep the codebase learnable. Everything else is production-quality.

---

## 2. NebulaDB Architecture Overview

```
HTTP API Layer (port 6380)
    ├── POST /collections              → create collection
    ├── POST /collections/{name}/points → upsert points
    ├── POST /collections/{name}/search → vector search
    ├── GET  /collections/{name}/points/{id} → get point
    └── DELETE /collections/{name}/points/{id} → delete point

Collection Manager
    ├── Collection "products"
    │   ├── HNSWIndex    (in-memory graph, persisted to disk)
    │   ├── VectorStore  (float32 arrays, mmap'd file)
    │   ├── PayloadStore (JSON objects, BoltDB)
    │   ├── PayloadIndex (keyword/numeric B-Trees)
    │   └── WAL          (write-ahead log file)
    └── Collection "articles"
        └── ... same structure ...

Persistence Layer
    data/
    ├── products/
    │   ├── hnsw.bin       (serialized HNSW graph)
    │   ├── vectors.bin    (flat float32 array)
    │   ├── payload.db     (BoltDB)
    │   └── wal.log        (write-ahead log)
    └── articles/
        └── ...
```

**Request flow for a vector search:**

```
1. HTTP POST /collections/products/search
2. Parse request: query_vector + filter + limit
3. Apply payload index → get candidate IDs (optional, for pre-filtering)
4. HNSW search with inline filter checking
5. Fetch payloads for top-k results
6. Return JSON response
```

---

## 3. Project Setup and Structure

```bash
mkdir nebuladb && cd nebuladb
go mod init nebuladb

go get go.etcd.io/bbolt    # BoltDB for payload storage
go get golang.org/x/sys    # for mmap
```

Directory layout:

```
nebuladb/
├── main.go
├── server/
│   ├── server.go         # HTTP handlers
│   └── middleware.go
├── collection/
│   ├── collection.go     # Collection struct (orchestrates everything)
│   ├── manager.go        # CollectionManager (owns all collections)
│   └── config.go         # CollectionConfig
├── hnsw/
│   ├── index.go          # HNSW index implementation
│   ├── graph.go          # Graph node and edge structures
│   └── distance.go       # Distance functions
├── storage/
│   ├── vector_store.go   # Float32 vector storage (mmap)
│   ├── payload_store.go  # BoltDB payload storage
│   └── wal.go            # Write-ahead log
├── index/
│   └── payload_index.go  # B-Tree payload indexes
└── types/
    └── types.go          # Core types: Point, Vector, Payload, Filter
```

---

## 4. Core Data Types

```go
// types/types.go
package types

import "encoding/json"

// PointID is either a uint64 or a UUID string
type PointID = uint64

// Vector is a float32 slice (the embedding)
type Vector = []float32

// Payload is arbitrary JSON metadata attached to a point
type Payload map[string]json.RawMessage

// Point is the fundamental unit stored in NebulaDB
type Point struct {
    ID      PointID `json:"id"`
    Vector  Vector  `json:"vector"`
    Payload Payload `json:"payload"`
}

// ScoredPoint is a search result with a similarity score
type ScoredPoint struct {
    Point
    Score float32 `json:"score"`
}

// Distance metric
type Distance string

const (
    Cosine    Distance = "Cosine"
    Euclidean Distance = "Euclidean"
    DotProduct Distance = "Dot"
)

// Filter represents a search filter on payload fields
type Filter struct {
    Must    []Condition `json:"must,omitempty"`
    MustNot []Condition `json:"must_not,omitempty"`
    Should  []Condition `json:"should,omitempty"`
}

// Condition is a single filter condition on one field
type Condition struct {
    Field string          `json:"field"`
    Match *MatchCondition `json:"match,omitempty"`
    Range *RangeCondition `json:"range,omitempty"`
}

type MatchCondition struct {
    Value any   `json:"value"`          // exact match: string, int, bool
    Any   []any `json:"any,omitempty"`  // match any of these values
}

type RangeCondition struct {
    Gt  *float64 `json:"gt,omitempty"`
    Gte *float64 `json:"gte,omitempty"`
    Lt  *float64 `json:"lt,omitempty"`
    Lte *float64 `json:"lte,omitempty"`
}

// SearchRequest is the body of a POST /search request
type SearchRequest struct {
    Vector      Vector  `json:"vector"`
    Filter      *Filter `json:"filter,omitempty"`
    Limit       int     `json:"limit"`
    WithPayload bool    `json:"with_payload"`
}

// UpsertRequest is the body of a POST /points request
type UpsertRequest struct {
    Points []Point `json:"points"`
}
```

---

## 5. The Collection Manager

The `CollectionManager` is the top-level orchestrator. It holds all collections and handles collection lifecycle.

```go
// collection/manager.go
package collection

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type CollectionManager struct {
    dataDir     string
    collections map[string]*Collection
    mu          sync.RWMutex
}

func NewCollectionManager(dataDir string) (*CollectionManager, error) {
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return nil, fmt.Errorf("create data dir: %w", err)
    }

    m := &CollectionManager{
        dataDir:     dataDir,
        collections: make(map[string]*Collection),
    }

    // Load existing collections from disk
    if err := m.loadExisting(); err != nil {
        return nil, err
    }

    return m, nil
}

func (m *CollectionManager) Create(name string, cfg Config) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.collections[name]; exists {
        return fmt.Errorf("collection %q already exists", name)
    }

    dir := filepath.Join(m.dataDir, name)
    c, err := newCollection(name, dir, cfg)
    if err != nil {
        return err
    }

    m.collections[name] = c
    return nil
}

func (m *CollectionManager) Get(name string) (*Collection, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    c, ok := m.collections[name]
    if !ok {
        return nil, fmt.Errorf("collection %q not found", name)
    }
    return c, nil
}

func (m *CollectionManager) Delete(name string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    c, ok := m.collections[name]
    if !ok {
        return fmt.Errorf("collection %q not found", name)
    }

    if err := c.Close(); err != nil {
        return err
    }
    delete(m.collections, name)
    return os.RemoveAll(filepath.Join(m.dataDir, name))
}

func (m *CollectionManager) List() []string {
    m.mu.RLock()
    defer m.mu.RUnlock()

    names := make([]string, 0, len(m.collections))
    for name := range m.collections {
        names = append(names, name)
    }
    return names
}

func (m *CollectionManager) loadExisting() error {
    entries, err := os.ReadDir(m.dataDir)
    if err != nil {
        return nil // empty data dir is fine
    }

    for _, e := range entries {
        if !e.IsDir() {
            continue
        }
        name := e.Name()
        dir := filepath.Join(m.dataDir, name)
        c, err := loadCollection(name, dir)
        if err != nil {
            return fmt.Errorf("load collection %q: %w", name, err)
        }
        m.collections[name] = c
    }
    return nil
}

func (m *CollectionManager) Close() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, c := range m.collections {
        c.Close()
    }
    return nil
}
```

---

## 6. Configuration and Startup

```go
// collection/config.go
package collection

import (
    "encoding/json"
    "os"
    "path/filepath"

    "nebuladb/types"
)

// Config holds the collection creation parameters
type Config struct {
    Dimension int            `json:"dimension"`
    Distance  types.Distance `json:"distance"`
    HNSW      HNSWConfig     `json:"hnsw"`
}

// HNSWConfig controls index quality vs build speed
type HNSWConfig struct {
    M              int `json:"m"`               // connections per node (default: 16)
    EfConstruction int `json:"ef_construction"` // build-time exploration (default: 200)
    EfSearch       int `json:"ef_search"`       // query-time exploration (default: 128)
}

func (cfg *HNSWConfig) withDefaults() HNSWConfig {
    out := *cfg
    if out.M == 0 {
        out.M = 16
    }
    if out.EfConstruction == 0 {
        out.EfConstruction = 200
    }
    if out.EfSearch == 0 {
        out.EfSearch = 128
    }
    return out
}

func saveConfig(dir string, cfg Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

func loadConfig(dir string) (Config, error) {
    data, err := os.ReadFile(filepath.Join(dir, "config.json"))
    if err != nil {
        return Config{}, err
    }
    var cfg Config
    return cfg, json.Unmarshal(data, &cfg)
}
```

`main.go`:

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "nebuladb/collection"
    "nebuladb/server"
)

func main() {
    dataDir := os.Getenv("NEBULADB_DATA_DIR")
    if dataDir == "" {
        dataDir = "./data"
    }

    mgr, err := collection.NewCollectionManager(dataDir)
    if err != nil {
        log.Fatal("init collection manager:", err)
    }
    defer mgr.Close()

    srv := server.New(mgr, ":6380")

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Println("NebulaDB listening on :6380")
        if err := srv.Start(); err != nil {
            log.Fatal(err)
        }
    }()

    <-stop
    log.Println("Shutting down NebulaDB...")
    srv.Stop()
}
```

---

## Summary

- NebulaDB is a Qdrant-inspired vector database with Collections, Points (id + vector + payload), HNSW index, payload filtering, WAL, and REST API — all built in Go.
- **Collections** group points with a fixed dimension and distance metric.
- **CollectionManager** owns all collections and handles lifecycle (create, get, delete, persist).
- The architecture separates concerns: HNSW index, vector storage, payload storage, and WAL are independent components composed by `Collection`.
- Core types (`Point`, `Filter`, `Condition`) are defined once and used everywhere — clean boundaries.

### Exercises

**Easy:** Draw the NebulaDB architecture as a box diagram. Which components communicate with each other? What data flows between them during a `upsert` vs a `search`?

**Medium:** Add a `CollectionStats` struct with fields `PointCount int64`, `IndexSize int64`, `PayloadSize int64`. Implement a `Stats()` method on `Collection` that returns this (with placeholder values for now). Hook it up to `GET /collections/{name}`.

**Hard:** Design an `UpdateConfig` operation that lets callers change `HNSWConfig.EfSearch` at runtime without requiring a collection recreation. What state needs to change? Does the index need to be rebuilt? Write the `UpdateConfig(cfg HNSWConfig) error` signature on `Collection` and explain your approach.
