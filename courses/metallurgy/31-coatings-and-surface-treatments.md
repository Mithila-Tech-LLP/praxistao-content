# Chapter 31: Coatings and Surface Treatments — Protecting the Surface

> **"Most materials fail from the outside in. The surface bears the brunt of corrosion, wear, fatigue crack initiation, and oxidation. Protecting or modifying the surface can extend component life by 10× at a fraction of the cost of changing the bulk material. A $10 zinc coating on a steel beam might prevent a $10,000 replacement for 50 years. A $500 thermal barrier coating on a turbine blade enables temperatures that the underlying metal couldn't survive for 5 minutes."**

---

## Table of Contents

1. [Why Coat? The Surface-Dominated Failure Philosophy](#1-why-coat-the-surface-dominated-failure-philosophy)
2. [Electroplating](#2-electroplating)
3. [Hot-Dip Galvanizing](#3-hot-dip-galvanizing)
4. [Anodizing](#4-anodizing)
5. [Conversion Coatings (Phosphate, Chromate)](#5-conversion-coatings-phosphate-chromate)
6. [Physical Vapor Deposition (PVD)](#6-physical-vapor-deposition-pvd)
7. [Chemical Vapor Deposition (CVD)](#7-chemical-vapor-deposition-cvd)
8. [Thermal Spray Coatings (HVOF, Plasma)](#8-thermal-spray-coatings-hvof-plasma)
9. [Pack Cementation (Aluminide Coatings)](#9-pack-cementation-aluminide-coatings)
10. [Organic Coatings — Paints and Polymers](#10-organic-coatings--paints-and-polymers)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Coat? The Surface-Dominated Failure Philosophy

**Mismatch between bulk needs and surface needs:**
- Bulk material: must carry loads → needs strength, toughness
- Surface: must resist environment → needs oxidation resistance, corrosion immunity, wear resistance, or fatigue strength

Often, no single material satisfies both. The solution: use a strong, tough bulk with a protective surface.

**Coating selection criteria:**
| Factor | Requirement |
|--------|-------------|
| Adhesion | Coating must bond to substrate |
| CTE compatibility | Mismatch → spallation under thermal cycling |
| Chemical stability | Must not react destructively with substrate |
| Service temperature | Coating must be stable at operating T |
| Deposition damage | Process must not degrade substrate (avoid grit blast on SX blades!) |

---

## 2. Electroplating

**Principle:** Metal ions in solution deposit on a cathode (the part being plated) when current is applied:
```
Plating bath: M^n+ + ne⁻ → M (deposits on cathode = workpiece)
Anode: M → M^n+ + ne⁻ (dissolves to replenish bath, or inert with periodic top-up)
```

**Common electroplated coatings:**

| Coating | Thickness | Purpose | Application |
|---------|-----------|---------|-------------|
| Chrome (Cr) | 10–250 μm | Hard, wear-resistant, corrosion | Hydraulic rods, engine cylinders |
| Nickel | 5–50 μm | Corrosion, decorative undercoat | Electronics, plumbing fixtures |
| Zinc | 5–25 μm | Sacrificial corrosion protection | Steel hardware, bolts |
| Cadmium | 5–25 μm | Corrosion on steel (aerospace) | Aircraft fasteners (being phased out — Cd toxic) |
| Gold | 0.1–5 μm | Contact resistance, biocompat. | Electronics, medical implants |
| Copper | 10–100 μm | Under-coat, electrical | Printed circuit boards |
| Platinum | 5–10 μm | Catalytic activity, on SX blades | Platinum-aluminide bond coat precursor |

**Electroless nickel:** Chemical reduction (no external current) → more uniform deposit on complex geometries:
Ni²⁺ + hypophosphite → Ni-P alloy coating. The P content controls hardness (>10%P → 1,000 HV after heat treatment).

**Environmental concern:** Hard chrome uses hexavalent Cr (Cr⁶⁺) — carcinogenic. Being replaced by:
- Trivalent chrome plating
- HVOF-sprayed WC-Co (harder and less environmentally problematic)
- Hard anodizing (for aluminum)

---

## 3. Hot-Dip Galvanizing

**The most widely used zinc coating for structural steel.** Steel is dipped in molten zinc at ~450°C:

```
Process:
  Steel (clean surface) → flux bath → molten Zn (450°C) → withdraw → cool → 
  FeZn alloy layers + pure Zn outer layer
  
Coating structure from substrate outward:
  Fe substrate
  │
  Γ layer (Fe₃Zn₁₀) — thin, hard
  δ layer (FeZn₇) — hard, brittle 
  ζ layer (FeZn₁₃) — harder
  η layer (pure Zn) — soft, ductile, provides sacrificial protection
```

**Coating thickness:** 45–100 μm typical; thicker for aggressive environments.

**Life expectancy in various environments:**
```
Rural atmosphere:   60–100 years
Industrial (moderate):  30–50 years  
Seawater:            10–20 years
```

**Zinc sacrificial protection:** Even where coating is scratched, zinc corrodes preferentially (Zn is anodic to Fe). This is unlike paint where any scratch exposes unprotected steel.

**Continuous galvanizing (galvannealed):** Steel strip through zinc bath → short anneal → Fe diffuses into Zn → Fe-Zn alloy coating (GA) → better paint adhesion. Used for auto body panels.

---

## 4. Anodizing

**Electrochemical oxidation of aluminum** to build a thicker, controlled Al₂O₃ layer (5–30 μm vs. 2–5 nm natural):

```
Anodizing process:
- Electrolyte: H₂SO₄ (most common), chromic acid, oxalic acid
- Al part is the ANODE: Al → Al³⁺ + 3e⁻
- O²⁻ from water reacts: 2Al³⁺ + 3O²⁻ → Al₂O₃
- Result: porous anodic oxide layer grows into and out of surface
```

**Anodic oxide properties:**
- Thickness: 5–30 μm (versus natural oxide ~2-5 nm)
- Harder: 250–400 HV (much harder than Al substrate at 80–150 HV)
- Porous (except for barrier anodizing): pores can be sealed (boiling water, Ni acetate) or dyed

**Types of anodizing:**
| Type | Electrolyte | Thickness | Hardness | Purpose |
|------|------------|-----------|---------|---------|
| Type I (chromic acid) | CrO₃ | 2–8 μm | Moderate | Aerospace (fatigue-sensitive, thin = less K_t) |
| Type II (sulfuric) | H₂SO₄ | 5–25 μm | Good | General corrosion, dyeing |
| Type III (hard) | Cold H₂SO₄ | 25–75 μm | 300–500 HV | Wear resistance, hydraulics |

**Sealing:** After anodizing, hydrate the oxide: Al₂O₃ + H₂O → Al₂O₃·H₂O (boehmite) → closes pores → maximum corrosion resistance.

---

## 5. Conversion Coatings (Phosphate, Chromate)

**Conversion coatings** convert the metal surface itself into a non-metallic compound through chemical reaction:

### Phosphate Coating (Parkerizing)

Steel (or Al) + dilute phosphoric acid + Zn/Fe/Mn salts:
- Fe surface reacts → zinc or iron phosphate crystals form
- Coating: 5–30 μm, porous, gray
- Excellent paint adhesion base → combined with paint = very effective corrosion protection
- Mild corrosion resistance alone → must be painted
- Also provides dimensional stability (auto body prep)

**Lubricating phosphate:** Porous phosphate + oil soak → running-in lubricant for gears, fasteners → reduces galling during assembly.

### Chromate Conversion Coating (Alodine/Iridite on Al)

Al + chromic acid → thin chromate film (1–3 μm):
- Moderate corrosion resistance (500+ hours salt spray)
- Excellent paint adhesion primer
- Used in aerospace on Al substructure (every 7075, 2024 part)

**Hexavalent chromate (Cr⁶⁺):** Traditional; excellent self-healing properties (Cr⁶⁺ can migrate to damaged areas). RESTRICTION: EU RoHS/REACH severely limits Cr⁶⁺ → industry switching to:
- Trivalent chromate (Cr³⁺): less effective self-healing
- Chrome-free: TCP (trivalent chromium process), Ti/Zr-based conversions

---

## 6. Physical Vapor Deposition (PVD)

**Thin, hard, wear-resistant coatings** deposited by evaporation or sputtering in vacuum.

**Process:**
```
Target material (Ti, Cr, Al, TiAl etc.) in vacuum chamber
   ↓ (evaporation by heat, e-beam, or sputtering by Ar⁺ ions)
Vapor phase atoms + reactive gas (N₂ for nitrides, C₂H₂ for carbides)
   ↓ (deposit on substrate)
Thin film on part
```

**Common PVD coatings:**
| Coating | Hardness (HV) | Max T (°C) | Application |
|---------|-------------|-----------|-------------|
| TiN (gold color) | 2,300 | 600 | Drill bits, tap-and-die, molds |
| TiAlN (dark) | 3,000 | 800 | High-speed cutting (dry machining) |
| CrN | 1,800 | 700 | Punches, forming dies, plastic molds |
| DLC (diamond-like) | 3,000–8,000 | 350 | Engine parts, medical, optics |
| AlCrN | 3,200 | 900 | High-temperature cutting |
| TiCN | 3,000 | 400 | High-hardness cutting |

**Thickness:** Typically 1–10 μm (much thinner than electroplate or thermal spray)
**Adhesion:** Excellent for machined substrate (compressive stress in coating)
**Substrate temperature:** 200–600°C during deposition → limits substrate material choice

---

## 7. Chemical Vapor Deposition (CVD)

**Thicker, very hard coatings** deposited by gas-phase chemical reaction at high temperature:

```
Typical CVD for TiC:
  TiCl₄(g) + CH₄(g) → TiC(s) + 4HCl(g)   at 950–1,050°C
  
Deposited on: cemented carbide cutting tools (WC-Co substrate)
Temperature: 950–1,100°C (much higher than PVD → not suitable for heat-treated steels)
Thickness: 5–20 μm
```

**Common CVD coatings on cutting tools:**
- TiC (base layer, high hardness, wear-resistant)
- Al₂O₃ (middle layer, chemical stability, heat-barrier)
- TiN (top layer, low friction, gold color indicator of wear)

**Multilayer CVD:** TiN/TiCN/Al₂O₃/TiN → optimized combination (modern carbide inserts = 3–7 layers)

**CVD in pack cementation (turbine blade aluminizing):** Covered separately in Ch 46.

**PECVD (Plasma-Enhanced CVD):** Lower T (100–500°C) → allows coating of polymers, heat-sensitive substrates → microelectronics (SiO₂, Si₃N₄ dielectrics).

---

## 8. Thermal Spray Coatings (HVOF, Plasma)

**Thick functional coatings deposited by melting powder and projecting it onto the surface.** Covered in detail in Ch 46 for turbine blade bond coats. Overview here:

**HVOF (High Velocity Oxy-Fuel):**
- Fuel + O₂ combustion → very high gas velocity (>600 m/s) + moderate temperature
- Dense, low-porosity coating (< 1%), high bond strength
- Best for: WC-Co (hard chrome replacement), MCrAlY bond coats

**Air Plasma Spray (APS):**
- Plasma torch melts powder → particles hit surface and splat
- Lamellar, more porous than HVOF
- Best for: TBC top coat (porosity beneficial for thermal insulation), abradables

**LPPS (Low Pressure Plasma Spray):**
- APS in controlled atmosphere → oxidation-free → best for MCrAlY bond coats (no oxide inclusions)

**Cold spray:**
- Kinetic energy (no melting) → particles deform plastically → bond
- Very dense, minimal oxidation, low heat input
- Best for: corrosion-sensitive areas, dimensional restoration, Cu/Ti coatings

---

## 9. Pack Cementation (Aluminide Coatings)

This is the dominant coating process for turbine blade protection — covered in detail in Ch 46. Summary:

**Process:** Blade packed in powder (Al source + activator halide + inert Al₂O₃):
- At 850–1,100°C, AlCl vapors transport Al to blade surface
- Al diffuses into Ni → NiAl (β) intermetallic forms (β-NiAl aluminide)
- Surface Al concentrations: 30–40 at% Al → forms Al₂O₃ on oxidation

**Results:**
- 25–75 μm aluminide coating on Ni superalloy
- Excellent oxidation protection (Al₂O₃ scale forms during service)
- Used on all turbine blades since 1960s

---

## 10. Organic Coatings — Paints and Polymers

The most widely applied corrosion protection system by volume.

**Paint system for steel structure:**
```
Substrate (cleaned steel)
│
Surface preparation (blast to Sa 2.5 or Sa 3)
│
Primer (zinc-rich epoxy: sacrificial protection + adhesion) ← 50–75 μm
│
Intermediate coat (epoxy: barrier) ← 75–125 μm
│
Topcoat (polyurethane: UV resistance, appearance) ← 50–75 μm
Total: 175–275 μm
```

**Surface preparation is critical:**
- ISO 8501-1: cleanliness grades (Sa = abrasive blasted, ST = manually cleaned)
- Sa 2.5 (near-white blasted) → minimum for offshore coating
- Profile: Rz = 40–100 μm (anchor profile for paint adhesion)

**Coating types:**
| Type | Mechanism | Temperature (°C) | Application |
|------|-----------|-----------------|-------------|
| Zinc-rich primer | Sacrificial | < 80 | Structural steel |
| Epoxy | Barrier | < 120 | Industrial equipment |
| Polyurethane | Barrier + UV | < 100 | Topcoats, aircraft |
| Alkyd | Barrier | < 80 | General purpose |
| Zinc silicate | Sacrificial inorganic | < 400 | Hot surfaces |
| Ceramic/silicone | High T barrier | < 700 | Exhaust, furnaces |

---

## Summary

| Method | Thickness | Hardness | Temperature | Best For |
|--------|-----------|---------|-------------|---------|
| Electroplating | 5–250 μm | Moderate–high | < 200°C dep. T | Corrosion, wear, electrical |
| Hot-dip galvanize | 45–100 μm | Low (Zn) | — | Sacrificial, structural steel |
| Anodizing | 5–75 μm | 250–500 HV | — | Al corrosion + wear |
| PVD (TiN, TiAlN) | 1–10 μm | 2,000–3,200 HV | Up to 900°C (coating) | Cutting tool wear |
| CVD (TiC, Al₂O₃) | 5–20 μm | 3,000 HV | > 900°C dep. T | Carbide inserts |
| HVOF (WC-Co) | 100–500 μm | 1,400 HV | — | Hard chrome replacement |
| Pack cementation | 25–75 μm | Moderate | 1,100°C service | Turbine blade oxidation |
| Organic paint | 100–300 μm | Low | < 120°C | Structural steel corrosion |

---

## Exercises

1. A hydraulic piston rod (steel 1045) is currently protected by hard chromium electroplate (250 HV, 100 μm). Due to RoHS regulations, you must find a replacement. Compare: (a) HVOF WC-17Co (1,250 HV, 200 μm), (b) HVOF WC-10Co-4Cr (1,300 HV, 200 μm), (c) PVD CrN (1,800 HV, 5 μm). For each: comment on hardness, corrosion resistance in hydraulic fluid, fatigue impact (coating introduces residual stress), and processing cost. Which would you recommend?

2. An aluminum structural component (2024-T3) is to be anodized and primed for aerospace service. Choose between Type I (chromic, 5 μm) and Type III (hard, 25 μm): (a) Which provides better wear resistance? (b) For a part with many small holes (< 0.5mm diameter), which process is preferred (Type I has better throwing power)? (c) Type III anodizing reduces fatigue limit by ~30% due to tensile residual stresses in the coating. If the part design allows σ_a = 120 MPa, does it still exceed 70% of the endurance limit for uncoated 2024-T3 (endurance ≈ 138 MPa)? (d) What additional surface treatment would you apply after anodizing?

3. Calculate the mass of zinc consumed per year on a galvanized bridge pier (10 m long, 0.5 m diameter cylinder) in a marine atmosphere with corrosion rate of 12 μm/year for zinc (density ρ_Zn = 7.14 g/cm³). (a) Calculate surface area. (b) Calculate zinc mass loss per year. (c) If initial coating thickness is 80 μm, when will the zinc coating be consumed? (d) What happens to the steel when zinc is consumed — will it corrode at the same rate as bare steel?

4. A jet engine HPT disk (Inconel 718) needs a hard wear-resistant coating on the fir-tree root where the blade and disk contact. Service temperature is 650°C. Compare: (a) TiAlN PVD (max service T = 800°C, HV = 3,000), (b) CrN PVD (max service T = 700°C, HV = 1,800), (c) CVD Al₂O₃ (max service T = 1,100°C, HV = 2,200, deposition T = 1,000°C). Which process cannot be used for IN718 (σ_y = 1,034 MPa at RT, drops to 700 MPa at 650°C, aging precipitation begins dissolving above 700°C during CVD)? Of the remaining options, which gives better wear resistance?

5. A 30-year-old offshore oil platform requires recoating. The existing coating has failed in 40% of the area — blistering and rust under the film. (a) Explain the mechanism: why does rust form under intact-looking paint? (hint: osmotic blistering, cathodic disbondment). (b) The platform also has ICCP (impressed current cathodic protection). If the coating fails in area A but ICCP is still running: what is the effect on: (i) the bare steel area A (hint: cathodically protected at -0.85V), (ii) the coated area adjacent to A (hint: cathodic disbondment from excess current). (c) Design the maintenance plan: surface preparation, coating system, and ICCP current adjustment.
