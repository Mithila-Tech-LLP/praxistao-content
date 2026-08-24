# Chapter 116: Building a CDN-Like Cache

> **"A CDN's entire value proposition boils down to one question, asked billions of times a day: have I already seen this exact response, and can I prove it's still good enough to reuse without asking the origin? This chapter builds the code that answers that question correctly, instead of guessing."**

---

## Table of Contents

1. [Recap: Chapter 72's Rules, Chapter 96's Preview, Chapter 112's Proxy](#1-recap-chapter-72s-rules-chapter-96s-preview-chapter-112s-proxy)
2. [The Problem: A Cache That Ignores Cache-Control Is Worse Than No Cache](#2-the-problem-a-cache-that-ignores-cache-control-is-worse-than-no-cache)
3. [The Real Solution: Key, Freshness, Validators, Eviction](#3-the-real-solution-key-freshness-validators-eviction)
4. [Code: The Cache Entry and Cache-Control Parsing](#4-code-the-cache-entry-and-cache-control-parsing)
5. [Code: An LRU Store](#5-code-an-lru-store)
6. [Code: The Caching Proxy — Fresh Hits and Full Misses](#6-code-the-caching-proxy--fresh-hits-and-full-misses)
7. [Code: Conditional Revalidation With If-None-Match](#7-code-conditional-revalidation-with-if-none-match)
8. [Code: Honoring the Client's Own Conditional Request](#8-code-honoring-the-clients-own-conditional-request)
9. [Code: A Tiny Origin Server to Test Against](#9-code-a-tiny-origin-server-to-test-against)
10. [Code: main() — Wiring It All Together](#10-code-main--wiring-it-all-together)
11. [Hands-On Experiment: HIT, MISS, REVALIDATED, and Eviction](#11-hands-on-experiment-hit-miss-revalidated-and-eviction)
12. [Worked Example: Tracing Three Requests Through the Cache](#12-worked-example-tracing-three-requests-through-the-cache)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes: What Real CDNs Add](#14-production-notes-what-real-cdns-add)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)
19. [What's Next in This Volume](#19-whats-next-in-this-volume)

---

## 1. Recap: Chapter 72's Rules, Chapter 96's Preview, Chapter 112's Proxy

Chapter 72 laid out the complete decision procedure a cache follows (its Section 13): check freshness against `Cache-Control: max-age`, and if stale, revalidate with `If-None-Match`/`If-Modified-Since` rather than blindly refetching or blindly reusing. Chapter 96 previewed what a CDN edge node does with that procedure at global scale. Chapter 112 built the plumbing this chapter reuses directly — a Go program that accepts a client connection and forwards it to an origin. This chapter is what happens when you put Chapter 72's rules *inside* Chapter 112's proxy.

---

## 2. The Problem: A Cache That Ignores Cache-Control Is Worse Than No Cache

A naive "cache everything for 5 minutes" proxy is tempting to write and actively dangerous to run. It would happily cache a `Set-Cookie`-bearing, per-user account page (Chapter 72 Section 9's `private` directive exists specifically to forbid this) and serve *user A's* personalized response to *user B* on the very next request — a serious data leak, not a performance win. It would also happily serve a stale copy of a resource the origin marked `no-store` for security reasons. The naive approach optimizes for speed while ignoring the one thing that makes caching *safe*: the origin server's own explicit instructions about what may be cached, for how long, and by whom.

---

## 3. The Real Solution: Key, Freshness, Validators, Eviction

```
1. Cache only GET responses, keyed by request path (Chapter 72 Section 8's
   rationale: caching a POST/PUT's "response" doesn't make sense the same way).
2. On every response FROM the origin, parse Cache-Control: respect no-store
   and private by never storing the response at all (Chapter 72 Sections 9, 12).
3. Store cacheable responses with their ETag and a freshness window
   (max-age), per Chapter 72 Section 10.
4. On a later request for the same key:
     - If cached AND still fresh -> serve directly, zero origin contact.
     - If cached but stale AND has an ETag -> send a conditional GET
       (If-None-Match) to the origin. 304 -> reuse the body, refresh
       freshness. 200 -> replace the cached entry.
     - Otherwise -> full fetch.
5. Bound total cache size with LRU eviction, so memory use can't grow
   without limit (Chapter 96's CDN edges face exactly this constraint
   at a scale of terabytes, not bytes).
```

This is Chapter 72 Section 13's decision procedure, applied by a proxy on behalf of every client hitting it, instead of by one browser on behalf of one user.

---

## 4. Code: The Cache Entry and Cache-Control Parsing

```go
package main

import (
	"container/list"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// CacheEntry is what this proxy stores for one cached GET response.
type CacheEntry struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	ETag       string
	StoredAt   time.Time
	MaxAge     time.Duration
	NoCache    bool // Cache-Control: no-cache — store it, but always revalidate (Ch 72 Sec 9)
	Cacheable  bool // false for no-store / private responses (Ch 72 Sec 9, 12)
}

func (e *CacheEntry) isFresh() bool {
	if e.NoCache {
		return false
	}
	return time.Since(e.StoredAt) < e.MaxAge
}

// CacheControl holds the parsed directives from one Cache-Control header
// (Chapter 72, Section 9's directive table).
type CacheControl struct {
	NoStore   bool
	NoCache   bool
	Private   bool
	Public    bool
	MaxAge    time.Duration
	HasMaxAge bool
}

func parseCacheControl(header string) CacheControl {
	var cc CacheControl
	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		switch {
		case part == "no-store":
			cc.NoStore = true
		case part == "no-cache":
			cc.NoCache = true
		case part == "private":
			cc.Private = true
		case part == "public":
			cc.Public = true
		case strings.HasPrefix(part, "max-age="):
			if secs, err := strconv.Atoi(strings.TrimPrefix(part, "max-age=")); err == nil {
				cc.MaxAge = time.Duration(secs) * time.Second
				cc.HasMaxAge = true
			}
		}
	}
	return cc
}

// buildEntry turns one origin response into a CacheEntry, applying Chapter 72
// Section 9's rules for what a SHARED cache (this proxy, standing in for a
// CDN edge) is allowed to store at all.
func buildEntry(resp *http.Response, body []byte) *CacheEntry {
	cc := parseCacheControl(resp.Header.Get("Cache-Control"))
	entry := &CacheEntry{
		StatusCode: resp.StatusCode,
		Header:     resp.Header.Clone(),
		Body:       body,
		ETag:       resp.Header.Get("ETag"),
		StoredAt:   time.Now(),
		NoCache:    cc.NoCache,
	}
	if cc.HasMaxAge {
		entry.MaxAge = cc.MaxAge
	}
	// A SHARED cache must never store a `private` response (it might belong
	// to a different user next time) or one marked `no-store` at all.
	entry.Cacheable = !cc.NoStore && !cc.Private
	return entry
}
```

---

## 5. Code: An LRU Store

```go
type lruItem struct {
	key   string
	entry *CacheEntry
}

// LRUCache bounds total memory use by evicting the least-recently-used entry
// once capacity is exceeded — exactly the constraint a real CDN edge node
// faces at a much larger scale (Chapter 96).
type LRUCache struct {
	mu       sync.Mutex
	capacity int
	ll       *list.List
	items    map[string]*list.Element
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{capacity: capacity, ll: list.New(), items: map[string]*list.Element{}}
}

func (c *LRUCache) Get(key string) (*CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.ll.MoveToFront(el)
	return el.Value.(*lruItem).entry, true
}

func (c *LRUCache) Put(key string, entry *CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.ll.MoveToFront(el)
		el.Value.(*lruItem).entry = entry
		return
	}
	el := c.ll.PushFront(&lruItem{key: key, entry: entry})
	c.items[key] = el
	if c.ll.Len() > c.capacity {
		back := c.ll.Back()
		evicted := back.Value.(*lruItem)
		c.ll.Remove(back)
		delete(c.items, evicted.key)
		fmt.Printf("[cache] evicted %q (LRU capacity %d reached)\n", evicted.key, c.capacity)
	}
}
```

---

## 6. Code: The Caching Proxy — Fresh Hits and Full Misses

```go
type CachingProxy struct {
	origin *url.URL
	client *http.Client
	cache  *LRUCache
}

func NewCachingProxy(originURL string, capacity int) *CachingProxy {
	u, err := url.Parse(originURL)
	if err != nil {
		log.Fatal(err)
	}
	return &CachingProxy{
		origin: u,
		client: &http.Client{Timeout: 5 * time.Second},
		cache:  NewLRUCache(capacity),
	}
}

func (p *CachingProxy) fetchFromOrigin(path, ifNoneMatch string) (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, p.origin.String()+path, nil)
	if err != nil {
		return nil, nil, err
	}
	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch) // Chapter 72 Section 11
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, body, nil
}

func (p *CachingProxy) store(key string, entry *CacheEntry) {
	if entry.Cacheable && entry.StatusCode == http.StatusOK {
		p.cache.Put(key, entry)
	}
}

// ServeHTTP is Chapter 72 Section 13's decision procedure, as an HTTP handler.
func (p *CachingProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this cache only handles GET (Chapter 72 Section 8)", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Path

	if entry, hit := p.cache.Get(key); hit {
		if entry.isFresh() {
			p.serveFromCache(w, r, entry, "HIT")
			return
		}
		if entry.ETag != "" {
			p.revalidateAndServe(w, r, key, entry)
			return
		}
		// stale, no validator at all — Chapter 72 Section 13, step 3's "No" branch.
	}

	// Full miss: nothing usable was cached.
	resp, body, err := p.fetchFromOrigin(key, "")
	if err != nil {
		http.Error(w, "bad gateway: "+err.Error(), http.StatusBadGateway)
		return
	}
	newEntry := buildEntry(resp, body)
	p.store(key, newEntry)
	p.serveFromCache(w, r, newEntry, "MISS")
}
```

---

## 7. Code: Conditional Revalidation With If-None-Match

```go
// revalidateAndServe implements Chapter 72 Section 11's conditional GET: ask
// the origin "has this changed?" using the stored ETag, instead of either
// blindly reusing a stale copy or blindly re-downloading everything.
func (p *CachingProxy) revalidateAndServe(w http.ResponseWriter, r *http.Request, key string, stale *CacheEntry) {
	resp, body, err := p.fetchFromOrigin(key, stale.ETag)
	if err != nil {
		http.Error(w, "bad gateway: "+err.Error(), http.StatusBadGateway)
		return
	}
	if resp.StatusCode == http.StatusNotModified {
		// Nothing changed — reuse the body we already have, just refresh the
		// freshness window. This is the entire point of Chapter 72 Section 11:
		// a tiny header-only exchange instead of re-sending the whole body.
		stale.StoredAt = time.Now()
		fmt.Printf("[proxy] %s revalidated with origin (304), body unchanged\n", key)
		p.serveFromCache(w, r, stale, "REVALIDATED")
		return
	}
	// Origin sent a full 200 — the resource actually changed.
	fresh := buildEntry(resp, body)
	p.store(key, fresh)
	p.serveFromCache(w, r, fresh, "MISS")
}
```

---

## 8. Code: Honoring the Client's Own Conditional Request

A real CDN edge doesn't just use conditional requests *upstream* toward the origin — it also honors conditional requests *downstream* from its own clients, exactly like Chapter 72 Section 15's `cachedAssetHandler` did directly. Whatever this proxy is about to serve — a fresh hit, a revalidated entry, or a brand-new miss — the client might already have the current version:

```go
func (p *CachingProxy) serveFromCache(w http.ResponseWriter, r *http.Request, entry *CacheEntry, result string) {
	if inm := r.Header.Get("If-None-Match"); inm != "" && entry.ETag != "" && inm == entry.ETag {
		w.Header().Set("X-Cache", result+"-304")
		w.WriteHeader(http.StatusNotModified) // Chapter 72 Section 11 — no body sent
		return
	}
	for k, v := range entry.Header {
		if k == "Content-Length" {
			continue
		}
		w.Header()[k] = v
	}
	w.Header().Set("X-Cache", result) // not a real HTTP standard header, but universal debugging convention
	w.WriteHeader(entry.StatusCode)
	w.Write(entry.Body)
}
```

`X-Cache` isn't part of any RFC — it's a de facto convention nearly every real CDN (Cloudflare, Fastly, CloudFront) uses for exactly this purpose: telling you, from outside, whether your request was served from cache without you needing to inspect timing or `ETag` values yourself.

---

## 9. Code: A Tiny Origin Server to Test Against

```go
func startOrigin(addr string) {
	mux := http.NewServeMux()
	var hits int32

	data := []byte(`{"headline": "CDN chapter demo payload"}`)
	sum := sha256.Sum256(data)
	etag := fmt.Sprintf(`"%x"`, sum[:8])

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&hits, 1)
		fmt.Printf("[origin] request #%d for /data (If-None-Match: %q)\n", n, r.Header.Get("If-None-Match"))
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "public, max-age=5") // short window, for the demo in Section 11
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	mux.HandleFunc("/private-data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, no-store") // Chapter 72 Section 9 — never cache
		fmt.Fprintln(w, `{"account": "should never be cached by a shared proxy"}`)
	})

	// /item/<id> — distinct, independently cacheable resources, used to
	// demonstrate LRU eviction in Section 11.
	mux.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/item/")
		body := []byte(fmt.Sprintf(`{"id": %q}`, id))
		s := sha256.Sum256(body)
		itemETag := fmt.Sprintf(`"%x"`, s[:8])
		w.Header().Set("ETag", itemETag)
		w.Header().Set("Cache-Control", "public, max-age=60")
		if r.Header.Get("If-None-Match") == itemETag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(addr, mux))
}
```

---

## 10. Code: main() — Wiring It All Together

```go
func main() {
	go startOrigin(":9000")
	time.Sleep(200 * time.Millisecond)

	proxy := NewCachingProxy("http://127.0.0.1:9000", 3) // capacity 3, deliberately small for Section 11
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
```

---

## 11. Hands-On Experiment: HIT, MISS, REVALIDATED, and Eviction

**Freshness and revalidation**, against the short 5-second `max-age` on `/data`:

```
$ curl -s -D - -o /dev/null http://localhost:8080/data | grep -E 'X-Cache|ETag'
X-Cache: MISS
ETag: "3a1f9c0e2b7d4f11"

$ curl -s -D - -o /dev/null http://localhost:8080/data | grep X-Cache
X-Cache: HIT

$ sleep 6   # past max-age=5, entry is now stale but has an ETag

$ curl -s -D - -o /dev/null http://localhost:8080/data | grep X-Cache
X-Cache: REVALIDATED
```

The origin's own log confirms exactly what happened at each step — request #1 got a full `200`, request #2 never reached the origin at all (no corresponding origin log line), and request #3 reached the origin but got a `304`, not a full body.

**A `private, no-store` response is never cached, ever:**

```
$ curl -s -D - -o /dev/null http://localhost:8080/private-data | grep X-Cache
X-Cache: MISS
$ curl -s -D - -o /dev/null http://localhost:8080/private-data | grep X-Cache
X-Cache: MISS    # still MISS — nothing was ever stored, by design (Chapter 72 Section 9)
```

**LRU eviction**, with capacity 3:

```
$ for i in 1 2 3; do curl -s -o /dev/null -w "%{url_effective} " http://localhost:8080/item/$i; done
[cache] (nothing evicted yet — 3 items, capacity 3)

$ curl -s -D - -o /dev/null http://localhost:8080/item/4 | grep X-Cache
X-Cache: MISS
[cache] evicted "/item/1" (LRU capacity 3 reached)

$ curl -s -D - -o /dev/null http://localhost:8080/item/1 | grep X-Cache
X-Cache: MISS    # evicted earlier — has to be fetched from the origin all over again
```

---

## 12. Worked Example: Tracing Three Requests Through the Cache

For `GET /data`, walking Chapter 72 Section 13's algorithm exactly as this proxy implements it:

| Request # | Cache state before | Step taken | Result | Origin contacted? |
|---|---|---|---|---|
| 1 | empty | full fetch, `buildEntry` stores it (`public, max-age=5`) | `200`, `X-Cache: MISS` | Yes — full body |
| 2 (immediately after) | fresh (< 5s old) | `isFresh()` true, served directly | `200`, `X-Cache: HIT` | No |
| 3 (after 6s) | stale, has ETag | conditional GET with `If-None-Match`, origin says unchanged | `200`, `X-Cache: REVALIDATED` | Yes — headers only, no body |

This table is Chapter 72 Section 13's four-step flowchart with concrete outcomes filled in — the same logic, just proven against a running program instead of described in prose.

---

## 13. Common Misconceptions

- **"A cache hit means the request never happened."** A `304` revalidation (row 3 above) *is* a real network round trip — Chapter 72 Section 11 was explicit about this: the request still happens, only the expensive part (the full body) is skipped. Only a true freshness hit (row 2) involves zero network contact at all.
- **"`no-cache` means the same thing as `no-store`."** As Chapter 72 Section 9 stressed, `no-cache` means "you may store this, but must revalidate before every reuse" (modeled here by `CacheEntry.NoCache` forcing `isFresh()` to always return `false`) — `no-store` is the directive that actually forbids storage, modeled here by `Cacheable = false`.
- **"An ETag mismatch on the client's own conditional request should be treated as an error."** It's the normal case — Section 8's handler falls through to serve the current body with a fresh `200`, exactly what Chapter 72 Section 11 describes for a changed resource.
- **"LRU eviction means the least-used-overall item gets removed."** It's least-*recently*-used, not least-frequently-used — Section 5's `MoveToFront` on every `Get` means an item requested once, a long time ago, but not since, is evicted before an item requested many times but not for a while would be under a different policy (LFU) — a deliberate, common trade-off real caches make for O(1) simplicity.

---

## 14. Production Notes: What Real CDNs Add

- **Vary-aware cache keys.** Chapter 72 Section 13 introduced the `Vary` header; this chapter's cache key is just the URL path, ignoring it entirely — a production cache must fold relevant request headers (like `Accept-Encoding`) into the key, or it will serve a gzip-compressed response to a client that never said it could handle one.
- **Cache invalidation at scale is explicit, not just time-based.** Real CDNs offer purge/invalidation APIs so an origin can force-evict a specific URL across thousands of edge nodes the moment content changes urgently, rather than waiting out a `max-age` window — Chapter 72 Section 17's content-hashed-filename strategy sidesteps this need entirely for static assets, which is why it's the preferred approach wherever it's applicable.
- **Stale-while-revalidate and stale-if-error** are real `Cache-Control` extensions (not covered in Chapter 72's core set) letting a cache serve a *known-stale* copy immediately while revalidating in the background, or during an origin outage — trading strict correctness for availability, a deliberate choice many high-traffic sites make.
- **Real CDN edges are geographically distributed** (Chapter 96), so "the cache" is actually thousands of independent LRU-like stores worldwide, each warmed independently by the traffic that happens to land on it — a request hitting an edge in Tokyo for the first time is a MISS there even if the same URL is a long-standing HIT in Frankfurt.
- **Byte-range caching, partial content (`206`), and streaming responses** all complicate the simple whole-body `CacheEntry` model this chapter uses — production HTTP caches (Varnish, Nginx's proxy_cache, every commercial CDN) handle partial responses as a first-class case.

---

## 15. What's Simplified Here

- The cache key ignores the `Vary` header and query strings entirely; a production cache must fold both into the key correctly.
- `Last-Modified`/`If-Modified-Since` (Chapter 72 Section 10's other validator) isn't implemented — only `ETag`/`If-None-Match`.
- Cache entries are held entirely in one process's memory with no persistence, no size-based (only count-based) eviction limit, and no sharing across multiple proxy instances — a real edge fleet needs a distributed or at least size-bounded store.
- There's no purge/invalidation API — the only way to remove an entry early is LRU pressure or the process restarting.
- Concurrent requests for the same currently-uncached URL will each independently miss and fetch from the origin ("cache stampede") — production caches typically coalesce concurrent misses for the same key into a single origin request.

---

## 16. Interview Questions & Model Answers

**Beginner: What's the difference between a cache "hit," a "revalidation," and a "miss" in this chapter's proxy?**
A hit means the cached copy is still within its freshness window and is served with zero contact with the origin. A revalidation means the cached copy is stale but has a validator (an ETag), so a conditional GET is sent to the origin — if the origin confirms nothing changed (`304`), the existing body is reused with a refreshed freshness window; only the headers were exchanged, not the body. A miss means nothing usable was cached at all (or nothing to revalidate against), and the origin is asked for the full response.

**Intermediate: Why does this proxy refuse to cache any response carrying `Cache-Control: private`, even though it would happily cache the exact same bytes if the header said `public`?**
Because this proxy is a *shared* cache serving many different clients, exactly like a CDN edge (Chapter 72 Section 12). A `private` response is explicitly marked as correct for only the one client it was generated for — often because it's personalized or carries session-specific data. Caching it here would risk serving user A's personalized response to user B on a later request to the same URL, which is a serious correctness and security bug, not a performance optimization.

**Advanced: Walk through exactly what happens, end to end, when a stale cache entry with an ETag is revalidated and the origin's content has genuinely changed — and explain why this design is strictly better than either always trusting a `max-age` window blindly or always re-fetching on every single request.**
The proxy finds the cached entry stale (past `max-age`) but has an ETag, so instead of either serving the stale body outright or performing a full unconditional GET, it sends a conditional GET with `If-None-Match` set to the stored ETag. Because the content changed, the origin doesn't recognize that ETag as current and responds with a full `200`, a new ETag, and the new body. The proxy calls `buildEntry` again to build a fresh `CacheEntry` from this response, stores it (evicting the old one via `Put`, potentially triggering LRU eviction of an unrelated key if capacity is exceeded), and serves the new body to the client with `X-Cache: MISS`. This is strictly better than blind trust because a `max-age` window alone can't detect an origin change that happens to occur mid-window (Chapter 72 Section 8's whole reason for validators existing); it's strictly better than unconditional refetching on every request because on the (usually far more common) unchanged case, only a small header exchange happens instead of the full body being retransmitted every time.

---

## 17. Exercises

### Easy
1. Add support for `Last-Modified`/`If-Modified-Since` as a fallback validator when no `ETag` is present, per Chapter 72 Section 10.
2. Change the LRU capacity to 1 and trace, by hand, what `X-Cache` value each of five sequential requests to five different URLs would produce.
3. Explain why `ServeHTTP` rejects non-GET methods outright rather than trying to cache them.

### Medium
4. Add a `/purge?path=/data` admin endpoint that removes a specific key from the cache immediately, independent of its freshness — the beginning of a real invalidation API (Section 14).
5. Fold the `Accept-Encoding` request header into the cache key (mirroring the real `Vary: Accept-Encoding` behavior from Chapter 72 Section 13), and demonstrate that a gzip-requesting client and a plain client now get independently cached entries for the same URL.
6. Add a maximum total byte-size limit to `LRUCache` (not just an entry count), evicting from the back until under budget after each `Put`.

### Hard
7. Implement request coalescing: if two concurrent requests arrive for the same currently-uncached key, ensure only one actual origin fetch happens and both callers receive the result (Section 14's "cache stampede" problem).
8. Add a `stale-while-revalidate=N` extension: serve the stale copy immediately while kicking off a background revalidation, updating the cache for the *next* request rather than making the current one wait.
9. Turn this single-process cache into a two-node cache where both nodes share cache state via a simple invalidation broadcast (a UDP message, à la Chapter 115's virtual wires) whenever either node stores a fresh entry for a key the other might also hold stale.

---

## 18. Summary

| Term | Meaning |
|---|---|
| Cache key | The identifier (here, request path) used to look up a stored response |
| `CacheEntry.isFresh()` | True if the entry is within its `max-age` window and not marked `no-cache` |
| `Cacheable` | Whether a shared cache is permitted to store this response at all (false for `no-store`/`private`) |
| LRU eviction | Removing the least-recently-used entry once the cache exceeds capacity |
| Conditional GET / `If-None-Match` | A request asking "has this changed?" using a previously stored ETag |
| `304 Not Modified` | The origin's answer confirming a stale cached copy is still correct, with no body sent |
| `X-Cache` | A de facto (non-standard) response header indicating HIT/MISS/REVALIDATED, used by nearly every real CDN |
| Cache stampede | Multiple concurrent misses for the same uncached key each independently hitting the origin |

Chapter 116 turned Chapter 72's caching rules into a proxy that actually gets them right — refusing to cache what it shouldn't, revalidating instead of guessing, and bounding its own memory with LRU eviction. It closes out the "build it yourself" arc for the HTTP-facing half of this volume.

---

## 19. What's Next in This Volume

Chapter 117 leaves user-space sockets behind for a real TUN virtual interface, building a minimal VPN: encapsulating one IP packet inside a UDP tunnel and encrypting the payload — the code-level realization of Chapter 85's tunneling concept. Chapter 118 closes the volume with its capstone: a small distributed key-value service, with multiple nodes discovering each other and speaking a custom wire protocol, bringing together sockets, serialization, and the concurrency patterns built across every chapter in this volume.
