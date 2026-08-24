package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
)

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

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

// serialize builds a deterministic byte representation of tx, for
// hashing into its ID. (Provided for you.)
func (tx *Transaction) serialize() []byte {
	var buf bytes.Buffer
	for _, in := range tx.Inputs {
		buf.Write(in.TxID)
		idx := make([]byte, 8)
		binary.BigEndian.PutUint64(idx, uint64(in.OutIndex))
		buf.Write(idx)
		buf.WriteString(in.Address)
	}
	for _, out := range tx.Outputs {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, uint64(out.Value))
		buf.Write(val)
		buf.WriteString(out.Address)
	}
	return buf.Bytes()
}

// isCoinbase reports whether tx is a coinbase (reward) transaction: one
// input, with an empty TxID.
func isCoinbase(tx *Transaction) bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0
}

// NewCoinbaseTX creates a reward transaction with no real inputs,
// paying reward coins to address.
func NewCoinbaseTX(address string, reward int64) *Transaction {
	panic("TODO: implement NewCoinbaseTX")
}

// FindUTXOs scans all transactions and returns every output currently
// unspent and locked to address.
func FindUTXOs(transactions []*Transaction, address string) []TxOutput {
	panic("TODO: implement FindUTXOs")
}

// BalanceOf sums the Value of every UTXO belonging to address.
func BalanceOf(transactions []*Transaction, address string) int64 {
	panic("TODO: implement BalanceOf using FindUTXOs")
}

// NewTransaction builds a transaction spending enough of from's UTXOs
// to cover amount, paying to, with a change output back to from if
// needed. Returns an error if from's balance is insufficient.
func NewTransaction(transactions []*Transaction, from, to string, amount int64) (*Transaction, error) {
	panic("TODO: implement NewTransaction")
}
