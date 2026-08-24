---
title: Read Replicas and Replication Lag
category: Software & Programming
tags: [Databases, Scaling]
duration: 7 min read
relatedCourses: [databases-and-async-go, senior-engineer-interview]
relatedProjects: []
relatedTopics: [acid-vs-base-consistency-models, cap-theorem-in-plain-terms]
---

## TL;DR

- A read replica is a copy of the primary database that stays continuously updated and serves read-only queries — used to scale read throughput beyond what one server can handle.
- Writes still go to exactly one primary; replicas apply the same writes asynchronously, shortly afterward. That gap between "written on the primary" and "visible on a replica" is **replication lag**.
- Lag is usually milliseconds, but under load, network issues, or a slow replica, it can grow to seconds or more — and any code that reads its own recent write from a replica needs to plan for the write not being there yet.
- The standard fixes: read-your-writes routing (route a user's own follow-up read to the primary, or to whichever replica already has their write), monotonic read guarantees, or simply tolerating brief staleness where it's genuinely fine.

## Why Read Replicas Exist

A single database server has a ceiling on how many queries per second it can serve, no matter how well-indexed and well-tuned. For most applications, reads vastly outnumber writes — a social feed, a product catalog, a dashboard are read far more often than they're written to. Read replicas exploit this: keep one primary handling all writes, and spin up additional read-only copies that handle the read traffic, spreading load across as many replicas as needed.

```
                 writes
Application ---------------> Primary
     |                          |
     |  reads                  | replication stream (async)
     v                          v
  Replica 1  <-------  Replica 2  <---------- ...
```

Applications route writes to the primary and reads to whichever replica is available (often via round robin or least-connections, same as any load-balanced pool) — see Load Balancing Strategies for the routing side of this.

## How Replication Actually Works

The primary records every write to a durable log (PostgreSQL's write-ahead log, MySQL's binlog). Replicas continuously stream that log and replay the same operations in the same order, keeping their own copy of the data in sync. This is asynchronous by default: the primary acknowledges a write to the client as soon as it's durably committed on the primary itself, **without waiting** for any replica to apply it.

This is precisely why replication lag exists at all — it's not a bug, it's the direct consequence of not making every write wait for every replica before responding to the client (which would trade write latency, and often availability, for consistency — the exact tradeoff CAP theorem describes).

## What Replication Lag Actually Looks Like

```
T+0ms:    write "balance = 500" committed on Primary
T+0ms:    client gets "OK" response
T+0ms:    client immediately reads from a Replica -> might still see "balance = 400" (stale!)
T+15ms:   replication stream delivers the write to Replica
T+15ms:   Replica now shows "balance = 500"
```

Under normal conditions this gap is single-digit milliseconds — usually invisible. It grows under real conditions worth knowing about: a replica falling behind because it's under heavy read load itself, a network blip between primary and replica, or a replica applying a particularly expensive write (a bulk update, a schema migration) that takes it longer to catch up. In degraded conditions, lag measured in seconds — occasionally much more — is a real, observed failure mode, not a hypothetical.

## The "Read Your Own Write" Problem

The most common way replication lag actually causes a visible bug: a user submits a form, gets redirected to a page that reads their just-written data back — from a replica that hasn't caught up yet — and sees their own change appear to have not happened.

```
User updates their profile bio -> write goes to Primary -> "Saved!"
User is redirected to /profile -> reads from a Replica -> shows the OLD bio
```

Three common fixes, in order of how often they're actually used:

1. **Route the follow-up read to the primary**, specifically for the user who just wrote, for some short window (or specifically for read-after-write-sensitive operations). Simple, but adds load back onto the primary for exactly the reads replicas were meant to offload.
2. **Track a "read-your-writes" token** — the write returns a marker (often the replication log position at the time of the write); the subsequent read is routed to a replica that has confirmed it's caught up to at least that position, waiting briefly if needed. More complex, but scales better than blanket-routing to the primary.
3. **Design the UI to not need it** — after a write, use the data the client already has (optimistic UI update) instead of immediately re-reading it from the database at all. Often the simplest real fix, when applicable.

## Monotonic Reads

A related, separate guarantee worth knowing by name: **monotonic read consistency** means that once a client has seen a value, it never sees an *older* value on a subsequent read — even if that subsequent read happens to land on a replica that's further behind than the one it read from before. Without this guarantee, a user could refresh a page and see their comment disappear and reappear, simply because consecutive requests happened to hit replicas at different lag levels. This is usually solved by consistently routing a given client's reads to the *same* replica (session affinity — see Consistent Hashing / Load Balancing Strategies) rather than round-robining across replicas that may be at different points in the replication stream.

## Synchronous Replication: The Alternative, at a Cost

Some systems configure at least one replica as **synchronous** — the primary waits for that replica to confirm it received the write before acknowledging the client. This eliminates replication lag for that specific replica (it's always caught up, by construction) at a direct latency cost on every write (now bounded by the round-trip to the synchronous replica, not just the primary's own disk). This is a deliberate, common choice specifically for durability (surviving the primary's outright failure without losing recently-committed writes), more than for solving the read-lag problem directly — most systems still keep the *majority* of their read replicas asynchronous for read scaling, and use synchronous replication selectively, for a small number of replicas, to protect against data loss on primary failure.

## Common Pitfalls

- **Assuming all replicas have identical lag at all times** — lag varies per replica and per moment; code that assumes uniform freshness across a replica pool will eventually be surprised by a replica that's meaningfully behind the others.
- **Reading your own write from a replica without any read-after-write strategy** — this is the most commonly-hit real bug, and it's directly caused by architecture, not a rare edge case.
- **Failing over to a replica during a primary outage without checking how far behind it was** — promoting a lagging replica to primary means any writes that hadn't yet replicated to it are lost. Failover logic needs to account for this, not just "promote whichever replica is available."
- **Treating replication lag as a fixed, ignorable constant** — it's a variable that degrades under exactly the conditions (high load, network issues) when the system is already under stress; monitoring actual replication lag (not assuming it) matters for exactly this reason.
