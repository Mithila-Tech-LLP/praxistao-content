---
title: TCP Listener
number: 1
difficulty: easy
duration: 15-20 minutes
concept: net.Listen, Accept, Conn
---

## What to Build

Implement `ListenAndServe` on the `Server` struct. It should open a TCP listener on `s.Addr`, accept one connection, read up to 1024 bytes (draining the request), write a fixed HTTP 200 response, and close the connection.

## Function Signature

```go
type Server struct{ Addr string }

func (s *Server) ListenAndServe() error
```

## Requirements

- Call `net.Listen("tcp", s.Addr)` to open the listener
- Call `listener.Accept()` to get a `net.Conn` (blocks until a client connects)
- Read up to 1024 bytes from the connection (ignoring partial reads for now)
- Write exactly `"HTTP/1.1 200 OK\r\n\r\nOK\n"` to the connection
- Close the connection and return nil

## Key Concept: TCP in Go

```go
ln, err := net.Listen("tcp", ":8080")
conn, err := ln.Accept() // blocks until a client connects
buf := make([]byte, 1024)
conn.Read(buf)
conn.Write([]byte("response"))
conn.Close()
```

## Hints

<details>
<summary>Hint 1: Getting the actual address</summary>

Use `ln.Addr().String()` after `net.Listen` to get the address including the assigned port (useful when using `:0` for a random port in tests).

</details>

<details>
<summary>Hint 2: Ignoring read errors for now</summary>

Call `conn.Read(buf)` and ignore the error — you just need to drain enough of the request before responding. Full request parsing comes in later tasks.

</details>

## How to Verify

```bash
lncli run
```
