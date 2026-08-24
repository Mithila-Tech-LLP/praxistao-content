# Chapter 101: Final Capstone — E-Commerce Platform

This is the culminating project of the course. You'll design and build a production-grade e-commerce platform that synthesizes everything from all 12 volumes: clean architecture, databases, async processing, microservices patterns, observability, and production engineering.

This chapter is a **design document** — the full implementation would take weeks. Use it as a blueprint, reference the earlier chapters for implementation details, and build it iteratively.

---

## System Requirements

**Functional:**
- Customers browse a product catalog, add items to cart, and place orders
- Inventory is deducted on order placement; restored on cancellation
- Orders flow through: Pending → Payment Processing → Confirmed → Fulfillment → Shipped → Delivered
- Customers receive email notifications at key order events
- Sellers manage their product listings and view sales analytics
- Search by keyword, filter by category/price/rating, faceted navigation

**Non-Functional:**
- 99.9% uptime SLA
- Product listing page: < 200ms P99
- Order placement: < 500ms P99 (can be async)
- Search: < 300ms P99
- Support 10,000 concurrent users at peak

---

## Architecture

```
                    ┌─────────────────────┐
                    │     API Gateway      │  Rate limiting, Auth, Routing
                    │   (Traefik/Kong)     │
                    └──────────┬──────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        ▼                      ▼                       ▼
┌──────────────┐    ┌─────────────────┐     ┌──────────────────┐
│   Product    │    │     Order       │     │      User        │
│   Service    │    │    Service      │     │    Service       │
│  (HTTP+gRPC) │    │  (HTTP+gRPC)    │     │  (HTTP+gRPC)    │
└──────┬───────┘    └────────┬────────┘     └────────┬─────────┘
       │                     │                        │
       ▼                     ▼                        ▼
 PostgreSQL +          PostgreSQL +              PostgreSQL
 OpenSearch           Redis (inventory)
 (full-text)          + asynq (payments)

                         Kafka Event Bus
                    ┌────────┴──────────┐
                    ▼                   ▼
            Notification Service   Analytics Service
            (email, push)          (TimescaleDB)
```

---

## Domain Model

### Product Service

```go
// Domain entities
type Product struct {
    ID          ProductID
    SellerID    SellerID
    Name        string
    Slug        string
    Description string
    Category    CategoryID
    Price       Money
    Images      []Image
    Attributes  map[string]any  // JSONB: electronics → {brand, storage, color}
    Stock       int
    Status      ProductStatus   // draft | active | archived
    Rating      float64
    ReviewCount int
    CreatedAt   time.Time
}

// Bounded contexts:
// Product Service "product" = full product with SEO, images, attributes
// Order Service "product" = snapshot at time of purchase (name, price, sku)
// Search Service "product" = searchable projection (name, category, price, tags)
```

### Order Service

```go
type Order struct {
    id          OrderID
    customerID  CustomerID
    items       []OrderItem    // snapshot of product at purchase time
    shipping    Address
    payment     PaymentInfo
    status      OrderStatus
    total       Money
    createdAt   time.Time
    events      []DomainEvent
}

// Order states and transitions
type OrderStatus string
const (
    OrderDraft     OrderStatus = "draft"
    OrderPending   OrderStatus = "pending"       // placed, awaiting payment
    OrderPaid      OrderStatus = "paid"           // payment confirmed
    OrderConfirmed OrderStatus = "confirmed"      // seller confirmed, in fulfillment
    OrderShipped   OrderStatus = "shipped"        // with carrier
    OrderDelivered OrderStatus = "delivered"
    OrderCancelled OrderStatus = "cancelled"
    OrderRefunded  OrderStatus = "refunded"
)
```

---

## Key Design Decisions

### 1. Inventory — Atomic Hold with Redis

Same pattern as Ch 95 (Ticket Booking). On `AddToCart`, hold inventory in Redis with a 30-minute TTL. On `PlaceOrder`, convert the hold to a confirmed reservation. If payment fails, release the hold.

```go
// Lua script: atomic check-and-hold
holdInventoryScript := `
    local qty = tonumber(redis.call("GET", KEYS[1]))
    if not qty or qty < tonumber(ARGV[1]) then
        return 0  -- insufficient stock
    end
    redis.call("DECRBY", KEYS[1], ARGV[1])
    return 1  -- held
`

// Key: "inventory:{product_id}"
// Value: available quantity
```

### 2. Search — Dual Write PostgreSQL + OpenSearch

Products are stored in PostgreSQL (source of truth). On create/update, an outbox event is written. The relay syncs to OpenSearch for full-text search, faceted filtering, and autocomplete.

```go
// Product search index mapping
type ProductSearchDoc struct {
    ID          string          `json:"id"`
    Name        string          `json:"name"`         // text: analyzed
    Description string          `json:"description"`   // text: analyzed
    Category    string          `json:"category"`      // keyword: exact
    Brand       string          `json:"brand"`         // keyword: exact
    Price       float64         `json:"price"`         // float: range queries
    Tags        []string        `json:"tags"`          // keyword array
    Rating      float64         `json:"rating"`        // float: sort
    InStock     bool            `json:"in_stock"`      // boolean: filter
}
```

