package main

import (
	"bytes"
	"testing"
)

func txData(n int) [][]byte {
	var data [][]byte
	for i := 0; i < n; i++ {
		data = append(data, []byte{byte('a' + i)})
	}
	return data
}

func TestMerkleRootDeterministic(t *testing.T) {
	data := txData(4)
	r1 := MerkleRoot(data)
	r2 := MerkleRoot(data)
	if !bytes.Equal(r1, r2) {
		t.Fatal("expected the same data to produce the same Merkle root every time")
	}
}

func TestMerkleRootChangesWithData(t *testing.T) {
	data1 := txData(4)
	data2 := txData(4)
	data2[2] = []byte("different")

	if bytes.Equal(MerkleRoot(data1), MerkleRoot(data2)) {
		t.Fatal("expected changing one leaf to change the Merkle root")
	}
}

func TestMerkleRootHandlesOddCount(t *testing.T) {
	data := txData(5)
	root := MerkleRoot(data)
	if len(root) == 0 {
		t.Fatal("expected a valid, non-empty root for an odd number of leaves")
	}
}

func TestProofVerifiesForEveryLeaf_EvenCount(t *testing.T) {
	data := txData(4)
	root := MerkleRoot(data)

	for i := range data {
		proof := BuildMerkleProof(data, i)
		if !VerifyMerkleProof(data[i], proof, root) {
			t.Fatalf("expected proof for leaf %d to verify against the real root", i)
		}
	}
}

func TestProofVerifiesForEveryLeaf_OddCount(t *testing.T) {
	data := txData(5)
	root := MerkleRoot(data)

	for i := range data {
		proof := BuildMerkleProof(data, i)
		if !VerifyMerkleProof(data[i], proof, root) {
			t.Fatalf("expected proof for leaf %d (of 5) to verify against the real root", i)
		}
	}
}

func TestProofFailsForWrongLeaf(t *testing.T) {
	data := txData(4)
	root := MerkleRoot(data)

	proof := BuildMerkleProof(data, 1)
	if VerifyMerkleProof([]byte("not the real leaf"), proof, root) {
		t.Fatal("expected verification to fail when the leaf doesn't match what the proof was built for")
	}
}

func TestProofFailsAgainstWrongRoot(t *testing.T) {
	data := txData(4)
	proof := BuildMerkleProof(data, 1)

	wrongRoot := Hash([]byte("not the real root"))
	if VerifyMerkleProof(data[1], proof, wrongRoot) {
		t.Fatal("expected verification to fail against an unrelated root")
	}
}
