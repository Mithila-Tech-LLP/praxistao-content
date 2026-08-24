---
title: Production Readiness
---
Writing an API that works on your laptop is the easy part. This section covers what it takes to run that API in production: clean architecture that survives change, observability so you know when something breaks, and the infrastructure to actually deploy it.

### Clean Architecture & Dependency Injection
As a service grows past a few handlers, structure starts to matter. Clean architecture keeps your business logic independent of your database and HTTP framework, so either one can change without a rewrite.

**Resources:**
- [Clean Architecture](course:go-programming#85-clean-architecture)
- [Dependency Injection with Wire](course:go-programming#87-dependency-injection-wire)

### Async Patterns: Queues & Kafka
> optional

Not every piece of work needs to happen inside the request/response cycle. Message queues and event logs like Kafka let you decouple slow or bursty work from the API that triggered it.

**Resources:**
- [Async Systems Overview](course:go-programming#95-async-systems-overview)
- [Kafka Fundamentals](course:go-programming#98-kafka-fundamentals)

### Observability: Logging, Metrics, Tracing
When something breaks in production, you don't get a debugger — you get logs, metrics, and dashboards. Learn to build a service that tells you what it's doing.

**Resources:**
- [Observability: Logging and Metrics](course:go-programming#118-observability-logging-metrics)
- [Grafana Dashboards](course:go-programming#121-grafana-dashboards)

### Docker & Kubernetes
Package your service so it runs the same way on your machine, your teammate's machine, and in production — then orchestrate many of them with Kubernetes.

**Resources:**
- [Docker & Kubernetes](course:go-programming#128-docker-kubernetes)

### System Design Fundamentals
Zoom out from a single service to how whole systems are put together — the tradeoffs behind scaling, caching, and reliability that show up in every senior engineering conversation.

**Resources:**
- [System Design Framework](course:senior-engineer-interview#35-system-design-framework)
- [CAP Theorem and Consistency](course:senior-engineer-interview#29-cap-theorem-consistency)

### Capstone: Ship a Production-Style Service
> branches-from: Docker & Kubernetes

Bring everything in this roadmap together in one project: an API, a real database, tests, and the deployment tooling to actually run it.

**Resources:**
- [Final Capstone: E-Commerce Platform](course:go-programming#136-final-capstone-ecommerce)
