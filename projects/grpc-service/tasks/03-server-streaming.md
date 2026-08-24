---
title: Server Streaming
number: 3
difficulty: medium
duration: 20-25 minutes
concept: Server-Side Streaming, Send Loop
---

## What You Need to Build

Implement `ListGreetings` — a server-streaming RPC that sends one reply per name.

## Service Addition

```proto
rpc ListGreetings(ListRequest) returns (stream HelloReply);
message ListRequest { repeated string names = 1; }
```

## What to Implement

```go
func (s *server) ListGreetings(req *ListRequest, stream Greeter_ListGreetingsServer) error {
    // TODO: for each name in req.Names, send stream.Send(&HelloReply{...})
}
```

## Requirements

- Send exactly one `HelloReply` per name in `req.Names`
- Message format: `"Hello, <name>!"`
- Empty names list → send nothing, return nil

## How to Verify

```bash
lncli run
```
