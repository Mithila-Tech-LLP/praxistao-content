---
title: In-Memory CRUD
task: 06
slug: in-memory-crud
concept: map, CRUD Operations, REST Conventions
difficulty: intermediate
---

## What You Will Build

Implement a complete CRUD (Create, Read, Update, Delete) API for Todo items backed by an in-memory map. This is the most common pattern in REST API development.

## Types and Package-Level State

```go
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var db     = map[int]Todo{}
var nextID = 1
```

## Function Signature

```go
func TodosHandler(w http.ResponseWriter, r *http.Request)
```

This single handler is registered at `/todos/` and must dispatch on both method and path.

## REST Route Table

| Method | Path         | Action              | Success Status | Error       |
|--------|--------------|---------------------|----------------|-------------|
| GET    | /todos       | List all todos      | 200 + JSON array | —         |
| POST   | /todos       | Create todo         | 201 + new Todo | —           |
| GET    | /todos/{id}  | Get one todo        | 200 + Todo     | 404 if missing |
| PUT    | /todos/{id}  | Replace todo        | 200 + updated Todo | 404 if missing |
| DELETE | /todos/{id}  | Delete todo         | 204            | 404 if missing |

## Key Concepts

**Dispatch logic** — parse `r.URL.Path` to determine if the request targets the collection (`/todos` or `/todos/`) or a single item (`/todos/42`). Then switch on `r.Method`.

**strconv.Atoi** — convert the string ID from the URL to an integer:

```go
id, err := strconv.Atoi(idStr)
if err != nil {
    http.Error(w, "invalid id", http.StatusBadRequest)
    return
}
```

**Ordered list from a map** — when listing all todos, convert the map to a slice. Maps iterate in random order; for predictable tests, sort by ID:

```go
var list []Todo
for _, t := range db {
    list = append(list, t)
}
sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
```

**404 for missing items** — `_, ok := db[id]` — if `!ok`, write `404` and return.

## Hints

<details>
<summary>Hint 1 — Detecting collection vs single-item path</summary>

```go
path := strings.TrimPrefix(r.URL.Path, "/todos")
path = strings.TrimPrefix(path, "/")
// path is "" for /todos, and "42" for /todos/42
```

Check `path == ""` for collection routes, otherwise parse the integer ID.
</details>

<details>
<summary>Hint 2 — Handling PUT (update)</summary>

```go
var input Todo
json.NewDecoder(r.Body).Decode(&input)
todo := db[id]
if input.Title != "" {
    todo.Title = input.Title
}
todo.Done = input.Done
db[id] = todo
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(todo)
```
</details>

<details>
<summary>Hint 3 — Incrementing nextID safely</summary>

```go
// POST /todos
var input struct{ Title string `json:"title"` }
json.NewDecoder(r.Body).Decode(&input)
todo := Todo{ID: nextID, Title: input.Title}
db[nextID] = todo
nextID++
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(todo)
```
</details>

## How to Verify

```bash
cd starter/task-06-in-memory-crud
go test ./...
```

The test resets `db` and `nextID`, then exercises all five operations in order:

1. POST two todos
2. GET list (expects both)
3. GET single (expects the first)
4. PUT (update title)
5. DELETE (then GET returns 404)
