# Go Programming — Complete Backend Engineer Course
## From Zero to Production

**Course Tagline**: Start with "Hello, World." End with a production-grade distributed system.

---

## What You Will Build

By the end of this course, you will have built:

| Project | What | Tech |
|---------|------|------|
| Mini 1 | CLI Task Manager | Go fundamentals |
| Mini 2 | Concurrent Log Processor | Goroutines, channels |
| Mini 3 | In-Memory Key-Value Store | Data structures |
| Mini 4 | Route Planner | Algorithms, graphs |
| Mini 5 | Blog REST API | HTTP, chi, JWT |
| Mini 6 | Product Catalog with Search | PostgreSQL, OpenSearch |
| Mini 7 | Fully Observable API | Prometheus, Grafana, Sentry |
| Major 1 | Food Delivery Backend (Modular Monolith) | Clean arch, DDD, CQRS |
| Major 2 | Ticket Booking System (Event-Driven) | Kafka, Watermill, Sagas |
| Major 3 | E-Commerce Microservices Platform | Microservices, gRPC, service mesh |
| **Final** | **Production SaaS Backend** | **Everything: k8s, Datadog, CI/CD** |

---

## Volume 0: Welcome and Setup (Ch 00–02)

- **Ch 00** — Table of Contents and Learning Roadmap ← *you are here*
- [**Ch 01**](01-what-is-go-and-why-learn-it.md) — What Is Go and Why Learn It
- [**Ch 02**](02-setting-up-your-development-environment.md) — Setting Up Your Development Environment

---

## Volume 1: Go Fundamentals (Ch 03–17)

*Learn Go from absolute zero — variables, functions, collections, structs, interfaces.*

- [**Ch 03**](03-hello-world-packages-and-modules.md) — Hello World, Packages, and Modules
- [**Ch 04**](04-variables-constants-and-basic-types.md) — Variables, Constants, and Basic Types
- [**Ch 05**](05-operators-expressions-and-type-conversion.md) — Operators, Expressions, and Type Conversion
- [**Ch 06**](06-control-flow-if-for-switch.md) — Control Flow — if, for, switch
- [**Ch 07**](07-functions-the-building-blocks.md) — Functions — Parameters, Return Values, Variadic
- [**Ch 08**](08-arrays-and-slices.md) — Arrays and Slices
- [**Ch 09**](09-maps.md) — Maps
- [**Ch 10**](10-structs.md) — Structs
- [**Ch 11**](11-methods.md) — Methods
- [**Ch 12**](12-interfaces.md) — Interfaces
- [**Ch 13**](13-pointers.md) — Pointers
- [**Ch 14**](14-error-handling.md) — Error Handling
- [**Ch 15**](15-defer-panic-recover.md) — Defer, Panic, and Recover
- [**Ch 16**](16-testing-in-go.md) — Testing in Go — Unit Tests, Table Tests, Benchmarks
- [**Ch 17**](17-mini-project-cli-task-manager.md) — 🔨 Mini Project 1: CLI Task Manager

---

## Volume 2: Advanced Go (Ch 18–28)

*Concurrency, generics, the Go runtime — Go's most powerful and distinctive features.*

- [**Ch 18**](18-goroutines.md) — Goroutines — Lightweight Concurrency
- [**Ch 19**](19-channels.md) — Channels — Communicating Between Goroutines
- [**Ch 20**](20-select-timeouts-nonblocking.md) — Select, Timeouts, and Non-Blocking Operations
- [**Ch 21**](21-sync-package.md) — sync Package — Mutex, WaitGroup, Once, RWMutex
- [**Ch 22**](22-context-and-cancellation.md) — Context — Cancellation, Deadlines, Values
- [**Ch 23**](23-generics.md) — Generics in Go
- [**Ch 24**](24-packages-modules-workspace.md) — Packages, Modules, and Go Workspace
- [**Ch 25**](25-go-runtime-and-memory.md) — Go Runtime — Scheduler, Garbage Collector, Memory Model
- [**Ch 26**](26-profiling-and-performance.md) — Profiling and Benchmarking with pprof
- [**Ch 27**](27-reflection.md) — Reflection and the unsafe Package
- [**Ch 28**](28-mini-project-2-concurrent-log-processor.md) — 🔨 Mini Project 2: Concurrent Log Processor

---

## Volume 3: Data Structures (Ch 29–44)

*Build every important data structure from scratch in Go.*

