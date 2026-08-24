---
title: JSON Request Body
task: 05
slug: json-request-body
concept: json.NewDecoder, json.NewEncoder
difficulty: intermediate
---

## What You Will Build

Write a handler that reads a JSON body from the incoming request, creates a new Todo item, stores it, and returns the created item. This is the heart of any POST endpoint in a REST API.

## Types and Package-Level State

```go
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var store []Todo
```

The `store` slice is shared across requests — it acts as an in-memory database for this task.

## Function Signature

```go
func CreateTodoHandler(w http.ResponseWriter, r *http.Request)
```

## What It Should Do

1. Decode the JSON request body into a struct (only `title` is sent by the client)
2. Assign `ID = len(store) + 1` to the new todo
3. Set `Done = false`
4. Append the todo to `store`
5. Respond with status `201 Created` and the created Todo as JSON

## Example

```
POST /todos
Content-Type: application/json
Body: {"title":"buy milk"}

→ 201 Created
→ Content-Type: application/json
→ {"id":1,"title":"buy milk","done":false}

POST /todos
Body: {"title":"read a book"}

→ 201 Created
→ {"id":2,"title":"read a book","done":false}
```

## Key Concepts

**json.NewDecoder** — reads JSON from an `io.Reader` (like `r.Body`):

```go
var input struct{ Title string `json:"title"` }
json.NewDecoder(r.Body).Decode(&input)
```

**json.NewEncoder** — writes JSON to an `io.Writer` (like `w`):

```go
json.NewEncoder(w).Encode(todo)
```

**Why use Decoder/Encoder over Marshal/Unmarshal?** — `Decoder` reads directly from a stream without loading the full body into a `[]byte` first. It's more memory-efficient and idiomatic in HTTP handlers.

**Resetting package-level state in tests** — because `store` is package-level, the test resets it with `store = nil` before each sub-test. Your implementation should work with any initial state of `store`.

## Hints

<details>
<summary>Hint 1 — Decoding the request body</summary>

```go
var input struct {
    Title string `json:"title"`
}
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
    http.Error(w, "bad request", http.StatusBadRequest)
    return
}
```
</details>

<details>
<summary>Hint 2 — Assigning the ID</summary>

```go
todo := Todo{
    ID:    len(store) + 1,
    Title: input.Title,
    Done:  false,
}
store = append(store, todo)
```

`len(store)` before the append gives you 0-based count, so adding 1 gives IDs starting at 1.
</details>

<details>
<summary>Hint 3 — Writing the 201 response</summary>

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(todo)
```

Set the header before `WriteHeader`, and `WriteHeader` before writing the body.
</details>

## How to Verify

```bash
cd starter/task-05-json-request-body
go test ./...
```

The test resets `store`, then creates two todos and checks:

- First response: status `201`, `id` is `1`, `title` matches, `done` is `false`
- Second response: status `201`, `id` is `2`, `title` matches, `done` is `false`
