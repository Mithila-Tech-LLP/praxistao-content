# Chapter 27: Building with Redis in Go

Complete production patterns for Redis in Go: connection management, all data types, pipelines, Lua scripts, and a real-time chat system.

## Table of Contents

1. Production Connection Setup
2. All Data Types in Go Code
3. Pipelining — Batching Redis Commands
4. Lua Scripts — Atomic Operations
5. Pub/Sub Real-Time Messaging
6. Mini Project: Real-Time Chat with Redis
7. Exercises

---

## 1. Production Connection Setup

```go
package cache

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

func NewClient(addr, password string, db int) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:            addr,
        Password:        password,
        DB:              db,
        PoolSize:        20,
        MinIdleConns:    5,
        MaxRetries:      3,
        DialTimeout:     5 * time.Second,
        ReadTimeout:     3 * time.Second,
        WriteTimeout:    3 * time.Second,
        PoolTimeout:     4 * time.Second,
        ConnMaxIdleTime: 30 * time.Minute,
        ConnMaxLifetime: time.Hour,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("ping redis: %w", err)
    }
    return client, nil
}

// NewClusterClient for Redis Cluster
func NewClusterClient(addrs []string) (*redis.ClusterClient, error) {
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:    addrs,
        PoolSize: 10,
    })
    ctx := context.Background()
    return client, client.Ping(ctx).Err()
}
```

---

## 2. All Data Types in Go Code

### Strings — Cache and Counters

```go
type Cache struct{ rdb *redis.Client }

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) GetString(ctx context.Context, key string) (string, bool, error) {
    val, err := c.rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", false, nil // cache miss
    }
    return val, true, err
}

func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
    return c.rdb.Incr(ctx, key).Result()
}

func (c *Cache) IncrementBy(ctx context.Context, key string, by int64) (int64, error) {
    return c.rdb.IncrBy(ctx, key, by).Result()
}
```

### Hashes — Structured Objects

```go
type Session struct {
    UserID   string
    Username string
    Role     string
}

func (c *Cache) SetSession(ctx context.Context, token string, s Session, ttl time.Duration) error {
    key := "session:" + token
    err := c.rdb.HSet(ctx, key,
        "user_id", s.UserID,
        "username", s.Username,
        "role", s.Role,
    ).Err()
    if err != nil {
        return err
    }
    return c.rdb.Expire(ctx, key, ttl).Err()
}

func (c *Cache) GetSession(ctx context.Context, token string) (*Session, error) {
    key := "session:" + token
    fields, err := c.rdb.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    if len(fields) == 0 {
        return nil, nil // not found
    }
    return &Session{
        UserID:   fields["user_id"],
        Username: fields["username"],
        Role:     fields["role"],
    }, nil
}
```

### Lists — Task Queues

```go
type Queue struct{ rdb *redis.Client }

func (q *Queue) Push(ctx context.Context, name string, payload string) error {
    return q.rdb.RPush(ctx, "queue:"+name, payload).Err()
}

// Pop blocks for up to `timeout` waiting for a job
func (q *Queue) Pop(ctx context.Context, name string, timeout time.Duration) (string, error) {
    result, err := q.rdb.BLPop(ctx, timeout, "queue:"+name).Result()
    if err == redis.Nil {
        return "", nil // timeout, no job
    }
    if err != nil {
        return "", err
    }
    return result[1], nil // result[0] is the queue name, result[1] is the value
}

func (q *Queue) Length(ctx context.Context, name string) (int64, error) {
    return q.rdb.LLen(ctx, "queue:"+name).Result()
}
```

### Sorted Sets — Leaderboards

```go
type Leaderboard struct{ rdb *redis.Client }

func (l *Leaderboard) AddScore(ctx context.Context, name, member string, score float64) error {
    return l.rdb.ZAdd(ctx, "lb:"+name, redis.Z{Score: score, Member: member}).Err()
}

func (l *Leaderboard) IncrScore(ctx context.Context, name, member string, by float64) (float64, error) {
    return l.rdb.ZIncrBy(ctx, "lb:"+name, by, member).Result()
}

func (l *Leaderboard) TopN(ctx context.Context, name string, n int64) ([]redis.Z, error) {
    return l.rdb.ZRevRangeWithScores(ctx, "lb:"+name, 0, n-1).Result()
}

func (l *Leaderboard) GetRank(ctx context.Context, name, member string) (int64, float64, error) {
    rank, err := l.rdb.ZRevRank(ctx, "lb:"+name, member).Result()
    if err == redis.Nil {
        return -1, 0, nil
    }
    score, _ := l.rdb.ZScore(ctx, "lb:"+name, member).Result()
    return rank + 1, score, err // +1 for 1-indexed rank
}
```

