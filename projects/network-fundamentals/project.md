---
title: Network Fundamentals
subtitle: Implement 10 core networking building blocks from scratch in Go — subnetting, checksums, routing, sockets, and DNS
category: Systems Programming
difficulty: intermediate
duration: 6-8 hours
accent: "#38bdf8"
technologies: [Go]
skills: [Subnetting, Checksums, Routing, TCP, UDP, DNS, Flow Control, Rate Limiting]
prerequisites: [basic-programming]
repo: network-fundamentals
outcomes:
  - Calculate subnet ranges and host counts from CIDR notation
  - Implement the Internet checksum algorithm used by IP, TCP, and UDP
  - Build a longest-prefix-match router lookup
  - Simulate a bounded ARP cache with LRU eviction
  - Write real TCP and UDP echo servers
  - Encode and decode DNS names in wire format
  - Simulate TCP's sliding window flow control
  - Implement a token-bucket rate limiter
  - Parse a raw IPv4 header byte-for-byte
---

## Overview

Every protocol you've ever used — HTTP, DNS, TLS, gRPC — is built out of a small set of primitives that show up again and again: address math, checksums, routing decisions, sockets, flow control, and rate limiting. Learn those primitives once, correctly, and every higher-level protocol you read about afterward stops looking like magic and starts looking like a specific arrangement of pieces you already understand.

This project isolates ten of those primitives as small, independently gradable Go programs. There's no framework here and nothing to glue together — each task is a self-contained function or type you implement and test on its own, the same way you'd verify a single idea in isolation before trusting it inside something larger.

Unlike `http-server-from-scratch`, there is no single running application that ties these tasks together. You won't build one router that also does checksums that also serves DNS. Instead, think of this as ten short, focused exercises — each one a hands-on companion to a specific set of chapters in the Computer Networks course, letting you turn what you read into code that actually runs and gets checked against real test vectors.

## What You Will Build

1. **Subnet Calculator** — compute network/broadcast addresses, host ranges, and host counts from CIDR notation
2. **Internet Checksum** — the RFC 1071 one's-complement checksum used by IP, TCP, and UDP headers
3. **Longest Prefix Match Router** — a routing table lookup that always prefers the most specific matching route
4. **Bounded ARP Cache** — an IP-to-MAC cache with LRU eviction
5. **TCP Echo Server** — a real `net.Listen`-based server with a persistent per-connection echo loop
6. **UDP Echo Server** — a connectionless echo server built on `net.ListenUDP`
7. **DNS Name Encoding** — encode and decode domain names in DNS's length-prefixed wire format
8. **Sliding Window Simulation** — TCP-style flow control with cumulative acknowledgments
9. **Token Bucket Rate Limiter** — a deterministic, time-injected rate limiter
10. **IPv4 Header Parser** — byte-for-byte parsing of a raw IPv4 header
