# Chapter 20: Tempering and Case Hardening

> **"Fresh martensite is like a coiled spring loaded to breaking point — so stressed, so full of trapped energy, that it will crack or shatter without warning. Tempering is the art of releasing that energy carefully, trading some hardness for the resilience to survive in service. And case hardening is the art of having it both ways: a surface hard enough to resist wear, a core tough enough to absorb shock."**

---

## Table of Contents

1. [Why Temper? The Problem With Fresh Martensite](#1-why-temper-the-problem-with-fresh-martensite)
2. [The Four Stages of Tempering](#2-the-four-stages-of-tempering)
3. [Tempering Temperature vs. Properties](#3-tempering-temperature-vs-properties)
4. [Temper Embrittlement](#4-temper-embrittlement)
5. [Double Tempering and Secondary Hardening](#5-double-tempering-and-secondary-hardening)
6. [Case Hardening Overview](#6-case-hardening-overview)
7. [Carburizing](#7-carburizing)
8. [Nitriding](#8-nitriding)
9. [Carbonitriding and Ferritic Nitrocarburizing](#9-carbonitriding-and-ferritic-nitrocarburizing)
10. [Induction Hardening and Flame Hardening](#10-induction-hardening-and-flame-hardening)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Temper? The Problem With Fresh Martensite

Fresh (as-quenched) martensite is:
- Very hard (up to 66 HRC at high C content)
- Very BRITTLE (essentially zero ductility at >0.4%C)
- Full of residual stresses (from rapid quenching + transformation volume change)
- At risk of spontaneous cracking if not tempered promptly

**Tempering** reduces brittleness and residual stress by heating the martensite to 100–700°C:
- Carbon precipitates as carbides → less distorted BCC (closer to equilibrium)
- Internal stresses relax
- Ductility improves dramatically
- Hardness decreases (the trade-off)

The final microstructure after tempering is called **tempered martensite** — the workhorse of high-strength structural steels.

---

## 2. The Four Stages of Tempering

As tempering temperature increases, four sequential transformations occur:

### Stage 1: 100–250°C — Carbon Segregation and ε-Carbide Formation
- C atoms migrate from octahedral sites to defects (dislocations, subgrain boundaries)
- Fine ε-carbide (Fe₂.₄C, hexagonal) precipitates
- Slight hardness reduction (2–4 HRC)
- Martensite tetragonality decreases

### Stage 2: 200–300°C — Retained Austenite Decomposition
- Retained austenite transforms to lower bainite-like structure
- Plate martensite → some conversion
- Hardness may INCREASE slightly (if significant retained austenite was present)

### Stage 3: 250–350°C — ε-Carbide → Cementite (Fe₃C) Conversion
- ε-carbide dissolves, cementite (Fe₃C) forms
- Martensite tetragonality fully disappears → cubic ferrite
- Major softening begins

### Stage 4: 350–700°C — Cementite Coarsening and Grain Recovery
- Fe₃C spheroidizes and coarsens (Ostwald ripening)
- Ferrite grains recover and recrystallize (at higher T)
- Significant softening → high ductility

```
Tempering temperature effects (0.4%C steel):

Hardness  Ductility  Toughness  Application
200°C: 56 HRC  ↓low      ↓low       Files, cutting tools
300°C: 50 HRC  ↑         ↑          Springs (note: avoid 350–550°C TE zone!)
400°C: 44 HRC  moderate  moderate   Axes, chisels
500°C: 38 HRC  good      good       Structural parts
600°C: 30 HRC  high      high       Heavy machinery shafts
650°C: 25 HRC  very high very high  Q&T structural steel
```

---

## 3. Tempering Temperature vs. Properties

**Tempering parameter (Hollomon-Jaffe parameter P):**
Captures the combined effect of temperature and time:
```
P = T × (log t + C)

where T = temperature in Kelvin
      t = time in hours
      C ≈ 14–20 (material constant, ~17 for many steels)
```

Higher P → more tempering effect (same as higher T or longer t).

**Tempered hardness prediction:**
```
For plain carbon steels:
HRC_tempered ≈ HRC_max × (1 - 0.001 × P)  (approximate)

More precisely: use tempering curves from steel manufacturers' datasheets
```

**Typical tempering practice:**
- Temper IMMEDIATELY after quench (before martensite cracks)
- Temperature ± 10°C for precision parts
- Minimum time: typically 1 hour per 25mm of thickness
- For critical aerospace parts: 2 × 2-hour tempers at temperature

---

## 4. Temper Embrittlement

Two distinct embrittlement phenomena:

### Tempered Martensite Embrittlement (TME) — "350°C embrittlement"
- Occurs in carbon steels tempered at 250–350°C
- Caused by formation of thin ε-carbide and Fe₃C films at prior austenite grain boundaries + twinned martensite boundaries
- Impact energy drops dramatically (Charpy)
- Avoided by: NOT tempering in this range; or tempering above 350°C

### Alloy Temper Embrittlement (ATE) — "Reversible temper embrittlement"
- Occurs in alloy steels (Cr-Ni, Mn, Cr-Mo steels) when tempered OR slowly cooled through 350–550°C
- Caused by P, Sb, Sn, As segregation to prior austenite grain boundaries
- Reversible: can be "de-embrittled" by re-heating above 600°C and fast cooling
- Assessed by: shift in DBTT (ΔT_Charpy ≥ 30°C = embrittled)
- Prevented by: minimizing P/Sb/Sn/As in composition; adding Mo (0.2–0.3%Mo counteracts P embrittlement); fast cooling through susceptibility range

**Practical consequence:**
Steam turbine rotors (Cr-Mo-V steel) operating at 550°C can be susceptible to ATE over years. Regular mechanical property verification required.

---

## 5. Double Tempering and Secondary Hardening

**Double tempering:** Some high-alloy tool steels and high-speed steels are tempered twice:
1. First temper: temper fresh martensite; retained austenite → fresh martensite
2. Second temper: temper the newly formed martensite from step 1

Without double tempering, the fresh martensite from step 1 remains brittle.

**Secondary hardening (e.g., high-speed steel M2, H13 tool steel):**
Alloy carbides (Mo₂C, W₂C, V₄C₃, Cr₇C₃) precipitate during high-temperature tempering (500–600°C):

```
Hardness vs. tempering temperature (H13 tool steel):

Hardness
  │ 
65│────── fresh martensite
  │      ↘
55│       ─────────── hardness drops then:
  │                ╱
50│               ╱  ← secondary hardness peak (500–600°C)
  │              ╱     from alloy carbide precipitation
  │             ╱
40│────────────╱─────────────────  continues dropping
  └───────────────────────────── Temperature →
       200   400   600   700
```

Secondary hardening peak at ~540°C for M2 (H67-68 maintained after 550°C) → tool steels can be used at 600°C while maintaining cutting hardness.

---

## 6. Case Hardening Overview

**Purpose:** Create a hard surface layer on a tough core for optimal performance in wear + impact applications.

**Applications:**
- Gears: tooth surface hardness 58–62 HRC; core toughness for shock loading
- Camshafts: lobe surface hardness; core toughness
- Crankshaft journals: surface wear resistance
- Roller bearing races: case to resist Hertzian contact fatigue

**Choosing the right method:**

| Method | Temp (°C) | Case Depth | Case Hardness | Core Change | Distortion |
|--------|-----------|------------|---------------|-------------|------------|
| Carburizing | 900–940 | 0.5–4 mm | 58–65 HRC | Yes (must quench) | Moderate |
| Nitriding | 480–540 | 0.1–0.6 mm | 65–72 HRC | No | Very low |
| Carbonitriding | 850–875 | 0.1–0.75 mm | 58–65 HRC | Yes | Low |
| Induction hardening | Surface only | 1–10 mm | Depends on C | No | Low |
| Flame hardening | Surface only | 1–6 mm | Depends on C | No | Moderate |

---

## 7. Carburizing

**Principle:** Introduce carbon into the surface of a LOW-carbon steel base material by exposing the hot (austenite) surface to a carbon-rich atmosphere. The core retains its original low-carbon composition and transforms to tough martensite/bainite on quenching.

**Common methods:**

**Gas carburizing** (most common in industry):
```
Steel in furnace + endothermic gas (CO + H₂) + enrichment gas (CH₄ or C₃H₈)
Temperature: 900–940°C (in austenite range)
Carbon potential: controlled by O₂ probe (target 0.8–1.0% C at surface)
Time: 2–12 hours depending on target case depth

Carbon diffuses inward according to Fick's Second Law:
∂C/∂t = D_C × ∂²C/∂x²

Diffusion distance ≈ √(D_C × t)  where D_C ≈ 10⁻¹¹ m²/s at 930°C
```

**Vacuum carburizing (low pressure):**
```
Propane or acetylene bursts in vacuum furnace → very pure, accurate C control
Better case uniformity; less intergranular oxidation than gas carburizing
Common in aerospace gears (F-22 gearbox, helicopter transmission)
```

**After carburizing:**
```
1. Direct quench (from carburizing temperature) → fastest, but more distortion
2. Slow cool + reheat + quench → better microstructure control, less distortion
3. Oil quench → martensite in case (high C) + martensite/bainite in core (low C)
4. Temper at 160–180°C → reduce brittleness; maintain >58 HRC
```

**Carburized steel carbon profile:**
```
%C
1.0│─────── surface (0.8–1.0%)
   │       ╲
0.4│        ╲─────── effective case depth (0.4% C boundary)
   │                ╲
0.2│                 ╲──────────────── core carbon
   └────────────────────────── depth (mm) →
        0    0.5    1.0    2.0    3.0
```

---

## 8. Nitriding

**Principle:** Introduce nitrogen into steel surface at relatively low temperature (480–540°C) — below the critical temperature, so NO austenitizing and NO quench required. Core properties are UNCHANGED.

**Why nitrogen hardens:**
Nitrogen forms fine nitride precipitates (Fe₄N, Fe₂N for plain steels; CrN, AlN for alloy steels) at dislocations → precipitation hardening in the surface layer. Also solid solution hardening of the surface.

**Ammonia gas nitriding (Nitralloy process):**
```
Atmosphere: NH₃ gas → dissociates at surface: NH₃ → [N] + H₂
Temperature: 500–525°C for 20–100 hours
Case depth: 0.1–0.6 mm (very slow — N diffusivity < C diffusivity at same T)
Surface hardness: 900–1,100 HV (65–70 HRC equivalent!)
```

**Nitridable steels:** Must contain Cr, Al, Mo, V (form stable alloy nitrides → higher hardness):
- Nitralloy 135M (Al-Cr-Mo steel): highest nitriding response
- H13 tool steel: good nitriding response
- 4340: moderate response

**Advantages:**
- No quench → minimal distortion → ideal for precision gears, crankshafts, gauge plates
- Hard white layer + diffusion zone (some applications machine off white layer)
- Excellent fatigue resistance (compressive residual stress from volume expansion during nitride formation)

**Plasma nitriding:** Glow-discharge plasma at 350–600°C → faster, better control, shorter times, less distortion. Used for stainless steels and difficult alloys.

---

## 9. Carbonitriding and Ferritic Nitrocarburizing

**Carbonitriding:** Introduce both C and N simultaneously at 850–875°C (austenite range) + oil quench. Hybrid process: some benefits of both.
- Thinner cases than carburizing
- Lower distortion than carburizing (lower temperature)
- Better hardenability (N reduces M_s → more martensite during quench)
- Used for: small gears, bolts, fasteners

**Ferritic Nitrocarburizing (FNC):** Introduce N and small amount C at 570–590°C (below A₁, in ferrite). Very popular recently:
- Very thin hard layer (10–20 μm compound layer + diffusion zone)
- NO quench required → minimal distortion
- Excellent wear and corrosion resistance
- Low cost (short cycle, gas or salt bath)
- Used for: cylinder bores, piston rings, stamping dies, crankshafts

---

## 10. Induction Hardening and Flame Hardening

Both methods heat only the SURFACE using external heat sources, then quench — without altering core composition.

### Induction Hardening

**Principle:** Eddy currents induced by alternating magnetic field → joule heating → surface heats to austenite range → spray quench → martensite.

```
Induction hardening setup:

     [AC power supply] → copper induction coil
                              ↓ alternating magnetic field
                    ════════════════ steel part
                              ↓ induced eddy currents → heat
                         spray quench applied
```

**Frequency determines heating depth:**
- High frequency (10–500 kHz): shallow case (0.5–3 mm) → gear teeth, bearing journals
- Medium frequency (3–10 kHz): deeper case (3–8 mm) → crankshaft journals
- Low frequency (60–3,000 Hz): very deep case → very large rolls

**Advantages:**
- Very fast (seconds per part)
- Localized hardening (gear tooth tips vs. roots can be controlled)
- Low distortion
- Easy to automate
- No atmosphere required
- Only hardens where needed → fuel economy in automotive

**Limitation:** Base steel must have C > 0.3% (otherwise insufficient martensite hardness).

### Flame Hardening

Similar to induction but uses oxyacetylene or oxygas flame:
- Simpler equipment, flexible shapes
- Less precise than induction
- Used for: large surfaces, gear teeth, guide ways, machine tool slides
- Typical case depth: 1–6 mm

---

## Summary

| Process | Mechanism | Case/Depth | Hardness | Distortion |
|---------|-----------|-----------|---------|------------|
| Tempering (200°C) | Carbide precipitation | Through | -5 HRC | None |
| Tempering (600°C) | Full recovery | Through | -25 HRC | None |
| Gas carburizing | C diffusion at 930°C | 0.5–4 mm | 58–65 HRC | Moderate |
| Gas nitriding | N diffusion at 520°C | 0.1–0.5 mm | 65–70 HRC | Very low |
| Carbonitriding | C+N at 860°C | 0.1–0.7 mm | 60–65 HRC | Low |
| Induction hardening | Austenitize + quench | 0.5–8 mm | 55–65 HRC | Low |

---

## Exercises

1. A 4340 steel shaft (0.4%C, critical gear) requires surface hardness 58 HRC, core hardness 35 HRC. Compare: (a) through hardening + temper vs. (b) induction hardening. For each: state the process parameters, expected hardness profile, distortion risk, and whether the core property requirement is achievable. Which process would you recommend and why?

2. A nitrided crankshaft shows the following hardness profile: surface 950 HV, at 0.3mm depth 600 HV, at 0.6mm 300 HV (= core). The specification requires effective case depth > 0.4mm (at 550 HV). Does this part meet specification? The nitriding was at 520°C for 30 hours. Using an Arrhenius relation, estimate the additional hours needed at 520°C to reach 0.5mm case depth (assume depth ∝ √t from Fick's Law).

3. A carburized gear has surface carbon = 0.9%C, case depth (at 0.35%C) = 1.5mm. After quench and temper at 160°C: calculate expected surface hardness (using Fig: 0.9%C martensite tempered at 160°C ≈ 64 HRC). A competitive process is carbonitriding to 0.7mm depth. Compare the two for a gear application: (a) cyclic contact fatigue life (deeper case = better), (b) core strength (same base steel), (c) total cycle time (carburize at 920°C for 8h vs. carbonitriding at 860°C for 2h), (d) dimensional distortion. Which process wins overall?

4. Explain how secondary hardening works in H13 tool steel (0.4%C-5%Cr-1.3%Mo-1%V): (a) What phases form during tempering at 560°C? (b) Why does hardness INCREASE at ~560°C instead of continuing to decrease? (c) After secondary hardening, the steel is used at 600°C. Does it soften in service? Why or why not (think about thermodynamic stability of alloy carbides vs. Fe₃C).

5. A company uses gas carburizing for transmission gears and is considering switching to vacuum carburizing. For each process parameter: (a) carbon potential control (compare atmosphere management), (b) intergranular oxidation at surface (explain why vacuum process avoids this problem), (c) distortion during heating (vacuum vs. batch furnace with fixtures), (d) furnace cost and per-part cost. Write a recommendation memo summarizing the trade-offs for a high-volume automotive manufacturer vs. a low-volume aerospace company.
