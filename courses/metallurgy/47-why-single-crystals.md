# Chapter 47: Why Single Crystals? — Eliminating Grain Boundaries

> **"Every grain boundary is a weakness. A crack path, a creep path, a corrosion path. The metallurgist's dream is a material with the composition of a superalloy but the perfection of a single crystal. In 1968, that dream became reality. Today, every high-pressure turbine blade in every commercial jet engine is grown as a single crystal."**

---

## Table of Contents

1. [The Problem With Polycrystals at High Temperature](#1-the-problem-with-polycrystals-at-high-temperature)
2. [Why Grain Boundaries Fail Under Creep](#2-why-grain-boundaries-fail-under-creep)
3. [Directional Solidification — The First Step](#3-directional-solidification--the-first-step)
4. [Single Crystals — The Complete Solution](#4-single-crystals--the-complete-solution)
5. [Quantifying the Benefit — Temperature Capability Comparison](#5-quantifying-the-benefit--temperature-capability-comparison)
6. [What You Gain by Eliminating Boundaries](#6-what-you-gain-by-eliminating-boundaries)
7. [What You Lose — Trade-offs of SX](#7-what-you-lose--trade-offs-of-sx)
8. [The Elastic Anisotropy Advantage](#8-the-elastic-anisotropy-advantage)
9. [Which Elements Can Now Be Removed?](#9-which-elements-can-now-be-removed)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Problem With Polycrystals at High Temperature

A **polycrystalline** metal is made of many individual grains, each with a different crystal orientation, separated by grain boundaries. At room temperature, grain boundaries are your friends:
- They block dislocation motion (Hall-Petch strengthening)
- They strengthen the material significantly

But at high temperature (T_H > 0.5), grain boundaries become **the weakest point in the structure**. They:
- Slide past each other under sustained stress
- Develop cavities at triple points
- Become pathways for accelerated diffusion
- Fail in an intergranular fracture mode with minimal warning

In the 1950s and early 1960s, jet engine turbine blades were polycrystalline investment castings of alloys like IN100 or Mar-M200. Engineers were hitting a wall — increasing operating temperature was being limited not by the bulk alloy's strength, but by grain boundary failure.

The question was: can we remove the weakest links — the grain boundaries — without sacrificing everything else?

---

## 2. Why Grain Boundaries Fail Under Creep

Chapter 16 introduced creep mechanisms. Let's focus specifically on grain-boundary damage in a turbine blade:

### 2.1 Grain Boundary Sliding (GBS)

At high temperature and stress, grain boundaries can slide — one grain shifts relative to its neighbor. 

The problem: **transverse grain boundaries** (perpendicular to the centrifugal stress direction) are directly pulled apart by the blade's centrifugal load. This creates opening mode cavitation.

```
Turbine blade with polycrystalline structure:

       ╔═══════╗ Centrifugal force (upward)
       ║ grain ║     ↑
  ─────╫───────╫─────← transverse grain boundary (⊥ to load)
       ║ grain ║     ← this boundary gets pulled open by centrifugal load
  ─────╫───────╫─────
       ║ grain ║
       ╚═══════╝
       ↓ disk attachment
```

Transverse grain boundaries → cavitation → cracks → intergranular fracture. This is the dominant blade failure mode in polycrystalline materials.

### 2.2 Triple-Point Cavitation

Where three grains meet (triple junction), geometric incompatibility during GBS creates stress concentrations → voids nucleate:

```
        Grain A
          ╱╲
         ╱  ╲
Grain B ╱    ╲ Grain C
       ●       ← triple junction under high local stress
       ↑
   void nucleates here during GBS
```

Once nucleated, voids grow by diffusion (Coble creep mechanism) and eventually link up → intergranular crack.

### 2.3 Grain Boundary Segregation and Precipitation

At service temperature, certain elements segregate to grain boundaries over time (especially S, P, B, and some carbide-forming elements). This:
- Changes local chemistry → local composition → altered phase stability
- Can cause low-melting-point phases to form at boundaries
- Can cause embrittlement

Grain boundary engineering (adding controlled amounts of B, Zr, Hf, C) can help, but it's fighting the fundamental problem rather than solving it.

---

## 3. Directional Solidification — The First Step

**Key insight (circa 1960):** Not all grain boundaries are equally harmful. **Transverse grain boundaries** (⊥ to stress) cause damage. **Longitudinal grain boundaries** (∥ to stress) experience shear rather than tension → far less damaging.

If you can't eliminate ALL grain boundaries, can you at least eliminate the transverse ones?

**Yes — by directional solidification (DS):**

In normal (conventional) casting, solidification begins at random nucleation sites throughout the mold → equiaxed grains in all orientations → many transverse boundaries.

In **directional solidification**, you control the solidification so it proceeds in ONE direction — along the blade length (centrifugal direction):

```
Conventional casting:           Directionally solidified:

  ╔═══════╗                      ╔═══════╗
  ║ grain ║← random orientations  ║ ∥ ∥ ∥ ║← columnar grains, all parallel
  ║ grain ║                       ║ ∥ ∥ ∥ ║   to blade axis
  ║ grain ║                       ║ ∥ ∥ ∥ ║
  ╚═══════╝                      ╚═══════╝
  
  Many transverse boundaries       No transverse boundaries!
  Creep limited by GBS             Longitudinal boundaries only
  Life: X                          Life: 2X or better
```

The DS process involves:
- Heating the mold above the alloy's liquidus
- Withdrawing it slowly from a furnace at a controlled rate into a cold zone
- Maintaining a steep thermal gradient so solidification sweeps upward
- Grains nucleate at the cold base plate and grow upward as long columnar crystals

**Result:** Multiple grains, but ALL elongated parallel to the blade axis. No transverse grain boundaries. Creep life approximately doubled at equivalent conditions.

First commercial DS turbine blades entered service in the **Pratt & Whitney JT9D in 1966** — blades for the Boeing 747's engine. This single process change improved temperature capability by ~30°C.

---

## 4. Single Crystals — The Complete Solution

DS was a great improvement. But you still have some longitudinal grain boundaries. Can you go further?

**Yes — single crystal (SX) casting:**

If you add a **grain selector** to the DS process (a helical spiral section at the bottom of the mold), only ONE grain can navigate through it and emerge into the main cavity. That single grain grows to fill the entire mold — the entire blade is ONE crystal, with no grain boundaries at all.

```
Single crystal blade cross-section:
  
  ╔═══════╗
  ║       ║ ← one continuous crystal lattice
  ║       ║   from root to tip, side to side
  ║       ║   identical atom arrangement
  ╚═══════╝
  
  ZERO grain boundaries
  ZERO GBS
  ZERO intergranular cracking
  ZERO triple-point cavitation
```

The first SX turbine blade was cast in a laboratory in 1968. It took until 1982 for the first commercial SX blade to enter service (PWA 1480 in the Pratt & Whitney F100 military engine). The technique was a closely guarded industrial secret through the 1980s.

Today, every commercial and military jet engine high-pressure turbine blade is a single crystal.

---

## 5. Quantifying the Benefit — Temperature Capability Comparison

The improvement in creep resistance is quantified by "equivalent temperature capability" — how much hotter can you run the material while maintaining the same creep life?

| Alloy | Structure | Stress Rupture Temp. at 200 MPa / 1000h |
|-------|-----------|----------------------------------------|
| IN100 | Polycrystal (PC) | ~930°C |
| Mar-M200 DS | Directionally solidified | ~960°C (+30°C) |
| PWA 1480 | Single crystal (Gen 1) | ~990°C (+30°C vs DS) |
| CMSX-4 | Single crystal (Gen 2, +3%Re) | ~1020°C (+30°C vs Gen 1) |
| CMSX-10 | Single crystal (Gen 3, +6%Re) | ~1045°C (+25°C vs Gen 2) |

Each major step (PC → DS → SX → Gen 2 → Gen 3) is worth ~25–30°C. In thermodynamics, +25°C at turbine inlet temperature improves engine thermal efficiency by ~0.5–1%. Over a commercial airline fleet of 3,000 aircraft, that's hundreds of millions of dollars in fuel savings per year.

The **temperature advantage of SX vs. PC** (holding alloy composition constant) is about 35–50°C from the elimination of grain boundaries alone.

---

## 6. What You Gain by Eliminating Boundaries

With no grain boundaries, all the damage mechanisms that operate at grain boundaries simply cannot occur:

| Mechanism Eliminated | Consequence |
|---------------------|-------------|
| Grain boundary sliding | No GBS cavitation |
| Nabarro-Herring creep | Eliminated (requires GB as vacancy source/sink) |
| Coble creep | Eliminated (no GB diffusion path) |
| Triple-point cavitation | Eliminated (no triple junctions) |
| Intergranular fracture | Eliminated |
| Grain boundary oxidation penetration | Eliminated |
| GB segregation embrittlement | Eliminated |

The remaining creep mechanisms are all dislocation-based (power-law creep) — operating in the grain interior. These are much more amenable to control by alloy design (Re addition, γ′ optimization) than grain-boundary mechanisms.

**New failure modes in SX:**
Eliminating grain boundaries doesn't make a blade perfect. New challenges appear:
- **Creep rafting and collapse**: eventual breakdown of γ/γ′ rafts
- **Stray grain formation**: accidental grain nucleation during casting (a major manufacturing challenge)
- **Oxidation**: the absence of grain boundaries slightly changes the oxide scale growth kinetics
- **Recrystallization**: during heat treatment, mechanical damage can trigger grain formation (must be avoided)

---

## 7. What You Lose — Trade-offs of SX

Single crystals are not without disadvantages:

### 7.1 Casting Difficulty and Yield Rate

Growing a perfect single crystal with no grain defects (stray grains, low-angle boundaries, freckles) is extremely difficult. Typical first-pass casting yields for SX blades are 50–70% — meaning 30–50% of blades are scrapped at casting. This significantly adds to cost.

Every parameter in the casting process (withdrawal rate, temperature gradient, alloy chemistry, mold design) must be tightly controlled to avoid defects.

### 7.2 Can't Use Grain Boundary Strengtheners

In polycrystalline alloys, grain boundary strengtheners — **boron (B), zirconium (Zr), carbon (C) at grain boundaries** — prevent embrittlement and improve ductility. These are essential in PC alloys.

In SX alloys, there are no grain boundaries, so these elements provide no benefit. They can even be harmful (B can form low-melting borides that cause incipient melting during heat treatment). SX alloy compositions intentionally have very low B, C, and Zr.

This means SX alloys can have **higher melting points** (no low-melting grain boundary phases) and can be **solution treated more aggressively** (higher temperatures) — allowing more complete homogenization and better γ′ dissolution.

### 7.3 Secondary Orientation Control

A single crystal is NOT isotropic — it has different properties in different crystallographic directions (elastic anisotropy — see §8). The primary [001] orientation along the blade is well controlled. But the **secondary orientation** (rotation around the [001] axis) matters for thermal fatigue and oxidation at leading/trailing edges. Controlling secondary orientation in casting is possible but adds complexity.

### 7.4 Cost

SX blades cost 3–5× more than equivalent PC blades, due to:
- Longer casting time (controlled withdrawal)
- Lower casting yield
- More complex post-cast inspection (X-ray, orientation checks)
- Specialized furnace equipment

This is why DS rather than SX is still used for some lower-temperature applications (LPT blades) where the cost is not justified by the temperature requirement.

---

## 8. The Elastic Anisotropy Advantage

Eliminating grain boundaries also allows the engineer to exploit **crystallographic anisotropy** — the variation of properties with direction in a single crystal.

For FCC Ni superalloys, the elastic modulus varies dramatically with direction:
- E[001] ≈ 125 GPa (cube face direction, most compliant)
- E[011] ≈ 220 GPa
- E[111] ≈ 294 GPa (cube diagonal, stiffest)

**Why [001] is chosen for blade primary orientation:**

Turbine blades experience large temperature gradients (cool leading edge, hot midchord) with every power cycle. These temperature gradients cause **thermal strains**. For the same temperature gradient:
```
Thermal stress = E × α × ΔT
```

If E is lower (as in [001] direction), the **thermal stress is lower** for the same thermal strain. Lower thermal stress → less thermal fatigue damage per cycle → more cycles before fatigue cracking.

In a polycrystalline blade, grains are randomly oriented → average E ≈ 207 GPa.  
In a [001] single crystal → E = 125 GPa → **40% lower** → thermal stresses are 40% lower → thermal fatigue life dramatically improved.

This elastic anisotropy advantage is ONLY available in a single crystal (polycrystal averages it out). This is an **additional** benefit of SX casting beyond grain boundary elimination.

---

## 9. Which Elements Can Now Be Removed?

Because SX alloys have no grain boundaries, certain alloying elements used specifically for grain boundary strengthening in PC alloys are no longer needed:

**Removed in SX:**
- **Carbon (C)**: Grain boundary carbides (MC, M₂₃C₆) strengthen PC alloys. In SX, they serve no purpose and reduce the solidus temperature. SX alloys: < 0.005 wt% C (vs. 0.05–0.2% in PC).
- **Boron (B)**: Segregates to grain boundaries, reduces intergranular fracture. Not needed in SX. < 0.001 wt% B.
- **Zirconium (Zr)**: Similar grain boundary benefit. Removed or minimized in SX.

**Benefit of removing C, B, Zr:**
The solidus temperature rises by ~30–50°C because these elements no longer form low-melting grain boundary films. This allows solution heat treatment at higher temperatures → better homogenization → better γ′ dissolution → cleaner microstructure → improved mechanical properties.

This is a non-obvious bonus: the alloy chemistry for SX is fundamentally different from PC alloys not just in the major elements but in these "minor" grain-boundary elements.

---

## Summary

- **Grain boundaries are the weakest point** in high-temperature creep: sliding, cavitation, intergranular fracture.
- **Directional solidification (DS)**: controls solidification direction → columnar grains ∥ to centrifugal stress → eliminates transverse boundaries → +30°C equivalent.
- **Single crystal (SX)**: adds a grain selector → one grain fills entire mold → ZERO grain boundaries → another +30–50°C equivalent.
- **Mechanisms eliminated**: GBS, N-H creep, Coble creep, intergranular fracture, grain boundary oxidation.
- **Elastic anisotropy**: [001] orientation chosen because E[001] = 125 GPa vs. 207 GPa for polycrystal → 40% lower thermal stresses per cycle.
- **SX alloy chemistry**: no B, C, Zr (not needed without grain boundaries) → higher solidus → better solution treatment.
- **Trade-offs**: lower casting yield (50–70%), new defect types (stray grains), secondary orientation control, higher cost.
- **Commercial reality**: all HPT rotating blades in modern jet engines are SX. First commercial use: 1982.

**Next chapter:** How do you actually make a single crystal? The Bridgman process and grain selector are the core technology — but the thermal gradient, withdrawal rate, and mold geometry must all be perfectly controlled to grow a defect-free crystal the size of your thumb.

---

## Exercises

1. A turbine blade experiences a temperature gradient causing ΔT = 200°C across the wall. Compare the thermal stress in: (a) a polycrystalline blade (E = 200 GPa), (b) a [001] single crystal blade (E = 125 GPa). Assume α = 13 × 10⁻⁶ /°C and use σ_thermal = E × α × ΔT. By what fraction does the SX reduce thermal stress?

2. In a DS blade, grain boundaries run parallel to the blade axis (centrifugal direction). The centrifugal stress is 200 MPa axial. What shear stress acts on the longitudinal grain boundaries? (Consider the resolved shear stress at 0°, 45°, and 90° to the applied stress direction.) Why are longitudinal boundaries much less damaging than transverse boundaries?

3. SX alloys deliberately have < 0.005% C vs. 0.05–0.15% C in polycrystalline alloys. Calculate the temperature increase in solidus achievable by removing carbon if each 0.01 wt% C depresses the solidus by approximately 3°C. How does a higher solidus help heat treatment?

4. Stray grain formation is a critical defect in SX casting. If a stray grain forms in the middle of the blade with a [111] orientation and the surrounding crystal is [001], what is the misorientation angle? Does this constitute a high-angle or low-angle grain boundary? What is the typical rejection criterion for SX blades (maximum misorientation from [001])? (Research: typical specification is 10° from [001] primary, and <15° secondary.)

5. Compare the cost structure of PC vs. DS vs. SX blades: if a PC blade costs $500 to produce and DS adds 50% cost with 80% yield, while SX further adds to $2000 per attempt with 60% yield, what is the actual cost per acceptable DS blade vs. SX blade? At what point does the performance improvement justify the cost?
