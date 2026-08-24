# Chapter 61: High-Entropy Alloys — Rethinking Alloy Design from the Ground Up

> **"For 150 years, alloy design meant taking a primary metal and adding small amounts of other elements. Steel is iron with a little carbon. Brass is copper with some zinc. Nimonic is nickel with some chromium and aluminum. In 2004, two independent research groups — Brian Cantor at Oxford and Jien-Wei Yeh at NTHU Taiwan — asked: what if you made an alloy with FIVE elements in equal amounts? What they found violated every prediction of classical alloy theory."**

---

## Table of Contents

1. [The Classical Alloy Design Paradigm — And Its Limits](#1-the-classical-alloy-design-paradigm--and-its-limits)
2. [The High-Entropy Concept — Configurational Entropy](#2-the-high-entropy-concept--configurational-entropy)
3. [The Four "Core Effects"](#3-the-four-core-effects)
4. [The Cantor Alloy — CoCrFeMnNi](#4-the-cantor-alloy--cocrfemnni)
5. [Properties of CoCrFeMnNi](#5-properties-of-cocrfemnni)
6. [Refractory High-Entropy Alloys (RHEA) — The Jet Engine Connection](#6-refractory-high-entropy-alloys-rhea--the-jet-engine-connection)
7. [RHEA Properties and Challenges](#7-rhea-properties-and-challenges)
8. [HEA vs. Superalloys — A Realistic Comparison](#8-hea-vs-superalloys--a-realistic-comparison)
9. [Precipitation-Hardened HEAs — Best of Both Worlds](#9-precipitation-hardened-heas--best-of-both-worlds)
10. [Composition Space — Exploring 10²⁵ Possible Alloys](#10-composition-space--exploring-10-possible-alloys)
11. [Where HEAs Are Going](#11-where-heas-are-going)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Classical Alloy Design Paradigm — And Its Limits

Classical metallurgy treats alloys as a **primary component with additions**:
- Steel = Fe + 0.1–2% C (+ optional Cr, Ni, Mo, V...)
- Brass = Cu + 10–40% Zn
- CMSX-4 = Ni + ~34% total alloying elements (but Ni is still the dominant component at ~66%)

This approach works well — 150 years of alloy development produced remarkable materials. But it has limits:
1. The primary element determines the basic crystal structure, melting point, and density
2. You can't change the fundamental character of the alloy without changing the primary element
3. The composition space around any one primary element is thoroughly explored
4. To go to higher temperatures (above Ni alloys), you need a different primary element

What if there were a fundamentally different approach?

---

## 2. The High-Entropy Concept — Configurational Entropy

**High-entropy alloys (HEA)** contain five or more principal elements, each in roughly equal concentrations (5–35 at% each). The defining characteristic is the **configurational entropy of mixing**:

```
ΔS_mix = -R × Σ(x_i × ln(x_i))

where x_i = mole fraction of element i
      R = gas constant (8.314 J/mol·K)
```

For equimolar alloys:
- 1 element: ΔS_mix = 0 (pure metal, no mixing entropy)
- 2 elements equimolar: ΔS_mix = R × ln(2) = 5.76 J/mol·K
- 3 elements equimolar: ΔS_mix = R × ln(3) = 9.13 J/mol·K
- 4 elements equimolar: ΔS_mix = R × ln(4) = 11.53 J/mol·K
- **5 elements equimolar: ΔS_mix = R × ln(5) = 13.38 J/mol·K**
- 10 elements equimolar: ΔS_mix = R × ln(10) = 19.14 J/mol·K

The **"high entropy" threshold** is typically defined as ΔS_mix > 1.5R ≈ 12.5 J/mol·K, which requires ≥5 equimolar elements.

### Why Entropy Matters

The Gibbs free energy of mixing determines whether a mixture forms one phase or separates:
```
ΔG_mix = ΔH_mix - T × ΔS_mix
```

For the mixture to be thermodynamically stable (single solid-solution phase), ΔG_mix < 0:
```
Need: T × ΔS_mix > ΔH_mix

If ΔS_mix is large (HEA), even if ΔH_mix is unfavorable (positive, meaning elements would rather be separate), at high temperature the entropy term dominates → single solid solution phase is stable
```

**The prediction:** HEAs should form simple solid solution phases (BCC, FCC, or HCP) rather than complex intermetallic compounds, because the high entropy stabilizes the disordered solid solution.

**The experimental surprise:** This prediction is largely correct — many HEAs DO form simple single-phase structures despite containing elements that would form many intermetallics in binary combinations.

---

## 3. The Four "Core Effects"

Yeh (2004) proposed that HEAs exhibit four fundamental effects not seen in conventional alloys:

### 1. High Entropy Effect
As described above: high ΔS_mix stabilizes simple solid solution phases → reduces tendency for intermetallic formation → simpler, more predictable microstructure.

### 2. Severe Lattice Distortion Effect
In a conventional alloy, a few "foreign" atoms (e.g., 5% Cr in Ni) cause minor lattice distortion. In HEA, EVERY atom is "foreign" in the sense that its neighbors are different elements with different atomic sizes. The lattice is severely distorted everywhere:

```
Conventional alloy: Ni - Ni - Ni - Cr - Ni - Ni (5% Cr)
                    uniform lattice, occasional small distortion

HEA: Co - Cr - Fe - Mn - Ni - Co - Cr - Fe (equimolar)
     EVERY site distorted relative to average lattice parameter
```

This severe distortion:
- Strongly impedes dislocation motion (dislocations must drag through a rough lattice)
- Produces high yield strength
- Reduces thermal conductivity (phonon scattering by distortion)

### 3. Sluggish Diffusion Effect
In conventional alloys, diffusion is limited by a few rate-controlling elements. In HEA, every atom must pass through an environment of 5+ different elements to diffuse — each diffusion jump has a different activation energy depending on the local configuration:

```
ΔG_activation varies at every lattice site → some sites are "easy" jumps, some are "hard" jumps
Overall diffusion = sum of many individually varying jumps → effectively slower
```

Result: diffusion is slower in HEAs → sluggish phase transformation → good for high-temperature stability.

### 4. Cocktail Effect
The final properties of the HEA are more than the sum of its elemental properties. Unexpected synergistic effects emerge from multi-principal-element combinations — properties that cannot be predicted by simple rule-of-mixtures.

**Note on "core effects":** These effects are now known to be variable — not all HEAs show all four effects equally strongly. They are useful conceptual guidelines, not universal laws.

---

## 4. The Cantor Alloy — CoCrFeMnNi

The first and most-studied HEA, identified by Brian Cantor at Oxford (2004):

**Composition:** Co₂₀Cr₂₀Fe₂₀Mn₂₀Ni₂₀ (equimolar, all at 20 at%)
**Crystal structure:** Single-phase FCC solid solution (at room temperature)
**Microstructure:** Similar to any FCC alloy — grains with grain boundaries, no precipitates

### Why CoCrFeMnNi Forms a Single Phase

All five elements (Co, Cr, Fe, Mn, Ni) have:
- Similar atomic radii (Ni: 124 pm, Co: 125 pm, Fe: 126 pm, Mn: 127 pm, Cr: 128 pm) — within 4% of each other
- Compatible crystal structures (FCC or BCC)
- Small mixing enthalpy between many pairs

These favorable conditions + high ΔS_mix → single-phase FCC stable from ~300°C to melting.

---

## 5. Properties of CoCrFeMnNi

CoCrFeMnNi has been extensively studied. Key properties:

### Room Temperature
- Yield strength: ~200 MPa (moderate — similar to austenitic stainless steel)
- UTS: ~500–600 MPa
- Ductility: >50% elongation (very ductile)
- Hardness: ~150 HV

**Not impressive** compared to age-hardened Ni superalloys (σ_y > 1,000 MPa). The single-phase solid solution is inherently moderate strength.

### Low Temperature — The Remarkable Property

**Cryogenic properties:** Unlike most FCC metals that become brittle at low temperature, CoCrFeMnNi shows INCREASING strength AND ductility as temperature decreases (from room T to -196°C):
- Yield strength: 200 MPa (RT) → 500 MPa (-196°C)
- Ductility: 50% (RT) → 70%+ (-196°C)

This is highly unusual — most metals either lose ductility or lose strength when cooled. The mechanism: at low temperatures, deformation twinning activates → provides additional deformation mechanism AND blocks dislocation motion → simultaneous strength + ductility improvement.

**Application implication:** CoCrFeMnNi is one of the best materials known for cryogenic vessels (LNG storage, liquid hydrogen, liquid nitrogen tanks). Better than austenitic stainless steel 316L at very low temperatures.

### Fracture Toughness

CoCrFeMnNi at -196°C: K_Ic > 200 MPa√m — among the highest fracture toughness of any material at any temperature. Exceptional for structural applications requiring defect tolerance.

### High Temperature Properties

**Below expectations for high-temperature applications:**
- Oxidation resistance: poor (no Al or Si → no protective oxide)
- Creep resistance: moderate → Mn evaporates above ~700°C (Mn is volatile)
- Yield strength at 800°C: ~100 MPa → inferior to any Ni superalloy

CoCrFeMnNi is NOT a high-temperature alloy.

---

## 6. Refractory High-Entropy Alloys (RHEA) — The Jet Engine Connection

The CoCrFeMnNi alloy system cannot compete with Ni superalloys at high temperature. But what if you replace the transition metal elements with **refractory elements** (W, Mo, Nb, Ta, V, Hf, Zr, Ti, Cr)?

**Refractory High-Entropy Alloys (RHEA):** HEAs based primarily on refractory elements with melting points > 1,650°C.

The concept: refractory elements have very high melting points → RHEA should maintain strength to higher temperatures than Ni superalloys (which are limited by Ni's melting point at 1,455°C).

**Target:** Structural alloys capable of service at 1,200–1,400°C metal temperature (vs. 1,050°C for best Ni SX).

### Early RHEA Systems

**MoNbTaW (Senkov, 2010):**
- BCC crystal structure
- Melting point: ~3,000°C (average of Mo 2,623°C + Nb 2,477°C + Ta 3,017°C + W 3,422°C)
- Yield strength at 1,600°C: ~400 MPa — remarkable! (Ni superalloys yield at ~100 MPa at this temperature)

**MoNbTaVW:**
- BCC
- Yield strength 1,600°C: ~477 MPa
- Density: 13.2 g/cm³ (much higher than Ni superalloys at ~8.7 g/cm³)

These results showed that RHEAs could indeed maintain high strength at temperatures far beyond Ni superalloys. But the density problem was severe.

---

## 7. RHEA Properties and Challenges

### High Strength at High Temperature — Real

Many RHEAs do show excellent yield strength retention at temperatures where Ni superalloys fail:

```
Strength comparison at 1,000°C, 1,200°C, 1,600°C (MPa yield strength):

Temperature:      1,000°C    1,200°C    1,600°C
CMSX-4 (Ni SX):    500        200         ~30 (near melting)
MoNbTaW (RHEA):    700        600         400
HfNbTaTiZr:        400        300         250
```

The RHEA advantage at 1,200°C+ is real and significant — ~3–5× higher specific yield strength than Ni superalloys.

### Challenges

**Challenge 1 — Brittleness at Room Temperature:**

Most BCC-based RHEAs have a **ductile-to-brittle transition temperature (DBTT)** above room temperature → they are brittle at room temperature → hard to process, handle, manufacture.

```
DBTT for key RHEAs:
MoNbTaW: ~500°C (brittle at room T — can't even bend without cracking)
WNbMoTa: ~500°C
NbMoTaW: >100°C (marginally workable)
HfNbTaTiZr: ~RT (ductile at room T — exceptional)
```

Tungsten and molybdenum rich alloys inherently have high DBTT → limits processing and component manufacturing.

**Challenge 2 — Oxidation:**

Refractory elements oxidize catastrophically:
- Mo: forms volatile MoO₃ above 700°C → rapid loss of material
- W: forms volatile WO₃ above 900°C → rapid oxidation
- Nb: forms Nb₂O₅, not very protective

An RHEA turbine blade operating at 1,200°C in combustion gas would oxidize and dissolve in hours without a sophisticated oxidation protection system. No existing coating system can protect RHEAs adequately in current form.

**Challenge 3 — Density:**

```
Density of common RHEAs:
MoNbTaW: 13.2 g/cm³ (vs. CMSX-4: 8.7 g/cm³ → 52% heavier)
WNbMoTa: 14.0 g/cm³
HfNbTaTiZr: 9.1 g/cm³ (acceptable — close to Ni SX)
```

High density → much higher centrifugal stress → would require massive disk redesign or shorter blades → net benefit from higher temperature capability may be marginal.

**Challenge 4 — Manufacturing:**

No existing investment casting process can cast RHEAs — refractory elements have melting points above ceramic crucible capabilities. Powder metallurgy or arc melting + machining are the only current options. Complex cooling-hole geometry (essential for turbine blades) is not producible.

---

## 8. HEA vs. Superalloys — A Realistic Comparison

| Property | Best Ni SX (CMSX-4) | Best FCC HEA (CoCrFeMnNi) | Best RHEA (HfNbTaTiZr) |
|----------|--------------------|--------------------------|-----------------------|
| Service T (max) | ~1,050°C | ~700°C | ~1,200°C (potential) |
| Yield strength (RT) | 1,100 MPa | 200 MPa | 800 MPa |
| Yield strength (1,000°C) | 500 MPa | ~80 MPa | 400 MPa |
| Oxidation resistance | Excellent (with coating) | Poor | Very poor |
| Castability | Excellent (SX) | Excellent | Very difficult |
| Density | 8.7 g/cm³ | 7.9 g/cm³ | 9.1 g/cm³ |
| Room T ductility | Low | Very high (>50%) | Moderate (~10%) |
| TRL (Technology Readiness) | 9 (in service) | 4–5 (lab) | 3–4 (concept) |

**Conclusion:** As of 2025, no HEA system is ready to replace Ni superalloys in turbine blades. RHEAs show promise for 2040+ applications, but require solutions to: oxidation protection, room-temperature ductility, and manufacturing scalability.

---

## 9. Precipitation-Hardened HEAs — Best of Both Worlds

The major weakness of single-phase HEAs is insufficient strength. Researchers are applying the γ/γ′ concept to HEAs:

**γ/γ′ HEA systems:**

Adding Al (and sometimes Ti or Nb) to CoCrFeNi systems:
```
(CoCrFeNi)₁₋ₓAlₓ → forms γ (FCC disordered) + γ′ (L1₂ ordered precipitates)
```

At Al = 10–15 at%: two-phase γ/γ′ microstructure forms — exactly like Ni superalloys!

**AlCoCrFeNi₂.₁ (prototype precipitate-hardened HEA):**
- γ/γ′ two-phase structure
- Yield strength: ~800 MPa (vs. ~200 for single-phase CoCrFeMnNi)
- Ductility: 15% elongation (acceptable)
- Density: ~7.5 g/cm³ (lower than Ni superalloys)

This is a genuinely interesting alloy: it has structural Ni superalloy strength but at lower density, and the "high entropy" in the γ matrix may provide additional solid-solution strengthening.

**Challenge:** Temperature capability limited by Ni-based γ′ stability → same fundamental limit as conventional Ni superalloys. No significant improvement in service temperature is achieved.

---

## 10. Composition Space — Exploring 10²⁵ Possible Alloys

One of the most profound insights from HEA research is the vastness of unexplored composition space:

**How many possible alloys exist?**

If we consider just 25 metallic elements from the periodic table, each at one of 10 possible compositions (0 to 45 at% in 5 at% steps), the number of possible 5-element equimolar alloys is:
```
C(25, 5) = 53,130 alloys (equimolar only)
With variable compositions: ~ 10²⁵ distinct compositions
```

Traditional metallurgy has explored perhaps 10,000–100,000 distinct compositions over 150 years. The vast majority of composition space has NEVER been explored.

**Computational screening + machine learning:**

Modern approach:
1. CALPHAD predicts phase stability for millions of compositions automatically
2. Machine learning models (trained on existing data) predict properties of unexplored compositions
3. High-throughput experiments (combinatorial deposition, arc melting libraries) rapidly screen many compositions experimentally
4. Optimal compositions identified computationally → targeted experimental validation

This "accelerated alloy discovery" paradigm is changing how metallurgists work — increasingly a data-science plus materials-science hybrid.

---

## 11. Where HEAs Are Going

**Near-term (2025–2030):**
- Cryogenic applications: CoCrFeMnNi-based alloys for LNG/LH₂ vessels → likely first commercial HEA products
- Wear-resistant coatings: hard HEA coatings (TiAlCrNbSi nitrides) on cutting tools → already in commercial products
- Medium-entropy alloys (MEA) as precursors to HEA in simpler systems

**Medium-term (2030–2040):**
- Precipitation-hardened HEAs for structural aerospace applications (at moderate temperatures)
- RHEAs for hypersonic applications (leading edges, re-entry vehicles) — where density penalty is acceptable
- First turbine test articles with RHEA airfoils (research engines)

**Long-term (2040+):**
- RHEAs with coating systems achieving turbine blade service conditions
- HEA-based composites (RHEA matrix + oxide reinforcement)
- Discovery of transformative new alloys from the vast unexplored composition space

---

## Summary

- **HEA definition**: ≥5 principal elements, each 5–35 at%, ΔS_mix > 1.5R = ~12.5 J/mol·K
- **High entropy stabilizes simple phases**: disordered FCC/BCC/HCP solid solutions instead of intermetallics
- **Four core effects**: high entropy → simple phases; severe lattice distortion → strong; sluggish diffusion → stable; cocktail → unpredictable bonuses
- **CoCrFeMnNi (Cantor alloy)**: FCC, exceptional cryogenic properties (K_Ic > 200 MPa√m at -196°C), moderate room-T strength → best for cryogenic applications
- **RHEAs (refractory HEA)**: W/Mo/Nb/Ta/Hf base → extremely high strength at 1,200°C+ → limited by oxidation, DBTT, density, manufacturing
- **Precipitation-hardened HEAs**: γ/γ′ structure in HEA → ~800 MPa RT strength, but no T-capability advantage over Ni SX
- **Vast unexplored space**: >10²⁵ possible compositions; computational HT screening + ML is the new discovery paradigm

---

## Exercises

1. Calculate ΔS_mix for: (a) equimolar CoCrFeMnNi, (b) equimolar MoNbTaW, (c) an alloy Ni-60Cr-20Mo-10Al-10Co (non-equimolar HEA — does it qualify as high-entropy?). Use ΔS_mix = -R × Σ(x_i × ln x_i).

2. CoCrFeMnNi at -196°C has yield strength 500 MPa and fracture toughness K_Ic = 210 MPa√m. (a) Using σ_y = K_Ic/√(π a_c), calculate the critical defect size a_c below which the alloy will yield before fracture. (b) Compare to CMSX-4 at RT (σ_y = 1,100 MPa, K_Ic = 40 MPa√m). Which alloy is more defect-tolerant?

3. MoNbTaW RHEA: density = 13.2 g/cm³, yield strength at 1,200°C = 600 MPa. CMSX-4: density = 8.7 g/cm³, yield strength at 1,200°C = 150 MPa. Calculate **specific yield strength** (yield strength / density) in units of MPa/(g/cm³) for each. Which is better, and by how much? How does this specific strength comparison change the picture from raw strength numbers?

4. Composition space: You are designing a 5-element RHEA for hypersonic leading edge applications (1,400°C service). Required properties: (a) melting point > 2,500°C, (b) density < 11 g/cm³, (c) at least one element providing oxidation resistance. Using elements W (3,422°C, 19.3 g/cm³), Mo (2,623°C, 10.2), Ta (3,017°C, 16.6), Nb (2,477°C, 8.6), Ti (1,668°C, 4.5), Cr (1,907°C, 7.2): propose an equimolar 5-element alloy and calculate its estimated melting point (rule of mixtures) and density (rule of mixtures). Does it meet all criteria?

5. Machine learning for alloy discovery: A ML model trained on 10,000 HEA compositions predicts that (CoCrFeMnNi)₉₀Al₁₀ has a γ′ volume fraction of 35% with a yield strength of 750 MPa at RT. (a) What experimental tests would you design to validate this prediction? List at least 4 measurements needed. (b) If ML screening can evaluate 100,000 compositions per day, and 1% of screened compositions are worth experimental testing, and each experimental test takes 3 days, how many researchers does it take to keep up with the ML screening rate? (c) Is this a resource bottleneck? What would you do?
