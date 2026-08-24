---
title: Persistent Connections
number: 7
difficulty: medium
duration: 25-30 minutes
concept: Connection: keep-alive, HTTP/1.1
---

## What to Build

Implement `ServeConn` which handles multiple HTTP requests on a single connection in a loop, dispatching each to the provided function and writing the response back.

## Function Signature

```go
func ServeConn(conn net.Conn, dispatch func(method, path string, headers map[string]string) Response)
```

## Requirements

- Wrap `conn` in a `bufio.Reader` once; use it for all reads in the loop
- Read the request line, then headers, on each iteration
- Dispatch the request and write the response to `conn`
- Break the loop when the `"connection"` header equals `"close"`
- Break the loop when reading fails (EOF or any error)

## Key Concept: Keep-Alive

HTTP/1.1 reuses the same TCP connection for multiple requests by default. The connection only closes when the client sends `Connection: close` or the connection errors/times out.

## Hints

<details>
<summary>Hint 1: Reading the request line in a loop</summary>

Call `reader.ReadString('\n')` at the top of each loop iteration to get the request line. An error here (including EOF) means the client disconnected.

</details>

<details>
<summary>Hint 2: Detecting EOF</summary>

`errors.Is(err, io.EOF)` checks for a clean client disconnect. Any non-nil error from Read should break the loop.

</details>

## How to Verify

```bash
lncli run
```
