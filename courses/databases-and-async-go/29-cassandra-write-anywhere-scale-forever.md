# Chapter 29: Cassandra — Write Anywhere, Scale Forever

Cassandra was built at Facebook to handle hundreds of millions of writes per day across data centers on different continents. It has no single point of failure, no primary/replica — every node is equal. This chapter explains how it achieves this and when to use it.

## Table of Contents

1. The Problem Cassandra Solves
2. Cassandra Architecture — The Peer-to-Peer Ring
3. Data Model — Keyspaces, Tables, Partition Keys
4. CQL — Cassandra Query Language
5. Consistency Levels
6. Cassandra Anti-Patterns
7. Building with Cassandra in Go
8. Exercises

---

## 1. The Problem Cassandra Solves

Imagine you're collecting IoT sensor readings: 10,000 sensors sending a temperature reading every second = 10,000 inserts/second. You need this data for the next 5 years.

PostgreSQL can handle this, but as data grows to billions of rows, problems appear:
- A single node becomes the bottleneck
- Adding nodes to PostgreSQL is complex and expensive
- If the primary goes down, writes stop until failover completes

Cassandra's solution: **no primary, no replicas — just nodes**. Every node accepts writes. If one node is down, the others continue without interruption. Need more capacity? Add a node — it starts accepting traffic immediately.

**Cassandra is designed for:**
- Very high write throughput (millions/second)
- Data that scales to petabytes
- Multi-data-center deployments
- No tolerance for single points of failure
- Append-mostly workloads (IoT, logs, activity streams, time-series)

---

## 2. Cassandra Architecture — The Peer-to-Peer Ring

```
               ┌──────────────────────────────────────────────────┐
               │              Cassandra Ring                       │
               │                                                   │
               │      Node 1              Node 2                  │
               │    (owns tokens          (owns tokens            │
               │     0 - 25%)             25% - 50%)              │
               │                                                   │
               │      Node 4              Node 3                  │
               │    (owns tokens          (owns tokens            │
               │     75% - 100%)          50% - 75%)              │
               │                                                   │
               └──────────────────────────────────────────────────┘
```

**Consistent hashing:** Each row's partition key is hashed to a number between 0 and 2^63. Each node owns a range of these numbers (token range). To find which node holds a row: hash the partition key and go to the node that owns that token.

**Replication factor (RF):** Data is stored on multiple nodes for fault tolerance. RF=3 means each row is on 3 nodes. If one node fails, the other two have the data.

**No primary:** Writes go to any node (the "coordinator"). The coordinator routes the write to the nodes responsible for that partition. All RF nodes receive the write — no leader election, no failover needed.

**Gossip protocol:** Nodes periodically share state with neighbors (like rumors spreading through a group). Each node knows the cluster topology without centralized coordination.

---

## 3. Data Model — Keyspaces, Tables, Partition Keys

Cassandra's data model is fundamentally different from SQL. **Design for your queries, not your entities.**

### Primary Key Components

```
PRIMARY KEY (partition_key, clustering_column1, clustering_column2)
```

- **Partition key:** Determines which node stores the data. All rows with the same partition key live on the same node.
- **Clustering columns:** Define the sort order within a partition. Data is physically stored in this order on disk.

### Example: IoT Sensor Data

```sql
-- Goal: "Get all readings from sensor 123 in the last hour, sorted by time"
CREATE TABLE sensor_readings (
    sensor_id   UUID,
    recorded_at TIMESTAMP,
    temperature FLOAT,
    humidity    FLOAT,
    PRIMARY KEY (sensor_id, recorded_at)
) WITH CLUSTERING ORDER BY (recorded_at DESC);
-- Data for sensor 123 is all on the same node, sorted by time → fast!
```

### Example: User Activity Log

```sql
-- Goal: "Get Alice's last 50 activities"
CREATE TABLE user_activities (
    user_id     UUID,
    activity_id TIMEUUID,  -- time-based UUID, auto-sorted by time
    event_type  TEXT,
    metadata    TEXT,
    PRIMARY KEY (user_id, activity_id)
) WITH CLUSTERING ORDER BY (activity_id DESC);
```

---

## 4. CQL — Cassandra Query Language

CQL looks like SQL but has strict limitations because of how data is distributed.

```sql
-- Create keyspace (database)
CREATE KEYSPACE myapp
WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'datacenter1': 3   -- 3 replicas in datacenter1
};

USE myapp;

-- Insert (Cassandra calls this "upsert" — same syntax for insert and update)
INSERT INTO sensor_readings (sensor_id, recorded_at, temperature, humidity)
VALUES (uuid(), toTimestamp(now()), 23.5, 60.2);

-- Query MUST include the full partition key
SELECT * FROM sensor_readings WHERE sensor_id = ? LIMIT 100;

-- Query with clustering column range
SELECT * FROM sensor_readings
WHERE sensor_id = ?
  AND recorded_at >= '2024-01-01'
  AND recorded_at <  '2024-01-02';

-- Update (Cassandra overwrites the columns you specify)
UPDATE sensor_readings
SET temperature = 24.0
WHERE sensor_id = ? AND recorded_at = ?;

-- Delete (Cassandra marks as tombstone, doesn't actually delete)
DELETE FROM sensor_readings
WHERE sensor_id = ? AND recorded_at = ?;

-- TTL: automatically expire rows after N seconds
INSERT INTO sensor_readings (sensor_id, recorded_at, temperature)
VALUES (?, ?, ?) USING TTL 86400;  -- expires after 24 hours
```

### What You Cannot Do in Cassandra

