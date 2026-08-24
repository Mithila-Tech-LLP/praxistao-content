---
title: Internet Checksum
number: 2
difficulty: easy
duration: 15-20 minutes
concept: RFC 1071 one's complement checksum
---

## What to Build

Implement `InternetChecksum`, the same one's-complement checksum algorithm used by IPv4, TCP, and UDP headers to detect corrupted data.

## Function Signature

```go
func InternetChecksum(data []byte) uint16
```

## Requirements

- Sum all 16-bit big-endian words in `data` into a wider accumulator
- If `data` has an odd length, pad the trailing byte with a zero low byte before summing it as a final word
- Fold any carry out of the low 16 bits back into the sum, repeating until there's no carry left
- Return the one's complement of the final sum

## Key Concept: RFC 1071 Checksum

The Internet checksum is deliberately simple: it's cheap to compute in software and catches the overwhelming majority of real-world bit errors, even though it's not cryptographically strong. It's a one's-complement sum — meaning when addition overflows past 16 bits, the overflow bit isn't discarded, it's added back in ("end-around carry"). That's what "fold the carry" means below.

This exact algorithm appears in the UDP header (Chapter 58) and the TCP header's checksum field (Chapter 65). Once you've implemented it once here, you'll recognize it instantly in both.

## Hints

<details>
<summary>Hint 1: Summing 16-bit words</summary>

Iterate over `data` two bytes at a time, combining each pair into a 16-bit word with a shift and OR. Use a `uint32` accumulator so you have headroom above 16 bits before you fold.

</details>

<details>
<summary>Hint 2: The odd-byte case</summary>

If `len(data)` is odd, there's one byte left over after the main loop. Treat it as the high byte of a word whose low byte is zero: `uint32(data[last]) << 8`.

</details>

<details>
<summary>Hint 3: Folding carries</summary>

Fold repeatedly while `sum` still has bits above position 16: add the overflow (`sum >> 16`) back into the low 16 bits, and loop, since one fold may not be enough. Finally, return `^uint16(sum)` — the bitwise complement of the low 16 bits.

</details>

## How to Verify

```bash
lncli run
```