- [**Ch 29**](29-complexity-analysis.md) — Complexity Analysis — Big O, Time and Space
- [**Ch 30**](30-slices-deep-dive.md) — Slices Deep Dive — Internals, Patterns, and Pitfalls
- [**Ch 31**](31-linked-lists.md) — Linked Lists — Singly and Doubly Linked
- [**Ch 32**](32-stacks-and-queues.md) — Stacks
- **Ch 33** — Queues and Deques
- [**Ch 34**](34-hash-tables.md) — Hash Maps and Hash Sets — How They Work Internally
- [**Ch 35**](35-trees-and-bst.md) — Binary Trees and Tree Traversals
- **Ch 36** — Binary Search Trees
- [**Ch 37**](37-balanced-trees.md) — Balanced Trees — AVL and Red-Black Trees
- [**Ch 38**](38-heaps-and-priority-queues.md) — Heaps and Priority Queues
- [**Ch 39**](39-graphs.md) — Graphs — Representation and Terminology
- [**Ch 40**](40-tries-and-advanced-ds.md) — Tries
- [**Ch 41**](41-bloom-filters.md) — Bloom Filters and Probabilistic Data Structures
- [**Ch 42**](42-lfu-cache.md) — LRU and LFU Cache
- [**Ch 43**](43-skip-lists.md) — Skip Lists and Concurrent Data Structures
- [**Ch 44**](44-mini-project-3-in-memory-key-value-store.md) — 🔨 Mini Project 3: In-Memory Key-Value Store

---

## Volume 4: Algorithms (Ch 45–58)

*From sorting to dynamic programming to graph algorithms — implement them all.*

- [**Ch 45**](45-sorting-algorithms.md) — Sorting Algorithms — Bubble to QuickSort to TimSort
- [**Ch 46**](46-searching-and-binary-search.md) — Binary Search and Its Variants
- [**Ch 47**](47-recursion-and-backtracking.md) — Recursion and Backtracking
- [**Ch 48**](48-two-pointers-sliding-window.md) — Two Pointers, Sliding Window, and Prefix Sum
- [**Ch 49**](49-divide-and-conquer.md) — Divide and Conquer
- [**Ch 50**](50-dynamic-programming.md) — Dynamic Programming — Fundamentals (1D DP)
- [**Ch 51**](51-dynamic-programming-advanced-patterns.md) — Dynamic Programming — Advanced Patterns (2D DP, Memoization)
- [**Ch 52**](52-greedy-algorithms.md) — Greedy Algorithms
- **Ch 53** — Graph Algorithms — BFS, DFS, Topological Sort
- [**Ch 54**](54-bellman-ford-floyd-warshall.md) — Shortest Path — Dijkstra, Bellman-Ford, Floyd-Warshall
- [**Ch 55**](55-minimum-spanning-trees.md) — Minimum Spanning Trees — Kruskal and Prim
- [**Ch 56**](56-string-algorithms.md) — String Algorithms — KMP, Rabin-Karp, Z-Algorithm
- [**Ch 57**](57-bit-manipulation.md) — Bit Manipulation
- [**Ch 58**](58-mini-project-4-route-planner.md) — 🔨 Mini Project 4: Route Planner with Dijkstra

---

## Volume 5: Web Development (Ch 59–72)

*Build real HTTP servers, REST APIs, gRPC services, and WebSocket apps.*

- [**Ch 59**](59-http-fundamentals.md) — HTTP Fundamentals and Go's net/http Package
- [**Ch 60**](60-building-rest-apis.md) — Building REST APIs with chi Router
- [**Ch 61**](61-middleware.md) — Middleware — Logging, Authentication, Recovery, Tracing
- [**Ch 62**](62-json-validation-serialization.md) — JSON, Validation, and Serialization
- [**Ch 63**](63-openapi-and-documentation.md) — OpenAPI and Code Generation with oapi-codegen
- [**Ch 64**](64-authentication-and-jwt.md) — Authentication — JWT, Sessions, and OAuth2
- [**Ch 65**](65-websockets.md) — WebSockets with gorilla/websocket
- [**Ch 66**](66-grpc.md) — gRPC with Protocol Buffers
- [**Ch 67**](67-file-uploads-multipart.md) — File Uploads, Multipart Forms, and Static Files
- [**Ch 68**](68-rate-limiting-throttling.md) — Rate Limiting and Throttling
- [**Ch 69**](69-cors-security-headers.md) — CORS and Security Headers
- [**Ch 70**](70-testing-http-services.md) — Testing HTTP Handlers — Unit, Integration, E2E
- [**Ch 71**](71-http-client.md) — HTTP Client — Calling External APIs Reliably
- [**Ch 72**](72-mini-project-5-blog-platform.md) — 🔨 Mini Project 5: Blog Platform REST API

---

## Volume 6: Databases and Persistence (Ch 73–84)

*PostgreSQL, Redis, MongoDB, OpenSearch — store and query data like a professional.*

