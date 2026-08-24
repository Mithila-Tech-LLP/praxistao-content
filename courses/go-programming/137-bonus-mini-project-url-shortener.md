# Chapter 26: Mini Project 2 — URL Shortener API

Time to apply the advanced Go concepts from Chapters 18–25. You'll build a production-quality URL shortener with an HTTP API, in-memory storage (with optional file persistence), goroutine-safe concurrent access, rate limiting, and proper error handling. This is your first real HTTP backend — the foundation for everything in the web development chapters ahead.

## Project Overview

**What you'll build:** A REST API for shortening URLs.

```bash
# Create a short URL:
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com/search?q=golang", "custom_code": "google"}'
# Response: {"code": "google", "short_url": "http://localhost:8080/r/google"}

# Redirect (browser follows this):
curl -L http://localhost:8080/r/google
# → redirects to https://www.google.com/...

# Stats:
curl http://localhost:8080/stats/google
# Response: {"code": "google", "original_url": "...", "hits": 42, "created_at": "..."}

# List all:
curl http://localhost:8080/links
```

**What you'll practice:**
- HTTP server with `net/http`
- Goroutine-safe map with `sync.RWMutex`
- JSON encoding/decoding
- Error handling across layers
- Context usage for request scoping
- Rate limiting with channels
- Atomic counters for stats
- Table-driven tests with `httptest`

---

## Project Structure

```
urlshortener/
├── main.go           ← HTTP server, routing
├── store.go          ← Thread-safe link storage
├── handler.go        ← HTTP handlers
├── model.go          ← Data types
├── ratelimit.go      ← Token bucket rate limiter
├── handler_test.go   ← Handler tests using httptest
├── store_test.go     ← Store tests
└── go.mod
```

---

## Step 1: Data Model (`model.go`)

```go
// model.go
package main

import (
    "fmt"
    "time"
)

type Link struct {
    Code        string    `json:"code"`
    OriginalURL string    `json:"original_url"`
    ShortURL    string    `json:"short_url"`
    Hits        int64     `json:"hits"`  // atomic, so int64
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func (l *Link) IsExpired() bool {
    if l.ExpiresAt == nil {
        return false
    }
    return time.Now().After(*l.ExpiresAt)
}

type CreateRequest struct {
    URL        string  `json:"url"`
    CustomCode string  `json:"custom_code,omitempty"`
    TTLSeconds *int    `json:"ttl_seconds,omitempty"`
}

func (r *CreateRequest) Validate() error {
    if r.URL == "" {
        return fmt.Errorf("url is required")
    }
    if len(r.URL) > 2048 {
        return fmt.Errorf("url too long (max 2048 chars)")
    }
    if r.CustomCode != "" {
        if len(r.CustomCode) < 2 || len(r.CustomCode) > 32 {
            return fmt.Errorf("custom_code must be 2-32 characters")
        }
        for _, c := range r.CustomCode {
            if !isAlphanumeric(c) && c != '-' && c != '_' {
                return fmt.Errorf("custom_code may only contain letters, digits, - and _")
            }
        }
    }
    return nil
}

func isAlphanumeric(c rune) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

type CreateResponse struct {
    Code     string `json:"code"`
    ShortURL string `json:"short_url"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

---

## Step 2: Thread-Safe Store (`store.go`)

