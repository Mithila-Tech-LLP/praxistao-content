# Chapter 27: Concurrent Mining with Goroutines

Every mining loop this course has written so far — `Run()` in Chapter 25, still `Run()` after Chapter 26's difficulty adjustment — does exactly one thing at a time, on exactly one CPU core: try a nonce, hash it, check it, try the next one. Meanwhile, nearly every real computer sold in the last decade has four, eight, or more CPU cores sitting idle during that search. This chapter fixes that specific waste: splitting the nonce search across every available core using goroutines, coordinating them with a channel and a `context.Context`, and measuring exactly how much faster that makes mining.

## Table of Contents

1. [One Core Is Leaving Performance on the Table](#1-one-core-is-leaving-performance-on-the-table)
2. [The Parallelization Strategy: Striding the Nonce Space](#2-the-parallelization-strategy-striding-the-nonce-space)
3. [Coordinating Goroutines With a Channel and context.Context](#3-coordinating-goroutines-with-a-channel-and-contextcontext)
4. [Implementing RunConcurrent](#4-implementing-runconcurrent)
5. [Avoiding Goroutine Leaks](#5-avoiding-goroutine-leaks)
6. [Wiring RunConcurrent Into MineBlock](#6-wiring-runconcurrent-into-mineblock)
7. [Benchmarking: Before and After](#7-benchmarking-before-and-after)
8. [Diminishing Returns and Amdahl's Law](#8-diminishing-returns-and-amdahls-law)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. One Core Is Leaving Performance on the Table

Recall Chapter 24's central fact about mining: there is no cleverness available to shortcut the nonce search, only brute force — try nonce after nonce until one happens to produce a qualifying hash. That fact hasn't changed and never will (it's the entire security property proof of work relies on). What *can* change is how many nonces get tried per second, and that number scales directly with how much hashing work happens in parallel.

Think of it like a search party looking for a lost hiker across a huge, gridded forest. One person walking the entire grid alone, cell by cell, is going to take a long time — but the *nature* of the search (walk cells, check each one) never gets any smarter no matter how carefully that one person plans their route, because there's no shortcut to "which cell has the hiker." The obvious fix isn't a smarter search — it's more searchers, each covering a different slice of the same grid simultaneously. Ten searchers covering non-overlapping tenths of the grid find the hiker, on average, roughly ten times faster than one searcher covering the whole thing alone.

Go's goroutines (Chapter 05) make this almost embarrassingly easy to build, because the nonce search has a property that makes it what computer scientists call **embarrassingly parallel**: every candidate nonce's hash computation is completely independent of every other candidate's. Trying nonce 41 tells you nothing about nonce 42 (the avalanche effect, again, from Chapter 08), and — crucially for parallelism specifically — computing `Hash(data + 41)` never needs the result of computing `Hash(data + 42)`, or vice versa. There is no shared state to coordinate, no ordering requirement, no risk of one attempt's result depending on another's. This is about as friendly a problem for concurrency as exists in real software.

---

## 2. The Parallelization Strategy: Striding the Nonce Space

Given `N` goroutines, the nonce space (0, 1, 2, 3, ... up to `maxNonce`) needs to be split into `N` non-overlapping pieces, so no two goroutines ever waste effort hashing the exact same nonce. GoChain uses **striding**: goroutine `i` (out of `N` total) tries every nonce whose value, divided by `N`, leaves a remainder of `i`.

```
4 goroutines (N = 4), striding the nonce space:

  Goroutine 0:  0,  4,  8, 12, 16, 20, 24, ...
  Goroutine 1:  1,  5,  9, 13, 17, 21, 25, ...
  Goroutine 2:  2,  6, 10, 14, 18, 22, 26, ...
  Goroutine 3:  3,  7, 11, 15, 19, 23, 27, ...

  Together, every nonce is tried by EXACTLY one goroutine, in order,
  with no gaps and no overlaps.
```

An alternative strategy would split the space into contiguous *ranges* instead — goroutine 0 tries 0 through 249,999, goroutine 1 tries 250,000 through 499,999, and so on. Striding is preferred here for two practical reasons. First, it needs no upfront estimate of how large a range each goroutine should get — with contiguous ranges, guessing a range too small means a goroutine runs out of nonces to try before a solution turns up anywhere; guessing too large wastes memory bookkeeping for nothing. Striding sidesteps the guess entirely: every goroutine's slice is implicitly "every Nth nonce, forever," with no boundary to size correctly. Second, since the winning nonce for any given block is, for all practical purposes, uniformly random across the whole space (Chapter 24, Section 3's geometric distribution), striding guarantees every goroutine has an equal, ongoing chance of being the one to find it — no goroutine is stuck exploring a "cold" region while another gets lucky early.

To see striding actually pay off, trace a tiny, hand-sized example: suppose (purely for illustration — real difficulty needs far more attempts) the winning nonce for some block happens to be 9, and 4 goroutines are searching with `stride = 4`.

```
Single-threaded Run():          RunConcurrent() with 4 goroutines:

  try 0 -> no                     Goroutine 0: try 0 -> no, try 4 -> no, try 8 -> no
  try 1 -> no                     Goroutine 1: try 1 -> no, try 5 -> no, try 9 -> YES!
  try 2 -> no                     Goroutine 2: try 2 -> no, try 6 -> no, try 10 -> (cancelled)
  try 3 -> no                     Goroutine 3: try 3 -> no, try 7 -> no, try 11 -> (cancelled)
  try 4 -> no
  try 5 -> no                     All 4 goroutines make their 3rd attempt at
  try 6 -> no                     ROUGHLY the same moment (they're running
  try 7 -> no                     truly in parallel, not in the interleaved
  try 8 -> no                     order shown above, which is only laid out
  try 9 -> YES!  (10 attempts)    this way for readability). Goroutine 1 wins
                                   on ITS 3rd attempt -- 3 attempts of real
                                   work, not 10, before a winner is found.
```

Single-threaded `Run()` needed 10 sequential attempts to reach nonce 9. The 4-goroutine version needed only 3 *sequential* attempts *per goroutine* running in parallel — each goroutine did less total work, and all of that work happened at the same time rather than one attempt after another. This is the entire mechanism behind Section 7's measured speedup, made concrete on numbers small enough to trace by hand instead of by benchmark.

---

## 3. Coordinating Goroutines With a Channel and context.Context

Splitting the work is the easy half. The harder half is *stopping* the work the moment any one goroutine wins — without that, every worker would keep burning CPU cycles searching for a nonce that no longer matters, for a block whose winner has already been decided. Two tools from Go's standard library do this job together:

- A **channel** (Chapter 05) lets whichever goroutine finds a valid nonce first report it back to whoever is waiting, exactly once.
- **`context.Context`** (new this chapter) is Go's standard mechanism for cancellation: one part of a program can create a `context.Context`, hand copies of it to several goroutines, and later call a `cancel` function that every one of those goroutines can detect — via `ctx.Done()`, a channel that closes the instant cancellation happens — and use as their signal to stop whatever they were doing and return.

```
                    context.WithCancel(...)
                            |
              +-------------+-------------+-------------+
              |             |             |             |
         Goroutine 0   Goroutine 1   Goroutine 2   Goroutine 3
         (stride 0)    (stride 1)    (stride 2)    (stride 3)
              |             |             |             |
              |         FINDS A VALID NONCE!             |
              |             |             |             |
              |      sends result on resultCh            |
              |             |             |             |
      main goroutine receives from resultCh, then
      calls cancel() -- ctx.Done() closes for EVERYONE
              |             |             |             |
              v             v             v             v
          notices        (already      notices       notices
          ctx.Done(),      returned)   ctx.Done(),   ctx.Done(),
          returns                      returns       returns
```

`context.Context` might look like overkill for something a shared `bool` flag plus a mutex could technically also accomplish — but it is the *idiomatic*, standard way to express "cancel this work" in Go specifically because it composes: a `context.Context` can be threaded through any number of layers of function calls (network requests, database queries, nested goroutines) without every layer needing its own bespoke cancellation flag, and it is exactly the tool Volume 7's networking code will reach for repeatedly once a node needs to cancel in-flight peer requests. Learning it here, on a problem simple enough to see the whole picture at once, pays off directly later.

---

## 4. Implementing RunConcurrent

```go
package consensus

import (
	"context"
	"math"
	"math/big"
	"runtime"
	"sync"

	"github.com/you/gochain/crypto"
)

// miningResult carries a winning nonce and the hash it produced back from
// whichever goroutine finds it first.
type miningResult struct {
	nonce uint64
	hash  []byte
}

// RunConcurrent behaves exactly like Run (Chapter 25) -- it returns a
// nonce whose hash satisfies pow.Target -- but splits the search across
// numWorkers goroutines running in parallel, returning as soon as ANY of
// them finds a solution. Passing numWorkers <= 0 defaults to one worker
// per available CPU core.
func (pow *ProofOfWork) RunConcurrent(numWorkers int) (nonce uint64, hash []byte) {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	// ctx is shared by every worker; cancel() is how we tell all of them
	// to stop, the instant any one of them wins. defer cancel() guarantees
	// this happens even if we return early or something panics.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Buffered with capacity 1: the FIRST worker to find a solution can
	// always send without blocking, even before anyone is guaranteed to
	// be listening yet. Section 5 explains why this buffer matters.
	resultCh := make(chan miningResult, 1)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pow.searchWorker(ctx, uint64(i), uint64(numWorkers), resultCh, &wg)
	}

	result := <-resultCh // blocks here until SOME worker reports a win
	cancel()             // tell every other worker to stop immediately
	wg.Wait()             // wait for all of them to actually notice and exit

	return result.nonce, result.hash
}

// searchWorker searches every nonce congruent to `start` modulo `stride`
// -- see Section 2's diagram -- reporting the first qualifying one it
// finds on resultCh, or returning early if ctx is cancelled by a
// different worker winning first.
func (pow *ProofOfWork) searchWorker(
	ctx context.Context,
	start, stride uint64,
	resultCh chan<- miningResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	var hashInt big.Int

	for nonce := start; nonce < math.MaxInt64; nonce += stride {
		// Checking ctx.Done() on literally every iteration would work,
		// but a channel-select has real (if small) overhead, and we're
		// potentially doing millions of iterations per second. Checking
		// every 1024 nonces instead trims that overhead to a rounding
		// error while still cancelling within a few dozen microseconds
		// of another worker winning -- effectively instant to a human.
		if nonce&1023 == 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		data := pow.prepareData(nonce)
		hash := crypto.Hash(data)
		hashInt.SetBytes(hash)

		if hashInt.Cmp(pow.Target) == -1 {
			select {
			case resultCh <- miningResult{nonce: nonce, hash: hash}:
			case <-ctx.Done():
				// Someone else already won a moment before we did; don't
				// block trying to send to a channel nobody will read again.
			}
			return
		}
	}
}
```

`RunConcurrent` is the public entry point: it decides how many workers to launch (`runtime.NumCPU()` by default — Go's standard library function reporting how many logical CPU cores are available), creates a cancellable `context.Context`, launches that many `searchWorker` goroutines each with a different `start` offset but the same `stride` (equal to `numWorkers`), and then blocks on `<-resultCh` until exactly one of them reports a win. The moment it receives a result, it calls `cancel()` — closing `ctx.Done()` for every remaining worker — and `wg.Wait()` blocks until every worker has actually noticed and returned, guaranteeing `RunConcurrent` never returns while a stray goroutine is still quietly running in the background. `searchWorker` is a direct parallel cousin of Chapter 25's `Run`: same `prepareData`/`crypto.Hash`/`big.Int` comparison logic, with two additions — the periodic `ctx.Done()` check, and a `select` around sending the winning result so a worker that wins a nonce-search but loses the race to report it doesn't block forever.

---

## 5. Avoiding Goroutine Leaks

A **goroutine leak** is a goroutine that never returns — it keeps running (and, in this specific case, keeps burning CPU) forever, invisible to the rest of the program, because nothing is left to receive whatever it eventually tries to send or nothing ever tells it to stop. Two separate design choices in Section 4 exist specifically to prevent this, and it's worth naming both explicitly.

**Choice one: `resultCh` is buffered with capacity 1, not unbuffered.** Consider what happens if two workers happen to find valid nonces at nearly the same instant (rare, given how sparse the qualifying sliver is per Chapter 24, but not impossible, especially with many workers running simultaneously). `RunConcurrent` only ever reads from `resultCh` *once*. If the channel were unbuffered, the first worker's send would succeed (matched immediately by that one read), but the *second* worker's send would then block forever — there's no second reader coming, ever — leaking that goroutine permanently. The buffer of size 1 means the first send always succeeds immediately, but critically, it does not fully solve the problem alone: a *third* almost-simultaneous winner would still find the buffer full and the one reader already gone.

**Choice two: the `select` around the send, racing against `ctx.Done()`.** This is what actually closes the gap the buffer leaves open. Once `RunConcurrent` receives its one result and calls `cancel()`, every other worker's next attempt to send on `resultCh` (buffer already full, or no reader left at all) will find `ctx.Done()` also ready to proceed — and Go's `select` statement, when multiple cases are ready, picks one at random, but here only one branch is actually reachable in each losing worker's case (send blocks forever otherwise, `ctx.Done()` doesn't), so the worker reliably takes the `ctx.Done()` branch and returns. Combined, these two choices guarantee **every worker goroutine `RunConcurrent` launches is guaranteed to return**, regardless of how many of them happen to find a winning nonce, and regardless of timing — which is exactly what `wg.Wait()` in Section 4 depends on to avoid `RunConcurrent` itself returning too early while a worker is still alive.

---

## 6. Wiring RunConcurrent Into MineBlock

`core.Blockchain.MineBlock` needs exactly one line changed from Chapter 26's version: swap `pow.Run()` for `pow.RunConcurrent()`.

```go
package core

import (
	"time"

	"github.com/you/gochain/consensus"
)

func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newHeight := prevBlock.Height + 1

	block := NewBlock(transactions, prevBlock.Hash, newHeight)
	block.Timestamp = time.Now().Unix()
	block.Bits = bc.currentDifficultyBits() // Chapter 26

	pow := consensus.NewProofOfWork(block, block.Bits)
	nonce, hash := pow.RunConcurrent(0) // 0 = use every available CPU core

	block.Nonce = nonce
	block.Hash = hash

	bc.blocks = append(bc.blocks, block)
	bc.tip = block.Hash

	return block
}
```

Nothing else about `MineBlock` changes — not the block construction, not the difficulty lookup, not how the winning nonce and hash get written back onto the block. This is a direct, satisfying payoff of Chapter 25's `Run()`/`Validate()` split: because mining and validating are already two separate methods, speeding up mining (this chapter) never touches validation at all. Any node checking a block mined this way calls `Validate()` exactly as before — one hash, one comparison — completely unaware of, and unaffected by, however many goroutines the miner who produced it happened to use.

---

## 7. Benchmarking: Before and After

Talk about speedup is cheap; Go's `testing.B` benchmark support (previewed here, covered fully in Chapter 28) lets us measure it for real.

```go
package consensus

import (
	"testing"

	"github.com/you/gochain/core"
)

func newBenchBlock() *core.Block {
	return &core.Block{
		Height:        1,
		Timestamp:     1234567890,
		PrevBlockHash: []byte("previous-hash-placeholder"),
		MerkleRoot:    []byte("merkle-root-placeholder"),
	}
}

func BenchmarkRun_SingleThreaded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pow := NewProofOfWork(newBenchBlock(), 20) // fixed difficulty for a fair comparison
		pow.Run()
	}
}

func BenchmarkRunConcurrent_AllCores(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pow := NewProofOfWork(newBenchBlock(), 20)
		pow.RunConcurrent(0) // 0 = runtime.NumCPU()
	}
}
```

Run both with `go test ./consensus/... -bench=. -benchtime=5x` (limiting to 5 iterations each, since difficulty-20 mining takes real, noticeable time per attempt). A representative result on an 8-core laptop looks like this:

```
BenchmarkRun_SingleThreaded-8            5    1042318917 ns/op   (~1.04 s/op)
BenchmarkRunConcurrent_AllCores-8        5     148903215 ns/op   (~0.15 s/op)
```

Roughly **7x faster** with 8 cores available — not a full 8x, and Section 8 explains exactly why that gap exists and why it's expected, not a bug. Chapter 28's mini project builds a proper, repeatable benchmarking harness that runs this same comparison across several difficulty levels and prints a clean table, rather than reading raw `go test -bench` output by eye.

---

## 8. Diminishing Returns and Amdahl's Law

It would be a mistake to expect exactly `N`x speedup from `N` goroutines, and it's worth understanding precisely why, rather than just noting "it's usually a bit less." **Amdahl's Law** is the general principle at work here: if a program's total work is split into a portion that *can* run in parallel and a portion that fundamentally *can't*, the maximum possible speedup is capped by the size of that non-parallelizable portion, no matter how many cores you throw at the parallel part.

In `RunConcurrent`'s case, a few real costs don't shrink just because more goroutines are running:

- **Goroutine creation and scheduling overhead.** Launching `N` goroutines and having Go's scheduler distribute them across OS threads and CPU cores isn't free, even though Go goroutines are famously cheap compared to OS threads.
- **The periodic `ctx.Done()` check.** Every 1024 iterations, every worker pays a small tax checking for cancellation — multiply that by `N` workers, and the *total* cancellation-checking work done across the whole program actually grows with `N`, even though each individual worker's share of it stays the same.
- **Memory and cache effects.** `N` goroutines simultaneously computing SHA-256 hashes compete for the same CPU's shared memory bandwidth and cache space in ways one goroutine alone never has to. Modern multi-core CPUs handle this reasonably well for a CPU-bound, independent workload like ours, but "reasonably well" is not "perfectly," which is part of why 8 cores measured closer to 7x than a clean 8x in Section 7's numbers.
- **Operating system noise.** Real machines are never running *only* your benchmark — background processes, other applications, and the OS scheduler itself all compete for the same cores, shaving a little more off real-world speedup versus a theoretical ideal.

None of this means concurrent mining isn't worth it — a consistent 6-7x speedup on an 8-core machine is a dramatic, real improvement over single-threaded mining, and it scales further on machines with even more cores. It does mean that doubling core count should be expected to produce *less than* double the mining speed, an increasingly noticeable gap as core counts grow very large — exactly the kind of measurement Chapter 28's benchmark table makes concrete rather than theoretical.

A quick, illustrative look at how speedup tends to taper off as `numWorkers` climbs, at a fixed difficulty, on the same hypothetical 8-core machine from Section 7:

```
numWorkers    measured speedup vs. single-threaded
    1                1.00x   (RunConcurrent(1) is just Run() with extra overhead)
    2                1.94x
    4                3.78x
    8                6.94x
   16                7.35x   (past the 8 physical cores -- extra goroutines now
                              compete for the SAME 8 cores instead of getting
                              their own, so returns flatten sharply)
   32                7.41x   (essentially flat from here on)
```

The pattern is exactly what Amdahl's Law predicts: speedup climbs quickly and nearly linearly while there are genuinely idle cores to hand work to (1 through 8 workers), then flattens hard the moment `numWorkers` exceeds the number of physical cores actually available, because at that point extra goroutines aren't adding new parallel capacity — they're just adding more contenders for the same fixed 8 cores, plus a little more scheduling overhead for no real gain. This is exactly why `RunConcurrent`'s default (`numWorkers <= 0`) is `runtime.NumCPU()` specifically, rather than some arbitrarily large number: past that point, more goroutines cost more than they help.

---

## Summary

- The nonce search is **embarrassingly parallel** — every candidate nonce's hash is fully independent of every other's — which makes it an unusually clean fit for goroutines, with no shared state to coordinate between attempts.
- **Striding** the nonce space (worker `i` of `N` tries nonces `i, i+N, i+2N, ...`) splits the search evenly across goroutines without needing to guess a range size upfront, and gives every worker an equal chance at the (effectively random) winning nonce.
- A buffered channel (`resultCh`) reports the first winning nonce back to the caller; `context.Context` and its `cancel()` function tell every other goroutine to stop searching the instant a winner is found.
- `RunConcurrent(numWorkers)` launches `numWorkers` (or `runtime.NumCPU()` by default) `searchWorker` goroutines, blocks until the first result arrives, cancels the rest, and waits for every goroutine to actually exit before returning.
- Combining a **buffered** result channel with a **`select` racing against `ctx.Done()`** on every send is what guarantees no goroutine leaks, even when multiple workers happen to find valid nonces nearly simultaneously.
- Wiring `RunConcurrent` into `MineBlock` required changing exactly one line — proof that keeping `Run`/`Validate` cleanly separated in Chapter 25 paid off, since validation is completely untouched by how mining got faster.
- A benchmark comparing `Run` against `RunConcurrent` at the same difficulty measured roughly a 7x speedup on 8 cores — real, but short of a perfect 8x, exactly as **Amdahl's Law** predicts once goroutine overhead, cancellation checks, and shared hardware resources are accounted for.

---

## Exercises

### Easy

1. Explain, in your own words, why the nonce search is "embarrassingly parallel," and give one example of a different problem (mining or otherwise) that would *not* have this property, because one step genuinely depends on the result of a previous step.
2. `searchWorker` checks `ctx.Done()` only every 1024 iterations rather than on every single one. Explain the tradeoff this makes, in terms of responsiveness to cancellation versus per-iteration overhead.
3. If `RunConcurrent(4)` is called on an 8-core machine, would you expect roughly the same speedup as `RunConcurrent(0)`? Why or why not?

### Medium

4. `resultCh` has a buffer of exactly 1. Modify `RunConcurrent` to accept up to `numWorkers` results without any worker ever blocking on its send (regardless of the `select`/`ctx.Done()` fallback), and explain whether this change affects correctness, performance, or neither.
5. Write a test, `TestRunConcurrent_MatchesRunResult`, that mines the *same* block (same fields, same difficulty) once with `Run()` and once with `RunConcurrent()`, and asserts that both returned nonces independently satisfy `pow.Validate()` — without assuming they'll be the *same* nonce (since striding means different workers may find different, equally valid solutions first).
6. Rewrite `searchWorker` to use contiguous ranges instead of striding (Section 2's alternative strategy), taking an explicit `[start, end)` range per worker instead of a stride. What new problem does `RunConcurrent` now have to solve that striding avoided for free?

### Hard

7. Instrument `searchWorker` (with an atomic counter or similar) to report exactly how many total hash attempts were made across all workers combined before a solution was found, for both `Run` and `RunConcurrent` at the same difficulty and same block data. Run each several times and discuss whether the *total* attempts across all workers in the concurrent version is roughly the same as the single-threaded version's attempt count, more, or less — and explain why, referencing the geometric distribution from Chapter 24, Section 3.
8. Profile `RunConcurrent` using Go's built-in `pprof` support (`go test -cpuprofile`) at a high difficulty (24+ bits) and identify where CPU time is actually being spent across all goroutines. Report whether `crypto.Hash` dominates as expected, or whether something else (channel operations, `big.Int` allocation, the `ctx.Done()` check) takes up a surprising fraction of total time.
9. Design and implement a variant, `RunConcurrentWithProgress(numWorkers int, progress chan<- uint64)`, that periodically reports the *total* number of nonces tried so far (summed across all workers) on a second channel, without meaningfully slowing down the search itself. Discuss what synchronization primitive you chose for the shared counter and why it's safe under concurrent increments from multiple goroutines.
