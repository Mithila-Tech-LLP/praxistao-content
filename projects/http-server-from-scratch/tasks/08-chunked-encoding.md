---
title: Chunked Transfer Encoding
number: 8
difficulty: hard
duration: 25-30 minutes
concept: Transfer-Encoding: chunked
---

## What to Build

Implement `WriteChunked` and `ReadChunked` for HTTP chunked transfer encoding — the mechanism HTTP uses to send a response body when the total length is not known in advance.

## Function Signature

```go
func WriteChunked(w io.Writer, chunks [][]byte) error
func ReadChunked(r *bufio.Reader) ([]byte, error)
```

## Requirements

- `WriteChunked`: for each chunk write `"<hex-length>\r\n<data>\r\n"`
- `WriteChunked`: after all chunks write the terminal `"0\r\n\r\n"`
- `WriteChunked`: empty `chunks` slice writes only the terminal chunk
- `ReadChunked`: read the hex size line, then read that many bytes plus `"\r\n"`
- `ReadChunked`: stop and return all accumulated data when chunk size is `0`

## Key Concept: Chunked Wire Format

```
5\r\n
Hello\r\n
6\r\n
 World\r\n
0\r\n
\r\n
```

Each chunk: hex size + CRLF + data + CRLF. Terminal chunk: `0` + CRLF + CRLF.

## Hints

<details>
<summary>Hint 1: Writing the hex size</summary>

Use `fmt.Fprintf` with `"%x\r\n"` to write the chunk size in lowercase hex followed by CRLF.

</details>

<details>
<summary>Hint 2: Reading the chunk size</summary>

Read the size line with `r.ReadString('\n')`, trim whitespace, then `strconv.ParseInt(sizeStr, 16, 64)` parses the hex value.

</details>

## How to Verify

```bash
lncli run
```
