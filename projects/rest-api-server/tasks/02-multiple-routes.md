---
title: Multiple Routes
task: 02
slug: multiple-routes
concept: HTTP Routing, ServeMux
difficulty: beginner
---

## What You Will Build

Build a router that maps three URL paths to three different JSON responses. This is how a real API separates concerns — each path has its own handler logic.

## Function Signature

```go
func NewRouter() http.Handler
```

Return an `http.Handler` (typically a configured `*http.ServeMux`) with all routes registered.

## Routes to Register

| Method | Path       | Response body             | Status |
|--------|-----------|--------------------------|--------|
| GET    | /hello    | `{"message":"hello"}`    | 200    |
| GET    | /ping     | `{"status":"ok"}`        | 200    |
| GET    | /version  | `{"version":"1.0"}`      | 200    |

All three routes must set `Content-Type: application/json`.

## Example

```
GET /ping HTTP/1.1
→ 200 OK
→ Content-Type: application/json
→ {"status":"ok"}

GET /version HTTP/1.1
→ 200 OK
→ Content-Type: application/json
→ {"version":"1.0"}
```

## Key Concepts

**http.ServeMux** — Go's built-in request multiplexer (router). You register patterns and handler functions with `mux.HandleFunc(pattern, handler)`, then pass the mux to `http.ListenAndServe`.

**Pattern matching** — `ServeMux` matches the most specific pattern. `/hello` matches exactly `/hello`. A trailing slash like `/api/` matches any path under `/api/`.

**Inline handlers** — you can register an anonymous function directly:

```go
mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    // handle /ping
})
```

## Hints

<details>
<summary>Hint 1 — Creating a ServeMux</summary>

```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", helloHandler)
mux.HandleFunc("/ping", pingHandler)
// ...
return mux
```
</details>

<details>
<summary>Hint 2 — Returning JSON from each handler</summary>

Each handler just needs to set the header and write a fixed string:

```go
func pingHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status":"ok"}`)
}
```
</details>

<details>
<summary>Hint 3 — Testing routes with httptest</summary>

The test calls your router like this:

```go
router := NewRouter()
req := httptest.NewRequest("GET", "/ping", nil)
rr  := httptest.NewRecorder()
router.ServeHTTP(rr, req)
// check rr.Code and rr.Body.String()
```

`ServeHTTP` dispatches the request to the correct handler just like a real server would.
</details>

## How to Verify

```bash
cd starter/task-02-multiple-routes
go test ./...
```

The test sends a request to each of the three routes and checks:

- Status code is `200` for each
- `Content-Type` is `application/json` for each
- Body matches the expected JSON string for each route