- [**Ch 73**](73-postgresql-and-sqlx.md) — PostgreSQL Fundamentals for Go Developers
- **Ch 74** — database/sql — The Standard Library Way
- [**Ch 75**](75-sqlc-type-safe-sql.md) — SQLC — Type-Safe SQL Queries from Go
- [**Ch 76**](76-database-migrations.md) — Database Migrations with golang-migrate
- [**Ch 77**](77-repository-pattern.md) — The Repository Pattern and Database Abstraction
- [**Ch 78**](78-transactions-locking-concurrency.md) — Transactions, Locking, and Concurrency Control
- [**Ch 79**](79-advanced-postgresql.md) — Advanced PostgreSQL — Full-Text Search, JSONB, Indexes
- [**Ch 80**](80-redis-caching-pubsub.md) — Redis — Caching, Pub/Sub, and Data Structures
- [**Ch 81**](81-mongodb-go-driver.md) — MongoDB with the Official Go Driver
- [**Ch 82**](82-opensearch-elasticsearch.md) — OpenSearch and Elasticsearch Integration
- [**Ch 83**](83-time-series-databases.md) — Time-Series Databases — InfluxDB and TimescaleDB
- [**Ch 84**](84-mini-project-6-product-catalog.md) — 🔨 Mini Project 6: Product Catalog with Full-Text Search

---

## Volume 7: Clean Architecture and Design Patterns (Ch 85–94)

*Write code that is maintainable, testable, and scales to large teams.*

- [**Ch 85**](85-clean-architecture.md) — Clean Architecture — Layers and Dependency Rule
- [**Ch 86**](86-domain-driven-design.md) — Domain-Driven Design (DDD) in Go
- [**Ch 87**](87-dependency-injection-wire.md) — Dependency Injection — Manual and with Wire
- [**Ch 88**](88-cqrs-event-sourcing.md) — CQRS — Command Query Responsibility Segregation
- [**Ch 89**](89-event-sourcing.md) — Event Sourcing — State as a Stream of Events
- [**Ch 90**](90-outbox-saga-patterns.md) — The Outbox Pattern — Reliable Event Publishing
- [**Ch 91**](91-saga-pattern.md) — Saga Pattern — Distributed Transactions
- **Ch 92** — Repository, Unit of Work, and Other Patterns
- [**Ch 93**](93-config-management-viper.md) — Configuration Management with Viper
- [**Ch 94**](94-major-project-1-food-delivery.md) — 🏗️ Major Project 1: Food Delivery Backend (Modular Monolith)

---

## Volume 8: Asynchronous and Event-Driven Systems (Ch 95–106)

*Go beyond request-response — build systems that react, retry, and recover.*

- [**Ch 95**](95-async-systems-overview.md) — Async Systems — Why, When, and Trade-offs
- [**Ch 96**](96-worker-pools-async-patterns.md) — Worker Pools, Job Queues, and Pipelines
- **Ch 97** — Redis Streams — Lightweight Event Streaming
- [**Ch 98**](98-kafka-fundamentals.md) — Apache Kafka with Go (kafka-go, Sarama)
- [**Ch 99**](99-rabbitmq.md) — RabbitMQ with Go (amqp091-go)
- [**Ch 100**](100-nats-jetstream.md) — NATS and NATS JetStream
- [**Ch 101**](101-watermill.md) — Watermill — The Go Event-Driven Framework
- [**Ch 102**](102-at-least-once-delivery-idempotency.md) — At-Least-Once Delivery and Idempotency
- [**Ch 103**](103-message-ordering-partitioning.md) — Message Ordering and Partitioning
- [**Ch 104**](104-asynq-task-queues.md) — Background Jobs and Task Queues with asynq
- [**Ch 105**](105-scheduled-tasks-cron.md) — Scheduled Tasks, Cron Jobs, and Distributed Schedulers
- [**Ch 106**](106-major-project-2-ticket-booking.md) — 🏗️ Major Project 2: Ticket Booking System (Event-Driven)

---

## Volume 9: System Architecture (Ch 107–117)

*Monoliths, microservices, distributed systems — understand the trade-offs and build the right system.*

- [**Ch 107**](107-monolithic-architecture.md) — Monolithic Architecture — When It's the Right Choice
- [**Ch 108**](108-microservices-patterns.md) — Microservices — Principles, Patterns, and Pitfalls
- [**Ch 109**](109-service-communication-rest-grpc-events.md) — Service Communication — REST vs gRPC vs Events
- **Ch 110** — Service Discovery and Load Balancing
- **Ch 111** — API Gateway with Traefik / Kong
- **Ch 112** — Circuit Breaker, Retry, and Bulkhead Patterns
- [**Ch 113**](113-distributed-caching-consistent-hashing.md) — Distributed Caching and Consistent Hashing
- **Ch 114** — Distributed Locks with Redis and etcd
- [**Ch 115**](115-cap-theorem-distributed-systems.md) — CAP Theorem, Consistency Models, and PACELC
- [**Ch 116**](116-data-partitioning-sharding.md) — Data Partitioning, Sharding, and Multi-Tenancy
- **Ch 117** — 🏗️ Major Project 3: E-Commerce Platform (Microservices)

