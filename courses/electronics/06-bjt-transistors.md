# Chapter 06: BJT Transistors — Bipolar Junction Transistors

> **"The transistor is the greatest invention of the 20th century. A small current controlling a large current — this simple principle underlies every amplifier, every digital chip, every radio, and every computer ever made."**

---

## Table of Contents
1. [History and Importance](#1-history-and-importance)
2. [BJT Structure](#2-bjt-structure)
3. [How BJT Works — Physical Operation](#3-how-bjt-works--physical-operation)
4. [Current Relationships](#4-current-relationships)
5. [Operating Regions](#5-operating-regions)
6. [BJT Characteristics Curves](#6-bjt-characteristics-curves)
7. [Three Configurations](#7-three-configurations)
8. [DC Biasing Circuits](#8-dc-biasing-circuits)
9. [BJT as a Switch](#9-bjt-as-a-switch)
10. [BJT as an Amplifier](#10-bjt-as-an-amplifier)
11. [Small Signal Model](#11-small-signal-model)
12. [Darlington Pair](#12-darlington-pair)
13. [PNP Transistor](#13-pnp-transistor)
14. [Common BJT Packages and Types](#14-common-bjt-packages-and-types)
15. [Summary](#15-summary)

---

## 1. History and Importance

**December 23, 1947** — Walter Brattain and John Bardeen at Bell Labs demonstrated the first working point-contact transistor, under the direction of William Shockley.

- Replaced vacuum tubes: smaller, lower power, more reliable, much cheaper
- 1956: Nobel Prize in Physics for Bardeen, Brattain, Shockley
- 1954: Texas Instruments makes first commercial silicon transistor
- 1958: First integrated circuit (Jack Kilby at TI) — multiple transistors on one chip
- By 2024: >10²² (10 sextillion) transistors manufactured per year globally

**The transistor enabled:**
- Portable radios (1954) → walkman → smartphones
- Mainframes → minicomputers → personal computers → laptops
- Every digital device ever made

---

## 2. BJT Structure

A **Bipolar Junction Transistor** has **three** semiconductor regions:

### NPN Transistor

```
     Collector (C)
          │
     ─────┴─────
    │     N     │  ← Collector: N-type, lightly doped, wide
    ├───────────┤  ← BC junction (normally reverse biased)
    │     P     │  ← Base: P-type, VERY THIN, lightly doped
    ├───────────┤  ← BE junction (normally forward biased)
    │     N     │  ← Emitter: N-type, HEAVILY doped
     ─────┬─────
          │
       Emitter (E)
          │
        Base (B) ←── (connected to thin P region)
```

### PNP Transistor

```
     Collector (C)
          │
     ─────┴─────
    │     P     │  ← Collector: P-type
    ├───────────┤
    │     N     │  ← Base: N-type, thin
    ├───────────┤
    │     P     │  ← Emitter: P-type, heavily doped
     ─────┬─────
          │
       Emitter (E)
```

### Circuit Symbols

**NPN:**
```
            C
            │
         ───┤
   B ───────┤   (arrow on emitter points OUT of base = NPN)
         ───┤
            │
            E
```

**PNP:**
```
            C
            │
         ───┤
   B ───────┤   (arrow on emitter points IN toward base = PNP)
         ───┤
            │
            E
```

**Memory tip:** "NPN = Not Pointing iN, PNP = Pointing iN Proudly" or
"Arrow points iN = PNP, arrow points Not iN = NPN"

### Key Physical Features

1. **Emitter:** Heavily doped (N+ for NPN) — large supply of majority carriers
2. **Base:** Very thin (~1μm in modern devices, was ~25μm in 1950s) and lightly doped
3. **Collector:** Moderately doped, physically larger area, handles more power

**Why the thin base is critical:**
- Electrons injected from emitter must travel across base
- If base is thin, most electrons make it to collector before recombining
- If base is thick, electrons recombine in base → only base current, no collector current → just a diode, not a transistor!

---

## 3. How BJT Works — Physical Operation

Let's trace what happens in an NPN transistor when properly biased (active region):

### Step-by-Step Physics

**Step 1: Forward bias the Base-Emitter (BE) junction**
```
Apply VBE ≈ +0.7V (forward bias)
→ BE junction turns on
→ Electrons injected from N+ emitter into thin P base
```

**Step 2: Electrons diffuse across thin base**
```
Base is P-type: holes are majority carriers
Injected electrons from emitter are MINORITY carriers in base
Electric field tries to pull them to base terminal
But base is SO THIN that most electrons diffuse across before recombining

Fraction recombining in base = determined by base width and minority carrier lifetime
Typical: only 1-5% of emitter electrons recombine in base (become base current)
95-99% make it through!
```

**Step 3: Electrons swept into collector**
```
The reverse-biased BC junction creates a strong electric field
This field SWEEPS electrons arriving at BC junction into collector
Even though BC is reverse biased, electrons (minority carriers in base-side of BC junction)
are swept in by the field — this is why collector current flows despite reverse bias!
```

**Step 4: Current amplification**
```
IB = recombination current in base (small — only the electrons that didn't make it)
IC = collector current (large — all the electrons that made it through)

IC/IB = β (current gain) = 20 to 500 typically

So: small IB controls large IC
    This is AMPLIFICATION!
```

### Energy Band Diagram

```
NPN in active region energy bands:

Emitter (N)  |   Base (P)   |  Collector (N)
─────────────┼──────────────┼───────────────
  CB         │              │    CB
 ─────       │              │   ─────
             │              │
  EF(E)─── →──→→→─ electrons tunnel/diffuse across → ─→→─ EF(C)
             │   (thin base: ~1μm)                │
  ─────       │              │   ─────
  VB         │              │    VB
```

---

## 4. Current Relationships

### Fundamental Equations

```
KCL at transistor: IE = IB + IC

(Emitter current = Base current + Collector current)
```

**Current gain β (also called hFE):**
```
β = IC/IB

Typical values:
  Small signal BJTs: β = 100-500
  Power BJTs: β = 20-100
  Darlington: β = β₁ × β₂ (can be 10,000+)
```

**Current gain α:**
```
α = IC/IE = β/(β+1)

Relationship:
  β = α/(1-α)
  α = β/(β+1)

For β = 100: α = 100/101 = 0.99 (99% of emitter current reaches collector)
For β = 300: α = 300/301 = 0.997 (99.7% efficiency)
```

**All currents from β:**
```
Given β and one current, find all others:

IC = β × IB
IE = IC + IB = β×IB + IB = (β+1)×IB
IB = IC/β = IE/(β+1)
```

**Example:**
```
β = 100, IB = 50μA

IC = 100 × 50μA = 5mA
IE = 5mA + 50μA = 5.05mA
```

---

## 5. Operating Regions

A BJT has four operating regions determined by the bias conditions of both junctions:

### Region 1: Cutoff (OFF state)

```
Conditions: VBE < 0.7V (BE junction not forward biased)
            VBC < 0 (BC junction reverse biased)

Result: No base current → IC ≈ 0 → transistor is OFF (open switch)
        Only tiny leakage current ICEO flows

Applications: Digital logic, switching circuits (when OFF)
```

### Region 2: Active (Amplification)

```
Conditions: VBE ≈ 0.7V (BE forward biased)
            VCE > 0.2V (BC reverse biased)
            VBC = VBE - VCE < 0

Result: IC = β × IB
        Linear amplification
        VCE can vary without changing IC much

Applications: Amplifiers, linear circuits
```

### Region 3: Saturation (ON state)

```
Conditions: VBE ≈ 0.7V (BE forward biased)
            VBC > 0 (BC forward biased too!)
            VCE(sat) ≈ 0.1-0.2V

Result: Transistor is fully ON (closed switch)
        IC is limited by EXTERNAL circuit, NOT by β×IB
        Base is "flooded" with excess minority carriers (charge storage)
        Need IB >> IC/β to ensure saturation

Applications: Digital logic, switching circuits (when ON)
```

### Region 4: Reverse Active (rarely used)

```
Emitter and collector swapped roles
Very low β in this region (typically 0.1-5)
Generally avoid this region
```

### Summary Table

| Region | VBE | VBC | VCE | IC | State |
|--------|-----|-----|-----|-----|-------|
| Cutoff | <0.5V | <0.5V | VCC | ≈0 | OFF |
| Active | ≈0.7V | <0 | >0.2V | β×IB | Amplifying |
| Saturation | ≈0.7V | ≈0.5-0.7V | ≈0.2V | circuit-limited | ON |

---

## 6. BJT Characteristics Curves

### Output Characteristics (Common Emitter: IC vs VCE)

```
IC (mA)
│
8├──────────────────────────── IB = 80μA
 │                             (β=100 → IC = 8mA)
6├──────────────────────────── IB = 60μA
 │
4├──────────────────────────── IB = 40μA
 │
2├──────────────────────────── IB = 20μA
 │
0├──────────────────────────── IB = 0μA
 ├────────────────────────────────── VCE (V)
  0.2  1   2   3   4   5  ...

Regions:
  0 to 0.2V: saturation (curves rising steeply)
  0.2V onwards: active (flat curves — IC independent of VCE)
  Breakdown voltage: where curves turn up sharply (BVCEO)

Early effect: curves slightly slope upward in active region
  IC = IS × e^(VBE/VT) × (1 + VCE/VA)
  VA = Early voltage (50-200V for BJTs)
```

### Input Characteristics (IB vs VBE)

```
IB (μA)
│
80│                    /
   │                  /
40│                 /
   │               /   ← exponential (like diode!)
 5│            /
   │────────────────────── VBE (V)
    0    0.5  0.6  0.7
```

- Same as forward-biased diode: IB = IS × e^(VBE/VT)
- Threshold ≈ 0.6-0.7V

---

## 7. Three Configurations

The BJT can be connected three ways — each with different properties:

### 1. Common Emitter (CE) — Most Widely Used

```
          VCC
           │
          [RC]
           │
           ├──────── Output (Vout)
           │
      ─────┤ C
Vin ──[RB]─┤ B    ← Input
      ─────┤ E
           │
          GND

Characteristics:
  Voltage gain: Av = -gm×RC = -RC/re (NEGATIVE — 180° phase inversion!)
  Current gain: Ai ≈ β
  Input impedance: Zi = RB || β×re  (moderate)
  Output impedance: Zo = RC  (moderate)
  Power gain: highest of the three
```

**Small signal analysis:**
```
re = VT/IC = 26mV/IC  (intrinsic emitter resistance)

For IC = 1mA: re = 26Ω
For IC = 10mA: re = 2.6Ω

Voltage gain: Av = -RC/re
Example: RC = 5kΩ, IC = 1mA → re = 26Ω → Av = -5000/26 = -192 V/V
```

### 2. Common Base (CB)

```
Input at emitter, output at collector, base grounded (for AC):

       VCC
        │
       [RC]
        │
        ├──── Output
        │
       [C]
       [B]──── GND (via bypass cap for AC)
Vin ──[E]
        │
      [-Vee]

Characteristics:
  Voltage gain: Av ≈ RC/re (positive, high)
  Current gain: Ai = α ≈ 1 (less than 1, no current gain)
  Input impedance: Zi = re  (very low, ~26Ω)
  Output impedance: Zo = RC  (high)
  Used at: high frequencies (no Miller effect)
```

### 3. Common Collector (CC) / Emitter Follower

```
          VCC
           │
      ─────┤ C
Vin ──[RB]─┤ B
      ─────┤ E
           │
          [RE]
           │
          GND
           │
           └── Output (Vout = Vin - 0.7V approximately)

Characteristics:
  Voltage gain: Av ≈ 1 (output follows input — "follower")
  Current gain: Ai ≈ β+1 (high)
  Input impedance: Zi = RB || (β+1)(RE+re)  (HIGH — good buffer!)
  Output impedance: Zo = RE || (RB+rs)/(β+1)  (LOW)
  No phase inversion
  Used as: buffer, impedance transformer, current amplifier
```

**Applications of emitter follower:**
- Impedance transformation: high-impedance source → low-impedance drive
- Drive speakers, motors, LEDs (without large current into logic pin)
- Class A audio output stages

---

## 8. DC Biasing Circuits

Biasing sets the **quiescent (Q) point** — the DC operating point before any signal is applied.

### Method 1: Fixed Bias (Base Resistor Bias)

```
      VCC
       │
      [RB]
       │
   B───┤
   C──[RC]──VCC
   E───┤
       │
      GND

Analysis:
  IB = (VCC - VBE) / RB = (VCC - 0.7) / RB
  IC = β × IB
  VCE = VCC - IC × RC

Problem: Very unstable! If β varies (device-to-device, temperature):
  IC = β × (VCC-VBE)/RB → IC depends strongly on β
  Temperature change: VBE decreases 2mV/°C → IB increases → IC increases
  This can cause "thermal runaway" and destroy transistor
```

### Method 2: Emitter Bias (Two Supplies)

```
      +VCC
       │
      [RC]
       │
   C───┤ NPN B ──[RB]──0V
   E───┤
       │
      [RE]
       │
      -VEE

Analysis:
  IB = (-VEE - (-VBE)) / (RB + (β+1)RE) ≈ (VEE-0.7) / ((β+1)RE)
  Better stability than fixed bias
```

### Method 3: Voltage Divider Bias (Best — Most Common!)

```
      VCC
       │
      [R1]
       │
       ├────[B] (NPN)
      [R2]   └─[C]──[RC]──VCC
       │       [E]──[RE]──GND
      GND

Analysis (standard approach):
  VB = VCC × R2/(R1+R2)      ← Thevenin equivalent
  VE = VB - VBE ≈ VB - 0.7V
  IE = VE/RE ≈ IC
  VC = VCC - IC×RC
  VCE = VC - VE

Stability criterion: if VB >> VBE (VB is large and stable)
  Then VE = VB - 0.7V ≈ stable
  IE ≈ IC = VE/RE ≈ stable
  → IC doesn't depend on β! (self-stabilizing)

Design rule: Make VE ≈ 10% of VCC (e.g., 0.5V for 5V supply)
             Make VB ≈ VE + 0.7V ≈ 1.2V (for 5V supply)
             Choose R1, R2 for VB, with R1||R2 = ~10×β×RE for stability
```

### Method 4: Collector-to-Base Bias (Feedback Bias)

```
      VCC
       │
      [RC]
       │
       ├──[RB]──[B] NPN
       │          [E]──GND
       │          [C]──(same node as RC junction above)

Analysis:
  IB = (VCC - VBE - VCE) / RB  (feedback: if IC increases, VCE drops, so IB decreases)
  Self-correcting — negative feedback
  IC × RC = VCC - VCE - IB×RB ≈ VCC - IC/β × RB - VCE

Less stable than voltage divider but simple
```

---

## 9. BJT as a Switch

This is the most important application in digital electronics.

### Switching Between Cutoff (OFF) and Saturation (ON)

**OFF state (Cutoff):**
```
VBE = 0V (or negative for PNP)
IB = 0 → IC ≈ 0
VCE = VCC (full supply voltage appears across transistor)
Output HIGH (if using VCE as output)
```

**ON state (Saturation):**
```
VBE = 0.7V (forward biased)
IC(sat) = (VCC - VCE(sat)) / RC ≈ VCC/RC
VCE(sat) ≈ 0.1-0.2V (small residual voltage)
Output LOW
```

**Required base current for saturation:**
```
Must satisfy: IB > IC(sat) / β

Rule: Use 5-10× more IB than minimum for guaranteed saturation:
  IB = IC(sat) / βmin × overdrive_factor (5-10)

Base resistor:
  RB = (Vin - VBE) / IB = (Vin - 0.7) / IB

Example: Switch a 100mA load with VCC=12V, RC=120Ω, β_min=50, Vin=5V
  IC(sat) = (12-0.2)/120 ≈ 98mA
  IB(min) = 98mA/50 = 1.96mA
  IB(design) = 1.96 × 5 = 9.8mA ≈ 10mA
  RB = (5-0.7)/10mA = 430Ω → use 470Ω
```

### Switching Speed

```
Turn-ON time: td (delay) + tr (rise time)
Turn-OFF time: ts (storage) + tf (fall time)

Storage time (ts) is the big problem:
  Saturated transistor has excess minority charge stored in base
  Must be removed before transistor turns off
  ts can be 1-10 μs for regular BJTs

Solutions to speed up switching:
1. Baker clamp: Schottky diode prevents deep saturation
   (keeps VBC < 0.4V so transistor never fully saturates)
2. Schottky transistors (TTL "S" family): built-in Baker clamp
3. Use MOSFET instead (no minority carrier storage)
```

### Application: Relay Driver

```
      VCC (relay coil supply)
       │
      [Relay Coil]
       │         ← freewheeling diode (D1) in parallel with coil (cathode to VCC)
      Collector
   B ──[RB]── GPIO (microcontroller output)
      Emitter
       │
      GND

When GPIO = HIGH (3.3V or 5V):
  IB flows → transistor saturates → relay coil energized → relay closes
When GPIO = LOW:
  IB = 0 → transistor off → relay de-energized → relay opens
  D1 clamps inductive kickback from coil (protects transistor)

For Arduino 5V output driving 12V relay with 200mA coil:
  Transistor: 2N2222 (600mA, β=75 typical)
  IC = 200mA, β = 75
  IB needed = 200mA/75 × 5 (overdrive) = 13.3mA
  RB = (5-0.7)/13.3mA = 323Ω → use 330Ω
```

---

## 10. BJT as an Amplifier

### Single-Stage Common Emitter Amplifier

```
        VCC (+12V)
         │
        [RC] 5.6kΩ
         │
         ├──────────────[Cout]────── Output
         │
   [C] NPN
   [B] ──────────────────── ←── Input (through Cin)
   [E]
    │
   [RE] 560Ω
   [CE] ← bypass capacitor (shorts RE for AC)
    │
   GND

R1 = 100kΩ ┬── VB
R2 = 22kΩ  ┴── GND   (voltage divider bias)
```

**Design for IC = 1mA, VCC = 12V:**
```
1. Choose VE = 1.2V → VCE(Q) = 12-1.2-IC×RC = comfortable midpoint
2. RE = VE/IC = 1.2/0.001 = 1.2kΩ → use 1.2kΩ
3. Choose VCE = 5V (leave room for signal swing)
   VRC = VCC - VCE - VE = 12-5-1.2 = 5.8V
   RC = VRC/IC = 5.8/0.001 = 5.8kΩ → use 5.6kΩ
4. VB = VE + VBE = 1.2 + 0.7 = 1.9V
   R1/R2 ratio: VB = VCC×R2/(R1+R2) → R2/(R1+R2) = 1.9/12 = 0.158
   Choose IB×10 = 1mA/100 × 10 = 100μA flows through divider
   R1+R2 = VCC/Idivider = 12/0.1mA = 120kΩ
   R2 = 0.158 × 120kΩ = 19kΩ → use 18kΩ
   R1 = 120-18 = 102kΩ → use 100kΩ

5. AC voltage gain (CE bypassed):
   re = 26mV/1mA = 26Ω
   Av = -RC/re = -5600/26 = -215 V/V (215x amplification!)
```

### Class A, B, AB, C Amplifiers

- **Class A:** transistor conducts full 360° of cycle. Low efficiency (~25%), low distortion. Audio preamplifiers.
- **Class B:** transistor conducts 180° (half cycle). Each of two transistors handles one half. Crossover distortion. ~78% efficiency.
- **Class AB:** each transistor conducts slightly more than 180°. Reduces crossover distortion. ~60-70% efficiency. Most audio power amplifiers.
- **Class C:** transistor conducts <180°. Very efficient (>90%). High distortion. Only for tuned RF amplifiers.
- **Class D:** transistor switches ON/OFF (PWM). >95% efficiency. Digital audio amplifiers (Class D switching amplifiers).

---

## 11. Small Signal Model

The **hybrid-π model** is used for small signal (AC) analysis around the Q-point:

```
Small signal equivalent circuit of BJT:

         rπ          Cπ
B ──────┤ ├──────┬──────────────── C
        │        │    ro
        vπ      gm×vπ ↑ (current source)
        │        │
E ──────────────────────────────── E

Parameters:
  gm = transconductance = IC/VT = IC/26mV
  rπ = β/gm = β×VT/IC = β×re
  ro = output resistance = VA/IC  (Early voltage VA, typically 50-200V)
  Cπ = base-emitter capacitance (diffusion + junction)
  Cμ = base-collector capacitance (junction, Miller effect)
```

**Example at IC = 2mA, β = 100, VA = 100V:**
```
gm = 2mA/26mV = 76.9 mA/V = 0.077 A/V
rπ = 100/0.077 = 1299Ω ≈ 1.3kΩ
ro = 100/0.002 = 50kΩ
```

### Frequency Response

The gain falls at high frequencies due to Cπ and Cμ:

**Unity gain bandwidth (fT):**
```
fT = gm/(2π(Cπ+Cμ))

At fT: |current gain| = 1 (0 dB)
Typical: fT = 100MHz to 10GHz for modern BJTs

fβ (3dB cutoff of current gain) = fT/β
```

---

## 12. Darlington Pair

Two transistors connected so that the first amplifies the base current for the second.

```
           C2 ──── Collector
      ┌────┤
      │    E1 ──── B2
B1 ───┤    ├───── E2 ──── Emitter
           │
        R_pull-down (optional, to speed up turn-off)
```

**Effective parameters:**
```
β_total = β₁ × β₂  (can be 10,000 to 500,000!)
VBE_total = VBE1 + VBE2 = 1.2-1.4V (two BE junctions in series)
VCE(sat) ≈ 0.9-1.2V (higher than single transistor)
```

**Use case:** Drive high-current loads from very weak control signals
- Example: ULN2003 Darlington array IC: 7 Darlington pairs, 500mA each, common used to drive stepper motors, relays from microcontroller GPIO

---

## 13. PNP Transistor

PNP is the complement of NPN — everything is reversed:

```
Differences from NPN:
  Arrow on emitter points INWARD (toward base)
  VEB ≈ +0.7V (emitter more positive than base)
  VEC = collector voltage (negative of NPN's VCE)
  Current arrows reversed

Biasing: Emitter connected to VCC, collector to load, base pulls current out

PNP circuit analysis rules:
  IB flows OUT of base (conventional current)
  IC flows OUT of collector (into load)
  IE = IC + IB (flows into emitter from VCC)
  IC = β × IB (same relationship, different directions)
```

**PNP switch (high-side switch):**
```
VCC ── Emitter
       [PNP]
       Base ──[RB]── Control (GPIO)
       Collector ──[Load]── GND

When GPIO = LOW (0V): VEB = VCC - 0 = VCC > 0.7V → ON
When GPIO = HIGH: VEB too small → OFF

Used when load needs to be connected to VCC side
```

**Complementary symmetry:** NPN+PNP pairs used in class AB push-pull amplifiers, CMOS (NMOS+PMOS).

---

## 14. Common BJT Packages and Types

### Through-Hole Packages
- **TO-92:** 3-pin plastic, for small signal BJTs (2N2222, BC547, 2N3904)
- **TO-18:** metal can, for microwave/RF BJTs
- **TO-220:** large plastic, for power transistors (TIP31, TIP41C, MJE3055)
- **TO-3:** metal diamond, for high-power transistors (2N3055)

### Surface Mount (SMD)
- **SOT-23:** 3-pin, for small signal (MMBT2222, BC817)
- **SOT-223, D-PAK, D²-PAK:** power transistors in SMD

### Common BJTs and Their Specs

| Part | Type | β | VCE max | IC max | fT | Package |
|------|------|---|---------|--------|-----|---------|
| 2N3904 | NPN | 100-300 | 40V | 200mA | 300MHz | TO-92 |
| 2N3906 | PNP | 100-300 | 40V | 200mA | 250MHz | TO-92 |
| 2N2222 | NPN | 75-300 | 40V | 600mA | 300MHz | TO-18/TO-92 |
| BC547 | NPN | 110-800 | 45V | 100mA | 300MHz | TO-92 |
| BC557 | PNP | 125-800 | 45V | 100mA | 150MHz | TO-92 |
| TIP31C | NPN | 25-50 | 100V | 3A | 3MHz | TO-220 |
| TIP32C | PNP | 25-50 | 100V | 3A | 3MHz | TO-220 |
| 2N3055 | NPN | 20-70 | 60V | 15A | 0.8MHz | TO-3 |
| MJE3055 | NPN | 20-70 | 60V | 10A | 2MHz | TO-220 |
| BC337 | NPN | 100-630 | 45V | 800mA | 100MHz | TO-92 |

### Selection Guide

For a **switching** application:
- Choose BJT with: IC(max) > 2× load current, VCEO > 2× supply, β > 20 (enough for driving)

For an **amplifying** application:
- Choose based on: frequency (need fT >> signal freq), noise figure (for RF), VCE range, power dissipation

---

## 15. Summary

```
BJT FUNDAMENTALS
════════════════

Structure:
  NPN: Emitter(N+) | Base(P, thin) | Collector(N)
  PNP: Emitter(P+) | Base(N, thin) | Collector(P)

Currents:
  IE = IB + IC
  IC = β × IB
  α = IC/IE = β/(β+1)

Operating Regions:
  Cutoff:     VBE < 0.7V → IC ≈ 0 (OFF)
  Active:     VBE = 0.7V, VCE > 0.2V → IC = β×IB (amplifier)
  Saturation: Both junctions forward → VCE ≈ 0.2V (ON)

As Switch:
  ON:  IB = VIN/RB, ensure IB > IC/β (overdrive 5-10×)
  OFF: VBE = 0V

As Amplifier (CE):
  Av = -RC/re, re = 26mV/IC
  gm = IC/VT, rπ = β/gm

Key Formula:
  VBE = 0.7V (Si NPN, active)
  IC = β × IB
  re = 26mV/IC
  Av = -RC/re (CE amp, bypassed RE)
```

---

**← Previous:** [Chapter 05: Basic Electronic Components](./05-basic-electronic-components.md)
**→ Next:** [Chapter 07: MOSFETs and Advanced Transistors](./07-mosfet-and-advanced-transistors.md)
