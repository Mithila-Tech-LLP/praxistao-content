package main

import (
	"crypto/sha256"
)

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

type MerkleNode struct {
	Left, Right *MerkleNode
	Hash        []byte
}

// combine hashes two child hashes together into their parent's hash.
// (Provided for you.)
func combine(a, b []byte) []byte {
	buf := make([]byte, 0, len(a)+len(b))
	buf = append(buf, a...)
	buf = append(buf, b...)
	return Hash(buf)
}

// BuildMerkleTree builds a tree from data (each entry treated as one
// leaf) and returns its root node. If len(data) is odd, duplicate the
// last leaf so every level has an even number of nodes.
func BuildMerkleTree(data [][]byte) *MerkleNode {
	panic("TODO: implement BuildMerkleTree")
}

// MerkleRoot returns the tree's root hash.
func MerkleRoot(data [][]byte) []byte {
	panic("TODO: implement MerkleRoot using BuildMerkleTree")
}

// MerkleProof is the ordered list of sibling hashes needed to
// reconstruct the root from one specific leaf.
type MerkleProof struct {
	LeafIndex int
	Siblings  [][]byte
}

// BuildMerkleProof returns the proof for the leaf at index i in data.
func BuildMerkleProof(data [][]byte, i int) MerkleProof {
	panic("TODO: implement BuildMerkleProof")
}

// VerifyMerkleProof reports whether leaf, combined with proof's sibling
// hashes in order, reconstructs root.
func VerifyMerkleProof(leaf []byte, proof MerkleProof, root []byte) bool {
	panic("TODO: implement VerifyMerkleProof")
}
