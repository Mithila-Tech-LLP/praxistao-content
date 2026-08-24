# Chapter 44: Heat Engines, Carnot Cycle, and Refrigerators

> "The motive power of heat is independent of the agents employed to realize it; its quantity is fixed solely by the temperatures of the bodies between which is effected, finally, the transfer of the caloric."
> — Sadi Carnot, 1824

---

## Table of Contents

1. What Is a Heat Engine?
2. Energy Flow in a Heat Engine
3. Efficiency of a Heat Engine
4. The Carnot Cycle: The Perfect Engine
5. Carnot Efficiency Formula
6. Why No Real Engine Can Reach Carnot Efficiency
7. Worked Example 1: Steam Turbine Efficiency
8. Worked Example 2: Carnot Efficiency Calculation
9. Real Engines: A History of Steam
10. The Otto Cycle: Your Petrol Car Engine
11. The Diesel Engine
12. Jet Engines and Gas Turbines
13. Refrigerators: Heat Engines in Reverse
14. Coefficient of Performance
15. Heat Pumps: Heating More Efficiently Than Electric Heaters
16. Worked Example 3: Refrigerator COP
17. Worked Example 4: Heat Pump Analysis
18. Summary
19. Key Equations

---

## 1. What Is a Heat Engine?

A **heat engine** is any device that converts heat energy into mechanical work. The basic idea is simple: you use a temperature difference to generate useful work.

Think about these examples:
- A steam locomotive burns coal, uses the heat to boil water, and the steam drives pistons that turn wheels
- A coal power station burns coal, makes steam, drives turbines connected to generators
- A petrol car engine burns fuel in cylinders, the expanding hot gases push pistons, turning the crankshaft
- A jet aircraft engine burns aviation fuel, the expanding hot exhaust provides thrust

All of these are heat engines. They all share the same basic structure:

    1. Absorb heat Q_H from a hot source (burning fuel, nuclear reaction, geothermal heat)
    2. Convert SOME of that heat into useful work W
    3. Reject the remaining heat Q_C to a cold "sink" (atmosphere, river, cooling tower)

The fundamental challenge — and the profound limitation from the Second Law — is that you can NEVER convert all of Q_H into work. Some Q_C must always be wasted.

### Why Do Heat Engines Need a Cold Reservoir?

Imagine trying to run a steam engine with just a boiler and no exhaust. The steam would push the piston once, but then you'd have nowhere for the steam to go. To reset the piston and do work again, you need to exhaust the steam somewhere cold — a condenser, the atmosphere, a river.

The cold reservoir is not a design flaw. It is a thermodynamic necessity. Without a temperature difference, there is no driving force for heat to flow, and no work can be extracted.

---

## 2. Energy Flow in a Heat Engine

Let us set up the standard notation:

    T_H = temperature of the hot reservoir (Kelvin)
    T_C = temperature of the cold reservoir (Kelvin)
    Q_H = heat absorbed from hot reservoir per cycle (Joules)
    Q_C = heat rejected to cold reservoir per cycle (Joules)
    W   = net work output per cycle (Joules)

The energy flow looks like this:

    HOT RESERVOIR (T_H)
          |
          | Q_H (heat flows in)
          v
    +------------------+
    |                  |
    |   HEAT ENGINE    |---> W (useful work out)
    |                  |
    +------------------+
          |
          | Q_C (waste heat flows out)
          v
    COLD RESERVOIR (T_C)

By conservation of energy (First Law of Thermodynamics), the energy going in must equal the energy going out:

    Q_H = W + Q_C

Rearranging:

    W = Q_H - Q_C

The useful work output is the difference between the heat absorbed from the hot reservoir and the heat rejected to the cold reservoir.

### What Counts as the "Hot" and "Cold" Reservoirs?

For a car engine:
- Hot reservoir: burning fuel in the cylinder (temperatures up to 2000°C)
- Cold reservoir: atmosphere (about 20°C) or coolant system