```go
// store.go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

var (
    ErrNotFound  = errors.New("link not found")
    ErrCodeTaken = errors.New("code already in use")
    ErrExpired   = errors.New("link has expired")
)

type Store struct {
    mu      sync.RWMutex
    links   map[string]*Link
    baseURL string
}

func NewStore(baseURL string) *Store {
    return &Store{
        links:   make(map[string]*Link),
        baseURL: baseURL,
    }
}

// Create adds a new link, generating a code if not provided.
func (s *Store) Create(req CreateRequest) (*Link, error) {
    code := req.CustomCode
    if code == "" {
        var err error
        code, err = generateCode(6)
        if err != nil {
            return nil, fmt.Errorf("generating code: %w", err)
        }
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.links[code]; exists {
        return nil, fmt.Errorf("code %q: %w", code, ErrCodeTaken)
    }

    link := &Link{
        Code:        code,
        OriginalURL: req.URL,
        ShortURL:    s.baseURL + "/r/" + code,
        CreatedAt:   time.Now(),
    }

    if req.TTLSeconds != nil && *req.TTLSeconds > 0 {
        exp := time.Now().Add(time.Duration(*req.TTLSeconds) * time.Second)
        link.ExpiresAt = &exp
    }

    s.links[code] = link
    return link, nil
}

// Get returns a link by code and increments its hit counter.
func (s *Store) Get(code string) (*Link, error) {
    s.mu.RLock()
    link, ok := s.links[code]
    s.mu.RUnlock()

    if !ok {
        return nil, fmt.Errorf("code %q: %w", code, ErrNotFound)
    }

    if link.IsExpired() {
        return nil, fmt.Errorf("code %q: %w", code, ErrExpired)
    }

    atomic.AddInt64(&link.Hits, 1)  // Thread-safe increment
    return link, nil
}

// Stats returns a link without incrementing hits.
func (s *Store) Stats(code string) (*Link, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    link, ok := s.links[code]
    if !ok {
        return nil, fmt.Errorf("code %q: %w", code, ErrNotFound)
    }
    return link, nil
}

// List returns all non-expired links.
func (s *Store) List() []*Link {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var result []*Link
    for _, link := range s.links {
        if !link.IsExpired() {
            result = append(result, link)
        }
    }
    return result
}

// Delete removes a link by code.
func (s *Store) Delete(code string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, ok := s.links[code]; !ok {
        return fmt.Errorf("code %q: %w", code, ErrNotFound)
    }

    delete(s.links, code)
    return nil
}

// generateCode creates a URL-safe random code.
func generateCode(length int) (string, error) {
    b := make([]byte, length)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    // URL-safe base64, trim padding, use only first `length` chars
    return base64.URLEncoding.EncodeToString(b)[:length], nil
}
```

---

## Step 3: Rate Limiter (`ratelimit.go`)

```go
// ratelimit.go
package main

import (
    "net/http"
    "sync"
    "time"
)

// IPRateLimiter limits requests per IP address using token bucket algorithm.
type IPRateLimiter struct {
    mu       sync.Mutex
    limiters map[string]*tokenBucket
    rate     int           // tokens per second
    burst    int           // max burst size
}

type tokenBucket struct {
    tokens     float64
    maxTokens  float64
    refillRate float64  // tokens per nanosecond
    lastRefill time.Time
}

func NewIPRateLimiter(ratePerSecond, burst int) *IPRateLimiter {
    return &IPRateLimiter{
        limiters: make(map[string]*tokenBucket),
        rate:     ratePerSecond,
        burst:    burst,
    }
}

func (l *IPRateLimiter) Allow(ip string) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    bucket, ok := l.limiters[ip]
    if !ok {
        bucket = &tokenBucket{
            tokens:     float64(l.burst),
            maxTokens:  float64(l.burst),
            refillRate: float64(l.rate) / 1e9, // tokens per nanosecond
            lastRefill: time.Now(),
        }
        l.limiters[ip] = bucket
    }

    // Refill tokens:
    now := time.Now()
    elapsed := now.Sub(bucket.lastRefill)
    bucket.tokens = min(bucket.maxTokens, bucket.tokens+float64(elapsed)*bucket.refillRate)
    bucket.lastRefill = now

    if bucket.tokens < 1 {
        return false
    }
    bucket.tokens--
    return true
}

func min(a, b float64) float64 {
    if a < b {
        return a
    }
    return b
}

// Middleware wraps a handler with rate limiting.
func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        if !l.Allow(ip) {
            http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## Step 4: HTTP Handlers (`handler.go`)

```go
// handler.go
package main

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strings"
)

type Handler struct {
    store *Store
}

func NewHandler(store *Store) *Handler {
    return &Handler{store: store}
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        log.Printf("writeJSON encode error: %v", err)
    }
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, ErrorResponse{Error: msg})
}

