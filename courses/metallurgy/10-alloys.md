# Chapter 10: Alloys — Mixtures That Outperform Pure Metals

> **"Pure iron has a yield strength of 130 MPa. Add 0.8% carbon, heat-treat, and you get 1,000 MPa. Add nickel, chromium, and cobalt in the right proportions and you approach 2,000 MPa at room temperature — and still maintain useful strength at 1,100°C. Alloying is the metallurgist's most powerful tool: combining elements to create properties that no single element can provide."**

---

## Table of Contents

1. [What Is an Alloy?](#1-what-is-an-alloy)
2. [Solid Solutions — Substitutional and Interstitial](#2-solid-solutions--substitutional-and-interstitial)
3. [Hume-Rothery Rules for Solid Solubility](#3-hume-rothery-rules-for-solid-solubility)
4. [Intermediate Phases and Intermetallic Compounds](#4-intermediate-phases-and-intermetallic-compounds)
5. [Ordered vs. Disordered Structures](#5-ordered-vs-disordered-structures)
6. [Solubility Limits and the Solvus Line](#6-solubility-limits-and-the-solvus-line)
7. [How Alloying Changes Properties](#7-how-alloying-changes-properties)
8. [Alloy Designation Systems](#8-alloy-designation-systems)
9. [Multi-Component Alloys — The High-Entropy Concept](#9-multi-component-alloys--the-high-entropy-concept)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is an Alloy?

An **alloy** is a metallic material composed of two or more elements, where at least one is a metal. Alloys are designed to have specific combinations of properties that pure metals cannot provide:

| Limitation of pure metal | Alloying solution |
|--------------------------|-------------------|
| Too soft (pure Al) | Add Cu, Mg, Zn, Si → precipitation-hardened alloys |
| Too reactive (pure Fe corrodes) | Add Cr → stainless steel (passive oxide) |
| Melts too low for gas turbines (pure Ni: 1,455°C) | Add Al, Ti, Ta, Re → raise effective operating T |
| Too brittle (pure tungsten) | Add Ni, Fe → tungsten heavy alloys |
| Poor castability (pure Al) | Add Si → hypereutectic Al-Si (good fluidity) |

**Components:** The individual chemical species (elements, or sometimes stable compounds) in the alloy.
**Phases:** Distinct homogeneous regions with different structure or composition (see Ch 06).

---

## 2. Solid Solutions — Substitutional and Interstitial

The simplest alloy structure is a **solid solution** — a single phase where solute atoms fit into the solvent's crystal structure.

### Substitutional Solid Solution

Solute atoms **replace** solvent atoms at lattice sites. Cu-Ni is a classic example: Ni and Cu have nearly identical atomic radii and crystal structures → complete solid solubility across all compositions.

```
Substitutional solid solution (Ni in Cu):

FCC lattice with random substitution:
 Cu Cu Ni Cu Cu Cu
 Cu Ni Cu Cu Ni Cu  ← Ni (●) randomly replacing Cu (○)
 Cu Cu Cu Cu Cu Ni
 
Properties change continuously with composition.
```

### Interstitial Solid Solution

Small solute atoms **fit into gaps (interstices)** between host atoms. Carbon in iron is the most important example:

```
Interstitial solution (C in Fe-BCC):

         ○  ─  ○
         │     │
○  ─  ●  ─  ○
│     │     │     ← ● = interstitial C in octahedral hole
○  ─  ○  ─  ○

BCC octahedral site radius = 0.036a ≈ 0.064 nm
C atom radius = 0.077 nm → DISTORTS the lattice → strong hardening!
Max solubility in α-Fe: 0.022% C at 727°C
Max solubility in γ-Fe: 2.14% C at 1,147°C (larger octahedral sites in FCC)
```

**Interstitial hardening is VERY effective:** Even tiny amounts of C (0.1–0.8%) dramatically strengthen steel because every C atom distorts the local lattice and interacts strongly with dislocations.

---

## 3. Hume-Rothery Rules for Solid Solubility

First stated by William Hume-Rothery (1934), these rules predict whether two metals will form extensive solid solutions:

### Rule 1: Atomic Size Factor
The atomic radii must differ by less than ±15%:
```
|r_solute - r_solvent| / r_solvent < 0.15
```
If size difference > 15% → large lattice distortion → limited solubility.

**Examples:**
- Ni (r=1.25Å) in Cu (r=1.28Å): 2.4% → extensive solubility ✓
- Cd (r=1.51Å) in Cu (r=1.28Å): 18% → limited solubility ✗

### Rule 2: Electronegativity Difference
Similar electronegativities → metallic bonding → solid solution.
Large difference → ionic or covalent bonds → intermetallic compound forms instead.

### Rule 3: Crystal Structure
Same crystal structure (both FCC, or both BCC) → more likely to form continuous solid solution.
Different structures → solubility limited to terminal solid solutions.

### Rule 4: Valence Difference
Metals with higher valence dissolve better in metals with lower valence than vice versa.
Large valence difference → strong tendency to form intermetallics.

**When ALL FOUR rules are satisfied:** Complete solid solubility possible (e.g., Cu-Ni, Au-Ag).
**When rules are violated:** Limited solid solubility, eutectic systems, or intermetallic phases form.

---

## 4. Intermediate Phases and Intermetallic Compounds

When elements don't follow the Hume-Rothery rules, they can form **intermediate phases** — new crystal structures that don't correspond to either element's structure:

### Electron Compounds (Hume-Rothery Phases)
Form at specific electron-to-atom ratios. The crystal structure is determined by the electron concentration:
- e/a = 3/2: β-brass (CuZn), β-NiAl → BCC structure
- e/a = 21/13: γ-brass (Cu₅Zn₈) → complex cubic
- e/a = 7/4: ε-phase (CuZn₃) → HCP

### Interstitial Compounds
Formed by small atoms (C, N, B, H) with transition metals:
- Fe₃C (cementite): important in steel (Ch 07)
- TiC, WC: extremely hard carbides used in cutting tools
- TiN: gold-colored wear coating

### True Intermetallic Compounds
Fixed compositions with strong ionic or covalent character:
- Ni₃Al (γ′) → the key strengthening phase in nickel superalloys (Ch 44)
- NiAl (β) → ordered B2 structure, used in bond coats (Ch 46)
- Ni₃Nb (γ″) → D0₂₂ structure in IN718
- Fe₃Al, FeAl → iron aluminides (potential structural alloys)
- MoSi₂ → silicide (extreme temperature applications)

**Intermetallic properties — trade-offs:**
| Property | Intermetallics | Conventional Alloys |
|----------|----------------|---------------------|
| Melting point | Very high (often > 1,500°C) | Moderate |
| Hardness/strength | Extremely high | Moderate |
| Oxidation resistance | Often excellent | Variable |
| Ductility (RT) | Often brittle (ordered structure resists dislocation motion) | Usually ductile |
| Fracture toughness | Low (5–20 MPa√m) | High (30–100+ MPa√m) |

---

## 5. Ordered vs. Disordered Structures

In a **disordered** solid solution, solute and solvent atoms are randomly distributed on the lattice sites. In an **ordered** structure, atoms arrange in a regular, periodic pattern.

```
Disordered vs. Ordered (schematic for A-B 50-50 alloy):

Disordered (high T):        Ordered (low T, = intermetallic):
  A B A B B A               A B A B A B
  B A A B A B               B A B A B A
  A B B A B A               A B A B A B
  
Atoms on any site randomly  Atoms on specific sublattices only
```

**Order-disorder transformation:** Many alloys are ordered at low temperature but disorder as temperature rises above an ordering temperature T_c. Below T_c, the ordered structure is thermodynamically stable.

**Example — Ni₃Al (γ′):**
- Below ~1,350°C: ordered L1₂ structure (Ni on face-center sites, Al on corner sites)
- This ordering makes dislocations move in pairs (superdislocations) → anomalous yield effect (strength INCREASES with temperature up to ~800°C, Ch 44)
- In a nickel superalloy at service temperature (900–1,000°C), γ′ remains ordered → this is why it works!

---

## 6. Solubility Limits and the Solvus Line

Solid solutions have a **solubility limit** at each temperature — the maximum amount of solute that can dissolve without forming a second phase. This limit is shown as the **solvus line** on a phase diagram.

```
Schematic phase diagram (solubility curve):

T (°C)
  │
  │    Liquid
  │───────────────────
  │   α solid solution │ α + L (two phase)
  │─────────────────────────────
  │   α solid solution │ α + β
  │        (single phase)    ← solvus line
  │─────────────────────────────────
  └──────────── Composition (% B) ──────→
```

**Precipitation hardening** exploits the decreasing solubility with temperature:
1. At high T: alloy is in single-phase region (all solute dissolved)
2. Quench to low T: alloy is supersaturated (more solute than equilibrium allows)
3. Age at intermediate T: fine precipitates nucleate and grow → EXTREMELY effective strengthening

This is covered in detail in Ch 21 (Precipitation Hardening) and is the basis for all high-strength aluminum alloys, many steels, and nickel superalloys.

---

## 7. How Alloying Changes Properties

### Solid Solution Hardening

Even fully dissolved solute atoms harden the matrix:
```
Δσ_ss = A × c^n   (n = 1/2 to 2/3 typically)

where c = solute concentration, A = constant depending on mismatch
```

**Mechanisms:**
- **Size mismatch:** Solute strains the lattice → interacts with dislocation strain fields → resists dislocation motion. Mo, W in nickel (large atoms) → strong solid solution strengthening.
- **Modulus mismatch:** Solute changes local elastic modulus → affects dislocation line energy
- **Electrical (Suzuki) effect:** Solute segregates to stacking faults → locks partial dislocations

**Size mismatch effectiveness:**
| Element in Ni | Atomic size mismatch | Strengthening (MPa per at%) |
|---------------|---------------------|----------------------------|
| Mo            | +12%                | ~40                         |
| W             | +10%                | ~35                         |
| Re            | +9%                 | ~30                         |
| Cr            | +4%                 | ~10                         |
| Co            | −3%                 | ~5                          |
| Al            | −13% (smaller)      | ~20                         |

### Effect of Alloying on Other Properties

| Property | Effect of Adding Solute |
|----------|-------------------------|
| Melting point | Usually lowers (depresses liquidus), but W/Re raise Ni's creep T |
| Electrical conductivity | Always decreases (solute scatters electrons) |
| Thermal conductivity | Usually decreases |
| Corrosion resistance | Depends: Cr → passivation; S → hotness; Ce → scale adhesion |
| Density | Depends on relative atomic weights (W heavy; Al light) |
| CTE (thermal expansion) | Usually changes; critical for coatings (must match) |

---

## 8. Alloy Designation Systems

Different industries use different systems to name alloys. Key ones to know:

**Steel:**
- AISI/SAE: 4 digits, e.g., 4340 = 4×xx molybdenum alloy steel, 0.40% C
- EN (European): 1.xxxx, e.g., 1.4404 = 316L stainless steel
- Stainless: 300 series (austenitic), 400 series (ferritic/martensitic)

**Aluminum:**
- 4-digit system: 1xxx (pure Al), 2xxx (Al-Cu), 3xxx (Al-Mn), 5xxx (Al-Mg), 6xxx (Al-Mg-Si), 7xxx (Al-Zn)
- T temper designations: T6 (peak aged), T7 (overaged), etc.

**Titanium:**
- Grade 1–4: commercially pure (CP)
- Grade 5: Ti-6Al-4V (most common), α+β alloy

**Nickel superalloys:**
- Proprietary designations: IN718, CMSX-4, PWA 1484, René N6, TMS-238
- UNS: N06600 = Inconel 600, etc.

**Copper:**
- Brass: Cu-Zn (C26000 = 70/30 brass)
- Bronze: Cu-Sn
- Beryllium copper: Cu-Be (C17200)

---

## 9. Multi-Component Alloys — The High-Entropy Concept

Traditional alloys have one dominant element (Fe, Al, Ni, Cu) with minor additions. But what if you use 5+ elements in roughly equal proportions?

**High-entropy alloys (HEA):** Equal or near-equal molar ratios of 5+ principal elements (Ch 61).

**Multi-principal element alloys (MPEA):** Broader term.

The key insight: with 5+ elements, the configurational entropy of mixing is so large that a single disordered FCC or BCC phase can be stable (entropy stabilizes the solid solution over intermetallics):
```
ΔS_mix = -R Σ(x_i ln x_i) 

For 5 equiatomic elements: ΔS_mix = R × ln(5) ≈ 1.61R = 13.4 J/mol·K
For 2 equiatomic: ΔS_mix = R × ln(2) ≈ 0.69R  = 5.7 J/mol·K
```

CMSX-4 is already a complex 9-component alloy, but it's still "Ni-base" — Ni is ~70 at%.
True HEAs have no dominant element.

Commercial nickel superalloys are moving toward more components (each generation adds elements — Re, Ru, Hf) — understanding multi-component phase equilibria via CALPHAD (Ch 64) is essential.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Solid solution | Solute atoms in solvent crystal; substitutional (replace) or interstitial (gaps) |
| Hume-Rothery rules | Size <15%, similar electronegativity/structure/valence → extensive solubility |
| Intermetallics | Fixed-composition, ordered structures; hard, strong, often brittle; γ′ is a benign exception |
| Ordered structures | Regular atomic arrangement; superdislocations; anomalous yield in Ni₃Al |
| Solid solution hardening | Δσ ∝ c^(1/2 to 2/3); size mismatch most effective; Mo/W/Re strongest in Ni |
| Precipitation hardening | Exploit solvus: solution treat → quench → age → fine precipitates → high strength |
| HEA | 5+ equiatomic elements; entropy stabilizes single phase; new alloy design paradigm |

---

## Exercises

1. Using Hume-Rothery rules, predict the solid solubility behavior for: (a) Ag in Cu (radii: Cu=1.28Å, Ag=1.44Å; both FCC; similar valence), (b) Zn in Cu (radii: Cu=1.28Å, Zn=1.33Å; Cu FCC, Zn HCP; Zn valence = 2, Cu = 1), (c) C in Fe (radii: Fe=1.26Å, C=0.77Å; BCC iron). For each, state which rule(s) are violated and what type of phase behavior to expect.

2. An aluminum alloy contains 2 at% Mg and 1 at% Cu. Calculate the solid solution hardening contribution from each element assuming: Mg causes Δσ = 25 MPa per at%, Cu causes Δσ = 40 MPa per at%. Are these contributions additive? What does this tell you about multi-component alloying design?

3. The Ni₃Al γ′ phase has an L1₂ ordered structure with lattice parameter a = 3.56 Å. Calculate: (a) the unit cell volume, (b) the atomic positions of the 4 Ni atoms and 1 Al atom per unit cell, (c) why a dislocation in the ordered phase must move as a pair (superdislocation) rather than a single dislocation.

4. A Cu-Ag alloy has a eutectic composition of 28 wt% Ag (T_eutectic = 779°C). An alloy with 10 wt% Ag is solution-treated at 900°C and quenched. On aging at 300°C, describe qualitatively: (a) what thermodynamic driving force exists for precipitation, (b) what microstructure you expect after under-aging vs. peak-aging vs. over-aging, (c) how you would confirm the presence of the precipitate phase using XRD (Ch 40).

5. Compare solid solution hardening vs. precipitation hardening as strengthening mechanisms for a high-temperature alloy operating at 1,000°C. Consider: (a) which mechanism is likely to remain effective at 1,000°C and why, (b) the risk of precipitate coarsening (Ch 44) for precipitation hardening, (c) why modern nickel superalloys use BOTH mechanisms simultaneously.
