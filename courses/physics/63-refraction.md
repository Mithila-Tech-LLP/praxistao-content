# Chapter 63: Refraction

> **"A straw looks bent in water. A swimming pool looks shallower than it is. A prism splits white light into a rainbow. All of this is refraction — light bending as it changes speed."**

---

## Table of Contents

1. [What is Refraction?](#1-what-is-refraction)
2. [Refractive Index](#2-refractive-index)
3. [Snell's Law](#3-snells-law)
4. [Worked Examples](#4-worked-examples)
5. [Everyday Refraction Phenomena](#5-everyday-refraction-phenomena)
6. [Atmospheric Refraction](#6-atmospheric-refraction)
7. [Dispersion: Why Prisms Split White Light](#7-dispersion-why-prisms-split-white-light)
8. [Rainbows](#8-rainbows)
9. [Refraction and Light Pipes (Preview of Fiber Optics)](#9-refraction-and-light-pipes-preview-of-fiber-optics)
10. [Summary](#summary)
11. [Key Equations](#key-equations)

---

## 1. What is Refraction?

Light travels at different speeds in different materials. In a vacuum, light zips along at about 3 × 10⁸ m/s — the ultimate speed limit of the universe. But when light enters glass, water, or any transparent medium, it slows down. This speed change does something fascinating: it causes the light to **bend**. That bending is called **refraction**.

### Why Does Slowing Down Cause Bending?

Think about a car driving on a smooth road that hits a patch of soft mud on one side. The wheel in the mud slows down first while the other wheel keeps its original speed. Because one side of the car is slower than the other, the car turns — it pivots toward the slower side. The car doesn't drive straight through; it changes direction.

Light behaves exactly the same way when it crosses from one medium into another at an angle.

Here's another helpful picture. Imagine a marching band in neat parallel rows walking across a field. They reach a line in the ground where the grass on one side is perfectly mown and the other side is thick, uncut grass. As each column of people crosses the line, they slow down. The people who hit the line first slow down while those who haven't reached it yet are still moving fast. This causes the entire column to swing around — the rows bend as they cross the boundary.

Light waves work the same way. The wavefronts (the rows in the marching-band picture) bend when they cross a boundary between two materials.

### The One Exception

If light hits the boundary at exactly 90° — meaning it's traveling straight into the surface with zero angle — then both sides of the wavefront hit the boundary at exactly the same time. There's no "one side slowing down first." The light passes straight through without bending. This is called **normal incidence**, and at normal incidence, no refraction (direction change) occurs, even though the light still changes speed.

### Visualizing Wavefronts at a Boundary

The diagram below shows light waves (wavefronts) crossing from a fast medium (like air) into a slow medium (like glass). Notice how the wavelength gets shorter in the denser medium because the wave is moving slower but the frequency stays the same.

```
     AIR (fast, low n)           GLASS (slow, high n)
     ------------------          ---------------------

     |      |      |              | | | |
     |      |      |      -->     | | | |
     |      |      |              | | | |
     |      |      |              | | | |
 wavefronts spread out       wavefronts bunched up
 (long wavelength)           (short wavelength)

     Direction of travel -->
```

When the wavefronts hit the boundary at an angle:

```
     AIR                         GLASS
                  /
     ----  ----  /  -- -- -- --
               /
    ----  ----/  -- -- -- --
            /|
      ----  / |  -- -- -- --
          /  |
       ---/  |   -- -- -- --
         /   |
  Wavefronts bend at the interface
  because part of each wave has
  already slowed down
```

The key insight: **refraction is caused by a change in wave speed at a boundary**.

---

## 2. Refractive Index

The **refractive index** (also called the **index of refraction**) of a material tells us how much slower light travels in that material compared to its speed in a vacuum.

**Formula:**

```
n = c / v
```

Where:
- **n** = refractive index (dimensionless — no units)
- **c** = speed of light in vacuum = 3 × 10⁸ m/s
- **v** = speed of light in the material (m/s)

Since light is always slower in a material than in a vacuum (nothing can make light go *faster* than its vacuum speed), we always have **n ≥ 1**. Vacuum has n = 1 exactly, by definition.

### Table of Common Refractive Indices

| Material         | Refractive Index (n) |
|------------------|----------------------|
| Vacuum           | 1.000 (exact)        |
| Air              | 1.0003               |
| Ice              | 1.31                 |
| Water            | 1.33                 |
| Crown glass      | 1.52                 |
| Dense flint glass| 1.62                 |
| Diamond          | 2.42                 |

Notice:
- Air is so close to vacuum that we almost always approximate n_air = 1.
- Diamond has the highest refractive index in the table, which is why diamonds sparkle so brilliantly — light bends a lot when it enters and exits, and much of it is trapped inside (total internal reflection, covered in the next chapter).

### What High n Means

A higher refractive index means:
- Light travels **slower** in that material.
- Light **bends more** when it enters that material from air.

We can calculate the actual speed of light in any material:

```
v = c / n
```

**Example: Speed of light in water**

```
v = (3 × 10⁸ m/s) / 1.33
v = 2.26 × 10⁸ m/s
```

Light in water travels at about 226 million meters per second — fast, but about 25% slower than in vacuum.

**Example: Speed of light in crown glass**

```
v = (3 × 10⁸ m/s) / 1.52
v = 1.97 × 10⁸ m/s
≈ 2 × 10⁸ m/s
```

In crown glass, light moves at roughly two-thirds of its vacuum speed.

**Example: Speed of light in diamond**

```
v = (3 × 10⁸ m/s) / 2.42
v = 1.24 × 10⁸ m/s
```

Light inside a diamond is only moving at about 40% of its vacuum speed!

---

## 3. Snell's Law

**Snell's Law** is the quantitative rule that tells us exactly how much light bends when it crosses a boundary. It was discovered experimentally by Willebrord Snellius around 1621 and is the foundation of all optics calculations.

**Snell's Law:**

```
n1 × sin(θ1) = n2 × sin(θ2)
```

Where:
- **n1** = refractive index of the first medium (where light is coming from)
- **n2** = refractive index of the second medium (where light is going into)
- **θ1** = angle of incidence (angle between the incoming ray and the **normal**)
- **θ2** = angle of refraction (angle between the refracted ray and the **normal**)

**IMPORTANT:** All angles are measured from the **normal** — an imaginary line perpendicular (90°) to the surface at the point where the ray hits. Do NOT measure angles from the surface itself.

### ASCII Diagram: Snell's Law Setup

```
             MEDIUM 1 (n1)
             
                  |  ← Normal (perpendicular to surface)
    Incident ray  |
         \        |
          \  θ1  |
           \     |
            \    |
             \   |
==============\==|================================  ← Boundary
               \ |
                \|  θ2
                 |\
                 | \
                 |  \  ← Refracted ray
                 |   \
             MEDIUM 2 (n2)
             
If n2 > n1: θ2 < θ1 (bends TOWARD normal)
If n2 < n1: θ2 > θ1 (bends AWAY from normal)
```

### Bending Toward or Away from the Normal?

This is one of the most important things to remember in refraction:

**Going into a denser medium (n2 > n1):**
Light slows down → bends TOWARD the normal → θ2 is smaller than θ1.

```
           AIR (n=1)               GLASS (n=1.5)

               |                       |
  Ray comes in |                       |
        \      |                       |
         \  30°|                       |
          \    |                       |
           \   |                       |
============\==|========================  Surface
             \ |
              \| 19.5°   ← Smaller angle!
               |\
               | \
               |  \
               |   Ray continues in glass (closer to normal)
```

**Going into a less dense medium (n2 < n1):**
Light speeds up → bends AWAY from the normal → θ2 is larger than θ1.

```
           GLASS (n=1.5)              AIR (n=1)

               |                       |
  Ray comes in |                       |
        \      |                       |
         \  20°|                       |
          \    |                       |
           \   |                       |
============\==|========================  Surface
             \ |
              \| 27°  ← Larger angle!
               |\
               | \
               |  \
               |   Ray in air (farther from normal)
```

### Special Case: Normal Incidence

If θ1 = 0° (ray hits boundary head-on, perpendicular to surface):

```
n1 × sin(0°) = n2 × sin(θ2)
n1 × 0 = n2 × sin(θ2)
0 = n2 × sin(θ2)
sin(θ2) = 0
θ2 = 0°
```

No bending. The ray goes straight through. This is consistent with our intuition — if light hits a surface dead-on, it doesn't need to bend.

### Where Does Snell's Law Come From?

Here's the conceptual explanation (you don't need the full derivation for this chapter). When a wave crosses a boundary, the frequency of the wave cannot change — it's determined by the source and stays constant. Since v = f × λ (wave speed = frequency × wavelength), if the speed changes but frequency stays the same, the wavelength must change.

At the boundary, the part of the wavefront along the surface must match on both sides (you can't have a "gap" or "tear" in the wave). This matching condition, when you work it through using Huygens' principle (every point on a wavefront acts as a new wave source), gives exactly Snell's Law.

---

## 4. Worked Examples

Let's practice using Snell's Law with a series of worked examples. Each one builds a different skill.

---

### Example 1: Air to Glass (Finding the Refraction Angle)

**Problem:** A ray of light in air (n = 1.00) hits a piece of glass (n = 1.50) at an angle of incidence of 30°. What is the angle of refraction?

**Given:**
```
n1 = 1.00  (air)
n2 = 1.50  (glass)
θ1 = 30°
θ2 = ?
```

**Apply Snell's Law:**
```
n1 × sin(θ1) = n2 × sin(θ2)
1.00 × sin(30°) = 1.50 × sin(θ2)
1.00 × 0.500 = 1.50 × sin(θ2)
0.500 = 1.50 × sin(θ2)
sin(θ2) = 0.500 / 1.50
sin(θ2) = 0.333
θ2 = arcsin(0.333)
θ2 = 19.5°
```

**Answer:** The refracted ray makes an angle of **19.5°** with the normal inside the glass.

**Check:** θ2 (19.5°) < θ1 (30°), so the light bent toward the normal. We're going from air (less dense, n=1) into glass (more dense, n=1.5), so this is correct. ✓

---

### Example 2: Water to Air (Bending Away from Normal)

**Problem:** A ray of light travels through water (n = 1.33) and hits the water-air boundary at an angle of 20° to the normal. Find the angle of refraction in air.

**Given:**
```
n1 = 1.33  (water)
n2 = 1.00  (air)
θ1 = 20°
θ2 = ?
```

**Apply Snell's Law:**
```
n1 × sin(θ1) = n2 × sin(θ2)
1.33 × sin(20°) = 1.00 × sin(θ2)
1.33 × 0.342 = 1.00 × sin(θ2)
0.455 = sin(θ2)
θ2 = arcsin(0.455)
θ2 = 27.1°
```

**Answer:** The refracted ray in air makes an angle of **27.1°** with the normal.

**Check:** θ2 (27.1°) > θ1 (20°), so light bent away from the normal. We're going from water (denser) to air (less dense), which is correct. ✓

---

### Example 3: Finding the Refractive Index of an Unknown Material

**Problem:** A physicist is testing an unknown liquid. A ray of light in air hits the liquid surface at 45° to the normal, and the refracted ray inside the liquid makes an angle of 28° with the normal. What is the refractive index of the liquid?

**Given:**
```
n1 = 1.00  (air)
θ1 = 45°
θ2 = 28°
n2 = ?
```

**Apply Snell's Law and solve for n2:**
```
n1 × sin(θ1) = n2 × sin(θ2)
1.00 × sin(45°) = n2 × sin(28°)
1.00 × 0.707 = n2 × 0.469
0.707 = n2 × 0.469
n2 = 0.707 / 0.469
n2 = 1.51
```

**Answer:** The refractive index of the unknown liquid is approximately **1.51**.

**Interpretation:** This value is very close to the refractive index of crown glass (1.52). The unknown liquid is probably a type of glass-like oil or glycerin-type substance.

---

### Example 4: A Ray Passing Through Multiple Interfaces

**Problem:** A ray of light in air hits a glass slab (n = 1.50) at 40°, passes through the glass, and exits back into air on the other side. What is the exit angle?

**Step 1: Air → Glass (entering the slab)**
```
n1 = 1.00, θ1 = 40°, n2 = 1.50, θ2 = ?
1.00 × sin(40°) = 1.50 × sin(θ2)
sin(θ2) = sin(40°) / 1.50 = 0.643 / 1.50 = 0.428
θ2 = arcsin(0.428) = 25.4°
```

**Step 2: Glass → Air (exiting the slab, other side is parallel)**
```
n1 = 1.50, θ1 = 25.4°, n2 = 1.00, θ3 = ?
1.50 × sin(25.4°) = 1.00 × sin(θ3)
sin(θ3) = 1.50 × 0.429 = 0.643
θ3 = arcsin(0.643) = 40°
```

**Answer:** The exit angle is **40°** — exactly the same as the entry angle!

**Key insight:** When light passes through a flat slab with parallel sides, it exits at the same angle it entered. The light is shifted sideways (displaced), but not angled differently. This is why looking through a flat window doesn't distort the direction things appear to be.

---

### Example 5: Calculating Apparent Depth

**Problem:** A coin sits at the bottom of a fish tank. The water is 60 cm deep (n_water = 1.33). How deep does the coin appear to be to someone looking straight down from above?

**Formula for apparent depth (derivation uses Snell's Law for small angles):**
```
Apparent depth = Real depth / n
Apparent depth = 60 cm / 1.33
Apparent depth = 45.1 cm
```

**Answer:** The coin appears to be only about **45 cm** deep, even though it's actually 60 cm down. The water makes it look closer than it really is.

---

## 5. Everyday Refraction Phenomena

Refraction isn't just a classroom concept — it explains dozens of things you see every day.

### The Bent Straw Illusion

When you put a straw or pencil in a glass of water, it looks broken or bent at the water surface. Why?

Light from the part of the straw underwater travels up to the surface and refracts (bends away from the normal) as it goes from water into air. Your eye traces this bent ray in a straight line back to where it *appears* to have come from — a position that's shifted from the straw's actual location.

```
    Your Eye
       ^
       |
       |   ← You trace the ray straight back
       |       (ignoring the bend)
       |
  ~~~~~|~~~~~~~~~~~~~~~~~~~~~~~~~~~~  Water Surface
       |  ↗ Refracted ray bends here
       | /
       |/
       *  ← Where the straw actually is
       
The straw appears to be at a different position than it really is.
This is why straws look bent at the water surface!
```

### Apparent Depth: Pools Look Shallower Than They Are

A swimming pool marked "3 m deep" looks much shallower from the surface. This is because light from the bottom refracts at the water surface when leaving, making the bottom appear to be closer than it actually is.

**Formula:**
```
Apparent Depth = Real Depth / n
```

**Example:** A pool is 3.0 m deep. n_water = 1.33.
```
Apparent Depth = 3.0 / 1.33 = 2.26 m
```

The pool appears to be only **2.26 m** deep — nearly a metre shallower than it really is! This is why you should never jump into a pool based on how deep it looks from the edge.

**Example:** A fish swims at a real depth of 4.0 m. An observer standing above (in air) sees the fish at:
```
Apparent Depth = 4.0 / 1.33 = 3.0 m
```

The fish appears to be only **3 m** down, even though it's actually 4 m down.

### Snell's Window: The Fish's View

Now consider the reverse — what does a fish see looking up toward the surface?

A fish looking upward can see the entire above-water world — the sky, birds, people — but only within a circular "window" of cone half-angle about 48.8°. This is the **critical angle** for water. Outside this window, the fish sees only the reflected image of the underwater world (total internal reflection, covered fully in the next chapter). From the fish's perspective, the entire 360° view above the water is compressed into a cone less than 100° wide. This cone is called **Snell's Window**.

### Refraction When Wading

When you wade into water, your legs appear shorter and bent at the surface. Light from your legs underwater refracts as it exits, making your submerged legs appear displaced. The effect makes it look like your legs end right at the water surface and shorter legs begin there — the classic "broken" look.

---

## 6. Atmospheric Refraction

Earth's atmosphere is thickest and densest near the ground, getting thinner and less dense as you go up. This means the refractive index of air changes gradually with altitude — it's higher near the ground and drops toward 1.000 in space. Light passing through this gradient of refractive index bends continuously, not just at one sharp boundary.

### Stars Appear Higher Than They Are

When starlight enters the atmosphere at a shallow angle, it bends — curving downward toward Earth. Your eye traces this curved ray back in a straight line, so the star appears to be in a slightly higher position than its true geometrical location. Near the horizon, this effect can raise the apparent position of a star by nearly 0.5° (about the width of the full Moon).

### Twinkling Stars (Scintillation)

Stars are so far away that they appear as pinpoint light sources. As starlight passes through the turbulent atmosphere, tiny pockets of warm and cool air with slightly different refractive indices constantly shift around. This causes the star's apparent position and brightness to jitter rapidly — we call this **scintillation**, or twinkling.

Planets, on the other hand, are much closer to us and appear as tiny disks (not pinpoints) in the sky. The twinkling effects average out across the disk, so planets generally appear steady while stars twinkle.

### The Sun and Moon at the Horizon

The Sun and Moon actually rise and set a few minutes earlier than they geometrically should, because atmospheric refraction "lifts" their image above the geometric horizon. In fact, when you see the Sun just sitting on the horizon, the actual disk of the Sun is already about 0.5° below the horizon geometrically. Refraction makes it visible a couple of minutes earlier than pure geometry would predict.

### Mirages

On a hot day, the air immediately above a road or desert surface is very hot and less dense (lower refractive index) than the slightly cooler air a few centimetres higher. This creates an upward gradient of refractive index near the ground.

Light from the sky traveling downward at a shallow angle bends gradually upward as it encounters progressively lower-n air near the ground. To your eye, it appears to come from the ground ahead — looking exactly like water reflecting the sky.

```
         Cool air (higher n)
     ___________________________________________
         Warmer air
     ___________________________________________
         Hot air (lower n)
     === Road surface =========================
     
  Light from     Light curves     You see
  the sky    →   upward       →   "water" ahead
  
         _____
        /     \          Observer
       |  Sky  |             👁
        \_____/              |
           |                 |
           |   ______________/
           |  /  (curved path)
    -------|-/------------------------------------------  Ground
    
The curved path makes it look like sky is reflected from
the ground, like a puddle of water (mirage).
```

This is the **inferior mirage** — the most common type, seen on hot roads and in deserts. It's not a hallucination; it's real light that has just taken a very unusual path to reach your eyes.

---

## 7. Dispersion: Why Prisms Split White Light

You might expect glass to have a single, fixed refractive index. But it turns out the refractive index of any transparent material depends slightly on the **wavelength** (color) of the light passing through it. This dependence is called **dispersion**.

### How Refractive Index Varies with Color

In glass (and most transparent materials), shorter wavelengths (toward violet and blue) are slowed down slightly more than longer wavelengths (toward orange and red). In other words:

- **Violet light** (λ ≈ 400 nm): higher n, more refraction, bends more
- **Red light** (λ ≈ 700 nm): lower n, less refraction, bends less

For crown glass, the refractive index ranges from about 1.523 for red light to about 1.531 for violet light. The difference is small (about 0.5%) but visible when light travels through a prism.

### How a Prism Splits White Light

White light is a mixture of all visible wavelengths. When it enters a glass prism at an angle:

1. Each color refracts by a slightly different amount at the first surface (violet bends most, red bends least).
2. The colors spread out inside the prism.
3. At the second surface, they refract again, spreading further.
4. A full rainbow spectrum emerges from the prism.

```
                          /|
                         / |
    White light →       /  |       ← Violet (bends most)
    ──────────────►    /   |       ← Blue
                      /    |       ← Green
                     /     |       ← Yellow
                    /      |       ← Orange
                   /       |       ← Red (bends least)
                  /        |
                 /  PRISM  |
                /__________|
                
     Entry face            Exit face
     (first refraction)    (second refraction, further spreading)
```

The order from top to bottom exiting the prism: **Violet, Indigo, Blue, Green, Yellow, Orange, Red** (VIBGYOR) — or reading it as a mnemonic in the other direction: **ROY G BIV** (Red, Orange, Yellow, Green, Blue, Indigo, Violet).

### Newton's Great Experiment

In 1666, Isaac Newton performed a famous experiment. He passed a beam of sunlight through a prism and got a spectrum (not surprising — people had done this before). But then he did something clever: he let only one color from the first prism hit a second prism. The second prism refracted that color but did NOT spread it any further. The color came out as a single color, not another spectrum.

This proved that **white light is composed of all the colors** and the prism is just separating what is already there — not adding color or changing the light. Newton also recombined the spectrum back into white light using a second prism (or a lens), conclusively showing that dispersion is reversible.

### Chromatic Aberration in Lenses

Because different colors refract by different amounts, a simple lens focuses different colors at different distances. Red light focuses slightly farther from the lens than blue light. This means a simple lens cannot bring all colors to a precise single focal point simultaneously.

This defect is called **chromatic aberration**. You can see it as colored fringes (usually purple/blue on one side and red/orange on the other side) around high-contrast images in cheap cameras and telescopes.

The fix is an **achromatic doublet**: two lenses made from different types of glass (with different dispersion properties) cemented together. The dispersions cancel each other out while still converging the light, giving a sharp image.

---

## 8. Rainbows

A rainbow is one of nature's most spectacular demonstrations of refraction and dispersion. It happens when sunlight passes through millions of spherical water droplets suspended in the air (usually rain or mist).

### The Path of Light Through a Droplet

Each individual water droplet acts like a tiny combination of prism and mirror. Here's what happens:

**Step 1 — Refraction entering the droplet:**
Sunlight (white light) hits the spherical surface of the droplet. Different colors refract by slightly different amounts (dispersion), just as in a prism. The light bends and slows as it enters the water (n_water ≈ 1.33).

**Step 2 — Internal reflection:**
The refracted light travels through the droplet and hits the back curved surface. At the correct angles, it undergoes an **internal reflection** (bounces back inside the droplet, like a mirror).

**Step 3 — Refraction exiting the droplet:**
The reflected light travels back to the front surface of the droplet and refracts again as it exits into air. This second refraction further separates the colors.

```
                     Sunlight in (white)
                           ↓
                     ______↓______
                    /      │      \
                   /    1. Refraction
                  /  (white → colors)
                 |                  |
                 |      2. Internal |
     Red exits  ←|←── reflection    |← Violet exits
     at ~42°     |   (back of drop) |    at ~40°
                  \                /
                   \_____________ /
                   
         Red exits at a slightly wider angle (42°)
         Violet exits at a slightly narrower angle (40°)
```

### Why We See an Arc

The key observation: **red light exits each droplet at about 42° from the original direction of sunlight, and violet light exits at about 40°.**

Now imagine standing with your back to the Sun, looking at a sky full of water droplets. Only the droplets that happen to be at exactly 42° from your anti-solar point (the point directly opposite the Sun from your perspective) will send their red light to your eye. Only those at 40° will send violet.

Droplets at the same angle from the antisolar point form a circle. So you see:
- A circular arc of red at 42°
- A circular arc of violet at 40°
- All other colors in between

This gives the classic arc shape of a rainbow, with **red on the outside and violet on the inside** of the primary rainbow.

```
        Sun behind you →   👤  ← You (observer)
                           |
              Anti-solar   |
              point        ↓
                           ×
              
                         *  *  *
                      *           *
                 42° /               \ 42°
                    *    RED ARC      *
                   40°               40°
                    *   VIOLET ARC   *
                      *           *
                         *  *  *
```

**Why you can't walk to the end of a rainbow:** The rainbow's position is always at the same angle relative to *your* eyes and the Sun. As you move, the rainbow moves with you. It has no fixed location in space — it's an angular phenomenon, not a physical object.

### Primary vs. Secondary Rainbow

**Primary rainbow:** One internal reflection in each droplet. Colors go red (top/outer, 42°) to violet (bottom/inner, 40°).

**Secondary rainbow:** Two internal reflections. It appears at about 50–53° from the antisolar point (higher in the sky). Because there's an extra reflection, the color order is **reversed**: violet is on the outside and red is on the inside of the secondary rainbow.

The secondary rainbow is fainter because:
1. Each reflection loses some light (some escapes the droplet rather than reflecting).
2. Two reflections means more light is lost than one.

Between the two rainbows, the sky appears noticeably darker — this region is called **Alexander's Dark Band**, and it occurs because no droplets in that angular range can direct light to your eye.

### The Sun Must Be Behind You

A rainbow always appears in the part of the sky opposite the Sun. This is because light enters the droplets on the Sun's side, bounces back, and exits toward you — so you always need to be facing away from the Sun to see a rainbow.

---

## 9. Refraction and Light Pipes (Preview of Fiber Optics)

We've seen that when light goes from a denser medium (like glass) to a less dense medium (like air), it bends away from the normal — the refracted angle θ2 is larger than the incident angle θ1.

Now here's a fascinating question: what if we keep increasing the angle of incidence? As θ1 increases, so does θ2. At some critical value of θ1, the refracted angle θ2 reaches 90° — meaning the refracted ray would travel along the boundary surface itself, not into the second medium at all.

If θ1 exceeds this critical angle, there is **no refracted ray** at all. All of the light reflects back inside the original medium. This is called **total internal reflection (TIR)**.

```
                Glass (n=1.5)         Air (n=1.0)
                
          θ1 small                         |
    Ray →   \      →  refracted ray goes into air
             \         (normal refraction)
--------------\-------------------------------  Boundary

          θ1 = critical angle              |
    Ray →    \    → refracted ray travels  |
              \         along boundary     |
--------------\----------|----------------  Boundary

          θ1 > critical angle
    Ray →    \  ← All light reflects back
              \ ← into glass (TIR!)
--------------\-------------------------------  Boundary
               \
                \
```

For glass with n = 1.5, the critical angle is about 41.8°. For water (n = 1.33), it's about 48.8°.

### Why This Matters: Optical Fibers

This principle is the basis for **fiber optic cables**, which carry internet, telephone, and television signals around the world. A glass fiber (with a high refractive index core surrounded by a lower refractive index cladding) allows light to travel along its length — even around curves — by bouncing off the walls repeatedly via total internal reflection.

```
       Light in
          →→→→
   _______________________
  |   Core glass (n=1.5)  |  → Light bounces inside
  |___/\/\/\/\/\/\/\/\____|      via total internal
  |  Cladding (n=1.45)   |      reflection, losing
  |_______________________|      very little energy
                           → Light out
                           
The core-cladding boundary is always hit at an angle
greater than the critical angle → light stays inside!
```

This preview topic — total internal reflection — is the subject of the next chapter.

---

## Summary

- **Refraction** is the bending of light (or any wave) when it crosses a boundary between two media with different speeds.

- Refraction is caused by a **change in wave speed** at the boundary. If the angle of incidence is 0° (normal incidence), no bending occurs.

- The **refractive index** n = c/v tells us how much slower light travels in a medium compared to vacuum. n ≥ 1 always.

- Higher refractive index means slower light, shorter wavelength (same frequency), and more bending when entering from a less dense medium.

- **Snell's Law** — n1 × sin(θ1) = n2 × sin(θ2) — gives the exact relationship between angles of incidence and refraction at any boundary.

- Going from a less dense to a more dense medium (n2 > n1): light bends **toward** the normal (θ2 < θ1).

- Going from a more dense to a less dense medium (n2 < n1): light bends **away** from the normal (θ2 > θ1).

- **Apparent depth** = Real depth / n. Objects underwater appear shallower than they really are.

- A flat-sided glass slab causes lateral displacement of a light ray but does not change the ray's direction overall — the exit angle equals the entry angle.

- **Atmospheric refraction** makes stars appear slightly higher than they really are, causes twinkling (scintillation) due to atmospheric turbulence, and creates mirages by bending light near hot surfaces.

- **Dispersion** is the variation of refractive index with wavelength. Violet light has higher n (bends more) and red light has lower n (bends less) in glass.

- A **prism** splits white light into a spectrum because of dispersion. Newton proved this by recombining the spectrum back into white light.

- **Chromatic aberration** in lenses results from dispersion; fixed with an achromatic doublet.

- **Rainbows** are formed by refraction (into droplet), internal reflection (at back of droplet), and refraction again (exiting droplet). Red exits at ~42°, violet at ~40°. Red is on the outside of the primary rainbow.

- The **secondary rainbow** (fainter, higher in sky) results from two internal reflections and has its colors reversed.

- **Total internal reflection** occurs when light in a dense medium hits a less dense medium at an angle exceeding the critical angle. This is the basis of fiber optic cables.

---

## Key Equations

**Refractive Index:**
```
n = c / v

where:
  n = refractive index (dimensionless)
  c = 3 × 10⁸ m/s (speed of light in vacuum)
  v = speed of light in the medium (m/s)
```

**Speed of Light in a Medium:**
```
v = c / n
```

**Snell's Law:**
```
n1 × sin(θ1) = n2 × sin(θ2)

where:
  n1 = refractive index of incident medium
  n2 = refractive index of refracted medium
  θ1 = angle of incidence (from normal)
  θ2 = angle of refraction (from normal)
```

**Apparent Depth:**
```
Apparent Depth = Real Depth / n

(valid for viewing from air, looking straight down)
```

**Critical Angle for Total Internal Reflection:**
```
sin(θc) = n2 / n1

where n1 > n2 (light going from dense to less dense)
For glass-air: sin(θc) = 1/1.5 = 0.667, so θc ≈ 41.8°
For water-air: sin(θc) = 1/1.33 = 0.752, so θc ≈ 48.8°
```

**Wavelength in a Medium:**
```
λ_medium = λ_vacuum / n

(frequency stays the same; wavelength decreases in denser media)
```

---

*Next Chapter: Total Internal Reflection and Fiber Optics — where we explore what happens when light tries to escape a dense medium at too steep an angle, and how this traps light inside glass fibers to carry information around the world.*
