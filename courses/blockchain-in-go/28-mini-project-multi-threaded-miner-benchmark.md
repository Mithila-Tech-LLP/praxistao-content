# Chapter 28: Mini Project — Multi-Threaded Miner Benchmark

Chapter 27 measured one speedup number, at one difficulty level, with one raw `go test -bench` run. That's a start, but it doesn't yet answer the question anyone actually cares about in practice: *how does the speedup from concurrency change as difficulty rises?* Does single-threaded mining fall further behind as difficulty increases, or does the gap stay roughly constant? This chapter builds `minebench`, a small, reusable benchmarking tool that answers exactly that — across several difficulty levels, single-threaded versus concurrent, with a clean printed comparison table you generate yourself rather than take on faith.

## Table of Contents

1. [What minebench Needs to Measure](#1-what-minebench-needs-to-measure)
2. [Hashes Per Second: A Fairer Unit Than Wall-Clock Time Alone](#2-hashes-per-second-a-fairer-unit-than-wall-clock-time-alone)
3. [Designing the Harness](#3-designing-the-harness)
4. [Instrumenting the Miner to Count Attempts](#4-instrumenting-the-miner-to-count-attempts)
5. [Building the Comparison Table](#5-building-the-comparison-table)
6. [Mini Project: minebench](#mini-project-minebench)
7. [Reading Your Own Results](#7-reading-your-own-results)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What minebench Needs to Measure

Three things need to vary independently for this benchmark to actually answer the motivating question: **difficulty** (how hard each block is to mine), **mode** (single-threaded `Run()` versus concurrent `RunConcurrent()`), and, implicitly, **hash rate** (how many attempts per second each mode actually achieves). The tool's job is to hold everything else fixed, vary difficulty across a small range, run both modes at each level, and report:

- **Total time** to mine one block at that difficulty, in that mode.
- **Hashes per second** achieved — Section 2 explains why this, not total time alone, is the number that actually generalizes across difficulty levels.
- The **speedup ratio** between the two modes at each difficulty.

```
minebench's job, in one picture:

  for each difficulty in {16, 18, 20, 22}:
      time0 := mine one block, single-threaded
      time1 := mine one block, concurrently (all cores)
      print: difficulty | single-threaded time | concurrent time | speedup
```

---

## 2. Hashes Per Second: A Fairer Unit Than Wall-Clock Time Alone

Wall-clock time to mine one block is a real, meaningful number, but it's a noisy one on its own: recall from Chapter 24, Section 5 that mining time is only an *expected* value — any single attempt, at any difficulty, might get lucky and finish in a fraction of the expected time, or unlucky and take several times longer, purely by chance. Comparing two single "mine one block" timings, one per mode, risks comparing two random draws rather than two genuinely different rates of work.

**Hashes per second** (often abbreviated H/s, the same unit real mining hardware is marketed and compared by) sidesteps most of this: it's computed as `attempts / elapsedTime`, and while any *one* mining run's attempt count is still randomly distributed around its expected value, dividing by elapsed time turns "how long did this particular lucky-or-unlucky run take" into "roughly how fast is this mode actually chewing through the search space" — a rate, not a single data point. Running several trials per difficulty level and averaging their H/s (Section 3 does exactly this) further smooths out the remaining randomness, giving a number that meaningfully compares across different difficulty levels, where raw "seconds to mine one block" would not, since higher difficulty means more attempts by design (Chapter 24, Section 3) — a longer time at higher difficulty doesn't necessarily mean a *slower* miner, but a raw time-only comparison would make it look that way.

```
Total time alone:                       Hashes per second:
  16 bits: 0.07s                          16 bits: ~940,000 H/s
  20 bits: 1.05s     <- looks "slower"    20 bits: ~998,000 H/s   <- actually
                        but isn't --                                 about the
                        it just had                                  SAME rate,
                        16x more work                                 as expected
                        to do (Ch 24)
```

---

## 3. Designing the Harness

`minebench` needs to mine a fresh, minimal block at a chosen difficulty, in a chosen mode, and report exactly how long it took and how many hash attempts it made. That second number — attempt count — doesn't come for free from `Run()` or `RunConcurrent()` as built in Chapters 25 and 27, since neither one currently returns it. Section 4 fixes that with a small, non-invasive instrumentation addition.

```go
package minebench

import "time"

// Result holds one trial's measurements: how long mining took, and how
// many total hash attempts were made across every goroutine involved
// (1, for single-threaded; numWorkers, for concurrent).
type Result struct {
	Mode       string        // "single" or "concurrent"
	Difficulty int           // leading zero bits required
	Elapsed    time.Duration
	Attempts   uint64
}

// HashesPerSecond is the rate this one trial achieved -- the fairer,
// difficulty-independent number Section 2 argued for.
func (r Result) HashesPerSecond() float64 {
	return float64(r.Attempts) / r.Elapsed.Seconds()
}
```

`Result` is deliberately minimal: enough fields to reconstruct every number Section 1 asked for (`HashesPerSecond` is a derived method, not a stored field, so it's always freshly computed from the two raw measurements that actually matter — `Attempts` and `Elapsed`).

---

## 4. Instrumenting the Miner to Count Attempts

Neither `Run()` (Chapter 25) nor `RunConcurrent()` (Chapter 27) reports how many nonces were actually tried — both simply return the winning nonce and hash. Since the winning nonce itself is a reasonable *proxy* for attempt count in the single-threaded case (nonce `N` winning means roughly `N + 1` attempts were made, since the search starts at 0 and increments by 1 every time), but is **not** a valid proxy in the concurrent case (a stride-4 worker winning at nonce 4,000 made only 1,000 attempts of its own, not 4,000, and three other workers made attempts too), `minebench` adds a small atomic counter rather than changing either method's signature or return values.

```go
package consensus

import (
	"context"
	"math"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/you/gochain/crypto"
)

// RunConcurrentCounted behaves exactly like RunConcurrent, but also
// returns the TOTAL number of hash attempts made across every worker
// combined -- the number RunConcurrent itself has no reason to compute
// in normal operation, but minebench needs for Section 2's H/s metric.
func (pow *ProofOfWork) RunConcurrentCounted(numWorkers int) (nonce uint64, hash []byte, attempts uint64) {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan miningResult, 1)
	var totalAttempts uint64 // shared, incremented via atomic ops only

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go pow.countedSearchWorker(ctx, uint64(i), uint64(numWorkers), resultCh, &wg, &totalAttempts)
	}

	result := <-resultCh
	cancel()
	wg.Wait()

	return result.nonce, result.hash, atomic.LoadUint64(&totalAttempts)
}

// countedSearchWorker is searchWorker (Chapter 27) with one addition: it
// increments a shared counter after every attempt, using atomic
// operations so concurrent increments from multiple goroutines never
// race or lose an update.
func (pow *ProofOfWork) countedSearchWorker(
	ctx context.Context,
	start, stride uint64,
	resultCh chan<- miningResult,
	wg *sync.WaitGroup,
	totalAttempts *uint64,
) {
	defer wg.Done()

	var hashInt big.Int
	var localAttempts uint64 // batched locally, flushed periodically -- see below

	for nonce := start; nonce < math.MaxInt64; nonce += stride {
		localAttempts++

		if nonce&1023 == 0 {
			// Flush the local counter into the shared one only every
			// 1024 iterations, not every single one -- an atomic add is
			// cheap, but not free, and we're doing millions of these per
			// second per worker. This is the same "check periodically,
			// not constantly" tradeoff Chapter 27 made for ctx.Done().
			atomic.AddUint64(totalAttempts, localAttempts)
			localAttempts = 0

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
			atomic.AddUint64(totalAttempts, localAttempts+1) // flush the final, partial batch
			select {
			case resultCh <- miningResult{nonce: nonce, hash: hash}:
			case <-ctx.Done():
			}
			return
		}
	}
}
```

`RunConcurrentCounted` is `RunConcurrent` with one extra return value and one extra shared variable, `totalAttempts`, updated only through `atomic.AddUint64` and read only through `atomic.LoadUint64` — Go's standard, lock-free way to safely share a single counter across goroutines without a full `sync.Mutex` (recall Chapter 05's warning about unprotected shared state: without atomics here, multiple goroutines incrementing the same plain `uint64` concurrently could silently lose updates). `countedSearchWorker` batches its own increments into a `localAttempts` variable and only flushes them into the shared counter every 1024 iterations — the same "batch the expensive shared operation, don't pay its cost every single time" idea Chapter 27 already used for checking `ctx.Done()`, applied here to the atomic add instead. The single-threaded equivalent needs no such machinery, since there's only one goroutine to begin with:

```go
// RunCounted is Run (Chapter 25) with an attempt count returned alongside
// the winning nonce and hash -- trivial here, since there's exactly one
// goroutine and the loop variable itself IS the attempt count.
func (pow *ProofOfWork) RunCounted() (nonce uint64, hash []byte, attempts uint64) {
	nonce, hash = pow.Run()
	return nonce, hash, nonce + 1 // nonce started at 0, so N wins after N+1 tries
}
```

---

## 5. Building the Comparison Table

With both counted variants in place, `minebench`'s harness runs a small matrix of trials — every combination of difficulty level and mode — and prints the results as a table.

```go
package minebench

import (
	"fmt"

	"github.com/you/gochain/consensus"
	"github.com/you/gochain/core"
)

// Difficulties is the set of leading-zero-bit levels minebench compares.
// Kept deliberately small and fast (16-22 bits) so the whole benchmark
// finishes in well under a minute on an ordinary laptop; Chapter 24's
// table already showed why 24+ bits would make this tool slow to run.
var Difficulties = []int{16, 18, 20, 22}

func benchBlock() *core.Block {
	return &core.Block{
		Height:        1,
		Timestamp:     1234567890,
		PrevBlockHash: []byte("previous-hash-placeholder"),
		MerkleRoot:    []byte("merkle-root-placeholder"),
	}
}

// RunAll mines one block at every difficulty in Difficulties, once
// single-threaded and once concurrently, and returns every trial's Result
// in the order they were run.
func RunAll(trialsPerDifficulty int) []Result {
	var results []Result

	for _, bits := range Difficulties {
		for trial := 0; trial < trialsPerDifficulty; trial++ {
			results = append(results, runSingleTrial(bits))
			results = append(results, runConcurrentTrial(bits))
		}
	}

	return results
}

func runSingleTrial(bits int) Result {
	pow := consensus.NewProofOfWork(benchBlock(), bits)

	start := time.Now()
	_, _, attempts := pow.RunCounted()
	elapsed := time.Since(start)

	return Result{Mode: "single", Difficulty: bits, Elapsed: elapsed, Attempts: attempts}
}

func runConcurrentTrial(bits int) Result {
	pow := consensus.NewProofOfWork(benchBlock(), bits)

	start := time.Now()
	_, _, attempts := pow.RunConcurrentCounted(0) // 0 = every available core
	elapsed := time.Since(start)

	return Result{Mode: "concurrent", Difficulty: bits, Elapsed: elapsed, Attempts: attempts}
}
```

`RunAll` is the whole benchmark matrix from Section 1's diagram, made real: for every difficulty level, run `trialsPerDifficulty` single-threaded trials and the same number of concurrent trials, collecting every individual `Result`. Running more than one trial per difficulty matters directly because of Section 2's point about randomness — a single trial's H/s is still one noisy sample, and Section 6's `main.go` averages across trials before printing, rather than reporting any one run's number as if it were exact.

---

## Mini Project: minebench

Here is the complete, runnable command-line tool tying every piece above together, plus the averaging and table-printing logic that turns raw `Result` values into the comparison table this chapter set out to build.

```go
// cmd/minebench/main.go
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/you/gochain/minebench"
)

func main() {
	trials := flag.Int("trials", 3, "number of trials to average per difficulty level")
	flag.Parse()

	fmt.Printf("minebench: mining at difficulties %v, %d trial(s) each, single-threaded vs. concurrent\n\n",
		minebench.Difficulties, *trials)

	results := minebench.RunAll(*trials)
	printTable(results, *trials)
}

// row holds one difficulty level's averaged single-threaded and
// concurrent measurements, ready to print as one table line.
type row struct {
	difficulty      int
	singleAvgTime   time.Duration
	singleAvgHashes float64
	concAvgTime     time.Duration
	concAvgHashes   float64
}

func printTable(results []minebench.Result, trials int) {
	rows := aggregate(results, trials)

	fmt.Printf("%-6s | %-14s | %-16s | %-14s | %-16s | %-8s\n",
		"Bits", "Single Time", "Single H/s", "Concur. Time", "Concur. H/s", "Speedup")
	fmt.Println("-------|----------------|------------------|----------------|------------------|---------")

	for _, r := range rows {
		speedup := r.concAvgHashes / r.singleAvgHashes
		fmt.Printf("%-6d | %-14s | %-16s | %-14s | %-16s | %.2fx\n",
			r.difficulty,
			r.singleAvgTime.Round(time.Millisecond),
			formatHashRate(r.singleAvgHashes),
			r.concAvgTime.Round(time.Millisecond),
			formatHashRate(r.concAvgHashes),
			speedup,
		)
	}
}

func aggregate(results []minebench.Result, trials int) []row {
	sums := make(map[int]*row)

	for _, res := range results {
		r, ok := sums[res.Difficulty]
		if !ok {
			r = &row{difficulty: res.Difficulty}
			sums[res.Difficulty] = r
		}
		if res.Mode == "single" {
			r.singleAvgTime += res.Elapsed
			r.singleAvgHashes += res.HashesPerSecond()
		} else {
			r.concAvgTime += res.Elapsed
			r.concAvgHashes += res.HashesPerSecond()
		}
	}

	var out []row
	for _, bits := range minebench.Difficulties {
		r := sums[bits]
		r.singleAvgTime /= time.Duration(trials)
		r.concAvgTime /= time.Duration(trials)
		r.singleAvgHashes /= float64(trials)
		r.concAvgHashes /= float64(trials)
		out = append(out, *r)
	}
	return out
}

func formatHashRate(hps float64) string {
	switch {
	case hps >= 1_000_000:
		return fmt.Sprintf("%.2f MH/s", hps/1_000_000)
	case hps >= 1_000:
		return fmt.Sprintf("%.2f kH/s", hps/1_000)
	default:
		return fmt.Sprintf("%.0f H/s", hps)
	}
}
```

Running `go run ./cmd/minebench -trials 5` on an 8-core laptop produces output similar to this:

```
minebench: mining at difficulties [16 18 20 22], 5 trial(s) each, single-threaded vs. concurrent

Bits   | Single Time    | Single H/s       | Concur. Time   | Concur. H/s      | Speedup
-------|----------------|------------------|----------------|------------------|---------
16     | 68ms           | 963.41 kH/s      | 11ms           | 5.82 MH/s        | 6.04x
18     | 267ms          | 981.76 kH/s      | 39ms           | 6.61 MH/s        | 6.73x
20     | 1.062s         | 987.64 kH/s      | 152ms          | 6.85 MH/s        | 6.94x
22     | 4.31s          | 990.15 kH/s      | 596ms          | 7.03 MH/s        | 7.10x
```

---

## 7. Reading Your Own Results

A few things worth checking whenever you run `minebench` for real, on your own machine, rather than trusting the sample table above:

- **Single-threaded H/s should be roughly flat across every difficulty row.** It's measuring the same underlying "how many hashes per second can one core compute," which doesn't actually depend on the target — only total *time* changes with difficulty, exactly as Section 2 predicted. If your single-threaded H/s column varies wildly between rows, that's a sign of measurement noise (too few trials — try raising `-trials`) rather than a real difference in per-hash cost.
- **Concurrent H/s should be roughly `numCPU()` times the single-threaded rate, minus some overhead.** If your machine has 8 logical cores and concurrent H/s is only 2-3x single-threaded, something other than raw core count is limiting you — check `runtime.NumCPU()` actually reports what you expect (some cloud VMs and containers report fewer usable cores than the physical machine has), and check nothing else CPU-intensive is running at the same time.
- **Speedup should trend slightly upward as difficulty rises**, as the sample table shows (6.04x at 16 bits, up to 7.10x at 22 bits). This is a real, if subtle, effect: at very low difficulty, a run finishes so quickly that goroutine startup and scheduling overhead (Chapter 27, Section 8) make up a proportionally larger slice of the total time; at higher difficulty, that fixed overhead is amortized over a much longer, steadier run, letting the concurrent advantage show through more cleanly.

---

## Summary

- `minebench` measures mining performance across several difficulty levels, in both single-threaded and concurrent mode, reporting time, hashes per second, and speedup as one table.
- **Hashes per second** (`attempts / elapsedTime`) is a fairer comparison unit than raw wall-clock time alone, because it isolates *rate of work* from the *amount* of expected work, which by design grows with difficulty (Chapter 24).
- Neither `Run()` nor `RunConcurrent()` reports an attempt count by default, so `minebench` adds counted variants — trivial for the single-threaded case (the winning nonce itself is the count), requiring an `atomic.AddUint64`-based shared counter for the concurrent case, batched every 1024 iterations to keep the atomic operation's overhead small.
- Running multiple trials per difficulty and averaging their hash rates smooths out the natural per-run randomness that Chapter 24's geometric distribution guarantees any single mining run will have.
- The complete `minebench` CLI mines a batch of trials, aggregates them by difficulty, and prints a formatted table with human-readable hash-rate units (H/s, kH/s, MH/s).
- A sample run showed roughly 6-7x speedup on 8 cores, trending slightly upward at higher difficulty as fixed goroutine-startup overhead gets amortized over a longer run — concrete numbers backing up Chapter 27, Section 8's Amdahl's Law discussion.

---

## Exercises

### Easy

1. Run `minebench` on your own machine with `-trials 10` and report your own table. Is your single-threaded H/s roughly flat across difficulty levels, as Section 7 predicts? If not, what might explain the variation?
2. `Difficulties` is currently a package-level `var`, hardcoded to `{16, 18, 20, 22}`. Add a `-difficulties` flag to `main.go` (accepting a comma-separated list like `"14,16,18,20,22,24"`) that overrides it, and re-run with a wider range.
3. Explain, in your own words, why `RunCounted`'s attempt count (`nonce + 1`) would be WRONG if applied naively to a result from `RunConcurrentCounted` instead — i.e., why the concurrent case genuinely needs the atomic counter rather than reusing the single-threaded trick.

### Medium

4. Add a `-workers` flag to `main.go` that lets a user override how many goroutines `runConcurrentTrial` uses (instead of always passing `0` for "every core"), and extend the printed table with an extra column showing results at both `numCPU()` and `numCPU()/2` workers, to see how speedup changes with fewer workers than cores available.
5. `aggregate` currently averages hash rates across trials by summing `HashesPerSecond()` values and dividing by trial count. Discuss whether averaging the *rates* this way gives the same answer as computing one combined rate from total attempts divided by total elapsed time across all trials, and if they can differ, explain under what conditions and which one you'd trust more.
6. Extend `minebench` with a `-csv` flag that writes every individual `Result` (not just the averaged table) to a CSV file, one row per trial, suitable for loading into a spreadsheet for your own graphing. Include a header row with column names.

### Hard

7. Benchmark `minebench` itself running on a machine with hyperthreading/SMT enabled (most modern desktop and laptop CPUs), comparing `runtime.NumCPU()` workers against exactly half that many. Report whether using every *logical* core (including hyperthreads) actually outperforms using only as many workers as *physical* cores, and explain what you observe in terms of shared execution resources between hyperthreads on the same physical core.
8. Add a third mining mode to `minebench` — a version using exactly 2 workers, regardless of `runtime.NumCPU()` — and produce a table showing speedup at 1, 2, 4, and every-available-core worker counts, at a single fixed high difficulty. Plot (a text table is fine) speedup versus worker count and discuss where the curve starts noticeably flattening out, connecting your answer back to Chapter 27, Section 8's Amdahl's Law discussion.
9. `minebench`'s trials currently mine the exact same, fixed placeholder block every time (only the nonce search differs). Real mining also involves the (comparatively cheap, but nonzero) cost of building a fresh block from pending transactions each round. Extend the harness to optionally include realistic block-building overhead (a `core.NewBlock` call with a plausible number of dummy transactions) in each trial's timing, and discuss whether this changes the *relative* speedup between single-threaded and concurrent mining, or only the *absolute* numbers.
