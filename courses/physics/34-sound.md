# Chapter 34: Sound

> "Without silence, there would be no music." — but without sound, there would be no silence to appreciate.

---

## Table of Contents

1. What Is Sound?
2. Sound as a Longitudinal Pressure Wave
3. Compressions and Rarefactions — In Detail
4. Speed of Sound
5. Why Sound Travels Faster in Denser Media
6. Pitch and Frequency
7. Loudness and Amplitude
8. The Decibel Scale
9. Common Sound Levels Table
10. The Human Hearing Range
11. Ultrasound — Above 20 kHz
12. Infrasound — Below 20 Hz
13. Resonance in Musical Instruments
14. Resonance in Strings
15. Resonance in Air Columns (Pipes)
16. Interference of Sound and Beats
17. Worked Example 1 — Wave Equation for Sound
18. Worked Example 2 — Beat Frequency
19. Worked Example 3 — Decibel Reasoning
20. Summary
21. Key Equations

---

## 1. What Is Sound?

**Sound** is a mechanical wave that travels through a medium (air, water, solids) as a series of pressure variations. It is produced by any vibrating object — a speaker cone, a vocal cord, a tuning fork, a drum skin.

Sound is the way your brain interprets rapid pressure fluctuations at your eardrum. The whole chain is:

```
    Vibrating source → Pressure waves through air → Eardrum vibrates
    → Tiny bones in ear amplify vibration → Cochlea converts to nerve signals
    → Brain interprets as sound
    
    Example: speaker cone
    
    Speaker:    → → →
    Cone moves  compressions     travel     arrive at     eardrum
    outward →   form in air  →  through →  your ear  →   vibrates
    
    Cone moves  rarefactions    travel     arrive at     eardrum
    inward  ←   form in air  →  through →  your ear  →   vibrates
```

---

## 2. Sound as a Longitudinal Pressure Wave

Sound is a **longitudinal wave**: air molecules move back and forth in the same direction as the wave travels. (Review Chapter 33 if needed.)

When a speaker cone pushes forward:
- It pushes air molecules together in front of it
- This creates a **compression** — a region of higher pressure
- The compression region propagates outward at the speed of sound

When the speaker cone pulls back:
- It creates a gap in front of it
- Air molecules rush back, creating a **rarefaction** — a region of lower pressure
- This rarefaction also propagates outward

```
    Speaker at left, wave traveling right:
    
    TIME 1:   Speaker pushes forward
              S|>>>||||  |  ||||  |  ||||
              speaker  comp  rare  comp  rare
    
    TIME 2:   Speaker pulls back
              S<<<|  ||||  |  ||||  |  ||||
    
    TIME 3:   Speaker pushes forward again
              S|>>>||||||  |  ||||||  |  ||||
    
    Key:
    ||||  = compression (high pressure, molecules close together)
    |  |  = rarefaction (low pressure, molecules spread out)
```

The pressure at any point fluctuates above and below the normal atmospheric pressure (about 101,325 Pa). For ordinary conversation, these fluctuations are incredibly small — just a few thousandths of a Pascal above and below atmospheric pressure!

---

## 3. Compressions and Rarefactions — In Detail

Let's zoom in and look at what individual air molecules are doing:

```
    Snapshot of a sound wave in air:
    
    Compression    Rarefaction    Compression    Rarefaction
    
    . . . ...  .  .   .      . . . ...  .  .   .
    . . ...... .  .   .      . . ...... .  .   .
    . . ...  . .  .          . . ...  . .  .
    ← more crowded →  ← less crowded →  ← more crowded →
    
    ← high P →   ← low P →   ← high P →   ← low P →
    
    Molecules in a compression: close together, pushing on neighbors.
    Molecules in a rarefaction: spread apart, less pressure on neighbors.
    
    Each MOLECULE only moves a tiny bit left and right.
    The WAVE PATTERN moves to the right.
```

The wavelength λ of a sound wave is the distance from one compression to the next compression (or one rarefaction to the next).

The amplitude of a sound wave is the maximum displacement of air molecules from their rest position — or equivalently, the maximum pressure variation from atmospheric pressure.

---

## 4. Speed of Sound

The speed of sound in air at room temperature (20°C) is approximately:

```
    v_sound ≈ 343 m/s  (at 20°C)
    v_sound ≈ 331 m/s  (at 0°C)
```

