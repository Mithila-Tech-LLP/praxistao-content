---
title: Parse Request Line
number: 2
difficulty: easy
duration: 15-20 minutes
concept: HTTP Request Format, strings.Fields
---

## What to Build

Implement `ParseRequestLine` which splits an HTTP request line like `"GET /path HTTP/1.1"` into its three components.

## Function Signature

```go
func ParseRequestLine(line string) (method, path, version string, err error)
```

## Requirements

- Split on whitespace; return exactly method, path, and version
- Return `ErrMalformedRequestLine` if the line does not have exactly 3 fields
- Return `ErrMalformedRequestLine` for an empty string
- Preserve the path exactly, including query strings (e.g. `/path?q=1`)
- Method and version are returned as-is (no normalization required)

## Key Concept: HTTP Request Line

Every HTTP request starts with a request line:

```
GET /path?query=1 HTTP/1.1
```

Three fields separated by a single space. `strings.Fields` handles multiple spaces and is the idiomatic way to split.

## Hints

<details>
<summary>Hint 1: strings.Fields</summary>

`strings.Fields(s)` splits on any whitespace and returns a slice. Check `len(fields) != 3` to detect malformed lines.

</details>

<details>
<summary>Hint 2: Returning named values</summary>

With named return values you can return the zero values plus `ErrMalformedRequestLine` explicitly, or assign to the named results and use a bare `return`.

</details>

## How to Verify

```bash
lncli run
```
