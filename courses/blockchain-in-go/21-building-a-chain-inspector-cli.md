# Chapter 21: Building a Chain Inspector CLI

Every chapter so far has inspected blocks by sprinkling `fmt.Println` calls through throwaway test programs. That works for one experiment, but it does not scale to the rest of this course — you will build, break, mine, and network dozens of chains between here and the final volume, and you need a real, reusable tool for looking inside one. This chapter builds `gochain inspect`: a proper command-line tool that prints every block in a chain readably, and a `--verify` flag that walks the whole thing checking every hash and every link using Chapter 19's `ValidateBlock`.

## Table of Contents

1. [Why a Debugging CLI Matters This Early](#1-why-a-debugging-cli-matters-this-early)
2. [Designing the inspect Command](#2-designing-the-inspect-command)
3. [Formatting a Block for Human Eyes](#3-formatting-a-block-for-human-eyes)
4. [Building the CLI With Go's flag Package](#4-building-the-cli-with-gos-flag-package)
5. [Wiring Up --verify](#5-wiring-up---verify)
6. [Running the Inspector — A Worked Session](#6-running-the-inspector--a-worked-session)
7. [Handling a Corrupted or Missing Chain Gracefully](#7-handling-a-corrupted-or-missing-chain-gracefully)
8. [Testing the CLI Logic](#8-testing-the-cli-logic)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why a Debugging CLI Matters This Early

Imagine a mechanic who could only ever look at a car engine while it was already running, with the hood welded shut — no dipstick, no dashboard lights, no way to pop the hood and actually look at any single part in isolation. Every diagnosis would be a guess made from vibration and sound alone. A **chain inspector** is GoChain's dipstick and dashboard combined: a tool that opens the hood, one block at a time, and reports exactly what is inside without you having to write a new throwaway Go program every single time you want to check something.

This matters earlier than it might seem to. From Volume 4 onward, chains get more complex — mining introduces real proof-of-work search loops, Volume 5 adds real signed transactions, Volume 7 has multiple independent nodes each holding their own copy of the chain. Every one of those chapters benefits from being able to type one command and see, in plain text, exactly what a chain currently contains and whether it is still internally consistent — rather than reaching for a debugger or adding print statements to already-working code. Building this tool now, while `Block` and `Blockchain` are still simple, means it is ready and trustworthy by the time the data it inspects gets genuinely complicated.

## 2. Designing the inspect Command

`gochain inspect` needs to do two related but separable jobs: **print** every block readably, and, only when asked, **verify** that the chain it just printed is actually valid. Splitting these into a plain listing and an opt-in check, rather than always verifying, matters because verification (Chapter 19's `ValidateBlock`, called once per block) is real, non-trivial work — recomputing a hash for every single block — and a quick "let me just look at what's in here" listing should not have to pay that cost every time.

```
gochain inspect [-db path] [-verify]

  -db string     path to the block file to inspect (default "chain.dat")
  -verify        walk the chain and validate every block and link
```

This is intentionally the simplest possible command surface: one subcommand, two flags. Chapter 74 replaces this ad-hoc, hand-rolled dispatch with the Cobra library once GoChain's CLI surface grows to a dozen commands across every volume — but for exactly two flags on exactly one subcommand, Go's standard `flag` package is more than enough, and reaching for a third-party library this early would add complexity this chapter has no use for yet.

## 3. Formatting a Block for Human Eyes

Chapter 16, Section 2's field-size table showed that a `Block`'s raw form — 32-byte hashes, a Unix timestamp integer — is not something a human should be expected to read directly. A dedicated formatting function turns one block into exactly the readable summary Chapter 16's own field-by-field diagram promised: height, timestamp, transaction count, hash, and previous hash.

```go
// cmd/gochain/inspect.go
package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/you/gochain/core"
)

// formatBlock renders b as a multi-line, human-readable summary.
// Returning a string rather than printing directly keeps this
// function trivially testable (Section 8) -- no need to capture
// stdout to check its output.
func formatBlock(b *core.Block) string {
	return fmt.Sprintf(
		"Block %d\n"+
			"  Timestamp:     %s\n"+
			"  Transactions:  %d\n"+
			"  Hash:          %s\n"+
			"  PrevBlockHash: %s\n"+
			"  MerkleRoot:    %s\n"+
			"  Nonce:         %d\n",
		b.Height,
		time.Unix(b.Timestamp, 0).UTC().Format(time.RFC3339),
		len(b.Transactions),
		hex.EncodeToString(b.Hash),
		hex.EncodeToString(b.PrevBlockHash),
		hex.EncodeToString(b.MerkleRoot),
		b.Nonce,
	)
}
```

Two small but deliberate choices here. First, `time.Unix(b.Timestamp, 0).UTC().Format(time.RFC3339)` turns the raw integer Chapter 16, Section 2 justified keeping as a plain `int64` back into a genuinely readable date — `2026-07-31T09:15:00Z` instead of `1785575700` — without changing anything about how `Block.Timestamp` is stored or hashed; formatting for humans and storing for machines are two separate concerns, and this function only ever touches the first one. Second, every hash is passed through `hex.EncodeToString`, exactly the Chapter 08/09 convention this entire course has used since Chapter 08, Section 1 — 64 readable hex characters, never raw, unprintable bytes dumped to a terminal.

## 4. Building the CLI With Go's flag Package

Go's standard `flag` package parses command-line arguments like `-db chain.dat` or `-verify` directly into typed Go variables, with almost no boilerplate. Since GoChain's CLI needs a subcommand (`inspect`) rather than a single flat set of flags, `main` dispatches on the first argument manually, then hands the rest to a dedicated `flag.FlagSet` scoped to that one subcommand:

```go
// cmd/gochain/main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gochain <command> [flags]")
		fmt.Println("  inspect   print every block in a chain, optionally verifying it")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "inspect":
		runInspect(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}
```

```go
// cmd/gochain/inspect.go
import "flag"

// runInspect parses inspect's own flags from args (everything after
// "inspect" on the command line), loads the chain at -db, prints
// every block, and -- if -verify was passed -- validates it too.
func runInspect(args []string) {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)
	dbPath := fs.String("db", "chain.dat", "path to the block file to inspect")
	verify := fs.Bool("verify", false, "walk the chain and validate every block and link")
	fs.Parse(args)

	chain, err := core.LoadBlockchain(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading chain at %s: %v\n", *dbPath, err)
		os.Exit(1)
	}

	for _, b := range chain.Blocks() {
		fmt.Print(formatBlock(b))
		fmt.Println()
	}

	if *verify {
		passed, report := verifyChain(chain)
		fmt.Print(report)
		if !passed {
			os.Exit(1)
		}
	}
}
```

`flag.NewFlagSet("inspect", flag.ExitOnError)` creates a self-contained set of flags just for this subcommand — separate from any flags a future `gochain wallet` or `gochain node` subcommand might define with the same names, since each subcommand gets its own independent `FlagSet`. `fs.String` and `fs.Bool` both return a *pointer* to the parsed value (`*string`, `*bool`), which is why every use of `dbPath` and `verify` below the `Parse` call dereferences them with a leading `*` — a small Go detail worth pausing on if it looks unfamiliar: the flag package needs a pointer so it has somewhere to write the parsed value into, before your code ever reads it.

## 5. Wiring Up --verify

`verifyChain` reuses Chapter 19's `ValidateBlock` directly — this CLI adds no new validation logic of its own, only a readable report of what Chapter 19's real checks find:

```go
// cmd/gochain/inspect.go

// verifyChain walks every block in chain and validates it using
// Chapter 19's ValidateBlock, building a line-by-line report. It
// returns false the moment it finds the first invalid block, matching
// ValidateChain's own "stop at the first failure" behavior.
func verifyChain(chain *core.Blockchain) (bool, string) {
	var report strings.Builder
	report.WriteString("Verifying chain...\n")

	for _, b := range chain.Blocks() {
		if err := chain.ValidateBlock(b); err != nil {
			report.WriteString(fmt.Sprintf("  Block %d: INVALID -- %v\n", b.Height, err))
			report.WriteString("Chain verification FAILED.\n")
			return false, report.String()
		}
		report.WriteString(fmt.Sprintf("  Block %d: OK\n", b.Height))
	}

	report.WriteString("Chain verification PASSED -- every block and link checks out.\n")
	return true, report.String()
}
```

Building the report into a `strings.Builder` rather than printing each line directly (with `fmt.Println`) is what makes `verifyChain` testable without capturing anything from the terminal — Section 8 calls it directly and inspects the returned string and boolean, exactly the way `formatBlock` was made testable in Section 3. `runInspect` is the only place that actually decides what to *do* with the result — print it, and exit with a non-zero status if verification failed, following the standard Unix convention that a command's exit code should reflect whether it succeeded.

## 6. Running the Inspector — A Worked Session

Build a small chain first, using the persistence tools from Chapter 20:

```go
// a one-off setup program, not part of cmd/gochain itself
chain, _ := core.OpenBlockchain("chain.dat")
for i := 1; i <= 3; i++ {
	tx := &core.Transaction{ID: []byte(fmt.Sprintf("tx-%d", i)), Timestamp: time.Now().Unix()}
	next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)
	chain.AddBlock(next)
	core.AppendBlockToFile("chain.dat", next)
}
```

Now run the inspector:

```
$ go run ./cmd/gochain inspect -db chain.dat

Block 0
  Timestamp:     2026-07-31T09:00:00Z
  Transactions:  1
  Hash:          7a3fc9e1d84f3b6e1a5c9d02f7e8b1a3c4d5e6f708192a3b4c5d6e7f8091a2b3
  PrevBlockHash: 0000000000000000000000000000000000000000000000000000000000000000
  MerkleRoot:    2f8c1a9d3e5b7c0f1a2b3c4d5e6f708192a3b4c5d6e7f8091a2b3c4d5e6f7081
  Nonce:         0

Block 1
  Timestamp:     2026-07-31T09:00:01Z
  Transactions:  1
  Hash:          c81e44b0f0a291d7e3b5c8a9f2013d4e5f6a7b8c9d0e1f203141516171819202
  PrevBlockHash: 7a3fc9e1d84f3b6e1a5c9d02f7e8b1a3c4d5e6f708192a3b4c5d6e7f8091a2b3
  MerkleRoot:    9c2e0a71d84f3b6e1a5c9d02f7e8b1a3c4d5e6f708192a3b4c5d6e7f8091a2b3

  Nonce:         0

... (blocks 2 and 3 follow the same shape)
```

Now run it again, adding `-verify`:

```
$ go run ./cmd/gochain inspect -db chain.dat -verify

Block 0
  ...
Block 3
  ...

Verifying chain...
  Block 0: OK
  Block 1: OK
  Block 2: OK
  Block 3: OK
Chain verification PASSED -- every block and link checks out.
```

This single command is now your everyday debugging workflow for the rest of this course. Every time you build something that produces or modifies a chain — a new mining loop, a networking bug, a storage migration — running `gochain inspect -verify` against its output is the fastest way to answer "did I just break something" without writing a single new line of throwaway code.

## 7. Handling a Corrupted or Missing Chain Gracefully

A debugging tool that itself panics with an unreadable stack trace the moment something is wrong with the data it is debugging is not a very good debugging tool. `runInspect` already handles the most common failure — a missing or corrupted file — through the ordinary `error` return from `core.LoadBlockchain`, printed as a clear, specific message on `os.Stderr` rather than left to crash:

```
$ go run ./cmd/gochain inspect -db does-not-exist.dat
error loading chain at does-not-exist.dat: opening block file does-not-exist.dat: open does-not-exist.dat: no such file or directory
```

The other interesting failure is a chain that loads fine — every record decodes successfully — but fails `-verify` because a block was tampered with, exactly the scenario Chapter 19's worked example built by hand:

```
$ go run ./cmd/gochain inspect -db tampered-chain.dat -verify

Block 0
  ...
Block 3
  ...

Verifying chain...
  Block 0: OK
  Block 1: INVALID -- block 1: stored MerkleRoot 9c2e0a71... does not match a freshly computed Merkle root 3d8f1a29... over its transactions
Chain verification FAILED.
```

Notice this is precisely the situation Chapter 19, Section 6 walked through by hand, now surfaced through a real, runnable tool instead of ad-hoc `fmt.Println` calls scattered through a test file. The exit code (`os.Exit(1)` on a failed verification, back in Section 4's `runInspect`) also matters beyond the printed text: a future Chapter 92 CI/CD pipeline can run `gochain inspect -verify` as an automated check and treat a non-zero exit code as "the build should fail," exactly the way any other command-line tool's exit status is used to gate automation.

## 8. Testing the CLI Logic

Because Sections 3 and 5 deliberately returned strings and booleans instead of printing and exiting directly, testing this chapter's real logic needs no special tricks for capturing stdout or intercepting `os.Exit` at all.

```go
// cmd/gochain/inspect_test.go
package main

import (
	"strings"
	"testing"

	"github.com/you/gochain/core"
)

func buildTestChain(t *testing.T, n int) *core.Blockchain {
	t.Helper()
	chain := core.NewBlockchain()
	for i := 1; i <= n; i++ {
		tx := &core.Transaction{ID: []byte("tx"), Timestamp: 1719792000}
		next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)
		if err := chain.AddBlock(next); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	return chain
}

func TestFormatBlock_ContainsExpectedFields(t *testing.T) {
	chain := buildTestChain(t, 1)
	out := formatBlock(chain.LastBlock())

	for _, want := range []string{"Block 1", "Timestamp:", "Transactions:", "Hash:", "PrevBlockHash:", "MerkleRoot:", "Nonce:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected formatted block to contain %q, got:\n%s", want, out)
		}
	}
}

func TestVerifyChain_HealthyChainPasses(t *testing.T) {
	chain := buildTestChain(t, 4)

	passed, report := verifyChain(chain)
	if !passed {
		t.Fatalf("expected a healthy chain to pass, got report:\n%s", report)
	}
	if !strings.Contains(report, "PASSED") {
		t.Errorf("expected report to mention PASSED, got:\n%s", report)
	}
}

func TestVerifyChain_TamperedChainFails(t *testing.T) {
	chain := buildTestChain(t, 3)

	tampered := chain.Blocks()[1]
	tampered.Transactions[0].ID = []byte("forged")
	// Hash deliberately not recomputed -- exactly Chapter 19's scenario.

	passed, report := verifyChain(chain)
	if passed {
		t.Fatal("expected a tampered chain to fail verification")
	}
	if !strings.Contains(report, "Block 1: INVALID") {
		t.Errorf("expected report to flag block 1 as invalid, got:\n%s", report)
	}
}
```

Run `go test ./cmd/gochain/...`. `TestVerifyChain_TamperedChainFails` is worth noticing in particular: it is the same tampering scenario Chapter 19, Section 6 walked through by hand, now checked automatically at the CLI layer too — proof that the inspector's `-verify` flag genuinely surfaces the exact failure Chapter 19's `ValidateBlock` was built to catch, rather than silently passing something it should not.

---

## Summary

- A chain inspector is a debugging tool worth building early: as chains grow more complex in later volumes, having one reliable command to print and validate a chain saves writing a new throwaway program every time.
- `gochain inspect` supports two flags: `-db` (which block file to open) and `-verify` (whether to also walk and validate the chain).
- `formatBlock` renders one block as a readable multi-line summary — height, human-formatted timestamp, transaction count, and hex-encoded `Hash`, `PrevBlockHash`, and `MerkleRoot` — returning a string rather than printing directly, which keeps it trivially testable.
- Go's standard `flag` package, via a `flag.FlagSet` scoped to the `inspect` subcommand, is enough for a CLI this small; Chapter 74 introduces Cobra once GoChain's command surface grows much larger.
- `verifyChain` adds no new validation logic — it calls Chapter 19's `ValidateBlock` on every block and formats the result into a readable, testable report.
- Errors from a missing or corrupted block file, and validation failures on a tampered chain, are both reported with clear messages and a non-zero exit code, rather than crashing with a raw stack trace.
- Returning strings and booleans from the core logic (`formatBlock`, `verifyChain`), rather than printing and exiting directly inside them, is what makes CLI code testable without needing to capture stdout or intercept `os.Exit`.

---

## Exercises

### Easy

1. Add a `-json` flag to `gochain inspect` that, when set, prints each block as a JSON object (using `encoding/json`) instead of `formatBlock`'s plain-text format. Test it against a small chain and confirm the output is valid JSON you can pipe into a tool like `jq`.
2. Run `gochain inspect` against a path that does not exist, and against a path that exists but is not a valid block file (for example, a text file containing the word "hello"). Compare the two error messages and explain, in your own words, why they differ.
3. `formatBlock` prints `len(b.Transactions)` but never prints the transactions themselves. Add a `-full` flag that, when set, also prints each transaction's `ID` (hex-encoded) indented under its block.

### Medium

4. Add a `-height` flag that, when set to a non-negative integer, prints only the single block at that height instead of the whole chain (returning a clear error if the height does not exist). Write a test covering both the found and not-found cases.
5. `verifyChain` currently stops at the first invalid block, matching `ValidateChain`'s behavior. Write an alternative `verifyChainFull` that continues past a failure, collecting and reporting *every* invalid block found, and add a `-continue-on-error` flag that switches between the two behaviors.
6. Extend the inspector with a second subcommand, `gochain summary`, that prints only chain-wide statistics: total block count, total transaction count across all blocks, and the timestamp range (earliest to latest block). Reuse `core.LoadBlockchain` and write the dispatch logic in `main.go` following the pattern already used for `inspect`.

### Hard

7. Rewrite `runInspect`'s error handling so that a corrupted trailing record (Chapter 20, Exercise 4's truncated-file scenario) produces a distinctly different, more specific error message than a completely missing file — for example, explicitly mentioning how many blocks were successfully read before the corruption was hit. This will likely require the "best effort" `LoadBlockchain` variant proposed in Chapter 20's exercises.
8. Design and implement a `-diff old.dat new.dat` mode for the inspector that loads two block files and reports exactly where they first diverge (the lowest height at which the two chains' blocks differ), printing both blocks at that height side by side. Consider what "diverge" should mean if one file is simply a shorter prefix of the other.
9. Package `gochain inspect -verify` as a step in a shell script that exits non-zero (and prints a clear failure message) if any `.dat` file in a given directory fails verification, suitable for wiring into a CI pipeline (a real version of this appears in Chapter 92). Test it against a directory containing one healthy and one deliberately tampered chain file, and confirm the script correctly identifies which file is broken.
