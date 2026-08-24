# Chapter 09: Grain Structure and Grain Boundaries

> **"A polycrystalline metal is a community of crystals, each with its own orientation, each separated from its neighbor by a thin zone of atomic disorder. These boundaries — a few atoms thick — control whether a metal is strong or brittle, whether it creeps slowly or fails suddenly, whether a turbine blade lasts 30,000 flights or fails on takeoff."**

---

## Table of Contents

1. [What Is a Grain?](#1-what-is-a-grain)
2. [How Grains Form — Nucleation and Competitive Growth](#2-how-grains-form--nucleation-and-competitive-growth)
3. [Grain Boundaries — Structure and Energy](#3-grain-boundaries--structure-and-energy)
4. [Types of Grain Boundaries](#4-types-of-grain-boundaries)
5. [Grain Growth — Why Grains Coarsen](#5-grain-growth--why-grains-coarsen)
6. [Hall-Petch Relationship — Grain Size and Strength](#6-hall-petch-relationship--grain-size-and-strength)
7. [ASTM Grain Size Number](#7-astm-grain-size-number)
8. [Grain Boundary Segregation](#8-grain-boundary-segregation)
9. [Columnar vs. Equiaxed Grain Structures](#9-columnar-vs-equiaxed-grain-structures)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is a Grain?

A **grain** is a single crystal region within a polycrystalline material. Each grain has:
- A unique crystallographic orientation (its crystal lattice points in a specific direction)
- A size typically 10 μm – 10 mm (depending on alloy and processing history)
- A shape influenced by how many neighbors it has and how it grew

```
Microstructure of a polycrystalline metal (schematic):

    ─────────────────────────────
   /      /      /    \    \
  / Grain /  Grain\  Grain \ Grain
 /   A   /    B    \   C   /
/───────/─────────  \─────/──────
│ Grain │   Grain    │  Grain   │
│   D   │     E      │    F     │
└───────────────────────────────┘

Each region = one grain (single crystal, one orientation)
Lines between = grain boundaries (zone of atomic disorder)
```

When you look at a polished and etched metal surface under a microscope, the grain boundaries appear as dark lines because the etchant attacks the disordered boundary region preferentially.

---

## 2. How Grains Form — Nucleation and Competitive Growth

During solidification (Ch 08), many nuclei form at different locations, each with a random crystallographic orientation. They grow outward until they impinge on each other. Where two growing crystals meet, a grain boundary forms.

**Competitive growth:** During columnar grain growth (directional solidification), grains with their fast-growth direction [001] aligned with the heat flow direction grow faster and "outcompete" (consume) misoriented neighbors:

```
Grain selection during DS:

Start: many small grains at base with random orientations
  ↓↓↓↓↓↓↓↓↓↓↓↓   (many competing grains)
  
  ↓  ↓↓ ↓  ↓↓↓
  
  ↓  ↓↓    ↓↓↓   (misoriented grains die out)
  
    ↓↓    ↓↓     (only well-aligned remain)
  
      ↓↓↓        (few dominant columnar grains)
```

**Grain refinement:** Adding inoculants (grain refiners like Al-Ti-B for aluminum) creates many nucleation sites → many more, smaller grains. This is desirable for most structural applications (Hall-Petch strengthening, §6).

---

## 3. Grain Boundaries — Structure and Energy

A grain boundary is the transition zone between two grains with different crystallographic orientations. The atoms in this zone are NOT in a perfect crystal lattice — they are in a disordered, higher-energy arrangement.

**Misorientation angle θ:** The angle between the crystal lattices of the two neighboring grains, measured about a specific rotation axis.

**Grain boundary energy γ_GB:** Energy per unit area of the boundary:
```
γ_GB ≈ 0.2–1.0 J/m²  (typical metals)

For comparison: γ_surface ≈ 1–3 J/m²
               γ_liquid/solid ≈ 0.1–0.5 J/m²
```

Higher misorientation → more atomic disorder → higher energy.

**Why grain boundaries matter:**
- Barrier to dislocation motion → strengthening (Hall-Petch)
- Fast diffusion paths (D_GB ~ 10⁴× D_bulk at low T)
- Segregation sites for impurities
- Nucleation sites for second phases
- Crack initiation under fatigue/creep loading
- For SX turbine blades: ZERO grain boundaries = eliminate all these problems (Ch 47)

---

## 4. Types of Grain Boundaries

### By Misorientation Angle

**Low-angle grain boundary (LAGB): θ < 15°**
Can be described as an array of dislocations. For a simple tilt boundary:
```
LAGB structure (tilt boundary):
        
        |  ↑  |  ↑  |  ↑  |
Grain A |  ⊥  |  ⊥  |  ⊥  | Grain B
        ↑ = edge dislocations in a vertical array
        
Dislocation spacing D = b / θ (for small θ)
where b = Burgers vector, θ = misorientation angle
```

At 5° misorientation, D ≈ 3 nm (dislocations very close together).

**High-angle grain boundary (HAGB): θ > 15°**
Dislocations overlap so much they lose identity. Truly disordered transition zone 2–5 atoms thick.

**Special grain boundaries:**

**Coincidence Site Lattice (CSL) boundaries:** When two crystal orientations happen to share a fraction (1/Σ) of lattice sites. For example, Σ3 (60° rotation about [111] in FCC) = twin boundary, where 1/3 of lattice sites coincide. Very low energy.

**Twin boundaries:** Very common in FCC metals (Cu, stainless steel). Created during annealing. Very low energy (γ ≈ 0.02–0.04 J/m²). Appear as straight parallel lines in the microstructure.

### By Plane Relative to Crystal

**Tilt boundary:** Rotation axis in the boundary plane
**Twist boundary:** Rotation axis perpendicular to boundary plane
**General boundary:** Mixed character

---

## 5. Grain Growth — Why Grains Coarsen

A system with many grain boundaries has higher energy than a system with fewer. Given enough thermal energy (high temperature), grain boundaries migrate to reduce total boundary area → grains coarsen.

**Driving force:** Reduction in total boundary energy
**Mechanism:** Boundary curvature drives migration — curved boundaries migrate toward their center of curvature. Larger grains (less curved boundaries) consume smaller grains.

**Grain growth kinetics:**
```
Normal grain growth: d² - d₀² = K × t × exp(-Q/RT)

where d = grain diameter at time t
      d₀ = initial grain diameter
      K = constant
      Q = activation energy (related to grain boundary diffusion)
      R = gas constant, T = temperature (K)
```

**Modified: d^n - d₀^n = K × t**

Experimentally, n ≈ 2–4 for most metals (theoretical = 2, second-phase particles slow growth → higher n).

**Abnormal grain growth:** A few grains grow very much faster than others → coarse mixed grain structure → BAD for fatigue life (stress concentrators at large grains).

**Grain growth inhibition by pinning particles (Zener pinning):**
Second-phase particles (e.g., TiC, NbC in steel; γ′ in superalloys) pin grain boundaries by opposing their migration. Maximum stable grain size:
```
d_max = (4/3) × r / f_v

where r = particle radius, f_v = volume fraction of particles
```

This is why steel is often microalloyed with Nb, Ti, V — their carbides/nitrides pin grain boundaries at hot rolling temperatures → fine grain size retained → better combination of strength and toughness.

---

## 6. Hall-Petch Relationship — Grain Size and Strength

One of the most important relationships in physical metallurgy:

```
σ_y = σ_0 + k_y / √d

where σ_y = yield strength (MPa)
      σ_0 = friction stress (intrinsic resistance to dislocation motion in a single crystal)
      k_y = Hall-Petch slope (MPa × √m), material constant
      d = grain diameter (m or mm)
```

**Why grain boundaries strengthen:** A dislocation moving through a grain is blocked at the grain boundary. The dislocation pile-up creates a stress concentration that must be large enough to generate a new dislocation in the adjacent grain. Finer grains = more boundaries = shorter slip distances = higher stress needed to propagate slip.

**Example for steel:**
```
σ_0 ≈ 70 MPa, k_y ≈ 0.74 MPa√m

d = 1 mm:   σ_y = 70 + 0.74/√0.001 = 70 + 23 = 93 MPa
d = 100 μm: σ_y = 70 + 0.74/√0.0001 = 70 + 74 = 144 MPa
d = 10 μm:  σ_y = 70 + 0.74/√0.00001 = 70 + 234 = 304 MPa
d = 1 μm:   σ_y = 70 + 0.74/√0.000001 = 70 + 740 = 810 MPa
```

10× reduction in grain diameter → 3× increase in strength (from Hall-Petch alone).

**IMPORTANT — Direction matters:** The Hall-Petch relationship applies to yielding (where grain boundaries block dislocations). But grain boundaries are WEAK POINTS for:
- Creep (grain boundary sliding) → single crystals avoid this (Ch 47)
- High-temperature fatigue (grain boundary oxidation/cracking)
- Hydrogen embrittlement (H concentrates at boundaries)
- Temper embrittlement (impurities segregate to boundaries)

This is why single-crystal turbine blades are stronger under creep despite having NO Hall-Petch strengthening — the benefits of eliminating grain boundaries outweigh the loss of grain-boundary strengthening at high temperature.

---

## 7. ASTM Grain Size Number

The American Society for Testing and Materials (ASTM) defines a grain size scale for standardized reporting:

```
G = grain size number (ASTM E112)
n = N × 2^(G-1)  where N = number of grains per in² at 100× magnification

G = 1:  very coarse (~1 grain/mm²)
G = 5:  moderate (~100 grains/mm²)  
G = 8:  fine (~1,000 grains/mm²)    ← common engineering target
G = 12: very fine (~100,000 grains/mm²)
```

**Mean chord length** (more rigorous):
```
d_mean = 1 / (N_L × M)

where N_L = number of grain boundary intercepts per unit length of test line
      M = magnification factor
```

**Relationship to actual diameter:**
```
G = 1 → d ≈ 1 mm
G = 5 → d ≈ 150 μm
G = 8 → d ≈ 22 μm
G = 10 → d ≈ 11 μm
G = 14 → d ≈ 2.8 μm
```

---

## 8. Grain Boundary Segregation

Elements with atomic sizes different from the solvent tend to migrate to grain boundaries where the disordered structure accommodates their misfit more easily:

**Typical segregating elements:**
- **Strengthening/beneficial:** C, N (in steel) → carbide/nitride precipitation at boundaries
- **Harmful:** P, S, Sn, As, Sb, Bi in steel → temper embrittlement
- **Hydrogen:** segregates to boundaries → hydrogen-induced cracking
- **Boron in superalloys (PC, DS alloys):** strengthens grain boundaries → ESSENTIAL for resistance to grain boundary cracking → deliberately added to PC/DS alloys but NOT to SX (no boundaries to strengthen)

**Temper embrittlement (practical example):**
Steel quenched and then tempered at 400–600°C → P, Sn, Sb slowly segregate to grain boundaries → embrittles the steel → impact toughness drops dramatically. Avoided by: rapid cooling through 300–600°C range, minimizing P/Sn/Sb in steel composition.

---

## 9. Columnar vs. Equiaxed Grain Structures

| Feature | Equiaxed | Columnar |
|---------|----------|---------|
| Grain shape | Roughly equal dimensions all directions | Elongated along growth/working direction |
| Origin | Many inoculant nuclei; light deformation + anneal | Directional solidification; heavy cold work |
| Strength | Isotropic (same in all directions) | Anisotropic (stronger along axis) |
| Ductility | More ductile transverse | Poor transverse ductility |
| Creep | More grain boundary area normal to stress → GBS | Fewer transverse boundaries → better creep |
| Applications | Most structural components | DS turbine blades, wire (cold drawn) |
| Grain size (ASTM) | G = 5–12 typical | Not applicable (very elongated) |

**Wrought vs. Cast grain structures:**
After solidification (columnar or mixed), hot/cold working breaks up the cast grain structure. Recrystallization during annealing creates new equiaxed grains. Multiple working + anneal cycles refine grain size progressively.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Grain | Single crystal region; bounded by grain boundaries; size 10 μm – 10 mm typically |
| LAGB | θ < 15°; array of dislocations; < 5° acceptable in SX castings |
| HAGB | θ > 15°; fully disordered; strongest barrier to dislocations; weakest creep link |
| Hall-Petch | σ_y = σ_0 + k_y/√d; fine grains → stronger; finer by 10× → stronger by ~3× |
| Grain growth | d² ∝ t × exp(-Q/RT); inhibited by second-phase particles (Zener pinning) |
| Segregation | P, S to grain boundaries → embrittlement; B added to strengthen boundaries in PC alloys |
| Columnar | Better creep (fewer transverse GBs); used in DS blades; anisotropic properties |

---

## Exercises

1. A steel has σ_0 = 70 MPa and k_y = 0.74 MPa√m. Calculate the yield strength for grain sizes of: (a) 200 μm (annealed), (b) 50 μm (normalized), (c) 15 μm (fine-grained). What is the percentage increase in yield strength going from 200 μm to 15 μm?

2. A steel's grain boundaries are pinned by NbC particles: radius r = 50 nm, volume fraction f_v = 0.002. Using d_max = (4/3) × r / f_v, calculate the maximum grain size that can be maintained. If the temperature increases and the NbC dissolves to f_v = 0.0005, what happens to grain size? Why is this relevant to hot rolling control?

3. An austenitic steel undergoes grain growth at 1,100°C. Initial grain size d₀ = 20 μm. Using d² - d₀² = K×t with K = 5 × 10⁻¹⁴ m²/s, calculate: (a) grain size after 1 hour, (b) grain size after 10 hours, (c) how long to reach d = 500 μm (turbine nozzle vane grain size). Does the grain size scale linearly with time? What does this mean for heat treatment window control?

4. Compare polycrystalline Ni (grain size 50 μm) and single crystal Ni oriented [001] for tensile loading. At room temperature: (a) predict which has higher tensile yield strength and why; (b) at 1,000°C in creep, which would you expect to perform better and why? (c) What does this tell you about why SX blades are used despite losing Hall-Petch strengthening?

5. A metallurgist observes that a batch of steel shows intergranular fracture (crack follows grain boundaries) after heat treatment at 450°C for 4 hours. List: (a) three likely metallurgical explanations (temper embrittlement, hydrogen embrittlement, grain boundary precipitation), (b) how to distinguish between them using fractography + chemistry, (c) what process changes could prevent each.
