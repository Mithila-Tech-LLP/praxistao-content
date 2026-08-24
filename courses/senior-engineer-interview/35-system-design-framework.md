# Chapter 35: System Design Framework & Interview Strategy

System design interviews test whether you can take a vague requirement ("design Twitter") and turn it into a concrete, working architecture. The framework in this chapter is used by engineers at every top company — it works for any system design question.

## Table of Contents

1. [The 45-Minute Framework](#1-the-45-minute-framework)
2. [Requirements Clarification — The Most Important Step](#2-requirements-clarification--the-most-important-step)
3. [Capacity Estimation](#3-capacity-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Data Modeling](#5-data-modeling)
6. [API Design](#6-api-design)
7. [Deep Dives — Where Seniority Shows](#7-deep-dives--where-seniority-shows)
8. [Common Traps & How to Avoid Them](#8-common-traps--how-to-avoid-them)
9. [The Universal Building Blocks](#9-the-universal-building-blocks)
10. [Interview Questions & Model Answers](#10-interview-questions--model-answers)

---

## 1. The 45-Minute Framework

```
0-5  min:  Requirements clarification — ask, don't assume
5-10 min:  Capacity estimation — numbers drive decisions
10-20 min: High-level design — boxes and arrows, major components
20-30 min: Data modeling + API design
30-45 min: Deep dives — the interviewer's focus areas

CRITICAL: Let the interviewer guide your deep dives. Ask "which part would you like to explore deeper?"
```

---

## 2. Requirements Clarification — The Most Important Step

Never start designing without understanding what you're building. A URL shortener with 100 users has a completely different architecture than one with 100M users.

### Questions to Ask

```
Functional requirements:
  "What are the CORE features we need to support?"
  "What features can we explicitly exclude?"
  "Are there edge cases we should handle? (e.g., expired links, custom aliases?)"

Scale:
  "How many users do we expect? Daily active users?"
  "What's the read/write ratio? (important for database design)"
  "What's the expected QPS (queries per second)?"

Performance:
  "What's the acceptable latency for reads? Writes?"
  "Are there SLA requirements? (e.g., 99.9% uptime)"
  "Does this need to be globally distributed?"

Data:
  "How long do we need to store data? Any retention requirements?"
  "Do we need to support analytics/reporting?"

Constraints:
  "Are there regulatory requirements? (GDPR, PCI-DSS)"
  "What's the team size? (affects operational complexity choices)"
```

---

## 3. Capacity Estimation

The goal is not precise numbers — it's to determine if you need a single database or 10, a CDN or not, caching or not.

### Key Numbers to Memorize

```
Time:
  1 day = 86,400 seconds ≈ 100K seconds
  1 month = 2.6M seconds ≈ 30 × 86,400
  
Throughput:
  1M users, 10 requests/day = 10M requests/day
  10M requests/day ÷ 100K seconds = 100 requests/second
  
Storage:
  1 character = 1 byte
  1 tweet = 280 chars ≈ 300 bytes
  1 photo = 200KB average (compressed)
  1 video (1 min, 1080p) ≈ 50MB
  1 TB = 10^12 bytes = 10^6 MB
  
Latency:
  L1 cache: 1 ns
  RAM access: 100 ns
  SSD random read: 100 µs (0.1 ms)
  HDD seek: 10 ms
  Network (same datacenter): 0.5 ms
  Network (cross-continent): 100-150 ms
```

### Estimation Example: Twitter-like System

```
Given:
  300M daily active users
  Each user sends 2 tweets/day
  Read/write ratio: 100:1

Writes:
  300M users × 2 tweets = 600M tweets/day
  600M ÷ 86,400 = ~7,000 writes/second peak

Reads:
  7,000 writes/second × 100 (read/write ratio) = 700,000 reads/second
  
Storage (tweets only):
  600M tweets/day × 300 bytes/tweet = 180GB/day
  365 days = 65TB/year of tweet text
  
Media:
  Assume 10% of tweets have images: 60M images/day
  60M × 200KB = 12TB/day of images
  This needs a CDN + object storage, not a database!

Conclusion:
  Write QPS is manageable (7K/s) — single primary DB can handle this
  Read QPS (700K/s) — needs caching (Redis) + read replicas
  Media storage → CDN + S3 (not in the database)
```

---

## 4. High-Level Design

Draw the major components and their connections. Start with the happy path.

### Standard Components to Consider

```
Client layer:
  Web browsers, mobile apps, third-party API consumers

Edge layer:
  DNS → Load Balancer → API Gateway
  CDN for static assets and media

Application layer:
  Stateless application servers (horizontally scalable)
  Multiple services or a monolith (discuss trade-offs)

Data layer:
  Primary database (PostgreSQL/MySQL for transactional data)
  Cache (Redis)
  Object storage (S3) for media/blobs
  Search engine (Elasticsearch) if needed
  Message queue (Kafka) for async processing

Support services:
  Auth service / Identity provider
  Notification service (email, push, SMS)
  Analytics pipeline (Kafka → data warehouse)
```

### Drawing the Diagram

```
[Client] → [CDN] → [Load Balancer] → [API Gateway]
                                           │
                           ┌───────────────┼───────────────┐
                           ▼               ▼               ▼
                    [User Service]  [Tweet Service]  [Feed Service]
                           │               │               │
                    [PostgreSQL]   [PostgreSQL]     [Redis Cache]
                                           │
                                    [Object Store]
                                    (media/images)
                                           │
                                    [Kafka Queue]
                                           │
                                   [Fan-out Worker]
                                   (updates follower feeds)
```

---

## 5. Data Modeling

Define the main entities and their schemas. Think about:
- What queries will be most common?
- What indexes are needed?
- What can be denormalized for performance?

```sql
-- Example: Twitter-like system
-- Users table
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tweets table
CREATE TABLE tweets (
    id BIGINT PRIMARY KEY,
    author_id BIGINT REFERENCES users(id),
    content VARCHAR(280),
    media_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_tweets_author_time ON tweets(author_id, created_at DESC);

-- Follows table
CREATE TABLE follows (
    follower_id BIGINT REFERENCES users(id),
    following_id BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

-- For high-scale: consider pre-computing feeds
-- feed_items table: (user_id, tweet_id, created_at) — denormalized for fast reads
CREATE TABLE feed_items (
    user_id BIGINT,
    tweet_id BIGINT,
    tweet_author_id BIGINT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id, tweet_id)
);
CREATE INDEX idx_feed_user_time ON feed_items(user_id, created_at DESC);
```

---

## 6. API Design

Define the main API endpoints. Be specific about request/response format.

```
POST /tweets
  Request:  { content: string, media?: string[] }
  Response: { id: string, content: string, author: User, created_at: string }
  
GET /feed?cursor=<timestamp>&limit=20
  Response: { tweets: Tweet[], next_cursor: string }

GET /users/{username}/tweets?page=1&limit=20
  Response: { tweets: Tweet[], total: int }

POST /follows/{user_id}
  Response: 204 No Content

DELETE /follows/{user_id}
  Response: 204 No Content

Note cursor-based pagination vs offset-based:
  Offset: LIMIT 20 OFFSET 40 → slow for large offsets (DB scans all previous rows)
  Cursor: WHERE created_at < :cursor ORDER BY created_at DESC LIMIT 20
         → fast regardless of position (uses index)
```

---

## 7. Deep Dives — Where Seniority Shows

This is where the interview distinguishes a senior engineer from a mid-level one. The interviewer picks an area to dig into.

### Common Deep Dive Areas

**Feed generation:**
```
Push model (fan-out on write):
  When Alice tweets, immediately write to all her followers' feeds
  Pros: reads are instant (pre-computed)
  Cons: if Alice has 100M followers, one tweet creates 100M writes!
  Solution: fan-out for normal users, pull for celebrities

Pull model (fan-out on read):
  When Bob opens the app, query all accounts Bob follows, merge and sort
  Pros: no write amplification
  Cons: reads are slow (must query N accounts per feed load)

Hybrid (used by Twitter):
  Fan-out on write for users with < 10K followers
  Pull on read for users with > 10K followers (celebrities)
  Merge at read time
```

**Database sharding:**
- Shard tweets by tweet_id or user_id?
- If by user_id: all of a user's tweets on one shard (good for user profile page)
- If by tweet_id: even distribution (good for global timeline)
- Celebrity problem: shard for Justin Bieber's tweets gets all traffic

**Caching:**
- What to cache: hot tweets, user profiles, follower lists
- Cache invalidation: when does the cache expire?
- Write-through vs cache-aside vs read-through

---

## 8. Common Traps & How to Avoid Them

```
Trap 1: Jumping to solution before understanding requirements
  Fix: spend 5 minutes on requirements first, even if it feels slow

Trap 2: Perfect from the start
  Fix: draw the simple version first, then iterate
  "Let me start with the simplest design that works, and then we can optimize"

Trap 3: Ignoring failure modes
  Fix: after each component, ask "what happens if this fails?"
  How does data re-replicate? How do you handle a shard going down?

Trap 4: Not discussing trade-offs
  Fix: for every decision, mention the alternative and why you chose this one
  "I'm using SQL here because we need transactions; NoSQL would scale better for pure reads"

Trap 5: Inventing jargon or wrong numbers
  Fix: if you don't know, say "I'd need to verify this but roughly..." and move on
```

---

## 9. The Universal Building Blocks

Every system design uses a subset of these. Know each one cold:

| Building Block | When to Use |
|---|---|
| **Load balancer** | Scale out stateless servers, health checks |
| **CDN** | Static assets, geographic distribution, video streaming |
| **API Gateway** | Auth, rate limiting, routing, protocol translation |
| **Object storage (S3)** | Images, video, documents, backups |
| **Cache (Redis)** | Frequently-read data, sessions, rate limit counters |
| **Message queue (Kafka)** | Async processing, decoupling, event streaming |
| **Search (Elasticsearch)** | Full-text search, faceted search |
| **Database (PostgreSQL)** | Transactional data, complex queries |
| **NoSQL (DynamoDB/Cassandra)** | High-scale, simple access patterns |
| **Data warehouse (BigQuery)** | Analytics, reporting, OLAP queries |

---

## 10. Interview Questions & Model Answers

**Q: How do you approach a system design interview?**

"I start by clarifying requirements — functional (what features?) and non-functional (scale, latency, consistency). Then I estimate capacity to understand if we're dealing with 100 QPS or 1M QPS, since that determines whether I need caching, sharding, etc. Then I draw a high-level design — major components and their connections. Then data model and APIs. Finally, I deep dive into the most interesting or challenging parts — the interviewer usually has a focus area they want to explore. Throughout, I'm explicit about trade-offs: why SQL vs NoSQL, push vs pull for feeds, synchronous vs async. The goal isn't the 'right' answer but demonstrating how you think through complex trade-offs."

---

## Summary

- **Framework:** Clarify → Estimate → High-level → Data model → API → Deep dives.
- Spend the first 5 minutes on requirements. Never design for unspecified scale.
- Know capacity estimation numbers: 100K seconds/day, tweet = 300 bytes, photo = 200KB.
- High-level design: client → CDN → load balancer → services → databases.
- Deep dives separate seniors: feed fan-out, sharding strategies, cache invalidation, failure modes.
- Every decision needs a trade-off discussion. "I chose X because... The alternative Y would be better if..."
- Push vs pull for feed generation. Cursor vs offset pagination. Cache-aside vs write-through.
