---
title: HTTP Routing
number: 5
difficulty: easy
duration: 20-25 minutes
concept: Handler Registration, Method Dispatch
---

## What to Build

Implement a simple `Mux` (multiplexer) that registers handlers by `"METHOD /path"` and dispatches incoming requests, returning a 404 `Response` for unknown routes.

## Function Signature

```go
type HandlerFunc func(method, path string, headers map[string]string) Response

func (m *Mux) Handle(method, path string, fn HandlerFunc)
func (m *Mux) Dispatch(method, path string, headers map[string]string) Response
```

## Requirements

- `Handle` stores the handler keyed by `"METHOD /path"` (e.g. `"GET /users"`)
- `Dispatch` looks up the key and calls the matching handler
- Unknown routes return `Response{Status: 404}`
- Method matching is exact and case-sensitive
- The same path with different methods is allowed (e.g. `GET /users` and `POST /users`)

## Key Concept: Route Tables

A route table maps a (method, path) pair to a handler function. The simplest implementation uses a `map[string]HandlerFunc` where the key is `"METHOD /path"`.

## Hints

<details>
<summary>Hint 1: Building the map key</summary>

`key := method + " " + path` gives a unique, readable key like `"GET /users"`.

</details>

<details>
<summary>Hint 2: Returning 404</summary>

`return Response{Status: 404}` satisfies the requirement for unknown routes. No body required.

</details>

## How to Verify

```bash
lncli run
```
