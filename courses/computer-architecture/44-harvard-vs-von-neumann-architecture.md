# Chapter 44: Harvard vs Von Neumann Architecture

Every computer is built on one of two fundamental organizational blueprints — or a hybrid of both. The Von Neumann architecture stores programs and data in the same memory, accessed through the same bus. The Harvard architecture uses separate memories and buses for programs and data. This distinction sounds simple, but it has profound implications for performance, security, and chip design. Modern computers live in between: their outer shell is Von Neumann (unified RAM), but their inner CPU pipelines are Harvard (split instruction and data caches). Understanding this distinction helps explain why CPUs work the way they do — and why some security vulnerabilities are fundamental to the Von Neumann model.

## Table of Contents

1. [Von Neumann Architecture — The Stored-Program Computer](#1-von-neumann-architecture--the-stored-program-computer)
2. [The Von Neumann Bottleneck](#2-the-von-neumann-bottleneck)
3. [Harvard Architecture — Separation of Concerns](#3-harvard-architecture--separation-of-concerns)
4. [Modified Harvard Architecture — The Best of Both](#4-modified-harvard-architecture--the-best-of-both)
5. [Security Implications: Von Neumann and Code Injection](#5-security-implications-von-neumann-and-code-injection)
6. [Where Each Architecture Appears Today](#6-where-each-architecture-appears-today)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Von Neumann Architecture — The Stored-Program Computer

Before 1945, most computers were **fixed-program machines** — you changed the computation by physically rewiring the machine. ENIAC (1945) required two to three days of rewiring to change programs.

**John von Neumann's key insight** (1945): Store the program itself in memory, just like data. The CPU reads instructions from memory, executes them, and can even modify the program (write to the same memory that holds instructions). This is the **stored-program model**.

```
Von Neumann Architecture:
  ┌──────────────────────────────────────┐
  │             CPU                       │
  │  ┌──────────┐    ┌─────────────┐    │
  │  │  ALU     │    │   Control   │    │
  │  │          │    │   Unit      │    │
  │  └──────────┘    └──────┬──────┘    │
  │                         │           │
  └─────────────────────────┼───────────┘
                            │ Single Bus
  ┌─────────────────────────┴───────────┐
  │          Memory (Unified)            │
  │  [ instructions ] [ data ]           │
  │  (both in the same memory space)     │
  └──────────────────────────────────────┘
```

The five components of Von Neumann architecture:
1. **Central Processing Unit (CPU)**: executes instructions
2. **Control Unit**: orchestrates execution
3. **ALU (Arithmetic Logic Unit)**: performs computation
4. **Memory**: stores both program and data
5. **I/O devices**: input/output

**Why this was revolutionary**: Programs became data. You could write a program that generated other programs (compilers!). You could load different programs from storage without rewiring. Every modern computer — from your phone to the world's fastest supercomputer — descends from this model.

### Quick Check
> 1. What was the main innovation of the Von Neumann stored-program model vs. fixed-program machines?
> 2. Draw a simple diagram of Von Neumann architecture showing the single bus connecting CPU to memory.
> 3. How does the stored-program model enable compilers?

---

## 2. The Von Neumann Bottleneck

The single bus connecting CPU to unified memory is both the genius and the limitation of the Von Neumann model. Every instruction fetch AND every data read/write must use the same bus. This creates a fundamental performance ceiling.

```
Von Neumann bottleneck:
  During instruction fetch:
    CPU reads instruction from memory → bus is busy
    Cannot also read data from memory simultaneously
  
  During data access:
    CPU reads/writes data from memory → bus is busy
    Next instruction cannot be fetched simultaneously
  
  Serial nature: fetch instruction, then fetch data, then fetch next instruction, ...
  One bus = one transaction at a time
```

John Backus (inventor of FORTRAN) named this the **Von Neumann bottleneck** in his famous 1977 Turing Award lecture. He described it as "a bottleneck between the CPU and memory...forcing information to flow through it one word at a time."

**The bottleneck in numbers**: At 3 GHz, the CPU wants to execute 3 billion instructions per second. At 64-bit bus width and 3.2 GHz DDR5: ~25 GB/s bandwidth. An instruction takes 4 bytes. A 3 GIPS CPU would need to fetch 12 GB/s of instructions alone — before any data reads/writes.

Solutions explored over 70 years:
1. **Caches**: Keep recently used data/instructions close to the CPU, reducing bus pressure
2. **Wider buses**: Multiple bytes per transfer (64-bit bus, dual-channel)
3. **Pipelining**: Overlap instruction fetch with prior instruction's execution
4. **Harvard architecture**: Separate buses for instructions and data (the main topic)

### Quick Check
> 1. What is the Von Neumann bottleneck?
> 2. At 3 GIPS and 4 bytes/instruction, what instruction fetch bandwidth is needed? Can DDR5 dual-channel support this?
> 3. Name three architectural solutions to the Von Neumann bottleneck.

---

## 3. Harvard Architecture — Separation of Concerns

The **Harvard architecture** uses physically separate memory storage and signal pathways for instructions and data. It was named after the Harvard Mark I electromechanical computer (1944), which stored programs on punched tape and data on registers/counters.

```
Harvard Architecture:
  ┌────────────────────────────────────────────────┐
  │                    CPU                          │
  │   ┌───────────┐         ┌──────────────────┐  │
  │   │  ALU      │         │  Control Unit    │  │
  │   └─────┬─────┘         └───────┬──────────┘  │
  │         │ Data Bus              │ Instruction Bus
  └─────────┼───────────────────────┼──────────────┘
            │                       │
  ┌─────────┴──────────┐  ┌─────────┴──────────┐
  │   Data Memory      │  │  Program Memory     │
  │ (read/write)       │  │ (read-only typically)│
  └────────────────────┘  └─────────────────────┘
```

**Key advantage**: The CPU can simultaneously fetch an instruction from program memory AND read/write data from data memory — **doubling potential throughput** over Von Neumann.

```
Harvard throughput:
  Cycle N:   Fetch instruction[N] from instruction memory
             Read data[N-1] from data memory    ← simultaneously
  
  Two buses operating in parallel → no bus contention
```

**Key constraints**:
1. **Program memory is often read-only** (ROM in microcontrollers): You can't modify the running program
2. **Separate address spaces**: A pointer to a data address is meaningless in the instruction address space
3. **Code cannot be generated at runtime**: no JIT compilation, no self-modifying code
4. **Higher hardware cost**: Two complete memory systems

Harvard architecture is the dominant design for **microcontrollers** (Chapter 46): program lives in Flash (read-only), data lives in SRAM. The MCU reads instructions from Flash and reads/writes data from SRAM simultaneously.

```
AVR microcontroller (Harvard):
  Program Memory: 32KB Flash (read-only, addressed by program counter)
  Data Memory:    2KB SRAM + I/O registers + external (read/write)
  
  Separate address spaces:
    LDS R16, 0x0100  ← loads from data address 0x0100 (SRAM)
    "instruction at 0x0100" is meaningless — different address space
```

### Quick Check
> 1. What is the key performance advantage of Harvard over Von Neumann?
> 2. Why is program memory often read-only in Harvard microcontrollers?
> 3. Name two things you cannot do in a strict Harvard architecture that you can in Von Neumann.

---

## 4. Modified Harvard Architecture — The Best of Both

Real modern processors don't use either pure model — they use the **Modified Harvard Architecture**: separate instruction and data caches internally, but a unified external memory (DRAM) externally.

```
Modified Harvard Architecture (modern CPU):
  
  ┌──────────────────────────────────────────────────────────┐
  │                      CPU CORE                             │
  │                                                           │
  │   ┌────────────────┐    ┌────────────────────────────┐  │
  │   │   Instruction   │    │      Data Cache (L1-D)     │  │
  │   │   Cache (L1-I)  │    │      (read/write)          │  │
  │   │   (read-only)   │    └────────────────────────────┘  │
  │   └────────────────┘                                      │
  │           │                           │                   │
  │           └──────────┬────────────────┘                  │
  │                      ↓                                    │
  │              ┌───────────────┐                            │
  │              │  L2 Cache     │ (unified)                  │
  │              └───────┬───────┘                            │
  │                      ↓                                    │
  │              ┌───────────────┐                            │
  │              │  L3 Cache     │ (unified, shared)          │
  │              └───────┬───────┘                            │
  └──────────────────────┼────────────────────────────────────┘
                         ↓
  ┌──────────────────────────────────────────────────────────┐
  │    DRAM (unified — instructions and data share memory)    │
  └──────────────────────────────────────────────────────────┘
```

**Inside the CPU**: Harvard (split L1 I-cache and D-cache). Eliminates structural hazard — the fetch unit reads from L1-I while the load/store unit reads/writes from L1-D simultaneously.

**Outside the CPU**: Von Neumann (unified DRAM, same address space for instructions and data). This allows:
- JIT compilation (generate code, write to memory, execute it)
- Self-modifying code (though rarely used)
- Dynamic linking and loading
- Unified virtual address space (no separate "instruction pointer space")

**Cache coherence complication**: When the CPU writes new instructions to DRAM (JIT compilation, dynamic linker), those instructions might be cached in L1-D (as data). The L1-I cache has the old version. The CPU must flush the instruction cache when new code is written. ARM requires explicit `isync` / `ic ivau` instructions for this. x86's strongly ordered memory model handles it automatically.

### Quick Check
> 1. In the Modified Harvard Architecture, what is Harvard about it and what is Von Neumann?
> 2. Why must the OS/JIT compiler flush the instruction cache after writing new code to memory?
> 3. What structural hazard does the split L1-I / L1-D cache eliminate?

---

## 5. Security Implications: Von Neumann and Code Injection

The Von Neumann model's "code and data in the same memory" is also a security liability. If an attacker can write data to memory and trick the CPU into executing it as code, they can run arbitrary code — this is the **code injection** attack.

**Buffer overflow → code injection (classic)**:
```c
void vulnerable(char* input) {
    char buf[128];
    strcpy(buf, input);    // no bounds check!
}

// If input is > 128 bytes, it overwrites beyond buf:
// [ 128-byte buf ][ saved return address ][ ... ]
//                  ← attacker writes shellcode here + overwrite return address
```

When the function returns, the program jumps to the attacker's shellcode (executable code the attacker injected as "data"). This worked because Von Neumann: both the stack (data) and code share the same address space and execution privileges.

**Defenses**:
1. **NX/XD bit (No-Execute/Execute Disable)**: Page table entries with the XD bit mark pages as "data-only, not executable." The CPU will fault if the instruction pointer reaches an NX page. Supported by Intel since 2004 (AMD since 2000). Implemented in the MMU.

2. **ASLR (Address Space Layout Randomization)**: Randomize where code, stack, and heap are loaded in virtual memory. The attacker doesn't know the address to jump to.

3. **Stack Canaries**: The compiler inserts a random "canary value" between local variables and the return address. Before returning, check if the canary was overwritten. If corrupted → abort.

4. **CFI (Control Flow Integrity)**: Compiler and hardware ensure control flow (calls, returns, indirect jumps) only targets valid, expected destinations. Prevents ROP (Return-Oriented Programming) attacks that chain existing code fragments.

**The fundamental tension**: The Von Neumann model allows JIT compilers, dynamic linkers, and self-modifying code — powerful capabilities. These require pages that are writable AND executable (W+X). Security requires W^X (writable XOR executable — never both). Modern OSes enforce W^X strictly; JIT compilers use `mprotect()` to temporarily change page permissions.

### Quick Check
> 1. How does the Von Neumann model enable code injection attacks?
> 2. What does the NX/XD bit do in hardware?
> 3. Explain W^X as a security principle. What does a JIT compiler need to do to comply?

---

## 6. Where Each Architecture Appears Today

```
Pure Von Neumann:
  - Von Neumann's original 1945 design
  - Historical mainframes
  - Some early microprocessors (Intel 8080)
  
Modified Harvard (mainstream modern CPUs):
  - Intel Core, AMD Ryzen: split L1-I/L1-D, unified L2/L3/DRAM
  - Apple M-series, Qualcomm Snapdragon
  - Any CPU with separate instruction and data L1 caches
  
Pure Harvard (microcontrollers/embedded):
  - AVR (Arduino): Harvard, separate Flash/SRAM address spaces
  - PIC microcontrollers: Harvard
  - ARM Cortex-M (some variants): Modified Harvard
  - DSPs: many use Harvard for predictable timing
  
Super-Harvard (DSPs):
  - Some DSPs add a THIRD memory system for coefficient tables
  - TI C6000 series: 3 separate memory buses
    (instruction memory, data memory, coefficient memory)
```

The pattern: as you move from high-performance general-purpose computing → real-time embedded → digital signal processing, the architecture moves from unified Von Neumann → modified Harvard → strict Harvard → super-Harvard.

### Quick Check
> 1. What architecture does an Arduino (AVR) use?
> 2. What architecture does an Intel Core i9 use internally?
> 3. What is a "super-Harvard" architecture and where is it used?

---

## Summary

- **Von Neumann architecture** (1945): unified memory for both instructions and data, single bus. Enables the stored-program model.
- The **Von Neumann bottleneck**: one bus limits simultaneous instruction fetch and data access.
- **Harvard architecture**: separate memories and buses for instructions and data. Better throughput, but program memory is typically read-only — no JIT, no self-modifying code.
- **Modified Harvard** (modern CPUs): split L1 instruction and data caches (Harvard inside), but unified DRAM (Von Neumann outside). Best of both worlds.
- **Security**: Von Neumann's unified address space enables code injection attacks. NX/XD bits, ASLR, W^X enforcement, and CFI are hardware/OS countermeasures.
- **Pure Harvard** survives in microcontrollers (AVR, PIC) and DSPs where program-in-Flash, data-in-SRAM is natural and performance is predictable.

---

## Exercises

### Easy
1. What was the revolutionary insight of the Von Neumann stored-program model?
2. Why does having separate instruction and data caches (Modified Harvard) eliminate a structural hazard in the pipeline?
3. Explain in one paragraph why a buffer overflow can lead to code execution on a Von Neumann machine.

### Medium
4. An AVR microcontroller (Harvard) has 32KB Flash (program memory) and 2KB SRAM (data memory). A C program defines a global constant string `"Hello"`. In Von Neumann (e.g., x86): where does this string live in memory? In Harvard (AVR): where does it live at compile time? How does the AVR C runtime copy it to SRAM at startup (and why)?
5. A JIT compiler (e.g., the JVM's JIT) generates native code at runtime. Describe the sequence of OS calls needed to: (a) allocate memory for the JIT-generated code, (b) write the machine code bytes to that memory, (c) make the memory executable (but not writable), (d) jump to the generated code. Why does step (c) require a memory protection change?
6. The NX bit prevents stack/heap execution. An attacker knows this and instead uses ROP (Return-Oriented Programming): they overwrite the return address to point to an existing code "gadget" (a short sequence ending in `ret`), which chains to another gadget, and so on. (a) Why does NX not stop ROP? (b) What defense specifically counters ROP? (c) How does ASLR make ROP harder but not impossible?

### Hard
7. Cache coherence in Modified Harvard: when a JIT compiler writes machine code to memory (data write → L1-D), the instruction cache may still hold old instructions from that address. Describe: (a) the hardware-level inconsistency this creates, (b) how ARM requires explicit cache maintenance (`DC CVAU` + `IC IVAU` + `DSB` + `ISB` sequence), (c) how x86 handles this implicitly via strong memory ordering, (d) the performance cost of explicit cache flushing for a JIT that generates 1MB of code.
8. Dataflow architectures represent a radical departure from both Von Neumann and Harvard: instead of a program counter stepping through sequential instructions, data tokens flow through a graph of operations — an operation executes as soon as all its input tokens arrive. (a) How does this eliminate the Von Neumann bottleneck entirely? (b) What does a "program" look like in a dataflow model? (c) Why didn't dataflow architectures succeed commercially in the 1980s despite theoretical advantages? (d) Where do we see echoes of dataflow thinking in modern hardware (hint: OOO execution, GPU warps, Google TPU systolic arrays)?
