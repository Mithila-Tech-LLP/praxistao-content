# Chapter 100: NATS and NATS JetStream

NATS is a lightweight, high-performance messaging system built around a simple subject-based pub/sub model. Core NATS offers at-most-once delivery — fast and fire-and-forget. NATS JetStream adds persistence, at-least-once delivery, consumer groups, and replay on top of the same infrastructure.

## Table of Contents

1. [NATS Core Concepts](#1-nats-core-concepts)
2. [Pub/Sub and Queue Groups in Go](#2-pubsub-and-queue-groups-in-go)
3. [Request-Reply Pattern](#3-request-reply-pattern)
4. [JetStream: Streams and Consumers](#4-jetstream-streams-and-consumers)
5. [JetStream Push Consumers](#5-jetstream-push-consumers)
6. [JetStream Pull Consumers](#6-jetstream-pull-consumers)
7. [Ack Strategies and Retry Backoff](#7-ack-strategies)
8. [NATS vs Kafka](#8-nats-vs-kafka)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. NATS Core Concepts

```
Subject hierarchy:  orders.created
                    orders.payment.completed
                    orders.>               (wildcard: all orders subjects)
                    orders.*              (wildcard: one token only)
```

| Concept | Description |
|---------|-------------|
| **Subject** | Hierarchical dot-separated string; the "address" of a message |
| **Publisher** | Sends a message to a subject |
| **Subscriber** | Receives messages matching a subject pattern |
| **Queue Group** | Multiple subscribers share a subject; only one receives each message |
| **At-most-once** | Core NATS has no persistence; if no subscriber is connected, the message is lost |
| **JetStream** | Layer on top of NATS that adds persistence, replay, and consumer groups |

### Subject wildcards

```
orders.*     → matches: orders.created, orders.shipped  (exactly one token)
              does NOT match: orders.payment.completed

orders.>     → matches: orders.created, orders.payment.completed, orders.a.b.c
              (one or more remaining tokens)
```

---

## 2. Pub/Sub and Queue Groups in Go

Install the Go client:

```
go get github.com/nats-io/nats.go
```

### Basic publish and subscribe

```go
package natsexample

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "time"

    "github.com/nats-io/nats.go"
)

type OrderEvent struct {
    EventType string    `json:"event_type"`
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    OccuredAt time.Time `json:"occurred_at"`
}

func main() {
    // Connect — defaults to nats://localhost:4222
    nc, err := nats.Connect(nats.DefaultURL,
        nats.ReconnectWait(2*time.Second),
        nats.MaxReconnects(10),
        nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
            slog.Warn("NATS disconnected", "err", err)
        }),
    )
    if err != nil {
        panic(err)
    }
    defer nc.Drain() // flush and close cleanly

    // Subscribe — fires handler on each matching message
    sub, err := nc.Subscribe("orders.>", func(msg *nats.Msg) {
        var event OrderEvent
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            slog.Error("unmarshal", "err", err)
            return
        }
        slog.Info("received", "subject", msg.Subject, "order_id", event.OrderID)
    })
    if err != nil {
        panic(err)
    }
    defer sub.Unsubscribe()

    // Publish
    event := OrderEvent{
        EventType: "order.created",
        OrderID:   "ord-001",
        UserID:    "usr-42",
        Amount:    149.99,
        OccuredAt: time.Now(),
    }
    data, _ := json.Marshal(event)
    if err := nc.Publish("orders.created", data); err != nil {
        panic(err)
    }

    time.Sleep(100 * time.Millisecond) // let async handler fire
}
```

### Queue subscribe — load-balanced consumers

In a queue group, NATS delivers each message to exactly one subscriber in the group. This is the NATS equivalent of a Kafka consumer group or a RabbitMQ competing consumers pattern.

```go
func startWorker(nc *nats.Conn, workerID int) {
    // All workers with the same group name share the load
    sub, err := nc.QueueSubscribe("orders.payment.>", "payment-workers", func(msg *nats.Msg) {
        var event OrderEvent
        json.Unmarshal(msg.Data, &event)
        slog.Info("processing payment",
            "worker", workerID,
            "order_id", event.OrderID,
        )
        processPayment(event)
    })
    if err != nil {
        slog.Error("queue subscribe", "err", err)
        return
    }
    // Keep sub alive until program exits
    _ = sub
}

// Start 3 workers — NATS will round-robin messages among them
func startPaymentWorkers(nc *nats.Conn) {
    for i := 1; i <= 3; i++ {
        startWorker(nc, i)
    }
}
```

---

## 3. Request-Reply Pattern

NATS has built-in support for synchronous RPC: the publisher sends a message with a reply subject, and the subscriber responds to that subject. NATS handles the temporary reply subject internally.

```go
// ---- Server side: the service that handles requests ----

func startPricingService(nc *nats.Conn) {
    nc.Subscribe("pricing.calculate", func(msg *nats.Msg) {
        var req struct {
            ProductID string `json:"product_id"`
            Quantity  int    `json:"quantity"`
        }
        json.Unmarshal(msg.Data, &req)

        price := calculatePrice(req.ProductID, req.Quantity)

        resp, _ := json.Marshal(map[string]float64{"price": price})
        // Respond to the ephemeral reply subject NATS set up for this call
        msg.Respond(resp)
    })
}

// ---- Client side: caller that wants a synchronous answer ----

func getPrice(nc *nats.Conn, productID string, qty int) (float64, error) {
    reqData, _ := json.Marshal(map[string]interface{}{
        "product_id": productID,
        "quantity":   qty,
    })

    // nc.Request sends the message and waits up to 3s for a reply
    reply, err := nc.Request("pricing.calculate", reqData, 3*time.Second)
    if err != nil {
        return 0, fmt.Errorf("request pricing: %w", err)
    }

    var resp struct {
        Price float64 `json:"price"`
    }
    json.Unmarshal(reply.Data, &resp)
    return resp.Price, nil
}
```

Request-Reply over NATS is useful for microservice calls where you need an answer immediately, without standing up a full HTTP server.

---

## 4. JetStream: Streams and Consumers

JetStream stores messages in **streams**. A stream captures messages from one or more subjects and retains them for a configurable period or size limit.

```
NATS Core (no persistence):
  Publisher → Subject → Subscriber (must be connected right now)

JetStream (persistent):
  Publisher → Subject → [Stream: ORDERS] → Consumer → Subscriber
                             │
                        stored on disk
                        replayed on demand
                        multiple consumers, each with own cursor
```

```go
package jetstream

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/nats-io/nats.go"
)

func connectJetStream(url string) (nats.JetStreamContext, *nats.Conn, error) {
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, nil, err
    }

    js, err := nc.JetStream()
    if err != nil {
        nc.Close()
        return nil, nil, err
    }
    return js, nc, nil
}

func setupOrderStream(js nats.JetStreamContext) error {
    // AddStream is idempotent — safe to call on every startup
    _, err := js.AddStream(&nats.StreamConfig{
        Name:     "ORDERS",
        Subjects: []string{"orders.>"},   // capture all orders.* subjects
        Storage:  nats.FileStorage,       // persist to disk (vs MemoryStorage)
        MaxAge:   7 * 24 * time.Hour,     // retain for 7 days
        MaxBytes: 1 << 30,               // max 1 GB total
        Replicas: 1,                      // 3 for production cluster
    })
    if err != nil {
        return fmt.Errorf("add stream: %w", err)
    }
    return nil
}

func publishOrder(js nats.JetStreamContext, event OrderEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    // PublishAsync is non-blocking; use Publish for synchronous confirmation
    ack, err := js.Publish("orders.created", data)
    if err != nil {
        return fmt.Errorf("js publish: %w", err)
    }
    // ack.Sequence is the sequence number assigned by the stream
    fmt.Printf("published to stream %s seq %d\n", ack.Stream, ack.Sequence)
    return nil
}
```

---

## 5. JetStream Push Consumers

A push consumer has the server deliver messages to a subscriber automatically. The consumer cursor (which messages have been processed) is maintained by the server.

```go
func startPushConsumer(js nats.JetStreamContext) error {
    // Durable name: the server remembers this consumer's position across restarts
    sub, err := js.Subscribe(
        "orders.>",
        func(msg *nats.Msg) {
            var event OrderEvent
            if err := json.Unmarshal(msg.Data, &event); err != nil {
                msg.Term() // permanent failure — don't retry
                return
            }

            if err := processOrderEvent(event); err != nil {
                // Retry with 5-second delay
                msg.NakWithDelay(5 * time.Second)
                return
            }

            msg.Ack()
        },
        nats.Durable("order-processor"),     // durable name — survives restarts
        nats.DeliverNew(),                    // start from new messages (use DeliverAll() to replay)
        nats.AckExplicit(),                   // require manual ack
        nats.MaxDeliver(5),                   // retry up to 5 times before going to DLQ
        nats.AckWait(30*time.Second),         // if not acked in 30s, redeliver
    )
    if err != nil {
        return fmt.Errorf("subscribe: %w", err)
    }
    _ = sub

    // Block until context cancelled
    select {}
}
```

### Multiple push consumers (fan-out)

```go
// Two different durable consumers on the same stream — each gets ALL messages independently.
// This is unlike a queue group where messages are shared.

func startNotificationConsumer(js nats.JetStreamContext) {
    js.Subscribe("orders.>",
        handleNotification,
        nats.Durable("notification-service"),
        nats.AckExplicit(),
    )
}

func startAuditConsumer(js nats.JetStreamContext) {
    js.Subscribe("orders.>",
        handleAudit,
        nats.Durable("audit-service"),
        nats.AckExplicit(),
    )
}
// Both notification-service and audit-service receive every orders.> message.
```

---

## 6. JetStream Pull Consumers

Pull consumers let the application control when to fetch messages. This gives back-pressure: the consumer only asks for as many messages as it can handle right now.

```go
func startPullWorker(js nats.JetStreamContext) error {
    // Create a durable pull subscriber
    sub, err := js.PullSubscribe(
        "orders.payment.>",
        "payment-worker-group", // durable name
    )
    if err != nil {
        return fmt.Errorf("pull subscribe: %w", err)
    }
    defer sub.Unsubscribe()

    for {
        // Fetch up to 10 messages; wait up to 1s if fewer are available
        msgs, err := sub.Fetch(10, nats.MaxWait(1*time.Second))
        if err != nil {
            if err == nats.ErrTimeout {
                continue // no messages right now — loop and try again
            }
            return fmt.Errorf("fetch: %w", err)
        }

        for _, msg := range msgs {
            if err := processPaymentMessage(msg); err != nil {
                msg.NakWithDelay(10 * time.Second)
                continue
            }
            msg.Ack()
        }
    }
}

func processPaymentMessage(msg *nats.Msg) error {
    var event OrderEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    // ... charge payment, update DB, etc.
    return chargePayment(event)
}
```

### Concurrent batch processing with pull

```go
func startConcurrentPullWorkers(js nats.JetStreamContext, workerCount int) {
    sub, _ := js.PullSubscribe("orders.>", "batch-processor")

    sem := make(chan struct{}, workerCount) // limit concurrency

    for {
        msgs, err := sub.Fetch(workerCount, nats.MaxWait(500*time.Millisecond))
        if err != nil {
            continue
        }

        var wg sync.WaitGroup
        for _, msg := range msgs {
            msg := msg
            wg.Add(1)
            sem <- struct{}{} // acquire slot
            go func() {
                defer wg.Done()
                defer func() { <-sem }() // release slot
                if err := processMessage(msg); err != nil {
                    msg.Nak()
                } else {
                    msg.Ack()
                }
            }()
        }
        wg.Wait()
    }
}
```

---

## 7. Ack Strategies

JetStream supports several ack types that give fine-grained control over retries.

```go
func handleMessage(msg *nats.Msg) {
    var event OrderEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        // Permanent failure — message is bad, do not retry
        // Moves to the consumer's "dead letter" (MaxDeliver exhausted → sent to advisory subject)
        msg.Term()
        return
    }

    // Check idempotency: have we already processed this?
    if alreadyProcessed(event.OrderID) {
        msg.Ack() // Ack to remove from delivery queue
        return
    }

    err := processOrder(event)
    switch {
    case err == nil:
        // Success
        msg.Ack()

    case isTransient(err):
        // Retry after 5 seconds (server will redeliver then)
        msg.NakWithDelay(5 * time.Second)

    case isPermanent(err):
        // No point retrying — terminate this message
        msg.Term()

    default:
        // Default Nak: use the consumer's configured AckWait before redelivery
        msg.Nak()
    }
}
```

| Ack Type | Call | Effect |
|----------|------|--------|
| **Ack** | `msg.Ack()` | Message processed successfully; remove from pending |
| **Nak** | `msg.Nak()` | Redeliver after AckWait period |
| **NakWithDelay** | `msg.NakWithDelay(d)` | Redeliver after specific delay (backoff) |
| **Term** | `msg.Term()` | Permanent failure; stop retrying; increment deliveries count to MaxDeliver |
| **InProgress** | `msg.InProgress()` | Reset the AckWait timer; message is still being processed |

### JetStream dead letter pattern

When `MaxDeliver` is exhausted, JetStream publishes an advisory to `$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.ORDERS.order-processor`. Subscribe to this to handle permanently-failed messages.

```go
// Listen for max-delivery advisories
nc.Subscribe("$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.ORDERS.*", func(msg *nats.Msg) {
    var advisory struct {
        Stream    string `json:"stream"`
        Consumer  string `json:"consumer"`
        Subject   string `json:"subject"`
        Sequences []uint64 `json:"stream_seq"`
    }
    json.Unmarshal(msg.Data, &advisory)
    slog.Error("message exhausted retries",
        "stream", advisory.Stream,
        "consumer", advisory.Consumer,
        "sequences", advisory.Sequences,
    )
    // Fetch from stream by sequence and move to a DLQ stream
})
```

---

## 8. NATS vs Kafka

| Dimension | NATS | Kafka |
|-----------|------|-------|
| **Delivery (core)** | At-most-once | At-least-once |
| **Delivery (JetStream)** | At-least-once | At-least-once / exactly-once |
| **Latency** | Sub-millisecond | Low millisecond |
| **Infrastructure** | Single binary, minimal config | Kafka + ZooKeeper/KRaft, more ops overhead |
| **Throughput** | Very high | Extremely high (millions/sec) |
| **Subject routing** | Rich wildcard subjects | Simple topic + partition key |
| **Long-term storage** | JetStream up to configured limits | Designed for very long retention |
| **Consumer groups** | Queue groups (core) or durable consumers (JetStream) | Native consumer groups with offset tracking |
| **Request-Reply** | Built-in (`nc.Request`) | Not native |
| **Best for** | Microservice messaging, IoT, low-latency pipelines | Large-scale event streaming, audit logs, data pipelines |

### Rule of thumb

- Use **NATS** when you want simple, fast messaging without the Kafka ops burden — especially for microservice-to-microservice communication or IoT.
- Use **Kafka** when you need very long message retention, complex consumer groups at scale, or tight ecosystem integrations (Kafka Connect, ksqlDB).

---

## Summary

- **Core NATS**: subject-based pub/sub, at-most-once, no persistence — messages are lost if no subscriber is connected
- **Subjects**: hierarchical dot-notation; `*` matches one token, `>` matches one or more
- **Queue groups**: `nc.QueueSubscribe` — one message delivered to exactly one subscriber in the group (load balancing)
- **Request-Reply**: `nc.Request` blocks until a subscriber responds — synchronous RPC without HTTP
- **JetStream streams**: capture subjects on disk with configurable retention; replay is possible
- **Push consumers**: server delivers messages automatically; `Durable` name preserves position across restarts
- **Pull consumers**: `sub.Fetch(n)` gives back-pressure control; application decides when it is ready for more
- **Ack strategies**: `Ack`, `Nak`, `NakWithDelay` (backoff), `Term` (permanent failure), `InProgress` (heartbeat)
- NATS is simpler and lower-latency than Kafka; Kafka excels at massive scale and long-term event storage

## Exercises

### Easy
1. Start NATS with Docker (`docker run -p 4222:4222 nats`). Write a publisher that sends 20 `OrderEvent` messages to subject `orders.created`. Write a subscriber on `orders.>` that prints the subject and order ID of every message received.
2. Start 3 queue subscribers on `orders.payment.>` with group name `"payment-workers"`. Publish 30 messages and verify each message goes to exactly one worker (by printing the goroutine ID or worker number in the handler).
3. Implement a Request-Reply service: a `pricing.calculate` handler that accepts `{product_id, quantity}` and returns `{price}`. Write a client that calls it 5 times and prints the results.

### Medium
4. Create a JetStream stream named `ORDERS` capturing `orders.>`. Publish 50 order events, then stop the publisher. Start a durable push consumer `nats.Durable("audit-consumer")` with `nats.DeliverAll()` — verify it replays all 50 messages from the beginning.
5. Implement a **pull consumer** that processes messages in batches of 5. If processing fails, call `msg.NakWithDelay(10 * time.Second)`. Add a counter tracking how many times each order ID was delivered; log a warning when an order is retried more than twice.
6. Simulate a slow consumer: make the handler `time.Sleep(200 * time.Millisecond)`. Use `msg.InProgress()` every 100ms to reset the AckWait timer and prevent redelivery. Verify the message is not re-delivered while being processed.

### Hard
7. Build a **NATS-based order processing pipeline** with 3 stages connected by subjects: `pipeline.stage1`, `pipeline.stage2`, `pipeline.stage3`. Each stage is a queue group consumer that enriches the event and publishes it to the next subject. Use JetStream so events survive a service restart mid-pipeline. Track end-to-end latency from stage 1 entry to stage 3 completion.
8. Implement an **idempotent JetStream consumer** backed by Redis: before processing, check a Redis set for the message sequence number. If present, ack and skip. After processing, add the sequence to Redis with a 24-hour TTL. Benchmark the overhead of the Redis check vs. total processing time.
