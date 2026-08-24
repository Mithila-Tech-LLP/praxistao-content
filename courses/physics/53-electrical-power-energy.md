# Chapter 53: Electrical Power and Energy

> **"Every time you flip a switch, run a motor, or charge a phone, energy is being converted. Understanding electrical power lets you calculate exactly how much energy is used and at what rate."**

---

## Table of Contents

- [53.1 Electric Power](#531-electric-power)
- [53.2 Power in Resistors](#532-power-in-resistors)
- [53.3 Electrical Energy](#533-electrical-energy)
- [53.4 Efficiency](#534-efficiency)
- [53.5 Power in AC Circuits](#535-power-in-ac-circuits)
- [53.6 Power Transmission and High Voltage](#536-power-transmission-and-high-voltage)
- [53.7 Worked Examples](#537-worked-examples)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 53.1 Electric Power

**Power** is the rate at which energy is transferred or converted:

```
P = W / t = energy / time

Unit: Watt (W) = Joule/second (J/s)
```

For an electric circuit, power is the rate at which electrical energy is converted to other forms (heat, light, mechanical energy, etc.):

```
P = V × I

where:
  P = power (W)
  V = voltage (V)
  I = current (A)
```

### Physical Meaning

```
Current I carries charge past a point.
Voltage V is energy per unit charge.

So: P = energy/time = (energy/charge) × (charge/time) = V × I

In 1 second:
  charge Q = I × 1 = I coulombs pass
  each carries V joules of energy
  total energy = Q × V = I × V joules
  
  Power = I × V joules per second = I × V watts
```

---

## 53.2 Power in Resistors

From Ohm's Law (V = IR), we can write power in three equivalent ways:

```
P = VI

Since V = IR:
  P = (IR) × I = I²R

Since I = V/R:
  P = V × (V/R) = V²/R

All three are equivalent:
  P = VI = I²R = V²/R
```

The most useful form depends on what you know:
- Know I and R: use P = I²R
- Know V and R: use P = V²/R
- Know V and I: use P = VI

### Where Does the Energy Go?

In a resistor, all electrical energy converts to **heat** (thermal energy). The process is called **Joule heating** or **resistive heating**.

This is:
- Useful: electric heaters, toasters, incandescent bulbs, soldering irons
- Unwanted: power line losses, computer chips overheating

### Worked Example 53.1

A 60 W incandescent bulb runs on 240 V.

(a) Find the current through it.
(b) Find its resistance.

**Solution:**

(a) P = VI → I = P/V = 60/240 = **0.25 A**

(b) R = V/I = 240/0.25 = **960 Ω**

Or: R = V²/P = 240²/60 = 57,600/60 = **960 Ω** ✓

---

## 53.3 Electrical Energy

Total energy consumed:

```
E = P × t

where:
  E = energy (J or kWh)
  P = power (W or kW)
  t = time (s or hours)
```

### The Kilowatt-Hour

Electricity bills use **kilowatt-hours (kWh)** not joules (joules are too small):

```
1 kWh = 1000 W × 3600 s = 3.6 × 10⁶ J = 3.6 MJ

A 100 W bulb running for 10 hours uses:
  E = P × t = 0.1 kW × 10 h = 1 kWh

If electricity costs ₹7 per kWh:
  Cost = 1 kWh × ₹7/kWh = ₹7
```

### Power Ratings of Common Devices

| Device | Typical Power |
|--------|--------------|
| LED lamp | 8 W |
| Incandescent bulb | 60 W |
| Desktop computer | 200-400 W |
| Microwave | 900-1200 W |
| Electric kettle | 1500-2000 W |
| Air conditioner | 1000-3500 W |
| Electric car charger | 7000-50,000 W |

### Worked Example 53.2

A household uses the following per day:
- 10 LED bulbs (8 W each) for 5 hours
- Refrigerator (150 W) running continuously (24 hours)
- TV (100 W) for 4 hours
- Air conditioner (1500 W) for 6 hours

Find the total daily energy consumption in kWh, and monthly cost at ₹8/kWh.

**Solution:**

Daily energy:
- Bulbs: 10 × 0.008 × 5 = 0.4 kWh
- Fridge: 0.15 × 24 = 3.6 kWh
- TV: 0.1 × 4 = 0.4 kWh
- AC: 1.5 × 6 = 9 kWh

Total: 0.4 + 3.6 + 0.4 + 9 = **13.4 kWh/day**

Monthly: 13.4 × 30 = 402 kWh

Cost: 402 × ₹8 = **₹3,216 per month**

---

## 53.4 Efficiency

No device converts energy with 100% efficiency. Some energy is always "lost" (converted to heat).

```
Efficiency = useful power output / total power input × 100%

η = P_out / P_in × 100%
```

Or using energy:

```
η = E_out / E_in × 100%
```

### Examples

```
LED bulb:      η ≈ 25-30%  (25% visible light, 75% heat)
Incandescent:  η ≈ 5%      (5% light, 95% heat!)
Electric motor: η ≈ 85-95%
Solar cell:    η ≈ 15-22%
Car engine:    η ≈ 25-35%
Steam turbine: η ≈ 40-45%
```

### Worked Example 53.3

An electric motor draws 500 W from the supply but delivers only 400 W of mechanical power.

(a) Find efficiency.
(b) Where does the remaining power go?

**Solution:**

(a) η = 400/500 × 100% = **80%**

(b) Remaining power = 500 - 400 = 100 W → converted to heat (in motor windings, due to resistance).

---

## 53.5 Power in AC Circuits

Household electricity is **alternating current (AC)**, not DC.

```
AC VOLTAGE:  v(t) = V₀ × sin(ωt)
AC CURRENT:  i(t) = I₀ × sin(ωt)  (for pure resistive load)

where V₀, I₀ are peak (maximum) values.
```

The instantaneous power: p(t) = v(t) × i(t) = V₀I₀ sin²(ωt)

This oscillates between 0 and V₀I₀.

### RMS Values

For power calculations, we use **root mean square (rms)** values:

```
V_rms = V₀ / √2 ≈ 0.707 × V₀

I_rms = I₀ / √2 ≈ 0.707 × I₀

Average power:
P_avg = V_rms × I_rms = (V₀/√2)(I₀/√2) = V₀I₀/2
```

In India, the mains supply is 230 V rms, 50 Hz.

```
V₀ = V_rms × √2 = 230 × 1.414 ≈ 325 V (peak voltage)
```

### Power Factor

For circuits with capacitors or inductors, voltage and current are out of phase:

```
P_avg = V_rms × I_rms × cos(φ)

where φ = phase angle between voltage and current.
cos(φ) = power factor

For pure resistor:  φ = 0, cos(φ) = 1  (maximum power)
For pure inductor:  φ = 90°, cos(φ) = 0  (zero average power)
For pure capacitor: φ = -90°, cos(φ) = 0  (zero average power)
```

---

## 53.6 Power Transmission and High Voltage

Why do power lines use very high voltages (400 kV in India)?

### The Problem: Power Loss in Lines

Transmission lines have resistance R. Power lost to heat:

```
P_loss = I² × R

For the same power P = VI:
  I = P/V

So: P_loss = (P/V)² × R = P²R / V²
```

Higher voltage → lower current → much less power loss!

```
EXAMPLE:
  Transmit 1000 MW over a line with R = 5 Ω:
  
  At 10 kV:
    I = P/V = 10⁹ / 10⁴ = 100,000 A
    P_loss = I²R = 10¹⁰ × 5 = 50,000 MW  (5× MORE than sent! Impossible.)
  
  At 400 kV:
    I = P/V = 10⁹ / 4×10⁵ = 2500 A
    P_loss = I²R = 6.25×10⁶ × 5 = 31.25 MW  (only 3.1% lost)
```

This is why transformers are crucial: they step voltage up for transmission, then step down for home use.

### Step-up and Step-down Transformers

```
GENERATION:    ~11 kV  →  step up  →  TRANSMISSION: 220–400 kV
                                    →  step down →  DISTRIBUTION: 33 kV
                                    →  step down →  LOCAL: 11 kV
                                    →  step down →  HOME: 230 V
```

---

## 53.7 Worked Examples

### Worked Example 53.4 — High Voltage Transmission

A power station generates 100 MW at 11 kV. A step-up transformer increases voltage to 275 kV for transmission.

The transmission line has resistance 4 Ω per conductor (8 Ω total for both conductors).

(a) Find current in the line.
(b) Find power lost in the line.
(c) Find efficiency of transmission.

**Solution:**

(a) P = VI → I = P/V = 100×10⁶ / 275×10³ = **364 A**

(b) P_loss = I²R = 364² × 8 = 132,496 × 8 ≈ **1.06 MW**

(c) η = (100 - 1.06)/100 × 100% = **98.9%** (only 1.1% lost)

### Worked Example 53.5 — Joule Heating

A 3 kW immersion heater is used to heat 10 kg of water from 20°C to 100°C.

(a) How long does it take?
(b) What is the cost at ₹5 per kWh?

Specific heat capacity of water = 4200 J/(kg·°C).

**Solution:**

(a) Energy needed: E = mcΔT = 10 × 4200 × 80 = 3,360,000 J = 3.36 MJ

    Time: t = E/P = 3,360,000 / 3000 = **1120 s ≈ 18.7 minutes**

(b) Energy in kWh = 3.36 MJ / 3.6 MJ/kWh ≈ 0.933 kWh

    Cost = 0.933 × ₹5 = **₹4.67**

---

## Summary

- **Power**: P = VI = I²R = V²/R; unit Watts (W = J/s)
- In a resistor: all power → heat (Joule heating)
- **Energy**: E = Pt; unit kWh (1 kWh = 3.6 MJ)
- **Efficiency**: η = P_out/P_in × 100%
- **AC power**: use rms values; V_rms = V₀/√2; P_avg = V_rms × I_rms × cos(φ)
- **High voltage transmission**: reduces current, reduces I²R losses dramatically
- Power loss ∝ I², so doubling voltage reduces power loss by factor 4

---

## Key Equations

```
Power:
  P = V × I
  P = I² × R     (power dissipated in resistor)
  P = V² / R     (power dissipated in resistor)

Energy:
  E = P × t
  1 kWh = 3.6 × 10⁶ J

Efficiency:
  η = P_out / P_in × 100%

AC rms values:
  V_rms = V₀ / √2
  I_rms = I₀ / √2

AC average power:
  P_avg = V_rms × I_rms × cos(φ)
  (For pure resistor: cos(φ) = 1)

Power line losses:
  P_loss = I²R = (P/V)²R
  Higher V → less loss
```
