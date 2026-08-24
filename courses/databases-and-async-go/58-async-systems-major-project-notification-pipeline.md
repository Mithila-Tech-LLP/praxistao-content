# Chapter 58: Async Systems — Major Project: Notification Pipeline

Build a complete, production-grade notification system: events flow from multiple services through Kafka, get processed by different consumers, and deliver notifications via email, SMS (simulated), and push notifications — with retries, dead letter queues, and monitoring.

## Table of Contents

1. System Design
2. Event Producer Service
3. Notification Router (Kafka Consumer)
4. Email Delivery Worker
5. SMS Delivery Worker
6. Monitoring Dashboard
7. Running the Full System
8. Exercises

---

## 1. System Design

```
┌──────────────┐     ┌─────────────────────────────────────┐
│ Order Service│────►│   Kafka: "user-events" (3 partitions)│
│ Auth Service │────►└──────────────┬──────────────────────┘
│ User Service │                    │
└──────────────┘                    │
                                    ▼
                     ┌─────────────────────────────────────┐
                     │  Notification Router Consumer       │
                     │  - Routes events to notification    │
                     │    type queues                      │
                     └────────────┬────────────────────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              ▼                   ▼                   ▼
      ┌───────────┐       ┌───────────┐       ┌───────────┐
      │  Kafka:   │       │  Kafka:   │       │  Kafka:   │
      │  email    │       │   sms     │       │   push    │
      │  queue    │       │   queue   │       │   queue   │
      └─────┬─────┘       └─────┬─────┘       └─────┬─────┘
            ▼                   ▼                   ▼
      ┌───────────┐       ┌───────────┐       ┌───────────┐
      │   Email   │       │   SMS     │       │   Push    │
      │  Worker   │       │  Worker   │       │  Worker   │
      └───────────┘       └───────────┘       └───────────┘
            │                   │
            ▼                   ▼
      ┌───────────────────────────┐
      │ DLQ: failed-notifications │
      └───────────────────────────┘
```

---

## 2. Event Producer Service

```go
// producer/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    kafka "github.com/segmentio/kafka-go"
)

type UserEvent struct {
    EventID   string    `json:"event_id"`
    Type      string    `json:"type"`
    UserID    string    `json:"user_id"`
    UserEmail string    `json:"user_email"`
    UserPhone string    `json:"user_phone"`
    Data      map[string]interface{} `json:"data"`
    OccurredAt time.Time `json:"occurred_at"`
}

var writer *kafka.Writer

func main() {
    writer = &kafka.Writer{
        Addr:         kafka.TCP("localhost:9092"),
        Topic:        "user-events",
        Balancer:     &kafka.Hash{},
        RequiredAcks: kafka.RequireOne,
    }
    defer writer.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("POST /events/order-placed", handleOrderPlaced)
    mux.HandleFunc("POST /events/user-signup", handleUserSignup)
    mux.HandleFunc("POST /events/password-reset", handlePasswordReset)
    mux.HandleFunc("POST /simulate", handleSimulate)

    log.Println("Event Producer API on :8081")
    log.Fatal(http.ListenAndServe(":8081", mux))
}

func handleOrderPlaced(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID  string  `json:"user_id"`
        Email   string  `json:"email"`
        OrderID string  `json:"order_id"`
        Amount  float64 `json:"amount"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    event := UserEvent{
        EventID:    fmt.Sprintf("evt-%d", time.Now().UnixNano()),
        Type:       "order.placed",
        UserID:     req.UserID,
        UserEmail:  req.Email,
        OccurredAt: time.Now(),
        Data: map[string]interface{}{
            "order_id": req.OrderID,
            "amount":   req.Amount,
        },
    }
    publishEvent(event)
    w.WriteHeader(202)
    json.NewEncoder(w).Encode(map[string]string{"event_id": event.EventID})
}

func handleUserSignup(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID string `json:"user_id"`
        Email  string `json:"email"`
        Name   string `json:"name"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    event := UserEvent{
        EventID:    fmt.Sprintf("evt-%d", time.Now().UnixNano()),
        Type:       "user.signup",
        UserID:     req.UserID,
        UserEmail:  req.Email,
        OccurredAt: time.Now(),
        Data:       map[string]interface{}{"name": req.Name},
    }
    publishEvent(event)
    w.WriteHeader(202)
}

func handlePasswordReset(w http.ResponseWriter, r *http.Request) {
    var req struct {
        UserID string `json:"user_id"`
        Email  string `json:"email"`
        Token  string `json:"token"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    event := UserEvent{
        EventID:    fmt.Sprintf("evt-%d", time.Now().UnixNano()),
        Type:       "user.password_reset",
        UserID:     req.UserID,
        UserEmail:  req.Email,
        OccurredAt: time.Now(),
        Data:       map[string]interface{}{"token": req.Token},
    }
    publishEvent(event)
    w.WriteHeader(202)
}

