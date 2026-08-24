# Chapter 71: Quantum Theory

> **"Quantum theory is the most successful physical theory ever devised — its predictions agree with experiment to 12 decimal places. Yet it upended everything we thought we knew about reality, determinism, and the nature of measurement itself."**

---

## Table of Contents

- [71.1 The Failure of Classical Physics](#711-the-failure-of-classical-physics)
- [71.2 Planck's Quantum Hypothesis](#712-plancks-quantum-hypothesis)
- [71.3 The Photoelectric Effect](#713-the-photoelectric-effect)
- [71.4 Photons and Energy](#714-photons-and-energy)
- [71.5 Atomic Spectra and Bohr's Model](#715-atomic-spectra-and-bohrs-model)
- [71.6 The de Broglie Wavelength](#716-the-de-broglie-wavelength)
- [71.7 The Heisenberg Uncertainty Principle](#717-the-heisenberg-uncertainty-principle)
- [71.8 The Schrödinger Equation and Wave Functions](#718-the-schrödinger-equation-and-wave-functions)
- [71.9 Quantum Tunneling](#719-quantum-tunneling)
- [71.10 Implications and the Copenhagen Interpretation](#7110-implications-and-the-copenhagen-interpretation)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 71.1 The Failure of Classical Physics

By 1900, classical (Newtonian) physics was enormously successful. But three experiments it couldn't explain led to a revolution.

### Problem 1: Blackbody Radiation

Every hot object emits radiation. Classical theory (the Rayleigh-Jeans law) predicted the intensity at high frequencies should be infinite — the "ultraviolet catastrophe."

```
CLASSICAL PREDICTION vs EXPERIMENT:

Intensity
   |               classical: goes to infinity!
   |              /
   |             /
   |       /\   /  (experiment: peaks and drops)
   |      /  \_/
   |     /
   |    /
   +-------------> frequency
```

Experimentally, hot objects emit a specific curve (Planck curve) that peaks at a characteristic frequency and drops off at high frequencies. Classical physics got this completely wrong.

### Problem 2: Photoelectric Effect

When light shines on a metal surface, electrons are ejected. But:
- **Below a threshold frequency, no electrons are ejected no matter how bright the light**
- **Above the threshold, electrons are ejected instantly, even if light is dim**
- **Brighter light ejects more electrons, but doesn't give them more energy**

Classical wave theory predicted brighter light should always eject electrons with more energy. It couldn't explain the threshold or the instantaneous ejection.

### Problem 3: Atomic Line Spectra

Heated hydrogen gas emits only specific wavelengths (a discrete line spectrum), not a continuous rainbow. Classical theory predicted it should emit all wavelengths.

```
HYDROGEN LINE SPECTRUM:
  
  Expected (classical): continuous rainbow
  ─────────────────────────────────────────
  
  Observed (experiment): only specific lines
  
  | |  | |  |
  H_α H_β H_γ ...
  
  (specific wavelengths only: 656 nm, 486 nm, 434 nm...)
  
  WHY??? Classical physics had no answer.
```

---

## 71.2 Planck's Quantum Hypothesis

In 1900, Max Planck solved the blackbody problem by making a radical assumption:

> Energy is not continuous. It comes in discrete chunks (quanta).

```
PLANCK'S HYPOTHESIS:

Energy of a quantum of radiation:
  E = h × f

where:
  h = Planck's constant = 6.626 × 10⁻³⁴ J·s
  f = frequency of oscillation
```

At high frequencies, hf is so large that thermal energy kT can't supply it — so high-frequency oscillations don't get excited. This naturally explained the drop-off at high frequencies.

Planck himself thought it was a mathematical trick. Einstein showed it was physical reality.

---

## 71.3 The Photoelectric Effect

In 1905, Einstein explained the photoelectric effect by treating light as consisting of discrete particles (**photons**), each with energy E = hf.

```
PHOTOELECTRIC EXPLANATION:

  Each electron is ejected by a single photon.
  If hf < work function Φ → not enough energy to free electron
  If hf > Φ → electron freed, excess becomes kinetic energy
  
  KE_max = hf - Φ

  Φ = work function = energy needed to escape the surface (depends on metal)
  
  Threshold frequency: f₀ = Φ/h  (minimum frequency to eject electrons)
```

```
GRAPH: Stopping Voltage vs Frequency

V_stop
   |         /
   |        /  slope = h/e
   |       /
  -Φ/e   /
   |-----X-----------> f
         f₀ (threshold)
```

This explained ALL the puzzling observations:
- Below threshold: hf < Φ, no energy to escape no matter how bright
- Brightness → more photons → more electrons ejected, but same KE per electron
- Instantaneous: one photon, one electron — no need to accumulate energy

Einstein won the Nobel Prize for this in 1921.

### Worked Example 71.1

Sodium has work function Φ = 2.3 eV. Find:
(a) The threshold frequency.
(b) The kinetic energy of electrons ejected by 400 nm light.

**Solution:**

(a) f₀ = Φ/h = (2.3 × 1.6×10⁻¹⁹) / (6.626×10⁻³⁴) = 3.68×10⁻¹⁹ / 6.626×10⁻³⁴ = **5.56 × 10¹⁴ Hz**

(b) f = c/λ = 3×10⁸ / 400×10⁻⁹ = 7.5×10¹⁴ Hz

    E_photon = hf = 6.626×10⁻³⁴ × 7.5×10¹⁴ = 4.97×10⁻¹⁹ J = **3.1 eV**

    KE_max = hf - Φ = 3.1 - 2.3 = **0.8 eV**

---

## 71.4 Photons and Energy

Light (all EM radiation) consists of photons — massless particles of electromagnetic energy.

```
PHOTON PROPERTIES:
  Energy:    E = hf = hc/λ
  Momentum:  p = E/c = hf/c = h/λ
  Mass:      0 (massless)
  Speed:     c = 3×10⁸ m/s (always, in vacuum)
```

The energy of a photon depends only on frequency:
- Radio wave photon: ~10⁻⁷ eV (very low energy)
- Visible light photon: ~2-3 eV
- X-ray photon: ~10,000 eV = 10 keV

---

## 71.5 Atomic Spectra and Bohr's Model

In 1913, Niels Bohr explained hydrogen's line spectrum by postulating:

1. Electrons orbit the nucleus in specific, allowed orbits at fixed energies
2. Electrons don't radiate while in these orbits
3. An electron can jump between orbits by absorbing or emitting a photon with exactly the right energy

```
BOHR MODEL OF HYDROGEN:

       n=3 ─────────────────────────  E₃ = -1.5 eV
       n=2 ─────────────────────────  E₂ = -3.4 eV
       
  photon emitted when electron
  drops from n=3 to n=2:
  E_photon = E₃ - E₂ = -1.5 - (-3.4) = 1.9 eV
  λ = hc/E = 656 nm  (red light!)
  
       n=1 ─────────────────────────  E₁ = -13.6 eV (ground state)
```

### Energy Levels in Hydrogen

```
E_n = -13.6 eV / n²   (n = 1, 2, 3, ...)

n=1: -13.6 eV (ground state)
n=2: -3.4 eV
n=3: -1.5 eV
n=∞: 0 eV (ionized = free electron)
```

Photon emitted when electron drops from n₂ to n₁:

```
E_photon = E_n₂ - E_n₁ = 13.6(1/n₁² - 1/n₂²) eV

λ = hc/E_photon
```

This perfectly explained all of hydrogen's spectral lines.

---

## 71.6 The de Broglie Wavelength

In 1924, Louis de Broglie proposed that if light (waves) can behave as particles, then particles can behave as waves.

```
DE BROGLIE WAVELENGTH:

λ = h / p = h / (mv)

where:
  λ = wavelength of matter wave
  h = Planck's constant
  p = momentum = mv
```

For macroscopic objects, λ is unimaginably small:
- Tennis ball (0.06 kg, 50 m/s): λ = 6.626×10⁻³⁴ / (0.06×50) = 2×10⁻³⁴ m (unmeasurable)

For electrons:
- Electron in hydrogen: λ ≈ 0.33 nm (comparable to atomic size!)
- This explains WHY only certain Bohr orbits are allowed: the electron wave must fit a whole number of wavelengths around the orbit!

```
BOHR ORBIT AS STANDING WAVE:

  n=1: one wavelength fits around orbit
  n=2: two wavelengths fit
  n=3: three wavelengths fit
  
  n=2.5: doesn't fit → not an allowed orbit
  
  2πr = nλ = nh/mv → Bohr's quantum condition!
```

---

## 71.7 The Heisenberg Uncertainty Principle

Werner Heisenberg (1927) showed there is a fundamental limit to how precisely we can simultaneously know certain pairs of properties:

```
HEISENBERG UNCERTAINTY PRINCIPLE:

Δx × Δp ≥ ħ/2

ΔE × Δt ≥ ħ/2

where:
  ħ = h/(2π) ≈ 1.055 × 10⁻³⁴ J·s  ("h-bar")
  Δx = uncertainty in position
  Δp = uncertainty in momentum
  ΔE = uncertainty in energy
  Δt = uncertainty in time
```

This is NOT because our measurement tools are imprecise. It is a fundamental property of nature.

### Physical Reason

A particle is described by a wave. To know its position precisely, you need a narrow, localized wave packet — but a narrow packet requires many different wavelengths superimposed, meaning a wide spread of momenta. Knowing position precisely → unknowing momentum precisely.

```
WAVE PACKET:

Spread out wave packet:          Narrow wave packet:
(precise momentum,               (precise position,
 uncertain position)              uncertain momentum)

 ~~~~~~~~~~~~~~                    ~~
 
Many wavelengths superimposed → position narrow → but many k-values → Δp large
```

### Worked Example 71.2

An electron is confined to a region of width Δx = 0.1 nm (roughly an atom).

Find the minimum uncertainty in its velocity.

**Solution:**

Δp ≥ ħ/2Δx = 1.055×10⁻³⁴ / (2 × 0.1×10⁻⁹) = 5.275×10⁻²⁵ kg·m/s

Δv = Δp/m = 5.275×10⁻²⁵ / 9.11×10⁻³¹ ≈ **5.8 × 10⁵ m/s**

This is 0.2% of the speed of light — electrons in atoms are inherently "fuzzy."

---

## 71.8 The Schrödinger Equation and Wave Functions

In 1926, Erwin Schrödinger formulated the full quantum wave equation:

```
SCHRÖDINGER EQUATION (time-dependent):

  iħ × ∂Ψ/∂t = -(ħ²/2m) ∂²Ψ/∂x² + V(x)Ψ

where:
  Ψ (psi) = the wave function of the particle
  V(x) = potential energy
  i = imaginary unit (√-1)
  
(For one dimension; extends to 3D)
```

### The Wave Function

The wave function Ψ contains all information about the quantum system.

The **probability** of finding the particle in a small region dx:

```
P(x) dx = |Ψ(x)|² dx

|Ψ|² = probability density
```

Ψ itself is complex-valued (involves imaginary numbers). Only |Ψ|² has physical meaning — it gives the probability of finding the particle at each location.

```
ELECTRON IN HYDROGEN ATOM:

  |Ψ|²: not a sharp orbit, but a "probability cloud"
  
  Densest where electron most likely to be found.
  
  |Ψ|² for 1s orbital:
  
  Probability
    |*
    | *
    |  *
    |   **
    |      ****
    |          ********
    +-------------------> r
    
  Most likely to find electron near nucleus (a₀ = 0.053 nm)
```

---

## 71.9 Quantum Tunneling

A classical particle cannot pass through a barrier higher than its total energy — it would need negative kinetic energy, which is impossible.

But a quantum particle's wave function extends into and through the barrier. There is a non-zero probability of finding the particle on the other side.

```
QUANTUM TUNNELING:

Classical:
  [particle →]  [||||barrier||||]  no way through

Quantum:
  [particle →]  [barrier]
                 Ψ decays inside   Ψ survives on other side!
                 
  There's a small but nonzero probability of the particle appearing on the far side.
```

The tunneling probability depends exponentially on:
- Barrier width (thinner barrier → more tunneling)
- Barrier height vs particle energy (closer to barrier height → more tunneling)
- Particle mass (lighter particles tunnel more easily)

### Applications

- **Radioactive alpha decay**: alpha particle tunnels out of the nucleus
- **Scanning Tunneling Microscope (STM)**: images individual atoms using electron tunneling current
- **Flash memory**: stores data using electron tunneling through an insulating layer
- **Nuclear fusion in stars**: protons tunnel through their mutual Coulomb repulsion
- **Tunnel diode**: electronic device using quantum tunneling for very fast switching

---

## 71.10 Implications and the Copenhagen Interpretation

Quantum mechanics raises deep philosophical questions.

### The Measurement Problem

Before measurement, a quantum system exists in a **superposition** of multiple states simultaneously. Measuring it "collapses" the wave function to a single definite outcome — but why? And what constitutes a "measurement"?

```
SCHRÖDINGER'S CAT (thought experiment):

  Cat in box with quantum radioactive atom and poison vial.
  If atom decays → releases poison → cat dies.
  
  Before opening box:
    atom = superposition of (decayed + undecayed)
    cat = superposition of (alive + dead)?
    
  After opening box (measuring): cat is definitively alive or dead.
  
  Quantum mechanics says the superposition is real, not just "we don't know."
  This is deeply strange.
```

### The Copenhagen Interpretation

The mainstream interpretation (Bohr, Heisenberg):
- The wave function is not a real physical wave — it represents our knowledge
- Measurement causes collapse to a definite state
- It makes no sense to ask what the particle was doing before measurement
- Focus on what can be measured; don't demand a deeper reality

### Other Interpretations

- **Many-worlds**: every quantum measurement branches reality into multiple worlds
- **Pilot wave (de Broglie-Bohm)**: particles are real, guided by a real wave
- **QBism**: wave function is personal belief, not objective reality

All interpretations give the same experimental predictions — the choice is philosophical.

---

## Summary

- Classical physics failed at blackbody radiation, photoelectric effect, atomic spectra
- **Planck**: energy is quantized; E = hf; h = 6.626 × 10⁻³⁴ J·s
- **Photoelectric effect**: explained by Einstein using photons; KE = hf - Φ
- **Bohr model**: electrons in fixed energy orbits; E_n = -13.6/n² eV for hydrogen; photon emitted on transition
- **de Broglie wavelength**: λ = h/p = h/mv; matter has wave nature
- **Heisenberg uncertainty**: Δx·Δp ≥ ħ/2; position and momentum cannot both be known precisely
- **Schrödinger equation**: wave function Ψ; |Ψ|² = probability density
- **Quantum tunneling**: particles can penetrate classically forbidden barriers
- Interpretation debates continue; Copenhagen is mainstream but alternatives exist

---

## Key Equations

```
Planck's quantum:
  E = hf = hc/λ
  h = 6.626 × 10⁻³⁴ J·s
  ħ = h/2π = 1.055 × 10⁻³⁴ J·s

Photoelectric effect:
  KE_max = hf - Φ
  f₀ = Φ/h  (threshold frequency)
  eV_stop = hf - Φ

Hydrogen energy levels:
  E_n = -13.6 eV / n²
  ΔE = 13.6(1/n₁² - 1/n₂²) eV

Photon momentum:
  p = hf/c = h/λ

de Broglie wavelength:
  λ = h/p = h/(mv)

Heisenberg uncertainty:
  Δx × Δp ≥ ħ/2
  ΔE × Δt ≥ ħ/2

Probability:
  P(x)dx = |Ψ(x)|² dx
```
