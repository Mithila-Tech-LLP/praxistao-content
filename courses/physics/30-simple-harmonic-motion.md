# Chapter 30: Simple Harmonic Motion

> **"SHM is everywhere — pendulums, guitar strings, molecules, atoms, and even the tides. Understanding it unlocks the deepest patterns in physics."**

---

## Table of Contents

- [30.1 What is Oscillation?](#301-what-is-oscillation)
- [30.2 Defining Simple Harmonic Motion](#302-defining-simple-harmonic-motion)
- [30.3 The Mathematics of SHM](#303-the-mathematics-of-shm)
- [30.4 Energy in SHM](#304-energy-in-shm)
- [30.5 The Spring-Mass System](#305-the-spring-mass-system)
- [30.6 The Simple Pendulum](#306-the-simple-pendulum)
- [30.7 Graphs in SHM](#307-graphs-in-shm)
- [30.8 Damped Oscillations](#308-damped-oscillations)
- [30.9 Resonance](#309-resonance)
- [30.10 Real-World SHM](#3010-real-world-shm)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 30.1 What is Oscillation?

An **oscillation** is a repeated back-and-forth motion about an equilibrium position.

```
EQUILIBRIUM POSITION
        |
        v
        *
       /|\
      / | \
     /  |  \
    /   |   \
LEFT <--+--> RIGHT
  extreme    extreme

The object repeatedly moves left → center → right → center → left...
```

Key terms:

- **Equilibrium position**: where the object would rest if undisturbed
- **Displacement (x)**: how far the object is from equilibrium (can be positive or negative)
- **Amplitude (A)**: the maximum displacement (always positive)
- **Period (T)**: time for one complete oscillation (seconds)
- **Frequency (f)**: oscillations per second (Hz); f = 1/T
- **Angular frequency (ω)**: ω = 2πf (rad/s)

---

## 30.2 Defining Simple Harmonic Motion

**Simple Harmonic Motion (SHM)** is a special type of oscillation where:

> The acceleration is always directed toward the equilibrium position, and its magnitude is proportional to the displacement from equilibrium.

In math:

```
a = -(ω²) × x

The minus sign shows that acceleration opposes displacement:
  - When x is positive (right), a is negative (pointing left)
  - When x is negative (left), a is positive (pointing right)
  - This is the restoring force that keeps bringing the object back
```

This defining condition is the key test for SHM: **a ∝ -x**.

### The Restoring Force

By Newton's 2nd Law:

```
F = m × a = m × (-ω²x) = -mω²x

So: F ∝ -x

Force always points back toward equilibrium — it "restores" the object.
```

For a spring: F = -kx (Hooke's Law) — this IS the restoring force for a spring-mass system.

---

## 30.3 The Mathematics of SHM

If at t = 0 the object is at maximum displacement (x = A), displacement varies as:

```
x(t) = A × cos(ωt)

where:
  x = displacement at time t
  A = amplitude
  ω = angular frequency = 2πf = 2π/T
  t = time
```

If at t = 0 the object is at equilibrium moving in positive direction:

```
x(t) = A × sin(ωt)
```

The general solution uses either, depending on initial conditions.

### Velocity in SHM

Taking the derivative of x(t) = A cos(ωt):

```
v(t) = -Aω × sin(ωt)

Maximum speed = Aω  (at equilibrium, x = 0)
Zero speed = 0 at    (at extremes, x = ±A)
```

### Acceleration in SHM

Differentiating velocity:

```
a(t) = -Aω² × cos(ωt) = -ω² × x

Maximum acceleration = Aω²  (at extremes, x = ±A)
Zero acceleration at equilibrium (x = 0)
```

### Summary of Phases

```
At equilibrium (x = 0):
  Displacement = 0
  Speed = maximum (Aω)
  Acceleration = 0
  
At extreme (x = A or x = -A):
  Displacement = maximum (A)
  Speed = 0
  Acceleration = maximum (Aω²)
```

---

## 30.4 Energy in SHM

In SHM, energy continuously swaps between kinetic and potential forms:

```
POSITION:      LEFT    MIDDLE    RIGHT
               extreme  (equil)  extreme

Displacement:   -A        0        +A
Speed:           0       max        0
KE:              0       max        0
PE:             max        0       max
Total E:        max       max       max

KE + PE = constant (Total mechanical energy)
```

### Energy Equations

```
Kinetic energy:
  KE = (1/2)mv²  = (1/2)m(Aω)²sin²(ωt)

Potential energy (for spring):
  PE = (1/2)kx² = (1/2)mω²x²  = (1/2)m(Aω)²cos²(ωt)

Total energy:
  E = (1/2)mA²ω² = (1/2)kA²  (constant)
```

Using sin²+cos² = 1 confirms that KE + PE = constant.

### Velocity at Any Position

From energy conservation:

```
(1/2)mv² + (1/2)kx² = (1/2)kA²

v² = (k/m)(A² - x²) = ω²(A² - x²)

v = ω × sqrt(A² - x²)
```

This gives speed at any displacement x.

---

## 30.5 The Spring-Mass System

A mass m attached to a spring of stiffness k undergoes SHM.

```
SPRING-MASS SYSTEM:

|||||  spring  [mass m]
            |
            v  displacement x

Restoring force: F = -kx
Newton's 2nd: ma = -kx
So: a = -(k/m)x

Comparing to SHM definition a = -ω²x:
  ω² = k/m
  ω = sqrt(k/m)
```

Therefore the period of a spring-mass system:

```
T = 2π / ω = 2π × sqrt(m/k)

f = (1/2π) × sqrt(k/m)
```

Key observations:
- Larger mass → longer period (harder to accelerate)
- Stiffer spring → shorter period (stronger restoring force)
- Amplitude does NOT affect period!

### Worked Example 30.1

A mass of 0.5 kg is attached to a spring (k = 200 N/m). Find the period and frequency of oscillation.

**Solution:**

ω = sqrt(k/m) = sqrt(200/0.5) = sqrt(400) = 20 rad/s

T = 2π/ω = 2π/20 ≈ **0.314 s**

f = 1/T ≈ **3.18 Hz**

### Worked Example 30.2

The same spring-mass system (T = 0.314 s) is displaced by 0.1 m and released.

Find: (a) maximum speed, (b) maximum acceleration, (c) speed when x = 0.06 m.

**Solution:**

ω = 20 rad/s, A = 0.1 m

(a) v_max = Aω = 0.1 × 20 = **2 m/s**

(b) a_max = Aω² = 0.1 × 400 = **40 m/s²**

(c) v = ω × sqrt(A² - x²) = 20 × sqrt(0.01 - 0.0036) = 20 × sqrt(0.0064) = 20 × 0.08 = **1.6 m/s**

---

## 30.6 The Simple Pendulum

A simple pendulum (bob on a string of length L) undergoes SHM for small angles (less than about 15°).

```
         pivot
          |
          |  L (length)
          |
         [m]
          
For small angle θ:
  Restoring force ≈ mg × θ  (where θ in radians)
  
  This gives: ω = sqrt(g/L)
```

Period of a simple pendulum:

```
T = 2π × sqrt(L/g)

f = (1/2π) × sqrt(g/L)
```

Key observations:
- Longer pendulum → longer period
- Stronger gravity → shorter period
- Mass does NOT affect period (same as free fall!)
- Amplitude does NOT affect period (for small angles)

### Worked Example 30.3

Find the length of a pendulum with period T = 2 s (a "seconds pendulum" — swings once per second).

**Solution:**

T = 2π × sqrt(L/g)

2 = 2π × sqrt(L/9.8)

1/π = sqrt(L/9.8)

1/π² = L/9.8

L = 9.8/π² ≈ **0.993 m ≈ 1 m**

This is why grandfather clocks have pendulums about 1 m long!

### Worked Example 30.4

A student times 20 oscillations of a pendulum and gets 36 s.

(a) Find T and f.
(b) Find the length of the pendulum.

**Solution:**

(a) T = 36/20 = 1.8 s, f = 1/1.8 ≈ **0.556 Hz**

(b) T = 2π × sqrt(L/g)
    T² = 4π² × L/g
    L = gT²/(4π²) = 9.8 × 3.24 / (4π²) ≈ **0.806 m**

---

## 30.7 Graphs in SHM

### For a system starting at x = A (released from rest at maximum displacement):

```
DISPLACEMENT vs TIME:  x = A cos(ωt)

x
 A|  *       *       *
  | * *     * *     * *
  |*   *   *   *   *
  0|    *   *       *
 -A|     * *         *
  +----------------------> t
    T/4  T/2  3T/4  T
```

```
VELOCITY vs TIME:  v = -Aω sin(ωt)

v
Aω|      *             *
  |     * *           * *
  |    *   *         *   *
  0|---*-----*------*-----*--> t
  |          *     *
-Aω|          * * *

Velocity peaks when displacement = 0 (at equilibrium).
```

```
ACCELERATION vs TIME:  a = -Aω² cos(ωt)

a
    |  *       *       *  <- at extremes (same shape as -x graph)
    |
  0 |----*--------*-------> t
    |        *
    |        
    
a is always 180° out of phase with x.
```

### Phase Relationships

- Velocity is 90° ahead of displacement
- Acceleration is 180° ahead of displacement (or 180° behind)
- Maximum v when x = 0; maximum a when x = ±A

---

## 30.8 Damped Oscillations

In real oscillators, friction and air resistance remove energy. The amplitude gradually decreases.

```
DAMPED OSCILLATION:

x
A|  *
  | * *
  |*   *       *
  |     *     * *
  |      *   *   *     *
  |       * *     *   * *
  |              *  * *  *
  +--------------------------> t

Envelope (dashed): A decreases exponentially over time
```

### Types of Damping

```
LIGHT DAMPING (underdamping):
  Object oscillates but amplitude shrinks slowly.
  Used in: clocks, musical instruments (we want sustained oscillation)

CRITICAL DAMPING:
  Object returns to equilibrium as fast as possible WITHOUT oscillating.
  Used in: car shock absorbers, door closers

HEAVY DAMPING (overdamping):
  Object returns to equilibrium slowly WITHOUT oscillating.
  
CRITICAL vs HEAVY:
  Both don't oscillate, but critical is fastest to return to rest.
  Car shock absorbers use critical damping — return to level fast without bouncing.
```

---

## 30.9 Resonance

Every oscillator has a **natural frequency** (f₀) at which it oscillates when disturbed.

**Resonance** occurs when a periodic driving force has the same frequency as the natural frequency. The amplitude grows very large.

```
RESONANCE GRAPH:

Amplitude
    |
    |           *   <- resonance peak
    |          * *
    |         *   *
    |       *       *
    |   *               *    *
    +--------------------------> driving frequency f
                f0
```

### Real-World Examples

```
TACOMA NARROWS BRIDGE (1940):
  Wind created oscillations near bridge's natural frequency.
  Bridge resonated — amplitude grew until it collapsed.
  
WINE GLASS:
  Sung at the right frequency → resonates → shatters!
  
MICROWAVE OVEN:
  Microwaves at the resonant frequency of water molecules.
  Water absorbs energy, heats up.
  
MRI MACHINE:
  Radio waves at resonant frequency of hydrogen atoms.
  Atoms absorb and re-emit energy → used for imaging.
```

---

## 30.10 Real-World SHM

### Molecular Vibrations

Atoms in a molecule vibrate in SHM. The frequency depends on bond stiffness and atomic mass. This is how infrared spectroscopy identifies molecules.

### Seismometers

Detect earthquakes by measuring oscillations in the ground. The pendulum inside oscillates at a known frequency; deviations reveal ground motion.

### Quartz Clocks

Quartz crystals vibrate at very precise frequencies (e.g., 32,768 Hz) when an electric field is applied. This vibration drives the clock mechanism.

### Sound Production

Guitar strings, piano strings, and speaker cones all undergo SHM. The frequency of vibration determines the pitch (frequency) of the sound.

```
GUITAR STRING SHM:
  
  ~~~~~~~~~~~~~~~~~~  vibrating string
  
  Amplitude → loudness
  Frequency → pitch
  
  f = (1/2L) × sqrt(T/μ)
  where L = length, T = tension, μ = mass per unit length
```

---

## Summary

- **SHM**: oscillation where a ∝ -x (acceleration proportional and opposite to displacement)
- **Defining equation**: a = -ω²x
- **Displacement**: x = A cos(ωt) or A sin(ωt)
- **Velocity**: v = ±ω√(A² - x²); maximum at equilibrium (Aω)
- **Acceleration**: a = -ω²x; maximum at extremes (Aω²)
- **Spring-mass**: ω = √(k/m), T = 2π√(m/k)
- **Pendulum**: ω = √(g/L), T = 2π√(L/g)
- Both spring and pendulum: period independent of amplitude
- **Energy**: total E = ½kA²; swaps between KE and PE
- **Damping**: reduces amplitude; critical damping returns fastest without oscillation
- **Resonance**: maximum amplitude when driving frequency = natural frequency

---

## Key Equations

```
Defining SHM:
  a = -ω²x   (acceleration ∝ -displacement)

Displacement:
  x = A cos(ωt)   (starting from maximum)

Velocity:
  v = Aω sin(ωt)
  v = ω × sqrt(A² - x²)   (at any position)
  v_max = Aω

Acceleration:
  a = Aω² cos(ωt)
  a_max = Aω²

Spring-mass system:
  ω = sqrt(k/m)
  T = 2π × sqrt(m/k)

Simple pendulum (small angles):
  ω = sqrt(g/L)
  T = 2π × sqrt(L/g)

Angular frequency:
  ω = 2πf = 2π/T

Total energy in SHM:
  E = (1/2) × k × A²
  E = (1/2) × m × A² × ω²
```
