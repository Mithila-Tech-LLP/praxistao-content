# Chapter 48: Directional Solidification — The Bridgman Process

> **"The most important advance in turbine blade manufacturing in the twentieth century was not a new alloy — it was a new way to control the direction in which that alloy solidifies. By making the solidification front move in exactly one direction, engineers eliminated the most dangerous microstructural feature in high-temperature creep: the transverse grain boundary."**

---

## Table of Contents

1. [Why Solidification Direction Matters](#1-why-solidification-direction-matters)
2. [Conventional Casting — Equiaxed Structure](#2-conventional-casting--equiaxed-structure)
3. [Principles of Directional Solidification](#3-principles-of-directional-solidification)
4. [The Bridgman Furnace — Detailed Operation](#4-the-bridgman-furnace--detailed-operation)
5. [The Solidification Front — What It Looks Like](#5-the-solidification-front--what-it-looks-like)
6. [Constitutional Supercooling — The Enemy of DS](#6-constitutional-supercooling--the-enemy-of-ds)
7. [Thermal Gradient Engineering](#7-thermal-gradient-engineering)
8. [The DS Microstructure — What You Get](#8-the-ds-microstructure--what-you-get)
9. [Radiation vs. Liquid Metal Cooling (LMC)](#9-radiation-vs-liquid-metal-cooling-lmc)
10. [DS vs. Conventional — Property Comparison](#10-ds-vs-conventional--property-comparison)
11. [Limitations of DS — Why We Needed SX](#11-limitations-of-ds--why-we-needed-sx)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Why Solidification Direction Matters

When metal solidifies in a mold, it forms grains — small crystals of random size and orientation. The **grain boundaries** between these crystals have fundamentally different properties from the grain interiors.

At room temperature, grain boundaries strengthen a metal (Hall-Petch effect). At high temperature (T_H > 0.5), grain boundaries become weak links:
- Grain boundary sliding (GBS)
- Intergranular creep cavity growth
- Preferential oxidation penetration along boundaries
- Fatigue cracks nucleate preferentially at boundaries

The critical insight: **not all grain boundaries are equally dangerous under centrifugal loading**. A turbine blade's dominant stress is centrifugal — along the blade length (radial direction). 

```
Turbine blade loading:

     ↑ centrifugal stress (primary load direction = blade axis)
     
     ┌─────────────────┐
     │  BLADE AIRFOIL  │
     │                 │
     │                 │
     └────────┬────────┘
              │
              ↓ blade root attachment to disk
```

**Transverse grain boundaries** (perpendicular to blade axis) → directly pulled apart by centrifugal load → GBS opens cavities → catastrophic.

**Longitudinal grain boundaries** (parallel to blade axis) → subjected to shear, not opening → much less damaging.

Conclusion: if you can make all grain boundaries longitudinal (parallel to blade axis), you eliminate the catastrophic failure mode while still having a polycrystalline material.

This is exactly what directional solidification achieves.

---

## 2. Conventional Casting — Equiaxed Structure

In **conventional investment casting** of a turbine blade:
1. The ceramic mold is preheated to above the alloy liquidus
2. Molten metal is poured in at ~1500°C
3. The mold is placed in a cooling environment
4. Solidification begins simultaneously at many points on the mold walls

Because the mold is preheated uniformly and cooling occurs from all surfaces, nucleation occurs everywhere → many grains → completely random orientations → **equiaxed** microstructure.

```
Equiaxed casting grain structure:

  ╔═══════════════════════╗
  ║  ◇  ◈  ◆  ◇  ◈  ◆   ║
  ║  ◈  ◆  ◇  ◈  ◆  ◇   ║← equiaxed grains
  ║  ◆  ◇  ◈  ◆  ◇  ◈   ║  random size, random orientation
  ║  ◇  ◈  ◆  ◇  ◈  ◆   ║  transverse AND longitudinal boundaries
  ╚═══════════════════════╝
```

Result:
- Grain size: 1–10 mm (large, because slow cooling of a massive mold)
- Grain orientation: completely random
- Grain boundaries: present in all orientations including transverse
- Creep life: limited by transverse grain boundary failure

This was the standard from the 1940s through the 1960s — and it was hitting a performance wall.

---

## 3. Principles of Directional Solidification

For directional solidification to produce columnar grains all aligned with the blade axis, three conditions must be maintained simultaneously throughout the casting:

**Condition 1 — Unidirectional heat flow:**
Heat must flow only in ONE direction — along the blade length (axial direction). No radial heat loss. Any radial cooling would nucleate grains on the mold walls with random orientations.

**Condition 2 — Steep thermal gradient at the solidification front:**
The temperature must drop steeply at the front. Above the front: metal is liquid (above liquidus). Below the front: metal is solid. The transition must be sharp — if the "mushy zone" (solid + liquid region) is too wide, secondary nucleation can occur within it.

**Condition 3 — Controlled withdrawal rate:**
The solidification front must move at a controlled rate — slow enough to maintain the steep gradient, fast enough to be economical. The front must sweep from the cold end (blade root) toward the hot end (blade tip) in ONE direction.

```
Directional Solidification — schematic progression:

t=0:         t=5min:       t=15min:      t=30min:
  HOT          HOT           HOT           HOT
  ▓▓▓▓▓▓       ▓▓▓▓▓▓        ▓▓▓▓▓▓        ▓▓▓▓▓▓
  ▓▓▓▓▓▓       ▓▓▓▓▓▓        ▓▓▓▓▓▓        ▓▓▓▓▓▓ ← all liquid (hot zone)
  ▓▓▓▓▓▓       ▓▓▓▓▓▓        ▓▓▓▓▓▓        ───── solidification front
  ▓▓▓▓▓▓       ▓▓▓▓▓▓        ───── front   ∥∥∥∥∥∥
  ▓▓▓▓▓▓       ─────  front  ∥∥∥∥∥∥        ∥∥∥∥∥∥ ← columnar grains (solid)
  ▓▓▓▓▓▓       ∥∥∥∥∥∥        ∥∥∥∥∥∥        ∥∥∥∥∥∥
  COLD         COLD          COLD          COLD
  chill plate  chill plate   chill plate   chill plate
```

The columnar grains (∥) all grow along [001] — the preferred growth direction for FCC metals — which aligns with the blade axis. No transverse boundaries form.

---

## 4. The Bridgman Furnace — Detailed Operation

The standard industrial method is the **Bridgman furnace** (named after Percy Bridgman, Nobel laureate 1946, though his original use was pressure science — the DS application was developed by others in the 1960s):

```
Complete Bridgman Furnace Cross-Section:

┌─────────────────────────────────────────────────┐
│  FURNACE CHAMBER (Vacuum < 10⁻³ Pa)             │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │  Radiation heaters (resistance/induction)│   │
│  │  Temperature: 1500–1600°C               │   │
│  │  Uniformity: ±5°C across zone           │   │
│  │  ┌─────────────────────────────────┐   │   │
│  │  │  Ceramic investment casting mold│   │   │
│  │  │  (blade shape + selector/seed)  │   │   │
│  │  │  Filled with liquid superalloy  │   │   │
│  │  └────────────┬────────────────────┘   │   │
│  └───────────────┼─────────────────────────┘   │
│                  │                              │
├──────────────────┼──────────────────────────────┤ ← BAFFLE PLATE
│                  │  Radiation shadow zone       │   (critical component)
│                  │  (already solidified blade)  │
│                  │                              │
│  Water-cooled withdrawal rod                    │
│  Motor-driven, precision rate control           │
│  ↓ (withdrawing downward)                       │
└─────────────────────────────────────────────────┘
```

### The Baffle Plate — Most Critical Component

The baffle plate (sometimes called a "radiation shield" or "susceptor") is a horizontal plate with a hole for the mold to pass through. It is:
- Made of refractory material (BN, graphite)
- Maintained at approximately the alloy liquidus temperature
- Positioned at the height where solidification should occur

**Its function:** Create a sharp transition from hot zone (above) to cold zone (below). Metal above the baffle: maintained liquid. Metal below the baffle: exposed to radiation from the cold environment → solidifies.

The sharpness of this transition determines the thermal gradient G (°C/mm). A well-designed baffle can achieve G = 40–70°C/mm. The solidification front sits AT the baffle position and moves upward relative to the mold (i.e., the mold moves downward through the fixed baffle).

### Process Sequence

1. **Mold loading:** Assembled ceramic shell + core placed on withdrawal rod
2. **Pump down:** Chamber evacuated to < 10⁻³ Pa (prevents reactive element oxidation)
3. **Heat up:** Heaters bring mold to above liquidus temperature (typically 1500–1600°C for Ni superalloys) → all metal melts
4. **Superheat soak:** Hold at temperature 30–60 minutes → complete melting, temperature uniformity
5. **Withdrawal:** Motor drives mold downward through the baffle at constant rate (2–10 mm/min)
6. **Solidification propagation:** Solidification front is stationary in space; the mold moves down → front appears to move upward through the mold
7. **End of withdrawal:** Entire blade solidified; final solidification at blade tip
8. **Cool down:** Allow to reach 200°C before breaking vacuum

Total cycle time: 2–5 hours depending on blade size and withdrawal rate.

---

## 5. The Solidification Front — What It Looks Like

The solidification front is NOT a sharp plane — it has structure:

```
Solidification Front Detail (cross-section through blade wall):

LIQUID (above liquidus)
─────────────────────────────  T = T_liquidus

  ╔═╗  ╔═╗  ╔═╗  ╔═╗  ← dendrite tips (first to solidify)
  ║ ║  ║ ║  ║ ║  ║ ║
  ║ ║  ║ ║  ║ ║  ║ ║  ← primary dendrite arms growing upward
  ║╔╝  ║╔╝  ║╔╝  ║╔╝  ← secondary arms growing sideways
  ║║   ║║   ║║   ║║
MUSHY ZONE (solid + liquid mixture)
  ║║   ║║   ║║   ║║
  ██   ██   ██   ██  ← fully solidified columnar dendrites

SOLID (below solidus)
```

The distance from the dendrite tips (liquidus temperature) to the fully solidified base (solidus temperature) is the **mushy zone length**:
```
Mushy zone length = ΔT_0 / G
```

Where ΔT_0 is the liquidus-solidus temperature interval and G is the thermal gradient.

For CMSX-4: ΔT_0 ≈ 50°C, G = 50°C/mm → mushy zone ≈ 1 mm
For IN100 (equiaxed alloy): ΔT_0 ≈ 80°C, G = 10°C/mm → mushy zone ≈ 8 mm

Shorter mushy zone = less time for secondary nucleation = better DS quality.

---

## 6. Constitutional Supercooling — The Enemy of DS

**Constitutional supercooling** is the most important concept for understanding DS quality. It explains why stray grains form and what limits the thermal gradient and withdrawal rate.

### What Causes It

As the dendrite tips advance into the liquid, they reject solute elements (those with partition coefficient k < 1) into the liquid ahead of them. This solute-enriched liquid has a **lower liquidus temperature** than the nominal alloy.

```
Liquidus temperature profile ahead of the solidification front:

Temperature    Actual temperature profile (imposed by furnace)
    |          /
    |         /← slope = G (thermal gradient)
T_L ──────────/──────────────────  ← liquidus of nominal alloy
    |        ●← solidification front
    |       /|
    |      / | ← THIS REGION: actual T > T_liquidus(local)
    |     /  |   local liquidus depressed by solute enrichment
T_L'────/──────────────────────   ← liquidus of solute-enriched liquid
    |  /   ← distance ahead of front →
```

If the actual temperature drops BELOW the local liquidus temperature (T_liquidus of the solute-enriched liquid), that region is **constitutionally supercooled** — it's below its melting point, even though it's above the nominal alloy's liquidus. This region can spontaneously nucleate new grains → stray grain formation → DS failure.

### The Stability Criterion

For a stable planar (or columnar) solidification front without constitutional supercooling:

```
G/V ≥ ΔT_0 / D_L

where:
  G = thermal gradient at the front (°C/mm)
  V = withdrawal rate (mm/min)
  ΔT_0 = liquidus-solidus interval (°C)
  D_L = diffusion coefficient of solute in liquid (mm²/min)
```

This is the fundamental design equation for DS casting. If G/V is too low (low gradient or too fast withdrawal), constitutional supercooling occurs → stray grains → rejected blade.

**Practical implications:**
- For a given alloy (fixed ΔT_0 and D_L), increasing G allows higher V (faster, more economical)
- Alloy design can help: lower ΔT_0 (narrower freezing range) → easier DS
- This is why DS superalloy compositions differ from wrought alloy compositions

---

## 7. Thermal Gradient Engineering

Achieving high G is the primary challenge in DS furnace design.

**Standard radiation Bridgman:** G ≈ 20–50°C/mm
**High-gradient Bridgman (optimized baffle):** G ≈ 40–70°C/mm
**Liquid metal cooling (LMC):** G ≈ 80–150°C/mm

Higher G provides:
1. Shorter mushy zone → less constitutional supercooling risk
2. Finer dendrite arm spacing (λ = C × G^(-0.5) × V^(-0.25)) → less microsegregation → better homogenization in heat treatment
3. Higher allowable withdrawal rate V (from G/V ≥ ΔT_0/D_L) → faster cycle time → lower cost

**Dendrite arm spacing (λ):**

Finer λ means:
- Less distance for elements to diffuse during homogenization heat treatment
- Shorter homogenization time needed
- More uniform microstructure after heat treatment → better mechanical properties

Typical λ in DS blades:
- Standard Bridgman: 200–400 μm
- High-gradient: 100–200 μm  
- LMC: 50–100 μm

---

## 8. The DS Microstructure — What You Get

After DS casting (before heat treatment), the microstructure contains:

### Columnar Grains

All elongated along [001] direction (parallel to blade axis). Grain diameter: 0.5–3 mm. Grain length: the full blade height (100–200 mm). All grain boundaries are longitudinal.

```
DS blade microstructure (longitudinal section):

∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥ ← grain 1 boundary
║                     ║
║    ← grain 1 →      ║ ← grain 1 (columnar, [001] up)
║                     ║
∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥
║                     ║
║    ← grain 2 →      ║ ← grain 2 (columnar, [001] up)
║                     ║
∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥∥

All boundaries are LONGITUDINAL (∥ to blade axis)
No transverse boundaries!
```

### Dendritic Microsegregation

As described in Chapter 49, the as-cast structure has chemical segregation on the dendrite arm spacing scale. DS blades require subsequent solution heat treatment to homogenize this segregation.

### As-Cast γ′ Distribution

Large, irregular γ′ precipitates in the interdendritic regions. Not the fine, uniform cuboidal γ′ needed for creep resistance. Heat treatment dissolves and re-precipitates γ′ in the desired morphology.

### Grain Orientation Spread

Individual columnar grains in DS blades are near-[001], but there is a spread (typically 5–15°). Because multiple grains are present, there is random secondary orientation (rotation around [001]) — this is different from SX where secondary orientation can be controlled.

---

## 9. Radiation vs. Liquid Metal Cooling (LMC)

The two Bridgman variants differ in how the solidified portion of the blade is cooled:

### Standard (Radiation-Based) Bridgman

As the mold exits the hot zone through the baffle, the solidified blade cools by radiation from the mold surface to the cold furnace walls. The cooling rate depends on the emissivity difference and view factor.

**Limitation:** Heat removal by radiation is proportional to T⁴ differences — at very high temperatures (just below solidification), radiation is reasonably efficient. As the blade cools, radiation becomes less effective. G is limited by the furnace geometry.

### Liquid Metal Cooling (LMC)

In LMC, the solidified blade is withdrawn into a bath of LIQUID metal (usually tin, Sn, at ~250°C, or aluminum):

```
LMC Bridgman:

  ┌─────────────────────────────────────────┐
  │  HOT ZONE (> liquidus)                 │
  │  ┌───────────────────────────────────┐ │
  │  │      Superalloy mold              │ │
  │  │      (liquid alloy inside)        │ │
  │  └───────────────┬───────────────────┘ │
  └──────────────────┼──────────────────────┘
  ═══════════════════╪════════════════════════ BAFFLE
  ┌──────────────────┼──────────────────────┐
  │  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
  │  ░░░░░Liquid Sn (~250°C) BATH░░░░░░░░░ │ ← mold immersed in molten Sn
  │  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │   dramatically higher G
  └─────────────────────────────────────────┘
  
  Mold withdraws downward into the liquid Sn bath
```

**Why LMC achieves higher G:**
- Liquid Sn has MUCH higher thermal conductivity (62 W/m·K) than gas/vacuum
- Direct contact heat transfer coefficient: ~1000–5000 W/m²K (vs. ~50 W/m²K for radiation at 1000°C)
- Cooling is 10–100× more effective → steeper gradient → G = 80–150°C/mm

**LMC advantages:**
- G up to 150°C/mm → can run faster V for same quality
- Finer dendrite arm spacing → less segregation → better properties
- Better SX quality (fewer stray grains)
- Particularly important for complex blade geometries (platforms, shrouds) where local geometry causes cold spots that can trigger stray grain nucleation

**LMC challenges:**
- Tin contamination of the alloy if the ceramic shell cracks → blade scrapped
- Process control more complex (Sn bath temperature, immersion depth)
- Equipment cost higher
- Sn must be prevented from solidifying on the blade → bath must be precisely maintained

LMC is the current state of the art for highest-quality SX blades.

---

## 10. DS vs. Conventional — Property Comparison

The improvement from conventional casting (CC) to directional solidification (DS) is substantial:

| Property | Conventional Casting | Directional Solidification | Change |
|----------|---------------------|---------------------------|--------|
| Creep life (1000°C, 150 MPa) | ~50h | ~100h+ | +100% |
| Equivalent temperature capability | reference | +25–30°C | +25–30°C |
| Fatigue strength (LCF, 900°C) | reference | +20–30% | Significant |
| Fracture mode at high T | Intergranular | Intergranular (longitudinal) | Changed character |
| Ductility at high T | Low | Higher | Better |
| Thermal fatigue resistance | Reference | +40% | Improved |

The creep improvement is directly attributable to eliminating transverse grain boundaries. Every other property improvement follows from this.

**Historical data point:** The first DS engine (JT9D, 1966) could run 30°C hotter than its predecessor using the SAME alloy, simply because of the microstructure change. This translates directly to better fuel efficiency and more thrust.

---

## 11. Limitations of DS — Why We Needed SX

Despite the major improvement over equiaxed casting, DS still has limitations:

### Longitudinal Grain Boundaries Still Present

DS eliminates transverse boundaries, but longitudinal ones remain. These can still:
- Slide under off-axis loads (real blades experience bending, vibration, not pure axial loads)
- Act as fatigue initiation sites under thermal cycling
- Provide oxidation pathways

### Multiple Grains = Crystallographic Variability

Different DS grains have slightly different orientations. The elastic modulus, creep resistance, and fatigue behavior vary slightly grain-to-grain. This creates internal stress concentrations at grain boundaries under thermal loading.

### Grain Boundary Strengtheners Still Needed

DS alloys still need B, C, Zr to strengthen the remaining grain boundaries. These elements limit the solidus temperature and depress the heat treatment window.

### The Next Step

By 1968, researchers at Pratt & Whitney recognized that if you could reduce the number of grains from "many columnar" to "one," you'd gain another step change in creep life. The grain selector makes this possible.

The transition from DS to SX (covered in Ch 47 and Ch 49) is conceptually simple but technically demanding. It required another 14 years of development (1968 lab demo → 1982 commercial service) to master the defect control required for production SX blades.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| DS principle | Unidirectional heat flow → columnar grains ∥ blade axis → no transverse grain boundaries |
| Bridgman furnace | Heater zone + baffle plate + controlled withdrawal; vacuum environment |
| Baffle plate | Sharp hot/cold transition; determines thermal gradient G |
| Thermal gradient G | High G required to prevent constitutional supercooling ahead of front |
| Constitutional supercooling | Solute-enriched liquid ahead of front has depressed T_liquidus → can nucleate stray grains |
| Stability criterion | G/V ≥ ΔT₀/D_L; governs maximum withdrawal rate for given gradient |
| Mushy zone | ΔT₀/G; shorter = better = less secondary nucleation risk |
| Dendrite arm spacing λ | Decreases with higher G and lower V; finer = less segregation |
| LMC | Liquid Sn bath instead of radiation cooling; G up to 150°C/mm; current state of art |
| DS vs. CC | +25–30°C temperature capability, +100% creep life at equivalent conditions |
| DS limitations | Longitudinal boundaries remain; grain boundary chemistry elements needed; leads to SX |

**Next chapter:** The ultimate step — adding a grain selector to the DS process to produce a single crystal with zero grain boundaries.

---

## Exercises

1. A Bridgman furnace has G = 40°C/mm and V = 4 mm/min. For CMSX-4 with ΔT₀ = 50°C and D_L = 5×10⁻³ mm²/min, does this satisfy the constitutional supercooling criterion G/V ≥ ΔT₀/D_L? What is the maximum withdrawal rate at this gradient?

2. Liquid metal cooling achieves G = 120°C/mm. For the same alloy as question 1, what maximum withdrawal rate is possible? By what factor does LMC improve throughput (parts per unit time) over the standard Bridgman at G = 40°C/mm?

3. Dendrite arm spacing follows λ = C × G^(-0.5) × V^(-0.25) where C is a material constant. If standard Bridgman (G=40°C/mm, V=4 mm/min) gives λ = 300 μm, what is λ for: (a) LMC with G=120°C/mm at the same V, (b) LMC with G=120°C/mm and V=12 mm/min?

4. A DS blade has 5 columnar grains across its width (20 mm wide, ~4 mm per grain). A conventional casting of the same blade has approximately 10 equiaxed grains across the width. If we define "transverse boundary density" as the number of transverse boundaries per mm of blade length, estimate this density for the equiaxed vs DS cases. (For equiaxed: grain aspect ratio ≈ 1, so there are as many transverse as longitudinal boundaries. For DS: all boundaries are longitudinal → zero transverse boundaries.)

5. An alloy designer wants to create a new DS alloy with a narrower freezing range (lower ΔT₀) to allow faster withdrawal. Currently ΔT₀ = 60°C. If they reformulate to ΔT₀ = 40°C, by what factor can they increase withdrawal rate while maintaining the same thermal gradient and constitutional supercooling stability margin?
