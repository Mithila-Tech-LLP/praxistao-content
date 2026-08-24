# Chapter 57: Thermal Barrier Coatings — The Ceramic Shield

> **"Yttria-stabilized zirconia is perhaps the most important ceramic in modern engineering. A 0.1mm layer of it separates the burning 1,600°C combustion gas from the 1,000°C metal. It is deposited as a rough, columnar, nano-porous microstructure that no one would design on paper — yet this seemingly imperfect structure gives it extraordinary resistance to thermal fatigue. When it fails, it fails suddenly, and the blade it protected dies in seconds."**

---

## Table of Contents

1. [What Is a TBC and Why Is It Needed?](#1-what-is-a-tbc-and-why-is-it-needed)
2. [The Two-Layer System: Bond Coat + Top Coat](#2-the-two-layer-system-bond-coat--top-coat)
3. [Yttria-Stabilized Zirconia (YSZ) — Why This Material?](#3-yttria-stabilized-zirconia-ysz--why-this-material)
4. [Deposition Methods: EB-PVD vs. APS](#4-deposition-methods-eb-pvd-vs-aps)
5. [EB-PVD Microstructure — Columnar Columns](#5-eb-pvd-microstructure--columnar-columns)
6. [APS Microstructure — Splats and Lamellae](#6-aps-microstructure--splats-and-lamellae)
7. [Thermally Grown Oxide (TGO) — The Critical Interface](#7-thermally-grown-oxide-tgo--the-critical-interface)
8. [TBC Failure Mechanisms](#8-tbc-failure-mechanisms)
9. [CMAS Attack — The Modern Threat](#9-cmas-attack--the-modern-threat)
10. [Next-Generation TBC Materials](#10-next-generation-tbc-materials)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Is a TBC and Why Is It Needed?

A **Thermal Barrier Coating (TBC)** is a multi-layer coating system applied to the exterior surface of turbine blades (and vanes) to:
1. **Reduce heat flux** into the blade metal → lower metal temperature → longer creep/fatigue life
2. **Allow higher turbine inlet temperature** at same metal temperature → better engine efficiency
3. **Reduce cooling air requirements** → better engine efficiency

**System of systems:**

The TBC does not work alone — it is part of a complete **thermal protection system (TPS)**:
```
Gas stream (1,600°C)
    ↓ heat flux
TBC top coat (YSZ ceramic, ~0.1–0.3 mm)       ← insulation layer
    ↓
Bond coat (MCrAlY or Pt-Al, ~0.05–0.15 mm)    ← oxidation protection + adhesion
    ↓
Thermally grown oxide (TGO, α-Al₂O₃, ~0–10 μm, grows during service)
    ↓
Superalloy substrate (SX blade, ~2–3 mm)       ← structural member
    ↓ heat flux
Internal cooling (600–700°C air)
```

**Without TBC**, blade metal temperature ≈ 1,100°C at current TIT levels.
**With TBC**, metal temperature ≈ 950–1,000°C — a savings of 100–150°C.

This temperature difference is CRITICAL:
- Creep rate scales as exp(-Q/RT) — at 950°C vs. 1,100°C, creep rate is 10–100× different
- Oxidation rate: 150°C lower → 3–10× slower oxidation → dramatically longer coating life

---

## 2. The Two-Layer System: Bond Coat + Top Coat

### Bond Coat

The bond coat serves two functions:
1. **Adhesion**: bridges the chemical and thermal expansion mismatch between the ceramic top coat and the alloy substrate
2. **Oxidation protection**: protects the superalloy from oxidation (forms its own protective α-Al₂O₃ scale — the thermally grown oxide, TGO)

**Two types of bond coats:**

**MCrAlY:** Metallic alloy (M = Ni, Co, or NiCo; Cr = 15–25%; Al = 8–12%; Y = 0.5%)
- Y (yttrium) pegs the Al₂O₃ scale → prevents scale spallation → much better scale adhesion
- Applied by: Low-pressure plasma spray (LPPS) or high-velocity oxy-fuel (HVOF)
- Thickness: 100–200 μm

**Pt-aluminide:** Nickel aluminide (NiAl) with platinum addition
- Formed by: electroplating Pt on blade surface, then pack cementation aluminizing at ~1000°C → interdiffusion forms (Ni,Pt)Al intermetallic
- Pt: dramatically improves scale adhesion (by segregating to oxide grain boundaries)
- Higher temperature capability than MCrAlY for oxidation
- Used on highest-temperature blades (1st stage HPT)

### Top Coat — YSZ

Yttria-stabilized zirconia (Y₂O₃-stabilized ZrO₂):
- Thickness: 100–300 μm on blade airfoils; up to 600 μm on vanes
- Thermal conductivity: 2.1 W/mK (20–30% of bond coat, 5× less than alloy)
- Applied by: EB-PVD (blades) or APS (vanes, combustors)

---

## 3. Yttria-Stabilized Zirconia (YSZ) — Why This Material?

Pure zirconia (ZrO₂) undergoes a phase transformation on cooling that makes it useless as a thermal barrier:
- At 2,370°C: Cubic → tetragonal transition
- At 1,170°C: Tetragonal → monoclinic transition (3–5% volume expansion)

The monoclinic transformation occurs on COOLING (during engine shutdown) and causes ~5% volume change → cracks → spallation. Unacceptable.

**Yttria stabilization:**

Adding 6–8 mol% Y₂O₃ stabilizes the tetragonal phase (called "non-transformable tetragonal" or t′ phase):
- Y³⁺ substitutes for Zr⁴⁺ → charge compensated by oxygen vacancies → modified phase stability
- The t′ phase is metastable below ~1,200°C but does NOT transform to monoclinic on cooling → no volume change → survives thermal cycling
- The metastability is maintained by the kinetics — the transformation is suppressed at temperatures below ~1,200°C

**Why 7 wt% Y₂O₃ specifically?**

```
ZrO₂-Y₂O₃ phase diagram:

Y₂O₃ (mol%):  0    2    4    7    8    15
Phase:          M    M+T  T'   T'  C+T'  C
                ↑unstable ↑OPTIMAL↑ stable
                monoclinic      cubic
                phase change
                on cooling

T' = "non-transformable tetragonal" — the target phase
M = monoclinic — bad (phase transformation on cooling)
C = cubic — stable BUT much lower thermal cycle life (no t' toughening)
```

7 wt% Y₂O₃ (= ~4.5 mol%) gives the t′ phase with the best combination of:
- No monoclinic transformation on cooling
- High strain tolerance (t′ can absorb stress by stress-induced transformation — "transformation toughening")
- Good phase stability up to ~1,200°C service

This is the "standard" TBC composition used in essentially all production turbine blades since the 1970s.

---

## 4. Deposition Methods: EB-PVD vs. APS

Two completely different deposition methods produce completely different microstructures with different properties:

### Electron Beam Physical Vapor Deposition (EB-PVD)

**Process:**
1. A ZrO₂+Y₂O₃ ceramic ingot is heated by a focused electron beam (>10 kW)
2. The ceramic evaporates into vapor
3. The rotating, heated blade is positioned in the vapor cloud
4. Ceramic condenses on the blade surface
5. High vacuum environment (10⁻³ Pa)

**Result:** **Columnar microstructure** (columns perpendicular to blade surface)

**Processing conditions:**
- Substrate temperature: 900–1,000°C (for good adhesion and columnar structure)
- Deposition rate: 5–20 μm/min
- Typical TBC thickness: 100–200 μm
- Deposition time per blade: 30–60 minutes

### Atmospheric Plasma Spray (APS)

**Process:**
1. Plasma gun generates plasma jet at 8,000–15,000°C
2. YSZ powder (20–100 μm) injected into plasma jet → melted into liquid droplets
3. Molten droplets impact blade surface at 100–300 m/s
4. Each droplet spreads and quenches instantly → "splat"
5. Layer builds up from millions of overlapping splats

**Result:** **Lamellar (layered) microstructure** with significant porosity

**Processing conditions:**
- Open air process (no vacuum required)
- Deposition rate: 50–200 μm/min
- Typical TBC thickness: 150–400 μm
- Much faster and cheaper than EB-PVD

---

## 5. EB-PVD Microstructure — Columnar Columns

The EB-PVD TBC microstructure is one of the most remarkable in all of engineering:

```
EB-PVD TBC cross-section (SEM view, schematic):

Gas stream surface (top):
  ─────────────────────────────────────────────────────
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │  ← columnar grains
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │    (~1–5 μm wide)
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │
  │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │ │
  ─────────────────────────────────────────────────────
TGO + bond coat interface (bottom, ~100 μm below top)

Column width: 1–5 μm
Column height: 100–200 μm (full TBC thickness)
Inter-column gap: ~0.1–0.5 μm (open or closed)
```

### Why Columnar Microstructure Is Superior for Thermal Cycling

The blade's thermal expansion coefficient:
- Superalloy substrate: α_alloy ≈ 13.5×10⁻⁶/K
- YSZ TBC: α_YSZ ≈ 10.7×10⁻⁶/K

This mismatch of ~3×10⁻⁶/K means that during thermal cycling:
- Heating: alloy expands more than TBC → TBC in COMPRESSION
- Cooling: alloy contracts more than TBC → TBC in TENSION

The induced strain per cycle:
```
Δε = Δα × ΔT = 3×10⁻⁶/K × 700K = 0.0021 (0.21%)
```

This cyclic strain would crack and spall a solid, rigid TBC. The columnar EB-PVD structure accommodates this strain via:
- **Inter-column gaps** open and close during cycling (columns slide laterally relative to each other) → strain accommodated without building up in-plane stress
- Individual columns are mechanically compliant in the lateral direction

This is why EB-PVD TBC survives 30,000+ thermal cycles while a dense, rigid TBC would spall after dozens of cycles.

**Feathery nano-porosity within columns** also helps: nano-pores within each column add compliance AND reduce thermal conductivity below the fully dense value.

---

## 6. APS Microstructure — Splats and Lamellae

APS TBC has a completely different structure:

```
APS TBC cross-section (schematic):

Gas stream surface (top):
  ════════════════════════════════════  ← splat (one impacted droplet)
  ════════════════════════════════════
  ════════════════════════════════════  ← layered "lamellar" structure
  ════════════════════════════════════
  ════════════════════════════════════  ← cracks between splats (horizontal)
  ════════════════════════════════════
  ════════════════════════════════════  ← porosity (1–15%, depending on spray params)
  ════════════════════════════════════
Bond coat interface (bottom)
```

**APS characteristics:**
- Porosity: 10–20% (much more porous than EB-PVD)
- Horizontal cracks between splats → delamination planes
- Lower thermal conductivity than EB-PVD (2.1 → 1.5–1.8 W/mK) due to high porosity
- Cheaper, faster to deposit

**Why APS has lower thermal cycling life:**

Horizontal splat boundaries act as delamination crack nucleation sites. When the CTE mismatch drives lateral stress during cycling, cracks grow along the splat boundaries → spallation.

APS TBC typically survives 1,000–3,000 cycles in accelerated burner testing. EB-PVD: 5,000–15,000+ cycles.

**Why APS is still used:**
- Much cheaper than EB-PVD ($50 vs $500 per blade)
- Thicker coatings possible → more temperature reduction
- Good enough for lower-temperature components (LPT blades, combustors, industrial GT vanes)

In aircraft HPT blades: EB-PVD is standard.
In industrial gas turbines and lower-temperature aerospace components: APS.

---

## 7. Thermally Grown Oxide (TGO) — The Critical Interface

During service, the bond coat oxidizes. A thin layer of α-Al₂O₃ (alumina) grows at the bond coat/TBC interface — the **thermally grown oxide (TGO)**:

```
TGO formation:

     YSZ TBC
     ────────────────────────────────────────
     TGO: α-Al₂O₃ (growing slowly)           ← ~0.1 μm thick initially, grows to ~10 μm
     ────────────────────────────────────────
     Bond coat (MCrAlY): Al depletes as TGO grows
     ────────────────────────────────────────
     Superalloy substrate
```

**TGO is essential — and also the life-limiting factor:**

**Why TGO is GOOD:**
- α-Al₂O₃ is the most protective oxide: very low growth rate, dense, adherent
- Protects the bond coat and substrate from further oxidation
- The TGO is a thermodynamic barrier: as long as TGO remains intact and α-Al₂O₃ composition, oxidation is minimal

**Why TGO causes FAILURE:**
- TGO grows by consuming Al from the bond coat
- TGO has its own CTE mismatch with both YSZ and bond coat
- As TGO thickens (0 → 10 μm over 30,000h), it accumulates stress
- When TGO reaches ~5–10 μm thick: **stresses become large enough to initiate cracks** at TGO/TBC interface → TBC delamination → spallation

This is the fundamental life-limiting mechanism of TBC systems:
```
TBC life limited by TGO growth → rate-limited by Al diffusion rate through TGO
```

Factors that accelerate TGO growth (and shorten TBC life):
- Higher temperature → faster Al diffusion → faster TGO growth
- Water vapor → volatile Al(OH)₃ → faster scale consumption
- High cycling rate → more CTE mismatch damage per unit time

**The yttrium role in bond coat:**

Without Y (yttrium) in the bond coat, Al₂O₃ scale grows as metastable θ-Al₂O₃, which has 50% higher growth rate than α-Al₂O₃, and is prone to "rumpling" (wavy interface) → severe spallation.

Y segregates to the oxide grain boundaries → pins them → prevents scale growth by grain boundary diffusion → slower TGO growth rate → longer TBC life.

This is why MCrAlY (with Y) dramatically outperforms MCrAl (without Y) for TBC bond coat.

---

## 8. TBC Failure Mechanisms

### Mechanism 1: TGO-Driven Delamination (Most Common in Service)

As described above: TGO grows → accumulates stress → cracks form between TGO and TBC → delamination → spallation.

Symptoms: Large sheets of TBC peeling off in hot zone. Sudden exposure of bond coat → rapid oxidation.

### Mechanism 2: Rumpling

At very high temperatures (>1,050°C bond coat temperature), the MCrAlY bond coat softens and can undergo plastic deformation from CTE mismatch stresses → the TGO/bond coat interface develops waves (rumples):

```
Before rumpling:         After rumpling:
─────────────────       ─────────────────
─────────────────   →   ~~~~~~~~~~~~~~~~~~~  ← bond coat top surface
═════════════════       ═══════════════════  ← bond coat (deformed)
```

The TBC tries to follow these rumples → additional stress concentrations → cracks at rumple peaks → spallation.

Pt-aluminide bond coats are more resistant to rumpling than MCrAlY (because NiAl is harder → less plastic deformation).

### Mechanism 3: Sintering

At high temperature (>1,150°C TBC service), the nano-porous, columnar YSZ structure **sinters** — the nano-pores close as surface diffusion reduces surface energy:
- Thermal conductivity INCREASES (more fully dense → less insulation → higher metal temperature)
- Compliance DECREASES (columns fuse together → can't accommodate CTE mismatch → spallation)

Sintering is the reason EB-PVD TBC has a practical temperature limit of ~1,150–1,200°C surface temperature.

### Mechanism 4: Phase Transformation at High Temperature

At T > 1,200°C, the metastable t′ phase (7 wt% YSZ) starts slowly transforming toward the equilibrium phases:
```
t′ (at service T > 1,200°C) → tetragonal (t) + cubic (C)
t → monoclinic (m) on COOLING → 5% volume change → cracking
```

This is called **destabilization** of the TBC. It sets an absolute upper temperature limit of ~1,200°C for standard 7YSZ. Above this temperature, YSZ eventually fails catastrophically.

---

## 9. CMAS Attack — The Modern Threat

**CMAS** = Calcium-Magnesium-Alumino-Silicate — the modern existential threat to TBCs.

**Source:** Fine particulates that enter the engine:
- **Volcanic ash** (well-known: Iceland Eyjafjallajökull 2010 → massive flight cancellations)
- **Desert sand** (Middle East operations)
- **Runway debris** (especially severe at low altitude during takeoff/landing)
- **Industrial dust** in power plant installations

At temperatures above ~1,240°C (which modern TBC surfaces exceed), sand and ash **melt into a glass**. This molten CMAS:
1. Wets the TBC surface and penetrates the inter-column gaps in EB-PVD TBC
2. CMAS melt dissolves Y₂O₃ from YSZ → Y is leached out → destabilizes the t′ phase → monoclinic transformation on cooling → cracking
3. CMAS fills the inter-column gaps → after solidification on cooling → the flexible columnar structure becomes RIGID → loses strain tolerance → spallation

**CMAS effect on APS TBC:** Penetrates through micro-cracks into the layered splat structure → seals cracks (removes compliance) → spallation on next thermal cycle.

**CMAS is increasingly important** as:
- TIT temperatures push higher → TBC surface temperatures exceed CMAS melting point more frequently
- Airlines increasingly operate in dust/ash-prone environments

**CMAS resistance strategies:**
- **Sacrificial reaction layers**: Rare-earth zirconates (Gd₂Zr₂O₇, Nd₂Zr₂O₇) react with CMAS preferentially → form a crystalline barrier that blocks further CMAS penetration
- **Dense top layer**: A thin, dense YSZ or alumina cap that prevents CMAS penetration into the columnar structure
- **Alternative TBC chemistries**: Some newer materials have higher crystallization temperature for CMAS interaction

---

## 10. Next-Generation TBC Materials

Standard 7YSZ (7 wt% Y₂O₃ in ZrO₂) has been the workhorse since the 1970s. Its limits:
- Temperature capability: ~1,200°C surface (above this: sintering + phase destabilization)
- Thermal conductivity: ~2.1 W/mK (as-deposited EB-PVD)

Engine TIT continues to push higher. New TBC materials are being developed:

### Rare-Earth Zirconates

**Gadolinium Zirconate (Gd₂Zr₂O₇, GZO):**
- Lower thermal conductivity: 1.5–1.7 W/mK (vs. 2.1 for YSZ)
- Better high-temperature stability: cubic phase stable to melt → no phase transformation issue
- Better CMAS resistance (reacts with CMAS → blocks penetration)
- **Limitation:** More brittle than YSZ → lower fracture toughness → worse thermal cycling life

Typically used as a **dual-layer system**: YSZ (inner, bonded layer) + GZO (outer layer) → best of both.

### Lanthanum Zirconate (La₂Zr₂O₇)

Similar to GZO: lower conductivity, better high-T stability. Less CMAS interaction than GZO.

### Pyrochlore Structures

**La₂Ce₂O₇, (La,Gd)₂(Zr,Hf)₂O₇**: Pyrochlore crystal structure (ordered variant of fluorite). Very low conductivity (1.0–1.5 W/mK). High temperature stability. Current research focus.

### Thermal Conductivity Comparison

| Material | k (W/mK) | T limit (°C) | CMAS resistance |
|----------|----------|-------------|-----------------|
| 7YSZ | 2.1 | ~1,200 | Poor |
| 7YSZ (columnar, nanostructured) | 1.7 | ~1,200 | Poor |
| Gd₂Zr₂O₇ | 1.5–1.7 | ~1,500 | Good |
| La₂Zr₂O₇ | 1.5 | ~1,500 | Moderate |
| 20YSZ | 1.8 | ~1,300 | Better than 7YSZ |
| Pyrochlore oxides | 1.0–1.5 | >1,400 | Good |

---

## Summary

| Topic | Key Point |
|-------|-----------|
| TBC purpose | Reduce metal temperature by 100–200°C; allows higher TIT at same blade life |
| System structure | TBC top coat + bond coat + TGO (grows during service) + superalloy substrate |
| YSZ composition | 7 wt% Y₂O₃ in ZrO₂; stabilizes non-transformable t′ phase to prevent monoclinic transformation |
| EB-PVD | Columnar microstructure; compliant during thermal cycling; 100–200 μm; used on HPT blades |
| APS | Lamellar (layered) microstructure; more porous, lower k; cheaper; lower cycling life; used on vanes/combustors |
| TGO | α-Al₂O₃ grows at bond coat surface during service; life-limiting when ~5–10 μm thick |
| Bond coat yttrium | Y segregates to oxide grain boundaries → slower TGO growth → longer TBC life |
| Failure modes | TGO-driven delamination, rumpling, sintering, phase destabilization |
| CMAS | Molten silicate deposits; infiltrate TBC; leach Y₂O₃; cause catastrophic failure; increasingly important |
| Next-gen TBCs | Rare-earth zirconates (Gd₂Zr₂O₇); lower k + better high-T stability + CMAS resistance |

**Next chapter:** How turbine blade life ends — the failure modes of HPT blades in service, and what the failure analysis reveals about which mechanism dominated.

---

## Exercises

1. Calculate the heat flux reduction achieved by a 0.15 mm TBC vs. no TBC. Boundary conditions: gas temperature 1,600°C, h_gas = 3,000 W/m²K; metal conductivity 11 W/mK, thickness 2.5 mm; h_internal = 2,500 W/m²K, coolant temperature 650°C. Set up the thermal resistance network for both cases and solve for heat flux and metal surface temperature. (Hint: R_total = 1/h_gas + t_TBC/k_TBC + t_metal/k_metal + 1/h_internal)

2. YSZ has α_YSZ = 10.7×10⁻⁶/K and the alloy has α_alloy = 13.5×10⁻⁶/K. During engine startup, the blade heats from 20°C to 1,000°C (ΔT = 980°C). Calculate: (a) the thermal strain mismatch Δε = Δα × ΔT, (b) if the TBC were fully bonded with no gap accommodation, what stress would develop? (Use E_TBC = 50 GPa), (c) how do the inter-column gaps in EB-PVD TBC prevent this stress from accumulating?

3. TGO grows at a parabolic rate: thickness h_TGO(t) = k_p × √t, where k_p = 0.1 μm/h^0.5 at 1,050°C. (a) What is TGO thickness after 100h, 1,000h, and 30,000h of service? (b) TBC fails when TGO reaches 8 μm. At what service time does this occur? (c) If TBC is applied to a blade used 6 hours per flight, how many flights before TBC replacement is needed?

4. CMAS melts at 1,240°C and the TBC outer surface is at 1,150°C during normal operation. The TIT increases from 1,550°C to 1,650°C (100°C increase). If TBC surface temperature scales proportionally with TIT minus metal temperature (which stays constant due to more cooling), estimate the new TBC surface temperature. Does it exceed the CMAS melting point? What are the implications?

5. Compare EB-PVD and APS TBC in terms of: (a) deposition cost, (b) thermal cycling life, (c) thermal conductivity, (d) appropriate applications. Design a blade coating strategy for: (i) a 1st-stage HPT blade that sees 30,000 flights, (ii) a 2nd-stage LPT blade that sees 30,000 flights at 300°C lower temperature, (iii) an industrial gas turbine blade that must last 100,000 hours at steady temperature. Justify your choice for each.
