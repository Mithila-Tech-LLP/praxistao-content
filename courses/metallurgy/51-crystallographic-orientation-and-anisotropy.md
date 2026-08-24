# Chapter 51: Crystallographic Orientation and Anisotropy in SX Blades

> **"A single crystal turbine blade is not just 'without grain boundaries.' The direction in which it grows, the specific crystallographic axis aligned with the blade's length, determines how it deforms elastically, how it creeps, how fatigue cracks prefer to grow, and even how it responds to the anisotropic thermal gradient in the engine. The [001] direction is not an arbitrary choice — it is the direction of minimum elastic modulus, minimum creep rate, and minimum thermal stress. Every single crystal blade is grown [001]-up for physics reasons, not convention."**

---

## Table of Contents

1. [Anisotropy in Cubic Crystals — Basics](#1-anisotropy-in-cubic-crystals--basics)
2. [Elastic Anisotropy of FCC Ni Superalloys](#2-elastic-anisotropy-of-fcc-ni-superalloys)
3. [Why [001] for Turbine Blades?](#3-why-001-for-turbine-blades)
4. [Primary and Secondary Orientation](#4-primary-and-secondary-orientation)
5. [Schmid Factor and Slip in SX](#5-schmid-factor-and-slip-in-sx)
6. [Creep Anisotropy in SX](#6-creep-anisotropy-in-sx)
7. [Fatigue Anisotropy in SX](#7-fatigue-anisotropy-in-sx)
8. [Measuring Orientation (Laue, EBSD)](#8-measuring-orientation-laue-ebsd)
9. [Orientation Tolerance and Misorientations](#9-orientation-tolerance-and-misorientations)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Anisotropy in Cubic Crystals — Basics

An **isotropic** material has the same properties in all directions (polycrystal with random texture).

A **single crystal** is the extreme case of anisotropy — properties depend strongly on crystallographic direction.

**Cubic symmetry (FCC):**
Even with cubic symmetry, elastic, plastic, and thermal properties are direction-dependent:
- Elastic modulus: E varies by direction
- Yield stress: depends on slip system orientation relative to load
- Creep rate: varies with load direction
- Thermal expansion: for cubic — actually isotropic! (cubic symmetry)

**Miller Indices review (cubic):**
- Direction [uvw]: vector components in the unit cell
- Plane (hkl): plane with intercepts 1/h, 1/k, 1/l
- Family of directions ⟨uvw⟩: all equivalent directions
- Family of planes {hkl}: all equivalent planes

**Key directions in FCC (Ni):**
```
[001] ─── cube edge direction (growth direction in SX blades)
[011] ─── face diagonal
[111] ─── body diagonal (slip direction in FCC)
⟨110⟩  ── family of 12 equivalent face diagonals
{111}  ── family of 8 equivalent octahedral planes (slip planes in FCC)
```

---

## 2. Elastic Anisotropy of FCC Ni Superalloys

**For a cubic crystal, three independent elastic constants:** C₁₁, C₁₂, C₄₄

**Elastic modulus in direction [uvw]:**
```
1/E[uvw] = 1/E[100] - 2A' × (l₁²l₂² + l₂²l₃² + l₃²l₁²)

where A' = S₁₁ - S₁₂ - S₄₄/2 (S = compliance tensor components)
      l₁, l₂, l₃ = direction cosines of [uvw]
```

**Zener anisotropy ratio:**
```
A = 2C₄₄ / (C₁₁ - C₁₂)

For isotropic solid: A = 1
For Ni: A ≈ 2.5 (significantly anisotropic)
For CMSX-4 (with γ′): A ≈ 2.2–2.8
```

**Elastic modulus values for CMSX-4:**
| Direction | E (GPa) |
|-----------|---------|
| [001] | ~130 GPa (lowest!) |
| [011] | ~196 GPa |
| [111] | ~292 GPa (highest!) |

The [001] direction has the LOWEST modulus. This is crucial for blade design.

**Physical reason for E[001] < E[111]:**
In FCC, atoms are most densely packed along ⟨110⟩ (face diagonals). The [001] direction has lower linear density → atoms are less constrained → more compliant.

---

## 3. Why [001] for Turbine Blades?

**Reason 1 — Minimum thermal stress:**
Thermal stress = E × α × ΔT (where α = coefficient of thermal expansion)

In a turbine blade, there is a large thermal gradient:
- Blade tip: cooler (impingement cooling)
- Blade mid-span: hotter
- Leading edge: hot; trailing edge: cooler

ΔT across the blade section can be 100–200°C → large thermal stress.

Since E[001] is the LOWEST (130 GPa vs 292 GPa for [111]):
- Using [001] growth direction → thermal stress is 2.2× lower than [111] growth
- Lower thermal stress → longer fatigue life → less thermal cracking

**Reason 2 — Minimum creep rate in [001]:**
FCC slip systems are {111}⟨110⟩. The Schmid factor for [001] loading is:
- m = cos(φ)cos(λ) for each slip system
- For [001] loading: cos(φ)cos(λ) = (1/√3)(1/√2) = 0.408 for each of 8 active systems
- But the key: at high T, creep involves dislocation climb + glide. [001] loading activates multiple slip systems SIMULTANEOUSLY → work hardening through dislocation interactions → SLOWER creep rate than other orientations
- This is called "elastic shielding" — the combination of low E + orientation hardening effect

**Experimental verification:**
```
Creep life (1,050°C / 150 MPa, time to 1% strain):
  [001]: 200 hours (reference)
  [011]: 80 hours (0.4×)
  [111]: 40 hours (0.2×)
  [013]: 120 hours (0.6×)
```

**Reason 3 — Minimum oxidation fatigue damage:**
Lower thermal stress → smaller crack opening displacement at cooling holes → longer time before thermal-mechanical fatigue initiation at holes.

---

## 4. Primary and Secondary Orientation

**Primary orientation:** The angle between [001] and the blade's length axis (span direction, also the centrifugal stress direction). This is what most people mean when they say "blade orientation."

**Secondary orientation:** Rotation of the crystal about the [001] axis. Even with perfect [001] primary alignment, the blade can be rotated (imagine a square cross-section: [100] can point in various directions around the perimeter).

**Why secondary orientation matters:**
- The blade cross-section is not square — it is an airfoil (asymmetric)
- Loading in the chord direction (bending from gas pressure) is different from span direction (centrifugal)
- For bending fatigue: secondary orientation determines modulus in the chord direction
- For fatigue crack growth from cooling holes: crack propagates on specific crystallographic planes → secondary orientation determines which plane is exposed at the hole surface

**Standard orientation convention (turbine blade):**
- Primary [001] axis: ∥ blade span (centrifugal axis) ± 10°
- Secondary [100] axis: ∥ chord (direction from leading edge to trailing edge) ± 10°

This is specified in the drawing and verified by Laue diffraction (Ch 40, 52).

---

## 5. Schmid Factor and Slip in SX

**Schmid's law:**
Slip begins on a specific system when the resolved shear stress τ_RSS reaches the critical resolved shear stress τ_CRSS:

```
τ_RSS = σ × cos(φ) × cos(λ) = σ × m

where φ = angle between stress axis and slip plane normal
      λ = angle between stress axis and slip direction
      m = Schmid factor (0 ≤ m ≤ 0.5)
```

**Slip systems in FCC (Ni superalloys):**
- Slip planes: {111} (four families: (111), (1̄11), (11̄1), (111̄))
- Slip directions: ⟨110⟩ (three in each {111} plane)
- Total: 4 × 3 = 12 slip systems

**Schmid factor for [001] loading:**
For each {111}⟨110⟩ system:
- φ = angle between [001] and {111} normal = arccos(1/√3) = 54.7°
- λ = angle between [001] and ⟨110⟩ = 45°
- m = cos(54.7°) × cos(45°) = 0.577 × 0.707 = 0.408

**Schmid factor for [111] loading:**
- The slip direction closest to [111] is one specific ⟨110⟩
- m_max = 0.272 for [111] loading (low!) — this is why [111] appears very strong elastically

**Multiple slip in [001] SX:**
Under [001] loading, 8 slip systems have EQUAL Schmid factor (0.408). They all activate simultaneously → strong dislocation interactions → high work hardening rate → resistance to creep deformation.

---

## 6. Creep Anisotropy in SX

SX creep is highly anisotropic — strong dependence on orientation:

**[001] creep (best):**
- 8 simultaneous slip systems → strong work hardening
- Lowest E → lowest thermal stress → lower driving force for creep under thermal loading
- Rafting: under [001] tension, γ′ rafts PERPENDICULAR to load → creates channels; at high T, the γ channels are obstacles to dislocation passage

**[011] creep (intermediate):**
- 4 primary slip systems → less work hardening
- Higher E → more stress for same strain

**[111] creep (worst at moderate T):**
- Only 2 slip systems active
- Very high resolved shear stress on those systems → easy glide → low ductility + poor creep

**At very high temperature (> 1,100°C):**
Creep transitions from {111}⟨110⟩ slip to {100}⟨110⟩ slip in some alloys.
The orientation dependence reverses at some conditions — [111] can become the best direction.

**Rafting and creep:**
```
Stage I creep (0–100 hours):
  γ′ is still cuboidal → creep by dislocation glide through γ channels

Stage II creep (100+ hours):
  γ′ rafts form (plates perpendicular to [001] stress axis)
  Creep continues by dislocation climb over γ′ rafts

Stage III:
  Rafts coarsen, lose coherency → creep accelerates → rupture
```

---

## 7. Fatigue Anisotropy in SX

**Fatigue crack growth in SX:**
Cracks propagate on specific crystallographic planes:

- **Stage I crack (initiation):** Propagates on {111} planes (maximum shear stress planes)
- **Stage II crack (propagation):** Usually on {001} planes (planes perpendicular to maximum tensile stress) for [001]-oriented blades

**Secondary orientation controls Stage I fatigue:**
At a film cooling hole, the stress concentration causes Stage I cracks to initiate on {111} planes. The angle between the {111} plane and the hole surface depends on secondary orientation:
- Certain secondary orientations: {111} planes are tangential to hole → crack runs along surface → shallow, less damaging
- Other secondary orientations: {111} planes are radial to hole → crack runs into the blade → rapidly deepens → dangerous

**Crack path selection:**
Once initiated, a fatigue crack follows the lowest energy path:
- At low ΔK: Stage I on {111} (crystallographic fracture)
- At high ΔK: Stage II on planes perpendicular to maximum stress (non-crystallographic)

**Fracture toughness anisotropy:**
K_Ic also varies with crack plane and direction:
- K_Ic on {001} ≈ 45 MPa√m
- K_Ic on {111} ≈ 35 MPa√m
- Cleavage planes in Ni superalloys are {001} (weakest planes at low T)

---

## 8. Measuring Orientation (Laue, EBSD)

Detailed in Ch 40. Summary specific to SX blades:

**Laue back-reflection (production standard):**
- 100% inspection of every SX blade
- 10-minute exposure per blade
- Film or area detector
- Determines: primary [001] deviation, secondary [100] deviation, stray grains

**EBSD (for R&D, failure analysis):**
- Map orientation variation within the crystal (local misorientation)
- Identify sub-grain boundaries (LAGBs) that Laue misses
- Sensitivity: 0.1° misorientation

**Synchrotron X-ray Laue tomography (R&D):**
- 3D orientation mapping throughout the blade volume
- Detect LAGB networks deep inside the airfoil
- Used for alloy development programs, not routine production

---

## 9. Orientation Tolerance and Misorientations

**Why is orientation tolerance not zero?**
The Bridgman process (Ch 49) cannot produce perfect [001] every time:
- Grain selector slightly misaligns the preferred growth direction
- Thermal fluctuations during solidification
- Practical casting yield requires ±10° tolerance (otherwise rejection rate too high)

**Effect of primary misorientation:**
At 10° misorientation from [001]:
- Elastic modulus increases from 130 GPa toward [011] value
- Thermal stress increases by ~5%
- Creep life decreases by ~10%
- This is deemed acceptable for production components

**Low-Angle Grain Boundaries (LAGBs) within SX:**
A blade may "pass" Laue inspection but contain internal LAGBs — boundaries between slightly misoriented regions (sub-grains). See Ch 52 for details.

**Critical misorientation angles for rejection:**
| Feature | Tolerance | Effect if exceeded |
|---------|----------|-------------------|
| Primary [001] deviation | ± 10° | Thermal stress increase; creep life reduction |
| Secondary [100] deviation | ± 10° | Fatigue crack path change |
| LAGB (internal) | > 2° concern, > 5° reject | Preferred creep/fatigue crack path |
| Stray grain | Any (0° tolerance) | Automatic reject |

---

## Summary

| Property | [001] | [011] | [111] | Reason |
|----------|-------|-------|-------|--------|
| Elastic modulus (GPa) | 130 | 196 | 292 | Atom packing density |
| Thermal stress (relative) | 1.0 | 1.5 | 2.2 | E × α × ΔT |
| Creep life (relative) | 1.0 | 0.4 | 0.2 | Multiple simultaneous slip |
| Schmid factor (max) | 0.408 | 0.408 | 0.272 | Geometry |
| Fatigue crack plane | {001} | Mixed | {111} | Crack plane selection |

---

## Exercises

1. CMSX-4 blade test specimens machined in [001], [011], and [111] directions. (a) Calculate the Schmid factor m for [001], [011], and [111] loading on the {111}[110] slip system (use cos(φ)cos(λ) formula). [Hint: for [001] tension, φ = 54.7°, λ = 45°; for [111] tension, φ = 70.5°, λ = 35.3°; for [011] tension, φ = 35.3°, λ = 60°.] (b) What is the CRSS (critical resolved shear stress) if the specimen yield strength is 880 MPa in the [001] direction? (c) Predict the yield strength in [111] using the same CRSS. (d) Experimentally, [111] yield strength ≈ 1,600 MPa (much higher). Explain why the simple Schmid factor calculation underpredicts [111] strength. (hint: multiple vs single active slip systems)

2. Thermal stress comparison in a [001] vs [111] SX blade: A 100°C temperature difference exists across the blade wall (10 mm thick section). (a) Calculate the thermal stress for [001] (E = 130 GPa) and [111] (E = 292 GPa) orientations, assuming α = 14 × 10⁻⁶/°C and σ_thermal = E × α × ΔT (free-expansion approximation). (b) If the fatigue endurance limit = 600 MPa for [001] and 500 MPa for [111] (lower due to lower ductility), what safety factor does each orientation provide against thermal fatigue? (c) In the Goodman diagram (Ch 14), add the thermal stress as mean stress. With alternating LCF stress = 200 MPa, calculate the effective safety factor for each orientation. (d) A [013] misoriented blade (10° from [001] toward [011]) has E ≈ 150 GPa. Estimate its thermal stress and relative creep life vs a perfect [001] blade.

3. Secondary orientation effect on cooling hole fatigue: A cylindrical film cooling hole (diameter 0.5 mm) is drilled perpendicular to the [001] growth axis. Two blades have different secondary orientations: (a) Blade A: [100] aligned with chord direction (standard); (b) Blade B: [110] aligned with chord direction (45° rotated). For each blade, identify which crystallographic plane ({111} family) is at the highest shear stress at the hole surface (tangential shear is maximum at hole). Use Schmid factor analysis for shear loading. (c) Blade A's {111} plane makes 35° angle with hole surface; Blade B's makes 10°. Which blade has the more damaging Stage I crack path? (d) Engine OEM specifies secondary tolerance ± 10°. Why is this secondary tolerance critical for hole-initiated fatigue?

4. Laue X-ray orientation measurement: A blade gives the following Laue pattern measurements (back-reflection geometry): (a) The (002) reflection spot is at angular position (ψ = 3°, φ = 180°) from the expected position for perfect [001]. Calculate the primary misorientation. Does the blade meet ±10° specification? (b) The (200) reflection is at (ψ = 6°, φ = 90°). Calculate secondary misorientation. (c) A second set of Laue spots appears, rotated 27° from the primary pattern. What does this indicate? What disposition action is required? (d) If the rejection rate for primary misorientation > 10° is 4% and for secondary > 10° is 3% (independent), and for stray grains is 8%, calculate the overall casting yield (fraction of blades that pass all three orientation requirements).

5. Creep anisotropy in the engine context: An HPT blade experiences: (a) Centrifugal stress along [001]: 150 MPa; (b) Bending stress from gas load along secondary [100] direction: 40 MPa; (c) Thermal stress along thickness direction [010]: 80 MPa. For isotropic analysis (polycrystal): use Mises criterion, σ_VM = √[(σ₁-σ₂)² + (σ₂-σ₃)² + (σ₃-σ₁)²] / √2. For SX: creep rate is governed by resolved shear stress on each {111}⟨110⟩ system. (i) Calculate σ_VM for the given stress state. (ii) For SX blade, the [001] direction has the slowest creep (minimum RSS). Calculate the RSS on the most active {111}⟨110⟩ system for the combined stress state. (iii) If τ_CRSS = 80 MPa (at service temperature), is the SX blade creeping? How does this compare to a DS or polycrystalline blade under the same loading?
