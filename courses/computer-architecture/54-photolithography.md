# Chapter 54: Photolithography — Printing Circuits at Nanoscale

Photolithography is the process by which circuit patterns are transferred from a design file onto a silicon wafer using light. It is the most critical, most expensive, and most technically challenging step in chip manufacturing. Every transistor gate, every metal wire, every contact must be patterned with sub-nanometer precision across an entire 300mm wafer — millions of times per chip. Understanding photolithography explains why chip manufacturing is so capital-intensive, why ASML is one of the most strategically important companies in the world, and why shrinking chips beyond 2nm requires fundamentally new approaches.

## Table of Contents

1. [How Photolithography Works](#1-how-photolithography-works)
2. [Resolution Limits: The Rayleigh Criterion](#2-resolution-limits-the-rayleigh-criterion)
3. [Immersion Lithography and ArF Scanners](#3-immersion-lithography-and-arf-scanners)
4. [EUV Lithography — The Breakthrough](#4-euv-lithography--the-breakthrough)
5. [Resolution Enhancement Techniques](#5-resolution-enhancement-techniques)
6. [Multi-Patterning — Cheating the Diffraction Limit](#6-multi-patterning--cheating-the-diffraction-limit)
7. [ASML: The Monopoly That Controls Chips](#7-asml-the-monopoly-that-controls-chips)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. How Photolithography Works

The basic photolithography step consists of five sub-steps:

**Step 1: Coat with photoresist**
Spin a thin layer (~50–150nm) of photoresist (light-sensitive polymer) onto the wafer at 3000 RPM. Surface tension creates a uniform film.

**Step 2: Soft bake**
Heat to ~90°C to drive off solvent and harden the resist.

**Step 3: Expose**
Shine UV/EUV light through a **reticle** (chrome-on-quartz glass plate carrying the circuit pattern at 4× magnification). The lens demagnifies 4× → 1× and projects the pattern onto the resist.

**Step 4: Post-exposure bake + Develop**
Heat activates acid in the resist (for chemically amplified resist). Dip in developer solution — either exposed regions dissolve (positive-tone resist) or unexposed regions dissolve (negative-tone).

**Step 5: Etch / Implant / Deposit**
The photoresist pattern now acts as a mask for the next process step (etch the layer below, implant dopants, or deposit material). After the process step, the resist is stripped.

```
Photolithography sequence:
  
  Wafer with            Coat resist:           Expose through mask:
  oxide layer           ▓▓▓▓▓▓▓▓▓▓             ▓▓ ▓▓▓▓▓ ▓▓
  ──────────            ──────────              ──────────
  (SiO₂)               (SiO₂)                 (exposed regions marked)
  
  Develop:              Etch:                  Strip resist:
  ▓▓▓   ▓▓▓             ──  ──────  ──         ── ──────── ──
  ──────────            (SiO₂ removed          (pattern in SiO₂)
  (resist removed        under gaps)
   where exposed)
  
  This forms the circuit feature (in this case, gaps in the oxide)
```

**Reticle (photomask)**: A 152mm × 152mm × 6.35mm fused silica plate with chrome patterns. For ArF lithography, absorber patterns define where light is blocked/transmitted. For EUV, it is a reflective multilayer mirror with an absorber pattern (EUV is absorbed by glass, not transmitted). A single chip design requires 40–80 reticles (one per process layer). A full reticle set costs $5–10 million.

### Quick Check
> 1. What is a reticle and what role does it play in photolithography?
> 2. What is the difference between positive-tone and negative-tone resist?
> 3. Why is the reticle pattern 4× larger than the final features on the wafer?

---

## 2. Resolution Limits: The Rayleigh Criterion

How small a feature can be printed is governed by the **Rayleigh criterion** from optics:

```
Resolution: R = k₁ × λ / NA

Depth of focus: DOF = k₂ × λ / NA²

Where:
  λ = wavelength of light used
  NA = Numerical Aperture of the lens (= n × sin θ, where n = refractive index of medium)
  k₁ = process factor (0.25 < k₁ < 0.8; lower = harder, requires more tricks)
  k₂ = depth of focus factor
```

To print smaller features, you can:
1. **Reduce λ** (shorter wavelength light)
2. **Increase NA** (larger lens aperture)
3. **Reduce k₁** (use resolution enhancement techniques — more complex process)

**Historical wavelength reduction:**
```
Era         Light source     Wavelength     Min feature
1970s       Mercury (g-line) 436nm          ~1µm
1980s       Mercury (i-line) 365nm          ~300nm
1990s       KrF excimer      248nm          ~130nm
2000s       ArF excimer      193nm          ~65nm
2010s       ArF + immersion  193nm, n=1.44  ~38nm (single patterning)
2020s       EUV              13.5nm         ~13nm (single patterning)
```

**Numerical Aperture**: The sine of the half-angle of the light cone entering the lens, multiplied by the refractive index of the medium. Dry lens: NA < 1 (limited by n_air = 1). Immersion: NA up to 1.35 (n_water = 1.44).

### Quick Check
> 1. Write the Rayleigh criterion formula for lithographic resolution.
> 2. How does increasing NA improve resolution?
> 3. Why did the industry move from mercury lamps to excimer lasers for lithography?

---

## 3. Immersion Lithography and ArF Scanners

The biggest single improvement in ArF (193nm) lithography was **immersion**: filling the gap between the lens and wafer with ultrapure water instead of air.

Water has refractive index n = 1.44, so the effective wavelength in water is λ/n = 193nm/1.44 = 134nm. More importantly, NA = n×sin(θ) can now exceed 1.0, enabling NA = 1.35 on production scanners.

```
Dry ArF (NA=0.93):         Immersion ArF (NA=1.35):
  
  Last lens element          Last lens element
  ─────────────────          ─────────────────
   air gap (n=1)              water film (n=1.44)
  ─────────────────          ─────────────────
  Wafer                      Wafer
  
  λ_eff = 193nm              λ_eff = 193/1.44 = 134nm
  R = 0.4 × 193 / 0.93      R = 0.4 × 193 / 1.35
    = 83nm                     = 57nm (better resolution)
```

**ASML TWINSCAN NXT (ArF immersion)**: The industry-standard ArF immersion scanner. Throughput: ~275 wafers/hour. Dose: ~30 mJ/cm². Stage precision: <1nm overlay accuracy. Cost: ~$80 million per system.

**Double/multi-patterning**: Even with immersion, 193nm cannot directly print features smaller than ~38nm in a single exposure. Below this, multiple exposures with different masks are used (covered in Section 6).

**The 193nm plateau**: The semiconductor industry used 193nm ArF immersion from 2007 through 2019 for leading-edge production (with multi-patterning). Moving to shorter wavelength (F₂ laser at 157nm) was tried and abandoned. EUV at 13.5nm was the jump that eventually came in 2019.

### Quick Check
> 1. Why does immersion lithography improve resolution compared to dry lithography?
> 2. What is NA and what is the maximum achievable NA with water immersion?
> 3. Why did the industry stay on 193nm for so long (2007–2019) rather than moving to shorter wavelengths?

---

## 4. EUV Lithography — The Breakthrough

**EUV (Extreme Ultraviolet) lithography** uses 13.5nm wavelength light — about 14× shorter than ArF. This single change allows printing features as small as ~13nm in a single exposure.

**Why EUV was so hard to develop:**

1. **EUV is absorbed by everything**: Unlike visible light or ArF that passes through glass, 13.5nm EUV is absorbed by air, lenses, and even the reticle. The entire optical path must be in vacuum, and lenses must be replaced by mirrors.

2. **Source power**: The EUV light source uses high-powered CO₂ laser pulses (40kW!) to vaporize tiny tin droplets (~30µm diameter, 50,000 times/second), creating tin plasma that emits 13.5nm radiation. Collecting ~5% of this light with a collector mirror. The source is spectacularly complex.

3. **Reflective optics only**: 6 ultra-polished multilayer mirrors (Mo/Si alternating layers, λ/4 thick) reflect 13.5nm light with ~67% reflectivity each. The mask (reticle) is also reflective. Each mirror is polished to <0.2nm rms roughness — if the mirror were scaled to Earth's surface, the tallest mountain would be 2.5cm.

4. **Photon shot noise**: At 13.5nm, fewer photons per dose means larger statistical fluctuations → stochastic defects (randomly rough edges). Higher dose fixes this but reduces throughput.

```
EUV system schematic:
  
  CO₂ laser → Tin droplets → EUV plasma (13.5nm)
       │
  Collector mirror (ellipsoidal, Mo/Si multilayer)
       │
  Intermediate focus (IF) aperture
       │
  Illumination system (6 mirrors)
       │
  Reflective reticle (EUV absorber + Mo/Si multilayer)
       │
  Projection optics (6 mirrors, NA=0.33 or 0.55)
       │
  Wafer (in vacuum)
  
  Total: 12+ mirrors. Each reflects ~67% → 0.67¹² ≈ 1% throughput
  Need enormous source power to compensate
```

**High-NA EUV (NA=0.55)**: ASML Twinscan EXE:5000, releasing 2025+. NA increased from 0.33 to 0.55 → resolution from 13nm to ~8nm. First customer: TSMC, Samsung. Cost: ~$350 million per tool.

**EUV production**: TSMC started volume EUV production in 2019 (N7+), now uses it extensively at N5/N4/N3/N2. Samsung used EUV at 5nm. Intel uses EUV at Intel 4/3/20A/18A.

### Quick Check
> 1. Why must EUV lithography use mirrors instead of lenses?
> 2. How does ASML's EUV system generate 13.5nm light?
> 3. What is "stochastic defectivity" in EUV and what causes it?

---

## 5. Resolution Enhancement Techniques

Even with EUV, optical tricks are needed to print patterns beyond the nominal diffraction limit. These **Resolution Enhancement Techniques (RETs)** manipulate the light to improve contrast and pattern fidelity.

**OPC (Optical Proximity Correction)**: The photomask pattern is pre-distorted to compensate for diffraction effects that would otherwise cause corners to round and lines to vary in width. A modern chip layout has trillions of OPC corrections applied computationally before mask writing.

```
Without OPC:          After OPC:
  Desired shape:         Mask shape (pre-distorted):
  ┌─────┐               ┌──────┐
  │     │               │  ╔╗  │
  │     │           →   │  ╔╝  │
  └─────┘               └──────┘
  
  After printing, diffraction rounds the corners.
  OPC adds "serifs" and "hammerheads" to compensate.
  Printed result ≈ desired shape.
```

**Phase-Shift Masks (PSM)**: Modify the mask to shift the phase of the light through some openings by 180°. Destructive interference occurs at edges → sharper edge definition. Alternating PSM (altPSM) and attenuated PSM (attPSM) are both used.

**Off-Axis Illumination (OAI)**: Instead of illuminating the mask on-axis (straight ahead), use annular, dipole, or quadrupole illumination patterns. The interference of diffraction orders improves resolution for specific pattern types.

**Source-Mask Optimization (SMO)**: Jointly optimize the illumination pattern AND the mask design using computational algorithms. State-of-the-art technique combining OAI + OPC.

**Inverse Lithography Technology (ILT)**: Computationally invert the lithography equation to find the mask pattern that produces the target printed pattern. Produces "curvilinear" mask shapes (not just Manhattan geometry) for maximum fidelity. Enabled by high-NA EUV computational lithography.

### Quick Check
> 1. What is OPC and why is it needed?
> 2. What is a Phase-Shift Mask and how does it improve resolution?
> 3. What is Source-Mask Optimization?

---

## 6. Multi-Patterning — Cheating the Diffraction Limit

When a single lithographic exposure cannot print a pattern fine enough, **multi-patterning** splits the pattern across multiple exposures.

**LELE (Litho-Etch-Litho-Etch) double patterning**:
1. Print alternating lines with exposure 1 (lines A)
2. Etch those lines into a hard mask
3. Print intervening lines with exposure 2 (lines B) — aligned precisely to lines A
4. Etch lines B

Result: line pitch is halved vs single exposure.

```
Single exposure:          Double patterning (LELE):
  Can print pitch P        Prints pitch P/2
  ──  ──  ──  ──  ──  ──   ──────────────────
  (limit of diffraction)   Exposure 1: A  A  A
                           Exposure 2:  B  B  B
                           Result:     ABABAB
                           Pitch = P/2 ✓
```

**SADP (Self-Aligned Double Patterning)**: Instead of a second lithography step, grow a thin spacer (sidewall material) around the first-patterned lines, then remove the original lines. The spacers form the final pattern at half the original pitch. Only one lithography step needed — self-aligned, no overlay error between the two sets of features.

**SAQP (Self-Aligned Quadruple Patterning)**: Apply SADP twice → 4× pitch reduction. Used extensively at 5nm and 7nm for dense metal layers.

```
SADP process:
  
  1. Pattern lines (pitch P):   2. Deposit spacer:    3. Remove lines:
  ┌──┐   ┌──┐   ┌──┐           ┌┐┌──┐┌┐  ┌┐┌──┐┌┐   ┌┐    ┌┐  ┌┐    ┌┐
  └──┘   └──┘   └──┘           └┘└──┘└┘  └┘└──┘└┘   └┘    └┘  └┘    └┘
  P                             spacer grown           only spacers remain
  
  Final pitch = P/2, defined entirely by spacer width (controlled by ALD thickness)
  No second lithography needed!
```

**Multi-patterning limitations**: More steps = longer cycle time, higher defect risk, overlay errors. SAQP at 5nm requires extraordinary overlay accuracy (<1nm). This is one major motivation for EUV — single EUV exposure can replace 3-4 patterning steps.

### Quick Check
> 1. What does LELE double patterning achieve and how many lithography steps does it use?
> 2. What is SADP and how does it halve the pitch without a second lithography step?
> 3. Why does multi-patterning increase manufacturing cost and defect risk?

---

## 7. ASML: The Monopoly That Controls Chips

**ASML** (Advanced Semiconductor Materials Lithography, Netherlands) is the only company in the world that makes EUV lithography machines. Every advanced chip in the world — Apple M-series, NVIDIA H100, AMD EPYC — was made on ASML equipment.

**How ASML got this monopoly:**
- EUV development started in the 1990s as a research project backed by Intel, AMD, and Micron through EUV LLC consortium
- Massive technical challenges required sustained $8B+ investment over 20 years
- Competing efforts (Canon, Nikon, SVG) failed to sustain investment
- ASML's unique supply chain: Carl Zeiss (optics), Cymer (laser source), Trumpf (CO₂ laser), hundreds of specialized suppliers
- First production EUV delivery: 2012; first volume production: 2019

**ASML's key products:**
- TWINSCAN NXT:2000i: ArF immersion, ~$80M, 275 wph
- TWINSCAN NXE:3600D: EUV (NA=0.33), ~$170M, 160 wph
- TWINSCAN EXE:5000: High-NA EUV (NA=0.55), ~$350M, 2025+ delivery

**Geopolitical dimension**: The Netherlands, at US pressure, has banned ASML from exporting EUV machines to China since 2019. China can obtain older immersion ArF scanners, but not EUV. This is a critical bottleneck in China's semiconductor ambitions — SMIC (China's leading foundry) is estimated to be 5–10 years behind TSMC partly because it cannot access EUV.

**The supply chain dependency graph:**
```
ASML EUV machine:
  ├── Zeiss optics (mirrors polished to 0.2nm rms) — Germany
  ├── Cymer/ASML laser source (Sn droplet + CO₂ laser) — USA/Netherlands
  ├── Trumpf CO₂ laser — Germany
  ├── Novellus vacuum systems — USA
  ├── Thousands of precision components from ~5,000 suppliers globally
  └── Total: ~100,000 parts per machine
```

### Quick Check
> 1. Why does ASML have a monopoly on EUV lithography machines?
> 2. What is the approximate cost of a High-NA EUV machine?
> 3. Why is the US export ban on EUV machines to China strategically significant?

---

## Summary

- **Photolithography** transfers circuit patterns to silicon using light through a patterned mask (reticle). Key steps: coat resist → expose → develop → etch/implant.
- **Resolution limit**: R = k₁ × λ / NA. Smaller λ and larger NA improve resolution.
- **Wavelength history**: 436nm (g-line) → 248nm (KrF) → 193nm (ArF) → 193nm immersion (NA=1.35) → 13.5nm EUV.
- **EUV**: Uses tin plasma to generate 13.5nm light; all reflective optics; vacuum required; complex source. TSMC production since 2019. High-NA EUV (NA=0.55) coming 2025.
- **RETs**: OPC (mask pre-distortion), PSM (phase shift), OAI, SMO improve resolution.
- **Multi-patterning**: LELE, SADP, SAQP print finer features by splitting exposures or using self-aligned spacers.
- **ASML**: Only EUV supplier; ~$170M for standard EUV, ~$350M for high-NA. Geopolitically significant (export restrictions to China).

---

## Exercises

### Easy
1. What is a photoresist and what happens when UV light hits it?
2. What is the Rayleigh criterion and what two parameters can you change to improve resolution?
3. Why can't EUV lithography use glass lenses like ArF lithography?

### Medium
4. Resolution calculation: An ArF immersion scanner has λ=193nm, NA=1.35, k₁=0.28. (a) Calculate the minimum resolvable feature. (b) For EUV with λ=13.5nm, NA=0.33, k₁=0.4: minimum feature? (c) For High-NA EUV with NA=0.55, k₁=0.35: minimum feature? (d) A chip design needs 5nm metal pitch (half-pitch = 2.5nm). Can standard EUV achieve this? What approach is needed?
5. Multi-patterning math: At 7nm node, the metal 2 (M2) pitch needs to be 40nm. (a) ArF immersion can print 80nm pitch (single exposure). How many patterning stages (LELE) are needed? (b) Using SADP instead: starting pitch 160nm, one SADP gives pitch 80nm. Two SADP (SAQP): what pitch? (c) EUV single exposure can achieve 40nm pitch directly. How many fewer process steps does this represent?
6. ASML cost per exposure: An ASML EUV scanner costs $170M, lasts 10 years, processes 160 wafers/hour (running 24/7). (a) Total wafers in 10 years (include 85% uptime). (b) Depreciation cost per wafer? (c) A 300mm wafer has ~883 chips at 80mm² die size. What is the lithography depreciation cost per chip just for the EUV step? (d) If the chip design needs 15 EUV layers + 50 ArF immersion layers (at $80M per ArF scanner), what is the total lithography equipment cost per chip (rough estimate)?

### Hard
7. EUV source physics: The ASML EUV source uses a 20kW CO₂ laser (10.6µm wavelength) pre-pulse and main pulse to vaporize a 30µm tin (Sn) droplet, creating hot plasma emitting 13.5nm EUV. (a) The CO₂ laser wavelength is 10.6µm — what temperature plasma radiates at 13.5nm? Use Wien's law: λ_max × T = 2.9×10⁻³ m·K. (b) The source converts CO₂ laser power to EUV with ~2% efficiency. To get 250W of EUV at intermediate focus: what CO₂ laser power is required? (c) The collector mirror collects ~50% of emitted EUV. Optical train from source to wafer has 10 mirrors each reflecting 67%. What fraction of EUV power reaches the wafer? (d) At 250W EUV dose and 20 mJ/cm² required dose: what wafer throughput (wafers/hour) is achievable? (Use resist area per exposure = 33mm × 26mm.)
8. China semiconductor dilemma: China's SMIC cannot obtain EUV machines. They have access to ArF immersion (193nm, NA=1.35). (a) With k₁=0.28 and SAQP (4× pitch reduction): what minimum pitch can SMIC achieve? (b) Compare to TSMC's N3 (using EUV, single exposure at k₁=0.35): what pitch does TSMC achieve? (c) China is investing $150B in domestic semiconductor equipment. SMIC has reportedly achieved 7nm-class chips using multi-patterning. What is the cost multiplier (process steps, yield loss, cycle time) of SAQP vs EUV for the same pitch? (d) How many years would it take China to develop domestic EUV capability if: lithography machine development takes ~15 years with $10B/year investment and requires specialized supply chain components currently unavailable domestically (Zeiss-quality optics)?
