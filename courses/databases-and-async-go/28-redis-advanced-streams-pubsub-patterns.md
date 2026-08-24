# Chapter 28: Redis Advanced — Streams, Pub/Sub, and Patterns

Redis Streams are Redis's answer to Kafka: an append-only log with consumer groups, offset tracking, and at-least-once delivery. This chapter covers advanced Redis features that power production systems.

## Table of Contents

1. Redis Streams — Kafka-Lite
2. Consumer Groups
3. HyperLogLog — Counting Unique Visitors
4. Distributed Locks Deep Dive
5. Redis as a Job Queue
6. Mini Project: Activity Feed with Redis Streams
7. Exercises

---

## 1. Redis Streams — Kafka-Lite

A Redis Stream is an append-only log of messages. Unlike Pub/Sub (ephemeral, fire-and-forget), streams persist messages and allow consumers to read from any point.

```
Stream: "events"
──────────────────────────────────────────
1704067200000-0  type=login  user=alice
1704067201000-0  type=view   page=/home  user=alice
1704067202000-0  type=click  button=buy  user=bob
1704067203000-0  type=login  user=carol
──────────────────────────────────────────
```

Each entry has an auto-generated ID (`timestamp-sequence`) and arbitrary key-value fields.

```
# Append to a stream (* = auto-generate ID)
XADD events * type login user alice
XADD events * type view page /home user alice
XADD events * type purchase amount 99 user alice

# Read all entries
XRANGE events - + COUNT 10
# - = oldest, + = newest

# Read new entries (like "tail -f")
XREAD COUNT 10 BLOCK 1000 STREAMS events $
# $ = only new entries, block 1000ms if empty

# Read from a specific ID onwards
XREAD COUNT 10 STREAMS events 1704067201000-0

# Stream length
XLEN events → 3

# Trim to last 1000 entries (cap the log)
XTRIM events MAXLEN ~ 1000  # ~ = approximate trim (faster)
```

---

## 2. Consumer Groups

Consumer groups allow multiple workers to process a stream without processing the same message twice.

```
Stream: "orders"
           ┌─────────────────────────────────────────┐
Consumer Group: "processors"
  Worker 1: reads messages 1, 4, 7, 10...
  Worker 2: reads messages 2, 5, 8, 11...
  Worker 3: reads messages 3, 6, 9, 12...
           └─────────────────────────────────────────┘
```

Each worker reads different messages. Unacknowledged messages go to the "pending" list for redelivery.

```
# Create consumer group (start from beginning with 0, or latest with $)
XGROUP CREATE orders processors $ MKSTREAM

# Worker 1 reads messages assigned to it
XREADGROUP GROUP processors worker-1 COUNT 5 BLOCK 2000 STREAMS orders >
# > = read next undelivered message

# Acknowledge successful processing
XACK orders processors 1704067200000-0

# See pending (unacknowledged) messages
XPENDING orders processors - + 10

# Re-claim messages pending > 30 seconds (for dead workers)
XAUTOCLAIM orders processors worker-2 30000 0-0
```

---

## 3. HyperLogLog — Counting Unique Visitors

How many unique visitors does your site have today? Storing every visitor's ID is expensive. HyperLogLog is a probabilistic data structure that estimates unique counts using only 12 KB of memory — with < 1% error.

```
PFADD visitors:2024-01-15 user123 user456 user789 user123  # user123 counted once
PFCOUNT visitors:2024-01-15 → 3 (approximate)

# Add more visitors
PFADD visitors:2024-01-15 user101 user202 user303

# Count across multiple days (unique visitors this week)
PFMERGE visitors:week visitors:2024-01-15 visitors:2024-01-16 visitors:2024-01-17
PFCOUNT visitors:week → ~total unique
```

In Go:

```go
func trackVisitor(ctx context.Context, rdb *redis.Client, date, userID string) error {
    key := "visitors:" + date
    err := rdb.PFAdd(ctx, key, userID).Err()
    if err != nil {
        return err
    }
    // Expire after 7 days
    rdb.Expire(ctx, key, 7*24*time.Hour)
    return nil
}

func uniqueVisitors(ctx context.Context, rdb *redis.Client, date string) (int64, error) {
    return rdb.PFCount(ctx, "visitors:"+date).Result()
}
```

---

## 4. Distributed Locks Deep Dive

A distributed lock prevents multiple servers from running critical code simultaneously.

### Simple Lock (Single Redis Node)

```go
import "github.com/go-redsync/redsync/v4"

// Redlock algorithm across multiple Redis nodes
func NewDistributedLock(rdb *redis.Client, resource string) *redsync.Mutex {
    rs := redsync.New(goredis.NewPool(rdb))
    return rs.NewMutex("lock:"+resource,
        redsync.WithExpiry(15*time.Second),
        redsync.WithRetryDelay(50*time.Millisecond),
        redsync.WithTries(3),
    )
}

func withLock(ctx context.Context, mutex *redsync.Mutex, fn func() error) error {
    if err := mutex.LockContext(ctx); err != nil {
        return fmt.Errorf("acquire lock: %w", err)
    }
    defer mutex.UnlockContext(ctx)
    return fn()
}

// Usage
mutex := NewDistributedLock(rdb, "payment:user123")
err := withLock(ctx, mutex, func() error {
    // Only one server runs this at a time
    return processPayment(userID, amount)
})
```

---

