# Chapter 32: Building Transactions in Go

Chapters 29-31 built the mental model — transactions as signed instructions, UTXOs as indivisible chunks of value, and why GoChain tracks them instead of running balances. This chapter turns all of it into real, compiling Go code: the `Transaction`, `TxInput`, and `TxOutput` types, a minimal `UTXOSet` that scans the chain, and `NewTransaction`, the function that builds a spendable, (as-yet-unsigned) transaction from a wallet, a recipient, and an amount.

## Table of Contents

1. [Recap: What a Transaction Needs to Contain](#1-recap-what-a-transaction-needs-to-contain)
2. [The Wallet Type, at Last](#2-the-wallet-type-at-last)
3. [Defining TxOutput](#3-defining-txoutput)
4. [Defining TxInput](#4-defining-txinput)
5. [Defining Transaction](#5-defining-transaction)
6. [Serializing and Hashing a Transaction](#6-serializing-and-hashing-a-transaction)
7. [A Minimal UTXOSet](#7-a-minimal-utxoset)
8. [Selecting UTXOs to Cover an Amount](#8-selecting-utxos-to-cover-an-amount)
9. [Building NewTransaction](#9-building-newtransaction)
10. [A Note on Coinbase Transactions](#10-a-note-on-coinbase-transactions)
11. [Trying It Out](#11-trying-it-out)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Recap: What a Transaction Needs to Contain

From Chapter 29, a trustworthy transaction needs: proof of where funds come from, proof of authorization (a signature — built in Chapter 33), and a destination. From Chapter 30, GoChain represents "where funds come from" as references to specific earlier UTXOs (**inputs**), and "destination" as newly created UTXOs (**outputs**), possibly including change back to the sender.

This chapter defines exactly three Go types to capture that shape — `core.TxOutput`, `core.TxInput`, and `core.Transaction` — matching this course's shared type contract precisely, since every later volume (mempool, fees, wallets, networking, storage indexing) depends on these exact field names and types never changing again.

---

## 2. The Wallet Type, at Last

Every transaction needs a sender, and every sender needs to be represented by *something* in Go code — a value with a key pair to sign with and an address to receive change at. This course's shared contract has quietly assumed a `*wallet.Wallet` type since Chapter 13, and Volume 6 later grows it into a full hierarchical-deterministic wallet with seed phrases and encrypted files — but nothing has actually written down its minimal shape yet. Before this chapter can build a single real transaction, that gap needs closing.

```go
package wallet

import (
	"github.com/you/gochain/crypto"
)

// Wallet is the minimal thing this volume (and Volume 6, which extends it
// with HD derivation, seed phrases, and encrypted storage) needs: a key
// pair to sign with, and a way to derive the human-readable address that
// key pair controls.
type Wallet struct {
	KeyPair *crypto.KeyPair
}

// New generates a brand-new key pair (Chapter 13) and wraps it in a Wallet.
func New() *Wallet {
	kp, err := crypto.NewKeyPair()
	if err != nil {
		// A failure here means the underlying crypto/rand source failed --
		// not something calling code can meaningfully recover from, so we
		// panic rather than return an error every caller would have to
		// check and immediately treat as fatal anyway.
		panic(err)
	}
	return &Wallet{KeyPair: kp}
}

// Address derives this wallet's Base58, checksummed address from its
// public key -- exactly the format Chapter 14 built: hash the public key,
// prepend a version byte, append a checksum, Base58-encode the result.
func (w *Wallet) Address() string {
	pubKeyHash := crypto.PubKeyHash(w.KeyPair.PublicKey)

	versionedPayload := append([]byte{0x00}, pubKeyHash...)
	checksum := crypto.Checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)

	return string(crypto.Base58Encode(fullPayload))
}

// HashPubKey is a wallet-facing name for crypto.PubKeyHash (Chapter 14),
// kept here so transaction code (Section 4 onward) can talk about "whose
// public key hash is this" without reaching into the crypto package
// directly for what is, conceptually, a wallet-identity question.
func HashPubKey(pubKey []byte) []byte {
	return crypto.PubKeyHash(pubKey)
}

// PubKeyHashFromAddress reverses Address above: given a human-readable
// GoChain address, it Base58-decodes it and strips the leading version
// byte and trailing 4-byte checksum added there, leaving the raw public
// key hash embedded inside. TxOutput.Lock (Section 3) uses this to turn
// "pay this address" into the PubKeyHash a TxOutput actually stores.
func PubKeyHashFromAddress(address string) []byte {
	versionedHash := crypto.Base58Decode([]byte(address))
	// versionedHash is [1 version byte][20-or-so byte pubkey hash][4 byte checksum].
	return versionedHash[1 : len(versionedHash)-4]
}
```

`Wallet` itself is deliberately thin: one field, `KeyPair`, holding exactly the `*crypto.KeyPair` Chapter 13 already knows how to generate, sign with, and verify against. `New()` is the constructor every later chapter calls (starting with Chapter 36's CLI) whenever a fresh identity is needed — it generates a key pair and does nothing else, which is exactly right for this volume; Volume 6 is where wallets grow seed phrases, multiple derived addresses, and encrypted files, all *around* this same core `KeyPair` field, not replacing it. `Address()` is the other half of Chapter 14's work made concrete on a real type: hash the public key, glue on a version byte and checksum, Base58-encode. `HashPubKey` and `PubKeyHashFromAddress` are small, named conveniences this chapter's transaction code (starting in Section 3) leans on constantly — the former whenever we have a raw public key and need the hash that identifies it, the latter whenever we have a human-typed address string and need to recover the raw hash locked inside it.

With `Wallet` now real, `core.Blockchain` can grow the one convenience method every later chapter's demos assume exists — a direct way to ask "how much does this address have?" without manually constructing a `UTXOSet` first:

```go
package core

// BalanceOf sums the Value of every UTXO currently belonging to address,
// by building this chapter's UTXOSet against the receiver and delegating
// to FindUTXO (Section 7). It exists purely for convenience -- callers
// that already have a *UTXOSet handy (as Chapter 56's indexed version
// does) can call FindUTXO directly and skip rebuilding one here.
func (bc *Blockchain) BalanceOf(address string) int64 {
	pubKeyHash := wallet.PubKeyHashFromAddress(address)
	utxoSet := UTXOSet{Blockchain: bc}

	var balance int64
	for _, out := range utxoSet.FindUTXO(pubKeyHash) {
		balance += out.Value
	}
	return balance
}
```

`Blockchain.BalanceOf` is a thin wrapper, not a new capability: it decodes the given address back into a public key hash, builds a `UTXOSet` pointed at the receiver blockchain, and sums whatever `FindUTXO` (Section 7) returns. Every later chapter that calls `bc.BalanceOf(address)` or `chain.BalanceOf(address)` — Chapter 36's CLI, Chapter 37's major project demo, and this course's own "What You Will Build" preview — is calling exactly this method. Volume 8 introduces a separate, *indexed* `storage.UTXOSet.BalanceOf`, which answers the identical question far faster by not rescanning the whole chain on every call; the two never conflict, because by the time Volume 8 replaces the underlying UTXO lookup, this method's simple scan-based version is exactly what gets swapped out underneath it, same as the rest of this chapter's `UTXOSet`.

---

## 3. Defining TxOutput

An output is the simpler of the two halves: it just says how much value it holds and who is allowed to spend it. "Who is allowed to spend it" is represented not by a plain address string, but by a **public key hash** — the same hashed-public-key value that Chapter 14 used to build GoChain's address format in the first place.

```go
package core

// TxOutput represents a chunk of value created by a transaction, locked to
// whoever controls the private key behind PubKeyHash. It is exactly what
// Chapter 30 called a UTXO once no later transaction has spent it.
type TxOutput struct {
	Value      int64  // amount, in gochips
	PubKeyHash []byte // identifies who can spend this output
}
```

`Value` is the amount in gochips this output holds — an ordinary `int64`, since gochips are indivisible whole units for GoChain's purposes (real systems like Bitcoin use fractional units; we keep GoChain's unit whole for simplicity). `PubKeyHash` is *not* the recipient's raw public key — it's the hash of it, the same value that gets Base58-encoded (with a version byte and checksum) into the human-readable address you built in Chapter 14. Storing the hash rather than the address string means we never have to re-decode a Base58 string just to check "does this signature's public key match who this output is locked to?" — we can compare raw bytes directly.

A couple of small helper methods make outputs pleasant to work with:

```go
// Lock assigns this output to whoever holds the private key behind the
// given address, by deriving and storing that address's public key hash.
// wallet.PubKeyHashFromAddress reverses the Base58-decode-and-checksum-strip
// steps built in Chapter 14 to recover the raw hash from a human-readable
// address string.
func (out *TxOutput) Lock(address string) {
	out.PubKeyHash = wallet.PubKeyHashFromAddress(address)
}

// IsLockedWithKey reports whether this output is spendable by whoever
// controls the given public key hash -- used when scanning the chain for
// UTXOs that belong to a particular address (Section 7).
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// NewTXOutput is a small convenience constructor: build and lock an output
// to an address in one call, instead of two separate steps.
func NewTXOutput(value int64, address string) *TxOutput {
	txo := &TxOutput{Value: value}
	txo.Lock(address)
	return txo
}
```

`Lock` and `NewTXOutput` both depend on `wallet.PubKeyHashFromAddress`, a small utility that already exists from Chapter 14's address work — it decodes a Base58 address, verifies its checksum, and returns the raw public key hash embedded inside it. `IsLockedWithKey` is a plain byte-slice comparison; we reach for it constantly once the UTXO set needs filtering by owner.

---

## 4. Defining TxInput

An input is a *reference*, not a value — it doesn't say how much it's worth (that would be redundant, since the output it points at already says that); it says *which* earlier output it's spending, and carries the proof that the spender is entitled to.

```go
// TxInput references one specific, previously-unspent output and proves the
// spender is authorized to consume it. Signature and PublicKey start out
// nil when a transaction is first built (Section 9) and are filled in by
// Sign() in Chapter 33 -- an unsigned transaction is only a plan, not yet
// an authorized instruction.
type TxInput struct {
	TxID      []byte // ID of the transaction holding the output being spent
	OutIndex  int    // which output in that transaction
	Signature []byte
	PublicKey []byte
}
```

`TxID` and `OutIndex` together are the "which UTXO" part — they say "the output at index `OutIndex` of the transaction whose ID is `TxID`," which is exactly the UTXO#1, UTXO#2, ... references Chapter 30's worked example used informally. `Signature` and `PublicKey` are the "proof of authorization" part: `PublicKey` is the raw public key of whoever is spending (not a hash — the *full* key is needed here so any node can run signature verification against it, then separately confirm that key's hash matches the output's `PubKeyHash`), and `Signature` is the ECDSA signature Chapter 33 computes.

One small helper pulls its weight repeatedly:

```go
// UsesKey reports whether this input is authorized by the holder of the
// given public key hash -- used, symmetrically with TxOutput.IsLockedWithKey,
// when scanning the chain to figure out which UTXOs a given address has
// already spent.
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PublicKey)
	return bytes.Equal(lockingHash, pubKeyHash)
}
```

`wallet.HashPubKey` is the same hashing step Chapter 14 used to build addresses in the first place, applied here to an input's public key so we can compare it directly against an output's stored `PubKeyHash`.

---

## 5. Defining Transaction

`Transaction` ties everything together: a list of inputs (what's being spent), a list of outputs (what's being created), a timestamp, and its own ID.

```go
// Transaction is a signed instruction moving value: its Inputs consume
// existing UTXOs, and its Outputs create new ones. ID uniquely identifies
// this transaction and is what later transactions' TxInput.TxID fields
// reference when spending one of its outputs.
type Transaction struct {
	ID        []byte
	Inputs    []TxInput
	Outputs   []TxOutput
	Timestamp int64
}
```

`ID` deserves a moment of attention: it is computed by hashing the transaction's own contents (Section 6), the same "canonical serialize, then hash" pattern `core.Block` has used since Chapter 09. Once a transaction has an ID, it becomes something other transactions can point back at — `TxInput.TxID` is exactly this value. `Timestamp` records when the transaction was built, following the same convention `core.Block.Timestamp` already uses.

---

## 6. Serializing and Hashing a Transaction

Before we can compute an ID, we need a way to turn a `Transaction` into bytes deterministically — the same canonical-serialization requirement Chapter 09 introduced for blocks applies here too: two transactions with identical contents must serialize to identical bytes, or their hashes (and therefore their IDs) would differ for no good reason.

```go
import (
	"bytes"
	"encoding/gob"

	"github.com/you/gochain/crypto"
)

// Serialize turns a transaction into a deterministic byte slice using gob,
// the same encoding Chapter 07 chose for GoChain's early, simplicity-first
// storage format.
func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		// Encoding a plain struct of slices and ints cannot realistically
		// fail; panicking here surfaces a programming error immediately
		// rather than silently producing a corrupt transaction.
		panic(err)
	}
	return buf.Bytes()
}

// Hash computes this transaction's content fingerprint, with its own ID
// field temporarily cleared -- otherwise a transaction's hash would depend
// on its own hash, which is circular. This is the same function used both
// to assign a fresh transaction's ID (Section 9) and, in Chapter 33, to
// compute what a signature actually signs.
func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	return crypto.Hash(txCopy.Serialize())
}
```

`Serialize` reuses `encoding/gob`, Chapter 07's chosen format for GoChain's early volumes, to turn the whole struct into bytes in one call. `Hash` is where the subtlety lives: it copies the transaction (so the original is untouched), blanks out the copy's `ID` field, serializes *that*, and hashes the result with `crypto.Hash` from Volume 2. Blanking `ID` first matters because otherwise you'd be trying to hash a value that includes its own hash — an obvious contradiction. We'll lean on this exact `Hash` method again in Chapter 33 for a related but distinct reason: computing what bytes a signature actually covers.

---

## 7. A Minimal UTXOSet

To build a new transaction, a wallet needs to answer two questions: "what UTXOs does this address currently own?" and "which of them should I use to cover this amount?" Answering both requires knowing the current UTXO set — but nothing has computed or stored one yet. `core.Blockchain` only knows how to hold and validate blocks (Volumes 3-4); it has no index of which outputs are still unspent.

For this volume, we build the simplest thing that works: a `UTXOSet` that, whenever asked, walks the *entire* chain from the newest block back to genesis, keeping track of which outputs have been referenced as inputs along the way, and returning whatever's left over. Volume 8 replaces this with a real, persistently-maintained index (so this expensive full scan doesn't have to happen on every lookup) — but the *meaning* of "UTXO set" defined here doesn't change; only its performance does.

```go
package core

import "encoding/hex"

// UTXOSet answers questions about currently-unspent outputs by scanning the
// underlying blockchain. This is intentionally the simplest correct
// implementation, not the fastest one -- Volume 8 (Chapter 56) replaces the
// full chain scan with an incrementally-maintained index, behind this same
// conceptual interface.
type UTXOSet struct {
	Blockchain *Blockchain
}

// findUnspentTransactions walks the chain from the tip back to genesis and
// returns, for each transaction, only the outputs that no later input (in
// any transaction anywhere on the chain) has referenced yet.
func (u UTXOSet) findUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int) // txID (hex) -> spent output indexes

	for _, block := range u.Blockchain.blocks {
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				// Skip this output if we've already seen, in a later
				// (already-scanned) block, an input that spends it.
				if spent, ok := spentTXOs[txID]; ok {
					for _, spentOut := range spent {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// A coinbase transaction (Chapter 37) has no real inputs to
			// spend, so skip the input bookkeeping for it entirely.
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutIndex)
					}
				}
			}
		}
	}

	return unspentTXs
}

// FindUTXO returns every currently-unspent TxOutput belonging to the given
// public key hash -- exactly the pile of "bills and coins" from Chapter 30,
// used to compute a balance by summing their Value fields.
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TxOutput {
	var utxos []TxOutput
	for _, tx := range u.findUnspentTransactions(pubKeyHash) {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				utxos = append(utxos, out)
			}
		}
	}
	return utxos
}
```

`findUnspentTransactions` is the workhorse: it walks every block, and within each block, every transaction. For each transaction's outputs, it checks whether some *already-recorded* spending input has claimed that exact output — if so, it's skipped, because it's no longer a UTXO. Then it records every input's own reference as "spent," so later (in scan order) transactions checking their own outputs against `spentTXOs` see it. `FindUTXO` is a thin wrapper that filters down to just the outputs actually belonging to `pubKeyHash` and returns them directly — this is precisely the sum-of-UTXOs balance computation from Chapter 30, Section 5, just made real: a wallet's `BalanceOf` would call this and add up every returned output's `Value`.

---

## 8. Selecting UTXOs to Cover an Amount

Computing a balance is only half of what a wallet needs. To *spend*, it needs to pick a specific subset of UTXOs whose combined value is at least the amount being sent — exactly the coin-selection problem Chapter 30's worked example glossed over informally.

```go
// FindSpendableOutputs greedily selects enough of pubKeyHash's UTXOs to
// cover amount, stopping as soon as it has enough. It returns the total
// value accumulated and a map from transaction ID (hex) to the specific
// output indexes selected within that transaction -- everything
// NewTransaction (Section 9) needs to build a set of TxInputs.
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int64) (int64, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	var accumulated int64

	unspentTXs := u.findUnspentTransactions(pubKeyHash)

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
```

This is a deliberately simple **greedy** selection strategy: walk the address's UTXOs in whatever order the scan produces them, adding each one to the selected set until the running total meets or exceeds `amount`, then stop immediately (the labeled `break Work` exits both the outer and inner loop at once, the moment enough value has been gathered). It's not the most sophisticated coin-selection algorithm real wallets use — Chapter 30's exercises touch on smarter strategies — but it is correct: it never selects more UTXOs than necessary to reach the target, and if the address's total balance is less than `amount`, `accumulated` simply comes back lower than `amount`, which `NewTransaction` checks explicitly next.

---

## 9. Building NewTransaction

Everything is now in place to build a new, unsigned transaction: look up enough spendable UTXOs from the sender's wallet, turn each selected UTXO into a `TxInput`, create an output paying the recipient, add a change output if there's leftover value, and compute the transaction's ID.

```go
package core

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/you/gochain/wallet"
)

// NewTransaction builds a new, unsigned transaction sending amount gochips
// from the wallet's address to the recipient's, spending whatever UTXOs are
// necessary from utxoSet. The returned transaction still needs Sign()
// (Chapter 33) called on it before any node will accept it -- an unsigned
// transaction is only a plan, not yet an authorized instruction.
func NewTransaction(from *wallet.Wallet, to string, amount int64, utxoSet *UTXOSet) (*Transaction, error) {
	var inputs []TxInput
	var outputs []TxOutput

	pubKeyHash := wallet.HashPubKey(from.KeyPair.PublicKey)
	accumulated, validOutputs := utxoSet.FindSpendableOutputs(pubKeyHash, amount)

	if accumulated < amount {
		return nil, fmt.Errorf("insufficient funds: have %d gochips, need %d", accumulated, amount)
	}

	// Turn every selected UTXO into a TxInput. Signature and PublicKey are
	// left nil here -- Sign() fills them in per input in Chapter 33.
	for txIDHex, outIdxs := range validOutputs {
		txID, err := hex.DecodeString(txIDHex)
		if err != nil {
			return nil, fmt.Errorf("decoding transaction id: %w", err)
		}
		for _, outIdx := range outIdxs {
			inputs = append(inputs, TxInput{
				TxID:     txID,
				OutIndex: outIdx,
			})
		}
	}

	// The actual payment.
	outputs = append(outputs, *NewTXOutput(amount, to))

	// Change: whatever's left over from the UTXOs we had to consume goes
	// back to the sender, exactly like getting $8 back from a $20 bill for
	// a $12 item (Chapter 30, Section 4).
	if accumulated > amount {
		change := accumulated - amount
		outputs = append(outputs, *NewTXOutput(change, from.Address()))
	}

	tx := &Transaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: time.Now().Unix(),
	}
	tx.ID = tx.Hash()

	return tx, nil
}
```

Walking through it: first, `wallet.HashPubKey(from.KeyPair.PublicKey)` derives the sender's public key hash directly from their key pair, so we can look up their UTXOs without needing to re-derive their address string. `utxoSet.FindSpendableOutputs` (Section 8) then does the coin selection, returning both how much value was gathered and exactly which outputs to spend. If that's less than `amount`, we return a clear error rather than building a transaction that could never be valid — better to fail immediately in `NewTransaction` than to let an unspendable transaction reach a miner or another node later. Each selected UTXO becomes a bare `TxInput` (no signature yet — that's Chapter 33's job entirely). We append one output paying the recipient using the `NewTXOutput` convenience constructor from Section 3, and, if we had to spend more than `amount`, a second output sending the leftover back to `from.Address()` as change. Finally, we compute the transaction's `ID` by calling `tx.Hash()` (Section 6) once every other field is set.

Note that `from` here is a `*wallet.Wallet` — the minimal type Section 2 just built. For our purposes, it needs exactly two things: a `KeyPair *crypto.KeyPair` field (so we can read the sender's public key) and an `Address() string` method (so we know where to send change). Volume 6 grows this same type substantially — seed phrases, HD derivation, encrypted files — but never changes these two essentials, so nothing built on top of them here ever needs to change later.

---

## 10. A Note on Coinbase Transactions

You may have noticed `findUnspentTransactions` in Section 7 checking `tx.IsCoinbase()` before processing a transaction's inputs. A **coinbase transaction** is a special kind of transaction, with *no inputs at all*, used to reward whoever mines a block with brand-new gochips — it's the very first transaction in any block, and the very first UTXO in this volume's worked example (Chapter 30, Section 6) implicitly came from one. Since it has no inputs, there is nothing to sign or verify against a previous output.

```go
// IsCoinbase reports whether this transaction is a coinbase (block reward)
// transaction: one with exactly one input, whose TxID is empty. Coinbase
// transactions create new gochips out of nothing, by design -- they don't
// spend any earlier output, so they need no signature verification against
// a previous output's PubKeyHash.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0
}

// NewCoinbaseTX builds a coinbase transaction paying reward gochips to
// toAddress, with a single input that references nothing (TxID and
// OutIndex are left at their zero values) since there is no earlier output
// being spent.
func NewCoinbaseTX(toAddress string, reward int64) *Transaction {
	txIn := TxInput{TxID: []byte{}, OutIndex: -1}
	txOut := NewTXOutput(reward, toAddress)

	tx := &Transaction{
		Inputs:    []TxInput{txIn},
		Outputs:   []TxOutput{*txOut},
		Timestamp: time.Now().Unix(),
	}
	tx.ID = tx.Hash()
	return tx
}
```

We define `NewCoinbaseTX` here, alongside the rest of `Transaction`'s constructors, so the type is complete and every later chapter can rely on it existing. Its full story — *why* mining a block rewards you with new gochips, and how this ties together the entire send-and-receive lifecycle — is properly told in Chapter 37, where you'll put it to work end to end. For now, just note its shape: one input with an empty `TxID` (nothing to reference), and one output paying the miner.

---

## 11. Trying It Out

A short test exercises the whole pipeline built in this chapter, short of signing (Chapter 33 finishes the job):

```go
package core

import "testing"

func TestNewTransaction_InsufficientFunds(t *testing.T) {
	bc := NewBlockchainForTest(t) // test helper: a fresh chain with one
	                              // coinbase-funded block already mined
	utxoSet := &UTXOSet{Blockchain: bc}

	sender := wallet.NewForTest(t) // test helper: a throwaway wallet with
	                                // no funds on this chain at all

	_, err := NewTransaction(sender, "someAddress...", 10, utxoSet)
	if err == nil {
		t.Fatal("expected an error building a transaction with no funds, got nil")
	}
}
```

`TestNewTransaction_InsufficientFunds` deliberately builds a transaction for a wallet with zero balance on the test chain, and asserts `NewTransaction` returns an error rather than a broken, unspendable transaction — exactly the check Section 9 added right after coin selection. A companion happy-path test (left as Exercise 6) would mine a coinbase reward to a wallet first, then build a transaction spending part of it and assert the resulting inputs, outputs, and change all match Chapter 30's arithmetic exactly.

---

## Summary

- `core.TxOutput` holds an amount and a `PubKeyHash` identifying who can spend it; `Lock` and `NewTXOutput` derive that hash from a human-readable address using Chapter 14's encoding.
- `core.TxInput` references a specific earlier output (`TxID` + `OutIndex`) and carries `Signature`/`PublicKey` fields that start out empty and are filled in by `Sign()` in Chapter 33.
- `core.Transaction` bundles inputs, outputs, a timestamp, and an `ID` computed by hashing a copy of itself with `ID` blanked out, avoiding a circular self-referential hash.
- `UTXOSet` is a minimal, correctness-first helper that scans the entire chain on every call to find which outputs are still unspent — Volume 8 replaces the scan with a real index without changing what "UTXO set" means.
- `FindSpendableOutputs` greedily selects just enough UTXOs to cover a requested amount, exactly answering the coin-selection question Chapter 30 left informal.
- `NewTransaction` ties it all together: select UTXOs, build inputs, create a payment output and (if needed) a change output, and compute the transaction's ID — returning an explicit error if the sender's balance is insufficient.
- `NewCoinbaseTX` and `IsCoinbase` are defined here to complete the type, but their full role in rewarding miners is told properly in Chapter 37.

---

## Exercises

### Easy

1. Explain why `Transaction.Hash()` has to blank out the `ID` field before serializing and hashing, rather than hashing the transaction exactly as it is.
2. `TxOutput.PubKeyHash` stores a hash of a public key, not a raw address string. Explain, in your own words, why comparing raw hashed bytes (as `IsLockedWithKey` does) is more convenient than comparing Base58-encoded address strings directly.
3. Trace through `FindSpendableOutputs` by hand for a UTXO set of `{20 -> Alice, 5 -> Alice, 30 -> Alice}` (in that scan order) and a requested amount of 22. Which UTXOs get selected, and what is the final `accumulated` value?

### Medium

4. `findUnspentTransactions` builds a `spentTXOs` map as it scans, but only ever *adds* entries to it — it never removes any. Explain why removal is unnecessary here, referencing how the function decides an output has already been spent.
5. Extend `NewTransaction` (on paper or in code) to also return the total fee implied by the transaction it built, where fee is defined as `accumulated - amount` minus whatever change output was created (in the current implementation, this fee is always zero, since `accumulated - amount` becomes exactly the change output). What additional parameter or logic would need to change for a non-zero fee to be possible? (Chapter 35 covers this properly — this is a chance to reason about it yourself first.)
6. Write the "happy path" companion test referenced in Section 11: build a test blockchain, mine a coinbase transaction to a test wallet, build a `UTXOSet` from that chain, call `NewTransaction` to send part of the balance to a second test wallet, and assert the resulting transaction's inputs, outputs, and computed change all match what Chapter 30's arithmetic would predict.

### Hard

7. `FindSpendableOutputs`'s greedy selection strategy can, depending on scan order, select more individual UTXOs than strictly necessary (for example, selecting three small UTXOs when one larger one alone would have covered the amount). Rewrite it to prefer the fewest possible inputs: first check if any single UTXO alone covers the amount, and if so use only that one; otherwise fall back to the greedy approach. Discuss one tradeoff your new version makes compared to the original.
8. The current `UTXOSet` rescans the entire blockchain from scratch on every single call to `FindSpendableOutputs` or `FindUTXO`. Estimate, with rough big-O reasoning, how this scales as the number of blocks and transactions on the chain grows, and describe in your own words (without writing the real implementation — that's Chapter 56) what kind of incremental index would let you avoid rescanning unchanged history on every call.
9. `NewTransaction` currently has no way to detect if it accidentally selects the *same* UTXO twice within one call (which should be impossible given how `FindSpendableOutputs` is written, but imagine a future bug introduces this). Write a validation check, added at the end of `NewTransaction` before returning, that would detect and reject a transaction whose `Inputs` slice references the same `(TxID, OutIndex)` pair more than once, and explain what real-world problem this specific bug would cause if it reached the network unnoticed.
