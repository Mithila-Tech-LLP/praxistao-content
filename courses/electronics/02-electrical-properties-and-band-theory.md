# Chapter 02: Electrical Properties and Band Theory

> **"Band theory is the quantum mechanical explanation for why some materials conduct electricity and others don't — and it tells us exactly how to engineer conductivity in semiconductors."**

---

## Table of Contents
1. [Electric Charge and Current](#1-electric-charge-and-current)
2. [Voltage, Resistance, and Ohm's Law](#2-voltage-resistance-and-ohms-law)
3. [Power and Energy](#3-power-and-energy)
4. [Energy Band Theory](#4-energy-band-theory)
5. [Classification by Band Structure](#5-classification-by-band-structure)
6. [Fermi-Dirac Distribution](#6-fermi-dirac-distribution)
7. [Carrier Concentration](#7-carrier-concentration)
8. [Conductivity and Mobility](#8-conductivity-and-mobility)
9. [Kirchhoff's Laws](#9-kirchhoffs-laws)
10. [Series and Parallel Circuits](#10-series-and-parallel-circuits)
11. [Wheatstone Bridge](#11-wheatstone-bridge)
12. [AC Fundamentals](#12-ac-fundamentals)
13. [Summary of Key Formulas](#13-summary-of-key-formulas)

---

## 1. Electric Charge and Current

### Electric Charge

**Electric charge** is a fundamental property of matter. It comes in two types:
- **Positive charge (+)**: protons, holes (in semiconductors)
- **Negative charge (−)**: electrons

The elementary charge:
```
q = 1.602 × 10⁻¹⁹ Coulombs (C)
```

**Coulomb's Law** (force between two charges):
```
       q₁ × q₂
F = k × ────────
          r²

Where:
  F  = force (Newtons)
  k  = 8.99×10⁹ N·m²/C² (Coulomb's constant)
  q₁, q₂ = charges (Coulombs)
  r  = distance between charges (meters)

Same sign charges: REPEL (positive force)
Opposite sign charges: ATTRACT (negative force)
```

### Electric Current

**Current (I)** is the rate of flow of charge:

```
     ΔQ
I = ────
     Δt

Or in differential form:

     dQ
I = ────
     dt

Units: Amperes (A) = Coulombs/second (C/s)
```

**Convention:**
- **Conventional current** flows from **+ to −** (positive charges moving)
- **Electron flow** (actual) is from **− to +** (electrons moving toward positive terminal)
- In all circuit analysis, we use conventional current direction

**Example:** A current of 2A means 2 Coulombs of charge passes a point every second.
- Number of electrons per second = 2 / (1.602×10⁻¹⁹) = **1.25×10¹⁹ electrons/second**

### Types of Current
- **Direct Current (DC)**: flows in one direction, constant magnitude (batteries, DC power supplies)
- **Alternating Current (AC)**: periodically reverses direction (household power, 50Hz/60Hz)
- **Pulsating DC**: DC that varies in magnitude but not direction (rectified AC)

---

## 2. Voltage, Resistance, and Ohm's Law

### Voltage (Electric Potential Difference)

**Voltage (V)** is the energy per unit charge — it's what "pushes" current through a circuit.

```
     W
V = ───
     Q

Units: Volts (V) = Joules/Coulomb (J/C)
```

- Think of voltage as "electrical pressure"
- Current flows from HIGH voltage to LOW voltage
- Voltage is always measured **between two points** (it's a difference)

**Common voltage references:**
- AA battery: 1.5V
- USB: 5V
- Household (India): 230V AC
- Household (USA): 120V AC
- Car battery: 12V
- Arduino logic: 5V
- ESP32 logic: 3.3V
- Modern CPU core: 0.7-1.1V

### Resistance

**Resistance (R)** is the opposition to current flow.

```
Units: Ohms (Ω)
1 Ω = 1 V/A
```

**Resistivity formula** (material property):
```
         ρ × L
R = ──────────────
          A

Where:
  ρ = resistivity (Ω·m) — material property
  L = length of conductor (m)
  A = cross-sectional area (m²)
```

**Resistivity values at 20°C:**

| Material | Resistivity (Ω·m) | Type |
|----------|------------------|------|
| Silver (Ag) | 1.59×10⁻⁸ | Conductor |
| Copper (Cu) | 1.72×10⁻⁸ | Conductor |
| Gold (Au) | 2.44×10⁻⁸ | Conductor |
| Aluminum (Al) | 2.82×10⁻⁸ | Conductor |
| Tungsten (W) | 5.60×10⁻⁸ | Conductor |
| Iron (Fe) | 1.00×10⁻⁷ | Conductor |
| Silicon (intrinsic) | 6.4×10² | Semiconductor |
| Germanium (intrinsic) | 4.6×10⁻¹ | Semiconductor |
| Silicon (doped) | 10⁻⁵ to 10⁰ | Semiconductor |
| Glass | 10¹⁰ to 10¹⁴ | Insulator |
| Rubber | 10¹³ to 10¹⁵ | Insulator |
| Teflon (PTFE) | 10²³ | Best insulator |
| Diamond | 10¹⁰ to 10¹⁴ | Insulator |

**Conductivity (σ)** = inverse of resistivity:
```
      1
σ = ─────
      ρ

Units: Siemens/meter (S/m)
```

### Ohm's Law

One of the most fundamental laws in electronics:

```
V = I × R

Or equivalently:
  I = V/R
  R = V/I
```

**Georg Simon Ohm** discovered this in 1827 — the voltage across a conductor is **directly proportional** to the current through it, provided temperature is constant.

**Ohm's Law Triangle:**
```
        V
       ─┬─
      I  ×  R

Cover what you want:
• V = I × R
• I = V / R
• R = V / I
```

**Example problems:**

*Q: What is the current through a 220Ω resistor with 5V across it?*
```
I = V/R = 5/220 = 0.0227 A = 22.7 mA
```

*Q: LED has 2V forward drop, powered by 5V supply. What resistor limits current to 20mA?*
```
Voltage across resistor = 5 - 2 = 3V
R = V/I = 3/0.020 = 150 Ω (use 150Ω or 220Ω standard value)
```

---

## 3. Power and Energy

### Electrical Power

**Power (P)** is the rate of energy conversion:

```
P = V × I

Substituting Ohm's Law:
  P = I² × R    (when you know I and R)
  P = V² / R    (when you know V and R)

Units: Watts (W) = Joules/second (J/s)
```

### Energy

```
E = P × t

Units: Joules (J) = Watt × seconds
Or: Watt-hours (Wh) = P × t (hours)
Or: kilowatt-hours (kWh) = unit used by electricity companies
```

**Examples:**
- LED: 0.02A × 3.3V = **66mW**
- Arduino Uno (full load): 5V × 0.05A = **250mW**
- ESP32 (WiFi active): 5V × 0.25A = **1.25W**
- Raspberry Pi 4: 5V × 0.6-3A = **3-15W**
- Desktop PC: **150-600W**
- Electric car motor: **50,000-500,000W (50-500kW)**

---

## 4. Energy Band Theory

This is the **heart of understanding semiconductors**. Band theory explains why:
- Copper conducts electricity
- Glass doesn't
- Silicon is "in between" and can be controlled

### From Atomic Orbitals to Energy Bands

When atoms are **isolated**, electrons occupy specific **discrete energy levels** (as seen in Chapter 01).

When billions of atoms come together to form a **crystal**, these discrete levels **split and spread** into continuous **energy bands** due to quantum mechanical interactions between neighboring atoms.

```
Isolated atoms → Crystal formation
─────────────────────────────────

Isolated Si atom:        Silicon crystal (many atoms):
   ───── 3p              ╔══════════════╗ ← Conduction Band (CB)
   ───── 3s                   Gap
   ───── 2p              ╚══════════════╝ ← Valence Band (VB)
   ───── 2s
   ───── 1s              (lower bands fully filled, not shown)
```

### Key Bands

**1. Valence Band (VB)**
- Highest energy band that is completely (or mostly) filled with electrons at 0K
- Electrons here are involved in covalent bonding
- These electrons are NOT free to conduct current (they're "used up" in bonds)

**2. Conduction Band (CB)**
- The next higher energy band, mostly empty at room temperature
- Electrons that reach this band are **FREE** to move and conduct current
- Separated from valence band by the "band gap"

**3. Forbidden Energy Gap / Band Gap (Eg)**
- Energy range where NO electron states exist
- Electrons cannot exist with energies in this range
- To move from valence band to conduction band, an electron needs energy ≥ Eg
- This energy comes from:
  - Thermal energy (heat)
  - Photons of light
  - Doping (adding impurity atoms)

```
Energy diagram:

CB (Conduction Band):  ════════════════  ← Free electrons can exist here
                            Eg (gap)
VB (Valence Band):     ════════════════  ← Bound electrons in bonds

At 0K: VB full, CB empty → no conduction
At 300K: some electrons thermally excited from VB to CB → some conduction
```

### Band Gap Values

| Material | Band Gap (eV) | Type | Color of light (if direct gap) |
|----------|--------------|------|-------------------------------|
| Ge | 0.67 | Indirect | — |
| Si | 1.12 | Indirect | — |
| GaAs | 1.42 | Direct | Near-IR (870nm) |
| InP | 1.34 | Direct | IR |
| GaP | 2.26 | Indirect | Green/Red |
| GaN | 3.40 | Direct | UV (365nm) |
| SiC | 3.26 | Indirect | — |
| ZnO | 3.37 | Direct | UV |
| Diamond | 5.47 | Indirect | — |
| SiO₂ | ~9 | Insulator | — |
| Al₂O₃ | ~8.8 | Insulator | — |
| Si₃N₄ | ~5 | Insulator | — |

**1 eV = 1.602×10⁻¹⁹ Joules** (the energy gained by 1 electron moving through 1 volt potential)

---

## 5. Classification by Band Structure

### Conductors

```
Energy:
  ┌──────────────────┐
  │   Conduction Band│ ← partially filled OR overlaps with VB
  │  ●●●●    (some   │
  │  empty spaces)   │
  │  or              │
  ├──────────────────┤ ← bands overlap (no gap!)
  │   Valence Band   │
  └──────────────────┘
```

- Valence band overlaps with conduction band **OR** conduction band is partially filled
- Electrons are free to move to slightly higher energy states
- Even a tiny electric field causes current flow
- **No band gap** (or negative gap — bands overlap)
- Resistivity: ~10⁻⁸ Ω·m
- Examples: Cu, Al, Ag, Au, Fe, all metals
- **Temperature effect:** resistance INCREASES with temperature
  (more thermal vibrations → more collisions → more resistance)
  - Resistance temperature coefficient (Cu): α = +0.00393/°C

### Insulators

```
Energy:
  ┌──────────────────┐
  │  Conduction Band │ ← completely empty
  │                  │
        Large Gap
       Eg > 5 eV
  │                  │
  │  Valence Band    │ ← completely full
  └──────────────────┘
```

- Huge band gap (> 5 eV)
- Thermal energy at room temperature (kT ≈ 26 meV) is FAR too small to excite electrons across
- Essentially no free carriers
- Resistivity: > 10⁸ Ω·m
- Examples: SiO₂, Al₂O₃, Si₃N₄, rubber, glass, diamond, air
- Used as: gate oxide in MOSFETs, wire insulation, PCB substrate

### Semiconductors

```
Energy:
  ┌──────────────────┐
  │  Conduction Band │ ← nearly empty (a few electrons at room temp)
  │  ●               │
        Small Gap
       Eg = 0.1-3 eV
  │                ○  │ ← holes (missing electrons)
  │  Valence Band    │ ← nearly full
  └──────────────────┘
```

- Small band gap (0.1 to ~3 eV)
- At room temperature, **some** electrons have enough thermal energy to cross the gap
- These electrons in CB + holes in VB = charge carriers
- Resistivity: 10⁻⁴ to 10⁴ Ω·m (tunable!)
- **Temperature effect:** resistance DECREASES with temperature
  (more thermal energy → more carriers → more conductance)
- Examples: Si, Ge, GaAs, InP, GaN, SiC

**The magic of semiconductors:** we can CONTROL how many carriers exist by:
1. **Temperature** (more heat → more carriers)
2. **Light** (photons excite electrons across gap)
3. **Doping** (adding impurities creates extra carriers)
4. **Electric field** (in MOSFET: gate voltage creates/destroys carrier channel)

---

## 6. Fermi-Dirac Distribution

The **Fermi-Dirac distribution** tells us the **probability** that a given energy state is occupied by an electron at thermal equilibrium.

```
           1
f(E) = ─────────────────
        1 + e^((E-EF)/kT)

Where:
  f(E) = probability of state at energy E being occupied
  EF   = Fermi energy (Fermi level)
  k    = Boltzmann's constant = 8.617 × 10⁻⁵ eV/K
  T    = Temperature in Kelvin
  kT   = 26 meV at T = 300K (room temperature)
```

**Important values:**

| Condition | Result |
|-----------|--------|
| E = EF | f(EF) = 1/(1+1) = **0.5 always** |
| E << EF (much below) | f(E) → **1** (almost certainly occupied) |
| E >> EF (much above) | f(E) → **0** (almost certainly empty) |
| At T = 0K | Perfect step: all below EF filled, all above empty |

**At 300K, kT = 0.026 eV:**
```
At E = EF + 0.1 eV (100meV above Fermi level):
f = 1/(1+e^(100/26)) = 1/(1+47.8) = 0.02 = 2%

At E = EF + 3kT (78meV above):
f = 1/(1+e³) = 1/(1+20.1) = 4.5%
```

### Fermi Level Position

The **Fermi level** position within the energy band structure determines the material type:

```
Conductor:      EF is inside partially filled band
                ═══════EF══════ ← electrons available at EF
                  many states with electrons near EF

Intrinsic SC:   EF is near midgap
                ═══════════════ ← CB (empty at 0K)
                    --- EF ---  ← Fermi level at midgap
                ═══════════════ ← VB (full at 0K)

N-type SC:      EF moves UP (closer to CB)
                ═══════════════ ← CB
                        EF ---  ← near conduction band
                ═══════════════ ← VB

P-type SC:      EF moves DOWN (closer to VB)
                ═══════════════ ← CB
                --- EF          ← near valence band
                ═══════════════ ← VB

Insulator:      EF in large forbidden gap, far from both bands
```

---

## 7. Carrier Concentration

### Intrinsic Carrier Concentration

For a **pure (intrinsic) semiconductor**, the number of electrons in the conduction band equals the number of holes in the valence band:

```
n = p = ni

where ni = intrinsic carrier concentration
```

The formula for ni:

```
ni = √(Nc × Nv) × exp(−Eg / 2kT)

Where:
  Nc = effective density of states in conduction band
       Nc(Si) = 2.8 × 10¹⁹ cm⁻³ at 300K
  Nv = effective density of states in valence band
       Nv(Si) = 1.04 × 10¹⁹ cm⁻³ at 300K
  Eg = band gap energy
  kT = 26 meV at 300K
```

**ni values at 300K:**

| Material | ni (cm⁻³) | Comment |
|----------|-----------|---------|
| Si | 1.5 × 10¹⁰ | Very low — pure Si is nearly insulator |
| Ge | 2.4 × 10¹³ | ~1000× more than Si (smaller gap) |
| GaAs | 1.8 × 10⁶ | Very low (large gap) |
| GaN | ~10⁻¹⁰ | Essentially none at room temp |

**Temperature dependence:** ni doubles for every ~11°C increase in Si temperature!

### Mass Action Law

In any semiconductor (doped or not), the **product** of electron and hole concentrations is always constant at a given temperature:

```
n × p = ni²

At 300K for Si: n × p = (1.5×10¹⁰)² = 2.25×10²⁰ cm⁻⁶
```

This fundamental relationship is used everywhere in semiconductor physics.

---

## 8. Conductivity and Mobility

### Carrier Drift in Electric Field

When an electric field **E** is applied to a semiconductor:
- Electrons drift in the **opposite** direction to E (toward + terminal)
- Holes drift in the **same** direction as E (toward − terminal)

**Drift velocity:**
```
vd = μ × E

Where:
  vd = drift velocity (cm/s)
  μ  = carrier mobility (cm²/V·s)
  E  = electric field (V/cm)
```

### Mobility Values at 300K

| Material | μn (electrons) | μp (holes) |
|----------|---------------|------------|
| Si | 1400 cm²/V·s | 450 cm²/V·s |
| Ge | 3900 cm²/V·s | 1900 cm²/V·s |
| GaAs | 8500 cm²/V·s | 400 cm²/V·s |

**Note:** electrons always have higher mobility than holes! That's why N-type devices are faster.

**GaAs electrons are ~6× more mobile than Si** → why GaAs is used in high-speed RF circuits.

### Conductivity Formula

```
σ = q × (n × μn + p × μp)

Where:
  σ  = conductivity (S/cm)
  q  = electron charge = 1.602×10⁻¹⁹ C
  n  = electron concentration (cm⁻³)
  p  = hole concentration (cm⁻³)
  μn = electron mobility (cm²/V·s)
  μp = hole mobility (cm²/V·s)
```

**Ohm's Law in field form:**
```
J = σ × E

Where:
  J = current density (A/cm²)
  σ = conductivity (S/cm)
  E = electric field (V/cm)
```

**Relating to circuit quantities:**
```
J = I/A    and    E = V/L

So: I/A = σ × V/L
    I = σ × A/L × V
    I = V / (L/σA) = V/R

This confirms Ohm's Law: V = IR where R = ρL/A
```

---

## 9. Kirchhoff's Laws

These two laws are the **foundation of all circuit analysis**. Every circuit simulator (SPICE, LTspice) is built on them.

### Kirchhoff's Current Law (KCL)

> **"The algebraic sum of all currents entering a node (junction) equals zero."**

Or equivalently: **Sum of currents entering a node = Sum of currents leaving a node**

```
Mathematical form: Σ I = 0

Example:
         I₁ = 3A →  ●  ← I₂ = 1A
                     ↓
                    I₃ = ?

KCL: I₁ + I₂ = I₃
     3 + 1 = I₃
     I₃ = 4A (leaving the node)
```

Based on **conservation of charge** — charge cannot be created or destroyed.

### Kirchhoff's Voltage Law (KVL)

> **"The algebraic sum of all voltages around any closed loop in a circuit equals zero."**

Or: **Sum of voltage rises = Sum of voltage drops**

```
Mathematical form: Σ V = 0 (around any closed loop)

Example: Simple series circuit
    +9V battery, R₁=100Ω, R₂=200Ω

KVL: +9 - V_R1 - V_R2 = 0
     9 - 100I - 200I = 0
     9 = 300I
     I = 30mA
     V_R1 = 100 × 0.03 = 3V
     V_R2 = 200 × 0.03 = 6V
     Check: 3 + 6 = 9V ✓
```

Based on **conservation of energy** — energy gained in a loop equals energy lost.

---

## 10. Series and Parallel Circuits

### Series Circuits

Elements connected **end-to-end**, single path for current:

```
+──[ R₁ ]──[ R₂ ]──[ R₃ ]──-
```

**Rules:**
- **Same current** flows through all: I₁ = I₂ = I₃ = I
- **Voltages add**: Vtotal = V₁ + V₂ + V₃
- **Resistances add**: Rtotal = R₁ + R₂ + R₃

**Voltage divider formula:**
```
         R₂
Vout = ────── × Vin
       R₁ + R₂

Example: 10kΩ and 10kΩ in series across 5V:
Vout (across R₂) = 10/(10+10) × 5 = 2.5V

Example: sensor output divider — if sensor R changes from 1kΩ to 10kΩ:
At 1kΩ: Vout = 1/(1+10) × 5 = 0.45V
At 10kΩ: Vout = 10/(10+10) × 5 = 2.5V
→ ADC reads changing voltage to detect sensor change
```

### Parallel Circuits

Elements connected **side-by-side**, multiple paths for current:

```
    +──┬──[ R₁ ]──┬──-
       ├──[ R₂ ]──┤
       └──[ R₃ ]──┘
```

**Rules:**
- **Same voltage** across all: V₁ = V₂ = V₃ = V
- **Currents add**: Itotal = I₁ + I₂ + I₃
- **Conductances add**: Gtotal = G₁ + G₂ + G₃
- **Resistances** combine as:
```
1/Rtotal = 1/R₁ + 1/R₂ + 1/R₃

For just 2 resistors (product-over-sum rule):
          R₁ × R₂
Rtotal = ──────────
          R₁ + R₂
```

**Current divider formula:**
```
         R₂
I₁ = ────────── × Itotal   (current through R₁)
     R₁ + R₂

Note: larger resistance gets LESS current!
```

### Mixed Series-Parallel Analysis

Steps:
1. Identify series and parallel groups
2. Simplify parallel groups to equivalent resistors
3. Now solve the simplified series circuit
4. Work backwards to find individual currents/voltages

---

## 11. Wheatstone Bridge

A classic circuit for **precise resistance measurement** (and the basis of many sensor circuits):

```
         +Vcc
          |
     ┌────┴────┐
    R₁         R₂
     ├────┬────┤
    R₃    │   R₄
     └────┴────┘
          |
          GND

     Vout = voltage between middle nodes
```

**Balance condition (Vout = 0):**
```
R₁/R₂ = R₃/R₄
Or equivalently: R₁ × R₄ = R₂ × R₃
```

**When balanced:** Replace one resistor (e.g., R₃) with a sensor (thermistor, strain gauge, RTD):
- If sensor changes, bridge becomes unbalanced
- Vout ≠ 0, and Vout ∝ change in sensor resistance
- Used in: load cells (weight measurement), temperature sensors, pressure sensors

---

## 12. AC Fundamentals

### Sinusoidal Waveform

AC voltage varies sinusoidally with time:
```
v(t) = Vm × sin(2πft + φ) = Vm × sin(ωt + φ)

Where:
  Vm = peak (amplitude) voltage
  f  = frequency (Hz) = cycles per second
  ω  = angular frequency = 2πf (rad/s)
  T  = period = 1/f (seconds)
  φ  = phase angle (radians or degrees)
```

**India household AC:** 230V RMS, 50Hz
- So Vm = 230 × √2 ≈ 325V peak
- T = 1/50 = 20ms
- ω = 2π × 50 = 314.16 rad/s

### RMS (Root Mean Square) Value

The **effective value** of AC — produces same heating effect as equivalent DC:

```
           Vm
Vrms = ─────── ≈ 0.707 × Vm    (for pure sine wave)
          √2
```

**Power in AC circuit:**
```
P = Vrms × Irms × cos(φ)    (real power, Watts)
Q = Vrms × Irms × sin(φ)    (reactive power, VAR)
S = Vrms × Irms              (apparent power, VA)
PF = cos(φ)                  (power factor, 0 to 1)
```

### Reactance

**Capacitive reactance** (capacitor's opposition to AC):
```
         1
Xc = ─────────
     2π × f × C

Units: Ohms
Note: Xc DECREASES as frequency INCREASES (capacitor passes high-freq, blocks DC)
```

**Inductive reactance** (inductor's opposition to AC):
```
XL = 2π × f × L

Units: Ohms
Note: XL INCREASES as frequency INCREASES (inductor passes DC, blocks high-freq)
```

**Impedance (Z):** complex resistance for AC circuits
```
Z = R + jX    (j = imaginary unit = √-1)

For RC circuit:
Z = R - j/ωC = √(R² + Xc²) at angle -arctan(Xc/R)

For RL circuit:
Z = R + jωL = √(R² + XL²) at angle +arctan(XL/R)
```

### Resonance

When XL = Xc, the circuit is at **resonance**:
```
         1
f₀ = ──────────
     2π√(L×C)

At resonance:
- Series RLC: impedance is minimum (= R), current maximum
- Parallel RLC: impedance is maximum, current minimum
- Phase angle = 0 (pure resistive)
```

---

## 13. Summary of Key Formulas

| Formula | Meaning | Units |
|---------|---------|-------|
| I = Q/t | Current = charge/time | A |
| V = W/Q | Voltage = energy/charge | V |
| V = IR | Ohm's Law | V, A, Ω |
| R = ρL/A | Resistance from geometry | Ω |
| σ = 1/ρ | Conductivity | S/m |
| P = VI = I²R = V²/R | Power | W |
| E = Pt | Energy | J |
| ni = √(NcNv)·e^(-Eg/2kT) | Intrinsic carrier concentration | cm⁻³ |
| n·p = ni² | Mass action law | cm⁻⁶ |
| σ = q(nμn + pμp) | Conductivity (carriers) | S/cm |
| J = σE | Current density | A/cm² |
| f(E) = 1/(1+e^((E-EF)/kT)) | Fermi-Dirac | probability |
| Xc = 1/(2πfC) | Capacitive reactance | Ω |
| XL = 2πfL | Inductive reactance | Ω |
| f₀ = 1/(2π√LC) | Resonant frequency | Hz |
| Vrms = Vm/√2 | RMS voltage (sinusoid) | V |

### Constants

| Constant | Symbol | Value |
|----------|--------|-------|
| Electron charge | q | 1.602×10⁻¹⁹ C |
| Boltzmann constant | k | 8.617×10⁻⁵ eV/K |
| Thermal voltage (300K) | kT/q | 26 mV |
| Permittivity of Si | εSi | 11.7 × ε₀ |
| Permittivity of free space | ε₀ | 8.854×10⁻¹² F/m |

---

**← Previous:** [Chapter 01: Atoms and Matter](./01-atoms-and-matter.md)
**→ Next:** [Chapter 03: Semiconductors Deep Dive](./03-semiconductors-deep-dive.md)
