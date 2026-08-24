# Chapter 72: Wave-Particle Duality

> **"Electrons behave like waves when nobody is watching, but like particles when measured. This isn't a failure of our experiments — it is the deepest truth about the universe."**

---

## Table of Contents

- [72.1 The Double-Slit Experiment](#721-the-double-slit-experiment)
- [72.2 Which-Way Information](#722-which-way-information)
- [72.3 Electrons as Waves](#723-electrons-as-waves)
- [72.4 Photons as Particles](#724-photons-as-particles)
- [72.5 Complementarity](#725-complementarity)
- [72.6 The Davisson-Germer Experiment](#726-the-davisson-germer-experiment)
- [72.7 Electron Diffraction in Practice](#727-electron-diffraction-in-practice)
- [72.8 Probability and Born's Rule](#728-probability-and-borns-rule)
- [72.9 Implications for Our Picture of Reality](#729-implications-for-our-picture-of-reality)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 72.1 The Double-Slit Experiment

The double-slit experiment is the central experiment of quantum mechanics. Richard Feynman called it "the only mystery" of quantum physics — everything else is a variation.

### Classical Waves (Water/Sound)

```
DOUBLE SLIT WITH WAVES:

  Source → [  | | ]  → Screen
               slits
  
  Waves pass through both slits simultaneously.
  Waves from the two slits interfere:
    - Constructive interference: bright bands
    - Destructive interference: dark bands
  
  PATTERN ON SCREEN:
  
  bright | dark | bright | dark | bright
  (alternating bright/dark bands — interference pattern)
```

### Classical Particles (Bullets)

```
DOUBLE SLIT WITH PARTICLES (classical bullets):

  Gun → [  | | ]  → Screen
               slits
  
  Each bullet goes through ONE slit.
  No interference.
  
  PATTERN ON SCREEN:
  
  Two bright bands (one behind each slit).
  No interference pattern.
```

### Electrons: Surprising Result

```
DOUBLE SLIT WITH ELECTRONS:

  Electron gun → [  | | ]  → Screen
  
  RESULT (low intensity, one electron at a time):
  
  After 1 electron:    random dot, somewhere
  After 100 electrons: some dots, random
  After 10,000 electrons: AN INTERFERENCE PATTERN!
  
  bright | dark | bright | dark | bright
  
  Each electron interferes with ITSELF!
  It seems to go through BOTH slits simultaneously.
```

The interference pattern proves electrons behave as waves going through both slits. But each electron hits the screen at a definite point — like a particle. It's both at once.

---

## 72.2 Which-Way Information

What happens if we watch which slit the electron goes through?

```
MODIFIED DOUBLE-SLIT (with detector at slits):

  Electron gun → [Detector|Slit 1 | Slit 2|Detector] → Screen
  
  We observe which slit each electron uses.
  
  RESULT:
  
  Two bands (one behind each slit).
  NO interference pattern!
  
  The act of observing the electron's path DESTROYS the interference.
```

This is the most profound result in all of physics:

> **Observing which path the electron takes changes the outcome. The interference pattern disappears when we know which slit the electron used.**

If we determine the electron's path (particle behavior), we lose the wave interference. If we allow both paths (wave behavior), we cannot know which slit it used.

**This is not about disturbance by the measurement apparatus.** Even the gentlest possible measurement that reveals which-way information destroys the interference. The information itself — not the mechanical disturbance — is what matters.

---

## 72.3 Electrons as Waves

Before measurement (detection at the screen), the electron is described by a wave function that spreads through both slits and interferes with itself.

```
ELECTRON WAVE FUNCTION (spreading):

          /==========\
  source=|     both   |==  → wave spreads through both slits
          \==========/ → interference pattern in probability |Ψ|²
  
  |Ψ|² at the screen:
  
  bright regions: high probability of finding electron
  dark regions:   low probability (nearly zero)
```

The wave function is not a physical wave like sound or water. It's a **probability amplitude wave** — its square gives the probability of finding the particle at each location.

---

## 72.4 Photons as Particles

Conversely, light shows particle behavior:

### Compton Scattering (1923)

Arthur Compton fired X-rays at electrons. The scattered X-rays had a longer wavelength than the incident ones — as if the X-ray photon had collided with an electron and lost energy (and momentum) to it:

```
COMPTON SCATTERING:

  Incident photon (λ₁) + stationary electron
  
       →
  [γ]     [e]  →  collision  →  [γ'] + [e']
  
  Scattered photon (λ₂ > λ₁)    Recoiling electron
  
  Change in wavelength:
  Δλ = λ₂ - λ₁ = (h/m_e c)(1 - cos θ)
  
  This proves the photon has MOMENTUM: p = h/λ
```

The wavelength shift depends on scattering angle — exactly what you'd expect from a particle collision. Waves don't have momentum in this particle-like way.

### Photon in Single-Photon Interference

Even a single photon (one at a time) builds up an interference pattern over many trials — each photon interferes with itself, just like the electron.

---

## 72.5 Complementarity

Niels Bohr's principle of **complementarity** (1927):

> **Wave behavior and particle behavior are complementary aspects of reality. They cannot both be observed simultaneously. Which aspect manifests depends on the type of experiment performed.**

```
WAVE BEHAVIOR:
  - Interference, diffraction
  - Spread over space
  - Can't say where it "is"
  - Both slits at once
  - NO which-way information

PARTICLE BEHAVIOR:
  - Definite location
  - Localized interaction
  - Definite trajectory
  - Goes through ONE slit
  - HAS which-way information

You can have ONE or the OTHER, never both at the same time.
```

This is not a limitation of our instruments. It is a fundamental feature of quantum reality.

---

## 72.6 The Davisson-Germer Experiment

In 1927, Clinton Davisson and Lester Germer directly confirmed de Broglie's wave hypothesis for electrons.

They fired electrons at a nickel crystal and observed the scattered electrons:

```
DAVISSON-GERMER EXPERIMENT:

  Electron beam → Nickel crystal → Detector
  
  The crystal acts as a diffraction grating for electrons.
  
  Expected (classical): scattered uniformly
  Observed: peaks at specific angles (diffraction pattern!)
  
  The angles matched exactly the de Broglie wavelength:
  λ = h/mv
  
  PROOF: electrons have wave properties.
```

The peak at 50° for 54 V electrons gave λ = 0.167 nm — exactly matching the de Broglie prediction.

---

## 72.7 Electron Diffraction in Practice

Electron waves are used in **electron microscopes** precisely because of their wave nature:

```
ELECTRON MICROSCOPE vs LIGHT MICROSCOPE:

  Light microscope: limited by λ_visible ≈ 400-700 nm
    → can't resolve features smaller than ~200 nm
    
  Electron microscope: electrons at 100 keV have λ ≈ 0.004 nm
    → can resolve features down to 0.1 nm = atomic scale!
    
  TEM (Transmission Electron Microscope): electrons pass through thin sample
  SEM (Scanning Electron Microscope): electrons scan across surface
  
  Used to: see viruses, individual protein molecules, crystal defects, atom columns
```

### X-ray Diffraction vs Electron Diffraction

Both are used to determine crystal structures. X-ray diffraction (used to solve DNA structure) exploits photon waves; electron diffraction exploits electron waves. Both work because λ is comparable to atomic spacing (~0.1-0.3 nm).

---

## 72.8 Probability and Born's Rule

Max Born (1926) proposed the correct interpretation of the wave function:

```
BORN'S RULE:

|Ψ(x)|² dx = probability of finding the particle in interval [x, x+dx]

Not a physical matter wave.
Not a charge distribution.
Just probability amplitude.
```

This is why quantum mechanics is probabilistic — even with perfect knowledge of the wave function (and hence everything knowable about the system), we can only predict probabilities, not definite outcomes.

```
EXAMPLE: Single electron heading toward the screen

  Wave function: spread over screen
  
  |Ψ|² at screen:
  
  Probability
    |
  P₁|    *
    |   * *
    |  *   *    probability to land here
  P₂|       *
    |        *
    |          ***
    +-------------> position x
  
  The electron will land at ONE specific spot — but which spot is random.
  All we can predict is the probability distribution |Ψ|².
```

---

## 72.9 Implications for Our Picture of Reality

Wave-particle duality forces us to abandon several classical assumptions:

### 1. Objects Don't Have Definite Properties Before Measurement

Before measuring which slit the electron went through, it didn't go through either slit in a definite sense — it was in a superposition of both paths.

### 2. The Measurement Creates the Outcome

Quantum mechanics doesn't describe what the electron IS doing between measurements; it only predicts probabilities of outcomes WHEN measured.

### 3. Identical Particles Can Be Indistinguishable

Two quantum particles of the same type are truly identical. Swapping them produces no observable change. This leads to new statistics (Bose-Einstein statistics for bosons, Fermi-Dirac statistics for fermions) that differ from classical statistics.

### 4. Non-Locality (Entanglement)

Two particles can be "entangled" — measuring one instantly affects the probability distribution of the other, even if separated by light-years. (But this cannot be used to send information faster than light.)

```
ENTANGLEMENT:
  
  Two electrons created together → wave functions entangled
  
  Measure electron 1 → instantly collapses both wave functions
  Electron 2 now has definite (correlated) properties
  
  But: you can't use this to signal anything, because the
       outcomes are random — you can't control what you measure.
```

---

## Summary

- **Double-slit experiment**: electrons (and all matter) show interference — wave behavior when unobserved
- **Which-way measurement**: knowing path destroys interference — particle behavior manifests
- **Complementarity**: wave and particle behavior are mutually exclusive aspects of the same entity
- **de Broglie wave**: λ = h/mv; confirmed by Davisson-Germer (electron diffraction)
- **Compton scattering**: proves photons have momentum p = h/λ; particle behavior of light
- **Born's rule**: |Ψ|² = probability density; quantum mechanics is fundamentally probabilistic
- **Consequences**: properties not definite before measurement; superposition; entanglement
- Wave-particle duality is not a limitation of experiments but the fundamental structure of reality

---

## Key Equations

```
de Broglie wavelength:
  λ = h / p = h / (mv)
  h = 6.626 × 10⁻³⁴ J·s

Photon momentum:
  p = h / λ = hf / c

Compton scattering:
  Δλ = λ₂ - λ₁ = (h / m_e c)(1 - cos θ)
  h/m_e c = 2.426 × 10⁻¹² m (Compton wavelength)

Born's rule:
  P(x)dx = |Ψ(x)|² dx

Energy in terms of wavelength:
  E = hf = hc/λ

For photon:
  E = pc,  p = h/λ
  
For massive particle:
  E = p²/2m (non-relativistic), p = mv, λ = h/mv
```
