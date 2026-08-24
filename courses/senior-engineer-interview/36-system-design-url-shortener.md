# Chapter 36: System Design — URL Shortener & Pastebin

URL shortener (like bit.ly) and Pastebin (like pastebin.com) are the most common starter system design questions. Mastering these teaches the core patterns that apply to every system design problem.

## Table of Contents

1. [Requirements](#1-requirements)
2. [Capacity Estimation](#2-capacity-estimation)
3. [High-Level Design](#3-high-level-design)
4. [Short URL Generation — The Core Problem](#4-short-url-generation--the-core-problem)
5. [Data Model & API](#5-data-model--api)
6. [Database & Caching Strategy](#6-database--caching-strategy)
7. [Analytics Tracking](#7-analytics-tracking)
8. [Deep Dives](#8-deep-dives)
9. [Pastebin Variations](#9-pastebin-variations)
10. [Summary](#summary)

---

## 1. Requirements

### Functional Requirements
```
Core:
  - Given a long URL, create a short URL (e.g., short.ly/abc123)
  - When a user visits the short URL, redirect to the long URL
  
Optional (clarify with interviewer):
  - Custom short URLs (user-chosen aliases)
  - Expiration dates for links
  - Analytics: click count, geographic breakdown, time-based data
  - Link deletion/deactivation
  - User accounts to manage their links
```

### Non-Functional Requirements
```
- High availability: a broken redirect = failed product
- Low latency: redirects should be fast (<100ms)
- Links should persist for years (high durability)
- Read-heavy: 100:1 read/write ratio (many more clicks than link creations)
```

---

## 2. Capacity Estimation

```
Scale assumptions:
  100M URLs created per day
  100:1 read/write ratio → 10B redirects per day

Writes:
  100M/day ÷ 86,400 s/day = ~1,200 writes/second

Reads:
  10B/day ÷ 86,400 = ~115,000 reads/second

Storage:
  Each URL record ≈ 500 bytes (short URL + long URL + metadata)
  100M URLs/day × 500 bytes = 50GB/day
  5 years: 50GB × 365 × 5 = ~90TB

Bandwidth:
  Reads: 115,000/s × 500 bytes = 57MB/s outbound
  Easily handled by a CDN/load balancer

Conclusion:
  Write volume (1,200/s) is easy — single primary DB
  Read volume (115,000/s) requires Redis caching
  Storage (90TB in 5 years) — use object storage or distributed DB
```

---

## 3. High-Level Design

```
[Browser] → [DNS] → [Load Balancer]
                           │
              ┌────────────┼────────────┐
              │            │            │
       [API Servers] [API Servers] [API Servers]
              │
    ┌─────────┼─────────┐
    │         │         │
[Redis   [PostgreSQL] [Analytics
 Cache]  (Primary)    Service]
              │
        [PostgreSQL
          Replicas]
```

**Create short URL flow:**
```
POST /api/v1/shorten
1. API server validates request
2. Generate unique short code
3. Write to PostgreSQL (primary)
4. Invalidate any related cache entries
5. Return short URL
```

**Redirect flow:**
```
GET /abc123
1. Check Redis cache for abc123 → long URL
2. If cache hit: HTTP 301/302 redirect (< 1ms)
3. If cache miss: query PostgreSQL → cache result → redirect
4. Increment click counter (async, via Kafka)
```

---

## 4. Short URL Generation — The Core Problem

The most interesting part of URL shortener design: how to generate a unique 7-character code for every URL.

### Approach 1: Random + Check for Collision

```go
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 62 characters, 7 length = 62^7 = 3.5 trillion unique codes

func generateShortCode() string {
    b := make([]byte, 7)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// Insert and handle collision:
func createShortURL(ctx context.Context, db *sql.DB, longURL string) (string, error) {
    for attempts := 0; attempts < 5; attempts++ {
        code := generateShortCode()
        _, err := db.ExecContext(ctx,
            "INSERT INTO urls (short_code, long_url) VALUES ($1, $2) ON CONFLICT DO NOTHING",
            code, longURL)
        if err == nil {
            return code, nil
        }
        // ON CONFLICT DO NOTHING means 0 rows affected = collision, retry
    }
    return "", errors.New("failed to generate unique code after 5 attempts")
}
```

**Collision probability:** With 3.5 trillion codes and 100M URLs, collision rate is extremely low. But "extremely low" != "never" in distributed systems.

### Approach 2: Hash-Based

```go
import "crypto/md5"

func hashToShortCode(longURL string) string {
    hash := md5.Sum([]byte(longURL))
    // MD5 = 128 bits, encode first 43 bits as base62
    encoded := base62Encode(hash[:])
    return encoded[:7] // take first 7 characters
}

// Problem: same long URL always generates same short code (deterministic)
// This is actually useful! Same URL won't have multiple short codes.
// But two different URLs could collide (MD5 collision, though rare)
```

### Approach 3: Base62 Counter (Best for High Scale)

```go
// Use a global counter (or per-shard counter) and convert to base62
// Counter 0 → "aaaaaaa"
// Counter 1 → "aaaaaab"
// ...
// Counter 3.5 trillion → "zzzzzzz"

// Counter stored in: database sequence, Redis INCR, or a counter service (Snowflake-like)

func counterToBase62(n int64) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    if n == 0 { return string(charset[0]) }
    
    result := []byte{}
    for n > 0 {
        result = append([]byte{charset[n%62]}, result...)
        n /= 62
    }
    return string(result)
}

// Get next ID from database sequence (atomic, no collision):
func nextShortCode(ctx context.Context, db *sql.DB) (string, error) {
    var id int64
    err := db.QueryRowContext(ctx, "SELECT nextval('url_id_seq')").Scan(&id)
    if err != nil { return "", err }
    return counterToBase62(id), nil
}
```

**Trade-offs:**
| Method | Pros | Cons |
|---|---|---|
| Random | Simple, no coordination | Collisions, needs check |
| Hash | Deterministic, same URL = same code | Hash collisions, hard to guarantee uniqueness |
| Counter | No collisions, predictable | Requires coordination (single point) |

---

## 5. Data Model & API

```sql
-- Core table:
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    user_id BIGINT,                        -- null for anonymous
    expires_at TIMESTAMP,                  -- null = never expires
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_short_code ON urls(short_code);

-- For analytics:
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL,
    clicked_at TIMESTAMP DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT
);
CREATE INDEX idx_clicks_code_time ON clicks(short_code, clicked_at);
```

```
APIs:
POST /api/v1/urls
  Request:  { long_url: string, custom_alias?: string, expires_in?: int }
  Response: { short_url: string, short_code: string, expires_at?: string }
  
GET /{short_code}   (the redirect endpoint)
  Response: 301 Moved Permanently with Location header
            (Use 302 for trackable analytics, 301 for browser caching)

DELETE /api/v1/urls/{short_code}
  Response: 204 No Content

GET /api/v1/urls/{short_code}/stats
  Response: { total_clicks: int, clicks_by_day: [], top_referrers: [] }
```

---

## 6. Database & Caching Strategy

```
Redis as cache for redirects:
  Key:   short_code (e.g., "abc123")
  Value: long_url
  TTL:   24 hours (refresh on access)
  
  Lookups: Redis → PostgreSQL (if miss) → cache → redirect
  Write-through: on URL creation, write to both DB and cache

Cache invalidation:
  On URL deletion: DEL short_code from Redis
  On URL expiry: use Redis TTL to auto-expire (align with URL expiry)

301 vs 302:
  301 Permanent Redirect: browser caches it → no future requests to our server
    Pros: reduces server load
    Cons: we can't update the redirect or track clicks (browser uses cache)
  
  302 Temporary Redirect: browser always checks our server
    Pros: can track every click, can change the destination
    Cons: more server load
  
  For URL shorteners with analytics: use 302.
  For pure redirect performance: use 301.
```

---

## 7. Analytics Tracking

```
Two approaches for tracking clicks:

1. Synchronous (simple but adds latency to redirects):
   user clicks → API writes click record → then redirects
   Problem: database write adds 5-10ms to every redirect
   
2. Asynchronous (recommended for high scale):
   user clicks → immediate 302 redirect
   API server fires async event to Kafka "clicks" topic
   Analytics consumers read from Kafka and write to ClickHouse/BigQuery
   
   Trade-off: slight delay in analytics, but redirect is instant
   
For real-time counters (total click count):
  Redis INCR short_code:clicks   # increment counter in Redis
  Async worker periodically flushes Redis counts to PostgreSQL
```

---

## 8. Deep Dives

### Custom Short URLs
```
User picks "my-company.com/promo" instead of random code
Challenges:
  - Conflict detection: "promo" might already be taken
  - Bad words/reserved names: "/admin", "/api", "/login" should be blocked
  - Length limits: still need to fit URL constraints

Implementation: reserve namespace with a blocklist; check existence before inserting
```

### URL Expiration
```
Option 1: TTL in database. Cron job queries WHERE expires_at < NOW() and deletes.
  Problem: scanning the whole table every minute is expensive.
  
Option 2: Partition by month. Drop entire partitions at month boundaries.
  Efficient bulk deletes. Only works if expiry is by month.
  
Option 3: Lazy expiration. Check expires_at at redirect time. If expired, return 410 Gone.
  Simple. Expired records accumulate in the database (run cleanup async).
  
Option 4: Redis TTL. Set TTL on cache entry = URL expiry. When it expires from cache,
  the redirect hits DB where we can check and return 410.
```

### Geographic Distribution
```
DNS → Route to nearest region (Route 53 geolocation routing)
Each region has its own API servers + Redis cache
Primary PostgreSQL in one region, read replicas in others
Short codes are globally unique (single counter service or UUID-based)
```

---

## 9. Pastebin Variations

Pastebin stores text (code snippets, logs) with a URL. Same architecture, different storage:

```
Difference from URL shortener:
  - Content is the text itself (not a redirect destination)
  - Content can be large (up to 10MB)
  - Need syntax highlighting (metadata: language)
  - Private pastes (access control)

Storage:
  Database: stores metadata (paste_id, title, language, created_at, expires_at, user_id)
  Object storage (S3): stores the actual paste content
  Why S3? Text can be MB+ sized — too large for database rows. S3 is cheap and scalable.

Serving pastes:
  Metadata (from PostgreSQL) + content (from S3) = two fetches
  Cache in Redis for hot pastes
  CDN for public pastes with long cache TTLs (most pastes don't change)
```

---

## Summary

- **Core flow:** generate unique short code → store in DB → cache in Redis → redirect via 302.
- **Short code generation:** counter + base62 is most reliable (no collisions, no coordination needed at moderate scale).
- **Caching:** Redis on redirect path eliminates DB reads for hot links. TTL = URL expiry.
- **Analytics:** async (Kafka + analytics DB) to avoid slowing down redirects.
- **301 vs 302:** 301 for performance (browser caches), 302 for analytics (server sees every click).
- **Pastebin:** same architecture but content goes to S3, metadata to PostgreSQL.
- Deep dive topics: custom aliases, URL expiry, geographic distribution, abuse prevention.
