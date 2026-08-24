# Chapter 40: Kinetic Theory of Gases

> **"Temperature is not a mysterious property — it's just the average kinetic energy of molecules bouncing around. The kinetic theory connects the microscopic world of atoms to the macroscopic world we measure."**

---

## Table of Contents

- [40.1 The Molecular Picture of a Gas](#401-the-molecular-picture-of-a-gas)
- [40.2 Assumptions of Kinetic Theory](#402-assumptions-of-kinetic-theory)
- [40.3 Pressure from Molecular Collisions](#403-pressure-from-molecular-collisions)
- [40.4 Temperature and Kinetic Energy](#404-temperature-and-kinetic-energy)
- [40.5 The rms Speed](#405-the-rms-speed)
- [40.6 The Maxwell-Boltzmann Distribution](#406-the-maxwell-boltzmann-distribution)
- [40.7 Mean Free Path](#407-mean-free-path)
- [40.8 Degrees of Freedom and Energy Equipartition](#408-degrees-of-freedom-and-energy-equipartition)
- [40.9 Real Gases vs Ideal Gases](#409-real-gases-vs-ideal-gases)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 40.1 The Molecular Picture of a Gas

A gas is made of an enormous number of tiny molecules moving rapidly in all directions.

```
GAS IN A BOX:

  +------------------+
  |   •    •         |
  |  •   •      •    |
  |     •    •       |
  |  •       •   •   |
  |     •  •         |
  +------------------+

Each dot = one molecule
Moving fast, in random directions
Constantly colliding with walls and each other
```

Key numbers to appreciate the scale:
- 1 mole of gas contains 6.022 × 10²³ molecules (Avogadro's number)
- A molecule moves at roughly 400-1000 m/s
- A molecule collides roughly 10⁹ times per second
- The average distance between collisions is about 100 nm

Despite this chaos, the macroscopic properties (pressure, temperature, volume) are predictable — they are statistical averages over trillions of molecules.

---

## 40.2 Assumptions of Kinetic Theory

The **ideal gas model** makes these simplifying assumptions:

1. **Very many molecules**: enough for statistical averages to be valid

2. **Negligible volume**: the molecules themselves are tiny compared to the space between them (the gas is mostly empty space)

3. **Random motion**: molecules move in all directions with a range of speeds; on average, equal numbers move in each direction

4. **Elastic collisions**: when molecules collide with the walls or each other, no kinetic energy is lost (perfectly bouncy)

5. **No intermolecular forces**: molecules don't attract or repel each other (except during collisions)

```
IDEAL GAS MOLECULE MODEL:

  Hard sphere:      ( )
  
  No stickiness.
  No pull toward each other.
  Just bouncing off walls and each other.
```

These assumptions are approximate but work well for real gases at normal temperatures and pressures.

---

## 40.3 Pressure from Molecular Collisions

Pressure arises from molecules hitting the walls of the container.

Each collision delivers a tiny impulse. With trillions of collisions per second, the average force is constant and measurable as pressure.

### Derivation (Simplified)

Consider a box of side length L with N molecules. Each molecule of mass m has speed v.

```
MOLECULE BOUNCING BETWEEN WALLS:

  |<-------- L ------->|
  |                    |
  [wall] <---m,v---> [wall]
  
  Time between hits on same wall: t = 2L/v
  
  Impulse per hit: delta_p = 2mv (reverses direction)
  
  Force from one molecule: F = 2mv / (2L/v) = mv²/L
  
  For N molecules moving in 3 directions (x, y, z):
  Only N/3 contribute to each wall direction.
  
  Total force: F = N/3 × mv²/L
  
  Pressure: P = F/A = F/L² = Nmv²/(3L³) = Nmv²/(3V)
  
  So: PV = (1/3)Nmv²
```

This is the fundamental equation of kinetic theory:

```
PV = (1/3)Nm<v²>

where <v²> is the mean square speed (average of v² over all molecules)
```

---

## 40.4 Temperature and Kinetic Energy

Comparing the kinetic theory result to the ideal gas law (PV = NkT where k is Boltzmann's constant):

```
(1/3)Nm<v²> = NkT

(1/2)m<v²> = (3/2)kT

Average KE per molecule = (3/2) × k × T
```

This is the profound result of kinetic theory:

> **Temperature is a measure of the average kinetic energy of the molecules.**

```
Average KE per molecule = (3/2) kT

where:
  k = 1.38 × 10⁻²³ J/K  (Boltzmann's constant)
  T = absolute temperature in Kelvin
```

Key implications:

- At T = 0 K (absolute zero), molecular KE = 0 — molecules stop moving
- Higher temperature = faster molecules
- All gases (regardless of type) have the SAME average KE at the same temperature
- Heavy molecules move slower (same KE, more mass → less speed)

### Worked Example 40.1

Find the average kinetic energy of a gas molecule at room temperature (T = 293 K).

**Solution:**

KE = (3/2) kT = (3/2) × 1.38 × 10⁻²³ × 293 = **6.07 × 10⁻²¹ J**

This is tiny, but with 10²³ molecules, the total energy is enormous.

---

## 40.5 The rms Speed

From (1/2)m<v²> = (3/2)kT:

```
<v²> = 3kT/m

v_rms = sqrt(<v²>) = sqrt(3kT/m)

where:
  v_rms = root mean square speed (m/s)
  k = 1.38 × 10⁻²³ J/K
  T = temperature (K)
  m = mass of one molecule (kg)
```

Using molar mass M (kg/mol) and gas constant R = kN_A:

```
v_rms = sqrt(3RT/M)

where:
  R = 8.314 J/(mol·K)
  M = molar mass (kg/mol)
```

### Worked Example 40.2

Find the rms speed of nitrogen molecules (M = 0.028 kg/mol) at 25°C.

**Solution:**

T = 25 + 273 = 298 K

v_rms = sqrt(3 × 8.314 × 298 / 0.028) = sqrt(266,000) ≈ **516 m/s**

This is about 1800 km/h — nitrogen molecules move incredibly fast!

### Comparison at Same Temperature

```
MOLECULE    Molar Mass (g/mol)    v_rms at 300K
----------  ------------------    --------------
H₂          2                     1934 m/s
He          4                     1368 m/s
N₂          28                    517 m/s
O₂          32                    484 m/s
CO₂         44                    412 m/s

Lighter molecules = faster speeds (same average KE).
```

This explains why hydrogen leaks out of balloons faster than helium.

---

## 40.6 The Maxwell-Boltzmann Distribution

Not all molecules in a gas move at the same speed. They have a **distribution** of speeds.

```
NUMBER OF MOLECULES vs SPEED:

 N(v)
   |
   |       *
   |      * *
   |     *   *
   |    *     *
   |   *       *
   |  *         *
   | *            *
   |*               *          *
   +----------------------------> speed v
         ^   ^     ^
         |   |     |
       v_p  v_avg  v_rms
       (most (mean) (rms)
       probable)
```

Features of the Maxwell-Boltzmann distribution:
- **Peak (most probable speed, v_p)**: the speed most molecules have
- **Mean speed (v_avg)**: arithmetic average
- **rms speed (v_rms)**: sqrt of average of v²

Relationships:
- v_p < v_avg < v_rms
- v_p = sqrt(2RT/M) ≈ 0.816 v_rms
- v_avg = sqrt(8RT/πM) ≈ 0.921 v_rms

### Effect of Temperature

```
LOW T:    narrow peak at low speed
          
  N(v)
   |  *
   | * *
   |*   **
   |      *****
   +---------> v

HIGH T:   wide, flat distribution at higher speed

  N(v)
   |    *
   |   * *
   |  *   *
   | *     ***
   |*         *******
   +---------> v
```

As temperature increases, the distribution broadens and shifts to higher speeds.

---

## 40.7 Mean Free Path

A molecule doesn't travel far before hitting another one.

The **mean free path (λ)** is the average distance a molecule travels between successive collisions:

```
λ = 1 / (sqrt(2) × π × d² × n)

where:
  d = molecular diameter
  n = number density (molecules per m³)
```

### Worked Example 40.3

For air at standard conditions:
- d ≈ 3 × 10⁻¹⁰ m
- n ≈ 2.7 × 10²⁵ molecules/m³

λ = 1 / (sqrt(2) × π × (3×10⁻¹⁰)² × 2.7×10²⁵)
  = 1 / (sqrt(2) × π × 9×10⁻²⁰ × 2.7×10²⁵)
  ≈ **68 nm** (about 200 molecular diameters)

At high altitude (lower pressure, lower n): mean free path is much larger.

---

## 40.8 Degrees of Freedom and Energy Equipartition

The **equipartition theorem** states: each degree of freedom has average energy (1/2)kT.

```
DEGREE OF FREEDOM: an independent way a molecule can store energy

MONATOMIC GAS (He, Ne, Ar):
  3 translational degrees (moving in x, y, z)
  Total KE per molecule = 3 × (1/2)kT = (3/2)kT  ✓ (matches our result)

DIATOMIC GAS (N₂, O₂):
  3 translational + 2 rotational = 5 degrees
  Total energy per molecule = (5/2)kT

LINEAR POLYATOMIC (CO₂):
  3 translational + 2 rotational = 5 degrees
  Plus vibrational at high temperatures

NON-LINEAR POLYATOMIC (H₂O):
  3 translational + 3 rotational = 6 degrees
  Total energy per molecule = (6/2)kT = 3kT
```

This affects the **heat capacity** of gases — diatomic gases need more energy to raise temperature because they store energy in rotation too.

---

## 40.9 Real Gases vs Ideal Gases

Ideal gas assumptions break down at:
- **High pressure**: molecules are packed close — molecular volume matters
- **Low temperature**: molecules move slowly — intermolecular attractions matter

### The Van der Waals Equation

A more accurate equation for real gases:

```
(P + a/V²)(V - b) = nRT

where:
  a = accounts for intermolecular attractions (reduces effective pressure)
  b = accounts for molecular volume (reduces effective volume)
```

At high temperature and low pressure (molecules far apart, moving fast), a/V² → 0 and b → 0, and this reduces to the ideal gas law.

### Behavior of Real Gases

```
COMPRESSIBILITY FACTOR Z = PV/(nRT):

For ideal gas: Z = 1 always

For real gas:
  Z > 1 at very high pressure (repulsion dominates)
  Z < 1 at moderate pressure (attraction dominates)
  Z → 1 at low pressure (ideal behavior)
```

---

## Summary

- **Kinetic theory**: gas = vast number of tiny molecules in rapid random motion
- **Pressure**: caused by molecules hitting walls; P = (1/3)Nmv²_rms / V
- **Temperature**: proportional to average kinetic energy per molecule: KE = (3/2)kT
- **rms speed**: v_rms = sqrt(3kT/m) = sqrt(3RT/M); lighter molecules are faster
- **Maxwell-Boltzmann distribution**: range of molecular speeds; shifts right and broadens with T
- **Mean free path**: average distance between collisions ≈ 68 nm in air at STP
- **Equipartition**: each degree of freedom stores (1/2)kT energy
- Monatomic: 3 degrees, diatomic: 5 degrees
- **Real gases** deviate from ideal at high pressure or low temperature

---

## Key Equations

```
Fundamental kinetic theory equation:
  PV = (1/3) × N × m × <v²>

Average KE per molecule:
  KE = (3/2) × k × T
  k = 1.38 × 10⁻²³ J/K

rms speed:
  v_rms = sqrt(3kT/m) = sqrt(3RT/M)
  R = 8.314 J/(mol·K)

Maxwell-Boltzmann most probable speed:
  v_p = sqrt(2kT/m)

Mean free path:
  λ = 1 / (sqrt(2) × π × d² × n)

Equipartition:
  Energy per degree of freedom = (1/2)kT
  Monatomic total KE = (3/2)kT
  Diatomic total KE  = (5/2)kT

Avogadro's number:
  N_A = 6.022 × 10²³ mol⁻¹
  R = k × N_A
```
