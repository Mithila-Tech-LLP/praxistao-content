# Chapter 73: NebulaDB — HTTP API

NebulaDB needs a network interface. Qdrant exposes both gRPC (port 6334) and REST (port 6333). We'll build a clean REST API in Go — no frameworks, just the standard library.

## Table of Contents

1. API Design
2. Server Setup
3. Collection Endpoints
4. Points Endpoints
5. Search Endpoint
6. Error Handling
7. Integration Tests
8. Exercises

---

## 1. API Design

```
POST   /collections                    Create a collection
GET    /collections                    List collections
GET    /collections/{name}             Get collection info
DELETE /collections/{name}             Delete collection
POST   /collections/{name}/snapshots   Create snapshot

PUT    /collections/{name}/points      Upsert points (batch)
GET    /collections/{name}/points/{id} Get a single point
DELETE /collections/{name}/points/{id} Delete a point
POST   /collections/{name}/points/scroll Scroll through all points (paginated)

POST   /collections/{name}/points/search Search (vector similarity)
POST   /collections/{name}/index       Create payload field index
```

Request/response format: JSON. Status codes follow REST conventions (200, 201, 400, 404, 500).

---

## 2. Server Setup

```go
// server/server.go
package server

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"

    "nebuladb/collection"
)

type Server struct {
    mgr    *collection.CollectionManager
    server *http.Server
}

func New(mgr *collection.CollectionManager, addr string) *Server {
    s := &Server{mgr: mgr}

    mux := http.NewServeMux()
    s.registerRoutes(mux)

    s.server = &http.Server{
        Addr:         addr,
        Handler:      withLogging(mux),
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
    return s
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
    // Collections
    mux.HandleFunc("POST /collections", s.createCollection)
    mux.HandleFunc("GET /collections", s.listCollections)
    mux.HandleFunc("GET /collections/{name}", s.getCollection)
    mux.HandleFunc("DELETE /collections/{name}", s.deleteCollection)
    mux.HandleFunc("POST /collections/{name}/snapshots", s.createSnapshot)

    // Points
    mux.HandleFunc("PUT /collections/{name}/points", s.upsertPoints)
    mux.HandleFunc("GET /collections/{name}/points/{id}", s.getPoint)
    mux.HandleFunc("DELETE /collections/{name}/points/{id}", s.deletePoint)
    mux.HandleFunc("POST /collections/{name}/points/scroll", s.scrollPoints)

    // Search & index
    mux.HandleFunc("POST /collections/{name}/points/search", s.search)
    mux.HandleFunc("POST /collections/{name}/index", s.createIndex)
}

func (s *Server) Start() error {
    return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return s.server.Shutdown(ctx)
}

// withLogging wraps the handler to log each request
func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := &responseWriter{ResponseWriter: w, status: 200}
        next.ServeHTTP(rw, r)
        log.Printf("%s %s → %d (%s)", r.Method, r.URL.Path, rw.status, time.Since(start))
    })
}

type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}
```

Helper functions for consistent responses:

```go
// server/helpers.go
package server

import (
    "encoding/json"
    "net/http"
)

type apiResponse struct {
    Result any    `json:"result,omitempty"`
    Status string `json:"status"`
    Error  string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, result any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(apiResponse{Result: result, Status: "ok"})
}

func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(apiResponse{Status: "error", Error: msg})
}

func decodeJSON(r *http.Request, dst any) bool {
    return json.NewDecoder(r.Body).Decode(dst) == nil
}
```

---

## 3. Collection Endpoints

```go
// server/collections.go
package server

import (
    "net/http"

    "nebuladb/collection"
)

type createCollectionReq struct {
    Name      string          `json:"name"`
    Dimension int             `json:"dimension"`
    Distance  string          `json:"distance"`
    HNSW      collection.HNSWConfig `json:"hnsw"`
}

func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
    var req createCollectionReq
    if !decodeJSON(r, &req) {
        writeError(w, 400, "invalid JSON")
        return
    }
    if req.Name == "" {
        writeError(w, 400, "name is required")
        return
    }
    if req.Dimension <= 0 {
        writeError(w, 400, "dimension must be > 0")
        return
    }

    cfg := collection.Config{
        Dimension: req.Dimension,
        Distance:  types.Distance(req.Distance),
        HNSW:      req.HNSW,
    }
    if cfg.Distance == "" {
        cfg.Distance = types.Cosine
    }

    if err := s.mgr.Create(req.Name, cfg); err != nil {
        writeError(w, 409, err.Error())
        return
    }

    writeJSON(w, 201, map[string]any{"name": req.Name, "status": "created"})
}

func (s *Server) listCollections(w http.ResponseWriter, r *http.Request) {
    names := s.mgr.List()
    writeJSON(w, 200, map[string]any{"collections": names})
}

func (s *Server) getCollection(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }
    writeJSON(w, 200, c.Info())
}

func (s *Server) deleteCollection(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    if err := s.mgr.Delete(name); err != nil {
        writeError(w, 404, err.Error())
        return
    }
    writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) createSnapshot(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    snapPath, err := c.Snapshot("./snapshots")
    if err != nil {
        writeError(w, 500, err.Error())
        return
    }
    writeJSON(w, 200, map[string]string{"path": snapPath})
}
```

