# Chapter 63: Additive Manufacturing of Metals — Building Layer by Layer

> **"Additive manufacturing does not make parts that conventional processes make faster or cheaper — at least not yet. What it does is make parts that could not be made before: cooling channels that curve and bifurcate like blood vessels; lattice structures with controlled porosity; blade tip repairs that add back lost material without welding; geometries optimized by algorithms with no regard for whether a tool can reach them. AM is not a replacement for investment casting. It is a new capability that changes what is possible."**

---

## Table of Contents

1. [What Is Metal AM and Why Does It Matter?](#1-what-is-metal-am-and-why-does-it-matter)
2. [Selective Laser Melting (SLM / LPBF)](#2-selective-laser-melting-slm--lpbf)
3. [Electron Beam Melting (EBM)](#3-electron-beam-melting-ebm)
4. [Directed Energy Deposition (DED)](#4-directed-energy-deposition-ded)
5. [AM Microstructure — What Makes It Different](#5-am-microstructure--what-makes-it-different)
6. [Defects in Metal AM](#6-defects-in-metal-am)
7. [Post-Processing Requirements](#7-post-processing-requirements)
8. [AM Materials: Alloys and Their Challenges](#8-am-materials-alloys-and-their-challenges)
9. [Applications in Aerospace and Gas Turbines](#9-applications-in-aerospace-and-gas-turbines)
10. [Design for AM — Topology Optimization](#10-design-for-am--topology-optimization)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Is Metal AM and Why Does It Matter?

**Additive manufacturing (AM)** builds parts by depositing material layer by layer from a 3D digital file:
- No subtractive material removal (minimal waste — critical for expensive alloys)
- No tooling, molds, or dies (reduces lead time from months to days for prototypes)
- Complex internal geometry possible (lattice structures, conformal cooling channels)
- Repair capability (add material back to worn or damaged parts)

**AM ≠ faster machining:**
AM is best justified when:
- Complexity is very high (channels, lattices, organic shapes)
- Volume is very low (prototypes, replacement parts)
- Material is very expensive (titanium, Inconel → buy-to-fly ratio of 20:1 in machining becomes 1.2:1 in AM)
- Lead time is critical

**Main metal AM processes:**
| Process | Heat source | Feedstock | Resolution | Best for |
|---------|------------|---------|-----------|---------|
| SLM / L-PBF | Laser | Powder bed | 50–200 μm | Complex, dense parts |
| EBM | Electron beam | Powder bed | 100–300 μm | Ti, Ni, near-net shape |
| DED/LENS | Laser + nozzle | Powder blown | 0.5–2 mm | Repair, large parts, clad |
| Wire Arc AM | Arc | Wire | 1–5 mm | Large Ti/steel structures |
| Binder jetting | None (sintering later) | Powder binder | 30–100 μm | Complex steel, Cu parts |

---

## 2. Selective Laser Melting (SLM / LPBF)

**Also called:** Laser Powder Bed Fusion (L-PBF), DMLS (Direct Metal Laser Sintering — proprietary term)

**Process:**
```
1. Powder layer deposited by recoater blade (20–60 μm thick)
2. Laser scans the cross-section → melts + solidifies
3. Build plate lowers by one layer thickness
4. Repeat for next layer
5. Remove unsintered powder (reused)
6. Part is surrounded by and resting on unfused powder
```

**Laser parameters:**
- Power: 100–500 W (single laser) → up to 4 × 1,000 W (multi-laser systems)
- Scan speed: 500–2,000 mm/s
- Layer thickness: 20–60 μm
- Hatch spacing: 60–120 μm

**Melt pool characteristics:**
```
Melt pool shape: elongated ellipse (~150 μm long, ~100 μm wide, ~100 μm deep)
Cooling rate: 10⁵–10⁷ °C/s (orders of magnitude faster than casting!)
Temperature gradient G: 10⁶–10⁷ °C/m
Solidification velocity R: 0.1–1 m/s
```

The extreme cooling rates produce very fine, non-equilibrium microstructures.

**Powder characteristics required for SLM:**
- Particle size: 10–45 μm (fine enough for thin layers)
- Spherical morphology (good flowability for even powder spreading)
- Low oxygen content: < 200 ppm (oxidation → porosity in melt)
- Gas-atomized powder (Ch 35) is the standard

**Chamber atmosphere:** Ar or N₂ inert atmosphere (O₂ < 100 ppm during build)

---

## 3. Electron Beam Melting (EBM)

**Process:** Same powder bed concept, but uses an electron beam instead of a laser.

**Key differences from SLM:**

| Feature | SLM (Laser) | EBM |
|---------|-------------|-----|
| Heat source | Laser (photons) | Electron beam |
| Vacuum | No (Ar/N₂ atmosphere) | Yes (< 10⁻⁴ Pa) |
| Build temperature | Near RT (small part heating) | High (700–1,000°C for Ti, Ni) |
| Residual stress | High (rapid cooling) | Lower (pre-heated bed) |
| Build speed | Slow | Faster (beam deflection: no inertia) |
| Surface finish | Better (Ra ~5–15 μm) | Rougher (Ra ~20–35 μm) |
| Resolution | Better (50–100 μm) | Less good (100–300 μm) |

**EBM advantages for Ti-6Al-4V:**
- Vacuum: no oxygen pickup (Ti is extremely reactive at high T)
- Pre-heated powder bed (700°C): reduces residual stress; allows α+β microstructure to develop (instead of martensitic as in SLM Ti)
- Result: EBM Ti-6Al-4V has α+β microstructure, better ductility than SLM

**EBM for Ni superalloys:**
At build temperature ~1,000°C, γ′ starts to precipitate during build → complex in-situ aging.

---

## 4. Directed Energy Deposition (DED)

**DED** melts material (powder or wire) as it is deposited, using a focused laser or electron beam:

```
LENS (Laser Engineered Net Shaping):
   Laser beam
      ↓
   Melt pool on substrate
   ↑
   Powder nozzle(s) (4-nozzle ring)
   
Build head moves → deposits bead → next layer on top
```

**Key DED characteristics:**
- Layer thickness: 0.5–3 mm (much coarser than SLM)
- Deposition rate: 10–100× faster than SLM
- Can deposit onto EXISTING parts (repair, cladding, hybrid manufacturing)
- Can deposit multiple materials in same build (graded composition)

**Applications of DED:**
- **Repair of turbine blades:** Add back material to worn blade tips, rebuild damaged leading edges
- **Repair of molds and dies:** Rebuild worn die steel areas
- **Cladding:** Apply corrosion-resistant layer on steel substrate
- **Functionally graded material (FGM):** Transition from one alloy to another within a single part

**Wire Arc AM (WAAM):**
- Uses welding wire + electric arc (MIG/TIG) → deposits beads → builds large structure
- Deposition rate: 1–10 kg/hour (vs 0.1 kg/hour for SLM)
- Resolution: ~5 mm
- Applications: large Ti-6Al-4V aerospace frames (Norsk Titanium, GKN)

---

## 5. AM Microstructure — What Makes It Different

**Three distinctive features of AM microstructure:**

### 1. Epitaxial (Layer-on-Layer) Columnar Grain Growth
- Each layer remelts the previous layer → epitaxial solidification (Ch 36 analogy)
- Columnar grains grow OPPOSITE to heat flow (downward through layers) → vertical columnar grains through build height
- For [001]-textured SLM Ni alloys: columnar grains ALONG build direction

**Consequence:**
SLM Ni superalloy (IN718, Hastelloy X) → strong texture → anisotropic properties:
- Tensile: E, σ_y differ in build direction vs. perpendicular
- Fatigue: crack growth rate differs with orientation

### 2. Extremely Fine, Non-Equilibrium Microstructure
- Cooling rates 10⁵–10⁷ °C/s → fine cellular/dendritic structure at nm–μm scale
- No macro-segregation (each melt pool = tiny isolated solidification event)
- Metastable phases: SLM Ti-6Al-4V → martensitic α' (not α+β from slow cooling)
- Fine cells → high dislocation density → work-hardened-like state → high strength

### 3. Thermal History Complexity
- EVERY layer heats the layers below it → effective "cycling" of lower layers
- Tempering, precipitation, stress relief all happen to each layer as subsequent layers are added
- Microstructure varies through build height

**AS-BUILT SLM IN718 microstructure:**
- Columnar grains, height >> width
- Fine cellular substructure (cell width 200–500 nm)
- NbC precipitates at cell boundaries (Nb segregates during rapid solidification)
- γ″ (Ni₃Nb) NOT fully precipitated (needs dedicated aging cycle)
- Residual stress: compressive at surface, tensile in bulk (from thermal contraction)

---

## 6. Defects in Metal AM

### Porosity (Main AM Defect)
**Type 1 — Lack of Fusion (LOF):**
- Cause: insufficient laser energy → not enough melting → unfused powder trapped
- Morphology: irregular, elongated voids (aligned with scan tracks)
- Prevention: increase energy density = Laser power / (scan speed × hatch spacing × layer thickness)
- Detection: X-ray CT, Archimedes density

**Type 2 — Keyhole Porosity:**
- Cause: too much laser energy → vapor pressure exceeds hydrostatic → keyhole collapses → spherical pore
- Morphology: spherical pores ~50–200 μm
- Prevention: reduce energy density

**Type 3 — Gas Porosity:**
- Cause: entrapped Ar in powder particles (from gas atomization) or moisture → H₂ in melt
- Morphology: spherical, small (< 50 μm)
- Prevention: powder drying; vacuum degassing of powder

### Residual Stress and Cracking
- SLM: very high temperature gradients → large residual stresses (can warp or crack thin walls)
- MITIGATION: scan strategy (rotate 67° between layers = minimizes stress build-up); pre-heat base plate; post-process HIP + stress relief anneal

### Solidification Cracking (Hot Cracking)
- SLM conditions (fast solidification) → same mechanism as welding hot cracking (Ch 36)
- Crack-susceptible alloys in AM: IN738, CMSX-4 (same crack-susceptible as in welding)
- More weldable alloys: IN718 (sluggish γ″ → less strain-age cracking), IN625 (low γ′)

---

## 7. Post-Processing Requirements

AS-BUILT AM parts are rarely used in aerospace directly. Post-processing is essential:

**1. Stress relief anneal:**
- For SLM: 800–1,000°C in inert atmosphere before part removal from build plate
- Reduces residual stress → prevents distortion on removal and machining

**2. Hot Isostatic Pressing (HIP):**
- Closes residual pores (both LOF and keyhole type)
- For aerospace: mandatory for fracture-critical parts
- Typical: 1,120°C / 150 MPa / 3 hours for IN718
- After HIP: porosity < 0.1 vol%

**3. Solution + Age Heat Treatment:**
- SLM IN718: as-built has low γ″ fraction → age at 720°C + 620°C to develop full strength
- SLM Ti-6Al-4V: as-built has α' martensite → anneal to convert to α+β
- Must follow certified HT specs (same as for wrought alloy)

**4. Machining of datum surfaces:**
- AM surfaces (Ra ~5–20 μm as-built) too rough for mating surfaces, seals, fillet radii
- Machine critical surfaces to drawing tolerance
- CNC 5-axis machining of AM parts → complex setup due to AM net shape

**5. Surface treatment:**
- Shot peen (where required): same as conventional → compressive RS → fatigue improvement
- Electropolish / abrasive flow machining: internal channels that can't be machined

---

## 8. AM Materials: Alloys and Their Challenges

**Titanium — Ti-6Al-4V:**
- Most common AM metal by volume
- SLM: α' martensite as-built → post-anneal to α+β → ductility recovered
- EBM: directly α+β (high build T)
- Challenge: hot cracking absent (α' solidification without segregation)
- Properties after HIP+HT: match or exceed wrought min

**IN718 (Nickel superalloy):**
- Most common AM Ni alloy
- SLM as-built: high residual stress, columnar grains, low γ″
- After HIP + standard aging: σ_y ≥ 1,100 MPa, UTS ≥ 1,300 MPa (AMS 5662 wrought min)
- Challenge: columnar texture → anisotropy in fatigue

**Stainless 316L:**
- Good SLM workability (austenitic, ductile)
- As-built: cellular substructure, high dislocation density → higher σ_y than wrought
- Challenge: for elevated T service, cellular structure coarsens → properties degrade

**Hastelloy X (IN625 variant):**
- Excellent laser AM workability (no hot cracking, weldable)
- Used: combustor liners, heat shields (complex shapes with film holes)
- Challenge: lower strength than aging-hardened alloys → limited HPT applications

**CMSX-4 and other SX alloys:**
- Extremely difficult in SLM: wide solidification range → hot cracking; high Re segregation
- Research only: special scan strategies + selective nucleation to achieve SX-like texture
- Not in production AM (as of 2025)

**Hardened steels (17-4PH, 15-5PH, H13):**
- SLM possible; near-wrought properties after HIP+age
- Tool and die repair by DED

---

## 9. Applications in Aerospace and Gas Turbines

**GE Aviation CF6/GE9X fuel nozzles:**
- First certified AM metal part in a commercial jet engine (2015)
- Previously 20 welded pieces → now 1 printed piece
- Material: Co-Cr alloy
- Benefit: 5× more durable, 25% lighter, one-third cost

**Airbus A350 titanium structural bracket (AM):**
- Topology-optimized bracket: 30% lighter than machined equivalent
- Ti-6Al-4V SLM
- Buy-to-fly ratio: 1.2 (vs 20 for machined from billet)

**GE Catalyst aircraft engine (Czech factory):**
- >35% of parts by AM
- Reduced part count 855 → 12 AM assemblies

**Turbine blade tip repair (DED):**
- Fan and LPT blade tip wear → DED adds back Ti-6Al-4V or IN718 material → machine to profile
- Replaces full blade replacement → saves $3,000–$8,000 per blade

**Combustor liners (SLM Hastelloy X):**
- Complex swirler geometry, multiple cooling holes, lightweight lattice walls
- Weight reduction 30–40%

**Future — AM HPT blades:**
- Current limitation: grain structure, porosity, and anisotropy limit HPT use
- Research goal: controlled epitaxial SLM to grow [001]-textured columnar grains (like DS casting)
- Expected: AM HPT blades in demonstrator engines by ~2030

---

## 10. Design for AM — Topology Optimization

**Topology optimization:** Mathematical algorithm finds the minimum-mass structure that satisfies stress, displacement, and eigenfrequency constraints.

**Process:**
```
Start: filled design domain
Apply: loads, boundary conditions, material properties
Algorithm: removes material from low-stress regions
           adds material (or maintains) in load-bearing paths
Result: organic, bone-like structure

Conventional machining: cannot make this (no tool access to internal voids)
AM: can make this directly from digital model
```

**Lattice structures:**
AM allows internal lattice structures (e.g., body-centered cubic lattice: BCC, face-centered cubic: F2CC):
- Fill volume fraction: 20–60%
- Function: support thermal insulation, reduce mass, absorb impact
- Applications: TBC-replacement concepts, heat exchangers, satellite brackets

**Conformal cooling channels:**
Injection molds and die casting dies need cooling:
- Conventional: straight-drilled circular cross-section channels
- AM: free-form channels following the mold surface contour → faster, more uniform cooling → shorter cycle time, better part quality

---

## Summary

| Process | Layer size | Build rate | Resolution | Best for |
|---------|-----------|-----------|-----------|---------|
| SLM (L-PBF) | 20–60 μm | Slow | Excellent | Complex, precise, small |
| EBM | 50–100 μm | Medium | Good | Ti, Ni (low residual stress) |
| DED/LENS | 0.5–3 mm | Fast | Poor | Repair, large parts, cladding |
| Wire Arc AM | 1–5 mm | Very fast | Poor | Large Ti/steel structures |

**Post-processing always needed:** Stress relief → HIP → Solution + age → Machining → NDT

---

## Exercises

1. SLM parameter optimization for IN718: (a) Volumetric energy density E_v = P / (v × h × t) where P = laser power (W), v = scan speed (mm/s), h = hatch spacing (mm), t = layer thickness (mm). For P = 200W, v = 1,000 mm/s, h = 0.1 mm, t = 0.04 mm, calculate E_v in J/mm³. (b) For full fusion of IN718, E_v should be 50–100 J/mm³. Is this value in range? (c) To achieve density > 99.5% (measured by Archimedes), research shows E_v needs to be 60–80 J/mm³. If power is fixed at 200W, what range of scan speeds achieves this? (d) Above E_v = 100 J/mm³, keyhole porosity forms. What scan speed would you avoid?

2. SLM Ti-6Al-4V microstructure comparison: (a) As-built SLM: α' martensite (HCP, metastable, fine plates, high dislocation density). After stress relief at 800°C/2h: α+β (HCP α + BCC β, equilibrium). EBM: α+β (α plates ~5 μm wide). Compare predicted tensile properties using Hall-Petch type analysis: SLM as-built (d = 1 μm α' width), SLM stress-relieved (d = 3 μm), EBM (d = 5 μm). (b) Fatigue crack growth rate in SLM Ti-6Al-4V is higher than wrought in the build direction but lower perpendicular to build direction. Explain this anisotropy in terms of crack path relative to columnar grain boundaries. (c) After HIP at 920°C / 100 MPa / 2h: α lath width = 8 μm (coarsened). How does this affect HCF life compared to SLM as-built?

3. DED turbine blade tip repair: A HPT blade (IN718) has measured tip wear of 1.5 mm across 10 mm chord width. Repair specification: replace 1.5 mm depth × full chord width. (a) DED deposition rate = 6 cm³/hour at tip area. Blade tip area = 20 × 80 mm (root to tip section). Estimate total volume added and repair time. (b) DED IN718 composition matches base material. But the HAZ in the existing blade is re-heated during deposition. What microstructural change occurs in the existing γ″-strengthened IN718 at 600–850°C? (c) After DED repair, the blade requires heat treatment. The standard aging (720°C/8h + 620°C/8h) is performed. But the DED region and base material may respond differently if the DED region's γ' solvus is slightly different. How would you verify the DED region achieved proper precipitation? (d) Calculate cost saving per blade: replacement cost = $8,000; repair cost (DED + HT + inspect) = $2,500. How many blades per year justify purchasing a $500,000 DED system if the system saves the cost difference per repair?

4. Topology optimization of an AM bracket (Ti-6Al-4V): The original machined bracket weighs 350g (100% dense). Topology optimization removes 55% of material, leaving a truss-like structure of 157g. (a) Verify the mass: density of Ti-6Al-4V = 4.43 g/cm³. What volume does the original part occupy? What volume after optimization? (b) The bracket experiences peak stress of 180 MPa in the thin struts of the topology-optimized design. If σ_y (SLM Ti-6Al-4V after HIP+HT) = 860 MPa, calculate safety factor. (c) A lattice infill with 30% density is used in a core section to further reduce mass by 70g while maintaining stiffness (lattice acts like a lower-modulus bulk material). Calculate equivalent elastic modulus of the lattice if E_Ti = 114 GPa and E_lattice = E_Ti × (ρ/ρ_solid)^2 (Gibson-Ashby relationship). (d) In a vibration test, the AM bracket's 1st natural frequency = 850 Hz, specification requires > 800 Hz. Does it pass?

5. AM defect detection: An SLM IN718 fuel nozzle requires fatigue certification. Specification: no pores > 0.3 mm diameter, no LOF > 0.5 mm. (a) X-ray CT scan (lab, 100 μm voxel) detects pores > 200 μm. Does this meet detection requirements for 0.3 mm pores? (b) After HIP (1,120°C / 150 MPa / 3h): Archimedes density = 8.18 g/cm³ vs theoretical 8.19 g/cm³. Calculate % porosity. (c) Fatigue test of 5 AM samples (after HIP+HT) to runout (10⁷ cycles) at σ_a = 350 MPa: 4 samples survive, 1 fails at 2×10⁶ cycles. Post-fracture: 0.25 mm pore at failure origin. The pore was below CT detection limit (200 μm). This illustrates the "detection gap." What would you specify as the maximum allowed fatigue stress level to ensure a 0.25 mm pore doesn't propagate? Use K_I = 1.1σ√(πa), K_threshold = 3.5 MPa√m. (d) Suggest an improved inspection protocol that reduces the "detection gap" for critical AM parts.
