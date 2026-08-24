# Chapter 19: Hardening and Quenching — Making Steel Hard

> **"The Japanese swordsmith heats the blade cherry-red, pulls it from the fire, plunges it into water. In that moment of transformation — the crackling, the steam, the color draining from the blade — centuries of empirical knowledge are at work. The blade is harder than before, but also more brittle. The next step, tempering, is the art of finding the balance: hard enough to hold an edge, tough enough not to shatter in battle."**

---

## Table of Contents

1. [The Goal — Converting Austenite to Martensite](#1-the-goal--converting-austenite-to-martensite)
2. [Austenitizing — Preparing for Quench](#2-austenitizing--preparing-for-quench)
3. [Martensite — The Hard Phase](#3-martensite--the-hard-phase)
4. [Quench Media and Cooling Rates](#4-quench-media-and-cooling-rates)
5. [Quenching Problems — Distortion and Cracking](#5-quenching-problems--distortion-and-cracking)
6. [Hardenability and the Jominy Test](#6-hardenability-and-the-jominy-test)
7. [Retained Austenite](#7-retained-austenite)
8. [Marquenching and Austempering](#8-marquenching-and-austempering)
9. [Surface vs. Through Hardening](#9-surface-vs-through-hardening)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Goal — Converting Austenite to Martensite

The hardening process aims to convert the FCC austenite phase into martensite — the hardest form of steel at a given carbon content. The sequence is:
```
Heat → Austenitize → Quench (fast cool) → Martensite → HARD STEEL
```

Martensite forms because:
- Quenching is too fast for carbon to diffuse (no time to form pearlite or bainite)
- Austenite → martensite is a **diffusionless** transformation (atoms move < one interatomic spacing)
- Carbon is "trapped" in the BCC lattice → lattice distorts to BCT (body-centered tetragonal) → very high internal stress → very high hardness

---

## 2. Austenitizing — Preparing for Quench

**Requirements for a good austenite:**
1. **Complete austenitization:** All carbides dissolved, all carbon in solution
2. **Uniform composition:** Avoid composition gradients (microsegregation)
3. **Correct grain size:** Fine austenite grain (ASTM 6–10) → fine martensite → better toughness

**Temperature selection:**
```
Plain carbon steels:
  Hypoeutectoid (<0.8%C): A₃ + 30°C    → 830–880°C
  Eutectoid (0.8%C):      A₁ + 50°C    → 780°C
  Hypereutectoid (>0.8%C): A₁ + 50°C   → 760–780°C 
    (heat to austenite + undissolved carbide → prevents austenite grain growth
     while retaining some carbide for wear resistance)

Alloy steels:
  4340, H13 tool steel: 830–950°C (higher to dissolve alloy carbides)
```

**Soaking time:** Approximately 1 hour per 25mm of thickness (rule of thumb). Time for:
- Furnace equalization
- Carbide dissolution
- Composition homogenization

**Atmosphere:** Protective (endogas, vacuum, salt bath) to prevent decarburization (loss of surface carbon → soft surface = catastrophic for bearing races, cutting tools).

---

## 3. Martensite — The Hard Phase

### Crystal Structure

Martensite is a supersaturated interstitial solid solution of carbon in BCC iron. The trapped carbon distorts the BCC lattice → BCT (Body-Centered Tetragonal):

```
BCC (pure Fe) → BCT (martensite with C):

     ○
    /│\
   / │ \
  ○──○──○  →  c/a > 1 (elongated in one direction)
  │  ●  │     ● = interstitial C atom in center of elongated cell
  ○──○──○
  
c/a ≈ 1 + 0.046 × [%C]   (tetragonality increases with C content)
```

### Martensite Hardness

Martensite hardness depends almost entirely on carbon content:
```
Hardness (HRC) vs. % C:
  0.1% C: ~30 HRC (320 HV)
  0.2% C: ~40 HRC (380 HV)
  0.4% C: ~55 HRC (570 HV)
  0.6% C: ~63 HRC (750 HV)
  0.8% C: ~65 HRC (800 HV)
  1.0% C: ~66 HRC (>800 HV)
  (above 0.6% C, retained austenite reduces apparent hardness)
```

### Why Is Martensite So Hard?

1. **Extreme solid solution hardening:** Interstitial C in BCT lattice creates massive lattice distortion
2. **High dislocation density:** The shear transformation (diffusionless) generates enormous internal stress → many dislocations
3. **Fine microstructure:** Martensite plates/laths are extremely small → many internal boundaries

### Martensite Morphology

| Type | Carbon Content | Appearance | Typical Steel |
|------|---------------|------------|---------------|
| Lath martensite | < 0.5% C | Parallel laths in packets | Low/medium carbon |
| Mixed | 0.5–0.8% C | Mixed morphology | Medium-high carbon |
| Plate martensite | > 0.8% C | Lens-shaped plates, midrib | High carbon |

Lath martensite is tougher (more ductile) than plate martensite.

---

## 4. Quench Media and Cooling Rates

The quench medium determines the cooling rate. The critical requirement: cool FAST enough to miss the pearlite nose (and usually the bainite nose too) to form martensite.

| Quench Medium | Relative Cooling Severity | Typical Use |
|---------------|--------------------------|-------------|
| Brine (10% NaCl + water) | Very fast (~2×water) | Deep hardening small parts, maximum speed |
| Water (25°C) | Fast (severity H ≈ 1.0) | High-hardenability alloy steels, some plain carbon |
| Water (65°C warm) | Fast-moderate | Alloy steels, reduces distortion |
| Fast quench oil | Moderate (H ≈ 0.3–0.4) | Most alloy steels (4340, etc.) |
| Conventional oil (60°C) | Moderate (H ≈ 0.25–0.35) | Alloy steels, less distortion |
| Polymer quench (Aqua-Quench) | Variable (0.5×–2×water) | Controlled, programmable, water-like but less distortion |
| Air (still) | Slow (H ≈ 0.02) | Very high-hardenability steels, tool steels |
| Pressurized gas (N₂) | Moderate-slow | Vacuum furnaces, tool steels, superalloys |

**Quenching severity H:** Defined such that H = ∞ is an ideal quench (surface at 0°C instantly).

---

## 5. Quenching Problems — Distortion and Cracking

Rapid quenching creates THERMAL GRADIENTS → differential contraction → thermal stresses + phase transformation stresses → distortion and cracking.

### Thermal Stress Mechanism

```
Surface cools first:
Stage 1: Surface contracts → tension on surface, compression in center
         (normal thermal contraction)
         
Stage 2: Surface transforms to martensite → EXPANDS (~4% volume increase!)
         Center still austenite
         → Surface expansion + center contraction creates:
         
Stage 3: Surface in COMPRESSION, center in TENSION after martensite transformation
         (opposite to thermal stress!)
         
Stage 4: Center transforms → expansion of center → surface put in TENSION
         → Risk of surface quench cracking
```

**Quench cracks:** Typically occur at:
- Sharp corners (stress concentrations)
- Thin sections adjacent to thick sections
- At the junction of carburized and uncarburized regions
- High-carbon steels (< 10°C thermal window before cracking)

**Distortion:** Non-uniform transformation and thermal contraction → parts bow, twist, or ovalize. Round bars become oval; flat plates curve.

**Minimizing quench cracking/distortion:**
- Use lowest severity quench that still achieves required hardness (oil over water if possible)
- Preheat before austenitizing (reduce ΔT on quench)
- Maintain uniform quench agitation
- Use marquenching (§8) for susceptible parts
- Correct part design: avoid sharp internal corners, uniform section thickness

---

## 6. Hardenability and the Jominy Test

**Hardenability:** The ability of a steel to form martensite throughout its cross-section. NOT the same as hardness (which is the actual hardness of martensite at a given C content).

A steel with HIGH hardenability:
- Forms martensite even at a SLOW cooling rate (deep in a thick section)
- TTT curve shifted far to the right (long times to start pearlite/bainite)

**Alloying elements that increase hardenability (in order of effectiveness):**
Mo > Mn = Cr > Ni > Si (per unit weight%)

These elements dissolve in austenite and slow the diffusion-controlled pearlite/bainite transformations.

**Jominy End-Quench Test (ASTM A255):**

```
1. Normalize test bar (25mm dia × 100mm long)
2. Austenitize at standard temperature
3. Spray cold water on ONE END only:
   
       WATER SPRAY
            ↓
   ─────────────────────────────
   │ hot end │────────────────→│ far end
   └─────────────────────────────┘
   1.6mm 12mm    25mm       100mm (distance from quench end)
   
4. Measure Rockwell hardness at every 1.6mm from quenched end
5. Plot hardness vs. distance
```

**Reading Jominy curves:**
- Near quenched end: fastest cooling → martensite → high hardness
- Far from quenched end: slow cooling → pearlite/bainite → lower hardness
- Steep drop: low hardenability (plain carbon steel)
- Flat curve: high hardenability (alloy steel)

---

## 7. Retained Austenite

Above ~0.5% C, the martensite finish temperature M_f drops below room temperature → some austenite is NOT transformed (retained) at RT:

```
% Retained Austenite ≈ equation varies, roughly:
  0.2% C: ~2% retained
  0.6% C: ~10% retained
  1.0% C: ~20–30% retained
  High-alloy tool steel: up to 40% retained
```

**Problems with retained austenite:**
- Soft spots in hardened steel
- Dimensional instability (transforms to martensite during service → volume change → distortion)
- CRITICAL for precision parts (bearings, gages, molds): must be <1%

**Elimination (sub-zero treatment):**
Cool to −70°C to −120°C within hours of quenching (before austenite stabilizes):
```
M_f ≈ M_s − 150°C for most steels
Sub-zero cool to −70°C → transforms most retained austenite
Liquid nitrogen (−196°C) for maximum conversion
```

Immediately followed by tempering to relieve brittleness.

---

## 8. Marquenching and Austempering

### Marquenching (Martempering)

**Purpose:** Reduce thermal gradient during martensite transformation → less distortion and cracking.

**Process:**
```
Austenitize
    ↓
Quench into molten salt bath or hot oil at JUST ABOVE M_s (e.g., 200°C)
    ↓
Hold until temperature equalizes throughout section
    ↓ (no transformation yet — still austenite)
Air cool slowly through martensite transformation range
    ↓ (slow, uniform martensite formation → much less distortion)
Temper
```

Result: Same martensite hardness, much less distortion. Used for complex gears, dies, precision parts.

### Austempering

**Purpose:** Form bainite directly. Bainite has good combination of hardness AND toughness (better than tempered martensite at same hardness level).

**Process:**
```
Austenitize
    ↓
Quench into salt bath held at BAINITE TEMPERATURE (260–400°C)
    ↓
Hold until transformation is COMPLETE (isothermal bainite)
    ↓
Air cool
    ↓
NO TEMPERING NEEDED (bainite is already stable)
```

**Advantages:**
- No distortion from martensite transformation
- Better ductility/toughness than tempered martensite at same hardness (HRC 45–55)
- No risk of quench cracking
- Used for: springs, small gears, clips, safety-critical springs (automotive)

---

## 9. Surface vs. Through Hardening

**Through hardening:** Entire cross-section austenitized and quenched. Used when hardness through section is required.

**Case hardening:** Only surface and subsurface regions hardened. Core remains soft and tough. Best combination for gears, camshafts: hard surface (wear) + soft tough core (impact resistance). Methods:
- Carburizing (Ch 20): add C to surface
- Nitriding (Ch 20): add N to surface
- Induction hardening: rapid surface heating by induction, then quench

```
Case-hardened microstructure:

Surface → Case: Martensite (hard, 60–65 HRC)
          Transition zone
          Core: original tough microstructure (25–30 HRC)
```

---

## Summary

| Step | Parameter | Goal |
|------|-----------|------|
| Austenitize | A₃ + 30°C, 1h/25mm | Fully dissolved austenite, fine grain |
| Quench | Fast enough to miss pearlite nose | Martensite formation |
| Hardness achieved | Depends on %C (not alloy content) | Up to 66 HRC for high C |
| Hardenability | Depends on alloy content (Mn, Cr, Mo) | Martensite through thick sections |
| Retained austenite | Sub-zero treat for high-C steels | Dimensional stability |
| Marquench | Hold above M_s before air cool | Reduce distortion |
| Austemper | Hold in bainite range | Bainite: good hardness + toughness |

---

## Exercises

1. A 0.6% C steel has M_s = 260°C, M_f = 110°C. (a) Calculate the approximate % retained austenite if quenched to room temperature (25°C). (b) If sub-zero treated to −70°C, estimate the retained austenite. (c) The part is a precision bearing race. What maximum retained austenite is acceptable, and what sub-zero temperature is required? (d) Why must sub-zero treatment be done within 1–2 hours of quenching (before austenite stabilization)?

2. Compare the Jominy hardenability curves for: (a) AISI 1040 plain carbon steel (steep drop, H=4 at J=12mm), (b) AISI 4140 alloy steel (H=4 at J=25mm), (c) AISI 4340 (H=4 at J=50mm). For a 40mm diameter shaft requiring 40 HRC at the center: which steels could be through-hardened using water quench? oil quench? Which steel would you choose for minimum quench distortion?

3. During quench cracking investigation, cracks are found at a sharp internal corner (radius 0.5 mm) in a high-carbon tool steel after oil quenching. List three changes to the manufacturing process that could prevent quench cracking: (a) one geometry change, (b) one heat treatment process change (without changing the quench medium), (c) one quench medium change. For each, explain the metallurgical mechanism by which it reduces quench crack risk.

4. Austempering is specified for a safety-critical automotive suspension spring (AISI 1080, 0.8%C). The TTT diagram shows bainite start at 320°C requires 10 seconds (nose) and finishes in 200 seconds. Design the austempering process: (a) austenitizing temperature and time (spring wire diameter = 4mm), (b) quench bath temperature, (c) required hold time in salt bath, (d) expected hardness (bainite at 320°C ≈ HRC 48). Compare to tempered martensite at the same hardness — which has better toughness and why?

5. A steel gear (module 3, face width 25mm) requires: surface hardness 58–62 HRC, core hardness 28–35 HRC. Design a complete case-hardening process including: (a) base steel selection (low-carbon grade, why), (b) carburizing to achieve surface 0.8%C, (c) quench to develop surface martensite + soft core, (d) temper to reduce brittleness while maintaining ≥58 HRC at surface. Sketch the expected hardness profile from surface to core.