func handleSimulate(w http.ResponseWriter, r *http.Request) {
    count := 50
    types := []string{"order.placed", "user.signup", "user.password_reset"}
    for i := 0; i < count; i++ {
        event := UserEvent{
            EventID:    fmt.Sprintf("evt-%d", time.Now().UnixNano()+int64(i)),
            Type:       types[rand.Intn(len(types))],
            UserID:     fmt.Sprintf("user-%d", rand.Intn(1000)),
            UserEmail:  fmt.Sprintf("user%d@example.com", rand.Intn(1000)),
            OccurredAt: time.Now(),
            Data:       map[string]interface{}{"simulated": true},
        }
        publishEvent(event)
    }
    json.NewEncoder(w).Encode(map[string]int{"published": count})
}

func publishEvent(event UserEvent) {
    data, _ := json.Marshal(event)
    if err := writer.WriteMessages(context.Background(), kafka.Message{
        Key:   []byte(event.UserID),
        Value: data,
    }); err != nil {
        log.Printf("publish error: %v", err)
    }
}
```

---

## 3. Notification Router

```go
// router/main.go — routes user-events to typed notification queues
package main

import (
    "context"
    "encoding/json"
    "log"

    kafka "github.com/segmentio/kafka-go"
)

type NotificationMessage struct {
    EventID       string                 `json:"event_id"`
    UserID        string                 `json:"user_id"`
    Email         string                 `json:"email"`
    Phone         string                 `json:"phone"`
    Template      string                 `json:"template"`
    TemplateData  map[string]interface{} `json:"template_data"`
    DeliveryChannel string               `json:"delivery_channel"` // "email", "sms", "push"
}

var producers = map[string]*kafka.Writer{
    "email": {Addr: kafka.TCP("localhost:9092"), Topic: "notifications.email"},
    "sms":   {Addr: kafka.TCP("localhost:9092"), Topic: "notifications.sms"},
    "push":  {Addr: kafka.TCP("localhost:9092"), Topic: "notifications.push"},
}

// Routing rules: event type → delivery channels + template
var routes = map[string][]struct {
    Channel  string
    Template string
}{
    "order.placed": {
        {"email", "order_confirmation"},
        {"push", "order_push"},
    },
    "user.signup": {
        {"email", "welcome"},
        {"sms", "welcome_sms"},
    },
    "user.password_reset": {
        {"email", "password_reset"},
    },
}

func main() {
    ctx := context.Background()

    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:    []string{"localhost:9092"},
        Topic:      "user-events",
        GroupID:    "notification-router",
        MaxWait:    500,
    })
    defer reader.Close()

    log.Println("Notification Router started")

    for {
        msg, err := reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            log.Println("fetch:", err)
            continue
        }

        var event UserEvent
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            log.Printf("bad event: %v", err)
            reader.CommitMessages(ctx, msg)
            continue
        }

        routedAny := false
        for _, route := range routes[event.Type] {
            notif := NotificationMessage{
                EventID:         event.EventID,
                UserID:          event.UserID,
                Email:           event.UserEmail,
                Phone:           event.UserPhone,
                Template:        route.Template,
                TemplateData:    event.Data,
                DeliveryChannel: route.Channel,
            }
            data, _ := json.Marshal(notif)
            producers[route.Channel].WriteMessages(ctx, kafka.Message{
                Key:   []byte(event.UserID),
                Value: data,
            })
            log.Printf("[Router] %s → %s/%s", event.Type, route.Channel, route.Template)
            routedAny = true
        }

        if !routedAny {
            log.Printf("[Router] No route for event type: %s", event.Type)
        }

        reader.CommitMessages(ctx, msg)
    }
}
```

---

## 4. Email Delivery Worker

```go
// workers/email/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    kafka "github.com/segmentio/kafka-go"
)

