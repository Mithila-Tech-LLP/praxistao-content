# Chapter 68: Polarization

> **"Polarization is a property only transverse waves can have. This simple fact explains why polaroid sunglasses reduce glare, how 3D movies work, and how LCD screens display images."**

---

## Table of Contents

- [1. Transverse vs Longitudinal Waves Recap](#1-transverse-vs-longitudinal-waves-recap)
- [2. Unpolarized Light](#2-unpolarized-light)
- [3. Polarized Light](#3-polarized-light)
- [4. Methods of Polarization](#4-methods-of-polarization)
  - [4a. Polaroid Filters (Absorption Method)](#4a-polaroid-filters-absorption-method)
  - [4b. Reflection at Brewster's Angle](#4b-reflection-at-brewsters-angle)
  - [4c. Scattering](#4c-scattering-why-the-sky-is-polarized)
  - [4d. Birefringence](#4d-birefringence-double-refraction)
- [5. Malus's Law](#5-maluss-law)
- [6. The Three-Polarizer Paradox](#6-the-three-polarizer-paradox)
- [7. Worked Example 1: Two Polarizers](#7-worked-example-1-two-polarizers)
- [8. Worked Example 2: Three Polarizers](#8-worked-example-2-three-polarizers)
- [9. 3D Movies](#9-3d-movies)
- [10. LCD Screens](#10-lcd-screens)
- [11. Polaroid Sunglasses](#11-polaroid-sunglasses)
- [12. Brewster's Angle Worked Example](#12-brewsters-angle-worked-example)
- [13. Optical Activity](#13-optical-activity)
- [14. Circular Polarization](#14-circular-polarization)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 1. Transverse vs Longitudinal Waves Recap

Before we can understand polarization, we need to be crystal-clear about two kinds of waves.

### Transverse Waves

In a **transverse wave**, the particles (or the field) oscillate **perpendicular** (at right angles) to the direction the wave is traveling.

Imagine you are holding one end of a long rope and your friend holds the other end. If you shake your end **up and down**, the rope creates a wave that travels horizontally toward your friend — but the rope itself moves **up and down**, which is perpendicular to the horizontal direction of travel. That is a transverse wave.

```
Direction of wave travel →

     ↑     ↑
     |     |
  ───┼─────┼─── wave moves this way →
     |     |
     ↓     ↓

Rope particles oscillate up ↕ and down,
wave itself travels to the right →
```

**Examples of transverse waves:**
- Light (electromagnetic waves)
- Waves on a string or rope
- Water surface waves (approximately)
- Seismic S-waves

### Longitudinal Waves

In a **longitudinal wave**, the particles oscillate **parallel** to the direction of travel — they compress and expand back and forth along the same line the wave moves.

Think of a spring (Slinky) lying on a table. Push one end forward, and the compression travels along the spring. The coils move forward and backward — the same direction the wave travels.

```
Direction of wave travel →

  ███  ░░░  ███  ░░░  ███   →
  (compression and rarefaction travel rightward)
  
  Particles move ← → along the same line
```

**Examples of longitudinal waves:**
- Sound waves
- Seismic P-waves
- Compression waves in a spring

### Why Polarization is Only for Transverse Waves

Here is the key insight. In a transverse wave, the oscillation can happen in **any direction** that is perpendicular to the direction of travel. For a wave traveling along the x-axis, the oscillation can be along y, along z, or any combination of the two. **Polarization is simply specifying which of these perpendicular directions the oscillation is confined to.**

In a longitudinal wave, the oscillation is always along the same single direction as travel. There are no "other perpendicular directions" to choose from. The oscillation is locked into one dimension. Therefore, **there is no concept of polarization for longitudinal waves.**

This is why **sound cannot be polarized.** Sound is a longitudinal wave — the air molecules compress and expand along the direction the sound travels. There is no second perpendicular dimension for the vibration to be "filtered" into. No matter what you put in the path of sound, you cannot restrict it to one plane of oscillation in the way we can with light.

Think of it this way: polarizing a longitudinal wave would be like asking a train to move sideways on its tracks. The tracks only allow one direction — there is nothing to polarize.

---

## 2. Unpolarized Light

Natural light — from the sun, a candle flame, or an incandescent bulb — is **unpolarized light**.

What does that mean? Light is an electromagnetic wave. It has an **electric field** (E-field) that oscillates perpendicular to the direction the light travels. But in unpolarized light, this electric field does not stick to one direction. It rapidly and randomly flips between all possible orientations in the plane perpendicular to travel.

Think of it like a hula-hoop. If a light wave travels forward (into the page), the electric field could be pointing up, then upper-right, then right, then lower-right, then down — changing orientation billions of times per second, completely randomly, sampling all orientations with equal probability.

**Analogy:** Imagine shaking a rope randomly — sometimes up-down, sometimes left-right, sometimes diagonally. The rope wave has no single plane. That is unpolarized.

```
Unpolarized light traveling toward you (into the page):
(Each arrow is the E-field direction at different instants)

            ↑
       ↖    |    ↗
         \  |  /
    ←  ───( · )───  →      · = direction of travel (into page)
         /  |  \
       ↙    |    ↘
            ↓

E-field vibrates in ALL directions perpendicular to travel,
randomly changing orientation at extremely high frequencies.
```

**Sources of unpolarized light:**
- The Sun
- Incandescent light bulbs
- Candle flames
- Fluorescent tube lights
- LEDs (mostly unpolarized)
- Fire

In all of these sources, atoms emit light independently and randomly, so there is no preferred direction for the electric field oscillation. The result is unpolarized light — a superposition of waves with all possible polarization directions.

---

## 3. Polarized Light

**Polarized light** is light in which the electric field vibrates in only ONE specific direction (or plane). The simplest and most common type is **linear polarization** (also called **plane polarization**).

In linearly polarized light, if the wave travels along the z-axis, the electric field might always oscillate along the y-axis (vertical) — never along x, never diagonally. Just up and down, repeatedly.

```
Comparison: Unpolarized vs Linearly Polarized Light
===================================================

UNPOLARIZED (side view):
Wave travels → → → →

  ↕ ↗ ↔ ↘ ↕ ↗ ↔ ↘
  |/ |\ |/ |\ |/ |\
──────────────────────→  direction of travel

  E-field wiggles in all perpendicular directions randomly.


LINEARLY POLARIZED (side view):
Wave travels → → → →

  ↕   ↕   ↕   ↕   ↕
  |   |   |   |   |
──────────────────────→  direction of travel

  E-field ONLY wiggles up-down (vertical plane).
  All oscillations are in the same plane.


HEAD-ON VIEW (looking into the beam):

  Unpolarized:             Linearly Polarized (vertical):

        ↑                           ↑
    ↖   |   ↗                       |
      \ | /                         |
  ←────●────→                  ─────●─────
      / | \                         |
    ↙   |   ↘                       |
        ↓                           ↓

  Arrows in all directions.    Arrow in one direction only.
```

The **plane of polarization** is the plane containing both the direction of travel and the direction of the electric field oscillation.

**Vertical polarization:** E-field oscillates up-down.
**Horizontal polarization:** E-field oscillates left-right.
**Diagonal polarization:** E-field oscillates at some angle between horizontal and vertical.

All of these are examples of linear polarization — the key feature is that the E-field is confined to a single plane.

---

## 4. Methods of Polarization

There are four main ways to produce polarized light from unpolarized light.

---

### 4a. Polaroid Filters (Absorption Method)

A **Polaroid filter** (invented by Edwin Land in 1929) is the most common way to polarize light.

**How it works:**

A Polaroid filter contains long-chain polymer molecules (like polyvinyl alcohol) that have been stretched so they all align parallel to one another. These long molecules conduct electrons along their length. When the electric field of the light wave is oriented parallel to these long molecules, it drives the electrons up and down the chains — the electrons absorb the energy. That component of the light is **absorbed**.

When the electric field is oriented perpendicular to the molecular chains, it cannot drive electrons along the chains (they are too short in that direction), so it passes through with very little absorption.

**Result:** Only the component of the E-field **perpendicular to the molecular chain direction** passes through. The output is linearly polarized light.

```
POLAROID FILTER — how it works:
================================

Unpolarized light (all E-field directions) entering the filter:

    ↑↗→↘↓↙←↖  (all directions)
         |
         ↓
   ┌─────────────┐
   │ | | | | | | │   ← long molecule chains run vertically
   │ | | | | | | │      (parallel components get absorbed)
   │ | | | | | | │
   └─────────────┘
         |
         ↓
    →   →   →        Only horizontal E-field components pass through.
                     (perpendicular to the chains)

The "transmission axis" of the polaroid is HORIZONTAL here.
The "absorption axis" is VERTICAL.
```

**Important facts about Polaroid filters:**

1. When unpolarized light passes through a single Polaroid filter, the output intensity is approximately **half** the input intensity. This is because the filter passes only one of two equal perpendicular components.

   I_after = I_before / 2    (for unpolarized light through one polaroid)

2. The filter has a **transmission axis** — the direction of the E-field that is allowed through.

3. A second Polaroid filter placed after the first is called an **analyzer**. Its effect depends on the angle between its transmission axis and that of the first polaroid — this is described by **Malus's Law** (Section 5).

---

### 4b. Reflection at Brewster's Angle

When light hits a smooth, flat surface (like glass, water, or a wet road), part of it reflects and part refracts (enters the material). The reflected and refracted beams can be partially or completely polarized.

**Brewster's angle** (θ_B) is a special angle of incidence at which the **reflected beam is 100% linearly polarized**, parallel to the surface.

The formula for Brewster's angle:

    tan(θ_B) = n2 / n1

Where:
- n1 = refractive index of the medium the light is coming FROM
- n2 = refractive index of the medium the light is going INTO

**Example:** Light going from air (n1 = 1.00) into glass (n2 = 1.50):

    tan(θ_B) = 1.50 / 1.00 = 1.50
    θ_B = arctan(1.50) ≈ 56.3°

At this angle, the reflected light is completely horizontally polarized. The refracted (transmitted) beam is only partially polarized.

**The geometric reason this works:**

At Brewster's angle, something special happens geometrically: the reflected ray and the refracted ray are **exactly perpendicular to each other** (they form a 90° angle). This is not a coincidence — it is the condition required for complete polarization of the reflected beam.

```
BREWSTER'S ANGLE:
=================

Incoming unpolarized light
         \
          \  θ_B ≈ 56°
           \
────────────●──────────────  glass surface (n=1.5)
           /|
          / |
         /  |  90° angle between reflected
        /   |  and refracted rays!
       ↗    ↓
  Reflected  Refracted
  (100%      (partially
  polarized) polarized)

The reflected beam has ONLY the E-field component
parallel to the surface (horizontal polarization).
```

**Why glare is polarized:**

When sunlight strikes a horizontal surface (water, wet road, car hood) at a shallow angle — which is close to Brewster's angle — the reflected light (glare) is predominantly horizontally polarized. This is why Polaroid sunglasses (with vertical transmission axis) are so effective at cutting glare.

---

### 4c. Scattering (Why the Sky is Polarized)

When sunlight travels through Earth's atmosphere, it collides with tiny gas molecules (mainly nitrogen and oxygen). The molecules absorb this energy and re-emit it as light in all directions — this is called **Rayleigh scattering**.

But here is the key: the re-emitted (scattered) light is **not the same in all directions**. The molecules oscillate as electric dipoles, and a dipole cannot radiate along its own axis. This means the scattered light is **partially or fully polarized**.

**Specifically:**
- Light scattered **sideways** (at 90° to the Sun) is **completely linearly polarized**.
- Light scattered in other directions is partially polarized.
- Light coming directly from the Sun is unpolarized.

```
SCATTERING POLARIZATION:
=========================

                       Sun (unpolarized light)
                            ★
                            |
                            |  (sunlight going downward)
                            |
              ←─────────────●──────────────→
                           gas              
                          molecule          
                    (scatters light sideways)

  Light scattered at 90° to the sun's rays         →
  is completely POLARIZED (horizontal in this case) →
```

**Bees and the polarized sky:**

Honeybees can detect the polarization of sky light with special cells in their eyes. Even on a cloudy day with a small patch of clear sky, a bee can determine the direction of the Sun and navigate back to the hive. This is one of nature's most elegant uses of polarization.

Ants, mantis shrimps, and some species of fish also use sky polarization for navigation.

---

### 4d. Birefringence (Double Refraction)

Some crystals — most famously **calcite** (Iceland spar) and **quartz** — have a property called **birefringence** or **double refraction**.

In an ordinary material, light travels at the same speed regardless of the direction of its E-field oscillation. The refractive index is the same for all polarization directions.

In a **birefringent crystal**, the arrangement of atoms is asymmetric. Light polarized along one crystal axis travels at one speed; light polarized along a perpendicular crystal axis travels at a different speed. The crystal has **two different refractive indices** for the two perpendicular polarization directions.

**What this means in practice:**

When a single beam of unpolarized light enters a birefringent crystal, it splits into **two separate beams:**
- The **ordinary ray (o-ray):** follows normal Snell's law refraction, has one specific polarization direction.
- The **extraordinary ray (e-ray):** refracts at a different angle (uses a different refractive index), has a polarization perpendicular to the ordinary ray.

The two emerging beams are linearly polarized in perpendicular directions.

```
BIREFRINGENCE in a calcite crystal:
=====================================

  Unpolarized light
         ↓
  ┌──────────────┐
  │   calcite    │
  │   crystal    │
  └──────────────┘
       ↙    ↘
      ↙      ↘
   o-ray    e-ray
  (polarized  (polarized
  vertically) horizontally)

Two separate spots visible on paper below the crystal.
Each spot is linearly polarized.
Each is polarized perpendicular to the other.
```

Birefringent crystals are widely used in optical equipment — polarizing prisms (Nicol prism, Wollaston prism), wave plates, and laser optical components all rely on birefringence.

---

## 5. Malus's Law

**Malus's Law** describes what happens when already-polarized light passes through a second polarizer (called an **analyzer**).

Suppose linearly polarized light of intensity I_0 passes through an analyzer whose transmission axis makes an angle **θ** with the polarization direction of the incoming light. The transmitted intensity I is:

```
  I = I_0 × cos²(θ)
```

This is **Malus's Law**, discovered by the French physicist Étienne-Louis Malus in 1808.

**Physical interpretation:**

The incoming E-field has amplitude E_0. Only the component of E_0 along the analyzer's transmission axis passes through. That component has amplitude E_0 × cos(θ). Since intensity is proportional to the square of amplitude:

    I ∝ E² = (E_0 cos θ)² = E_0² cos²θ = I_0 cos²θ

```
MALUS'S LAW — Setup:
=====================

  Polarized light                          Analyzer
  (intensity I_0)                          (transmission axis
        |                                   at angle θ)
        ↓                                      |
  ┌──────────┐      I_0 (polarized)      ┌──────────┐      I
  │ Polarizer│  ─────────────────────→   │ Analyzer │  ─────────→
  │ (creates │                           │ (at angle│
  │polarized │                           │  θ to    │
  │  light)  │                           │ first)   │
  └──────────┘                           └──────────┘
        ↑
  Transmission axis
  of first polarizer                  I = I_0 × cos²(θ)


  Transmission axes visualized (head-on):

  First polarizer:     Analyzer at θ=0°:    Analyzer at θ=45°:   Analyzer at θ=90°:
        ↑                    ↑                      ↗                     →
        |                    |                    /                       
        |                    |                  /                         
        |                    |                /                           
   I = I_0             I = I_0           I = I_0/2               I = 0
                      (full pass)       (half passes)          (blocked)
```

**Key special cases:**

| Angle θ | cos²(θ) | Transmitted Intensity |
|---------|---------|----------------------|
| 0°      | 1       | I = I_0 (all light passes) |
| 30°     | 3/4     | I = 0.75 × I_0 |
| 45°     | 1/2     | I = 0.5 × I_0 |
| 60°     | 1/4     | I = 0.25 × I_0 |
| 90°     | 0       | I = 0 (no light passes) |

**Important:** Malus's Law applies when the **incoming light is already polarized**. If the incoming light is unpolarized, you cannot directly use I = I_0 cos²θ. Instead, an unpolarized beam passing through the first polaroid loses half its intensity (I = I_0/2), and then Malus's Law applies for any additional polarizers after that.

---

## 6. The Three-Polarizer Paradox

Here is one of the most counterintuitive results in optics — the kind of result that makes physics genuinely surprising and beautiful.

**The setup:**
- Take two Polaroid filters with their transmission axes perpendicular (crossed at 90°).
- Shine light through both.
- Result: No light comes out. The second filter blocks everything. (cos²90° = 0)

This makes sense. Two crossed polarizers = complete darkness.

**Now try this:**
- Take those same two crossed polarizers.
- Insert a THIRD polarizer BETWEEN them, at 45° to each of the first two.

You might expect this to make things even darker — after all, you are adding something that absorbs light. But something remarkable happens: **light comes through!**

**Why? Let's trace the intensity step by step.**

Let the initial unpolarized light have intensity I_initial.

```
THREE POLARIZER SETUP:
=======================

Unpolarized      Polarizer 1     Polarizer 2      Polarizer 3
  light          (vertical,      (at 45° to       (horizontal,
                  0°)            first, 45°)       90° to first)

I_initial  ─────→ [P1] ─────→ [P2] ─────→ [P3] ─────→  I_final

     Transmission axes:
     P1: ↑ (vertical)
     P2: ↗ (45° diagonal)
     P3: → (horizontal)
```

**Step 1:** Unpolarized light through Polarizer 1 (vertical):

    I_1 = I_initial / 2

(A polaroid always cuts unpolarized light to half its intensity.)

**Step 2:** Polarized light (vertical) through Polarizer 2 at 45°:

    Angle between P1 and P2 = 45°
    I_2 = I_1 × cos²(45°)
    I_2 = (I_initial / 2) × (1/√2)²
    I_2 = (I_initial / 2) × (1/2)
    I_2 = I_initial / 4

The light coming out of P2 is now polarized at 45°.

**Step 3:** Light polarized at 45° through Polarizer 3 (horizontal = 90° from vertical):

    Angle between P2 and P3 = 45°
    (P2 is at 45°, P3 is at 90° = horizontal, so the angle between them is 45°)
    I_3 = I_2 × cos²(45°)
    I_3 = (I_initial / 4) × (1/2)
    I_3 = I_initial / 8

**Final result:** I_final = I_initial / 8

Without the middle polarizer: **no light.** With the middle polarizer: **1/8 of the original intensity.**

Adding a polarizer increased the transmitted light from zero to a positive value!

**Why does this work?** The middle polarizer "rotates" the polarization state by 45° in two steps. It essentially acts as a bridge, converting the vertical polarization into a diagonal, and then the horizontal filter can extract a component from that diagonal. The key is that each polarizer only blocks the perpendicular component — it does not completely block any light that is not exactly perpendicular.

This result has deep implications in **quantum mechanics** and is related to why quantum measurements affect the system being measured.

---

## 7. Worked Example 1: Two Polarizers

**Problem:**

Unpolarized light of intensity **800 W/m²** passes through two polarizers in sequence.
- The first polarizer has a vertical transmission axis.
- The second polarizer (analyzer) has its axis at **30°** to the first.

Find:
(a) The intensity after the first polarizer.
(b) The intensity after the second polarizer.
(c) The percentage of original light that makes it through.

---

**Solution:**

**Given:**
- Initial intensity: I_0 = 800 W/m²
- Angle between polarizers: θ = 30°
- Light is initially unpolarized.

---

**Part (a): Intensity after the first polarizer**

When unpolarized light passes through a Polaroid filter, only half the intensity gets through. This is because unpolarized light has equal energy in all polarization directions. The filter selects one linear component and discards the rest, retaining about half.

    I_1 = I_0 / 2
    I_1 = 800 / 2
    I_1 = 400 W/m²

The light exiting the first polarizer is **400 W/m²**, and it is now **linearly polarized** in the vertical direction.

---

**Part (b): Intensity after the second polarizer**

Now we apply Malus's Law. The incoming light is polarized (vertical), and the second polarizer's axis is at 30° to it.

    I_2 = I_1 × cos²(θ)
    I_2 = 400 × cos²(30°)

What is cos(30°)?

    cos(30°) = √3/2 ≈ 0.866

So:

    cos²(30°) = (√3/2)² = 3/4 = 0.75

Therefore:

    I_2 = 400 × 0.75
    I_2 = 300 W/m²

The intensity after the second polarizer is **300 W/m²**, and the light is now polarized at 30° (along the second polarizer's transmission axis).

---

**Part (c): Percentage transmitted**

    Percentage = (I_2 / I_0) × 100
    Percentage = (300 / 800) × 100
    Percentage = 37.5%

Only **37.5%** of the original light makes it through both polarizers.

---

**Summary of results:**

```
  Unpolarized      After P1           After P2
  light            (vertical)         (at 30°)
  800 W/m²   →    400 W/m²     →    300 W/m²

  Intensity:       ÷2 (always)       × cos²(30°) = × 0.75
```

---

## 8. Worked Example 2: Three Polarizers

**Problem:**

Light passes through three polarizers in sequence.
- After the first polarizer, the intensity is I_0 (light is linearly polarized, vertical).
- The second polarizer is at **60°** to the first.
- The third polarizer is at **90°** to the first (horizontal).

Find:
(a) Intensity after the second polarizer.
(b) Intensity after the third polarizer.
(c) What would happen if the second polarizer were removed?

---

**Solution:**

**Given:**
- Intensity entering the problem: I_0 (already polarized by first polarizer)
- P1 axis: 0° (vertical)
- P2 axis: 60° (to P1)
- P3 axis: 90° (to P1) = horizontal

---

**Part (a): Intensity after the second polarizer**

Applying Malus's Law between P1 and P2:

    Angle between P1 and P2 = 60°
    I_2 = I_0 × cos²(60°)

What is cos(60°)?

    cos(60°) = 1/2 = 0.5

So:

    cos²(60°) = (1/2)² = 1/4 = 0.25

Therefore:

    I_2 = I_0 × 0.25
    I_2 = I_0 / 4

The light is now polarized at **60°** to vertical (along P2's axis) with intensity **I_0/4**.

---

**Part (b): Intensity after the third polarizer**

Now we apply Malus's Law between P2 and P3.

P2 is at 60° from vertical. P3 is at 90° from vertical (horizontal).

The angle **between P2 and P3** is:

    θ_23 = 90° - 60° = 30°

Applying Malus's Law:

    I_3 = I_2 × cos²(θ_23)
    I_3 = (I_0 / 4) × cos²(30°)
    I_3 = (I_0 / 4) × (3/4)
    I_3 = 3 × I_0 / 16
    I_3 ≈ 0.1875 × I_0

The final intensity is **3I_0/16 ≈ 18.75% of I_0**.

---

**Part (c): If the second polarizer is removed**

Without P2, the light goes from P1 (vertical, intensity I_0) directly to P3 (horizontal, 90° to P1).

    Angle between P1 and P3 = 90°
    I = I_0 × cos²(90°) = I_0 × 0 = 0

**No light comes through.**

Removing the middle polarizer causes the output to drop from I_0/16 × 3 to zero. This is the three-polarizer paradox we discussed in Section 6 — adding a polarizer actually **increases** the transmitted intensity.

---

**Summary of results:**

```
  After P1:        After P2 (at 60°):   After P3 (at 90° from P1):
  I_0          →   I_0/4          →     3I_0/16

  (Without P2, light goes from I_0 at P1 directly to 0 at P3.)
```

---

## 9. 3D Movies

Three-dimensional movies create the illusion of depth by presenting slightly different images to your left and right eyes — just like how your two eyes naturally see slightly different views of real objects (this is called **binocular disparity**).

**The fundamental challenge:** You need the left eye to see image A and the right eye to see image B, even though both eyes are looking at the same screen.

**How polarization solves this:**

```
3D CINEMA — Polarization Method:
==================================

  Projector A  ──→ [Horizontal polarizer] ──→ ╲
                                                ╲
                                                 [ SCREEN ]
                                                ╱
  Projector B  ──→ [Vertical polarizer]   ──→ ╱

  (Both images are projected onto the screen simultaneously)

  LEFT LENS of glasses: Horizontal polarizer
  RIGHT LENS of glasses: Vertical polarizer

  Left eye sees only image A (horizontal polarization matches left lens).
  Right eye sees only image B (vertical polarization matches right lens).
  Brain combines them → 3D depth perception!
```

**Old-school linear polarization cinema:**

In older systems (IMAX, some 3D cinemas), two images are projected through **linearly polarized** filters at perpendicular orientations — one horizontal, one vertical. The audience wears glasses with matching polaroid lenses. The screen must be a special metallic screen that preserves polarization (regular screens scatter light and randomize polarization).

**The tilt problem with linear polarization:**

If you tilt your head sideways, the polarization axes of your glasses tilt too, but the screen images do not. This can cause "ghosting" — each eye faintly sees the image meant for the other eye.

**Modern RealD circular polarization:**

Modern RealD cinema uses **circular polarization** (Section 14) instead of linear. With circular polarization, head-tilting does not cause ghosting. One image is projected with left-circular polarization, the other with right-circular polarization. The glasses have matching circular polarization lenses. Regular screens can be used, though they still need to be somewhat reflective.

**At home (passive 3D TVs):**

Some passive 3D televisions use the same polarization principle. Odd rows of pixels emit horizontally polarized light (for one eye) and even rows emit vertically polarized light (for the other eye). Resolution is halved vertically for each eye, but no battery is needed in the glasses.

Active 3D systems (active shutter glasses) use a different approach — LCD shutters that alternately block each eye in sync with the TV alternating between left and right images. These glasses require batteries and electronics.

---

## 10. LCD Screens

**Liquid Crystal Display (LCD)** technology is built entirely on polarization. Understanding LCDs is a wonderful application of everything we have covered.

### The Basic Components

An LCD pixel has, from back to front:
1. A backlight (bright white light source)
2. A **rear polarizer** (polarizes backlight in one direction, say horizontal)
3. A layer of **liquid crystals**
4. A **front polarizer** (with its axis perpendicular to the rear polarizer — vertical)
5. A color filter (red, green, or blue) for colored pixels

### The OFF State (no electric voltage applied)

When no voltage is applied to the liquid crystal layer, the liquid crystal molecules **spontaneously arrange themselves in a twisted structure** — they gradually rotate from one alignment at the back to a 90°-rotated alignment at the front.

As horizontally polarized light passes through this twisted structure, the twist **guides the polarization, rotating it 90°** — from horizontal at the back to vertical at the front.

Now the light arrives at the front polarizer, which has a vertical transmission axis. Since the light is now vertically polarized, it **passes through** the front polarizer.

**OFF state = bright pixel (white or colored)**

### The ON State (voltage applied)

When an electric voltage is applied to the liquid crystal layer, the electric field forces the liquid crystal molecules to **align with the field** — they all stand up perpendicular to the glass, losing their twisted arrangement.

Now the light passes through the liquid crystals without any rotation. It remains horizontally polarized when it exits the liquid crystal layer.

The front polarizer is vertically oriented. Horizontal polarization is perpendicular to vertical — so Malus's Law gives cos²(90°) = 0. **The light is completely blocked.**

**ON state = dark pixel (black)**

```
LCD OPERATION:
===============

    BACKLIGHT → REAR POLARIZER → LIQUID CRYSTALS → FRONT POLARIZER → VIEWER

OFF (no voltage):                        ON (voltage applied):
──────────────────                       ─────────────────────
                                         
 White      →→→   Horizontal polarized    White      →→→   Horizontal polarized
 backlight         light enters LC layer  backlight         light enters LC layer
                         |                                       |
                    LC molecules                           LC molecules
                    twisted 90°                            untwisted (aligned
                         |                                  with E-field)
                         ↓                                       |
               Vertical polarized                                ↓
               light exits LC layer                   Still horizontal polarized
                         |                                       |
                  ═══════╪═══════                        ═══════╪═══════
                  Front polarizer                        Front polarizer
                  (vertical axis)                        (vertical axis)
                  PASSES LIGHT ✓                         BLOCKS LIGHT ✗
                         |
                         ↓
                    BRIGHT PIXEL                           DARK PIXEL
```

### Color and Subpixels

Each full pixel is made of three **subpixels** — one covered by a red filter, one by a green filter, one by a blue filter. Each subpixel can be independently controlled to any brightness level between fully off (0) and fully on (255). By mixing red, green, and blue at different intensities, any color can be produced.

### Viewing Angle Issues

LCD screens can look dim or washed-out when viewed from the side. This is because liquid crystals only rotate polarization effectively when viewed head-on. At steep angles, the optical path length through the liquid crystal layer changes, the polarization rotation is incomplete, and contrast suffers. IPS (In-Plane Switching) panels address this by changing how the liquid crystals are arranged, improving viewing angles significantly.

This is a direct consequence of how polarization interacts with the anisotropic (direction-dependent) liquid crystal molecules.

---

## 11. Polaroid Sunglasses

Polaroid sunglasses are one of the most practical everyday applications of polarization. Let us understand exactly why they work.

### The Source of Glare

When sunlight strikes a horizontal surface — a lake, a wet road, the hood of a car, a snow field — it reflects. The reflected light is the glare that causes eye strain and visual discomfort.

As we covered in Section 4b, when light reflects off a non-metallic surface at or near Brewster's angle, the reflected light is predominantly or completely polarized **parallel to the surface** — which for horizontal surfaces means **horizontal polarization**.

Think of it: sunlight comes down at an angle, hits a horizontal surface, and bounces up into your eyes. The glare is mostly horizontally polarized (the E-field oscillates left-right, parallel to the water or road surface).

```
GLARE GEOMETRY:
================

         Sun
          ★
          |\
          | \  (sunlight comes down at an angle)
          |  \
          |   \
          |    ● ← reflection point on water surface
          |   /
          |  /   (reflected glare goes up into eyes)
          | /
          |/
         👁  ← your eyes

  The reflected glare is HORIZONTALLY POLARIZED
  (E-field oscillates parallel to the water surface: ←→)
```

### How Polaroid Sunglasses Block Glare

Polaroid lenses are manufactured with their **transmission axis oriented vertically**. They only allow vertically polarized light (↕) to pass through. They block horizontally polarized light (↔).

Since glare from horizontal surfaces is mostly horizontally polarized, and the sunglasses transmit only vertical polarization:

    Malus's Law: I = I_0 × cos²(θ)
    θ = 90° (glare is horizontal, glasses pass vertical)
    I = I_0 × cos²(90°) = I_0 × 0 = 0

The glare is **completely blocked** (or nearly so, since real reflections are not perfectly at Brewster's angle).

Regular non-polarized sunglasses reduce overall light intensity equally for all directions — they make everything darker. Polaroid sunglasses specifically target the glare while still allowing through diffuse light (which is less polarized and comes from all directions).

### Fishermen and Polaroid Glasses

```
WITHOUT polaroid glasses:        WITH polaroid glasses:
========================         =======================

      ★ Sun                            ★ Sun
      |                                |
      |                                |
  ~~~~|~~~~  ←glare→   👁           ~~~|~~~              👁
  ~~~~|~~~~  (bright   (sees        ~~~|~~~  (glare     (sees
  ~~~~|~~~~   glare)    only        ~~~|~~~   blocked)   fish !)
      |                glare)          |
      🐟 (invisible                    🐟 (now visible
          under glare)                     under the surface!)
```

When you look into a river without polaroid glasses, the glare off the surface is so bright it overwhelms the faint reflected light from underwater objects. A fisherman wearing polaroid glasses cuts out the glare entirely and can clearly see fish swimming below the surface.

This is why fishing is probably the single biggest commercial market for polaroid sunglasses.

---

## 12. Brewster's Angle Worked Example

**Problem:**

Monochromatic light travels from **air** (n_1 = 1.00) into **water** (n_2 = 1.33).

(a) Calculate Brewster's angle for this interface.
(b) Find the angle of refraction when the light is incident at Brewster's angle.
(c) Verify that the reflected and refracted rays are perpendicular (90° apart).
(d) What does this perpendicularity tell us about polarization?

---

**Solution:**

**Part (a): Calculate Brewster's angle**

Using the Brewster's angle formula:

    tan(θ_B) = n_2 / n_1
    tan(θ_B) = 1.33 / 1.00
    tan(θ_B) = 1.33

Taking the inverse tangent:

    θ_B = arctan(1.33)
    θ_B ≈ 53.1°

So the Brewster's angle for an air-water interface is approximately **53.1°**.

---

**Part (b): Angle of refraction at Brewster's angle**

At θ_i = 53.1°, use Snell's Law to find the refraction angle θ_r:

    n_1 × sin(θ_i) = n_2 × sin(θ_r)
    1.00 × sin(53.1°) = 1.33 × sin(θ_r)

What is sin(53.1°)?

    sin(53.1°) ≈ 0.7997 ≈ 0.80

Therefore:

    1.00 × 0.80 = 1.33 × sin(θ_r)
    sin(θ_r) = 0.80 / 1.33
    sin(θ_r) = 0.6015
    θ_r = arcsin(0.6015)
    θ_r ≈ 37.0°

The refracted ray enters the water at approximately **37.0°** from the normal.

---

**Part (c): Verify the reflected and refracted rays are perpendicular**

The reflected ray leaves the surface at the same angle as the incident ray (law of reflection):

    Angle of reflected ray from normal = θ_B = 53.1°

The refracted ray is at:

    θ_r = 37.0° from the normal (on the other side of the surface)

The angle between the reflected ray and the refracted ray:

The reflected ray goes up at 53.1° from the normal (on the same side as the incoming ray). The refracted ray goes down at 37.0° from the normal (on the transmission side).

The total angle between them (measured from one to the other, going through the surface material):

    θ_reflected-from-surface + θ_refracted-from-surface
    = (90° - 53.1°) + (90° - 37.0°)

Wait — let us think about this more carefully.

The reflected ray makes angle 53.1° with the surface normal, so it makes (90° - 53.1°) = 36.9° with the surface itself.

The refracted ray makes 37.0° with the normal, so it makes (90° - 37.0°) = 53.0° with the surface.

Now the angle between the reflected ray and the refracted ray:

    Angle = 180° - 53.1° - 37.0°
    Angle = 180° - 90.1°
    Angle ≈ 90° ✓

The small discrepancy (0.1°) is due to rounding. In exact arithmetic:

    θ_B + θ_r = arctan(n) + arcsin(sin(arctan(n))/n)

This always equals 90° at Brewster's angle.

**The reflected and refracted rays are indeed perpendicular.**

---

**Part (d): Why does perpendicularity cause polarization?**

Think of the molecules in the water near the surface as tiny antennas. When the electric field of light hits them, they oscillate and re-radiate light (this is the origin of reflection).

An oscillating dipole (antenna) cannot radiate along its own axis — just like a radio antenna does not transmit directly above or below itself.

At Brewster's angle, the reflected beam would have to come from the refracted beam's electric field components pointing along the reflected ray direction. But at exactly this angle (where reflected and refracted are perpendicular), those components would require the dipoles to radiate along their own axis — which they cannot do.

Therefore, the component of the E-field that would cause vibration **parallel to the reflected ray direction** (the p-component, or transverse-magnetic component) cannot be reflected. Only the s-component (perpendicular to the plane of incidence, parallel to the surface) can be reflected. The reflected light is **completely s-polarized** — polarized parallel to the surface.

---

## 13. Optical Activity

Some transparent materials have a remarkable property: they rotate the plane of polarization of light as it passes through them. This property is called **optical activity**.

### What is Optical Activity?

When linearly polarized light enters an **optically active** substance, the plane of polarization continuously rotates. After traveling a certain distance through the material, the polarization has turned by a definite angle. The outgoing light is still linearly polarized — just in a different direction than it entered.

```
OPTICAL ACTIVITY:
==================

Linearly polarized                      Linearly polarized
light enters:                           light exits:
    |                                       \
    | (vertical)    → OPTICALLY ACTIVE →   \ (rotated by angle α)
    |                   MATERIAL            \
    
The polarization plane has rotated by angle α.
```

### Dextrorotatory and Levorotatory

Optically active materials rotate polarization in two directions:

- **Dextrorotatory (+):** Rotates polarization **clockwise** when looking toward the light source. Often labeled (+) or (d). Example: natural glucose (D-glucose).
- **Levorotatory (-):** Rotates polarization **counterclockwise** when looking toward the light source. Often labeled (-) or (l). Example: fructose, L-glucose.

The (+) and (-) molecules are **mirror images** of each other (like left and right hands) — they are called **enantiomers**. The ability to distinguish them using polarized light is one of the most important tools in chemistry and pharmaceutical science.

### The Formula for Rotation

The angle of rotation α (in degrees) depends on:

    α = [α] × c × l

Where:
- [α] = specific rotation (a property of the substance, degrees·mL/(g·dm))
- c = concentration of the solution (g/mL)
- l = path length through the solution (in decimeters, dm)

### Saccharimetry: Measuring Sugar with Polarized Light

One of the most important practical uses of optical activity is **saccharimetry** — measuring the concentration of sugar (sucrose, glucose, fructose) in a solution using polarized light.

**How a saccharimeter (polarimeter) works:**

```
SACCHARIMETER (Polarimeter):
==============================

  Light source → [Polarizer] → [Sugar solution tube] → [Analyzer] → Detector/eye
  
                    ↑                   ↑                    ↑
               Creates             Rotates the          Rotate this
               polarized           polarization         until maximum
               light               by angle α           brightness
                                   (proportional         to read off α
                                    to concentration)
```

1. Polarize the light with the first polaroid.
2. Pass it through a tube filled with the sugar solution.
3. The sugar rotates the polarization by angle α.
4. Rotate the analyzer until maximum brightness is seen.
5. Read the angle — that is α.
6. Use α = [α] × c × l to calculate concentration c.

**Applications:**
- **Food industry:** Quality control of sugar syrups, molasses, honey.
- **Medicine:** Measuring blood glucose concentration (though modern methods are more commonly used today).
- **Pharmaceutical industry:** Checking purity of drug molecules (since many drugs have optically active forms, only one of which is medically effective).
- **Wine making:** Measuring sugar content during fermentation.

### Examples of Optically Active Substances

| Substance | Rotation Direction | [α] (approximate) |
|-----------|-------------------|-------------------|
| D-glucose (dextrose) | + (dextro) | +52.7° |
| D-fructose | - (levo) | -92° |
| Sucrose (table sugar) | + (dextro) | +66.5° |
| L-amino acids | - (levo) | varies |
| Turpentine | - (levo) | -37° |
| Quartz | + or - (both forms exist) | depends on crystal form |

### Worked Mini-Example: Glucose Concentration

A 10 cm (1 dm) tube is filled with glucose solution. The polarimeter shows a rotation of +26.4°. The specific rotation of glucose is +52.7° mL/(g·dm). Find the concentration.

    α = [α] × c × l
    26.4 = 52.7 × c × 1
    c = 26.4 / 52.7
    c ≈ 0.501 g/mL ≈ 500 g/L

The glucose concentration is approximately **0.5 g/mL**.

---

## 14. Circular Polarization

We have been discussing **linear polarization** — the electric field oscillates in a single plane. But there is another important type: **circular polarization**.

### What is Circular Polarization?

Imagine two linearly polarized waves traveling in the same direction (say, along the z-axis):
- Wave 1: polarized along the x-axis (horizontal)
- Wave 2: polarized along the y-axis (vertical)

If these two waves have the same amplitude and are **90° out of phase** with each other, their combination produces a wave where the electric field vector **rotates** as the wave travels — like a corkscrew.

```
CIRCULAR POLARIZATION — E-field tip traces a helix:
====================================================

Looking head-on (along direction of travel):

  Left-circular:           Right-circular:
  
  E-field rotates          E-field rotates
  counterclockwise:        clockwise:
  
       ↑                        ↑
    ←  ·  →                  ←  ·  →
  rotates ↺                 rotates ↻
       ↓                        ↓


Side view (both components visible):
  
  Horizontal component: ─── ─── ─── ───   (sine wave)
  Vertical component:   │ │ │ │ │ │ │ │   (cosine wave, 90° phase shifted)
  
  Combined E-field tip traces a spiral (helix) as it travels.
```

If the two component waves have unequal amplitudes, the result is **elliptical polarization** — the E-field tip traces an ellipse rather than a circle.

### How to Produce Circular Polarization

A **quarter-wave plate** (λ/4 plate) is a birefringent crystal cut to a specific thickness. When linearly polarized light enters it at 45° to the crystal's axes, one component travels slightly slower than the other, creating exactly a 90° phase difference. The output is circularly polarized light.

### Uses of Circular Polarization

**Modern 3D cinema (RealD):**

As mentioned in Section 9, RealD cinema uses circular polarization. One projector image is left-circularly polarized, the other right-circularly polarized. The glasses have a left-circular polarizer for one eye and a right-circular polarizer for the other. The key advantage over linear polarization: head tilting does not cause the 3D effect to break down.

**Satellite and radio communications:**

Circular polarization is widely used in satellite communications. Because the satellite moves relative to the ground station, and because the ionosphere can rotate linear polarization (Faraday rotation), circular polarization provides a more stable signal regardless of orientation.

**Photography (circular polarizer filters):**

Camera filters for controlling reflections and glare are usually circular polarizers (not linear). This is because camera autofocus and metering systems can be confused by linearly polarized light due to the beam splitters inside the camera. A circular polarizer converts the linear polarization back into circular after it passes through the linear polarizer layer, ensuring the camera electronics work correctly.

**Chiral molecules and life:**

Nature itself uses circular polarization at the molecular level. The amino acids in all living organisms are exclusively L-amino acids (levorotatory), and the sugars are D-sugars (dextrorotatory). This universal handedness of biological molecules, called **homochirality**, remains one of the great unsolved puzzles in the origin of life. It is studied using circularly polarized light via a technique called **circular dichroism spectroscopy**.

---

## Summary

- **Polarization** is a property exclusive to **transverse waves**. Sound (longitudinal) cannot be polarized.

- **Unpolarized light** has its electric field vibrating in all possible orientations perpendicular to the direction of travel, rapidly and randomly.

- **Linearly polarized light** has its electric field confined to a single plane. It can be produced by:
  - **Polaroid filters** (absorption of one component by aligned molecular chains)
  - **Reflection at Brewster's angle** (tan θ_B = n2/n1)
  - **Scattering** (sky polarization, used by bees for navigation)
  - **Birefringence** (double refraction in crystals, producing two polarized beams)

- When unpolarized light passes through a single Polaroid, intensity drops to **half**: I = I_0/2.

- **Malus's Law:** When polarized light of intensity I_0 passes through a polarizer at angle θ, transmitted intensity is **I = I_0 cos²(θ)**.

- Two **crossed polarizers** (θ = 90°) block all light. Inserting a **third polarizer at 45°** between them allows 1/8 of the original intensity through — a beautiful paradox demonstrating that adding an absorber can increase transmitted light.

- **3D movies** use two images projected through perpendicular polarizers; glasses with matching polaroid lenses ensure each eye sees only its intended image.

- **LCD screens** use two crossed polarizers sandwiching a liquid crystal layer. Voltage untwists the crystals to block light (dark pixel); no voltage allows the crystals to twist polarization 90° so light passes (bright pixel).

- **Polaroid sunglasses** work because glare from horizontal surfaces (roads, water) is predominantly horizontally polarized; the vertical transmission axis of the lens blocks it.

- At **Brewster's angle**, reflected and refracted rays are perpendicular; the reflected beam is completely polarized parallel to the surface.

- **Optical activity** occurs in substances like glucose and quartz, rotating the plane of polarization by an angle proportional to concentration and path length (used in saccharimetry).

- **Circular polarization** arises when two perpendicular linearly polarized waves are superposed with 90° phase difference; the E-field vector rotates helically. Used in modern 3D cinema (RealD), satellite communications, and camera filters.

---

## Key Equations

**Intensity after one Polaroid (unpolarized input):**

    I = I_0 / 2

**Malus's Law (polarized light through analyzer at angle θ):**

    I = I_0 × cos²(θ)

**Three-polarizer result (unpolarized through P1 at 0°, P2 at 45°, P3 at 90°):**

    I_final = I_initial / 8

**Brewster's angle:**

    tan(θ_B) = n_2 / n_1

Where n_1 is the index of the incident medium, n_2 is the index of the transmission medium.

**At Brewster's angle (geometric property):**

    θ_B + θ_refracted = 90°

**Optical rotation (saccharimetry):**

    α = [α] × c × l

Where:
- α = observed rotation (degrees)
- [α] = specific rotation (degrees · mL / g · dm)
- c = concentration (g/mL)
- l = path length (dm)

**Snell's Law (used with Brewster's angle):**

    n_1 × sin(θ_i) = n_2 × sin(θ_r)

**Relationship between reflected and refracted angle at Brewster's angle:**

    θ_refracted = 90° - θ_B
    Therefore: θ_B + θ_refracted = 90°

**Amplitude component through polarizer (used to derive Malus's Law):**

    E_transmitted = E_0 × cos(θ)
    I_transmitted = I_0 × cos²(θ)    (since I ∝ E²)

---

*End of Chapter 68: Polarization*
