# Chapter 61: Middleware

Middleware is the backbone of any HTTP service. It separates cross-cutting concerns — logging, authentication, rate limiting, tracing — from business logic. A middleware is just a function that wraps a handler: it runs code before, calls the next handler, then optionally runs code after. This chapter shows how to build production-quality middleware and compose them correctly.

## Table of Contents

1. [The Middleware Pattern](#1-the-middleware-pattern)
2. [Logging Middleware](#2-logging-middleware)
3. [Authentication Middleware](#3-authentication-middleware)
4. [Rate Limiting Middleware](#4-rate-limiting-middleware)
5. [CORS Middleware](#5-cors-middleware)
6. [Request ID and Tracing](#6-request-id-and-tracing)
7. [Panic Recovery](#7-panic-recovery)
8. [Composing Middleware Chains](#8-composing-middleware-chains)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Middleware Pattern

```go
// A middleware wraps an http.Handler:
type Middleware func(http.Handler) http.Handler

// Basic structure:
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before: preprocessing
        log.Println("before handler")

        next.ServeHTTP(w, r)  // Call the wrapped handler

        // After: post-processing
        log.Println("after handler")
    })
}

// Apply to a single handler:
http.Handle("/foo", MyMiddleware(myHandler))

// Apply to a mux (all routes):
mux := http.NewServeMux()
mux.HandleFunc("/foo", myHandler)
http.ListenAndServe(":8080", MyMiddleware(mux))
```

**Execution order:** middleware wraps are applied inside-out, but they execute outside-in.
```
Apply:   Outer(Middle(Inner(handler)))
Execute: Outer → Middle → Inner → handler → Inner-after → Middle-after → Outer-after
```

**Reading the status code after the handler runs** requires a response writer wrapper:
```go
type statusRecorder struct {
    http.ResponseWriter
    status  int
    written int64
}

func (r *statusRecorder) WriteHeader(status int) {
    r.status = status
    r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
    n, err := r.ResponseWriter.Write(b)
    r.written += int64(n)
    return n, err
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
    return &statusRecorder{ResponseWriter: w, status: http.StatusOK}
}
```

---

## 2. Logging Middleware

```go
package middleware

import (
    "log/slog"
    "net/http"
    "time"
)

// Logger logs method, path, status, duration, and bytes for every request.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            rec := newStatusRecorder(w)
            next.ServeHTTP(rec, r)

            logger.Info("http request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", rec.status,
                "duration", time.Since(start),
                "bytes", rec.written,
                "ip", r.RemoteAddr,
                "userAgent", r.UserAgent(),
                "requestID", RequestIDFromContext(r.Context()),
            )
        })
    }
}

// Structured log output:
// {"time":"...","level":"INFO","msg":"http request","method":"GET","path":"/users","status":200,"duration":"1.2ms"}
```

---

## 3. Authentication Middleware

```go
package middleware

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string

const userKey contextKey = "user"

type Claims struct {
    UserID int64
    Email  string
    Roles  []string
}

// Authenticate validates a Bearer token and stores claims in context.
// Token validation logic (JWT) is covered in Chapter 64.
func Authenticate(tokenValidator func(string) (*Claims, error)) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            auth := r.Header.Get("Authorization")
            if !strings.HasPrefix(auth, "Bearer ") {
                writeJSON(w, http.StatusUnauthorized, map[string]string{
                    "error": "missing or invalid Authorization header",
                })
                return
            }

            token := strings.TrimPrefix(auth, "Bearer ")
            claims, err := tokenValidator(token)
            if err != nil {
                writeJSON(w, http.StatusUnauthorized, map[string]string{
                    "error": "invalid token: " + err.Error(),
                })
                return
            }

            ctx := context.WithValue(r.Context(), userKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Optional auth — continues even without token:
func OptionalAuth(tokenValidator func(string) (*Claims, error)) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            auth := r.Header.Get("Authorization")
            if strings.HasPrefix(auth, "Bearer ") {
                token := strings.TrimPrefix(auth, "Bearer ")
                if claims, err := tokenValidator(token); err == nil {
                    ctx := context.WithValue(r.Context(), userKey, claims)
                    r = r.WithContext(ctx)
                }
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Extract claims from context:
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
    c, ok := ctx.Value(userKey).(*Claims)
    return c, ok
}

// Require specific roles:
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    roleSet := make(map[string]bool)
    for _, r := range roles { roleSet[r] = true }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := ClaimsFromContext(r.Context())
            if !ok {
                writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
                return
            }
            for _, role := range claims.Roles {
                if roleSet[role] {
                    next.ServeHTTP(w, r)
                    return
                }
            }
            writeJSON(w, http.StatusForbidden, map[string]string{"error": "insufficient permissions"})
        })
    }
}
```

---

## 4. Rate Limiting Middleware

```go
package middleware

import (
    "net/http"
    "strings"
    "sync"
    "time"
)

type tokenBucket struct {
    tokens     float64
    maxTokens  float64
    refillRate float64  // Tokens per second
    lastRefill time.Time
    mu         sync.Mutex
}

func newBucket(rps float64) *tokenBucket {
    return &tokenBucket{
        tokens:     rps,
        maxTokens:  rps,
        refillRate: rps,
        lastRefill: time.Now(),
    }
}

func (b *tokenBucket) allow() bool {
    b.mu.Lock()
    defer b.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(b.lastRefill).Seconds()
    b.tokens += elapsed * b.refillRate
    if b.tokens > b.maxTokens { b.tokens = b.maxTokens }
    b.lastRefill = now

    if b.tokens < 1 { return false }
    b.tokens--
    return true
}

// RateLimit limits requests per IP. rps = requests per second.
func RateLimit(rps float64) func(http.Handler) http.Handler {
    buckets := sync.Map{}

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            // Strip port from "IP:port":
            if idx := strings.LastIndex(ip, ":"); idx != -1 { ip = ip[:idx] }

            val, _ := buckets.LoadOrStore(ip, newBucket(rps))
            bucket := val.(*tokenBucket)

            if !bucket.allow() {
                w.Header().Set("Retry-After", "1")
                writeJSON(w, http.StatusTooManyRequests, map[string]string{
                    "error": "rate limit exceeded",
                })
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Periodic cleanup to avoid memory leak from abandoned IPs:
func startBucketCleanup(buckets *sync.Map, interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            buckets.Range(func(k, v any) bool {
                b := v.(*tokenBucket)
                b.mu.Lock()
                idle := time.Since(b.lastRefill)
                b.mu.Unlock()
                if idle > 5*time.Minute { buckets.Delete(k) }
                return true
            })
        }
    }()
}
```

---

## 5. CORS Middleware

```go
package middleware

import (
    "net/http"
    "strconv"
    "strings"
)

type CORSConfig struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
    AllowCredentials bool
    MaxAge          int
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
    allowedOrigins := make(map[string]bool)
    for _, o := range cfg.AllowedOrigins { allowedOrigins[o] = true }

    methods := strings.Join(cfg.AllowedMethods, ", ")
    headers := strings.Join(cfg.AllowedHeaders, ", ")
    maxAge := strconv.Itoa(cfg.MaxAge)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            w.Header().Add("Vary", "Origin")  // Response differs per origin — tell caches

            if origin != "" && (allowedOrigins["*"] || allowedOrigins[origin]) {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                if cfg.AllowCredentials {
                    w.Header().Set("Access-Control-Allow-Credentials", "true")
                }
            }

            // Preflight request:
            if r.Method == http.MethodOptions {
                w.Header().Set("Access-Control-Allow-Methods", methods)
                w.Header().Set("Access-Control-Allow-Headers", headers)
                if cfg.MaxAge > 0 {
                    w.Header().Set("Access-Control-Max-Age", maxAge)
                }
                w.WriteHeader(http.StatusNoContent)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// Common configuration:
var DefaultCORS = CORSConfig{
    AllowedOrigins:   []string{"https://app.example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
    AllowCredentials: true,
    MaxAge:           3600,
}
```

---

## 6. Request ID and Tracing

```go
package middleware

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "net/http"
)

const requestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Reuse if provided by upstream (proxy, load balancer):
        id := r.Header.Get("X-Request-ID")
        if id == "" { id = generateID() }

        w.Header().Set("X-Request-ID", id)
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func generateID() string {
    b := make([]byte, 8)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func RequestIDFromContext(ctx context.Context) string {
    id, _ := ctx.Value(requestIDKey).(string)
    return id
}
```

---

## 7. Panic Recovery

```go
package middleware

import (
    "log/slog"
    "net/http"
    "runtime/debug"
)

// Recover catches panics in handlers and returns a 500.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    logger.Error("panic recovered",
                        "panic", rec,
                        "stack", string(debug.Stack()),
                        "requestID", RequestIDFromContext(r.Context()),
                        "path", r.URL.Path,
                    )
                    // Headers may have already been sent — best effort:
                    http.Error(w, http.StatusText(http.StatusInternalServerError),
                        http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 8. Composing Middleware Chains

```go
package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    r := chi.NewRouter()

    // Global middleware — applies to ALL routes:
    r.Use(middleware.Recover(logger))
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger(logger))
    r.Use(middleware.CORS(middleware.DefaultCORS))
    r.Use(middleware.RateLimit(100))  // 100 req/sec per IP

    // Public routes — no auth:
    r.Group(func(r chi.Router) {
        r.Get("/health", healthHandler)
        r.Post("/auth/login", loginHandler)
        r.Post("/auth/register", registerHandler)
    })

    // Protected routes — require valid token:
    r.Group(func(r chi.Router) {
        r.Use(middleware.Authenticate(validateJWT))
        r.Route("/api/v1", func(r chi.Router) {
            r.Mount("/users", userHandler.Routes())

            // Admin-only:
            r.Group(func(r chi.Router) {
                r.Use(middleware.RequireRole("admin"))
                r.Mount("/admin", adminHandler.Routes())
            })
        })
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
    logger.Info("starting server", "addr", srv.Addr)
    srv.ListenAndServe()
}
```

**Middleware chain diagram:**
```
Request → Recover → RequestID → Logger → CORS → RateLimit → Authenticate → RequireRole → Handler
          ↑ panic                                                                         │
          └─────────────────────────────────────────────────────────────────────────────┘ Response
```

---

## Summary

- Middleware = `func(http.Handler) http.Handler` — wraps handlers to add behavior
- Execution is outside-in; apply middleware from most general (recovery) to most specific (auth)
- Use `statusRecorder` to capture the response status code after the handler runs
- Rate limiting: token bucket per IP in `sync.Map` with periodic cleanup
- CORS: handle `OPTIONS` preflight, set `Access-Control-*` headers based on config
- `RequestID`: inject at the boundary, propagate through context, include in all log lines
- Panic recovery must be the outermost (first applied) middleware

---

## Exercises

### Easy
1. Build a `ContentType` middleware that rejects requests with body (POST, PUT, PATCH) if `Content-Type` is not `application/json`. Return `415 Unsupported Media Type`.
2. Build a `CacheControl` middleware that sets `Cache-Control: no-store` on all responses by default. Apply `public, max-age=3600` for `GET /static/*` paths (check `r.URL.Path`).
3. Write a test for the `Logger` middleware using `httptest.NewRecorder`. Verify that after a request, the recorder has the correct status code and the log output (captured with a custom `slog.Handler`) contains `method`, `path`, and `status` fields.

### Medium
4. **Timeout middleware**: implement `Timeout(d time.Duration) func(http.Handler) http.Handler` that cancels the request context after d. Use `context.WithTimeout(r.Context(), d)` and inject it with `r.WithContext(ctx)`. Test it by making a handler that sleeps and verifying it returns `503 Service Unavailable` after the timeout.
5. **IP allowlist middleware**: `IPAllowlist(cidrs []string) func(http.Handler) http.Handler` that blocks requests from IPs not in any of the given CIDR ranges. Parse CIDRs with `net.ParseCIDR`. Return `403 Forbidden` for blocked IPs. Handle `X-Forwarded-For` header for proxied requests.
6. **Retry middleware for HTTP client**: build a `RetryTransport` that wraps `http.RoundTripper`. On 5xx responses or network errors, retry up to 3 times with exponential backoff (100ms, 200ms, 400ms). Don't retry 4xx or successful responses. Use `context.Context` deadline to abort early if the caller's timeout expires.

### Hard
7. **Distributed rate limiting with Redis**: modify the token bucket rate limiter to use Redis (`go get github.com/redis/go-redis/v9`) instead of an in-memory map. Use a Lua script for atomic check-and-decrement: load tokens count, compute refill, clamp, check, decrement. This makes the rate limit work correctly across multiple server instances.
8. **Middleware metrics**: build `Metrics(reg prometheus.Registerer) func(http.Handler) http.Handler` that records `http_requests_total{method, path, status}` (counter) and `http_request_duration_seconds{method, path}` (histogram). Use the normalized path pattern (e.g., `/users/{id}` not `/users/42`) as the path label — extract it from `chi.RouteContext(r.Context()).RoutePattern()`. Chapter 120 (Prometheus) builds on this.
