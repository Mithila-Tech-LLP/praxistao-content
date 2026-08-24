# Chapter 53: RabbitMQ — Flexible Message Routing

If Kafka is a highway (fast, simple, one direction), RabbitMQ is a post office (flexible routing, acknowledgements, different delivery strategies). RabbitMQ is perfect when you need complex routing logic, per-message acknowledgement, or classic work queues.

## Table of Contents

1. RabbitMQ vs Kafka — Key Differences
2. Core Concepts: Exchanges, Queues, Bindings
3. Exchange Types
4. Docker Setup
5. Building with RabbitMQ in Go
6. Acknowledgements and Reliability
7. Dead Letter Queues
8. Mini Project: Task Processing System
9. Exercises

---

## 1. RabbitMQ vs Kafka — Key Differences

| Aspect | Kafka | RabbitMQ |
|--------|-------|----------|
| Model | Pull (consumers fetch) | Push (broker delivers) |
| Storage | Disk-first, long retention | RAM-first, short retention |
| Routing | By partition key | Complex routing rules |
| Consumer tracking | Consumer tracks offset | Broker tracks per-message |
| Order guarantee | Per partition | Per queue |
| Best for | Event streaming, replay | Task queues, RPC, complex routing |
| Throughput | Millions/sec | 50,000-100,000/sec |

**RabbitMQ shines when:**
- You need complex message routing (route by content, pattern matching)
- Tasks must be acknowledged per-message (not just by offset)
- You need priority queues, message TTL, or delayed delivery
- Request-reply (RPC) patterns

---

## 2. Core Concepts: Exchanges, Queues, Bindings

Unlike Kafka where producers write directly to topics, in RabbitMQ:

```
Producer → Exchange → (Binding rules) → Queue → Consumer
```

**Exchange:** The entry point for messages. Decides which queue(s) to route to.

**Queue:** Stores messages until a consumer processes them. Messages deleted after acknowledgement.

**Binding:** Rules connecting an exchange to a queue. Based on routing key, headers, etc.

```
Producer sends:  exchange="orders", routing_key="order.placed"

Exchange "orders" has bindings:
  routing_key "order.*"    → queue "order-processor"
  routing_key "order.placed" → queue "email-service"
  routing_key "#"           → queue "audit-log"

Result: message goes to all three queues matching the routing key!
```

---

## 3. Exchange Types

**Direct Exchange:** Route by exact routing key match.
```
exchange: "payments"
binding:  routing_key="payment.success" → queue "success-handler"
binding:  routing_key="payment.failed"  → queue "failure-handler"
```

**Fanout Exchange:** Broadcast to all bound queues, ignoring routing key.
```
exchange: "system-events" (fanout)
  → queue "email-service"
  → queue "analytics"
  → queue "notifications"
All three get every message.
```

**Topic Exchange:** Routing key with wildcards (`*` = one word, `#` = zero or more words).
```
binding: "order.placed.#"  → matches "order.placed", "order.placed.eu", "order.placed.us.west"
binding: "*.placed.*"      → matches "order.placed.eu" but NOT "order.placed"
binding: "#"               → matches everything
```

**Headers Exchange:** Route based on message headers, not routing key.
```
binding: headers{"format": "json", "x-match": "all"} → queue "json-processor"
binding: headers{"format": "xml", "x-match": "all"}  → queue "xml-processor"
```

---

## 4. Docker Setup

```bash
docker run -d \
  --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3-management

# Management UI: http://localhost:15672 (guest/guest)
```

```bash
go get github.com/rabbitmq/amqp091-go
```

---

## 5. Building with RabbitMQ in Go

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

func connect() (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, nil, fmt.Errorf("connect: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, nil, fmt.Errorf("channel: %w", err)
    }

    return conn, ch, nil
}

func declareWorkQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
    return ch.QueueDeclare(
        name,
        true,  // durable: survives broker restart
        false, // auto-delete: keep even when no consumers
        false, // exclusive: accessible by other connections
        false, // no-wait
        amqp.Table{
            "x-dead-letter-exchange": "dlx", // DLQ routing
        },
    )
}
```

**Producer:**

```go
func producer(ch *amqp.Channel, queue amqp.Queue, body string) error {
    return ch.PublishWithContext(
        context.Background(),
        "",         // exchange: "" = default direct exchange (routes to queue by name)
        queue.Name, // routing key = queue name for default exchange
        true,       // mandatory: fail if no route found
        false,      // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            DeliveryMode: amqp.Persistent, // survives broker restart
            Timestamp:    time.Now(),
            Body:         []byte(body),
        },
    )
}
```

**Consumer with manual acknowledgement:**

```go
func consumer(ch *amqp.Channel, queue amqp.Queue) error {
    // Set prefetch: don't deliver more than 5 messages at once per consumer
    ch.Qos(5, 0, false)

    msgs, err := ch.Consume(
        queue.Name,
        "",    // consumer tag: auto-generated
        false, // auto-ack: false = manual ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,
    )
    if err != nil {
        return err
    }

    for msg := range msgs {
        if err := processMessage(msg.Body); err != nil {
            log.Printf("process error: %v, rejecting message", err)
            // Reject and requeue (will be redelivered later)
            msg.Nack(false, true) // multiple=false, requeue=true
            continue
        }

        // Success: acknowledge the message
        msg.Ack(false) // multiple=false
    }
    return nil
}

func processMessage(body []byte) error {
    fmt.Printf("Processing: %s\n", body)
    time.Sleep(100 * time.Millisecond) // simulate work
    return nil
}
```

---

## 6. Acknowledgements and Reliability

RabbitMQ tracks individual message delivery:

```
Consumer connects → "give me messages"
Broker sends message
Consumer processes
Consumer sends ACK → broker deletes message

If consumer crashes before ACK:
Broker re-queues message → another consumer picks it up
```

**QoS (Quality of Service) / Prefetch:**

```go
// Only send up to 10 unacknowledged messages to this consumer
ch.Qos(10, 0, false)
```

Without QoS, RabbitMQ dumps all messages to the consumer at once. If the consumer is slow, it buffers millions of messages in memory and crashes. With QoS=10, the consumer gets at most 10 at a time.

**Publisher confirms (equivalent to Kafka acks):**

```go
// Enable publisher confirms
ch.Confirm(false)

confirmations := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

ch.Publish(...) // send the message

select {
case confirm := <-confirmations:
    if !confirm.Ack {
        log.Println("message not acknowledged by broker!")
    }
case <-time.After(5 * time.Second):
    log.Println("confirmation timeout!")
}
```

---

## 7. Dead Letter Queues

When a message can't be processed (wrong format, max retries exceeded), send it to a Dead Letter Queue (DLQ) for manual inspection.

```go
func setupDLQ(ch *amqp.Channel) error {
    // Declare the DLX (Dead Letter Exchange)
    ch.ExchangeDeclare("dlx", "fanout", true, false, false, false, nil)

    // Declare the DLQ
    ch.QueueDeclare("dead-letter-queue", true, false, false, false, nil)

    // Bind DLQ to DLX
    ch.QueueBind("dead-letter-queue", "", "dlx", false, nil)

    // Main queue with DLX routing
    ch.QueueDeclare("tasks", true, false, false, false, amqp.Table{
        "x-dead-letter-exchange": "dlx",  // failed messages go here
        "x-message-ttl":          30000,   // expire after 30s if not processed
    })

    return nil
}

