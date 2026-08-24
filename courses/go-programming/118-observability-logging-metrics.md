# Chapter 118: Observability — Logging, Metrics, and Tracing

Observability answers the question: "what is my system doing right now and why?" Without it, distributed systems are impossible to debug in production. The three pillars are logs, metrics, and traces. This chapter covers all three with Go's standard `slog`, Prometheus metrics, and OpenTelemetry tracing.

## Table of Contents

1. [The Three Pillars](#1-the-three-pillars)
2. [Structured Logging with slog](#2-structured-logging-with-slog)
3. [Prometheus Metrics](#3-prometheus-metrics)
4. [OpenTelemetry Tracing](#4-opentelemetry-tracing)
5. [Connecting the Three](#5-connecting-the-three)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. The Three Pillars

| Pillar | Answers | Tool |
|--------|---------|------|
| **Logs** | What happened? Event stream with details. | slog → Loki/OpenSearch |
| **Metrics** | How is it performing? Aggregated numbers. | Prometheus → Grafana |
| **Traces** | Why did this request take 800ms? | OpenTelemetry → Jaeger/Tempo |

These work together:
- Metrics alert you that latency spiked at 14:30
- Logs show you what errors occurred at 14:30
- Traces show you which service call took 750ms

---

## 2. Structured Logging with slog

Go 1.21 added `log/slog` — structured, leveled logging built into the standard library.

```go
import "log/slog"

// Setup: JSON output for production, text for development
func newLogger(env string) *slog.Logger {
    var handler slog.Handler
    opts := &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: true, // add file:line to every log entry
    }
    
    if env == "production" {
        handler = slog.NewJSONHandler(os.Stdout, opts)
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
    }
    
    return slog.New(handler)
}

// Replace the default logger globally
func init() {
    slog.SetDefault(newLogger(os.Getenv("ENV")))
}
```

### Structured key-value logging

```go
// Use structured fields, not fmt.Sprintf
slog.Info("user registered",
    "user_id",  user.ID,
    "email",    user.Email,
    "duration", time.Since(start),
)

// Avoid: loses the structure
slog.Info(fmt.Sprintf("user %d registered email %s", user.ID, user.Email))

// Groups: namespace related fields
slog.Info("order placed",
    slog.Group("order",
        "id",    order.ID,
        "total", order.Total,
        "items", len(order.Items),
    ),
    slog.Group("customer",
        "id",    user.ID,
        "email", user.Email,
    ),
)
```

### Context-carrying logger

```go
type contextKey struct{}

func LoggerFromContext(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok { return l }
    return slog.Default()
}

func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, contextKey{}, l)
}

// statusRecorder wraps http.ResponseWriter to capture the status code
// and bytes written (the ResponseWriter alone doesn't expose them)
type statusRecorder struct {
    http.ResponseWriter
    code  int
    bytes int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.code = code
    r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
    n, err := r.ResponseWriter.Write(b)
    r.bytes += n
    return n, err
}

// Middleware: add request-scoped fields to every log in this request
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" { requestID = uuid.New().String() }
            
            reqLogger := logger.With(
                "request_id", requestID,
                "method",     r.Method,
                "path",       r.URL.Path,
                "remote_ip",  r.RemoteAddr,
            )
            
            ctx := WithLogger(r.Context(), reqLogger)
            rec := &statusRecorder{ResponseWriter: w, code: 200}
            start := time.Now()
            
            next.ServeHTTP(rec, r.WithContext(ctx))
            
            reqLogger.Info("request completed",
                "status",   rec.code,
                "duration", time.Since(start),
                "bytes",    rec.bytes,
            )
        })
    }
}

// Usage in handler
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    log := LoggerFromContext(r.Context())
    log.Info("creating order", "customer_id", getCustomerID(r))
    
    order, err := h.svc.CreateOrder(r.Context(), req)
    if err != nil {
        log.Error("create order failed", "err", err)
        http.Error(w, "internal error", 500)
        return
    }
    log.Info("order created", "order_id", order.ID)
    writeJSON(w, 201, order)
}
```

---

## 3. Prometheus Metrics

Prometheus scrapes metrics from your service periodically. You expose a `/metrics` endpoint; Prometheus polls it.

```go
import "github.com/prometheus/client_golang/prometheus"
import "github.com/prometheus/client_golang/prometheus/promauto"
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Define metrics at package level
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
        },
        []string{"method", "path"},
    )
    
    activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "active_connections",
        Help: "Number of active connections",
    })
    
    orderProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "order_processing_duration_seconds",
        Help:    "Order processing duration",
        Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
    })
    
    ordersTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "orders_total",
        Help: "Total orders",
    }, []string{"status"})
)

// Metrics middleware
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        activeConnections.Inc()
        defer activeConnections.Dec()
        
        rec := &statusRecorder{ResponseWriter: w, code: 200}
        next.ServeHTTP(rec, r)
        
        status := strconv.Itoa(rec.code)
        path := sanitizePath(r.URL.Path) // avoid cardinality explosion from /users/1, /users/2...
        
        httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
        httpRequestDuration.WithLabelValues(r.Method, path).Observe(time.Since(start).Seconds())
    })
}

// Expose metrics endpoint
r.Handle("/metrics", promhttp.Handler())
```

### Business metrics

```go
// Track business events, not just infrastructure
func (s *OrderService) PlaceOrder(ctx context.Context, in PlaceOrderInput) (*Order, error) {
    start := time.Now()
    defer func() {
        orderProcessingDuration.Observe(time.Since(start).Seconds())
    }()
    
    order, err := s.doPlaceOrder(ctx, in)
    if err != nil {
        ordersTotal.WithLabelValues("failed").Inc()
        return nil, err
    }
    
    ordersTotal.WithLabelValues("success").Inc()
    return order, nil
}
```

---

## 4. OpenTelemetry Tracing

A trace shows the full journey of a request across services. Each step is a **span**. Spans are linked by a **trace ID**.

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
    "go.opentelemetry.io/otel/trace"
)

// Setup tracing
func setupTracing(ctx context.Context, serviceName, endpoint string) (func(), error) {
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(endpoint), // "otel-collector:4317"
        otlptracegrpc.WithInsecure(),
    )
    if err != nil { return nil, err }
    
    res, _ := resource.New(ctx,
        resource.WithAttributes(semconv.ServiceName(serviceName)),
        resource.WithAttributes(semconv.ServiceVersion("1.0.0")),
    )
    
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
        sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.1))), // 10% sampling
    )
    
    otel.SetTracerProvider(tp)
    
    return func() { tp.Shutdown(context.Background()) }, nil
}

