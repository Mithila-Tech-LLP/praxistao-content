# Chapter 65: Lenses

> **"A lens is a curve of glass that bends light in a controlled way. The lens in your eye focuses light onto your retina. The lens in a camera captures an image. The lens in a telescope brings distant stars close. Understanding lenses means understanding how we see."**

---

## Table of Contents

- [65.1 How Lenses Work: Refraction](#651-how-lenses-work-refraction)
- [65.2 Types of Lenses](#652-types-of-lenses)
- [65.3 Key Terms and Definitions](#653-key-terms-and-definitions)
- [65.4 Ray Diagrams for Convex Lenses](#654-ray-diagrams-for-convex-lenses)
- [65.5 Ray Diagrams for Concave Lenses](#655-ray-diagrams-for-concave-lenses)
- [65.6 The Lens Equation](#656-the-lens-equation)
- [65.7 Magnification](#657-magnification)
- [65.8 The Power of a Lens](#658-the-power-of-a-lens)
- [65.9 Lens Defects (Aberrations)](#659-lens-defects-aberrations)
- [65.10 The Human Eye](#6510-the-human-eye)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 65.1 How Lenses Work: Refraction

Lenses work because of **refraction** — light bends when it passes from one medium to another at an angle, following Snell's Law: n₁ sin(θ₁) = n₂ sin(θ₂).

```
LIGHT ENTERING GLASS:
  
  Air (n=1)       Glass (n=1.5)
  
    \              |
     \  θ₁=40°   |
      \            |
       +-----------+  (interface)
        \          |
         \ θ₂=25° |
          \        |
  
  sin(40°)/sin(25°) = n_glass/n_air
  
  Light bends TOWARD the normal when entering a denser medium.
  Light bends AWAY from the normal when leaving a denser medium.
```

A lens is curved glass. The curved surface means each part of the lens bends light by a different amount, converging or diverging the rays at a point.

---

## 65.2 Types of Lenses

```
CONVEX (CONVERGING) LENS:       CONCAVE (DIVERGING) LENS:

     |  |                            | |
    || ||                           /| |\
   |   |   |                       / |   | \
  |    |    |                     |  |    |  |
  |  -----  |  <- thicker         |  |    |  |  <- thinner
  |    |    |     in middle       |  |    |  |    in middle
   |   |   |                       \ |   | /
    || ||                            \| |/
     |  |                             | |

Brings parallel rays together         Spreads rays apart
at the focal point                    (they appear to come from
                                       a virtual focal point)
```

### Real-World Examples

- **Convex**: magnifying glass, camera lens, eye, projector, telescope objective
- **Concave**: flashlight (to spread light), glasses for nearsightedness, telescope eyepiece

---

## 65.3 Key Terms and Definitions

```
PRINCIPAL AXIS: the straight horizontal line through the center of the lens

OPTICAL CENTRE (O): the center of the lens (light passes through without bending)

FOCAL POINT (F): 
  Convex: point where parallel incident rays converge after passing through lens
  Concave: point from which parallel incident rays appear to diverge after lens
  
  Each lens has two focal points (one on each side), both at distance f.

FOCAL LENGTH (f):
  Distance from optical center to focal point.
  Convex: f is positive (real focus)
  Concave: f is negative (virtual focus)

OBJECT DISTANCE (u): distance from object to lens (always positive for real objects)

IMAGE DISTANCE (v): distance from lens to image
  Positive = image on far side of lens (real image — can be projected on screen)
  Negative = image on same side as object (virtual image — only seen by eye)
```

---

## 65.4 Ray Diagrams for Convex Lenses

Use three principal rays to locate the image:

```
RAY RULES FOR CONVEX LENS:

Ray 1: Parallel to axis → passes through far focal point F
Ray 2: Through optical center O → continues straight, undeviated
Ray 3: Through near focal point F → emerges parallel to axis

        F        O        F
        |         |         |
        |         |         |
object  |         |         |     image
    →   |    →   |    →   |  →
        |         |         |
```

### Case 1: Object Beyond 2F (u > 2f)

```
       F      O    F
       |       |    |
  [obj]|       |    |       [img]
       |       |    |          (smaller, inverted, real)
```

Image: real, inverted, smaller, between F and 2F on far side.
Applications: camera, eye (objects far away)

### Case 2: Object at 2F (u = 2f)

Image: real, inverted, **same size**, at 2F on far side.

### Case 3: Object Between F and 2F (f < u < 2f)

```
       F      O   F
       |       |   |
       |       |   |    [obj]     [img] (beyond 2F, larger, inverted, real)
```

Image: real, inverted, larger, beyond 2F on far side.
Applications: projector, microscope objective

### Case 4: Object at F (u = f)

```
Parallel rays emerge from other side → image at infinity (never converges).
```

### Case 5: Object Inside F (u < f)

```
       F      O   F
       |       |   |
       |       |   | [obj]   
       |       |   |       
       
  Rays diverge; extended back, appear to come from virtual image (same side as object)
```

Image: **virtual, upright, magnified** — on same side as object.
Application: magnifying glass

### Summary Table: Convex Lens

| Object position | Image | Type | Orientation | Size |
|----------------|-------|------|-------------|------|
| Beyond 2F | Between F and 2F | Real | Inverted | Smaller |
| At 2F | At 2F | Real | Inverted | Same |
| Between F and 2F | Beyond 2F | Real | Inverted | Larger |
| At F | Infinity | — | — | — |
| Inside F | Same side | Virtual | Upright | Larger |

---

## 65.5 Ray Diagrams for Concave Lenses

**Concave lenses always form virtual, upright, smaller images** regardless of object position.

```
RAY RULES FOR CONCAVE LENS:

Ray 1: Parallel to axis → diverges as if coming from far focal point F
Ray 2: Through optical center → continues straight
Ray 3: Aimed at far F → emerges parallel to axis

       F       O       F
       |        |       |
  [obj]|        |       |   [virtual img] (smaller, same side as object)
       |        |       |
```

Image: always virtual, upright, smaller, on the same side as the object.

Application: glasses for myopia (nearsightedness), diverging part of telescope eyepiece.

---

## 65.6 The Lens Equation

The **thin lens equation** relates object distance, image distance, and focal length:

```
1/v - 1/u = 1/f

Or equivalently:  1/f = 1/v + 1/u   (using different sign convention)

The most common form (real is positive convention):
  1/f = 1/v + 1/u  (where all are positive for real object/image)
```

**Sign conventions (real-is-positive):**
- u: object distance, always positive for real objects
- v: positive for real image (on far side of lens), negative for virtual image
- f: positive for converging lens, negative for diverging lens

### Worked Example 65.1

An object is placed 30 cm from a convex lens of focal length 10 cm.

Find image distance and state type of image.

**Solution:**

1/f = 1/v + 1/u
1/10 = 1/v + 1/30
1/v = 1/10 - 1/30 = 3/30 - 1/30 = 2/30
v = **15 cm** (positive → real, on far side)

Image is real, inverted, and smaller (v < u).

### Worked Example 65.2

Object 6 cm from a concave lens, f = -12 cm.

**Solution:**

1/v = 1/f - 1/u = 1/(-12) - 1/6 = -1/12 - 2/12 = -3/12 = -1/4

v = **-4 cm** (negative → virtual, same side as object)

---

## 65.7 Magnification

**Magnification (m)** is the ratio of image size to object size:

```
m = image height / object height = v / u

(with signs following the lens convention)
```

- m > 0: image is upright (virtual images from single lenses)
- m < 0: image is inverted (real images from convex lenses)
- |m| > 1: image is magnified
- |m| < 1: image is diminished

### Worked Example 65.3

For Example 65.1 above (u = 30 cm, v = 15 cm):

m = v/u = 15/30 = **0.5**

Image is half the size of object, inverted (negative sign from convention; often written m = -0.5).

---

## 65.8 The Power of a Lens

The **power** of a lens measures how strongly it bends light:

```
P = 1/f

Units: Diopters (D) = m⁻¹

where f must be in METERS
```

- Converging lens: positive power
- Diverging lens: negative power

```
EXAMPLES:
  f = 0.1 m → P = 10 D  (strong converging)
  f = 1 m   → P = 1 D   (weak converging)
  f = -0.5 m → P = -2 D (diverging, glasses for myopia)
  f = 0.25 m → P = 4 D  (reading glasses for hyperopia)
```

When thin lenses are in contact:

```
P_total = P₁ + P₂ + P₃ + ...
```

This is why power (diopters) is more useful than focal length for describing optical systems: powers just add.

---

## 65.9 Lens Defects (Aberrations)

### Spherical Aberration

Rays passing through the edge of a spherical lens focus at a slightly different point from rays through the center.

```
SPHERICAL ABERRATION:

  Edge rays:  → ↘ F_edge
  Center rays: → → F_center
  
  F_edge ≠ F_center → blurry image
  
  Fix: use aspheric lenses (cameras), stop down aperture (reduce to central rays)
```

### Chromatic Aberration

Different colors (wavelengths) refract by slightly different amounts → different focal lengths.

```
CHROMATIC ABERRATION:

  White light →  lens
  
  Violet: f_violet  (shorter focal length, bends more)
  Red:    f_red     (longer focal length, bends less)
  
  → Color fringing around images (visible in cheap lenses/telescopes)
  
  Fix: achromatic doublet — two lenses of different glass types, cancel each other's dispersion
```

---

## 65.10 The Human Eye

The eye is a remarkable optical instrument:

```
EYE STRUCTURE:

  Cornea (does ~2/3 of focusing)
  Pupil (adjustable aperture)
  Lens (adjustable, does ~1/3 of focusing)
  Retina (image sensor)
  
  Total power: ~60 D
  Range of focus: from ~25 cm (near point) to infinity
```

### Accommodation

The lens can change shape (become more or less curved) to adjust focal length — this is called **accommodation**:

- Distant object: ciliary muscles relax, lens becomes flatter (longer f, less power)
- Near object: ciliary muscles contract, lens becomes rounder (shorter f, more power)

### Eye Defects

**Myopia (nearsightedness)**: eyeball too long or lens too curved. Images focus in front of retina. Can see near, not far.
- Correction: concave (diverging) lens, negative power

**Hyperopia (farsightedness)**: eyeball too short. Images focus behind retina. Can see far, not near.
- Correction: convex (converging) lens, positive power

**Presbyopia**: age-related loss of accommodation (lens stiffens). Reading becomes difficult.
- Correction: reading glasses (convex)

**Astigmatism**: cornea is not perfectly spherical (slightly oval). Images distorted.
- Correction: cylindrical lens

---

## Summary

- **Convex (converging) lens**: thicker in middle; brings rays together; f > 0
- **Concave (diverging) lens**: thinner in middle; spreads rays apart; f < 0
- **Ray diagrams**: use 3 principal rays (parallel to axis → F, through center, through F → parallel)
- **Convex lens images**: depend on object distance; real/inverted when u > f; virtual/upright/magnified when u < f
- **Concave lens**: always virtual, upright, smaller
- **Thin lens equation**: 1/f = 1/v + 1/u
- **Magnification**: m = v/u; positive → upright, negative → inverted; |m|>1 → magnified
- **Power**: P = 1/f (diopters); powers of thin lenses in contact add
- **Aberrations**: spherical (edge/center focus differ); chromatic (color dispersion)
- **Eye**: ~60 D total; myopia corrected with concave lens, hyperopia with convex

---

## Key Equations

```
Snell's Law:
  n₁ sin(θ₁) = n₂ sin(θ₂)

Thin lens equation:
  1/f = 1/v + 1/u

Magnification:
  m = image height / object height = v / u

Power of a lens:
  P = 1/f   (f in meters, P in diopters D)

Lenses in contact:
  P_total = P₁ + P₂ + ...

Sign conventions:
  Real object: u > 0
  Real image: v > 0
  Virtual image: v < 0
  Converging lens: f > 0
  Diverging lens: f < 0
```