// HandleCreate handles POST /shorten
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    var req CreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
        return
    }

    if err := req.Validate(); err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    link, err := h.store.Create(req)
    if err != nil {
        if errors.Is(err, ErrCodeTaken) {
            writeError(w, http.StatusConflict, "code already in use")
            return
        }
        log.Printf("store.Create error: %v", err)
        writeError(w, http.StatusInternalServerError, "failed to create link")
        return
    }

    writeJSON(w, http.StatusCreated, CreateResponse{
        Code:     link.Code,
        ShortURL: link.ShortURL,
    })
}

// HandleRedirect handles GET /r/{code}
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
    // Extract code from path: /r/{code}
    code := strings.TrimPrefix(r.URL.Path, "/r/")
    if code == "" {
        writeError(w, http.StatusBadRequest, "missing code")
        return
    }

    link, err := h.store.Get(code)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            writeError(w, http.StatusNotFound, "link not found")
            return
        }
        if errors.Is(err, ErrExpired) {
            writeError(w, http.StatusGone, "link has expired")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

// HandleStats handles GET /stats/{code}
func (h *Handler) HandleStats(w http.ResponseWriter, r *http.Request) {
    code := strings.TrimPrefix(r.URL.Path, "/stats/")
    if code == "" {
        writeError(w, http.StatusBadRequest, "missing code")
        return
    }

    link, err := h.store.Stats(code)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            writeError(w, http.StatusNotFound, "link not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    writeJSON(w, http.StatusOK, link)
}

// HandleList handles GET /links
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
    links := h.store.List()
    if links == nil {
        links = []*Link{}  // Return empty array, not null
    }
    writeJSON(w, http.StatusOK, links)
}

// HandleDelete handles DELETE /links/{code}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }

    code := strings.TrimPrefix(r.URL.Path, "/links/")
    if err := h.store.Delete(code); err != nil {
        if errors.Is(err, ErrNotFound) {
            writeError(w, http.StatusNotFound, "link not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "internal error")
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

---

## Step 5: Server (`main.go`)

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
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:" + port
    }

    store := NewStore(baseURL)
    handler := NewHandler(store)
    limiter := NewIPRateLimiter(10, 30)  // 10 req/s, burst of 30

    mux := http.NewServeMux()
    mux.HandleFunc("/shorten", handler.HandleCreate)
    mux.HandleFunc("/r/", handler.HandleRedirect)
    mux.HandleFunc("/stats/", handler.HandleStats)
    mux.HandleFunc("/links", handler.HandleList)
    mux.HandleFunc("/links/", handler.HandleDelete)
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
    })

    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      limiter.Middleware(mux),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Graceful shutdown:
    done := make(chan os.Signal, 1)
    signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Printf("Server starting on %s (base URL: %s)", srv.Addr, baseURL)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe: %v", err)
        }
    }()

    <-done
    log.Println("Shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped cleanly")
}
```

---

## Step 6: Tests (`handler_test.go`)

```go
// handler_test.go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func setupTestServer(t *testing.T) (*Handler, *Store) {
    t.Helper()
    store := NewStore("http://test.example.com")
    handler := NewHandler(store)
    return handler, store
}

func postJSON(t *testing.T, handler http.HandlerFunc, body any) *http.Response {
    t.Helper()
    b, _ := json.Marshal(body)
    req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    handler(w, req)
    return w.Result()
}

func TestHandleCreate_Success(t *testing.T) {
    h, _ := setupTestServer(t)

    resp := postJSON(t, h.HandleCreate, CreateRequest{URL: "https://example.com"})

    if resp.StatusCode != http.StatusCreated {
        t.Errorf("status = %d; want 201", resp.StatusCode)
    }

    var result CreateResponse
    json.NewDecoder(resp.Body).Decode(&result)

    if result.Code == "" {
        t.Error("expected a code in response")
    }
    if result.ShortURL == "" {
        t.Error("expected a short_url in response")
    }
}

func TestHandleCreate_InvalidURL(t *testing.T) {
    h, _ := setupTestServer(t)

    resp := postJSON(t, h.HandleCreate, CreateRequest{URL: ""})

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("status = %d; want 400", resp.StatusCode)
    }
}

func TestHandleCreate_DuplicateCode(t *testing.T) {
    h, _ := setupTestServer(t)

    postJSON(t, h.HandleCreate, CreateRequest{URL: "https://a.com", CustomCode: "test"})
    resp := postJSON(t, h.HandleCreate, CreateRequest{URL: "https://b.com", CustomCode: "test"})

    if resp.StatusCode != http.StatusConflict {
        t.Errorf("status = %d; want 409", resp.StatusCode)
    }
}