// Create spans in your code
var tracer = otel.Tracer("myapp")

func (s *OrderService) PlaceOrder(ctx context.Context, in PlaceOrderInput) (*Order, error) {
    ctx, span := tracer.Start(ctx, "OrderService.PlaceOrder",
        trace.WithAttributes(
            attribute.String("customer_id", in.CustomerID),
            attribute.Int("items_count", len(in.Items)),
        ),
    )
    defer span.End()
    
    // Reserve inventory
    ctx, inventorySpan := tracer.Start(ctx, "inventory.reserve")
    err := s.inventory.Reserve(ctx, in.Items)
    inventorySpan.End()
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    // Charge payment
    ctx, paymentSpan := tracer.Start(ctx, "payment.charge")
    _, err = s.payments.Charge(ctx, in.CustomerID, calculateTotal(in.Items))
    paymentSpan.End()
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    // Persist the order
    order, err := s.repo.Save(ctx, in)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }
    
    span.SetAttributes(attribute.String("order_id", order.ID))
    span.SetStatus(codes.Ok, "")
    return order, nil
}

// HTTP middleware: extract/inject trace context from headers
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

// Automatic instrumentation for HTTP handlers
handler = otelhttp.NewHandler(handler, "http.server")

// Propagate trace context in outbound HTTP calls
client = &http.Client{
    Transport: otelhttp.NewTransport(http.DefaultTransport),
}
```

---

## 5. Connecting the Three

Add trace IDs to logs so you can jump from metrics → logs → trace:

```go
// Add trace ID to every log line in this request
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            span := trace.SpanFromContext(r.Context())
            traceID := span.SpanContext().TraceID().String()
            spanID  := span.SpanContext().SpanID().String()
            
            reqLogger := logger.With(
                "trace_id", traceID,
                "span_id",  spanID,
                "request_id", r.Header.Get("X-Request-ID"),
            )
            ctx := WithLogger(r.Context(), reqLogger)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Exemplar: link a Prometheus metric data point to a specific trace ID
// Grafana can click from a spike in a chart to the trace that caused it
obs := httpRequestDuration.With(prometheus.Labels{"method": "POST", "path": "/orders"})
if eo, ok := obs.(prometheus.ExemplarObserver); ok {
    eo.ObserveWithExemplar(
        time.Since(start).Seconds(),
        prometheus.Labels{"trace_id": traceID},
    )
}

// Note: exemplars are only exposed in the OpenMetrics format —
// enable it on the /metrics endpoint:
// promhttp.HandlerFor(prometheus.DefaultGatherer,
//     promhttp.HandlerOpts{EnableOpenMetrics: true})
```

---

## Summary

- **slog**: structured logging, context-aware; use JSON in production, text in development
- Always log with structured fields (key-value) — avoid unstructured `fmt.Sprintf` messages
- **Prometheus**: counter (ever-increasing totals), gauge (current value), histogram (distribution)
- Avoid high-cardinality labels (`user_id` → millions of label combinations); use `sanitizePath`
- **OpenTelemetry**: create spans for every significant operation; propagate context through HTTP headers
- **Connect them**: add `trace_id` to log lines; use exemplars to link metric data points to traces
- **Sampling**: 100% tracing is too expensive at scale; sample 10-100% based on traffic volume

## Exercises

### Easy
1. Add `slog` to a small HTTP service. Set up JSON logging for production and text for development. Add a middleware that logs every request with method, path, status, and duration.
2. Add three Prometheus metrics to a handler: `requests_total` (counter), `request_duration_seconds` (histogram), and `active_requests` (gauge). Expose them at `/metrics`.
3. Instrument a function with an OpenTelemetry span. Run a local Jaeger instance (`docker run jaegertracing/all-in-one`) and verify the trace appears in the UI.

### Medium
4. Implement a **correlation ID middleware**: generate a UUID for each request, add it to response headers as `X-Request-ID`, and include it in every `slog` log line. Enable log querying by correlation ID in a log aggregator.
5. Build a **Grafana dashboard** (JSON config) with 4 panels: request rate, error rate, P99 latency, and active connections. Write the PromQL queries for each.
6. Add tracing to a service that makes an outbound HTTP call. Verify that the trace shows two spans: the incoming request and the outbound call, with the outbound span as a child.

### Hard
7. Implement **adaptive sampling**: start at 100% sampling rate. If request rate exceeds 1000/s, drop to 10%. If error rate exceeds 5%, boost back to 100% for those error traces. Use a custom `otel.Sampler` implementation.
8. Build a **SLO dashboard**: define Service Level Objectives (99.9% of requests under 500ms, error rate < 0.1%). Track the error budget burn rate. Alert when the budget is consumed at > 5× the expected rate (using Prometheus alert rules).
