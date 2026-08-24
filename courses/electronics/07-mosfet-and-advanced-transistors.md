# Chapter 07: MOSFETs and Advanced Transistors

> **"The MOSFET is the workhorse of modern electronics. Every digital chip — from the microcontroller in your toothbrush to the processor in a supercomputer — is built almost entirely from MOSFETs. Understanding the MOSFET is understanding modern computing."**

---

## Table of Contents
1. [MOSFET Fundamentals](#1-mosfet-fundamentals)
2. [MOSFET Types](#2-mosfet-types)
3. [MOSFET Operation — N-channel Enhancement](#3-mosfet-operation--n-channel-enhancement)
4. [MOSFET Equations and Regions](#4-mosfet-equations-and-regions)
5. [CMOS Technology](#5-cmos-technology)
6. [MOSFET as a Switch](#6-mosfet-as-a-switch)
7. [Power MOSFETs](#7-power-mosfets)
8. [JFET](#8-jfet-junction-field-effect-transistor)
9. [IGBT](#9-igbt-insulated-gate-bipolar-transistor)
10. [Advanced MOSFET Structures](#10-advanced-mosfet-structures)
11. [Transistor Comparison](#11-transistor-comparison)
12. [Transistor Numbering Systems](#12-transistor-numbering-systems)
13. [Summary](#13-summary)

---

## 1. MOSFET Fundamentals

**MOSFET = Metal Oxide Semiconductor Field Effect Transistor**

Invented simultaneously (disputed) by:
- Dawon Kahng and Martin Atalla at Bell Labs (1959) — first practical MOSFET
- Applied Fairchild's planar process to make it manufacturable

**Why MOSFETs dominate digital circuits (vs BJT):**
| Property | MOSFET | BJT |
|----------|--------|-----|
| Control mechanism | Voltage (gate voltage) | Current (base current) |
| Steady-state gate current | ~0 (essentially none!) | IB = IC/β (must supply current) |
| Switching speed | Very fast | Slower (minority storage) |
| Input impedance | >10¹⁴ Ω (insulated gate!) | 1kΩ-1MΩ |
| Static power in CMOS | Near zero | Always some |
| Scalability | Excellent (scales to <2nm) | Poor at small sizes |
| Density | Very high | Lower |

### MOSFET Structure

```
            Gate (G) ← Metal (aluminum, then polysilicon, now metal again)
               │
           ────┴──── ← Gate oxide (SiO₂ or high-k dielectric)
          │        │
     ─────┴────────┴─────
    │  S           D    │ ← P-substrate (for NMOS)
    │  N+         N+    │ ← Source (S) and Drain (D): heavily doped N+
    │   └──channel──┘   │
     ─────────────────────
          Substrate/Body (B)

4 terminals: Gate (G), Source (S), Drain (D), Body/Substrate (B)
(Often 3-terminal: body tied to source)
```

**Gate is INSULATED from channel by gate oxide!**
- Thickness: ~100nm (old) → ~1.2nm (modern 5nm node!)
- This insulation is why gate draws NO DC current
- Gate voltage creates electric field → controls channel

---

## 2. MOSFET Types

### Classification by Channel Polarity

**N-channel MOSFET (NMOS):**
- Channel carries electrons (negative carriers)
- Source and Drain: N+ doped
- Substrate: P-type
- Positive gate voltage creates N-channel
- Faster (electrons have higher mobility than holes)
- Lower on-resistance for same size

**P-channel MOSFET (PMOS):**
- Channel carries holes (positive carriers)
- Source and Drain: P+ doped
- Substrate: N-type
- Negative gate voltage creates P-channel
- Slower than NMOS

### Classification by Mode

**Enhancement Mode (normally OFF):**
- No channel at VGS = 0V — gate must create channel
- Requires gate voltage > threshold to turn ON
- **Most common type** in digital circuits
- N-channel: need VGS > Vth (positive voltage)
- P-channel: need VGS < Vth (negative voltage)

**Depletion Mode (normally ON):**
- Built-in channel exists at VGS = 0V
- Need gate voltage to REDUCE or DEPLETE the channel
- Used in: JFETs, some RF transistors, self-biasing circuits
- Less common

### The Four Types

```
Type                   Symbol behavior
─────────────────────────────────────────────────────────
N-channel Enhancement: OFF when VGS=0, ON when VGS > Vth(+)
N-channel Depletion:   ON when VGS=0, can turn off with VGS < 0
P-channel Enhancement: OFF when VGS=0, ON when VGS < Vth(-)
P-channel Depletion:   ON when VGS=0, can turn off with VGS > 0
```

---

## 3. MOSFET Operation — N-channel Enhancement

Let's trace exactly what happens in an N-channel enhancement MOSFET:

### Step 1: VGS = 0V (No channel)

```
   G (0V)
   │
   ├── Gate oxide (SiO₂) ────────────────────────
   │
   S (0V)              D (positive)
   │                    │
  [N+] ─── P-substrate ─── [N+]

No channel between S and D → no current → OFF state
```

### Step 2: VGS > 0, VGS < Vth (Depletion only)

- Positive gate voltage pushes holes away from surface (P-substrate)
- Creates depletion region under gate
- Still no electrons → still no channel

### Step 3: VGS > Vth (Inversion — channel forms!)

```
   G (+5V)
   │
   ├── Gate oxide ────────────────────────────
              ↑ Electric field
   n n n n n ← INVERSION LAYER (electrons attracted from N+ S/D and thermally generated)

   S (0V)                              D (+5V)
   │                                    │
  [N+] ── depletion ── inversion ── depletion ── [N+]
          region     channel(N)      region

Electrons flow from S to D (conventional current D to S)!
```

**Threshold voltage (Vth) for NMOS:**
```
Vth = Vfb + 2φF + Qd/Cox

Where:
  Vfb = flat-band voltage (work function difference + oxide charges)
  φF  = Fermi potential = (kT/q)×ln(NA/ni)
  Qd  = depletion charge density = q×NA×Xd_max
  Cox = gate oxide capacitance per unit area = ε_ox/t_ox

Typical Vth for NMOS: +0.5V to +3V
Typical Vth for PMOS: -0.5V to -3V
Logic-level MOSFETs: Vth < 2V (can be driven by 3.3V/5V logic)
```

---

## 4. MOSFET Equations and Regions

### Three Operating Regions

**Region 1: Cutoff**
```
Condition: VGS < Vth
Result:    ID = 0  (no channel, no current)
           Transistor is OFF
```

**Region 2: Linear (Ohmic / Triode)**
```
Condition: VGS > Vth  AND  VDS < (VGS - Vth) = Vov

ID = μn × Cox × (W/L) × [(VGS-Vth)×VDS - VDS²/2]

Where:
  μn   = electron mobility in channel (~270 cm²/V·s in inversion layer)
  Cox  = gate oxide capacitance = ε_ox/t_ox
  W    = channel width (perpendicular to current flow)
  L    = channel length (source to drain)
  Vov  = overdrive voltage = VGS - Vth
  VDS  = drain-source voltage

For small VDS: VDS² term negligible:
  ID ≈ μn×Cox×(W/L)×(VGS-Vth)×VDS

In this region: acts like voltage-controlled resistor (linear)
Ron = VDS/ID = 1/(μn×Cox×(W/L)×(VGS-Vth))
```

**Region 3: Saturation (Pinch-off)**
```
Condition: VGS > Vth  AND  VDS ≥ (VGS - Vth)

ID = (μn×Cox/2) × (W/L) × (VGS-Vth)²  = K×(VGS-Vth)²/2

Where K = μn×Cox×W/L = process transconductance × geometry

Including channel length modulation:
ID = (K/2)×(VGS-Vth)² × (1 + λ×VDS)

λ = channel length modulation parameter (~0.01-0.1 V⁻¹)
VA = 1/λ = Early voltage (analogous to BJT VA)

Key: ID is independent of VDS in saturation! (controlled only by VGS)
Used for: amplifiers, current mirrors, constant current sources
```

### Transconductance

```
gm = ∂ID/∂VGS|VDS=const = μn×Cox×(W/L)×(VGS-Vth) = √(2K×ID)

For voltage gain in common-source amplifier:
Av = -gm×RD (negative = 180° phase shift)
```

### W/L Ratio — Design Parameter

```
W/L ratio determines:
  - Current driving capability
  - On-resistance (RDS_on)
  - Speed (larger W = more parasitic capacitance)

Wider (larger W): more current, lower RDS_on, slower switching
Shorter (smaller L): more current, faster, but harder to manufacture
```

---

## 5. CMOS Technology

**CMOS = Complementary MOS** — pairs of NMOS + PMOS

### CMOS Inverter (NOT Gate)

```
            VDD
             │
            [PMOS] ← (gate connected to input Vin)
             │
             ├──────── Output (Vout)
             │
            [NMOS] ← (gate connected to input Vin)
             │
            GND

Operation:
  Vin = 0V (LOW):
    PMOS: VGS = 0-VDD = -VDD < Vtp (p-channel threshold) → PMOS ON
    NMOS: VGS = 0V < Vtn → NMOS OFF
    → Output connected to VDD → Vout = HIGH = VDD

  Vin = VDD (HIGH):
    PMOS: VGS = VDD-VDD = 0V > Vtp → PMOS OFF
    NMOS: VGS = VDD > Vtn → NMOS ON
    → Output connected to GND → Vout = LOW = 0V
```

**Why CMOS is revolutionary:**
1. **Near-zero static power:** when output is stable, either PMOS or NMOS is OFF → no DC path from VDD to GND → essentially zero static power!
2. **Full voltage swing:** output swings rail-to-rail (0V to VDD)
3. **High noise margin:** since output is VDD or GND (not somewhere in between)
4. **Scales beautifully:** smaller transistors → higher density → same power

### CMOS Dynamic Power Dissipation

Even though static power ≈ 0, switching has dynamic power:

```
P_dynamic = α × C × V² × f

Where:
  α = activity factor (0 to 1, fraction of clock cycles where node switches)
  C = load capacitance (sum of wire + transistor gate capacitances)
  V = supply voltage
  f = clock frequency

This is why:
  Lower voltage saves power quadratically (V² term!) — from 5V to 1V = 25× less power
  Higher frequency linearly increases power
  Smaller transistors: less C but higher f → net depends
  Modern chips: massive power from billions of transistors switching at GHz rates
```

### CMOS NAND Gate

```
           VDD
            │
     ┌──────┴──────┐
    [PA]          [PB]    ← PMOS in PARALLEL
     │              │
     └──────┬───────┘
            │
           Vout
            │
           [NA]           ← NMOS in SERIES
            │
           [NB]
            │
           GND

For both A=1, B=1:
  NA and NB both ON (series = ON)
  PA and PB both OFF (parallel = all OFF)
  Vout = 0 (LOW)

For A=0 or B=0:
  At least one of NA or NB is OFF → series path broken
  At least one of PA or PB is ON → Vout = VDD (HIGH)
→ NAND behavior!
```

**Rule for CMOS gate implementation:**
- Pull-up network (PUN): PMOS transistors — mirror of pull-down
  - NMOS in series → PMOS in parallel (NAND)
  - NMOS in parallel → PMOS in series (NOR)
- Pull-down network (PDN): NMOS transistors — follows Boolean function directly

---

## 6. MOSFET as a Switch

MOSFET is preferred over BJT for most modern switching applications.

### N-channel MOSFET Switch (Low-Side)

```
      VDD (load supply)
       │
      [Load] (motor, LED, relay, etc.)
       │
       Drain
   G ──┤ NMOS
       Source
       │
      GND

Control signal → Gate (via series resistor 10Ω-1kΩ to limit current spikes)
VGS > Vth → transistor ON → current flows through load
VGS < Vth → transistor OFF
```

**Gate drive considerations:**
- Gate capacitance Ciss must be charged to turn on (and discharged to turn off)
- For fast switching: gate driver IC provides high peak current (1-10A for power MOSFETs)
- Gate charge Qg = total charge needed to swing gate
- Turn-on time: ≈ Qg / Idriver

### P-channel MOSFET Switch (High-Side)

```
      VDD
       │
       Source
   G ──┤ PMOS  ← VGS must be negative: VG < VS = VDD
       Drain
       │
      [Load]
       │
      GND
```

For VGS = −10V on a P-channel: need VG = VDD − 10V
- If VDD = 12V: need VG = 2V → can drive from 5V logic with proper circuit

### Gate Driver ICs

For power MOSFETs, use gate driver ICs for fast, efficient switching:
- **TC4427/TC4428:** dual 1.5A gate driver, non-inverting/inverting
- **IR2110:** high-side and low-side driver, up to 500V
- **DRV8313:** full H-bridge driver for motors
- **TPS2828:** synchronous buck gate driver

---

## 7. Power MOSFETs

Designed for high current and/or high voltage switching applications.

### Structure

Power MOSFETs use a **DMOS (Double-diffused MOS)** structure:

```
                  Gate
                   │
    ┌─────────────────────────────┐
    │ Source (N+) many cells      │ Source metal on top
    │ P-body                      │ Gate poly
    │ N-drift (lightly doped)     │ High voltage held here
    │ N+ substrate                │
    │                             │ Drain (bottom contact)
    └─────────────────────────────┘

Current flows vertically (source top to drain bottom)
Many parallel cells → handles high current
N-drift region → handles high voltage
```

### Key Specifications

**RDS(on) — On-state resistance:**
```
Power dissipated when ON:
  P = ID² × RDS(on)

Typical values:
  IRF540N: 44mΩ (55V, 33A) — popular hobbyist MOSFET
  IRLZ44N: 28mΩ (55V, 47A) — logic-level version
  SiC MOSFET: very low RDS(on) at high voltage (Wolfspeed, STMicro SiC)
```

**Breakdown voltage BVDSS:**
- Must exceed maximum VDS in circuit (typically use 80% of rated)

**Gate threshold voltage Vth:**
- Logic-level MOSFETs: Vth < 2.5V (can be driven directly from 3.3V or 5V microcontroller)
- Standard MOSFETs: Vth up to 4V (need 10V+ gate drive from dedicated driver)

**Maximum current ID:**
- Limited by heating (I²×RDS_on losses)
- Junction temperature limit: 150-175°C

**Gate charge Qg:**
- Total charge to swing gate from 0 to VGS
- Qg = Qgs + Qgd + remaining Qgd
- Lower Qg → faster switching → less switching loss

### Popular Power MOSFETs

| Part | VDS(max) | ID(max) | RDS(on) | Logic Level? | Package |
|------|---------|---------|---------|-------------|---------|
| IRLZ44N | 55V | 47A | 28mΩ @ 10V | YES (5V) | TO-220 |
| IRF540N | 100V | 33A | 44mΩ | No (10V) | TO-220 |
| IRF3205 | 55V | 110A | 8mΩ | No | TO-220 |
| Si2302 | 20V | 3A | 100mΩ | YES | SOT-23 |
| AO3400 | 30V | 5.7A | 52mΩ | YES | SOT-23 |
| FDMS7672S | 100V | 20A | 11mΩ | No | PowerPAK |

### Freewheeling Diode

All power MOSFETs have an **internal body diode** (from drain to source in NMOS):
- Formed by the PN junction between P-body and N-drain
- Conducts when inductor current tries to reverse (freewheeling)
- Usually adequate for motor/inductor circuits
- For very fast switching: external Schottky diode in parallel (body diode is slow)

---

## 8. JFET (Junction Field Effect Transistor)

**Older than MOSFET**, uses reverse-biased PN junction to control channel.

### Structure

```
N-channel JFET:

    Gate (P+) ─────┐
                    │
    Source (N) ─── [N channel] ─── Drain (N)
                    │
    Gate (P+) ─────┘

Two gate regions (both connected) pinch the N channel
```

### Operation

- **VGS = 0V:** Maximum current IDSS flows (full channel open)
- **VGS negative:** P-gate reverse biased, depletion expands, narrows channel
- **VGS = VP (pinch-off voltage):** Channel completely pinched off, ID = 0

**Drain current equation (saturation):**
```
ID = IDSS × (1 - VGS/VP)²

Where:
  IDSS = saturation current at VGS = 0V (typically 1-50mA)
  VP   = pinch-off voltage (negative for N-channel, typically -1 to -8V)
```

**Key JFET properties:**
- Depletion mode (normally ON)
- Very high input impedance (gate is reverse-biased PN junction, essentially no current)
- Low noise (no channel oxide → no oxide noise)
- Used in: first stage of sensitive amplifiers (audio preamps, scientific instruments)

**Self-biasing (depletion mode):**
```
RS (source resistor) creates negative VGS automatically:
  VGS = -ID × RS
No separate bias supply needed!
```

**Popular JFETs:** 2N5457/5458 (N-channel audio), J201, MPF102, BF245

---

## 9. IGBT (Insulated Gate Bipolar Transistor)

An IGBT combines the best of MOSFET and BJT:
- **MOSFET input:** high input impedance gate, voltage-controlled, easy to drive
- **BJT output:** lower VCE(sat) than MOSFET for high voltage/current

### Structure

```
      Gate (G)                     Collector (C)
         │                              ↑
    ─────┴───────────────────────────────
    │ Gate oxide                        │
    │ N-channel (induced by gate)        │
    │ P-body                            │
    │ N-drift                           │
    │ P-substrate ← (this is what makes it bipolar)
    │                                   │
    ─────┬───────────────────────────────
         │
       Emitter (E)
```

The P-substrate injects holes into the drift region (like BJT action), reducing VCE(sat).

### Comparison to MOSFET for High Voltage

At high voltages (>600V), IGBT has lower VCE(sat) than equivalent MOSFET:
```
MOSFET: VDS(on) = ID × RDS(on) (rises with voltage rating — RDS scales as BV^2.5)
IGBT:   VCE(sat) ≈ 1.5-3V (relatively constant, from PN junction offset)

At 1200V, 100A:
  MOSFET: losses might be 100² × 0.05Ω = 500W (very high!)
  IGBT: losses ≈ 100 × 2V = 200W (much better)
```

### IGBT Drawbacks vs MOSFET
- **Slower switching:** minority carrier tail current (like BJT saturation)
- **No reverse blocking** (for that, use reverse-blocking IGBT or add diode)
- Can't switch as fast as MOSFET (max ~100kHz, vs MHz for MOSFET)

### IGBT Applications
- Motor drives (variable frequency drives, VFDs)
- EV/HEV traction inverters (Tesla Model 3 rear motor)
- UPS systems
- Induction heating
- Grid-tied solar inverters
- Welding machines
- Typical voltage range: 600V, 1200V, 1700V, 3300V, 6500V

### Popular IGBTs
- IHW40N60R3 (600V, 40A, TO-247)
- FGH30N60 (600V, 60A, TO-247)
- CM450DX-24S (1200V, 450A, IGBT module)

---

## 10. Advanced MOSFET Structures

As transistors shrink below ~100nm, classical planar MOSFET faces serious problems:
- **Short channel effects:** VT depends on channel length (bad!)
- **Drain-induced barrier lowering (DIBL):** drain voltage reduces threshold voltage
- **Subthreshold leakage:** significant current even when "OFF"
- **Gate oxide tunneling:** at <2nm oxide, electrons tunnel through — massive leakage!

Solutions:

### FinFET (3D Fin Field Effect Transistor)

```
Standard planar:           FinFET:
    Gate ─────                  Gate
   ┌────────────┐             ─────┐
   │  channel   │            │    │  ← gate wraps around 3 sides!
   └────────────┘            │ fin│
   S      D                  │    │
                              ─────┘
                              S   D

Fin = thin vertical silicon slab (~7-10nm wide, 30-40nm tall)
Gate wraps around 3 sides of fin
→ Much better electrostatic control
→ Greatly reduces short-channel effects
→ Lower leakage, steeper subthreshold slope
```

**History:**
- Intel introduced FinFET commercially at 22nm (2011) — called "Tri-Gate"
- TSMC/Samsung followed at 16/14nm (2014-2015)
- Used at: 22nm, 14nm, 10nm, 7nm, 5nm, 4nm

**Fin dimensions are quantized:** you can use 1 fin, 2 fins, 3 fins (not arbitrary width)
→ Current per transistor comes in discrete steps

### Gate-All-Around (GAA) / Nanosheet FET

```
     Gate surrounds ALL FOUR sides:

     ─────────── Gate ──────────
     │  ┌─────────────────┐     │
     │  │  nanosheet chan │     │ ← gate on all 4 sides
     │  └─────────────────┘     │
     ──────────────────────────
     │  ┌─────────────────┐     │
     │  │  nanosheet chan │     │ ← stacked nanosheets
     │  └─────────────────┘     │
     ───────────────────────────
     S                           D
```

- Gate wraps all 4 sides = ultimate gate control
- Stacked nanosheets: multiple channels vertically = more current per footprint
- **Samsung**: introduced at 3nm in 2022 (called MBCFET — Multi-Bridge Channel FET)
- **TSMC**: N2 (2nm, 2025) uses nanosheets
- **Intel 20A/18A**: uses RibbonFET (GAA)

### FDSOI (Fully Depleted Silicon on Insulator)

```
      Gate
       │
   ────┴──── Gate oxide
   [channel] ← ultra-thin Si (5-10nm)
   [buried oxide (BOX)] ← insulates channel completely
   [Si substrate]
```

- Thin Si layer completely depleted → better control
- Buried oxide eliminates bulk leakage
- Back-gate (substrate) can fine-tune Vth dynamically — allows saving power
- Used by STMicroelectronics (28nm FD-SOI for IoT), GlobalFoundries (22FDX)
- Excellent for low-power wireless chips (BLE, NB-IoT)

### High-k / Metal Gate (HKMG)

Problem: SiO₂ gate oxide at <1.2nm leaks electrons by quantum tunneling.

Solution: Use **high-k dielectric** (higher permittivity = physically thicker oxide for same capacitance):

```
Same Cox = ε_ox/t_ox:
  SiO₂ (k=3.9):  t = 1nm → severe tunneling!
  HfO₂ (k=25):   t = 6.5nm → same C, much less tunneling!
```

**High-k materials:** HfO₂ (Hafnium dioxide), HfSiO, HfON, ZrO₂

**Metal gate:** Polysilicon gate replaced with metal (TiN, TaN, W) to eliminate:
- Polysilicon depletion (effective oxide thickness increase)
- Boron penetration from P+ poly gate into channel

Intel introduced HKMG at 45nm (2007) — all chips 45nm and below use this.

---

## 11. Transistor Comparison

| Property | NMOS | PMOS | NPN BJT | JFET | IGBT |
|----------|------|------|---------|------|------|
| Control | Voltage (VGS) | Voltage (VGS) | Current (IB) | Voltage (VGS) | Voltage (VGE) |
| Gate current | ~0 | ~0 | IB = IC/β | ~0 | ~0 |
| Input Z | >10¹⁴Ω | >10¹⁴Ω | ~kΩ | >10¹⁰Ω | >10¹⁴Ω |
| Carriers | Electrons | Holes | Both | Electrons | Both |
| Speed | Very fast | Fast | Medium | Fast | Slow |
| Max voltage | ~1200V | ~600V | ~1500V | ~60V | ~6500V |
| Max current | ~1000A (mod) | ~500A | ~200A | ~1A | ~3600A (mod) |
| Low VCE/VDS | Good (RDS_on) | Poorer | 0.7V offset | Good | Good |
| Scalability | Excellent | Good | Poor | Poor | N/A |
| Use in digital | YES (mainly) | YES (CMOS) | Older circuits | No | No |

---

## 12. Transistor Numbering Systems

### JEDEC (USA)
- Format: 2Nxxxx
- "2N" means 2-junction device (transistor)
- Examples: 2N2222, 2N3904, 2N3055, 2N7000 (MOSFET)

### Pro-Electron (European)
- Format: Two letters + digits
- First letter: semiconductor material (B=Si, A=Ge, C=GaAs)
- Second letter: type (C=audio, D=power audio, F=RF small signal, P=power, T=switching)
- Examples: BC547, BD135, BF245, BC817

### Japanese (JIS)
- Format: 2SAxxxx, 2SBxxxx, 2SCxxxx, 2SDxxxx
- 2SA: PNP high frequency
- 2SB: PNP audio
- 2SC: NPN high frequency/audio
- 2SD: NPN power
- Examples: 2SC945, 2SA1015, 2SC1815, 2SD1047

### SMD Codes
- Component too small for full part number → short code stamped on it
- Cross-reference table needed to find full part number
- Example: "1P" on SOT-23 = BC846 (NPN), "A1" = BC807 (PNP)
- Use manufacturer's SMD codebook or online tools

---

## 13. Summary

```
MOSFET FUNDAMENTALS
═══════════════════

Structure: G-S-D + Body
Gate insulated by oxide → no gate current
Voltage-controlled current source

Regions of operation:
  Cutoff:     VGS < Vth → ID = 0
  Linear:     VGS > Vth, VDS < Vov → ID = K[(VGS-Vth)VDS - VDS²/2]
  Saturation: VGS > Vth, VDS ≥ Vov → ID = (K/2)(VGS-Vth)²

Threshold voltage: Vth ≈ +1-3V (NMOS), -1 to -3V (PMOS)
Transconductance: gm = K×(VGS-Vth) = √(2K×ID)

CMOS Logic:
  NMOS + PMOS complementary pairs
  Static power ≈ 0 (no DC path VDD to GND)
  Dynamic power = α×C×V²×f

As switch:
  VGS > Vth: ON, RDS(on) = 1/(K×Vov)
  VGS < Vth: OFF, near-zero current
  Advantages over BJT: no gate current, faster, voltage controlled

Advanced:
  FinFET (14nm-5nm): gate wraps 3 sides of silicon fin
  GAA/Nanosheet (3nm-2nm): gate wraps all 4 sides
  High-k dielectric: HfO₂ replaces SiO₂ (45nm and below)
```

### Key Formulas

| Formula | Meaning |
|---------|---------|
| ID = (K/2)(VGS-Vth)² | Saturation drain current |
| K = μn×Cox×W/L | Process + geometry transconductance |
| gm = √(2K×ID) | Small-signal transconductance |
| Av = -gm×RD | Voltage gain (common source) |
| P = α×C×V²×f | CMOS dynamic power |
| RDS(on) = 1/(K×Vov) | On-resistance |

---

**← Previous:** [Chapter 06: BJT Transistors](./06-bjt-transistors.md)
**→ Next:** [Chapter 08: Logic Gates and Boolean Algebra](./08-logic-gates-and-boolean-algebra.md)
