# Chapter 04: Inside a Computer — Meet the Hardware

> **"Opening a computer is like opening a watch — suddenly the mystery becomes a collection of understandable parts. Each part has one job. Together they make the whole thing work."**

---

## Table of Contents

1. [The Motherboard — The Foundation](#1-the-motherboard--the-foundation)
2. [The CPU — The Brain](#2-the-cpu--the-brain)
3. [RAM — The Desk](#3-ram--the-desk)
4. [Storage — The Filing Cabinet](#4-storage--the-filing-cabinet)
5. [The Power Supply — The Electricity Source](#5-the-power-supply--the-electricity-source)
6. [The GPU — The Graphics Card](#6-the-gpu--the-graphics-card)
7. [Cooling — Keeping Things Cold](#7-cooling--keeping-things-cold)
8. [How All the Parts Connect](#8-how-all-the-parts-connect)
9. [Summary](#summary)

---

## 1. The Motherboard — The Foundation

The **motherboard** is a large circuit board that all other components plug into. It's the "city grid" of the computer — every part communicates through it.

```
What the motherboard contains:
  
  ┌────────────────────────────────────────────────────────┐
  │                    MOTHERBOARD                          │
  │                                                        │
  │  ┌──────┐   ┌──────────┐   ┌──────────────────────┐   │
  │  │ CPU  │   │  RAM     │   │  PCIe slots           │   │
  │  │ slot │   │  slots   │   │  (GPU, extra cards)   │   │
  │  └──────┘   └──────────┘   └──────────────────────┘   │
  │                                                        │
  │  ┌──────────┐   ┌────────────┐   ┌────────────────┐   │
  │  │ Storage  │   │ BIOS chip  │   │ USB / HDMI /   │   │
  │  │ ports    │   │ (firmware) │   │ Audio ports    │   │
  │  └──────────┘   └────────────┘   └────────────────┘   │
  │                                                        │
  │  Tiny copper traces (like roads) connect everything    │
  └────────────────────────────────────────────────────────┘
```

Think of the motherboard as the city, and each component as a building in that city. The copper traces (circuits) are the roads connecting them.

---

## 2. The CPU — The Brain

The **CPU** (Central Processing Unit) is the most important part of a computer. It's a small chip (about the size of a postage stamp) that does ALL the thinking.

```
What the CPU looks like:
  
  ┌─────────────────────┐
  │  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │  ← Top surface (metal heat spreader)
  │  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │
  │  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │  About 40mm × 40mm
  │  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  │  
  └─────────────────────┘
  
  Underneath: hundreds of tiny gold pins that plug into motherboard
```

**What the CPU does:**
- Runs instructions — billions of them per second
- Does arithmetic: adds, subtracts, multiplies, divides
- Makes decisions: "if this, then that"
- Manages all other components

**CPU speed** is measured in GHz (gigahertz):
- 1 GHz = 1 billion cycles per second
- A modern CPU runs at 3–5 GHz
- A cycle is one small operation (like adding two numbers)

**CPU cores:**
Modern CPUs have multiple "cores" — each core is essentially a complete CPU on its own.
- 2 cores: can do 2 things at once
- 8 cores: can do 8 things at once
- Apple M4 has 10 cores; high-end server CPUs have 128+ cores

---

## 3. RAM — The Desk

**RAM** (Random Access Memory) is your computer's working space — the desk where it puts everything it's currently using.

```
Analogy: Imagine you're writing a report.
  
  Your desk  →  RAM
  ┌─────────────────────────────────────────────────┐
  │  Browser    Word Doc    Music player    Email   │
  │  (open)     (open)      (playing)       (open)  │
  └─────────────────────────────────────────────────┘
  
  Your filing cabinet  →  Hard Drive (Storage)
  ┌─────────────────────────────────────────────────┐
  │  Old documents, photos, installed apps, files   │
  │  (stored away, but reachable)                   │
  └─────────────────────────────────────────────────┘
```

**Key facts about RAM:**
- Very fast to access (1,000× faster than a hard drive)
- **Temporary** — everything in RAM is wiped when you turn off the computer
- More RAM = more apps open at once without slowing down

**RAM sizes:**
- 4GB — minimal (phone-level)
- 8GB — fine for everyday use
- 16GB — comfortable for most tasks
- 32GB+ — video editing, programming, heavy use

**Physical appearance:**
RAM comes as "sticks" — long thin green circuit boards that slot into the motherboard.

---

## 4. Storage — The Filing Cabinet

Storage is where your data lives **permanently** — even when the power is off.

### Hard Disk Drive (HDD)

```
┌──────────────────────────────────────┐
│  Inside an HDD:                      │
│                                      │
│    ┌────┐     Spinning magnetic disk │
│    │ ●● │  ← (like a CD, but stores data)│
│    │ ●● │     Spins at 5,400–7,200 RPM│
│    └────┘                            │
│      ↕                               │
│    Read/write head (like a record    │
│    player needle) moves over disk    │
└──────────────────────────────────────┘
```

- Stores data magnetically
- Cheaper per gigabyte
- Slower (moving parts)
- Can fail if dropped (physical spinning parts)

### Solid State Drive (SSD)

```
┌──────────────────────────────────────┐
│  Inside an SSD:                      │
│                                      │
│    ┌──────────────────────────┐      │
│    │  Flash memory chips      │      │
│    │  (like a USB stick,      │      │
│    │   but much faster)       │      │
│    └──────────────────────────┘      │
│                                      │
│  No moving parts — stores data       │
│  electrically in memory cells        │
└──────────────────────────────────────┘
```

- 3–10× faster than HDD
- Silent (no moving parts)
- Survives drops
- More expensive per gigabyte
- All modern laptops and phones use flash storage

**Storage sizes:**
- 128GB — small (a few thousand photos)
- 256GB — typical phone
- 1TB (1,000GB) — typical laptop
- 4TB+ — desktop with lots of files/games

---

## 5. The Power Supply — The Electricity Source

The **PSU** (Power Supply Unit) converts the electricity from your wall outlet (AC power) into the type computers need (DC power, at specific voltages).

```
Wall outlet: 120V AC (alternating current)
    ↓
PSU converts it to:
  3.3V DC → logic circuits
  5V DC   → USB ports, some chips
  12V DC  → CPU, GPU, fans, drives
```

Laptops have a small adapter brick that does the same job. Phones charge via USB which provides 5V (or up to 20V for fast charging with newer standards).

---

## 6. The GPU — The Graphics Card

The **GPU** (Graphics Processing Unit) is a specialized processor designed to handle graphics and visual output.

```
Why not just use the CPU for graphics?
  
  Rendering a single video frame requires:
  2,073,600 pixels on a 1080p screen
  × 3 color calculations each (Red, Green, Blue)
  = 6,220,800 calculations per frame
  × 60 frames per second
  = 373,248,000 calculations per second — just for video
  
  The CPU is good at doing complex things one at a time.
  The GPU is good at doing millions of simple things simultaneously.
  
  CPU:  8 very powerful cores
  GPU:  10,000 smaller cores (each weaker, but many in parallel)
```

GPUs are also used for:
- Machine learning / AI training
- Scientific simulations
- Bitcoin/cryptocurrency mining (which is just doing math really fast)

**Types:**
- **Integrated GPU** — built into the CPU chip (most laptops)
- **Discrete GPU** — separate card (gaming computers, workstations)

---

## 7. Cooling — Keeping Things Cold

All these chips generate heat. Left uncooled, a CPU would reach 100°C+ and destroy itself within seconds.

```
Methods of cooling:
  
  Air cooling (most common):
    Heat sink: metal fins that spread heat over a large surface area
    Fan: blows air over the fins to carry heat away
    
  Liquid cooling (high-performance):
    Liquid circulates through a metal block on the CPU
    Liquid carries heat to a "radiator" at the back of the case
    Fan blows through the radiator
    
  Thermal paste:
    A gray paste between CPU and cooler
    Fills microscopic air gaps that would trap heat
    Without it, the CPU would overheat immediately
```

Your laptop's fan spinning up loudly = the CPU is working hard and generating heat that needs to be removed.

---

## 8. How All the Parts Connect

```
Information flows between components via "buses" — high-speed data paths.

CPU ←──── PCIe bus (fastest) ─────→ GPU
 ↕
Memory bus
 ↕
RAM (fast, temporary)

 ↕
SATA / NVMe bus (fast-ish)
 ↕
Storage (slower, permanent)

 ↕
USB / I2C / etc.
 ↕
Keyboard, Mouse, Screen, Webcam

Everything is connected. The CPU is in the middle — it orchestrates all of it.
```

**The key insight:** A computer is just components passing information to each other at high speed, coordinated by the CPU, managed by software.

---

## Summary

| Component | Job | Analogy |
|-----------|-----|---------|
| Motherboard | Connects all parts | The city (roads + buildings) |
| CPU | Does all processing | The brain |
| RAM | Working memory (temporary) | Your desk |
| Storage (HDD/SSD) | Permanent memory | Your filing cabinet |
| GPU | Handles graphics | A specialized artist |
| PSU | Provides electricity | The food (power) |
| Cooling | Removes heat | The air conditioning |

**You now know what's inside every computer. Next: let's look more closely at the CPU — how does it actually "think"?**
