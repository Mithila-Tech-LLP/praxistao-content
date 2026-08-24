# Chapter 127: Mini Project 7 — Fully Observable API Service

Every production API eventually hits a moment where something breaks silently — latency climbs, errors spike, and you have no idea why. Observability prevents that. This chapter builds an order management API that is fully instrumented from day one: every HTTP request is logged, every metric is tracked, every error is captured, and every slow database query is traceable end-to-end.

---

## 1. Goal and Tech Stack

You are building a small but production-grade order management API. The goal is not the business logic — it is the instrumentation layer that wraps it.

| Concern | Tool |
|---|---|
| HTTP routing | chi v5 |
| Database | PostgreSQL via pgx v5 |
| Structured logging | slog (stdlib, JSON handler) |
| Metrics | Prometheus (prometheus/client_golang) |
| Distributed tracing | OpenTelemetry SDK + OTLP exporter to Jaeger |
| Error capture | Sentry Go SDK |
| Infrastructure | Docker Compose |

The API exposes four routes:

```
POST /orders         create an order
GET  /orders/{id}    get an order by ID
GET  /healthz        liveness probe
GET  /readyz         readiness probe (checks DB)
```

Every request flows through three middleware layers — logging, metrics, tracing — before reaching the handler. Panics and unhandled errors flow through a Sentry recovery layer.

---

## 2. Project Structure

```
observable-api/
  cmd/
    api/
      main.go
  internal/
    orders/
      handler.go
      service.go
      repository.go
    observability/
      logging.go
      metrics.go
      tracing.go
      sentry.go
  docker-compose.yml
  prometheus.yml
```

The `observability` package owns all instrumentation setup. The `orders` package owns business logic and is instrumented by calling observability helpers, not by importing observability libraries directly into every file.

---

## 3. Observability Bootstrap

### observability/logging.go

```go
package observability

import (
    "context"
    "log/slog"
    "os"

    "go.opentelemetry.io/otel/trace"
)

// NewLogger returns a JSON slog logger writing to stdout.
func NewLogger() *slog.Logger {
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}

// LoggerFromContext returns a logger enriched with trace_id and span_id
// extracted from the active span in ctx. Falls back to the base logger.
func LoggerFromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
    span := trace.SpanFromContext(ctx)
    if !span.IsRecording() {
        return base
    }
    sc := span.SpanContext()
    return base.With(
        slog.String("trace_id", sc.TraceID().String()),
        slog.String("span_id", sc.SpanID().String()),
    )
}
```

### observability/metrics.go

```go
package observability

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // RequestsTotal counts every HTTP request by method, path, and status code.
    RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "requests_total",
        Help: "Total HTTP requests.",
    }, []string{"method", "path", "status"})

    // RequestDuration tracks latency per route.
    RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "request_duration_seconds",
        Help:    "HTTP request latency.",
        Buckets: prometheus.DefBuckets,
    }, []string{"method", "path", "status"})

    // ActiveRequests is the number of in-flight requests.
    ActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "active_requests",
        Help: "In-flight HTTP requests.",
    })

    // DBQueryDuration tracks database query latency by operation name.
    DBQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "db_query_duration_seconds",
        Help:    "Database query latency.",
        Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.5},
    }, []string{"operation"})
)
```

### observability/tracing.go

```go
package observability

import (
    "context"
    "fmt"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
    "go.opentelemetry.io/otel/trace"
)

// InitTracer sets up an OTLP gRPC exporter pointing at Jaeger.
// Call the returned shutdown function in main via defer.
func InitTracer(ctx context.Context, serviceName, jaegerEndpoint string) (func(context.Context) error, error) {
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(jaegerEndpoint), // e.g. "jaeger:4317"
        otlptracegrpc.WithInsecure(),
    )
    if err != nil {
        return nil, fmt.Errorf("create otlp exporter: %w", err)
    }

    res, _ := resource.New(ctx,
        resource.WithAttributes(semconv.ServiceName(serviceName)),
    )

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)

    return tp.Shutdown, nil
}

// Tracer returns the named tracer for this service.
func Tracer() trace.Tracer {
    return otel.Tracer("observable-api")
}
```

### observability/sentry.go

