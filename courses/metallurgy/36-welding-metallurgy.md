# Chapter 36: Welding Metallurgy — Joining by Fusion

> **"Welding is casting in miniature, done at 1,500°C, with the rest of the component as your mold, using nothing but an electric arc. The weld pool solidifies in milliseconds; the heat-affected zone goes through every heat treatment phenomenon in seconds. Understanding welding metallurgy means understanding all of the microstructure transformation science from the previous chapters — compressed into a few centimeters and a few seconds."**

---

## Table of Contents

1. [Welding Processes Overview](#1-welding-processes-overview)
2. [Weld Pool Solidification](#2-weld-pool-solidification)
3. [Heat-Affected Zone (HAZ) Microstructure](#3-heat-affected-zone-haz-microstructure)
4. [Solidification Cracking (Hot Cracking)](#4-solidification-cracking-hot-cracking)
5. [Hydrogen-Induced Cracking (Cold Cracking)](#5-hydrogen-induced-cracking-cold-cracking)
6. [Weldability — Concept and Carbon Equivalent](#6-weldability--concept-and-carbon-equivalent)
7. [Welding of Specific Alloys](#7-welding-of-specific-alloys)
8. [Post-Weld Heat Treatment (PWHT)](#8-post-weld-heat-treatment-pwht)
9. [Weld Defects and Inspection](#9-weld-defects-and-inspection)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Welding Processes Overview

All fusion welding processes create a weld pool (molten zone) surrounded by a solid heat-affected zone (HAZ).

**Key welding processes:**

| Process | Heat Source | Applications |
|---------|------------|-------------|
| SMAW (Stick) | Electric arc + coated electrode | General fabrication, repair |
| GMAW (MIG) | Arc + wire electrode (gas shielded) | Automotive, fabrication |
| GTAW (TIG) | Non-consumable W electrode (inert gas) | Al, Ti, stainless, thin sections |
| SAW | Arc submerged under flux | Heavy structural, pressure vessels |
| FCAW | Flux-cored wire | Outdoor construction, wind towers |
| EBW (Electron Beam) | Focused electron beam (vacuum) | Aerospace, precision, deep narrow welds |
| LBW (Laser) | Laser beam | Automotive (tailored blanks), precision |
| Friction Stir (FSW) | Plastic deformation (no fusion!) | Al alloys, 2xxx/7xxx (no solidification crack) |

**Key zones in a weld cross-section:**
```
  Fusion zone (FZ) = weld metal: completely melted + solidified
  ↕  (solidus line)
  Partially Melted Zone (PMZ): some grain boundary liquid
  ↕
  HAZ (Heat-Affected Zone): not melted, but microstructure changed by heat
  ↕  (temperature below which no change occurs)
  Base metal: unchanged
```

---

## 2. Weld Pool Solidification

The weld pool solidifies by **epitaxial solidification** — grains grow from existing base metal grains at the fusion line:
- No nucleation required (substrate grains continue to grow)
- Solidification direction: perpendicular to the fusion line (maximum temperature gradient)
- Produces **columnar grains** growing toward the weld centerline

**Solidification modes in welds:**
The ratio G/R (thermal gradient / growth rate) determines morphology — same as in casting (Ch 08):
- High G/R at fusion line → planar solidification
- Lower G/R toward centerline → columnar dendritic (most common)
- Very low G/R at centerline → equiaxed (can be induced by grain refiner additives)

**Segregation in weld metal:**
- Same Scheil equation as casting: C_s = k×C_0×(1-f_s)^(k-1)
- Low-k elements (Nb, S, P) segregate to grain boundaries + interdendritic regions
- Cause: hot cracking (see §4)
- Filler metal composition chosen to minimize low-m_p elements

**Multi-pass welding:**
- Each weld pass reheats underlying passes
- Previous weld metal recrystallizes or transforms → refined microstructure
- Martensite in HAZ of previous pass can be tempered by subsequent pass

---

## 3. Heat-Affected Zone (HAZ) Microstructure

The HAZ spans the region from the fusion line to where peak temperature caused no microstructural change. Different sub-zones form:

**For low-alloy steel (e.g., 4140) welded:**
```
Distance from fusion line:

← FUSION ZONE │ COARSE GRAIN HAZ │ FINE GRAIN HAZ │ INTERCRITICAL HAZ │ SUBCRITICAL │ BASE METAL →

Temp reached:
> T_liq  │ > T_liq → 1,100°C │ ~1,000°C     │ A_c1 to A_c3    │ ~500–700°C  │ RT
         │ Grain growth:      │ Austenite     │ partial re-      │ Tempered    │ 
         │ austenite grain    │ forms near    │ austenitize       │ carbides    │ 
         │ 100–500 μm        │ grain bdry    │                   │ recover     │
```

**Coarse-Grain HAZ (CGHAZ):**
- Peak T: 1,100–1,350°C (below T_liquidus)
- Austenite grain growth: d >> 100 μm (prior austenite grain size control lost)
- On cooling: large martensite packets → brittle
- Toughness: worst of any HAZ zone
- Risk: HAZ cold cracking + brittle fracture in service

**Fine-Grain HAZ (FGHAZ):**
- Peak T: ~950–1,100°C (just above Ac3)
- Re-austenitize + recrystallize → fine grain
- Often the toughest zone in the HAZ

**Intercritical HAZ:**
- Peak T: Ac1–Ac3 (partial austenitization)
- Mixed ferrite + austenite → on cooling: mixed ferrite + hard islands
- Potential for MA (martensite/austenite) constituents → embrittlement in some conditions

---

## 4. Solidification Cracking (Hot Cracking)

**Hot cracking** occurs when tensile stresses (from contraction on solidification) exceed the strength of the mushy zone (semi-solid region) — same mechanism as in castings (Ch 08), but in welding it occurs in seconds.

**Conditions:**
1. Wide solidification range (large T_liq - T_sol): more time in brittle temperature range (BTR)
2. Low-melting-point films at grain boundaries (Ni-S, Ni-P, Al-Si eutectic)
3. Tensile stress during solidification (from weld contraction)

**Sensitivity:** Solute content of S, P, B, Si:
- These elements have very low k (partition coefficient) → segregate to solidification front → form low-mp liquid films
- Once films wet grain boundaries → grain boundary strength → 0 → crack

**Weld cracking indices:**
For stainless steel, Schaffler diagram / Creq/Nieq ratio:
- Weld metal with > 5–8% ferrite: less susceptible to hot cracking (ferrite dissolves S, P)
- Pure austenite: highly susceptible

**Prevention:**
- Keep S + P < 0.010% in base metal and filler
- Use fillers with higher Mn/S ratio (Mn forms MnS — higher mp than FeS)
- Minimize heat input (faster solidification)
- Reduce restraint (reduces tensile stress during solidification)
- Choose welding parameters: narrower weld bead → less segregation

---

## 5. Hydrogen-Induced Cracking (Cold Cracking)

**Cold cracking** (also: hydrogen embrittlement, underbead cracking, delayed cracking) occurs AFTER the weld has cooled — sometimes days later.

**Requirements (all three must be present):**
1. Hydrogen: from moisture in electrodes, flux, base metal surface
2. Susceptible microstructure: martensite (>~0.30% C equivalent)
3. Tensile stress: from restraint + thermal shrinkage

**Mechanism:**
```
H₂ diffuses into weld metal and HAZ at high T → dissolved as H atoms
On cooling → martensite forms (hard, low diffusivity for H)
H accumulates at stress concentrations (pores, inclusions, boundaries)
H reduces local fracture toughness → crack initiates + propagates
```

**Why delayed:** Hydrogen diffusion is slow at room temperature → may take hours to days to accumulate at critical location.

**Hydrogen levels:**
| Electrode type | Diffusible H (mL/100g weld metal) | Cold crack risk |
|---|---|---|
| Non-low-hydrogen | >15 | High |
| Basic low-H (E7018) | <5 | Lower |
| Very-low H (dried) | <2 | Low |

**Prevention:**
- Low-hydrogen electrodes (baked at 300–400°C before use)
- Preheat (30–200°C) → H diffuses out before martensite forms, slower cooling → less martensite
- Post-weld hydrogen bake-out (200–300°C for 1–4 hours)
- Choose base metal with low CE (carbon equivalent, see §6)
- Interpass temperature control (keep between passes warm)

---

## 6. Weldability — Concept and Carbon Equivalent

**Weldability:** Ability of a material to be welded under specific conditions to meet required properties without defects.

**Carbon Equivalent (CE):**
Predicts susceptibility to HAZ hardening (martensite) → cold cracking:
```
CE = %C + %Mn/6 + (%Cr+%Mo+%V)/5 + (%Ni+%Cu)/15

(IIW formula — International Institute of Welding)
```

| CE value | HAZ Hardness | Cold crack risk | Preheat needed |
|---------|-------------|----------------|---------------|
| < 0.35 | < 350 HV | Very low | None |
| 0.35–0.45 | 350–450 HV | Low | 50–100°C |
| 0.45–0.60 | 450–550 HV | Medium | 100–200°C |
| > 0.60 | > 550 HV | High | >200°C + PWHT |

**Simplified cold cracking susceptibility index (Pcm):**
```
Pcm = %C + %Si/30 + (%Mn+%Cu+%Cr)/20 + %Ni/60 + %Mo/15 + %V/10 + 5B
```
More sensitive for low-C steels.

**Effect of preheat:**
- Reduces thermal gradient → slower cooling → less martensite
- Allows H to diffuse out before martensite forms
- T_preheat from: T_preheat = 350 × √(CE - 0.25) - 25 (°C), approximate formula

---

## 7. Welding of Specific Alloys

### Austenitic Stainless Steel (304, 316)
- Susceptible to: (a) hot cracking (austenite), (b) sensitization in HAZ (Ch 30)
- Solution: Use 308L/316L filler (low C) → < 0.03%C → no sensitization at the very narrow sensitization zone
- Maintain < 5% δ-ferrite in weld: avoids hot cracking (ferrite absorbs S, P)
- Heat input control: lower T avoids sensitization; must not keep weld joint at 450–850°C too long

### Duplex Stainless Steel (2205)
- HAZ can have excess ferrite (loss of austenite → loss of toughness, SCC resistance)
- Preheat: NOT recommended (promotes σ-phase precipitation)
- Fast cooling: good; high heat input: allows time for austenite to re-form

### Aluminum Alloys
- TIG or MIG with matching/similar filler
- 2xxx (Al-Cu): hot crack susceptible → FSW preferred (avoids solidification)
- 6xxx: HAZ softening (overaging in HAZ) → strength loss ~40%
- 7xxx (Al-Zn): FSW preferred; fusion welding possible only with specific fillers (4043 alloy - Al-Si)

### Titanium Alloys
- Must use inert gas shielding (Ar) from all sides: TIG welding in glove box for Ti
- HAZ O₂ contamination → α-case → brittle
- No preheat; maintain clean surfaces (no grease, water)
- Ti-6Al-4V: can be TIG welded; post-weld anneal to relieve residual stress

### Nickel Superalloys
- Higher-strength superalloys (> ~30% γ′): effectively UNWELDABLE by fusion methods
  - Strain-age cracking: on PWHT (aging), γ′ precipitates form while HAZ is still stressed → crack
  - Low-γ′ alloys (IN625, IN625, IN617): weldable
  - IN718: Nb-bearing → sluggish γ″/δ precipitation → better weldability than others
- TBC coatings and disk alloys: repaired by brazing or diffusion bonding instead of welding

---

## 8. Post-Weld Heat Treatment (PWHT)

**PWHT** is performed to:
- Relieve residual stresses (stress relief anneal)
- Temper martensite in HAZ (for carbon steels)
- Restore toughness
- Reduce distortion risk in service

**Stress relief for carbon steels:**
- 550–650°C for 1 hour per 25 mm thickness
- Reduces residual stress by ~80%
- Softens hard HAZ martensite

**PWHT for pressure vessel code (ASME Sec. VIII):**
- Required when wall thickness > 38 mm OR CE > threshold OR service environment requires it (H₂, H₂S service)
- Specific T and time per material P-number classification

**PWHT for sensitization risk (stainless):**
- Solution anneal at 1,050–1,100°C + water quench: dissolves Cr carbides → prevents sensitization
- Only if entire component can be heat treated; otherwise: use L-grade (low C) filler

---

## 9. Weld Defects and Inspection

**Common weld defects:**

| Defect | Cause | Detection | Prevention |
|--------|-------|-----------|-----------|
| Porosity | Gas (H₂, CO) trapped in weld pool | X-ray, UT | Dry consumables, clean base metal |
| Hot cracks | Segregation at grain boundaries | Dye-PT, X-ray | Low S+P, proper filler, reduce restraint |
| Cold cracks | H in martensite + stress | Magnetic PT (UT if embedded) | Low-H electrodes, preheat |
| Lack of fusion | Incomplete melting | UT, X-ray | Proper heat input, correct torch angle |
| Undercut | Excessive arc: melts into base | Visual | Correct parameters |
| Distortion | Thermal contraction | Dimensional check | Clamping, sequence, pre-setting |
| Slag inclusions | Entrapped flux/slag | X-ray | Interpass cleaning |
| Overlap | Metal flowing over cold base | Visual | Correct parameters |

**NDT for welds:**
- Visual inspection: surface defects (always first)
- Dye penetrant (PT): surface-breaking cracks (non-magnetic materials)
- Magnetic particle (MT): surface + near-surface cracks (ferromagnetic materials only)
- Radiographic testing (RT): volumetric (porosity, inclusions, cracks)
- Ultrasonic testing (UT): best for embedded cracks; phased array UT → automated scanning

---

## Summary

| Weld Zone/Issue | Key Mechanism | Critical Variable | Prevention/Control |
|----------------|--------------|-------------------|-------------------|
| CGHAZ | Grain growth in austenite at high T | Peak temperature + time | Minimize heat input; low-CE base metal |
| Hot cracking | Low-mp boundary films in mushy zone | S, P, Si content | Low-S filler; maintain >5% ferrite in stainless |
| Cold cracking | H in martensite + residual stress | CE, H content | Preheat, low-H electrode, PWHT |
| Sensitization | Cr₂₃C₆ precipitation at GB | 450–850°C exposure time | Low-C filler; quick heat input pass |
| FSW (Al) | Solid state — no solidification | Tool design, speed | Best route for 2xxx/7xxx Al |

---

## Exercises

1. A 25 mm plate of AISI 4140 steel (0.40%C, 0.90%Mn, 1.0%Cr, 0.20%Mo) is being welded with SMAW using E7018 electrodes. (a) Calculate CE (IIW formula). (b) Is preheat required? What temperature? (c) The CGHAZ cools from 1,200°C to 300°C in 15 seconds. Using Figure 19.X (schematic TTT): does this cooling rate produce martensite, bainite, or pearlite in the CGHAZ? (d) Without preheat, a cold crack forms 6 hours after welding. List the three conditions that allowed cracking. (e) Specify the complete welding procedure (preheat, electrode, interpass T, PWHT) to prevent recurrence.

2. Austenitic stainless steel 316 plate is welded using 316L filler (0.02%C). A sensitization test shows the HAZ adjacent to the fusion line is immune to intergranular corrosion, but the region 3–5 mm away failed the sensitization test. (a) Explain why the zone nearest the fusion line is NOT sensitized, while the zone 3–5 mm away IS sensitized, using a time-temperature schematic. (b) Calculate the approximate carbon content at which Cr₂₃C₆ precipitation can be prevented at 700°C for 10 minutes (use: 12%Cr minimum for protection; assume all C ties up as Cr₂₃C₆ with Cr:C ratio = 4:1 by weight). (c) Why does using 316L (0.02%C max) prevent sensitization in the narrow sensitized zone, but NOT if the weld is multi-pass and the zone spends cumulatively 30 minutes at 700°C?

3. Friction Stir Welding of 7075-T6 aluminum plate: (a) Unlike fusion welding, FSW produces NO solidification. What defect mechanisms are therefore absent? (b) The FSW tool produces a nugget (recrystallized) zone + TMAZ (thermomechanically affected) + HAZ. In the HAZ, the peak temperature reaches ~250°C for ~5 minutes. 7075 naturally overages at 180°C. What precipitation change occurs in the HAZ? What property is lost? (c) 7075-T73 (overaged) has better SCC resistance than T6. After FSW (which overages the HAZ through thermal exposure), is the HAZ of a 7075-T6 weld more or less SCC-resistant than the base metal? Explain. (d) Can the base metal T6 strength be restored by PWHT after FSW? What concerns exist?

4. Hot cracking susceptibility: Two Ni-based alloys are being compared for weld repair: (a) IN625 (0.10%C, 22%Cr, 9%Mo, 3.6%Nb) — relatively low γ′ content. (b) Mar-M-247 (0.15%C, 8.25%Cr, 10%W, 5.5%Al, 3%Ti) — ~65% γ′. Why is IN625 commonly used as a filler for weld repair while Mar-M-247 cannot be fusion welded? Consider: (i) γ′ fraction effect on strain-age cracking, (ii) solidification range of each alloy (IN625: narrow; Mar-M-247: wide due to W, Ta segregation), (iii) residual stress state after cooling.

5. In a hydrogen-induced cracking test, four 20 mm thick steel plates (CE = 0.50) are welded with: (a) E6010 (non-low-H, H = 18 mL/100g); (b) E7018 as-received (H ≈ 8 mL/100g); (c) E7018 dried at 350°C (H = 2 mL/100g); (d) E7018 dried + preheat to 150°C. Rank the probability of cold cracking (1 = highest risk, 4 = lowest). For cases (a) and (d), calculate the CGHAZ hardness expected (use empirical formula HV_CGHAZ ≈ 90 + 1,050 × CE_pcm where Pcm = CE/1.6 approximately) and state whether this exceeds 350 HV (critical threshold for H cracking).
