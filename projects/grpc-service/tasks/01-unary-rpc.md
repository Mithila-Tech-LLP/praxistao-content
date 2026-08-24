---
title: Unary RPC
number: 1
difficulty: easy
duration: 15-20 minutes
concept: RPC, Request/Response, gRPC Server
---

## What You Need to Build

Implement a `GreeterServer` that handles a simple `SayHello` unary RPC call.

## Service Definition

```proto
service Greeter {
  rpc SayHello(HelloRequest) returns (HelloReply);
}
message HelloRequest { string name = 1; }
message HelloReply   { string message = 1; }
```

The generated types and service interface are already in the starter — you just implement the method.

## What to Implement

```go
type server struct {
    UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
    // TODO: return &HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}
```

## Requirements

- Empty name → return `&HelloReply{Message: "Hello, !"}`
- Any name → return `"Hello, <name>!"`

## How to Verify

```bash
lncli run
```
