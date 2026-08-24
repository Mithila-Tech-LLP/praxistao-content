---
title: Middleware
task: 07
slug: middleware
concept: http.Handler wrapping, Middleware chains
difficulty: intermediate
---

## What You Will Build

Write two middleware functions that wrap any `http.Handler`. Middleware is how you add cross-cutting concerns — logging, authentication, rate limiting, CORS — to every route without duplicating code in each handler.

## Package-Level State

```go
var Logs []string
```

`LoggingMiddleware` appends a log entry to this slice after each request.

## Function Signatures

```go
func LoggingMiddleware(next http.Handler) http.Handler

func AuthMiddleware(next http.Handler) http.Handler
```

## LoggingMiddleware

- Wraps `next`
- After calling `next.ServeHTTP(...)`, appends to `Logs` in the format:
  ```
  "GET /path 200"
  ```
  (method, space, path, space, status code)
- Uses a `responseRecorder` to capture the status code written by `next`

## AuthMiddleware

- Reads the `X-API-Key` request header
- If the value is not exactly `"secret"`, writes `401 Unauthorized` with body `{"error":"unauthorized"}` and returns — do NOT call `next`
- If the key is correct, calls `next.ServeHTTP(w, r)` normally

## Middleware Chain

Middlewares are composable. Wrapping from outside in:

```go
handler := LoggingMiddleware(AuthMiddleware(myHandler))
```

Request flows: LoggingMiddleware → AuthMiddleware → myHandler → back.

## Key Concepts

**The middleware pattern** — a function that takes an `http.Handler` and returns a new `http.Handler`. The returned handler does something before and/or after calling `next.ServeHTTP(w, r)`.

```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        next.ServeHTTP(w, r)
        // after
    })
}
```

**Capturing status code** — `http.ResponseWriter` doesn't expose the written status code. Wrap it with your own struct that records it:

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}
func (sr *statusRecorder) WriteHeader(code int) {
    sr.status = code
    sr.ResponseWriter.WriteHeader(code)
}
```

## Hints

<details>
<summary>Hint 1 — statusRecorder for LoggingMiddleware</summary>

```go
type statusRecorder struct {
    http.ResponseWriter
    status int
}
func (sr *statusRecorder) WriteHeader(code int) {
    sr.status = code
    sr.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
        next.ServeHTTP(rec, r)
        Logs = append(Logs, fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, rec.status))
    })
}
```
</details>

<details>
<summary>Hint 2 — AuthMiddleware</summary>

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-API-Key") != "secret" {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprint(w, `{"error":"unauthorized"}`)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```
</details>

<details>
<summary>Hint 3 — Why status defaults to 200</summary>

If a handler calls `w.Write(...)` without calling `w.WriteHeader(...)` first, Go implicitly sends status 200. Your `statusRecorder` initialises `status: http.StatusOK` to handle this case — `WriteHeader` only gets called explicitly for non-200 responses.
</details>

## How to Verify

```bash
cd starter/task-07-middleware
go test ./...
```

The test:

1. Wraps a simple `200` handler with both middlewares
2. Sends a request with the correct API key — expects `200` and a log entry `"GET /test 200"`
3. Resets `Logs`, sends a request with a wrong key — expects `401` and log entry `"GET /test 401"`
4. Sends a request with no key — expects `401`
