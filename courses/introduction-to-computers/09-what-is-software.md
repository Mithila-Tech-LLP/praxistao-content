# Chapter 09: What Is Software?

> **"Hardware is the stage. Software is the play. A computer without software is just expensive metal. A song without instruments is just a melody in someone's head. You need both."**

---

## Table of Contents

1. [The Definition](#1-the-definition)
2. [Software is Instructions](#2-software-is-instructions)
3. [Types of Software](#3-types-of-software)
4. [How Software Gets on Your Computer](#4-how-software-gets-on-your-computer)
5. [Open Source vs. Closed Source](#5-open-source-vs-closed-source)
6. [How Software is Made (Overview)](#6-how-software-is-made-overview)
7. [Summary](#summary)

---

## 1. The Definition

**Software** is any set of instructions that tells a computer what to do.

Software is not a physical thing. You can't pick it up. You can't feel it. It's just data — patterns of 0s and 1s stored on your drive — but when the CPU reads those patterns, it does things: displays windows, plays music, sends messages.

```
Analogy:
  A recipe is software.
  A kitchen is hardware.
  
  The recipe is just text on paper — not physical in the way a pot is.
  But when a cook follows the recipe, amazing food is produced.
  
  Change the recipe → different food (same kitchen, same equipment)
  Change the software → different app (same CPU, same RAM)
  
  This is the power of software: one piece of hardware can do an 
  unlimited number of different things by loading different software.
```

---

## 2. Software is Instructions

Every piece of software is ultimately a long list of instructions for the CPU.

```
When you press "Play" on Spotify:
  
  The Spotify app contains instructions like:
  1. Check the song ID from the playlist
  2. Connect to Spotify's servers over the internet
  3. Request the audio data for that song
  4. Receive audio data in chunks
  5. Decode the compressed audio (MP3/OGG)
  6. Send decoded audio to the sound driver
  7. Update the progress bar on screen every 500ms
  8. Watch for user tapping pause/skip
  
  Millions of lines of code, millions of instructions,
  all happening in milliseconds.
```

A modern app like Google Chrome has about **35 million lines of code**. A Boeing 777 aircraft has about 7 million lines. The US tax system: ~10 million lines. All are just long lists of instructions telling machines what to do.

---

## 3. Types of Software

Software comes in several layers:

```
┌────────────────────────────────────────────────────────┐
│                 USER LAYER                             │
│                                                        │
│   Apps you use:  Browser, Word, Spotify, Games        │
│   (Application Software)                              │
├────────────────────────────────────────────────────────┤
│                 SYSTEM LAYER                           │
│                                                        │
│   Operating System: Windows, macOS, Linux, Android    │
│   (Manages hardware, runs apps)                        │
├────────────────────────────────────────────────────────┤
│                 FIRMWARE LAYER                         │
│                                                        │
│   BIOS/UEFI: First software to run when PC turns on   │
│   (Built into hardware, rarely changes)               │
├────────────────────────────────────────────────────────┤
│                 HARDWARE LAYER                         │
│                                                        │
│   CPU, RAM, Storage, Screen, Keyboard                 │
└────────────────────────────────────────────────────────┘
```

**Application software** — what you use directly:
- Productivity: Microsoft Word, Excel, Google Docs
- Communication: WhatsApp, Email, Zoom
- Entertainment: Spotify, Netflix, Games
- Browsers: Chrome, Safari, Firefox
- Utilities: antivirus, file manager, calculator

**System software** — runs the machine:
- Operating system (Chapter 10)
- Device drivers (let the OS talk to hardware)
- Programming tools (compilers, code editors)

**Firmware** — software burned into hardware:
- BIOS/UEFI on your motherboard (starts the boot process)
- Software in your keyboard, mouse, printer
- Software controlling your car's engine
- Rarely updated, very low-level

---

## 4. How Software Gets on Your Computer

```
Method 1: Pre-installed
  Comes with the device. Windows comes with Edge browser.
  macOS comes with Safari, Mail, Photos, Calendar.
  
Method 2: Downloaded from the internet
  Go to a website → download .exe (Windows) or .dmg (Mac)
  Run the installer → files are copied to your drive
  Example: installing VLC, Discord, Zoom
  
Method 3: App Store
  Apple App Store, Google Play Store, Microsoft Store
  Browse, click Install, app downloads automatically
  App Store verifies apps for safety before allowing them
  
Method 4: Package manager (for programmers)
  brew install vlc  (Mac)
  apt install vlc   (Linux)
  Downloads and installs from the internet, command-line style
  
Method 5: Web app (nothing to install)
  Gmail, Google Docs, Figma — these run in your browser
  No download, no install
  Software runs on the company's servers; you just see the result
```

---

## 5. Open Source vs. Closed Source

**Closed source (proprietary):**
- Source code (the human-readable instructions) is kept secret
- Only the compiled binary is released
- You can use it but can't see how it works or modify it
- Examples: Windows, macOS, Chrome, Microsoft Office, most commercial software

**Open source:**
- Source code is publicly available
- Anyone can read, modify, and contribute
- Usually free (though companies can sell services around it)
- Examples: Linux, Firefox, LibreOffice, VLC, Git, Python, Android

```
Why open source matters:
  Security: anyone can inspect for vulnerabilities — many eyes, fewer bugs
  Trust: you can verify the software doesn't secretly spy on you
  Innovation: anyone builds on top of others' work
  Cost: no license fees
  
Why closed source exists:
  Business model: code is the product
  Competitive advantage: hide your secret sauce
  Control: ensure a consistent experience
```

Most of the internet runs on open source software (Linux, Apache, nginx, MySQL, PostgreSQL, Python). Many companies build on open source and contribute back.

---

## 6. How Software Is Made (Overview)

Software creation is a craft. Here's the very brief overview (we'll go deeper in Chapter 21):

```
1. Problem / Idea
   "I want an app that helps people track their water intake"
   
2. Design
   What should it look like? What features? How does it work?
   Designers make wireframes (rough sketches of each screen)
   
3. Programming
   Developers write code (instructions for the computer)
   This takes the most time — weeks, months, years
   
4. Testing
   Is it buggy? Does it crash? Does it work on all devices?
   Testers find problems; developers fix them
   
5. Release
   App is published to the App Store or as a download
   
6. Maintenance
   Users find bugs. Developers patch them.
   New features are added.
   Software is never really "finished"
```

The people who write software are called **programmers**, **developers**, or **software engineers**. We'll look at what they actually do starting in Chapter 19.

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Software | Instructions telling hardware what to do |
| Application | Software you use directly (apps, games) |
| Operating system | Software that manages everything else |
| Firmware | Software baked into hardware |
| Open source | Code that anyone can read and modify |
| Closed source | Code kept secret (proprietary) |
| Web app | Software that runs in your browser, nothing to install |

**Now you know what software is. The most important piece of software on any computer is the operating system — let's meet it.**
