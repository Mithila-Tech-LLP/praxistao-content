# Chapter 48: Capacitors

> **"A capacitor is a device that stores electrical energy in an electric field. From the tiny capacitors in your phone to the huge banks in power plants — they are fundamental to all electronics."**

---

## Table of Contents

- [48.1 What is a Capacitor?](#481-what-is-a-capacitor)
- [48.2 Capacitance](#482-capacitance)
- [48.3 Parallel Plate Capacitor](#483-parallel-plate-capacitor)
- [48.4 Dielectrics](#484-dielectrics)
- [48.5 Energy Stored in a Capacitor](#485-energy-stored-in-a-capacitor)
- [48.6 Capacitors in Series](#486-capacitors-in-series)
- [48.7 Capacitors in Parallel](#487-capacitors-in-parallel)
- [48.8 Charging and Discharging](#488-charging-and-discharging)
- [48.9 Applications of Capacitors](#489-applications-of-capacitors)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 48.1 What is a Capacitor?

A **capacitor** is any device that stores electric charge and energy — consisting of two conductors (plates) separated by an insulator (dielectric).

```
CAPACITOR SYMBOL:

  ---|  |---

  (two lines representing the two plates)

PHYSICAL STRUCTURE (parallel plate):

  + + + + + + +   (positive plate, connected to + terminal)
  -  -  -  -  -   (insulator/dielectric)
  - - - - - - -   (negative plate, connected to - terminal)
  
  When connected to a battery:
  - Electrons flow off the + plate
  - Electrons flow onto the - plate
  - Opposite charges create an electric field between plates
```

How is charge held? The + plate attracts the negative charge on the other plate, and vice versa. The charges are "stuck" in place by the attraction — they can't cross the insulator.

---

## 48.2 Capacitance

**Capacitance** (C) measures how much charge a capacitor can store per volt of potential difference:

```
C = Q / V

where:
  C = capacitance (Farads, F)
  Q = charge stored on one plate (C)
  V = voltage across the capacitor (V)
```

A larger C means more charge stored for the same voltage.

### SI Unit: Farad

The Farad is a huge unit — real capacitors are typically:
- Microfarads: 1 μF = 10⁻⁶ F (common in electronics)
- Nanofarads: 1 nF = 10⁻⁹ F (radio/RF circuits)
- Picofarads: 1 pF = 10⁻¹² F (high-frequency circuits)
- Supercapacitors: up to thousands of Farads (energy storage)

---

## 48.3 Parallel Plate Capacitor

The simplest capacitor is two parallel conducting plates:

```
Area A
+------------------+
|+ + + + + + + + +|  d (separation)
|                  |
|- - - - - - - - -|
+------------------+
Area A
```

For a parallel plate capacitor:

```
C = ε₀ × A / d

where:
  ε₀ = 8.85 × 10⁻¹² F/m  (permittivity of free space)
  A = area of each plate (m²)
  d = separation between plates (m)
```

To increase capacitance:
- Increase plate area A (more charge per unit area)
- Decrease separation d (stronger electric field, more charge attracted)
- Add a dielectric (see next section)

### Worked Example 48.1

A parallel plate capacitor has plate area 0.01 m² and separation 1 mm. Find C.

**Solution:**

C = ε₀A/d = 8.85×10⁻¹² × 0.01 / 0.001 = **88.5 pF**

How much charge is stored at 12 V?

Q = CV = 88.5×10⁻¹² × 12 = **1.06 × 10⁻⁹ C = 1.06 nC**

---

## 48.4 Dielectrics

A **dielectric** is an insulating material placed between the capacitor plates.

When a dielectric is inserted:

```
WITH VACUUM:                     WITH DIELECTRIC:
  + + + + + + +                    + + + + + + +
     E₀ →                            [dielectric]
  - - - - - - -                       E₀/κ →
                                    - - - - - - -
  
  C₀ = ε₀A/d                       C = κε₀A/d
```

The **dielectric constant** κ (kappa) ≥ 1 tells how much the capacitance is increased.

```
C = κ × C₀ = κ × ε₀A/d = ε_r × ε₀A/d

where κ = ε_r = relative permittivity (dielectric constant)
```

### Why Does a Dielectric Increase Capacitance?

```
POLARIZATION OF DIELECTRIC:

  External E field →→→→→→→
  
  Molecules in dielectric align:
  
  ⊕-⊖  ⊕-⊖  ⊕-⊖   (dipoles align with field)
  
  The - sides of dipoles face the + plate
  The + sides face the - plate
  
  This creates a surface charge that PARTIALLY CANCELS
  the field → you need more charge on plates to maintain
  the same voltage → higher capacitance!
```

### Dielectric Constants

| Material | κ |
|---------|---|
| Vacuum  | 1 |
| Air     | 1.00059 |
| Paper   | 3.5 |
| Mica    | 5-7 |
| Water   | 80 |
| BaTiO₃  | ~1000-5000 |

---

## 48.5 Energy Stored in a Capacitor

A charged capacitor stores energy in its electric field:

```
Energy = (1/2) × C × V²
        = (1/2) × Q × V
        = Q² / (2C)
```

All three forms are equivalent; use whichever variables you know.

### Worked Example 48.2

A 100 μF capacitor is charged to 9 V.

(a) Find the charge stored.
(b) Find the energy stored.

**Solution:**

(a) Q = CV = 100×10⁻⁶ × 9 = **9×10⁻⁴ C = 0.9 mC**

(b) E = (1/2)CV² = (1/2) × 100×10⁻⁶ × 81 = **4.05×10⁻³ J = 4.05 mJ**

### Energy Density in the Field

```
Energy density: u = (1/2) × ε₀ × E²

Total energy = u × Volume = (1/2)ε₀E² × A × d
             = (1/2)ε₀(V/d)² × Ad
             = (1/2)(ε₀A/d)V²
             = (1/2)CV²  ✓
```

The energy is stored in the electric field between the plates — not in the plates themselves.

---

## 48.6 Capacitors in Series

When capacitors are connected end-to-end (in series):

```
---| C₁ |---| C₂ |---| C₃ |---

Each capacitor gets the SAME charge Q.
The voltages add up: V = V₁ + V₂ + V₃

Since Q = CV:
  V = Q/C₁ + Q/C₂ + Q/C₃
  V/Q = 1/C₁ + 1/C₂ + 1/C₃

So: 1/C_total = 1/C₁ + 1/C₂ + 1/C₃ + ...
```

Series capacitors: total capacitance **less than the smallest** individual capacitor.

### Worked Example 48.3

Three capacitors: C₁ = 6 μF, C₂ = 3 μF, C₃ = 2 μF in series. Find C_total.

**Solution:**

1/C = 1/6 + 1/3 + 1/2 = 1/6 + 2/6 + 3/6 = 6/6 = 1

C_total = **1 μF**

---

## 48.7 Capacitors in Parallel

When capacitors are connected with both plates connected (parallel):

```
      +--| C₁ |--+
      |           |
V ----|--| C₂ |--|---- V
      |           |
      +--| C₃ |--+

Each capacitor gets the SAME voltage V.
The charges add: Q = Q₁ + Q₂ + Q₃

Since Q = CV:
  C_total × V = C₁V + C₂V + C₃V

So: C_total = C₁ + C₂ + C₃ + ...
```

Parallel capacitors: total capacitance is the **sum** of all capacitors.

### Worked Example 48.4

Same three capacitors (6 μF, 3 μF, 2 μF) in parallel:

C_total = 6 + 3 + 2 = **11 μF**

### Memory Aid

```
Resistors:   series = sum,   parallel = reciprocal sum
Capacitors:  series = reciprocal sum,   parallel = sum

(Exactly backwards from resistors!)
```

---

## 48.8 Charging and Discharging

When a capacitor charges through a resistor:

```
CHARGING CIRCUIT:
  
  Battery V₀ ----|--- R ---|--- C ---|--- battery
  
  At t=0: C is empty, voltage across C = 0, all voltage across R
  As time goes on: C fills up, voltage across C rises, current falls
  At t=∞: C fully charged to V₀, no current flows
```

The voltage across the capacitor during charging:

```
V_C(t) = V₀ × (1 - e^(-t/τ))

During discharging:
V_C(t) = V₀ × e^(-t/τ)

where τ = RC = time constant (seconds)
```

The time constant τ = RC tells how fast the capacitor charges/discharges:

```
After t = τ:    63% charged
After t = 2τ:   86% charged
After t = 3τ:   95% charged
After t = 5τ:   99.3% charged (essentially fully charged)
```

```
CHARGING GRAPH:
  V_C
  V₀ |          ________
     |        /
  0.63V₀ |    *  <- at t=τ
     |  /
     | /
     |/
  0  +---------> t
         τ
```

### Worked Example 48.5

A 100 μF capacitor charges through a 10 kΩ resistor to 12 V.

(a) Find the time constant.
(b) Find voltage after 2 s.
(c) How long to reach 10 V?

**Solution:**

(a) τ = RC = 10,000 × 100×10⁻⁶ = **1 second**

(b) V = 12(1 - e^(-2/1)) = 12(1 - e^(-2)) = 12(1 - 0.135) = **10.4 V**

(c) 10 = 12(1 - e^(-t))
    e^(-t) = 2/12 = 0.167
    -t = ln(0.167) = -1.79
    t = **1.79 s**

---

## 48.9 Applications of Capacitors

### 1. Energy Storage (Camera Flash)

A camera flash capacitor (a few thousand μF) stores energy, then dumps it in ~1 ms through the flash bulb — a peak power of thousands of watts from a tiny battery.

```
Flash circuit:
  Battery slowly charges capacitor → flash trigger →
  capacitor dumps all energy in 1 ms → bright flash
```

### 2. Smoothing (Power Supplies)

Capacitors smooth out the pulsating output of a rectifier in power supplies:

```
Rectified DC (pulsating):          Smoothed DC:
  |                                  _____________
  |  ||| |||                        /             \
  | / \ / \                        /               \
  |/   X   \                      /                 \
  +-----------> t                 +-----------------> t
  
  Capacitor charges on peak,       Capacitor slowly discharges
  discharges between peaks         between peaks — smoother!
```

### 3. Timing Circuits

RC circuits (resistor + capacitor) create precise time delays. Used in:
- Oscillator circuits
- Timing pulses for microprocessors
- Delay timers

### 4. Memory (DRAM)

Dynamic RAM stores each bit as a charge on a tiny capacitor (about 25 fF = 25×10⁻¹⁵ F). Your computer's RAM is billions of these.

```
Memory cell:
  Transistor (switch) --- Capacitor
  
  Charged = 1
  Discharged = 0
  
  Leaks slowly → must refresh every few milliseconds → "dynamic"
```

### 5. Touchscreens

The glass of a smartphone has a grid of transparent capacitors. Your finger (a conductor) changes the capacitance of nearby cells, telling the phone where you touched.

### 6. Defibrillators

Medical defibrillators store ~300 J in a capacitor, then discharge it in ~10 ms through the patient's chest — enough to reset the heart rhythm.

```
Energy: E = (1/2)CV² = 300 J
At V = 3000 V: C = 2×300/3000² = 67 μF
```

---

## Summary

- **Capacitor**: two conductors separated by insulator; stores charge and energy
- **Capacitance**: C = Q/V; unit Farads; larger C = more charge per volt
- **Parallel plate**: C = ε₀A/d; larger area or smaller gap → larger C
- **Dielectric**: insulating material between plates; multiplies C by κ (dielectric constant)
- **Energy stored**: E = ½CV² = ½QV = Q²/2C; stored in the electric field
- **In series**: 1/C_total = 1/C₁ + 1/C₂ + ... (total less than smallest)
- **In parallel**: C_total = C₁ + C₂ + ... (total is sum)
- **Charging**: V_C = V₀(1 - e^(-t/τ)); discharging: V_C = V₀e^(-t/τ); τ = RC
- Applications: flash photography, power supply smoothing, timing, memory, touchscreens, defibrillators

---

## Key Equations

```
Capacitance:
  C = Q / V

Parallel plate capacitor:
  C = ε₀ × A / d  (vacuum)
  C = κ × ε₀A/d   (with dielectric κ)

Energy stored:
  E = (1/2) × C × V²
  E = (1/2) × Q × V
  E = Q² / (2C)

Capacitors in series:
  1/C_total = 1/C₁ + 1/C₂ + 1/C₃ + ...

Capacitors in parallel:
  C_total = C₁ + C₂ + C₃ + ...

RC charging/discharging:
  τ = R × C  (time constant)
  Charging:    V(t) = V₀(1 - e^(-t/τ))
  Discharging: V(t) = V₀ × e^(-t/τ)

ε₀ = 8.85 × 10⁻¹² F/m
```
