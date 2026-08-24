---
title: Metadata
number: 7
difficulty: medium
duration: 20-25 minutes
concept: gRPC Metadata, Request Headers
---

## What You Need to Build

Implement `WhoAmI` — reads a user ID from request metadata.

## Service Addition

```proto
rpc WhoAmI(google.protobuf.Empty) returns (WhoAmIReply);
message WhoAmIReply { string user_id = 1; }
```

## What to Implement

```go
func (s *server) WhoAmI(ctx context.Context, _ *emptypb.Empty) (*WhoAmIReply, error) {
    // TODO: read "x-user-id" from incoming metadata
    // If missing: return codes.Unauthenticated error
    // If present: return &WhoAmIReply{UserId: value}
}
```

## Key Concept: Reading Metadata

```go
import "google.golang.org/grpc/metadata"

md, ok := metadata.FromIncomingContext(ctx)
vals := md.Get("x-user-id")
```

## How to Verify

```bash
lncli run
```
