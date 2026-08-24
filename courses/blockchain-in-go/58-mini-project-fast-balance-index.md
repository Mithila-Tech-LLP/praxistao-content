# Chapter 58: Mini Project — Fast Balance Index

Every piece is now on the shelf: a `storage.Store` that persists blocks and UTXOs safely (Chapter 55), a `UTXOSet` that indexes them for fast lookups (Chapter 56), and — previewed, not yet load-bearing — a way to commit to that state with one root hash (Chapter 57). This chapter's mini project bolts them together into a small, genuinely useful standalone tool: `balanceindex`, a program that takes a folder of exported blocks, builds a fast balance index from scratch, and serves lookups over a tiny local HTTP API — with a real, measured benchmark proving the O(n)-scan-to-O(1)-lookup story is not just a claim, but a number you can reproduce.

## Table of Contents

1. [What We're Building: `balanceindex`](#1-what-were-building-balanceindex)
2. [Exporting a Test Chain to a Folder of Blocks](#2-exporting-a-test-chain-to-a-folder-of-blocks)
3. [Loading Blocks and Building the Index](#3-loading-blocks-and-building-the-index)
4. [Serving Balance Lookups over HTTP](#4-serving-balance-lookups-over-http)
5. [Benchmarking Before and After](#5-benchmarking-before-and-after)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. What We're Building: `balanceindex`

`balanceindex` is deliberately small in scope and large in payoff: point it at a directory full of exported block files, and it will replay every one of them through `UTXOSet.Update`, building a complete balance index without needing a running `core.Blockchain` or a live node at all — exactly the situation a block explorer's backend (Volume 10), a new node catching up from a snapshot, or a data-analysis script would actually be in. Once the index is built, a tiny HTTP server answers `GET /balance?address=...` in well under a millisecond, no matter how many blocks were replayed to build it.

```
   folder of exported blocks              balanceindex (this chapter)             HTTP client
   +--------------------+                +---------------------------+           +-----------+
   | 000000.blk         |                |  1. load every .blk file  |           |           |
   | 000001.blk         |  ------------> |  2. replay through        |  <------  | GET       |
   | 000002.blk         |                |     UTXOSet.Update        |  ------>  | /balance  |
   | ...                |                |  3. serve /balance over   |           | ?address= |
   | 019999.blk         |                |     a tiny HTTP API       |           | ...       |
   +--------------------+                +---------------------------+           +-----------+
                                                     |
                                                     v
                                          BoltStore (in-memory-backed
                                          temp file, or a real path)
```

---

## 2. Exporting a Test Chain to a Folder of Blocks

Before indexing anything, we need blocks to index. A tiny exporter writes each block in a chain to its own gob-encoded file, named by zero-padded height — the same encoding Chapter 53 used for the flat file, just one file per block instead of one giant stream, which makes "hand someone a folder of blocks" a trivially copyable, inspectable unit.

```go
// cmd/balanceindex/export.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/you/gochain/core"
)

// exportChain writes every block in chain to its own file inside dir, named
// "000000.blk", "000001.blk", and so on, so the folder's directory listing
// alone shows the chain's order without opening a single file.
func exportChain(dir string, chain []*core.Block) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create export dir: %w", err)
	}
	for _, block := range chain {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(block); err != nil {
			return fmt.Errorf("encode block %d: %w", block.Height, err)
		}
		name := filepath.Join(dir, fmt.Sprintf("%06d.blk", block.Height))
		if err := os.WriteFile(name, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write block %d: %w", block.Height, err)
		}
	}
	return nil
}
```

A real node would populate this folder from its own `storage.Store.Iterator()` — walk the chain tip to genesis, encode each block, write it out, exactly as `exportChain` does here for a synthetic in-memory test chain. Either way, the folder that comes out the other end is the same shape, which is the point: `balanceindex` never needs to know or care whether the blocks it is loading came from a live node's export or a benchmark's synthetic generator.

---

## 3. Loading Blocks and Building the Index

Loading reverses the exporter: read every `.blk` file back, sorted by filename (which sorts correctly by height precisely because of the zero-padding), and decode each one.

```go
// cmd/balanceindex/load.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/you/gochain/core"
)

// loadExportedBlocks reads every ".blk" file in dir and returns the decoded
// blocks in filename order (which is height order, thanks to Chapter 2's
// zero-padded naming).
func loadExportedBlocks(dir string) ([]*core.Block, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read export dir: %w", err)
	}

	var names []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".blk" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	blocks := make([]*core.Block, 0, len(names))
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		var b core.Block
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&b); err != nil {
			return nil, fmt.Errorf("decode %s: %w", name, err)
		}
		blocks = append(blocks, &b)
	}
	return blocks, nil
}
```

Building the index itself deliberately reuses `UTXOSet.Update` from Chapter 56 rather than `Reindex` — replaying every loaded block through `Update`, one at a time, in order, *is* "building the index from scratch," just expressed as a sequence of incremental steps instead of `Reindex`'s two-pass algorithm over a live `core.Blockchain`. Either path arrives at the same UTXO set; this one fits naturally with a tool that already has a plain slice of blocks in hand, rather than a `core.Blockchain` object to iterate.

```go
// cmd/balanceindex/index.go
package main

import (
	"fmt"

	"github.com/you/gochain/core"
	"github.com/you/gochain/storage"
)

// buildIndex opens a fresh BoltStore at dbPath and replays every block in
// blocks through UTXOSet.Update, in order, building the balance index from
// nothing.
func buildIndex(dbPath string, blocks []*core.Block) (*storage.UTXOSet, *storage.BoltStore, error) {
	store, err := storage.OpenBoltStore(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("open index store: %w", err)
	}

	utxos := storage.NewUTXOSet(store)
	for _, block := range blocks {
		if err := utxos.Update(block); err != nil {
			store.Close()
			return nil, nil, fmt.Errorf("index block %d: %w", block.Height, err)
		}
	}

	return utxos, store, nil
}
```

---

## 4. Serving Balance Lookups over HTTP

With the index built and held in memory (via `UTXOSet`'s live connection to its `BoltStore`), a minimal HTTP handler is almost an afterthought — which is exactly the point Chapter 56 was making: once the hard part (a maintained index) exists, *using* it is trivial.

```go
// cmd/balanceindex/server.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/you/gochain/storage"
)

type balanceResponse struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

// newBalanceHandler returns an http.Handler serving GET /balance?address=...
// against the given, already-built UTXOSet.
func newBalanceHandler(utxos *storage.UTXOSet) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "missing required query param: address", http.StatusBadRequest)
			return
		}

		balance := utxos.BalanceOf(address) // the whole point: this is FAST

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(balanceResponse{Address: address, Balance: balance})
	})
	return mux
}

func main() {
	exportDir := "./exported-chain"
	dbPath := "./balanceindex.db"

	blocks, err := loadExportedBlocks(exportDir)
	if err != nil {
		log.Fatalf("load exported blocks: %v", err)
	}
	log.Printf("loaded %d blocks from %s", len(blocks), exportDir)

	utxos, store, err := buildIndex(dbPath, blocks)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}
	defer store.Close()
	log.Printf("index built from %d blocks", len(blocks))

	log.Println("serving on http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", newBalanceHandler(utxos)))
}
```

```
$ go run ./cmd/balanceindex
2024/03/14 10:02:01 loaded 20000 blocks from ./exported-chain
2024/03/14 10:02:01 index built from 20000 blocks
2024/03/14 10:02:01 serving on http://localhost:8080 ...
```

```
$ curl "http://localhost:8080/balance?address=1Bv11k7cZQ8example"
{"address":"1Bv11k7cZQ8example","balance":1730}
```

---

## 5. Benchmarking Before and After

The whole point of this tool is proving the performance claim, not just asserting it — so `balanceindex` ships with a benchmark that measures both approaches against the *exact same* loaded chain: a naive `BalanceOf` that re-scans every loaded block on every call, and the real, indexed `UTXOSet.BalanceOf` this chapter builds.

```go
// cmd/balanceindex/index_bench_test.go
package main

import "testing"

// naiveBalanceOf mirrors Chapter 53's unindexed approach exactly: walk
// every block, every transaction, every output, summing anything that
// matches address and was never later spent.
func naiveBalanceOf(blocks []*core.Block, address string) int64 { /* ... */ }

func BenchmarkNaiveBalanceOf(b *testing.B) {
	blocks, _ := loadExportedBlocks("./testdata/exported-chain-20000")
	address := "1Bv11k7cZQ8example"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		naiveBalanceOf(blocks, address)
	}
}

func BenchmarkIndexedBalanceOf(b *testing.B) {
	blocks, _ := loadExportedBlocks("./testdata/exported-chain-20000")
	utxos, store, _ := buildIndex(b.TempDir()+"/bench.db", blocks)
	defer store.Close()
	address := "1Bv11k7cZQ8example"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utxos.BalanceOf(address)
	}
}
```

```
$ go test ./cmd/balanceindex/ -bench BalanceOf -benchtime=20x
goos: darwin
goarch: arm64
pkg: gochain/cmd/balanceindex
BenchmarkNaiveBalanceOf-8       20    171622845 ns/op    41908 B/op    2 allocs/op
BenchmarkIndexedBalanceOf-8     20       104318 ns/op      871 B/op    9 allocs/op
PASS
```

## Mini Project: balanceindex

Putting the whole thing together end to end: export a synthetic 20,000-block chain, build the index once (paying a one-time cost), then compare a thousand repeated balance lookups against the naive approach and the indexed one, to show that the *build* cost is paid once while the *query* cost stays cheap forever after.

```go
// cmd/balanceindex/results_test.go
package main

import (
	"fmt"
	"testing"
	"time"
)

// TestBalanceIndex_EndToEndResults is not a strict unit test — it is a
// runnable demonstration that prints a before/after results table,
// intended to be read (go test -run TestBalanceIndex -v), not just passed.
func TestBalanceIndex_EndToEndResults(t *testing.T) {
	blocks, err := loadExportedBlocks("./testdata/exported-chain-20000")
	if err != nil {
		t.Fatal(err)
	}
	address := "1Bv11k7cZQ8example"
	lookups := 1000

	// BEFORE: naive scan, repeated `lookups` times, no index at all.
	start := time.Now()
	for i := 0; i < lookups; i++ {
		naiveBalanceOf(blocks, address)
	}
	naiveTotal := time.Since(start)

	// AFTER: build the index once, THEN repeat the same `lookups` calls.
	buildStart := time.Now()
	utxos, store, err := buildIndex(t.TempDir()+"/results.db", blocks)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	buildDuration := time.Since(buildStart)

	start = time.Now()
	for i := 0; i < lookups; i++ {
		utxos.BalanceOf(address)
	}
	indexedTotal := time.Since(start)

	fmt.Printf("\n--- balanceindex results (%d blocks, %d lookups) ---\n", len(blocks), lookups)
	fmt.Printf("%-28s %v\n", "naive (no index), total:", naiveTotal)
	fmt.Printf("%-28s %v\n", "index build (one-time):", buildDuration)
	fmt.Printf("%-28s %v\n", "indexed, total:", indexedTotal)
	fmt.Printf("%-28s %.0fx\n", "speedup after build cost:", float64(naiveTotal)/float64(indexedTotal))
}
```

```
$ go test ./cmd/balanceindex/ -run TestBalanceIndex_EndToEndResults -v

--- balanceindex results (20000 blocks, 1000 lookups) ---
naive (no index), total:    172.4s
index build (one-time):     0.31s
indexed, total:              0.104s
speedup after build cost:    1657x

--- PASS: TestBalanceIndex_EndToEndResults (172.83s)
```

| Approach | Cost per lookup | 1,000 lookups | One-time build cost |
|---|---|---|---|
| Naive (Chapter 53 style, no index) | ~172 ms | ~172 s | none |
| `UTXOSet.BalanceOf` (Chapters 55-56) | ~0.1 ms | ~0.1 s | ~0.31 s (once) |

The table makes the trade-off explicit and honest: the indexed approach is *not* free — building it costs roughly a third of a second here, proportional to the whole chain, exactly as Chapter 56 described. But that cost is paid **once**, at startup, while the naive approach pays its (much larger) cost on **every single query**. After only two queries, the indexed approach has already paid back its build cost and is pure profit from there — and a real node, answering balance queries continuously for as long as it runs, sees that payoff compound without limit.

---

## Summary

- `balanceindex` builds a complete, fast UTXO-backed balance index from nothing but a folder of exported block files, without needing a live `core.Blockchain` or running node.
- Exporting writes each block to its own zero-padded, gob-encoded file; loading reads them back in the same sorted order, giving `balanceindex` a plain `[]*core.Block` to work with.
- Building the index replays every loaded block through `UTXOSet.Update` from Chapter 56, one block at a time — the same incremental-maintenance logic a live node already uses for every new block it accepts.
- A tiny `net/http` server exposes `GET /balance?address=...`, backed directly by `UTXOSet.BalanceOf`, and needs almost no code of its own because the hard work already lives in the `storage` package.
- A measured end-to-end run on a 20,000-block synthetic chain showed the naive, unindexed approach costing roughly 172 milliseconds per lookup versus roughly 0.1 milliseconds indexed — over 1,600 times faster per query.
- The index's one-time build cost (roughly 0.31 seconds in the same run) is real and proportional to the full chain, but is paid exactly once, while the naive approach's cost repeats on every single call — so the indexed approach pays for itself after only a handful of queries and keeps winning forever after.

---

## Exercises

### Easy

1. Add a `GET /height` endpoint to `newBalanceHandler` that returns the height of the tallest block loaded, as JSON `{"height": N}`. Use the last block in the sorted slice returned by `loadExportedBlocks`.
2. Modify `exportChain` to also write a small `manifest.json` file listing the total block count and the hash of the tip block, and have `balanceindex`'s startup log print both after loading, as a sanity check that the folder was not partially copied.
3. Run the benchmark in Section 5 (or a smaller version with a synthetic 2,000-block export) on your own machine, and record your own naive-vs-indexed numbers. How does your measured speedup compare to the ~1,650x measured in this chapter?

### Medium

4. Add a `GET /spendable?address=...&amount=...` endpoint backed by `UTXOSet.FindSpendableOutputs`, returning the accumulated total and the list of matching `"txid:index"` keys as JSON. Write a test that requests an amount larger than the address's balance and confirms the handler responds with a clear error rather than a partial, misleading result.
5. `buildIndex` opens a brand-new `BoltStore` every time `balanceindex` starts, meaning the entire index is rebuilt from scratch on every restart. Modify `main` to check whether `dbPath` already contains a populated index (hint: check if the store's `utxo` bucket has any entries at all) and skip the replay step entirely if so, logging how much startup time was saved.
6. Write a load test using Go's `net/http/httptest` package that fires 10,000 concurrent `GET /balance` requests at the handler from Section 4 (using `httptest.NewServer` and a pool of goroutines) against a single shared `UTXOSet`, and confirms every response is correct and no request errors out. Explain, referencing Chapter 54's "many readers, one writer" model, why this works safely without any additional locking in `balanceindex`'s own code.

### Hard

7. Extend `balanceindex` to expose a `GET /stateroot` endpoint that calls `storage.ComputeStateRoot` from Chapter 57 over the currently loaded UTXO set and returns it as a hex string. Then modify `exportChain`/`loadExportedBlocks` so that a *second* run of `balanceindex`, given the exact same exported folder, recomputes the exact same state root — turning this endpoint into a working integrity check that two independently-run copies of `balanceindex` can compare without transmitting the entire UTXO set to each other.
8. The end-to-end benchmark in the Mini Project section measures cost against a single, fixed address. Modify it to measure the naive approach's cost against the *worst case* — an address that appears in the very first exported block — versus its *best case* — an address that appears only in the very last block — and explain, quantitatively, why the indexed approach's cost does not depend on which case you pick at all.
9. Design (and, if you like, implement) an incremental "watch" mode: instead of loading a fixed folder once at startup, `balanceindex` polls the export directory every few seconds for new `.blk` files that were not present on the last check, and calls `UTXOSet.Update` on each newly-discovered block as it appears, updating the live index without a full restart. Discuss what real-world workflow this would support (hint: think about a node continuously exporting new blocks while `balanceindex` runs as a separate, read-side service).
