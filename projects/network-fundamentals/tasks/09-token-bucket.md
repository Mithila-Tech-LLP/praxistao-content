---
title: Token Bucket Rate Limiter
number: 9
difficulty: medium
duration: 20-30 minutes
concept: rate limiting, deterministic time injection
---

## What to Build

Implement `TokenBucket`, a classic rate limiter: a bucket holds tokens up to some capacity, refills at a steady rate over time, and each allowed action consumes one token.

## Function Signature

```go
type TokenBucket struct {
    capacity   float64
    tokens     float64
    refillRate float64 // tokens per second
    lastRefill time.Time
}

func NewTokenBucket(capacity, refillRate float64, now time.Time) *TokenBucket
func (b *TokenBucket) Allow(now time.Time) bool // refills based on elapsed time since lastRefill, then consumes 1 token if available
```

## Requirements

- `NewTokenBucket` starts the bucket full, i.e. `tokens = capacity`
- Both functions take `now time.Time` as an explicit parameter — never call `time.Now()` internally, so tests are fully deterministic with no real sleeping
- `Allow(now)` computes `elapsed := now.Sub(b.lastRefill)`, adds `elapsed.Seconds() * refillRate` tokens (capped at `capacity`), and updates `lastRefill = now`
- After refilling, if `tokens >= 1`, subtract 1 and return `true`; otherwise return `false`

## Key Concept: Rate Limiting with Injected Time

Token bucket is the standard algorithm behind API rate limits, traffic shaping, and (conceptually) TCP's own sending discipline covered in Chapter 62's congestion control. The key engineering discipline in this task isn't the algorithm itself — it's making it testable: by passing `now` in explicitly instead of reaching for `time.Now()`, the whole bucket becomes a pure function of its inputs, so a test can simulate two seconds passing instantly instead of actually sleeping for two seconds.

## Hints

<details>
<summary>Hint 1: Refill happens on every call, not on a timer</summary>

There's no background goroutine ticking away. Instead, every call to `Allow` first "catches up" the bucket based on however much time has passed since the last call, then evaluates whether there's a token to spend.

</details>

<details>
<summary>Hint 2: Capping at capacity</summary>

Tokens shouldn't accumulate forever if nothing consumes them — after adding the refill amount, clamp `tokens` back down to `capacity` if the refill pushed it over.

</details>

<details>
<summary>Hint 3: Test the boolean results, not the float internals</summary>

Floating point tokens can accumulate small rounding errors, so don't assert on the exact value of `b.tokens` in a test. Assert on what `Allow` returns (`true`/`false`) at each step instead — that's the actually observable, deterministic behavior.

</details>

## How to Verify

```bash
lncli run
```
