# Chapter 54: NATS — Cloud-Native Messaging

NATS is the simplest, fastest message broker you can use. Where Kafka needs a cluster and ZooKeeper, NATS is a single binary. Where RabbitMQ needs exchanges and bindings, NATS uses simple subject strings. Yet NATS JetStream (its persistent streaming layer) competes with Kafka on durability.

## Table of Contents

1. Why NATS?
2. Core NATS: Pub/Sub and Request/Reply
3. NATS JetStream: Persistent Streams
4. Docker Setup
5. Building with NATS in Go
6. JetStream in Go
7. NATS vs Kafka vs RabbitMQ
8. Mini Project: Microservice Event Bus
9. Exercises

---

## 1. Why NATS?

**NATS Philosophy:**
- Simple: one binary, no dependencies
- Fast: 10+ million messages/second on a single server
- Lightweight: uses < 20 MB RAM
- Cloud-native: built for containers, Kubernetes, distributed systems

**When to use NATS:**
- Microservices communication (service mesh alternative)
- IoT devices sending telemetry
- Real-time features: notifications, presence, collaborative editing
- Simple event bus that doesn't need Kafka's complexity

**When NOT to use NATS (core):**
- Long message retention (use JetStream or Kafka)
- Guaranteed delivery if consumers are offline (use JetStream or Kafka)

---

## 2. Core NATS: Pub/Sub and Request/Reply

**Subjects:** NATS uses dot-separated subjects instead of topics/queues.
```
"orders.placed"          → exact match
"orders.*"               → wildcard: matches "orders.placed", "orders.shipped"
"orders.>"               → wildcard: matches any suffix: "orders.placed.eu.london"
```

**Core NATS delivery semantics:**
- At-most-once: if a subscriber is offline when a message is published, it's lost.
- Fire and forget.
- Perfect for real-time events where missing one doesn't matter (presence updates, metric pushes).

**Request/Reply:**

NATS has built-in request-reply — one of its most powerful features:

```
Client: publish to "orders.create", reply-to="_INBOX.abc123"
Server: subscribes to "orders.create", receives message, publishes response to "_INBOX.abc123"
Client: receives response from "_INBOX.abc123"
```

This enables synchronous RPC over async pub/sub.

---

## 3. NATS JetStream: Persistent Streams

JetStream is NATS's persistence layer — think of it as NATS + Kafka semantics:

- **Streams:** Named, persistent logs of messages. Subjects flow into streams.
- **Consumers:** Read from streams at their own pace, with acknowledgements.
- **Retention policies:** Limits (time, count, bytes), and compaction modes.
- **Delivery guarantees:** At-least-once (with ACK), or exactly-once (with deduplication).

```
JetStream Stream "ORDERS":
  Subject: "orders.>"  → captures all order events
  Retention: 7 days
  MaxBytes: 100 GB
  Replicas: 3

Consumer "EmailService":
  Filter: "orders.placed"
  AckPolicy: explicit
  DeliverPolicy: all (from beginning)

Consumer "Analytics":
  Filter: "orders.>"
  AckPolicy: none (fire-and-forget)
  DeliverPolicy: new (only future messages)
```

---

## 4. Docker Setup

```bash
docker run -d \
  --name nats \
  -p 4222:4222 \
  -p 8222:8222 \
  nats:2.10 \
  -js \     # enable JetStream
  -m 8222   # enable HTTP monitoring

# Monitor: http://localhost:8222
```

```bash
go get github.com/nats-io/nats.go
```

---

## 5. Building with NATS in Go

**Basic Pub/Sub:**

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"
)

