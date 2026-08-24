# Chapter 32: Wave Properties

> "The surface of the sea is not the sea. Waves travel across it, but the water stays." — anon

---

## Table of Contents

1. What Is a Wave?
2. The Mexican Wave Analogy
3. Key Wave Properties — Definitions
4. Labelled ASCII Wave Diagram
5. The Wave Equation: v = fλ
6. Worked Example 1 — Finding Wave Speed
7. Worked Example 2 — Finding Wavelength
8. Worked Example 3 — Finding Frequency
9. Wavefronts and Rays
10. The Superposition Principle
11. Interference — Constructive
12. Interference — Destructive
13. Reflection — Fixed End vs. Free End
14. Refraction
15. Diffraction
16. Summary
17. Key Equations

---

## 1. What Is a Wave?

A **wave** is a disturbance that transfers **energy** from one place to another without transferring **matter**.

This is the defining characteristic of waves, and it is counterintuitive at first. Let's unpack it carefully.

When you drop a stone into a still pond:
- Ripples spread outward in circles
- A leaf floating on the water bobs up and down — it does NOT travel outward with the ripples
- Energy from the stone's impact travels outward
- Water molecules move up and down (locally) but do not travel from the center outward

```
    Stone drops:         Ripples spread out:       Leaf bobs:
    
        *                  ~~*~~                    ~~~*~~~
       / \                ~~ * ~~                  ~~[leaf]~~
      /   \              ~~  *  ~~                 ~~  *  ~~
                        ~~   *   ~~               ~~   *   ~~
    
    Energy in:          Energy travels out:        Leaf moves UP/DOWN,
    impact              as a wave                  does not move outward
```

The matter (water) just oscillates locally. The energy moves through the medium.

### What Is a Medium?

A **medium** is the material through which a wave travels. 
- Sound travels through air (medium = air)
- Seismic waves travel through rock (medium = rock)
- Water waves travel through water (medium = water)

Some waves do NOT need a medium — electromagnetic waves (light, radio, X-rays) can travel through the vacuum of space. We'll cover this in Chapter 33.

---

## 2. The Mexican Wave Analogy

The **Mexican wave** (or stadium wave) is a perfect everyday analogy for waves.

Picture a football stadium with 50,000 fans:
1. One fan in section A stands up and raises their arms
2. This triggers fans in section B to stand, then section C, then D...
3. A visible wave travels around the entire stadium

What is traveling? Not the fans — each fan just stands up and sits back down. They return to their seat. The **disturbance** (the "standing-up" action) travels around the stadium. The **information** (or in physics, the **energy**) propagates.

```
    Sections:   A    B    C    D    E    F
    
    Time 1:    [^]  [ ]  [ ]  [ ]  [ ]  [ ]   (A standing)
    Time 2:    [ ]  [^]  [ ]  [ ]  [ ]  [ ]   (B standing)
    Time 3:    [ ]  [ ]  [^]  [ ]  [ ]  [ ]   (C standing)
    Time 4:    [ ]  [ ]  [ ]  [^]  [ ]  [ ]   (D standing)
    
    The WAVE moves right. Each FAN just moves up then down.
    
    [^] = standing    [ ] = seated
```

This is exactly what happens with a mechanical wave in a medium — each particle oscillates about its rest position, but the wave pattern travels forward.

---

## 3. Key Wave Properties — Definitions

### Wavelength (λ — Greek letter "lambda")

**Wavelength** is the distance between two successive identical points on a wave — most easily measured from one crest to the next, or from one trough to the next, or from any point to the next point that is in exactly the same phase of the oscillation.

- Symbol: λ (lambda)
- Unit: meters (m)
- "The length of one complete wave"

### Frequency (f)

**Frequency** is the number of complete wave cycles that pass a fixed point per second.

- Symbol: f
- Unit: hertz (Hz) = cycles per second
- 1 Hz = 1 cycle per second
- A 440 Hz sound means 440 complete wave cycles pass your ear each second

