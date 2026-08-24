# Chapter 75: Semiconductors and Electronics

> **"The transistor — invented in 1947 — is the most important invention of the 20th century. Every smartphone contains billions of them. Semiconductors turned quantum physics into civilization-changing technology."**

---

## Table of Contents

- [75.1 Conductors, Insulators, and Semiconductors](#751-conductors-insulators-and-semiconductors)
- [75.2 Band Theory of Solids](#752-band-theory-of-solids)
- [75.3 Intrinsic Semiconductors](#753-intrinsic-semiconductors)
- [75.4 Doping: N-type and P-type Semiconductors](#754-doping-n-type-and-p-type-semiconductors)
- [75.5 The P-N Junction](#755-the-p-n-junction)
- [75.6 The Diode](#756-the-diode)
- [75.7 The Transistor](#757-the-transistor)
- [75.8 Integrated Circuits and Moore's Law](#758-integrated-circuits-and-moores-law)
- [75.9 Light-Emitting Diodes (LEDs)](#759-light-emitting-diodes-leds)
- [75.10 Solar Cells](#7510-solar-cells)
- [75.11 Modern Semiconductor Technology](#7511-modern-semiconductor-technology)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 75.1 Conductors, Insulators, and Semiconductors

All solids can be classified by how well they conduct electricity:

```
CONDUCTORS:
  Copper, silver, aluminum, gold
  Resistivity: ~10⁻⁸ Ω·m
  Free electrons easily drift → low resistance
  Examples: wires, circuit boards, heat sinks

INSULATORS:
  Glass, rubber, plastic, diamond
  Resistivity: ~10¹²+ Ω·m
  Electrons tightly bound → very high resistance
  Examples: wire coating, circuit board substrate

SEMICONDUCTORS:
  Silicon (Si), Germanium (Ge), Gallium Arsenide (GaAs)
  Resistivity: ~10⁻³ to 10³ Ω·m (between conductors and insulators)
  Key property: resistance changes dramatically with:
    - Temperature (decreases as T increases — opposite to metals)
    - Light (photoconductivity)
    - Impurities (doping — the key to electronics!)
```

---

## 75.2 Band Theory of Solids

In isolated atoms, electrons occupy discrete energy levels. In a solid, billions of atoms' levels merge into continuous **energy bands**:

```
ENERGY BAND DIAGRAM:

                    Conduction band (free electrons here)
   ─────────────────
       Band gap Eg
   ─────────────────
                    Valence band (electrons bound to atoms)
   ─────────────────────────────
   
   Below valence band: all filled (core electrons)

CONDUCTOR:           SEMICONDUCTOR:        INSULATOR:
  
  [Conduction]         [Conduction]          [Conduction]
  ── empty/partial──   ── empty ────         ── empty ──
  
  ── filled ──         ── gap Eg ──          ── large gap ──
  [Valence]            ── filled ──          ── filled ──
                       [Valence]             [Valence]
  
  Bands overlap        Eg = 1-2 eV           Eg = 5+ eV
  (easy conduction)    (thermal energy can   (too large for
                       bridge gap)           thermal promotion)
```

**Band gap Eg** is the minimum energy to move an electron from the valence band to the conduction band.

Silicon: Eg = 1.12 eV
Germanium: Eg = 0.67 eV
GaAs: Eg = 1.42 eV
Diamond: Eg = 5.5 eV (insulator)

At room temperature (kT ≈ 0.025 eV), silicon has a small number of electrons thermally promoted to the conduction band.

---

## 75.3 Intrinsic Semiconductors

A **pure (intrinsic) semiconductor** like pure silicon:

```
SILICON CRYSTAL STRUCTURE:

  Each Si atom has 4 valence electrons.
  Forms covalent bonds with 4 neighbors.
  All electrons are used in bonds.
  
  Si - Si - Si
  |    |    |
  Si - Si - Si
  |    |    |
  Si - Si - Si
  
  Pure silicon: very few free charge carriers at room temperature.
  Resistivity ≈ 640 Ω·m (much more resistive than copper at 1.7×10⁻⁸ Ω·m)
```

When a valence electron gains enough energy (from heat or light) to jump to the conduction band:

1. It becomes a free electron that can carry current
2. It leaves behind a **hole** — a missing electron, effectively a positive charge carrier

```
ELECTRONS AND HOLES:
  
  In conduction band: electrons (negative, move in direction opposite to E field)
  In valence band: holes (positive, move in direction of E field)
  
  Both contribute to current in intrinsic semiconductor.
  n (electrons) = p (holes) in intrinsic semiconductor
```

**Intrinsic carrier concentration** at room temperature for Si: ni ≈ 1.5 × 10¹⁰ /cm³ (compared to ~5 × 10²² Si atoms/cm³ — only 1 in 3 trillion atoms contributes!)

---

## 75.4 Doping: N-type and P-type Semiconductors

**Doping** = adding tiny amounts of impurity atoms to dramatically change conductivity.

### N-type Semiconductor

Add a **donor** impurity (Group V atom: phosphorus, arsenic, antimony) to silicon:

```
N-TYPE DOPING (with Phosphorus, 5 valence electrons):

  Si - Si - Si
  |    |    |
  Si - P  - Si   ← Phosphorus has 5 valence electrons
  |    |    |       4 used in bonds, 1 FREE electron!
  Si - Si - Si
  
  Each P atom donates one free electron.
  Result: MANY extra electrons (majority carriers)
  
  Electrically neutral overall (same number of + and - charges)
  but electrons are the majority mobile carriers → N-type
```

### P-type Semiconductor

Add an **acceptor** impurity (Group III: boron, aluminum):

```
P-TYPE DOPING (with Boron, 3 valence electrons):

  Si - Si - Si
  |    |    |
  Si - B  - Si   ← Boron has 3 valence electrons
  |    |    |       3 used in bonds, 1 bond MISSING → hole!
  Si - Si - Si
  
  Each B atom creates one hole.
  Result: MANY extra holes (majority carriers)
  
  Holes are the majority mobile carriers → P-type
```

```
COMPARISON:

N-type: majority carriers = electrons (from donor atoms)
        minority carriers = holes
        
P-type: majority carriers = holes (from acceptor atoms)
        minority carriers = electrons
        
Typical doping level: 1 dopant per 10⁷ silicon atoms → increases 
conductivity by factor ~10⁷ !
```

---

## 75.5 The P-N Junction

Join P-type and N-type silicon together. What happens at the interface?

```
P-N JUNCTION FORMATION:

P region          N region
- - - - |  + + + + +
- - - - |  + + + + +   (majority carriers)
- - - - |  + + + + +

Electrons from N diffuse into P.
Holes from P diffuse into N.

A depletion region forms near the junction:
(no free carriers — they've all recombined)

P region    DEPLETION    N region
             REGION
- - - [- - |+ + +] + + +
       ← negative ions near P side
           → positive ions near N side

These fixed ions create a built-in ELECTRIC FIELD across the junction.
This field opposes further diffusion → equilibrium reached.

Built-in potential (contact potential): V₀ ≈ 0.6-0.7 V for silicon
```

---

## 75.6 The Diode

A **diode** is a P-N junction device that allows current to flow in one direction only.

```
FORWARD BIAS:
  Apply + to P side, - to N side.
  
  P[- - -|+ + +]N
  ← V_forward (0.6-0.7V for Si)
  
  Reduces the barrier → charges flow freely → large current
  
  "On" state: diode conducts (resistance very low)

REVERSE BIAS:
  Apply + to N side, - to P side.
  
  P[- - -|+ + +]N
                 → V_reverse
  
  Increases the barrier → no current can flow
  (except tiny leakage current)
  
  "Off" state: diode blocks current (resistance very high)
```

```
I-V CURVE OF A DIODE:

Current (mA)
   |
   |              / <- forward bias, large current
   |             /
   |            /
   |           / V_threshold (~0.6V)
   |    ← 0 → /
---+----------/----> Voltage (V)
   |  reverse /
   |  bias   /   <- tiny reverse leakage
   |
   
   Diode conducts strongly only in forward direction.
```

### Rectification

A diode converts AC to pulsating DC (half-wave rectification):

```
AC input:         Output after diode:
  /\   /\           /\
 /  \ /  \      → /  \     /\
             \  /       \  /
              \/
              
Only positive half-cycles pass through.
```

**Full-wave rectification** uses a bridge of 4 diodes to use both halves of the AC cycle.

---

## 75.7 The Transistor

The **transistor** (1947, Shockley, Bardeen, Brattain) is the fundamental building block of modern electronics.

### BJT (Bipolar Junction Transistor)

Two types: NPN and PNP. 

```
NPN TRANSISTOR:

  Emitter (N) | Base (P) | Collector (N)
  
  Structure: N-P-N junction
  
  When small current flows into BASE:
    → controls large current from EMITTER to COLLECTOR
    
  → CURRENT AMPLIFIER: I_C ≈ β × I_B  (β = gain, typically 50-500)
  
  Also used as SWITCH:
    I_B = 0 → transistor OFF (no collector current)
    I_B = small → transistor ON (large collector current)
```

```
TRANSISTOR AS AMPLIFIER:

  Input (small signal) → Base → Output (amplified signal) → Collector
  
  A 1 mA input → 200 mA output (β = 200)
  
  Used in: microphone amplifiers, radio receivers, audio systems
  
TRANSISTOR AS SWITCH (digital logic):

  I_B = 0 → I_C = 0  → Output = HIGH voltage (OFF) = logic '1'
  I_B = high → I_C = large  → Output ≈ 0 V (ON) = logic '0'
  
  Used in: computers, phones, digital logic, memory
```

### MOSFET (Metal Oxide Semiconductor Field Effect Transistor)

The dominant transistor in modern chips. Voltage at the **gate** controls current between **source** and **drain**:

```
MOSFET STRUCTURE:

           Gate
            |
  Source—[channel]—Drain
            |
        Gate oxide (insulating layer)
        Silicon substrate
        
  Gate voltage creates/destroys conducting channel.
  Nearly infinite input impedance (gate is insulated → no DC gate current).
  Easier to miniaturize than BJT → used in integrated circuits.
```

---

## 75.8 Integrated Circuits and Moore's Law

An **integrated circuit (IC)** is many transistors (and other components) fabricated on a single chip of silicon.

```
CHIP FABRICATION (simplified):

1. Start with silicon wafer
2. Grow thin SiO₂ layer
3. Coat with photoresist
4. Expose with UV light through a mask (photolithography)
5. Develop: exposed resist dissolves
6. Etch or deposit material in the pattern
7. Repeat many times (30+ layers for modern chips)
8. Cut, package, test
```

### Moore's Law

Gordon Moore observed (1965) that the number of transistors per chip doubled roughly every 2 years:

```
TRANSISTOR COUNT OVER TIME:

Year:  1971  1980  1990  2000  2010  2020
Count: 2,300  30k  1.2M  42M   2.6B  50B
                        
Intel 4004: 2,300 transistors (1971)
iPhone chip: ~10 billion transistors (2020)
```

Current state (2024): transistors at 3 nm (3 nanometer process = feature size ~3 nm = ~15 silicon atoms). Physical limits are approaching.

---

## 75.9 Light-Emitting Diodes (LEDs)

When an electron crosses from the conduction band to the valence band (recombines with a hole), it can emit a photon:

```
LED PRINCIPLE:

  Forward biased P-N junction.
  
  Electrons in N side + holes in P side →
  electrons cross junction → recombine with holes →
  emit photon:  E_photon = E_g (band gap energy)
  
  λ = hc/E_g
```

Different semiconductor materials → different band gaps → different colors:

```
COLOR      MATERIAL        Eg (eV)    λ (nm)
Red        GaAsP           1.9        660
Orange     GaAsP           2.0        620
Yellow     GaP             2.3        590
Green      InGaN           2.5        520
Blue       InGaN           2.9        430
White      Blue LED + phosphor
```

**White LED** (used in lighting): a blue InGaN LED coated with yellow phosphor. Blue light excites phosphor → emits yellow → blue + yellow = white light.

Efficiency: ~50-70% (vs. 5% for incandescent, 20% for fluorescent) — LEDs are the most efficient white light source.

---

## 75.10 Solar Cells

A solar cell is a P-N junction that converts light to electricity — the reverse of an LED.

```
SOLAR CELL PRINCIPLE:

  1. Photon with E > Eg absorbed by semiconductor
  2. Creates electron-hole pair (photoelectric effect in semiconductor)
  3. Built-in electric field at P-N junction:
      - sweeps electrons toward N side
      - sweeps holes toward P side
  4. This charge separation → potential difference → current when connected to circuit
  
  P side → (-) terminal
  N side → (+) terminal
```

```
SOLAR CELL CIRCUIT:
  
  Sunlight
    ↓
  [P-N junction]
     |     |
     -     +
     |     |
  external load (gets electricity)
```

### Efficiency

- Single-crystal silicon: ~22-26%
- Multi-junction cells (3+ junctions): ~47% (in concentration systems)
- The theoretical limit for a single-junction cell (Shockley-Queisser limit) is ~33%

---

## 75.11 Modern Semiconductor Technology

### Memory

- **DRAM**: capacitors store charge (0/1); fast but needs refresh every few ms
- **SRAM**: 6 transistors per bit; very fast, no refresh; expensive (used in CPU cache)
- **Flash memory**: electrons stored on floating gate; non-volatile (data survives power-off)

### Processors

Modern CPUs and GPUs contain billions of transistors performing logic operations:
- Logic gates (AND, OR, NOT) from transistors
- Flip-flops (memory) from logic gates
- Adders, multipliers, ALUs from flip-flops and logic
- CPUs from these building blocks

### Quantum Computing

Emerging technology using quantum states (qubits) instead of classical bits. Qubits can be in superposition (0 and 1 simultaneously) and entangled:

```
CLASSICAL BIT:   0 or 1

QUBIT:           |0⟩ + |1⟩ (superposition — both at once!)

N qubits: can represent 2^N states simultaneously
For N=300: more states than atoms in the observable universe

Potential applications: cryptography, drug discovery, optimization
Current state: 1000+ qubit processors, but error rates still high
```

---

## Summary

- **Semiconductors**: Eg ≈ 1-2 eV; conductivity between metals and insulators; tunable
- **Band theory**: valence band (filled), conduction band (empty), band gap Eg
- **Intrinsic**: pure semiconductor; n = p = ni; few carriers at room T
- **N-type doping**: donor impurities (Group V); electrons are majority carriers
- **P-type doping**: acceptor impurities (Group III); holes are majority carriers
- **P-N junction**: depletion region forms; built-in potential ~0.6 V (Si)
- **Diode**: forward biased → conducts; reverse biased → blocks; used for rectification
- **BJT transistor**: small base current controls large collector current; switch/amplifier
- **MOSFET**: gate voltage controls channel; dominant in ICs; nearly infinite input impedance
- **Integrated circuits**: billions of transistors on a chip; Moore's Law
- **LED**: electron-hole recombination → photon; color set by band gap
- **Solar cell**: photon → electron-hole pair → current; efficiency ~25% for Si

---

## Key Equations

```
Band gap and photon:
  E_photon = E_g = hf = hc/λ

Diode current (Shockley equation):
  I = I₀(e^(eV/kT) - 1)
  e = 1.6 × 10⁻¹⁹ C
  k = 1.38 × 10⁻²³ J/K
  kT/e ≈ 0.026 V at room temperature

Transistor (BJT):
  I_C = β × I_B
  β = current gain (50-500 typical)

LED wavelength:
  λ = hc/E_g

Solar cell:
  P_max = FF × V_oc × I_sc
  (FF = fill factor, V_oc = open circuit voltage, I_sc = short circuit current)

Moore's Law (observed, not fundamental):
  Transistor count doubles every ~2 years

Energy scales:
  Room temperature: kT ≈ 0.025 eV
  Silicon band gap: 1.12 eV
  Ge band gap: 0.67 eV
  GaAs band gap: 1.42 eV
```
