# Chapter 49: Single Crystal Growth — Seeds, Selectors, and the Art of Casting

> **"Growing a single-crystal turbine blade is one of the most demanding manufacturing operations in the world. You are controlling the atomic-scale solidification front of a 10-element alloy, inside a ceramic mold the shape of an airfoil, to grow a perfect crystal 15 cm long with no defects — and you must do it reproducibly, 30,000 times per year."**

---

## Table of Contents

1. [Two Methods: Grain Selector vs. Seed Crystal](#1-two-methods-grain-selector-vs-seed-crystal)
2. [The Bridgman Furnace — Hardware Overview](#2-the-bridgman-furnace--hardware-overview)
3. [The Grain Selector — Pigtail and Spiral Geometries](#3-the-grain-selector--pigtail-and-spiral-geometries)
4. [Seed Crystal Method — Controlling Orientation Precisely](#4-seed-crystal-method--controlling-orientation-precisely)
5. [Thermal Gradient and Withdrawal Rate — The Critical Parameters](#5-thermal-gradient-and-withdrawal-rate--the-critical-parameters)
6. [Dendrite Solidification and Microsegregation](#6-dendrite-solidification-and-microsegregation)
7. [The Investment Casting Mold with Ceramic Core](#7-the-investment-casting-mold-with-ceramic-core)
8. [Post-Solidification Processing](#8-post-solidification-processing)
9. [Orientation Measurement and Qualification](#9-orientation-measurement-and-qualification)
10. [Common Defects and Their Causes](#10-common-defects-and-their-causes)
11. [Industrial Scale and Process Control](#11-industrial-scale-and-process-control)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Two Methods: Grain Selector vs. Seed Crystal

Two fundamentally different methods exist to produce single-crystal turbine blades:

### Method 1 — Grain Selector (Most Common in Industry)

A **geometric grain selector** — typically a spiral or pigtail passage — is attached to the bottom of the investment casting mold. During solidification:
1. Many grains nucleate at the cold base
2. Only grains with favorable growth orientation for the selector geometry can navigate through
3. Typically, ONE grain wins the selection race
4. That grain grows into the full blade mold

**Advantages:**
- No pre-made seed crystal needed
- Self-selecting → simpler logistically
- Naturally selects near-[001] orientations (detailed in §3)

**Disadvantages:**
- Cannot control secondary orientation (rotation of crystal around primary axis)
- Statistical process → small probability of multiple grains winning selection → defects

### Method 2 — Seed Crystal (Research and Premium Applications)

A **pre-oriented seed crystal** of the alloy is placed at the base of the mold:
1. The seed crystal (a small piece of previously grown SX, carefully oriented)
2. The melt above is directionally solidified so it epitaxially continues the seed's orientation
3. The result is a crystal with controlled primary AND secondary orientation

**Advantages:**
- Full crystallographic control (primary and secondary orientations)
- No "lucky" grain selection — guaranteed orientation
- Reproducible from part to part

**Disadvantages:**
- Must maintain a library of qualified seed crystals
- Seed-blade interface can be a weak point if not handled correctly
- Higher process complexity

Most industrial blade production uses grain selectors. Research and some premium applications (especially where secondary orientation matters for fatigue) use seed crystals.

---

## 2. The Bridgman Furnace — Hardware Overview

The **Bridgman process** (also called the Bridgman-Stockbarger process) is the standard industrial method for directional solidification of superalloy blades:

```
Bridgman Furnace Schematic:

    ┌─────────────────────────────────────┐
    │  Resistance or induction heaters    │ ← Zone 1: Preheat (~1500°C)
    │  ≈ 1550°C (above alloy liquidus)   │
    │                                     │
    │  [Ceramic mold + metal inside]      │ ← Mold assembly on withdrawal rod
    │                                     │
    ├─────────────────────────────────────┤ ← Baffle plate (radiation shield)
    │  Chill zone (below baffle)          │
    │  Water-cooled base / chill plate    │
    └─────────────────────────────────────┘
    
    ↓ Withdrawal direction (mold moves DOWN at controlled rate)
    
Key components:
1. Heater zone: keeps metal above liquidus
2. Radiation baffle: creates steep thermal gradient at solidification front
3. Chill plate / withdrawal rod: pulls mold downward through baffle at 1–10 mm/min
4. Vacuum environment: protects reactive alloy from oxidation
```

**Vacuum requirement:** Superalloys contain reactive elements (Al, Ti, Ta) that oxidize instantly in air at 1500°C. The entire furnace operates under high vacuum (< 10⁻³ Pa) or protective inert atmosphere.

**Radiation baffle (susceptor):** The most critical element for SX quality. It creates an abrupt transition from hot (above) to cold (below). The sharper the gradient at the solidification front, the less likely stray grain nucleation, the more perfect the crystal.

---

## 3. The Grain Selector — Pigtail and Spiral Geometries

The grain selector is a narrow passage attached below the blade mold. As grains grow up from the base plate, only grains with a specific orientation can navigate through the geometry.

### Spiral Selector

```
    ┌──────────────────┐ ← Main blade cavity
    │                  │
    │    Blade mold    │
    │                  │
    └────────┬─────────┘
             │  ← Transition to selector
           ┌─┤
           │ │ ← Conical starter
           └─│────────────────╮  ← Spiral coil
             │   spiral turn  │     (3/4 to 1 full revolution)
             │                │
             └────────────────╯
                ↕ Chill plate
```

The spiral geometry selects grains based on **growth orientation**: grains that grow most rapidly along the centrifugal direction (the [001] direction in FCC superalloys grows fastest — lowest surface energy, preferred dendrite direction) navigate the spiral fastest. All other grains are geometrically blocked at the curved walls and stop growing.

### Why [001] Grows Fastest (Dendrite Orientation)

FCC metals grow preferentially in the ⟨001⟩ directions because these are the crystallographically "easy" directions for heat extraction along dendrite arms. In the spiral, the grain whose [001] direction most closely aligns with the spiral's axis will navigate farthest up the spiral. Over the spiral length of ~20–30 mm, competing grains are naturally eliminated.

The result: the grain emerging from the spiral top is typically within 10–15° of [001] — acceptable for most applications. With careful design, this can be narrowed to < 10°.

### Pigtail Selector

An alternative geometry: a sharp 180° turn followed by another turn. Only grains that can round the tight curve (those growing with [001] perpendicular to the curve — i.e., aligned with the main blade direction) survive. Similar selection mechanism.

---

## 4. Seed Crystal Method — Controlling Orientation Precisely

In the seed crystal method:

```
    ┌──────────────────┐
    │   Blade mold     │
    │    (empty)       │
    ├──────────────────┤
    │   Seed crystal   │ ← Small SX piece, precisely oriented
    │   (pre-oriented) │   [001] ± 5° from blade axis
    │                  │   known secondary orientation
    └──────────────────┘
          ↓ chill plate
```

**Process:**
1. The seed crystal is placed at the base, held in a ceramic holder
2. Mold is heated until ONLY the top of the seed melts slightly (partial remelt of seed)
3. The melted region creates a clean interface
4. As the mold is slowly withdrawn, the melt re-solidifies epitaxially from the partially-melted seed
5. The crystal grows up exactly continuing the seed's orientation

**Critical step — partial remelt:** You want to melt the TOP of the seed (1–2 mm) to create a fresh interface without contamination. But if you melt TOO much of the seed, you lose the orientation reference. Too little → oxide/contamination layer at the interface → grain nucleation.

Temperature control accuracy for seeded casting: ±1–2°C at the seed.

**Where seed casting is preferred:**
- Military blades where secondary orientation strongly affects fatigue life (e.g., trailing edge fatigue)
- Research and development of new alloys (need precise orientation for property measurement)
- Components where tight crystallographic tolerances are specified

---

## 5. Thermal Gradient and Withdrawal Rate — The Critical Parameters

Two parameters dominate SX casting quality:

### Thermal Gradient (G, °C/mm)

The temperature gradient at the solidification front:
```
G = (T_liquidus - T_solidus) / distance
```

Typical range: 30–100°C/mm in industrial Bridgman furnaces.

**Why high G is critical:**
- High G → steep temperature rise ahead of the solidification front → any ahead-of-front nucleus quickly melts back (too hot)
- Low G → shallow temperature gradient → ahead-of-front region is close to liquidus → constitutional supercooling → stray grain nucleation ahead of the front

For SX, the region ahead of the solidification front must be **superheated** (above liquidus temperature) at all times. Any undercooled region can nucleate a stray grain.

Modern Bridgman furnaces achieve G up to 50–70°C/mm. **Liquid metal cooling (LMC)** technology (using a molten Sn or Al bath as the cooling medium instead of radiation) achieves G up to 100–150°C/mm → significantly better SX quality → fewer stray grains → higher casting yield.

### Withdrawal Rate (V, mm/min)

The speed at which the mold is pulled through the baffle.

**Critical parameter: G × V (cooling rate) vs. G/V (temperature gradient per unit velocity):**

The stability criterion for planar solidification front (no constitutional supercooling):
```
G/V ≥ ΔT_0 / D_L

where:
  ΔT_0 = liquidus-solidus temperature interval (function of composition)
  D_L = diffusion coefficient of solute in liquid
```

**Too fast V:** Constitutional supercooling ahead of front → dendritic instability → stray grain nucleation
**Too slow V:** Acceptable for quality but economically impractical (long cycle time); also allows more time for segregation

Typical withdrawal rates: 1–5 mm/min for research; 3–8 mm/min for industrial.

**For a turbine blade ~150 mm tall, at 5 mm/min: 30 minutes of withdrawal time per blade.**

---

## 6. Dendrite Solidification and Microsegregation

Even in a perfect single crystal, solidification is not perfectly homogeneous. The solidification proceeds by **dendritic growth**:

```
Dendrite tree structure (one crystal, many dendrite arms):

     Primary arm (⟨001⟩, upward):
         ↑
         |── Secondary arm (⟨001⟩, sideways)
         |      |── Tertiary arm
         |
         |── Secondary arm (other side)
         |
```

The region between dendrite arms — the **interdendritic region** — solidifies last. Because each element has a different partitioning coefficient (k = concentration in solid / concentration in liquid), the interdendritic region is enriched in elements with k < 1 and depleted in those with k > 1.

**Partition coefficients in CMSX-4 (approx.):**
| Element | k = C_solid/C_liquid | Effect |
|---------|---------------------|--------|
| Ni | ~1.0 | Neutral |
| Cr | ~1.1 | Partitions to solid (dendrite core) |
| Al | ~0.9 | Slightly to interdendritic |
| Re | ~1.3 | Strongly to dendrite core |
| W | ~1.4 | Strongly to dendrite core |
| Ta | ~0.8 | To interdendritic |
| Mo | ~1.1 | To dendrite core |

This creates **chemical segregation** (microsegregation) on the scale of the dendrite arm spacing (~100–500 μm):
- Dendrite cores: rich in Re, W, Cr
- Interdendritic: rich in Ta, Al, (and eutectic phases)

The as-cast microstructure also contains:
- **Eutectic pools**: regions that solidified last are pushed to eutectic composition → (γ + γ′) eutectic pockets with larger γ′
- **Coarse irregular γ′**: a few μm in size, not the fine cuboidal γ′ we want
- **Residual porosity**: microshrinkage from solidification contraction

This as-cast microstructure is chemically and microstructurally far from the service-ready condition.

---

## 7. The Investment Casting Mold with Ceramic Core

The hollow internal cooling architecture of modern turbine blades (Chapter 56) requires internal passages that must be formed during casting — you can't drill them into a single crystal afterward (would cause recrystallization).

### The Ceramic Core

A **ceramic core** is made first (silica-based, alumina-based, or combination):
- The core's shape defines the internal cooling passages
- It must withstand the 1550°C casting temperature without melting or reacting with the alloy
- It must be removable after casting — typically dissolved in NaOH or HF solutions

```
Blade cross-section:
╔════════════════════╗
║  alloy wall        ║
║  ┌──────────────┐  ║
║  │ ceramic core │  ║ ← this space becomes the cooling channel
║  └──────────────┘  ║
╚════════════════════╝
```

Ceramic core design is one of the most challenging aspects of turbine blade manufacturing. Cores can have wall thicknesses of 0.3–0.5 mm and extremely complex 3D geometries (Z-shaped channels, pin fin arrays, etc.) in a fragile ceramic body.

### The Complete Mold Assembly

**Investment casting process (lost-wax):**
1. **Wax injection**: molten wax injected around the ceramic core → wax "blade" with internal core
2. **Shell building**: wax assembly repeatedly dipped in ceramic slurry and stucco, building up 5–10mm ceramic shell over 7–10 dipping cycles
3. **Dewax (autoclave)**: steam or oven at 170°C melts the wax out → hollow ceramic shell containing the ceramic core
4. **Pre-fire**: ceramic shell fired at ~1000°C to develop strength
5. **Metal pour**: done in vacuum Bridgman furnace; alloy superheated to ~1550°C; mold filled by gravity or vacuum assist

**The spiral selector is part of the wax pattern** — it gets coated along with the blade and becomes a cavity in the shell that guides solidification. It doesn't survive as hardware; its function is one-use.

---

## 8. Post-Solidification Processing

The as-cast blade needs significant processing before it can be used:

### 1. Knockout / Shell Removal
The outer ceramic shell is mechanically or vibrationally knocked off. This is aggressive — must not damage the blade.

### 2. Core Removal
Chemical leaching in hot NaOH or HF solution dissolves the silica/alumina core, leaving the hollow cooling channels. This takes hours to days depending on core chemistry and channel complexity.

### 3. Hot Isostatic Pressing (HIP)
At ~1170°C, 175 MPa argon for 4 hours:
- Closes any residual casting porosity
- Pressure collapses pores; temperature enables diffusion bonding of pore walls
- Mandatory step for all modern SX blades

### 4. Solution Heat Treatment
Above γ′ solvus (typically 1280–1335°C) for 2–6 hours:
- Dissolves as-cast eutectic γ′ and segregated regions
- Homogenizes chemistry across dendrite arm spacing
- Must NOT cause recrystallization (requires no mechanical damage)

### 5. Aging Heat Treatments
As described in Chapter 44 — primary and secondary aging to develop optimal γ′ size and distribution.

### 6. Inspection
Fluorescent penetrant inspection (FPI), X-ray (for porosity and core remnants), Laue X-ray diffraction (orientation measurement).

---

## 9. Orientation Measurement and Qualification

Each blade must be measured to verify crystallographic orientation. Two measurements are made:

**Primary orientation:** Angle of [001] direction from the blade axis.
- Specification: typically < 10° for aerospace HPT blades
- Measured by: Laue back-reflection X-ray diffraction on the blade root

**Secondary orientation:** Rotation of crystal around the [001] axis.
- Determines which crystallographic direction faces the leading edge, trailing edge, pressure face, suction face
- Important for thermal fatigue because stress concentrations at edges activate different slip systems
- Specification varies; typically ±15–25° range around the preferred secondary orientation

**Laue X-ray measurement:**
```
X-ray beam → hits crystal face → back-diffracted spots form pattern on film/detector
Pattern of spots reveals: primary orientation (spot position), secondary orientation (spot rotation)
Measurement time: ~2–5 minutes per blade
```

Every SX turbine blade in commercial aviation is individually measured and has a serial number traceable to its orientation measurement. This is a regulatory requirement.

---

## 10. Common Defects and Their Causes

### Stray Grain (Most Critical)

A grain with different orientation growing inside or adjacent to the crystal.

**Causes:**
- Constitutional supercooling ahead of the solidification front → new nucleus forms
- Dendrite fragmentation (mechanical disturbance, convection in liquid)
- Solidification in a geometrically complex region (blade airfoil platform corner) where thermal gradient is locally low

**Detection:** X-ray Laue diffraction; etching reveals grain boundaries under optical microscope.

**Consequence:** Reduced creep and fatigue life at the boundary → blade rejected.

### Low-Angle Grain Boundary (LAGB)

A region where the crystal orientation gradually changes by 2–15° over a small distance. Structurally an array of dislocations (Chapter 5).

**Cause:** Two slightly misoriented grains grew toward each other and merged.

**Specification:** Maximum allowable misorientation typically 2–5° for HPT blades. Above 5°, blade rejected.

**Consequence:** LAGB is less harmful than a true grain boundary but can accelerate creep if misorientation exceeds specification.

### Freckles

Strings of equiaxed grains or compositionally segregated channels running parallel to the withdrawal direction.

**Cause:** **Thermosolutal convection** (also called freckling) — buoyancy-driven convection in the mushy zone. Regions of low-density, Cr/Re-depleted liquid rise through the dendritic network and solidify as equiaxed grains.

**Prevention:** Low withdrawal rate, appropriate alloy composition (minimize density difference between liquid fractions), optimized mold geometry.

### Porosity

Microshrinkage pores at the last-to-solidify interdendritic regions.

**Cause:** Shrinkage of liquid on solidification without adequate feeding.

**Resolution:** HIP closes porosity. Remaining unhealed porosity → rejection if above specification size/density.

---

## 11. Industrial Scale and Process Control

A major jet engine manufacturer (GE Aviation, Rolls-Royce, Pratt & Whitney, IHI, Safran) producing a commercial engine at scale:

**HPT blade production data (approximate):**
- Large turbofan HPT: 60–80 blades per stage, 2 stages
- Fleet requirements: thousands of engines in service
- Annual blade replacement need: hundreds of thousands of blades

**Per-blade production time:**
- Shell building: 7–10 days
- Casting (furnace time): 2–4 hours per batch
- Solution + aging heat treatment: 1–2 days
- Inspection: 1–2 days
- Total: 2–3 weeks from mold to inspected blade

**Statistical yield rates (approximate):**
- Shell cracking/leaks: 5%
- Casting defects (stray grain, LAGB, freckle): 15–25%
- Dimensional out-of-spec: 5%
- Overall first-pass yield: ~65–75%

Process control systems monitor:
- Furnace temperature at 20+ zones (thermocouple array)
- Withdrawal rate (servo-controlled motor, ±0.01 mm/min)
- Vacuum level (< 0.001 mbar)
- Cooling gas flow rates
- Real-time thermal imaging of baffle region

Each casting run generates gigabytes of process data, all linked to the blade serial numbers for traceability.

---

## Summary

| Topic | Key Point |
|-------|-----------|
| Grain selector | Geometric spiral/pigtail; selects near-[001] orientation by dendritic competition |
| Seed crystal | Pre-oriented seed at mold base; epitaxial growth; precise primary + secondary control |
| Bridgman furnace | Resistance heater + radiation baffle + controlled withdrawal; vacuum environment |
| Thermal gradient G | Must be high to prevent stray grain nucleation ahead of front |
| Withdrawal rate V | Must satisfy G/V ≥ ΔT₀/D_L for stable planar-front solidification |
| Dendritic segregation | Re/W to dendrite cores; Ta/Al to interdendritic; heterogeneous as-cast microstructure |
| Ceramic core | Defines hollow cooling channels; removed by chemical leaching post-cast |
| Post-processing | HIP → solution treat → age → inspect |
| Orientation spec | Primary < 10° from [001]; secondary ±15–25° range |
| Defects | Stray grains, LAGB, freckles, porosity → rejection criteria strictly enforced |
| Industrial yield | ~65–75% first-pass; each blade individually tracked |

**Next chapter:** We examine the specific alloy compositions that define each generation of single-crystal superalloys — the CMSX series, René series, PWA series, and TMS series — and understand how each element addition contributes to the generational improvements in temperature capability.

---

## Exercises

1. In a Bridgman furnace, the gradient ahead of the solidification front is G = 50°C/mm and withdrawal rate is V = 5 mm/min. The diffusion coefficient of Re in liquid Ni at 1450°C is approximately 3×10⁻⁹ m²/s, and the liquidus-solidus interval is ~50°C. Does the solidification satisfy the constitutional supercooling criterion G/V ≥ ΔT₀/D_L? (Careful with units: convert V to m/s.)

2. A blade 150 mm long is cast by Bridgman at 5 mm/min withdrawal rate. The crystal starts growing at the grain selector exit. Sketch the temperature profile in the blade at t = 15 min. How long is the solidified vs. liquid portion of the blade? How long does the full withdrawal take?

3. CMSX-4 has Re partitioning to dendrite cores with k_Re ≈ 1.3. If the nominal alloy composition is 3 wt% Re, estimate the Re concentration in: (a) the dendrite core, (b) the interdendritic region, using the Scheil equation C_solid = k×C₀×(1-f_s)^(k-1) where f_s is the solid fraction. Take f_s = 0.1 for the core and f_s = 0.9 for near-complete solidification.

4. You are designing a spiral grain selector for a 10 mm diameter spiral, 25 mm long. A grain growing at angle θ from [001] will traverse a path through the spiral proportional to 1/cos θ. If the maximum selector length is 25 mm, what is the maximum grain misorientation θ that can survive to the blade? (Hint: a grain at θ = 30° travels the same path length as a straight-through grain travels 25/cos30° = 28.9 mm — longer than the selector.) What does this mean for the typical scatter in final orientation?

5. Explain why the secondary orientation of a single-crystal blade matters for thermal fatigue at the leading edge, even though primary [001] determines the creep performance. What crystallographic planes are active at the leading edge stress concentrations, and why does secondary orientation affect which of those planes is favorably oriented for slip?
