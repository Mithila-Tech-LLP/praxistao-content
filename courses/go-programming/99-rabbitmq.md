# Chapter 99: RabbitMQ with Go (amqp091-go)

RabbitMQ is a message broker built on the AMQP protocol. Unlike Kafka — which retains messages in a log — RabbitMQ routes and delivers messages to consumers and removes them once acknowledged. This makes it ideal for task queues, work distribution, and complex routing scenarios.

## Table of Contents

1. [Core Concepts](#1-core-concepts)
2. [Exchange Types](#2-exchange-types)
3. [Producer in Go](#3-producer-in-go)
4. [Consumer in Go](#4-consumer-in-go)
5. [Dead Letter Exchange (DLX)](#5-dead-letter-exchange-dlx)
6. [Publisher Confirms](#6-publisher-confirms)
7. [Connection and Channel Pooling](#7-connection-and-channel-pooling)
8. [RabbitMQ vs Kafka](#8-rabbitmq-vs-kafka)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Core Concepts

AMQP (Advanced Message Queuing Protocol) defines a network protocol and a set of broker-side concepts.

```
Producer → Exchange → (Binding with routing key) → Queue → Consumer
```

| Concept | Description |
|---------|-------------|
| **Broker** | The RabbitMQ server — receives, routes, and stores messages |
| **Exchange** | Entry point for messages; decides which queues get them |
| **Queue** | Buffer that stores messages until a consumer picks them up |
| **Binding** | A link from an exchange to a queue with an optional routing key |
| **Routing Key** | A string the producer attaches; the exchange uses it to route |
| **Channel** | Lightweight virtual connection inside a TCP connection |
| **Acknowledgement** | Consumer signals "I processed this" — broker deletes the message |

A **connection** is one TCP connection. Each connection can have many **channels**. Creating a channel is cheap; creating a connection is expensive.

---

## 2. Exchange Types

### Direct Exchange — exact routing key match

```
Producer publishes key="orders.payment"
                       │
              [direct exchange: orders]
                       │
          ┌────────────┴────────────┐
          │                         │
   binding:"orders.payment"  binding:"orders.shipping"
          │                         │
   [queue: payment-svc]      [queue: shipping-svc]
          ▼
   only payment-svc receives it
```

### Fanout Exchange — broadcast to all bound queues

```
Producer publishes (routing key ignored)
                       │
             [fanout exchange: logs]
                       │
       ┌───────────────┼───────────────┐
       ▼               ▼               ▼
[queue: log-db]  [queue: log-es]  [queue: log-console]
   (all three receive every message)
```

### Topic Exchange — wildcard routing

```
Routing key pattern: <level>.<service>.<action>

Producer publishes:
  "error.payment-svc.timeout"
  "info.order-svc.created"
  "error.auth-svc.failed"

Bindings:
  "error.#"    → queue: alerts        (all error messages)
  "#.created"  → queue: audit-log     (all created events)
  "*.payment-svc.*" → queue: payment  (any level, any action)

Wildcards:
  *  matches exactly one word
  #  matches zero or more words
```

### Headers Exchange

Headers exchange ignores the routing key; it matches on message headers instead.

```
Binding: {x-match: "all", service: "payment", version: "2"}
Message headers: {service: "payment", version: "2", region: "us-east"}
→ Matches (all required headers present)

Binding: {x-match: "any", service: "payment", service: "order"}
→ Matches if either header is present
```

---

## 3. Producer in Go

Install the driver:

```
go get github.com/rabbitmq/amqp091-go
```

```go
package rabbitmq

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewProducer(url string) (*Producer, error) {
    conn, err := amqp.Dial(url) // e.g. "amqp://guest:guest@localhost:5672/"
    if err != nil {
        return nil, fmt.Errorf("dial: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("open channel: %w", err)
    }

    // Declare the exchange (idempotent — safe to call on every startup)
    err = ch.ExchangeDeclare(
        "orders",  // name
        "topic",   // kind: direct | fanout | topic | headers
        true,      // durable: survives broker restart
        false,     // auto-delete
        false,     // internal
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("exchange declare: %w", err)
    }

    return &Producer{conn: conn, channel: ch}, nil
}

type OrderEvent struct {
    EventType string    `json:"event_type"`
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    OccuredAt time.Time `json:"occurred_at"`
}

func (p *Producer) Publish(ctx context.Context, routingKey string, event OrderEvent) error {
    body, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal: %w", err)
    }

    return p.channel.PublishWithContext(
        ctx,
        "orders",   // exchange name
        routingKey, // e.g. "orders.payment.completed"
        false,      // mandatory: return if no queue bound
        false,      // immediate: return if no consumer ready
        amqp.Publishing{
            ContentType:  "application/json",
            DeliveryMode: amqp.Persistent, // survives broker restart (2 = persistent)
            Timestamp:    time.Now(),
            MessageId:    fmt.Sprintf("%s-%d", event.OrderID, time.Now().UnixNano()),
            Body:         body,
        },
    )
}

func (p *Producer) Close() {
    p.channel.Close()
    p.conn.Close()
}
```

### Publishing to different exchange types

```go
// Direct exchange: exact match
p.channel.PublishWithContext(ctx, "orders-direct", "payment", false, false, msg)

// Fanout exchange: routing key is ignored
p.channel.PublishWithContext(ctx, "notifications-fanout", "", false, false, msg)

// Topic exchange: wildcard matching
p.channel.PublishWithContext(ctx, "logs-topic", "error.payment-svc.timeout", false, false, msg)
```

---

## 4. Consumer in Go

```go
package rabbitmq

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    logger  *slog.Logger
}

func NewConsumer(url string, logger *slog.Logger) (*Consumer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, fmt.Errorf("dial: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("open channel: %w", err)
    }

    // Limit how many unacknowledged messages this consumer holds at once
    // Without this, RabbitMQ pushes ALL queued messages to one fast consumer
    if err := ch.Qos(10, 0, false); err != nil {
        return nil, fmt.Errorf("qos: %w", err)
    }

    return &Consumer{conn: conn, channel: ch, logger: logger}, nil
}

func (c *Consumer) Setup(exchange, queue, routingKey string) error {
    // Declare exchange
    if err := c.channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
        return fmt.Errorf("exchange declare: %w", err)
    }

    // Declare queue — durable so it survives broker restart
    q, err := c.channel.QueueDeclare(
        queue, // name
        true,  // durable
        false, // auto-delete when unused
        false, // exclusive to this connection
        false, // no-wait
        nil,   // arguments (use amqp.Table for DLX, TTL etc.)
    )
    if err != nil {
        return fmt.Errorf("queue declare: %w", err)
    }

    // Bind queue to exchange with routing key pattern
    return c.channel.QueueBind(
        q.Name,     // queue name
        routingKey, // binding key, e.g. "orders.payment.#"
        exchange,   // exchange name
        false,
        nil,
    )
}

type Handler func(ctx context.Context, event OrderEvent) error

func (c *Consumer) Consume(ctx context.Context, queue string, handler Handler) error {
    msgs, err := c.channel.Consume(
        queue,   // queue name
        "",      // consumer tag (auto-generated)
        false,   // auto-ack: false means we ack manually
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,
    )
    if err != nil {
        return fmt.Errorf("consume: %w", err)
    }

    for {
        select {
        case <-ctx.Done():
            return nil
        case d, ok := <-msgs:
            if !ok {
                return fmt.Errorf("channel closed")
            }
            c.processDelivery(ctx, d, handler)
        }
    }
}

func (c *Consumer) processDelivery(ctx context.Context, d amqp.Delivery, handler Handler) {
    var event OrderEvent
    if err := json.Unmarshal(d.Body, &event); err != nil {
        c.logger.Error("unmarshal failed", "err", err, "msg_id", d.MessageId)
        // Reject and do NOT requeue — malformed messages will loop forever
        _ = d.Nack(false, false)
        return
    }

    if err := handler(ctx, event); err != nil {
        c.logger.Error("handler failed",
            "err", err,
            "order_id", event.OrderID,
            "delivery_tag", d.DeliveryTag,
        )
        // Nack with requeue=true — message goes back to the front of the queue
        // WARNING: if the error is permanent, this creates an infinite loop; use DLX instead
        _ = d.Nack(false, true)
        return
    }

    // Ack: false means ack only this message (not all prior unacked)
    if err := d.Ack(false); err != nil {
        c.logger.Warn("ack failed", "err", err)
    }
}

func (c *Consumer) Close() {
    c.channel.Close()
    c.conn.Close()
}
```

---

## 5. Dead Letter Exchange (DLX)

When a message is rejected (Nack with requeue=false), expires, or the queue is full, RabbitMQ can route it to a **Dead Letter Exchange** instead of discarding it.

```
Normal flow:
  Producer → [exchange: orders] → [queue: payment-tasks] → Consumer
                                                                │
                                                    handler fails 3 times
                                                                │
                                                        Nack(requeue=false)
                                                                │
                                                                ▼
                                  [DLX exchange: orders.dlx] → [queue: payment-tasks.dlq]
                                                                │
                                                     DLQ consumer investigates
```

### Configure DLX at queue declaration time

```go
func (c *Consumer) SetupWithDLX(exchange, queue, routingKey string) error {
    // Step 1: declare the dead letter exchange
    if err := c.channel.ExchangeDeclare("orders.dlx", "direct", true, false, false, false, nil); err != nil {
        return err
    }

    // Step 2: declare the dead letter queue
    if _, err := c.channel.QueueDeclare("payment-tasks.dlq", true, false, false, false, nil); err != nil {
        return err
    }

    // Step 3: bind DLQ to DLX
    if err := c.channel.QueueBind("payment-tasks.dlq", "payment-tasks", "orders.dlx", false, nil); err != nil {
        return err
    }

    // Step 4: declare the main queue with DLX argument
    _, err := c.channel.QueueDeclare(
        queue,
        true, false, false, false,
        amqp.Table{
            "x-dead-letter-exchange":    "orders.dlx",
            "x-dead-letter-routing-key": queue,  // DLX routing key for the dead message
            "x-message-ttl":             int32(30000), // optional: expire after 30s
        },
    )
    return err
}
```

### Consumer with retry counter

```go
const maxRetries = 3

func (c *Consumer) processWithRetryHeader(ctx context.Context, d amqp.Delivery, handler Handler) {
    var event OrderEvent
    if err := json.Unmarshal(d.Body, &event); err != nil {
        _ = d.Nack(false, false) // malformed — send to DLQ immediately
        return
    }

    // Read retry count from headers (RabbitMQ adds x-death header automatically)
    retries := 0
    if xDeath, ok := d.Headers["x-death"].([]interface{}); ok && len(xDeath) > 0 {
        if table, ok := xDeath[0].(amqp.Table); ok {
            if count, ok := table["count"].(int64); ok {
                retries = int(count)
            }
        }
    }

    if err := handler(ctx, event); err != nil {
        c.logger.Error("handler failed", "retries", retries, "err", err)
        if retries >= maxRetries {
            // Exhausted retries — send to DLQ, do not requeue
            _ = d.Nack(false, false)
        } else {
            // Requeue for another attempt
            _ = d.Nack(false, true)
        }
        return
    }

    _ = d.Ack(false)
}
```

---

## 6. Publisher Confirms

By default, `Publish` returns as soon as the message leaves the Go process. The broker might still lose it before writing to disk. **Publisher confirms** make the broker send an ack back.

```go
func NewReliableProducer(url string) (*Producer, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    // Put channel in confirm mode
    if err := ch.Confirm(false); err != nil {
        return nil, fmt.Errorf("confirm mode: %w", err)
    }

    return &Producer{conn: conn, channel: ch}, nil
}

func (p *Producer) PublishReliable(ctx context.Context, routingKey string, event OrderEvent) error {
    // Get a channel for broker confirmations — must be done BEFORE publishing
    confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

    body, _ := json.Marshal(event)
    if err := p.channel.PublishWithContext(ctx, "orders", routingKey, false, false,
        amqp.Publishing{
            ContentType:  "application/json",
            DeliveryMode: amqp.Persistent,
            Body:         body,
        },
    ); err != nil {
        return fmt.Errorf("publish: %w", err)
    }

    // Wait for broker acknowledgement
    select {
    case confirm := <-confirms:
        if !confirm.Ack {
            return fmt.Errorf("broker nacked message (delivery tag %d)", confirm.DeliveryTag)
        }
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## 7. Connection and Channel Pooling

Creating a connection is expensive (~100ms). Channels are cheap but not goroutine-safe — each goroutine needs its own channel.

```go
package rabbitmq

import (
    "fmt"
    "sync"

    amqp "github.com/rabbitmq/amqp091-go"
)

// ChannelPool holds a single connection and issues channels from it.
// Each caller gets its own channel (channels are not goroutine-safe).
type ChannelPool struct {
    mu      sync.Mutex
    conn    *amqp.Connection
    url     string
}

func NewChannelPool(url string) (*ChannelPool, error) {
    p := &ChannelPool{url: url}
    if err := p.reconnect(); err != nil {
        return nil, err
    }
    return p, nil
}

func (p *ChannelPool) reconnect() error {
    conn, err := amqp.Dial(p.url)
    if err != nil {
        return fmt.Errorf("dial: %w", err)
    }
    p.conn = conn
    return nil
}

// Acquire returns a new channel. Caller must close it after use.
func (p *ChannelPool) Acquire() (*amqp.Channel, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.conn.IsClosed() {
        if err := p.reconnect(); err != nil {
            return nil, err
        }
    }

    ch, err := p.conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("open channel: %w", err)
    }
    return ch, nil
}

func (p *ChannelPool) Close() error {
    return p.conn.Close()
}

// Usage: one goroutine, one channel
func worker(pool *ChannelPool, job Job) error {
    ch, err := pool.Acquire()
    if err != nil {
        return err
    }
    defer ch.Close()

    return ch.PublishWithContext(context.Background(), "jobs", job.Type, false, false,
        amqp.Publishing{Body: job.Payload, DeliveryMode: amqp.Persistent},
    )
}
```

---

## 8. RabbitMQ vs Kafka

| Dimension | RabbitMQ | Kafka |
|-----------|----------|-------|
| **Delivery model** | Push — broker pushes messages to consumers | Pull — consumers poll the broker |
| **Message retention** | Deleted after successful ack | Retained for a configurable period (default 7 days) |
| **Replay** | Not supported (message is gone after ack) | Supported — reset consumer group offset |
| **Routing** | Rich: direct, fanout, topic, headers | Simple: topic + partition key |
| **Ordering** | Per-queue FIFO | Per-partition ordering |
| **Throughput** | ~50K msgs/sec per node | Millions of msgs/sec |
| **Use case** | Task queues, work distribution, complex routing | Event streaming, audit logs, large-scale pipelines |
| **Consumer groups** | Competing consumers on one queue | Consumer groups — each group reads all messages independently |
| **Best for** | "Do this work exactly once" | "Record that this happened, let anyone replay" |

### Rule of thumb

- Use **RabbitMQ** when you need flexible routing, task queues (job workers), or request-reply patterns.
- Use **Kafka** when you need event streaming, audit logs, or multiple independent services consuming the same stream.

---

## Summary

- **AMQP flow**: Producer → Exchange → (Binding) → Queue → Consumer
- **Exchange types**: Direct (exact key), Fanout (broadcast), Topic (wildcard `*`/`#`), Headers (attribute match)
- **Persistent messages**: set `DeliveryMode: amqp.Persistent` and declare durable queues to survive restarts
- **Manual ack**: call `d.Ack(false)` after successful processing; `d.Nack(false, true)` to requeue; `d.Nack(false, false)` to send to DLQ
- **DLX**: declare a queue with `x-dead-letter-exchange` argument — failed/expired messages route there automatically
- **Publisher confirms**: `ch.Confirm(false)` + wait on `NotifyPublish` channel — ensures the broker persisted the message
- **Channels are not goroutine-safe**: use one channel per goroutine; pool the connection
- **RabbitMQ** suits task queues and routing; **Kafka** suits event streaming and replay

## Exercises

### Easy
1. Start RabbitMQ with Docker (`docker run -p 5672:5672 -p 15672:15672 rabbitmq:3-management`). Write a producer that publishes 10 `OrderEvent` messages to a direct exchange `orders` with routing key `payment`. Write a consumer that reads from a queue bound to that exchange and prints each event.
2. Change the exchange to `fanout` and bind two queues: `notifications` and `audit-log`. Verify both queues receive every published message.
3. Declare a queue with `DeliveryMode: amqp.Transient` on the producer and `durable: false` on the queue. Restart RabbitMQ and observe the messages disappear. Then switch both to persistent/durable and verify messages survive a restart.

### Medium
4. Implement a **topic exchange** with routing key `logs.<level>.<service>`. Create three queues: `errors-only` (binding `logs.error.#`), `payment-logs` (binding `logs.*.payment-svc`), and `all-logs` (binding `logs.#`). Publish events at different levels and services; verify each queue receives the right subset.
5. Set up a **Dead Letter Exchange** for a `payment-tasks` queue. Write a consumer that calls `d.Nack(false, false)` on every third message. Verify those messages appear in the `payment-tasks.dlq` queue via the RabbitMQ management UI.
6. Implement **publisher confirms**: publish 100 messages and collect the confirmation channel; report how many were acked vs nacked. Introduce a simulated broker failure (close the channel mid-way) and verify your code returns an error instead of silently dropping messages.

### Hard
7. Build a **reliable worker pool**: use `ch.Qos(1, 0, false)` so each worker gets one message at a time. Spawn 5 goroutines, each with its own channel. Track processing time per message. If a worker takes more than 2 seconds, nack with requeue and log a timeout. Show that work distributes evenly across workers.
8. Implement a **retry queue with exponential backoff** using two queues: a main queue (`orders.work`) and a delay queue (`orders.retry`) with `x-message-ttl` set to 5 seconds and `x-dead-letter-exchange` pointing back to `orders.work`. When a job fails, publish to `orders.retry` with a header tracking the attempt count; after 3 attempts, send to `orders.dlq` permanently.