### 3. Order Events → Kafka → Multiple Consumers

```
OrderService → Kafka "orders" topic
  ├── EmailService:    sends order confirmation, shipping notification
  ├── InventoryService: confirms reservation, updates stock counts
  ├── AnalyticsService: increments category revenue counters
  └── RecommendationService: updates "customers also bought" model
```

### 4. Payment — Outbox Pattern

```go
// In one transaction:
// 1. Create order with status=pending
// 2. Write PaymentRequested event to outbox
// → Payment worker reads outbox, calls Stripe/Razorpay
// → On success: update order status + write PaymentSucceeded to outbox
// → PaymentSucceeded relayed to Kafka → all consumers notified
```

---

## Observability Stack

```yaml
Logs:    slog → Loki → Grafana
Metrics: Prometheus → Grafana (request rate, error rate, P99, business metrics)
Traces:  OpenTelemetry → Jaeger/Tempo

Key Dashboards:
  1. Business: orders/hour, revenue/hour, cart conversion rate
  2. API: request rate by service, error rates, P50/P95/P99 latency
  3. Infrastructure: CPU, memory, DB connections, Redis memory
  4. Alerts: error rate > 1%, P99 > 2s, DB connection pool exhausted
```

---

## Data Layer

```sql
-- Product Service DB
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    category_id INT NOT NULL,
    price_cents BIGINT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft',
    attributes JSONB DEFAULT '{}',
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector('english', name || ' ' || coalesce(description, ''))
    ) STORED,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_products_search ON products USING GIN (search_vector);
CREATE INDEX idx_products_category ON products (category_id, status, price_cents);

-- Order Service DB
CREATE TABLE orders (id TEXT PRIMARY KEY, ...);
CREATE TABLE order_items (order_id TEXT, product_id BIGINT, name TEXT, price_cents BIGINT, qty INT, ...);

-- Outbox (both services)
CREATE TABLE outbox (id BIGSERIAL, event_type TEXT, payload JSONB, published_at TIMESTAMPTZ);
CREATE INDEX idx_outbox_unpublished ON outbox (id) WHERE published_at IS NULL;
```

---

## CI/CD Pipeline

```
Push to main:
  1. Unit tests (go test -race ./...)
  2. Integration tests (testcontainers)
  3. govulncheck + gosec
  4. golangci-lint
  5. Build Docker image → push to registry

Auto-deploy to staging:
  6. kubectl rollout to staging
  7. k6 smoke tests (100 VU, 2 min, < 1% errors, P99 < 500ms)
  8. If pass → notify team

Manual approval → deploy to production:
  9. kubectl rollout to production (rolling update, maxUnavailable=0)
  10. Monitor error rate and latency for 5 min
  11. Auto-rollback if error rate > 2%
```

---

## Implementation Roadmap

### Phase 1 — Core (Week 1-2)
- [ ] Product Service: CRUD, PostgreSQL FTS, OpenSearch sync
- [ ] User Service: register, login, JWT auth
- [ ] Basic HTTP API with chi router

### Phase 2 — Orders (Week 3-4)
- [ ] Cart with Redis inventory holds
- [ ] Order placement with outbox
- [ ] asynq payment worker (mock payment gateway)

### Phase 3 — Events (Week 5)
- [ ] Kafka consumer: email notifications (asynq)
- [ ] Kafka consumer: analytics projection (TimescaleDB)
- [ ] Order status updates via WebSocket

### Phase 4 — Search and Recommendations (Week 6)
- [ ] OpenSearch autocomplete
- [ ] Faceted search (category, brand, price range, rating)
- [ ] "Similar products" based on category/tags

### Phase 5 — Production Hardening (Week 7-8)
- [ ] Full observability: slog + Prometheus + OpenTelemetry
- [ ] Grafana dashboards + alert rules
- [ ] k6 load tests + performance optimization
- [ ] Docker + Kubernetes manifests
- [ ] GitHub Actions CI/CD pipeline

---

## You Made It

This final project represents the full arc of the course:

| Volume | Contribution |
|--------|-------------|
| Vol 1: Go Basics | Go syntax, types, error handling |
| Vol 2: Concurrency | Worker pools, channels for async processing |
| Vol 3: Data Structures | Sorted sets, caches, probabilistic structures |
| Vol 4: Algorithms | Efficient search, sorting, path finding |
| Vol 5: Web | HTTP API, auth middleware, file uploads |
| Vol 6: Databases | PostgreSQL, Redis, MongoDB, OpenSearch, TimescaleDB |
| Vol 7: Architecture | Clean arch, DDD, CQRS, outbox, saga, config |
| Vol 8: Async | Kafka, asynq, worker pools, event-driven design |
| Vol 9: Microservices | gRPC, service discovery, circuit breakers, CAP theorem |
| Vol 10: Observability | slog, Prometheus, OpenTelemetry |
| Vol 11: Production | Docker, Kubernetes, CI/CD, security, load testing |

**What to build next after this course:**
- Contribute to an open-source Go project
- Build a SaaS side project using this stack
- Study distributed systems deeper: DDIA (Designing Data-Intensive Applications)
- Go internals: read the Go compiler source, understand the runtime scheduler