---

## 3. Pipelining — Batching Redis Commands

Each Redis command is a network round-trip. Pipelining sends multiple commands in one go:

```go
// Without pipelining: 3 network round-trips
rdb.Set(ctx, "key1", "val1", 0)
rdb.Set(ctx, "key2", "val2", 0)
rdb.Set(ctx, "key3", "val3", 0)

// With pipelining: 1 network round-trip
pipe := rdb.Pipeline()
pipe.Set(ctx, "key1", "val1", 0)
pipe.Set(ctx, "key2", "val2", 0)
pipe.Set(ctx, "key3", "val3", 0)
_, err := pipe.Exec(ctx)
```

Pipelining with reading results:

```go
pipe := rdb.Pipeline()

// Queue up commands
get1 := pipe.Get(ctx, "user:1")
get2 := pipe.Get(ctx, "user:2")
get3 := pipe.Get(ctx, "user:3")

// Execute all at once
_, err := pipe.Exec(ctx)
if err != nil && err != redis.Nil {
    return err
}

// Read results after Exec
val1, _ := get1.Result()
val2, _ := get2.Result()
val3, _ := get3.Result()
```

---

## 4. Lua Scripts — Atomic Operations

Lua scripts run atomically in Redis — no other command executes between lines. Essential for operations that must be atomic but require multiple commands.

```go
// Atomic check-and-set with expiry (can't do this with SETNX + EXPIRE atomically)
var lockScript = redis.NewScript(`
    local current = redis.call('GET', KEYS[1])
    if current == false then
        redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[2])
        return 1
    end
    return 0
`)

func acquireLock(ctx context.Context, rdb *redis.Client, key, owner string, ttl time.Duration) (bool, error) {
    result, err := lockScript.Run(ctx, rdb,
        []string{key},             // KEYS
        owner,                     // ARGV[1]
        ttl.Milliseconds(),        // ARGV[2]
    ).Int()
    return result == 1, err
}

// Atomic release lock (only if we own it)
var releaseLockScript = redis.NewScript(`
    if redis.call('GET', KEYS[1]) == ARGV[1] then
        return redis.call('DEL', KEYS[1])
    end
    return 0
`)

func releaseLock(ctx context.Context, rdb *redis.Client, key, owner string) error {
    return releaseLockScript.Run(ctx, rdb, []string{key}, owner).Err()
}
```

Usage:

```go
acquired, err := acquireLock(ctx, rdb, "lock:payment:user123", "my-process-id", 30*time.Second)
if !acquired {
    return errors.New("another process is handling this payment")
}
defer releaseLock(ctx, rdb, "lock:payment:user123", "my-process-id")

// Critical section — only one process runs this at a time
processPayment(...)
```

---

## 5. Pub/Sub Real-Time Messaging

```go
// Publisher: send messages to a channel
func PublishMessage(ctx context.Context, rdb *redis.Client, channel, message string) error {
    return rdb.Publish(ctx, channel, message).Err()
}

// Subscriber: receive messages
func Subscribe(ctx context.Context, rdb *redis.Client, channel string, handler func(string)) {
    pubsub := rdb.Subscribe(ctx, channel)
    defer pubsub.Close()

    ch := pubsub.Channel()
    for {
        select {
        case msg := <-ch:
            handler(msg.Payload)
        case <-ctx.Done():
            return
        }
    }
}

// Pattern subscribe: receive messages from all matching channels
func SubscribePattern(ctx context.Context, rdb *redis.Client, pattern string, handler func(channel, payload string)) {
    pubsub := rdb.PSubscribe(ctx, pattern)
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        handler(msg.Channel, msg.Payload)
    }
}
```

---

