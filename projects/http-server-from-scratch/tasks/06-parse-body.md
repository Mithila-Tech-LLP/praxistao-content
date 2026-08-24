---
title: Parse Request Body
number: 6
difficulty: medium
duration: 20-25 minutes
concept: Content-Length, io.ReadFull
---

## What to Build

Implement `ParseBody` which reads exactly `Content-Length` bytes from a buffered reader, or returns an empty slice when the header is absent.

## Function Signature

```go
func ParseBody(r *bufio.Reader, headers map[string]string) ([]byte, error)
```

## Requirements

- If `"content-length"` header is absent, return `[]byte{}`, nil
- If `"content-length"` is `"0"`, return `[]byte{}`, nil
- Parse the value with `strconv.Atoi`; return an error if it is not a valid integer
- Read exactly that many bytes using `io.ReadFull`
- Propagate any read error from `io.ReadFull`

## Key Concept: Content-Length

The `Content-Length` header tells the receiver exactly how many bytes to read. Without it, a receiver cannot determine where the body ends.

```
POST /submit HTTP/1.1\r\n
Content-Length: 16\r\n
\r\n
{"name":"Alice"}
```

## Hints

<details>
<summary>Hint 1: io.ReadFull</summary>

`io.ReadFull(r, buf)` reads exactly `len(buf)` bytes or returns an error. Allocate `buf := make([]byte, n)` first.

</details>

<details>
<summary>Hint 2: Header key case</summary>

`ParseHeaders` lowercases all keys, so look up `"content-length"` not `"Content-Length"`.

</details>

## How to Verify

```bash
lncli run
```
