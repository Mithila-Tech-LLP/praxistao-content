# Chapter 29: CAP Theorem & Consistency Models

CAP theorem is one of the most tested distributed systems concepts in senior interviews. But it's also one of the most misunderstood. This chapter gives you both the theory and the practical reality.

## Table of Contents

1. [CAP Theorem Explained](#1-cap-theorem-explained)
2. [Consistency Models — A Spectrum](#2-consistency-models--a-spectrum)
3. [PACELC — A Better Model](#3-pacelc--a-better-model)
4. [Consistency in Practice](#4-consistency-in-practice)
5. [Strong Consistency Patterns in Go](#5-strong-consistency-patterns-in-go)
6. [Interview Questions & Model Answers](#6-interview-questions--model-answers)
7. [Summary](#summary)

---

## 1. CAP Theorem Explained

CAP theorem: in a distributed system, you can only guarantee 2 of these 3 properties simultaneously:

- **C — Consistency:** every read receives the most recent write (or an error), never stale data
- **A — Availability:** every request receives a response (not necessarily the most recent data)
- **P — Partition Tolerance:** the system continues to operate even when network partitions occur

**The key insight:** network partitions WILL happen. You can't opt out of P. So the real choice is **C vs A during a partition**:

```
Network partition happens (nodes can't communicate):

Choose Consistency (CP):
  System refuses requests (or returns errors) rather than risk stale data
  Examples: etcd, ZooKeeper, HBase
  Use when: financial transactions, leader election, anything where wrong data > no data

Choose Availability (AP):
  System continues serving requests, but some nodes may return stale data
  Examples: Cassandra, DynamoDB (eventually consistent), CouchDB
  Use when: social media feeds, product catalogs, anything where stale is ok
```

### What CAP DOESN'T Mean

```
Common mistake: "MongoDB is CP, Cassandra is AP"
Reality: these databases offer tunable consistency!

Cassandra:
  CONSISTENCY ONE   → AP (fast, potentially stale)
  CONSISTENCY QUORUM → CP-ish (majority must agree)
  CONSISTENCY ALL   → CP (all nodes must agree, but very slow)

PostgreSQL with async replicas:
  Writes to primary → CP on primary
  Reads from replicas → AP (might see stale data)
```

---

## 2. Consistency Models — A Spectrum

From strongest to weakest:

```
LINEARIZABILITY (strongest):
  All operations appear to happen instantaneously at a single point in time.
  Every read returns the value of the most recent completed write.
  Example: etcd, single-node databases
  
SEQUENTIAL CONSISTENCY:
  All operations appear in some sequential order. All nodes see the SAME order.
  But that order might not match real time (global clock not required).
  
CAUSAL CONSISTENCY:
  If operation A causally precedes operation B, all nodes see A before B.
  "If you see my comment reply, you should see the original post."
  
READ-YOUR-OWN-WRITES:
  After a write, your subsequent reads will reflect that write.
  "After posting, you see your own post."
  Not guaranteed on other users' reads.
  
MONOTONIC READS:
  If you've read a value, subsequent reads return the same or newer value.
  You won't "go back in time" within a session.
  
EVENTUAL CONSISTENCY (weakest):
  If no new writes, all replicas will eventually converge to the same value.
  No guarantee on when, or what you'll read in the meantime.
  Example: DNS, S3, social media likes
```

---

## 3. PACELC — A Better Model

CAP only considers behavior during a partition. But partitions are rare! PACELC extends CAP to describe normal operation too:

```
PACELC: If Partition, then choose Availability or Consistency (CA tradeoff)
         ELse (normal operation), choose Latency or Consistency (LC tradeoff)

Examples:
  DynamoDB:  PA/EL = Available during partition, Low Latency in normal operation
  etcd:      PC/EC = Consistent during partition, Consistent in normal operation (but higher latency)
  PostgreSQL: PC/EC (single primary) but PA/EL with async replicas and reads

The LC tradeoff is what you face every day:
  More consistent → must wait for more replicas to confirm → higher latency
  More available  → return cached/local data → lower latency, potentially stale
```

---

## 4. Consistency in Practice

### Bank Transfer Example

```
Bank transfer: deduct from A, credit to B
Requires LINEARIZABILITY — you cannot accept "eventual consistency" here
If network partition: reject the transfer rather than risk double-spend

Pattern: use a single authoritative node (primary) for writes
Read balance from primary (not replica) before transferring
Hold a database transaction across both operations
```

### Social Media Feed Example

```
User posts a photo. 1 billion users might see it.
Acceptable: the photo takes a few seconds to appear in feeds globally.
Not acceptable: operations are slow because we wait for global consistency.

Pattern: eventual consistency is fine
Write to primary, async replicate to edge nodes/CDN
Users in different regions see the post at slightly different times
That's OK — nobody expects instantaneous global consistency
```

### Shopping Cart Example

```
User adds item to cart. Cart must be accurate at checkout.
For browsing: eventual consistency (stale cart is fine)
At checkout: read from primary (strongly consistent) before charging

Pattern: tiered consistency
Default reads → eventually consistent replica (fast)
Payment critical reads → primary (strongly consistent)
```

---

## 5. Strong Consistency Patterns in Go

```go
// Pattern 1: Always write to primary, read back from primary for critical paths
type OrderService struct {
    primaryDB  *sql.DB
    replicaDB  *sql.DB
}

func (s *OrderService) PlaceOrder(ctx context.Context, order *Order) error {
    tx, _ := s.primaryDB.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // All writes AND the confirmation read go to primary
    _, err := tx.ExecContext(ctx, "INSERT INTO orders ...")
    if err != nil { return err }
    
    _, err = tx.ExecContext(ctx, "UPDATE inventory SET quantity -= 1 WHERE product_id = $1", order.ProductID)
    if err != nil { return err }
    
    return tx.Commit()
}

func (s *OrderService) GetOrderForPayment(ctx context.Context, orderID int64) (*Order, error) {
    // Payment path: read from primary only
    var o Order
    s.primaryDB.QueryRowContext(ctx, "SELECT * FROM orders WHERE id = $1", orderID).Scan(&o)
    return &o, nil
}

func (s *OrderService) GetOrderHistory(ctx context.Context, userID int64) ([]*Order, error) {
    // Non-critical path: read from replica (stale OK)
    rows, _ := s.replicaDB.QueryContext(ctx, "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT 20", userID)
    // ...
}

// Pattern 2: Version/token-based consistency
// Client gets a version token after a write, sends it with subsequent reads
// Server can use this to route to primary or wait for replica to catch up
type VersionToken struct {
    WALPosition string // PostgreSQL LSN (Log Sequence Number)
}

func (s *OrderService) CreateOrderWithToken(ctx context.Context, order *Order) (*VersionToken, error) {
    tx, _ := s.primaryDB.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    _, _ = tx.ExecContext(ctx, "INSERT INTO orders ...")
    tx.Commit()
    
    // Return WAL position so client can verify replica has caught up
    var lsn string
    s.primaryDB.QueryRowContext(ctx, "SELECT pg_current_wal_lsn()::text").Scan(&lsn)
    return &VersionToken{WALPosition: lsn}, nil
}
```

---

## 6. Interview Questions & Model Answers

**Q: Explain CAP theorem. Is it actually useful?**

"CAP theorem says you can only have 2 of: Consistency, Availability, Partition Tolerance. Since partition tolerance isn't optional in distributed systems (networks do fail), the real choice is: during a network partition, do you sacrifice consistency or availability? CP systems reject requests to avoid stale data — good for financial systems. AP systems continue serving, possibly with stale data — fine for social feeds. The practical limitation of CAP is it's binary, describing only extreme failures. PACELC is more useful: it also captures the Latency vs Consistency trade-off in normal operation, which is the trade-off you face every day."

**Q: What is eventual consistency and when is it acceptable?**

"Eventual consistency means: if you stop writing, all replicas will eventually converge to the same value — but there's no guarantee on when, and reads during that window may return stale data. It's acceptable when: (1) staleness has low business impact (social media likes, product view counts), (2) the alternative (global coordination) would be too slow or expensive, (3) the system can handle compensating logic when conflicts are detected. It's NOT acceptable for: financial balances, inventory counts, or anything where seeing old data could cause double-spend or data loss."

---

## Summary

- **CAP:** in a partition, choose consistency (reject requests) or availability (serve stale data). Partition tolerance is not optional.
- **PACELC:** extends CAP to include the Latency vs Consistency trade-off in normal operation.
- **Consistency spectrum:** Linearizable → Sequential → Causal → Read-your-own-writes → Eventual.
- In practice: use strong consistency for financial transactions, write to primary for critical reads, allow eventual consistency for social/analytics workloads.
- Most databases offer tunable consistency — the "CP or AP" label is a default, not a hard constraint.
- Design tier by tier: identify which operations require strong consistency and route only those to the primary.
