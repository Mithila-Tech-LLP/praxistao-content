# Chapter 33: Signing and Verifying Transactions

`NewTransaction` from Chapter 32 builds a transaction's shape — inputs, outputs, an ID — but leaves every input's `Signature` and `PublicKey` empty. An unsigned transaction is just a plan; anyone could claim any UTXO as their own input if nothing forced them to prove it. This chapter closes that gap: `Transaction.Sign()` and `Transaction.Verify()`, built on the ECDSA primitives from Volume 2, plus a careful look at exactly which bytes get signed and why it has to be a trimmed copy rather than the whole transaction.

## Table of Contents

1. [Why Signing Every Input Matters](#1-why-signing-every-input-matters)
2. [What Exactly Gets Signed](#2-what-exactly-gets-signed)
3. [Trimmed Copy vs. Full Transaction](#3-trimmed-copy-vs-full-transaction)
4. [The Malleability Problem This Avoids](#4-the-malleability-problem-this-avoids)
5. [Implementing TrimmedCopy](#5-implementing-trimmedcopy)
6. [Implementing Transaction.Sign](#6-implementing-transactionsign)
7. [Implementing Transaction.Verify](#7-implementing-transactionverify)
8. [A Full Walkthrough](#8-a-full-walkthrough)
9. [What's Next](#9-whats-next)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Signing Every Input Matters

Recall Chapter 32's `TxInput`: it references a UTXO by `TxID` and `OutIndex`, and carries `Signature` and `PublicKey` fields that `NewTransaction` leaves empty. Without those two fields filled in and checked, here is exactly what would go wrong: anyone on the network could construct a transaction whose input simply *names* one of your UTXOs — say, "spend Alice's 30-gochip output, pay it to my own address" — with no proof whatsoever that they control Alice's private key. Since a UTXO's only defense is "whoever can prove they control it," an unsigned or unverified input is equivalent to an unlocked door: the reference alone is enough to walk in.

**Signing** an input means using the private key that controls the referenced UTXO to produce a signature over specific data related to this transaction — proof that whoever built this transaction genuinely holds that private key. **Verifying** means any node, holding only the corresponding *public* key (which the input itself carries, openly, in its `PublicKey` field), can confirm that signature is valid without ever needing the private key. This is the exact ECDSA sign/verify relationship Chapter 12 explained conceptually and Chapter 13 implemented for real, in `gochain/crypto` — this chapter is where transactions finally put those primitives to work for their intended purpose.

Every input must be signed *separately*, even within the same transaction, because different inputs can reference UTXOs owned by different addresses (recall Chapter 30's Transaction 4, which combined two of Carol's UTXOs — if a transaction ever combined inputs from *different* owners, each owner's signature would need to prove authorization independently for their own input specifically).

---

## 2. What Exactly Gets Signed

The naive answer — "just sign the whole transaction" — turns out to be almost right, but wrong in one specific, important way. Let's build up to the correct answer carefully.

A signature is only meaningful relative to *some specific bytes*. `crypto.Sign(priv, data)` takes a private key and a byte slice, and produces a signature that `crypto.Verify(pubKey, data, signature)` can later confirm was produced by the private key matching `pubKey`, over that *exact* `data`. If even one byte of `data` changes between signing and verifying, verification fails — this is precisely the "avalanche effect" property hashing gave us back in Chapter 08, now doing its job inside a signature scheme.

So: what should `data` be, for a transaction input? Two requirements pull in the same direction:

- It must include enough of the transaction's own content that the signature actually commits to *this specific transaction* — its inputs and outputs — not just "some transaction, we don't know which."
- It must **not** include the very `Signature` bytes being produced, because a signature cannot include itself as part of what it signs (this is circular, exactly like `Transaction.Hash()` in Chapter 32 had to blank out `ID` before hashing, for the same underlying reason).

The answer GoChain uses, matching Bitcoin's original design, is to sign a **trimmed copy** of the transaction: a version with every input's `Signature` and `PublicKey` fields cleared out, except that, for the one input currently being signed, its `PublicKey` field is temporarily set to the `PubKeyHash` of the output it's spending. We hash *that* trimmed structure and sign the resulting hash.

---

## 3. Trimmed Copy vs. Full Transaction

```
                     FULL TRANSACTION (after signing)

  +----------------------------------------------------------------+
  |  ID: 0xAB12...                                                 |
  |  Inputs:                                                        |
  |    [0] TxID: 0x55.., OutIndex: 0,                               |
  |         Signature: 0x9F3E...  <-- present, real signature       |
  |         PublicKey:  0x04A1... <-- present, spender's real key   |
  |  Outputs:                                                        |
  |    [0] Value: 12, PubKeyHash: 0x77..                             |
  |    [1] Value: 8,  PubKeyHash: 0x22..                             |
  |  Timestamp: 177020...                                            |
  +----------------------------------------------------------------+


           TRIMMED COPY (what actually gets hashed and signed,
                 for input [0], during Sign() specifically)

  +----------------------------------------------------------------+
  |  ID: (blank -- not part of what's signed either)                |
  |  Inputs:                                                        |
  |    [0] TxID: 0x55.., OutIndex: 0,                               |
  |         Signature: nil        <-- cleared: can't sign itself     |
  |         PublicKey:  0x77..    <-- SET to the *previous output's* |
  |                                    PubKeyHash, standing in for   |
  |                                    the real key just for hashing |
  |  Outputs:                                                        |
  |    [0] Value: 12, PubKeyHash: 0x77..                             |
  |    [1] Value: 8,  PubKeyHash: 0x22..                             |
  |  Timestamp: 177020...                                            |
  +----------------------------------------------------------------+
```

Notice two things happening at once here, and it's easy to conflate them, so keep them separate in your head:

1. **Every input's `Signature` is cleared** in the trimmed copy — obviously, since we're computing what to sign, no signature exists yet.
2. **The input currently being processed gets its `PublicKey` field temporarily replaced** with the `PubKeyHash` from the output it's spending (not the spender's real public key — that field is used, for this one step only, to stand in for "which specific output's locking condition this input is satisfying"). Every *other* input's `PublicKey` stays cleared to `nil` during this step.

This construction means the signature for input `[0]` commits to: this exact set of outputs, this exact set of input references (`TxID`/`OutIndex` pairs, untouched), and specifically which output `[0]` claims to be unlocking — all without needing the final signature bytes to exist yet. Once the hash of this trimmed structure is computed and signed, the real `Signature` and `PublicKey` are written back into the *actual* transaction's input `[0]` (not the trimmed copy, which is thrown away after each input is processed).

---

## 4. The Malleability Problem This Avoids

**Transaction malleability** is the name for a specific class of bug: an attacker takes a valid, already-signed transaction and modifies *some* part of it — without invalidating its signature — producing a different-looking transaction that still spends the same funds the same way. If a transaction's ID depends on parts of the transaction that a signature doesn't actually cover, an attacker can change those parts, changing the transaction's ID (since `Hash()` covers the whole thing), while the signature — computed over some smaller, unrelated set of bytes — still checks out fine.

Here's concretely how signing "the whole transaction, signature fields included" would create this bug: suppose (incorrectly) that `Sign()` computed a signature over the transaction's serialized bytes *as they'd finally appear*, signature fields and all. To even get a value to sign, you'd need some placeholder for the signature before it exists — and whatever convention you pick for that placeholder, an attacker intercepting a broadcast-but-not-yet-mined transaction could potentially craft a *different*, still-validly-signed encoding of logically the same spend (for example, by exploiting how signature encoding itself has some flexibility — real ECDSA signatures have this property unless additional canonical-encoding rules are enforced) and rebroadcast it. Now two different transaction IDs exist that both spend the identical UTXOs to the identical destinations, and anything that was tracking "transaction ID 0xAB12 hasn't confirmed yet, so let's resend it" gets confused when a *different*-ID transaction (0xCD34, functionally identical) shows up mined instead.

Signing a trimmed copy with signature fields always cleared sidesteps this specific problem at its root: the signature never depends on itself, and by only ever populating `Signature` *after* hashing a version that has that field blanked, there's no ambiguity about what "the data being signed" was — it's always the input's own reference plus the spent output's `PubKeyHash`, plain and unambiguous, with no signature-encoding wiggle room feeding back into what gets hashed.

---

## 5. Implementing TrimmedCopy

```go
package core

// TrimmedCopy returns a copy of tx with every input's Signature and
// PublicKey cleared. Sign and Verify both build one of these, then
// temporarily set just one input's PublicKey at a time (Sections 6-7)
// before hashing -- never touching the real transaction's fields directly
// mid-computation.
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{
			TxID:     in.TxID,
			OutIndex: in.OutIndex,
			// Signature and PublicKey deliberately left nil.
		})
	}
	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{
			Value:      out.Value,
			PubKeyHash: out.PubKeyHash,
		})
	}

	return Transaction{
		ID:        tx.ID,
		Inputs:    inputs,
		Outputs:   outputs,
		Timestamp: tx.Timestamp,
	}
}
```

`TrimmedCopy` builds a brand-new `Transaction` value (not a pointer into the original, so mutating the copy can never accidentally corrupt the real transaction) with every input's `Signature` and `PublicKey` left at their zero value (`nil`), while outputs are copied over unchanged, since outputs never carry proof-of-authorization fields to begin with — only inputs do.

---

## 6. Implementing Transaction.Sign

```go
import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/you/gochain/crypto"
)

// Sign signs every input of tx using privKey, given prevTXs -- a lookup
// from hex-encoded transaction ID to the full Transaction that produced the
// output each input spends. prevTXs is how Sign discovers the PubKeyHash
// each input is claiming to unlock, without needing the whole blockchain in
// scope.
func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	// A coinbase transaction has no real inputs to authorize -- see
	// Chapter 32, Section 9.
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, input := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(input.TxID)]

		// Stand in the spent output's PubKeyHash as this input's
		// PublicKey field, just for computing this one input's hash --
		// see the trimmed-copy diagram in Section 3.
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PublicKey = prevTX.Outputs[input.OutIndex].PubKeyHash

		dataToSign := txCopy.Hash()

		// Restore the trimmed copy's input to a clean slate before moving
		// on to the next input in the loop.
		txCopy.Inputs[inID].PublicKey = nil

		signature, err := crypto.Sign(privKey, dataToSign)
		if err != nil {
			// In production code this error should propagate to the
			// caller; we keep Sign's signature error-free, matching this
			// course's shared type contract, and panic here since a
			// signing failure during otherwise-valid input means
			// something is fundamentally wrong with the key material.
			panic(err)
		}

		// Write the real signature and public key onto the ACTUAL
		// transaction's input -- never onto the throwaway trimmed copy.
		tx.Inputs[inID].Signature = signature
		tx.Inputs[inID].PublicKey = crypto.PublicKeyToBytes(&privKey.PublicKey)
	}
}
```

Walking through the loop: for each input, we look up `prevTX`, the earlier transaction holding the output this input spends, using `prevTXs` (a map the caller assembles beforehand — typically by looking up each referenced `TxID` in the blockchain or UTXO set). We temporarily set the trimmed copy's *matching* input's `PublicKey` to `prevTX.Outputs[input.OutIndex].PubKeyHash` — exactly the substitution the Section 3 diagram showed — then call `txCopy.Hash()` to get `dataToSign`. Immediately after, we clear that field back to `nil` so the *next* iteration starts from a clean trimmed copy (only one input's `PublicKey` is ever set at a time). We sign `dataToSign` with `crypto.Sign`, and write the resulting signature — along with the spender's actual public key, `crypto.PublicKeyToBytes(&privKey.PublicKey)` — onto the *real* transaction's input, never the trimmed copy, which is discarded once the loop ends.

---

## 7. Implementing Transaction.Verify

```go
// Verify checks that every input of tx carries a valid signature,
// authorizing the spend of the output it references, given the same
// prevTXs lookup Sign used to build those signatures in the first place.
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true // nothing to verify: no inputs, so no signatures needed
	}

	txCopy := tx.TrimmedCopy()

	for inID, input := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(input.TxID)]

		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PublicKey = prevTX.Outputs[input.OutIndex].PubKeyHash

		dataThatWasSigned := txCopy.Hash()

		txCopy.Inputs[inID].PublicKey = nil

		if !crypto.Verify(input.PublicKey, dataThatWasSigned, input.Signature) {
			return false
		}
	}

	return true
}
```

`Verify` mirrors `Sign` almost exactly, which is intentional and important: a verifier must reconstruct the *identical* `dataToSign` bytes that `Sign` originally hashed, or verification would fail even for a perfectly legitimate signature. It builds the same trimmed copy, substitutes the same spent-output `PubKeyHash` into the same input's `PublicKey` field, and hashes it — reproducing `dataToSign` from Section 6 exactly. Then it calls `crypto.Verify(input.PublicKey, dataThatWasSigned, input.Signature)`, using the *real* public key stored on the actual (untrimmed) input, to confirm the signature is mathematically valid for that data. If any single input fails this check, the whole transaction is rejected — `Verify` returns `false` the moment one input doesn't check out, since a transaction with even one unauthorized input is not trustworthy as a whole.

Note what `Verify` deliberately does *not* check: whether the referenced UTXO is currently *unspent*. That's a separate question — Chapter 34's mempool double-spend check — from whether the *signature* is valid. A transaction can carry a perfectly valid signature over an output that's already been spent by some other transaction; `Verify` alone can't and shouldn't know that, since it only has `prevTXs` (the transactions that *created* the referenced outputs) in scope, not the live UTXO set.

---

## 8. A Full Walkthrough

Let's trace signing and verifying one concrete transaction end to end, continuing Chapter 30's worked example: Alice spends her 30-gochip UTXO (from Transaction 1's change output) to pay Carol 30 gochips exactly (Chapter 30's Transaction 3).

```go
// Alice already has her key pair from wallet setup (Volume 6 covers this
// properly; for now, from.KeyPair.PrivateKey is simply available).
tx, err := core.NewTransaction(alice, carolAddress, 30, utxoSet)
if err != nil {
	log.Fatal(err)
}

// prevTXs maps this transaction's referenced input TxIDs to the full
// transactions that created those outputs -- here, just the one
// transaction that produced Alice's 30-gochip change output.
prevTXs := map[string]Transaction{
	hex.EncodeToString(tx.Inputs[0].TxID): *transactionThatCreatedAlicesChange,
}

tx.Sign(alice.KeyPair.PrivateKey, prevTXs)

fmt.Println("signature present:", tx.Inputs[0].Signature != nil) // true
fmt.Println("verifies:", tx.Verify(prevTXs))                     // true

// Now simulate tampering: an attacker changes the payment amount after
// the fact, hoping nobody notices.
tx.Outputs[0].Value = 3000 // attacker tries to redirect a much larger sum

fmt.Println("still verifies after tampering:", tx.Verify(prevTXs)) // false
```

The first `Verify` call succeeds because `tx.Hash()` (via the trimmed copy) reproduces exactly the bytes `Sign` originally signed, and `crypto.Verify` confirms Alice's signature is valid over them. After the attacker mutates `tx.Outputs[0].Value`, the trimmed copy's hash changes too — outputs are part of what gets hashed and signed (Section 2) — so `dataThatWasSigned` inside `Verify` no longer matches what Alice actually signed, and `crypto.Verify` correctly reports the signature as invalid. This is the entire point: a signature doesn't just prove "Alice approved *a* transaction" — it proves "Alice approved *these exact* inputs and outputs," and any change to either, however small, is caught.

---

## 9. What's Next

`Sign` and `Verify` complete the "proof of authorization" leg of Chapter 29's three-part trustworthiness checklist — combined with Chapter 32's inputs (proof of source) and outputs (destination), `core.Transaction` is now a fully self-contained, independently verifiable unit. What it still can't do on its own is stop the *same* legitimately-signed UTXO from being spent twice by two different, both perfectly-valid-looking transactions racing each other — that's a question about the current state of the UTXO set, not about any single transaction's internal correctness. Chapter 34 picks up exactly there, building the mempool and the double-spend check that watches for precisely this.

---

## Summary

- Every `TxInput` must carry a signature proving whoever built the transaction controls the private key behind the UTXO it references — otherwise anyone could name any address's UTXO as their own input.
- Signatures cannot cover their own bytes, so `Sign` and `Verify` both operate on a **trimmed copy** of the transaction, with every input's `Signature` cleared.
- For the specific input being processed, the trimmed copy's `PublicKey` field is temporarily set to the *spent output's* `PubKeyHash`, so the signature commits to exactly which output this input claims to unlock.
- Signing the whole transaction (signature bytes included) would create a **malleability** bug: an attacker could produce a different-looking, still-validly-signed encoding of the same spend, changing the transaction's ID without invalidating its signature.
- `TrimmedCopy` builds the stripped-down structure; `Sign` hashes it per-input (substituting in the right `PubKeyHash` each time) and writes the real signature and public key onto the actual transaction; `Verify` reconstructs the identical hash and checks it against each input's stored signature and public key.
- `Verify` only checks signature validity — it does not (and cannot, from the data it has) check whether a referenced UTXO has already been spent elsewhere; that's the mempool's job, next.
- The full walkthrough demonstrated that even a single-field tamper (changing an output's `Value` after signing) is caught immediately by `Verify`, because outputs are part of what the signature covers.

---

## Exercises

### Easy

1. Explain, in your own words, why `Sign` clears a trimmed copy's `Signature` field before hashing, rather than hashing the transaction with its final signature already in place.
2. What specific field does the trimmed copy temporarily set on the input currently being processed, and what value does it get set to? Why that value, specifically, instead of the spender's real public key?
3. In the Section 8 walkthrough, changing `tx.Outputs[0].Value` after signing made `Verify` fail. Would changing `tx.Timestamp` after signing also make `Verify` fail? Justify your answer by referring to what `TrimmedCopy` includes.

### Medium

4. Trace through `Sign` by hand for a transaction with *two* inputs, spending UTXOs from two different previous transactions. Write out, step by step, what the trimmed copy looks like at the moment each input's signature is being computed (which fields are set, which are nil).
5. `Verify` returns `false` immediately on the first invalid input rather than checking all inputs and reporting how many failed. Discuss one advantage and one disadvantage of this "fail fast" behavior compared to checking every input and returning a full list of which ones failed.
6. Suppose someone proposed simplifying `Sign`/`Verify` by hashing the *entire untrimmed* transaction, including whatever signatures already exist on *other* inputs at the time each input is processed (rather than always clearing all signatures first). Explain what would go wrong, referencing the trimmed-copy diagram in Section 3.

### Hard

7. Section 4 describes malleability abstractly. Design a concrete (hypothetical) scenario, using GoChain's own types, where a signature that only covered a transaction's inputs (not its outputs) would let an attacker redirect funds to a different address without invalidating the signature. Be specific about what the attacker changes and why the (flawed) signature would still check out.
8. Research (general background knowledge is fine) how real Bitcoin eventually addressed a related malleability issue with a change called SegWit (Segregated Witness) — you don't need deep protocol details, just the core idea of separating signature data from the data used to compute a transaction's ID. Write 200-300 words comparing that approach's goal to what this chapter's trimmed-copy approach already achieves for GoChain, and note one way SegWit's approach differs from ours.
9. Extend `Transaction.Verify` (on paper or in code) to also confirm that each input's `PublicKey`, when hashed, actually matches the `PubKeyHash` on the output it claims to spend (currently, `Verify` only checks that the *signature* is valid for the given public key — it does not separately confirm that public key is the *right* one for that output). Explain, with a concrete scenario, what attack this additional check specifically prevents that plain signature verification alone would miss.
