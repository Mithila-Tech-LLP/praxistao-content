package main

import "testing"

func TestCoinbaseCreatesUTXO(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	balance := BalanceOf(txs, "alice")
	if balance != 50 {
		t.Fatalf("expected alice's balance to be 50, got %d", balance)
	}
}

func TestBalanceZeroForUnknownAddress(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	if BalanceOf(txs, "bob") != 0 {
		t.Fatal("expected bob's balance to be 0 before receiving anything")
	}
}

func TestNewTransactionInsufficientFunds(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	_, err := NewTransaction(txs, "alice", "bob", 1000)
	if err == nil {
		t.Fatal("expected an error building a transaction for more than the sender's balance")
	}
}

func TestNewTransactionExactAmountNoChange(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	tx, err := NewTransaction(txs, "alice", "bob", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tx.Outputs) != 1 {
		t.Fatalf("expected exactly 1 output (no change) when spending the exact balance, got %d", len(tx.Outputs))
	}
	if tx.Outputs[0].Value != 50 || tx.Outputs[0].Address != "bob" {
		t.Fatalf("expected a single 50-value output to bob, got %+v", tx.Outputs[0])
	}
}

func TestNewTransactionWithChange(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	tx, err := NewTransaction(txs, "alice", "bob", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tx.Outputs) != 2 {
		t.Fatalf("expected 2 outputs (payment + change), got %d", len(tx.Outputs))
	}

	var paidBob, changeToAlice int64
	for _, out := range tx.Outputs {
		if out.Address == "bob" {
			paidBob = out.Value
		}
		if out.Address == "alice" {
			changeToAlice = out.Value
		}
	}
	if paidBob != 20 {
		t.Fatalf("expected bob to receive 20, got %d", paidBob)
	}
	if changeToAlice != 30 {
		t.Fatalf("expected alice's change to be 30, got %d", changeToAlice)
	}
}

func TestBalancesUpdateAfterTransaction(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	tx, err := NewTransaction(txs, "alice", "bob", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	txs = append(txs, tx)

	if got := BalanceOf(txs, "alice"); got != 30 {
		t.Fatalf("expected alice's balance to be 30 after sending 20 of 50, got %d", got)
	}
	if got := BalanceOf(txs, "bob"); got != 20 {
		t.Fatalf("expected bob's balance to be 20 after receiving it, got %d", got)
	}
}

func TestSpentOutputNotDoubleCounted(t *testing.T) {
	cb := NewCoinbaseTX("alice", 50)
	txs := []*Transaction{cb}

	tx1, err := NewTransaction(txs, "alice", "bob", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	txs = append(txs, tx1)

	// Alice's coinbase output is now fully spent -- she has nothing left.
	if got := BalanceOf(txs, "alice"); got != 0 {
		t.Fatalf("expected alice's balance to be 0 after spending her entire coinbase reward, got %d", got)
	}

	_, err = NewTransaction(txs, "alice", "bob", 1)
	if err == nil {
		t.Fatal("expected building a new transaction from alice (now with 0 balance) to fail")
	}
}
