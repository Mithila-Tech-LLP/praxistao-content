# Task 07: Merkle Trees

## What you will build

A Merkle tree over a list of transactions, plus a function that produces and verifies a Merkle proof — a way to prove one specific transaction is included in a set without needing every other transaction present.

## Concepts

### Hashing pairs, repeatedly, up to one root

Instead of hashing a whole list of transactions as one blob, a Merkle tree hashes pairs of leaves together, then hashes pairs of *those* hashes, and so on, until only one hash remains: the Merkle root.

```
   leaf: H(tx1)   H(tx2)   H(tx3)   H(tx4)
              \    /            \    /
           H(H(tx1)+H(tx2))  H(H(tx3)+H(tx4))
                      \            /
                   H( ... + ... )  <- Merkle root
```

### Why bother — the proof

Once you have a tree, you can prove "transaction 2 is definitely in this set" by handing over just the *sibling* hashes along the path from that leaf up to the root (here, `H(tx1)` and `H(H(tx3)+H(tx4))`), rather than the other three full transactions. Anyone can recombine those siblings with the leaf they already have and check the result matches the known root.

## Interface to implement

```go
type MerkleNode struct {
	Left, Right *MerkleNode
	Hash        []byte
}

// BuildMerkleTree builds a tree from data (each entry treated as one
// leaf) and returns its root node. If len(data) is odd, duplicate the
// last leaf so every level has an even number of nodes.
func BuildMerkleTree(data [][]byte) *MerkleNode

// MerkleRoot returns the tree's root hash.
func MerkleRoot(data [][]byte) []byte

// MerkleProof is the ordered list of sibling hashes needed to
// reconstruct the root from one specific leaf.
type MerkleProof struct {
	LeafIndex int
	Siblings  [][]byte
}

// BuildMerkleProof returns the proof for the leaf at index i in data.
func BuildMerkleProof(data [][]byte, i int) MerkleProof

// VerifyMerkleProof reports whether leaf, combined with proof's
// sibling hashes in order, reconstructs root.
func VerifyMerkleProof(leaf []byte, proof MerkleProof, root []byte) bool
```

## Hints

- Hash each leaf once up front (`Hash(data[i])`), then build the tree over those leaf hashes, not the raw data.
- Combining two nodes: `Hash(append(left.Hash, right.Hash...))`.
- For `BuildMerkleProof`, the simplest correct approach is to rebuild the tree level by level, and at each level record whichever neighbor (left or right) of your current node was NOT the one you came from.
- `VerifyMerkleProof` needs to know, at each step, whether the sibling goes on the left or the right when recombining — `LeafIndex`'s bits (even/odd at each level, as you divide the index by 2 going up) tell you this: if the current position is even, the sibling is on the right; if odd, the sibling is on the left.
- Test with 4 and 5 leaves (even and odd counts) to make sure your odd-leaf duplication logic is correct, and verify a proof for every leaf index against the real root — not just leaf 0.

## Run the tests

```bash
cd starter/task-07-merkle-trees
go test ./...
```

All tests must pass before moving to Task 08.
