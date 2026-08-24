# Chapter 119: Structured Logging with slog, zerolog, and zap

Chapter 118 introduced the three pillars and showed you enough `slog` to get structured JSON out of a service. This chapter is the deep dive on the logging pillar itself: how `slog` actually works under the hood (loggers, records, handlers), how to control levels at runtime, how to redact secrets automatically, and how to write a custom handler that stamps every log line with a request ID. Then we look at the two most popular third-party loggers — zerolog and zap — learn their APIs, and build a decision framework for choosing between the three. We close with the operational topics that separate hobby logging from production logging: sampling, correlation, and the list of things you must never write to a log file.

## Table of Contents

1. [How slog Works — Logger, Record, Handler](#1-how-slog-works--logger-record-handler)
2. [Levels — Built-in, Custom, and Dynamic](#2-levels--built-in-custom-and-dynamic)
3. [Attrs — The Fast Path](#3-attrs--the-fast-path)
4. [Groups and Child Loggers](#4-groups-and-child-loggers)
5. [LogValuer — Lazy Values and Automatic Redaction](#5-logvaluer--lazy-values-and-automatic-redaction)
6. [ReplaceAttr — Rewriting Output](#6-replaceattr--rewriting-output)
7. [Custom Handlers — Context-Aware Logging](#7-custom-handlers--context-aware-logging)
8. [zerolog](#8-zerolog)
9. [zap](#9-zap)
10. [Choosing a Logger](#10-choosing-a-logger)
11. [Log Sampling](#11-log-sampling)
12. [Correlation — Following One Request Across Services](#12-correlation--following-one-request-across-services)
13. [What Not to Log](#13-what-not-to-log)
14. [Summary](#summary)
15. [Exercises](#exercises)

---

## 1. How slog Works — Logger, Record, Handler

`log/slog` splits logging into a **front end** (the API you call) and a **back end** (the handler that formats and writes). Understanding this split explains almost everything else in the package.

```
  your code                      front end                back end
  ─────────                      ─────────                ────────
  logger.Info("msg", ...)  ──►   slog.Logger  ──Record──► slog.Handler ──► io.Writer
                                 (builds a                (formats: JSON,
                                  slog.Record)             text, or custom)
```

- **`slog.Logger`** is what you call. It is cheap to copy and safe for concurrent use.
- **`slog.Record`** is the in-flight log event: time, level, message, and a list of attributes.
- **`slog.Handler`** is an interface. `NewJSONHandler` and `NewTextHandler` are the two built-in implementations, but anyone can write one — that is how third-party backends, test capture, and multi-destination logging plug in.

```go
type Handler interface {
    Enabled(ctx context.Context, level Level) bool  // is this level on?
    Handle(ctx context.Context, r Record) error     // format and write
    WithAttrs(attrs []Attr) Handler                 // returns handler with pre-set attrs
    WithGroup(name string) Handler                  // returns handler with a group prefix
}
```

The front end calls `Enabled` **first** — if the level is off, the record is never even built. That is why a `Debug` call in production costs almost nothing when the level is `Info`.

One design consequence worth internalizing: because the handler is an interface, `slog` is not just a logger — it is a logging *API*. You can keep `slog.Logger` in all of your application code and swap the backend (JSON, zap, zerolog, a test recorder) without touching a single call site.

---

## 2. Levels — Built-in, Custom, and Dynamic

`slog` defines four named levels with numeric gaps between them:

| Level | Value | Use for |
|-------|-------|---------|
| `LevelDebug` | -4 | Developer detail: cache hits, SQL statements, retry attempts |
| `LevelInfo` | 0 | Normal operations: request completed, job finished |
| `LevelWarn` | 4 | Something odd but handled: retry succeeded, deprecated call used |
| `LevelError` | 8 | Something failed and needs attention |

The gaps are deliberate — you can define your own levels in between:

```go
const (
    LevelTrace = slog.Level(-8) // finer than Debug
    LevelFatal = slog.Level(12) // coarser than Error
)

logger.Log(ctx, LevelTrace, "entering handler", "path", r.URL.Path)
```

By default a custom level prints as `DEBUG-4` or `ERROR+4`. Give it a proper name with `ReplaceAttr` (section 6).

### Changing the level at runtime

Hardcoding the level means redeploying to turn on debug logs — exactly when you least want to redeploy. Use a `slog.LevelVar`, which is safe to change from another goroutine:

```go
var programLevel = new(slog.LevelVar) // defaults to Info

func main() {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: programLevel,
    })
    slog.SetDefault(slog.New(handler))

    // Admin endpoint: flip the whole service to debug logging, live
    http.HandleFunc("POST /admin/loglevel", func(w http.ResponseWriter, r *http.Request) {
        var lvl slog.Level
        if err := lvl.UnmarshalText([]byte(r.URL.Query().Get("level"))); err != nil {
            http.Error(w, "bad level (use DEBUG, INFO, WARN, ERROR)", 400)
            return
        }
        programLevel.Set(lvl)
        slog.Info("log level changed", "level", lvl)
    })
}
```

Now `curl -X POST 'localhost:8080/admin/loglevel?level=DEBUG'` turns on debug logs during an incident, and setting it back to `INFO` turns them off — no restart, no redeploy. (Protect the endpoint with auth, obviously.)

---

## 3. Attrs — The Fast Path

Chapter 118 used the "sugar" form — alternating keys and values:

```go
slog.Info("order placed", "order_id", id, "total", total)
```

This is convenient but has two costs: every value is boxed into an `any` (which usually allocates), and a missing key or value is only caught at runtime (`slog` inserts a `!BADKEY` attr instead of crashing). The typed form avoids both:

```go
slog.Info("order placed",
    slog.String("order_id", id),
    slog.Float64("total", total),
    slog.Int("items", len(items)),
    slog.Duration("elapsed", time.Since(start)),
    slog.Time("placed_at", now),
    slog.Bool("gift", isGift),
    slog.Any("address", addr), // fallback for arbitrary types
)
```

`slog.String`, `slog.Int64`, and friends build a `slog.Attr` — a key plus a `slog.Value`, which is a compact union type that stores small values (ints, bools, durations) **without heap allocation**.

For the hottest code paths there is one more step, `LogAttrs`, which skips the varargs-of-`any` conversion entirely:

```go
logger.LogAttrs(ctx, slog.LevelInfo, "cache hit",
    slog.String("key", key),
    slog.Int("size_bytes", n),
)
```

Practical guidance: use the sugar form everywhere by default; switch to typed attrs in middleware and per-request hot paths; reach for `LogAttrs` only when a profiler tells you logging shows up.

---

## 4. Groups and Child Loggers

Two tools keep large log entries organized. **Groups** namespace related fields inside one log call — Chapter 118 showed `slog.Group`. **`WithGroup`** goes further: it prefixes *every subsequent attribute* on a logger:

```go
dbLogger := logger.WithGroup("db")
dbLogger.Info("query executed", "table", "orders", "rows", 42)
```

```json
{"level":"INFO","msg":"query executed","db":{"table":"orders","rows":42}}
```

**Child loggers** via `With` pre-bind fields once instead of repeating them:

```go
// Bind service-wide fields at startup...
logger := slog.New(handler).With(
    "service", "order-service",
    "version", version,
    "env", os.Getenv("ENV"),
)

// ...and request-scoped fields in middleware (as in Chapter 118)
reqLogger := logger.With("request_id", requestID)
```

`With` is efficient by design: the handler pre-formats the bound attrs **once** (that is what the `WithAttrs` method on the handler is for), so a logger with ten bound fields costs the same per log call as one with none.

---

## 5. LogValuer — Lazy Values and Automatic Redaction

`slog.LogValuer` is a small interface with a big payoff:

```go
type LogValuer interface {
    LogValue() slog.Value
}
```

If a logged value implements it, `slog` calls `LogValue()` **only when the record is actually emitted**. That gives you two superpowers.

**1. Lazy expensive values.** The computation is skipped entirely when the level is disabled:

```go
type queryStats struct{ db *sql.DB }

func (q queryStats) LogValue() slog.Value {
    s := q.db.Stats() // only runs if this log line is actually written
    return slog.GroupValue(
        slog.Int("open", s.OpenConnections),
        slog.Int("in_use", s.InUse),
        slog.Int("idle", s.Idle),
    )
}

slog.Debug("pool state", "pool", queryStats{db}) // free when level is Info
```

**2. Redaction at the type level.** Make it *impossible* to accidentally log a secret:

```go
type Password string

func (Password) LogValue() slog.Value {
    return slog.StringValue("[REDACTED]")
}

type User struct {
    ID       int
    Email    string
    Password Password
}

// Control exactly how the whole struct is logged
func (u User) LogValue() slog.Value {
    return slog.GroupValue(
        slog.Int("id", u.ID),
        slog.String("email", maskEmail(u.Email)), // "a***@example.com"
        // Password deliberately omitted
    )
}

slog.Info("user logged in", "user", u)
// {"msg":"user logged in","user":{"id":42,"email":"a***@example.com"}}
```

Now even a careless `slog.Info("debug", "pw", user.Password)` prints `[REDACTED]`. The safety lives in the type, not in every developer's memory. Use this pattern for tokens, API keys, and any PII-carrying domain struct.

---

## 6. ReplaceAttr — Rewriting Output

`HandlerOptions.ReplaceAttr` is a hook that sees every attribute (including the built-in `time`, `level`, `msg`, and `source`) just before it is written. Use it to match your company's log schema, name custom levels, and add a last-line-of-defense redaction filter:

```go
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: programLevel,
    ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
        switch a.Key {
        case slog.TimeKey:
            a.Key = "timestamp" // rename time → timestamp
            a.Value = slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339Nano))
        case slog.MessageKey:
            a.Key = "message"   // rename msg → message
        case slog.LevelKey:
            // Name our custom levels properly
            if lvl := a.Value.Any().(slog.Level); lvl == LevelTrace {
                a.Value = slog.StringValue("TRACE")
            }
        case "authorization", "cookie", "password", "token":
            a.Value = slog.StringValue("[REDACTED]") // belt and suspenders
        }
        return a
    },
})
```

Two caveats: `ReplaceAttr` runs on **every attribute of every log line**, so keep it fast (no regexes over values); and it does not see attrs inside `Group` values unless you recurse into them yourself — which is why `LogValuer` (redaction at the source) is the primary defense and `ReplaceAttr` the backstop.

---

## 7. Custom Handlers — Context-Aware Logging

Chapter 118 stored a request-scoped logger **in** the context and pulled it out in handlers. The inverse pattern is often nicer: keep one global logger, and teach the *handler* to pull values **out of** the context. Notice that `Handle` receives the `context.Context` — that is exactly what it is for.

```go
type ctxKey string

const requestIDKey ctxKey = "request_id"

// ContextHandler wraps another handler and injects request-scoped
// attributes from the context into every record.
type ContextHandler struct {
    slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
    if reqID, ok := ctx.Value(requestIDKey).(string); ok {
        r.AddAttrs(slog.String("request_id", reqID))
    }
    return h.Handler.Handle(ctx, r)
}

// IMPORTANT: re-wrap in WithAttrs/WithGroup, or child loggers created
// with logger.With(...) will silently bypass ContextHandler.
func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return ContextHandler{h.Handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
    return ContextHandler{h.Handler.WithGroup(name)}
}
```

The `WithAttrs`/`WithGroup` overrides are the classic gotcha: if you rely on the embedded handler's methods, they return the **inner** handler type, and every logger derived via `logger.With(...)` loses your wrapper. Always re-wrap.

Wire it up and use the `...Context` logging methods so the context actually reaches the handler:

```go
func main() {
    base := slog.NewJSONHandler(os.Stdout, nil)
    slog.SetDefault(slog.New(ContextHandler{base}))
}

func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = uuid.New().String()
        }
        w.Header().Set("X-Request-ID", id)
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
    // No logger plumbing at all — just pass the context
    slog.InfoContext(r.Context(), "fetching order", "order_id", r.PathValue("id"))
    // {"msg":"fetching order","order_id":"o-7","request_id":"1f3b..."}
}
```

In Chapter 122 we will extend this exact handler to also inject the OpenTelemetry `trace_id` and `span_id` — one wrapper, and every log line in your codebase becomes trace-correlated.

---

## 8. zerolog

[`github.com/rs/zerolog`](https://github.com/rs/zerolog) predates `slog` and remains the fastest mainstream Go logger. Its trademark is a **fluent, chained API** that writes JSON directly with zero allocations:

```go
package main

import (
    "os"
    "time"

    "github.com/rs/zerolog"
)

func main() {
    logger := zerolog.New(os.Stdout).With().
        Timestamp().
        Str("service", "order-service").
        Logger()

    logger.Info().
        Str("order_id", "o-42").
        Int("items", 3).
        Dur("elapsed", 27*time.Millisecond).
        Msg("order placed")
    // {"level":"info","service":"order-service","order_id":"o-42",
    //  "items":3,"elapsed":27,"time":"2026-07-03T10:00:00Z","message":"order placed"}
}
```

Key API points:

```go
// Levels: a global floor plus per-logger levels
zerolog.SetGlobalLevel(zerolog.InfoLevel)
debugLogger := logger.Level(zerolog.DebugLevel)

// Errors get first-class treatment
logger.Error().Err(err).Str("order_id", id).Msg("payment failed")

// Child loggers, same idea as slog's With
reqLogger := logger.With().Str("request_id", reqID).Logger()

// Pretty, colorized output for local development
dev := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).
    With().Timestamp().Logger()

// Context integration: store and retrieve a logger
ctx = reqLogger.WithContext(ctx)
zerolog.Ctx(ctx).Info().Msg("deep in the call stack")
```

Why it is fast: each chained call (`Str`, `Int`, `Dur`) appends bytes straight into a pooled buffer — there is no intermediate `Record`, no reflection, no `any` boxing. The price is the API's one sharp edge: **if you forget the final `.Msg()` (or `.Send()`), nothing is logged at all.** No error, no panic, just silence. Linters (`zerologlint`, bundled in `golangci-lint`) catch this — enable one if you adopt zerolog.

---

## 9. zap

[`go.uber.org/zap`](https://github.com/uber-go/zap) is the most widely deployed structured logger in the Go ecosystem — you will meet it in Kubernetes controllers, etcd, and half the CNCF landscape. It offers **two APIs over one core**:

```go
package main

import (
    "time"

    "go.uber.org/zap"
)

func main() {
    logger, err := zap.NewProduction() // JSON, Info level, sampling on
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // flush buffered logs on shutdown

    // 1) The typed API — fastest, everything is a strongly-typed Field
    logger.Info("order placed",
        zap.String("order_id", "o-42"),
        zap.Int("items", 3),
        zap.Duration("elapsed", 27*time.Millisecond),
    )

    // 2) The sugared API — looser, printf-style and key-value pairs
    sugar := logger.Sugar()
    sugar.Infow("order placed", "order_id", "o-42", "items", 3)
    sugar.Infof("worker %d finished in %s", 4, 27*time.Millisecond)
}
```

`zap.NewDevelopment()` gives human-readable console output with `Debug` enabled. For real deployments, build from a config so the level, encoding, and sampling are explicit:

```go
cfg := zap.NewProductionConfig()
cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // AtomicLevel = runtime-changeable
cfg.EncoderConfig.TimeKey = "timestamp"
cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
cfg.Sampling = &zap.SamplingConfig{Initial: 100, Thereafter: 100} // see §11

logger, err := cfg.Build()

// Child loggers
reqLogger := logger.With(zap.String("request_id", reqID))

// Errors carry the message and (in development mode) a stacktrace
logger.Error("payment failed", zap.Error(err), zap.String("order_id", id))
```

Under the hood sits `zapcore`: an encoder (JSON or console), a syncer (where bytes go), and a level, composable into trees — for example "JSON to stdout at Info **and** everything at Debug to a file" is a five-line `zapcore.NewTee`. And if your team standardizes on the `slog` API, zap ships an official bridge: `go.uber.org/zap/exp/zapslog` exposes any zap logger as a `slog.Handler`, so libraries log via `slog` while zap does the writing.

---

## 10. Choosing a Logger

| | `slog` | `zerolog` | `zap` |
|---|--------|-----------|-------|
| Dependency | none (stdlib) | one small module | module + zapcore |
| Speed (typed API) | good | fastest | very fast |
| Allocations | low (typed attrs) | zero on hot path | zero/near-zero |
| API style | key-value / typed attrs | fluent chain, `.Msg()` required | typed fields + sugared |
| Runtime level change | `LevelVar` | `SetGlobalLevel` / per-logger | `AtomicLevel` |
| Built-in sampling | no (handler's job) | yes (per-logger samplers) | yes (config) |
| Extensibility | `Handler` interface | `Hook` interface | `zapcore` composition |
| Can back the slog API | — | via community bridge | official `zapslog` |

Honest numbers: in typical benchmarks logging ten fields, zerolog and zap land in the ~50–100ns range with zero allocations, while `slog`'s JSON handler is a few times slower with a couple of allocations. That sounds dramatic; in a service doing 1,000 requests/second with five log lines each, the difference is well under 0.1% of your CPU budget. Logger speed almost never decides system performance — I/O and log *volume* do.

So the decision framework:

- **Default to `slog`.** It is in the standard library, it is the API third-party libraries increasingly target, and its handler model means you can swap backends later without a rewrite.
- **Pick zerolog** when logging genuinely is your hot path (proxies, high-frequency event pipelines) or your team likes the fluent style.
- **Pick zap** when you are joining an ecosystem that already uses it, or you need its mature `zapcore` composition (multi-output trees, custom encoders) — and consider fronting it with `zapslog`.

Whichever you choose: choose **one** per service, inject it (or its `slog.Handler`) at startup, and never mix two logging libraries in one binary — you will get two formats, two level configs, and twice the confusion.

---

## 11. Log Sampling

A crash-looping client or a hot retry loop can emit the same log line 50,000 times per second. That floods your log storage (which you pay for per GB — see Chapter 123) and drowns the one line you actually need. **Sampling** keeps a representative subset.

**zerolog** attaches samplers per logger:

```go
// Log every 10th event
sampled := logger.Sample(&zerolog.BasicSampler{N: 10})

// Better: let bursts through, then throttle.
// First 5 events per second pass; after that, 1 in every 100.
sampled = logger.Sample(&zerolog.BurstSampler{
    Burst:       5,
    Period:      time.Second,
    NextSampler: &zerolog.BasicSampler{N: 100},
})

// Sample noisy levels only — errors always get through
prodLogger := logger.Sample(zerolog.LevelSampler{
    DebugSampler: &zerolog.BurstSampler{Burst: 5, Period: time.Second,
        NextSampler: &zerolog.BasicSampler{N: 100}},
    InfoSampler: &zerolog.BasicSampler{N: 10},
    // Warn/Error: no sampler set → never sampled
})
```

**zap** samples by `(level, message)` per second, configured once:

```go
cfg.Sampling = &zap.SamplingConfig{
    Initial:    100, // first 100 identical entries per second pass
    Thereafter: 100, // then 1 of every 100
}
```

Because zap keys on the *message text*, identical repeated lines get squashed while distinct messages flow freely — a good fit for the "same error 50,000 times" failure mode.

**slog** has no built-in sampling; you implement it as a wrapping handler (that is Exercise 7) or sample downstream in your log pipeline.

Rules of thumb: never sample `Error` and above; sample `Info` only on genuinely high-volume lines (per-request logs); and log an occasional counter (`"suppressed": 4900`) so you know sampling happened.

---

## 12. Correlation — Following One Request Across Services

Chapter 118 built the request-ID middleware. The piece that makes it work across a *distributed* system is **propagation**: every outbound call must forward the ID, and every service must prefer an incoming ID over generating a fresh one.

```go
// Outbound: forward the request ID to downstream services
func (c *PaymentClient) Charge(ctx context.Context, amount int64) error {
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/charge", body)
    if err != nil {
        return err
    }
    if reqID, ok := ctx.Value(requestIDKey).(string); ok {
        req.Header.Set("X-Request-ID", reqID) // same ID travels onward
    }
    resp, err := c.http.Do(req)
    // ...
}
```

With the `ContextHandler` from section 7 installed in every service, the flow becomes:

```
 client ──X-Request-ID: 1f3b──►  order-svc ──1f3b──► payment-svc ──1f3b──► bank-gw
              │                       │                    │
              ▼                       ▼                    ▼
        logs: request_id=1f3b   request_id=1f3b      request_id=1f3b
```

One search — `request_id:"1f3b"` in Kibana (Chapter 123) — returns the complete story of that request across every service, in order. Chapter 122 upgrades this idea to full distributed tracing, where the propagated ID also carries parent/child timing structure; when you adopt OpenTelemetry, log the `trace_id` alongside (or instead of) the request ID so logs and traces join on the same key.

---

## 13. What Not to Log

Logs feel private but they are not: they flow to third-party storage, get broad read access ("everyone in engineering"), stick around for 30–90 days, and are routinely pasted into tickets and Slack. Treat every log line as semi-public.

**Never log:**

| Category | Examples | Why |
|----------|----------|-----|
| Credentials | passwords (even wrong ones), API keys, private keys | instant account takeover if leaked |
| Tokens | JWTs, session cookies, OAuth tokens, password-reset links | a logged JWT *is* the session |
| Payment data | full card numbers (PAN), CVV | PCI-DSS violation; last 4 digits only |
| Sensitive PII | government IDs, health data, precise location | GDPR/law; deletion requests can't reach logs |
| Raw payloads | full request/response bodies, `Authorization`/`Cookie` headers | they smuggle in all of the above |

**Be careful with:** email addresses and phone numbers (mask them: `a***@example.com`), user IDs (fine internally, but they make logs subject to data-deletion regulation), and `err.Error()` strings from libraries that echo their input (a URL with `?token=...` in an error message is a classic leak).

**Defense in depth**, in order of reliability:

1. **Types** — `LogValuer` on secrets and domain structs (section 5). Compile-time-adjacent safety; works even when a developer is careless.
2. **Handler backstop** — `ReplaceAttr` redacting known-bad keys (section 6).
3. **Pipeline scrubbing** — regex filters for card-number and JWT patterns in Fluent Bit/Logstash (Chapter 123), catching whatever slipped through.
4. **Review habits** — grep new code for `dump`, `body`, `%+v` on request structs.

And the flip side — do log *enough*: a log line that says `"payment failed"` with no order ID, amount, or error detail is as useless as no log at all. The goal is maximum diagnostic value with zero secret material.

---

## Summary

- `slog` separates the **front end** (`Logger` builds a `Record`) from the **back end** (`Handler` formats and writes) — code against `slog`, and any backend can do the writing.
- Use `slog.LevelVar` (or zap's `AtomicLevel`) to change log levels at runtime via an admin endpoint — no redeploy during incidents.
- Prefer typed attrs (`slog.String`, `slog.Int`) in hot paths; they avoid `any` boxing and catch mistakes the sugar form cannot.
- `LogValuer` gives lazy evaluation and type-level redaction — the most reliable way to keep secrets out of logs.
- `ReplaceAttr` rewrites output at the handler: rename built-in keys, name custom levels, redact known-bad keys.
- A custom wrapping handler can inject `request_id` (and later `trace_id`) from the context into every line — remember to re-wrap in `WithAttrs`/`WithGroup`.
- **zerolog**: fastest, zero-alloc fluent chain — but a forgotten `.Msg()` silently drops the line (use a linter).
- **zap**: typed + sugared APIs, `zapcore` composition, built-in `(level, message)` sampling, official `zapslog` bridge.
- Choose one logger per service; default to `slog` unless benchmarks or ecosystem pull you elsewhere.
- Sample noisy `Debug`/`Info` lines (burst-then-throttle), never sample errors.
- Correlation = generate an ID at the edge, propagate it on every outbound call, stamp it on every log line.
- Never log credentials, tokens, card data, or raw request bodies; enforce with types first, filters second.

---

## Exercises

### Easy

1. Configure a `slog` JSON handler whose output uses `timestamp` (RFC3339, UTC) instead of `time` and `message` instead of `msg`, using `ReplaceAttr`. Verify the output with a few log calls.
2. Define a `LevelTrace = slog.Level(-8)` custom level, make it print as `"TRACE"`, and add a `slog.LevelVar`-backed admin endpoint that can switch the service between `TRACE`, `DEBUG`, and `INFO` at runtime.
3. Write the same "order placed" log line (order ID, item count, duration, error case) three times: with `slog` typed attrs, with zerolog's fluent API, and with zap's typed fields. Compare the JSON each produces.

### Medium

4. Implement `LogValue()` for a `Customer` struct containing `ID`, `Email`, `Phone`, and an embedded `Card` (number, expiry, CVV). The logged form must show the ID, a masked email, a masked phone (`+91******1234`), the card's last 4 digits, and no CVV under any circumstances. Add a test that logs the struct to a buffer and asserts the raw output contains none of the forbidden values.
5. Extend the `ContextHandler` from section 7 to inject both `request_id` and a `tenant_id` from the context, and prove the `WithAttrs` gotcha: comment out your `WithAttrs` override, create a child logger with `logger.With("component", "billing")`, and show that `request_id` disappears from its output — then restore the override and show it come back.
6. Benchmark the three loggers with `go test -bench . -benchmem`: one message with 8 fields, (a) at an enabled level and (b) at a disabled level. Report ns/op and allocs/op for `slog` sugar, `slog.LogAttrs`, zerolog, and zap typed. Explain why the disabled-level numbers are so much closer together.

### Hard

7. Write a `SamplingHandler` that wraps any `slog.Handler` and implements burst-then-throttle sampling per `(level, message)` key: the first N identical records per second pass, then 1 in M. `Error` and above must never be sampled, and every 60 seconds it should emit a synthetic record per suppressed key with a `suppressed_count` attr. Make it concurrency-safe and write a test that hammers it from 10 goroutines.
8. Build a two-service demo of end-to-end correlation: service A receives a request, logs, and calls service B over HTTP; service B logs and returns. Requirements: a shared middleware package that extracts-or-generates `X-Request-ID`, a `ContextHandler` stamping it on every line, an HTTP client wrapper (implement `http.RoundTripper`) that forwards it automatically, and a script that makes one request and then greps both services' stdout to reconstruct the full timeline from the shared ID.