For a power station:
- Hot reservoir: steam from boiler (typically 500-600°C)
- Cold reservoir: river water or cooling towers (typically 20-40°C)

The bigger the temperature difference (T_H - T_C), the more work you can potentially extract.

---

## 3. Efficiency of a Heat Engine

**Thermal efficiency** (eta) measures what fraction of the absorbed heat is converted to useful work:

    eta = W / Q_H

This is the ratio of what you get (useful work) to what you pay for (heat from the hot source).

Since W = Q_H - Q_C:

    eta = (Q_H - Q_C) / Q_H = 1 - Q_C/Q_H

Expressed as a percentage, this ranges from 0% (no work produced, all heat wasted) to 100% (all heat converted to work — but the Second Law forbids this!).

### Real Engine Efficiencies

- Coal power plant: about 35-45%
- Combined-cycle gas power plant: up to 60%
- Petrol car engine: about 25-30%
- Diesel engine: about 40-50%
- Steam locomotive: about 6-10% (terrible!)
- Jet engine: about 35-45% thermal efficiency

These are all well below 100%. Why? Because of the Second Law and irreversible losses.

---

## 4. The Carnot Cycle: The Perfect Engine

In 1824, a young French engineer named **Sadi Carnot** (1796–1832) asked a remarkable question: what is the MAXIMUM possible efficiency for a heat engine operating between two temperature reservoirs?

He answered it by imagining a theoretical ideal engine now called the **Carnot engine**, which operates on the **Carnot cycle**.

The Carnot cycle has four steps, all perfectly reversible (no friction, no heat flow across temperature differences):

    Step 1: Isothermal Expansion at T_H
    ====================================
    The gas is in contact with the hot reservoir at T_H.
    The gas expands slowly and isothermally (at constant temperature T_H).
    It absorbs heat Q_H from the hot reservoir.
    The gas does work as it expands.
    
    Step 2: Adiabatic Expansion
    ============================
    The gas is thermally isolated (no heat exchange).
    The gas continues to expand, doing work.
    As it expands, it cools from T_H down to T_C.
    No heat flows; all work comes from the gas's internal energy.
    
    Step 3: Isothermal Compression at T_C
    =======================================
    The gas is in contact with the cold reservoir at T_C.
    The gas is compressed slowly and isothermally (at constant T_C).
    It rejects heat Q_C to the cold reservoir.
    Work is done ON the gas.
    
    Step 4: Adiabatic Compression
    ==============================
    The gas is thermally isolated again.
    The gas is compressed back to its original state.
    As it is compressed, it heats up from T_C back to T_H.
    No heat flows; work is done on the gas.

After step 4, the gas is back to its original temperature, pressure, and volume. The cycle repeats.

### Pressure-Volume Diagram of the Carnot Cycle

    Pressure
       ^
       |
    P1 |****               (1)
       |    ****
       |        *****      Isothermal expansion (T_H, absorbs Q_H)
    P2 |             *  (2)
       |              *\
       |                \  Adiabatic expansion
       |                 \
    P3 |                  * (3)
       |                 ****
       |            *****     Isothermal compression (T_C, rejects Q_C)
    P4 |      **(4)
       |      |
       |      Adiabatic compression (returns to start)
       +-------------------------------------> Volume
            V1  V2            V3  V4

The area enclosed by the cycle on this P-V diagram equals the net work W done per cycle.

### Key Properties of the Carnot Cycle

- All four steps are reversible (idealized, frictionless, quasi-static)
- Maximum possible efficiency for given T_H and T_C
- Delta S_universe = 0 (it is reversible)
- No engine operating between the same two temperatures can be more efficient

---

## 5. Carnot Efficiency Formula

For the Carnot cycle, the ratio Q_C/Q_H equals the ratio T_C/T_H:

    Q_C/Q_H = T_C/T_H    (temperatures in KELVIN!)

Therefore, the Carnot efficiency is:

    eta_Carnot = 1 - T_C/T_H    (KELVIN temperatures!)

