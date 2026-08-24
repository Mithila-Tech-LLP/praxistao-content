# Chapter 52: Kirchhoff's Laws

> **"Kirchhoff's laws are just conservation of energy and conservation of charge applied to circuits. Once you understand that, they're not rules to memorize — they're truths about the universe."**

---

## Table of Contents

- [52.1 Why We Need Kirchhoff's Laws](#521-why-we-need-kirchhoffs-laws)
- [52.2 Kirchhoff's Current Law (KCL)](#522-kirchhoffs-current-law-kcl)
- [52.3 Kirchhoff's Voltage Law (KVL)](#523-kirchhoffs-voltage-law-kvl)
- [52.4 Setting Up Equations](#524-setting-up-equations)
- [52.5 Solving Multi-Loop Circuits](#525-solving-multi-loop-circuits)
- [52.6 Branch Currents Method](#526-branch-currents-method)
- [52.7 Worked Examples](#527-worked-examples)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 52.1 Why We Need Kirchhoff's Laws

Simple circuits with one battery and resistors in series or parallel can be solved with Ohm's Law and basic series/parallel rules.

But what about this?

```
    R₁        R₂
+---/\/\/---+---/\/\/---+
|           |           |
|           R₃          |
|           |           |
+--- V₁ ---+--- V₂ ---+

Two batteries, three resistors, multiple loops.
```

Series/parallel rules don't apply — the circuit has multiple independent loops. Kirchhoff's Laws (KCL and KVL) solve ANY circuit, no matter how complex.

---

## 52.2 Kirchhoff's Current Law (KCL)

**KCL (Kirchhoff's Current Law):** The sum of all currents entering a node equals the sum of all currents leaving the node.

Or equivalently: the algebraic sum of all currents at a node = 0.

```
PHYSICAL BASIS: Conservation of charge
  Charge cannot accumulate at a node (steady state).
  Whatever current flows in must flow out.
```

### Visual Example

```
NODE A:
          I₁ (2A in)
          ↓
     -----A-----
     |         |
  I₂↓         ↑I₃
  (out)        (out)
  1.5A         ?

KCL: I₁ = I₂ + I₃
  2 = 1.5 + I₃
  I₃ = 0.5A
```

### Convention

Taking currents INTO the node as positive:

```
Sum of all currents (with sign) = 0

At node A: I₁ - I₂ - I₃ = 0
           I₁ = I₂ + I₃
```

---

## 52.3 Kirchhoff's Voltage Law (KVL)

**KVL (Kirchhoff's Voltage Law):** The sum of all voltage changes around any closed loop equals zero.

```
PHYSICAL BASIS: Conservation of energy
  If you travel around a loop and return to the start,
  you must end at the same potential you started with.
  All gains = all losses.
```

### Voltage Changes: Signs

When traversing a loop, assign signs to voltage changes:

```
THROUGH A RESISTOR:
  If you traverse in the direction of current:     voltage DROPS (-)
  If you traverse against the current:             voltage RISES (+)
  
THROUGH A BATTERY:
  If you traverse from - to + terminal:            voltage RISES (+)
  If you traverse from + to - terminal:            voltage DROPS (-)
```

```
VISUAL:

  Traversing loop:    +--- V₀ (battery) ---+---/\/\R/\/\---+
                      |    (+ to -)        |   (in direction |
                      |    DROPS (-V₀)     |    of current)  |
                      |                    |    DROPS (-IR)  |
                      +--------------------+                 |
                      |                                      |
                      +--------------------------------------+
  
  KVL: -V₀ + IR = 0  →  V₀ = IR  (Ohm's Law!)
```

---

## 52.4 Setting Up Equations

**Step-by-step process:**

1. **Label currents**: assign a symbol (I₁, I₂, etc.) and direction to each branch current. If you guess wrong, the answer will be negative (which is fine — it just means the current flows the other way).

2. **Apply KCL** at each independent node to get equations relating currents.

3. **Choose loops**: pick independent closed loops. The number of loops needed = number of unknown currents minus number of KCL equations.

4. **Apply KVL** around each chosen loop: sum voltage rises and drops = 0.

5. **Solve** the system of simultaneous equations.

### Counting Equations Needed

For a circuit with:
- N nodes: gives (N-1) independent KCL equations
- B branches (with unknown currents): need B unknowns
- Need B - (N-1) KVL equations

---

## 52.5 Solving Multi-Loop Circuits

### Example Circuit

```
        R₁ = 4Ω        R₂ = 2Ω
   +---/\/\/---+---/\/\/---+
   |           |           |
V₁=12V        R₃=6Ω     V₂=6V
   |           |           |
   +-----------+-----------+
               A
```

**Step 1: Label currents**

```
        R₁ = 4Ω  I₁→      R₂ = 2Ω  I₂→
   +---/\/\/---+---/\/\/---+
   |           |           |
V₁=12V        R₃=6Ω     V₂=6V
   |           |↓ I₃       |
   +-----------A-----------+
```

**Step 2: KCL at node A**

Current in = current out:
I₁ = I₂ + I₃ ... (we'll set up properly below)

Actually, let's define: I₁ flows left-to-right through R₁ (top left), I₂ flows right-to-left through R₂ (top right), and I₃ flows downward through R₃.

At node A: I₁ - I₂ - I₃ = 0
→ I₁ = I₂ + I₃ ... (1)

**Step 3: KVL around left loop (clockwise)**

Starting at bottom-left, going clockwise:
+V₁ - I₁R₁ - I₃R₃ = 0
+12 - 4I₁ - 6I₃ = 0 ... (2)

**Step 4: KVL around right loop (clockwise)**

Starting at bottom-right, going clockwise:
+I₃R₃ - I₂R₂ - V₂ = 0
+6I₃ - 2I₂ - 6 = 0 ... (3)

**Step 5: Solve**

From (1): I₂ = I₁ - I₃

Substitute into (3): 6I₃ - 2(I₁ - I₃) - 6 = 0
6I₃ - 2I₁ + 2I₃ - 6 = 0
-2I₁ + 8I₃ = 6 ... (3')

From (2): 4I₁ + 6I₃ = 12 → 2I₁ + 3I₃ = 6 ... (2')

Adding (2') and (3'): 11I₃ = 12
I₃ = 12/11 ≈ **1.09 A**

From (2'): I₁ = (6 - 3×12/11)/2 = (6 - 36/11)/2 = (66/11 - 36/11)/2 = (30/11)/2 = 15/11 ≈ **1.36 A**

I₂ = I₁ - I₃ = 15/11 - 12/11 = 3/11 ≈ **0.27 A**

---

## 52.6 Branch Currents Method

An alternative systematic approach:

1. Assign a current variable to each branch
2. Write KCL equations for all nodes except one
3. Write KVL equations for all independent loops
4. Solve the system

**Tip for loop selection:** Choose loops that together cover every branch at least once. Meshes (smallest possible loops, with no interior branches) are a natural choice.

---

## 52.7 Worked Examples

### Worked Example 52.1 — Simple Two-Battery Circuit

```
    12V battery        6Ω resistor
+---[+|battery|-]------/\/\/------+
|                                 |
|    4Ω resistor      9V battery  |
+-------/\/\/------[+|battery|-]--+
```

More formally:

```
   E₁=12V    R₁=6Ω
+----||--------/\/\/----+
|                       |
+----/\/\/----||--------+
   R₂=4Ω    E₂=9V
```

Label current I flowing clockwise (one loop).

KVL clockwise: +E₁ - IR₁ - E₂ - IR₂ = 0
12 - 6I - 9 - 4I = 0
3 = 10I
I = **0.3 A** clockwise

Voltage across R₁ = IR₁ = 0.3 × 6 = **1.8 V**
Voltage across R₂ = IR₂ = 0.3 × 4 = **1.2 V**

Check: 12 - 1.8 - 9 - 1.2 = 0 ✓

### Worked Example 52.2 — Node Voltage

```
        I₁         I₂
   +---→---A---→---+
   |       |       |
  6V     3Ω I₃↓  12V
   |       |       |
   +-------+-------+
```

At node A, apply KCL: I₁ + I₂ - I₃ = 0 (I₁ and I₂ in, I₃ out)

Wait, let's use node voltage method instead.

Let V_A = potential at node A (bottom = 0 V reference).

I₁ = (6 - V_A) / (assume R₁ = 2Ω) ... let's use specific values.

**The key insight**: set one node as reference (ground, 0 V), write KCL in terms of node voltages.

---

## Summary

- **KCL**: sum of currents into any node = sum out; based on conservation of charge
- **KVL**: sum of voltage changes around any closed loop = 0; based on conservation of energy
- **Voltage sign rules**: 
  - Battery traversed +→-: drops voltage (-)
  - Battery traversed -→+: gains voltage (+)
  - Resistor traversed in current direction: drops voltage (-IR)
  - Resistor traversed against current: gains voltage (+IR)
- **Procedure**: label currents → write KCL at nodes → write KVL around loops → solve system
- Number of KCL equations: N-1 (where N = number of nodes)
- Number of KVL equations: B-(N-1) (where B = branches)
- If a current turns out negative, it means it actually flows opposite to your assumed direction

---

## Key Equations

```
Kirchhoff's Current Law (KCL):
  Sum of currents into node = sum out
  Algebraically: Σ I = 0  at any node

Kirchhoff's Voltage Law (KVL):
  Sum of voltage changes around any closed loop = 0
  Algebraically: Σ V = 0  around any loop

Sign conventions for KVL:
  Resistor (with current):     ΔV = -IR  (voltage drop)
  Resistor (against current):  ΔV = +IR  (voltage rise)
  Battery (- to +):            ΔV = +ε   (voltage rise)
  Battery (+ to -):            ΔV = -ε   (voltage drop)

Independent equations needed:
  For B branches (unknown currents):
    KCL gives (N-1) equations
    KVL gives B-(N-1) equations
    Total B equations for B unknowns
```
