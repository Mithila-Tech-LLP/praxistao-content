# Chapter 47: Electric Potential and Voltage

> **"Voltage is to electricity what height is to gravity. A stone falls from high to low; current flows from high to low voltage. The analogy is almost perfect."**

---

## Table of Contents

- [47.1 Electric Potential Energy](#471-electric-potential-energy)
- [47.2 Electric Potential (Voltage)](#472-electric-potential-voltage)
- [47.3 Potential Difference](#473-potential-difference)
- [47.4 Equipotential Surfaces](#474-equipotential-surfaces)
- [47.5 Potential Due to a Point Charge](#475-potential-due-to-a-point-charge)
- [47.6 Relationship Between E and V](#476-relationship-between-e-and-v)
- [47.7 Work Done by Electric Force](#477-work-done-by-electric-force)
- [47.8 Electron Volt](#478-electron-volt)
- [47.9 Capacitors and Stored Energy](#479-capacitors-and-stored-energy)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 47.1 Electric Potential Energy

A charge in an electric field has **electric potential energy**, just as a mass in a gravitational field has gravitational potential energy.

### Analogy with Gravity

```
GRAVITY:                         ELECTRICITY:

mass m                           positive charge +q
at height h                      near positive charge Q
  
  h |  ↑ (against gravity)         ↑ (against repulsion)
    |  |                           |
    |  m                           +q
    |
  ground (reference)            infinity (reference)
  
PE_grav = mgh                   PE_elec = kqQ/r

Lift mass up: do work,          Push charge toward + source:
PE increases                    do work, PE increases
```

Both potential energies are stored energy that can be released.

### Sign of Electric PE

```
SAME-SIGN CHARGES (repulsion):
  Bringing them together (decreasing r) INCREASES PE
  (like compressing a spring)
  
OPPOSITE-SIGN CHARGES (attraction):
  Bringing them together DECREASES PE
  (like dropping a mass — releases energy)
```

---

## 47.2 Electric Potential (Voltage)

The **electric potential** V at a point is the electric potential energy per unit charge:

```
V = PE / q

Or: PE = q × V

Units: Volts (V) = Joules per Coulomb (J/C)
```

Just as electric field E = F/q (force per unit charge), potential V = PE/q (energy per unit charge).

Electric potential is a **scalar** (not a vector) — much easier to work with than the electric field in many situations.

### Absolute Potential

The absolute potential is defined with reference to infinity (where PE = 0):

```
V at a point = work done per unit charge to bring a positive test charge
               from infinity to that point
```

Usually we only care about **potential differences** (see next section).

---

## 47.3 Potential Difference

The **potential difference** between two points A and B is:

```
V_AB = V_A - V_B = W_AB / q

where W_AB = work done by the electric force in moving charge q from A to B.
```

This is what we measure with a voltmeter and call **voltage**.

### The Gravity Analogy Revisited

```
GRAVITY:                         ELECTRICITY:
  
  Ground floor                     Low potential (-)
  |                                |
  |                                |
  |  mgh = potential               |  qV = potential
  |  energy                        |  energy
  |                                |
  Top floor                        High potential (+)
  
A ball falls from top to ground:  A positive charge "falls" from high
releases mgh of energy             to low potential; releases qΔV energy.
                                   
A ball must be LIFTED to go up:   A positive charge must be PUSHED against
requires mgh of work               the field to go from low to high; requires work.
```

If a positive charge moves from high to low potential:
- It moves in the direction of E
- Electric force does positive work
- KE increases (energy released)

---

## 47.4 Equipotential Surfaces

An **equipotential surface** is a surface where all points have the same potential. No work is done in moving a charge along an equipotential surface.

```
EQUIPOTENTIAL LINES AROUND A POINT CHARGE:

  (concentric circles, like contour lines on a map)
  
       V=100V
    V=50V   .|.   V=25V
          .|   |.
         .|     |.
         |   +   |
         .|     |.
          .|   |.
            .|.
            
  Closer to + charge → higher potential
  Further away → lower potential
```

### Relationship Between Equipotentials and Field Lines

```
Field lines and equipotentials are ALWAYS PERPENDICULAR:

  →→→→→→→   (field lines, pointing right)
  
  |  |  |   (equipotential surfaces, vertical lines)
  |  |  |   
  |  |  |   
  
  Field is the gradient of potential.
  Moving along equipotential = zero work (perpendicular to force).
```

This is exactly analogous to a contour map: contour lines = equipotentials, direction of steepest slope = field direction.

---

## 47.5 Potential Due to a Point Charge

From the definition and Coulomb's law, the electric potential at distance r from point charge Q:

```
V = k × Q / r

where:
  V = electric potential (V)
  k = 8.99 × 10⁹ N·m²/C²
  Q = source charge (C, positive or negative)
  r = distance from charge (m)
```

Note:
- For positive Q: V is positive and decreases with distance
- For negative Q: V is negative and becomes less negative with distance
- V → 0 as r → infinity

### Superposition for Multiple Charges

```
V_total = V₁ + V₂ + V₃ + ...   (scalar addition — much easier than vector!)
        = k(Q₁/r₁ + Q₂/r₂ + Q₃/r₃ + ...)
```

Because V is a scalar, adding contributions is simpler than adding electric field vectors.

### Worked Example 47.1

Find the potential at a point 0.3 m from charge +5 μC and 0.4 m from charge -3 μC.

**Solution:**

V₁ = kQ₁/r₁ = 8.99×10⁹ × 5×10⁻⁶ / 0.3 = **149,833 V ≈ 150 kV**

V₂ = kQ₂/r₂ = 8.99×10⁹ × (-3×10⁻⁶) / 0.4 = **-67,425 V ≈ -67.4 kV**

V_total = 150,000 + (-67,400) = **82,600 V ≈ 82.6 kV**

---

## 47.6 Relationship Between E and V

The electric field and potential are closely related:

```
E = -ΔV / Δx   (in one dimension)

Or in 3D: E = -grad(V)  (E is the negative gradient of V)

For uniform field:
  E = V / d   (where d is the separation and V is the potential difference)
```

The minus sign means: E points from high V to low V (downhill, like gravity).

```
POTENTIAL MAP:
  
  V = 100V    V = 60V    V = 20V
     |            |          |
     |    →  E   |    →  E  |   → E
     |            |          |
     
  Field points from high to low potential.
  Equipotential surfaces are perpendicular to E.
```

### Worked Example 47.2

The potential difference between two parallel plates is 300 V, and they are 5 cm apart.

Find the electric field between them.

**Solution:**

E = V/d = 300 / 0.05 = **6000 V/m = 6 kN/C**

---

## 47.7 Work Done by Electric Force

The work done by the electric force when a charge q moves from point A to point B:

```
W = q × (V_A - V_B) = q × ΔV

where ΔV = V_A - V_B is the potential drop.
```

By the work-energy theorem:
- If W > 0: KE increases (charge gains speed)
- If W < 0: KE decreases (charge slows down)

### Worked Example 47.3

An electron (charge = -e = -1.6×10⁻¹⁹ C) accelerates through a potential difference of 500 V.

Find the kinetic energy gained.

**Solution:**

W = q × ΔV = (-1.6×10⁻¹⁹) × (-500) = **+8×10⁻¹⁷ J**

(The electron moves from low to high potential (toward the + plate), so ΔV is negative, but since q is negative, W is positive.)

KE gained = 8×10⁻¹⁷ J

Speed: v = sqrt(2KE/m) = sqrt(2 × 8×10⁻¹⁷ / 9.11×10⁻³¹) ≈ **1.33 × 10⁷ m/s**

(About 4% of the speed of light!)

---

## 47.8 Electron Volt

The **electron volt (eV)** is a convenient unit of energy in atomic and nuclear physics.

```
1 eV = energy gained by an electron accelerated through 1 volt

1 eV = e × 1V = 1.6×10⁻¹⁹ C × 1V = 1.6×10⁻¹⁹ J

Conversions:
  1 keV = 10³ eV = 1.6×10⁻¹⁶ J
  1 MeV = 10⁶ eV = 1.6×10⁻¹³ J
  1 GeV = 10⁹ eV = 1.6×10⁻¹⁰ J
```

Examples:
- Electron in a TV screen: ~15,000 eV = 15 keV
- X-ray photon: 1–100 keV
- Proton at Large Hadron Collider: 6.5 TeV = 6.5 × 10¹² eV

---

## 47.9 Capacitors and Stored Energy

A **capacitor** stores charge and electric potential energy. It consists of two conductors separated by an insulator.

### Capacitance

The capacitance C of a capacitor is:

```
C = Q / V

where:
  C = capacitance (Farads, F)
  Q = charge stored on each plate (C)
  V = potential difference between plates (V)
```

For a parallel-plate capacitor:

```
C = ε₀ × A / d

where:
  A = area of each plate (m²)
  d = separation (m)
  ε₀ = 8.85 × 10⁻¹² F/m
```

### Energy Stored in a Capacitor

```
E = (1/2) × C × V²
  = (1/2) × Q × V
  = Q² / (2C)
```

### Where is the Energy?

The energy is stored in the **electric field** between the plates:

```
Energy density = (1/2) × ε₀ × E²   (J/m³)
```

This is a profound result — empty space can store energy in the form of an electric field. This is the basis for understanding electromagnetic waves carrying energy through space.

---

## Summary

- **Electric potential energy**: PE = kqQ/r; positive for same-sign charges, negative for opposite
- **Electric potential**: V = PE/q = work per unit charge; scalar (V = J/C)
- **Point charge potential**: V = kQ/r; positive for + charge, negative for - charge
- **Superposition**: V_total = V₁ + V₂ + ... (scalar — easier than field!)
- **Potential difference (voltage)**: V = V_A - V_B; what voltmeters measure
- **Equipotential surfaces**: perpendicular to field lines; no work done moving along them
- **E and V**: E = -ΔV/Δx; field points from high to low potential
- **Work by electric force**: W = qΔV; if positive, KE increases
- **Electron volt**: 1 eV = 1.6×10⁻¹⁹ J; useful unit for atomic/nuclear energies
- **Capacitance**: C = Q/V; energy stored = ½CV²

---

## Key Equations

```
Electric potential:
  V = k × Q / r  (point charge)
  V = PE / q

Potential difference:
  ΔV = V_A - V_B = W/q

E-V relationship:
  E = -ΔV/Δx   (uniform field: E = V/d)

Work done by electric field:
  W = q × ΔV = q(V_A - V_B)
  W = ΔKE  (work-energy theorem)

Electron volt:
  1 eV = 1.6 × 10⁻¹⁹ J

Capacitance:
  C = Q / V
  C = ε₀A/d  (parallel plate)

Energy in capacitor:
  E = (1/2)CV² = (1/2)QV = Q²/(2C)

Energy density:
  u = (1/2)ε₀E²
```
