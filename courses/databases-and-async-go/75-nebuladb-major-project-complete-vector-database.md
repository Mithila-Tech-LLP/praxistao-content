# Chapter 75: NebulaDB — Major Project: Complete Vector Database

You've built every component. This chapter assembles them into a fully working NebulaDB, builds a real-world application on top of it, and compares the result to Qdrant. Then we reflect on what you've learned.

## Table of Contents

1. The Complete NebulaDB
2. Project: Semantic Code Search Engine
3. Deploying NebulaDB with Docker
4. What You'd Add to Make NebulaDB Production-Grade
5. NebulaDB vs Qdrant vs pgvector — When to Use What
6. Exercises — The Final Challenge

---

## 1. The Complete NebulaDB

Here's the final wiring in `collection/collection.go` — all components integrated:

```go
// collection/collection.go — final version
package collection

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "sync/atomic"

    "nebuladb/hnsw"
    "nebuladb/index"
    "nebuladb/storage"
    "nebuladb/types"
)

type Collection struct {
    name         string
    dir          string
    config       Config
    hnswIndex    *hnsw.Index
    vectorStore  *storage.VectorStore
    payloadStore *storage.PayloadStore
    indexMgr     *index.Manager
    wal          *storage.WAL
    count        atomic.Int64
    mu           sync.RWMutex
}

// Info returns metadata about the collection for GET /collections/{name}
func (c *Collection) Info() map[string]any {
    return map[string]any{
        "name":        c.name,
        "dimension":   c.config.Dimension,
        "distance":    c.config.Distance,
        "point_count": c.count.Load(),
        "hnsw":        c.config.HNSW,
    }
}

// Upsert inserts or updates a point (WAL → stores → HNSW)
func (c *Collection) Upsert(point types.Point) error {
    if len(point.Vector) != c.config.Dimension {
        return fmt.Errorf("vector dimension %d != collection dimension %d",
            len(point.Vector), c.config.Dimension)
    }
    if err := c.wal.WriteUpsert(point); err != nil {
        return fmt.Errorf("wal: %w", err)
    }
    if err := c.vectorStore.Set(point.ID, point.Vector); err != nil {
        return fmt.Errorf("vectors: %w", err)
    }
    if err := c.payloadStore.Set(point.ID, point.Payload); err != nil {
        return fmt.Errorf("payload: %w", err)
    }
    c.indexMgr.Index(point.ID, point.Payload)
    c.hnswIndex.Insert(point.ID, point.Vector)
    c.count.Add(1)
    return nil
}

// GetPoint retrieves a point's vector and payload
func (c *Collection) GetPoint(id uint64) (*types.Point, error) {
    vec, err := c.vectorStore.Get(id)
    if err != nil {
        return nil, fmt.Errorf("point %d not found", id)
    }
    payload, _ := c.payloadStore.Get(id)
    return &types.Point{ID: id, Vector: vec, Payload: payload}, nil
}

// DeletePoint removes a point (lazy delete from HNSW)
func (c *Collection) DeletePoint(id uint64) error {
    if err := c.wal.WriteDelete(id); err != nil {
        return err
    }
    payload, _ := c.payloadStore.Get(id)
    if payload != nil {
        c.indexMgr.Remove(id, payload)
    }
    c.payloadStore.Delete(id)
    c.vectorStore.Delete(id)
    c.hnswIndex.Delete(id)
    c.count.Add(-1)
    return nil
}

// Search performs vector similarity search with optional filtering
func (c *Collection) Search(req types.SearchRequest) ([]types.ScoredPoint, error) {
    if len(req.Vector) != c.config.Dimension {
        return nil, fmt.Errorf("query dimension %d != collection dimension %d",
            len(req.Vector), c.config.Dimension)
    }
    if req.Limit <= 0 {
        req.Limit = 10
    }

    results, err := c.search(req.Vector, req.Filter, req.Limit)
    if err != nil {
        return nil, err
    }

    scored := make([]types.ScoredPoint, 0, len(results))
    for _, r := range results {
        sp := types.ScoredPoint{
            Point: types.Point{ID: r.ID},
            Score: r.Score,
        }
        if req.WithPayload {
            sp.Payload, _ = c.payloadStore.Get(r.ID)
        }
        scored = append(scored, sp)
    }
    return scored, nil
}

// Scroll iterates through points in ID order (for export/migration)
func (c *Collection) Scroll(offset uint64, limit int) ([]types.Point, *uint64) {
    return c.vectorStore.Scan(offset, limit, c.payloadStore)
}

// CreateFieldIndex registers a payload field for indexing
func (c *Collection) CreateFieldIndex(field string) {
    c.indexMgr.CreateIndex(field)
    // Re-index all existing points for this field
    c.payloadStore.ForEach(func(id uint64, payload types.Payload) {
        c.indexMgr.Index(id, payload)
    })
}
```

---

## 2. Project: Semantic Code Search Engine

