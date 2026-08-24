# Chapter 57: Merkle-Patricia Tries and State

`UTXOSet` answers "what does this address own" quickly. It says nothing about a much stronger question: "can I *prove*, with a single small fingerprint, that the entire UTXO set is exactly what I claim it is — with nothing added, removed, or altered?" That is the question Ethereum's account-state design answers with one data structure, the **Merkle-Patricia trie**, and it is the same question GoChain will need answered the moment Volume 9 lets a contract read and write persistent storage. This chapter builds a simplified Merkle-Patricia trie in Go, computes a single root hash over GoChain's entire UTXO set, and previews exactly where that root hash slots into a block header.

## Table of Contents

1. [Why "Store It in BoltDB" Is Not "Commit to It"](#1-why-store-it-in-boltdb-is-not-commit-to-it)
2. [Patricia Tries: Sharing Prefixes for Free](#2-patricia-tries-sharing-prefixes-for-free)
3. [Merkle Trees, Recapped, and the Hybrid Idea](#3-merkle-trees-recapped-and-the-hybrid-idea)
4. [Designing a Simplified Merkle-Patricia Trie](#4-designing-a-simplified-merkle-patricia-trie)
5. [Implementing the Trie in Go](#5-implementing-the-trie-in-go)
6. [Inserting Keys and Watching the Root Change](#6-inserting-keys-and-watching-the-root-change)
7. [Computing a State Root Over GoChain's UTXO Set](#7-computing-a-state-root-over-gochains-utxo-set)
8. [Where the Root Hash Would Live: the Block Header](#8-where-the-root-hash-would-live-the-block-header)
9. [What Real Ethereum Adds That This Chapter Skips](#9-what-real-ethereum-adds-that-this-chapter-skips)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why "Store It in BoltDB" Is Not "Commit to It"

`UTXOSet` from Chapter 56 is a completely honest, correct index — but it makes no promise about its own integrity. If a single byte in the `utxo` bucket were corrupted, or a malicious operator quietly edited one entry directly in the `.db` file while the node was offline, nothing in `BoltStore` or `UTXOSet` would notice. Compare that to what Chapter 19 built for blocks themselves: every block's hash is computed *from* its contents, so tampering with a transaction inside an old block breaks that block's hash, which breaks every hash after it — tampering is detectable, exactly because the data structure was designed to make it detectable.

The UTXO set has no equivalent property, because it is "just" a key-value bucket, not a chained, hash-linked structure. What we want is a single hash — a **state root** — that changes if, and only if, the underlying set of key-value pairs changes in any way at all, no matter how small. Get that one hash right, and you can put it in a block header, and every node that recomputes the same root from the same UTXO set and gets the same answer has just cryptographically confirmed that everyone agrees on the exact same state, without transmitting the entire state to prove it.

```
WITHOUT a state root:                      WITH a state root:

"Trust that the utxo bucket is correct"    block.StateRoot = 0xa93f...
     no way to verify this compactly            |
     without re-reading the entire set          | anyone can recompute this
                                                 | from their own utxo bucket
                                                 v
                                          if their root != block.StateRoot,
                                          their state has DIVERGED from
                                          consensus — detectable instantly,
                                          without comparing every key
```

---

## 2. Patricia Tries: Sharing Prefixes for Free

A **trie** (the name comes from "re**trie**val," pronounced "try" to avoid confusion with "tree") is a tree where each edge represents one symbol of a key, and following a path from the root spells out a key one symbol at a time. A **Patricia trie** — the name stands for "Practical Algorithm to Retrieve Information Coded in Alphanumeric," though almost nobody calls it that in conversation — is a trie optimized specifically for keys that share long common prefixes, which addresses and hashes do constantly, since they are themselves just long hex strings.

```
Three keys sharing prefixes, stored as plain strings:
   "1a2b3c"
   "1a2b4d"
   "1a9f00"

As a Patricia trie over hex characters (each edge = one hex digit):

                    root
                     |
                    "1a"          <- shared by all three keys, stored ONCE
                   /     \
                "2b"      "9f"
               /    \        \
             "3c"   "4d"     "00"
           (leaf)  (leaf)   (leaf)

The shared prefix "1a" is walked through exactly once, no matter how many
keys start with it — this is the entire point of a Patricia trie over a
plain sorted list or a flat hash map, where each key's full length is
stored independently even when 90% of it duplicates a neighbor's.
```

For GoChain's UTXO keys (`"<txID-hex>:<index>"`) and, more importantly, for Ethereum-style account addresses in general, this matters a great deal in practice: real blockchain state tries hold millions of keys that share long common prefixes purely by chance (any two random hex strings of reasonable length share a few leading characters far more often than intuition suggests), and a Patricia trie's structure exploits that sharing automatically, without anyone designing for it explicitly.

---

## 3. Merkle Trees, Recapped, and the Hybrid Idea

Chapter 10 built a Merkle tree: hash pairs of leaves together, then hash pairs of *those* hashes, repeatedly, up to one Merkle root — and changing any single leaf changes every hash on the path from that leaf to the root, all the way up. That is the tamper-evidence property.

A **Merkle-Patricia trie** fuses the two ideas: it has the Patricia trie's prefix-sharing branch structure, but instead of each node simply holding pointers to its children, every node's own identity *is* a hash — computed from that node's contents and, recursively, from the hashes of everything beneath it. Changing any single key's value anywhere in the trie changes that leaf's hash, which changes its parent's hash (since the parent's hash is computed from its children), which changes *its* parent's hash, all the way up to the root — exactly the propagation Chapter 10's Merkle tree gave you, but layered on top of a structure that also collapses shared prefixes.

```
MERKLE TREE (Chapter 10)              MERKLE-PATRICIA TRIE (this chapter)
- fixed shape: always pairs           - shape follows the KEYS themselves
- leaves are a flat ORDERED list      - leaves are found by walking a
- great for "prove tx #4 is in        prefix path, one symbol at a time
  this ordered block"                 - great for "prove address 0x1a2b...
                                       has exactly this balance, in this
                                       state" — a keyed lookup, not a
                                       position in a list
```

---

## 4. Designing a Simplified Merkle-Patricia Trie

A production-grade Merkle-Patricia trie (the kind `go-ethereum` actually ships) has three distinct node types — **branch** nodes (up to 16 children, one per hex nibble), **extension** nodes (compress a run of single-child branches into one edge labeled with multiple nibbles at once), and **leaf** nodes (hold an actual value) — plus a careful RLP-based encoding for exactly how each node type hashes. That is real, production complexity, and Section 9 names it explicitly so you know what to reach for later. This chapter builds a deliberately simplified version that keeps the two properties that actually matter for understanding *why* this structure exists: **prefix-sharing branches** and **hash propagation to one root** — while dropping the extension-node compression as an explicit simplification (Exercise 7 asks you to add it back).

```
SIMPLIFIED TRIE — every node is a 16-way branch (one child per hex nibble),
any node may ALSO hold a value directly (this chapter merges "branch" and
"leaf" into one node type, unlike real Ethereum's design):

                         root (hash = H(all children's hashes))
                    0..9,a..f children, only some populated
                     /                              \
                  nibble '1'                     nibble '9'
                 node.hash = H(...)              node.hash = H(...)
                 /          \
            nibble 'a'    nibble ...
            node.hash=H(   ...
              value + children)
                 |
            (holds a value here — this is a "leaf" in effect,
             even though it's the same Node type as every branch)
```

---

## 5. Implementing the Trie in Go

Each key is walked one **nibble** (4 bits — half a byte, one hex digit) at a time, so every node has exactly 16 possible children. This is the standard choice for Merkle-Patricia tries specifically because keys are usually already hashes or hex-friendly identifiers, and 16-way branching keeps the trie shallow (a 32-byte hash key is only 64 nibbles deep, worst case) without the wasted space of a 256-way branch (one per full byte).

```go
// trie/trie.go
package trie

import (
	"crypto/sha256"
)

// Node is one node in a simplified Merkle-Patricia trie. Every node has up
// to 16 children — one per hex nibble — and may optionally hold a value,
// merging what a real Merkle-Patricia trie would split into separate
// branch and leaf node types.
type Node struct {
	Children [16]*Node
	Value     []byte // nil unless a key ends exactly at this node
	cachedHash []byte // memoized; cleared by any mutation beneath this node
}

// Trie is a simplified Merkle-Patricia trie: a Patricia-style prefix trie
// over hex nibbles, where every node's hash also commits to everything
// beneath it, Merkle-tree style.
type Trie struct {
	root *Node
}

// New returns an empty trie with a single, valueless root node.
func New() *Trie {
	return &Trie{root: &Node{}}
}

// toNibbles splits key into its hex nibbles — two per byte, high bit first —
// which is what lets the trie branch 16 ways per byte of key instead of 2.
func toNibbles(key []byte) []byte {
	nibbles := make([]byte, len(key)*2)
	for i, b := range key {
		nibbles[i*2] = b >> 4
		nibbles[i*2+1] = b & 0x0f
	}
	return nibbles
}

// Put inserts or overwrites the value stored at key, walking (and creating,
// as needed) one child per nibble.
func (t *Trie) Put(key, value []byte) {
	t.root.put(toNibbles(key), value)
}

func (n *Node) put(nibbles []byte, value []byte) {
	n.cachedHash = nil // this node's hash is now stale, all the way to the root
	if len(nibbles) == 0 {
		n.Value = value
		return
	}
	next := nibbles[0]
	if n.Children[next] == nil {
		n.Children[next] = &Node{}
	}
	n.Children[next].put(nibbles[1:], value)
}

// Get returns the value stored at key, and whether it was found at all.
func (t *Trie) Get(key []byte) ([]byte, bool) {
	return t.root.get(toNibbles(key))
}

func (n *Node) get(nibbles []byte) ([]byte, bool) {
	if len(nibbles) == 0 {
		return n.Value, n.Value != nil
	}
	next := n.Children[nibbles[0]]
	if next == nil {
		return nil, false
	}
	return next.get(nibbles[1:])
}

// RootHash returns the trie's current root hash — a single fingerprint
// that changes if, and only if, ANY key or value anywhere in the trie
// changes. Recomputes only the parts of the trie whose cached hash was
// invalidated by a Put since the last call.
func (t *Trie) RootHash() []byte {
	return t.root.hash()
}

// hash computes (or returns the cached) hash of this node: the hash of its
// own value, if any, concatenated with the hash of every child (or a fixed
// zero value for an absent child, so a node's hash reflects EXACTLY which
// of its 16 slots are populated, not just which are non-empty).
func (n *Node) hash() []byte {
	if n.cachedHash != nil {
		return n.cachedHash
	}

	h := sha256.New()
	if n.Value != nil {
		h.Write([]byte("leaf:"))
		h.Write(n.Value)
	}
	for i, child := range n.Children {
		if child == nil {
			h.Write(zeroHash[:])
			continue
		}
		h.Write([]byte{byte(i)})
		h.Write(child.hash())
	}

	sum := h.Sum(nil)
	n.cachedHash = sum
	return sum
}

var zeroHash = sha256.Sum256(nil)
```

A subtle correctness detail worth naming explicitly: hashing an absent child as a fixed `zeroHash` value, rather than simply skipping it, matters because two *different* sets of populated children could otherwise hash identically if the hasher only ever wrote the children that existed — a node with a child at nibble `3` and nothing else could become indistinguishable from a node with a child at nibble `3` and a *different* child at nibble `7`, if the empty slots contributed nothing to the hash at all and slot positions were not written either. Writing the slot index (`byte(i)`) before each present child's hash closes that gap directly, since every one of the 16 slots is now accounted for in a fixed, position-dependent way.

---

## 6. Inserting Keys and Watching the Root Change

A short demo makes the tamper-evidence property concrete rather than theoretical:

```go
// cmd/trietest/main.go
package main

import (
	"fmt"

	"github.com/you/gochain/trie"
)

func main() {
	t := trie.New()
	t.Put([]byte("addr-alice"), []byte("balance:100"))
	t.Put([]byte("addr-bob"), []byte("balance:50"))
	fmt.Printf("root after 2 inserts:  %x\n", t.RootHash())

	// "Tamper" with Alice's balance — change one value, nothing else.
	t.Put([]byte("addr-alice"), []byte("balance:999"))
	fmt.Printf("root after tampering:  %x\n", t.RootHash())

	// Insert Bob's original value back — does the root return to normal?
	t.Put([]byte("addr-alice"), []byte("balance:100"))
	fmt.Printf("root after reverting:  %x\n", t.RootHash())
}
```

```
$ go run ./cmd/trietest
root after 2 inserts:  9c14af...02e3
root after tampering:  4b7d91...c8a0
root after reverting:  9c14af...02e3
```

The third line is the important one: reverting Alice's balance back to its original value reproduces the *exact same root hash* as before the tamper — because the trie's root is a pure function of its current contents, with no hidden history or ordering dependence. This is a stronger, more useful property than it might look at first glance: it means two independently-built tries, constructed by two different nodes that processed the same set of updates in a completely different order, will still converge on the identical root hash the moment their contents actually match — exactly the property a distributed network of nodes needs to agree "we all have the same state" without comparing every key by hand.

---

## 7. Computing a State Root Over GoChain's UTXO Set

Wiring the trie to GoChain's actual UTXO set is now a matter of walking every entry once and inserting it — reusing the `ForEachUTXO` optional capability Chapter 56 defined:

```go
// storage/state_root.go
package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/you/gochain/core"
	"github.com/you/gochain/trie"
)

// ComputeStateRoot walks every entry currently in the UTXO set and inserts
// it into a fresh trie, keyed by the same "<txID-hex>:<index>" string
// Chapter 56 uses, returning the trie's root hash as a single fingerprint
// over the ENTIRE current UTXO set.
func ComputeStateRoot(store Store) ([]byte, error) {
	iterable, ok := store.(utxoIterable)
	if !ok {
		return nil, fmt.Errorf("compute state root: store does not support iteration")
	}

	t := trie.New()
	err := iterable.ForEachUTXO(func(key string, output *core.TxOutput) error {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(output); err != nil {
			return fmt.Errorf("encode output for key %s: %w", key, err)
		}
		t.Put([]byte(key), buf.Bytes())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("compute state root: %w", err)
	}

	return t.RootHash(), nil
}
```

Run this once against a freshly-`Reindex`ed `UTXOSet` and you get one 32-byte hash that represents the *entire* current set of spendable outputs. Change a single satoshi's worth of value anywhere in that set — spend one output, receive one payment — and recomputing this function produces a completely different root, with no partial resemblance to the old one, exactly the avalanche property Chapter 8 introduced for plain hash functions and Chapter 10 extended to whole trees.

---

## 8. Where the Root Hash Would Live: the Block Header

Here is the forward-looking piece this chapter previews rather than fully wires in: in a design like Ethereum's, every block header carries a `StateRoot` field, computed *after* applying that block's transactions, so the header itself becomes a commitment to "the entire state of the world, as of this block" — not just "the transactions this block happened to include," which is all a Merkle root over transactions (Chapter 10) actually proves.

```
GoChain's block header TODAY (Chapters 16-17):          WHAT A STATE ROOT ADDS:

Height          | 4102                                  Height          | 4102
Timestamp       | 1712345678                            Timestamp       | 1712345678
PrevBlockHash   | 7f3a...                               PrevBlockHash   | 7f3a...
MerkleRoot      | c81e...   (commits to THIS             MerkleRoot      | c81e...
                              block's transactions only)  StateRoot      | 9c14...   <-- NEW
Nonce           | 88213481                                                  (commits to the ENTIRE
Hash            | 4b7d...                                                    UTXO set, after this
                                                                               block's transactions
                                                                               have been applied)
Nonce           | 88213481
Hash            | 4b7d...
```

A `MerkleRoot` proves "these specific transactions were included in this specific block." A `StateRoot` proves something categorically bigger: "given every block up to and including this one, the resulting UTXO set is *exactly* this, and nothing else" — which is precisely the guarantee Volume 9's smart contracts will lean on hard, since a contract's persistent storage (Chapter 66) is exactly the kind of state a root hash needs to commit to, the moment contracts can hold data that outlives any single transaction. This chapter computes the value and shows exactly where it belongs; wiring `StateRoot` into `core.Block` as a field that consensus actually validates on every block is deliberately left for Volume 9, once there is contract state worth committing to beyond the UTXO set alone.

---

## 9. What Real Ethereum Adds That This Chapter Skips

Being explicit about the gap between this chapter's trie and a production one is worth doing honestly, rather than implying they are the same thing at a smaller scale:

- **Extension nodes.** Real Merkle-Patricia tries compress a chain of single-child branch nodes into one "extension" node holding multiple nibbles at once, so a long unique suffix does not cost one full node per hex digit. Our simplified version pays that cost, trading a slightly larger trie for meaningfully simpler code.
- **RLP encoding.** Ethereum encodes each node's contents using RLP (Recursive Length Prefix) before hashing, a specific, canonical byte format chosen for compactness and unambiguous decoding — we use a simpler, direct `sha256.New()`-and-`Write` scheme, which is correct for our purposes but not wire-compatible with any real Ethereum client.
- **Persistent, on-disk tries.** Ethereum's actual state trie lives in its own database (keyed by node hash, so identical subtrees across different states are stored exactly once), not rebuilt fully in memory from a UTXO set on demand. Our version is deliberately in-memory and rebuilt from scratch, which is fine for previewing the concept but would not scale to a state trie with millions of entries.
- **Merkle proofs.** A real state trie lets you prove a single key's value against the root without revealing the rest of the trie — a "state proof" analogous to Chapter 10's transaction Merkle proofs, and the exact mechanism light clients use to verify an account balance without downloading the entire chain. Exercise 8 sketches this.

None of these gaps change the core idea this chapter set out to teach: a trie whose every node's hash commits to everything beneath it turns "does everyone agree on the entire state" into a single 32-byte comparison.

---

## Summary

- A UTXO set stored in a plain key-value bucket has no built-in way to prove its own integrity — corruption or tampering would go completely unnoticed.
- A Patricia trie shares common key prefixes structurally, walking each shared prefix once no matter how many keys begin with it — a real efficiency win for hex-like keys such as hashes and addresses.
- A Merkle-Patricia trie layers Merkle-style hash propagation on top of that prefix structure: every node's hash commits to its own value and to every child's hash, so any single change anywhere propagates all the way to the root.
- This chapter's simplified trie merges what real implementations split into branch, extension, and leaf node types into one 16-way `Node`, trading some efficiency for much simpler code, and is explicit about that trade-off.
- `ComputeStateRoot` walks the entire UTXO set via the `ForEachUTXO` optional capability from Chapter 56 and inserts every entry into a trie, producing one root hash over the entire current state.
- A `StateRoot` field in a block header would prove something a transaction `MerkleRoot` cannot: that the resulting state, after applying every block up to this one, is exactly this and nothing else — exactly what Volume 9's smart contract storage will need.
- Real Merkle-Patricia tries add extension-node compression, RLP encoding, on-disk persistence keyed by node hash, and Merkle-style proofs for individual keys — all named explicitly here as the gap between this chapter's teaching version and a production implementation.

---

## Exercises

### Easy

1. In `Node.hash()`, explain in your own words why writing `byte(i)` before each present child's hash is necessary for correctness, using a concrete example of two different trie shapes that would otherwise hash identically.
2. Modify the demo in Section 6 to insert a third address, print the root hash, then delete Alice's entry entirely (hint: `Put` with a `nil` value is not quite a delete — you will need to think about what "delete" should mean for a trie node that might still have children). Discuss what your chosen behavior implies for `Get` afterward.
3. Compute, by hand or with a short script, how many trie nodes in the worst case a 32-byte (64-nibble) key requires to insert into an empty trie, assuming no other keys share any prefix with it at all.

### Medium

4. Write a test, `TestTrie_OrderIndependence`, that inserts the same five key-value pairs into two separate `Trie` instances in two different orders, and asserts both end up with identical root hashes — directly demonstrating the property Section 6 claims informally.
5. `RootHash()` currently recomputes hashes lazily via `cachedHash`, invalidated on every `Put` along the path from root to the modified leaf. Add a benchmark that inserts 10,000 keys one at a time, calling `RootHash()` after every single insert, and compare its total time against calling `RootHash()` only once at the very end. Explain the difference in terms of how much of the trie each call actually needs to recompute.
6. Extend `ComputeStateRoot` into `ComputeStateRootAt(store Store, bc *core.Blockchain, height int64) ([]byte, error)` that computes what the state root *would have been* immediately after a specific historical block height, by using `UTXOSet.Reindex` against a blockchain "truncated" to that height (you may assume a helper that returns such a truncated `*core.Blockchain` exists). Discuss why this operation is inherently expensive with this chapter's in-memory, rebuild-from-scratch design, and how a real, persisted-per-block trie would avoid that cost.

### Hard

7. Implement extension-node compression: modify `Node` (or add a new node kind) so that a chain of single-child branch nodes collapses into one node holding a multi-nibble path directly, only branching into a real 16-way `Node` where two or more keys actually diverge. Verify, with a test, that your compressed trie produces the *same* root hash as this chapter's uncompressed version for keys that share no prefixes at all, but uses measurably fewer total node allocations for keys that share long prefixes.
8. Design (in comments, or fully implemented if you want the challenge) a `Prove(key []byte) ([][]byte, error)` method that returns the list of node hashes along the path from the trie's root to key's leaf — a Merkle proof for one specific key — and a corresponding `VerifyProof(rootHash, key, value []byte, proof [][]byte) bool` that checks the proof without needing the rest of the trie at all. Explain how this would let a "light client" verify one address's balance against a block's `StateRoot` without downloading the entire UTXO set.
9. Research how `go-ethereum`'s actual `trie` package persists nodes to its underlying database — specifically, how it avoids re-storing an unchanged subtree every time a *different* part of the trie changes. Write a short explanation of the technique (hint: it relates directly to why every node's hash is also its storage key), and explain why this chapter's fully-in-memory, rebuild-on-demand trie sidesteps needing that technique at all, at the cost of scalability.
