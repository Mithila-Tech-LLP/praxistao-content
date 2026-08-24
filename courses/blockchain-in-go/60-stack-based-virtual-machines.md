# Chapter 60: Stack-Based Virtual Machines

Chapter 59 established that a smart contract's code must run identically everywhere, which rules out letting contracts execute as ordinary native programs. This chapter introduces the alternative every major blockchain uses: a small, restricted **virtual machine** built around a single stack of values, and we trace a tiny program through it by hand, one instruction at a time, before writing a single line of Go.

## Table of Contents

1. [Why Not Just Run Native Code?](#1-why-not-just-run-native-code)
2. [What a Virtual Machine Is](#2-what-a-virtual-machine-is)
3. [The Stack: One Pile, Two Operations](#3-the-stack-one-pile-two-operations)
4. [A Stack-Based Machine in One Paragraph](#4-a-stack-based-machine-in-one-paragraph)
5. [Tracing a Tiny Program by Hand](#5-tracing-a-tiny-program-by-hand)
6. [Why This Style of Machine Fits Blockchains So Well](#6-why-this-style-of-machine-fits-blockchains-so-well)
7. [Bitcoin Script and the EVM, Briefly](#7-bitcoin-script-and-the-evm-briefly)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Not Just Run Native Code?

The most obvious way to "run a program stored on the blockchain" would be to compile it to native machine code and execute it directly, the way your operating system runs any other program. This is exactly what GoChain will *not* do, for two decisive reasons:

- **Non-determinism.** Native code can do anything the operating system allows: read files, open network connections, check the system clock, spawn threads that finish in unpredictable order. Chapter 59 already established why any one of these breaks the "every node agrees on the result" guarantee a contract needs.
- **Danger.** Native code, run with the full privileges of whatever process executes it, can crash the node, read memory it should not, or (in the worst case) let a malicious contract author execute arbitrary code on every machine in the network. Letting untrusted, internet-supplied code run as unrestricted native instructions is one of the most dangerous things a piece of software can do.

The fix used by every major smart-contract blockchain — Bitcoin, Ethereum, and the many chains modeled after them — is the same: define a small, custom **instruction set** (a fixed, deliberately limited list of operations a program is allowed to express), and run programs written in that instruction set on a **virtual machine** that you, the blockchain's author, control completely. If the instruction set has no "read a file" instruction, no program running on that machine can ever read a file — not because of a rule someone has to remember to follow, but because the operation simply does not exist to be called.

```
NATIVE EXECUTION                          VIRTUAL MACHINE EXECUTION

 contract code                             contract code
      |                                          |
      v                                          v
 compiled to CPU instructions              compiled to a SMALL, FIXED
 (full access to OS, files,                instruction set (OpAdd, OpPush,
  network, memory, time)                   OpJump, ...) with NO access to
      |                                    anything outside a stack and
      v                                    (later) a storage slot.
 runs directly on hardware                       |
                                                  v
 -- dangerous, non-deterministic --        runs inside a Go program (the VM)
                                            that controls every single thing
                                            the code is capable of doing.

                                            -- safe, deterministic by design --
```

---

## 2. What a Virtual Machine Is

A **virtual machine** (VM), in this context, is a program that simulates a computer — it has its own tiny "CPU" logic, its own way of storing values, and its own instruction set — implemented entirely in software, running as an ordinary part of your Go program. It is "virtual" because there is no physical chip anywhere shaped like this machine; GoChain's `vm.VM` is just a Go struct with an `Execute()` method, and the "instructions" it runs are just data (a slice of `Instruction` values) that the VM's Go code interprets one at a time.

This might sound like it adds a layer of indirection for no reason, but that indirection is the entire safety mechanism. Every operation the VM performs is a case in a Go `switch` statement that you wrote and can review completely. There is no way for contract code to do anything the switch statement does not explicitly implement — unlike native code, which can do anything the underlying hardware and operating system allow.

---

## 3. The Stack: One Pile, Two Operations

A **stack** is a simple data structure that holds a sequence of values with exactly two operations: **push** (add a value to the top) and **pop** (remove and return the value currently on top). It follows LIFO order — Last In, First Out — meaning whatever was added most recently is the first thing to come back out.

The everyday analogy: a stack of plates in a kitchen cupboard. You always add a new clean plate to the *top* of the stack, and you always take a plate off the *top* to use it. You cannot politely pull a plate out of the middle of the stack without disturbing everything above it — a real stack data structure enforces this same restriction in code.

```
PUSH 7                     PUSH 3                     POP -> returns 3

  |     |                    |     |                    |     |
  |     |                    |  3  |  <- top              |     |
  |  7  |  <- top            |  7  |                    |  7  |  <- top
  +-----+                    +-----+                    +-----+
  (stack: [7])               (stack: [7, 3])            (stack: [7])
```

A **stack-based virtual machine** keeps exactly one such stack of values as its working memory. Every instruction either pushes something new onto it, pops one or more values off it (to use as inputs), or both — pop some inputs, push back a result. There are no named variables, no registers to keep track of; the stack itself is the only place values live while a program runs. This radical simplicity is a large part of why stack machines are such a good fit for a blockchain: there is very little state to define, serialize, and reason about, which keeps the "what exactly does this program do" question answerable by a human reviewer.

---

## 4. A Stack-Based Machine in One Paragraph

Here is the entire idea, compressed: a program is a list of instructions. Each instruction is a small operation — like "push the number 5," "add the top two numbers," or "jump to instruction 10 if the top of the stack is zero." A **program counter** keeps track of which instruction runs next (starting at 0, normally advancing by one after each instruction — except for jump instructions, which can move it anywhere). The machine repeats one loop: look at the instruction the program counter points to, perform whatever it says to the stack (and possibly the program counter itself), advance, and repeat — until it hits a "stop" instruction, runs out of instructions, or hits an error.

```
                +------------------------------+
                |         PROGRAM               |
                |  [0] OpPush 5                 |
                |  [1] OpPush 3                 |
                |  [2] OpAdd                     |
                |  [3] OpHalt                    |
                +------------------------------+
                              |
                    program counter (pc)
                    points at instruction [pc]
                              |
                              v
        +---------------------------------------------+
        |  loop:                                        |
        |    instr = program[pc]                         |
        |    perform instr's effect on the STACK          |
        |    (and/or move pc directly, for jumps)        |
        |    pc = pc + 1   (unless a jump already moved it)|
        |    repeat until OpHalt, end of program, or error|
        +---------------------------------------------+
```

This is genuinely the whole model. Every opcode GoChain's VM will support — arithmetic, comparisons, control flow, even blockchain-specific operations like verifying a signature — is just a different rule for "what does this instruction do to the stack (and maybe the program counter)."

---

## 5. Tracing a Tiny Program by Hand

The best way to build real intuition for a stack machine is to run one in your head, instruction by instruction, before any Go code exists to do it for you. Consider this tiny program, written informally (Chapter 61 formalizes the exact opcode names and encoding):

```
[0] PUSH 5
[1] PUSH 3
[2] ADD
[3] PUSH 8
[4] EQUAL
[5] HALT
```

In plain language, this program computes `5 + 3` and checks whether the result equals `8`. Let's trace the stack's exact contents after every single instruction:

```
START                      pc=0, stack: [ ]

--- [0] PUSH 5 ---         pc=1, stack: [ 5 ]
    Push the constant 5 onto the stack. Nothing else happens.

      +---+
      | 5 |  <- top
      +---+

--- [1] PUSH 3 ---         pc=2, stack: [ 5, 3 ]
    Push the constant 3. The 5 from before is still underneath it.

      +---+
      | 3 |  <- top
      +---+
      | 5 |
      +---+

--- [2] ADD ---             pc=3, stack: [ 8 ]
    Pop 3 (top), pop 5 (next), push their sum: 5 + 3 = 8.
    Both inputs are GONE from the stack; only the result remains.

      +---+
      | 8 |  <- top
      +---+

--- [3] PUSH 8 ---          pc=4, stack: [ 8, 8 ]
    Push the constant 8. Now the stack holds our computed
    result (8) and a fresh constant (8) to compare it against.

      +---+
      | 8 |  <- top
      +---+
      | 8 |
      +---+

--- [4] EQUAL ---           pc=5, stack: [ 1 ]
    Pop 8 (top), pop 8 (next), compare them: they are equal,
    so push 1 (true). If they had NOT been equal, this would
    push 0 (false) instead.

      +---+
      | 1 |  <- top   (1 means "true": 5 + 3 really does equal 8)
      +---+

--- [5] HALT ---            EXECUTION STOPS
    The machine sees the halt instruction and stops immediately.
    Whatever remains on the stack ( [1] ) is the program's final
    result — in this case, confirmation that the equality check
    passed.
```

Notice exactly what happened at every step: an instruction either added something new to the stack, or consumed some number of values from the top and replaced them with a result. Nothing about this trace required a computer — you just did, by hand, precisely what `vm.VM.Execute()` will do in Go starting in Chapter 62. That is the entire point of a stack machine: its behavior is simple enough to trace on paper, which makes it auditable, which is exactly what a machine running untrusted, financially consequential code needs to be.

---

## 6. Why This Style of Machine Fits Blockchains So Well

A stack-based design is not the only possible choice — some real systems use register-based virtual machines instead, where values live in a fixed set of named slots rather than a single stack. But stack machines have specific properties that make them an especially good fit for blockchains:

- **Tiny, easy-to-specify state.** A stack machine's entire working state, at any instant, is: the stack's contents, the program counter, and (starting in Chapter 64) how much gas has been used. That is a short, precise list — easy to reason about, easy to test exhaustively, easy for a security reviewer to fully understand.
- **Instructions are naturally simple.** Because every instruction only ever touches the top of one stack, there is no way to write an instruction that reaches into arbitrary memory or corrupts unrelated state. Compare this to a native program, where a single bug can corrupt memory anywhere in the process.
- **Trivial to serialize.** A program is just a list of instructions (opcode + optional argument), and the stack is just a list of byte slices. Both are exactly the kind of data GoChain already knows how to hash, sign, and store (Volumes 2, 3, and 8) — no new serialization machinery is needed.
- **Easy to meter.** Because each instruction does one small, bounded unit of work, it is straightforward to assign each one a fixed cost and add them up as execution proceeds — exactly what Chapter 64's gas system depends on.

---

## 7. Bitcoin Script and the EVM, Briefly

Two real systems make this concrete, and GoChain's VM (built starting in Chapter 61) deliberately sits conceptually between them:

- **Bitcoin Script** is a genuinely tiny stack machine: a handful of opcodes, deliberately *not* Turing-complete (it has no general-purpose loop construct), used almost entirely to express "who is allowed to spend this output." Chapter 63's locking/unlocking scripts are directly modeled on this idea.
- **The Ethereum Virtual Machine (EVM)** is a much larger stack machine — hundreds of opcodes, a full general-purpose instruction set (including loops), its own persistent storage per contract, and a gas system to keep execution bounded despite that generality. GoChain's Chapter 61 opcode set, Chapter 64 gas system, and Chapter 66 contract storage are all small, learnable versions of ideas the EVM popularized.

GoChain's own `vm.VM` will land in between: general enough to express real conditional logic and loops (unlike Bitcoin Script), but with a deliberately small opcode table (unlike the EVM's much larger one) — enough to build a real token contract in Chapter 65-69 without needing hundreds of pages of opcode reference.

---

## Summary

- Running contract code as native machine code is unsafe and non-deterministic; blockchains instead run a small, custom virtual machine with a deliberately restricted instruction set.
- A virtual machine, here, is just a Go program (`vm.VM`) that simulates a tiny computer: it has its own "instructions," its own stack, and an execution loop you fully control.
- A stack is a Last-In-First-Out structure with two operations, push and pop; a stack-based VM keeps exactly one stack as its only working memory.
- A stack machine's core loop is simple: fetch the instruction at the program counter, apply its effect to the stack (and possibly the program counter), advance, repeat until halt.
- Tracing a tiny `PUSH 5, PUSH 3, ADD, PUSH 8, EQUAL, HALT` program by hand shows exactly how each instruction consumes and produces stack values — the same mechanics `vm.VM.Execute()` implements in Go from Chapter 62 onward.
- Stack machines fit blockchains especially well because their state is small and easy to specify, their instructions cannot corrupt unrelated memory, they serialize trivially, and they are easy to meter with gas.
- Bitcoin Script (tiny, not Turing-complete) and the EVM (large, general-purpose, gas-metered) bracket the design space; GoChain's VM sits deliberately in between — small but general enough for a real token contract.

---

## Exercises

### Easy

1. In your own words, explain the two reasons a blockchain does not run contract code as ordinary native machine code.
2. Define, in one or two sentences each: virtual machine, stack, program counter, push, pop.
3. What does LIFO stand for, and how does the kitchen-plates analogy illustrate it?

### Medium

4. Trace this program by hand, the way Section 5 did, showing the stack's exact contents after every instruction: `PUSH 10, PUSH 4, SUB, PUSH 6, EQUAL, HALT`. What is the final value left on the stack?
5. Explain why a stack machine's instructions "cannot corrupt unrelated memory," contrasting this with how a bug in a native program can corrupt arbitrary memory elsewhere in the process.
6. A friend suggests GoChain should just let contracts be written in ordinary Go and run directly, "since Go is a safe, memory-safe language anyway." Explain what problem this would still fail to solve, even though Go is indeed memory-safe.

### Hard

7. Trace this program, which uses a comparison result to decide what happens next (informally written; formal jump semantics arrive in Chapter 61): `PUSH 5, PUSH 5, EQUAL` followed by "if the top of the stack is 0, push the number 100; otherwise push the number 200." Show the stack after each step and state the final answer.
8. Research, at a high level, one real difference between Bitcoin Script and the EVM beyond "small vs. large" (for example: Bitcoin Script's lack of general loops, or how each one handles persistent storage). Explain in your own words why that specific difference exists, given what each system is trying to accomplish.
9. Design a tiny stack-machine program (5-8 instructions, using only the informal `PUSH`, `ADD`, `SUB`, `EQUAL` operations from this chapter) that computes whether a number pushed at the start is exactly double another number pushed at the start. Trace it by hand and show your work.
