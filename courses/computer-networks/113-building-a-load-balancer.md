# Chapter 113: Building a Load Balancer

> **"A reverse proxy always knows where it's sending a request. A load balancer has to decide — and then has to notice, on its own, when that decision was wrong."**

---

## Table of Contents

1. [Recap: Chapter 112's Proxy, and Chapter 95's L4/L7 Distinction](#1-recap-chapter-112s-proxy-and-chapter-95s-l4l7-distinction)
2. [The Problem: One Backend Doesn't Scale, and Neither Does One That's Silently Dead](#2-the-problem-one-backend-doesnt-scale-and-neither-does-one-thats-silently-dead)
3. [The Naive Shortcut We're Deliberately Not Taking](#3-the-naive-shortcut-were-deliberately-not-taking)
4. [The Real Solution: A Backend Pool, a Selection Strategy, and Active Health Checks](#4-the-real-solution-a-backend-pool-a-selection-strategy-and-active-health-checks)
5. [Code: The Backend Pool and Round-Robin Selection](#5-code-the-backend-pool-and-round-robin-selection)
6. [Code: Least-Connections Selection](#6-code-least-connections-selection)
7. [Code: Active Health Checking in a Background Goroutine](#7-code-active-health-checking-in-a-background-goroutine)
8. [Code: The Complete Load Balancer](#8-code-the-complete-load-balancer)
9. [Hands-On Experiment: Rotation, Failure, and Recovery](#9-hands-on-experiment-rotation-failure-and-recovery)
10. [Worked Example: Round-Robin vs. Least-Connections Under Uneven Load](#10-worked-example-round-robin-vs-least-connections-under-uneven-load)
11. [Common Pitfalls in Hand-Rolled Load Balancing](#11-common-pitfalls-in-hand-rolled-load-balancing)
12. [Production Notes: What Real Load Balancers Add](#12-production-notes-what-real-load-balancers-add)
13. [What's Simplified Here](#13-whats-simplified-here)
14. [Interview Questions & Model Answers](#interview-questions--model-answers)
15. [Exercises](#exercises)
16. [Summary](#summary)

---

## 1. Recap: Chapter 112's Proxy, and Chapter 95's L4/L7 Distinction

Chapter 112 built a reverse proxy that always forwarded to the same one backend, decided by a hardcoded constant. Chapter 95 described, in prose, what a real load balancer adds on top of that idea: a *pool* of interchangeable backends, a *strategy* for picking one per request (round-robin, least-connections, and others), and *health checks* that keep the pool honest by removing backends that stop responding and re-adding them once they recover. This chapter builds all three, directly extending Chapter 112's code — this is the Layer 7 half of Chapter 95's L4/L7 split, built on `httputil.ReverseProxy` per backend rather than on Chapter 112's fully manual byte relay, because Layer 7 balancing decisions (by definition) require understanding HTTP well enough to hand a request to a specific backend's `ReverseProxy` instance.

---

## 2. The Problem: One Backend Doesn't Scale, and Neither Does One That's Silently Dead

Two separate problems, both understated by calling this "just add more backends." First, the scaling problem Chapter 95 opened with: one server has a ceiling on concurrent requests it can serve, and adding more servers only helps if something in front of them actually spreads requests across all of them — otherwise "adding a server" does nothing, because nothing ever sends it traffic. Second, and easy to underestimate until it happens in production: a backend can fail *silently* from the load balancer's point of view. A crashed process, a full disk, a deadlocked goroutine — the load balancer has no way to know unless it actively checks, and until it does, it keeps sending some fraction of all traffic straight into a black hole.

---

## 3. The Naive Shortcut We're Deliberately Not Taking

A first attempt might cycle through a fixed slice of backend addresses with a plain counter and nothing else:

```go
backends := []string{"a:9091", "b:9092", "c:9093"}
next := backends[i % len(backends)]
i++
```

This handles the *scaling* half of Section 2's problem, but does nothing for the *silent failure* half — if `b:9092` crashes, exactly one out of every three requests keeps getting routed straight to a dead backend forever, with no mechanism to notice or correct it. Chapter 95's health-check concept is not an optional add-on to load balancing; without it, a naive round-robin implementation actively makes an outage worse by continuing to hand a fixed fraction of traffic to a backend that will never answer.

---

## 4. The Real Solution: A Backend Pool, a Selection Strategy, and Active Health Checks

```
1. Model each backend as a struct carrying: its address, a live "is this
   backend healthy right now" flag, an active-connection counter, and its
   own pre-built ReverseProxy (Ch 112 Sec 7).
2. On every incoming request, ask the current strategy (round-robin or
   least-connections) to pick ONE currently-healthy backend.
3. Forward the request through that backend's ReverseProxy, tracking the
   in-flight connection count for least-connections to use.
4. Independently of request traffic, run a background goroutine that
   periodically probes every backend's health endpoint directly, updating
   each backend's healthy flag based on real, current results — not
   inferred from whether traffic happens to be flowing to it right now.
```

Step 4 is genuinely independent of steps 1-3: a backend can be marked unhealthy even if no client request has touched it recently, and — just as importantly — can be marked healthy again automatically once it starts responding, with no manual intervention.

---

## 5. Code: The Backend Pool and Round-Robin Selection

```go
type Backend struct {
	URL         *url.URL
	Proxy       *httputil.ReverseProxy
	Alive       int32 // 0 or 1, read/written atomically from multiple goroutines
	ActiveConns int64 // in-flight request count, for least-connections (Section 6)
}

func (b *Backend) IsAlive() bool { return atomic.LoadInt32(&b.Alive) == 1 }

func (b *Backend) SetAlive(alive bool) {
	var v int32
	if alive {
		v = 1
	}
	atomic.StoreInt32(&b.Alive, v)
}

type LoadBalancer struct {
	backends []*Backend
	counter  uint64 // monotonically increasing; round-robin position is counter % len(backends)
	strategy string // "round-robin" or "least-connections"
}

// pickRoundRobin cycles through backends in order, skipping unhealthy ones,
// and checks at most len(backends) candidates before giving up.
func (lb *LoadBalancer) pickRoundRobin() *Backend {
	n := len(lb.backends)
	for i := 0; i < n; i++ {
		idx := int(atomic.AddUint64(&lb.counter, 1) % uint64(n))
		if b := lb.backends[idx]; b.IsAlive() {
			return b
		}
	}
	return nil // every backend is currently unhealthy
}
```

`atomic.AddUint64` makes the counter safe under concurrent requests from many goroutines simultaneously (Chapter 106's goroutine-per-connection model, applied here to `http.Server`'s own goroutine-per-request handling) without needing an explicit mutex just to hand out the next index.

---

## 6. Code: Least-Connections Selection

Round-robin assumes every request costs roughly the same amount of backend time — true for many workloads, false whenever request cost varies widely (a slow report-generation endpoint next to a fast health check, for instance). Least-connections instead tracks, in real time, how many requests each backend is *currently* handling, and always routes to whichever healthy backend has the fewest:

```go
// pickLeastConnections scans every healthy backend and returns the one
// currently handling the fewest in-flight requests.
func (lb *LoadBalancer) pickLeastConnections() *Backend {
	var best *Backend
	bestConns := int64(math.MaxInt64)
	for _, b := range lb.backends {
		if !b.IsAlive() {
			continue
		}
		conns := atomic.LoadInt64(&b.ActiveConns)
		if conns < bestConns {
			bestConns, best = conns, b
		}
	}
	return best
}
```

This only works correctly if `ActiveConns` is incremented the instant a request is handed to a backend and decremented the instant that request completes — Section 8's `ServeHTTP` does exactly that with a `defer`, mirroring the same "increment on start, decrement on cleanup, guaranteed even on error" pattern Chapter 106 used for connection-count tracking.

---

## 7. Code: Active Health Checking in a Background Goroutine

```go
// healthCheckLoop pings every backend's /healthz endpoint on a fixed
// interval, independent of whatever real client traffic is or isn't
// flowing right now, and flips each backend's Alive flag based on the result.
func (lb *LoadBalancer) healthCheckLoop(interval time.Duration) {
	client := &http.Client{Timeout: 2 * time.Second}
	for {
		for _, b := range lb.backends {
			go func(b *Backend) {
				resp, err := client.Get(strings.TrimRight(b.URL.String(), "/") + "/healthz")
				alive := err == nil && resp.StatusCode == http.StatusOK
				if resp != nil {
					resp.Body.Close()
				}
				wasAlive := b.IsAlive()
				b.SetAlive(alive)
				if alive != wasAlive {
					log.Printf("health check: backend %s is now %s", b.URL, statusWord(alive))
				}
			}(b)
		}
		time.Sleep(interval)
	}
}

func statusWord(alive bool) string {
	if alive {
		return "UP"
	}
	return "DOWN"
}
```

Every backend is checked concurrently (each in its own goroutine) rather than one after another, so one slow or hung backend's 2-second timeout doesn't delay checking the others — a real consequence of Chapter 106's concurrency model applied to a background maintenance task instead of a client connection.

---

## 8. Code: The Complete Load Balancer

```go
// loadbalancer.go
package main

import (
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

// -- Backend, LoadBalancer, pickRoundRobin: Section 5 --
// -- pickLeastConnections: Section 6 --
// -- healthCheckLoop, statusWord: Section 7 --

func NewLoadBalancer(backendURLs []string, strategy string) *LoadBalancer {
	lb := &LoadBalancer{strategy: strategy}
	for _, raw := range backendURLs {
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("bad backend URL %q: %v", raw, err)
		}
		proxy := httputil.NewSingleHostReverseProxy(u)
		backend := &Backend{URL: u, Proxy: proxy}
		backend.SetAlive(true) // optimistic until the first health check (Section 7) says otherwise

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("backend %s error: %v", backend.URL, err)
			backend.SetAlive(false) // a failed proxied request is itself a health signal
			w.WriteHeader(http.StatusBadGateway)
		}
		lb.backends = append(lb.backends, backend)
	}
	return lb
}

func (lb *LoadBalancer) pick() *Backend {
	if lb.strategy == "least-connections" {
		return lb.pickLeastConnections()
	}
	return lb.pickRoundRobin()
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.pick()
	if backend == nil {
		http.Error(w, "503 Service Unavailable: no healthy backends", http.StatusServiceUnavailable)
		return
	}
	atomic.AddInt64(&backend.ActiveConns, 1)
	defer atomic.AddInt64(&backend.ActiveConns, -1)

	w.Header().Set("X-Balanced-By", backend.URL.Host) // visible proof of routing, Section 9
	backend.Proxy.ServeHTTP(w, r)
}

func main() {
	backends := []string{
		"http://localhost:9091",
		"http://localhost:9092",
		"http://localhost:9093",
	}
	lb := NewLoadBalancer(backends, "round-robin")
	go lb.healthCheckLoop(2 * time.Second)

	log.Println("load balancer listening on :8080, strategy:", lb.strategy)
	log.Fatal(http.ListenAndServe(":8080", lb))
}
```

`ServeHTTP`'s three lines around `backend.Proxy.ServeHTTP(w, r)` are the whole point of this chapter: increment the connection counter before forwarding, decrement it (via `defer`, guaranteed even if the backend request panics or errors) after, and let each backend's own pre-built `ReverseProxy` — already carrying Chapter 112's header-rewriting behavior — do the actual forwarding work.

---

## 9. Hands-On Experiment: Rotation, Failure, and Recovery

**Step 1 — a backend that reports its own port and a health endpoint:**

```go
// backend.go
package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "response from backend on port %s\n", port)
	})
	fmt.Println("backend listening on :" + port)
	http.ListenAndServe(":"+port, nil)
}
```

**Step 2 — start three backends and the load balancer:**

```
$ PORT=9091 go run backend.go &
$ PORT=9092 go run backend.go &
$ PORT=9093 go run backend.go &
$ go run loadbalancer.go
load balancer listening on :8080, strategy: round-robin
```

**Step 3 — send several requests and watch round-robin cycle through all three:**

```
$ for i in 1 2 3 4 5 6; do curl -s http://localhost:8080/; done
response from backend on port 9091
response from backend on port 9092
response from backend on port 9093
response from backend on port 9091
response from backend on port 9092
response from backend on port 9093
```

**Step 4 — kill the 9092 backend (`Ctrl+C` in its terminal) and wait a couple of seconds for the health check interval to catch it:**

```
health check: backend http://localhost:9092 is now DOWN
```

```
$ for i in 1 2 3 4; do curl -s http://localhost:8080/; done
response from backend on port 9091
response from backend on port 9093
response from backend on port 9091
response from backend on port 9093
```

`pickRoundRobin`'s health check inside the loop (Section 5) is what makes this work — the counter still advances past index 1 (`9092`) every third call, but `b.IsAlive()` returns `false` for it, so the loop keeps advancing to the next candidate instead of returning a dead backend.

**Step 5 — restart the 9092 backend and confirm it rejoins automatically, with no restart of the load balancer itself:**

```
$ PORT=9092 go run backend.go &
```
```
health check: backend http://localhost:9092 is now UP
```
```
$ for i in 1 2 3; do curl -s http://localhost:8080/; done
response from backend on port 9091
response from backend on port 9092
response from backend on port 9093
```

This is Chapter 95's health-check concept made fully concrete: no human intervention, no load balancer restart — a backend's own responsiveness is the only thing that ever changes its membership in the rotation.

---

## 10. Worked Example: Round-Robin vs. Least-Connections Under Uneven Load

Round-robin and least-connections agree when every request costs the same amount of backend time. They diverge sharply when it doesn't. Suppose three backends are already handling requests of very different durations, and a new request arrives:

| Backend | Currently handling | Round-robin's choice | Least-connections' choice |
|---|---|---|---|
| A (`:9091`) | 1 slow request (still running) | — | — |
| B (`:9092`) | 5 fast requests (still running) | — | — |
| C (`:9093`) | 0 requests | — | — |

Round-robin, mid-cycle, simply sends the new request to whichever backend is next in sequence — if A was last, B is next, *regardless of the fact that B already has five times the load of A or C*. Least-connections, reading `ActiveConns` for each backend via `pickLeastConnections` (Section 6), sees `A=1, B=5, C=0` and correctly routes the new request to C — the backend genuinely least busy right now. This is exactly why Chapter 95 introduced least-connections as a distinct strategy rather than a refinement of round-robin: it requires tracking real, live state (`ActiveConns`) that round-robin's stateless counter never needs, in exchange for materially better behavior whenever request costs vary.

Play the same scenario forward for ten more incoming requests, alternating instantaneously (an unrealistic but illustrative simplification — real requests take nonzero time to finish): round-robin keeps cycling `A, B, C, A, B, C, ...` regardless of outcome, so B keeps accumulating requests it's already struggling to drain, while least-connections continuously re-reads `ActiveConns` before every single decision and naturally steers subsequent requests toward whichever of A or C finishes first and drops back to a lower count — the strategy actively self-corrects toward balance, where round-robin has no feedback loop at all. The trade-off is cost: `pickRoundRobin` is an O(1) counter increment, while `pickLeastConnections` is an O(n) scan over every backend on every single request — for a small backend pool this is irrelevant, but a load balancer fronting hundreds of backends would need a more efficient data structure (a min-heap keyed by connection count, updated on each increment/decrement) to keep least-connections' per-request cost from growing with pool size.

---

## 11. Common Pitfalls in Hand-Rolled Load Balancing

- **Checking health only via passive observation of request failures, never actively.** `ErrorHandler`'s `backend.SetAlive(false)` (Section 8) is a useful *additional* signal, but relying on it alone means a backend that's been silently dead for hours with zero traffic routed to it (e.g., during a quiet period) is never detected until a real client request happens to hit it and fail — which is precisely the outage Section 2 described. `healthCheckLoop`'s independent, traffic-agnostic probing (Section 7) is what actually closes that gap.
- **Forgetting to decrement `ActiveConns` on every exit path.** Without the `defer` in `ServeHTTP` (Section 8), a panic or early return inside `Proxy.ServeHTTP` would leave the counter permanently inflated for that backend, making least-connections increasingly avoid a perfectly healthy backend forever — a slow, silent form of self-inflicted imbalance.
- **A round-robin loop with no bound on iterations.** `pickRoundRobin`'s `for i := 0; i < n; i++` cap (Section 5) is deliberate — without it, if every backend were simultaneously unhealthy, an unbounded loop incrementing the shared counter would spin forever instead of correctly returning "no healthy backends" so `ServeHTTP` can answer with `503`.
- **Racing the health check against `NewLoadBalancer`'s startup.** `backend.SetAlive(true)` optimistically before the first health check completes means a genuinely dead backend added at startup gets some (bounded, brief) amount of real traffic before the first health check interval elapses — an intentional trade-off (favor availability over waiting) that's worth calling out explicitly rather than leaving implicit.
- **Sharing one `http.Client` with a long timeout for health checks.** A health check timeout that's too long (or absent) makes the whole health-check loop, and therefore the promptness of removing a truly dead backend, only as fast as its slowest possible failure — Section 7's 2-second `client.Timeout` bounds this explicitly.

---

## 12. Production Notes: What Real Load Balancers Add

- **L4 vs. L7, revisited concretely (Chapter 95).** This chapter's load balancer operates entirely at Layer 7 — it terminates the client's HTTP connection, reads the request, and issues a *new* HTTP request to the chosen backend via `ReverseProxy`. A Layer 4 load balancer (like a plain TCP/UDP load balancer, or a cloud "Network Load Balancer") instead forwards raw packets or relays a TCP connection without ever parsing HTTP — faster and protocol-agnostic, but unable to make routing decisions based on URL path, headers, or cookies the way this chapter's L7 balancer could be extended to.
- **Sticky sessions.** Some applications need a given client to keep hitting the *same* backend across multiple requests (e.g., in-memory session state not shared across backends). Production load balancers implement this via a cookie or client-IP hash consulted before falling back to round-robin/least-connections — this chapter's balancer has no such mechanism and would need one added as a third strategy.
- **Weighted strategies.** Real deployments often run backends of different capacities (a bigger machine should get proportionally more traffic) — weighted round-robin and weighted least-connections extend Sections 5-6's logic with a per-backend weight multiplier. Concretely, a weighted version of `pickRoundRobin` doesn't cycle through the raw `backends` slice — it cycles through an expanded sequence built from each backend's weight, e.g. a weight-3 backend `A` and a weight-1 backend `B` produce the selection sequence `A, A, A, B, A, A, A, B, ...` (interleaved rather than grouped, so a burst of consecutive requests still spreads across backends instead of hammering `A` three times in a row), regenerated whenever a backend's health or weight changes.
- **Graceful draining.** Removing a backend for planned maintenance should stop sending it *new* requests while letting its *in-flight* requests finish — this chapter's health check is binary (up/down) and doesn't distinguish "stop routing new traffic here" from "kill everything immediately."
- **Passive health checks combined with active ones**, circuit-breaker patterns that back off retrying a recently-failed backend with increasing delay, and integration with service discovery (Chapter 101) so the backend pool itself updates automatically as instances scale up or down — all real, common production additions on top of this chapter's fixed backend list.

---

## 13. What's Simplified Here

This load balancer supports exactly two strategies (a real one might offer half a dozen), has a fixed, hardcoded backend list (no dynamic registration/deregistration via an API or service discovery), performs no TLS termination, has no sticky-session support, and treats health as strictly binary rather than tracking response-time-based "degraded" states. Each of these mirrors a named, real production feature in Section 12 rather than an arbitrary omission.

---

## Interview Questions & Model Answers

**Beginner: Why is round-robin alone, without health checks, actively dangerous rather than just suboptimal?**

Without health checks, round-robin has no way to know a backend has failed, so it keeps sending it an equal share of traffic forever — turning what should be "one backend down, the others absorb the load" into "some fixed fraction of all client requests fail permanently," which is worse for users than if the dead backend had simply never existed. Health checks are what let the load balancer notice a failure and route around it automatically.

**Intermediate: Using this chapter's code, explain exactly how `pickLeastConnections` avoids routing traffic to a backend that just crashed, even in the brief window before the next health check runs.**

`pickLeastConnections` only considers backends where `b.IsAlive()` returns true (Section 6) — it skips a backend entirely once its `Alive` flag is `false`, regardless of its current `ActiveConns` count. A backend that crashed mid-request will have its `ErrorHandler` invoked when `ReverseProxy` fails to get a response from it (Section 8's `NewLoadBalancer`), which calls `backend.SetAlive(false)` immediately as a *passive* signal — faster than waiting for the next scheduled active health check (Section 7) to catch it. Combining both mechanisms means a crash gets detected either the moment a real request happens to hit it (passive) or, at worst, within one health-check interval even with zero real traffic (active) — the two pitfalls named in Section 11 addressed together rather than relying on either alone.

**Advanced: This chapter's health check hits `/healthz` and only checks for a `200 OK` status code. Describe two concrete failure modes this would miss, and how you would extend the health check to catch them.**

First, a backend process can be running and returning `200 OK` from `/healthz` while the actual application logic behind other routes is broken — for example, a database connection pool exhausted or a critical dependency down, while the process itself and its HTTP listener are perfectly alive. A more meaningful health check would have `/healthz` itself verify the backend's critical dependencies (a lightweight database ping, for instance) before returning `200`, rather than being a static handler that always succeeds. Second, a backend can be alive and correctly answering `/healthz` quickly while being severely overloaded and responding to *real* traffic routes with multi-second latency — a binary up/down health signal never detects this "technically alive but effectively useless" state. Extending the health check to also track recent real-request latency or error rate (a "degraded" state between fully healthy and fully down) and factoring that into `pick()` — for example, deprioritizing but not fully excluding a degraded backend — would catch both gaps this simple binary check misses.

---

## Exercises

### Easy
1. Add a fourth backend to the pool and confirm round-robin's cycle length changes to match.
2. Change `main()`'s `strategy` to `"least-connections"` and re-run Section 9's experiment, observing that routing no longer strictly alternates in fixed order.
3. Print the current `ActiveConns` value for every backend once per health-check interval, to watch the counts change live as you send concurrent requests.

### Medium
4. Add a `/stats` endpoint on the load balancer itself that reports each backend's URL, alive status, and current `ActiveConns` as JSON.
5. Implement weighted round-robin: give each backend an integer weight, and make `pickRoundRobin` favor higher-weight backends proportionally (e.g., a weight-3 backend should be picked roughly 3x as often as a weight-1 backend).
6. Add a "degraded" third health state (distinct from alive/dead) based on a rolling average of recent response times, and have `pick()` avoid degraded backends only when at least one fully healthy backend is available.

### Hard
7. Implement sticky sessions via a cookie: on a backend's first response to a new client, set a cookie identifying which backend served it, and have `pick()` honor that cookie on subsequent requests from the same client as long as the backend remains healthy.
8. Add graceful draining: a `Drain(backend)` method that immediately stops new routing to a backend (skip it in `pick()`) while letting its current `ActiveConns` count fall to zero naturally before fully removing it, with a maximum wait timeout.
9. Replace the fixed backend list with dynamic registration: an HTTP endpoint (`POST /backends`) that adds a new backend to the live pool at runtime, safely under concurrent access from `ServeHTTP` and `healthCheckLoop` — consider what data structure changes are needed beyond a plain slice to do this without a lock contended on every single request.

---

## Summary

| Term | Meaning |
|---|---|
| Backend pool | A set of interchangeable servers a load balancer can route requests to |
| Round-robin | Cycling through backends in fixed order, skipping unhealthy ones |
| Least-connections | Always routing to whichever healthy backend currently has the fewest in-flight requests |
| Active health check | A background probe (this chapter: `GET /healthz` on a timer) independent of real client traffic |
| Passive health check | Marking a backend unhealthy because a real proxied request to it just failed |
| `ActiveConns` | A live, atomically-updated counter of in-flight requests per backend, read by least-connections |
| Sticky session | Routing a given client's requests to the same backend repeatedly, for session-state reasons (not implemented here — Section 12) |

You've now built the code-level realization of Chapter 95's load balancing concepts end to end: a pool of backends, two real selection strategies, and a health-check loop that adds and removes backends from rotation with no human in the loop — all built directly on Chapter 112's reverse proxy. Every chapter in this volume so far has worked entirely at or above the transport layer, taking IP and Ethernet framing for granted. Chapter 114 drops below that assumption for the first time, building a packet sniffer that captures and decodes raw Ethernet, IP, and TCP/UDP headers straight off a network interface — putting Volumes 5, 6, and 9's byte-level header diagrams directly into running code.
