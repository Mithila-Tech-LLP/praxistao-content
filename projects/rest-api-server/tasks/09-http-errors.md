---
title: HTTP Errors
task: 09
slug: http-errors
concept: Error responses, http.Error, custom error types
difficulty: intermediate
---

## What You Will Build

Write helper functions for consistent JSON responses and a handler that carefully validates input and returns meaningful error codes. Good error handling separates a professional API from a toy.

## Function Signatures

```go
func WriteJSON(w http.ResponseWriter, status int, v any)

func WriteError(w http.ResponseWriter, status int, message string)

func SafeDivideHandler(w http.ResponseWriter, r *http.Request)
```

## WriteJSON

Sets `Content-Type: application/json`, writes the given status code, then encodes `v` as JSON.

## WriteError

Calls `WriteJSON` with `map[string]string{"error": message}`.

## SafeDivideHandler

Reads query parameters `?a=` and `?b=` and divides `a / b`.

| Condition                     | Status | Body                              |
|-------------------------------|--------|-----------------------------------|
| `a` or `b` missing            | 400    | `{"error":"missing parameter"}`   |
| `a` or `b` not a valid number | 400    | `{"error":"invalid number"}`      |
| `b == 0`                      | 422    | `{"error":"division by zero"}`    |
| success                       | 200    | `{"result":N}`                    |

Check in this order: missing → invalid → zero.

## Examples

```
GET /divide?a=10&b=2   → 200 {"result":5}
GET /divide?a=10&b=0   → 422 {"error":"division by zero"}
GET /divide?a=foo&b=2  → 400 {"error":"invalid number"}
GET /divide?a=10       → 400 {"error":"missing parameter"}
```

Note: `result` is an integer (integer division: `10/3 = 3`).

## Key Concepts

**Consistent error format** — every error response should have the same shape (`{"error":"..."}`) so clients can handle them uniformly. This is why `WriteError` is worth writing once and reusing.

**strconv.Atoi** — convert a query string value to an integer:

```go
a, err := strconv.Atoi(r.URL.Query().Get("a"))
```

**422 Unprocessable Entity** — used when the request is syntactically valid but semantically wrong. Division by zero is a valid request (the numbers parsed correctly) but the operation cannot be performed — 422 fits better than 400.

**Order of validation** — always check presence first, then validity, then business rules. This gives clearer error messages.

## Hints

<details>
<summary>Hint 1 — Implementing WriteJSON and WriteError</summary>

```go
func WriteJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, message string) {
    WriteJSON(w, status, map[string]string{"error": message})
}
```
</details>

<details>
<summary>Hint 2 — Checking for missing parameters</summary>

```go
aStr := r.URL.Query().Get("a")
bStr := r.URL.Query().Get("b")
if aStr == "" || bStr == "" {
    WriteError(w, http.StatusBadRequest, "missing parameter")
    return
}
```

`Query().Get(key)` returns `""` both when the key is absent and when its value is an empty string. Both should be treated as missing.
</details>

<details>
<summary>Hint 3 — Returning the result</summary>

```go
WriteJSON(w, http.StatusOK, map[string]int{"result": a / b})
```

This encodes to `{"result":5}` — note no quotes around the number.
</details>

## How to Verify

```bash
cd starter/task-09-http-errors
go test ./...
```

The test covers all four branches:

- `?a=10&b=2` → 200 + `{"result":5}`
- `?a=10` (missing b) → 400 + `{"error":"missing parameter"}`
- `?a=10&b=0` → 422 + `{"error":"division by zero"}`
- `?a=foo&b=2` → 400 + `{"error":"invalid number"}`