For comparison:
- A fast car travels at ~33 m/s (120 km/h) — about 10× slower than sound
- A commercial airplane travels at ~250 m/s — still slower than sound
- A supersonic jet (Mach 1+) travels faster than sound

### Effect of Temperature

Sound travels faster in warmer air because warmer molecules move faster and can pass compressions along more quickly:

```
    v ≈ 331 × sqrt(1 + T/273)   (in m/s, T in degrees Celsius)
    
    At 20°C: v = 331 × sqrt(1 + 20/273) = 331 × sqrt(1.073) = 343 m/s
    At 0°C:  v = 331 × sqrt(1 + 0/273) = 331 × sqrt(1.000) = 331 m/s
    At -20°C: v = 331 × sqrt(1 + (-20)/273) = 331 × sqrt(0.927) ≈ 318 m/s
```

The rule of thumb: sound speed in air increases by about 0.6 m/s for every degree Celsius rise.

### Sound Speed in Different Media

```
    Medium                  Speed of sound
    -----------------------------------------
    Air (20°C)              343 m/s
    Air (0°C)               331 m/s
    Carbon dioxide          259 m/s   (slower than air — heavier)
    Water (25°C)            1,497 m/s
    Seawater                1,560 m/s (dissolved salts increase speed)
    Wood (oak)              ~3,850 m/s
    Aluminum                6,420 m/s
    Steel                   5,960 m/s
    Glass                   ~5,600 m/s
    -----------------------------------------
```

Sound travels about 4× faster in water than air, and about 17× faster in steel than air.

---

## 5. Why Sound Travels Faster in Denser Media

This seems counterintuitive — "denser" usually means heavier, which should be slower, right?

The key is that **stiffness** matters more than density for sound speed. In solids and liquids:
- Molecules are much closer together
- The bonding forces (stiffness) between molecules are very strong
- When one molecule is pushed, it immediately pushes its neighbor with great force
- The compression is transmitted extremely rapidly

In gases:
- Molecules are far apart
- A molecule must travel some distance before bumping into a neighbor
- The compression is transmitted much more slowly

```
    In a gas:       o . . . . . o . . . . . o    (far apart)
                    push → wait → bump → wait → bump → (slow)
    
    In a solid:     o-o-o-o-o-o-o-o-o-o-o-o    (bonded together)
                    push → instant push through bonds → fast!
    
    The stiffness of the bonds overwhelms the effect of greater density.
```

A more accurate statement: sound travels faster when the restoring force (stiffness) is higher relative to the density.

```
    v_sound = sqrt(Elastic modulus / density)
    
    Steel: high elastic modulus (stiff), high density → very fast
    Air: low elastic modulus (soft), low density → slow
    The stiffness wins.
```

---

## 6. Pitch and Frequency

**Pitch** is the perceptual quality of sound that we describe as "high" or "low" — like the difference between a mouse squeak and a bass drum.

Pitch is determined by **frequency**:
- High frequency → high pitch (squeaky voice, violin)
- Low frequency → low pitch (bass voice, tuba)

```
    Low pitch (low frequency):        High pitch (high frequency):
    
    |  |    |  |    |  |    |         ||||||||||||||||||||||||
    |  |    |  |    |  |    |         ||||||||||||||||||||||||
    
    Few compressions per second       Many compressions per second
    Long wavelength                   Short wavelength
    Bass drum, thunder, tuba          Flute, whistle, bat sonar
```

### Musical Notes and Frequencies

```
    Note         Frequency (Hz)
    ---------------------------------
    A4 (concert A)  440 Hz
    Middle C (C4)   261.6 Hz
    C3 (one octave below) 130.8 Hz
    C5 (one octave above) 523.2 Hz
    
    An OCTAVE = doubling of frequency.
    C5 is exactly twice the frequency of C4.
```

---

## 7. Loudness and Amplitude

**Loudness** is the perceptual experience of sound volume. Physically, it corresponds to the **amplitude** of the sound wave — how large the pressure variations are.

- Large amplitude (large pressure swings) → loud sound
- Small amplitude (tiny pressure swings) → quiet sound

```
    Quiet sound (small amplitude):
    
    atmospheric pressure: ─────────────────────────────
    pressure:             ─────/\─────/\─────/\─────
                                      (small ripples)
    
    Loud sound (large amplitude):
    
    atmospheric pressure: ─────────────────────────────
    pressure:             ─────/\──────/\──────/\──────
                         /     \      /  \    /  \
                        /       \    /    \  /    \
    ───────────────────/         \──/      \/      \──
                                      (large waves)
```

