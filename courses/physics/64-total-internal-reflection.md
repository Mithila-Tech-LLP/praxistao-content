# Chapter 64: Total Internal Reflection and Optical Fibers

> **"A beam of light can be trapped inside a glass fiber and guided thousands of kilometers around the world. This is the backbone of the internet — and it works because of a simple physics phenomenon."**

---

## Table of Contents

1. [Recap: Refraction at a Dense-to-Less-Dense Boundary](#1-recap-refraction-at-a-dense-to-less-dense-boundary)
2. [The Critical Angle](#2-the-critical-angle)
3. [Total Internal Reflection (TIR)](#3-total-internal-reflection-tir)
4. [Worked Examples: Finding Critical Angles](#4-worked-examples-finding-critical-angles)
5. [Optical Fibers — Structure and How They Work](#5-optical-fibers--structure-and-how-they-work)
6. [Types of Optical Fiber: Single-Mode and Multi-Mode](#6-types-of-optical-fiber-single-mode-and-multi-mode)
7. [Fiber Optic Applications](#7-fiber-optic-applications)
8. [How the Internet Uses Fiber Optics](#8-how-the-internet-uses-fiber-optics)
9. [Diamonds and Sparkle](#9-diamonds-and-sparkle)
10. [Prisms and TIR](#10-prisms-and-tir)
11. [Other Applications of TIR](#11-other-applications-of-tir)
12. [Summary](#12-summary)
13. [Key Equations](#13-key-equations)

---

## 1. Recap: Refraction at a Dense-to-Less-Dense Boundary

In the previous chapters, we studied **refraction** — the bending of light as it crosses the boundary between two materials with different **refractive indices**. We used **Snell's Law**:

    n1 × sin(θ1) = n2 × sin(θ2)

Where n1 and n2 are the refractive indices of the two media, and θ1 and θ2 are the angles of incidence and refraction measured from the **normal** (a line perpendicular to the boundary).

Now, something very special happens when light travels from a **denser medium** (higher n) into a **less dense medium** (lower n).

When n1 > n2, Snell's Law tells us that sin(θ2) > sin(θ1), which means **θ2 > θ1**. The refracted ray bends *away* from the normal.

Let's watch what happens as we gradually increase the angle of incidence θ1:

### Case 1: Small Angle of Incidence

```
      Medium 1 (glass, n=1.5)   |   Medium 2 (air, n=1.0)
                                 |
     Incident ray                |    Refracted ray
          \                      |         \
           \  θ1 (small)         |          \  θ2 (slightly larger)
            \                    |           \
  -----------\---------BOUNDARY--|-----------\-----------
              \                  |
               \ (reflected ray, |
                  small and      |
                  faint)         |
```

At small angles, there is a refracted ray that exits into the less dense medium. The refracted angle θ2 is larger than θ1, but both are modest. Most light is transmitted; a small portion reflects.

### Case 2: Medium Angle of Incidence

```
      Medium 1 (glass, n=1.5)   |   Medium 2 (air, n=1.0)
                                 |
    Incident ray                 |         Refracted ray
         \                       |         /  (nearly grazing
          \  θ1 (medium)         |        /    the surface)
           \                     |       /  θ2 ≈ 80°
            \                    |      /
  ----------\---------BOUNDARY---|-----/------------------
              \                  |
               \                 |
```

As θ1 increases, θ2 increases even faster. The refracted ray is now nearly parallel to the boundary — it is "grazing" the surface.

### Case 3: At the Critical Angle

```
      Medium 1 (glass, n=1.5)   |   Medium 2 (air, n=1.0)
                                 |
   Incident ray                  |
        \                        |
         \  θ1 = θ_c             |
          \   (critical angle)   |
           \                     |
  ----------\--------BOUNDARY----|============================
              \                  |
               \                 |  Refracted ray travels
                \                |  EXACTLY along the boundary
                                    θ2 = 90°
```

At one specific angle of incidence called the **critical angle** (θ_c), the refracted ray bends so far that it travels exactly along the boundary between the two media. The angle of refraction is 90°.

And what happens if we increase θ1 even further past the critical angle? Something remarkable — read on.

---

## 2. The Critical Angle

The **critical angle** is the angle of incidence at which the refracted ray travels exactly along the boundary (θ2 = 90°).

To find it, we substitute θ2 = 90° into Snell's Law:

    n1 × sin(θ_c) = n2 × sin(90°)

Since sin(90°) = 1:

    n1 × sin(θ_c) = n2

    sin(θ_c) = n2 / n1

This is the formula for the critical angle.

**Important conditions:**
- This only works when n1 > n2 (light going from denser to less dense medium)
- If n1 < n2, there is no critical angle — the refracted ray always bends toward the normal

**Special case — when Medium 2 is air (n ≈ 1.0):**

    sin(θ_c) = 1 / n1

    θ_c = arcsin(1 / n1)

This simplification works whenever light is traveling from any medium into air or vacuum.

**Intuition:** The higher the refractive index n1 of the denser medium, the smaller the critical angle. This means light in a very dense material (like diamond) can be totally internally reflected at surprisingly shallow angles.

---

## 3. Total Internal Reflection (TIR)

When the angle of incidence θ1 **exceeds** the critical angle θ_c, something dramatic happens: **no refracted ray exists at all**. Every single photon is reflected back into the denser medium.

This is called **Total Internal Reflection (TIR)**.

```
      Medium 1 (glass, n=1.5)   |   Medium 2 (air, n=1.0)
                                 |
  Incident ray     Reflected ray |     NO refracted ray!
        \              /         |
         \  θ1 > θ_c  /          |      (No light escapes)
          \          /           |
           \        /            |
  ----------\------/-BOUNDARY----|-----------------------------
```

The reflected ray obeys the **law of reflection**: the angle of reflection equals the angle of incidence (measured from the normal). The reflection is perfect.

Why is TIR called "total"? Because **100% of the light energy is reflected**. No light is lost. Ordinary mirrors reflect only about 88–95% of incoming light (some is absorbed by the metal coating). But TIR reflects 100%, with no metal needed. This makes TIR mirrors and optical fibers extremely efficient.

**The three regimes summarized:**

| Angle of Incidence (θ1) | What Happens |
|---|---|
| θ1 < θ_c | Most light refracts into medium 2; some reflects |
| θ1 = θ_c | Refracted ray travels along boundary (θ2 = 90°) |
| θ1 > θ_c | **Total Internal Reflection** — all light stays in medium 1 |

---

## 4. Worked Examples: Finding Critical Angles

### Example 1: Glass to Air (n_glass = 1.5, n_air = 1.0)

**Given:**
- n1 = 1.5 (glass)
- n2 = 1.0 (air)

**Find:** Critical angle θ_c

**Solution:**

    sin(θ_c) = n2 / n1 = 1.0 / 1.5 = 0.667

    θ_c = arcsin(0.667) = 41.8°

**Interpretation:** Any ray of light inside glass that hits a glass-air boundary at an angle greater than 41.8° (from the normal) will be completely reflected back into the glass. This is the basis of optical fibers made of glass.

---

### Example 2: Water to Air (Snell's Window for Fish)

**Given:**
- n1 = 1.33 (water)
- n2 = 1.0 (air)

**Find:** Critical angle θ_c

**Solution:**

    sin(θ_c) = 1.0 / 1.33 = 0.752

    θ_c = arcsin(0.752) = 48.8°

**Interpretation:** From beneath the water's surface, a fish can see the entire above-water world (sky, trees, people) compressed into a circular cone of half-angle 48.8°. This is called **Snell's Window** or the **Snell's cone**. Everything outside that 48.8° cone appears completely dark (TIR from below) or shows reflections of the underwater scene.

If you've ever looked at the surface of a swimming pool from underwater, you've seen this effect — there's a bright circular "window" in the center, surrounded by a mirror-like reflection of the pool floor.

---

### Example 3: Diamond to Air (n_diamond = 2.42)

**Given:**
- n1 = 2.42 (diamond)
- n2 = 1.0 (air)

**Find:** Critical angle θ_c

**Solution:**

    sin(θ_c) = 1.0 / 2.42 = 0.413

    θ_c = arcsin(0.413) = 24.4°

**Interpretation:** Diamond has an extremely small critical angle of just 24.4°. This means light rattling around inside a diamond is almost always hitting a surface at an angle greater than 24.4° — and is therefore totally internally reflected. A skilled diamond cutter uses this property deliberately. We'll explore this more in Section 9.

---

### Example 4: Glass to Water (Not Glass to Air)

**Given:**
- n1 = 1.5 (glass)
- n2 = 1.33 (water)

**Find:** Critical angle θ_c

**Solution:**

    sin(θ_c) = n2 / n1 = 1.33 / 1.5 = 0.887

    θ_c = arcsin(0.887) = 62.5°

**Interpretation:** When glass is in contact with water instead of air, the critical angle is much larger (62.5° vs 41.8°). This is because the refractive index difference between glass and water is smaller than between glass and air. The smaller the n1 / n2 ratio, the larger the critical angle.

**Key insight:** TIR is harder to achieve when the two media have similar refractive indices. You need a larger angle of incidence to achieve it.

---

### Example 5: Does TIR occur? (Verification problem)

A ray of light travels inside glass (n = 1.5) and hits the glass-air boundary at an angle of incidence of 50°. Does TIR occur?

**Solution:**

First, find the critical angle:

    sin(θ_c) = 1.0 / 1.5 = 0.667
    θ_c = 41.8°

Since the angle of incidence (50°) is greater than the critical angle (41.8°), **yes, TIR occurs**. The light is completely reflected back into the glass.

---

## 5. Optical Fibers — Structure and How They Work

An **optical fiber** is a thin strand of transparent material (glass or plastic) that can guide light along its length by using total internal reflection, even around curves.

### Structure of an Optical Fiber

```
     Cross-section view (end-on):

     +----------------------------------+
     |           JACKET                 |  ← Protective outer coating
     |   +------------------------+     |     (plastic, ~250 µm total)
     |   |       CLADDING         |     |  ← Outer glass layer, n ≈ 1.45
     |   |   +--------------+    |     |
     |   |   |              |    |     |
     |   |   |    CORE      |    |     |  ← Inner glass/plastic
     |   |   |   (n ≈ 1.5)  |    |     |     (9–62 µm diameter)
     |   |   |              |    |     |
     |   |   +--------------+    |     |
     |   +------------------------+     |
     +----------------------------------+
```

- **Core:** The central region. High refractive index (n ≈ 1.5 for glass). This is where the light travels.
- **Cladding:** A surrounding layer of glass or plastic with a slightly lower refractive index (n ≈ 1.45). The core-cladding boundary is where TIR occurs.
- **Jacket:** An outer plastic coating that protects the fiber mechanically. Does not carry light.

### How Light Travels Through the Fiber

```
     Side-view (fiber cut lengthwise):

     ==========CLADDING (n=1.45)===========
          ------CORE (n=1.5)------
     
     Light enters -->  *        *        *        *--> exits far end
                        \      / \      / \      /
                         \    /   \    /   \    /
                          \  /     \  /     \  /
                    TIR →  \/   TIR \/   TIR \/
     
          ------CORE (n=1.5)------
     ==========CLADDING (n=1.45)===========
     
                  ↑ Each reflection is TIR at core-cladding boundary
                    angle > critical angle θ_c ≈ 73° (from normal)
```

1. Light enters the core at one end.
2. It travels at a slight angle toward the core-cladding boundary.
3. It hits the boundary at an angle **greater than** the critical angle → **TIR**.
4. It reflects to the other side of the core.
5. It hits that boundary → **TIR** again.
6. This repeats thousands of times per meter, all the way to the other end.
7. Light exits the far end, having traveled its entire journey inside the core.

Even if the fiber bends (within limits), TIR continues to work because the geometry of each reflection is approximately preserved. This is why optical fibers can be coiled, routed around corners, and strung under oceans.

### Why Use Cladding Instead of Just Air?

One might wonder: why not just use a bare glass wire in air? The critical angle for glass-air is 41.8°, which is actually smaller — so TIR is even easier to achieve! But cladding is used for several crucial reasons:

1. **Mechanical protection:** Bare glass is fragile. Cladding gives the fiber structural integrity.
2. **Surface cleanliness:** Dust, moisture, or fingerprints on a bare glass surface would scatter light and degrade the signal. Cladding keeps the optical surface permanently clean.
3. **Prevents crosstalk:** In a bundle of bare fibers touching each other, light could "leak" from one fiber into adjacent ones. Cladding prevents this.
4. **Consistent performance:** The n ratio between core and cladding is precisely controlled during manufacturing, ensuring predictable TIR characteristics.

### The Acceptance Cone (Numerical Aperture)

Light must enter the fiber at the right angle to undergo TIR inside. If it enters at too steep an angle (relative to the fiber axis), it will hit the core-cladding boundary at too small an angle — and escape through the cladding instead of reflecting.

The maximum angle at which light can enter the fiber and still undergo TIR is called the **acceptance angle** (θ_max).

The **Numerical Aperture (NA)** describes this property:

    NA = n_outside × sin(θ_max)

For most single-mode glass fibers, NA ≈ 0.12 (acceptance cone half-angle ≈ 6.9°).
For multi-mode fibers, NA ≈ 0.2 to 0.5 (wider acceptance cone).

---

## 6. Types of Optical Fiber: Single-Mode and Multi-Mode

Not all optical fibers are the same. The core diameter makes a huge difference in how the fiber behaves.

### Single-Mode Fiber

```
     Core diameter: ~9 micrometers (about 1/10 the width of a human hair)
     
     ========================CLADDING========================
     ==========CORE (~9 µm)==========
     
     Light path:  --------->  (travels nearly straight down the axis)
     
     ========================CLADDING========================
     
     Only ONE mode (path) of propagation
     Used for: long-distance telecommunications (undersea cables, backbone networks)
```

- **Core diameter:** ~9 micrometers (µm)
- **How light travels:** In a single straight path down the fiber axis, with very little bouncing. The core is so narrow that only one "mode" (path of propagation) fits.
- **Advantages:** Almost no **modal dispersion** (see below). Signals travel without spreading out. Can carry data over hundreds of kilometers without amplification.
- **Disadvantages:** Requires a very precise, expensive laser to couple light into such a narrow core. More expensive to manufacture and install.

### Multi-Mode Fiber

```
     Core diameter: ~50-62.5 micrometers
     
     ========================CLADDING========================
     ============CORE (~50 µm)============
     
     Ray 1:  ------>  (travels along axis)
     Ray 2:   \    /\    /\    /-->  (bounces at shallow angle)
               \/    \/    \/
     Ray 3:  \/\/\/\/\/\/\/\/-->  (bounces at steep angle — LONGER path!)
     
     ========================CLADDING========================
     
     Many modes (paths) of propagation
```

- **Core diameter:** ~50 to 62.5 micrometers
- **How light travels:** Multiple different bounce angles are possible. Different rays travel different total path lengths.
- **Modal dispersion:** Rays taking steeper bounce angles travel a longer zigzag path and arrive slightly later than rays traveling straight. A short pulse of light spreads out in time as it travels — this is **modal dispersion**. It limits the bandwidth and maximum distance.
- **Advantages:** Larger core means easier to couple light in (can use cheaper LEDs instead of lasers). Less precise alignment needed. Cheaper.
- **Disadvantages:** Modal dispersion limits distance. Typically used for short-distance connections (within buildings, data centers — up to ~550 meters at 10 Gbps).

### Graded-Index Multi-Mode Fiber

A clever solution to modal dispersion: instead of a sharp core-cladding boundary, the refractive index **gradually decreases** from the center outward. Rays traveling near the axis see higher n and travel slower; rays bouncing near the edges see lower n and travel faster. This compensates for the different path lengths, keeping pulse spread small. Graded-index multi-mode fiber can achieve ~2 km range at 10 Gbps.

---

## 7. Fiber Optic Applications

### Telecommunications

The most important application of optical fiber is carrying telephone calls, internet data, and television signals. Optical fibers have revolutionized global communication.

**Why fiber beats copper wire:**

| Property | Copper Wire | Optical Fiber |
|---|---|---|
| Signal carrier | Electrons | Photons (light pulses) |
| Speed | Limited by electrical capacitance | ~0.67c (2×10⁸ m/s) |
| Bandwidth | ~10 Gbps max (over short distances) | Terabits/s (multiple wavelengths) |
| Attenuation | ~10 dB/km (signal weakens fast) | ~0.2 dB/km (extremely low loss) |
| Maximum unamplified distance | ~1 km | ~80–100 km |
| Electromagnetic interference | Yes (susceptible to EMI) | None (light not affected by EMI) |
| Weight | Heavy | Very light |
| Diameter | Thicker | Thinner (~125 µm with cladding) |

**Transoceanic fiber optic cables** lie on the ocean floor, connecting continents. The first transatlantic fiber cable (TAT-8) was laid in 1988. Today's cables carry multiple terabits per second. For example, the MAREA cable between the US and Spain carries 160 terabits per second — enough for roughly 71 million simultaneous HD video streams.

### Medical Endoscopes

An **endoscope** is a thin, flexible tube containing thousands of optical fibers, used by doctors to look inside the human body without surgery.

```
     Endoscope fiber bundle (simplified):

     Doctor's eyepiece             Patient's body
     or camera                     (stomach, colon, etc.)
         |                                 |
     [viewer]----[fiber bundle]----[tiny lens & light]
     
     Illumination fibers: carry bright light INTO body
     Image fibers: carry image OUT from inside body
     
     Each fiber carries one "pixel" of the image
     Thousands of fibers → full image
```

- **Illumination fibers** carry light from an external source (typically a bright LED or xenon lamp) into the body cavity.
- **Image fibers** carry the reflected light image back out to a camera or eyepiece.
- The fiber bundle is flexible, so it can bend through the curved pathways of the digestive tract.
- Used in **colonoscopy** (colon inspection), **gastroscopy** (stomach), **bronchoscopy** (airways), **laparoscopy** (abdominal cavity), and many other procedures.
- Can pass tiny surgical instruments through a central channel, allowing minor surgery without large incisions (**keyhole surgery** or **minimally invasive surgery**).

### Decorative and Architectural Lighting

Fiber optic cables can pipe light from a single bright source to many distant endpoints. This is used in:
- **Starfield ceilings** in theaters and luxury cars
- **Decorative lamps** (fiber flowers, fountains of light)
- **Signage and displays** where electrical wiring would be hazardous (near water features, in bathrooms)
- Museum exhibit lighting (light can be kept away from UV-sensitive artifacts)

### Sensors

Changes to an optical fiber — slight bending, temperature variation, pressure — subtly affect how light travels through it. This makes fibers useful as extremely sensitive sensors:

- **Distributed temperature sensing:** measure temperature at thousands of points along a fiber simultaneously
- **Structural health monitoring:** optical fibers embedded in bridges or aircraft wings detect tiny strains and cracks
- **Hydrophones:** fiber cables on the ocean floor detect pressure waves (earthquakes, submarines)
- **Gyroscopes:** fiber optic gyroscopes (FOG) use the interference of light traveling in opposite directions around a coiled fiber to measure rotation extremely precisely — used in aircraft, submarines, and spacecraft navigation

---

## 8. How the Internet Uses Fiber Optics

### Encoding Data as Light Pulses

Information is encoded in binary: 0s and 1s. In a fiber optic cable:
- **1** = a pulse of laser light (light on)
- **0** = no pulse (light off)

Modern systems switch the laser on and off billions of times per second, encoding billions of bits per second on a single wavelength of light.

### Wavelength Division Multiplexing (WDM)

```
     One fiber, many colors:
     
     [Laser λ1 (1530nm)] ─┐
     [Laser λ2 (1550nm)] ─┤
     [Laser λ3 (1570nm)] ─┤──[Multiplexer]──[FIBER]──[Demultiplexer]──> receivers
     [Laser λ4 (1590nm)] ─┤
     [Laser λ5 (1610nm)] ─┘
     
     Each color carries an independent data stream
     All travel through the same fiber simultaneously
     At the receiving end, a prism/grating separates the colors again
```

**Wavelength Division Multiplexing (WDM)** allows many separate streams of data to share the same fiber, each using a different wavelength (color) of light. A modern **Dense WDM (DWDM)** system can carry 80 or more wavelengths in one fiber, multiplying the capacity 80-fold. With multiple fiber pairs in a cable, total capacities reach petabits per second.

### Optical Amplifiers

Light pulses gradually weaken (attenuate) as they travel through the fiber. Every ~80–100 km, the signal needs to be boosted.

Rather than converting the light signal to electricity, amplifying it electronically, and converting back (which is slow and expensive), modern systems use **Erbium-Doped Fiber Amplifiers (EDFA)**:

- A short segment of fiber is "doped" with erbium atoms
- A separate pump laser excites the erbium atoms to a higher energy state
- When the weakened signal passes through, the excited erbium atoms release energy as stimulated emission — amplifying the light signal directly
- This works for all WDM wavelengths simultaneously
- EDFA maintains signal quality over transoceanic distances

### Global Undersea Fiber Network

- Over **400 submarine cable systems** are currently in service or planned
- Total length exceeds **1.3 million km** — enough to wrap around Earth more than 32 times
- These cables carry approximately **99% of all international internet traffic** (satellite handles the remaining ~1%)
- A single modern cable segment can carry 10–200+ terabits per second
- Cables are buried under the seabed near coasts (to protect from anchors and fishing trawls) and rest on the ocean floor in deeper water
- The deepest undersea cables operate at depths of ~8,000 meters

---

## 9. Diamonds and Sparkle

Why does a diamond sparkle so brilliantly? The answer is almost entirely due to total internal reflection.

### The Physics of Diamond Brilliance

Recall from Example 3: diamond has n = 2.42, giving a critical angle of only **24.4°**.

This incredibly small critical angle means that light inside a diamond will undergo TIR at almost any angle it hits a facet. A masterfully cut diamond exploits this.

```
     Cross-section of a "brilliant cut" diamond:

                     [Table facet — top, flat]
                   /‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾\
                  /   Crown facets              \
                 /                               \
                /                                 \
               /─────────────────────────────────\
               |      Girdle (widest point)       |
               \                                 /
                \   Pavilion facets             /
                 \         ↗ ↖               /
                  \       ↗   ↖            /
                   \     ↗ TIR ↖         /
                    \   ↗       ↖       /
                     \ ↗    *    ↖     /
                      ↘     TIR  ↗
                       \   ↙↗  /
                        \ ↙  ↗/
                         V   /
                              
         Light enters top → bounces via TIR through pavilion facets
         → returns up and exits through the table or crown facets
         → appears as a bright flash from above
```

The process:
1. Light enters through the **table** (large flat top facet) or **crown facets** (angled facets above the girdle).
2. Inside the diamond, it hits a **pavilion facet** (angled facet below the girdle) at an angle > 24.4° → **TIR**.
3. It reflects to the opposite pavilion facet → **TIR** again.
4. Now traveling upward, it exits through the top of the diamond with great intensity.
5. The **depth** and **angle** of the pavilion determines where the light exits. A well-cut diamond returns light directly back to the viewer's eye.

A poorly cut diamond (too deep or too shallow pavilion) may lose light through the bottom or sides, appearing dark in the center — a defect called a **"fish-eye"** or **"nail-head."**

### Dispersion: The "Fire" of a Diamond

Beyond simple TIR, diamond has very high **dispersion** — its refractive index varies significantly with wavelength (color). Blue light has n ≈ 2.45; red light has n ≈ 2.41. This means different colors refract at slightly different angles, separating white light into its spectrum.

After multiple TIR bounces, the separated colors exit through different facets at different angles. This produces the flashes of spectral color — reds, blues, greens — visible as a diamond is moved. This is called the **"fire"** of a diamond, and it is distinct from its **"brilliance"** (overall brightness).

### Comparing Gemstones

| Gemstone | Refractive Index n | Critical Angle |
|---|---|---|
| Glass | 1.5 | 41.8° |
| Cubic Zirconia (CZ) | 2.15 | 27.8° |
| Moissanite | 2.65 | 22.2° |
| Diamond | 2.42 | 24.4° |
| Ruby/Sapphire | 1.77 | 34.4° |

Moissanite actually has a slightly higher n than diamond, giving it an even smaller critical angle. It also has higher dispersion (more "fire"), though some people find it "too flashy." Cubic zirconia is the most common diamond simulant but has lower dispersion and brilliance due to its smaller n.

---

## 10. Prisms and TIR

Prisms can use TIR to redirect light with perfect (100%) efficiency — no metallic mirror coating needed.

### Right-Angle Prism as a Mirror

A right-angle prism with a 45°–45°–90° geometry uses TIR on its hypotenuse face to reflect light through exactly 90°:

```
     Right-angle prism — side view:

         |\
         | \
         |  \
     →→→ |   \
     →→→ |    \  ← Hypotenuse face
     →→→ |     \    (light hits here at 45° from normal)
         |   TIR\   (45° > 41.8° critical angle for glass ✓)
         |_______\
             ↓↓↓
        Light exits downward
        
     45° angle of incidence > 41.8° critical angle → TIR occurs
     100% of light is reflected (vs ~92% for the best metallic mirrors)
```

A right-angle prism can also reflect light through 180° (retroreflect it straight back) if light enters the hypotenuse face at normal incidence.

### Porro Prisms in Binoculars

Binoculars face two problems: (1) a straight tube long enough for the required magnification would be extremely long and unwieldy; (2) a simple telescope produces an inverted image.

**Porro prisms** solve both problems using TIR:

```
     Binocular optical path with Porro prisms:

     Objective lens                        Eyepiece
          |                                    |
          |   PRISM 1 (folds path down)         |
          |   _____                             |
     →→→→|→ |     | → (TIR 90°)                 |
          |  |     ↓                             |
          |  |_____|                             |
          |     ↓ (TIR 90°)                      |
          |     → →→→ PRISM 2                   |
          |            _____                    |
          |           |     | → (TIR 90°)       |
          |           |     ↓                   |
          |           |_____|                   |
          |              ↓ (TIR 90°)            |
          |              →→→→→→→→→→→→→→→→→→→→→→|
     
     Net effect: optical path folded back and forth
     Image is right-side-up and left-right correct
     Binocular body is compact
     All reflections are TIR → virtually no light loss
```

Each Porro prism contains two TIR reflections. Two prisms per optical path give four reflections total, folding the optical path to fit within a short body and correcting the image orientation. Because TIR is used rather than metallic mirrors, the image quality is excellent with minimal light loss.

### Retroreflectors

Three right-angle prism faces arranged mutually perpendicular (a **corner cube reflector**) will reflect any incoming ray of light exactly back toward its source, regardless of the angle of incidence. This is used in:

- **Road studs** ("cats' eyes") — retroreflect car headlights back to the driver
- **Bicycle reflectors**
- **Surveying targets** for laser rangefinders
- **Lunar laser ranging** — corner cube reflectors placed on the Moon by Apollo astronauts are still used today to measure the Earth-Moon distance to millimeter precision

---

## 11. Other Applications of TIR

### Evanescent Waves and TIRF Microscopy

When TIR occurs, the light does not simply stop at the boundary. An electromagnetic wave called the **evanescent wave** penetrates a tiny distance into the less-dense medium — typically about the same order as the wavelength of light (~100–500 nm). The evanescent wave decays exponentially in intensity with distance from the boundary.

This tiny penetration depth is used in **Total Internal Reflection Fluorescence (TIRF) microscopy** — a technique used in biology to image molecules within a few hundred nanometers of a cell membrane, without illuminating the entire cell. This gives extremely sharp, low-background images of processes at the cell surface.

**Frustrated TIR** occurs when a second surface is brought within the evanescent wave range (~1 wavelength) of the first. Some light "tunnels" through the gap and appears in the second medium. This is an optical analog of quantum tunneling.

### Optical Waveguides on Chips

The same TIR principle that guides light through optical fibers is used to guide light through tiny channels (waveguides) on silicon chips — **photonic integrated circuits (PICs)**. Light is guided through narrow ridges of silicon (n ≈ 3.5) surrounded by silicon dioxide (n ≈ 1.45), with TIR at the silicon-SiO₂ boundary.

PICs are used in:
- **Silicon photonics** for high-speed data center interconnects
- **Optical sensors** (biosensors, chemical sensors)
- **LIDAR** (Light Detection and Ranging) for autonomous vehicles
- Quantum computing experiments (photonic qubits)

The ability to integrate multiple optical functions (lasers, modulators, filters, detectors) on a single chip using waveguides is creating a revolution in computing and sensing analogous to what the transistor did for electronics.

### Optical Fiber Sensors

Beyond simple data transmission, the sensitivity of optical fibers to their environment makes them superb sensors:

- **Distributed Acoustic Sensing (DAS):** vibrations along a fiber cause tiny Doppler shifts and backscatter changes, allowing the entire fiber length to act as a microphone array. Can detect earthquakes, traffic, trains, and even footsteps from a single interrogation unit connected at one end.

- **Oil and gas pipeline monitoring:** fibers attached to pipelines detect leaks, pressure changes, and corrosion.

- **Railway track monitoring:** fibers along rail lines detect train position, speed, wheel defects, and track deformation in real-time.

---

## 12. Summary

- When light travels from a **denser medium** (higher n) to a **less dense medium** (lower n), the refracted ray bends away from the normal (θ2 > θ1)

- The **critical angle** (θ_c) is the angle of incidence at which the refracted ray travels exactly along the boundary (θ2 = 90°). It is given by: sin(θ_c) = n2 / n1

- The critical angle only exists when n1 > n2. The higher n1 relative to n2, the smaller the critical angle

- When θ1 > θ_c, **Total Internal Reflection** (TIR) occurs: 100% of light is reflected back into the denser medium. No light escapes into the less dense medium

- TIR reflects 100% of light — more efficient than any metallic mirror

- Critical angles for common situations:
  - Glass (n=1.5) to air: 41.8°
  - Water (n=1.33) to air: 48.8°
  - Diamond (n=2.42) to air: 24.4°

- **Optical fibers** use TIR to guide light along their length. The **core** (high n) is surrounded by **cladding** (slightly lower n). Light bouncing between the core-cladding boundaries undergoes TIR and cannot escape

- **Single-mode fibers** (~9 µm core) carry one propagation mode, have no modal dispersion, and are used for long-distance telecommunications

- **Multi-mode fibers** (~50 µm core) carry many modes, have modal dispersion (limiting distance and bandwidth), but are cheaper and easier to connect

- Fiber optic cables carry ~99% of international internet data. One fiber can carry terabits per second via **Wavelength Division Multiplexing (WDM)**

- **Medical endoscopes** use bundles of thousands of fibers to illuminate and image inside the body without surgery

- **Diamonds** sparkle because their very small critical angle (24.4°) causes light to bounce many times inside by TIR before exiting through the top, producing brilliant brightness. High dispersion produces colored "fire"

- **Right-angle prisms** use TIR to reflect light at 90° or 180° with 100% efficiency. **Porro prisms** in binoculars use TIR to fold the optical path and correct image orientation

- **Evanescent waves** (exponentially decaying waves that briefly penetrate the less-dense medium during TIR) are used in TIRF microscopy and fiber coupling

- **Photonic integrated circuits** use TIR in waveguides on silicon chips, enabling optical computing and sensing

---

## 13. Key Equations

**Snell's Law (general):**

    n1 × sin(θ1) = n2 × sin(θ2)

**Critical angle (general, requires n1 > n2):**

    sin(θ_c) = n2 / n1

    θ_c = arcsin(n2 / n1)

**Critical angle when medium 2 is air (n2 = 1):**

    sin(θ_c) = 1 / n1

    θ_c = arcsin(1 / n1)

**Condition for Total Internal Reflection:**

    θ1 > θ_c   AND   n1 > n2

**Critical angles for common materials (to air):**

    Glass (n=1.50):     θ_c = 41.8°
    Water (n=1.33):     θ_c = 48.8°
    Crown glass (n=1.52): θ_c = 41.1°
    Diamond (n=2.42):   θ_c = 24.4°
    Cubic zirconia (n=2.15): θ_c = 27.8°

**Numerical Aperture of an optical fiber:**

    NA = n_core × sin(θ_acceptance)

    Or equivalently:
    NA = sqrt(n_core² - n_cladding²)

**Evanescent wave penetration depth (d):**

    d ≈ λ / (4π × sqrt(n1² × sin²(θ1) - n2²))

    Where λ is the wavelength of light in vacuum and θ1 > θ_c

---

*Next Chapter — Chapter 65: Lenses and Image Formation — explores how curved glass surfaces bend light to form images, and how the eye, cameras, telescopes, and microscopes work.*
