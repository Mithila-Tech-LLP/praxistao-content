# Chapter 22: Steels — Iron's Greatest Achievement

> **"Steel built the modern world. It is the skeleton of every skyscraper, the rail that carries every train, the casing of every bearing, the armor of every submarine, the spring in every pen. It is estimated that over 1.8 billion tonnes of steel are produced every year — more than the combined production of all other metals. Yet steel is not one material. It is a family of thousands of alloys, united by iron and carbon, differentiated by everything else."**

---

## Table of Contents

1. [What Is Steel?](#1-what-is-steel)
2. [AISI/SAE Steel Classification System](#2-aisisae-steel-classification-system)
3. [Plain Carbon Steels](#3-plain-carbon-steels)
4. [Alloy Steels](#4-alloy-steels)
5. [Stainless Steels — Passivation and Corrosion Resistance](#5-stainless-steels--passivation-and-corrosion-resistance)
6. [Tool Steels](#6-tool-steels)
7. [HSLA (High-Strength Low-Alloy) Steels](#7-hsla-high-strength-low-alloy-steels)
8. [Maraging Steels](#8-maraging-steels)
9. [Steel Processing — From Ingot to Product](#9-steel-processing--from-ingot-to-product)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is Steel?

Steel is an iron-carbon alloy containing < 2.14% carbon (above this is cast iron, Ch 23). In practice, most steels contain 0.05–1.5% C, with the bulk of production in the 0.1–0.5% range.

**Why iron-carbon is special:**
- Iron's FCC → BCC transformation allows dramatic property changes via heat treatment
- Carbon dramatically strengthens iron (solid solution + precipitation as carbides)
- Wide composition range → tunable properties from soft (0.05%C) to extremely hard (1.2%C)
- Iron is abundant (4th most abundant element in Earth's crust) and cheap to produce

---

## 2. AISI/SAE Steel Classification System

The four-digit AISI (American Iron and Steel Institute) / SAE (Society of Automotive Engineers) system:

```
XXXX
│││└─ Last two digits = carbon content × 100 (e.g., 40 = 0.40% C)
││└── Minor designation (sometimes used for subgroups)
│└─── Major designation within series
└──── Series (principal alloying element)
```

| Series | Principal Alloying Element |
|--------|---------------------------|
| 1xxx | Carbon steels |
| 10xx | Plain carbon (no major alloys) |
| 11xx | Resulfurized (free machining, S added) |
| 12xx | Resulfurized + rephosphorized |
| 13xx | Manganese (1.75% Mn) |
| 2xxx | Nickel steels |
| 3xxx | Nickel-chromium steels |
| 4xxx | Molybdenum steels (4Mo, Cr-Mo, Ni-Mo) |
| 41xx | Chromium-molybdenum (Cr-Mo) |
| 43xx | Nickel-chromium-molybdenum (Ni-Cr-Mo) → 4340 |
| 5xxx | Chromium steels |
| 52100 | High-carbon chromium bearing steel |
| 6xxx | Chromium-vanadium steels |
| 8xxx | Nickel-chromium-molybdenum |
| 9xxx | Silicon-manganese |

**Example: 4340 steel**
- 43xx series = Ni-Cr-Mo alloy
- Last two digits = 40 → 0.40% C
- Full composition: 0.38–0.43%C, 0.60–0.80%Mn, 1.65–2.00%Ni, 0.70–0.90%Cr, 0.20–0.30%Mo
- A classic structural alloy steel for landing gear, shafts, gears

---

## 3. Plain Carbon Steels

Plain carbon steels have only C (and small amounts of Mn, Si for deoxidation) as alloying elements:

### Low-Carbon Steel (< 0.25% C)
- **Properties:** Soft, ductile, easily formed and welded
- **Microstructure:** Ferrite + small amount of pearlite
- **Examples:** AISI 1010 (auto body panels), 1020 (structural shapes, rivets)
- **Use:** Sheet metal, structural sections, rods, fasteners
- **Limitations:** Low strength (σ_y = 180–280 MPa); not heat-treatable for high strength

### Medium-Carbon Steel (0.25–0.60% C)
- **Properties:** Higher strength, reduced ductility; can be heat-treated
- **Microstructure:** Normalized: mixed ferrite+pearlite; Q&T: tempered martensite
- **Examples:** 1040 (gears, shafts), 1045 (axles), 1060 (springs)
- **Use:** Machine parts, crankshafts, axles, springs
- **Heat treatment:** Q&T to σ_y = 600–1,000 MPa

### High-Carbon Steel (0.60–1.25% C)
- **Properties:** High strength and hardness after heat treatment; lower ductility/toughness
- **Microstructure:** After Q&T: tempered martensite (hard); after spheroidize: spheroidized carbide in ferrite (soft, machinable)
- **Examples:** 1080 (rails, wire rope), 1095 (springs, knives), 52100 (bearings)
- **Use:** Cutting tools, bearings, springs, high-hardness wear parts

---

## 4. Alloy Steels

Adding alloying elements to steel improves:
- **Hardenability** (Mn, Cr, Mo, Ni shift TTT to longer times)
- **High-temperature strength** (Mo, V form stable carbides)
- **Toughness** (Ni)
- **Corrosion resistance** (Cr, Ni)
- **Wear resistance** (Cr, V, W carbides)

**Key alloy steels:**

| Steel | Composition | Key Properties | Applications |
|-------|-------------|---------------|--------------|
| 4140 | Cr-Mo (0.40%C) | Hardenability, toughness | Shafts, gears, bolts |
| 4340 | Ni-Cr-Mo (0.40%C) | Deep hardenability, high K_Ic | Landing gear, crankshafts |
| 8620 | Ni-Cr-Mo (0.20%C) | Carburizable, tough core | Case-hardened gears |
| 300M | Modified 4340 + Si, V | K_Ic = 120 MPa√m at 1,930 MPa UTS | Airframe (ultra-high strength) |
| AF1410 | Co-Ni + Cr, Mo, C | K_Ic > 154 MPa√m at 1,520 MPa | Navy submarine hulls |

---

## 5. Stainless Steels — Passivation and Corrosion Resistance

**Definition:** Steel containing at least 10.5% Cr (by weight). The Cr forms a thin, tenacious Cr₂O₃ passive film on the surface that prevents further oxidation and corrosion.

```
Stainless steel family tree:

Stainless steels
├── Austenitic (300 series): FCC, highest corrosion resistance
├── Ferritic (400 series): BCC, moderate corrosion, magnetic
├── Martensitic (400 series): BCC, heat-treatable, lower corrosion
├── Duplex: two-phase austenite + ferrite, excellent SCC resistance
└── Precipitation hardening (PH): high strength + corrosion resistance
```

**Austenitic stainless steels (304, 316, 310):**
- Cr 16–26%, Ni 8–22% → stabilizes austenite (FCC) at RT
- Not heat-treatable for hardness (no martensite)
- Strengthened by cold work only
- Non-magnetic (useful for MRI equipment)
- **304 (18-8):** 18%Cr-8%Ni; most widely used; food, chemical equipment
- **316 (18-12-2Mo):** +2%Mo; better pitting resistance; marine, pharmaceutical
- **310:** 25%Cr-20%Ni; high-temperature oxidation resistance to 1,050°C

**Sensitization problem:** Holding 304 at 450–850°C → Cr₂₃C₆ precipitates at grain boundaries → Cr depletion near boundaries → susceptible to intergranular corrosion.
- **Fix:** Use 304L (0.03% max C) or 316L — less C → less Cr₂₃C₆ → less sensitization
- Or use 321 (Ti-stabilized) or 347 (Nb-stabilized) — Ti/Nb scavenge C preferentially

**Ferritic stainless steels (430, 444):**
- 16–30% Cr, low Ni, no austenite → cannot be hardened by Q&T
- Magnetic
- Lower cost than austenitic
- Used for: automotive trim, kitchen sinks, exhaust systems

**Martensitic stainless steels (410, 420, 440C):**
- 12–18% Cr, low Ni, 0.15–1.0% C → can be hardened by Q&T
- Magnetic; moderate corrosion resistance
- 440C: 17%Cr, 1.0%C → ball bearings, knives, high hardness

**Duplex stainless steels (2205, 2507):**
- ~22-25%Cr, 4-7%Ni, 3%Mo → 50% austenite + 50% ferrite
- Twice the yield strength of 304 + better SCC resistance in chlorides
- Used for: offshore platforms, desalination plants, chemical reactors

---

## 6. Tool Steels

Tool steels must maintain hardness at cutting temperature (up to 600°C for high-speed cutting) and resist wear:

| Group | Grade | Key Alloying | Properties | Use |
|-------|-------|-------------|-----------|-----|
| W (water quench) | W1 | High C only | Hard, brittle, cheap | Files, drills (light use) |
| O (oil quench) | O1 | Mn, Cr, W | Good wear | Precision tools |
| A (air hardening) | A2 | Cr, Mo, V | Low distortion | Punches, dies |
| D (high-C, high-Cr) | D2 | 12%Cr, 1.5%C | Exceptional wear | Blanking dies, gauges |
| H (hot work) | H13 | Cr, Mo, V | Hot hardness, thermal fatigue | Die casting dies, forge dies |
| T (tungsten HS) | T1 | W, Cr, V | Red hardness (600°C) | Lathe cutters |
| M (molybdenum HS) | M2 | Mo, W, Cr, V | Red hardness (cheaper) | Drills, end mills, hobs |

**High-speed steel (HSS, M2):** Secondary hardening (Ch 20) from Mo₂C, W₂C, V₄C₃ precipitation at 550°C → maintains 62+ HRC at red heat.

---

## 7. HSLA (High-Strength Low-Alloy) Steels

**Design philosophy:** Achieve higher strength than plain carbon steel using LESS carbon (better weldability) by adding small amounts of Nb, V, Ti:
- Grain refinement (TiN, NbC particles pin austenite grain boundaries during rolling)
- Precipitation hardening from NbC, V₄C₃ in ferrite
- No expensive Ni, Cr additions needed

```
HSLA vs. plain carbon steel:
              σ_y (MPa)  σ_UTS (MPa)  %C   Weldability
Plain C 1040:   380        590        0.40   Moderate (preheat)
HSLA A572 Gr50: 345        450        0.23   Excellent
HSLA A514:      700        760        0.18   Good
HSLA 100:      690        760        0.10   Very good
```

**Applications:** Ship hulls, bridges, pipelines, auto frames, wind turbine towers.

**TMCP (Thermomechanical Controlled Processing):**
Roll steel at low temperature to refine austenite grains → transform to very fine ferrite → σ_y = 500–700 MPa without heat treatment.

---

## 8. Maraging Steels

Covered in Ch 21 (precipitation hardening). Ultra-high-strength steels for aerospace:
- σ_y up to 2,070 MPa (grade 300)
- K_Ic = 80 MPa√m (higher than conventional ultra-high-strength steels)
- Applications: rocket motor cases, aircraft landing gear, tooling

---

## 9. Steel Processing — From Ingot to Product

```
Iron ore + coke + limestone
       ↓
  Blast furnace → Pig iron (4% C, Si, Mn, P, S impurities)
       ↓
  Basic Oxygen Furnace (BOF) or Electric Arc Furnace (EAF)
  → reduce C to < 2%, remove impurities
  → add alloying elements
       ↓
  Secondary metallurgy (ladle furnace)
  → precise composition adjustment
  → desulfurization, degassing (vacuum degassing for H₂)
       ↓
  Continuous casting → slabs, blooms, billets
       ↓
  Hot rolling → structural shapes, plates, sheet, bar
       ↓
  Cold rolling (sheet, strip) → precise dimensions + work hardened
       ↓
  Heat treatment as required (annealing, Q&T, normalizing)
       ↓
  Finished product
```

**Cleanliness:** Modern steels (especially bearing steel 52100, spring steel) are produced to extremely low inclusion levels using:
- Vacuum degassing (reduces H₂ → prevents hydrogen embrittlement)
- Ladle desulfurization (Ca injection → CaS inclusions, globular not stringer)
- Secondary refining → S < 0.002%, O < 15 ppm

---

## Summary

| Steel Family | C Range | Key Feature | Main Applications |
|-------------|---------|-------------|------------------|
| Low-C plain | 0.05–0.25% | Formable, weldable | Sheet, structural |
| Med-C plain | 0.25–0.60% | Q&T capable | Shafts, gears |
| High-C plain | 0.6–1.2% | High hardness | Bearings, springs |
| Alloy (4340) | 0.3–0.5% | Deep hardenability | Landing gear |
| Austenitic SS | Low C | Corrosion resistance | Food, chemical |
| Duplex SS | Low C | SCC resistance, high strength | Offshore |
| Tool steel (M2) | 0.8–1.2% | Red hardness | Cutting tools |
| HSLA | < 0.25% | Grain-refined, weldable, strong | Pipelines, ships |
| Maraging | Very low | Precipitation in martensite | Rocket cases |

---

## Exercises

1. A bridge designer needs to choose between: (a) 1040 plain carbon steel, (b) HSLA A572 Gr50, (c) 4340 alloy steel. Requirements: σ_y > 350 MPa, excellent field weldability, Charpy CVN > 50J at −20°C. Which steel meets all requirements? Calculate the weight savings per 1,000 kg of 1040 steel structure if switched to A572 Gr50 (same load → proportion to σ_y ratio).

2. Stainless steel 304 is welded into a chemical tank. After welding, the HAZ shows intergranular corrosion in an oxalic acid test. (a) Explain the sensitization mechanism including the Cr-depletion zone width (~5–10nm from boundary). (b) Three solutions: specify 304L, stabilize with Ti (321), or post-weld anneal at 1,050°C + quench — for each, explain the metallurgical basis. (c) The customer wants to switch to duplex 2205. What properties do they gain and lose compared to 304L?

3. High-speed steel M2 (0.85%C-6%W-5%Mo-4%Cr-2%V) is solution treated at 1,220°C, then triple tempered at 560°C for 1 hour each. (a) Why is the solution treatment temperature so high compared to plain carbon steel? (b) Sketch the hardness vs. tempering temperature curve for M2, marking the secondary hardening peak. (c) What alloy carbide precipitates are responsible for secondary hardening? (d) Why is triple tempering better than single 3-hour temper?

4. AISI 4340 steel (0.40%C, 1.8%Ni, 0.80%Cr, 0.25%Mo) is being considered for a 75mm diameter landing gear pin requiring through-hardening to achieve 38 HRC at center. The Jominy D_I (ideal quench diameter for 50% martensite) for 4340 is approximately 100mm. (a) Is 75mm within the D_I for 4340? Will oil quenching achieve martensite through the section? (b) For comparison, 1040 plain carbon steel has D_I ≈ 25mm. Would 1040 work for this application? (c) What tempering temperature and time would achieve 38 HRC ± 2 HRC?

5. A metallurgist needs to select a steel for offshore oil platform risers (long tubes connecting subsea wellhead to surface). Requirements: σ_y > 550 MPa, SCC resistance in seawater at −5°C, excellent weldability, CVN > 70J at −40°C. Compare: (a) HSLA X70 pipeline steel, (b) duplex stainless 2205, (c) austenitic 316L stainless. For each: state if it meets requirements, primary failure mode risk, and estimated cost per tonne. Which would you recommend and why?
