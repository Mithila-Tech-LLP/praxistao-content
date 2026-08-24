# Chapter 37: When to Use Which Database — The Grand Decision Framework

You've now learned 10 different databases. The most important skill isn't knowing each one — it's knowing which one to pick. Wrong choice = years of pain. Right choice = system runs smoothly at scale.

## Table of Contents

1. The Decision Framework
2. By Data Shape
3. By Access Pattern
4. By Scale Requirement
5. By Team and Operational Constraints
6. Real-World Scenarios
7. The Multi-Database Architecture
8. Summary Decision Matrix

---

## 1. The Decision Framework

Ask these questions in order:

```
1. What shape is my data?
   → Relational (tables with FK relationships)?
   → Documents (nested, variable structure)?
   → Time-series / append-only?
   → Graph (highly connected)?
   → Vectors (similarity search)?
   → Key-value (simple lookups)?

2. What's my primary access pattern?
   → Complex analytical queries?
   → Simple CRUD by primary key?
   → Similarity search?
   → Write-heavy (millions/sec)?
   → Read-heavy with complex filters?

3. How much data? How fast does it grow?
   → < 100 GB: Almost anything works
   → 100 GB - 10 TB: Need to think about indexes and partitioning
   → > 10 TB: Need distributed storage from day one

4. What are your operational constraints?
   → Team familiarity (don't pick Neo4j if nobody knows it)
   → Ops complexity tolerance (managed cloud vs self-hosted)
   → Budget (open source vs commercial)
```

---

## 2. By Data Shape

### Relational Data → PostgreSQL or MySQL

Your data has clear entities with relationships, foreign keys, and you query across multiple entities:

```
Users have Orders have LineItems have Products
```

```
Use PostgreSQL when:
✓ You need ACID transactions across tables
✓ Complex JOIN queries
✓ Row-level security
✓ Need PostGIS, pgvector, or other extensions
✓ Not sure yet — PostgreSQL is the safe default

Use MySQL when:
✓ Your team knows MySQL deeply
✓ You need MySQL-specific replication features
✓ WordPress, Magento, or other MySQL-centric ecosystems
```

### Document Data → MongoDB or PostgreSQL JSONB

Your data is hierarchical, variable-structure, or changes schema frequently:

```
Blog posts with arbitrary metadata
Product catalogs where each category has different attributes
User settings with deeply nested preferences
```

```
Use MongoDB when:
✓ Schema changes frequently and unpredictably
✓ Data is naturally document-shaped (deeply nested)
✓ You need flexible queries on arbitrary fields
✓ Team knows MongoDB

Use PostgreSQL JSONB when:
✓ Mostly relational but some fields need document flexibility
✓ You want one database for everything
✓ Need transactional consistency with relational data
```

### Time-Series / Append-Only → ClickHouse, Cassandra, or TimescaleDB

Data is generated continuously with timestamps, rarely updated:

```
Metrics, logs, IoT sensor readings, financial tick data, events
```

```
Use ClickHouse when:
✓ Analytical queries (GROUP BY, aggregations, time windows)
✓ Business intelligence / dashboards
✓ Less than ~1 billion writes/day per node

Use Cassandra when:
✓ Write throughput is extreme (millions/second)
✓ Data must span multiple data centers
✓ High availability is non-negotiable
✓ You can design specific query patterns upfront

Use TimescaleDB (PostgreSQL extension) when:
✓ Time-series + SQL + familiar tooling
✓ Continuous aggregates (automatic roll-ups)
✓ Your scale fits one server
```

### Graph Data → SurrealDB, Neo4j, or PostgreSQL

Highly connected data where traversing relationships is the main query pattern:

```
Social networks, recommendation engines, fraud detection, knowledge graphs
```

```
Use SurrealDB when:
✓ New project, want graph + document + relational in one
✓ Real-time apps (live queries)

Use Neo4j when:
✓ Graph traversal is the core product feature
✓ Need the mature Cypher query language
✓ Graph algorithms (shortest path, community detection)
```

### Vector Data → pgvector or Qdrant

Similarity search, semantic search, AI embeddings:

```
Use pgvector when:
✓ Already using PostgreSQL
✓ < 10 million vectors
✓ Want to combine with relational data in one query

Use Qdrant when:
✓ Vectors are the primary data model
✓ > 10 million vectors
✓ Need advanced filtering + payload search
```

---

## 3. By Access Pattern

### Read-Heavy, Complex Queries
→ **PostgreSQL** with proper indexes. Add read replicas if needed.

### Write-Heavy, Simple Lookups
→ **Cassandra** (millions writes/sec, no complex queries) or **Redis** (microsecond reads/writes for hot data).

### Caching / Sessions / Ephemeral Data
→ **Redis**. Always. No debate.

### Full-Text Search
→ **PostgreSQL** (pg_trgm, tsvector) for small scale. **Elasticsearch** or **Typesense** for dedicated search.

### Leaderboards / Sorted Sets
→ **Redis** with sorted sets.

