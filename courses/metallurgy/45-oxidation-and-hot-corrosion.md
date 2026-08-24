# Chapter 45: Oxidation and Hot Corrosion in Superalloys

> **"A superalloy that cannot protect itself from its environment is just expensive scrap. The alloy strength at 1,000°C is irrelevant if the material dissolves in the gas stream within 1,000 hours. The art of superalloy environmental protection is making the alloy grow exactly the right oxide — thin, adherent, self-healing — and keeping that oxide intact for 30,000 hours in one of the most chemically aggressive atmospheres on earth."**

---

## Table of Contents

1. [The High-Temperature Oxidation Problem](#1-the-high-temperature-oxidation-problem)
2. [Oxidation Kinetics — Parabolic and Linear Laws](#2-oxidation-kinetics--parabolic-and-linear-laws)
3. [Oxides Formed by Superalloy Elements](#3-oxides-formed-by-superalloy-elements)
4. [Selective Oxidation — The Key Principle](#4-selective-oxidation--the-key-principle)
5. [The Al₂O₃ Scale — Why It's the Target](#5-the-alo-scale--why-its-the-target)
6. [Transient Oxidation — Getting to the Steady State](#6-transient-oxidation--getting-to-the-steady-state)
7. [Spallation — The Failure Mode of Protective Scales](#7-spallation--the-failure-mode-of-protective-scales)
8. [Type I Hot Corrosion — Sulfidation Attack](#8-type-i-hot-corrosion--sulfidation-attack)
9. [Type II Hot Corrosion — Low-Temperature Pitting](#9-type-ii-hot-corrosion--low-temperature-pitting)
10. [MCrAlY Coatings — The Protection System](#10-mcraly-coatings--the-protection-system)
11. [Aluminide Coatings](#11-aluminide-coatings)
12. [Platinum Modification](#12-platinum-modification)
13. [Interactions: Oxidation + Mechanical Loading](#13-interactions-oxidation--mechanical-loading)
14. [Summary](#summary)
15. [Exercises](#exercises)

---

## 1. The High-Temperature Oxidation Problem

All metals exposed to oxygen at high temperature will oxidize. The question is not WHETHER oxidation occurs, but at what RATE and with what CONSEQUENCES.

In a jet engine HPT:
- Gas temperature: 1,500–1,700°C
- Oxygen partial pressure: 15–20% of combustion gas
- Water vapor: 5–15% (combustion product) — accelerates many oxidation reactions
- Sulfur compounds: SO₂, SO₃ from fuel sulfur
- Temperature cycling: every flight is a thermal cycle

**Unprotected CMSX-4 at 1,100°C:**
Would oxidize at ~1 mg/cm²/h → lose ~500 μm of wall in 5,000 hours. A 2.5 mm wall would be 20% consumed → blade fails structurally.

**Properly protected blade at 1,100°C (TBC + MCrAlY):**
Metal temperature is reduced to ~950–1,000°C, protected by MCrAlY + Al₂O₃ scale:
- Oxidation rate: ~0.001 mg/cm²/h → ~0.5 μm in 5,000 hours → negligible

The protection system extends blade life from ~1,000 hours (unprotected) to 30,000 hours (fully coated) — a 30× improvement.

---

## 2. Oxidation Kinetics — Parabolic and Linear Laws

### Parabolic Oxidation (Protective Scales)

When the oxide layer acts as a diffusion barrier:
```
x² = k_p × t

where x = oxide thickness (or mass gain Δm/A)
      k_p = parabolic rate constant (temperature-dependent)
      t = time
```

The rate DECREASES with time because the growing oxide layer makes it harder for reactants to diffuse to the reaction site. This is the desirable regime — protective scales grow parabolically.

Parabolic rate constant temperature dependence:
```
k_p = A × exp(-Q_p / RT)

where Q_p = activation energy for diffusion through the oxide
```

For Al₂O₃ scale growth: Q_p ≈ 250–350 kJ/mol → very slow growth at T < 900°C.

### Linear Oxidation (Non-Protective)

When the oxide does not form a coherent barrier (e.g., it cracks, peels, or is volatile):
```
x = k_l × t
```

Rate is CONSTANT with time — the metal is consumed at a steady rate. This is the catastrophic case (Type I hot corrosion, unprotected Fe/Ni at high T, volatile oxidation products).

### Breakaway Oxidation

Transition from parabolic to linear — the protective scale breaks down:
```
x        parabolic                  linear
         ─────────────         ────────────────────
t:       safe regime           ← breakaway → catastrophic
```

Breakaway occurs when the protective Al₂O₃ scale spalls off (from thermal cycling or scale growth stresses), exposing fresh alloy → rapid oxidation starts again from scratch → if breakaway happens every cycle, effective rate approaches linear.

---

## 3. Oxides Formed by Superalloy Elements

Each element in CMSX-4 has a tendency to form its own oxide:

| Element | Oxide | Growth Rate | Protectiveness | Note |
|---------|-------|------------|----------------|------|
| Ni | NiO | Fast | Poor | Forms at >600°C |
| Cr | Cr₂O₃ | Moderate | Good | Volatile as CrO₃ above 1,000°C |
| Al | Al₂O₃ (α) | Very slow | Excellent | The target oxide |
| Co | CoO | Fast | Poor | Similar to NiO |
| Ti | TiO₂ | Moderate | Poor | Grows fast, rutile structure |
| Ta | Ta₂O₅ | Moderate | Moderate | Mixed scale |
| W | WO₃ | Moderate | Volatile above 900°C | WO₃ volatile → can't protect |
| Mo | MoO₃ | Fast | Volatile above 700°C | Catastrophically non-protective |

**The problem:** CMSX-4 contains Ni, Co, W, Mo, Re (all poor oxidation resistors) and only limited Al (5.6%) and Cr (6.5%). Unprotected, the alloy would form a mix of NiO + CoO + Cr₂O₃ + WO₃(volatile) → rapidly degrading, non-protective scale.

**The solution:** Selective oxidation of Al to form ONLY Al₂O₃.

---

## 4. Selective Oxidation — The Key Principle

**Wagner's criterion for selective oxidation of component B in alloy A-B:**

Component B (Al) will form a continuous external scale if its concentration is above a minimum:
```
N_B* = (π g* / 12) × √(N_O^(s) × D_O / ν × D_B)

where N_B* = minimum mole fraction of B for external oxide formation
      D_O, D_B = diffusivities of oxygen and B in the alloy
      N_O^(s) = oxygen concentration at surface
      ν = stoichiometric coefficient
      g* = critical volume fraction for oxide percolation
```

For practical purposes, the key insight is:
- **Aluminum activity must be high enough** to maintain the Al chemical potential above the Al₂O₃ stability threshold
- **Aluminum must diffuse to the surface fast enough** to repair any scale damage before base-metal oxides form
- **Chromium "getters" oxygen** (forms Cr₂O₃ below the outer Al₂O₃) reducing the oxygen activity at the alloy surface → helps Al form its oxide at lower Al concentrations

The **Cr + Al synergy** is fundamental: Cr₂O₃ reduces oxygen activity at the metal surface, making it easier for Al (which has even stronger oxide affinity than Cr) to selectively oxidize. This allows alloys with 5–6% Al (too low for Al₂O₃ formation without Cr) to still form protective Al₂O₃ scales in the presence of 5–10% Cr.

---

## 5. The Al₂O₃ Scale — Why It's the Target

**α-Al₂O₃ (corundum) properties:**

| Property | Value | Significance |
|----------|-------|-------------|
| Crystal structure | Rhombohedral (corundum, same as ruby/sapphire) | Dense → good barrier |
| Density | 3.98 g/cm³ | Dense, no porosity |
| Growth rate at 1,100°C | ~0.05 μm/100h | 100× slower than NiO |
| Thermal expansion | 8.1×10⁻⁶/K | Good match with Ni alloys |
| Mechanical strength | Very high (HV ~2,000) | Resists abrasion |
| Adhesion to MCrAlY | Excellent (with Y) | Spallation-resistant |
| Oxygen diffusivity in Al₂O₃ | ~10⁻²⁰ m²/s at 1,100°C | Extreme barrier to O transport |

**Why Al₂O₃ grows so slowly:**

Oxygen must diffuse through the growing Al₂O₃ scale to reach the alloy surface (where it reacts with Al). The oxygen diffusivity in α-Al₂O₃ is extremely low (~10⁻²⁰ m²/s) because:
- Al₂O₃ is an ionic crystal with large bandgap
- Very few oxygen vacancies (necessary for diffusion) at equilibrium
- Yttrium additions further reduce vacancy concentration

This low diffusivity directly gives the parabolic rate constant k_p its small value → slow scale growth → long protection lifetime.

**Comparison at 1,100°C:**

| Oxide | Oxygen diffusivity | Scale thickness after 5,000h | Protective? |
|-------|-------------------|------------------------------|-------------|
| NiO | ~10⁻¹⁴ m²/s | ~500 μm | No |
| Cr₂O₃ | ~10⁻¹⁷ m²/s | ~50 μm | Marginal (volatile above 1,000°C) |
| α-Al₂O₃ | ~10⁻²⁰ m²/s | ~5 μm | Yes — excellent |

---

## 6. Transient Oxidation — Getting to the Steady State

When a fresh alloy surface is first exposed to oxygen, it does NOT immediately form pure Al₂O₃. There is a **transient oxidation** period:

```
Initial exposure → all elements at surface compete to oxidize:
  - NiO, CoO, Cr₂O₃, Al₂O₃, TiO₂ all nucleate and grow simultaneously
  
This mixed transient oxide is much less protective than pure Al₂O₃.
Al gradually "wins" the thermodynamic competition as NiO and CoO reduce 
  (Al₂O₃ more stable) and Cr₂O₃ getter effect activates.

Steady state (hours to days):
  - Continuous Al₂O₃ scale covers the surface
  - Other oxides disappear from surface (reduced by Al internally)
  - Protection is now excellent
```

This transient period is critical at startup after repair or TBC spallation. The blade is most vulnerable during transient oxidation — the protective state hasn't been established yet.

**Implication for engine operation:** Repeated starts/stops (many cycles) cause more transient oxidation damage than the same time in steady operation, because:
- Each shutdown destroys/removes some Al₂O₃ by scale spallation
- Each startup requires transient oxidation period before re-establishing protection

---

## 7. Spallation — The Failure Mode of Protective Scales

**Scale spallation** = the protective oxide detaches from the metal surface. This immediately exposes fresh, unoxidized alloy → rapid oxidation restarts → accelerated life consumption.

### Why Scales Spall

The CTE mismatch between Al₂O₃ (8.1×10⁻⁶/K) and the alloy (13.5×10⁻⁶/K) creates thermal stress during cooling:

```
Δα = 13.5 - 8.1 = 5.4 × 10⁻⁶/K

After cooling from 1,100°C to 20°C (ΔT = 1,080K):
Δε = Δα × ΔT = 5.4×10⁻⁶ × 1,080 = 0.0058 (0.58% strain mismatch)

Stress in scale (tensile in plane, since alloy contracts more than scale):
σ = E_scale × Δε = 400 GPa × 0.0058 = 2,320 MPa
```

This 2.3 GPa stress is well above the fracture strength of Al₂O₃ (typically 200–500 MPa in tension for a thin film). However, Al₂O₃ is in COMPRESSION during cooling (the alloy contracts faster → squeezes the scale), which actually helps — compression is favorable.

**Spallation typically occurs during COOLING or at room temperature**, when the scale is in compression and buckles → delamination.

**Yttrium effect on spallation resistance:**

Without Y in the bond coat: Al₂O₃ grows as columnar grains with grain boundaries perpendicular to the surface → void formation at grain boundaries → scale detaches easily at grain boundaries.

With Y: Y segregates to oxide grain boundaries → pins them → prevents void formation → scale remains adherent → dramatic improvement in spallation resistance (10–100× longer scale adherence life with Y).

---

## 8. Type I Hot Corrosion — Sulfidation Attack

**Temperature range:** 800–950°C (with peak at ~850–900°C)

**Mechanism:**

1. NaCl from sea air (or NaOH from HCl ingestion) + SO₃ from fuel sulfur → Na₂SO₄ deposits on blade surface
2. Na₂SO₄ is liquid above 884°C → molten salt deposit
3. Molten Na₂SO₄ reacts with and dissolves the protective Cr₂O₃ scale:
   ```
   Na₂SO₄ + Cr₂O₃ → 2 NaCrO₂ + SO₃  (or 2 Na₂CrO₄ + S² depending on conditions)
   ```
4. With Cr₂O₃ dissolved, fresh alloy contacts the Na₂SO₄ melt
5. Sulfur from Na₂SO₄ reacts with alloy:
   ```
   Ni → NiS₂ (sulfide)
   Cr → Cr₂S₃ (sulfide)
   ```
6. These sulfides have low melting points → liquid sulfide phases form within the alloy → grain boundary penetration

**Result:** Rapid, catastrophic oxidation + sulfidation. Rate is 10–100× faster than pure oxidation.

**Microstructural signature:** Deep sulfide particles (CrS, NiS, Co₄S₃) found BELOW the main oxide scale, embedded in the alloy → unambiguous sign of Type I hot corrosion in failure investigation.

**Threshold conditions:**
- Must have both: Na source AND S source → Na₂SO₄ deposit
- Na levels: > 0.008 ppm Na in the gas stream can cause Type I
- S levels: > 0.5 ppm S

**Naval aircraft environment:** Seawater aerosol → high NaCl ingestion → severe Type I hot corrosion risk. Marine gas turbines require >10% Cr in alloys/coatings vs. 5–8% for aviation.

---

## 9. Type II Hot Corrosion — Low-Temperature Pitting

**Temperature range:** 650–800°C (peak ~680–750°C)

**Mechanism:**

Type II doesn't require liquid Na₂SO₄ (solid at these temperatures). Instead:
- Na₂SO₄ + SO₃ at lower temperatures → NiSO₄ (liquid at 672°C, melting point) forms at the scale/metal interface
- The NiSO₄ liquid dissolves the alloy at pits
- Morphology: discrete circular pits rather than general attack

**Signature:** Discrete, circular pits in the blade surface, typically with white/yellow sulfate deposits in the pit. No deep sulfide penetration (distinguish from Type I).

**Different control factors from Type I:**
- SO₃ partial pressure: higher SO₃ → more severe Type II
- Na activity: same requirement
- Temperature: must be in 650–800°C range

**Aircraft vulnerability:** Low-altitude operation (lower gas temperature in some regions of the gas path), especially initial and final approach profiles with high SO₃ levels from industrial areas.

---

## 10. MCrAlY Coatings — The Protection System

Since CMSX-4 itself cannot provide adequate oxidation/hot corrosion resistance at elevated temperatures, an external **overlay coating** is applied.

**MCrAlY composition:**
- M = Ni, Co, or NiCo (matrix metal)
- Cr = 15–25% (for hot corrosion resistance)
- Al = 8–15% (reservoir for Al₂O₃ scale formation)
- Y = 0.3–1.0% (for scale adhesion)

**Why higher Cr than alloy:**
MCrAlY can contain 15–25% Cr (vs. 6.5% in CMSX-4) because MCrAlY doesn't need the creep resistance that limits Cr in the base alloy. This extra Cr provides:
- Better hot corrosion resistance (Cr₂O₃ primary defense against Type I)
- Better "getter" effect → helps Al₂O₃ form at lower Al levels

**MCrAlY as Al reservoir:**

The MCrAlY coating is a **sacrificial layer** — its purpose is to contain a large reservoir of Al. During service:
1. Al diffuses to the coating surface → forms α-Al₂O₃ (the TGO)
2. Coating loses Al over time
3. When Al in coating drops below ~5–6% → can no longer maintain Al₂O₃ → NiO, CoO form → rapid oxidation
4. Coating is "exhausted" → must replace blade at this point

Typical MCrAlY Al reservoir: 10 wt% Al in 100 μm thick coating = 100 μg Al per cm² of blade surface. This supports ~30,000 hours of α-Al₂O₃ scale growth.

**MCrAlY deposition methods:**

**Low-pressure plasma spray (LPPS):**
- Plasma jet at vacuum (< 50 mbar) or low pressure
- MCrAlY powder (~50 μm) injected → melts → impacts blade → adheres
- Dense, low-oxide coating (vacuum prevents oxidation during spray)
- Standard for MCrAlY bond coats in aviation

**High-velocity oxy-fuel (HVOF):**
- Fuel combustion drives powder at supersonic velocity
- Dense, low-porosity coating
- Used when LPPS not available or for industrial GT applications

**Electron beam physical vapor deposition (EB-PVD) of MCrAlY:**
- Rare for MCrAlY; more common for aluminides
- Used when very dense, thin coating needed

---

## 11. Aluminide Coatings

An alternative to MCrAlY for oxidation protection: **aluminide coatings** formed by diffusing aluminum into the blade surface.

**Pack aluminizing process:**
1. Blades packed in powder mixture: Al source (AlCl₃ vapor from Al + activator like NH₄Cl), ceramic filler (Al₂O₃ or Al powder)
2. Heated to 700–1,100°C: AlCl₃ vapor diffuses to blade surface → Al deposits → Ni from alloy and Al co-diffuse → NiAl intermetallic forms
3. Result: ~50–100 μm thick NiAl coating on blade surface

**NiAl intermetallic properties:**
- Strong Al₂O₃ former (high Al activity)
- High melting point: 1,638°C (higher than alloy)
- Excellent oxidation resistance (similar to MCrAlY)
- **Limitation:** Very brittle at room temperature → susceptible to mechanical damage during handling

**Platinum-aluminide (Pt-Al):**
Adding Pt to aluminide dramatically improves performance:
1. Electroplate 5–10 μm Pt on blade
2. Pack aluminize: Pt interdiffuses with Al and Ni → (Ni,Pt)Al intermetallic
3. Pt at grain boundaries → improves scale adhesion (Pt in grain boundaries prevents void formation)
4. Result: Much longer coating life (2–3× vs plain aluminide)

Pt-aluminide is used on the highest-temperature HPT blades in most commercial engines.

---

## 12. Platinum Modification

Platinum's role in improving oxidation/coating performance is multifaceted:

**Mechanism 1 — Scale adhesion:**
Pt segregates to the α-Al₂O₃ grain boundaries during scale growth. This:
- Increases the grain boundary cohesive energy → less void formation at boundaries → improved scale adhesion (5–10× fewer voids)
- The "Reactive Element Effect" (similar to Y in MCrAlY, but Pt works differently — it doesn't react, it pins)

**Mechanism 2 — Alloy reservoir effect:**
Pt in the bond coat raises the Al activity → the thermodynamic driving force for selective Al oxidation increases → Al₂O₃ forms more readily even as Al is depleted from the coating.

**Mechanism 3 — Phase stabilization:**
Pt stabilizes β-NiAl phase in the bond coat (more Al-rich than γ/γ′ alloy) → the β phase provides a sustained high-Al reservoir throughout coating service.

**Cost of Pt:**
Pt costs ~$30,000–35,000/kg. A blade with 5–10 μm of Pt has ~0.02–0.04 g of Pt per blade = ~$0.6–1.2 per blade in Pt cost. Negligible vs. $10,000+ blade value, but significant in aggregate for a production run.

---

## 13. Interactions: Oxidation + Mechanical Loading

Oxidation and mechanical loading are not independent — they interact to produce accelerated damage:

**Oxide-Induced Crack Closure (OICC):**
Oxide forms inside fatigue cracks → wedges the crack open → prevents crack closure → higher effective stress intensity at crack tip → faster crack growth. Especially important for HCF cracks at film holes in oxidizing environment.

**Stress-Assisted Grain Boundary Oxidation (SAGBO):**
Tensile stress at grain boundaries (in PC/DS alloys) accelerates oxygen penetration along grain boundaries. Not relevant for SX blades (no grain boundaries), but critical for DS blades where longitudinal boundaries exist.

**Hot Corrosion + Fatigue:**
Corrosion pits (from Na₂SO₄ attack) act as fatigue crack initiation sites. The pit stress concentration factor Kt ≈ 3–5 → dramatically reduces the stress amplitude needed for fatigue crack nucleation.

**Oxidation + Creep (Oxidation-Enhanced Creep):**
Oxygen diffusing into the alloy dissolves the γ′ phase locally (Al preferentially oxidizes → depletes Al → γ′ dissolves). γ′-depleted zones have much lower creep resistance → create soft zones → enhanced local creep deformation → crack formation.

---

## Summary

- **Oxidation** of bare Ni superalloy is rapid and catastrophic → must form PROTECTIVE scale
- **Target oxide**: α-Al₂O₃ — extremely slow growth rate, dense, adherent (with Y modification)
- **Selective oxidation**: achieved by Cr + Al synergy — Cr reduces oxygen activity at surface → Al₂O₃ forms selectively
- **Parabolic kinetics**: protective scale; rate decreases with time as scale thickens
- **Scale spallation**: CTE mismatch drives spallation on cooling → yttrium in bond coat prevents by pinning oxide grain boundaries
- **Type I hot corrosion**: 800–950°C; molten Na₂SO₄ dissolves Cr₂O₃ → sulfidation attack; requires Cr content >10% in coating
- **Type II hot corrosion**: 650–800°C; pitting by liquid NiSO₄; discrete pits, different from Type I
- **MCrAlY bond coat**: overlay coating with 15–25% Cr + 8–12% Al + Y; sacrificial Al reservoir; deposited by LPPS or HVOF
- **Aluminide coatings**: diffusion coating → NiAl layer; brittle but good Al₂O₃ former
- **Platinum modification**: improves scale adhesion 5–10×; enables longer life

**Next chapter:** The complete coating system — how MCrAlY + TBC work together, the bond coat failure sequence, and the next-generation coating development.

---

## Exercises

1. Parabolic oxidation: α-Al₂O₃ scale grows on MCrAlY at 1,050°C with k_p = 0.002 (μm)²/h. (a) What is scale thickness after 500h, 5,000h, 30,000h? (b) If the coating has 10 wt% Al in 120 μm thick layer (density 6.5 g/cm³), and each μm of α-Al₂O₃ consumes 0.8 mg/cm² of Al, what is the Al reservoir in the coating? (c) At 30,000h scale thickness, how much Al has been consumed? Has the Al reservoir been exhausted?

2. Type I hot corrosion threshold: Na₂SO₄ deposits form when gas Na > 0.008 ppm AND S > 0.5 ppm. A naval helicopter ingests 0.5 ppm Na from sea spray and operates near industrial zones with 2 ppm SO₂. (a) Are conditions met for Type I hot corrosion? (b) The engine has a component operating at 870°C — is this in the Type I temperature range? (c) What alloying or coating strategy would you recommend?

3. CTE mismatch calculation: Bond coat α = 13.5×10⁻⁶/K, α-Al₂O₃ scale α = 8.1×10⁻⁶/K. Engine operates at 1,100°C and cycles to 25°C. (a) Calculate the strain mismatch Δε during a full cooling cycle. (b) If the scale is 5 μm thick and acts as a thin film (E = 400 GPa), what is the maximum in-plane compressive stress in the scale at room temperature? (c) Buckling of a thin film occurs when σ_compressive > (π²/3) × E_film × (t/a)² where t = film thickness and a = half-wavelength of buckle. For a = 50 μm and t = 5 μm, what stress is needed for buckling? Does the thermal stress exceed this?

4. MCrAlY life assessment: A MCrAlY coating starts with 10 wt% Al. After 15,000 hours at 1,000°C, the Al content has dropped to 5.5% (measured by EPMA). Below 4% Al, the coating can no longer form protective Al₂O₃. (a) Using linear depletion rate (reasonable approximation for diffusion-limited depletion), estimate remaining coating life. (b) If the operator wants a 4,000h safety margin, when should the blade be removed for inspection/recoating?

5. Selective oxidation: CMSX-4 has 5.6 wt% Al and 6.5 wt% Cr. A researcher proposes a new alloy with 3 wt% Al and 15 wt% Cr. Using the concept that Cr's "getter effect" allows protective Al₂O₃ formation at lower Al contents, explain qualitatively whether this alloy would likely form protective Al₂O₃, a mixed scale, or Cr₂O₃. What would be the trade-off in terms of mechanical properties vs. oxidation resistance?
