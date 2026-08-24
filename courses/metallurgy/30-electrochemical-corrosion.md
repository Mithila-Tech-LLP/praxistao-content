# Chapter 30: Electrochemical Corrosion — Forms and Protection

> **"The rusting of iron is not a single reaction — it is an electrochemical cell. The bridge that crosses a harbor corrodes faster at the waterline than in air, faster in winter than summer, and can be protected by simply attaching a zinc block to its legs. Understanding corrosion means understanding electricity, chemistry, and materials science simultaneously. But understanding it also means billion-dollar savings and countless lives protected."**

---

## Table of Contents

1. [The Galvanic Series](#1-the-galvanic-series)
2. [Galvanic Corrosion](#2-galvanic-corrosion)
3. [Pitting Corrosion](#3-pitting-corrosion)
4. [Crevice Corrosion](#4-crevice-corrosion)
5. [Intergranular Corrosion](#5-intergranular-corrosion)
6. [Stress Corrosion Cracking (SCC)](#6-stress-corrosion-cracking-scc)
7. [Corrosion Fatigue](#7-corrosion-fatigue)
8. [Erosion-Corrosion and Cavitation](#8-erosion-corrosion-and-cavitation)
9. [Cathodic Protection](#9-cathodic-protection)
10. [Anodic Protection and Inhibitors](#10-anodic-protection-and-inhibitors)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Galvanic Series

The **galvanic series** ranks metals by their tendency to corrode in seawater (or another specific electrolyte). More negative potential = more anodic = corrodes preferentially:

```
ACTIVE (anodic, corrode)
─────────────────────────────────────
Magnesium (−1.73 V vs SHE)
Zinc      (−0.76 V)
Aluminum  (−0.66 V) [protected by passive film; actual behavior closer to −0.4V]
Cadmium   (−0.40 V)
Mild steel, iron (−0.44 V)
Stainless steel 304, 316 (active)
Lead      (−0.13 V)
Tin       (−0.14 V)
Nickel (active)  
Brass (70/30)
Bronze    
Cupronickel  
Copper (+0.34 V)
Stainless steel 316 (passive) (+0.2 V typical in seawater)
Silver (+0.80 V)
Titanium    
Platinum (+1.19 V)
Gold (+1.50 V)
─────────────────────────────────────
NOBLE (cathodic, protected)
```

**Key rule:** When two dissimilar metals are electrically connected in an electrolyte, the more active (higher on list) becomes the anode and corrodes; the more noble (lower on list) is the cathode and is protected.

---

## 2. Galvanic Corrosion

When dissimilar metals are in electrical contact in an electrolyte, the less noble metal corrodes at an accelerated rate.

**Classic example — Steel bolt in copper fittings:**
```
Cu fitting (cathode, noble) ── electrolyte ── Steel bolt (anode, active)
                         └── electrons → ┘
Steel bolt corrodes → failure; copper fitting protected
```

**Factors controlling severity:**
1. **Potential difference:** Larger ΔE → more severe (Mg near steel much worse than Cu near steel)
2. **Area ratio:** Large cathode / small anode = very BAD (concentrated attack on small anode area)
   Small cathode / large anode = mild (diffuse attack on large area)
3. **Electrolyte conductivity:** High conductivity (seawater) → worse than low conductivity (fresh water, humid air)
4. **Distance:** Corrosion concentrated near contact zone (worst at contact, diminishes with distance)

**Galvanic series area ratio example:**
A steel bolt in a large copper plate in seawater → enormous cathode area driving rapid corrosion of small steel bolt → bolt fails in months vs. years in non-galvanic situation.

**Design rules:**
- Avoid connecting dissimilar metals (especially > 0.25 V difference in electrolyte)
- If necessary: use larger anodic metal (or same metal)
- Insulate the junction (plastic washers, paint)
- Apply protective coating to the CATHODE (not the anode — breaks down and causes concentrated attack)
- Use sacrificial anode (intentionally make something more active corrode instead)

**CFRP + aluminum:** Carbon fiber is noble (+0.1 to +0.3 V); aluminum is active (−0.66 V). The Al corrodes at joints unless insulated. Boeing 787 uses Ti fasteners (similar potential to CFRP) and sealant to insulate the interface.

---

## 3. Pitting Corrosion

**Pitting** is localized breakdown of the passive film followed by rapid localized attack. Characteristic of passive metals (stainless steel, Al) in chloride environments.

**Mechanism:**
```
Passive film (Cr₂O₃ on stainless)
     ↓  Cl⁻ adsorbs on film defect (inclusion, scratch)
Film penetrated at specific site
     ↓
Metal dissolves locally: Fe → Fe²⁺ + 2e⁻ (inside pit)
O₂ reduction: O₂ + 2H₂O + 4e⁻ → 4OH⁻ (outside pit, on passive surface)
     ↓
Inside pit: acidic (H⁺ generated), Cl⁻ concentrated (to balance Fe²⁺)
→ self-catalyzing (autocatalytic) growth
     ↓
Outside pit: passive (stays protected)
→ small pit area = very HIGH local corrosion rate
```

**Pitting potential E_pit:** Above this potential, pitting initiates. Below it, passive.

**Mo effect on pitting:** Mo (in 316L, 317L, Hastelloy) raises E_pit → harder to pit → better pitting resistance.

**Pitting corrosion index (PREN):**
```
PREN = %Cr + 3.3×%Mo + 16×%N

Higher PREN → better pitting resistance in chlorides:
316L:      PREN = 24
317LMN:    PREN = 32
25Cr duplex 2507: PREN = 42
Super duplex (254 SMO): PREN = 43
Ti Grade 2: immune (PbO₂ passive film)
```

---

## 4. Crevice Corrosion

Crevice corrosion attacks passive metals in tight gaps (under bolt heads, gaskets, biofilm deposits) where stagnant electrolyte accumulates:

**Mechanism:**
```
    ─────────────────────
    │  Metal    │ Crevice (narrow gap, stagnant)
    │──────────────────
    │ metal │ gasket │ metal
    
1. O₂ depletes inside crevice (consumed, slow diffusion in)
2. Metal ions accumulate (Fe²⁺, Cr³⁺)
3. To maintain charge balance: Cl⁻ migrates IN
4. Hydrolysis: FeCl₂ + H₂O → Fe(OH)₂ + 2HCl → ACIDIC inside
5. Passive film dissolves in acid → active corrosion begins
6. Self-accelerating: more dissolution → more acid → more Cl⁻ migration
```

**Results:** Deep pitting inside crevice while adjacent open surface remains passive. Often worse than plain pitting.

**Prevention:**
- Eliminate crevices by design (full-penetration welds, no lap joints)
- Use higher alloy (higher PREN)
- Change to Ti or HDPE in extreme service
- Drain and dry periodically
- Use PTFE or elastomer gaskets that don't form metallic crevice

---

## 5. Intergranular Corrosion

Attack at grain boundaries (which may be depleted in protective element or enriched in susceptible phases).

**Sensitization of 304 stainless steel:**
- Heating 450–850°C → Cr₂₃C₆ precipitates at grain boundaries
- Adjacent metal depleted in Cr to < 12% → loses passivity
- Grain boundaries corrode → "knife-line attack" → intergranular failure in corrosive environment

```
Grain boundary cross-section:

Far from boundary:    Near boundary (sensitized):
   Cr = 18%             Cr = 6%
                         ← Cr depleted zone → ← Cr₂₃C₆ precipitates
```

**Test: ASTM A262 (Strauss test):** Immerse in CuSO₄+H₂SO₄ → bent sample → if cracks appear = sensitized.

**Prevention:** L-grade (304L, 316L), Ti/Nb-stabilized grades, or solution anneal + fast quench after welding.

**Weld decay:** HAZ of 304 weld at 450–800°C → sensitized zone adjacent to weld → corrosion zone slightly away from weld bead.

---

## 6. Stress Corrosion Cracking (SCC)

**SCC:** Cracking caused by the COMBINATION of:
1. Tensile stress (can be residual or applied)
2. Specific corrosive environment
3. Susceptible material

Neither stress alone nor environment alone causes cracking — requires both.

**Classic SCC systems:**

| Material | Environment | Mechanism |
|----------|-------------|-----------|
| Austenitic SS | Chlorides, >50°C | Cl-induced passivation breakdown + stress |
| Brass (70/30) | Ammonia/amines | NH₃ + residual stress (season cracking) |
| Al 7xxx T6 | Salt + humidity | GBcorrosion + stress |
| High-strength steel | H₂SO₄, H₂S, cathodic | Hydrogen embrittlement (H-assisted SCC) |
| Ti alloys | Methanol, HNO₃ | Specific environment sensitivity |
| Magnesium | NaCl | Na, Cl attack |

**Mechanics:** SCC crack grows at K_I < K_Ic — specifically at K_ISCC (threshold for SCC):
```
K_ISCC << K_Ic for susceptible systems
(e.g., 7075-T6 in NaCl: K_Ic = 29 MPa√m but K_ISCC ≈ 3–6 MPa√m!)
```

**Prevention:**
- Remove one of the three requirements (material, stress, environment)
- Change alloy: T73 temper for 7075 (overaged → less susceptible); duplex SS in chlorides
- Reduce stress: stress relief heat treatment, shot peening (compressive surface stress)
- Environment: dehumidification, inhibitors, pH control
- Coatings/isolation

---

## 7. Corrosion Fatigue

**Corrosion + cyclic stress → MUCH faster crack initiation and growth than either alone.**

**Effects:**
- Fatigue limit eliminated (no safe stress below endurance limit in corrosive environment)
- Paris Law exponent m increases (crack grows faster per ΔK cycle)
- Crack initiates from corrosion pits (which are stress concentrators K_t ≈ 2–4)

```
S-N curves:
                    ─── inert environment (has endurance limit)
Stress
amplitude  ─────────────────────────────────────────
           ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─  (endurance limit)
           ─────────────────────────────────────────── corrosive environment
                                                       (no endurance limit)
           ─────────────────────────────────────────── log N
```

**Important for:**
- Ship propeller shafts (seawater + rotating bending stress)
- Offshore platform welds (seawater + wave loading)
- Aircraft fasteners (aluminum + moisture)

**Prevention:** Same as SCC + fatigue control (shot peening, low stress concentrations, surface coatings).

---

## 8. Erosion-Corrosion and Cavitation

**Erosion-corrosion:** High-velocity fluid removes the protective passive film → fresh active metal exposed → rapid corrosion → more film removal → accelerating attack.

**Critical parameters:**
- Velocity threshold: below some V, passive film intact; above, erosion removes it
- Turbulent flow: more damaging than laminar
- Particulates: sand/silt in water dramatically worsen erosion-corrosion
- Geometry: bends, elbows, valves are worst-hit (high turbulence, impingement)

**Cavitation corrosion:** Rapid formation and collapse of vapor bubbles in fast-flowing liquid → collapse produces local pressure spikes of 500+ MPa → fatigue-like damage + film removal → deep pitting.

Common in: marine propellers, hydraulic turbines, ship pump impellers, engine cylinder liners.

**Prevention:** Use harder, more noble alloys (Cu-Al-Ni bronze for propellers, 316L for pump impellers), reduce flow velocity, streamline geometry.

---

## 9. Cathodic Protection

**The most important corrosion prevention method for buried/submerged steel.**

**Principle:** Drive the steel's potential below E_corr (into the immunity region on the Pourbaix diagram) by supplying electrons. Two methods:

### Sacrificial Anode Cathodic Protection (SACP)

Connect a more active metal (Zn, Mg, Al alloy anode) to the steel:
```
Zn (anode) ──────── electrolyte ──────── Steel (cathode, protected)
                                         
Zn → Zn²⁺ + 2e⁻  → electrons → Steel: O₂ + H₂O + 2e⁻ → 2OH⁻
(Zn corrodes sacrificially)       (steel protected — cathode)
```

**Anode materials:**
- Zinc alloy (97%Zn-0.1%Cd-0.003Al): seawater and brackish water
- Aluminum alloy: offshore platforms (deepwater, current generation = long life)
- Magnesium: fresh water (soil, pipelines)

**Application:** Ship hulls (zinc blocks), offshore platforms, buried pipelines (Mg anodes), water tanks.

**Design:** Anode must supply enough current to shift entire structure below protection potential:
```
I_protection = A_structure × i_protection (A/m²)
i_protection ≈ 0.05–0.15 A/m² for steel in seawater
Anode output = (M_anode × F) / (e × M_element) (Faraday's law)
Anode life = total charge output / required current
```

### Impressed Current Cathodic Protection (ICCP)

External DC power source connected to an inert anode (mixed metal oxide, graphite):
- Power supply drives current through electrolyte → steel is always cathodic
- Allows continuous monitoring and adjustment
- Better suited for large structures (pipelines thousands of km long)
- More expensive per installation, cheaper for large structures

**Pipeline protection:** Trans-Alaska Pipeline uses ICCP for ~1,300 km → monitoring every few km.

**Protection criteria (NACE):**
- Steel in soil: E ≤ −0.85 V vs Cu/CuSO₄
- Steel in seawater: E ≤ −0.80 V vs Ag/AgCl

---

## 10. Anodic Protection and Inhibitors

### Anodic Protection

Applied in specific cases where driving the metal MORE positive (anodic direction) moves it into the passive region:
- Used only for metals that passivate (stainless steel in H₂SO₄)
- Maintain potential in passive window: if too negative → corrosion; too positive → transpassive
- Chemical storage tanks for H₂SO₄

### Corrosion Inhibitors

Chemicals added to the electrolyte that reduce corrosion rate:

| Type | Mechanism | Examples |
|------|-----------|---------|
| Anodic inhibitors | Block anodic sites | Chromates, nitrites, molybdates |
| Cathodic inhibitors | Block cathode reaction (O₂ reduction or H evolution) | Amines, phosphates |
| Mixed inhibitors | Both | Benzotriazole (Cu protection) |
| Film-forming | Adsorb on metal surface | Organic amines in steel |

**Engine coolant inhibitors:** Ethylene glycol + package of inhibitors (molybdates, benzotriazole, carboxylates) → protects Al, Cu, solder, and steel in the cooling system simultaneously.

---

## Summary

| Corrosion Form | Driving Factor | Prevention |
|---------------|---------------|-----------|
| Galvanic | Potential difference between metals | Insulate, sacrifice, same-metal use |
| Pitting | Cl⁻ breaks passive film locally | Higher PREN alloy, remove Cl⁻ |
| Crevice | Stagnant, O₂-depleted zone | Design out crevices, higher alloy |
| Intergranular | Cr depletion by sensitization | L-grade, stabilized, anneal |
| SCC | Stress + corrosive environment | Reduce stress, change alloy/environment |
| Erosion-corrosion | High velocity removes passive film | Harder alloy, reduce velocity |
| Prevention (cathodic) | Drive potential below E_corr | Sacrificial anodes or ICCP |

---

## Exercises

1. A steel pipe (−0.44 V vs SHE) is connected to a copper valve (+0.34 V). The pipe area = 10 m², valve area = 0.01 m². (a) Which is the anode? (b) The potential difference = 0.78 V → estimate galvanic current using empirical rule: i_gal ≈ 1 mA/cm² per 0.1V ΔE (rough estimate). (c) Calculate total current and current density at the pipe (anode), then estimate corrosion rate (use Faraday: i × M / (n × F)). (d) How does replacing the copper valve with Ti change the corrosion situation?

2. An austenitic 304 stainless storage tank was welded at a fabrication shop, then put into service containing 3.5% NaCl at 80°C. After 6 months, intergranular corrosion appeared in the weld HAZ. (a) Describe the metallurgical cause. (b) The sensitized zone shows Cr < 11% near grain boundaries — explain why this is below the minimum for passivation. (c) If you could redo the specification: would you use 304L, 316L, or 321? Justify each choice. (d) Can the existing tank be repaired? Describe the thermal treatment needed and its limitations.

3. A 316L pipeline carries seawater at 50°C. After 3 years, deep pitting (up to 6mm deep) is found at flange connections. (a) Identify whether this is pitting or crevice corrosion, and explain your reasoning. (b) The specification says "316L is immune to pitting in seawater at 25°C." Why has failure occurred at a flange/gasket at 50°C? (c) Suggest three design changes to prevent recurrence without changing the base material. (d) Superaustenitic alloy 254 SMO (PREN = 43) is proposed as an upgrade. Calculate its PREN from %Cr = 20, %Mo = 6.1, %N = 0.18.

4. An offshore oil platform has 10,000 m² of steel hull below the waterline. Design a sacrificial anode cathodic protection system: (a) calculate total current needed at 0.1 A/m², (b) aluminum alloy anode capacity = 2,500 A·h/kg; if anodes weigh 50 kg each, how many anodes are needed for a 5-year life? (c) if zinc anodes (780 A·h/kg) were used instead, how many 50kg zinc anodes? (d) Compare total anode mass for Al vs. Zn systems.

5. SCC testing of 7075-T6 aluminum in 3.5% NaCl reveals that cracks grow when K > K_ISCC = 4 MPa√m. In service, the alloy is under 200 MPa stress. (a) Calculate the critical crack size a_c,SCC using K_ISCC = 4 = F×σ×√(πa) with F=1.0. (b) Compare to K_Ic = 29 MPa√m: what is a_c,fracture? (c) Without corrosion protection, the part will fail at a_c,SCC which is much smaller than NDE detection limit (0.5mm). This means the part fails without warning — discuss implications for aircraft inspection interval. (d) How would switching to 7075-T73 (K_ISCC ≈ 18 MPa√m) change a_c,SCC?
