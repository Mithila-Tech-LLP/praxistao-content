# Chapter 32: High-Temperature Oxidation — When Metal Meets Hot Gas

> **"A turbine blade in a jet engine is continuously bathed in gas at 1,500°C, far above the melting point of any metal. Yet the blade survives because of a thin, dense oxide layer — barely a micrometer thick — that forms between the metal and the hot gas. This Al₂O₃ layer is the reason modern aviation is possible. Understanding how it forms, why it fails, and what damages it is the heart of high-temperature materials science."**

---

## Table of Contents

1. [Why High-Temperature Oxidation Is Different](#1-why-high-temperature-oxidation-is-different)
2. [Wagner's Theory of Oxidation Kinetics](#2-wagners-theory-of-oxidation-kinetics)
3. [Growth of the TGO — A Deep Dive](#3-growth-of-the-tgo--a-deep-dive)
4. [Selective Oxidation — Why Al₂O₃ Forms](#4-selective-oxidation--why-alo-forms)
5. [Scale Adhesion and Spallation](#5-scale-adhesion-and-spallation)
6. [Hot Corrosion (Type I and Type II)](#6-hot-corrosion-type-i-and-type-ii)
7. [Oxidation in Superalloys — System View](#7-oxidation-in-superalloys--system-view)
8. [Protective Oxide Ranking](#8-protective-oxide-ranking)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why High-Temperature Oxidation Is Different

High-temperature (T > 500°C) oxidation differs from ambient corrosion:
- Solid-state ion diffusion through oxide scale (not aqueous electrochemistry)
- Scale must physically protect the metal at high temperature
- Thermal cycling causes scale cracking/spallation (differential thermal expansion)
- Multiple oxidizing species: O₂, H₂O, SO₂, NaCl → different attack modes
- Oxide scales can melt (MoO₃ bp = 1,155°C, NiO mp = 1,955°C)
- Gas composition affects equilibrium oxide stability (pO₂, pS₂)

**The fundamental challenge:** The oxide must be:
- Thermodynamically stable (ΔG_f very negative)
- Kinetically slow-growing (low diffusivity through oxide)
- Mechanically adherent (doesn't spall)
- Maintained during service cycles and environmental attacks

---

## 2. Wagner's Theory of Oxidation Kinetics

Carl Wagner (1933) derived the parabolic rate constant from first principles:

**The oxidation model:**
Oxide grows by diffusion of ions through the existing oxide layer. The rate is limited by:
- Metal cations diffusing outward through the oxide to react with O₂ at the outer surface
- O²⁻ diffusing inward to react with metal at the metal/oxide interface
- Electronic conduction (oxide must be a semiconductor for charged species to move)

**Parabolic rate law derivation:**
```
Scale thickness x grows by:
  dx/dt = k / x  (rate = diffusion flux × area = D × dc/dx ≈ D × ΔC/x)
  
Integrating:
  x² = 2k × t = k_p × t   ← parabolic law
  
k_p = 2 × D_ion × C_metal × V_molar / N_A
```

**Temperature dependence of k_p (Arrhenius):**
```
k_p = A × exp(-Q / RT)

where Q = activation energy for ion diffusion in oxide
      A = pre-exponential constant

For Al₂O₃: Q ≈ 420 kJ/mol  (high → very slow growth even at 1,000°C)
For NiO:   Q ≈ 160 kJ/mol  (lower → faster growth)
For FeO:   Q ≈ 120 kJ/mol  (fast growth)
```

**Key insight:** Alloys that form Al₂O₃ (Q = 420 kJ/mol) are 100–1,000× more oxidation-resistant than those forming NiO or FeO at the same temperature. This is WHY nickel superalloys are aluminized.

---

## 3. Growth of the TGO — A Deep Dive

The **Thermally Grown Oxide (TGO)** on MCrAlY bond coats was introduced in Ch 46. Here we examine the kinetics and failure in detail:

**Two-stage TGO growth:**

**Stage 1 (< 200 hours):** θ-Al₂O₃ (metastable) forms first:
- θ-Al₂O₃: needle/platelet morphology; mixed inward/outward growth
- Grows faster than α-Al₂O₃
- Creates an irregular, mechanically imperfect layer

**Stage 2 (onset: 200–1,000 hours depending on T):** θ → α-Al₂O₃ transformation:
- α-Al₂O₃: corundum structure; grows primarily inward (O²⁻ diffusion controlled)
- Much slower growth rate (10–100× slower than θ)
- Dense, equiaxed grains
- **This is the stable, protective phase**

```
TGO thickness vs. time (simplified):

x (μm)
5─┼──────────────────── θ→α transition region
  │        α phase (slowly thickening)
3─┼───────────────────────────────────────── LIFE LIMIT ~5 μm
  │       ← θ phase
1─┼──────
  └────────────────────────────────── t (hours)
  0   200  500  1000  5000  10000  20000  30000
```

**TGO failure mechanism:**
- Bond coat oxidizes → Al depleted → eventually: NiO, Cr₂O₃ spinels form → less protective
- TGO stresses: in-plane compressive stress (due to oxide volume > metal consumed → PBR effect); at service T: creep relaxation; on cooling: large compressive stress (oxide thermo-elastic mismatch)
- After many thermal cycles: rumpling, ridging → TBC detaches
- Life-limiting thickness: ~5–7 μm (interface energy decreases below thermal stress → delamination)

---

## 4. Selective Oxidation — Why Al₂O₃ Forms

In a multi-component alloy (CMSX-4: Ni-Cr-Co-W-Ta-Re-Al-Ti-Hf), WHY does Al₂O₃ form preferentially instead of NiO or Cr₂O₃?

**Thermodynamics:** At 1,000°C, ΔG_f per mole O₂ for:
```
Al₂O₃:  -836 kJ/mol O₂  ← most negative (most stable)
Cr₂O₃:  -636 kJ/mol O₂
NiO:    -434 kJ/mol O₂
```

Thermodynamics says Al₂O₃ should form first → but the alloy has only ~12 at% Al (in CMSX-4) vs 70+ at% Ni.

**How sufficient Al activity is maintained:** The MCrAlY bond coat provides a high Al reservoir (8–12 wt% Al). If Al activity > critical value (a_Al,critical), Al₂O₃ nucleates preferentially.

**Transition from Al₂O₃ to NiO formation:**
As bond coat ages, Al depletes. Once Al concentration drops below β-phase → only γ matrix (5–6 at% Al) → insufficient thermodynamic driving force to maintain exclusive Al₂O₃ → mixed NiAl₂O₄ spinels + NiO form → "breakaway oxidation" → rapid attack.

**The 3%Al Rule:** 3 wt% Al minimum in the alloy at the surface is roughly needed to form a continuous Al₂O₃ scale. Below this, Cr₂O₃ + NiO form instead.

---

## 5. Scale Adhesion and Spallation

Even if the oxide is slow-growing and thermodynamically stable, it must STAY ON the metal.

**Spallation:** Scale detaches from metal → fresh metal exposed → rapid reoxidation → metal consumed rapidly.

**Mechanisms of spallation:**

1. **Thermal cycling (most common):**
   CTE of Al₂O₃ (8.3 × 10⁻⁶/°C) vs. NiAl bond coat (13–14 × 10⁻⁶/°C) → on cooling: oxide under compressive stress = oxide wants to buckle → wrinkles form → eventually delamination.

2. **Growth stress:**
   New oxide forms within the scale → volume increase → growth stress → can buckle/crack scale.

3. **Void formation:**
   If metal oxidizes faster than O²⁻ can diffuse inward, Kirkendall voids form at metal/oxide interface → weakens adhesion.

**The Reactive Element Effect (Y, Hf, Ce):**
Adding small amounts of Y (0.1–1%) or Hf (0.1–0.5%) to MCrAlY alloys DRAMATICALLY improves scale adhesion:
- **Mechanism 1 (grain boundary segregation):** Y segregates to Al₂O₃ grain boundaries → reduces grain boundary diffusion → slows scale growth
- **Mechanism 2 (anchor pegs):** Y₂O₃ or HfO₂ "pegs" form at metal/oxide interface → mechanical interlocking → resist spallation
- **Mechanism 3 (change growth mechanism):** Suppresses outward cation diffusion → only inward O²⁻ diffusion → oxide grows at metal/oxide interface → less growth stress

**Effect:** 0.1% Y can increase oxide life from 100 hours to 1,000+ hours at 1,100°C.

---

## 6. Hot Corrosion (Type I and Type II)

**Hot corrosion** = accelerated attack when molten sulfate deposits (Na₂SO₄) are present on the blade surface (from sea salt or S-contaminated fuel + Na in air). Covered quantitatively in Ch 45, but mechanism summary here:

### Type I Hot Corrosion (750–950°C)

**Na₂SO₄ deposits are liquid** above their melting point (884°C):
- Liquid Na₂SO₄ **fluxes** the normally protective Cr₂O₃ scale:
  - Acidic fluxing: Cr₂O₃ + SO₃ → 2CrO₃ + SO₂ (at high pSO₃)
  - Basic fluxing: Cr₂O₃ + Na₂SO₄ → 2NaCrO₂ + SO₃ (at high Na₂O activity)
- Once Cr₂O₃ is removed → accelerated sulfidation attack
- S penetrates into alloy → Cr₂S₃, NiS form → complex sulfide needles deep in alloy
- **Signature:** Deep internal sulfide precipitates in SEM cross-section

### Type II Hot Corrosion (650–800°C)

Lower temperature, different deposits:
- Na₂SO₄ + CoSO₄ or NiSO₄ form low-melting eutectics (some melt as low as 565°C)
- Pitting attack without the deep sulfidation
- **Signature:** Hemispherical pits on surface

**Protection against hot corrosion:**
- High-Cr alloys or coatings (Cr₂O₃ more resistant than Al₂O₃ to Type I fluxing — up to a point)
- Aluminide + Pt modification → (Ni,Pt)Al → slows Na₂SO₄ dissolution
- MCrAlY with high Cr (>20%) → better Type I resistance
- Avoid S-contaminated fuel + reduce Na ingestion (inlet air filtration)

---

## 7. Oxidation in Superalloys — System View

The oxidation protection of turbine blades is a multilayer system working together:

```
Gas (1,500°C, containing O₂, H₂O, SO₂, Na...)
     ↓
TBC (YSZ, 7%Y₂O₃-ZrO₂): 150–200 μm thermal barrier; T at TBC-TGO interface ≈ 1,000°C
     ↓
TGO (α-Al₂O₃, 1–10 μm): primary oxidation barrier; grows at TBC/bond coat interface
     ↓
Bond coat (MCrAlY or PtAl): 100–200 μm; Al reservoir; Cr for hot corrosion protection
     ↓
Interdiffusion zone: bond coat elements diffuse into blade alloy; SRZ forms (TCP risk)
     ↓
CMSX-4 blade alloy (Ni-based superalloy): structural material
```

**Each layer fails differently:**
- TBC: spalls from TGO stress concentration
- TGO: reaches critical thickness (5–7 μm) → delamination energy > adhesion energy
- Bond coat: Al depletion → breakaway oxidation
- Interdiffusion: SRZ brittle zone forms near bond coat/alloy interface

---

## 8. Protective Oxide Ranking

Ranking of protective oxide scales from best to worst for high-temperature service:

```
BEST (slowest growth, most adherent):
1. α-Al₂O₃ (alumina): k_p ≈ 10⁻¹⁶ g²/cm⁴·s at 1,000°C; excellent adherence with Y
2. Cr₂O₃ (chromia): k_p ≈ 10⁻¹⁴ at 1,000°C; volatile CrO₃ above 1,000°C
3. SiO₂ (silica): excellent but forms vitreous glass → susceptible to H₂O at high T (→ Si(OH)₄)
4. TiO₂: moderate, dissolves in alloy = less protective
5. NiO: faster growth; protective at moderate T (<850°C)
6. FeO/Fe₂O₃/Fe₃O₄: mixed; worst of engineering metals
WORST (fastest, least protective):
```

**Transition metals for > 1,200°C protection:**
Above 1,200°C, CrO₃ becomes volatile → Cr₂O₃ protection fails.
Only Al₂O₃ and SiO₂ remain stable and slow-growing.
Next generation coatings: Gd₂Zr₂O₇, La₂Zr₂O₇ pyrochlores — stable at 1,300–1,500°C (for next-generation TBC).

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Wagner kinetics | x² = k_p×t; k_p = A×exp(-Q/RT); Al₂O₃ has highest Q → slowest growth |
| TGO growth | θ-Al₂O₃ first (fast), then α-Al₂O₃ (slow, protective); life limit ~5-7 μm |
| Selective oxidation | ΔG_f most negative for Al₂O₃ → Al oxidizes preferentially if activity sufficient |
| Reactive elements | Y, Hf, Ce → improve adhesion via grain boundary segregation + anchor pegs |
| Hot corrosion | Na₂SO₄ fluxes protective oxide → Type I (750-950°C) sulfidation; Type II (650-800°C) pitting |
| Scale spallation | CTE mismatch on cooling; growth stress; Kirkendall voids → relieved by Y additions |
| Best oxide | α-Al₂O₃ > Cr₂O₃ > SiO₂ > NiO > FeO |

---

## Exercises

1. Calculate the TGO growth rate for α-Al₂O₃ at 1,000°C and 1,100°C using k_p = 10⁻¹⁶ × exp(+Q/RT) and Q = 420 kJ/mol (use R = 8.314 J/mol·K). (a) Time to reach 1 μm at 1,000°C and 1,100°C. (b) Time to reach 5 μm (life limit) at each temperature. (c) What is the ratio of service life at 1,000°C vs 1,100°C? (d) This illustrates why every 25°C TIT increase is so difficult — calculate the % increase in TGO growth rate for +25°C at 1,000°C.

2. A bond coat has 8 wt% Al initially. During 20,000 hours at 1,050°C, Al diffuses into the blade alloy AND is consumed forming Al₂O₃. Assume Al consumption rate = 0.3 mg/cm² per 1,000 hours (from TGO growth) + diffusion flux = 0.5 mg/cm² per 1,000 hours. Bond coat thickness = 150 μm, density = 7.5 g/cm³. (a) Calculate total Al available in bond coat (g/cm²). (b) Total Al consumed in 20,000 hours. (c) When does the bond coat reach the critical minimum Al (3 wt% → breakaway oxidation)? Is 20,000 hours achievable?

3. Type I hot corrosion test: Ni-based alloy coupon coated with Na₂SO₄ at 900°C in air for 24 hours. SEM shows deep sulfide penetration (Cr₂S₃, NiS) to 200 μm. (a) Write the reaction for basic sulfate fluxing of Cr₂O₃ by Na₂SO₄. (b) Why does S penetrate so deeply into the alloy? (c) Compare: same alloy at 700°C (Type II range) — predict the morphology difference. (d) Adding 30% Cr to the alloy (from 20% to 30%): explain why this improves Type I resistance. (e) What fuel composition change reduces hot corrosion risk?

4. An MCrAlY bond coat without reactive element additions spalls after 500 thermal cycles from 1,100°C to 25°C. The same composition with 0.4%Y survives 3,000 cycles. (a) Describe three mechanisms by which Y improves adhesion. (b) If the turbine operates at 12,500 ft cruising, each flight = 1 thermal cycle → without Y, blade life = 500 flights. With Y, life = 3,000 flights. Calculate fleet cost saving if each new blade costs $25,000 and a 200-engine fleet has 80 HPT blades per engine with 50% replacement at each overhaul (overhaul at blade life limit).

5. Compare the oxidation resistance of the following at 1,100°C: (a) bare CMSX-4 (12 at% Al, 6 at% Cr), (b) aluminide-coated CMSX-4 (NiAl with 30 at% Al at surface), (c) MCrAlY-coated CMSX-4 (15%Cr, 10%Al). For each: predict which oxide forms (Al₂O₃, Cr₂O₃, or mixed), estimate relative k_p (use: k_p ∝ exp(-Q/RT), Q_Al2O3 = 420, Q_Cr2O3 = 250 kJ/mol), and calculate relative oxide thickness after 1,000 hours. Rank the three systems by oxidation life.
