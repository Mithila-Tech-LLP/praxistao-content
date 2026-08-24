# Chapter 61: Light and Its Properties

> **"Light is the fastest thing in the universe, the carrier of information across cosmic distances, and the tool that lets us see everything. Understanding light is understanding sight itself."**

---

## Table of Contents

1. [What is Light?](#1-what-is-light)
2. [The Visible Spectrum](#2-the-visible-spectrum)
3. [Sources of Light](#3-sources-of-light)
4. [Light Travels in Straight Lines](#4-light-travels-in-straight-lines)
5. [Shadows: Umbra and Penumbra](#5-shadows-umbra-and-penumbra)
6. [Pinhole Camera (Camera Obscura)](#6-pinhole-camera-camera-obscura)
7. [Speed of Light in a Medium](#7-speed-of-light-in-a-medium)
8. [Coherent vs Incoherent Light](#8-coherent-vs-incoherent-light)
9. [Polarization (Preview)](#9-polarization-preview)
10. [Summary](#summary)
11. [Key Equations](#key-equations)

---

## 1. What is Light?

Light is one of the most fundamental things in the universe — it is how we see everything around us, how the Sun heats the Earth, how your microwave oven works, how your phone receives Wi-Fi, and how doctors take X-rays. All of these are forms of the same thing: **electromagnetic radiation**.

### Light as an Electromagnetic Wave

To understand light, you need to picture two things moving through space at the same time:

- An **electric field** — an invisible influence that pushes or pulls electric charges
- A **magnetic field** — an invisible influence related to moving charges and magnets

In a light wave, these two fields oscillate (wiggle back and forth) at right angles to each other, and both are perpendicular to the direction the light is traveling. This is called a **transverse wave**.

Here is a picture of that:

```
         Direction of travel
         ─────────────────────────────────────────►

         ↑ Electric Field (E)
         │     /\          /\          /\
         │    /  \        /  \        /  \
         │───/────\──────/────\──────/────\───────►
         │  /      \    /      \    /      \
         │ /        \  /        \  /        \
         │/          \/          \/          \

         Magnetic Field (B) — oscillates in the
         plane coming out of / going into the page,
         perpendicular to E and to direction of travel
```

This combination — electric and magnetic fields oscillating together — is an **electromagnetic wave**. Light is an electromagnetic wave. No medium (no air, no water, no wire) is needed for it to travel. It propagates through the vacuum of space perfectly well. That is how sunlight reaches us across 150 million km of empty space.

### Key Wave Properties

Every wave has three core properties:

| Property | Symbol | Unit | What it means |
|----------|--------|------|----------------|
| **Frequency** | f | Hz (hertz) | How many full wave cycles pass a point per second |
| **Wavelength** | λ (lambda) | meters (m) | Distance from one wave crest to the next |
| **Speed** | c | m/s | How fast the wave travels through space |

In a vacuum, all electromagnetic waves (radio, light, X-rays, everything) travel at the same speed. This is the **speed of light**:

```
c = 3 × 10⁸ m/s   (approximately)
```

The exact value is 299,792,458 m/s. Physicists use c as the symbol because the Latin word for speed is *celeritas*.

### The Wave Equation

These three properties are connected by a simple but powerful equation:

```
c = f × λ
```

This says: the speed of a wave equals its frequency multiplied by its wavelength. If you know any two, you can find the third.

- Rearranged to find frequency:    f = c / λ
- Rearranged to find wavelength:   λ = c / f

> **Key insight:** When light travels from one medium to another (say, from air into glass), its **frequency stays constant** but its **wavelength changes**. The wave equation still holds, but with a different speed in the medium. More on this in Section 7.

---

### Worked Example 1: Finding Frequency from Wavelength

**Problem:** Green light has a wavelength of 500 nm. What is its frequency?

**Step 1 — Convert wavelength to meters:**

```
500 nm = 500 × 10⁻⁹ m = 5 × 10⁻⁷ m
```

**Step 2 — Use the wave equation:**

```
f = c / λ
f = (3 × 10⁸ m/s) / (5 × 10⁻⁷ m)
f = (3 / 5) × 10⁸⁺⁷ Hz
f = 0.6 × 10¹⁵ Hz
f = 6 × 10¹⁴ Hz
```

**Answer:** Green light oscillates at 600 trillion times per second. Your eye detects this frequency and your brain interprets it as the color green.

---

### Light as a Photon

Waves are not the complete picture. In the early 1900s, Albert Einstein and Max Planck showed that light also behaves like a stream of tiny packets of energy called **photons**. Each photon carries a fixed amount of energy that depends on its frequency:

```
E = h × f
```

Where:
- **E** = energy of one photon (joules, J)
- **h** = Planck's constant = 6.63 × 10⁻³⁴ J·s
- **f** = frequency (Hz)

Higher frequency = higher energy per photon. This is why UV light can damage your skin (high energy photons) but radio waves cannot (very low energy photons). We will explore photons in detail in later chapters on quantum physics. For now, remember: light is both a wave and a particle — this is called **wave-particle duality**.

---

## 2. The Visible Spectrum

Human eyes can only detect a tiny slice of the full electromagnetic spectrum. This slice is called **visible light**, and it spans wavelengths from roughly **400 nm to 700 nm**.

### Colors and Wavelengths

Different wavelengths within visible light are perceived as different colors:

```
Violet │ Indigo │ Blue │ Green │ Yellow │ Orange │ Red
  400nm ─────────────────────────────────────────── 700nm
(higher f)                                    (lower f)
(more energy)                             (less energy)
```

The relationship to remember:

- **Shorter wavelength → higher frequency → more energy per photon**
- **Longer wavelength → lower frequency → less energy per photon**

Red light has about half the energy per photon compared to violet light, even though both are visible to us.

A useful mnemonic for the order of colors from long to short wavelength is **ROY G BIV** (Red, Orange, Yellow, Green, Blue, Indigo, Violet), which is the order you see in a rainbow.

### Beyond the Visible: The Full EM Spectrum

Visible light is only a tiny window. Here is the full picture:

```
The Electromagnetic Spectrum
─────────────────────────────────────────────────────────────────────────

 LONGER WAVELENGTH ◄──────────────────────────────► SHORTER WAVELENGTH
 LOWER FREQUENCY   ◄──────────────────────────────► HIGHER FREQUENCY
 LOWER ENERGY      ◄──────────────────────────────► HIGHER ENERGY

  Radio    Microwave  Infrared  VISIBLE   UV      X-ray   Gamma Ray
  ┌──────┐  ┌───────┐  ┌──────┐  ┌─────┐  ┌────┐  ┌────┐  ┌────────┐
  │ >1 m │  │1m–1mm │  │1mm– │  │700– │  │400-│  │10nm│  │<0.01nm │
  │      │  │       │  │700nm│  │400nm│  │10nm│  │–   │  │        │
  └──────┘  └───────┘  └──────┘  └─────┘  └────┘  │0.01│  └────────┘
                                                    │nm  │
                                                    └────┘

 Examples:
 Radio   → FM/AM radio, TV broadcasts
 Micro   → Wi-Fi, microwave ovens, radar
 IR      → TV remotes, heat lamps, night vision
 VISIBLE → what our eyes see
 UV      → causes sunburn, germicidal lamps, fluorescent dyes
 X-ray   → medical imaging, security scanners
 Gamma   → nuclear reactions, cancer treatment, sterilization

─────────────────────────────────────────────────────────────────────────
```

### Full EM Spectrum Reference Table

| Type | Wavelength Range | Frequency Range | Key Uses |
|------|-----------------|-----------------|----------|
| Radio | > 1 m | < 3 × 10⁸ Hz | Broadcasting, communications |
| Microwave | 1 mm – 1 m | 3×10⁸ – 3×10¹¹ Hz | Wi-Fi, radar, ovens |
| Infrared (IR) | 700 nm – 1 mm | 3×10¹¹ – 4×10¹⁴ Hz | Heat sensing, remote controls |
| Visible | 400 – 700 nm | 4×10¹⁴ – 7.5×10¹⁴ Hz | Human vision |
| Ultraviolet (UV) | 10 – 400 nm | 7.5×10¹⁴ – 3×10¹⁶ Hz | Sterilization, tanning |
| X-ray | 0.01 – 10 nm | 3×10¹⁶ – 3×10¹⁹ Hz | Medical imaging |
| Gamma ray | < 0.01 nm | > 3×10¹⁹ Hz | Nuclear physics, cancer therapy |

---

## 3. Sources of Light

Where does light come from? There are several different physical mechanisms that produce light. Understanding them helps us see why different light sources look and behave differently.

### Incandescent Bulbs (Thermal Emission)

An incandescent bulb works by passing electric current through a thin tungsten wire called a **filament**. The electrical resistance of the wire causes it to heat up to around **2700 K** (about 2400°C). Any object at this temperature radiates light — this is called **thermal emission** or **blackbody radiation**.

- Produces a **continuous spectrum**: all wavelengths from red through violet (and beyond into infrared and UV)
- Very **inefficient**: only about 5% of the electrical energy becomes visible light; the other ~95% becomes heat (infrared radiation)
- The color of the light is warm (yellowish-white), because at 2700 K the peak emission is in the infrared/red region

Thermal emission is why things glow when heated: a metal poker in a fire first glows dull red, then bright orange-yellow as it gets hotter. The Sun, at ~5778 K, peaks in the yellow-green part of the visible spectrum.

### Fluorescent Tubes

A fluorescent tube contains **mercury vapor** at low pressure. When electrical current passes through it, mercury atoms get excited and emit **ultraviolet (UV) light**. The UV light then hits a **phosphor coating** on the inside of the glass tube. The phosphor absorbs UV photons and re-emits **visible light** — a process called **fluorescence**.

- More efficient than incandescent bulbs (~25% electrical energy → visible light)
- Produces a somewhat **line spectrum** (specific colors), though the phosphor blend smooths this out
- The characteristic flicker and slight color difference versus incandescent come from this mechanism

### LED (Light Emitting Diode)

LEDs work on a completely different principle called **electroluminescence**. An LED is made of a **semiconductor** — a material like gallium nitride or gallium arsenide. The semiconductor has a **band gap**: a gap between energy levels that electrons can occupy.

When electrons in the LED fall across this band gap, they release energy in the form of a photon. The energy of the photon (and therefore its wavelength/color) is determined by the size of the band gap:

```
E_photon = h × f = band gap energy
```

- Different semiconductor materials have different band gaps → different colors of LED
- Blue LEDs + phosphor coating → white LED (the key invention that earned the 2014 Nobel Prize in Physics)
- Very **efficient**: >90% of electrical energy can become light
- Very **long-lasting**: no filament to burn out, no glass to break
- Compact, durable, and can be made in virtually any color

### Laser

**LASER** stands for **Light Amplification by Stimulated Emission of Radiation**. Lasers produce light through a process called **stimulated emission**: an incoming photon triggers an excited atom to release a second photon with exactly the same energy, direction, and phase.

Laser light has three special properties that ordinary light does not:

1. **Coherent**: All the waves are in phase (their crests and troughs line up perfectly). This allows lasers to form interference patterns.
2. **Monochromatic**: All the photons have exactly the same wavelength (one pure color).
3. **Collimated**: The beam is extremely parallel — it barely spreads out over long distances.

```
Ordinary Lamp:                 Laser:
~~~~~→                         ──────────────────────►
  ~~~→  (random phases,        ──────────────────────►  (all same phase,
 ~~~~→   many wavelengths,     ──────────────────────►   one wavelength,
~~~~→    spreads in all        ──────────────────────►   tight beam)
  ~~~→   directions)
```

Applications of lasers:
- Barcode scanners in supermarkets
- Reading CDs, DVDs, Blu-ray discs
- Laser surgery (LASIK eye surgery, cancer removal)
- Fiber optic communications (pulses of laser light carry internet data)
- Laser pointers, laser levels (construction)
- Scientific measurement (LIDAR, distance measurement, spectroscopy)

### Natural Sources

- **The Sun**: Energy comes from **nuclear fusion** — hydrogen atoms fuse to form helium in the core, releasing enormous energy. This energy eventually reaches the surface as thermal radiation, producing the continuous spectrum of sunlight.
- **Fire**: Thermal emission from hot soot particles (giving yellow-orange glow) and chemiluminescence from excited molecules in the flame (giving the blue base).
- **Bioluminescence**: Some living organisms (fireflies, deep-sea fish, certain jellyfish) produce light through **chemical reactions** in their bodies. The reaction produces excited molecules that release photons as they fall back to their ground state. No heat is produced — it is called "cold light."

---

## 4. Light Travels in Straight Lines

One of the most fundamental properties of light — and the basis of all of geometric optics — is that **light travels in straight lines** in a uniform medium.

This is called the **rectilinear propagation of light**.

### Why Does This Work?

Strictly speaking, light (being a wave) can diffract (bend around corners). However, diffraction is only significant when the size of the opening or obstacle is comparable to the wavelength of the light. Since visible light has wavelengths of only 400–700 nm (nanometers), everyday objects (millimeters to meters in size) are enormously larger than the wavelength. In this regime, light behaves as if it travels in straight lines — this is the **geometric optics approximation**.

When you need to deal with very tiny slits or holes (comparable to light's wavelength), you must use **wave optics** instead. We will cover that in later chapters.

### Evidence for Straight-Line Propagation

1. **Laser beams**: You can see the beam travel in a perfectly straight line (especially visible in dusty air or fog).

2. **Shadows**: Objects block light and cast sharp shadows only because light cannot bend around them. If light curved around corners, there would be no shadows.

3. **The pinhole camera** (see Section 6): The fact that a pinhole camera forms an image — and an inverted one at that — is direct proof that light travels in straight lines.

### Light Rays

In geometric optics, we draw **light rays** — straight arrows representing the direction light is traveling. A ray is an idealization (light is really a wave, not a geometric line), but it is an extremely useful model for predicting where light goes.

```
     Light source
         ★
          \  ← ray 1
           \
            \
         ────□──────── (obstacle with small hole)
             |
             | ← ray continues in same direction
             ↓
          (screen)
          dot of light
```

---

## 5. Shadows: Umbra and Penumbra

A **shadow** is a region where light is blocked by an opaque object. The nature of the shadow depends on whether the light source is a **point source** or an **extended source**.

### Point Source → Sharp Shadow

A **point source** is a light source so small it can be treated as a single point. All rays from it spread out from one spot.

When a point source illuminates an opaque object, the shadow has only one region: the **umbra**, where no light reaches at all. The edge of the shadow is perfectly sharp.

```
Point Source Shadow Formation:
─────────────────────────────

  ★  ← (point source)
  |\ 
  | \
  |  \
  |   ■ ■ ■ ← opaque object
  |  /|    \
  | / |     \
  |/  |      \
  |   ↓       ↓
  |  UMBRA   fully lit
  |  (no light reaches here)
```

### Extended Source → Umbra + Penumbra

The **Sun** and most practical light sources are **extended sources** — they have a physical size. Different parts of the extended source illuminate different parts of the space around the object.

- **Umbra** (Latin for "shadow"): The region that receives light from **no part** of the source. Total darkness.
- **Penumbra** (Latin for "almost shadow"): The region that receives light from **some but not all** parts of the source. Partial shadow, getting gradually brighter as you move away from the umbra.

```
Extended Source Shadow Formation:
──────────────────────────────────────────────────

  ┌────────────┐
  │  Extended  │  ← Light source (e.g., the Sun)
  │   Source   │
  └────────────┘
   \          /
    \        /
     \      /  ← rays from edges of source
      \    /
       ■■■■  ← opaque object (e.g., the Moon)
      /|  |\
     / |  | \
    /  |  |  \
   /  UMBRA   \
  /  (no light)\
 /──────────────\
│  PENUMBRA     │  ← partial shadow, some light gets through
│               │     from one side of the source
```

### Real-World Example: Solar Eclipse

The Moon passing between the Earth and the Sun creates:

- **Umbra**: A narrow cone (~100–170 km wide on Earth's surface). Observers here see the Sun **completely blocked** — a **total solar eclipse**. It goes dark enough to see stars.
- **Penumbra**: A much larger region thousands of km wide. Observers here see the Sun **partially blocked** — a **partial solar eclipse**. The sky dims but does not go dark.

Because the Moon's orbit is slightly elliptical, sometimes the Moon is farther from Earth and its umbra doesn't quite reach the surface — observers on the centerline see an **annular eclipse** ("ring of fire"), where the Moon's disk is completely inside the Sun's disk.

### Why Shadows Change in Size

When an object moves **closer to the light source** and **farther from the screen**, its shadow grows larger. Here is why:

Light rays spread out from the source. The farther the object is from the screen (or the closer it is to the source), the more the shadow-rays have spread by the time they hit the screen, creating a bigger shadow.

```
   ★  ←─── light source

   ■  ←─── object close to source
   │
   │
   ────────  ← screen
   (large shadow)

   vs.

   ★  ←─── same light source


         ■  ←─── object close to screen
         ────────  ← screen
   (small shadow)
```

---

## 6. Pinhole Camera (Camera Obscura)

The **pinhole camera** (from the Latin *camera obscura*, meaning "dark room") is one of the oldest optical devices known to humanity. It works entirely because light travels in straight lines.

### How It Works

Take a light-tight box with a tiny hole (pinhole) in one side and a white screen on the opposite side. When you point the pinhole at an object in bright light:

- Light from every point on the object travels in all directions.
- Only one ray from each point on the object can pass through the tiny pinhole.
- That ray continues in a straight line and hits a specific point on the screen.
- The point that was at the top of the object sends a ray through the hole that hits the **bottom** of the screen.
- The point at the bottom of the object sends a ray through the hole that hits the **top** of the screen.

**Result**: The image is **inverted** — upside down AND left-right flipped.

```
Pinhole Camera Ray Diagram:
────────────────────────────────────────────────────────

     Object        Pinhole          Screen
       │               │                │
  ─────┼─────          │           ─────┼─────
  A    │               │                │    A' (top of object → bottom of image)
       │               │                │
  ─────┼─────          │           ─────┼─────
  B ───│──────────────►●────────────────►────  B' (center → center)
       │               │                │
  ─────┼─────          │           ─────┼─────
  C    │               │                │    C' (bottom of object → top of image)
       │               │                │

  A (top)    →  ray goes down through hole →  lands at bottom  →  A' (bottom)
  C (bottom) →  ray goes up through hole  →  lands at top     →  C' (top)

  Image is inverted (upside down and left-right reversed).
```

### Hole Size Matters

- **Smaller hole → sharper image**: Each point on the object sends only one ray through the tiny hole, making a sharp point on the screen. But less light gets through → dimmer image.
- **Larger hole → brighter but blurred image**: Each point on the object can now send many rays through the larger hole (a cone of rays). These land on a small circular region on the screen instead of a single point → blurred.
- **Modern cameras** solve this with a lens, which bends all the rays from one point back to a single focus, giving both sharpness and brightness.

### Magnification Formula

The **magnification** of a pinhole camera is simply the ratio of image size to object size, and it equals the ratio of distances:

```
Magnification (M) = image height (h_i) / object height (h_o)
                  = image distance (v) / object distance (u)
```

Or to find the image height directly:

```
h_i = h_o × (v / u)
```

---

### Worked Example 2: Pinhole Camera Image Size

**Problem:** A tree is 4 m tall and stands 10 m from a pinhole camera. The screen inside the camera is 0.4 m behind the pinhole. How tall is the image of the tree on the screen?

**Given:**
- Object height h_o = 4 m
- Object distance u = 10 m
- Image distance v = 0.4 m

**Using the magnification formula:**

```
h_i = h_o × (v / u)
h_i = 4 × (0.4 / 10)
h_i = 4 × 0.04
h_i = 0.16 m = 16 cm
```

**Answer:** The image of the tree on the screen is 16 cm tall — and it is upside down.

---

### Worked Example 3: Finding Object Distance

**Problem:** A pinhole camera has a screen 25 cm from the pinhole. A person 1.8 m tall produces an image 9 cm tall. How far is the person from the pinhole?

**Given:**
- Image height h_i = 9 cm = 0.09 m
- Object height h_o = 1.8 m
- Image distance v = 25 cm = 0.25 m

**Rearranging the formula:**

```
h_i / h_o = v / u
u = v × (h_o / h_i)
u = 0.25 × (1.8 / 0.09)
u = 0.25 × 20
u = 5 m
```

**Answer:** The person is 5 m from the pinhole.

---

## 7. Speed of Light in a Medium

### Speed in a Vacuum

In a perfect vacuum, all electromagnetic waves travel at:

```
c = 3 × 10⁸ m/s  (exact: 299,792,458 m/s)
```

This is the universal speed limit. Nothing with mass can reach this speed; only massless particles like photons can travel at exactly c.

### Slowing Down in a Medium

When light enters a material (glass, water, diamond, etc.), it interacts with the atoms of that material. The light is constantly being absorbed and re-emitted by atoms, which slows the overall propagation. The speed of light in a medium is:

```
v = c / n
```

Where **n** is the **refractive index** (also called the **index of refraction**) of the material. The refractive index is always greater than or equal to 1.

- **n = 1** means light travels at the full speed c (vacuum).
- **n = 1.5** means light travels at c/1.5 = two-thirds the speed of light.
- **Higher n → slower light**.

### Refractive Index Values for Common Materials

| Material | Refractive Index (n) | Speed of Light (v) |
|----------|---------------------|-------------------|
| Vacuum | 1.000 | 3.00 × 10⁸ m/s |
| Air | ≈ 1.0003 | ≈ 3.00 × 10⁸ m/s |
| Water (liquid) | 1.33 | 2.26 × 10⁸ m/s |
| Crown glass | 1.52 | 1.97 × 10⁸ m/s |
| Flint glass | 1.62 | 1.85 × 10⁸ m/s |
| Diamond | 2.42 | 1.24 × 10⁸ m/s |

Air is so close to vacuum (n ≈ 1.0003) that for most purposes we treat it as if light travels at full speed c in air.

### What Changes and What Stays the Same?

When light enters a medium from vacuum (or air):

- **Frequency stays the same** — the number of oscillations per second doesn't change. This makes physical sense: every wave crest that enters the medium must exit the other side, so the frequency at the boundary must match on both sides.

- **Speed decreases** (v = c/n, slower)

- **Wavelength decreases** — since c = f × λ still applies (just with speed v and wavelength λ_medium):

```
v = f × λ_medium
λ_medium = v / f = (c/n) / f = λ_vacuum / n
```

So the wavelength in the medium is shorter by a factor of n.

This is important in optics and is why **prisms disperse white light into a spectrum**: different wavelengths have slightly different refractive indices in glass (n is slightly higher for violet than red), so they slow down by different amounts and bend at slightly different angles.

---

### Worked Example 4: Speed of Light in Diamond

**Problem:** The refractive index of diamond is 2.42. What is the speed of light inside a diamond?

**Using:** v = c / n

```
v = (3 × 10⁸ m/s) / 2.42
v = 1.24 × 10⁸ m/s
```

**Answer:** Light travels at about 1.24 × 10⁸ m/s in diamond — less than half its vacuum speed. This large refractive index is what gives diamonds their brilliant sparkle (total internal reflection, covered in a later chapter).

---

### Worked Example 5: Wavelength in Water

**Problem:** Yellow light has a wavelength of 590 nm in air. What is its wavelength inside water (n = 1.33)?

**Using:** λ_medium = λ_vacuum / n

```
λ_water = 590 nm / 1.33
λ_water = 443 nm
```

**Answer:** The wavelength shortens to 443 nm inside water, but the frequency (and therefore the color your eye perceives) stays the same.

---

## 8. Coherent vs Incoherent Light

### Incoherent Light

Most everyday light sources — incandescent bulbs, fluorescent lamps, the Sun, candles — produce **incoherent light**. This means:

- They emit **many different wavelengths** (many colors mixed together)
- The waves have **random phase relationships** — different parts of the source emit waves that are out of step with each other in an unpredictable way

```
Incoherent light:
~~~~~►
  ~~~►     (different wavelengths,
   ~~~~►    random phases,
 ~~~~~►     spreads in all
    ~~~►     directions)
```

Because the phases are random and constantly changing, the waves from different parts of the source cannot consistently reinforce or cancel each other. The result is a steady, uniform illumination with no interference pattern.

### Coherent Light

**Coherent light** has two key properties:

1. **Same wavelength** (monochromatic): All waves have the same frequency/color.
2. **Fixed phase relationship**: The waves all have the same phase (or a constant, predictable phase difference). They are "in step" with each other.

```
Coherent light (laser):
──────────────────────────────────────►
──────────────────────────────────────►  (same wavelength,
──────────────────────────────────────►   same phase,
──────────────────────────────────────►   all waves aligned)
```

A **laser** is the primary source of coherent light. The process of stimulated emission (see Section 3) ensures all photons are emitted with the same wavelength and phase.

### Why Coherence Matters

Coherent light can produce **interference patterns** — regions of constructive interference (bright bands) and destructive interference (dark bands) when two coherent beams overlap. This is only possible because the phase relationship is stable and predictable.

Applications that require coherence:

- **Holograms**: A hologram is recorded by splitting a laser beam, directing one part at an object and combining it with the other part. The interference pattern records a 3D image. Incoherent light cannot make holograms.
- **Interferometry**: Measuring tiny distances (down to fractions of a nanometer) by detecting interference fringes. Used in LIGO (gravitational wave detector) to measure distances 1/1000th the width of a proton.
- **Fiber optic communications**: Laser pulses are used because they stay focused and intense over long distances.
- **Optical coherence tomography (OCT)**: A medical imaging tool that uses coherent light to image the retina of the eye in 3D.

**Spatial coherence** refers to whether different points across the beam are in phase with each other. A laser has high spatial coherence; a light bulb does not.

**Temporal coherence** refers to how long a wave train maintains its phase — related to how monochromatic the light is. A laser has high temporal coherence; a white light source has very low temporal coherence.

---

## 9. Polarization (Preview)

### How Normal Light Oscillates

Recall that in a light wave, the electric field oscillates perpendicular to the direction of travel. But perpendicular means there are infinitely many possible orientations — up-down, left-right, diagonal, and everything in between.

In **unpolarized light**, the electric field oscillates in **all directions perpendicular to travel** at random, constantly changing. If you looked head-on at a beam of unpolarized light, the electric field would be pointing in all directions simultaneously (or rather, averaging over all directions very rapidly).

```
Unpolarized light (viewed head-on):
             ↑
        ↖    │    ↗
        ←────●────→
        ↙    │    ↘
             ↓
   (electric field in all directions)
```

### Polarized Light

**Polarized light** is light where the electric field oscillates in only **one specific plane** (or in a controlled pattern). The direction of this oscillation is called the **plane of polarization** or **polarization direction**.

```
Polarized light (viewed head-on):
        ↑
        │
        │
        ●         (electric field only up-down)
        │
        │
        ↓

   or

        ←────●────→  (electric field only left-right)
```

### How to Polarize Light: Polaroid Filters

A **polaroid filter** (or polarizer) contains long-chain molecules aligned in one direction. It transmits only the component of the electric field that is aligned with a specific axis (called the **transmission axis**) and absorbs the perpendicular component.

When unpolarized light passes through a polarizer:
- About 50% of the light is transmitted (the component aligned with the transmission axis)
- The transmitted light is now polarized along the transmission axis

```
Unpolarized →  Polarizer  → Polarized (vertical)
  light       (vertical      light
               axis)

  ↑↓←→↗↘      │ filter │      ↑
   (all         │       │      │   only vertical
  directions)   └───────┘      ↓   component passes
```

### Crossed Polarizers

If you place two polaroid filters one after the other with their transmission axes at **90° to each other** ("crossed"), **no light gets through**:

- First polarizer: allows only vertically polarized light
- Second polarizer (axis horizontal, 90° from first): blocks all vertically polarized light

```
Unpolarized →  Polarizer 1   → Polarized  →  Polarizer 2  →  No light
  light        (vertical)       (vertical)    (horizontal,    passes
                                              90° to first)   through
```

If you rotate the second polarizer, the transmitted intensity follows **Malus's Law**:

```
I = I_0 × cos²(θ)
```

Where θ is the angle between the two polarizer axes, and I_0 is the intensity after the first polarizer. At θ = 0° (aligned), all light passes; at θ = 90°, no light passes.

### How Light Becomes Polarized Naturally

Light can become partially or fully polarized without a filter:

1. **Reflection**: Light reflected off a flat surface (water, glass, wet roads) becomes partially horizontally polarized. This is the glare you see off puddles and car hoods.

2. **Scattering**: Light scattered by the atmosphere (Rayleigh scattering) is partially polarized. This is why the sky looks different through polarized sunglasses when you look at different parts of it.

### Applications of Polarization

- **Polarized sunglasses**: The lenses have a vertical transmission axis. Glare from horizontal surfaces is horizontally polarized, so polarized lenses block it. This dramatically reduces glare from roads, water, and snow.

- **LCD screens** (Liquid Crystal Displays): Every pixel in your phone or computer screen uses two crossed polarizers with a liquid crystal layer between them. Applying a voltage to a pixel rotates the polarization of light, allowing it to pass through (white/bright pixel) or be blocked (dark pixel). LCD screens don't work if you look at them through a polarized filter at the wrong angle — you may have noticed your phone screen going black when wearing polarized sunglasses and tilting your head.

- **Photography**: A **circular polarizing filter** on a camera lens reduces reflections from glass and water, makes blue skies appear deeper, and increases color saturation.

- **3D cinema**: Old red-green anaglyphic 3D glasses used color filtering; modern "Real D 3D" cinema uses circularly polarized light — the left and right projections use opposite circular polarizations. The glasses for each eye filter out one polarization, so each eye sees a slightly different image, creating the illusion of depth.

- **Stress analysis**: When transparent materials (like plastic or glass) are under stress, they rotate polarized light differently in different regions. Engineers place stressed models between crossed polarizers to see colorful stress patterns that reveal where a structure might break.

---

## Summary

- **Light is an electromagnetic wave**: oscillating electric and magnetic fields, perpendicular to each other and to the direction of travel. It also behaves as photons (wave-particle duality).

- **Wave properties**: frequency (f, Hz), wavelength (λ, m), speed (c = 3 × 10⁸ m/s in vacuum). Related by c = f × λ.

- **Photon energy**: E = h × f, where h = 6.63 × 10⁻³⁴ J·s. Higher frequency = more energetic photons.

- **Visible spectrum**: 400 nm (violet) to 700 nm (red). ROY G BIV from long to short wavelength. UV is beyond violet (higher energy), IR is beyond red (lower energy).

- **Full EM spectrum** ranges from radio waves (meters long) to gamma rays (sub-nanometer), all traveling at c in vacuum.

- **Sources of light**: Incandescent (thermal emission, inefficient), Fluorescent (UV → phosphor → visible, more efficient), LED (electroluminescence across semiconductor band gap, very efficient), Laser (stimulated emission, coherent + monochromatic + collimated). Natural: Sun (fusion), fire (thermal + chemical), bioluminescence (chemical).

- **Rectilinear propagation**: Light travels in straight lines in a uniform medium (geometric optics approximation valid when objects >> λ).

- **Shadows**: Point source → sharp umbra only. Extended source → umbra (total shadow) + penumbra (partial shadow). Larger distance from screen → larger shadow.

- **Pinhole camera**: Light through a tiny hole forms an inverted image on a screen. Image height = object height × (image distance / object distance). Smaller hole = sharper but dimmer image.

- **Refractive index**: n = c / v. Light slows to v = c/n in a medium. Frequency stays constant; wavelength shortens to λ/n. Higher n → slower, shorter wavelength.

- **Coherent vs incoherent**: Coherent light (laser) has same wavelength + fixed phase; can form interference patterns. Incoherent light (lamp) has random phases and multiple wavelengths.

- **Polarization**: Unpolarized light has electric field in all perpendicular directions. Polarized light has electric field in one plane. Polaroid filters transmit one polarization axis. Crossed polarizers block all light. Applications: sunglasses, LCD screens, 3D movies, photography.

---

## Key Equations

**Wave equation (speed, frequency, wavelength):**
```
c = f × λ

f = c / λ

λ = c / f
```

**Speed of light in vacuum:**
```
c = 3 × 10⁸ m/s  (exact: 299,792,458 m/s)
```

**Photon energy:**
```
E = h × f

h = 6.63 × 10⁻³⁴ J·s  (Planck's constant)
```

**Speed of light in a medium:**
```
v = c / n

n = c / v
```

**Wavelength in a medium:**
```
λ_medium = λ_vacuum / n
```

**Pinhole camera magnification:**
```
M = h_i / h_o = v / u

h_i = h_o × (v / u)
```

**Malus's Law (intensity through a polarizer):**
```
I = I_0 × cos²(θ)

where θ = angle between polarizer transmission axes
```

---

**Quick Reference: Refractive Indices**

| Material | n |
|----------|---|
| Vacuum | 1.000 |
| Air | ~1.000 |
| Water | 1.33 |
| Crown glass | 1.52 |
| Flint glass | 1.62 |
| Diamond | 2.42 |

---

**Quick Reference: Visible Spectrum Colors**

| Color | Approximate Wavelength | Approximate Frequency |
|-------|----------------------|----------------------|
| Violet | 400–420 nm | 7.1–7.5 × 10¹⁴ Hz |
| Indigo | 420–445 nm | 6.7–7.1 × 10¹⁴ Hz |
| Blue | 445–500 nm | 6.0–6.7 × 10¹⁴ Hz |
| Green | 500–565 nm | 5.3–6.0 × 10¹⁴ Hz |
| Yellow | 565–590 nm | 5.1–5.3 × 10¹⁴ Hz |
| Orange | 590–625 nm | 4.8–5.1 × 10¹⁴ Hz |
| Red | 625–700 nm | 4.3–4.8 × 10¹⁴ Hz |

---

*End of Chapter 61. Next chapter: Reflection of Light — plane mirrors, laws of reflection, image formation.*
