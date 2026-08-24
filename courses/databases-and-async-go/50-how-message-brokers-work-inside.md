# Chapter 50: How Message Brokers Work Inside

You know WHY async systems exist. Now let's understand HOW message brokers implement them. What happens when you call `producer.Send("order_placed", data)`? How does Kafka store millions of messages per second without losing any?

## Table of Contents

1. The Broker's Job
2. The Commit Log — Kafka's Core Data Structure
3. Topics, Partitions, and Offsets
4. How Producers Write
5. How Consumers Read
6. Replication for Durability
7. Consumer Groups — Work Distribution
8. Exercises

---

## 1. The Broker's Job

A broker is a server that:
1. Accepts messages from producers (writes)
2. Stores them durably
3. Delivers them to consumers (reads)
4. Tracks which messages each consumer has processed

Think of it as a very specialized database for sequential data — optimized for high-throughput writes and reads, not random access or complex queries.

---

## 2. The Commit Log — Kafka's Core Data Structure

Kafka stores messages in an **append-only commit log** (a file that only grows). This is almost identical to the WAL we built in VaultDB!

```
messages.log
┌─────────────────────────────────────────────────────────────────────┐
│ Offset 0: [len][key][value][timestamp][headers]                     │
│ Offset 1: [len][key][value][timestamp][headers]                     │
│ Offset 2: [len][key][value][timestamp][headers]                     │
│ Offset 3: [len][key][value][timestamp][headers]                     │
│ Offset 4: [len][key][value][timestamp][headers]   ← latest          │
└─────────────────────────────────────────────────────────────────────┘
```

**Why append-only is fast:**
- Appending to a file = sequential write. HDDs do this at 100-500 MB/s.
- No disk seeks (jumping to random positions). Sequential is 100-1000x faster than random.
- Modern OS page cache keeps recently written data in RAM, making even reads fast.

**Message retention:** Unlike a queue, the log doesn't delete messages after consumption. Messages expire after a retention period (e.g., 7 days) or when log size exceeds a limit (e.g., 100 GB).

---

## 3. Topics, Partitions, and Offsets

**Topic:** A named category of messages. Like a database table for events.

**Partition:** A topic is split into N partitions, each on a different broker. This enables:
- **Parallel writes:** Each producer can write to a different partition simultaneously.
- **Parallel reads:** Each consumer reads from a different partition simultaneously.
- **Scalability:** Add partitions to increase throughput.

```
Topic: "orders" (3 partitions across 2 brokers)

Broker 1:
  Partition 0: [msg0, msg3, msg6, msg9, ...]
  Partition 1: [msg1, msg4, msg7, msg10, ...]

Broker 2:
  Partition 2: [msg2, msg5, msg8, msg11, ...]
```

**Offset:** The position of a message within its partition. Offsets start at 0 and increase monotonically.

```
Partition 0: offset 0, offset 1, offset 2, ...
Partition 1: offset 0, offset 1, offset 2, ...
```

**Key insight:** Offset is the consumer's bookmark. A consumer reads "partition 0, from offset 42." Next time it reads from 43.

**Partition assignment for a message:**
- No key: round-robin across partitions.
- With key: `hash(key) % num_partitions`. Same key always goes to same partition → ordered delivery per key.

```go
// Same order ID always goes to same partition → order events arrive in order
producer.Send("orders", orderID, data) // key=orderID ensures ordering
```

---

## 4. How Producers Write

**The write path:**

```
1. Producer serializes message
2. Producer sends to broker: {topic: "orders", partition: 2, key: "order-123", value: {...}}
3. Broker appends to partition 2's log file
4. Broker optionally replicates to follower brokers
5. Broker acknowledges: "message at offset 447 committed"
6. Producer proceeds
```

**Acknowledgement levels (acks):**

- `acks=0`: Fire and forget. Fastest, but data may be lost if broker crashes.
- `acks=1`: Wait for the leader broker to write. Lost if leader crashes before replication.
- `acks=all` (or `-1`): Wait for all replicas to acknowledge. Safest, slightly slower.

For financial data, always use `acks=all`.

**Batching for throughput:**

Producers buffer messages in memory and send them in batches. Instead of one network round-trip per message:

```
Producer buffer: [msg1, msg2, msg3, ...msg100]
Single network call: send 100 messages to broker
Broker: appends all 100 to log, sends one ACK
```

This gives 10-100x higher throughput than sending messages one by one.

---

## 5. How Consumers Read

**The read path:**

```
1. Consumer asks: "Give me up to 500 messages from topic 'orders', partition 2, starting at offset 100"
2. Broker reads from the log file starting at offset 100
3. Broker returns messages 100-599 (if available)
4. Consumer processes each message
5. Consumer commits: "I've processed up to offset 599"
```

**Commit strategies:**

```go
// Auto-commit (risky): SDK commits every 5 seconds
// If the process crashes at second 4, the last 4 seconds of work is lost

// Manual commit (safe): commit only after processing
messages := consumer.Poll(500 * time.Millisecond)
for _, msg := range messages {
    processMessage(msg)
}
consumer.CommitSync() // commit after processing
```

**What "commit" means in Kafka:** The consumer group's offset for this partition is stored back in Kafka itself (in the `__consumer_offsets` internal topic). This means if the consumer crashes and restarts, it can resume exactly where it left off.

---

## 6. Replication for Durability

Data is precious. A single broker can fail (disk failure, network issue, power outage). Replication copies each partition to multiple brokers:

```
Partition 0:
  Leader: Broker 1 (accepts reads and writes)
  Follower: Broker 2 (copies from leader)
  Follower: Broker 3 (copies from leader)
```

**ISR (In-Sync Replicas):** Set of replicas that are fully caught up with the leader. A write is committed only when all ISR replicas have received it (with `acks=all`).

**If the leader fails:** Kafka automatically elects a new leader from the ISR. This takes < 30 seconds. Consumers and producers automatically reconnect to the new leader.

**Replication factor:** `RF=3` means data survives 2 simultaneous broker failures. Production minimum: RF=3.

---

## 7. Consumer Groups — Work Distribution

Consumer groups solve: "I have 10 million messages/sec coming in. One consumer can process 1 million/sec. How do I process all of them?"

```
Topic: "orders" (10 partitions)

Consumer Group: "order-processors" (5 consumers)

Consumer 1 → Partitions 0, 1
Consumer 2 → Partitions 2, 3
Consumer 3 → Partitions 4, 5
Consumer 4 → Partitions 6, 7
Consumer 5 → Partitions 8, 9
```

**Rules:**
- Each partition is assigned to exactly one consumer in the group.
- Multiple consumer groups can all read from the same topic independently — each group gets all messages.
- Consumer count > partition count = some consumers idle (no benefit to adding more).
- Consumer count < partition count = some consumers handle multiple partitions.
- Rebalancing: if a consumer joins or leaves, Kafka redistributes partitions automatically.

```
                       Topic: "orders"
                      ┌────────────────────┐
                      │  Partition 0       │◄──── Consumer Group A: Consumer 1
                      │  Partition 1       │◄──── Consumer Group A: Consumer 2
                      │  Partition 2       │◄──── Consumer Group A: Consumer 3
                      │  Partition 0       │◄──── Consumer Group B: Consumer 1
                      │  Partition 1       │◄──── Consumer Group B: Consumer 1
                      └────────────────────┘

Group A = order processors (each consumer processes a subset)
Group B = analytics service (reads all messages independently)
```

---

## Summary

- Message brokers use an append-only commit log for storage — exactly like a database WAL.
- Topics are divided into partitions. Each partition is an independent log with sequential offsets.
- Producers batch messages for throughput. `acks=all` for durability.
- Consumers read from a partition at their current offset. Manual commit after processing = no data loss.
- Replication (RF=3) ensures data survives broker failures.
- Consumer groups distribute work: each partition to one consumer, but multiple groups get all messages.

### Exercises

**Easy:** Draw the read/write path for this scenario: 3 producers writing to a topic with 6 partitions, replication factor 3, 2 consumer groups each with 3 consumers. How many partitions does each consumer handle?

**Medium:** Research Kafka's log compaction. How does it differ from log retention? When would you use compaction instead of retention? What problem does it solve?

**Hard:** Explain the "consumer group rebalance" process. What happens when a consumer crashes? When a new consumer joins? What is a "rebalance storm" and how do "static group membership" (KIP-345) and "incremental cooperative rebalancing" mitigate it?