---

## 4. Points Endpoints

```go
// server/points.go
package server

import (
    "net/http"
    "strconv"

    "nebuladb/types"
)

func (s *Server) upsertPoints(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    var req types.UpsertRequest
    if !decodeJSON(r, &req) {
        writeError(w, 400, "invalid JSON")
        return
    }
    if len(req.Points) == 0 {
        writeError(w, 400, "points array is empty")
        return
    }

    for _, point := range req.Points {
        if err := c.Upsert(point); err != nil {
            writeError(w, 500, err.Error())
            return
        }
    }

    writeJSON(w, 200, map[string]any{
        "upserted": len(req.Points),
        "status":   "acknowledged",
    })
}

func (s *Server) getPoint(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    idStr := r.PathValue("id")

    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        writeError(w, 400, "invalid point id")
        return
    }

    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    point, err := c.GetPoint(id)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    writeJSON(w, 200, point)
}

func (s *Server) deletePoint(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    idStr := r.PathValue("id")

    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        writeError(w, 400, "invalid point id")
        return
    }

    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    if err := c.DeletePoint(id); err != nil {
        writeError(w, 404, err.Error())
        return
    }

    writeJSON(w, 200, map[string]string{"status": "deleted"})
}

type scrollRequest struct {
    Limit  int    `json:"limit"`
    Offset uint64 `json:"offset"`
}

func (s *Server) scrollPoints(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    var req scrollRequest
    decodeJSON(r, &req)
    if req.Limit == 0 {
        req.Limit = 100
    }

    points, nextOffset := c.Scroll(req.Offset, req.Limit)
    writeJSON(w, 200, map[string]any{
        "points":      points,
        "next_offset": nextOffset,
    })
}
```

---

## 5. Search Endpoint

```go
// server/search.go
package server

import (
    "net/http"

    "nebuladb/types"
)

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    var req types.SearchRequest
    if !decodeJSON(r, &req) {
        writeError(w, 400, "invalid JSON")
        return
    }
    if len(req.Vector) == 0 {
        writeError(w, 400, "vector is required")
        return
    }
    if req.Limit == 0 {
        req.Limit = 10
    }

    results, err := c.Search(req)
    if err != nil {
        writeError(w, 400, err.Error())
        return
    }

    writeJSON(w, 200, map[string]any{
        "results": results,
        "count":   len(results),
    })
}

func (s *Server) createIndex(w http.ResponseWriter, r *http.Request) {
    name := r.PathValue("name")
    c, err := s.mgr.Get(name)
    if err != nil {
        writeError(w, 404, err.Error())
        return
    }

    var req struct {
        FieldName string `json:"field_name"`
    }
    if !decodeJSON(r, &req) || req.FieldName == "" {
        writeError(w, 400, "field_name is required")
        return
    }

    c.CreateFieldIndex(req.FieldName)
    writeJSON(w, 200, map[string]string{"status": "created", "field": req.FieldName})
}
```

---

## 6. Integration Tests

