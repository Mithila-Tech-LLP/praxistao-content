# Chapter 10: Operating Systems — Windows, macOS, Linux, Android, iOS

> **"The operating system is the manager of a very busy office. It doesn't do the actual work itself — the apps do that. But without the manager, the apps would all fight over the printer, steal each other's resources, and the whole office would collapse."**

---

## Table of Contents

1. [What an Operating System Does](#1-what-an-operating-system-does)
2. [Windows](#2-windows)
3. [macOS](#3-macos)
4. [Linux](#4-linux)
5. [Android](#5-android)
6. [iOS and iPadOS](#6-ios-and-ipados)
7. [Why Are There So Many? Which Is Best?](#7-why-are-there-so-many-which-is-best)
8. [The Key Jobs of Any OS](#8-the-key-jobs-of-any-os)
9. [Summary](#summary)

---

## 1. What an Operating System Does

Without an operating system, every app would have to:
- Figure out how to talk to your specific screen
- Manage its own RAM without bumping into other apps
- Write its own code to talk to the keyboard
- Implement its own file system

That would be impossible to maintain. Instead, the OS handles all of that, and apps just use the OS's services.

```
Without OS:                    With OS:
  
  App → talks directly       App → calls OS
        to hardware          OS → talks to hardware
        
  Problem: every app        Benefit: apps are simple;
  must know every           hardware complexity
  piece of hardware         is hidden from apps
```

The OS provides:
- A **file system** (folders and files)
- A **window manager** (the graphical interface)
- **Process management** (running multiple apps at once)
- **Memory management** (RAM for each app)
- **Device drivers** (talking to keyboard, screen, USB)
- **Networking** (connecting to the internet)
- **Security** (permissions, login, virus protection)

---

## 2. Windows

```
Made by:   Microsoft
First version: 1985 (Windows 1.0)
Current:   Windows 11 (2021)
Users:     ~1.6 billion active devices (72% of desktops/laptops)
Used on:   Desktops, laptops (most PCs)
```

**History:**
- 1985: Windows 1.0 — a graphical layer on top of MS-DOS
- 1995: Windows 95 — the first "real" Windows (Start button, taskbar)
- 2001: Windows XP — beloved, ran for 13 years before people moved on
- 2007: Windows Vista — widely disliked for being slow and annoying
- 2009: Windows 7 — excellent, very popular
- 2012: Windows 8 — touch-focused, removed Start button, users hated it
- 2015: Windows 10 — Start button back, free upgrade, very successful
- 2021: Windows 11 — new design, requires newer hardware

**Strengths:**
- Runs on the most hardware (almost any PC)
- Most software is made for Windows first
- Best for gaming (most games are Windows-only)
- Required in many businesses

**Weaknesses:**
- Historically more vulnerable to viruses (biggest target)
- Updates can be disruptive
- Less polish than macOS on the same hardware

---

## 3. macOS

```
Made by:   Apple
First version: Mac System 1 (1984), modern macOS from 2001
Current:   macOS Sequoia 15 (2024)
Users:     ~100 million active Macs (about 16% of desktops/laptops)
Used on:   MacBook, iMac, Mac mini, Mac Pro, Mac Studio
```

**History:**
- 1984: Original Mac — first popular computer with a GUI
- 2001: Mac OS X — rebuilt on Unix (FreeBSD) foundation
- 2006: Intel Macs — switched from PowerPC to Intel chips
- 2020: Apple Silicon — switched from Intel to their own M-series chips
- macOS versions named after California landmarks: Ventura, Sonoma, Sequoia

**Strengths:**
- Tightly integrated with Apple hardware (optimized)
- Beautiful, consistent design
- Excellent for creative work (design, video, music production)
- Unix underpinnings make it great for programming
- Very good security
- Seamless with iPhone and iPad

**Weaknesses:**
- Expensive (only runs on Apple hardware)
- Less gaming support than Windows
- Less software availability in some niche areas
- Harder to upgrade hardware (especially on newer models)

---

## 4. Linux

```
Made by:   Linus Torvalds (1991) + thousands of contributors
Current:   Linux 6.x (updates constantly)
Users:     ~3–4% of desktops, but ~97% of servers and cloud
Used on:   Servers, supercomputers, Android, smart TVs, routers, IoT
```

**What makes Linux different:**
Linux is the kernel (core of the OS). People take it and build complete operating systems (called "distributions" or "distros") around it.

```
Popular Linux distributions:
  Ubuntu    → easiest for beginners, very popular
  Fedora    → cutting-edge, used by developers
  Debian    → stable, the base of Ubuntu
  Arch      → minimalist, for advanced users who like control
  Mint      → Windows-like feel, good for switchers
  Kali      → security testing tools built-in
  Android   → mobile Linux (more on this below)
```

**Strengths:**
- Free (no license cost)
- Very stable and secure
- Enormous customization
- Preferred by developers and sysadmins
- Runs on almost anything (old PCs, tiny devices, supercomputers)

**Weaknesses:**
- Steeper learning curve
- Less "just works" polish
- Some software not available (Adobe Creative Suite, Microsoft Office native)
- Gaming improving but not as good as Windows

---

## 5. Android

```
Made by:   Google (acquired Android Inc. in 2005)
Current:   Android 15 (2024)
Users:     ~3 billion active devices (73% of smartphones globally)
Used on:   Smartphones, tablets, TVs, Wear OS watches, car displays
```

**Android is Linux:**
Android's core is built on Linux. The Linux kernel manages hardware. On top of it, Google added:
- Android Runtime (ART) — runs Java/Kotlin apps
- Android API — what apps use to access the phone's features
- Google Play Services — payment, notifications, Maps integration
- Material Design — the visual language

**How apps work:**
Apps are written in Java or Kotlin, compiled to bytecode, and run inside a virtual machine. This is why the same app can run on phones from different manufacturers.

**Open source but not entirely:**
The Android core (AOSP: Android Open Source Project) is open source. But Google's apps (Maps, Gmail, Play Store) are closed source. Phones without Google services (Huawei, some China-specific phones) use AOSP but without Google apps.

---

## 6. iOS and iPadOS

```
Made by:   Apple
First:     2007 (with the original iPhone)
Current:   iOS 18 (2024)
Users:     ~1.6 billion iPhones + iPads
Used on:   iPhone, iPod Touch (discontinued), iPad, iPadOS
```

**Strengths:**
- Most secure mobile OS
- Seamless hardware + software integration
- Best app quality (developers often release iOS first)
- 5–7 years of software updates (Android phones often get 2–3)
- Privacy features (App Tracking Transparency, Privacy Labels)

**Weaknesses:**
- Walled garden: Apple controls everything (only allowed on Apple devices)
- Less customization than Android
- Sideloading apps (installing outside App Store) difficult (recently opened in EU)
- No physical back button (feature or bug, depending on preference)

---

## 7. Why Are There So Many? Which Is Best?

```
For a desktop/laptop:
  Gaming?              → Windows
  Creative work?       → macOS
  Programming/Server?  → Linux or macOS
  Budget?              → Linux (free) or Windows
  Privacy?             → Linux
  Just "it works"?     → Windows or macOS
  
For a phone:
  Privacy + ecosystem?  → iPhone (iOS)
  Customization?        → Android
  Budget?               → Android (huge range of prices)
  Best camera?          → Either (Google Pixel, iPhone both excellent)
```

There is no universally "best" OS. They have different strengths for different uses.

---

## 8. The Key Jobs of Any OS

No matter which OS you use, they all do these same fundamental things:

```
1. Boot:
   Load itself when the computer starts
   Initialize all hardware
   
2. Run programs:
   Load apps from storage into RAM
   Give each app CPU time
   Keep apps from crashing each other
   
3. Manage memory:
   Allocate RAM to each running app
   Reclaim RAM when app closes
   
4. Manage files:
   Organize data in folders/files on storage
   Control who can read/write which files
   
5. Handle input:
   Receive keyboard, mouse, touch events
   Route them to the correct app
   
6. Display output:
   Tell the screen what to show
   Composite multiple app windows together
   
7. Network:
   Connect to Wi-Fi, cellular, ethernet
   Let apps send/receive internet data
   
8. Security:
   Password/PIN/biometric login
   Permission system (ask before accessing camera/mic)
   App isolation (App A can't read App B's data)
```

---

## Summary

| OS | Maker | Best For | Market |
|----|-------|---------|--------|
| Windows | Microsoft | Gaming, business, general PC | 72% of PCs |
| macOS | Apple | Creative work, developers | 16% of PCs |
| Linux | Community | Servers, developers, privacy | 97% of servers |
| Android | Google | Smartphones (affordable to premium) | 73% of phones |
| iOS | Apple | iPhones (premium, security) | 27% of phones |

**The operating system is the foundation everything else sits on. Next: how does your OS organize your files?**
