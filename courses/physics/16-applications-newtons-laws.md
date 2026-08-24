# Chapter 16: Applications of Newton's Laws

> "Give me a lever long enough and a place to stand, and I will move the world." — Archimedes

---

## Table of Contents

1. [The Four-Step Method for Any Dynamics Problem](#1-the-four-step-method-for-any-dynamics-problem)
2. [Free Body Diagrams — Your Most Important Tool](#2-free-body-diagrams--your-most-important-tool)
3. [Connected Objects — Two Boxes on a Rope](#3-connected-objects--two-boxes-on-a-rope)
4. [The Atwood Machine — Two Masses Over a Pulley](#4-the-atwood-machine--two-masses-over-a-pulley)
5. [Block on a Frictionless Inclined Plane](#5-block-on-a-frictionless-inclined-plane)
6. [Block on an Inclined Plane WITH Friction](#6-block-on-an-inclined-plane-with-friction)
7. [Elevator Problems — When Do You Feel Heavy or Light?](#7-elevator-problems--when-do-you-feel-heavy-or-light)
8. [Box on a Table Connected to a Hanging Box](#8-box-on-a-table-connected-to-a-hanging-box)
9. [Common Mistakes to Avoid](#9-common-mistakes-to-avoid)
10. [Summary](#summary)
11. [Key Equations](#key-equations)

---

## 1. The Four-Step Method for Any Dynamics Problem

Newton's Second Law — F_net = m × a — is a simple equation. But real problems involve multiple forces, multiple objects, and multiple directions. Without a systematic method you will make mistakes almost every time.

Here is the four-step method that professional physicists and engineers use. Follow it every single time, even for problems that seem easy. The habit will save you on hard problems.

---

### Step 1 — Draw a Free Body Diagram (FBD)

A **Free Body Diagram** is a sketch showing ONE object in isolation, with arrows representing every force acting ON that object. The arrows show the direction and relative size of each force.

Rules for a good FBD:
- Draw the object as a simple box or dot. You do not need artistic skill.
- Include ONLY forces acting ON the chosen object. Do not draw forces it exerts on other things.
- Label each force clearly: weight W (or mg), normal force N, tension T, friction f, applied force F, etc.
- The arrow starts at the point where the force is applied and points in the direction the force acts.

---

### Step 2 — Define a Positive Direction

Choose which direction you will call "positive." Common choices:
- For horizontal problems: rightward is positive.
- For vertical problems: upward is positive.
- For incline problems: direction of motion (down the slope) is positive.

Write your choice down. This prevents sign errors that ruin otherwise correct work.

---

### Step 3 — Write F = ma for Each Axis

Sum up all forces in the x-direction and set equal to m × a_x.
Sum up all forces in the y-direction and set equal to m × a_y.

If the object does not accelerate in one direction (for example, a box sitting on a flat table does not accelerate vertically), then a = 0 in that direction, and the forces must balance.

---

### Step 4 — Solve the Algebra

Now that you have equations, it is just mathematics. Solve for the unknown(s).

---

## 2. Free Body Diagrams — Your Most Important Tool

Let us practice drawing FBDs before tackling the worked examples.

### Example: A box sitting still on a flat table

```
           FBD of BOX
                |
           N (up, normal force from table)
                |
   +y           |
    ^        [  BOX  ]
    |            |
    |           mg (down, gravity)
    +---> +x     |
```

Forces acting on the box:
- Weight mg pointing straight down (Earth pulls the box down)
- Normal force N pointing straight up (table pushes the box up)

Since the box is not accelerating: N - mg = 0, so N = mg.

---

### Example: A box being pushed horizontally across a rough floor

```
           FBD of BOX

          N (up)
          |
F (push) -->[BOX]--> f (friction, opposes motion, pointing left)
          |
         mg (down)
```

Forces:
- Weight mg downward
- Normal force N upward
- Applied force F to the right
- Friction force f to the left (opposes motion)

Vertical: N - mg = 0 → N = mg  
Horizontal: F - f = ma → a = (F - f) / m

---

## 3. Connected Objects — Two Boxes on a Rope

### The Key Concept: Tension

When two objects are connected by a **rope** (or string), the rope pulls BOTH objects toward each other. The magnitude of this pulling force is called the **tension** T.

An important assumption: if the rope is light (we say "massless"), the tension is the SAME throughout the entire rope. The rope does not have extra mass to accelerate, so the pull is transmitted perfectly from one end to the other.

Think of it like a chain of people pushing a stalled car. If the chain is light, the push force passes through unchanged.

---

### Setting Up a Connected-Objects Problem

When you have two (or more) objects connected and moving together, they ALL share the same acceleration. But you draw a separate FBD for each object.

Then you write F = ma for each object separately, which gives you two equations. You combine those equations to find the two unknowns: acceleration a and tension T.

---

### Worked Example 3.1 — Two Boxes, Frictionless Surface

**Problem:** Box A has mass 3 kg. Box B has mass 5 kg. They are connected by a light rope. An applied force F = 16 N pulls Box A to the right on a frictionless surface. Find: (a) the acceleration of the system, and (b) the tension T in the rope.

```
     F = 16 N          T          T
   -------->[  A  ]--rope--[  B  ]
    3 kg                      5 kg
   
   Both boxes accelerate together to the RIGHT.
```

**Step 1 — Draw FBDs:**

```
FBD of Box A:               FBD of Box B:

     N_A (up)                    N_B (up)
      |                           |
F=16N -->[A]-->T (rope pulls B    T-->[B]
      |        forward)           |
     m_A × g (down)             m_B × g (down)
```

Wait — let us be careful about directions.

Box A is being pulled RIGHT by F = 16 N, and the rope connects to Box B on the RIGHT side of A. The rope therefore pulls Box A to the RIGHT? No — the rope pulls Box A BACKWARD (to the left), because Box B is in front and Box A is trying to accelerate faster than Box B; the rope restrains Box A.

Let us think physically. The applied force F is trying to accelerate the entire system. The rope carries the force that pulls Box B along. So:

- The rope pulls Box B FORWARD (to the right) with force T.
- The rope pulls Box A BACKWARD (to the left) with force T.

This is Newton's Third Law: the rope exerts equal and opposite forces on each box.

**Revised FBDs:**

```
FBD of Box A:               FBD of Box B:

      N_A                        N_B
       |                          |
F=16N->[A]<-T             T--->[B]
       |                          |
     m_A×g                     m_B×g

(rope pulls A left)      (rope pulls B right)
```

**Step 2 — Define positive direction:** Rightward is positive.

**Step 3 — Write F = ma:**

For Box A (horizontal):
  F - T = m_A × a
  16 - T = 3 × a     ... (equation 1)

For Box B (horizontal):
  T = m_B × a
  T = 5 × a          ... (equation 2)

**Step 4 — Solve:**

Substitute equation 2 into equation 1:
  16 - 5a = 3a
  16 = 8a
  a = 2 m/s²

Now find T using equation 2:
  T = 5 × 2 = 10 N

**Answers:**
- Acceleration a = 2 m/s²
- Tension T = 10 N

**Check:** Does this make sense? The total mass is 3 + 5 = 8 kg. The net external force is 16 N (no friction). So a = F / m_total = 16 / 8 = 2 m/s². Correct!

The tension 10 N is what Box B "feels." And 16 - 10 = 6 N is the net force on Box A, and 6 = 3 × 2. Also correct!

---

### Shortcut: Treat the System as One Object

For finding acceleration only, you can treat the entire connected system as a single object:

  F_net = m_total × a
  16 = (3 + 5) × a
  a = 16 / 8 = 2 m/s²

This works because internal forces (like tension) cancel out when you treat the system as a whole. To find tension, you still need to isolate one object.

---

## 4. The Atwood Machine — Two Masses Over a Pulley

### What Is an Atwood Machine?

An **Atwood machine** is one of the classic physics setups. It consists of two masses hanging on either side of a pulley connected by a rope. When you release them, the heavier mass falls and pulls the lighter mass up.

```
         [Pulley]
         /       \
        /         \
   T ↑ /           \ ↑ T
      /             \
   [ m₁ ]         [ m₂ ]
      ↓               ↑
   (m₁ falls)    (m₂ rises)
   
   Assume m₁ > m₂
```

The beauty of the Atwood machine is that it lets you study gravity and acceleration with control. Instead of an object falling freely at 9.8 m/s², the acceleration is reduced because both masses must be accelerated (one up, one down).

---

### Deriving the Formulas

Let m₁ be the heavier mass and m₂ be the lighter mass. Both experience tension T in the rope.

**Assumption:** the pulley is frictionless and massless (ideal pulley). This means the tension is the same on both sides.

Define positive direction: m₁ moves DOWN (positive), so m₂ moves UP (positive from its own perspective).

**FBDs:**

```
FBD of m₁ (heavier, going down):     FBD of m₂ (lighter, going up):

    T ↑                                   ↑ T
    |                                     |
   [m₁]                                 [m₂]
    |                                     |
  m₁g ↓ (gravity wins over T)          m₂g ↓ (T wins over gravity)
```

For m₁ (positive direction = down):
  m₁g - T = m₁ × a      ... (equation 1)

For m₂ (positive direction = up):
  T - m₂g = m₂ × a      ... (equation 2)

**Add the two equations** (this eliminates T):
  m₁g - T + T - m₂g = m₁a + m₂a
  m₁g - m₂g = (m₁ + m₂) × a
  (m₁ - m₂) × g = (m₁ + m₂) × a

**Solve for acceleration:**

  a = (m₁ - m₂) × g / (m₁ + m₂)

**Solve for tension:** Substitute back into equation 2:
  T = m₂(g + a)
  T = m₂ × g + m₂ × (m₁ - m₂)g / (m₁ + m₂)

After simplification (you can work through the algebra):
  T = 2 × m₁ × m₂ × g / (m₁ + m₂)

---

### Worked Example 4.1 — Atwood Machine

**Problem:** m₁ = 4 kg, m₂ = 6 kg. Which mass falls? Find the acceleration and the tension.

Since m₂ = 6 kg > m₁ = 4 kg, mass m₂ falls and m₁ rises.

Redefine: let m₁ = 6 kg (heavier) and m₂ = 4 kg (lighter) for consistency with our formula.

**Acceleration:**
  a = (m₁ - m₂) × g / (m₁ + m₂)
  a = (6 - 4) × 9.8 / (6 + 4)
  a = 2 × 9.8 / 10
  a = 19.6 / 10
  a = 1.96 m/s²

**Tension:**
  T = 2 × m₁ × m₂ × g / (m₁ + m₂)
  T = 2 × 6 × 4 × 9.8 / (6 + 4)
  T = 2 × 24 × 9.8 / 10
  T = 470.4 / 10
  T = 47.04 N

**Answers:**
- The 6 kg mass falls, the 4 kg mass rises.
- Acceleration a ≈ 1.96 m/s²
- Tension T ≈ 47.0 N

**Sanity check:**
- If both masses were equal, a = 0. No motion. Makes sense.
- If m₂ were zero, a = g = 9.8 m/s². That would be free fall. Makes sense.
- Tension T = 47 N is between 4 × 9.8 = 39.2 N (weight of lighter mass) and 6 × 9.8 = 58.8 N (weight of heavier mass). This makes sense: T must be greater than m₂g to pull m₂ up, and T must be less than m₁g to allow m₁ to fall.

---

## 5. Block on a Frictionless Inclined Plane

### Why Inclines Are Interesting

An **inclined plane** (a ramp) redirects gravity. Instead of an object falling straight down, the slope guides it at an angle. This means you must break forces into components — one along the slope and one perpendicular to the slope.

```
         /|
        / |
       /  |  height h
      /   |
     /θ   |
    /_____|
   
   θ = angle of the slope from the horizontal
```

### Breaking Gravity Into Components

The weight of the block is mg, pointing straight down. On a slope, it is convenient to use axes aligned WITH and PERPENDICULAR TO the slope.

```
        /
       / ← block sits here
      /
     /     
    /  θ    
   /______  

   
Weight mg splits into two components:

   Along the slope (down the slope):   mg sin θ
   Perpendicular to slope (into slope): mg cos θ

        ↑ N (normal force, perpendicular to slope)
        |
   [BLOCK]
        \  ← mg cos θ (into slope, balanced by N)
         \
          \--> mg sin θ (down the slope, causes acceleration)
```

**Perpendicular to slope:** The block does not sink into or fly off the slope, so acceleration in this direction = 0.
  N - mg cos θ = 0
  N = mg cos θ

**Along the slope (positive = down the slope):** Only mg sin θ acts (no friction).
  mg sin θ = m × a
  a = g sin θ

Notice that mass cancels! All objects accelerate at the same rate down a frictionless slope, regardless of how heavy they are.

---

### Worked Example 5.1 — Frictionless Ramp

**Problem:** A block slides down a frictionless ramp inclined at 30° to the horizontal. Find:
(a) The acceleration of the block.
(b) The velocity of the block after it has slid 2 m from rest.

**Part (a) — Acceleration:**
  a = g sin θ
  a = 9.8 × sin 30°
  a = 9.8 × 0.5
  a = 4.9 m/s²

**Part (b) — Velocity after 2 m from rest:**

Use the kinematics equation: v² = u² + 2as
  u = 0 (starts from rest)
  a = 4.9 m/s²
  s = 2 m

  v² = 0 + 2 × 4.9 × 2
  v² = 19.6
  v = √19.6
  v ≈ 4.43 m/s

**Answers:**
- Acceleration = 4.9 m/s²
- Speed after 2 m ≈ 4.43 m/s (down the slope)

---

### What Happens at Different Angles?

```
θ = 0°   (flat surface): a = g sin 0° = 0   (no motion, as expected)
θ = 30°:  a = g × 0.50 = 4.9 m/s²
θ = 45°:  a = g × 0.71 = 6.9 m/s²
θ = 60°:  a = g × 0.87 = 8.5 m/s²
θ = 90°  (vertical cliff): a = g × 1.0 = 9.8 m/s² (free fall, as expected)
```

As the slope gets steeper, acceleration increases from 0 to g. Beautiful.

---

## 6. Block on an Inclined Plane WITH Friction

### Adding Friction to the Mix

Now the slope is rough. When the block slides down, **kinetic friction** acts up the slope, opposing the motion.

```
        /   ← block slides DOWN
       /
      / [BLOCK]
     /   ↑ f (friction, up the slope)
    /    ↑↑ mg sin θ is DOWN, friction is UP
   /___________

Forces along the slope:
  Down the slope:   mg sin θ
  Up the slope:     f = μk × N = μk × mg cos θ

Net force down the slope = mg sin θ - μk × mg cos θ
```

**Along the slope (positive = down):**
  mg sin θ - μk × mg cos θ = m × a

Factor out mg:
  mg(sin θ - μk cos θ) = m × a

Mass cancels:
  a = g(sin θ - μk cos θ)

**This is an extremely important formula.** Learn it well.

---

### When Does the Block NOT Slide?

The block will only slide if gravity wins over static friction:
  mg sin θ > μs × mg cos θ
  tan θ > μs
  θ > arctan(μs)

If the slope angle is less than arctan(μs), the block stays put.

---

### Worked Example 6.1 — Ramp With Friction

**Problem:** A block slides down a rough inclined plane at 40° to the horizontal. The coefficient of kinetic friction between the block and surface is μk = 0.2. Find the acceleration.

**Step 1 — Identify forces:**
  Along slope (down): mg sin 40°
  Along slope (up): friction = μk × N = μk × mg cos 40°

**Step 2 — Write F = ma along slope:**
  mg sin 40° - μk × mg cos 40° = m × a
  a = g(sin 40° - μk cos 40°)

**Step 3 — Plug in numbers:**
  sin 40° ≈ 0.643
  cos 40° ≈ 0.766

  a = 9.8 × (0.643 - 0.2 × 0.766)
  a = 9.8 × (0.643 - 0.153)
  a = 9.8 × 0.490
  a = 4.80 m/s²

**Answer:** a ≈ 4.80 m/s² down the slope.

Compare to frictionless 40° ramp: a = g sin 40° = 9.8 × 0.643 = 6.3 m/s². Friction reduced the acceleration from 6.3 to 4.8 m/s². Makes sense.

---

### Full FBD for Block on Rough Incline

```
                   N (perpendicular to slope, pointing away from slope)
                   ↑
                   |
                   | (perpendicular axis)
    f (friction)   |
    <----------[BLOCK]----------> (along slope axis)
                   |
                   ↓
               mg sin θ (down the slope, along slope axis)

   Also: mg cos θ pushes block INTO slope (perpendicular axis),
         balanced by N.

   View rotated to align with slope:

       slope surface
   ________________________
        ↑ N
        |
   f<--[B]-->mg sin θ
        |
   (mg cos θ into slope, balanced by N)
```

---

## 7. Elevator Problems — When Do You Feel Heavy or Light?

### The Apparent Weight Concept

Have you ever noticed how you feel slightly heavier when an elevator starts moving up, and lighter when it slows down at the top? This is not your imagination — it is Newton's Second Law in action.

Your **true weight** is always mg, the gravitational force pulling you down. It never changes (unless you travel to another planet).

Your **apparent weight** is what a scale under your feet reads. It equals the Normal force N that the floor (or scale) pushes up on you. This CAN change when you accelerate.

---

### Four Cases for an Elevator

**Case 1 — Elevator at rest (or moving at constant speed):**

No acceleration, so a = 0.
  N - mg = 0
  N = mg

The scale reads your true weight. You feel normal.

---

**Case 2 — Elevator accelerating UPWARD:**

The elevator pushes you upward faster than gravity pulls you down. You are accelerating upward, so the net force must be upward.

```
      ↑ N (floor pushes you up)
      |
    [YOU]
      |
      ↓ mg (gravity pulls you down)
      
   Net upward: N - mg = m × a (a is positive, upward)
   Therefore: N = m(g + a)
```

N > mg. The scale reads MORE than your true weight. You feel HEAVIER.

---

**Case 3 — Elevator accelerating DOWNWARD:**

The elevator drops faster, so the net force is downward.

```
      ↑ N
      |
    [YOU]
      |
      ↓ mg
      
   Net force is downward: mg - N = m × a (a is positive, downward)
   Therefore: N = m(g - a)
```

N < mg. The scale reads LESS than your true weight. You feel LIGHTER.

---

**Case 4 — Elevator in free fall (a = g downward):**

The cable snaps and the elevator falls freely. Now a = g.
  N = m(g - g) = 0

The scale reads ZERO. You are **weightless** — you float inside the elevator!

This is exactly how astronauts experience weightlessness in the International Space Station. The ISS is in free fall around Earth (it just moves so fast horizontally that it keeps missing the ground). The astronauts and the station fall together, so they float.

---

### Worked Example 7.1 — Person in an Elevator

**Problem:** A person of mass 60 kg is standing on a scale inside an elevator. The elevator accelerates upward at 2 m/s². What does the scale read? Express in both Newtons and kilograms.

**Step 1 — FBD:**
```
   ↑ N (scale pushes person up)
   |
 [PERSON, 60 kg]
   |
   ↓ mg = 60 × 9.8 = 588 N
```

**Step 2 — F = ma (upward positive):**
  N - mg = m × a
  N = m(g + a)
  N = 60 × (9.8 + 2)
  N = 60 × 11.8
  N = 708 N

**Step 3 — Convert to apparent kg:**
  Scales often display in kg by dividing by g:
  Apparent mass = N / g = 708 / 9.8 ≈ 72.2 kg

**Answer:** The scale reads 708 N, or equivalently about 72.2 kg. The person appears to weigh about 12.2 kg more than their actual 60 kg.

---

### Worked Example 7.2 — Elevator Decelerating

**Problem:** Same 60 kg person. Elevator is moving upward but slowing down at 3 m/s².

When an elevator moves UP but slows down, its acceleration points DOWNWARD (deceleration = negative acceleration in upward direction).

  N = m(g - a)     [a = 3 m/s² downward]
  N = 60 × (9.8 - 3)
  N = 60 × 6.8
  N = 408 N

The person feels lighter — only about 41.6 kg equivalent — even though they are actually 60 kg.

---

### Summary of Elevator Cases

```
Situation                         N = ?            Feel?
---------------------------------------------------------------------
At rest / constant velocity       N = mg           Normal
Accelerating upward               N = m(g + a)     Heavier
Accelerating downward             N = m(g - a)     Lighter
Free fall (a = g downward)        N = 0            Weightless
```

---

## 8. Box on a Table Connected to a Hanging Box

### The Setup

This is a very common problem. One box sits on a table and is connected by a rope that passes over a pulley at the table edge. The other end of the rope has a hanging box.

```
   [Box A, on table]----rope----[Pulley at edge]
                                       |
                                       | (rope hangs down)
                                       |
                                   [Box B, hanging]
```

When released, Box B pulls Box A horizontally across the table. They are connected by the same rope, so they have the SAME acceleration (in magnitude).

---

### Worked Example 8.1 — Table and Hanging Box

**Problem:** Box A (mass 4 kg) sits on a frictionless table. It is connected by a light rope over a frictionless pulley to Box B (mass 2 kg) which hangs off the edge. Find the acceleration and tension.

**FBD of Box A (horizontal motion):**
```
      ↑ N_A
      |
T -->[A]     (T pulls it toward the right, toward pulley)
      |
    m_A × g = 39.2 N (down, balanced by N_A)
```

Horizontal for Box A (positive = right, toward pulley):
  T = m_A × a       ... (equation 1)
  T = 4a

**FBD of Box B (vertical motion):**
```
      ↑ T (rope pulls it up)
      |
    [B]
      |
      ↓ m_B × g = 2 × 9.8 = 19.6 N
```

Vertical for Box B (positive = down, direction of motion):
  m_B × g - T = m_B × a
  19.6 - T = 2a       ... (equation 2)

**Solve:**

Substitute equation 1 into equation 2:
  19.6 - 4a = 2a
  19.6 = 6a
  a = 19.6 / 6
  a ≈ 3.27 m/s²

Find T:
  T = 4 × 3.27 ≈ 13.1 N

**Answers:**
- Acceleration ≈ 3.27 m/s²
- Tension ≈ 13.1 N

**Check:** Is T < m_B × g? Yes, 13.1 < 19.6. Box B must have a net downward force to accelerate down, which requires T < m_B × g. Correct.

---

### What If the Table Has Friction?

If the surface has kinetic friction μk, then Box A also experiences a friction force opposing its motion (pointing LEFT, away from the pulley).

Friction on Box A:
  f = μk × N_A = μk × m_A × g

Equation for Box A becomes:
  T - f = m_A × a
  T - μk × m_A × g = m_A × a

Equation for Box B stays the same:
  m_B × g - T = m_B × a

Add both equations:
  m_B × g - μk × m_A × g = (m_A + m_B) × a
  a = (m_B - μk × m_A) × g / (m_A + m_B)

This makes sense: friction reduces the acceleration. If friction is large enough, the system might not move at all.

---

## 9. Common Mistakes to Avoid

### Mistake 1 — Including the wrong forces in the FBD

**Wrong:** Including forces that the box EXERTS on other objects.
**Correct:** Include ONLY forces acting ON the box you chose.

For example, when drawing the FBD of a box on a table:
- Include: the weight of the box (Earth pulling box down) and normal force (table pushing box up).
- Do NOT include: the weight of the table, the force the box pushes down on the table (that belongs in the table's FBD).

---

### Mistake 2 — Forgetting to decompose forces on inclines

**Wrong:** Applying F = ma directly with the weight mg along the slope without splitting.
**Correct:** The component of weight along the slope is mg sin θ. The full weight mg is NOT along the slope.

```
   WRONG: "The force along the slope = mg"   ← This would be true only on a vertical cliff.

   CORRECT:
      Weight vector mg points straight DOWN.
      Component along slope = mg sin θ
      Component into slope  = mg cos θ
```

---

### Mistake 3 — Sign errors

**Wrong:** Not defining a positive direction and mixing up signs.
**Correct:** Before writing any equation, write: "Positive direction = [direction]." Then check every force sign against this definition.

In an Atwood machine, if you call "m₁ downward" as positive, then the tension on m₁ acts UPWARD, so its contribution is NEGATIVE: m₁g - T = m₁a. If you accidentally write m₁g + T = m₁a, your answer will be completely wrong.

---

### Mistake 4 — Using the wrong friction force

**Wrong:** Using static friction when the object is moving, or applying friction in the wrong direction.
**Correct:**
- If the object is moving: use kinetic friction f = μk × N.
- If the object is stationary and you are checking whether it will move: compare the applied force to maximum static friction f_max = μs × N.
- Friction always opposes the direction of motion (or potential motion).

On an incline, friction acts UP the slope when the object slides DOWN the slope.

---

### Mistake 5 — Confusing acceleration of the system with acceleration of individual parts

In an Atwood machine, both masses have the same magnitude of acceleration (one up, one down). In connected horizontal boxes, both have the same acceleration. But their equations are separate — you must write F = ma for each.

---

### Mistake 6 — Forgetting that the normal force is NOT always equal to mg

On an incline: N = mg cos θ (not mg).
In an elevator: N = m(g ± a) (not mg).
On a flat surface with a vertical applied force (pushing down or pulling up): N changes accordingly.

---

## Summary

- The **four-step method** for dynamics: (1) Draw FBD, (2) Define positive direction, (3) Write F = ma for each axis, (4) Solve algebra.

- In **connected objects** on a frictionless surface: all objects share the same acceleration. Tension T is the same throughout a massless rope. Find a by treating the system as one object: a = F_total / m_total. Find T by isolating one object.

- In the **Atwood machine**: a = (m₁ - m₂)g / (m₁ + m₂) and T = 2m₁m₂g / (m₁ + m₂).

- On a **frictionless incline**: a = g sin θ. Normal force N = mg cos θ. Mass does not matter for acceleration.

- On a **rough incline**: a = g(sin θ - μk cos θ). Friction reduces acceleration.

- In an **elevator accelerating upward**: apparent weight N = m(g + a) — you feel heavier. Accelerating downward: N = m(g - a) — you feel lighter. In free fall: N = 0 — weightless.

- For a **table-and-hanging-box** system: treat each object separately, then add equations to eliminate tension and solve for a.

- Common mistakes: wrong forces in FBD, forgetting to decompose weight on inclines, sign errors, wrong friction direction.

---

## Key Equations

```
Newton's Second Law (general):
   F_net = m × a

Connected objects (frictionless):
   a = F_applied / (m_A + m_B)
   T = m_B × a   (for the pulled object)

Atwood machine:
   a = (m₁ - m₂) × g / (m₁ + m₂)
   T = 2 × m₁ × m₂ × g / (m₁ + m₂)

Frictionless incline:
   a = g × sin θ
   N = m × g × cos θ

Rough incline (sliding):
   a = g × (sin θ - μk × cos θ)
   N = m × g × cos θ
   f = μk × N

Elevator (upward acceleration):
   N = m × (g + a)

Elevator (downward acceleration):
   N = m × (g - a)

Friction force:
   f = μ × N    (μk for kinetic, μs for static)

Weight:
   W = m × g
```
