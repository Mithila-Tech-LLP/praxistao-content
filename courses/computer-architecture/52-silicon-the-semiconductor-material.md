# Chapter 52: Silicon — The Semiconductor Material

Every chip in this course — CPUs, GPUs, FPGAs, ASICs — is made of silicon. Not because silicon is the best semiconductor theoretically, but because it occupies a unique sweet spot of properties that make it ideal for mass-produced transistors: abundant (second most common element in Earth's crust), grows an excellent native oxide (SiO₂ = glass), compatible with CMOS manufacturing, and its electrical properties can be precisely controlled with dopants. This chapter explains the physics of semiconductors, why silicon was chosen, how dopants create p-type and n-type regions, how a MOSFET works at the transistor level, and where silicon might be replaced by new materials.

## Table of Contents

1. [What Is a Semiconductor?](#1-what-is-a-semiconductor)
2. [Doping: Creating n-type and p-type Silicon](#2-doping-creating-n-type-and-p-type-silicon)
3. [The MOSFET — The Transistor in Your CPU](#3-the-mosfet--the-transistor-in-your-cpu)
4. [CMOS: Putting n-MOSFET and p-MOSFET Together](#4-cmos-putting-n-mosfet-and-p-mosfet-together)
5. [Silicon vs Other Semiconductors](#5-silicon-vs-other-semiconductors)
6. [The Limits of Silicon](#6-the-limits-of-silicon)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What Is a Semiconductor?

Materials fall into three electrical categories:

**Conductors** (copper, aluminum): Electrons flow freely. Resistivity ~10⁻⁸ Ω·m.
**Insulators** (glass, rubber): Electrons are tightly bound. Resistivity ~10¹⁰–10²⁰ Ω·m.
**Semiconductors** (silicon, germanium): In between. Resistivity ~10⁻⁴–10⁴ Ω·m. Controllable.

The magic of semiconductors is that their conductivity can be controlled — by temperature, light, electric field, or adding impurity atoms (doping). This controllability is what makes transistors possible.

**Band theory**: In a crystal, quantum mechanics causes electron energy levels to form continuous "bands." The **valence band** is filled with electrons; the **conduction band** is empty (at absolute zero). The **band gap** is the energy gap between them.

```
Band structure:
  
  Insulators:          Semiconductors:         Conductors:
  ┌──────────┐         ┌──────────┐            ┌──────────┐
  │ Conduction│  Large  │Conduction│  Small      │Conduction│
  │  Band    │  gap    │  Band    │  gap        │  Band    │ Overlap
  └──────────┘         └──────────┘            └──────────┘
  Large gap            Small gap               No gap
  ~5-9 eV             ~1-3 eV                 (conduction)
  
  Silicon band gap: 1.12 eV at room temperature
  Germanium: 0.67 eV
  Diamond: 5.47 eV (insulator)
  GaN: 3.4 eV (wide-bandgap semiconductor)
```

**Silicon crystal structure**: Silicon atoms form a **diamond cubic** crystal lattice — each silicon atom is covalently bonded to 4 neighbors in a tetrahedral arrangement. At room temperature, thermal energy breaks some covalent bonds, creating free electron-hole pairs that enable electrical conduction.

### Quick Check
> 1. What is a semiconductor and how does it differ from a conductor and insulator?
> 2. What is a "band gap" and what is silicon's band gap?
> 3. What is a "hole" in semiconductor physics?

---

## 2. Doping: Creating n-type and p-type Silicon

Pure silicon is a poor conductor at room temperature (intrinsic carrier concentration ~10¹⁰ cm⁻³). By adding tiny amounts of impurity atoms (**dopants**), we can increase conductivity by factors of 10⁶ and control the type of charge carriers.

**n-type silicon** (negative carriers = electrons):
- Add **phosphorus (P)** or **arsenic (As)** atoms (Group V: 5 valence electrons)
- Silicon needs 4 bonds; Group V atom has 4 bonds + one leftover electron
- That extra electron is weakly bound and easily freed to conduct electricity
- Result: excess **free electrons** as majority carriers

**p-type silicon** (positive carriers = holes):
- Add **boron (B)** atoms (Group III: 3 valence electrons)
- Silicon needs 4 bonds; boron has only 3 — leaves a "hole" (missing electron)
- Neighboring electrons can fill the hole, making the hole appear to "move"
- Result: excess **holes** as majority carriers

```
Doping visualization:
  
  n-type (phosphorus):          p-type (boron):
  
  Si - Si - Si - Si             Si - Si - Si - Si
  |    |    |    |              |    |    |    |
  Si - Si - P  - Si             Si - Si - B  - Si
  |    |  e⁻ |    |             |    |  ○  |    |
  Si - Si - Si - Si             Si - Si - Si - Si
         (free electron)               (hole ○)
  
  Typical dopant concentration: 10¹⁵–10¹⁸ cm⁻³
  Pure silicon: 5×10²² Si atoms/cm³
  Dopant fraction: 1 in 10⁴–10⁷ atoms
```

**p-n junction**: When p-type and n-type silicon are joined, free electrons from n-side diffuse into p-side (and holes diffuse the other way), creating a **depletion region** with no free carriers and a built-in electric field (~0.7V for silicon). This is a diode — current flows easily in one direction, blocked in the other. The p-n junction is the fundamental building block of all transistors.

### Quick Check
> 1. What dopant creates n-type silicon and why?
> 2. What is a "hole" in p-type silicon and how does it conduct electricity?
> 3. What is a p-n junction and what does the depletion region represent?

---

## 3. The MOSFET — The Transistor in Your CPU

The **MOSFET (Metal-Oxide-Semiconductor Field-Effect Transistor)** is the transistor used in virtually all modern digital chips. The largest modern chips (like Apple's M-series Max chips or multi-die GPUs) contain tens of billions to over 100 billion MOSFETs.

An n-channel MOSFET (nMOS) has:
- **Source**: n-type region (source of electrons)
- **Drain**: n-type region (destination of electrons)
- **Gate**: metal (or polysilicon) electrode above a thin oxide (SiO₂) layer
- **Body/Substrate**: p-type silicon below

```
nMOS transistor cross-section:
  
  Gate (metal)
       │
  ─────┴─────────────────────
  │  Gate oxide (SiO₂, 1nm) │
  ─────────────────────────────
  n+Source │  p-substrate  │ n+Drain
           │   (depletion) │
     
  ─────────────────────────────
              Substrate (p-type)
  
  Two n-type islands in p-type body = two back-to-back p-n junctions
  Normally: no current flows (like two diodes in series, blocking in both dirs)
```

**How a MOSFET works:**

1. **Gate voltage = 0V (OFF state)**: No electric field. The p-type body blocks electron flow between source and drain. OFF (no current).

2. **Gate voltage > V_threshold (ON state)**: The positive gate voltage attracts electrons to the surface of the p-type body, forming an **inversion layer** (a thin channel of n-type-like electrons). Now n-source, n-channel, n-drain are connected. Current flows.

```
MOSFET operation:
  
  V_G = 0V (OFF):           V_G > V_th (ON):
  
  Gate ─── 0V              Gate ─── +1V
  ─────────────            ─────────────
  ▓▓▓ oxide ▓▓▓            ▓▓▓ oxide ▓▓▓
  ─────────────            ════════════  ← inversion layer (electrons)
  n+ │  p  │ n+            n+ │channel│ n+
  S  │block│ D             S  ──────── D
     no current                current flows!
```

**Gate oxide thickness**: The oxide under the gate is what creates the field-effect. It must be thin enough for the electric field to reach through and create the inversion channel. In modern chips, this is ~1–2 nm thick (about 4–8 silicon atom layers). Making it thinner improves switching speed but causes **gate leakage current** (quantum tunneling through the oxide).

**MOSFET as switch**: Logic 1 (high voltage) = ON, Logic 0 (low voltage) = OFF. This is the physical basis of digital logic.

### Quick Check
> 1. What are the four terminals of a MOSFET?
> 2. What happens in the MOSFET when the gate voltage exceeds the threshold voltage?
> 3. What is the "inversion layer" and why is it important?

---

## 4. CMOS: Putting n-MOSFET and p-MOSFET Together

**CMOS (Complementary MOS)** uses both nMOS and pMOS transistors in complementary pairs. This is what almost every modern digital chip uses.

A **pMOS** transistor is the complement of nMOS: it uses p-type source/drain in an n-type body, and turns ON when gate voltage is LOW (near 0V), OFF when gate is HIGH.

```
CMOS NOT gate (inverter) — the simplest CMOS gate:
  
  V_DD (+1V) ─── pMOS source
                  │
                  └── pMOS gate ──── Input
                  │
               pMOS drain ─── Output ─── nMOS drain
                                              │
               nMOS gate ──── Input      nMOS source
                                              │
                                            GND (0V)
  
  When Input = 0V:
    pMOS: gate=0V, VGS = -1V < -Vth → pMOS ON → Output connected to VDD
    nMOS: gate=0V, VGS = 0V < Vth → nMOS OFF → Output disconnected from GND
    Result: Output = 1V (logic 1)
  
  When Input = 1V:
    pMOS: gate=1V, VGS = 0V → pMOS OFF
    nMOS: gate=1V, VGS = 1V > Vth → nMOS ON → Output connected to GND
    Result: Output = 0V (logic 0)
```

**Why CMOS is power-efficient**: In the steady state (input stable), either the pMOS or nMOS is OFF — there is no direct path from VDD to GND. Power is only consumed during switching (charging/discharging the capacitance of the next gate). This is **dynamic power consumption**.

**Dynamic power**: P = α × C × V² × f
- α = activity factor (fraction of cycles a gate switches)
- C = capacitance of gate output (load)
- V = supply voltage
- f = clock frequency

This formula explains why reducing supply voltage V dramatically reduces power (squared relationship), why higher clock frequency f directly increases power, and why power-gating (α≈0) saves power.

### Quick Check
> 1. What is the key structural difference between an nMOS and a pMOS transistor?
> 2. Why does a CMOS inverter consume almost no static power?
> 3. The dynamic power formula P = αCV²f explains why reducing voltage is so effective. If you halve the voltage: by what factor does dynamic power decrease?

---

## 5. Silicon vs Other Semiconductors

Silicon is not the only option. Other semiconductors offer different trade-offs:

**Germanium (Ge)**:
- Band gap: 0.67 eV (vs Si 1.12 eV) → lower threshold voltage, faster at low voltage
- Electron mobility: 2× higher than silicon → faster switches
- Problems: Poor native oxide, weaker thermally (lower melting point), expensive
- Current use: Intel FinFETs and gate-all-around FETs use SiGe alloy in the source/drain regions to improve carrier mobility

**Gallium Arsenide (GaAs)**:
- Electron mobility: 6× higher than silicon
- Direct band gap: enables light emission (LEDs, lasers) → optoelectronics
- Problems: No good native oxide, expensive wafers (vs $100 Si wafer), brittle
- Used in: RF amplifiers, solar cells, high-speed ICs for telecom, LEDs

**Gallium Nitride (GaN)**:
- Wide band gap: 3.4 eV → high breakdown voltage, operates at high temperatures
- Used in: Power electronics (EV chargers, motor drives), RF power amplifiers (5G base stations)
- GaN-on-silicon: growing GaN devices on silicon wafers to reduce cost
- Not suitable for general-purpose logic (band gap too wide)

**Silicon Carbide (SiC)**:
- Very wide band gap: 3.26 eV, extremely high breakdown voltage
- Used in: High-power, high-temperature electronics (EV power inverters, industrial)
- Tesla Model 3 uses SiC MOSFETs in its drive inverter

**Carbon Nanotubes (CNTs)** and **Graphene**:
- Research materials; extraordinary electron mobility
- Challenge: Cannot yet grow with required uniformity for mass production

**Indium Gallium Arsenide (InGaAs)**:
- Very high electron mobility; used in high-speed RF and photonic devices
- Long researched (including by Intel) as a possible high-mobility transistor channel material, but not yet used in production CPU logic

### Quick Check
> 1. Why does GaN succeed silicon in power electronics but not in logic chips?
> 2. Why is GaAs used for RF amplifiers but not for mass-market CPUs?
> 3. What is SiGe and why do modern Intel/AMD chips use it?

---

## 6. The Limits of Silicon

Silicon CMOS is approaching fundamental physical limits:

**Quantum tunneling**: Gate oxide is now ~1nm — a few atomic layers. Electrons tunnel through quantum mechanically, causing leakage current even when the transistor should be OFF. High-K dielectrics (HfO₂) solve this by providing better electrical insulation at greater physical thickness.

**Short-channel effects**: When the transistor channel length approaches 5nm, the source and drain electric fields interfere, making it hard to fully turn the transistor OFF (subthreshold leakage). FinFET and gate-all-around architectures surround the channel from multiple sides to regain control.

**Heat**: At high transistor density (100 million transistors/mm²), heat generation limits clock speed. Most modern CPUs are thermally limited, not electrically limited.

**Dopant variability**: A 5nm transistor has only a few hundred dopant atoms in the channel. Statistical fluctuations in dopant placement cause transistors on the same chip to behave differently.

```
Scaling limits:
  
  1970: 10µm gate length     → well above atomic scale
  1990: 0.5µm = 500nm        → still well above
  2000: 130nm                → room for much more scaling
  2010: 32nm                 → quantum effects beginning
  2015: 14nm FinFET          → 3D transistor required
  2020: 5nm                  → actual gate length ~18nm (~80 atoms)
  2025: 2nm-class GAA        → actual gate length ~12–16nm
        (TSMC N2, Intel 18A)
  
  Note: node names like "5nm" and "2nm" are marketing labels —
  no feature on the chip is actually that size (see Chapter 56).
  
  A silicon atom is 0.22nm wide
  A gate length of ~10nm = ~45 silicon atoms — and squeezing much
  below that runs into quantum effects.
  Physics hard limit approaching
```

**Solutions**: New transistor architectures (GAA nanosheets, CFET), new materials (2D materials like MoS₂), new computing paradigms (3D stacking, specialized processors). Chapter 56 covers process nodes in detail.

### Quick Check
> 1. What is quantum tunneling and how does it affect transistors as they get smaller?
> 2. What is a "short-channel effect"?
> 3. Is any feature on a "2nm" chip actually 2nm wide? Roughly how many silicon atoms span a modern transistor's gate?

---

## Summary

- **Semiconductors** have a controllable band gap (~1.12 eV for silicon) that enables transistors.
- **Doping**: adding Group V atoms (phosphorus) creates n-type; Group III atoms (boron) creates p-type. Dopants control the number and type of charge carriers.
- **MOSFET**: gate voltage controls an inversion channel between source and drain. Below threshold voltage: OFF. Above threshold: ON.
- **CMOS**: complementary nMOS + pMOS pairs. No static current draw; power only during switching. P = αCV²f.
- **Silicon advantages**: abundant, great native oxide (SiO₂), mature manufacturing. GaAs (fast, optoelectronics), GaN (power/RF), SiC (power/high-temp) are specialized alternatives.
- **Physical limits**: quantum tunneling through gate oxide, short-channel effects, heat density, dopant variability at single-digit nm scale.

---

## Exercises

### Easy
1. What is a semiconductor and why is silicon useful for making transistors?
2. What does "doping" do to silicon and why is it necessary?
3. What are the four terminals of a MOSFET and what does each one do?

### Medium
4. CMOS power calculation: A 3 GHz processor has 10 billion CMOS gates. Each gate has capacitance C = 1 fF, switches with activity factor α = 0.1, and supply voltage V = 1V. (a) Calculate total dynamic power using P = αCV²f. (b) If voltage is reduced from 1V to 0.7V: what is the new power? What is the percentage savings? (c) If clock frequency doubles from 3 GHz to 6 GHz (same voltage): new power? (d) Why can modern laptops switch between 3 GHz and 1 GHz modes to save battery?
5. Dopant statistics: A 5nm-class transistor has a channel region approximately 5nm × 10nm × 3nm = 150nm³. (Careful with unit conversion: 1 cm³ = 10²¹ nm³.) (a) Volume of the channel region in cm³ (5nm × 10nm × 3nm). (b) At silicon atomic density 5×10²² cm⁻³, how many silicon atoms are in this volume? (c) At n-type doping of 10¹⁸ cm⁻³, how many dopant atoms are in this channel? (d) What is the statistical standard deviation (√N) of the dopant count, and what percentage threshold voltage variation does this cause?
6. p-n junction rectification: A silicon diode has forward voltage drop of 0.7V. In a circuit with 5V supply and 1kΩ series resistor: (a) forward biased: what current flows? (b) reverse biased at -5V: what current flows (assume ideal)? (c) In a transistor, two p-n junctions are formed (source-body and body-drain). When V_GS < V_th: the body-drain junction is reverse biased. What does this mean for current flow? (d) This is why MOSFET is "off" without gate voltage — explain using p-n junction physics.

### Hard
7. Transistor scaling consequences: Intel's 1µm process (1989) had: Vth = 0.9V, tox = 20nm, Lmin = 1µm. A modern 5nm-class process (e.g., TSMC N5, 2020) has roughly: Vth = 0.25V, tox = 1nm (equivalent oxide thickness, using HfO₂), gate length ~18nm. (a) What is the ratio of gate oxide thickness? How has this changed the tunneling current? (b) At Vth = 0.25V, how close is threshold voltage to thermal noise kT/q = 26mV at room temperature? What does this mean for transistor OFF-state leakage? (c) Why do modern chips have both high-Vth (logic OFF) and low-Vth (logic speed) transistors mixed on the same chip? (d) What is the "subthreshold slope" and why is its theoretical minimum 60 mV/decade important?
8. Alternative material trade-offs: Your company is designing: (a) a 5G millimeter-wave base station power amplifier, (b) a 1kW motor drive inverter for an electric vehicle, (c) a high-performance mobile CPU, (d) a chip for LED lighting control. For each application: select the best semiconductor material (Si, GaAs, GaN, SiC, Ge) and justify your choice with specific properties. What are the cost, temperature, and efficiency trade-offs?
