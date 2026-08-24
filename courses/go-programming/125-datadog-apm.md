# Chapter 125: Datadog — APM, Infrastructure, and Logs

Datadog is an all-in-one observability platform: metrics, distributed traces, logs, infrastructure monitoring, and real user monitoring (RUM) all live under one roof and are correlated by design. This chapter covers how to instrument a Go service with Datadog APM, emit custom metrics via DogStatsD, and wire logs to traces so you can debug a slow production request from any starting point.

## Table of Contents

1. [Datadog vs DIY Observability Stack](#1-datadog-vs-diy-observability-stack)
2. [Tracer Setup](#2-tracer-setup)
3. [Instrument HTTP with dd-trace-go](#3-instrument-http-with-dd-trace-go)
4. [Create Custom Spans](#4-create-custom-spans)
5. [Instrument gRPC](#5-instrument-grpc)
6. [Datadog Metrics via DogStatsD](#6-datadog-metrics-via-dogstatsd)
7. [Unified Service Tagging](#7-unified-service-tagging)
8. [Log Correlation](#8-log-correlation)
9. [APM Features Overview](#9-apm-features-overview)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Datadog vs DIY Observability Stack

Both approaches cover the same three pillars — metrics, traces, logs — but the trade-offs differ.

```
DIY Stack                          Datadog
─────────────────────              ──────────────────────────────
Prometheus  → Grafana              Metrics  ─┐
OpenTelemetry → Jaeger/Tempo       Traces   ──┤  One UI, one agent,
slog → Loki / OpenSearch           Logs     ──┤  one billing relationship
                                   Infra    ─┘
You wire the correlations.         Correlations built in.
Free. More operational work.       Paid. Less operational work.
```

| Concern | DIY Stack | Datadog |
|---------|-----------|---------|
| Cost | Free (hosting costs) | Per-host or per-GB pricing |
| Setup time | Days to weeks | Hours |
| Trace-to-log linking | Manual (trace ID in log fields + Grafana config) | Automatic when using dd-trace-go + Datadog Agent |
| Vendor lock-in | None (OpenTelemetry is vendor-neutral) | High |
| Flame graphs, service maps | Jaeger/Grafana Tempo | Built-in APM |

`dd-trace-go` is Datadog's Go tracing library. It instruments HTTP, gRPC, database drivers, Redis clients, and more via integrations packages under `gopkg.in/DataDog/dd-trace-go.v1/contrib/`.

---

## 2. Tracer Setup

Install the tracer:

```bash
go get gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer
```

Start the tracer in `main.go` before your server starts:

```go
package main

import (
    "os"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
    tracer.Start(
        tracer.WithService("order-service"),
        tracer.WithEnv("production"),
        tracer.WithServiceVersion("1.4.2"),
    )
    defer tracer.Stop()

    // ... start HTTP server
}
```

`tracer.Stop()` flushes any in-flight spans to the Datadog Agent before the process exits.

**Environment variable alternative.** Instead of hardcoding values in code, set them on the process:

```bash
DD_SERVICE=order-service
DD_ENV=production
DD_VERSION=1.4.2
DD_AGENT_HOST=datadog-agent   # hostname of the Datadog Agent sidecar
```

When these variables are set, `tracer.Start()` with no options picks them up automatically. The env-var approach is preferred in Kubernetes because values come from the pod spec, not from code.

The Datadog Agent runs as a DaemonSet (or sidecar) and receives trace data from your application on port 8126. Your app never talks to the Datadog backend directly.

```
App (dd-trace-go)  --port 8126-->  Datadog Agent  --HTTPS-->  Datadog backend
```

---

## 3. Instrument HTTP with dd-trace-go

Two approaches depending on your router.

**Standard library: `ddhttp.WrapHandler`**

```bash
go get gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http
```

```go
import (
    ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/orders", handleOrders)

    // WrapHandler wraps the entire mux — every request gets a trace span.
    wrapped := ddhttp.WrapHandler(mux, "order-service", "/")

    http.ListenAndServe(":8080", wrapped)
}
```

**Chi router: `chitrace.Middleware`**

```bash
go get gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5
```

```go
import (
    "github.com/go-chi/chi/v5"
    chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func newRouter() *chi.Mux {
    r := chi.NewRouter()

    // Add Datadog tracing middleware before your routes.
    r.Use(chitrace.Middleware(chitrace.WithServiceName("order-service")))

    r.Post("/orders", handleCreateOrder)
    r.Get("/orders/{id}", handleGetOrder)

    return r
}
```

The middleware extracts trace context from incoming headers (e.g. `X-Datadog-Trace-Id`), so distributed traces stitch together correctly when a load balancer or upstream service propagates headers.

```
Load Balancer
     |
     |  X-Datadog-Trace-Id: 9823741234
     |  X-Datadog-Parent-Id: 1234567890
     v
[order-service]  <-- chitrace.Middleware extracts context here
     |
     |  continues the same trace
     v
[inventory-service]  <-- propagated via outgoing HTTP call
```

---

## 4. Create Custom Spans

A span represents a unit of work. Wrap any operation you want to measure — a database call, a cache lookup, a third-party API call — in a span.

```go
import (
    "context"
    "database/sql"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// GetOrder queries the database. It accepts ctx so the span appears as a child
// of the incoming HTTP trace.
func GetOrder(ctx context.Context, db *sql.DB, orderID string) (*Order, error) {
    span, ctx := tracer.StartSpanFromContext(ctx, "db.query",
        tracer.ResourceName("SELECT orders"),
    )
    defer span.Finish()

    query := "SELECT id, user_id, total FROM orders WHERE id = $1"
    span.SetTag(ext.DBType, "postgresql")
    span.SetTag(ext.DBStatement, query)
    span.SetTag("order.id", orderID)

    var o Order
    err := db.QueryRowContext(ctx, query, orderID).Scan(&o.ID, &o.UserID, &o.Total)
    if err != nil {
        // Mark the span as an error — visible in APM error tracking.
        span.SetTag(ext.Error, err)
        return nil, err
    }

    return &o, nil
}
```

Key points:

- `tracer.StartSpanFromContext` creates a child span of whatever span is already in `ctx`. This is how spans nest into a trace tree.
- The second return value is a new `ctx` with the child span injected. Pass this new `ctx` to any further calls.
- `defer span.Finish()` records the end time and sends the span to the Agent.
- `ext.Error` on a span marks it red in the APM UI and surfaces it in error tracking with the full stack trace.

```
HTTP span (chitrace.Middleware)
  └── db.query span (GetOrder)
        └── (future child span if GetOrder calls another service)
```

---

## 5. Instrument gRPC

```bash
go get gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc
```

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

func newGRPCServer() *grpc.Server {
    // UnaryServerInterceptor creates a span for each incoming unary RPC call.
    si := grpctrace.UnaryServerInterceptor(
        grpctrace.WithServiceName("order-grpc"),
    )

    server := grpc.NewServer(
        grpc.UnaryInterceptor(si),
    )

    // Register your services here.
    pb.RegisterOrderServiceServer(server, &OrderServiceServer{})

    return server
}

// On the client side, propagate the trace to downstream gRPC services.
func newGRPCClient(target string) (*grpc.ClientConn, error) {
    ci := grpctrace.UnaryClientInterceptor(
        grpctrace.WithServiceName("order-grpc-client"),
    )

    return grpc.Dial(target,
        grpc.WithUnaryInterceptor(ci),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
}
```

The server interceptor extracts trace context from gRPC metadata on incoming calls. The client interceptor injects trace context into outgoing metadata. Together they stitch a single distributed trace across gRPC boundaries.

---

## 6. Datadog Metrics via DogStatsD

APM traces are great for latency analysis, but counters and gauges are cheaper to query at scale. Use DogStatsD for custom business metrics.

```bash
go get github.com/DataDog/datadog-go/v5/statsd
```

```go
import (
    "time"

    "github.com/DataDog/datadog-go/v5/statsd"
)

// NewStatsdClient connects to the Datadog Agent's StatsD port.
func NewStatsdClient() (*statsd.Client, error) {
    return statsd.New("localhost:8125",
        statsd.WithNamespace("order_service."),  // prefix for all metrics
    )
}

func PlaceOrder(ctx context.Context, sc *statsd.Client, order Order) error {
    tags := []string{
        "env:prod",
        "service:order-service",
        "region:us-east-1",
    }

    start := time.Now()

    if err := processPayment(ctx, order); err != nil {
        sc.Incr("orders.failed", tags, 1)
        return err
    }

    // Increment: count of successfully placed orders.
    sc.Incr("orders.placed", tags, 1)

    // Distribution: payment processing duration in milliseconds.
    // Use Distribution (not Histogram) for accurate p99 percentiles at scale.
    sc.Distribution("orders.payment_duration_ms",
        float64(time.Since(start).Milliseconds()),
        tags, 1,
    )

    // Gauge: current queue depth (a point-in-time value).
    queueDepth := float64(getQueueDepth())
    sc.Gauge("orders.queue_depth", queueDepth, tags, 1)

    return nil
}
```

**Distribution vs Histogram**

| Type | Aggregated where | Accurate percentiles at scale |
|------|-----------------|-------------------------------|
| Histogram | Client-side, then flushed | No — pre-bucketed before the Agent sees them |
| Distribution | Agent-side across all sources | Yes — Agent aggregates raw values from all pods |

At high throughput with many pods, use `Distribution` for latency metrics. It gives accurate p50/p95/p99 because the Agent merges samples from all instances before computing percentiles.

---

## 7. Unified Service Tagging

Unified Service Tagging is the practice of setting `DD_SERVICE`, `DD_ENV`, and `DD_VERSION` on every workload. When all three are consistent, Datadog links metrics, traces, and logs automatically. You can click from a spike in a latency metric directly to the deployment version that caused it.

**Kubernetes Deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
spec:
  template:
    spec:
      containers:
        - name: order-service
          image: order-service:1.4.2
          env:
            - name: DD_SERVICE
              value: "order-service"
            - name: DD_ENV
              value: "production"
            - name: DD_VERSION
              value: "1.4.2"
            - name: DD_AGENT_HOST
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP   # DaemonSet agent on the node
```

Set `DD_VERSION` to your image tag or Git SHA. When you deploy `1.4.3` and latency spikes, APM shows the version tag on every affected trace, and you can compare p99 before and after the deploy in a single graph.

---

## 8. Log Correlation

To link a log entry to its trace, add `dd.trace_id` and `dd.span_id` fields to every structured log line. Datadog's log pipeline reads these fields and creates a clickable link from the log to the APM trace.

```go
import (
    "context"
    "log/slog"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// traceAttrs extracts Datadog trace and span IDs from the context.
// Returns slog attributes you can attach to any log record.
func traceAttrs(ctx context.Context) []any {
    span, ok := tracer.SpanFromContext(ctx)
    if !ok {
        return nil
    }
    return []any{
        slog.Uint64("dd.trace_id", span.Context().TraceID()),
        slog.Uint64("dd.span_id", span.Context().SpanID()),
        slog.String("dd.service", "order-service"),
        slog.String("dd.env", "production"),
        slog.String("dd.version", "1.4.2"),
    }
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    order, err := parseOrder(r)
    if err != nil {
        attrs := traceAttrs(ctx)
        slog.ErrorContext(ctx, "failed to parse order body",
            append(attrs, slog.String("error", err.Error()))...,
        )
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    slog.InfoContext(ctx, "order received",
        append(traceAttrs(ctx),
            slog.String("order_id", order.ID),
            slog.String("user_id", order.UserID),
        )...,
    )

    // ... process order
}
```

A correlated log line in JSON looks like:

```json
{
  "time": "2026-06-30T10:14:33Z",
  "level": "INFO",
  "msg": "order received",
  "dd.trace_id": 9823741234567890,
  "dd.span_id": 1234567890123456,
  "dd.service": "order-service",
  "dd.env": "production",
  "dd.version": "1.4.2",
  "order_id": "ord_8821",
  "user_id": "usr_441"
}
```

In the Datadog Log Explorer, clicking the trace link on this log opens the full APM trace — flame graph, child spans, errors — for the exact request that generated the log.

---

## 9. APM Features Overview

These features require no additional instrumentation beyond what is covered above.

**Flame graphs.** APM renders every trace as a flame graph showing time spent in each span. Useful for finding which database query or downstream call accounts for most of a request's latency.

**Service maps.** Datadog automatically maps service dependencies from trace data. The map shows which services call which, request rates, error rates, and latency between edges. Useful for understanding blast radius before a deployment.

**Error tracking.** Spans tagged with `ext.Error` are grouped by error type and stack trace fingerprint. APM shows the first and last occurrence, affected versions, and impacted users — without requiring you to search logs manually.

**Watchdog.** Datadog's anomaly detection engine (Watchdog) monitors your APM data automatically. It surfaces unexpected spikes in error rate or latency, regressions after a new deployment, and unusual patterns — without requiring you to write alert rules.

---

## Summary

- Datadog is a paid all-in-one platform. The DIY alternative (Prometheus + Grafana + Jaeger + Loki) is free but requires more operational work to wire correlations together.
- `tracer.Start()` in `main.go` and `defer tracer.Stop()` are the minimum required to enable APM. Use `DD_SERVICE`, `DD_ENV`, `DD_VERSION` env vars instead of hardcoding options.
- Use `ddhttp.WrapHandler` or `chitrace.Middleware` to auto-instrument HTTP. Use `grpctrace.UnaryServerInterceptor` for gRPC.
- `tracer.StartSpanFromContext` creates child spans. Always pass the returned `ctx` downstream. Call `defer span.Finish()`. Tag errors with `span.SetTag(ext.Error, err)`.
- DogStatsD (`statsd.New`) sends custom business metrics. Use `Distribution` instead of `Histogram` for accurate p99 percentiles across many pods. Always tag with `env`, `service`, and `version`.
- Unified Service Tagging links metrics, traces, and logs. Set `DD_SERVICE`, `DD_ENV`, `DD_VERSION` on every Kubernetes pod.
- Add `dd.trace_id` and `dd.span_id` to every structured log line. Extract them from the span via `tracer.SpanFromContext(ctx)`. This makes every log in Datadog clickable to its trace.

---

## Exercises

### Easy

Start a Datadog tracer in `main.go`. Wrap a single `GET /health` HTTP handler with `ddhttp.WrapHandler`. Run the server and confirm it starts without errors. Check that `tracer.Start()` is called before the server starts and `tracer.Stop()` is deferred.

### Medium

Add DogStatsD metrics to an order placement flow. Create a `statsd.Client` connected to `localhost:8125`. In a `PlaceOrder` function:

1. Increment an `orders.placed` counter with tags `env:staging` and `service:order-service` on success.
2. Record the end-to-end payment processing time as a `Distribution` metric named `orders.payment_duration_ms`.
3. On failure, increment an `orders.failed` counter with the same tags plus an `error_type` tag set to the error kind (e.g. `error_type:payment_declined`).

### Hard

Build a fully correlated service. Wire up all three signal types so any slow request is debuggable from any starting point:

1. HTTP handler uses `chitrace.Middleware` — every request has a root span.
2. The handler calls `GetOrder(ctx, db, id)` which creates a child span named `db.query` via `tracer.StartSpanFromContext`. The child span is tagged with the SQL statement and marks an error if the query fails.
3. After the DB call, the handler logs the result using `slog.InfoContext` with `dd.trace_id` and `dd.span_id` fields extracted from the context span.
4. The handler emits a `Distribution` metric named `orders.get_duration_ms` tagged with `env`, `service`, and `status` (either `hit` or `miss`).

The end state: a p99 latency spike appears in the metric, you click to a slow trace in APM, you see the flame graph showing the DB span is slow, and you click to the correlated log that shows the exact order ID and user ID involved.
