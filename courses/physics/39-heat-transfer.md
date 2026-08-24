# Chapter 39: Heat Transfer

> "Heat can be transferred in three ways — like a rumour spreading: by direct contact, by mass movement, or by shouting across a room."
> — Anonymous physics teacher

---

## Table of Contents

1. [Three Ways Heat Travels](#three-ways-heat-travels)
2. [Mechanism 1: Conduction](#mechanism-1-conduction)
3. [Fourier's Law of Conduction](#fouriers-law)
4. [Thermal Conductivity Values](#thermal-conductivity-values)
5. [Conductors and Insulators](#conductors-and-insulators)
6. [Double-Pane Windows](#double-pane-windows)
7. [Worked Example 1: Heat Loss Through a Wall](#worked-example-1)
8. [Mechanism 2: Convection](#mechanism-2-convection)
9. [Natural and Forced Convection](#natural-and-forced-convection)
10. [Convection in Nature](#convection-in-nature)
11. [How Radiators Heat Rooms](#how-radiators-heat-rooms)
12. [Worked Example 2: Convection Currents in a Kettle](#worked-example-2)
13. [Mechanism 3: Radiation](#mechanism-3-radiation)
14. [Stefan-Boltzmann Law](#stefan-boltzmann-law)
15. [Dark vs Light Surfaces](#dark-vs-light-surfaces)
16. [The Greenhouse Effect](#the-greenhouse-effect)
17. [Worked Example 3: Radiation from a Human Body](#worked-example-3)
18. [The Thermos Flask](#the-thermos-flask)
19. [Summary](#summary)
20. [Key Equations](#key-equations)

---

## Three Ways Heat Travels

In Chapter 37 we learned that heat is the transfer of thermal energy. But exactly how does that energy travel from one place to another?

There are three distinct mechanisms:

```
┌──────────────────────────────────────────────────────────────┐
│                   THREE HEAT TRANSFER MECHANISMS              │
│                                                              │
│   CONDUCTION         CONVECTION          RADIATION           │
│                                                              │
│  ████░░░░░░░░      ↑↑  ↑↑  ↑↑          )))))))            │
│  ████░░░░░░░░    ↑↑  ↑↑  ↑↑           )))))))            │
│  HOT    COLD       warm fluid           )))))))            │
│   |  →  |          rises               HEAT (infrared)     │
│  vibration       cool fluid             travels through    │
│  passes           sinks                 empty space        │
│  along                                                      │
│                                                              │
│  Needs solid      Needs fluid          Needs nothing        │
│  or fluid         (liquid or gas)      (works in vacuum)    │
└──────────────────────────────────────────────────────────────┘
```

Each mechanism works differently and dominates in different situations. Understanding all three is essential for understanding insulation, weather, heating systems, cooking, and even the energy balance of the entire Earth.

---

## Mechanism 1: Conduction

**Conduction** is the transfer of heat through a material by the passing of vibrational energy from atom to atom, without any bulk movement of the material itself.

Imagine a long metal rod with one end in a fire. The atoms at the hot end vibrate violently. They collide with their neighbours, passing some energy along. Those neighbours pass energy to their neighbours. Gradually, the vibration (and thus the heat) travels along the rod.

```
Heat Conduction Along a Metal Rod:

  HOT end                                           COLD end
  (fire)                                            (in air)
  
  ████████████████████████████████████████████████████████
  ←←← heat flows this way ←←←←←←←←←←←←←←←←←←←←←←←←←
     
  Temperature profile along the rod:
  
  Temp
   |
  T₁|█
   | ██
   |  ███
   |    ████
   |       █████
   |           ██████
  T₂|               ████████████████
   |
   +──────────────────────────────────→ position along rod
   hot end                           cold end
```

In **metals**, there is an additional mechanism: free electrons. Metals have electrons that are not tightly bound to atoms and can move freely through the material. These free electrons are excellent energy carriers — they absorb energy at the hot end, travel quickly through the metal, and deposit energy at the cooler end. This is why metals are such good conductors of both heat and electricity.

In **non-metals** (wood, plastic, brick, glass), there are no free electrons. Heat can only travel via atomic vibrations (phonons), which is slower and less efficient. These materials are poor conductors — they are **insulators**.

In **gases**, atoms are far apart and rarely collide. Conduction through gases is very poor.

---

## Fourier's Law

The rate of heat flow by conduction through a material depends on:
- How large the temperature difference is (ΔT)
- How thick the material is (L) — thicker means slower
- How large the cross-sectional area is (A) — larger area means more heat flow
- What the material is made of (thermal conductivity k)

**Fourier's Law of Heat Conduction:**

```
Q/t = k × A × ΔT / L

Where:
  Q/t = rate of heat flow (watts = joules per second)
  k   = thermal conductivity of the material (W/m·°C or W/m·K)
  A   = cross-sectional area perpendicular to heat flow (m²)
  ΔT  = temperature difference across the material (°C or K)
  L   = thickness of the material (m)
```

This formula tells you: heat flows faster through a material that has high k (good conductor), large area, large temperature difference, and small thickness.

```
Fourier's Law Illustrated:

        Area A
        ┌───────────────────┐
        │                   │
  T₁   │   material with   │   T₂
(hot)  │   thickness L     │  (cold)
        │   conductivity k  │
        │                   │
        └───────────────────┘
        
        Heat flow rate = k × A × (T₁ - T₂) / L
        Direction: from hot (T₁) to cold (T₂)
```

---

## Thermal Conductivity Values

```
Material                 | k (W/m·°C) | Category
-------------------------|------------|----------
Silver                   |  429       | Excellent conductor
Copper                   |  385       | Excellent conductor
Gold                     |  315       | Excellent conductor
Aluminium                |  205       | Very good conductor
Steel (carbon)           |   50       | Good conductor
Lead                     |   35       | Moderate conductor
Stainless steel          |   16       | Moderate conductor
-------------------------|------------|----------
Glass                    |    1.0     | Poor conductor
Water (liquid)           |    0.60    | Poor conductor
Brick                    |    0.72    | Poor conductor
Concrete                 |    1.4     | Poor conductor
-------------------------|------------|----------
Wood (oak)               |    0.17    | Insulator
Wool / felt              |    0.04    | Good insulator
Styrofoam                |    0.033   | Excellent insulator
Fiberglass insulation    |    0.044   | Excellent insulator
Air (still)              |    0.025   | Excellent insulator
Aerogel                  |    0.015   | Best known insulator (solid)
-------------------------|------------|----------
Human tissue             |    0.50    | Poor conductor
```

Notice that **still air** is an excellent insulator (k = 0.025 W/m·°C). This is why:
- Wool and down are warm — they trap tiny pockets of air
- Double-pane windows work — the air gap between the panes insulates
- Styrofoam is a good insulator — it is mostly air bubbles
- Fluffy snow insulates the ground in winter — it traps air

---

## Conductors and Insulators

Why does a metal feel cold to the touch even when it is at room temperature, while a wooden table feels warm?

The metal and the wood are the same temperature. But metal conducts heat away from your hand very quickly, so your skin cools and you sense "cold." Wood conducts heat away slowly, so your hand stays warm and you sense "warm."

```
Your hand touching metal:
                              
  HAND (37°C) ───→ METAL (20°C) ───→ rest of metal
  
  Heat flows quickly away from hand → hand feels cold
  
Your hand touching wood:
  
  HAND (37°C) ─→ WOOD (20°C)   (slow conduction)
  
  Heat stays near skin → hand feels warm
```

The sensation of temperature is actually a sensation of heat flow rate, not temperature itself. Metal and wood at the same temperature feel different because they conduct heat at different rates.

This is also why wool socks feel warm even though the wool itself is not a heat source — the wool just slows heat loss from your feet.

---

## Double-Pane Windows

A single pane of glass conducts heat quite well. Houses lose a lot of heat through windows in winter. **Double-pane windows** (also called double glazing or insulated glass units) use a sealed air gap between two panes of glass to dramatically reduce heat loss.

```
Double-Pane Window Cross-Section:

  OUTSIDE                                     INSIDE
  (cold)                                      (warm)
  
  ║ glass pane ║▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓║ glass pane ║
               ←── air gap ──→
               (sealed, still air)
               
  Heat must cross:  glass → air → glass
  
  Still air has very low k (0.025 W/m·°C)
  → Heat loss through window is greatly reduced
```

Some high-performance windows use **argon gas** in the gap instead of air. Argon has slightly lower thermal conductivity than air (k ≈ 0.017 W/m·°C) and dramatically reduces convection within the gap. Even better: **krypton gas** (k ≈ 0.0088 W/m·°C) is used in premium triple-pane windows.

Low-emissivity (low-e) coatings on the glass also reduce heat loss by radiation (the third mechanism, discussed later).

---

## Worked Example 1

**Problem:** A brick wall is 20 cm thick and has an area of 15 m². The outside temperature is 5°C and the inside is 20°C. How much heat is conducted through the wall per hour? (k for brick = 0.72 W/m·°C)

**Given:**
- k = 0.72 W/m·°C
- A = 15 m²
- ΔT = 20 - 5 = 15°C
- L = 0.20 m
- Time = 1 hour = 3600 seconds

**Rate of heat flow:**

```
Q/t = k × A × ΔT / L
Q/t = 0.72 × 15 × 15 / 0.20
Q/t = 0.72 × 15 × 75
Q/t = 810 W = 810 J/s
```

**Total heat lost in 1 hour:**

```
Q = (Q/t) × t
Q = 810 W × 3600 s
Q = 2,916,000 J = 2916 kJ ≈ 2.9 MJ
```

**Answer:** The wall loses about 810 watts continuously, or 2.9 megajoules per hour. This is why insulation (which reduces k_effective) is so important for energy efficiency.

---

## Mechanism 2: Convection

**Convection** is the transfer of heat by the actual bulk movement of a fluid (liquid or gas) carrying thermal energy with it from one place to another.

Unlike conduction (where vibrations pass energy along and the material stays put), in convection the hot material itself moves.

### How Natural Convection Works

When a fluid is heated, it expands (recall Chapter 38). Expanded fluid is less dense than the surrounding cooler fluid. Less dense fluid floats upward (buoyancy). Cooler, denser fluid sinks. This creates a **convection current** — a circulating loop of fluid that continuously carries heat from hot regions to cold regions.

```
Natural Convection Current:

    cool fluid descends        cool fluid descends
    ↓                ↓         ↓               ↓
    ↓    ←←←←←←←←←←←←←←←←←←↓               ↓
    ↓    ←── cool fluid ←←←  ↓               ↓
    ↓                         ↓               ↓
 ──────────────────────────────────────────────────
          HEATER (hot surface)
 ──────────────────────────────────────────────────
    ↑                         ↑               ↑
    ↑    hot fluid rises →    ↑               ↑
    ↑    →→→→→→→→→→→→→→→→→→  ↑               ↑
    
    Arrows show the circulation pattern.
    Hot fluid rises, spreads, cools, sinks, returns to heater.
```

The density difference due to thermal expansion drives the entire convection loop. Remove gravity (e.g., in space) and natural convection stops — the hot fluid just sits there.

---

## Natural and Forced Convection

**Natural (free) convection**: driven purely by density differences due to temperature. No external mechanism needed. Examples: the movement of air over a hot radiator, ocean currents, the atmosphere.

**Forced convection**: a pump, fan, or other external mechanism forces the fluid to move. Examples: a car's cooling system (water pump circulates coolant), a hair dryer (fan forces air over heating element), a fan oven (fan circulates hot air for even cooking), the human circulatory system (heart pumps blood to distribute body heat).

Forced convection is generally much more efficient than natural convection because the faster-moving fluid carries heat away more quickly.

```
Natural vs Forced Convection:

Natural (gentle, slow):
  
      ↑↑ warm air rising slowly ↑↑
      
      ════════════════════════
              radiator
      ════════════════════════

Forced (vigorous, fast):
  
  → → → → → → → → → → → →
  → → → → → → → → → → → →    ← fan blowing air across
  → → → → → → → → → → → →
      ════════════════════════
              heater
      ════════════════════════
```

---

## Convection in Nature

### Ocean Currents

The ocean is driven by a massive global convection system called the **thermohaline circulation** (or ocean conveyor belt). Cold, salty water near the poles is denser and sinks to the ocean floor. Warm water from the tropics flows at the surface to replace it. This global circulation carries enormous amounts of heat around the planet.

The **Gulf Stream** is a prominent example: warm surface water from the Gulf of Mexico flows northeast across the Atlantic, keeping northwestern Europe much warmer than its latitude would suggest. London (51°N) has milder winters than Winnipeg (50°N) largely because of this oceanic heat transport.

### Atmospheric Convection

The Sun heats the Earth's surface, which heats the air above it. Hot air rises (creating low-pressure zones), cools at altitude, and sinks back down (creating high-pressure zones). This drives winds. Thunderstorms are intense local convection cells. The Hadley cells, Ferrel cells, and Polar cells are large-scale atmospheric convection patterns that drive global weather.

### Inside the Earth

The Earth's mantle is solid but behaves like a very viscous fluid over millions of years. Heat from radioactive decay in the core drives slow convection currents in the mantle. These currents drag the tectonic plates, causing continental drift, earthquakes, and volcanic activity.

### In Stars

The outer layers of the Sun are in constant convection. Hot plasma rises from the radiative zone, reaches the surface, dumps energy as light and heat, cools, and sinks. The "granulation" visible on the Sun's surface (with telescopes) is the tops of these convection cells — each cell is about 1000 km across.

---

## How Radiators Heat Rooms

Despite their name, **radiators** (the old cast-iron type used in central heating) actually heat rooms primarily by **convection**, not radiation.

Hot water from the boiler flows through the radiator. The radiator heats the air touching it. This warm air becomes less dense and rises. Cooler room air flows in at the bottom of the radiator to replace it. A convection loop forms in the room.

```
How a Room Radiator Heats a Room (Mostly Convection):

    ↑  warm air rises along wall
    ↑  spreads across ceiling
    ↑  cools, sinks far wall
    ↑  returns along floor to radiator
    
    (ceiling)
    ──────────────────────────────────────
    ↑ ← ← ← ← ← ← ← ← ← ← ← ← ← ← ↓
    ↑                                   ↓
    ↑                                   ↓
    ↑                                   ↓
    ════════════                        ↓
      radiator     → → → → → → → → → ↓ ← sinks
    ════════════                        
    ──────────────────────────────────────
    (floor)
```

The actual heat output is roughly 70% convection and 30% radiation for a typical cast-iron radiator. Modern flat-panel radiators are similar; fan-assisted radiators are almost entirely convection.

---

## Worked Example 2

**Conceptual Problem:** Why does a single hot drink cool faster when placed near an open window compared to in a still room, even if both locations are the same temperature?

**Answer:**

Near the open window, air currents (forced convection) continuously replace the warm air layer immediately around the cup with fresh cool air. The warm air boundary layer is constantly swept away. This maintains a steep temperature gradient (large ΔT) between the cup surface and the surrounding air at all times.

In a still room, the layer of air immediately around the cup gradually warms up. This reduces the effective temperature difference between the cup surface and its immediate surroundings, slowing down further convective heat loss.

The key principle: **forced convection is more efficient than natural convection** because it constantly renews the temperature gradient. This is also why blowing on hot food cools it — you are replacing the warm air layer with cooler air via forced convection.

---

## Mechanism 3: Radiation

**Radiation** (also called thermal radiation or infrared radiation) is the transfer of energy by electromagnetic waves. Unlike conduction and convection, radiation requires **no medium** — it can travel through a perfect vacuum.

All objects with a temperature above absolute zero (0 K) emit thermal radiation. The hotter the object, the more radiation it emits, and the shorter the wavelength (higher frequency, higher energy) of that radiation.

```
Thermal Radiation from Objects at Different Temperatures:

  Object          | Temperature | Peak radiation wavelength | What you see/feel
  ----------------|-------------|---------------------------|------------------
  Your body       |   37°C = 310 K | ~9,400 nm (infrared) | Felt as warmth
  A warm radiator |  ~80°C = 353 K | ~8,200 nm (infrared) | Felt as warmth
  A red-hot iron  | ~800°C = 1073K | ~2,700 nm (near-IR) | Glows red
  A candle flame  | ~1700°C = 1973K| ~1,470 nm (orange) | Orange-yellow light
  Lightbulb filament|~2700°C=2973K | ~1,000 nm (white) | Bright white light
  Surface of Sun  | ~5500°C = 5778K| ~502 nm (green) | Full white spectrum
```

The peak wavelength shifts to shorter, higher-energy radiation as temperature rises. This is described by **Wien's Displacement Law**: λ_peak = b/T, where b = 2.898×10⁻³ m·K. Stars hotter than the Sun peak in blue or ultraviolet; cooler stars peak in red or infrared.

```
Radiation vs Temperature:

  Low temperature: →→→→→→→→  (long wavelength, infrared — invisible)
  
  Medium temperature: →→→→→→  (near infrared, red glow)
  
  High temperature: →→→→  (visible light — glowing hot metal)
  
  Very high temperature: →→  (blue/white — very hot stars)
```

---

## Stefan-Boltzmann Law

The total power (energy per second) radiated by an object depends on its temperature, size, and surface type:

```
P = ε × σ × A × T⁴

Where:
  P = power radiated (watts)
  ε = emissivity (dimensionless, 0 to 1)
  σ = Stefan-Boltzmann constant = 5.67×10⁻⁸ W/m²·K⁴
  A = surface area of the object (m²)
  T = absolute temperature in KELVIN (not Celsius!)
```

**Emissivity ε** (Greek letter epsilon) measures how well a surface emits radiation compared to a perfect emitter (called a **blackbody** with ε = 1).

```
Surface type               | Emissivity ε
---------------------------|-------------
Blackbody (theoretical)    |  1.00 (perfect emitter)
Human skin (all colours)   |  0.97-0.99
Flat black paint           |  0.97
Black asphalt              |  0.93
Brick                      |  0.90
Water                      |  0.96
White paint                |  0.90
Glass                      |  0.90
Polished copper            |  0.03
Polished aluminium         |  0.05
Polished silver            |  0.02
```

Notice something surprising: **the colour of a surface matters less than you might think**. White paint and black paint have similar emissivities! This seems counterintuitive. The difference is that black paint absorbs and emits visible light well, while both black and white paint are good emitters of infrared radiation (the relevant wavelength for thermal emission at room temperature).

The distinction matters mainly at visible wavelengths (solar radiation) not infrared (thermal emission). White surfaces reflect visible sunlight and absorb less solar energy (useful in hot climates), but they emit infrared as well as black surfaces.

The T⁴ dependence makes radiation extremely sensitive to temperature. Double the temperature and the radiated power increases by 2⁴ = 16 times.

---

## Dark vs Light Surfaces

Good emitters of radiation are also good **absorbers** of radiation. Poor emitters (shiny metal surfaces) are also poor absorbers — they reflect radiation instead.

This has practical consequences:

**Staying cool in summer:**
- Wear light-coloured clothing to reflect solar radiation
- Paint buildings white in hot climates (Mediterranean, Middle East)
- White roofs reflect sunlight and absorb less solar heat

**Staying warm:**
- Dark surfaces absorb more radiation from the Sun
- Solar collectors are painted matt black (high absorptivity/emissivity)
- Thermal cameras detect people and animals by their infrared emission

**Space applications:**
- Spacecraft are often wrapped in gold-coloured (polished metal) thermal blankets — very low ε means minimal heat loss by radiation in the cold of space
- The dark side of the Moon (when in shadow) cools to -173°C; the lit side reaches 127°C

---

## The Greenhouse Effect

The **greenhouse effect** is a natural phenomenon that keeps Earth warm enough for life, mediated by thermal radiation.

The Earth receives solar radiation mainly in visible and near-infrared wavelengths (because the Sun is hot — about 5778 K). The Earth absorbs this energy and warms up. As a warm body, the Earth re-emits radiation — but at its own temperature (~288 K, 15°C average), which means it emits in the **far infrared** (long wavelength).

Greenhouse gases in the atmosphere — primarily water vapour (H₂O), carbon dioxide (CO₂), methane (CH₄), and nitrous oxide (N₂O) — are transparent to visible light but absorb and re-emit far infrared radiation. They act like a thermal blanket around the Earth.

```
The Greenhouse Effect:

                     Space
                      ↑↑
                ← → re-emitted infrared
               ↗              ↗
    ╔═══════════════════════════════════╗
    ║   ATMOSPHERE                      ║  ← greenhouse gases
    ║   (transparent to visible,         ║    absorb & re-emit
    ║    opaque to far infrared)         ║    infrared
    ╚═══════════════════════════════════╝
               ↓                  ↑ far infrared
      visible light            re-radiated by Earth
         ↓
    ╔═══════════════════════════════════╗
    ║        EARTH'S SURFACE            ║
    ║   absorbs visible, heats up        ║
    ╚═══════════════════════════════════╝
```

Without any greenhouse effect, Earth's average temperature would be about -18°C. With the natural greenhouse effect, it is about +15°C — a difference of 33°C that makes liquid water (and life) possible.

The concern about climate change is the **enhanced greenhouse effect**: human activities (burning fossil fuels, deforestation) increase atmospheric CO₂ and other greenhouse gases, trapping more outgoing infrared radiation and causing additional warming beyond the natural baseline.

---

## Worked Example 3

**Problem:** Estimate the power radiated by a human body. Assume:
- Surface area A = 1.8 m²
- Skin temperature = 34°C (skin is slightly cooler than core body temperature)
- Emissivity of skin ε = 0.97
- Room temperature = 20°C

**Part 1: Power radiated by the body**

Temperature of body: T_body = 34 + 273.15 = 307.15 K ≈ 307 K

```
P_radiated = ε × σ × A × T⁴
P_radiated = 0.97 × 5.67×10⁻⁸ × 1.8 × (307)⁴
```

Calculate (307)⁴:
```
307² = 94,249
307⁴ = (94,249)² ≈ 8.883×10⁹
```

```
P_radiated = 0.97 × 5.67×10⁻⁸ × 1.8 × 8.883×10⁹
P_radiated = 0.97 × 5.67×10⁻⁸ × 1.599×10¹⁰
P_radiated = 0.97 × 906.6
P_radiated ≈ 879 W
```

**Part 2: Power absorbed from room radiation**

Temperature of room: T_room = 20 + 273.15 = 293.15 K ≈ 293 K

```
P_absorbed = ε × σ × A × T_room⁴
293⁴: 293² = 85,849 ; 293⁴ = (85,849)² ≈ 7.370×10⁹

P_absorbed = 0.97 × 5.67×10⁻⁸ × 1.8 × 7.370×10⁹
P_absorbed = 0.97 × 5.67×10⁻⁸ × 1.327×10¹⁰
P_absorbed ≈ 730 W
```

**Part 3: Net radiation loss**

```
P_net = P_radiated - P_absorbed
P_net = 879 - 730 = 149 W
```

**Answer:** The human body radiates about 149 W of net power by radiation alone at room temperature. The body also loses heat by conduction and convection, and must generate about 80-100 watts of metabolic power just to maintain body temperature. This is why a room full of people gets warm — each person is effectively a 80-150 W heater.

---

## The Thermos Flask

A **thermos flask** (vacuum flask, Dewar flask) is an elegant engineering device that minimises all three forms of heat transfer simultaneously.

```
Thermos Flask Cross-Section:

          Outer casing (plastic or metal)
          │
          │  ┌──────────────────────┐
          │  │      Vacuum          │  ← eliminates conduction
          │  │     (no air)         │     and convection
          │  │                      │
          └─→│  Inner flask         │
             │  (silvered surface)  │  ← eliminates radiation
             │                      │     (highly reflective)
             │     HOT or           │
             │     COLD liquid      │
             │                      │
             │  (silvered surface)  │
             └──────────────────────┘
             
             Outer silvered surface also reflects radiation from room
```

### How It Fights Each Heat Transfer Mechanism

**Conduction:** The inner flask is separated from the outer casing by a vacuum (or very thin glass supports). Nothing conducts heat through a vacuum. The neck of the flask is made as narrow as possible to minimise conduction through the glass walls at the top.

**Convection:** Convection requires a fluid to carry heat. In a vacuum, there is no fluid. Convection through the neck is minimised by the narrow opening.

**Radiation:** The inner and outer surfaces of the vacuum gap are silvered (highly polished, mirror-like). Polished metal has very low emissivity (ε ≈ 0.02-0.05). Both surfaces reflect most of the infrared radiation back rather than absorbing it. Hot liquid tries to radiate energy outward; the silvered inner wall reflects most of it back. Cold air outside tries to radiate inward; the silvered outer wall reflects most of it away.

The result: a well-designed thermos flask can keep coffee hot for 12+ hours or liquid nitrogen cold for days.

### Summary of Thermos Design:

```
Problem               | Solution
----------------------|----------------------------------------
Conduction outward    | Vacuum gap (no material to conduct through)
Convection outward    | Vacuum gap (no fluid to convect through)
Radiation outward     | Silvered surfaces (low emissivity, high reflectivity)
Heat loss at top      | Narrow neck, insulating stopper
```

---

## Summary

Heat is transferred by three distinct mechanisms:

**Conduction** occurs through direct contact by atom-to-atom vibration transfer (and free electrons in metals). Rate described by Fourier's Law: Q/t = kAΔT/L. Metals are excellent conductors (k up to 400 W/m·°C); wood, air, and foam are excellent insulators (k ≈ 0.02-0.17 W/m·°C). Double-pane windows trap still air to reduce conduction.

**Convection** occurs through bulk movement of a fluid. Hot fluid rises (less dense), cool fluid sinks (denser), creating circulating convection currents. Natural convection is driven by density differences; forced convection uses pumps or fans. Drives ocean currents, atmospheric circulation, and room heating by radiators.

**Radiation** occurs through electromagnetic waves (mainly infrared) and requires no medium — it works in vacuum. All objects above 0 K emit radiation. Power described by Stefan-Boltzmann Law: P = εσAT⁴. Dark, rough surfaces are better emitters and absorbers. Polished metal surfaces are poor emitters (low ε). The greenhouse effect is caused by atmospheric gases absorbing outgoing infrared radiation.

**The thermos flask** minimises all three: vacuum eliminates conduction and convection; silvered surfaces reduce radiation.

---

## Key Equations

```
Fourier's Law of Conduction:
    Q/t = k × A × ΔT / L
    
    k = thermal conductivity (W/m·°C)
    A = cross-sectional area (m²)
    ΔT = temperature difference (°C or K)
    L = thickness (m)

Stefan-Boltzmann Law (Radiation):
    P = ε × σ × A × T⁴
    
    P = radiated power (W)
    ε = emissivity (0 to 1)
    σ = Stefan-Boltzmann constant = 5.67×10⁻⁸ W/m²·K⁴
    A = surface area (m²)
    T = temperature in KELVIN

Stefan-Boltzmann constant:
    σ = 5.67×10⁻⁸ W/m²·K⁴

Net radiation between two objects:
    P_net = ε × σ × A × (T_hot⁴ - T_cold⁴)

Wien's Displacement Law (peak wavelength):
    λ_peak = b / T
    b = 2.898×10⁻³ m·K
```

---

*Next Chapter: Chapter 40 — Kinetic Theory of Gases: What Is Inside That Gas?*