```sql
-- NO: full table scan without partition key
SELECT * FROM sensor_readings WHERE temperature > 25;  -- ERROR

-- NO: joins
SELECT * FROM sensor_readings JOIN sensors ON ...;  -- NOT SUPPORTED

-- NO: arbitrary ORDER BY (must match clustering order)
SELECT * FROM sensor_readings WHERE sensor_id = ?
ORDER BY temperature DESC;  -- ERROR (not a clustering column)
```

These limitations exist because of distribution — Cassandra can't efficiently scan all nodes for a query that doesn't specify the partition.

---

## 5. Consistency Levels

Because data is on multiple nodes, reads and writes can be configured for different consistency:

| Level | Quorum? | Meaning |
|-------|---------|---------|
| `ONE` | No | Fastest: contact 1 replica. May read stale data. |
| `QUORUM` | Yes | Contact majority (RF/2 + 1). Strong consistency if also written with QUORUM. |
| `ALL` | Yes | Contact all replicas. Strongest consistency, but fails if any node is down. |
| `LOCAL_QUORUM` | Yes | Quorum within local datacenter only. |

```sql
-- Set consistency level for a query
CONSISTENCY QUORUM;
SELECT * FROM sensor_readings WHERE sensor_id = ?;
```

**The rule for strong consistency:** Write with QUORUM + Read with QUORUM = you always read your own writes.

---

## 6. Cassandra Anti-Patterns

### Too Many Tombstones

Cassandra's DELETE creates a **tombstone** — a marker saying "this row was deleted." Tombstones accumulate until compaction removes them. Too many tombstones slow down reads.

**Avoid:** deleting lots of individual rows. Instead, use TTL (rows expire naturally) or time-bucketed tables (drop old bucket tables entirely).

### Unbounded Partitions

A partition can't be split across nodes. If you put too much data in one partition, that node becomes a hotspot.

**Bad:**
```sql
-- All events in one partition — gets huge!
PRIMARY KEY (type, event_time)  -- "click" partition has billions of rows
```

**Good:**
```sql
-- Bucket by day: partition stays bounded
PRIMARY KEY ((type, date_bucket), event_time)
-- "click:2024-01-15" partition has one day of clicks
```

### Secondary Indexes

Cassandra's secondary indexes are local (each node indexes only its data) and are slow at scale. Avoid for high-cardinality columns. Instead, create a separate table for each query pattern.

---

## 7. Building with Cassandra in Go

```bash
docker run -d --name cassandra -p 9042:9042 cassandra:4.1
go get github.com/gocql/gocql
```

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/gocql/gocql"
)

func main() {
    cluster := gocql.NewCluster("localhost")
    cluster.Keyspace = "iot"
    cluster.Consistency = gocql.Quorum
    cluster.ProtoVersion = 4
    cluster.ConnectTimeout = 10 * time.Second
    cluster.Timeout = 5 * time.Second

    session, err := cluster.CreateSession()
    if err != nil {
        log.Fatal("connect:", err)
    }
    defer session.Close()

    // Create keyspace and table
    session.Query(`CREATE KEYSPACE IF NOT EXISTS iot
        WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec()
    session.Query(`CREATE TABLE IF NOT EXISTS iot.sensor_readings (
        sensor_id   UUID,
        recorded_at TIMESTAMP,
        temperature FLOAT,
        PRIMARY KEY (sensor_id, recorded_at)
    ) WITH CLUSTERING ORDER BY (recorded_at DESC)`).Exec()

    // Insert a reading
    sensorID := gocql.TimeUUID()
    err = session.Query(
        "INSERT INTO iot.sensor_readings (sensor_id, recorded_at, temperature) VALUES (?, ?, ?)",
        sensorID, time.Now(), 23.5,
    ).Exec()
    if err != nil {
        log.Fatal("insert:", err)
    }

    // Query readings
    var temp float32
    var ts time.Time
    var id gocql.UUID

    iter := session.Query(
        "SELECT sensor_id, recorded_at, temperature FROM iot.sensor_readings WHERE sensor_id = ? LIMIT 10",
        sensorID,
    ).Iter()

    for iter.Scan(&id, &ts, &temp) {
        fmt.Printf("Sensor %s at %s: %.1f°C\n", id, ts.Format(time.RFC3339), temp)
    }
    if err := iter.Close(); err != nil {
        log.Fatal("iter:", err)
    }

    // Batch insert (for performance)
    batch := session.NewBatch(gocql.UnloggedBatch)
    for i := 0; i < 100; i++ {
        batch.Query(
            "INSERT INTO iot.sensor_readings (sensor_id, recorded_at, temperature) VALUES (?, ?, ?)",
            sensorID, time.Now().Add(time.Duration(i)*time.Second), float32(20+i%10),
        )
    }
    if err := session.ExecuteBatch(batch); err != nil {
        log.Fatal("batch:", err)
    }
    fmt.Println("Batch inserted 100 readings")
}
```

---

## Summary

- Cassandra is a peer-to-peer database — no primary, no single point of failure. Every node accepts writes.
- Data is distributed by **partition key** (which node) and sorted by **clustering columns** (order within node).
- Design tables for specific query patterns — you cannot do ad-hoc queries or JOINs.
- Consistency levels trade off speed vs accuracy: ONE (fast, stale), QUORUM (balanced), ALL (safe, slow).
- Avoid: unbounded partitions, too many tombstones, secondary indexes on high-cardinality columns.

### Exercises

**Easy:** Create a `user_activity_log` table in Cassandra. Insert 10 events for user A and 5 for user B. Query only user A's events.

**Medium:** Design a time-series schema for tracking website page views by URL and day. Write a Go function that inserts a view and another that returns the view count for a given URL on a given day.

**Hard:** Implement the "time bucket" pattern: a `click_events` table partitioned by `(event_type, date_bucket)` where `date_bucket = toDate(event_time)`. Write a Go function that queries clicks across multiple days by iterating over date buckets.
