---
title: HTTP Methods
task: 04
slug: http-methods
concept: r.Method, Method Dispatch
difficulty: beginner
---

## What You Will Build

Write a handler that responds differently based on the HTTP method. REST APIs use methods to express intent: GET reads, POST creates, PUT updates, DELETE removes.

## Function Signature

```go
func ItemsHandler(w http.ResponseWriter, r *http.Request)
```

## Method Dispatch Table

| Method | Status | Body                         |
|--------|--------|------------------------------|
| GET    | 200    | `[]`                         |
| POST   | 201    | echo the request body as-is  |
| DELETE | 204    | (empty body)                 |
| other  | 405    | `{"error":"method not allowed"}` |

For POST: read whatever bytes are in the request body and write them back verbatim as the response body (also set `Content-Type: application/json`).

For DELETE: status 204 means "No Content" — no body is sent.

## Examples

```
GET /items
→ 200 OK, Content-Type: application/json
→ []

POST /items
Body: {"title":"buy milk"}
→ 201 Created, Content-Type: application/json
→ {"title":"buy milk"}

DELETE /items
→ 204 No Content

PATCH /items
→ 405 Method Not Allowed, Content-Type: application/json
→ {"error":"method not allowed"}
```

## Key Concepts

**r.Method** — a string like `"GET"`, `"POST"`, `"DELETE"`. Use a `switch` statement to dispatch:

```go
switch r.Method {
case http.MethodGet:
    // ...
case http.MethodPost:
    // ...
}
```

**http.MethodXxx constants** — use `http.MethodGet`, `http.MethodPost`, etc. instead of raw strings to avoid typos.

**io.Copy for echoing body** — `io.Copy(w, r.Body)` streams the request body directly to the response writer without loading it all into memory.

**204 No Content** — do NOT write a body after `w.WriteHeader(http.StatusNoContent)`. The HTTP spec says 204 responses must not have a body.

## Hints

<details>
<summary>Hint 1 — Switch on r.Method</summary>

```go
switch r.Method {
case http.MethodGet:
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, "[]")
case http.MethodPost:
    // ...
case http.MethodDelete:
    w.WriteHeader(http.StatusNoContent)
default:
    // 405
}
```
</details>

<details>
<summary>Hint 2 — Echoing the POST body</summary>

```go
case http.MethodPost:
    body, _ := io.ReadAll(r.Body)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(body)
```

`io.ReadAll` reads the entire body into a byte slice. You can also use `io.Copy(w, r.Body)` to stream without buffering.
</details>

<details>
<summary>Hint 3 — Returning 405</summary>

The HTTP spec says 405 responses should include an `Allow` header listing valid methods:

```go
default:
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Allow", "GET, POST, DELETE")
    w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprint(w, `{"error":"method not allowed"}`)
```
</details>

## How to Verify

```bash
cd starter/task-04-http-methods
go test ./...
```

The test sends GET, POST (with body), DELETE, and PATCH requests and checks:

- GET → 200 + `[]`
- POST → 201 + echo of body
- DELETE → 204 + empty body
- PATCH → 405 + error JSON