Build a code search engine that finds functions by what they *do*, not by their name.

**How it works:**
1. Parse a Go repository, extract function signatures + docstrings + body summaries
2. Embed each function with a text embedding model
3. Store in NebulaDB with payload: `{file, name, signature, package, lines}`
4. HTTP server: query with natural language → find semantically similar functions

```go
// cmd/codesearch/main.go
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "nebuladb/collection"
    "nebuladb/server"
    "nebuladb/types"
)

type FunctionInfo struct {
    File      string `json:"file"`
    Package   string `json:"package"`
    Name      string `json:"name"`
    Signature string `json:"signature"`
    Doc       string `json:"doc"`
    StartLine int    `json:"start_line"`
    EndLine   int    `json:"end_line"`
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: codesearch <repo-path>")
    }
    repoPath := os.Args[1]

    mgr, err := collection.NewCollectionManager("./codesearch_data")
    if err != nil {
        log.Fatal(err)
    }
    defer mgr.Close()

    // Create collection (1536 dims for text-embedding-3-small)
    mgr.Create("functions", collection.Config{
        Dimension: 1536,
        Distance:  types.Cosine,
    })
    c, _ := mgr.Get("functions")

    // Create payload indexes for filtering
    c.CreateFieldIndex("package")
    c.CreateFieldIndex("file")

    // Index all Go functions in the repository
    log.Printf("Indexing functions in %s...", repoPath)
    indexed := indexRepo(repoPath, c)
    log.Printf("Indexed %d functions", indexed)

    // Start HTTP server
    srv := server.New(mgr, ":7000")
    log.Println("Code search running on :7000")
    log.Fatal(srv.Start())
}

func indexRepo(repoPath string, c *collection.Collection) int {
    var id uint64
    filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil || !strings.HasSuffix(path, ".go") {
            return err
        }

        fset := token.NewFileSet()
        f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
        if err != nil {
            return nil
        }

        for _, decl := range f.Decls {
            fd, ok := decl.(*ast.FuncDecl)
            if !ok {
                continue
            }

            doc := ""
            if fd.Doc != nil {
                doc = fd.Doc.Text()
            }

            sig := funcSignature(fd)
            text := fmt.Sprintf("Function: %s\nPackage: %s\nSignature: %s\nDoc: %s",
                fd.Name.Name, f.Name.Name, sig, doc)

            embedding, err := getEmbedding(text)
            if err != nil {
                log.Printf("embedding error for %s: %v", fd.Name.Name, err)
                continue
            }

            rel, _ := filepath.Rel(repoPath, path)
            payload := encodePayload(FunctionInfo{
                File:      rel,
                Package:   f.Name.Name,
                Name:      fd.Name.Name,
                Signature: sig,
                Doc:       doc,
                StartLine: fset.Position(fd.Pos()).Line,
                EndLine:   fset.Position(fd.End()).Line,
            })

            c.Upsert(types.Point{ID: id, Vector: embedding, Payload: payload})
            id++
        }
        return nil
    })
    return int(id)
}

func funcSignature(fd *ast.FuncDecl) string {
    var buf strings.Builder
    buf.WriteString("func ")
    if fd.Recv != nil {
        buf.WriteString("(receiver) ")
    }
    buf.WriteString(fd.Name.Name)
    buf.WriteString("(")
    if fd.Type.Params != nil {
        buf.WriteString("...params...")
    }
    buf.WriteString(")")
    return buf.String()
}

func encodePayload(fn FunctionInfo) types.Payload {
    p := types.Payload{}
    fields := map[string]any{
        "file": fn.File, "package": fn.Package, "name": fn.Name,
        "signature": fn.Signature, "doc": fn.Doc,
        "start_line": fn.StartLine, "end_line": fn.EndLine,
    }
    for k, v := range fields {
        data, _ := json.Marshal(v)
        p[k] = data
    }
    return p
}

func getEmbedding(text string) ([]float32, error) {
    body, _ := json.Marshal(map[string]any{
        "model": "text-embedding-3-small",
        "input": text,
    })
    req, _ := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Data []struct {
            Embedding []float32 `json:"embedding"`
        } `json:"data"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    if len(result.Data) == 0 {
        return nil, fmt.Errorf("no embeddings")
    }
    return result.Data[0].Embedding, nil
}
```

**Query:**

```bash
# Search for functions related to "parse HTTP request body"
curl -X POST localhost:7000/collections/functions/points/search \
  -H "Content-Type: application/json" \
  -d '{
    "vector": <embedding of "parse HTTP request body">,
    "limit": 5,
    "with_payload": true,
    "filter": {
      "must": [{"field": "package", "match": {"value": "server"}}]
    }
  }'
```

This finds relevant functions even if they're named `ReadBody`, `DecodeRequest`, or `parseJSON`.

---

## 3. Deploying NebulaDB with Docker

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o nebuladb ./cmd/nebuladb

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/nebuladb .
VOLUME /data
ENV NEBULADB_DATA_DIR=/data
EXPOSE 6380
CMD ["./nebuladb"]
```

