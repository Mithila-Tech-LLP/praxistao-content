package main

import "testing"

func TestAddAndHas(t *testing.T) {
	mp := NewMempool()
	tx := &Transaction{ID: []byte("tx1"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}}

	if err := mp.Add(tx); err != nil {
		t.Fatalf("unexpected error adding a fresh transaction: %v", err)
	}
	if !mp.Has(tx.ID) {
		t.Fatal("expected Has to report true right after Add")
	}
}

func TestConflictingSpendRejected(t *testing.T) {
	mp := NewMempool()
	tx1 := &Transaction{ID: []byte("tx1"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}}
	tx2 := &Transaction{ID: []byte("tx2"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}} // same UTXO!

	if err := mp.Add(tx1); err != nil {
		t.Fatalf("unexpected error adding tx1: %v", err)
	}
	err := mp.Add(tx2)
	if err == nil {
		t.Fatal("expected Add to reject tx2, which spends the same output as tx1")
	}
	if !mp.Has(tx1.ID) {
		t.Fatal("expected tx1 to remain in the mempool after tx2 was rejected")
	}
	if mp.Has(tx2.ID) {
		t.Fatal("expected tx2 to NOT be in the mempool after being rejected")
	}
}

func TestNonConflictingSpendsBothAccepted(t *testing.T) {
	mp := NewMempool()
	tx1 := &Transaction{ID: []byte("tx1"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}}
	tx2 := &Transaction{ID: []byte("tx2"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 1}}} // different output

	if err := mp.Add(tx1); err != nil {
		t.Fatalf("unexpected error adding tx1: %v", err)
	}
	if err := mp.Add(tx2); err != nil {
		t.Fatalf("unexpected error adding tx2 (spends a different output): %v", err)
	}
	if len(mp.Pending()) != 2 {
		t.Fatalf("expected 2 pending transactions, got %d", len(mp.Pending()))
	}
}

func TestRemoveFreesClaimedOutput(t *testing.T) {
	mp := NewMempool()
	tx1 := &Transaction{ID: []byte("tx1"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}}
	tx2 := &Transaction{ID: []byte("tx2"), Inputs: []TxInput{{TxID: []byte("prev"), OutIndex: 0}}}

	if err := mp.Add(tx1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mp.Remove(tx1.ID)

	if mp.Has(tx1.ID) {
		t.Fatal("expected tx1 to be gone after Remove")
	}
	if err := mp.Add(tx2); err != nil {
		t.Fatalf("expected tx2 to be accepted after tx1 (which claimed the same output) was removed: %v", err)
	}
}

func TestPendingPreservesInsertionOrder(t *testing.T) {
	mp := NewMempool()
	tx1 := &Transaction{ID: []byte("tx1"), Inputs: []TxInput{{TxID: []byte("a"), OutIndex: 0}}}
	tx2 := &Transaction{ID: []byte("tx2"), Inputs: []TxInput{{TxID: []byte("b"), OutIndex: 0}}}
	tx3 := &Transaction{ID: []byte("tx3"), Inputs: []TxInput{{TxID: []byte("c"), OutIndex: 0}}}

	mp.Add(tx1)
	mp.Add(tx2)
	mp.Add(tx3)

	pending := mp.Pending()
	if len(pending) != 3 {
		t.Fatalf("expected 3 pending transactions, got %d", len(pending))
	}
	if string(pending[0].ID) != "tx1" || string(pending[1].ID) != "tx2" || string(pending[2].ID) != "tx3" {
		t.Fatal("expected Pending() to preserve insertion order")
	}
}
