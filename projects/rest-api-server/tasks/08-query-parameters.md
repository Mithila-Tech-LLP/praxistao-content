---
title: Query Parameters
task: 08
slug: query-parameters
concept: r.URL.Query(), URL parsing, filtering
difficulty: beginner
---

## What You Will Build

Write a handler that filters a list of items based on a URL query parameter. Query parameters are the `?key=value` parts of a URL — used for filtering, sorting, pagination, and search in real APIs.

## Types and Package-Level State

```go
type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var items = []Item{
    {ID: 1, Name: "apple"},
    {ID: 2, Name: "banana"},
    {ID: 3, Name: "apricot"},
}
```

## Function Signature

```go
func SearchHandler(w http.ResponseWriter, r *http.Request)
```

## What It Should Do

- Read the `q` query parameter from the URL
- If `q` is present and non-empty: return only items whose `Name` contains `q` as a **case-insensitive** substring
- If `q` is absent or empty: return all items
- Always return a JSON array (even if empty: `[]`)

## Examples

```
GET /search?q=ap
→ [{"id":1,"name":"apple"},{"id":3,"name":"apricot"}]

GET /search?q=ban
→ [{"id":2,"name":"banana"}]

GET /search?q=xyz
→ []

GET /search
→ [{"id":1,"name":"apple"},{"id":2,"name":"banana"},{"id":3,"name":"apricot"}]
```

## Key Concepts

**r.URL.Query()** — returns a `url.Values` map of all query parameters:

```go
params := r.URL.Query()
q := params.Get("q")   // "" if not present
```

**strings.Contains and strings.ToLower** — for case-insensitive matching:

```go
strings.Contains(strings.ToLower(item.Name), strings.ToLower(q))
```

**Always return an array** — even when there are zero results. Returning `null` instead of `[]` breaks many JSON clients. Initialise the result slice as `[]Item{}` (not `var result []Item`) so it encodes to `[]` instead of `null`.

## Hints

<details>
<summary>Hint 1 — Reading the query parameter</summary>

```go
q := r.URL.Query().Get("q")
// q is "" if the parameter is not in the URL
```
</details>

<details>
<summary>Hint 2 — Filtering with case-insensitive match</summary>

```go
result := []Item{}
for _, item := range items {
    if q == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(q)) {
        result = append(result, item)
    }
}
```
</details>

<details>
<summary>Hint 3 — Encoding an empty array vs null</summary>

In Go, a nil slice encodes to `null` in JSON. To always get `[]`:

```go
result := []Item{}    // non-nil empty slice → encodes to []
// NOT: var result []Item  → nil slice → encodes to null
```
</details>

## How to Verify

```bash
cd starter/task-08-query-parameters
go test ./...
```

The test covers four cases:

- `?q=ap` — matches apple and apricot (2 results, in order)
- `?q=ban` — matches banana (1 result)
- `?q=xyz` — no matches (`[]`)
- no `q` parameter — all 3 items returned
