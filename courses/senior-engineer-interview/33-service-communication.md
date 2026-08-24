# Chapter 33: Service Communication — REST, gRPC, WebSockets & Event-Driven

Choosing the right communication protocol between services is a key architectural decision. Senior engineers must know when each protocol fits, how each works internally, and the trade-offs.

## Table of Contents

1. [REST — Principles and Best Practices](#1-rest--principles-and-best-practices)
2. [gRPC — High-Performance RPC](#2-grpc--high-performance-rpc)
3. [WebSockets — Full-Duplex Communication](#3-websockets--full-duplex-communication)
4. [Server-Sent Events (SSE)](#4-server-sent-events-sse)
5. [GraphQL — Flexible Querying](#5-graphql--flexible-querying)
6. [Event-Driven Communication](#6-event-driven-communication)
7. [Protocol Comparison Matrix](#7-protocol-comparison-matrix)
8. [Go Implementation Patterns](#8-go-implementation-patterns)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. REST — Principles and Best Practices

REST (Representational State Transfer) is the most common service communication style. It uses HTTP and URLs to represent resources and actions.

### REST Constraints

```
1. Uniform Interface: resources identified by URLs, operations via HTTP verbs
2. Stateless: each request contains all information needed (server holds no session)
3. Client-Server: separation of concerns
4. Cacheable: responses can be cached (use Cache-Control headers)
5. Layered System: proxies, CDNs, gateways are transparent to client
```

### URL Design

```
Resources are NOUNS, not verbs:
  Good:  GET /users/123
  Bad:   GET /getUser?id=123

  Good:  POST /orders
  Bad:   POST /createOrder

  Good:  DELETE /sessions/abc123   (logout)
  Bad:   POST /logout

Nested resources:
  GET  /users/123/orders          (list all orders for user 123)
  GET  /users/123/orders/456      (get specific order for user 123)
  POST /users/123/orders          (create order for user 123)

Filtering/sorting: use query params
  GET /orders?status=pending&sort=created_at&order=desc&page=2&limit=20
```

### HTTP Methods and Status Codes

```
GET:     read resource          → 200 OK, 404 Not Found
POST:    create resource        → 201 Created (Location header), 400 Bad Request
PUT:     replace resource       → 200 OK, 204 No Content
PATCH:   partial update         → 200 OK, 422 Unprocessable Entity
DELETE:  delete resource        → 204 No Content, 404 Not Found

Status codes:
  2xx Success:
    200 OK           — general success
    201 Created      — resource created (POST)
    204 No Content   — success, no body (DELETE, PUT)
  
  4xx Client Error:
    400 Bad Request          — invalid request body/params
    401 Unauthorized         — not authenticated
    403 Forbidden            — authenticated but not authorized
    404 Not Found            — resource doesn't exist
    409 Conflict             — duplicate key, version conflict
    422 Unprocessable Entity — valid format but business logic failed
    429 Too Many Requests    — rate limited
  
  5xx Server Error:
    500 Internal Server Error — unexpected server error
    502 Bad Gateway          — upstream service error
    503 Service Unavailable  — server overloaded or maintenance
    504 Gateway Timeout      — upstream timeout
```

### REST Best Practices

```go
// Version your API:
// v1: /api/v1/users
// v2: /api/v2/users (breaking changes go in v2)

// Return consistent error format:
type APIError struct {
    Code    string `json:"code"`    // machine-readable error code
    Message string `json:"message"` // human-readable message
    Details any    `json:"details,omitempty"`
}

// Use pagination for list endpoints:
type PagedResponse struct {
    Data     []any  `json:"data"`
    Total    int    `json:"total"`     // total count
    Page     int    `json:"page"`
    PageSize int    `json:"page_size"`
    NextPage *int   `json:"next_page,omitempty"`
}

// Use ETags for conditional requests (cache validation):
// Response: ETag: "abc123"
// Next request: If-None-Match: "abc123"
// If unchanged: 304 Not Modified (no body, bandwidth savings)
```

---

## 2. gRPC — High-Performance RPC

gRPC is a high-performance RPC framework that uses Protocol Buffers for serialization and HTTP/2 for transport. It's used for internal service-to-service communication.

### Protocol Buffers

```protobuf
// user.proto
syntax = "proto3";

package user;
option go_package = "github.com/example/api/user";

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);    // server streaming
    rpc BulkCreateUsers(stream CreateUserRequest) returns (BulkCreateResponse); // client streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage); // bidirectional streaming
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    google.protobuf.Timestamp created_at = 4;
}

message GetUserRequest {
    string id = 1;
}
```

### gRPC in Go

```go
// Server implementation:
type UserServer struct {
    pb.UnimplementedUserServiceServer
    db *sql.DB
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    var u pb.User
    err := s.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", req.Id).
        Scan(&u.Id, &u.Name, &u.Email)
    if err == sql.ErrNoRows {
        return nil, status.Errorf(codes.NotFound, "user %s not found", req.Id)
    }
    if err != nil {
        return nil, status.Errorf(codes.Internal, "database error: %v", err)
    }
    return &u, nil
}

// Server streaming: send multiple responses
func (s *UserServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    rows, err := s.db.QueryContext(stream.Context(), "SELECT id, name, email FROM users LIMIT 1000")
    if err != nil { return status.Error(codes.Internal, err.Error()) }
    defer rows.Close()
    
    for rows.Next() {
        u := &pb.User{}
        rows.Scan(&u.Id, &u.Name, &u.Email)
        if err := stream.Send(u); err != nil {
            return err // client disconnected
        }
    }
    return nil
}

// Start gRPC server:
func main() {
    lis, _ := net.Listen("tcp", ":50051")
    s := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            loggingInterceptor,
            authInterceptor,
        ),
    )
    pb.RegisterUserServiceServer(s, &UserServer{})
    reflection.Register(s) // for grpcurl debugging
    s.Serve(lis)
}

// Client:
conn, _ := grpc.Dial("user-service:50051",
    grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
    grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
)
client := pb.NewUserServiceClient(conn)

user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "123"})
```

### Why gRPC Over REST?

```
gRPC advantages:
  - Protobuf is ~5-10x more compact than JSON (binary, not text)
  - HTTP/2: multiplexing (many requests over one connection), header compression
  - Generated client/server code: no manual parsing, type-safe
  - Bidirectional streaming
  - Interceptors for cross-cutting concerns (auth, logging, tracing)

gRPC disadvantages:
  - Not browser-native (need grpc-web proxy for browser clients)
  - Harder to debug (binary format; need grpcurl or Bloom RPC)
  - Not human-readable
  
Use REST for: public APIs, browser-facing APIs
Use gRPC for: internal service communication, high-throughput microservices
```

---

## 3. WebSockets — Full-Duplex Communication

WebSockets provide a persistent, bidirectional connection between client and server. Used for: real-time chat, live updates, collaborative editing, gaming.

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // validate origin in production!
    },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil { return }
    defer conn.Close()
    
    // Set read deadline to detect dead connections
    conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    conn.SetPongHandler(func(string) error {
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    // Send pings to keep connection alive
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }()
    
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        // Echo back:
        conn.WriteMessage(messageType, message)
    }
}
```

### WebSocket Hub (Broadcast Pattern)

```go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default: // client send buffer full — disconnect slow clients
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
```

---

## 4. Server-Sent Events (SSE)

SSE is simpler than WebSockets when you only need server-to-client streaming (no client-to-server messages):

```go
func handleSSE(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    flusher := w.(http.Flusher)
    
    for {
        select {
        case <-r.Context().Done(): // client disconnected
            return
        case event := <-eventStream:
            fmt.Fprintf(w, "data: %s\n\n", event)
            flusher.Flush()
        }
    }
}
```

---

## 5. GraphQL — Flexible Querying

GraphQL lets clients request exactly the fields they need. Used by GitHub, Facebook, Shopify:

```graphql
# Client requests only what it needs:
query GetUserWithOrders {
  user(id: "123") {
    name
    email
    orders(limit: 5, status: PENDING) {
      id
      total
      items {
        productName
        quantity
      }
    }
  }
}
```

```
Pros:
  - No over-fetching (client gets only requested fields)
  - No under-fetching (one request for related data vs multiple REST calls)
  - Strongly typed schema — documentation is always up-to-date
  - Subscriptions for real-time updates

Cons:
  - N+1 problem is worse — requires DataLoader pattern
  - HTTP caching doesn't work well (all POST to /graphql)
  - Complex queries can be expensive to resolve
  - Learning curve for backend implementation
```

---

## 6. Event-Driven Communication

Asynchronous communication via events (Kafka, SQS, RabbitMQ) for decoupling services:

```
Synchronous (gRPC/REST):
  Client ──request──▶ Server
  Client ◀──response── Server
  Client blocks until server responds
  
Asynchronous (Event-driven):
  Service A ──event──▶ Kafka Topic
  Service B reads from Topic (independently, at its own pace)
  Service A doesn't wait for B to process
  
When to use events:
  - Long-running operations (don't make client wait for email to send)
  - Fan-out: one event triggers multiple consumers
  - Decoupling: service A shouldn't know about B, C, D
  - Resilience: if B is down, events queue up and are processed when B recovers
```

---

## 7. Protocol Comparison Matrix

| Protocol | Direction | Transport | Format | Best For |
|---|---|---|---|---|
| REST | Request/response | HTTP/1.1 | JSON | Public APIs, simple CRUD |
| gRPC | Request/response + streaming | HTTP/2 | Protobuf | Internal services, high performance |
| WebSocket | Bidirectional | TCP | Any | Real-time: chat, gaming, collaboration |
| SSE | Server → Client | HTTP | Text | Live updates: feeds, notifications |
| GraphQL | Request/response | HTTP | JSON | Flexible querying, mobile APIs |
| Kafka/async | Fire-and-forget | TCP | Any | Decoupling, event streaming |

---

## 8. Go Implementation Patterns

```go
// HTTP middleware pattern (works for any handler):
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// gRPC unary interceptor:
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("gRPC %s %v %v", info.FullMethod, time.Since(start), err)
    return resp, err
}
```

---

## 9. Interview Questions & Model Answers

**Q: When would you choose gRPC over REST?**

"gRPC is better for internal service-to-service communication where performance matters. Protobuf serialization is 5-10x more compact than JSON, and HTTP/2 multiplexing reduces connection overhead. The generated client/server code eliminates manual marshaling and provides type safety across services. I'd use REST for public-facing APIs because it's universally supported by browsers and HTTP clients without special libraries, and JSON is human-readable for debugging. For a mix, I'd expose REST at the API gateway for external clients and use gRPC for internal microservice communication."

**Q: How do WebSockets differ from HTTP long polling?**

"HTTP long polling is a workaround for HTTP's request-response nature: the client sends a request, and the server holds it open until it has data to send, then the client immediately sends another request. It's compatible with any HTTP infrastructure but has overhead — each 'event' requires a new HTTP request with headers. WebSockets upgrade the HTTP connection to a persistent TCP connection after the initial handshake. After that, frames can be sent in both directions at any time with minimal overhead (2-10 byte header per frame). WebSockets are more efficient for high-frequency bidirectional communication like chat or gaming, but long polling is simpler and works with standard HTTP proxies that may not support WebSocket upgrades."

---

## Summary

- **REST:** URLs as resources, HTTP verbs as actions. Use for public APIs. Stick to proper status codes.
- **gRPC:** Protobuf + HTTP/2. More efficient than REST. Use for internal microservices.
- **WebSockets:** persistent bidirectional TCP connection. Use for chat, gaming, real-time collaboration.
- **SSE:** server-to-client only streaming over HTTP. Simpler than WebSockets for one-way live updates.
- **GraphQL:** client-defined queries. Eliminates over/under-fetching. Requires DataLoader to avoid N+1.
- **Event-driven:** Kafka/SQS for async decoupling. Fire-and-forget. Consumers process independently.
- The right choice depends on: directionality (uni vs bidirectional), latency requirements, client type (browser vs service), and operational complexity tolerance.
