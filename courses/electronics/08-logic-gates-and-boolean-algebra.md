# Chapter 08: Logic Gates and Boolean Algebra

> **"Boolean algebra is the mathematics of TRUE and FALSE, 1 and 0, ON and OFF. Every program ever written, every computation ever done, reduces to these simple rules. And each rule is implemented by just a few transistors."**

---

## Table of Contents
1. [Number Systems](#1-number-systems)
2. [Binary Arithmetic](#2-binary-arithmetic)
3. [Two's Complement (Negative Numbers)](#3-twos-complement-negative-numbers)
4. [Other Codes](#4-other-codes)
5. [Boolean Algebra](#5-boolean-algebra)
6. [Logic Gates — Each in Full Detail](#6-logic-gates--each-in-full-detail)
7. [Logic Gate Families](#7-logic-gate-families)
8. [Karnaugh Maps (K-Maps)](#8-karnaugh-maps-k-maps)
9. [Logic Minimization](#9-logic-minimization)
10. [Summary](#10-summary)

---

## 1. Number Systems

### Decimal (Base-10) — Our everyday system

```
Digits: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9

Place values (powers of 10):
... 10³  10²  10¹  10⁰
... 1000  100   10    1

Example: 4726 = 4×1000 + 7×100 + 2×10 + 6×1
```

### Binary (Base-2) — Computer language

```
Digits: 0, 1  (only two!)

Place values (powers of 2):
... 2⁷  2⁶  2⁵  2⁴  2³  2²  2¹  2⁰
... 128  64   32   16   8   4   2   1

Example: 1011₂ = 1×8 + 0×4 + 1×2 + 1×1 = 11₁₀

Converting Decimal to Binary (divide by 2 method):
  23 ÷ 2 = 11 R 1 ← LSB (Least Significant Bit)
  11 ÷ 2 = 5  R 1
   5 ÷ 2 = 2  R 1
   2 ÷ 2 = 1  R 0
   1 ÷ 2 = 0  R 1 ← MSB (Most Significant Bit)

Read remainders bottom-up: 23₁₀ = 10111₂

Verify: 16+0+4+2+1 = 23 ✓
```

**Binary terminology:**
- **Bit:** single binary digit (0 or 1)
- **Nibble:** 4 bits (e.g., 1010)
- **Byte:** 8 bits (e.g., 10110110) — 0 to 255 range
- **Word:** processor's natural data size (16, 32, or 64 bits depending on CPU)
- **Kilobyte (KB):** 2¹⁰ = 1024 bytes (often approximated as 1000 in storage)
- **Megabyte (MB):** 2²⁰ = 1,048,576 bytes
- **Gigabyte (GB):** 2³⁰ = 1,073,741,824 bytes

### Hexadecimal (Base-16) — Convenient for binary

```
Digits: 0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F

A=10, B=11, C=12, D=13, E=14, F=15

Each hex digit represents exactly 4 binary bits (1 nibble)!

Hex table:
Dec | Hex | Binary
  0 |  0  | 0000
  1 |  1  | 0001
  2 |  2  | 0010
  3 |  3  | 0011
  4 |  4  | 0100
  5 |  5  | 0101
  6 |  6  | 0110
  7 |  7  | 0111
  8 |  8  | 1000
  9 |  9  | 1001
 10 |  A  | 1010
 11 |  B  | 1011
 12 |  C  | 1100
 13 |  D  | 1101
 14 |  E  | 1110
 15 |  F  | 1111

Example: 0xAF = 1010 1111₂ = 10×16 + 15 = 175₁₀

Converting hex to binary: replace each hex digit with its 4-bit binary:
  0xFF = 1111 1111 = 255 (all 8 bits set)
  0x1F = 0001 1111 = 31
  0xDEAD = 1101 1110 1010 1101₂ = 57,005₁₀ (common hex "word")
```

**Prefix notation:**
- Binary: 1010₂ or 0b1010
- Hex: 0xFF or FFh or 0xFF
- Octal: 017₈ or 017

### Octal (Base-8)

```
Digits: 0-7
Each octal digit = 3 binary bits
Example: 0177₈ = 001 111 111₂ = 127₁₀
Unix file permissions: chmod 755 = 111 101 101₂ (rwxr-xr-x)
```

---

## 2. Binary Arithmetic

### Addition

```
Binary addition rules:
  0+0 = 0
  0+1 = 1
  1+0 = 1
  1+1 = 10  (0, carry 1)
  1+1+1 = 11 (1, carry 1)

Example: 1011 + 0110
    1011   (11)
  + 0110   ( 6)
  ──────
   10001   (17) ← 16+1=17 ✓

Shows carry propagation:
  Column 0: 1+0 = 1, no carry
  Column 1: 1+1 = 0, carry 1
  Column 2: 0+1+1(carry) = 0, carry 1
  Column 3: 1+0+1(carry) = 0, carry 1
  Column 4: 0+0+1(carry) = 1
  Result: 10001
```

### Subtraction

```
Binary subtraction rules (borrowing):
  0-0 = 0
  1-0 = 1
  1-1 = 0
  0-1 = borrow from next column (10-1=1, and next column reduced by 1)

Example: 1010 - 0110
    1010   (10)
  - 0110   ( 6)
  ──────
    0100   ( 4) ← correct!
```

### Multiplication

```
Same as decimal but simpler (only multiply by 0 or 1):
  1011 × 1010

  1011 × 0 = 0000    (0 × LSB)
  1011 × 1 = 1011    (1 × next bit, shift 1 left)
  1011 × 0 = 0000    (0 × next bit, shift 2 left)
  1011 × 1 = 1011    (1 × MSB, shift 3 left)

  Sum: 0000000
       0101100
      00011000
    + 10110000
    ──────────
      1101110₂ = 64+32+8+4+2 = 110₁₀  ✓ (11×10=110)
```

---

## 3. Two's Complement (Negative Numbers)

**Problem:** How to represent negative numbers in binary?

**Two's complement** is the universal method used in all modern computers.

### How It Works

For an **n-bit** number:
- **Positive numbers:** 0 to 2^(n-1)-1 (normal binary)
- **Negative numbers:** represent -x as (2^n - x)

**Range of n-bit two's complement:**
```
Most negative: -2^(n-1)
Most positive: +2^(n-1) - 1

For 8-bit:  -128 to +127
For 16-bit: -32,768 to +32,767
For 32-bit: -2,147,483,648 to +2,147,483,647
For 64-bit: -9,223,372,036,854,775,808 to +9,223,372,036,854,775,807
```

### Converting Decimal to Two's Complement

Method: **Invert all bits, then add 1**

```
Example: -7 in 8-bit two's complement

Step 1: Write +7 in binary:  0000 0111
Step 2: Invert all bits:     1111 1000
Step 3: Add 1:               1111 1001

So -7 = 1111 1001₂

Verification: 1111 1001
  = -128 + 64 + 32 + 16 + 8 + 0 + 0 + 1
  = -128 + 121 = -7 ✓

Or: the MSB (bit 7) has NEGATIVE weight: -2^7 = -128
    1111 1001: -128 + 64 + 32 + 16 + 8 + 1 = -7
```

### Key Property

**Subtraction = Addition of negative:**
```
10 - 6 = 10 + (-6)

+10 = 0000 1010
 -6 = 1111 1010  (two's complement of 6)

  0000 1010
+ 1111 1010
──────────
 10000 0100  ← ignore overflow bit (carry out of MSB)
=  0000 0100 = 4 ✓

This is why computers only need an ADDER to do both addition and subtraction!
```

### Sign Detection

- **MSB = 0:** positive number
- **MSB = 1:** negative number (in two's complement)

```
0111 1111 = +127 (positive, MSB=0)
1000 0000 = -128 (negative, MSB=1)
1111 1111 = -1   (negative, MSB=1)
0000 0000 = 0    (zero)
```

---

## 4. Other Codes

### BCD (Binary Coded Decimal)

```
Represent each decimal digit separately as 4-bit binary:

Example: 97₁₀
  9 = 1001
  7 = 0111
  BCD: 1001 0111

Note: Only 0-9 valid (10-15 never used, so not efficient)
But: Easy to convert to/from decimal display

Used in: calculators, 7-segment displays, digital clocks
```

### Gray Code

Adjacent values differ by exactly **one bit**:

```
Decimal | Binary | Gray
  0     |  0000  | 0000
  1     |  0001  | 0001
  2     |  0010  | 0011
  3     |  0011  | 0010
  4     |  0100  | 0110
  5     |  0101  | 0111
  6     |  0110  | 0101
  7     |  0111  | 0100
  8     |  1000  | 1100
  9     |  1001  | 1101

Binary to Gray: G[n] = B[n] ⊕ B[n+1]  (XOR adjacent bits)
Gray to Binary: B[n] = XOR of all Gray bits from MSB to position n
```

**Why Gray code?** In mechanical encoders (shaft encoders, optical encoders), only ONE bit changes at a time between adjacent positions. Prevents glitches where multiple bits might momentarily be wrong during transition.

### ASCII (American Standard Code for Information Interchange)

7-bit code (128 characters) for text:
```
'A' = 65 = 0x41 = 0100 0001
'a' = 97 = 0x61 = 0110 0001  (lowercase = uppercase + 32)
'0' = 48 = 0x30
'1' = 49 = 0x31
' ' = 32 = 0x20 (space)
'\n' = 10 = 0x0A (newline)
'\r' = 13 = 0x0D (carriage return)
```

**Extended ASCII:** 8-bit (256 characters) — includes special symbols, different per region
**Unicode/UTF-8:** modern standard, handles ALL languages (up to 4 bytes per character)

---

## 5. Boolean Algebra

**George Boole** (1854) invented a mathematics of logic using just two values: TRUE/FALSE (1/0).

**Claude Shannon** (1937) showed Boolean algebra maps perfectly to electronic switching circuits — the birth of digital electronics.

### Basic Boolean Operations

**NOT (Complement):**
```
A' or Ā or ¬A

Truth table:
A | A'
0 |  1
1 |  0
```

**AND:**
```
A · B or A AND B or AB

Truth table:
A B | AB
0 0 |  0
0 1 |  0
1 0 |  0
1 1 |  1

Result is 1 ONLY when BOTH inputs are 1
```

**OR:**
```
A + B or A OR B

Truth table:
A B | A+B
0 0 |  0
0 1 |  1
1 0 |  1
1 1 |  1

Result is 1 when AT LEAST ONE input is 1
```

### Boolean Postulates and Theorems

```
Identity:
  A + 0 = A        A · 1 = A

Null/Dominance:
  A + 1 = 1        A · 0 = 0

Idempotent:
  A + A = A        A · A = A

Complement:
  A + A' = 1       A · A' = 0

Double complement:
  (A')' = A

Commutative:
  A + B = B + A    A · B = B · A

Associative:
  (A+B)+C = A+(B+C)    (AB)C = A(BC)

Distributive:
  A(B+C) = AB+AC    A+BC = (A+B)(A+C)

Absorption:
  A + AB = A       A(A+B) = A

Consensus:
  AB + A'C + BC = AB + A'C  (BC is redundant!)
```

### De Morgan's Theorems (VERY IMPORTANT)

```
Theorem 1:  (A + B)' = A' · B'
  "NOT(A OR B) = (NOT A) AND (NOT B)"
  NAND form: complement of OR = AND of complements

Theorem 2:  (A · B)' = A' + B'
  "NOT(A AND B) = (NOT A) OR (NOT B)"
  NOR form: complement of AND = OR of complements

Extended:   (A + B + C)' = A' · B' · C'
            (A · B · C)' = A' + B' + C'
```

**Practical use of De Morgan's:**
```
If you need a NOR gate but only have NAND gates:
  (A+B)' = A'·B' = NOT(A AND B with A NOT and B NOT inputs)
  So: NOR(A,B) = NAND(NOT A, NOT B) = NAND(A',B')
```

### Duality Principle

Every Boolean theorem has a **dual** — obtained by swapping AND↔OR and 0↔1:
- Any valid equation remains valid when you make this swap
- The dual of (A+0=A) is (A·1=A) — both are valid!

---

## 6. Logic Gates — Each in Full Detail

### NOT Gate (Inverter)

**Boolean:** Y = A'
**Function:** Output is always opposite of input

```
Truth Table:
A | Y
0 | 1
1 | 0

Symbol:
A ──▷○── Y  (triangle with bubble at output)

CMOS Implementation:
        VDD
         │
        [PMOS] ← gate=A
         │
         ├──── Y
         │
        [NMOS] ← gate=A
         │
        GND

When A=0: PMOS ON, NMOS OFF → Y=VDD=1
When A=1: PMOS OFF, NMOS ON → Y=GND=0

IC: 74HC04 (hex inverter, 6 NOT gates)
```

### AND Gate

**Boolean:** Y = A · B (written as AB)
**Function:** Output HIGH only when ALL inputs HIGH

```
Truth Table (2-input):
A B | Y
0 0 | 0
0 1 | 0
1 0 | 0
1 1 | 1

Symbol:
A ──|
    D── Y   (flat back, curved front)
B ──|

3-input AND: Y = A·B·C (1 only when A=B=C=1)

CMOS: NAND + Inverter (CMOS AND = NAND followed by NOT)
      Direct AND is less efficient than NAND in CMOS!

IC: 74HC08 (quad 2-input AND)
    74HC11 (triple 3-input AND)
    74HC21 (dual 4-input AND)
```

### OR Gate

**Boolean:** Y = A + B
**Function:** Output HIGH when ANY input is HIGH

```
Truth Table:
A B | Y
0 0 | 0
0 1 | 1
1 0 | 1
1 1 | 1

Symbol:
A ──|
    )── Y   (curved on both sides)
B ──|

IC: 74HC32 (quad 2-input OR)
    74HC4075 (triple 3-input OR)
```

### NAND Gate — UNIVERSAL GATE 1

**Boolean:** Y = (A · B)' = AB̄
**Function:** NOT-AND — output LOW only when ALL inputs HIGH

```
Truth Table:
A B | Y
0 0 | 1
0 1 | 1
1 0 | 1
1 1 | 0   ← only case where output is 0

Symbol:
A ──|
    D○── Y   (AND symbol with bubble at output)
B ──|

CMOS: Only 4 transistors! (2 PMOS parallel + 2 NMOS series)
      Faster and simpler than AND gate in CMOS

IC: 74HC00 (quad 2-input NAND) — MOST USED LOGIC IC
    74HC10 (triple 3-input NAND)
    74HC20 (dual 4-input NAND)

WHY IT'S UNIVERSAL — Can build any other gate from NAND only:

NOT from NAND:      A──|NAND|──(A=B)── A'
                    A──┘

AND from NAND:      A──|NAND|──|NAND|── AB
                    B──┘       (both tied together = NOT)

OR from NAND:       A──|NAND|──|NAND|──┐
                    A──┘       ↑A'      ├──|NAND|── A+B
                    B──|NAND|──|NAND|──┘
                    B──┘       ↑B'
                    (Apply De Morgan: (A'·B')' = A+B)
```

### NOR Gate — UNIVERSAL GATE 2

**Boolean:** Y = (A + B)' = (A+B)̄
**Function:** NOT-OR — output HIGH only when ALL inputs LOW

```
Truth Table:
A B | Y
0 0 | 1   ← only case where output is 1
0 1 | 0
1 0 | 0
1 1 | 0

Symbol:
A ──|
    )○── Y   (OR symbol with bubble)
B ──|

CMOS: 4 transistors (2 PMOS series + 2 NMOS parallel)

IC: 74HC02 (quad 2-input NOR)

UNIVERSAL: Can build any gate from NOR only (similar to NAND)
NOT from NOR:   A──|NOR|── A' (both inputs tied to A)
AND from NOR:   (A'+B')' = A·B using De Morgan
```

### XOR Gate (Exclusive OR)

**Boolean:** Y = A ⊕ B = A'B + AB'
**Function:** Output HIGH when inputs are DIFFERENT

```
Truth Table:
A B | Y
0 0 | 0   ← same inputs
0 1 | 1   ← different
1 0 | 1   ← different
1 1 | 0   ← same inputs

Symbol:
A ──|
    )── Y   (OR symbol with extra curved line at input)
B ──|

Properties:
  A ⊕ 0 = A        (XOR with 0 = identity)
  A ⊕ 1 = A'       (XOR with 1 = inversion)
  A ⊕ A = 0        (XOR with itself = always 0)
  A ⊕ A' = 1       (XOR with complement = always 1)
  Commutative: A⊕B = B⊕A
  Associative: (A⊕B)⊕C = A⊕(B⊕C)

Applications:
  Parity generation/checking
  Half adder (sum bit)
  Full adder
  Encryption (one-time pad: data XOR key)
  CRC (cyclic redundancy check)
  Comparators

IC: 74HC86 (quad 2-input XOR)
```

### XNOR Gate (Exclusive NOR)

**Boolean:** Y = (A ⊕ B)' = AB + A'B'
**Function:** Output HIGH when inputs are SAME (equality detector)

```
Truth Table:
A B | Y
0 0 | 1   ← same inputs
0 1 | 0
1 0 | 0
1 1 | 1   ← same inputs

Applications:
  Equality detection
  Error checking
  Digital comparators
```

---

## 7. Logic Gate Families

### TTL (Transistor-Transistor Logic) — 74xx series

- **Supply voltage:** VCC = 5V ±5%
- **Logic levels:**
  - Output LOW: 0 to 0.4V (spec), but typically < 0.2V
  - Output HIGH: 2.4V to VCC (spec), typically ~3.4V
  - Input threshold: ~1.3V
- **Speed:** propagation delay ~10ns (standard TTL)
- **Power:** ~10mW per gate (static)
- **Fan-out:** 10 (one output drives 10 inputs)
- Subfamilies:
  - 74: Standard TTL (original, 1960s)
  - 74S: Schottky TTL (faster, ~3ns)
  - 74LS: Low-power Schottky (~9mW, ~10ns) — very popular
  - 74AS: Advanced Schottky (very fast, 1.5ns)
  - 74ALS: Advanced Low-power Schottky

### CMOS — 74HCxx, 74HCTxx series

- **Supply voltage:** 2V to 6V (HC) or 4.5-5.5V (HCT, TTL compatible)
- **Logic levels:**
  - Output LOW: <0.1V (essentially 0V) — much better than TTL!
  - Output HIGH: >4.9V (essentially VCC)
  - Noise margin: nearly 50% of supply (vs ~30% for TTL)
- **Speed:** ~7ns (74HC), comparable to 74LS
- **Power:** near zero static, depends on frequency (C×V²×f)
- **Fan-out:** 50+ (output can drive many inputs)
- Subfamilies:
  - 74HC: High-speed CMOS, 5V supply
  - 74HCT: TTL-compatible (3.3V logic level input threshold)
  - 74AC: Advanced CMOS (faster, ~5ns)
  - 74ACT: TTL-compatible advanced CMOS
  - 74AHC: Advanced high-speed CMOS
  - 74LVC: Low voltage CMOS (1.65V to 3.6V) — for 3.3V systems
  - 74ALVC: Advanced low voltage CMOS (1.65-3.3V, very fast)

### Level Translating

When mixing 5V TTL/CMOS with 3.3V systems:
- **3.3V to 5V:** usually safe (3.3V output still recognized as HIGH by TTL threshold of 2V)
- **5V to 3.3V:** DANGEROUS — 5V can destroy 3.3V device! Use:
  - 74LVC245 (level translator)
  - BSS138 MOSFET translator
  - TXS0108 / TXS0102 (bi-directional translator)
  - Resistor divider (slow but simple)

### ECL (Emitter-Coupled Logic)

- Fastest logic family: <0.5ns propagation delay
- Differential signaling: PECL, LVPECL
- High power consumption
- Used in: 10 Gbps+ serial communications, high-speed ADCs

---

## 8. Karnaugh Maps (K-Maps)

K-maps are a visual method for simplifying Boolean expressions.

### Why K-Maps?

Boolean algebra manipulation requires skill and creativity. K-maps provide a **systematic visual** method — adjacent cells in the K-map differ by exactly 1 variable (Gray code ordering!), so grouping adjacent cells eliminates a variable.

### 2-Variable K-Map

```
         A
      0     1
    ┌─────┬─────┐
  0 │  m0 │  m2 │
B   ├─────┼─────┤
  1 │  m1 │  m3 │
    └─────┴─────┘

m0 = A'B' (A=0, B=0)
m1 = A'B  (A=0, B=1)
m2 = AB'  (A=1, B=0)
m3 = AB   (A=1, B=1)
```

### 3-Variable K-Map

```
          AB
      00   01   11   10
    ┌────┬────┬────┬────┐
  0 │ m0 │ m2 │ m6 │ m4 │
C   ├────┼────┼────┼────┤
  1 │ m1 │ m3 │ m7 │ m5 │
    └────┴────┴────┴────┘

Note the column order: 00, 01, 11, 10 (Gray code — adjacent differ by 1 bit)
```

### 4-Variable K-Map

```
          AB
      00   01   11   10
    ┌────┬────┬────┬────┐
 00 │  0 │  4 │ 12 │  8 │
    ├────┼────┼────┼────┤
CD 01 │  1 │  5 │ 13 │  9 │
    ├────┼────┼────┼────┤
 11 │  3 │  7 │ 15 │ 11 │
    ├────┼────┼────┼────┤
 10 │  2 │  6 │ 14 │ 10 │
    └────┴────┴────┴────┘
```

### K-Map Rules

1. Fill in 1s (or don't cares) for the minterms of the function
2. Form groups of 1s (powers of 2: 1, 2, 4, 8, 16...)
3. Groups can wrap around edges (the map is toroidal — top connects to bottom, left to right)
4. Each cell can be used in multiple groups
5. Choose fewest, largest groups that cover all 1s
6. Each group → one term in SOP expression (variables that don't change = include, that do change = eliminate)
7. If all 1s covered → done!

### K-Map Example (4-variable)

```
Function: F = Σm(0,2,4,6,8,10) = minterms where output = 1

K-Map:
          AB
      00   01   11   10
    ┌────┬────┬────┬────┐
 00 │  1 │  0 │  0 │  1 │
    ├────┼────┼────┼────┤
CD 01 │  0 │  0 │  0 │  0 │
    ├────┼────┼────┼────┤
 11 │  0 │  0 │  0 │  0 │
    ├────┼────┼────┼────┤
 10 │  1 │  0 │  0 │  1 │
    └────┴────┴────┴────┘

Groups:
  Group 1: All four corners (0,4,8,10 → but wait, 0,8 and 4,2 and 2,10...)

  Let's see: positions are m0(AB=00,CD=00), m2(AB=00,CD=10), m4(AB=01,CD=00), m6(AB=01,CD=10)

  Wait, let me redo this. m0=0,m2=2,m4=4,m6=6,m8=8,m10=10

  Pattern: all minterms where D=0!
  Group all 8 cells where D=0 → F = D' (!)

  Simplified: F = D' (output is 1 whenever D is 0)
  Without K-map: might write complex expression, K-map reveals the simplicity
```

### Don't Care Conditions (X)

Some input combinations might never occur (e.g., states 10-15 in a decimal counter).
Mark these as 'X' — use them in groups if they help simplify, ignore them otherwise.

---

## 9. Logic Minimization

### Sum of Products (SOP)

Express function as OR of AND terms (minterms):
```
F = A'B'C + A'BC' + ABC
```

**Canonical (full) SOP:** each minterm includes ALL variables
**Minimal SOP:** fewest terms, fewest literals (from K-map or algebra)

### Product of Sums (POS)

Express function as AND of OR terms (maxterms):
```
F = (A+B+C)(A+B'+C')(A'+B+C')
```

### Quine-McCluskey (Tabulation) Method

For many variables where K-map is impractical (>6 variables):
1. List all minterms in binary
2. Group by number of 1s
3. Combine groups that differ by exactly 1 bit (mark with dash)
4. Repeat until no more combinations possible
5. Find minimum cover using a prime implicant chart
6. Result: minimum SOP expression

---

## 10. Summary

```
NUMBER SYSTEMS
══════════════
Binary:  base 2, digits 0,1
Hex:     base 16, digits 0-F, 1 hex = 4 binary bits
BCD:     each decimal digit as 4-bit binary
Gray:    adjacent values differ by 1 bit (encoder use)

TWO'S COMPLEMENT
════════════════
Negative numbers: invert all bits + add 1
Range (n-bit): -2^(n-1) to +2^(n-1)-1
MSB=0: positive, MSB=1: negative

BOOLEAN ALGEBRA KEY THEOREMS
══════════════════════════════
De Morgan: (AB)' = A'+B',  (A+B)' = A'B'
Absorption: A+AB=A, A(A+B)=A
Complement: A+A'=1, A·A'=0

LOGIC GATES
═══════════
NOT:  Y = A'
AND:  Y = AB  (1 only if ALL inputs 1)
OR:   Y = A+B (1 if ANY input 1)
NAND: Y = (AB)' (universal — can make any gate)
NOR:  Y = (A+B)' (universal — can make any gate)
XOR:  Y = A⊕B (1 if inputs DIFFER)
XNOR: Y = (A⊕B)' (1 if inputs SAME)

GATE FAMILIES
══════════════
TTL (74LS): VCC=5V, 10mW/gate
CMOS (74HC): VCC=2-6V, near-zero static power
CMOS wins in digital circuits: low power, full swing, scalable
```

---

**← Previous:** [Chapter 07: MOSFETs and Advanced Transistors](./07-mosfet-and-advanced-transistors.md)
**→ Next:** [Chapter 09: Combinational Circuits](./09-combinational-circuits.md)
