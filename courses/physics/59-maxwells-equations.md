# Chapter 59: Maxwell's Equations

> **"In 1865, James Clerk Maxwell wrote four equations that unified electricity, magnetism, and light. Einstein called them the most profound and fruitful thing physics has produced. From these four equations, all of classical electromagnetism follows."**

---

## Table of Contents

- [59.1 The Story of Maxwell's Equations](#591-the-story-of-maxwells-equations)
- [59.2 Maxwell's Equations: Overview](#592-maxwells-equations-overview)
- [59.3 Gauss's Law for Electricity](#593-gausss-law-for-electricity)
- [59.4 Gauss's Law for Magnetism](#594-gausss-law-for-magnetism)
- [59.5 Faraday's Law](#595-faradays-law)
- [59.6 Ampère's Law with Maxwell's Correction](#596-ampères-law-with-maxwells-correction)
- [59.7 The Displacement Current](#597-the-displacement-current)
- [59.8 Electromagnetic Waves from Maxwell's Equations](#598-electromagnetic-waves-from-maxwells-equations)
- [59.9 The Speed of Light](#599-the-speed-of-light)
- [59.10 Implications and Significance](#5910-implications-and-significance)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 59.1 The Story of Maxwell's Equations

By 1860, four experimental laws had been established:

- **Coulomb/Gauss**: charged objects create electric fields
- **Magnetic Gauss**: no magnetic monopoles exist
- **Faraday**: changing magnetic fields create electric fields
- **Ampère**: electric currents create magnetic fields

James Clerk Maxwell (1831–1879) noticed an inconsistency in Ampère's law and added a correction term. This seemingly small mathematical fix led to the prediction that electromagnetic waves exist and travel at the speed of light.

When the speed came out as c = 3 × 10⁸ m/s — exactly the measured speed of light — Maxwell recognized that **light IS an electromagnetic wave**.

This was one of the greatest unifications in science.

---

## 59.2 Maxwell's Equations: Overview

In integral form (the most intuitive version):

```
MAXWELL'S FOUR EQUATIONS:

1. GAUSS'S LAW (electricity):
   ∮ E·dA = Q_enc/ε₀
   "Electric fields diverge from charges"

2. GAUSS'S LAW (magnetism):
   ∮ B·dA = 0
   "Magnetic fields have no sources or sinks (no monopoles)"

3. FARADAY'S LAW:
   ∮ E·dl = -dΦ_B/dt
   "Changing magnetic field creates electric field"

4. AMPÈRE-MAXWELL LAW:
   ∮ B·dl = μ₀I_enc + μ₀ε₀(dΦ_E/dt)
   "Electric current AND changing electric field create magnetic field"
```

Each equation tells a story. Together, they describe ALL electromagnetic phenomena.

---

## 59.3 Gauss's Law for Electricity

```
∮ E·dA = Q_enc/ε₀

In words: The total electric flux through any closed surface equals
          the net charge enclosed, divided by ε₀.
```

### What It Means

Electric field lines "diverge" (spread out) from positive charges and "converge" (flow in) toward negative charges. The number of field lines crossing any surface surrounding a charge depends only on the enclosed charge.

```
POSITIVE CHARGE +Q:

  Any closed surface
  surrounding +Q:
  
      .·|·.
    /   |   \
   | →→→+→→→ |   ← field lines crossing outward
    \   |   /
      '·|·'
  
  Net flux = Q/ε₀ (positive = outward)
  
  Surface with no charge: flux in = flux out → net = 0
```

### Using Gauss's Law

For symmetric charge distributions:

```
If E is uniform over a Gaussian surface of area A:
  E × A = Q_enc/ε₀
  E = Q_enc/(ε₀A)

For a sphere around point charge:
  E × 4πr² = Q/ε₀
  E = Q/(4πε₀r²) = kQ/r²  → recovers Coulomb's Law!
```

---

## 59.4 Gauss's Law for Magnetism

```
∮ B·dA = 0

In words: The total magnetic flux through any closed surface is zero.
```

### What It Means

Magnetic field lines always form closed loops — they never start or end. There are no magnetic monopoles (no isolated N or S poles).

```
COMPARE:

ELECTRIC:                        MAGNETIC:
  
  + → → → → →  flux exits       N
  (positive charge is source)     ↑
                                  |
  ← ← ← ← ← -  flux enters     ↓
  (negative charge is sink)      S
  
  Net flux ≠ 0 if charge inside   Field lines form closed loops
                                   Net flux through any surface = 0
```

This equation tells us magnets always come in dipoles (N and S together). If you cut a bar magnet in half, you get two smaller bar magnets, each with N and S poles.

---

## 59.5 Faraday's Law

```
∮ E·dl = -dΦ_B/dt

In words: A changing magnetic flux through a surface creates
          an EMF (circulation of E) around the boundary of that surface.
```

### What It Means

A time-varying magnetic field induces an electric field — even in empty space, with no charges or currents.

```
FARADAY'S LAW:

   Closed loop (wire):
      +--------+
      |        |
      |  B(t)  |  ← changing B inside
      |        |
      +--------+
           ↓
   EMF = -dΦ/dt   (Faraday)
   Current flows in wire (Lenz: opposing the change)
```

The minus sign encodes Lenz's Law: the induced field opposes the change.

### Connection to Electromagnetic Induction

This is exactly what we used in electromagnetic induction — a changing B creates an E that drives current. But Faraday's Law also applies in free space (no conductor needed):

A changing B field creates a circulating E field in space itself.

---

## 59.6 Ampère's Law with Maxwell's Correction

The original Ampère's Law:

```
∮ B·dl = μ₀I_enc
"A current creates a circulating magnetic field"
```

Maxwell found this was incomplete. Consider a capacitor charging:

```
CAPACITOR CHARGING:

  +plate ----I_wire----> [CAP] <----I_wire---- -plate
  
  Current flows through wires but NOT through the gap between plates.
  Yet B field circles around the whole circuit (including the gap?).
  
  Ampère's original law gives different B values depending on which
  surface you use to "catch" the current → INCONSISTENCY!
```

Maxwell's fix: add the **displacement current** term:

```
∮ B·dl = μ₀I_enc + μ₀ε₀(dΦ_E/dt)

The second term μ₀ε₀(dΦ_E/dt) is Maxwell's correction.
A changing electric field (like in the capacitor gap) creates a magnetic field
just as real current does.
```

---

## 59.7 The Displacement Current

The **displacement current density** is:

```
J_d = ε₀ × dE/dt

Displacement current: I_d = ε₀ × dΦ_E/dt
```

It's not a real current (no charges moving), but it produces a magnetic field just like a real current.

### In the Capacitor Gap

```
CHARGING CAPACITOR:

  +I_wire→    [  |  ]    ←I_wire-
  
  Real current in wires creates B field.
  
  In the gap: E increases as capacitor charges.
  dE/dt ≠ 0 → displacement current I_d = ε₀ dΦ_E/dt
  
  I_d = I_wire (they're equal!)
  
  So B field is continuous around the entire circuit,
  including the gap. Problem solved!
```

---

## 59.8 Electromagnetic Waves from Maxwell's Equations

Maxwell's great insight: combine equations 3 and 4.

**Faraday**: changing B creates E
**Ampère-Maxwell**: changing E creates B

So:
- Changing B creates E
- That changing E creates B
- That changing B creates E
- And so on...

This is a **self-sustaining propagating wave** — an electromagnetic wave!

```
ELECTROMAGNETIC WAVE:
  
  Direction of →→→→→→→→→→→→→→→→→→ propagation
  
                    ↑ E field (oscillating up/down)
       ___          |         ___
     /     \        |       /     \
    /       \       |      /       \
---+----+----+--x---+--x--+----+----+--- x
          ↔ B field (oscillating in/out of page)
  
  E and B are perpendicular to each other
  and perpendicular to the direction of travel.
  
  Both E and B oscillate in phase (peaks at same time and place).
```

### The Wave Equation

From Maxwell's equations, the wave equation for E:

```
∂²E/∂x² = μ₀ε₀ × ∂²E/∂t²

This is a wave equation with wave speed:

v = 1/√(μ₀ε₀)
```

---

## 59.9 The Speed of Light

Substituting the measured values:

```
μ₀ = 4π × 10⁻⁷ T·m/A
ε₀ = 8.85 × 10⁻¹² C²/(N·m²)

c = 1/√(μ₀ε₀) = 1/√(4π×10⁻⁷ × 8.85×10⁻¹²)

c = 1/√(11.12×10⁻¹⁸)

c = 1/(3.33×10⁻⁹) = 2.998 × 10⁸ m/s ≈ 3 × 10⁸ m/s
```

This is exactly the measured speed of light.

Maxwell's conclusion: **light is an electromagnetic wave.**

```
"We can scarcely avoid the conclusion that light consists in
the transverse undulations of the same medium which is the
cause of electric and magnetic phenomena."
— James Clerk Maxwell, 1865
```

This was one of the most stunning predictions in the history of physics, later confirmed by Heinrich Hertz (1887) who produced radio waves in the lab using sparks.

---

## 59.10 Implications and Significance

### Everything Light Does Is Explained by Maxwell's Equations

- Reflection, refraction, diffraction, interference — all derivable from Maxwell's equations
- The refractive index of a material: n = √(ε_r × μ_r) — from the material's electric and magnetic properties
- Speed of light in a medium: v = c/n

### Relativity

Einstein noticed that Maxwell's equations predict a specific speed c for light, independent of any reference frame. This was incompatible with Newtonian mechanics.

The resolution was **special relativity** (1905): the laws of physics are the same in all inertial frames, and c is the same in all frames. This overthrew Newtonian mechanics and completely changed our understanding of space and time.

Maxwell's equations are already relativistically correct — Newton's laws needed modification.

### Wireless Communication

Maxwell's equations predict radio waves, microwaves, infrared, visible light, X-rays — all the same thing, just different frequencies:

```
Frequency and wavelength of EM waves:
  c = f × λ
  
  f = 10² Hz:  Radio (λ = 3×10⁶ m)
  f = 10⁹ Hz:  Microwave (λ = 0.3 m)
  f = 10¹² Hz: Infrared (λ = 3×10⁻⁴ m)
  f = 5×10¹⁴ Hz: Visible light (λ = 600 nm)
  f = 10¹⁸ Hz: X-rays (λ = 3×10⁻¹⁰ m)
```

Hertz produced radio waves in 1887 → Marconi sent them across the Atlantic in 1901 → every phone, WiFi, GPS, TV broadcast today.

---

## Summary

- Maxwell's four equations are the complete description of classical electromagnetism
- **Gauss's Law (E)**: electric flux ∝ enclosed charge; E field diverges from charges
- **Gauss's Law (B)**: no magnetic monopoles; B field lines always close
- **Faraday's Law**: changing B creates E; basis of electromagnetic induction
- **Ampère-Maxwell Law**: current AND changing E both create B
- Maxwell's addition: **displacement current** — changing electric field creates magnetic field
- Together: self-sustaining EM wave propagates with speed c = 1/√(μ₀ε₀)
- Light IS an electromagnetic wave (Maxwell 1865)
- Prediction confirmed by Hertz (1887); led to radio, TV, WiFi, microwave ovens, etc.
- Maxwell's equations are already relativistically correct; led Einstein to develop special relativity

---

## Key Equations

```
Maxwell's Equations (integral form):

1. Gauss (electricity): ∮ E·dA = Q_enc/ε₀

2. Gauss (magnetism): ∮ B·dA = 0

3. Faraday: ∮ E·dl = -dΦ_B/dt

4. Ampère-Maxwell: ∮ B·dl = μ₀I_enc + μ₀ε₀(dΦ_E/dt)

Speed of light:
  c = 1/√(μ₀ε₀) ≈ 3 × 10⁸ m/s

Electromagnetic wave:
  c = f × λ
  E and B perpendicular to each other and to direction of propagation
  Both E and B oscillate in phase

Constants:
  μ₀ = 4π × 10⁻⁷ T·m/A
  ε₀ = 8.85 × 10⁻¹² C²/(N·m²)
  c = 2.998 × 10⁸ m/s
```
