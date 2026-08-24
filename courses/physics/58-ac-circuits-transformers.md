# Chapter 58: AC Circuits and Transformers

> **"Alternating current powers the world — every light bulb, motor, and device in your home runs on it. Understanding AC requires understanding how resistors, capacitors, and inductors each respond differently to oscillating voltages."**

---

## Table of Contents

- [58.1 AC Basics: Frequency, Phase, and RMS](#581-ac-basics-frequency-phase-and-rms)
- [58.2 Resistors in AC Circuits](#582-resistors-in-ac-circuits)
- [58.3 Capacitors in AC Circuits](#583-capacitors-in-ac-circuits)
- [58.4 Inductors in AC Circuits](#584-inductors-in-ac-circuits)
- [58.5 Series RLC Circuit](#585-series-rlc-circuit)
- [58.6 Resonance in RLC Circuits](#586-resonance-in-rlc-circuits)
- [58.7 Power in AC Circuits](#587-power-in-ac-circuits)
- [58.8 Transformers](#588-transformers)
- [58.9 The Power Grid](#589-the-power-grid)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 58.1 AC Basics: Frequency, Phase, and RMS

**Alternating current (AC)** is current that reverses direction periodically.

```
AC VOLTAGE:  v(t) = V₀ × sin(ωt)

V₀ = peak (maximum) voltage
ω  = 2πf = angular frequency (rad/s)
f  = frequency (Hz)
T  = 1/f = period (s)
```

```
v(t)
V₀ |   *     *
   |  * *   * *
   | *   * *   *
   |*     *     *
---+--------------> t
   |*           *
  -V₀
```

Household AC:
- India/UK: 230 V rms, 50 Hz
- USA: 120 V rms, 60 Hz

### RMS Values

For sinusoidal AC:

```
V_rms = V₀ / √2 ≈ 0.707 × V₀

I_rms = I₀ / √2 ≈ 0.707 × I₀

"230 V mains" means V_rms = 230 V, so V₀ = 230√2 ≈ 325 V
```

RMS values are used because they give the **equivalent DC value** for power calculations:

```
P = V_rms × I_rms = V₀I₀/2   (for resistive load)
```

---

## 58.2 Resistors in AC Circuits

For a resistor R, Ohm's Law applies at every instant:

```
v(t) = R × i(t)

If v = V₀ sin(ωt):
   i = (V₀/R) sin(ωt) = I₀ sin(ωt)

Voltage and current are IN PHASE (peaks at same time).
```

```
v and i vs time:
  
  v| * (voltage)
   |* *
   |  * (current, same phase)
   |   *
   +--------> t
   
  No phase difference!
```

The **impedance** of a resistor (AC analog of resistance): Z_R = R

Power dissipated: P = I²_rms × R = V²_rms / R (same as DC)

---

## 58.3 Capacitors in AC Circuits

A capacitor charges and discharges with the AC voltage. Current flows as long as voltage is changing:

```
i = C × dv/dt

If v = V₀ sin(ωt):
   i = C × ωV₀ cos(ωt) = I₀ cos(ωt)

Current LEADS voltage by 90° (current peaks before voltage does).
```

```
v and i vs time:

  i| *                     <- current leads by 90°
   |* *     *
   |  * * * *    <- voltage
   |    *
   +--------> t
```

### Capacitive Reactance

The AC "resistance" of a capacitor:

```
X_C = 1 / (ωC) = 1 / (2πfC)

Units: Ohms (Ω)
```

Key property: **higher frequency → lower reactance** (capacitor passes high-frequency AC more easily)

- At DC (f = 0): X_C = ∞ (capacitor blocks DC completely)
- At high frequency: X_C → 0 (capacitor acts like a short circuit)

```
X_C vs frequency:
  
  X_C
  |*
  | *
  |  **
  |    **
  |       ***
  |           *****
  +----------------> f
  
  (hyperbola — decreasing with frequency)
```

---

## 58.4 Inductors in AC Circuits

An inductor opposes changes in current. The voltage across it is proportional to the rate of current change:

```
v = L × di/dt

If i = I₀ sin(ωt):
   v = L × I₀ω cos(ωt) = V₀ cos(ωt)

Voltage LEADS current by 90° (or current lags voltage by 90°).
```

```
v and i vs time:

  v| *          <- voltage leads
   |* *
   |   * * *    <- current
   |       *
   +--------> t
```

### Inductive Reactance

```
X_L = ωL = 2πfL

Units: Ohms (Ω)
```

Key property: **higher frequency → higher reactance** (inductor blocks high-frequency AC)

- At DC (f = 0): X_L = 0 (inductor acts like a plain wire)
- At high frequency: X_L → ∞ (inductor blocks high-frequency signals)

```
MEMORY AID:
  "ELI the ICE man"
  
  E = EMF (voltage), L = inductor, I = current
  In an inductor (L): E (voltage) leads I (current) → ELI
  
  I = current, C = capacitor, E = EMF (voltage)
  In a capacitor (C): I (current) leads E (voltage) → ICE
```

---

## 58.5 Series RLC Circuit

A series circuit with resistor R, inductor L, and capacitor C:

```
  ----R----L----C----
  |                 |
  |   AC source     |
  ---v(t)---
```

The total impedance:

```
Z = √(R² + (X_L - X_C)²)

where:
  X_L = ωL   (inductive reactance)
  X_C = 1/ωC (capacitive reactance)
```

The current:

```
I = V / Z

I_rms = V_rms / Z
```

Phase angle between voltage and current:

```
tan(φ) = (X_L - X_C) / R

If X_L > X_C: inductive (current lags voltage, φ > 0)
If X_L < X_C: capacitive (current leads voltage, φ < 0)
If X_L = X_C: purely resistive (current in phase with voltage, φ = 0) → RESONANCE
```

### Phasor Diagram

A **phasor diagram** shows voltages as vectors (phasors) rotating in the complex plane:

```
PHASOR DIAGRAM for series RLC:

     V_L     ^
      |      | V_L - V_C
      |      |
 -----+-------> V_R (reference, along real axis)
      |
      | V_C
      v

V_total = √(V_R² + (V_L - V_C)²)
```

---

## 58.6 Resonance in RLC Circuits

When X_L = X_C:

```
ωL = 1/(ωC)

ω₀ = 1/√(LC)

f₀ = 1/(2π√(LC))
```

At resonance:
- Impedance Z = R (minimum)
- Current is maximum (V/R)
- Phase angle φ = 0 (purely resistive)
- V_L and V_C cancel each other

```
IMPEDANCE vs FREQUENCY:

Z
|        R (minimum at resonance)
|       *|*
|      * | *
|     *  |  *
|    *   |   *
|   *    |    *
|  *     |     *
| *      |      *
|*       |       *
+--------+---------> f
          f₀

Below f₀: capacitive (X_C dominates)
Above f₀: inductive (X_L dominates)
At f₀: minimum Z, maximum current
```

### Worked Example 58.1

An RLC series circuit: R = 100 Ω, L = 0.1 H, C = 10 μF, V_rms = 50 V.

(a) Find the resonant frequency.
(b) Find the current at resonance.

**Solution:**

(a) f₀ = 1/(2π√LC) = 1/(2π√(0.1 × 10⁻⁵)) = 1/(2π × 10⁻³) ≈ **159 Hz**

(b) At resonance, Z = R = 100 Ω
    I_rms = V_rms/R = 50/100 = **0.5 A**

---

## 58.7 Power in AC Circuits

In AC circuits with reactance (capacitors/inductors), voltage and current are out of phase. Average power depends on the phase difference:

```
P_avg = V_rms × I_rms × cos(φ)

cos(φ) = power factor
```

- **Purely resistive** (R only): cos(φ) = 1, P = V_rms × I_rms (maximum power)
- **Purely inductive** (L only): cos(φ) = 0, P = 0 (inductor stores/returns energy, net power = 0)
- **Purely capacitive** (C only): cos(φ) = 0, P = 0 (capacitor stores/returns energy, net power = 0)

```
REACTIVE POWER (VAR):
  In inductors and capacitors, energy oscillates in and out of the device.
  This is called reactive power Q_r (in VA reactive, VAR).
  
  Apparent power: S = V_rms × I_rms (VA)
  Real power: P = S × cos(φ) (W) = useful power
  Reactive power: Q_r = S × sin(φ) (VAR) = "wasted" power
  
  S² = P² + Q_r²
```

---

## 58.8 Transformers

A transformer is an AC device that changes voltage levels using mutual inductance between two coils.

```
TRANSFORMER STRUCTURE:

  Primary coil (N₁ turns)    Secondary coil (N₂ turns)
  
  AC input V₁ →→→→ [iron core] →→→→ V₂ output
  
  Changing current in primary → changing flux in core →
  induces EMF in secondary
```

### Transformer Equations

```
V₂/V₁ = N₂/N₁   (turns ratio)

For ideal transformer (no losses):
  P_in = P_out
  V₁I₁ = V₂I₂
  I₂/I₁ = N₁/N₂  (current inversely proportional to turns)
```

### Types of Transformers

```
STEP-UP TRANSFORMER: N₂ > N₁
  V₂ > V₁, I₂ < I₁
  Used at power stations (increase voltage for transmission)

STEP-DOWN TRANSFORMER: N₂ < N₁
  V₂ < V₁, I₂ > I₁
  Used to reduce voltage for homes (230V → 5V for phone chargers)

ISOLATION TRANSFORMER: N₁ = N₂
  V₂ = V₁ (no voltage change)
  Used to break galvanic connection while transmitting power
  (safety/interference reduction)
```

### Real Transformer Losses

Real transformers are not 100% efficient due to:
1. **Copper losses**: I²R heating in coil windings
2. **Iron losses**: eddy currents and hysteresis in the iron core
3. **Flux leakage**: some magnetic flux doesn't couple both coils

Modern power transformers have >99% efficiency.

### Worked Example 58.2

A step-down transformer has 4000 primary turns and 200 secondary turns.

Primary: 240 V AC at 0.5 A.

(a) Secondary voltage?
(b) Secondary current (assume ideal)?
(c) Power transferred?

**Solution:**

(a) V₂ = V₁ × N₂/N₁ = 240 × 200/4000 = **12 V**

(b) I₂ = I₁ × N₁/N₂ = 0.5 × 4000/200 = **10 A**

(c) P = V₁I₁ = 240 × 0.5 = **120 W** (or V₂I₂ = 12 × 10 = 120 W ✓)

---

## 58.9 The Power Grid

The AC power grid uses transformers at multiple stages to efficiently transmit power:

```
POWER GRID:

GENERATION (~11 kV, 50 Hz)
    |
    v  Step-up transformer (×20)
    |
GRID TRANSMISSION (220-400 kV)
    | 
    v  Step-down transformer (÷10)
    |
REGIONAL DISTRIBUTION (33 kV)
    |
    v  Step-down transformer (÷3)
    |
LOCAL SUBSTATION (11 kV)
    |
    v  Step-down transformer (÷50)
    |
HOME (230 V)
```

Why do we use AC (not DC) for power transmission?
- **AC can be transformed** (transformers only work with AC)
- DC would require power electronics (transistors, converters) for voltage changes
- Historical: in the late 1800s, Tesla's AC system won the "War of Currents" against Edison's DC system

Modern HVDC (High Voltage DC) transmission is now competitive for very long distances (>600 km) due to lower losses, but transformers still use AC.

---

## Summary

- **AC voltage**: v(t) = V₀sin(ωt); V_rms = V₀/√2; frequency f = ω/2π
- **Resistor**: voltage and current in phase; Z = R; P = I²_rms × R
- **Capacitor**: current leads voltage by 90°; X_C = 1/(ωC); blocks DC, passes high-frequency AC
- **Inductor**: voltage leads current by 90°; X_L = ωL; blocks high-frequency AC, passes DC
- **RLC series**: Z = √(R² + (X_L - X_C)²); phase angle φ = arctan((X_L - X_C)/R)
- **Resonance**: at ω₀ = 1/√(LC); Z = R (minimum); maximum current
- **AC power**: P = V_rms × I_rms × cos(φ); power factor = cos(φ)
- **Transformer**: V₂/V₁ = N₂/N₁; I₂/I₁ = N₁/N₂; P_in = P_out (ideal)
- Power grid: step up for transmission, step down for distribution

---

## Key Equations

```
AC voltage/current:
  v(t) = V₀ sin(ωt)
  V_rms = V₀/√2
  ω = 2πf

Reactances:
  X_C = 1/(ωC)   (decreases with frequency)
  X_L = ωL       (increases with frequency)

RLC impedance:
  Z = √(R² + (X_L - X_C)²)
  I = V/Z

Resonant frequency:
  ω₀ = 1/√(LC)
  f₀ = 1/(2π√(LC))

AC power:
  P = V_rms I_rms cos(φ)
  cos(φ) = R/Z = power factor

Transformer:
  V₂/V₁ = N₂/N₁
  I₁V₁ = I₂V₂  (ideal)
```
