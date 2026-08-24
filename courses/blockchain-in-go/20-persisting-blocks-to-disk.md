# Chapter 20: Persisting Blocks to Disk

Every chain built in Chapters 17 through 19 lives entirely in RAM. The moment the program exits — on purpose, by crashing, or because you closed your laptop — every block, every hash, every carefully-verified link vanishes with it. A blockchain nobody can reopen tomorrow is not much of a ledger. This chapter gives GoChain its first real storage layer: a simple append-only file that saves every block as it is created, and rebuilds the exact same in-memory chain from that file the next time the program starts.

## Table of Contents

1. [Why an In-Memory Chain Disappears](#1-why-an-in-memory-chain-disappears)
2. [Choosing gob for GoChain's First Storage Format](#2-choosing-gob-for-gochains-first-storage-format)
3. [Designing an Append-Only Block File](#3-designing-an-append-only-block-file)
4. [AppendBlockToFile — Writing One Block to Disk](#4-appendblocktofile--writing-one-block-to-disk)
5. [LoadBlockchain — Rebuilding the Chain on Startup](#5-loadblockchain--rebuilding-the-chain-on-startup)
6. [OpenBlockchain — Tying It Together](#6-openblockchain--tying-it-together)
7. [Trying It — Save, Restart, Reload](#7-trying-it--save-restart-reload)
8. [Limitations, Named Honestly](#8-limitations-named-honestly)
9. [Testing Persistence](#9-testing-persistence)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why an In-Memory Chain Disappears

Everything `core.Blockchain` has done so far lives inside a Go slice, `bc.blocks`, sitting in your program's RAM. RAM is fast, but it is also **volatile** — a term borrowed from electronics meaning its contents are lost the instant power stops flowing to it. Close your terminal, let your program `panic`, or simply reach the end of `main()`, and Go's runtime frees that memory for reuse by something else entirely. There is no "chain" left anywhere on your computer, only whatever you happened to print to the screen before it disappeared.

Think of writing an important letter directly onto a whiteboard instead of paper. It looks perfectly readable while it's up there, and you can revise it freely — but the moment someone wipes the board (or the office simply closes for the night and the cleaning crew comes through), it is gone completely, with no way to recover so much as the sentence you wrote first. Paper — or in this case, a file on disk — survives exactly the kind of "the room got cleaned" event a whiteboard does not.

**New term — persistence:** data that survives after the program that created it stops running, typically by being written to disk (or some other non-volatile storage) rather than kept only in RAM. Every chapter before this one built a blockchain with zero persistence: real, correct, tamper-evident — and gone forever the moment `main()` returns.

## 2. Choosing gob for GoChain's First Storage Format

Chapter 07 compared three ways to turn Go data into bytes: `encoding/json` (readable, larger), `encoding/gob` (Go-native, fast, and — critically for this chapter — able to handle a struct's *entire* shape, including nested slices and pointers, with almost no code), and hand-rolled binary encoding (smallest, but the most work to write correctly). Chapter 07 picked `gob` for GoChain's early chapters, and this is exactly where that choice pays off.

It is worth being precise about *why* `gob` is safe to reach for here, when Chapter 09, Section 6 spent real effort warning against a generic encoder for **hashing**. The two situations have different requirements. Hashing needs byte-for-byte **canonical** output — the same logical value must always produce identical bytes, because two different-looking but logically-equal byte sequences would break the hash comparisons Chapter 19 depends on. Chapter 17's `Block.Serialize()` still handles that job by hand, exactly as before — nothing about it changes in this chapter. Disk storage has a different, easier requirement: it only needs to be **round-trippable** — whatever bytes you write, reading them back must reconstruct an equal Go value. `gob` is built for exactly that, and does it correctly for every field on `Block`, `Transaction`, `TxInput`, and `TxOutput`, since every one of those fields is exported (capitalized), which is `gob`'s only real requirement.

```
   Hashing (Chapter 17's Serialize)        Disk storage (this chapter)

   Needs: canonical bytes                  Needs: round-trippable bytes
   ("same logical value always              ("write it, read it back,
   produces the SAME bytes")                get an equal value")

   Hand-written, field by field             gob.Encode / gob.Decode,
                                             the whole struct at once
```

## 3. Designing an Append-Only Block File

An **append-only file** is exactly what it sounds like: new data is always added to the *end* of the file, and existing bytes already written are never modified or moved. This is a deliberately simple, restrictive design, and that simplicity is the point — it maps directly onto how a blockchain actually grows (new blocks are added to the end; old blocks are never edited, per Chapter 19's tamper-evidence rules), and it sidesteps an entire category of bugs that come from editing bytes in the middle of an existing file.

The one design problem an append-only file has to solve is **framing**: `gob.Encode`, called once per block, produces one chunk of bytes, but multiple blocks appended back to back are just one long, undifferentiated stream of bytes on disk. Something has to mark where one block's encoded bytes end and the next one's begin — exactly the field-boundary problem Chapter 09, Section 7 and Chapter 17, Section 4 already solved for individual fields, now showing up one level higher, at the level of whole records in a file.

The fix is the same one you already know: a **length prefix**. Every record on disk is written as a 4-byte big-endian integer (the length of the encoded block that follows) immediately followed by that many bytes of `gob`-encoded block data.

```
  blocks.dat, growing one record at a time:

  [len0: 4 bytes][gob(block0): len0 bytes]
  [len1: 4 bytes][gob(block1): len1 bytes]
  [len2: 4 bytes][gob(block2): len2 bytes]
  ...

  Reading back: read 4 bytes -> know exactly how many more bytes
  make up the next block -> read exactly that many -> decode ->
  repeat until the file runs out of bytes.
```

This is precisely the same idea Chapter 44 later uses for TCP message framing — reading a stream of bytes off the wire has the exact same "where does one message end and the next begin" problem a growing file does, and the exact same length-prefix fix solves both.

## 4. AppendBlockToFile — Writing One Block to Disk

```go
// core/storage.go
package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
)

// AppendBlockToFile appends one block's gob-encoded, length-prefixed
// record to the file at path, creating the file if it does not exist
// yet. This is GoChain's entire storage engine for this volume --
// Chapter 53 explains, in detail, exactly why that will not remain
// true for long.
func AppendBlockToFile(path string, b *Block) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening block file %s: %w", path, err)
	}
	defer f.Close()

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(b); err != nil {
		return fmt.Errorf("gob-encoding block %d: %w", b.Height, err)
	}

	var lenPrefix [4]byte
	binary.BigEndian.PutUint32(lenPrefix[:], uint32(encoded.Len()))

	if _, err := f.Write(lenPrefix[:]); err != nil {
		return fmt.Errorf("writing block %d length prefix: %w", b.Height, err)
	}
	if _, err := f.Write(encoded.Bytes()); err != nil {
		return fmt.Errorf("writing block %d data: %w", b.Height, err)
	}

	return nil
}
```

Three `os.OpenFile` flags matter here, each doing exactly one job: `os.O_APPEND` guarantees every write lands at the current end of the file, even if some other part of the program (or another process) also has it open; `os.O_CREATE` makes the very first call, before the file exists at all, succeed instead of erroring out; `os.O_WRONLY` opens it write-only, since this function never needs to read. `defer f.Close()` ensures the file handle is released as soon as the function returns, on every code path, success or error alike — a pattern Chapter 07 already introduced for resource cleanup in Go.

Note that `Encode` writes to an in-memory `bytes.Buffer` first, not directly to the file. This is deliberate: it lets `AppendBlockToFile` know the encoded block's exact length *before* writing anything to disk at all, so the length prefix can be computed correctly and written first, exactly matching the format Section 3 designed.

## 5. LoadBlockchain — Rebuilding the Chain on Startup

Reading back reverses Section 4's process exactly: read 4 bytes, learn how many more bytes make up the next block, read exactly that many, decode, and repeat until the file has nothing left to give.

```go
// core/storage.go
import "io"

// LoadBlockchain rebuilds a Blockchain entirely from the records
// stored at path, in the order they were written. Because blocks are
// only ever appended in increasing height order (Chapter 18, Section
// 6; Chapter 19, Section 4), reading the file front to back
// reconstructs the exact same chain that was saved.
func LoadBlockchain(path string) (*Blockchain, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening block file %s: %w", path, err)
	}
	defer f.Close()

	bc := &Blockchain{}

	for {
		var lenPrefix [4]byte
		if _, err := io.ReadFull(f, lenPrefix[:]); err != nil {
			if err == io.EOF {
				break // clean end of file -- every record read successfully
			}
			return nil, fmt.Errorf("reading record length: %w", err)
		}

		recordLen := binary.BigEndian.Uint32(lenPrefix[:])
		data := make([]byte, recordLen)
		if _, err := io.ReadFull(f, data); err != nil {
			return nil, fmt.Errorf("reading %d bytes of block data: %w", recordLen, err)
		}

		var b Block
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&b); err != nil {
			return nil, fmt.Errorf("gob-decoding block: %w", err)
		}

		bc.blocks = append(bc.blocks, &b)
		bc.tip = b.Hash
	}

	if len(bc.blocks) == 0 {
		return nil, fmt.Errorf("block file %s contains no blocks", path)
	}

	return bc, nil
}
```

`io.ReadFull` is the right tool here instead of a plain `f.Read`, because a single `Read` call on a file is not guaranteed to fill the entire buffer you gave it in one go — `io.ReadFull` keeps reading until either the buffer is completely full or an error (including `io.EOF`) occurs, which is exactly the guarantee this loop needs when reading a fixed-size length prefix. The loop's termination condition is worth sitting with: `io.ReadFull` returning `io.EOF` on the *first* byte of an attempted read means the file ended cleanly, right where a new record was expected to begin — precisely what a well-formed file looks like after its very last record.

## 6. OpenBlockchain — Tying It Together

A convenience function spares every caller from having to remember "check if the file exists, and build a fresh genesis block if not":

```go
// core/storage.go

// OpenBlockchain loads an existing chain from path, or, if no file
// exists there yet, creates a brand new chain (Chapter 18's
// NewBlockchain) and persists its genesis block immediately, so the
// file on disk and the in-memory chain never disagree about what
// exists.
func OpenBlockchain(path string) (*Blockchain, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		bc := NewBlockchain()
		if err := AppendBlockToFile(path, bc.LastBlock()); err != nil {
			return nil, fmt.Errorf("persisting genesis block: %w", err)
		}
		return bc, nil
	}

	return LoadBlockchain(path)
}
```

Persistence is not automatic on every `AddBlock` call — that method's signature, fixed since Chapter 18, takes no file path and writes to memory only. Instead, callers explicitly persist each block right after adding it:

```go
chain, err := core.OpenBlockchain("chain.dat")
if err != nil {
	log.Fatal(err)
}

tx := &core.Transaction{ID: []byte("tx-a"), Timestamp: time.Now().Unix()}
next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)

if err := chain.AddBlock(next); err != nil {
	log.Fatal(err)
}
if err := core.AppendBlockToFile("chain.dat", next); err != nil {
	log.Fatal(err)
}
```

This explicitness is deliberate, not an oversight: it makes the exact moment data hits disk visible in the calling code, which matters directly for Section 8's honest discussion of what this design does and does not protect against.

## 7. Trying It — Save, Restart, Reload

Here is the entire round trip, in one program you can actually run twice to see persistence work:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/you/gochain/core"
)

func main() {
	const path = "chain.dat"

	chain, err := core.OpenBlockchain(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded chain at height %d\n", chain.Height())

	tx := &core.Transaction{
		ID:        []byte(fmt.Sprintf("tx-at-%d", time.Now().Unix())),
		Timestamp: time.Now().Unix(),
	}
	next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)

	if err := chain.AddBlock(next); err != nil {
		log.Fatal(err)
	}
	if err := core.AppendBlockToFile(path, next); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added block %d. New height: %d\n", next.Height, chain.Height())
	_ = os.Args // placeholder; a real CLI parses flags here (Chapter 21)
}
```

Run it once: `Loaded chain at height 0` (a fresh genesis-only chain, since `chain.dat` did not exist yet), then `Added block 1. New height: 1`. Run it again, without deleting `chain.dat`: `Loaded chain at height 1` — the exact chain from the previous run, rebuilt entirely from disk, with `chain.Height()` correctly reporting 1 before a single new block is added. Run it a third time and you will see height 2, then height 3, and so on — the program has genuine memory now, across restarts, something nothing before this chapter could claim.

## 8. Limitations, Named Honestly

This storage layer is deliberately the simplest thing that could work, and it is worth being completely honest about three specific ways it falls short of what a real system needs — not as a confession of failure, but as the exact motivation Chapter 53 picks back up in detail when Volume 8 replaces this whole file with a real embedded database.

**Crash-safety.** `AppendBlockToFile` performs two separate `Write` calls — one for the length prefix, one for the encoded block. If the process crashes, or the machine loses power, between those two writes (or partway through the second one), `chain.dat` ends up with a length prefix promising N more bytes that are not actually all present. `LoadBlockchain`'s `io.ReadFull` call for that final, truncated record will fail with an error that is neither a clean success nor a clean `io.EOF` — the file is left in a state this code was never designed to recover from gracefully. A real database uses techniques like write-ahead logging and atomic commits specifically to guarantee a crash mid-write never leaves data in an unreadable, half-written state; this chapter's file format has none of that.

**Concurrent writes.** `os.O_APPEND` guarantees each individual `Write` call lands atomically at the file's current end — but nothing stops two goroutines (or two separate processes) from both opening `chain.dat` and calling `AppendBlockToFile` at nearly the same moment, potentially interleaving their length-prefix and data writes in a way that corrupts both records. Chapter 05's discussion of race conditions applies directly here, just at the level of a file instead of a shared Go variable — and nothing in this chapter's code takes a lock of any kind to prevent it.

**Scan performance.** `LoadBlockchain` reads the *entire* file, decoding every single block, every time the program starts — and there is no way to jump straight to "the block at height 500,000" without first reading and decoding blocks 0 through 499,999 to find it. For a chain with a handful of blocks this is instant; for a real chain with millions of blocks, this becomes a slow, linear scan on every single startup, and there is no efficient way to look up one specific block, one specific transaction, or one address's balance without a full scan every time.

```
                    THIS CHAPTER              REAL DATABASE (Volume 8)

  Crash mid-write     Corrupted trailing        Write-ahead log, atomic
                       record, unrecoverable      commits
  Concurrent writers  Not safe, no locking       Handled by the engine
  Find block N        O(n) -- scan from the      O(log n) or O(1) via
                       start every time            an index
```

None of this makes the flat-file approach *wrong* for this stage of the course — it is genuinely, honestly the simplest thing that demonstrates real persistence, and Chapter 53 reproduces exactly these three problems hands-on, against this very code, before Chapter 54 introduces BoltDB as the fix. Naming a simplification's limits clearly, rather than pretending it has none, is itself a habit worth carrying into every real system you build after this course.

## 9. Testing Persistence

The property that matters most here is a round trip: build a chain, save it, load it back from a different `Blockchain` value entirely, and confirm the two are equivalent.

```go
// core/storage_test.go
package core

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPersistence_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.dat")

	original, err := OpenBlockchain(path)
	if err != nil {
		t.Fatalf("OpenBlockchain: %v", err)
	}

	for i := 1; i <= 3; i++ {
		next := NewBlock([]*Transaction{testTx("tx")}, original.Tip(), original.Height()+1)
		if err := original.AddBlock(next); err != nil {
			t.Fatalf("AddBlock: %v", err)
		}
		if err := AppendBlockToFile(path, next); err != nil {
			t.Fatalf("AppendBlockToFile: %v", err)
		}
	}

	reloaded, err := LoadBlockchain(path)
	if err != nil {
		t.Fatalf("LoadBlockchain: %v", err)
	}

	if reloaded.Height() != original.Height() {
		t.Fatalf("expected height %d after reload, got %d", original.Height(), reloaded.Height())
	}
	if !bytes.Equal(reloaded.Tip(), original.Tip()) {
		t.Fatal("expected reloaded tip to match original tip")
	}
	if err := reloaded.ValidateChain(); err != nil {
		t.Fatalf("expected reloaded chain to validate, got: %v", err)
	}
}

func TestOpenBlockchain_CreatesFileIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist-yet.dat")

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to not exist before OpenBlockchain")
	}

	bc, err := OpenBlockchain(path)
	if err != nil {
		t.Fatalf("OpenBlockchain: %v", err)
	}
	if bc.Height() != 0 {
		t.Fatalf("expected a fresh genesis-only chain, got height %d", bc.Height())
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("expected OpenBlockchain to have created the file")
	}
}
```

`t.TempDir()` gives each test its own throwaway directory, automatically cleaned up when the test finishes — the right tool whenever a test needs to touch the real filesystem without leaving files scattered around your repository. `TestPersistence_RoundTrip` is the test that matters most in this chapter: it is the automated proof that everything Section 7 demonstrated by running a program twice by hand also holds true as real, repeatable, `go test`-verified behavior — including, satisfyingly, that the reloaded chain still passes Chapter 19's `ValidateChain()` untouched.

---

## Summary

- An in-memory `Blockchain` lives only in volatile RAM and disappears completely when the program exits — **persistence** means writing data somewhere, like disk, that survives past the process that created it.
- GoChain's first storage format uses `encoding/gob` to encode whole blocks, because disk storage only needs round-trippable bytes, a weaker (and easier) requirement than the canonical bytes Chapter 17's `Serialize()` still provides for hashing.
- An **append-only file** only ever adds new bytes at the end, mirroring how a real blockchain grows and avoiding an entire class of "editing bytes in the middle of a file" bugs.
- Each on-disk record is a 4-byte length prefix followed by that many bytes of gob-encoded block data — the same length-prefixing idea from Chapters 09 and 17, now applied to whole records instead of individual fields.
- `AppendBlockToFile` writes one block's record to the end of a file; `LoadBlockchain` reads records back in order, using `io.ReadFull` and the length prefix to know exactly where each record ends.
- `OpenBlockchain` ties both together: load an existing file, or create and persist a fresh genesis block if none exists yet.
- This design has three honestly-named limitations — no crash-safety across a partial write, no protection against concurrent writers, and a full linear scan required on every startup — that Chapter 53 revisits in depth before Chapter 54 replaces this whole approach with a real embedded database.

---

## Exercises

### Easy

1. Run the Section 7 program three times in a row against a fresh `chain.dat`, recording the printed height before and after each run. Then delete `chain.dat` and run it once more, confirming the chain starts over from height 0.
2. Explain, in your own words, why `AppendBlockToFile` writes the encoded block into an in-memory `bytes.Buffer` first, rather than encoding directly to the open file handle.
3. `LoadBlockchain` returns an error if the file contains zero blocks. Explain why this should never actually happen in normal use, given how `OpenBlockchain` is written, and describe one abnormal scenario (outside of normal `OpenBlockchain`/`AppendBlockToFile` use) that could still produce an empty file.

### Medium

4. Write a small program that calls `AppendBlockToFile` to write a valid block, then manually truncates the last few bytes off the resulting file (using `os.Truncate` or by editing the file directly) to simulate a crash mid-write. Call `LoadBlockchain` on the truncated file and report the exact error you get. Does any partial chain data survive, or does the whole load fail?
5. Extend `LoadBlockchain` with a "best effort" mode: instead of returning an error the moment it hits a corrupted trailing record, it should return every block successfully read *before* the corruption, plus a count of how many bytes of unreadable data were skipped at the end. Write a test proving this against a deliberately truncated file.
6. Benchmark `LoadBlockchain` against block files containing 100, 10,000, and 1,000,000 blocks (you can generate placeholder blocks with minimal transactions for this). Report how load time scales with block count, and relate your results to Section 8's "scan performance" limitation in concrete numbers rather than just the word "slow."

### Hard

7. Implement a simple file lock (using a `.lock` file whose existence signals "another process is writing," created and removed around each `AppendBlockToFile` call) to address Section 8's concurrent-writes limitation. Write a test using two goroutines calling `AppendBlockToFile` on the same path concurrently, with and without your lock, and demonstrate the difference in outcome.
8. Design (in a written proposal, 300-500 words) a checksum scheme for each record in `blocks.dat` — for example, storing a SHA-256 hash of each record's raw bytes alongside its length prefix — that would let `LoadBlockchain` detect (not necessarily recover from) silent bit-level corruption in an old record, not just a truncated final one. Explain what new failure mode this catches that the length-prefix format alone does not.
9. Research BoltDB's or SQLite's actual on-disk crash-safety mechanism (write-ahead logging, in particular) at a conceptual level, and write a comparison (300-400 words) between it and this chapter's append-only file, specifically addressing: what happens, in each system, if the process is killed (`kill -9`, not a graceful shutdown) in the middle of writing one record, and why your answer differs between the two systems.
