---
title: Rate Limiter
task: 10
slug: rate-limiter
concept: Token Bucket, sync.Mutex, rate limiting
difficulty: advanced
---

## What You Will Build

Write a middleware factory that enforces a request rate limit using the token bucket algorithm. Rate limiting is how APIs prevent abuse — if a client sends too many requests too quickly, they get a `429 Too Many Requests` response.

## Function Signature

```go
func NewRateLimiter(requestsPerSecond int) func(http.Handler) http.Handler
```

Returns a middleware constructor. The returned function takes any `http.Handler` and wraps it with rate limiting.

## How It Should Work

**Token Bucket Algorithm:**

1. Create a bucket with capacity equal to `requestsPerSecond` tokens
2. Start a background goroutine that adds one token every `(1 second / requestsPerSecond)` interval (never exceeding capacity)
3. When a request arrives, attempt to take one token:
   - If tokens > 0: take one, call `next.ServeHTTP` (allow the request)
   - If tokens == 0: return `429 Too Many Requests` immediately (do not call `next`)
4. Use `sync.Mutex` to protect the token counter from concurrent access

## Example Usage

```go
limiter := NewRateLimiter(3)  // allow 3 req/s
handler := limiter(myHandler)
// first 3 rapid requests → 200
// 4th rapid request → 429
```

## Key Concepts

**Token Bucket** — a classic algorithm for rate limiting. It allows short bursts (up to bucket capacity) while enforcing a long-term average rate. Simple alternative to the "leaky bucket" or sliding window algorithms.

**sync.Mutex** — protects shared state from concurrent goroutines:

```go
var mu sync.Mutex
mu.Lock()
// safely read/write tokens
mu.Unlock()
```

**time.NewTicker** — fires on a channel at regular intervals:

```go
ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
for range ticker.C {
    // add a token
}
```

**Closure over mutable state** — each call to `NewRateLimiter` creates its own independent bucket with its own goroutine. The closure captures the `tokens` variable and `mu`.

## Hints

<details>
<summary>Hint 1 — Struct vs closure for state</summary>

You can store state in a closure (captures local variables) or in a struct. A struct is often cleaner:

```go
type rateLimiter struct {
    tokens int
    max    int
    mu     sync.Mutex
}

func (rl *rateLimiter) start(rps int) {
    interval := time.Second / time.Duration(rps)
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            rl.mu.Lock()
            if rl.tokens < rl.max {
                rl.tokens++
            }
            rl.mu.Unlock()
        }
    }()
}
```
</details>

<details>
<summary>Hint 2 — The middleware wrapper</summary>

```go
func NewRateLimiter(rps int) func(http.Handler) http.Handler {
    rl := &rateLimiter{tokens: rps, max: rps}
    rl.start(rps)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            rl.mu.Lock()
            if rl.tokens <= 0 {
                rl.mu.Unlock()
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            rl.tokens--
            rl.mu.Unlock()
            next.ServeHTTP(w, r)
        })
    }
}
```
</details>

<details>
<summary>Hint 3 — Why the test works without sleeps</summary>

The test sends 5 requests immediately in a tight loop. At 3 req/s, the refill interval is ~333ms. Because all 5 requests happen in microseconds, the refill goroutine won't fire in between. The bucket starts with 3 tokens, requests 1-3 consume them, and requests 4-5 see 0 tokens → 429. No sleep needed.
</details>

## How to Verify

```bash
cd starter/task-10-rate-limiter
go test ./...
```

The test creates a `3 req/s` limiter and fires 5 requests in immediate succession:

- Requests 1, 2, 3 → status `200`
- Requests 4, 5 → status `429`
