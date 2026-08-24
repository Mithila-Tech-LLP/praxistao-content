---
title: Consistent Hashing, Explained
category: Software & Programming
tags: [System Design, Distributed Systems]
duration: 8 min read
relatedCourses: [senior-engineer-interview, databases-and-async-go]
relatedProjects: [key-value-store]
relatedTopics: [cap-theorem-in-plain-terms, load-balancing-strategies]
---

## TL;DR

- Plain hashing (`hash(key) % N`) breaks badly when `N` (the number of servers) changes — nearly every key remaps to a different server, causing a massive cache/data reshuffle.
- Consistent hashing places both servers and keys on the same conceptual ring (a hash space that wraps around); a key belongs to the next server clockwise from it.
- Adding or removing one server only remaps the keys between it and its neighbor — roughly `1/N` of all keys, not nearly all of them.
- **Virtual nodes** (multiple ring positions per physical server) fix consistent hashing's own weak spot: uneven key distribution when there are few servers.

## The Problem With `hash(key) % N`

Say you have 4 cache servers and route each key with `hash(key) % 4`. This works fine — until you add a 5th server. Now it's `hash(key) % 5`, and for almost every key, the result changes:

```
key "user:42" -> hash = 1042
  % 4 = 2  (server 2)
  % 5 = 2  (server 2, coincidentally same — rare)

key "user:99" -> hash = 1099
  % 4 = 3  (server 3)
  % 5 = 4  (server 4 — different!)
```

In general, changing `N` changes the remainder for the vast majority of keys, not just the ones that "should" move to the new server. For a cache, this means a near-total cache miss storm right after any scaling event — every client suddenly asks the wrong server, which doesn't have the data, so it all gets re-fetched from the origin at once. For a sharded database, it means redistributing almost all the data, not just a proportional slice.

## The Ring

Consistent hashing fixes this by hashing both the **servers** and the **keys** into the same numeric space, and imagining that space as a circle (the hash function's output range wraps from its maximum back to zero).

```
                    0/2^32
                     |
          Server C --+-- Server A
             \        |        /
              \       |       /
               \      |      /
                key "x" (hashes here)
               /      |      \
              /       |       \
             /        |        \
          Server B  --+--
```

To find which server owns a key: hash the key to get a point on the ring, then walk clockwise until you hit the first server. That server owns the key.

```
function getServer(key, servers):
    keyHash = hash(key)
    // servers sorted by their own ring position
    for server in servers sorted by hash(server):
        if hash(server) >= keyHash:
            return server
    return servers[0] // wrapped around past the highest server hash
```

## Why This Fixes the Remapping Problem

When a server is added, it takes a new position on the ring. Only the keys that fall between the new server and the *previous* server going clockwise (the ones that used to belong to that previous server) get remapped to the new one. Every other key's "next server clockwise" is completely unaffected.

```
Before adding Server D:
  ring positions: A(50), B(150), C(280)
  key at 200 -> belongs to C (next clockwise from 200)

After adding Server D at position 220:
  key at 200 -> now belongs to D (D is now the next clockwise server)
  key at 100 -> still belongs to B (unaffected — D is nowhere near it)
```

With `N` servers, adding or removing one server remaps roughly `1/N` of all keys on average — a proportional, bounded amount of churn instead of "almost everything."

## The Uneven Distribution Problem (and Virtual Nodes)

Plain consistent hashing has its own flaw: with only a few servers, their ring positions might land unevenly by chance, giving one server a much larger arc of the ring (and therefore a much larger share of keys) than another.

```
Server A at position 10
Server B at position 20
Server C at position 500  (huge arc from 20 to 500 belongs to... wait, C owns 20->500)
```

The fix used in virtually every real consistent-hashing implementation (this is how Amazon's Dynamo, Cassandra, and most distributed caches actually do it) is **virtual nodes**: instead of one ring position per physical server, hash each server under many different labels (`serverA-0`, `serverA-1`, ... `serverA-150`) so each physical server owns many small, scattered arcs instead of one large contiguous one. With enough virtual nodes per server (typically 100-200), the law of large numbers evens out the distribution close to `1/N` per server, regardless of where any individual server's hash happens to land.

```
function getServer(key, virtualNodes): // virtualNodes: sorted list of (ringPos, physicalServer)
    keyHash = hash(key)
    for (pos, server) in virtualNodes:
        if pos >= keyHash:
            return server
    return virtualNodes[0].server
```

Virtual nodes also make server removal smoother in practice — since a removed server's many small arcs are scattered across the ring, its load gets redistributed across many other servers roughly evenly, instead of dumping its entire share onto whichever single server happens to be next.

## Where This Actually Shows Up

- **Distributed caches** (Memcached client libraries, Redis Cluster's hash-slot approach is a close relative) — so scaling the cache fleet up or down doesn't invalidate almost everything at once.
- **Sharded databases and key-value stores** — deciding which shard/node owns a given key, with the same "add a node without reshuffling everything" property.
- **CDN request routing** — mapping a resource key to one of many edge servers/origins.
- **Load balancers with session affinity** — routing a given client consistently to the same backend, in a way that survives backends being added/removed.

## Common Pitfalls

- **Using a weak or non-uniform hash function** — if the hash function clusters outputs instead of distributing them uniformly across the ring, you get the same uneven-distribution problem virtual nodes are meant to fix, just from a different cause. A good general-purpose hash (e.g., a solid 32/64-bit hash, not something ad hoc) matters.
- **Too few virtual nodes** — with only a handful of virtual nodes per server, you still get meaningfully uneven distribution; 100+ per server is the typical production range.
- **Forgetting replication** — consistent hashing alone tells you the *primary* owner of a key; real systems typically also store the key on the next R-1 servers clockwise from the primary, for redundancy. That replication logic is layered on top of the ring, not part of the basic hashing scheme itself.
