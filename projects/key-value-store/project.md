---
title: Build a Key-Value Store
subtitle: Build a Redis-inspired in-memory store with strings, lists, sets, and pub/sub
category: Systems Programming
difficulty: intermediate
duration: 5-8 hours
accent: "#f97316"
technologies: [Go]
skills: [Maps, Concurrency, sync.RWMutex, Pub/Sub, Transactions, TTL]
prerequisites: [basic-programming]
repo: key-value-store
outcomes:
  - Implement a thread-safe in-memory key-value store
  - Support TTL-based key expiration
  - Build list and set data type operations
  - Implement a publish/subscribe message bus
  - Write atomic multi-command transactions
  - Implement LRU eviction when capacity is exceeded
---

## Overview

Redis is one of the most widely used pieces of infrastructure in the world — powering caches, queues, pub/sub systems, rate limiters, and session stores at companies like Twitter, GitHub, and Airbnb. In this project you will build a simplified version of it from scratch.

Each task adds a new capability. By the end you will have a concurrent, TTL-aware, feature-rich in-memory store that you can use as the backbone for real applications.
