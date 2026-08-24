# Chapter 00: How Computers Work — For the Complete Beginner

*Before you can hack a computer, you need to understand what a computer actually is. This chapter assumes you know absolutely nothing about how computers work. By the end, you'll understand what's really happening when you click a button or open a website.*

---

## What Is a Computer?

A computer is a machine that:
1. Stores information (data)
2. Processes that information (runs programs)
3. Takes input (keyboard, mouse, network)
4. Produces output (screen, files, network)

Everything a computer does — playing games, browsing the internet, running your bank's servers — is ultimately just: read data, do math on it, write data back.

That's it. The complexity comes from the scale and speed at which this happens.

---

## The Physical Parts — Hardware

### The CPU (Central Processing Unit)

The CPU is the "brain" of the computer. It executes instructions — billions of them per second.

An instruction is something like:
- "Add these two numbers"
- "Copy this value from memory to here"
- "If this number is zero, jump to instruction 5,000"

Modern CPUs (Intel, AMD, Apple Silicon) execute 3-5 billion instructions per second on each core. A typical CPU has 4-16 cores.

**Why this matters for security:** Everything a running program does — including malware — is instructions being executed by the CPU. Understanding the CPU helps you understand how programs work and how they can be exploited.

### RAM (Random Access Memory)

RAM is where programs live while they're running. It's fast, temporary storage.

When you open Chrome, it loads from your hard drive into RAM. Chrome runs in RAM. When you close Chrome, its RAM is freed.

Think of RAM like a desk. Your hard drive is a filing cabinet. You work with papers on your desk (RAM). Papers you're not currently using go in the filing cabinet (hard drive).

**Size:** A typical computer has 8-32 GB of RAM. Each byte in RAM has an address — a number from 0 to (RAM size).

**Why this matters for security:** Most attacks — buffer overflows, injection attacks, privilege escalation — work by manipulating what's in RAM. Understanding memory is fundamental to understanding exploitation.

### Storage (Hard Drive / SSD)

Long-term storage. Persists when the computer is off.

Your programs, files, and operating system live here. An SSD (Solid State Drive) is fast flash storage; a traditional hard drive uses spinning magnetic disks.

**Why this matters for security:** Evidence of attacks is on disk (log files, malware files). File system forensics reads the disk to understand what happened.

### Network Card (NIC)

Connects the computer to a network. Converts between computer data (0s and 1s) and network signals (electrical, optical, or radio waves).

Every network card has a MAC address — a globally unique 48-bit hardware identifier (e.g., `AA:BB:CC:DD:EE:FF`).

### GPU (Graphics Processing Unit)

Originally for graphics. Now also used for machine learning and — importantly for security — password cracking. A high-end GPU can try billions of password hashes per second.

---

## The Software Stack

Software is the programs that run on hardware. There are layers:

```
┌─────────────────────────────────┐
│     Your Application (Chrome)   │
├─────────────────────────────────┤
│     Programming Libraries       │
├─────────────────────────────────┤
│     Operating System (Linux)    │
├─────────────────────────────────┤
│     BIOS/UEFI Firmware          │
├─────────────────────────────────┤
│     Hardware (CPU, RAM, Disk)   │
└─────────────────────────────────┘
```

### The Operating System (OS)

The operating system (Linux, Windows, macOS) is the master program that:
- Manages hardware (allocates CPU time, RAM, disk access)
- Runs other programs
- Controls who can access what (security)
- Manages the file system (folders and files)
- Handles network connections

**Why Linux?** Most servers, cloud infrastructure, and IoT devices run Linux. Hacking tools are primarily built for Linux. Android is Linux-based. We'll use Linux throughout this course.

### The Kernel

The most privileged part of the OS. Runs in "kernel space" — it has direct hardware access. Everything else runs in "user space" with restricted access.

Programs running in user space must ask the kernel to do privileged things (access the network, read files, create processes) through "system calls."

