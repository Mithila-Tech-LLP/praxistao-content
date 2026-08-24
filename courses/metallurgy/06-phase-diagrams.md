# Chapter 06: Phase Diagrams — The Maps of Materials

> **"A phase diagram is a thermodynamic GPS for the metallurgist. Tell it the temperature and composition, and it tells you exactly what phases are present, how much of each, and what the chemistry of each phase is. It is one of the most powerful tools in all of materials science."**

---

## Table of Contents

1. [What Is a Phase?](#1-what-is-a-phase)
2. [Why Phase Diagrams?](#2-why-phase-diagrams)
3. [One-Component (Unary) Diagrams](#3-one-component-unary-diagrams)
4. [Two-Component (Binary) Diagrams — The Basics](#4-two-component-binary-diagrams--the-basics)
5. [Isomorphous (Complete Solid Solubility) Systems](#5-isomorphous-complete-solid-solubility-systems)
6. [The Lever Rule — How Much of Each Phase?](#6-the-lever-rule--how-much-of-each-phase)
7. [Eutectic Systems](#7-eutectic-systems)
8. [Peritectic Systems](#8-peritectic-systems)
9. [Intermediate Phases and Intermetallics](#9-intermediate-phases-and-intermetallics)
10. [Reading a Phase Diagram — Step by Step](#10-reading-a-phase-diagram--step-by-step)
11. [Ternary Systems (Introduction)](#11-ternary-systems-introduction)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What Is a Phase?

A **phase** is a chemically and physically homogeneous portion of a system, bounded by surfaces across which properties change discontinuously.

Think of water at 0°C with ice floating in it: you have two phases (ice and liquid water) with the same chemical composition but different structure and properties. At 100°C, liquid water and steam (vapor) are two more phases.

For metals, examples of phases:
- **α-ferrite**: BCC iron with < 0.022 wt% C dissolved
- **Austenite (γ)**: FCC iron with up to 2.14 wt% C dissolved
- **Cementite (Fe₃C)**: the hard iron carbide compound
- **Liquid**: molten metal
- **γ′ (Ni₃Al)**: the ordered intermetallic precipitate in superalloys

The number of phases present at equilibrium is governed by the **Gibbs Phase Rule:**
```
F = C - P + 2

where:
  F = degrees of freedom (temperature, pressure, composition variables that can be changed
      without changing the number of phases)
  C = number of components (elements)
  P = number of phases present
  2 accounts for temperature and pressure
```

For most metallurgical systems (at constant atmospheric pressure): **F = C - P + 1**

At a **eutectic point** (one temperature, one composition, three phases in a two-component system):
F = 2 - 3 + 1 = 0 → no degrees of freedom → you're fixed in T and composition. This is an **invariant point**.

---

## 2. Why Phase Diagrams?

Phase diagrams encode enormous amounts of information:

1. **What phases are stable** at any T and composition
2. **How much of each phase** (via the lever rule)
3. **What chemistry each phase has** (via the tie line)
4. **Where transformations occur** (liquidus, solidus, solvus lines)
5. **Maximum solubility** of elements in each phase
6. **Reactions** — eutectic, peritectic, eutectoid, peritectoid

With a phase diagram, you can:
- Design heat treatment schedules (heat to single-phase region, cool to two-phase → controlled precipitation)
- Predict what happens when two alloys are welded
- Understand why a casting solidifies the way it does
- Design alloys with specific phase fractions

---

## 3. One-Component (Unary) Diagrams

With only one component, the variables are **temperature** and **pressure**. The phase diagram shows which phases are stable at each T-P combination.

**Water phase diagram:**
```
Pressure
  |       SOLID         LIQUID
  |       (ice)         (water)
  |          \    ←— melting curve
  |           \
  |            \
  |   Triple    ●──────────────────→  vapor curve
  |   point     (0.006 atm, 0.01°C)
  |
  └──────────────────────────────→ Temperature
```

At the **triple point**, all three phases coexist. One degree of freedom is used for each phase beyond the first — so three phases = F = 1 - 3 + 2 = 0 → invariant point.

**Iron unary diagram (simplified):**
```
Pressure                 
  1 atm ─────────────────────────────────────────
                    α-Fe  |  γ-Fe   |  δ-Fe  | Liquid
                    (BCC) | (FCC)   | (BCC)  |
                          |         |         |
                         912°C    1394°C    1538°C     Temperature
```

At atmospheric pressure, iron exists as α (BCC), γ (FCC), δ (BCC), or liquid depending on temperature. This allotropy is the foundation of all iron and steel metallurgy.

---

## 4. Two-Component (Binary) Diagrams — The Basics

With two components (A and B), we usually fix pressure at 1 atm and plot **temperature vs. composition**. Composition is expressed as:
- Weight percent (wt%): mass of B / total mass × 100%
- Atomic percent (at%): atoms of B / total atoms × 100%

**Key features of any binary phase diagram:**

```
Temperature
  |                      Liquidus (above = all liquid)
  |                ─────────────────────────
  |           ─────         L          ─────
  |      ─────        α+L         β+L       ─────
  |  ─────                                         ─────
  |           Solidus (below = all solid)
  |     α            α+β             β
  |                  ← two-phase region →
  |
  └──────────────────────────────────────────────────→ Composition (wt% B)
  0% B                                              100% B
  (pure A)                                         (pure B)
```

**Lines:**
- **Liquidus**: above this, all liquid
- **Solidus**: below this, all solid (in a simple system)
- **Solvus**: boundary of solid solubility in single-phase solid regions

**Regions:**
- Single-phase regions (liquid, α, β): only one phase present
- Two-phase regions (α+L, β+L, α+β): two phases coexist
- Three-phase reactions: at specific temperature lines (eutectic, peritectic, etc.)

---

## 5. Isomorphous (Complete Solid Solubility) Systems

The simplest binary diagram is the **isomorphous system**: two metals that are **completely soluble** in each other in both solid and liquid states.

Requirements (Hume-Rothery rules):
1. Atomic size difference < 15%
2. Same crystal structure
3. Similar electronegativity
4. Same valence (or similar)

**Cu-Ni system** — the classic example:

```
Temperature (°C)
  1455 ┤─────────────────────────────────── Ni melting point
       |          LIQUID
  1200 ┤   ┌──── Liquidus
       |   |    α+L region (two-phase!)
  1085 ┤   |    └──── Solidus
       |        α (Cu-Ni solid solution) — single phase, FCC
     0 ┤──────────────────────────────────
       0% Ni                           100% Ni
      (pure Cu)                       (pure Ni)
```

The "lens-shaped" two-phase region between liquidus and solidus is characteristic of all isomorphous systems.

**Reading this diagram:** If you have a 70 wt% Ni alloy at 1300°C, you're in the two-phase (α + L) region. The solid that's present doesn't have 70% Ni — it's richer in Ni (the higher-melting component). The liquid is leaner in Ni. This **composition difference between solid and liquid** is the origin of **solidification segregation**.

Other isomorphous systems: Fe-Ni, Fe-Cr (above certain temperatures), Au-Ag.

---

## 6. The Lever Rule — How Much of Each Phase?

When you're in a two-phase region, you know *which* phases are present from the diagram. But how much of each? The **Lever Rule** tells you.

**Setup:** Draw a horizontal **tie line** at the temperature of interest, across the two-phase region.

```
                  α                α+β           β
                  |                 |            |
T ────────────────●─────────────────────────────●────
                  |←─── Wα×(Cβ-C₀) ────→|←──── Wβ×(C₀-Cα)
                 Cα        C₀                  Cβ
```

The lever rule (imagine the tie line as a lever balanced at the overall composition C₀):

```
Weight fraction of β phase:  Wβ = (C₀ - Cα) / (Cβ - Cα)
Weight fraction of α phase:  Wα = (Cβ - C₀) / (Cβ - Cα)

Verify: Wα + Wβ = 1 ✓
```

**Example:** You have a Cu-40wt%Ni alloy at 1300°C.
- From the phase diagram, Cα (solid, Ni-rich) ≈ 58 wt% Ni
- Cβ (liquid, Ni-lean) ≈ 32 wt% Ni
- C₀ = 40 wt% Ni

```
Wsolid = (40 - 32) / (58 - 32) = 8/26 ≈ 0.31 (31% solid)
Wliquid = (58 - 40) / (58 - 32) = 18/26 ≈ 0.69 (69% liquid)
```

The lever rule applies in **any** two-phase field in any binary diagram.

---

## 7. Eutectic Systems

Many alloy systems have **limited solid solubility**: you can dissolve some B in A, and some A in B, but there's a maximum — the **solvus** boundary.

When solubility is limited, the phase diagram develops a **eutectic reaction**:

**L → α + β** (at constant temperature, from one liquid to two solids simultaneously)

```
Temperature
  ┤                  LIQUID
  ┤         ──────────────────────────
  ┤     ─── ← liquidus             ─── 
  ┤   α+L       ●Eutectic           β+L
  ┤─────────────────────────────────────  ← Eutectic temperature (invariant!)
  ┤      α              α+β              β
  ┤      
  └──────────────────────────────────────→ Composition
      0%B           Eutectic             100%B
                  composition
```

**The eutectic reaction:** At the eutectic temperature, the liquid of eutectic composition transforms into a **two-phase mixture** of α + β simultaneously. This reaction is invariant — it occurs at a single, fixed temperature.

**Why "eutectic"?** Greek for "easily melted" — the eutectic point is the lowest melting composition in the system.

**The eutectic microstructure:** The α and β phases form an alternating lamellar (plate-like) or fibrous pattern. The spacing of the lamellae depends on the cooling rate:
- Fast cooling → fine lamellae → harder (more interface area = more dislocation barriers)
- Slow cooling → coarse lamellae → softer

**Key eutectic systems in metallurgy:**
| System | Eutectic Comp. | Eutectic T | Application |
|--------|---------------|------------|-------------|
| Pb-Sn | 61.9 wt% Sn | 183°C | Soldering (traditional) |
| Al-Si | 12.6 wt% Si | 577°C | Die casting Al alloys |
| Fe-C | 4.3 wt% C | 1147°C | Cast iron |
| Au-Si | 97.1 wt% Si | 363°C | Electronics die-attach |

**Hypoeutectic vs. Hypereutectic:**
- **Hypoeutectic**: composition left of the eutectic point → primary α forms first, then eutectic
- **Hypereutectic**: composition right of the eutectic point → primary β forms first, then eutectic

---

## 8. Peritectic Systems

A **peritectic reaction** occurs when a solid phase reacts with liquid to form a different solid phase on cooling:

**L + α → β** (liquid + one solid → another solid)

```
Temperature
  ┤          LIQUID
  ┤       ──────────────
  ┤     α+L              β+L
  ┤─────────────────────────── ← Peritectic line
  ┤          α + β
  └──────────────────────────→
             ↑
         Peritectic point
```

Peritectic reactions are common in many engineering alloy systems and often cause difficulties during solidification because the β shell around α can kinetically block completion of the reaction.

The Fe-C system contains a peritectic at 1493°C:
**L + δ-Fe → γ-Fe** (at 0.18 wt% C in liquid, 0.09 wt% C in δ, forms γ at 0.17 wt% C)

---

## 9. Intermediate Phases and Intermetallics

Real binary systems often have **intermediate phases** that appear between the two pure elements:

- **Intermediate solid solutions**: chemically disordered, but stoichiometry is roughly fixed
- **Intermetallic compounds**: chemically ordered, specific stoichiometry (A₃B, AB, AB₂ etc.)

**Example — Fe-C system's intermediate phases:**
- **Cementite (Fe₃C)**: very hard, brittle intermetallic; the dark plates in pearlite
- Appears as a "pure component" on the right edge of the practical Fe-C diagram

**Intermetallics in superalloys:**
- **γ′ (Ni₃Al)**: the beneficial strengthening phase — coherent, cuboidal, stable
- **γ″ (Ni₃Nb)**: in IN718, the primary strengthener — metastable, transforms to δ above 650°C
- **TCP phases** (sigma σ, mu μ, Laves): harmful phases that appear in overaged or mis-designed superalloys — needles of brittle intermetallic that cause embrittlement

On phase diagrams, intermetallics appear as vertical lines (fixed composition) or narrow single-phase fields between the terminal solid solutions.

---

## 10. Reading a Phase Diagram — Step by Step

Let's practice reading a phase diagram systematically. Use the Pb-Sn eutectic system:

```
Temperature (°C)
  327  ┤─────●────────────────────────────── Pb melting point
       |  Liquid region
  232  ┤─────────────────────────────────●── Sn melting point
       |  α+L            ●eutectic    β+L
  183  ┤─────────────────────────────────────← eutectic T (invariant line)
       |      α              α+β         β
       | (Pb-rich)                   (Sn-rich)
    0  ┤─────────────────────────────────────
       0                  61.9              100
       Pure Pb          wt% Sn           Pure Sn
```

**Solvus:**
- Maximum solubility of Sn in Pb (α): 19.2 wt% Sn at 183°C; decreases at lower T
- Maximum solubility of Pb in Sn (β): ~2.5 wt% Pb at 183°C

**Step-by-step: What is a Pb-50wt%Sn alloy doing at 200°C?**
1. Find the composition (50 wt% Sn) on the x-axis
2. Go up to 200°C
3. You're in the **α + L** region (between liquidus and eutectic isotherm)
4. Draw a horizontal tie line at 200°C
5. Left end of tie line: ~18 wt% Sn → this is the α phase composition
6. Right end of tie line: ~56 wt% Sn → this is the L (liquid) phase composition
7. Lever rule: Wα = (56-50)/(56-18) = 6/38 ≈ 16%; WL = 84%
8. So the alloy is 84% liquid and 16% solid (α phase rich in Pb)

**On further cooling to 183°C:**
- The eutectic reaction occurs: remaining liquid (of eutectic composition 61.9 wt% Sn) → α + β
- Final microstructure: proeutectic α grains + eutectic mixture (alternating fine α and β lamellae)

**As you cool below 183°C:**
- The solvus restricts solubility. Any Sn in excess of the solvus precipitates as β within the α grains (or vice versa)
- This is precipitation — if done deliberately (solution treat then age) it's **precipitation hardening**

---

## 11. Ternary Systems (Introduction)

Real engineering alloys have many components. Nickel superalloys have 10–15 elements. The full phase diagram would be impossible to draw.

**Binary phase diagrams** (2 components) → easy to draw, 2D.
**Ternary phase diagrams** (3 components) → 3D representation (T on vertical axis, composition as a triangle). Can be shown as isothermal sections or vertical sections (pseudobinary sections).

**Ternary invariant reactions:**
In a ternary system: F = C - P + 1 = 3 - P + 1 = 4 - P.
For F = 0 (invariant): P = 4 phases → ternary eutectic at a specific T and composition.

**CALPHAD method:**
For 10-15 element superalloys, experimental determination of the full phase diagram is impossible. Instead:
- Thermodynamic data (Gibbs energies) for each binary and ternary pair are assessed from experimental data
- Software (Thermo-Calc, Pandat) computes the equilibrium phases at any T and composition by minimizing total Gibbs free energy
- This is the **CALPHAD** approach — covered fully in Chapter 64

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Phase | Homogeneous region with uniform properties, bounded by interfaces |
| Gibbs Phase Rule | F = C - P + 1 (at constant P). Invariant reactions have F=0 |
| Isomorphous | Complete solid solubility; lens-shaped two-phase region; Cu-Ni classic |
| Lever rule | Wβ = (C₀ - Cα)/(Cβ - Cα); phase fraction from composition on tie line |
| Eutectic | L → α + β; invariant; lowest melting composition; lamellar microstructure |
| Peritectic | L + α → β; solid reacts with liquid; can cause kinetic complications |
| Solvus | Boundary of solid solubility; decreasing solubility → precipitation on cooling |
| Intermetallics | Fixed stoichiometry ordered compounds; vertical lines on diagram |
| Ternary | 3-component; shown as isothermal triangles or pseudobinary sections |
| CALPHAD | Computational phase diagram calculation; essential for multicomponent alloys |

**Next chapter:** The most important phase diagram in all of metallurgy — the **iron-carbon diagram** — which underlies every steel that has ever been made.

---

## Exercises

1. A Cu-40 wt% Ni alloy is heated to 1350°C. Using the Cu-Ni diagram (liquidus ≈ 1340°C, solidus ≈ 1270°C at 40 wt% Ni), determine: (a) is the alloy fully liquid, fully solid, or two-phase? (b) What are the compositions of the two phases? (c) What fraction is liquid (lever rule)?

2. In the Pb-Sn system, starting from a fully liquid alloy of 30 wt% Sn, describe step by step what happens as you cool from 300°C to room temperature. What is the final microstructure?

3. Using the Gibbs Phase Rule, verify that an eutectic point (2-component system, 3 phases: α, β, liquid) is truly invariant (F=0) at constant pressure.

4. A steel with 4.3 wt% C is the eutectic composition (cast iron eutectic). What reaction occurs at 1147°C? What is the name of the product phases?

5. CALPHAD question: Why can't you simply read off the stable phases in a 14-component Ni superalloy from existing binary diagrams? What thermodynamic quantity must be minimized to find the equilibrium phase assemblage?