Remember: E ∝ A² — the energy of a sound wave is proportional to the square of the amplitude. To double the loudness perceived by humans, you need to increase intensity by about 10 times (because the human ear responds logarithmically).

---

## 8. The Decibel Scale

The range of sound intensities that humans can hear is enormous:
- The quietest audible sound (threshold of hearing): I₀ = 10⁻¹² W/m²
- The threshold of pain: I_pain ≈ 1 W/m²
- Ratio: 1 / 10⁻¹² = 10¹² = **one trillion!**

A linear scale would be absurd for this range. Instead, we use the **decibel (dB) scale**, which is logarithmic:

```
    Sound level (dB) = 10 × log₁₀(I / I₀)
    
    Where:
    I  = intensity of the sound (W/m²)
    I₀ = 10⁻¹² W/m²  (reference intensity, threshold of hearing)
    log₁₀ = logarithm base 10
```

### Key Properties of the Decibel Scale

```
    Every 10 dB increase = 10× increase in INTENSITY
    Every 10 dB increase = roughly 2× increase in PERCEIVED LOUDNESS
    
    Examples:
    
    0 dB:   I = 10⁻¹² W/m²   = I₀ × 10^0   = I₀ × 1
    10 dB:  I = 10⁻¹¹ W/m²   = I₀ × 10^1   = I₀ × 10
    20 dB:  I = 10⁻¹⁰ W/m²   = I₀ × 10^2   = I₀ × 100
    30 dB:  I = 10⁻⁹ W/m²    = I₀ × 10^3   = I₀ × 1,000
    60 dB:  I = 10⁻⁶ W/m²    = I₀ × 10^6   = I₀ × 1,000,000
    120 dB: I = 10⁰ = 1 W/m² = I₀ × 10^12  = I₀ × 1,000,000,000,000
```

The threshold of pain (120 dB) is a trillion times more intense than the threshold of hearing (0 dB) — yet the decibel scale compresses this to just 120 units. That is the power of a logarithmic scale.

---

## 9. Common Sound Levels Table

```
    Sound Level      dB      Description / Example
    ─────────────────────────────────────────────────────────────
    0 dB             0       Threshold of hearing — barely audible
    Rustling leaves  10      Extremely quiet room at night
    Quiet whisper    20      Hearing test quiet booth
    Library          30      Very quiet
    Quiet office     40      Typical quiet office environment
    Normal conversation 60   Face-to-face chat at 1 meter
    Busy restaurant  70      Background noise in city cafe
    Lawn mower       80      Prolonged exposure causes damage
    Heavy traffic    90      City street with trucks
    Subway train    100      Inside a subway car
    Rock concert    110      Front row at a live concert
    Threshold of pain 120    Pain begins; temporary hearing damage
    Jet engine (nearby) 130  Permanent hearing damage possible
    Rocket launch   140-180  Lethal at close range
    ─────────────────────────────────────────────────────────────
    
    Safe limit for prolonged exposure: below 85 dB
    For every 3 dB increase above 85, safe exposure time halves.
```

### Why the Decibel Scale Feels Natural

Human hearing is approximately logarithmic. We perceive equal increases in dB as equal increases in apparent loudness — even though the actual intensity changes by powers of 10. Evolution tuned our ears to detect a massive range of sound levels by compressing the signal logarithmically.

---

## 10. The Human Hearing Range

Humans can hear sounds with frequencies from roughly:

```
    20 Hz  ─────────────────────────────────────────────  20,000 Hz
    (20 Hz)                                            (20 kHz)
    
    Low pitch:                                        High pitch:
    Bass drum rumble                                  Bat sonar, dog whistle
    Thunder                                           
    Lowest piano note (27.5 Hz)     Highest piano note (4,186 Hz)
```

This range varies significantly with age:
- Children: can hear up to 20,000 Hz or beyond
- By age 25: typically lose ability to hear above ~16,000 Hz
- By age 60: typical upper limit is ~8,000–12,000 Hz

Below 20 Hz: **infrasound**
Above 20,000 Hz (20 kHz): **ultrasound**

---

## 11. Ultrasound — Above 20 kHz

**Ultrasound** is sound with frequency above 20 kHz (too high for humans to hear). Despite being inaudible to us, it has many important applications.

### Medical Ultrasound Imaging

Frequencies used: 1–20 MHz (1,000,000–20,000,000 Hz)

