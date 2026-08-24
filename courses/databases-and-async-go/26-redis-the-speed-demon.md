# Chapter 26: Redis — The Speed Demon

Redis is a database that lives in RAM. That's its secret: no disk seeks, no B-tree lookups, just direct memory access. It can handle millions of operations per second on a single server. This chapter covers everything from Redis basics to production patterns.

## Table of Contents

1. What Redis Is and Why It's Fast
2. Redis Data Types
3. Expiration and Eviction
4. Redis Persistence
5. Redis Clustering and High Availability
6. Common Redis Patterns
7. Building with Redis in Go
8. Exercises

---

## 1. What Redis Is and Why It's Fast

Redis (Remote Dictionary Server) is an in-memory key-value store. Every piece of data lives in RAM. Operations are O(1) or O(log n) — no disk I/O in the hot path.

**Why RAM is fast:**

| Operation | Approximate Time |
|-----------|-----------------|
| RAM access | ~100 nanoseconds |
| SSD read | ~100,000 nanoseconds (1,000x slower) |
| HDD read | ~10,000,000 nanoseconds (100,000x slower) |

Redis also uses a single-threaded event loop (like Node.js). No locking, no thread contention. All operations are atomic by default.

**Throughput on a single Redis server:** 100,000 to 1,000,000 operations/second depending on operation type and data size.

---

## 2. Redis Data Types

### Strings — The Basic Building Block

```
SET name "Alice"
GET name → "Alice"
DEL name

SET counter 0
INCR counter → 1   (atomic increment)
INCRBY counter 5 → 6

# With expiration (TTL)
SET session:abc123 "user_data" EX 3600   # expires in 1 hour
TTL session:abc123 → 3598               # seconds remaining
```

### Lists — Ordered Sequences

```
RPUSH queue job1 job2 job3  # push to right
LPUSH queue job0            # push to left
LRANGE queue 0 -1           # get all: [job0, job1, job2, job3]
LPOP queue                  # pop from left: job0
BRPOP queue 30              # blocking pop, wait up to 30 seconds
LLEN queue → 3
```

Use case: task queues (push jobs, workers pop them).

### Hashes — Dictionaries

```
HSET user:1 name "Alice" email "alice@example.com" age 30
HGET user:1 name → "Alice"
HGETALL user:1 → {name: Alice, email: alice@example.com, age: 30}
HINCRBY user:1 age 1  # atomic increment of one field
HDEL user:1 age
```

Use case: storing user sessions, object attributes.

### Sets — Unique Collections

```
SADD online_users user1 user2 user3
SISMEMBER online_users user1 → 1 (true)
SMEMBERS online_users → {user1, user2, user3}
SCARD online_users → 3

SADD user1_friends user2 user3 user4
SADD user2_friends user1 user3 user5
SINTER user1_friends user2_friends → {user3}  # common friends
SUNION user1_friends user2_friends → all friends
```

Use case: unique visitors, mutual friends, set operations.

### Sorted Sets — Sets with Scores

```
ZADD leaderboard 1500 "alice"
ZADD leaderboard 2300 "bob"
ZADD leaderboard 1800 "carol"

ZRANGE leaderboard 0 -1 WITHSCORES   # sorted by score ascending
ZREVRANGE leaderboard 0 2 WITHSCORES # top 3 (descending)
ZRANK leaderboard "alice" → 0        # rank (0-indexed)
ZSCORE leaderboard "bob" → 2300
ZINCRBY leaderboard 100 "alice"      # atomic: alice's score → 1600
```

Use case: leaderboards, priority queues, time-ordered events.

### Streams — Append-Only Log

```
XADD events * type click user_id 42   # * = auto-generate ID
XADD events * type purchase amount 99

XRANGE events - + COUNT 10            # read oldest 10 events
XREAD COUNT 10 BLOCK 0 STREAMS events $ # read new events as they arrive
```

Streams are Redis's Kafka-like feature. We cover them in Chapter 28.

---

## 3. Expiration and Eviction

### Key Expiration (TTL)

```
SET session:xyz "data" EX 3600    # expires after 3600 seconds
SET temp_key "val" PX 5000        # expires after 5000 milliseconds

EXPIRE session:xyz 7200           # reset TTL to 7200 seconds
PERSIST session:xyz               # remove TTL (make permanent)

TTL session:xyz → 3598            # seconds until expiry (-1 = no expiry, -2 = expired)
```

Redis checks keys lazily (when accessed) and actively (background sampling). You don't need to clean up expired keys manually.

### Eviction When Memory Is Full

When Redis runs out of memory, it must evict keys. Configure the policy:

```
# In redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru    # evict least-recently-used keys
```

Eviction policies:
- `allkeys-lru` — evict LRU from all keys (best for caching)
- `volatile-lru` — evict LRU only from keys with expiry
- `allkeys-random` — evict random keys
- `volatile-ttl` — evict keys with shortest remaining TTL
- `noeviction` — return error when full (best for queues)

---

## 4. Redis Persistence

Redis is in-memory, but you can persist to disk for recovery.

### RDB Snapshots

Periodically saves the entire dataset to disk as a binary snapshot (`dump.rdb`).

```
# redis.conf
save 60 1000    # save every 60 seconds if at least 1000 keys changed
save 300 10     # save every 5 minutes if at least 10 keys changed
```

**Pros:** Small file, fast startup. **Cons:** Potential data loss between snapshots (up to minutes).

### AOF (Append-Only File)

Logs every write operation. On restart, replays all operations.

```
# redis.conf
appendonly yes
appendfsync everysec   # flush to disk every second (balance: safety vs speed)
# appendfsync always   # safest: flush every write (slowest)
# appendfsync no       # fastest: let OS flush (risk of data loss)
```

