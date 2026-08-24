---
title: gRPC Status Codes
number: 2
difficulty: easy
duration: 15-20 minutes
concept: gRPC Status Codes, google.golang.org/grpc/codes
---

## What You Need to Build

Extend `SayHello` to return proper gRPC error codes for invalid inputs.

## Requirements

| Condition | Status Code | Message |
|-----------|-------------|---------|
| `name == ""` | `codes.InvalidArgument` | `"name is required"` |
| `name == "banned"` | `codes.PermissionDenied` | `"user is banned"` |
| otherwise | success | `"Hello, <name>!"` |

## Key Concept: gRPC Status Codes

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

return nil, status.Errorf(codes.InvalidArgument, "name is required")
```

## How to Verify

```bash
lncli run
```
