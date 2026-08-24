package main

import "fmt"

type TxOutput struct {
	Value   int64
	Address string
}

type TxInput struct {
	TxID     []byte
	OutIndex int
	Address  string
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type Mempool struct {
	pending map[string]*Transaction // hex(txID) -> transaction
	claimed map[string]bool         // "hex(txID):outIndex" -> claimed
	order   []string                // preserves insertion order
}

func NewMempool() *Mempool {
	return &Mempool{
		pending: make(map[string]*Transaction),
		claimed: make(map[string]bool),
	}
}

// claimKey builds the map key identifying one specific (txID, outIndex)
// pair. (Provided for you.)
func claimKey(txID []byte, outIndex int) string {
	return fmt.Sprintf("%x:%d", txID, outIndex)
}

// Add validates tx against everything currently pending and, if nothing
// conflicts, adds it. Returns an error if tx tries to spend a
// (TxID, OutIndex) pair some other pending transaction already claims.
func (mp *Mempool) Add(tx *Transaction) error {
	panic("TODO: implement Add")
}

// Remove takes a transaction out of the mempool (called after it's
// been mined into a block).
func (mp *Mempool) Remove(txID []byte) {
	panic("TODO: implement Remove")
}

// Pending returns every transaction currently waiting, in the order
// they were added.
func (mp *Mempool) Pending() []*Transaction {
	panic("TODO: implement Pending")
}

// Has reports whether a transaction with the given ID is currently
// pending.
func (mp *Mempool) Has(txID []byte) bool {
	panic("TODO: implement Has")
}
