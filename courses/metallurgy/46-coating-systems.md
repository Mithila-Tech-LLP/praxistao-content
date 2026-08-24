# Chapter 46: Coating Systems — MCrAlY, Aluminides, and TBC Integration

> **"The turbine blade coating system is an engineered system within a system. Three distinct layers — aluminide or MCrAlY bond coat, thermally grown oxide, and YSZ thermal barrier — each with a different composition, structure, and function, must work together for 30,000 hours under cyclic thermal loading. Each layer evolves during service. The system fails when these evolutions become incompatible."**

---

## Table of Contents

1. [The Complete Coating Stack](#1-the-complete-coating-stack)
2. [Historical Development — From Aluminides to TBC Systems](#2-historical-development--from-aluminides-to-tbc-systems)
3. [Bond Coat Types and Selection](#3-bond-coat-types-and-selection)
4. [Bond Coat Deposition — LPPS vs. HVOF](#4-bond-coat-deposition--lpps-vs-hvof)
5. [The TGO — Interface Chemistry During Service](#5-the-tgo--interface-chemistry-during-service)
6. [YSZ TBC — Deposition on Bond Coat](#6-ysz-tbc--deposition-on-bond-coat)
7. [System Evolution During Service — What Changes](#7-system-evolution-during-service--what-changes)
8. [System Life Prediction](#8-system-life-prediction)
9. [Next-Generation Coating Concepts](#9-next-generation-coating-concepts)
10. [Complete Process Sequence — Blade to Engine](#10-complete-process-sequence--blade-to-engine)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Complete Coating Stack

A modern HPT blade coating system consists of:

```
COMPLETE COATING CROSS-SECTION (not to scale):

┌──────────────────────────────────────────────────────────────────┐
│  Gas stream: T_gas = 1,550–1,700°C                              │
└──────────────────────────────────────────────────────────────────┘
                          ↓ heat flux
═══════════════════════════════════════════════════════════════════
  YSZ TBC (EB-PVD):  0.1–0.2 mm     T_surface: 1,050–1,150°C
  Columnar, porous                   k = 2.1 W/mK
═══════════════════════════════════════════════════════════════════
  TGO (α-Al₂O₃):    0–8 μm          T: 950–1,050°C
  Grows during service               Grows parabolicaly
═══════════════════════════════════════════════════════════════════
  Bond coat:         0.05–0.15 mm    T: 900–1,000°C
  MCrAlY or Pt-Al                    Sacrificial Al reservoir
═══════════════════════════════════════════════════════════════════
  Superalloy (SX):   2–3 mm          T_metal: 950–1,050°C
  CMSX-4                             Structure, strength
═══════════════════════════════════════════════════════════════════
  Cooling passages   (inside blade)  T_coolant: 650–750°C
  Hollow, complex geometry
═══════════════════════════════════════════════════════════════════
```

**Each layer's function:**

| Layer | Primary function | Secondary function |
|-------|-----------------|-------------------|
| YSZ TBC | Thermal insulation | Oxidation protection minor at surface |
| TGO | Thermodynamic stability | Chemical bridge TBC ↔ bond coat |
| MCrAlY | Oxidation protection (Al₂O₃) | Hot corrosion resistance (Cr), TBC adhesion |
| Superalloy | Structural strength | — |

---

## 2. Historical Development — From Aluminides to TBC Systems

**1950s–1960s: Aluminide coatings**

First generation: simple diffusion aluminides. Blade surface enriched in Al → NiAl β-phase → forms Al₂O₃ scale. Extends life vs. uncoated blade but limited temperature capability.

**1970s: MCrAlY overlay coatings**

Plasma spray of MCrAlY overcomes the Al limitation of diffusion coatings:
- Higher Al content possible
- Higher Cr content for hot corrosion
- Yttrium for scale adhesion → major spallation resistance improvement

First MCrAlY coatings (Co-Cr-Al-Y, Pratt & Whitney, ~1975) increased blade life 3–5× vs. aluminide alone.

**1975–1985: First TBC systems**

NASA Glenn Research Center (then Lewis) developed early TBC on plasma-sprayed MCrAlY in the mid-1970s. First engine test: Pratt & Whitney F100, 1976. Initial TBCs were thin (~0.1 mm) and APS-deposited. Life was short (hundreds of cycles) but demonstrated the concept.

**1985–1990s: EB-PVD TBC**

Rolls-Royce and MTU demonstrated EB-PVD TBC on turbine blades (superior thermal cycling life). By the mid-1990s, EB-PVD TBC became standard on HPT first-stage blades in all major engine programs.

**2000s–present: Optimization and new materials**

Pt-modified aluminide bond coats become standard on highest-temperature blades. TBC thickness optimization. Initial testing of Gd₂Zr₂O₇ and other novel TBCs. CMAS-resistant coatings under development.

---

## 3. Bond Coat Types and Selection

### MCrAlY Overlay Coatings

**Composition variants:**

| Type | Composition | Best for |
|------|-------------|---------|
| NiCrAlY | Ni-20Cr-10Al-1Y | High oxidation resistance |
| CoCrAlY | Co-25Cr-10Al-0.5Y | Better hot corrosion resistance |
| NiCoCrAlY | Ni-22Co-17Cr-12Al-0.6Y | Best balance; most widely used |
| NiCoCrAlYHf | + 0.2Hf | Scale adhesion + TGO morphology control |

The addition of Hf (hafnium) to MCrAlY (making it MCrAlYHf):
- Hf segregates to α-Al₂O₃ grain boundaries (like Y, but different mechanism)
- Further reduces TGO growth rate
- Reduces TGO rumpling tendency
- Used on highest-performance blade systems

### Diffusion Aluminide Coatings

**Simple NiAl aluminide:**
- Formed by pack cementation or chemical vapor deposition (CVD)
- Inward growing (Ni diffuses out to alloy surface, Al deposits) or outward growing
- 50–100 μm thick β-NiAl layer
- Brittle below ~700°C → ductile above
- Good oxidation resistance but poor hot corrosion (low Cr)
- Used on: cooled vanes, low-temperature blades, compressor blades

**Platinum-modified aluminide (Pt-Al, PtNiAl):**
- Pt electroplated first (5–10 μm) → then aluminized
- Result: (Ni,Pt)Al β phase
- Pt dramatically improves scale adhesion and spallation resistance
- Used as bond coat under TBC on HPT blades in all modern commercial engines
- Pratt & Whitney, GE, Rolls-Royce all use Pt-Al or MCrAlY variants for HPT

### Selection Criteria

| Parameter | MCrAlY | Pt-aluminide |
|-----------|--------|-------------|
| Al reservoir | High (deposited by spray) | Moderate (limited by interdiffusion) |
| Oxidation life | Very long (30,000h+ possible) | Long (20,000–30,000h) |
| Hot corrosion resistance | Excellent (high Cr) | Moderate (lower Cr) |
| TBC adhesion | Excellent | Excellent (with Pt) |
| Processing | LPPS/HVOF spray (adds 100–150 μm) | Electroplate + diffuse (thin, precise) |
| Cost | Moderate | Higher (Pt cost + processing) |
| Rumpling resistance | Good (depends on composition) | Better (NiAl harder than MCrAlY at T) |
| Preferred application | Vanes, LPT blades, industrial GT | HPT blades (first stage, hottest) |

---

## 4. Bond Coat Deposition — LPPS vs. HVOF

### Low-Pressure Plasma Spray (LPPS)

Also called Vacuum Plasma Spray (VPS):

**Process:**
1. Chamber pumped to 50–150 mbar (5–15 kPa) of Ar
2. Plasma gun (Ar/H₂ or Ar/He plasma): 8,000–15,000°C plasma jet
3. MCrAlY powder (15–75 μm diameter, spherical gas-atomized) injected into jet
4. Particles melt → impact blade at 150–300 m/s → quench → adhere

**Advantages:**
- Vacuum prevents oxide inclusions in coating (no air → no O₂)
- Dense coating (< 1% porosity)
- Good adhesion strength (> 50 MPa bond strength)

**Disadvantages:**
- Expensive chamber equipment and maintenance
- Slower deposition (one blade per 10–20 min)
- Line-of-sight limitation (complex blade geometries difficult to coat uniformly)

### High-Velocity Oxy-Fuel (HVOF)

**Process:**
1. Fuel (kerosene, H₂) + O₂ → combustion in nozzle at 3,000°C
2. Supersonic gas jet (1,500–2,000 m/s)
3. MCrAlY powder injected → partially melts → impacts at very high velocity
4. Dense coating formed primarily by mechanical impact energy (kinetic energy = density)

**Advantages:**
- Very high particle velocity → very dense coating (< 0.5% porosity even without vacuum)
- Atmospheric process → lower equipment cost
- High deposition rate

**Disadvantages:**
- Some oxide inclusions (process in air)
- More complex for complex blade geometries (spray angle sensitivity)

Both methods are used in production: LPPS for highest-performance HPT blades, HVOF increasingly competitive due to lower cost with nearly equivalent quality.

---

## 5. The TGO — Interface Chemistry During Service

The thermally grown oxide (TGO) is the most dynamic part of the coating system. It:
- Doesn't exist when the blade is new (or is just a few nm thick from the pre-heat steps)
- Grows during service by Al oxidation at the bond coat surface
- Reaches 5–10 μm after 30,000h
- Accumulates stress as it grows → ultimately causes TBC spallation

### TGO Phase Evolution

The initially formed oxide is NOT α-Al₂O₃:
1. **Stage 1 (first exposure):** θ-Al₂O₃ forms (metastable, needle-like) — fast growing
2. **Stage 2 (hours–days at 1,000°C):** θ-Al₂O₃ → α-Al₂O₃ transformation (irreversible) — this transformation involves ~7% volume contraction → can nucleate microcracks
3. **Stage 3 (long-term service):** Steady α-Al₂O₃ growth (parabolic, slow)

The initial θ → α transformation is potentially damaging if it occurs after TBC is deposited (the volume change cracks the TBC/TGO interface). This is why:
- Blades are pre-oxidized after bond coat deposition (before TBC) to establish α-Al₂O₃
- Pre-oxidation conditions (temperature, time) are carefully controlled

### TGO Stress State

The TGO grows at high temperature and cools with the blade. Because of CTE mismatch:
- At service temperature: TGO in low stress (near zero → slight compression)
- On cooling to room temperature: TGO in biaxial compression (~5 GPa at RT)
- On heating back up: residual stress released → near-zero stress at temperature
- The stress is PATH DEPENDENT (different on heating vs. cooling)

This biaxial compressive stress is the stored elastic energy that drives delamination.

---

## 6. YSZ TBC — Deposition on Bond Coat

### Pre-Deposition Steps

Before TBC deposition:
1. **Grit blast bond coat surface:** Roughen to Ra = 3–5 μm → mechanical interlocking for first TBC atomic layers
2. **Pre-oxidize:** Heat coated blade in air at 1,000–1,050°C for 1–5h → establishes a continuous α-Al₂O₃ TGO layer → provides clean chemical surface for TBC adhesion
3. **Inspect bond coat:** Surface quality verification (no cracks, correct surface chemistry)

### EB-PVD Deposition Conditions

- Substrate temperature: 950–1,000°C (critical for columnar structure formation)
- Deposition rate: 5–15 μm/min
- Target TBC thickness: 100–200 μm
- Rotation: blade rotated continuously → even coating on all faces
- Multiple ingots: separate ingots for ZrO₂ and Y₂O₃ (or pre-alloyed YSZ ingots)

**Column formation mechanism:**

At 950°C substrate temperature, adatoms on the growing TBC surface have sufficient mobility to diffuse laterally → find low-energy nucleation sites → deposit preferentially on already-growing columns → columns grow, gaps between them persist → columnar structure.

At lower T: random deposition → dense, equiaxed → no column structure → no compliance → would spall immediately.

The 950°C substrate temperature is therefore CRITICAL — too cool → wrong structure → product failure.

---

## 7. System Evolution During Service — What Changes

During 30,000 hours of service, every layer in the coating system changes:

**Bond coat:**
- Al decreases (consumed to TGO): 10% → ~5% remaining after 30,000h
- Phase changes: β-NiAl (Al-rich) → γ/γ′ (Al-depleted) → when all β consumed, no more Al reservoir
- Interdiffusion with substrate: bond coat elements diffuse into alloy, alloy elements diffuse into bond coat → "interdiffusion zone" forms (can nucleate harmful phases)
- Some TCP phase formation in severely overtemperature cases

**TGO:**
- Grows from 0 → 8–10 μm (over 30,000h at 1,000°C bond coat temperature)
- Interface morphology changes: initially smooth → "rumpling" develops in some systems
- Chemical composition: pure α-Al₂O₃ → some Cr₂O₃, spinel (NiAl₂O₄) as Al is depleted

**YSZ TBC:**
- Sintering of nano-pores (above 1,200°C TBC temperature) → thermal conductivity increases
- Phase evolution: t′ metastable → starts transforming above 1,200°C sustained temperature
- CMAS penetration (if operating in sandy/volcanic environment)
- Possible stabilizer leaching (in harsh environments)

**Substrate:**
- γ′ coarsening: 400 nm → 1–3 μm after 30,000h at 980°C → creep resistance reduction
- Creep elongation accumulates
- Interdiffusion: elements from bond coat diffuse in (Al, Cr) → can modify γ′ volume fraction near surface → "secondary reaction zone" (SRZ) — potential issue with Pt-Al bond coats on modern alloys

### Secondary Reaction Zone (SRZ)

A subtle but important degradation mechanism unique to SX blades with Pt-Al coatings:

Interdiffusion of Al from bond coat into the alloy at high temperature (> 1,050°C) → local Al enrichment near the blade surface → precipitation of needle-like TCP phases (P-phase, σ-phase) in a thin zone below the coating.

This SRZ:
- Can extend 30–100 μm into the alloy
- Reduces local creep strength
- Can nucleate cracks
- Requires careful blade design to keep SRZ within tolerances

Prevention: minimize bond coat Al activity; Re addition to alloy retards SRZ growth; careful temperature management.

---

## 8. System Life Prediction

**TBC life prediction models:**

The primary failure mode (TGO-driven spallation) is modeled as:

```
N_spall = f(TGO thickness, TGO growth stress, TBC/TGO fracture toughness, ΔT per cycle)

Simplified: N_spall ∝ K_Ic² / (σ_TGO² × a_TGO)

where K_Ic = fracture toughness of TBC/TGO interface
      σ_TGO = stress in TGO (function of TGO thickness and CTE mismatch)
      a_TGO = TGO thickness (crack size)
```

As TGO grows: a_TGO increases → K approaches K_Ic → final fracture

**Probabilistic life assessment:**

Because blade-to-blade variability in:
- Bond coat surface roughness
- TGO morphology
- Film hole location relative to hot spots

Individual blade TBC life has a distribution (Weibull statistics are used). Engine certification is based on the Pth percentile (P = 0.1%–1%) surviving to the design life — not the average.

---

## 9. Next-Generation Coating Concepts

### Multilayer TBC

Instead of single YSZ layer, stack:
- **Bottom (inner) layer**: dense, erosion-resistant YSZ (no gaps → good erosion resistance, good adhesion)
- **Top (outer) layer**: porous YSZ or novel low-k ceramic (low conductivity, CMAS resistance)

The dense inner layer protects TGO from CMAS infiltration; the porous outer layer provides most insulation.

### "Strain-Tolerant" Dense Coatings

Vertically cracked (VC) APS TBC: deliberately crack the coating during deposition → vertical crack pattern that mimics EB-PVD compliance at lower cost. Currently used on some industrial GT vane applications.

### CMAS-Resistant Coatings

Gadolinium zirconate (Gd₂Zr₂O₇) outer layer: GZO reacts with CMAS melt → forms an apatite (Ca₂Gd₈(SiO₄)₆O₂) crystalline layer that blocks further CMAS penetration. Under active development for Middle East and volcanic environments.

### Self-Healing TBC

Concept: embed microcapsules of "healing agents" (oxide precursors) in TBC. If a crack forms, capsules rupture → healing agent fills crack → re-oxidizes → seals the crack. Still at research stage.

---

## 10. Complete Process Sequence — Blade to Engine

From raw CMSX-4 casting to flight-ready coated blade:

```
1. AS-CAST BLADE (from DS/SX casting)
   ↓
2. SHELL KNOCKOUT + CORE REMOVAL
   (remove ceramic mold + dissolve ceramic core)
   ↓
3. HOT ISOSTATIC PRESSING (HIP)
   1170°C, 175 MPa, 4h → close porosity
   ↓
4. SOLUTION HEAT TREATMENT
   1280–1335°C, 2–6h, vacuum → homogenize chemistry, dissolve eutectic γ′
   ↓
5. AGING HEAT TREATMENTS
   1140°C/6h + 870°C/20h → develop fine cuboidal γ′
   ↓
6. INSPECTION (BASELINE)
   Fluorescent penetrant, X-ray, Laue diffraction (orientation check)
   ↓
7. PLATINUM ELECTROPLATING (for Pt-Al bond coat)
   5–10 μm Pt in sulfuric acid plating bath
   ↓
8. ALUMINIZING (pack cementation or CVD)
   850–1,100°C; AlCl₃ atmosphere; 2–8h → form β-NiAl or (Ni,Pt)Al
   ↓  [OR: LPPS/HVOF MCrAlY SPRAY instead of steps 7–8]
9. PRE-OXIDATION
   1,000–1,050°C, 1–5h, air → establish α-Al₂O₃ TGO layer
   ↓
10. EB-PVD TBC DEPOSITION
    950°C substrate, 100–200 μm YSZ, 30–60 min per blade
    ↓
11. LASER DRILLING OF FILM COOLING HOLES (if not pre-drilled)
    300–600 holes, 0.3–0.5 mm dia., through TBC + metal wall
    ↓
12. FINAL INSPECTION
    Visual, dimensional, X-ray (hole breakthrough), thermal imaging
    ↓
13. ENGINE ASSEMBLY AND FLIGHT SERVICE
```

**Total blade processing time (from casting to engine-ready):** 4–8 weeks.

---

## Summary

| Component | Material | Function | Life Limiter |
|-----------|----------|----------|-------------|
| YSZ TBC (EB-PVD) | 7 wt% Y₂O₃-ZrO₂ | Thermal insulation; 100–150°C temp. reduction | TGO growth → spallation; sintering; CMAS |
| TGO | α-Al₂O₃ | Thermodynamic barrier; chemical stability | Thickness → stress → delamination at ~8 μm |
| Bond coat | MCrAlY or Pt-NiAl | Oxidation protection; Al reservoir; TBC adhesion | Al depletion; rumpling; SRZ |
| Superalloy | CMSX-4 | Structural strength | Creep; LCF; oxidation through failed coating |

**System interaction:** All four layers must maintain compatibility throughout service. The coating system fails when the TGO grows thick enough to accumulate sufficient strain energy to drive delamination — this is set by the TGO parabolic growth rate, the CTE mismatch, and the number of thermal cycles.

**Key innovations enabling 30,000h life:** (1) Y in bond coat → slower TGO growth; (2) EB-PVD columnar TBC → cycle-tolerant; (3) Pt modification → scale adhesion; (4) Pre-oxidation to establish α-Al₂O₃ before TBC.

---

## Exercises

1. A Pt-aluminide bond coat has 5 μm Pt (density 21.4 g/cm³) on a blade surface area of 100 cm². (a) What mass of Pt is on one blade? (b) At $32,000/kg for Pt, what is the Pt cost per blade? (c) For an engine with 80 HPT blades, what is the total Pt investment per engine?

2. TBC system thermal resistance: A blade has MCrAlY bond coat (k=15 W/mK, t=0.12 mm), TGO (k=30 W/mK, t=5 μm), and YSZ TBC (k=2.1 W/mK, t=0.15 mm). The gas-side h=3,000 W/m²K (T_gas = 1,580°C), internal h=2,800 W/m²K (T_coolant=670°C), metal k=11 W/mK, t=2.2 mm. Set up the full thermal resistance network and calculate: (a) heat flux through the blade, (b) temperature at each interface. Which layer provides most thermal resistance?

3. Pre-oxidation converts θ-Al₂O₃ to α-Al₂O₃ with 7% volume contraction. If the TGO is 0.5 μm thick at the time of conversion, what in-plane strain results from this contraction? Could this crack the TGO/TBC interface? (K_Ic of YSZ/TGO interface ≈ 0.5 MPa√m; assume crack size = TGO thickness; E_TGO = 400 GPa; plane-stress: σ_c = K_Ic/√(πa))

4. Secondary reaction zone (SRZ): Al diffuses from Pt-aluminide bond coat (10 at% Al) into CMSX-4 (nominally 13 at% Al). The SRZ grows by diffusion: d_SRZ = C × √(D_Al × t), where D_Al in Ni ≈ 10⁻¹⁶ m²/s at 1,100°C. After 10,000 hours: (a) What is d_SRZ? (b) If blade wall is 2.5 mm and SRZ reduces local strength by 30%, what fraction of the wall's load-carrying capability is lost? (c) Is this within acceptable limits (typically < 100 μm SRZ)?

5. CMAS infiltration: A volcanic eruption introduces CMAS particles (melting point 1,240°C) into the engine. TBC surface temperature is 1,120°C normally, but brief excursions to 1,280°C occur during high-power climb. (a) Can CMAS melt during normal operation? During excursions? (b) EB-PVD TBC column gaps are ~0.1 μm wide and 150 μm deep. CMAS viscosity at 1,280°C ≈ 10 Pa·s. Using the Washburn equation for capillary penetration depth d = √(γ r t / 2η) where γ = surface tension ~0.3 N/m, r = gap radius ~0.05 μm, estimate penetration depth after 1 hour of excursion. (c) What design change would reduce CMAS penetration risk?