### Pub/Sub / Message Queue
→ **Redis** (small scale, simple) or **Kafka/RabbitMQ** (see Part 2 of this course).

---

## 4. By Scale Requirement

```
Tier 1: Startup (< 10 GB, < 10K users)
→ PostgreSQL for everything. Add Redis for caching when slow.
→ Run on a single $50/month server.
→ Don't over-engineer.

Tier 2: Growing (10 GB - 1 TB, 100K users)
→ PostgreSQL with read replicas.
→ Redis for sessions and hot data.
→ Consider ClickHouse if you're building analytics.
→ PgBouncer for connection pooling.

Tier 3: Scale (1 TB - 100 TB, millions of users)
→ PostgreSQL with sharding or Citus extension for primary data.
→ Cassandra for write-heavy append data.
→ ClickHouse for analytics.
→ Redis Cluster for caching.
→ Qdrant/pgvector for AI features.

Tier 4: Hyperscale (> 100 TB, 100M+ users)
→ Custom solutions. You're Google/Amazon/Meta territory.
→ DynamoDB, Bigtable, Spanner — managed cloud-scale databases.
```

---

## 5. Real-World Scenarios

### Scenario A: E-Commerce Platform

```
Orders, users, inventory → PostgreSQL (ACID transactions critical)
Product search → PostgreSQL full-text or Elasticsearch
Sessions, cart → Redis (TTL, fast access)
Product recommendations → pgvector (similar products)
Analytics dashboard → ClickHouse (order history, revenue trends)
```

### Scenario B: Social Media App

```
Users, posts, core data → PostgreSQL
Social graph (follows, friends) → PostgreSQL with proper indexes (< 100M users)
  or → Neo4j/SurrealDB (> 100M users with complex graph queries)
Feed cache → Redis sorted sets
Real-time notifications → Redis Pub/Sub or Kafka
Media metadata → MongoDB (variable structure)
Search → Elasticsearch
```

### Scenario C: IoT Platform

```
Device registry → PostgreSQL
Sensor readings (time-series) → ClickHouse or Cassandra
Real-time alerts → Redis streams
Device state (latest reading) → Redis hash
Long-term archive → ClickHouse
```

### Scenario D: AI Chatbot

```
User accounts → PostgreSQL
Conversation history → MongoDB (nested documents, variable length)
Document embeddings → pgvector or Qdrant
LLM response cache → Redis (expensive calls cached by query hash)
Analytics → ClickHouse
```

---

## 6. The Multi-Database Architecture

Almost every production system uses 2-5 databases. That's normal. The key is: **use each database for what it's best at.**

```
┌─────────────────────────────────────────────────────────────────┐
│                    Your Application                             │
└──────────────────────────┬──────────────────────────────────────┘
          ┌────────────────┼────────────────────┐
          ▼                ▼                    ▼
    ┌───────────┐    ┌──────────┐         ┌──────────┐
    │PostgreSQL │    │  Redis   │         │ClickHouse│
    │(source of │    │(cache,   │         │(analytics│
    │  truth)   │    │sessions) │         │ queries) │
    └───────────┘    └──────────┘         └──────────┘
          │
    ┌─────┴─────┐
    ▼           ▼
┌────────┐ ┌─────────┐
│pgvector│ │  Kafka  │
│ (AI    │ │(events) │
│ search)│ │         │
└────────┘ └─────────┘
```

**The golden rule:** PostgreSQL (or MySQL) is your **source of truth**. Everything else is derived data that can be rebuilt from PostgreSQL if needed.

---

## Summary Decision Matrix

| Need | Primary Choice | Alternative |
|------|---------------|-------------|
| Relational, ACID | PostgreSQL | MySQL |
| Documents | MongoDB | PostgreSQL JSONB |
| Analytics | ClickHouse | TimescaleDB |
| High write scale | Cassandra | DynamoDB |
| Cache / sessions | Redis | Memcached |
| Vector search | pgvector | Qdrant |
| Multi-model / graph | SurrealDB | Neo4j |
| Full-text search | PostgreSQL tsvector | Elasticsearch |
| Unknown / new project | PostgreSQL | Add others when you hit limits |

**The most important lesson:** Start with PostgreSQL. Migrate to specialized databases only when PostgreSQL demonstrably can't handle your specific bottleneck. Premature optimization with exotic databases kills startups.

### Exercises

**Easy:** For each of these scenarios, name the best primary database and one reason why: (1) A bank's transaction ledger, (2) A Spotify-like music streaming app's song metadata, (3) A weather station collecting temperature every second.

**Medium:** Design the database architecture for a food delivery app (Uber Eats / DoorDash clone). List each entity/data type and which database you'd use. Justify each choice.

**Hard:** You're CTO of a startup. You have 50K daily active users, 200 GB of PostgreSQL data, and the site is slow. Profile which queries are slow. Design an architecture change (add Redis? add read replicas? move analytics to ClickHouse?) with cost and complexity estimates.
