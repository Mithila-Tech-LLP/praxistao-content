# Chapter 18: Annealing, Normalizing, and Stress Relief

> **"A steel part fresh from forging or cold-rolling is internally tortured — locked stresses, tangled dislocations, irregular grain shapes, hard second phases. Annealing is the release. It is the metallurgist's method of resetting the material — not to its original state, but to a well-defined, controlled state from which the next step in manufacturing can proceed."**

---

## Table of Contents

1. [Overview — Why Soften or Relieve?](#1-overview--why-soften-or-relieve)
2. [Recovery — The First Stage of Softening](#2-recovery--the-first-stage-of-softening)
3. [Recrystallization — New Grains Grow](#3-recrystallization--new-grains-grow)
4. [Grain Growth After Recrystallization](#4-grain-growth-after-recrystallization)
5. [Full Annealing (Steel)](#5-full-annealing-steel)
6. [Subcritical Annealing and Spheroidize Annealing](#6-subcritical-annealing-and-spheroidize-annealing)
7. [Process Annealing and Intermediate Annealing](#7-process-annealing-and-intermediate-annealing)
8. [Normalizing](#8-normalizing)
9. [Stress Relief Annealing](#9-stress-relief-annealing)
10. [Annealing of Non-Ferrous Metals](#10-annealing-of-non-ferrous-metals)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Overview — Why Soften or Relieve?

After cold working, machining, welding, or casting, metals often contain:
- High dislocation density → high hardness + low ductility
- Residual stresses → distortion on further machining, stress corrosion risk
- Non-equilibrium microstructures (martensite from quenching, microsegregation from casting)
- Large, irregular grains from solidification

We use thermal treatments to reverse these conditions:

| Condition to Remove | Treatment | Temperature |
|--------------------|-----------|-------------|
| Cold work (high dislocation density) | Full anneal or process anneal | 0.6–0.8 T_m |
| Residual stress only | Stress relief | 0.3–0.5 T_m |
| Hard martensite → soft | Subcritical anneal + spheroidize | 650–700°C |
| Casting microsegregation | Homogenization anneal | 0.85–0.95 T_m |
| Coarse equiaxed to fine | Normalizing | A₃ + 50°C, air cool |

---

## 2. Recovery — The First Stage of Softening

When a cold-worked metal is heated to relatively low temperatures (0.3–0.4 T_m), the first changes occur:

**What happens in recovery:**
1. Vacancies migrate and annihilate (at grain boundaries, surfaces, dislocations)
2. Dislocation rearrangement: random tangles organize into lower-energy configurations (subgrain walls)
3. Dislocation density drops slightly but grain boundaries are unchanged

**Subgrain formation (polygonization):**
Edge dislocations with the same sign line up vertically → form low-angle subgrain boundaries (tilt walls). The crystal bends (curved dislocations) straightens into a mosaic of slightly misoriented blocks:

```
Before recovery (bent lattice):    After recovery (polygonized):

         ↕ bent crystal                  │ │ │ subgrain 
         ↕                               │ │ │ walls
         ↕ ← curved dislocations         │ │ │ (low-angle 
                                         │ │ │ boundaries)
```

**Properties after recovery:**
- Hardness drops slightly (~10–20%)
- Ductility increases slightly
- Residual stresses are mostly relieved (major effect)
- Grain size unchanged
- No recrystallization yet

---

## 3. Recrystallization — New Grains Grow

At higher temperatures (0.4–0.5 T_m for most metals), new, dislocation-free grains nucleate and grow at the expense of the deformed (high dislocation density) matrix:

**Driving force:** Reduction in stored energy from dislocations (1–10 J/g stored in heavily cold-worked metals).

**Recrystallization temperature:** The temperature at which 50% recrystallization occurs in 1 hour:
```
Approximate T_recryst ≈ 0.4 × T_m (K)  for pure metals

Examples:
  Lead (T_m = 327°C):  T_recryst ≈ -17°C → recrystallizes at RT!
  Tin:                  T_recryst ≈ -7°C
  Aluminum (660°C):     T_recryst ≈ 148°C
  Iron (1,538°C):       T_recryst ≈ 454°C
  Nickel (1,455°C):     T_recryst ≈ 430°C
  Tungsten (3,422°C):   T_recryst ≈ 1,205°C
```

Alloying elements RAISE T_recryst (solute drag on grain boundary migration).

### Recrystallization Kinetics (Avrami Equation):

```
f_recryst = 1 - exp(-k × t^n)

where f = fraction recrystallized
      t = time
      k, n = constants depending on T and material
```

S-shaped curve in time: slow start (nucleation), fast middle (growth), slow finish (impingement).

### Recrystallization Temperature Factors:

| Factor | Effect |
|--------|--------|
| Higher cold work % | Lower T_recryst (more stored energy, easier nucleation) |
| Higher temperature | Faster (both nucleation and growth rates increase) |
| Alloying elements | Raise T_recryst (solute drag) |
| Fine original grain | Lower T_recryst (more GB area → more nucleation sites) |
| Longer time | Complete at lower T |

### Critical Deformation

Small amounts of cold work (1–5%) can cause **abnormal grain growth** (exaggerated grain growth after recrystallization) because few nuclei form. The recrystallized grains are enormous and non-uniform. DANGEROUS for turbine blade castings: any surface damage (grit blasting) → recrystallization → grain boundaries in SX → REJECTION (Ch 59).

---

## 4. Grain Growth After Recrystallization

Once recrystallization is complete (all deformed structure consumed), further holding at high temperature causes grain growth (Ch 09):
```
d² - d₀² = K × t × exp(-Q/RT)
```

**Controlling grain size:**
- Stop annealing when recrystallization is complete (don't hold longer at high T)
- Use grain-growth inhibiting particles (NbC, TiN) — "microalloyed" steels
- Control soaking temperature: higher T → faster grain growth

---

## 5. Full Annealing (Steel)

**Purpose:** Produce maximum softness, ductility, and eliminate all residual stress. Used before machining complex shapes.

**Process:**
```
1. Heat to above A₃ (or A₁ for hypereutectoid): 30–50°C above A₃
   → forms austenite (or austenite + cementite for hypereutectoid)
2. Hold until fully austenitized (1 hour per 25mm section)  
3. Cool SLOWLY in furnace (< 30°C/hour)
   → austenite transforms to coarse pearlite + ferrite or pearlite + cementite
   → equilibrium microstructure
```

**Resulting microstructure:** Coarse pearlite (thick lamellae) + proeutectoid ferrite (hypoeutectoid steel). Very soft.

**Typical hardness:** HRB 60–85 (very soft, 100–160 HV)

**Limitation:** Slow furnace cooling is time-consuming (24–48 hours for thick sections).

---

## 6. Subcritical Annealing and Spheroidize Annealing

### Subcritical Annealing

Heat to just BELOW A₁ (600–700°C) and hold:
- No austenite forms
- Recovery + partial recrystallization occurs
- Faster than full annealing
- Softens cold-worked steel but doesn't produce equilibrium structure

### Spheroidize Annealing

Purpose: convert the lamellar cementite in pearlite into spherical (spheroidized) cementite particles in a ferrite matrix — MAXIMUM machinability for high-carbon steels.

**Process:**
```
Method 1: Hold just below A₁ (700–720°C) for 15–24 hours
  → cementite lamellae are thermodynamically unstable at near-A₁ →
  → carbide edges dissolve, spheres are more stable → spheroidization

Method 2: Cycle above and below A₁
  → repeatedly form and dissolve fine carbides → faster spheroidization
```

**Why spheroidize:**
- Ball bearings (52100 steel): machined soft (spheroidized, ~200 HV), then hardened
- Hypereutectoid steels (tool steels): continuous cementite network → brittle; must spheroidize first to dissolve network

**Result:**
```
Pearlite (lamellar):        Spheroidized:
  ─────────── Fe₃C           ●  ●  ●  ← Fe₃C spheres in ferrite
  ─────────── Fe              ●  ●  ●
  ─────────── Fe₃C
  
Hardness: 200–240 HV         Hardness: 160–190 HV
```

---

## 7. Process Annealing and Intermediate Annealing

**Process annealing** (also "cold-work annealing" or "bright annealing"): used between cold-working passes to restore ductility for further working.

- Temperature: 0.5–0.7 T_m (below A₁ for steels → no austenite forms)
- Purpose: recrystallize cold-worked structure → restore ductility → can cold-work again
- Atmosphere: controlled (H₂ or N₂) → bright surface finish

**Example — Copper/brass strip:**
```
Cold roll → harden → bright anneal (hydrogen atmosphere, 500–700°C) → recrystallize → soft again → cold roll again
Repeat until final thickness and temper achieved.
```

**Temper designations (e.g., copper alloys):**
- O (annealed/soft): recrystallized, minimum hardness
- H01 (¼ hard): 10–15% cold work after anneal
- H02 (½ hard): 20–25% cold work
- H04 (hard): 35–40% cold work
- H08 (spring): maximum cold work

---

## 8. Normalizing

**Purpose:** Produce a more refined, uniform grain size than as-cast or as-forged structure. Improve toughness. Ensure consistent properties throughout a component.

**Process:**
```
1. Heat to A₃ + 30–60°C (typically 850–950°C for carbon steels)
2. Hold until uniform temperature throughout
3. AIR COOL (not furnace cool)
   → faster than full anneal → finer pearlite
   → austenite → fine pearlite + ferrite
```

**Key difference from annealing:** Air cooling is faster → finer pearlite spacing → harder, higher tensile strength, better toughness.

**Normalized vs. annealed comparison:**
| Property | Annealed | Normalized |
|---------|----------|------------|
| Tensile strength | Lower | Higher by 10–20% |
| Yield strength | Lower | Higher |
| Hardness | Lower | Higher |
| Toughness | Lower | Higher (finer pearlite) |
| Machinability | Best | Good |
| Ductility | Highest | Slightly lower |

**Application:** Structural steel forgings and castings are routinely normalized as a quality-assurance step before delivery to ensure uniform microstructure regardless of section size variations during forge cooling.

---

## 9. Stress Relief Annealing

**Purpose:** Reduce residual stresses WITHOUT changing the microstructure (no phase transformations, no significant grain growth, no significant softening).

**Sources of residual stress:**
- Welding (HAZ contracts differently from base metal)
- Machining (surface layers compressed by tool forces)
- Casting (outer surface cools and contracts faster than center)
- Shot peening (deliberate beneficial compressive residual stress)
- Forming and bending (outer fibers tensile, inner fibers compressive)

**Process:**
```
Temperature: 0.3–0.5 T_m (typically 550–700°C for steel, 150–200°C for cold-worked Al)
Hold: 1–4 hours (depending on section size and target residual stress level)
Cool: Slowly (< 100°C/hour) to avoid introducing NEW thermal stresses during cooling!
```

**Effect:** Recovery annihilates vacancy clusters and allows slight dislocation rearrangement. Most residual stress is relieved without recrystallization (which would change grain size and properties).

**Post-weld heat treatment (PWHT):** Mandatory for many structural welds:
- Pressure vessels: ASME code requirements
- Steam lines: 593–650°C (Cr-Mo steels) for 1 hour/inch thickness
- Nuclear components: strict PWHT requirements

**Caution — sensitization in stainless steels:** Heating 304 SS at 450–850°C for >10 minutes causes chromium carbide precipitation at grain boundaries → depletes Cr near boundaries → susceptible to intergranular corrosion. Use L-grade (304L, 316L, max 0.03%C) or stabilized grades (321, 347) for post-weld service.

---

## 10. Annealing of Non-Ferrous Metals

### Aluminum Alloys

- Solution annealing: heat to solid solution field → quench → metastable supersaturated solid solution
- Recrystallization anneal: heat non-heat-treatable alloys (1xxx, 3xxx, 5xxx) to 300–420°C → restore ductility
- Homogenization of ingots: 500–600°C for 12–24 hours → eliminate dendritic microsegregation

### Copper Alloys

- Recrystallization of cold-worked brass (70/30): 300–600°C → grain size control
- Stress relief of bent copper tubes: 150–300°C → eliminate residual stress → prevent stress corrosion cracking in ammonia environments

### Titanium Alloys

- Stress relief: 480–650°C (below beta transus)
- Annealing: 700–815°C (full anneal) → equiaxed α + β structure
- Recrystallization annealing: near beta transus → larger, equiaxed β

### Nickel Superalloys

- "Solution treatment" of γ′-hardened alloys: 1,100–1,320°C → dissolve γ′ → quench → then re-age
- The process is so critical that it is discussed separately in Ch 49 (SX post-processing)

---

## Summary

| Treatment | Temperature | Key Effect | Microstructure |
|-----------|-------------|-----------|----------------|
| Stress relief | 0.3–0.5 T_m | Remove residual stress, no phase change | Unchanged + less stress |
| Recovery | 0.3–0.4 T_m | Dislocation rearrangement, slight softening | Same grains, subgrains form |
| Recrystallization | 0.4–0.5 T_m | New strain-free grains, major softening | New fine equiaxed grains |
| Full annealing | A₃+30°C, furnace cool | Maximum softness, coarse pearlite | Coarse pearlite + ferrite |
| Spheroidize | 700°C, 15+ hours | Maximum machinability | Fe₃C spheres in ferrite |
| Normalizing | A₃+50°C, air cool | Fine pearlite, better toughness than anneal | Fine pearlite + ferrite |

---

## Exercises

1. A 50mm thick steel plate has been welded. Residual tensile stress at the weld = 400 MPa (near yield strength). Describe the PWHT procedure to reduce residual stress to < 100 MPa. Include: (a) temperature, (b) heating rate, (c) hold time, (d) cooling rate, (e) what would happen if you cooled too fast. Calculate the required hold time if guidance says 1 hour per 25mm thickness.

2. Copper wire is cold-drawn 60% (area reduction from 5 mm² to 2 mm²). Tensile strength is now 500 MPa, elongation = 5%. You need to restore ductility (elongation > 25%) before further cold drawing. Design an annealing process: (a) target temperature range (T_m,Cu = 1,085°C; T_recryst ≈ 200°C for heavily worked Cu), (b) how to check that recrystallization is complete (hardness testing), (c) why controlled atmosphere (H₂ or N₂) is preferred over air annealing for final-pass wire.

3. A hypereutectoid steel bearing ring (1.5% C) must be machined before hardening. Current microstructure: pearlite + proeutectoid Fe₃C network at grain boundaries. Machinability is poor and there is risk of cracking during machining. Design a spheroidize anneal to: (a) dissolve the continuous Fe₃C network, (b) achieve spheroidized carbides for maximum machinability, (c) estimate the hardness before and after spheroidizing (use: pearlite+cementite ≈ 280 HV, spheroidized ≈ 190 HV).

4. Calculate the recrystallization temperature for: (a) aluminum (T_m = 660°C), (b) iron (T_m = 1,538°C), (c) tungsten (T_m = 3,422°C) using T_recryst ≈ 0.4 T_m in Kelvin. Express in °C. Explain why: (i) aluminum window frames slightly soften after years of service, (ii) tungsten light bulb filaments can operate at 2,500°C without recrystallizing immediately (but do fail eventually), (iii) lead pipe connections made by soldering can creep at room temperature.

5. During SX turbine blade manufacturing, the ceramic mold is removed by grit blasting. Explain: (a) what metallurgical damage grit blasting introduces to the SX blade surface (hint: local plastic deformation, dislocation density, residual stress), (b) why this is dangerous during the subsequent solution heat treatment at 1,300°C, (c) what alternative mold-removal technique avoids this problem, (d) how a recrystallized zone is detected during quality inspection.