### Period (T)

**Period** is the time for one complete wave cycle to pass a fixed point. It is the reciprocal of frequency.

```
    T = 1 / f        f = 1 / T
```

- Symbol: T
- Unit: seconds (s)
- A 440 Hz wave has period T = 1/440 ≈ 0.00227 s ≈ 2.27 ms

### Amplitude (A)

**Amplitude** is the maximum displacement of the medium from its equilibrium (rest) position.

- Symbol: A
- Unit: meters (m), or appropriate unit for the medium
- Amplitude determines the **energy** carried by the wave — a bigger amplitude means a more energetic wave (energy ∝ A²)
- For sound: amplitude ↔ loudness
- For light: amplitude ↔ brightness

### Wave Speed (v)

**Wave speed** is how fast the wave pattern moves through the medium.

- Symbol: v
- Unit: meters per second (m/s)
- NOT the same as how fast individual particles of the medium move (called particle velocity)

---

## 4. Labelled ASCII Wave Diagram

Study this diagram carefully. It shows all the key properties of a wave at one instant in time:

```
    Displacement
    (meters)
    ^
    |
  A |        *****                    *****
    |      **     **                **     **
    |     *         *              *         *
    |    *           *            *           *
    | ***             *          *             *
    |                  *        *               *
  0 +-------------------*------*-----------------> Distance (meters)
    |                    *    *
    |                     *  *
 -A |                      **
    |
    
    |<--- λ (wavelength) --->|
    
    Example: A = 0.5 m, λ = 4 m
    
    Key points:
    - CREST: maximum positive displacement (top of wave)
    - TROUGH: maximum negative displacement (bottom of wave)
    - AMPLITUDE: distance from center to crest (= A)
    - WAVELENGTH: distance from crest to crest (= λ)
    - The EQUILIBRIUM LINE is at displacement = 0
```

Now showing what a wave looks like over time at a FIXED POINT:

```
    Displacement
    ^
    |
  A |   *       *       *       *
    |  * *     * *     * *     * *
    | *   *   *   *   *   *   *   *
    |*     * *     * *     * *     *
    +-------*-------*-------*-------> Time (seconds)
    |
 -A |
    
    |<----T (period)---->|
    
    At this fixed point, the medium oscillates up and down with period T.
    Frequency f = 1/T tells you how many crests pass per second.
```

---

## 5. The Wave Equation: v = fλ

The most important equation in wave physics connects speed, frequency, and wavelength:

```
    v = f × λ
```

Where:
- v = wave speed (m/s)
- f = frequency (Hz)
- λ = wavelength (m)

### Deriving It From First Principles

Imagine watching waves on a rope. One full wavelength (λ meters) passes a fixed point in exactly one period (T seconds). So:

```
    speed = distance / time
    v = λ / T
    
    Since T = 1/f:
    
    v = λ × (1/T) = λ × f = fλ
```

### The Crucial Insight: v is Fixed by the Medium

In most situations, the wave speed v is determined by the medium (material) the wave is in — not by the source. For example:
- Sound in air at room temperature: always ~340 m/s regardless of pitch
- Light in a vacuum: always exactly 3.00 × 10⁸ m/s

This means if you change the frequency (pitch, for sound), the wavelength must adjust accordingly so that v = fλ stays constant.

- Higher frequency → shorter wavelength (same speed)
- Lower frequency → longer wavelength (same speed)

```
    Same wave speed:
    
    Low frequency, long wavelength:
    ~~~~~~~~~~~~~~~~~~~~~~~~  (few wide ripples per second)
    
    High frequency, short wavelength:
    \/\/\/\/\/\/\/\/\/\/\/\/  (many narrow ripples per second)
    
    Both travel at the same speed in the same medium.
```

---

## 6. Worked Example 1 — Finding Wave Speed

**Problem**: A water wave has a frequency of 2 Hz and a wavelength of 3 m. What is the wave speed?

