# Chapter 68: Testing Smart Contracts

Every package this course has built so far has shipped with tests — Chapter 62 tested every opcode, Chapter 55 tested `BoltStore` against a temporary directory. Smart contracts deserve more than that same baseline, and this chapter explains why before writing a single test: a bug in a contract's `withdraw` function is not a bug you patch tomorrow. Once a contract is deployed and other people's funds are sitting in its storage, its code is effectively permanent — Chapter 67's entire DAO story exists because a bug that would have been an ordinary embarrassing pull request in most software instead became a fork of an entire blockchain. This chapter writes the test suite the token contract (Chapter 65) and the storage layer (Chapter 66) deserve: one test per operation, an adversarial reentrancy test built directly on Chapter 67's attacker, and gas-consumption assertions that turn "this got slower" into a red test rather than a surprise three months later.

## Table of Contents

1. [Recap: Why Contract Bugs Are Different](#1-recap-why-contract-bugs-are-different)
2. [The Token Contract's Public Shape](#2-the-token-contracts-public-shape)
3. [Unit Tests for Mint, Transfer, and BalanceOf](#3-unit-tests-for-mint-transfer-and-balanceof)
4. [Table-Driven Edge Cases](#4-table-driven-edge-cases)
5. [Adversarial Test: Reentrancy Must Now Fail Safely](#5-adversarial-test-reentrancy-must-now-fail-safely)
6. [Gas-Consumption Assertions](#6-gas-consumption-assertions)
7. [Why Contract Bugs Are Especially Costly](#7-why-contract-bugs-are-especially-costly)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Recap: Why Contract Bugs Are Different

Ordinary software has a safety net most engineers take for granted: if a bug ships, you fix it and deploy again. A blockchain's whole reason for existing works against that safety net. Once a contract is deployed at an address, and other users' balances are sitting in its storage (Chapter 66), the contract's code is what every node on the network has already agreed to trust — you cannot quietly swap in corrected bytecode the way you would push a hotfix to a web server, without either convincing every node to accept a brand-new, differently-addressed contract (leaving the old, buggy one still reachable) or, in the most extreme case Chapter 67 described, convincing the entire network to rewrite history. A contract's test suite is, in a very real sense, its only line of defense that exists *before* the cost of a mistake becomes irreversible.

```
ORDINARY WEB SERVICE                     DEPLOYED SMART CONTRACT

  bug shipped                              bug shipped
      |                                        |
      v                                        v
  user notices / monitoring alerts         attacker notices first,
      |                                    usually before you do
      v                                        |
  hotfix deployed in minutes                   v
      |                                    funds already moved,
      v                                    irreversibly, before
  incident over, cost = downtime          any "hotfix" is possible
                                                |
                                                v
                                    cost = the funds themselves,
                                    not just engineering time

  Tests catch a bug BEFORE deploy.  For a contract, that "before" is
  the only chance you get — there is no equivalent of the left-hand
  column's "hotfix deployed in minutes" step once real value is at
  stake in storage other people already trust.
```

This is not an argument for testing contracts *differently* from other code — every technique in this chapter is ordinary Go testing, nothing exotic. It is an argument for testing them **at least as thoroughly as any other financial system you would ever be asked to build**, and treating a passing test suite as a precondition for deployment, not a nice-to-have.

---

## 2. The Token Contract's Public Shape

Chapter 65 designed the token contract's three operations — `mint`, `transfer`, `balanceOf` — as small VM programs, dispatched by method name through `Contract.Call` exactly the way Chapter 66 Section 7 built it to work. Each operation's actual bytecode is Chapter 65's job, opcode by opcode, the same way Chapter 62 built each of `execAdd` and `execSub` one method at a time. What this chapter's test suite needs is the *behavior* that bytecode promises — and, following the same pattern Chapter 67 used for `VulnerableBank` and `FixedBank`, the cleanest way to pin that behavior down for testing is a small Go type built directly on `ContractStore` (Chapter 66), exposing exactly the operations the deployed bytecode implements:

```go
// vm/token.go
package vm

import "errors"

// ErrInsufficientTokenBalance means a transfer asked to move more than
// the sender's recorded balance — the one check every transfer performs
// before touching storage at all.
var ErrInsufficientTokenBalance = errors.New("vm: insufficient token balance")

// Token is the fungible token contract's storage-backed surface: the
// same mint/transfer/balanceOf behavior Chapter 65's VM bytecode
// implements opcode by opcode, expressed here directly against
// ContractStore so it can be tested with the same clarity Chapter 67
// tested VulnerableBank and FixedBank.
type Token struct {
	Address string
	store   *ContractStore
}

// NewToken derives a contract address (the same way DeployContract does
// in Chapter 66) and binds it to store.
func NewToken(address string, store *ContractStore) *Token {
	return &Token{Address: address, store: store}
}

func (t *Token) BalanceOf(address string) uint64 {
	v, err := t.store.Get(t.Address, []byte(address))
	if err != nil || len(v) == 0 {
		return 0
	}
	return bytesToUint64(v)
}

func (t *Token) Mint(to string, amount uint64) error {
	newBalance := t.BalanceOf(to) + amount
	return t.store.Set(t.Address, []byte(to), uint64ToBytes(newBalance))
}

func (t *Token) Transfer(from, to string, amount uint64) error {
	fromBalance := t.BalanceOf(from)
	if fromBalance < amount {
		return ErrInsufficientTokenBalance
	}
	if err := t.store.Set(t.Address, []byte(from), uint64ToBytes(fromBalance-amount)); err != nil {
		return err
	}
	toBalance := t.BalanceOf(to)
	return t.store.Set(t.Address, []byte(to), uint64ToBytes(toBalance+amount))
}
```

Every method here does at the Go level exactly what Chapter 65's dispatcher does at the bytecode level: `BalanceOf` is one `OpSLoad`; `Mint` is a load, an add, and a store; `Transfer` checks the sender's balance, then performs two load-modify-store sequences. Testing against this surface tests the behavior a real deployed contract promises, without this chapter needing to re-derive Chapter 65's exact opcode sequence to do it — the same separation of concerns Chapter 55's `storage.Store` interface established between "what a caller needs" and "how it happens to be implemented underneath."

---

## 3. Unit Tests for Mint, Transfer, and BalanceOf

Each operation gets a focused test, following the same one-behavior-per-test discipline Chapter 62 used for opcodes:

```go
// vm/token_test.go
package vm

import "testing"

func newTestToken(t *testing.T) *Token {
	t.Helper()
	cs := NewContractStore(newFakeStore())
	return NewToken("token-contract-address", cs)
}

func TestToken_Mint(t *testing.T) {
	token := newTestToken(t)

	if err := token.Mint("alice", 1000); err != nil {
		t.Fatalf("mint: %v", err)
	}

	if balance := token.BalanceOf("alice"); balance != 1000 {
		t.Fatalf("expected balance 1000 after mint, got %d", balance)
	}
}

func TestToken_Transfer(t *testing.T) {
	token := newTestToken(t)

	if err := token.Mint("alice", 1000); err != nil {
		t.Fatalf("mint: %v", err)
	}
	if err := token.Transfer("alice", "bob", 300); err != nil {
		t.Fatalf("transfer: %v", err)
	}

	if balance := token.BalanceOf("alice"); balance != 700 {
		t.Fatalf("expected alice's balance to be 700 after sending 300, got %d", balance)
	}
	if balance := token.BalanceOf("bob"); balance != 300 {
		t.Fatalf("expected bob's balance to be 300 after receiving it, got %d", balance)
	}
}

func TestToken_BalanceOf_NeverMinted(t *testing.T) {
	token := newTestToken(t)

	if balance := token.BalanceOf("nobody"); balance != 0 {
		t.Fatalf("expected 0 for an address that was never minted to, got %d", balance)
	}
}
```

`TestToken_BalanceOf_NeverMinted` matters as much as the other two: it confirms Chapter 66's "an unwritten storage slot reads back as empty, not an error" guarantee holds all the way up through the token contract's own public behavior, not just at `ContractStore`'s own level.

---

## 4. Table-Driven Edge Cases

Every contract operation has at least one way to be asked to do something it should refuse. A table-driven test keeps these compact and, just as importantly, makes it obvious at a glance which edge cases have and have not been covered yet:

```go
func TestToken_Transfer_RejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name       string
		mintTo     string
		mintAmount uint64
		from, to   string
		amount     uint64
		wantErr    bool
	}{
		{
			name: "insufficient balance", mintTo: "alice", mintAmount: 100,
			from: "alice", to: "bob", amount: 500, wantErr: true,
		},
		{
			name: "transfer from an address that never received any tokens",
			from: "carol", to: "bob", amount: 1, wantErr: true,
		},
		{
			name: "transfer of exactly the full balance succeeds", mintTo: "dave", mintAmount: 50,
			from: "dave", to: "erin", amount: 50, wantErr: false,
		},
		{
			name: "transfer of zero is allowed and a no-op", mintTo: "frank", mintAmount: 10,
			from: "frank", to: "grace", amount: 0, wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := newTestToken(t)
			if tt.mintTo != "" {
				if err := token.Mint(tt.mintTo, tt.mintAmount); err != nil {
					t.Fatalf("mint setup: %v", err)
				}
			}

			err := token.Transfer(tt.from, tt.to, tt.amount)
			if tt.wantErr && err == nil {
				t.Fatalf("expected an error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
```

The "transfer of zero is allowed and a no-op" case is worth calling out: it would be easy to assume zero-amount transfers are meaningless and skip testing them, but a contract that mishandled a zero transfer (for instance, by treating a zero balance check as "insufficient" even though zero minus zero is a perfectly valid zero) would fail silently in exactly the kind of edge case a real user could trigger by accident — one more reason table-driven tests earn their place: they make "did we think about zero?" a permanent, visible line in the test file instead of a question someone has to remember to ask.

`Mint` deserves the identical treatment, not just `Transfer` — it is tempting to assume the "simpler" operation needs fewer tests, but `Mint` is exactly where a token's total supply is decided, and a single mistaken edge case here (minting a negative-looking amount that somehow wraps, or silently overwriting rather than accumulating a balance) is just as costly as a transfer bug:

```go
func TestToken_Mint_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		mints       []struct {
			to     string
			amount uint64
		}
		checkAddr string
		wantTotal uint64
	}{
		{
			name: "a single mint sets the balance directly",
			mints: []struct {
				to     string
				amount uint64
			}{{"alice", 500}},
			checkAddr: "alice", wantTotal: 500,
		},
		{
			name: "repeated mints to the same address accumulate, not overwrite",
			mints: []struct {
				to     string
				amount uint64
			}{{"bob", 100}, {"bob", 50}, {"bob", 25}},
			checkAddr: "bob", wantTotal: 175,
		},
		{
			name: "minting zero is allowed and leaves the balance unchanged",
			mints: []struct {
				to     string
				amount uint64
			}{{"carol", 40}, {"carol", 0}},
			checkAddr: "carol", wantTotal: 40,
		},
		{
			name: "mints to different addresses never bleed into each other",
			mints: []struct {
				to     string
				amount uint64
			}{{"dave", 1000}, {"erin", 1}},
			checkAddr: "dave", wantTotal: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := newTestToken(t)
			for _, m := range tt.mints {
				if err := token.Mint(m.to, m.amount); err != nil {
					t.Fatalf("mint(%s, %d): %v", m.to, m.amount, err)
				}
			}
			if got := token.BalanceOf(tt.checkAddr); got != tt.wantTotal {
				t.Fatalf("expected %s's balance to be %d, got %d", tt.checkAddr, tt.wantTotal, got)
			}
		})
	}
}
```

The "mints to different addresses never bleed into each other" case connects directly back to Chapter 66 Section 3's storage-isolation guarantee — but proven here at the token contract's own level, not just at `ContractStore`'s. A regression in either layer would be caught by one of these two test suites, which is exactly the kind of overlapping coverage a financial system deserves rather than relying on a single layer to catch every possible mistake.

---

## 5. Adversarial Test: Reentrancy Must Now Fail Safely

The token contract's `transfer` operation, as built in Sections 2 through 4, never calls out to anyone else's code — so it is not itself reentrancy-vulnerable, and there is no bug to reproduce inside it. What *does* need an adversarial test, explicitly and by name, is Chapter 67's fix — because a fix that is never tested against the exact attack it claims to defeat is only a fix by assertion, not by evidence. This test belongs in the same test suite as the token contract's own tests precisely because "does this codebase's security fix still hold" is exactly the kind of question a contract's test suite exists to answer, permanently, on every future change:

```go
// vm/token_test.go (continued) — reusing Chapter 67's FixedBank and
// AttackerContract directly, so this test breaks immediately if a future
// change to FixedBank.Withdraw ever reintroduces the ordering bug.
func TestFixedBank_ReentrancyAdversarialRegression(t *testing.T) {
	cs := NewContractStore(newFakeStore())
	bank := NewFixedBank("bank-address", cs)

	const depositor = "alice"
	if err := bank.Deposit(depositor, 10); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	attacker := &AttackerContract{MaxReentries: 5}
	if err := bank.Withdraw(depositor, 10, attacker); err != nil {
		t.Fatalf("the legitimate withdrawal should still succeed: %v", err)
	}

	// The whole point of this test: assert the attack FAILS, not that it
	// happens to produce some particular number. If a future refactor
	// accidentally reorders FixedBank.Withdraw back to Chapter 67's
	// vulnerable shape, attacker.Drained jumps from 10 to 50 and this
	// assertion catches it immediately.
	if attacker.Drained != 10 {
		t.Fatalf("REGRESSION: reentrancy attack drained %d, expected exactly 10 (attack must fail safely)", attacker.Drained)
	}
	if attacker.LastErr != ErrInsufficientBalance {
		t.Fatalf("REGRESSION: expected the reentrant call to fail with ErrInsufficientBalance, got %v", attacker.LastErr)
	}
}
```

This test is deliberately named and commented as a *regression* test, not merely a correctness test: its entire job, forever, is to fail loudly the instant anyone (including a future version of yourself, six months from now, "simplifying" `FixedBank.Withdraw`) reintroduces the exact ordering bug Chapter 67 fixed.

---

## 6. Gas-Consumption Assertions

Chapter 64 gave every opcode a real gas cost; without a test watching it, gas usage can silently grow — someone adds one extra `OpDup`/`OpPop` pair "just to be safe" during a refactor, and now every call to `transfer` costs a little more than it used to, with nothing announcing the change until a much later chapter's benchmark (or a real user's unexpectedly expensive transaction) notices. A gas regression test turns that silent drift into an immediate, specific test failure.

```go
// GasUsed exposes how much gas the most recently completed Execute call
// consumed — the first exported way to read vm.gasUsed back out, added
// specifically so tests like this one can assert on it directly.
func (vm *VM) GasUsed() uint64 {
	return vm.gasUsed
}
```

A gas assertion needs real VM bytecode under measurement, not the `Token` Go type from Section 2 — `GasUsed()` only means something for code that actually ran through `Execute()`. Chapter 66 Section 8's small "store a value, then read it back" contract is exactly one storage write plus one storage read — the same shape of work every one of the token contract's own operations does at least once — which makes it a fair stand-in for measuring "does one representative storage round-trip still cost roughly what it used to":

```go
func TestContract_StorageRoundTrip_GasUnderCeiling(t *testing.T) {
	// The exact "set, then get" contract from Chapter 66, Section 8 —
	// OpPop, OpPush, OpSStore, OpPush, OpSLoad, OpHalt: one representative
	// storage write and one representative storage read.
	code := []Instruction{
		{Op: OpPop},
		{Op: OpPush, Arg: []byte("balance")},
		{Op: OpSStore},
		{Op: OpPush, Arg: []byte("balance")},
		{Op: OpSLoad},
		{Op: OpHalt},
	}
	contract := DeployContract(code)

	cs := NewContractStore(newFakeStore())
	v := NewVM(nil, 100000)
	v.AttachStorage(cs)

	if _, err := contract.Call(v, "set", [][]byte{uint64ToBytes(42)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// This ceiling is deliberately generous, not a tight bound — its job
	// is to catch an accidental 5x or 10x regression (e.g., an
	// accidentally introduced unbounded loop, or a duplicated storage
	// operation), not to enforce micro-optimization.
	const gasCeiling = 50
	if used := v.GasUsed(); used > gasCeiling {
		t.Fatalf("REGRESSION: one storage round-trip used %d gas, expected at most %d — check for accidentally added opcodes", used, gasCeiling)
	}
}
```

A ceiling this generous might look pointless — surely nothing this six-instruction program is anywhere close to 50 gas units. That gap is exactly the point: the test is not trying to catch "you added one extra `OpDup`," it is trying to catch "someone accidentally introduced an unbounded loop, or duplicated an entire storage operation" — the kind of mistake that turns a 6-gas operation into a 5,000-gas one, and that a developer staring at a diff might not notice, but a red test cannot miss. Every one of the token contract's real operations (Section 2) does at least one storage round-trip like this one, so a regression here is a regression everywhere that pattern repeats.

---

## 7. Why Contract Bugs Are Especially Costly

Three properties compound to make contract testing worth this much more effort than an equivalent web-application test suite:

- **No quiet patch.** A web server's bug is fixed by a deploy; a deployed contract's bytecode, once other users' balances live in its storage, cannot simply be edited in place — fixing it means either accepting the bug's consequences or, in the most extreme real-world cases, coordinating an entire network around an emergency change (exactly what Chapter 67's DAO story describes).
- **Adversarial by default.** A web application's users are mostly not actively trying to break it. Every contract deployed to a real, public chain should be assumed to have at least one person actively looking for a way to profit from a mistake in it — which is exactly why Section 5's adversarial test matters as much as Sections 3 and 4's ordinary-case tests.
- **Financial, not cosmetic.** A rendering bug in a web page is embarrassing. A rounding, ordering, or access-control bug in a contract is measured directly in funds that do not come back, which is the entire reason this chapter treats gas regressions and reentrancy regressions with the same seriousness as an incorrect balance.

None of this means contract tests need exotic new tools — every test in this chapter uses the same `testing.T`, table-driven pattern, and `t.Run` subtests Chapter 07 introduced in Volume 1. What changes is not the tooling; it is how much is riding on the tests actually being there, and actually being run, before code reaches a real chain.

---

## Summary

- Contract bugs are effectively irreversible once deployed and funded — there is no quiet hotfix, which is why this chapter treats contract testing as a first-class chapter of its own, not an afterthought bolted onto Chapter 65.
- The token contract's `mint`, `transfer`, and `balanceOf` operations each get a dedicated unit test, following the one-behavior-per-test discipline established since Chapter 62.
- Table-driven tests cover the edge cases a hand-written list of individual tests tends to forget: insufficient balance, an address that never received tokens, transferring an entire balance, and transferring exactly zero.
- The reentrancy fix from Chapter 67 gets its own explicit, named regression test — `TestFixedBank_ReentrancyAdversarialRegression` — reusing `AttackerContract` directly so any future reintroduction of the ordering bug fails immediately and loudly.
- `VM.GasUsed()` exposes gas consumption for the first time, letting a test assert a generous ceiling on an operation's cost, catching accidental regressions like an unbounded loop rather than enforcing micro-optimization.
- Contract bugs are costlier than ordinary software bugs for three compounding reasons: no quiet patching after deployment, an adversarial user base by default, and financial (not cosmetic) consequences.

---

## Exercises

### Easy

1. `TestToken_BalanceOf_NeverMinted` asserts a balance of 0 for an address that was never minted to. Trace through `ContractStore.Get` (Chapter 66) and explain, in your own words, exactly which line makes this the correct behavior rather than an error.
2. Why is "transfer of zero is allowed and a no-op" included as its own table-driven case in Section 4, rather than being assumed to obviously work? What specific implementation mistake would this case catch that a nonzero-only test suite would miss?
3. `TestFixedBank_ReentrancyAdversarialRegression` is described as a regression test rather than a correctness test. What is the practical difference in how you should react if this specific test starts failing, versus how you'd react if `TestToken_Mint` started failing?

### Medium

4. Write `TestToken_Mint_AccumulatesAcrossMultipleCalls` confirming that minting to the same address three times in a row (100, then 50, then 25) leaves that address with a balance of 175, not just the most recent mint amount.
5. The gas ceiling in Section 6 is described as "deliberately generous." Rewrite `TestContract_StorageRoundTrip_GasUnderCeiling` as a table-driven test covering a program with one storage write only, one storage read only, and the combined round-trip, each with their own appropriately generous (but real) ceiling, and justify each number you choose in a comment.
6. Section 2 implements `Token.Mint` and `Token.Transfer` directly against `ContractStore`, standing in for Chapter 65's actual VM bytecode. Try hand-deriving `mint`'s bytecode yourself using only `OpDup`, `OpPop`, `OpPush`, `OpAdd`, `OpSLoad`, `OpSStore`, and `OpHalt` (the opcodes fixed since Chapter 61) — specifically, once you have popped the recipient address and computed `existingBalance + amount` with `OpAdd`, is the recipient address still anywhere on the stack for `OpSStore` to use as a key? If not, what additional opcode (hint: something Ethereum's own stack machine relies on constantly) would you need to add to the instruction set to make this possible, and what would its exact stack effect be?

### Hard

7. Design and implement a "fuzz"-style test (using Go's built-in `testing.F` fuzzing support) for `Token.Transfer` that generates random mint amounts and transfer amounts (including amounts larger than `uint64` values a real balance could plausibly reach) and asserts one invariant that must always hold regardless of input: the total token supply (sum of every address's balance) never changes across a transfer, only across a mint. Explain what class of bug this invariant would catch that Section 3's fixed-input tests would not.
8. Chapter 64's gas costs are presumably different for `OpSLoad`/`OpSStore` (storage operations) than for `OpAdd`/`OpDup` (pure stack operations), reflecting that real storage writes are much more expensive on an actual blockchain than in-memory stack math. Without access to Chapter 64's exact cost table, write a test that empirically determines (by running two minimal programs and comparing `GasUsed()`) whether your project's actual implementation reflects this expectation, and explain what it would mean, from a security standpoint, if storage operations were priced the same as stack operations.
9. Propose and implement a `TestToken_ConcurrentTransfers` test using multiple goroutines, each holding its own `*VM` but sharing one `ContractStore` (and therefore one underlying `storage.Store`), all transferring from the same sender's balance simultaneously. Determine experimentally whether the sender's final balance is always correct, and if it is not, explain — referencing Chapter 66 Exercise 7 — exactly where the race condition lives and what would need to change in `ContractStore` or `Contract.Call` to fix it.
