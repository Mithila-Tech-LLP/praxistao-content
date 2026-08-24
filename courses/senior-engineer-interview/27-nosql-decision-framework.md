# Chapter 27: NoSQL Decision Framework — MongoDB, Redis, DynamoDB & Cassandra

Knowing when to use SQL vs NoSQL, and which NoSQL to choose, is a senior engineer skill. This chapter gives you a framework for that decision plus deep internals for each major NoSQL database.

## Table of Contents

1. [SQL vs NoSQL — When to Use Which](#1-sql-vs-nosql--when-to-use-which)
2. [Redis — In-Memory Data Structure Store](#2-redis--in-memory-data-structure-store)
3. [MongoDB — Document Database](#3-mongodb--document-database)
4. [DynamoDB — Managed Key-Value/Document at Scale](#4-dynamodb--managed-key-valuedocument-at-scale)
5. [Cassandra — Wide-Column for Massive Write Scale](#5-cassandra--wide-column-for-massive-write-scale)
6. [Decision Matrix](#6-decision-matrix)
7. [Interview Questions & Model Answers](#7-interview-questions--model-answers)
8. [Summary](#summary)

---

## 1. SQL vs NoSQL — When to Use Which

**Use PostgreSQL/SQL when:**
- Your data is relational and you need JOINs
- You need transactions (ACID guarantees across multiple rows/tables)
- Schema is relatively stable
- You need complex queries (window functions, aggregations, CTEs)
- Compliance/audit trails require strong consistency

**Use NoSQL when:**
- Massive write scale (millions of writes/second)
- Flexible/evolving schemas (user-generated content, logs, events)
- Specific access patterns that SQL handles poorly (time-series, caching, search)
- Horizontal sharding is required from day one

The most common mistake: reaching for NoSQL because "it scales better" before proving you've exhausted what your SQL database can do with proper indexing and caching.

---

## 2. Redis — In-Memory Data Structure Store

Redis stores everything in RAM, making it 10-100x faster than disk-based databases for reads and writes. It's used as a cache, message broker, session store, rate limiter, and distributed lock.

### Data Structures

```
String:     SET user:1:name "Alice"     GET user:1:name
            INCR page:views             SETEX session:abc 3600 "{...}"  (TTL)

List:       RPUSH queue:jobs "job1"     LPOP queue:jobs
            LRANGE feed:user:1 0 9      (recent 10 items)

Hash:       HSET user:1 name Alice age 30   HGETALL user:1
            HINCRBY user:1 login_count 1

Set:        SADD followers:user1 user2 user3    SISMEMBER followers:user1 user2
            SUNION followers:user1 followers:user2   (union of sets)

Sorted Set: ZADD leaderboard 1500 "alice"   ZRANGEBYSCORE leaderboard 0 2000
            ZRANK leaderboard "alice"        ZREVRANGE leaderboard 0 9  (top 10)

Bitmap:     SETBIT active_users:2024-01-01 userid 1    BITCOUNT active_users:2024-01-01
            (used for tracking daily active users with very low memory)
```

### Redis as a Cache — Cache Strategies

```
Cache-Aside (most common):
  1. Read: check Redis first
  2. If miss: read from PostgreSQL, write to Redis with TTL
  3. On write: update PostgreSQL, invalidate Redis key
  
Write-Through:
  Write to Redis AND PostgreSQL synchronously. Redis always has fresh data. Slower writes.

Write-Behind (Write-Back):
  Write to Redis only, flush to PostgreSQL asynchronously. Faster writes, risk of data loss.
```

```go
// Cache-aside pattern in Go
func getUser(ctx context.Context, rdb *redis.Client, db *sql.DB, id int) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // 1. Try cache
    cached, err := rdb.Get(ctx, key).Result()
    if err == nil {
        var u User
        json.Unmarshal([]byte(cached), &u)
        return &u, nil
    }
    
    // 2. Cache miss: read from DB
    u := &User{}
    db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Name, &u.Email)
    
    // 3. Write to cache with TTL
    data, _ := json.Marshal(u)
    rdb.Set(ctx, key, data, 15*time.Minute)
    
    return u, nil
}
```

### Redis Persistence

Redis is in-memory but offers persistence options:
- **RDB (Redis Database):** point-in-time snapshots. Fast restart but potential data loss since last snapshot.
- **AOF (Append-Only File):** logs every write command. Slower but more durable.
- **No persistence:** pure cache, data is lost on restart. Acceptable if your source of truth is elsewhere.

### Redis Cluster

Redis Cluster shards keys across multiple nodes using consistent hashing:
```
Keys are distributed across 16384 slots
hash_slot = CRC16(key) % 16384
Each node owns a range of slots
Automatic failover with replicas
```

---

## 3. MongoDB — Document Database

MongoDB stores JSON-like BSON documents. Each document can have a different structure. Good for user-generated content, catalogs, and content management.

### Document Model

```javascript
// A MongoDB document:
{
  "_id": ObjectId("..."),
  "name": "Alice",
  "email": "alice@example.com",
  "address": {
    "city": "New York",
    "zip": "10001"
  },
  "tags": ["premium", "early-adopter"],
  "orders": [              // embedded subdocuments (denormalization)
    { "product": "iPhone", "total": 999 },
    { "product": "AirPods", "total": 199 }
  ]
}
```

### Indexing & Queries

```javascript
// Find with filter:
db.users.find({ "address.city": "New York", "tags": "premium" })

// Create index:
db.users.createIndex({ email: 1 })                    // single field
db.users.createIndex({ "address.city": 1, name: 1 }) // compound
db.users.createIndex({ tags: 1 })                    // multikey (array field)

// Aggregation pipeline (equivalent to SQL GROUP BY + JOIN):
db.orders.aggregate([
  { $match: { status: "completed" } },                // WHERE
  { $group: { _id: "$user_id", total: { $sum: "$amount" } } }, // GROUP BY, SUM
  { $sort: { total: -1 } },                          // ORDER BY
  { $limit: 10 }                                     // LIMIT
])

// Lookup (LEFT JOIN equivalent):
db.orders.aggregate([
  {
    $lookup: {
      from: "users",
      localField: "user_id",
      foreignField: "_id",
      as: "user"
    }
  }
])
```

### MongoDB Limitations

- No multi-document transactions until version 4.0 (and still have overhead)
- JOINs ($lookup) are expensive — MongoDB is designed for denormalized data
- Schema validation is optional (a feature AND a danger)
- Not great for complex analytics that need arbitrary aggregations

---

## 4. DynamoDB — Managed Key-Value/Document at Scale

DynamoDB is AWS's fully managed NoSQL. It scales to any load, has predictable single-digit millisecond latency, and requires zero operational overhead. The trade-off: you must design your access patterns upfront.

### Key Concepts

```
Table: collection of items
Item: a row (like a JSON document)
Primary Key: identifies each item uniquely
  - Partition Key (PK): determines which partition holds the item
  - Sort Key (SK): secondary ordering within a partition (optional)
```

```python
# Item example:
{
  "PK": "USER#alice",          # partition key
  "SK": "PROFILE",             # sort key
  "name": "Alice",
  "email": "alice@example.com"
}

{
  "PK": "USER#alice",
  "SK": "ORDER#2024-01-15#abc123",
  "total": 99.99,
  "status": "shipped"
}

# Both items share the same partition key (USER#alice)
# Querying: all orders for alice → PK = "USER#alice", SK begins_with "ORDER#"
# This pattern is called "single-table design"
```

### Single-Table Design

The DynamoDB best practice: put all entity types in ONE table, using generic PK/SK names:

```
ACCESS PATTERNS → DRIVE TABLE DESIGN (not the other way around!)

Pattern 1: Get user profile           → PK=USER#id, SK=PROFILE
Pattern 2: Get user's orders           → PK=USER#id, SK begins_with ORDER#
Pattern 3: Get order by ID             → PK=ORDER#id, SK=METADATA  
Pattern 4: Get all orders for a date   → GSI: PK=DATE#2024-01-15, SK=ORDER#id

GSI (Global Secondary Index): allows querying on different attribute combinations
```

---

## 5. Cassandra — Wide-Column for Massive Write Scale

Cassandra is designed for write-heavy workloads at massive scale. Instagram, Netflix, and Apple use it to store billions of events per day. It uses a distributed ring architecture with no single point of failure.

### Data Model

```sql
-- Cassandra tables are designed around query patterns:
CREATE TABLE user_events (
    user_id UUID,
    event_time TIMESTAMP,
    event_type TEXT,
    data TEXT,
    PRIMARY KEY (user_id, event_time)  -- user_id = partition key, event_time = clustering key
) WITH CLUSTERING ORDER BY (event_time DESC);

-- This allows:
SELECT * FROM user_events WHERE user_id = ? LIMIT 100;
-- But NOT:
SELECT * FROM user_events WHERE event_type = 'click';  -- no partition key!
```

### How Cassandra Handles Writes

```
Write path:
1. Client sends write to any node (coordinator)
2. Coordinator hashes the partition key to find the ring position
3. Write is sent to N=3 replica nodes (replication factor)
4. Write is acknowledged when quorum (QUORUM = majority) of replicas respond
5. Each replica writes to a commit log (for durability) and a memtable (in-memory)
6. Memtable is periodically flushed to an SSTable file on disk

QUORUM writes + QUORUM reads = strong consistency
ONE writes + ONE reads = eventual consistency (higher throughput)
```

### Cassandra vs PostgreSQL

| | Cassandra | PostgreSQL |
|---|---|---|
| Write throughput | Millions/sec | ~10K-100K/sec |
| Query flexibility | Limited (partition key required) | Full SQL |
| Consistency | Tunable (eventual by default) | ACID |
| Joins | None | Full |
| Best for | Time-series, event logs, IoT | Transactional data, reporting |

---

## 6. Decision Matrix

| Requirement | Best Choice |
|---|---|
| Caching, sessions, leaderboards, pub/sub | Redis |
| Transactional data, complex queries, reporting | PostgreSQL |
| User-generated content, flexible schema, catalogs | MongoDB |
| Managed, serverless, AWS-native, massive scale | DynamoDB |
| Massive write scale (billions/day), time-series | Cassandra |
| Full-text search | Elasticsearch |
| Graph relationships | Neo4j |
| Time-series metrics | InfluxDB or TimescaleDB |

---

## 7. Interview Questions & Model Answers

**Q: When would you choose Redis over a traditional database?**

"Redis is ideal as a cache layer in front of a primary database — storing frequently-read data in RAM reduces latency by 100x. It's also great for ephemeral data with TTLs (sessions, OTPs, rate limit counters), leaderboards using sorted sets, pub/sub messaging, and distributed locks. I wouldn't use Redis as my primary database for critical data unless it's combined with persistence (AOF/RDB) and replication, because it's an in-memory system and RAM is expensive."

**Q: What is DynamoDB single-table design and why is it recommended?**

"In DynamoDB, each query requires the partition key, and secondary indexes are expensive. Single-table design puts all entity types in one table with generic PK/SK columns. For example, PK='USER#123' SK='PROFILE' holds user data, while PK='USER#123' SK='ORDER#2024-01...' holds orders. This lets you retrieve a user and their recent orders in a single query using the `begins_with` operator on the sort key. You design the table around your access patterns rather than your entity types — the opposite of relational design."

---

## Summary

- **SQL:** relational data, transactions, complex queries, stable schema. Don't abandon it prematurely.
- **Redis:** in-memory, microsecond latency. Cache, sessions, leaderboards, pub/sub, rate limiting.
- **MongoDB:** flexible JSON documents, no JOINs needed, evolving schema. Good for content.
- **DynamoDB:** managed AWS NoSQL. Design around access patterns (single-table design). Scales to any load.
- **Cassandra:** extreme write scale, time-series data. Partition key required for all queries.
- Decision order: Prove PostgreSQL can't do it → pick the NoSQL that matches your specific access pattern.
