# Chapter 02: A Brief History of Operating Systems

> **"Every design decision in modern operating systems is a scar from a past failure. To understand why systems work the way they do today, you need to see what went wrong before."**

---

## Table of Contents

1. [The Era Before Operating Systems (1940s–1955)](#1-the-era-before-operating-systems-1940s1955)
2. [Batch Systems — The First OS (1955–1965)](#2-batch-systems--the-first-os-19551965)
3. [Time-Sharing — Multiple Users at Once (1960s–1970s)](#3-time-sharing--multiple-users-at-once-1960s1970s)
4. [Unix — The OS That Changed Everything (1969)](#4-unix--the-os-that-changed-everything-1969)
5. [Personal Computers — DOS and Early Windows (1970s–1990s)](#5-personal-computers--dos-and-early-windows-1970s1990s)
6. [The GUI Revolution — Macintosh and Windows (1984–1995)](#6-the-gui-revolution--macintosh-and-windows-19841995)
7. [Linux — The Open Source Revolution (1991)](#7-linux--the-open-source-revolution-1991)
8. [The Internet Era — Networked OS (1990s–2000s)](#8-the-internet-era--networked-os-1990s2000s)
9. [Mobile OS — iOS and Android (2007–present)](#9-mobile-os--ios-and-android-2007present)
10. [The Cloud Era and Beyond (2010s–present)](#10-the-cloud-era-and-beyond-2010spresent)
11. [Key Lessons From History](#11-key-lessons-from-history)
12. [Summary](#summary)

---

## 1. The Era Before Operating Systems (1940s–1955)

In the beginning, computers were programmed by physically rewiring them.

**ENIAC (1945):** One of the first general-purpose electronic computers. It had no OS. Programming it meant:
- Physically plugging cables into a panel (like an old telephone switchboard)
- Setting 3,000 switches by hand
- Running a single computation, then rewiring for the next

**The workflow:**
1. Programmer writes a mathematical formula on paper
2. Operators spend hours wiring the machine for that specific computation
3. Machine computes the result (in seconds or minutes)
4. Operators spend hours wiring for the NEXT computation

The computer was busy for maybe 30 minutes out of each day. The rest of the time humans were setting it up.

**Punch cards (early 1950s):**
IBM introduced stored-program computers where programs were punched onto cards. You'd hand a stack of punch cards to an operator, come back hours later, and collect your output (printed on paper). If you had a typo in your program, you'd wait another day to fix it.

**The problem:** Computers were enormously expensive ($1M–$10M in 1950s dollars). Having them sit idle while humans set up jobs was a massive waste.

---

## 2. Batch Systems — The First OS (1955–1965)

**The insight:** Instead of having humans setup one job at a time, collect many jobs and run them one after another, automatically.

This required software to control the transition from one job to the next — the very first operating system.

**FMS (Fortran Monitor System) and IBSYS (IBM):**
- Programs were submitted as batches of punch cards
- The OS (called a "monitor") loaded the next job automatically when the previous finished
- No human intervention needed between jobs

```
Batch processing workflow:
  Cards arrive → sorted into batches
  Batch queued on magnetic tape
  Computer reads job 1 → runs → writes output to tape
  Computer automatically reads job 2 → runs → writes output
  ...repeat for all jobs in batch...
  Output printed → given back to programmer hours later
```

**What this first "OS" provided:**
- **Job sequencing:** Automatically running the next program
- **I/O routines:** Standard code for reading cards, writing output (programmers didn't write their own I/O)
- **Error handling:** If a program crashed, the monitor recovered and ran the next job

**Drawback:** Still incredibly slow for the programmer. Submit a job in the morning, get results in the afternoon, find a bug, submit again, wait until tomorrow.

**CPU utilization problem:** Even with batching, the CPU was often idle while waiting for I/O (reading cards, writing output). I/O was 1,000× slower than the CPU.

---

## 3. Time-Sharing — Multiple Users at Once (1960s–1970s)

**The revolutionary idea:** Instead of running one program to completion, switch between multiple programs so quickly that each user THINKS they have the whole machine to themselves.

This is the concept of **multitasking** that every modern OS uses.

**MIT CTSS (Compatible Time-Sharing System, 1961):**
- Could support 30 simultaneous users on a single IBM 7090
- Each user got a terminal (keyboard + printer)
- OS switched between users every fraction of a second
- From each user's perspective, they had the whole computer

**Key technical challenges this introduced:**
- **Process management:** Track which program is running and what state it's in
- **Memory protection:** User A's program can't read User B's data
- **Scheduling:** Decide fairly which user runs next
- **File system:** Users need persistent storage for their programs and data

**Multics (1964–1969):**
MIT, Bell Labs, and GE built the most ambitious OS of its era:
- Multi-user, time-sharing
- Hierarchical file system (the idea of directories within directories)
- Virtual memory (programs could be larger than physical RAM)
- Online editing and compilation
- Security and access control

Multics was huge, complex, and expensive. Bell Labs pulled out of the project. But the ideas were enormously influential.

**Two Bell Labs researchers who left Multics decided to build something simpler...**

---

## 4. Unix — The OS That Changed Everything (1969)

**Ken Thompson** and **Dennis Ritchie** at Bell Labs built Unix on a discarded PDP-7 computer in 1969. They wanted the good ideas from Multics but without the complexity.

**Unix principles (still guide OS design today):**
1. **Everything is a file** — devices, pipes, network connections — all accessed like files
2. **Small programs that do one thing well** — ls just lists files; grep just searches
3. **Programs work together via pipes** — `ls | grep .c | wc -l` counts C source files
4. **Write programs as if humans will need to maintain them**

**1973: Dennis Ritchie rewrites Unix in C (not assembly):**
This was revolutionary. Before this, ALL OS code was written in assembly language — machine-specific, extremely difficult to port.

Writing Unix in C meant:
- You could take the Unix source code, recompile it for a new computer, and have a working OS
- Same OS could run on different hardware
- The OS was readable and maintainable by humans

This single decision (write OS in C) shaped computing for the next 50+ years. Linux, macOS, and even Windows kernel code is still mostly C today.

**Unix's influence:**
- **Every Linux system** is a Unix descendant (or follows Unix design principles)
- **macOS and iOS** run on XNU, which is directly derived from BSD Unix
- **Android** runs on a Linux kernel
- **The C programming language** itself was created to write Unix
- **POSIX standard** defines the Unix interface all compatible systems follow

**The Unix fork:**
In the 1970s–80s, Unix split into many variants:
- **BSD Unix** (Berkeley Software Distribution) — developed at UC Berkeley
- **System V Unix** — AT&T's commercial version
- **Solaris** (Sun Microsystems)
- **AIX** (IBM)
- **HP-UX** (HP)

---

## 5. Personal Computers — DOS and Early Windows (1970s–1990s)

By the mid-1970s, computers became small enough for individuals to own. The **microcomputer revolution** began.

**Apple II (1977), CP/M:**
Early personal computers ran simple OSes. CP/M was popular — just enough to load a program from disk and run it. No multitasking, no protection, no multi-user.

**IBM PC + MS-DOS (1981):**
IBM hired Microsoft to provide an OS for the IBM PC. Bill Gates bought QDOS (Quick and Dirty Operating System) for $50,000 and renamed it MS-DOS.

**MS-DOS design:**
- Single-user, single-tasking (runs one program at a time)
- No memory protection (any program could overwrite any memory)
- No graphical interface — command line only
- 640KB memory limit (because the designers thought nobody would ever need more)
- Programs communicated with hardware directly (no hardware abstraction!)

DOS was simple and worked. But it had a fundamental problem: **no protection**. Any program could crash the entire system by writing to arbitrary memory.

---

## 6. The GUI Revolution — Macintosh and Windows (1984–1995)

**Xerox PARC (1970s):**
Xerox researchers invented the graphical user interface: windows, icons, mouse, pointer (WIMP). They built the Alto (1973) as a research machine.

**Apple Macintosh (1984):**
Steve Jobs visited Xerox PARC, was inspired, and built the Macintosh — the first commercially successful GUI computer.
- Mac OS had GUI with overlapping windows, icons, menus
- Mouse-driven
- Still single-tasking, no memory protection
- Target: ordinary people, not technical users

**Microsoft Windows:**
- Windows 1.0 (1985): A GUI on top of DOS. Not a real OS — just a program.
- Windows 3.1 (1992): Still DOS underneath. Became hugely popular.
- **Windows NT 3.1 (1993):** This was the real OS. A completely new kernel, NOT based on DOS:
  - 32-bit preemptive multitasking
  - Memory protection
  - Multiple users
  - Designed for stability and security
  - Written by the VMS team (Dave Cutler) that Microsoft hired from DEC

**Windows NT became the foundation for everything modern Windows:**
- Windows 2000, XP (NT 5.x)
- Windows Vista, 7, 8, 10, 11 (NT 6.x, 10.x)

---

## 7. Linux — The Open Source Revolution (1991)

**Linus Torvalds**, a 21-year-old Finnish university student, started Linux in 1991 as a hobby project. He wanted to run a Unix-like system on his new 386 PC.

**August 25, 1991 — The famous newsgroup post:**
```
From: torvalds@klaava.Helsinki.FI (Linus Benedict Torvalds)
Subject: What would you like to see most in minix?

Hello everybody out there using minix -

I'm doing a (free) operating system (just a hobby, won't be big
and professional like gnu) for 386(486) AT clones. This has been
brewing since april, and is starting to get ready. I'd like any
feedback on things people like/dislike in minix, as my OS 
resembles it somewhat...
```

"Won't be big and professional" — the understatement of computing history.

**Linux is revolutionary because:**
- **Free and open source:** Anyone can read, modify, and distribute the code
- **Unix-compatible:** Followed Unix design principles
- **Community-driven:** Thousands of developers worldwide contributed
- **Portable:** Runs on x86, ARM, MIPS, RISC-V, PowerPC, and more

**Today Linux runs on:**
- ~70% of web servers (Apache, Nginx — all running on Linux)
- ~100% of the top 500 supercomputers
- ~100% of Android phones (Linux kernel)
- Most of the cloud (AWS, Google Cloud, Azure all run Linux underneath)
- Embedded systems (routers, TVs, cars, IoT devices)
- The International Space Station

---

## 8. The Internet Era — Networked OS (1990s–2000s)

The internet changed what OSes needed to do:

**TCP/IP stacks:** Every OS needed to include networking. The Berkeley Sockets API (from BSD Unix) became the standard — the same API is used today in every OS.

**Security became critical:**
- Multi-user systems connected to the internet were a target
- Buffer overflows, privilege escalation, network attacks
- OSes added: firewalls, memory protection, sandboxing

**Windows XP (2001):** The most successful OS version ever. Unified home and business Windows on the NT kernel. Over 1 billion PCs ran it. Also notoriously insecure (designed before internet security was a priority).

**macOS X (2001):** Apple killed the old Mac OS and built macOS X on a Unix (Darwin/XNU) foundation. Elegant UI on top of a solid Unix core.

---

## 9. Mobile OS — iOS and Android (2007–present)

Mobile devices required rethinking OS design:
- Battery life: wake up, do work, sleep immediately
- Touch input: completely different UI model
- App sandboxing: apps can't read each other's data
- Push notifications: efficient long-running communication
- Continuous location, camera, sensors

**iOS (2007):**
- Built on XNU (same kernel as macOS) + Darwin
- Massively simplified (no direct file system for users, no background processes for early versions)
- App sandbox: each app in its own isolated container
- Tight hardware-software integration

**Android (2008):**
- Linux kernel (modified, with changes like Binder IPC and Wakelocks)
- Dalvik VM / ART runtime for Java/Kotlin apps
- Open ecosystem (various manufacturers can use it)
- Over 3 billion active devices

**Mobile OS changed security thinking:**
- Apps are sandboxed by default (unlike desktop apps)
- Permission system (location, camera, microphone — user must approve)
- Signed apps from controlled stores
- Encrypted storage by default

---

## 10. The Cloud Era and Beyond (2010s–present)

**Containers (Docker, 2013):**
Instead of running full VMs, containers share the host OS kernel but have isolated user-space. Much lighter and faster. Linux namespaces + cgroups made this possible.

**Kubernetes (2014):**
OS-like management system for thousands of containers across thousands of machines.

**Serverless / FaaS:**
Programs don't even see an OS. They define a function; the cloud platform decides where and how to run it.

**WASM (WebAssembly):**
A new execution model — programs compiled to a safe intermediate format, run with near-native speed, fully sandboxed. Could change how we think about OS-application boundaries.

**Unikernels:**
Instead of a general-purpose OS, compile only the OS pieces your application needs into a single image. Smaller attack surface, faster boot, better performance.

**seL4 (Secure Microkernel):**
A formally verified microkernel — mathematically proven to be correct. Used in safety-critical systems.

---

## 11. Key Lessons From History

Looking at 80 years of OS development, clear patterns emerge:

**Lesson 1: Abstractions win long-term.**
DOS didn't abstract hardware → every program broke when hardware changed.
Unix abstracted everything → programs written in 1970s still compile today.

**Lesson 2: Security is retrofitted at great cost.**
DOS had no security → internet connected it → disaster (viruses, worms).
Better to design security in from the beginning (see: iOS).

**Lesson 3: Simple wins over complex.**
Multics was more technically ambitious than Unix. Unix won.
The "worse is better" philosophy (simpler, slightly less correct, but more widely deployable) often beats perfect but complex.

**Lesson 4: Portability is essential.**
Writing Unix in C made it portable → ran everywhere → became dominant.
Assembly-only OSes were tied to one hardware architecture.

**Lesson 5: Open source is more resilient.**
Linux survives because no single company controls it. The code belongs to everyone.

---

## Summary

| Era | Key Development | Why It Mattered |
|-----|----------------|----------------|
| 1940s–55 | No OS; manual wiring | Computers existed but were hard to use |
| 1955–65 | Batch systems | Automated job-to-job transition; first real OS |
| 1960s–70s | Time-sharing (Multics, CTSS) | Multiple users, preemptive multitasking |
| 1969 | Unix | Portable, clean design, C language, everything is a file |
| 1981 | MS-DOS | Personal computers; simple, no protection |
| 1984–91 | GUI (Mac, Windows), Linux | Graphical UI; open-source OS revolution |
| 2001 | macOS X, Windows XP | Unix foundation for desktop; NT for business |
| 2007–08 | iOS, Android | Mobile-first OS design; sandboxing standard |
| 2013+ | Containers, cloud, WASM | OS boundaries dissolving at the edges |
