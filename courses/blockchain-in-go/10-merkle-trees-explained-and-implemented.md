# Chapter 10: Merkle Trees — Explained and Implemented

Chapter 09 gave you `crypto.Hash`, a way to fingerprint a single canonical byte slice. A block, however, does not hold one piece of data — it holds a whole list of transactions, and Chapter 16 will need a way to fingerprint that entire list with one hash, the same way a single `Note` got fingerprinted last chapter. The naive approach — concatenate every transaction's bytes into one giant blob and call `crypto.Hash` once — works, but it throws away something valuable: the ability to prove that *one* transaction belongs to that list without handing over every other transaction alongside it. A **Merkle tree** solves this by hashing the list in a structured, pairwise way instead of as one flat blob, and this chapter builds one, piece by piece, in Go.

## Table of Contents

1. [Why Not Just Hash the Whole List as One Blob?](#1-why-not-just-hash-the-whole-list-as-one-blob)
2. [The Merkle Tree Idea, Level by Level](#2-the-merkle-tree-idea-level-by-level)
3. [Building a Tree by Hand — 4 Leaves, Step by Step](#3-building-a-tree-by-hand--4-leaves-step-by-step)
4. [The Go Shape — MerkleTree and MerkleNode](#4-the-go-shape--merkletree-and-merklenode)
5. [Implementing NewMerkleTree Step by Step](#5-implementing-newmerkletree-step-by-step)
6. [Handling an Odd Number of Leaves](#6-handling-an-odd-number-of-leaves)
7. [Computing and Reading the Merkle Root](#7-computing-and-reading-the-merkle-root)
8. [Merkle Proofs — Proving Inclusion Without the Whole List](#8-merkle-proofs--proving-inclusion-without-the-whole-list)
9. [Implementing a Merkle Proof Verifier in Go](#9-implementing-a-merkle-proof-verifier-in-go)
10. [Why Light Clients Need This](#10-why-light-clients-need-this)
11. [Where Merkle Trees Fit in GoChain](#11-where-merkle-trees-fit-in-gochain)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Why Not Just Hash the Whole List as One Blob?

Imagine a block containing 2,000 transactions — an ordinary amount for a busy chain. Following Chapter 09's pattern exactly, you could concatenate all 2,000 transactions' serialized bytes end to end and call `crypto.Hash` once on the result. This works, and it is deterministic, and it would detect tampering — change one transaction, and the avalanche effect (Chapter 08, Section 3) makes the combined hash unrecognizable. So far, so good.

Now suppose a mobile wallet — call it Bob's phone — wants to confirm that a specific transaction paying Bob 25 gochips is really included in that block, without downloading and storing all 2,000 other transactions Bob has no personal interest in. With the flat-blob approach, Bob's phone has no choice: to recompute the block's single hash and check it matches, it needs *every* transaction, because the hash function has no internal structure — it is one giant, opaque computation over the whole input. There is no way to check "is X included" without redoing the entire computation over everything.

```
  Flat blob approach:

  tx1 || tx2 || tx3 || ... || tx2000  ──▶  Hash  ──▶  one fingerprint

  To verify tx847 is included, Bob's phone needs all 2000
  transactions' bytes -- there is no partial-computation shortcut.
```

This is the shipping-crate analogy from Chapter 08 pushed past its limit: the barcode machine fingerprints the whole crate at once, but it gives you no way to prove one specific item is inside the crate without opening the whole thing and checking by hand. A Merkle tree is a different way of organizing the *same* hashing work so that proving one item's membership only requires a small, fixed number of extra hashes — not the entire list. Sections 2 and 3 build the structure that makes this possible.

## 2. The Merkle Tree Idea, Level by Level

Instead of hashing every item together in one pass, a Merkle tree hashes items in pairs, then hashes pairs of *those* hashes together, and repeats this pairing-and-hashing process level by level until only one hash remains — the **Merkle root**.

```
  Level 0 (leaves):     hash each item individually

  Level 1:               hash each adjacent PAIR of level-0 hashes

  Level 2:               hash each adjacent PAIR of level-1 hashes

    ...                  ... repeat until one hash remains ...

  Top:                   the Merkle root
```

Think of a single-elimination sports tournament bracket. Sixteen teams (Level 0) play eight matches, producing eight winners (Level 1). Those eight winners play four matches, producing four winners (Level 2), and so on, until one final match crowns a single champion. A Merkle tree runs the exact same shape of computation, except every "match" is a hash of the two inputs concatenated together, and the "champion" at the top is the Merkle root — a single fingerprint that depends on every item at the bottom, but was computed through a structured series of small, pairwise steps rather than one giant pass.

That structure — a tree of pairwise hashes rather than one flat hash — is precisely what makes the Section 1 problem solvable: to check that one leaf near the bottom belongs under a given root, you do not need every other leaf, only the small handful of sibling hashes standing next to your leaf's path as it climbs, match by match, to the top. Section 8 makes this precise; first, Section 3 builds the tree itself, concretely, with real numbers.

## 3. Building a Tree by Hand — 4 Leaves, Step by Step

Take four transactions, referred to here simply as `tx1`, `tx2`, `tx3`, and `tx4` (in a real block these would be full, serialized `Transaction` values from Chapter 32; for this chapter, plain strings are enough to demonstrate the mechanics). Every hash below is a real SHA-256 output, computed exactly the way `crypto.Hash` from Chapter 09 computes it, so you can reproduce every one of these values yourself.

**Stage 1 — hash each leaf individually:**

```
  H1 = Hash("tx1") = 709b55bd3da0f5a8...79a201b
  H2 = Hash("tx2") = 27ca64c092a959c7...c8a779a1e3
  H3 = Hash("tx3") = 1f3cb18e896256d7...d1e2c8b7a9
  H4 = Hash("tx4") = 41b637cfd9eb3e2f...9a086ac78

  Level 0 (leaves):   [H1]      [H2]      [H3]      [H4]
```

**Stage 2 — hash each adjacent pair of leaf hashes together:**

```
  H12 = Hash(H1 || H2) = bbea820f07f7f89a...371cd0c48
  H34 = Hash(H3 || H4) = 5709445d10349996...9104487ca

  Level 0 (leaves):    H1         H2         H3         H4
                        └────┬────┘           └────┬────┘
  Level 1 (pairs):         H12                    H34
```

Here, `||` means byte concatenation: `H12` is the hash of `H1`'s 32 raw bytes followed immediately by `H2`'s 32 raw bytes — 64 bytes of input producing one new 32-byte hash. Note that `H12` depends on *both* `H1` and `H2`; the avalanche effect (Chapter 08, Section 3) means changing a single byte anywhere in `tx1` or `tx2` changes `H1` or `H2`, which changes `H12` completely and unpredictably.

**Stage 3 — hash the two remaining pair-hashes together, producing the root:**

```
  Root = Hash(H12 || H34) = ea59a369466be42d...a4b9af145

  Level 0 (leaves):    H1         H2         H3         H4
                        └────┬────┘           └────┬────┘
  Level 1 (pairs):         H12                    H34
                             └───────────┬───────────┘
  Level 2 (root):                     Root
```

Four leaves collapsed into two pairs, and two pairs collapsed into one root, in exactly two pairing rounds. This single 32-byte `Root` now stands in for the entire ordered list of four transactions, the same way one `Hash()` call stood in for one `Note` in Chapter 09 — except this root was built through a structure that Section 8 will exploit to prove inclusion cheaply. Every value shown above is a real, verifiable SHA-256 output; Sections 4 and 5 turn this exact by-hand process into Go code that produces these same numbers.

## 4. The Go Shape — MerkleTree and MerkleNode

`gochain/crypto` represents this tree with two small types, added alongside `Hash` from Chapter 09:

```go
// crypto/merkle.go
package crypto

// MerkleNode is a single node in a Merkle tree. A leaf node has Left
// and Right both nil, and Data is the hash of its raw input data. An
// internal node has Left and Right set, and Data is the hash of its
// two children's Data, concatenated in a fixed left-then-right order
// -- exactly the H12 = Hash(H1 || H2) computation from Section 3.
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree is a binary tree of hashes built over a list of data
// items, collapsing down to a single root hash (RootNode.Data) that
// fingerprints the entire list.
type MerkleTree struct {
	RootNode *MerkleNode
}
```

Notice what `MerkleNode.Data` means changes depending on where the node sits in the tree: at the leaves, it is the hash of a real transaction's bytes; everywhere above the leaves, it is the hash of two children's hashes. This dual meaning is deliberate and is what makes the pairing-and-hashing process in Section 2 uniform all the way up the tree — the same node type, and the same hashing rule, works at every level.

## 5. Implementing NewMerkleTree Step by Step

Building the tree needs one small constructor for a single node, and then a loop that repeatedly pairs up a level's nodes to build the level above it — mechanically the same process Section 3 walked through by hand.

```go
// crypto/merkle.go

// NewMerkleNode builds a single node. If left and right are both nil,
// this is a leaf node, and data is hashed directly to become its
// Data. Otherwise, this is an internal node: data is ignored, and its
// Data is the hash of its two children's Data concatenated together,
// left then right.
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{Left: left, Right: right}

	if left == nil && right == nil {
		node.Data = Hash(data)
		return node
	}

	combined := append(append([]byte{}, left.Data...), right.Data...)
	node.Data = Hash(combined)
	return node
}

// NewMerkleTree builds a Merkle tree over data, one leaf per element,
// in the order given, and returns the tree rooted at the single
// combined hash.
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{RootNode: NewMerkleNode(nil, nil, []byte{})}
	}

	var nodes []*MerkleNode
	for _, d := range data {
		nodes = append(nodes, NewMerkleNode(nil, nil, d))
	}

	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1]) // Section 6
		}

		var level []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			level = append(level, NewMerkleNode(nodes[i], nodes[i+1], nil))
		}
		nodes = level
	}

	return &MerkleTree{RootNode: nodes[0]}
}
```

Walk through what happens for the four-leaf example from Section 3. The first loop builds four leaf nodes: `nodes = [H1, H2, H3, H4]` (writing each node by the hash it holds, for readability). The outer `for len(nodes) > 1` loop then runs twice. On its first pass, `len(nodes)` is 4 (even, so Section 6's duplication does not trigger), and the inner loop pairs `(H1,H2)` into `H12` and `(H3,H4)` into `H34`, leaving `nodes = [H12, H34]`. On the second pass, `len(nodes)` is 2, and the inner loop pairs `(H12,H34)` into `Root`, leaving `nodes = [Root]`. The outer loop's condition `len(nodes) > 1` is now false, so it exits, and `NewMerkleTree` returns a `MerkleTree` whose `RootNode` is exactly the `Root` computed by hand in Section 3.

## 6. Handling an Odd Number of Leaves

Section 5's pairing loop assumes an even number of nodes at every level — but a block will not always contain a number of transactions that happens to be a power of two, or even just even. `NewMerkleTree`'s line `if len(nodes)%2 != 0 { nodes = append(nodes, nodes[len(nodes)-1]) }` handles this by **duplicating the last node** whenever a level has an odd count, pairing it with a copy of itself — the same convention Bitcoin uses.

```
  3 leaves:              H1        H2        H3
                          └───┬────┘    (H3 has no partner)

  Duplicate H3:           H1        H2        H3       H3(copy)
                          └───┬────┘          └───┬────┘

  Pair up:                  H12                  H33
                              └─────────┬─────────┘
                                      Root
```

Here `H33 = Hash(H3 || H3)` — the same node hashed with itself. This keeps the pairing rule from Section 5 uniform (every level always has an even count by the time it is paired), at the cost of one detail worth remembering for Section 9: a leaf that ends up duplicated this way has a sibling that is a copy of itself, not a different transaction — the proof-generation code needs to produce that same duplicate, or a proof for the last item in an odd-sized list will not verify.

## 7. Computing and Reading the Merkle Root

A small accessor method, following Chapter 09's `HashHex` pattern for readability, rounds out the basic tree:

```go
// crypto/merkle.go

// MerkleRoot returns the tree's root hash -- the single fingerprint
// that stands in for the entire list of leaves it was built from.
func (t *MerkleTree) MerkleRoot() []byte {
	return t.RootNode.Data
}
```

```go
data := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
tree := NewMerkleTree(data)
fmt.Println(hex.EncodeToString(tree.MerkleRoot()))
// ea59a369466be42d1a4783f09ae0721a5a157d6dba9c4b053d407b5a4b9af145
```

That printed value is the exact `Root` from Section 3's by-hand walkthrough — running this code yourself is a good way to confirm the diagrams in this chapter are not just illustrations but genuine, reproducible SHA-256 computations. Chapter 17 stores exactly this value in `core.Block.MerkleRoot`, one hash standing in for however many thousands of transactions a block actually contains.

## 8. Merkle Proofs — Proving Inclusion Without the Whole List

Return to Bob's phone from Section 1, now armed with the tree structure from Sections 2 through 6. Bob wants to prove `tx2` is included under the `Root` computed in Section 3, without needing `tx1`, `tx3`, or `tx4`'s full data. Look again at how `Root` was built: `Root = Hash(H12 || H34)`, and `H12 = Hash(H1 || H2)`. To recompute `Root` starting from `H2` alone, Bob needs exactly two extra pieces of information: `H1` (to combine with `H2` and reproduce `H12`), and `H34` (to combine with the recomputed `H12` and reproduce `Root`).

```
                              Root
                        ┌──────┴──────┐
                       H12            H34   ◀── need this to finish
                  ┌─────┴─────┐
                 H1           H2   ◀── H1 needed to start;
              (sibling,                H2 is our own leaf
               provide this)

  Path from H2 to the root, with the two sibling hashes marked:
  a Merkle proof for tx2 is exactly {H1, H34}, in that order,
  plus a note of which side each sibling sits on (H1 is to H2's
  LEFT; H34 is to H12's RIGHT).
```

This pair of hashes — `H1` and `H34` — is a **Merkle proof** for `tx2`. Whoever wants to verify Bob's claim needs only: `tx2`'s own raw data, this list of two sibling hashes (each tagged left or right), and the `Root` they are checking against. Notice how little that is compared to the full list — for four leaves the savings look modest, but a Merkle tree's height grows only as the *logarithm* of the number of leaves, so a block with a million transactions needs a proof only about twenty hashes long, not a million. This is the concrete payoff Section 1 promised: proving membership scales with the *height* of the tree, not its *width*.

## 9. Implementing a Merkle Proof Verifier in Go

A proof is naturally represented as an ordered list of steps, each carrying a sibling hash and which side it sits on:

```go
// crypto/merkle.go

// ProofStep is one hop on the path from a leaf up to the Merkle root:
// the sibling hash needed at that level, and whether the sibling sits
// to the left or right of the hash being carried upward.
type ProofStep struct {
	Hash   []byte
	IsLeft bool
}
```

Generating a proof rebuilds the tree level by level, exactly as `NewMerkleTree` does, but instead of discarding each level once the next is built, it records the one sibling hash the target leaf needs at that level, and tracks the target's index as it folds in half on every pass:

```go
// crypto/merkle.go
import "fmt"

// GenerateMerkleProof builds a Merkle proof for the leaf at index,
// given the same ordered data NewMerkleTree would build a tree from.
func GenerateMerkleProof(data [][]byte, index int) ([]ProofStep, error) {
	if index < 0 || index >= len(data) {
		return nil, fmt.Errorf("index %d out of range for %d leaves", index, len(data))
	}

	var nodes []*MerkleNode
	for _, d := range data {
		nodes = append(nodes, NewMerkleNode(nil, nil, d))
	}

	var proof []ProofStep

	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		if index%2 == 0 {
			// our node is the LEFT half of its pair; the sibling
			// to its right is what we need to record.
			proof = append(proof, ProofStep{Hash: nodes[index+1].Data, IsLeft: false})
		} else {
			// our node is the RIGHT half of its pair; the sibling
			// to its left is what we need to record.
			proof = append(proof, ProofStep{Hash: nodes[index-1].Data, IsLeft: true})
		}

		var level []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			level = append(level, NewMerkleNode(nodes[i], nodes[i+1], nil))
		}
		nodes = level
		index = index / 2 // our node's position in the level above
	}

	return proof, nil
}
```

Verifying a proof runs the same combining step Section 3 used to build the tree in the first place, just following the recorded path instead of the whole structure:

```go
// crypto/merkle.go
import "bytes"

// VerifyMerkleProof checks that leafData, combined with the sibling
// hashes in proof, in order, produces exactly root. It needs only the
// original leaf's data and its short list of sibling hashes -- never
// any of the tree's other leaves.
func VerifyMerkleProof(leafData []byte, proof []ProofStep, root []byte) bool {
	current := Hash(leafData)

	for _, step := range proof {
		if step.IsLeft {
			current = Hash(append(append([]byte{}, step.Hash...), current...))
		} else {
			current = Hash(append(append([]byte{}, current...), step.Hash...))
		}
	}

	return bytes.Equal(current, root)
}
```

Run it against the Section 3 example, and the numbers match exactly:

```go
data := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
tree := NewMerkleTree(data)

proof, _ := GenerateMerkleProof(data, 1) // tx2 is index 1
// proof[0]: IsLeft=true,  Hash=H1  (709b55bd...79a201b)
// proof[1]: IsLeft=false, Hash=H34 (5709445d...9104487ca)

ok := VerifyMerkleProof([]byte("tx2"), proof, tree.MerkleRoot())
fmt.Println(ok) // true

tampered := VerifyMerkleProof([]byte("tx2-tampered"), proof, tree.MerkleRoot())
fmt.Println(tampered) // false
```

The first `ok` matches the hand-traced diagram in Section 8 exactly: `Hash(H1 || Hash("tx2"))` reproduces `H12`, then `Hash(H12 || H34)` reproduces `Root`. The second call demonstrates the property that makes a proof trustworthy rather than just convenient: feed the verifier data that was not actually in the tree, and the avalanche effect (Chapter 08, Section 3) means the recomputed root comes out completely different from the real one, so `VerifyMerkleProof` correctly returns `false`. A short table-driven test belongs in `crypto/merkle_test.go` covering exactly these two cases, plus the odd-leaf-count case from Section 6 — generate a proof for the *duplicated* last leaf in a 3-item list, and confirm it still verifies correctly against that tree's root.

## 10. Why Light Clients Need This

A **light client** — a mobile wallet, a browser extension, any program that wants to interact with a blockchain without storing the entire thing — cannot afford to download every transaction in every block just to answer "did I receive this payment?" Bitcoin's original whitepaper calls this exact technique **Simplified Payment Verification (SPV)**: a light client downloads only block headers (which include each block's Merkle root, tiny and cheap to store by the thousands) and, when it needs to confirm one specific transaction, asks a fuller node for that transaction's Merkle proof against the relevant header's root.

```
   Full node (has every transaction in every block)
        │
        │  "here is tx2's data, plus its 2-hash Merkle proof"
        ▼
   Light client (Bob's phone -- has only block headers)
        │
        │  recomputes the root from tx2 + the proof (Section 9)
        ▼
   Compare recomputed root to the root already stored in the
   block header Bob's phone downloaded.  Match?  tx2 is
   genuinely, provably included -- without Bob's phone ever
   needing the other 1,999 transactions in that block.
```

This is the same trust model this entire course has been building since Chapter 08: nobody has to *ask* the full node to be honest about `tx2`'s inclusion and simply believe it. The full node's claim is checked mathematically, cheaply, by anyone, the same way every hash in this book has been checkable rather than merely trusted. A light client that verifies proofs this way gets real security guarantees with a tiny fraction of the storage and bandwidth a full node needs — which is exactly why every mobile crypto wallet in existence leans on some version of this technique, rather than each one running a full copy of the entire chain on your phone.

## 11. Where Merkle Trees Fit in GoChain

`core.Block`, defined in Chapter 17, stores a `MerkleRoot` field computed exactly the way Section 7 computes one — over that block's list of transactions — rather than storing or hashing the full transaction list directly inside the block header. This keeps block headers small and fast to verify, and sets up two payoffs that only become fully visible later: Chapter 49's blockchain synchronization can let a syncing node fetch and verify transactions incrementally against a known-good root, and Volume 10's block explorer can offer "prove this transaction is in this block" as a real, checkable feature rather than a claim you would have to trust the explorer's own database for.

```
   Block N
  ┌──────────────────────────┐
  │ Height, Timestamp,        │
  │ PrevBlockHash             │
  │ MerkleRoot  ◀── this chapter's Root, over Block N's txs
  │ Hash, Nonce                │
  ├──────────────────────────┤
  │ Transactions: [tx1..txN]   │  ◀── the actual data the
  └──────────────────────────┘      MerkleRoot fingerprints
```

Chapter 09 gave every GoChain type its own `Serialize()` for turning one value into canonical bytes; this chapter gives GoChain a way to fingerprint an entire *list* of values with one hash while keeping individual membership provable. Together, they are the two hashing tools every later chapter builds with — Chapter 11 turns to the other half of Volume 2's cryptography, public and private keys, before Volume 3 puts both tools to work inside a real, running `core.Block`.

---

## Summary

- Hashing an entire transaction list as one flat blob works for tamper-detection but makes it impossible to prove one transaction's inclusion without handing over every other transaction too.
- A **Merkle tree** hashes leaves in pairs, then hashes pairs of those hashes, repeating level by level until a single **Merkle root** remains — a structured alternative to one flat hash.
- GoChain's shape is `MerkleTree{RootNode *MerkleNode}` and `MerkleNode{Left, Right *MerkleNode; Data []byte}`: leaves hash raw data directly, internal nodes hash their two children's `Data` concatenated together.
- `NewMerkleTree` builds the tree bottom-up, pairing nodes level by level; an odd node count at any level is handled by duplicating the last node, the same convention Bitcoin uses.
- A **Merkle proof** for a leaf is the short list of sibling hashes (each tagged left or right) along the path from that leaf up to the root — its length grows with the tree's height (logarithmic in the number of leaves), not its width.
- `VerifyMerkleProof` recombines a claimed leaf's hash with its proof's sibling hashes, in order, and checks the result against the known root; tampering with the leaf data or any sibling hash makes the recomputed root come out unrecognizable, thanks to the avalanche effect.
- **Light clients** (mobile wallets, SPV clients) store only small block headers and Merkle roots, and use Merkle proofs to verify specific transactions are genuinely included without downloading and storing the entire blockchain.
- `core.Block.MerkleRoot`, built starting in Chapter 17, is exactly this chapter's Merkle root, computed over a block's transaction list, keeping block headers small while still making every transaction's inclusion independently provable.

---

## Exercises

### Easy

1. By hand or with a script, compute `H1 = Hash("tx1")` through `H4 = Hash("tx4")` yourself using any SHA-256 tool, and confirm your results match the values given in Section 3. Then compute `H12 = Hash(H1 || H2)` (concatenating the raw bytes, not the hex strings) and confirm it matches the chapter's stated value.

2. Using `NewMerkleTree` and `MerkleRoot`, build a tree over five string leaves of your choosing and print the root as hex. Then change a single character in just one leaf, rebuild the tree, and print the new root. Confirm the two roots share no visible resemblance, and explain which property from Chapter 08 guarantees this.

3. Explain, in your own words and without code, why `MerkleNode.Data` means something different for a leaf node versus an internal node, and why that difference does not cause any problem for the pairing-and-hashing process described in Section 2.

### Medium

4. Using `GenerateMerkleProof` and `VerifyMerkleProof` on the four-leaf example from this chapter, generate and verify proofs for all four leaves (index 0 through 3), and print each proof's two `ProofStep` values. Confirm by hand that the `IsLeft` flags for `tx1`'s proof are the mirror image of `tx4`'s proof, and explain why that symmetry makes sense given the tree's shape.

5. Build a Merkle tree over 7 leaves (an odd number that is not one less than a power of two) and trace, level by level as in Section 6, exactly which nodes get duplicated at each level. Draw the resulting tree as an ASCII diagram in the style of Section 3.

6. Write a test that generates a valid Merkle proof for some leaf, then deliberately corrupts one byte inside one `ProofStep.Hash` (not the leaf data itself) before calling `VerifyMerkleProof`. Confirm verification fails, and explain in 3-4 sentences why corrupting a *sibling* hash breaks verification just as thoroughly as corrupting the leaf itself.

### Hard

7. The chapter states that a Merkle proof's length grows with the *logarithm* of the number of leaves. For a block containing 1,048,576 (2²⁰) transactions, calculate exactly how many `ProofStep` entries a Merkle proof would contain, and compare that to how many transactions' full data a naive "prove inclusion" scheme (handing over every transaction) would require. Express the savings as a ratio.

8. Research how Bitcoin's real Merkle tree implementation differs from this chapter's in one specific, documented way (for example: the exact byte order used when concatenating hashes, or a historical vulnerability related to duplicate transactions in an unbalanced tree, sometimes called a "CVE-2012-2459"-style issue). Write a 250-350 word explanation of the difference and why it matters in practice.

9. Design (in a short written proposal, 300-500 words, with at least one diagram) an extension to `GenerateMerkleProof` and `VerifyMerkleProof` that could prove a *range* of consecutive leaves are all included under a root using fewer total hashes than generating and verifying each leaf's proof independently. You do not need to implement it — explain which sibling hashes would be shared across the range and why reusing them saves work.