## 5. Redis as a Job Queue

Redis makes an excellent job queue with visibility timeouts (messages reappear if not acknowledged within N seconds):

```go
type JobQueue struct {
    rdb      *redis.Client
    name     string
    timeout  time.Duration
}

type Job struct {
    ID      string          `json:"id"`
    Payload json.RawMessage `json:"payload"`
}

func (q *JobQueue) Enqueue(ctx context.Context, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    job := Job{
        ID:      fmt.Sprintf("job:%d", time.Now().UnixNano()),
        Payload: data,
    }
    jobData, _ := json.Marshal(job)
    return q.rdb.RPush(ctx, q.name, jobData).Err()
}

func (q *JobQueue) Dequeue(ctx context.Context) (*Job, error) {
    // Move from main queue to processing queue (atomic)
    result, err := q.rdb.BLMove(ctx,
        q.name,         // source
        q.name+":proc", // destination
        "LEFT", "RIGHT",
        q.timeout,
    ).Result()
    if err == redis.Nil {
        return nil, nil // no job available
    }
    if err != nil {
        return nil, err
    }

    var job Job
    json.Unmarshal([]byte(result), &job)
    return &job, nil
}

func (q *JobQueue) Acknowledge(ctx context.Context, jobData string) error {
    // Remove from processing queue
    return q.rdb.LRem(ctx, q.name+":proc", 1, jobData).Err()
}

// RequeueStuck moves jobs stuck in processing back to the main queue
func (q *JobQueue) RequeueStuck(ctx context.Context) error {
    for {
        result := q.rdb.RPopLPush(ctx, q.name+":proc", q.name)
        if result.Err() == redis.Nil {
            break
        }
    }
    return nil
}
```

---

## 6. Mini Project: Activity Feed with Redis Streams

An activity feed that captures user events and lets other services consume them.

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

const (
    streamName = "user-events"
    groupName  = "analytics"
)

type Event struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Metadata  string    `json:"metadata"`
    OccurredAt time.Time `json:"occurred_at"`
}

func main() {
    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    ctx := context.Background()

    // Setup consumer group
    rdb.XGroupCreateMkStream(ctx, streamName, groupName, "0")

    // Producer: publish events
    go runProducer(ctx, rdb)

    // Consumer: process events
    runConsumer(ctx, rdb, "worker-1")
}

func runProducer(ctx context.Context, rdb *redis.Client) {
    events := []Event{
        {UserID: "alice", EventType: "login"},
        {UserID: "bob",   EventType: "purchase", Metadata: `{"amount":99}`},
        {UserID: "alice", EventType: "view",     Metadata: `{"page":"/products"}`},
    }

    for _, e := range events {
        e.OccurredAt = time.Now()
        meta, _ := json.Marshal(e)

        id, err := rdb.XAdd(ctx, &redis.XAddArgs{
            Stream: streamName,
            Values: map[string]interface{}{
                "user_id":    e.UserID,
                "event_type": e.EventType,
                "metadata":   string(meta),
            },
        }).Result()
        if err != nil {
            log.Printf("publish error: %v", err)
            continue
        }
        fmt.Printf("Published event %s: %s\n", id, e.EventType)
        time.Sleep(500 * time.Millisecond)
    }
}

func runConsumer(ctx context.Context, rdb *redis.Client, workerID string) {
    fmt.Printf("Consumer %s started\n", workerID)

    for {
        results, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    groupName,
            Consumer: workerID,
            Streams:  []string{streamName, ">"},
            Count:    10,
            Block:    2 * time.Second,
        }).Result()
        if err == redis.Nil {
            continue // timeout, no new messages
        }
        if err != nil {
            log.Printf("read error: %v", err)
            time.Sleep(time.Second)
            continue
        }

        for _, stream := range results {
            for _, msg := range stream.Messages {
                processEvent(ctx, rdb, workerID, msg)
            }
        }
    }
}

func processEvent(ctx context.Context, rdb *redis.Client, workerID string, msg redis.XMessage) {
    userID := msg.Values["user_id"].(string)
    eventType := msg.Values["event_type"].(string)

    fmt.Printf("[%s] Processing %s: user=%s type=%s\n",
        workerID, msg.ID, userID, eventType)

    // Your processing logic here...
    // e.g., update user stats, send to ClickHouse, trigger emails

    // Acknowledge: message won't be redelivered
    if err := rdb.XAck(ctx, streamName, groupName, msg.ID).Err(); err != nil {
        log.Printf("ack error: %v", err)
    }
}
```

---

## Summary

- Redis Streams = persistent append-only log with consumer groups, like a lightweight Kafka.
- Consumer groups distribute stream processing across multiple workers without duplicates.
- `XACK` marks messages as processed; unacknowledged messages can be reclaimed.
- HyperLogLog counts unique items with < 1% error using only 12 KB of memory.
- Distributed locks via Lua scripts or Redlock prevent concurrent critical-section execution.

### Exercises

**Easy:** Create a Redis Stream for "page view" events. Write a producer that adds 100 events. Write a consumer that reads them and prints each one.

**Medium:** Implement a consumer group with 3 workers reading from the same stream. Add artificial random failures and implement the recovery loop that re-claims stuck messages after 10 seconds.

**Hard:** Build a "dead letter queue" for your Redis stream: if a message fails processing 3 times (track attempt count in message metadata), move it to a separate `failed-events` stream for manual review.