```go
package observability

import (
    "fmt"
    "time"

    "github.com/getsentry/sentry-go"
)

// InitSentry configures the Sentry SDK. DSN comes from environment.
func InitSentry(dsn, environment string) error {
    err := sentry.Init(sentry.ClientOptions{
        Dsn:              dsn,
        Environment:      environment,
        TracesSampleRate: 0.2,
    })
    if err != nil {
        return fmt.Errorf("sentry init: %w", err)
    }
    return nil
}

// FlushSentry flushes buffered events before process exit.
func FlushSentry() {
    sentry.Flush(2 * time.Second)
}
```

---

## 4. Step 1 — Structured Logging Middleware

```go
// internal/orders/handler.go (middleware section)

package orders

import (
    "log/slog"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5/middleware"
    "github.com/google/uuid"
    "github.com/yourorg/observable-api/internal/observability"
)

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := uuid.NewString()

            // ctx already carries the trace span started by otelhttp,
            // since this middleware runs after it.
            ctx := r.Context()
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

            next.ServeHTTP(ww, r.WithContext(ctx))

            // After the handler runs, the span is still active in ctx.
            log := observability.LoggerFromContext(ctx, logger)
            log.InfoContext(ctx, "request completed",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.Int("status", ww.Status()),
                slog.Duration("duration", time.Since(start)),
                slog.String("request_id", requestID),
            )
        })
    }
}
```

Every log line now looks like:

```json
{
  "time": "2026-06-30T10:15:22Z",
  "level": "INFO",
  "msg": "request completed",
  "method": "POST",
  "path": "/orders",
  "status": 201,
  "duration": "4.2ms",
  "request_id": "a3f1...",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7"
}
```

The `trace_id` field lets you jump directly from a log line to the matching Jaeger trace.

---

## 5. Step 2 — Prometheus Metrics Middleware

To capture the HTTP status code after the handler runs, wrap the `ResponseWriter` with chi's `WrapResponseWriter`.

```go
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        observability.ActiveRequests.Inc()
        defer observability.ActiveRequests.Dec()

        start := time.Now()
        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

        next.ServeHTTP(ww, r)

        status := fmt.Sprintf("%d", ww.Status())
        observability.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        observability.RequestDuration.
            WithLabelValues(r.Method, r.URL.Path, status).
            Observe(time.Since(start).Seconds())
    })
}
```

Register the `/metrics` endpoint in `main.go`:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

r.Handle("/metrics", promhttp.Handler())
```

---

## 6. Step 3 — OpenTelemetry Tracing

Traces flow through three layers: handler -> service -> repository. Each layer creates a child span from the incoming context.

```
HTTP Request
    │
    ▼
[HTTP Handler span]  ← root span, started by otelhttp middleware
    │
    ├── [service.CreateOrder span]
    │       │
    │       └── [repository.Insert span]  ← annotated with db.statement
```

### Handler — extract and start span

```go
// internal/orders/handler.go

import (
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel/attribute"
)

