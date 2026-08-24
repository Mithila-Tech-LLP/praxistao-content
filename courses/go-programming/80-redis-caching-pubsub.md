# Chapter 80: Redis — Caching, Pub/Sub, and Streams

Redis is an in-memory data structure server. "Caching layer" is its most common use case, but it also has native pub/sub, sorted sets (leaderboards, rate limiting), and Streams (durable event log). This chapter covers all four with the `go-redis/v9` client.

## Table of Contents

1. [Connecting and Basic Operations](#1-connecting-and-basic-operations)
2. [Caching Patterns](#2-caching-patterns)
3. [Distributed Locks](#3-distributed-locks)
4. [Sorted Sets — Leaderboards and Rate Limiting](#4-sorted-sets--leaderboards-and-rate-limiting)
5. [Pub/Sub](#5-pubsub)
6. [Redis Streams](#6-redis-streams)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Connecting and Basic Operations

```go
import "github.com/redis/go-redis/v9"

func newRedisClient(addr, password string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         addr,      // "localhost:6379"
        Password:     password,  // "" for no auth
        DB:           db,        // 0-15
        PoolSize:     20,
        MinIdleConns: 5,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
}

// Health check
func ping(ctx context.Context, rdb *redis.Client) error {
    return rdb.Ping(ctx).Err()
}
```

### Core data types

```go
ctx := context.Background()

// Strings
rdb.Set(ctx, "name", "alice", 5*time.Minute)
val, _ := rdb.Get(ctx, "name").Result()  // "alice"
rdb.Incr(ctx, "counter")                // atomic increment

// Hashes (like a dictionary stored at a single key)
rdb.HSet(ctx, "user:1", "name", "Alice", "email", "alice@example.com")
name, _ := rdb.HGet(ctx, "user:1", "name").Result()
fields, _ := rdb.HGetAll(ctx, "user:1").Result() // map[string]string

// Lists
rdb.LPush(ctx, "queue", "item1", "item2")
item, _ := rdb.RPop(ctx, "queue").Result()  // FIFO queue
rdb.BRPop(ctx, 0, "queue")                 // blocking pop (wait for item)

// Sets
rdb.SAdd(ctx, "online_users", "user:1", "user:2")
rdb.SIsMember(ctx, "online_users", "user:1")  // bool
rdb.SMembers(ctx, "online_users")             // []string
rdb.SCard(ctx, "online_users")               // count

// Sorted sets (score + member)
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 1500, Member: "alice"})
rdb.ZRangeWithScores(ctx, "leaderboard", 0, 9) // top 10
```

---

## 2. Caching Patterns

### Cache-aside (lazy loading)

```go
type UserCache struct {
    rdb *redis.Client
    ttl time.Duration
    db  UserDB
}

func (c *UserCache) Get(ctx context.Context, id int64) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // Try cache first
    data, err := c.rdb.Get(ctx, key).Bytes()
    if err == nil {
        var u User
        if json.Unmarshal(data, &u) == nil { return &u, nil }
    }
    
    // Cache miss — load from database
    u, err := c.db.GetUser(ctx, id)
    if err != nil { return nil, err }
    
    // Populate cache
    if data, err := json.Marshal(u); err == nil {
        c.rdb.Set(ctx, key, data, c.ttl)
    }
    return u, nil
}

func (c *UserCache) Invalidate(ctx context.Context, id int64) {
    c.rdb.Del(ctx, fmt.Sprintf("user:%d", id))
}
```

### Cache stampede prevention with SETNX

When many requests miss the cache simultaneously and all hit the database at once:

```go
func (c *UserCache) GetWithLock(ctx context.Context, id int64) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    lockKey := key + ":lock"
    
    for {
        // Try cache
        data, err := c.rdb.Get(ctx, key).Bytes()
        if err == nil {
            var u User
            if json.Unmarshal(data, &u) == nil { return &u, nil }
        }
        
        // Try to acquire the fill lock (1 second TTL)
        acquired, err := c.rdb.SetNX(ctx, lockKey, "1", time.Second).Result()
        if err != nil { return nil, err }
        
        if acquired {
            // We got the lock — load from DB and populate cache
            defer c.rdb.Del(ctx, lockKey)
            
            u, err := c.db.GetUser(ctx, id)
            if err != nil { return nil, err }
            
            if data, err := json.Marshal(u); err == nil {
                c.rdb.Set(ctx, key, data, c.ttl)
            }
            return u, nil
        }
        
        // Another goroutine is filling the cache — wait briefly
        time.Sleep(5 * time.Millisecond)
    }
}
```

### Write-through cache

```go
func (c *UserCache) Update(ctx context.Context, u *User) error {
    // Write to DB first
    if err := c.db.UpdateUser(ctx, u); err != nil { return err }
    
    // Then update cache (write-through)
    if data, err := json.Marshal(u); err == nil {
        c.rdb.Set(ctx, fmt.Sprintf("user:%d", u.ID), data, c.ttl)
    }
    return nil
}
```

### Pipelining — batch multiple commands

```go
// Single round trip for multiple operations
pipe := rdb.Pipeline()
cmds := make([]*redis.StringCmd, len(userIDs))
for i, id := range userIDs {
    cmds[i] = pipe.Get(ctx, fmt.Sprintf("user:%d", id))
}
pipe.Exec(ctx)

users := make([]*User, 0, len(userIDs))
for i, cmd := range cmds {
    data, err := cmd.Bytes()
    if err == redis.Nil { continue } // cache miss
    if err != nil { continue }
    var u User
    if json.Unmarshal(data, &u) == nil {
        users = append(users, &u)
    }
    _ = i
}
```

---

## 3. Distributed Locks

For mutual exclusion across multiple service instances:

```go
type Lock struct {
    rdb   *redis.Client
    key   string
    token string // unique per lock acquisition
    ttl   time.Duration
}

func AcquireLock(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) (*Lock, error) {
    token := uuid.New().String()
    key = "lock:" + key
    
    ok, err := rdb.SetNX(ctx, key, token, ttl).Result()
    if err != nil { return nil, fmt.Errorf("acquire lock: %w", err) }
    if !ok       { return nil, fmt.Errorf("lock %q is held by another process", key) }
    
    return &Lock{rdb: rdb, key: key, token: token, ttl: ttl}, nil
}

func (l *Lock) Release(ctx context.Context) error {
    // Lua script: only delete if our token matches (prevents releasing someone else's lock)
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `)
    return script.Run(ctx, l.rdb, []string{l.key}, l.token).Err()
}

// Usage
func processOrder(ctx context.Context, rdb *redis.Client, orderID int64) error {
    lock, err := AcquireLock(ctx, rdb, fmt.Sprintf("order:%d", orderID), 30*time.Second)
    if err != nil { return fmt.Errorf("could not lock order: %w", err) }
    defer lock.Release(ctx)
    
    // only one instance processes this order at a time
    return doProcess(ctx, orderID)
}
```

---

## 4. Sorted Sets — Leaderboards and Rate Limiting

### Leaderboard

```go
type LeaderboardEntry struct {
    Rank   int
    Member string
    Score  float64
}

func (lb *Leaderboard) AddScore(ctx context.Context, userID string, delta float64) (float64, error) {
    return lb.rdb.ZIncrBy(ctx, lb.key, delta, userID).Result()
}

func (lb *Leaderboard) TopN(ctx context.Context, n int) ([]LeaderboardEntry, error) {
    // ZRevRangeWithScores returns highest scores first
    results, err := lb.rdb.ZRevRangeWithScores(ctx, lb.key, 0, int64(n-1)).Result()
    if err != nil { return nil, err }
    
    entries := make([]LeaderboardEntry, len(results))
    for i, z := range results {
        entries[i] = LeaderboardEntry{
            Rank:   i + 1,
            Member: z.Member.(string),
            Score:  z.Score,
        }
    }
    return entries, nil
}

func (lb *Leaderboard) Rank(ctx context.Context, userID string) (int, float64, error) {
    rank, err := lb.rdb.ZRevRank(ctx, lb.key, userID).Result()
    if err != nil { return 0, 0, err }
    score, _ := lb.rdb.ZScore(ctx, lb.key, userID).Result()
    return int(rank) + 1, score, nil
}
```

### Sliding window rate limiter

```go
type RateLimiter struct {
    rdb    *redis.Client
    limit  int
    window time.Duration
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, int, error) {
    now := time.Now().UnixMilli()
    windowStart := now - r.window.Milliseconds()
    fullKey := "rate:" + key
    
    pipe := r.rdb.Pipeline()
    // Remove entries older than window
    pipe.ZRemRangeByScore(ctx, fullKey, "-inf", strconv.FormatInt(windowStart, 10))
    // Count entries in current window
    countCmd := pipe.ZCard(ctx, fullKey)
    // Add current request
    pipe.ZAdd(ctx, fullKey, redis.Z{Score: float64(now), Member: strconv.FormatInt(now, 10)})
    // Set TTL
    pipe.Expire(ctx, fullKey, r.window+time.Second)
    
    if _, err := pipe.Exec(ctx); err != nil { return false, 0, err }
    
    count := int(countCmd.Val()) + 1
    return count <= r.limit, r.limit - count, nil
}

// Usage in HTTP middleware
func rateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            key := r.Header.Get("X-API-Key")
            allowed, remaining, err := limiter.Allow(r.Context(), key)
            if err != nil {
                // Fail open: log and allow
                next.ServeHTTP(w, r)
                return
            }
            w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
            if !allowed {
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## 5. Pub/Sub

Redis Pub/Sub delivers messages to all current subscribers. It does NOT persist messages — if no one is subscribed when a message is published, it's lost.

```go
// Publisher
func publishEvent(ctx context.Context, rdb *redis.Client, channel string, event any) error {
    data, err := json.Marshal(event)
    if err != nil { return err }
    return rdb.Publish(ctx, channel, data).Err()
}

// Subscriber (runs in its own goroutine)
func subscribeToEvents(ctx context.Context, rdb *redis.Client, channel string, handler func([]byte)) {
    sub := rdb.Subscribe(ctx, channel)
    defer sub.Close()
    
    ch := sub.Channel()
    for {
        select {
        case msg, ok := <-ch:
            if !ok { return }
            handler([]byte(msg.Payload))
        case <-ctx.Done():
            return
        }
    }
}

// Pattern subscription: subscribe to "events.*"
func subscribePattern(ctx context.Context, rdb *redis.Client) {
    psub := rdb.PSubscribe(ctx, "events.*")
    defer psub.Close()
    
    for msg := range psub.Channel() {
        fmt.Printf("channel=%s payload=%s\n", msg.Channel, msg.Payload)
    }
}
```

---

## 6. Redis Streams

Streams are an append-only log with consumer groups. Unlike Pub/Sub:
- Messages are persisted
- Consumer groups allow multiple workers to share a stream without duplicate processing
- Failed messages can be re-delivered

```go
// Producer: append to stream
type OrderEvent struct {
    OrderID  string `json:"order_id"`
    UserID   string `json:"user_id"`
    Amount   int64  `json:"amount"`
    Status   string `json:"status"`
}

func (p *EventProducer) Publish(ctx context.Context, event OrderEvent) error {
    _, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
        Stream: "orders",
        MaxLen: 100_000, // trim stream to ~100k messages
        Approx: true,    // MAXLEN ~ is faster (approximate trim)
        Values: map[string]any{
            "order_id": event.OrderID,
            "user_id":  event.UserID,
            "amount":   event.Amount,
            "status":   event.Status,
        },
    }).Err()
    return err
}

// Consumer with consumer groups
type OrderConsumer struct {
    rdb      *redis.Client
    stream   string
    group    string
    consumer string
}

func NewOrderConsumer(rdb *redis.Client, stream, group, consumer string) *OrderConsumer {
    c := &OrderConsumer{rdb: rdb, stream: stream, group: group, consumer: consumer}
    // Create group if it doesn't exist ($ = start from new messages)
    rdb.XGroupCreateMkStream(context.Background(), stream, group, "$")
    return c
}

func (c *OrderConsumer) Run(ctx context.Context, handler func(OrderEvent) error) {
    for {
        // Read up to 10 messages, block for 2 seconds if stream is empty
        msgs, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    c.group,
            Consumer: c.consumer,
            Streams:  []string{c.stream, ">"},
            Count:    10,
            Block:    2 * time.Second,
        }).Result()
        
        if err == redis.Nil { continue } // timeout, no messages
        if err != nil {
            if ctx.Err() != nil { return }
            time.Sleep(time.Second)
            continue
        }
        
        for _, stream := range msgs {
            for _, msg := range stream.Messages {
                event := OrderEvent{
                    OrderID: msg.Values["order_id"].(string),
                    UserID:  msg.Values["user_id"].(string),
                    Status:  msg.Values["status"].(string),
                }
                
                if err := handler(event); err != nil {
                    // Leave message in PEL (pending entries list) for retry
                    continue
                }
                
                // Acknowledge processed message
                c.rdb.XAck(ctx, c.stream, c.group, msg.ID)
            }
        }
    }
}

// Recover pending messages from a crashed consumer
func (c *OrderConsumer) RecoverPending(ctx context.Context, handler func(OrderEvent) error) error {
    // Claim messages pending for > 30 seconds (worker likely crashed)
    msgs, _, err := c.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
        Stream:   c.stream,
        Group:    c.group,
        Consumer: c.consumer,
        MinIdle:  30 * time.Second,
        Start:    "0",
        Count:    100,
    }).Result()
    if err != nil { return err }
    
    for _, msg := range msgs {
        event := OrderEvent{
            OrderID: msg.Values["order_id"].(string),
            Status:  msg.Values["status"].(string),
        }
        if err := handler(event); err != nil { continue }
        c.rdb.XAck(ctx, c.stream, c.group, msg.ID)
    }
    return nil
}
```

---

## Summary

| Feature | Use Case | Persistence |
|---------|----------|-------------|
| Strings | Cache, counters, feature flags | TTL-based |
| Hashes | User sessions, config | TTL-based |
| Sorted sets | Leaderboards, rate limiting | Until explicitly deleted |
| Pub/Sub | Real-time notifications | None — fire and forget |
| Streams | Durable event log, task queues | Until trimmed |
| SETNX | Distributed locks | TTL-based |

- Use **pipelining** to batch multiple commands into one round trip
- Use **Lua scripts** (`redis.NewScript`) for atomic multi-step operations
- Use **Streams** instead of Pub/Sub when you need message persistence and at-least-once delivery
- Always **set a TTL** on cached data — unbounded caches fill memory
- Handle `redis.Nil` errors explicitly — it means key not found, not a real error

## Exercises

### Easy
1. Implement a `SessionStore` using Redis hashes: `HSET session:<token> user_id 1 created_at 1234 expires_at 9999`. Support `Create`, `Get`, `Delete`, and automatic expiry using `EXPIRE`.
2. Build a `VisitCounter` using Redis INCR: count page views per URL per day. Key: `views:2026-06-01:/products/123`. Return the count for any URL+date combination.
3. Implement a `SetWithJSON` / `GetWithJSON` helper that marshals/unmarshals Go structs automatically when writing to and reading from Redis.

### Medium
4. Build a **token bucket rate limiter** using a Redis Lua script: each bucket has a capacity and a refill rate. The Lua script atomically checks the current token count, refills based on elapsed time, and either allows or denies the request.
5. Implement a **pub/sub chat system**: a `ChatRoom` type with `Publish(message)` and `Subscribe() <-chan Message`. Use Redis pub/sub so multiple instances of the server can broadcast to all connected clients.
6. Write a **stream-based task queue**: producers push jobs to a Redis Stream. A pool of workers (`XReadGroup`) processes jobs. If a worker crashes, `XAutoClaim` re-delivers stuck messages to healthy workers after a timeout.

### Hard
7. Implement a **distributed semaphore** using a Redis sorted set: allow at most N concurrent operations. Use `ZADD ... NX` to acquire, `ZREM` to release, and `ZREMRANGEBYSCORE` to clean up expired holders. All operations should be atomic via Lua.
8. Build a **cache stampede test**: create a benchmark that simulates 1000 concurrent goroutines requesting the same cache key simultaneously after it expires. Compare three strategies: naive (no lock), SETNX lock, and probabilistic early expiration (refresh the cache before it expires based on remaining TTL / time-to-recompute ratio).
