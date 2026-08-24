# Chapter 04: Crystal Structures — How Atoms Pack

> **"A metal is not a random jumble of atoms. It is an exquisitely ordered 3D array, repeated billions of times. That perfect order is why metals have predictable, reproducible properties — and why their 'mistakes' (defects) matter so much."**

---

## Table of Contents

1. [What Is a Crystal?](#1-what-is-a-crystal)
2. [The Unit Cell Concept](#2-the-unit-cell-concept)
3. [The Three Common Metal Crystal Structures](#3-the-three-common-metal-crystal-structures)
4. [Body-Centered Cubic (BCC)](#4-body-centered-cubic-bcc)
5. [Face-Centered Cubic (FCC)](#5-face-centered-cubic-fcc)
6. [Hexagonal Close-Packed (HCP)](#6-hexagonal-close-packed-hcp)
7. [Atomic Packing Factor and Density](#7-atomic-packing-factor-and-density)
8. [Coordination Number — How Many Neighbors?](#8-coordination-number--how-many-neighbors)
9. [Crystal Directions and Planes (Miller Indices)](#9-crystal-directions-and-planes-miller-indices)
10. [Polymorphism — One Metal, Multiple Structures](#10-polymorphism--one-metal-multiple-structures)
11. [Allotropes of Iron — Why Steel Heat Treatment Works](#11-allotropes-of-iron--why-steel-heat-treatment-works)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What Is a Crystal?

A **crystal** is a solid in which atoms are arranged in a **regular, repeating, three-dimensional pattern** that extends throughout the material.

Contrast with:
- **Amorphous solid** (glass, some polymers): atoms are randomly arranged, no long-range order
- **Liquid**: atoms are close but randomly arranged and mobile
- **Gas**: atoms far apart, random positions

Virtually all metals are **crystalline** under normal conditions. When you melt a metal and let it solidify, the atoms spontaneously arrange into a crystal — it's thermodynamically favorable.

### Why Crystalline Structure Matters

The crystal structure determines:
- **Density** — how efficiently atoms pack
- **Ductility** — which directions atoms can slide past each other (slip systems)
- **Magnetic properties** — the arrangement determines spin alignment
- **Elastic anisotropy** — stiffness varies with direction
- **Phase transformations** — steel hardening only works because iron changes crystal structure with temperature

---

## 2. The Unit Cell Concept

The entire crystal can be described by its **unit cell** — the smallest repeating unit that contains all the symmetry of the crystal:

```
        Imagine a crystal as a 3D wallpaper pattern.
        The unit cell is one repeat of the motif.
        
        ┌──┬──┬──┬──┬──┐
        │  │  │  │  │  │    ← 2D example: each square is a unit cell
        ├──┼──┼──┼──┼──┤
        │  │  │  │  │  │
        └──┴──┴──┴──┴──┘
```

In 3D, the unit cell is a parallelepiped defined by three edge lengths (a, b, c) and three angles (α, β, γ).

There are **7 crystal systems** and **14 Bravais lattices**. For metals, three structures dominate:
- **BCC** (body-centered cubic): a = b = c, all angles 90°, one extra atom at center
- **FCC** (face-centered cubic): a = b = c, all angles 90°, extra atoms at face centers
- **HCP** (hexagonal close-packed): a = b ≠ c, angles 90°, 90°, 120°

---

## 3. The Three Common Metal Crystal Structures

Roughly 95% of all engineering metals have one of these three structures at room temperature:

| Structure | Abbreviation | Examples | APF | Coord. No. |
|-----------|-------------|---------|-----|-----------|
| Body-centered cubic | BCC | Fe (α), Cr, Mo, W, V, Ta, Nb | 0.68 | 8 |
| Face-centered cubic | FCC | Fe (γ), Al, Cu, Ni, Au, Ag, Pb, Co (β) | 0.74 | 12 |
| Hexagonal close-packed | HCP | Mg, Ti (α), Zn, Co (α), Zr, Be | 0.74 | 12 |

**The key number: APF (Atomic Packing Factor)** = fraction of the unit cell volume actually occupied by atom spheres. Higher APF = more efficiently packed.

---

## 4. Body-Centered Cubic (BCC)

```
        BCC Unit Cell:
        
             •─────────•
            /|         /|
           / |        / |
          •─────────•  |
          |  |      |  |
          |  •───── | ─•
          | /       | /
          |/        |/
          •─────────•
          
          Plus one atom right in the CENTER
```

**Atoms per unit cell:**
- 8 corner atoms × (1/8 each) = 1
- 1 body-center atom × 1 = 1
- **Total: 2 atoms per unit cell**

**Lattice parameter:** The atoms touch along the **body diagonal** (corner-to-center-to-corner = 4r):
```
Body diagonal = a√3 = 4r
Therefore: a = 4r/√3
```

**Key BCC metals:**
- **α-Iron (Fe)**: stable below 912°C — the basis of carbon steel
- **Chromium (Cr)**: BCC, gives corrosion resistance to stainless steel
- **Molybdenum (Mo)**: BCC refractory metal, solid-solution strengthener in steels and superalloys
- **Tungsten (W)**: BCC, highest melting metal (3422°C)
- **Vanadium (V)**: BCC, microalloying addition in high-strength steels
- **Tantalum (Ta)** and **Niobium (Nb)**: BCC refractories, key superalloy additions

**Characteristic properties of BCC:**
- Less densely packed (APF = 0.68) → more room for interstitial atoms (C, N)
- Fewer slip systems (48 possible, but only {110}⟨111⟩ and {112}⟨111⟩ active) → harder
- Exhibits **ductile-to-brittle transition** (DBTT) at low temperatures — important for structural steels in cold climates
- Often harder and stronger than FCC at room temperature

---

## 5. Face-Centered Cubic (FCC)

```
        FCC Unit Cell:
        
             •─────────•
            /• •       /•
           / |  •     / |
          •─────────•  |
          •| |    • •| |
          | •─────── | ─•
          |/ •       |/ •
          •─────────•
          
          Corner atoms + face-center atoms (one per face, 6 faces)
```

**Atoms per unit cell:**
- 8 corner atoms × (1/8) = 1
- 6 face-center atoms × (1/2) = 3
- **Total: 4 atoms per unit cell**

**Lattice parameter:** Atoms touch along the **face diagonal**:
```
Face diagonal = a√2 = 4r
Therefore: a = 4r/√2 = 2r√2
```

**Key FCC metals:**
- **γ-Iron (Austenite)**: iron above 912°C — critical for steel heat treatment
- **Nickel (Ni)**: FCC always — the base of all nickel superalloys
- **Aluminum (Al)**: FCC — lightweight structural metal
- **Copper (Cu)**: FCC — best common conductor
- **Gold (Au) and Silver (Ag)**: FCC — soft precious metals
- **Lead (Pb)**: FCC — very soft, used for vibration damping
- **β-Cobalt**: FCC above 421°C

**Characteristic properties of FCC:**
- More densely packed (APF = 0.74) → less room for interstitials
- **12 slip systems** ({111}⟨110⟩: 4 planes × 3 directions) → highly ductile
- No DBTT — remains ductile at cryogenic temperatures (important for LNG tanks, space applications)
- Austenitic stainless steels are FCC → good cryogenic toughness

**Why FCC is more ductile than BCC:**
FCC has 12 easy slip systems; BCC has fewer active ones at low stress. More slip systems means more directions for dislocations to move → greater ductility.

This is why austenitic stainless steel (FCC) is far more ductile than martensitic steel (BCT, similar to BCC).

---

## 6. Hexagonal Close-Packed (HCP)

HCP is also close-packed (APF = 0.74, same as FCC), but the stacking sequence of layers is different:

**FCC stacking:** ABCABC...  
**HCP stacking:** ABABAB...

Think of stacking oranges: in the third layer, you can place each orange either directly above the first layer (→ ABAB = HCP) or offset from both layers 1 and 2 (→ ABCABC = FCC).

```
        HCP Unit Cell:
        
        Top layer (A):    • • •
                         • • •   (hexagon top)
        
        Middle layer (B): • •     (3 atoms in hollows)
        
        Bottom layer (A): • • •
                         • • •   (hexagon bottom)
```

**Atoms per unit cell:**
- 12 corner atoms × (1/6) = 2
- 2 face atoms × (1/2) = 1
- 3 internal atoms × 1 = 3
- **Total: 6 atoms per unit cell**

**c/a ratio:** Ideal close-packed sphere packing gives c/a = 1.633. Real metals deviate slightly:
| Metal | c/a ratio |
|-------|-----------|
| Ideal | 1.633 |
| Mg | 1.623 (close to ideal) |
| Zn | 1.856 (significantly non-ideal) |
| Ti | 1.588 (slightly low) |
| Be | 1.568 (low) |

Non-ideal c/a affects mechanical behavior significantly.

**Key HCP metals:**
- **α-Titanium**: HCP below 883°C — the base of most titanium alloys used in aerospace
- **Magnesium (Mg)**: lightest structural metal, HCP
- **Zinc (Zn)**: protective coatings (galvanizing), HCP
- **α-Cobalt**: HCP below 421°C
- **Zirconium (Zr)**: nuclear fuel cladding, HCP

**Characteristic properties of HCP:**
- Only **3 easy slip systems** at room temperature (basal plane {0001}⟨1120⟩)
- Much less ductile than FCC; more limited formability
- Can sometimes activate prismatic or pyramidal slip at high temperature → more ductile at elevated temperature
- Mg and Zn are notoriously hard to form at room temperature — you must warm them up

---

## 7. Atomic Packing Factor and Density

**Atomic Packing Factor (APF):**
```
APF = (number of atoms per unit cell × volume of one atom) / volume of unit cell

For FCC (a = 2r√2):
APF = 4 × (4/3)πr³ / (2r√2)³ = 4 × (4/3)πr³ / (8√2 r³)
APF = 4π / (3 × 8√2) = 16π / (3 × 8√2) = π/(3√2) ≈ 0.7405
```

For BCC: APF = π√3/8 ≈ 0.6802

**Calculating theoretical density:**
```
ρ = (n × A) / (V_c × N_A)

where:
  n   = atoms per unit cell
  A   = atomic mass (g/mol)
  V_c = volume of unit cell (cm³)
  N_A = Avogadro's number (6.022 × 10²³ atoms/mol)
```

**Example — Copper (FCC, a = 0.3615 nm, A = 63.55 g/mol):**
```
n = 4 (FCC)
V_c = (3.615 × 10⁻⁸ cm)³ = 4.724 × 10⁻²³ cm³
ρ = (4 × 63.55) / (4.724 × 10⁻²³ × 6.022 × 10²³)
ρ = 254.2 / 28.44
ρ = 8.94 g/cm³
```
Measured value: 8.96 g/cm³ — excellent agreement!

---

## 8. Coordination Number — How Many Neighbors?

The **coordination number** is the number of nearest-neighbor atoms surrounding any given atom:

| Structure | Coordination Number |
|-----------|-------------------|
| BCC | **8** (8 atoms at cube corners touch body-center) |
| FCC | **12** (6 face atoms + 6 face atoms of adjacent cells) |
| HCP | **12** (6 in same layer + 3 above + 3 below) |

Higher coordination number = more neighbors = more bonding = generally higher melting point and density.

FCC and HCP both have CN=12 — they are both "close-packed" in the mathematical sense. This is also why they have the same APF (0.74).

---

## 9. Crystal Directions and Planes (Miller Indices)

To describe directions and planes within a crystal, we use **Miller indices** — essential for specifying slip systems, cleavage planes, growth directions, and X-ray diffraction.

### Crystal Directions: [uvw]

Use integers along the a, b, c axes of the unit cell. Steps to find [uvw]:
1. Set origin at the tail of the vector
2. Read the components along each axis
3. Convert to smallest integers (multiply out fractions)
4. Enclose in square brackets: **[uvw]**

Common directions in a cubic crystal:
- Along cube edge: **[100]**, [010], [001] — and their negatives [1̄00], [01̄0], [001̄]
- Along face diagonal: **[110]**, [101], [011]...
- Along body diagonal: **[111]**, [1̄11], [11̄1]...

A family of equivalent directions uses angle brackets: ⟨100⟩ means [100], [010], [001] and their negatives.

### Crystal Planes: (hkl)

The Miller index for a plane is the reciprocal of its intercepts on the three axes:

1. Find where the plane cuts the a, b, c axes (in units of lattice parameter)
2. Take reciprocals
3. Clear fractions, convert to integers
4. Enclose in parentheses: **(hkl)**

Common planes in a cubic crystal:
- Cube face: **(100)** — the face of the cube; intercepts at 1, ∞, ∞
- Diagonal plane: **(110)** — intercepts at 1, 1, ∞
- Close-packed plane in FCC: **(111)** — the plane cutting all three axes at 1

A family of equivalent planes uses curly braces: {111} means all four equivalent diagonal planes.

### Slip Systems = Slip Plane + Slip Direction

A **slip system** is a combination of the plane on which dislocations move and the direction in which they move. The most important in metals:

| Structure | Slip Plane | Slip Direction | # Systems |
|-----------|-----------|----------------|-----------|
| FCC | {111} | ⟨110⟩ | 4 × 3 = **12** |
| BCC | {110} | ⟨111⟩ | 6 × 2 = 12 (but fewer active) |
| HCP | {0001} (basal) | ⟨1120⟩ | 1 × 3 = **3** (limited!) |

FCC has 12 easy slip systems → most ductile. HCP has only 3 → most limited ductility. This explains everything from why aluminum foil can be pressed into any shape, to why Mg alloys must be warm-formed, to why Ti alloys are intermediate.

**Single crystal [001] orientation for turbine blades:**
The [001] direction in a single-crystal Ni superalloy has the lowest elastic modulus (~125 GPa vs ~220 GPa in [111]). This means the blade is most compliant in the centrifugal direction → thermal stresses are lower (less stiff = less stress for the same thermal strain). This is the mechanical reason why [001] single crystals outperform polycrystals. Chapter 51 covers this in detail.

---

## 10. Polymorphism — One Metal, Multiple Structures

**Polymorphism** (also called **allotropy** for pure elements): the same element existing in different crystal structures depending on temperature or pressure.

| Metal | Low-T form | Transition T | High-T form | Transition T | Higher-T form |
|-------|-----------|-------------|------------|-------------|--------------|
| **Iron (Fe)** | α-Fe (BCC) | 912°C | γ-Fe (FCC) | 1394°C | δ-Fe (BCC) |
| **Titanium (Ti)** | α-Ti (HCP) | 883°C | β-Ti (BCC) | — | — |
| **Cobalt (Co)** | α-Co (HCP) | 421°C | β-Co (FCC) | — | — |
| **Tin (Sn)** | α-Sn (diamond cubic) | 13.2°C | β-Sn (tetragonal) | 231°C | liquid |
| **Zirconium (Zr)** | α-Zr (HCP) | 863°C | β-Zr (BCC) | — | — |

Polymorphism is enormously important for processing because you can:
1. Heat the metal into one crystal structure (different properties, solubility)
2. Process it
3. Cool it back — sometimes the transformation can be controlled to produce microstructures not possible in the simpler single-phase material

---

## 11. Allotropes of Iron — Why Steel Heat Treatment Works

Iron is the most important example of polymorphism in engineering:

```
Temperature (°C):
1538 ─────────── Melting point (liquid)
     ─────────── L → δ-Fe (BCC) solidifies
1394 ─────────── δ-Fe (BCC) → γ-Fe (FCC) = AUSTENITE
                 ← can dissolve 2.14 wt% C in FCC (much more room)
912  ─────────── γ-Fe (FCC) → α-Fe (BCC) = FERRITE  
                 ← can only dissolve 0.022 wt% C in BCC
770  ─────────── Curie temperature: α-Fe becomes non-magnetic above this
  0  ─────────── Room temperature: α-Fe (BCC)
```

**Why this matters for steel:**

When you heat steel to 900°C (austenite, FCC), carbon dissolves into the structure — up to 2.14%. When you quench (cool rapidly to room temperature), the carbon has no time to escape. The iron wants to transform from FCC back to BCC, but BCC can't hold all that carbon. The result is a distorted, body-centered **tetragonal** (BCT) structure: **martensite** — the hardest, most brittle form of steel.

This is why quenching hardens steel. The entire art of steel heat treatment exploits iron's crystal structure transformations. Without allotropy of iron, there would be no hardened steel tools, no springs, no ball bearings.

Titanium's α→β transformation is equally important for titanium alloy processing — we'll cover it in Chapter 25.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Crystal | Atoms in regular, repeating 3D array |
| Unit cell | Smallest repeating unit; described by a, b, c, α, β, γ |
| BCC | 2 atoms/cell, APF=0.68, CN=8; Fe(α), Cr, W, Mo, V |
| FCC | 4 atoms/cell, APF=0.74, CN=12; Ni, Al, Cu, Fe(γ); 12 slip systems → ductile |
| HCP | 6 atoms/cell, APF=0.74, CN=12; Mg, Ti(α), Zn; only 3 slip systems |
| APF | Fraction of space occupied; FCC=HCP=0.74 > BCC=0.68 |
| Miller indices | [uvw] for directions, (hkl) for planes; {hkl} family, ⟨uvw⟩ family |
| Slip system | Plane + direction for dislocation motion; more = more ductile |
| Polymorphism | Same element, different crystal structures at different T; Fe: α(BCC)→γ(FCC)→δ(BCC) |
| Iron allotropy | γ-Fe (FCC, austenite) holds much more C than α-Fe (BCC, ferrite) → basis of steel hardening |

**Next chapter:** Perfect crystals are useful for understanding ideal properties. But real metals are full of imperfections — vacancies, dislocations, grain boundaries. These "defects" are not failures; they are what makes metallurgy possible. The controlled management of defects is the entire art of alloy design.

---

## Exercises

1. Iron at room temperature is BCC with lattice parameter a = 0.2866 nm. Calculate: (a) the atomic radius of Fe; (b) the APF (verify it ≈ 0.68); (c) the theoretical density (use A = 55.85 g/mol).

2. Draw (using ASCII art or describe) the (110) plane in a BCC unit cell. Which atoms does it pass through?

3. A steel engineer wants to use austenite (FCC iron, γ-Fe) because it's more ductile. At what temperatures is γ-Fe stable? What problem occurs if they try to use it at room temperature?

4. List all 12 slip systems for FCC (4 {111} planes × 3 ⟨110⟩ directions each). For one plane (111), what are the three slip directions?

5. Magnesium (HCP) is very difficult to roll at room temperature but forms easily at 200°C. Explain this in terms of slip systems and how temperature activates additional ones.

6. Calculate the linear atomic density along the [110] direction in FCC aluminum (a = 0.4049 nm). Compare to the [100] direction. Which is more densely packed?