```yaml
# docker-compose.yml
services:
  nebuladb:
    build: .
    ports:
      - "6380:6380"
    volumes:
      - nebuladb_data:/data
    environment:
      NEBULADB_DATA_DIR: /data
    restart: unless-stopped

volumes:
  nebuladb_data:
```

```bash
docker compose up -d

# Create a collection
curl -X POST localhost:6380/collections \
  -H "Content-Type: application/json" \
  -d '{"name":"test","dimension":128,"distance":"Cosine"}'
```

---

## 4. What You'd Add to Make NebulaDB Production-Grade

NebulaDB is a learning project. Qdrant adds many features to become production-ready:

**Quantization:**
Store vectors as int8 (scalar) or binary. 4-32x less RAM. We covered the concept in Ch 67.

**Distributed mode:**
Shard collections across nodes. Replicate for fault tolerance. Requires Raft consensus for metadata, gossip for health, consistent hashing for shard assignment.

**Named vectors:**
Multiple embeddings per point (text + image). NebulaDB could add this by making `VectorStore` hold a `map[string][]float32` per point.

**gRPC API:**
REST is convenient but slow for high-throughput. Qdrant uses gRPC (Protocol Buffers) which is 3-5x more efficient for vector payloads. Add with `google.golang.org/grpc`.

**Memory-mapped HNSW:**
Our HNSW stores nodes in a Go `map`. Qdrant uses mmap'd files with a fixed-size node layout. This means the index doesn't need to fit in RAM — the OS pages it.

**Sparse vector support:**
BM25/SPLADE keyword embeddings for hybrid search. Requires a different index structure (inverted index, not HNSW).

**On-disk vector storage:**
Our `VectorStore` loads the entire id→offset map into memory (~8 bytes/point). For 100M points, that's 800 MB just for the index. Qdrant uses a sorted array on disk + binary search.

---

## 5. NebulaDB vs Qdrant vs pgvector — When to Use What

| Scenario | Recommendation | Why |
|----------|---------------|-----|
| Already using PostgreSQL, <10M vectors | pgvector | Zero new infrastructure |
| New project, <100M vectors, production | Qdrant | Purpose-built, fast, easy ops |
| Learning how vector databases work | NebulaDB | You built it, you understand it |
| Multi-billion vectors, strict SLA | Milvus or Qdrant cluster | Distributed, battle-tested |
| Offline/embedded (no server) | Chroma (Python) or LanceDB | Embedded, no infra |
| Need SQL + vector in one query | pgvector | `SELECT ... WHERE ... ORDER BY embedding <=> ...` |

---

## 6. Exercises — The Final Challenge

**Easy:** Add `GET /collections/{name}/points/count` that returns `{"count": N}` — the number of points in the collection.

**Medium:** Implement collection **aliasing**: `PUT /aliases` accepts `{"alias": "search", "collection": "products_v2"}`. `GET /collections/search/points/search` routes to `products_v2`. This lets you swap collections without changing clients — the same trick Qdrant uses for zero-downtime reindexing.

**Hard (Production Challenge):** Implement **online reindexing**: a user changes HNSW parameters (e.g., M from 16 to 32 for better recall). Instead of taking NebulaDB offline, build a new index in the background while the old one serves traffic. Once the new index is complete, atomically swap to it. Hint: use an `atomic.Pointer[hnsw.Index]` for the live index pointer and build the new one in a goroutine. What happens to writes that arrive during the rebuild?

---

## NebulaDB: What You Built

```
7 chapters. ~3,000 lines of Go.

  HNSW index       → the graph algorithm at the core of every vector DB
  Payload storage  → BoltDB-backed JSON with B-Tree field indexes
  Filter engine    → evaluates complex conditions, auto-selects strategy
  Vector store     → flat binary file with in-memory offset index
  WAL              → CRC32-protected append-only log for crash recovery
  Snapshots        → gzip'd tar archive for backups
  HTTP API         → REST interface with logging middleware

Understanding the internals means you can:
- Debug slow queries (which strategy did the filter engine choose?)
- Tune HNSW parameters with confidence (not just trial and error)
- Read Qdrant's source code and understand every line
- Design the right schema for your use case (indexed fields, distance metric)
- Interview for database engineering roles
```

---

## Summary

- NebulaDB is a working vector database: create collections, upsert points with payloads, search with filters, crash-safe, HTTP API, Docker deployable.
- The semantic code search project shows a real application: parse a Go codebase, embed functions, search by natural language intent.
- Qdrant adds quantization, distribution, gRPC, and memory-mapped HNSW to the same architecture.
- Use pgvector when you already have Postgres. Use Qdrant for production AI applications. Use NebulaDB to understand what both of them do under the hood.

**Congratulations on finishing the NebulaDB project — and the full databases course.**
