---
title: TCP Echo Server
number: 5
difficulty: medium
duration: 25-35 minutes
concept: net.Listen, Accept loop, io.Copy
---

## What to Build

Implement `EchoServer.ListenAndServe`, a real TCP server that echoes back every byte it reads on a connection, for as long as the connection stays open.

## Function Signature

```go
type EchoServer struct{ Addr string }

func (s *EchoServer) ListenAndServe() (net.Listener, error)
```

## Requirements

- Open a TCP listener on `s.Addr`
- Return the listener immediately, so the caller (or a test) can call `ln.Close()` to stop the server
- In a goroutine, run an accept loop: call `Accept()` repeatedly, and handle each connection in its own goroutine
- For each connection, copy every byte read back out to the same connection (`io.Copy(conn, conn)` or an explicit read/write loop) until the client closes the connection or an error occurs, then close the connection

## Key Concept: The Accept Loop

A TCP server is fundamentally a loop: `Listen` creates a socket bound to an address, then `Accept` blocks until a client connects, handing back a new `net.Conn` for that specific connection. Because `Accept` only returns one connection at a time, a real server spins it off into its own goroutine immediately and loops back to `Accept` again — otherwise it could only ever talk to one client at once. This pattern is introduced in Chapter 106 (building a TCP server and client) and underpins the three-way handshake mechanics from Chapter 59 — every connection you accept here completed a real handshake before your code ever saw it.

## Hints

<details>
<summary>Hint 1: Don't block the caller</summary>

`ListenAndServe` must return the `net.Listener` right away — the accept loop needs to run in its own goroutine (`go func() { ... }()`), not inline, or the function would never return until the listener closes.

</details>

<details>
<summary>Hint 2: Handling Accept's error on shutdown</summary>

When the test calls `ln.Close()`, the blocked `Accept()` call will return an error. Treat any error from `Accept` as a signal to stop the loop and return from the goroutine — don't retry forever.

</details>

<details>
<summary>Hint 3: io.Copy as the echo loop</summary>

`io.Copy(conn, conn)` reads from `conn` and writes whatever it reads back to `conn`, looping until it hits EOF or an error — which is exactly "echo everything back until the client disconnects." Close the connection afterward with `defer conn.Close()` inside the per-connection goroutine.

</details>

## How to Verify

```bash
lncli run
```
