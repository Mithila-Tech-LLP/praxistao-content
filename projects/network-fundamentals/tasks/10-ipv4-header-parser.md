---
title: IPv4 Header Parser
number: 10
difficulty: hard
duration: 30-40 minutes
concept: IPv4 header layout, byte offsets
---

## What to Build

Implement `ParseIPv4Header`, which reads the fixed fields of a raw IPv4 header directly out of its byte layout — the same bytes a packet sniffer like Wireshark or tcpdump would show you, decoded by hand.

## Function Signature

```go
func ParseIPv4Header(data []byte) (version, ihl, totalLength, ttl, protocol int, srcIP, dstIP net.IP, err error)
```

## Requirements

- Assume no IP options, so a valid header is exactly 20 bytes minimum — return a non-nil `err` if `len(data) < 20`
- `version` is the high nibble of byte 0
- `ihl` (Internet Header Length, in 32-bit words) is the low nibble of byte 0
- `totalLength` is the big-endian `uint16` at bytes 2-3
- `ttl` is byte 8
- `protocol` is byte 9
- `srcIP` is the 4 bytes at offset 12-15, as a `net.IP`
- `dstIP` is the 4 bytes at offset 16-19, as a `net.IP`

## Key Concept: The IPv4 Header Byte Layout

Every field in an IPv4 header lives at a fixed byte offset — there's no delimiter or parsing ambiguity, just "read these specific bytes." The version and IHL famously share a single byte, one nibble each, because IPv4 was designed when every bit of header overhead mattered. Chapter 36 covers what these addresses mean, Chapter 65 walks through this same fixed-offset parsing pattern for the TCP header, and Chapter 114 (building a packet sniffer) is where you'd use exactly this kind of parser on live captured traffic. Chapter 28's Ethernet frame is a useful contrast: Ethernet framing is comparatively simple, and IPv4's packed, bit-level layout is part of why routing and addressing get their own extended treatment.

## Hints

<details>
<summary>Hint 1: Splitting a byte into two nibbles</summary>

`version := int(data[0] >> 4)` shifts the high nibble down into the low 4 bits. `ihl := int(data[0] & 0x0F)` masks off everything but the low nibble. Both come from the same single byte.

</details>

<details>
<summary>Hint 2: Reading a big-endian uint16 by hand</summary>

A big-endian `uint16` is just `high<<8 | low` — reconstruct `totalLength` from `data[2]` and `data[3]` that way, without needing `encoding/binary`.

</details>

<details>
<summary>Hint 3: Slicing out the IP addresses</summary>

`net.IP(data[12:16])` gives you a 4-byte IP address directly from the header bytes. Be careful that this slice aliases the original `data` slice — if that matters in your implementation, copy the bytes instead of slicing directly.

</details>

## How to Verify

```bash
lncli run
```
