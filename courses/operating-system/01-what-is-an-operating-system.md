# Chapter 01: What Is an Operating System?

> **"An operating system is the layer of software that turns a box of electronic components into something you can actually use."**

---

## Table of Contents

1. [The Problem Without an OS](#1-the-problem-without-an-os)
2. [What an OS Actually Does — Three Big Jobs](#2-what-an-os-actually-does--three-big-jobs)
3. [The OS as a Resource Manager](#3-the-os-as-a-resource-manager)
4. [The OS as an Abstraction Layer](#4-the-os-as-an-abstraction-layer)
5. [The OS as a Guardian — Protection and Security](#5-the-os-as-a-guardian--protection-and-security)
6. [A Day in the Life — What the OS Does Every Second](#6-a-day-in-the-life--what-the-os-does-every-second)
7. [What Is NOT the OS](#7-what-is-not-the-os)
8. [OS vs. Kernel — Are They the Same?](#8-os-vs-kernel--are-they-the-same)
9. [The Dual-Mode Operation](#9-the-dual-mode-operation)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Problem Without an OS

Imagine you have a brand new computer — just the hardware: CPU, RAM, hard drive, keyboard, monitor. You have a music player program you want to run.

**Without an OS, your music player would need to:**
- Know exactly which pins on the hard drive controller to toggle to read a music file
- Know exactly which memory addresses are safe to use (and not conflict with other programs)
- Know how to communicate with the sound card at the hardware level
- Know how to draw pixels on the screen using the display controller's registers
- Do all of this while ALSO doing the actual work of decoding MP3 audio

Every program would need to be an expert in every piece of hardware. Every new hard drive model, every new graphics card — every program would need to be rewritten.

This is exactly how early computers worked. Programs were written for specific machines. If you changed a hardware component, the programs often stopped working.

**The solution:** Put a layer of software between programs and hardware. This layer:
- Talks to the hardware on behalf of all programs
- Provides programs with a clean, simple interface
- Makes sure programs don't interfere with each other

That layer is the **operating system**.

---

## 2. What an OS Actually Does — Three Big Jobs

An operating system has exactly three fundamental jobs:

```
┌─────────────────────────────────────────────────────────────────┐
│                    USER PROGRAMS                                 │
│  (Chrome, Spotify, Word, your Python script, games, ...)       │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │  "I need memory, I need a file, I need the screen"
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                 OPERATING SYSTEM                                 │
│                                                                  │
│  Job 1: Manage Resources      (CPU time, RAM, disk, network)   │
│  Job 2: Provide Abstractions  (files, processes, sockets...)   │
│  Job 3: Enforce Protection    (program A can't read program B) │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │  "Here's what you asked for"
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                    HARDWARE                                      │
│  (CPU, RAM, Disk, Keyboard, Monitor, Network Card, GPU, ...)   │
└─────────────────────────────────────────────────────────────────┘
```

Let's understand each job:

---

## 3. The OS as a Resource Manager

Your computer has limited resources: one CPU (or a few cores), a fixed amount of RAM, one disk, one network card. Many programs want to use these at the same time.

**Who decides which program gets the CPU right now?**
The OS. It runs a scheduler that gives each program a turn.

**Who decides how much RAM each program gets?**
The OS. It has a memory manager that allocates and tracks memory.

**Who handles a program writing to disk at the same time another program is reading?**
The OS. It queues disk operations and makes sure they don't conflict.

**The OS is like a hotel manager:**
- You (a program) arrive and ask for a room (memory), access to the restaurant (CPU), and use of the phone (network)
- The manager allocates these resources, makes sure two guests don't get the same room, and handles conflicts

Without a resource manager, programs would fight over resources and crash each other. With a resource manager (the OS), they coexist peacefully.

**Resources the OS manages:**
| Resource | OS component that manages it |
|---------|------------------------------|
| CPU time | Process Scheduler |
| Physical RAM | Memory Manager |
| Hard disk / SSD | File System + Block I/O layer |
| Network (Ethernet, WiFi) | Network Stack |
| GPU | Graphics subsystem |
| USB, Bluetooth | Device driver subsystem |
| Keyboard/Mouse | Input subsystem |

---

## 4. The OS as an Abstraction Layer

The second job of an OS is to **hide complexity** behind simple interfaces.

**Example: Reading a file**

Without OS abstraction, to read a text file from an SSD, a program would need to:
1. Issue an NVMe command to the SSD controller
2. Specify the LBA (Logical Block Address) of the storage sectors
3. Set up a DMA buffer in RAM
4. Wait for an interrupt
5. Decode the NAND flash wear leveling translation
6. ...and dozens more hardware-specific steps

**With OS abstraction:**
```c
// This is all the program sees:
FILE *f = fopen("hello.txt", "r");
char buf[100];
fread(buf, 1, 100, f);
fclose(f);
```

The OS translates this simple request into all the hardware operations needed. The program doesn't know or care whether the file is on an SSD, HDD, USB drive, or network share.

**Key abstractions the OS provides:**

| What Programs See | What's Really Happening |
|-------------------|------------------------|
| **File** | Sectors on a disk arranged by a file system |
| **Process** | A program loaded into memory with CPU time slices |
| **Thread** | A sequence of instructions with its own stack |
| **Socket** | A TCP/UDP connection over a network card |
| **Memory pointer** | A virtual address mapped to physical RAM by the MMU |
| **Screen** | A framebuffer that the display controller reads |
| **Keyboard input** | Electrical signals from key presses, decoded and buffered |

**The power of abstraction:** A program written on Linux can use `write(fd, buf, n)` to write to a file, a pipe, a network socket, or a terminal — all the same call. The OS handles the difference underneath.

---

## 5. The OS as a Guardian — Protection and Security

The third job of an OS: make sure programs can't interfere with each other or with the OS itself.

**Why protection is necessary:**

Imagine you have a bank app and a music player running at the same time. Without protection:
- The music player could accidentally (or intentionally) read the bank app's memory — stealing your password
- A buggy program could write to the OS's own memory — crashing the entire system
- A program could hog 100% of the CPU forever — freezing everything else

**How the OS protects:**

1. **Memory isolation:** Each process gets its own virtual address space. Process A cannot access Process B's memory. If it tries, the CPU raises a fault → OS kills the misbehaving process.

2. **Privilege levels:** The CPU has hardware-enforced "rings." The OS runs in Ring 0 (highest privilege). Programs run in Ring 3 (lowest privilege). Programs physically cannot execute privileged instructions (like disabling interrupts or accessing hardware directly).

3. **Controlled access to hardware:** Programs can't talk to hardware directly. They must ask the OS via system calls. The OS checks permissions before doing anything.

4. **File permissions:** The file system stores who can read/write/execute each file. The OS enforces these permissions.

---

## 6. A Day in the Life — What the OS Does Every Second

Let's trace one second of your computer running, showing what the OS is doing constantly:

```
You're watching a YouTube video in Chrome:

EVERY 1ms (1000 times per second):
  Timer interrupt fires
  → OS scheduler runs
  → Decides: "Chrome gets the next 10ms, Spotify gets next 2ms, 
               background updates get 0.5ms..."
  → Switches CPU from one process to another

EVERY TIME Chrome needs a video frame:
  Chrome asks OS: "give me 4MB of memory for this frame"
  OS allocates virtual pages → maps to physical RAM
  Chrome decodes frame → OS sends to GPU driver → pixels appear

EVERY TIME audio plays:
  Audio driver DMA transfer → circular buffer
  OS wakes up audio thread every 5ms for more data
  Spotify fills buffer → OS hands to sound card

EVERY keypress:
  Keyboard controller raises interrupt (IRQ 1)
  OS keyboard ISR runs → reads scan code → converts to key event
  OS puts key event in Chrome's input queue
  Chrome's event loop wakes up → processes keystroke

EVERY network packet arriving:
  Network card raises interrupt
  OS NIC driver runs → copies packet to socket buffer
  OS notifies Chrome (via epoll/select) that data arrived
  Chrome reads from socket

ALL OF THIS HAPPENS SIMULTANEOUSLY, invisibly, in the background.
```

The OS is running millions of these tiny operations per second. You don't see any of it. That's the whole point.

---

## 7. What Is NOT the OS

People often confuse the OS with other things. Let's be precise:

**The OS is NOT the graphical user interface (GUI):**
- On Windows, the "OS" is Windows NT kernel. The Explorer shell, taskbar, windows — that's a GUI built on top of the OS.
- Linux kernel is the OS. GNOME, KDE, i3 — those are GUI environments that run ON TOP of the OS.
- macOS: the XNU kernel is the OS. The Finder, Dock — those are applications.

**The OS is NOT the applications:**
- Chrome, Word, Spotify — NOT part of the OS. They're programs that run on top of the OS.

**The OS is NOT the C library (libc):**
- glibc, musl — these are user-space libraries that make it easier to call OS functions.
- They're not the OS, but they're very closely tied to it.

**The OS IS:**
- The kernel (the core software that talks to hardware directly)
- The system call interface
- Basic drivers for essential hardware
- The scheduler, memory manager, file system

---

## 8. OS vs. Kernel — Are They the Same?

You'll often hear "OS" and "kernel" used interchangeably. They're related but not identical:

**Kernel:**
The kernel is the core of the OS — the piece that runs with full hardware privilege, manages all resources, and provides the foundation. It's one program that loads at boot and stays in memory forever.

**Operating System:**
Usually means the kernel PLUS system software needed to make a working environment:
- Standard libraries (C library, math library)
- System daemons (cron job scheduler, syslog, networking daemon)
- Device drivers (though these can be inside or outside the kernel)
- Basic utilities (shell, ls, cp, init system)

```
"Linux" usually means:

  Linux kernel            ← The actual kernel (Linus Torvalds' code)
  +
  GNU utilities           ← ls, cp, bash, gcc, etc. (GNU project)
  +
  Package manager         ← apt, yum, etc.
  +
  Init system             ← systemd or SysV init
  +
  Desktop environment     ← GNOME, KDE (optional)
  =
  "Linux distribution" (Ubuntu, Fedora, Debian, etc.)
```

In this course, when we say "OS," we usually mean the kernel and the essential system software around it. When we build our OS, we'll build primarily the kernel.

---

## 9. The Dual-Mode Operation

This is a crucial concept you'll encounter throughout this course:

**Modern CPUs have (at least) two operating modes:**

**Kernel Mode (Privileged Mode, Ring 0):**
- Can execute ANY instruction
- Can access ANY memory address
- Can enable/disable interrupts
- Can directly access hardware
- The OS kernel runs here

**User Mode (Unprivileged Mode, Ring 3):**
- CANNOT execute privileged instructions
- CANNOT access OS memory (access fault → exception)
- CANNOT directly touch hardware
- All user programs run here

**How does a program get OS services then?**
It uses a **system call** — a special instruction that switches from user mode to kernel mode in a controlled way:

```
Program in user mode:
  "I need to read a file"
  
  → executes: syscall (or int 0x80 on old x86)
  
CPU automatically:
  1. Switches to kernel mode
  2. Saves program's state
  3. Jumps to OS system call handler
  
OS handler:
  "This is read() syscall with fd=3, buf=0x7fff..., count=100"
  - Checks if fd is valid
  - Reads from disk
  - Copies data to program's buffer
  - Returns result code
  
CPU:
  4. Switches back to user mode
  5. Restores program's state
  
Program continues:
  "OK, I got 100 bytes in my buffer"
```

The boundary between kernel mode and user mode is the most important boundary in any operating system. Crossing it always goes through a controlled gate (system call, exception, interrupt). This is what makes the OS secure.

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Operating System | Software layer between hardware and user programs |
| Resource Management | OS allocates CPU, RAM, disk, network among programs |
| Abstraction | OS hides hardware complexity behind simple interfaces |
| Protection | OS ensures programs can't corrupt each other or the OS |
| Kernel | The core of the OS; runs with full hardware privilege |
| User Mode | Restricted mode where programs run; can't touch hardware |
| Kernel Mode | Privileged mode where the OS runs; can do anything |
| System Call | The controlled gate from user mode into kernel mode |

---

## Exercises

1. **Understanding abstractions:** You open a file in Python with `open("data.txt", "r")`. List at least 5 things the OS is doing behind the scenes that the Python program doesn't know about (think: file system lookup, disk I/O, memory allocation, etc.).

2. **Resource management thinking:** On a computer with 8GB RAM, Chrome needs 2GB, Photoshop needs 1.5GB, and a video game needs 3GB. That's 6.5GB total which fits. Now Chrome needs another 2GB. Describe what the OS might do: (a) refuse the request, (b) use swap space, (c) ask another program to free memory. What are the tradeoffs of each?

3. **Dual-mode protection:** Explain why it would be dangerous if user programs could run in kernel mode. Give two specific attack scenarios that would be possible.

4. **Abstraction vs. performance:** Abstractions make programming easier but add overhead. A database application might want to bypass the OS file system and talk to the disk directly for maximum performance. What is this called? (Research: "raw disk access" or "direct I/O"). What does the OS need to provide to allow this safely?

5. **OS identification:** For each of the following, identify whether it's the OS kernel, a system library, or a user application: (a) Linux kernel, (b) Microsoft Word, (c) glibc, (d) systemd, (e) GNOME desktop, (f) the VFS layer in Linux, (g) /bin/bash (the shell), (h) the NTFS driver in Windows.
