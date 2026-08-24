# Chapter 124: Sentry — Error Tracking and Performance Monitoring

Production services fail silently. A user hits a 500, closes the tab, and never files a bug report. Your logs might show the error, but logs require someone to be watching. Sentry bridges that gap: it captures errors the moment they happen, enriches them with context, and surfaces them where your team can act.

## Table of Contents

1. [Why Sentry?](#1-why-sentry)
2. [Setup and Initialization](#2-setup-and-initialization)
3. [Capture Errors and Messages](#3-capture-errors-and-messages)
4. [Add Context with Scopes](#4-add-context-with-scopes)
5. [HTTP Middleware](#5-http-middleware)
6. [Breadcrumbs](#6-breadcrumbs)
7. [Performance Monitoring](#7-performance-monitoring)
8. [Grouping and Fingerprinting](#8-grouping-and-fingerprinting)
9. [Alert Rules](#9-alert-rules)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Sentry?

Logs and Sentry serve different purposes.

Logs tell you what the system did. They capture every database query, every cache miss, every request duration. They are comprehensive and chronological. But logs require you to already know something went wrong before you go look.

Sentry tells you what users experienced. When an error occurs, Sentry captures it instantly with a full stack trace, the request that triggered it, the user who was affected, and the chain of events leading up to the failure. It groups repeated occurrences of the same error into a single issue, so you see "this crash happened 847 times affecting 212 users" rather than 847 separate log lines.

The practical difference:

- Logs: you discover an error at 9am when you review last night's output
- Sentry: you get paged at 2am because error rate spiked above threshold

Sentry also tracks which errors are new (never seen before), which are regressions (fixed, then broken again), and how many users each issue affects. These dimensions are not possible to reconstruct from flat log files without significant tooling.

Use both. Logs for forensic investigation, Sentry for real-time alerting and user impact.

---

## 2. Setup and Initialization

Install the SDK:

```bash
go get github.com/getsentry/sentry-go
go get github.com/getsentry/sentry-go/http
```

Initialize Sentry once at program startup, before any handlers run:

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/getsentry/sentry-go"
    sentryhttp "github.com/getsentry/sentry-go/http"
    "github.com/go-chi/chi/v5"
)

func main() {
    err := sentry.Init(sentry.ClientOptions{
        Dsn:              "https://examplePublicKey@o0.ingest.sentry.io/0",
        Environment:      "production",
        Release:          "orders-service@1.4.2",
        TracesSampleRate: 0.2,  // capture 20% of transactions for performance
        SampleRate:       1.0,  // capture 100% of errors
    })
    if err != nil {
        log.Fatalf("sentry.Init: %v", err)
    }
    // Flush buffered events before the program exits.
    // 2 seconds is enough for most network conditions.
    defer sentry.Flush(2 * time.Second)

    r := chi.NewRouter()

    sentryMiddleware := sentryhttp.New(sentryhttp.Options{
        Repanic: true,
    })
    r.Use(sentryMiddleware.Handle)

    r.Get("/orders/{id}", getOrderHandler)

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

Key fields in `ClientOptions`:

| Field | Purpose |
|---|---|
| `Dsn` | Your project's data source name from the Sentry dashboard |
| `Environment` | Separates production, staging, dev in the Sentry UI |
| `Release` | Ties errors to a specific deploy; enables regression detection |
| `SampleRate` | Fraction of error events to send (1.0 = all errors) |
| `TracesSampleRate` | Fraction of requests to trace for performance data |

`sentry.Flush` is critical in short-lived processes (CLIs, Lambda functions). The SDK buffers events and sends them in the background. Without flushing, the program exits before the network call completes and events are lost.

---

## 3. Capture Errors and Messages

Two functions cover most use cases:

```go
sentry.CaptureException(err)    // for Go error values
sentry.CaptureMessage("text")   // for informational alerts without an error
```

A real handler using both:

```go
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "id")

    order, err := db.GetOrder(r.Context(), orderID)
    if err != nil {
        // This is an unexpected failure — capture it as an exception
        // so we get a full stack trace and error details in Sentry.
        sentry.CaptureException(err)
        http.Error(w, "failed to fetch order", http.StatusInternalServerError)
        return
    }

    if order == nil {
        // Not an error, but worth tracking in Sentry as a message.
        // Useful if you're seeing unexpected 404s in production.
        sentry.CaptureMessage("order not found: " + orderID)
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    writeJSON(w, order)
}
```

`CaptureException` extracts the error message. If the error carries a stack trace (as errors from `github.com/pkg/errors` do), Sentry uses it; plain errors from `errors.New` or `fmt.Errorf` (even with `%w`) do not carry one, so Sentry captures the stack at the point where `CaptureException` was called instead. `CaptureMessage` creates a Sentry event without an associated error — useful for tracking unexpected-but-non-fatal conditions.

---

## 4. Add Context with Scopes

A raw error in Sentry shows you what broke. Scope context shows you who was affected and why.

`sentry.WithScope` creates a temporary scope for a single capture. Changes to this scope do not leak to other goroutines or future calls.

```go
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "id")
    userID  := r.Header.Get("X-User-ID")
    email   := r.Header.Get("X-User-Email")

    payload, _ := io.ReadAll(r.Body)

    order, err := db.GetOrder(r.Context(), orderID)
    if err != nil {
        sentry.WithScope(func(scope *sentry.Scope) {
            // Tag the affected user so you can see "this error
            // affected 34 distinct users" in the Sentry UI.
            scope.SetUser(sentry.User{
                ID:    userID,
                Email: email,
            })

            // Tags are indexed — you can filter and search by them.
            scope.SetTag("order_id", orderID)
            scope.SetTag("region", "us-east-1")

            // Extras are not indexed but appear in the event detail.
            // Good for large blobs like request bodies.
            scope.SetExtra("payload", string(payload))

            sentry.CaptureException(err)
        })

        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    writeJSON(w, order)
}
```

When to use tags vs extras:

- Tags: short, discrete values you want to filter by (region, tier, feature flag)
- Extras: large or unstructured data you want visible in the event but not searchable (request body, config snapshot)

The scope is local to the `WithScope` closure. If you need persistent context for all events in a request (user, request ID), set it on the hub attached to the request context instead — the sentryhttp middleware does this automatically when you call `sentry.GetHubFromContext(r.Context())`.

---

## 5. HTTP Middleware

The `sentryhttp` package provides a middleware that:

- Recovers from panics and sends them to Sentry before re-panicking
- Attaches request metadata (URL, method, headers, IP) to every event captured during that request
- Creates a per-request hub so concurrent requests do not share scope state

```go
sentryMiddleware := sentryhttp.New(sentryhttp.Options{
    Repanic:         true,   // re-panic after capturing, so your recovery middleware still runs
    WaitForDelivery: false,  // do not block the response waiting for Sentry
})

r.Use(sentryMiddleware.Handle)
```

`Repanic: true` is almost always what you want. Without it, Sentry catches the panic and swallows it — your other middleware (logging, metrics) never sees it. With `Repanic: true`, Sentry records the panic and then re-panics so the rest of your middleware stack handles it normally.

Request flow with sentryhttp in place:

```
Incoming HTTP Request
        |
        v
+------------------+
|  sentryhttp       |
|  middleware       |  <-- attaches hub to context, wraps handler in recover
+------------------+
        |
        v
+------------------+
|  auth middleware  |
+------------------+
        |
        v
+------------------+
|  your handler    |  <-- calls sentry.CaptureException(err)
|                  |      Sentry reads hub from context
|                  |      attaches request URL, method, headers
+------------------+
        |
   (panic occurs)
        |
        v
+------------------+
|  sentryhttp       |  <-- recover() fires, captures panic to Sentry
|  (deferred)      |      then re-panics (Repanic: true)
+------------------+
        |
        v
+------------------+
|  chi panic       |
|  recovery        |  <-- writes 500 to client
+------------------+
```

Wire it with chi:

```go
r := chi.NewRouter()

sentryMiddleware := sentryhttp.New(sentryhttp.Options{Repanic: true})
r.Use(sentryMiddleware.Handle)

// your other middleware
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)  // chi's built-in recovery runs after sentryhttp re-panics

r.Get("/orders/{id}", getOrderHandler)
r.Post("/orders", createOrderHandler)
```

Inside a handler, access the request-scoped hub to add context that sentryhttp will include automatically:

```go
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
    hub := sentry.GetHubFromContext(r.Context())
    if hub != nil {
        hub.Scope().SetTag("order_id", chi.URLParam(r, "id"))
    }
    // ...
}
```

---

## 6. Breadcrumbs

Breadcrumbs are a trail of events recorded during a request. They are not sent to Sentry on their own — they attach to the next error or message that is captured. When you open an issue in Sentry, you see the breadcrumb trail leading up to it.

```go
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "id")

    sentry.AddBreadcrumb(&sentry.Breadcrumb{
        Category: "handler",
        Message:  "getOrderHandler called",
        Level:    sentry.LevelInfo,
        Data: map[string]any{
            "order_id": orderID,
        },
    })

    order, err := db.GetOrder(r.Context(), orderID)
    if err != nil {
        sentry.AddBreadcrumb(&sentry.Breadcrumb{
            Category: "db",
            Message:  "GetOrder failed",
            Level:    sentry.LevelError,
            Data: map[string]any{
                "order_id": orderID,
                "error":    err.Error(),
            },
        })

        sentry.CaptureException(err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    sentry.AddBreadcrumb(&sentry.Breadcrumb{
        Category: "db",
        Message:  "GetOrder succeeded",
        Level:    sentry.LevelInfo,
    })

    fulfillment, err := fulfillmentClient.Check(r.Context(), order.ID)
    if err != nil {
        // When this error is captured, Sentry shows the full breadcrumb
        // trail: handler called, db query succeeded, fulfillment check failed.
        sentry.CaptureException(err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    writeJSON(w, mergeOrderFulfillment(order, fulfillment))
}
```

Breadcrumbs vs logs:

| | Breadcrumbs | Logs |
|---|---|---|
| Always emitted | No | Yes |
| Attach to error | Yes | No |
| Searchable | No (part of event) | Yes (log aggregator) |
| Useful for | Reconstructing error context | Auditing, debugging |

Use breadcrumbs for events that only matter when something goes wrong. Use logs for events you always need (audit trails, compliance, metrics).

---

## 7. Performance Monitoring

Sentry performance monitoring tracks how long operations take in production. It uses a transaction/span model: a transaction is one user-visible operation (an HTTP request, a background job), and spans are the sub-operations within it (database queries, external HTTP calls, cache lookups).

`sentryhttp` creates a transaction automatically for each HTTP request when `TracesSampleRate > 0`. You add child spans for internal operations:

```go
package db

import (
    "context"
    "database/sql"

    "github.com/getsentry/sentry-go"
)

func GetOrder(ctx context.Context, orderID string) (*Order, error) {
    // Start a span as a child of whatever transaction is in ctx.
    // If there is no transaction (TracesSampleRate caused this request
    // to be unsampled), StartSpan is a no-op.
    span := sentry.StartSpan(ctx, "db.query",
        sentry.WithDescription("SELECT orders WHERE id = ?"),
    )
    defer span.Finish()

    var order Order
    err := db.QueryRowContext(span.Context(), // pass span's context forward
        "SELECT id, user_id, status, total FROM orders WHERE id = $1",
        orderID,
    ).Scan(&order.ID, &order.UserID, &order.Status, &order.Total)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        span.Status = sentry.SpanStatusInternalError
        return nil, err
    }

    span.Status = sentry.SpanStatusOK
    return &order, nil
}
```

A handler that creates its own transaction with multiple child spans:

```go
func processOrderHandler(w http.ResponseWriter, r *http.Request) {
    // sentryhttp creates the root transaction; StartSpan attaches to it.
    span := sentry.StartSpan(r.Context(), "order.process")
    defer span.Finish()
    ctx := span.Context()

    // Child span for the DB read.
    order, err := db.GetOrder(ctx, chi.URLParam(r, "id"))
    if err != nil {
        sentry.CaptureException(err)
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }

    // Child span for the downstream service call.
    fulfillSpan := sentry.StartSpan(ctx, "http.client",
        sentry.WithDescription("GET fulfillment-service/check"),
    )
    fulfillment, err := fulfillmentClient.Check(ctx, order.ID)
    fulfillSpan.Finish()

    if err != nil {
        sentry.CaptureException(err)
        http.Error(w, "error", http.StatusInternalServerError)
        return
    }

    writeJSON(w, mergeOrderFulfillment(order, fulfillment))
}
```

In the Sentry UI, performance monitoring shows:

- P50, P75, P95, P99 latency per endpoint
- Which spans are slowest (database? external calls?)
- Transaction traces for individual slow requests
- Apdex score over time

`TracesSampleRate` at 1.0 traces every request. In production, 0.1–0.2 (10–20%) gives useful data without excessive overhead or cost.

---

## 8. Grouping and Fingerprinting

Sentry automatically groups events into issues. The default algorithm uses the error type, message, and stack trace. Two occurrences of the same `sql: no rows` error at the same call site become one issue with an occurrence count.

Sometimes the default grouping is wrong. The same root cause manifests at different call sites, creating dozens of separate issues. Or different root causes have identical stack traces. Fix this with fingerprinting.

```go
// All database errors for a specific table go into one issue,
// regardless of which function triggered them.
func captureDBError(err error, table string) {
    sentry.WithScope(func(scope *sentry.Scope) {
        scope.SetFingerprint([]string{"database-error", table})
        scope.SetTag("table", table)
        sentry.CaptureException(err)
    })
}

// Force all payment failures into one issue.
func capturePaymentError(err error, provider string) {
    sentry.WithScope(func(scope *sentry.Scope) {
        scope.SetFingerprint([]string{"payment-failure", provider})
        scope.SetTag("payment_provider", provider)
        sentry.CaptureException(err)
    })
}
```

Fingerprint values are arbitrary strings. Sentry groups events with identical fingerprint arrays into the same issue. Use specific values when you want fine-grained grouping, and coarse values when you want aggregation.

When to override fingerprinting:

- A generic error message ("context deadline exceeded") occurs in many places that are actually one systemic problem
- A noisy third-party library emits errors from varying call stacks for the same root cause
- You want to track "any payment failure" as one business-level issue across multiple providers

---

## 9. Alert Rules

After errors flow into Sentry, configure alert rules in the Sentry dashboard to route notifications to your team.

**Issue alerts** fire on event conditions:

- New issue: triggers once when an error type is seen for the first time. Useful for catching new bugs introduced by a deploy.
- Regression: triggers when an issue marked resolved starts occurring again. Catches when a fix breaks in a new release.
- Error rate spike: triggers when events per minute for an issue crosses a threshold.

**Metric alerts** fire on aggregate conditions:

- Error rate for a specific transaction exceeds 1%
- P95 latency for `/orders/{id}` exceeds 500ms

**Integration targets:**

- Slack: post alert to a channel, mention on-call
- PagerDuty: create an incident, page on-call engineer
- Email: send to team DL
- GitHub/Jira: auto-create issues from Sentry events

**Issue owners** map code ownership to Sentry alerts. Define `CODEOWNERS` in your repository and link it in Sentry project settings. When an error occurs in `internal/payments/`, Sentry knows to route it to the payments team rather than broadcasting to everyone.

Configure alert rules per environment. A new issue in staging sends a Slack message. The same issue in production pages PagerDuty.

---

## Summary

Sentry gives you visibility into what users experience in production — not what the system logs, but what actually broke for real people.

- `sentry.Init` with `Dsn`, `Environment`, `Release`, and `TracesSampleRate` bootstraps the SDK. Call `defer sentry.Flush(2*time.Second)` before your server starts listening.
- `sentry.CaptureException(err)` sends errors with stack traces. `sentry.CaptureMessage` sends informational events.
- `sentry.WithScope` attaches user identity, tags, and extra data to a single capture without polluting other events.
- `sentryhttp.New(sentryhttp.Options{Repanic: true}).Handle` auto-captures panics and attaches HTTP request context to every event.
- `sentry.AddBreadcrumb` leaves a trail that attaches to the next captured event, giving you the sequence of operations before the failure.
- `sentry.StartSpan(ctx, "operation")` creates performance spans. Child spans nest under the root transaction created by sentryhttp. `TracesSampleRate` controls sampling overhead.
- `scope.SetFingerprint` overrides Sentry's default grouping. Use it when the same root cause appears under different stack traces, or when you want to merge related issues.
- Alert rules in the Sentry UI route notifications to Slack, PagerDuty, or email. Use issue owners to route alerts to the right team automatically.

---

## Exercises

### Easy

Initialize Sentry in a small HTTP server. The server has one endpoint: `GET /test`. When that endpoint is hit, call `sentry.CaptureMessage("test endpoint hit")` and return `200 OK`. Verify the message appears in your Sentry project's Issues list.

Requirements:
- `sentry.Init` with a real DSN from a Sentry project you create
- `defer sentry.Flush(2 * time.Second)` before the server starts
- The `/test` handler calls `sentry.CaptureMessage` with any descriptive string
- Run the server, hit the endpoint with `curl`, confirm the event in the Sentry dashboard

### Medium

Add `sentryhttp` middleware to a chi router. Create a handler that panics intentionally:

```go
r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
    panic("intentional test panic")
})
```

Also add chi's `middleware.Recoverer` after the sentryhttp middleware. Hit `/panic` and verify:

1. The server does not crash (Recoverer handles the re-panic)
2. The event appears in Sentry with the full stack trace
3. The Sentry event includes the request URL, method, and headers
4. The server returns a 500 response to the client

Explain in a comment why `Repanic: true` is necessary for both the recovery middleware and Sentry to work correctly.

### Hard

Implement a full observability middleware for a chi router that handles an order processing endpoint. The middleware and handler together must:

1. Create a Sentry transaction for each request using `sentry.StartSpan(r.Context(), "http.request")`
2. For every SQL query executed during the request, add a breadcrumb with `Category: "db"`, the query string, and execution duration in the `Data` map
3. For every SQL query, also create a child span under the request transaction
4. If the handler returns an error, capture it with `sentry.WithScope` that includes:
   - User ID and email from request headers
   - Order ID from the URL parameter as a tag
   - The raw request body as an extra field
5. At the end of the request, finish the root transaction span

Structure your solution as:
- A `ObservabilityMiddleware` chi middleware that sets up the transaction and defers span finishing
- A `TracedDB` wrapper around `*sql.DB` that adds breadcrumbs and spans for each query
- An `OrderHandler` struct that uses `TracedDB` and calls `captureWithUserContext` on error

The goal is that every error in Sentry has: the user, the order, the request body, the breadcrumb trail of DB queries, and the performance trace showing how long each DB call took.
