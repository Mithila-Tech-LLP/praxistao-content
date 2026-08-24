# Chapter 46: Electric Fields

> **"An electric field is a map of how space itself is altered by the presence of charges. Place a test charge anywhere, and the field tells you exactly what force it will feel."**

---

## Table of Contents

- [46.1 Why We Need the Field Concept](#461-why-we-need-the-field-concept)
- [46.2 Defining the Electric Field](#462-defining-the-electric-field)
- [46.3 Electric Field of a Point Charge](#463-electric-field-of-a-point-charge)
- [46.4 Electric Field Lines](#464-electric-field-lines)
- [46.5 Electric Field Due to Multiple Charges](#465-electric-field-due-to-multiple-charges)
- [46.6 Uniform Electric Field](#466-uniform-electric-field)
- [46.7 Motion of Charges in Electric Fields](#467-motion-of-charges-in-electric-fields)
- [46.8 Conductors in Electric Fields](#468-conductors-in-electric-fields)
- [46.9 Gauss's Law (Introduction)](#469-gausss-law-introduction)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 46.1 Why We Need the Field Concept

Coulomb's Law describes the force between two charges directly: F = kq₁q₂/r².

But there's a conceptual problem: how does charge A "know" that charge B is there? How does force propagate through empty space?

The **field concept** solves this. Instead of saying "charge A acts on charge B", we say:

1. **Charge A creates an electric field throughout space**
2. **That field acts on charge B wherever B happens to be**

The field is a real physical entity. It carries energy and can exist even with no charges around (as electromagnetic waves — light!).

```
DIRECT ACTION VIEW (problematic):
  
  [A] ----forces----> [B]  (how? what is the mechanism?)
  
  
FIELD VIEW (better):
  
  [A] ----creates---> [field E everywhere in space] ---acts on---> [B]
  
  The field is the intermediary. It's real and physical.
```

---

## 46.2 Defining the Electric Field

The electric field **E** at a point is defined as the force per unit charge experienced by a positive test charge placed at that point:

```
E = F / q₀

Or equivalently: F = q₀ × E

where:
  E  = electric field (N/C or V/m)
  F  = force on test charge (N)
  q₀ = small positive test charge (C) — "test" means it's small enough not to disturb the field being measured
```

The electric field is a **vector** — it has magnitude and direction.

Direction of E = direction of force on a **positive** test charge.

So:
- Near a positive charge: E points away (positive charge repels positive test charge)
- Near a negative charge: E points toward (negative charge attracts positive test charge)

---

## 46.3 Electric Field of a Point Charge

From Coulomb's Law, the force on test charge q₀ at distance r from point charge Q is:

```
F = kQq₀ / r²

So: E = F/q₀ = kQ / r²

     k × Q
E = -------
      r²
```

Direction: 
- Away from Q if Q is positive
- Toward Q if Q is negative

```
POSITIVE CHARGE Q:        NEGATIVE CHARGE Q:

       ↑                        ↑
    ↖  |  ↗                  ↙  |  ↘
  ←    Q    →              →    Q    ←
    ↙  |  ↘                  ↖  |  ↗
       ↓                        ↓
       
  Field points outward      Field points inward
  in all directions          from all directions
```

### Worked Example 46.1

Find the electric field at a point 0.5 m from a +2 μC charge.

**Solution:**

E = kQ/r² = 8.99×10⁹ × 2×10⁻⁶ / (0.5)² = 17,980 / 0.25 = **71,920 N/C ≈ 72 kN/C**

Direction: pointing away from the positive charge.

---

## 46.4 Electric Field Lines

**Electric field lines** are a visual tool to represent the electric field in a region of space.

### Rules for Drawing Field Lines

1. Field lines start on positive charges and end on negative charges (or go to infinity)
2. The direction of the field at any point is **tangent to the field line**
3. The **density** of field lines (how many per unit area) indicates the **strength** of the field
4. Field lines **never cross** (the field has one direction at each point)

```
ISOLATED + CHARGE:           + AND - PAIR (dipole):

      |                          +    -
    --|--                       /  \/ \
   |  |  |                     /    \/  \
  -|  +  |-                   /    / \   \
   |     |                   |    /   \   |
    --|--                   /    |     |   \
      |                    /      \   /    \
                                   \ /
                                    
Lines radiate outward.        Lines curve from + to -
                              (arrows point from + toward -)
```

### Information from Field Line Diagrams

```
WIDELY SPACED lines:     CLOSELY SPACED lines:

  |         |               ||||
  |         |               ||||
  
  Weak field               Strong field
  
PARALLEL EQUALLY-SPACED lines:

  |||||||||||||||
  |||||||||||||||
  |||||||||||||||
  
  Uniform field (same everywhere)
  Like the field between parallel plate capacitor
```

---

## 46.5 Electric Field Due to Multiple Charges

By the superposition principle, the total electric field at any point is the vector sum of fields from each individual charge:

```
E_total = E₁ + E₂ + E₃ + ...   (vector sum)
```

### Worked Example 46.2

Two charges: +3 μC at origin, -3 μC at x = 0.4 m. Find E at x = 0.2 m (midpoint).

```
     0.2 m    0.2 m
+3μC -------- P -------- -3μC
(x=0)      (x=0.2)      (x=0.4)
```

**Field from +3μC at P** (distance = 0.2 m):
E₁ = k × 3×10⁻⁶ / (0.2)² = 8.99×10⁹ × 3×10⁻⁶ / 0.04 = **674 kN/C pointing RIGHT** (away from + charge)

**Field from -3μC at P** (distance = 0.2 m):
E₂ = k × 3×10⁻⁶ / (0.2)² = **674 kN/C pointing RIGHT** (toward - charge)

**Total field = E₁ + E₂ = 674 + 674 = 1348 kN/C pointing RIGHT**

Both fields point in the same direction (right) at the midpoint of a dipole!

---

## 46.6 Uniform Electric Field

A **uniform electric field** has the same magnitude and direction everywhere in a region.

The most common source is two parallel plates carrying equal and opposite charges:

```
+ + + + + + + + + +   (positive plate)

→ → → → → → → → →   (uniform field between plates)
→ → → → → → → → →   (all arrows same direction and length)
→ → → → → → → → →

- - - - - - - - - -   (negative plate)
```

The field points from + plate to - plate.

The magnitude of the field between the plates:

```
E = V / d

where:
  V = potential difference (voltage) between plates (V)
  d = separation of plates (m)

Units: N/C = V/m  (these are equivalent)
```

This is one of the most useful configurations in physics — used in capacitors, cathode ray tubes, electron guns, and mass spectrometers.

---

## 46.7 Motion of Charges in Electric Fields

A charged particle in an electric field experiences a force F = qE.

By Newton's Second Law: a = F/m = qE/m

```
For positive charge:     acceleration is in the direction of E
For negative charge:     acceleration is opposite to E
```

### Uniform Field — Projectile-Like Motion

In a uniform electric field between two plates, the motion resembles projectile motion:

```
ELECTRON GUN:

 +plate         -plate
  |    →→→→→→   |
  |    →→→→→→   |
  |  e- →→→→→→ |-->  electron exits
  |    →→→→→→   |
  
Motion parallel to field: uniform acceleration (like gravity)
Motion perpendicular to field: constant velocity (if no other force)
```

This is exactly projectile motion with "gravity" replaced by the electric force.

### Worked Example 46.3

An electron (m = 9.11×10⁻³¹ kg, q = -1.6×10⁻¹⁹ C) is placed in a uniform E = 1000 N/C pointing right.

(a) Find the acceleration.
(b) Find the velocity after 10⁻⁸ s (starting from rest).

**Solution:**

(a) Force on electron: F = |q|E = 1.6×10⁻¹⁹ × 1000 = 1.6×10⁻¹⁶ N (pointing LEFT, since electron is negative)

    a = F/m = 1.6×10⁻¹⁶ / 9.11×10⁻³¹ = **1.76 × 10¹⁴ m/s²** (pointing left)

(b) v = at = 1.76×10¹⁴ × 10⁻⁸ = **1.76 × 10⁶ m/s** (to the left)

(About 0.6% of the speed of light — electrons are extremely light!)

---

## 46.8 Conductors in Electric Fields

When a conductor is placed in an external electric field, a very interesting thing happens:

```
CONDUCTOR IN EXTERNAL FIELD:

Before:                     After (equilibrium):

  →→→→[conductor]→→→         → → |+     -| → →
  →→→→           →→→         → → |+     -| → →
  →→→→           →→→         → → |+     -| → →
  
  External field              Electrons in conductor
                              redistribute. Surface has
                              induced charge. NET field
                              INSIDE = 0.
```

**Key properties of a conductor in electrostatics:**

1. **E = 0 inside a conductor** — free electrons redistribute until they completely cancel the external field inside

2. **All excess charge resides on the surface**

3. **E is perpendicular to the conductor surface** just outside (any parallel component would drive currents)

4. **Shielding**: a hollow conductor shields its interior from external electric fields (Faraday cage)

```
FARADAY CAGE:
  
  External field →→→→→
                     |+-+-|
  →→→→→ →→→→→       | 0  |   <-- interior shielded
                     |+-+-|
  →→→→→ →→→→→
  
  Used in: microwave ovens (prevent microwaves escaping),
           MRI rooms (shield from radio interference),
           cars (protect occupants from lightning)
```

---

## 46.9 Gauss's Law (Introduction)

**Gauss's Law** is a more elegant way to compute electric fields for symmetric charge distributions.

It states: the total electric flux through any closed surface equals the enclosed charge divided by ε₀.

```
Φ_E = Q_enc / ε₀

where:
  Φ_E = electric flux through the closed surface (N·m²/C)
  Q_enc = total charge enclosed by the surface (C)
  ε₀ = 8.85 × 10⁻¹² C²/(N·m²)
```

**Electric flux** = E × A × cos(θ), where A is surface area and θ is angle between E and the surface normal.

For a sphere around a point charge:

```
Φ = E × 4πr²  (sphere surface area)
   = Q/ε₀  (Gauss's Law)

Solving for E:
  E = Q/(4πε₀r²) = kQ/r²  ← recovers Coulomb's Law!
```

Gauss's Law is most useful for problems with **symmetry**: spherical, cylindrical, or planar distributions.

---

## Summary

- **Electric field E**: force per unit charge; E = F/q; vector quantity (N/C or V/m)
- Point charge: E = kQ/r²; direction away from + charge, toward - charge
- **Field lines**: tangent = field direction; density = field strength; start on +, end on -; never cross
- **Superposition**: E_total = vector sum of individual fields
- **Uniform field**: between parallel plates; E = V/d; field lines parallel, equally spaced
- **Force on charge**: F = qE; positive charge accelerates along E, negative charge opposite
- **Conductors in E field**: free electrons redistribute; E = 0 inside; charge on surface; Faraday cage
- **Gauss's Law**: flux through closed surface = Q_enc/ε₀; recovers Coulomb's Law for symmetric cases

---

## Key Equations

```
Electric field definition:
  E = F / q₀
  F = q × E

Point charge field:
  E = k × Q / r²
  k = 8.99 × 10⁹ N·m²/C²
  (or E = Q / (4πε₀r²))

Uniform field between plates:
  E = V / d

Force on charge in field:
  F = q × E
  a = qE / m

Superposition:
  E_total = E₁ + E₂ + ... (vectors)

Gauss's Law:
  Φ_E = E × A × cos(θ) = Q_enc / ε₀
  ε₀ = 8.85 × 10⁻¹² C²/(N·m²)
```
