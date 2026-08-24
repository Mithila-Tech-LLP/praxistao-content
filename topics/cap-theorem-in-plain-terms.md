---
title: CAP Theorem, in Plain Terms
category: Software & Programming
tags: [System Design, Distributed Systems]
duration: 7 min read
relatedCourses: [senior-engineer-interview, databases-and-async-go]
relatedProjects: []
relatedTopics: [read-replicas-and-replication-lag, consistent-hashing-explained]
---

## TL;DR

- CAP says: when a network partition happens between nodes, a distributed system must choose between staying **Consistent** (every read sees the latest write) or staying **Available** (every request gets a response) — it cannot guarantee both during that partition.
- It does *not* say "pick 2 of 3 always" — Partition tolerance isn't really optional for any system spread across more than one machine, because networks fail regardless of what you'd prefer. The real choice is CP vs AP, and only *during* an actual partition.
- Outside of a partition, most real systems are both consistent and available — CAP is specifically about the failure case, not steady-state operation.
- In practice, most systems land somewhere on a consistency *spectrum* (strong, eventual, causal, read-your-writes...) rather than a binary CP/AP choice — CAP is a useful mental model, not a literal specification.

## The Three Letters

- **Consistency (C)**: every node that responds to a read returns the most recent write, or an error — never stale data. (Note: this is a *different* "consistency" than the C in ACID, which is about database invariants like foreign keys — a common point of confusion.)
- **Availability (A)**: every request to a non-failing node receives a (non-error) response — it just might not be the most recent value.
- **Partition tolerance (P)**: the system keeps operating even when network messages between nodes are dropped or delayed.

## The Actual Claim

The theorem, as originally proven by Eric Brewer and formalized by Gilbert and Lynch, is specifically about what happens **during a network partition** — some nodes can't talk to others. In that moment, if a node receives a write, it has exactly two options:

1. **Refuse to respond** (or respond with an error) until it can confirm the value is consistent with the rest of the cluster — this preserves Consistency, sacrifices Availability.
2. **Respond anyway**, using whatever value it locally has — this preserves Availability, sacrifices Consistency (some other node might have a newer value it doesn't know about yet).

There is no third option that preserves both, *during the partition*. That's the entire theorem. It says nothing about what happens when there's no partition — most systems are fully consistent and available under normal operation, and only have to make this tradeoff when things are already going wrong.

## Why "Partition Tolerance" Isn't Really a Choice

A common misreading of CAP is "pick any 2 of C, A, P." In practice, P isn't optional for any system that runs on more than one machine connected by a real network — because that network *will* eventually drop or delay messages, whether or not your design "chooses" to tolerate it. A system that isn't partition-tolerant just means: when a partition happens, it breaks in an undefined way, rather than degrading in a chosen, deliberate direction.

So the real-world decision is really **CP vs AP**, and only meaningful *during* an actual partition:

```
Partition happens between Node A and Node B.
A client writes to Node A. Can Node B's replica of that data
be read from right now?

  CP system: Node B refuses the read (or blocks) until it can
             confirm it has the latest value. Consistent, not Available.

  AP system: Node B serves its last known value immediately.
             Available, but possibly stale (not Consistent).
```

## Where Real Systems Land

- **CP-leaning**: systems like ZooKeeper, etcd, and traditional single-primary relational databases (during a failover) prioritize correctness over uptime — they'd rather return an error than a wrong answer, because the things they coordinate (leader election, configuration) break badly if two nodes disagree about the truth.
- **AP-leaning**: systems like DynamoDB, Cassandra, and most CDNs prioritize staying responsive — a stale product listing or a slightly-out-of-date follower count is a far smaller problem than the whole site going down.
- **Most real systems aren't purely one or the other** — a system might be CP for writes (a single leader accepts them) and AP for reads (any replica can serve a possibly-stale read), which is exactly how most "eventually consistent" read replica setups work in practice.

## Why "Eventual Consistency" Is the More Useful Everyday Concept

CAP is a binary, worst-case framing (what happens during a partition). Most engineering conversations about real systems are better served by talking about the actual **consistency model** a system offers day-to-day:

- **Strong consistency**: every read sees the latest write, always (what a CP choice buys you).
- **Eventual consistency**: reads may return stale data temporarily, but all replicas converge to the same value once writes stop propagating.
- **Read-your-writes consistency**: a specific, weaker guarantee — a client is guaranteed to see its *own* writes, even if it might see other clients' writes with some delay.
- **Causal consistency**: writes that are causally related (B happened after seeing A) are seen by everyone in that same order; unrelated writes can be seen in different orders by different readers.

Naming which of these a system actually provides is almost always more useful in a design discussion than saying "it's AP."

## Common Pitfalls

- **Treating CAP as "pick 2 of 3" as a permanent architectural decision** — the tradeoff only actually applies during a partition; conflating it with everyday behavior overstates what the theorem claims.
- **Confusing CAP's "Consistency" with ACID's "Consistency"** — ACID consistency means the database enforces its own declared invariants (constraints, foreign keys); CAP consistency means all nodes agree on the current value of a piece of data. They're unrelated ideas that happen to share a word.
- **Assuming a single-node database is somehow "outside" CAP** — CAP is about distributed systems specifically; a single-node database with no replicas doesn't have a partition to reason about at all, so the theorem doesn't apply until you add replication.
