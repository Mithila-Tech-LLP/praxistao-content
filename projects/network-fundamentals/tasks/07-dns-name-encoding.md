---
title: DNS Name Encoding
number: 7
difficulty: hard
duration: 30-40 minutes
concept: DNS wire format, length-prefixed labels
---

## What to Build

Implement `EncodeDNSName` and `DecodeDNSName`, which convert domain names between their human-readable dotted form (`"example.com"`) and the length-prefixed wire format DNS actually sends over the network.

## Function Signature

```go
func EncodeDNSName(name string) []byte
func DecodeDNSName(data []byte, offset int) (name string, bytesRead int, err error)
```

## Requirements

- `EncodeDNSName` splits the domain on `"."`, and for each label writes a single length byte followed by that label's raw bytes, terminated by a single zero byte
- `DecodeDNSName` does the reverse starting at `offset` in `data`: read a length byte, read that many raw bytes as a label, repeat until a zero-length byte is read, joining the labels with `"."`
- `DecodeDNSName` returns the reconstructed name, the total number of bytes consumed (including the final zero byte), and a non-nil error if `data` runs out before a terminating zero byte is found
- DNS compression pointers (the `0xC0` prefix used to reference an earlier name in the packet) are explicitly out of scope — this task only handles plain, uncompressed names

## Key Concept: DNS Wire Format

DNS doesn't send domain names as plain strings with dots — it uses a length-prefixed label format, where each label (the text between dots) is preceded by a single byte giving its length, and the whole name ends with a zero-length "label" acting as a terminator. Chapters 66-69 cover why DNS exists, its hierarchy, resolvers and caching, and record types — this task is the low-level encoding those higher-level concepts all ride on top of. It's also the exact parsing logic you'd need for Chapter 111's DNS resolver, just isolated here as its own standalone piece.

## Hints

<details>
<summary>Hint 1: Encoding is a straightforward loop</summary>

`strings.Split(name, ".")` gives you the labels. For each one, append `byte(len(label))` followed by the label's bytes (`[]byte(label)`). After the loop, append a single trailing `0` byte.

</details>

<details>
<summary>Hint 2: Decoding tracks a moving position</summary>

Keep a `pos` variable starting at `offset`. Read `data[pos]` as the length, advance `pos` by one, then if the length is nonzero, read that many bytes as a label and advance `pos` by the label's length. Stop when you read a length of zero. `bytesRead` is `pos - offset` after consuming the terminating zero.

</details>

<details>
<summary>Hint 3: Guard against running off the end</summary>

Before reading a length byte or a label's bytes, check that there's actually enough data left (`pos < len(data)`, and `pos+length <= len(data)`). If not, return an error instead of panicking with an out-of-bounds slice access.

</details>

## How to Verify

```bash
lncli run
```
