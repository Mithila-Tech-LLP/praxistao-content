---
title: Bidirectional Streaming
number: 5
difficulty: medium
duration: 25-30 minutes
concept: Bidirectional Streaming, Concurrent Send/Recv
---

## What You Need to Build

Implement `Chat` — a bidirectional streaming RPC that echoes each message back.

## Service Addition

```proto
rpc Chat(stream HelloRequest) returns (stream HelloReply);
```

## What to Implement

```go
func (s *server) Chat(stream Greeter_ChatServer) error {
    // TODO: for each received request, immediately send back "Echo: <name>"
}
```

## Requirements

- For each received `HelloRequest{Name: "X"}`, send `HelloReply{Message: "Echo: X"}`
- Continue until client closes the stream (recv returns io.EOF)
- Return nil on EOF

## How to Verify

```bash
lncli run
```