## 6. Mini Project: Real-Time Chat with Redis

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type Message struct {
    Room    string    `json:"room"`
    User    string    `json:"user"`
    Text    string    `json:"text"`
    SentAt  time.Time `json:"sent_at"`
}

type SSEClient struct {
    ch chan string
}

var (
    mu      sync.RWMutex
    clients = map[string][]SSEClient{} // room → clients
)

func main() {
    rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    ctx := context.Background()

    // Subscribe to all room channels and forward to SSE clients
    pubsub := rdb.PSubscribe(ctx, "chat:*")
    go func() {
        for msg := range pubsub.Channel() {
            mu.RLock()
            roomClients := clients[msg.Channel]
            mu.RUnlock()
            for _, c := range roomClients {
                select {
                case c.ch <- msg.Payload:
                default: // client too slow, skip
                }
            }
        }
    }()

    http.HandleFunc("POST /chat/{room}", handleSend)
    http.HandleFunc("GET /chat/{room}/stream", handleStream)
    http.HandleFunc("GET /chat/{room}/history", handleHistory)

    log.Println("Chat server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
    room := r.PathValue("room")
    user := r.URL.Query().Get("user")
    if user == "" {
        http.Error(w, "user required", 400)
        return
    }

    var req struct{ Text string `json:"text"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }

    msg := Message{
        Room:   room,
        User:   user,
        Text:   req.Text,
        SentAt: time.Now(),
    }
    data, _ := json.Marshal(msg)

    ctx := r.Context()
    // Store in history (last 100 messages)
    rdb.LPush(ctx, "history:"+room, data)
    rdb.LTrim(ctx, "history:"+room, 0, 99)

    // Publish to subscribers
    rdb.Publish(ctx, "chat:"+room, string(data))

    w.WriteHeader(204)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
    room := r.PathValue("room")

    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    client := SSEClient{ch: make(chan string, 10)}
    mu.Lock()
    clients["chat:"+room] = append(clients["chat:"+room], client)
    mu.Unlock()

    defer func() {
        mu.Lock()
        roomClients := clients["chat:"+room]
        for i, c := range roomClients {
            if c.ch == client.ch {
                clients["chat:"+room] = append(roomClients[:i], roomClients[i+1:]...)
                break
            }
        }
        mu.Unlock()
    }()

    flusher, _ := w.(http.Flusher)
    for {
        select {
        case msg := <-client.ch:
            fmt.Fprintf(w, "data: %s\n\n", msg)
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
    room := r.PathValue("room")
    msgs, _ := rdb.LRange(r.Context(), "history:"+room, 0, 49).Result()

    var messages []Message
    for i := len(msgs) - 1; i >= 0; i-- { // reverse (newest first in Redis, oldest first for history)
        var msg Message
        json.Unmarshal([]byte(msgs[i]), &msg)
        messages = append(messages, msg)
    }
    json.NewEncoder(w).Encode(messages)
}
```

Test:
```bash
# Send a message
curl -X POST "localhost:8080/chat/general?user=Alice" \
  -d '{"text":"Hello everyone!"}'

# Stream messages (SSE)
curl -N "localhost:8080/chat/general/stream"

# Get history
curl "localhost:8080/chat/general/history"
```

---

## Summary

- `BLPop` creates blocking queues — workers wait for jobs without polling.
- Pipelining reduces round-trips by batching commands — critical for bulk operations.
- Lua scripts are atomic — use them for operations that require checking then setting.
- Pub/Sub is ephemeral — no message history. Use `LPUSH/LRANGE` + Pub/Sub together for real-time + history.
- Distributed locks with Lua scripts prevent race conditions across multiple processes.

### Exercises

**Easy:** Build a Go task queue: producer pushes job payloads with `RPUSH`, consumer polls with `BLPOP`. Add worker count as a CLI flag to run multiple workers concurrently.

**Medium:** Build a "trending topics" feature: every time a hashtag is mentioned, increment its score. Return the top 10 trending hashtags. Reset scores hourly using a time-based key.

**Hard:** Implement a distributed rate limiter using sliding window algorithm: for each user, maintain a sorted set of request timestamps. On each request, remove timestamps older than 1 minute, count remaining, and add current timestamp. Reject if count > 100.