func main() {
    nc, err := nats.Connect("nats://localhost:4222",
        nats.MaxReconnects(10),
        nats.ReconnectWait(2*time.Second),
        nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
            log.Printf("NATS disconnected: %v", err)
        }),
        nats.ReconnectHandler(func(nc *nats.Conn) {
            log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Drain() // wait for pending messages before closing

    // Subscribe to all order events
    nc.Subscribe("orders.*", func(msg *nats.Msg) {
        fmt.Printf("Received on %s: %s\n", msg.Subject, string(msg.Data))
    })

    // Subscribe with queue group (only one subscriber in the group gets each message)
    nc.QueueSubscribe("tasks", "workers", func(msg *nats.Msg) {
        fmt.Printf("Worker processing: %s\n", string(msg.Data))
    })

    // Publish messages
    nc.Publish("orders.placed", []byte(`{"order_id":"123","amount":99.99}`))
    nc.Publish("orders.shipped", []byte(`{"order_id":"123","tracking":"ABC"}`))
    nc.Publish("tasks", []byte(`{"task":"resize-image","file":"photo.jpg"}`))

    time.Sleep(time.Second)
}
```

**Request/Reply:**

```go
// Server side: handle requests
nc.Subscribe("math.add", func(msg *nats.Msg) {
    // Parse request
    var req struct{ A, B int }
    json.Unmarshal(msg.Data, &req)

    // Send response back to msg.Reply subject
    result := req.A + req.B
    response, _ := json.Marshal(map[string]int{"result": result})
    nc.Publish(msg.Reply, response)
})

// Client side: send request and wait for reply
var resp struct{ Result int }
msg, err := nc.Request("math.add",
    []byte(`{"a":5,"b":3}`),
    2*time.Second,  // timeout
)
if err != nil {
    log.Fatal("request timeout:", err)
}
json.Unmarshal(msg.Data, &resp)
fmt.Printf("5 + 3 = %d\n", resp.Result) // 5 + 3 = 8
```

---

## 6. JetStream in Go

```go
// Connect and enable JetStream
nc, _ := nats.Connect("nats://localhost:4222")
js, err := nc.JetStream()
if err != nil {
    log.Fatal(err)
}

// Create a stream
_, err = js.AddStream(&nats.StreamConfig{
    Name:     "ORDERS",
    Subjects: []string{"orders.>"},
    MaxAge:   7 * 24 * time.Hour,  // retain for 7 days
    Replicas: 1,                   // 3 for production
})
if err != nil {
    log.Println("stream may already exist:", err)
}

// Publish to JetStream (persistent, acknowledged)
ack, err := js.Publish("orders.placed", []byte(`{"order_id":"456","amount":149.99}`))
if err != nil {
    log.Fatal("publish:", err)
}
fmt.Printf("Published to stream %s, seq %d\n", ack.Stream, ack.Sequence)

// Subscribe with push consumer (messages pushed to us)
sub, err := js.Subscribe("orders.placed",
    func(msg *nats.Msg) {
        fmt.Printf("Processing: %s\n", string(msg.Data))
        msg.Ack() // acknowledge: tell JetStream we processed it
    },
    nats.Durable("email-service"),         // named consumer — resumes on restart
    nats.AckExplicit(),                    // manual ack required
    nats.DeliverAll(),                     // deliver all messages from beginning
    nats.MaxDeliver(3),                    // max 3 delivery attempts
    nats.AckWait(30*time.Second),          // redelivery if no ack within 30s
)
if err != nil {
    log.Fatal("subscribe:", err)
}
defer sub.Unsubscribe()

// Pull consumer (we fetch messages on demand)
pullSub, err := js.PullSubscribe("orders.>", "analytics",
    nats.BindStream("ORDERS"))
if err != nil {
    log.Fatal("pull subscribe:", err)
}

// Fetch up to 10 messages, wait up to 500ms
msgs, err := pullSub.Fetch(10, nats.MaxWait(500*time.Millisecond))
if err != nil && err != nats.ErrTimeout {
    log.Println("fetch:", err)
} else {
    for _, m := range msgs {
        fmt.Printf("Pulled: %s\n", string(m.Data))
        m.Ack()
    }
}
```

---

## 7. NATS vs Kafka vs RabbitMQ

| Criteria | NATS Core | NATS JetStream | Kafka | RabbitMQ |
|----------|-----------|----------------|-------|----------|
| Persistence | No | Yes | Yes | Yes |
| Replay | No | Yes | Yes | No |
| Throughput | 10M+/s | 1M+/s | 1M+/s | 100K/s |
| Latency | < 100µs | < 1ms | ~1ms | 1-10ms |
| Complexity | Very low | Medium | High | Medium |
| Best for | Real-time, IoT | Event streaming | Big data, high volume | Task queues, routing |

---

## 8. Mini Project: Microservice Event Bus

A simple event bus using NATS for microservice communication:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"
)

// EventBus wraps NATS for service communication
type EventBus struct {
    nc *nats.Conn
    js nats.JetStreamContext
}

func NewEventBus(url string) (*EventBus, error) {
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }
    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    // Create streams
    js.AddStream(&nats.StreamConfig{
        Name:     "EVENTS",
        Subjects: []string{"service.>"},
        MaxAge:   24 * time.Hour,
    })

    return &EventBus{nc: nc, js: js}, nil
}

func (bus *EventBus) Publish(subject string, data interface{}) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return err
    }
    _, err = bus.js.Publish(subject, payload)
    return err
}

func (bus *EventBus) Subscribe(subject, groupName string, handler func([]byte)) error {
    _, err := bus.js.QueueSubscribe(subject, groupName,
        func(msg *nats.Msg) {
            handler(msg.Data)
            msg.Ack()
        },
        nats.Durable(groupName),
        nats.AckExplicit(),
    )
    return err
}

func (bus *EventBus) Request(subject string, data interface{}, timeout time.Duration) ([]byte, error) {
    payload, _ := json.Marshal(data)
    msg, err := bus.nc.Request(subject, payload, timeout)
    if err != nil {
        return nil, err
    }
    return msg.Data, nil
}

// Example: user service
func userService(bus *EventBus) {
    // Handle user creation requests
    bus.nc.Subscribe("users.create", func(msg *nats.Msg) {
        var req struct{ Name, Email string }
        json.Unmarshal(msg.Data, &req)

        userID := fmt.Sprintf("user-%d", time.Now().UnixNano())
        fmt.Printf("[UserService] Created user %s (%s)\n", userID, req.Name)

        resp, _ := json.Marshal(map[string]string{"user_id": userID})
        bus.nc.Publish(msg.Reply, resp)

        // Emit event
        bus.Publish("service.users.created", map[string]string{
            "user_id": userID,
            "email":   req.Email,
        })
    })
}

// Example: email service
func emailService(bus *EventBus) {
    bus.Subscribe("service.users.created", "email-service", func(data []byte) {
        var event struct {
            UserID string `json:"user_id"`
            Email  string `json:"email"`
        }
        json.Unmarshal(data, &event)
        fmt.Printf("[EmailService] Sending welcome email to %s\n", event.Email)
    })
}

func main() {
    bus, err := NewEventBus("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }

    userService(bus)
    emailService(bus)

    // Test: create a user
    resp, err := bus.Request("users.create",
        map[string]string{"name": "Alice", "email": "alice@example.com"},
        2*time.Second,
    )
    if err != nil {
        log.Fatal("request:", err)
    }

    var result map[string]string
    json.Unmarshal(resp, &result)
    fmt.Printf("Created user: %s\n", result["user_id"])

    time.Sleep(time.Second)
}
```

---

## Summary

- NATS uses dot-separated subjects with `*` (one token) and `>` (remaining tokens) wildcards.
- Core NATS = at-most-once pub/sub. Simple, fast, perfect for real-time features.
- Queue subscriptions = multiple subscribers sharing a subject like a work queue.
- Request/Reply is built-in — perfect for microservice RPC.
- JetStream adds persistence, consumer groups, acknowledgements, and replay — Kafka-like guarantees.

### Exercises

**Easy:** Set up NATS and write a pub/sub example: publisher sends 10 messages to "news.headlines", subscriber prints each. Add a wildcard subscriber to "news.>" that counts all messages.

**Medium:** Implement a microservice mesh: 3 services (users, orders, inventory) communicating via NATS Request/Reply. Each service has an in-memory store. Placing an order calls the inventory service (reserve stock) and the user service (get email), then returns the combined result.

**Hard:** Implement a JetStream-based event sourcing system for a bank account: events are "deposited", "withdrawn", "transferred". Each event is published to JetStream. A consumer replays all events to reconstruct the current balance. Implement `GetBalance(accountID)` as a NATS Request that triggers a balance replay.
