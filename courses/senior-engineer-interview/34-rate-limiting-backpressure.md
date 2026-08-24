# Chapter 34: Rate Limiting, Backpressure & Graceful Degradation

Senior engineers are expected to build systems that stay healthy under load spikes. Rate limiting protects your services. Backpressure prevents cascading overload. Graceful degradation keeps core features working when dependencies fail.

## Table of Contents

1. [Why Rate Limiting Matters](#1-why-rate-limiting-matters)
2. [Rate Limiting Algorithms](#2-rate-limiting-algorithms)
3. [Distributed Rate Limiting with Redis](#3-distributed-rate-limiting-with-redis)
4. [Go Implementation](#4-go-implementation--in-process-rate-limiter)
5. [Backpressure Patterns](#5-backpressure-patterns)
6. [Load Shedding](#6-load-shedding)
7. [Graceful Degradation](#7-graceful-degradation)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Why Rate Limiting Matters

```
Without rate limiting:
  One bad client sends 100K requests/second
  → Server CPU maxes out
  → All other clients experience degraded performance
  → Server crashes → everyone is affected

With rate limiting:
  Client exceeds their quota → 429 Too Many Requests
  Other clients unaffected
  Server remains healthy

Use cases:
  - API quotas: free tier = 100 req/hour, paid = 10,000 req/hour
  - DDoS protection: block IPs that exceed limits
  - Fair usage: prevent one tenant from starving others
  - Cost protection: prevent runaway jobs from infinite API calls
```

---

## 2. Rate Limiting Algorithms

### Token Bucket

Imagine a bucket that fills with tokens at a fixed rate. Each request consumes one token. If no tokens, reject the request.

```
Bucket capacity: 100 tokens
Refill rate: 10 tokens/second

Second 0:  100 tokens available
Second 1:  Client sends 50 requests → 50 tokens remain → all served
Second 2:  Refill 10 tokens → 60 tokens
Second 3:  Client sends 100 requests → 60 served, 40 rejected (429)

Properties:
  - Allows bursting up to capacity
  - Smooth average rate over time
  - Most common algorithm for API rate limiting
```

### Leaky Bucket

Requests enter a queue (bucket). A fixed number drain (leak) per second regardless of request rate. Excess requests overflow (rejected).

```
Queue capacity: 100 requests
Drain rate: 10 requests/second (smoothed output rate)

Properties:
  - No bursting — output is perfectly smooth
  - Good for shaping traffic into consistent rate
  - Used for network traffic shaping
```

### Fixed Window Counter

Count requests in a fixed time window. Reset counter at window boundary.

```
Window: 1 minute
Limit: 100 requests/minute

00:00 - 00:59: counter = 0
00:30: 100 requests arrive → all served (counter = 100)
00:59: 100 MORE requests → rejected (counter already 100)
01:00: counter resets to 0
01:00: 100 requests arrive → all served
→ 100 + 100 = 200 requests in 2 seconds around 00:59-01:01!

Problem: the "boundary burst" — 2x the rate limit around window resets.
```

### Sliding Window Log

Keep a log of timestamps of recent requests. Count how many fall in the last N seconds.

```
At 00:30:00, check last 60 seconds:
  Timestamps: [00:00:01, 00:00:05, ..., 00:30:00]
  Count timestamps in [00:29:00, 00:30:00]

Properties:
  - Perfectly accurate sliding window
  - Expensive: must store and count timestamps per user
  - Good for low-traffic, high-accuracy requirements
```

### Sliding Window Counter (best for distributed)

Approximate sliding window using two fixed window counters:

```
current window count + previous window count × (overlap fraction)

Example at 00:30 (30% into current window, 70% overlap with previous window):
  previous window count = 80
  current window count = 30
  estimated count = 30 + 80 × 0.7 = 86

Properties:
  - O(1) storage per user (just two counters)
  - Very close approximation of true sliding window
  - Commonly used in distributed systems (Redis-backed)
```

---

## 3. Distributed Rate Limiting with Redis

For multi-instance services, rate limiting state must be shared:

```go
// Token bucket using Redis
func allowRequest(ctx context.Context, rdb *redis.Client, userID string, limit int, window time.Duration) (bool, error) {
    // Lua script runs atomically on Redis
    script := redis.NewScript(`
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])  -- in seconds
        local now = tonumber(ARGV[3])
        
        -- Remove timestamps outside the sliding window
        redis.call("ZREMRANGEBYSCORE", key, 0, now - window)
        
        -- Count requests in current window
        local count = redis.call("ZCARD", key)
        
        if count < limit then
            -- Add this request's timestamp
            redis.call("ZADD", key, now, now .. "-" .. math.random())
            redis.call("EXPIRE", key, window + 1)
            return 1  -- allowed
        end
        return 0  -- rate limited
    `)
    
    key := fmt.Sprintf("rate_limit:%s", userID)
    now := float64(time.Now().UnixNano()) / 1e9
    
    result, err := script.Run(ctx, rdb,
        []string{key},
        limit,
        int(window.Seconds()),
        now).Int()
    
    return result == 1, err
}

// HTTP middleware using token bucket:
func rateLimitMiddleware(rdb *redis.Client, limit int) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID := r.Header.Get("X-User-ID")
            if userID == "" {
                userID = r.RemoteAddr // fallback to IP
            }
            
            allowed, err := allowRequest(r.Context(), rdb, userID, limit, time.Minute)
            if err != nil {
                // Fail open: if Redis is down, allow the request
                next.ServeHTTP(w, r)
                return
            }
            
            if !allowed {
                w.Header().Set("Retry-After", "60")
                w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 4. Go Implementation — In-Process Rate Limiter

```go
import "golang.org/x/time/rate"

// Token bucket rate limiter — 10 requests/second, burst of 30
limiter := rate.NewLimiter(rate.Limit(10), 30)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Wait blocks until a token is available or context is done:
    if err := limiter.Wait(r.Context()); err != nil {
        http.Error(w, "rate limited", http.StatusTooManyRequests)
        return
    }
    // Process request...
}

// Non-blocking check:
if !limiter.Allow() {
    http.Error(w, "rate limited", http.StatusTooManyRequests)
    return
}

// Per-user limiters using sync.Map:
type UserLimiter struct {
    mu       sync.Mutex
    limiters map[string]*rate.Limiter
}

func (ul *UserLimiter) Get(userID string) *rate.Limiter {
    ul.mu.Lock()
    defer ul.mu.Unlock()
    
    if l, exists := ul.limiters[userID]; exists {
        return l
    }
    l := rate.NewLimiter(rate.Limit(10), 30) // 10 req/s per user
    ul.limiters[userID] = l
    return l
}
```

---

## 5. Backpressure Patterns

Backpressure is signaling to upstream producers to slow down when downstream consumers are overwhelmed.

```go
// Pattern 1: Bounded channel as backpressure mechanism
jobs := make(chan Job, 100) // buffer of 100

// Producer: sends work into the channel
go func() {
    for _, job := range pendingJobs {
        select {
        case jobs <- job:
            // sent successfully
        default:
            // channel full → backpressure applied → drop or reject
            log.Printf("job queue full, dropping job %d", job.ID)
        }
    }
}()

// Consumers: drain the channel at their own pace
for i := 0; i < workerCount; i++ {
    go func() {
        for job := range jobs {
            processJob(job)
        }
    }()
}

// Pattern 2: Return 503 Service Unavailable when queue is full
func handleRequest(w http.ResponseWriter, r *http.Request) {
    select {
    case workerPool.queue <- r:
        // queued successfully
    default:
        // queue full — reject and let client retry
        w.Header().Set("Retry-After", "5")
        http.Error(w, "service overloaded", http.StatusServiceUnavailable)
    }
}
```

---

## 6. Load Shedding

Proactively drop low-priority requests when the system is under extreme load to protect critical paths:

```go
// Admission control based on system load:
type LoadShedder struct {
    maxConcurrent int64
    current       atomic.Int64
}

func (ls *LoadShedder) Allow() bool {
    current := ls.current.Load()
    return current < ls.maxConcurrent
}

func (ls *LoadShedder) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // High-priority paths (health checks, payment) bypass shedding:
        if r.URL.Path == "/health" || r.URL.Path == "/api/v1/payments" {
            next.ServeHTTP(w, r)
            return
        }
        
        ls.current.Add(1)
        defer ls.current.Add(-1)
        
        if !ls.Allow() {
            http.Error(w, "service overloaded", http.StatusServiceUnavailable)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## 7. Graceful Degradation

When a dependency fails, return a degraded-but-functional response rather than failing completely:

```go
// Example: product recommendations
func getRecommendations(ctx context.Context, userID string) []Product {
    // Try personalized recommendations first
    recs, err := recommendationSvc.GetPersonalized(ctx, userID)
    if err == nil {
        return recs
    }
    
    // Recommendation service is down: fall back to trending products
    log.Printf("recommendation service unavailable, falling back to trending: %v", err)
    trending, err := cache.GetTrending(ctx)
    if err == nil {
        return trending
    }
    
    // Cache is also down: return empty (page still works, just no recommendations)
    log.Printf("trending cache unavailable: %v", err)
    return []Product{}
}

// Example: user profile with feature flags
func getUserProfile(ctx context.Context, userID string) *UserProfile {
    profile, err := db.GetUser(ctx, userID)
    if err != nil {
        return &UserProfile{Name: "Guest"} // degraded response, not error
    }
    
    // Try to enrich with premium status (non-critical)
    isPremium, err := premiumSvc.Check(ctx, userID)
    if err != nil {
        isPremium = false // default to non-premium if service is down
    }
    
    profile.IsPremium = isPremium
    return profile
}
```

---

## 8. Interview Questions & Model Answers

**Q: How would you design a rate limiter for a public API?**

"I'd use a token bucket algorithm with a Redis backend for distributed rate limiting. Each user/API key gets a bucket tracked with a sorted set in Redis — add a timestamp on each request, remove old timestamps outside the window, count remaining. I'd expose rate limit headers in responses (X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After on 429). For tiers: free tier gets 100 req/hour, paid gets 10,000 req/hour, tracked by API key. I'd implement the Redis Lua script atomically to avoid race conditions. If Redis is unavailable, I'd fail open (allow requests) rather than blocking everyone — rate limiting should be a best-effort protection, not a hard dependency."

**Q: What is backpressure and why is it important?**

"Backpressure is when a downstream system signals to an upstream producer to slow down. Without it, a slow consumer and fast producer leads to unbounded buffer growth, then OOM crashes, then cascading failures. In Go, this manifests as bounded channels — when a channel buffer is full, the sender blocks (or gets rejected if using `select default`). At the HTTP level, returning 429 or 503 is backpressure — telling the client to slow down. In Kafka, consumer lag grows when consumers can't keep up — this is backpressure being applied (messages queue, consumers catch up at their own pace). The key insight: it's better to reject requests at the edge than to have them pile up internally and cause the whole system to collapse."

---

## Summary

- **Token bucket:** allow bursting up to capacity, refill at fixed rate. Most common API rate limiting algorithm.
- **Sliding window:** accurate but expensive. Approximate with two fixed counters for efficiency.
- **Redis Lua scripts:** atomic rate limit operations across distributed instances.
- **Fail open when Redis is down:** rate limiting is a protection, not a hard dependency.
- **Backpressure:** bounded channels, 503 responses. Reject at the edge before internal queues fill.
- **Load shedding:** drop low-priority requests proactively when under extreme load. Protect critical paths.
- **Graceful degradation:** return degraded-but-functional responses when dependencies fail. Empty recommendations > error page.