```go
// server/server_test.go
package server_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "nebuladb/collection"
    "nebuladb/server"
)

func setup(t *testing.T) (*server.Server, func()) {
    t.Helper()
    dir := t.TempDir()
    mgr, err := collection.NewCollectionManager(dir)
    if err != nil {
        t.Fatal(err)
    }
    srv := server.New(mgr, ":0")
    return srv, func() { mgr.Close() }
}

func post(t *testing.T, handler http.Handler, path string, body any) *httptest.ResponseRecorder {
    t.Helper()
    data, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", path, bytes.NewReader(data))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)
    return w
}

func TestCreateAndSearch(t *testing.T) {
    srv, cleanup := setup(t)
    defer cleanup()

    handler := srv.Handler()

    // Create collection
    resp := post(t, handler, "/collections", map[string]any{
        "name":      "test",
        "dimension": 3,
        "distance":  "Cosine",
    })
    if resp.Code != 201 {
        t.Fatalf("create: %d %s", resp.Code, resp.Body)
    }

    // Upsert points
    resp = put(t, handler, "/collections/test/points", map[string]any{
        "points": []map[string]any{
            {"id": 1, "vector": []float32{1, 0, 0}, "payload": map[string]any{"name": "A"}},
            {"id": 2, "vector": []float32{0, 1, 0}, "payload": map[string]any{"name": "B"}},
            {"id": 3, "vector": []float32{0.9, 0.1, 0}, "payload": map[string]any{"name": "C"}},
        },
    })
    if resp.Code != 200 {
        t.Fatalf("upsert: %d %s", resp.Code, resp.Body)
    }

    // Search: query close to [1,0,0] should return A and C
    resp = post(t, handler, "/collections/test/points/search", map[string]any{
        "vector":       []float32{0.99, 0.01, 0},
        "limit":        2,
        "with_payload": true,
    })
    if resp.Code != 200 {
        t.Fatalf("search: %d %s", resp.Code, resp.Body)
    }

    var result struct {
        Result struct {
            Results []struct {
                ID    uint64          `json:"id"`
                Score float32         `json:"score"`
            } `json:"results"`
        } `json:"result"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    if len(result.Result.Results) != 2 {
        t.Fatalf("expected 2 results, got %d", len(result.Result.Results))
    }

    // First result should be ID 1 (closest to [1,0,0])
    if result.Result.Results[0].ID != 1 {
        t.Errorf("expected ID 1 first, got %d", result.Result.Results[0].ID)
    }
}

func TestFilteredSearch(t *testing.T) {
    srv, cleanup := setup(t)
    defer cleanup()

    handler := srv.Handler()

    post(t, handler, "/collections", map[string]any{
        "name": "filtered", "dimension": 3, "distance": "Cosine",
    })

    put(t, handler, "/collections/filtered/points", map[string]any{
        "points": []map[string]any{
            {"id": 1, "vector": []float32{1, 0, 0}, "payload": map[string]any{"cat": "A", "price": 10.0}},
            {"id": 2, "vector": []float32{0.9, 0.1, 0}, "payload": map[string]any{"cat": "B", "price": 200.0}},
            {"id": 3, "vector": []float32{0.95, 0.05, 0}, "payload": map[string]any{"cat": "A", "price": 50.0}},
        },
    })

    // Search with filter: category="A" AND price < 100
    resp := post(t, handler, "/collections/filtered/points/search", map[string]any{
        "vector": []float32{1, 0, 0},
        "limit":  10,
        "filter": map[string]any{
            "must": []map[string]any{
                {"field": "cat", "match": map[string]any{"value": "A"}},
                {"field": "price", "range": map[string]any{"lt": 100.0}},
            },
        },
    })

    var result struct {
        Result struct {
            Results []struct{ ID uint64 } `json:"results"`
        } `json:"result"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    // Should return IDs 1 and 3 (cat=A, price<100); not ID 2 (price=200)
    if len(result.Result.Results) != 2 {
        t.Errorf("expected 2 filtered results, got %d", len(result.Result.Results))
    }
}
```

---

## Summary

- The HTTP server uses Go 1.22's new `{name}` path parameters — no external router needed.
- All endpoints return a consistent `{"result": ..., "status": "ok"}` envelope.
- `writeJSON`/`writeError` helpers ensure every response has the right Content-Type and status code.
- `withLogging` middleware logs every request with method, path, status, and duration.
- Integration tests use `httptest.NewRecorder()` — no network, no ports, fast and isolated.

### Exercises

**Easy:** Add `GET /health` that returns `{"status": "ok", "collections": N}`. No authentication needed — this is for load balancers to check liveness.

**Medium:** Add request validation middleware that rejects requests with `Content-Type` other than `application/json` for POST/PUT endpoints. Return a 415 Unsupported Media Type error with a helpful message.

**Hard:** Add basic API key authentication. When `NEBULADB_API_KEY` environment variable is set, all write endpoints (POST/PUT/DELETE) require an `Authorization: Bearer <key>` header. Read endpoints remain public. Write a test that verifies: (1) reads work without auth, (2) writes fail without auth, (3) writes succeed with correct auth.