// Email templates
var templates = map[string]string{
    "order_confirmation": "Thank you for your order! Order ID: {{order_id}}, Amount: ${{amount}}",
    "welcome":            "Welcome to our platform, {{name}}! Your account is ready.",
    "password_reset":     "Reset your password using this token: {{token}} (expires in 1 hour)",
}

var (
    dlqWriter *kafka.Writer
    delivered int
    failed    int
)

func main() {
    dlqWriter = &kafka.Writer{
        Addr:  kafka.TCP("localhost:9092"),
        Topic: "notifications.dead-letter",
    }
    defer dlqWriter.Close()

    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "notifications.email",
        GroupID: "email-worker",
        MaxWait: 500,
    })
    defer reader.Close()

    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer stop()

    log.Println("Email Worker started")

    // Stats printer
    go func() {
        for {
            time.Sleep(5 * time.Second)
            log.Printf("[Email] Stats: delivered=%d, failed=%d", delivered, failed)
        }
    }()

    for {
        msg, err := reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                log.Printf("[Email] Shutting down. Delivered: %d, Failed: %d", delivered, failed)
                return
            }
            continue
        }

        var notif NotificationMessage
        if err := json.Unmarshal(msg.Value, &notif); err != nil {
            sendToDLQ(ctx, msg, "invalid JSON: "+err.Error())
            reader.CommitMessages(ctx, msg)
            continue
        }

        if err := sendEmail(notif); err != nil {
            sendToDLQ(ctx, msg, err.Error())
            failed++
        } else {
            delivered++
        }

        reader.CommitMessages(ctx, msg)
    }
}

func sendEmail(notif NotificationMessage) error {
    template, ok := templates[notif.Template]
    if !ok {
        return &permanentError{fmt.Errorf("unknown template: %s", notif.Template)}
    }

    // Render template (simplified)
    body := template
    for k, v := range notif.TemplateData {
        body = strings.ReplaceAll(body, "{{"+k+"}}", fmt.Sprintf("%v", v))
    }

    // In production: use SMTP or SendGrid API
    log.Printf("[Email] To: %s | Template: %s | Body: %s",
        notif.Email, notif.Template, body)

    // Simulate occasional failures
    if rand.Float64() < 0.05 { // 5% failure rate
        return fmt.Errorf("SMTP timeout")
    }
    return nil
}

type permanentError struct{ err error }
func (e *permanentError) Error() string { return e.err.Error() }

func sendToDLQ(ctx context.Context, msg kafka.Message, reason string) {
    dlqMsg := map[string]interface{}{
        "original_topic": msg.Topic,
        "payload":        string(msg.Value),
        "reason":         reason,
        "failed_at":      time.Now(),
    }
    data, _ := json.Marshal(dlqMsg)
    dlqWriter.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: data})
    log.Printf("[Email] Sent to DLQ: %s", reason)
}
```

---

## 5. Running the Full System

```bash
# Start infrastructure
docker compose up -d

# Create Kafka topics
for topic in user-events notifications.email notifications.sms notifications.push notifications.dead-letter; do
  docker exec kafka /opt/kafka/bin/kafka-topics.sh \
    --bootstrap-server localhost:9092 \
    --create --topic $topic \
    --partitions 3 --replication-factor 1
done

# Start services (in separate terminals)
cd producer && go run main.go    # :8081
cd router && go run main.go
cd workers/email && go run main.go
cd workers/sms && go run main.go

# Test: simulate 50 events
curl -X POST localhost:8081/simulate

# Check DLQ
docker exec kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic notifications.dead-letter \
  --from-beginning
```

---

## Summary

You've built a production-grade notification pipeline with:
- Event producers publishing to Kafka
- A routing consumer that dispatches to specialized notification queues
- Per-channel delivery workers with failure handling
- Dead letter queue for failed messages
- Stats monitoring

This is the pattern behind every major notification system (Twilio, SendGrid, Firebase, Mailchimp).

### Exercises

**Easy:** Add a monitoring endpoint to the email worker: `GET /stats` that returns `{delivered, failed, dlq_count}` as JSON.

**Medium:** Add retry logic to the email worker: if `sendEmail` returns a transient error (not `permanentError`), retry up to 3 times with 2s, 4s, 8s backoff before sending to DLQ.

**Hard:** Add end-to-end tracing using OpenTelemetry. Propagate a trace ID from the event producer through the router and into the email worker. Log the full trace path for each notification delivered.
