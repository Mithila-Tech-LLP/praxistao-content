package main

import "testing"

func TestMineRewardsMiner(t *testing.T) {
	chain := NewChain()
	miner := NewWallet()

	chain.Mine(miner.Address(), 50, 12)

	if got := chain.Balance(miner.Address()); got != 50 {
		t.Fatalf("expected miner's balance to be 50 after mining, got %d", got)
	}
}

func TestSendRejectsInsufficientFunds(t *testing.T) {
	chain := NewChain()
	alice := NewWallet()
	bob := NewWallet()

	chain.Mine(alice.Address(), 50, 12)

	_, err := chain.Send(alice.PrivateKey, alice.Address(), bob.Address(), 1000)
	if err == nil {
		t.Fatal("expected Send to reject an amount exceeding alice's balance")
	}
}

func TestSendRejectsMismatchedKey(t *testing.T) {
	chain := NewChain()
	alice := NewWallet()
	mallory := NewWallet()
	bob := NewWallet()

	chain.Mine(alice.Address(), 50, 12)

	// Mallory's private key does not match alice's address.
	_, err := chain.Send(mallory.PrivateKey, alice.Address(), bob.Address(), 10)
	if err == nil {
		t.Fatal("expected Send to reject a private key that doesn't match the from address")
	}
}

func TestFullSendAndMineLifecycle(t *testing.T) {
	chain := NewChain()
	alice := NewWallet()
	bob := NewWallet()

	// Mine a reward for alice.
	chain.Mine(alice.Address(), 50, 12)
	if got := chain.Balance(alice.Address()); got != 50 {
		t.Fatalf("expected alice's balance to be 50 after the first mine, got %d", got)
	}

	// Send 20 from alice to bob -- still pending, shouldn't move balances yet.
	tx, err := chain.Send(alice.PrivateKey, alice.Address(), bob.Address(), 20)
	if err != nil {
		t.Fatalf("unexpected error sending: %v", err)
	}
	if tx == nil {
		t.Fatal("expected Send to return the built transaction")
	}
	if got := chain.Balance(bob.Address()); got != 0 {
		t.Fatalf("expected bob's balance to still be 0 before mining, got %d", got)
	}

	// Mine again -- this should confirm the pending transaction.
	chain.Mine(alice.Address(), 50, 12)

	if got := chain.Balance(alice.Address()); got != 80 {
		// 50 (original) - 20 (sent) + 50 (second mining reward) = 80
		t.Fatalf("expected alice's balance to be 80 after sending 20 and mining a second reward, got %d", got)
	}
	if got := chain.Balance(bob.Address()); got != 20 {
		t.Fatalf("expected bob's balance to be 20 after the transaction was mined, got %d", got)
	}
}

func TestMineClearsMempool(t *testing.T) {
	chain := NewChain()
	alice := NewWallet()
	bob := NewWallet()

	chain.Mine(alice.Address(), 50, 12)
	chain.Send(alice.PrivateKey, alice.Address(), bob.Address(), 10)

	if len(chain.mempool.Pending()) != 1 {
		t.Fatalf("expected 1 pending transaction before mining, got %d", len(chain.mempool.Pending()))
	}

	chain.Mine(alice.Address(), 50, 12)

	if len(chain.mempool.Pending()) != 0 {
		t.Fatalf("expected the mempool to be empty after mining, got %d pending", len(chain.mempool.Pending()))
	}
}

func TestDoubleSpendRejectedByMempool(t *testing.T) {
	chain := NewChain()
	alice := NewWallet()
	bob := NewWallet()
	carol := NewWallet()

	chain.Mine(alice.Address(), 50, 12)

	// First spend of alice's only UTXO succeeds.
	if _, err := chain.Send(alice.PrivateKey, alice.Address(), bob.Address(), 50); err != nil {
		t.Fatalf("unexpected error on first send: %v", err)
	}

	// A second transaction trying to spend the SAME (now-claimed) UTXO
	// must be rejected by the mempool before it's ever mined.
	_, err := chain.Send(alice.PrivateKey, alice.Address(), carol.Address(), 50)
	if err == nil {
		t.Fatal("expected the mempool to reject a conflicting second spend of the same UTXO")
	}
}
