# Chapter 101: Watermill — The Go Event-Driven Framework

Watermill is an abstraction layer over message brokers. Instead of writing Kafka-specific or RabbitMQ-specific code, you write against Watermill's Publisher/Subscriber interfaces. Swapping brokers means changing two lines of config, not rewriting your business logic.

## Table of Contents

1. [What Is Watermill?](#1-what-is-watermill)
2. [Core Concepts](#2-core-concepts)
3. [Publisher and Subscriber](#3-publisher-and-subscriber)
4. [The Router](#4-the-router)
5. [Middleware](#5-middleware)
6. [Switching Backends](#6-switching-backends)
7. [CQRS with Watermill](#7-cqrs-with-watermill)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is Watermill?

```
Your code        Watermill layer         Broker
─────────        ───────────────         ──────
Handler ──────→  Router                  Kafka
                 Publisher interface ──→ RabbitMQ (AMQP)
                 Subscriber interface ←─ NATS JetStream
                 Middleware             Redis Streams
                 CQRS helpers          GoChannel (in-memory, for tests)
```

Without Watermill, your business logic is tangled with Kafka consumer loops, AMQP channel management, and manual retry logic. With Watermill, you write a handler function that receives a `*message.Message` and returns `([]*message.Message, error)`. The framework handles the rest.

Install:

```
go get github.com/ThreeDotsLabs/watermill
go get github.com/ThreeDotsLabs/watermill-kafka/v2        # Kafka backend
go get github.com/ThreeDotsLabs/watermill-amqp/v2         # RabbitMQ backend
go get github.com/ThreeDotsLabs/watermill-nats/v2          # NATS backend
```

---

## 2. Core Concepts

| Concept | Description |
|---------|-------------|
| **Message** | UUID + `[]byte` payload + metadata map. The unit of work. |
| **Publisher** | Sends messages to a topic |
| **Subscriber** | Provides a channel of messages from a topic |
| **Router** | Connects subscribers to handlers; runs middleware; manages goroutines |
| **Handler** | A function: `func(msg *message.Message) ([]*message.Message, error)` |
| **Middleware** | Wraps handlers to add retry, logging, metrics, poison queue, etc. |
| **Topic** | Broker-specific name — a Kafka topic, a RabbitMQ queue, a NATS subject |

### The Message type

```go
import "github.com/ThreeDotsLabs/watermill/message"

// Create a message
uuid := watermill.NewUUID()  // generates a UUID string
msg := message.NewMessage(uuid, payload) // payload is []byte

// Metadata is a key-value map attached to the message (like HTTP headers)
msg.Metadata.Set("correlation_id", correlationID)
msg.Metadata.Set("event_type", "order.created")

// Access metadata
corrID := msg.Metadata.Get("correlation_id")

// Payload
fmt.Println(string(msg.Payload))
```

---

## 3. Publisher and Subscriber

### In-memory (GoChannel) — good for tests and local dev

```go
package watermill_example

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/ThreeDotsLabs/watermill"
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type OrderEvent struct {
    EventType string    `json:"event_type"`
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    OccuredAt time.Time `json:"occurred_at"`
}

func basicPubSub() {
    logger := watermill.NewStdLogger(false, false)

    // GoChannel is an in-memory pub/sub — useful for testing
    pubSub := gochannel.NewGoChannel(
        gochannel.Config{
            OutputChannelBuffer: 100,
            Persistent:          false, // messages lost if no subscriber
        },
        logger,
    )

    // Subscribe
    messages, err := pubSub.Subscribe(context.Background(), "orders.created")
    if err != nil {
        panic(err)
    }

    // Process in a goroutine
    go func() {
        for msg := range messages {
            var event OrderEvent
            json.Unmarshal(msg.Payload, &event)
            fmt.Printf("received order %s for $%.2f\n", event.OrderID, event.Amount)

            // Ack signals Watermill that processing succeeded
            msg.Ack()
            // Use msg.Nack() to signal failure — Watermill will redeliver (if backend supports it)
        }
    }()

    // Publish
    event := OrderEvent{
        EventType: "order.created",
        OrderID:   "ord-001",
        UserID:    "usr-42",
        Amount:    149.99,
        OccuredAt: time.Now(),
    }
    payload, _ := json.Marshal(event)
    msg := message.NewMessage(watermill.NewUUID(), payload)
    msg.Metadata.Set("event_type", event.EventType)

    if err := pubSub.Publish("orders.created", msg); err != nil {
        panic(err)
    }

    time.Sleep(100 * time.Millisecond)
}
```

---

## 4. The Router

The Router is the heart of Watermill. It connects subscriber topics to handler functions, applies middleware, and manages the goroutine lifecycle. Think of it as an HTTP router but for messages.

```go
package watermill_example

import (
    "context"
    "encoding/json"

    "github.com/ThreeDotsLabs/watermill"
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/ThreeDotsLabs/watermill/message/router/middleware"
    "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func buildRouter(pubSub *gochannel.GoChannel) (*message.Router, error) {
    logger := watermill.NewStdLogger(false, false)

    router, err := message.NewRouter(message.RouterConfig{}, logger)
    if err != nil {
        return nil, err
    }

    // AddHandler: consume from inTopic, transform, publish result to outTopic
    // Signature: name, inTopic, subscriber, outTopic, publisher, handlerFunc
    router.AddHandler(
        "order-enricher",         // handler name (unique within router)
        "orders.created",         // input topic
        pubSub,                   // subscriber
        "orders.enriched",        // output topic (publish return value here)
        pubSub,                   // publisher
        enrichOrderHandler,       // the handler function
    )

    // AddNoPublisherHandler: consume only, no output topic
    // Use for side-effect handlers (send email, write to DB, etc.)
    router.AddNoPublisherHandler(
        "notification-sender",
        "orders.enriched",
        pubSub,
        sendNotificationHandler,
    )

    return router, nil
}

// Handler that transforms a message and returns output messages
func enrichOrderHandler(msg *message.Message) ([]*message.Message, error) {
    var event OrderEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        return nil, err
    }

    // Enrich: fetch user details, add shipping estimate, etc.
    enriched := struct {
        OrderEvent
        CustomerName string `json:"customer_name"`
        ShippingDays int    `json:"shipping_days"`
    }{
        OrderEvent:   event,
        CustomerName: fetchCustomerName(event.UserID),
        ShippingDays: estimateShipping(event.OrderID),
    }

    payload, err := json.Marshal(enriched)
    if err != nil {
        return nil, err
    }

    outMsg := message.NewMessage(watermill.NewUUID(), payload)
    outMsg.Metadata.Set("correlation_id", msg.Metadata.Get("correlation_id"))
    return []*message.Message{outMsg}, nil
}

// No-publisher handler: pure side effect
func sendNotificationHandler(msg *message.Message) error {
    var event OrderEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        return err
    }
    return sendEmail(event.UserID, fmt.Sprintf("Your order %s is being processed", event.OrderID))
}

func runRouter(router *message.Router) {
    ctx := context.Background()
    if err := router.Run(ctx); err != nil {
        panic(err)
    }
}
```

### Wiring multiple handlers into a pipeline

```
[orders.created] → enrichOrderHandler → [orders.enriched]
                                               │
                           ┌───────────────────┤
                           ▼                   ▼
               sendNotificationHandler   auditLogHandler
               (no publisher)            (no publisher)
```

---

## 5. Middleware

Middleware wraps handlers. Watermill ships with a useful standard library of middleware in `github.com/ThreeDotsLabs/watermill/message/router/middleware`.

```go
package watermill_example

import (
    "time"

    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/ThreeDotsLabs/watermill/message/router/middleware"
    "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func buildRouterWithMiddleware(pubSub *gochannel.GoChannel, poisonPub message.Publisher) (*message.Router, error) {
    logger := watermill.NewStdLogger(false, false)
    router, _ := message.NewRouter(message.RouterConfig{}, logger)

    // ---- Global middleware (applies to all handlers) ----

    // Recoverer: catch panics in handlers, convert to errors, don't crash the router
    router.AddMiddleware(middleware.Recoverer)

    // CorrelationID: propagate correlation_id from incoming message to outgoing messages
    router.AddMiddleware(middleware.CorrelationID)

    // Throttle: limit to 100 messages per second across all handlers
    router.AddMiddleware(middleware.NewThrottle(100, time.Second).Middleware)

    // ---- Per-handler middleware ----

    router.AddHandler(
        "order-processor",
        "orders.created",
        pubSub,
        "orders.processed",
        pubSub,
        processOrderHandler,
    ).AddMiddleware(
        // Retry: retry up to 3 times with exponential backoff before giving up
        middleware.Retry{
            MaxRetries:      3,
            InitialInterval: 100 * time.Millisecond,
            Multiplier:      2,          // 100ms → 200ms → 400ms
            MaxInterval:     5 * time.Second,
            Logger:          logger,
        }.Middleware,

        // Poison: if the message still fails after all retries, send to poison queue
        // instead of blocking the handler forever
        middleware.NewPoisonQueue(poisonPub, "orders.poison").Middleware,
    )

    return router, nil
}

func processOrderHandler(msg *message.Message) ([]*message.Message, error) {
    var event OrderEvent
    if err := json.Unmarshal(msg.Payload, &event); err != nil {
        // Returning an error causes Retry middleware to retry
        return nil, fmt.Errorf("unmarshal: %w", err)
    }

    if err := reserveInventory(event); err != nil {
        return nil, fmt.Errorf("reserve inventory: %w", err)
    }

    result := message.NewMessage(watermill.NewUUID(), msg.Payload)
    return []*message.Message{result}, nil
}
```

### Middleware execution order

```
Incoming message
       │
       ▼
  CorrelationID   ← outermost: runs first on the way in, last on the way out
       │
       ▼
  Retry           ← wraps the inner call, catches errors and retries
       │
       ▼
  Poison          ← last resort: if Retry gives up, sends to poison queue
       │
       ▼
  Your Handler    ← innermost: your actual business logic
```

---

## 6. Switching Backends

This is Watermill's superpower. The handler code never changes; only the subscriber/publisher construction changes.

### Kafka backend

```go
import (
    "github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
    "github.com/Shopify/sarama"
)

func newKafkaPublisher(brokers []string) (message.Publisher, error) {
    saramaConfig := kafka.DefaultSaramaSyncPublisherConfig()
    return kafka.NewPublisher(
        kafka.PublisherConfig{
            Brokers:   brokers,
            Marshaler: kafka.DefaultMarshaler{},
        },
        watermill.NewStdLogger(false, false),
    )
}

func newKafkaSubscriber(brokers []string, consumerGroup string) (message.Subscriber, error) {
    saramaConfig := kafka.DefaultSaramaSubscriberConfig()
    saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
    return kafka.NewSubscriber(
        kafka.SubscriberConfig{
            Brokers:               brokers,
            Unmarshaler:           kafka.DefaultMarshaler{},
            ConsumerGroup:         consumerGroup,
            OverwriteSaramaConfig: saramaConfig,
        },
        watermill.NewStdLogger(false, false),
    )
}
```

### AMQP (RabbitMQ) backend

```go
import "github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"

func newAMQPPublisher(amqpURL string) (message.Publisher, error) {
    config := amqp.NewDurablePubSubConfig(amqpURL, nil)
    return amqp.NewPublisher(config, watermill.NewStdLogger(false, false))
}

func newAMQPSubscriber(amqpURL string) (message.Subscriber, error) {
    config := amqp.NewDurablePubSubConfig(amqpURL, amqp.GenerateQueueNameTopicNameWithSuffix("svc"))
    return amqp.NewSubscriber(config, watermill.NewStdLogger(false, false))
}
```

### Wiring handlers: only the constructor changes

```go
func buildApp(cfg Config) (*message.Router, error) {
    logger := watermill.NewStdLogger(false, false)

    var pub message.Publisher
    var sub message.Subscriber
    var err error

    switch cfg.Backend {
    case "kafka":
        pub, err = newKafkaPublisher(cfg.KafkaBrokers)
        sub, err = newKafkaSubscriber(cfg.KafkaBrokers, "order-svc")
    case "amqp":
        pub, err = newAMQPPublisher(cfg.AMQPURL)
        sub, err = newAMQPSubscriber(cfg.AMQPURL)
    default:
        // In-memory for local dev and tests
        goChannel := gochannel.NewGoChannel(gochannel.Config{}, logger)
        pub, sub = goChannel, goChannel
    }
    if err != nil {
        return nil, err
    }

    router, _ := message.NewRouter(message.RouterConfig{}, logger)
    router.AddMiddleware(middleware.Recoverer, middleware.CorrelationID)

    // These handlers are identical regardless of backend:
    router.AddHandler("order-enricher", "orders.created", sub, "orders.enriched", pub, enrichOrderHandler)
    router.AddNoPublisherHandler("notification-sender", "orders.enriched", sub, sendNotificationHandler)

    return router, nil
}
```

---

## 7. CQRS with Watermill

Watermill ships a `cqrs` package that builds Commands and Events buses on top of its Router. This lets you build a proper CQRS/Event Sourcing architecture without reinventing the plumbing.

```go
package cqrs_example

import (
    "context"

    "github.com/ThreeDotsLabs/watermill"
    "github.com/ThreeDotsLabs/watermill/components/cqrs"
    "github.com/ThreeDotsLabs/watermill/message"
    "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// ---- Commands ----

type PlaceOrderCommand struct {
    OrderID string  `json:"order_id"`
    UserID  string  `json:"user_id"`
    Amount  float64 `json:"amount"`
}

type PlaceOrderHandler struct {
    eventBus *cqrs.EventBus
}

func (h *PlaceOrderHandler) HandlerName() string { return "PlaceOrderHandler" }
func (h *PlaceOrderHandler) NewCommand() interface{} { return &PlaceOrderCommand{} }

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd interface{}) error {
    c := cmd.(*PlaceOrderCommand)

    // Business logic: validate, persist to DB, etc.
    if err := persistOrder(c); err != nil {
        return err
    }

    // Emit a domain event after the command succeeds
    return h.eventBus.Publish(ctx, &OrderPlacedEvent{
        OrderID: c.OrderID,
        UserID:  c.UserID,
        Amount:  c.Amount,
    })
}

// ---- Events ----

type OrderPlacedEvent struct {
    OrderID string  `json:"order_id"`
    UserID  string  `json:"user_id"`
    Amount  float64 `json:"amount"`
}

type SendConfirmationEmailHandler struct{}

func (h *SendConfirmationEmailHandler) HandlerName() string  { return "SendConfirmationEmail" }
func (h *SendConfirmationEmailHandler) NewEvent() interface{} { return &OrderPlacedEvent{} }

func (h *SendConfirmationEmailHandler) Handle(ctx context.Context, event interface{}) error {
    e := event.(*OrderPlacedEvent)
    return sendEmail(e.UserID, fmt.Sprintf("Order %s confirmed!", e.OrderID))
}

// ---- Wiring ----

func buildCQRS() error {
    logger := watermill.NewStdLogger(false, false)
    pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

    router, _ := message.NewRouter(message.RouterConfig{}, logger)
    marshaler := cqrs.JSONMarshaler{}

    // EventBus: publishes events to topics named after the event type
    eventBus, err := cqrs.NewEventBusWithConfig(pubSub, cqrs.EventBusConfig{
        GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
            return params.EventName, nil // topic name = event struct name
        },
        Marshaler: marshaler,
        Logger:    logger,
    })
    if err != nil {
        return err
    }

    // CommandBus: routes commands to their handlers
    commandBus, err := cqrs.NewCommandBusWithConfig(pubSub, cqrs.CommandBusConfig{
        GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
            return params.CommandName, nil
        },
        Marshaler: marshaler,
        Logger:    logger,
    })
    if err != nil {
        return err
    }

    // CommandProcessor: subscribes to command topics and dispatches to handlers
    commandProcessor, _ := cqrs.NewCommandProcessorWithConfig(router, cqrs.CommandProcessorConfig{
        GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
            return params.CommandName, nil
        },
        SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
            return pubSub, nil
        },
        Marshaler: marshaler,
        Logger:    logger,
    })
    commandProcessor.AddHandlers(&PlaceOrderHandler{eventBus: eventBus})

    // EventProcessor: subscribes to event topics and dispatches to handlers
    eventProcessor, _ := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
        GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
            return params.EventName, nil
        },
        SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
            return pubSub, nil
        },
        Marshaler: marshaler,
        Logger:    logger,
    })
    eventProcessor.AddHandlers(&SendConfirmationEmailHandler{})

    ctx := context.Background()
    go router.Run(ctx)
    <-router.Running() // wait until router is ready

    // Send a command — flows through CommandBus → PlaceOrderHandler → EventBus → SendConfirmationEmailHandler
    return commandBus.Send(ctx, &PlaceOrderCommand{
        OrderID: "ord-001",
        UserID:  "usr-42",
        Amount:  149.99,
    })
}
```

### CQRS flow

```
commandBus.Send(PlaceOrderCommand)
           │
           ▼
  [topic: PlaceOrderCommand]
           │
           ▼
  PlaceOrderHandler.Handle()
      ├── persistOrder()
      └── eventBus.Publish(OrderPlacedEvent)
                   │
                   ▼
          [topic: OrderPlacedEvent]
                   │
           ┌───────┴──────────────┐
           ▼                      ▼
  SendConfirmationEmail   UpdateInventoryHandler
  Handler.Handle()        Handler.Handle()
```

---

## Summary

- **Watermill** provides `Publisher` and `Subscriber` interfaces — swap brokers by changing the constructor, not the handler code
- **Message**: UUID + `[]byte` payload + metadata map; `msg.Ack()` signals success, `msg.Nack()` signals failure
- **Router**: the main runtime component — connects subscribers to handlers, manages goroutines, applies middleware
- **AddHandler**: input topic → handler → output topic (transforms a message); **AddNoPublisherHandler**: pure consumer (side effects only)
- **Middleware stacks**: `Recoverer` (catch panics), `Retry` (exponential backoff), `Poison` (DLQ after retries exhausted), `CorrelationID` (trace propagation), `Throttle` (rate limiting)
- **Backends**: GoChannel (in-memory/tests), Kafka, AMQP (RabbitMQ), NATS, Redis Streams — same handler code for all
- **CQRS package**: `cqrs.CommandBus` + `cqrs.EventBus` built on the same Router — enables CQRS/ES with minimal boilerplate
- Watermill is a good default for Go services that need broker-agnostic event-driven architecture

## Exercises

### Easy
1. Use GoChannel as the backend. Build a two-handler pipeline: `order-validator` reads from `orders.raw`, validates that `Amount > 0`, and publishes valid orders to `orders.valid`. A `order-logger` NoPublisherHandler reads from `orders.valid` and logs each order. Publish 10 orders (2 with Amount = 0) and verify only 8 reach the logger.
2. Add `middleware.Recoverer` and a handler that panics on every other message. Verify the router continues running and processes the remaining messages instead of crashing.
3. Add `middleware.CorrelationID` and log the correlation ID in two handlers. Publish a message and verify the same correlation ID appears in both handler logs.

### Medium
4. Add `middleware.Retry{MaxRetries: 3, InitialInterval: 50 * time.Millisecond, Multiplier: 2}` to a handler. Make the handler fail the first 2 times (use a counter in a closure), succeed on the 3rd. Verify the message is eventually processed and that the retry intervals grow exponentially.
5. Add `middleware.NewPoisonQueue(poisonPub, "orders.poison")` and make a handler always return an error. After MaxRetries are exhausted, verify the message appears in `orders.poison`. Write a poison queue consumer that logs the failed message and the reason.
6. Swap the backend from GoChannel to a real Kafka broker running locally. Change only the `newPublisher` and `newSubscriber` constructor calls; keep all handler code identical. Verify the pipeline works the same way.

### Hard
7. Build a full **CQRS order system** using Watermill's `cqrs` package: `PlaceOrderCommand` → `PlaceOrderHandler` (persist + emit `OrderPlacedEvent`), `OrderPlacedEvent` → two event handlers in parallel: `SendEmailHandler` and `ChargePaymentHandler`. Use GoChannel so no broker is needed. Write a table-driven test that sends 5 commands and verifies both handlers fired for each.
8. Implement a **saga orchestrator** using Watermill: when `OrderPlacedEvent` is received, the saga sends `ReserveInventoryCommand` and `ChargePaymentCommand` in sequence. If `ChargePaymentCommand` fails, the saga sends a `ReleaseInventoryCommand` (compensation). Model each step as a Watermill handler. Use a state machine pattern (pending → inventory-reserved → payment-charged → completed / rolling-back → compensated) stored in-memory, keyed by order ID.