This is one of the most important equations in thermodynamics. Notice what it says:

- The efficiency depends ONLY on the temperatures of the hot and cold reservoirs
- The efficiency is higher when T_H is large (very hot source)
- The efficiency is higher when T_C is small (very cold sink)
- You can only reach 100% efficiency if T_C = 0 K (absolute zero) — which is impossible
- Therefore, no real engine operating at achievable temperatures can ever be 100% efficient

### Crucial Warning: Use Kelvin!

The Carnot formula ONLY works with absolute temperatures in Kelvin.

    T(K) = T(°C) + 273

A common mistake is to plug in Celsius temperatures. This gives completely wrong answers.

For example: A steam engine with hot steam at 200°C and cold exhaust at 20°C.

WRONG (using Celsius): eta = 1 - 20/200 = 0.90 = 90% (incorrect!)
RIGHT (using Kelvin): eta = 1 - 293/473 = 0.38 = 38% (correct)

Always convert to Kelvin first.

---

## 6. Why No Real Engine Can Reach Carnot Efficiency

The Carnot cycle is an idealization that assumes:
- Perfectly reversible processes (no friction, no turbulence)
- Infinitely slow operation (quasi-static)
- Perfect thermal contact for isothermal steps
- Perfect thermal isolation for adiabatic steps

None of these can be achieved in practice:

**Friction**: Any moving parts (pistons, turbine blades, bearings) produce friction, converting kinetic energy to heat. Irreversible. Reduces efficiency.

**Finite temperature differences**: In real engines, heat must flow across finite temperature differences to do so at a practical rate. But heat flow across a finite temperature gap is irreversible (increases entropy). This is why the Carnot cycle requires infinitely slow operation for its isothermal steps — but infinitely slow means zero power output, useless in practice.

**Turbulence and non-equilibrium effects**: Real gas flows are turbulent and non-uniform.

**Heat leaks**: No perfect thermal insulator exists. Adiabatic steps always have some heat leakage.

**Combustion**: Burning fuel is an irreversible chemical reaction.

The Carnot efficiency is a theoretical UPPER BOUND. Real engines always fall below it. The gap between Carnot efficiency and actual efficiency quantifies how much irreversibility the engine has.

---

## 7. Worked Example 1: Steam Turbine in a Power Station

A coal power station has steam at 550°C entering the turbine and exhausts at 40°C into a condenser cooled by river water.

(a) What is the maximum (Carnot) efficiency?
(b) The actual efficiency is 42%. What fraction of the maximum efficiency is achieved?
(c) If the power station generates 1000 MW of electrical power, how much heat is absorbed from the boiler per second?

**Solution:**

Convert temperatures to Kelvin:
    T_H = 550 + 273 = 823 K
    T_C = 40 + 273 = 313 K

(a) Carnot efficiency:
    eta_Carnot = 1 - T_C/T_H = 1 - 313/823 = 1 - 0.380 = 0.620 = 62.0%

(b) Fraction of maximum efficiency:
    Fraction = eta_actual / eta_Carnot = 0.42 / 0.62 = 0.677 = 67.7%

The power station achieves about 68% of the theoretical maximum efficiency.

(c) Power output W = 1000 MW = 1 × 10^9 W

Using eta = W/Q_H:
    Q_H = W / eta = (1 × 10^9) / 0.42 = 2.38 × 10^9 W

The boiler must supply 2.38 × 10^9 J of heat per second (2380 MW).

Of this, 1000 MW becomes electricity and the remaining 1380 MW is dumped as waste heat into the river and cooling towers. This is why power stations need large amounts of cooling water.

---

## 8. Worked Example 2: Analyzing Engine Performance

A heat engine absorbs 5000 J of heat per cycle from a hot reservoir at 600 K and rejects 3000 J of heat per cycle to a cold reservoir at 300 K.

(a) What is the work output per cycle?
(b) What is the actual thermal efficiency?
(c) What is the Carnot efficiency for these temperatures?
(d) Does this engine violate the Second Law?

