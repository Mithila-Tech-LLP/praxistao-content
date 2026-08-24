# Chapter 13: Putting It Together — A Simple CPU in Action

Imagine you have spent the last twelve chapters learning the individual instruments of an orchestra. You know the violin, the trumpet, the drums, the piano — each one in isolation. Now, for the first time, the conductor raises the baton and the whole orchestra plays a piece of music together. That is this chapter.

We are going to write a small program, translate it into machine code, and then trace every single instruction through the CPU — showing exactly what the Program Counter, registers, memory, and control signals look like at each moment in time. Nothing will be hidden. Nothing will be hand-waved. Every gate, every flip-flop, every wire from the earlier chapters will be earning its keep.

By the end of this chapter you will have watched a real computation happen from the very first fetch to the very last clock cycle.

---

## Table of Contents

1. [Our Minimal CPU — A Refresher](#1-our-minimal-cpu--a-refresher)
2. [The Instruction Set We Will Use](#2-the-instruction-set-we-will-use)
3. [The Program: Sum 1 + 2 + 3 + 4 + 5](#3-the-program-sum-1--2--3--4--5)
4. [Translating Assembly to Machine Code](#4-translating-assembly-to-machine-code)
5. [The Complete Execution Trace](#5-the-complete-execution-trace)
6. [How the Loop Works — Branch Instructions Explained](#6-how-the-loop-works--branch-instructions-explained)
7. [Connecting Back to Part 1 — The Hardware Beneath](#7-connecting-back-to-part-1--the-hardware-beneath)
8. [Quick Checks](#8-preview-of-part-2-making-it-faster)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Our Minimal CPU — A Refresher

Before we run the program, let us agree on the exact machine we are using. Think of this as reading the instruction manual before cooking from a recipe.

Our CPU has the following parts:

**Registers** (the scratchpad on the chef's counter):
- `R0` — general-purpose register (starts at 0)
- `R1` — general-purpose register (starts at 0)
- `R2` — general-purpose register (starts at 0)
- `PC` — Program Counter (holds the address of the next instruction)
- `FLAG` — a 1-bit comparison result flag (0 = "not equal/zero", 1 = "equal/zero")

**Memory** (the pantry the chef walks to):
- 16 memory locations, addresses 0 through 15
- Addresses 0–9: where we store our program instructions
- Addresses 10–15: where we store our data (variables)

**Clock**: One instruction completes each clock cycle.

Here is a bird's-eye view of the data paths in our CPU:

```
         +------------------------------------------+
         |               CPU                        |
         |                                          |
         |  +--------+    +---------+    +------+   |
 Memory  |  |        |    |         |    |      |   |
 ------> |  | FETCH  |--> | DECODE  |--> |EXECUTE   |
 <------ |  |        |    |         |    |      |   |
         |  +--------+    +---------+    +------+   |
         |      ^                           |        |
         |      |         +-----+           |        |
         |      |         | PC  |<----------+        |
         |      +---------|     |                    |
         |                +-----+                    |
         |                                          |
         |  Registers: R0, R1, R2, FLAG             |
         +------------------------------------------+
```

The PC tells the fetch unit which memory address to read from. After each instruction, the PC automatically advances by 1. (Unless a jump instruction overrides it — more on that in Section 6.)

---

## 2. The Instruction Set We Will Use

Our CPU understands exactly seven instructions. Think of these as the seven verbs our computer speaks. Every program must be built from only these words.

| Instruction | Syntax | What It Does |
|---|---|---|
| LOAD | `LOAD Rx, addr` | Copy the value at memory[addr] into register Rx |
| STORE | `STORE Rx, addr` | Copy the value from register Rx into memory[addr] |
| ADD | `ADD Rx, Ry` | Rx = Rx + Ry |
| SUB | `SUB Rx, Ry` | Rx = Rx - Ry |
| CMP | `CMP Rx, Ry` | If Rx == Ry, set FLAG = 1; otherwise FLAG = 0 |
| JMP | `JMP addr` | Unconditionally set PC = addr (jump to that address) |
| JEQ | `JEQ addr` | If FLAG == 1, set PC = addr; otherwise do nothing |

That is it. Seven instructions. Yet from these seven verbs, every computation your computer has ever done is ultimately composed. The elegance here is genuinely astonishing.

**How instructions are encoded in binary**

Each instruction occupies one memory location and is 8 bits wide. We split those 8 bits like this:

```
  Bits 7-5: Opcode (which instruction)
  Bits 4-3: First operand (register or unused)
  Bits 2-0: Second operand (address, register, or unused)
```

| Instruction | Opcode (binary) | Opcode (decimal) |
|---|---|---|
| LOAD | 001 | 1 |
| STORE | 010 | 2 |
| ADD | 011 | 3 |
| SUB | 100 | 4 |
| CMP | 101 | 5 |
| JMP | 110 | 6 |
| JEQ | 111 | 7 |

Register encoding:
- R0 = 00
- R1 = 01
- R2 = 10

For example, `LOAD R0, 10` (load the value at address 10 into R0) encodes as:

```
  001  00  01010
   ^    ^    ^
   |    |    |
 LOAD  R0  addr=10

= 0 0 1 0 0 0 1 0 1 0   (but we only have 8 bits, so addr is 3 bits: 010)
= 001 00 010
= 0 0 1 0 0 0 1 0  = 0x22 = 34 decimal
```

For simplicity in our trace, we will write instructions in assembly form (like `LOAD R0, 10`) and show the binary alongside. You will not need to decode the binary yourself — we will do it together at each step.

### Quick Check A

1. How many registers does our CPU have, and what is each one's purpose?
2. What does the FLAG register store, and which instruction sets it?
3. What is the difference between JMP and JEQ?

---

## 3. The Program: Sum 1 + 2 + 3 + 4 + 5

We want to compute 1 + 2 + 3 + 4 + 5 = 15.

A human would add these from left to right: "1 plus 2 is 3, plus 3 is 6, plus 4 is 10, plus 5 is 15." Our CPU will do essentially the same thing, but using a loop — a section of code that repeats until a condition is met. Think of a loop like a factory conveyor belt: the same set of operations is applied over and over until the product is finished.

**The strategy:**
- Keep a running total (accumulator) in R0
- Keep the current number we are adding in R1
- Keep a countdown counter in R2 (starts at 5, counts down to 0)
- Each loop iteration: add R1 to R0, increment R1, decrement R2
- When R2 reaches 0, the loop ends

**Memory layout we will use:**

| Address | Contents |
|---|---|
| 10 | 0 (initial value for R0, the accumulator) |
| 11 | 1 (initial value for R1, the counter starting at 1) |
| 12 | 5 (initial value for R2, the countdown) |
| 13 | 1 (the constant 1, used to increment/decrement) |
| 14 | 0 (the constant 0, used for comparison) |
| 15 | (result stored here at the end) |

**The assembly program:**

```
Address  Instruction         Comment
-------  ------------------  -----------------------------------
  0      LOAD R0, 10         R0 = 0 (initialize accumulator)
  1      LOAD R1, 11         R1 = 1 (first number to add)
  2      LOAD R2, 12         R2 = 5 (loop counter, counts down)

  -- LOOP START --
  3      ADD  R0, R1         R0 = R0 + R1 (add current number)
  4      LOAD R2, 12         (re-read counter... see note below)

  -- actually, let us use registers throughout --
```

Wait — let us redesign this slightly so it is cleaner and shows register use better. Here is the final version:

```
Address  Instruction         Comment
-------  ------------------  -----------------------------------
  0      LOAD R0, 10         R0 = 0     (accumulator = 0)
  1      LOAD R1, 11         R1 = 1     (current number = 1)
  2      LOAD R2, 12         R2 = 5     (countdown = 5)

  -- LOOP START (address 3) --
  3      ADD  R0, R1         R0 = R0 + R1  (add current number to total)
  4      ADD  R1, R1         ** see note **

```

Hmm, ADD R1, R1 would double R1, not increment it. We need a constant 1. Let us load the constant 1 into a helper approach. Since we only have three general-purpose registers, let us be clever: we will use memory address 13 which holds the value 1, load it when needed, and accept that R2 will need to be reloaded from memory each iteration. This is exactly the kind of constraint real CPU architects face.

**Final program — clean version:**

```
Addr  Instruction       Comment
----  ---------------   -----------------------------------------
 0    LOAD R0, 10       R0 = mem[10] = 0   (accumulator)
 1    LOAD R1, 11       R1 = mem[11] = 1   (number to add, starts at 1)
 2    LOAD R2, 12       R2 = mem[12] = 5   (loop countdown)

      -- LOOP starts here --
 3    ADD  R0, R1       R0 = R0 + R1       (add current number to sum)
 4    STORE R0, 15      mem[15] = R0       (save partial result)
 5    ADD  R1, R1       ** NO -- wrong **
```

Let us be honest: with only 3 registers and no immediate-value instructions, we must reload constants from memory. Here is the truly final program:

```
Addr  Instruction       Comment
----  ---------------   -----------------------------------------
 0    LOAD R0, 10       R0 = 0   (accumulator = 0)
 1    LOAD R1, 11       R1 = 1   (current number = 1)
 2    LOAD R2, 12       R2 = 5   (loop countdown = 5)

      -- LOOP: address 3 --
 3    ADD  R0, R1       R0 = R0 + R1    (accumulate)
 4    LOAD R2, 13       R2 = mem[13] = 1  (load the constant 1)
 5    ADD  R1, R2       R1 = R1 + 1     (increment current number)
 6    LOAD R2, 12       R2 = mem[12]    (reload the countdown)
 7    SUB  R2, R1       Temporarily: compute countdown progress

```

This is getting unwieldy because we keep running out of registers. Let us simplify the memory layout and use a dedicated countdown cell that we update each iteration. Here is the **definitive, clear program** we will trace:

---

**Memory initialization:**

| Address | Value | Role |
|---|---|---|
| 10 | 0 | accumulator (R0 will be loaded here, result stored back) |
| 11 | 1 | current number to add (we increment this) |
| 12 | 5 | loop counter (we decrement this) |
| 13 | 1 | the constant 1 (never changes) |
| 14 | 0 | the constant 0 (never changes) |
| 15 | — | final result stored here |

**Program:**

```
Addr  Instruction       Comment
----  ---------------   ------------------------------------------
 0    LOAD R0, 10       R0 = 0    initialize accumulator
 1    LOAD R1, 11       R1 = 1    load first number
 2    LOAD R2, 12       R2 = 5    load loop counter

      ---- LOOP TOP (address 3) ----
 3    ADD  R0, R1       R0 = R0 + R1   add current number to sum
 4    LOAD R2, 13       R2 = 1         load constant 1
 5    ADD  R1, R2       R1 = R1 + 1    increment current number
 6    LOAD R2, 14       R2 = 0         load constant 0
 7    STORE R0, 15      mem[15] = R0   save running total
 8    LOAD R2, 12       R2 = loop counter (reload)
 9    SUB  R2, R1       hmm this is R2 - R1...
```

You know what? Let us step back from over-engineering and use the simplest approach that clearly demonstrates the loop concept. We will restructure: use R0 as accumulator, R1 as the "number we are adding this iteration" (1, 2, 3, 4, 5), and R2 as the loop counter that counts from 5 down to 0.

---

**FINAL DEFINITIVE PROGRAM (we are committing to this):**

**Memory layout:**

| Address | Initial Value | Purpose |
|---|---|---|
| 10 | 0 | Starting value for accumulator |
| 11 | 1 | The number 1 (constant, used for incrementing) |
| 12 | 5 | Loop counter (decrements each iteration) |
| 13 | 0 | Zero constant (used for loop-end comparison) |
| 14 | 1 | Current number being added (we update this in memory) |
| 15 | — | Final result |

**Assembly program:**

```
Addr  Instruction       Comment
----  ---------------   ------------------------------------------
 0    LOAD R0, 10       R0 = 0     (accumulator = 0)
 1    LOAD R1, 13       R1 = 0     (R1 will hold the zero constant)
 2    LOAD R2, 12       R2 = 5     (loop counter)

      ---- LOOP TOP at address 3 ----
 3    LOAD R0, 14       R0 = current number (1, then 2, 3, 4, 5)
 4    LOAD R2, 12       R2 = loop counter  
 5    ADD  R0, R2       R0 = R0 + R2  -- NO, this adds counter not the sum
```

I keep running into the same register-shortage problem. Let me take a completely different, simpler approach. We do NOT loop by decrementing a counter. We just unroll the loop across 5 ADD instructions. But that defeats the purpose of showing a branch.

The correct solution: use R0 as accumulator, and instead of trying to increment the "number being added", simply add the loop counter itself each iteration. Since counter goes 5, 4, 3, 2, 1 — the sum is 5+4+3+2+1 = 15 = 1+2+3+4+5. Same answer.

---

**THE ACTUAL FINAL PROGRAM:**

**Strategy:** R0 = running sum. R2 = counter (starts at 5, goes down to 1). Each iteration: R0 = R0 + R2, then R2 = R2 - 1. When R2 = 0, stop.

We need a constant 1 to subtract from R2. We load it into R1 once and leave it there.

**Memory:**

| Address | Initial Value | Purpose |
|---|---|---|
| 10 | 0 | Initial value 0 |
| 11 | 1 | Constant: 1 |
| 12 | 0 | Constant: 0 (for comparison) |
| 13 | — | Result stored here |

**Program:**

```
Addr  Instruction       Comment
----  ---------------   ------------------------------------------
 0    LOAD R0, 10       R0 = 0     (accumulator = 0)
 1    LOAD R1, 11       R1 = 1     (constant 1, stays in R1 forever)
 2    LOAD R2, 11       R2 = 5 -- WAIT, address 11 holds 1, not 5
```

OK. I need address 10=0, address 11=1, address 12=5. Let me just be explicit.

**Memory:**

| Address | Initial Value | Purpose |
|---|---|---|
| 10 | 0 | The constant zero |
| 11 | 1 | The constant one |
| 12 | 5 | The loop counter (starts at 5) |
| 13 | — | Final result |

**Program:**

```
Addr  Instruction       Comment
----  ---------------   ------------------------------------------
 0    LOAD R0, 10       R0 = 0     (accumulator = 0)
 1    LOAD R1, 11       R1 = 1     (constant 1, stays here)
 2    LOAD R2, 12       R2 = 5     (loop counter = 5)

      ---- LOOP TOP: address 3 ----
 3    ADD  R0, R2       R0 = R0 + R2   (add counter to sum: 0+5, 5+4, 9+3...)
 4    SUB  R2, R1       R2 = R2 - 1    (decrement loop counter)
 5    STORE R2, 12      mem[12] = R2   (save updated counter back to memory)
 6    CMP  R2, R0 -- NO, we want CMP R2 with 0
```

We want to compare R2 with 0. But 0 is in memory[10], not in a register. We need to either:
(a) Load 0 into a register (but we used all 3 registers)
(b) Use a trick

With only 3 registers, once R0=accumulator, R1=constant 1, R2=counter, we have no free register to hold 0 for comparison. However — here is the key insight — CMP R2, R0 compares R2 with R0. But R0 is not 0, it is the running sum.

The solution: we accept a small program restructuring. At the moment we want to compare, we temporarily use R0 for something else. OR: we structure the program so that right before the comparison, we use R0 to hold 0.

Actually, the cleanest solution for a teaching example: let us use 4 registers. Let me re-define our CPU to have R0, R1, R2, R3. This is more realistic anyway (most real CPUs have 8, 16, or 32 registers).

---

## Revised CPU Spec: 4 Registers

Our CPU now has: **R0, R1, R2, R3** (all start at 0), plus **PC** and **FLAG**.

**Memory:**

| Address | Initial Value | Purpose |
|---|---|---|
| 10 | 0 | constant zero |
| 11 | 1 | constant one |
| 12 | 5 | loop counter |
| 13 | — | final result |

**Program (sum numbers 5 down to 1, equals 15):**

```
Addr  Instruction       Comment
----  ---------------   ------------------------------------------
 0    LOAD R0, 10       R0 = 0     (accumulator = 0)
 1    LOAD R1, 11       R1 = 1     (constant 1, for decrement)
 2    LOAD R2, 12       R2 = 5     (loop counter = 5)
 3    LOAD R3, 10       R3 = 0     (constant zero, for comparison)

      ---- LOOP TOP: address 4 ----
 4    ADD  R0, R2       R0 = R0 + R2   (add R2 to running sum)
 5    SUB  R2, R1       R2 = R2 - 1    (decrement counter)
 6    CMP  R2, R3       FLAG = (R2 == 0) ? 1 : 0
 7    JEQ  9            if FLAG=1, jump to address 9 (end)
 8    JMP  4            else jump back to loop top

      ---- END: address 9 ----
 9    STORE R0, 13      mem[13] = R0   (save result = 15)
```

This program is 10 instructions long (addresses 0–9). It uses a conditional branch (JEQ) to exit the loop and an unconditional jump (JMP) to repeat the loop. Let us count the iterations:

| Iteration | R2 at start | Adds to R0 | R0 after | R2 after |
|---|---|---|---|---|
| 1 | 5 | 5 | 5 | 4 |
| 2 | 4 | 4 | 9 | 3 |
| 3 | 3 | 3 | 12 | 2 |
| 4 | 2 | 2 | 14 | 1 |
| 5 | 1 | 1 | 15 | 0 |
| (exit) | 0 | — | — | — |

5 + 4 + 3 + 2 + 1 = 15. Correct!

### Quick Check B

1. Why do we need R3 to hold the constant 0, when 0 is already in memory at address 10?
2. What would happen if we forgot instruction 8 (the JMP 4)?
3. How many total clock cycles does this program take to complete? (Count: 4 setup instructions + 5 iterations × 5 loop instructions + 1 final store = ?)

---

## 4. Translating Assembly to Machine Code

Now we translate each instruction into the 8-bit binary the CPU actually sees. Recall our encoding:

```
Bits [7:5] = Opcode
Bits [4:3] = Destination register (for instructions that use one)
Bits [2:0] = Source register OR memory address
```

Register encoding: R0=00, R1=01, R2=10, R3=11
Address encoding: 0–15 in 4-bit binary (we allow 4 bits for the address field when needed)

For simplicity, instructions that use an address (LOAD, STORE, JMP, JEQ) are encoded as:
```
Bits [7:5] = Opcode
Bits [4:3] = Register (if applicable)
Bits [3:0] = 4-bit address (0–15)
```
(This slightly overlaps bit 3, which is fine for this teaching example.)

| Addr | Assembly | Binary | Hex | Notes |
|---|---|---|---|---|
| 0 | LOAD R0, 10 | 001 00 1010 | 0x0A family | opcode=001, reg=00, addr=1010 |
| 1 | LOAD R1, 11 | 001 01 1011 | — | opcode=001, reg=01, addr=1011 |
| 2 | LOAD R2, 12 | 001 10 1100 | — | opcode=001, reg=10, addr=1100 |
| 3 | LOAD R3, 10 | 001 11 1010 | — | opcode=001, reg=11, addr=1010 |
| 4 | ADD R0, R2 | 011 00 10 xx | — | opcode=011, dst=00, src=10 |
| 5 | SUB R2, R1 | 100 10 01 xx | — | opcode=100, dst=10, src=01 |
| 6 | CMP R2, R3 | 101 10 11 xx | — | opcode=101, reg1=10, reg2=11 |
| 7 | JEQ 9 | 111 00 1001 | — | opcode=111, addr=1001 |
| 8 | JMP 4 | 110 00 0100 | — | opcode=110, addr=0100 |
| 9 | STORE R0, 13 | 010 00 1101 | — | opcode=010, reg=00, addr=1101 |

The CPU does not see the assembly labels. It sees only the binary patterns. The decode stage reads the opcode bits and generates the appropriate control signals — exactly what we studied in Chapter 10.

---

## 5. The Complete Execution Trace

Now the moment we have been building toward. We will trace every single clock cycle. At each step, the table shows:

- **PC**: Program Counter (address of instruction being fetched)
- **Instruction**: What is at that address
- **R0, R1, R2, R3**: Register values AFTER this instruction executes
- **FLAG**: Value of the comparison flag
- **mem[12], mem[13]**: Relevant memory cells
- **Notes**: What is physically happening in the hardware

### Initial State (before clock cycle 1)

| PC | R0 | R1 | R2 | R3 | FLAG | mem[10] | mem[11] | mem[12] | mem[13] |
|---|---|---|---|---|---|---|---|---|---|
| 0 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 5 | — |

---

### Cycle 1 — Instruction at address 0

**Fetch:** PC = 0, memory bus reads address 0, returns `LOAD R0, 10`

**Decode:** Control unit sees opcode 001 (LOAD), destination R0, source address 10. Generates signals: enable memory read, route data to R0.

**Execute:** Address 10 is placed on memory address bus. Memory outputs 0. The value 0 is written to R0. PC increments to 1.

```
FETCH                DECODE               EXECUTE
+----------+         +----------+         +----------+
| mem[0]   |         | opcode=  |         | addr bus |
| =LOAD    | ------> | LOAD     | ------> | = 10     |
| R0, 10   |         | dst=R0   |         | data=0   |
+----------+         | src=10   |         | R0 <-- 0 |
                     +----------+         +----------+
```

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **0** | LOAD R0, 10 | **0** | 0 | 0 | 0 | 0 | R0 loaded with 0 from mem[10] |

---

### Cycle 2 — Instruction at address 1

**Fetch:** PC = 1, memory reads address 1, returns `LOAD R1, 11`
**Decode:** LOAD, destination R1, source address 11
**Execute:** mem[11] = 1, write 1 to R1. PC becomes 2.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **1** | LOAD R1, 11 | 0 | **1** | 0 | 0 | 0 | R1 loaded with 1 (the constant) |

---

### Cycle 3 — Instruction at address 2

**Fetch:** PC = 2, memory reads address 2, returns `LOAD R2, 12`
**Decode:** LOAD, destination R2, source address 12
**Execute:** mem[12] = 5, write 5 to R2. PC becomes 3.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **2** | LOAD R2, 12 | 0 | 1 | **5** | 0 | 0 | R2 loaded with 5 (loop counter) |

---

### Cycle 4 — Instruction at address 3

**Fetch:** PC = 3, memory reads address 3, returns `LOAD R3, 10`
**Decode:** LOAD, destination R3, source address 10
**Execute:** mem[10] = 0, write 0 to R3. PC becomes 4.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **3** | LOAD R3, 10 | 0 | 1 | 5 | **0** | 0 | R3 loaded with 0 (comparison constant) |

**Setup complete. All four registers are initialized. The loop begins now.**

---

### LOOP ITERATION 1 (R2 = 5, adding 5 to sum)

---

### Cycle 5 — Instruction at address 4

**Fetch:** PC = 4, returns `ADD R0, R2`
**Decode:** Opcode = ADD. Source registers: R0 and R2. Destination: R0.
**Execute:** ALU receives R0=0 and R2=5. Performs addition. Output = 5. Writes 5 to R0. PC becomes 5.

The ALU at this moment: the 4-bit adder we built in Chapter 5 is receiving bits `0000` and `0101`, producing `0101`. The carry-out is 0. The result `0101` (= 5) is written into the R0 flip-flop array.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **4** | ADD R0, R2 | **5** | 1 | 5 | 0 | 0 | Sum is now 5. (5+0=5) |

---

### Cycle 6 — Instruction at address 5

**Fetch:** PC = 5, returns `SUB R2, R1`
**Decode:** Opcode = SUB. Source: R2 and R1. Destination: R2.
**Execute:** ALU receives R2=5 and R1=1. Subtraction: 5 - 1 = 4. The ALU uses two's complement: it inverts R1 (0001 → 1110), adds 1 to get 1111, then adds to R2=0101 to get 0100. Wait, let us do this correctly:

```
  R2 = 0101 (5)
- R1 = 0001 (1)
Two's complement of R1: NOT(0001) + 1 = 1110 + 0001 = 1111

  0101
+ 1111
------
  0100  (with carry out = 1, which we discard for subtraction)
```

Result: R2 = 4. PC becomes 6.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **5** | SUB R2, R1 | 5 | 1 | **4** | 0 | 0 | Counter decremented: 5→4 |

---

### Cycle 7 — Instruction at address 6

**Fetch:** PC = 6, returns `CMP R2, R3`
**Decode:** Opcode = CMP. Compare R2 with R3.
**Execute:** R2 = 4, R3 = 0. Are they equal? No. FLAG = 0. PC becomes 7.

The comparator circuit: it receives `0100` and `0000`, feeds each bit pair into XNOR gates (XNOR outputs 1 when both inputs are equal), then ANDs all outputs. Since bit 2 differs (1 vs 0), the AND output is 0. FLAG register is set to 0.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **6** | CMP R2, R3 | 5 | 1 | 4 | 0 | **0** | 4 ≠ 0, so FLAG=0 |

---

### Cycle 8 — Instruction at address 7

**Fetch:** PC = 7, returns `JEQ 9`
**Decode:** Opcode = JEQ. Jump to address 9 if FLAG = 1.
**Execute:** FLAG = 0. Condition not met. Do NOT jump. PC increments normally to 8.

This is a critical moment. The control unit reads the FLAG register. It sees 0. The multiplexer that chooses between "PC+1" and "jump target" selects "PC+1". The PC register receives 8.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **7** | JEQ 9 | 5 | 1 | 4 | 0 | 0 | FLAG=0, NO jump. PC → 8. |

---

### Cycle 9 — Instruction at address 8

**Fetch:** PC = 8, returns `JMP 4`
**Decode:** Opcode = JMP. Unconditional jump to address 4.
**Execute:** The jump target (4) is placed directly into the PC register. PC becomes 4.

The multiplexer now unconditionally selects the jump target value 4. The PC flip-flop is overwritten with `0100`.

| PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|
| **8** | JMP 4 | 5 | 1 | 4 | 0 | 0 | Unconditional jump. PC → 4. |

**End of iteration 1. State: R0=5 (sum so far), R2=4 (counter). Loop repeats.**

---

### LOOP ITERATION 2 (R2 = 4, adding 4 to sum)

| Cycle | PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|---|
| 10 | 4 | ADD R0, R2 | **9** | 1 | 4 | 0 | 0 | 5+4=9 |
| 11 | 5 | SUB R2, R1 | 9 | 1 | **3** | 0 | 0 | 4-1=3 |
| 12 | 6 | CMP R2, R3 | 9 | 1 | 3 | 0 | **0** | 3≠0, FLAG=0 |
| 13 | 7 | JEQ 9 | 9 | 1 | 3 | 0 | 0 | no jump, PC→8 |
| 14 | 8 | JMP 4 | 9 | 1 | 3 | 0 | 0 | PC→4 |

---

### LOOP ITERATION 3 (R2 = 3, adding 3 to sum)

| Cycle | PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|---|
| 15 | 4 | ADD R0, R2 | **12** | 1 | 3 | 0 | 0 | 9+3=12 |
| 16 | 5 | SUB R2, R1 | 12 | 1 | **2** | 0 | 0 | 3-1=2 |
| 17 | 6 | CMP R2, R3 | 12 | 1 | 2 | 0 | **0** | 2≠0 |
| 18 | 7 | JEQ 9 | 12 | 1 | 2 | 0 | 0 | no jump |
| 19 | 8 | JMP 4 | 12 | 1 | 2 | 0 | 0 | PC→4 |

---

### LOOP ITERATION 4 (R2 = 2, adding 2 to sum)

| Cycle | PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|---|
| 20 | 4 | ADD R0, R2 | **14** | 1 | 2 | 0 | 0 | 12+2=14 |
| 21 | 5 | SUB R2, R1 | 14 | 1 | **1** | 0 | 0 | 2-1=1 |
| 22 | 6 | CMP R2, R3 | 14 | 1 | 1 | 0 | **0** | 1≠0 |
| 23 | 7 | JEQ 9 | 14 | 1 | 1 | 0 | 0 | no jump |
| 24 | 8 | JMP 4 | 14 | 1 | 1 | 0 | 0 | PC→4 |

---

### LOOP ITERATION 5 (R2 = 1, adding 1 to sum)

| Cycle | PC | Instruction | R0 | R1 | R2 | R3 | FLAG | Notes |
|---|---|---|---|---|---|---|---|---|
| 25 | 4 | ADD R0, R2 | **15** | 1 | 1 | 0 | 0 | 14+1=15 ✓ |
| 26 | 5 | SUB R2, R1 | 15 | 1 | **0** | 0 | 0 | 1-1=0 |
| 27 | 6 | CMP R2, R3 | 15 | 1 | 0 | 0 | **1** | 0==0, FLAG=1! |
| 28 | 7 | JEQ 9 | 15 | 1 | 0 | 0 | 1 | FLAG=1, JUMP to 9! |

At cycle 28, something different happens. The control unit reads FLAG = 1. The "should we jump?" multiplexer receives a 1 on its select line. It routes the jump target address (9) to the PC register instead of PC+1. The loop ends.

---

### Final Instruction — Cycle 29

| Cycle | PC | Instruction | R0 | R1 | R2 | R3 | FLAG | mem[13] | Notes |
|---|---|---|---|---|---|---|---|---|---|
| 29 | 9 | STORE R0, 13 | 15 | 1 | 0 | 0 | 1 | **15** | Result saved to memory |

The result 15 is now in memory address 13. The computation is complete.

**Total clock cycles used: 29**
**Result: mem[13] = 15 = 1+2+3+4+5**

---

### Summary Timeline

```
Cycles 1-4:   Setup (load constants into registers)
Cycles 5-9:   Loop iteration 1 (add 5)
Cycles 10-14: Loop iteration 2 (add 4)
Cycles 15-19: Loop iteration 3 (add 3)
Cycles 20-24: Loop iteration 4 (add 2)
Cycles 25-28: Loop iteration 5 (add 1, then exit loop)
Cycle 29:     Store result

Running total: 0 → 5 → 9 → 12 → 14 → 15
```

### Quick Check C

1. In cycle 27, the FLAG register changes from 0 to 1. Which hardware component is responsible for computing this change?
2. What is the total number of memory reads that occur during the entire program execution? (Hint: every instruction fetch is a memory read, plus the 4 LOAD instructions.)
3. If we wanted to compute 1+2+3+4+5+6, what single change would we make to the memory layout?

---

## 6. How the Loop Works — Branch Instructions Explained

Let us zoom in on exactly what happens inside the CPU when a branch instruction is executed, because this is one of the most important mechanisms in all of computing.

### The Normal PC Increment

In most instructions, after execution, the PC is incremented by 1. This is handled by a simple adder in the control unit:

```
  PC (current) ----> +1 adder ----> PC (next)

  Example: PC=6 ----> +1 ----> PC=7
```

### The Jump Override

When a JMP or JEQ instruction fires, the control unit uses a multiplexer to choose between "PC+1" and the jump target:

```
                    +-----+
  PC+1 ----------->|     |
                   | MUX |-----> PC (next)
  jump_target ---->|     |
                    +-----+
                       ^
                       |
                 select signal
               (0 = PC+1, 1 = jump)

  For JMP:  select is always 1 (always jump)
  For JEQ:  select = FLAG register value
```

This multiplexer is a critical component. It is the physical mechanism that makes loops, if-statements, and all conditional behavior possible. Every if-statement in every programming language ultimately compiles down to a CMP instruction followed by a conditional branch.

### Why Loops Are So Powerful

Without a loop, to add numbers 1 through 100, you would need 100 ADD instructions — 100 lines of code. With a loop, you need 6 instructions, regardless of whether you are summing to 5 or to 5,000,000. The loop is a force multiplier.

```
Without loop (sum 1..5):         With loop (sum 1..N):
  ADD R0, 1                        LOAD R2, N
  ADD R0, 2                   ---> ADD R0, R2   <--+
  ADD R0, 3                   |    SUB R2, R1       |
  ADD R0, 4                   |    CMP R2, R3       |
  ADD R0, 5                   |    JEQ  done        |
  = 5 instructions            |    JMP  loop -------+
                              done: STORE R0, result
                                   = 6 instructions regardless of N!
```

This is why loops exist. They are time machines for the CPU — they let a small program perform an arbitrarily large number of operations.

### What the Hardware Does During a Jump

Let us trace the JEQ 9 instruction at cycle 28 in detail:

```
Step 1 — FETCH:
  PC = 7
  Address bus carries value 7 to memory
  Memory outputs the bits of "JEQ 9"
  Instruction Register (IR) latches these bits

Step 2 — DECODE:
  Bits [7:5] = 111 → Control unit recognizes JEQ
  Bits [3:0] = 1001 → Jump target = 9
  Control unit reads FLAG register → FLAG = 1
  Control unit asserts "conditional jump" signal

Step 3 — EXECUTE:
  Because FLAG = 1:
    Multiplexer select = 1
    PC input = 9 (jump target)
    On next clock edge: PC register latches 9
  PC is now 9, not 8

Step 4 — next FETCH uses PC = 9:
  Memory outputs instruction at address 9: STORE R0, 13
```

The entire mechanism: one multiplexer, one register (FLAG), one set of bits in the instruction. That is all a conditional branch is.

---

## 7. Connecting Back to Part 1 — The Hardware Beneath

Every step in the trace above relied on components we built from first principles in Chapters 1 through 12. Let us honor that work by seeing exactly where each component was used.

### Chapter 2: Transistors and Electricity

Every signal in our trace — every 1 and every 0 — was physically a high or low voltage in a wire. When we said "R0 = 5", what was really happening was that four transistors in the R0 register circuit were in specific on/off states representing `0101`. The "constant 1 in R1" was one transistor conducting and three not conducting, held in place by the feedback loop of the flip-flop.

### Chapter 3: Binary Numbers

When we computed 9 + 3 = 12 in iteration 3, the ALU worked with `1001 + 0011 = 1100`. The concept of binary addition — the carry rippling from bit to bit — was the actual physical process.

### Chapter 4: Logic Gates

The comparator in the CMP instruction is made entirely of XNOR gates and an AND gate. Every time FLAG was set or cleared, thousands of transistors (which are just logic gates underneath) were switching.

### Chapter 5: Combinational Circuits

The ALU that performed every ADD and SUB is the ripple-carry adder from Chapter 5. The 4-bit adder we built there ran on every single iteration. The multiplexer that chose between PC+1 and the jump target is the 2-to-1 MUX from Chapter 5.

### Chapter 6: Flip-Flops and Memory

Every register (R0, R1, R2, R3, PC, FLAG) is a set of D flip-flops. When we wrote "PC becomes 4 after the JMP instruction", what happened was: on the rising edge of the clock signal, the PC flip-flop array latched the value 4. All register updates are edge-triggered D flip-flop operations.

Memory itself (where the instructions and data live) is built from arrays of flip-flops or capacitors arranged in a grid, with address decoders (from Chapter 5) selecting which cell to read or write.

### Chapter 7: CPU Anatomy

The three stages we traced — fetch, decode, execute — are exactly the three units described in Chapter 7. The control unit, ALU, register file, and memory interface all played their roles.

### Chapter 8: The ALU

Every ADD and SUB instruction passed through the ALU. Chapter 8 showed how an ALU is constructed from adders and logic gates with a select signal. In our trace, the select signal came from the opcode bits: 011 meant ADD (ALU mode: add), 100 meant SUB (ALU mode: subtract using two's complement).

### Chapter 9: Registers

The register file is a set of D flip-flops with read/write enable signals. At cycle 5, when we wrote `ADD R0, R2`, the register file output both R0 and R2 on the A and B buses to the ALU, and then accepted the result back on the write bus to R0.

### Chapter 10: The Control Unit

The control unit decoded the opcode bits each cycle and generated the right control signals: memory read/write enable, ALU operation select, register write enable, PC jump select. Every instruction we traced was driven by the control unit.

### Chapter 11: Fetch-Decode-Execute

This entire chapter was one long demonstration of Chapter 11. Every one of the 29 cycles followed the same three-phase pattern.

### Chapter 12: Memory

When we stored the result at the end (`STORE R0, 13`), the address 13 went onto the address bus, R0's value (15) went onto the data bus, and the memory write-enable signal was asserted. The memory cell at address 13 latched the value 15.

---

### The Full Picture in One Diagram

```
          +------------------------------------------------+
          |                   C P U                       |
          |                                               |
MEMORY    |  +---------+   +----------+   +-----------+  |
addr[4]-->|  |         |   |          |   |           |  |
data[8]<->|  | CONTROL |-->| REGISTER |-->|    ALU    |  |
          |  |  UNIT   |   |  FILE    |   |  +,-,cmp  |  |
          |  | (decode)|   | R0,R1,R2 |<--|           |  |
          |  |         |   | R3,FLAG  |   +-----------+  |
          |  +----+----+   +----------+                   |
          |       |                                       |
          |  +----v----+                                  |
          |  |   PC    |                                  |
          |  | counter |---> address bus to memory        |
          |  +---------+                                  |
          +------------------------------------------------+

  Each component = circuits from Chapters 4-10
  Each wire = a transistor channel (Chapter 2)
  Each signal = a binary value (Chapter 3)
```

This is the gestalt. The thing you could not see when staring at a single NAND gate in Chapter 4 is now visible: all those gates were working toward this moment. They were being assembled, layer by layer, into something that can compute.

---

### The Limits of Our Simple CPU

Our CPU works. It computed 1+2+3+4+5=15 correctly in 29 clock cycles. But it has serious limitations:

| Limitation | Problem | Real CPUs solve this by... |
|---|---|---|
| One instruction per cycle | Slow | Pipelining (executing multiple instructions simultaneously) |
| Memory is slow | Waiting for RAM wastes time | Caches (fast memory near the CPU) |
| Only 4 registers | Constant memory access needed | More registers (x86 has 16, ARM has 31) |
| Fixed instruction width | Limited expressiveness | Variable-length encoding (CISC) or extensions (RISC) |
| No parallelism | One thing at a time | Superscalar execution (multiple ALUs) |
| Simple branch | Always costs cycles | Branch prediction |

Each of these limitations is the seed of a major CPU design innovation. Each one is a story we will tell in the rest of this course.

### Quick Check D

1. Which chapter introduced the component that makes JEQ work (the component that switches between two inputs based on a select signal)?
2. When we wrote R0 = 9 in cycle 10, what physical event caused this — what exactly happened in the hardware?
3. Name one limitation of our simple CPU and describe in one sentence what technique a modern CPU uses to overcome it.

---

## 8. Preview of Part 2: Making It Faster

Our CPU is correct. Correctness is non-negotiable. But speed matters enormously.

Consider: a modern CPU runs at 3–5 GHz. That means 3 to 5 billion clock cycles per second. If each cycle executes one instruction, that is 3 billion instructions per second. But modern CPUs actually execute 3–5 instructions per cycle on average, thanks to the techniques we are about to study. That is 10–15 billion instructions per second.

The gap between our simple CPU (1 instruction per cycle, 29 cycles for our tiny program) and a modern processor is staggering. How did engineers bridge it?

**The answer is the second half of this course.**

Here is a taste of what is coming:

**Pipelining (Chapter 20):** Instead of waiting for one instruction to finish all three stages before starting the next, overlap the stages. While instruction 5 is executing, instruction 6 is decoding, and instruction 7 is being fetched. Like an assembly line in a factory — multiple cars being assembled simultaneously at different stations.

**Caches (Chapter 21):** Our CPU had to go to main memory for every instruction and every data load. Modern CPUs have small, extremely fast memory banks (caches) right next to the processor core. Most memory accesses hit the cache and take 1–4 cycles instead of 100+ cycles for main RAM.

**Branch Prediction (Chapter 22):** Every time our CPU hit a JEQ, it had to wait to see if the jump would be taken. Modern CPUs guess (predict) which way the branch will go and start executing speculatively. They are right about 95% of the time.

**Superscalar Execution (Chapter 23):** Our ALU performed one operation per cycle. Modern CPUs have multiple ALUs running in parallel, finding independent instructions in the stream and executing them simultaneously.

**Out-of-Order Execution (Chapter 24):** If instruction 5 depends on instruction 4 (which is still computing), and instruction 6 is independent of both — execute instruction 6 first while waiting for 4 to finish. The CPU reorders instructions dynamically, like a waiter who brings you your drinks while your food is still cooking.

All of these techniques are built on the same foundation we just traced through. The transistors are still there. The flip-flops are still there. The ALU is still there. But they are organized with extraordinary cleverness to extract maximum speed from the hardware.

The question that drives the rest of this course is deceptively simple:

**We know how to make a CPU that is correct. How do we make it fast?**

---

## Summary

In this chapter, we:

1. **Defined a complete minimal CPU** with 4 registers (R0–R3), a PC, a FLAG register, and a 7-instruction instruction set (LOAD, STORE, ADD, SUB, CMP, JMP, JEQ).

2. **Wrote a real program** to compute the sum 5+4+3+2+1=15 using a loop, 10 instructions, and 16 bytes of memory.

3. **Traced all 29 clock cycles** in full detail, showing register values, memory state, FLAG changes, and the physical operations at each step.

4. **Examined the branch mechanism** — the multiplexer that chooses between PC+1 and a jump target, controlled by the FLAG register.

5. **Connected every instruction to earlier chapters**: transistors, binary, gates, adders, flip-flops, the ALU, control unit, and fetch-decode-execute cycle.

6. **Identified the limitations** of our simple CPU and previewed the techniques (pipelining, caching, branch prediction, superscalar execution) that make modern CPUs fast.

The key insight: **correctness and speed are separate problems**. Our CPU is correct. The rest of the course is about speed.

---

## Exercises

### Easy

1. Modify the program to compute the sum 1+2+3 (loop only 3 times). Write out the memory layout changes needed and predict the final value in the result register.

2. Trace the CMP instruction at cycle 12 (iteration 2, comparing R2=3 with R3=0) in full detail. Describe what the comparator circuit receives and outputs, and what value the FLAG register holds afterward.

3. Count how many total times each instruction type is executed over the entire 29-cycle program (e.g., how many ADD, how many CMP, etc.).

### Medium

4. Our loop runs from 5 down to 1 (5 iterations). Redesign the loop to count UP from 1 to 5 instead. What would you need to change in memory? What additional instruction or constant would you need?

5. Suppose the clock speed doubles (runs twice as fast). How many real-world seconds does our 29-cycle program take at (a) 1 MHz clock speed and (b) 1 GHz clock speed? What does this tell you about why clock speed matters?

6. We used 4 registers. Redraw the register file (as a diagram) showing the flip-flops for a 4-register, 4-bit-wide register file. How many D flip-flops are needed in total?

### Hard

7. Add a new instruction to our ISA: `INC Rx` (increment register Rx by 1, without needing to use another register). Write out the new opcode encoding, describe what hardware the control unit would need to generate, and rewrite the loop program to use INC instead of ADD R2, R1.

8. The JEQ instruction wastes a cycle every time the condition is false (the loop body always ends with JEQ then JMP). Design a new instruction `JNE addr` (jump if NOT equal, i.e., if FLAG=0) that would let you eliminate the JMP instruction entirely. Rewrite the loop program using JNE. How many cycles does the optimized version take?

9. Our CPU executes 29 instructions to sum 5 numbers. A pipelined version of this same CPU (studied in Chapter 20) could theoretically execute one instruction per cycle after filling the pipeline. However, each branch instruction (JEQ, JMP) causes the pipeline to stall for 2 extra cycles (called a "branch penalty"). Count how many branches our program executes, calculate the total pipeline stalls, and compute the effective CPI (cycles per instruction) for the pipelined version. Is it better or worse than our simple CPU?

---

*End of Chapter 13.*

*Part 2 begins with Chapter 14: "What Is an Instruction Set Architecture?" — where we move from our toy ISA to the real-world instruction sets that power every computer you have ever used.*