**Given**:
- f = 2 Hz
- λ = 3 m

**Formula**: v = fλ

**Solution**:
```
    v = f × λ
    v = 2 Hz × 3 m
    v = 6 m/s
```

**Answer**: The wave travels at 6 m/s.

**Sense check**: 2 waves pass per second, each 3 m long — in 1 second, 6 m worth of wave passes. Makes perfect sense!

---

## 7. Worked Example 2 — Finding Wavelength

**Problem**: A sound wave in air travels at 340 m/s with a frequency of 680 Hz. What is its wavelength?

**Given**:
- v = 340 m/s
- f = 680 Hz

**Rearrange**: v = fλ → λ = v/f

**Solution**:
```
    λ = v / f
    λ = 340 m/s ÷ 680 Hz
    λ = 0.5 m
```

**Answer**: The wavelength is 0.5 m (50 cm).

**Sense check**: 680 Hz is a high-pitched sound (roughly the E above middle E on a piano). A wavelength of 50 cm is plausible — human-audible sound wavelengths range from a few centimeters to many meters.

---

## 8. Worked Example 3 — Finding Frequency

**Problem**: A radio wave travels at 3.0 × 10⁸ m/s and has a wavelength of 3 m (this is in the FM radio band, roughly). What is the frequency?

**Given**:
- v = 3.0 × 10⁸ m/s
- λ = 3 m

**Rearrange**: v = fλ → f = v/λ

**Solution**:
```
    f = v / λ
    f = (3.0 × 10⁸ m/s) / 3 m
    f = 1.0 × 10⁸ Hz
    f = 100 MHz
```

**Answer**: 100 MHz — which is exactly the center of the FM radio band (88–108 MHz)!

**Note**: Radio waves and light are both electromagnetic waves traveling at c = 3 × 10⁸ m/s. The only difference is their frequency (and therefore wavelength).

---

## 9. Wavefronts and Rays

When we describe how a wave spreads out, we use two complementary pictures:

### Wavefronts

A **wavefront** is an imaginary surface connecting all points on a wave that are at the same phase (e.g., all at a crest at the same moment). For a point source, wavefronts are circles (in 2D) or spheres (in 3D).

```
    Point source → circular wavefronts:
    
                    )))
                  )     )       (each circle = 
                )    .    )      one wavefront,
                  )     )        one wavelength apart)
                    )))
                    
    Distant source → plane wavefronts (nearly straight lines):
    
    ||||||||  (straight lines, all in phase)
    ||||||||
    ||||||||  ---> direction of travel
    ||||||||
```

### Rays

A **ray** is an arrow drawn perpendicular to the wavefronts, showing the direction the wave is traveling. Rays are useful for tracking where the wave goes (especially in optics with light).

```
    >>>     ray (arrow showing direction)
    
    Rays and wavefronts are always perpendicular:
    
    wavefronts:  |  |  |  |  |
    rays:        →  →  →  →  →
```

---

## 10. The Superposition Principle

When two waves meet in the same medium, they **add together** while they overlap and then pass through each other unchanged. This is the **principle of superposition**:

```
    Resultant displacement = sum of individual displacements

    If wave 1 displaces medium by +3 cm
    And wave 2 displaces medium by +2 cm at the same point
    Then: resultant = +3 + (+2) = +5 cm

    If wave 1 displaces medium by +3 cm
    And wave 2 displaces medium by -3 cm at the same point
    Then: resultant = +3 + (-3) = 0 cm
```

After they pass through each other, each wave continues as if nothing happened — waves pass through each other without distortion.

---

## 11. Interference — Constructive

**Constructive interference** occurs when two waves arrive at the same point in phase — both at crests at the same time, or both at troughs at the same time. Their displacements add together.