func TestHandleRedirect_Success(t *testing.T) {
    h, store := setupTestServer(t)

    link, _ := store.Create(CreateRequest{URL: "https://go.dev", CustomCode: "gohome"})

    req := httptest.NewRequest(http.MethodGet, "/r/gohome", nil)
    w := httptest.NewRecorder()
    h.HandleRedirect(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusFound {
        t.Errorf("status = %d; want 302", resp.StatusCode)
    }
    if loc := resp.Header.Get("Location"); loc != link.OriginalURL {
        t.Errorf("Location = %q; want %q", loc, link.OriginalURL)
    }
}

func TestHandleRedirect_NotFound(t *testing.T) {
    h, _ := setupTestServer(t)

    req := httptest.NewRequest(http.MethodGet, "/r/nosuchcode", nil)
    w := httptest.NewRecorder()
    h.HandleRedirect(w, req)

    if w.Code != http.StatusNotFound {
        t.Errorf("status = %d; want 404", w.Code)
    }
}

func TestHandleStats_HitCount(t *testing.T) {
    h, store := setupTestServer(t)

    store.Create(CreateRequest{URL: "https://example.com", CustomCode: "mylink"})

    // Hit the link 3 times:
    for i := 0; i < 3; i++ {
        req := httptest.NewRequest(http.MethodGet, "/r/mylink", nil)
        w := httptest.NewRecorder()
        h.HandleRedirect(w, req)
    }

    // Check stats:
    req := httptest.NewRequest(http.MethodGet, "/stats/mylink", nil)
    w := httptest.NewRecorder()
    h.HandleStats(w, req)

    var link Link
    json.NewDecoder(w.Result().Body).Decode(&link)

    if link.Hits != 3 {
        t.Errorf("Hits = %d; want 3", link.Hits)
    }
}
```

---

## Running the Project

```bash
cd urlshortener
go mod init urlshortener
go run .

# Test it:
# Create a short URL:
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://go.dev/doc/effective_go","custom_code":"effectivego"}'

# Redirect:
curl -L http://localhost:8080/r/effectivego

# Stats:
curl http://localhost:8080/stats/effectivego

# List all:
curl http://localhost:8080/links

# Run tests:
go test ./... -v

# With race detector:
go test -race ./...
```

---

## What You Just Practiced

| Concept | Where used |
|---------|-----------|
| Goroutines | Graceful shutdown goroutine |
| Channels | `done` signal channel |
| Context | `context.WithTimeout` for shutdown |
| sync.RWMutex | Store: concurrent read/write safety |
| sync/atomic | `link.Hits` increment |
| Error wrapping | Throughout store and handlers |
| Sentinel errors | `ErrNotFound`, `ErrCodeTaken`, `ErrExpired` |
| errors.Is | Handler error type detection |
| JSON encode/decode | All handlers |
| httptest | All handler tests |
| Graceful shutdown | `signal.Notify` + `srv.Shutdown` |

---

## Extension Challenges

1. **Persistence**: On startup, load links from a JSON file. On each create/delete, write to file atomically (temp file + rename). Implement `Store.Save(path string)` and `Store.Load(path string)`.

2. **Redis backend**: Replace the in-memory map with Redis using `go-redis`. Store links as JSON strings with TTL support (`SET key value EX seconds`). Increment hits atomically with `INCR`.

3. **Analytics**: Track per-link stats: hits per day (last 7 days), referrer breakdown, geo-location from IP (use a free IP database like MaxMind GeoLite2). Return from `/stats/{code}`.

4. **Custom domains**: Allow `{"url": "...", "domain": "short.mycompany.com"}`. The server looks up which domain a request came from and resolves codes in that domain's namespace.

5. **Authentication**: Add API key authentication middleware. Keys are stored in a map. Protected routes: `POST /shorten`, `DELETE /links/{code}`. Public routes: redirects and stats.

6. **QR codes**: Add `GET /qr/{code}` that returns a PNG QR code for the short URL (use `github.com/skip2/go-qrcode`).
