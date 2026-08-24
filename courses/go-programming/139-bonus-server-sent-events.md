# Chapter 48: Server-Sent Events and Real-Time Patterns

Not every real-time feature needs WebSockets. Server-Sent Events (SSE) is a one-way push channel from server to browser — simpler, HTTP-native, and auto-reconnecting. Long polling is even simpler. Knowing which tool fits which problem is the skill.

## Table of Contents

1. [Comparing Real-Time Approaches](#1-comparing-real-time-approaches)
2. [Server-Sent Events](#2-server-sent-events)
3. [Long Polling](#3-long-polling)
4. [Fan-out Architecture](#4-fan-out-architecture)
5. [Backpressure and Slow Consumers](#5-backpressure-and-slow-consumers)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Comparing Real-Time Approaches

| | HTTP Polling | Long Polling | SSE | WebSocket |
|--|-------------|--------------|-----|-----------|
| Direction | Client→Server | Client→Server | Server→Client | Bidirectional |
| Protocol | HTTP | HTTP | HTTP | WS (upgrade) |
| Browser support | All | All | All (not IE) | All |
| Reconnect | Manual | Manual | Automatic | Manual |
| Complexity | Low | Medium | Low | High |
| Overhead | High (repeat requests) | Medium | Low | Low |
| Use cases | Status polling | Notifications | Feeds, live data | Chat, collaboration |

**Decision guide:**
- Need bidirectional real-time? → WebSocket
- Server pushes to browser, browser only sends via REST? → SSE
- Simple polling with low frequency? → Long polling or polling
- Notification push from backend? → SSE or WebPush

---

## 2. Server-Sent Events

SSE is an HTTP response that never ends. The server keeps writing `data:` lines; the browser's `EventSource` API reads them automatically and reconnects if the connection drops.

**Wire format:**
```
HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

data: {"type":"update","value":42}\n\n

event: user-joined\n
data: {"name":"Alice"}\n\n

id: 123\n
data: ping\n\n
```

### SSE Server
```go
package sse

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Event represents a single SSE message.
type Event struct {
    ID    string
    Type  string  // Empty = "message" (default)
    Data  any
    Retry int  // Reconnect delay in ms (0 = don't set)
}

func (e Event) Write(w http.ResponseWriter) error {
    if e.ID != "" {
        fmt.Fprintf(w, "id: %s\n", e.ID)
    }
    if e.Type != "" {
        fmt.Fprintf(w, "event: %s\n", e.Type)
    }
    if e.Retry > 0 {
        fmt.Fprintf(w, "retry: %d\n", e.Retry)
    }

    var dataStr string
    switch v := e.Data.(type) {
    case string:
        dataStr = v
    case []byte:
        dataStr = string(v)
    default:
        b, err := json.Marshal(v)
        if err != nil { return err }
        dataStr = string(b)
    }
    fmt.Fprintf(w, "data: %s\n\n", dataStr)  // Double newline ends the event
    return nil
}

// Broker manages SSE connections and broadcasts events.
type Broker struct {
    clients    map[chan Event]bool
    register   chan chan Event
    unregister chan chan Event
    broadcast  chan Event
    mu         sync.RWMutex
}

func NewBroker() *Broker {
    b := &Broker{
        clients:    make(map[chan Event]bool),
        register:   make(chan chan Event),
        unregister: make(chan chan Event),
        broadcast:  make(chan Event, 256),
    }
    go b.run()
    return b
}

func (b *Broker) run() {
    for {
        select {
        case ch := <-b.register:
            b.mu.Lock()
            b.clients[ch] = true
            b.mu.Unlock()

        case ch := <-b.unregister:
            b.mu.Lock()
            if _, ok := b.clients[ch]; ok {
                delete(b.clients, ch)
                close(ch)
            }
            b.mu.Unlock()

        case event := <-b.broadcast:
            b.mu.RLock()
            for ch := range b.clients {
                select {
                case ch <- event:
                default:
                    // Slow client — skip this event
                }
            }
            b.mu.RUnlock()
        }
    }
}

func (b *Broker) Publish(event Event) {
    b.broadcast <- event
}

func (b *Broker) ClientCount() int {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return len(b.clients)
}

// Handler serves the SSE stream.
func (b *Broker) Handler(w http.ResponseWriter, r *http.Request) {
    // Verify the client supports SSE:
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming not supported", http.StatusInternalServerError)
        return
    }

    // SSE headers:
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no")  // Disable Nginx buffering

    // Register this client:
    eventCh := make(chan Event, 16)
    b.register <- eventCh
    defer func() { b.unregister <- eventCh }()

    // Send initial "connected" event:
    Event{Type: "connected", Data: map[string]int{"clients": b.ClientCount()}}.Write(w)
    flusher.Flush()

    // Set retry hint — browser reconnects after 3s if connection drops:
    Event{Retry: 3000}.Write(w)

    // Stream events until client disconnects:
    for {
        select {
        case event, ok := <-eventCh:
            if !ok { return }
            event.Write(w)
            flusher.Flush()

        case <-r.Context().Done():
            return  // Client disconnected

        case <-time.After(30 * time.Second):
            // Keep-alive comment (prevents proxy from closing idle connections):
            fmt.Fprint(w, ": keep-alive\n\n")
            flusher.Flush()
        }
    }
}
```

### SSE Client in Go (for testing)
```go
func consumeSSE(ctx context.Context, url string, handler func(Event)) error {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("Accept", "text/event-stream")
    req.Header.Set("Cache-Control", "no-cache")

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()

    scanner := bufio.NewScanner(resp.Body)
    var event Event
    var dataLines []string

    for scanner.Scan() {
        line := scanner.Text()
        if line == "" {
            // Empty line = end of event
            if len(dataLines) > 0 {
                event.Data = strings.Join(dataLines, "\n")
                handler(event)
                event = Event{}
                dataLines = nil
            }
            continue
        }
        if strings.HasPrefix(line, "data:") {
            dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
        } else if strings.HasPrefix(line, "event:") {
            event.Type = strings.TrimPrefix(line, "event: ")
        } else if strings.HasPrefix(line, "id:") {
            event.ID = strings.TrimPrefix(line, "id: ")
        }
    }
    return scanner.Err()
}
```

### JavaScript Client (reference)
```javascript
// Browser EventSource — auto-reconnects on disconnect:
const es = new EventSource('/events');

es.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log('Update:', data);
};

es.addEventListener('user-joined', (e) => {
    console.log('User joined:', JSON.parse(e.data));
});

es.onerror = (e) => console.error('SSE error:', e);

// Send last event ID on reconnect — server uses this to replay missed events:
// Browser automatically sends "Last-Event-ID" header when reconnecting
```

---

## 3. Long Polling

Long polling holds the request open until there's something to return.

```go
type LongPollBroker struct {
    waiters map[string]chan any  // topic → waiting channels
    mu      sync.Mutex
}

func NewLongPollBroker() *LongPollBroker {
    return &LongPollBroker{waiters: make(map[string]chan any)}
}

// Wait returns the next event for a topic, blocking until timeout or event.
func (b *LongPollBroker) Wait(ctx context.Context, topic string, timeout time.Duration) (any, error) {
    ch := make(chan any, 1)

    b.mu.Lock()
    b.waiters[topic] = ch
    b.mu.Unlock()

    defer func() {
        b.mu.Lock()
        delete(b.waiters, topic)
        b.mu.Unlock()
    }()

    timer := time.NewTimer(timeout)
    defer timer.Stop()

    select {
    case event := <-ch:
        return event, nil
    case <-timer.C:
        return nil, nil  // Timeout — client should retry
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Publish sends an event to a waiting client.
func (b *LongPollBroker) Publish(topic string, event any) {
    b.mu.Lock()
    ch, ok := b.waiters[topic]
    b.mu.Unlock()
    if ok {
        select {
        case ch <- event:
        default:
        }
    }
}

// Long-poll handler:
func (b *LongPollBroker) Handler(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    if topic == "" {
        http.Error(w, "topic required", http.StatusBadRequest)
        return
    }

    event, err := b.Wait(r.Context(), topic, 30*time.Second)
    if err != nil {
        http.Error(w, "context cancelled", http.StatusServiceUnavailable)
        return
    }
    if event == nil {
        // Timeout — return empty 204 so client retries:
        w.WriteHeader(http.StatusNoContent)
        return
    }
    writeJSON(w, http.StatusOK, event)
}
```

---

## 4. Fan-out Architecture

When events come from multiple sources (database changes, queue messages, user actions), a fan-out layer distributes them to all connected SSE/WS clients.

```go
// EventBus: an in-process pub/sub system.
type EventBus struct {
    subscribers map[string][]chan Event  // topic → subscriber channels
    mu          sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{subscribers: make(map[string][]chan Event)}
}

func (b *EventBus) Subscribe(topic string, bufSize int) (<-chan Event, func()) {
    ch := make(chan Event, bufSize)
    b.mu.Lock()
    b.subscribers[topic] = append(b.subscribers[topic], ch)
    b.mu.Unlock()

    unsubscribe := func() {
        b.mu.Lock()
        subs := b.subscribers[topic]
        for i, s := range subs {
            if s == ch {
                b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
                break
            }
        }
        b.mu.Unlock()
        close(ch)
    }
    return ch, unsubscribe
}

func (b *EventBus) Publish(topic string, event Event) {
    b.mu.RLock()
    subs := b.subscribers[topic]
    b.mu.RUnlock()

    for _, ch := range subs {
        select {
        case ch <- event:
        default:  // Non-blocking: skip slow subscribers
        }
    }
}

// Wiring SSE with EventBus:
func sseHandler(bus *EventBus, w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")

    flusher := w.(http.Flusher)
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")

    // Subscribe to user-specific and global events:
    userEvents, unsubUser := bus.Subscribe("user:"+userID, 16)
    globalEvents, unsubGlobal := bus.Subscribe("global", 16)
    defer unsubUser()
    defer unsubGlobal()

    for {
        select {
        case event := <-userEvents:
            event.Write(w); flusher.Flush()
        case event := <-globalEvents:
            event.Write(w); flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}
```

---

## 5. Backpressure and Slow Consumers

A slow consumer that can't keep up with the event rate can cause memory to grow without bound.

```go
// Strategy 1: Non-blocking send with drop (good for non-critical events)
select {
case ch <- event:
default:
    droppedEvents.Add(1)  // Metric
}

// Strategy 2: Drop oldest (ring buffer per client)
type RingBuffer[T any] struct {
    buf  []T
    head int
    size int
    mu   sync.Mutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
    return &RingBuffer[T]{buf: make([]T, capacity)}
}

func (r *RingBuffer[T]) Push(v T) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.buf[r.head%len(r.buf)] = v
    r.head++
    if r.size < len(r.buf) { r.size++ }
}

func (r *RingBuffer[T]) Drain() []T {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.size == 0 { return nil }
    result := make([]T, r.size)
    start := (r.head - r.size + len(r.buf)) % len(r.buf)
    for i := 0; i < r.size; i++ {
        result[i] = r.buf[(start+i)%len(r.buf)]
    }
    r.size = 0
    r.head = 0
    return result
}

// Strategy 3: Disconnect slow clients after buffer fills
func (b *Broker) sendOrDisconnect(client *Client, event Event) {
    select {
    case client.ch <- event:
    case <-time.After(100 * time.Millisecond):
        b.unregister <- client  // Disconnect laggy client
    }
}
```

---

## Summary

- **SSE**: one-way server→client push over HTTP; auto-reconnect; simpler than WebSocket for live feeds
- SSE format: `data: ...\n\n` per event; `event:` names custom events; `id:` enables last-event replay
- Always set `Content-Type: text/event-stream`, `Cache-Control: no-cache`, and periodic keep-alive comments
- Use `http.Flusher` to flush each event immediately (don't wait for buffer to fill)
- **Long polling**: hold request open until data ready or timeout; client retries on `204` or timeout
- **EventBus**: decouple event producers from SSE/WS connections; subscribe per topic
- **Backpressure**: non-blocking send + drop, ring buffer, or disconnect slow consumers

---

## Exercises

### Easy
1. Build a live counter: `GET /events` streams SSE events. A background goroutine increments a counter every second and publishes `{"count": N}` to the broker. Connect with `curl -N http://localhost:8080/events` and watch the stream.
2. Add `Last-Event-ID` support to the SSE handler: read `r.Header.Get("Last-Event-ID")`, look up missed events from a ring buffer of the last 100 events, and replay them before starting the live stream.
3. Build a long-poll notification service: `POST /notify/{topic}` publishes a message; `GET /wait/{topic}` blocks until a message arrives (max 30s). Test with two curl windows: one waiting and one posting.

### Medium
4. Build a **live order book**: stock orders arrive as events (buy/sell at a price). The order book maintains bids (buy orders, highest price first) and asks (sell orders, lowest price first). On each change, publish the current top 5 bids and asks via SSE. Multiple browser clients all see the same state.
5. **SSE with auth**: add JWT auth to the SSE handler. The token is passed as a query parameter (`/events?token=...`) since `EventSource` in browsers can't set headers. Validate the token before upgrading to SSE. If the token expires mid-stream, close the connection with a `401` event: `event: error\ndata: {"code":401}\n\n`.
6. Implement **topic-based SSE fan-out** where each user subscribes to their own channel (`/events?userID=42`) and also a global channel. Publish a user-specific event and verify only that user's SSE stream receives it; publish a global event and verify all streams receive it.

### Hard
7. **Distributed SSE with Redis Pub/Sub**: when multiple server instances run, a client connected to server A misses events published to server B. Implement `RedisBroker`: on `Publish`, write to a Redis channel; each server subscribes to Redis and forwards to its local SSE clients. Handle reconnects: store last 1000 events per topic in a Redis list and replay on reconnect.
8. **SSE load test**: write a Go load tester that opens 10,000 SSE connections to your server and measures: time-to-first-event, events-per-second throughput, memory per connection, and tail latency (p99). Identify the bottleneck (goroutines? channels? GC?) using `pprof`. Optimize until you can sustain 10,000 connections at 10 events/sec on a single server.
