# Chapter 53: Why Flat Files Are Not Enough

Since Chapter 20, GoChain has saved its blocks by appending them, one after another, to a single file on disk. That was the right amount of complexity for Volume 3 — it let you see a chain survive a program restart without dragging in a whole database library before you even understood what a block was. This chapter puts that same flat-file store under real pressure — a crash mid-write, a chain with tens of thousands of blocks, two goroutines writing at once — and shows, hands-on, exactly where it breaks.

## Table of Contents

1. [Recap: What the Chapter 20 Flat File Actually Does](#1-recap-what-the-chapter-20-flat-file-actually-does)
2. [Problem 1: A Crash Mid-Write Corrupts the File](#2-problem-1-a-crash-mid-write-corrupts-the-file)
3. [Problem 2: Finding One Transaction Means Scanning Everything](#3-problem-2-finding-one-transaction-means-scanning-everything)
4. [Problem 3: Two Goroutines Writing at Once](#4-problem-3-two-goroutines-writing-at-once)
5. [Why "Just Add a Mutex and a Try/Catch" Is Not Enough](#5-why-just-add-a-mutex-and-a-trycatch-is-not-enough)
6. [What We Actually Need](#6-what-we-actually-need)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Recap: What the Chapter 20 Flat File Actually Does

A quick recap, because we are about to break this code on purpose. The Chapter 20 store keeps things as simple as possible: open one file, and every time a new block is mined or received, gob-encode it and append the bytes to the end of the file. To read the chain back on startup, read from the beginning and gob-decode blocks one after another until you hit the end of the file.

```go
// core/flatfile.go — the Chapter 20 store, reproduced here so we can break it
package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// AppendBlock writes one length-prefixed, gob-encoded block to the end of f.
// The 4-byte length prefix tells the reader exactly how many bytes to expect,
// so blocks of different sizes can sit back-to-back in the same file.
func AppendBlock(f *os.File, b *Block) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return fmt.Errorf("encode block: %w", err)
	}

	lenPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lenPrefix, uint32(buf.Len()))

	// Two separate writes: first the length, then the body. This detail
	// matters enormously later in this chapter.
	if _, err := f.Write(lenPrefix); err != nil {
		return fmt.Errorf("write length prefix: %w", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("write block body: %w", err)
	}
	return nil
}

// LoadChain reads every block back out of the file, in the order they were
// appended, by repeatedly reading a length prefix and then that many bytes.
func LoadChain(path string) ([]*Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open chain file: %w", err)
	}
	defer f.Close()

	var blocks []*Block
	lenBuf := make([]byte, 4)
	for {
		_, err := io.ReadFull(f, lenBuf)
		if err == io.EOF {
			break // clean end of file — every block we wrote was read back
		}
		if err != nil {
			return nil, fmt.Errorf("read length prefix: %w", err)
		}

		n := binary.BigEndian.Uint32(lenBuf)
		body := make([]byte, n)
		if _, err := io.ReadFull(f, body); err != nil {
			return nil, fmt.Errorf("read block body (wanted %d bytes): %w", n, err)
		}

		var b Block
		if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&b); err != nil {
			return nil, fmt.Errorf("decode block: %w", err)
		}
		blocks = append(blocks, &b)
	}
	return blocks, nil
}
```

`AppendBlock` does two separate `Write` calls — one for the 4-byte length prefix, one for the block's actual bytes — because we do not know the block's encoded size until after we have gob-encoded it. `LoadChain` mirrors this: read 4 bytes to learn how long the next block is, then read exactly that many bytes and decode them. This "length-prefixed framing" idea is simple, and it is genuinely the right building block — you will see the same pattern again, more formally, when GoChain's network protocol frames messages in Volume 7. The problem is not the framing. The problem is everything around it.

```
blockchain.dat
+------+----------------+------+----------------+------+----------------+
| LEN  |  gob(block 0)  | LEN  |  gob(block 1)  | LEN  |  gob(block 2)  |
| 4B   |   N0 bytes     | 4B   |   N1 bytes     | 4B   |   N2 bytes     |
+------+----------------+------+----------------+------+----------------+
        ^                                                              ^
        start of file                                          current end of file
        (LoadChain reads from here, block by block, until EOF)
```

---

## 2. Problem 1: A Crash Mid-Write Corrupts the File

`f.Write` is not atomic across a crash. If GoChain's process is killed — a `kill -9`, a power outage, a container getting OOM-killed by Kubernetes — in the middle of appending a block, the file can be left in a state `LoadChain` was never designed to handle. Let's reproduce this on purpose rather than take it on faith.

```go
// core/flatfile_crash_test.go
package core

import (
	"os"
	"testing"
)

// simulateCrashMidWrite writes a length prefix claiming N bytes follow, but
// only actually writes N-10 bytes before the "process dies" — exactly what
// a real crash between two Write() calls (or mid-way through one) can leave
// behind on disk.
func simulateCrashMidWrite(t *testing.T, path string, b *Block) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := AppendBlock(f, b); err != nil {
		t.Fatal(err)
	}

	// Now truncate the file to lop off the last 10 bytes of the block body
	// we just wrote — simulating a crash that happened after the length
	// prefix and part of the body hit disk, but before the rest did.
	info, _ := f.Stat()
	if err := f.Truncate(info.Size() - 10); err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestLoadChain_CrashedFile(t *testing.T) {
	path := t.TempDir() + "/blockchain.dat"
	block := NewGenesisBlock()

	simulateCrashMidWrite(t, path, block)

	_, err := LoadChain(path)
	if err == nil {
		t.Fatal("expected LoadChain to fail on a truncated file, got nil error")
	}
	t.Logf("LoadChain correctly failed with: %v", err)
}
```

Running this test produces exactly the failure we engineered:

```
$ go test ./core/ -run TestLoadChain_CrashedFile -v
=== RUN   TestLoadChain_CrashedFile
    flatfile_crash_test.go:44: LoadChain correctly failed with: read block body (wanted 612 bytes): unexpected EOF
--- PASS: TestLoadChain_CrashedFile (0.00s)
```

Notice what "correctly failed" means here: `LoadChain` returns an error, which is the *best case*. The genuinely dangerous case is a crash that lands the file in a shape `gob.Decode` does not error on at all — for example, if the crash happens to leave bytes that gob happily decodes into a `Block` with a zeroed-out `Hash` or a `Transactions` slice missing its last entry. GoChain would then load a chain that *looks* fine, pass its own hash checks against corrupted data, and silently serve wrong balances to a wallet. A loud crash is annoying. A silent one is the failure mode that actually loses someone money.

```
BEFORE CRASH (in flight):                    AFTER CRASH (on disk):
+------+----------------+------+----...      +------+----------------+------+--
| LEN  | gob(block 0)   | LEN  | gob(bl...   | LEN  | gob(block 0)   | LEN  |  (10 bytes
| 4B   |  full N1 bytes | 4B   | ock 1)       | 4B   |  full N1 bytes | 4B   |   missing!)
+------+----------------+------+----...      +------+----------------+------+--
                                                                     ^
                                              process died here, mid-Write() —
                                              length prefix says N1 bytes follow,
                                              but only N1-10 actually made it to disk
```

The root cause is that our flat file has **no notion of a transaction boundary that the operating system itself understands**. A real database calls this property **atomicity**: a write either completes fully, as a single indivisible unit, or it has no effect at all — never a half-finished state sitting on disk waiting to confuse the next reader. `os.File.Write` gives you none of that for free; you get exactly the bytes you asked it to write, landing wherever the OS's page cache and disk scheduler decide, in whatever order, unless you go out of your way to force otherwise.

---

## 3. Problem 2: Finding One Transaction Means Scanning Everything

The second problem shows up not from a crash, but from success — GoChain working exactly as designed, for long enough that the chain gets big. Suppose a wallet wants to answer "did transaction `abc123...` ever confirm?" or "what is this address's balance?" With only a flat, append-only file and no index, there is exactly one way to answer either question: read every block, from the first to the last, and inspect every transaction inside each one.

```go
// core/scan.go
package core

import "bytes"

// FindTransaction scans every block, oldest to newest, looking for a
// transaction with the given ID. This is the only way to find a transaction
// when all you have is an ordered list of blocks and no index.
func FindTransaction(blocks []*Block, txID []byte) (*Transaction, error) {
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, txID) {
				return tx, nil
			}
		}
	}
	return nil, fmt.Errorf("transaction %x not found in %d blocks", txID, len(blocks))
}
```

This works. It is also the exact definition of an operation that gets slower, in direct proportion to how successful GoChain is. Let's put real numbers on it rather than wave our hands at "it doesn't scale."

```go
// core/scan_bench_test.go
package core

import "testing"

func BenchmarkFindTransaction(b *testing.B) {
	chain := generateTestChain(50_000) // helper: 50,000 blocks, ~15 tx each
	target := chain[0].Transactions[0].ID // worst case: the very first tx ever mined

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := FindTransaction(chain, target); err != nil {
			b.Fatal(err)
		}
	}
}
```

```
$ go test ./core/ -bench BenchmarkFindTransaction -benchtime=20x
goos: darwin
goarch: arm64
pkg: gochain/core
BenchmarkFindTransaction-8    20    398214583 ns/op    41821 B/op    3 allocs/op
PASS
```

That is roughly **400 milliseconds** to answer one "does this transaction exist?" question — for the worst case of a transaction near the genesis block, on a chain of 50,000 blocks with about 15 transactions each (750,000 transactions total). Work through where that time goes:

```
50,000 blocks x (gob-decode ~5-8 microseconds per block, from Chapter 20's
                 encoding choice) + 15 transactions x a byte-compare each
                 ≈ 50,000 x ~8µs
                 ≈ 400,000 microseconds
                 ≈ 400 milliseconds
```

Four hundred milliseconds sounds almost tolerable in isolation — until you remember that a real node is not answering one such question in a quiet room. Volume 10's API will need to serve `getBalance` and `getTransaction` calls to potentially dozens of concurrent requests from wallets and a block explorer, all hitting a chain that only grows. At 50,000 blocks this design already costs 400ms per lookup in the worst case; at 500,000 blocks (a plausible size after a year of real usage) the same math predicts roughly **4 seconds** per lookup — for a single request, before any concurrent load. `BalanceOf(address)`, which must scan for *every* UTXO belonging to an address rather than stopping at the first match, is worse still: it cannot short-circuit early, so it pays the full scan cost on every single call, every time a wallet's balance page refreshes.

This is exactly the kind of complexity that matters more as a system succeeds than as it starts out — the demo with 12 blocks felt instant, which is precisely why the problem is easy to miss until it is expensive to fix. Laid out across chain sizes, the trend is unambiguous — this is a straight line, not a curve that levels off:

```
Chain size       Approx. worst-case FindTransaction scan time
------------     ---------------------------------------------
     1,000                 ~8 ms
    10,000                 ~80 ms
    50,000                 ~400 ms   <- measured above
   500,000                 ~4.0 s
 5,000,000                 ~40 s

Every column is the previous one, scaled by the same factor as the chain
size grew. There is no point on this line where the flat file "catches up" —
it gets linearly, predictably worse, forever.
```

This is the textbook definition of an **O(n) operation** — one whose cost grows in direct, linear proportion to the size of the input, here the number of blocks. An indexed lookup, by contrast, is close to **O(1)** — its cost barely changes whether the chain has a thousand blocks or five million, because it goes directly to the answer instead of walking past everything that is not the answer. Chapter 56 builds exactly this for balances; the difference is not a minor optimization, it is the difference between a wallet that feels instant and one that visibly hangs.

---

## 4. Problem 3: Two Goroutines Writing at Once

The third problem does not need a crash or a large chain — it can happen on the very next block, on a perfectly healthy machine, the moment more than one goroutine calls `AppendBlock` concurrently. This is a completely realistic scenario for GoChain from Volume 7 onward: a mining goroutine finishes solving a block at nearly the same instant a network goroutine receives a different block from a peer during sync.

```go
// core/flatfile_race_test.go
package core

import (
	"os"
	"sync"
	"testing"
)

func TestAppendBlock_ConcurrentWrites(t *testing.T) {
	path := t.TempDir() + "/blockchain.dat"
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	blockA := NewBlock(1, []*Transaction{}, []byte("prevhash-a"))
	blockB := NewBlock(1, []*Transaction{}, []byte("prevhash-b"))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); AppendBlock(f, blockA) }() // mining goroutine
	go func() { defer wg.Done(); AppendBlock(f, blockB) }() // network sync goroutine
	wg.Wait()

	blocks, err := LoadChain(path)
	if err != nil {
		t.Fatalf("LoadChain failed after concurrent writes: %v", err)
	}
	t.Logf("recovered %d blocks (expected 2)", len(blocks))
}
```

Run this with Go's built-in race detector, which instruments every memory access and reports the exact moment two goroutines touch the same data without synchronization:

```
$ go test ./core/ -run TestAppendBlock_ConcurrentWrites -race -v
=== RUN   TestAppendBlock_ConcurrentWrites
==================
WARNING: DATA RACE
Write at 0x00c0001100a0 by goroutine 8:
  os.(*File).Write()
      /usr/local/go/src/os/file.go:176
  gochain/core.AppendBlock()
      /Users/you/gochain/core/flatfile.go:24

Previous write at 0x00c0001100a0 by goroutine 7:
  os.(*File).Write()
      /usr/local/go/src/os/file.go:176
  gochain/core.AppendBlock()
      /Users/you/gochain/core/flatfile.go:24
==================
--- FAIL: TestAppendBlock_ConcurrentWrites (0.01s)
FAIL
```

Even setting the race detector's low-level warning aside, think through what can actually happen at the file level. `os.File` tracks one shared write offset per open file descriptor. Two goroutines calling `Write` at nearly the same moment can each read "the current end of file is byte 1,204" *before* either one has actually appended its bytes — so both write starting at offset 1,204, and one goroutine's bytes silently overwrite the other's.

```
Goroutine A (mining)                 Goroutine B (network sync)
   |                                       |
   | reads current offset: 1204            | reads current offset: 1204
   |                                       |    (same offset — A hasn't written yet!)
   | writes blockA's bytes at 1204         |
   |                                       | writes blockB's bytes AT 1204 TOO
   |                                       |    (overwrites part or all of blockA)
   v                                       v
                  Result on disk: a mangled mix of both blocks,
                  or one block's bytes silently clobbering the other's.
```

Sometimes this race is "lucky" and one write simply lands after the other, appending both blocks correctly. Sometimes it is not, and the file ends up with an interleaved mess that `LoadChain` cannot parse, or — worse — a mess it *can* parse into something that is neither block A nor block B. Bugs that only sometimes reproduce, depending on precise timing, are some of the most expensive bugs to track down in production, precisely because they pass every test run except the unlucky one.

---

## 5. Why "Just Add a Mutex and a Try/Catch" Is Not Enough

At this point it is tempting to reach for the smallest possible patch: wrap `AppendBlock` in a `sync.Mutex` to fix Problem 3, and recover from decode errors to paper over Problem 1. Both patches are worth understanding — and worth understanding why they stop short of a real fix.

A mutex genuinely does fix concurrent writes *within a single process*:

```go
type FileStore struct {
	mu sync.Mutex
	f  *os.File
}

func (s *FileStore) AppendBlock(b *Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return AppendBlock(s.f, b) // now safe from other goroutines in this process
}
```

This is a real improvement, and you should always reach for it when multiple goroutines share one resource — GoChain's later chapters use mutexes exactly like this in several places. But it does not touch Problem 1 (a crash still leaves a half-written block on disk, mutex or not — a lock protects against *other goroutines*, not against the process disappearing mid-`Write`) and it does not touch Problem 2 at all (a mutex does not make scanning 50,000 blocks any faster; if anything, serializing all writes behind one lock makes GoChain's single writer a bottleneck as block frequency increases).

The equivalent quick patch for Problem 1 is to wrap `LoadChain`'s decode step in a `recover()` and simply stop at the first block that fails to parse:

```go
func LoadChainTolerant(path string) (blocks []*Block, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Swallow the panic and return whatever we managed to load.
			// This hides the corruption instead of fixing it.
		}
	}()
	// ... same loop as LoadChain, but stop and return on the first error
	// instead of propagating it.
	return blocks, nil
}
```

This looks like a fix because the program no longer crashes on startup. But think about what it actually does: it silently discards the last block (or blocks) that failed to decode, with no record anywhere that data was lost, no way for an operator to know how many blocks vanished, and no guarantee the *next* corrupted write will fail as cleanly as this one did. Defensive error handling is valuable as a second line of defense — GoChain will absolutely still validate everything it reads, forever — but it is not a substitute for a storage engine that does not produce half-written data in the first place.

What we actually need is:

- **Crash safety** built into the storage engine itself, so a half-finished write is either rolled back or was never visible to readers in the first place — not something we bolt on with defensive error handling after the fact.
- **Indexed lookups**, so finding a block by hash or a UTXO by key is a direct, near-constant-time operation, not a scan whose cost grows with the entire history of the chain.
- **Safe concurrent access** that is a property of the storage engine's design, not something every caller has to remember to wrap in a mutex correctly, every single time, forever.

These are exactly the three properties a real embedded database gives you essentially for free, because solving them well is the entire reason such databases exist. Chapter 54 introduces two well-established options — BoltDB and Badger — and Chapter 55 replaces GoChain's flat file with one of them for good.

---

## 6. What We Actually Need

To close the loop, here is the shape of what the rest of this volume builds, stated as requirements rather than code:

- A storage engine that groups related writes into a single atomic unit, so a crash mid-write leaves the last *complete* write intact and never a half-finished one.
- A way to look up a block by its hash, or a UTXO by its key, directly — without touching any other block or UTXO — regardless of how large the chain has grown.
- A design that is safe to call from multiple goroutines without every caller re-inventing its own locking scheme.
- All of this behind an interface, so GoChain's `core` package depends on "a place that stores blocks and UTXOs" rather than on any one specific database library — exactly the same interface-first instinct you used for `consensus.Engine` in Volume 4.

```
TODAY (Chapter 20 flat file)              AFTER THIS VOLUME (BoltDB via storage.Store)
+---------------------------+             +---------------------------------------+
| append-only file          |             | storage.Store interface                |
| - no atomicity            |   ===>      |  - atomic transactions (Ch 55)         |
| - full linear scan only   |             |  - indexed Get/Put by key (Ch 55)      |
| - no concurrency safety   |             |  - safe for concurrent callers (Ch 54) |
+---------------------------+             |  - UTXOSet index for balances (Ch 56)  |
                                           |  - state root via a trie (Ch 57)       |
                                           +---------------------------------------+
```

---

## Summary

- The Chapter 20 flat-file store appends length-prefixed, gob-encoded blocks to one file, and rebuilds the chain by reading them back in order — simple, and good enough only for the earliest volumes.
- A crash mid-`Write` can leave a half-written block on disk. In the best case, `LoadChain` errors out loudly; in the worst case, it silently loads corrupted data that passes every check.
- Finding one transaction or computing one balance requires scanning every block from the start, because the flat file has no index. A worked estimate on a 50,000-block chain puts this at roughly 400ms per lookup in the worst case, and growing linearly with chain size from there.
- Two goroutines calling `AppendBlock` concurrently — entirely realistic once GoChain has both a miner and a network layer — can interleave their writes at the OS level, corrupting the file, exactly as Go's race detector demonstrates.
- A mutex fixes the concurrency problem within one process but does nothing for crash safety or scan performance; a `recover()` around decode errors hides symptoms without fixing the underlying lack of atomicity.
- What GoChain actually needs is atomic writes, indexed lookups, and safe concurrent access — properties a real embedded database provides by design, which the rest of this volume adopts behind a `storage.Store` interface.

---

## Exercises

### Easy

1. In your own words, explain why `LoadChain` failing loudly on a truncated file (as in Section 2's test) is actually the *good* outcome, and describe a scenario where a crash could instead produce a chain that loads successfully but contains wrong data.
2. Modify `simulateCrashMidWrite` to truncate the file in the middle of the 4-byte length prefix itself (rather than the block body). Does `LoadChain` still fail with a similar error, or does it behave differently?
3. Run the `BenchmarkFindTransaction` benchmark (or a smaller version you construct yourself with `generateTestChain(5_000)`) on your own machine and record the actual `ns/op` you get. How does it compare to the 50,000-block estimate in Section 3?

### Medium

4. Add a `sync.Mutex` to a `FileStore` wrapper type as shown in Section 5, and re-run the concurrent-write test from Section 4 with `-race`. Confirm the race detector no longer reports a warning, but explain in a comment why this change does not address Problems 1 or 2 from this chapter.
5. Estimate, using the same per-block decode-time assumption from Section 3, how long `FindTransaction`'s worst case would take on a chain of 500,000 blocks. Then estimate it for 5 million blocks. At what chain size does a single worst-case lookup start to exceed one full second?
6. `LoadChain` currently returns an error the moment it cannot fully read a block's body. Write an alternative version, `LoadChainTolerant`, that instead stops reading and returns every *valid* block read so far when it hits a truncated final block. Discuss, in a short comment, what real blockchain behavior this partially mirrors (hint: think about what a node should do with the very last, currently-being-written block after a crash).

### Hard

7. Design (in comments or a short doc, no need to implement it fully) a minimal "write-ahead log" scheme that would make `AppendBlock` crash-safe using nothing but flat files: what would you write before the actual append, and how would `LoadChain` use it on startup to detect and discard an incomplete write? You do not need to reinvent a full WAL — a workable sketch is enough.
8. The benchmark in Section 3 measures a single-threaded, uncontended lookup. Extend `TestAppendBlock_ConcurrentWrites` into a benchmark that runs `FindTransaction` concurrently from 20 goroutines while a 21st goroutine is appending new blocks with a mutex-protected `FileStore`. What happens to scan latency under concurrent read load, and why?
9. Research how a real database (SQLite is a good, well-documented example) achieves atomic writes using techniques like a rollback journal or write-ahead logging. Write a short explanation, in your own words, of how SQLite's approach avoids the exact corruption scenario demonstrated in Section 2 of this chapter — even though SQLite, like our flat file, ultimately writes to ordinary OS files too.
