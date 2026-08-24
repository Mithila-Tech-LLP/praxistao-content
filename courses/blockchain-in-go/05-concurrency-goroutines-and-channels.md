# Chapter 05: Concurrency — Goroutines and Channels

A running blockchain node juggles mining, listening for peers, receiving transactions, and answering API requests, all at the same time. This chapter teaches you Go's two core tools for that — goroutines and channels — by building a small, complete, working simulation: fake transactions arriving every second and getting collected into a mini mempool, the exact pattern GoChain's real mempool will use starting in Volume 5.

## Table of Contents

1. [Why a Blockchain Node Needs to Do Many Things at Once](#1-why-a-blockchain-node-needs-to-do-many-things-at-once)
2. [Goroutines: Lightweight, Cheap Concurrency](#2-goroutines-lightweight-cheap-concurrency)
3. [Channels: Safe Handoff Between Goroutines](#3-channels-safe-handoff-between-goroutines)
4. [The Producer-Consumer Pattern](#4-the-producer-consumer-pattern)
5. [A Full Worked Example: The Fake Transaction Feed](#5-a-full-worked-example-the-fake-transaction-feed)
6. [Race Conditions, and Why Channels Avoid Them](#6-race-conditions-and-why-channels-avoid-them)
7. [The select Statement](#7-the-select-statement)
8. [Shutting Down Cleanly](#8-shutting-down-cleanly)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why a Blockchain Node Needs to Do Many Things at Once

Picture a single GoChain node running on your laptop, fully operational. At any given moment, it might be:

- Trying nonces as fast as possible, searching for a valid proof-of-work solution (Volume 4).
- Listening on a network port for new peers trying to connect (Volume 7).
- Receiving a freshly submitted transaction from a wallet and needing to validate it and add it to the mempool (Volume 5).
- Answering an HTTP request from a block explorer asking for the current chain height (Volume 10).

None of these tasks are naturally "one after another" — a peer can connect at any moment, independent of whether mining is mid-attempt or an API request just arrived. A program that could only do one thing at a time would have to awkwardly interleave all of this manually, checking "is there a new connection? no? is there a new transaction? no? okay, try one more mining nonce" in a tight loop — fragile, hard to reason about, and wasteful.

**Concurrency** is a program's ability to structure and manage multiple tasks that are logically independent of each other, whether or not they're running at the exact same physical instant. Go was designed, from its very first release, around making concurrency easy to write correctly — and this chapter covers the two tools that make it possible: **goroutines** and **channels**.

---

## 2. Goroutines: Lightweight, Cheap Concurrency

A **goroutine** is a function that runs concurrently with the rest of your program. You start one by writing the `go` keyword directly in front of a function call:

```go
package main

import (
	"fmt"
	"time"
)

func mineForever() {
	nonce := 0
	for {
		nonce++
		if nonce%1_000_000 == 0 {
			fmt.Println("still mining, tried", nonce, "nonces")
		}
	}
}

func main() {
	go mineForever() // starts running concurrently, in the background

	fmt.Println("main goroutine continues immediately, without waiting")
	time.Sleep(2 * time.Second) // give mineForever a moment to print something
	fmt.Println("main is done")
}
```

`go mineForever()` launches `mineForever` as a new goroutine and returns **immediately** — `main` does not wait for `mineForever` to finish (it never does, since it loops forever) before moving on to the next line. This is the entire point: `main` and `mineForever` now run concurrently, side by side, each making progress independently.

Every Go program already has at least one goroutine running from the start — the one executing `main` itself. `go someFunc()` simply adds another one alongside it. Go's runtime schedules potentially many thousands of goroutines onto a much smaller number of real operating system threads automatically, which is exactly why Chapter 02 could claim that spawning one goroutine per connected network peer is a completely reasonable design, even with thousands of peers.

One crucial gotcha worth calling out immediately: if `main` returns before a goroutine finishes its work, the whole program exits immediately, and any goroutines still running are simply cut off mid-execution — there's no "wait for background work to finish" behavior by default. The `time.Sleep(2 * time.Second)` in the example above is a crude, temporary way to give `mineForever` a moment to run before `main` exits; Section 4 and beyond show the proper way to coordinate this using channels instead of guessing at a sleep duration.

---

## 3. Channels: Safe Handoff Between Goroutines

Goroutines running side by side are useful on their own, but they usually need to *communicate* — one goroutine discovers something, and another goroutine needs to receive it, safely, without both of them reading and writing the same data at the exact same moment and corrupting it.

A **channel** is a typed pipe: one goroutine sends a value into one end, and another goroutine receives it from the other end. Go's own design philosophy, often quoted directly from the language's documentation, puts it memorably: *"Do not communicate by sharing memory; instead, share memory by communicating."* Rather than two goroutines both reaching into the same shared variable (risky, as Section 6 demonstrates), one goroutine hands a value *through* a channel to another, and Go's runtime guarantees that handoff is safe.

```go
package main

import "fmt"

func main() {
	messages := make(chan string) // a channel that carries string values

	go func() {
		messages <- "Alice lent Bob 10 gochips" // SEND a value into the channel
	}()

	received := <-messages // RECEIVE a value from the channel (blocks until one arrives)
	fmt.Println(received)
}
```

`make(chan string)` creates a new, unbuffered channel that carries `string` values. Inside the anonymous goroutine (`func() { ... }()`, a function defined and immediately launched with `go` in one expression), `messages <- "..."` **sends** a value into the channel. Back in `main`, `<-messages` **receives** a value from the channel — critically, this receive operation *blocks* (pauses `main`'s execution) until a value actually arrives, which is exactly what lets `main` safely wait for the goroutine's result without an arbitrary `time.Sleep` guess.

```
   goroutine                    channel                    main goroutine
+------------+              +-----------+              +----------------+
|  computes  |   send  -->  |  "Alice   |  <-- receive |   blocks until  |
|  a message |              |  lent Bob |              |   a value is    |
|            |              |  10 gochips" |            |   available    |
+------------+              +-----------+              +----------------+
```

An **unbuffered channel** (the kind created by plain `make(chan string)`) has no internal storage at all — a send only completes once a matching receive is ready to accept it, and vice versa, which is why this handoff is inherently synchronized. Go also supports **buffered channels**, created with `make(chan string, capacity)`, which can hold a limited number of values without a receiver being immediately ready; Section 5 uses one of these for the mempool feed.

---

## 4. The Producer-Consumer Pattern

The **producer-consumer pattern** is one of the most common and useful concurrency shapes: one or more goroutines ("producers") generate items and send them into a channel, while one or more separate goroutines ("consumers") receive those items from the channel and do something with each one. Neither side needs to know how many of the other kind exist, or exactly when they'll produce or consume — the channel handles all of the coordination.

This maps directly onto a real blockchain node's mempool: new transactions arrive from many different sources (local wallet submissions, gossiped transactions from peers) — these are the producers — and the mempool itself is a consumer, continuously receiving new transactions and adding them to its pending pool, ready to be picked up by the miner later.

```
   Producer 1 ---\
                  \
   Producer 2 -----> [ channel ] -----> Consumer (collects into a slice)
                  /
   Producer 3 ---/
```

---

## 5. A Full Worked Example: The Fake Transaction Feed

Let's build this pattern end to end, with a producer goroutine that "discovers" a fake transaction once per second, sending each one over a channel to a consumer goroutine that collects them into a slice acting as a tiny mempool.

**Note on scope:** this is explicitly a preview and a simplification. The real `core.Mempool`, arriving in Volume 5, handles actual signed `core.Transaction` values, checks for double-spends, and integrates with the miner. What we're building here is the concurrency *shape* of that system — the plumbing — using a fake, simplified transaction type, so that when the real mempool arrives, the pattern is already familiar.

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// fakeTransaction is a deliberately simplified stand-in for core.Transaction,
// which arrives for real in Volume 5 with fields like Inputs, Outputs, and
// a proper cryptographic ID.
type fakeTransaction struct {
	ID     int
	Amount int64 // amount in gochips
}

// produceTransactions simulates transactions arriving from the outside
// world once per second, sending each one into txs. It closes txs when
// stop is closed, signaling the consumer that no more values are coming.
func produceTransactions(txs chan<- fakeTransaction, stop <-chan struct{}) {
	id := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			close(txs) // no more values will ever be sent; let the consumer know
			return
		case <-ticker.C:
			id++
			tx := fakeTransaction{ID: id, Amount: int64(rand.Intn(100) + 1)}
			txs <- tx // SEND the new fake transaction into the channel
			fmt.Printf("[producer] discovered tx #%d for %d gochips\n", tx.ID, tx.Amount)
		}
	}
}

// collectMempool is the consumer: it receives every fake transaction sent
// into txs and appends it to an in-memory slice, acting as a miniature
// mempool. It signals completion on done once txs is closed and drained.
func collectMempool(txs <-chan fakeTransaction, done chan<- []fakeTransaction) {
	var mempool []fakeTransaction
	for tx := range txs { // loops until txs is closed AND fully drained
		mempool = append(mempool, tx)
	}
	done <- mempool
}

func main() {
	txs := make(chan fakeTransaction)      // producer -> consumer handoff
	stop := make(chan struct{})            // main -> producer shutdown signal
	done := make(chan []fakeTransaction)   // consumer -> main final result

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		produceTransactions(txs, stop)
	}()

	go collectMempool(txs, done)

	time.Sleep(5500 * time.Millisecond) // let a handful of fake txs arrive
	close(stop)                          // tell the producer to wind down
	wg.Wait()                            // wait for the producer goroutine to actually exit

	mempool := <-done // wait for the consumer to hand back the final slice
	fmt.Printf("\nfinal mempool has %d pending transactions:\n", len(mempool))
	for _, tx := range mempool {
		fmt.Printf("  tx #%d: %d gochips\n", tx.ID, tx.Amount)
	}
}
```

Let's walk through the moving pieces carefully, since there's a lot packed into this example.

`txs chan<- fakeTransaction` in `produceTransactions`'s signature and `txs <-chan fakeTransaction` in `collectMempool`'s signature are **directional channel types** — `chan<-` means "this function may only send into this channel," and `<-chan` means "this function may only receive from it." This isn't strictly required (you could pass a plain bidirectional `chan fakeTransaction` to both), but it's good practice: it lets the compiler catch a mistake like accidentally trying to receive from a channel a function is only supposed to produce into.

`time.NewTicker(1 * time.Second)` creates a ticker that sends a value on its `.C` channel once per second, forever, until stopped — exactly what simulates "a new transaction arrives roughly once a second." `select { case <-stop: ... case <-ticker.C: ... }` (covered properly in Section 7) lets `produceTransactions` wait on *either* the shutdown signal *or* the next tick, whichever happens first.

`close(txs)` closes the channel, signaling "no more values will ever be sent here." On the consumer side, `for tx := range txs` is a clean, idiomatic way to receive values from a channel repeatedly until it's closed *and* fully drained of any values sent before the close — the loop then exits naturally, no manual "did I get a special stop signal" check required.

`sync.WaitGroup` is a small counter used to wait for a group of goroutines to finish: `wg.Add(1)` increments the counter, `wg.Done()` (called via `defer` so it always runs, even if the goroutine were to panic) decrements it, and `wg.Wait()` blocks until the counter returns to zero. Here it's used specifically so `main` can be sure the producer goroutine has genuinely finished (and therefore has called `close(txs)`) before moving on.

Running this program prints one line per second as fake transactions arrive, then, after `main` closes `stop`, prints the final collected mempool contents — a complete, working, concurrent producer-consumer pipeline in well under 60 lines. This is precisely the shape `core.Mempool` will take starting in Volume 5, except the producer there will be real network and API code delivering genuinely signed `core.Transaction` values, not a `time.Ticker` generating random amounts.

---

## 6. Race Conditions, and Why Channels Avoid Them

A **race condition** happens when two or more goroutines access the same piece of shared memory at the same time, at least one of them writing to it, without any coordination — the final result becomes unpredictable, depending on the exact, essentially random timing of which goroutine happens to run first. Let's see this go wrong directly, using a shared map instead of a channel:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	mempool := make(map[int]int64) // tx ID -> amount, shared across goroutines

	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mempool[id] = int64(id * 10) // UNSAFE: concurrent map write, no lock
		}(i)
	}
	wg.Wait()

	fmt.Println("mempool size:", len(mempool)) // unreliable — may even crash
}
```

This program spawns 1,000 goroutines, each writing to the *same* `mempool` map concurrently, with no protection whatsoever. In practice, running this will often crash outright with Go's runtime reporting `"fatal error: concurrent map writes"` — Go's built-in maps are explicitly documented as unsafe for concurrent writes, precisely to force this class of bug into the open loudly rather than silently corrupting data. Even in the rare case it doesn't crash, the *contents* of the map afterward are not reliably guaranteed to be all 1,000 entries — some writes can be lost when two goroutines happen to collide at exactly the wrong moment.

The traditional fix, if you must share memory directly like this, is a **mutex** (mutual exclusion lock) — a lock one goroutine holds while writing, that forces every other goroutine to wait its turn:

```go
var mu sync.Mutex
// ...
go func(id int) {
	defer wg.Done()
	mu.Lock()
	mempool[id] = int64(id * 10) // now safe: only one goroutine at a time
	mu.Unlock()
}(i)
```

This works, but it puts the burden entirely on the programmer to remember to lock and unlock correctly, every single time, everywhere the shared map is touched — forget one `mu.Lock()` call anywhere in a large codebase, and the bug is back, often intermittently and hard to reproduce.

This is exactly the problem channels are designed to sidestep. In Section 5's example, `mempool` (the slice built up inside `collectMempool`) is only ever touched by *one* goroutine — the consumer. Producers never reach into it directly; they only ever send values *through* the channel, and the channel itself guarantees that handoff is safe. There is no shared memory to protect with a mutex in the first place, because there's no memory being shared — only messages being passed. This is the concrete meaning behind "share memory by communicating" from Section 3: instead of many goroutines fighting over one mutable map, one goroutine owns the data exclusively, and everyone else talks to it through a channel.

```
  UNSAFE (shared map, no lock):          SAFE (channel handoff):

  G1 --\                                 G1 --\
  G2 ---+--> mempool map <--- CRASH      G2 ---+--> channel --> single owner
  G3 --/     (concurrent writes)         G3 --/     goroutine owns the slice
```

---

## 7. The select Statement

The `select` statement, already used in Section 5, lets a goroutine wait on *multiple* channel operations at once, proceeding with whichever one becomes ready first — think of it as a `switch` statement, but for channels instead of values.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	peerMessages := make(chan string)
	apiRequests := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		peerMessages <- "new block from peer 7"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		apiRequests <- "GET /balance/alice"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-peerMessages:
			fmt.Println("handling peer message:", msg)
		case req := <-apiRequests:
			fmt.Println("handling API request:", req)
		}
	}
}
```

`select` blocks until *at least one* of its `case` branches has a value ready, then runs that branch's code — if the peer message arrives before the API request (as it does here, after 1 second versus 2), the `peerMessages` case fires first. If multiple channels happen to be ready at the exact same moment, `select` picks one at random among them, rather than always favoring the first `case` written — a deliberate design choice that prevents code from accidentally depending on the order `case`s happen to be listed.

`select` can also include a `default` case, which runs immediately if *no* channel is ready yet, letting a goroutine check "is anything waiting for me right now?" without blocking at all — useful for a node that wants to do a bit of other work (like continuing a mining attempt) whenever there's nothing more urgent to handle.

---

## 8. Shutting Down Cleanly

Section 5's example used a `stop chan struct{}`, closed with `close(stop)`, as a shutdown signal. This is a common and idiomatic Go pattern worth naming explicitly: `struct{}` (an empty struct) is used here specifically because it carries no data at all — the *only* thing that matters is whether the channel has been closed, not any value sent through it. Closing a channel causes every goroutine currently receiving from it (via `<-stop` inside a `select`, as in `produceTransactions`) to immediately unblock, letting many goroutines all learn about a shutdown request from a single `close` call.

This pattern — a dedicated "done" or "stop" channel, checked inside a `select` alongside a goroutine's normal work — is the foundation GoChain's real node shutdown logic will build on starting in Volume 7, where a node needs to cleanly stop mining, close every open peer connection, and flush any pending writes to disk, all without leaving something half-finished. Go's standard library also provides `context.Context`, a more feature-rich tool for exactly this kind of cancellation signaling, which Chapter 27 introduces properly once concurrent mining makes the need for it concrete.

---

## Summary

- A blockchain node genuinely needs to do many independent things at once — mining, networking, receiving transactions, serving API requests — which is why concurrency is a first-class concern from the very first working node.
- A **goroutine**, started with `go someFunc()`, is a lightweight, concurrently running function; `main` does not wait for it unless you explicitly make it wait.
- A **channel** (`make(chan T)`) is a typed pipe for safely handing values between goroutines; sends and receives on an unbuffered channel block until a matching partner is ready.
- The **producer-consumer pattern** — one or more goroutines generating items into a channel, one or more goroutines receiving and processing them — is exactly the shape of GoChain's real mempool, previewed here with a fake, once-per-second transaction feed.
- A **race condition** occurs when multiple goroutines access shared memory (like a map) concurrently with at least one write and no coordination; Go's map writes are explicitly unsafe for this and will often crash loudly rather than silently corrupt data.
- Channels avoid race conditions by ensuring only one goroutine ever owns a given piece of mutable state directly, with everyone else communicating through the channel instead of reaching into shared memory.
- `select` lets a goroutine wait on multiple channel operations at once, proceeding with whichever becomes ready first; a `default` case makes a `select` non-blocking.
- A closed `stop chan struct{}`, checked inside a `select`, is the standard, lightweight pattern for signaling shutdown to one or many goroutines at once, and is the direct ancestor of the cleaner-shutdown logic GoChain's real node will need starting in Volume 7.

---

## Exercises

### Easy

1. **Modify the goroutine example in Section 2** so that `mineForever` also accepts and increments a shared `*int64` counter of total nonces tried, passed in as a parameter, and have `main` print that counter's value every half second using a loop and `time.Sleep`. Note (but don't yet fix) that reading the counter from `main` while `mineForever` writes to it concurrently is technically a data race — we'll return to this class of problem in Exercise 6.

2. **Write a program with two goroutines**: one sends the numbers 1 through 10 (one at a time) into an unbuffered `chan int`, and the other receives and prints each one, prefixed with `"received: "`. Add a `sync.WaitGroup` so `main` correctly waits for both goroutines to finish before exiting.

3. **Take the fake transaction feed from Section 5** and change the producer's interval from once per second to once every 250 milliseconds, and change `main`'s `time.Sleep` so roughly 12-15 fake transactions are collected before shutdown. Run it and confirm the final printed count matches your expectation.

---

### Medium

4. **Add a second, independent producer goroutine** to the Section 5 example, simulating transactions arriving from a different "source" (give it a different ID range or a distinguishing field), sending into the *same* `txs` channel the original producer uses. Make sure `close(txs)` is only called once both producers have genuinely finished (hint: you'll need a `sync.WaitGroup` around both producer goroutines, closing `txs` only after `wg.Wait()` returns), and confirm the consumer still correctly collects transactions from both sources into one mempool slice.

5. **Reproduce the race condition from Section 6 on purpose**, then run it with Go's built-in race detector: `go run -race yourfile.go`. Paste (or describe in detail) the race detector's output, and explain, in your own words, what specific two lines of code it identifies as racing with each other and why.

6. **Fix the race condition from Exercise 5 two different ways**: once using a `sync.Mutex` around the map access, and once by redesigning the code to use a channel-based handoff instead (a single goroutine owns the map, and the 1,000 goroutines each send their `(id, amount)` pair through a channel instead of writing directly). Run `go run -race` against both fixed versions and confirm neither reports a race anymore. Write two or three sentences comparing which fix felt more natural to write correctly.

---

### Hard

7. **Build a "multi-miner" simulation**: launch 4 goroutines that each independently generate random guesses (simulate nonce attempts as random `int` values) looking for one that satisfies a toy condition (say, divisible by 100,000), sending only the *first* winning guess found across all 4 goroutines into a shared `chan int` with capacity 1, and having the other goroutines stop trying once a winner has been sent (use a `stop chan struct{}`, closed once, checked inside each miner's loop via `select` with a `default` case). Report how many total attempts (summed across all 4 goroutines) it took to find a winner, run the whole thing 5 times, and report the range of attempt counts you observed. This exercise directly previews the concurrent mining design covered fully in Chapter 27.

8. **Extend the fake transaction feed from Section 5 to support a maximum mempool size of 20.** Once the mempool (inside `collectMempool`) reaches 20 pending transactions, additional incoming transactions should be dropped, with a printed message noting each drop and why — a simplified version of a real mempool's eviction policy under load, foreshadowing Chapter 34's more complete treatment of the topic. Make sure this size check happens safely, without introducing any new race condition (there should still be exactly one goroutine that owns and modifies the mempool slice).

9. **Design and implement a `select`-based "priority" scheduler**: given three channels — `highPriority`, `normalPriority`, and `lowPriority` — all carrying string messages, write a goroutine that always prefers to process a `highPriority` message if one is available, falls back to `normalPriority` if not, and only processes `lowPriority` if neither of the other two has anything waiting (hint: you'll need a `select` with a `default` case nested inside an outer `select`, or an explicit priority-checking loop — research and choose whichever approach you think is cleaner, and justify your choice in a short comment). Feed all three channels from separate goroutines with different timings and demonstrate, with printed output, that high-priority messages are consistently handled before lower-priority ones arrive, even when a low-priority message technically arrived first.
