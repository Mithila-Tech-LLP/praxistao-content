# Chapter 36: Standing Waves and Resonance

> **"Standing waves are where two waves collide and create something completely different — a pattern that looks frozen in place, yet carries energy everywhere."**

---

## Table of Contents

- [36.1 Superposition of Waves](#361-superposition-of-waves)
- [36.2 How Standing Waves Form](#362-how-standing-waves-form)
- [36.3 Nodes and Antinodes](#363-nodes-and-antinodes)
- [36.4 Standing Waves on a String](#364-standing-waves-on-a-string)
- [36.5 Harmonics and Overtones](#365-harmonics-and-overtones)
- [36.6 Standing Waves in Pipes](#366-standing-waves-in-pipes)
- [36.7 Resonance](#367-resonance)
- [36.8 Beats](#368-beats)
- [36.9 Applications of Standing Waves](#369-applications-of-standing-waves)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 36.1 Superposition of Waves

When two or more waves overlap in the same medium at the same time, the **principle of superposition** applies:

> The resultant displacement at any point is the algebraic sum of the displacements of the individual waves.

```
CONSTRUCTIVE INTERFERENCE:
  Two waves in phase (crests align with crests):
  
  Wave 1:    /\   /\   /\
  Wave 2:    /\   /\   /\
  Result:   /  \ /  \ /  \   <- double amplitude
           /    /    /
  
  Amplitude = A1 + A2

DESTRUCTIVE INTERFERENCE:
  Two waves out of phase (crests align with troughs):
  
  Wave 1:    /\   /\   /\
  Wave 2:    \/   \/   \/
  Result:    —————————————   <- flat line (zero)
  
  If A1 = A2: complete cancellation
```

---

## 36.2 How Standing Waves Form

A **standing wave** forms when two waves of equal frequency and amplitude travel in opposite directions through the same medium.

This happens most commonly when a wave reflects off a fixed boundary and interferes with the incoming wave.

```
FORMING A STANDING WAVE:

Step 1: Send a wave right along a string fixed at both ends
  ~~~~~~~~~~~~~~~~~~~~>

Step 2: Wave hits right wall, reflects back
  <~~~~~~~~~~~~~~~~~~~~

Step 3: Right-traveling + left-traveling waves overlap:
  ~~~~ SUPERPOSITION ~~~~

Step 4: Result — a STANDING WAVE:

  time 1:    /\/\/\/\/\
  time 2:    ----------   <- everything flat
  time 3:    \/\/\/\/\/
  time 4:    ----------
  
  The pattern appears to "stand still" — no energy propagates left or right
  (though energy is stored in the wave).
```

Why does it look stationary? The interference pattern creates permanent positions of zero displacement (nodes) and maximum displacement (antinodes) that don't move.

---

## 36.3 Nodes and Antinodes

```
STANDING WAVE PATTERN:

     A       A       A       A
      |       |       |       |
  ----*---*---*---*---*---*---*----  (string)
      N       N       N       N
  
  N = NODE: point that never moves (always zero displacement)
  A = ANTINODE: point with maximum amplitude (moves most)
  
  Adjacent N to N (or A to A) = half a wavelength (λ/2)
```

**Node**: a point where the two waves always cancel (destructive interference). Displacement is always zero.

**Antinode**: a point where the two waves always reinforce. Displacement oscillates between maximum positive and maximum negative.

```
Distance between adjacent nodes = λ/2
Distance between node and adjacent antinode = λ/4
```

---

## 36.4 Standing Waves on a String

A string fixed at both ends must have **nodes at both ends** (since the string cannot move at the fixed points).

Only certain wavelengths can fit. These are the **resonant modes** or **harmonics**.

```
FUNDAMENTAL MODE (1st harmonic, n=1):
  One antinode in the middle.
  
  N---------A---------N
  |<------- L ------->|
  
  Half wavelength fits:  L = λ1/2
  So: λ1 = 2L
  f1 = v/λ1 = v/(2L)

2nd HARMONIC (n=2):
  Two antinodes.
  
  N----A----N----A----N
  |<------- L ------->|
  
  One full wavelength: L = λ2
  So: λ2 = L
  f2 = v/L = 2f1

3rd HARMONIC (n=3):
  
  N--A--N--A--N--A--N
  |<------- L ------->|
  
  Three half-wavelengths: L = 3λ3/2
  So: λ3 = 2L/3
  f3 = 3v/(2L) = 3f1
```

### General Rule for String

```
For nth harmonic:
  λ_n = 2L/n
  f_n = n × v/(2L) = n × f1

where:
  n = 1, 2, 3, 4, ...
  L = length of string
  v = wave speed on string
  f1 = fundamental frequency
```

Wave speed on a string depends on tension T and linear mass density μ:

```
v = sqrt(T/μ)

where μ = mass per unit length (kg/m)
```

### Worked Example 36.1

A guitar string is 0.65 m long. The wave speed on it is 520 m/s.

Find the fundamental frequency and the first three harmonics.

**Solution:**

f1 = v/(2L) = 520/(2 × 0.65) = 520/1.3 = **400 Hz**

f2 = 2f1 = **800 Hz**
f3 = 3f1 = **1200 Hz**
f4 = 4f1 = **1600 Hz**

---

## 36.5 Harmonics and Overtones

Every string or pipe can vibrate at multiple frequencies simultaneously. These are called **harmonics** or **overtones**.

```
MUSICAL TERMS vs PHYSICS TERMS:

Physics:     1st harmonic    2nd harmonic    3rd harmonic
             (fundamental)
Music:       fundamental    1st overtone    2nd overtone
             
Note:  1st overtone = 2nd harmonic (one higher than fundamental)
```

The **timbre** (quality/color) of a musical note depends on which harmonics are present and their relative amplitudes.

```
PURE TONE: only f1

     /\   /\   /\
    /  \ /  \ /  \

RICH SOUND: f1 + f2 + f3 + ...

    complicated waveform
    /\/\/\/\/\/\/\/\/
    
Musical instruments produce characteristic blends of harmonics.
Flute: mostly fundamental (pure)
Violin: rich in overtones (complex timbre)
```

---

## 36.6 Standing Waves in Pipes

Sound waves in pipes form standing waves too — but now instead of transverse waves on a string, we have longitudinal pressure waves in air.

### Pipes Open at Both Ends

At an open end, there must be an **antinode** (the air can move freely).
At a closed end, there must be a **node** (air cannot move past the wall).

**Both ends open (flute-like):**

```
A---------A          (fundamental: L = λ/2)
A----N----A----N----A  (2nd harmonic: L = λ)
A--N--A--N--A--N--A    (3rd harmonic: L = 3λ/2)
```

All harmonics present: f_n = n × v/(2L) for n = 1, 2, 3, ...

**One end closed, one open (clarinet-like):**

```
N---------A          (fundamental: L = λ/4)
N----A----N----A     (3rd harmonic: L = 3λ/4)
N--A--N--A--N--A     (5th harmonic: L = 5λ/4)
```

Only ODD harmonics present: f_n = n × v/(4L) for n = 1, 3, 5, ...

```
Summary:
  Pipe open at both ends:   all harmonics (n=1,2,3,4,...)
  Pipe closed at one end:   odd harmonics only (n=1,3,5,...)
  
This is why a clarinet (closed at reed end, open at bell) sounds different 
from a flute (open at both ends) — different harmonic content!
```

### Worked Example 36.2

A pipe closed at one end has a fundamental frequency of 150 Hz.

(a) What is its length? (v_sound = 340 m/s)
(b) What are the next two resonant frequencies?

**Solution:**

(a) f1 = v/(4L) → L = v/(4f1) = 340/(4 × 150) = **0.567 m**

(b) Only odd harmonics:
    f3 = 3f1 = 3 × 150 = **450 Hz**
    f5 = 5f1 = 5 × 150 = **750 Hz**

---

## 36.7 Resonance

**Resonance** occurs when a system is driven at its natural (resonant) frequency, causing large-amplitude oscillations.

For a string or pipe, resonance occurs when the driving frequency matches one of the harmonic frequencies.

```
RESONANCE CONDITION:
  Driving frequency = natural frequency
  
  Standing wave builds up → amplitude grows large

OFF RESONANCE:
  Driving frequency ≠ natural frequency
  
  No sustained pattern → small amplitude, wastes energy
```

### Forced Vibration

When you hold one end of a string and shake it, you're "driving" the string. Most frequencies produce a messy, small-amplitude response. But at resonant frequencies, a clear standing wave pattern emerges.

### Resonance Examples

```
WINE GLASS:
  Each glass has a natural frequency.
  Sing/play that exact note → glass resonates → may shatter.
  
TACOMA NARROWS BRIDGE:
  Wind-driven oscillations matched bridge's natural frequency.
  Resonance → growing amplitude → collapse (1940).
  
MUSICAL INSTRUMENTS:
  Strings/pipes resonate at specific harmonics.
  Plucking/blowing drives the resonant modes.
  
RADIO TUNING:
  LC circuit tuned to resonate at specific frequency.
  Selects only that station from all radio waves present.
```

---

## 36.8 Beats

**Beats** occur when two waves of slightly different frequencies interfere. The result is a wave whose amplitude oscillates (pulses) at the **beat frequency**:

```
f_beat = |f1 - f2|

where f1 and f2 are the two slightly different frequencies.
```

```
BEAT FORMATION:
  
  Wave 1 (f1 = 100 Hz):  /\/\/\/\/\/\/\/\/\/
  Wave 2 (f2 = 102 Hz):  /\/\/\/\/\/\/\/\/\/  (slightly faster)
  
  Sometimes in phase (constructive) — LOUD
  Sometimes out of phase (destructive) — QUIET
  
  Combined:  /\/\  ——  /\/\  ——  /\/\
             LOUD quiet LOUD quiet LOUD
  
  Beats per second = |100 - 102| = 2 beats/second
```

### Applications of Beats

**Tuning instruments**: Musicians tune by listening for beats. As two strings approach the same frequency, beats slow down and disappear. Zero beats = perfectly in tune.

**Stethoscope and medical devices**: Some diagnostic tools use beat frequencies.

**Heterodyne radio receivers**: Two frequencies mixed to produce a beat at a convenient intermediate frequency for amplification.

---

## 36.9 Applications of Standing Waves

### Musical Instruments

All string and wind instruments work by creating standing waves:

```
VIOLIN:
  Strings fixed at both ends.
  Bow sets string vibrating.
  Fundamental + harmonics create rich tone.
  Pitch controlled by shortening string (raising f).
  
FLUTE:
  Tube open at both ends.
  Breath creates pressure wave.
  Holes change effective length → change pitch.
  
ORGAN PIPE:
  Some pipes open both ends, some closed at one end.
  Different lengths produce different fundamental frequencies.
```

### Laser Cavity

A laser is basically a standing wave cavity for light. The mirrors at each end reflect the light, and only frequencies that form standing waves (match the cavity length) are amplified.

### Microwave Ovens (Hot Spots)

Microwave ovens create standing wave patterns inside. This creates **nodes** (cold spots) and **antinodes** (hot spots). That's why food heats unevenly — and why the turntable rotates: to average out the standing wave pattern.

### Electron Orbitals

The shapes of electron orbitals in atoms are standing wave patterns of the electron's quantum wavefunction. The allowed energy levels correspond to allowed standing wave modes.

---

## Summary

- **Superposition**: waves add algebraically; constructive (+) or destructive (-) interference
- **Standing waves**: formed by two identical waves traveling in opposite directions
- **Node**: fixed point of zero displacement; **antinode**: maximum displacement
- Adjacent nodes separated by λ/2
- **String (both ends fixed)**: all harmonics; f_n = nv/(2L) for n = 1, 2, 3, ...
- **Pipe (both ends open)**: all harmonics; f_n = nv/(2L)
- **Pipe (one end closed)**: odd harmonics only; f_n = nv/(4L) for n = 1, 3, 5, ...
- **Resonance**: large amplitude when driving frequency = natural frequency
- **Beats**: |f1 - f2| — amplitude oscillation when two close frequencies interfere
- Applications: music, lasers, microwave ovens, radio tuning, atom structure

---

## Key Equations

```
Standing waves on string (both ends fixed):
  Harmonics: f_n = n × v / (2L)   n = 1, 2, 3, ...
  Wavelengths: λ_n = 2L / n
  v_string = sqrt(T/μ)

Standing waves in pipe (both ends open):
  f_n = n × v / (2L)   n = 1, 2, 3, ...

Standing waves in pipe (one end closed):
  f_n = n × v / (4L)   n = 1, 3, 5, ...

Beat frequency:
  f_beat = |f1 - f2|

Node spacing:
  Δx (node to node) = λ/2

Node to antinode spacing:
  λ/4
```