```
    Wave 1:    +  ++  +
              /\/\/\/\/\
    
    Wave 2:    +  ++  +
              /\/\/\/\/\
    
    Result:   + ++ +++ ++ +
             /\/\/\/\/\/\/\    (DOUBLE the amplitude!)
    
    Crest meets crest → double amplitude
    Trough meets trough → double depth
```

The resulting amplitude is the sum of the two individual amplitudes. For two identical waves completely in phase:

```
    Resultant amplitude = A₁ + A₂
    
    For two equal waves: A_result = 2A
```

### Real-World Constructive Interference

- **Noise-cancelling headphones** work by understanding interference — they actually destructively interfere with unwanted noise, but constructive interference explains why combining two identical signals makes them louder.
- **Loud spots** in a concert hall where sound is enhanced because waves from multiple speakers arrive in phase.

---

## 12. Interference — Destructive

**Destructive interference** occurs when two waves arrive at the same point exactly out of phase — when one is at a crest and the other is at a trough. Their displacements cancel.

```
    Wave 1:    +  ++  +
              /\/\/\/\/\
    
    Wave 2:   - -- -- -- (inverted — shifted by half a wavelength)
              \/\/\/\/\/
    
    Result:   ----------   (ZERO — they cancel!)
    
    Crest meets trough → cancellation
```

For two identical waves completely out of phase:
```
    Resultant amplitude = A₁ - A₂ = A - A = 0
```

### Real-World Destructive Interference

- **Noise-cancelling headphones**: microphones detect ambient noise, circuits create an inverted copy, headphone speakers play both — they destructively interfere, silencing the noise.
- **Dead spots** in concert halls or nightclubs where sound waves from different directions cancel.
- **Anti-reflection coatings** on camera lenses: a thin layer causes destructive interference for reflected light, reducing glare.

### Partial Interference

When waves are neither perfectly in phase nor perfectly out of phase, partial constructive or destructive interference occurs, producing intermediate amplitudes.

```
    Constructive:        Partial:         Destructive:
    
     /\/\                 /\               ----
    /    \     +         /  \      +
    
    = /\/\/\/\           = /\        =  (nothing)
    
    Double amplitude     Some boost      Cancellation
```

---

## 13. Reflection — Fixed End vs. Free End

When a wave on a rope reaches the end of the rope, it **reflects** back. But there are two cases:

### Reflection at a Fixed End

If the end of the rope is tied to a wall (fixed), the wave reflects and is **inverted** — a crest becomes a trough.

```
    Incoming wave:           Reflected wave:
    
        /\/\               \/\/   (inverted!)
    ---/    \---wall    ---    ---
```

