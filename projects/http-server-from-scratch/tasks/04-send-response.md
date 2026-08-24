---
title: Send HTTP Response
number: 4
difficulty: easy
duration: 15-20 minutes
concept: HTTP Response Format, Status Lines
---

## What to Build

Implement `Write` on the `Response` struct, which serializes a complete HTTP/1.1 response to any `io.Writer`.

## Function Signature

```go
type Response struct {
    Status  int
    Headers map[string]string
    Body    []byte
}

func (r *Response) Write(w io.Writer) error
```

## Requirements

- Write the status line: `"HTTP/1.1 <code> <text>\r\n"`
- Write each header as `"Key: Value\r\n"`
- Write a blank line `"\r\n"` after all headers
- Write the body bytes (if any)
- Return any write error immediately

## Key Concept: HTTP Response Format

```
HTTP/1.1 200 OK\r\n
Content-Type: application/json\r\n
Content-Length: 13\r\n
\r\n
{"ok": true}\n
```

Status line, then headers, then blank line, then body. All line endings are CRLF.

## Hints

<details>
<summary>Hint 1: fmt.Fprintf for the status line</summary>

`fmt.Fprintf` writes the status line in one call — format it as `"HTTP/1.1 %d %s\r\n"` with the status code and its reason phrase.

</details>

<details>
<summary>Hint 2: Writing the body</summary>

Use `w.Write(r.Body)` after the blank line. Calling `w.Write(nil)` writes zero bytes and is safe.

</details>

## How to Verify

```bash
lncli run
```
