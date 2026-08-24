---
title: Client Streaming
number: 4
difficulty: medium
duration: 20-25 minutes
concept: Client-Side Streaming, Recv Loop
---

## What You Need to Build

Implement `CollectNames` — a client-streaming RPC that collects all names then replies once.

## Service Addition

```proto
rpc CollectNames(stream HelloRequest) returns (HelloReply);
```

## What to Implement

```go
func (s *server) CollectNames(stream Greeter_CollectNamesServer) error {
    // TODO: recv all requests, collect names, send one reply
}
```

## Requirements

- Receive all `HelloRequest` messages (until `io.EOF`)
- Join names with `", "`: `"Hello, Alice, Bob, Charlie!"`
- Zero names → `"Hello, !"`

## How to Verify

```bash
lncli run
```
