---
title: Build a gRPC Service
subtitle: Learn RPC, streaming, interceptors, and metadata with real gRPC in Go
category: Web APIs
difficulty: intermediate
duration: 4-6 hours
accent: "#818cf8"
technologies: [Go, gRPC, Protocol Buffers]
skills: [Unary RPC, Streaming, Interceptors, Metadata, Error Handling]
prerequisites: [rest-api-server]
repo: grpc-service
outcomes:
  - Implement unary, server-streaming, client-streaming, and bidirectional RPC
  - Return proper gRPC status codes for different error conditions
  - Write server-side interceptors for logging and authentication
  - Pass and read metadata in gRPC requests
  - Test gRPC services end-to-end with a real in-process server
---

## Overview

gRPC is Google's open-source RPC framework used at Netflix, Cloudflare, and Dropbox. It uses Protocol Buffers for serialization and HTTP/2 for transport, giving you type-safe APIs, bidirectional streaming, and generated client/server code.

In this project the proto definitions and generated code are pre-included in each starter directory — so you focus on implementing the service logic.

**First run note:** The first `lncli run` downloads gRPC dependencies (~30s). Subsequent runs are instant.

## What You Will Build

A Greeter service with progressive complexity: unary calls, all four streaming modes, error handling, interceptors, and metadata-based authentication.
