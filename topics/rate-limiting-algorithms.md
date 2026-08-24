---
title: Rate Limiting Algorithms
category: Software & Programming
tags: [System Design, APIs]
duration: 8 min read
relatedCourses: [senior-engineer-interview, go-programming]
relatedProjects: [rest-api-server]
relatedTopics: [idempotency-in-apis, load-balancing-strategies]
---

## TL;DR

- **Fixed window**: simplest, but allows a burst of 2x the limit right at the window boundary.
- **Sliding window log**: exact, but stores every request timestamp — memory cost scales with request rate.
- **Sliding window counter**: a practical approximation of the log approach, using two fixed windows and interpolation, at constant memory.
- **Token bucket**: allows controlled bursts up to bucket capacity, then throttles to a steady refill rate — the most common choice for public APIs.
- **Leaky bucket**: smooths bursts into a perfectly constant output rate, at the cost of added latency for bursty traffic.

## Why Rate Limiting Exists

Without it, one client — buggy, malicious, or just enthusiastic — can consume enough of a shared resource (CPU, database connections, downstream API quota) to degrade the service for everyone else. Rate limiting caps how much of that shared resource any one identity (a user, an API key, an IP) can consume in a given time window.

## Fixed Window Counter

Divide time into fixed windows (e.g., one per minute). Count requests in the current window; reset the count to zero at each window boundary.

```
Window: [0:00-0:01) limit=100

count = 0
on request:
  if now is a new window: count = 0
  count++
  if count > 100: reject
  else: allow
```

**The boundary burst problem**: if a client sends 100 requests at 0:00:59 and another 100 at 0:01:00, both batches are individually within the limit for their respective windows — but the client just sent 200 requests in a two-second span. Fixed windows don't see across the boundary they just crossed.

Despite this flaw, fixed window is extremely common in practice because it's trivial to implement (often just one counter with a TTL in Redis) and the boundary burst is a real but bounded, well-understood edge case most systems can tolerate.

## Sliding Window Log

Store the timestamp of every request. To check if a new request is allowed, drop timestamps older than the window, then count what's left.

```
requests = [] // timestamps of allowed requests

on request at time t:
  requests = requests.filter(ts => ts > t - windowSize)
  if len(requests) >= limit: reject
  else:
    requests.append(t)
    allow
```

This is exact — no boundary burst problem at all, because the "window" is always the trailing N seconds relative to *now*, not a fixed clock-aligned bucket. The cost is memory: you're storing a timestamp per request, per rate-limited identity, which doesn't scale well if you have many high-volume clients.

## Sliding Window Counter (the practical middle ground)

Keep two fixed-window counters (current and previous), and estimate the sliding window's count as a weighted combination of both:

```
estimated_count = current_window_count +
                  previous_window_count * (1 - elapsed_fraction_of_current_window)
```

If you're 30% of the way through the current window, you count all of the current window's requests plus 70% of the previous window's — approximating "how many requests happened in the trailing 60 seconds" without storing individual timestamps. This is what most production rate limiters (including the ones behind Cloudflare, Redis-based limiters like `redis-cell`) actually implement, because it gets sliding-window accuracy at fixed-window memory cost.

## Token Bucket

A bucket holds up to *N* tokens. Tokens refill at a fixed rate (e.g., 10/second) up to the bucket's capacity. Each request consumes one token; if the bucket is empty, the request is rejected (or queued, depending on design).

```
capacity = 100
refillRate = 10 // tokens per second
tokens = 100
lastRefill = now()

on request:
  elapsed = now() - lastRefill
  tokens = min(capacity, tokens + elapsed * refillRate)
  lastRefill = now()
  if tokens < 1: reject
  else:
    tokens -= 1
    allow
```

The key property: a client that's been idle can **burst** up to the full bucket capacity all at once, then has to slow down to the refill rate. This matches real traffic patterns well — a user opening an app and firing off a batch of requests, then going quiet — which is why token bucket is the default choice for most public-facing APIs (AWS, Stripe, and GitHub's APIs all use variants of this).

## Leaky Bucket

Conceptually the inverse of token bucket: requests go into a queue (the "bucket"); they're processed ("leak out") at a strictly constant rate, regardless of how bursty the input was. If the queue fills up, new requests are rejected.

```
Incoming (bursty):  ||||    |  ||||||        |
Outgoing (steady):  | | | | | | | | | | | | | |
```

The difference from token bucket matters in practice: token bucket allows the *output* rate to burst (up to capacity) when there's been idle time to accumulate tokens; leaky bucket enforces a constant output rate no matter what, which smooths traffic but adds queuing latency for anything arriving faster than the leak rate. Leaky bucket suits systems where a perfectly steady downstream rate matters more than low latency for bursty clients — e.g., shaping traffic before it hits a fixed-capacity downstream service.

## Choosing One

| Need | Choice |
|---|---|
| Simplicity, and boundary bursts are tolerable | Fixed window |
| Exact limiting, low request volume per key | Sliding window log |
| Exact-ish limiting at scale | Sliding window counter |
| Allow legitimate bursts, then throttle | Token bucket |
| Perfectly smooth output rate | Leaky bucket |

## Common Pitfalls

- **Rate limiting by IP alone** — many legitimate users share an IP (corporate NAT, mobile carrier NAT), and it's trivial for an attacker to rotate IPs. Rate limit by authenticated identity (API key, user ID) wherever you can, and treat IP-based limiting as a coarser secondary layer.
- **Not returning `Retry-After` / rate-limit headers** — a rejected request without `X-RateLimit-Remaining`/`Retry-After` headers forces clients to guess when to retry, which usually makes them retry immediately and makes the problem worse.
- **Applying one global limit across a distributed fleet without a shared store** — if each server instance tracks its own in-memory counter, the *effective* limit becomes (per-instance limit × number of instances), not the intended limit. A shared store (Redis is the common choice) is required for an accurate limit across multiple servers.
- **Forgetting that the refill/leak rate itself needs tuning per endpoint** — a single global rate limit config is rarely right; a cheap read endpoint and an expensive write endpoint usually need very different limits.
