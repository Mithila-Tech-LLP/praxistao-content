# Chapter 15: Friction and Drag

> **"Without friction, you couldn't walk, drive, or hold a pen. With too much of it, machines overheat and wear out. Friction is one of physics's great double-edged swords."**

---

## Table of Contents

1. [What Is This Chapter About?](#1-what-is-this-chapter-about)
2. [What Is Friction?](#2-what-is-friction)
3. [The Microscopic Picture — Why Surfaces Have Friction](#3-the-microscopic-picture--why-surfaces-have-friction)
4. [The Normal Force and Why It Matters for Friction](#4-the-normal-force-and-why-it-matters-for-friction)
5. [Static Friction — Preventing Motion from Starting](#5-static-friction--preventing-motion-from-starting)
6. [Kinetic Friction — Opposing Sliding Motion](#6-kinetic-friction--opposing-sliding-motion)
7. [The Coefficients of Friction](#7-the-coefficients-of-friction)
8. [Key Properties of Kinetic Friction](#8-key-properties-of-kinetic-friction)
9. [Worked Example 1: Box on a Floor — Kinetic Friction](#9-worked-example-1-box-on-a-floor--kinetic-friction)
10. [Worked Example 2: Minimum Force to Start Moving](#10-worked-example-2-minimum-force-to-start-moving)
11. [Friction on an Inclined Plane](#11-friction-on-an-inclined-plane)
12. [Worked Example 3: Block on a Ramp — Does It Slide?](#12-worked-example-3-block-on-a-ramp--does-it-slide)
13. [The Critical Angle for Sliding](#13-the-critical-angle-for-sliding)
14. [Air Resistance — Friction From the Air](#14-air-resistance--friction-from-the-air)
15. [Terminal Velocity — When Drag Equals Gravity](#15-terminal-velocity--when-drag-equals-gravity)
16. [Worked Example 4: Estimating Terminal Velocity](#16-worked-example-4-estimating-terminal-velocity)
17. [How Parachutes Work](#17-how-parachutes-work)
18. [Useful Friction vs Harmful Friction](#18-useful-friction-vs-harmful-friction)
19. [Ways to Reduce Friction](#19-ways-to-reduce-friction)
20. [Common Mistakes to Avoid](#20-common-mistakes-to-avoid)
21. [Practice Problems](#21-practice-problems)
22. [Summary](#22-summary)
23. [Key Equations](#23-key-equations)

---

## 1. What Is This Chapter About?

Every time you take a step, write with a pen, drive a car, or throw a ball — friction and drag are at work. They are the forces that slow things down, generate heat, wear out surfaces, and yet, paradoxically, make most motion possible in the first place.

Try to imagine a world without friction:
- Your feet would slip on the floor with every step.
- Your pen would skate across the paper without making a mark.
- Cars could not brake. Nails and screws would not hold.
- Even sitting on a chair would be dangerous — you'd slide off.

And yet, too much friction causes problems:
- Engines overheat and wear out.
- Aircraft burn extra fuel fighting air resistance.
- Machinery grinds to a halt.

Understanding friction means understanding when to increase it (tyres, brakes, grips) and when to reduce it (engines, aircraft, ball bearings).

By the end of this chapter you will be able to:

- Explain what friction is and where it comes from at the microscopic level
- Distinguish between static and kinetic friction
- Use the equations f_s_max = μ_s × N and f_k = μ_k × N to solve problems
- Analyse friction on flat surfaces and inclined planes
- Explain air resistance and terminal velocity
- Calculate how a parachute lowers terminal velocity

---

## 2. What Is Friction?

**Friction** is the force that opposes the relative motion — or the tendency of relative motion — between two surfaces that are in contact.

Key words in that definition:

- **Opposes**: friction always acts in the direction that resists what the surfaces are trying to do relative to each other.
- **Relative motion**: what matters is how the surfaces move relative to EACH OTHER, not how either one moves relative to the ground.
- **Tendency of motion**: even when nothing is moving, if a force is trying to push one surface over another, friction resists that too.

Friction is a **contact force** — it only exists where the two surfaces actually touch. Remove the contact and friction vanishes.

```
EXAMPLES OF FRICTION:

  Situation                      | Direction of friction force
  -------------------------------|-------------------------------
  Box pushed right across floor  | Acts LEFT on box (opposes motion)
  Car tyre rolling forward       | Acts FORWARD on car (enables motion!)
  Hand sliding on a rope         | Acts opposite to sliding direction
  Falling skydiver in air        | Acts UPWARD on skydiver (air on body)
  Book about to slide off table  | Prevents sliding (static friction)
```

Note the third row: for a driving wheel, friction from the road acts **forward** on the car. Without this friction, the wheel would just spin in place. Friction is what converts wheel rotation into forward motion. That is why friction is not always harmful.

---

## 3. The Microscopic Picture — Why Surfaces Have Friction

When you look at a smooth surface — polished metal, glass, a plastic table — it looks flat and featureless. But zoom in to the microscopic scale (millionths of a metre) and every surface looks like a mountain range.

```
MACROSCOPIC VIEW:
  ___________________
  ___________________   Looks perfectly flat

MICROSCOPIC VIEW (hugely magnified):
     /\/\/\/\/\/\/\
  __/ \  /\ / \/  \_   Peaks and valleys everywhere!
  ____\/  \/   \___
```

When two surfaces are pressed together, the tiny peaks (called **asperities**) on one surface interlock with the peaks of the other surface. This mechanical interlocking creates resistance to sliding.

In addition, at the atomic level, molecules of one surface are attracted to molecules of the other surface. This **adhesion** also contributes to friction.

```
INTERLOCKING ASPERITIES:

   Surface 1: \/\/\/\/\/\/\/\/\
              /\/\/\/\/\/\/\/\/\
   Surface 2:

   The peaks of one surface fit into the valleys of the other.
   To slide, you must either:
     (a) lift one surface over the other's peaks (requires energy), OR
     (b) wear down the peaks (creates heat and debris)
```

This is why:
- **Rougher surfaces have more friction** — more interlocking.
- **Heavier loads have more friction** — more force presses peaks together.
- **Lubrication reduces friction** — oil fills the gaps and separates the surfaces, preventing interlocking.
- **Surfaces wear out from friction** — the asperities break off over time.

---

## 4. The Normal Force and Why It Matters for Friction

Before we write down the friction equations, we need to understand the **normal force** again — because friction depends on it directly.

The **normal force** (N) is the contact force a surface exerts perpendicular (normal = at 90°) to itself on the object resting on it. It represents how hard the two surfaces are pressed together.

**On a flat horizontal surface:**

```
              N (upward, perpendicular to floor)
              ↑
              |
   -----------+----------
   |   Box    |          mass = m
   -----------+----------
              |
              ↓
           W = mg

For a stationary box, ΣF_y = 0:
N = mg
```

The harder the surfaces press against each other (higher N), the more the asperities interlock, and the larger the friction force.

**On an inclined plane at angle θ:**

```
                       |  (vertical)
                       |
              N ───────+  (perpendicular to ramp surface)
             /         |
            /          |
           /  θ        |
          /            |
         [BLOCK]       |
        / |            |
       /  ↓            |
      /  W = mg (straight down)

The component of W perpendicular to ramp:
   W_perp = mg × cos θ

The normal force: N = mg × cos θ   (not mg!)

The component of W along the ramp (trying to make it slide down):
   W_parallel = mg × sin θ
```

This is why the normal force (and therefore friction) decreases as the ramp gets steeper — less of the weight presses perpendicularly into the surface.

---

## 5. Static Friction — Preventing Motion from Starting

**Static friction** (f_s) is the friction force that acts between two surfaces that are NOT sliding relative to each other. It prevents motion from starting.

Think about pushing a heavy filing cabinet. You push lightly — it does not move. Push harder — still no movement. Push even harder — it JUST starts to slide.

What is happening? As you push harder, static friction increases to match your push exactly — right up to a maximum value. Once you exceed that maximum, the cabinet starts sliding.

```
STATIC FRICTION GRAPH:

  Applied Force (N)
  ^
  |      f_s matches F_applied
  |      exactly (no motion)
  |                  ← Maximum static friction
  f_s_max - - - - - - - - *──────────────────
  |                  *     \  kinetic friction (lower)
  |               *         \___________________
  |            *
  |         *
  |      *
  |   *
  +─────────────────────────────────────────→ Applied Force
                    ^
                    Motion begins here
```

**Key behaviour:**
- When F_applied < f_s_max: f_s = F_applied (they are equal, box stays still)
- When F_applied = f_s_max: the box is on the verge of sliding
- When F_applied > f_s_max: the box starts sliding (now kinetic friction applies)

**The formula for maximum static friction:**

```
f_s_max = μ_s × N

where:
  f_s_max = maximum static friction force (N)
  μ_s     = coefficient of static friction (dimensionless, no units)
  N       = normal force (N)
```

The actual static friction at any moment can be anywhere from 0 to f_s_max:

```
0 ≤ f_s ≤ μ_s × N
```

If F_applied = 10 N and f_s_max = 20 N, then f_s = 10 N. The box stays still. If F_applied increases to 25 N, the box slides.

---

## 6. Kinetic Friction — Opposing Sliding Motion

**Kinetic friction** (f_k, also called sliding friction or dynamic friction) acts when two surfaces ARE sliding relative to each other. It opposes the sliding motion.

Unlike static friction, kinetic friction is not a variable — it has a fixed value for given surfaces at a given normal force:

```
f_k = μ_k × N

where:
  f_k = kinetic friction force (N)
  μ_k = coefficient of kinetic friction (dimensionless)
  N   = normal force (N)
```

Kinetic friction is always in the direction **opposite to the direction of motion** of the sliding surface.

```
BOX SLIDING RIGHT ACROSS FLOOR:

              N
              ↑
              |
  f_k ←──────+────────→ F_applied
  (friction  [BOX]
  opposes
  motion)
              |
              ↓
           W = mg

f_k = μ_k × N = μ_k × mg   (on flat surface)

F_net = F_applied - f_k
a = F_net / m
```

**The crucial fact:** Kinetic friction is almost always **less than maximum static friction**.

```
μ_k < μ_s    (always)

This means:
  It takes MORE force to START sliding something
  than to KEEP it sliding once it's already moving.
```

You experience this when you rearrange heavy furniture. It takes a big initial push to get it moving, but less effort to keep it moving once it's going. That initial peak is f_s_max. Once it breaks free, you only fight f_k, which is lower.

---

## 7. The Coefficients of Friction

The **coefficient of friction** (μ, the Greek letter "mu") is a dimensionless number that describes how "grippy" or "slippery" two surfaces are with each other. A higher μ means more friction.

μ depends on:
- The materials of both surfaces
- The surface finish (rough vs. polished)
- Whether surfaces are dry, wet, or lubricated
- NOT on the area of contact (this surprises people — see Section 8)

**Typical values of μ_s and μ_k:**

```
Material Pair             |  μ_s (static)  |  μ_k (kinetic)
--------------------------|----------------|----------------
Rubber on dry concrete    |   0.6 – 0.8    |   0.5 – 0.7
Rubber on wet concrete    |   0.4 – 0.6    |   0.35 – 0.5
Wood on wood (dry)        |   0.25 – 0.5   |   0.2 – 0.4
Wood on wood (oiled)      |   0.1 – 0.2    |   0.05 – 0.15
Steel on steel (dry)      |   0.15 – 0.6   |   0.1 – 0.4
Steel on steel (lubricated)|  0.05 – 0.1   |   0.03 – 0.08
Ice on ice                |   0.03 – 0.1   |   0.01 – 0.05
Teflon on steel           |   0.04 – 0.06  |   0.04 – 0.05
Bone on cartilage (joints)|   0.005 – 0.01 |   0.003 – 0.01
```

Note:
- Rubber on dry concrete has a high μ — essential for car tyres gripping the road.
- Ice on ice has a very low μ — why hockey pucks slide so far.
- Teflon (PTFE) has extremely low friction — used in non-stick pans and industrial bearings.
- Bone-on-cartilage has the lowest friction of any natural or man-made material — your knee and hip joints are incredibly slippery thanks to synovial fluid.

**μ_s > μ_k always**: The coefficient of static friction is always greater than kinetic friction for the same surface pair. This is a universal rule.

---

## 8. Key Properties of Kinetic Friction

These three properties define kinetic friction and are somewhat counter-intuitive:

### Property 1: Kinetic friction is proportional to the normal force

Double the normal force → double the kinetic friction.

```
f_k = μ_k × N

If N doubles: f_k = μ_k × (2N) = 2 × μ_k × N  (doubled)

Physically: pressing harder → more asperity interlocking → more friction
```

### Property 2: Kinetic friction is independent of surface area

This surprises most people. A wide flat box and a narrow tall box with the same mass and same surface material have the SAME friction force on the same floor.

```
Why? Larger area → each unit of area carries less force (less asperity contact per spot)
     Smaller area → each unit of area carries more force (more contact per spot)
     The effects cancel out. Only the total N matters.

     Example:
     Brick lying flat (large area):   N = mg,  f_k = μ_k × mg
     Brick standing on end (small area): N = mg,  f_k = μ_k × mg
     SAME friction force!
```

This is why wider tyres do NOT necessarily grip better (in dry conditions). The advantage of wide tyres is different — they can handle heat better and provide more contact area on uneven surfaces. But for dry friction on flat surfaces, width alone does not change the friction force.

### Property 3: Kinetic friction is approximately independent of speed

Over a normal range of speeds, kinetic friction force does not change with how fast the surfaces are sliding. (At very high speeds or extreme conditions, this breaks down, but for everyday physics it holds.)

```
Sliding slowly:   f_k = μ_k × N
Sliding fast:     f_k ≈ μ_k × N   (same!)

This is an approximation — at very high speeds or temperatures,
μ_k can change, especially for rubber on road.
```

---

## 9. Worked Example 1: Box on a Floor — Kinetic Friction

**Problem:** A 20 kg box sits on a wooden floor. The coefficient of kinetic friction between the box and floor is μ_k = 0.3. A horizontal force of 80 N is applied to the box to the right. Find: (a) the friction force, (b) the net force, and (c) the acceleration of the box.

**Step 1: Draw the Free Body Diagram**

```
              N
              ↑
              |
f_k ←─────────+──────────────→ F_A = 80 N
              | 20 kg         |
              └───────────────┘
              |
              ↓
           W = mg = 20 × 9.8 = 196 N

(Positive direction: → to the right)
(The box IS moving — so kinetic friction applies)
```

**Step 2: Find the Normal Force**

Vertical direction (no vertical acceleration):
```
ΣF_y = 0
N - W = 0
N = W = mg = 20 × 9.8 = 196 N
```

**Step 3: Find the Kinetic Friction Force**

```
f_k = μ_k × N
f_k = 0.3 × 196
f_k = 58.8 N   (to the left, opposing motion)
```

**Step 4: Find the Net Force**

```
ΣF_x = F_A - f_k
ΣF_x = 80 - 58.8
ΣF_x = 21.2 N   (to the right)
```

**Step 5: Find Acceleration**

```
a = F_net / m
a = 21.2 / 20
a = 1.06 m/s²
```

**Answers:**
- Friction force: **58.8 N** (to the left)
- Net force: **21.2 N** (to the right)
- Acceleration: **1.06 m/s²** (to the right)

**Sense check:** Without friction, a = 80/20 = 4 m/s². With friction eating up most of the force, we get 1.06 m/s². Makes sense — friction significantly reduces the acceleration.

---

## 10. Worked Example 2: Minimum Force to Start Moving

**Problem:** The same 20 kg box is now at rest. The coefficient of static friction is μ_s = 0.4. What is the minimum horizontal force needed to START the box moving?

**Explanation:** The box will start moving when the applied force JUST exceeds the maximum static friction. So we need F_applied = f_s_max.

**Step 1: Find Normal Force (same as before)**

```
N = mg = 20 × 9.8 = 196 N
```

**Step 2: Calculate Maximum Static Friction**

```
f_s_max = μ_s × N
f_s_max = 0.4 × 196
f_s_max = 78.4 N
```

**Step 3: Minimum force to start moving**

```
F_minimum = f_s_max = 78.4 N
```

**Answer:** You need at least **78.4 N** to START the box moving.

**Compare the two examples:**
- To START moving (static): need > 78.4 N
- Once moving (kinetic): friction drops to 58.8 N

This means: you need 78.4 N to break it free, then only 58.8 N to keep it moving at constant velocity, and 80 N gives an acceleration of 1.06 m/s² while sliding.

This matches the everyday experience: getting something moving takes more effort than keeping it moving.

---

## 11. Friction on an Inclined Plane

Inclined planes (ramps) add an extra step: you must first resolve the weight into components along and perpendicular to the slope.

```
BLOCK ON A RAMP AT ANGLE θ:

                         N (perpendicular to ramp)
                         ↑
                        /
               f_s ← [BLOCK]
                      /  |  \
                     /   |   \
                    /    ↓    \
                   /   mg      \
                  θ             \
                 /                \
────────────────/──────────────────\──────

The weight mg acts STRAIGHT DOWN.
We split it into two components:
  - mg sin θ : component ALONG the ramp (tries to make block slide DOWN)
  - mg cos θ : component PERPENDICULAR to ramp (presses block into ramp)
```

**The component diagram:**

```
                     (perpendicular to ramp)
                          ↑
                          | mg cos θ
                          |
        mg (straight ─────+─────→ mg sin θ (along ramp, downward along slope)
            down)         |
                          |
                          ↓
                     (into ground)

By geometry (angle θ between the ramp and horizontal):
  Component along ramp:          mg sin θ
  Component perpendicular to ramp: mg cos θ
```

**Normal Force on the ramp:**

Since there is no acceleration perpendicular to the ramp:

```
ΣF_perpendicular = 0
N - mg cos θ = 0
N = mg cos θ
```

Notice: N = mg cos θ, which is LESS than mg (because cos θ < 1 for θ > 0). As the ramp gets steeper, N decreases.

**Friction force on the ramp:**

```
f_k = μ_k × N = μ_k × mg cos θ   (if sliding)
f_s_max = μ_s × N = μ_s × mg cos θ   (if stationary)
```

**Net force along the ramp (if sliding down):**

```
F_net = mg sin θ - f_k
F_net = mg sin θ - μ_k × mg cos θ
F_net = mg(sin θ - μ_k cos θ)

If this is positive: block accelerates down the ramp
If this is zero: block slides at constant velocity
If F_applied < f_s_max: block stays still (static friction holds it)
```

---

## 12. Worked Example 3: Block on a Ramp — Does It Slide?

**Problem:** A 5 kg block sits on a ramp inclined at 30° to the horizontal. The coefficient of static friction between block and ramp is μ_s = 0.4. Does the block slide?

**Step 1: Draw the Free Body Diagram for the block on the ramp**

```
               N
               ↑ (perpendicular to ramp)
               |
        f_s ←──+──────────────────────────
               | [BLOCK]                  \
               |/                    θ=30° \
              /↓\                           \
         mg cos θ  mg sin θ →                \
                   (down slope)               \

(We're checking if static friction can hold the block.
 If mg sin θ ≤ f_s_max, block stays. Otherwise it slides.)
```

**Step 2: Calculate weight components**

```
mg = 5 × 9.8 = 49 N

Component down the slope:
   W_parallel = mg sin 30° = 49 × 0.5 = 24.5 N

Component into the ramp:
   W_perp = mg cos 30° = 49 × 0.866 = 42.4 N
```

**Step 3: Find Normal Force**

```
N = W_perp = mg cos 30° = 42.4 N
```

**Step 4: Find Maximum Static Friction**

```
f_s_max = μ_s × N = 0.4 × 42.4 = 16.96 N
```

**Step 5: Compare**

```
Force trying to slide block down:  W_parallel = 24.5 N
Maximum friction holding it:       f_s_max    = 16.96 N

Since 24.5 N > 16.96 N, static friction CANNOT hold the block.

The block SLIDES!
```

**Bonus: Find the acceleration once sliding begins.**

Using kinetic friction (assume μ_k = 0.3):

```
f_k = μ_k × N = 0.3 × 42.4 = 12.72 N

F_net = mg sin θ - f_k = 24.5 - 12.72 = 11.78 N

a = F_net / m = 11.78 / 5 = 2.36 m/s² (down the slope)
```

**Answer:** The block slides. Once sliding, it accelerates down the ramp at **2.36 m/s²**.

---

## 13. The Critical Angle for Sliding

There is a neat result: you can find the critical angle θ_c at which a block JUST begins to slide, and it gives you the coefficient of static friction directly.

**Derivation:**

At the critical angle, the component down the slope exactly equals maximum static friction:

```
mg sin θ_c = f_s_max = μ_s × N = μ_s × mg cos θ_c

Divide both sides by mg cos θ_c:

sin θ_c / cos θ_c = μ_s

tan θ_c = μ_s
```

So:

```
┌────────────────────────────┐
│   tan θ_c = μ_s            │
│   θ_c = arctan(μ_s)        │
└────────────────────────────┘
```

**This is a remarkably useful result:** if you slowly tilt a surface with an object on it and measure the angle at which it just starts to slide, you have directly measured the coefficient of static friction!

**Examples:**

```
μ_s = 0.4:    θ_c = arctan(0.4) = 21.8°
μ_s = 0.577:  θ_c = arctan(0.577) = 30°
μ_s = 1.0:    θ_c = arctan(1.0) = 45°
μ_s = 1.7:    θ_c = arctan(1.7) = 60°

Some surfaces (rubber on rough concrete) have μ > 1,
meaning the critical angle exceeds 45°.
```

**Quick reference table:**

```
θ (degrees) | tan θ | This is μ_s for which objects start sliding at this angle
------------|-------|----------------------------------------------------------
     10°    |  0.18 | Very slippery (ice, lubricated metal)
     15°    |  0.27 | Slippery floor
     22°    |  0.40 | Typical wood on wood
     30°    |  0.58 | Moderately rough surface
     45°    |  1.00 | Very rough surface
```

---

## 14. Air Resistance — Friction From the Air

So far we have looked at friction between solid surfaces. But fluids (liquids and gases) also exert a resistive force on objects moving through them. This is called **drag** or **air resistance** (when the fluid is air).

**What is drag?** When an object moves through air, it pushes air molecules out of the way. By Newton's Third Law, those air molecules push back on the object — this is drag.

```
OBJECT MOVING THROUGH AIR:

           →→→ Object moves this way →→→

      Air pushed aside          Air pushed aside
          ↗                           ↘
    ──────────────────────────────────────────
    ||||    →→→  [OBJECT]  →→→          |||||
    ──────────────────────────────────────────
          ↘                           ↗
      Air pushed aside          Air pushed aside

   Drag force acts BACKWARD on the object (opposing motion).
   The object pushes air forward/sideways.
   Air pushes object backward (Newton 3rd Law).
```

**The drag force equation:**

```
F_drag = ½ × ρ × C_d × A × v²

where:
  F_drag = drag force (N)
  ρ      = density of air (kg/m³) ≈ 1.2 kg/m³ at sea level
  C_d    = drag coefficient (dimensionless) — depends on shape
  A      = cross-sectional area facing the direction of motion (m²)
  v      = speed of object (m/s)
  ½      = a constant from the physics of fluid flow
```

**Key features of drag:**

1. **Drag increases with speed SQUARED** (v²). Double the speed → four times the drag. Triple the speed → nine times the drag. This is why cars need far more fuel at highway speeds than in city traffic.

2. **Drag depends on shape** (C_d). A streamlined shape has a low C_d. A blunt shape has a high C_d.

```
Typical drag coefficients:

Shape                     |  C_d
--------------------------|------
Sphere                    |  0.47
Cube (facing flat side)   |  1.05
Streamlined car           |  0.25 – 0.35
Racing car (with wings)   |  0.7 – 1.0
Cyclist (upright)         |  0.9 – 1.1
Cyclist (racing crouch)   |  0.7 – 0.8
Skydiver (spread eagle)   |  1.0 – 1.3
Skydiver (head-down)      |  0.7 – 0.8
```

3. **Drag depends on cross-sectional area** (A). A larger object facing the airflow experiences more drag. This is why parachutes work: they enormously increase A.

4. **Drag depends on fluid density** (ρ). Air is much less dense than water (about 800 times less dense). This is why drag in water is so much stronger than drag in air — why swimmers feel much more resistance than runners.

---

## 15. Terminal Velocity — When Drag Equals Gravity

When an object falls through air, gravity pulls it downward and drag pushes it upward. Initially (when v is small), gravity wins and the object accelerates. But as v increases, drag increases (proportional to v²). Eventually drag equals gravity and there is no more net force — the object stops accelerating and falls at constant velocity.

This constant velocity is called **terminal velocity**.

```
THE JOURNEY TO TERMINAL VELOCITY:

Just released (v = 0):
   ↓ W = mg
   ↑ F_drag ≈ 0 (no movement yet)
   F_net = mg downward → large acceleration

Falling at medium speed:
   ↓ W = mg
   ↑ F_drag = ½ρC_dAv² (some drag, smaller than W)
   F_net = mg - F_drag → some acceleration (still speeding up)

At terminal velocity (v = v_t):
   ↓ W = mg
   ↑ F_drag = mg  (drag NOW equals weight!)
   F_net = 0 → zero acceleration → constant velocity!

             ┌──────────────────────────────────────────┐
             │  At terminal velocity: F_drag = W = mg   │
             │  ½ × ρ × C_d × A × v_t² = mg            │
             └──────────────────────────────────────────┘
```

**Speed vs time graph for a falling object:**

```
Speed (m/s)
^
|                                                    ─────── v_terminal
|                                          ─────────/
|                                   ──────/
|                              ────/
|                         ────/
|                    ────/
|               ────/
|          ────/
|     ────/
|────/
+──────────────────────────────────────────────────────→ Time (s)
 0

The curve rises steeply at first (large net force, large acceleration)
then levels off as drag increases
and approaches terminal velocity asymptotically.
```

**Typical terminal velocities:**

```
Object               |  Terminal Velocity (approximate)
---------------------|------------------------------------
Large raindrop       |  9 m/s  (about 32 km/h)
Skydiver (spread)    |  55 m/s (about 200 km/h)
Skydiver (head-down) |  90 m/s (about 320 km/h)
Skydiver (with open parachute) | 5–6 m/s (safe landing speed!)
A falling ant        |  ~1 m/s (essentially unharmed by falls!)
A mouse              |  ~10 m/s (survives most falls)
A cat                |  ~20 m/s (can survive with injury)
A human without parachute | ~55 m/s (fatal impact)
```

**Why ants survive falls from any height:** An ant is so light and small that its terminal velocity is only about 1 m/s. Air resistance brings it to a gentle landing even after a fall from a tall building.

**Why large animals are more vulnerable:** Larger animals have more mass (weight scales with volume, i.e., with length³) but cross-sectional area only scales with length². As animals get bigger, weight grows faster than drag. Terminal velocity therefore increases with size.

---

## 16. Worked Example 4: Estimating Terminal Velocity

**Problem:** A skydiver has a mass of 80 kg and spreads out in a stable position so that their cross-sectional area is A = 0.7 m². The drag coefficient is C_d = 1.0. Take air density as ρ = 1.2 kg/m². Find the terminal velocity.

**Step 1: At terminal velocity, drag equals weight**

```
F_drag = W
½ × ρ × C_d × A × v_t² = mg
```

**Step 2: Substitute values**

```
½ × 1.2 × 1.0 × 0.7 × v_t² = 80 × 9.8
0.42 × v_t² = 784
v_t² = 784 / 0.42
v_t² = 1866.7
v_t = √1866.7
v_t ≈ 43 m/s
```

**Answer:** The terminal velocity is approximately **43 m/s** ≈ 155 km/h.

(Real skydivers typically report 50–60 m/s in spread position — the slight discrepancy is because we simplified the body shape and assumed uniform air density.)

**What if the skydiver curls into a ball (smaller area, lower C_d)?**

Say A = 0.2 m², C_d = 0.8:

```
½ × 1.2 × 0.8 × 0.2 × v_t² = 784
0.096 × v_t² = 784
v_t² = 8166.7
v_t = √8166.7 ≈ 90 m/s
```

More than double the terminal velocity! Skydivers in "head-down" position can reach ~90 m/s (300+ km/h) this way.

---

## 17. How Parachutes Work

A parachute dramatically reduces terminal velocity by increasing the cross-sectional area A and increasing the drag coefficient C_d.

```
WITHOUT PARACHUTE:
                    [Skydiver]
                        ↓  W = mg
                        ↑  F_drag (small A, small C_d)
                    Terminal velocity ~ 55 m/s  ← fatal impact

WITH PARACHUTE:
              ─────────────────────
             /  PARACHUTE (huge A)  \
            /    C_d ≈ 1.3           \
           /─────────────────────────\
                        |
                    [Skydiver]
                        ↓  W = mg (same as before)
                        ↑  F_drag (MUCH larger — large A)
                    Terminal velocity ~ 5–6 m/s  ← safe landing
```

**Why parachutes have holes:** A small hole in the top of the parachute allows some air through. This actually increases stability — without the hole, the parachute would oscillate and tumble wildly. The slight reduction in drag area is worth the massive gain in stability.

**The numbers:**

```
Without parachute:
  A ≈ 0.7 m², C_d ≈ 1.0
  v_t ≈ 55 m/s

With round parachute (diameter ~8 m):
  A = π × r² = π × 4² ≈ 50 m²
  C_d ≈ 1.3 (parachutes are blunt — high C_d)

  ½ × 1.2 × 1.3 × 50 × v_t² = 80 × 9.8
  39 × v_t² = 784
  v_t² = 20.1
  v_t ≈ 4.5 m/s  ← gentle walking speed
```

A parachute turns a potentially fatal 55 m/s into a survivable 4.5 m/s by increasing the effective area from 0.7 m² to 50 m² — a factor of ~70 in area.

Since v_t ∝ 1/√A (from the terminal velocity formula), a 70× increase in area gives:
```
v_t_new / v_t_old = 1 / √70 ≈ 1 / 8.4

55 / 8.4 ≈ 6.5 m/s  ← matches our calculation!
```

---

## 18. Useful Friction vs Harmful Friction

Friction is not inherently good or bad — it depends entirely on the application.

### Useful Friction

**Walking and running:**
Without friction between shoes and ground, walking would be impossible. Ice and polished marble floors feel dangerous because their low μ removes this friction.

**Car braking:**
Braking pads press against rotating discs. Friction converts kinetic energy into heat, slowing the car. Anti-lock braking systems (ABS) prevent wheels locking up because a spinning tyre (kinetic friction) has less braking force than a tyre on the verge of slipping (static friction). μ_s > μ_k, so keeping the tyre just below lock-up maximises braking force.

**Writing with a pen:**
Friction between pen tip and paper allows ink to transfer. Try writing on glass — the pen slips and barely marks it.

**Holding and gripping:**
Friction between fingers and objects allows you to grip cups, tools, and bottles. Grip tape on handles increases μ to prevent slipping.

**Nails and screws:**
Nails rely on friction between nail shaft and wood. Screws use both friction and mechanical interlocking via threads.

**Tyres on roads:**
Friction between tyres and road allows cornering, acceleration, and braking. Wet roads reduce μ by providing a lubricating film between tyre and road — reason for lower speed limits in rain.

### Harmful Friction

**Engine wear:**
Inside a car engine, metal surfaces (pistons, crankshaft, valves) slide against each other at high speed and temperature. Without lubrication, they would quickly grind down to dust. Engine oil maintains a thin film between surfaces.

**Energy loss in machines:**
Every bearing, gear, belt, and shaft in a machine loses energy to friction as heat. High-performance machines (spacecraft, precision instruments) go to great lengths to minimise friction losses.

**Brake heat and fade:**
Braking converts kinetic energy to heat in the brake pads. In heavy braking (racing cars, long downhill grades), pads can overheat, reducing μ and causing "brake fade" — the car cannot stop as quickly.

**Wear on surfaces:**
Any sliding surfaces eventually wear down. Roads degrade from tyre friction. Running shoe soles thin out. Industrial machinery must be regularly maintained.

---

## 19. Ways to Reduce Friction

Engineers have developed many techniques to reduce friction where it is harmful.

### Lubrication

A thin layer of oil, grease, or other liquid between surfaces prevents direct contact between asperities. The surfaces now slide on the fluid layer rather than on each other. Effective lubrication can reduce μ by a factor of 10 or more.

```
WITHOUT OIL:          WITH OIL:
 /\/\/\/\/\/\          /\/\/\/\/\
 \/\/\/\/\/\/          __________ (oil layer)
 /\/\/\/\/\/\          /\/\/\/\/\
Asperities interlock  Asperities float apart
High friction          Low friction
```

**Types of lubricants:**
- Oils (mineral or synthetic): engine oil, hydraulic fluid
- Greases: oils thickened with a soap agent, used in wheel bearings
- Dry lubricants: graphite, Teflon (PTFE) powder — used where liquid lubricants would contaminate (food processing, vacuum systems)
- Air bearings: a thin layer of pressurised air completely separates the surfaces (used in precision instruments, hard drives)

### Ball Bearings and Roller Bearings

Instead of sliding friction (surface on surface), rolling elements (balls or rollers) convert friction into rolling — which has much less resistance.

```
SLIDING:       ROLLING (ball bearing):
 ──────────       ─────────────────
 ─────↔─────       O O O O O O O O
 ──────────       ─────────────────

 High sliding friction   Balls roll between surfaces
                         Much lower friction
```

### Streamlining

Streamlined shapes (teardrop, aerofoil) have a much lower drag coefficient C_d than blunt shapes. Air flows smoothly around them rather than creating turbulent eddies.

```
BLUNT SHAPE:                STREAMLINED SHAPE:
   ┌─────┐                      ╱─────────╲
   │     │  lots of turbulence  │ smooth    ╲
   │     │  →→→XXXX             │ flow       ╲──→
   └─────┘                      ╲           ╱
                                  ╲────────╱
   C_d ≈ 1.0                      C_d ≈ 0.04
```

### Surface Finishing

Polishing surfaces to a mirror finish reduces the number and height of asperities, reducing friction. Used in precision bearings, hydraulic cylinders, and piston bores.

### Materials with Low μ

Teflon (PTFE, polytetrafluoroethylene) has one of the lowest friction coefficients of any solid material (μ ≈ 0.04). Used in:
- Non-stick cookware
- Medical implants
- Industrial liners
- Plumbing tape

---

## 20. Common Mistakes to Avoid

### Mistake 1: Using f = μ × mg instead of f = μ × N

On a flat horizontal surface, N = mg, so this works. But on a ramp, N = mg cos θ. Always find N first, then calculate friction from f = μ × N.

```
WRONG: f_k = μ_k × mg   (on a ramp)
RIGHT: N = mg cos θ, then f_k = μ_k × N = μ_k × mg cos θ
```

### Mistake 2: Assuming static friction always equals f_s_max

Static friction only equals f_s_max when the object is on the verge of sliding. If the applied force is less, static friction adjusts to match:

```
If F_applied = 10 N and f_s_max = 30 N:
  f_s = 10 N (not 30 N)
  Object stays still.
```

### Mistake 3: Thinking μ_k = μ_s

Always μ_k < μ_s. They are different values. Using μ_s when you should use μ_k (or vice versa) gives wrong answers.

### Mistake 4: Thinking bigger area = more friction

Kinetic friction is independent of area. Only N matters. A heavy wide brick and a heavy narrow brick have the same friction if they have the same weight and same material.

### Mistake 5: Thinking friction always opposes the applied force

Friction opposes **motion** (or attempted motion), not the applied force direction. For a driving wheel, friction from the road acts FORWARD (same direction as travel), which is opposite to the direction the tyre surface would slip (backward).

### Mistake 6: Ignoring drag because "air seems light"

Air at 1.2 kg/m³ might seem negligible, but at high speeds (v²!) drag becomes enormous. A car at 30 m/s (108 km/h) experiences far more drag than at 10 m/s — nine times more, in fact.

---

## 21. Practice Problems

**Problem 1 (Easy)**
A 15 kg box sits on a floor. μ_s = 0.5. What force is needed to just start it moving?

**Problem 2 (Easy)**
A 10 kg box slides across a floor at constant velocity under an applied force of 25 N. What is μ_k? (Hint: constant velocity means F_net = 0)

**Problem 3 (Medium)**
A 30 kg box is pushed with 120 N across a floor. μ_k = 0.25. Find: (a) the normal force, (b) the friction force, (c) the net force, (d) the acceleration.

**Problem 4 (Medium)**
A 6 kg block is on a 25° ramp. μ_s = 0.35. Does it slide?
(sin 25° = 0.42, cos 25° = 0.91)

**Problem 5 (Medium)**
A 90 kg skydiver falls at terminal velocity with A = 0.8 m² and C_d = 1.1. What is their terminal velocity? (ρ_air = 1.2 kg/m³)

**Problem 6 (Challenging)**
A 5 kg block is placed on a ramp (μ_k = 0.2). The ramp is tilted at 40°. Find the acceleration of the block as it slides down.
(sin 40° = 0.643, cos 40° = 0.766)

---

**Answers:**

1. f_s_max = μ_s × N = μ_s × mg = 0.5 × 15 × 9.8 = **73.5 N**

2. Constant velocity → F_net = 0 → F_applied = f_k → f_k = 25 N
   N = mg = 10 × 9.8 = 98 N
   μ_k = f_k / N = 25 / 98 = **0.255**

3. (a) N = mg = 30 × 9.8 = **294 N**
   (b) f_k = μ_k × N = 0.25 × 294 = **73.5 N**
   (c) F_net = 120 - 73.5 = **46.5 N**
   (d) a = F_net / m = 46.5 / 30 = **1.55 m/s²**

4. mg sin 25° = 6 × 9.8 × 0.42 = 24.7 N
   N = mg cos 25° = 6 × 9.8 × 0.91 = 53.5 N
   f_s_max = μ_s × N = 0.35 × 53.5 = 18.7 N
   24.7 N > 18.7 N → **YES, it slides**

5. ½ × 1.2 × 1.1 × 0.8 × v_t² = 90 × 9.8
   0.528 × v_t² = 882
   v_t² = 1670.5
   v_t = **40.9 m/s** ≈ 41 m/s

6. N = mg cos 40° = 5 × 9.8 × 0.766 = 37.5 N
   f_k = μ_k × N = 0.2 × 37.5 = 7.5 N
   F_net = mg sin 40° - f_k = 5 × 9.8 × 0.643 - 7.5 = 31.5 - 7.5 = 24.0 N
   a = F_net / m = 24.0 / 5 = **4.8 m/s²** (down the slope)

---

## 22. Summary

- **Friction** is the force opposing relative motion (or tendency of motion) between surfaces in contact. It arises from microscopic interlocking of surface asperities and molecular adhesion.

- **Normal force** (N) is the perpendicular contact force between surfaces. Friction depends on N, not directly on mass or weight. On a flat surface: N = mg. On a ramp: N = mg cos θ.

- **Static friction** (f_s) opposes the start of motion. Its magnitude adjusts from 0 up to a maximum: f_s_max = μ_s × N. Until F_applied exceeds this maximum, the object does not move.

- **Kinetic friction** (f_k) opposes sliding motion. It has a fixed value: f_k = μ_k × N. It is always less than maximum static friction: μ_k < μ_s.

- **The coefficients** μ_s and μ_k are dimensionless numbers that depend on the materials and surface condition. Common values range from ≈ 0.01 (ice on ice) to ≈ 0.8 (rubber on dry concrete).

- Kinetic friction is: (1) proportional to N, (2) independent of contact area, (3) approximately independent of sliding speed.

- On an inclined plane at angle θ: N = mg cos θ, the force down the slope is mg sin θ, and friction = μ × mg cos θ. The critical angle where sliding begins: tan θ_c = μ_s.

- **Drag** (air resistance) is the friction force from a fluid (air, water) on a moving object. F_drag = ½ × ρ × C_d × A × v². Drag increases with the SQUARE of speed.

- **Terminal velocity** occurs when drag equals gravity: the object no longer accelerates and falls at constant speed. v_t = √(2mg / ρ C_d A).

- **Parachutes** reduce terminal velocity by dramatically increasing A, thus increasing drag to match weight at a much lower speed (~5 m/s vs ~55 m/s without parachute).

- **Useful friction** includes: walking, driving, braking, writing, gripping. **Harmful friction** includes: machine wear, energy loss, heat generation.

- **Friction reduction methods**: lubrication, ball/roller bearings, streamlined shapes, surface polishing, low-μ materials (Teflon).

---

## 23. Key Equations

```
Maximum static friction:
   f_s_max = μ_s × N
   (actual f_s can be any value from 0 to f_s_max)

Kinetic friction:
   f_k = μ_k × N
   (always less than f_s_max; μ_k < μ_s)

Relation between coefficients:
   μ_k < μ_s   (always — harder to start than to continue)

Normal force on a flat surface:
   N = mg

Normal force on a ramp (angle θ):
   N = mg cos θ

Force components on a ramp:
   Along slope (down):        mg sin θ
   Perpendicular to slope:    mg cos θ = N

Friction on a ramp:
   f_k = μ_k × mg cos θ
   f_s_max = μ_s × mg cos θ

Critical angle for sliding:
   tan θ_c = μ_s
   θ_c = arctan(μ_s)

Drag force:
   F_drag = ½ × ρ × C_d × A × v²
   where ρ = fluid density, C_d = drag coefficient,
         A = cross-sectional area, v = speed

Terminal velocity (when F_drag = mg):
   ½ × ρ × C_d × A × v_t² = mg
   v_t = √(2mg / ρ × C_d × A)

Net force on ramp (sliding down):
   F_net = mg sin θ - μ_k × mg cos θ
   F_net = mg(sin θ - μ_k cos θ)
   a = g(sin θ - μ_k cos θ)
```

---

*Next chapter: Chapter 16 — Work, Energy, and Power*
