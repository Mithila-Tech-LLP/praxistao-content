# Task 05: Transactions and the UTXO Model

## What you will build

A transaction type using the UTXO (Unspent Transaction Output) model — the same accounting approach Bitcoin uses — plus a function that scans a set of transactions to compute how much a given address currently owns.

## Concepts

### Coins as discrete outputs, not a running balance

Instead of tracking "Alice has 50 coins" as a single mutable number, the UTXO model tracks individual, indivisible outputs: "this specific 50-coin output belongs to Alice." Spending means consuming one or more outputs entirely as inputs to a new transaction, and creating new outputs — including a change output back to yourself if you spent more than you needed.

```
  Alice has one UTXO worth 50.
  She wants to pay Bob 20.

  Transaction:
    Input:   the 50-coin UTXO (consumed entirely)
    Outputs: 20 -> Bob        (the payment)
             30 -> Alice      (change, just like getting $30 back from a $50 bill)
```

### Balance = sum of your unspent outputs

To find out how much an address owns, walk every transaction ever made, track which outputs have been referenced as inputs (spent) somewhere, and sum the `Value` of whatever's left that's locked to that address.

## Interface to implement

```go
type TxOutput struct {
	Value   int64
	Address string // simplified: the owner's address string directly
}

type TxInput struct {
	TxID     []byte
	OutIndex int
	Address  string // simplified: who is spending, no signature yet
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// NewCoinbaseTX creates a reward transaction with no real inputs,
// paying reward coins to address -- used to introduce new coins
// when a block is mined.
func NewCoinbaseTX(address string, reward int64) *Transaction

// FindUTXOs scans all transactions and returns every output currently
// unspent and locked to address.
func FindUTXOs(transactions []*Transaction, address string) []TxOutput

// BalanceOf sums the Value of every UTXO belonging to address.
func BalanceOf(transactions []*Transaction, address string) int64

// NewTransaction builds a transaction spending enough of from's UTXOs to
// cover amount, paying to, with a change output back to from if needed.
// Returns an error if from's balance is insufficient.
func NewTransaction(transactions []*Transaction, from, to string, amount int64) (*Transaction, error)
```

## Hints

- A coinbase transaction has exactly one input with an empty `TxID` — that's how you recognize "this one doesn't spend anything real."
- To find unspent outputs: build a set of `(txID, outIndex)` pairs referenced by every input across all transactions, then return only the outputs *not* in that set.
- `NewTransaction`'s coin selection can be simple and greedy: walk the address's UTXOs in any order, accumulate until you've covered `amount`, then stop.
- Compute each new transaction's `ID` with the `Hash` function from Task 01, applied to a deterministic serialization of its inputs and outputs.

## Run the tests

```bash
cd starter/task-05-transactions-and-utxos
go test ./...
```

All tests must pass before moving to Task 06.
