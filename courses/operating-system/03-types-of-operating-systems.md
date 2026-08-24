# Chapter 03: Types of Operating Systems

> **"There is no one-size-fits-all OS. A medical device OS that controls an insulin pump has completely different requirements from an OS on a gaming PC. Understanding the type of OS helps you understand every design decision within it."**

---

## Table of Contents

1. [Why Different Types Exist](#1-why-different-types-exist)
2. [Desktop / Workstation OS](#2-desktop--workstation-os)
3. [Server OS](#3-server-os)
4. [Mobile OS](#4-mobile-os)
5. [Real-Time OS (RTOS)](#5-real-time-os-rtos)
6. [Embedded OS](#6-embedded-os)
7. [Distributed OS](#7-distributed-os)
8. [Network OS](#8-network-os)
9. [Batch OS (Legacy but important)](#9-batch-os-legacy-but-important)
10. [Hypervisors / Type-1 OS](#10-hypervisors--type-1-os)
11. [OS Kernels by Architecture: Monolithic vs. Microkernel vs. Hybrid](#11-os-kernels-by-architecture)
12. [Summary Table](#summary)

---

## 1. Why Different Types Exist

Every OS makes design tradeoffs. The right tradeoffs depend on what the OS needs to do:

| Design Goal | What It Means | Who Needs It |
|-------------|--------------|-------------|
| **Responsiveness** | React to events within milliseconds | Gaming, GUI, RTOS |
| **Throughput** | Complete as many tasks as possible | Servers, batch processing |
| **Reliability** | Never crash, ever | Medical, aviation, industrial |
| **Security** | Prevent unauthorized access | Banking, military, mobile |
| **Power efficiency** | Use as little battery as possible | Mobile, IoT, laptops |
| **Compatibility** | Run old software unchanged | Windows, enterprise |
| **Simplicity** | Fit in tiny RAM/flash | Microcontrollers, IoT |

No single OS can optimize ALL of these simultaneously. Tradeoffs must be made. That's why different types exist.

---

## 2. Desktop / Workstation OS

**Purpose:** General-purpose computing for a human user sitting at a computer.

**Examples:** Windows 10/11, macOS, Ubuntu Linux, Fedora, ChromeOS

**Design priorities:**
- Rich graphical user interface
- Wide application compatibility (run many different programs)
- Responsiveness (UI should feel instant)
- Multitasking (Chrome, Spotify, Word all at once)
- Hardware support (thousands of devices: printers, webcams, GPUs)
- Relatively good security (but not military-grade)

**What these OSes sacrifice:**
- Hard real-time guarantees (Windows may take 50ms to respond to a keypress if busy)
- Minimal footprint (Windows takes 30GB+ of disk space)
- Absolute security (tradeoff: run arbitrary user apps)

**Key characteristics:**
- Preemptive multitasking (OS forcibly takes CPU from one process and gives to another)
- Virtual memory (each process has its own address space)
- Hardware abstraction layer (run on many hardware configurations)
- User account system with basic permissions

---

## 3. Server OS

**Purpose:** Run services for many clients; maximize throughput and uptime.

**Examples:** Ubuntu Server, Red Hat Enterprise Linux (RHEL), CentOS, Windows Server, FreeBSD

**Design priorities:**
- **Throughput over responsiveness:** Serving 10,000 web requests per second matters more than fast GUI
- **Stability:** Must run for months/years without reboot (5-nines uptime: 99.999% = ~5 min downtime/year)
- **Security:** Exposed to the internet; must be hardened
- **Scalability:** Handle more load by adding more CPU/RAM
- **Remote management:** No physical keyboard — managed over SSH

**What server OSes drop:**
- Usually no GUI (run headless — no monitor attached)
- No media players, games, desktop apps
- Often minimal default services (reduce attack surface)

**Server OS design choices:**
- CPU scheduler optimized for throughput (run processes longer before switching)
- Large buffer caches for disk I/O
- Networking tuned for high connection counts
- Memory overcommit enabled (can allocate more memory than physically exists — works statistically)

**Linux dominates here.** Why? Free, stable, highly configurable, excellent networking.

---

## 4. Mobile OS

**Purpose:** Computing on battery-powered handheld devices with touchscreens.

**Examples:** iOS (Apple), Android (Google), HarmonyOS (Huawei)

**Design priorities:**
- **Battery life:** Every CPU cycle costs battery; aggressively sleep when not in use
- **Touch UI:** No mouse; finger input with different ergonomics
- **App security:** Apps from unknown developers must be sandboxed
- **Connectivity:** Always-on network (background sync, push notifications)
- **Sensors:** GPS, accelerometer, camera, microphone — integrated at OS level
- **Instant-on:** User picks up phone; must be responsive immediately

**Revolutionary design decisions in mobile OS:**
- **App sandbox:** Each app has its own isolated directory; can't read other apps' files
- **Permission system:** User explicitly grants camera/location/contacts access
- **Background restrictions:** Apps can't run freely in background (kills battery)
- **OTA updates:** OS updates delivered silently, over the air
- **Trusted app store:** Apps reviewed before distribution (iOS); or signed (Android)

**Power management is the core challenge:**
- CPU has multiple power states (C0 through C8 — idle states)
- Mobile CPU can clock down to MHz range when idle
- OS tracks all wakeup sources ("wakelocks" in Android)
- Even a poorly written app checking the network every second kills battery

---

## 5. Real-Time OS (RTOS)

**Purpose:** Guarantee response to events within a specific time deadline.

**Examples:** FreeRTOS, VxWorks, QNX, Zephyr, RTEMS, LynxOS

**The crucial concept — what "real-time" means:**

"Real-time" does NOT mean "fast." It means **predictable**.

- A non-real-time OS might respond to a sensor reading in 1ms on average, but occasionally take 50ms (if the scheduler was busy).
- A real-time OS guarantees it will ALWAYS respond within, say, 5ms — even in the worst case.

**Why this matters:**

| Application | What happens if deadline missed |
|-------------|--------------------------------|
| Airbag controller | Airbag deploys too late → passenger injury |
| Pacemaker | Heart gets wrong timing signal → dangerous |
| Anti-lock brakes | ABS can't respond fast enough → longer stopping distance |
| Industrial robot | Wrong timing → physical collision, safety hazard |
| Power grid control | Wrong switching timing → equipment damage |

**RTOS design principles:**
- **Deterministic scheduler:** Given a task priority, worst-case response time is bounded and calculable
- **Priority inversion avoidance:** High-priority task must never wait for low-priority task indefinitely
- **Minimal interrupt latency:** Hardware interrupt → ISR runs within microseconds
- **Small and predictable:** No garbage collectors, no virtual memory tricks that cause unpredictable delays

**Hard RTOS vs. Soft RTOS:**
- **Hard real-time:** MUST meet deadline. Missing = system failure. (Medical devices, aircraft)
- **Soft real-time:** Should meet deadline; occasional miss is degraded but not catastrophic. (Video streaming — one dropped frame is OK)

**FreeRTOS example (runs on Arduino-class microcontrollers):**
- 9KB of RAM is enough to run FreeRTOS
- Deterministic priority-based preemptive scheduler
- Used in billions of IoT devices, medical equipment, industrial controllers

---

## 6. Embedded OS

**Purpose:** Run on resource-constrained hardware inside a specific product.

**Examples:** FreeRTOS, Contiki, RIOT, ThreadX, embedded Linux (Yocto), bare-metal (no OS)

**"Embedded" means the OS is embedded inside a product:**
- The OS in your microwave
- The OS in your car's ECU (engine control unit)
- The OS in your router
- The OS in a smart thermostat

**Constraints:**
- **RAM:** Often 4KB–256KB (vs 8GB on a laptop)
- **Flash:** Often 32KB–4MB for the OS + application
- **No filesystem:** Some embedded systems have no file system at all
- **No MMU:** Many microcontrollers (like Cortex-M0) have no Memory Management Unit → no virtual memory, no memory protection
- **Power:** May run on a battery for 10 years
- **No display/keyboard:** Devices have a specific function, not general-purpose UI

**Embedded OS vs. RTOS:**
These overlap but aren't the same:
- An embedded OS focuses on small footprint and specific hardware
- An RTOS focuses on timing guarantees
- Many embedded systems use both (FreeRTOS is both)

**Bare-metal programming:**
Some embedded systems have NO OS at all. The program runs directly on the hardware:
- Arduino programs (before RTOS is added): one main() loop, no OS
- Simple sensor nodes
- Boot ROM code

---

## 7. Distributed OS

**Purpose:** Coordinate multiple networked computers to appear as a single system.

**Examples:** Google's Borg (runs Google's services), Plan 9 (research OS), MOSIX

**The idea:**
Instead of one powerful computer, use many ordinary computers. The OS makes them look and act like one computer to applications.

**Challenges unique to distributed OS:**
- **Network is not reliable:** Messages can be lost, delayed, duplicated
- **Clocks are not synchronized:** Different machines disagree on what time it is
- **Partial failure:** Some machines can fail while others continue
- **Consistency:** Is the file I see on machine A the same as what machine B sees?

**CAP Theorem** (fundamental limit of distributed systems):
A distributed system can have at most TWO of: Consistency (everyone sees same data), Availability (system always responds), Partition tolerance (works even if network splits).

**Examples of distributed systems you use daily:**
- Google Search (thousands of servers appear as one service)
- Netflix streaming (distributed storage + delivery)
- WhatsApp (distributed messaging at scale)
- Bitcoin (no central server — fully distributed)

---

## 8. Network OS

**Purpose:** Manage a network of computers; focus on network services.

**Examples:** Cisco IOS (router OS), Junos (Juniper routers), network-attached storage (NAS) OS

**Network OS ≠ Networked OS:**
- A **network OS** is an OS built specifically to run network infrastructure (routers, switches, firewalls)
- A **networked OS** is any OS with networking support (almost all modern OSes)

**Cisco IOS (Internetwork Operating System):**
- Runs in every Cisco router and many switches
- CLI-based configuration
- Not Linux (proprietary)
- Real-time packet forwarding at line rate (millions of packets per second)
- Focus: routing protocols (OSPF, BGP), packet forwarding, QoS

---

## 9. Batch OS (Legacy but Important)

**Purpose:** Process large volumes of jobs without human interaction, in order.

**Modern examples:** Mainframe OS (IBM z/OS), job schedulers (PBS, SLURM for HPC clusters)

**Batch OS is still used for:**
- Payroll processing (run overnight, process 100,000 employee records)
- Bank statement generation
- Scientific simulations (HPC clusters with SLURM job scheduler)
- Data warehouse ETL (Extract, Transform, Load) jobs

**Characteristics:**
- Jobs queued → executed in order (or by priority)
- No interactive user during execution
- High throughput (process as many jobs as possible)
- Maximize CPU utilization

---

## 10. Hypervisors / Type-1 OS

**Purpose:** Run multiple complete operating systems simultaneously on one machine.

**Examples:** VMware ESXi, Microsoft Hyper-V, Xen, KVM (Linux)

**Two types:**

**Type 1 (Bare-metal hypervisor):**
```
┌─────────────────────────────────────────────────────────┐
│  VM1 (Windows)  │  VM2 (Linux)  │  VM3 (macOS)         │
│  (full OS)      │  (full OS)    │  (full OS)            │
├─────────────────────────────────────────────────────────┤
│              HYPERVISOR (Type 1)                        │
├─────────────────────────────────────────────────────────┤
│                 HARDWARE                                 │
└─────────────────────────────────────────────────────────┘
```
Runs directly on hardware. No host OS.
Examples: VMware ESXi (used in data centers), Xen (used by AWS EC2), Hyper-V

**Type 2 (Hosted hypervisor):**
```
┌─────────────────────────────────────────────────────────┐
│   VM1 (Windows)       │  VM2 (Linux)                   │
├───────────────────────┴────────────────────────────────┤
│              HYPERVISOR (Type 2)                        │
├─────────────────────────────────────────────────────────┤
│              HOST OS (macOS / Windows / Linux)          │
├─────────────────────────────────────────────────────────┤
│                 HARDWARE                                 │
└─────────────────────────────────────────────────────────┘
```
Runs as an application on top of a host OS.
Examples: VirtualBox, VMware Workstation, Parallels (macOS)

**Why hypervisors matter:**
- Cloud computing: one physical machine runs hundreds of customer VMs
- Development: test on Windows while running macOS
- Security: run untrusted code in isolated VM
- Consolidation: 20 underutilized servers → 1 server running 20 VMs

---

## 11. OS Kernels by Architecture

Within any OS type, the kernel itself can be designed in different ways:

**Monolithic Kernel:**
```
User Space:  | App | App | App |
-------------|-----|-----|-----|
Kernel Space:  [Everything: scheduler, memory, FS, drivers, network]
              all running with full kernel privilege in one big program
```
- All kernel code runs in kernel space
- Subsystems communicate by direct function calls (fast!)
- One bug can crash the entire OS
- **Examples:** Linux kernel, original Unix

**Microkernel:**
```
User Space:  | App | FS Server | Device Driver | Network Stack |
-------------|-----|-----------|---------------|---------------|
Kernel Space: [Only: IPC, basic scheduling, address spaces]
              Tiny kernel; everything else is a user-space server
```
- Kernel does ONLY: memory protection, IPC, basic scheduling
- Everything else (file system, drivers, network) runs as user-space servers
- One driver crashing doesn't crash the kernel
- IPC overhead makes it slower
- **Examples:** Minix, QNX, L4, seL4, GNU Hurd

**Hybrid Kernel:**
```
User Space:  | App | Some subsystems |
-------------|-----|-----------------|
Kernel Space: [Core + some drivers/subsystems for performance]
```
- Mix of monolithic and microkernel approaches
- Key services in kernel, but structured for maintainability
- **Examples:** Windows NT kernel, macOS XNU (Mach microkernel + BSD monolithic)

---

## Summary

| OS Type | Examples | Key Priority | Sacrifice |
|---------|---------|-------------|----------|
| Desktop | Windows, macOS, Ubuntu | Rich UI, app compat | Hard real-time, minimal size |
| Server | RHEL, Ubuntu Server | Throughput, uptime | GUI, interactivity |
| Mobile | iOS, Android | Battery, security, touch | Background processes, flexibility |
| RTOS | FreeRTOS, VxWorks, QNX | Timing guarantees | Complexity, features |
| Embedded | Contiki, RIOT, ThreadX | Tiny size, low power | General-purpose features |
| Distributed | Plan 9, Borg | Scale across machines | Consistency complexity |
| Hypervisor | ESXi, Hyper-V, Xen | Isolation of multiple OSes | Performance overhead |

| Kernel Type | Examples | Advantage | Disadvantage |
|------------|---------|-----------|-------------|
| Monolithic | Linux, BSD | Fast (direct function calls) | One bug = crash |
| Microkernel | seL4, QNX, Minix | Fault isolated | IPC overhead |
| Hybrid | Windows NT, macOS | Balance | Complexity |
