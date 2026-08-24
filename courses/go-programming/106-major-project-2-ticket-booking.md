# Chapter 95: Major Project 2 — Ticket Booking System

This project synthesizes Vol 8 (Async & Event-Driven). You'll build a concert ticket booking system where inventory management, payment processing, and ticket delivery are all asynchronous and event-driven.

**What makes this hard**: when 10,000 users try to book the last 100 tickets simultaneously, you need to avoid overselling while remaining fast. The solution: Redis atomic operations for inventory, asynq for background processing, Kafka for event streaming.

---

## System Design

```
Client → HTTP API → Redis inventory check (atomic) → Create booking (pending)
                  → Enqueue payment job (asynq)
                                          ↓
                                   Payment Worker
                                   ├── Charge card
                                   ├── On success: confirm booking, publish BookingConfirmed event
                                   └── On failure: release inventory, publish BookingFailed event
                                                              ↓
                                                     Kafka Consumer
                                                     ├── Send confirmation email
                                                     ├── Generate ticket PDF
                                                     └── Update analytics
```

---

## Domain Model

```go
// domain/booking.go
package domain

type BookingStatus string
const (
    BookingPending   BookingStatus = "pending"
    BookingConfirmed BookingStatus = "confirmed"
    BookingFailed    BookingStatus = "failed"
    BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
    ID         string
    EventID    string
    UserID     string
    SeatIDs    []string
    TotalPrice float64
    Status     BookingStatus
    CreatedAt  time.Time
    UpdatedAt  time.Time
    FailReason string
}

type Event struct {
    ID           string
    Name         string
    Venue        string
    StartsAt     time.Time
    TotalSeats   int
    AvailableSeats int
    PricePerSeat float64
}

type Seat struct {
    ID      string
    EventID string
    Row     string
    Number  int
    Held    bool // reserved but not confirmed
    Sold    bool
}
```

---

## Inventory Management with Redis

The key challenge: prevent overselling without a slow database transaction.

```go
// infrastructure/redis/inventory.go
package redis

// Hold seats atomically using a Lua script
// Returns true if all seats were successfully held, false if any seat is taken
var holdSeatsScript = redis.NewScript(`
    local event_key = KEYS[1]
    local hold_ttl = ARGV[1]
    
    for i = 2, #ARGV do
        local seat_key = KEYS[i+1]
        -- seat_key exists if seat is already held or sold
        if redis.call("EXISTS", seat_key) == 1 then
            -- Roll back already-held seats
            for j = 2, i-1 do
                redis.call("DEL", KEYS[j+1])
            end
            return 0  -- seat taken
        end
        -- Hold this seat for 5 minutes
        redis.call("SET", seat_key, ARGV[i+1], "EX", hold_ttl)
    end
    return 1  -- all seats held
`)

type InventoryService struct {
    rdb *redis.Client
}

const holdTTL = 300 // 5 minutes in seconds

func (s *InventoryService) HoldSeats(ctx context.Context, eventID string, seatIDs []string, bookingID string) (bool, error) {
    keys := make([]string, 1+len(seatIDs))
    keys[0] = "event:" + eventID
    for i, seatID := range seatIDs {
        keys[i+1] = fmt.Sprintf("seat:%s:%s", eventID, seatID)
    }
    
    args := make([]interface{}, 1+len(seatIDs))
    args[0] = holdTTL
    for i := range seatIDs {
        args[i+1] = bookingID // value = who holds it
    }
    
    result, err := holdSeatsScript.Run(ctx, s.rdb, keys, args...).Int()
    if err != nil { return false, err }
    return result == 1, nil
}

func (s *InventoryService) ConfirmSeats(ctx context.Context, eventID string, seatIDs []string) error {
    pipe := s.rdb.Pipeline()
    for _, seatID := range seatIDs {
        key := fmt.Sprintf("seat:%s:%s", eventID, seatID)
        // Mark as permanently sold (no TTL)
        pipe.Set(ctx, key, "sold", 0)
    }
    _, err := pipe.Exec(ctx)
    return err
}

func (s *InventoryService) ReleaseSeats(ctx context.Context, eventID string, seatIDs []string) error {
    keys := make([]string, len(seatIDs))
    for i, seatID := range seatIDs {
        keys[i] = fmt.Sprintf("seat:%s:%s", eventID, seatID)
    }
    return s.rdb.Del(ctx, keys...).Err()
}
```

