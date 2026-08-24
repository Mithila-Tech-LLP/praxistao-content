# Chapter 29: Corrosion Fundamentals — Why Metals Decay

> **"Every year, corrosion costs the United States approximately $270 billion — 3.1% of GDP. Globally, the cost approaches $2.5 trillion annually. The rusting of a pipe, the greenish patina on copper, the white powder on aluminum: these are metals returning to their natural, lower-energy state. The metallurgist's job is to delay this process long enough to justify the cost of making the metal in the first place."**

---

## Table of Contents

1. [What Is Corrosion?](#1-what-is-corrosion)
2. [Thermodynamics — Why Metals Want to Corrode](#2-thermodynamics--why-metals-want-to-corrode)
3. [Dry Corrosion — Gas-Metal Reactions](#3-dry-corrosion--gas-metal-reactions)
4. [The Pilling-Bedworth Ratio](#4-the-pilling-bedworth-ratio)
5. [Wet (Electrochemical) Corrosion — Overview](#5-wet-electrochemical-corrosion--overview)
6. [Passivation — The Metal's Defense](#6-passivation--the-metals-defense)
7. [Pourbaix Diagrams — Stability in Aqueous Environments](#7-pourbaix-diagrams--stability-in-aqueous-environments)
8. [Corrosion Rate Measurement](#8-corrosion-rate-measurement)
9. [Practical Corrosion in Engineering](#9-practical-corrosion-in-engineering)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is Corrosion?

Corrosion is the chemical or electrochemical degradation of a metal through reaction with its environment. It is, fundamentally, the reverse of metal extraction — metals want to return to their oxidized state (ores).

**Two main types:**

**Dry corrosion (high-temperature):** Metal reacts with hot gas (oxygen, sulfur, etc.) → metal oxide, sulfide, etc. forms on the surface. Relevant for jet engines, furnaces, power plants.

**Wet (aqueous) corrosion:** Metal reacts with water and dissolved species at lower temperatures. Relevant for infrastructure, marine, chemical plant, biological environments.

Both are oxidation reactions: **M → M^n+ + n e⁻** (metal loses electrons = is oxidized).

---

## 2. Thermodynamics — Why Metals Want to Corrode

**Gibbs free energy of oxide formation:**

Metals are thermodynamically unstable in oxygen (most metals have negative ΔG_f for oxide formation):
```
Fe + ½O₂ → FeO    ΔG_f = -251 kJ/mol  ← iron WANTS to form FeO
Al + ¾O₂ → ½Al₂O₃ ΔG_f = -792 kJ/mol  ← aluminum VERY strongly wants to oxidize
Cu + ½O₂ → CuO    ΔG_f = -128 kJ/mol  ← copper less eager, hence "noble-ish"
Au + ½O₂ → ½Au₂O  ΔG_f = +77 kJ/mol   ← POSITIVE → Au does NOT oxidize spontaneously!
```

**Ellingham Diagram:** Plots ΔG_f of metal oxides vs. temperature. Metals at the BOTTOM are the most "oxidizable"; metals at the TOP are the most "noble":

```
Ellingham Diagram (schematic):

ΔG_f     Noble metals (Au, Pt) — top
(kJ/mol  ─────────────────────── Cu, Ag
of O₂)   ─────────────────────── Fe, Ni, Co
         ─────────────────────── Cr
         ─────────────────────── Ti
         ─────────────────────── Si, Al
         ─────────────────────── Mg, Ca
Base     ─────────────────────── La, Y, Ce (bottom)
metals
         ────────────────────────────────── Temperature →
```

**Key insight from Ellingham diagram:** A metal at the bottom can reduce (steal oxygen from) any oxide above it. Aluminum can reduce Fe₂O₃ (thermite reaction). Cr can reduce Y₂O₃ is debatable — but Y can reduce Cr₂O₃, hence Y additions improve oxidation resistance by forming stable Y₂O₃ at the scale-metal interface.

---

## 3. Dry Corrosion — Gas-Metal Reactions

At high temperatures, metals react with oxygen (and H₂S, SO₂, Cl₂ in combustion environments) to form surface scale:

**General reaction:** M(s) + O₂(g) → MO₂(s) (or other stoichiometry)

**Scale growth kinetics:**

The scale growth rate depends on whether oxygen or metal ions can diffuse through the scale:

**Parabolic growth law (compact, adherent scale):**
```
x² = k_p × t   (or x = √(k_p × t))

where x = scale thickness
      k_p = parabolic rate constant (function of T and alloy)
      t = time
```

Parabolic kinetics = GOOD. The thicker the scale, the more protection it provides (longer diffusion path). This is how Cr₂O₃ and Al₂O₃ protect.

**Linear growth (no protection, scale spalls or is porous):**
```
x = k_L × t

Linear kinetics = BAD. Constant corrosion rate regardless of how thick the scale gets.
Porous iron rust (Fe₂O₃) can follow linear kinetics — doesn't protect.
```

**Logarithmic growth (very thin, dense films at low T):**
Rapid initial growth then slows dramatically. Characteristic of passive films on Al, Ti, Ni at room temperature.

---

## 4. The Pilling-Bedworth Ratio

The **Pilling-Bedworth ratio (PBR)** predicts whether an oxide scale will be protective:

```
PBR = Volume of oxide formed / Volume of metal consumed

PBR = (M_oxide / ρ_oxide × n) / (M_metal / ρ_metal)

where M = molecular weight, ρ = density, n = number of metal atoms per formula unit
```

**Interpretation:**
```
PBR < 1:  Oxide volume < metal volume → oxide layer is POROUS → not protective
          Examples: Mg (PBR=0.81), Na, K
          
PBR = 1–2: Oxide well-matched → PROTECTIVE, compact, adherent
          Examples: Cr (PBR=2.0), Al (PBR=1.28), Ti (PBR=1.76)
          Best range for protective scales!
          
PBR > 2:  Oxide volume >> metal volume → scale too large → COMPRESSIVE STRESSES → 
          scale CRACKS and SPALLS → not fully protective (fresh metal exposed)
          Examples: Fe₂O₃ (PBR=2.1) — marginal; V₂O₅ (PBR=3.25) — bad
```

**Practical consequences:**
- Pure Al forms Al₂O₃ (PBR=1.28): excellent protection → doesn't rust at room T
- Pure Mg forms MgO (PBR=0.81): powdery, porous scale → corrodes in moist air
- Iron forms Fe₂O₃ / Fe₃O₄ (mixed): partially protective but eventually spalls

---

## 5. Wet (Electrochemical) Corrosion — Overview

Wet corrosion is an **electrochemical** process. It requires:
1. An **anode** (oxidation site): M → M^n+ + ne⁻ (metal dissolution)
2. A **cathode** (reduction site): e⁻ consumed (O₂ reduction or H⁺ reduction)
3. An **electrolyte** connecting them (water with dissolved ions)
4. A **metallic path** for electrons to flow from anode to cathode

```
Wet corrosion cell (iron in aerated water):

      Electrolyte (water + O₂)
     ────────────────────────────
   Fe²⁺ →│         │← O₂ + H₂O
         │ ANODE   │ CATHODE
         │ Fe→Fe²⁺+2e⁻  │ O₂+2H₂O+4e⁻→4OH⁻
    Fe²⁺ + 2OH⁻ → Fe(OH)₂ → Fe₂O₃·H₂O (rust)
     ────────────────────────────
     ← electrons flow through metal →
```

**Common cathodic reactions:**
```
In neutral/alkaline, with dissolved O₂:
O₂ + 2H₂O + 4e⁻ → 4OH⁻  (oxygen reduction — most common in atmospheric corrosion)

In acidic solution:
2H⁺ + 2e⁻ → H₂  (hydrogen evolution — common in acid attack)
```

**Rate-limiting step:**
- Often limited by diffusion of O₂ to cathode (in neutral solution)
- Or by metal ion transport away from anode (in concentrated solution)
- Or by ohmic resistance of electrolyte

---

## 6. Passivation — The Metal's Defense

The most important concept in practical corrosion resistance:

**Passive film:** A thin (2–10 nm), dense, adherent oxide layer that forms spontaneously on certain metals (Al, Ti, Cr, Ni, Co, stainless steel). This film:
- Prevents further corrosion by blocking ion diffusion
- Has very high electrical resistance → stops electrochemical reaction
- Reforms instantly if mechanically damaged

**Stainless steel passivation:**
- ≥ 10.5% Cr → Cr₂O₃ film (enriched at surface due to preferential Cr oxidation)
- Film forms within milliseconds of exposure to oxygen
- pH-dependent: below ~pH 4 → film dissolves → active corrosion
- Cl⁻ → locally destroys film → pitting corrosion

**Flade potential:** The potential at which active → passive transition occurs. Above E_passive: passive. Below: active.

**The paradox of stainless steel:**
Stainless steel is MORE susceptible to pitting than regular steel in some environments because:
- Regular steel corrodes uniformly (passive film absent) → slow, predictable
- Stainless has passive film → if broken locally (by Cl⁻) → rapid localized attack (pitting) → unexpected failure

---

## 7. Pourbaix Diagrams — Stability in Aqueous Environments

A **Pourbaix diagram** maps the thermodynamically stable phase of a metal as a function of electrode potential E and pH of the solution:

```
Pourbaix diagram for iron (simplified):

E (V vs. SHE)
  │
1.0│                      Fe₂O₃  (passive)
   │                     ────────────────
0.5│         Fe³⁺        │ Fe₂O₃ forms │
   │          (dissolved) │ protective   │
 0 │──────────────────────│ scale        │────────────
   │   Fe²⁺              │              │
-0.5│   (dissolved,       │              │
   │    active           │              │
-1.0│    corrosion)       │              │
   │                                    
-1.5│─────────────────────────────────── HER line
   └─────────────────────────────────── pH →
     0       4       7      10      14
```

**Three regions:**
- **Immunity:** Metal thermodynamically stable (very negative E) — no corrosion
- **Corrosion:** Metal ions stable in solution — active dissolution
- **Passive:** Oxide/hydroxide stable → protective film → corrosion arrested

**Practical use:** Know your service environment (pH, potential) → find region on Pourbaix diagram → predict behavior.

---

## 8. Corrosion Rate Measurement

**Weight loss / coupon testing:**
```
Corrosion rate = (W_initial - W_final) / (Area × time)
Units: mg/cm²/day, or "mils per year" (mpy) in US industry

Typical rates:
  Mild steel in seawater: 0.1–0.5 mpy (with cathodic protection)
  Mild steel bare in seawater: 3–10 mpy
  Stainless 316L in seawater: 0.01–0.1 mpy (pitting may be worse)
  Ti Grade 2 in seawater: < 0.001 mpy (essentially immune)
```

**Electrochemical methods:**
- **Tafel plot:** Measure current vs. applied potential → extrapolate to corrosion potential → corrosion current → corrosion rate
- **Linear polarization resistance (LPR):** Fast, non-destructive; relates polarization resistance to corrosion current
- **Electrochemical impedance spectroscopy (EIS):** Characterize passive film thickness and defect density

**E_corr and i_corr:**
```
Mixed potential theory: E_corr = corrosion potential (no net current)
i_corr = exchange current at E_corr → proportional to corrosion rate

i_corr = B / R_p  (where B = Stern-Geary constant ≈ 26 mV, R_p = polarization resistance)
```

---

## 9. Practical Corrosion in Engineering

**Most important practical forms (detailed in Ch 30):**
- Uniform corrosion
- Galvanic corrosion (dissimilar metals)
- Pitting corrosion
- Crevice corrosion
- Intergranular corrosion
- Stress corrosion cracking (SCC)
- Erosion-corrosion
- High-temperature oxidation / hot corrosion (Ch 45)

**Economic impact by sector:**
| Sector | Annual US Cost |
|--------|---------------|
| Infrastructure | $22 billion |
| Transportation | $29 billion |
| Utilities | $47 billion |
| Government | $20 billion |
| Oil/gas | $17 billion |
| Total | ~$270 billion |

**The economic opportunity:**
NACE International estimates that 25–30% of corrosion costs could be avoided with existing technology and proper design/maintenance practices. The gap is knowledge and implementation, not science.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Thermodynamic driving force | ΔG_f for oxide; Ellingham diagram; most metals want to corrode |
| Dry corrosion kinetics | Parabolic (protective scale) vs. linear (porous scale) |
| Pilling-Bedworth ratio | 1–2 = protective scale; < 1 = porous; > 2 = spalling |
| Electrochemical corrosion | Anode (oxidation) + cathode (reduction) + electrolyte + path |
| Passivation | Dense 2–10 nm oxide film on Al, Ti, Cr, Ni, stainless → immunity |
| Pourbaix diagram | E vs. pH → immunity / corrosion / passive regions |
| Corrosion rate | mpy, or i_corr from electrochemistry |

---

## Exercises

1. Calculate the Pilling-Bedworth ratio for: (a) Mg forming MgO (M_Mg = 24.3, M_MgO = 40.3, ρ_Mg = 1.74, ρ_MgO = 3.60 g/cm³), (b) Al forming Al₂O₃ (M_Al = 27.0, M_Al₂O₃ = 102.0, ρ_Al = 2.70, ρ_Al₂O₃ = 3.99 g/cm³), (c) Fe forming Fe₂O₃ (M_Fe = 55.8, M_Fe₂O₃ = 160, ρ_Fe = 7.87, ρ_Fe₂O₃ = 5.24 g/cm³). Based on PBR alone, predict the corrosion behavior of each.

2. A steel pipe corrodes parabolicaly: x² = k_p × t with k_p = 10⁻¹⁴ m²/s at 600°C. Calculate: (a) scale thickness after 1 hour, 1 day, 1 month, 1 year. (b) The pipe wall is 5mm thick. If scale growth beyond 2mm causes scale spalling (PBR problem), how long before spalling starts? (c) At 700°C, k_p = 10⁻¹² m²/s. How does time to spalling compare at 600°C vs 700°C?

3. Mild steel and stainless steel 316L are both used in a seawater pipeline. Using the simplified Pourbaix diagram: (a) At pH 8 (seawater), E = -0.3 V vs. SHE for mild steel — in what region does it fall? (b) What cathodic protection potential (-0.80 V vs. SHE for steel in seawater) moves mild steel into which region? (c) Stainless steel 316L pits in seawater at Cl⁻ = 19,000 ppm. What property of Cl⁻ breaks the passive film? How does Mo (3% in 316L) help resist pitting?

4. Two coupons are tested in aerated 3.5% NaCl (seawater simulant) for 30 days at 25°C: (a) Carbon steel 1020: weight loss = 6.2 g over 100 cm². Calculate corrosion rate in mg/cm²/day and mpy (1 mpy = 0.0254 mm/year; steel density = 7.87 g/cm³). (b) Aluminum 5083: weight loss = 0.002 g over 100 cm². Calculate corrosion rate. (c) Does Al actually corrode slower than steel, or does the surface oxide protect it while superficially dissolving? How would you distinguish between uniform corrosion and pitting?

5. The Ellingham diagram shows ΔG_f for Al₂O₃ is more negative than Cr₂O₃ at all temperatures. (a) Thermodynamically, should Al reduce Cr₂O₃ (i.e., can Al steal oxygen from Cr oxide)? (b) Yet in CMSX-4 superalloy at 1,000°C (which contains both Al and Cr), what oxide actually forms on the surface? (c) Explain this apparent paradox (hint: kinetics vs. thermodynamics, and Al activity in the alloy vs. Cr activity). (d) What is the practical consequence for the oxidation protection system in turbine blades?
