# Chapter 08: The Arithmetic Logic Unit — The Math Engine

Imagine a factory that produces every computed result your computer ever generates. Every time you add two numbers in a spreadsheet, compare two values in a database, or shift bits in a graphics pipeline — that factory is a single circuit inside your CPU called the **Arithmetic Logic Unit**, or ALU. It is small enough to fit on a sliver of silicon, yet powerful enough to perform billions of operations per second. This chapter tears the factory open and shows you every machine inside.

---

## Table of Contents

1. [What the ALU Does](#1-what-the-alu-does)
2. [Inputs and Outputs: The ALU's Interface](#2-inputs-and-outputs-the-alus-interface)
3. [Building a 4-Bit ALU from Gates](#3-building-a-4-bit-alu-from-gates)
4. [Arithmetic Operations: Add and Subtract](#4-arithmetic-operations-add-and-subtract)
5. [Logical Operations: AND, OR, XOR, NOT](#5-logical-operations-and-or-xor-not)
6. [Shift Operations: Multiply and Divide by 2](#6-shift-operations-multiply-and-divide-by-2)
7. [The Status Flags](#7-the-status-flags)
8. [Concrete Examples: 5+7, 15-8, Flag Settings](#8-concrete-examples-57-15-8-flag-settings)
9. [How Flags Enable Conditional Branching](#9-how-flags-enable-conditional-branching)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What the ALU Does

The ALU is the part of the CPU that actually *computes*. Everything else in the CPU — the control unit, the registers, the buses — exists to feed the ALU with data and to ferry its results away. The ALU itself has one job: take numbers in, do something to them, hand results out.

### The Kitchen Analogy

Think of the CPU as a professional kitchen. Ingredients (data) sit in the pantry (memory) or on the counter (registers). The head chef (control unit) decides what dish to make next and calls out orders. The actual cooking — chopping, mixing, heating — happens at a single station called the **ALU stove**. Every dish passes through it.

The ALU stove can do three kinds of cooking:

| Category | What it does | Examples |
|---|---|---|
| Arithmetic | Math with numbers | Add, Subtract |
| Logical | Bitwise logic | AND, OR, XOR, NOT |
| Shift | Move bits left or right | Left shift, Right shift |

That is the complete menu. Every computation in every program ever written boils down to combinations of those three categories.

### Why One Unit?

You might wonder: why not have a separate adder, a separate AND unit, a separate shifter? The answer is silicon real estate and speed. A single unified ALU is cheaper to build, simpler to route signals to, and shares transistors between operations. Modern CPUs actually have multiple ALUs running in parallel — but each individual ALU is still this same unified design.

---

### Quick Check 1

1. Name the three categories of operations an ALU can perform.
2. In the kitchen analogy, what represents the ALU?
3. Why do CPUs use a single unified ALU rather than separate circuits for each operation?

---

## 2. Inputs and Outputs: The ALU's Interface

An ALU has a precisely defined interface — a specific set of wires going in and wires going out. Understanding this interface is like understanding the plugs and sockets of an appliance before you try to open it up.

### Inputs

An ALU has three main inputs:

```
     ┌─────────────────────────────┐
     │                             │
 A ──┤  Operand A (N bits)         │
     │                             │
 B ──┤  Operand B (N bits)         │
     │                             │
OP ──┤  Opcode (selects operation) │
     │                             │
     └─────────────────────────────┘
```

- **Operand A**: The first number (N bits wide — 4, 8, 16, 32, or 64 bits depending on the CPU).
- **Operand B**: The second number (same width as A).
- **Opcode**: A short code that tells the ALU *which* operation to perform. For a simple ALU, this might be 3 bits, allowing 8 different operations. A full CPU ALU might have a 6-bit opcode for 64 operations.

**Operand** is just a fancy word for "thing we are operating on." If you add 5 + 7, then 5 and 7 are the operands. The operation is addition.

### Outputs

```
     ┌─────────────────────────────┐
     │                             │
     │                         ───┤ Result (N bits)
     │                             │
     │                         ───┤ Zero flag (Z)
     │                             │
     │                         ───┤ Carry flag (C)
     │                             │
     │                         ───┤ Negative flag (N)
     │                             │
     │                         ───┤ Overflow flag (V)
     └─────────────────────────────┘
```

- **Result**: The N-bit answer to the computation.
- **Status flags**: Single bits that report *special conditions* about the result. We will dedicate a full section to these — they are more important than they look.

### The Complete Picture

```
                        Opcode (e.g. 3 bits)
                              │
                              ▼
    A (N bits) ──────────┬──[ALU]──────────── Result (N bits)
                         │   │
    B (N bits) ───────────┘  └── Flags: Z, C, N, V
```

The opcode selects the operation. A and B are the raw material. Result and flags are the output. That is the complete interface. Simple on the outside — astonishing on the inside.

---

### Quick Check 2

1. What are the three inputs to an ALU?
2. What is an "operand"?
3. List the four status flags an ALU typically outputs.

---

## 3. Building a 4-Bit ALU from Gates

Now let us build one. We will construct a 4-bit ALU — one that operates on 4-bit numbers (values 0 through 15). This is small enough to reason about completely, yet complex enough to show all the real principles.

### Step 1: The Full Adder (The Core Building Block)

In Chapter 05 we built a Full Adder — a circuit that adds three 1-bit inputs (A, B, and Carry-In) and produces a 1-bit Sum and a 1-bit Carry-Out.

```
        A ──┐
            ├── [Full Adder] ── Sum
        B ──┘        │
                     └── Carry-Out
       Cin ──────────────────────
```

To add two 4-bit numbers, we chain four Full Adders together — the carry-out of each feeds into the carry-in of the next:

```
     Bit 3            Bit 2            Bit 1            Bit 0
  ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐
  │         │      │         │      │         │      │         │
A3┤         │    A2┤         │    A1┤         │    A0┤         │
B3┤   FA3   │    B2┤   FA2   │    B1┤   FA1   │    B0┤   FA0   │
  │         │      │         │      │         │      │         │
  │  Sum3   │      │  Sum2   │      │  Sum1   │      │  Sum0   │
  │  Cout───┼─────►│  Cout───┼─────►│  Cout───┼─────►│  Cout  ◄── Cin=0
  └─────────┘      └─────────┘      └─────────┘      └─────────┘
     │                 │                 │                 │
    S3               S2               S1               S0
```

This is called a **Ripple Carry Adder** because the carry "ripples" from the least significant bit (bit 0) to the most significant bit (bit 3). The carry-in of the entire chain (at FA0) starts at 0.

### Step 2: Adding Logic Operations

We need to handle AND, OR, XOR, and NOT. These are simpler — they operate bit by bit with no carry:

```
     A3─┬─[AND]──  A2─┬─[AND]──  A1─┬─[AND]──  A0─┬─[AND]──
     B3─┘          B2─┘          B1─┘          B0─┘
```

Each bit position gets its own AND gate, OR gate, and XOR gate. We compute all three results simultaneously and then use a **multiplexer** (MUX) to select which result passes through to the output.

### Step 3: The Operation Selector (MUX)

A multiplexer is like a multi-way switch. The opcode selects which input gets routed to the output:

```
         ┌──────────┐
ADD ─────┤ 0        │
AND ─────┤ 1        ├─── Result
OR  ─────┤ 2        │
XOR ─────┤ 3        │
         └────▲─────┘
              │
           Opcode
```

With a 3-bit opcode, we can select among 8 operations. For our 4-bit ALU, a simple opcode table might look like:

| Opcode | Operation | Notes |
|--------|-----------|-------|
| 000 | ADD | A + B |
| 001 | SUB | A - B (using two's complement) |
| 010 | AND | A AND B, bitwise |
| 011 | OR | A OR B, bitwise |
| 100 | XOR | A XOR B, bitwise |
| 101 | NOT | NOT A (B ignored) |
| 110 | SHL | Shift A left by 1 |
| 111 | SHR | Shift A right by 1 |

### The Complete 4-Bit ALU Structure

```
        A[3:0]        B[3:0]       Opcode[2:0]
           │              │              │
           ▼              ▼              │
    ┌──────────────────────────┐         │
    │                          │         │
    │   ┌──────────────────┐   │         │
    │   │  4-bit Adder     │───┤ ADD     │
    │   └──────────────────┘   │         │
    │                          │         │
    │   ┌──────────────────┐   │         │
    │   │  4-bit AND unit  │───┤ AND     │
    │   └──────────────────┘   │         ▼
    │                          │   ┌──────────┐
    │   ┌──────────────────┐   │   │          │
    │   │  4-bit OR unit   │───┤───►  4-to-1  ├── Result[3:0]
    │   └──────────────────┘   │   │   MUX    │
    │                          │   │          │
    │   ┌──────────────────┐   │   └──────────┘
    │   │  4-bit XOR unit  │───┤ XOR
    │   └──────────────────┘   │
    │                          │
    │   ┌──────────────────┐   │
    │   │  Shift unit      │───┤ SHL/SHR
    │   └──────────────────┘   │
    │                          │
    │   ┌──────────────────┐   │
    │   │  Flag generator  │───┼── Z, C, N, V flags
    │   └──────────────────┘   │
    └──────────────────────────┘
```

All the sub-units compute in parallel — simultaneously. The multiplexer just picks which answer to deliver. This parallelism is part of what makes hardware fast.

---

### Quick Check 3

1. What is a Ripple Carry Adder and why is it called "ripple"?
2. What is the role of the multiplexer in an ALU?
3. If an opcode is 3 bits wide, how many different operations can it select?

---

## 4. Arithmetic Operations: Add and Subtract

### Addition

Addition in binary follows the same rules as addition in decimal — you add column by column and carry when the result exceeds the base.

In decimal:   In binary:
  1 8           1 0 0 1 0
+ 1 3         + 0 1 1 0 1
------        ---------
  3 1           1 1 1 1 1

Binary addition rules for each bit pair:

| A | B | Cin | Sum | Cout |
|---|---|-----|-----|------|
| 0 | 0 |  0  |  0  |  0   |
| 0 | 0 |  1  |  1  |  0   |
| 0 | 1 |  0  |  1  |  0   |
| 0 | 1 |  1  |  0  |  1   |
| 1 | 0 |  0  |  1  |  0   |
| 1 | 0 |  1  |  0  |  1   |
| 1 | 1 |  0  |  0  |  1   |
| 1 | 1 |  1  |  1  |  1   |

This is exactly what one Full Adder implements. Chain four of them and you can add two 4-bit numbers.

### Subtraction: The Two's Complement Trick

Here is where it gets clever. Subtraction is *not* implemented with a separate subtracter circuit. Instead, the ALU reuses the adder with a clever mathematical trick called **two's complement**.

**Two's complement** is a way of representing negative numbers in binary such that:

> A - B = A + (-B)   and   -B = (NOT B) + 1

The rule is: to negate a binary number, flip all its bits and add 1.

**Why does this work?** Think about it on a number wheel. An 8-bit number has 256 possible values. If you add 1 to 255 (binary `11111111`), you get 256 — but that overflows and wraps around to 0. So 255 behaves like -1 on the wheel. 254 behaves like -2. And so on. The "complement" of a number is its "negative partner" on this wheel.

**The hardware consequence**: To subtract B from A:

1. Flip all bits of B (bitwise NOT)
2. Add 1 (by setting the carry-in of the adder to 1 instead of 0)
3. Add A

The same adder circuit handles both addition and subtraction. The only hardware difference is:
- **ADD**: send B straight in, Cin = 0
- **SUB**: send NOT(B) in, Cin = 1

```
                    ADD mode: B passes straight through
                    SUB mode: B is inverted (NOT B)
                         │
              B[3:0] ────┤
                    NOT gates (controlled by SUB signal)
                         │
                         ▼
              ┌──────────────────┐
    A[3:0] ──►│  4-bit Adder     ├── Result[3:0]
              │                  │
   Cin ───────┤ (0 for ADD,      ├── Carry-out
              │  1 for SUB)      │
              └──────────────────┘
```

This is one of the most elegant ideas in computer architecture. One adder circuit performs two operations. Hardware designers love this kind of reuse.

---

### Quick Check 4

1. What are the two steps to negate a number using two's complement?
2. How does the ALU perform subtraction without a separate subtraction circuit?
3. What value is fed into the carry-in wire for subtraction (instead of the normal 0)?

---

## 5. Logical Operations: AND, OR, XOR, NOT

Logical operations work on individual bits independently. There is no carrying, no borrowing — each bit position is completely isolated.

### AND

AND outputs 1 only when *both* inputs are 1. Think of it as a strict check: "are both conditions true?"

```
  A: 1 0 1 1
  B: 1 1 0 1
     -------
A&B: 1 0 0 1
```

**Use case**: Masking — keeping only certain bits. To extract the lower 4 bits of a byte, AND it with `0000 1111`. Every bit in the "mask" that is 0 forces the output to 0; every bit that is 1 lets the original bit through.

### OR

OR outputs 1 when *at least one* input is 1. Think of it as a permissive check: "is either condition true?"

```
  A: 1 0 1 0
  B: 0 1 1 0
     -------
A|B: 1 1 1 0
```

**Use case**: Setting bits. To force bit 3 of a byte to 1 without touching the others, OR the byte with `0000 1000`.

### XOR (Exclusive OR)

XOR outputs 1 when the inputs are *different* from each other. It is the "one or the other, but not both" operation.

```
  A: 1 0 1 1
  B: 1 1 0 1
     -------
A^B: 0 1 1 0
```

**Use case 1**: Toggling bits. XOR with `0000 0001` flips bit 0. Each 1 in the XOR mask flips the corresponding bit; each 0 leaves it alone.

**Use case 2**: Checking equality. If A XOR B = 0, then A and B are identical (every bit was the same, so every XOR was 0).

**Use case 3**: Simple encryption. XOR a message with a key to encrypt; XOR again with the same key to decrypt.

### NOT

NOT flips every bit (also called "bitwise inversion" or "one's complement"):

```
  A: 1 0 1 1 0 0 1 0
NOT: 0 1 0 0 1 1 0 1
```

**Use case**: Part of the two's complement negation we just studied. NOT is also used to invert control signals and flags.

### Logical Operations Side by Side

| Operation | Symbol | Rule | Example (A=1010, B=1100) | Result |
|-----------|--------|------|---------------------------|--------|
| AND | & | 1 only if both 1 | 1010 AND 1100 | 1000 |
| OR | \| | 1 if either 1 | 1010 OR 1100 | 1110 |
| XOR | ^ | 1 if different | 1010 XOR 1100 | 0110 |
| NOT | ~ | flip every bit | NOT 1010 | 0101 |

---

### Quick Check 5

1. What does AND do to a bit when the mask bit is 0?
2. How would you use OR to set bit 5 of a byte to 1?
3. What is the result of XOR-ing any value with itself?

---

## 6. Shift Operations: Multiply and Divide by 2

Shifting moves all bits left or right. It is one of the most useful operations in computing because it connects directly to multiplication and division by powers of two.

### Left Shift

A left shift moves every bit one position to the left. The rightmost bit is filled with 0. The leftmost bit is lost (it falls off the edge).

```
Original:   0 0 1 0 1 1 0 1   (= 45 in decimal)
Left shift: 0 1 0 1 1 0 1 0   (= 90 in decimal)
                         └─ new 0 inserted here
                   └─── original bits move left
```

**Left shift by 1 = multiply by 2**. This works for the same reason that shifting a decimal number left multiplies by 10: the positional values of all bits double when they move one position left.

Left shifting by N positions multiplies by 2^N:
- Shift left 1: multiply by 2
- Shift left 2: multiply by 4
- Shift left 3: multiply by 8

### Right Shift

A right shift moves every bit one position to the right. The leftmost bit is filled with 0 (for unsigned numbers) or with the sign bit (for signed numbers). The rightmost bit is lost.

```
Original:    0 0 1 0 1 1 0 0   (= 44 in decimal)
Right shift: 0 0 0 1 0 1 1 0   (= 22 in decimal)
             └─ new 0 inserted here
                  bits move right ──┘
```

**Right shift by 1 = divide by 2** (integer division — the remainder is discarded).

### Logical vs. Arithmetic Right Shift

For unsigned numbers, we always insert a 0 on the left. This is a **logical right shift**.

For signed numbers (using two's complement), we want dividing a negative number by 2 to give us a negative result. So we fill the leftmost bit with the original sign bit (0 for positive, 1 for negative). This is called an **arithmetic right shift**.

```
Signed number:  1 0 1 1 0 0 0 0  (= -64 in two's complement)
Arithmetic SHR: 1 1 0 1 1 0 0 0  (= -32) -- sign bit preserved
Logical SHR:    0 1 0 1 1 0 0 0  (= +88) -- wrong sign!
```

Modern CPUs provide both. The instruction set specifies which to use.

### Why Shifts Are Valuable

Hardware multiplication is expensive — a full multiplier circuit uses many thousands of gates. But a shift is essentially free — it is just rewiring. Compilers exploit this aggressively:

- `x * 8` becomes `x << 3` (shift left 3)
- `x / 4` becomes `x >> 2` (arithmetic shift right 2)
- `x * 10` becomes `(x << 3) + (x << 1)` (8x + 2x)

| Shift | Equivalent operation |
|-------|----------------------|
| `<< 1` | multiply by 2 |
| `<< 2` | multiply by 4 |
| `<< 3` | multiply by 8 |
| `<< N` | multiply by 2^N |
| `>> 1` | divide by 2 |
| `>> 2` | divide by 4 |
| `>> N` | divide by 2^N |

---

### Quick Check 6

1. What happens to the value of a number when you left shift it by 2 positions?
2. What is the difference between a logical right shift and an arithmetic right shift?
3. How would a compiler calculate `x * 12` using only shifts and addition?

---

## 7. The Status Flags

Status flags are single-bit outputs of the ALU. Each flag reports a specific condition about the result of the most recent operation. They are stored in a special register called the **status register**, **flags register**, or **condition code register** (different CPUs use different names, but the idea is the same).

Think of flags like indicator lights on a car dashboard. The engine (ALU) finishes its work, and then a set of lights illuminate to tell you what happened: "fuel low," "engine hot," "seatbelt off." The flags are the ALU's dashboard.

### The Four Core Flags

#### Zero Flag (Z)

The Zero flag is 1 if the result is exactly zero. It is 0 otherwise.

```
Result = 0000 0000  →  Z = 1
Result = 0000 0001  →  Z = 0
Result = 1111 1111  →  Z = 0
```

**What it means**: The two operands were equal (if we were subtracting), or the result was zero (if we were adding or ANDing).

**When it matters**: Almost every conditional jump instruction checks this flag. "Jump if zero" (JZ) and "jump if not zero" (JNZ) are among the most common instructions in all programs.

#### Carry Flag (C)

The Carry flag is 1 if the addition produced a carry-out from the most significant bit position. It is 0 otherwise.

```
  1111 1111
+ 0000 0001
-----------
  0000 0000   Result
          1   Carry flag = 1  (the addition "overflowed" past 8 bits)
```

**What it means**: For unsigned addition, a carry-out means the true result did not fit in N bits — it is larger than the maximum unsigned value.

**For subtraction**: The carry flag is inverted and called the **borrow flag** on some architectures. If A < B, the subtraction needs to borrow, and the carry/borrow flag signals this.

**Multi-precision arithmetic**: The carry flag enables adding numbers larger than the register width. To add two 64-bit numbers on a 32-bit CPU, you add the lower 32 bits, save the carry, then add the upper 32 bits plus the carry.

#### Negative Flag (N) — Also Called Sign Flag (S)

The Negative flag is simply a copy of the most significant bit (the sign bit) of the result.

```
Result = 0111 1111  →  N = 0  (positive, since bit 7 is 0)
Result = 1000 0000  →  N = 1  (negative in two's complement, bit 7 is 1)
```

**What it means**: In two's complement arithmetic, the most significant bit is the sign bit. N=1 means the result is negative.

**When it matters**: "Jump if negative" (JN/JMI) checks this flag. Sorting algorithms, comparison functions, and anything involving signed comparisons rely on it.

#### Overflow Flag (V)

The Overflow flag is the trickiest. It is 1 when the result of a *signed* arithmetic operation is incorrect — when the result was too large or too small to fit in the signed range.

For an N-bit two's complement number:
- Maximum positive value: 2^(N-1) - 1   (e.g., 127 for 8-bit)
- Minimum negative value: -2^(N-1)       (e.g., -128 for 8-bit)

Overflow happens when:
- A positive + positive = negative (result wrapped around the top)
- A negative + negative = positive (result wrapped around the bottom)

```
8-bit example (max positive = 127):
  0111 1111  (+127)
+ 0000 0001  (+1)
-----------
  1000 0000  (-128 in two's complement!)  ← WRONG!
  
  V = 1 (signed overflow occurred)
```

**Detecting overflow**: The hardware detects overflow by checking whether the carry into the sign bit differs from the carry out of the sign bit. If they differ, overflow occurred.

```
V = Carry_into_MSB XOR Carry_out_of_MSB
```

#### Summary Table of Flags

| Flag | Full Name | Set when... | Common use |
|------|-----------|-------------|------------|
| Z | Zero | Result == 0 | Check equality |
| C | Carry | Unsigned overflow or borrow | Multi-precision math, unsigned compare |
| N | Negative | Result MSB == 1 | Signed comparisons |
| V | Overflow | Signed result out of range | Detect signed arithmetic errors |

---

### Quick Check 7

1. When is the Zero flag set?
2. What is the difference between the Carry flag and the Overflow flag?
3. How is the Negative flag determined from the result?

---

## 8. Concrete Examples: 5+7, 15-8, Flag Settings

Let us run the ALU through three complete examples and trace every flag.

### Example 1: 5 + 7 (4-bit unsigned addition)

In 4-bit binary:
- 5 = `0101`
- 7 = `0111`

```
  Step-by-step ripple carry:

  Bit:   3    2    1    0
  A:     0    1    0    1   (= 5)
  B:     0    1    1    1   (= 7)
  Cin:   ?    ?    ?    0   (starts at 0 for addition)

  Bit 0:  1 + 1 + 0 = 10 binary  →  Sum=0, Cout=1
  Bit 1:  0 + 1 + 1 = 10 binary  →  Sum=0, Cout=1
  Bit 2:  1 + 1 + 1 = 11 binary  →  Sum=1, Cout=1
  Bit 3:  0 + 0 + 1 = 01 binary  →  Sum=1, Cout=0

  Result: 1100 (= 12 in decimal) ✓
```

Flag settings:
- **Z = 0**: Result (12) is not zero
- **C = 0**: No carry out from bit 3 (the result fits in 4 bits)
- **N = 1**: Bit 3 of result is 1 — but wait! In *unsigned* 4-bit arithmetic, 12 is valid. This flag is meaningful only for signed interpretation.
- **V = 1**: In signed 4-bit, range is -8 to +7. 5 + 7 = 12, which exceeds +7 — so signed overflow occurred (two positives added to give a "negative" bit pattern). If you were treating these as signed numbers, the result is wrong.

This illustrates something crucial: the ALU sets all flags every time. Whether the flags are *meaningful* depends on whether you are treating your numbers as signed or unsigned. That is the programmer's (or compiler's) responsibility, not the hardware's.

### Example 2: 15 - 8 (4-bit subtraction using two's complement)

In 4-bit binary:
- 15 = `1111`
- 8  = `1000`

Step 1: Compute NOT(8) = NOT(`1000`) = `0111`

Step 2: Add A + NOT(B) + 1 (carry-in = 1):

```
  A:        1 1 1 1   (= 15)
  NOT(B):   0 1 1 1   (NOT of 8)
  Cin:      1         (set to 1 for subtraction)

  Bit 0:  1 + 1 + 1 = 11 binary  →  Sum=1, Cout=1
  Bit 1:  1 + 1 + 1 = 11 binary  →  Sum=1, Cout=1
  Bit 2:  1 + 1 + 1 = 11 binary  →  Sum=1, Cout=1
  Bit 3:  1 + 0 + 1 = 10 binary  →  Sum=0, Cout=1

  Result: 0111 (= 7 in decimal) ✓  (15 - 8 = 7)
```

Flag settings:
- **Z = 0**: Result (7) is not zero
- **C = 1**: Carry-out from bit 3 is 1. For subtraction, C=1 means no borrow (A >= B). This is correct: 15 >= 8.
- **N = 0**: Bit 3 is 0, result is non-negative
- **V = 0**: 7 is within the 4-bit signed range (-8 to +7), no overflow

### Example 3: 8 - 15 (subtraction with a negative result)

- 8  = `1000`
- 15 = `1111`

Step 1: NOT(15) = NOT(`1111`) = `0000`

Step 2: Add 8 + NOT(15) + 1:

```
  A:        1 0 0 0   (= 8)
  NOT(B):   0 0 0 0   (NOT of 15)
  Cin:      1

  Bit 0:  0 + 0 + 1 = 01  →  Sum=1, Cout=0
  Bit 1:  0 + 0 + 0 = 00  →  Sum=0, Cout=0
  Bit 2:  0 + 0 + 0 = 00  →  Sum=0, Cout=0
  Bit 3:  1 + 0 + 0 = 01  →  Sum=1, Cout=0

  Result: 1001 (= -7 in 4-bit two's complement) ✓
```

Check: 8 - 15 = -7. In 4-bit two's complement, -7 is represented as `1001`. Correct!

Flag settings:
- **Z = 0**: Result is not zero
- **C = 0**: Carry-out is 0. For subtraction, C=0 means a borrow occurred (A < B). Correct: 8 < 15.
- **N = 1**: Bit 3 is 1, result is negative. Correct: -7 is negative.
- **V = 0**: -7 is within the 4-bit signed range, no overflow.

---

### Quick Check 8

1. In the example 5 + 7 on a 4-bit ALU, why might the Overflow flag be set even though the unsigned result (12) is correct?
2. What does C = 0 mean after a subtraction operation?
3. What 4-bit binary pattern represents -7 in two's complement?

---

## 9. How Flags Enable Conditional Branching

Flags are useless without something to read them. Their whole purpose is to feed into the **control unit**, which uses them to make decisions about program flow.

### The Compare Instruction

Most CPUs have a compare instruction (CMP on x86, SUBS on ARM). It subtracts two values and sets the flags — but discards the result. It is subtraction used purely for its side effects (the flags).

```assembly
; x86-style assembly example (conceptual)
MOV AX, 15        ; Load 15 into register AX
MOV BX, 8         ; Load 8 into register BX
CMP AX, BX        ; Compute 15 - 8, set flags, discard result
; Flags now: Z=0, C=1 (no borrow), N=0, V=0
JG  positive_path ; Jump if Greater (checks N, V, Z flags)
; ... if we reach here, AX <= BX
```

### The Conditional Jump Family

Every conditional jump instruction tests a specific flag combination:

| Instruction | Meaning | Flags tested |
|-------------|---------|--------------|
| JZ / JE | Jump if zero / equal | Z = 1 |
| JNZ / JNE | Jump if not zero / not equal | Z = 0 |
| JC | Jump if carry | C = 1 |
| JNC | Jump if no carry | C = 0 |
| JN / JMI | Jump if negative | N = 1 |
| JP / JPL | Jump if positive | N = 0 |
| JO | Jump if overflow | V = 1 |
| JNO | Jump if no overflow | V = 0 |
| JL / JLT | Jump if less than (signed) | N != V |
| JG / JGT | Jump if greater than (signed) | Z=0 and N=V |
| JB / JLO | Jump if below (unsigned) | C = 1 |
| JA / JHI | Jump if above (unsigned) | C=0 and Z=0 |

### A Complete Conditional Example

Here is how a simple `if (a > b)` in a high-level language becomes ALU operations and flag checks:

```
High-level code:
    if (a > b) {
        result = a - b;
    }

In the CPU:
    1. Load a into register R1
    2. Load b into register R2
    3. ALU: CMP R1, R2          (compute R1 - R2, set flags)
       - If R1 > R2 (signed): Z=0, N=V (N and V equal each other)
    4. Control unit checks flags: is it "greater than"?
       - Yes: continue to the subtraction
       - No:  jump past the subtraction
    5. ALU: SUB R1, R2          (compute actual a - b)
    6. Store result
```

### The Feedback Loop

```
        ┌─────────┐
        │   ALU   │──── Z, C, N, V flags ──┐
        └─────────┘                         │
             ▲                              ▼
             │                     ┌──────────────────┐
     operands│                     │  Control Unit     │
             │                     │  reads flags and  │
        ┌────┴────┐                │  decides: normal  │
        │Registers│                │  flow or branch?  │
        └────▲────┘                └────────┬──────────┘
             │                              │
             │      next instruction        │
             └──────────────────────────────┘
```

Every conditional statement in every program ever written — every `if`, every `while`, every `for` loop — ultimately reduces to this cycle: the ALU computes, sets flags, the control unit reads the flags, decides whether to branch, and fetches the next instruction accordingly.

This is the machinery that makes programs *programmable*. Without flags and conditional jumps, a CPU could only execute fixed sequences of operations. Flags are what give programs the ability to make decisions.

---

### Quick Check 9

1. What does the CMP instruction do with the result of its subtraction?
2. What flag condition represents "less than" in signed comparison?
3. What is the difference between JL (jump if less, signed) and JB (jump if below, unsigned)?

---

## Summary

The ALU is the engine at the heart of all computation. Here is what we covered:

**Interface**: The ALU takes two N-bit operands (A and B) and an opcode, and produces an N-bit result plus four status flags (Z, C, N, V).

**Building an ALU**: A 4-bit ALU chains four Full Adders for arithmetic, uses parallel AND/OR/XOR/NOT gates for logic, adds a shift unit, and selects between them with a multiplexer controlled by the opcode. All sub-units compute in parallel — the MUX picks the winner.

**Addition**: Ripple carry adder chains Full Adders from LSB to MSB. Carry propagates upward.

**Subtraction**: Reuses the adder. To compute A - B, the ALU computes A + NOT(B) + 1 (two's complement negation). No separate subtractor needed.

**Logical operations**: AND, OR, XOR, NOT operate independently on each bit pair. AND masks bits, OR sets bits, XOR toggles bits, NOT flips all bits.

**Shift operations**: Left shift multiplies by 2 (per position); right shift divides by 2. Arithmetic right shift preserves the sign bit for signed numbers.

**Flags**:
- **Z (Zero)**: Result is zero
- **C (Carry)**: Unsigned overflow or borrow
- **N (Negative)**: MSB of result is 1 (result is negative in two's complement)
- **V (Overflow)**: Signed arithmetic produced an out-of-range result

**Branching**: Flags feed the control unit. Conditional jump instructions test flag combinations to implement `if/else`, loops, and all other control flow in programs.

The wonder of the ALU is its economy. A handful of gate types, chained and multiplexed, can perform every arithmetic and logical operation a program needs. The architecture of modern CPUs, with their billions of transistors, still rests on this same foundation: carry chains, two's complement, and a multiplexer selecting among parallel results.

---

## Exercises

### Easy

1. Compute 6 + 9 in 4-bit binary, showing the carry chain at each bit position. What are the Z, C, N flags?

2. What is the two's complement representation of -5 in 8 bits? Show the NOT-then-add-1 process step by step.

3. Compute the following 8-bit operations and give the result in both binary and hex:
   - `1010 1010` AND `1111 0000`
   - `1010 1010` OR `0000 1111`
   - `1010 1010` XOR `1111 1111`

### Medium

4. A 4-bit ALU computes 7 + 7. Show the binary addition, identify the result, and determine the settings of all four flags (Z, C, N, V). Explain why the Overflow flag is set even though the unsigned result fits in 4 bits.

5. An 8-bit ALU computes 200 + 100 (unsigned values).
   - Show the binary addition.
   - The result wraps around — what is the 8-bit result?
   - Which flag signals that this happened?
   - What is the true mathematical result, and how could you recover it?

6. Design a truth table for the Overflow flag detector. The inputs are: Carry into MSB (C_in_MSB), Carry out of MSB (C_out_MSB). The output is V. Verify your table matches the formula `V = C_in_MSB XOR C_out_MSB`.

### Hard

7. A compiler encounters the expression `x * 36`. Show how to compute this using only shift and add operations (no multiply instruction). What is the minimum number of shift and add operations needed? Write out the sequence.

8. Explain why a Ripple Carry Adder gets slower as you add more bits. If each Full Adder takes 2 nanoseconds to propagate a carry, how long does a 64-bit ripple carry addition take in the worst case? Research the term "Carry Lookahead Adder" and explain in your own words how it solves this problem.

9. The flags from one ALU operation are: Z=0, C=0, N=1, V=1. An instruction JL (jump if less than signed) is about to execute.
   - The condition for JL is N != V. Does the jump occur?
   - What combination of A and B values (in 8-bit signed arithmetic) could have produced these exact flag settings?
   - If the ALU had computed A - B to get these flags, is A less than B? Explain using the flag logic.

---

*Next chapter: Registers — The CPU's Working Memory. We will look at how registers are built from flip-flops, why they are so much faster than RAM, and how the register file is designed to feed the ALU with operands at full speed.*
