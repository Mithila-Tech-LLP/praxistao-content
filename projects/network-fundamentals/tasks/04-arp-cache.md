---
title: Bounded ARP Cache with LRU Eviction
number: 4
difficulty: easy-medium
duration: 20-25 minutes
concept: bounded cache, LRU eviction
---

## What to Build

Implement `ARPCache`, a bounded IP-to-MAC-address cache that evicts its least-recently-used entry once it's full — mirroring how a real operating system keeps its ARP table small while still holding onto actively-used entries.

## Function Signature

```go
type ARPCache struct {
    capacity int
    // your fields
}

func NewARPCache(capacity int) *ARPCache
func (c *ARPCache) Set(ip, mac string)
func (c *ARPCache) Lookup(ip string) (mac string, ok bool)
```

## Requirements

- A bounded cache mapping IP address to MAC address
- When capacity is exceeded, evict the LEAST RECENTLY USED entry
- Both `Set` on a new key and `Lookup` of an existing key count as "use" and must refresh that entry's recency
- Keep it deterministic — no time/TTL involved
- `Lookup` of a key that was never set returns `ok=false`

## Key Concept: The ARP Cache

Every device that communicates on a LAN keeps an ARP cache mapping IP addresses to the MAC addresses they resolve to, so it doesn't have to broadcast an ARP request for every single packet. That cache can't grow unbounded, so real implementations cap its size and evict old entries — typically favoring entries that are still being actively used. Chapter 53 covers ARP and the ARP cache in depth; this task builds the eviction policy a real cache needs once it hits its size limit.

## Hints

<details>
<summary>Hint 1: A doubly-linked list plus a map</summary>

The classic LRU structure is a `container/list.List` (for recency order) paired with a `map[string]*list.Element` (for O(1) lookup by key). The front of the list is most-recently-used; the back is least-recently-used.

</details>

<details>
<summary>Hint 2: Refreshing on both Set and Lookup</summary>

Don't forget: `Lookup` of an existing key is a "use" too, and must move that entry to the front of the list — not just `Set`. This is the detail that's easy to miss and easy to get wrong.

</details>

<details>
<summary>Hint 3: Eviction only happens on a new key</summary>

If `Set` is called with a key that already exists, that's an update, not an insertion — it shouldn't trigger eviction. Eviction only needs to happen when you're about to insert a genuinely new key and the cache is already at capacity.

</details>

## How to Verify

```bash
lncli run
```
