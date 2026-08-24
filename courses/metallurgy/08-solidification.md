# Chapter 08: Solidification — How Metals Freeze

> **"Every metal component you have ever touched was once a liquid. The way that liquid transformed into a solid — over seconds or hours, in a mold or in space — determined everything about the grain size, segregation, porosity, and residual stress of the final part. Solidification is not a passive event; it is a controlled transformation that the metallurgist engineers to produce the desired microstructure."**

---

## Table of Contents

1. [Why Solidification Matters](#1-why-solidification-matters)
2. [Nucleation — Getting the First Crystal Started](#2-nucleation--getting-the-first-crystal-started)
3. [Homogeneous vs. Heterogeneous Nucleation](#3-homogeneous-vs-heterogeneous-nucleation)
4. [Growth — How Crystals Expand](#4-growth--how-crystals-expand)
5. [Dendrites — The Tree-Like Growth Morphology](#5-dendrites--the-tree-like-growth-morphology)
6. [Constitutional Supercooling and the Mushy Zone](#6-constitutional-supercooling-and-the-mushy-zone)
7. [Grain Structure — Equiaxed vs. Columnar](#7-grain-structure--equiaxed-vs-columnar)
8. [Segregation — Uneven Distribution of Alloying Elements](#8-segregation--uneven-distribution-of-alloying-elements)
9. [Porosity — The Void Problem](#9-porosity--the-void-problem)
10. [Hot Cracking — Solidification Cracking](#10-hot-cracking--solidification-cracking)
11. [Cooling Rate Effects](#11-cooling-rate-effects)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Why Solidification Matters

Virtually all metals begin their life as a liquid poured into a mold. The solidification process determines:
- **Grain size and shape** → affects strength, ductility, fatigue life
- **Segregation** → non-uniform composition → affects heat treatment and properties
- **Porosity** → voids in the casting → weakens the part
- **Residual stresses** → locked-in stresses from differential cooling
- **Inclusions** → trapped slag or oxide particles

Understanding solidification is essential for:
- Designing casting processes (sand casting, die casting, investment casting)
- Understanding how welding works (weld pool = local solidification)
- Understanding DS and SX turbine blade casting (Ch 48, 49)

---

## 2. Nucleation — Getting the First Crystal Started

When liquid metal cools below its melting point (T_m), it doesn't solidify instantly. The atoms must organize into a crystal lattice — and forming a tiny crystal requires overcoming an energy barrier.

### Thermodynamic Driving Force

Below T_m, the solid phase is more stable (lower Gibbs free energy) than the liquid:
```
ΔG_volume = ΔG_v × (4/3 πr³)  ← negative (driving force for solidification)
ΔG_surface = γ_SL × (4πr²)     ← positive (cost of creating solid/liquid interface)

Total: ΔG_total = (4/3)πr³ × ΔG_v + 4πr² × γ_SL
```

At small radius r, the surface term dominates → ΔG_total is positive → the nucleus is unstable and dissolves.

At a **critical radius r***:
```
r* = -2γ_SL / ΔG_v = 2γ_SL × T_m / (L_f × ΔT)

where L_f = latent heat of fusion, ΔT = undercooling = T_m - T_actual
```

For r < r*: nucleus dissolves (thermodynamically unfavorable)
For r > r*: nucleus grows spontaneously

**The undercooling requirement:** The larger the undercooling ΔT, the smaller r* becomes, and the easier nucleation is. In practice, metals need only 0.1–10°C of undercooling for heterogeneous nucleation, but 100–300°C for homogeneous nucleation.

---

## 3. Homogeneous vs. Heterogeneous Nucleation

### Homogeneous Nucleation

Pure liquid, no surfaces or impurities. Requires large undercooling (~100–300°C) because the entire surface energy cost must be paid from scratch.

Rarely observed in engineering metals — impurities, mold walls, and dissolved gases always provide nucleation sites.

### Heterogeneous Nucleation

On existing surfaces (mold walls, oxide particles, inclusions, grain refiners):

```
Contact angle θ between liquid, solid, and substrate:

     Liquid
    ─────────────────────
      \  θ  /
       \ | /
        \|/
    ─────────────────────  Substrate (mold wall)
    
θ < 90°: solid "wets" substrate → low nucleation barrier
θ → 0°: perfect wetting → essentially zero barrier → nucleation at very small ΔT
```

The heterogeneous nucleation barrier:
```
ΔG*_het = ΔG*_hom × f(θ) where f(θ) = (2 + cosθ)(1 - cosθ)² / 4

For θ = 0°: f = 0 → no barrier (spontaneous)
For θ = 180°: f = 1 → same as homogeneous
```

**Grain refiners:** Deliberately added particles (e.g., Al-Ti-B for aluminum alloys) that are good nucleation substrates (small θ) → many nuclei form → fine equiaxed grain structure → better properties.

---

## 4. Growth — How Crystals Expand

Once a stable nucleus forms, it grows by atoms attaching from the liquid to the solid/liquid interface.

**Heat extraction is critical:** Solidification releases latent heat (L_f). This heat must be conducted away for freezing to continue. The rate of solidification is limited by how fast heat can be extracted.

**Growth morphology depends on temperature gradient in the liquid:**

**Planar front:** If the liquid ahead of the front is always superheated (above T_m), any protrusion into the liquid would immediately be in hotter liquid → protrusion melts back → flat interface maintained. Occurs only in carefully controlled DS casting.

**Dendritic growth:** If liquid ahead is constitutionally supercooled (see §6), protrusions grow faster → dendritic morphology develops.

---

## 5. Dendrites — The Tree-Like Growth Morphology

The most characteristic solidification structure is the **dendrite** — a tree-like crystal with a primary trunk and secondary and tertiary branches:

```
Dendrite structure:
         |  ← primary arm (fast growth direction)
    ─────┼─────
    ─────┼─────  ← secondary arms (perpendicular to primary)
    ─────┼─────
         |
```

**Why dendrites form:**
- Crystal grows fastest in preferred crystallographic directions ([001] for FCC, [100] for BCC)
- Constitutional supercooling ahead of the tip accelerates growth in "finger" directions
- Heat is extracted fastest from sharp protrusions (larger surface area) → tips grow fastest
- Secondary arms nucleate on primary arms at specific crystallographic angles

**Dendrite arm spacing (DAS):**
The distance between secondary dendrite arms λ depends on cooling rate:
```
λ = C × (cooling rate)^(-n)   where n ≈ 0.3–0.5, C is alloy-dependent
```

Faster cooling → smaller DAS → finer microstructure → shorter diffusion distances → better properties. This is why rapidly cooled castings (die casting) have better properties than slowly cooled sand castings.

---

## 6. Constitutional Supercooling and the Mushy Zone

In alloys (not pure metals), the solidification front is not a sharp plane but a **mushy zone** — a region of solid + liquid coexisting.

As discussed in Chapter 48, the alloy's liquidus and solidus temperatures differ (ΔT_0 = T_liquidus - T_solidus). Within the mushy zone, partial solidification has occurred.

**Constitutional supercooling:** The liquid ahead of the solidification front has been enriched in rejected solute (elements with partition coefficient k < 1 partition to the liquid). This solute-enriched liquid has a LOWER liquidus temperature → even though the actual temperature is above the nominal liquidus, the local liquid is below ITS liquidus → "constitutionally supercooled" → drives dendritic instability.

```
Constitutional supercooling diagram:

Temperature
    |   T_liquidus of nominal alloy
    |──────────────────────────────────
    |            ↗ actual temperature profile
    |           /  (imposed by cooling)
    |          /
    |─────────────────────────────────── T_liquidus of enriched liquid
    |                    ← supercooled region →
    |
    └──────────────── distance from solidification front →
```

The region where actual T < local T_liquidus is constitutionally supercooled → any perturbation of the interface grows into this region → dendrites form.

---

## 7. Grain Structure — Equiaxed vs. Columnar

The macro-scale grain structure of a casting depends on how many nucleation sites are active and what the temperature gradients are:

```
Cross-section of a simple casting:

    ┌─────────────────────────────────────┐
    │ Chill zone (fine equiaxed grains)   │← rapid cooling at mold wall
    │ ─────────────────────────────────── │
    │ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║  │← columnar grains (grew inward)
    │ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║ ║  │   heat extracted through walls
    │ ─────────────────────────────────── │
    │   Equiaxed zone (center)            │← bulk nucleation, random growth
    └─────────────────────────────────────┘
```

**Chill zone:** Fine equiaxed grains at mold wall from rapid heat extraction + many heterogeneous nuclei (mold surface).

**Columnar zone:** Grains that nucleated at the wall grew inward. Fast-growing [001] grains outcompete others → selection → fewer, elongated columnar grains pointing toward center (heat flow direction).

**Equiaxed center:** If the center liquid cools slowly enough with enough nucleation sites (dendrite fragments, inoculants), equiaxed grains form throughout.

**Control:** Grain refiners → more equiaxed; directional heat extraction → more columnar; DS/SX process → fully columnar or single grain.

---

## 8. Segregation — Uneven Distribution of Alloying Elements

When an alloy solidifies, elements distribute unevenly between the solid and liquid phases according to their **partition coefficient k**:
```
k = C_solid / C_liquid (at equilibrium)
```

Elements with k < 1 (most common) concentrate in the LAST liquid to solidify (interdendritic regions).
Elements with k > 1 concentrate in the FIRST solid (dendrite cores).

**Microsegregation:** Composition variation on the scale of dendrite arm spacing (~50–500 μm). Eliminated by homogenization heat treatment.

**Macrosegregation:** Large-scale composition variation across the casting. Caused by liquid convection (denser/lighter liquid moving), shrinkage-driven flow. More difficult to eliminate.

**The Scheil equation** (no diffusion in solid, perfect mixing in liquid):
```
C_s = k × C_0 × (1 - f_s)^(k-1)

where C_s = solid composition at fraction f_s solidified
      C_0 = nominal alloy composition
      k = partition coefficient
```

For k < 1: as f_s → 1 (last to solidify), C_s → very high (severe enrichment). This is why the last interdendritic liquid is highly enriched in low-k elements → can form low-melting eutectic phases → hot cracking risk (§10).

---

## 9. Porosity — The Void Problem

Casting porosity is one of the most common defects. Two main types:

### Shrinkage Porosity

Liquid metal is less dense than solid → as it freezes, volume decreases (~3–7% for most metals). If liquid can't feed the solidifying region (long mushy zone, isolated hot spots), a void forms:
```
Pipe (large shrinkage cavity):        Microporosity (interdendritic):
                                       
  ┌─────────┐                          ●●●●   ← pores between dendrite arms
  │  void   │                          ●●●
  │  (top)  │                        common in alloys with wide freezing range
  └─────────┘
  common in pure metals/
  narrow-freezing-range alloys
```

**Prevention:**
- Risers (feeders): large reservoir of liquid at top that feeds shrinkage as part solidifies
- HIP (Hot Isostatic Pressing): pressure + temperature closes pores after casting
- Controlled solidification direction (solidify bottom-up so liquid always feeds from above)

### Gas Porosity

Gases dissolved in liquid metal (H₂, N₂, CO) come out of solution during solidification (solubility decreases sharply on solidification) → spherical gas pores.

**Aluminum alloys:** Hydrogen from moisture → H₂ gas porosity very common
**Steel:** N₂, CO from dissolved carbon + oxygen reactions
**Prevention:** Degassing of liquid metal, controlled atmosphere casting

---

## 10. Hot Cracking — Solidification Cracking

In the final stages of solidification, when only thin liquid films remain between dendrites, the solid skeleton contracts thermally and cannot be fed by liquid. If thermal stresses exceed the strength of the thin liquid films → solidification cracking ("hot tearing"):

```
Hot tear formation:
         
  ─ solid ─    ─ solid ─
           │  │  ← thin liquid film
           │  │
  ← pulling apart from thermal contraction → 
           │  │
           ↓  ↓
           CRACK
```

**Alloys susceptible to hot cracking:**
- Wide solidification range (T_liquidus - T_solidus large) → long mushy zone → thin liquid films persist for long time → more risk
- High alloying: aluminum alloys (2xxx, 7xxx), nickel superalloys (IN718)
- Restrained castings or welds

**Prevention:**
- Alloy composition adjustment: avoid compositions near maximum solidification range
- Preheat and control heat input in welding
- Grain refinement: finer equiaxed grains accommodate strain better

---

## 11. Cooling Rate Effects

Cooling rate during solidification profoundly affects microstructure:

| Feature | Slow Cooling (sand cast) | Fast Cooling (die cast) |
|---------|--------------------------|-------------------------|
| Grain size | Large (mm) | Fine (μm) |
| DAS | 200–500 μm | 20–50 μm |
| Segregation extent | Severe | Mild |
| Porosity | Shrinkage dominant | Gas pores possible |
| Phases formed | Equilibrium | Metastable phases possible |
| Typical yield strength | Lower | Higher |

**Very rapid solidification (splat quenching, atomization):** > 10⁶ °C/s → amorphous metals (metallic glasses) possible, extremely fine nanocrystalline structures. Used for transformer cores (Fe-Si metallic glass), AM powder production.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Nucleation | Critical radius r* = 2γT_m/(L_f × ΔT); small undercooling needed for heterogeneous nucleation |
| Dendrites | Tree-like growth along preferred crystal directions; secondary DAS controls diffusion distances |
| Constitutional supercooling | Solute-enriched liquid ahead of front has depressed liquidus → drives dendritic morphology |
| Segregation | k < 1 elements enrich interdendritic region; Scheil equation; removed by homogenization |
| Porosity | Shrinkage (volume decrease on freezing) and gas (dissolved gases expelled); cured by HIP |
| Hot cracking | Thermal stresses tear thin liquid films in late-stage mushy zone |
| Cooling rate | Faster → finer DAS → less segregation → better properties |

---

## Exercises

1. An alloy has T_m = 1,450°C, latent heat L_f = 280 kJ/mol, and solid-liquid interfacial energy γ_SL = 0.20 J/m². Calculate the critical nucleus radius r* for undercoolings of: (a) 1°C (typical heterogeneous), (b) 100°C (homogeneous nucleation range). Convert answers to nanometers.

2. An aluminum casting has DAS = 150 μm in the center and DAS = 30 μm at the surface. If diffusion distance scales as L = √(D × t) with D_Al at homogenization T = 10⁻¹² m²/s, how long does homogenization at 500°C need to be to eliminate microsegregation in: (a) center (half-spacing = 75 μm), (b) surface (half-spacing = 15 μm)?

3. An alloy has C_0 = 4 wt% Cu, k_Cu = 0.17 for Al-Cu alloy. Using the Scheil equation C_s = k × C_0 × (1 - f_s)^(k-1): calculate C_s at f_s = 0, 0.5, 0.9, and 0.99. At what f_s does C_s reach the eutectic composition (33 wt% Cu)? This marks where eutectic phases begin forming in the interdendritic regions.

4. A steel casting has a hot spot (last-to-freeze region) of volume 50 cm³. Steel shrinks 3% on solidification (density 7,000 → 7,870 kg/m³ roughly). (a) What is the volume of shrinkage cavity that forms if no feeding occurs? (b) If a riser is designed to feed this shrinkage, what minimum riser volume is needed assuming the riser itself has 10% waste metal?

5. Explain why rapid solidification (die casting) produces better mechanical properties than slow solidification (sand casting) for the same aluminum alloy. Consider at least three microstructural differences and their effects on: (a) yield strength (Hall-Petch + precipitation), (b) ductility (segregation effects), (c) fatigue life (surface porosity).
