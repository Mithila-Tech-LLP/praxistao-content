---
title: Auth Interceptor
number: 8
difficulty: medium
duration: 20-25 minutes
concept: Authentication via Interceptor, Bearer Tokens
---

## What You Need to Build

Write an authentication interceptor that validates a Bearer token in metadata.

## What to Implement

```go
func AuthInterceptor(token string) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
        // TODO: read "authorization" from metadata
        // Expected value: "Bearer <token>"
        // Missing or wrong → codes.Unauthenticated
        // Correct → call handler
    }
}
```

## Requirements

- Missing `authorization` metadata → `codes.Unauthenticated`
- Value not `"Bearer <expected_token>"` → `codes.Unauthenticated`
- Correct token → proceed normally

## How to Verify

```bash
lncli run
```
