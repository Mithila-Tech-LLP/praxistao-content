# Chapter 107: Monolithic Architecture — When It's the Right Choice

Most teams jump to microservices before they need them. The result: operational complexity that kills velocity at exactly the stage when shipping fast matters most. Monoliths are not a stepping stone to be embarrassed about — they are the right tool for most systems at most stages. This chapter covers what a well-structured monolith looks like in Go, when it is worth keeping, and how to extract services later without a big-bang rewrite.

## Table of Contents

1. [The Monolith Is Not Dead](#1-the-monolith-is-not-dead)
2. [Types of Monoliths](#2-types-of-monoliths)
3. [Modular Monolith in Go](#3-modular-monolith-in-go)
4. [Directory Structure and Module Boundaries](#4-directory-structure-and-module-boundaries)
5. [When to Stay Monolithic](#5-when-to-stay-monolithic)
6. [When to Split: Signals That Actually Matter](#6-when-to-split-signals-that-actually-matter)
7. [Strangler Fig: Extracting Services Safely](#7-strangler-fig-extracting-services-safely)
8. [The Real Cost of Microservices](#8-the-real-cost-of-microservices)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Monolith Is Not Dead

Shopify ran on a Rails monolith past $1B in revenue. Stack Overflow still serves millions of requests per day from a single ASP.NET monolith on remarkably modest hardware. Basecamp has run on a monolith for its entire history and explicitly chose to stay there. These are not cautionary tales — they are engineering successes.

A monolith is **one deployable unit**. That is a feature, not a bug: one deploy, one rollback, one process to inspect, and no network hop between your components.

```
  MONOLITH (one process)              MICROSERVICES

  +-----------------------------+     +----------+   +----------+
  |                             |     |  Orders  |-->| Inventory|
  |  Orders | Inventory | Users |     +----------+   +----------+
  |                             |          |              |
  +-----------------------------+     +----------+   +----------+
          |                           |  Users   |   | Payments |
        [DB]                          +----------+   +----------+
                                           |              |
                                         [DB]           [DB]
```

The left diagram deploys in one `go build`. The right diagram requires container orchestration, service discovery, distributed tracing, and careful handling of network failures — before you write a single line of business logic.

---

## 2. Types of Monoliths

Not all monoliths are equal. The structure inside the single binary determines whether it stays maintainable over time.

| Type | Structure | Coupling | Testability | Refactoring ease |
|---|---|---|---|---|
| **Big Ball of Mud** | None — everything calls everything | Maximum: any change breaks anything | Very hard — no seams | Nightmare |
| **Layered Monolith** | Handler → Service → Repository layers | Medium: layers respected, but domains bleed across | Better — layers mockable | Moderate |
| **Modular Monolith** | Explicit domain packages with enforced boundaries | Low: domains communicate via interfaces or events | High — domains independently testable | Straightforward: boundaries already define service cut lines |

The **Modular Monolith** is the target. Domain packages (`orders`, `inventory`, `users`) own their internals. They do not import each other's repositories. They communicate through exported interfaces and, for async flows, an in-process event bus. When you later need to split a domain into its own service, the boundary already exists — you replace the interface with an HTTP/gRPC client.

---

## 3. Modular Monolith in Go

The key rule: `orders/service` may depend on an `inventory.Checker` interface that it defines locally. It never imports `inventory/repository` directly. This keeps coupling explicit and reversible.

```go
// internal/inventory/domain/product.go
package domain

type Stock struct {
    ProductID string
    Available int
}
```

```go
// internal/orders/service/order_service.go
package service

import (
    "context"
    "fmt"

    ordersdomain "myapp/internal/orders/domain"
    "myapp/pkg/events"
)

// Checker is defined here, in the orders package.
// inventory/service satisfies it — but orders never imports inventory internals.
type InventoryChecker interface {
    Check(ctx context.Context, productID string, qty int) (bool, error)
}

type OrderService struct {
    repo      OrderRepository
    inventory InventoryChecker
    bus       *events.Bus
}

func NewOrderService(repo OrderRepository, inv InventoryChecker, bus *events.Bus) *OrderService {
    return &OrderService{repo: repo, inventory: inv, bus: bus}
}

func (s *OrderService) PlaceOrder(ctx context.Context, productID string, qty int) (*ordersdomain.Order, error) {
    ok, err := s.inventory.Check(ctx, productID, qty)
    if err != nil {
        return nil, fmt.Errorf("inventory check: %w", err)
    }
    if !ok {
        return nil, fmt.Errorf("insufficient stock for product %s", productID)
    }

    order := &ordersdomain.Order{
        ProductID: productID,
        Quantity:  qty,
        Status:    ordersdomain.StatusPending,
    }
    if err := s.repo.Save(ctx, order); err != nil {
        return nil, fmt.Errorf("save order: %w", err)
    }

    // Publish an in-process event — inventory module subscribes to this.
    s.bus.Publish(events.Event{
        Type:    "order.placed",
        Payload: map[string]any{"order_id": order.ID, "product_id": productID, "qty": qty},
    })

    return order, nil
}
```

```go
// pkg/events/bus.go
package events

import "sync"

type Event struct {
    Type    string
    Payload map[string]any
}

type Handler func(Event)

type Bus struct {
    mu       sync.RWMutex
    handlers map[string][]Handler
}

func NewBus() *Bus {
    return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Subscribe(eventType string, h Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], h)
}

func (b *Bus) Publish(e Event) {
    b.mu.RLock()
    handlers := b.handlers[e.Type]
    b.mu.RUnlock()

    for _, h := range handlers {
        h(e) // synchronous; swap for goroutine if async is needed
    }
}
```

```go
// cmd/server/main.go — wires every module together
package main

import (
    "log"
    "net/http"

    inventoryrepo "myapp/internal/inventory/repository"
    inventorysvc  "myapp/internal/inventory/service"
    ordersrepo    "myapp/internal/orders/repository"
    orderssvc     "myapp/internal/orders/service"
    usersrepo     "myapp/internal/users/repository"
    userssvc      "myapp/internal/users/service"
    "myapp/pkg/events"
)

func main() {
    bus := events.NewBus()

    // Build each module independently.
    invRepo    := inventoryrepo.NewPostgres(db)
    invService := inventorysvc.NewInventoryService(invRepo)

    ordRepo    := ordersrepo.NewPostgres(db)
    ordService := orderssvc.NewOrderService(ordRepo, invService, bus)

    usrRepo    := usersrepo.NewPostgres(db)
    usrService := userssvc.NewUserService(usrRepo)

    // Inventory module reacts to order events.
    bus.Subscribe("order.placed", invService.OnOrderPlaced)

    mux := http.NewServeMux()
    registerOrderHandlers(mux, ordService)
    registerUserHandlers(mux, usrService)

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Notice: `main.go` is the only file that imports all three domain services. No domain package reaches into another's internals.

---

## 4. Directory Structure and Module Boundaries

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # wires everything together
├── internal/
│   ├── orders/
│   │   ├── domain/          # Order, OrderItem types
│   │   ├── service/         # OrderService + InventoryChecker interface
│   │   └── repository/      # PostgreSQL implementation
│   ├── inventory/
│   │   ├── domain/          # Product, Stock types
│   │   ├── service/         # InventoryService (also satisfies orders.InventoryChecker)
│   │   └── repository/
│   └── users/
│       ├── domain/          # User type
│       ├── service/         # UserService
│       └── repository/
├── pkg/
│   └── events/              # shared in-process event bus
└── go.mod
```

**The boundary rules:**

- `internal/orders/service` → can import `internal/orders/domain` and `internal/orders/repository`. Fine.
- `internal/orders/service` → can import `internal/inventory/domain` for shared types if unavoidable.
- `internal/orders/service` → **cannot** import `internal/inventory/repository`. That would couple two domains at the persistence layer, making them impossible to split later.
- `internal/orders/service` → **cannot** import `internal/inventory/service` directly. Use the local `InventoryChecker` interface instead and inject it from `main.go`.

Go's `internal` package visibility rule prevents code outside `myapp` from importing these packages, but it does **not** enforce cross-domain boundaries inside `myapp`. That discipline is yours to maintain — or automate with a linter (see Exercise 8).

---

## 5. When to Stay Monolithic

| Signal | Monolith is right | Microservices are right |
|---|---|---|
| **Team size** | < 15 engineers | > 50 engineers, multiple autonomous teams |
| **Deployment frequency** | Once a week or less | Multiple deploys per day per component |
| **Scaling needs** | All components scale together | Wildly different resource profiles per component |
| **Domain knowledge** | Still evolving; bounded contexts unclear | Stable, well-understood bounded contexts |
| **Operational maturity** | No dedicated SRE or platform team | Dedicated platform/SRE team, mature on-call |
| **Latency budget** | Inter-component calls must be fast | Components can tolerate ~1ms+ network latency |

If you tick three or more boxes in the "Monolith is right" column, a well-structured monolith is the better call. Complexity is not free — every dollar of engineering time spent on service meshes is a dollar not spent on features.

---

## 6. When to Split: Signals That Actually Matter

**Legitimate reasons to extract a service:**

- **Independent scaling**: your image-processing pipeline needs 32-core GPU VMs while your API runs happily on 2-core instances. Running them in the same process means paying for 32 cores everywhere.
- **Team autonomy**: two teams keep blocking each other by deploying the same binary. Merge conflicts in CI become a daily tax.
- **Different deployment lifecycles**: an ML model needs hourly retraining and redeployment; the core API deploys weekly. Coupling their release cycle is pure friction.
- **Technology diversity**: one component genuinely needs Python (ML inference, numerical libraries) or Rust (near-zero-latency path). Not a decision to make lightly.

**Not legitimate reasons:**

- "The codebase is getting big." A big, well-structured modular monolith is fine. `grep` and `go to definition` still work.
- "Tech conference talks recommend microservices." That is survivorship bias — you see the Netflixes, not the 500 teams that added microservices complexity and slowed to a crawl.
- "We might need to scale later." Premature decomposition is as wasteful as premature optimization.

---

## 7. Strangler Fig: Extracting Services Safely

Martin Fowler named this pattern after the fig tree that grows around a host tree, gradually takes over its role, and eventually replaces it entirely. Apply it to services: extract incrementally, route traffic gradually, never do a big-bang rewrite.

```
Migration stages:

Stage 1: All traffic → Monolith
  [Client] ──────────────────────> [Monolith]

Stage 2: Proxy in front; new service handles /v2/orders
  [Client] ──> [StranglerProxy] ──> [Monolith]     (all routes except /v2/orders)
                     └──────────────> [OrderService] (/v2/orders)

Stage 3: New service handles all order routes
  [Client] ──> [StranglerProxy] ──> [OrderService]  (all /orders/*)
                     └──────────────> [Monolith]    (everything else)

Stage 4: Remove proxy; client points directly to new service + stripped monolith
  [Client] ──> [OrderService]
  [Client] ──> [Monolith (orders removed)]
```

```go
// pkg/strangler/proxy.go
package strangler

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
)

// StranglerProxy routes matched path prefixes to a new service;
// everything else falls through to the monolith.
type StranglerProxy struct {
    monolith   *httputil.ReverseProxy
    newService *httputil.ReverseProxy
    routes     []string // path prefixes handled by newService
}

func NewStranglerProxy(monolithURL, newServiceURL string, routes []string) (*StranglerProxy, error) {
    mURL, err := url.Parse(monolithURL)
    if err != nil {
        return nil, err
    }
    sURL, err := url.Parse(newServiceURL)
    if err != nil {
        return nil, err
    }
    return &StranglerProxy{
        monolith:   httputil.NewSingleHostReverseProxy(mURL),
        newService: httputil.NewSingleHostReverseProxy(sURL),
        routes:     routes,
    }, nil
}

func (p *StranglerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    for _, prefix := range p.routes {
        if strings.HasPrefix(r.URL.Path, prefix) {
            p.newService.ServeHTTP(w, r)
            return
        }
    }
    p.monolith.ServeHTTP(w, r)
}
```

```go
// cmd/proxy/main.go
func main() {
    proxy, err := strangler.NewStranglerProxy(
        "http://monolith:8080",
        "http://order-service:9090",
        []string{"/v2/orders", "/v2/order-items"},
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Println("strangler proxy listening on :8000")
    log.Fatal(http.ListenAndServe(":8000", proxy))
}
```

Migration checklist:
1. Start with **read-only** endpoints (GET). Dual-write is not yet needed.
2. Add **write** endpoints once reads are stable and the new service's data store is in sync.
3. Switch **100%** of traffic to the new service; keep the monolith path in place for one release cycle as a fallback.
4. Delete the monolith code for that domain. Remove the proxy route. Declare victory.

---

## 8. The Real Cost of Microservices

| Cost | Monolith | Microservices |
|---|---|---|
| **Function call latency** | ~1 ns | ~1 ms over HTTP (1,000,000× slower) |
| **Debugging one request** | One log stream, one process | Requires distributed tracing (Jaeger, Zipkin) across N services |
| **Transactions** | Single DB ACID transaction | Must use sagas or two-phase commit — eventual consistency |
| **Deployment** | One binary, one pipeline | One pipeline per service, container images, k8s manifests |
| **Service discovery** | Not needed | Consul, k8s DNS, or a service mesh |
| **Failure modes** | Process crash (obvious) | Timeout cascades, retry storms, partial failures (subtle) |

The 1 ns vs 1 ms gap is not theoretical. A tight inner loop calling a service 1,000 times per request goes from **1 µs** to **1 s** — a second of latency you did not have when it was a function call. Every cross-service call must be treated as a distributed systems problem: timeouts, retries with backoff, circuit breakers. In a monolith all of that is free.

None of this means microservices are wrong. It means they are expensive, and the expense must be paid with a real return.

---

## Summary

- **Monoliths are not failures**: Shopify, Stack Overflow, and Basecamp built large, successful products on monoliths. One deploy, one rollback, zero network overhead.
- **Structure matters inside the monolith**: Big Ball of Mud → Layered → Modular. Only the Modular Monolith stays maintainable at scale.
- **Enforce domain boundaries with interfaces**: `orders/service` defines its own `InventoryChecker` interface; it never imports `inventory/repository`. `main.go` wires the concrete types together.
- **The event bus keeps modules decoupled**: in-process `events.Bus` gives you async communication without a message broker. Swap the handler for a Kafka producer when you extract the service.
- **Stay monolithic when**: team < 15, domains still evolving, no SRE team, scaling needs are uniform.
- **Split when**: genuinely different scaling profiles, team autonomy blocked by shared deployments, different release lifecycles, or different runtime requirements.
- **Strangler Fig is the safe extraction path**: proxy in front, migrate one endpoint prefix at a time, start with reads, delete old code only after the new service is battle-tested.
- **Microservices cost is real**: 1,000,000× slower inter-component calls, distributed tracing required, no cross-service transactions, per-service CI/CD, service discovery.

---

## Exercises

### Easy
1. Draw the directory tree for a modular monolith with three domains: `billing`, `notifications`, and `subscriptions`. Add `domain`, `service`, and `repository` sub-packages to each. Where does the shared event bus live?
2. Define a `PaymentChecker` interface inside `subscriptions/service` that `billing/service` will satisfy. Write the struct and a stub implementation. Wire them together in a `main.go` that does not import `billing/repository` from `subscriptions`.
3. Implement the `events.Bus` from Section 3 and write a unit test that subscribes two handlers to `"order.placed"`, publishes one event, and asserts both handlers were called.

### Medium
4. Build a working modular monolith with two domains (`products` and `cart`). `CartService.AddItem` must call a `ProductChecker` interface to verify the product exists. Wire it in `main.go` with an in-memory repository for each. Expose two HTTP endpoints: `POST /cart/items` and `GET /cart`.
5. Implement the `StranglerProxy` from Section 7. Write an integration test using `httptest.NewServer` for both the "monolith" and the "new service". Verify that requests to `/v2/orders` reach the new service and all other paths reach the monolith.
6. Add a module boundary enforcement test: write a Go test in `internal/orders/service` that uses `go list -f '{{.Imports}}' ./...` via `os/exec` and fails if `internal/inventory/repository` appears in the import list of any `orders` package.

### Hard
7. Migrate a feature from a monolith to a service using the Strangler Fig with **feature-flag routing**. Add a `X-Feature-Flag: new-orders` HTTP header check in the proxy: if the flag is present, route to the new service; otherwise route to the monolith. Implement a simple in-memory flag store that can be toggled via a `POST /admin/flags` endpoint. Write an integration test that demonstrates the cutover.
8. Implement a **module dependency linter** using `go/ast` and `go/parser`. It should read all `.go` files under `internal/`, detect the module each file belongs to (by its directory path), and reject any import that crosses a forbidden boundary (e.g., `orders/*` importing `inventory/repository`). Define the forbidden rules in a `depguard.yaml` config file. Run it as a `go generate` step and as a CI check.