// Consumer: send to DLQ after N failures
func consumerWithDLQ(ch *amqp.Channel) {
    msgs, _ := ch.Consume("tasks", "", false, false, false, false, nil)

    for msg := range msgs {
        retryCount, _ := msg.Headers["retry-count"].(int32)

        if err := processMessage(msg.Body); err != nil {
            if retryCount >= 3 {
                // Too many retries: send to DLQ
                log.Printf("message failed %d times, sending to DLQ", retryCount)
                msg.Nack(false, false) // requeue=false → goes to DLX/DLQ
            } else {
                // Increment retry count and requeue
                msg.Nack(false, true) // requeue=true
            }
            continue
        }

        msg.Ack(false)
    }
}
```

---

## 8. Mini Project: Task Processing System

A job queue: HTTP server accepts tasks, workers process them.

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "sync"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Task struct {
    ID       string `json:"id"`
    Type     string `json:"type"`
    Payload  string `json:"payload"`
    Priority int    `json:"priority"`
}

var (
    conn *amqp.Connection
    ch   *amqp.Channel
    mu   sync.Mutex
)

func main() {
    var err error
    conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err = conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    // Setup queues
    ch.ExchangeDeclare("tasks", "direct", true, false, false, false, nil)
    for _, queueName := range []string{"tasks.high", "tasks.normal", "tasks.low"} {
        q, _ := ch.QueueDeclare(queueName, true, false, false, false, nil)
        ch.QueueBind(q.Name, q.Name, "tasks", false, nil)
    }

    // Start 3 workers
    for i := 0; i < 3; i++ {
        go worker(i)
    }

    // HTTP server for submitting tasks
    http.HandleFunc("POST /tasks", handleSubmitTask)
    log.Println("Task API on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSubmitTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "invalid JSON", 400)
        return
    }
    task.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())

    // Route to queue based on priority
    queueName := "tasks.normal"
    if task.Priority > 7 {
        queueName = "tasks.high"
    } else if task.Priority < 3 {
        queueName = "tasks.low"
    }

    data, _ := json.Marshal(task)
    mu.Lock()
    err := ch.Publish("tasks", queueName, false, false, amqp.Publishing{
        ContentType:  "application/json",
        DeliveryMode: amqp.Persistent,
        Body:         data,
    })
    mu.Unlock()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.WriteHeader(202) // Accepted
    json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID, "queue": queueName})
}

func worker(id int) {
    wConn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    defer wConn.Close()
    wCh, _ := wConn.Channel()
    defer wCh.Close()

    wCh.Qos(5, 0, false)

    // Each worker reads from all priority queues
    for _, queueName := range []string{"tasks.high", "tasks.normal", "tasks.low"} {
        go func(q string) {
            msgs, _ := wCh.Consume(q, "", false, false, false, false, nil)
            for msg := range msgs {
                var task Task
                json.Unmarshal(msg.Body, &task)

                // Simulate work
                duration := time.Duration(100+rand.Intn(900)) * time.Millisecond
                fmt.Printf("[Worker %d] Processing %s (%s) — %v\n",
                    id, task.ID, task.Type, duration)
                time.Sleep(duration)

                msg.Ack(false)
            }
        }(queueName)
    }

    // Block forever
    select {}
}
```

Test:
```bash
# Submit tasks
curl -X POST localhost:8080/tasks \
  -d '{"type":"email","payload":"send welcome email","priority":8}'

curl -X POST localhost:8080/tasks \
  -d '{"type":"report","payload":"generate monthly report","priority":2}'
```

---

## Summary

- RabbitMQ uses exchanges + queues + bindings for flexible message routing.
- Four exchange types: direct (exact key), fanout (broadcast), topic (wildcard), headers.
- Always `ch.Qos(N, 0, false)` to limit unacknowledged messages per consumer.
- `msg.Ack(false)` after successful processing, `msg.Nack(false, requeue)` on failure.
- Dead Letter Queues catch messages that fail after N retries.

### Exercises

**Easy:** Set up a fanout exchange "events". Create 3 queues bound to it: "logging", "analytics", "alerts". Publish 5 messages and verify all 3 queues receive each message.

**Medium:** Implement a priority queue in RabbitMQ: queues `tasks.high`, `tasks.normal`, `tasks.low`. A single worker drains high-priority tasks first, then normal, then low. Verify ordering with 3 tasks of each priority.

**Hard:** Implement RPC (Remote Procedure Call) over RabbitMQ: client sends a request with a `reply_to` queue and `correlation_id`, server processes and replies to the reply_to queue, client waits for the response. Implement `Add(a, b int) int` as an RPC call.
