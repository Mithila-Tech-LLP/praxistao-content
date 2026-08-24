# Chapter 04: Logic Gates — Building Decisions from Electricity

You know that a transistor is a switch. You know binary represents information. Now we combine these two ideas: how do you build a circuit that takes binary inputs and produces a binary output according to some logical rule? The answer is **logic gates** — the actual building blocks of every digital circuit ever made.

## Table of Contents

1. [What Is a Logic Gate?](#1-what-is-a-logic-gate)
2. [The NOT Gate — The Inverter](#2-the-not-gate)
3. [The AND Gate](#3-the-and-gate)
4. [The OR Gate](#4-the-or-gate)
5. [The NAND Gate — The Universal Gate](#5-the-nand-gate--the-universal-gate)
6. [The NOR Gate](#6-the-nor-gate)
7. [The XOR Gate — Exclusive Or](#7-the-xor-gate--exclusive-or)
8. [Building Gates from Transistors](#8-building-gates-from-transistors)
9. [Truth Tables and Boolean Algebra](#9-truth-tables-and-boolean-algebra)
10. [Why NAND is Universal](#10-why-nand-is-universal)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Is a Logic Gate?

A logic gate is an electronic circuit that takes one or more binary inputs and produces a single binary output, based on a logical rule.

Think of it this way: you are making a decision about whether to open an umbrella. The decision depends on inputs: Is it raining? Did I forget my umbrella? A logic gate is the hardware equivalent of such a decision.

Every logic gate has:
- **Inputs**: one or more wires carrying 0 (low voltage) or 1 (high voltage)
- **Output**: one wire that is 0 or 1 based on the inputs and the gate type
- **A rule**: the logical relationship between inputs and output

The rule for every possible combination of inputs is described in a **truth table** — a table listing every possible input combination and its corresponding output.

Here are the six fundamental gates. Everything in your computer is built from combinations of these:

---

## 2. The NOT Gate

The simplest gate: one input, one output. The output is the **opposite** of the input.

```
    Symbol:          Truth Table:
    
    A ──[>o]── Q    ┌───┬───┐
                    │ A │ Q │
                    ├───┼───┤
                    │ 0 │ 1 │
                    │ 1 │ 0 │
                    └───┴───┘
```

The NOT gate is also called an **inverter**. In Boolean algebra, NOT is written as Ā (A-bar) or ¬A or ~A.

**Real-world analogy**: A NOT gate is like a room's smoke alarm: when there is no smoke (0), the alarm is active/ready (1). When there IS smoke (1), it triggers output (0 = no alarm?). OK, that doesn't work perfectly — but the logic is: the output is always the opposite of the input.

**Implementation**: A CMOS NOT gate uses exactly one PMOS and one NMOS transistor (as we saw in Chapter 2). When input is HIGH: NMOS on, PMOS off → output connected to ground → LOW. When input is LOW: PMOS on, NMOS off → output connected to VDD → HIGH.

### Quick Check

> 1. What does a NOT gate do to its input?
> 2. If input A = 1, what is the output of a NOT gate?
> 3. Why is the NOT gate also called an "inverter"?

---

## 3. The AND Gate

Two inputs. Output is 1 **only if BOTH inputs are 1**. Otherwise output is 0.

```
    Symbol:          Truth Table:
    
    A ──┐            ┌───┬───┬───┐
        ├──[ ]── Q   │ A │ B │ Q │
    B ──┘            ├───┼───┼───┤
                     │ 0 │ 0 │ 0 │
                     │ 0 │ 1 │ 0 │
                     │ 1 │ 0 │ 0 │
                     │ 1 │ 1 │ 1 │
                     └───┴───┴───┘
```

The AND gate has a curved right side and flat left side in standard notation. In Boolean algebra: Q = A · B or Q = AB.

**Real-world analogy**: A combination lock that requires both a fingerprint AND a PIN. You get access only when BOTH conditions are true.

**Common uses**:
- Enabling/disabling a signal: Q = A AND Enable. When Enable=0, Q=0 regardless of A. When Enable=1, Q follows A.
- Checking if two conditions are both satisfied
- Masking bits: AND with a mask to clear specific bits (e.g., 11001100 AND 11110000 = 11000000 — clears the lower 4 bits)

### The AND Condition Table Mnemonic

"**A**ND requires **A**ll inputs to be 1."

### Quick Check

> 1. What is the output of an AND gate when inputs are A=1 and B=0?
> 2. If you AND any value with 0, what do you always get?
> 3. If you AND any value with 1, what do you get?

---

## 4. The OR Gate

Two inputs. Output is 1 if **at least one input is 1**.

```
    Symbol:          Truth Table:
    
    A ──┐            ┌───┬───┬───┐
        ├──D── Q     │ A │ B │ Q │
    B ──┘            ├───┼───┼───┤
                     │ 0 │ 0 │ 0 │
                     │ 0 │ 1 │ 1 │
                     │ 1 │ 0 │ 1 │
                     │ 1 │ 1 │ 1 │
                     └───┴───┴───┘
```

The OR gate has a curved right side and a curved "bow" on the left. In Boolean algebra: Q = A + B.

**Real-world analogy**: A light switch where either of two switches (at either end of a hallway) can turn the light on. If EITHER switch is in the ON position, the light is ON.

**Common uses**:
- Combining signals: turn on an alert if any of several error conditions occur
- Setting bits: OR with a mask to set specific bits to 1 (e.g., 11001100 OR 00001111 = 11001111)

### Quick Check

> 1. What is the output of an OR gate when both inputs are 0?
> 2. If you OR any value with 1, what do you always get?
> 3. Give a real-world example where OR logic is natural.

---

## 5. The NAND Gate — The Universal Gate

NAND = NOT + AND. The output is the **opposite of AND** — it is 0 only when ALL inputs are 1.

```
    Symbol:          Truth Table:
    
    A ──┐            ┌───┬───┬───┐
        ├──[ ]o─ Q   │ A │ B │ Q │
    B ──┘            ├───┼───┼───┤
                     │ 0 │ 0 │ 1 │
                     │ 0 │ 1 │ 1 │
                     │ 1 │ 0 │ 1 │
                     │ 1 │ 1 │ 0 │
                     └───┴───┴───┘
```

The small circle (bubble) at the output means "inversion" — it inverts the AND gate. In Boolean algebra: Q = ̄(AB) (NAND of A and B).

**Why NAND is special**: NAND is a **universal gate** — you can build ANY other gate using only NAND gates. This is profoundly important: chip manufacturers only need one type of gate to build everything.

**Implementation**: The CMOS NAND gate uses **2 PMOS in parallel** (output HIGH when either input is LOW) and **2 NMOS in series** (output LOW only when BOTH inputs are HIGH). This is simpler than AND (which would need NAND + inverter) — which is why real chips often use NAND rather than AND.

### Quick Check

> 1. Write the truth table for NAND and compare it to AND.
> 2. If inputs are A=1, B=1, what is the NAND output?
> 3. If inputs are A=0, B=1, what is the NAND output?

---

## 6. The NOR Gate

NOR = NOT + OR. The output is 0 if **any input is 1**. The output is 1 only when ALL inputs are 0.

```
    Symbol:          Truth Table:
    
    A ──┐            ┌───┬───┬───┐
        ├──Do── Q    │ A │ B │ Q │
    B ──┘            ├───┼───┼───┤
                     │ 0 │ 0 │ 1 │
                     │ 0 │ 1 │ 0 │
                     │ 1 │ 0 │ 0 │
                     │ 1 │ 1 │ 0 │
                     └───┴───┴───┘
```

In Boolean algebra: Q = ̄(A + B).

Like NAND, **NOR is also a universal gate** — you can build any circuit from NOR gates alone. The first commercially popular integrated circuits (like the Apollo Guidance Computer) were built entirely from NOR gates.

### Quick Check

> 1. Under what input conditions is NOR output 1?
> 2. NOR is sometimes called the "all-must-be-zero gate." Why?

---

## 7. The XOR Gate — Exclusive Or

XOR = "Exclusive OR." Output is 1 if **exactly one input is 1** (inputs are different). Output is 0 if both inputs are the same.

```
    Symbol:          Truth Table:
    
    A ──┐            ┌───┬───┬───┐
        ├──⊕── Q     │ A │ B │ Q │
    B ──┘            ├───┼───┼───┤
                     │ 0 │ 0 │ 0 │
                     │ 0 │ 1 │ 1 │
                     │ 1 │ 0 │ 1 │
                     │ 1 │ 1 │ 0 │
                     └───┴───┴───┘
```

In Boolean algebra: Q = A ⊕ B.

**Real-world analogy**: A "debate" gate. Output is 1 (interesting) when the two people disagree (A≠B). Output is 0 (boring) when they agree (A=B).

**Critical uses in computing**:
- **Binary addition**: XOR gives the sum bit when adding two 1-bit numbers (1+1=0 carry 1, which is exactly what XOR gives — we'll see this in Chapter 5)
- **Comparison**: A XOR B = 0 exactly when A = B (useful for equality checking)
- **Cryptography**: XOR is the core operation of many encryption schemes — XOR with a key to encrypt, XOR again to decrypt
- **Error detection**: XOR all the bytes in a message → if any byte changes, the XOR checksum changes

**XNOR**: The complement of XOR. Output is 1 when inputs are the same (equality test). Sometimes called the "equivalence gate."

### Quick Check

> 1. What is the output of XOR when A=1, B=1?
> 2. If A XOR B = 0, what can you conclude about A and B?
> 3. Why is XOR useful for addition? (Hint: what is 1+1 in binary without the carry?)

---

## 8. Building Gates from Transistors

Let us see exactly how CMOS transistors form these gates.

### CMOS NOT Gate

```
VDD (1V)
 │
[P]  ← PMOS: ON when input LOW, OFF when input HIGH
 │
─┼── Output
 │
[N]  ← NMOS: OFF when input LOW, ON when input HIGH
 │
GND (0V)

Input LOW (0): PMOS=ON, NMOS=OFF → output connected to VDD → Output HIGH (1)
Input HIGH (1): PMOS=OFF, NMOS=ON → output connected to GND → Output LOW (0)
```

### CMOS NAND Gate (2 PMOS parallel, 2 NMOS series)

```
VDD
 ├────────[P-A]────────┐
 └────────[P-B]────────┤
                        │── Output
               [N-A]──┤
                       │
               [N-B]──┘
                       │
                      GND

When A=0 OR B=0: at least one PMOS is ON → output connected to VDD → HIGH (1)
When A=1 AND B=1: both PMOS OFF, both NMOS ON → output connected to GND → LOW (0)
This is exactly the NAND truth table.
```

### CMOS NOR Gate (2 PMOS series, 2 NMOS parallel)

```
VDD
 │
[P-A]
 │
[P-B]  ← series: BOTH must be ON (both inputs LOW) for output to be HIGH
 │──── Output
 ├────[N-A]────┐
 └────[N-B]────┘── GND

When A=0 AND B=0: both PMOS ON → output HIGH
When A=1 OR B=1: at least one NMOS ON, at least one PMOS OFF → output LOW
This is exactly the NOR truth table.
```

### Transistor Count

| Gate | PMOS transistors | NMOS transistors | Total |
|------|-----------------|-----------------|-------|
| NOT  | 1 | 1 | 2 |
| NAND | 2 | 2 | 4 |
| NOR  | 2 | 2 | 4 |
| AND  | 3 | 3 | 6 (NAND + NOT) |
| OR   | 3 | 3 | 6 (NOR + NOT) |
| XOR  | ~8-12 | ~8-12 | ~16-24 |

This is why chip designers prefer NAND/NOR over AND/OR — they use fewer transistors.

### Quick Check

> 1. How many transistors does a CMOS NAND gate use?
> 2. Why do chip designers prefer NAND gates over AND gates in their implementations?
> 3. In the CMOS NAND gate, why are the PMOS transistors in parallel (not series)?

---

## 9. Truth Tables and Boolean Algebra

Logic gates follow the rules of **Boolean algebra** — a mathematical system developed by George Boole in 1854, long before electronic computers. Boolean algebra has two values (true/false, 1/0) and three basic operations (AND, OR, NOT).

### Key Laws of Boolean Algebra

These laws let you simplify complex circuits:

**Identity laws:**
- A AND 1 = A
- A OR 0 = A

**Null laws:**
- A AND 0 = 0
- A OR 1 = 1

**Idempotent laws:**
- A AND A = A
- A OR A = A

**Complement laws:**
- A AND (NOT A) = 0
- A OR (NOT A) = 1

**Commutative laws:**
- A AND B = B AND A
- A OR B = B OR A

**Associative laws:**
- A AND (B AND C) = (A AND B) AND C
- A OR (B OR C) = (A OR B) OR C

**Distributive laws:**
- A AND (B OR C) = (A AND B) OR (A AND C)
- A OR (B AND C) = (A OR B) AND (A OR C)

**De Morgan's Laws** — among the most useful:
- NOT (A AND B) = (NOT A) OR (NOT B)
- NOT (A OR B) = (NOT A) AND (NOT B)

De Morgan's Laws convert between NAND and NOR, which is powerful for circuit simplification.

### Boolean Algebra in Programming

Boolean algebra is not just for hardware — it is in every program you write:

```
// In Python:
if is_admin and not is_suspended:      # AND and NOT
    grant_access()

if has_error or time_exceeded:         # OR
    abort_operation()

# XOR check (are they different?)
if (user_a_online) != (user_b_online): # XOR equivalent
    notify_mismatch()
```

The logic gates inside the CPU implement these same operations in hardware.

### Quick Check

> 1. Simplify using Boolean algebra: A AND (A OR B)
> 2. Apply De Morgan's Law: NOT (A AND B) = ?
> 3. In Python, what does `True and False or True` evaluate to? Trace through using AND/OR precedence.

---

## 10. Why NAND Is Universal

A **universal gate** is one from which you can build ANY other gate. NAND and NOR are both universal.

Here is how to build all other gates from NAND:

### NOT from NAND
Connect both inputs of a NAND to the same signal A:
```
NAND(A, A) = NOT(A AND A) = NOT(A) = A̅
```

### AND from NAND
NOT of NAND:
```
AND(A,B) = NOT(NAND(A,B)) = NAND(NAND(A,B), NAND(A,B))
```

### OR from NAND (De Morgan!)
```
OR(A,B) = NOT(NOT(A) AND NOT(B))   [De Morgan]
        = NAND(NOT(A), NOT(B))
        = NAND(NAND(A,A), NAND(B,B))
```

This means a chip foundry only needs to perfect one cell (the NAND gate) and can implement everything else from it. In practice, standard cell libraries contain dozens of optimized gate types for efficiency, but the universality of NAND means you are never stuck.

### Historical Significance: The Apollo Computer

The Apollo Guidance Computer (AGC) that took humans to the Moon in 1969 was built using approximately 4,100 NOR gates — and nothing else. Every circuit in the computer — adders, memory control, timing — was built from NOR gates alone. This was a practical manufacturing decision: using only one component type simplified production, testing, and spare parts.

### Quick Check

> 1. Build an OR gate using only NAND gates. Draw the circuit.
> 2. Why is it practically useful that NAND is a universal gate?
> 3. The Apollo Guidance Computer used only NOR gates. Why might a designer choose to use only one type of gate even if it makes circuits less efficient?

---

## Summary

- A **logic gate** takes binary inputs and produces a binary output according to a logical rule.
- The six fundamental gates: **NOT** (inverts input), **AND** (1 only if all inputs are 1), **OR** (0 only if all inputs are 0), **NAND** (NOT of AND), **NOR** (NOT of OR), **XOR** (1 when inputs differ).
- Every gate is fully described by its **truth table**.
- CMOS gates are built from PMOS (network connecting output to VDD) and NMOS (network connecting output to GND) transistors. NAND uses 4 transistors; NOR uses 4 transistors; NOT uses 2.
- **Boolean algebra** provides rules for simplifying logic expressions. De Morgan's Laws convert AND↔OR under inversion.
- **NAND and NOR are universal gates** — you can build any logic function from either one alone.
- Everything in your computer — arithmetic, memory control, decisions — is built from combinations of these gates.

---

## Exercises

### Easy

1. Fill in the output column for each of the six gates when A=1 and B=0:
   - NOT(A) = ?
   - AND(A,B) = ?
   - OR(A,B) = ?
   - NAND(A,B) = ?
   - NOR(A,B) = ?
   - XOR(A,B) = ?

2. Write the truth table for a 3-input AND gate (inputs A, B, C). How many rows does it have?

3. Apply De Morgan's Law to simplify NOT(A OR B). What gate is this equivalent to?

### Medium

4. A security system has three sensors: door (D), window (W), and motion (M). The alarm should sound if:
   - The door OR the window is open, AND
   - The motion sensor is triggered.
   Write the Boolean expression for the alarm. Then draw a gate circuit implementing it.

5. Show step by step how to build an XOR gate using only NAND gates. (Hint: XOR(A,B) = (A AND NOT B) OR (NOT A AND B). Now replace each AND and OR with NAND equivalents.)

6. Using Boolean algebra, prove that XOR(A,B) = (A OR B) AND NOT(A AND B). (Hint: build truth tables for both sides and show they match.)

### Hard

7. A 2-to-1 multiplexer (MUX) selects one of two input signals based on a select signal:
   - When S=0, output Q = input A
   - When S=1, output Q = input B
   
   Write the Boolean expression for Q in terms of A, B, and S. Draw a gate circuit. Then implement it using only NAND gates.

8. The parity bit is a simple error detection technique. Compute the XOR of all bits in a byte; if it is 1, append a 0 parity bit; if it is 0, append a 1 parity bit. The receiver XORs all received bits including the parity bit — if the result is not 0, an error occurred.
   (a) Compute the parity bit for the byte 10110101.
   (b) Draw a gate circuit to compute the XOR of an 8-bit number (chain of XOR gates).
   (c) What type of errors can parity detect? What type of errors can it NOT detect?
