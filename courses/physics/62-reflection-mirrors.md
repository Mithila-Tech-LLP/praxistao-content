# Chapter 62: Reflection and Mirrors

> **"Mirrors seem like magic — but they're just physics. The law of reflection is so simple a child can state it, yet it explains everything from bathroom mirrors to giant telescopes."**

---

## Table of Contents

1. [The Law of Reflection](#1-the-law-of-reflection)
2. [Specular vs Diffuse Reflection](#2-specular-vs-diffuse-reflection)
3. [Plane Mirrors](#3-plane-mirrors)
4. [Curved Mirrors: Introduction](#4-curved-mirrors-introduction)
5. [Concave Mirrors (Converging)](#5-concave-mirrors-converging)
6. [Convex Mirrors (Diverging)](#6-convex-mirrors-diverging)
7. [The Mirror Equation and Sign Convention](#7-the-mirror-equation-and-sign-convention)
8. [Worked Examples](#8-worked-examples)
9. [Applications of Mirrors](#9-applications-of-mirrors)
10. [Summary](#summary)
11. [Key Equations](#key-equations)

---

## 1. The Law of Reflection

### What Happens When Light Hits a Surface?

When a ray of light travels through air and strikes a smooth surface — such as a mirror, a still pond, or a polished piece of metal — it does not pass through. Instead, it **bounces back**. This bouncing of light off a surface is called **reflection**.

Reflection is not random. It follows a precise and elegant rule that never breaks: the **Law of Reflection**.

---

### The Normal Line

Before we state the law, we need to understand one crucial concept: the **normal**.

The **normal** is an imaginary line drawn **perpendicular (at 90°) to the mirror surface** at the exact point where the incoming light ray strikes it. Think of it as a flagpole sticking straight out of the mirror at the point of contact.

Every angle in reflection is measured **from the normal**, not from the surface itself. This is one of the most common errors beginners make, so it's worth emphasizing right away.

---

### Stating the Law

When a ray of light hits a reflective surface:

- The **angle of incidence (θᵢ)** is the angle between the incoming ray and the normal.
- The **angle of reflection (θᵣ)** is the angle between the reflected ray and the normal.
- The **incident ray**, the **normal**, and the **reflected ray** all lie in the same flat plane.

**The Law of Reflection: θᵢ = θᵣ**

That's it. The angle coming in equals the angle going out — both measured from the normal.

---

### ASCII Diagram 1: The Law of Reflection

```
         Incident Ray          Normal          Reflected Ray
                 \                |                /
                  \               |               /
                   \      θᵢ     |     θᵣ      /
                    \      ___   |   ___      /
                     \  __|   |  |  |   |__ /
                      \ |  θᵢ | |  | θᵣ  | /
                       \|_____|  |  |_____|/
                        \        |        /
                         \       |       /
                          \      |      /
                           \_____|_____/
    ====================================================  ← Mirror Surface
```

A cleaner view:

```
                     NORMAL
                       |
  Incident Ray         |         Reflected Ray
          \            |            /
           \    θᵢ     |    θᵣ    /
            \          |          /
             \         |         /
              \        |        /
               \       |       /
                \      |      /
  _______________\_____|_____/_______________
                  Mirror Surface
                  
  θᵢ = angle of incidence  (measured from normal)
  θᵣ = angle of reflection (measured from normal)
  Law: θᵢ = θᵣ  (always!)
```

---

### Common Mistake: Measuring From the Surface

Many students instinctively measure angles from the mirror surface, not the normal. This gives wrong answers. Here's why the normal is used:

Suppose a ray hits a mirror at 30° from the surface. That means it's 60° from the normal. The reflected ray will also be 60° from the normal — meaning 30° from the surface on the other side. The angles from the surface are equal **too**, but the normal convention is universal because it works even for curved surfaces (where the surface is not flat, but the normal at any point is always well-defined).

**Rule to remember:** Normal = perpendicular to surface. All reflection angles measured from the normal.

---

### Does Color Matter?

No. The Law of Reflection applies to **all wavelengths** of light — red, blue, ultraviolet, infrared, radio waves. The angle of reflection equals the angle of incidence regardless of the color or frequency of the light. This universality makes reflection extremely predictable and useful.

---

## 2. Specular vs Diffuse Reflection

Not all surfaces reflect the same way. The key difference is **smoothness**.

---

### Specular Reflection (Mirror-Like)

When a surface is **very smooth** — like a polished mirror, calm water, or a shiny metal surface — all the incoming parallel rays hit surfaces whose normals all point in the same direction. Each ray reflects at the same angle. The reflected rays are also parallel.

This is **specular reflection**. Because all reflected rays travel in a coherent, organized direction, your eye can trace them back to a clear source — you see a sharp, clear **image**. This is how mirrors work.

---

### Diffuse Reflection

When a surface is **rough** — like paper, a painted wall, cloth, or skin — the surface looks smooth to the naked eye, but at the microscopic level it's full of tiny bumps and valleys. Each tiny patch of surface faces a slightly different direction, so each tiny patch has its own normal pointing in a different direction.

Parallel incoming rays hit these different micro-patches and reflect off in **many different directions**. This is **diffuse reflection**.

---

### ASCII Diagram 2: Specular vs Diffuse Reflection

```
  SPECULAR REFLECTION (smooth surface)        DIFFUSE REFLECTION (rough surface)
  
  Incoming rays (parallel):                   Incoming rays (parallel):
    ↓  ↓  ↓  ↓  ↓                              ↓  ↓  ↓  ↓  ↓
    |  |  |  |  |                              |  |  |  |  |
  --+--+--+--+--+--  ← smooth mirror         /\/\/\/\/\/\/\/\  ← rough surface
    ↑  ↑  ↑  ↑  ↑                            ↗ ↖ ↗ ↘ ↙ ↗ ↘ ↙
    
  Reflected rays are                          Reflected rays go in
  parallel → clear image                      ALL directions → no image,
                                              surface looks uniformly lit
```

---

### Why Can We See Non-Luminous Objects?

Most objects around you — books, walls, people — don't emit their own light. Yet you can see them. How?

They work by **diffuse reflection**. Sunlight or a lamp shines on them. Their rough surfaces scatter the light in all directions. Some of that scattered light enters your eye, and your brain perceives the object. Because the light goes in all directions, you can see the object from any angle — unlike a mirror where you only see the reflection from specific angles.

This is the main reason you can see the world around you. Nearly everything you observe (other than the sun, lamps, and screens) is visible through diffuse reflection.

---

### Why Rough Surfaces Don't Form Images

Each tiny patch on a rough surface does obey the Law of Reflection perfectly. The problem is that neighboring patches face different directions, so reflected rays from a single incoming ray go in chaotic directions. There is no organized convergence to form an image. The surface appears uniformly illuminated rather than showing a reflection.

---

## 3. Plane Mirrors

A **plane mirror** is simply a flat, smooth, reflective surface. Your bathroom mirror is a plane mirror. Despite its simplicity, it has some fascinating and precise optical properties.

---

### Properties of the Image in a Plane Mirror

When you look into a plane mirror, the image you see has these exact properties:

1. **Virtual** — The image appears to exist behind the mirror, but no light actually travels there. You cannot project it onto a screen.

2. **Upright** — The image is the same way up as the object. Your head appears at the top, your feet at the bottom.

3. **Same size as the object** — The **magnification** is exactly 1. If you're 1.7 m tall, your image is also 1.7 m tall.

4. **Same distance behind the mirror as the object is in front** — If you stand 2 m from the mirror, your image appears to be 2 m behind the mirror's surface. The total distance between you and your image is 4 m.

5. **Laterally inverted** — Left and right are swapped. Your right hand appears to be the image's left hand. This is why text in a mirror appears backwards.

---

### ASCII Diagram 3: Image Formation in a Plane Mirror

```
                          PLANE MIRROR
                               |
  Object (arrow)               |              Image (virtual arrow)
                               |
         ↑                     |                    ↑
         |  O (top of          |          O' (image of top) |
         |   arrow)            |                            |
         |           ray 1 →   |   ← ray 1 (reflected)      |
         |           - - - - - ↑ - - - - - - - - - - →      |
         |                     |  \       (extended back)   |
         |           ray 2 →   |    \                       |
         |           - - - - ↗ |     \  - - - - - - - →    |
         |                  ↗  |      \                     |
    Eye ←←←←←←←←←←←←←←←←←   |       - - - - - → (Image O')
                               |
         |←   d (object)  →|   |   |←  d (image)  →|
         
  The two reflected rays diverge as they leave the mirror.
  Your eye traces them back (dotted extensions) and they 
  converge behind the mirror → that's where the image appears.
  The image is VIRTUAL: no real light exists behind the mirror.
```

---

### Why Is the Image Virtual?

After reflecting off the mirror, the two rays **diverge** (spread apart). Your eye receives these diverging rays. Your brain automatically traces them back in straight lines until they meet — and they meet behind the mirror at the image point. But there is no actual light at that point. The image is a **virtual image**.

Compare this to a **real image**, which would be formed by rays that actually converge at a point (as we'll see with concave mirrors). A real image can be projected onto a screen. A virtual image cannot.

---

### The Left-Right Reversal Mystery

Why do mirrors reverse left-right but not up-down? This puzzle confuses many people.

The honest answer: mirrors don't reverse left-right. They reverse **front-to-back (depth)**. When you face a mirror, your image faces you — meaning its front is pointing toward you, while your front points toward the mirror. That's a front-to-back reversal.

The left-right confusion arises because when a person turns around (which is how you imagine "stepping into" your reflection), they rotate about their vertical axis — swapping left and right. But they don't somersault (rotate about the horizontal axis) — that's why up-down stays the same.

If you held a page of text flat and lowered it face-down onto a mirror, you'd see the text flipped top-to-bottom. The mirror treats every axis the same — it's our imagination of "walking into" the mirror that introduces asymmetry.

---

### Multiple Mirror Images

When two plane mirrors are placed at an angle θ to each other, multiple images are formed due to repeated reflections.

**Formula:** Number of images = (360° / θ) − 1    (when 360°/θ is a whole number)

**Worked Example:**
Two mirrors are placed at 60° to each other. How many images does an object placed between them produce?

- Number of images = 360° / 60° − 1 = 6 − 1 = **5 images**

At 90°: images = 360/90 − 1 = 3 images (seen in dressing table mirrors)
At 45°: images = 360/45 − 1 = 7 images (kaleidoscope effect)
At 0° (parallel mirrors): infinite images, as in a hotel elevator with mirrors on both walls.

---

## 4. Curved Mirrors: Introduction

Flat mirrors are limited — they always produce same-size images at the same distance. **Curved mirrors** allow us to produce magnified, diminished, real, and virtual images by changing the mirror's curvature. They are among the most powerful optical tools in existence.

---

### Key Terms for Curved Mirrors

Curved mirrors are sections of a **sphere**. Imagine cutting a small piece from the inside or outside of a large ball — you get a curved mirror.

- **Centre of Curvature (C):** The center of the full sphere that the mirror is a part of. All radii of the sphere pass through C.

- **Radius of Curvature (R):** The radius of that sphere — the distance from C to any point on the mirror's surface.

- **Pole (P):** The geometric center of the mirror's curved surface. The point where the principal axis meets the mirror.

- **Principal Axis:** The straight line passing through P and C (and beyond). This is the mirror's axis of symmetry.

- **Focal Point (F):** The special point on the principal axis where parallel rays converge (or appear to diverge from) after reflection. It lies exactly halfway between P and C.

- **Focal Length (f):** The distance PF = R / 2. This is the most important single number describing a curved mirror.

---

### ASCII Diagram 4: Concave Mirror — Key Points Labeled

```
                              Principal Axis
    ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←→
                                                      
    C             F        P
    |             |        |
    ●             ●        |)  ← concave mirror surface
    |             |        |
    |←    R/2   →|←  R/2 →|
    |←         R          →|
    
    Where:
      P = Pole (center of mirror surface)
      F = Focal point (midpoint of PC)
      C = Centre of curvature
      R = Radius of curvature
      f = Focal length = PF = R/2
```

---

## 5. Concave Mirrors (Converging)

### What Is a Concave Mirror?

A **concave mirror** curves **inward**, like the inside of a bowl or a cave (the word "concave" comes from the Latin for hollow). If you look at a spoon from the inside of the bowl, you see a concave reflecting surface.

The key property: **parallel rays converge at the focal point** after reflecting off a concave mirror. This is why concave mirrors are also called **converging mirrors**.

The focal length f of a concave mirror is taken as **positive** in the sign convention we use.

---

### The Three Principal Rays

To draw a ray diagram for a concave mirror, you use three special rays. Any two of them are enough to find the image; the third serves as a check.

**Ray 1 (Parallel Ray):**
A ray traveling parallel to the principal axis → after reflection, it passes through the focal point F.

**Ray 2 (Focal Ray):**
A ray passing through (or heading toward) the focal point F → after reflection, it travels parallel to the principal axis.

**Ray 3 (Center Ray):**
A ray passing through the center of curvature C → it hits the mirror perpendicularly (along the normal), so it reflects straight back through C.

---

### Image Formation Table for Concave Mirror

| Object Position | Image Position | Nature | Orientation | Size |
|----------------|---------------|--------|-------------|------|
| Beyond C (u > 2f) | Between F and C | Real | Inverted | Diminished |
| At C (u = 2f) | At C | Real | Inverted | Same size |
| Between C and F (f < u < 2f) | Beyond C | Real | Inverted | Magnified |
| At F (u = f) | At infinity | Real | Inverted | Infinitely large |
| Inside F (u < f) | Behind mirror | Virtual | Upright | Magnified |

---

### ASCII Diagram 5: Concave Mirror — Object Beyond C (Real, Diminished Image)

```
Object                                           Concave Mirror
  |                                                    /)
  |  (Object arrow                                    / )
  ↑  pointing up)        Ray 1 (parallel) →          /  )
  |  ←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→  /   )
  |  ←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←↙         )
  |                                       ↙ (through F) )
  |                                      ↙              )
  |  ←→←→←→←→←→←→←→←→←→←→ (Ray 3 through C)          )
  |  ←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→←→  \   )
  |                                             \    )
C●              F●              P●              \ ← Reflected back through C
|               |               |               \  )
←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←  ←←←←

Simplified diagram (object beyond C):

               u > 2f                   
       ↑                                    /)
       |   Ray 1 →→→→→→→→→→→→→→→→→→→→→↘  /
       |                                  ↘ / → through F
       |   Ray 2 →→→→→→→→→→→→→→→→→→→↗  /
       |                             ↗  /  ← came from F, now goes parallel
     Object                     Image (real, inverted, smaller)
       ←   →                   ↓
  C          F       P    
  ●----------●-------●---------↓ (image forms here, between F and C)
  
  Real image → rays ACTUALLY cross → can project onto screen
```

---

### ASCII Diagram 6: Concave Mirror — Object Inside F (Virtual, Magnified Image)

```
Object inside F:

                Ray 1 →→→→→→→→→→→→→→→→→→→→→→→→→→→→→↘  /)
                                                       ↘ ) → reflects through F
                                                    ←←←←)← diverging reflected ray 1
                                                  ↙     )
  Image                            Object         )
  (virtual,       ←←←←←←←←←←←←←←←  ↑     Ray 2)/
  upright,                            | ←←←←←←← ) → reflects parallel
  magnified)                          |           )
  behind mirror                       |      P    )
  ●                                   ●-----------●
          C                   F       ↑     
          ●                   ●       Object
          
  The reflected rays DIVERGE.
  Trace them BACK (dotted lines) behind the mirror → they meet at the VIRTUAL IMAGE.
  
  Dotted backward extensions:
        ←- - - - - - ←←←← diverging ray 1 (extended back)
  Image ●
        ←- - - - - - ←←←← diverging ray 2 (extended back)
  
  This is how a magnifying/makeup mirror works: object inside F → large virtual image
```

---

## 6. Convex Mirrors (Diverging)

### What Is a Convex Mirror?

A **convex mirror** curves **outward**, like the outside of a ball or the back of a spoon. When you look at a spoon from the back (the convex side), you see your reflection — small and upright. That's a convex mirror at work.

The key property: **parallel rays diverge** after reflecting off a convex mirror, but they appear to come from a focal point **behind** the mirror (a virtual focal point). This is why convex mirrors are also called **diverging mirrors**.

The focal length f is taken as **negative** for a convex mirror.

---

### Always Virtual, Upright, and Diminished

No matter where you place an object in front of a convex mirror, the image is always:
- **Virtual** (behind the mirror)
- **Upright** (same orientation as the object)
- **Diminished** (smaller than the object)

This makes convex mirrors ideal where you want a wide-angle view — such as car rear-view mirrors, security mirrors in shops, and road safety mirrors at blind corners.

---

### ASCII Diagram 7: Convex Mirror Ray Diagram

```
                                  Convex Mirror
                                       (
Parallel Ray 1 →→→→→→→→→→→→→→→→→→→→→ ( ↗ diverges outward
                                       (
Parallel Ray 2 →→→→→→→→→→→→→→→→→→→→→ ( ↗ diverges outward
                                       (
                                  P ●  (
                                       (
     Behind mirror (virtual):         (
                   F ●                (  R ●
                   (virtual focus)         (centre of curvature)
  
  Reflected rays DIVERGE. Extended backward (dotted lines):
  
  →→→→→→→→→→→→→→→→→→→→ (   ↗  ← actual reflected ray
                         (
                    - - -●- - - → (both dotted extensions meet at F behind mirror)
                         (
  →→→→→→→→→→→→→→→→→→→→ (   ↗  ← actual reflected ray
  
  F is BEHIND the mirror → it's a VIRTUAL focus → f is NEGATIVE
  
  Object and Image (convex mirror):
  
  Object ↑                               (
         |   →→→→→→→→→→→→→→→→→→→→→→→→  ( ↗ diverges
         |                               (
         |   →→→→→→→→→→→→→→→→→→→→→→→→  ( ↗ diverges
         |                               (
         ●             ●            P ●  (
      Object           Image             
                    (virtual,            
                    upright,             
                    diminished)          
         ←    u    →←  |v|  →←    →
```

---

## 7. The Mirror Equation and Sign Convention

### The Mirror Equation

The relationship between object distance, image distance, and focal length is captured in the **mirror equation** (also called the mirror formula):

```
    1     1     1
   --- = --- + ---
    f     v     u
```

Where:
- **f** = focal length (distance from pole to focal point)
- **u** = object distance (distance from pole to object)
- **v** = image distance (distance from pole to image)

This equation works for both concave and convex mirrors — as long as you apply the sign convention correctly.

---

### Sign Convention (Real is Positive)

We use a simple and consistent sign convention throughout:

| Quantity | Sign Rule |
|---------|-----------|
| Object distance (u) | Always **positive** for real objects (in front of mirror) |
| Image distance (v) | **Positive** if image is real (in front of mirror); **negative** if virtual (behind mirror) |
| Focal length (f) | **Positive** for concave mirror; **negative** for convex mirror |
| Radius of curvature (R) | Same sign as f |

---

### Magnification

The **linear magnification** (m) tells you the size of the image relative to the object:

```
         image height       -v
    m = -------------- = ------
        object height       u
```

Interpreting the sign of m:
- **m is negative** → image is **inverted** (upside down) relative to object
- **m is positive** → image is **upright** (same way up as object)

Interpreting the magnitude of m:
- **|m| > 1** → image is **magnified** (larger than object)
- **|m| < 1** → image is **diminished** (smaller than object)
- **|m| = 1** → image is the **same size** as the object

---

### Quick Reference

```
Concave mirror:  f > 0  (positive focal length)
Convex mirror:   f < 0  (negative focal length)

Real image:      v > 0  (in front of mirror, real rays converge there)
Virtual image:   v < 0  (behind mirror, only extended rays meet there)

Inverted image:  m < 0
Upright image:   m > 0
```

---

## 8. Worked Examples

### Example 1: Concave Mirror — Object Beyond F (Real Image)

**Problem:** A concave mirror has a focal length of 10 cm. An object is placed 30 cm in front of the mirror. Find the image distance and magnification. Describe the image.

**Given:**
- f = +10 cm (concave mirror → positive)
- u = +30 cm (real object → positive)

**Step 1: Apply the mirror equation**

```
1/f = 1/v + 1/u

1/10 = 1/v + 1/30

1/v = 1/10 - 1/30
```

**Step 2: Calculate 1/v**

```
1/v = 3/30 - 1/30 = 2/30
```

**Step 3: Find v**

```
v = 30/2 = 15 cm
```

v is **positive** → the image is **real** and forms **15 cm in front of the mirror**.

**Step 4: Find magnification**

```
m = -v/u = -15/30 = -0.5
```

**Interpretation:**
- m = -0.5 → **negative** means the image is **inverted**
- |m| = 0.5 < 1 → image is **diminished** (half the size)
- v = +15 cm → **real** image (in front of mirror)

**Summary:** The image is **real, inverted, diminished (half size), and located 15 cm in front of the mirror.** (Object was at 30 cm, which is beyond C at 20 cm. This matches our table: object beyond C → real, inverted, diminished, between F and C.)

---

### Example 2: Concave Mirror — Object Inside F (Virtual Image / Magnifying Mirror)

**Problem:** A concave mirror has a focal length of 15 cm. An object is placed 10 cm in front of the mirror. Find the image distance and magnification. Describe the image.

**Given:**
- f = +15 cm (concave mirror → positive)
- u = +10 cm (real object → positive)
- Note: u = 10 cm < f = 15 cm → object is inside F

**Step 1: Apply the mirror equation**

```
1/f = 1/v + 1/u

1/15 = 1/v + 1/10

1/v = 1/15 - 1/10
```

**Step 2: Calculate 1/v**

```
1/v = 2/30 - 3/30 = -1/30
```

**Step 3: Find v**

```
v = -30 cm
```

v is **negative** → the image is **virtual** and forms **30 cm behind the mirror**.

**Step 4: Find magnification**

```
m = -v/u = -(-30)/10 = +30/10 = +3
```

**Interpretation:**
- m = +3 → **positive** means the image is **upright**
- |m| = 3 > 1 → image is **magnified** (three times the size)
- v = -30 cm → **virtual** image (behind mirror)

**Summary:** The image is **virtual, upright, and magnified three times, appearing 30 cm behind the mirror.** This is exactly how a **magnifying makeup mirror** or a **dental mirror** works — the object (your face or a tooth) is placed inside the focal length, producing a large, upright virtual image.

---

### Example 3: Convex Mirror — Always Virtual

**Problem:** A convex mirror used as a car wing mirror has a focal length of 20 cm. A car behind is 30 cm from the mirror. Find the image distance and magnification. Describe the image.

**Given:**
- f = -20 cm (convex mirror → negative)
- u = +30 cm (real object → positive)

**Step 1: Apply the mirror equation**

```
1/f = 1/v + 1/u

1/(-20) = 1/v + 1/30

1/v = 1/(-20) - 1/30
```

**Step 2: Calculate 1/v**

```
1/v = -1/20 - 1/30
    = -3/60 - 2/60
    = -5/60
```

**Step 3: Find v**

```
v = -60/5 = -12 cm
```

v is **negative** → the image is **virtual** and forms **12 cm behind the mirror**.

**Step 4: Find magnification**

```
m = -v/u = -(-12)/30 = +12/30 ≈ +0.4
```

**Interpretation:**
- m = +0.4 → **positive** means the image is **upright**
- |m| = 0.4 < 1 → image is **diminished** (about 40% of actual size)
- v = -12 cm → **virtual** image (behind mirror)

**Summary:** The image is **virtual, upright, and diminished.** Objects appear smaller than they really are — and therefore **closer than they actually are**. This is exactly why car wing mirrors carry the warning: **"Objects in mirror are closer than they appear."**

The benefit: a convex mirror gives a wide **field of view**. Because the mirror curves outward, it captures light from a large angular range and shows you a wide panorama — much wider than a flat mirror of the same size would show.

---

### Example 4: Finding the Focal Length from Image Data

**Problem:** An object placed 20 cm in front of a mirror produces a real image 60 cm in front of the mirror. Is the mirror concave or convex? What is its focal length?

**Given:**
- u = +20 cm (real object)
- v = +60 cm (real image, in front of mirror → positive)

**Step 1: Apply mirror equation**

```
1/f = 1/v + 1/u = 1/60 + 1/20
```

**Step 2: Calculate**

```
1/f = 1/60 + 3/60 = 4/60 = 1/15

f = 15 cm
```

f = +15 cm → **positive** → the mirror is **concave** (only a concave mirror forms real images of real objects).

**Bonus — magnification:**

```
m = -v/u = -60/20 = -3
```

The image is inverted and 3 times the size of the object.

---

### Example 5: Radius of Curvature

**Problem:** A concave mirror has a radius of curvature of 36 cm. An object is placed 27 cm in front. Find the image.

**Given:**
- R = 36 cm → f = R/2 = 18 cm → f = +18 cm (concave)
- u = +27 cm

**Step 1: Mirror equation**

```
1/f = 1/v + 1/u

1/18 = 1/v + 1/27

1/v = 1/18 - 1/27
```

**Step 2: Common denominator (LCM of 18 and 27 = 54)**

```
1/v = 3/54 - 2/54 = 1/54

v = 54 cm
```

**Step 3: Magnification**

```
m = -v/u = -54/27 = -2
```

**Summary:** Real image, 54 cm in front of mirror, inverted, magnified twice. (Object between F and C → real, inverted, magnified → confirmed.)

---

## 9. Applications of Mirrors

### Concave Mirror Applications

**1. Reflecting Telescopes**
Large telescopes (like the Hubble Space Telescope) use giant concave mirrors instead of lenses to collect light from distant stars and galaxies. Mirrors can be made much larger than lenses and don't suffer from **chromatic aberration** (color fringing). The larger the mirror, the more light collected and the dimmer the objects you can see. The Hubble's primary mirror is 2.4 m in diameter.

**2. Car Headlights and Flashlights**
A light bulb placed at the focal point F of a concave mirror produces a **parallel beam** of light (the reverse of focusing — rays from F reflect as parallel rays). This creates the powerful, directed beam of a flashlight or car headlight.

**3. Dental and Makeup Mirrors**
The dentist's small mirror and the magnifying bathroom mirror both work by placing the object (a tooth, your face) **inside the focal length**. As we saw in Example 2, this produces a **virtual, upright, magnified** image — making small details easier to see.

**4. Solar Furnaces**
Large concave mirrors focus sunlight to the focal point. The enormous concentration of light energy heats the focal point to thousands of degrees Celsius. Solar furnaces in France and other countries reach temperatures exceeding 3000°C without burning any fuel.

**5. Satellite Dish Antennas**
Satellite dishes are parabolic (more on that below) but work on the same principle — parallel incoming radio waves are focused to a point where the receiver sits.

---

### Convex Mirror Applications

**1. Car Wing (Side) Mirrors**
Convex mirrors give a **wide field of view** — you can see a much larger area of traffic than a flat mirror would show. The trade-off is that images appear smaller and closer to the mirror than the actual objects. The warning "Objects in mirror are closer than they appear" is legally required on passenger-side mirrors in many countries.

**2. Security and Surveillance Mirrors**
The round, convex mirrors you see at the corners of shop aisles and parking garages give security staff a wide view of the entire area from one mirror. Their wide field of view is their entire purpose.

**3. Road Safety Mirrors**
At sharp bends in narrow roads, convex mirrors are installed to let drivers see oncoming traffic around the corner before it's visible to the naked eye.

---

### Parabolic Mirrors vs Spherical Mirrors

All our calculations assumed **spherical mirrors** (sections of a sphere). Spherical mirrors have a slight problem called **spherical aberration**: rays hitting the outer edges of a large concave mirror do not focus at exactly the same point as rays hitting near the center. This gives a slightly blurry image.

**Parabolic mirrors** (shaped like a paraboloid of revolution) eliminate this problem completely. Parallel rays hitting any part of a parabolic mirror all focus at exactly one point. For this reason:
- Large research telescopes use **parabolic primary mirrors**
- Flashlights and headlights use **parabolic reflectors**
- Satellite dishes are **parabolic**

For small mirrors or low precision, spherical mirrors are good enough and much cheaper to manufacture. For the highest optical quality, parabolic is the choice.

---

### Interesting Curiosities

**Retroreflectors:** The reflectors on bikes and road signs are designed to send light straight back to its source, regardless of the angle of incidence. They use corner-cube geometry (three mutually perpendicular mirrors). Astronauts left retroreflectors on the Moon during Apollo missions — scientists still bounce laser beams off them today to measure the Earth-Moon distance to millimeter precision.

**The Hubble Disaster:** When the Hubble Space Telescope was launched in 1990, its primary mirror had been ground to the wrong shape — just 2.2 micrometers (about 1/50th the width of a human hair) off. This tiny error caused severe spherical aberration and blurry images. The mirror was later corrected with additional optics installed by astronauts. This incident famously showed how precise mirror shapes must be in high-quality optics.

---

## Summary

- The **Law of Reflection** states that the angle of incidence equals the angle of reflection (θᵢ = θᵣ), both measured from the **normal** to the surface.

- **Specular reflection** occurs at smooth surfaces (mirrors, calm water) — parallel rays reflect in an organized way, forming clear images.

- **Diffuse reflection** occurs at rough surfaces (paper, walls) — rays scatter in all directions, no image forms, but objects are visible from any direction.

- In a **plane mirror**, the image is **virtual, upright, same size, laterally inverted**, and as far behind the mirror as the object is in front.

- **Curved mirrors** are sections of a sphere, characterized by: Pole (P), Centre of Curvature (C), Focal Point (F = midpoint of PC), Radius (R), and Focal Length (f = R/2).

- **Concave mirrors** (converging) have positive focal length. They can form both real and virtual images depending on where the object is placed.

- When the object is beyond F in a concave mirror, the image is **real and inverted**. When the object is inside F, the image is **virtual and upright (magnified)**.

- **Convex mirrors** (diverging) have negative focal length. They always produce **virtual, upright, diminished** images, regardless of object position.

- The **mirror equation** 1/f = 1/v + 1/u relates focal length, object distance, and image distance.

- **Magnification** m = -v/u tells you the size and orientation: negative m = inverted, positive m = upright; |m| > 1 = magnified, |m| < 1 = diminished.

- Sign convention (Real is Positive): real objects and real images have positive distances; virtual images have negative distances; f positive for concave, negative for convex.

- **Concave mirror uses:** telescopes, headlights, makeup mirrors, dental mirrors, solar furnaces.

- **Convex mirror uses:** car wing mirrors, security/surveillance mirrors, road safety mirrors.

- **Parabolic mirrors** have perfect focus with no spherical aberration — used in the best telescopes and precision optics.

---

## Key Equations

**Law of Reflection:**
```
θᵢ = θᵣ
```
(angle of incidence = angle of reflection, both from the normal)

**Focal Length and Radius of Curvature:**
```
f = R / 2
```

**Mirror Equation:**
```
1/f = 1/v + 1/u
```
Or equivalently:
```
1/v = 1/f - 1/u
```
Or:
```
1/u = 1/f - 1/v
```

**Linear Magnification:**
```
m = image height / object height = -v / u
```

**Number of Images Between Two Mirrors at Angle θ:**
```
Number of images = (360° / θ) - 1
```
(valid when 360°/θ is a whole number)

**Sign Convention Summary:**
```
Concave mirror:  f = positive (+)
Convex mirror:   f = negative (-)

Real object:     u = positive (+)
Real image:      v = positive (+)    → image in front of mirror
Virtual image:   v = negative (-)    → image behind mirror

Inverted image:  m = negative (-)
Upright image:   m = positive (+)
Magnified:       |m| > 1
Diminished:      |m| < 1
Same size:       |m| = 1
```

**Useful rearrangements for solving problems:**
```
Finding v:     v = uf / (u - f)
Finding u:     u = vf / (v - f)
Finding f:     f = uv / (u + v)
```

---

*End of Chapter 62: Reflection and Mirrors*

*Next Chapter: Chapter 63 — Refraction and Lenses (when light bends as it crosses from one medium to another, and how lenses form images)*