**Pros:** Near-zero data loss. **Cons:** Larger files, slower startup.

### Recommended: RDB + AOF

Use both for production:
```
save 900 1
appendonly yes
appendfsync everysec
```

---

## 5. Redis Clustering and High Availability

### Redis Sentinel

Sentinel provides automatic failover for a single Redis primary.

```
Primary ← Sentinel 1
         ← Sentinel 2
         ← Sentinel 3
Replica 1
Replica 2
```

If the primary fails, sentinels vote to promote a replica. Client connects through Sentinel to always find the current primary.

### Redis Cluster

Cluster shards data across multiple nodes. Data is split into 16,384 hash slots — each node owns a range.

```
Node 1: slots 0-5460       (keys that hash to these slots)
Node 2: slots 5461-10922
Node 3: slots 10923-16383
Each node also has replicas for failover
```

The Go client handles routing transparently.

---

## 6. Common Redis Patterns

### Cache-Aside Pattern

```
function getUser(id):
  user = redis.GET("user:" + id)
  if user != nil:
    return deserialize(user)   # cache hit
  user = postgres.query("SELECT * FROM users WHERE id = ?", id)
  redis.SET("user:" + id, serialize(user), EX 300)  # cache for 5 min
  return user
```

### Distributed Lock (Redlock)

```
SET lock:resource "owner_id" NX PX 5000
# NX = only set if key doesn't exist
# PX 5000 = expire after 5 seconds

# Returns OK if lock acquired, nil if already locked
```

### Rate Limiter

```
key = "ratelimit:" + user_id + ":" + floor(now / 60)
count = INCR key
if count == 1:
    EXPIRE key 60       # set TTL on first call only
if count > 100:
    return error("rate limit exceeded")
```

### Pub/Sub

```
# Publisher
PUBLISH channel:news "breaking news message"

# Subscriber (blocks, waits for messages)
SUBSCRIBE channel:news
```

---

## 7. Building with Redis in Go

```bash
go get github.com/redis/go-redis/v9
```

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func main() {
    rdb = redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        PoolSize:     20,
        MinIdleConns: 5,
    })

    ctx := context.Background()
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal("cannot connect to Redis:", err)
    }
    fmt.Println("Connected to Redis!")

    // String operations
    rdb.Set(ctx, "greeting", "Hello, Redis!", 5*time.Minute)
    val, _ := rdb.Get(ctx, "greeting").Result()
    fmt.Println("greeting:", val)

    // Struct caching
    type User struct {
        ID    int    `json:"id"`
        Email string `json:"email"`
    }

    // Cache a user
    user := User{ID: 1, Email: "alice@example.com"}
    data, _ := json.Marshal(user)
    rdb.Set(ctx, "user:1", data, 5*time.Minute)

    // Retrieve cached user
    data, err := rdb.Get(ctx, "user:1").Bytes()
    if err == redis.Nil {
        fmt.Println("cache miss")
    } else if err != nil {
        log.Fatal(err)
    } else {
        var cached User
        json.Unmarshal(data, &cached)
        fmt.Printf("cached user: %+v\n", cached)
    }

    // Sorted set leaderboard
    rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 1500, Member: "alice"})
    rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 2300, Member: "bob"})
    rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 1800, Member: "carol"})

    top3, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
    fmt.Println("Top 3:")
    for i, z := range top3 {
        fmt.Printf("  %d. %s (%.0f)\n", i+1, z.Member, z.Score)
    }

    // Rate limiting
    key := fmt.Sprintf("ratelimit:user1:%d", time.Now().Unix()/60)
    count, _ := rdb.Incr(ctx, key).Result()
    if count == 1 {
        rdb.Expire(ctx, key, time.Minute)
    }
    fmt.Printf("Requests this minute: %d\n", count)
}
```

### Cache-Aside in Go

```go
func GetUser(ctx context.Context, db DB, id int64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // Try cache first
    data, err := rdb.Get(ctx, cacheKey).Bytes()
    if err == nil {
        var u User
        json.Unmarshal(data, &u)
        return &u, nil // cache hit!
    }
    if err != redis.Nil {
        // Redis error (not a miss) — log but continue to DB
        log.Printf("redis get %s: %v", cacheKey, err)
    }

    // Cache miss: go to database
    u, err := db.GetUser(ctx, id)
    if err != nil {
        return nil, err
    }
    if u == nil {
        return nil, nil // user doesn't exist
    }

    // Store in cache
    if data, err := json.Marshal(u); err == nil {
        rdb.Set(ctx, cacheKey, data, 5*time.Minute)
    }

    return u, nil
}
```

---

## Summary

- Redis stores data in RAM → microsecond latency, millions of ops/second.
- Data types: string, list, hash, set, sorted set, stream — each optimized for different use cases.
- TTL-based expiration and eviction policies handle cache management automatically.
- RDB + AOF together provide the best balance of durability and performance.
- Sentinel provides automatic failover; Cluster provides horizontal sharding.
- The cache-aside pattern is the most common way to use Redis with a relational database.

### Exercises

**Easy:** Build a simple session store using Redis hashes. `SET session:{token} user_id=1 username=alice` with 30-minute TTL. Read and delete sessions.

**Medium:** Implement a rate limiter in Go that allows 100 requests per minute per IP. Use Redis INCR + EXPIRE.

**Hard:** Implement a leaderboard API: `POST /score` to add/update a user's score, `GET /leaderboard?limit=10` for the top N, `GET /rank?user=alice` for a specific user's rank. Use Redis sorted sets.
