# Chapter 68: Rate Limiting and Throttling in Go HTTP Services

Every public-facing API needs rate limiting. Without it, a single misbehaving client can exhaust your server, a misconfigured script can cause accidental DoS, and a bad actor can brute-force authentication. Rate limiting enforces fair usage, protects SLAs, and keeps your service alive.

## Table of Contents

1. [Why Rate Limiting?](#1-why-rate-limiting)
2. [Fixed Window Counter](#2-fixed-window-counter)
3. [Sliding Window Log](#3-sliding-window-log)
4. [Token Bucket](#4-token-bucket)
5. [Leaky Bucket](#5-leaky-bucket)
6. [Per-IP Middleware with Token Bucket](#6-per-ip-middleware-with-token-bucket)
7. [Distributed Rate Limiting with Redis](#7-distributed-rate-limiting-with-redis)
8. [`golang.org/x/time/rate`](#8-golangorgxtimerate)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Rate Limiting?

```
Without rate limiting:
  Client A sends 50,000 requests/sec → server OOMs
  Attacker enumerates /login → brute-forces passwords
  A retry bug in a microservice cascades into a thundering herd
  One premium client monopolizes capacity, degrading others

With rate limiting:
  Each client gets a fair share
  Abuse is capped before it becomes an outage
  SLA guarantees (e.g. "1000 req/min per API key") are enforceable
  Downstream services get predictable load
```

**The four main algorithms**, in order of increasing sophistication:

```
Fixed Window   → simplest, burst problem at boundary
Sliding Window → accurate, high memory
Token Bucket   → allows controlled bursts, most common
Leaky Bucket   → smooths output to constant rate
```

---

## 2. Fixed Window Counter

Divide time into fixed windows (e.g. 1-minute buckets). Count requests per window. Reject if over limit.

```go
package ratelimit

import (
    "sync"
    "time"
)

// FixedWindow allows up to limit requests per window duration.
type FixedWindow struct {
    mu       sync.Mutex
    limit    int
    window   time.Duration
    count    int
    windowStart time.Time
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
    return &FixedWindow{
        limit:       limit,
        window:      window,
        windowStart: time.Now(),
    }
}

// Allow returns true if the request is permitted.
func (fw *FixedWindow) Allow() bool {
    fw.mu.Lock()
    defer fw.mu.Unlock()

    now := time.Now()
    if now.Sub(fw.windowStart) >= fw.window {
        // New window: reset counter
        fw.windowStart = now
        fw.count = 0
    }

    if fw.count >= fw.limit {
        return false
    }
    fw.count++
    return true
}
```

**The burst problem at the window boundary:**

```
Limit: 100 req/min

  11:59:50 → 100 requests arrive → allowed (window 1)
  12:00:05 → 100 requests arrive → allowed (window 2 just started)

In 15 seconds, 200 requests went through — 2× the rate limit!
```

This is the fundamental flaw of fixed windows: a client can burst 2× the limit by straddling the boundary.

---

## 3. Sliding Window Log

Keep a timestamped log of every request. Count only requests in the past window duration. Accurate, but memory-heavy for high-traffic clients.

```go
package ratelimit

import (
    "sync"
    "time"
)

// SlidingWindowLog is accurate but stores one timestamp per request.
// For 1000 req/min limit with 1M clients this is 1000×1M = 1B entries — impractical.
// Use for low-volume, high-accuracy scenarios (e.g. per-user webhook deliveries).
type SlidingWindowLog struct {
    mu     sync.Mutex
    limit  int
    window time.Duration
    log    []time.Time
}

func NewSlidingWindowLog(limit int, window time.Duration) *SlidingWindowLog {
    return &SlidingWindowLog{limit: limit, window: window}
}

func (s *SlidingWindowLog) Allow() bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-s.window)

    // Evict timestamps outside the window
    i := 0
    for i < len(s.log) && s.log[i].Before(cutoff) {
        i++
    }
    s.log = s.log[i:]

    if len(s.log) >= s.limit {
        return false
    }
    s.log = append(s.log, now)
    return true
}
```

---

## 4. Token Bucket

The **token bucket** is the most practical algorithm. Think of a bucket that refills at a constant rate (the rate limit) and holds up to a maximum number of tokens (the burst capacity). Each request consumes one token. If the bucket is empty, the request is rejected.

```
Properties:
  - Allows short bursts up to the bucket capacity
  - Long-term average throughput is capped at the refill rate
  - Unused capacity accumulates (up to burst limit), enabling legitimate bursts
```

```go
package ratelimit

import (
    "sync"
    "time"
)

// TokenBucket implements a token bucket rate limiter.
type TokenBucket struct {
    mu       sync.Mutex
    tokens   float64   // current token count (float to handle fractional refills)
    capacity float64   // max tokens (burst size)
    rate     float64   // tokens added per second
    lastRefill time.Time
}

// NewTokenBucket creates a bucket with the given burst capacity and refill rate.
// Example: NewTokenBucket(10, 2) → burst of 10, refill 2 tokens/sec (120 req/min steady state).
func NewTokenBucket(capacity int, ratePerSec float64) *TokenBucket {
    return &TokenBucket{
        tokens:     float64(capacity),
        capacity:   float64(capacity),
        rate:       ratePerSec,
        lastRefill: time.Now(),
    }
}

// Allow returns true if a token is available. Thread-safe.
func (tb *TokenBucket) Allow() bool {
    return tb.AllowN(1)
}

// AllowN returns true if n tokens are available (for variable-cost requests).
func (tb *TokenBucket) AllowN(n float64) bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()

    now := time.Now()
    elapsed := now.Sub(tb.lastRefill).Seconds()
    tb.lastRefill = now

    // Add tokens for elapsed time, capped at capacity
    tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.rate)

    if tb.tokens < n {
        return false
    }
    tb.tokens -= n
    return true
}

func min(a, b float64) float64 {
    if a < b { return a }
    return b
}
```

**Trace: bucket capacity=5, rate=1 token/sec, 8 requests 0.3s apart:**

```
t=0.0s: tokens=5. Request → allow, tokens=4
t=0.3s: tokens=4+0.3=4.3. Request → allow, tokens=3.3
t=0.6s: tokens=3.3+0.3=3.6. Request → allow, tokens=2.6
t=0.9s: tokens=2.6+0.3=2.9. Request → allow, tokens=1.9
t=1.2s: tokens=1.9+0.3=2.2. Request → allow, tokens=1.2
t=1.5s: tokens=1.2+0.3=1.5. Request → allow, tokens=0.5
t=1.8s: tokens=0.5+0.3=0.8. Request → allow, tokens=-0.2 → DENY
t=2.1s: tokens=-0.2+0.3=0.1 → DENY (still recovering)
```

---

## 5. Leaky Bucket

The **leaky bucket** processes requests at a constant output rate. Incoming requests queue up; if the queue is full, new requests are dropped. Unlike the token bucket, there are no bursts — output is always smooth.

```go
package ratelimit

import (
    "context"
    "time"
)

// LeakyBucket processes at most ratePerSec requests per second.
// Excess requests queue up to queueSize; beyond that they are dropped.
type LeakyBucket struct {
    queue    chan struct{}
    ticker   *time.Ticker
    done     chan struct{}
}

func NewLeakyBucket(ratePerSec int, queueSize int) *LeakyBucket {
    lb := &LeakyBucket{
        queue:  make(chan struct{}, queueSize),
        ticker: time.NewTicker(time.Second / time.Duration(ratePerSec)),
        done:   make(chan struct{}),
    }
    go lb.drain()
    return lb
}

func (lb *LeakyBucket) drain() {
    for {
        select {
        case <-lb.ticker.C:
            select {
            case <-lb.queue: // process one request
            default:         // queue is empty, nothing to do
            }
        case <-lb.done:
            return
        }
    }
}

// Submit enqueues a request. Returns false if queue is full (dropped).
func (lb *LeakyBucket) Submit() bool {
    select {
    case lb.queue <- struct{}{}:
        return true
    default:
        return false // queue full
    }
}

// Wait blocks until a slot opens in the queue or context is cancelled.
func (lb *LeakyBucket) Wait(ctx context.Context) error {
    select {
    case lb.queue <- struct{}{}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (lb *LeakyBucket) Stop() { lb.ticker.Stop(); close(lb.done) }
```

**Token bucket vs leaky bucket:**

```
Scenario: API sends 10 requests at once, limit = 2 req/sec

Token bucket (capacity=5):
  t=0s: 5 tokens → 5 requests go through immediately (burst allowed)
  t=1s: 2 more tokens → 2 more requests
  t=2s: 2 more tokens → last 3 requests (3 remain)

Leaky bucket (rate=2/sec, queue=5):
  t=0s: queue 5 (capacity), 5 dropped
  t=0.5s: process 1
  t=1.0s: process 1
  t=1.5s: process 1
  t=2.0s: process 1
  t=2.5s: process 1 → done after 2.5 seconds, perfectly smooth
```

Use the leaky bucket when downstream systems need predictable input rates (e.g. a database that can handle exactly N writes/sec).

---

## 6. Per-IP Middleware with Token Bucket

A real HTTP middleware that rate-limits each client IP independently.

```go
package middleware

import (
    "fmt"
    "net"
    "net/http"
    "sync"
    "time"
)

// IPRateLimiter holds a token bucket per client IP.
type IPRateLimiter struct {
    mu       sync.RWMutex
    limiters map[string]*tokenBucketEntry
    capacity int
    rate     float64 // tokens per second
}

type tokenBucketEntry struct {
    bucket   *TokenBucket
    lastSeen time.Time
}

// NewIPRateLimiter creates a limiter. capacity = burst size, rate = sustained req/sec.
func NewIPRateLimiter(capacity int, ratePerSec float64) *IPRateLimiter {
    rl := &IPRateLimiter{
        limiters: make(map[string]*tokenBucketEntry),
        capacity: capacity,
        rate:     ratePerSec,
    }
    go rl.cleanupLoop() // evict stale entries to prevent unbounded memory growth
    return rl
}

func (rl *IPRateLimiter) getBucket(ip string) *TokenBucket {
    // Fast path: read lock
    rl.mu.RLock()
    entry, ok := rl.limiters[ip]
    rl.mu.RUnlock()
    if ok {
        entry.lastSeen = time.Now()
        return entry.bucket
    }

    // Slow path: write lock to create new bucket
    rl.mu.Lock()
    defer rl.mu.Unlock()
    // Double-check after acquiring write lock
    if entry, ok = rl.limiters[ip]; ok {
        return entry.bucket
    }
    bucket := NewTokenBucket(rl.capacity, rl.rate)
    rl.limiters[ip] = &tokenBucketEntry{bucket: bucket, lastSeen: time.Now()}
    return bucket
}

// cleanupLoop evicts IPs not seen in the last 5 minutes.
func (rl *IPRateLimiter) cleanupLoop() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        cutoff := time.Now().Add(-5 * time.Minute)
        rl.mu.Lock()
        for ip, entry := range rl.limiters {
            if entry.lastSeen.Before(cutoff) {
                delete(rl.limiters, ip)
            }
        }
        rl.mu.Unlock()
    }
}

// Middleware returns an http.Handler middleware that enforces per-IP rate limiting.
func (rl *IPRateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := clientIP(r)
        bucket := rl.getBucket(ip)

        if !bucket.Allow() {
            // Retry-After: how long until one token refills
            retryAfter := int(1.0 / rl.rate)
            if retryAfter < 1 { retryAfter = 1 }
            w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.capacity))
            w.Header().Set("X-RateLimit-Remaining", "0")
            http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// clientIP extracts the real client IP, respecting X-Forwarded-For behind a proxy.
func clientIP(r *http.Request) string {
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        // X-Forwarded-For: client, proxy1, proxy2 — take the first (leftmost)
        // WARNING: only trust this header if your load balancer strips it first
        if ip, _, err := net.SplitHostPort(xff); err == nil {
            return ip
        }
        return xff
    }
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil { return r.RemoteAddr }
    return ip
}
```

**Wiring it into a chi router:**

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "yourapp/middleware"
)

func main() {
    rl := middleware.NewIPRateLimiter(
        10,  // burst: 10 requests at once
        2.0, // steady state: 2 requests/second = 120/minute
    )

    r := chi.NewRouter()
    r.Use(rl.Middleware) // apply to all routes

    r.Get("/api/products", listProducts)
    r.Post("/api/orders", createOrder)

    http.ListenAndServe(":8080", r)
}
```

---

## 7. Distributed Rate Limiting with Redis

A single in-memory limiter works for one server. In a horizontally scaled service with multiple instances, you need a shared counter. Redis sorted sets provide atomic, expiring, per-client counters.

**Pattern: sliding window log with Redis sorted sets**

```go
package ratelimit

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisLimiter implements a sliding window log using a Redis sorted set.
// The set key is per-client. Each member is a unique request ID (timestamp+nano).
// The score is the Unix timestamp in milliseconds — used for range queries.
type RedisLimiter struct {
    client *redis.Client
    limit  int
    window time.Duration
}

func NewRedisLimiter(client *redis.Client, limit int, window time.Duration) *RedisLimiter {
    return &RedisLimiter{client: client, limit: limit, window: window}
}

// Allow returns (allowed bool, remaining int, err error).
// key is the rate-limit identifier, e.g. "ratelimit:ip:192.168.1.1".
func (rl *RedisLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
    now := time.Now()
    windowStart := now.Add(-rl.window).UnixMilli()
    nowMs := now.UnixMilli()
    member := fmt.Sprintf("%d-%d", nowMs, now.UnixNano()) // unique per request

    pipe := rl.client.Pipeline()

    // Remove entries older than the window
    pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", windowStart))

    // Count entries in the current window
    countCmd := pipe.ZCard(ctx, key)

    // Add the current request
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(nowMs), Member: member})

    // Set TTL so the key auto-expires (avoid orphaned keys)
    pipe.Expire(ctx, key, rl.window+time.Second)

    if _, err := pipe.Exec(ctx); err != nil {
        return false, 0, fmt.Errorf("redis pipeline: %w", err)
    }

    count := int(countCmd.Val())
    if count >= rl.limit {
        // Remove the entry we just added — request is rejected
        rl.client.ZRem(ctx, key, member)
        return false, 0, nil
    }

    remaining := rl.limit - count - 1
    return true, remaining, nil
}
```

**Usage in middleware:**

```go
func RedisRateLimitMiddleware(limiter *RedisLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Use API key if authenticated, fall back to IP
            key := "ratelimit:ip:" + clientIP(r)
            if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
                key = "ratelimit:key:" + apiKey
            }

            allowed, remaining, err := limiter.Allow(r.Context(), key)
            if err != nil {
                // Fail open: if Redis is down, let the request through
                // (fail closed would be: reject and return 503)
                next.ServeHTTP(w, r)
                return
            }

            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
            if !allowed {
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**The Redis commands in detail:**

```
ZREMRANGEBYSCORE key -inf <window_start_ms>
  → evict entries older than the window

ZCARD key
  → count entries still in the window

ZADD key <now_ms> <unique_member>
  → record this request

EXPIRE key <window + 1s>
  → auto-cleanup when a client goes quiet
```

This is atomic per pipeline execution and scales across any number of app instances.

---

## 8. `golang.org/x/time/rate`

The standard way to rate-limit in Go. It implements a token bucket and is used inside the Go standard library itself.

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "golang.org/x/time/rate"
)

func main() {
    // rate.NewLimiter(r, b): r = events per second, b = burst size
    // rate.Every(d) converts a duration to a rate: Every(500ms) = 2/sec
    limiter := rate.NewLimiter(rate.Every(500*time.Millisecond), 5)

    // Allow: non-blocking check
    if limiter.Allow() {
        fmt.Println("request allowed")
    }

    // Wait: blocks until a token is available or context is cancelled
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    if err := limiter.Wait(ctx); err != nil {
        fmt.Println("wait failed:", err)
    }

    // Reserve: reserve a token for use in the future
    r := limiter.Reserve()
    if !r.OK() {
        fmt.Println("rate limit would never be satisfied")
    }
    time.Sleep(r.Delay()) // wait the appropriate time before proceeding
    fmt.Println("proceeding after reservation delay")
}
```

**Per-client limiter map using `x/time/rate`:**

```go
package main

import (
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

type PerClientLimiter struct {
    mu       sync.Mutex
    clients  map[string]*rate.Limiter
    r        rate.Limit
    b        int
}

func NewPerClientLimiter(r rate.Limit, b int) *PerClientLimiter {
    return &PerClientLimiter{
        clients: make(map[string]*rate.Limiter),
        r: r,
        b: b,
    }
}

func (pcl *PerClientLimiter) Get(key string) *rate.Limiter {
    pcl.mu.Lock()
    defer pcl.mu.Unlock()
    if l, ok := pcl.clients[key]; ok {
        return l
    }
    l := rate.NewLimiter(pcl.r, pcl.b)
    pcl.clients[key] = l
    return l
}

func (pcl *PerClientLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        limiter := pcl.Get(clientIP(r))
        if !limiter.Allow() {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    pcl := NewPerClientLimiter(
        rate.Every(time.Second), // 1 token per second
        5,                        // burst of 5
    )
    http.ListenAndServe(":8080", pcl.Middleware(http.DefaultServeMux))
}
```

**When to use `x/time/rate` vs rolling your own:**
- Use `x/time/rate` for single-process rate limiting — it is well-tested and production-proven.
- Roll your own token bucket when you need custom behaviour: variable-cost requests (`AllowN`), exposing `Retry-After`, or embedding into a larger data structure.
- Use Redis for multi-instance distributed limiting.

---

## Summary

- **Fixed window**: simple but allows 2× burst at window boundaries.
- **Sliding window log**: accurate but O(requests) memory per client.
- **Token bucket**: the right default. Allows bursts up to capacity, steady-state capped at refill rate.
- **Leaky bucket**: smooths output to constant rate. Use when downstream systems need predictable input.
- Per-IP in-memory limiting: use `sync.Map` or a guarded `map` with periodic cleanup to avoid unbounded growth.
- Always return `429 Too Many Requests` with `Retry-After` and `X-RateLimit-Remaining` headers.
- Distributed limiting: Redis sorted sets with `ZREMRANGEBYSCORE` + `ZCARD` + `ZADD` in a pipeline.
- Production shortcut: `golang.org/x/time/rate` for single-process token bucket.

---

## Exercises

### Easy

1. Modify `TokenBucket.AllowN` to return the number of tokens remaining after the request. Use this to set the `X-RateLimit-Remaining` response header.

2. Write a test for the fixed-window burst problem: send 100 requests at `t=window-1ms`, then 100 more at `t=window+1ms`. Assert that 200 requests got through despite a limit of 100 per window.

3. Implement a simple **rate limit exceeded** response that includes a JSON body: `{"error": "rate limit exceeded", "retry_after_seconds": 3, "limit": 100, "window": "1m"}`.

### Medium

4. Extend `IPRateLimiter` to support **tiered limits**: free-tier clients get 10 req/sec, premium clients get 100 req/sec. Accept a `tierFunc func(r *http.Request) string` that returns the tier name. Store separate `TokenBucket` configs per tier.

5. Implement a **backoff-aware client**: given a `*rate.Limiter`, write an HTTP client that automatically retries on `429` responses, reading the `Retry-After` header and sleeping accordingly. Use exponential backoff as a fallback when `Retry-After` is absent.

6. Implement **rate limiting by route**: different endpoints have different limits (e.g. `POST /api/checkout` = 5/min, `GET /api/products` = 100/min). Create a `RouteRateLimiter` that takes a map of route patterns to limits and applies them selectively.

### Hard

7. Implement **token bucket with Redis** for distributed use: instead of an in-memory `float64` counter, store `(tokens, last_refill_time)` in a Redis hash. Use a Lua script to make the refill-and-consume operation atomic. This is more complex than the sliding window approach but uses O(1) memory per client instead of O(requests).

8. Build a **rate-limiting load test**: using `sync.WaitGroup` and goroutines, simulate 1000 concurrent clients each sending requests as fast as possible against a test HTTP server protected by your `IPRateLimiter`. Measure: total requests served, total rejected, and verify the rejection rate matches expectations. Plot a histogram of request latency (use `time.Since` + a `[]time.Duration` slice) to see the effect of queueing.