---

## Booking Use Case

```go
// usecase/booking_service.go
package usecase

type CreateBookingInput struct {
    EventID string
    UserID  string
    SeatIDs []string
}

type BookingService struct {
    bookings  domain.BookingRepository
    inventory *redis.InventoryService
    taskQueue *asynq.Client
    idgen     IDGenerator
}

func (s *BookingService) CreateBooking(ctx context.Context, in CreateBookingInput) (*domain.Booking, error) {
    bookingID := s.idgen.New()
    
    // 1. Atomically hold seats in Redis (fast, no DB)
    held, err := s.inventory.HoldSeats(ctx, in.EventID, in.SeatIDs, bookingID)
    if err != nil { return nil, fmt.Errorf("hold seats: %w", err) }
    if !held      { return nil, domain.ErrSeatsUnavailable }
    
    // 2. Create pending booking in DB
    booking := &domain.Booking{
        ID:        bookingID,
        EventID:   in.EventID,
        UserID:    in.UserID,
        SeatIDs:   in.SeatIDs,
        Status:    domain.BookingPending,
        CreatedAt: time.Now(),
    }
    if err := s.bookings.Create(ctx, booking); err != nil {
        s.inventory.ReleaseSeats(ctx, in.EventID, in.SeatIDs) // rollback hold
        return nil, err
    }
    
    // 3. Enqueue payment processing (async)
    task := asynq.NewTask("payment:process", mustMarshal(ProcessPaymentPayload{
        BookingID: bookingID,
        UserID:    in.UserID,
        Amount:    booking.TotalPrice,
    }))
    
    _, err = s.taskQueue.EnqueueContext(ctx, task,
        asynq.Queue("payments"),
        asynq.MaxRetry(3),
        asynq.Timeout(2*time.Minute),
    )
    if err != nil {
        // Can't enqueue payment — release seats and fail booking
        s.inventory.ReleaseSeats(ctx, in.EventID, in.SeatIDs)
        s.bookings.UpdateStatus(ctx, bookingID, domain.BookingFailed, "failed to enqueue payment")
        return nil, fmt.Errorf("enqueue payment: %w", err)
    }
    
    return booking, nil
}
```

---

## Payment Worker

```go
// worker/payment_handler.go
package worker

type ProcessPaymentPayload struct {
    BookingID string  `json:"booking_id"`
    UserID    string  `json:"user_id"`
    Amount    float64 `json:"amount"`
}

type PaymentHandler struct {
    bookings  domain.BookingRepository
    inventory *redis.InventoryService
    payments  PaymentGateway
    events    *kafka.Writer
}

func (h *PaymentHandler) Handle(ctx context.Context, t *asynq.Task) error {
    var p ProcessPaymentPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil { return err }
    
    booking, err := h.bookings.GetByID(ctx, p.BookingID)
    if err != nil { return err }
    if booking.Status != domain.BookingPending {
        return nil // already processed (idempotent)
    }
    
    // Attempt payment
    paymentResult, err := h.payments.Charge(ctx, PaymentRequest{
        UserID:    p.UserID,
        Amount:    p.Amount,
        Reference: p.BookingID,
    })
    
    if err != nil || !paymentResult.Success {
        reason := "payment failed"
        if err != nil { reason = err.Error() }
        
        // Release seat holds
        h.inventory.ReleaseSeats(ctx, booking.EventID, booking.SeatIDs)
        
        // Update booking status
        h.bookings.UpdateStatus(ctx, p.BookingID, domain.BookingFailed, reason)
        
        // Publish failure event
        h.publishEvent(ctx, "booking.failed", BookingFailedEvent{
            BookingID: p.BookingID,
            UserID:    p.UserID,
            Reason:    reason,
        })
        
        return nil // don't retry payment failures
    }
    
    // Confirm seats permanently
    if err := h.inventory.ConfirmSeats(ctx, booking.EventID, booking.SeatIDs); err != nil {
        return fmt.Errorf("confirm seats: %w", err) // retry
    }
    
    // Update booking to confirmed
    if err := h.bookings.UpdateStatus(ctx, p.BookingID, domain.BookingConfirmed, ""); err != nil {
        return fmt.Errorf("update booking: %w", err) // retry
    }
    
    // Publish success event
    h.publishEvent(ctx, "booking.confirmed", BookingConfirmedEvent{
        BookingID:  p.BookingID,
        UserID:     p.UserID,
        EventID:    booking.EventID,
        SeatIDs:    booking.SeatIDs,
        TotalPrice: p.Amount,
    })
    
    return nil
}

func (h *PaymentHandler) publishEvent(ctx context.Context, eventType string, payload any) {
    data, _ := json.Marshal(payload)
    h.events.WriteMessages(ctx, kafka.Message{
        Topic: "bookings",
        Key:   []byte(eventType),
        Value: data,
    })
}
```