A transducer (piezoelectric crystal) emits a short pulse of ultrasound into the body. The pulse reflects off boundaries between tissues of different densities (e.g., between soft tissue and bone, or between a fluid-filled cyst and surrounding tissue). The time delay between emitting and receiving each reflection tells the machine how deep the boundary is:

```
    Ultrasound transducer (probe):
    
    [transducer]───pulse→───────────────────────
                            body tissues
                            
    ←─reflected echo──────────────────────────
    
    Time delay tells depth:
    depth = (v_sound × time) / 2
    (divide by 2 because pulse travels there AND back)
    
    By scanning across the body, a 2D image is built up.
```

### Why Ultrasound Instead of X-rays for Babies?

- No ionizing radiation — completely safe
- Can show soft tissues that X-rays cannot (X-rays show bone well but soft tissue poorly)
- Real-time imaging (can see the baby moving, heartbeat, etc.)

### Sonar (Sound Navigation And Ranging)

Used by:
- **Ships** to map the ocean floor (echo sounding)
- **Submarines** for navigation and detecting objects
- **Bats** for echolocation (they emit ~20–200 kHz pulses)
- **Dolphins** for echolocation (~100–150 kHz)

```
    Ship sonar:
    
    [ship]──────────────surface─────────────────────
           |
           | pulse down
           |
           ↓
    ───────────────────ocean floor────────────────
    
    Reflected pulse travels back up.
    Time measured → depth calculated.
    depth = v × t / 2
```

### Industrial Ultrasound

- **Non-destructive testing**: detect cracks inside metal parts (aircraft engines, pipelines) without cutting them open
- **Cleaning**: ultrasonic cleaners vibrate tiny bubbles in liquid that scrub surfaces (jewelry, medical instruments)
- **Welding plastics**: ultrasound at the joint melts and fuses plastic parts

---

## 12. Infrasound — Below 20 Hz

**Infrasound** is sound with frequency below 20 Hz (too low for humans to hear). It can travel enormous distances with little absorption.

### Natural Infrasound Sources

- **Elephants**: communicate using infrasound at 14–35 Hz over distances of several kilometers. The low frequencies travel well through the ground. Elephants feel infrasound vibrations through their feet!
- **Earthquakes and volcanic eruptions**: generate infrasound that can travel around the entire globe
- **Whales**: blue whales communicate at 10–40 Hz, carrying calls thousands of kilometers across oceans
- **Storms and weather**: severe thunderstorms and tornadoes generate characteristic infrasound patterns
- **Auroras**: the Northern Lights generate infrasound

### Detection of Nuclear Tests

The Comprehensive Nuclear-Test-Ban Treaty Organization (CTBTO) operates a global network of infrasound detectors to detect clandestine nuclear weapon tests. Nuclear explosions generate distinctive infrasound that can be detected thousands of kilometers away.

### Why Can't We Hear Infrasound?

The basilar membrane in our cochlea is tuned to resonate at different frequencies along its length. Below 20 Hz, the displacements of the membrane are too slow for our auditory system to resolve as distinct oscillations. We might feel very intense infrasound as vibration (like standing next to a massive subwoofer), but we cannot hear it as a pitched tone.

### Human Sensitivity to Infrasound

Interestingly, very intense infrasound near 18–19 Hz has been reported to cause feelings of unease, dread, or the sensation that the room is haunted. One proposed mechanism: resonance with eyeballs. This remains somewhat controversial scientifically.

---

## 13. Resonance in Musical Instruments

All musical instruments use **resonance** to amplify sound. A vibrating object (plucked string, vibrating reed, struck drum skin) makes very little sound on its own — not enough surface area to push much air. The instrument's resonating body (wooden box, pipe column, drum shell) is designed to resonate at the same frequencies, amplifying the vibrations dramatically.

The object vibrates at its **natural frequencies** (also called **harmonics** or **overtones**). These are specific frequencies that form standing waves in the instrument — covered in detail in Chapter 36. Here, we introduce the basic concepts for strings and pipes.

---

## 14. Resonance in Strings

When a string fixed at both ends is plucked, it can only vibrate at frequencies where the standing wave fits exactly. The lowest such frequency is the **fundamental frequency** (f₁).

