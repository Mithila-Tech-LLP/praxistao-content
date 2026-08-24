---
title: ACID vs BASE Consistency Models
category: Software & Programming
tags: [Databases, Distributed Systems]
duration: 7 min read
relatedCourses: [databases-and-async-go, senior-engineer-interview]
relatedProjects: []
relatedTopics: [cap-theorem-in-plain-terms, optimistic-vs-pessimistic-locking]
---

## TL;DR

- **ACID** (Atomicity, Consistency, Isolation, Durability) describes the strong transactional guarantees traditional relational databases give you — the whole transaction happens or none of it does, and once committed, it's safely durable.
- **BASE** (Basically Available, Soft state, Eventual consistency) describes the looser guarantees many distributed/NoSQL systems make instead, deliberately trading strict consistency for availability and horizontal scale.
- These aren't "old way vs new way" — they're different tradeoffs for different problems. A bank ledger needs ACID; a social media like-counter is fine with BASE.
- Many modern systems mix both: ACID within a single node or a bounded transaction scope, BASE across replicas or services.

## ACID, One Letter at a Time

**Atomicity** — a transaction's operations happen as an indivisible unit: all of them succeed, or none do. Transferring money between two accounts (debit one, credit the other) must never leave the system in a state where the debit happened but the credit didn't, even if the process crashes mid-transaction.

```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT; -- both updates take effect together, or (on any failure before COMMIT) neither does
```

**Consistency** — a transaction can only move the database from one valid state to another, according to its own declared rules (foreign keys, unique constraints, check constraints). A transaction that would violate a `NOT NULL` or a foreign key constraint is rejected outright, not partially applied.

**Isolation** — concurrent transactions don't see each other's uncommitted, in-progress changes. Exactly how strictly this is enforced is itself tunable (see isolation levels below) — full isolation is the strongest, most expensive option, not the only one.

**Durability** — once a transaction commits, it survives a crash immediately afterward. This is typically implemented via a write-ahead log: the change is durably recorded on disk before the commit is acknowledged to the client, so even a power loss the instant after doesn't lose it.

## Isolation Levels: How Much Isolation You're Actually Paying For

"Isolation" isn't binary — SQL defines four standard levels, each allowing (or preventing) specific concurrency anomalies, at increasing cost:

| Level | Prevents | Allows |
|---|---|---|
| Read Uncommitted | nothing | dirty reads (seeing another transaction's uncommitted changes) |
| Read Committed | dirty reads | non-repeatable reads (a value changes between two reads in the same transaction) |
| Repeatable Read | dirty + non-repeatable reads | phantom reads (a query's result set changes between two runs in the same transaction) |
| Serializable | all of the above | — behaves as if transactions ran one at a time |

Most production systems default to Read Committed (PostgreSQL, Oracle) or Repeatable Read (MySQL/InnoDB) rather than Serializable, precisely because full serializability is the most expensive to actually enforce (typically via more aggressive locking or transaction aborts-and-retries), and most applications don't need the strongest guarantee for most of their transactions.

## BASE: The Deliberate Alternative

BASE isn't a formal acronym in the same rigorous sense as ACID — it was coined specifically as a contrast, to name what large-scale distributed systems (the ones that inspired the CAP theorem framing) were actually doing instead:

- **Basically Available**: the system responds to requests most of the time, even during a partial failure — prioritizing uptime over strict correctness (this is the "A" side of the CAP tradeoff).
- **Soft state**: the system's state may change over time even without new input, as replicas converge toward consistency in the background.
- **Eventual consistency**: if no new writes happen, all replicas will *eventually* converge to the same value — but there's no guarantee about exactly when, and a read in the meantime might return a stale value.

A concrete example: a "like count" on a post, replicated across multiple regions. A like registered in one region might take a few hundred milliseconds (or, under network issues, longer) to propagate to other regions' replicas. During that window, different users might see slightly different counts. This is an explicit, accepted tradeoff — the count converging *eventually* is fine, because a like count being off by a handful for a moment has essentially no real cost, and demanding strict consistency here would mean either much higher latency (waiting for all replicas to agree before responding) or reduced availability (rejecting writes during any partition).

## Why the Choice Actually Matters

The two models aren't interchangeable — using the wrong one for a given piece of data has real consequences:

- **A bank balance needs ACID.** Eventual consistency on account balances means two concurrent withdrawals could both succeed against a stale balance, overdrawing the account — a correctness bug with direct financial consequences, not a cosmetic glitch.
- **A social media like-counter, a view count, a "user is typing..." indicator are fine with BASE.** Demanding ACID-level consistency for these — say, requiring every replica worldwide to agree before showing a count — would add latency and reduce availability for a guarantee the use case doesn't actually need.

## Mixing Both in the Same System

Most real, large systems aren't purely one or the other. A common shape: the system of record for critical data (orders, payments, inventory counts that must never go negative) uses a traditional ACID-compliant relational database, often on a single primary with synchronous replication for durability. Meanwhile, derived or replicated views of that data used for reads at scale (a search index, a cache, read replicas serving a recommendation feed) are BASE — eventually consistent, refreshed asynchronously, prioritizing availability and read throughput over strict up-to-the-millisecond accuracy.

## Common Pitfalls

- **Assuming "NoSQL" automatically means BASE and "SQL" automatically means ACID** — this correlation is common but not a rule; some NoSQL databases offer strong consistency options, and some relational setups (read replicas with async replication) introduce eventual consistency despite being "SQL."
- **Applying eventual consistency to data where staleness has a real correctness cost** — inventory counts, account balances, anything where "briefly wrong" translates directly into a business or safety problem, need the stronger guarantee, not the more scalable one.
- **Assuming Serializable isolation is "free correctness" with no downside** — it's the right choice for genuinely correctness-critical transactions, but applying it universally adds real performance cost (often via transaction retries on conflict) that most queries in a system don't actually need.
- **Confusing "eventual consistency" with "no consistency"** — eventual consistency is still a real guarantee (convergence, given no further writes); it's a weaker promise than strong consistency, not an absence of one.
