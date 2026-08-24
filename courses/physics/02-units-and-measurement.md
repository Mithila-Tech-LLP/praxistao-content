# Chapter 02: Units and Measurement

> **"Without a standard unit of measurement, science is just opinion. A metre is a metre everywhere on Earth — and that changes everything."**

---

## Table of Contents

- [2.1 Why Standard Units Matter](#21-why-standard-units-matter)
- [2.2 The SI System — Seven Base Units](#22-the-si-system--seven-base-units)
- [2.3 Derived Units — Built From Base Units](#23-derived-units--built-from-base-units)
- [2.4 Prefixes — From Nano to Tera](#24-prefixes--from-nano-to-tera)
- [2.5 Unit Conversion — A Step-by-Step Method](#25-unit-conversion--a-step-by-step-method)
- [2.6 Scientific Notation](#26-scientific-notation)
- [2.7 Significant Figures](#27-significant-figures)
- [2.8 Accuracy vs Precision](#28-accuracy-vs-precision)
- [2.9 Dimensional Analysis](#29-dimensional-analysis)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 2.1 Why Standard Units Matter

### The $327 Million Mistake

On September 23, 1999, NASA lost a spacecraft worth $327 million dollars. Not because of a technical failure. Not because of a bad rocket. Because two teams used different units of measurement.

The Mars Climate Orbiter had been travelling through space for 286 days. It was supposed to enter orbit around Mars. Instead, it came in at the wrong angle, hit the Martian atmosphere, and burned up.

The cause? One engineering team at NASA used **metric units** (Newton-seconds). The other team, at Lockheed Martin, used **imperial units** (pound-force-seconds). The navigation software received data in one unit and treated it as if it were the other.

The difference was tiny in the raw numbers. The consequences were catastrophic.

```
THE MARS CLIMATE ORBITER DISASTER

        Earth                                    Mars
          |                                        |
     [Launch]                               [Intended orbit]
          |                                        |
          |-------- 286 days of travel ----------->|
          |                                        |
     "We sent                               "We expected
      thrust data                            metric units"
      in pounds"                                   |
                                            [WRONG ANGLE]
                                                   |
                                            [Burns up in
                                             atmosphere]
                                                   X
                                           $327M lost
```

This story is not an isolated incident. Throughout history, confusion over units has caused:

- Aircraft running out of fuel mid-flight (Air Canada Flight 143, 1983 — the plane was fuelled in pounds instead of kilograms)
- Bridge collapses from incorrect load calculations
- Drug overdoses in hospitals due to mg vs mcg confusion
- Construction errors when metric and imperial blueprints are mixed

The lesson is simple: **units are not optional decoration on a number. They are part of the number itself.**

### What is a Unit?

When you say "the room is 5 wide," that sentence is meaningless. Five what? Five centimetres? Five metres? Five football fields?

A **unit** is a standard quantity used to express and compare measurements. When you say "the room is 5 metres wide," everyone on Earth knows exactly how wide you mean — because a metre is defined precisely and consistently everywhere.

Without agreed-upon units:
- Scientists in different countries could not share results meaningfully
- Engineers in different companies could not build systems together
- Doctors could not prescribe doses safely across different healthcare systems
- Space missions would crash

### A Brief History of Units

Before standardisation, units were chaotic:
- A **foot** was literally the length of a king's foot (which changed with each new king)
- A **cubit** was the length of a forearm from elbow to fingertip (which varied person to person)
- A **furlong** was the distance a team of oxen could plough without resting
- Different cities in the same country used different units for the same quantity

In 1799, France introduced the **metric system** — the first attempt at a universal, logical measurement system based on powers of ten. By 1960, the world had adopted the **International System of Units** (SI, from the French *Système International d'Unités*), which is what the scientific world uses today.

The SI system has two beautiful properties:
1. It is **decimal** — everything is based on powers of ten, making conversions easy
2. It is **universal** — one and only one unit for each quantity

---

## 2.2 The SI System — Seven Base Units

The entire structure of physics measurement is built on just **seven base units**. Everything else — every measurement you will ever make in physics — is derived from these seven.

Think of them like the primary colours of measurement. Just as you can mix red, blue, and yellow to make any colour, you can combine these seven units to describe anything in the universe.

```
THE SEVEN BASE UNITS OF THE SI SYSTEM

+------------------+------------------+----------+
|   QUANTITY       |   UNIT NAME      |  SYMBOL  |
+------------------+------------------+----------+
|   Length         |   metre          |    m     |
|   Mass           |   kilogram       |    kg    |
|   Time           |   second         |    s     |
|   Electric       |   ampere         |    A     |
|   current        |                  |          |
|   Temperature    |   kelvin         |    K     |
|   Amount of      |   mole           |    mol   |
|   substance      |                  |          |
|   Luminous       |   candela        |    cd    |
|   intensity      |                  |          |
+------------------+------------------+----------+
```

Let's understand each one with real examples and a feel for the scale.

### The Metre (m) — Length

The **metre** is the base unit of length, distance, and displacement.

Today, the metre is defined as the distance light travels in a vacuum in exactly 1/299,792,458 of a second. This definition ties the metre to the speed of light — a universal constant — making it extraordinarily precise and reproducible anywhere in the universe.

Real-world scale:
- Thickness of a human hair: about 0.00007 m (70 micrometres)
- Length of a pen: about 0.15 m
- Height of a door: about 2 m
- Length of a football pitch: about 100 m
- Height of Mount Everest: 8,849 m
- Distance from Mumbai to Delhi: about 1,150,000 m (1,150 km)

```
LENGTH SCALE IN METRES

          0.000000001 m  ← diameter of an atom (~0.1 nm)
          0.000001 m     ← bacterium (~1 micrometre)
          0.001 m        ← thickness of a credit card (~1 mm)
          0.01 m         ← width of a thumbnail (~1 cm)
          1 m            ← your arm's length (roughly)
          100 m          ← length of a football pitch
          1,000 m        ← 1 km — a short walk
          1,000,000 m    ← 1000 km — roughly Mumbai to Delhi
          400,000,000 m  ← Earth-Moon distance
```

### The Kilogram (kg) — Mass

The **kilogram** is the base unit of mass.

Important distinction: **mass** is not the same as weight.
- **Mass** is the amount of matter in an object. It is measured in kilograms.
- **Weight** is the force that gravity pulls on that mass. It is measured in Newtons.

Your mass is the same on Earth and on the Moon. Your weight is different (about 6 times less on the Moon, because the Moon's gravity is weaker).

Since 2019, the kilogram is defined using a fundamental constant of nature called Planck's constant (h = 6.626 × 10⁻³⁴ J·s), making it perfectly reproducible anywhere in the universe — no longer tied to a physical object.

Real-world scale:
- A paper clip: about 0.001 kg (1 gram)
- A teaspoon of sugar: about 0.004 kg (4 grams)
- A litre of water: exactly 1 kg (by the original metric definition)
- A bag of rice: about 5 kg
- An adult human: about 60–80 kg
- A small car: about 1,200–1,500 kg

### The Second (s) — Time

The **second** is the base unit of time.

Today, the second is defined using the vibrations of a caesium-133 atom. A caesium atomic clock "ticks" exactly 9,192,631,770 times per second. These clocks are so accurate they would not gain or lose one second in 300 million years.

Real-world scale:
- A heartbeat: about 1 second
- Time to say "one Mississippi": about 1 second
- The time to read this sentence: about 4–6 seconds
- A school period: about 2,700–3,600 seconds (45–60 minutes)
- One day: 86,400 seconds
- One year: about 31,536,000 seconds (3.15 × 10⁷ s)

### The Ampere (A) — Electric Current

The **ampere** (often shortened to "amp") is the base unit of electric current.

**Electric current** is the flow of electric charge through a conductor. When electrons move through a wire, that is an electric current. When you plug in a phone charger, current flows through the cable.

Real-world scale:
- A phone charger: about 1–2 A
- A household light bulb (100 W): about 0.45 A
- A hair dryer: about 10 A
- A car starter motor: about 100–200 A
- A lightning bolt: about 20,000–30,000 A (but only for about 0.2 seconds)

### The Kelvin (K) — Temperature

The **kelvin** is the base unit of temperature in physics.

You are probably familiar with Celsius (°C). The kelvin scale works the same way — one degree of change is the same size on both scales — but it starts at **absolute zero**, the coldest possible temperature, where all atomic and molecular motion stops completely.

The conversion is simple:
```
K = °C + 273.15
°C = K − 273.15
```

Real-world scale:
- Absolute zero: 0 K = −273.15 °C
- Liquid nitrogen: 77 K = −196 °C
- Dry ice (solid CO₂): 195 K = −78 °C
- Water freezes: 273.15 K = 0 °C
- Room temperature: about 293 K = 20 °C
- Human body temperature: about 310 K = 37 °C
- Water boils: 373.15 K = 100 °C
- Surface of the Sun: about 5,778 K

Note: We write "273 kelvin" (not "273 degrees kelvin"). There is no degree symbol for kelvin. Kelvin is an absolute scale, not a relative one.

### The Mole (mol) — Amount of Substance

The **mole** is a counting unit for incredibly tiny particles like atoms and molecules.

One mole = 6.022 × 10²³ particles. This enormous number is called **Avogadro's number** (Nₐ).

Why such a huge number? Because atoms are unimaginably small. If you counted one atom per second, it would take you about 19 quadrillion years to count a single mole — roughly a million times longer than the current age of the universe.

The mole is defined so that one mole of carbon-12 atoms has a mass of exactly 12 grams. This makes it very practical for chemistry and materials science.

Real-world perspective:
- One mole of water (H₂O) = 18 grams (a small sip)
- One mole of iron = 56 grams (a small piece)
- One mole of table salt (NaCl) = 58.4 grams

### The Candela (cd) — Luminous Intensity

The **candela** measures the intensity of light as perceived by the human eye.

The name comes from an old standard: one candela is approximately the brightness of one candle. More precisely, it is calibrated to the sensitivity of the human eye — which is most sensitive to yellow-green light.

Real-world scale:
- A single candle: about 1 cd
- A flashlight: about 10–100 cd
- A car headlight (low beam): about 700 cd
- A projector: about 1,000–5,000 cd
- The Sun (as seen from Earth's surface): about 2.7 × 10²⁷ cd

---

## 2.3 Derived Units — Built From Base Units

Most measurements in everyday physics use **derived units** — units formed by combining the seven base units through multiplication and division.

Think of it like building with LEGO bricks. The seven base units are your fundamental bricks. Derived units are structures you build by connecting them in different ways.

### Velocity — metres per second (m/s)

**Velocity** (and speed) tell us how much distance is covered per unit of time.

```
velocity = distance / time
v = d / t

Units: metres / seconds = m/s
```

Worked example: A car travels 150 metres in 6 seconds.
```
v = 150 m / 6 s = 25 m/s
```

### Acceleration — metres per second squared (m/s²)

**Acceleration** tells us how quickly velocity changes per unit of time.

```
acceleration = change in velocity / time
a = Δv / t

Units: (m/s) / s = m/s²
```

Worked example: A car goes from 0 m/s to 20 m/s in 4 seconds.
```
a = (20 − 0) m/s / 4 s = 5 m/s²
```

### Force — Newton (N)

**Force** is what causes objects to accelerate. Isaac Newton's second law states:

```
Force = mass × acceleration
F = m × a

Units: kg × (m/s²) = kg·m/s²
```

This combination appears so often it gets its own name: the **Newton (N)**.

```
1 Newton = 1 kg·m/s²
```

A Newton is roughly the force needed to hold a small apple (100 g) against gravity. The apple pulls down with about 1 N of force.

### Energy — Joule (J)

**Energy** is the capacity to do work. Work is done when a force moves something through a distance.

```
Energy (work) = force × distance
E = F × d

Units: Newton × metre = N·m = kg·m²/s²
```

This combination is called the **Joule (J)**.

```
1 Joule = 1 N·m = 1 kg·m²/s²
```

Real-world scale of Joules:
- Lifting a 100 g apple 1 metre upward: about 1 J
- Lifting a textbook (500 g) onto a 1 m shelf: about 5 J
- A 60 W light bulb uses 60 J every second
- One food calorie (kcal): about 4,184 J
- A lightning bolt: about 5 billion J (5 × 10⁹ J)

### Power — Watt (W)

**Power** is how fast energy is transferred or used.

```
Power = Energy / time
P = E / t

Units: Joule / second = J/s = kg·m²/s³
```

This is called the **Watt (W)**.

```
1 Watt = 1 J/s
```

Real-world Watts:
- A human at rest: about 80 W
- A bicycle rider: about 200–400 W
- A hair dryer: about 1,500 W (1.5 kW)
- A car engine: about 100,000–200,000 W (100–200 kW)

### Pressure — Pascal (Pa)

**Pressure** is the force applied per unit area.

```
Pressure = Force / Area
P = F / A

Units: Newton / metre² = N/m² = kg/(m·s²)
```

This is called the **Pascal (Pa)**.

```
1 Pascal = 1 N/m²
```

Real-world Pascals:
- A butterfly landing on your hand: about 1 Pa
- Atmospheric pressure at sea level: 101,325 Pa (about 101 kPa)
- A car tyre: about 200,000–300,000 Pa (200–300 kPa)
- Pressure at the deepest ocean point (Mariana Trench): about 110,000,000 Pa (110 MPa)

### Summary Table of Common Derived Units

```
+------------+-----------+--------+------------------------+
|  QUANTITY  |   UNIT    | SYMBOL |   IN BASE UNITS        |
+------------+-----------+--------+------------------------+
| Velocity   | m/s       |  m/s   |   m·s⁻¹               |
| Accel.     | m/s²      |  m/s²  |   m·s⁻²               |
| Force      | Newton    |  N     |   kg·m·s⁻²             |
| Energy     | Joule     |  J     |   kg·m²·s⁻²           |
| Power      | Watt      |  W     |   kg·m²·s⁻³           |
| Pressure   | Pascal    |  Pa    |   kg·m⁻¹·s⁻²          |
| Frequency  | Hertz     |  Hz    |   s⁻¹                  |
| Charge     | Coulomb   |  C     |   A·s                  |
| Voltage    | Volt      |  V     |   kg·m²·A⁻¹·s⁻³       |
+------------+-----------+--------+------------------------+
```

---

## 2.4 Prefixes — From Nano to Tera

Physics deals with an enormous range of scales. The diameter of a proton is about 0.000000000000001 metres (10⁻¹⁵ m). The observable universe is about 880,000,000,000,000,000,000,000,000 metres (8.8 × 10²⁶ m) across.

Writing all those zeros is tedious and error-prone. So we use **prefixes** — standard syllables added to the front of any unit name to multiply or divide it by a power of ten.

### The Full Prefix Table

```
+--------+--------+--------+---------------------+----------------------------+
| PREFIX | SYMBOL | POWER  |       VALUE         |   EVERYDAY EXAMPLE         |
+--------+--------+--------+---------------------+----------------------------+
| tera   |   T    |  10¹²  | 1,000,000,000,000   | Hard drive: 2 TB           |
+--------+--------+--------+---------------------+----------------------------+
| giga   |   G    |  10⁹   | 1,000,000,000       | Phone storage: 256 GB      |
|        |        |        |                     | CPU speed: 3.5 GHz         |
+--------+--------+--------+---------------------+----------------------------+
| mega   |   M    |  10⁶   | 1,000,000           | Song file: 5 MB            |
|        |        |        |                     | City distance: 2 Mm away   |
+--------+--------+--------+---------------------+----------------------------+
| kilo   |   k    |  10³   | 1,000               | 1 km = 1,000 m             |
|        |        |        |                     | 1 kg = 1,000 g             |
+--------+--------+--------+---------------------+----------------------------+
| (none) |  ---   |  10⁰   | 1                   | The base unit itself       |
+--------+--------+--------+---------------------+----------------------------+
| centi  |   c    |  10⁻²  | 0.01                | 1 cm = 0.01 m              |
|        |        |        |                     | Ruler markings             |
+--------+--------+--------+---------------------+----------------------------+
| milli  |   m    |  10⁻³  | 0.001               | 1 mm = 0.001 m             |
|        |        |        |                     | Thickness of a coin (~2mm) |
+--------+--------+--------+---------------------+----------------------------+
| micro  |   μ    |  10⁻⁶  | 0.000001            | 1 μs = one-millionth of    |
|        |        |        |                     | a second                   |
|        |        |        |                     | Red blood cell: ~8 μm wide |
+--------+--------+--------+---------------------+----------------------------+
| nano   |   n    |  10⁻⁹  | 0.000000001         | CPU transistor: ~3–5 nm    |
|        |        |        |                     | Wavelength of violet light |
|        |        |        |                     | ~400 nm                    |
+--------+--------+--------+---------------------+----------------------------+
| pico   |   p    |  10⁻¹² | 0.000000000001      | Diameter of an atom:       |
|        |        |        |                     | ~100–300 pm                |
+--------+--------+--------+---------------------+----------------------------+
```

### Visual Scale: Powers of 10

```
SMALLER  <--------------------------------------------->  LARGER

  pm       nm       μm      mm      cm      m      km     Mm
 10⁻¹²   10⁻⁹    10⁻⁶   10⁻³   10⁻²   10⁰    10³    10⁶

[atom] [virus] [cell] [hair] [nail] [you] [city] [continent]
```

### Everyday Prefix Examples in Depth

**nano (n, 10⁻⁹) — The scale of the very small:**
Modern computer chips have transistors measuring just 3–5 nanometres across — a scale so small you could fit roughly 20,000 of them across a single human hair. DNA has a double helix width of about 2 nm. Visible light has wavelengths between 400 nm (violet) and 700 nm (red).

**micro (μ, 10⁻⁶) — The scale of cells:**
A human red blood cell is about 6–8 micrometres in diameter. Bacteria range from 1–10 micrometres. A human hair is about 70 micrometres wide. Your microwave oven works at a frequency of about 2.45 gigahertz — a wavelength of about 12 centimetres (not micro!).

**milli (m, 10⁻³) — The nearly-visible scale:**
A raindrop is typically 1–5 mm in diameter. A standard aspirin tablet is 500 mg. Millimetres are the smallest division on most school rulers. A hummingbird egg is about 12 mm long.

**centi (c, 10⁻²) — Human-scale measurement:**
The centimetre is the unit most people use for everyday measurements. Body height, fabric lengths, box dimensions. A standard sheet of A4 paper is 21 cm × 29.7 cm. A centimetre is about the width of a fingernail.

**kilo (k, 10³) — The scale of distances and large masses:**
Road distances are in kilometres. Grocery shopping uses kilograms. Computer files are often in kilobytes. A kilowatt-hour (kWh) is the standard unit on your electricity bill — it's 3,600,000 joules.

**mega (M, 10⁶) — Large-scale quantities:**
Music and photo files are in megabytes. Large explosions are measured in megatons (1 megaton = 4.18 × 10¹⁵ J). The frequency of FM radio stations is in megahertz (88–108 MHz).

**giga (G, 10⁹) — Computing scale:**
Modern storage devices are in gigabytes. Processor speeds are in gigahertz. The human brain is estimated to have about 86 billion (86 × 10⁹) neurons.

**tera (T, 10¹²) — The scale of big data:**
Large hard drives are 1–4 terabytes. The global internet carries roughly 5 exabytes (5 × 10¹⁸ bytes) of data per day — about 5 million terabytes.

---

## 2.5 Unit Conversion — A Step-by-Step Method

One of the most critical practical skills in physics is converting between units. There is a single foolproof method that always works: **multiply by a conversion factor equal to 1**.

### The Core Idea

Any fraction where the numerator and denominator represent the same quantity is equal to 1.

```
1000 m
------ = 1     (1000 m and 1 km are the same distance)
 1 km

60 s
---- = 1       (60 seconds and 1 minute are the same time)
1 min
```

When you multiply any measurement by such a fraction, you don't change its value — you only change how it is expressed. Choose the fraction so that the unit you want to eliminate appears in the opposite position (numerator or denominator) so it cancels out.

### Step-by-Step Method

```
STEP 1: Write down the number with its unit
STEP 2: Identify the conversion relationship (e.g., 1 km = 1000 m)
STEP 3: Write it as a fraction equal to 1
STEP 4: Arrange so the old unit will cancel (old unit on opposite side)
STEP 5: Multiply the numbers
STEP 6: Write the new unit (whatever was left after cancellation)
STEP 7: Check — does the answer make intuitive sense?
```

### Worked Example 1: Convert 90 km/h to m/s

This is probably the single most common unit conversion in a physics course.

We need two conversion factors:
- km → m: 1 km = 1000 m
- h → s: 1 h = 3600 s

```
Step 1: 90 km/h

Step 2–4: Set up conversion fractions so km and h cancel:

    90 km    1000 m    1 h
    ----- × -------- × ------
      h       1 km    3600 s

         ↑             ↑
    km in denominator matches km in numerator (cancels)
    h in denominator matches h in numerator (cancels)

Step 5: Multiply the numbers:
    Numerator:   90 × 1000 × 1 = 90,000
    Denominator: 1 × 1 × 3600 = 3,600

    90,000 / 3,600 = 25

Step 6: Remaining units: m/s

ANSWER: 90 km/h = 25 m/s
```

Quick sanity check: 90 km in 1 hour. That's 90,000 m in 3,600 s. 90,000 ÷ 3,600 = 25. Yes. ✓

**Golden shortcut:** To convert km/h to m/s, always divide by 3.6.
To convert m/s to km/h, multiply by 3.6.

```
km/h ÷ 3.6 = m/s
m/s × 3.6 = km/h

Examples:
72 km/h ÷ 3.6 = 20 m/s
108 km/h ÷ 3.6 = 30 m/s
10 m/s × 3.6 = 36 km/h
```

### Worked Example 2: Convert 5 kg to grams

```
Step 1: 5 kg

Step 2: 1 kg = 1000 g

Step 3–4:
    5 kg × (1000 g / 1 kg)
               ↑
          kg cancels ✓

Step 5: 5 × 1000 = 5000

ANSWER: 5 kg = 5000 g
```

### Worked Example 3: Convert 2.5 hours to seconds

```
Step 1: 2.5 hours

Step 2: 1 hour = 60 minutes; 1 minute = 60 seconds

Step 3–4: Chain two conversion factors:
    2.5 h × (60 min / 1 h) × (60 s / 1 min)
                 ↑                  ↑
            h cancels           min cancels ✓

Step 5: 2.5 × 60 × 60 = 2.5 × 3600 = 9000

ANSWER: 2.5 hours = 9000 seconds
```

### Worked Example 4: Convert 500 mg to kg

```
Step 1: 500 mg (milligrams)

Step 2: 1 mg = 10⁻³ g; 1 g = 10⁻³ kg

Step 3–4:
    500 mg × (10⁻³ g / 1 mg) × (10⁻³ kg / 1 g)
                  ↑                   ↑
             mg cancels           g cancels ✓

Step 5: 500 × 10⁻³ × 10⁻³ kg
      = 500 × 10⁻⁶ kg
      = 5 × 10⁻⁴ kg

ANSWER: 500 mg = 0.0005 kg = 5 × 10⁻⁴ kg
```

### Worked Example 5: Convert 60 m/s² to cm/ms²

This is more advanced but shows how the method scales to complex conversions.

```
Need: m → cm (× 100); s → ms (× 1000, but it's squared so × 1,000,000)

    60 m    100 cm     1 s  ²
    ---- × ------- × (----)
    s²      1 m      1000 ms

= 60 × 100 × (1/1000)² cm/ms²
= 60 × 100 × (1/1,000,000) cm/ms²
= 6000 / 1,000,000 cm/ms²
= 0.006 cm/ms²
= 6 × 10⁻³ cm/ms²
```

---

## 2.6 Scientific Notation

Physics deals with numbers that are either incredibly large or incredibly small. Writing them out in full is impractical and error-prone.

**Scientific notation** is a way of writing any number as:

```
a × 10ⁿ

where:  1 ≤ a < 10   (exactly one non-zero digit before the decimal point)
        n is any whole number (positive, negative, or zero)
```

### Converting TO Scientific Notation

**For large numbers — move the decimal LEFT:**

Count how many places you move the decimal point to the left. That count becomes a positive exponent.

```
6,500,000
         ↑
         Decimal is here (implied after last digit)

Move decimal LEFT 6 places: 6.5 → 6.500000

Result: 6.5 × 10⁶

Verification: 6.5 × 10⁶ = 6.5 × 1,000,000 = 6,500,000 ✓
```

**For small numbers — move the decimal RIGHT:**

Count how many places you move the decimal point to the right. That count becomes a negative exponent.

```
0.000000123
   ↑
   Decimal is here

Move decimal RIGHT 7 places: .000000123 → 1.23

Result: 1.23 × 10⁻⁷

Verification: 1.23 × 10⁻⁷ = 1.23 / 10,000,000 = 0.000000123 ✓
```

### Common Scientific Notation Examples

```
+--------------------------------+---------------------+
|   STANDARD FORM                |  SCIENTIFIC         |
|                                |  NOTATION           |
+--------------------------------+---------------------+
| 299,792,458 m/s                | 3.00 × 10⁸ m/s     |
| (speed of light)               |                     |
+--------------------------------+---------------------+
| 6,022,000,000,000,000,         | 6.022 × 10²³        |
| 000,000,000 particles/mol      | particles/mol       |
| (Avogadro's number)            |                     |
+--------------------------------+---------------------+
| 0.000000000000911 kg           | 9.11 × 10⁻³¹ kg    |
| (mass of an electron)          |                     |
+--------------------------------+---------------------+
| 0.0000000016 C                 | 1.6 × 10⁻¹⁹ C      |
| (charge of a proton)           |                     |
+--------------------------------+---------------------+
| 1,989,000,000,000,000,         | 1.989 × 10³⁰ kg    |
| 000,000,000,000,000,000 kg     |                     |
| (mass of the Sun)              |                     |
+--------------------------------+---------------------+
| 0.000000000000000001 m         | 1 × 10⁻¹⁸ m        |
| (diameter of a quark)          |                     |
+--------------------------------+---------------------+
```

### Arithmetic with Scientific Notation

**Multiplying:**
```
Rule: multiply the 'a' values, ADD the exponents

(a × 10ⁿ) × (b × 10ᵐ) = (a × b) × 10^(n+m)

Example 1: (3 × 10⁴) × (2 × 10⁵)
= (3 × 2) × 10^(4+5)
= 6 × 10⁹

Example 2: (4.5 × 10³) × (2 × 10⁻⁷)
= (4.5 × 2) × 10^(3 + (−7))
= 9.0 × 10⁻⁴

Example 3: (6.0 × 10⁸) × (5.0 × 10⁻³)
= 30.0 × 10⁵
= 3.0 × 10⁶  (adjust so coefficient is between 1 and 10)
```

**Dividing:**
```
Rule: divide the 'a' values, SUBTRACT the exponents

(a × 10ⁿ) / (b × 10ᵐ) = (a / b) × 10^(n−m)

Example 1: (6 × 10⁸) / (2 × 10³)
= (6/2) × 10^(8−3)
= 3 × 10⁵

Example 2: (9 × 10⁻²) / (3 × 10⁴)
= (9/3) × 10^(−2−4)
= 3 × 10⁻⁶
```

**Adding and Subtracting:**
```
Rule: make exponents EQUAL first, then add/subtract coefficients

Example: (3.0 × 10⁴) + (2.0 × 10³)
= (3.0 × 10⁴) + (0.2 × 10⁴)   ← convert to same power
= (3.0 + 0.2) × 10⁴
= 3.2 × 10⁴
```

---

## 2.7 Significant Figures

When you measure something, you can only ever know it to a certain precision. A ruler marked in millimetres cannot tell you something is exactly 5.43217 cm long. **Significant figures** (also called **significant digits**, or just "sig figs") represent how precisely a measurement is actually known.

### The Core Idea

Consider a balance that measures to the nearest gram. You weigh an object and get 47 g. Saying the mass is "47.000000 g" would be dishonest — you don't know all those decimal places. The true reading is 47 g, which has 2 significant figures.

Significant figures are a way of communicating the precision of a measurement in the number itself.

### Rules for Counting Significant Figures

```
RULE 1: All non-zero digits are ALWAYS significant.
        4.23     → 3 sig figs  (4, 2, 3 are all non-zero)
        987      → 3 sig figs
        1,234.5  → 5 sig figs

RULE 2: Zeros BETWEEN non-zero digits are significant.
        ("trapped zeros" cannot be zero — so they must have been measured)
        4.023    → 4 sig figs  (4, 0, 2, 3)
        10.05    → 4 sig figs  (1, 0, 0, 5)
        30,407   → 5 sig figs

RULE 3: LEADING zeros (before any non-zero digit) are NOT significant.
        They only show where the decimal is — they carry no measurement info.
        0.00720  → 3 sig figs  (7, 2, 0 — the leading 00s don't count)
        0.005    → 1 sig fig   (only the 5)
        0.042    → 2 sig figs  (4, 2)

RULE 4: TRAILING zeros AFTER a decimal point ARE significant.
        They were written deliberately to show precision.
        3.140    → 4 sig figs  (3, 1, 4, AND the trailing zero)
        2.500    → 4 sig figs
        10.0     → 3 sig figs
        100.00   → 5 sig figs

RULE 5: Trailing zeros WITHOUT a decimal point are AMBIGUOUS.
        4600 — is it 2 sig figs (measured to nearest 100)?
                    3 sig figs (measured to nearest 10)?
                    4 sig figs (measured exactly)?
        Use scientific notation to be unambiguous:
        4.6 × 10³   → 2 sig figs
        4.60 × 10³  → 3 sig figs
        4.600 × 10³ → 4 sig figs
```

### Worked Examples: Counting Sig Figs

```
+---------------------+--------+--------------------------------------------+
|   NUMBER            | SIG    |   REASONING                                |
|                     | FIGS   |                                            |
+---------------------+--------+--------------------------------------------+
| 3.140               |   4    | 3, 1, 4 are non-zero. Trailing 0 after     |
|                     |        | decimal = significant.                     |
+---------------------+--------+--------------------------------------------+
| 0.00720             |   3    | Leading 00s don't count. 7, 2, and         |
|                     |        | trailing 0 after decimal = significant.    |
+---------------------+--------+--------------------------------------------+
| 12,000              |   2?   | Ambiguous. Use 1.2 × 10⁴ for clarity.     |
+---------------------+--------+--------------------------------------------+
| 1.2000 × 10⁴        |   5    | All trailing zeros after decimal count.    |
+---------------------+--------+--------------------------------------------+
| 0.000300            |   3    | 3, 0, 0 (last two zeros after decimal).    |
+---------------------+--------+--------------------------------------------+
| 100.0               |   4    | Decimal present; all four digits count.    |
+---------------------+--------+--------------------------------------------+
| 304.00              |   5    | 3, 0, 4, 0, 0 — all significant.           |
+---------------------+--------+--------------------------------------------+
| 0.0070050           |   5    | 7, 0, 0, 5, 0 (leading zeros don't count). |
+---------------------+--------+--------------------------------------------+
```

### Sig Figs in Calculations

**Multiplication and Division:**
The answer has the same number of sig figs as the measurement with the FEWEST sig figs.

```
Example: Calculate the area of a rectangle.
Length = 4.52 m     → 3 sig figs
Width  = 3.1 m      → 2 sig figs  (fewest)

Area = 4.52 × 3.1 = 14.012 m²

Round to 2 sig figs: 14 m²

ANSWER: Area = 14 m²   (NOT 14.012 m²)
```

```
Example: Calculate speed.
Distance = 2.50 × 10² m    → 3 sig figs
Time     = 12.5 s          → 3 sig figs

Speed = 250 / 12.5 = 20.0 m/s

Round to 3 sig figs: 20.0 m/s

ANSWER: v = 20.0 m/s
```

**Addition and Subtraction:**
The answer has the same number of DECIMAL PLACES as the measurement with the FEWEST decimal places.

```
Example:
  24.52 m    (2 decimal places)
+  0.3 m     (1 decimal place) ← fewest decimal places
= 24.82 m

Round to 1 decimal place: 24.8 m

ANSWER: 24.8 m
```

### Rounding Rules

```
To round a number, look at the first digit being dropped:
  If it is < 5:  leave the last kept digit unchanged (round down)
  If it is > 5:  increase the last kept digit by 1 (round up)
  If it is = 5:  if followed by non-zero digits, round up
                 if exactly 5 (or 5 followed by zeros), round to even

Examples:
  3.145 → 3.14   (drop 5, previous digit 4 is even, keep 4)
  3.135 → 3.14   (drop 5, previous digit 3 is odd, round up to 4)
  3.146 → 3.15   (drop 6 > 5, round up)
  3.143 → 3.14   (drop 3 < 5, keep 4)
```

---

## 2.8 Accuracy vs Precision

Two words that beginners frequently confuse — but they describe completely different properties of measurements.

**Accuracy** is how close a measurement is to the **true, accepted value**.

**Precision** is how close **repeated measurements** are to each other, regardless of whether they're near the true value.

You can have precision without accuracy, accuracy without precision, neither, or both.

### The Dartboard Analogy

```
TARGET: The bullseye represents the TRUE VALUE of a measurement.
Each dart represents one measurement attempt.

=================================================================
CASE 1: Low Accuracy AND Low Precision
=================================================================

     +---------------------------+
     |                           |
     |     *          *          |
     |                           |
     |  *       +         *      |   + = bullseye
     |                           |
     |              *    *       |
     |                           |
     +---------------------------+

  Darts scattered randomly across the board.
  Neither close to the bullseye nor close to each other.
  This is the worst possible outcome.
  Cause: large random errors AND/OR large systematic errors.

=================================================================
CASE 2: High Precision, LOW Accuracy
=================================================================

     +---------------------------+
     |                           |
     | * *                       |
     | * *                       |   + = bullseye (far away)
     |  *         +              |
     |                           |
     |                           |
     |                           |
     +---------------------------+

  Darts grouped tightly together — but in the WRONG place.
  Consistent results, but consistently WRONG.
  Cause: a SYSTEMATIC ERROR (e.g., a badly calibrated instrument).
  Taking more measurements won't help — they'll all be wrong the same way.

=================================================================
CASE 3: High Accuracy, Low Precision
=================================================================

     +---------------------------+
     |    *                      |
     |                *          |
     |          +                |   + = bullseye
     |      *         *          |
     |            *              |
     |  *                        |
     +---------------------------+

  Darts roughly centred on the bullseye, but scattered.
  The AVERAGE of many measurements is close to the true value,
  but individual measurements vary a lot.
  Cause: large RANDOM errors.
  Taking MORE measurements and averaging HELPS.

=================================================================
CASE 4: High Accuracy AND High Precision (THE GOAL)
=================================================================

     +---------------------------+
     |                           |
     |          * *              |
     |         * + *             |   + = bullseye
     |          * *              |
     |                           |
     |                           |
     +---------------------------+

  Darts tightly clustered right at the bullseye.
  Every individual measurement is close to the true value,
  AND close to every other measurement.
  This is what good experimental physics looks like.
=================================================================
```

### Types of Errors

**Systematic errors** affect accuracy. They push all measurements in the same direction — either all too high or all too low.

Common causes of systematic errors:
- A scale that reads 0.5 kg when nothing is on it (zero error)
- A ruler that has worn away at the end — so zero is not where you think it is
- A thermometer consistently calibrated 2°C too high
- Parallax error — always reading a scale from the same wrong angle

Systematic errors CANNOT be reduced by taking more measurements. You must identify and correct the source of the bias.

**Random errors** affect precision. They cause measurements to scatter around the true value — sometimes reading too high, sometimes too low, with no consistent pattern.

Common causes of random errors:
- Slight vibrations in a laboratory
- Human reaction time when using a stopwatch
- Reading a scale from slightly different positions each time
- Electrical noise in a measuring circuit
- Slight variations in the physical conditions between measurements

Random errors CAN be reduced by taking more measurements and calculating the average. As you take more readings, the random high and low values tend to cancel out.

```
+----------------+--------------------+---------------------+
| TYPE OF ERROR  | AFFECTS            | HOW TO REDUCE       |
+----------------+--------------------+---------------------+
| Systematic     | Accuracy           | Identify and fix    |
|                | (biases all        | the source;         |
|                | readings one way)  | recalibrate         |
+----------------+--------------------+---------------------+
| Random         | Precision          | Take more readings; |
|                | (scatters readings | calculate the mean  |
|                | around true value) |                     |
+----------------+--------------------+---------------------+
```

---

## 2.9 Dimensional Analysis

**Dimensional analysis** is a technique for checking whether equations are physically valid — and sometimes for deriving unknown relationships — purely by tracking the dimensions (types of units) on both sides of an equation.

### Dimensions vs Units

A **dimension** is the fundamental type of a quantity. A **unit** is one particular way of measuring it.

```
QUANTITY    DIMENSION    POSSIBLE UNITS
Length      [L]          metres, centimetres, miles, feet, light-years
Mass        [M]          kilograms, grams, pounds, tonnes
Time        [T]          seconds, minutes, hours, years
```

We write dimensions in square brackets: [L], [M], [T], [I], [Θ], [N], [J].

### The Fundamental Rule of Dimensional Analysis

**Any valid equation must have the same dimensions on both sides.**

This is as absolute as any law of nature. You cannot add mass to time. You cannot equate velocity to force. If the dimensions don't balance, the equation is wrong — no matter how clever the derivation looks.

```
[Left side] = [Right side]    ← MUST be true for any valid equation
```

### Worked Example 1: Verify F = m × a

```
Force = mass × acceleration

Left side:  F   →   [M][L][T]⁻²
            (Newton = kg·m/s²)

Right side: m × a
            m: [M]
            a: [L][T]⁻²  (acceleration = m/s² = length per time²)

            m × a: [M] × [L][T]⁻² = [M][L][T]⁻²

Left side dimensions: [M][L][T]⁻²
Right side dimensions: [M][L][T]⁻²

They match. ✓   F = m × a is dimensionally correct.
```

### Worked Example 2: Verify v = u + at

Where v and u are velocities, a is acceleration, t is time.

```
Left side:  v  →  [L][T]⁻¹  (velocity = m/s)

Right side: u + at
            u:   [L][T]⁻¹              (velocity)
            a×t: [L][T]⁻² × [T] = [L][T]⁻¹   (accel × time = velocity)

            u and at have same dimensions — they can be added ✓
            u + at: [L][T]⁻¹

Left = Right ✓   v = u + at is dimensionally correct.
```

### Worked Example 3: Verify E = ½mv²

```
Left side:  E  →  [M][L]²[T]⁻²  (energy = Joules = kg·m²/s²)

Right side: ½mv²
            ½: dimensionless (pure number)
            m: [M]
            v²: ([L][T]⁻¹)² = [L]²[T]⁻²

            ½mv²: [M] × [L]²[T]⁻² = [M][L]²[T]⁻²

Left side: [M][L]²[T]⁻²
Right side: [M][L]²[T]⁻²

They match. ✓   E = ½mv² is dimensionally correct.
```

### Worked Example 4: Catching a Wrong Equation

Suppose someone claims: "The kinetic energy of an object equals its mass plus its velocity squared."

```
In equation form: E = m + v²

Left side:  E  →  [M][L]²[T]⁻²

Right side: m + v²
            m:  [M]
            v²: [L]²[T]⁻²

            [M] + [L]²[T]⁻²  ← DIFFERENT dimensions!
            You cannot add kilograms to m²/s²

Dimensions do NOT match. ✗   This equation is physically impossible.
```

Dimensional analysis caught the error instantly — no calculation needed.

### Worked Example 5: Deriving an Unknown Relationship

**Question:** The time period T of a simple pendulum depends on its length L and the gravitational acceleration g. What is the relationship?

```
We know:
  T:  [T]
  L:  [L]
  g:  [L][T]⁻²  (g = 9.8 m/s², so dimensions are length/time²)

Assume: T = k × Lᵃ × gᵇ
where k is a dimensionless constant we cannot find from dimensions alone.

Write in dimensions:
  [T]¹ = [L]ᵃ × ([L][T]⁻²)ᵇ
  [T]¹ = [L]^(a+b) × [T]^(−2b)

Match the powers of each dimension:

  Power of [T]:  1 = −2b    →   b = −1/2

  Power of [L]:  0 = a + b  →   a = −b = +1/2

Therefore: T = k × L^(1/2) × g^(−1/2)
           T = k × √(L/g)
           T ∝ √(L/g)

The actual formula is: T = 2π√(L/g)
Dimensional analysis told us the relationship; experiment gives us k = 2π.
```

This is a powerful result: just from knowing which quantities are involved and their dimensions, we derived the form of the equation — without solving a single differential equation.

---

## Summary

- **Units matter absolutely** — the $327 million Mars Climate Orbiter crashed because two engineering teams used different units. Units are part of every measurement, not decorative labels.

- The **SI system** has seven base units: metre (m) for length, kilogram (kg) for mass, second (s) for time, ampere (A) for electric current, kelvin (K) for temperature, mole (mol) for amount of substance, and candela (cd) for luminous intensity. Everything else is derived from these.

- **Derived units** are built by combining base units. Key ones: velocity (m/s), force (Newton N = kg·m/s²), energy (Joule J = kg·m²/s²), power (Watt W = J/s), pressure (Pascal Pa = N/m²).

- **Prefixes** allow compact notation: pico (10⁻¹²), nano (10⁻⁹), micro (10⁻⁶), milli (10⁻³), centi (10⁻²), kilo (10³), mega (10⁶), giga (10⁹), tera (10¹²).

- **Unit conversion** uses the multiply-by-one technique: write the conversion as a fraction equal to 1, arrange it so the unwanted unit cancels. Shortcut: km/h ÷ 3.6 = m/s.

- **Scientific notation** writes numbers as a × 10ⁿ where 1 ≤ a < 10. Multiply by adding exponents; divide by subtracting them.

- **Significant figures** communicate the precision of a measurement. Non-zero digits are always significant. Leading zeros are not. Trailing zeros after a decimal point are significant. Trailing zeros without a decimal point are ambiguous — use scientific notation to clarify.

- **Accuracy** is closeness to the true value; **precision** is consistency of repeated measurements. Systematic errors reduce accuracy; random errors reduce precision. The dartboard analogy shows all four combinations clearly.

- **Dimensional analysis** checks equations by verifying both sides have identical dimensions — [M], [L], [T], etc. It can also derive unknown relationships between physical quantities.

---

## Key Equations

```
UNIT CONVERSION (speed):
    km/h ÷ 3.6 = m/s
    m/s × 3.6 = km/h

TEMPERATURE CONVERSION:
    K = °C + 273.15
    °C = K − 273.15

SCIENTIFIC NOTATION ARITHMETIC:
    Multiply: (a × 10ⁿ) × (b × 10ᵐ) = (a × b) × 10^(n+m)
    Divide:   (a × 10ⁿ) / (b × 10ᵐ) = (a / b) × 10^(n−m)

KEY DERIVED UNIT DEFINITIONS:
    Velocity:  v = d / t              [m/s]
    Accel.:    a = Δv / t             [m/s²]
    Force:     F = m × a              [N = kg·m/s²]
    Energy:    E = F × d              [J = kg·m²/s² = N·m]
    Power:     P = E / t              [W = J/s]
    Pressure:  P = F / A              [Pa = N/m²]

DIMENSIONAL NOTATION:
    [M] = mass dimension
    [L] = length dimension
    [T] = time dimension
    Force:     [M][L][T]⁻²
    Energy:    [M][L]²[T]⁻²
    Velocity:  [L][T]⁻¹
    Accel.:    [L][T]⁻²
    Pressure:  [M][L]⁻¹[T]⁻²
```
