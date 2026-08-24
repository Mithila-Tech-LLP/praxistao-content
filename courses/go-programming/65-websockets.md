# Chapter 44: WebSockets

HTTP is request-response: the client asks, the server answers, the connection closes. WebSockets flip this model — the connection stays open and either side can send messages at any time. This enables real-time features: chat, live dashboards, collaborative editing, and push notifications.

## Table of Contents

1. [WebSocket Protocol](#1-websocket-protocol)
2. [Server with gorilla/websocket](#2-server-with-gorillawebsocket)
3. [Connection Hub — Broadcasting](#3-connection-hub--broadcasting)
4. [Rooms and Namespaces](#4-rooms-and-namespaces)
5. [Presence and State](#5-presence-and-state)
6. [Authentication over WebSockets](#6-authentication-over-websockets)
7. [Testing WebSockets](#7-testing-websockets)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. WebSocket Protocol

**The handshake**: WebSocket starts as an HTTP request with `Upgrade: websocket` header. The server responds with `101 Switching Protocols`, and from that point the TCP connection is used for the WebSocket protocol instead of HTTP.

```
Client → Server:
GET /ws HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13

Server → Client:
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

**Message types**: text (UTF-8 string) or binary (raw bytes). Control frames: ping, pong, close.

```bash
go get github.com/gorilla/websocket
```

---

## 2. Server with gorilla/websocket

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // In production: check r.Header.Get("Origin") against allowlist
        return true
    },
}

// Client wraps a WebSocket connection with message queuing.
type Client struct {
    conn   *websocket.Conn
    send   chan []byte
    userID int64
}

const (
    writeWait      = 10 * time.Second    // Time to write message to client
    pongWait       = 60 * time.Second    // Time to wait for pong response
    pingPeriod     = (pongWait * 9) / 10 // Ping interval (< pongWait)
    maxMessageSize = 512 * 1024          // 512 KB
)

// readPump pumps messages from the WebSocket connection to the hub.
// Each connection runs one readPump in its own goroutine.
func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err,
                websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        hub.broadcast <- message
    }
}

// writePump pumps messages from the send channel to the WebSocket connection.
// Each connection runs one writePump in its own goroutine.
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Hub closed the channel:
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil { return }
            w.Write(message)

            // Batch any queued messages:
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte("\n"))
                w.Write(<-c.send)
            }
            if err := w.Close(); err != nil { return }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    client := &Client{conn: conn, send: make(chan []byte, 256)}
    hub.register <- client

    // Start goroutines:
    go client.writePump()
    go client.readPump(hub)
}
```

---

## 3. Connection Hub — Broadcasting

The hub is the central goroutine that manages all connections and routes messages. **All state mutations happen in a single goroutine** — no mutexes needed.

```go
// Hub maintains the set of active clients and broadcasts to them.
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("client connected, total: %d", len(h.clients))

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("client disconnected, total: %d", len(h.clients))
            }

        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    // Client's buffer is full — assume disconnected:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

**The full server:**
```go
func main() {
    hub := NewHub()
    go hub.Run()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWS(hub, w, r)
    })
    http.HandleFunc("/", serveHome)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 4. Rooms and Namespaces

```go
// Room-based hub for a chat application:
type Message struct {
    Type    string `json:"type"`     // "message", "join", "leave"
    Room    string `json:"room"`
    Content string `json:"content"`
    From    string `json:"from"`
    Time    string `json:"time"`
}

type RoomHub struct {
    rooms      map[string]map[*Client]bool  // room → set of clients
    register   chan *RoomJoin
    unregister chan *Client
    broadcast  chan *RoomMessage
    mu         sync.RWMutex  // Protecting rooms map
}

type RoomJoin struct {
    client *Client
    room   string
}

type RoomMessage struct {
    room    string
    message []byte
    from    *Client  // Exclude sender from broadcast (optional)
}

func NewRoomHub() *RoomHub {
    return &RoomHub{
        rooms:      make(map[string]map[*Client]bool),
        register:   make(chan *RoomJoin),
        unregister: make(chan *Client),
        broadcast:  make(chan *RoomMessage, 256),
    }
}

func (h *RoomHub) Run() {
    for {
        select {
        case join := <-h.register:
            if h.rooms[join.room] == nil {
                h.rooms[join.room] = make(map[*Client]bool)
            }
            h.rooms[join.room][join.client] = true

        case client := <-h.unregister:
            for room, clients := range h.rooms {
                if _, ok := clients[client]; ok {
                    delete(clients, client)
                    close(client.send)
                    if len(clients) == 0 { delete(h.rooms, room) }
                }
            }

        case msg := <-h.broadcast:
            for client := range h.rooms[msg.room] {
                if client == msg.from { continue }  // Don't send back to sender
                select {
                case client.send <- msg.message:
                default:
                    close(client.send)
                    delete(h.rooms[msg.room], client)
                }
            }
        }
    }
}

// Client with room support:
func (c *Client) joinRoom(hub *RoomHub, room string) {
    hub.register <- &RoomJoin{client: c, room: room}
}

func (c *Client) sendToRoom(hub *RoomHub, room string, msg []byte) {
    hub.broadcast <- &RoomMessage{room: room, message: msg, from: c}
}
```

---

## 5. Presence and State

```go
// Typing indicators and online presence:

type PresenceEvent struct {
    Type   string `json:"type"`   // "online", "offline", "typing", "stopped_typing"
    UserID int64  `json:"userId"`
    Room   string `json:"room,omitempty"`
}

type PresenceHub struct {
    *RoomHub
    online map[int64]*Client  // userID → client
}

func (h *PresenceHub) OnConnect(client *Client) {
    h.online[client.userID] = client
    event, _ := json.Marshal(PresenceEvent{Type: "online", UserID: client.userID})
    h.broadcast <- &RoomMessage{room: "global", message: event}
}

func (h *PresenceHub) OnDisconnect(client *Client) {
    delete(h.online, client.userID)
    event, _ := json.Marshal(PresenceEvent{Type: "offline", UserID: client.userID})
    h.broadcast <- &RoomMessage{room: "global", message: event}
}

func (h *PresenceHub) IsOnline(userID int64) bool {
    _, ok := h.online[userID]
    return ok
}

// Throttle typing events to avoid flooding:
type TypingDebouncer struct {
    timers map[string]*time.Timer
    mu     sync.Mutex
}

func (d *TypingDebouncer) Debounce(key string, delay time.Duration, fn func()) {
    d.mu.Lock()
    defer d.mu.Unlock()
    if t, ok := d.timers[key]; ok { t.Stop() }
    d.timers[key] = time.AfterFunc(delay, fn)
}
```

---

## 6. Authentication over WebSockets

```go
// Option 1: JWT in query parameter (simple, but token in logs):
func serveWSAuth(hub *Hub, jwtSvc *JWTService, w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    claims, err := jwtSvc.ValidateToken(token)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil { return }

    client := &Client{conn: conn, send: make(chan []byte, 256), userID: claims.UserID}
    hub.register <- client
    go client.writePump()
    go client.readPump(hub)
}

// Option 2: Token in first message after connection (cleaner):
func (c *Client) authenticate(jwtSvc *JWTService) (*Claims, error) {
    c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    _, msg, err := c.conn.ReadMessage()
    if err != nil { return nil, err }

    var authMsg struct { Token string `json:"token"` }
    if err := json.Unmarshal(msg, &authMsg); err != nil {
        return nil, errors.New("expected auth message")
    }
    return jwtSvc.ValidateToken(authMsg.Token)
}

// Option 3: Validate HTTP-only cookie during upgrade (most secure):
func serveWSCookieAuth(hub *Hub, sessionStore SessionStore, w http.ResponseWriter, r *http.Request) {
    session, err := sessionStore.Get(r, "session")
    if err != nil || session.Values["userID"] == nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    // Proceed with upgrade...
}
```

---

## 7. Testing WebSockets

```go
import "github.com/gorilla/websocket"

func TestWebSocketBroadcast(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        serveWS(hub, w, r)
    }))
    defer server.Close()

    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

    // Connect two clients:
    conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    require.NoError(t, err)
    defer conn1.Close()

    conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    require.NoError(t, err)
    defer conn2.Close()

    // Wait for registration:
    time.Sleep(50 * time.Millisecond)

    // Send from conn1:
    msg := []byte(`{"type":"message","content":"hello"}`)
    err = conn1.WriteMessage(websocket.TextMessage, msg)
    require.NoError(t, err)

    // Both clients should receive it (broadcast):
    conn2.SetReadDeadline(time.Now().Add(time.Second))
    _, received, err := conn2.ReadMessage()
    require.NoError(t, err)
    assert.Equal(t, msg, received)

    conn1.SetReadDeadline(time.Now().Add(time.Second))
    _, received, err = conn1.ReadMessage()
    require.NoError(t, err)
    assert.Equal(t, msg, received)
}

func TestWebSocketPingPong(t *testing.T) {
    // Verify connection stays alive with ping/pong:
    server := httptest.NewServer(http.HandlerFunc(wsHandler))
    defer server.Close()

    conn, _, _ := websocket.DefaultDialer.Dial(
        "ws" + strings.TrimPrefix(server.URL, "http"), nil)
    defer conn.Close()

    // Server should send a ping within pingPeriod:
    conn.SetPingHandler(func(appData string) error {
        return conn.WriteControl(websocket.PongMessage, []byte(appData),
            time.Now().Add(time.Second))
    })

    // Stay alive for 2 pings:
    conn.SetReadDeadline(time.Now().Add(2 * pingPeriod + time.Second))
    _, _, err := conn.ReadMessage()
    // Should not get an error — ping/pong kept the connection alive
    if !websocket.IsCloseError(err, websocket.CloseNormalClosure) && err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

---

## Summary

- WebSocket = persistent bidirectional connection via HTTP upgrade (`101 Switching Protocols`)
- Each client: two goroutines — `readPump` (reads from WS, sends to hub) and `writePump` (reads from channel, writes to WS)
- Hub: single goroutine managing all clients — no mutexes on the client map
- Always set `ReadDeadline` refreshed by pong handler; send pings on a timer
- Broadcast: send to `client.send` channel; if full, assume disconnected and close
- Authentication: validate JWT in query param, first message, or cookie during upgrade
- Testing: `httptest.NewServer` + `websocket.DefaultDialer.Dial`

---

## Exercises

### Easy
1. Build a simple echo WebSocket server that sends back exactly what it receives, prefixed with `"Echo: "`. Connect with a browser or `wscat` and verify bidirectional communication.
2. Add a `GET /ws/stats` HTTP endpoint (not WebSocket) that returns the current number of connected clients: `{"connected": N}`. The hub should expose a `Count() int` method that's safe to call from HTTP handlers.
3. Write a WebSocket client in Go that connects to the echo server, sends 10 messages, reads 10 responses, and verifies each response equals `"Echo: " + original`. Use goroutines for concurrent send/receive.

### Medium
4. Build a **live feed**: server pushes a new random stock price update every second to all connected clients. Message format: `{"symbol": "AAPL", "price": 182.34, "change": +0.52}`. The client should print updates until the connection is closed (Ctrl+C sends close frame). Test with two browser tabs open simultaneously.
5. Build a **collaborative cursor tracker**: clients join a room and broadcast their mouse position (`{"x": 100, "y": 200}`). Each client's position is relayed to all OTHER clients in the room (not back to the sender). Maintain a server-side presence map of `userID → {x, y}`. New connections receive a snapshot of all current cursors on join.
6. Implement **message persistence**: when a client connects to a room, replay the last 50 messages from an in-memory ring buffer before forwarding live messages. `RingBuffer[T]` with size 50. Accessing the buffer is safe because it's only touched from the hub's goroutine.

### Hard
7. **Horizontal scaling with Redis Pub/Sub**: when running multiple server instances, a message from a client on server A must reach clients on server B. Implement a `RedisBus` that subscribes to a Redis channel per room, publishes broadcast messages to Redis, and receives them to forward to local clients. Test with two server instances. (Preview of Chapter 56 — Redis.)
8. **Graceful reconnection protocol**: design and implement a reconnection protocol: (1) each client gets a session ID on first connect, (2) messages are assigned sequence numbers per room, (3) on reconnect with the same session ID, replay any missed messages (up to a max gap). The server stores the last 1000 messages per room in memory. The client sends its last received sequence number on reconnect.