func NewRouter(h *Handler, logger *slog.Logger) http.Handler {
    r := chi.NewRouter()

    // otelhttp wraps the entire router, creating a root span per request.
    r.Use(otelhttp.NewMiddleware("observable-api"))
    r.Use(LoggingMiddleware(logger))
    r.Use(MetricsMiddleware)
    r.Use(observability.SentryMiddleware) // lives in the observability package

    r.Post("/orders", h.CreateOrder)
    r.Get("/orders/{id}", h.GetOrder)
    r.Get("/healthz", h.Healthz)
    r.Get("/readyz", h.Readyz)
    r.Handle("/metrics", promhttp.Handler())

    return r
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // otelhttp already started a root span. Add attributes.
    span := trace.SpanFromContext(ctx)
    span.SetAttributes(attribute.String("http.body_size", r.Header.Get("Content-Length")))

    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    order, err := h.svc.CreateOrder(ctx, req)
    if err != nil {
        sentry.CaptureException(err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}
```

### Service — child span

```go
// internal/orders/service.go

func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    ctx, span := observability.Tracer().Start(ctx, "orders.CreateOrder")
    defer span.End()

    span.SetAttributes(attribute.String("order.customer_id", req.CustomerID))

    order := &Order{
        ID:         uuid.NewString(),
        CustomerID: req.CustomerID,
        Items:      req.Items,
        Status:     "pending",
    }

    if err := s.repo.Insert(ctx, order); err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    return order, nil
}
```

### Repository — child span with DB attributes

```go
// internal/orders/repository.go

func (r *Repository) Insert(ctx context.Context, order *Order) error {
    ctx, span := observability.Tracer().Start(ctx, "orders.repository.Insert")
    defer span.End()

    span.SetAttributes(
        attribute.String("db.system", "postgresql"),
        attribute.String("db.operation", "INSERT"),
        attribute.String("db.name", "orders"),
    )

    start := time.Now()
    _, err := r.db.Exec(ctx,
        `INSERT INTO orders (id, customer_id, status) VALUES ($1, $2, $3)`,
        order.ID, order.CustomerID, order.Status,
    )
    observability.DBQueryDuration.WithLabelValues("insert_order").Observe(time.Since(start).Seconds())

    if err != nil {
        span.RecordError(err)
        return fmt.Errorf("insert order: %w", err)
    }
    return nil
}
```

The ctx flows from handler to service to repository unchanged. Each `Start` call looks up the active span in ctx and makes the new span a child automatically.

---

## 7. Step 4 — Sentry Error Capture Middleware

```go
// internal/observability/sentry.go (middleware)

func SentryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // sentryHub binds this request to a Sentry hub clone with request context.
        hub := sentry.GetHubFromContext(ctx)
        if hub == nil {
            hub = sentry.CurrentHub().Clone()
            ctx = sentry.SetHubOnContext(ctx, hub)
        }

        hub.Scope().SetRequest(r)

        defer func() {
            if err := recover(); err != nil {
                hub.RecoverWithContext(ctx, err)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

For non-panic errors, call `sentry.CaptureException(err)` directly in the handler after checking the error return. To associate errors with a user:

```go
hub := sentry.GetHubFromContext(ctx)
hub.Scope().SetUser(sentry.User{ID: userID, Email: userEmail})
hub.CaptureException(err)
```

---

## 8. Step 5 — Health Endpoints

```go
// Healthz responds immediately — used by Kubernetes liveness probe.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}

// Readyz pings the database — used by Kubernetes readiness probe.
// Returns 503 if the DB is unreachable so the load balancer stops routing traffic.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    if err := h.db.Ping(ctx); err != nil {
        http.Error(w, "database unavailable", http.StatusServiceUnavailable)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ready"))
}
```

The split matters: a crashing pod should fail liveness (triggering a restart), while a pod with a broken DB connection should fail readiness (stopping traffic routing) without restarting.

---

## 9. Docker Compose

### docker-compose.yml

```yaml
version: "3.9"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://app:secret@postgres:5432/orders?sslmode=disable
      JAEGER_ENDPOINT: jaeger:4317
      SENTRY_DSN: ${SENTRY_DSN}
      APP_ENV: development
    depends_on:
      postgres:
        condition: service_healthy
      jaeger:
        condition: service_started

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: orders
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app -d orders"]
      interval: 5s
      timeout: 5s
      retries: 5

  prometheus:
    image: prom/prometheus:v2.51.0
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:10.3.3
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    depends_on:
      - prometheus

  jaeger:
    image: jaegertracing/all-in-one:1.56
    ports:
      - "16686:16686"   # Jaeger UI
      - "4317:4317"     # OTLP gRPC receiver
    environment:
      COLLECTOR_OTLP_ENABLED: "true"
```

### prometheus.yml

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "observable-api"
    static_configs:
      - targets: ["app:8080"]
    metrics_path: /metrics
```

---

## 10. Grafana Dashboard Panels

After starting Docker Compose, open Grafana at `http://localhost:3000` and add Prometheus as a data source (`http://prometheus:9090`). Create a new dashboard with these four panels.

**Panel 1 — Request Rate**
PromQL: `rate(requests_total[1m])`
Shows requests per second over a 1-minute window, split by `method` and `path`. Use a time series graph. This is the first thing to check when an incident starts.

**Panel 2 — Error Rate**
PromQL: `rate(requests_total{status=~"5.."}[1m]) / rate(requests_total[1m])`
Shows the fraction of requests returning 5xx. Alert when this crosses 0.01 (1%).

**Panel 3 — P99 Latency**
PromQL: `histogram_quantile(0.99, rate(request_duration_seconds_bucket[5m]))`
Shows the 99th percentile request latency over a 5-minute window. A climbing P99 is the first sign of a slow query or resource exhaustion.

**Panel 4 — DB Query Duration (P95)**
PromQL: `histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))`
Shows the 95th percentile database query latency split by `operation` label. Useful for identifying which query is the bottleneck.

---

## 11. Verification

Start everything:

```bash
docker-compose up --build
```

Create an order:

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id": "cust-1", "items": ["item-a", "item-b"]}'
```

Check readiness:

```bash
curl http://localhost:8080/readyz
# ready

curl http://localhost:8080/healthz
# ok
```

View the raw Prometheus metrics:

```bash
curl http://localhost:8080/metrics | grep requests_total
```

**Jaeger UI** — open `http://localhost:16686`, select service `observable-api`, click Find Traces. Each trace shows three spans: the root HTTP span, the service span, and the repository span. Click any span to see its attributes including `db.operation` and `db.name`.

**Prometheus** — open `http://localhost:9090/graph` and query `requests_total`. You should see counters incrementing after each curl.

**Grafana** — open `http://localhost:3000`, log in with `admin/admin`. The four panels update in real time.

**Sentry** — trigger a deliberate error by sending a malformed request. The exception appears in your Sentry project within seconds.

### Load Test with k6

```javascript
// loadtest.js
import http from "k6/http";
import { sleep } from "k6";

export const options = {
  vus: 50,
  duration: "30s",
};

export default function () {
  const payload = JSON.stringify({
    customer_id: `cust-${Math.floor(Math.random() * 100)}`,
    items: ["item-a"],
  });
  http.post("http://localhost:8080/orders", payload, {
    headers: { "Content-Type": "application/json" },
  });
  sleep(0.02); // ~50 req/s per VU
}
```

Run with `k6 run loadtest.js`. Watch the P99 panel in Grafana climb and stabilize. Check Jaeger for traces with slow DB spans.

---

## 12. Summary

```
HTTP Request
     │
     ▼
[otelhttp middleware]   ← creates root trace span
     │
[LoggingMiddleware]     ← logs at end with trace_id + span_id
     │
[MetricsMiddleware]     ← tracks active requests, latency, status codes
     │
[SentryMiddleware]      ← recovers panics, scopes user context
     │
[Handler]               ← decodes request, calls service
     │
[Service span]          ← child span, records business-level attributes
     │
[Repository span]       ← child span, records db.operation, db.name
     │                    measures DBQueryDuration histogram
[PostgreSQL]
```

The key principle is that observability is wired once at the middleware layer and propagated through `context.Context`. Individual handlers and services do not need to import logging or metrics libraries to participate — they get trace context from `ctx`, log with the enriched logger, and their DB calls are automatically measured by the repository layer.

The three signals are correlated via `trace_id`: a Prometheus alert fires, you find the log line for that time window, the log line contains the `trace_id`, and you paste it into Jaeger to see exactly which span was slow and why.

---

## Extension Challenges

**Easy — Expose and verify metrics endpoint**

Add a `/metrics` endpoint to the chi router that serves the default Prometheus registry via `promhttp.Handler()`. Update `prometheus.yml` to include `metrics_path: /metrics` for the `observable-api` scrape config. Send 10 requests to `POST /orders` and then query `http://localhost:9090/graph` with the expression `requests_total{path="/orders"}`. Verify the counter value equals 10. Also confirm that `active_requests` returns to 0 after all requests complete.

**Medium — Instrument an in-memory cache layer**

Add a cache layer to `OrderService` using a `sync.RWMutex`-protected `map[string]*Order`. Before calling the repository on a `GetOrder` request, check the map. Add two new Prometheus counters: `cache_hits_total` and `cache_misses_total`, both with a `method` label. Add an OpenTelemetry span around the cache lookup named `orders.cache.Get`. Set a span attribute `cache.hit` to true or false. In Jaeger, confirm that a cache hit produces a shorter overall trace with the repository span absent. In Grafana, add a panel showing the cache hit ratio: `rate(cache_hits_total[1m]) / (rate(cache_hits_total[1m]) + rate(cache_misses_total[1m]))`.

**Hard — Distributed tracing across two services**

Create a second service, `payment-service`, as a separate Go binary that exposes one endpoint: `POST /payments`. In the order service, after inserting the order, call the payment service over HTTP using a standard `http.Client`. Propagate the trace context using the W3C TraceContext format by injecting the outgoing request headers with `otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))`. In the payment service, extract the trace context from the incoming request headers with `otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))` before starting your first span. Both services should use the same Jaeger instance. In the Jaeger UI, search for a trace that started in `observable-api`. You should see two services in a single trace: `observable-api` as the parent and `payment-service` as a child. This confirms that the W3C TraceContext headers are being correctly propagated across the HTTP boundary.
