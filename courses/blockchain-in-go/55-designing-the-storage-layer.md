# Chapter 55: Designing the Storage Layer

Chapter 53 broke the flat file on purpose. Chapter 54 picked BoltDB as the replacement and showed you its raw API: buckets, `db.Update`, `db.View`. Neither chapter, though, wired anything into GoChain itself. This chapter does exactly that — it defines the `storage.Store` interface the rest of GoChain will depend on, designs the three BoltDB buckets that back it, and implements `BoltStore`, a complete, working `Store` built on `go.etcd.io/bbolt`. By the end of this chapter, `core`, `consensus`, and every future package that needs to save or load a block will import `gochain/storage` and never once mention the word "Bolt."

## Table of Contents

1. [Recap: What Chapters 53 and 54 Established](#1-recap-what-chapters-53-and-54-established)
2. [Why Store Is an Interface, Not `*bolt.DB`](#2-why-store-is-an-interface-not-boltdb)
3. [Designing `storage.Store` and `storage.Iterator`](#3-designing-storagestore-and-storageiterator)
4. [Designing BoltDB's Buckets: blocks, utxo, meta](#4-designing-boltdbs-buckets-blocks-utxo-meta)
5. [Implementing `OpenBoltStore`](#5-implementing-openboltstore)
6. [Implementing `PutBlock` and `GetBlock`](#6-implementing-putblock-and-getblock)
7. [Implementing `PutUTXO`, `GetUTXO`, and `DeleteUTXO`](#7-implementing-pututxo-getutxo-and-deleteutxo)
8. [Implementing the Block Iterator](#8-implementing-the-block-iterator)
9. [Hands-On: Storing and Retrieving a Small Chain](#9-hands-on-storing-and-retrieving-a-small-chain)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Recap: What Chapters 53 and 54 Established

Chapter 53 identified three concrete failures of the Chapter 20 flat file: a crash mid-write can corrupt the file, finding one transaction (or one balance) means scanning every block ever mined, and two goroutines appending at once can silently interleave their writes. Chapter 54 explained why an *embedded database* — a library compiled into your own process, with no separate server — fixes all three at once, and chose **BoltDB** (via its maintained fork, `go.etcd.io/bbolt`) as GoChain's storage engine: single-writer, ACID transactions, one file on disk, organized into **buckets** — independent, sorted namespaces, like drawers in a filing cabinet.

What Chapter 54 did *not* do is decide how GoChain's own types — `*core.Block`, `*core.TxOutput` — map onto BoltDB's raw byte-key, byte-value world. That mapping is this chapter's job, and it starts with a decision that matters more than the mapping itself: GoChain's `core` package will never import `go.etcd.io/bbolt` directly.

---

## 2. Why Store Is an Interface, Not `*bolt.DB`

Imagine, for a moment, that `core.Blockchain` held a `*bolt.DB` field and called `db.Update(...)` and `tx.Bucket(...)` directly, scattered across a dozen methods. That code would work. It would also mean: every test for `core.Blockchain` needs a real BoltDB file on disk (slow, and awkward to run in parallel). Swapping to Badger — the LSM-tree alternative from Chapter 54, the honest answer for a much higher-throughput production node — would mean rewriting every one of those call sites, in every package that ever touched storage, by hand. And a lightweight, in-memory fake for fast unit tests would not exist at all, because nothing defines what "a fake" would even need to implement.

An **interface** solves this by naming the *behavior* GoChain's other packages actually need — "store a block by hash," "fetch a UTXO by key" — without naming *how* that behavior is provided. This is the exact same instinct you used for `consensus.Engine` back in Volume 4: `core.Blockchain` did not care whether proof of work or, later, proof of stake solved the puzzle, only that *something* implementing `consensus.Engine` did. Here, nothing above the `storage` package will care whether BoltDB, Badger, or an in-memory map answers `GetBlock` — only that *something* implementing `storage.Store` does.

```
WITHOUT AN INTERFACE                          WITH storage.Store
(core imports go.etcd.io/bbolt directly)      (core imports only gochain/storage)

core.Blockchain ----> *bolt.DB                core.Blockchain ----> storage.Store (interface)
     |                                                                    ^
     | every method calls                                                |
     | db.Update(), tx.Bucket()                          +---------------+---------------+
     v                                                    |                               |
 tightly coupled to BoltDB's                        storage.BoltStore              storage.badgerStore
 exact API, forever                                 (this chapter)                 (a future chapter,
                                                                                     zero changes to core)
```

This is not a hypothetical concern reserved for "real" production codebases. It is the concrete reason Chapter 58's mini project can build a fast, throwaway test suite for `UTXOSet` without ever touching a real `.db` file on disk, and the concrete reason a future contributor could hand you a Badger-backed `Store` implementation as a self-contained pull request that touches exactly one new file.

---

## 3. Designing `storage.Store` and `storage.Iterator`

With that motivation settled, here is the interface itself — five methods for reading and writing individual blocks and UTXOs, plus one method for walking the whole chain in order:

```go
// storage/store.go
package storage

import "github.com/you/gochain/core"

// Store is the storage contract every other GoChain package depends on.
// core, consensus, and wallet never import a specific database library —
// they call these six methods and nothing else.
type Store interface {
	// PutBlock saves a block under its hash. Calling it again with the same
	// hash overwrites the previous value — blocks are content-addressed, so
	// in practice the bytes never actually change, but the interface does
	// not forbid it.
	PutBlock(hash []byte, block *core.Block) error

	// GetBlock returns the block stored under hash, or an error if no block
	// with that hash has ever been saved.
	GetBlock(hash []byte) (*core.Block, error)

	// PutUTXO saves an unspent transaction output under key — by convention,
	// "<txID-in-hex>:<output-index>", established in Chapter 56.
	PutUTXO(key string, output *core.TxOutput) error

	// GetUTXO returns the output stored under key, or an error if that
	// output does not exist (either it was never created, or it has
	// already been spent and removed by DeleteUTXO).
	GetUTXO(key string) (*core.TxOutput, error)

	// DeleteUTXO removes an output — called the moment some transaction's
	// input spends it, so the UTXO set only ever tracks unspent outputs.
	DeleteUTXO(key string) error

	// Iterator returns a fresh Iterator positioned at the chain's current
	// tip, ready to walk backward to the genesis block.
	Iterator() Iterator
}

// Iterator walks a chain of blocks one at a time, from the newest block
// back to genesis, the same direction core.Blockchain has always used.
type Iterator interface {
	// Next returns the current block and advances the iterator to that
	// block's predecessor. Calling it past genesis is an error.
	Next() (*core.Block, error)

	// HasNext reports whether Next can still be called — false once the
	// iterator has already returned the genesis block.
	HasNext() bool
}
```

Two design choices are worth calling out explicitly, because they are easy to gloss over on a first read. First, every method that can fail returns an `error` — even `GetBlock` and `GetUTXO`, whose "not found" case is genuinely different from "the database is corrupted" or "the disk is full," but a caller can distinguish those by inspecting the error's text or, in a more defensive codebase, by defining a sentinel `ErrNotFound` (left as Exercise 6). Second, `Iterator()` returns a *new* iterator every time it is called, rather than one shared cursor — two goroutines can each walk the chain independently, at their own pace, without stepping on each other, which matters the moment Volume 10's API and Volume 7's sync logic both want to read the chain at the same time.

---

## 4. Designing BoltDB's Buckets: blocks, utxo, meta

`storage.Store` says nothing about buckets — that is deliberately an implementation detail of `BoltStore`, invisible from outside the `storage` package. Internally, `BoltStore` uses three buckets, each with a narrow, single-purpose job:

```
                       blockchain.db (one physical BoltDB file)
   +--------------------------------------------------------------------+
   |                         FILING CABINET                             |
   |  +---------------------+  +---------------------+  +-------------+ |
   |  |  DRAWER: blocks     |  |  DRAWER: utxo        |  | DRAWER: meta| |
   |  |---------------------|  |---------------------- |  |------------| |
   |  |  key:   block hash  |  |  key:  "txid:index"  |  |  key: "tip"| |
   |  |  value: gob(Block)  |  |  value: gob(TxOutput)|  |  value:    | |
   |  |                     |  |                       |  |   tip hash | |
   |  |  one entry per      |  |  one entry per        |  +-------------+ |
   |  |  mined/received     |  |  currently-unspent    |                 |
   |  |  block, forever     |  |  output, added and    |                 |
   |  |                     |  |  removed as coins     |                 |
   |  |                     |  |  move                 |                 |
   |  +---------------------+  +---------------------+                  |
   +--------------------------------------------------------------------+
```

The **blocks** bucket is a permanent, append-mostly log: `hash -> gob(Block)`. Nothing ever gets deleted from it — a block, once mined and accepted, is part of history forever (barring the chain reorganizations Volume 7 handles separately). The **utxo** bucket is the opposite: a live, constantly-churning index that Chapter 56 builds and maintains, where an entry's *presence* means "spendable right now" and `DeleteUTXO` is called just as often as `PutUTXO`. The **meta** bucket is the smallest of the three — right now it holds exactly one key, `"tip"`, whose value is the hash of the most recently added block, which is precisely what lets `Iterator()` know where to start walking without needing a separate index of "which block is newest."

Splitting these into separate buckets, rather than one shared bucket with prefixed keys (`"block:"` + hash vs. `"utxo:"` + key), is not just tidiness. BoltDB keeps each bucket's keys sorted independently as its own B-tree (Chapter 54, Section 2), so a bucket with a smaller, more uniform key space stays a shallower, faster tree — and, just as importantly, it means a bug that accidentally reuses a key string in one bucket can never collide with an unrelated key in another.

---

## 5. Implementing `OpenBoltStore`

`BoltStore` wraps a `*bolt.DB` and satisfies `storage.Store`. Opening one creates the file (if it does not exist yet) and makes sure all three buckets exist, using the same `CreateBucketIfNotExists` idempotence Chapter 54 introduced — safe to call on every single startup, first run or the hundredth.

```go
// storage/bolt_store.go
package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"

	"github.com/you/gochain/core"
)

var (
	blocksBucket = []byte("blocks")
	utxoBucket   = []byte("utxo")
	metaBucket   = []byte("meta")
	tipKey       = []byte("tip")
)

// BoltStore is the BoltDB-backed implementation of storage.Store.
type BoltStore struct {
	db *bolt.DB
}

// OpenBoltStore opens (or creates) a BoltDB file at path and makes sure the
// blocks, utxo, and meta buckets all exist before returning.
func OpenBoltStore(path string) (*BoltStore, error) {
	// A one-second open timeout keeps a stuck process from hanging forever
	// if some other process already holds the file's exclusive lock —
	// BoltDB enforces one writer per file at the OS level, across processes,
	// not just across goroutines within one process.
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt db at %s: %w", path, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{blocksBucket, utxoBucket, metaBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize buckets: %w", err)
	}

	return &BoltStore{db: db}, nil
}

// Close releases the file lock and flushes any pending writes. Not part of
// the Store interface — callers that hold a concrete *BoltStore (typically
// just main.go, at startup) are responsible for calling it on shutdown.
func (s *BoltStore) Close() error {
	return s.db.Close()
}

// encodeBlock and decodeBlock mirror the gob framing Chapter 53 used for the
// flat file — the encoding hasn't changed, only where the bytes end up.
func encodeBlock(b *core.Block) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return nil, fmt.Errorf("gob-encode block: %w", err)
	}
	return buf.Bytes(), nil
}

func decodeBlock(data []byte) (*core.Block, error) {
	var b core.Block
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&b); err != nil {
		return nil, fmt.Errorf("gob-decode block: %w", err)
	}
	return &b, nil
}
```

Notice what did *not* need to change from Chapter 53: blocks are still gob-encoded, byte for byte. What changed is everywhere those bytes end up living — instead of an append-only stream with no index, they land in a bucket keyed by hash, addressable directly, and wrapped in a transaction that either commits completely or leaves the previous state untouched.

---

## 6. Implementing `PutBlock` and `GetBlock`

`PutBlock` does two things inside one atomic write transaction: save the block's bytes under its hash, and update `meta["tip"]` if (and only if) this block is now the tallest one the store has ever seen. Doing both inside the same `db.Update` call matters — if the process crashes between "save the block" and "update the tip," BoltDB's atomicity guarantees that either *both* changes are visible after restart, or *neither* is. There is no window where the block exists but the tip pointer has gone stale, or vice versa.

```go
// storage/bolt_store.go (continued)

func (s *BoltStore) PutBlock(hash []byte, block *core.Block) error {
	encoded, err := encodeBlock(block)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		blocks := tx.Bucket(blocksBucket)
		if err := blocks.Put(hash, encoded); err != nil {
			return fmt.Errorf("put block %x: %w", hash, err)
		}

		meta := tx.Bucket(metaBucket)
		currentTip := meta.Get(tipKey)
		if currentTip == nil {
			// First block ever stored (genesis) — it is the tip by definition.
			return meta.Put(tipKey, hash)
		}

		tipBlock, err := decodeBlock(blocks.Get(currentTip))
		if err != nil {
			return fmt.Errorf("decode current tip block: %w", err)
		}
		if block.Height > tipBlock.Height {
			return meta.Put(tipKey, hash)
		}
		// A shorter or equal-height block (e.g. one arriving late from a
		// peer during sync) is stored, but does not become the new tip.
		return nil
	})
}

func (s *BoltStore) GetBlock(hash []byte) (*core.Block, error) {
	var block *core.Block
	err := s.db.View(func(tx *bolt.Tx) error {
		encoded := tx.Bucket(blocksBucket).Get(hash)
		if encoded == nil {
			return fmt.Errorf("block %x not found", hash)
		}
		decoded, err := decodeBlock(encoded)
		if err != nil {
			return err
		}
		block = decoded
		return nil
	})
	if err != nil {
		return nil, err
	}
	return block, nil
}
```

`GetBlock` uses `db.View`, BoltDB's read-only transaction — which, per Chapter 54, can run concurrently with any number of other `View` calls, including while a `PutBlock` from another goroutine is queued up waiting for its turn to write. This is the "safe concurrent access" property Chapter 53 asked for, delivered entirely by BoltDB's own locking, with zero mutexes written by us.

---

## 7. Implementing `PutUTXO`, `GetUTXO`, and `DeleteUTXO`

The UTXO methods follow the identical shape, against the `utxo` bucket instead of `blocks`, with `string` keys converted to `[]byte` for BoltDB's API (which only ever deals in byte slices):

```go
// storage/bolt_store.go (continued)

func (s *BoltStore) PutUTXO(key string, output *core.TxOutput) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(output); err != nil {
		return fmt.Errorf("gob-encode utxo %s: %w", key, err)
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(utxoBucket).Put([]byte(key), buf.Bytes())
	})
}

func (s *BoltStore) GetUTXO(key string) (*core.TxOutput, error) {
	var output *core.TxOutput
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(utxoBucket).Get([]byte(key))
		if data == nil {
			return fmt.Errorf("utxo %s not found", key)
		}
		var out core.TxOutput
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&out); err != nil {
			return fmt.Errorf("gob-decode utxo %s: %w", key, err)
		}
		output = &out
		return nil
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *BoltStore) DeleteUTXO(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(utxoBucket).Delete([]byte(key))
	})
}
```

`Delete` on a BoltDB bucket is not an error if the key never existed — it simply does nothing, which is exactly the behavior `UTXOSet.Update` in Chapter 56 wants when it unconditionally tries to remove every output an incoming transaction's inputs reference.

---

## 8. Implementing the Block Iterator

`Iterator()` needs to answer "give me the current tip, then let me walk backward one block at a time." The tip's hash lives in the `meta` bucket; each subsequent block's predecessor is its own `PrevBlockHash` field, the same field that has linked GoChain's chain together since Chapter 18.

```go
// storage/bolt_store.go (continued)

// blockIterator implements storage.Iterator against a *BoltStore.
type blockIterator struct {
	store       *BoltStore
	currentHash []byte
}

func (s *BoltStore) Iterator() Iterator {
	var tip []byte
	s.db.View(func(tx *bolt.Tx) error {
		// Copy the bytes out — BoltDB's Get() returns a slice that is only
		// valid for the lifetime of this transaction, and this transaction
		// ends the moment View() returns.
		if v := tx.Bucket(metaBucket).Get(tipKey); v != nil {
			tip = append([]byte(nil), v...)
		}
		return nil
	})
	return &blockIterator{store: s, currentHash: tip}
}

func (it *blockIterator) HasNext() bool {
	return len(it.currentHash) > 0 && !isZeroHash(it.currentHash)
}

func (it *blockIterator) Next() (*core.Block, error) {
	block, err := it.store.GetBlock(it.currentHash)
	if err != nil {
		return nil, fmt.Errorf("iterator: %w", err)
	}
	it.currentHash = block.PrevBlockHash
	return block, nil
}

// isZeroHash reports whether hash is the genesis block's conventional
// all-zero PrevBlockHash, established back in Chapter 18 — the signal that
// there is no predecessor left to visit.
func isZeroHash(hash []byte) bool {
	for _, b := range hash {
		if b != 0 {
			return false
		}
	}
	return true
}
```

The copy in `Iterator()` — `append([]byte(nil), v...)` — is not decoration. BoltDB's memory-mapped design means a `[]byte` returned from inside a transaction points directly into the mapped file; once that transaction ends, the mapping can be reused, and holding onto the unsafe slice afterward is undefined behavior. Copying the bytes out before the `View` call returns is the correct, idiomatic pattern any time a value needs to outlive the transaction that produced it — the same reason `decodeBlock` fully deserializes into a fresh `core.Block` rather than trying to hold onto the raw bytes.

---

## 9. Hands-On: Storing and Retrieving a Small Chain

Put the whole thing together with a short program that opens a store, saves three blocks, and walks them back out with the iterator:

```go
// cmd/storagedemo/main.go
package main

import (
	"fmt"
	"log"

	"github.com/you/gochain/core"
	"github.com/you/gochain/storage"
)

func main() {
	store, err := storage.OpenBoltStore("gochain.db")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer store.Close()

	genesis := core.NewGenesisBlock()
	if err := store.PutBlock(genesis.Hash, genesis); err != nil {
		log.Fatalf("put genesis: %v", err)
	}

	prev := genesis
	for i := 0; i < 2; i++ {
		next := core.NewBlock(prev.Height+1, nil, prev.Hash)
		if err := store.PutBlock(next.Hash, next); err != nil {
			log.Fatalf("put block %d: %v", next.Height, err)
		}
		prev = next
	}

	fmt.Println("Chain, tip to genesis:")
	it := store.Iterator()
	for it.HasNext() {
		block, err := it.Next()
		if err != nil {
			log.Fatalf("iterate: %v", err)
		}
		fmt.Printf("  height=%d hash=%x\n", block.Height, block.Hash)
	}
}
```

```
$ go run ./cmd/storagedemo
Chain, tip to genesis:
  height=2 hash=7f3a1c...e91b
  height=1 hash=4b8d02...c317
  height=0 hash=000000...0000
```

Run the program a second time without deleting `gochain.db` first and watch `PutBlock` happily add three *more* blocks on top — a small reminder that `storage.Store` alone does not stop you from calling it incorrectly; it is `core.Blockchain`'s job, from Volume 3 onward, to decide *when* a new block should be appended. This chapter's contribution is making that append safe, indexed, and concurrent by construction, exactly as Chapter 53 asked for.

---

## Summary

- `storage.Store` is an interface with six methods — `PutBlock`, `GetBlock`, `PutUTXO`, `GetUTXO`, `DeleteUTXO`, and `Iterator` — that every other GoChain package depends on instead of depending on BoltDB directly.
- The interface-first design means a future Badger-backed store, or an in-memory fake for fast tests, is a self-contained new file in the `storage` package, not a rewrite of `core`, `consensus`, or anything else.
- `BoltStore` organizes GoChain's data into three buckets: `blocks` (a permanent hash-to-block log), `utxo` (a constantly-updated index of spendable outputs), and `meta` (currently just the chain's tip hash).
- `PutBlock` saves a block and, inside the *same* atomic transaction, updates the tip pointer if the new block is taller — so a crash can never leave the tip pointing at a block that was never actually saved.
- `PutUTXO`, `GetUTXO`, and `DeleteUTXO` map GoChain's `"txid:index"` string keys onto the `utxo` bucket's byte-slice keys, gob-encoding `*core.TxOutput` the same way blocks are encoded.
- The block iterator walks from the tip backward via each block's `PrevBlockHash`, stopping at the genesis block's conventional all-zero predecessor hash.
- Bytes read out of a BoltDB transaction must be copied before the transaction ends, since they point directly into memory-mapped file pages that become invalid once the transaction closes.

---

## Exercises

### Easy

1. Add a `Has(hash []byte) bool` helper method to `*BoltStore` (not part of the `Store` interface) that reports whether a block with that hash has been stored, without decoding the block's full contents. Use `tx.Bucket(blocksBucket).Get(hash) != nil`.
2. Modify the hands-on demo in Section 9 to print the total number of blocks visited by the iterator, using a counter incremented inside the `for it.HasNext()` loop.
3. `PutBlock`'s tip-update logic uses `block.Height > tipBlock.Height`. Explain, in a comment, what would go wrong if this were `>=` instead, in the specific case of Chapter 53's network scenario where a mining goroutine and a network-sync goroutine each try to store a different block at the same height.

### Medium

4. Write a table-driven test, `TestBoltStore_PutGetBlock`, that opens a `BoltStore` in a temporary directory (`t.TempDir()`), stores three blocks with different heights in a random order, and asserts `GetBlock` returns each one correctly and that `Iterator()` walks them tallest-first regardless of insertion order.
5. `GetBlock` currently returns a generic `fmt.Errorf` when a hash is not found. Define a sentinel error, `var ErrNotFound = errors.New("storage: not found")`, wrap it with `%w` in `GetBlock` and `GetUTXO`, and write a test that uses `errors.Is` to detect it specifically, distinct from a decode failure.
6. Add a `PutBlocks(blocks []*core.Block) error` convenience method to `*BoltStore` that writes an entire slice of blocks inside a single `db.Update` transaction, rather than one transaction per block. Explain, using Chapter 54 Section 4's copy-on-write explanation, why batching writes like this into fewer transactions is measurably faster than one transaction per block.

### Hard

7. Design and implement an in-memory `Store` (backed by Go maps and a `sync.RWMutex`, no BoltDB involved) that satisfies the same `storage.Store` interface. Run any tests you wrote in Exercise 4 against both `BoltStore` and your in-memory store by parameterizing the test over a `func() storage.Store` factory — this is exactly the payoff the interface promised in Section 2.
8. `BoltStore.Iterator()` reads the tip once, at creation time. Write a test that creates an iterator, then calls `PutBlock` with a new, taller block on the *same* store before finishing the iteration. Does the in-progress iterator see the new block? Explain why, referencing BoltDB's copy-on-write page model from Chapter 54.
9. Research `bolt.DB.Backup` (or the equivalent `Tx.WriteTo`) and write a small `Backup(destPath string) error` method on `*BoltStore` that copies the entire database to a new file while it remains open and in use. Explain, in a short comment, why this operation is itself safe to run concurrently with ongoing reads and writes, given everything Chapter 54 explained about BoltDB's internals.
