# Chapter 21: Precipitation Hardening — Aging to Maximum Strength

> **"In 1906, Alfred Wilm discovered something strange: an aluminum-copper alloy that had been quenched and left sitting overnight was much harder the next morning. He had stumbled upon precipitation hardening — one of the most powerful strengthening mechanisms in metals. Today it is the foundation of aerospace aluminum alloys, stainless steels, and nickel superalloys worth hundreds of billions of dollars."**

---

## Table of Contents

1. [The Principle — Why Precipitates Strengthen](#1-the-principle--why-precipitates-strengthen)
2. [Requirements for Precipitation Hardening](#2-requirements-for-precipitation-hardening)
3. [The Three-Step Process](#3-the-three-step-process)
4. [GP Zones — The First Precipitates](#4-gp-zones--the-first-precipitates)
5. [Metastable Precipitates — Peak Hardness](#5-metastable-precipitates--peak-hardness)
6. [Overaging and Stable Precipitates](#6-overaging-and-stable-precipitates)
7. [The Aging Curve and Peak Aging Condition](#7-the-aging-curve-and-peak-aging-condition)
8. [Precipitation in Major Alloy Systems](#8-precipitation-in-major-alloy-systems)
9. [Two-Step Aging and Thermomechanical Processing](#9-two-step-aging-and-thermomechanical-processing)
10. [Precipitation Hardening vs. Temperature](#10-precipitation-hardening-vs-temperature)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Principle — Why Precipitates Strengthen

Precipitation hardening works by creating a dense dispersion of fine, coherent or semi-coherent precipitates in the matrix. These precipitates:
- Create coherency strains (misfit between precipitate and matrix lattice) → stress field → impedes dislocation motion
- In ordered precipitates (like γ′ in Ni): dislocations create antiphase boundaries → APB energy → resists cutting
- At larger sizes: force dislocations to loop around (Orowan mechanism)

**Summary of mechanism evolution during aging:**
```
Early aging:         Peak aging:          Overaging:
(under-aged)         (optimal)            (over-aged)
                     
GP zones (few nm)    θ″ or γ′ (10–30nm)   θ (>100nm)
coherent             coherent/semicoher.   incoherent
                     
Low cutting          Max APB+Orowan        Pure Orowan
resistance           combined effect       (widely spaced)
                     
→ Soft              → Hardest              → Softer again
```

---

## 2. Requirements for Precipitation Hardening

Not all alloy systems can be precipitation hardened. Requirements:
1. **Decreasing solid solubility with decreasing temperature** (solvus must slope outward with T)
2. **A second phase that can precipitate** with a known crystal structure compatible with the matrix
3. **Sufficient coherency** between precipitate and matrix lattice (for early stages; incoherence develops on aging)
4. **Kinetics that are controllable** — the transformation must be slow enough at aging temperature to control, but fast enough to be practical

**Examples of suitable alloy systems:**
- Al-Cu (2xxx): Cu-rich θ phase series (GP zones → θ″ → θ′ → θ)
- Al-Mg-Si (6xxx): Mg₂Si (β″, β′, β)
- Al-Zn-Mg-Cu (7xxx): MgZn₂ (η″, η′, η)
- Ni-Al (superalloys): Ni₃Al (L1₂) = γ′
- Fe-Ni-Ti (maraging steel): Ni₃Ti, Ni₃Mo, Fe₂Mo
- PH stainless steels: Cu-rich, NiAl, or TiC precipitates

---

## 3. The Three-Step Process

### Step 1: Solution Treatment (Solutionizing)
```
Heat to T_solvus + 20-50°C → ALL second phase dissolves into solid solution
Hold until uniform composition throughout (1 hour / 25mm)
Atmosphere: air (Al alloys), vacuum (Ni superalloys, Ti alloys)
```

**Goal:** Single-phase supersaturated solid solution with all alloying elements dissolved.

**Critical requirement:** Achieve complete dissolution WITHOUT:
- Grain growth (too high T or too long time)
- Incipient melting (if T_solution > T_liquidus → catastrophic)
- Scale or decarburization (for steels, use controlled atmosphere)

### Step 2: Quench (Rapid Cooling)
```
Cool fast to room temperature (or intermediate quench temperature)
Media: water quench (Al, some steels), oil quench, air quench, polymer quench
```

**Goal:** Trap the alloying elements in solid solution at RT — create a SUPERSATURATED solid solution. If cooling is too slow, coarse equilibrium precipitates form → depletes matrix → less hardening potential → undesirable.

**Quench sensitivity:** Some alloys (7xxx Al, PH stainless) are very sensitive — even a 1°C/s slower cooling allows coarse precipitates and reduces peak hardness significantly.

### Step 3: Aging (Precipitation)
```
Heat to aging temperature (below solvus) and hold
Aging T: 100–200°C for Al alloys; 700–870°C for Ni superalloys
Aging time: minutes to hundreds of hours depending on T and alloy
```

**Goal:** Nucleate and grow fine precipitates — but NOT let them coarsen past the peak strength condition.

---

## 4. GP Zones — The First Precipitates

**Guinier-Preston (GP) zones:** The very first stage of precipitation — solute atom clusters that are FULLY COHERENT with the matrix and so small they have no distinct crystal structure yet.

**In Al-Cu alloys:**
- GP zones = 1–2 atomic layer thick platelets of Cu atoms on {100} planes of Al matrix
- Size: 1–5 nm diameter
- Coherent: Al lattice continues through them (but strained due to Cu being larger)
- Form at room temperature (natural aging) or at low aging temperatures (<100°C)

```
GP zone in Al-Cu (view along [001]):

● ● ● ● ● ● ● ← GP zone (Cu atoms, ● larger than Al, ○)
○ ○ ○ ○ ○ ○ ○ ← Al matrix
○ ○ ○ ○ ○ ○ ○
● ● ● ● ● ● ● ← another GP zone
○ ○ ○ ○ ○ ○ ○
```

**Properties:** Small size → high number density → multiple obstacles per dislocation path → good strengthening, but not as good as θ″.

**Natural aging (T6 vs T4):**
- T4 temper: quenched and naturally aged to stable condition at room temperature
- GP zones form over days-weeks at RT
- 2024-T4: good combination of strength and formability

---

## 5. Metastable Precipitates — Peak Hardness

As aging proceeds (higher T or longer time), GP zones evolve into metastable transition precipitates with distinct crystal structures:

### Al-Cu System:
```
GP zones → θ″ → θ′ → θ (CuAl₂)

θ″: tetragonal, fully coherent, ~5–10nm → MAXIMUM HARDNESS (T6 condition)
θ′: semi-coherent, larger (~50nm) → still hard but declining
θ:  incoherent, stable, large → overaged, soft
```

### Nickel Superalloy (Ni-Al):
```
Disordered γ′ clusters → ordered L1₂ γ′ → coarsened γ′

L1₂ γ′ at peak size (~10–30nm, 65% volume fraction) → MAXIMUM STRENGTH
(via combination of cutting + Orowan mechanisms)
```

The ordered L1₂ structure of γ′ is unique:
- Dislocations must move as pairs (superdislocations) to preserve order
- Single dislocation passing = antiphase boundary (APB) → energy cost
- APB energy γ_APB ≈ 180 mJ/m² for Ni₃Al → this is why γ/γ′ is so strong

---

## 6. Overaging and Stable Precipitates

If aging continues too long or at too high a temperature, precipitates coarsen (Ostwald ripening):
- Large particles grow at the expense of small ones
- Number density decreases, size increases, spacing λ increases
- Orowan stress τ = Gb/λ → decreases as λ increases → SOFTENS

**Coarsening kinetics (LSW theory):**
```
r³(t) - r³(0) = K × t

where r = mean precipitate radius
      K = rate constant = (2/9) × D × C_∞ × γ / (RT)
      D = diffusion coefficient
      γ = precipitate-matrix interface energy
```

**Intentional overaging (T7 temper):**
For some applications, a slightly overaged condition is PREFERRED:
- Better stress corrosion cracking resistance (7xxx Al alloys): T73 temper
- Better dimensional stability
- Lower residual stress
- Sacrifice 10–15% of peak strength for much better corrosion resistance

---

## 7. The Aging Curve and Peak Aging Condition

The aging curve maps hardness vs. time at a given aging temperature:

```
Hardness vs. aging time at fixed temperature (schematic):

Hardness
  │                    ╭──── peak age
  │                   ╱  ╲
  │                  ╱    ╲
  │                 ╱      ╲──────── overaged
  │                ╱
  │               ╱   θ″/γ′ growing
  │──────────────╱  GP zones → metastable precipitates
  │  As-quenched
  └──────────────────────────── log(aging time) →
     1hr    10hr   100hr  1000hr
```

**Effect of aging temperature:**
Higher T → faster kinetics → peak reached sooner → lower peak hardness (coarser precipitates at same volume fraction → wider spacing → lower Orowan stress):

```
       Low T aging:    long time, fine precipitates, HIGHEST hardness
       Medium T aging: medium time, optimal balance
       High T aging:   short time, coarser precipitates, lower hardness
```

For 7075 Al:
- T6 (peak): 120°C for 24 hours → σ_y = 503 MPa, σ_UTS = 572 MPa
- T73 (overaged): 120°C + 175°C for 8 hours → σ_y = 435 MPa, but better SCC resistance

---

## 8. Precipitation in Major Alloy Systems

### Al-Cu (2xxx series) — 2024, 2219
| Stage | Precipitate | Structure | Size | Condition |
|-------|-------------|---------|------|-----------|
| 1 | GP zones | {100} monolayers | 1–5 nm | T4 |
| 2 | θ″ | Tetragonal, coherent | 5–20 nm | T6 (peak) |
| 3 | θ′ | Semi-coherent | 20–100 nm | T8 |
| 4 | θ (CuAl₂) | Incoherent | >100 nm | Overaged |

### Al-Mg-Si (6xxx series) — 6061, 6082
GP zones → β″ (Mg₅Si₆, monoclinic, coherent) → β′ → β (Mg₂Si)
6061-T6: σ_y = 276 MPa (peak β″); 6061-T4: 145 MPa (GP zones only)

### Al-Zn-Mg-Cu (7xxx series) — 7075, 7068
GP zones → η″ → η′ → η (MgZn₂); Zn/Mg ratio critical
7075-T6: σ_y = 503 MPa (highest strength commercial Al alloy)

### Maraging Steels (Fe-18Ni + Co, Mo, Ti)
- Solution treat at 820°C → air cool (18%Ni suppresses martensite to RT = SOFT martensite)
- Age at 480°C for 3 hours
- Ni₃Mo + Ni₃Ti precipitate in the martensite matrix → HUGE strengthening
- Peak σ_y: Grade 250 = 1,720 MPa; Grade 300 = 2,070 MPa!
- Used for: rocket motor cases, landing gear, tool inserts

### PH Stainless Steels (17-4 PH, 15-5 PH)
- Martensitic PH stainless: solution treat → quench → age at 470–620°C
- Cu-rich ε-Cu precipitates in martensite matrix
- Combines corrosion resistance + high strength (σ_y = 700–1,200 MPa)
- Used for: surgical instruments, aerospace fasteners, pump shafts

### Nickel Superalloys (γ/γ′)
- Covered in detail in Ch 44 — the most important precipitation hardening system for high-temperature service

---

## 9. Two-Step Aging and Thermomechanical Processing

**Two-step (retrogression and re-aging, RRA):**
Used for 7xxx Al alloys to improve SCC resistance while maintaining near-T6 strength:
1. Age to T6 condition
2. "Retrogress" by brief heating to 200°C (dissolves GP zones near grain boundaries that cause SCC)
3. Re-age at 120°C → T6 strength restored, but better SCC resistance

**Thermomechanical processing (TMCP):**
Cold or warm working BEFORE or DURING aging → dislocations act as additional nucleation sites for precipitates → finer, more uniform precipitate distribution → higher peak hardness:

```
T8 temper (2xxx Al): Solution treat → cold work 6% → age
Result: higher peak hardness than T6 (no cold work)
```

**Used in turbine disk alloys:** Forging in γ/γ′ two-phase field → deformed γ′ particles → on subsequent aging, highly heterogeneous nucleation → fine γ′ → better creep resistance.

---

## 10. Precipitation Hardening vs. Temperature

**The critical weakness:** Precipitates coarsen at elevated temperatures. Service temperature must stay well below the solvus to prevent overaging in service.

```
Maximum service temperatures (rough guidelines):
Al alloys (7075-T6):    120–150°C sustained
Al alloys (2219):       180°C (better than 7075)
Cu-Be alloys:           250°C
Maraging steel 300:     350°C
17-4 PH stainless:      400°C
IN718 (γ″):             650°C (γ″ dissolves above 650°C)
Nickel superalloys (γ′): 1,050°C (γ′ thermally stable to ~1,300°C, but coarsens above 1,000°C)
```

The γ′ system in nickel superalloys is extraordinary — it maintains precipitation hardening up to 1,050°C because:
1. γ′ has anomalous yield increase up to 800°C (unique)
2. Low γ/γ′ misfit in modern alloys → slow coarsening
3. Re additions (CMSX-4) retard coarsening further

---

## Summary

| Step | Purpose | Key Parameter |
|------|---------|--------------|
| Solution treat | Dissolve all precipitates → single phase | T > T_solvus, avoid grain growth |
| Quench | Trap supersaturation | Fast enough to avoid coarse precipitate on cooling |
| Age | Nucleate fine precipitates | T and time to reach peak |
| GP zones | Early stage, coherent clusters | Form at RT or low T; T4 condition |
| Metastable precipitate | Peak hardness | θ″, β″, γ′ — coherent/semi-coherent |
| Overaging | Progressive coarsening and softening | Orowan spacing increases |
| Peak age (T6) | Maximum strength | Balance cutting vs. Orowan mechanisms |
| Overaged (T73) | Better SCC resistance at lower strength | Deliberately coarsened |

---

## Exercises

1. An Al-4%Cu alloy is solution treated at 550°C and quenched. Natural aging occurs at 25°C. (a) Why does Cu not precipitate immediately on quenching? (b) Using GP zone coherency argument, explain why GP zones nucleate on {100} planes preferentially. (c) After 4 days at 25°C, the material is artificially aged at 130°C. Draw the expected hardness vs. time curve. Mark the GP zone dissolving stage (brief softening before rebound), θ″ peak, and θ′ overaging region.

2. 7075-T6 aluminum has σ_y = 503 MPa; 7075-T73 has σ_y = 435 MPa (14% lower). The T73 alloy is immune to stress corrosion cracking in aircraft wing skins exposed to salt spray. Calculate: (a) weight increase required to achieve same load-carrying capacity in T73 vs T6 (area must increase by σ_y(T6)/σ_y(T73) - 1), (b) if wing skin is 2 mm × 500 mm × 3,000 mm (one panel), density 2.81 g/cm³: mass of T6 panel vs T73 panel, (c) for a 100-panel aircraft, total extra mass with T73. Does the corrosion resistance justify the weight penalty in your view?

3. A maraging steel 300 component is aged at 480°C. Using the LSW coarsening equation r³(t) = r³(0) + K×t with K = 3×10⁻²⁷ m³/s and initial precipitate radius r₀ = 2 nm: (a) Calculate r at t = 1h, 10h, 100h, 1000h. (b) If Orowan stress τ ∝ 1/r, how much does the Orowan contribution change from 1h to 1000h? (c) Why is maraging steel NOT suitable for 700°C service even though it has excellent room-temperature properties?

4. Nickel superalloy CMSX-4 requires a two-step aging after solution treatment: (a) first age 1,140°C × 6h, (b) then age 870°C × 20h. Explain the purpose of each aging step: (i) what precipitate sizes/morphologies each step produces (primary coarse γ′ vs. fine secondary γ′), (ii) why a bimodal γ′ distribution (two different sizes) gives better creep resistance than a monomodal distribution, (iii) what happens if the 1,140°C age is too long (hint: γ′ coarsening + TCP phase risk).

5. A designer wants to use 2024-T3 aluminum (cold-worked after quench) for a structural bracket that will see temperatures up to 150°C during service. (a) What aging stage is the T3 condition? (b) What microstructural change occurs during service at 150°C? (c) Estimate the change in yield strength after 1,000 hours at 150°C (assume aging curve peak at 1 hour, then overaging follows σ_y ∝ t^(-0.15) from peak). (d) Should this material be used in this application, or should 7075-T73 be used instead? Compare the trade-offs.
