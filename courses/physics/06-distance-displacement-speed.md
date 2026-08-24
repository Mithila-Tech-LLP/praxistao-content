# Chapter 06: Distance, Displacement and Speed

> **"Not all who wander are lost — but in physics, we always need to know exactly how far they went and where they ended up."**

---

## Table of Contents

- [1. Let's Start With a Story](#1-lets-start-with-a-story)
- [2. What is Distance?](#2-what-is-distance)
- [3. What is Displacement?](#3-what-is-displacement)
- [4. The Crucial Difference: Distance vs Displacement](#4-the-crucial-difference-distance-vs-displacement)
- [5. Reference Points and Coordinate Systems](#5-reference-points-and-coordinate-systems)
- [6. Position Notation: Using x and y](#6-position-notation-using-x-and-y)
- [7. What is Speed?](#7-what-is-speed)
- [8. The Speed Formula](#8-the-speed-formula)
- [9. Units of Speed](#9-units-of-speed)
- [10. Average Speed vs Instantaneous Speed](#10-average-speed-vs-instantaneous-speed)
- [11. Unit Conversion: km/h and m/s](#11-unit-conversion-kmh-and-ms)
- [12. Worked Examples](#12-worked-examples)
- [13. Scalars and Vectors: A First Look](#13-scalars-and-vectors-a-first-look)
- [14. Odometer vs GPS: Real Life Distance vs Displacement](#14-odometer-vs-gps-real-life-distance-vs-displacement)
- [15. Relative Position: Where is Alice Compared to Bob?](#15-relative-position-where-is-alice-compared-to-bob)
- [16. Common Mistakes and How to Avoid Them](#16-common-mistakes-and-how-to-avoid-them)
- [17. Summary](#17-summary)
- [18. Key Equations](#18-key-equations)

---

## 1. Let's Start With a Story

Imagine it is a Saturday morning. You decide to go for a walk from your home to the nearby park, but on the way you also stop at the corner store to buy a snack.

Here is a rough map of your journey:

```
HOME -----> STORE -----> PARK
  |<-- 300m -->|<-- 200m -->|

Total path walked = 300 + 200 = 500 metres
```

You walked **500 metres** in total. But here is the interesting part:

If you drew a straight line from HOME directly to PARK, that line would be only **360 metres** (we will calculate this later). You walked 500 metres, but you are only 360 metres away from where you started.

This gap — the difference between **how much you walked** and **how far you are from your starting point** — is the heart of this entire chapter.

Physics has two different words for these two different ideas:
- The total path you walked = **Distance**
- How far you are from your starting point = **Displacement**

These two concepts might seem similar, but they are completely different in physics, and mixing them up is one of the most common mistakes beginners make. By the end of this chapter, you will understand them perfectly.

---

## 2. What is Distance?

**Distance** is the total length of the path you travel, no matter which direction you go.

Think of distance as measuring every single step you take, from start to finish. You follow the road, the path, the corridor — every twist, every turn, every backtrack — and you add it all up. That total is your distance.

### Everyday Examples of Distance

- You walk from your bedroom to the kitchen through the hallway: the distance is the length of the hallway.
- A race car drives around an oval track 10 times: the distance is 10 × the length of the track.
- You go from your house to school, then later come back home: if each trip is 2 km, your total distance for the day is 4 km.

### Key Properties of Distance

- Distance is **always positive** (or zero). You can never walk a "negative" distance.
- Distance **keeps adding up**. Every step you take increases the distance.
- Distance does **not care about direction**. Going left or going right — both count the same toward your total distance.
- Distance is measured in units of **length**: metres (m), kilometres (km), centimetres (cm), etc.

### Distance is Like a Taxi Meter

Imagine a taxi meter. Every metre the taxi travels, the meter ticks up. It does not care if you went left or right, north or south, up the hill or down the hill. It just keeps counting the total length of road covered.

That is exactly what distance is — the taxi meter reading of your journey.

```
Start your journey:
Taxi meter = 0.0 km

Drive 3 km north:
Taxi meter = 3.0 km

Turn around, drive 3 km south (back to start):
Taxi meter = 6.0 km

Distance = 6.0 km  (even though you're back where you started!)
```

---

## 3. What is Displacement?

**Displacement** is how far you are from your starting point, measured in a straight line, and it includes the direction.

Think of displacement as drawing an arrow — with a ruler — directly from where you started to where you ended up right now. The length of that arrow is the magnitude (size) of your displacement, and the direction the arrow points tells you the direction of your displacement.

### The "As the Crow Flies" Idea

You may have heard the expression "as the crow flies." It means the straight-line distance between two places, ignoring roads and paths. That is exactly what displacement is — the crow's-eye-view, straight-line measurement.

```
Your actual path (winding road):
START ~~~ up ~~~ left ~~~ down ~~~ right ~~~ END
        total distance = 8 km

Displacement (straight line, crow flies):
START ---------------------------------> END
        displacement = 5 km East
```

### Displacement Can Be Zero

Here is something surprising: your displacement can be zero even after a very long journey!

If you walk around an entire city block and come back to your starting point:

```
        200m North
         ___________
        |           |
  200m  |           |  200m
  West  |           |  East
        |___________|
        200m South

Distance = 200 + 200 + 200 + 200 = 800 metres
Displacement = 0 metres  (you're back where you started!)
```

You walked 800 metres — that is real effort! But your displacement is zero because you ended up exactly where you began.

### Displacement Has Direction

This is very important. When you describe displacement, you must say both:
1. **How far** (a number with units, like 5 km)
2. **Which direction** (North, East, up, to the right, at 30 degrees, etc.)

Just saying "my displacement is 5 km" is incomplete in physics. You must say "5 km East" or "5 km in the direction of the park."

We will talk more about this when we discuss vectors later in the chapter.

---

## 4. The Crucial Difference: Distance vs Displacement

Let us put both concepts side by side with a clear example.

### The Park Example (Revisited)

```
          N
          |
          |
   W------+------E
          |
          |
          S

HOME is at the origin (0, 0)
STORE is 300m East of HOME
PARK is 200m North of STORE

Map view:
         PARK
          *
          |
          | 200m
          |
HOME *----* STORE
     300m
```

**Path 1: HOME → STORE → PARK**

- You walk 300m East to the store.
- Then 200m North to the park.
- **Distance** = 300m + 200m = **500 metres**

**Displacement: HOME to PARK (straight line)**

The straight line from HOME to PARK forms the hypotenuse of a right-angled triangle:
- One side = 300m (East)
- Other side = 200m (North)

Using the Pythagorean theorem (a² + b² = c²):
Displacement = √(300² + 200²) = √(90000 + 40000) = √130000 ≈ **360.6 metres, in a Northeast direction**

Notice:
- Distance = 500 metres
- Displacement = 360.6 metres

The distance is bigger because you took a non-straight path. The displacement is smaller because it is the shortest possible route.

### Comparison Table

| Property | Distance | Displacement |
|---|---|---|
| What it measures | Total path length | Straight-line gap from start to end |
| Considers direction? | No | Yes |
| Can it be zero? | Only if you never move | Yes — if you return to start |
| Can it be negative? | No | Yes (opposite direction) |
| Type (physics term) | Scalar | Vector |
| Symbol | d | s or Δx |
| Example | 500 m | 360 m Northeast |

### More Everyday Examples

**Example A: Running Laps**

A school athlete runs 4 laps around a 400m track.

```
         _______
        /       \
       |         |
       |  400m   |  <- one full lap
       |  track  |
        \_______/

Distance = 4 × 400m = 1600m = 1.6 km
Displacement = 0m (back at start after 4 complete laps)
```

**Example B: Elevator**

You are on the 5th floor of a building. Each floor is 3 metres tall.
- You go UP to the 8th floor.
- Then you come DOWN to the 2nd floor.

```
Floor 8 ---- highest point
  |
  | went up 3 floors = 3 × 3m = 9m
  |
Floor 5 ---- starting point
  |
  | went down 6 floors = 6 × 3m = 18m
  |
Floor 2 ---- ending point
```

- Distance = 9m + 18m = **27 metres**
- Displacement = Floor 2 - Floor 5 = 3 floors downward = **9 metres downward** (or -9m if we call up positive)

---

## 5. Reference Points and Coordinate Systems

### The Problem of "Where Are You?"

Imagine someone calls you on the phone and asks: "Where are you right now?"

You might say: "I'm at the coffee shop."

But that only makes sense if they know where the coffee shop is! If they are from another city, they would have no idea.

In physics, we have this exact problem. To describe anyone's **position** (where they are), we need a reference — a fixed point that everyone agrees on as the "starting point" or "zero point."

This fixed reference is called the **origin**.

### What is an Origin?

The **origin** is a fixed reference point that we choose as our "zero" location. All positions are measured relative to (compared to) the origin.

You can choose any point as your origin — physics doesn't care which one you pick. But once you pick it, you must be consistent.

Common choices for origin:
- The starting point of a journey
- The entrance of a building
- The center of a city
- The position of a certain object at a certain moment

### The Number Line (1D Coordinate System)

For objects moving along a straight line (like a car on a road, or a ball thrown straight up), we only need one number to describe position.

We set up a **number line**:

```
<----|----|----|----|----|----|----|----|---->
    -4   -3   -2   -1    0    1    2    3    4
                         ^
                       ORIGIN
                       (zero point)

Positive direction -->
<-- Negative direction
```

- Positions to the **right** of the origin are **positive** (e.g., +2 metres, or just 2 metres)
- Positions to the **left** of the origin are **negative** (e.g., -3 metres)

### Example: Car on a Road

A car starts at position 0 (the origin) and moves to the right.

```
Start:        Car at position x = 0
              |
<------0------*------2------4------6------>
                     ^
                   Car moves here
                   Position x = +2m
```

Later, the car reverses past the starting point:

```
<------0------*------2------4------6------>
       ^
     Car is now here
     Position x = -2m
     (2 metres to the LEFT of origin)
```

### The 2D Coordinate System

When objects move across a flat surface (like a ball rolling on the floor, or a person walking on a field), we need **two** numbers to describe position.

We use an **x-axis** (horizontal) and a **y-axis** (vertical on paper, but think of it as depth/forward on the ground).

```
          y-axis
          |
        4 |        * (2, 4)
          |
        3 |
          |
        2 |    * (1, 2)
          |
        1 |
          |
----------+-----------> x-axis
          0    1    2    3    4

Origin is at (0, 0)
```

A position like **(2, 4)** means:
- 2 units along the x-axis (to the right)
- 4 units along the y-axis (upward)

This is called a **coordinate** — a pair of numbers (x, y) that pinpoints exactly where something is.

---

## 6. Position Notation: Using x and y

### Describing Where You Are

In physics, we write position using coordinates. For 1D motion (along a line), we just write x. For 2D motion (on a flat surface), we write (x, y).

### Notation for Change in Position

The symbol **Δ** (the Greek letter Delta) means "change in." It is used everywhere in physics to mean "final value minus initial value."

So:
- **Δx** means "change in x" = final x position minus initial x position
- **Δy** means "change in y" = final y position minus initial y position

This is also exactly what **displacement** is in 1D:

```
Displacement = Δx = x_final - x_initial
```

### Example: Finding Displacement from Positions

A person starts at position x = 2m and walks to position x = 7m.

```
Start: x = 2m
End:   x = 7m

<-----|-----|-----|-----|-----|-----|-----|---->
      0     1     2     3     4     5     6     7
                  ^                             ^
               x = 2m                        x = 7m
               (start)                        (end)

Displacement = Δx = 7 - 2 = +5m (5 metres to the right)
Distance = 5m (walked straight from 2 to 7, no backtracking)
```

### Example: When You Backtrack

Same person starts at x = 2m, walks to x = 7m, then walks back to x = 4m.

```
Step 1: x = 2m --> x = 7m  (moved +5m to the right)
Step 2: x = 7m --> x = 4m  (moved -3m to the left, i.e., backtracked 3m)

<-----|-----|-----|-----|-----|-----|-----|---->
      0     1     2     3     4     5     6     7
                  ^           ^                 ^
               x = 2m      x = 4m           x = 7m
               (start)     (end)            (turned here)

Distance = 5m + 3m = 8m  (total path length)
Displacement = x_final - x_initial = 4 - 2 = +2m (net move to the right)
```

Distance and displacement are different because the person backtracked.

---

## 7. What is Speed?

### The Basic Idea

**Speed** tells you how quickly something is covering distance.

If two cars both travel 100 km, but one takes 1 hour and the other takes 2 hours, they clearly moved differently. The first car was faster — it covered the same distance in less time.

Speed captures this idea: how much distance is covered per unit of time.

### Everyday Feel for Speed

Let us build intuition first:

- A person walking leisurely: about 1 to 1.5 metres per second (m/s)
- A person jogging: about 2 to 3 m/s
- A bicycle: about 5 to 7 m/s
- A car on a city road: about 14 m/s (roughly 50 km/h)
- A car on a highway: about 28 m/s (roughly 100 km/h)
- A commercial airplane: about 250 m/s (roughly 900 km/h)
- Sound in air: about 343 m/s
- Light: 300,000,000 m/s (the fastest possible speed in the universe!)

```
Speed comparison (on a scale):

0 m/s                                          343 m/s (sound)
|---walk---|---jog---|---bike---|---car----|------>
0     1    2    3    4    5    6         28   343
```

### Speed is a Rate

In everyday language, a "rate" is how much of something happens per unit of something else.

Speed is the rate of distance covered per unit of time:

```
Speed = Distance covered per unit of time

If you cover 10 metres every second, your speed = 10 m/s
If you cover 60 km every hour, your speed = 60 km/h
```

---

## 8. The Speed Formula

The formula for speed is one of the simplest and most important in all of physics:

```
Speed = Distance / Time

Or in shorthand:

v = d / t

Where:
  v = speed (v stands for velocity — we'll use it for speed for now)
  d = distance
  t = time taken
```

### Rearranging the Formula

From `v = d / t`, you can get two other useful formulas:

```
To find Distance:   d = v × t
To find Time:       t = d / v
```

A helpful way to remember these: use the **formula triangle**.

```
     +-------+
     |   d   |
     +---+---+
     | v | t |
     +---+---+

Cover what you want to find:

Want d? Cover d: remaining is v × t, so d = v × t
Want v? Cover v: remaining is d / t, so v = d / t
Want t? Cover t: remaining is d / v, so t = d / v
```

### Units Matter!

The units of your answer depend on the units you put in.

```
If distance is in metres (m) and time is in seconds (s):
  Speed = metres / seconds = m/s

If distance is in kilometres (km) and time is in hours (h):
  Speed = kilometres / hours = km/h
```

You must be consistent. Do not mix km with seconds, or metres with hours, unless you convert first.

---

## 9. Units of Speed

### The Standard Unit: m/s

In science, the standard unit of speed is **metres per second**, written as **m/s** or **ms⁻¹**.

This means: "how many metres does the object travel in one second?"

- Speed = 5 m/s means the object travels 5 metres every second.
- Speed = 0.5 m/s means the object travels 0.5 metres (half a metre) every second.

### The Everyday Unit: km/h

For everyday use — cars, trains, planes — we use **kilometres per hour**, written as **km/h** or **kmph**.

- A car at 60 km/h travels 60 kilometres in one hour.
- A train at 200 km/h travels 200 kilometres in one hour.

### Other Units of Speed

| Unit | Written as | Used for |
|---|---|---|
| Metres per second | m/s | Science, physics problems |
| Kilometres per hour | km/h | Cars, weather, daily life |
| Miles per hour | mph | Countries using imperial system (USA, UK) |
| Knots | kt | Ships and aircraft navigation |
| Mach | Mach 1, Mach 2... | Very fast aircraft (1 Mach ≈ 340 m/s) |

---

## 10. Average Speed vs Instantaneous Speed

This is a distinction that trips up many beginners, so let us be very careful here.

### Average Speed: The Big Picture

**Average speed** is the total distance traveled divided by the total time taken for an entire journey.

```
Average Speed = Total Distance / Total Time

v_avg = d_total / t_total
```

It does not matter whether you went fast or slow at any particular moment. You just look at the whole journey from start to finish.

### Example: The Morning Commute

You drive from home to work. The journey is 30 km. You leave at 8:00 AM and arrive at 9:00 AM (1 hour).

```
Average speed = 30 km / 1 hour = 30 km/h
```

But during that hour, you definitely were not going exactly 30 km/h the entire time. Sometimes you were stuck at a red light (0 km/h), sometimes zooming on the highway (80 km/h). The average speed hides all those details.

### Instantaneous Speed: The Snapshot

**Instantaneous speed** is your speed at one particular moment in time — an instant.

This is what your car's speedometer shows you. It does not care about where you came from or where you are going. It just tells you: "right now, at this exact moment, you are traveling at X km/h."

```
Speedometer showing different readings at different moments:

8:00 AM departure     8:30 AM (highway)    8:55 AM (city traffic)
  v = 0 km/h            v = 80 km/h            v = 15 km/h
     |                       |                       |
     |                       |                       |
  [start]               [highway]               [near work]

Average speed for whole trip = 30 km/h
```

### Why the Distinction Matters

Consider this scenario:

A student cycles to school. The school is 6 km away. The student takes 30 minutes.

Average speed = 6 km / 0.5 hours = **12 km/h**

But the student's actual journey looked like this:

```
Time 0:00 to 0:10 min:  Flat road, going at 16 km/h
Time 0:10 to 0:15 min:  Uphill, slowing to 6 km/h  
Time 0:15 to 0:20 min:  Stopped at traffic light, 0 km/h
Time 0:20 to 0:30 min:  Downhill, fast at 20 km/h

Speed chart:
        20 |    *                         *
           |   * *                       * *
km/h    16 |  *   *     *               *   *
           | *     *   * *             *     *
         6 |*       * *   *           *       
           |         *     *         *        
         0 |          *     *-------*          
           +----+----+----+----+----+----
                5   10   15   20   25   30
                        Time (minutes)
```

The average speed is 12 km/h, but the student went between 0 and 20 km/h throughout the journey. The average smooths all of that out.

### Summary: Average vs Instantaneous

| | Average Speed | Instantaneous Speed |
|---|---|---|
| What it tells you | Overall rate for a whole journey | Speed at one specific moment |
| How to calculate | Total distance ÷ Total time | Reading of speedometer at that instant |
| Changes over journey? | Just one number for the whole trip | Changes constantly |
| Example | "I averaged 60 km/h on the highway" | "I was doing 80 km/h when I passed that sign" |

---

## 11. Unit Conversion: km/h and m/s

Converting between km/h and m/s is a very common task in physics. Let us learn it once, clearly, and never get confused again.

### Understanding What the Units Mean

```
1 km/h means: 1 kilometre per 1 hour

Let's break it down:
1 kilometre = 1000 metres
1 hour = 3600 seconds (60 minutes × 60 seconds)

So:
1 km/h = 1000 metres / 3600 seconds
       = 1000/3600 m/s
       = 1/3.6 m/s
       ≈ 0.278 m/s
```

### The Conversion Factor

```
To convert km/h to m/s:
  Divide by 3.6

To convert m/s to km/h:
  Multiply by 3.6
```

You can also think of it as:

```
km/h × (1000 m / 1 km) × (1 h / 3600 s) = m/s
km/h × 1000/3600 = m/s
km/h / 3.6 = m/s
```

### Worked Conversion Examples

**Example 1: Convert 72 km/h to m/s**

```
72 km/h ÷ 3.6 = 20 m/s

Check: 72 km/h means 72,000 metres in 3600 seconds
       72,000 / 3600 = 20 m/s  ✓
```

**Example 2: Convert 15 m/s to km/h**

```
15 m/s × 3.6 = 54 km/h

Check: 15 m/s means 15 metres every second
       In one hour (3600 seconds): 15 × 3600 = 54,000 metres = 54 km  ✓
```

**Example 3: Convert 100 km/h to m/s**

```
100 km/h ÷ 3.6 = 27.78 m/s ≈ 27.8 m/s

This is the legal highway speed in many countries.
It means about 28 metres every single second — about the length of a bus!
```

**Example 4: Convert 340 m/s to km/h (speed of sound)**

```
340 m/s × 3.6 = 1224 km/h

The speed of sound in air is about 1224 km/h!
```

### Handy Reference Table

| km/h | m/s | Real-world example |
|---|---|---|
| 3.6 | 1.0 | Slow walking pace |
| 18 | 5.0 | Brisk jogging |
| 36 | 10.0 | Fast cycling |
| 54 | 15.0 | Slow city traffic |
| 72 | 20.0 | Normal city driving |
| 90 | 25.0 | Fast city road |
| 108 | 30.0 | Highway driving |
| 180 | 50.0 | Fast train |
| 360 | 100.0 | High-speed train |
| 1224 | 340.0 | Speed of sound |

---

## 12. Worked Examples

Let us now work through a set of problems step by step. Physics problems become much easier if you follow the same structured approach every time.

### Problem-Solving Steps

```
Step 1: Read the problem carefully
Step 2: Write down what you KNOW (given values with units)
Step 3: Write down what you WANT (what to find)
Step 4: Choose the right formula
Step 5: Substitute numbers into the formula
Step 6: Calculate the answer with units
Step 7: Check if the answer makes sense
```

---

### Worked Example 1: Finding Speed

**Problem:** Rohan runs a 100-metre race in 12.5 seconds. What is his average speed?

```
Given:
  Distance (d) = 100 m
  Time (t)     = 12.5 s

Find:
  Speed (v) = ?

Formula:
  v = d / t

Calculation:
  v = 100 m / 12.5 s
  v = 8 m/s

Answer: Rohan's average speed is 8 m/s
```

Reality check: 8 m/s = 28.8 km/h. That is a fast sprint for a school student — sounds reasonable!

---

### Worked Example 2: Finding Distance

**Problem:** A train travels at a constant speed of 120 km/h for 2.5 hours. How far does it travel?

```
Given:
  Speed (v) = 120 km/h
  Time (t)  = 2.5 h

Find:
  Distance (d) = ?

Formula:
  d = v × t

Calculation:
  d = 120 km/h × 2.5 h
  d = 300 km

Answer: The train travels 300 km
```

Reality check: A fast train in 2.5 hours covering 300 km sounds perfectly right.

---

### Worked Example 3: Finding Time

**Problem:** A cyclist needs to travel 45 km. She can maintain an average speed of 15 km/h. How long will the journey take?

```
Given:
  Distance (d) = 45 km
  Speed (v)    = 15 km/h

Find:
  Time (t) = ?

Formula:
  t = d / v

Calculation:
  t = 45 km / 15 km/h
  t = 3 hours

Answer: The journey will take 3 hours
```

---

### Worked Example 4: Mixing Units (Careful!)

**Problem:** A car travels 200 m in 10 seconds. What is its speed in km/h?

```
Method 1: Calculate in m/s, then convert

  Given:
    Distance = 200 m
    Time     = 10 s

  v = 200 m / 10 s = 20 m/s

  Convert to km/h:
  v = 20 × 3.6 = 72 km/h

Method 2: Convert units first, then calculate

  200 m = 0.2 km
  10 s = 10/3600 hours = 0.00278 h

  v = 0.2 km / 0.00278 h ≈ 72 km/h  ✓

Answer: The car's speed is 72 km/h (or 20 m/s)
```

Both methods give the same answer. Method 1 is usually easier.

---

### Worked Example 5: Average Speed Over Two Legs

**Problem:** Priya drives from City A to City B, a distance of 120 km, at 60 km/h. She then drives from City B to City C, a distance of 80 km, at 40 km/h. What is her average speed for the entire journey?

This is a very common trick question! Many students add the speeds and divide by 2. That is WRONG. Here is the correct approach:

```
Given:
  Leg 1: d1 = 120 km, v1 = 60 km/h
  Leg 2: d2 = 80 km,  v2 = 40 km/h

Step 1: Find total distance
  d_total = d1 + d2 = 120 + 80 = 200 km

Step 2: Find time for each leg
  t1 = d1 / v1 = 120 / 60 = 2 hours
  t2 = d2 / v2 = 80 / 40  = 2 hours

Step 3: Find total time
  t_total = t1 + t2 = 2 + 2 = 4 hours

Step 4: Calculate average speed
  v_avg = d_total / t_total = 200 km / 4 h = 50 km/h

Answer: Priya's average speed is 50 km/h
```

Wrong answer if you just averaged: (60 + 40) / 2 = 50 km/h — wait, it happened to be the same here because both legs took the same time! But in general, you must always use total distance / total time.

---

### Worked Example 6: Displacement Calculation

**Problem:** A person walks 4 km East, then 3 km North. What is their displacement from the starting point?

```
     N
     |
  3km|    * END
     |    |
     |    |
     +----*-------> E
  START  4km

This forms a right triangle.
The displacement is the hypotenuse.

Using Pythagorean theorem:
Displacement² = (4 km)² + (3 km)²
Displacement² = 16 + 9
Displacement² = 25
Displacement  = √25 = 5 km

The direction is Northeast (specifically, the angle from East
toward North = arctan(3/4) = about 37 degrees North of East)

Answer: Displacement = 5 km, in a direction 37° North of East
Distance = 4 + 3 = 7 km (the path taken)
```

---

### Worked Example 7: Speed with Displacement vs Distance

**Problem:** An ant starts at point A, crawls 8 cm North to point B, then 6 cm East to point C. The ant takes 7 seconds. Find: (a) average speed, (b) displacement, (c) magnitude of average velocity (displacement/time).

```
         C
    *----*
    |    6cm
    |
   8cm (North)
    |
    *
    A (start)

(a) Distance = 8 + 6 = 14 cm
    Average speed = 14 cm / 7 s = 2 cm/s

(b) Displacement (straight line from A to C):
    Displacement = √(8² + 6²) = √(64 + 36) = √100 = 10 cm
    Direction = Northeast

(c) Average velocity magnitude = Displacement / Time
    = 10 cm / 7 s
    ≈ 1.43 cm/s

Notice: Average speed (2 cm/s) ≠ average velocity magnitude (1.43 cm/s)
This is because the path was not straight!
```

---

## 13. Scalars and Vectors: A First Look

This section is a preview — we will go much deeper into vectors in later chapters. But it is important to introduce the idea here because it explains why distance and displacement behave so differently.

### What is a Scalar?

A **scalar** is a physical quantity that is fully described by just a number (and its unit).

It has **magnitude** (size) but **no direction**.

Examples of scalars:
- Temperature: "It is 28°C" — no direction needed
- Mass: "The bag weighs 5 kg" — no direction needed
- Time: "It took 3 hours" — no direction needed
- Distance: "I walked 500 m" — just a total length, no direction needed
- Speed: "The car was doing 60 km/h" — just how fast, no direction needed

```
Scalar quantities: just a number + unit

Distance = 500 m          (not "500 m northward")
Speed = 60 km/h           (not "60 km/h eastward")
Temperature = 37°C        (not "37°C toward the window")
```

### What is a Vector?

A **vector** is a physical quantity that needs both a **magnitude** (size) AND a **direction** to be fully described.

Examples of vectors:
- Displacement: "5 km to the North" — you need the direction
- Velocity: "30 m/s East" — you need the direction
- Force: "Push of 10 N toward the right" — direction matters
- Acceleration: "2 m/s² downward" — direction matters

```
Vector quantities: number + unit + direction

Displacement = 360 m, Northeast
Velocity = 20 m/s, due East
Force = 50 N, upward
```

### Why Does Direction Matter for Some Quantities?

Imagine you push a box. If you push it to the right with 50 N of force, it moves right. If you push it to the left with 50 N, it moves left. The size of the force is the same, but the direction is different — and it makes a HUGE difference to what happens!

That is why force is a vector — direction changes the outcome.

Now imagine you measure how much time has passed. Does it matter which "direction" time is flowing? No! Time just increases. That is why time is a scalar.

### Distance vs Displacement: The Scalar/Vector Connection

```
DISTANCE:
- Just tells you total path length
- No direction
- 500 m is just 500 m
- SCALAR ✓

DISPLACEMENT:
- Tells you how far from start to end AND in which direction
- Direction is essential — 5 km North is completely different from 5 km South
- VECTOR ✓

SPEED:
- Just tells you how fast you are going
- No direction
- 60 km/h is just 60 km/h
- SCALAR ✓

VELOCITY (the vector version of speed):
- Tells you how fast AND in which direction
- 60 km/h North is different from 60 km/h South
- VECTOR ✓
```

### Visualizing Vectors with Arrows

In physics, we draw vectors as arrows:
- The **length** of the arrow = the magnitude (how big)
- The **direction** the arrow points = the direction

```
Displacement of 3m East:   ----->
Displacement of 5m East:   --------->
Displacement of 3m West:   <-----
Displacement of 3m North:  ^
                            |
                            |
Displacement of 5m at 45°:   /
                             (diagonal arrow)
```

Notice that displacement of 3m East and 3m West look like opposites — and in physics, they are! One is +3m and the other is -3m (if we define East as positive).

### The Big Picture So Far

| Quantity | Type | Needs Direction? | Formula Link |
|---|---|---|---|
| Distance | Scalar | No | v = d/t |
| Speed | Scalar | No | v = d/t |
| Displacement | Vector | Yes | related to Δx |
| Velocity | Vector | Yes | Δx/t |

We will use velocity extensively in the next chapter. For now, just know it is speed with direction added.

---

## 14. Odometer vs GPS: Real Life Distance vs Displacement

This is a fantastic real-world example that you can observe the next time you are in a car.

### The Odometer

The **odometer** is the instrument in a car that counts total kilometres traveled. Every metre of road the car rolls over, the odometer counts it.

```
Car dashboard:

   ODOMETER        SPEEDOMETER
  +----------+    +----------+
  | 12345 km |    |  60 km/h |
  +----------+    +----------+
  Total distance  Current speed
  since the car   right now
  was new
```

The odometer measures **distance** — total path length.

It does not care if you went straight, around a circle, or backwards. Every metre of road movement adds to the odometer count.

When you go for a drive and come back home, your odometer will show more kilometres than when you left — even though you are back where you started (displacement = 0).

### GPS Navigation

A **GPS** (Global Positioning System) device tracks your position at every moment using satellites. It knows exactly where you are in terms of coordinates (latitude and longitude).

When a GPS says "your destination is 5.3 km away" — it is giving you **displacement** (straight-line distance to your destination).

When it says "you have traveled 12.7 km on this route" — that is **distance** (total path followed).

```
Example: Driving through a city center

Odometer shows:      14.2 km  (total distance driven on winding city roads)
GPS "as-crow-flies": 8.7 km   (straight-line displacement from start to end)

Start *                             * End
       \                           /
        \  roads through city     /
         ~~~~~~~~~~~~~~~~~~~~    /
          ~~~~~~~~~~~~~~~~~~~   /
           ~~~~~~~~~~~~~~~~~~  /
            ~~~~~~~~~~~~~~~~~ /
             Actual road path /
             (like spaghetti!)
```

The odometer count is always greater than or equal to the GPS straight-line displacement.
- Equal: only when you travel in a perfectly straight line.
- Greater: whenever your path curves or bends.

### Why GPS Uses Displacement for "Distance to Destination"

When your GPS says "in 2.4 km, turn right" — it is not telling you the straight-line displacement to the turn. It is giving you the distance along the road to the turn. 

But when it says "your destination is 10 km away" on the initial overview, it often means the direct-line displacement to give you a rough idea of how far the destination is. Then as you drive, it tracks the actual road distance.

This shows that even in everyday technology, both concepts (distance and displacement) are used — just in different contexts.

---

## 15. Relative Position: Where is Alice Compared to Bob?

### The Concept of Relative Position

Until now, we have described position relative to a fixed origin. But in real life, we often want to know where one person or object is **relative to another person or object**.

This is called **relative position**.

**Relative position of A with respect to B = Position of A − Position of B**

### Alice and Bob: A Simple Example

Let us set up a number line representing a road. The town center is at position 0.

```
<------|----|----|----|----|----|----|----|---->
      -4   -3   -2   -1    0    1    2    3    4
                            ^
                       Town Center
                        (origin)

Alice is at position: +3 m (3 metres east of town center)
Bob is at position:   +1 m (1 metre east of town center)
```

**Where is Alice relative to Bob?**

```
Alice's position relative to Bob = Alice's position − Bob's position
                                 = +3 m − (+1 m)
                                 = +2 m

This means: Alice is 2 metres to the East (positive direction) of Bob.
```

**Where is Bob relative to Alice?**

```
Bob's position relative to Alice = Bob's position − Alice's position
                                 = +1 m − (+3 m)
                                 = −2 m

This means: Bob is 2 metres to the West (negative direction) of Alice.
```

Notice: the two relative positions are opposite in sign. That makes sense — if Alice is to the right of Bob, then Bob must be to the left of Alice.

### Alice and Bob: Both Walking

Now let us say they start walking. Alice starts at x = 0 and walks East. Bob starts at x = 10m and walks West (toward Alice).

```
Initial positions:
Alice (A) at x = 0      Bob (B) at x = 10m

<----|----|----|----|----|----|----|----|----|---->
     0         A                   B        10

After 2 seconds (Alice moves at 2 m/s East, Bob moves at 3 m/s West):
  Alice: x = 0 + 2×2 = 4m
  Bob:   x = 10 - 3×2 = 4m

<----|----|----|----|----|----|----|----|----|---->
     0              AB                      10

They have met! Both at x = 4m.
```

**Alice's position relative to Bob over time:**

```
At t = 0:  Alice relative to Bob = 0 - 10 = -10m (Alice is 10m West of Bob)
At t = 1:  Alice = 2m, Bob = 7m. Alice relative to Bob = 2 - 7 = -5m
At t = 2:  Alice = 4m, Bob = 4m. Alice relative to Bob = 4 - 4 = 0m (they meet!)
At t = 3:  Alice = 6m, Bob = 1m. Alice relative to Bob = 6 - 1 = +5m (Alice has passed Bob, is now East)
```

### Relative Position: A More Visual Example

Imagine two students, Meera and Arjun, standing in a gym. The coach (origin) stands at the center.

```
Coach at center = origin (0)
Meera is at +5m (5 metres to coach's right)
Arjun is at -3m (3 metres to coach's left)

<----|----|----|----|----|----|----|----|---->
    -4   -3   -2   -1    0    1    2    3    4    5
           ^                  ^               ^
         Arjun              Coach           Meera
         x=-3m              x=0m            x=+5m

Meera's position relative to Arjun:
= Meera's position - Arjun's position
= +5 - (-3)
= +5 + 3
= +8m

Meera is 8 metres to the RIGHT of Arjun.

Arjun's position relative to Meera:
= -3 - (+5) = -8m (Arjun is 8 metres to the LEFT of Meera)
```

This concept of relative position becomes extremely important when we study relative motion, frames of reference, and later topics. For now, just remember: to find where A is relative to B, subtract B's position from A's position.

---

## 16. Common Mistakes and How to Avoid Them

Over many years of teaching physics, certain mistakes come up again and again. Here is a guide to avoiding them:

### Mistake 1: Using Distance When You Should Use Displacement

**Wrong thinking:** "The car traveled 100 km, so its displacement is 100 km."

**Why it is wrong:** The car may have driven a curved route, or doubled back. Displacement is only the straight-line separation from start to end.

**Fix:** Ask yourself: did the object travel in a perfectly straight line without backtracking? If yes, distance = displacement magnitude. If no, they are different.

---

### Mistake 2: Forgetting Units

**Wrong answer:** "Speed = 60"

**Why it is wrong:** 60 what? 60 m/s is a race car. 60 km/h is a city driver. They are completely different.

**Fix:** Always write units in every step of your calculation. Units are part of the answer.

---

### Mistake 3: Mixing Units in One Calculation

**Wrong calculation:**
```
Speed = Distance / Time = 10 km / 30 seconds = 0.33 ??? 
```
(What is km/second? Not a standard unit!)

**Fix:** Convert everything to the same system before calculating.
```
Convert 30 seconds to hours: 30/3600 = 0.00833 hours
Speed = 10 km / 0.00833 h = 1200 km/h

OR:

Convert 10 km to metres: 10,000 m
Speed = 10,000 m / 30 s = 333 m/s
```

---

### Mistake 4: Adding Speeds Directly for Average Speed

**Wrong calculation:** 
"I drove at 40 km/h for the first part and 60 km/h for the second part. Average = (40+60)/2 = 50 km/h."

**Why it is wrong:** This only works if both parts took equal time. In general, you must use total distance / total time.

**Fix:** Always use v_avg = total distance / total time.

---

### Mistake 5: Thinking Displacement Can Never Be Negative

**Wrong thinking:** "Displacement is always positive, just like distance."

**Why it is wrong:** Displacement has direction. If you walk 5 metres in the negative direction (West, if East is positive), your displacement is -5m.

**Fix:** Remember that displacement = final position − initial position. If final position is less than initial position (on a number line), displacement is negative.

---

### Mistake 6: Confusing Speed and Velocity

These are related but not the same:
- **Speed** = distance / time (scalar, no direction)
- **Velocity** = displacement / time (vector, has direction)

A car going around a circular track at 60 km/h has constant speed but changing velocity (because direction keeps changing).

---

## 17. Summary

You have covered a lot of ground in this chapter — quite literally! Here are the key takeaways:

### Core Concepts

- **Distance** is the total length of the path you travel. It is always positive and keeps adding up no matter which direction you move. It is a **scalar** quantity.

- **Displacement** is the straight-line distance from your starting point to your ending point, measured in a specific direction. It can be positive, negative, or zero. It is a **vector** quantity.

- The difference: distance is about the journey (every step of the path), displacement is about the outcome (where you ended up compared to where you started).

- Displacement can be zero (if you return to your starting point) while distance is always greater than zero for any motion.

- A **reference point** (origin) is needed to define positions. All positions are measured relative to the origin.

- **Position** is described using coordinates (x in 1D, x and y in 2D).

- **Speed** is the rate at which distance is covered: Speed = Distance / Time.

- The standard unit for speed in science is **m/s**. The common everyday unit is **km/h**.

- To convert: km/h ÷ 3.6 = m/s, and m/s × 3.6 = km/h.

- **Average speed** = total distance ÷ total time (for a whole journey).

- **Instantaneous speed** = speed at one specific moment (what a speedometer shows).

- **Scalars** have magnitude only (distance, speed, time, mass, temperature).

- **Vectors** have magnitude and direction (displacement, velocity, force).

- **Velocity** is the vector version of speed: Velocity = Displacement / Time.

- **Relative position** of A with respect to B = position of A − position of B.

- An **odometer** measures distance (total path). A **GPS** uses both — displacement for "how far to destination" and distance for "how far you have traveled on the route."

---

## 18. Key Equations

Here are all the important formulas from this chapter, collected in one place for easy reference:

### Speed, Distance, and Time

```
Speed = Distance / Time
v = d / t

Distance = Speed × Time
d = v × t

Time = Distance / Speed
t = d / v
```

### Average Speed

```
Average Speed = Total Distance / Total Time
v_avg = d_total / t_total
```

### Displacement (1D)

```
Displacement = Final Position - Initial Position
Δx = x_final - x_initial
```

### Displacement (2D, Pythagorean Theorem)

```
Displacement magnitude = √(Δx² + Δy²)
(where Δx = horizontal change, Δy = vertical change)
```

### Unit Conversion

```
km/h to m/s:   divide by 3.6
m/s to km/h:   multiply by 3.6

Exact factor:
1 km/h = 1000/3600 m/s = 5/18 m/s ≈ 0.2778 m/s
1 m/s  = 3600/1000 km/h = 18/5 km/h = 3.6 km/h
```

### Relative Position

```
Position of A relative to B = Position of A - Position of B
x_A_rel_B = x_A - x_B
```

### Average Velocity (preview)

```
Average Velocity = Displacement / Time
v_avg = Δx / t
```

---

*You have now built a solid foundation in the language of motion. In the next chapter, we will introduce velocity (the vector form of speed) and acceleration (how speed changes over time), taking everything you learned here to the next level.*