```
    Guitar string fixed at both ends (length L):
    
    Fundamental (n=1):
    
    ╔══════════════════════╗
    ║  ╱‾‾‾‾‾‾‾‾‾‾‾‾‾‾╲   ║
    ║ ╱                 ╲  ║
    ║                      ║
    ║╲                  ╱  ║
    ╚══╲════════════╱════╝
    
    One half-wavelength fits in the string: L = λ/2 → λ = 2L
    
    Second harmonic (n=2):
    
    ╔══════════════════════╗
    ║  ╱‾‾‾‾╲  ╱‾‾‾‾╲     ║
    ║ ╱      ╲╱      ╲    ║
    ╚═══════════════════╝
    
    One full wavelength fits: L = λ → λ = L
    
    Third harmonic (n=3):
    
    ╔══════════════════════╗
    ║  ╱‾‾╲  ╱‾‾╲  ╱‾‾╲   ║
    ║ ╱    ╲╱    ╲╱    ╲   ║
    ╚═══════════════════╝
    
    3/2 wavelengths: L = 3λ/2 → λ = 2L/3
```

The fundamental frequency for a string is:
```
    f₁ = v_string / (2L)
```

Where v_string depends on the tension and mass per unit length of the string. Higher tension → higher frequency (higher pitch). That's why tuning a guitar means adjusting string tension.

---

## 15. Resonance in Air Columns (Pipes)

Wind instruments (flutes, clarinets, trumpets, organ pipes) create sound using standing waves in air columns. The shape of the pipe determines which harmonics are present.

### Open Pipe (open at both ends — like a flute)

```
    Open pipe (length L):
    
    Both ends are antinodes (maximum movement — air can move freely at open ends).
    
    Fundamental:
    
    ←────────L────────→
    A────────────────A    (antinodes at both ends)
         N             (node in middle)
    
    λ/2 = L → λ = 2L
    f₁ = v/2L
    
    All harmonics are present: f₁, 2f₁, 3f₁, 4f₁ ...
```

### Closed Pipe (closed at one end — like a clarinet's lower register)

```
    Closed pipe (length L):
    
    Closed end: NODE (air cannot move)
    Open end:   ANTINODE (air moves freely)
    
    Fundamental:
    
    ←────────L────────→
    N────────────────A
    
    λ/4 = L → λ = 4L
    f₁ = v/4L
    
    Only ODD harmonics are present: f₁, 3f₁, 5f₁ ...
    (no 2f₁, 4f₁, etc.)
    
    This is why a clarinet sounds different from a flute —
    the missing even harmonics give it a hollow, reedy tone.
```

---

## 16. Interference of Sound and Beats

When two sound waves of slightly different frequencies reach your ear simultaneously, they interfere. Because the frequencies are slightly different, they are sometimes in phase (constructive) and sometimes out of phase (destructive), alternating regularly.

This produces **beats** — a throbbing, wavering loudness at regular intervals.

```
    Wave 1: f₁ = 440 Hz
    .....|.....|.....|.....|.....|.....| (rapid)
    
    Wave 2: f₂ = 443 Hz  (slightly higher)
    ....|....|....|....|....|....|....|
    
    Sometimes in phase           Sometimes out of phase
    (loud)                       (quiet)
    
    What you HEAR:
    
          LOUD    quiet    LOUD    quiet    LOUD
    ──────/\─────────────/\─────────────/\─────
    
    Beat frequency = |f₁ - f₂| = |440 - 443| = 3 Hz
    You hear 3 loud-soft cycles per second = 3 beats per second.
```

**Beat frequency formula**:
```
    f_beat = |f₁ - f₂|
    
    (absolute value — always positive)
```

### Applications of Beats

- **Tuning musical instruments**: a piano tuner listens for beats between a tuning fork and the piano string. When the beat frequency goes to zero (no beats), the string is in tune.
- **Hearing tests (audiometry)**: beats can help assess hearing sensitivity.
- **Radio receivers**: the principle of heterodyne detection uses beats in electronics.

### Why Piano Tuners Listen for Beats

A concert A tuning fork vibrates at exactly 440 Hz. The piano string's current pitch might be slightly off, say 437 Hz. The tuner holds the fork near the piano and plays the A key — they hear 3 beats per second (|440-437|=3). They tighten the string until the beats slow down and disappear. Zero beats = perfectly in tune!

---

## 17. Worked Example 1 — Wave Equation for Sound

**Problem**: A sound wave in air (v = 343 m/s) has a frequency of 1,000 Hz. What is the wavelength?

**Given**:
- v = 343 m/s
- f = 1,000 Hz

