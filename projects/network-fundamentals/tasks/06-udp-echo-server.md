---
title: UDP Echo Server
number: 6
difficulty: easy-medium
duration: 20-30 minutes
concept: net.ListenUDP, ReadFromUDP/WriteToUDP, connectionless
---

## What to Build

Implement `StartUDPEchoServer`, a UDP server that reads datagrams and sends each one straight back to whoever sent it — no connection setup involved, because UDP doesn't have connections.

## Function Signature

```go
func StartUDPEchoServer(addr string) (*net.UDPConn, error)
```

## Requirements

- Open a UDP socket via `net.ListenUDP` bound to `addr`
- Return the `*net.UDPConn` immediately — the caller stops the server by calling `conn.Close()`
- In a goroutine, loop forever: `ReadFromUDP` into a buffer, then `WriteToUDP` the same bytes back to the sender's address
- Must keep looping across multiple datagrams, not stop after the first one

## Key Concept: UDP is Connectionless

Unlike TCP, a UDP socket doesn't `Accept` individual connections — a single socket can receive datagrams from any number of different senders, and each read tells you exactly who sent it via the returned address. That's the crucial difference this task makes concrete: there's no handshake, no per-client state, just "receive a packet, learn its sender's address, reply directly to that address." Chapter 58 covers UDP itself, and Chapter 107 (building a UDP server and client) is the direct companion to this task.

## Hints

<details>
<summary>Hint 1: ReadFromUDP tells you who sent it</summary>

`ReadFromUDP` gives you both the bytes received and the sender's `*net.UDPAddr` — you need that address to reply, since there's no persistent connection to write back to.

</details>

<details>
<summary>Hint 2: Let the OS pick a port in tests</summary>

Binding to `"127.0.0.1:0"` lets the OS assign a free port. After `ListenUDP` succeeds, read the actual chosen address back out with `conn.LocalAddr()` so your test (or a client) knows where to send datagrams.

</details>

<details>
<summary>Hint 3: The loop must not stop after one packet</summary>

A common mistake is to `ReadFromUDP` once and return. Wrap the read/write pair in a `for { }` loop inside the goroutine so the server keeps echoing every subsequent datagram it receives.

</details>

## How to Verify

```bash
lncli run
```
