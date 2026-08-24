# Chapter 05: Basic Electronic Components — Resistors, Capacitors, Inductors

> **"Every electronic circuit is built from just a handful of fundamental components. Master these and you can analyze any circuit in the world."**

---

## Table of Contents
1. [Resistors](#1-resistors)
2. [Capacitors](#2-capacitors)
3. [Inductors](#3-inductors)
4. [Transformers](#4-transformers)
5. [LC Circuits and Resonance](#5-lc-circuits-and-resonance)
6. [AC Circuit Analysis](#6-ac-circuit-analysis)
7. [Filters](#7-filters)
8. [Summary of Formulas](#8-summary-of-formulas)

---

## 1. Resistors

### What is Resistance?

A **resistor** is a component that opposes the flow of current. Energy is dissipated as heat.

**Ohm's Law:** V = I × R

**Resistivity formula:**
```
        ρ × L
R = ───────────
          A

Where:
  ρ = resistivity (Ω·m) — material property
  L = length (m)
  A = cross-sectional area (m²)
```

### Types of Resistors

#### 1. Carbon Composition
- Ground carbon mixed with binder, pressed into rod
- Poor tolerance: ±5-20%
- High noise, non-linear at high frequencies
- Good surge capability
- Used in: vintage circuits, spark gaps, some RF

#### 2. Carbon Film
- Thin carbon film on ceramic core, spiral groove cut for value
- Tolerance: ±5% typical
- Lower noise than carbon composition
- Common, cheap
- Used in: general purpose

#### 3. Metal Film (Most Common for Precision)
- Metal alloy (NiCr, MnCu) film on ceramic
- Tolerance: ±0.1% to ±1%
- Low noise, excellent temperature stability
- TCR (Temperature Coefficient of Resistance): ±25 to ±100 ppm/°C
- Used in: precision circuits, audio, instrumentation

#### 4. Metal Oxide Film
- Tin oxide on ceramic
- Better power handling than metal film
- Tolerance: ±5%
- Used in: high voltage, higher temperature applications

#### 5. Wire Wound
- Resistance wire (NiCr) wound on ceramic core
- Very accurate: ±0.01%
- High power: 1W to thousands of watts
- But: inductance! (coil shape) — not for high-frequency use
- Used in: power circuits, precision decade boxes, current shunts

#### 6. SMD Chip Resistors (Surface Mount)
- Thick film or thin film on ceramic substrate
- Package sizes:
  ```
  0201: 0.6mm × 0.3mm  (very tiny, phone boards)
  0402: 1.0mm × 0.5mm  (small)
  0603: 1.6mm × 0.8mm  (most common)
  0805: 2.0mm × 1.25mm (easy to hand-solder)
  1206: 3.2mm × 1.6mm  (hand-solderable, higher power)
  2010: 5.0mm × 2.5mm  (high power)
  2512: 6.4mm × 3.2mm  (high power)
  ```
- Tolerance: ±1% (standard), ±0.1% (precision)
- 3-digit or 4-digit code for value (or EIA-96 code)
  ```
  3-digit: 104 = 10 × 10⁴ = 100,000Ω = 100kΩ
  4-digit: 1002 = 100 × 10² = 10,000Ω = 10kΩ
  EIA-96: 01A = 100Ω, uses letter for multiplier
  ```

#### 7. Variable Resistors
- **Potentiometer (pot):** 3-terminal, full resistance between ends, wiper taps off a fraction
  - Linear taper: resistance changes linearly with rotation
  - Audio (log) taper: logarithmic — matches human hearing
  - Used for: volume control, position sensing

- **Rheostat:** 2-terminal, variable series resistance
  - Used for: current control, motor speed (basic)

- **Trimmer/Trimpot:** small preset potentiometer for one-time adjustment
  - Used for: calibration, circuit tuning

- **Digital Potentiometer:** IC with digitally controlled wiper
  - Examples: MCP4131, AD5242
  - Used for: digitally-controlled audio, calibration

#### 8. Thermistor

Resistance varies strongly with temperature.

**NTC (Negative Temperature Coefficient):**
- Resistance DECREASES as temperature increases
- Made of metal oxide ceramics
- Used in: temperature measurement, inrush current limiting, temperature compensation
- Steinhart-Hart equation:
  ```
  1/T = A + B×ln(R) + C×(ln(R))³

  Simplified β (beta) equation:
  R(T) = R₀ × exp(β × (1/T - 1/T₀))

  Where:
    R₀ = resistance at reference temperature T₀ (usually 25°C = 298K)
    β  = material constant (typically 3000-5000K for NTC)
    T  = temperature in Kelvin
  ```
- **10kΩ NTC at 25°C** is the most common (used in 3D printers, thermostats)
- Reading temperature from NTC:
  ```
  Temperature (°C) from voltage divider reading:

  R_NTC = R_fixed × (VCC/V_out - 1)  (if NTC is bottom resistor)

  T = β / ln(R_NTC/R₀) + 1/T₀  (in Kelvin, convert to Celsius: T-273.15)
  ```

**PTC (Positive Temperature Coefficient):**
- Resistance INCREASES with temperature
- Used in: self-resetting fuses (polyfuses), heaters
- Polyfuse: trips (high resistance) when overloaded, resets when cooled

#### 9. LDR (Light Dependent Resistor / Photoresistor)

- CdS (Cadmium Sulfide) or CdSe sensitive material
- Resistance drops dramatically in light
- Dark resistance: 1MΩ to 10MΩ
- Bright light resistance: 100Ω to 10kΩ
- Slow response: ~100ms (not for fast switching)
- Used in: automatic street lights, camera light metering

### Resistor Color Code

**4-Band Resistor:**
```
Band 1: First significant digit
Band 2: Second significant digit
Band 3: Multiplier (×10^n)
Band 4: Tolerance (gold=±5%, silver=±10%)

Color → Digit table:
Black  = 0
Brown  = 1
Red    = 2
Orange = 3
Yellow = 4
Green  = 5
Blue   = 6
Violet = 7
Gray   = 8
White  = 9
Gold   = × 0.1 (multiplier) / ±5% (tolerance)
Silver = × 0.01 (multiplier) / ±10% (tolerance)

Example: Red-Red-Orange-Gold
  = 2-2-×1000-±5%
  = 22,000Ω ±5% = 22kΩ ±5%

Example: Brown-Black-Black-Gold
  = 1-0-×1-±5%
  = 10Ω ±5%
```

**5-Band Resistor (higher precision):**
```
Band 1,2,3: Three significant digits
Band 4: Multiplier
Band 5: Tolerance (brown=±1%, red=±2%, green=±0.5%)

Example: Brown-Black-Black-Black-Brown
  = 1-0-0-×1-±1%
  = 100Ω ±1%
```

### Series and Parallel Resistors

**Series:**
```
R_total = R₁ + R₂ + R₃ + ... + Rn

Current same everywhere: I₁ = I₂ = I₃ = I
Voltage divides: V₁ = I×R₁, V₂ = I×R₂, etc.

Voltage divider:
       R₂
Vo = ───── × Vin
     R₁+R₂
```

**Parallel:**
```
1/R_total = 1/R₁ + 1/R₂ + 1/R₃

For 2 resistors:
         R₁ × R₂
R_total = ─────────
          R₁ + R₂

Voltage same: V₁ = V₂ = V
Current divides: I₁ = V/R₁, I₂ = V/R₂

Current divider (from source I into R₁ and R₂ in parallel):
       R₂
I₁ = ───── × I  ← current through R₁ (larger R₁ gets less current!)
     R₁+R₂
```

### Resistor Power Rating

**Power dissipated in resistor:**
```
P = I² × R = V²/R = V × I

Standard power ratings: 1/8W, 1/4W (most common), 1/2W, 1W, 2W, 5W, 10W...
```

**Rule of thumb:** Always use resistors at < 50% of rated power for reliability.

**Example:** 100Ω resistor with 5V across it:
```
P = V²/R = 25/100 = 0.25W = 250mW → need at least 1/2W rated resistor
```

### E-Series Standard Values

Resistors come in standard values (E-series):
- **E12 series (±10%):** 10, 12, 15, 18, 22, 27, 33, 39, 47, 56, 68, 82 (and ×10ⁿ multiples)
- **E24 series (±5%):** 24 values per decade
- **E48 series (±2%):** 48 values per decade
- **E96 series (±1%):** 96 values per decade
- **E192 series (±0.5%):** 192 values per decade

---

## 2. Capacitors

### What is Capacitance?

A **capacitor** stores electrical energy in an **electric field** between two conductors separated by an insulator (dielectric).

**Charge stored:**
```
Q = C × V

Where:
  Q = charge (Coulombs)
  C = capacitance (Farads)
  V = voltage (Volts)
```

**Capacitance of parallel plate capacitor:**
```
         ε₀ × εr × A
C = ─────────────────
             d

Where:
  ε₀ = 8.854×10⁻¹² F/m (permittivity of free space)
  εr = relative permittivity (dielectric constant) of material
  A  = plate area (m²)
  d  = separation between plates (m)

Higher εr → larger capacitance:
  Air: εr = 1
  Glass: εr = 4-10
  Ceramic (X7R): εr = 2000-5000
  Aluminum oxide: εr = 8-10
  Tantalum oxide: εr = 27
  Barium titanate (BaTiO₃): εr > 10,000
```

**Energy stored:**
```
E = ½ × C × V²

Units: Joules
```

### Capacitor Charging and Discharging

**Charging (RC circuit, capacitor starts at 0V):**
```
v(t) = VS × (1 - e^(-t/RC))
i(t) = (VS/R) × e^(-t/RC)

Time constant: τ = R × C

Timeline:
  t = 1τ: v = 63.2% of VS
  t = 2τ: v = 86.5% of VS
  t = 3τ: v = 95.0% of VS
  t = 4τ: v = 98.2% of VS
  t = 5τ: v = 99.3% of VS  ← considered fully charged
```

**Discharging (from V₀ through R to GND):**
```
v(t) = V₀ × e^(-t/RC)
i(t) = (V₀/R) × e^(-t/RC)

At t = τ: v = 36.8% of V₀ (decayed to 1/e)
At t = 5τ: v = 0.7% of V₀ (considered fully discharged)
```

**Example:** RC timer (used in 555 timer circuit)
```
τ = R × C = 10kΩ × 100μF = 1 second
At t=1s: capacitor charged to 63.2% of supply voltage
```

### Capacitor Types

#### 1. Ceramic (MLCC — Multi-Layer Ceramic Capacitor)
- **Most common type** on PCBs
- Made of alternating ceramic dielectric and metal layers
- Non-polarized (can be connected either way)
- Very wide range: 1pF to 100μF
- Temperature characteristics (dielectric types):
  ```
  NP0/C0G: ±30ppm/°C, very stable, <1nF, precision timing/filters
  X5R: ±15%, -55°C to +85°C, up to 10μF, decoupling
  X7R: ±15%, -55°C to +125°C, up to 4.7μF, decoupling
  Y5V: +22%/-82%, -30°C to +85°C, cheapest, for non-critical bypass
  ```
- High-frequency performance: excellent (no inductance issues, very low ESR)
- **Decoupling capacitors** next to ICs: 100nF X7R ceramic (standard)
- Note: some MLCCs are **microphonic** (generate voltage from vibration)

#### 2. Electrolytic (Aluminum)
- **Highest capacitance per cost**
- **Polarized** — must connect + terminal to positive voltage!
- Made of aluminum oxide dielectric (formed electrolytically)
- Large: 1μF to 100,000μF
- High ESR (Equivalent Series Resistance) — affects filtering
- Limited frequency range: good up to ~100kHz
- **Aging:** capacitance decreases over time (~20% per decade of years)
- Temperature limit: typically 85°C or 105°C rated
- **Failure modes:** venting (splits top), exploding (backward connection!), dried electrolyte
- Used in: power supply filtering, audio coupling

**Low-ESR electrolytics:** for switching power supplies (higher ripple current capability)

#### 3. Tantalum
- **Polarized** (very important — can fail violently if reversed!)
- Tantalum pentoxide dielectric
- Smaller than equivalent aluminum electrolytic
- Better high-frequency performance (lower ESR)
- More stable with temperature
- Range: 0.1μF to 1000μF
- Maximum voltage rating usually lower than equivalent aluminum
- **Failure mode:** short circuit (can catch fire!) — always use at <50-70% of rated voltage
- Used in: portable electronics, anywhere small size important

#### 4. Film Capacitors
- Plastic film dielectric (polyester, polypropylene, polystyrene)
- Non-polarized
- No aging
- Low leakage
- Types:
  - **Polyester (PET/Mylar):** cheap, general purpose, 1nF-10μF
  - **Polypropylene (PP):** very low losses, excellent for audio and timing, up to 1μF
  - **Polystyrene:** excellent precision and stability, up to 10nF — being discontinued
  - **PTFE (Teflon):** extreme precision, RF applications
- Used in: audio circuits, timing, power factor correction (high voltage)

#### 5. Supercapacitor (Ultracapacitor, EDLC)
- Uses electric double-layer at electrode/electrolyte interface
- Capacitance: 1F to 3000F (!) — compared to mF for electrolytic
- But low voltage: max 2.7V per cell (series stacks used for higher voltage)
- Energy density: 1-30 Wh/kg (between battery and capacitor)
- Power density: very high — can charge/discharge in seconds
- Essentially unlimited cycle life (no chemical reaction)
- Used in: energy storage for regenerative braking, UPS backup, burst power

#### 6. Variable Capacitors
- Air-dielectric, mechanically adjustable
- Used in: AM radio tuning (historic), antenna matching
- Trimmer capacitor: small preset, used for one-time tuning

### Capacitors in Circuits

**Series:**
```
1/C_total = 1/C₁ + 1/C₂ + 1/C₃
(Total less than smallest — rarely used intentionally)
Voltage divides across capacitors
```

**Parallel:**
```
C_total = C₁ + C₂ + C₃
(Values add — used to increase capacitance)
Same voltage across all
```

**Capacitive reactance (AC opposition):**
```
         1
Xc = ──────────
     2π × f × C

At low frequency: Xc is HIGH (blocks low freq and DC)
At high frequency: Xc is LOW (passes high freq)
→ Capacitor passes AC, blocks DC

Impedance: Z = 1/(jωC) or Z = Xc ∠-90°
Phase: current LEADS voltage by 90°
```

### Capacitor Applications

1. **Decoupling/Bypass:** 100nF ceramic across every IC power pin — absorbs high-frequency current spikes, prevents them from spreading
2. **Bulk filtering:** large electrolytic (1000μF+) across power supply output — reduces low-frequency ripple
3. **Coupling/AC coupling:** blocking DC, passing AC signal (e.g., audio stages)
4. **Timing:** RC networks for timers, oscillators (555 timer)
5. **Sampling and hold:** holds voltage for ADC conversion
6. **Motor start capacitor:** phase shift for single-phase AC motor starting
7. **Power factor correction:** in AC power systems
8. **Snubber:** R-C across relay contacts or transistor to suppress spikes

---

## 3. Inductors

### What is Inductance?

An **inductor** stores energy in a **magnetic field** created by current flowing through a coil of wire.

**Faraday's Law:**
```
V = L × di/dt

The inductor opposes CHANGES in current
When current changes: inductor generates voltage to oppose that change

Units of inductance: Henry (H), mH, μH, nH
```

**Energy stored:**
```
E = ½ × L × I²

Units: Joules
```

**Inductance of a solenoid coil:**
```
         μ₀ × μr × N² × A
L = ──────────────────────
                l

Where:
  μ₀ = 4π×10⁻⁷ H/m (permeability of free space)
  μr = relative permeability of core material
  N  = number of turns
  A  = cross-sectional area of core (m²)
  l  = length of coil (m)

μr values:
  Air:        1
  Iron:       200-5000
  Ferrite:    100-10,000
  Permalloy:  100,000+
```

### RL Circuit Transient Response

**Building up current (switch closed at t=0):**
```
i(t) = (V/R) × (1 - e^(-t/τ))

Where τ = L/R (time constant in seconds)

At t = 1τ: i = 63.2% of final V/R
At t = 5τ: i = 99.3% of V/R (fully energized)
```

**Collapsing current (switch opened):**
```
Inductor tries to maintain current → generates very high voltage spike!
This "inductive kickback" can destroy transistors if not clamped.
Solution: freewheeling diode in parallel with inductor

v_spike = L × di/dt (can be thousands of volts if di/dt is large)
```

### Inductive Reactance

```
XL = 2π × f × L = ωL

At low frequency: XL is LOW (passes DC and low freq)
At high frequency: XL is HIGH (blocks high freq)
→ Opposite of capacitor!

Impedance: Z = jωL or Z = XL ∠+90°
Phase: current LAGS voltage by 90°
```

### Types of Inductors

#### Air Core
- No magnetic material (just coils of wire)
- μr = 1 → low inductance per turn
- Very low losses (no core loss)
- Stable — inductance doesn't change with frequency or current
- Used in: RF circuits (10nH to 10μH), RF filters, antenna matching

#### Iron Core
- Iron/silicon steel laminations (to reduce eddy currents)
- High μr → very high inductance in small size
- But: saturates at high current, loses at high frequency
- Used in: power transformers, audio chokes (up to 60Hz power)

#### Ferrite Core
- Iron oxide + metal oxide ceramics (MnZn or NiZn ferrite)
- High μr, low eddy current losses up to MHz range
- Different ferrites for different frequency ranges:
  - MnZn: good to ~1MHz (power supplies, EMI filters)
  - NiZn: good to >100MHz (RF suppression)
- Used in: switching power supply inductors, EMI chokes, ferrite beads, transformers

#### Toroidal Core
- Donut-shaped core
- Flux contained within core — minimal EMI radiation
- More efficient than rod core
- Used in: audio transformers, power inductors in SMPS

#### Ferrite Beads
- Not really an inductor — a lossy element
- Converts high-frequency energy to heat
- Acts as low-pass filter with frequency
- Common values: 600Ω at 100MHz impedance
- Placed in series on power/signal lines to suppress EMI
- Used on: USB lines, power rails, I/O lines in EMC-sensitive designs

### Inductors in Circuits

**Series:**
```
L_total = L₁ + L₂ + L₃   (if no mutual coupling)
```

**Parallel:**
```
1/L_total = 1/L₁ + 1/L₂ + 1/L₃
```

### Key Inductor Specifications

- **Inductance (L):** value in μH or mH
- **DC Resistance (DCR):** wire resistance, causes power loss
- **Saturation current (Isat):** current at which L drops 20-30% — must never exceed this!
- **RMS current rating:** heating limit
- **Self-Resonant Frequency (SRF):** above SRF, acts as capacitor — use below SRF

---

## 4. Transformers

A **transformer** transfers AC power between circuits through magnetic coupling, allowing voltage step-up or step-down with corresponding current change.

**Turns ratio:**
```
Vs   Ns   Ip
── = ── = ──
Vp   Np   Is

Where:
  Vp, Vs = primary, secondary voltage
  Np, Ns = primary, secondary turns
  Ip, Is = primary, secondary current

Power conservation (ideal):
  Pp = Ps → Vp × Ip = Vs × Is
```

**Types:**
- **Step-down:** Ns < Np → lower secondary voltage, higher current
- **Step-up:** Ns > Np → higher secondary voltage, lower current
- **Isolation:** Np = Ns, same voltage but galvanic isolation (safety)
- **Center-tapped:** secondary has center connection, gives ±Vs/2 or full-wave rectification

**Efficiency:** Real transformers have losses:
```
η = Pout/Pin = 1 - (copper losses + iron losses)/Pin

Power transformers: 95-99% efficient
SMPS transformers: 85-98% efficient
```

---

## 5. LC Circuits and Resonance

### Resonant Frequency

When a capacitor and inductor are combined, energy oscillates between electric field (capacitor) and magnetic field (inductor).

**Natural resonant frequency:**
```
           1
f₀ = ─────────────
      2π × √(L×C)

Or in angular frequency:
ω₀ = 1/√(LC)   radians/second
```

**Example:** L = 10μH, C = 100pF:
```
f₀ = 1/(2π × √(10×10⁻⁶ × 100×10⁻¹²))
   = 1/(2π × √(10⁻¹⁵))
   = 1/(2π × 31.6×10⁻⁹)
   = 5.03 MHz
```

### Series RLC Circuit

At resonance:
- XL = Xc → impedance is minimum = R
- Current is maximum
- Voltage across L and C can be **Q times** the source voltage!

**Quality factor Q:**
```
Q = f₀/BW = ω₀×L/R = 1/(ω₀×C×R) = (1/R)×√(L/C)

Bandwidth: BW = f₀/Q

High Q → sharp resonance, narrow bandwidth
Low Q → broad resonance, wide bandwidth
```

### Parallel RLC Circuit

At resonance:
- Impedance is maximum = R×Q² (or just R in ideal)
- Current from source is minimum
- Large circulating current within L and C

### Applications of LC Resonance

1. **Band-pass filters:** select frequency range (RF receivers, IF filters)
2. **Oscillators:** LC tank circuit sustains oscillation at f₀
3. **Impedance matching:** match antenna to transmitter
4. **Crystal oscillators:** quartz crystal as very high-Q LC equivalent
   - Q = 10,000 to 1,000,000 — incredibly stable frequency reference

---

## 6. AC Circuit Analysis

### Phasor Notation

For sinusoidal circuits at one frequency, use **phasors** (complex numbers) instead of time-domain equations.

```
V = Vm∠φ = Vm(cosφ + j×sinφ)

j = √(-1) (imaginary unit, called 'j' in engineering not 'i')
```

**Impedances in phasor form:**
```
Resistor:  Z_R = R         (purely real — no phase shift)
Inductor:  Z_L = jωL       (purely imaginary, +90°)
Capacitor: Z_C = 1/(jωC)   (purely imaginary, -90°)
```

**Series RLC impedance:**
```
Z = R + j(ωL - 1/ωC) = R + j(XL - Xc)

|Z| = √(R² + (XL-Xc)²)
φ = arctan((XL-Xc)/R)
```

### Power in AC Circuits

```
Real power (active):    P = V×I×cos(φ) = I²×R    [Watts, W]
Reactive power:         Q = V×I×sin(φ) = I²×X    [VAR, volt-ampere reactive]
Apparent power:         S = V×I                   [VA, volt-ampere]
Power triangle:         S² = P² + Q²
Power factor:           PF = cos(φ) = P/S
```

**Power factor correction:**
- Inductive loads (motors, transformers) have lagging PF (φ > 0, Q > 0)
- Adding capacitors cancels reactive power: improves PF toward 1
- Industrial utility companies charge penalties for poor PF

---

## 7. Filters

Filters select which frequencies pass and which are blocked.

### RC Low-Pass Filter

```
Input ──[ R ]──┬── Output
               [C]
                │
               GND
```

**Transfer function:**
```
H(jω) = 1/(1 + jωRC)

Cutoff frequency (−3dB point):
f_c = 1/(2πRC)

At f << fc: |H| ≈ 1 (passes)
At f = fc:  |H| = 1/√2 = 0.707 (−3dB)
At f >> fc: |H| ≈ fc/f (blocks, -20dB/decade rolloff)
```

**Example:** RC LPF for audio anti-aliasing before ADC
```
f_c = 20kHz, C = 10nF
R = 1/(2π × 20000 × 10×10⁻⁹) = 796Ω → use 820Ω
```

### RC High-Pass Filter

```
Input ──[C]──┬── Output
             [R]
              │
             GND
```

**Cutoff frequency:** same formula f_c = 1/(2πRC)
- Blocks DC and low frequencies
- Passes high frequencies
- -20dB/decade rolloff below cutoff

### LC Filters (Better Roll-Off)

- 2nd order: −40dB/decade (much sharper than RC's −20dB/decade)
- Butterworth: maximally flat passband
- Chebyshev: sharper rolloff, ripple in passband
- Bessel: maximally flat group delay (best for pulse signals)

### Filter Order

Each LC stage (or op-amp section) adds one "pole":
- 1st order: −20dB/decade
- 2nd order: −40dB/decade
- 3rd order: −60dB/decade
- nth order: −20n dB/decade

---

## 8. Summary of Formulas

### Resistors
| Formula | Meaning |
|---------|---------|
| V = IR | Ohm's Law |
| R = ρL/A | Resistance from geometry |
| Rtotal = R₁+R₂+... | Series combination |
| 1/Rtotal = Σ(1/Ri) | Parallel combination |
| Vout = Vin×R₂/(R₁+R₂) | Voltage divider |
| P = V²/R = I²R = VI | Power dissipation |

### Capacitors
| Formula | Meaning |
|---------|---------|
| Q = CV | Charge stored |
| C = ε₀εrA/d | Parallel plate capacitance |
| E = ½CV² | Energy stored |
| τ = RC | Time constant |
| v(t) = Vs(1-e^(-t/τ)) | Charging voltage |
| Xc = 1/(2πfC) | Capacitive reactance |
| Ctotal = C₁+C₂+... | Parallel combination |
| 1/Ctotal = Σ(1/Ci) | Series combination |

### Inductors
| Formula | Meaning |
|---------|---------|
| V = L·di/dt | Inductor voltage |
| E = ½LI² | Energy stored |
| τ = L/R | RL time constant |
| XL = 2πfL | Inductive reactance |
| Ltotal = L₁+L₂+... | Series combination |

### Resonance
| Formula | Meaning |
|---------|---------|
| f₀ = 1/(2π√LC) | Resonant frequency |
| Q = f₀/BW | Quality factor |
| BW = f₀/Q | Bandwidth |
| fc = 1/(2πRC) | RC filter cutoff |

### Transformer
| Formula | Meaning |
|---------|---------|
| Vs/Vp = Ns/Np = Ip/Is | Turns ratio |
| Pp = Ps (ideal) | Power conservation |

---

**← Previous:** [Chapter 04: PN Junction and Diodes](./04-pn-junction-and-diodes.md)
**→ Next:** [Chapter 06: BJT Transistors](./06-bjt-transistors.md)
