# Chapter 61: Designing an Instruction Set

Chapter 60 traced a program through an informal stack machine by hand. This chapter makes it exact: we design GoChain VM's real opcode table, specify precisely what each opcode pops and pushes, and settle the encoding conventions (how numbers and booleans are represented as bytes) that every later chapter in this volume — and the token contract built in the second half of this volume — depends on getting right.

## Table of Contents

1. [What an Instruction Set Actually Needs to Express](#1-what-an-instruction-set-actually-needs-to-express)
2. [Encoding Conventions: Numbers and Booleans as Bytes](#2-encoding-conventions-numbers-and-booleans-as-bytes)
3. [Stack Effect Notation](#3-stack-effect-notation)
4. [Arithmetic Opcodes](#4-arithmetic-opcodes)
5. [Stack Manipulation Opcodes](#5-stack-manipulation-opcodes)
6. [Comparison Opcodes](#6-comparison-opcodes)
7. [Control Flow Opcodes](#7-control-flow-opcodes)
8. [Blockchain-Specific Opcodes](#8-blockchain-specific-opcodes)
9. [The Complete Opcode Table](#9-the-complete-opcode-table)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What an Instruction Set Actually Needs to Express

Before listing individual opcodes, it helps to name the *categories* of work any useful contract logic needs, because GoChain VM's opcode table is organized around exactly these categories:

- **Arithmetic** — combine numbers (add, subtract) to compute amounts, balances, totals.
- **Stack manipulation** — duplicate or discard values, so a program can reuse a value without recomputing it, or clean up values it no longer needs.
- **Comparison** — ask yes/no questions about two values (are they equal? is one bigger?), producing a boolean result other instructions can act on.
- **Control flow** — change *which* instruction runs next, based on a computed condition, so a program can express "if" logic and loops instead of always running top to bottom.
- **Blockchain-specific operations** — things a general-purpose calculator has no reason to support, but a blockchain contract absolutely needs: verifying a cryptographic signature from inside the running program, and (Chapter 66) reading and writing the contract's own persistent storage.

Every opcode GoChain VM supports falls into exactly one of these five buckets. This chapter designs all of them; Chapter 62 implements every one in Go.

---

## 2. Encoding Conventions: Numbers and Booleans as Bytes

Every value on GoChain VM's stack is a `[]byte` — a raw slice of bytes, with no built-in notion of "this one is a number" versus "this one is text." This mirrors the `Instruction.Arg []byte` field from the shared `vm` package shape, and it keeps the VM itself completely generic: it does not need a type system, because every opcode's handler decides for itself how to interpret the bytes it receives.

Two conventions, fixed now so every later chapter agrees on them:

- **Numbers** are encoded as big-endian, unsigned 64-bit integers using Go's `encoding/binary` package — the same 8-byte representation used elsewhere. `OpAdd`, `OpSub`, and `OpGreaterThan` all interpret their operands this way. (A production VM would support arbitrary-precision integers; GoChain VM uses `uint64` to keep the implementation approachable, and this limitation is called out explicitly rather than hidden.)
- **Booleans** are single bytes: `[]byte{0x01}` means true, `[]byte{0x00}` means false. A value is considered "truthy" for `OpJumpIfFalse` if it is non-empty and contains at least one non-zero byte; anything else (including an empty slice) is "falsy." `OpEqual`, `OpGreaterThan`, and `OpCheckSig` all push exactly one of these two canonical values as their result.

```
NUMBER 5 as bytes (big-endian uint64):
  [0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x05]

BOOLEAN true:   [0x01]
BOOLEAN false:  [0x00]
```

---

## 3. Stack Effect Notation

Each opcode below is specified with a **stack effect**: a compact notation showing what it consumes from the stack (read left to right, top of stack listed last, immediately before the arrow) and what it leaves behind. For example:

```
OpAdd:  ( a b -- sum )
```

reads as: "given a stack with `a` below `b` on top, `OpAdd` consumes both and leaves `sum` in their place." Anything not mentioned (the rest of the stack below `a`) is left completely untouched. This notation comes from the Forth programming language tradition, which is itself a stack-based language — it is the standard, precise way stack machine opcodes are documented, and we use it for every opcode from here on.

---

## 4. Arithmetic Opcodes

| Opcode | Stack Effect | Behavior |
|---|---|---|
| `OpAdd` | `( a b -- sum )` | Pop `b` (top), pop `a`, push `a + b` as a big-endian uint64. |
| `OpSub` | `( a b -- diff )` | Pop `b` (top), pop `a`, push `a - b` as a big-endian uint64. |

Both opcodes pop exactly two values and push exactly one: the stack shrinks by one item each time either runs. Subtraction that would go negative is a deliberate error case — GoChain VM uses unsigned integers, so `OpSub` on `3 - 5` cannot represent `-2`; Chapter 62's implementation makes this an execution error rather than silently wrapping around, because a wrapped-around balance is exactly the kind of bug that costs real money in a financial system.

---

## 5. Stack Manipulation Opcodes

| Opcode | Stack Effect | Behavior |
|---|---|---|
| `OpDup` | `( a -- a a )` | Peek the top value without removing it, and push a second copy on top. |
| `OpPop` | `( a -- )` | Pop the top value and discard it. |

`OpDup` is essential the moment a program needs to use the same value twice — for example, checking a balance is above zero *and* subtracting from it, without having to push the same constant twice by hand. `OpPop` is the cleanup opcode: it discards a value a program no longer needs, keeping the stack from accumulating leftover results between logical steps.

```
Before OpDup:        After OpDup:         After OpPop (on the OpDup result):

  +---+                +---+                +---+
  | 5 | <- top          | 5 | <- top (copy)   | 5 | <- top (the original 5)
  +---+                +---+                +---+
                        | 5 |
                        +---+
```

---

## 6. Comparison Opcodes

| Opcode | Stack Effect | Behavior |
|---|---|---|
| `OpEqual` | `( a b -- result )` | Pop `b` (top), pop `a`; push `[]byte{1}` if their raw bytes are identical, else `[]byte{0}`. |
| `OpGreaterThan` | `( a b -- result )` | Pop `b` (top), pop `a`; push `[]byte{1}` if `a > b` (as big-endian uint64), else `[]byte{0}`. |

`OpEqual` compares raw bytes directly (using Go's `bytes.Equal`), which means it works correctly on any value — numbers, booleans, or arbitrary byte strings like addresses — not just numbers. `OpGreaterThan` specifically decodes both operands as `uint64`, so it only makes sense for numeric comparisons. Both opcodes push exactly the canonical boolean encoding from Section 2, so their result can flow directly into `OpJumpIfFalse` without any conversion step.

---

## 7. Control Flow Opcodes

| Opcode | Stack Effect | Behavior |
|---|---|---|
| `OpJump` | `( -- )` | Unconditionally set the program counter to the instruction index encoded in `Instruction.Arg`. Does not touch the stack at all. |
| `OpJumpIfFalse` | `( cond -- )` | Pop `cond` (top). If it is "falsy" (Section 2), set the program counter to the index encoded in `Instruction.Arg`; otherwise advance normally to the next instruction. |

This is the one place the opcode table departs from "everything comes from the stack": the jump *target* comes from the instruction's `Arg` field, not from a popped stack value, exactly matching the shared `Instruction` shape (`Arg []byte // operand for OpPush/OpJump/etc`). Only the *condition* for `OpJumpIfFalse` is popped from the stack. This design mirrors real bytecode formats (including the EVM's `JUMP`/`JUMPI`), where jump targets are baked into the instruction stream by whatever assembled the program, rather than computed at runtime from arbitrary stack values — which keeps `Execute()`'s bounds-checking simple and predictable.

```
[0] OpPush 1                  If the value on top is truthy, keep running
[1] OpJumpIfFalse -> [4]      normally into [2]. If it's falsy, skip straight
[2] OpPush 100                to instruction [4], skipping [2] and [3].
[3] OpJump -> [5]
[4] OpPush 200
[5] OpHalt
```

Together, `OpJump` and `OpJumpIfFalse` are enough to express both `if`/`else` branching and loops (a loop is simply a jump backward to an earlier instruction index) — which is exactly how Chapter 64 constructs its deliberate infinite-loop test.

---

## 8. Blockchain-Specific Opcodes

| Opcode | Stack Effect | Behavior |
|---|---|---|
| `OpCheckSig` | `( data signature pubkey -- result )` | Pop `pubkey` (top), pop `signature`, pop `data`. Push `[]byte{1}` if `crypto.Verify(pubkey, data, signature)` succeeds, else `[]byte{0}`. |
| `OpSLoad` | `( key -- value )` | Pop `key` (top). Push the value currently stored at that key in this contract's storage, or an empty `[]byte{}` if nothing has been stored there yet. (Fully implemented in Chapter 66.) |
| `OpSStore` | `( value key -- )` | Pop `key` (top), pop `value`. Persist `value` at that key in this contract's storage, overwriting whatever was there. (Fully implemented in Chapter 66.) |

`OpCheckSig` is the opcode that lets a contract answer "did the right person authorize this?" from *inside* running code, by calling straight into the `gochain/crypto.Verify` function Volume 2 already built — no new cryptography, just a new place to invoke it from. This is exactly the opcode Chapter 63 uses to implement "pay to public key hash" locking scripts, and the reason a plain payment and a smart contract can share one execution model: a plain payment is just a program whose only instruction is, in effect, `OpCheckSig`.

Notice the pop order for `OpCheckSig`: `pubkey` is popped first (it is the top of the stack), meaning the convention for building a program that calls it is to push `data`, then `signature`, then `pubkey`, in that order — the last thing pushed is the first thing popped. Chapter 63 shows exactly how a locking script and an unlocking script cooperate to leave the stack in this shape by the time `OpCheckSig` runs.

`OpSLoad` and `OpSStore` are previewed here for completeness of the opcode table (they appear starting in Chapter 66, written in this volume's second half), but they follow the same key-on-top convention as everything else: the storage *key* is always the value nearest the top of the stack, whether it is being read (`OpSLoad`) or written (`OpSStore`).

---

## 9. The Complete Opcode Table

Here is the full instruction set, in the exact order the shared `vm.OpCode` enum defines it, gathered into one reference table:

| Opcode | Category | Stack Effect | Uses `Arg`? |
|---|---|---|---|
| `OpPush` | constant | `( -- value )` | Yes — the constant to push |
| `OpDup` | stack manipulation | `( a -- a a )` | No |
| `OpPop` | stack manipulation | `( a -- )` | No |
| `OpAdd` | arithmetic | `( a b -- sum )` | No |
| `OpSub` | arithmetic | `( a b -- diff )` | No |
| `OpEqual` | comparison | `( a b -- result )` | No |
| `OpGreaterThan` | comparison | `( a b -- result )` | No |
| `OpJump` | control flow | `( -- )` | Yes — target instruction index |
| `OpJumpIfFalse` | control flow | `( cond -- )` | Yes — target instruction index |
| `OpCheckSig` | blockchain-specific | `( data signature pubkey -- result )` | No |
| `OpSLoad` | blockchain-specific | `( key -- value )` | No |
| `OpSStore` | blockchain-specific | `( value key -- )` | No |
| `OpHalt` | control flow | `( -- )` | No |

This table is the exact contract Chapter 62's Go implementation must satisfy, opcode for opcode, and it is what Chapter 63's scripting language and the token contract in Chapters 65-69 are both built directly on top of. Any future opcode this course (or a reader extending GoChain on their own) adds should be specified the same way: name, category, stack effect, and whether it reads `Arg`, before a single line of Go is written.

---

## Summary

- GoChain VM's opcodes fall into five categories: arithmetic, stack manipulation, comparison, control flow, and blockchain-specific operations.
- Every stack value is a raw `[]byte`; numbers are big-endian `uint64`, booleans are the canonical `[]byte{1}`/`[]byte{0}`, fixed now so every later chapter agrees.
- Stack effect notation `( inputs -- outputs )`, with the top of stack listed last before the arrow, precisely documents what each opcode pops and pushes.
- `OpAdd`/`OpSub` pop two numbers and push one; `OpDup`/`OpPop` manipulate the stack without arithmetic; `OpEqual`/`OpGreaterThan` pop two values and push a boolean.
- `OpJump` and `OpJumpIfFalse` get their jump target from the instruction's `Arg` field, not the stack — only `OpJumpIfFalse`'s condition comes from a pop — and together they're enough to express `if`/`else` and loops.
- `OpCheckSig` pops `data`, `signature`, `pubkey` (in that order, pubkey on top) and calls straight into `gochain/crypto.Verify`; `OpSLoad`/`OpSStore` (fully built in Chapter 66) always treat the storage key as the value nearest the top.
- The complete opcode table in Section 9 is the exact contract Chapter 62 implements and every later chapter in this volume depends on without modification.

---

## Exercises

### Easy

1. Without looking back at Section 9, write out the stack effect notation for `OpAdd`, `OpDup`, and `OpEqual` from memory, then check your answers.
2. Explain, in your own words, why `OpJump`'s target comes from `Instruction.Arg` rather than being popped off the stack.
3. What does the boolean value `[]byte{0x00}` mean, and which opcode is most directly responsible for deciding whether a value counts as "truthy" or "falsy"?

### Medium

4. `OpSub` on GoChain VM is defined over unsigned 64-bit integers. Explain, using a concrete example, why `PUSH 3, PUSH 5, OpSub` is a problem, and propose two different reasonable ways `Execute()` could handle it (you do not need to pick one — Chapter 62 will).
5. Using the stack effect notation from this chapter, write out (on paper, no Go needed) a short instruction sequence that computes `(7 + 3) > 8` and leaves the boolean result on top of the stack. Show the stack's contents after every instruction, the way Chapter 60's trace did.
6. Explain why `OpCheckSig`'s three operands are popped in the specific order `pubkey`, then `signature`, then `data`, and what that implies about the order a program must push them in beforehand.

### Hard

7. Using only `OpPush`, `OpDup`, `OpGreaterThan`, `OpJumpIfFalse`, and `OpJump`, design (on paper) a small program that computes the larger of two pushed numbers and leaves it on top of the stack. Specify exact instruction indices for your jump targets.
8. `OpEqual` compares raw bytes with `bytes.Equal`, while `OpGreaterThan` decodes both operands as `uint64` first. Explain a scenario where using `OpGreaterThan` on non-numeric data (like an address) would produce a meaningless or misleading result, and why `OpEqual` does not have this problem.
9. Propose a brand-new opcode this instruction set does not have (for example, one that checks the current block height, or multiplies two numbers) and specify it completely: its name, category, exact stack effect notation, whether it needs `Arg`, and one sentence justifying why it would be safe and deterministic to add.
