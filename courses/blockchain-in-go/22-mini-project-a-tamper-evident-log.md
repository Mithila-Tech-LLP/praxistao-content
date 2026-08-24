# Chapter 22: Mini Project — A Tamper-Evident Log

Every chapter in this volume built toward one goal: a chain of blocks where tampering with any old entry is unmistakably detectable. Nothing about that mechanism actually requires money, mining, or a cryptocurrency at all — it requires only a sequence of records, each one carrying a fingerprint of the record before it. This chapter proves that by ripping the mechanism out of its cryptocurrency context entirely and pointing it at something mundane: a tamper-evident activity log for a fictional file-sharing app, built directly on `core.Block` and `core.Blockchain`, with zero changes to either.

## Table of Contents

1. [Beyond Cryptocurrency — Block-Chaining as a General Tool](#1-beyond-cryptocurrency--block-chaining-as-a-general-tool)
2. [Designing DropVault's Event Model](#2-designing-dropvaults-event-model)
3. [Action — One Auditable Event, Encoded as Bytes](#3-action--one-auditable-event-encoded-as-bytes)
4. [Log — Wrapping core.Blockchain for a Non-Financial Use Case](#4-log--wrapping-coreblockchain-for-a-non-financial-use-case)
5. [The auditlog CLI](#5-the-auditlog-cli)

---

## 1. Beyond Cryptocurrency — Block-Chaining as a General Tool

Strip away everything this course has associated with "blockchain" so far — coins, wallets, mining, transactions moving value between people — and what remains is a much smaller, more general idea: **append new records one at a time, and make each new record's fingerprint depend on the one before it.** Chapter 01 introduced blockchains through the analogy of a shared notebook and a chain of wax-sealed envelopes; those analogies were never really about money either. They were about one specific problem — *how do you know a written history has not been quietly rewritten?* — that shows up constantly outside cryptocurrency entirely.

Think about a hospital's patient record system, a company's internal changelog of who approved which expense, or a courthouse's log of every time a piece of digital evidence was accessed. Every one of these is a sequence of "something happened" records where the exact same worry applies: could someone with access to the underlying storage quietly edit an old entry, and would anyone notice? Chapters 17 through 21 already built the entire toolkit needed to answer that worry with real, checkable proof instead of institutional trust — `core.Block`, `core.Blockchain`, `ValidateBlock`, `ValidateChain`, persistence, and an inspector. This chapter's mini project, **auditlog**, applies every one of those tools, completely unmodified, to a scenario with no coins in it whatsoever.

```
   Cryptocurrency use (Volume 3 so far)     General tamper-evident log (this chapter)

   Block.Transactions holds transfers        Block.Transactions holds application events
   ("Alice paid Bob 5 gochips")               ("alice uploaded report.pdf")

   PrevBlockHash links blocks so a            PrevBlockHash links blocks so an audit
   rewritten balance is detectable            trail entry cannot be quietly edited

   ValidateChain proves the ledger's          ValidateChain proves the log's history
   history has not been altered               has not been altered

               Same core.Block. Same core.Blockchain. Same guarantees.
```

## 2. Designing DropVault's Event Model

**DropVault** is this chapter's fictional file-sharing app — the kind of product where "who did what to which file, and when" genuinely matters: a user should be able to trust that DropVault's own activity log has not been edited after the fact, whether by a bug, an over-eager support engineer, or an attacker trying to cover their tracks. DropVault needs to record four kinds of events:

- **upload** — a user adds a new file.
- **delete** — a user removes a file.
- **rename** — a user changes a file's name.
- **share** — a user grants another user access to a file.

Every one of these events shares the same shape: *someone* (`Actor`) did *something* (`Type`) to *some file* (`Target`), possibly with one extra piece of context (`Detail` — a rename's new name, or a share's recipient), at *some time* (`Timestamp`). That shape is small and generic enough to reuse across all four event types without needing four different Go structs, which keeps this mini project's domain model — deliberately — much simpler than `core.Transaction`'s eventual UTXO design in Volume 5.

```
  Action{ Type: "delete", Actor: "alice", Target: "report.pdf", Detail: "", Timestamp: ... }

  reads as: "alice deleted report.pdf, at [timestamp]"
```

## 3. Action — One Auditable Event, Encoded as Bytes

`core.Transaction` (Chapter 17, Section 2) is deliberately opaque outside `core` at this stage of the course — a placeholder with an `ID`, some inputs and outputs meant for Volume 5's UTXO model, and a `Timestamp`. `auditlog` does not need any of the UTXO fields at all; it only needs *some* way to carry arbitrary application data through `core.Block`'s existing hashing and linking machinery. The simplest option available today, without touching `core` at all, is to encode each `Action` as bytes and store those bytes directly in a transaction's `ID` field — repurposing it, honestly and deliberately, as a general-purpose payload rather than a real cryptographic transaction identifier.

```go
// auditlog/action.go
package auditlog

import "encoding/json"

// Action is one auditable event in DropVault: someone did something
// to some file, at some time. It carries no cryptographic meaning of
// its own -- it is ordinary application data, encoded into bytes and
// handed to core.Block exactly the way any payload would be.
type Action struct {
	Type      string `json:"type"`      // "upload", "delete", "rename", "share"
	Actor     string `json:"actor"`     // who did it
	Target    string `json:"target"`    // which file
	Detail    string `json:"detail"`    // rename's new name, or share's recipient
	Timestamp int64  `json:"timestamp"`
}

// Encode turns an Action into JSON bytes. JSON is a deliberate choice
// here, not a contradiction of Chapter 09's canonical-serialization
// rules: those rules matter when bytes are the input to a hash
// computation you need to reproduce exactly (core.Block.Serialize()
// still follows them, unchanged, underneath this entire chapter).
// Here, Action's bytes are just opaque application payload being
// carried *through* that hashing machinery, the same way a real
// Transaction's inputs and outputs will be in Volume 5 -- so ordinary
// JSON, easy for a human to read directly off disk, is the right
// tool for this particular job.
func (a Action) Encode() []byte {
	data, _ := json.Marshal(a)
	return data
}

// DecodeAction reverses Encode. It returns an error for anything
// that is not valid Action JSON -- notably, the genesis block's
// placeholder coinbase transaction from Chapter 18, which Section 4
// below learns to skip gracefully rather than treat as a real event.
func DecodeAction(data []byte) (Action, error) {
	var a Action
	err := json.Unmarshal(data, &a)
	return a, err
}
```

Storing an `Action`'s JSON bytes inside `Transaction.ID` is a repurposing worth being honest about, not a design pattern to imitate blindly elsewhere in GoChain: `ID` will mean something precise and cryptographic once Volume 5 gives `Transaction` a real identity scheme. For `auditlog`'s purposes today, it is simply the one field on the placeholder `Transaction` struct available to carry a payload, and that is a perfectly legitimate way to explore what `core.Block`'s hashing and linking guarantees can do for data that has nothing to do with money.

## 4. Log — Wrapping core.Blockchain for a Non-Financial Use Case

`Log` is a thin wrapper — a few dozen lines — around exactly the tools Chapters 18 through 20 already built: `core.OpenBlockchain`, `AddBlock`, `AppendBlockToFile`, and `ValidateChain`. Nothing here reimplements anything; it only gives those general-purpose tools a vocabulary that fits DropVault.

```go
// auditlog/log.go
package auditlog

import (
	"fmt"

	"github.com/you/gochain/core"
)

// Log is a tamper-evident audit trail. Every recorded Action becomes
// one core.Block, linked to the block before it exactly the way
// Volume 3 built for cryptocurrency -- reused here, unmodified, for
// an entirely different purpose.
type Log struct {
	chain *core.Blockchain
	path  string
}

// Open loads an existing audit log from path (Chapter 20's
// OpenBlockchain), or starts a fresh one -- complete with its own
// genesis block -- if none exists yet.
func Open(path string) (*Log, error) {
	chain, err := core.OpenBlockchain(path)
	if err != nil {
		return nil, fmt.Errorf("opening audit log at %s: %w", path, err)
	}
	return &Log{chain: chain, path: path}, nil
}

// Record appends a to the log as a new block and persists it
// immediately, so every recorded action survives a restart the
// moment it is recorded -- there is no separate "save" step to forget.
func (l *Log) Record(a Action) (*core.Block, error) {
	tx := &core.Transaction{ID: a.Encode(), Timestamp: a.Timestamp}
	next := core.NewBlock([]*core.Transaction{tx}, l.chain.Tip(), l.chain.Height()+1)

	if err := l.chain.AddBlock(next); err != nil {
		return nil, fmt.Errorf("recording action: %w", err)
	}
	if err := core.AppendBlockToFile(l.path, next); err != nil {
		return nil, fmt.Errorf("persisting action: %w", err)
	}
	return next, nil
}

// History returns every recorded Action, oldest first. The genesis
// block's placeholder coinbase transaction is not valid Action JSON,
// so DecodeAction's error on it is used here to skip it silently
// rather than treat it as a real DropVault event.
func (l *Log) History() []Action {
	var actions []Action
	for _, b := range l.chain.Blocks() {
		if len(b.Transactions) == 0 {
			continue
		}
		a, err := DecodeAction(b.Transactions[0].ID)
		if err != nil {
			continue
		}
		actions = append(actions, a)
	}
	return actions
}

// Verify walks the entire log, reusing Chapter 19's ValidateChain
// without modification, and confirms every block -- and every link
// between blocks -- is still exactly what it claims to be.
func (l *Log) Verify() error {
	return l.chain.ValidateChain()
}
```

Every one of `Log`'s four methods delegates almost immediately to `core`. This is deliberate, and worth noticing as the actual point of the whole chapter: `auditlog` is not a new tamper-evidence mechanism, it is the *same* mechanism, wearing different vocabulary. `Record` is `NewBlock` plus `AddBlock` plus `AppendBlockToFile`, in the exact same order Chapter 20's examples used. `Verify` is a one-line call straight through to `ValidateChain`. If Chapter 19's tamper-evidence guarantee holds for `core.Blockchain` in general — and Chapter 19 spent an entire worked example proving it does — it holds here too, automatically, with no additional proof required.

## 5. The auditlog CLI

Following Chapter 21's exact pattern — a small `flag`-based dispatcher, one `FlagSet` per subcommand — `auditlog` gets three commands: `record`, `history`, and `verify`.

```go
// cmd/auditlog/main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "record":
		runRecord(os.Args[2:])
	case "history":
		runHistory(os.Args[2:])
	case "verify":
		runVerify(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("usage: auditlog <command> [flags]")
	fmt.Println("  record   -db path -type T -actor A -target F [-detail D]")
	fmt.Println("  history  -db path")
	fmt.Println("  verify   -db path")
}
```

```go
// cmd/auditlog/record.go
import (
	"flag"
	"time"

	"github.com/you/gochain/auditlog"
)

func runRecord(args []string) {
	fs := flag.NewFlagSet("record", flag.ExitOnError)
	dbPath := fs.String("db", "auditlog.dat", "path to the audit log file")
	actionType := fs.String("type", "", "upload | delete | rename | share")
	actor := fs.String("actor", "", "who performed the action")
	target := fs.String("target", "", "which file was affected")
	detail := fs.String("detail", "", "optional extra detail (new name, recipient, ...)")
	fs.Parse(args)

	log, err := auditlog.Open(*dbPath)
	if err != nil {
		fatal(err)
	}

	action := auditlog.Action{
		Type: *actionType, Actor: *actor, Target: *target,
		Detail: *detail, Timestamp: time.Now().Unix(),
	}

	block, err := log.Record(action)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("recorded as block %d (hash %x...)\n", block.Height, block.Hash[:8])
}
```

`history` and `verify` follow the exact same shape as Chapter 21's `runInspect`, so they are worth writing yourself as an exercise (Exercise 4) before checking them against the full listing below.

---

## Mini Project: auditlog

Here is the complete, runnable program, and the demonstration this whole chapter has been building toward: a malicious edit to the log file, caught cold by `verify`.

### Full Listing

```go
// cmd/auditlog/history.go
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/you/gochain/auditlog"
)

func runHistory(args []string) {
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	dbPath := fs.String("db", "auditlog.dat", "path to the audit log file")
	fs.Parse(args)

	log, err := auditlog.Open(*dbPath)
	if err != nil {
		fatal(err)
	}

	for _, a := range log.History() {
		fmt.Printf("[%s] %s %s %s %s\n",
			time.Unix(a.Timestamp, 0).UTC().Format(time.RFC3339),
			a.Actor, a.Type, a.Target, a.Detail)
	}
}
```

```go
// cmd/auditlog/verify.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/you/gochain/auditlog"
)

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	dbPath := fs.String("db", "auditlog.dat", "path to the audit log file")
	fs.Parse(args)

	log, err := auditlog.Open(*dbPath)
	if err != nil {
		fatal(err)
	}

	if err := log.Verify(); err != nil {
		fmt.Printf("TAMPERING DETECTED: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Audit log verified -- no tampering detected.")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
```

### Demo — Recording DropVault Activity

```
$ go run ./cmd/auditlog record -db audit.dat -type upload -actor alice -target report.pdf
recorded as block 1 (hash c81e44b0...)

$ go run ./cmd/auditlog record -db audit.dat -type share -actor alice -target report.pdf -detail bob
recorded as block 2 (hash 44b0f0a2...)

$ go run ./cmd/auditlog record -db audit.dat -type delete -actor alice -target report.pdf
recorded as block 3 (hash 9de2a5c1...)

$ go run ./cmd/auditlog history -db audit.dat
[2026-07-31T09:00:01Z] alice upload report.pdf
[2026-07-31T09:00:02Z] alice share report.pdf bob
[2026-07-31T09:00:03Z] alice delete report.pdf

$ go run ./cmd/auditlog verify -db audit.dat
Audit log verified -- no tampering detected.
```

Three real events, each its own block, linked exactly the way Chapter 18 diagrammed — and a clean bill of health from `verify`, backed by Chapter 19's three real checks, not a rubber stamp.

### Demo — Editing the Log File Directly, and Getting Caught

Now play the part of someone with direct filesystem access to `audit.dat` — no `auditlog` API, no cooperation from the program at all, exactly the threat model Chapter 19's worked example simulated. Suppose "alice" wants to hide that *she* deleted `report.pdf`, and pins the blame on a coworker instead, by editing the raw bytes of the log file with a plain text/hex editor and swapping the ASCII text `alice` for a same-length fake name, `bobby`, inside block 3's record:

```go
// simulate.go -- reproducing "opened it in a hex editor" programmatically
data, _ := os.ReadFile("audit.dat")
tampered := bytes.Replace(data, []byte("alice"), []byte("bobby"), 1) // first match only
os.WriteFile("audit.dat", tampered, 0644)
```

This works precisely because `Action.Encode()` (Section 3) stores plain, human-readable JSON — the literal ASCII text `"actor":"alice"` really does sit inside `audit.dat` as findable bytes, and swapping it for another 5-letter name does not break the surrounding gob framing or the length prefixes Chapter 20, Section 3 built, since the total byte length of the record never changes. The file remains perfectly well-formed. `history` will even decode it without a single error:

```
$ go run ./cmd/auditlog history -db audit.dat
[2026-07-31T09:00:01Z] alice upload report.pdf
[2026-07-31T09:00:02Z] alice share report.pdf bob
[2026-07-31T09:00:03Z] bobby delete report.pdf
```

At a glance, this looks like it worked — "bobby," not "alice," now appears to have deleted the file. But run `verify`:

```
$ go run ./cmd/auditlog verify -db audit.dat
TAMPERING DETECTED: block 3: stored MerkleRoot 2f8c1a9d... does not match a freshly computed Merkle root 7b1e4c93... over its transactions
```

This is Chapter 19, Section 2's second check, firing exactly as designed: block 3's `MerkleRoot` field (stored, from when the block was honestly built for "alice") no longer matches a freshly recomputed Merkle root over the transaction now sitting in `b.Transactions` (edited to say "bobby"). Editing raw bytes with a text editor changed what `history` *displays*, but it could not touch what `verify` *recomputes* — the two are independent, and only one of them is a mathematical guarantee. Exactly like Chapter 19's original worked example, an attacker determined to fully hide this change would need to recompute block 3's `MerkleRoot` and `Hash` to match the forgery, which would then break block 3's `PrevBlockHash` relationship with whatever comes after it in a longer log — the same cascading, cannot-stop-partway-through rewrite Chapter 19, Section 6 demonstrated by hand, now shown against a completely different application with zero new proof required.

---

## Summary

- Block-chaining's core guarantee — tampering with an old record is detectable — has nothing to do with cryptocurrency specifically; it applies to any append-only sequence of records where each one's fingerprint depends on the one before it.
- **auditlog** applies `core.Block` and `core.Blockchain`, completely unmodified, to DropVault's activity log: uploads, deletes, renames, and shares, each recorded as one linked block.
- `Action` encodes each event as JSON and stores it inside a placeholder `core.Transaction`'s `ID` field — a deliberate, honestly-named repurposing of a field that will carry real cryptographic meaning starting in Volume 5.
- `Log.Record`, `Log.History`, and `Log.Verify` are thin wrappers around `core.NewBlock`/`AddBlock`/`AppendBlockToFile` and `core.Blockchain.ValidateChain` — no new tamper-evidence logic was written for this chapter at all.
- Editing the log file's raw bytes directly (simulating a hex editor) can silently change what a plain listing (`history`) displays, because that listing only decodes and trusts whatever bytes are on disk.
- `verify` catches the exact same edit immediately, because it recomputes a Merkle root over the actual transaction bytes and compares it to the block's stored `MerkleRoot` — a check that does not care whether the bytes came from an honest `Record` call or a text editor.
- This mini project is the clearest possible demonstration that Chapter 19's tamper-evidence guarantee is a property of the *data structure*, not of any cryptocurrency-specific rule — the same three checks that protect a ledger of coin transfers protect a log of file uploads with zero modification.

---

## Exercises

### Easy

1. Record five DropVault events of your choosing (mixing all four `Type` values) into a fresh audit log, print its `history`, and confirm `verify` passes.
2. Explain, in your own words, why `History()` silently skips the genesis block's placeholder transaction instead of returning an error the first time it encounters something that is not valid `Action` JSON.
3. The demo's tampering example specifically swaps `"alice"` for a same-length fake name, `"bobby"`. Try swapping it for a *different*-length name instead (e.g., `"al"`) and describe, concretely, what breaks and why — referencing Chapter 20, Section 3's length-prefix framing.

### Medium

4. Write `runHistory` and `runVerify` yourself, from the descriptions in Section 5, before checking them against the Mini Project's full listing. Note any differences between your version and the one shown, and explain whether your differences are meaningful or purely stylistic.
5. Add a fifth `Action.Type`, `"restore"` (undoing a delete), and a `-restore` convenience flag to `record` that automatically sets `Type` to `"restore"`. Record a delete followed by a restore for the same file and confirm both appear correctly, in order, in `history`.
6. `Log.Record` calls `AddBlock` and `AppendBlockToFile` as two separate steps. Construct a scenario (you can do this by editing the code temporarily) where `AddBlock` succeeds but the subsequent `AppendBlockToFile` call fails (for example, by making the log file's directory read-only right before that call). Describe the resulting inconsistency between the in-memory `Log` and the on-disk file, and propose a fix.

### Hard

7. Extend `auditlog` with a `-since TIMESTAMP` flag on `history` that only prints actions recorded after a given Unix timestamp, without needing to change `Log.History()`'s signature — filter in the CLI layer instead. Then discuss, in a short paragraph, why this filtering is safe to do *after* decoding the full history rather than needing to be baked into `ValidateChain` or `Verify` itself.
8. Design and implement a `Log.VerifyRange(fromHeight, toHeight int64) error` method that validates only a contiguous sub-range of blocks, still correctly checking that the range's own internal links hold, but explicitly documenting (in a comment and a test) that it does *not* prove the range in question is actually linked back to genesis. Explain, in 200-300 words, why this weaker guarantee might still be useful for a very long-running audit log where re-verifying the entire history on every check becomes expensive.
9. DropVault, as designed, trusts whatever `-actor` string is passed to `record` with no authentication at all — anyone with access to the `auditlog` binary can record an action claiming to be anyone. Write a design proposal (350-500 words) for adding real accountability to `auditlog` using tools already built earlier in this course (Chapter 12's ECDSA signatures, Chapter 13's `crypto.Sign`/`Verify`), such that `Verify` could also confirm each recorded action was genuinely signed by the actor it claims, not just that the log's block-linking has not been tampered with. Be explicit about what new field(s) `Action` would need, and what new check `Log.Verify` would need to add.
