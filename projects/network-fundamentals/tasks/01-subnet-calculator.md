---
title: Subnet Calculator
number: 1
difficulty: easy
duration: 20-30 minutes
concept: CIDR, net.ParseCIDR, IPv4 addressing
---

## What to Build

Implement `SubnetInfo`, which takes a CIDR string like `"192.168.1.0/24"` and computes everything you'd normally reach for a subnet calculator website to find: the network address, the broadcast address, the first and last usable host addresses, and the number of usable hosts.

## Function Signature

```go
func SubnetInfo(cidr string) (network, broadcast, firstHost, lastHost net.IP, hostCount int, err error)
```

## Requirements

- Parse the CIDR string with `net.ParseCIDR`
- Return the network address, broadcast address, first usable host, and last usable host as 4-byte IPv4 `net.IP` values
- `hostCount` is the number of usable hosts: `2^hostBits - 2` (subtracting the network and broadcast addresses)
- Return a non-nil `err` for an invalid CIDR string (e.g. `"not-a-cidr"`)
- The input IP does not need to be network-aligned — `"192.168.1.100/24"` must still resolve to the containing network `192.168.1.0/24`

## Key Concept: CIDR and Subnet Math

An IPv4 address is 32 bits split into a network portion and a host portion. CIDR notation (`/24`) tells you where that split happens. Once you know the mask, the network address is `IP AND mask`, and the broadcast address is `network OR (NOT mask)` — flipping every host bit to 1. The first usable host is one past the network address; the last usable host is one before the broadcast address.

This is exactly the arithmetic covered in Computer Networks Chapters 36-39: IPv4 addressing, subnetting from first principles, and CIDR. If those chapters felt abstract, this task is where the bit math becomes muscle memory.

## Hints

<details>
<summary>Hint 1: Working with net.ParseCIDR</summary>

`net.ParseCIDR(cidr)` returns an IP and an `*net.IPNet`. The `ipnet.IP` is already the network address (Go zeroes out the host bits for you), and `ipnet.Mask` is the subnet mask. Call `.To4()` to make sure you're working with 4-byte representations, not 16-byte IPv6-mapped ones.

</details>

<details>
<summary>Hint 2: Treat addresses as 32-bit integers</summary>

It's much easier to do this arithmetic on `uint32` than on byte slices. `binary.BigEndian` from the `encoding/binary` package converts a 4-byte `net.IP` to and from a `uint32` — do your math on the integer, then convert back.

</details>

<details>
<summary>Hint 3: Host count formula</summary>

`ipnet.Mask.Size()` gives you the prefix length and total bits (32 for IPv4). Host bits are the difference between them, and usable hosts are `2^hostBits - 2` — except watch out for degenerate cases like `/31` or `/32` where that formula alone isn't meaningful. The test vectors here stick to normal subnets, so you don't need to special-case those.

</details>

## How to Verify

```bash
lncli run
```
