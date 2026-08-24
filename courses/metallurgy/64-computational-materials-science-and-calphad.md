# Chapter 64: Computational Materials Science and CALPHAD — Designing Alloys by Computer

> **"Every SX turbine blade alloy from CMSX-4 to TMS-238 was designed with pencil and paper, intuition, and thousands of casting trials over decades. Today, a graduate student with a laptop running CALPHAD software can predict in seconds whether a new alloy will form TCP phases, what temperature its γ′ will dissolve, and how its lattice parameter compares to pure Ni. Combined with machine learning trained on experimental databases, we can now screen millions of hypothetical alloys in days. This doesn't replace metallurgical intuition — it amplifies it by orders of magnitude."**

---

## Table of Contents

1. [Why Computational Materials Science?](#1-why-computational-materials-science)
2. [CALPHAD — Computational Thermodynamics](#2-calphad--computational-thermodynamics)
3. [Phase Diagram Calculations with CALPHAD](#3-phase-diagram-calculations-with-calphad)
4. [CALPHAD for Superalloy Design](#4-calphad-for-superalloy-design)
5. [Density Functional Theory (DFT)](#5-density-functional-theory-dft)
6. [Molecular Dynamics (MD) Simulations](#6-molecular-dynamics-md-simulations)
7. [Phase Field Modeling](#7-phase-field-modeling)
8. [Machine Learning in Materials Science](#8-machine-learning-in-materials-science)
9. [Integrated Computational Materials Engineering (ICME)](#9-integrated-computational-materials-engineering-icme)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Computational Materials Science?

**Traditional alloy development:**
- Try a composition → cast it → test properties → modify → repeat
- 10+ years, $50M+ to develop and qualify a new turbine alloy
- Exploratory space: 10-component alloy with 10 possible additions at 5 concentration levels = 5^10 × combinatorial problem = too large to explore experimentally

**Computational approach:**
- Screen compositions computationally → only do experiments on the most promising candidates
- Predict properties before synthesis (phase stability, lattice parameter, elastic modulus)
- Understand mechanisms at atomic scale (why does Re slow creep? — DFT answers this)
- Reduce development time from 10 years to 2–3 years

**Hierarchy of computational methods:**
```
Scale:       Atoms            nm            μm             mm
             ↑                ↑              ↑               ↑
Method:    DFT/MD       Monte Carlo    Phase Field     FEM/CALPHAD
             |                |              |               |
Predicts:  Bond          Short-range    Precipitate      Phase
           strengths       order         coarsening       diagrams
                                         and growth       in service
```

No single method spans all scales → ICME integrates them.

---

## 2. CALPHAD — Computational Thermodynamics

**CALPHAD** = CALculation of PHAse Diagrams (developed by Kaufman and Bernstein in the 1970s).

**Core principle:**
Represent the Gibbs energy of each phase in the system as a function of composition and temperature using empirically-fitted parameters. Then minimize the total Gibbs energy of the system to find the stable phase assemblage.

**Gibbs energy of a phase φ:**
```
G_φ(T, x) = Σ x_i × G°_i(T)           (reference states — pure elements)
           + RT Σ x_i × ln(x_i)         (ideal mixing entropy)
           + G^exc_φ(T, x)               (excess (non-ideal) interaction term)
           + G^magnetic_φ(T, x)          (magnetic contribution for Fe alloys)

where:
  x_i = mole fraction of element i
  G°_i(T) = Gibbs energy of pure element i in phase φ
  G^exc = fitted polynomial in x_i and T (Redlich-Kister series)
```

**Equilibrium condition:**
At equilibrium, the chemical potential of each component is equal in all co-existing phases:
```
μ_i^α = μ_i^β = μ_i^γ   (for all components i)
```
Solved numerically by Gibbs energy minimization → gives phase fractions and compositions.

**CALPHAD databases:**
All the G_φ(T, x) parameters come from critically-evaluated databases:
- NIST: thermodynamic data for pure elements
- Commercial databases: TCNI (for Ni alloys, Thermo-Calc), COST507 (Al alloys), SGTE
- These are built from experimental measurements (calorimetry, phase diagrams) over decades

**Commercial CALPHAD software:**
- Thermo-Calc (most widely used for alloy design)
- FactSage (metallurgical systems)
- Pandat (CompuTherm)
- MTDATA (National Physical Laboratory)

---

## 3. Phase Diagram Calculations with CALPHAD

**Binary phase diagram calculation:**
For Ni-Al binary:
- Input: database + temperature range + composition range
- Output: phase boundaries, solvus lines, liquidus/solidus temperatures

```
CALPHAD Ni-Al output (relevant to turbine alloys):

T (°C)
1400─┤ Liquid
     │     
1300─┤        γ + L
     │
1200─┤   γ + L        γ (FCC solid solution)
     │
1100─┤────────────────── γ solvus for γ′ precipitation ~~
1000─┤              γ + γ′
 900─┤
     │
 800─┤
     └──────────────────────────────
     0             10%Al        20%Al
```

**Multi-component phase diagram:**
For a 10-component superalloy (Ni-Cr-Co-W-Re-Al-Ti-Ta-Mo-Hf):
- Cannot display as a 2D diagram
- Instead: "property diagrams" — plot one variable (T) vs one composition at all other compositions fixed
- Or: "pseudo-binary" at a cut through the multi-component space

**Isopleth (section through multi-component diagram):**
Fix all elements except one varying element → plot T vs that element's content → isopleth.
Use: find γ′ solvus temperature as function of Al/Ti/Ta for alloy screening.

---

## 4. CALPHAD for Superalloy Design

**What CALPHAD predicts for turbine alloy design:**

**1. γ′ volume fraction at temperature:**
```
Thermo-Calc output for CMSX-4 at 900°C:
  γ (FCC) fraction: 0.32 (32%)
  γ′ (L1₂) fraction: 0.68 (68%)
  
  → High γ′ fraction → good creep resistance (consistent with experimental 65–70%)
```

**2. γ′ solvus temperature:**
- T at which γ′ fully dissolves → determines solution treatment temperature
- CALPHAD-predicted T_solvus for CMSX-4 ≈ 1,280°C (experimental: 1,277°C — excellent agreement)

**3. TCP phase stability:**
- Calculate stability of σ (sigma), μ, P phases vs T and composition
- If CALPHAD shows σ-phase Gibbs energy < γ → σ is stable → alloy rejected
- PHACOMP (earlier method) was based on electron vacancy number — CALPHAD is more accurate

**4. Density:**
- Calculate molar volume from atomic radii + composition → density
- Directly checks centrifugal stress implications

**5. Lattice parameter and misfit:**
- Calculate lattice parameter of γ and γ′ → misfit δ = (a_γ - a_γ′)/a_γ
- Target: small negative misfit (-0.1 to -0.5%) → good coherency → slow coarsening

**Alloy design workflow using CALPHAD:**
```
1. Define design targets: σ_creep, ρ, T_solvus, no TCP
2. Screen: 10,000 compositions using CALPHAD → filter by targets
3. Down-select to ~100 promising compositions
4. Cast pilot alloys (10g in arc furnace)
5. Measure: DSC (T_solvus), X-ray (lattice parameter), short-term creep test
6. Iterate: compare measured vs predicted → refine database parameters
7. Scale up: 100g VIM melt → DS casting → full property testing
8. Select best composition → qualification testing (AMS, OEM specs)
```

---

## 5. Density Functional Theory (DFT)

**DFT** calculates material properties from quantum mechanics — no adjustable parameters (ab initio).

**Schrödinger equation for N electrons:**
```
Ĥ Ψ = E Ψ

Exact solution: impossible for N > 3
DFT approximation (Hohenberg-Kohn, 1964): 
  Replace many-body problem with single-electron problem
  in an effective potential (exchange-correlation functional)
  → Kohn-Sham equations (solved iteratively)
```

**What DFT can predict (for perfect crystals):**
- Lattice parameters (accuracy: 0.5–2%)
- Elastic constants C₁₁, C₁₂, C₄₄ (accuracy: 5–10%)
- Stacking fault energies (→ determines dislocation width)
- Surface energies (→ γ/γ′ interfacial energy → coarsening rate)
- Migration barriers for diffusion (→ diffusivity → creep rate)
- Electronic structure (bonding character)

**DFT for Re effect in creep:**
DFT calculations show:
- Re atoms in γ Ni have much higher migration barriers than other elements
- Re forms short-range "clusters" that are energetically stable
- Both effects → reduced diffusivity → slower dislocation climb → explains why Re improves creep

This was the first COMPUTATIONAL explanation of the Re effect (experimental observation came first, mechanism was explained by DFT later).

**DFT limitations:**
- T = 0 K (temperature effects require further models — phonons, thermal expansion)
- Perfect crystal (no grain boundaries, dislocations) — approximate methods for defects
- System size: max ~1,000 atoms (small compared to real microstructural features)
- Computationally expensive: days to weeks on supercomputers for large systems

---

## 6. Molecular Dynamics (MD) Simulations

**MD** models atomic motion using classical potentials:
- Assigns a force law (potential energy function) to each pair/group of atoms
- Solves Newton's equations of motion for each atom: F = ma → updates positions + velocities at each timestep (1 fs = 10⁻¹⁵ s)
- Can simulate temperature effects, mechanical loading, crack propagation

**Potentials used for metals:**
- Embedded Atom Method (EAM): standard for FCC metals (Ni, Al, Cu)
- Modified EAM (MEAM): better for multi-component alloys
- Machine-learning potentials: trained on DFT data → high accuracy for specific alloy systems

**System size:**
- MD: 10⁶–10⁹ atoms practical on modern GPUs
- Still small compared to grain size (~10¹² atoms/grain for 10 μm grain)

**MD applications in metallurgy:**

**Dislocation mechanics:**
- Watch how a dislocation moves through a crystal under stress
- Observe it when it reaches a γ/γ′ interface: does it cut γ′ or bypass it?
- Measure the critical resolved shear stress (τ_CRSS) directly

**Grain boundary structure:**
- Simulate HAGB and LAGB structures
- Calculate grain boundary energy
- Watch how atoms move during grain boundary migration (recrystallization)

**Creep mechanisms:**
- At high T, simulate dislocations climbing over γ′ → measure climb rate vs Re content

**Shock and impact:**
- Simulate FOD impact → how plasticity and crack nucleation occur in the first nanoseconds

---

## 7. Phase Field Modeling

**Phase field modeling** describes microstructure evolution by tracking a smooth order parameter field φ(x,t) that varies between 0 (one phase) and 1 (another phase):

**Phase field equation (Allen-Cahn or Cahn-Hilliard):**
```
∂φ/∂t = -M_φ × δF/δφ   (gradient descent on free energy functional F)

where F includes:
  - Bulk free energy (drives phase separation)
  - Gradient energy (penalizes sharp interfaces)
  - Elastic energy (for coherent precipitates → misfit effects)
  - Magnetic/electric energy (if relevant)
```

**What phase field models predict:**

**γ′ coarsening (Ostwald ripening):**
- Simulate N γ′ particles in a γ matrix
- Track how small particles dissolve and large ones grow
- Compute coarsening rate K in r³(t) = r₀³ + Kt → validates Lifshitz-Slyozov-Wagner (LSW) theory

**γ′ rafting under stress:**
- Apply stress + temperature → watch γ′ cubes elongate into plates (rafts)
- Understand which misfit sign leads to N-type vs P-type rafting
- Quantify rafting kinetics vs alloy composition

**Solidification microstructure:**
- Simulate dendritic growth → predict dendrite arm spacing λ vs G × R
- Simulate microsegregation during solidification → composition profiles

**Recrystallization:**
- Model nucleation and growth of new grains after deformation
- Predict grain size after recrystallization anneal vs temperature + strain

---

## 8. Machine Learning in Materials Science

**ML approaches complement CALPHAD + DFT + MD by:**
- Learning correlations from large experimental datasets
- Predicting properties without understanding mechanism (purely data-driven)
- Accelerating DFT by replacing expensive quantum calculations with trained ML potentials

**ML workflow for alloy design:**
```
1. Compile database: 500+ alloy compositions with measured properties
   (fatigue life, creep rate, oxidation rate, hardness, etc.)
   
2. Feature engineering: represent each alloy as a vector of descriptors
   (atomic radii, electronegativities, melting points, valence electrons...)
   
3. Train ML model (Random Forest, Neural Network, Gaussian Process)
   → Model: property = f(descriptors)
   
4. Predict: evaluate model on new compositions (not in training set)
   → Screen millions of hypothetical alloys in seconds
   
5. Select top candidates → experimental validation → iterate
```

**Key ML successes in metallurgy:**
- Predicting austenite stability in steels from composition → validated new TRIP steels
- Predicting creep strength of multi-principal-element alloys (HEA) → guided design of new Ni-Co-Fe-Al-Ti alloys
- Predicting TBC phase stability from composition → identified new pyrochlore TBC materials

**ML-interatomic potentials (MLIPs):**
- Train a neural network on 100,000 DFT calculations → learns the energy surface
- Run MD with ML potential → DFT-accuracy at 10,000× lower cost
- Applications: simulate γ/γ′ interfaces in CMSX-4 with full compositional complexity

**Caution — ML limitations:**
- "Garbage in, garbage out" — bad training data → bad predictions
- Extrapolation beyond training set can be catastrophically wrong
- Black-box → doesn't provide physical understanding
- Always need experimental validation of ML predictions

---

## 9. Integrated Computational Materials Engineering (ICME)

**ICME** links multiple computational methods across length scales into a unified prediction framework:

```
DFT (atoms)
    ↓ (provides: γ/γ′ interfacial energy, diffusivity, elastic constants)
CALPHAD (equilibrium thermodynamics)
    ↓ (provides: phase fractions, compositions, γ′ solvus)
Phase field (microstructure evolution)
    ↓ (provides: γ′ size, morphology, rafting)
Crystal plasticity FEM (deformation at grain scale)
    ↓ (provides: stress-strain response, Schmid factors, slip activity)
Macroscopic FEM (component level)
    ↓ (provides: stress distribution, temperature field, life prediction)
Component Life Prediction (LCF, HCF, creep)
```

**ICME example: designing a new turbine alloy composition**

Step 1 (CALPHAD): Screen 50,000 compositions → 500 have γ′ fraction > 60%, no TCP, density < 9 g/cm³, T_solvus > 1,280°C

Step 2 (DFT): For 500 candidates → calculate misfit, stacking fault energy, Re clustering tendency → down-select to 50

Step 3 (Phase field): Simulate γ′ coarsening at 1,000°C for 1,000 hours → down-select to 20 with slowest coarsening

Step 4 (Crystal plasticity): Predict creep rate → down-select to 5

Step 5 (Experimental): Cast 5 alloys → measure actual creep → choose 1 winner

**Time comparison:**
- Traditional: 30 alloys × 2 years each = decades
- ICME: 5 experimental alloys × 2 years = 2 years (from 50,000 screened)

**ICME in industry:**
- GE Aviation, Rolls-Royce, NASA: ICME programs for next-gen turbine alloys
- DARPA Accelerated Insertion of Materials (AIM) program: demonstrated ICME for new Ti alloy in aircraft in 5 years vs typical 15 years
- Materials Genome Initiative (US DoE): national ICME database infrastructure

---

## Summary

| Method | Scale | Predicts | Accuracy | Cost |
|--------|-------|---------|---------|------|
| DFT | Atom (Å) | E, lattice, diffusion barriers | High (1–5%) | Days/weeks |
| MD | nm | Dislocation motion, GBs, defects | Medium | Hours-days |
| Monte Carlo | nm | Phase stability, short-range order | High | Hours |
| Phase field | μm | Microstructure evolution, coarsening | Medium | Hours |
| CALPHAD | Macroscopic | Phase diagrams, fractions, T_solvus | High (3–10%) | Seconds |
| Machine learning | Any | Any (data-limited) | Variable | Seconds after training |
| Crystal plasticity | μm–mm | Stress-strain, slip, texture | Medium | Hours |
| ICME | All | Full process-structure-property | Medium | Days (multi-tool) |

---

## Exercises

1. CALPHAD calculation for alloy design: Using the lever rule in a binary Ni-Al system at 900°C: (a) From a simplified CALPHAD-predicted phase diagram: at 8 at% Al total, the γ solvus is at 4 at% Al and the γ/γ′ two-phase region extends to 20 at% Al (γ′). Calculate the γ′ volume fraction using the lever rule: f_γ′ = (x_Al - x_γ) / (x_γ′ - x_γ) where x_Al = 8%, x_γ = 4%, x_γ′ = 20%. (b) For CMSX-4, CALPHAD predicts T_solvus = 1,280°C. The solution treatment is performed at 1,295°C (+15°C above T_solvus). What microstructural state exists during solution treatment? What happens if T exceeds 1,315°C (incipient melting)? (c) CALPHAD for CMSX-4 shows σ-phase becomes stable below 800°C for long exposures at 3% Re. How does this guide the lower temperature limit for the aging treatment? (d) How would you use a Thermo-Calc isopleth (T vs %Re, with all other elements fixed at CMSX-4 values) to find the maximum Re content before TCP phase forms at 900°C service temperature?

2. DFT prediction of Re effect: DFT calculations for impurity diffusion in FCC Ni matrix (vacancy-mediated diffusion) give migration barriers (E_m) for various solutes: (a) Ni self-diffusion: E_m = 1.08 eV; W: E_m = 1.35 eV; Re: E_m = 1.61 eV; Cr: E_m = 0.96 eV. Diffusivity D ∝ exp(-E_m/kT). Calculate the ratio D_Ni/D_Re at T = 1,050°C (use k = 8.62×10⁻⁵ eV/K). (b) If creep rate ∝ D_diffusion (for climb-controlled creep), estimate how much Re slows creep compared to pure Ni diffusivity — qualitatively, does this magnitude explain the 2–3× improvement in creep life from Re addition? (c) DFT shows Re has a local binding energy of -0.25 eV to a vacancy (Re-vacancy pair is slightly more stable). What does this mean for Re transport under stress? Does Re cluster near dislocations?

3. Phase field simulation of γ′ coarsening: The Lifshitz-Slyozov-Wagner (LSW) theory (equivalent to phase field results) gives: r³(t) - r₀³ = K×t where K = (8/9) × γ_γγ′ × D × C∞ × V_m² / (RT × (C_eq^γ′ - C_eq^γ)²). Given: γ_γγ′ = 0.04 J/m², D = 4×10⁻¹⁶ m²/s (Al diffusivity in γ at 950°C), C∞ = 0.08 (molar fraction), V_m = 6.6×10⁻⁶ m³/mol, T = 1,223 K, R = 8.314 J/mol·K, (C_eq^γ′ - C_eq^γ) = 0.12. (a) Calculate K in m³/s. (b) Starting from r₀ = 200 nm, calculate γ′ radius after 1,000 hours, 10,000 hours, and 100,000 hours at 950°C. (c) Orowan bypass becomes the dominant strengthening mechanism when r > r* (typically ~200 nm for γ′). At what service time does γ′ grow beyond r* (assuming starting from r₀ = 200 nm means Orowan is already dominant)? (d) The addition of 6% Re reduces D by a factor of 8 (see Exercise 2). How does this affect K and the predicted coarsening time to reach r = 500 nm?

4. ICME workflow for a new Co-based superalloy: A new L1₂-strengthened Co-Al-W alloy is proposed (discovered in 2006: Sato et al.). The γ′ phase is Co₃(Al,W). Design target: σ_y > 700 MPa at 900°C, ρ < 9.5 g/cm³, T_solvus > 1,050°C. (a) CALPHAD prediction for Co-10Al-10W (at%): T_solvus = 1,070°C, γ′ fraction at 900°C = 0.65, density = 9.0 g/cm³. Does it meet targets? (b) DFT predicts Co₃(Al,W) has anti-phase boundary energy γ_APB = 0.26 J/m² (compare to Ni₃Al γ_APB = 0.17 J/m²). Higher APB energy means dislocations cut γ′ with HIGHER energy penalty. Qualitatively, does this give STRONGER or WEAKER resistance to deformation? (c) Phase field simulation of γ′ coarsening at 900°C shows K (Co-Al-W) = 2×K(CMSX-4). How does this affect long-term stability? (d) The alloy has T_melt ≈ 1,440°C (lower than Ni alloys ~1,380°C after considering all elements). This means T_solvus/T_melt ratio = 1,070/1,440 = 0.74. For CMSX-4: 1,277/1,380 = 0.93. What does a lower ratio mean for the processing window for solution treatment? Is this favorable or unfavorable?

5. Machine learning for fatigue prediction: A Random Forest model is trained on a dataset of 320 Ni superalloy compositions with measured LCF life (cycles to failure at 650°C, Δε = 0.8%). Features: atomic% of 12 elements (Ni, Cr, Co, W, Re, Mo, Al, Ti, Ta, Nb, Hf, Ru) + γ′ fraction + lattice misfit. (a) The model predicts LCF life for IN718 = 12,000 cycles; measured = 11,500 cycles. Error = 4.3%. For a new alloy not in training set, the model predicts 15,000 cycles. What additional concern exists beyond model accuracy for this extrapolation? (b) Feature importance analysis shows: γ′ fraction (importance: 0.28), Re content (0.22), Co content (0.18), lattice misfit (0.15), other elements (<0.03 each). What does this tell you about the dominant LCF mechanisms? (c) The model is used to screen 50,000 alloy compositions. It identifies 23 with predicted LCF > 20,000 cycles, density < 9 g/cm³, and T_solvus > 1,270°C. But 18 of 23 are predicted to have TCP phases by CALPHAD screening (run separately). How many candidates remain? (d) Why is it important to combine ML predictions with CALPHAD phase stability calculations rather than using ML alone for alloy screening?