**Solution:**

(a) Work output per cycle:
    W = Q_H - Q_C = 5000 - 3000 = 2000 J

(b) Actual thermal efficiency:
    eta = W / Q_H = 2000 / 5000 = 0.40 = 40%

(c) Carnot efficiency:
    eta_Carnot = 1 - T_C/T_H = 1 - 300/600 = 1 - 0.50 = 0.50 = 50%

(d) The actual efficiency (40%) is less than the Carnot efficiency (50%). The engine does not violate the Second Law. In fact, any engine with efficiency less than or equal to Carnot is thermodynamically permissible.

Check entropy: Delta S_universe = Q_C/T_C - Q_H/T_H = 3000/300 - 5000/600 = 10 - 8.33 = +1.67 J/K
Since Delta S_universe > 0, this is an irreversible process. The Second Law is satisfied.

---

## 9. Real Engines: A History of Steam

### Thomas Savery (1698): The First Steam Engine

The first practical steam engine was patented by Thomas Savery. It was used to pump water from flooded mines. It had no moving mechanical parts — just steam condensing to create a vacuum that sucked water up. Efficiency: horrific (well under 1%).

### Thomas Newcomen (1712): Atmospheric Engine

Newcomen's engine improved on Savery's design with a piston. Steam pushed the piston up; then cold water was injected to condense the steam, creating a vacuum that atmospheric pressure pushed the piston back down. Still very inefficient — about 0.5%.

### James Watt (1769): The Revolution

James Watt's great innovation was the separate condenser. Instead of cooling the main cylinder (wasting the heat needed to re-heat it), Watt added a separate cold chamber connected by a valve. The steam condensed there, keeping the main cylinder always hot.

This dramatically increased efficiency to around 3-5% — a huge improvement. Watt's engine powered the Industrial Revolution.

Other Watt innovations: rotary motion (converting the up-down piston motion to rotation), centrifugal governor (automatic speed control), double-acting cylinder (steam pushes both ways).

The unit of power, the **Watt (W)**, is named in his honor.

### Compound Steam Engines and Turbines

Later engineers used compound cylinders (steam expands through multiple cylinders at decreasing pressures) to extract more work. Steam turbines (where steam spins a rotor rather than pushing a piston) eventually replaced reciprocating engines for large-scale power generation.

Modern steam turbines in power stations operate with steam at 550-600°C and pressures of 150-300 bar. They achieve about 40-45% thermal efficiency, far higher than early steam engines but still limited by the Carnot constraint.

---

## 10. The Otto Cycle: Your Petrol Car Engine

The **Otto cycle** was developed by Nikolaus Otto in 1876 and remains the basis of the modern petrol (gasoline) engine. It is a four-stroke cycle:

### The Four Strokes

    STROKE 1: INTAKE
    =================
    Piston moves DOWN.
    Intake valve opens, exhaust valve closed.
    Air-fuel mixture drawn into the cylinder.
    
    Piston at bottom dead center: cylinder full of air-fuel mixture
    
         INTAKE
           |
           v
    -------+-------
    |   AIR+FUEL  |
    |             |
    |             |
    +---piston----+
             |
             |  (piston moves down)
    
    STROKE 2: COMPRESSION
    ======================
    Piston moves UP.
    Both valves closed.
    Air-fuel mixture is compressed to about 1/8 to 1/10 of original volume
    (compression ratio r = 8-10 for petrol engines).
    Temperature rises due to compression (adiabatic compression).
    
    STROKE 3: POWER (COMBUSTION)
    =============================
    Piston at top dead center.
    Spark plug fires, igniting the compressed air-fuel mixture.
    Rapid combustion raises temperature and pressure dramatically.
    High-pressure hot gases push piston DOWN — this is where work is extracted.
    
    STROKE 4: EXHAUST
    ==================
    Piston moves UP.
    Exhaust valve opens, intake valve closed.
    Hot exhaust gases are pushed out of the cylinder.
    Cycle repeats.

