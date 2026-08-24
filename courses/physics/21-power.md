# Chapter 21: Power

> "It is not enough to do the work — what matters is how quickly you do it."
> — Engineering proverb

---

## Table of Contents

- [Introduction: What Is Power?](#introduction-what-is-power)
- [The Definition of Power](#the-definition-of-power)
  - [Power as Rate of Doing Work](#power-as-rate-of-doing-work)
  - [The Watt: The Unit of Power](#the-watt-the-unit-of-power)
  - [Horsepower: The Old Unit](#horsepower-the-old-unit)
- [Calculating Power: P = W / t](#calculating-power-p--w--t)
- [The Relationship P = F × v](#the-relationship-p--f--v)
  - [Derivation of P = F × v](#derivation-of-p--f--v)
  - [When to Use P = F × v](#when-to-use-p--f--v)
- [Power in Everyday Life](#power-in-everyday-life)
  - [A Large Power Comparison Table](#a-large-power-comparison-table)
  - [The Human Body as a Machine](#the-human-body-as-a-machine)
  - [Car Engines](#car-engines)
- [Efficiency: How Much Power Is Useful?](#efficiency-how-much-power-is-useful)
  - [The Efficiency Formula](#the-efficiency-formula)
  - [Why No Machine Is 100% Efficient](#why-no-machine-is-100-efficient)
  - [Efficiency of Common Devices](#efficiency-of-common-devices)
- [Electricity Bills and Kilowatt-Hours](#electricity-bills-and-kilowatt-hours)
  - [What Is a Kilowatt-Hour?](#what-is-a-kilowatt-hour)
  - [How to Calculate Your Electricity Bill](#how-to-calculate-your-electricity-bill)
- [Worked Examples](#worked-examples)
  - [Example 1: Kettle Power](#example-1-kettle-power)
  - [Example 2: Car Engine on the Motorway](#example-2-car-engine-on-the-motorway)
  - [Example 3: Human Climbing Stairs](#example-3-human-climbing-stairs)
  - [Example 4: Efficiency of a Light Bulb](#example-4-efficiency-of-a-light-bulb)
- [Common Misconceptions](#common-misconceptions)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: What Is Power?

Here is a puzzle. Two people are moving an identical pile of boxes from a delivery truck into a warehouse. Both of them move exactly the same number of boxes, and each box weighs exactly the same. One person takes 10 minutes; the other takes 30 minutes.

Did they both do the same amount of work? In physics, **yes** — they moved the same mass through the same distance against gravity, so they did the same total work.

But is there a difference? Clearly, yes! The first person worked three times faster. They got the job done in one third of the time. This idea — how fast work gets done, or how fast energy is transferred — is what physicists call **power**.

Power is an absolutely essential concept in engineering, technology, and everyday life. When you buy a car, the power (horsepower or kilowatts) tells you how fast the engine can do work — which determines acceleration. When you buy an appliance, the power (in watts) tells you how fast it consumes electrical energy — which determines your electricity bill. When a doctor tests your fitness, power tells them how efficiently your body converts food energy into motion.

In this chapter, you will learn:
- The precise definition of power and its formula
- The unit of power (Watts) and how it relates to Horsepower
- How to calculate the power needed to move a vehicle at constant speed
- What efficiency means and why nothing is ever 100% efficient
- How to read and understand your electricity bill

---

## The Definition of Power

### Power as Rate of Doing Work

**Power** is defined as the **rate of doing work** or the **rate of transferring energy**.

In plain English: power measures how much work is done per second (or per any unit of time).

A powerful machine does a lot of work in a short time. A less powerful machine does the same work but takes longer.

```
UNDERSTANDING POWER
====================

Low power machine:
[=====work done=====]      takes 10 seconds
|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|
  t=1   t=2   t=3   t=4   t=5   t=6   t=7   t=8   t=9  t=10

High power machine:
[=====same work done=====]  takes 2 seconds
|-----|-----|
  t=1   t=2

Both machines did the same WORK.
The high power machine transferred energy 5x faster.
```

### The Watt: The Unit of Power

The **Watt** (abbreviated **W**) is the SI unit of power, named after Scottish engineer James Watt (1736–1819), who greatly improved the steam engine.

```
1 Watt = 1 Joule per second = 1 J/s
```

So a machine running at 1 Watt is doing 1 Joule of work every second.

Multiples of the Watt:

```
WATT MULTIPLES
===============

1 milliwatt (mW)   = 0.001 W       = 10⁻³ W
1 Watt (W)         = 1 W
1 kilowatt (kW)    = 1,000 W       = 10³ W
1 megawatt (MW)    = 1,000,000 W   = 10⁶ W
1 gigawatt (GW)    = 10⁹ W

A power station generates about 1–3 GW.
A kettle uses about 2 kW.
A phone charger uses about 5–20 W.
A LED light bulb uses about 10 W.
```

### Horsepower: The Old Unit

**Horsepower (hp)** is an older unit of power, developed by James Watt himself to compare the output of his steam engines to the power of draft horses (to convince customers his engines were worth buying).

```
1 horsepower (hp) = 746 Watts ≈ 750 W
```

So:
- A 100 hp car engine = 100 × 746 = 74,600 W ≈ 74.6 kW
- A 200 hp car engine = 200 × 746 = 149,200 W ≈ 149 kW

Car manufacturers still use horsepower in many countries (particularly the US), but the rest of the world uses kilowatts. You can convert between them easily:

```
hp → kW:  multiply by 0.746
kW → hp:  multiply by 1.341
```

---

## Calculating Power: P = W / t

The formula for power is straightforward:

```
P = W / t

Where:
  P = power, measured in Watts (W)
  W = work done (or energy transferred), measured in Joules (J)
  t = time taken, measured in seconds (s)
```

This can also be written as:

```
P = Energy / time
P = E / t
```

Rearrangements of this formula:

```
P = W / t     --> find power given work and time
W = P × t     --> find work (energy) given power and time
t = W / P     --> find time given work and power
```

**Simple example:** A motor lifts a 10 kg box to a height of 3 m in 5 seconds.

First, find the work done:
```
W = F × d = m × g × h = 10 × 9.8 × 3 = 294 J
```

Then find the power:
```
P = W / t = 294 / 5 = 58.8 W
```

The motor operates at about 59 Watts.

---

## The Relationship P = F × v

There is another extremely useful form of the power equation that connects power to force and velocity (speed).

### Derivation of P = F × v

Start with P = W / t.

Work is defined as W = F × d (force times distance).

Substitute:
```
P = W / t = (F × d) / t
```

But d / t is just speed (v = distance / time):

```
P = F × (d / t) = F × v

Therefore:
P = F × v

Where:
  P = power (Watts, W)
  F = force applied (Newtons, N)
  v = speed (velocity) of the object (m/s)
```

This formula is used when a force is applied to an object that is moving at a constant speed.

### When to Use P = F × v

This formula is especially useful when:
- An object is moving at **constant speed** (forces are balanced)
- You know the **driving force** and the **speed**
- You want to find the engine power needed for a vehicle

**Example:** A car moves at a constant 20 m/s. The driving force from the engine is 2,000 N. What power is the engine producing?

```
P = F × v = 2,000 × 20 = 40,000 W = 40 kW
```

**Why constant speed?** When the car moves at constant speed, acceleration is zero, so the net force is zero. The engine force equals the drag force (friction + air resistance). In this situation, all the engine's power goes into overcoming drag — the energy is transferred directly into heat in the air and tyres rather than gaining kinetic energy.

---

## Power in Everyday Life

### A Large Power Comparison Table

This table shows the power of common devices, processes, and phenomena. Study it to develop intuition for power values.

```
POWER COMPARISON TABLE
=======================

DEVICE / PROCESS                   POWER (Watts)
------------------------------------------------
LED light bulb (modern)            5 – 15 W
Smartphone (screen on)             3 – 6 W
Smartphone charger                 5 – 20 W
Laptop computer (in use)           20 – 60 W
Desktop computer + monitor         100 – 300 W
Incandescent light bulb            40 – 100 W
Human brain (thinking hard)        ~20 W
Human body (resting)               ~80 W
Human body (walking)               ~250 W
Human body (cycling moderately)    ~100 – 200 W
Human body (maximum effort sprint) ~500 – 1,000 W
Professional cyclist (1 hour)      ~300 – 400 W
Kitchen kettle                     1,500 – 3,000 W  (1.5 – 3 kW)
Microwave oven                     600 – 1,200 W
Electric kettle (UK, 230V)         ~2,200 W  (2.2 kW)
Hair dryer                         1,000 – 2,000 W
Toaster                            800 – 1,500 W
Washing machine                    2,000 – 2,500 W
Electric oven                      2,000 – 5,000 W
Family car engine (petrol)         60,000 – 150,000 W (60–150 kW)
Sports car engine                  200,000 – 500,000 W (200–500 kW)
Formula 1 car                      ~750,000 W (750 kW ≈ 1,000 hp)
Commercial jet engine (each)       ~30,000,000 W (30 MW)
Large wind turbine                 2,000,000 – 5,000,000 W (2–5 MW)
Nuclear power station              1,000,000,000 W (1 GW)
Lightning bolt (peak)              ~1,000,000,000,000 W (1 TW) but for only microseconds
Sun's total output                 ~3.8 × 10²⁶ W (incomprehensible)
------------------------------------------------
```

Notice the enormous range — from a few watts for a light bulb to gigawatts for a power station. This is why we use the kilo/mega/giga prefixes.

### The Human Body as a Machine

Your body is a remarkably versatile machine, and understanding its power output gives a sense of scale.

```
HUMAN POWER OUTPUT
===================

Activity                 Approximate Power Output
-------------------------------------------------
Sleeping                 ~70 W
Sitting quietly          ~80 W
Light office work        ~100 W
Standing, light work     ~120 W
Walking at normal pace   ~200 – 300 W
Cycling (leisurely)      ~100 – 150 W
Cycling (competitive)    ~250 – 400 W
Swimming                 ~200 – 400 W
Running (jogging)        ~300 – 500 W
Running (sprinting)      ~700 – 1,000 W
Maximum brief effort     ~1,500 – 2,000 W (for a second or two)
-------------------------------------------------

Your body is about as powerful as a bright light bulb
when resting, and as powerful as a small appliance
when exercising moderately.

The food you eat provides the energy for all this.
An average person uses ~2,000 food Calories per day
= 2,000 × 4,186 J = 8,372,000 J ≈ 8.4 MJ

Average power = 8,400,000 J / 86,400 s (seconds in a day)
             ≈ 97 W  (roughly 100 W average, all day)
```

### Car Engines

When a car is moving at constant speed on a motorway, the engine is not accelerating the car — it is just overcoming air resistance and friction. The required engine power depends on the car's speed:

```
REQUIRED ENGINE POWER vs SPEED (typical family car)
====================================================

At 30 mph (13 m/s):     ~4 kW    = 5 hp
At 60 mph (27 m/s):     ~18 kW   = 24 hp
At 70 mph (31 m/s):     ~27 kW   = 37 hp
At 100 mph (45 m/s):    ~70 kW   = 94 hp
At 130 mph (58 m/s):    ~150 kW  = 200 hp

Why does power increase so dramatically?
Because air resistance force grows with the SQUARE of speed,
and power = F × v, so power grows with the CUBE of speed!

Going from 60 mph to 120 mph uses ~8x more power, not 2x.
This is why high speeds are so fuel-inefficient.
```

---

## Efficiency: How Much Power Is Useful?

No machine converts all of its input energy into useful output energy. Some always escapes as heat, sound, or other unwanted forms. The fraction that is useful is called **efficiency**.

### The Efficiency Formula

```
Efficiency = (Useful output energy / Total input energy) × 100%

Or equivalently (since P = E/t, and time is the same for both):

Efficiency = (Useful output power / Total input power) × 100%
```

Efficiency is measured as a **percentage** (%). The maximum possible efficiency is 100% (perfect machine — nothing wasted). In reality, all machines are below 100%.

**Example:** An electric motor uses 500 J of electrical energy and produces 400 J of useful mechanical work. The other 100 J becomes heat. What is its efficiency?

```
Efficiency = (400 / 500) × 100%
           = 80%
```

This is a reasonably good motor. Only 20% was wasted.

### Why No Machine Is 100% Efficient

The second law of thermodynamics (which you will study in a later chapter) states that whenever energy is converted from one form to another, some energy always escapes as heat that cannot be fully recovered. This is a fundamental law of nature, not an engineering limitation.

Reasons energy is lost:
- **Friction**: surfaces sliding against each other warm up
- **Air resistance**: objects moving through air push air molecules and heat them
- **Electrical resistance**: electric current flowing through wires generates heat
- **Sound**: vibrating parts create noise (sound carries away energy)
- **Radiation**: hot objects emit infrared radiation (heat radiated away)

```
ENERGY LOSSES IN A PETROL CAR ENGINE
======================================

100% input energy (from burning petrol)
              |
              v
   [Petrol burns in cylinders]
              |
    __________|___________
   |          |           |
  ~30%      ~35%        ~35%
Useful     Exhaust      Cooling
mechanical  heat         water
 energy    (out of      heat
(to wheels) exhaust     (engine
            pipe)       stays cool)
              |
              v
   Of the 30% mechanical energy:
     - Some lost to transmission friction
     - Some lost to tyre rolling resistance
     - Some lost to air resistance (at speed)
     - Only ~20-25% actually propels the car

Overall efficiency of a petrol engine: about 20-30%
(meaning 70-80% of petrol's energy is wasted as heat!)
```

### Efficiency of Common Devices

```
EFFICIENCY OF COMMON DEVICES
==============================

Device                        Efficiency     Main waste
-----------------------------------------------------
LED light bulb                ~95%           5% heat
Electric motor (good)         85–95%         5–15% heat
Electric kettle               ~98%           2% heat to air  *
Hydroelectric power station   ~85–90%        turbine friction
Solar panel (silicon)         15–23%         77–85% becomes heat
Petrol car engine             20–30%         70–80% heat/exhaust
Diesel car engine             35–45%         55–65% heat/exhaust
Steam power station           ~35%           ~65% cooling water
Human muscle                  ~25%           ~75% body heat
Fluorescent tube light        ~25%           ~75% heat
Incandescent light bulb       ~5%            ~95% heat (very poor!)
-----------------------------------------------------

* The kettle is nearly 100% efficient at HEATING WATER, but
  the power station that generates the electricity is only ~35%
  efficient. So the overall chain is about 34% efficient.
```

**Important note on the incandescent bulb:** It is 5% efficient at producing *light*. But 95% of its energy becomes heat. It is actually 100% efficient at converting electrical energy to *energy of some kind* — it just gives mostly the wrong kind (heat instead of light).

---

## Electricity Bills and Kilowatt-Hours

### What Is a Kilowatt-Hour?

When you look at your electricity bill, you do not see "Joules used" — you see **kilowatt-hours (kWh)**.

A **kilowatt-hour** is a unit of energy (not power). It is the amount of energy used by a 1-kilowatt device running for 1 hour.

```
1 kWh = 1 kW × 1 hour
      = 1,000 W × 3,600 s
      = 3,600,000 J
      = 3.6 × 10⁶ J = 3.6 MJ
```

So one kWh is 3.6 million Joules. That is a lot of energy — it is why a kWh is a more convenient unit than a Joule for everyday electricity.

**Why use kWh instead of Joules?**
Because Joules are tiny for household purposes. Running a kettle for one minute uses about 132,000 J. Saying "you used 0.037 kWh" is much more manageable for billing.

### How to Calculate Your Electricity Bill

```
ELECTRICITY BILL CALCULATION
==============================

Step 1: Find the device's power in kW.
  Power = 2,200 W = 2.2 kW (kettle)

Step 2: Find how many hours per day it runs.
  Kettle: about 10 minutes (0.167 hours) per day

Step 3: Calculate energy used per day.
  Energy = Power × Time = 2.2 kW × 0.167 h = 0.367 kWh per day

Step 4: Scale to a month (30 days).
  Monthly energy = 0.367 × 30 = 11 kWh

Step 5: Multiply by price per kWh.
  UK price ≈ 24p per kWh (varies)
  Monthly cost = 11 × 24p = 264p = £2.64 per month for the kettle

TOTAL HOME ELECTRICITY (approximate UK averages)
-------------------------------------------------
Small flat:      ~1,500 kWh per year  ≈ £360/year
Average house:   ~3,500 kWh per year  ≈ £840/year
Large house:     ~6,000 kWh per year  ≈ £1,440/year
```

**Biggest energy users in a typical home:**

```
HOME ENERGY CONSUMERS
======================

Appliance           Power    Hours/day   kWh/day   Approx Cost/month
---------------------------------------------------------------------
Electric heating    ~2 kW    8 h/day     16 kWh    £115 (major cost!)
Fridge/freezer      150 W    24 h        3.6 kWh   £26 (always on)
Washing machine     2 kW     0.5 h       1 kWh     £7
Tumble dryer        2.5 kW   0.5 h       1.25 kWh  £9
Oven                2 kW     0.5 h       1 kWh     £7
Kettle              2.2 kW   0.2 h       0.44 kWh  £3
TV (large screen)   100 W    4 h         0.4 kWh   £3
LED lighting        50 W     5 h         0.25 kWh  £2
Phone charging      10 W     2 h         0.02 kWh  £0.14
---------------------------------------------------------------------
* Prices approximate based on 24p/kWh
```

Heating dominates home energy costs. This is why insulation, heat pumps, and efficient boilers matter so much for energy bills.

---

## Worked Examples

### Example 1: Kettle Power

**Question:** A kettle uses 2,200 J of electrical energy every second to boil water. 
(a) What is the power of the kettle in Watts?
(b) What is it in kilowatts?
(c) How much energy (in kWh) does it use if you boil it for 2 minutes?

**Solution (a):**

Power is energy per second:
```
P = Energy / time = 2,200 J / 1 s = 2,200 W
```

**Solution (b):**
```
P = 2,200 W ÷ 1,000 = 2.2 kW
```

**Solution (c):**

First convert time: 2 minutes = 2/60 hours = 0.0333 hours

```
Energy = Power × time
Energy = 2.2 kW × 0.0333 h
Energy = 0.0733 kWh
```

Or equivalently in Joules:
```
Energy = 2,200 W × 120 s = 264,000 J = 264 kJ
```

**Answer:** The kettle uses 0.073 kWh (or 264,000 J) when boiled for 2 minutes.

**Cost check:** At 24p per kWh, this boiling costs: 0.073 × 24p ≈ 1.75p per boil. Cheap! But if you boil it 10 times a day for a year: 10 × 1.75p × 365 = £63.88/year. Less cheap.

---

### Example 2: Car Engine on the Motorway

**Question:** A car of mass 1,400 kg is travelling at a constant speed of 30 m/s on a motorway. The total drag force (air resistance + friction) is 800 N.

(a) What is the engine's output power?
(b) Convert this to horsepower.
(c) If the car's engine is 30% efficient, what is the rate of fuel energy consumption?

**Solution (a):**

At constant speed, engine force = drag force = 800 N.

Using P = F × v:
```
P = F × v = 800 × 30 = 24,000 W = 24 kW
```

**Solution (b):**
```
24,000 W ÷ 746 = 32.2 hp
```

So only about 32 horsepower is needed for constant motorway driving. The car may have 100 hp or more — the extra power is for acceleration and overtaking.

**Solution (c):**

Efficiency = 30%, meaning only 30% of fuel energy becomes useful engine output.

If output power is 24 kW and this is 30% of input:
```
Input power = Output power / Efficiency
Input power = 24,000 / 0.30
Input power = 80,000 W = 80 kW
```

The engine consumes fuel at a rate equivalent to 80 kW — but only 24 kW of that actually moves the car. The other 56 kW (70%) goes out as exhaust heat and engine cooling.

---

### Example 3: Human Climbing Stairs

**Question:** A person of mass 60 kg climbs a staircase that rises 4 m in 8 seconds.

(a) How much work do they do against gravity?
(b) What is their power output (in Watts)?
(c) How does this compare to a 60-watt light bulb?
(d) If their muscles are 25% efficient, how much chemical energy (from food) do they use?

**Solution (a):**

Work done against gravity = gain in PE:
```
W = m × g × h = 60 × 9.8 × 4 = 2,352 J
```

**Solution (b):**
```
P = W / t = 2,352 / 8 = 294 W
```

**Solution (c):**

294 W vs 60 W light bulb: the person is producing about **5 times** the power of a 60-watt bulb. If you could convert all that mechanical power to light, you could power five light bulbs just by climbing stairs!

**Solution (d):**

Muscle efficiency = 25%, so:
```
Chemical energy used = Useful mechanical energy / efficiency
Chemical energy used = 2,352 / 0.25
Chemical energy used = 9,408 J
```

Or the chemical energy consumption rate:
```
Chemical power input = 294 / 0.25 = 1,176 W
```

So while producing 294 W mechanically, the body is burning food energy at a rate of about 1,176 W. The other 882 W heats your body — which is why you get warm when exercising!

---

### Example 4: Efficiency of a Light Bulb

**Question:** A traditional incandescent light bulb is rated at 60 W.
- It produces 3 W of useful light energy.
- The rest becomes heat.

(a) How much power is wasted as heat?
(b) What is its efficiency?
(c) An LED replacement bulb produces the same amount of light (3 W) using only 5 W of electrical power. What is its efficiency?
(d) If the bulb is on for 6 hours per day, calculate the yearly electricity cost for each bulb at 24p/kWh.

**Solution (a):**
```
Power wasted as heat = 60 - 3 = 57 W
```

**Solution (b): Incandescent efficiency:**
```
Efficiency = (Useful output / Total input) × 100%
           = (3 / 60) × 100%
           = 5%

Only 5% of the electricity becomes light!
95% becomes heat.
```

**Solution (c): LED efficiency:**
```
Efficiency = (3 / 5) × 100% = 60%

The LED is 60% efficient at producing light.
It produces the same light as the incandescent bulb
but uses 12 times less electricity.
```

**Solution (d): Annual electricity cost:**

Incandescent:
```
Hours per year = 6 h/day × 365 days = 2,190 h
Energy = 60 W × 2,190 h = 131,400 Wh = 131.4 kWh
Cost = 131.4 × 24p = 3,154p = £31.54 per year
```

LED:
```
Energy = 5 W × 2,190 h = 10,950 Wh = 10.95 kWh
Cost = 10.95 × 24p = 263p = £2.63 per year
```

**Savings: £31.54 - £2.63 = £28.91 per year per bulb.**

A typical home has 20+ bulbs. Switching all to LEDs saves hundreds of pounds per year. This is why governments have banned incandescent bulbs.

```
VISUAL COMPARISON
==================

Incandescent 60W:
INPUT: [#################################### 60 W ]
         |                                       |
     3 W light                              57 W heat
     [##]                          [##################################]
     (5%)                                   (95%)

LED 5W:
INPUT: [#####  5 W  ]
         |           |
     3 W light    2 W heat
     [###]        [##]
     (60%)        (40%)

Same light output. 12x less power.
```

---

## Common Misconceptions

**Misconception 1:** "Power and energy are the same thing."
**Correction:** Power is the *rate* of energy transfer. Energy is the total amount. Power is energy divided by time. A powerful engine can do a lot of work *quickly*. An energy-rich battery stores a lot of energy but may deliver it slowly (low power) or quickly (high power).

**Misconception 2:** "A high-power appliance always costs more to run."
**Correction:** Running cost depends on energy used, which is power × time. A 2,000 W kettle running for 2 minutes uses less energy than a 100 W light bulb running for 8 hours. Check both power AND time.

**Misconception 3:** "Efficient machines are the most powerful."
**Correction:** Efficiency and power are completely different. A very efficient machine wastes little energy — but it might be powerful or weak. A very efficient machine converts nearly all input energy to useful output, regardless of how much that output is.

**Misconception 4:** "Cars need maximum engine power to drive at constant speed."
**Correction:** Cars need maximum power for rapid acceleration. At constant speed, the only power needed is to overcome drag. A 100 hp car driving steadily at 60 mph may only be using 20–30 hp. The rest sits unused.

**Misconception 5:** "100% efficiency is theoretically possible with better engineering."
**Correction:** The second law of thermodynamics proves that 100% efficiency is impossible for heat engines (engines that work because of a temperature difference). Some other machines (like an ideal electric motor) can theoretically approach 100%, but in practice, every real machine has some losses.

---

## Summary

- **Power** is the rate of doing work or transferring energy.

- The formula for power is:
  - P = W / t  (work done divided by time)
  - P = E / t  (energy transferred divided by time)

- The unit of power is the **Watt (W)**:
  - 1 W = 1 J/s (one Joule per second)
  - 1 kW = 1,000 W
  - 1 MW = 1,000,000 W

- **Horsepower** is an older unit: 1 hp = 746 W ≈ 750 W

- For an object moving at constant speed: **P = F × v** (power = force × velocity). This is derived from P = W/t and W = F × d.

- **Efficiency** measures how much input energy becomes useful output:
  - Efficiency = (Useful output / Total input) × 100%
  - No machine reaches 100% — some energy always escapes as heat

- **Kilowatt-hours (kWh)** are the unit on electricity bills:
  - 1 kWh = 1,000 W used for 1 hour = 3,600,000 J
  - Cost = Power (kW) × Time (hours) × Price per kWh

- A human body at rest outputs about **80–100 W** — similar to a bright light bulb

- A family car engine on the motorway uses about **20–30 kW** at constant speed

- A traditional kettle uses about **2–3 kW**, making it one of the highest-power common appliances

---

## Key Equations

```
POWER (basic)
  P = W / t
  P = E / t
  P in Watts (W)
  W or E in Joules (J)
  t in seconds (s)

POWER (for moving objects)
  P = F × v
  P in Watts (W)
  F in Newtons (N)
  v in metres per second (m/s)
  (valid when force and velocity are in the same direction)

REARRANGEMENTS
  W = P × t   (work/energy from power and time)
  t = W / P   (time from work and power)
  v = P / F   (speed from power and force)
  F = P / v   (force from power and speed)

EFFICIENCY
  Efficiency = (Useful output energy / Total input energy) × 100%
  Efficiency = (Useful output power / Total input power) × 100%
  Efficiency is between 0% and 100%

ENERGY FROM POWER AND TIME (for electricity bills)
  E (kWh) = P (kW) × t (hours)
  Cost = E (kWh) × price per kWh

UNIT CONVERSIONS
  1 W = 1 J/s
  1 kW = 1,000 W
  1 MW = 1,000,000 W
  1 hp = 746 W ≈ 750 W
  1 kWh = 3,600,000 J = 3.6 MJ
```

---

*Next chapter: Momentum and Impulse — why a small fast object can be as difficult to stop as a large slow object, and how car safety features use physics to save lives.*
