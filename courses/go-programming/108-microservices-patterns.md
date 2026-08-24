# Chapter 96: Microservices — Patterns and Communication

Microservices split a system into independently deployable services. The core tradeoff: operational complexity (many services to deploy, monitor, and keep in sync) in exchange for independent scalability, independent deployments, and team autonomy.

## Table of Contents

1. [When to Split](#1-when-to-split)
2. [Service Communication](#2-service-communication)
3. [gRPC in Go](#3-grpc-in-go)
4. [Service Discovery](#4-service-discovery)
5. [API Gateway](#5-api-gateway)
6. [Circuit Breaker and Retry](#6-circuit-breaker-and-retry)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. When to Split

**Start as a monolith.** Split when:
- Teams are stepping on each other (coupling causes long PR review cycles)
- A specific component needs independent scaling (payment service gets 10× traffic on sale days)
- Deployment bottleneck: one team's unrelated change blocks everyone's release
- One part needs different runtime requirements (ML inference in Python, core logic in Go)

**Don't split when:**
- You're still figuring out the domain model — premature decomposition forces you to refactor across services
- The team is small (< 10 engineers) — operational overhead outweighs the benefits
- You don't have solid observability — distributed systems are impossible to debug without tracing

**Wrong reasons to split:**
- "Microservices scale better" — a well-optimized monolith scales fine for most startups
- "It's the modern way" — it's a solution to specific organizational and scale problems

---

## 2. Service Communication

```
Synchronous (request-response):
  REST (HTTP/JSON)    ← simple, human-readable, tooling everywhere
  gRPC (HTTP/2+protobuf) ← fast, type-safe, bi-directional streaming

Asynchronous (event-driven):
  Kafka               ← durable, fan-out, replay
  RabbitMQ/NATS       ← simpler, lower-latency messaging
  asynq (Redis)       ← background jobs
```

### Choosing communication style

- **REST**: external-facing APIs, simple request-response between services
- **gRPC**: internal service-to-service where performance matters; strongly typed contracts
- **Async events**: when the caller doesn't need an immediate response (email, notifications, analytics)

Never do synchronous chains longer than 2-3 hops. `A→B→C→D` means A's p99 = B's p99 + C's p99 + D's p99.

---

## 3. gRPC in Go

gRPC uses Protocol Buffers for serialization (smaller, faster than JSON) and HTTP/2 (multiplexed, streaming).

### Define the service contract

```protobuf
// proto/order/v1/order.proto
syntax = "proto3";
package order.v1;
option go_package = "github.com/myapp/proto/order/v1;orderv1";

service OrderService {
    rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
    rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
    rpc WatchOrderStatus(WatchOrderStatusRequest) returns (stream OrderStatusUpdate);
}

message GetOrderRequest  { string order_id = 1; }
message GetOrderResponse { Order order = 1; }

message Order {
    string id          = 1;
    string customer_id = 2;
    string status      = 3;
    double total       = 4;
    int64  created_at  = 5;  // unix timestamp
}

message WatchOrderStatusRequest { string order_id = 1; }
message OrderStatusUpdate {
    string order_id = 1;
    string status   = 2;
    int64  updated_at = 3;
}
```

Generate Go code:
```bash
protoc --go_out=. --go-grpc_out=. proto/order/v1/order.proto
```

### gRPC server

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    orderv1 "github.com/myapp/proto/order/v1"
)

type OrderServer struct {
    orderv1.UnimplementedOrderServiceServer
    orders domain.OrderRepository
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
    order, err := s.orders.GetByID(ctx, req.OrderId)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return nil, status.Errorf(codes.NotFound, "order %s not found", req.OrderId)
        }
        return nil, status.Errorf(codes.Internal, "internal error")
    }
    return &orderv1.GetOrderResponse{Order: toProto(order)}, nil
}

func (s *OrderServer) WatchOrderStatus(req *orderv1.WatchOrderStatusRequest, stream orderv1.OrderService_WatchOrderStatusServer) error {
    // Send initial status
    order, _ := s.orders.GetByID(stream.Context(), req.OrderId)
    stream.Send(&orderv1.OrderStatusUpdate{
        OrderId:  order.ID,
        Status:   string(order.Status),
        UpdatedAt: order.UpdatedAt.Unix(),
    })
    
    // Poll for status changes (in production: use Pub/Sub or Watch)
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    prev := order.Status
    for {
        select {
        case <-stream.Context().Done(): return nil
        case <-ticker.C:
            order, _ = s.orders.GetByID(stream.Context(), req.OrderId)
            if order.Status != prev {
                prev = order.Status
                stream.Send(&orderv1.OrderStatusUpdate{
                    OrderId: order.ID, Status: string(order.Status),
                    UpdatedAt: time.Now().Unix(),
                })
            }
        }
    }
}

func startGRPCServer(port int, orders domain.OrderRepository) *grpc.Server {
    srv := grpc.NewServer(
        grpc.UnaryInterceptor(grpc.ChainUnaryInterceptor(
            grpcLoggingInterceptor,
            grpcAuthInterceptor,
        )),
    )
    orderv1.RegisterOrderServiceServer(srv, &OrderServer{orders: orders})
    
    lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
    go srv.Serve(lis)
    return srv
}
```

### gRPC client

```go
func newOrderClient(addr string) (orderv1.OrderServiceClient, *grpc.ClientConn, error) {
    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()), // use TLS in production
        grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10<<20)),
        grpc.WithKeepaliveParams(keepalive.ClientParameters{
            Time:    10 * time.Second,
            Timeout: 5 * time.Second,
        }),
    )
    if err != nil { return nil, nil, err }
    return orderv1.NewOrderServiceClient(conn), conn, nil
}