### Otto Cycle on a P-V Diagram

    Pressure
       ^
       |
    P3 |      * (3) Peak combustion
       |     / \
    P2 |    /   \  (2-3: constant volume combustion)
       |   /     \
       |  /  Adiabatic \
       | / compression  \ expansion
       |/                \
    P1 |*(1)              *(4)
       |  \              /
       |   \------------/  (1→4: constant volume, exhaust/intake valve change)
       +----*-----------*-----------> Volume
           V2          V1
            Compressed  Expanded

### Efficiency of the Otto Cycle

The ideal (theoretical) efficiency of the Otto cycle is:

    eta_Otto = 1 - 1/r^(gamma-1)

where:
- r = compression ratio = V_max / V_min (typically 8-10)
- gamma = ratio of specific heats = 1.4 for air

For r = 8 and gamma = 1.4:
    eta_Otto = 1 - 1/8^0.4 = 1 - 1/2.297 = 1 - 0.435 = 0.565 = 56.5%

This is the theoretical ideal. Real petrol engines achieve about 25-30% because of friction, heat losses, incomplete combustion, valve timing losses, and other irreversibilities.

### Why Not Use Higher Compression Ratios?

Higher r → higher efficiency. So why not make r = 20 or 30?

Because at high compression ratios, the fuel-air mixture can get hot enough to spontaneously ignite before the spark plug fires — this is called **knock** or **detonation**. The premature explosion can damage the engine. Petrol engines are limited to about r = 8-12 by the fuel's octane rating (resistance to autoignition).

---

## 11. The Diesel Engine

**Rudolf Diesel** patented his engine in 1892. The key difference from the petrol engine:

In a diesel engine, only air is drawn in and compressed (not an air-fuel mixture). The compression ratio is much higher — typically r = 14-24 — because there is no fuel to spontaneously ignite.

At the top of the compression stroke, diesel fuel is directly injected into the very hot compressed air. The fuel ignites spontaneously (no spark plug needed) because the compressed air is so hot (typically 700-900°C after compression).

### Advantages of Diesel

**Higher efficiency**: because of the higher compression ratio. Diesel engines achieve about 40-50% thermal efficiency.

**Better fuel economy**: diesel fuel has higher energy density and the engine converts more of it to work.

**Torque**: diesel engines produce high torque at low RPM, good for trucks and heavy vehicles.

### Diesel Cycle Characteristics

The diesel cycle differs from the Otto cycle in step 3:
- **Otto**: combustion at constant volume (spark ignition of premixed charge)
- **Diesel**: combustion at approximately constant pressure (fuel injected gradually)

The constant-pressure combustion allows for the much higher compression ratios.

### Why Diesel Exhaust Is Problematic

Diesel combustion is diffusion combustion — the fuel droplets burn as they mix with air. Some regions can be oxygen-starved, producing soot (particulates) and complex hydrocarbons. Diesel engines also produce more NOx (nitrogen oxides) due to their higher combustion temperatures. These are serious health and environmental concerns driving the shift away from diesel vehicles.

---

## 12. Jet Engines and Gas Turbines

A **jet engine** is a continuous-flow heat engine. Instead of reciprocating pistons, it uses rotating compressors and turbines.

