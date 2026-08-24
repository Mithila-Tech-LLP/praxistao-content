# Chapter 64: Gas and Execution Limits

Chapter 62's `Execute()` loop already has `gasUsed` and `gasLimit` fields, and a placeholder that charges a flat 1 gas per instruction "just enough to prove the mechanism works." That placeholder has a hole big enough to hang the entire network: nothing about it says any opcode costs *more* than any other, and nothing yet stops a program from looping forever, one cheap instruction at a time, tying up every node that tries to validate it. This chapter replaces the placeholder with a real, per-opcode cost table, and proves — with a program built specifically to never stop on its own — that gas is what makes "never stop on its own" survivable instead of catastrophic.

## Table of Contents

1. [The Problem Gas Solves](#1-the-problem-gas-solves)
2. [What Gas Actually Measures](#2-what-gas-actually-measures)
3. [Designing a Per-Opcode Cost Table](#3-designing-a-per-opcode-cost-table)
4. [Wiring Costs Into the Execute Loop](#4-wiring-costs-into-the-execute-loop)
5. [ErrOutOfGas and Failing Predictably](#5-erroutofgas-and-failing-predictably)
6. [Testing an Ordinary Program's Gas Cost](#6-testing-an-ordinary-programs-gas-cost)
7. [Testing a Deliberate Infinite Loop](#7-testing-a-deliberate-infinite-loop)
8. [Choosing Gas Limits in Practice](#8-choosing-gas-limits-in-practice)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Problem Gas Solves

Chapter 61 pointed out, almost in passing, that `OpJump` and `OpJumpIfFalse` together are "enough to express both `if`/`else` branching and loops (a loop is simply a jump backward to an earlier instruction index)." That is a genuinely useful capability — Chapter 65's token contract needs it, and any nontrivial contract logic needs it. But it is also a loaded gun: nothing about `OpJump` stops a program from jumping backward to instruction 0, forever.

```go
// A perfectly legal GoChain VM program, as far as Chapter 62's Execute()
// is concerned:
program := []Instruction{
    {Op: OpJump, Arg: uint64ToBytes(0)}, // jump to myself, forever
}
```

Without a limit, calling `NewVM(program, someLimit).Execute()` on this program would never return. Every node in the network validating a block that contained a transaction invoking this "contract" would hang at the exact same instruction, forever, the moment it tried to reach consensus on whether the transaction was valid. One buggy or malicious program, submitted once, would be enough to stop the entire network from making progress — a far worse outcome than a single failed transaction.

**Gas** is the fix: every opcode costs a small, fixed amount of an abstract resource called gas. Before submitting a transaction (or, for a locking/unlocking script pair, before a node bothers validating one), the caller states a **gas limit** — the most gas they are willing to let this particular execution consume. `Execute()` tracks how much gas has been spent so far, instruction by instruction, and the moment spending would exceed the limit, it stops immediately with an error, no matter what the program was in the middle of doing. "Runs forever" becomes "fails safely and predictably, after consuming exactly the resources the caller agreed to risk" — turning an availability catastrophe into an ordinary, expected failure mode.

```
WITHOUT GAS                              WITH GAS

  OpJump -> 0                              OpJump -> 0   (costs, say, 8 gas)
      |                                        |
      v                                        v
  runs forever, this node                  gasUsed += 8, check against
  (and every node validating               gasLimit; once gasUsed > gasLimit,
  the same transaction) hangs               Execute() returns ErrOutOfGas
  indefinitely                              and the node moves on
```

---

## 2. What Gas Actually Measures

Gas is not money, and it is not real computer time measured in milliseconds — it is an abstract unit of "how much work an opcode does," assigned by whoever designs the instruction set, calibrated so that opcodes doing genuinely more work cost genuinely more gas. Three properties make it work:

- **It is deterministic.** Every node computes the identical gas cost for the identical program, because the cost table is a fixed, hardcoded map from opcode to a constant number — never something that depends on the specific machine running it (unlike, say, wall-clock execution time, which Chapter 59 already ruled out as a source of non-determinism for exactly this reason).
- **It is charged before the work happens, not after.** `Execute()` adds an instruction's cost to `gasUsed` and checks the limit *before* (or, as Section 4 shows, in the same step as) actually performing that instruction's effect — so a program can never "sneak in" one more expensive operation by starting it before the gas check catches up.
- **It is cumulative and one-directional.** Gas only ever goes up as execution proceeds; there is no way for a program to "earn back" gas partway through. This is what guarantees termination: since every single instruction costs at least some positive amount (this chapter's table has no zero-cost opcode except the one that stops execution entirely), and `gasLimit` is a fixed finite number chosen before execution starts, there is a hard ceiling on how many instructions *any* program, however it loops or branches, can ever execute.

That last point is the whole safety argument in one sentence: bounded gas per instruction, plus a finite gas limit, mathematically guarantees `Execute()` terminates — regardless of what the program's author intended, or made a mistake about.

---

## 3. Designing a Per-Opcode Cost Table

Real costs should roughly track real work. GoChain VM's opcodes fall into three tiers:

- **Cheap, pure stack manipulation** — pushing a constant, duplicating, discarding, arithmetic, comparisons. These touch a handful of byte slices already sitting in memory and do essentially no computation beyond that.
- **Slightly pricier control flow** — jumps. Not because moving the program counter is expensive by itself, but because unrestricted jumping is exactly the mechanism a runaway loop uses, so pricing it a little above the cheapest tier discourages needlessly loop-heavy programs without being punitive for legitimate `if`/`else` logic or bounded loops.
- **Expensive, "reach outside the stack" operations** — `OpCheckSig` (a real elliptic-curve signature verification, computationally many times more expensive than any stack operation) and `OpSLoad`/`OpSStore` (from Chapter 66 onward, a real disk read or write against persistent contract storage — I/O being far slower than anything happening purely in memory).

```go
package vm

// gasCost assigns a fixed cost to every opcode. Every entry here is a
// deliberate design decision, not an arbitrary number: cheap stack
// operations cost a few gas, control flow costs a bit more (loops are
// exactly what runs away, so a small premium discourages needless
// looping), and anything reaching outside the stack -- cryptography,
// storage -- costs dramatically more, proportional to the real work it
// represents.
var gasCost = map[OpCode]uint64{
	OpPush:        3, // push a constant already sitting in Instruction.Arg
	OpDup:         3, // copy a value already on the stack
	OpPop:         2, // discard a value; strictly less work than a copy
	OpAdd:         3, // one addition over already-decoded operands
	OpSub:         3, // one subtraction, plus the negative-result check
	OpEqual:       3, // one byte-slice comparison
	OpGreaterThan: 3, // one decode-and-compare
	OpJump:        8, // moves the program counter -- the mechanism a
	                  // runaway loop relies on, priced above pure stack ops
	OpJumpIfFalse: 8, // same reasoning as OpJump, plus a pop
	OpCheckSig:    100, // a real ECDSA signature verification -- orders of
	                    // magnitude more CPU work than any stack operation
	OpSLoad:       200, // Chapter 66: a real storage read (disk-backed)
	OpSStore:      200, // Chapter 66: a real storage write (disk-backed)
	OpHalt:        0,   // stops execution; there is no more work left to do
}
```

`OpHalt` is the one deliberate zero: it does not perform any work on the stack at all, it just tells `Execute()` to stop, so charging it anything would be charging for work that never happens. Every opcode that touches the stack, jumps, or reaches outside the VM costs a strictly positive amount — closing the loophole a zero-cost opcode inside a loop body would otherwise open (a `while` loop built entirely out of free instructions would never run out of gas, no matter how large the limit).

Laid out as a quick-reference table, the same costs read like this — useful later, when Chapter 65 needs to estimate a gas limit for a whole token contract before running it:

| Opcode | Gas | Tier |
|---|---|---|
| `OpHalt` | 0 | stops execution |
| `OpPop` | 2 | stack manipulation |
| `OpPush`, `OpDup`, `OpAdd`, `OpSub`, `OpEqual`, `OpGreaterThan` | 3 | stack manipulation / arithmetic / comparison |
| `OpJump`, `OpJumpIfFalse` | 8 | control flow |
| `OpCheckSig` | 100 | cryptography |
| `OpSLoad`, `OpSStore` | 200 | persistent storage (Chapter 66) |

The jump between tiers is deliberately steep — roughly 3x from stack manipulation to control flow, then more than 10x again to reach cryptography, and 2x further still to reach storage — because the *real* costs those operations represent are themselves wildly different orders of magnitude apart, and a gas table that pretended otherwise would either make cryptography artificially cheap (inviting programs that call `OpCheckSig` far more often than the real computational cost justifies) or make ordinary stack arithmetic artificially expensive (making even trivial programs hit their gas limit too easily).

---

## 4. Wiring Costs Into the Execute Loop

Chapter 62's placeholder lived right at the top of the fetch-dispatch loop, before the `switch`:

```go
// Chapter 62's placeholder -- being replaced in this section:
vm.gasUsed++
if vm.gasUsed > vm.gasLimit {
    return ErrOutOfGas
}
```

The real version looks up each instruction's actual cost instead of assuming 1, but the shape of the check — increment, then compare — stays identical:

```go
package vm

import "fmt"

func (vm *VM) Execute() error {
	for {
		if vm.pc < 0 || vm.pc >= len(vm.program) {
			return ErrProgramCounterOutOfRange
		}

		instr := vm.program[vm.pc]

		// Real, per-opcode gas accounting, replacing Chapter 62's flat
		// placeholder. Every opcode must have an entry in gasCost -- an
		// opcode with no cost defined is a bug in the VM itself, not a
		// program error, so it gets its own explicit failure rather than
		// silently charging zero.
		cost, ok := gasCost[instr.Op]
		if !ok {
			return fmt.Errorf("vm: no gas cost defined for opcode %d", instr.Op)
		}
		vm.gasUsed += cost
		if vm.gasUsed > vm.gasLimit {
			return ErrOutOfGas
		}

		switch instr.Op {
		case OpPush:
			vm.execPush(instr)
		case OpDup:
			if err := vm.execDup(); err != nil {
				return err
			}
		case OpPop:
			if err := vm.execPop(); err != nil {
				return err
			}
		case OpAdd:
			if err := vm.execAdd(); err != nil {
				return err
			}
		case OpSub:
			if err := vm.execSub(); err != nil {
				return err
			}
		case OpEqual:
			if err := vm.execEqual(); err != nil {
				return err
			}
		case OpGreaterThan:
			if err := vm.execGreaterThan(); err != nil {
				return err
			}
		case OpJump:
			vm.execJump(instr)
			continue
		case OpJumpIfFalse:
			jumped, err := vm.execJumpIfFalse(instr)
			if err != nil {
				return err
			}
			if jumped {
				continue
			}
		case OpCheckSig:
			if err := vm.execCheckSig(); err != nil {
				return err
			}
		case OpSLoad:
			if err := vm.execSLoad(); err != nil {
				return err
			}
		case OpSStore:
			if err := vm.execSStore(); err != nil {
				return err
			}
		case OpHalt:
			return nil
		default:
			return fmt.Errorf("vm: unknown opcode %d at pc=%d", instr.Op, vm.pc)
		}

		vm.pc++
	}
}
```

Only the gas block changed; every opcode handler from Chapter 62 is untouched. This matters: gas accounting is entirely a property of the *loop*, not of any individual opcode's logic, so adding real costs required editing exactly one part of the file.

One small, useful addition alongside this: an exported way to read back how much gas an execution actually used, which the tests in Sections 6 and 7 rely on.

```go
// GasUsed reports how much gas has been consumed so far. Meaningful to
// call after Execute() returns, whether it returned nil, ErrOutOfGas, or
// any other error -- gasUsed always reflects exactly the work actually
// performed before execution stopped.
func (vm *VM) GasUsed() uint64 {
	return vm.gasUsed
}
```

---

## 5. ErrOutOfGas and Failing Predictably

```go
package vm

import "errors"

// ErrOutOfGas means execution would have exceeded its gas limit. It is
// returned the instant the overrun would happen -- the offending
// instruction's effect on the stack never actually runs, so a program
// that runs out of gas mid-computation never leaves a half-applied
// result behind.
var ErrOutOfGas = errors.New("vm: out of gas")
```

The ordering in Section 4's loop matters here: gas is charged and checked *before* the `switch` that actually performs the instruction's effect. This means an instruction that would push `Execute()` over the limit never runs at all — not even partially. A node that sees `ErrOutOfGas` can treat the whole execution as cleanly rejected, with no need to reason about whatever partial stack or storage state a half-finished instruction might otherwise have left behind.

---

## 6. Testing an Ordinary Program's Gas Cost

Before testing the failure case, confirm gas accounting does not get in the way of normal programs — reusing Chapter 62's full `5 + 3 == 8` program and checking its exact, hand-computable gas total:

```go
package vm

import "testing"

func TestExecute_GasAccounting_NormalProgram(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(5)},  // 3 gas
		{Op: OpPush, Arg: uint64ToBytes(3)},  // 3 gas
		{Op: OpAdd},                          // 3 gas
		{Op: OpPush, Arg: uint64ToBytes(8)},  // 3 gas
		{Op: OpEqual},                        // 3 gas
		{Op: OpHalt},                         // 0 gas
	}
	// Hand-computed total: 3+3+3+3+3+0 = 15.
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.GasUsed() != 15 {
		t.Fatalf("expected exactly 15 gas used, got %d", v.GasUsed())
	}
}

func TestExecute_GasLimitTooLow(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(5)}, // 3 gas -- fits
		{Op: OpPush, Arg: uint64ToBytes(3)}, // 3 gas -- total 6, still fits a limit of 5? No.
		{Op: OpAdd},
		{Op: OpHalt},
	}
	// A limit of 5 lets the first OpPush (3 gas) through, but the second
	// OpPush would bring the total to 6, exceeding the limit.
	v := NewVM(program, 5)
	err := v.Execute()
	if err != ErrOutOfGas {
		t.Fatalf("expected ErrOutOfGas, got %v", err)
	}
	if v.GasUsed() != 6 {
		t.Fatalf("expected gasUsed to reflect the instruction that tipped it over (6), got %d", v.GasUsed())
	}
}
```

```
=== RUN   TestExecute_GasAccounting_NormalProgram
--- PASS: TestExecute_GasAccounting_NormalProgram (0.00s)
=== RUN   TestExecute_GasLimitTooLow
--- PASS: TestExecute_GasLimitTooLow (0.00s)
PASS
ok      github.com/you/gochain/vm    0.003s
```

`TestExecute_GasLimitTooLow` is worth reading closely: `gasUsed` ends at 6, one instruction's worth past the limit of 5, exactly matching Section 5's promise that the check happens before the offending instruction's effect runs — the second `OpPush` never actually adds anything to the stack, but its cost is still the one that is measured against the limit and found to exceed it.

---

## 7. Testing a Deliberate Infinite Loop

This is the test that justifies the whole chapter: a program that, by construction, never reaches `OpHalt` on its own.

```go
package vm

import "testing"

func TestExecute_InfiniteLoop_RunsOutOfGas(t *testing.T) {
	// A program with exactly one instruction: jump to itself. Without gas,
	// Execute() would never return -- this is precisely the "hangs every
	// node forever" scenario Section 1 described, made concrete.
	program := []Instruction{
		{Op: OpJump, Arg: uint64ToBytes(0)}, // jump to index 0 -- itself
	}

	const gasLimit = 1000
	v := NewVM(program, gasLimit)

	err := v.Execute()
	if err != ErrOutOfGas {
		t.Fatalf("expected ErrOutOfGas, got %v", err)
	}

	// OpJump costs 8 gas per iteration (Section 3). 1000 / 8 = 125 exactly,
	// meaning gasUsed reaches exactly 1000 after 125 iterations (not yet
	// over the limit), and the 126th iteration pushes it to 1008 -- the
	// first value that exceeds 1000, which is when Execute() actually stops.
	if v.GasUsed() != 1008 {
		t.Fatalf("expected gasUsed to be 1008 (126 iterations of 8 gas), got %d", v.GasUsed())
	}
}
```

```
=== RUN   TestExecute_InfiniteLoop_RunsOutOfGas
--- PASS: TestExecute_InfiniteLoop_RunsOutOfGas (0.00s)
PASS
ok      github.com/you/gochain/vm    0.021s
```

Nothing about this program's *logic* changed since Chapter 61 first observed that a backward jump is a loop — it is still, semantically, "jump to instruction 0, forever." What changed is the *outcome*: instead of the test (and, in a real node, the entire validation process) hanging forever, `Execute()` returns a clean, ordinary Go error after a bounded, predictable amount of work — exactly 126 iterations, exactly 1008 gas, every single time this exact program runs with this exact gas limit, on any machine, because gas accounting is deterministic (Section 2). The test finishes in milliseconds specifically *because* gas exists; deleting the gas check from Section 4 and re-running this test would hang the test suite indefinitely.

---

## 8. Choosing Gas Limits in Practice

A gas *limit* only helps if something reasonable actually gets chosen for it, and getting this number right involves a real trade-off:

- **Too low**, and legitimate programs — even ones that would have finished perfectly correctly — get rejected with `ErrOutOfGas` before they can complete. A token transfer (Chapter 65) that needs, say, 40 opcodes to run should not be handed a limit of 50 gas, since a single `OpCheckSig` alone costs 100.
- **Too high**, and a single malicious or buggy program can still consume an enormous, if finite, amount of work before failing — the termination guarantee from Section 2 holds regardless, but "guaranteed to eventually stop" and "stops quickly enough to be practical" are different promises.

Real blockchains solve this by making the caller pay for gas with real value (Ethereum's "gas price," paid in the network's native currency), which naturally discourages setting limits far higher than a program actually needs — an idea previewed here and left for Volume 5's fee mechanisms and later chapters to develop further. For now, GoChain simply requires every execution to state an explicit, finite `gasLimit` up front, which is already enough to prove the core safety property this chapter set out to establish: no program, however it is written, can ever force a node to compute forever.

As a concrete example of the estimation this section is describing, take Chapter 63's plain P2PKH spend check: one `OpPush` for the sighash, two more for the unlocking script's signature and public key, one `OpPush` for the engine-computed hash glue, one `OpPush` for the locking script's `PubKeyHash` constant, one `OpEqual`, one `OpJumpIfFalse`, one `OpCheckSig`, and one `OpHalt` on the success path. Using Section 3's table: six pushes at 3 gas each (18), one `OpEqual` at 3, one `OpJumpIfFalse` at 8, one `OpCheckSig` at 100, one `OpHalt` at 0 — a total of exactly 129 gas for a successful spend. A caller offering a `gasLimit` of, say, 1,000 for a plain payment leaves an enormous, deliberate safety margin over the 129 gas an honest spend actually needs — comfortable enough to tolerate the rejected-spend paths from Chapter 63 (which stop earlier, and therefore cost *less* than 129), while still being nowhere near large enough to let a genuinely runaway program do meaningful damage before `ErrOutOfGas` cuts it off.

---

## Summary

- Without gas, a program that jumps backward to itself (or otherwise never reaches `OpHalt`) would hang every node that tries to validate it, forever — an availability catastrophe, not just a failed transaction.
- Gas is a deterministic, cumulative, one-directional cost charged per opcode; combined with a finite `gasLimit`, it mathematically guarantees `Execute()` terminates, regardless of what the program's author intended.
- The cost table prices opcodes by real work: cheap stack operations (`OpPush`, `OpAdd`, `OpEqual`, ...) cost 2-3 gas, control flow (`OpJump`, `OpJumpIfFalse`) costs a bit more at 8 gas since it's the mechanism a runaway loop relies on, and operations reaching outside the stack (`OpCheckSig` at 100, `OpSLoad`/`OpSStore` at 200) cost dramatically more, proportional to real cryptographic or I/O work.
- Gas is charged and checked *before* an instruction's effect runs, so `ErrOutOfGas` always means the offending instruction never partially executed — no half-applied stack or storage state to reason about.
- `Execute()`'s per-instruction placeholder from Chapter 62 was replaced with a `gasCost` map lookup; every other opcode handler is untouched, because gas accounting is a property of the loop, not of any individual opcode.
- A deliberately infinite program (`OpJump` to itself) was tested end to end: it fails predictably with `ErrOutOfGas` after exactly 126 iterations at a 1000 gas limit, instead of hanging the test suite (or a real node) forever.
- Choosing a good `gasLimit` is a genuine trade-off between rejecting legitimate programs (too low) and letting a bad program burn a large-but-finite amount of work before failing (too high); GoChain requires an explicit limit now and previews fee-based limit-setting as a future refinement.

---

## Exercises

### Easy

1. In your own words, explain why gas needs to be deterministic — what would go wrong if two different nodes computed different gas costs for the identical program?
2. Why does `OpHalt` cost 0 gas while every other opcode costs a strictly positive amount? What would happen to the termination guarantee if some *other* opcode were priced at 0?
3. `TestExecute_InfiniteLoop_RunsOutOfGas` expects `gasUsed` to end at exactly 1008, not 1000. Explain, referencing Section 4's ordering of the gas check, why it overshoots the limit rather than stopping exactly at it.

### Medium

4. Compute, by hand, the exact gas cost of this program, instruction by instruction: `OpPush, OpPush, OpCheckSig, OpPush, OpSStore, OpHalt`. Then write a test confirming your hand-computed total matches `GasUsed()`.
5. Modify `TestExecute_InfiniteLoop_RunsOutOfGas`'s gas limit to 10,000 instead of 1,000, recompute by hand exactly how many iterations and what final `gasUsed` value you expect, and update the test's assertions to match.
6. `gasCost` is a `map[OpCode]uint64` looked up on every single instruction. For a very long-running, gas-heavy program, would you expect this map lookup itself to become a performance concern? Propose an alternative data structure (hint: `OpCode` is a `byte`) that would avoid map lookups entirely, and explain the trade-off.

### Hard

7. Design a gas cost for a hypothetical `OpMul` (multiplication) opcode. Should it cost the same as `OpAdd`, more, or less? Justify your answer with reference to how CPUs actually perform multiplication versus addition, and specify the exact number you would assign.
8. `TestExecute_GasLimitTooLow` demonstrates that `Execute()` can exceed `gasLimit` by up to one opcode's full cost before stopping (it never stops mid-instruction). For `OpCheckSig` (100 gas) and `OpSLoad`/`OpSStore` (200 gas), this means the true "worst case" gas consumed can be up to 199 gas over the stated limit. Is this a real safety problem for a node with a fixed budget of, say, 10 million gas per block? Explain your reasoning and propose a design change (if you think one is needed) that would cap the overshoot precisely at the limit instead.
9. Write a test program that deliberately alternates cheap and expensive opcodes (for example, `OpPush`, `OpCheckSig`, `OpPush`, `OpCheckSig`, ...) to reach a specific target `gasUsed` value of your choosing (e.g., exactly 500), using the fewest possible instructions. Show your arithmetic and the resulting instruction slice.
