---
title: Circuit Breakers and Retries
category: Software & Programming
tags: [System Design, Reliability]
duration: 8 min read
relatedCourses: [senior-engineer-interview, go-programming]
relatedProjects: [grpc-service, rest-api-server]
relatedTopics: [idempotency-in-apis, load-balancing-strategies]
---

## TL;DR

- A **retry** handles a transient failure by trying again — but naive retries under sustained failure make things *worse*, not better, by piling more load onto an already-struggling service.
- **Exponential backoff with jitter** is what makes retries safe at scale: increasing delay between attempts, plus randomization so many clients don't retry in synchronized waves.
- A **circuit breaker** stops calling a failing dependency altogether for a while, instead of retrying forever — it trades "keep trying" for "fail fast and give the dependency room to recover."
- These two patterns solve different problems and are meant to be combined, not chosen between.

## The Retry Storm Problem

A service calls a downstream dependency that starts failing — maybe it's overloaded, maybe it's mid-deploy. Every caller, doing the sensible-looking thing, retries the failed call immediately. That retry traffic lands on the downstream service on top of whatever load already caused it to start failing — making the overload worse, which causes more failures, which causes more retries. This feedback loop, a **retry storm**, is one of the most common ways a small, transient blip turns into a full outage.

## Exponential Backoff

Instead of retrying immediately, wait progressively longer between attempts:

```go
func retryWithBackoff(fn func() error, maxAttempts int) error {
    var err error
    for attempt := 0; attempt < maxAttempts; attempt++ {
        if err = fn(); err == nil {
            return nil
        }
        backoff := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
        time.Sleep(backoff)
    }
    return err
}
// attempt 0: fails, wait 100ms
// attempt 1: fails, wait 200ms
// attempt 2: fails, wait 400ms
// attempt 3: fails, wait 800ms
```

This spreads out retry pressure over time instead of hammering the dependency at full rate. It should always be paired with a **maximum backoff cap** and a **maximum attempt count** — unbounded exponential growth just means an unbounded wait, which is its own kind of failure (a request that "eventually" succeeds after 10 minutes is often no better than one that fails outright, from the caller's perspective).

## Jitter: Why Backoff Alone Isn't Enough

If 1,000 clients all start retrying a failed call at the same moment, plain exponential backoff has them all waiting *exactly* 100ms, then *exactly* 200ms, then *exactly* 400ms — still perfectly synchronized, just delayed. When they all retry together at each step, that's still a thundering-herd spike, just a periodic one instead of continuous.

**Jitter** adds randomization to break this synchronization:

```go
func backoffWithJitter(attempt int) time.Duration {
    base := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
    jitter := time.Duration(rand.Int63n(int64(base)))
    return jitter // "full jitter": random value between 0 and base, not base+jitter
}
```

"Full jitter" (picking a random delay between 0 and the computed backoff, rather than the backoff plus a small random offset) is the specific variant AWS's architecture blog popularized after finding it outperformed simpler jitter schemes at actually breaking up synchronized retry waves — worth knowing by name if this comes up in a design discussion.

## What Should Actually Be Retried

Not every failure should trigger a retry — this is where idempotency matters directly:

- **Safe to retry automatically**: timeouts, connection resets, 503/429 responses — these usually indicate a transient condition, and retrying a `GET` or an idempotent operation (see Idempotency in APIs) doesn't risk duplicate side effects.
- **Not safe to retry blindly**: a `POST` that isn't idempotent, or any error that indicates the request itself was invalid (400, 404, 401) — retrying an invalid request just fails the same way again, and retrying a non-idempotent write risks doing it twice.

## Circuit Breakers: Failing Fast Instead of Retrying Forever

Retries assume the failure is transient and will resolve soon. A circuit breaker handles the case where it isn't — a dependency that's genuinely down for an extended period. Instead of every caller retrying (with backoff or not) against a service that simply isn't coming back soon, the circuit breaker "trips" and makes calls fail immediately, without even attempting the network call, for a cooldown period.

Three states, matching the electrical metaphor:

```
CLOSED (normal) --[failure rate exceeds threshold]--> OPEN (failing fast)
   ^                                                        |
   |                                                [cooldown timer elapses]
   |                                                        v
   +----[a trial request in HALF-OPEN succeeds]---- HALF-OPEN (one trial request)
```

- **Closed**: requests flow through normally; the breaker tracks recent failure rate.
- **Open**: once failures cross a threshold (e.g., 50% of the last 20 requests failed), the breaker stops even attempting calls — it returns an error immediately for a fixed cooldown window. This is the key benefit: callers get a fast, cheap failure instead of waiting out a timeout on every single request, and the struggling dependency gets a period with reduced load to actually recover.
- **Half-open**: after the cooldown, the breaker allows a small number of trial requests through. If they succeed, it closes again (back to normal). If they fail, it reopens and waits another cooldown period.

```go
type CircuitBreaker struct {
    state         string // "closed", "open", "half-open"
    failures      int
    threshold     int
    openedAt      time.Time
    cooldown      time.Duration
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.openedAt) < cb.cooldown {
            return errors.New("circuit open — failing fast")
        }
        cb.state = "half-open"
    }

    err := fn()
    if err != nil {
        cb.failures++
        if cb.failures >= cb.threshold {
            cb.state = "open"
            cb.openedAt = time.Now()
        }
        return err
    }

    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

## Why Both, Together

Retries and circuit breakers solve different halves of the same problem:

- **Retries** handle *brief* transient blips — a single dropped packet, a momentary GC pause on the other end — where trying again shortly after is likely to succeed.
- **Circuit breakers** handle *sustained* failure — a dependency that's actually down — where continuing to retry (even with backoff) just keeps generating load against something that needs time to recover, and keeps every caller waiting on timeouts in the meantime.

A well-built client wraps a call with both: retry a few times with backoff+jitter for transient blips, but if the *overall* failure rate stays high across many calls, trip the circuit breaker so subsequent calls fail immediately instead of each one going through its own retry cycle against a dependency that's clearly not coming back in the next few hundred milliseconds.

## Common Pitfalls

- **Retrying without a cap, "just in case it eventually works"** — always bound both the number of attempts and the total time spent retrying.
- **No jitter** — backoff alone still produces synchronized retry waves across many clients; jitter is what actually de-synchronizes them.
- **A circuit breaker with no half-open trial state** — a breaker that just flips back to closed after the cooldown timer, without a cautious trial, dumps full traffic back onto a dependency that might have only *just* recovered, potentially re-tripping it immediately.
- **Retrying non-idempotent operations without an idempotency key** — this is exactly the scenario that turns "helpful reliability pattern" into "double-charged customer."
