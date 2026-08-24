# Chapter 74: Nuclear Physics and Radioactivity

> **"The nucleus is a tiny speck at the center of an atom — yet it contains 99.97% of the atom's mass, holds enough energy to power cities or destroy them, and its natural decay has been ticking away since the Big Bang."**

---

## Table of Contents

- [74.1 The Nucleus: Protons and Neutrons](#741-the-nucleus-protons-and-neutrons)
- [74.2 Nuclear Forces and Stability](#742-nuclear-forces-and-stability)
- [74.3 Nuclear Binding Energy](#743-nuclear-binding-energy)
- [74.4 Radioactive Decay](#744-radioactive-decay)
- [74.5 Alpha Decay](#745-alpha-decay)
- [74.6 Beta Decay](#746-beta-decay)
- [74.7 Gamma Decay](#747-gamma-decay)
- [74.8 Radioactive Decay Law](#748-radioactive-decay-law)
- [74.9 Half-Life](#749-half-life)
- [74.10 Nuclear Fission](#7410-nuclear-fission)
- [74.11 Nuclear Fusion](#7411-nuclear-fusion)
- [74.12 Applications](#7412-applications)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 74.1 The Nucleus: Protons and Neutrons

The atomic nucleus consists of two types of particles — **nucleons**:

```
PROTON:
  Charge: +e = +1.602 × 10⁻¹⁹ C
  Mass:   1.6726 × 10⁻²⁷ kg = 1.007276 u
  Symbol: p

NEUTRON:
  Charge: 0 (neutral)
  Mass:   1.6749 × 10⁻²⁷ kg = 1.008665 u
  Symbol: n

ATOMIC MASS UNIT (u):
  1 u = 1.66054 × 10⁻²⁷ kg
  1 u ≈ mass of one proton or neutron
```

### Notation

```
     A
     X  (element symbol X, mass number A, atomic number Z)
     Z

A = mass number = total nucleons = protons + neutrons
Z = atomic number = number of protons (defines the element)
N = neutron number = A - Z

Examples:
  ¹H: hydrogen (1 proton, 0 neutrons)
  ¹H₂: deuterium (1 proton, 1 neutron)
  ¹H₃: tritium (1 proton, 2 neutrons)
  ⁴He: helium-4 (2 protons, 2 neutrons)
  ¹²C: carbon-12 (6 protons, 6 neutrons)
  ²³⁵U: uranium-235 (92 protons, 143 neutrons)
```

**Isotopes**: nuclei with same Z (same element) but different N.
- Carbon-12 and Carbon-14 are both carbon (Z=6) but different masses.

---

## 74.2 Nuclear Forces and Stability

Why do protons stay together in the nucleus despite their mutual electromagnetic repulsion?

```
FORCES IN THE NUCLEUS:

Electromagnetic force: protons repel each other strongly
  (all protons, long range, 1/r²)

Strong nuclear force: binds nucleons together
  (all nucleons, short range, ~1 fm)
  
  The strong force is far more powerful than EM at short range.
  But it drops to zero beyond ~2-3 fm (very short range).
```

**Nuclear stability** depends on the balance of these forces:

```
STABILITY CURVE (N vs Z):

N
(neutrons)
   |              * stable nuclei
   |           *
   |         *  <- N/Z increases for heavy nuclei
   |       *
   |     * 
   |  *
   | *
   |*  (N=Z line)
   +-----------> Z (protons)
   
Light nuclei: N ≈ Z is stable (N/Z ≈ 1)
Heavy nuclei: need more neutrons to dilute proton repulsion (N/Z > 1)
Nuclei far from the curve: unstable → radioactive decay
```

---

## 74.3 Nuclear Binding Energy

The mass of a nucleus is LESS than the sum of its free protons and neutrons. This "missing mass" is the **mass defect**, and via E = mc², it represents the **binding energy**.

```
BINDING ENERGY:

Mass defect: Δm = Z × m_p + N × m_n - m_nucleus

Binding energy: BE = Δm × c²

This is the energy needed to completely separate the nucleus into
individual protons and neutrons.

Example: Helium-4 (2p + 2n):
  Mass of free particles = 2(1.00728) + 2(1.00867) = 4.03190 u
  Actual mass of He-4     = 4.00260 u
  Mass defect Δm          = 0.03031 u
  BE = 0.03031 × 931.5 MeV/u = 28.3 MeV
```

```
BINDING ENERGY PER NUCLEON (BE/A):

BE/A
(MeV)
   |
 9 |           ****Fe (most stable)
   |        **      **
 8 |     ***            *
   |  **                  *
 7 |**                      *
   |                          *
   +-------------------------> A (mass number)
   1       10       100      200
   H
   
Iron-56 has the highest binding energy per nucleon → most stable nucleus.
Lighter nuclei: can release energy by FUSION (building up to iron).
Heavier nuclei: can release energy by FISSION (splitting down toward iron).
```

---

## 74.4 Radioactive Decay

Unstable nuclei undergo **radioactive decay** — they spontaneously transform to reach a more stable configuration.

Three main types:
- **Alpha (α) decay**: nucleus emits a helium-4 nucleus
- **Beta (β) decay**: a neutron or proton changes type
- **Gamma (γ) decay**: nucleus emits a high-energy photon

All decay processes conserve:
- Total mass number A (nucleons)
- Total charge Z (protons)
- Energy and momentum

---

## 74.5 Alpha Decay

The nucleus emits an **alpha particle** (⁴He nucleus = 2 protons + 2 neutrons):

```
ALPHA DECAY:

   A      →    A-4      +    4
   X             Y           He
   Z             Z-2          2

Mass number decreases by 4.
Atomic number decreases by 2.

Example: Radium-226 → Radon-222 + Helium-4

  ²²⁶Ra  →  ²²²Rn  +  ⁴He
   88        86         2
```

**Properties of alpha radiation:**
- Doubly charged (+2) → strongly ionizing
- Massive → stopped by a sheet of paper or a few cm of air
- Dangerous if inhaled/ingested (inside body, great damage)
- Example source: smoke detectors (Americium-241)

---

## 74.6 Beta Decay

### Beta-minus (β⁻) Decay

A neutron converts to a proton:

```
n → p + e⁻ + ν̄_e  (antineutrino emitted)

Nuclear equation:
   A      →    A       +   0
   X             Y         e
   Z             Z+1        -1

Mass number: unchanged
Atomic number: increases by 1

Example: Carbon-14 → Nitrogen-14 + electron + antineutrino

  ¹⁴C  →  ¹⁴N  +  e⁻  +  ν̄_e
   6        7
```

### Beta-plus (β⁺) Decay (Positron emission)

A proton converts to a neutron:

```
p → n + e⁺ + ν_e  (neutrino emitted)

Atomic number: decreases by 1

Example: Carbon-11 → Boron-11 + positron + neutrino

  ¹¹C  →  ¹¹B  +  e⁺  +  ν_e
   6        5
```

**Properties of beta radiation:**
- Charged particles → ionizing
- More penetrating than alpha: stopped by ~3 mm aluminum or ~3 m air
- Used in: cancer treatment (beta emitters), PET scans (positron emitters)

**The neutrino** was proposed to conserve energy and momentum in beta decay. Nearly massless, they barely interact — trillions pass through your body every second from the Sun.

---

## 74.7 Gamma Decay

After alpha or beta decay, the daughter nucleus is often in an **excited state**. It transitions to the ground state by emitting a gamma photon.

```
GAMMA DECAY:

  X*  →  X  +  γ
  (excited)
  
  No change in A or Z — same element, just loses energy.
  
  Example: ⁵⁷Co undergoes beta decay → ⁵⁷Fe* (excited)
           ⁵⁷Fe* → ⁵⁷Fe + γ (0.14 MeV)
```

**Properties of gamma radiation:**
- No charge → weakly ionizing but very penetrating
- Stopped by: several cm of lead, or ~1 m of concrete
- Used in: cancer radiotherapy, sterilization, Mössbauer spectroscopy

### Summary Table: Types of Radiation

| Type | Particle | Charge | Penetration | Stopped by |
|------|---------|--------|-------------|------------|
| Alpha α | ⁴He nucleus | +2 | Low | Paper, skin |
| Beta β⁻ | Electron | -1 | Medium | 3 mm Al |
| Beta β⁺ | Positron | +1 | Medium | 3 mm Al |
| Gamma γ | Photon | 0 | High | Lead/concrete |

---

## 74.8 Radioactive Decay Law

Radioactive decay is a random quantum process. Each nucleus has a fixed probability per unit time of decaying, called the **decay constant λ** (not wavelength here).

The number of nuclei decreases exponentially:

```
dN/dt = -λN

Solution:
N(t) = N₀ × e^(-λt)

Activity (decays per second):
A(t) = λN = λN₀ × e^(-λt) = A₀ × e^(-λt)

Unit of activity: Becquerel (Bq) = 1 decay per second
                  Also: Curie (Ci) = 3.7 × 10¹⁰ Bq
```

---

## 74.9 Half-Life

The **half-life** T₁/₂ is the time for half the nuclei to decay:

```
At t = T₁/₂: N = N₀/2

From N = N₀e^(-λT₁/₂) = N₀/2:
  e^(-λT₁/₂) = 1/2
  -λT₁/₂ = ln(1/2) = -ln(2)
  
  T₁/₂ = ln(2)/λ = 0.693/λ
```

```
DECAY CURVE:

N
   |*
N₀ |
   | *
N₀/2|   *
   |    *
N₀/4|        *
   |              *
   +----+----+----+----> t
   0   T₁/₂  2T₁/₂  3T₁/₂

After each half-life, half the remaining nuclei have decayed.
```

### Half-Lives Vary Enormously

| Nucleus | Half-life | Application |
|---------|-----------|-------------|
| Uranium-238 | 4.5 × 10⁹ years | Geological dating |
| Carbon-14 | 5730 years | Archaeological dating |
| Cobalt-60 | 5.27 years | Cancer radiotherapy |
| Iodine-131 | 8 days | Medical thyroid treatment |
| Francium-223 | 22 minutes | Research |
| Beryllium-8 | 8.2 × 10⁻¹⁷ s | Stellar helium burning |

### Radiocarbon Dating

Carbon-14 (half-life 5730 years) is continuously produced in the atmosphere by cosmic rays and incorporated into living organisms. When an organism dies, no new ¹⁴C is added, and the ¹⁴C decays.

```
DATING CALCULATION:

  Measure ¹⁴C/¹²C ratio in sample vs. modern standard.
  
  If ratio = 50% of modern: 1 half-life has passed → age ≈ 5730 years
  If ratio = 25%: 2 half-lives → age ≈ 11,460 years
  If ratio = 12.5%: 3 half-lives → age ≈ 17,190 years
  
  Accurate to ~50,000 years (beyond that, ¹⁴C too low to measure)
```

### Worked Example 74.1

A sample initially has 80 g of radioactive iodine-131 (T₁/₂ = 8 days).

(a) Find the decay constant.
(b) How much remains after 24 days?

**Solution:**

(a) λ = 0.693/T₁/₂ = 0.693/8 = **0.0866 day⁻¹**

(b) 24 days = 3 half-lives

    N = N₀ × (1/2)³ = 80 × (1/8) = **10 g**

---

## 74.10 Nuclear Fission

**Nuclear fission**: a heavy nucleus splits into two smaller nuclei plus neutrons and energy.

```
FISSION OF URANIUM-235:

  ²³⁵U + ¹n → ²³⁶U* → ⁹²Kr + ¹⁴¹Ba + 3¹n + energy
   92            92      36       56
   
  Three neutrons released → can cause more fissions → CHAIN REACTION!
```

```
CHAIN REACTION:
  
  1 fission → 3 neutrons
  3 neutrons → 3 fissions → 9 neutrons
  9 fissions → 27 neutrons
  
  Exponential growth → EXPLOSION (bomb) if uncontrolled
  Controlled rate → NUCLEAR POWER PLANT
```

### Energy from Fission

The binding energy released per fission event:

- Uranium-235 fission releases ~200 MeV
- Burning coal: ~few eV per reaction
- Fission is ~50 MILLION times more energetic per atom!

```
NUCLEAR POWER PLANT:
  
  Fission heats water → steam → turbine → generator → electricity
  
  Control rods (boron): absorb neutrons → reduce reaction rate
  Moderator (water): slows neutrons → makes them more likely to cause fission
  Coolant (water): carries heat away
```

---

## 74.11 Nuclear Fusion

**Nuclear fusion**: light nuclei combine to form a heavier nucleus, releasing energy.

```
FUSION REACTIONS (in the Sun):

Step 1: ¹H + ¹H → ²H + e⁺ + ν  (proton-proton fusion)
Step 2: ²H + ¹H → ³He + γ
Step 3: ³He + ³He → ⁴He + ¹H + ¹H

NET: 4 protons → 1 helium-4 + energy

Energy released: ~26 MeV per He-4 formed
Mass converted: Δm = 4(1.00728) - 4.00260 = 0.02853 u
               = 0.02853 × 931.5 = 26.6 MeV ✓
```

Fusion requires extreme conditions:
- Temperature: ~10⁷ to 10⁸ K (so nuclei have enough energy to overcome repulsion)
- Pressure: enormous (to get nuclei close enough)

The Sun's core: T = 1.5 × 10⁷ K, P = 2.5 × 10¹⁶ Pa.

### Fusion vs Fission

```
COMPARISON:

             FISSION              FUSION
Fuel:        Uranium, Plutonium   Hydrogen isotopes (deuterium, tritium)
Products:    Radioactive waste    Helium (mostly harmless)
Trigger:     Neutron absorption   High temperature (plasma)
Status:      In use worldwide     Not yet commercially viable
Promise:     Current power        "Infinite" clean energy
```

Fusion power research (ITER in France): trying to confine plasma magnetically at 150,000,000 K. Not yet achieved net energy gain for sustained periods.

---

## 74.12 Applications

### Nuclear Power

~10% of world electricity (higher in some countries: France ~70%). Produces no CO₂ during operation.

### Medical Applications

- **PET Scan**: positron emitters (¹⁸F) injected → positrons annihilate with electrons → gamma rays → 3D image of metabolic activity
- **Radiotherapy**: gamma rays, beta emitters, proton beams to destroy tumors
- **Medical diagnostics**: Technetium-99m (6 hour half-life) for gamma camera imaging
- **Radiation therapy**: Cobalt-60, Iridium-192 for brachytherapy

### Radiometric Dating

- ¹⁴C: up to 50,000 years
- Potassium-40 / Argon-40: millions to billions of years (geological dating)
- Uranium/Lead: billions of years (age of Earth = 4.54 billion years)

### Nuclear Weapons

- **Fission bomb (A-bomb)**: supercritical mass of U-235 or Pu-239 → uncontrolled chain reaction
- **Fusion bomb (H-bomb)**: fission trigger compresses and heats fusion fuel → much larger yield

---

## Summary

- **Nucleus**: protons (Z) + neutrons (N); A = Z + N; isotopes have same Z, different N
- **Nuclear stability**: strong force (short range) vs electromagnetic repulsion; light nuclei N≈Z, heavy need N>Z
- **Binding energy**: BE = Δm × c²; iron-56 most stable; fusion/fission both release energy moving toward iron
- **Alpha decay**: emits ⁴He; A-4, Z-2; stopped by paper
- **Beta⁻ decay**: n→p+e⁻+antineutrino; A same, Z+1; stopped by mm Al
- **Gamma decay**: excited nucleus emits photon; A,Z unchanged; stopped by lead
- **Decay law**: N = N₀e^(-λt); A = A₀e^(-λt)
- **Half-life**: T₁/₂ = 0.693/λ; time for half to decay
- **Fission**: heavy nucleus splits; ~200 MeV per event; chain reaction; nuclear power and weapons
- **Fusion**: light nuclei combine; powers stars; deuterium + tritium most promising; requires extreme T

---

## Key Equations

```
Nuclear notation:
  A = Z + N  (mass number = protons + neutrons)

Binding energy:
  Δm = Z·m_p + N·m_n - m_nucleus
  BE = Δm × c²
  1 u × c² = 931.5 MeV

Radioactive decay law:
  N(t) = N₀ × e^(-λt)
  A(t) = A₀ × e^(-λt)  (activity)

Half-life:
  T₁/₂ = ln(2)/λ = 0.693/λ

Decay constant:
  λ = 0.693/T₁/₂

After n half-lives:
  N = N₀ × (1/2)^n

Einstein's mass-energy:
  E = Δm × c²
  c = 3 × 10⁸ m/s
  1 u = 1.66054 × 10⁻²⁷ kg = 931.5 MeV/c²
```