**Formula**: v = fλ → λ = v/f

**Solution**:
```
    λ = v / f
    λ = 343 m/s ÷ 1,000 Hz
    λ = 0.343 m = 34.3 cm
```

**Answer**: The wavelength is 0.343 m (about 34 cm).

**Context**: 1,000 Hz (1 kHz) is a mid-range frequency — roughly the pitch of the phone dial tone. A wavelength of 34 cm is about the size of a large book — reasonable for a sound we can easily hear.

---

## 18. Worked Example 2 — Beat Frequency

**Problem**: Two guitar strings are plucked simultaneously. One vibrates at 329 Hz and the other at 332 Hz. How many beats per second do you hear?

**Given**:
- f₁ = 329 Hz
- f₂ = 332 Hz

**Formula**: f_beat = |f₁ - f₂|

**Solution**:
```
    f_beat = |329 - 332|
    f_beat = |-3|
    f_beat = 3 Hz
```

**Answer**: You hear 3 beats per second.

**Context**: 3 beats per second is noticeable and clearly audible — most musicians can easily detect beat frequencies below about 6–8 Hz. A musician would tighten or loosen one string slightly until the beats disappear.

---

## 19. Worked Example 3 — Decibel Reasoning

**Problem**: A factory machine produces a sound level of 90 dB. A second identical machine is turned on nearby. What is the new sound level?

**Key insight**: Adding two identical sources doubles the intensity. A doubling of intensity corresponds to a 3 dB increase on the decibel scale.

**Reasoning**:
```
    Intensity from one machine:  I
    Intensity from two machines: 2I
    
    Change in dB = 10 × log₁₀(2I / I) = 10 × log₁₀(2) = 10 × 0.301 = 3.01 dB
    
    New sound level = 90 dB + 3 dB = 93 dB
```

**Answer**: 93 dB.

**Note**: This is why hearing protection matters in factories. Adding machines doesn't just add dB linearly — but even a 3 dB increase represents a doubling of energy entering your ears, and the effects on hearing damage are significant over time.

---

## 20. Summary

- **Sound** is a mechanical longitudinal wave: compressions and rarefactions travel through air (or another medium), carrying energy.

- Sound cannot travel through a vacuum — it needs a medium.

- **Speed of sound** in air at 20°C is ~343 m/s. Sound travels faster in liquids (~1,480 m/s in water) and solids (~5,960 m/s in steel) because the molecules are closer together and more tightly bonded.

- **Pitch** corresponds to frequency. Higher frequency = higher pitch.

- **Loudness** corresponds to amplitude. Greater amplitude = louder sound.

- **Decibel scale**: logarithmic scale for sound intensity. 0 dB = threshold of hearing. 120 dB = threshold of pain. Every 10 dB = 10× intensity change.

- **Human hearing range**: approximately 20 Hz to 20,000 Hz (20 kHz).

- **Ultrasound** (>20 kHz): used in medical imaging (1–20 MHz), sonar, and industrial testing. Echoes give depth information via depth = vt/2.

- **Infrasound** (<20 Hz): used by elephants and whales for long-distance communication. Used by CTBTO to detect nuclear tests.

- **Resonance in instruments**: strings and air columns vibrate at specific natural frequencies (harmonics). Open pipes support all harmonics; closed pipes support only odd harmonics.

- **Beats**: two slightly different frequencies create a periodic loud-soft pattern. f_beat = |f₁ - f₂|. Used for tuning instruments.

---

## 21. Key Equations

```
Wave equation:
    v = f × λ
    λ = v / f
    f = v / λ

Speed of sound in air (approximate rule of thumb):
    v ≈ 343 m/s  at 20°C
    v increases ~0.6 m/s per °C rise

Decibel scale:
    dB = 10 × log₁₀(I / I₀)
    I₀ = 10⁻¹² W/m²  (threshold of hearing)

Decibel for intensity doubling:
    ΔdB = 10 × log₁₀(2) ≈ 3 dB

Beat frequency:
    f_beat = |f₁ - f₂|

Depth from sonar/ultrasound echo:
    depth = v × t / 2
    (t = time for echo to return, v = speed of sound in medium)

Fundamental frequency of a string (fixed both ends):
    f₁ = v_string / (2L)
    
Fundamental frequency of open pipe:
    f₁ = v_sound / (2L)

Fundamental frequency of closed pipe (one end closed):
    f₁ = v_sound / (4L)
```
