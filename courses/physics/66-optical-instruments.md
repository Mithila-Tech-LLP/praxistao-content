# Chapter 66: Optical Instruments

> **"A microscope extends sight downward into the invisible world of cells and bacteria. A telescope extends it outward to distant stars and galaxies. Both use the same lens physics — just arranged differently."**

---

## Table of Contents

- [66.1 The Simple Magnifier](#661-the-simple-magnifier)
- [66.2 The Compound Microscope](#662-the-compound-microscope)
- [66.3 The Astronomical Telescope](#663-the-astronomical-telescope)
- [66.4 The Galilean (Opera) Telescope](#664-the-galilean-opera-telescope)
- [66.5 The Reflecting Telescope](#665-the-reflecting-telescope)
- [66.6 Cameras](#666-cameras)
- [66.7 Projectors](#667-projectors)
- [66.8 The Human Eye with Spectacles](#668-the-human-eye-with-spectacles)
- [66.9 Resolution and Diffraction](#669-resolution-and-diffraction)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 66.1 The Simple Magnifier

A single convex lens used as a **magnifying glass**.

Hold the object inside the focal length (u < f). The lens produces a virtual, upright, magnified image.

### Angular Magnification

The "power" of a magnifier is described by its **angular magnification** — the ratio of the angle subtended at the eye with the lens to without it.

```
ANGULAR MAGNIFICATION:

Without magnifier: object at near point D (25 cm for normal eye)
  
                angle α
  [object] ←--D=25cm--→ [eye]
  
With magnifier: object at focal point F of lens
  
  [object]←f→[lens]          [eye]
              angle β (larger)
  
  M = β/α = D/f  (approximately, when final image at infinity)
  M = D/f = 25 cm / f
```

For comfortable viewing (image at infinity):

```
M = D/f = 25/f (cm)

where f is in cm and D = 25 cm (near point of normal eye)
```

A lens with f = 5 cm gives M = 25/5 = 5× magnification.

### Worked Example 66.1

A jeweler's loupe has f = 2.5 cm. What magnification does it provide?

M = 25/2.5 = **10×**

---

## 66.2 The Compound Microscope

For greater magnification than a simple magnifier, use TWO lenses — an objective and an eyepiece.

```
COMPOUND MICROSCOPE:

  Objective lens (short focal length f_obj):
    - Object placed just outside f_obj
    - Produces real, inverted, highly magnified intermediate image
    
  Eyepiece (short focal length f_eye):
    - Acts as a simple magnifier
    - Intermediate image placed at f_eye
    - Produces virtual, upright final image for the eye to view
```

```
MICROSCOPE DIAGRAM:

  [object] → [objective lens] → [intermediate image] → [eyepiece] → [eye]
  
    ←  u  →←        L        →←          v_e         →

  u: object distance (slightly > f_obj)
  L: tube length (distance between lenses minus their focal lengths)
  v_e: image distance in eyepiece
```

### Total Magnification

```
M_total = M_objective × M_eyepiece

M_objective = L / f_obj   (approximately, for object near f_obj)

M_eyepiece = D / f_eye = 25 / f_eye

M_total = (L × D) / (f_obj × f_eye) = L/f_obj × 25/f_eye
```

### Worked Example 66.2

A microscope has objective f = 0.5 cm, eyepiece f = 2.5 cm, tube length L = 15 cm.

M_total = (L/f_obj) × (25/f_eye) = (15/0.5) × (25/2.5) = 30 × 10 = **300×**

The final image is inverted (inverted by objective, then magnified but still inverted by eyepiece).

---

## 66.3 The Astronomical Telescope

Used to observe distant objects. Brings them closer in angular size.

```
ASTRONOMICAL TELESCOPE:

  Objective lens (long focal length):
    - Parallel rays from distant object → form real image at focal plane
    
  Eyepiece (short focal length):
    - Acts as magnifier of the intermediate image
    
  Separation = f_obj + f_eye (both objects at infinity → image at infinity = relaxed eye)
```

```
TELESCOPE DIAGRAM:

  Parallel rays →→→→ [objective f_obj] → [intermediate image at F_obj]
                                              ↕ (this is the "object" for eyepiece)
                                         [eyepiece f_eye] →→→→→→ [eye]
                                          ↑
                              F_obj and F_eye coincide here
```

### Angular Magnification

```
M = f_obj / f_eye

where f_obj >> f_eye (objective has long focal length, eyepiece has short focal length)
```

The image is inverted (can be accepted for astronomy — stars don't have an "up").

### Worked Example 66.3

An astronomical telescope: objective f = 100 cm, eyepiece f = 5 cm.

M = f_obj / f_eye = 100/5 = **20×**

The moon (angular diameter 0.5°) would appear as 20 × 0.5° = 10° in diameter — much larger!

---

## 66.4 The Galilean (Opera) Telescope

Uses a convex objective lens and a **concave** eyepiece lens. Produces an upright, smaller image. Used for opera glasses, binoculars (with prisms to fold the path).

```
GALILEAN TELESCOPE:

  Parallel rays → [convex objective] → (converging) → [concave eyepiece] → [eye]
  
  The concave eyepiece intercepts the converging rays before they meet,
  producing an upright, virtual image.
  
  M = f_obj / f_eye  (same formula, but f_eye is negative → M < f_obj/|f_eye|)
```

Advantages:
- Shorter tube length (tubes overlap)
- Upright image (important for terrestrial use)

---

## 66.5 The Reflecting Telescope

Large astronomical telescopes use a **concave mirror** instead of an objective lens:

```
REFLECTOR TELESCOPE (Newtonian):

  Parallel rays →→→→ [concave primary mirror] → reflects and converges
                                               ↓
                                      [small flat secondary mirror]
                                               ↓
                                       (diverts to side)
                                               ↓
                                         [eyepiece at side]
```

Advantages over refractors:
- No chromatic aberration (mirrors don't disperse colors)
- Cheaper to make large mirrors than large lenses
- Mirror supported from behind → can be very large
- World's largest telescopes (VLT, Keck, ELT) are all reflectors

The **focal ratio** (f-number = f/D) describes the "speed" of a telescope:
- Small f-number: wide field, good for astrophotography
- Large f-number: narrow field, high magnification

---

## 66.6 Cameras

A camera uses a convex lens to form a real, inverted image on a sensor (film or digital detector).

```
CAMERA:

  [scene at large distance u] → [lens f] → [sensor at v]
  
  1/f = 1/v + 1/u
  
  For far objects: u >> f, so v ≈ f (image at focal length)
  For closer objects: v > f (lens moves forward)
  
  "Focus" = adjust lens distance from sensor until v matches object distance
```

### Aperture (f-number)

The **aperture** controls how much light enters:

```
f-number (N) = f / D   (focal length / aperture diameter)

Small N (f/1.4, f/2):  large aperture → more light → shallow depth of field
Large N (f/11, f/16):  small aperture → less light → deep depth of field

Exposure time needed ∝ N²  (double N → 4× longer exposure needed)
```

### Depth of Field

At small aperture (large N), objects at different distances are all in focus simultaneously — large **depth of field**. Useful for landscapes.

At large aperture (small N), only objects at one distance are sharp — shallow depth of field. Used for portrait photography (blurry background).

### Worked Example 66.4

A camera with f = 50 mm lens photographs an object 1 m away.

How far should the sensor be from the lens?

1/v = 1/f - 1/u = 1/0.05 - 1/1 = 20 - 1 = 19

v = 1/19 ≈ **0.0526 m = 52.6 mm**

The sensor is slightly farther than the focal length (50 mm vs 52.6 mm).

---

## 66.7 Projectors

A **projector** creates a large, real, inverted image on a screen from a small, bright object (slide, LCD panel).

```
PROJECTOR:

  [small bright slide near 2f] → [lens] → [large image on screen far away]
  
  Object distance: just outside f (but close to f)
  Image distance: very large
  
  Magnification: m = v/u (very large, since v >> u)
```

The slide must be placed inverted in the projector — so the projected image appears the right way up.

### Modern Projectors

Digital projectors use an LCD or DLP chip as the "slide":
- DLP (Digital Light Processing): millions of tiny mirrors that tilt to direct light
- LCD: liquid crystal panel that controls light transmission
- Each pixel in the display corresponds to a pixel on the projected image

---

## 66.8 The Human Eye with Spectacles

For a person with a vision defect, spectacles (glasses) modify the incoming light so the eye can focus it properly.

### Myopia (Nearsightedness)

```
MYOPIA:
  
  Far object → light converges in FRONT of retina (eyeball too long)
  
  Correction: concave lens spreads the rays → moves focus back to retina
  
  Power needed: P_lens = -(1/far_point_in_meters)
  
  Example: Far point = 2 m (can't see beyond 2 m):
    P = -1/2 = -0.5 D
```

### Hyperopia (Farsightedness)

```
HYPEROPIA:
  
  Near object → light converges BEHIND retina (eyeball too short)
  
  Correction: convex lens converges rays more → moves focus forward to retina
  
  Power needed: more complex (depends on near point), typically +1 to +5 D
```

### Presbyopia

The eye lens stiffens with age, reducing accommodation (ability to change focal length). Reading glasses (convex, typically +1 to +3 D) bring near objects into focus.

---

## 66.9 Resolution and Diffraction

An optical instrument's ability to distinguish fine detail is limited by **diffraction** — light bends around the edges of the aperture, creating a diffraction pattern.

### Rayleigh Criterion

Two point sources are just resolved when the central maximum of one overlaps with the first minimum of the other:

```
Angular resolution (minimum resolvable angle):

θ_min = 1.22 × λ / D

where:
  λ = wavelength of light
  D = diameter of aperture (lens or mirror)
  
RESOLUTION LIMIT:
  Smaller θ_min = better resolution (can distinguish closer objects)
  
  To improve resolution:
    - Use shorter λ (blue light, UV, X-rays)
    - Use larger D (bigger mirror/lens)
```

### Examples

```
HUMAN EYE: D ≈ 2-8 mm (pupil)
  θ_min ≈ 1.22 × 500×10⁻⁹ / 0.003 ≈ 2×10⁻⁴ rad ≈ 0.01°
  At 25 cm: can resolve features 0.05 mm apart (about 50 μm)

HUBBLE SPACE TELESCOPE: D = 2.4 m
  θ_min ≈ 1.22 × 500×10⁻⁹ / 2.4 ≈ 2.5×10⁻⁷ rad
  Much sharper than any ground-based telescope!
  
RADIO TELESCOPE (100 m dish, λ = 1 m):
  θ_min ≈ 1.22 × 1 / 100 ≈ 0.01 rad ≈ 0.6°  (poor resolution)
  
  Solution: Very Long Baseline Interferometry (VLBI) — use multiple telescopes
  thousands of km apart → D = Earth's diameter → amazing resolution!
```

---

## Summary

- **Simple magnifier**: single convex lens; M = D/f where D = 25 cm; virtual, upright, magnified image
- **Compound microscope**: two lenses; M = (L/f_obj) × (D/f_eye); high magnification, inverted image
- **Astronomical telescope**: two lenses or objective mirror; M = f_obj/f_eye; inverted, for stars
- **Galilean telescope**: convex objective + concave eyepiece; upright image; used in opera glasses/binoculars
- **Reflecting telescope**: concave primary mirror; no chromatic aberration; allows very large apertures
- **Camera**: forms real, inverted image on sensor; aperture controls light and depth of field
- **Projector**: magnified real image from small bright source; object placed just outside F
- **Eye correction**: myopia (concave, negative P); hyperopia (convex, positive P)
- **Resolution limit**: θ_min = 1.22λ/D; larger aperture or shorter wavelength → better resolution

---

## Key Equations

```
Simple magnifier:
  M = D/f  (D = 25 cm for normal eye)

Compound microscope:
  M = M_obj × M_eye = (L/f_obj) × (D/f_eye)

Astronomical telescope:
  M = f_obj / f_eye

Camera f-number:
  N = f / D  (aperture diameter D)

Rayleigh resolution criterion:
  θ_min = 1.22 × λ / D

Thin lens equation (applies to all):
  1/f = 1/v + 1/u

Power (diopters):
  P = 1/f  (f in meters)

Correction for myopia (far point = far_point meters):
  P_lens = -1/far_point
```