### How a Jet Engine Works

    AIR IN
       |
       v
    [COMPRESSOR]
    (fans compress air to high pressure, raises temperature)
       |
       v
    [COMBUSTION CHAMBER]
    (fuel injected and burned continuously, greatly increases temperature)
       |
       v
    [TURBINE]
    (hot high-pressure gases expand through turbine blades,
     driving the turbine which drives the compressor — the turbine
     extracts just enough work to run the compressor)
       |
       v
    HOT EXHAUST GASES exit at high velocity
    (the momentum of the exhaust provides THRUST by Newton's 3rd law)

The jet engine follows the **Brayton cycle**:
1. Adiabatic compression (compressor)
2. Constant-pressure heat addition (combustion)
3. Adiabatic expansion (turbine)
4. Constant-pressure heat rejection (exhaust to atmosphere)

### Jet Engine Efficiency

Compressor inlet: ~25°C = 298 K (at altitude, cooler — better efficiency!)
Combustion temperature: ~1400°C = 1673 K

Carnot efficiency for these temperatures:
    eta_Carnot = 1 - 298/1673 = 1 - 0.178 = 82.2%

Actual thermal efficiency: about 35-45% (large irreversibilities from combustion, turbulence)

Modern high-bypass turbofan engines (used on commercial aircraft) achieve about 40-45% thermal efficiency and overall propulsive efficiency of about 35-40%.

### Combined-Cycle Power Plants

The most efficient power plants use a **combined cycle**: a gas turbine (similar to a jet engine) generates power and its hot exhaust gases (still at ~500°C) are used to generate steam for a second steam turbine. This uses the waste heat from the gas turbine and dramatically improves overall efficiency.

Modern combined-cycle plants achieve 58-62% overall thermal efficiency — much better than either steam or gas turbine alone.

---

## 13. Refrigerators: Heat Engines in Reverse

A **refrigerator** is essentially a heat engine running backwards. Instead of using heat to produce work, it uses work to pump heat from a cold place to a hot place.

    HOT RESERVOIR (room at T_H ~ 25°C = 298 K)
          ^
          | Q_H (heat rejected to room)
          |
    +------------------+
    |                  |
    |  REFRIGERATOR    |<--- W (electrical work input)
    |                  |
    +------------------+
          ^
          | Q_C (heat absorbed from fridge interior)
          |
    COLD RESERVOIR (fridge interior at T_C ~ 4°C = 277 K)

Energy balance:
    Q_H = Q_C + W

The refrigerator absorbs heat Q_C from the cold interior, uses electrical work W, and dumps Q_H = Q_C + W into the room.

### The Refrigeration Cycle

Modern refrigerators use a **vapor-compression cycle** with a working fluid (refrigerant, such as R-134a):

    [LOW PRESSURE LIQUID] ---> [EVAPORATOR] ---> [LOW PRESSURE GAS]
           ^                     (inside fridge)             |
           |                     absorbs Q_C                 |
           |                                                  v
    [EXPANSION VALVE]                              [COMPRESSOR]
    (throttle: reduces pressure)                   (uses work W)
           ^                                                  |
           |                                                  v
    [HIGH PRESSURE LIQUID] <--- [CONDENSER] <--- [HIGH PRESSURE GAS]
                                (back of fridge/              
                                 radiator coils)
                                 rejects Q_H to room

The refrigerant absorbs heat from inside the fridge (evaporates, becomes gas), is compressed (work input), releases heat to the room (condenses, becomes liquid), expands through a valve, and returns to the evaporator. 

### Does a Refrigerator Heat or Cool a Room?

Here is an interesting question: if you leave the fridge door open in a hot room, will it cool the room?

No! It will actually slightly heat the room. Why?

The fridge absorbs Q_C from the room air (cool). The compressor adds work W. Then Q_H = Q_C + W is dumped back into the room. Net effect: the room gains W worth of energy (from the electrical work). The room heats up, not cools down.

---

## 14. Coefficient of Performance

For refrigerators and heat pumps, we use **Coefficient of Performance (COP)** rather than efficiency, because COP can be greater than 1.

### COP for a Refrigerator

The goal is to remove heat Q_C from the cold space using work input W:

    COP_refrigerator = Q_C / W = Q_C / (Q_H - Q_C)

For the ideal (Carnot) refrigerator:
    COP_Carnot_fridge = T_C / (T_H - T_C)

A typical household refrigerator has COP around 2-4, meaning it removes 2-4 joules of heat from inside for every joule of electrical energy.

---

## 15. Heat Pumps: Heating More Efficiently Than Electric Heaters

A **heat pump** is also a refrigerator-like device, but the goal is different:
- Refrigerator goal: keep the cold space cold (focus on Q_C)
- Heat pump goal: keep a warm space warm (focus on Q_H)

A heat pump extracts heat from outside (cold air, ground, water) and delivers it inside a building.

    OUTSIDE (cold, T_C)                  INSIDE (warm, T_H)
          |                                      ^
          | Q_C absorbed from outside            | Q_H delivered inside
          v                                      |
    +--------------------------------------------------+
    |                                                  |
    |                   HEAT PUMP                      |
    |                                                  |
    +--------------------------------------------------+
                            ^
                            | W (electrical work input)

Energy balance:
    Q_H = Q_C + W

The heat pump delivers Q_H = Q_C + W inside, where Q_C is free heat from outside and W is paid electrical energy.

### COP for a Heat Pump

    COP_heat_pump = Q_H / W = Q_H / (Q_H - Q_C)

For the ideal (Carnot) heat pump:
    COP_Carnot_heat_pump = T_H / (T_H - T_C)

Note: COP_heat_pump = COP_refrigerator + 1

**This is remarkable**: a heat pump can deliver MORE heat energy than the electrical energy it consumes! This does not violate energy conservation because it is also drawing free heat from the outside air.

For example:
- Outside temperature: 5°C = 278 K
- Inside temperature: 20°C = 293 K

Carnot COP for heat pump:
    COP_Carnot = 293 / (293 - 278) = 293 / 15 = 19.5

In practice, real heat pumps achieve COP = 3-5 in mild weather (still much more than 1).

Compare with electric resistance heating: COP = 1 exactly (100% of electrical energy becomes heat). A heat pump with COP = 3 provides 3 joules of heat per joule of electricity, making it 3× more energy-efficient.

This is why heat pumps are central to decarbonizing home heating.

### When Heat Pumps Work Best

Heat pumps are most efficient when the temperature difference is small (mild weather). As outside temperature drops, the heat pump must do more work to pump heat across a larger temperature gap, reducing COP.

Below about -15°C to -20°C, simple air-source heat pumps become less effective. Ground-source heat pumps (which draw heat from the ground, which stays at ~10°C year-round) maintain higher efficiency in cold climates.

---

## 16. Worked Example 3: Refrigerator Analysis

A refrigerator operates between an interior temperature of 5°C and a room temperature of 25°C. It is found to absorb 200 J of heat from the interior per cycle and use 80 J of electrical work per cycle.

(a) How much heat is rejected to the room per cycle?
(b) What is the actual COP?
(c) What is the Carnot COP for these temperatures?
(d) What is the Second Law efficiency (ratio of actual to Carnot COP)?

**Solution:**

Temperatures in Kelvin:
    T_C = 5 + 273 = 278 K
    T_H = 25 + 273 = 298 K

(a) Heat rejected to room:
    Q_H = Q_C + W = 200 + 80 = 280 J

(b) Actual COP:
    COP_actual = Q_C / W = 200 / 80 = 2.5

(c) Carnot COP:
    COP_Carnot = T_C / (T_H - T_C) = 278 / (298 - 278) = 278 / 20 = 13.9

(d) Second Law efficiency:
    Second Law efficiency = COP_actual / COP_Carnot = 2.5 / 13.9 = 0.18 = 18%

The refrigerator achieves only 18% of the theoretical maximum COP. There is large room for improvement — the real refrigerator is quite inefficient relative to the Carnot ideal.

Note: The Carnot COP is very high (13.9) because the temperature difference between inside (5°C) and outside (25°C) is very small (only 20°C). This means an ideal refrigerator could be very efficient here, but real irreversibilities drag the actual COP down considerably.

---

## 17. Worked Example 4: Heat Pump in Cold Weather

A ground-source heat pump draws heat from the ground at 8°C and delivers it to a building at 22°C. The heat pump needs to deliver 10,000 W of heat to the building.

(a) What is the Carnot COP of this heat pump?
(b) If the heat pump achieves 40% of Carnot COP, what is its actual COP?
(c) What electrical power does it consume?
(d) How much heat is extracted from the ground per second?
(e) Compare with electric resistance heating.

**Solution:**

Temperatures in Kelvin:
    T_C = 8 + 273 = 281 K
    T_H = 22 + 273 = 295 K

(a) Carnot COP for heat pump:
    COP_Carnot = T_H / (T_H - T_C) = 295 / (295 - 281) = 295 / 14 = 21.1

(b) Actual COP:
    COP_actual = 0.40 × 21.1 = 8.4

(c) Electrical power consumed:
    COP = Q_H / W,  so W = Q_H / COP = 10,000 / 8.4 = 1190 W

The heat pump uses only 1190 W of electricity to deliver 10,000 W of heat!

(d) Heat extracted from ground per second:
    Q_C = Q_H - W = 10,000 - 1190 = 8810 W

The heat pump extracts 8810 W from the ground and uses 1190 W of electricity, delivering 10,000 W total to the building.

(e) Electric resistance heater comparison:
    Electric heater: COP = 1, so 10,000 W of heat requires 10,000 W of electricity

    Savings: 10,000 - 1190 = 8810 W of electrical power saved
    
    The heat pump uses 8.4× less electricity than a resistance heater for the same heat output.

This dramatic efficiency advantage explains why governments are encouraging the switch from gas boilers and electric heaters to heat pumps for decarbonizing building heating.

---

## Summary

A **heat engine** converts heat into work by operating between a hot reservoir at T_H and a cold reservoir at T_C. It absorbs Q_H, produces work W, and rejects waste heat Q_C to the cold reservoir. Energy conservation gives Q_H = W + Q_C.

**Thermal efficiency** is eta = W/Q_H = 1 - Q_C/Q_H.

The **Carnot cycle** is the ideal reversible cycle giving the maximum possible efficiency for given temperatures:

    eta_Carnot = 1 - T_C/T_H   (temperatures in KELVIN)

No real engine can exceed Carnot efficiency. Real engines are limited by friction, finite temperature differences, and other irreversibilities.

**Real engines** include the Otto cycle (petrol car, about 25-30% actual efficiency), diesel engine (40-50%), steam turbines (35-45%), and combined-cycle gas turbines (up to 62%).

A **refrigerator** is a heat engine run in reverse: it uses work W to pump heat Q_C from a cold space and rejects Q_H = Q_C + W to a hot space. Its performance is measured by COP_fridge = Q_C/W.

A **heat pump** focuses on delivering heat Q_H to a warm space and achieves COP_pump = Q_H/W, which can be much greater than 1 — making it far more efficient than electric resistance heating.

The Carnot COP for a refrigerator is T_C/(T_H - T_C) and for a heat pump is T_H/(T_H - T_C).

---

## Key Equations

**First Law for a heat engine:**

    Q_H = W + Q_C    or equivalently    W = Q_H - Q_C

**Thermal efficiency:**

    eta = W / Q_H = 1 - Q_C / Q_H

**Carnot efficiency (temperatures in Kelvin!):**

    eta_Carnot = 1 - T_C / T_H

**Carnot heat ratios:**

    Q_C / Q_H = T_C / T_H

**COP of refrigerator (actual):**

    COP_fridge = Q_C / W = Q_C / (Q_H - Q_C)

**COP of refrigerator (Carnot limit):**

    COP_Carnot_fridge = T_C / (T_H - T_C)

**COP of heat pump (actual):**

    COP_pump = Q_H / W = Q_H / (Q_H - Q_C)

**COP of heat pump (Carnot limit):**

    COP_Carnot_pump = T_H / (T_H - T_C)

**Relationship between COP values:**

    COP_pump = COP_fridge + 1

**Ideal Otto cycle efficiency:**

    eta_Otto = 1 - 1 / r^(gamma - 1)

where r = compression ratio (typically 8-10 for petrol), gamma = 1.4 for air.

**Temperature conversion (ALWAYS use Kelvin in thermodynamics):**

    T(K) = T(°C) + 273.15
