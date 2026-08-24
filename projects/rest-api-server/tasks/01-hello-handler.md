---
title: Hello Handler
task: 01
slug: hello-handler
concept: net/http, http.HandlerFunc, ResponseWriter
difficulty: beginner
---

## What You Will Build

Write your first HTTP handler — a function that responds to any incoming request with a JSON object containing a greeting message.

This is the foundation of every web API: a function that receives a request and writes a response.

## Function Signature

```go
func HelloHandler(w http.ResponseWriter, r *http.Request)
```

`w` is the `ResponseWriter` — you write your response into it.  
`r` is the incoming `*Request` — you can read headers, the URL, body, etc. from it.

## What It Should Do

1. Set the `Content-Type` response header to `application/json`
2. Write HTTP status `200 OK`
3. Write the JSON body `{"message":"hello"}`

## Example

```
GET /hello HTTP/1.1
→ 200 OK
→ Content-Type: application/json
→ {"message":"hello"}
```

## Key Concepts

**ResponseWriter** — an interface with three methods: `Header()` to set headers, `WriteHeader(status)` to send the status code, and `Write([]byte)` to send the body. You must call `Header()` before `WriteHeader`, and `WriteHeader` before `Write`.

**Content-Type header** — tells the client how to interpret the body. For JSON APIs, always set `Content-Type: application/json`.

**json.Marshal vs fmt.Fprintf** — you can write `{"message":"hello"}` as a literal string with `fmt.Fprintf`, or encode a struct/map with `json.NewEncoder(w).Encode(...)`. Both work; encoding a struct is safer for complex data.

## Hints

<details>
<summary>Hint 1 — Setting a response header</summary>

Use `w.Header().Set(key, value)` before calling `w.WriteHeader(...)` or `w.Write(...)`.

```go
w.Header().Set("Content-Type", "application/json")
```
</details>

<details>
<summary>Hint 2 — Writing the status and body</summary>

After setting headers, write the status code, then the body:

```go
w.WriteHeader(http.StatusOK)
fmt.Fprintf(w, `{"message":"hello"}`)
```

If you call `w.Write(...)` without calling `w.WriteHeader(...)` first, Go automatically sends status 200. But being explicit is good practice.
</details>

<details>
<summary>Hint 3 — Using json.NewEncoder</summary>

For encoding structs or maps:

```go
type Response struct {
    Message string `json:"message"`
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(Response{Message: "hello"})
```

`json.NewEncoder(w).Encode(v)` writes the JSON directly to the `ResponseWriter` and appends a newline.
</details>

## How to Verify

Run the tests from inside the starter directory:

```bash
cd starter/task-01-hello-handler
go test ./...
```

The test sends a `GET` request to your handler using `httptest.NewRecorder()` and checks:

- Status code is `200`
- `Content-Type` header contains `application/json`
- Body is exactly `{"message":"hello"}`

All three assertions must pass.
