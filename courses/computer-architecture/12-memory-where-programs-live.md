# Chapter 12: Memory — Where Programs Live

Imagine checking into a hotel with a billion rooms. Every room has a number on the door. You can walk up to any room, knock, and either drop something off or pick something up. That is RAM — your computer's main memory. The CPU is the guest who uses the rooms. The room numbers are addresses. The items left inside are data.

Without memory, a CPU is a calculator that forgets everything the moment it finishes a calculation. Memory is where programs wait to be run, where data waits to be processed, and where results wait to be stored. It is the stage on which all computation takes place.

This chapter goes deep into how memory actually works — from the electrons to the address space, from a single DRAM cell to the entire map of how your computer carves up gigabytes of address space.

---

## Table of Contents

1. [What Is RAM?](#1-what-is-ram)
2. [Every Byte Has an Address](#2-every-byte-has-an-address)
3. [32-Bit vs 64-Bit Address Spaces](#3-32-bit-vs-64-bit-address-spaces)
4. [The Memory Bus](#4-the-memory-bus)
5. [How the CPU Reads from Memory](#5-how-the-cpu-reads-from-memory)
6. [How the CPU Writes to Memory](#6-how-the-cpu-writes-to-memory)
7. [DRAM: Dynamic RAM and the Refresh Requirement](#7-dram-dynamic-ram-and-the-refresh-requirement)
8. [ROM: Memory That Does Not Forget](#8-rom-memory-that-does-not-forget)
9. [Flash Storage and SSDs](#9-flash-storage-and-ssds)
10. [The Memory Hierarchy Preview](#10-the-memory-hierarchy-preview)
11. [The Memory Map](#11-the-memory-map)
12. [Stack Memory](#12-stack-memory)
13. [Heap Memory](#13-heap-memory)
14. [Summary](#summary)
15. [Exercises](#exercises)

---

## 1. What Is RAM?

**RAM** stands for **Random Access Memory**. The word "random" here does not mean unpredictable — it means you can access any location in any order, in roughly the same amount of time. This is the opposite of something like a cassette tape, where you have to rewind or fast-forward to reach the part you want (called **sequential access**).

RAM is the computer's main working memory. When you open a program, the operating system copies it from your SSD or hard drive into RAM. When the CPU runs the program, it fetches instructions from RAM and reads and writes data in RAM. Everything currently "happening" on your computer is happening in RAM.

### The Hotel Analogy

Think of RAM as a giant hotel:

```
+-----------------------------------------------+
|               THE RAM HOTEL                   |
|                                               |
|  [Room 0] [Room 1] [Room 2] [Room 3] ...      |
|  [Room 4] [Room 5] [Room 6] [Room 7] ...      |
|  ...                                          |
|  [Room 4,294,967,295]  <-- last room (32-bit) |
+-----------------------------------------------+

Each room holds exactly 8 bits (1 byte) of data.
The room number is the address.
The CPU is the guest.
The CPU's registers are the guest's pockets.
```

Key properties of RAM:
- **Volatile**: When power is cut, all data is lost. The hotel empties when the lights go out.
- **Fast**: The CPU can access any location in nanoseconds (billionths of a second).
- **Random access**: Room 4,000,000 takes the same time to reach as Room 1.
- **Byte-addressable**: Every individual byte has its own unique address.

### Why RAM Is Not Enough on Its Own

RAM is fast but expensive and power-hungry. You cannot store your entire movie library in RAM — it would cost thousands of dollars and drain the battery in minutes. RAM is the expensive short-stay hotel. Your SSD is the cheaper long-term storage unit across town. The CPU needs both.

> **Quick Check 1.1:** What does "volatile" mean for RAM?
> **Quick Check 1.2:** What does "random access" mean? How is it different from sequential access?
> **Quick Check 1.3:** In the hotel analogy, what are the CPU's registers?

---

## 2. Every Byte Has an Address

The fundamental organizing principle of memory is simple: **every byte has a unique address**, expressed as an integer starting at 0.

```
Address  | Contents (8 bits)
---------+-----------------
0x0000   |  0100 1000   ('H')
0x0001   |  0110 0101   ('e')
0x0002   |  0110 1100   ('l')
0x0003   |  0110 1100   ('l')
0x0004   |  0110 1111   ('o')
0x0005   |  0000 0000   (null)
...
```

Addresses are almost always written in **hexadecimal** (base 16) because a 32-bit address would be an unwieldy 10-digit decimal number, but only 8 hex digits. The prefix `0x` signals hexadecimal notation.

### Multi-Byte Values: Endianness

Most data types are larger than one byte. An integer is 4 bytes; a 64-bit double is 8 bytes. When you store a multi-byte value in memory, you must decide: which byte goes in the lowest-numbered address?

There are two conventions, named after a famous debate in Gulliver's Travels about which end of an egg to crack open:

| Endianness | Meaning | Used By |
|------------|---------|---------|
| **Big-endian** | Most significant byte at lowest address | Network protocols, IBM mainframes |
| **Little-endian** | Least significant byte at lowest address | x86, ARM (usually), RISC-V |

Example: storing the 32-bit value `0x12345678` starting at address `0x100`:

```
Address   | Big-endian | Little-endian
----------+------------+--------------
0x100     |    0x12    |    0x78
0x101     |    0x34    |    0x56
0x102     |    0x56    |    0x34
0x103     |    0x78    |    0x12
```

Neither is more correct — it is a convention. But mixing them (as happens when you send data over a network from a little-endian machine to a big-endian one) causes bugs. This is why network protocols specify **network byte order** (big-endian) as a standard.

### Alignment

CPUs often prefer — or require — that multi-byte values be **aligned**: stored at addresses that are multiples of their size.

- A 4-byte integer should ideally start at an address divisible by 4 (0x100, 0x104, 0x108...).
- An 8-byte double should start at an address divisible by 8.

An **unaligned** access (e.g., a 4-byte integer starting at address `0x101`) may require two memory reads and extra shifting work, or may cause a hardware fault on strict architectures. Compilers handle alignment automatically in most cases.

> **Quick Check 2.1:** If address 0x200 contains `0x78`, address 0x201 contains `0x56`, 0x202 contains `0x34`, and 0x203 contains `0x12`, what is the 32-bit value stored here in little-endian format?
> **Quick Check 2.2:** Why are addresses written in hexadecimal rather than decimal?
> **Quick Check 2.3:** What is memory alignment and why does it matter?

---

## 3. 32-Bit vs 64-Bit Address Spaces

The size of an address — how many bits the CPU uses to represent it — directly determines how much memory the CPU can talk to. This is not a software limit; it is a physical limit baked into the hardware design.

### 32-Bit Address Space

A 32-bit CPU has 32 address lines. Each line is either 0 or 1. The number of unique addresses you can form with 32 bits is:

```
2^32 = 4,294,967,296 addresses
     = 4,294,967,296 bytes
     = 4 gigabytes (GB)
```

This was enormous in the 1990s. By the mid-2000s, it became a serious limitation. Computers started shipping with more than 4 GB of RAM, but a 32-bit CPU literally could not address it — like a hotel with only 10-digit room numbers trying to add an 11th digit to their system.

Workarounds like **PAE (Physical Address Extension)** allowed 32-bit CPUs to access more than 4 GB physically, but each individual process was still limited to a 4 GB virtual address space.

### 64-Bit Address Space

A 64-bit CPU uses 64 address lines:

```
2^64 = 18,446,744,073,709,551,616 addresses
     = 18 exabytes (EB)
     = 18 billion gigabytes
```

In practice, modern CPUs only implement 48 of the 64 address bits (for a 256 TB limit), because nobody makes that much RAM yet and implementing all 64 would add cost and complexity for no benefit today. But the architecture is ready when the time comes.

### Visualizing the Difference

```
32-BIT ADDRESS SPACE (4 GB)
|===========================================|
0                                    4 GB

64-BIT ADDRESS SPACE (theoretical 16 EB)
|================================================================...
0                                               16,000,000,000 GB
    (current implementations use ~256 TB of this space)
```

| Property | 32-bit | 64-bit |
|----------|--------|--------|
| Address size | 4 bytes | 8 bytes |
| Max addressable RAM | 4 GB | 16 EB (theoretical) |
| Practical RAM limit | 4 GB | 256 TB (common today) |
| Era of dominance | 1985-2005 | 2003-present |
| Pointer size | 32 bits | 64 bits |

The transition to 64-bit also doubled the size of **pointers** (variables that hold addresses) from 4 bytes to 8 bytes, which slightly increases memory usage in programs — but the tradeoff is completely worth it.

> **Quick Check 3.1:** Why can a 32-bit CPU not use more than 4 GB of RAM?
> **Quick Check 3.2:** How many bytes is 2^32?
> **Quick Check 3.3:** A 64-bit CPU could theoretically address 16 exabytes, but practical CPUs today only use 48 address bits. What is 2^48 bytes in more familiar units?

---

## 4. The Memory Bus

The CPU and RAM are separate chips on the motherboard. They communicate through a set of wires called the **memory bus**. The word "bus" in electronics means a shared set of wires — like a city bus route that multiple passengers share, rather than each having a private car.

The memory bus has three distinct groups of wires, each serving a different purpose:

```
                    MEMORY BUS
CPU                                         RAM
+-------+   Address Lines (32 or 64 wires)  +-------+
|       |=================================> |       |
|       |   Data Lines (64 wires typical)   |       |
|       |<================================> |       |
|       |   Control Lines (~5-10 wires)     |       |
|       |=================================> |       |
+-------+                                   +-------+
```

### Address Lines

These wires carry the address — which byte (or word) the CPU wants to read or write. On a 32-bit system there are 32 address lines; on a 64-bit system, up to 64 (though often fewer in practice). The address lines are **unidirectional**: they carry information from the CPU to RAM (the CPU tells RAM which room to look in, not the other way around).

### Data Lines

These wires carry the actual data being transferred. Modern memory buses are typically 64 bits (8 bytes) wide, meaning 8 bytes travel in parallel on every transfer. The data bus is **bidirectional**: during a read, data flows from RAM to CPU; during a write, data flows from CPU to RAM.

### Control Lines

These wires carry signals that coordinate the transaction:

| Control Signal | Meaning |
|---------------|---------|
| `MEM_READ` (or `OE` — Output Enable) | CPU wants to read from memory |
| `MEM_WRITE` (or `WE` — Write Enable) | CPU wants to write to memory |
| `CHIP_SELECT` | Activates a specific RAM chip |
| `RAS` (Row Address Strobe) | Selects which row of DRAM cells to activate |
| `CAS` (Column Address Strobe) | Selects which column within that row |
| `READY`/`ACK` | RAM signals that data is ready |

### Bus Width and Bandwidth

The data bus width determines how many bytes transfer per clock cycle. If the bus is 64 bits wide and runs at 3200 MHz (DDR4-3200), the theoretical bandwidth is:

```
Bandwidth = (bus width in bytes) x (transfers per second)
          = 8 bytes x 3,200,000,000
          = 25,600,000,000 bytes/second
          = ~25.6 GB/s
```

This is why "DDR4-3200" specifications matter — the number directly relates to available memory bandwidth.

> **Quick Check 4.1:** What are the three groups of wires in a memory bus?
> **Quick Check 4.2:** Which direction does data flow on the data lines during a memory read?
> **Quick Check 4.3:** If a data bus is 64 bits wide, how many bytes transfer per transaction?

---

## 5. How the CPU Reads from Memory

A memory read is a carefully choreographed dance between the CPU and the RAM chip. Every step must happen in the right order. Here is the full sequence:

### Step-by-Step Memory Read

```
CPU needs value at address 0x1000

Step 1: CPU places address on address lines
        Address bus: 0x1000
        Control:     MEM_READ = 1, MEM_WRITE = 0

Step 2: RAM receives address, activates row decoder
        The DRAM internally selects which row of cells
        contains address 0x1000

Step 3: RAM activates column decoder
        Within the selected row, the column containing
        the target byte(s) is selected

Step 4: RAM drives data onto data lines
        Data bus: 0x42 (the byte stored at 0x1000)

Step 5: CPU reads data lines and stores value
        The CPU's internal circuitry captures 0x42
        into a register

Step 6: Bus returns to idle state
        Control lines deasserted, bus freed for
        next transaction

TIME: ~50-100 nanoseconds total (DRAM latency)
```

Visually:

```
     CPU                  ADDRESS BUS                  RAM
      |                                                  |
      |------ 0x1000 ---------------------------------> |  (1) address sent
      |<------------------------------------- row sel -- |  (2) row decoded
      |<------------------------------------- col sel -- |  (3) col decoded
      |<---- 0x42 -------------------------------------- |  (4) data returned
      |                                                  |
      Register now holds 0x42
```

### The Latency Problem

Notice steps 2 and 3: the RAM must physically activate row and column decoders. This is slow compared to the CPU's clock speed. A modern CPU runs at 3-5 GHz (1 clock = 0.2-0.33 nanoseconds). A DRAM access takes 50-100 nanoseconds — roughly **150-300 clock cycles** of waiting.

This gap — the CPU's speed vs. memory's latency — is called the **memory wall**, and it is one of the most important architectural challenges in computing. Caches (which we preview in Section 10) exist precisely to bridge this gap.

### Memory Read in Assembly

When assembly code executes a load instruction, the CPU triggers exactly this sequence:

```asm
; x86 assembly: load value at memory address stored in RBX into RAX
MOV RAX, [RBX]

; RISC-V: load word from address in x1 into x2
LW x2, 0(x1)
```

Each of these instructions triggers one complete memory read transaction on the bus.

> **Quick Check 5.1:** Approximately how many CPU clock cycles does a single DRAM access take on a modern system?
> **Quick Check 5.2:** What is the memory wall?
> **Quick Check 5.3:** List the steps in a memory read in order.

---

## 6. How the CPU Writes to Memory

A memory write follows a similar sequence, but the direction of data reverses.

### Step-by-Step Memory Write

```
CPU wants to store value 0x99 at address 0x2000

Step 1: CPU places address on address lines
        Address bus: 0x2000

Step 2: CPU places data on data lines
        Data bus: 0x99

Step 3: CPU asserts write control signal
        Control: MEM_WRITE = 1, MEM_READ = 0

Step 4: RAM detects write signal, reads address and data lines
        DRAM internally selects row and column for address 0x2000

Step 5: RAM writes 0x99 into the selected cell(s)

Step 6: RAM signals completion (READY/ACK)

Step 7: CPU releases bus
        Control lines deasserted
```

Visually:

```
     CPU                  BUS                         RAM
      |                                                 |
      |------ 0x2000 (address) ----------------------> |
      |------ 0x99   (data)   -----------------------> |
      |------ WRITE (control) -----------------------> |
      |                                    write done - |
      |<----- ACK ------------------------------------ |
      |                                                 |
      Memory location 0x2000 now holds 0x99
```

### Read vs. Write: Comparison

| Aspect | Read | Write |
|--------|------|-------|
| Data direction | RAM to CPU | CPU to RAM |
| Control signal | MEM_READ | MEM_WRITE |
| Who drives data bus | RAM | CPU |
| Result | Value loaded into register | Memory cell updated |
| Typical latency | ~50-100 ns | ~50-100 ns |

Writes can sometimes be slightly faster in practice because modern CPUs use **write buffers** — the CPU drops the data into a buffer and continues executing, while the buffer drains to RAM in the background. From the CPU's perspective, the write "finishes" instantly. This is an optimization detail we will revisit in cache chapters.

> **Quick Check 6.1:** During a memory write, which side drives the data lines: the CPU or RAM?
> **Quick Check 6.2:** What is a write buffer and what problem does it solve?
> **Quick Check 6.3:** What are the control signals involved in a memory write?

---

## 7. DRAM: Dynamic RAM and the Refresh Requirement

We have been saying "RAM" as if it is a single thing. There are actually several types of RAM. The kind used in your computer's main memory is **DRAM — Dynamic Random Access Memory**.

### How a DRAM Cell Works

Each DRAM cell stores one bit using the simplest possible structure: a **capacitor** and a **transistor**.

```
ONE DRAM BIT CELL

    Word Line (row select)
         |
         |
    +----+----+
    |  NMOS   |   <-- transistor (the "gate")
    |Transistor|
    +----+----+
         |
    +----+----+
    |Capacitor |  <-- stores the charge (1=charged, 0=discharged)
    +----------+
         |
        GND
```

- **Bit stored as charge**: A charged capacitor = logic 1. A discharged capacitor = logic 0.
- **Transistor as gate**: When the word line (row select) is asserted, the transistor opens, connecting the capacitor to the bit line so its charge can be read or changed.

This design is remarkably compact — one transistor, one capacitor per bit. This is why DRAM can pack so many cells into a small chip. A modern 16 GB DRAM chip contains roughly **137 billion** of these cells.

### The Refresh Problem

Here is the catch: **capacitors leak charge**. A charged capacitor representing a `1` will slowly discharge toward `0` over time, even with no intentional read or write. Within a few milliseconds, all the 1s would turn to 0s and your data would be gone.

The solution is **refresh**: periodically re-reading every cell and writing the correct value back, before the charge leaks too far. Modern DRAM must be refreshed every **64 milliseconds** — all rows must be cycled through in that window.

```
REFRESH CYCLE (simplified)

Every 64ms:
  For each row (0 to 8191 in a typical chip):
    1. Assert row address
    2. Sense amplifiers detect charge level (1 or 0)
    3. Write the full-strength charge back
    4. Move to next row

This happens automatically inside the DRAM chip,
managed by a component called the Memory Controller.
```

The word "dynamic" in DRAM refers to this dynamic (ongoing, time-varying) requirement for refresh. It contrasts with **SRAM (Static RAM)**, which uses a different cell design (a flip-flop, which you learned about in Chapter 6) that holds its state without refreshing.

| Property | DRAM | SRAM |
|----------|------|------|
| Cell design | 1 transistor + 1 capacitor | 6 transistors |
| Requires refresh? | Yes (every ~64 ms) | No |
| Density | Very high (cheap per GB) | Lower (expensive per GB) |
| Speed | Slower (~50-100 ns) | Faster (~1-5 ns) |
| Power | Lower (per bit stored) | Higher |
| Used for | Main memory (RAM) | CPU caches (L1, L2, L3) |

SRAM's higher speed comes at the cost of 6x more transistors per bit — making it far too expensive to use for gigabytes of main memory. But it is perfect for the small, fast caches right next to the CPU.

### DDR: Double Data Rate

Modern systems use **DDR DRAM** (Double Data Rate), which transfers data on both the rising edge and the falling edge of the clock signal — effectively doubling bandwidth compared to older single-rate DRAM.

```
Clock:    _|--|_|--|_|--|_
           ^ ^  ^ ^  ^ ^
           | |  | |  | |
           transfers happen at EVERY edge (DDR)
           vs. only rising edge (old SDR)
```

DDR versions (DDR1 through DDR5) each roughly doubled bandwidth by increasing clock frequency, bus width, or prefetch buffer size. DDR5 (2020+) runs at speeds up to 6400 MT/s with a 64-bit bus, achieving peak bandwidths exceeding 50 GB/s per channel.

> **Quick Check 7.1:** What two components make up a single DRAM cell?
> **Quick Check 7.2:** Why does DRAM need to be refreshed? What happens if refresh stops?
> **Quick Check 7.3:** Why is SRAM used for CPU caches rather than DRAM?

---

## 8. ROM: Memory That Does Not Forget

**ROM** stands for **Read-Only Memory**. As the name says, you can read from it but (in its classic form) cannot write to it after manufacturing.

### The Boot Problem

When you press the power button, RAM is empty. The CPU needs instructions to execute immediately — but who puts them there? ROM solves this chicken-and-egg problem: it contains code that is permanently stored and ready to execute the instant power is applied.

```
POWER ON SEQUENCE

1. CPU resets, sets PC to a fixed address (e.g., 0xFFFFFFF0 in x86)
2. That address maps to ROM, not RAM
3. ROM contains the firmware (BIOS/UEFI)
4. Firmware runs:
   - Tests hardware (POST -- Power On Self Test)
   - Initializes RAM, keyboard, display
   - Finds boot device (SSD/USB)
   - Copies bootloader from storage into RAM
   - Jumps to bootloader in RAM
5. Bootloader loads the OS kernel
6. From here, everything runs from RAM
```

### Types of ROM

Over time, "ROM" evolved into a family of technologies with varying degrees of rewritability:

| Type | Full Name | Writable? | Erase Method | Use Case |
|------|-----------|-----------|--------------|---------|
| **Mask ROM** | — | Never | (not possible) | Fixed data in cheap devices |
| **PROM** | Programmable ROM | Once only | (not possible) | Factory programming |
| **EPROM** | Erasable PROM | Yes | UV light (15-30 min) | Development/prototyping |
| **EEPROM** | Electrically Erasable PROM | Yes | Electrical pulse (slow) | BIOS/UEFI chips, smart cards |
| **Flash** | Flash Memory | Yes | Electrical, block-by-block | SSDs, USB drives, phones |

Modern computers use **EEPROM** or **Flash** for their firmware (the BIOS or UEFI chip on the motherboard). This allows manufacturers to release firmware updates — you can "flash" your BIOS to fix bugs or add features. The firmware update process is literally reprogramming the ROM chip.

### Firmware Today: UEFI

The classic BIOS (Basic Input/Output System) has been replaced on most modern computers by **UEFI (Unified Extensible Firmware Interface)**. UEFI is stored in a flash chip, supports large disks (> 2 TB), provides a graphical interface, and enables features like Secure Boot. But the fundamental principle is the same: it is code in non-volatile memory that runs before the OS.

> **Quick Check 8.1:** Why does a computer need ROM at all — why not just use RAM for everything?
> **Quick Check 8.2:** What is the difference between EPROM and EEPROM?
> **Quick Check 8.3:** What does POST stand for and when does it run?

---

## 9. Flash Storage and SSDs

**Flash memory** deserves its own section because it has transformed computing. SSDs (Solid State Drives) use flash memory and have nearly replaced magnetic hard drives in consumer computers.

### How Flash Memory Stores Data

Flash uses a modified transistor called a **Floating Gate Transistor** (or, in modern NAND flash, a **Charge Trap Transistor**). The key innovation: it can trap electrons in an isolated region with no power needed to maintain the charge, unlike DRAM's leaky capacitors.

```
FLOATING GATE TRANSISTOR

    Control Gate
        |
    ----+----
   | Dielectric|   <-- insulator (electrons cannot escape)
   |           |
   |Floating   |   <-- trapped electrons = stored 1 or 0
   |  Gate     |
   |           |
   | Dielectric|   <-- another insulator
    ----+----
        |
   Source ---- Drain
```

Because the electrons are trapped by insulators on all sides, flash memory is **non-volatile** — it retains data without power, for 10+ years. This makes it ideal for storage.

### SSD vs. DRAM vs. HDD

| Property | DRAM (RAM) | Flash (SSD) | Magnetic HDD |
|----------|-----------|-------------|--------------|
| Volatile? | Yes | No | No |
| Read latency | ~100 ns | ~100 us | ~5-10 ms |
| Write latency | ~100 ns | ~100 us (reads old block first) | ~5-10 ms |
| Endurance | Unlimited | ~3,000-100,000 writes per cell | Very high |
| Density | Medium | High | Very high |
| Cost per GB | ~$10 | ~$0.10 | ~$0.02 |
| Moving parts | None | None | Yes (spinning platters) |

Flash is about 1,000x slower than DRAM but 50-100x faster than a hard drive, and it is non-volatile. It sits between them in the memory hierarchy.

### Why Flash Has Limited Write Endurance

Each write cycle to a flash cell slightly degrades the insulating layers. After a certain number of program-erase cycles (ranging from ~3,000 for consumer MLC NAND to ~100,000 for enterprise SLC NAND), the cell wears out and can no longer reliably store data.

SSD controllers manage this with **wear leveling** — spreading writes evenly across all cells rather than hammering the same cells repeatedly. This is why the controller inside an SSD is a sophisticated computer in its own right.

> **Quick Check 9.1:** What makes flash memory non-volatile, unlike DRAM?
> **Quick Check 9.2:** Approximately how much slower is flash compared to DRAM in terms of read latency?
> **Quick Check 9.3:** What is wear leveling and why is it necessary?

---

## 10. The Memory Hierarchy Preview

We have now seen five different types of memory with wildly different properties. Computer architects arrange them in a **hierarchy** based on speed, cost, and capacity — closer to the CPU means faster and more expensive, farther away means slower and cheaper.

```
THE MEMORY HIERARCHY

                    SPEED    SIZE     COST
  CPU               |
  +------------+    | fastest  smallest  most expensive
  | Registers  |    |  ~0.3ns   ~1 KB     ~$1000/GB
  +------------+    |
  +------------+    |
  |  L1 Cache  |    |  ~1ns     ~32-64 KB
  +------------+    |
  +------------+    |
  |  L2 Cache  |    |  ~5ns    ~256 KB - 1 MB
  +------------+    |
  +------------+    |
  |  L3 Cache  |    |  ~30ns   ~4-64 MB
  +------------+    |
  +------------+    |
  |    DRAM    |    |  ~100ns  ~8-128 GB   ~$10/GB
  |   (RAM)    |    |
  +------------+    |
  +------------+    |
  |    SSD     |    |  ~100us  ~256GB-8TB  ~$0.10/GB
  +------------+    |
  +------------+    | slowest  largest    cheapest
  |    HDD     |    ~10ms    ~1TB-20TB   ~$0.02/GB
  +------------+
```

The hierarchy exploits a principle called **locality of reference**:
- **Temporal locality**: If you used a piece of data recently, you will probably use it again soon.
- **Spatial locality**: If you used address X, you will probably soon use addresses near X.

Caches use these principles to automatically keep frequently-used data close to the CPU, so that the 150-cycle DRAM latency is hit only for "cold" data not already in cache. We will explore caches in detail in Chapters 20-22.

### The Analogy Extended

| Layer | Hotel Analogy | Computer Equivalent |
|-------|--------------|---------------------|
| Registers | Items in your pockets | Registers (~16 values) |
| L1 Cache | Items on your desk | 32-64 KB, ~1 ns |
| L2 Cache | Items in your room | 256 KB - 1 MB, ~5 ns |
| L3 Cache | Items at hotel reception | 4-64 MB, ~30 ns |
| DRAM | Items in hotel storage | 8-128 GB, ~100 ns |
| SSD | Items in a nearby warehouse | TBs, ~100 us |
| HDD | Items in a distant archive | TBs, ~10 ms |

> **Quick Check 10.1:** What is the memory hierarchy and why does it exist?
> **Quick Check 10.2:** What is temporal locality? Give an example.
> **Quick Check 10.3:** Why do caches work — what principle makes them effective?

---

## 11. The Memory Map

The CPU has an address space — a range of addresses from 0 to some maximum. Not all of this space is connected to RAM. The address space is divided up — **mapped** — to different hardware and purposes. This division is called the **memory map**.

Think of the address space as a city's postal grid. Most addresses are houses (RAM). But some addresses are government buildings (ROM/firmware), some are warehouses (I/O devices), some are parks (unmapped/reserved). The postal system (CPU's memory controller) knows what is at each address.

### A Typical 32-bit Memory Map

```
32-BIT ADDRESS SPACE (4 GB total)

0xFFFFFFFF  +----------------------------------+
            |  BIOS/UEFI ROM (firmware)        | ~2MB
0xFFFE0000  +----------------------------------+
            |  ACPI tables, interrupt vectors  |
0xFFF00000  +----------------------------------+
            |  Memory-mapped I/O               | ~256MB
            |  (PCI, PCIe device registers)    |
0xC0000000  +----------------------------------+ <- "3GB hole"
            |  OS kernel + kernel data         |
0x80000000  +----------------------------------+
            |                                  |
            |  User process address space      | ~2-3 GB
            |  (code, data, heap, stack)       |
            |                                  |
0x00000000  +----------------------------------+
```

### Memory-Mapped I/O

Many devices (video cards, network cards, sound cards) have registers that the CPU reads and writes to control them. Rather than having a separate set of I/O instructions, these devices are given addresses in the regular address space. When the CPU writes to address `0xFD000000`, it might actually be writing a pixel color to the GPU — the memory controller routes this address not to RAM but to the GPU over the PCIe bus.

This clever trick — **memory-mapped I/O (MMIO)** — means device drivers can use ordinary load/store instructions to talk to hardware.

### Virtual vs. Physical Addresses

There is one more critical layer: modern CPUs do not use physical memory addresses directly. Each running process sees its own **virtual address space** — a private 4 GB (32-bit) or 256 TB (64-bit) address space. The CPU's **Memory Management Unit (MMU)** translates virtual addresses to physical addresses in real time, using a data structure called a **page table**.

```
PROCESS A's view:          PHYSICAL MEMORY:
Virtual 0x1000 ---------> Physical 0x45678000 (actual RAM chip location)
Virtual 0x2000 ---------> Physical 0x89ABC000

PROCESS B's view:
Virtual 0x1000 ---------> Physical 0x12340000 (different physical location!)
Virtual 0x2000 ---------> Physical 0xFEDCB000
```

This means two processes can both "think" they own address `0x1000`, and they are both right — they map to different physical locations. This is how modern operating systems safely isolate processes from each other. It is also why a crashed program does not corrupt another program's memory. We will cover virtual memory deeply in Chapters 23-25.

> **Quick Check 11.1:** What is a memory map?
> **Quick Check 11.2:** What is memory-mapped I/O? Give an example of a device that uses it.
> **Quick Check 11.3:** What is the difference between a virtual address and a physical address?

---

## 12. Stack Memory

Every running program has two important regions of memory automatically managed for it: the **stack** and the **heap**. They serve different purposes and grow in opposite directions.

### What Is the Stack?

The **stack** is a region of memory used to manage function calls and local variables. It operates on a LIFO (Last In, First Out) principle — exactly like a stack of plates. The most recently added item is always the first to be removed.

A special register called the **Stack Pointer (SP)** always points to the top of the stack.

### Stack Grows Downward

By convention on virtually all architectures, the stack starts at a high address and grows toward lower addresses as you push data onto it.

```
HIGH ADDRESS
+-----------------------------+ <- initial stack pointer (e.g., 0xFFFF0000)
|                             |
|   main()'s stack frame      |
|   - return address          |
|   - saved registers         |
|   - local variables (a, b)  |
|                             |
+-----------------------------+ <- SP after main() called foo()
|                             |
|   foo()'s stack frame       |
|   - return address          |
|   - saved registers         |
|   - local variables (x, y)  |
|                             |
+-----------------------------+ <- SP after foo() called bar()
|                             |
|   bar()'s stack frame       |
|   ...                       |
|                             |
+-----------------------------+ <- current SP (top of stack)
|  (free space)               |
|                             |
|                             |
LOW ADDRESS
```

### Stack Frames

Each function call creates a **stack frame** (also called an **activation record**) on the stack. The frame contains:

1. The **return address** — where to jump back when the function returns
2. The **saved frame pointer** — the previous function's frame base
3. **Function arguments** (if not passed in registers)
4. **Local variables**
5. **Saved registers** (registers the function must preserve)

```asm
; x86-64 function prologue (setup)
PUSH RBP           ; save old frame pointer
MOV  RBP, RSP      ; set frame pointer to current stack top
SUB  RSP, 32       ; allocate 32 bytes for local variables

; ... function body ...

; x86-64 function epilogue (teardown)
MOV  RSP, RBP      ; restore stack pointer (frees locals)
POP  RBP           ; restore old frame pointer
RET                ; pop return address, jump to it
```

### Stack Overflow

The stack has a limited size (typically 1-8 MB for a thread). If a program uses too much stack space — usually through infinite or deeply nested recursive function calls — it runs off the bottom of the stack region into unmapped memory. The OS detects this and kills the program with a **stack overflow** error. Yes, that is exactly where the famous programming website gets its name.

```
Normal call stack (ok):
  main -> processFile -> readLine -> parseToken
  (4 frames deep, using maybe 1 KB of stack)

Infinite recursion (BAD):
  countdown(1000) -> countdown(999) -> countdown(998) -> ...
  -> countdown(0) -> countdown(-1) -> ...
  (eventually uses all stack space, CRASH!)
```

> **Quick Check 12.1:** In what direction does the stack grow — toward higher or lower addresses?
> **Quick Check 12.2:** What is a stack frame and what does it contain?
> **Quick Check 12.3:** What causes a stack overflow?

---

## 13. Heap Memory

The **heap** is the other major memory region in a process. Unlike the stack, which is automatically managed by function call/return mechanics, the heap is **manually managed** (or managed by a garbage collector in languages like Java or Python).

### What Is the Heap For?

The stack is great for local variables — but local variables disappear when the function returns. What if you need data that:
- Outlives the function that created it?
- Has a size determined at runtime (not compile time)?
- Is simply too large for the stack?

This is what the heap is for. When you call `malloc()` in C, `new` in C++/Java, or create a `list` in Python, the runtime allocates memory from the heap.

### Heap Grows Upward

The heap starts at a low address and grows toward higher addresses, growing up toward the stack (which is growing down).

```
PROCESS MEMORY LAYOUT (simplified)

HIGH ADDRESS  +----------------------+
0xFFFF0000    |  Stack               | v grows down
              |  ...                 |
              |  (free space)        |
              |  ...                 |
0x10000000    |  Heap                | ^ grows up
              +----------------------+
0x00800000    |  BSS (uninit globals)|
              +----------------------+
0x00600000    |  Data (init globals) |
              +----------------------+
0x00400000    |  Code (text segment) |
LOW ADDRESS   +----------------------+
```

If the heap and stack ever grow into each other, the program crashes — a condition called a **heap-stack collision**, which is rare on 64-bit systems with vast address spaces but was a real concern on older 32-bit systems.

### Dynamic Allocation

When a program requests heap memory:

```
MALLOC REQUEST:  "I need 100 bytes"

  Heap Manager checks free list:
  +------------+  +------------+  +------------+
  | free: 64B  |->| free: 200B |->| free: 32B  |
  +------------+  +------------+  +------------+
           Found 200B block -- split it:
  +------------+  +------------+  +------------+
  | free: 64B  |->| USED: 100B |  | free: 100B |
  +------------+  +------------+  +------------+
  Returns pointer to the 100B block to the program.
```

The program must later **free** this memory (in C/C++) or the runtime's garbage collector must reclaim it. Failing to free heap memory is a **memory leak** — the heap grows and grows until the program crashes or the system runs out of memory.

### Heap Fragmentation

Over time, as allocations and frees happen in different orders, the heap can become **fragmented** — lots of small free blocks that cannot satisfy a large allocation request, even if total free space is sufficient.

```
FRAGMENTED HEAP

[used 50B][FREE 10B][used 100B][FREE 15B][used 30B][FREE 8B][used 200B]

Total free: 10+15+8 = 33 bytes
But you cannot allocate a single 25-byte block from this!
(No contiguous free block large enough exists)
```

Memory allocators use sophisticated algorithms (first-fit, best-fit, buddy system) and **compaction** strategies to fight fragmentation.

### Stack vs. Heap: Summary

| Property | Stack | Heap |
|----------|-------|------|
| Management | Automatic (push/pop) | Manual or GC |
| Size | Fixed, small (1-8 MB) | Limited by RAM |
| Speed | Very fast (just increment SP) | Slower (find free block) |
| Lifetime | Tied to function scope | Until explicitly freed |
| Grows toward | Lower addresses (down) | Higher addresses (up) |
| Overflow | Stack overflow crash | Memory leak (gradual) |
| Use case | Local variables, call frames | Dynamic objects, large data |

> **Quick Check 13.1:** What is the heap used for? Give two scenarios where you would use heap instead of stack.
> **Quick Check 13.2:** What is a memory leak?
> **Quick Check 13.3:** What is heap fragmentation?

---

## Summary

Memory is the stage on which all computation takes place. Here is what we covered:

**RAM (DRAM)** is the main memory — fast, volatile, byte-addressable. Every byte has a unique address. A 32-bit address space tops out at 4 GB; a 64-bit address space extends to 16 exabytes (practically ~256 TB today).

**The memory bus** connects the CPU to RAM via three groups of wires: address lines (which byte?), data lines (what value?), and control lines (read or write?). The bus width determines bandwidth; latency is the big bottleneck.

**Reading from memory** involves the CPU placing an address on the bus, RAM decoding that address and driving data onto the data bus, and the CPU capturing that value into a register. The whole process takes ~50-100 ns — 150-300 CPU clock cycles.

**Writing to memory** is the reverse: the CPU drives both the address and data onto the bus while asserting the write control signal, and RAM captures and stores the value.

**DRAM cells** store bits as charge on capacitors. Because capacitors leak, DRAM requires periodic refresh every 64 ms. SRAM uses flip-flops instead — faster but larger and more expensive, making it suitable for caches but not main memory.

**ROM** provides non-volatile storage for firmware (BIOS/UEFI) that must survive power loss and execute before RAM is initialized. Modern systems use Flash-based ROM that can be updated electrically.

**Flash storage** uses floating-gate transistors to trap electrons semi-permanently, enabling non-volatile storage that is ~1,000x slower than DRAM but 50-100x faster than HDDs.

**The memory hierarchy** (registers -> L1 -> L2 -> L3 -> DRAM -> SSD -> HDD) exploits the tradeoff between speed and cost. Caches exploit locality of reference to make most accesses fast.

**The memory map** divides the address space between user code, data, OS, ROM, and memory-mapped I/O. Virtual addresses (what programs see) are translated to physical addresses by the MMU.

**Stack memory** handles function calls automatically — each call pushes a frame (return address, locals, saved registers), each return pops it. The stack grows downward from a high address.

**Heap memory** is for dynamic allocation — objects whose size or lifetime is not known at compile time. It grows upward, must be explicitly freed (or garbage collected), and is prone to fragmentation and leaks.

---

## Exercises

### Easy

1. **Address Space Math**: A CPU has 24 address lines. How many bytes of memory can it address? Convert your answer to megabytes.

2. **Memory Bus Roles**: Label each of the following as "address line", "data line", or "control line":
   - Carries the value being read from memory
   - Tells RAM whether the CPU wants to read or write
   - Carries which memory location the CPU wants

3. **Endianness**: The 32-bit hexadecimal value `0xDEADBEEF` is stored in little-endian format starting at address 0x1000. Write out the contents of addresses 0x1000, 0x1001, 0x1002, and 0x1003.

### Medium

4. **DRAM Refresh Calculation**: A DRAM chip has 8,192 rows and must be fully refreshed every 64 milliseconds. How many row refreshes per second does this require? What is the average time between refreshes for a single row?

5. **Memory Hierarchy Latency**: A program accesses data that is:
   - 90% of the time in L1 cache (1 ns)
   - 8% of the time in L2 cache (5 ns)
   - 1.9% of the time in L3 cache (30 ns)
   - 0.1% of the time in DRAM (100 ns)

   Calculate the average memory access time (AMAT). How much faster is this than if there were no cache at all?

6. **Stack Trace**: A program calls functions in this order: `main()` calls `processData()` which calls `sortArray()` which calls `swap()`. Draw the stack layout at the moment `swap()` is executing. Label each frame, the stack pointer, and the frame pointer.

### Hard

7. **Memory Bandwidth Analysis**: A DDR4-3200 system has a 64-bit (8-byte) data bus running at 3200 MT/s (megatransfers per second, with DDR meaning 2 transfers per clock). A program must copy 1 GB of data from one RAM location to another — this requires both reading 1 GB and writing 1 GB. Assuming the memory bus is the bottleneck and runs at 100% efficiency (unrealistic but useful for calculation), how long would this copy take? Why would the real time be higher?

8. **Stack Overflow Analysis**: A recursive function uses 64 bytes of stack space per call (return address, saved registers, local variables). A typical thread has an 8 MB stack. How many recursive calls can be made before a stack overflow? If each call does a constant amount of work and takes 10 ns, what is the maximum depth of recursion before overflow, and how long would that recursion chain take to execute?

9. **Virtual Memory Thought Experiment**: Two processes, A and B, are both running. Process A has a pointer at virtual address `0x00401234`. Process B also has a pointer at virtual address `0x00401234`. Can these pointers point to the same physical RAM location? Can they point to different physical RAM locations? Under what circumstances would each case occur? Why is this property (called **process isolation**) critical for operating system security?