**Why inverted?** The wall exerts an equal and opposite (Newton's Third Law) force back on the rope. This creates an inverted pulse.

### Reflection at a Free End

If the end of the rope is free to move (e.g., tied to a ring that slides on a rod), the wave reflects **without inversion** — a crest stays a crest.

```
    Incoming wave:           Reflected wave:
    
        /\/\                     /\/\   (same orientation!)
    ---/    \===end          ===      ---
```

**Why not inverted?** The free end can move freely — it overshoots upward, creating a reflected wave in the same direction.

### Why This Matters

The difference between fixed-end and free-end reflection determines whether standing waves form with a **node** (zero displacement) or **antinode** (maximum displacement) at the endpoint. This is crucial for understanding musical instruments (Chapter 36).

---

## 14. Refraction

**Refraction** is the change in direction (bending) of a wave when it passes from one medium to another where it travels at a different speed.

The key: wave frequency stays constant when crossing a boundary. But if the speed changes, the wavelength must change (v = fλ). This change in wavelength means the wavefronts bunch up or spread out, causing the wave to bend.

```
    Fast medium (e.g., deep water):    Slow medium (e.g., shallow water):
    
    |  |  |  |  |                      ||||||||
    |  |  |  |  |                      ||||||||
    |  |  |  |  |   ----boundary-->    ||||||||
    
    Wavefronts far apart               Wavefronts crowd together
    (long wavelength, fast)            (short wavelength, slow)
    
    Wave slows down and bends TOWARD the boundary normal
```

### Snell's Law (preview)

The bending follows Snell's Law: the wave bends toward the normal (perpendicular to boundary) when slowing down, and away from the normal when speeding up.

### Everyday Examples

- Light bending when entering glass or water (a straw looks bent in a glass of water)
- Ocean waves bending to face the beach as they approach shallow water (all beaches face their waves — it's refraction!)
- Sound bending over long distances due to temperature gradients in the atmosphere

---

## 15. Diffraction

**Diffraction** is the bending and spreading of waves around obstacles or through gaps. All waves diffract.

```
    Wave hits a wall with a small gap:
    
    ===========|   |===========   (wall with gap)
               | gap |
    
    Incoming:   Spreads out after gap:
    
    |||||||     |||  )))
    |||||||     ||| )   )
    |||||||     |||)     )
    |||||||     |||       )
    
    The wave spreads into the shadow region!
```

### Diffraction is Most Noticeable When:

```
    Wavelength ≈ Gap size
```

- If the gap is much larger than the wavelength → little diffraction, mostly straight through
- If the gap is about the same size as the wavelength → significant diffraction, spreads widely
- If the gap is much smaller than the wavelength → spreads almost as a new circular wave

### Why Can You Hear Around Corners But Not See?

- Sound: wavelengths 1 cm to 17 m. Doorways (~1 m wide) are comparable to sound wavelengths → sound diffracts around corners.
- Light: wavelengths ~400–700 nm. Doorways are billions of times bigger than light wavelengths → light does not noticeably diffract around doors, so shadows are sharp.

```
    Sound diffraction around corner:     Light — no diffraction:
    
    WALL                                 WALL
    ||||                                 ||||
    ||||  sound                          |||| SHARP SHADOW
    ||||  source                         ||||
    ||||   *                             ||||   *
    ||||     ))                          ||||    (no bending)
    ||||       ))  ← bends around        ||||
    ||||         ))   corner             ||||
    
    You can hear around a corner         You cannot see around a corner
```

---

## 16. Summary

- A **wave** transfers energy without transferring matter. Each particle oscillates locally while the wave pattern travels.

- The **Mexican wave** analogy: fans stand and sit (local oscillation), but the wave pattern travels around the stadium.

- **Key properties**: wavelength (λ = distance between successive identical points), frequency (f = cycles per second), period (T = 1/f), amplitude (A = maximum displacement), wave speed (v).

- **The wave equation**: v = fλ. Speed, frequency, and wavelength are linked. Speed is set by the medium; changing frequency changes wavelength proportionally.

- **Wavefronts** are surfaces of equal phase; **rays** are perpendicular arrows showing direction of travel.

- **Superposition principle**: displacements add when waves overlap.

- **Constructive interference**: waves in phase → amplitudes add → louder/brighter.

- **Destructive interference**: waves out of phase → amplitudes cancel → quiet/dark.

- **Reflection at fixed end**: inverted. **Reflection at free end**: not inverted.

- **Refraction**: wave bends when speed changes crossing a boundary. Frequency stays constant, wavelength changes.

- **Diffraction**: waves spread around obstacles or through gaps. Most significant when wavelength ≈ gap size.

---

## 17. Key Equations

```
Wave Equation:
    v = f × λ
    (v in m/s, f in Hz, λ in m)

Rearranged forms:
    f = v / λ
    λ = v / f

Period-Frequency Relationship:
    T = 1 / f        f = 1 / T

Energy proportionality:
    E ∝ A²
    (energy proportional to amplitude squared)

Superposition:
    y_total = y₁ + y₂
    (displacements add algebraically)

Constructive interference (two equal waves, in phase):
    A_result = A₁ + A₂

Destructive interference (two equal waves, out of phase):
    A_result = |A₁ - A₂|  (= 0 for equal amplitudes)
```
