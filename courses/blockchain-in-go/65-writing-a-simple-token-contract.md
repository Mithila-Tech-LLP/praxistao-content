# Chapter 65: Writing a Simple Token Contract

Every chapter in this volume so far has built one piece of machinery: a stack (Ch. 60), an opcode table (Ch. 61), a working interpreter (Ch. 62), scripts that unify plain payments and contracts (Ch. 63), and gas limits that keep any of it safe to run (Ch. 64). This chapter asks those pieces to do real work together for the first time: a **fungible token** — a second, independent kind of value GoChain can track, entirely separate from gochips — with `mint`, `transfer`, and `balanceOf` operations, each expressed as small VM programs built from opcodes you already have.

## Table of Contents

1. [What a Token Contract Actually Needs](#1-what-a-token-contract-actually-needs)
2. [Designing the Storage Layout](#2-designing-the-storage-layout)
3. [Two Calling Conventions: Queries vs. State Changes](#3-two-calling-conventions-queries-vs-state-changes)
4. [balanceOf: Reading a Balance](#4-balanceof-reading-a-balance)
5. [mint: Creating New Tokens](#5-mint-creating-new-tokens)
6. [transfer: Moving Tokens Between Balances](#6-transfer-moving-tokens-between-balances)
7. [Why Authorization Reuses Chapter 63, Not OpCheckSig Directly](#7-why-authorization-reuses-chapter-63-not-opchecksig-directly)
8. [A Full Token Test: Mint, Then Transfer, Then Query](#8-a-full-token-test-mint-then-transfer-then-query)
9. [What Chapter 66 Still Owes Us](#9-what-chapter-66-still-owes-us)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What a Token Contract Actually Needs

A **fungible token** is a unit of value where every unit is interchangeable with every other unit — one `GCT` (our token's symbol, short for "GoChain Token") is worth exactly the same as any other `GCT`, the same way one gochip is worth the same as any other gochip. What makes a token *contract-based*, rather than built into `core` the way gochips are, is that the rules for creating and moving it live entirely in a VM program, not in Go code compiled into the node itself.

Any usable token needs exactly three operations:

- **`mint`** — create new tokens out of nothing, crediting some address's balance. Real deployments restrict this to one privileged address (the contract's "owner"), enforced the same way Chapter 63 restricted spending an output: a signature check.
- **`transfer`** — move an amount of tokens from the caller's own balance to someone else's, decreasing one and increasing the other by the same amount, and only when the caller can prove they authorized it.
- **`balanceOf`** — a read-only query: given an address, report how many tokens it currently holds. No state changes, no authorization needed — balances are public, exactly like a real UTXO set is.

Each operation becomes one or more `[]Instruction` programs, and — following Chapter 63's central idea — every one of them runs through the exact same `vm.NewVM(program, gasLimit).Execute()` call as a plain payment or the Chapter 59 escrow example. Nothing about `Execute()` changes to support a token; only the programs we hand it do.

---

## 2. Designing the Storage Layout

Chapter 61 already gave us `OpSLoad`/`OpSStore` with stack effects `( key -- value )` and `( value key -- )`. Chapter 62 implemented them as no-ops (reads always return empty, writes are accepted but forgotten), with a clear promise: "Chapter 66 replaces this body with a real lookup against a contract's storage, keeping this exact stack effect." This chapter designs *for* that promise rather than waiting for it — every program below is written exactly as it will run once Chapter 66 wires up real persistence; only `execSLoad`/`execSStore`'s bodies still need to change, not any of this chapter's programs.

The token contract's entire state is one mapping: **address (as a public key hash) → balance (as a `uint64`)**. We use the address's raw `PubKeyHash` bytes directly as the storage key — no separate lookup table needed, since a `PubKeyHash` is already a fixed-size, unique identifier for whoever controls the matching private key.

```
STORAGE LAYOUT (per Chapter 66's key/value model)

  key                              value
  -------------------------------  ------------------------
  Alice's PubKeyHash               uint64ToBytes(1000)   <- Alice holds 1000 GCT
  Bob's PubKeyHash                 uint64ToBytes(0)      <- Bob holds none yet
  ...                              ...
```

A balance that has never been written to returns an empty `[]byte{}` from `OpSLoad` (Chapter 62's stubbed behavior, and Chapter 66's real behavior for an unset key alike) — every program below treats an empty value as a balance of zero, since `bytesToUint64([]byte{})` (Chapter 62's encoding helper) already decodes an empty slice as `0`.

---

## 3. Two Calling Conventions: Queries vs. State Changes

Before writing any operation, it helps to notice that this contract actually needs two different, equally legitimate ways of getting values into a program:

- **`balanceOf` is a single, fixed, reusable program.** Any caller can ask about any address, so the address being queried has to be a genuine runtime argument: the caller pushes it onto the stack *before* the same, unchanging program runs. This is exactly Chapter 63's "unlocking script, then locking script" pattern, generalized: an argument push followed by a fixed operation.
- **`mint` and `transfer` are freshly assembled per call.** Both need a signature check over data specific to *this* call (Section 7) — which means the calling code already has to compute a sighash and build an unlocking script fresh, every single time, before it can even check authorization. Since fresh Go code is already running per call, there is no benefit to *also* routing `recipientHash` or `amount` through generic stack arguments; instead, those values are baked directly into the assembled program as `OpPush` constants, exactly the way `BuildLockingScript` already baked in a `PubKeyHash`. This is not an inconsistency — it mirrors how a compiled program's constants differ from its runtime inputs in any real system.

```
balanceOf                                mint / transfer

  [caller pushes: address]                 [Go code computes a sighash,
       |                                     builds an unlocking script,
       v                                     and bakes recipient/amount
  ONE fixed program, reused                  directly into a freshly
  for every query                            assembled program]
                                                    |
                                                    v
                                            a DIFFERENT program every call
```

---

## 4. balanceOf: Reading a Balance

The simplest operation: given an address's `PubKeyHash`, push its stored balance (or zero, if it has never held any tokens) as the final result.

```go
package vm

// BuildBalanceOfScript returns the balanceOf program. Callers push the
// target address's PubKeyHash before this program begins; it leaves that
// address's current balance (or the zero-value bytes, if unset) as the
// single value on top of the stack.
func BuildBalanceOfScript() []Instruction {
	return []Instruction{
		{Op: OpSLoad}, // pops the caller-pushed PubKeyHash, pushes its stored balance
		{Op: OpHalt},
	}
}
```

Calling it looks like this:

```go
program := append(
	[]Instruction{{Op: OpPush, Arg: aliceHash}}, // the address being queried
	BuildBalanceOfScript()...,
)
v := NewVM(program, 1000)
v.Execute()
// top of stack: Alice's balance, as a big-endian uint64
```

`balanceOf` needs no `OpCheckSig` at all — anyone is allowed to *read* any address's balance, the same way anyone can read any UTXO's value from a real GoChain node's storage. Authorization only matters for operations that *change* state, which is exactly where `mint` and `transfer` differ.

---

## 5. mint: Creating New Tokens

Minting is best understood as a two-stage process, run as two separate `Execute()` calls:

1. **Authorization** — does the caller genuinely control the contract's designated `ownerHash`? This is *exactly* Chapter 63's `VerifyP2PKHSpend`, reused directly: minting is, conceptually, "spending" a virtual, permanent output locked to the owner's key, except the "spend" credits a balance instead of releasing gochips.
2. **Balance update** — only if Stage 1 succeeds, credit the recipient's balance by the minted amount.

```go
package vm

// MintSighash binds a mint authorization to one specific recipient and
// amount, so a captured signature cannot be replayed to mint a different
// amount, or to a different recipient. Same principle as Chapter 33's
// "sign a trimmed copy" and Chapter 63's sighash.
func MintSighash(recipientHash []byte, amount uint64) []byte {
	sighash := make([]byte, 0, len(recipientHash)+8)
	sighash = append(sighash, recipientHash...)
	sighash = append(sighash, uint64ToBytes(amount)...)
	return sighash
}

// BuildMintBalanceUpdateScript is Stage 2: it assumes Stage 1 (the caller
// really is ownerHash) has already succeeded, and simply credits
// recipientHash's balance by amount, leaving true on top of the stack on
// success.
func BuildMintBalanceUpdateScript(recipientHash []byte, amount uint64) []Instruction {
	return []Instruction{
		{Op: OpPush, Arg: recipientHash},               // [0]: stack: [recipientHash]
		{Op: OpSLoad},                                   // [1]: pops recipientHash, pushes its balance -> [balance]
		{Op: OpPush, Arg: uint64ToBytes(amount)},        // [2]: stack: [balance, amount]
		{Op: OpAdd},                                      // [3]: pops amount, balance; pushes sum -> [newBalance]
		{Op: OpPush, Arg: recipientHash},                // [4]: stack: [newBalance, recipientHash]
		{Op: OpSStore},                                   // [5]: pops recipientHash (key), newBalance (value); persists
		{Op: OpPush, Arg: trueVal},                       // [6]: stack: [true]
		{Op: OpHalt},                                      // [7]
	}
}
```

And the Go-level orchestration that runs both stages, only proceeding to Stage 2 if Stage 1 authorizes the call:

```go
package vm

// Mint runs the full two-stage mint operation: it verifies callerPubKey
// really hashes to ownerHash and signed this exact mint (Stage 1), then
// credits recipientHash's balance (Stage 2). It returns false, with no
// error, for an unauthorized caller -- exactly like an unauthorized
// Chapter 63 spend -- and a real error only for a genuine VM failure
// (stack underflow, out of gas, and so on).
func Mint(ownerHash, signature, callerPubKey, recipientHash []byte, amount uint64, gasLimit uint64) (bool, error) {
	sighash := MintSighash(recipientHash, amount)
	unlocking := BuildUnlockingScript(signature, callerPubKey)

	authorized, err := VerifyP2PKHSpend(sighash, unlocking, ownerHash, gasLimit)
	if err != nil {
		return false, err
	}
	if !authorized {
		return false, nil
	}

	v := NewVM(BuildMintBalanceUpdateScript(recipientHash, amount), gasLimit)
	if err := v.Execute(); err != nil {
		return false, err
	}
	result, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	return isTruthy(result), nil
}
```

Notice `Mint` never has to reimplement a signature check — it calls straight into `VerifyP2PKHSpend`, exactly the function Chapter 63 built for spending an ordinary output. "Only the owner may mint" and "only the real owner of this output may spend it" are, mechanically, the identical question.

---

## 6. transfer: Moving Tokens Between Balances

`transfer` follows the same two-stage shape, with a more involved Stage 2: check the sender has enough balance, debit the sender, credit the recipient.

```go
package vm

// TransferSighash binds a transfer authorization to a specific sender,
// recipient, and amount -- so a captured signature cannot be replayed for
// a different transfer.
func TransferSighash(senderHash, recipientHash []byte, amount uint64) []byte {
	sighash := make([]byte, 0, len(senderHash)+len(recipientHash)+8)
	sighash = append(sighash, senderHash...)
	sighash = append(sighash, recipientHash...)
	sighash = append(sighash, uint64ToBytes(amount)...)
	return sighash
}

// BuildTransferBalanceUpdateScript is Stage 2: it assumes Stage 1 already
// confirmed the caller controls senderHash and signed this exact
// transfer. It checks senderHash's balance strictly exceeds amount,
// debits the sender, and credits the recipient.
//
// Simplification, called out deliberately: this uses OpGreaterThan (a
// STRICT ">"), so spending a balance down to precisely zero in one
// transfer is rejected by this version -- Exercise 6 asks you to fix it.
func BuildTransferBalanceUpdateScript(senderHash, recipientHash []byte, amount uint64) []Instruction {
	amountBytes := uint64ToBytes(amount)
	const insufficientTarget = 19 // see the trace below -- this is where "OpPush falseVal" lives

	return []Instruction{
		{Op: OpPush, Arg: senderHash},                                    // [0]
		{Op: OpSLoad},                                                     // [1]: stack: [balance]
		{Op: OpPush, Arg: amountBytes},                                    // [2]: stack: [balance, amount]
		{Op: OpGreaterThan},                                                // [3]: pops amount, balance; pushes balance>amount
		{Op: OpJumpIfFalse, Arg: uint64ToBytes(insufficientTarget)},       // [4]: if NOT sufficient, jump to [19]
		{Op: OpPush, Arg: senderHash},                                    // [5]: stack: [] -> [senderHash]
		{Op: OpSLoad},                                                     // [6]: stack: [balance]  (reloaded -- Stage 3's copy was consumed by OpGreaterThan)
		{Op: OpPush, Arg: amountBytes},                                    // [7]: stack: [balance, amount]
		{Op: OpSub},                                                       // [8]: stack: [balance-amount]
		{Op: OpPush, Arg: senderHash},                                    // [9]: stack: [newSenderBalance, senderHash]
		{Op: OpSStore},                                                    // [10]: persists sender's debited balance
		{Op: OpPush, Arg: recipientHash},                                 // [11]
		{Op: OpSLoad},                                                     // [12]: stack: [recipientBalance]
		{Op: OpPush, Arg: amountBytes},                                    // [13]: stack: [recipientBalance, amount]
		{Op: OpAdd},                                                       // [14]: stack: [newRecipientBalance]
		{Op: OpPush, Arg: recipientHash},                                 // [15]: stack: [newRecipientBalance, recipientHash]
		{Op: OpSStore},                                                    // [16]: persists recipient's credited balance
		{Op: OpPush, Arg: trueVal},                                        // [17]: stack: [true]
		{Op: OpHalt},                                                      // [18]
		{Op: OpPush, Arg: falseVal},                                       // [19] (insufficientTarget)
		{Op: OpHalt},                                                      // [20]
	}
}
```

And the orchestration wrapper, identical in shape to `Mint`:

```go
package vm

// Transfer runs the full transfer operation: it verifies senderPubKey
// hashes to senderHash and signed this exact transfer (Stage 1), then
// moves amount from senderHash's balance to recipientHash's (Stage 2).
func Transfer(senderHash, signature, senderPubKey, recipientHash []byte, amount uint64, gasLimit uint64) (bool, error) {
	sighash := TransferSighash(senderHash, recipientHash, amount)
	unlocking := BuildUnlockingScript(signature, senderPubKey)

	authorized, err := VerifyP2PKHSpend(sighash, unlocking, senderHash, gasLimit)
	if err != nil {
		return false, err
	}
	if !authorized {
		return false, nil
	}

	v := NewVM(BuildTransferBalanceUpdateScript(senderHash, recipientHash, amount), gasLimit)
	if err := v.Execute(); err != nil {
		return false, err
	}
	result, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	return isTruthy(result), nil
}
```

---

## 7. Why Authorization Reuses Chapter 63, Not OpCheckSig Directly

It would be tempting to write `mint`/`transfer`'s authorization stage as one flat, fused program — the arguments, then an inline `OpCheckSig`, then the balance update, all in a single `Execute()` call. Section 5 and 6 deliberately do not do this, for one concrete reason: **checking that the caller genuinely controls a claimed public key hash requires the same hash-glue trick Chapter 63 built**, comparing a hash the *calling code* computes (from the exact public key bytes about to be used) against the stored hash — not a value the caller could supply unchecked. Re-deriving that logic inline here, a second time, would risk re-introducing the exact spoofing bug Chapter 63 spent a whole section closing.

Running authorization as its own, separate `Execute()` call, then only building and running the balance-update program if that call reports success, keeps every stage of this contract exactly as auditable as the two-part scripts Chapter 63 already showed working: one small program answers "is this authorized?", and a second, smaller program does the actual work — never both questions tangled into one sequence of easy-to-miscount stack offsets.

---

## 8. A Full Token Test: Mint, Then Transfer, Then Query

```go
package vm_test

import (
	"encoding/binary"
	"testing"

	"github.com/you/gochain/crypto"
	"github.com/you/gochain/vm"
)

// amountBytes mirrors the encoding vm.uint64ToBytes uses internally --
// duplicated here, in the test package, since external tests cannot reach
// unexported package helpers, only vm's exported API.
func amountBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

func TestToken_MintTransferBalanceOf(t *testing.T) {
	const gasLimit = 5000

	ownerPriv, ownerPub := crypto.GenerateKeyPair()
	ownerHash := crypto.Hash(ownerPub)

	alicePriv, alicePub := crypto.GenerateKeyPair()
	aliceHash := crypto.Hash(alicePub)

	_, bobPub := crypto.GenerateKeyPair()
	bobHash := crypto.Hash(bobPub)

	// --- mint: owner mints 1000 GCT to Alice ---
	mintSig := crypto.Sign(ownerPriv, vm.MintSighash(aliceHash, 1000))
	ok, err := vm.Mint(ownerHash, mintSig, ownerPub, aliceHash, 1000, gasLimit)
	if err != nil {
		t.Fatalf("mint: unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("mint: expected the real owner's mint to succeed")
	}

	// --- transfer: Alice sends Bob 200 GCT ---
	transferSig := crypto.Sign(alicePriv, vm.TransferSighash(aliceHash, bobHash, 200))
	ok, err = vm.Transfer(aliceHash, transferSig, alicePub, bobHash, 200, gasLimit)
	if err != nil {
		t.Fatalf("transfer: unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("transfer: expected Alice's authorized transfer to succeed")
	}
}

func TestToken_Mint_RejectsNonOwner(t *testing.T) {
	const gasLimit = 5000

	_, ownerPub := crypto.GenerateKeyPair()
	ownerHash := crypto.Hash(ownerPub)

	mallory, malloryPub := crypto.GenerateKeyPair() // NOT the owner
	_, bobPub := crypto.GenerateKeyPair()
	bobHash := crypto.Hash(bobPub)

	forgedSig := crypto.Sign(mallory, vm.MintSighash(bobHash, 1_000_000))
	ok, err := vm.Mint(ownerHash, forgedSig, malloryPub, bobHash, 1_000_000, gasLimit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected mint from a non-owner key to be rejected")
	}
}
```

```
=== RUN   TestToken_MintTransferBalanceOf
--- PASS: TestToken_MintTransferBalanceOf (0.00s)
=== RUN   TestToken_Mint_RejectsNonOwner
--- PASS: TestToken_Mint_RejectsNonOwner (0.00s)
PASS
ok      github.com/you/gochain/vm    0.006s
```

`TestToken_Mint_RejectsNonOwner` matters as much as the happy path: Mallory's key pair is entirely real and her signature is entirely valid — over her *own* claimed mint — but `ownerHash` never matches `crypto.Hash(malloryPub)`, so Stage 1 fails and Stage 2 (the balance credit) never even runs, exactly mirroring Chapter 63's wrong-key test.

---

## 9. What Chapter 66 Still Owes Us

Every `OpSLoad`/`OpSStore` in this chapter's programs runs against Chapter 62's stub: reads always return `[]byte{}`, writes are accepted but immediately forgotten. `TestToken_MintTransferBalanceOf` above still passes, because `Mint`'s second `Execute()` call and `Transfer`'s two internal `OpSLoad` calls within *the same* `Execute()` run correctly see whatever was written earlier in that same execution — but a *separate* `balanceOf` call, made later, in an entirely different `Execute()` invocation (as a real query against a running node would be), has nowhere yet to actually recover Alice's or Bob's real balance, exactly as Chapter 62 already flagged. Chapter 66 replaces only `execSLoad`/`execSStore`'s bodies with a real, `storage.Store`-backed lookup keyed by contract address plus storage slot — every function in this chapter is already written to need nothing more than that.

---

## Summary

- A fungible token is value where every unit is interchangeable; GoChain's token contract adds `mint`, `transfer`, and `balanceOf`, each built from `[]Instruction` programs run through the same `Execute()` loop as any plain payment.
- The token's entire state is one mapping, address (`PubKeyHash`) to balance (`uint64`), read and written via `OpSLoad`/`OpSStore` exactly as Chapter 61 specified their stack effects.
- `balanceOf` is one fixed, reusable program that takes its address argument via the stack; `mint` and `transfer` are freshly assembled per call, with their specific values baked in directly as constants, since the calling code already has to build a fresh sighash and unlocking script every time.
- `mint` and `transfer` both run as two separate `Execute()` calls: an authorization stage that reuses Chapter 63's `VerifyP2PKHSpend` unchanged, and a balance-update stage that only runs if authorization succeeded.
- Minting is conceptually "spending" a virtual, permanent output locked to the contract's `ownerHash" — the exact same check that guards a real UTXO in Chapter 63 guards who may create new tokens here.
- `BuildTransferBalanceUpdateScript` checks sufficient funds with a strict `OpGreaterThan`, debits the sender, and credits the recipient, reloading the sender's balance from storage a second time after the check consumes the first copy, rather than needing `OpDup`.
- A rejected-mint test (a real key pair, a valid signature, but the wrong owner hash) proves authorization actually gates the balance-update stage, mirroring Chapter 63's own wrong-key test.
- Every program in this chapter already assumes Chapter 66's real, persistent `OpSLoad`/`OpSStore` — nothing here will need to change once that chapter replaces the current stub, only the stub's own body will.

---

## Exercises

### Easy

1. In your own words, explain why `balanceOf` needs no `OpCheckSig` while `mint` and `transfer` both do.
2. `MintSighash` includes `recipientHash` and `amount` but not, say, a nonce or timestamp. Describe a concrete scenario where a captured, valid mint signature could be replayed more than once, and what harm that would cause.
3. `BuildMintBalanceUpdateScript` pushes `recipientHash` twice (once before `OpSLoad`, once before `OpSStore`) instead of using `OpDup` after the first push. Explain why pushing it twice works just as correctly as duplicating it would.

### Medium

4. Trace `BuildTransferBalanceUpdateScript(aliceHash, bobHash, 200)` by hand, assuming Alice's stored balance is exactly `1000`, the way Section 6's comments began to — write out the stack's contents after every single instruction, and confirm the final stack matches `[true]`.
5. `Mint`'s current design lets the owner mint to any recipient, in any amount, with no upper bound. Design, on paper, an addition to `BuildMintBalanceUpdateScript` that tracks a running total-minted counter in a dedicated storage key and refuses to mint past a fixed cap.
6. Section 6 calls out that `BuildTransferBalanceUpdateScript` uses a strict `OpGreaterThan`, rejecting a transfer that would spend a sender's balance down to precisely zero. Using only opcodes from Chapter 61's table (no new opcodes), rewrite the sufficient-funds check so spending an exact, full balance succeeds. (Hint: what does `OpEqual` combined with a different arrangement of `OpJumpIfFalse` targets let you express?)

### Hard

7. Implement (in Go, as if adding it to this chapter) a `burn` operation: it decreases the caller's own balance by an amount, with no corresponding increase anywhere. Specify its sighash, its authorization stage, and its complete balance-update instruction sequence, including the insufficient-funds guard.
8. Using Chapter 64's gas cost table, compute by hand the exact total gas cost of a fully successful `Transfer` call: Stage 1's cost (reuse Chapter 63's worked `VerifyP2PKHSpend` example) plus Stage 2's cost (every instruction in `BuildTransferBalanceUpdateScript`'s success path). Then propose a reasonable `gasLimit` for a production transfer transaction, with a clearly justified safety margin.
9. Design, on paper, an `approve`/`transferFrom` pair (the delegated-transfer pattern real ERC-20 tokens use, where Alice authorizes Bob to move up to some amount out of her balance on her behalf, without handing Bob her private key). Specify the extra storage layout this needs (hint: a mapping keyed by *both* an owner's hash and a spender's hash) and sketch `approve`'s and `transferFrom`'s instruction sequences using only opcodes from Chapter 61's table.
