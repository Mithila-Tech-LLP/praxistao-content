# Chapter 54: Embedded Databases — BoltDB and Badger

Chapter 53 showed exactly what goes wrong with a flat file. This chapter introduces the tool that fixes all three problems at once: an embedded database. We compare the two most popular pure-Go options — BoltDB and Badger — and settle on which one GoChain will use for the rest of this course.

## Table of Contents

1. [What "Embedded Database" Means](#1-what-embedded-database-means)
2. [B-Trees and LSM-Trees, Conceptually](#2-b-trees-and-lsm-trees-conceptually)
3. [BoltDB: Simple, Single-Writer, Battle-Tested](#3-boltdb-simple-single-writer-battle-tested)
4. [How BoltDB Guarantees Crash Safety](#4-how-boltdb-guarantees-crash-safety)
5. [Badger: Built for Write Throughput](#5-badger-built-for-write-throughput)
6. [Choosing BoltDB for GoChain](#6-choosing-boltdb-for-gochain)
7. [Buckets as Filing Cabinet Drawers](#7-buckets-as-filing-cabinet-drawers)
8. [Hands-On: A BoltDB Hello World](#8-hands-on-a-boltdb-hello-world)
9. [Updating, Deleting, and Iterating Keys](#9-updating-deleting-and-iterating-keys)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What "Embedded Database" Means

When most people hear "database," they picture something like PostgreSQL or MySQL: a separate server process, running independently, that your application talks to over a network connection (even if that connection is just `localhost`). You install it separately, start it as its own service, and your Go program is a *client* of it.

An **embedded database** is a different shape entirely: it is a Go library, compiled directly into your program, with no separate server, no network connection, and no separate process to install or manage. It reads and writes a single file (or a small folder of files) on disk, directly, from inside your own process's memory space. Think of the difference the way you would think about a filing cabinet in your own office versus calling a records-management company: with the filing cabinet, there is no phone call, no waiting for someone else's system to respond — you open the drawer yourself.

```
CLIENT-SERVER DATABASE                       EMBEDDED DATABASE
(e.g. PostgreSQL)                            (e.g. BoltDB, Badger)

+-------------+     network      +--------+   +--------------------------+
| Your Go     | ---------------> |Postgres|   | Your Go process          |
| process     | <--------------- | server |   |  +--------------------+  |
+-------------+                  +--------+   |  | BoltDB / Badger    |  |
      |                              |        |  | (just a library,   |  |
      v                              v        |  |  no separate       |  |
  runs on your                runs as its      |  |  process)          |  |
  machine                     own service,     |  +--------------------+  |
                              installed        |         |                |
                              separately        +--------------------------+
                                                          |
                                                          v
                                                  one file on disk
```

Despite having no server to manage, an embedded database still gives you the properties Chapter 53 asked for: **crash-safe writes** (a write either fully commits or has no effect, even across a crash), **indexed lookups** (get a value by key directly, without scanning anything else), and **safe concurrent access** (built into the library, not something every caller reinvents). Go has two dominant, mature, pure-Go embedded key-value databases that provide exactly this: **BoltDB** and **Badger**. Both let you store arbitrary byte-slice keys mapped to arbitrary byte-slice values — which is all GoChain actually needs, since we already know how to serialize blocks, transactions, and UTXOs into bytes.

---

## 2. B-Trees and LSM-Trees, Conceptually

BoltDB and Badger differ in one fundamental way: the internal data structure each uses to keep keys sorted and searchable on disk. You do not need to implement either structure yourself to use these databases — but understanding the shape of each explains *why* the two libraries behave so differently under load, and why GoChain picks the one it does.

**A B-tree** (the structure BoltDB uses) keeps keys sorted in a tree of fixed-size disk pages, where every page has several children, arranged so the tree stays short and wide rather than tall and thin. Reading a key means walking down from the root, comparing your key against a handful of keys per page, and following the branch that must contain it — usually only 3-4 disk pages need to be touched even for millions of keys. Writing a key means finding the right leaf page, inserting it in sorted order, and occasionally splitting a page that has become full. This makes B-trees excellent for reads, because the tree is always kept in a directly-searchable, sorted shape.

```
B-TREE (BoltDB) — always kept sorted and searchable
                        [ M ]
                      /       \
              [ D, H ]         [ R, W ]
             /   |    \       /    |    \
         [A,B] [E,F] [I,J]  [N,O] [S,T] [X,Y]

Every write updates the tree IN PLACE, keeping it sorted at all times.
A read walks straight down: compare, branch, compare, branch, done.
```

**An LSM-tree** ("Log-Structured Merge-tree", the structure Badger uses) takes the opposite approach: instead of updating a sorted structure in place on every write, it buffers recent writes in memory (in a structure often called a "memtable"), and periodically flushes that buffer to disk as a new, immutable, sorted file. Over time, many such files accumulate, and a background process ("compaction") periodically merges them back together, discarding overwritten or deleted keys. Writing is extremely fast, because it is just an in-memory append followed by an occasional sequential disk write — no searching for the right spot, no rewriting existing pages. Reading is more work, because a key might exist in the in-memory buffer, or in any of several on-disk sorted files, so a read may need to check multiple places.

```
LSM-TREE (Badger) — writes buffer in memory, flush to disk as sorted "runs"
memory:  [ new writes buffer here, fast in-memory append ]
                              |
                      (buffer fills up, flush)
                              v
disk:   [sorted run 3 (newest)] [sorted run 2] [sorted run 1 (oldest)]
                                                          |
                                        background compaction periodically
                                        merges runs, drops overwritten keys

A read may have to check the memory buffer, then run 3, then run 2, etc.,
until it finds the key (or confirms it is not there).
```

Put simply: **B-trees favor read speed and keep the on-disk structure simple at all times; LSM-trees favor write speed by deferring the cost of keeping things tidy to a background compaction process.** Neither is "better" in the abstract — they are different trade-offs, and which one wins depends on your workload's actual read/write balance.

---

## 3. BoltDB: Simple, Single-Writer, Battle-Tested

**BoltDB** (the package GoChain uses is its actively maintained fork, `go.etcd.io/bbolt`, since the original `boltdb/bolt` repository is archived) is a pure-Go, embedded, B-tree-based key-value store. Its design is deliberately narrow, and that narrowness is a feature, not a limitation for our purposes:

- **One writer at a time.** BoltDB allows any number of concurrent *readers*, but only one write transaction may be open at any moment — enforced by the library itself, so you never have to build your own locking scheme to avoid the exact concurrent-write corruption from Chapter 53.
- **ACID transactions.** Every read and every write happens inside a transaction. A write transaction either commits completely (all changes durable on disk) or, if anything goes wrong, none of its changes take effect at all — solving Chapter 53's crash-corruption problem at the storage-engine level, for free.
- **A single file.** The entire database lives in one `.db` file, using memory-mapped I/O, which makes backup as simple as copying that one file (safely, using BoltDB's own `Backup` support, while the database is open).
- **Simplicity over raw throughput.** Because every write is a full B-tree update with a transaction commit, BoltDB's write throughput is lower than an LSM-tree design under heavy sustained writes. For GoChain's actual workload — mining or receiving a new block every few seconds, not thousands of writes per second — this ceiling is nowhere close to being a real constraint.

BoltDB has been used in production for years, notably as the storage engine underneath `etcd` (the distributed key-value store Kubernetes itself is built on), and its API is small enough to read in an afternoon: everything happens through **buckets** (independent namespaces, like separate tables) accessed inside **transactions**.

---

## 4. How BoltDB Guarantees Crash Safety

It is worth being concrete about *how* BoltDB delivers on the crash-safety promise from Section 3, because "it's ACID" can otherwise feel like a magic word. BoltDB uses two techniques together, and both map directly onto Chapter 53's corruption scenario.

The first is **memory-mapped I/O**: instead of the usual pattern of `Read`/`Write` calls copying bytes between your program and the OS's file-system buffers, BoltDB maps the entire database file directly into your process's address space with `mmap`. Reading a page of the B-tree is then just reading memory — no system call per read, which is a large part of why BoltDB reads are fast.

The second, more important technique is **copy-on-write pages combined with a single atomic pointer swap**. When a write transaction modifies the B-tree, BoltDB never overwrites a page that is part of the *current*, committed tree in place. Instead, it writes the *new* version of any changed page to a fresh location in the file, leaving the old page untouched. Only once every changed page has been written out and flushed to disk does BoltDB perform the very last step of a commit: updating a single "meta page" that records the root of the tree — and that one small write is what atomically makes the new tree the current one.

```
BEFORE the write transaction commits:
  meta page --> points to OLD root --> OLD tree pages (unchanged, still valid)

DURING the write transaction:
  new/changed pages are written to FRESH locations in the file
  (the old pages are never touched — they still form a complete, valid tree)

THE COMMIT ITSELF:
  meta page is overwritten to point to the NEW root — one small write

IF THE PROCESS CRASHES BEFORE THIS LAST WRITE:
  meta page still points to the OLD root — the new, half-written pages
  are simply ignored the next time the file is opened. The database looks
  exactly as it did before the transaction ever started. Nothing is corrupt,
  because nothing that was already committed was ever overwritten in place.
```

This is exactly the "half-finished write" scenario from Chapter 53, Section 2 — except here it cannot corrupt anything, because the *old*, fully-valid tree is never modified in place while a new one is being built. A crash at any point before the final meta-page write simply leaves the old, complete tree as the current one; a crash after it leaves the new, complete tree as the current one. There is no window in between where a reader could observe a half-written state, which is precisely the guarantee our flat file could not make.

---

## 5. Badger: Built for Write Throughput

**Badger** (`github.com/dgraph-io/badger`) is also a pure-Go, embedded key-value store, but it is built on an LSM-tree design specifically to handle much higher sustained write volume than BoltDB — its own benchmarks show write throughput many times higher than BoltDB under heavy concurrent write load, at the cost of a larger, more complex on-disk footprint (multiple SST files plus a value log, rather than one single file) and periodic background compaction work competing for CPU and disk I/O.

Badger is the right choice when your workload genuinely writes a very large volume of data very frequently — it was built by Dgraph specifically for their graph database's needs, and it is a completely reasonable choice for a blockchain node that needs to ingest a firehose of transactions (a busy mempool receiving thousands of transactions per second, for instance, or a full archival node re-indexing years of history as fast as possible).

```
                       BoltDB (B-tree)              Badger (LSM-tree)
Write pattern      in-place, sorted updates    buffer + append, background merge
Write throughput   good, single-writer         high, built for heavy concurrent writes
Read pattern       always one sorted tree      may check memory + multiple sorted files
On-disk footprint  one file                    multiple SST files + a value log
Concurrency model  many readers, 1 writer      many readers, many writers
Operational shape  simplest to reason about    more moving parts (compaction, GC)
```

---

## 6. Choosing BoltDB for GoChain

For GoChain, we choose **BoltDB** as the storage engine for the rest of this course, for three concrete reasons that matter more to a learning project — and, frankly, to most real single-node blockchain clients — than raw write throughput:

- **GoChain's actual write rate is low.** A new block arrives every few seconds at most (Volume 4's difficulty target), not thousands of times per second. BoltDB's single-writer, B-tree design has no trouble keeping up with "one write transaction every few seconds," even with a large batch of transactions inside each one.
- **Its simplicity matches this course's goal of understanding every layer.** A single file, one writer at a time, and a small, direct API make it far easier to reason about exactly what is happening on disk — which matters enormously in a course whose entire purpose is building genuine understanding, not just wiring up a fast black box.
- **Its guarantees are exactly what Chapter 53 asked for.** ACID transactions solve the crash-corruption problem directly, and buckets plus B-tree indexing solve the full-scan problem directly, with no extra design work required from us.

Badger remains the honest, name-it-explicitly answer to "what would you reach for in production, at scale, with a much higher transaction volume than a teaching project generates?" — and because Chapter 55 defines `storage.Store` as an interface rather than coupling GoChain directly to BoltDB's API, swapping in a Badger-backed implementation later is a contained, one-package change, not a rewrite. We name that design decision now precisely so it is not a surprise later.

---

## 7. Buckets as Filing Cabinet Drawers

BoltDB organizes all of its keys and values into **buckets** — independent, named namespaces, roughly analogous to tables in a relational database, or to separate drawers in one filing cabinet. Two different buckets can have keys with the exact same bytes without colliding, because a key is only ever looked up *within* a specific bucket.

```
                      blockchain.db (one physical file)
        +------------------------------------------------------+
        |                    FILING CABINET                    |
        |  +------------------+  +------------------+          |
        |  |  DRAWER: blocks  |  |  DRAWER: utxo    |          |
        |  |  ---------------|  |  ---------------- |          |
        |  |  hash -> block  |  |  "txid:idx" ->    |          |
        |  |  (sorted by key,|  |     TxOutput       |          |
        |  |   B-tree inside)|  |  (sorted by key)   |          |
        |  +------------------+  +------------------+          |
        |             +------------------+                     |
        |             |  DRAWER: meta    |                     |
        |             |  ---------------- |                     |
        |             |  "tip" -> last   |                     |
        |             |    block hash    |                     |
        |             +------------------+                     |
        +------------------------------------------------------+
```

Each drawer keeps its own contents sorted and independently searchable — exactly the structure Chapter 55 will use to design GoChain's `blocks`, `utxo`, and `meta` buckets.

---

## 8. Hands-On: A BoltDB Hello World

Before designing anything GoChain-specific, get comfortable with BoltDB's raw API: open a database, create a bucket, put a key, and get it back.

```bash
go get go.etcd.io/bbolt
```

```go
// cmd/boltdemo/main.go
package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	// Open (or create, if it doesn't exist) a single database file.
	// 0600 permissions: only the owning user can read or write it.
	db, err := bolt.Open("hello.db", 0600, nil)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Every write goes through a write transaction. Update() runs the given
	// function inside one, and automatically commits it if the function
	// returns nil, or rolls it back entirely if it returns an error.
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("greetings"))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		return bucket.Put([]byte("hello"), []byte("world"))
	})
	if err != nil {
		log.Fatalf("write: %v", err)
	}

	// Reads go through a read-only transaction. View() can run concurrently
	// with other View() calls, and never blocks on (or is blocked by) them.
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("greetings"))
		if bucket == nil {
			return fmt.Errorf("bucket does not exist")
		}
		value := bucket.Get([]byte("hello"))
		fmt.Printf("hello -> %s\n", value)
		return nil
	})
	if err != nil {
		log.Fatalf("read: %v", err)
	}
}
```

Running this prints exactly what you would expect:

```
$ go run ./cmd/boltdemo
hello -> world
```

A few details worth naming explicitly, since every one of them reappears in Chapter 55's real implementation. `bolt.Open` takes a file path and creates it if it does not already exist — this single file *is* the entire database. `CreateBucketIfNotExists` is idempotent: calling it every time the program starts is completely safe, and is the standard way to make sure a bucket exists before using it, without needing separate "first run" logic. `db.Update` wraps a function in a write transaction — if that function returns a non-nil error at any point, BoltDB automatically rolls back every change made inside it, so a bucket's contents can never end up half-updated. `db.View` does the same for read-only access, and importantly can run concurrently with other `View` calls from other goroutines, since BoltDB's "single writer, many readers" rule only restricts writes.

Run this program twice in a row and check that the second run still prints `hello -> world` — the file persisted the write across the two separate process runs, which is the entire point.

---

## 9. Updating, Deleting, and Iterating Keys

The hello-world example only ever wrote one key once. Chapter 55's real bucket design needs three more operations constantly: overwriting an existing key, removing one, and walking every key in a bucket in order — all three still go through `db.Update`/`db.View`, just with different bucket methods.

```go
// Overwriting: Put with an existing key simply replaces its value.
// There is no separate "update" method -- Put always means "set this
// key to this value, whether or not it already existed."
err = db.Update(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("greetings"))
	return bucket.Put([]byte("hello"), []byte("world, again"))
})

// Deleting a key that may or may not exist is always safe -- Delete
// never errors just because the key was already absent.
err = db.Update(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("greetings"))
	return bucket.Delete([]byte("hello"))
})

// Iterating every key in a bucket, in sorted order, using a Cursor.
// This is the single most important BoltDB operation for GoChain:
// Chapter 56's UTXO reindexing and Chapter 58's balance index both
// walk an entire bucket exactly this way.
err = db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("greetings"))
	c := bucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("%s -> %s\n", k, v)
	}
	return nil
})
```

Three things are worth calling out. First, `Put` doubles as both "insert" and "update" — there is no separate method, and no error for overwriting an existing key, which matches how a real UTXO set behaves: crediting or updating a balance is just "write this key again with a new value." Second, `Delete` is equally forgiving — deleting a key that is already gone is a no-op, not an error, which matters in Chapter 56 when a block might (in edge cases, like reprocessing during a resync) ask to remove a UTXO that a previous pass already removed. Third, and most important architecturally: `Cursor()` walks keys in **sorted byte order**, not insertion order — this is a direct, visible consequence of BoltDB being a B-tree underneath, and it's exactly why Chapter 56 can implement "give me every UTXO for this address" efficiently by choosing key formats where everything belonging to one address sorts together.

A cursor can also seek directly to a specific starting point rather than always beginning at `First()`:

```go
// Seek jumps straight to the first key >= the given prefix, instead of
// scanning from the very beginning -- useful once Chapter 56 designs
// keys so a specific address's UTXOs all share a common prefix.
err = db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte("greetings"))
	c := bucket.Cursor()
	for k, v := c.Seek([]byte("h")); k != nil && bytes.HasPrefix(k, []byte("h")); k, v = c.Next() {
		fmt.Printf("%s -> %s\n", k, v)
	}
	return nil
})
```

`Seek` combined with a `bytes.HasPrefix` check while advancing is the standard BoltDB pattern for "give me every key starting with X" — precisely the operation Chapter 56 needs to answer "every UTXO belonging to this one address" quickly, without scanning keys belonging to every other address in the same bucket.

---

## Summary

- An embedded database runs as a library inside your own process, with no separate server, storing everything in a single file (or small set of files) and giving you crash-safe, indexed key-value storage without building any of that yourself.
- BoltDB uses a B-tree: writes update a sorted structure in place, favoring fast, simple, always-sorted reads at the cost of lower peak write throughput.
- Badger uses an LSM-tree: writes buffer in memory and flush as sorted, immutable files that get merged later by background compaction, favoring much higher write throughput at the cost of more operational complexity.
- GoChain chooses BoltDB because its write rate (one block every few seconds) is nowhere near BoltDB's ceiling, its single-file, single-writer design is far easier to reason about, and its ACID transactions and B-tree indexing directly solve both of Chapter 53's storage problems.
- Badger remains the honest answer for a much higher-throughput production system, and because GoChain will depend on a `storage.Store` interface rather than BoltDB's concrete API, swapping it in later stays a contained change.
- BoltDB organizes data into buckets — independent namespaces, like drawers in a filing cabinet — each internally sorted and searchable as its own B-tree.
- `db.Update` runs an all-or-nothing write transaction; `db.View` runs a concurrent-safe read-only transaction; both are the foundation Chapter 55 builds `storage.BoltStore` on top of.
- `Put` doubles as insert-or-update and `Delete` is always safe to call on an already-missing key; a `Cursor` walks a bucket in sorted byte order, and `Seek` plus a prefix check is the standard pattern for "every key belonging to one address" that Chapter 56's UTXO lookups depend on.

---

## Exercises

### Easy

1. In your own words, explain the difference between a client-server database and an embedded database, using the filing-cabinet-versus-records-company analogy or one of your own.
2. Modify the hello-world program to store three key-value pairs in the `greetings` bucket instead of one, then use `bucket.ForEach` to print all of them. Consult the `bbolt` documentation for `ForEach`'s signature.
3. What happens if you call `db.View` with a bucket name that was never created? Run the hello-world program against a fresh, empty database file with the `Update` call commented out, and observe (and explain) the result.

### Medium

4. Explain, using the B-tree and LSM-tree diagrams in Section 2, why a Badger-style database needs a background "compaction" process but a BoltDB-style database does not. What would happen to Badger's read performance over time if compaction were disabled entirely?
5. `db.Update` rolls back all changes if the passed function returns an error. Write a small BoltDB program that puts a key, then deliberately returns an error from inside the same `Update` call after the `Put`, and verify (with a follow-up `View`) that the key was never actually persisted.
6. BoltDB allows many concurrent readers but only one writer at a time. Write a short program that starts a long-running read transaction (hold it open with a `time.Sleep`) in one goroutine while attempting a write in another, and time how long the write takes to complete. What does this tell you about how BoltDB's single-writer rule interacts with long-lived read transactions?

### Hard

7. Read the `bbolt` package documentation for `Tx.Cursor()` and use it to iterate over every key in the `greetings` bucket in sorted order, printing each one. Explain why this iteration order is guaranteed (hint: think about what a B-tree's structure implies about key ordering).
8. Research Badger's "value log" design (separate from its LSM-tree index) and explain, in a short paragraph, what problem it solves for large values specifically, and why BoltDB — which stores keys and values together in the same B-tree pages — does not need an equivalent structure for GoChain's block and UTXO data sizes.
9. Benchmark raw BoltDB write throughput on your own machine: write a loop that inserts 10,000 key-value pairs, once with each `Put` in its own `db.Update` transaction, and once with all 10,000 `Put` calls batched inside a single `db.Update` transaction. Report both timings and explain the difference in terms of what a transaction commit actually costs.
