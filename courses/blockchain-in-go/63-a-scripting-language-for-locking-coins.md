# Chapter 63: A Scripting Language for Locking Coins

Chapter 32 gave every `TxOutput` a `PubKeyHash` field and a `bytes.Equal`-based `IsLockedWithKey` check. That was enough for a world where "can this output be spent?" only ever meant "does a signature match this one specific key." Chapters 59 through 62 built something more general: a virtual machine that can run *any* small, deterministic program. This chapter closes the gap between those two worlds. We redefine "is this output spendable?" as "does a small VM program, built from this output, evaluate to true?" — and show that a plain payment and a full smart contract are answered by the exact same question, run through the exact same `Execute()` call.

## Table of Contents

1. [Recap: Locking as a Byte Comparison](#1-recap-locking-as-a-byte-comparison)
2. [Locking and Unlocking Scripts, Defined](#2-locking-and-unlocking-scripts-defined)
3. [The One Snag: No Hash Opcode](#3-the-one-snag-no-hash-opcode)
4. [Building an Unlocking Script](#4-building-an-unlocking-script)
5. [Building a Locking Script](#5-building-a-locking-script)
6. [Combining Two Scripts Into One Program](#6-combining-two-scripts-into-one-program)
7. [Verifying a Full Spend, Step by Step](#7-verifying-a-full-spend-step-by-step)
8. [Wiring This Into core.TxOutput](#8-wiring-this-into-coretxoutput)
9. [Unifying Plain Transactions and Smart Contracts](#9-unifying-plain-transactions-and-smart-contracts)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Recap: Locking as a Byte Comparison

Chapter 32 defined `core.TxOutput` like this:

```go
type TxOutput struct {
	Value      int64  // amount, in gochips
	PubKeyHash []byte // identifies who can spend this output
}
```

and gave it `IsLockedWithKey`, a plain `bytes.Equal` check against a candidate public key hash. That check answers one narrow question well — "does this hash match?" — but it has no way to also ask "and did the real owner actually authorize *this* specific spend?" Chapter 33 answered that second question separately, with `Transaction.Sign()`/`Transaction.Verify()`, calling straight into `gochain/crypto`.

So, as of Chapter 33, spending an output required two entirely separate pieces of logic living in two different places: a hash comparison in `core`, and a signature check in `core` calling out to `crypto`. Neither piece knew about the other, and neither could be swapped out, extended, or replaced with something more elaborate (like an escrow condition) without rewriting `core` itself.

This chapter merges both checks into one artifact: a small VM program. Once an output's spendability is "run this program and see if it leaves `true` on the stack," `core` no longer needs to know *how* an output decides who can spend it — only that it can ask the VM to decide. That single change is what lets Chapter 65's token contract and this chapter's plain payments share the same execution path.

---

## 2. Locking and Unlocking Scripts, Defined

Real UTXO-based blockchains (Bitcoin among them) split "can this output be spent?" into two small programs that are **concatenated and run together**:

- A **locking script** (Bitcoin calls this a `scriptPubKey`) is attached to the output itself, at the moment it is created. It encodes the *rule* — "whoever provides a valid signature for this specific public key hash may spend this."
- An **unlocking script** (Bitcoin calls this a `scriptSig`) is supplied by whoever wants to spend the output, at spend time. It encodes the *proof* — "here is my signature, and here is my public key."

Neither script means anything by itself. The locking script has no signature to check yet; the unlocking script has no rule to satisfy yet. Only when the unlocking script runs *first*, leaving values on the stack, followed immediately by the locking script consuming them, does the combined program answer the question "is this spend authorized?"

```
                    UNLOCKING SCRIPT              LOCKING SCRIPT
                    (from the spender,             (baked into the
                     supplied now)                  output, long ago)
                         |                                |
                         v                                v
                +------------------+           +----------------------+
                | OpPush signature |           | ... checks that use   |
                | OpPush pubkey    |  ------>  | whatever the unlocking |
                +------------------+           | script left on the    |
                                                | stack ...              |
                                                +----------------------+
                         |________________________________|
                                        |
                                        v
                          ONE combined program, run through
                          the exact same vm.VM.Execute() loop
                          Chapter 62 already built.
```

This is the whole idea of the chapter: locking and unlocking scripts are not two different kinds of thing running two different kinds of machine. They are two *fragments* of one ordinary VM program, assembled at spend time, executed by the identical `Execute()` loop that runs a Chapter 65 token contract.

---

## 3. The One Snag: No Hash Opcode

A standard **pay-to-public-key-hash** (P2PKH) locking script needs to check two things about whatever the unlocking script provides: (1) does the provided public key actually hash to the output's stored `PubKeyHash`, and (2) is the provided signature valid for that public key over the data being spent. Real Bitcoin Script expresses check (1) *inside* the script itself, using an opcode (`OP_HASH160`) that hashes whatever is currently on top of the stack.

Chapter 61's opcode table, deliberately kept small, has no such opcode — only `OpCheckSig` reaches into cryptography, and it only verifies signatures, never computes a hash. This is not an oversight to patch quietly; it is worth being explicit about, because the fix matters for security:

- **What will not work:** having the spender push their own claimed "hash of my public key" onto the stack, then comparing it against `PubKeyHash` with `OpEqual`. Nothing stops a dishonest spender from pushing an unrelated, fabricated value that happens to equal the target `PubKeyHash`, while pushing a *different* public key (one they control the private key for) for the actual signature check. The comparison would pass despite the spender never having proven anything about the real owner's key.
- **What does work, and is what this chapter builds:** the code that *combines* the unlocking script and the locking script — not the VM itself — computes `crypto.Hash` of the public key the unlocking script actually provides, using the exact same bytes that will be fed into the VM a moment later. Because this computation reads directly from the instruction it is about to include, there is no way for a spender to make the VM see one public key while a different one gets hash-checked. The hashing happens in ordinary, already-trusted Go code — the same `crypto.Hash` function Chapter 09 built and Chapter 14 already uses for addresses — not inside the restricted VM sandbox, because computing a hash needs no secret information and every honest node computes the identical result regardless of where the code runs. What *does* need the VM's sandboxed, gas-metered machinery is the signature check, because that is the one operation Chapter 61 gave a dedicated, careful opcode for: `OpCheckSig`.

We will revisit this exact limitation as an exercise at the end of the chapter — it is a legitimate design question, not a shortcut this chapter is hiding.

---

## 4. Building an Unlocking Script

An unlocking script for a P2PKH output is the simplest possible fragment: push the signature, then push the public key, in that order, so that by the time the locking script runs, the public key sits nearest the top of the stack — matching `OpCheckSig`'s expected pop order from Chapter 61 (`pubkey` popped first).

```go
package vm

// BuildUnlockingScript returns the instructions a spender contributes when
// satisfying a pay-to-public-key-hash locking script: their signature,
// then their public key. By convention, the last instruction in an
// unlocking script always pushes the public key -- Section 6 relies on
// this to find it later.
func BuildUnlockingScript(signature, pubKey []byte) []Instruction {
	return []Instruction{
		{Op: OpPush, Arg: signature},
		{Op: OpPush, Arg: pubKey},
	}
}
```

Nothing here is specific to any one output — the same two-instruction shape works for spending *any* P2PKH output, because all the output-specific logic (Section 5) lives entirely on the locking side.

---

## 5. Building a Locking Script

A locking script is created once, when the output itself is created, and it is parameterized only by that output's `PubKeyHash` — nobody yet knows which specific public key will eventually spend it. The script's job: compare a hash (which Section 6's combining step will have already pushed) against the output's own `PubKeyHash`, bail out on a mismatch, and otherwise fall through to `OpCheckSig`.

```go
package vm

// falseVal and uint64ToBytes are already defined in vm.go (Chapter 62);
// reused here so a failed locking script leaves an explicit, unambiguous
// false result, and the jump target is encoded the same way every other
// opcode's Arg already is.

// BuildLockingScript returns the standard pay-to-public-key-hash program
// for an output locked to pubKeyHash. startIndex is the absolute index
// this script will occupy once combined with a sighash push and an
// unlocking script (Section 6) -- needed so its OpJumpIfFalse can point at
// the correct fail-path instruction within the final, assembled program.
func BuildLockingScript(pubKeyHash []byte, startIndex int) []Instruction {
	failTarget := uint64(startIndex + 5)
	return []Instruction{
		{Op: OpPush, Arg: pubKeyHash},                       // startIndex + 0
		{Op: OpEqual},                                        // startIndex + 1
		{Op: OpJumpIfFalse, Arg: uint64ToBytes(failTarget)},  // startIndex + 2
		{Op: OpCheckSig},                                     // startIndex + 3
		{Op: OpHalt},                                         // startIndex + 4
		{Op: OpPush, Arg: falseVal},                          // startIndex + 5 (fail target)
		{Op: OpHalt},                                         // startIndex + 6
	}
}
```

Walking through it: instruction 0 pushes the constant this output was locked with; instruction 1 compares it against whatever sits on top of the stack (which, by the time this script runs, will be the hash of the spender's provided public key); instruction 2 bails out to the fail path the instant that comparison is false; instruction 3, reached only if the hashes matched, hands off to `OpCheckSig`, which does the actual cryptographic work; instruction 4 halts cleanly on success. Instructions 5-6 are the fail path: push an explicit `false`, then halt — so a caller inspecting the top of the stack after `Execute()` always finds an unambiguous boolean, never garbage.

---

## 6. Combining Two Scripts Into One Program

This is the piece that actually assembles a runnable VM program out of three ingredients: the data being authorized (the **sighash** — the exact bytes Chapter 33's `Sign`/`Verify` operate over), the spender's unlocking script, and the output's locking script.

```go
package vm

import (
	"errors"

	"github.com/you/gochain/crypto"
)

// ErrEmptyUnlockingScript means a spender supplied no instructions at all
// -- there is no public key to check anything against.
var ErrEmptyUnlockingScript = errors.New("vm: unlocking script must push at least a public key")

// VerifyP2PKHSpend combines sighash, unlockingScript, and a locking script
// built from pubKeyHash into one VM program and reports whether the
// combined program authorizes the spend. This is the single function
// every P2PKH spend check in GoChain funnels through.
func VerifyP2PKHSpend(sighash []byte, unlockingScript []Instruction, pubKeyHash []byte, gasLimit uint64) (bool, error) {
	if len(unlockingScript) == 0 {
		return false, ErrEmptyUnlockingScript
	}

	// By convention (Section 4), the LAST instruction an unlocking script
	// pushes is the spender's public key. We read its raw bytes directly
	// -- the same bytes about to enter the VM -- so the hash we compute
	// next cannot be spoofed by the spender pushing a mismatched value.
	providedPubKey := unlockingScript[len(unlockingScript)-1].Arg
	pubKeyHashGlue := crypto.Hash(providedPubKey)

	program := make([]Instruction, 0, 2+len(unlockingScript)+7)
	program = append(program, Instruction{Op: OpPush, Arg: sighash})        // [0] the data being signed
	program = append(program, unlockingScript...)                          // spender's signature, pubkey
	program = append(program, Instruction{Op: OpPush, Arg: pubKeyHashGlue}) // engine-computed, not spender-supplied

	lockingStart := len(program)
	program = append(program, BuildLockingScript(pubKeyHash, lockingStart)...)

	v := NewVM(program, gasLimit)
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

Notice exactly one instruction in this whole assembly is *neither* part of the unlocking script *nor* part of the locking script: the `OpPush pubKeyHashGlue` line. It is inserted by the combining code itself, computed from the exact public-key bytes that immediately precede it in the very same program — which is precisely the guarantee Section 3 said was needed.

---

## 7. Verifying a Full Spend, Step by Step

Trace a concrete spend through the combined program, the same way Chapter 60 traced arithmetic. Say Alice's output was locked to `pkh` (her public key hash), and she now wants to spend it by proving she holds the matching private key. Symbols stand in for real byte slices: `sighash` (the message), `sig` (Alice's real signature over it), `pub` (Alice's real public key), `H(pub)` (its hash, computed by the combining code), and `pkh` (the value baked into the output at creation time — in a valid spend, `H(pub) == pkh`).

```
[0] OpPush sighash          stack: [ sighash ]
[1] OpPush sig              stack: [ sighash, sig ]
[2] OpPush pub              stack: [ sighash, sig, pub ]
[3] OpPush H(pub)           stack: [ sighash, sig, pub, H(pub) ]
[4] OpPush pkh              stack: [ sighash, sig, pub, H(pub), pkh ]
[5] OpEqual                 pops pkh, H(pub); pushes true (they match)
                             stack: [ sighash, sig, pub, true ]
[6] OpJumpIfFalse -> [9]    pops true; does NOT jump; falls through
                             stack: [ sighash, sig, pub ]
[7] OpCheckSig              pops pub, sig, sighash; verifies; pushes true
                             stack: [ true ]
[8] OpHalt                  EXECUTION STOPS -- top of stack: true
```

And the rejected case, where Mallory tries to spend Alice's output using her own key pair (`pub2`, `sig2`), so `H(pub2) != pkh`:

```
[0] OpPush sighash          stack: [ sighash ]
[1] OpPush sig2             stack: [ sighash, sig2 ]
[2] OpPush pub2             stack: [ sighash, sig2, pub2 ]
[3] OpPush H(pub2)          stack: [ sighash, sig2, pub2, H(pub2) ]
[4] OpPush pkh              stack: [ sighash, sig2, pub2, H(pub2), pkh ]
[5] OpEqual                 pops pkh, H(pub2); they DIFFER; pushes false
                             stack: [ sighash, sig2, pub2, false ]
[6] OpJumpIfFalse -> [9]    pops false; JUMPS to instruction 9
[9] OpPush falseVal         stack: [ sighash, sig2, pub2, false ]
[10] OpHalt                 EXECUTION STOPS -- top of stack: false
```

Mallory's own signature is never even checked — the program never reaches `OpCheckSig` at all, because the hash comparison fails first. This mirrors exactly how Bitcoin Script short-circuits: cheaper checks run first, and a failure early on skips the more expensive cryptographic work entirely (a detail Chapter 64's gas accounting makes concrete: Mallory's failed spend costs far less gas than a real signature check would).

Now let's turn this trace into real, runnable tests:

```go
package vm_test

import (
	"testing"

	"github.com/you/gochain/crypto"
	"github.com/you/gochain/vm"
)

func TestVerifyP2PKHSpend_Valid(t *testing.T) {
	alicePriv, alicePub := crypto.GenerateKeyPair()
	aliceHash := crypto.Hash(alicePub) // what Alice's output is locked to

	sighash := []byte("pay Bob 5 gochips")
	sig := crypto.Sign(alicePriv, sighash)

	unlocking := vm.BuildUnlockingScript(sig, alicePub)
	ok, err := vm.VerifyP2PKHSpend(sighash, unlocking, aliceHash, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected a valid, correctly-keyed spend to verify")
	}
}

func TestVerifyP2PKHSpend_WrongKey(t *testing.T) {
	alicePub := mustPubKey(t)
	aliceHash := crypto.Hash(alicePub) // Alice's output

	malloryPriv, malloryPub := crypto.GenerateKeyPair() // a different key entirely
	sighash := []byte("pay Bob 5 gochips")
	sig := crypto.Sign(malloryPriv, sighash) // a perfectly valid signature -- just the wrong key

	unlocking := vm.BuildUnlockingScript(sig, malloryPub)
	ok, err := vm.VerifyP2PKHSpend(sighash, unlocking, aliceHash, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected mismatched public key hash to fail, regardless of a valid signature")
	}
}

func TestVerifyP2PKHSpend_TamperedSignature(t *testing.T) {
	alicePriv, alicePub := crypto.GenerateKeyPair()
	aliceHash := crypto.Hash(alicePub)

	sighash := []byte("pay Bob 5 gochips")
	sig := crypto.Sign(alicePriv, []byte("pay Bob 500 gochips")) // signed the WRONG message

	unlocking := vm.BuildUnlockingScript(sig, alicePub)
	ok, err := vm.VerifyP2PKHSpend(sighash, unlocking, aliceHash, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatalf("expected a signature over the wrong message to fail")
	}
}

func mustPubKey(t *testing.T) []byte {
	t.Helper()
	_, pub := crypto.GenerateKeyPair()
	return pub
}
```

```
=== RUN   TestVerifyP2PKHSpend_Valid
--- PASS: TestVerifyP2PKHSpend_Valid (0.00s)
=== RUN   TestVerifyP2PKHSpend_WrongKey
--- PASS: TestVerifyP2PKHSpend_WrongKey (0.00s)
=== RUN   TestVerifyP2PKHSpend_TamperedSignature
--- PASS: TestVerifyP2PKHSpend_TamperedSignature (0.00s)
PASS
ok      github.com/you/gochain/vm    0.005s
```

The wrong-key test proves the hash check alone stops an unrelated key pair, even with an otherwise-valid signature. The tampered-signature test proves the hash check passing is not sufficient by itself — `OpCheckSig` still has the final word.

---

## 8. Wiring This Into core.TxOutput

`gochain/vm` deliberately never imports `gochain/core` — Chapter 62 already established that the VM package only depends on `crypto` (and, from Chapter 66, `storage`), precisely so that `core` can safely depend on `vm` without creating an import cycle. The wiring therefore lives on the `core` side:

```go
package core

import "github.com/you/gochain/vm"

// VerifySpend checks whether signature and pubKey authorize spending out,
// given sighash -- the exact bytes Chapter 33's Sign/Verify already
// operate over. This single function is where a plain GoChain payment and
// (from Chapter 65 onward) a smart contract call both funnel through the
// same vm.Execute() call.
func VerifySpend(out *TxOutput, sighash, signature, pubKey []byte, gasLimit uint64) (bool, error) {
	unlocking := vm.BuildUnlockingScript(signature, pubKey)
	return vm.VerifyP2PKHSpend(sighash, unlocking, out.PubKeyHash, gasLimit)
}
```

`core.VerifySpend` is a thin wrapper on purpose. It takes exactly the pieces `core` already has lying around — a `*TxOutput` from Chapter 32, a signature and public key from a `TxInput` (Chapter 33's `Sign()` already produces both) — and hands them straight to the VM, without `core` needing to know a single detail of *how* the check is implemented. Chapter 33's older, ECDSA-only `Transaction.Verify()` can now be reimplemented in terms of `VerifySpend`, without changing what it returns to its own callers.

---

## 9. Unifying Plain Transactions and Smart Contracts

Recall Chapter 59's escrow example: release funds to Devan if Priya confirms delivery, or refund Priya if a deadline passes. That contract's locking script would be a *longer* program than `BuildLockingScript` — using `OpJump`, `OpGreaterThan` against the current block height, and (from Chapter 66) `OpSLoad`/`OpSStore` to remember whether a confirmation already arrived — but it is still just a `[]Instruction` slice, still combined with an unlocking script, and still run through the identical `NewVM(program, gasLimit).Execute()` call this chapter used for a plain payment.

```
PLAIN PAYMENT                              ESCROW CONTRACT (Ch. 59, sketch)

BuildLockingScript(pubKeyHash, i)          buildEscrowScript(deadline, seller, i)
  OpPush pubKeyHash                          OpSLoad  <confirmed?>
  OpEqual                                     OpJumpIfFalse -> check deadline
  OpJumpIfFalse -> fail                       ... release to seller ...
  OpCheckSig                                  ... or refund to buyer ...
  OpHalt                                      OpHalt

        \                                          /
         \________________________________________/
                            |
                            v
              vm.NewVM(program, gasLimit).Execute()
              -- the SAME loop, either way --
```

This is what Chapter 59 meant by "unifying transactions and contracts under one execution model": `core` never has a special case for "this output is a plain payment" versus "this output is a smart contract." Every output is spendable exactly when its locking script, combined with a spender's unlocking script, evaluates to `true` under the VM — a plain payment is simply the smallest possible contract, one whose only real logic is a single `OpCheckSig`.

---

## Summary

- Chapter 32's `IsLockedWithKey` and Chapter 33's signature check were two separate mechanisms; this chapter merges them into one VM program every output's spendability is decided by.
- A **locking script** is attached to an output at creation time and encodes a rule; an **unlocking script** is supplied by the spender at spend time and encodes proof; concatenated, they form one runnable program.
- Chapter 61's opcode table has no hash opcode, so the public-key-hash check cannot be baked entirely into a static locking script the way Bitcoin's `OP_HASH160` does; instead, the combining code (`VerifyP2PKHSpend`) computes the hash itself, directly from the exact bytes about to enter the VM, closing the spoofing gap a naive "spender pushes their own hash" design would leave open.
- `BuildUnlockingScript` pushes a signature then a public key; `BuildLockingScript` compares a hash against the output's `PubKeyHash`, short-circuits to an explicit `false` on mismatch, and otherwise falls through to `OpCheckSig`.
- `core.VerifySpend` wires this into `core.TxOutput` without `gochain/vm` ever importing `gochain/core`, preserving the one-directional dependency Chapter 62 established.
- A wrong public key is rejected before any signature is even checked; a right public key with a tampered signature is rejected by `OpCheckSig` itself — both failure paths were tested directly.
- The chapter's real payoff is conceptual: a plain payment and Chapter 59's escrow contract are both just `[]Instruction` slices run through the identical `Execute()` loop — "smart contract" and "plain transaction" describe *what a program does*, not *which machine runs it*.

---

## Exercises

### Easy

1. In your own words, explain why an unlocking script by itself and a locking script by itself are each meaningless — what has to happen for either one to mean anything?
2. Trace `BuildLockingScript(pkh, 4)` by hand and write out its 7 instructions with their absolute indices, the way Section 5 did, without looking back.
3. Why does `VerifyP2PKHSpend` compute `crypto.Hash(providedPubKey)` itself, instead of letting the unlocking script push a pre-computed hash value directly?

### Medium

4. `TestVerifyP2PKHSpend_WrongKey` uses a signature that is completely valid (just for the wrong key). Modify the test so instead Mallory reuses Alice's real public key but signs with her own private key, and confirm the spend still correctly fails at `OpCheckSig` rather than at the hash check. Explain, in one sentence, which stage rejects it and why.
5. `BuildLockingScript`'s fail path leaves `[sighash, sig, pub, false]` on the stack rather than just `[false]`. Rewrite `VerifyP2PKHSpend` (or explain in words why you would rather not) so the fail path also cleans up the leftover values below the top, and discuss whether this cleanup actually matters given how `VerifyP2PKHSpend` only ever inspects the top of the stack.
6. `VerifyP2PKHSpend` currently returns `(false, nil)` for a wrong-key spend but `(false, err)` for a stack underflow or gas exhaustion. Explain why callers of `core.VerifySpend` need to treat these two cases differently, and what a caller that ignored the distinction (treating any non-nil-error-or-false result the same) might get wrong.

### Hard

7. Design, on paper, a new opcode `OpHash` with stack effect `( a -- hash(a) )` that would let a locking script compute a public key's hash itself, entirely inside the VM, removing the need for `VerifyP2PKHSpend`'s glue-instruction trick. Rewrite `BuildLockingScript` (in Go, as if adding it to this chapter) to use it, and explain what changes about who is responsible for security once the hash check moves fully inside the sandboxed program.
8. `BuildLockingScript` hardcodes a very specific 7-instruction shape. Suppose a future chapter wants an output spendable by *either* of two different public key hashes (a simple 2-of-2 "OR" condition). Sketch the instruction sequence (using only opcodes from Chapter 61) that would express this, being careful with your `OpJumpIfFalse`/`OpJump` targets.
9. `VerifyP2PKHSpend` calls `v.stack.Pop()` directly on the unexported `stack` field, which only works because this function lives inside `package vm` itself. Suppose a hypothetical `gochain/contracts` package (outside `vm`) wanted to run its own combined program and inspect the final stack the same way. What would you need to add to `vm.VM`'s exported API to support that, without exposing the `Stack` type's internals more broadly than necessary?
