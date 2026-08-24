# Chapter 57: Lenz's Law and Inductors

> **"Lenz's Law says that nature resists change. An inductor is a device that embodies this principle — it opposes any change in the current flowing through it, like a mechanical flywheel resists changes in rotation."**

---

## Table of Contents

- [57.1 Electromagnetic Induction Recap](#571-electromagnetic-induction-recap)
- [57.2 Lenz's Law](#572-lenzs-law)
- [57.3 Eddy Currents](#573-eddy-currents)
- [57.4 Inductance](#574-inductance)
- [57.5 Self-Inductance of a Solenoid](#575-self-inductance-of-a-solenoid)
- [57.6 Energy Stored in an Inductor](#576-energy-stored-in-an-inductor)
- [57.7 RL Circuits](#577-rl-circuits)
- [57.8 Mutual Inductance and Transformers](#578-mutual-inductance-and-transformers)
- [57.9 Applications of Inductors](#579-applications-of-inductors)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 57.1 Electromagnetic Induction Recap

Faraday's Law (from the previous chapter): a changing magnetic flux through a circuit induces an EMF in that circuit.

```
EMF = -dΦ/dt = -d(BA)/dt

Flux Φ = B × A × cos(θ)
```

The magnitude of the induced EMF equals the rate of change of flux.

The minus sign leads us to Lenz's Law.

---

## 57.2 Lenz's Law

**Lenz's Law:** The direction of the induced current is always such that it opposes the change in magnetic flux that caused it.

Or more memorably: **the induced effect opposes the cause.**

```
LENZ'S LAW IN ACTION:

Case 1: MAGNET APPROACHING (flux increasing into loop):

  N pole →→→→ [loop]
  
  Flux into loop is INCREASING.
  
  Induced current must create B that OPPOSES increase → 
  induced B points OUT of loop → current is COUNTER-CLOCKWISE (right-hand rule).
  
  The loop acts like a magnet with N pole facing the approaching bar magnet.
  → REPULSION (opposes the magnet coming closer).

Case 2: MAGNET RECEDING (flux decreasing into loop):

  N pole ←←←← [loop]
  
  Flux into loop is DECREASING.
  
  Induced current must CREATE flux to maintain it → 
  induced B points INTO loop → current is CLOCKWISE.
  
  The loop acts like an S pole facing the receding N → ATTRACTION 
  (opposes the magnet moving away).
```

In both cases, the induced current creates a force that OPPOSES the motion of the magnet. Work must be done to move the magnet against this force.

### Lenz's Law and Energy Conservation

Lenz's Law IS energy conservation. If the induced current helped the change (instead of opposing it), you'd get more energy out than you put in — violating conservation of energy.

```
ANALOGY:
  Push a magnet toward a loop.
  Lenz's Law creates repulsion → you must push harder.
  You do work against repulsion.
  That work appears as electrical energy (and then heat) in the circuit.
  
  Energy in = Energy out. Conservation is satisfied.
```

---

## 57.3 Eddy Currents

When a conducting material (not just a thin wire) moves through a magnetic field, currents are induced throughout the volume of the conductor. These are **eddy currents** (also called **Foucault currents**).

```
METAL PLATE MOVING THROUGH FIELD:

      N  [field region]  S
      
  →→→→ [solid metal plate] →→→→  (moving right)
  
  As the plate enters the field:
    eddy currents induced in plate
    eddy currents create B opposing entry (Lenz's Law)
    → braking force on plate (opposes motion)
    
  The plate slows down. Kinetic energy → heat in the plate.
```

### Effects and Applications

**Damping (useful):**
- **Magnetic braking**: trains use eddy current brakes — silent, no wear, powerful
- **Damping in meters**: moving-coil instruments have aluminum frames to damp oscillations

**Unwanted heating (harmful):**
- Transformer cores have eddy currents → heat loss
- Solution: **laminate the core** (thin layers of iron insulated from each other — interrupts eddy current loops)

```
LAMINATED vs SOLID CORE:

SOLID CORE:              LAMINATED CORE:
  
  Large eddy current       Eddy currents confined
  loops through core       to thin layers → much smaller
  → lots of heating        → much less heating
  
  +---------+             |  |  |  |  |  |
  |   ~~~~  |             |~~|~~|~~|~~|~~|
  |         |             |  |  |  |  |  |
  +---------+             |  |  |  |  |  |
```

---

## 57.4 Inductance

**Inductance** is the property of a circuit element that opposes changes in current.

When current changes in a coil, the changing magnetic flux induces an EMF in the same coil that opposes the change. This is **self-inductance**.

```
SELF-INDUCTION:
  
  Increasing current → increasing flux → induced EMF opposes increase
  Decreasing current → decreasing flux → induced EMF opposes decrease
  
  The inductor is like an electrical flywheel — resists current changes.
```

The induced EMF (back EMF) is:

```
EMF = -L × (dI/dt)

where:
  L = self-inductance (Henry, H)
  dI/dt = rate of change of current (A/s)
```

The Henry is a large unit: 1 H = 1 Wb/A = 1 V·s/A

Practical inductors: μH (microhenry) to H (Henry) range.

---

## 57.5 Self-Inductance of a Solenoid

For a solenoid with N turns, length ℓ, cross-section area A:

```
L = μ₀ × n² × A × ℓ

or equivalently:

L = μ₀ × N² × A / ℓ

where n = N/ℓ = turns per unit length
```

For an iron-core solenoid: replace μ₀ with μ = μ₀ × μ_r (μ_r can be thousands).

### Worked Example 57.1

A solenoid has 500 turns, length 20 cm, diameter 2 cm. Find L.

**Solution:**

A = π(0.01)² = 3.14 × 10⁻⁴ m²

L = μ₀N²A/ℓ = 4π×10⁻⁷ × 500² × 3.14×10⁻⁴ / 0.2

L = 4π×10⁻⁷ × 250,000 × 3.14×10⁻⁴ / 0.2

L ≈ **493 μH**

---

## 57.6 Energy Stored in an Inductor

An inductor carrying current stores energy in its magnetic field:

```
E = (1/2) × L × I²

where:
  E = energy stored (J)
  L = inductance (H)
  I = current (A)
```

Compare to capacitor: E = (1/2)CV² (energy in electric field)

The energy is stored in the **magnetic field** inside the inductor:

```
Energy density in magnetic field:
  u = B² / (2μ₀)   (J/m³)
```

### Worked Example 57.2

An inductor of 50 mH carries 2 A. Find the energy stored.

E = (1/2)LI² = (1/2) × 50×10⁻³ × 4 = **0.1 J**

---

## 57.7 RL Circuits

When a battery, resistor, and inductor are in series, the current doesn't instantly reach its final value — the inductor opposes the change.

```
SWITCHING ON (RL circuit):

   + --- R --- L --- (battery V₀) --- +
   |                                  |
   +----------------------------------+
   
   At t=0: switch closed.
   At t=0: I=0 (inductor resists sudden change)
   At t=∞: I = V₀/R (steady state, inductor acts like a wire)
```

Current growth:

```
I(t) = (V₀/R) × (1 - e^(-t/τ))

Time constant: τ = L/R

After t = τ: 63% of final current
After t = 5τ: 99.3% of final current (essentially steady)
```

Current decay (when battery disconnected):

```
I(t) = I₀ × e^(-t/τ)
```

```
CURRENT vs TIME (RL circuit switching on):

I
I_max|               _____________
     |             /
0.63 |           *  <- at t=τ
I_max|         /
     |       /
     |     /
     |   /
  0  +-----------> t
          τ

(Same shape as capacitor charging, with τ = L/R instead of RC)
```

### Worked Example 57.3

A 12 V battery is connected to a 100 Ω resistor and 50 mH inductor in series.

(a) Find the time constant.
(b) Find the current after 1 ms.
(c) Find the final current.

**Solution:**

(a) τ = L/R = 0.050/100 = **5 × 10⁻⁴ s = 0.5 ms**

(b) I = (V₀/R)(1 - e^(-t/τ)) = (12/100)(1 - e^(-1/0.5)) = 0.12(1 - e^(-2)) = 0.12(1 - 0.135) = **0.104 A**

(c) I_final = V₀/R = 12/100 = **0.12 A**

---

## 57.8 Mutual Inductance and Transformers

**Mutual inductance**: when two coils are near each other, changing current in one induces EMF in the other.

```
MUTUAL INDUCTION:

  Coil 1 (primary)    Coil 2 (secondary)
  
   I₁ changing   →   [iron core]   →   EMF₂ induced
   
   EMF₂ = -M × (dI₁/dt)
   
   M = mutual inductance (Henry)
```

This is the basis of the **transformer**:

```
TRANSFORMER:

  Primary: N₁ turns, voltage V₁
  Secondary: N₂ turns, voltage V₂
  
  V₂/V₁ = N₂/N₁
  
  Step-up: N₂ > N₁ → V₂ > V₁
  Step-down: N₂ < N₁ → V₂ < V₁
  
  For ideal transformer (no losses):
  P₁ = P₂ → V₁I₁ = V₂I₂ → I₂/I₁ = N₁/N₂
  
  Higher voltage → lower current (and vice versa).
```

### Worked Example 57.4

A transformer has 2000 primary turns and 100 secondary turns. Primary voltage = 240 V.

(a) Secondary voltage?
(b) If secondary current = 5 A, primary current?

**Solution:**

(a) V₂ = V₁ × N₂/N₁ = 240 × 100/2000 = **12 V** (step-down)

(b) I₁ = I₂ × N₂/N₁ = 5 × 100/2000 = **0.25 A**

---

## 57.9 Applications of Inductors

### Filters (Chokes)

Inductors pass DC and block AC (high-frequency signals). Used to "choke" high-frequency interference from power supplies.

```
Inductive reactance: X_L = ωL = 2πfL

Higher frequency → higher impedance → inductor blocks it more.
```

### Switched Mode Power Supplies (SMPS)

All modern phone chargers, laptop chargers use inductors. The inductor stores energy during each switching cycle and smooths out the current.

### Fluorescent Lamps

The ballast (choke) in fluorescent lights limits current and provides the high voltage spike needed to start the discharge.

### Metal Detectors

The detector coil creates an alternating magnetic field. Metal objects create eddy currents which change the inductance — detected by the circuit.

### Induction Cooktop

An alternating current in a coil under the glass top creates a rapidly changing magnetic field. This induces eddy currents in the iron pan → pan heats up. The glass top stays cool!

---

## Summary

- **Lenz's Law**: induced current opposes the change in flux that caused it; consequence of energy conservation
- **Eddy currents**: induced in bulk conductors moving through fields; cause braking/heating; reduced by lamination
- **Inductance L**: opposes changes in current; EMF = -L(dI/dt); unit Henry (H)
- **Solenoid**: L = μ₀N²A/ℓ; increased by iron core (×μ_r)
- **Energy in inductor**: E = ½LI² (stored in magnetic field)
- **RL circuit**: τ = L/R; I rises/falls exponentially with this time constant
- **Mutual inductance**: changing current in one coil induces EMF in nearby coil
- **Transformer**: V₂/V₁ = N₂/N₁; ideal: V₁I₁ = V₂I₂
- Applications: power supplies, filters, metal detectors, induction cooking, transformers

---

## Key Equations

```
Lenz's Law (direction rule):
  Induced current opposes the change in flux

Faraday's Law:
  EMF = -N × dΦ/dt

Self-inductance EMF:
  EMF = -L × dI/dt

Solenoid inductance:
  L = μ₀N²A/ℓ = μ₀n²Aℓ

Energy in inductor:
  E = (1/2) × L × I²

RL time constant:
  τ = L/R

RL current growth:
  I(t) = (V₀/R)(1 - e^(-t/τ))

RL current decay:
  I(t) = I₀ × e^(-t/τ)

Transformer:
  V₂/V₁ = N₂/N₁
  I₂/I₁ = N₁/N₂  (ideal)
  V₁I₁ = V₂I₂    (ideal, no losses)

Inductive reactance:
  X_L = 2πfL  (Ω)
```
