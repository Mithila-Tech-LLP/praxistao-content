# Chapter 04: How Your Computer Really Works

> **"The OS can only manage what it understands. Before writing a single line of OS code, you need to know exactly what hardware the OS is managing — what every component does and how they talk to each other."**

---

## Table of Contents

1. [The Big Picture — What's Inside a Computer](#1-the-big-picture--whats-inside-a-computer)
2. [The CPU — The Brain](#2-the-cpu--the-brain)
3. [Registers — The CPU's Working Memory](#3-registers--the-cpus-working-memory)
4. [The Memory Hierarchy](#4-the-memory-hierarchy)
5. [RAM — Main Memory](#5-ram--main-memory)
6. [The Bus System — How Components Communicate](#6-the-bus-system--how-components-communicate)
7. [The Memory Map — What's at Which Address](#7-the-memory-map--whats-at-which-address)
8. [The Stack and the Heap](#8-the-stack-and-the-heap)
9. [I/O — Input and Output](#9-io--input-and-output)
10. [DMA — Direct Memory Access](#10-dma--direct-memory-access)
11. [Interrupts — Hardware Events](#11-interrupts--hardware-events)
12. [BIOS, UEFI, and Firmware](#12-bios-uefi-and-firmware)
13. [Summary](#summary)

---

## 1. The Big Picture — What's Inside a Computer

Every computer — whether a laptop, phone, server, or microcontroller — is built from the same fundamental components:

```
┌──────────────────────────────────────────────────────────────────────┐
│                         COMPUTER SYSTEM                               │
│                                                                        │
│  ┌──────────┐   ┌───────────────────────────────────────────────────┐│
│  │          │   │                  MAIN MEMORY (RAM)                 ││
│  │   CPU    │←→ │  [OS code][OS data][process 1][process 2][stack]  ││
│  │          │   └───────────────────────────────────────────────────┘│
│  └─────┬────┘                                                         │
│        │  System Bus (address + data + control lines)                │
│        ↓                                                               │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │                      I/O SUBSYSTEM                          │     │
│  │  ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐  ┌──────┐ │     │
│  │  │ Disk   │  │ NIC    │  │ GPU    │  │ USB    │  │ kbd  │ │     │
│  │  │(SSD)   │  │(Net)   │  │(Video) │  │        │  │      │ │     │
│  │  └────────┘  └────────┘  └────────┘  └────────┘  └──────┘ │     │
│  └─────────────────────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────────────────────┘
```

Every component has a role. Understanding each one is essential to understanding the OS.

---

## 2. The CPU — The Brain

The **CPU (Central Processing Unit)** is the only component that actually executes instructions. Everything else — RAM, disk, GPU — exists to serve the CPU.

**What the CPU does in one clock cycle (a "fetch-decode-execute" cycle):**

```
Step 1: FETCH
  CPU reads the instruction at address stored in the Program Counter (PC/RIP)
  PC = 0x00401000  →  reads bytes at that address from memory

Step 2: DECODE
  CPU figures out what instruction this is
  Bytes: 48 89 E5 → "mov rbp, rsp" (copy stack pointer to base pointer)

Step 3: EXECUTE
  CPU performs the operation
  → writes value of RSP register into RBP register

Step 4: PC UPDATE
  PC = PC + instruction_length  → points to the next instruction

Repeat 3,000,000,000 times per second (3 GHz CPU)
```

**Modern CPUs do this for multiple instructions at once (pipelining, out-of-order execution, superscalar), but the conceptual model is: fetch → decode → execute.**

**CPU cores:**
A modern CPU has multiple cores. Each core can independently execute instructions. A quad-core CPU can execute 4 instruction streams simultaneously. The OS scheduler assigns different processes to different cores.

**CPU modes (critical for OS):**

| Mode | x86 Name | What it can do |
|------|----------|----------------|
| Most privileged | Ring 0 (kernel mode) | Everything: all instructions, all memory |
| Least privileged | Ring 3 (user mode) | Limited: no hardware access, no privileged instructions |

When the computer boots, the CPU starts in Ring 0. The OS stays in Ring 0. User programs run in Ring 3. (Rings 1 and 2 existed in theory; modern OSes don't use them.)

---

## 3. Registers — The CPU's Working Memory

Registers are tiny, ultra-fast memory cells inside the CPU itself. There's no latency — they're literally part of the CPU's circuitry.

**x86-64 General Purpose Registers:**

| Register | 64-bit name | Traditional use |
|----------|-------------|----------------|
| Accumulator | RAX | Return values, arithmetic |
| Base | RBX | Base pointer (preserved across calls) |
| Counter | RCX | Loop counters, 4th argument |
| Data | RDX | I/O operations, 3rd argument |
| Source index | RSI | Source of memory copies, 2nd argument |
| Destination index | RDI | Destination of memory copies, 1st argument |
| Stack pointer | RSP | Points to top of stack (always!) |
| Base pointer | RBP | Points to stack frame base |
| R8–R15 | R8–R15 | Additional general purpose |

**Special registers (critical for OS):**

| Register | Purpose |
|----------|---------|
| **RIP (Instruction Pointer)** | Address of the NEXT instruction to execute |
| **RFLAGS** | Status bits: zero flag, sign flag, carry flag, overflow flag, etc. |
| **CR0** | Control register: enable protected mode, paging |
| **CR3** | Page directory base register (points to page table) |
| **CS, DS, SS** | Segment registers (selectors into GDT) |
| **GDTR, IDTR** | Points to Global Descriptor Table, Interrupt Descriptor Table |
| **MSR (Model Specific Registers)** | Hundreds of registers for advanced CPU features |

**Why registers matter for the OS:**
When the OS switches from one process to another (context switch), it must:
1. Save ALL of the current process's register values (RIP, RSP, general purpose, RFLAGS)
2. Load the saved register values of the next process
3. Jump to that process's saved RIP

The CPU immediately "becomes" the new process — because ALL CPU state is in the registers.

---

## 4. The Memory Hierarchy

Computers have multiple levels of memory, each faster but smaller (and more expensive):

```
         Faster, smaller, more expensive
                    ↑
         ┌──────────────────────┐
         │   CPU REGISTERS      │  ~16 bytes,   0 cycles (in CPU)
         ├──────────────────────┤
         │     L1 CACHE         │  32–64 KB,    4 cycles (in CPU)
         ├──────────────────────┤
         │     L2 CACHE         │  256 KB–1 MB, 12 cycles (in CPU)
         ├──────────────────────┤
         │     L3 CACHE         │  4–32 MB,     40 cycles (on chip)
         ├──────────────────────┤
         │   MAIN MEMORY (RAM)  │  4–64 GB,     100+ cycles (off chip)
         ├──────────────────────┤
         │     SSD / NVMe       │  256 GB–4 TB, ~100,000 cycles
         ├──────────────────────┤
         │     HDD              │  1–20 TB,     ~10,000,000 cycles
         └──────────────────────┘
                    ↓
         Slower, larger, cheaper
```

**What this means for the OS:**
- The OS tries to keep frequently used data in cache (automatically done by hardware)
- When the OS writes to a memory address, it goes into cache first, then eventually to RAM
- When the OS needs to store more data than fits in RAM, it uses the SSD as "swap space" — but it's 1000× slower than RAM!

**Cache coherency:** On multi-core CPUs, each core has its own L1/L2 cache. If Core 1 writes to address 0x1000 and Core 2 reads address 0x1000, they need to see the same value. The hardware cache coherency protocol (MESI) handles this — but it's one reason multi-core programming is complex.

---

## 5. RAM — Main Memory

RAM (Random Access Memory) is where the OS, all running programs, and their data live while the computer is on.

**Key properties:**
- **Byte-addressable:** Every single byte has a unique address (0, 1, 2, 3, ...)
- **Volatile:** Contents are lost when power is cut
- **Uniform access time:** Every address takes the same time to read (unlike disk)
- **Random access:** Can jump to any address instantly (unlike a tape)

**Physical vs. virtual addresses:**
- **Physical address:** The actual address on the RAM chips (0 to RAM_SIZE-1)
- **Virtual address:** The address a program sees — a fiction created by the OS and MMU

Modern programs never see physical addresses. They see virtual addresses. The MMU (Memory Management Unit, part of the CPU) translates virtual → physical.

```
Program thinks it's reading address: 0x00007fff_cafe1234
MMU translates:                       physical RAM address: 0x0000_0001_ab34_5678
CPU reads from:                       physical chip at the translated address
```

This is how:
- Process A and Process B can both think they have address 0x1000 — they map to different physical addresses
- A process can have more virtual address space than physical RAM (using swap)
- Security: process A can't access process B's memory (different mappings)

---

## 6. The Bus System — How Components Communicate

Components inside a computer communicate over **buses** — shared electrical pathways.

**Three types of signals on a bus:**

```
ADDRESS BUS:  Which location? (CPU says: "I want address 0x4000")
              ← one direction (CPU to everything else)

DATA BUS:     What is the data? (RAM says: "here's the bytes: 0x48 0x89 0xE5")
              ← two directions (read or write)

CONTROL BUS:  What operation? (CPU signals: "READ" or "WRITE")
              ← various signals (memory read, I/O write, interrupt acknowledge, etc.)
```

**Modern computer bus architecture:**

```
CPU ←→ [PCIe] ←→ GPU, NVMe SSD, fast NICs
CPU ←→ [DMI]  ←→ Chipset (PCH)
              ↓
         [USB, SATA, audio, slow I/O]
```

Modern CPUs have a direct memory controller — no separate "memory bus controller". RAM talks directly to the CPU via the memory controller built into the CPU die.

**I/O Ports vs. Memory-Mapped I/O:**
There are two ways to communicate with hardware:

1. **Port-mapped I/O (PMIO):** x86 has special `in` and `out` instructions that talk to a separate "I/O address space" (ports 0x0000–0xFFFF). Example: PC speaker is at port 0x61.

2. **Memory-mapped I/O (MMIO):** Device registers are placed at specific physical memory addresses. You read/write them like RAM. Example: VGA framebuffer is at physical address 0xB8000.

Most modern devices use MMIO. Legacy devices (keyboard controller, PIC) use PMIO.

---

## 7. The Memory Map — What's at Which Address

When a PC boots in legacy BIOS mode, the physical memory layout looks like this:

```
Physical Address (hex)    Contents
─────────────────────────────────────────────────────
0x00000000 – 0x000003FF   Interrupt Vector Table (IVT) — 1 KB
0x00000400 – 0x000004FF   BIOS Data Area (BDA)
0x00000500 – 0x00007BFF   Usable RAM (first 29 KB)
0x00007C00 – 0x00007DFF   Boot sector (512 bytes — loaded by BIOS here)
0x00007E00 – 0x0007FFFF   More usable RAM
0x00080000 – 0x0009FFFF   Mostly usable (extended BIOS data)
0x000A0000 – 0x000BFFFF   VGA framebuffer (write pixels here → screen changes)
0x000C0000 – 0x000FFFFF   BIOS ROM, hardware BIOS extensions
0x00100000 – ...          Extended memory (RAM above 1MB, for OS and programs)
```

**Key addresses every OS developer must know:**
- `0x7C00` — the BIOS loads the boot sector here
- `0xB8000` — write ASCII + color bytes here to display text on the screen (VGA text mode)
- `0x100000` (1MB) — where the kernel typically loads

When we build our OS, we'll use these exact addresses.

---

## 8. The Stack and the Heap

Every running process has its memory organized into regions:

```
High address ↑
┌───────────────────────┐
│   Stack               │ ← grows downward ↓
│   (function calls,    │   (local variables, return addresses)
│    local variables)   │
├───────────────────────┤
│   (free space)        │
├───────────────────────┤
│   Heap                │ ← grows upward ↑
│   (malloc/free)       │   (dynamically allocated data)
├───────────────────────┤
│   BSS                 │   (uninitialized global variables — zeroed)
├───────────────────────┤
│   Data segment        │   (initialized global/static variables)
├───────────────────────┤
│   Text segment        │   (compiled program code — read-only)
└───────────────────────┘
Low address ↓
```

**The stack explained:**
The stack is a Last-In-First-Out structure. RSP (stack pointer) always points to the "top" (lowest address — remember, it grows downward).

```c
void foo() {
    int x = 42;     // push 4 bytes onto stack, RSP -= 4
    bar();          // push return address onto stack, RSP -= 8
}
// when foo() returns: RSP += 8 + 4 (pop frame)
```

**The heap explained:**
When a program calls `malloc(100)`, the OS's memory manager finds 100 bytes of free heap space and returns a pointer. When `free()` is called, that space is marked available again. This is entirely managed by user-space code (C library) calling the OS to get big chunks (via `brk()` or `mmap()` system calls).

**The kernel has its own stack too.** Each kernel thread has a separate kernel stack — about 8–16KB — for when it's running kernel code.

---

## 9. I/O — Input and Output

I/O means any communication with devices that aren't the CPU or RAM: disk, keyboard, network, USB, GPU.

**The fundamental problem with I/O:**
The CPU runs at 3 GHz. A disk read takes 100 microseconds (NVMe) to 10 milliseconds (HDD). That's 300,000 to 30,000,000 CPU cycles of waiting.

**Three ways the OS can handle I/O:**

**1. Programmed I/O (polling):**
```
while (device_not_ready)
    check device status register;
read data register;
```
CPU wastes cycles spinning. Only acceptable for very fast devices or embedded systems.

**2. Interrupt-driven I/O:**
```
Start I/O operation;
CPU continues doing other work;
...later...
Hardware interrupt fires when I/O completes;
CPU switches to ISR → processes completed data;
CPU returns to previous work;
```
Most efficient for medium-speed I/O. Keyboard, mouse, serial ports use this.

**3. DMA (Direct Memory Access):**
```
Configure DMA: "read 4KB from disk sector X into memory at address Y"
CPU is free; DMA controller does the transfer directly to/from RAM
When done, DMA raises interrupt
CPU processes result
```
For bulk transfers (disk, network). CPU is 100% free while the transfer happens.

---

## 10. DMA — Direct Memory Access

DMA allows devices to transfer data directly to/from RAM without involving the CPU for each byte.

**Without DMA (CPU does the transfer):**
```
CPU: read 1 byte from disk → write to RAM[0]
CPU: read 1 byte from disk → write to RAM[1]
... repeat 4096 times for a 4KB block ...
CPU was busy for entire transfer
```

**With DMA:**
```
CPU tells DMA controller:
  "Transfer 4096 bytes from disk sector 1234 to RAM address 0x100000"
CPU continues running other processes
DMA controller directly moves data to RAM over the bus
When done, DMA signals CPU with interrupt
CPU spends ~1 microsecond handling the interrupt
Total CPU cost: 1 microsecond instead of 4096 microseconds
```

Modern disks, network cards, and GPUs all use DMA. The OS configures DMA transfers and then gets out of the way.

---

## 11. Interrupts — Hardware Events

**Interrupts** are how hardware tells the CPU "something happened."

Without interrupts, the CPU would need to constantly check every device — "Is a key pressed? Is the network card ready? Is the disk done?" — wasting enormous time.

With interrupts, the CPU runs normally until a device needs attention, then handles it immediately.

**The interrupt mechanism:**

```
Normal execution:
  CPU running process A at instruction 0x401234

Hardware event:
  User presses a key → keyboard controller raises interrupt line IRQ1

CPU responds:
  1. Finishes current instruction
  2. Saves current state (RIP = 0x401238, flags, etc.) to the stack
  3. Switches to kernel mode
  4. Looks up IRQ1 in the IDT (Interrupt Descriptor Table)
  5. Jumps to the keyboard ISR (Interrupt Service Routine)

OS kernel's ISR runs:
  Reads the scan code from keyboard port 0x60
  Converts to ASCII character
  Puts it in the keyboard buffer
  Sends End-Of-Interrupt (EOI) to PIC

CPU restores state:
  Loads saved RIP (0x401238), flags
  Returns to user mode
  Process A continues from exactly where it was
```

**Types of interrupts:**

| Type | What causes it | Example |
|------|---------------|---------|
| **Hardware IRQ** | External device | Keyboard press, timer tick, disk complete |
| **Software interrupt** | `int 0x80` instruction | System call (legacy Linux method) |
| **Exception** | CPU detected error | Division by zero, page fault, invalid instruction |

**The Timer Interrupt:**
The most important interrupt for OS scheduling. The programmable timer fires ~100–1000 times per second. The OS uses this tick to:
- Preempt a running process (take away the CPU)
- Update the scheduler
- Update the system clock
- Run deferred work

This is how multitasking works: every timer tick, the OS may switch to a different process.

---

## 12. BIOS, UEFI, and Firmware

**Firmware** is software stored in ROM (non-volatile chip) on the motherboard. It runs before any OS.

**BIOS (Basic Input/Output System) — the classic:**
- Has existed since the IBM PC (1981)
- Stored in a ROM chip on the motherboard
- Runs at power-on: tests hardware (POST — Power-On Self Test)
- Finds a boot device (disk, USB, network)
- Loads the first 512 bytes of the boot device into RAM at address 0x7C00
- Jumps to 0x7C00 and executes the bootloader

BIOS runs in 16-bit "real mode" — the ancient x86 mode from the 1980s.

**UEFI (Unified Extensible Firmware Interface) — modern:**
- Replaced BIOS starting around 2012
- Runs in 32-bit or 64-bit mode
- Much more capable: can read FAT32 partitions, show a GUI setup screen, load multiple boot loaders
- Secure Boot: cryptographically verifies the boot loader before running it
- Still provides hardware discovery, but through ACPI tables and DeviceTree

**Why it matters for OS development:**
When we build our OS in Volume 9, we'll start with BIOS/Multiboot2 (simpler for learning). Real production OSes like Linux support both BIOS and UEFI.

---

## Summary

| Component | What it does | OS relevance |
|-----------|-------------|-------------|
| CPU | Executes instructions; has privilege modes | OS runs in Ring 0; controls mode transitions |
| Registers | Ultra-fast CPU storage; track program state | Context switch saves/restores all registers |
| Cache | Fast L1/L2/L3 memory in CPU | OS tries to stay cache-friendly |
| RAM | Working memory for OS + processes | OS manages who owns which physical pages |
| Bus | Data highway between components | OS configures DMA, I/O over bus |
| Memory map | Fixed physical address layout | OS knows where VGA buffer, ROM, RAM are |
| Stack | LIFO memory for function calls | Each process/kernel thread has its own |
| Heap | Dynamic memory allocation | OS provides pages; allocator manages them |
| I/O | Communication with devices | OS provides drivers; polling/IRQ/DMA |
| DMA | Device ↔ RAM transfer without CPU | OS configures DMA; CPU handles interrupt |
| Interrupts | Hardware → CPU notification | Foundation of preemptive multitasking |
| BIOS/UEFI | Firmware at boot | Loads bootloader; provides initial services |
