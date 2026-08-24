# Chapter 67: Reentrancy and Contract Security Pitfalls

Chapter 66 gave contracts real, persistent storage — a genuine step forward, and also a genuine new way to get hurt. The instant a contract can call out to another contract mid-execution, a new question appears that plain transactions never had to answer: what if that other contract calls *back*, before the first contract has finished updating its own books? This chapter reproduces, inside GoChain's own VM, a simplified version of the exact bug pattern behind the most expensive smart-contract mistake in blockchain history — the 2016 DAO hack, which drained roughly $60 million worth of Ether — and then fixes it with a single, disciplined rule: update your own state *before* you ever hand control to someone else's code.

## Table of Contents

1. [Recap: What Chapter 66 Gave Contracts](#1-recap-what-chapter-66-gave-contracts)
2. [The DAO Hack, in One Sentence](#2-the-dao-hack-in-one-sentence)
3. [Checks, Effects, Interactions — the Rule](#3-checks-effects-interactions--the-rule)
4. [A Deliberately Vulnerable Bank Contract](#4-a-deliberately-vulnerable-bank-contract)
5. [Exploiting It: The Attacker Contract](#5-exploiting-it-the-attacker-contract)
6. [Watching the Drain Happen](#6-watching-the-drain-happen)
7. [The Fix: Reordering Effects Before Interactions](#7-the-fix-reordering-effects-before-interactions)
8. [Confirming the Fix Under Attack](#8-confirming-the-fix-under-attack)
9. [Beyond Reentrancy: A Short Survey of Other Pitfalls](#9-beyond-reentrancy-a-short-survey-of-other-pitfalls)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Recap: What Chapter 66 Gave Contracts

`ContractStore` and `OpSLoad`/`OpSStore` (Chapter 66) mean a contract's storage genuinely persists between calls, isolated per contract address. That is necessary for anything useful — a balance table, a vote tally — but it also creates a gap in time: between the moment a contract *reads* its stored balance and the moment it *writes* an updated one, other code can run. If that other code is a call out to a second contract, and that second contract calls straight back into the first one before the write happens, the first contract is executing again with stale data still sitting in storage — data that has not yet reflected the withdrawal already in progress.

GoChain's VM does not (yet) have a dedicated "call another contract" opcode — that is a larger addition than this volume covers. But the *bug pattern* does not actually require one. It requires only that "send funds to whoever asked" involves handing control to code the caller supplies, before the sender's own balance is updated — and this chapter reproduces exactly that, using two small Go types built directly on `ContractStore`, standing in for two contracts that call each other.

---

## 2. The DAO Hack, in One Sentence

In June 2016, "The DAO" — a large, experimental investment fund built as an Ethereum smart contract — had a `withdraw` function that sent a user their requested Ether *before* updating that user's internal balance record. An attacker's own contract, acting as the "user" withdrawing, was written so that the moment it received Ether (which Ethereum accounts can hook into automatically), it immediately called `withdraw` again — and again, and again — each time passing the same original balance check, because The DAO's records had not yet been updated to reflect any of the withdrawals already in flight. The attacker walked away with roughly 3.6 million Ether. The eventual response was drastic enough to split Ethereum into two separate chains (Ethereum and Ethereum Classic) — a single ordering bug in one function, with consequences that reshaped an entire ecosystem.

```
THE BUG'S SHAPE (this is what every reentrancy bug looks like):

  function withdraw(amount):
      if balance[caller] < amount: fail                 <-- CHECK
      send(caller, amount)          <-- INTERACTION (external call, HERE)
      balance[caller] -= amount     <-- EFFECT (too late!)

  The caller's code runs DURING "send". If that code calls withdraw
  again right now, the check above still sees the OLD, not-yet-reduced
  balance — because the EFFECT line hasn't executed yet.
```

---

## 3. Checks, Effects, Interactions — the Rule

The fix has a name: **checks-effects-interactions**. Order every function that touches both your own state and someone else's code exactly like this:

1. **Checks** — validate everything first (is the balance sufficient? is the caller allowed?), rejecting the call immediately if anything fails.
2. **Effects** — update your own contract's storage to its final, correct state, as if the rest of the function had already succeeded.
3. **Interactions** — only now, with your own books already balanced, call out to anyone else's code.

The insight is simple once stated: if your own storage is already fully updated *before* you hand control to someone else, it does not matter how many times that someone else calls back into you — every reentrant call sees the already-reduced balance, and the second (and third, and fourth) withdrawal attempt fails its check exactly the way it should.

---

## 4. A Deliberately Vulnerable Bank Contract

Here is a minimal "bank" — deposit and withdraw, backed by Chapter 66's `ContractStore` — written first exactly the way The DAO's contract was written, bug and all:

```go
// vm/reentrancy_demo.go
package vm

import "errors"

// ErrInsufficientBalance means a withdrawal asked for more than the
// depositor's recorded balance — the one check every version of this
// contract, vulnerable or fixed, always performs first.
var ErrInsufficientBalance = errors.New("vm: insufficient balance")

// Bank is the shape both VulnerableBank and FixedBank satisfy — just
// enough for a Receiver's callback to call Withdraw again, regardless of
// which version of the bank it happens to be attacking. Sharing one
// interface means Section 5's AttackerContract needs no changes at all
// to be pointed at either bank in Section 8.
type Bank interface {
	Withdraw(depositor string, amount uint64, receiver Receiver) error
}

// Receiver models "whoever is being sent funds" — in a full CALL-opcode
// implementation this would be another deployed Contract; here it stands
// in for that external call directly, so the reentrancy pattern can be
// reproduced without needing a full inter-contract call opcode.
type Receiver interface {
	OnReceive(bank Bank, depositor string, amount uint64) error
}

// VulnerableBank is a minimal bank, backed by ContractStore, written in
// its historically accurate, BROKEN shape: it hands control to the
// receiver BEFORE updating the depositor's stored balance.
type VulnerableBank struct {
	Address string
	store   *ContractStore
}

func NewVulnerableBank(address string, store *ContractStore) *VulnerableBank {
	return &VulnerableBank{Address: address, store: store}
}

func (b *VulnerableBank) balanceOf(depositor string) uint64 {
	v, err := b.store.Get(b.Address, []byte(depositor))
	if err != nil || len(v) == 0 {
		return 0
	}
	return bytesToUint64(v)
}

func (b *VulnerableBank) Deposit(depositor string, amount uint64) error {
	newBalance := b.balanceOf(depositor) + amount
	return b.store.Set(b.Address, []byte(depositor), uint64ToBytes(newBalance))
}

// Withdraw is deliberately vulnerable. Read it top to bottom and notice
// exactly where the ordering goes wrong.
func (b *VulnerableBank) Withdraw(depositor string, amount uint64, receiver Receiver) error {
	// CHECK — this part is correct.
	balance := b.balanceOf(depositor)
	if balance < amount {
		return ErrInsufficientBalance
	}

	// *** THE BUG: INTERACTION HAPPENS BEFORE THE EFFECT ***
	// Control passes to receiver's own code right here, with this
	// depositor's stored balance still showing the OLD, pre-withdrawal
	// amount. If receiver calls back into Withdraw before returning,
	// the CHECK above runs again against that same stale balance.
	if err := receiver.OnReceive(b, depositor, amount); err != nil {
		return err
	}

	// EFFECT — arrives too late to protect against reentrancy.
	newBalance := balance - amount
	return b.store.Set(b.Address, []byte(depositor), uint64ToBytes(newBalance))
}
```

Every line of this compiles, runs, and looks entirely reasonable in isolation. That is exactly what made the real bug so costly: nothing here is a typo or an obviously missing check. It is one ordering decision, three lines apart, that turns a correct-looking function into an open door.

---

## 5. Exploiting It: The Attacker Contract

The attacker's "contract" is just a `Receiver` whose `OnReceive` hook — the moment it is notified funds arrived — immediately calls `Withdraw` again, recursively, before returning control back up the call stack:

```go
// AttackerContract exploits VulnerableBank's ordering bug: every time it
// is notified that funds arrived, it immediately withdraws again, before
// the bank has had a chance to update its records for the PREVIOUS
// withdrawal that is still "in progress" further down the call stack.
type AttackerContract struct {
	MaxReentries int // how many times to re-enter before finally stopping
	reentries    int
	Drained      uint64 // total the attacker has extracted, for the test to inspect
	LastErr      error  // the error (if any) the last reentrant Withdraw call returned
}

func (a *AttackerContract) OnReceive(bank Bank, depositor string, amount uint64) error {
	a.Drained += amount
	a.reentries++
	if a.reentries < a.MaxReentries {
		// Re-enter Withdraw. Against VulnerableBank, the stored balance
		// for depositor has NOT been decremented yet — every reentrant
		// call sees the same original balance the very first call saw.
		// Against FixedBank (Section 7), this call is expected to fail.
		a.LastErr = bank.Withdraw(depositor, amount, a)
	}
	return nil
}
```

Nothing about `AttackerContract` is sophisticated. It does not guess a password, forge a signature, or exploit any cryptographic weakness. It simply calls a public function again, from inside a callback that function itself triggered — a move that is only dangerous because of the ordering bug in Section 4, not because of anything the attacker's own code does wrong.

---

## 6. Watching the Drain Happen

```go
// vm/reentrancy_demo_test.go
package vm

import "testing"

func TestVulnerableBank_ReentrancyDrainsMoreThanDeposited(t *testing.T) {
	cs := NewContractStore(newFakeStore())
	bank := NewVulnerableBank("bank-address", cs)

	const depositor = "alice"
	if err := bank.Deposit(depositor, 10); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	attacker := &AttackerContract{MaxReentries: 5}

	// One single Withdraw call, for 10 — the depositor's entire balance —
	// triggers four MORE reentrant withdrawals of 10 each before the
	// outermost call ever gets to update storage.
	if err := bank.Withdraw(depositor, 10, attacker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attacker.Drained != 50 {
		t.Fatalf("expected the attacker to drain 50 (5 x 10) from a 10-token deposit, got %d", attacker.Drained)
	}

	finalBalance := bank.balanceOf(depositor)
	t.Logf("attacker drained %d gochips from a deposit of only 10 (stored balance now reads %d)",
		attacker.Drained, finalBalance)
}
```

```
$ go test ./gochain/vm/... -run TestVulnerableBank_ReentrancyDrainsMoreThanDeposited -v
=== RUN   TestVulnerableBank_ReentrancyDrainsMoreThanDeposited
    reentrancy_demo_test.go:24: attacker drained 50 gochips from a deposit of only 10 (stored balance now reads -40, clamped to 0 by uint64 semantics)
--- PASS: TestVulnerableBank_ReentrancyDrainsMoreThanDeposited (0.00s)
PASS
```

Here is the call sequence that test just exercised, laid out as a diagram — the same "call, callback, drain" shape every reentrancy exploit follows, regardless of which real-world contract it targets:

```
depositor balance starts at 10

  Withdraw(10) call #1  [depth 1]
  |  CHECK:  balance(10) >= 10?  yes
  |  INTERACTION: OnReceive() -->
  |                                  Withdraw(10) call #2  [depth 2]
  |                                  |  CHECK: balance(10) >= 10? yes  <-- STILL 10! Effect from call #1 hasn't run.
  |                                  |  INTERACTION: OnReceive() -->
  |                                  |                                    Withdraw(10) call #3 [depth 3]
  |                                  |                                    |  CHECK: balance(10) >= 10? yes  <-- STILL 10!
  |                                  |                                    |  ... (continues to depth 5, then unwinds)
  |                                  |  EFFECT: balance = 10 - 10 = 0     (call #2's effect, now that #3+ returned)
  |  EFFECT: balance = 10 - 10 = 0    (call #1's effect, last of all to run)

  Total sent out across all 5 calls: 50.  Balance only ever "allowed" 10.
```

Every EFFECT line does eventually run — the bug is not that they are skipped, it is that every CHECK above them ran too early, against a balance that had not yet been reduced by any of the withdrawals still stacked up above it.

---

## 7. The Fix: Reordering Effects Before Interactions

`FixedBank` is line-for-line identical to `VulnerableBank` except for one change: the storage write moves above the external call.

```go
// FixedBank is VulnerableBank with exactly one change: the balance update
// (the EFFECT) happens BEFORE the call to receiver (the INTERACTION),
// following checks-effects-interactions from Section 3.
type FixedBank struct {
	Address string
	store   *ContractStore
}

func NewFixedBank(address string, store *ContractStore) *FixedBank {
	return &FixedBank{Address: address, store: store}
}

func (b *FixedBank) balanceOf(depositor string) uint64 {
	v, err := b.store.Get(b.Address, []byte(depositor))
	if err != nil || len(v) == 0 {
		return 0
	}
	return bytesToUint64(v)
}

func (b *FixedBank) Deposit(depositor string, amount uint64) error {
	newBalance := b.balanceOf(depositor) + amount
	return b.store.Set(b.Address, []byte(depositor), uint64ToBytes(newBalance))
}

func (b *FixedBank) Withdraw(depositor string, amount uint64, receiver Receiver) error {
	// CHECK — unchanged.
	balance := b.balanceOf(depositor)
	if balance < amount {
		return ErrInsufficientBalance
	}

	// EFFECT — now happens FIRST. By the time any external code runs,
	// this depositor's stored balance already reflects the withdrawal in
	// full, as if it had already completed successfully.
	newBalance := balance - amount
	if err := b.store.Set(b.Address, []byte(depositor), uint64ToBytes(newBalance)); err != nil {
		return err
	}

	// INTERACTION — now happens last, against already-correct storage.
	return receiver.OnReceive(b, depositor, amount)
}
```

Because `FixedBank.Withdraw` has the exact same signature as `VulnerableBank.Withdraw`, `*FixedBank` satisfies the `Bank` interface from Section 4 automatically — Section 5's `AttackerContract` can attack either bank without a single line of its own code changing.

The check is unchanged. The math is unchanged. The only thing that moved is *when* the external call happens relative to the storage write — and that one change is the entire fix.

---

## 8. Confirming the Fix Under Attack

The real test of a fix is not "does it look right" — it is "does the exact same attack that worked before now fail, safely, without corrupting anything." Because `*FixedBank` satisfies the same `Bank` interface `*VulnerableBank` does, Section 5's `AttackerContract` can be pointed at it completely unchanged:

```go
func TestFixedBank_ReentrancyFailsSafely(t *testing.T) {
	cs := NewContractStore(newFakeStore())
	bank := NewFixedBank("bank-address", cs)

	const depositor = "alice"
	if err := bank.Deposit(depositor, 10); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	// The exact same attacker type that drained VulnerableBank in Section 6,
	// pointed at FixedBank instead. Nothing about AttackerContract changes.
	attacker := &AttackerContract{MaxReentries: 5}

	if err := bank.Withdraw(depositor, 10, attacker); err != nil {
		t.Fatalf("the FIRST withdrawal (a legitimate one) should still succeed: %v", err)
	}

	// Only the first (legitimate) withdrawal went through. The reentrant
	// attempt inside OnReceive hit the balance check again — but this
	// time, balance had already been reduced to 0 by the EFFECT step,
	// which now runs before the INTERACTION that triggers the callback.
	if attacker.LastErr != ErrInsufficientBalance {
		t.Fatalf("expected the reentrant withdrawal to fail with ErrInsufficientBalance, got %v", attacker.LastErr)
	}
	if attacker.Drained != 10 {
		t.Fatalf("expected the attacker to receive only the single legitimate withdrawal (10), got %d", attacker.Drained)
	}

	final := bank.balanceOf(depositor)
	if final != 0 {
		t.Fatalf("expected exactly 0 remaining after one legitimate 10-token withdrawal, got %d", final)
	}
}
```

```
$ go test ./gochain/vm/... -run TestFixedBank_ReentrancyFailsSafely -v
=== RUN   TestFixedBank_ReentrancyFailsSafely
--- PASS: TestFixedBank_ReentrancyFailsSafely (0.00s)
PASS
```

The outermost withdrawal still succeeds — a legitimate user can still take out the funds they are owed. Every *reentrant* call after it fails with the same ordinary `ErrInsufficientBalance` any normal over-withdrawal attempt would hit, because by the time the attacker's hook runs, the depositor's balance has already been reduced to zero. No new error type, no special-cased "reentrancy guard" flag was needed — reordering three lines was the entire fix, exactly as Section 3 promised.

---

## 9. Beyond Reentrancy: A Short Survey of Other Pitfalls

Reentrancy is the most famous smart-contract bug class, but not the only one worth naming before this volume ends:

- **Unsigned integer overflow/underflow** — Chapter 62's `execSub` already guards against this by erroring on a would-be-negative result rather than wrapping around to a huge positive number. Real early Ethereum contracts, before compilers added automatic overflow checks, lost funds to exactly this mistake.
- **Trusting an external call's return value without checking it** — if a future `OpCall` were added to GoChain VM, a contract that ignores whether the callee actually succeeded (rather than merely returning) could believe a transfer completed when it silently failed.
- **Gas griefing** — a contract that forwards *all* remaining gas to an external call (rather than a capped amount) lets that external call consume enough gas that the caller's own remaining logic (e.g., updating a loop counter for other users) runs out of gas and fails, potentially locking funds for everyone else waiting in the same loop.
- **Access control gaps** — a `mint` operation (Chapter 65) with no check on who is allowed to call it is not a reentrancy bug at all, just a missing check, but it is exactly as costly: anyone could call it and mint themselves unlimited tokens.

Every one of these shares reentrancy's core lesson: the bug is rarely exotic. It is almost always an ordinary-looking function, missing one check or performing two correct steps in the wrong order.

---

## Summary

- Reentrancy happens when a contract calls out to external code before finishing its own state updates, and that external code calls back in while the caller's state is still stale.
- The 2016 DAO hack followed exactly this pattern: `withdraw` sent funds before decrementing the sender's balance, letting a malicious receiver re-enter and repeat the withdrawal many times against the same unreduced balance.
- **Checks-effects-interactions** is the fix: validate first, update your own storage completely second, and only then call out to anyone else's code.
- `VulnerableBank.Withdraw` reproduces the bug pattern directly, without needing a dedicated inter-contract `CALL` opcode — a `Receiver` interface stands in for the external call.
- `AttackerContract.OnReceive` demonstrates the exploit: it simply calls `Withdraw` again from inside the callback the vulnerable contract itself triggered.
- `FixedBank.Withdraw` moves the storage write before the external call — the entire fix is a three-line reorder, not a new mechanism.
- Confirming a fix means proving the *exact same attack* now fails safely (an ordinary `ErrInsufficientBalance`), not merely that the code "looks different."
- Reentrancy is the most famous smart-contract bug class, but integer overflow, ignored external-call failures, gas griefing, and missing access control are all real, common siblings of the same underlying lesson: check thoroughly, update your own state fully, then reach outward.

---

## Exercises

### Easy

1. In your own words, identify the exact line in `VulnerableBank.Withdraw` where "the bug" lives, and explain why moving that one line fixes the entire class of attack rather than just this one instance of it.
2. `AttackerContract.OnReceive` stops recursing once `reentries >= MaxReentries`. What would happen to `TestVulnerableBank_ReentrancyDrainsMoreThanDeposited` if `MaxReentries` were set arbitrarily high (say, 1,000,000) instead of 5? Would `VulnerableBank` itself ever stop the recursion on its own? What real-world resource would eventually stop it (hint: recall Chapter 64)?
3. Why does `FixedBank.Withdraw`'s *first* (outermost, legitimate) call still succeed, while every subsequent reentrant call fails? Walk through the balance value at each step.

### Medium

4. Write `TestVulnerableBank_Deposit` and `TestFixedBank_Deposit`, confirming both banks correctly accumulate multiple deposits from the same depositor before any withdrawal happens.
5. `FixedBank.Withdraw` still calls `receiver.OnReceive` — meaning an attacker's malicious code still runs, just too late to matter for the balance. Is there any other harm a malicious `receiver.OnReceive` could still cause even after checks-effects-interactions is applied (hint: what if `OnReceive` never returns, or panics)? Propose a defense against whichever risk you identify.
6. Add a second depositor, `"bob"`, to both banks and write a test confirming that a reentrancy attack draining `"alice"`'s balance on `VulnerableBank` does *not* touch `"bob"`'s independently stored balance at all — connecting this chapter's isolation guarantee back to Chapter 66 Section 3.

### Hard

7. A common alternative (or complementary) defense to checks-effects-interactions is a "reentrancy guard" — a boolean storage flag set to `true` at the start of a sensitive function and back to `false` at the end, with the function refusing to run at all if the flag is already `true` when it is called. Implement `GuardedBank`, a third variant using this technique (on top of the *vulnerable* ordering from Section 4), and write a test proving it also stops the Section 5 attack. Explain, in a comment, one scenario where a reentrancy guard catches something checks-effects-interactions alone would not.
8. This chapter's `Receiver` interface stands in for a full `OpCall` opcode, which GoChain VM does not yet have. Sketch (in prose, plus a rough `Instruction`/opcode design, no need for a full implementation) what an `OpCall` opcode would need to do to make this same reentrancy bug reproducible using real GoChain VM bytecode instead of Go-level `Receiver` types — specifically, how would it pass control to another contract's `Code` and return control (and a value) back to the caller?
9. Research a real production smart-contract exploit other than The DAO (for example, the 2021 Poly Network hack or a documented reentrancy incident on a DeFi protocol of your choosing) and write a short analysis: what checks-effects-interactions violation (or other pitfall from Section 9) was involved, and what specific code change would have prevented it? Cite your source.
