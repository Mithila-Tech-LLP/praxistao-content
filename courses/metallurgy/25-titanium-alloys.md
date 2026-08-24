# Chapter 25: Titanium Alloys — Strength, Light Weight, and Heat

> **"Titanium has an astonishing combination: a density of 4.51 g/cm³ (60% of steel), a corrosion resistance that rivals platinum, and the ability to retain useful strength up to 600°C. It is the preferred metal when weight matters and steel is too heavy, when corrosion matters and aluminum is too reactive, or when temperature matters and aluminum is too soft. The one thing that limits titanium is cost — roughly 10× aluminum for the same weight."**

---

## Table of Contents

1. [Titanium Properties and Why They Matter](#1-titanium-properties-and-why-they-matter)
2. [Crystal Structures — α and β Titanium](#2-crystal-structures--α-and-β-titanium)
3. [Titanium Alloy Classification](#3-titanium-alloy-classification)
4. [Commercially Pure Titanium (α)](#4-commercially-pure-titanium-α)
5. [Ti-6Al-4V — The Workhorse Alloy](#5-ti-6al-4v--the-workhorse-alloy)
6. [Near-Alpha Alloys (High Temperature)](#6-near-alpha-alloys-high-temperature)
7. [Beta Alloys (High Strength)](#7-beta-alloys-high-strength)
8. [Titanium Processing — The Kroll Process and Beyond](#8-titanium-processing--the-kroll-process-and-beyond)
9. [Titanium in Aerospace and Medical Applications](#9-titanium-in-aerospace-and-medical-applications)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Titanium Properties and Why They Matter

**Titanium's key numbers:**

| Property | Titanium | Steel | Aluminum |
|---------|----------|-------|----------|
| Density (g/cm³) | 4.51 | 7.87 | 2.70 |
| E (GPa) | 116 | 207 | 69 |
| σ_y best alloy (MPa) | 1,200 | 2,000 | 600 |
| Specific strength (kN·m/kg) | 270 | 250 | 220 |
| Max temperature (°C) | 600 (600°C alloys) | 500+ | 180 |
| Corrosion in seawater | Excellent | Poor | Good |
| Cost ($/kg approx) | $25–100 | $1–3 | $3–8 |

**Exceptional corrosion resistance:** TiO₂ passive film is extremely stable — Ti is essentially immune to seawater, chlorides, and many acids. Used in:
- Chemical industry (chloride environments that destroy stainless)
- Offshore/marine applications
- Biomedical implants (bone-compatible, no rejection)

---

## 2. Crystal Structures — α and β Titanium

Pure titanium undergoes an allotropic transformation:
```
> 882°C: β-titanium (BCC) — high temperature phase
< 882°C: α-titanium (HCP) — low temperature phase
```

**Alloying elements shift the β-transus temperature:**
- **α-stabilizers** (raise transus): Al, O, N, C → stabilize HCP α
- **β-stabilizers** (lower transus): V, Mo, Fe, Cr, Mn, Nb → stabilize BCC β

```
Effect on phase diagram (schematic):

T
│  β         β      β
│         ─────   ─────
│  α+β         β+α
│         α+β     
│  α         α
└──────────────────── %β-stabilizer →

α-stabilized: two-phase region shifts right
β-stabilized: two-phase region widens, β persists to RT
```

---

## 3. Titanium Alloy Classification

```
Titanium alloys
├── Commercially pure (CP) titanium: very low alloy content, mostly α
├── Alpha (α) alloys: all HCP; good creep, weldable; no heat-treat hardening
├── Near-alpha: mostly α + small β (<10%); high-T aerospace
├── Alpha-beta (α+β): two-phase; heat-treatable; Ti-6Al-4V
└── Beta (β) alloys: mostly/all BCC; heat-treatable; highest strength
```

---

## 4. Commercially Pure Titanium (α)

**Grades 1–4:** Increasing O content (O is an interstitial strengthener):
```
Grade 1: 0.18% O max → σ_y = 170 MPa → maximum ductility, formability
Grade 2: 0.25% O max → σ_y = 275 MPa → standard CP Ti (most common)
Grade 3: 0.35% O max → σ_y = 380 MPa
Grade 4: 0.40% O max → σ_y = 480 MPa → highest-strength CP grade
```

**Key application — chemical industry:** Grade 2 is used for heat exchanger tubes, reaction vessels, pump parts in aggressive chemical environments (HNO₃, chloride solutions) where stainless steel would corrode.

**Medical: Grade 1-2 for implant hardware** where maximum ductility is needed (dental implant roots, surgical fasteners).

---

## 5. Ti-6Al-4V — The Workhorse Alloy

**The most important titanium alloy.** Approximately 50% of all titanium production. An α+β alloy:
- **6% Al:** α-stabilizer; solid solution hardening; raises α-β transus to ~1,000°C
- **4% V:** β-stabilizer; retains some β phase at RT for strength

**β-transus:** ~995°C (vs. 882°C for pure Ti)

**Microstructures (depends on processing temperature relative to β-transus):**

```
Processed BELOW β-transus (α+β field, e.g., 950°C):
→ Equiaxed α (60-80%) + intergranular β
→ Typical for forgings (mill-annealed)
→ Good ductility, toughness, fatigue

Processed ABOVE β-transus (β field, e.g., 1,050°C):
→ Widmanstätten (basketweave) α laths in β matrix
→ Better creep resistance, worse fatigue
→ Used for engine disks requiring creep performance
```

**Heat treatment (α+β alloys):**
- **Mill annealed (MA):** Common reference condition, equiaxed α+β
- **Solution Treat + Age (STA):** Solution treat just below β-transus → quench → age at 500–600°C → α plates form in β matrix → higher strength but lower K_Ic
- **STA maximum:** σ_y ≈ 1,100 MPa, K_Ic ≈ 55 MPa√m

**Ti-6Al-4V properties (typical):**
| Condition | σ_y (MPa) | σ_UTS (MPa) | K_Ic (MPa√m) | % El |
|-----------|-----------|-------------|-------------|------|
| Annealed | 828 | 900 | 75 | 14 |
| STA | 1,100 | 1,170 | 55 | 8 |
| High toughness | 760 | 860 | 100+ | 18 |

**Applications:**
- Airframe: wing spars, bulkheads, engine pylons (Boeing 787: ~15% Ti by weight)
- Jet engine: fan blades (large-chord for CFM56, GEnx), compressor blades
- Medical: hip implants, spinal fusion cages, dental implants
- Sports: bicycle frames, golf club heads

---

## 6. Near-Alpha Alloys (High Temperature)

Near-alpha alloys contain mostly α-stabilizers with small amounts of β-stabilizers. Key property: **creep resistance up to 600°C**.

**IMI 834 (Ti-5.8Al-4Sn-3.5Zr-0.7Nb-0.5Mo-0.35Si):**
- Maximum use temperature: 600°C
- Used for: high-pressure compressor disks in modern engines (Trent 700, GE90 early stages)
- Creep resistant due to HCP structure (fewer slip systems than FCC/BCC)

**Ti-1100 (Ti-6Al-2.75Sn-4Zr-0.4Mo-0.45Si):**
- Maximum use temperature: 600°C
- Used in: high-pressure compressor blades

**Why titanium is used instead of nickel for compressor stages:**
- Ti is 60% of Ni density → ~40% weight savings
- Up to stage 4-6 of compressor: temperature < 600°C → titanium works
- Stage 7+ and turbine: temperature > 600°C → must switch to nickel superalloys

---

## 7. Beta Alloys (High Strength)

Beta alloys are stabilized with enough β-stabilizer to retain BCC structure to room temperature. They can be heat-treated to very high strength.

**Ti-10V-2Fe-3Al (Ti-10-2-3):**
- Solution treat in β field → quench → age → fine α precipitates in β matrix
- STA: σ_y = 1,200 MPa, K_Ic = 45 MPa√m
- Excellent forging characteristics (β alloys are easier to forge than α+β)
- Used for: Boeing 777 main landing gear (replaced 300M steel, 15% lighter)

**Ti-15V-3Cr-3Al-3Sn:**
- Highly formable (β alloys are cold-rollable unlike α+β)
- Used for: springs, strip, thin sheet for aircraft hydraulic tubing

**Ti-3Al-8V-6Cr-4Mo-4Zr (Beta-C):**
- Most corrosion-resistant Ti alloy in chloride environments
- Hydrogen industry, oil/gas, chemical processing

**Trade-off of beta alloys:** Higher strength but lower K_Ic, lower creep resistance (BCC has more slip systems → easier dislocation motion at elevated T), and higher density (more β-stabilizers like V, Mo → heavier).

---

## 8. Titanium Processing — The Kroll Process and Beyond

**Primary production (Kroll process):**
```
TiO₂ + 2Cl₂ + 2C → TiCl₄ + 2CO (chlorination in fluidized bed)
TiCl₄ + 2Mg → Ti (sponge) + 2MgCl₂  (Mg reduction in sealed retort)
MgCl₂ → Mg + Cl₂  (electrolysis to recycle Mg and Cl₂)
Ti sponge → vacuum melt → ingot
```

**Why titanium is expensive:**
- Kroll process is batch, labor-intensive, energy-intensive
- Pure Mg atmosphere required (even ppm O₂ contamination → TiO₂ embrittlement)
- Multiple vacuum arc remelts (VAR) required for homogeneity
- Extremely reactive in liquid form → cannot use conventional refractories

**Processing challenges:**
- Oxygen/nitrogen contamination: TiO₂ and TiN form at surface during hot working → must remove by machining (α case removal)
- Beta transus must be known precisely for heat treatment
- Hydrogen pickup: Ti absorbs H₂ → must vacuum degas below 550°C

**Powder metallurgy (PM) titanium:**
- PREP (plasma rotating electrode process) or HDH (hydride-dehydride) powder
- HIP at 920°C, 100 MPa → near-net-shape compressor disks
- Direct metal laser sintering (DMLS/SLM) for complex AM parts (Ch 63)

---

## 9. Titanium in Aerospace and Medical Applications

**Jet Engine Fan and Compressor:**
```
Stage 1-2 compressor blades: Ti-6Al-4V (< 400°C)
Stage 3-6 compressor blades: Ti-6Al-4V or IMI 834 (< 600°C)
Stage 7+ compressor blades:  Nickel superalloys (> 600°C)
Fan blades: Large chord Ti-6Al-4V (GE90: 22 blades 1.1m long)
Compressor disks: Ti-6Al-4V or forged Ti-10-2-3
```

**Airframe structures:**
- Bulkheads: Ti-6Al-4V forgings (500 mm × 1,000 mm single forgings)
- Boeing 787: 15% Ti by weight → combines with 50% CFRP
- Ti fasteners everywhere that CFRP contacts aluminum (prevent galvanic corrosion)

**Medical implants:**
Ti-6Al-4V ELI (Extra Low Interstitials: O < 0.13%):
- Hip replacements: femoral stem (Ti-6Al-4V) + acetabular cup
- Spinal implants: cages, rods, screws
- Dental implants: screw + abutment (CP-Ti most biocompatible)
- Advantage: osseointegration (bone grows directly into Ti surface → no cement needed)

---

## Summary

| Alloy Type | Example | Max T | σ_y (MPa) | Key Property | Application |
|-----------|---------|-------|-----------|-------------|-------------|
| CP Ti (α) | Grade 2 | ~300°C | 275 | Corrosion resistance | Chemical plant |
| α+β | Ti-6Al-4V | ~350°C | 828–1,100 | Balanced properties | Aerospace, medical |
| Near-α | IMI 834 | 600°C | 900 | Creep resistance | HP compressor |
| β alloy | Ti-10-2-3 | ~300°C | 1,200 | High strength | Landing gear |

---

## Exercises

1. Ti-6Al-4V is used for a compressor blade (stage 3, maximum metal temperature 420°C). The blade must have σ_y > 700 MPa and K_Ic > 60 MPa√m. Can Ti-6Al-4V mill-annealed condition meet both requirements? If not, what processing change would help, and what property would be traded? Calculate the weight saving versus a 4340 steel blade (σ_y = 900 MPa, ρ = 7.87 g/cm³) that meets both requirements.

2. A titanium heat exchanger tube (Grade 2, 25mm OD, 2mm wall) will carry chlorinated cooling water at 90°C. (a) What corrosion mechanism will attack 316L stainless steel in this environment that Ti resists? (b) At 90°C, TiO₂ reform rapidly if scratched — estimate time to form a 2 nm passive film (use parabolic oxidation: x² = k×t, k ≈ 10⁻²⁰ m²/s at 90°C). (c) What is the weight per meter of Ti tube vs. 316L tube of same dimensions (ρ_Ti = 4.51, ρ_SS = 7.93 g/cm³)?

3. A landing gear strut for a large civil aircraft must withstand 500 kN tensile load with safety factor 1.5 at minimum weight. Compare: (a) 300M ultra-high-strength steel (σ_y = 1,650 MPa, ρ = 7.85 g/cm³), (b) Ti-10-2-3 STA (σ_y = 1,200 MPa, ρ = 4.65 g/cm³), (c) 7075-T6 Al (σ_y = 503 MPa, ρ = 2.80 g/cm³). Calculate minimum cross-section area and strut volume per meter for each. Which gives minimum mass for a 2-meter strut?

4. Titanium fatigue crack growth follows Paris Law: da/dN = C(ΔK)^m with C = 10⁻¹¹ m/cycle per (MPa√m)^m and m = 4 for Ti-6Al-4V. A compressor disk has an assumed initial crack a_i = 1mm (NDE limit). Operating stress amplitude = 400 MPa. (a) Calculate ΔK at a_i. (b) K_Ic = 65 MPa√m for this alloy: calculate a_c (critical crack size). (c) Integrate Paris Law to find cycles from a_i to a_c. (d) If each flight has 10 stress cycles for this disk, how many flights before inspection is mandatory (use N/3 safety factor)?

5. The Boeing 787 uses carbon fiber composite (CFRP) for the primary airframe and titanium fasteners throughout. (a) Why can't aluminum fasteners be used where CFRP contacts titanium? (b) Why can't steel fasteners be used where CFRP contacts steel? (hint: galvanic corrosion EMF) (c) What is the galvanic series position of Ti relative to CFRP (carbon fiber)? Is the combination safe? (d) Estimate the number of Ti fasteners on a 787 and their total mass contribution (rough estimate: 500,000 fasteners × 5g each).
