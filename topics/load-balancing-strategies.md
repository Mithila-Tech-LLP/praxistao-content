---
title: Load Balancing Strategies
category: Software & Programming
tags: [System Design, Networking]
duration: 7 min read
relatedCourses: [senior-engineer-interview, go-programming]
relatedProjects: [grpc-service]
relatedTopics: [consistent-hashing-explained, circuit-breakers-and-retries]
---

## TL;DR

- A load balancer's job is deciding, for each incoming request, which of several backend instances should handle it — the algorithm it uses to decide is the actual design choice.
- **Round robin** and **weighted round robin** are simple and stateless but ignore each backend's actual current load.
- **Least connections** adapts to real load but needs the balancer to track connection counts per backend.
- **Consistent hashing** trades perfect load evenness for *session affinity* — the same client (or key) reliably lands on the same backend.
- Layer 4 (TCP-level) vs Layer 7 (HTTP-aware) load balancing is an orthogonal, equally important decision — it determines what information the balancer can actually use to route.

## Round Robin

Send request 1 to backend A, request 2 to B, request 3 to C, request 4 back to A, and so on — a fixed rotation.

```go
type RoundRobin struct {
    backends []string
    next     int
    mu       sync.Mutex
}

func (r *RoundRobin) Pick() string {
    r.mu.Lock()
    defer r.mu.Unlock()
    b := r.backends[r.next]
    r.next = (r.next + 1) % len(r.backends)
    return b
}
```

Simple, stateless (beyond the rotation index), and works well when every backend is equally capable and every request costs roughly the same to handle. It falls apart when either of those isn't true: a backend that's slower or already overloaded still gets exactly its "fair share" of new requests under plain round robin, piling on further.

**Weighted round robin** fixes the "not all backends are equally capable" half of that problem by giving faster/bigger backends a higher share:

```
weights: A=5, B=3, C=2  (out of 10)
-> roughly 50% of requests to A, 30% to B, 20% to C
```

It still doesn't account for *current* load, only *declared* capacity.

## Least Connections

Route each new request to whichever backend currently has the fewest active connections. This directly adapts to real-time load rather than a fixed rotation or declared weight — a backend that's slow to finish requests (for any reason: a slow query, a GC pause, whatever) naturally accumulates connections and gets fewer new ones routed to it until it catches up.

```go
func (lb *LeastConn) Pick() string {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    best := lb.backends[0]
    for _, b := range lb.backends[1:] {
        if lb.activeConns[b] < lb.activeConns[best] {
            best = b
        }
    }
    return best
}
```

The tradeoff: the balancer now needs to track connection state per backend and keep it in sync, which is more bookkeeping than a stateless round robin counter — and "active connections" is a proxy for load, not load itself (a backend could have few connections but each one doing something expensive).

## Consistent Hashing (for Session Affinity)

Sometimes even distribution isn't actually the goal — you need the *same* client (or the same cache key, the same user session) to reliably land on the *same* backend across requests, so that backend can keep useful local state (an in-memory session, a warm cache entry) instead of every request needing to fetch that state fresh from a shared store.

```
backend = ring.getServer(hash(clientIP)) // or hash(sessionID), hash(userID)...
```

This is the same ring structure used for distributed caching and sharding (see Consistent Hashing in Related) — applied here so that adding or removing a backend only reshuffles a proportional slice of client-to-backend mappings, not everyone's session at once. The tradeoff versus least-connections/round-robin: perfectly even load distribution is no longer the primary goal, so a backend can end up with more than its even share if it happens to own a disproportionate slice of active clients.

## Layer 4 vs Layer 7 Load Balancing

This is a separate axis from the *algorithm* above — it's about *what the balancer can see and act on*:

- **Layer 4 (transport layer)**: the balancer only sees TCP/UDP connection info (source/destination IP and port) — it doesn't parse HTTP at all. It's fast and protocol-agnostic, but it can only route based on connection-level information, not anything about the actual request (URL path, headers, cookies).
- **Layer 7 (application layer)**: the balancer terminates and parses HTTP itself, so it can route based on the URL path (`/api/*` to one backend pool, `/static/*` to another), headers, cookies, or even request body content. This costs more CPU per request (actually parsing HTTP) but enables routing decisions Layer 4 fundamentally can't make.

Most modern reverse proxies (NGINX, Envoy, cloud load balancers) support both, and a real production setup often layers them — a Layer 4 balancer distributing across a fleet of Layer 7 proxies, which then do path-based routing to the actual services.

## Health Checks Are Not Optional

Every strategy above assumes the load balancer knows which backends are actually healthy. Without active health checks (periodic `GET /health` requests) or passive ones (tracking recent error rates per backend), a load balancer will happily keep sending traffic to a backend that's crashed, out of memory, or stuck — the algorithm choice doesn't matter at all if half your "backends" are actually dead and still receiving requests.

## Common Pitfalls

- **Choosing consistent hashing for pure load distribution when there's no actual affinity requirement** — it optimizes for a different goal (client stickiness) at some cost to even distribution; don't reach for it unless you specifically need the stickiness.
- **Least connections without accounting for connection *cost*** — a websocket connection open for an hour and a connection that completes in 5ms count identically as "one active connection," which can badly mislead a naive least-connections balancer.
- **Sticky sessions (consistent hashing / session affinity) without a fallback for backend failure** — if the backend a client is "stuck" to goes down, the balancer needs a defined behavior (route elsewhere, and accept the client loses its session-local state) rather than an undefined one.
- **No health checks, or health checks that don't actually exercise the failure mode that matters** — a `/health` endpoint that just returns 200 unconditionally tells you the process is running, not that it can actually serve real requests (e.g., that its database connection is up).
