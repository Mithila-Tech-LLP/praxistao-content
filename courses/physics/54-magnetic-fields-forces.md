# Chapter 54: Magnetic Fields and Forces

> **"Electricity and magnetism are two faces of the same phenomenon. Every moving charge creates a magnetic field. Every magnetic field exerts a force on moving charges. Understanding this connection changed civilization."**

---

## Table of Contents

- [54.1 What are Magnetic Fields?](#541-what-are-magnetic-fields)
- [54.2 Magnetic Field Lines](#542-magnetic-field-lines)
- [54.3 Force on a Moving Charge](#543-force-on-a-moving-charge)
- [54.4 Force on a Current-Carrying Wire](#544-force-on-a-current-carrying-wire)
- [54.5 Torque on a Current Loop](#545-torque-on-a-current-loop)
- [54.6 Magnetic Fields Created by Currents](#546-magnetic-fields-created-by-currents)
- [54.7 The Solenoid](#547-the-solenoid)
- [54.8 The Earth's Magnetic Field](#548-the-earths-magnetic-field)
- [54.9 Applications](#549-applications)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 54.1 What are Magnetic Fields?

A **magnetic field** is a region of space where magnetic forces act.

Magnetic fields are created by:
1. **Permanent magnets** (due to aligned magnetic moments of electrons)
2. **Moving charges / electric currents** (any current creates a magnetic field)
3. **Changing electric fields** (Maxwell's discovery — leads to electromagnetic waves)

The magnetic field is denoted **B** (sometimes called **H** in different contexts) and measured in **Tesla (T)** or **Gauss (G)** (1 T = 10,000 G).

### Typical Magnetic Field Strengths

| Source | B (T) |
|--------|-------|
| Earth's magnetic field | ~5 × 10⁻⁵ T = 0.5 G |
| Bar magnet | ~10⁻² to 10⁻¹ T |
| MRI machine | 1.5 to 7 T |
| Strongest lab magnet | ~45 T |
| Neutron star | ~10⁸ T |

---

## 54.2 Magnetic Field Lines

Like electric field lines, **magnetic field lines** show the direction and strength of the magnetic field.

```
BAR MAGNET FIELD LINES:

      N                 S
      |                 |
   ←  |  →          ←   |  →
  ←   |   →        ←    |    →
 ←    N    →      ←     S     →
  ←   |   →        ←    |    →
   ←  |  →          ←   |  →
      |       loop       |

Field lines:
- Exit from North pole
- Enter South pole
- Form closed loops (NEVER start or end — unlike electric field lines!)
- Closer together = stronger field
- Direction = direction N pole of small compass would point
```

**Key difference from electric fields**: Magnetic field lines always close on themselves (no magnetic monopoles — no isolated N or S poles).

### The Right-Hand Rule for Magnets

```
COMPASS IN A MAGNETIC FIELD:

  N pole of compass → direction of B field

Inside a bar magnet, field goes from S to N (opposite to outside).
```

---

## 54.3 Force on a Moving Charge

A charge moving in a magnetic field experiences a force — the **Lorentz force**:

```
F = q × v × B × sin(θ)

Or in vector form: F = q(v × B)

where:
  F = force (N)
  q = charge (C)
  v = velocity (m/s)
  B = magnetic field strength (T)
  θ = angle between v and B
```

### Key Properties of the Magnetic Force

```
1. F is perpendicular to BOTH v and B
   (always perpendicular to motion → does no work → doesn't change speed)

2. F = 0 if the charge is stationary (v = 0)

3. F = 0 if v is parallel or anti-parallel to B (θ = 0 or 180°)

4. F is maximum when v ⊥ B (θ = 90°)
```

### Right-Hand Rule for F = qv × B

For positive charge:
1. Point fingers in direction of velocity v
2. Curl them toward direction of B
3. Thumb points in direction of force F

For negative charge: force is opposite.

```
VISUAL:
  
  B pointing into page: ⊗ ⊗ ⊗
                        ⊗ ⊗ ⊗
  
  Positive charge moving right: --→ +q
  
  F = qvB sin(90°) = qvB   (maximum)
  Direction (right-hand rule): UPWARD
  
  ⊗ ⊗ ⊗
  ⊗ →+q ⊗   F ↑
  ⊗ ⊗ ⊗
```

### Circular Motion in a Magnetic Field

Since F is always perpendicular to v, a charge moving in a uniform magnetic field moves in a **circle**:

```
CIRCULAR MOTION:
  
  B pointing out of page: ⊙ ⊙ ⊙
                          ⊙ ⊙ ⊙
  
  Positive charge moving right:
  
    ⊙ ⊙ ⊙
    ⊙ ↑  ⊙
    ⊙ →+q ⊙
    ⊙ ⊙ ⊙
  
  Force always perpendicular to velocity → circular orbit!
```

The radius of circular motion:

```
Centripetal force = magnetic force
mv²/r = qvB

r = mv/(qB)

where:
  r = radius (m)
  m = mass of particle (kg)
  v = speed (m/s)
  q = charge (C)
  B = field strength (T)
```

This is the basis of mass spectrometers: different masses curve to different radii.

### Worked Example 54.1

An electron (m = 9.11×10⁻³¹ kg, q = 1.6×10⁻¹⁹ C) moves at 2×10⁶ m/s perpendicular to a 0.1 T magnetic field.

Find the radius of its circular orbit.

**Solution:**

r = mv/(qB) = (9.11×10⁻³¹ × 2×10⁶) / (1.6×10⁻¹⁹ × 0.1)

r = 1.822×10⁻²⁴ / 1.6×10⁻²⁰ = **1.14 × 10⁻⁴ m = 0.114 mm**

---

## 54.4 Force on a Current-Carrying Wire

A current is just moving charges. So a current-carrying wire in a magnetic field experiences a force:

```
F = B × I × L × sin(θ)

where:
  F = force on wire (N)
  B = magnetic field strength (T)
  I = current (A)
  L = length of wire in field (m)
  θ = angle between wire and B
```

Direction: right-hand rule, with v replaced by current direction.

```
FORCE ON WIRE IN MAGNETIC FIELD:

  B pointing into page: ⊗ ⊗ ⊗ ⊗ ⊗
                        ⊗ ⊗ ⊗ ⊗ ⊗
  
  Current flowing right in wire:
  
     ⊗ ⊗ ⊗ ⊗ ⊗
     ⊗ I→ I→ I→ ⊗    Force (↑) on wire
     ⊗ ⊗ ⊗ ⊗ ⊗
     
  Force is UPWARD (right-hand rule: fingers right, curl into page → thumb up)
```

### Worked Example 54.2

A 0.5 m long wire carrying 3 A is placed perpendicular to a 0.2 T magnetic field.

F = BIL sin(90°) = 0.2 × 3 × 0.5 × 1 = **0.3 N**

---

## 54.5 Torque on a Current Loop

A rectangular loop of wire carrying current, placed in a magnetic field, experiences a **torque** that tends to rotate it:

```
CURRENT LOOP IN FIELD:

     B→→→→→→→→
     
     +------+
  F↑ |      | F↓   (two sides experience equal opposite forces → torque)
     |  I   |
     +------+
     
  (Top and bottom sides: force parallel to plane — no torque contribution)
```

Torque:
```
τ = NIAB sin(θ)

where:
  N = number of turns
  I = current
  A = area of loop
  B = magnetic field
  θ = angle between B and normal to loop
```

This is how **electric motors** work! The rotating loop is the armature.

---

## 54.6 Magnetic Fields Created by Currents

Electric currents create magnetic fields. The direction is given by the **right-hand rule for wires**:

```
LONG STRAIGHT WIRE:
  
  Point right thumb in direction of current.
  Fingers curl around wire in direction of B field.
  
  Current going right:
  
  →→→→→→→→→→→→→→  wire
  
  B circles around the wire:
  - Above wire: B pointing OUT of page
  - Below wire: B pointing INTO page
```

Field strength for a long straight wire:

```
B = μ₀ × I / (2π × r)

where:
  μ₀ = 4π × 10⁻⁷ T·m/A  (permeability of free space)
  I = current (A)
  r = distance from wire (m)
```

### Worked Example 54.3

Find B at 0.1 m from a wire carrying 10 A.

B = (4π×10⁻⁷ × 10) / (2π × 0.1) = (4π×10⁻⁶) / (0.2π) = 2×10⁻⁵ T = **20 μT**

---

## 54.7 The Solenoid

A **solenoid** is a coil of wire with many turns. Inside the solenoid, the magnetic fields from each turn reinforce each other to create a strong, uniform field.

```
SOLENOID (cross-section):

  →→→→→→→→→→→→→→→→→→→→→  (field inside, uniform and strong)
  
  ×                        ×  (current flowing into page on top)
  × × × × × × × × × × × ×
  ○ ○ ○ ○ ○ ○ ○ ○ ○ ○ ○ ○  (current coming out of page on bottom)
  ○                        ○
  
  →→→→→→→→→→→→→→→→→→→→→  (field lines loop around outside)
```

Magnetic field inside a solenoid:

```
B = μ₀ × n × I

where:
  n = number of turns per unit length (turns/m)
  I = current (A)
  B is uniform throughout the interior
```

With an iron core, B is multiplied by the relative permeability μ_r (up to ~10,000 for iron) — making an **electromagnet**.

---

## 54.8 The Earth's Magnetic Field

The Earth acts like a giant bar magnet (actually due to convection currents of molten iron in the outer core).

```
EARTH'S MAGNETIC FIELD:

  Geographic North Pole
        |
        |  ← geographic north
        |
    S (magnetic)   (Earth's magnetic south is near geographic north!)
    
  Field lines:
    enter Earth near geographic north
    exit Earth near geographic south
    
  A compass needle's N pole points toward geographic north
  → the magnetic pole near geographic north is actually a SOUTH magnetic pole!
```

Key facts:
- Magnetic poles are not aligned with geographic poles (~11° tilt)
- The poles wander slowly over time
- The field has reversed many times in geological history

### Applications

- **Navigation**: compasses
- **Animal navigation**: birds, fish, bacteria have magnetite crystals
- **Cosmic ray deflection**: the field protects Earth from harmful solar particles

---

## 54.9 Applications

### Electric Motor

Uses torque on a current loop. Current direction reversed every half turn (by commutator) to maintain rotation in one direction.

```
MOTOR PRINCIPLE:
  Current loop + magnetic field → torque → rotation
  
  Commutator reverses current direction every 180° → continuous rotation
```

### Loudspeaker

A coil of wire attached to a cone sits in a magnetic field. AC current through the coil creates alternating forces → cone vibrates → produces sound.

### Mass Spectrometer

Ions accelerated through voltage, then curved by magnetic field. Radius of curvature r = mv/(qB) depends on mass-to-charge ratio. Different masses separate out → identify composition.

### Magnetic Levitation (Maglev Trains)

Superconducting electromagnets repel track magnets → train floats → no friction → very fast (600 km/h+).

### MRI (Magnetic Resonance Imaging)

Strong magnetic field (1.5–7 T) aligns hydrogen nuclei. Radio wave pulse at resonant frequency tips them. As they relax, they emit radio waves → detected → create image.

---

## Summary

- **Magnetic field B**: created by moving charges/currents; measured in Tesla (T)
- **Field lines**: form closed loops (no monopoles); exit N pole, enter S pole
- **Force on moving charge**: F = qvB sin(θ); perpendicular to both v and B
- **Right-hand rule**: determines direction of F, or direction of B around wire
- Magnetic force does **no work** (always perpendicular to motion)
- **Circular motion**: r = mv/(qB); basis of mass spectrometers
- **Force on wire**: F = BIL sin(θ); basis of electric motors
- **Torque on loop**: τ = NIAB sin(θ); how motors rotate
- **Long straight wire**: B = μ₀I/(2πr)
- **Solenoid**: B = μ₀nI; uniform field inside
- Applications: motors, speakers, mass spectrometers, MRI, maglev trains

---

## Key Equations

```
Force on moving charge:
  F = qvB sin(θ)   (magnitude)
  F = q(v × B)     (vector form)

Circular orbit:
  r = mv / (qB)

Force on current wire:
  F = BIL sin(θ)

Torque on current loop:
  τ = NIAB sin(θ)

Long straight wire:
  B = μ₀I / (2πr)

Solenoid:
  B = μ₀nI  (n = turns/meter)

Constants:
  μ₀ = 4π × 10⁻⁷ T·m/A
  1 Tesla = 1 Wb/m² = 1 kg/(A·s²)
```
