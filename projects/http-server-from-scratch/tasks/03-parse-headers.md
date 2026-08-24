---
title: Parse Headers
number: 3
difficulty: easy
duration: 20-25 minutes
concept: HTTP Headers, CRLF, bufio.Reader
---

## What to Build

Implement `ParseHeaders` which reads HTTP headers line by line from a `bufio.Reader` until a blank line, returning a map with lowercased keys.

## Function Signature

```go
func ParseHeaders(r *bufio.Reader) (map[string]string, error)
```

## Requirements

- Read lines until an empty line (just `"\r\n"` or `"\n"`)
- Split each non-blank line on the first `":"` only
- Lowercase all header keys
- Trim leading/trailing whitespace from both key and value
- Return `ErrMalformedHeader` for any line that does not contain `":"`

## Key Concept: HTTP Header Format

```
Content-Type: application/json\r\n
Host: example.com\r\n
\r\n
```

Each header is `Key: Value\r\n`. A blank line marks the end of the header section.

## Hints

<details>
<summary>Hint 1: bufio.Reader.ReadString</summary>

Use `r.ReadString('\n')` to read one line at a time. Trim `"\r\n"` with `strings.TrimRight(line, "\r\n")` to get the raw content.

</details>

<details>
<summary>Hint 2: Splitting on the first colon only</summary>

`strings.SplitN(line, ":", 2)` splits into at most 2 parts, so colons in values (e.g. `Date: Mon, 01 Jan`) are preserved correctly.

</details>

## How to Verify

```bash
lncli run
```
