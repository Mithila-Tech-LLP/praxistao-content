---
title: Longest Prefix Match Router
number: 3
difficulty: medium
duration: 25-35 minutes
concept: routing tables, CIDR containment, prefix length
---

## What to Build

Implement a minimal router lookup table: `Router` stores a list of CIDR routes, each pointing at a next hop, and `Lookup` finds the correct next hop for a given IP address using longest-prefix-match — the same rule every real router uses.

## Function Signature

```go
type Route struct {
    CIDR    string
    NextHop string
}

type Router struct {
    routes []Route
}

func NewRouter() *Router
func (r *Router) AddRoute(cidr, nextHop string)
func (r *Router) Lookup(ip string) (nextHop string, found bool)
```

## Requirements

- `AddRoute` appends a route to the table
- `Lookup` must find the route with the LONGEST matching prefix — the most specific `/N` wins, not the first match added
- Use `net.ParseCIDR` to parse each route and `ipnet.Contains(net.ParseIP(ip))` to test containment
- Return `found=false` if no route matches (there may be no default route)

## Key Concept: Longest Prefix Match

A routing table often has multiple entries that could match a destination IP — a broad default route (`0.0.0.0/0`) alongside increasingly specific routes for particular subnets. Real routers always prefer the most specific match, because it represents more precise knowledge about how to reach that destination. This is exactly the forwarding decision covered in Chapter 45 (forwarding, next hop, and longest prefix match) — you're implementing the core of what a router's forwarding table does on every packet.

## Hints

<details>
<summary>Hint 1: Checking containment</summary>

For each route, parse its CIDR with `net.ParseCIDR` to get an `*net.IPNet`. `ipnet.Contains(parsedIP)` tells you whether the IP falls inside that network — this correctly handles the mask arithmetic for you.

</details>

<details>
<summary>Hint 2: Comparing specificity</summary>

`ones, _ := ipnet.Mask.Size()` gives you the prefix length of a matching route. Track the best (highest) `ones` seen so far as you iterate over every route, and keep its next hop.

</details>

<details>
<summary>Hint 3: No match found</summary>

Don't assume there's always a default route. If you iterate every route and none contain the IP, `found` should be `false` and `nextHop` should be the zero value.

</details>

## How to Verify

```bash
lncli run
```