**Why this matters for security:** The kernel is the ultimate gatekeeper. Escalating privileges means getting access to kernel-level powers. Rootkits hide themselves by modifying the kernel. EDR tools (which you'll build) often hook kernel calls to detect attacks.

---

## How a Program Runs

Let's trace what happens when you run a program:

1. **You double-click an executable** (e.g., `chrome.exe`)

2. **OS reads the file from disk** and loads it into RAM

3. **OS creates a process** — a running instance of the program with its own memory space

4. **CPU starts executing instructions** from the start of the program

5. **Program makes system calls** to ask the OS for things:
   - "Open this file"
   - "Connect to this network address"
   - "Create a new thread"

6. **CPU continues executing** instructions, following branches and loops

7. **Program exits** — memory is freed, process is removed

**The key insight:** Every program, including malware, does exactly this. A ransomware program opens your files (system call), encrypts them (CPU instructions), and writes them back (system call). Understanding this flow lets you detect it.

---

## How Memory Works

When a program runs, its memory is organized into regions:

```
High addresses
┌─────────────┐
│    Stack    │ ← local variables, function calls
├─────────────┤
│      ↓      │
│      ↑      │
├─────────────┤
│    Heap     │ ← dynamically allocated memory
├─────────────┤
│    Data     │ ← global variables
├─────────────┤
│    Code     │ ← program instructions
└─────────────┘
Low addresses
```

**Stack:** Where local variables live. Grows down as functions are called, shrinks as they return. The stack is where buffer overflows happen.

**Heap:** Memory you request dynamically (like `malloc()` in C). Grows up as you allocate.

**Why this matters for security:**
- Buffer overflows overflow the stack, overwriting return addresses
- Use-after-free bugs misuse freed heap memory
- Shellcode is injected into writable, executable memory regions

---

## Bits, Bytes, and Representation

Everything in a computer is stored as binary — 0s and 1s.

**Bit:** A single 0 or 1
**Byte:** 8 bits together (can represent values 0-255)
**Kilobyte (KB):** 1,024 bytes
**Megabyte (MB):** 1,024 KB
**Gigabyte (GB):** 1,024 MB

The letter 'A' in the most common encoding (ASCII) is stored as the byte value 65 (binary: 01000001).

**Why this matters:** When you analyze malware, network traffic, or memory dumps, you're looking at raw bytes. Understanding binary and hex (we cover this in Chapter 01) lets you read what you're looking at.

---

## Input and Output

Computers receive input from:
- **Keyboard/Mouse:** User interaction
- **Network:** Data from other computers
- **Disk:** Files being read
- **Sensors:** (IoT devices)

Computers produce output through:
- **Display:** What you see on screen
- **Network:** Data sent to other computers
- **Disk:** Files being written

**The security insight:** Almost all attacks enter through input channels — network packets, file opens, keyboard input. Defenses focus on validating and sanitizing input.

---

## What a "Hack" Actually Is

When someone "hacks" a computer, they're doing one or more of:

1. **Exploiting a bug:** Making a program behave in an unintended way by providing unexpected input
2. **Misusing a legitimate feature:** Using a feature correctly but for a malicious purpose
3. **Social engineering:** Tricking a human rather than exploiting a machine
4. **Credential theft:** Obtaining usernames/passwords to authenticate legitimately
5. **Configuration errors:** Exploiting systems that are misconfigured by their administrators

None of this is magic. It's all based on understanding how the system works and finding where that understanding breaks down or can be abused.

---

## Summary

| Component | What it does | Security relevance |
|-----------|-------------|-------------------|
| CPU | Executes instructions | Executes both legitimate code and malware |
| RAM | Holds running programs | Target for exploitation, holds evidence |
| Disk | Persistent storage | Contains malware files, forensic evidence |
| Network card | Connects to network | Entry point for most attacks |
| OS Kernel | Manages everything | Ultimate privilege target for attackers |

---

## What's Next

In the next chapter, you'll learn binary and hexadecimal — the languages computers actually speak. Then we'll cover networking, Linux, and start writing Go code to build security tools.

**Exercise:** Open your computer's task manager (Windows: Ctrl+Shift+Esc, Linux: `top` or `htop`, Mac: Activity Monitor). Find a familiar program. What do you notice about CPU and memory usage?