// Use the client
client, conn, _ := newOrderClient("order-service:50051")
defer conn.Close()

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.GetOrder(ctx, &orderv1.GetOrderRequest{OrderId: orderID})
```

---

## 4. Service Discovery

Services need to find each other. Don't hardcode IP addresses — use service discovery.

```go
// Option 1: DNS-based discovery (simplest — Kubernetes does this automatically)
// order-service:50051 resolves to the current pod IPs via kube-dns

// Option 2: Consul
import "github.com/hashicorp/consul/api"

func registerService(name, addr string, port int, tags []string) error {
    client, _ := api.NewClient(api.DefaultConfig())
    
    return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
        Name:    name,
        Address: addr,
        Port:    port,
        Tags:    tags,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/healthz", addr, port),
            Interval: "10s",
            Timeout:  "2s",
            DeregisterCriticalServiceAfter: "60s",
        },
    })
}

func discoverService(name string) (string, error) {
    client, _ := api.NewClient(api.DefaultConfig())
    
    services, _, err := client.Health().Service(name, "", true, nil)
    if err != nil || len(services) == 0 {
        return "", fmt.Errorf("no healthy instances of %s", name)
    }
    
    // Random selection for load balancing
    s := services[rand.Intn(len(services))]
    return fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port), nil
}
```

---

## 5. API Gateway

The API Gateway is the single entry point for external clients. It handles routing, auth, rate limiting, and protocol translation.

```go
// Traefik configuration (traefik.yml)
/*
http:
  routers:
    orders:
      rule: "PathPrefix(`/api/orders`)"
      service: order-service
      middlewares:
        - auth
        - rate-limit
    products:
      rule: "PathPrefix(`/api/products`)"
      service: product-service

  services:
    order-service:
      loadBalancer:
        servers:
          - url: "http://order-service:8080"
    product-service:
      loadBalancer:
        servers:
          - url: "http://product-service:8080"

  middlewares:
    auth:
      forwardAuth:
        address: "http://auth-service:8080/verify"
    rate-limit:
      rateLimit:
        average: 100
        burst: 50
*/

// Or build a simple gateway in Go
type Gateway struct {
    routes map[string]*httputil.ReverseProxy
    auth   AuthMiddleware
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. Auth
    userID, err := g.auth.Verify(r)
    if err != nil { http.Error(w, "unauthorized", 401); return }
    
    // 2. Route
    for prefix, proxy := range g.routes {
        if strings.HasPrefix(r.URL.Path, prefix) {
            r.Header.Set("X-User-ID", userID)
            proxy.ServeHTTP(w, r)
            return
        }
    }
    http.NotFound(w, r)
}
```

---

## 6. Circuit Breaker and Retry

```go
import "github.com/sony/gobreaker"

// Circuit breaker: open after 5 consecutive failures, retry after 60s
cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "order-service",
    MaxRequests: 5,                 // max calls when half-open
    Interval:    60 * time.Second,  // reset counts after this period
    Timeout:     60 * time.Second,  // how long to stay open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
    OnStateChange: func(name string, from, to gobreaker.State) {
        slog.Warn("circuit breaker state change",
            "name", name, "from", from, "to", to)
    },
})

func getOrderWithBreaker(ctx context.Context, client orderv1.OrderServiceClient, id string) (*orderv1.Order, error) {
    result, err := cb.Execute(func() (interface{}, error) {
        ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
        defer cancel()
        resp, err := client.GetOrder(ctx, &orderv1.GetOrderRequest{OrderId: id})
        if err != nil { return nil, err }
        return resp.Order, nil
    })
    if err != nil {
        if err == gobreaker.ErrOpenState {
            return nil, errors.New("order service unavailable (circuit open)")
        }
        return nil, err
    }
    return result.(*orderv1.Order), nil
}
```

---

## Summary

- Start as a monolith; split when teams, scaling, or deployment independence demands it
- **gRPC**: protobuf contracts, HTTP/2, strongly typed; use for internal service communication
- **Service discovery**: DNS (Kubernetes) or Consul; never hardcode service addresses
- **API Gateway**: single entry point; handles routing, auth, rate limiting, SSL termination
- **Circuit breaker**: fail fast when a downstream service is down; prevent cascade failures
- Avoid synchronous chains > 2-3 hops; use async events for fire-and-forget

## Exercises

### Easy
1. Define a gRPC service `ProductService` with `GetProduct(id)` and `SearchProducts(query, page, pageSize)`. Generate Go code. Implement the server and a test that calls it.
2. Add a gRPC interceptor that logs every request with method name, duration, and status code. Test it with 10 calls and verify the logs.
3. Implement a basic API gateway in Go: route `/api/orders/*` to `order-service:8080` and `/api/products/*` to `product-service:8081` using `httputil.ReverseProxy`.

### Medium
4. Add **TLS to gRPC**: generate a self-signed certificate with `openssl`, configure the server to use it, and configure the client to trust it. Test that unencrypted connections are rejected.
5. Implement **Consul service registration and discovery**: register two instances of a service, discover them with health filtering, and implement round-robin load balancing across healthy instances.
6. Build a **bulkhead pattern**: separate connection pools for critical and non-critical downstream calls. If the non-critical pool is saturated, non-critical calls fail immediately rather than blocking threads needed for critical operations.

### Hard
7. Implement a **service mesh sidecar** (simplified): a proxy that intercepts all outbound HTTP calls, adds tracing headers, enforces circuit breaking, and logs call statistics. Wire it into a test application using `http.Transport` wrapping.
8. Build a **blue-green deployment test harness**: route 10% of traffic to the "green" service (new version) and 90% to "blue" (old version) via the gateway. Monitor error rates. If green's error rate exceeds blue's by 2×, automatically shift all traffic back to blue.