---

## Volume 10: Observability (Ch 118–127)

*See inside your running system — logs, metrics, traces, alerts.*

- [**Ch 118**](118-observability-logging-metrics.md) — Observability — The Three Pillars (Logs, Metrics, Traces)
- [**Ch 119**](119-structured-logging-slog-zerolog-zap.md) — Structured Logging with slog, zerolog, and zap
- **Ch 120** — Prometheus — Metrics Collection and Exposition
- [**Ch 121**](121-grafana-dashboards.md) — Grafana — Dashboards, Panels, and Alerts
- **Ch 122** — OpenTelemetry — Distributed Tracing in Go
- [**Ch 123**](123-kibana-opensearch-logs.md) — Kibana and OpenSearch — Log Management
- [**Ch 124**](124-sentry-error-tracking.md) — Sentry — Error Tracking and Performance Monitoring
- [**Ch 125**](125-datadog-apm.md) — Datadog — APM, Infrastructure, and Logs
- [**Ch 126**](126-alerting-pagerduty-on-call.md) — Alerting — Rules, Thresholds, PagerDuty, On-Call
- [**Ch 127**](127-mini-project-7-observable-api.md) — 🔨 Mini Project 7: Fully Observable API Service

---

## Volume 11: Production Engineering (Ch 128–136)

*Ship and operate your Go applications in production.*

- [**Ch 128**](128-docker-kubernetes.md) — Docker and Containers for Go Applications
- [**Ch 129**](129-kubernetes-for-backend-engineers.md) — Kubernetes for Backend Engineers
- [**Ch 130**](130-cicd-security-performance.md) — CI/CD with GitHub Actions
- [**Ch 131**](131-secrets-management.md) — Secrets Management — Vault, AWS Secrets Manager, Doppler
- **Ch 132** — Security Best Practices for Go Applications
- **Ch 133** — Performance Testing with k6 and vegeta
- [**Ch 134**](134-high-availability-graceful-shutdown.md) — High Availability, Graceful Shutdown, and Zero-Downtime Deploys
- [**Ch 135**](135-database-at-scale.md) — Database at Scale — Read Replicas, pgBouncer, Connection Pooling
- [**Ch 136**](136-final-capstone-ecommerce.md) — 🏆 Final Major Project: Production SaaS Backend

---

## Course Statistics

| Metric | Value |
|--------|-------|
| Total Chapters | 137 (Ch 00–136) |
| Volumes | 12 (Vol 0–11) |
| Mini Projects | 7 |
| Major Projects | 3 |
| Final Capstone | 1 |
| Total Projects | 11 |

---

## How to Use This Course

**Each chapter includes:**
- Clear explanation with analogies
- Code examples you can run
- ASCII diagrams where helpful
- ✅ Quick Checks after each section
- Chapter Summary
- Easy / Medium / Hard exercises

**Projects include:**
- Complete requirements
- Step-by-step guidance
- Starter code structure
- What to implement
- How to verify it works
- Extension challenges

**Recommended pace:**
- Volume 1 (Fundamentals): 2 weeks
- Volume 2 (Advanced Go): 1 week
- Volumes 3–4 (DSA): 3–4 weeks
- Volumes 5–6 (Web + Databases): 3 weeks
- Volumes 7–8 (Architecture + Async): 3 weeks
- Volumes 9–10 (Systems + Observability): 3 weeks
- Volume 11 (Production): 2 weeks
- **Total: ~18 weeks** (4–5 months at 2 hours/day)

---

## Prerequisites

- You know **how to use a computer** and **have some programming experience** (any language is fine — Python, JavaScript, Java, etc.)
- You do **not** need to know Go
- You do **not** need to know backend engineering or distributed systems

## What You Will Know at the End

- Write idiomatic, production-quality Go code
- Implement all common data structures and algorithms
- Build REST APIs, gRPC services, and WebSocket servers
- Design and build clean, testable backend systems
- Work with PostgreSQL, Redis, MongoDB, and OpenSearch
- Build event-driven systems with Kafka, RabbitMQ, and NATS
- Design monolithic and microservices architectures
- Instrument your system with Prometheus, Grafana, OpenTelemetry
- Monitor errors with Sentry, APM with Datadog
- Ship and operate Go applications in Docker and Kubernetes
- Write CI/CD pipelines and handle production incidents
