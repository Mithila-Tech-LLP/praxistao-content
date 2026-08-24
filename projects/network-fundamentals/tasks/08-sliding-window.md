---
title: TCP Sliding Window Simulation
number: 8
difficulty: medium
duration: 25-35 minutes
concept: TCP flow control, cumulative ACK
---

## What to Build

Implement `SlidingWindow`, a simplified simulation of TCP's flow control mechanism: a sender can have multiple unacknowledged segments in flight, up to the window size, and a single cumulative ACK can acknowledge several segments at once.

## Function Signature

```go
type SlidingWindow struct {
    base int // oldest unacknowledged sequence number
    next int // next sequence number to send
    size int // window size
}

func NewSlidingWindow(size int) *SlidingWindow
func (w *SlidingWindow) Send() (seq int, ok bool) // returns the next seq number sent, or ok=false if window is full
func (w *SlidingWindow) Ack(seq int)              // cumulative ack: acknowledges everything up to and including seq
```

## Requirements

- `Send()` may only succeed while `next < base+size`; on success it returns the current `next` and then increments it
- `Send()` returns `ok=false` when the window is full — no more room until an ACK arrives
- `Ack(seq)` is CUMULATIVE, like real TCP: it advances `base` to `seq+1`
- `Ack(seq)` must ignore stale or out-of-order acks that would move `base` backward — only apply it if it's forward progress

## Key Concept: Sliding Window Flow Control

TCP doesn't send one segment and wait for its ACK before sending the next — that would waste most of the available bandwidth on every round trip. Instead, it keeps a window of segments in flight simultaneously, sliding the window forward as ACKs arrive. Because TCP ACKs are cumulative, one ACK can confirm several segments at once, immediately freeing up that much room in the window. This is the exact mechanism from Chapter 61 (flow control and the sliding window) — `base` here is what that chapter calls the left edge of the window, and `next` is the right edge of what's been sent.

## Hints

<details>
<summary>Hint 1: Trace the state by hand first</summary>

Before writing any test assertions, walk through the sequence of calls step by step on paper: what are `base`, `next`, and `size` after each `Send()` and `Ack()`? This task's grading is unforgiving of off-by-one mistakes, and they're much easier to catch on paper than in a failing test.

</details>

<details>
<summary>Hint 2: The window condition</summary>

The check is `next < base + size`, not `next <= base + size`. A window of size 4 with `base=0` allows sequence numbers 0, 1, 2, 3 to be in flight — that's 4 segments, and the 5th `Send()` (which would set `next=4`) must fail since `4 < 0+4` is false.

</details>

<details>
<summary>Hint 3: Guarding against stale acks</summary>

An ack for `seq` should set `base = seq + 1`, but only if `seq + 1 > base` already. If a duplicate or out-of-order ack arrives for a segment already acknowledged, silently ignore it rather than letting `base` move backward.

</details>

## How to Verify

```bash
lncli run
```
