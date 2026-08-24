# Chapter 26: Rotational Kinematics

> **"Everything in rotational motion has an exact equivalent in linear motion. Once you see the pattern, you'll never find rotation confusing again."**

---

## Table of Contents

- [26.1 From Linear to Rotational](#261-from-linear-to-rotational)
- [26.2 Angle and Angular Displacement](#262-angle-and-angular-displacement)
- [26.3 Angular Velocity](#263-angular-velocity)
- [26.4 Angular Acceleration](#264-angular-acceleration)
- [26.5 The Rotational Equations of Motion](#265-the-rotational-equations-of-motion)
- [26.6 Relating Linear and Rotational Quantities](#266-relating-linear-and-rotational-quantities)
- [26.7 Torque](#267-torque)
- [26.8 Moment of Inertia](#268-moment-of-inertia)
- [26.9 Newton's Second Law for Rotation](#269-newtons-second-law-for-rotation)
- [26.10 Rotational Kinetic Energy](#2610-rotational-kinetic-energy)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 26.1 From Linear to Rotational

One of the beautiful things about physics is that rotation is completely parallel to linear motion. Every linear quantity has a rotational twin:

```
LINEAR              ROTATIONAL
----------          ----------
Position x          Angle θ (theta)
Velocity v          Angular velocity ω (omega)
Acceleration a      Angular acceleration α (alpha)
Mass m              Moment of inertia I
Force F             Torque τ (tau)
Linear momentum p   Angular momentum L
```

If you understand linear motion, you already understand rotation — just swap the symbols.

---

## 26.2 Angle and Angular Displacement

### Radians — The Natural Unit for Angles

In everyday life we use degrees. In physics, we use **radians (rad)**.

```
RADIAN DEFINITION:
  
  An angle of 1 radian is the angle formed when the arc
  length equals the radius.
  
         arc = r
      .--------.
    /            \
   |    angle     |  radius r
   |      1 rad   |
    \            /
      `--------'
  
  arc length = r × θ (in radians)
```

Converting:
- Full circle = 360° = 2π radians
- Half circle = 180° = π radians
- Quarter circle = 90° = π/2 radians
- 1 radian ≈ 57.3°

```
To convert degrees to radians: multiply by π/180
To convert radians to degrees: multiply by 180/π

Example:
  90° × (π/180) = π/2 ≈ 1.57 rad
  1 rad × (180/π) ≈ 57.3°
```

### Angular Displacement (θ)

**Angular displacement** is the angle swept by a rotating object.

```
θ_f - θ_i = Δθ (angular displacement, in radians)
```

Direction convention:
- Counter-clockwise (CCW) = positive
- Clockwise (CW) = negative

---

## 26.3 Angular Velocity

**Angular velocity (ω)** is the rate of change of angle:

```
ω = Δθ / Δt

units: rad/s
```

For uniform circular motion (constant speed), ω is constant.

### Relationship to Frequency and Period

If an object completes f rotations per second (frequency in Hz):

```
ω = 2π × f

Because one full rotation = 2π radians, so f rotations/second = 2πf rad/s.
```

Period T = 1/f (time for one full rotation), so:

```
ω = 2π / T
```

### Worked Example 26.1

A wheel spins at 1200 rpm (revolutions per minute). Find ω in rad/s.

**Solution:**
- f = 1200 rpm = 1200/60 = 20 rev/s = 20 Hz
- ω = 2π × f = 2π × 20 = 40π ≈ **125.7 rad/s**

---

## 26.4 Angular Acceleration

**Angular acceleration (α)** is the rate of change of angular velocity:

```
α = Δω / Δt

units: rad/s²
```

Positive α = object speeding up (if rotating CCW) or slowing down (if rotating CW).

---

## 26.5 The Rotational Equations of Motion

The four rotational equations of motion are identical in form to the linear equations (see Chapter 08), with the variable substitutions shown:

| Linear | Rotational |
|--------|-----------|
| x = displacement | θ = angular displacement |
| v = velocity | ω = angular velocity |
| a = acceleration | α = angular acceleration |
| v = u + at | ω = ω₀ + αt |
| x = ut + ½at² | θ = ω₀t + ½αt² |
| v² = u² + 2ax | ω² = ω₀² + 2αθ |
| x = ½(u+v)t | θ = ½(ω₀+ω)t |

The pattern is exact — just replace x→θ, v→ω, a→α, u→ω₀.

### Worked Example 26.2

A motor starts from rest and reaches ω = 100 rad/s in 5 seconds with constant angular acceleration.

(a) Find the angular acceleration.
(b) Find the number of rotations.

**Solution:**

(a) α = (ω - ω₀) / t = (100 - 0) / 5 = **20 rad/s²**

(b) θ = ω₀t + ½αt²
    θ = 0 + ½ × 20 × 25 = 250 rad

Number of rotations = θ / 2π = 250 / (2π) ≈ **39.8 rotations**

### Worked Example 26.3

A spinning top has ω₀ = 80 rad/s and decelerates at α = -4 rad/s².

How long until it stops? How many radians does it turn?

**Solution:**

Time to stop: ω = ω₀ + αt → 0 = 80 + (-4)t → t = **20 s**

Angle turned: ω² = ω₀² + 2αθ
0 = 80² + 2(-4)θ
8θ = 6400
θ = **800 rad** ≈ 127 rotations

---

## 26.6 Relating Linear and Rotational Quantities

Every point on a rotating object also has linear (tangential) motion. These are linked to the angular quantities through the radius r.

```
ROTATING DISK:

             v = ωr
          P•---------->
         /
        /  r
       /
    center (axis)
    
Point P at radius r:
  - Arc length swept: s = r × θ
  - Tangential speed: v = r × ω
  - Tangential acceleration: a_t = r × α
  - Centripetal acceleration: a_c = v²/r = ω²r
```

### Summary of Relationships

```
s = r × θ          (arc length)
v = r × ω          (tangential speed)
a_t = r × α        (tangential acceleration)
a_c = ω² × r       (centripetal acceleration)
```

Note: **v is not the same as ω**. Two points at different radii on the same rotating disk have the SAME ω but DIFFERENT v.

### Worked Example 26.4

A merry-go-round has radius 2 m and rotates at 3 rad/s.

(a) Find the speed of a child sitting at the edge.
(b) Find the centripetal acceleration.

**Solution:**

(a) v = r × ω = 2 × 3 = **6 m/s**

(b) a_c = ω² × r = 9 × 2 = **18 m/s²** (directed toward center)

---

## 26.7 Torque

**Torque (τ)** is the rotational equivalent of force. It is what causes angular acceleration.

```
τ = r × F × sin(θ)

where:
  τ = torque (N·m)
  r = distance from pivot to point of force application (m)
  F = force applied (N)
  θ = angle between F and the lever arm r
```

For maximum torque, apply force perpendicular to the lever arm (θ = 90°, sin = 1):

```
τ = r × F    (maximum torque)
```

### Visual Understanding of Torque

```
DOOR HINGE (pivot) at left.
Force F applied at right edge.

          r = 1 m
  [HINGE]-----------x --> F

  Torque = r × F

If the same force is applied closer to the hinge:

          r = 0.2 m
  [HINGE]----x --> F

  Torque = 0.2F   (much smaller!)
```

This is why door handles are at the far edge from the hinge.

### Direction of Torque

- Force that causes counter-clockwise rotation = positive torque
- Force that causes clockwise rotation = negative torque

### Worked Example 26.5

A spanner (wrench) of length 0.3 m applies a force of 50 N perpendicular to the handle.

Torque = r × F = 0.3 × 50 = **15 N·m**

---

## 26.8 Moment of Inertia

**Moment of inertia (I)** is the rotational equivalent of mass. It measures how hard it is to change an object's rotational speed.

```
For a point mass: I = m × r²

For an extended object: I = sum of (m_i × r_i²)

units: kg·m²
```

Crucially, I depends on how mass is distributed relative to the rotation axis — mass far from the axis contributes much more than mass near it.

### Common Moments of Inertia

```
SOLID CYLINDER (axis through center):
  I = (1/2) × m × r²

HOLLOW CYLINDER (thin wall):
  I = m × r²

SOLID SPHERE:
  I = (2/5) × m × r²

THIN ROD (axis through center):
  I = (1/12) × m × L²

THIN ROD (axis through end):
  I = (1/3) × m × L²

POINT MASS at distance r:
  I = m × r²
```

```
COMPARISON:
  Two cylinders, same mass m, same radius r:
  
  Hollow cylinder: I = mr²
  Solid cylinder: I = (1/2)mr²
  
  Hollow cylinder is HARDER to spin up (more I) because
  all its mass is at maximum distance r from axis.
  
  Solid cylinder has mass spread from 0 to r, so average
  r² is less.
```

---

## 26.9 Newton's Second Law for Rotation

The rotational form of Newton's Second Law:

```
τ_net = I × α

(Compare: F_net = m × a)

where:
  τ_net = net torque (N·m)
  I = moment of inertia (kg·m²)
  α = angular acceleration (rad/s²)
```

### Worked Example 26.6

A solid cylinder (m = 2 kg, r = 0.5 m) is free to rotate about its axis.

A tangential force of 10 N is applied at the rim. Find the angular acceleration.

**Solution:**

Moment of inertia: I = (1/2)mr² = (1/2)(2)(0.5²) = 0.25 kg·m²

Torque: τ = r × F = 0.5 × 10 = 5 N·m

Angular acceleration: α = τ/I = 5/0.25 = **20 rad/s²**

---

## 26.10 Rotational Kinetic Energy

A rotating object has kinetic energy due to its rotation:

```
KE_rot = (1/2) × I × ω²

(Compare: KE_linear = (1/2) × m × v²)
```

### Rolling Objects

An object rolling without slipping has BOTH translational and rotational kinetic energy:

```
KE_total = KE_trans + KE_rot
         = (1/2)mv² + (1/2)Iω²

Since v = rω for rolling without slipping:

KE_total = (1/2)mv² + (1/2)I(v/r)²
         = (1/2)v²(m + I/r²)
```

### Why Hollow Cylinders Roll Slower Down a Ramp

Both start with the same gravitational potential energy. But:

- Solid cylinder: I = (1/2)mr², so more energy goes to translational KE → faster
- Hollow cylinder: I = mr², so more energy goes to rotational KE → slower

```
RACE DOWN A RAMP:

              solid cylinder -----> (faster)
              hollow cylinder ---->  (slower)
              
Both have same mass and radius, but different I.
Hollow has more rotational KE, less translational KE.
```

---

## Summary

- **Rotational quantities**: θ (angle), ω (angular velocity), α (angular acceleration) — exact parallels to linear x, v, a
- **Angle in radians**: arc length s = rθ; full circle = 2π rad
- **Rotational equations of motion**: same form as linear, with variable substitutions
- **Linking linear and rotational**: v = rω, a_t = rα, a_c = ω²r
- **Torque**: τ = r × F × sin(θ) — rotational force; maximum when F is perpendicular to lever arm
- **Moment of inertia**: I = Σmr² — depends on mass distribution; I = (1/2)mr² for solid cylinder
- **Newton's 2nd Law for rotation**: τ_net = Iα
- **Rotational KE**: (1/2)Iω²
- Rolling objects have both translational and rotational KE

---

## Key Equations

```
Angle (radians):
  2π rad = 360°
  s = r × θ

Angular velocity:
  ω = Δθ / Δt
  ω = 2πf = 2π/T

Angular acceleration:
  α = Δω / Δt

Rotational equations of motion:
  ω = ω₀ + αt
  θ = ω₀t + (1/2)αt²
  ω² = ω₀² + 2αθ

Linking linear and rotational:
  v = rω
  a_t = rα
  a_c = ω²r

Torque:
  τ = r × F × sin(θ)
  τ_net = I × α

Moment of inertia:
  I = mr²  (point mass)
  I = (1/2)mr²  (solid cylinder)

Rotational KE:
  KE_rot = (1/2) × I × ω²
```
