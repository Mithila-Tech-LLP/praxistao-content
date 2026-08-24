# Chapter 04: PN Junction and Diodes

> **"The PN junction is the most fundamental structure in all of semiconductor electronics. Every transistor, every solar cell, every LED, every integrated circuit — all contain PN junctions."**

---

## Table of Contents
1. [PN Junction Formation](#1-pn-junction-formation)
2. [Depletion Region and Built-in Potential](#2-depletion-region-and-built-in-potential)
3. [Forward Bias](#3-forward-bias)
4. [Reverse Bias](#4-reverse-bias)
5. [Shockley Diode Equation](#5-shockley-diode-equation)
6. [Diode Characteristics](#6-diode-characteristics)
7. [Diode Types in Detail](#7-diode-types-in-detail)
8. [Rectifier Circuits](#8-rectifier-circuits)
9. [Clipper and Clamper Circuits](#9-clipper-and-clamper-circuits)
10. [Diode Specifications](#10-diode-specifications)
11. [Summary](#11-summary)

---

## 1. PN Junction Formation

A **PN junction** is formed when a P-type semiconductor is placed in contact with an N-type semiconductor. This can be done by:
- Diffusing dopants into one side of a crystal
- Epitaxially growing one type on the other
- Ion implantation

```
BEFORE contact:

P-type side:          |    N-type side:
  ⊖ ○ ○ ⊖ ○          |      ⊕ ● ● ⊕ ●
  ○ ⊖ ○ ○ ⊖          |      ● ⊕ ● ● ⊕
  ⊖ ○ ○ ⊖ ○          |      ● ⊕ ● ⊕ ●

P-side symbols:
  ⊖ = fixed acceptor ion (negative, e.g. B⁻ — cannot move)
  ○ = mobile hole (positive, majority carrier in P)

N-side symbols:
  ⊕ = fixed donor ion (positive, e.g. P⁺ — cannot move)
  ● = mobile electron (negative, majority carrier in N)

(Each side is overall neutral: the mobile carriers exactly balance the
 fixed dopant ions.)
```

### What Happens at the Junction?

When P and N are joined, **concentration gradients drive diffusion**:
1. **Holes** from P-side diffuse into N-side (high hole concentration → low)
2. **Electrons** from N-side diffuse into P-side (high electron concentration → low)

As they cross and recombine:
- P-side near junction: holes leave → exposes negatively charged acceptor ions (B⁻)
- N-side near junction: electrons leave → exposes positively charged donor ions (P⁺)
- This creates a region of **immobile (fixed) charges** = **depletion region**

```
AFTER contact — PN Junction at equilibrium:

P-type  | Depletion Region |  N-type
        |                  |
⊕⊕⊕⊕   |⊖⊖⊖⊖|⊕⊕⊕⊕|   ⊕⊕⊕⊕
⊕⊕⊕⊕   |⊖⊖⊖⊖|⊕⊕⊕⊕|   ⊕⊕⊕⊕
⊕⊕⊕⊕   |⊖⊖⊖⊖|⊕⊕⊕⊕|   ⊕⊕⊕⊕
         |← xp →|← xn →|

⊖ = exposed acceptor ions (negative, on P-side of depletion)
⊕ = exposed donor ions (positive, on N-side of depletion)
```

---

## 2. Depletion Region and Built-in Potential

### The Space Charge Region

The **depletion region** (also called space charge region or SCR):
- Contains NO mobile carriers (they've all recombined or been swept away)
- Contains ONLY fixed ionized atoms
- P-side of SCR: negative charge (from B⁻ acceptor ions)
- N-side of SCR: positive charge (from P⁺ donor ions)

This charge distribution creates an **electric field** pointing from N to P (opposing further diffusion):

```
Electric field E:

P-side  ←←←←←← E-field ←←←←←←  N-side
                (from + to −)
```

### Built-in Potential (Contact Potential)

The electric field creates a **built-in potential barrier** V₀ (also called Vbi or contact potential) that opposes further diffusion:

```
         kT     NA × ND
V₀ = ───── × ln(────────)
          q       ni²

At 300K (kT/q = 26mV):
  For Si with NA = ND = 10¹⁶ cm⁻³:
  V₀ = 0.026 × ln((10¹⁶ × 10¹⁶)/(1.5×10¹⁰)²)
  V₀ = 0.026 × ln(10³²/2.25×10²⁰)
  V₀ = 0.026 × ln(4.44×10¹¹)
  V₀ = 0.026 × 26.8
  V₀ ≈ 0.696 V ≈ 0.7 V

This is why silicon diodes have ~0.7V forward drop!

For Ge (ni = 2.4×10¹³):
  V₀ ≈ 0.3V → Ge diodes have ~0.3V forward drop

For GaAs (ni = 1.8×10⁶):
  V₀ ≈ 1.2V → GaAs diodes higher forward drop
```

### Depletion Width

The depletion region extends on both sides of the junction:

```
Total width: W = xp + xn

Where:
  xp = penetration into P-side = W × ND/(NA + ND)
  xn = penetration into N-side = W × NA/(NA + ND)

Total W = √(2ε₀εr(NA + ND)V₀ / (q × NA × ND))

Where:
  ε₀ = 8.854×10⁻¹² F/m (permittivity of free space)
  εr = 11.7 for Silicon
  V₀ = built-in potential

The depletion region extends MORE into the lightly doped side!
(More lightly doped → fewer ions → must deplete more to get same charge)
```

**Example:** N-side 100× more doped than P-side
- xp/xn = ND/NA = 100
- Depletion extends mostly into P-side

---

## 3. Forward Bias

**Forward bias:** applying + voltage to P-side, − voltage to N-side.

```
External circuit:

   (+) ──── P | N ──── (−)
   Battery
```

**What happens:**
1. External voltage **reduces** the barrier: Vbarrier = V₀ − VF
2. With lower barrier, many more carriers have enough energy to cross
3. Holes injected from P into N side (minority carriers in N)
4. Electrons injected from N into P side (minority carriers in P)
5. These minority carriers diffuse across and recombine
6. Current flows!

**Key thresholds:**
- **Silicon:** significant current starts at VF ≈ **0.6-0.7V**
- **Germanium:** VF ≈ **0.25-0.3V**
- **Schottky:** VF ≈ **0.2-0.4V**
- **LED (Red):** VF ≈ **1.8-2.2V**
- **LED (Green):** VF ≈ **2.0-2.5V**
- **LED (Blue/White):** VF ≈ **3.0-3.5V**

**Current increases exponentially with voltage:**
- At 0.5V: I ≈ 10μA (small)
- At 0.6V: I ≈ 1mA (significant — ~100× increase for 0.1V!)
- At 0.7V: I ≈ 100mA (large)
- Each ~60mV increase → 10× more current (decade/60mV rule)

---

## 4. Reverse Bias

**Reverse bias:** applying − voltage to P-side, + voltage to N-side.

```
External circuit:

   (−) ──── P | N ──── (+)
   Battery
```

**What happens:**
1. External voltage **increases** the barrier: Vbarrier = V₀ + VR
2. Majority carriers are pulled AWAY from junction → depletion region widens
3. No majority carriers can cross the junction
4. Only tiny **reverse saturation current (IS)** flows — from thermally generated minority carriers
5. IS is very small: nA to μA range
6. IS essentially independent of reverse voltage (hence "saturation")

**Widened depletion region in reverse bias:**
```
W(VR) = W₀ × √(1 + VR/V₀)

Wider depletion → larger capacitance decrease (used in varactor diodes!)
```

### Breakdown

At sufficiently high reverse voltage, **breakdown** occurs:

**1. Zener Breakdown (< ~5V in Si):**
- High electric field directly breaks covalent bonds
- Quantum tunneling of electrons through thin depletion region
- Occurs in heavily doped junctions (thin depletion region)
- Sharp, well-defined breakdown voltage
- Used deliberately in Zener diodes!

**2. Avalanche Breakdown (> ~7V in Si):**
- Thermally generated carrier gains enough energy from field to ionize another atom
- Chain reaction: 1 → 2 → 4 → 8 → 16 carriers (avalanche!)
- Occurs in lightly doped junctions (wide depletion, high field)
- Also sharp breakdown voltage but more temperature sensitive

**Note:** 5-7V range can be either mechanism. Below 5V = mostly Zener. Above 7V = mostly avalanche.

---

## 5. Shockley Diode Equation

William Shockley (1950) derived the fundamental equation for diode current:

```
I = IS × (e^(V/nVT) - 1)

Where:
  I  = diode current (A)
  IS = reverse saturation current (A), typically 10⁻¹² to 10⁻⁶ A
  V  = voltage across diode (V) — positive = forward bias
  n  = ideality factor (emission coefficient)
       n = 1: ideal diode (diffusion current dominant)
       n = 2: recombination current dominant (real devices at low forward voltage)
       n = 1 to 2 in practice
  VT = thermal voltage = kT/q = 26mV at 300K

```

**Physical meaning of IS:**
```
IS = Aqni² × (Dn/ND×Ln + Dp/NA×Lp)

Where:
  A  = junction area (cm²)
  Dn,Dp = diffusion coefficients
  Ln,Lp = diffusion lengths
  ND,NA = doping concentrations

IS ∝ ni² → IS doubles every ~5-6°C temperature increase!
    → At 60°C: IS ≈ 1000× greater than at 20°C
    → Diode current is very temperature sensitive
```

**Approximations:**

For forward bias (V >> VT ≈ 26mV):
```
I ≈ IS × e^(V/nVT)    (dominant term)
```

For reverse bias (V << -VT):
```
I ≈ -IS    (reverse saturation — essentially constant)
```

**Temperature effect on forward voltage:**
```
dVF/dT ≈ -2 mV/°C    (for Si at constant current)

So:
  At 0°C (273K):   VF ≈ 0.7 + 25 × 0.002 ≈ 0.75V
  At 25°C (298K):  VF ≈ 0.7V
  At 100°C (373K): VF ≈ 0.7 - 75 × 0.002 ≈ 0.55V
```

### Dynamic Resistance

The small-signal resistance of a diode at operating point:
```
        nVT    n × 26mV
rd = ─────── = ──────────
        ID         ID

Example: diode carrying 26mA, n=1:
  rd = 26mV/26mA = 1Ω

Example: diode carrying 1mA, n=1:
  rd = 26mV/1mA = 26Ω
```

---

## 6. Diode Characteristics

### I-V Curve

```
         I (mA)
         │
    100  │                              /
         │                            /
     50  │                          /
         │                        /
         │                      /   ← Forward bias region
      1  │─────────────────── /
         │                  /
──────────●────────────────/──────────── V (volts)
    -50  -40  -30   -20  -10   0  0.7
         │     ← Reverse bias     │
    -IS  │────────────────────────│       (tiny, ~μA)
         │                        │
      Breakdown
      voltage VBR
```

**Important points on I-V curve:**
- **Knee voltage / cut-in voltage:** ~0.7V for Si — below this, almost no current
- **Forward voltage drop (VF):** voltage across conducting diode, ~0.7V at rated current
- **Reverse saturation current (IS):** tiny current in reverse, ~nA
- **Breakdown voltage (VBR):** reverse voltage at which breakdown occurs

---

## 7. Diode Types in Detail

### 1. General Purpose Rectifier Diode

**Examples:** 1N4001, 1N4007, 1N5400, 1N5408

| Parameter | 1N4001 | 1N4007 |
|-----------|--------|--------|
| Peak Inverse Voltage (PIV) | 50V | 1000V |
| Average forward current IF(av) | 1A | 1A |
| Peak forward surge current | 30A | 30A |
| Forward voltage VF | 1.1V @ 1A | 1.1V @ 1A |
| Reverse current IR | 5μA | 5μA |
| Package | DO-41 | DO-41 |

**Applications:** Power supply rectification, voltage clipping, reverse polarity protection

### 2. Zener Diode

Designed to operate in **reverse breakdown** at a precise, stable voltage.

**Operation:**
- Reverse biased
- At Zener voltage VZ: breakdown occurs, current flows
- Voltage stays essentially constant despite current changes
- VZ ranges from 1.8V to >200V

**Zener voltage regulation circuit:**
```
            R (series resistor)
VIN ────[ R ]────┬──── VOUT
                 │
               [ZD]  (Zener diode, cathode up)
                 │
               GND

Design formulas:
  VOUT = VZ (regulated)
  IZ = (VIN - VZ) / R - IL

  For regulation: IZ must stay > IZmin (minimum Zener current, ~5mA)

  R = (VIN(min) - VZ) / (IZ(min) + IL(max))

  Power dissipation of Zener: PZ = VZ × IZ(max)
```

**Example:** Regulate 5V from 9V supply, load draws 50mA
```
  VZ = 5V Zener (1N4733A, 5.1V)
  IZ(min) = 5mA (minimum for regulation)
  IL = 50mA

  R = (9 - 5) / (0.005 + 0.050) = 4/0.055 = 72.7Ω → use 68Ω

  IZ(max) when IL=0: IZ = (9-5)/68 = 58.8mA
  PZ(max) = 5 × 58.8mA = 294mW → need >500mW Zener
```

**Zener vs Avalanche breakdown:**
- Below ~5V: Zener mechanism, **negative temperature coefficient** (VZ decreases with temp)
- Above ~7V: avalanche mechanism, **positive temperature coefficient** (VZ increases with temp)
- ~5-6V: nearly zero temperature coefficient (very stable) — used in precision references

### 3. Schottky Diode (Metal-Semiconductor Junction)

- Junction between metal (Au, Pt, Ni) and N-type Si
- **No minority carrier storage** → extremely fast switching
- **Very low forward voltage:** 0.15-0.45V (vs 0.7V for PN diode)
- Higher reverse leakage current than PN diodes
- Lower reverse breakdown voltage (typically < 100V)

**Why no minority carrier storage?**
- Only majority carrier (electron) device
- No holes injected into metal, no diffusion storage
- Reverse recovery time ≈ **nanoseconds** (vs microseconds for PN)

**Applications:**
- High-frequency rectification (switching power supplies)
- Protection diodes (Schottky on Arduino inputs)
- Preventing BJT from saturating (Baker clamp)
- Detector diodes for RF circuits
- Solar cells (contact diode)

**Popular Schottkys:** 1N5817-5819 (1A, 20-40V), BAT54 (SMD), SS14 (1A, 40V SMA)

### 4. LED (Light Emitting Diode)

When forward biased, electron-hole recombination emits photons (light).

**Requires direct bandgap** semiconductor — in indirect gap materials (Si, Ge) recombination is mostly non-radiative (heat).

**Wavelength from bandgap:**
```
         hc      1240 nm·eV
λ = ────────── = ────────────
          Eg          Eg(eV)
```

**LED materials and colors:**

| Color | Wavelength | Material | VF |
|-------|-----------|---------|-----|
| Infrared | 850-950nm | GaAs, AlGaAs | 1.2-1.6V |
| Red | 620-750nm | GaAsP, AlGaAs | 1.8-2.2V |
| Orange | 590-620nm | GaAsP | 2.0-2.3V |
| Yellow | 570-590nm | GaP:N | 2.1-2.4V |
| Green (old) | 520-570nm | GaP:N | 2.1-2.5V |
| Green (modern) | 520-570nm | InGaN | 2.9-3.5V |
| Blue | 450-500nm | InGaN/GaN | 3.0-3.5V |
| White | broadband | Blue LED + phosphor | 3.0-3.5V |
| UV | 365-400nm | GaN, AlGaN | 3.5-4.5V |

**LED current limiting resistor:**
```
         VCC - VF
RLED = ────────────
            IF

Where IF = desired LED current (typically 5-20mA)

Example: Red LED with VF=2V, IF=10mA, VCC=5V:
  R = (5-2)/0.01 = 300Ω → use 270Ω or 330Ω standard value
```

**LED efficiency:**
- Luminous efficacy: lumens per watt
- Early LEDs: ~1 lm/W
- Modern high-power LEDs: 150-220 lm/W
- Compare: incandescent bulb ~15 lm/W, fluorescent ~70 lm/W

### 5. Photodiode

A **reverse-biased diode** that generates current when illuminated.

**Operation:**
- Light photons with E > Eg create electron-hole pairs in depletion region
- Electric field separates e-h pairs → current flows
- Photocurrent: `IP = R × P_optical`
  - R = responsivity (A/W), typically 0.4-0.8 A/W for Si at 850nm
  - P = optical power (Watts)

**Two operating modes:**

*Photoconductive mode (reverse biased):*
- Faster response (wider depletion region)
- Small dark current (leakage)
- Used for: high-speed photodetection, fiber optics

*Photovoltaic mode (no bias = 0V):*
- Very low dark current
- Used for: solar cells, precision measurement
- Solar cell equation: I = IL - IS(e^(V/VT) - 1) where IL = light-generated current

### 6. Varactor Diode (Varicap)

Exploits the **voltage-dependent capacitance** of a reverse-biased junction:

```
Junction capacitance:
            ε × A
CJ = ────────────────
       W(VR)

Where W(VR) = √(2ε(V₀+VR)/qN)

Simplified model:
         CJ0
CJ = ────────────
     (1 + VR/V₀)^m

Where:
  CJ0 = zero-bias junction capacitance
  VR  = reverse voltage (positive number)
  V₀  = built-in potential
  m   = grading coefficient (0.33 to 0.5)

As VR increases → W increases → CJ decreases
```

**Applications:**
- Voltage-controlled oscillators (VCO) in PLLs
- FM radio tuning (varactor-tuned LC oscillator)
- Phase locked loops (PLLs)
- Parametric amplifiers

### 7. PIN Diode (P-Intrinsic-N)

Structure: P-type | Intrinsic (undoped) | N-type

- Intrinsic region: very lightly doped, wide depletion region
- At high frequencies: acts as variable resistor controlled by DC bias
  - Forward bias: low resistance (stored charge in I-region)
  - Reverse bias: high resistance
- Very fast switching at microwave frequencies

**Applications:**
- RF switches (phone antenna switching between LTE bands)
- Attenuators, phase shifters
- High voltage rectifiers (wide depletion → high breakdown)
- Radiation detectors (nuclear, X-ray)

### 8. Tunnel Diode

- Very heavily doped PN junction (both sides > 10¹⁹ cm⁻³)
- Quantum mechanical tunneling of electrons through thin depletion region
- Shows **negative differential resistance (NDR)** — unique property!

```
I-V curve:
      I
      │   ← peak current (VP, IP)
  IP  │───╮
      │    ╰─────── ← valley current (VV, IV)
  IV  │              ╰───────────────────── normal diode region
      │
      └─────────────────── V
          VP VV

In NDR region (VP < V < VV): as V increases, I DECREASES
```

**Applications:**
- Ultra-high frequency oscillators (microwave)
- Very fast switches
- Analog-to-digital converters

### 9. Avalanche Photodiode (APD)

- Reverse biased near breakdown voltage
- Photogenerated carriers undergo avalanche multiplication
- Internal gain: M = 10-1000×
- More sensitive than regular photodiode
- Used in: fiber optic receivers, LIDAR, medical imaging

### 10. TVS Diode (Transient Voltage Suppressor)

- Like Zener but designed for high peak power (kW for microseconds)
- Protects circuits from voltage spikes (ESD, lightning, inductive kickback)
- Very fast response: < 1ps for automotive
- Unidirectional or bidirectional (for AC signals)
- Examples: 1.5KE series, SMAJ series, P6KE series
- Used on: every microcontroller I/O pin protection, USB data lines, motor driver outputs

---

## 8. Rectifier Circuits

Rectifiers convert AC to DC.

### Half-Wave Rectifier

```
      D1
  ───►|───┬─── VOUT
AC in      │
  ────────┴─── GND

Output: only positive half-cycles pass
```

**Analysis:**
```
Vdc (average output) = Vm/π = 0.318 × Vm
Vrms output = Vm/2 = 0.5 × Vm

Where Vm = peak AC voltage
Ripple factor = 1.21 (very high ripple — poor DC)
Ripple frequency = fin (same as input)
```

**Example:** 9V transformer (RMS)
- Vm = 9 × √2 = 12.73V (peak)
- Diode drop: 0.7V
- Vm(output) = 12.73 - 0.7 = 12.03V
- Vdc = 12.03/π = 3.83V (very low — only used for simple applications)

### Full-Wave Rectifier (Center-Tap)

```
  ┌───►|─D1─┬─── VOUT
  │          │
Transformer  │
  │    D2   │
  └───►|────┘
  (center tap) ─── GND
```

**Analysis:**
```
Vdc = 2Vm/π = 0.636 × Vm
Vrms output = Vm/√2 = 0.707 × Vm
Ripple factor = 0.482 (better than half-wave)
Ripple frequency = 2 × fin
```

### Full-Wave Bridge Rectifier (Most Common!)

```
         D1      D3
     ┌───►|──┬───►|──┐
     │        │        │
AC ──┤       VOUT      ├── AC
     │        │        │
     └──|◄───┴──|◄───┘
         D4      D2

(D1,D2 conduct on positive half; D3,D4 on negative half)
```

**Analysis:**
```
Vdc = 2Vm/π - 2VD = 0.636Vm - 1.4V (two diode drops)
Vrms output = Vm/√2 - 2VD
Ripple factor = 0.482
Ripple frequency = 2 × fin

No center-tap transformer needed → simpler, cheaper, full secondary voltage used
```

### Filter Capacitor

Add capacitor in parallel with load to smooth output:

```
     ─────►|──────┬───── VOUT
                   │
                  [C]  (filter capacitor)
                   │
     ──────────────┴───── GND
```

**Capacitor selection:**
```
         IL
C = ──────────
      f × ΔV

Where:
  IL = load current (A)
  f  = ripple frequency (Hz) — 100Hz for bridge, 50Hz for half-wave (India)
  ΔV = acceptable ripple voltage (V)

Example: IL = 500mA, f = 100Hz, ΔV = 0.5V:
  C = 0.5 / (100 × 0.5) = 0.01F = 10,000μF → use 10,000μF electrolytic

After filtering: VOUT ≈ Vm - VD(drop) - ΔV/2
```

**Peak Inverse Voltage (PIV) across diodes:**
- Half-wave: PIV = Vm (peak AC voltage)
- Full-wave bridge: PIV = Vm (each diode)
- Full-wave center-tap: PIV = 2Vm

---

## 9. Clipper and Clamper Circuits

### Clipper (Limiter)

Removes part of waveform above/below a threshold.

**Series clipper (positive):**
```
Input → ─[D]─ → Output
            │
           GND

Removes negative half-cycles (diode blocks them)
```

**Shunt clipper:**
```
Input → ──┬── Output
           [D]
            │
           GND

Clips positive peaks (diode conducts when Vin > 0.7V → clamps output to 0.7V)
```

**Zener clipper (limits both peaks):**
```
Input → ──┬── Output
          [D1] in series with [ZD]
            │
           GND

Clips at +VZ+0.7V (forward ZD + forward D1)
Clips at −0.7V−VZ (reverse D1 breaks down at VZ, forward D2)
```

### Clamper (DC Restorer)

Shifts DC level of waveform without changing its shape.

**Positive clamper:**
```
Input ─[C]─┬── Output
            [D]
             │
            GND

Charges C to peak of input → shifts output UP
Output peak = 2×Vm (for sinusoid)
Output minimum = 0V
```

---

## 10. Diode Specifications

### Key Parameters to Understand

| Parameter | Symbol | Meaning |
|-----------|--------|---------|
| Average forward current | IF(AV) | Max DC current the diode can handle |
| Peak surge current | IFSM | Max peak current (very short pulses) |
| Forward voltage | VF | Voltage drop when conducting |
| Reverse voltage | VR, VRM, PIV | Max reverse voltage before breakdown |
| Reverse current | IR | Leakage current in reverse bias |
| Junction capacitance | Cj | Parasitic capacitance (limits high-freq) |
| Reverse recovery time | trr | Time to switch from forward to reverse |
| Storage temperature | Tstg | Storage temperature range |
| Junction temperature | TJ | Maximum junction temp (typically 150°C) |
| Thermal resistance | θJA | °C/W from junction to ambient |

### Reverse Recovery Time (trr)

Critical for switching applications:

```
When forward-biased diode is suddenly reverse-biased:

Current:
     │
  IF ├──────────────────╮
     │                   ╰──────────────────
   0 ├────────────────────────────────────── time
     │                         ↑
 -IR ├──────────────────────────╯
                    trr ←→

trr = time for minority carriers to be swept out/recombine
```

- **General purpose (1N4007):** trr = 1-10 μs — slow! Not for switching supplies
- **Fast recovery:** trr = 50-500 ns (MUR460, UF4007)
- **Ultra-fast:** trr = 20-75 ns (UF5401)
- **Schottky:** trr ≈ 0 (no minority storage — only majority carriers!)

### Derating

Always **derate** diode current in practice:
- Typically use 50-70% of rated current
- IF(rated) = 1A → use diode for max 0.5-0.7A
- Provides safety margin and improves reliability

---

## 11. Summary

```
PN JUNCTION FUNDAMENTALS
═════════════════════════

Depletion region:
  Built-in potential: V₀ = (kT/q)ln(NAND/ni²) ≈ 0.7V for Si

Forward bias (V > 0):
  I = IS × (e^(V/nVT) - 1) ≈ IS × e^(V/nVT) for V >> VT
  Turns on at ~0.7V (Si), ~0.3V (Ge), ~1.8-3.5V (LED)

Reverse bias (V < 0):
  I ≈ -IS (tiny leakage, ~nA for Si)
  Breakdown at VBR (Zener <5V, Avalanche >7V)

Dynamic resistance: rd = nVT/ID = 26mV/ID (n=1)

Diode types and uses:
  Rectifier:  power supply rectification
  Zener:      voltage regulation, reference
  Schottky:   fast switching, low Vf
  LED:        light emission (direct gap)
  Photodiode: light detection
  Varactor:   voltage-controlled capacitance
  TVS:        overvoltage protection
```

### Diode Quick Reference

| Diode Type | VF | Speed | Key Feature | Use |
|-----------|-----|-------|-------------|-----|
| Si rectifier | 0.7V | Slow | Cheap, high V | Power supplies |
| Ge diode | 0.3V | Medium | Low VF | Detection |
| Schottky | 0.2-0.4V | Very fast | No recovery | SMPS, RF |
| LED | 1.8-3.5V | N/A | Emits light | Indicators, lighting |
| Zener | VZ (spec) | Medium | Regulated V | Voltage regulation |
| TVS | VBR | Very fast | High power | Protection |
| Varactor | N/A | N/A | C vs V | RF tuning, VCO |
| PIN | Variable | Very fast | RF switch | Antenna switching |

---

**← Previous:** [Chapter 03: Semiconductors Deep Dive](./03-semiconductors-deep-dive.md)
**→ Next:** [Chapter 05: Basic Electronic Components](./05-basic-electronic-components.md)
