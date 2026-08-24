---
title: Path Parameters
task: 03
slug: path-parameters
concept: URL Parsing, strings.TrimPrefix
difficulty: beginner
---

## What You Will Build

Write a handler that reads a dynamic segment from the URL path — the "path parameter". Real APIs use this all the time: `GET /users/42` means "get user with ID 42".

## Function Signature

```go
func UserHandler(w http.ResponseWriter, r *http.Request)
```

This handler is registered at `/users/`, so `r.URL.Path` will be something like `/users/42`.

## What It Should Do

1. Extract the ID from the path (everything after `/users/`)
2. Return a JSON response:
   ```json
   {"id":"42","name":"User 42"}
   ```
   where both `id` and `name` use the extracted ID string.

## Examples

```
GET /users/42
→ 200 OK
→ {"id":"42","name":"User 42"}

GET /users/abc
→ 200 OK
→ {"id":"abc","name":"User abc"}

GET /users/hello-world
→ 200 OK
→ {"id":"hello-world","name":"User hello-world"}
```

## Key Concepts

**strings.TrimPrefix** — the simplest way to strip a known prefix from a string:

```go
id := strings.TrimPrefix(r.URL.Path, "/users/")
// r.URL.Path = "/users/42"  →  id = "42"
```

**Why no built-in path params?** — Go's `net/http` doesn't parse `{id}` style parameters automatically. You extract them manually from `r.URL.Path`. Popular frameworks like `chi` or `gin` add this, but it's easy to do yourself for simple cases.

**fmt.Sprintf for building strings** — use `fmt.Sprintf("User %s", id)` to build the name field.

## Hints

<details>
<summary>Hint 1 — Extracting the ID</summary>

```go
id := strings.TrimPrefix(r.URL.Path, "/users/")
```

This removes exactly the prefix `/users/` and leaves the rest. If the path is `/users/42`, `id` becomes `"42"`.
</details>

<details>
<summary>Hint 2 — Building the JSON response</summary>

Use a struct with JSON tags or `fmt.Sprintf`:

```go
type UserResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

resp := UserResponse{
    ID:   id,
    Name: fmt.Sprintf("User %s", id),
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(resp)
```
</details>

<details>
<summary>Hint 3 — Registering the handler</summary>

In `net/http`, registering `/users/` (with trailing slash) matches all paths starting with `/users/`:

```go
mux.HandleFunc("/users/", UserHandler)
```

Your handler then extracts whatever comes after `/users/`.
</details>

## How to Verify

```bash
cd starter/task-03-path-parameters
go test ./...
```

The test sends requests to `/users/42` and `/users/abc` and checks:

- Status `200` for both
- Body for `/users/42` is `{"id":"42","name":"User 42"}`
- Body for `/users/abc` is `{"id":"abc","name":"User abc"}`
