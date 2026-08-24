# Chapter 67: Wave Optics — Interference and Diffraction

> **"Young's double-slit experiment didn't just demonstrate interference — it proved light is a wave. It remains one of the most beautiful experiments in all of physics."**

---

## Table of Contents

- [1. Wave Optics vs Geometric Optics](#1-wave-optics-vs-geometric-optics)
- [2. Coherent Sources — The Secret Ingredient](#2-coherent-sources--the-secret-ingredient)
- [3. Young's Double-Slit Experiment](#3-youngs-double-slit-experiment)
- [4. Constructive Interference — Bright Fringes](#4-constructive-interference--bright-fringes)
- [5. Destructive Interference — Dark Fringes](#5-destructive-interference--dark-fringes)
- [6. Fringe Spacing](#6-fringe-spacing)
- [7. Worked Example 1 — Double-Slit Calculation](#7-worked-example-1--double-slit-calculation)
- [8. Thin Film Interference](#8-thin-film-interference)
- [9. Soap Bubble Colors](#9-soap-bubble-colors)
- [10. Anti-Reflection Coatings](#10-anti-reflection-coatings)
- [11. Diffraction — Waves Bend Around Corners](#11-diffraction--waves-bend-around-corners)
- [12. Single-Slit Diffraction](#12-single-slit-diffraction)
- [13. Diffraction Gratings](#13-diffraction-gratings)
- [14. Worked Example 2 — Diffraction Grating Calculation](#14-worked-example-2--diffraction-grating-calculation)
- [15. CDs and DVDs as Diffraction Gratings](#15-cds-and-dvds-as-diffraction-gratings)
- [16. Applications of Wave Optics](#16-applications-of-wave-optics)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 1. Wave Optics vs Geometric Optics

To understand wave optics, we first need to understand when it matters — and when it does not.

### Geometric Optics (Ray Optics)

In the chapters on mirrors and lenses, we treated light as **rays** — straight lines that bounce off mirrors and bend at lens surfaces. That model works wonderfully when the objects and openings that light interacts with are **much larger than the wavelength of light**.

Visible light has wavelengths between about 380 nm (violet) and 700 nm (red). One nanometre is 0.000 000 001 metres. So a wavelength of 550 nm is roughly 0.00055 mm — impossibly tiny compared to a mirror, a lens, or a window.

When λ << (size of obstacle or opening), geometric optics gives accurate, reliable results. We can ignore the wave nature of light entirely.

### Wave Optics

But here is the thing: light IS a wave. It is an oscillating electromagnetic field that travels through space. When the obstacles or openings become **comparable in size to the wavelength**, the wave nature of light asserts itself dramatically.

Two wave phenomena become important:

- **Interference** — two or more waves overlap and add together (or cancel). If crests meet crests, the result is a brighter wave (**constructive interference**). If crests meet troughs, they cancel (**destructive interference**).
- **Diffraction** — a wave bends around the edges of obstacles and spreads out through narrow openings, instead of traveling in perfectly straight lines.

### An Analogy: Sound vs Light Around Corners

Think about sound. Sound waves have wavelengths of roughly 1 cm to 10 m. You can easily hear someone talking around a corner because sound diffracts around the corner — the wavelength is comparable to the corner's size.

Light has a wavelength of ~500 nm. Everyday corners are millions of times larger. So light does NOT noticeably diffract around a doorframe. It seems to travel in perfectly straight rays — geometric optics describes this well.

BUT — put light through a slit that is only a few micrometres wide, or through two tiny slits separated by a fraction of a millimetre, and suddenly the wave nature of light produces spectacular patterns of light and dark bands on a screen. That is wave optics in action.

---

## 2. Coherent Sources — The Secret Ingredient

Wave interference is real and always happening. But to **see** a stable, visible interference pattern, the sources of light must be **coherent**.

### What is Coherence?

Two light sources are **coherent** if:

1. They have the **same frequency** (same wavelength, same color).
2. They maintain a **constant phase relationship** — the phase difference between them does not randomly fluctuate over time.

Condition 2 is the tricky one. Think of two friends waving their hands up and down. If they wave perfectly in sync — both up at the same moment, both down at the same moment — they are "in phase." If one is always exactly half a wave behind the other, they are "out of phase." The key word is **always**. The relationship must stay constant.

### Why Normal Light Bulbs Don't Produce Interference

A light bulb emits light from millions of atoms. Each atom emits a tiny "wave packet" — a short burst of light lasting only about 10^-8 seconds. These atoms emit randomly, independently of each other. The phase relationship between light from two separate bulbs changes billions of times per second — far faster than any detector or human eye can respond.

The result: the bright and dark regions of any interference pattern flash on and off so rapidly that you only see an average — which is uniform brightness everywhere. The interference washes out.

Two separate light bulbs are **incoherent** sources. They cannot produce a stable interference pattern.

### How to Get Coherent Sources

The classic trick — used by Thomas Young in 1801 — is to take **one** source and split it into two. Since both "sub-sources" come from the same original source, they maintain a fixed phase relationship. That is exactly what the double-slit experiment does.

Modern coherent light sources include **lasers**, which emit light where all photons are in phase. With a laser, you do not even need a double slit to get interference — the coherence length is enormous.

---

## 3. Young's Double-Slit Experiment

In 1801, Thomas Young performed an experiment that settled a long-running debate: is light a particle (Newton's view) or a wave (Huygens's view)? The result was stunning — light produces interference fringes, exactly as a wave would.

### The Setup

```
                      BARRIER WITH             SCREEN
MONOCHROMATIC         TWO SLITS
LIGHT SOURCE
                         |  <- slit S1          |  (bright)
     *  -------->        |                      |  (dark)
   (source)              |                      |  (bright)   <- central maximum
                         |  <- slit S2          |  (dark)
                         |                      |  (bright)
                         |
                    <----D---->
                 (distance from slits to screen)

  Slit separation = d (center-to-center, typically 0.1 mm to 1 mm)
  Screen distance = D (typically 0.5 m to 2 m)
  Wavelength = λ (monochromatic light: single color)
```

Monochromatic (single-wavelength) light hits a barrier with two narrow slits, S1 and S2, separated by a distance **d**. Both slits act as new coherent sources (because they are both fed by the same incoming wave — they are always in phase with each other at the slits).

The two sets of waves fan out from S1 and S2 and overlap in the region between the barrier and the screen. Where they overlap, they interfere.

### Path Difference — The Key Idea

Consider a point P on the screen at height y above the center.

```
                 S1 ----r1----
                  \           \
                   \           \
                    \           P  (on screen, height y above center)
                   /           /
                  /           /
                 S2 ----r2----

  r1 = distance from S1 to P
  r2 = distance from S2 to P
  Path difference = Δ = r2 - r1  (assuming S2 is below center, P is above)
```

For a point at height y on the screen (with the screen at distance D, and slits separated by d):

**Path difference: Δ = d * y / D** (valid for small angles, i.e., y << D)

More precisely, the path difference is Δ = d * sinθ, where θ is the angle from the center of the barrier to the point P. For small angles, sinθ ≈ tanθ ≈ y/D.

The result at P — whether bright or dark — depends entirely on this path difference Δ.

---

## 4. Constructive Interference — Bright Fringes

When the path difference Δ equals a **whole number of wavelengths**, the two waves arrive at P perfectly in phase — crest meets crest, trough meets trough. They reinforce each other, producing maximum brightness.

**Condition for constructive interference (bright fringe):**

    Δ = m * λ,   where m = 0, ±1, ±2, ±3, ...

Here m is called the **order number** or fringe order.

- m = 0: the central bright fringe (at y = 0, on the centerline)
- m = +1, -1: the first-order bright fringes, one above and one below center
- m = +2, -2: the second-order bright fringes, and so on

Substituting Δ = d * y / D:

    d * y_m / D = m * λ

Solving for the position of the m-th bright fringe:

**y_m = m * λ * D / d**

---

## 5. Destructive Interference — Dark Fringes

When the path difference is a **half-integer number of wavelengths**, the two waves arrive exactly out of phase — crest meets trough. They cancel completely, producing a dark region.

**Condition for destructive interference (dark fringe):**

    Δ = (m + 1/2) * λ,   where m = 0, ±1, ±2, ±3, ...

This gives path differences of λ/2, 3λ/2, 5λ/2, and so on.

Position of the m-th dark fringe:

**y_m = (m + 1/2) * λ * D / d**

For m = 0: y = λD/(2d) — first dark fringe just above center
For m = -1: y = -λD/(2d) — first dark fringe just below center

The result is alternating bright and dark bands — called **interference fringes** — across the screen.

---

## 6. Fringe Spacing

The distance between two adjacent bright fringes (or two adjacent dark fringes) is constant across the screen (for small angles). Let us compute it:

Position of m-th bright fringe:   y_m = m * λD/d
Position of (m+1)-th bright fringe:  y_(m+1) = (m+1) * λD/d

Spacing:

    Δy = y_(m+1) - y_m = (m+1) * λD/d - m * λD/d

**Fringe spacing: Δy = λD/d**

Key observations:
- Increase wavelength λ → wider fringe spacing (red light gives wider fringes than blue)
- Increase slit separation d → narrower fringes (slits closer together → wider fringes)
- Increase screen distance D → wider fringes

This formula is one of the most useful results in wave optics. It allows us to measure the wavelength of light with a ruler!

---

## 7. Worked Example 1 — Double-Slit Calculation

**Problem:** Light of wavelength 550 nm passes through two slits separated by 0.25 mm. The screen is placed 1.2 m away.

**(a)** Find the fringe spacing Δy.
**(b)** Find the position of the 3rd bright fringe from the center.
**(c)** Find the position of the 2nd dark fringe above the center.

---

**Given information — first, convert to consistent units (metres):**

    λ = 550 nm = 550 × 10^-9 m = 5.50 × 10^-7 m
    d = 0.25 mm = 0.25 × 10^-3 m = 2.50 × 10^-4 m
    D = 1.2 m

---

**Part (a): Fringe spacing**

Using the formula Δy = λD/d:

    Δy = (5.50 × 10^-7 m × 1.2 m) / (2.50 × 10^-4 m)

    Numerator: 5.50 × 10^-7 × 1.2 = 6.60 × 10^-7 m²

    Δy = (6.60 × 10^-7) / (2.50 × 10^-4)
       = (6.60 / 2.50) × 10^(-7-(-4))
       = 2.64 × 10^-3 m
       = 2.64 mm

**Answer (a): The fringe spacing is 2.64 mm.**

This means every 2.64 mm along the screen, there is another bright fringe. With a screen distance of 1.2 m, about 450 bright fringes would be visible (across a 1.2 m screen), though in practice the outer fringes get dimmer.

---

**Part (b): Position of the 3rd bright fringe (m = 3)**

Using y_m = mλD/d:

    y_3 = 3 × λD/d = 3 × Δy

    y_3 = 3 × 2.64 × 10^-3 m
        = 7.92 × 10^-3 m
        = 7.92 mm

**Answer (b): The 3rd bright fringe is 7.92 mm from the center.**

(This makes sense — it is 3 fringe spacings away from center, which is 3 × 2.64 = 7.92 mm. Good check!)

---

**Part (c): Position of the 2nd dark fringe above center (m = 1)**

Wait — let us be careful about the indexing. The formula is:

    y_dark = (m + 1/2) * λD/d

For the 1st dark fringe above center: m = 0, giving y = (0 + 1/2)λD/d = 0.5 × Δy = 1.32 mm
For the 2nd dark fringe above center: m = 1, giving y = (1 + 1/2)λD/d = 1.5 × Δy

    y = 1.5 × 2.64 × 10^-3 m
      = 3.96 × 10^-3 m
      = 3.96 mm

**Answer (c): The 2nd dark fringe above center is 3.96 mm from the center.**

Notice it sits exactly halfway between the 1st bright fringe (2.64 mm) and the 2nd bright fringe (5.28 mm). Dark fringes always sit midway between bright fringes. Consistent!

---

## 8. Thin Film Interference

Have you ever seen the swirling rainbow of colors on a soap bubble, a puddle of oil on wet pavement, or a camera lens? That is **thin film interference** — one of the most beautiful and practical consequences of wave optics.

### The Basic Mechanism

When light hits a thin transparent film (like a soap film or an oil slick on water), it partially reflects from the **top surface** and partially transmits through to the **bottom surface**, where it partially reflects again.

```
         Incident light
              |
              v
   ~~~~~~~~~~~~~~~~~~~~~~~~~~~  <-- top surface of film (medium 1 to medium 2)
   |         Film             |     (refractive index n_film)
   ~~~~~~~~~~~~~~~~~~~~~~~~~~~  <-- bottom surface (medium 2 to medium 3)
              |
         Transmitted
```

The ray reflected from the top surface and the ray reflected from the bottom surface both come back upward — and they travel different distances. The bottom-surface ray travels an extra distance of **2t** (down through the film and back up), where t is the film thickness.

So the path difference is approximately: **Δ = 2 * n_film * t** (the factor n_film accounts for the shorter wavelength inside the film — more on this shortly).

Actually, for light hitting at near-normal incidence (straight on), the optical path difference is:

    Δ = 2 * n * t

where n is the refractive index of the film.

### The Phase Shift Complication

Here is the subtle part that many beginners miss. When a light wave reflects off a **denser medium** (higher refractive index), it undergoes a **180° phase shift** — equivalent to adding half a wavelength to the path. When it reflects off a **less dense medium** (lower refractive index), there is **no phase shift**.

Think of it like a transverse wave on a string:
- String tied to a fixed wall → wave reflects inverted (180° shift)
- String connected to a lighter, freely moving rope → wave reflects upright (no shift)

This phase shift rule is crucial for thin films:

| Reflection at boundary | Phase shift? |
|------------------------|--------------|
| Going from low-n to high-n (e.g., air to glass) | YES, 180° shift |
| Going from high-n to low-n (e.g., glass to air) | NO shift |

### Applying to a Soap Film in Air (n_air = 1, n_soap ≈ 1.33)

- Top surface: air → soap (low-n to high-n) → **180° phase shift**
- Bottom surface: soap → air (high-n to low-n) → **no phase shift**

Net phase shifts: **one** of the two reflected rays has a 180° phase shift. That is like one ray getting an extra half-wavelength of path difference.

So the effective path difference = 2nt + λ/2 (from the phase shift)

**Conditions for a soap film:**

Constructive interference (bright, for wavelength λ in air):
    2 * n * t = (m + 1/2) * λ,   m = 0, 1, 2, ...

Destructive interference (dark):
    2 * n * t = m * λ,   m = 0, 1, 2, ...

Note the reversal from what you might expect! The phase shift "flips" the conditions.

For m = 0 (t ≈ 0, extremely thin film): 2nt = 0 which is a whole number (0) of wavelengths → **destructive** → the film looks **dark**, not bright. This is why the very top of a soap bubble (thinnest part) appears black just before it bursts.

---

## 9. Soap Bubble Colors

A soap bubble is a thin spherical film. The bubble is thicker near the equator and bottom (gravity pulls the soapy water down) and thinner near the top.

Different thicknesses t satisfy the constructive interference condition for different wavelengths:

    2 * n * t = (m + 1/2) * λ_constructive

At each point on the bubble surface, the thickness t is fixed. Different wavelengths (colors) of white light hit that spot. Only the wavelengths that satisfy the constructive interference condition reflect strongly — those are the colors you see.

As t changes across the bubble surface:
- Very thin film (top): all visible wavelengths are suppressed → **appears black or dark**
- Slightly thicker: might constructively reinforce blue-green → appears **blue-green**
- Thicker still: reinforces yellow-red → appears **golden**
- Even thicker: reinforces red, then loops back to violet, then second-order blue, etc.

The result is concentric rings of color — a gorgeous interference rainbow.

The same phenomenon occurs with an oil film on water (n_oil ≈ 1.46, n_water ≈ 1.33). Here the phase shifts at the two surfaces are different: air→oil gives a phase shift, but oil→water also gives a phase shift (since n_water < n_oil? No — n_water ≈ 1.33 < n_oil ≈ 1.46, so oil→water is high-n to low-n → no phase shift). Wait, let us redo:

- Air (n=1) → Oil (n≈1.46): low to high → **180° phase shift**
- Oil (n≈1.46) → Water (n≈1.33): high to low → **no phase shift**

Same situation as soap film in air: one phase shift. Same formulas apply.

---

## 10. Anti-Reflection Coatings

Camera lenses, eyeglasses, and binoculars are often coated with a thin layer of material to **reduce reflection**. This is a direct application of thin film interference.

### The Problem

Plain glass (n ≈ 1.5) reflects about 4% of incident light at each surface. A camera lens with 10 glass-air surfaces can lose 30-40% of light to reflections, reducing image brightness and creating internal glare.

### The Solution

Coat the glass with a thin film of **magnesium fluoride (MgF2)** with n ≈ 1.38.

Now we have three layers: air (n=1.00) → MgF2 (n=1.38) → glass (n=1.50)

Phase shifts:
- Air → MgF2: low to high (1.00 → 1.38) → **180° phase shift** (ray 1)
- MgF2 → glass: low to high (1.38 → 1.50) → **180° phase shift** (ray 2)

Both reflected rays get a 180° phase shift — they cancel each other out in terms of the phase-shift effect. The net effect of the phase shifts: **zero extra path** (both shifted equally, net difference = 0).

So the condition for destructive interference in the coating (to kill reflection) becomes:

    2 * n_coating * t = (m + 1/2) * λ

For m = 0 (thinnest useful coating):

    t = λ / (4 * n_coating)

This is called a **quarter-wave coating** (the film is one quarter of the wavelength thick, measured inside the film material).

**Example calculation:**

For λ = 550 nm (green light, middle of the visible spectrum) and n = 1.38:

    t = 550 nm / (4 × 1.38)
      = 550 / 5.52
      ≈ 99.6 nm
      ≈ 100 nm

That is about 100 atoms thick. Modern coating technology deposits films with this precision routinely.

### Why the Lens Looks Purple-Blue

The coating is optimized for green (λ ≈ 550 nm). But the red and blue ends of the spectrum are not perfectly cancelled. Some red and blue light still reflects, while green is suppressed. Red + blue = purple/magenta. That is why anti-reflection coated lenses have a characteristic **blue-purple sheen** when you look at them — it is constructive interference of the imperfectly cancelled ends of the spectrum.

---

## 11. Diffraction — Waves Bend Around Corners

**Diffraction** is the bending and spreading of waves when they encounter an obstacle or pass through an opening.

### Huygens' Principle

Every point on a wavefront can be treated as a **source of new spherical wavelets**. The new wavefront is the envelope of all these wavelets.

When a wave passes through a narrow slit, the edges of the slit act as new sources. The wavelets spread out in all directions — including sideways. The wave does not just go straight through; it fans out. The narrower the slit (relative to λ), the more the wave spreads.

### Why Sound Diffracts Easily But Light Seems Not To

Sound in a room: λ ≈ 0.1 to 3 m. A doorway is about 1 m wide. Since λ ~ doorway width, sound diffracts significantly around corners — you can hear through doorways even when you can't see through them.

Light: λ ≈ 500 nm = 5 × 10^-7 m. A doorway is 1 m = 10^6 nm wide. Since λ << doorway width, light barely diffracts at all — it essentially goes straight through like rays.

BUT — make the opening only a few micrometres wide (comparable to λ), and light diffracts dramatically. This is the regime of single-slit diffraction.

### Everyday Diffraction of Light

Even with "big" openings, you can see diffraction effects:

- **Halos around streetlights** on a foggy night — tiny water droplets diffract the light
- **Iridescence on a butterfly wing** — microscopic structures create diffraction patterns
- **The colors of a CD** — the spiral track grooves act as a diffraction grating (more on this soon)
- **Squinting your eyes** — the narrow gap between your eyelashes creates a diffraction pattern, making bright lights streak

---

## 12. Single-Slit Diffraction

When a wave passes through a single slit of width **a**, it does not produce a sharp shadow. Instead, it produces a **diffraction pattern** on a screen — a wide central bright region flanked by weaker maxima on either side, with dark minima between them.

### The Pattern

```
   Intensity
       |
       |         *****
       |        *     *
       |       *       *
       |      *         *
       |*****             *****           *****
       +--+--+--+--+--+--+--+--+--+--+--+--> Position on screen
         -3  -2  -1   0   1   2   3    (in units of λ/a)
              
   Central maximum: very bright, width ≈ 2λ/a
   First secondary maxima: much dimmer (about 4.7% of central peak)
   Second secondary maxima: even dimmer (about 1.6% of central peak)
```

### Dark Fringes in Single-Slit Pattern

Dark fringes (minima) occur at angles:

**sin(θ) = m * λ / a,   m = ±1, ±2, ±3, ...**  (but NOT m = 0)

Notice: m = 0 is NOT a minimum — it is the center of the central maximum.

The **angular half-width** of the central maximum (from center to first dark fringe) is:

    θ_1 = arcsin(λ/a)  ≈  λ/a  (for small angles)

A narrower slit (smaller a) gives a wider central maximum. This is the wave nature of light in action — forcing light through a narrower gap causes it to spread out MORE, not less. It seems counterintuitive but is a fundamental property of waves.

### Position of Dark Fringes on Screen

For small angles (y << D):

    y_dark = m * λ * D / a,   m = ±1, ±2, ...

### Comparison: Single-Slit vs Double-Slit Patterns

The double-slit pattern (rapid, equal-brightness fringes) is actually **modulated** by the single-slit envelope. In a real double-slit experiment, the individual slits have finite width a, and the overall intensity is the double-slit interference pattern multiplied by the single-slit diffraction envelope. Fringes near the edges are dimmer than fringes near the center because of this modulation.

---

## 13. Diffraction Gratings

A **diffraction grating** is a surface with many equally spaced parallel slits (or grooves). Typical gratings have 300 to 1200 lines per millimetre — so thousands of slits per centimeter.

### Why More Slits Is Better

Recall the double-slit experiment: two coherent sources produce broad, blurry bright fringes. Add a third slit — the pattern sharpens. Add a fourth — sharper still. With thousands of slits, the bright maxima become extremely sharp, narrow, bright lines, with almost completely dark regions in between.

This is the power of the diffraction grating: it separates light of different wavelengths into very distinct, sharp **spectral lines**, making it an ideal tool for precision spectroscopy.

### The Grating Equation

For a grating with slit spacing **d** (also called the grating constant), the **principal maxima** (bright lines) appear at angles satisfying:

**d * sin(θ) = m * λ,   m = 0, ±1, ±2, ±3, ...**

This looks just like the double-slit condition! But the pattern is much sharper.

- m = 0: the **zeroth-order** maximum (straight through, all wavelengths overlap here — it looks white for white light)
- m = ±1: the **first-order** maxima (one on each side)
- m = ±2: the **second-order** maxima, and so on

The grating spacing d is related to the number of lines per metre N by:

    d = 1 / N

For a grating with 600 lines/mm:
    N = 600 lines/mm = 600,000 lines/m
    d = 1 / 600,000 m = 1.667 × 10^-6 m = 1667 nm

### Grating as a Spectrum Splitter

```
                              m = +2  (blue closer to center than red)
                         m = +1
                         
   White light -----> | GRATING | -----> m = 0  (white, all colors together)
                         
                         m = -1
                              m = -2
   
   For each order m ≠ 0, different wavelengths appear at different angles.
   Blue (short λ) bends less than red (long λ) — opposite of a prism!
```

Wait — let us verify that "blue bends less." From d sinθ = mλ, for fixed m and d:

    sinθ = mλ/d

Larger λ → larger sinθ → larger θ. So red (λ ≈ 700 nm) bends MORE than violet (λ ≈ 400 nm) in a grating. This is opposite to a prism, where blue bends more due to dispersion.

### The Grating Spectrometer

A grating spectrometer shines light from an unknown source through a narrow slit onto a grating. Each element in the source (hydrogen, sodium, mercury, etc.) emits specific wavelengths — its **spectral fingerprint**. The grating separates these wavelengths spatially, allowing measurement of each. From the angle θ and the known grating spacing d, you can compute λ precisely.

This is how astronomers determine the chemical composition of distant stars — from the absorption lines in the star's spectrum.

---

### ASCII Diagram: Diffraction Grating

```
                  |  m = +3  (third order)
                  |
              ~~~~|~~~~  second order (m = +2)
           ~~~    |    ~~~
        ~~~       |       ~~~   m = +1 (first order)
     ~~~          |          ~~~
====|=============|=============|====  m = 0 (zeroth order, straight through)
     ~~~          |          ~~~
        ~~~       |       ~~~   m = -1
           ~~~    |    ~~~
              ~~~~|~~~~  m = -2
                  |
                  |  m = -3

   GRATING (vertical bar) with thousands of slits
   Each "~~~" line represents an angle θ for a principal maximum
   Different wavelengths (colors) appear at slightly different angles within each order
```

---

## 14. Worked Example 2 — Diffraction Grating Calculation

**Problem:** A diffraction grating has 600 lines/mm. When illuminated with monochromatic light, a 2nd-order maximum is observed at an angle of θ = 42.5°. Find the wavelength of the light.

---

**Step 1: Find the grating spacing d.**

    Number of lines per mm = 600
    Number of lines per metre = 600 × 1000 = 600,000 = 6.00 × 10^5 lines/m
    
    d = 1 / (6.00 × 10^5)
    d = 1.667 × 10^-6 m
    d = 1667 nm

---

**Step 2: Apply the grating equation.**

    d * sin(θ) = m * λ
    
    Given: θ = 42.5°, m = 2
    
    sin(42.5°) = 0.6756  (use a calculator or table)
    
    λ = d * sin(θ) / m
    λ = (1.667 × 10^-6 m) × 0.6756 / 2
    λ = (1.667 × 0.6756 / 2) × 10^-6 m
    λ = (1.1262 / 2) × 10^-6 m
    λ = 0.5631 × 10^-6 m
    λ = 563.1 nm

---

**Step 3: Check if this is physically reasonable.**

563 nm is in the green-yellow part of the visible spectrum. Visible light runs from about 380 nm (violet) to 700 nm (red). So 563 nm is plausible. ✓

Also check: could there be a 3rd-order maximum for this λ?

    sin(θ_3) = 3λ/d = 3 × 563 / 1667 = 1689/1667 = 1.013

sinθ cannot exceed 1, so no 3rd-order maximum exists for this wavelength and grating. ✓ (It would exist for shorter wavelengths, but just barely misses for 563 nm.)

**Answer: The wavelength is approximately 563 nm (green-yellow light).**

---

### Bonus: Can We Find the 1st-Order Maximum Angle?

    sin(θ_1) = 1 × λ / d = 563 nm / 1667 nm = 0.3378
    θ_1 = arcsin(0.3378) = 19.7°

So the 1st-order maximum is at about 19.7° and the 2nd-order is at 42.5°. Makes sense — the angle increases with order number (but not linearly, because of the arcsin).

---

## 15. CDs and DVDs as Diffraction Gratings

Have you ever held a CD or DVD in sunlight and seen a brilliant rainbow? That is diffraction — and the disc is acting as a **reflection diffraction grating**.

### How CD Tracks Work as a Grating

A CD stores data as microscopic pits and bumps in a spiral track. The tracks are arranged in a tight spiral, but adjacent tracks are separated by a precise distance:

    CD track spacing: d ≈ 1.6 µm = 1600 nm
    DVD track spacing: d ≈ 0.74 µm = 740 nm

These spacings are comparable to the wavelengths of visible light (380–700 nm). So when white light hits the disc, different wavelengths satisfy the grating equation d sinθ = mλ at different angles → rainbow of colors.

### Verification

For a CD (d = 1600 nm) and green light (λ = 550 nm), the first-order maximum angle:

    sinθ = λ/d = 550/1600 = 0.344
    θ = 20.1°

For red light (λ = 650 nm):

    sinθ = 650/1600 = 0.406
    θ = 24.0°

So red and green appear at different angles — you see them as different colors in the reflected rainbow. Blue is at a smaller angle, red at a larger angle (opposite to a prism, as noted earlier).

DVDs have finer track spacing (740 nm vs 1600 nm), which pushes the diffraction angles larger. DVDs still show rainbow colors but the pattern is slightly different. Blu-ray has even finer tracks (320 nm) — so only violet light (λ ≈ 400 nm) can satisfy the grating equation for first order (sinθ = 400/320 > 1 for most visible wavelengths!). Blu-ray discs look less colorful in white light than CDs.

---

## 16. Applications of Wave Optics

Wave optics is not just a physics curiosity — it underpins many important technologies and scientific tools.

### Spectroscopy

Diffraction gratings are the heart of **spectroscopes** and **spectrometers**. By measuring the precise angles at which known spectral lines appear, scientists can:

- Identify elements and molecules by their unique spectral fingerprints
- Measure the temperature of stars (from the shape of their spectra)
- Detect the velocity of galaxies via the **Doppler redshift** of spectral lines (Hubble's discovery that the universe is expanding came from spectroscopy)
- Analyze the composition of atmospheric pollutants, pharmaceuticals, and materials

### Holography

A **hologram** uses interference to record a 3D image. A laser beam is split into two: one illuminates the object (object beam), the other goes directly to a photographic plate (reference beam). The two beams interfere at the plate, creating a complex pattern of fringes that encodes the 3D phase information of the object beam. Illuminating the developed plate with a similar laser reconstructs the original wavefront — creating a 3D image.

Credit cards, passports, and banknotes use holograms as security features.

### Optical Coherence Tomography (OCT)

OCT is used in medicine (ophthalmology, cardiology) to create high-resolution cross-sectional images of tissue using infrared light interference. It is like ultrasound but with light, achieving resolution of ~1-10 µm. Eye doctors use OCT to image the retina in extraordinary detail.

### Interferometry in Science

The famous **LIGO** (Laser Interferometer Gravitational-Wave Observatory) uses a Michelson interferometer with 4-km-long arms. Gravitational waves from merging black holes or neutron stars stretch and compress space, changing the length of LIGO's arms by less than 10^-18 m — about 1/1000th the diameter of a proton. LIGO detects this by observing interference pattern shifts. This is one of the most precise measurements ever made, and it relies entirely on wave optics.

### Fiber Optic Sensors

Fiber Bragg gratings — periodic variations in the refractive index of optical fiber — act as diffraction gratings for light traveling along the fiber. They reflect specific wavelengths and transmit others, enabling precise temperature and strain sensors used in bridges, aircraft, and medical devices.

### Thin Film Coatings in Everyday Life

Beyond camera lenses, anti-reflection coatings appear on:
- **Eyeglasses** — reducing glare and ghosting
- **Solar cells** — reducing reflection to increase energy absorption
- **Display screens** — reducing reflections for outdoor readability
- **Telescope mirrors** — maximizing reflected light

High-reflection coatings (the opposite: choosing constructive interference) are used in laser cavities and mirrors. Stacking multiple thin film layers allows engineers to create mirrors that reflect 99.999% of light at a specific wavelength.

---

## Summary

- **Wave optics** describes light as a wave and explains phenomena that geometric (ray) optics cannot: interference and diffraction.

- Wave effects are significant when the **wavelength λ is comparable to the size of the obstacle or opening**. For visible light (λ ≈ 400–700 nm), this means very small features (micron scale or below).

- **Coherent sources** have the same frequency and a constant phase relationship. They are needed for stable, observable interference patterns. Lasers and derived beams (like double-slit) are coherent; separate light bulbs are not.

- In **Young's double-slit experiment**, two coherent sources separated by distance d produce alternating bright and dark **fringes** on a screen at distance D.

- **Bright fringe positions**: y_m = mλD/d (m = 0, ±1, ±2, ...)

- **Dark fringe positions**: y_m = (m + 1/2)λD/d (m = 0, ±1, ±2, ...)

- **Fringe spacing**: Δy = λD/d — equal spacing, proportional to λ and D, inversely proportional to d.

- **Thin film interference** occurs when light reflects from the top and bottom surfaces of a thin film. Path difference = 2nt. Phase shifts at boundaries (180° when going from low-n to high-n) modify the interference conditions.

- **Soap bubble colors** arise because different film thicknesses constructively reinforce different wavelengths of white light.

- **Anti-reflection coatings** use thin films (quarter-wave thickness) to cause destructive interference of reflected light, reducing glare and increasing transmission.

- **Diffraction** is the bending and spreading of waves around obstacles and through openings. It is significant when λ ~ opening size.

- **Single-slit diffraction** produces a wide central maximum with dark fringes at sinθ = mλ/a (m = ±1, ±2, ...). Narrower slit → wider diffraction pattern.

- **Diffraction gratings** (many equally spaced slits) produce sharp, bright principal maxima at angles given by d sinθ = mλ. They are used in spectrometers to separate and measure wavelengths precisely.

- **CDs and DVDs** act as reflection diffraction gratings, producing rainbow colors because their track spacing (1.6 µm for CDs) is comparable to visible wavelengths.

- Applications include: spectroscopy, holography, anti-reflection coatings, OCT medical imaging, LIGO gravitational wave detection, and thin film sensors.

---

## Key Equations

**Young's Double-Slit — Path Difference (small angle approximation):**

    Δ = d * sin(θ)  ≈  d * y / D

**Constructive interference (bright fringe):**

    Δ = m * λ,   m = 0, ±1, ±2, ...

**Destructive interference (dark fringe):**

    Δ = (m + 1/2) * λ,   m = 0, ±1, ±2, ...

**Bright fringe position on screen:**

    y_m = m * λ * D / d

**Dark fringe position on screen:**

    y_m = (m + 1/2) * λ * D / d

**Fringe spacing (distance between adjacent bright or dark fringes):**

    Δy = λ * D / d

**Thin film optical path difference:**

    Δ = 2 * n * t   (normal incidence)

**Quarter-wave anti-reflection coating thickness:**

    t = λ / (4 * n_coating)

**Single-slit diffraction — dark fringe angles:**

    sin(θ) = m * λ / a,   m = ±1, ±2, ...   (m = 0 is NOT a dark fringe)

**Angular half-width of central maximum (single slit):**

    θ_half ≈ λ / a   (in radians, for small angles)

**Diffraction grating equation (principal maxima):**

    d * sin(θ) = m * λ,   m = 0, ±1, ±2, ...

**Grating spacing from lines per metre:**

    d = 1 / N   (where N = number of lines per metre)

---

*End of Chapter 67*
