# Chapter 05: Combining Gates — Adders, Multiplexers, and Decoders

Individual logic gates are useful, but the real power comes from combining them into larger circuits. In this chapter, you will build the most important building blocks of any CPU: circuits that do addition, select between inputs, decode addresses, and compare values. By the end, you will have built — in principle — the core of a calculator.

## Table of Contents

1. [The Half Adder — Adding Two Bits](#1-the-half-adder--adding-two-bits)
2. [The Full Adder — Adding with Carry](#2-the-full-adder--adding-with-carry)
3. [The Ripple-Carry Adder — Adding Multi-Bit Numbers](#3-the-ripple-carry-adder--adding-multi-bit-numbers)
4. [The Subtractor](#4-the-subtractor)
5. [The Multiplexer (MUX)](#5-the-multiplexer-mux)
6. [The Demultiplexer (DEMUX)](#6-the-demultiplexer-demux)
7. [The Decoder](#7-the-decoder)
8. [The Encoder](#8-the-encoder)
9. [The Comparator](#9-the-comparator)
10. [Putting It Together — The Simple ALU](#10-putting-it-together--the-simple-alu)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Half Adder — Adding Two Bits

Let us build a circuit to add two 1-bit numbers. In binary, there are only four possible additions:

```
0 + 0 = 0    (result: 0, carry: 0)
0 + 1 = 1    (result: 1, carry: 0)
1 + 0 = 1    (result: 1, carry: 0)
1 + 1 = 10   (result: 0, carry: 1)  ← 1+1 = 2 in decimal = 10 in binary
```

So adding two bits produces two outputs: a **Sum** bit and a **Carry** bit.

Looking at the truth table:

```
┌───┬───┬──────────┬───────┐
│ A │ B │  Sum (S) │ Carry │
├───┼───┼──────────┼───────┤
│ 0 │ 0 │    0     │   0   │
│ 0 │ 1 │    1     │   0   │
│ 1 │ 0 │    1     │   0   │
│ 1 │ 1 │    0     │   1   │
└───┴───┴──────────┴───────┘
```

The Sum column is exactly the XOR truth table.
The Carry column is exactly the AND truth table.

Therefore:
- **Sum = A XOR B**
- **Carry = A AND B**

```
Circuit:
         ┌──[XOR]──── Sum
A ──┬────┤
    │    └──[AND]──── Carry
B ──┘
```

This two-gate circuit is the **Half Adder**. It is called "half" because it can only handle two inputs — it cannot handle an incoming carry from a previous addition. That is the full adder's job.

### Quick Check

> 1. What is the Sum output when A=1 and B=1 in a half adder?
> 2. What gate produces the Carry output?
> 3. Why is it called a "half" adder?

---

## 2. The Full Adder — Adding with Carry

When adding multi-bit numbers, each column position must also add the carry from the previous column. A **Full Adder** adds three bits: A, B, and Carry-In (Cin), producing Sum and Carry-Out (Cout).

```
Truth Table:
┌───┬───┬─────┬──────┬──────┐
│ A │ B │ Cin │ Sum  │ Cout │
├───┼───┼─────┼──────┼──────┤
│ 0 │ 0 │  0  │  0   │  0   │
│ 0 │ 0 │  1  │  1   │  0   │
│ 0 │ 1 │  0  │  1   │  0   │
│ 0 │ 1 │  1  │  0   │  1   │
│ 1 │ 0 │  0  │  1   │  0   │
│ 1 │ 0 │  1  │  0   │  1   │
│ 1 │ 1 │  0  │  0   │  1   │
│ 1 │ 1 │  1  │  1   │  1   │
└───┴───┴─────┴──────┴──────┘
```

The Sum is 1 when an odd number of inputs are 1 — which is the XOR of all three inputs:
- **Sum = A XOR B XOR Cin**

The Carry-Out is 1 when at least two inputs are 1:
- **Cout = (A AND B) OR (B AND Cin) OR (A AND Cin)**

### Building a Full Adder from Two Half Adders

```
Full Adder = Half Adder + Half Adder + OR gate

       ┌──── Half Adder 1 ────┐
A ─────┤ A            Sum ────┤──── Half Adder 2 ────┐
B ─────┤ B           Carry ──►│ A            Sum ────── Sum (final)
       └──────────────────────┤ B           Carry ──┐
Cin ────────────────────────── └─────────────────── │
                                                    OR ── Cout
```

A full adder uses 5 gates: 2 XOR, 2 AND, 1 OR — or equivalently, two half adders plus one OR gate.

### Quick Check

> 1. A full adder receives A=1, B=1, Cin=1. What are Sum and Cout?
> 2. How many gates does a full adder use?
> 3. What does Cin represent in multi-bit addition?

---

## 3. The Ripple-Carry Adder — Adding Multi-Bit Numbers

To add two 8-bit numbers (like 10110101 + 01100011), you need 8 full adders — one per bit position, chained so each one's Carry-Out feeds the next one's Carry-In.

```
Adding 8-bit numbers A[7:0] and B[7:0]:

Bit 0:  FA₀  ← A[0], B[0], Cin=0  → Sum[0], Cout₀
Bit 1:  FA₁  ← A[1], B[1], Cout₀  → Sum[1], Cout₁
Bit 2:  FA₂  ← A[2], B[2], Cout₁  → Sum[2], Cout₂
...
Bit 7:  FA₇  ← A[7], B[7], Cout₆  → Sum[7], Cout₇

If Cout₇ = 1, there is an overflow (result doesn't fit in 8 bits)
```

This is called a **Ripple-Carry Adder** because the carry "ripples" from bit 0 all the way to bit 7 before the final result is ready. The critical path — the longest path that determines the circuit's speed — runs through all the carry chain.

### The Speed Problem

For an n-bit ripple-carry adder, the worst case requires the carry to propagate through all n stages. With each stage taking ~1 gate delay, an 8-bit adder has ~8 gate delays. A 64-bit adder would need ~64 gate delays — this is too slow for modern CPUs running at billions of operations per second.

### The Carry-Lookahead Adder

Real CPUs use a **Carry-Lookahead Adder (CLA)**, which computes all carry bits in parallel rather than waiting for them to ripple. The key insight: for each bit position i:
- **Generate (G)**: bit position i will produce a carry regardless of Cin: G = A AND B
- **Propagate (P)**: bit position i will pass along an incoming carry: P = A XOR B

The carry into position i+1 is: Cin[i+1] = G[i] OR (P[i] AND Cin[i])

By expanding this formula recursively, you can compute all carry bits simultaneously using a tree of AND and OR gates with only ~4 levels of gate delay regardless of word size. Real 64-bit ALUs use this technique.

### Quick Check

> 1. How many full adders do you need to add two 32-bit numbers?
> 2. Why is a ripple-carry adder called "ripple-carry"?
> 3. What is the advantage of a carry-lookahead adder over a ripple-carry adder?

---

## 4. The Subtractor

Computers subtract by adding the negative of a number. Recall from Chapter 3: to negate a number in two's complement, invert all bits and add 1.

So: A - B = A + (NOT B) + 1

We can build this from an adder:
1. Invert all bits of B (using NOT gates)
2. Set the initial Carry-In to 1 (which adds the +1)
3. Feed into a regular adder

```
A  ─────────────────────────────────────┐
B  ──[NOT]────────────────────────────── ╔══════════════╗
                                         ║  N-bit Adder  ║──── Result (A - B)
Cin = 1 (subtract mode) ──────────────── ╚══════════════╝
```

Most CPU ALUs have a single subtract control signal: when it is 1, B is inverted and Cin is set to 1, turning the adder into a subtractor. The same hardware does both addition and subtraction — beautiful design economy.

### Quick Check

> 1. How does a computer perform subtraction without building a separate subtractor circuit?
> 2. What two operations are performed to negate a two's complement number?

---

## 5. The Multiplexer (MUX)

A **Multiplexer** (MUX) is a data selector — it picks one of several input signals and routes it to the output, based on a "select" signal.

A 2-to-1 MUX has:
- 2 data inputs: A and B
- 1 select input: S
- 1 output: Q

```
Truth Table:       Logic: Q = (NOT S AND A) OR (S AND B)
┌───┬───┬───┬───┐
│ S │ A │ B │ Q │  Circuit:
├───┼───┼───┼───┤
│ 0 │ 0 │ X │ 0 │  A ──[AND]──┐
│ 0 │ 1 │ X │ 1 │     ↑        ├──[OR]──── Q
│ 1 │ X │ 0 │ 0 │  NOT(S)      │
│ 1 │ X │ 1 │ 1 │              │
└───┴───┴───┴───┘  B ──[AND]──┘
                        ↑
                        S
```

A **4-to-1 MUX** has 4 data inputs and 2 select lines (2 bits can select among 4 options: 00, 01, 10, 11).

**MUX uses in CPUs:**
- The ALU uses MUXes to select which operation to perform based on the opcode
- The register file uses MUXes to select which register value to read
- Pipeline stages use MUXes to select between normal operands and forwarded values

### Quick Check

> 1. A 4-to-1 MUX has 4 data inputs. How many select lines does it need?
> 2. A 8-to-1 MUX? 16-to-1?
> 3. What does the MUX output when S=1, A=0, B=1 in a 2-to-1 MUX?

---

## 6. The Demultiplexer (DEMUX)

A **Demultiplexer** is the reverse of a MUX: it takes one input and routes it to one of several outputs based on the select signal.

```
1-to-4 DEMUX:
        ┌── Y0 (active when S=00)
Input ──┼── Y1 (active when S=01)
 S[1:0] ┼── Y2 (active when S=10)
        └── Y3 (active when S=11)
```

DEMUXes appear in memory systems (selecting which memory bank to access), data bus routing (steering data to the right register), and display controllers (selecting which output to drive).

---

## 7. The Decoder

A **Decoder** takes an n-bit binary input and activates exactly one of its 2ⁿ outputs — the one corresponding to the binary number represented by the input.

A 2-to-4 decoder:

```
Inputs A,B:    Output Y0=A̅B̅  Y1=A̅B  Y2=AB̅  Y3=AB

When A=0,B=0:  Y0=1, Y1=0, Y2=0, Y3=0
When A=0,B=1:  Y0=0, Y1=1, Y2=0, Y3=0
When A=1,B=0:  Y0=0, Y1=0, Y2=1, Y3=0
When A=1,B=1:  Y0=0, Y1=0, Y2=0, Y3=1
```

**Decoder uses in CPUs:**
- **Instruction decoding**: the opcode (e.g., 6 bits = 64 possible instructions) is decoded to activate exactly one operation circuit
- **Memory address decoding**: the memory address activates exactly one memory row (wordline) to read/write
- **Register file read**: the register number (5 bits for 32 registers) selects exactly one register's output

---

## 8. The Encoder

An **Encoder** is the inverse of a decoder: exactly one of 2ⁿ inputs is active, and the encoder outputs the n-bit binary code of which input is active.

An 8-to-3 priority encoder:

```
Inputs:  I7 I6 I5 I4 I3 I2 I1 I0
Outputs: A2 A1 A0 (binary number of the highest active input)

If I5=1 (and I6=0, I7=0): output = 101 (binary for 5)
If I3=1 (and I4=I5=I6=I7=0): output = 011 (binary for 3)
```

"Priority" means if multiple inputs are active, the highest-numbered one wins.

**Uses**: keyboard encoders (each key pressed outputs its binary code), interrupt controllers (each device requesting service gets encoded to a number).

---

## 9. The Comparator

A **comparator** checks the relationship between two binary numbers: Are they equal? Is A greater than B? Is A less than B?

**1-bit equality comparator**: A = B when A XOR B = 0. So:

```
Equal = NOT(A XOR B) = XNOR(A, B)
```

**n-bit equality comparator**: Two n-bit numbers are equal if and only if all their corresponding bits are equal. Chain n XNOR gates and AND their outputs:

```
A[7:0] equals B[7:0] when:
XNOR(A[7],B[7]) AND XNOR(A[6],B[6]) AND ... AND XNOR(A[0],B[0]) = 1
```

**Magnitude comparator**: To check if A > B in binary, compare from the most significant bit down: find the first bit position where they differ; the number with a 1 at that position is larger.

**Comparator uses**: 
- Cache tag comparison (is this the right cache line?)
- ALU comparison operations (for if/else branches)
- Sorting networks

### Quick Check

> 1. What gate produces "1 only when two bits are equal"?
> 2. To compare two 8-bit numbers for equality, how many gates does the circuit need?

---

## 10. Putting It Together — The Simple ALU

Now we can assemble all these building blocks into a simple **Arithmetic Logic Unit (ALU)** — the computing heart of a CPU.

A 4-bit ALU that supports 4 operations (Add, Subtract, AND, OR) selected by a 2-bit opcode:

```
               ┌────────────────────────────────────────────┐
               │                4-bit ALU                    │
               │                                             │
A[3:0] ───────►│──────┬──► 4-bit Adder  ─────────────────┐  │
               │      │   (ADD / SUB)                    │  │
B[3:0] ───────►│──────┤──► 4-bit AND ──────────────────► MUX─►─ Result[3:0]
               │      └──► 4-bit OR  ──────────────────► │  │
               │           Zero Gate ──────────────────►  │  │
Op[1:0] ──────►│────────────────────────────────────────► │  │
               │         (selects which result to output)  │  │
               └────────────────────────────────────────────┘
               
Flags output: Zero (result is 0), Carry, Overflow, Negative
```

The opcode selects one of four operation results via a 4-to-1 MUX. The selected result becomes the output. Flag circuits detect special conditions (zero result, overflow, etc.) and set status bits.

This simple circuit is the ancestor of every ALU ever built. The ALU in an Apple M4 is this same idea, but with 64-bit operands, dozens of operations, and extraordinary optimization — yet the core concept is identical.

### Quick Check

> 1. What does the MUX in the ALU select?
> 2. Why does the ALU output "flags" in addition to the result?
> 3. How would you extend this 4-operation ALU to support 8 operations?

---

## Summary

- A **half adder** adds two bits using XOR (Sum) and AND (Carry). It cannot handle a carry input.
- A **full adder** adds three bits (A, B, Carry-In) using two half adders plus an OR gate. It handles carry chains.
- A **ripple-carry adder** chains n full adders to add n-bit numbers. A carry-lookahead adder computes carries in parallel for speed.
- **Subtraction** is implemented by inverting B and setting Cin=1, turning the adder into a subtractor.
- A **MUX** selects one of n inputs based on a select signal. A **DEMUX** routes one input to one of n outputs.
- A **decoder** converts an n-bit binary code to one of 2ⁿ active output lines.
- An **encoder** converts one of 2ⁿ active inputs to an n-bit binary code.
- A **comparator** checks equality or magnitude relationships between two numbers.
- An **ALU** combines these building blocks: an adder/subtractor, logic operation circuits, and a MUX to select the desired result — forming the computing heart of a CPU.

---

## Exercises

### Easy

1. Trace through the full adder for the input A=1, B=0, Cin=1. What are Sum and Cout?

2. Draw a 4-to-1 MUX gate circuit. It has inputs I0, I1, I2, I3 and select lines S1, S0. When should I3 be selected?

3. A 3-to-8 decoder has 3 input bits and 8 output lines. When the input is 101 (binary for 5), which output line is active?

### Medium

4. Design an 8-bit ripple-carry adder. How many full adders, AND gates, XOR gates, and OR gates does it contain in total?

5. Build a 4-bit comparator that outputs 1 when two 4-bit numbers A and B are equal, 0 otherwise. How many gates does it use?

6. A 4-bit ALU with 4 operations (add, subtract, AND, OR) uses a 2-bit opcode. How would you extend it to support 8 operations? What changes in the circuit?

### Hard

7. A carry-lookahead adder computes all carries simultaneously. For a 4-bit adder with inputs A[3:0], B[3:0], and Cin:
   - Generate G[i] = A[i] AND B[i] for each bit
   - Propagate P[i] = A[i] XOR B[i] for each bit
   - Carry C[1] = G[0] OR (P[0] AND Cin)
   - Carry C[2] = G[1] OR (P[1] AND G[0]) OR (P[1] AND P[0] AND Cin)
   - Write the formula for C[3] and C[4].
   - How many gate delays does this take compared to a ripple-carry adder?

8. **Build a 4-bit ALU in logic simulator**: Use a free tool like Logisim Evolution (open-source) or Digital (open-source) to build a 4-bit ALU that supports: ADD, SUB, AND, OR, XOR, NOT, SHL (shift left), SHR (shift right) — 8 operations total, needing a 3-bit opcode. Screenshot your working design with test inputs.