---

## Downstream Consumers

```go
// worker/booking_consumers.go

// Email consumer: sends confirmation or failure email
type BookingEmailConsumer struct {
    reader *kafka.Reader
    mailer Mailer
}

func (c *BookingEmailConsumer) Run(ctx context.Context) {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil { return }
        
        eventType := string(msg.Key)
        switch eventType {
        case "booking.confirmed":
            var e BookingConfirmedEvent
            json.Unmarshal(msg.Value, &e)
            c.mailer.SendBookingConfirmation(ctx, e)
        case "booking.failed":
            var e BookingFailedEvent
            json.Unmarshal(msg.Value, &e)
            c.mailer.SendBookingFailure(ctx, e)
        }
        
        c.reader.CommitMessages(ctx, msg)
    }
}

// Ticket generator consumer: creates PDF tickets
type TicketGeneratorConsumer struct {
    reader  *kafka.Reader
    storage Storage
}

func (c *TicketGeneratorConsumer) Run(ctx context.Context) {
    for {
        msg, _ := c.reader.FetchMessage(ctx)
        if string(msg.Key) != "booking.confirmed" { 
            c.reader.CommitMessages(ctx, msg)
            continue
        }
        
        var e BookingConfirmedEvent
        json.Unmarshal(msg.Value, &e)
        
        // Generate PDF ticket
        pdf := generateTicketPDF(e)
        url, _ := c.storage.Upload(ctx, fmt.Sprintf("tickets/%s.pdf", e.BookingID), pdf)
        
        slog.Info("ticket generated", "booking_id", e.BookingID, "url", url)
        c.reader.CommitMessages(ctx, msg)
    }
}
```

---

## HTTP API

```
POST /bookings            → CreateBooking (returns pending booking immediately)
GET  /bookings/:id        → GetBooking (client polls for status update)
DELETE /bookings/:id      → CancelBooking

GET  /events              → ListEvents
GET  /events/:id          → GetEvent
GET  /events/:id/seats    → GetSeatMap (available/held/sold)

WebSocket /bookings/:id/status  → real-time status updates
```

---

## What You Built

1. **Atomic inventory**: Redis Lua script holds seats without race conditions
2. **Async payment**: asynq worker processes payments in the background; HTTP returns immediately
3. **Event streaming**: Kafka fan-out notifies email service, ticket generator, and analytics independently
4. **Idempotent processing**: payment handler checks booking status before processing
5. **Rollback on failure**: seat holds released if payment fails

---

## Extension Challenges

1. **WebSocket status updates**: when booking status changes (pending → confirmed/failed), push update to the waiting client. Use Redis Pub/Sub: payment worker publishes, WebSocket handler subscribes.

2. **Waitlist**: when all seats are sold, allow users to join a waitlist. When a booking is cancelled, release the seat and auto-book the first waitlisted user via an asynq task.

3. **Surge pricing**: during high demand (>80% sold), increase ticket price by 20%. Check demand level on each booking attempt using a Redis counter.

4. **Analytics projection**: consume booking events from Kafka and maintain `event_stats(event_id, total_bookings, total_revenue, conversion_rate)`. Expose as `GET /events/:id/analytics`.
