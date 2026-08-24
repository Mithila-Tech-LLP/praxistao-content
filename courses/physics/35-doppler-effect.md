# Chapter 35: The Doppler Effect

> **"The Doppler effect is why ambulance sirens seem to change pitch as they pass you — and why astronomers know the universe is expanding."**

---

## Table of Contents

- [35.1 What is the Doppler Effect?](#351-what-is-the-doppler-effect)
- [35.2 The Physics Behind It](#352-the-physics-behind-it)
- [35.3 The Doppler Formula](#353-the-doppler-formula)
- [35.4 Moving Source, Stationary Observer](#354-moving-source-stationary-observer)
- [35.5 Stationary Source, Moving Observer](#355-stationary-source-moving-observer)
- [35.6 Both Moving](#356-both-moving)
- [35.7 The Sonic Boom](#357-the-sonic-boom)
- [35.8 Doppler Effect with Light (Redshift)](#358-doppler-effect-with-light-redshift)
- [35.9 Real-World Applications](#359-real-world-applications)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 35.1 What is the Doppler Effect?

Have you ever stood near a road and listened to a car pass by with its horn honking?

```
CAR APPROACHING:
  [horn: hoooooOOOOOONK!] ---> (you)
  
CAR RECEDING:
  (you) <--- [horn: HOOOOOoooooonk...]
  
The pitch goes HIGH as it approaches, then drops LOW as it moves away.
Yet the driver hears a constant pitch the whole time!
```

This is the **Doppler effect**: the apparent change in the frequency (pitch) of a wave due to relative motion between the source and the observer.

It applies to all waves — sound, light, radio waves, water waves.

---

## 35.2 The Physics Behind It

Think about a source emitting one wave crest per second (f = 1 Hz):

```
STATIONARY SOURCE:
  
  ... [crest] [crest] [crest] [source] [crest] [crest] [crest] ...
  
  Equal spacing in all directions.
  Observer at left and right both hear f = 1 Hz.

SOURCE MOVING RIGHT:
  
  [crest][crest][source>][crest][crest]
  
  <- left observer             right observer ->
  
  Wave crests PILE UP ahead (right) of the moving source:
    → shorter wavelength ahead → higher frequency
    
  Wave crests STRETCH OUT behind (left):
    → longer wavelength behind → lower frequency
```

Key insight: **the source emits at the same frequency, but motion changes how crests pile up or spread out for the observer.**

---

## 35.3 The Doppler Formula

The general formula (for sound, when both source and observer may be moving):

```
f_observed = f_source × (v_wave + v_observer) / (v_wave + v_source)

where:
  f_observed = frequency heard by observer
  f_source   = frequency emitted by source
  v_wave     = speed of sound (≈ 343 m/s in air at 20°C)
  v_observer = speed of observer (positive if moving toward source)
  v_source   = speed of source (positive if moving away from observer)
```

**Sign convention:**
- Observer moving TOWARD source: v_observer is positive (numerator increases → f increases)
- Observer moving AWAY from source: v_observer is negative
- Source moving AWAY from observer: v_source is positive (denominator increases → f decreases)
- Source moving TOWARD observer: v_source is negative

---

## 35.4 Moving Source, Stationary Observer

When the observer is stationary (v_observer = 0):

```
f_observed = f_source × v_wave / (v_wave + v_source)

where v_source is positive when moving away, negative when moving toward observer.
```

**Moving toward observer:**

```
v_source is negative (use -v_s):
f_observed = f_source × v / (v - v_s)

Since denominator < v, f_observed > f_source → higher pitch
```

**Moving away from observer:**

```
v_source is positive:
f_observed = f_source × v / (v + v_s)

Since denominator > v, f_observed < f_source → lower pitch
```

### Worked Example 35.1

An ambulance siren emits 800 Hz. It approaches at 30 m/s. Speed of sound = 340 m/s.

**Moving toward:**
f = 800 × 340 / (340 - 30) = 800 × 340/310 = **877 Hz** (higher pitch)

**Moving away:**
f = 800 × 340 / (340 + 30) = 800 × 340/370 = **735 Hz** (lower pitch)

The driver hears 800 Hz the whole time — only the observer notices the change!

---

## 35.5 Stationary Source, Moving Observer

When the source is stationary (v_source = 0):

```
f_observed = f_source × (v_wave + v_observer) / v_wave

where v_observer is positive when moving toward source.
```

**Moving toward source:**

```
f_observed = f_source × (v + v_o) / v > f_source → higher pitch
```

**Moving away:**

```
f_observed = f_source × (v - v_o) / v < f_source → lower pitch
```

### Worked Example 35.2

A stationary siren emits 1000 Hz. An observer runs toward it at 5 m/s. Speed of sound = 340 m/s.

f = 1000 × (340 + 5) / 340 = 1000 × 345/340 ≈ **1015 Hz**

A small effect — because the observer is much slower than sound.

---

## 35.6 Both Moving

When both source and observer move, apply both corrections:

```
f_observed = f_source × (v_wave + v_observer) / (v_wave + v_source)
```

Be careful with sign conventions!

### Worked Example 35.3

A train emits 500 Hz. Both train and observer move toward each other.
Train speed: 40 m/s. Observer speed: 10 m/s. Speed of sound: 340 m/s.

Using conventions:
- v_observer = +10 (moving toward source)
- v_source = -40 (moving toward observer)

f = 500 × (340 + 10) / (340 + (-40)) = 500 × 350/300 ≈ **583 Hz**

---

## 35.7 The Sonic Boom

What happens when a source moves at the speed of sound?

```
SOURCE AT SPEED OF SOUND:

All wavefronts pile up AT the source, creating a giant pressure wall:

  O O O O O [source->]  <- all crests on top of each other

This is the "sound barrier." The source always stays at its own wavefront.

SOURCE FASTER THAN SOUND (supersonic):

  [source>]
           \   <- wake (Mach cone)
            \
             O
              O
               O

The source outpaces its own waves. All wave crests form a cone-shaped 
shockwave behind the source.
```

When this cone passes an observer on the ground, they hear a sudden, explosive **sonic boom** — the pressure discontinuity of many wave crests arriving simultaneously.

### Mach Number

```
Mach number = v_source / v_sound

Mach 1 = speed of sound
Mach 2 = twice speed of sound
```

Supersonic aircraft traveling at Mach 1.5 create a continuous sonic boom trailing behind them.

---

## 35.8 Doppler Effect with Light (Redshift)

The Doppler effect also applies to light and other electromagnetic waves. However, since light requires no medium, the formula is different (based on relativity):

```
For low speeds (v << c):
  f_observed ≈ f_source × (1 ± v/c)

  + if source and observer moving toward each other (blueshift)
  - if moving apart (redshift)
```

### Redshift and Blueshift

```
SOURCE MOVING AWAY from observer:
  Wavelength stretches out → longer wavelength → shifts toward red end of spectrum
  Called "REDSHIFT"
  
  ||||| original   →→→→→→→→→→   stretched    ||||||||

SOURCE MOVING TOWARD observer:
  Wavelength compressed → shorter wavelength → shifts toward blue end
  Called "BLUESHIFT"
  
  ||||||||| original   →→→→  compressed  ||||
```

### Hubble's Discovery

In 1929, Edwin Hubble observed that almost all galaxies show **redshift** — they are moving away from us. Moreover, the farther the galaxy, the greater the redshift.

This is the direct evidence that the universe is **expanding**.

```
GALAXY REDSHIFT:
  
  Expected spectrum:      |  |   |   ||
  Observed spectrum:      |    |     |    ||   (shifted right = longer λ)
  
  The shift tells us: how fast the galaxy is moving away.
  Distance tells us:  how far away it is.
  
  Result: v = H₀ × d  (Hubble's Law)
  where H₀ ≈ 70 km/s per Megaparsec
```

---

## 35.9 Real-World Applications

### Speed Radar (Traffic and Weather)

Police radar guns send out radio waves. These reflect off a moving car and return Doppler-shifted.

```
RADAR GUN:
  
  --- radio wave f₀ --->  [car moving at v]
  <--- reflected wave f₁ ---
  
  Shift (f₁ - f₀) tells us car's speed exactly.
  
  Same principle in weather radar: Doppler shift of
  reflected radio waves tells meteorologists wind speed inside storms.
```

### Medical Ultrasound

Doppler ultrasound measures blood flow. High-frequency sound waves reflect off moving red blood cells, and the Doppler shift indicates flow speed and direction.

```
DOPPLER ULTRASOUND:
  
  Probe sends 5 MHz sound into body.
  Reflected from blood cells moving at v.
  Frequency shift = 2v × f₀ / c_sound
  
  Displayed as: direction (color) and speed (brightness) of blood flow.
  Used to: detect blocked arteries, check fetal heart, find blood clots.
```

### Astronomy: Stellar Velocities

By analyzing the spectral lines of a star (which have known rest frequencies), astronomers can determine:
- Whether the star is moving toward or away from us
- Exactly how fast
- For double stars: the orbital period from cyclical Doppler shifts

### Sonar (Submarine Detection)

Navy sonar systems use the Doppler effect to determine submarine speed and direction. Active sonar sends pulses and measures the Doppler shift of the return.

---

## Summary

- **Doppler effect**: apparent change in frequency due to relative motion between source and observer
- Source moving **toward** observer → higher frequency (pitch goes up)
- Source moving **away** from observer → lower frequency (pitch goes down)
- Formula (for sound): f_obs = f_src × (v ± v_obs) / (v ∓ v_src)
- **Sonic boom**: when source reaches/exceeds speed of sound, creates shockwave cone
- **Redshift**: galaxy/star moving away → wavelength increases (toward red)
- **Blueshift**: approaching → wavelength decreases (toward blue)
- Evidence for expanding universe: nearly all galaxies are redshifted
- Applications: speed radar, Doppler ultrasound, weather radar, astronomy

---

## Key Equations

```
Doppler effect (general):
  f_obs = f_src × (v_wave + v_obs) / (v_wave + v_src)

Sign convention:
  v_obs: positive when moving toward source
  v_src: positive when moving away from observer

Special cases:
  Moving source, stationary observer:
    Toward: f_obs = f_src × v / (v - v_s)
    Away:   f_obs = f_src × v / (v + v_s)

  Moving observer, stationary source:
    Toward: f_obs = f_src × (v + v_o) / v
    Away:   f_obs = f_src × (v - v_o) / v

Mach number:
  M = v_source / v_sound

Redshift (low speed approximation):
  Δf / f ≈ v / c  (recession velocity)

Speed of sound in air:
  v ≈ 343 m/s at 20°C
  v increases with temperature
```
