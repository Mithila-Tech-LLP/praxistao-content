# Chapter 03: Semiconductors — Deep Dive

> **"Doping is the art of deliberately contaminating silicon — adding one foreign atom per billion silicon atoms to transform an insulator into a conductor. This precise control of impurity is what makes modern electronics possible."**

---

## Table of Contents
1. [Silicon Crystal Structure](#1-silicon-crystal-structure)
2. [Intrinsic Semiconductor Behavior](#2-intrinsic-semiconductor-behavior)
3. [The Hole Concept](#3-the-hole-concept)
4. [Doping — N-Type Semiconductors](#4-doping--n-type-semiconductors)
5. [Doping — P-Type Semiconductors](#5-doping--p-type-semiconductors)
6. [Carrier Concentrations in Doped Semiconductors](#6-carrier-concentrations-in-doped-semiconductors)
7. [Drift Current](#7-drift-current)
8. [Diffusion Current](#8-diffusion-current)
9. [Einstein Relation](#9-einstein-relation)
10. [Hall Effect](#10-hall-effect)
11. [Recombination and Generation](#11-recombination-and-generation)
12. [Semiconductor vs. Semiconductor Comparison](#12-semiconductor-vs-semiconductor-comparison)
13. [Compound Semiconductors](#13-compound-semiconductors)
14. [Summary](#14-summary)

---

## 1. Silicon Crystal Structure

### Why Silicon?

Silicon (Si) is element #14 with electron configuration **2, 8, 4**. The **4 valence electrons** are key:
- Silicon needs 4 more electrons to complete its outer shell (octet rule)
- It achieves this by **sharing** one electron with each of 4 neighboring silicon atoms
- This creates 4 **covalent bonds** per atom
- Result: a rigid, symmetric, 3D crystal structure

### Diamond Cubic Crystal Lattice

Silicon forms a **diamond cubic crystal structure** — identical to the structure of diamond (carbon):

```
Unit cell of Silicon diamond cubic:

       z
       │
       │    Si─────────────Si
       │   /│              /│
       │  / │             / │
       Si─────────────Si   │
       │  │   Si       │   │
       │  │  /●   ●\   │   Si
       │  │ / ●   ● \  │  /
       │  Si─────────────Si/
       │                 /
       └────────────────── x
      /
     y

● = additional Si atoms at (¼,¼,¼), (¾,¾,¼), (¾,¼,¾), (¼,¾,¾) positions
```

**Key numbers:**
- **Lattice constant:** a = **5.431 Å** (0.5431 nm) at 300K
- **Nearest-neighbor distance:** a√3/4 = **2.35 Å**
- **Atoms per unit cell:** 8 (4 from corners/faces + 4 interior)
- **Atom density:** 5.0 × 10²² atoms/cm³
- **Coordination number:** 4 (each Si has exactly 4 nearest neighbors)
- **Packing efficiency:** 34% (fairly open structure)

### The Covalent Bond in Silicon

Each silicon atom forms **4 covalent bonds**, sharing one electron with each neighbor:

```
         Si
         │
    Si───Si───Si    (2D representation of 3D tetrahedral bonds)
         │
         Si

Each line represents a covalent bond (2 shared electrons)
```

In 3D, the bonds point to the corners of a **tetrahedron** (bond angles = 109.47°).

**Bonding electrons are NOT free** — they're locked in bonds. This is why perfect silicon at 0K is an insulator.

### How Silicon Wafers Are Made

1. **Quartzite (SiO₂) → Metallurgical Grade Silicon (~98% pure)**
   - Reaction in arc furnace: SiO₂ + 2C → Si + 2CO (at ~2000°C)

2. **Metallurgical Grade Si → Electronic Grade Si (~99.9999999% = 9N pure)**
   - Siemens process: Si + 3HCl → SiHCl₃ (trichlorosilane) + H₂
   - Then: SiHCl₃ + H₂ → Si + 3HCl (decomposed on hot rod)
   - One impurity per **billion** silicon atoms!

3. **Czochralski Process (CZ) — growing single crystal boule:**
   - Melt electronic grade Si in quartz crucible at 1414°C
   - Touch seed crystal to melt surface
   - Slowly rotate and pull upward as Si crystallizes
   - Boule: **300mm diameter** (standard today), 1-2m long, 250kg
   - Dopant (B for P-type, P for N-type) added to melt to set substrate doping

4. **Wafer Preparation:**
   - Diamond wire saw slices boule into ~775μm thick wafers
   - Lapping (mechanical grinding) to uniform thickness
   - Chemical etch to remove damage
   - Chemical Mechanical Polishing (CMP): mirror finish, sub-nm roughness
   - Inspection and cleaning

---

## 2. Intrinsic Semiconductor Behavior

A **pure (undoped) semiconductor** is called **intrinsic**.

### At Absolute Zero (T = 0K)

- Every electron is locked in a covalent bond
- Valence band: completely full
- Conduction band: completely empty
- **No free carriers → no conductivity → perfect insulator**

### At Room Temperature (T = 300K)

- Thermal energy (kT = 26 meV) breaks some covalent bonds
- A freed electron goes to the conduction band (**free electron**)
- The broken bond leaves an empty position in the valence band (**hole**)
- Both electron and hole can move → **both contribute to current**

```
Bond breaking process:

Before: Si─●─Si  (bond intact)
             ↑
             thermal energy (kT)

After:  Si─  ─Si  +  e⁻ (free electron in conduction band)
           ↑
           hole (positive, in valence band)
```

### Intrinsic Carrier Concentration

```
ni = √(Nc × Nv) × exp(-Eg / 2kT)

For Silicon at 300K:
  Nc = 2.8 × 10¹⁹ cm⁻³
  Nv = 1.04 × 10¹⁹ cm⁻³
  Eg = 1.12 eV
  kT = 0.026 eV

ni(Si, 300K) ≈ 1.5 × 10¹⁰ cm⁻³

To put this in perspective:
  Si atomic density = 5 × 10²² atoms/cm³
  Free electrons at 300K = 1.5 × 10¹⁰ cm⁻³

  Ratio = 1.5×10¹⁰ / 5×10²² = 3 × 10⁻¹³

Only 3 out of every 10 TRILLION silicon atoms have a free electron at room temperature!
This is why pure silicon is such a poor conductor.
```

**Temperature dependence of ni (approximate doubling rule):**
- ni doubles for every ~11°C increase in temperature
- At 200°C, ni(Si) ≈ 10¹⁵ cm⁻³ → heavy current leakage
- This limits Si devices to < ~150°C operation

---

## 3. The Hole Concept

The **hole** is one of the most important (and confusing) concepts in semiconductor physics.

### What is a Hole?

When an electron breaks free from a covalent bond, it leaves behind a **missing electron in the bond**. This missing electron site is called a **hole**.

```
Si lattice with a hole:

    Si   Si   Si   Si
     \  /  \  /  \  /
      Si    Si    Si
     /  \  /  \  /  \
    Si   Si   [○]   Si   ← ○ = hole (missing electron)
     \  /  \  /  \  /
      Si    Si    Si
```

### Hole Movement

A hole doesn't actually move — instead, **neighboring electrons jump** into the hole, which makes the hole appear to move in the opposite direction:

```
Step 1:   Si─Si─[○]─Si    hole at position 3

Step 2:   Si─Si─[○]─Si    electron from position 2 jumps to fill position 3
          (e⁻ moves right → hole moves left)

Step 3:   Si─[○]─Si─Si    hole now at position 2
```

**Net result:** hole appears to move like a **positive charge carrier** with:
- Charge: **+q = +1.602×10⁻¹⁹ C** (positive!)
- Mass: **effective mass** (different from free electron mass)
  - In Si: mp* ≈ 0.56 m₀ (hole is "heavier" → moves slower)
  - mn* ≈ 0.26 m₀ (electron effective mass in Si)
- Mobility: μp < μn (holes slower than electrons)

### In an Electric Field

When E-field applied:
- **Electrons** drift **opposite** to E-field (toward positive terminal)
- **Holes** drift **same** direction as E-field (toward negative terminal)
- **Both** carry current in the **same conventional current direction**!

```
Applied E-field: ───────────→ (left to right)

Electrons: ←←←←← (move against field)
Holes:     →→→→→ (move with field)

Conventional current: →→→→→ (both contribute the same direction)
```

---

## 4. Doping — N-Type Semiconductors

**Doping** = intentionally adding tiny amounts of impurity atoms to silicon to dramatically change its conductivity.

### N-Type Doping (adding electrons)

Add a **Group V element** (5 valence electrons) to silicon lattice:
- **Phosphorus (P)** — most common N-type dopant
- **Arsenic (As)** — used in older processes, high solubility in Si
- **Antimony (Sb)** — slow diffuser, used for buried layers

**What happens when phosphorus replaces a silicon atom:**

```
Silicon lattice with phosphorus dopant:

    Si   Si   Si   Si
     \  /  \  /  \  /
      Si    Si    Si
     /  \  /  \  /  \
    Si   Si   [P]   Si   ← P has 5 valence electrons
     \  /  \  /  \  /    4 are used in bonds with Si neighbors
      Si    Si    Si      1 is EXTRA — loosely bound to P
                          ↑
                          This extra electron is easily freed!
```

**The extra electron:**
- At room temperature, thermal energy (~26 meV) easily frees it
  (ionization energy of P in Si is only ~45 meV — far less than Eg = 1.12 eV!)
- This freed electron goes to the conduction band
- The phosphorus atom becomes P⁺ (positively ionized **donor** atom)
- P⁺ is fixed in the lattice — it CANNOT move

**Result:**
- Each phosphorus atom donates **1 free electron** to the conduction band
- Called **donor atoms** (they donate electrons)
- Material is called **N-type** (negative charge carriers = electrons)

**Donor energy level in band diagram:**
```
Energy:
  ══════════════ CB (conduction band)
  -- -- -- -- --  ← Donor level Ed (just 45meV below CB for P in Si)
                     ← almost all donors ionized at room temperature
  ══════════════ VB (valence band)
```

The donor level is so close to CB that at 300K, essentially **100% of donors are ionized**.

### N-type Majority and Minority Carriers

If we add ND phosphorus atoms per cm³:

```
After ionization at room temperature:
  n ≈ ND    (electrons are majority carriers)
  p = ni²/ND    (holes are minority carriers — from mass action law)

Example: ND = 10¹⁶ cm⁻³
  n ≈ 10¹⁶ cm⁻³  (electrons)
  p = (1.5×10¹⁰)² / 10¹⁶ = 2.25×10²⁰ / 10¹⁶ = 2.25×10⁴ cm⁻³  (holes)

Compare: n/p = 10¹⁶/2.25×10⁴ ≈ 4.4×10¹¹
Electrons outnumber holes by 440 BILLION to 1!
```

---

## 5. Doping — P-Type Semiconductors

### P-Type Doping (creating holes)

Add a **Group III element** (3 valence electrons) to silicon lattice:
- **Boron (B)** — most common P-type dopant (small atom, fast diffuser)
- **Aluminum (Al)** — less common
- **Gallium (Ga)** — occasionally used
- **Indium (In)** — for specific applications

**What happens when boron replaces a silicon atom:**

```
Silicon lattice with boron dopant:

    Si   Si   Si   Si
     \  /  \  /  \  /
      Si    Si    Si
     /  \  /  \  /  \
    Si   Si   [B]   Si   ← B has only 3 valence electrons
     \  /  \  /  \  /    3 bonds made with Si neighbors
      Si    Si    Si      1 bond is INCOMPLETE (a hole!)
                          ↑
                          This hole accepts an electron from a neighboring bond
```

**The missing bond:**
- Boron needs 1 more electron to complete 4 bonds
- It can easily accept an electron from a neighboring Si-Si bond
  (ionization energy ~45 meV for B in Si)
- When B accepts an electron: B becomes B⁻ (negatively ionized **acceptor** atom)
- The Si-Si bond that donated the electron now has a **hole**
- B⁻ is fixed in lattice, cannot move

**Result:**
- Each boron atom accepts 1 electron → creates 1 free hole in valence band
- Called **acceptor atoms** (they accept electrons, creating holes)
- Material is called **P-type** (positive charge carriers = holes)

**Acceptor energy level in band diagram:**
```
Energy:
  ══════════════ CB (conduction band)


  -- -- -- -- --  ← Acceptor level Ea (just ~45meV above VB for B in Si)
  ══════════════ VB (valence band)
```

### P-type Majority and Minority Carriers

If we add NA boron atoms per cm³:

```
After ionization at room temperature:
  p ≈ NA    (holes are majority carriers)
  n = ni²/NA    (electrons are minority carriers)

Example: NA = 10¹⁷ cm⁻³
  p ≈ 10¹⁷ cm⁻³  (holes)
  n = (1.5×10¹⁰)² / 10¹⁷ = 2.25×10³ cm⁻³  (electrons)
```

---

## 6. Carrier Concentrations in Doped Semiconductors

### Charge Neutrality Condition

In any semiconductor, **total positive charges = total negative charges**:

```
p + ND⁺ = n + NA⁻

Where:
  p   = hole concentration
  ND⁺ = ionized donor concentration (positive)
  n   = electron concentration
  NA⁻ = ionized acceptor concentration (negative)

At room temperature, nearly all dopants ionized:
  ND⁺ ≈ ND (donor concentration)
  NA⁻ ≈ NA (acceptor concentration)

So: p + ND = n + NA
```

### General Solution

Combined with mass action law n·p = ni²:

**For N-type (ND >> NA):**
```
         ND - NA + √((ND-NA)² + 4ni²)
n =  ─────────────────────────────────
                    2

If ND >> ni: n ≈ ND
             p = ni²/ND
```

**For P-type (NA >> ND):**
```
         NA - ND + √((NA-ND)² + 4ni²)
p =  ─────────────────────────────────
                    2

If NA >> ni: p ≈ NA
             n = ni²/NA
```

### Degenerate vs Non-Degenerate Semiconductors

**Non-degenerate** (normal doping):
- ND or NA << Nc, Nv (effective density of states)
- Fermi level stays within the band gap
- Mass action law n·p = ni² holds
- Normal semiconductor behavior

**Degenerate** (very heavy doping):
- ND or NA comparable to or greater than Nc, Nv
- Fermi level enters the conduction band (N+) or valence band (P+)
- Material behaves more metal-like
- Mass action law no longer valid
- Used for: ohmic contacts to metal, tunnel diodes, heavily doped regions in BJTs

**Notation:**
- N+ or P+ = heavily doped (10¹⁸-10²⁰ cm⁻³)
- N- or P- = lightly doped (10¹⁴-10¹⁶ cm⁻³)

---

## 7. Drift Current

**Drift current** results from carrier movement under an applied **electric field**.

### Electron Drift Current Density

```
Jn,drift = q × n × μn × E

Where:
  Jn,drift = electron drift current density (A/cm²)
  q        = electron charge = 1.602×10⁻¹⁹ C
  n        = electron concentration (cm⁻³)
  μn       = electron mobility (cm²/V·s)
  E        = electric field (V/cm)
```

### Hole Drift Current Density

```
Jp,drift = q × p × μp × E

Where:
  p  = hole concentration (cm⁻³)
  μp = hole mobility (cm²/V·s)
```

### Total Drift Current Density

```
J_drift = Jn,drift + Jp,drift = q(nμn + pμp) × E = σ × E
```

### Velocity Saturation

At high electric fields, carriers can't keep accelerating — they reach a **saturation velocity**:

```
At low E:    vd = μ × E    (linear relationship)
At high E:   vd → vsat      (saturates)

vsat for Si:
  Electrons: ~10⁷ cm/s (10⁵ m/s)
  Holes:     ~8×10⁶ cm/s
```

This is why simply making transistors smaller (increasing E-field) doesn't keep increasing speed indefinitely.

---

## 8. Diffusion Current

**Diffusion current** results from carrier movement from **high concentration to low concentration** (like perfume spreading in a room) — **without** any electric field needed!

### Fick's First Law Applied to Carriers

**Electron diffusion current density:**
```
Jn,diff = q × Dn × (dn/dx)

Where:
  Dn  = electron diffusion coefficient (cm²/s)
  dn/dx = electron concentration gradient (cm⁻⁴)
  Sign: electrons diffuse from high→low concentration (positive current flows opposite)
```

**Hole diffusion current density:**
```
Jp,diff = -q × Dp × (dp/dx)

Where:
  Dp = hole diffusion coefficient (cm²/s)
  Negative sign: holes diffuse high→low, conventional current follows
```

### Total Current Density

```
Total: J = Jn + Jp

Jn = q×n×μn×E + q×Dn×(dn/dx)  ← drift + diffusion
Jp = q×p×μp×E - q×Dp×(dp/dx)  ← drift + diffusion

At equilibrium (no net current):
  drift + diffusion components cancel exactly
```

---

## 9. Einstein Relation

The **diffusion coefficient** and **mobility** are NOT independent — they're related by the **Einstein relation** (also called Einstein-Smoluchowski relation):

```
Dn     kT
──── = ──── = VT
μn      q

Dp     kT
──── = ──── = VT
μp      q

Where:
  VT = kT/q = thermal voltage ≈ 26 mV at 300K
```

**Numerical values at 300K:**
```
Silicon:
  μn = 1400 cm²/V·s  →  Dn = 1400 × 0.026 = 36.4 cm²/s
  μp = 450 cm²/V·s   →  Dp = 450 × 0.026 = 11.7 cm²/s

GaAs:
  μn = 8500 cm²/V·s  →  Dn = 8500 × 0.026 = 221 cm²/s
```

This relation is fundamental — it connects quantum statistics (kT term) to classical transport.

**Diffusion length** — average distance minority carriers travel before recombining:
```
Ln = √(Dn × τn)    (minority electron diffusion length in P-type)
Lp = √(Dp × τp)    (minority hole diffusion length in N-type)

Where τn, τp = minority carrier lifetime (typically 1μs to 1ms in Si)

Example: τn = 1μs, Dn = 36 cm²/s
  Ln = √(36 × 10⁻⁶) = √(36×10⁻⁶) ≈ 6×10⁻³ cm = 60μm
```

---

## 10. Hall Effect

The **Hall effect** is a phenomenon where a **magnetic field perpendicular to current flow** creates a **transverse voltage** — and it's used to measure carrier concentration and type.

### Physics

```
Current flow: →→→ (x-direction)
Magnetic field: ↑ (z-direction, out of page)

Electrons moving left (current right) experience magnetic force:
  F = q × v × B  (Lorentz force)
  Force on electrons: downward (y-direction)

Electrons accumulate at bottom → creates electric field
Equilibrium when: qE_Hall = qvB
                  E_Hall = vB
```

### Hall Voltage Formula

```
         IB
VH = ─────────
       q × n × t

Where:
  VH = Hall voltage (V)
  I  = current (A)
  B  = magnetic field (T)
  n  = carrier concentration (cm⁻³)
  t  = thickness of sample (cm)
  q  = electron charge

Hall coefficient: RH = 1/(q×n)  for n-type
                  RH = 1/(q×p)  for p-type (opposite sign!)
```

### Applications of Hall Effect

1. **Measuring carrier concentration:** measure VH → calculate n or p
2. **Determining carrier type:** sign of VH tells N-type vs P-type
3. **Hall effect sensors** (Hall ICs):
   - Used in: brushless DC motor control (position sensing)
   - Wheel speed sensors in cars (ABS brakes)
   - Current sensors (non-contact)
   - Position sensors in joysticks
   - Popular Hall ICs: SS49E (linear), A3144 (digital on/off), DRV5023

---

## 11. Recombination and Generation

### Generation

**Generation** = creating electron-hole pairs (breaking bonds)

Mechanisms:
1. **Thermal generation**: thermal energy breaks bonds (always happening at T > 0K)
   - Rate increases exponentially with temperature
2. **Optical generation**: photon with E > Eg absorbed → e-h pair created
   - Used in: photodiodes, solar cells, phototransistors, CCDs
3. **Impact ionization**: high-energy carrier collides with atom → creates e-h pair
   - Used in: avalanche photodiodes, occurs at breakdown in diodes/transistors

### Recombination

**Recombination** = electron meets hole → both disappear (annihilate)

Energy is released as:
- **Photon (light)**: in direct bandgap materials (GaAs, GaN, InP)
  - Basis of: LEDs, laser diodes!
  - Si is indirect bandgap → recombination mostly non-radiative (heat)
- **Heat (phonons)**: via lattice vibrations — common in Si (indirect gap)
- **Auger recombination**: energy transferred to another carrier (high doping)

**Recombination types:**
1. **Direct (band-to-band)**: electron jumps directly from CB to VB
   - Dominant in direct-gap semiconductors (GaAs, GaN)
2. **Indirect (via traps)**: via impurity/defect energy levels in gap (Shockley-Hall-Read)
   - Dominant in Si (indirect bandgap)
   - Traps can be intentionally added (Au in Si) to speed up recombination for fast switches
3. **Surface recombination**: defects at semiconductor surface act as traps
   - Si surface passivation with SiO₂ or Si₃N₄ reduces this

### Minority Carrier Lifetime (τ)

After creating excess carriers (by light flash or injection), they recombine exponentially:

```
Δn(t) = Δn₀ × exp(-t/τ)

Where:
  Δn(t) = excess electron concentration at time t
  τ      = minority carrier lifetime

Typical τ values in Si:
  High quality float-zone Si: up to 10ms
  Device-grade CZ Si: 100μs - 1ms
  Heavily doped Si: <1μs (more impurities → more recombination centers)
```

**Why lifetime matters:** In a BJT transistor, the stored charge in base = IB × τ. Short lifetime → faster switch-off. That's why some power transistors have gold doping to reduce lifetime.

---

## 12. Semiconductor vs Semiconductor Comparison

| Property | Si | Ge | GaAs | GaN | SiC |
|----------|----|----|------|-----|-----|
| Band gap (eV) | 1.12 | 0.67 | 1.42 | 3.4 | 3.26 |
| Gap type | Indirect | Indirect | Direct | Direct | Indirect |
| ni at 300K (cm⁻³) | 1.5×10¹⁰ | 2.4×10¹³ | 1.8×10⁶ | ~10⁻¹⁰ | ~10⁻⁶ |
| μn (cm²/V·s) | 1400 | 3900 | 8500 | 1000-2000 | 700 |
| μp (cm²/V·s) | 450 | 1900 | 400 | 30 | 120 |
| Max field (V/cm) | 3×10⁵ | 10⁵ | 4×10⁵ | 3.3×10⁶ | 3×10⁶ |
| Thermal cond (W/cm·K) | 1.5 | 0.6 | 0.46 | 1.3 | 4.9 |
| Max temp (°C) | 150 | 70 | 350 | 600 | 600+ |
| Substrate cost | Low | Medium | High | High | Very High |
| Mature tech? | Very | Yes | Yes | Growing | Growing |

### When to Use Which:

- **Si**: everything digital, low-frequency analog, power at moderate voltage/current
- **Ge**: high-speed RF (SiGe BiCMOS), IR detectors, some premium audio
- **GaAs**: RF/microwave MMICs, phased-array radar, LEDs (red/IR), laser diodes
- **GaN**: 5G base stations, power amplifiers, high-voltage power switches (EV chargers), blue LEDs
- **SiC**: very high voltage/current power electronics (EV traction inverters, solar inverters, grid)

---

## 13. Compound Semiconductors

### III-V Semiconductors (Group III + Group V)

These are **direct bandgap** materials — perfect for light emission.

**Gallium Arsenide (GaAs):**
- Eg = 1.42 eV (direct) → 870nm infrared light
- Used in: IR LEDs, laser diodes, solar cells, RF transistors (HBT, PHEMT)
- Zinc blende crystal structure (like Si but alternating Ga/As)
- Can form alloys: AlGaAs, InGaAs, GaAlAs

**Indium Phosphide (InP):**
- Eg = 1.34 eV (direct)
- Very high electron mobility
- Used in: long-haul fiber optic lasers (1310nm, 1550nm), high-speed electronics

**Gallium Nitride (GaN):**
- Eg = 3.4 eV (direct) → UV (365nm)
- Blue LEDs (Nobel Prize 2014 — Akasaki, Amano, Nakamura)
- White LEDs = blue LED + yellow phosphor
- Power electronics: GaN FETs for power supplies

**Indium Gallium Nitride (InGaN):**
- Alloy of InN and GaN
- Adjustable bandgap 0.7 to 3.4 eV depending on In fraction
- Red/green/blue/white LEDs by varying In content

### II-VI Semiconductors (Group II + Group VI)

**Zinc Selenide (ZnSe):** blue LEDs (older technology)
**Cadmium Telluride (CdTe):** solar cells (thin film)
**Mercury Cadmium Telluride (HgCdTe):** infrared detectors (thermal cameras)
**Zinc Oxide (ZnO):** UV LEDs, transparent conductor (ITO replacement)

### Silicon Carbide (SiC)

- Covalent compound (not III-V), exists in many polytypes (3C, 4H, 6H)
- 4H-SiC is most used for electronics
- Very high breakdown field → can make 10,000V power devices
- High thermal conductivity → power dense designs
- Used in: Tesla electric vehicle inverters, industrial motor drives

### 2D Materials (Next Generation)

**Graphene (single layer carbon):**
- Zero bandgap (semimetal) → can't be switched off → not ideal for transistors
- Extremely high mobility (>100,000 cm²/V·s)
- Transparent, flexible
- Applications: interconnects, electrodes, sensors

**Molybdenum Disulfide (MoS₂):**
- Single-layer MoS₂ has direct bandgap of 1.8 eV
- Very thin (0.65 nm) → ultimate gate control
- Potential future transistor channel material beyond Si

---

## 14. Summary

```
SEMICONDUCTOR FUNDAMENTALS
══════════════════════════

Crystal:
  Si → diamond cubic lattice, a = 5.431Å
  Each Si: 4 covalent bonds with neighbors

Intrinsic (pure Si):
  ni = 1.5×10¹⁰ cm⁻³ at 300K
  n = p = ni
  n × p = ni² (mass action law, ALWAYS true)

N-Type (add Group V: P, As, Sb):
  n ≈ ND >> p
  EF moves toward CB
  More electrons → better electron conductivity

P-Type (add Group III: B, Al, Ga):
  p ≈ NA >> n
  EF moves toward VB
  More holes → better hole conductivity

Current types:
  Drift:    J = q(nμn + pμp)E
  Diffusion: J = qDn(dn/dx) - qDp(dp/dx)

Einstein relation:
  D/μ = kT/q = VT = 26mV at 300K

Diffusion length:
  L = √(D×τ)  ← key BJT transistor parameter
```

### Key Numbers to Remember

| Quantity | Value |
|---------|-------|
| Si band gap | 1.12 eV |
| Si ni at 300K | 1.5×10¹⁰ cm⁻³ |
| Si lattice constant | 5.431 Å |
| Si atom density | 5×10²² /cm³ |
| μn (Si) | 1400 cm²/V·s |
| μp (Si) | 450 cm²/V·s |
| kT/q at 300K | 26 mV |
| Dn (Si) | 36 cm²/s |
| Dp (Si) | 11.7 cm²/s |
| Typical τ (Si) | 1-100 μs |

---

**← Previous:** [Chapter 02: Electrical Properties and Band Theory](./02-electrical-properties-and-band-theory.md)
**→ Next:** [Chapter 04: PN Junction and Diodes](./04-pn-junction-and-diodes.md)
