---
title: Server Interceptor
number: 6
difficulty: medium
duration: 20-25 minutes
concept: Middleware for gRPC, UnaryServerInterceptor
---

## What You Need to Build

Write a logging interceptor that records every RPC call.

## What to Implement

```go
var CallLog []string

func LoggingInterceptor(
    ctx context.Context,
    req any,
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (any, error) {
    // TODO: call handler, then append "method:status_code" to CallLog
}
```

Register it: `grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))`

## Requirements

- Log format: `"<method>:OK"` on success, `"<method>:InvalidArgument"` on error
- Method is the short name from `info.FullMethod` (e.g. `"SayHello"`)

## How to Verify

```bash
lncli run
```
