# Computer Architecture
### From Sand and Electricity to Designing Your Own Processor

---

## What This Course Is About

Every smartphone, laptop, server, and smartwatch runs on a processor. That processor is, at its core, billions of tiny switches made from silicon — a material found in ordinary beach sand. This course is the story of how those switches become the most complex machines humanity has ever built.

You will start at the very bottom: atoms, electrons, and transistors. You will build your way up through logic gates, arithmetic circuits, memory cells, and eventually to a complete understanding of how a modern CPU works — the same kind of chip inside your phone or computer right now. Along the way you will meet the great architectures of computing history, learn how chips are actually fabricated in factories, and study the Shakti processor — an open-source CPU designed right here in India by IIT Madras.

By the end of this course, you will understand:

- Why computers use binary (0s and 1s) and how they do arithmetic
- How transistors, logic gates, and flip-flops combine to form a processor
- What an Instruction Set Architecture (ISA) is and why it matters
- The difference between RISC and CISC and why ARM conquered mobile
- How modern CPUs achieve extraordinary performance through pipelining, caching, and out-of-order execution
- How a chip goes from a designer's idea to physical silicon in a factory
- What makes the IIT Madras SHAKTI processor special and why open-source hardware matters
- Where computing is heading — chiplets, AI accelerators, quantum computers, and beyond

---

## Who Is This Course For

Anyone. Seriously.

If you have never opened a computer before, this course starts from scratch. If you are a programmer who has always wondered what actually happens when your code runs, this course answers that question at every level. If you are an electronics enthusiast, a student preparing for competitive exams, or just someone fascinated by technology, everything here is written to be understood without prior knowledge.

Every new term is defined the moment it appears. Every concept is connected to something you already know from everyday life. Nothing is assumed except curiosity.

---

## How Long Will This Take

```
PART 1 — THE FOUNDATION (Volumes 1-2)        ~  30 hours
PART 2 — INSTRUCTION SET ARCHITECTURE (Vol 3) ~  20 hours
PART 3 — MICROARCHITECTURE (Volume 4)         ~  25 hours
PART 4 — MODERN PROCESSORS (Volume 5)         ~  20 hours
PART 5 — SPECIALIZED ARCHITECTURES (Volume 6) ~  20 hours
PART 6 — HOW CHIPS ARE MADE (Volume 7)        ~  20 hours
PART 7 — THE SHAKTI PROCESSOR (Volume 8)      ~  15 hours
PART 8 — THE FUTURE (Volume 9)                ~  10 hours
                                              ----------
GRAND TOTAL:                                  ~160 hours
```

At one hour per day, that is about five and a half months. Many readers find they move faster because the material is genuinely fascinating.

---

## Full Table of Contents

---

# PART 1 — THE FOUNDATION

---

## Volume 1: From Atoms to Logic

> Before there are computers, there is physics. Before there is physics, there is curiosity. In this volume you will understand how the universe's rules — specifically, how electrons move through certain materials — gave humanity the most powerful tool ever created.

### [Chapter 01: What Is a Computer? The Big Picture](01-what-is-a-computer-the-big-picture.md)

What does the word "computer" actually mean? Before the machine, a "computer" was a person — someone whose job was to perform calculations. We trace the evolution from human computers to mechanical calculators to electronic machines. The five fundamental operations any computer must perform: input, storage, processing, output, and control. A modern computer mapped to these five operations, using a smartphone as the example.

**Key topics:** definition of a computer, five operations, brief history of computing, von Neumann model overview.

---

### [Chapter 02: Electricity and the Transistor — The Switch That Changed Everything](02-electricity-and-the-transistor.md)

Everything in a computer runs on electricity. But what is electricity, exactly? A simple model: electrons flowing through metal like water through a pipe. The concept of voltage (pressure), current (flow rate), and resistance (obstruction). Conductors, insulators, and semiconductors — the three types of materials. Why semiconductors are special: you can control whether they conduct or not. The transistor: a semiconductor switch that turns on or off based on an electrical signal. Why the transistor, invented in 1947, is arguably the most important invention in history.

**Key topics:** voltage, current, resistance, conductors/insulators/semiconductors, the transistor, n-type and p-type silicon, MOSFET basics.

---

### [Chapter 03: Binary — Why Computers Think in 0s and 1s](03-binary-why-computers-think-in-0s-and-1s.md)

A transistor has two states: on and off. This maps perfectly to the number system called binary, which uses only two digits: 0 and 1. Why is binary used instead of decimal (0-9)? Because it is much easier to build reliable electronics with two states than ten. Binary counting: 0, 1, 10, 11, 100... Converting between binary and decimal. Bits and bytes — the units of digital information. Hexadecimal: a shorthand for writing binary more compactly. Representing characters, colors, and sounds as binary numbers.

**Key topics:** binary counting, binary↔decimal conversion, bits, bytes, kilobytes/megabytes/gigabytes, hexadecimal, ASCII.

---

### [Chapter 04: Logic Gates — Building Decisions from Electricity](04-logic-gates-building-decisions-from-electricity.md)

A transistor switches based on one input signal. Connect two transistors cleverly and you can make a circuit that switches based on two input signals combined in interesting ways. This gives us logic gates: AND (output is 1 only if BOTH inputs are 1), OR (output is 1 if EITHER input is 1), NOT (output is the opposite of input), NAND, NOR, and XOR. Every logic gate is drawn with a standard symbol. Truth tables show every possible input combination and the corresponding output. The remarkable fact: every digital circuit in every computer — no matter how complex — is built entirely from these simple gates.

**Key topics:** AND, OR, NOT, NAND, NOR, XOR gates, truth tables, gate symbols, building gates from transistors.

---

### [Chapter 05: Combining Gates — Addition, Comparison, and Selection](05-combining-gates-adders-multiplexers-decoders.md)

With individual gates understood, we combine them into useful circuits. A half adder adds two 1-bit numbers. A full adder adds three bits (two inputs plus a carry from a previous addition). Chain eight full adders together and you have an 8-bit adder that adds two 8-bit numbers — the core of all computer arithmetic. Comparators tell you whether two numbers are equal or which is larger. Multiplexers (MUX) select one of several inputs to pass through, like a railroad switch. Demultiplexers (DEMUX) route one input to one of several outputs. These building blocks appear everywhere inside a CPU.

**Key topics:** half adder, full adder, ripple-carry adder, comparator, multiplexer, demultiplexer, decoder, encoder.

---

### [Chapter 06: Memory from Logic — Flip-Flops and Registers](06-memory-from-logic-flip-flops-and-registers.md)

Logic gates alone cannot remember anything — they only respond to current inputs. To store a bit, we need a circuit whose output feeds back into its input, locking it in a stable state. The SR latch is the simplest such circuit: it has a Set input (force output to 1) and a Reset input (force output to 0). The D flip-flop improves on this: it captures its input value on a clock edge and holds it until the next clock edge. The clock is a square wave signal that oscillates at billions of cycles per second, synchronizing all activity inside the CPU. Eight D flip-flops wired together form an 8-bit register — a storage location that holds one byte.

**Key topics:** SR latch, D flip-flop, clock signal, clock frequency (Hz, GHz), registers, why memory needs feedback.

---

## Volume 2: Building a Simple CPU

> You now have the building blocks. In this volume you assemble them into a working computer. Not a modern supercharged processor — a simple, clean, understandable CPU that does real computation. Modern processors are this same machine, turbocharged.

### [Chapter 07: What Is a CPU? Anatomy of a Processor](07-what-is-a-cpu-anatomy-of-a-processor.md)

The Central Processing Unit is the brain of a computer. Every program you run is eventually broken down into a sequence of simple instructions — add, subtract, compare, jump — and the CPU executes them one by one, billions of times per second. Inside the CPU: the Arithmetic Logic Unit (ALU) performs calculations, registers hold the data being worked on, the control unit directs traffic, and the program counter tracks which instruction comes next. A first look at what an instruction looks like in machine code: a sequence of bits encoding what operation to do and what data to operate on.

**Key topics:** CPU components overview, ALU, registers, control unit, program counter, instruction register, machine code format.

---

### [Chapter 08: The Arithmetic Logic Unit — The Math Engine](08-the-arithmetic-logic-unit.md)

The ALU is the part of the CPU that actually does math and logic. It takes two input values, an operation code (opcode), and produces a result plus a set of status flags. An 8-bit ALU can add, subtract, AND, OR, XOR, NOT, and shift numbers. We build a simplified 4-bit ALU from gates and see exactly how it works. The carry flag (was there overflow?), the zero flag (was the result zero?), and the negative flag (was the result negative?) — these flags are how the CPU makes decisions. The connection between arithmetic and comparison: subtraction and checking flags is how a CPU implements "if A > B".

**Key topics:** ALU operations, opcodes, status flags (carry, zero, negative, overflow), building an ALU from gates, signed vs unsigned numbers, two's complement.

---

### [Chapter 09: Registers — The CPU's Working Memory](09-registers-the-cpus-working-memory.md)

A register is the fastest storage in the entire computer — it sits directly inside the CPU chip, just a few electrons' travel from the ALU. Modern CPUs have between 8 and 32 general-purpose registers. Why so few? Because registers are expensive to build at this speed. Each register holds one "word" — the native data size of the CPU (8-bit, 16-bit, 32-bit, or 64-bit depending on the processor). Special-purpose registers: the Program Counter (PC) holds the address of the next instruction to execute, the Stack Pointer (SP) tracks the top of the call stack, and the Flags Register holds the ALU status bits.

**Key topics:** general-purpose registers, register width (word size), program counter, stack pointer, flags register, register file, why registers matter for speed.

---

### [Chapter 10: The Control Unit — The Conductor](10-the-control-unit-the-conductor.md)

The control unit is the CPU's orchestrator. It reads an instruction, figures out what it means, and generates the electrical signals that make everything else do the right thing at the right time. Hardwired control: the instruction bits directly connect to gates that produce control signals — fast but inflexible. Microprogrammed control: the control unit itself contains a tiny internal memory that maps instructions to sequences of simpler micro-operations — slower but easier to design and modify. Modern CPUs use a hybrid approach. The instruction cycle: how the control unit sequences through fetch, decode, and execute in lockstep with the clock.

**Key topics:** control signals, hardwired vs microprogrammed control, the instruction cycle, timing and clocking, how the control unit reads opcodes.

---

### [Chapter 11: The Fetch-Decode-Execute Cycle](11-the-fetch-decode-execute-cycle.md)

This is the heartbeat of every computer. Every instruction executed by every program in history has passed through these three steps: Fetch (read the next instruction from memory into the instruction register), Decode (figure out what the instruction means — what operation, what data), Execute (perform the operation and store the result). We trace a concrete example — the simple program "add 5 and 7, store the result" — through every step of the cycle, showing exactly what happens in the CPU on each clock tick. The cycle repeats billions of times per second, working through the program instruction by instruction.

**Key topics:** fetch stage, decode stage, execute stage, memory read, instruction register, PC increment, concrete execution trace.

---

### [Chapter 12: Memory — Where Programs Live](12-memory-where-programs-live.md)

RAM (Random Access Memory) stores both the program's instructions and its data while the program is running. Every byte in RAM has a unique address — a number identifying its location, like a house number on a very long street. The CPU reads from and writes to RAM using the memory bus: a set of wires for the address, another set for the data, and control wires indicating read or write. Stack memory: the region of RAM used for function calls and local variables, managed by the stack pointer. Heap memory: the region used for dynamically allocated data. ROM (Read-Only Memory): memory that cannot be changed, used to store the firmware that starts the computer.

**Key topics:** RAM, ROM, memory address, address space, memory bus, address lines, data lines, stack vs heap, memory map.

---

### [Chapter 13: Putting It Together — A Simple CPU in Action](13-putting-it-together-a-simple-cpu-in-action.md)

We now trace a complete, realistic program through our simple CPU: a loop that sums the numbers from 1 to 10. We write the program in assembly language (human-readable machine code) and trace every instruction, every register change, and every memory access. By the end of this chapter you will have seen a real program execute from first instruction to last, using only the concepts you have already learned. This chapter ties together all of Part 1 and sets up the question that drives Part 2: what is the right set of instructions to give a CPU, and why do different CPUs have different instruction sets?

**Key topics:** a complete execution trace, assembly language preview, loop implementation, how programs terminate, connecting all components.

---

# PART 2 — INSTRUCTION SET ARCHITECTURE

---

## Volume 3: The Language of Processors

> Every CPU speaks a language — a precisely defined set of instructions it knows how to execute. That language is the Instruction Set Architecture (ISA). In this volume you will learn what an ISA is, why it matters enormously, and meet the major ISAs that power today's world.

### [Chapter 14: What Is an Instruction Set Architecture?](14-what-is-an-instruction-set-architecture.md)

The ISA is the contract between hardware and software. It defines every instruction the CPU can execute, the data types it operates on, the registers it provides, how memory is addressed, and how exceptions are handled. This contract is critical: software written for one ISA runs only on processors that implement that ISA. Intel x86, ARM, RISC-V, and MIPS are all different ISAs — they speak entirely different languages. The ISA is the most important design decision a processor architect makes. An ISA can outlive its implementation by decades: x86 is 46 years old and still runs every Windows program ever written.

**Key topics:** ISA definition, ISA vs microarchitecture, compatibility and portability, major ISAs, why ISAs persist so long, ABI (Application Binary Interface).

---

### [Chapter 15: CISC vs RISC — The Great Debate](15-cisc-vs-risc-philosophy-wars.md)

For decades, processor designers argued about the right philosophy for an ISA. CISC (Complex Instruction Set Computer) says: give the programmer many powerful instructions, even complex ones like "load a value from memory, multiply it by a constant, and store the result." Intel x86 is CISC. RISC (Reduced Instruction Set Computer) says: keep instructions simple and fast — each should execute in one clock cycle. Let the compiler combine simple instructions to do complex things. ARM is RISC. The surprising winner: RISC won the logical argument, but modern x86 processors secretly convert their CISC instructions into RISC-like micro-operations internally. Today the distinction is less important than it once was.

**Key topics:** CISC philosophy, RISC philosophy, the RISC revolution of the 1980s, x86 as CISC, ARM as RISC, micro-operations, why RISC won technically but CISC survived commercially.

---

### [Chapter 16: Assembly Language — Talking Directly to the CPU](16-assembly-language-speaking-to-the-cpu.md)

Assembly language is a human-readable form of machine code. Each line is one CPU instruction, written with a mnemonic (MOV for move, ADD for add, JMP for jump) instead of binary. An assembler converts assembly to the binary machine code the CPU actually runs. We write and trace short assembly programs for a simplified CPU: moving values into registers, arithmetic, branching (if/else), and loops. Assembly is almost never used for application programming today — compilers do it for us — but understanding it gives deep insight into how any program runs, and it is still used for performance-critical code and hardware drivers.

**Key topics:** assembly syntax, mnemonics, registers in assembly, MOV/ADD/SUB/MUL/DIV/CMP/JMP/JE/JNE instructions, assembler, object code, the compiler toolchain.

---

### [Chapter 17: The x86 Architecture — Intel's Legacy](17-x86-the-architecture-that-refused-to-die.md)

x86 is the ISA inside every Windows laptop, every desktop PC, and every server in most data centers. It began as the 8086 processor in 1978, with 16-bit instructions. It grew to 32-bit (called IA-32 or x86-32) and then to 64-bit (called x86-64 or AMD64 — yes, AMD designed the 64-bit extension). x86 has accumulated enormous complexity over 46 years — it still supports instructions from the original 1978 chip. The x86 register set: EAX, EBX, ECX, EDX (32-bit), RAX, RBX, RCX, RDX (64-bit). The x86 memory model. SIMD extensions: MMX, SSE, AVX — special instructions that process 4, 8, or 16 numbers simultaneously for multimedia and AI workloads.

**Key topics:** x86 history, 8086/286/386/486/Pentium lineage, 32-bit and 64-bit modes, x86 registers, SIMD, why x86 is everywhere.

---

### [Chapter 18: ARM — The King of Mobile and Embedded](18-arm-the-architecture-that-conquered-mobile.md)

ARM (Advanced RISC Machines) is the ISA in virtually every smartphone, tablet, and embedded device on earth. Apple's iPhone, every Android phone, Raspberry Pi, the Nintendo Switch — all ARM. ARM is a British company that does not make chips; it designs ISAs and sells licenses. Chip makers (Apple, Qualcomm, Samsung) pay for the license and design their own implementations. The result: hundreds of different ARM chips all running the same code. ARM64 (also called AArch64) is the modern 64-bit version. ARM's defining advantage: power efficiency. An ARM chip doing the same work as an x86 chip uses significantly less electricity — critical for battery-powered devices.

**Key topics:** ARM history, ARM vs x86 in mobile, ARM licensing model, ARM64 / AArch64, Thumb/Thumb-2 instruction encoding, ARM Cortex families (A, R, M), why ARM wins on power efficiency.

---

### [Chapter 19: RISC-V — The Open Source Revolution](19-risc-v-the-open-source-revolution.md)

RISC-V (pronounced "RISC Five") is a new ISA designed at UC Berkeley in 2010 and released completely free and open-source. Anyone can implement a RISC-V processor without paying any license fees. This is revolutionary — x86 costs billions, ARM licensing costs millions. RISC-V is now gaining momentum in academia, startups, and large companies alike. SiFive, Andes Technology, and Western Digital ship RISC-V chips. Google, Samsung, and NVIDIA use RISC-V internally. India's SHAKTI processor (which gets its own volume) is based on RISC-V. The design: a clean, modular ISA where the base integer instruction set is small and optional extensions (floating point, multiplication, vector) can be added as needed.

**Key topics:** RISC-V history and motivation, open ISA vs proprietary ISA, RISC-V base ISA and extensions (I, M, A, F, D, C), growing adoption, RISC-V vs ARM comparison.

---

### [Chapter 20: MIPS — The Teaching Architecture](20-addressing-modes-and-memory-models.md)

MIPS (Microprocessor without Interlocked Pipeline Stages) was designed at Stanford in 1981 and became the dominant workstation and SGI graphics workstation processor in the 1980s-1990s. Today MIPS powers routers (Cisco), set-top boxes, and some embedded systems. Most importantly, MIPS is the architecture taught in virtually every computer architecture university course worldwide because of its elegant, clean design. MIPS has 32 registers (all general-purpose except for a few conventions), a fixed 32-bit instruction width, three instruction formats (R-type, I-type, J-type), and a load/store architecture (only dedicated load and store instructions touch memory; all other instructions operate on registers).

**Key topics:** MIPS history, 32 registers, three instruction formats, load/store architecture, MIPS in education, MIPS R-type/I-type/J-type, delay slots.

---

### [Chapter 21: Addressing Modes — How CPUs Find Data in Memory](21-calling-conventions-and-the-abi.md)

An instruction needs to specify where its operands are: in a register, embedded directly in the instruction, stored in memory, or calculated from a combination. These different ways of specifying data location are called addressing modes. Immediate: the value is in the instruction itself (ADD R0, #5 adds 5 directly). Register: the value is in a register. Direct: the instruction contains the memory address. Indirect: the instruction contains a register whose value is the memory address (like a pointer). Indexed: address = base register + offset (used for arrays). Each mode trades speed for flexibility. Modern ISAs usually support several modes. RISC architectures tend to have fewer, simpler modes; CISC architectures often have many complex modes.

**Key topics:** immediate, register, direct, indirect, indexed, base+offset, PC-relative addressing, why different modes exist, RISC vs CISC addressing philosophy.

---

### [Chapter 22: The Function Call — Stack, Frame, and Convention](22-instruction-encoding-and-decoding.md)

How does a program call a function? When `main` calls `sort`, the CPU must: save the current position (so it can return), pass arguments to `sort`, jump to `sort`'s code, execute it, receive the return value, and jump back. All of this is managed with the stack: a region of memory where values can be pushed (added) and popped (removed) in Last-In-First-Out order. A stack frame is the block of stack memory a function uses for its local variables and saved registers. A calling convention is the agreed set of rules for which registers hold arguments, which hold return values, and which must be preserved across calls. Every ISA has its own calling convention.

**Key topics:** call stack, stack pointer, stack frame (activation record), calling conventions, leaf vs non-leaf functions, the CALL and RET instructions, function prologue/epilogue.

---

# PART 3 — MICROARCHITECTURE

---

## Volume 4: How Modern CPUs Achieve Their Speed

> Knowing the ISA tells you what a CPU can do. Microarchitecture is how it does it fast. This volume reveals the remarkable engineering inside modern processors that turns a simple fetch-decode-execute cycle into a multi-billion-instruction-per-second machine.

### [Chapter 23: Pipelining — Doing Multiple Things at Once](23-pipelining-doing-many-things-at-once.md)

A car factory does not build one car at a time, finish it completely, and then start the next. It has an assembly line: while one car gets its engine installed, another gets its doors attached, and a third gets painted. This is pipelining, and modern CPUs do the same thing with instructions. Instead of completing one instruction before starting the next, the CPU breaks instruction processing into stages and works on multiple instructions simultaneously — each at a different stage. A classic 5-stage pipeline: Fetch, Decode, Execute, Memory Access, Writeback. With this pipeline, the CPU can start one new instruction every clock cycle instead of one every five.

**Key topics:** pipeline stages, pipeline depth, throughput vs latency, ideal speedup, pipeline registers (inter-stage buffers), pipeline diagram, why pipelining matters.

---

### [Chapter 24: Pipeline Hazards — When the Pipeline Stalls](24-pipeline-hazards-when-the-pipeline-stalls.md)

The ideal pipeline always has an instruction in every stage. Reality is messier. Three types of hazards break the pipeline's flow. Data hazards: instruction B needs a result that instruction A has not finished computing yet — B must wait. Control hazards: a branch instruction forces the CPU to fetch instructions before knowing which path to take — if it guesses wrong, wasted work must be discarded. Structural hazards: two instructions need the same hardware resource at the same time — one must wait. Solutions: stalling (inserting bubbles — empty pipeline slots), forwarding (feeding a result directly to the next instruction without waiting for it to be written to a register), and pipeline reordering (a compiler or hardware reorders instructions to avoid hazards).

**Key topics:** data hazards (RAW, WAR, WAW), control hazards, structural hazards, pipeline stalls (bubbles), data forwarding/bypassing, hazard detection unit, CPI (cycles per instruction).

---

### [Chapter 25: Branch Prediction — Guessing the Future](25-branch-prediction-guessing-the-future.md)

Every if/else and loop in your program generates a branch instruction. When the CPU encounters a branch, it does not yet know which path to take — it would have to fully evaluate the condition first, which takes cycles. Meanwhile, the pipeline is fetching the next instructions. If the CPU guesses the branch outcome wrong, it has fetched and partially executed the wrong instructions and must discard that work (a branch misprediction penalty of 15-20 clock cycles on modern CPUs). Branch predictors are hardware circuits that study patterns in which branches are taken. A simple predictor: always predict loops are taken. Modern predictors use sophisticated pattern-matching with accuracy above 99% for typical programs.

**Key topics:** why branches are expensive, branch misprediction penalty, static prediction, dynamic prediction (BHT), two-bit saturating counters, global history, the branch target buffer (BTB), speculative execution.

---

### [Chapter 26: Out-of-Order Execution — Working Ahead](26-out-of-order-execution.md)

In-order execution: the CPU processes instructions in the exact order the program lists them, waiting whenever there is a dependency. Out-of-order execution: the CPU looks ahead at upcoming instructions and executes any that are ready — even if earlier instructions are still waiting. If instruction 5 does not depend on instructions 3 or 4, the CPU can execute instruction 5 first. This dramatically increases the percentage of time the CPU's execution units are busy. The critical data structure: the Reorder Buffer (ROB) — a circular queue that tracks in-flight instructions and ensures results are committed to architectural state in program order even if executed out-of-order.

**Key topics:** in-order vs out-of-order, instruction window, reorder buffer (ROB), reservation stations, Tomasulo's algorithm, register renaming (eliminating false dependencies), retirement/commit.

---

### [Chapter 27: Superscalar Execution — Multiple Pipelines](27-superscalar-processors.md)

A scalar processor issues one instruction per clock cycle. A superscalar processor issues multiple instructions per clock cycle. A 4-wide superscalar can issue four instructions simultaneously to four independent execution units. This requires quadrupling much of the pipeline: four fetch units, four decode units, multiple ALUs, multiple floating-point units. The challenge: finding enough independent instructions to keep all execution units busy. IPC (Instructions Per Clock) is the key metric — modern high-performance CPUs achieve 3-5 IPC. The scheduler (dispatch unit) continuously scans ready instructions and issues them to whichever execution unit is free.

**Key topics:** ILP (instruction-level parallelism), superscalar width, execution units (integer ALU, FPU, load/store unit, branch unit), dispatch/issue logic, IPC, theoretical vs practical ILP.

---

### [Chapter 28: Cache Memory — The Speed Gap Solution](28-cache-memory-the-speed-gap-solution.md)

The speed gap between CPU registers (1 cycle access) and RAM (hundreds of cycles access) would make the CPU wait constantly for data. The solution is cache: small, fast memory placed between the CPU and RAM that stores copies of recently accessed data. The key insight: locality. Programs tend to access the same data repeatedly (temporal locality) and to access data near recently accessed data (spatial locality). When the CPU needs data, it first checks the cache (a cache hit — data found, fast). If not found (a cache miss), it fetches from RAM and stores a copy in cache for next time. Cache capacity is measured in KB or MB; access time is measured in nanoseconds.

**Key topics:** cache motivation, temporal and spatial locality, cache hit and miss, cache lines, cache size, cache access time, direct-mapped vs set-associative vs fully-associative cache, LRU replacement.

---

### [Chapter 29: Cache Hierarchy — L1, L2, L3](29-cache-hierarchy-l1-l2-l3.md)

Modern CPUs have multiple levels of cache, each larger and slower than the previous. L1 cache: 32-64 KB, accessed in 4 cycles, split into instruction cache (L1-I) and data cache (L1-D), one per CPU core. L2 cache: 256 KB - 1 MB, accessed in 12 cycles, usually per-core. L3 cache (Last Level Cache, LLC): 8-64 MB, accessed in 30-40 cycles, shared across all cores on the chip. When the CPU needs data, it checks L1 first, then L2, then L3, then RAM. Each miss is more expensive than the last. Cache coherence: the problem of keeping cache copies consistent when multiple CPU cores each have their own L1 and L2 caches but share the same memory.

**Key topics:** L1/L2/L3 cache sizes and latencies, inclusion vs exclusion, victim cache, prefetching, cache coherence problem, MESI protocol (Modified/Exclusive/Shared/Invalid), cache thrashing.

---

### [Chapter 30: Virtual Memory and the MMU](30-virtual-memory-and-the-mmu.md)

Every program thinks it has the entire computer's memory to itself — its own address space starting from 0. This is virtual memory: an illusion maintained by the OS and hardware. The CPU's Memory Management Unit (MMU) translates virtual addresses (what the program uses) to physical addresses (actual RAM locations) using a data structure called a page table. Virtual memory enables: process isolation (each process has its own address space, cannot access other processes' memory), memory overcommitment (use more virtual memory than physical RAM, swapping inactive pages to disk), and simplified program loading. The Translation Lookaside Buffer (TLB) caches recent address translations so the MMU does not need to walk the page table for every memory access.

**Key topics:** virtual vs physical address, page table, page size (4KB), page fault, TLB, TLB miss, address space layout, OS kernel vs user space, segmentation vs paging.

---

### [Chapter 31: The Memory Bus and I/O](31-the-memory-bus-and-io.md)

The CPU communicates with RAM, storage, graphics cards, and peripherals through buses — shared communication channels. The front-side bus (old) and modern equivalents (PCIe, HyperTransport, Infinity Fabric) connect the CPU to memory and peripherals. Memory controllers (once on a separate chip called the northbridge, now integrated into the CPU) manage the communication with RAM. PCIe (Peripheral Component Interconnect Express) is the modern high-speed bus for graphics cards, SSDs, and network cards. Device drivers are software that let the OS speak to specific hardware devices. DMA (Direct Memory Access) lets devices transfer data to/from RAM without involving the CPU for each byte.

**Key topics:** system bus, PCIe, memory controller, northbridge/southbridge (historical), DRAM channels, DMA, memory-mapped I/O, port-mapped I/O, device drivers.

---

### [Chapter 32: Interrupts and Exceptions — Handling the Unexpected](32-interrupts-and-exceptions.md)

The CPU does not poll every peripheral constantly to ask "do you have data for me?" Instead, peripherals signal the CPU by raising an interrupt: an electrical signal that causes the CPU to pause its current work, save its state, and execute an interrupt service routine (ISR) — a special piece of code that handles the event. Timer interrupts let the OS scheduler switch between programs. Hardware interrupts signal events (keyboard press, network packet arrived, disk read complete). Software interrupts (also called traps or system calls) let user programs request services from the OS. Exceptions are internally-generated interrupts for exceptional conditions: division by zero, page fault, illegal instruction.

**Key topics:** interrupt vector table, ISR, hardware vs software interrupts, exceptions, interrupt latency, interrupt masking, system calls, privilege levels (ring 0 vs ring 3).

---

# PART 4 — MODERN PROCESSORS

---

## Volume 5: The Chips That Power Today's World

> Abstract architecture has become concrete silicon. This volume tours the real processors inside the devices you use every day — Intel, AMD, Apple, Qualcomm, and the GPU giants — explaining both what makes them extraordinary and how they connect to everything you have learned.

### [Chapter 33: Intel x86 — From 8086 to Core Ultra](33-intel-x86-from-8086-to-core-ultra.md)

Intel invented the x86 ISA in 1978 with the 8086. Forty-six years of iterative improvement produced the most commercially successful processor line in history. The key milestones: 8086 (1978, 16-bit), 386 (1985, 32-bit, protected mode), Pentium (1993, first superscalar x86), Core 2 (2006, modern microarchitecture reboot), Sandy Bridge (2011, integrated graphics), Haswell (2013, deep out-of-order), Skylake (2015), and the modern hybrid architecture with Performance and Efficiency cores (Alder Lake 2021, Raptor Lake 2022, Meteor Lake 2023, Core Ultra). Intel's 10nm and 7nm process nodes. Ring Bus and Mesh interconnect. The Ring 0 / Ring 3 privilege system. Why x86 remains dominant in laptops and servers despite ARM's rise.

**Key topics:** Intel microarchitecture history, P-core vs E-core hybrid design, Intel's process node naming, Ring Bus, mesh interconnect, Hyper-Threading (SMT), Intel's challenges post-2015.

---

### [Chapter 34: AMD — From Near-Bankruptcy to Zen Dominance](34-amd-from-near-bankruptcy-to-zen-dominance.md)

AMD was founded in 1969 to compete with Intel. It licensed x86 from Intel, then created its own 64-bit extension (AMD64 / x86-64). By 2012, AMD was nearly bankrupt. Then came Zen (2017) — a complete microarchitecture redesign led by Jim Keller that caught up with Intel's performance within two years. Zen 2 (2019), Zen 3 (2020), and Zen 4 (2022) pushed AMD ahead in many workloads. AMD's strategic advantage: chiplets. Instead of one monolithic die, a Zen 2+ CPU uses multiple smaller dies (one per 8-core "chiplet") connected by AMD's Infinity Fabric. This lets AMD use the most advanced process node for the chiplets while using cheaper nodes for I/O. Threadripper and EPYC for workstations and servers.

**Key topics:** AMD64 (x86-64) history, Zen microarchitecture, chiplet strategy, Infinity Fabric, EPYC servers, AMD vs Intel timeline, TSMC partnership.

---

### [Chapter 35: ARM Cortex — Mobile Dominance](35-arm-the-architecture-that-conquered-mobile.md)

ARM licenses its ISA but also designs reference implementations called Cortex cores, which licensees can use directly or modify. Cortex-A series: high-performance application cores for smartphones and tablets (A53, A55, A75, A78, A710). Cortex-R series: real-time cores for safety-critical systems (automotive ECUs, hard drive controllers). Cortex-M series: microcontroller cores for embedded devices (M0, M3, M4, M7) — found in everything from thermostats to pacemakers. big.LITTLE: ARM's heterogeneous CPU design that pairs fast, power-hungry Cortex-A cores with small, efficient Cortex-A cores, running programs on whichever is appropriate at the moment. This is the architecture inside every Android flagship.

**Key topics:** Cortex-A/R/M series, big.LITTLE architecture, DynamIQ, GlobalFoundries and TSMC for ARM, ARM Mali GPU, ARM TrustZone security, ARMv8/v9 ISA.

---

### [Chapter 36: Apple Silicon — The M-Series Revolution](36-apple-silicon-the-m-series-revolution.md)

In 2020 Apple announced it was switching its Mac computers from Intel x86 to ARM-based chips it designed itself. The M1 was a revelation: it matched or beat Intel's best in performance while using a fraction of the power. How? Apple had been designing ARM chips for iPhones (the A-series) for a decade and had learned to integrate CPU, GPU, neural engine, and memory onto a single chip with an incredibly fast unified memory architecture. The M1's secret weapon: unified memory. Instead of separate CPU RAM and GPU VRAM, all components share a single high-bandwidth pool, eliminating data copies. M2, M3, and M4 improved further. The M-series has forced Intel and AMD to completely rethink their roadmaps.

**Key topics:** Apple Silicon history, M1 unified memory architecture, CPU cluster design (P-cores + E-cores), Apple GPU architecture, Neural Engine, Secure Enclave, Apple's in-house fabrication at TSMC, why Apple Silicon changed the industry.

---

### [Chapter 37: Qualcomm Snapdragon — The Smartphone Chip](37-qualcomm-snapdragon-the-smartphone-chip.md)

Qualcomm's Snapdragon is the chip inside most flagship Android phones (Samsung Galaxy, Google Pixel, Sony Xperia). Unlike Apple, Qualcomm licenses the ARM ISA and designs its own cores — the Kryo series — rather than using ARM reference designs. The Snapdragon 8 Gen 3 (current flagship) uses a prime cluster (1 top performance core), a performance cluster (4 cores), and an efficiency cluster (3 cores). Snapdragon integrates a modem (5G/4G/3G wireless), a GPU (Adreno — widely considered the best mobile GPU), a DSP (Hexagon — for audio, camera, and AI), and a Neural Processing Unit. The integration of the modem directly onto the SoC is a major competitive advantage.

**Key topics:** Qualcomm Kryo cores, Adreno GPU, Hexagon DSP, integrated 5G modem, Snapdragon vs Apple A-series, Snapdragon for laptops (Snapdragon X Elite), custom vs reference ARM cores.

---

### [Chapter 38: GPUs — Parallel Processing Giants](38-gpus-parallel-processing-giants.md)

A CPU has 8-32 powerful cores optimized for sequential, complex tasks. A GPU has thousands of simpler cores optimized for doing the same operation on thousands of data points simultaneously. This is SIMT: Single Instruction Multiple Threads. A GPU can execute the same instruction on 32 or 64 pieces of data at the same time. Why does this matter? 3D graphics requires applying the same transformation (rotation, lighting, shading) to millions of vertices and pixels independently — perfect for SIMT. And it turns out training neural networks is the same kind of math: matrix multiplication applied to billions of numbers. The GPU, designed for games, became the engine of the AI revolution.

**Key topics:** GPU vs CPU architecture philosophy, SIMT, shader cores, warps/wavefronts, GPU memory (GDDR/HBM), CUDA (NVIDIA), OpenCL, the GPU programming model, latency-optimized vs throughput-optimized.

---

### [Chapter 39: NVIDIA GPU Architecture — From Graphics to AI](39-nvidia-gpu-architecture.md)

NVIDIA dominates the GPU market for both gaming and AI. The key milestones: GeForce 256 (1999, first GPU with hardware T&L), Tesla (2006, CUDA launched, GPUs for general computation), Kepler/Maxwell/Pascal (2010s, gaming dominance), Volta (2017, Tensor Cores for matrix multiply, targeting AI), Turing (2018, real-time ray tracing), Ampere (2020, massive AI compute for training), Hopper (2022, H100 — the chip powering ChatGPT's training), and Blackwell (2024). The H100: 80 billion transistors, 700W power consumption, built at TSMC 4nm. NVLink: NVIDIA's proprietary interconnect for tying multiple GPUs together into a single logical unit. Why NVIDIA controls AI: first-mover with CUDA in 2006.

**Key topics:** NVIDIA architecture evolution, Tensor Cores, CUDA cores, NVLink, HBM memory, GPU cluster (DGX), A100 vs H100, CUDA as a moat, NVIDIA's software ecosystem.

---

### [Chapter 40: AMD GPU Architecture and the Competitive Landscape](40-amd-gpu-and-competitive-landscape.md)

AMD (via acquisition of ATI in 2006) is NVIDIA's only serious competitor in discrete GPUs. AMD uses the RDNA architecture for gaming GPUs (RX 7900 series) and CDNA for compute/AI GPUs (Instinct MI300). AMD's Infinity Fabric (the same interconnect used in their CPUs) connects CPU and GPU dies in the MI300X APU — a direct competitor to the NVIDIA H100. ROCm is AMD's alternative to CUDA — open-source GPU compute software. The problem: CUDA has 18 years of ecosystem head start. Intel's Arc GPU: a newcomer entering the discrete GPU market, currently mid-range. The integrated GPU question: Apple's M-series has the best integrated GPU ever, challenging mid-range discrete cards.

**Key topics:** RDNA vs CDNA, ROCm vs CUDA ecosystem gap, MI300 for AI, Intel Arc GPU, integrated vs discrete GPU, AMD's strategy.

---

### [Chapter 41: NPUs and AI Accelerators — Intelligence in Silicon](41-npus-and-ai-accelerators.md)

Neural networks require one operation repeatedly: multiply two numbers and add the result to an accumulator. This is called a MAC (Multiply-Accumulate). A processor designed purely to do MACs as fast as possible will train and run neural networks much more efficiently than a general-purpose CPU or GPU. Google's TPU (Tensor Processing Unit): a custom chip with a massive 2D array of MAC units — a systolic array. Used internally in Google datacenters since 2016. Apple's Neural Engine: 16 cores, 38 trillion operations per second, handles Siri, Face ID, and on-device AI. Qualcomm's Hexagon NPU. The Groq LPU (Language Processing Unit): designed specifically for running LLM inference fast. Why custom silicon is 10-100x more efficient than general-purpose hardware for AI.

**Key topics:** MAC operation, systolic array, Google TPU, Apple Neural Engine, systolic array, TOPS (tera-operations per second), inference vs training, edge AI vs cloud AI, the race to build the best AI chip.

---

### [Chapter 42: RISC-V in Production — Open Hardware Growing Up](42-risc-v-in-production.md)

RISC-V has moved from academic curiosity to production deployment in five years. Western Digital ships hard drives using RISC-V cores (SweRV). SiFive designs commercial RISC-V chips. Alibaba's T-Head semiconductor arm ships RISC-V server chips in China (avoiding US export restrictions on ARM). NVIDIA embeds RISC-V cores as management controllers inside their GPUs. Google uses RISC-V cores in Pixel phone embedded controllers. The ecosystem: the RISC-V International Foundation governs the ISA and ratifies extensions. GCC and LLVM support RISC-V. Linux runs on RISC-V. The OpenHW Group designs open-source RISC-V cores (CVA6, CV32E40P). Why RISC-V matters geopolitically: countries that cannot access ARM or x86 licenses can design their own RISC-V chips.

**Key topics:** RISC-V adoption examples, SiFive, Western Digital, Alibaba T-Head, RISC-V ecosystem maturity, geopolitical significance, open-source hardware cores, RISC-V vs ARM for new designs.

---

### [Chapter 43: IBM POWER — The Server Powerhouse](43-ibm-power-the-server-powerhouse.md)

IBM's POWER architecture is the third major server ISA (after x86 and ARM). POWER (Performance Optimization with Enhanced RISC) runs IBM's mainframes and high-end UNIX servers. POWER10 (2021): built at Samsung 7nm, up to 15 cores per chip, 120 MB of L3 cache, exceptional integer throughput, and MMA (Matrix Math Accelerator) units for AI workloads. POWER runs Linux and AIX (IBM's UNIX). OpenPOWER: IBM opened the POWER ISA in 2019, allowing anyone to build POWER-compatible chips without licensing fees. Used in banking, telecommunications, and national laboratories. The historical context: POWER chips ran every Apple Mac from 1994-2005 before Apple switched to Intel.

**Key topics:** POWER ISA history, OpenPOWER, POWER10 specs, SMT (Simultaneous Multi-Threading up to 8 threads per core), IBM mainframes, POWER in banking/finance, POWER vs x86 servers.

---

# PART 5 — SPECIALIZED ARCHITECTURES

---

## Volume 6: Beyond the General-Purpose Processor

> Not every computing problem needs a general-purpose CPU. In this volume you will explore the specialized architectures designed for specific jobs — from the tiny microcontrollers in your toaster to the reconfigurable silicon of FPGAs to the quantum computers that exist at the edge of physical possibility.

### [Chapter 44: Harvard vs Von Neumann Architecture](44-harvard-vs-von-neumann-architecture.md)

Most modern computers use the Von Neumann architecture: a single memory space holds both program instructions and data, accessed over the same bus. The Von Neumann bottleneck: because instructions and data share the bus, you can never fetch an instruction and read data simultaneously. The Harvard architecture uses separate memory spaces and buses for instructions and data — eliminating the bottleneck. Modern microcontrollers (Arduino, PIC, AVR) often use Harvard architecture. Modern CPUs use a Modified Harvard architecture: physically they have unified RAM, but the L1 cache is split into separate instruction and data caches, giving Harvard's performance without its inflexibility.

**Key topics:** Von Neumann model, Von Neumann bottleneck, Harvard architecture, modified Harvard, Harvard in microcontrollers, split L1 cache as modified Harvard, tradeoffs.

---

### [Chapter 45: System on a Chip — Everything in One Package](45-system-on-a-chip.md)

A System on a Chip (SoC) integrates a CPU, GPU, memory controller, modem, DSP, USB/HDMI/camera/audio controllers, and sometimes RAM itself — all on a single piece of silicon. Your smartphone is an SoC. The iPhone's A18 Pro SoC contains 19 billion transistors. The advantages: no communication latency between components (they share high-speed on-chip buses instead of slow board-level connectors), dramatically lower power consumption, and smaller form factor. The challenge: designing all these components is enormously complex. SoC design requires hundreds of engineers and costs billions of dollars. TSMC and Samsung provide IP blocks (verified hardware designs for USB, PCIe, etc.) that SoC designers license to speed up development.

**Key topics:** SoC definition, smartphone SoC examples, NoC (Network on Chip), IP licensing, on-chip bus (AXI, AHB), SoC vs discrete components, SoC design flow.

---

### [Chapter 46: Microcontrollers — The Processors in Everything](46-microcontrollers-processors-in-everything.md)

A microcontroller (MCU) is a tiny, complete computer on a single chip: CPU + RAM + Flash storage + I/O peripherals + timers + ADC — everything. They run at 8 MHz to 400 MHz, have kilobytes to megabytes of RAM, and consume milliwatts of power. They are everywhere: washing machines, keyboards, car engine management, medical devices, toys, thermostats. The Arduino uses an ATmega328P (AVR 8-bit microcontroller). The Raspberry Pi Pico uses an RP2040 (ARM Cortex-M0+ 32-bit microcontroller). STM32 (ARM Cortex-M3/M4/M7/M33) dominates industrial embedded systems. Programming microcontrollers in C without an operating system: direct register manipulation, interrupt-driven I/O, real-time constraints.

**Key topics:** MCU vs CPU distinction, MCU components, AVR, PIC, ARM Cortex-M, STM32, Arduino, Raspberry Pi Pico, bare-metal programming, RTOS (Real-Time Operating System).

---

### [Chapter 47: Digital Signal Processors — Masters of Math](47-digital-signal-processors.md)

A DSP (Digital Signal Processor) is optimized for repeatedly applying mathematical operations to streams of numbers — exactly what happens when processing audio, video, radar signals, or communications. DSPs have hardware multiply-accumulate (MAC) units that perform a multiply and an add in a single cycle (a general CPU needs multiple cycles). They have special addressing modes for circular buffers (used in streaming processing) and bit-reversal (used in Fast Fourier Transforms). They can process multiple data samples per cycle (SIMD-like operations). You find DSPs in: every phone's audio system (noise cancellation, echo suppression), every 4G/5G modem, every hearing aid, camera image signal processors, and automotive radar.

**Key topics:** DSP applications, MAC unit, circular buffer addressing, FFT hardware support, fixed-point vs floating-point DSP, TI C6000/C28x, Analog Devices SHARC, DSP vs GPU for signal processing.

---

### [Chapter 48: FPGAs — Reconfigurable Hardware](48-fpgas-reconfigurable-hardware.md)

An FPGA (Field-Programmable Gate Array) is an integrated circuit whose logic can be reconfigured after manufacturing. Instead of fixed logic circuits (like a CPU or DSP), an FPGA contains thousands of programmable logic blocks connected by programmable interconnects — you literally program what circuits it contains. FPGAs are programmed using Hardware Description Languages (VHDL or Verilog) — languages that describe hardware rather than software. Use cases: FPGA prototyping of ASIC designs before fabrication, hardware acceleration (100-1000x faster than software for specific algorithms), low-latency trading systems, communications infrastructure (5G base stations), aerospace and defense. The trade-off: FPGAs are 10-100x less power-efficient than a custom ASIC doing the same job.

**Key topics:** FPGA architecture (LUTs, flip-flops, DSP blocks, BRAM), Xilinx (now AMD) and Intel (Altera) FPGA families, FPGA vs ASIC vs CPU, HDL introduction, partial reconfiguration, FPGAs in cloud (AWS F1, Azure FPGA).

---

### [Chapter 49: ASICs — Custom Silicon for One Purpose](49-asics-custom-silicon.md)

An ASIC (Application-Specific Integrated Circuit) is a chip designed to do one thing. Unlike a CPU (which runs any program) or an FPGA (which can be reprogrammed), an ASIC's function is permanently baked into its silicon. The trade-offs: extremely high non-recurring engineering (NRE) cost (tens to hundreds of millions of dollars to design and tape out), but once made, ASICs are the fastest and most power-efficient way to implement any function. Famous ASICs: Bitcoin mining chips (ASICs make mining profitable where GPUs cannot), Google's TPU (custom AI matrix multiply), Apple's Neural Engine, the AirPods H1/H2 (custom audio processing chip), every smartphone modem.

**Key topics:** ASIC vs FPGA vs CPU trade-offs, NRE cost, ASIC design flow, foundry (TSMC, Samsung, GlobalFoundries), standard cells, tape-out, why ASICs dominate high-volume applications.

---

### [Chapter 50: Quantum Computing — A Different Kind of Computer](50-quantum-computing.md)

Classical computers use bits: each is 0 or 1. Quantum computers use qubits: each can be in a superposition of 0 and 1 simultaneously until measured. This is not the same as "both at once" — it is a fundamentally different kind of information that follows quantum mechanical rules. Entanglement allows qubits to have correlated states across any distance. Quantum interference lets algorithms amplify correct answers and cancel wrong ones. Quantum algorithms like Shor's algorithm can factor large numbers exponentially faster than any classical algorithm — threatening current encryption. But quantum computers are not universally faster: they excel at specific problems (factoring, simulation, optimization) and fail at others (video games, web browsers). Current state: IBM has 1,121-qubit processors; Google claims quantum supremacy for specific tasks; error rates remain high.

**Key topics:** qubit, superposition, entanglement, quantum gates, quantum circuits, Shor's algorithm, Grover's algorithm, quantum decoherence, error correction, current limitations, D-Wave vs gate-based quantum computers.

---

### [Chapter 51: Neuromorphic Computing — Computing Like a Brain](51-neuromorphic-computing.md)

Traditional computers execute instructions sequentially on static data. The human brain processes information through billions of neurons firing in parallel, with connections (synapses) that strengthen or weaken based on use. Neuromorphic computing tries to mimic this architecture. Intel's Loihi 2 chip contains 1 million artificial neurons and 120 million synapses, all running in parallel and consuming only 1 watt — learning and adapting without explicit programming. IBM's TrueNorth: 1 million neurons, 256 million synapses in 70 mW. Current applications: event-based sensing (cameras that only respond to changes, not frames), edge inference, robotics. The promise: 1000x more energy efficient than traditional AI for certain tasks.

**Key topics:** biological neuron basics, spiking neural networks (SNN), Intel Loihi, IBM TrueNorth, event-based cameras, neuromorphic vs deep learning, current limitations, future potential.

---

---

# PART 6 — HOW CHIPS ARE MADE

---

## Volume 7: From Sand to Silicon

> Every chip starts as raw sand on a beach. In this volume you will follow the extraordinary journey from silicon dioxide to the most precisely manufactured objects humanity has ever produced — and meet the tools used to design them.

### [Chapter 52: Silicon — The Semiconductor Material](52-silicon-the-semiconductor-material.md)

Silicon is the second most abundant element in the Earth's crust (after oxygen) and the foundation of virtually all modern electronics. Pure silicon is a semiconductor: an insulator at low temperatures that conducts electricity when energy is added or when it is "doped" with impurities. N-type silicon: doped with phosphorus or arsenic, which donate extra electrons (negative charge carriers). P-type silicon: doped with boron or gallium, which create "holes" (missing electrons that act as positive charge carriers). A p-n junction — where p-type and n-type silicon meet — is a diode: current flows easily in one direction but not the other. Two p-n junctions make a bipolar transistor. A field-effect transistor (MOSFET) adds a gate electrode that controls the p-n junction electrically.

**Key topics:** silicon properties, doping, n-type and p-type silicon, p-n junction, diode, bipolar transistor vs MOSFET, CMOS (Complementary MOS), why CMOS dominates.

---

### [Chapter 53: The CMOS Process — Turning Silicon into Logic](53-the-cmos-process.md)

CMOS (Complementary Metal-Oxide-Semiconductor) is the process technology used to make virtually all modern chips. Every logic gate is built from pairs of complementary transistors: one n-type MOSFET and one p-type MOSFET. CMOS gates draw power only when switching (not when holding a state), making them far more power-efficient than older technologies. The silicon wafer: a thin circular slice (300mm diameter for modern chips) cut from a large cylindrical crystal of ultra-pure silicon. The chip manufacturing process starts with this wafer and adds hundreds of layers of materials using deposition, photolithography, and etching — each step precisely controlled to create circuits with features measured in nanometers.

**Key topics:** CMOS inverter/NAND/NOR from complementary transistors, static vs dynamic power, silicon wafer, wafer diameter and die yield, process layers, why CMOS is so efficient.

---

### [Chapter 54: Photolithography — Printing Circuits at Nanoscale](54-photolithography.md)

Photolithography is the process of "printing" circuit patterns onto silicon — like a very precise photographic development process but with features measured in billionths of a meter. The silicon wafer is coated with a light-sensitive material (photoresist). A light source (originally UV, now extreme ultraviolet — EUV — with 13.5nm wavelength) shines through a photomask (a stencil of the circuit pattern). The photoresist exposed to light becomes chemically different and can be washed away, leaving the pattern. This is repeated hundreds of times with different masks to build up the layers of the chip. ASML (a Dutch company) is the only manufacturer of EUV lithography machines, which cost $370 million each and are crucial for sub-5nm chips.

**Key topics:** photoresist, photomask, UV vs EUV lithography, ASML EUV machine, step-and-scan, diffraction limit, immersion lithography, multi-patterning, how EUV enables sub-5nm features.

---

### [Chapter 55: Moore's Law — The Driving Force of the Digital Age](55-moores-law.md)

In 1965 Gordon Moore (co-founder of Intel) observed that the number of transistors on a chip was doubling roughly every two years. This became Moore's Law — not a law of physics but an economic and engineering commitment that the semiconductor industry treated as a target. The consequences: every two years, chips became twice as powerful at the same cost. This sustained exponential improvement drove five decades of computing progress. Moore's Law is slowing: at 5nm and below, quantum tunneling effects cause transistors to leak current even when "off." The performance gains from shrinking transistors are diminishing. The industry response: 3D stacking (building chips vertically), chiplets, new transistor architectures (FinFET, GAA — Gate-All-Around), and new materials.

**Key topics:** Moore's Law history and numbers, Dennard scaling (why smaller = faster AND cooler), the end of Dennard scaling (2006), Moore's Law slowing, FinFET transistors, Gate-All-Around (GAA) nanosheet transistors, beyond-silicon materials (GaN, SiC, InGaAs).

---

### [Chapter 56: Process Nodes — 5nm, 3nm, 2nm](56-process-nodes.md)

You have seen chips advertised as "5nm" or "3nm." These numbers are marketing labels, not literal measurements of anything — the actual transistor dimensions are different (and often larger) than the number suggests. The node name reflects the generation of process technology, not a physical measurement. What actually matters: transistor density (how many fit per mm²), performance (how fast the transistors switch), and power efficiency (how much energy each switch uses). TSMC's N3 (3nm class) packs 290 million transistors per mm². Intel's Intel 20A introduces PowerVia (backside power delivery). Samsung 3GAE introduced Gate-All-Around nanosheet transistors in production. TSMC's N2 (2nm, 2025) uses GAA nanosheets for the first time.

**Key topics:** process node naming confusion, transistor density (MTr/mm²), TSMC N3/N2, Samsung 3GAE, Intel 20A/18A, FinFET vs GAA nanosheet, backside power delivery, leading-edge vs trailing-edge nodes.

---

### [Chapter 57: EDA Tools and the Chip Design Flow](57-eda-tools-chip-design-flow.md)

Designing a modern chip with billions of transistors by hand is impossible. Electronic Design Automation (EDA) tools are specialized software that assists chip designers at every stage. Schematic entry and RTL (Register Transfer Level) design: describing the chip's behavior in Verilog or VHDL. Logic synthesis: converting RTL to a network of logic gates (netlist) targeting a specific standard cell library. Place and route: determining where each gate goes on the physical silicon die and routing the wires connecting them. Static timing analysis (STA): verifying that all signals arrive within their timing budgets. DRC/LVS: checking that the layout obeys the foundry's design rules. The EDA market: Synopsys, Cadence, and Mentor (Siemens EDA) — three companies whose tools are used for virtually every chip made.

**Key topics:** EDA toolchain overview, Synopsys Design Compiler, Cadence Genus, standard cell library, placement and routing, clock tree synthesis, static timing analysis, sign-off checks, tape-out.

---

### [Chapter 58: Hardware Description Languages — Verilog and VHDL](58-verilog-and-vhdl.md)

Hardware description languages (HDLs) are programming languages for describing hardware. Instead of describing a sequence of operations (like C), they describe the structure and behavior of circuits. Verilog (1984, from Cadence) and VHDL (1987, from the US Department of Defense) are the two dominant HDLs. A module in Verilog describes a hardware component: its input/output ports, its internal logic. RTL (Register Transfer Level) coding style: describing how data moves between registers on each clock edge — this style can be synthesized into real gates. Behavioral simulation: running the HDL on a computer to verify it works before building real hardware. SystemVerilog: Verilog extended with modern programming constructs for verification (interfaces, assertions, classes).

**Key topics:** Verilog syntax basics (module, always, assign, if/else), VHDL overview, RTL vs behavioral coding, simulation vs synthesis, testbenches, SystemVerilog for verification, open-source HDL tools (Verilator, iverilog, Yosys).

---

### [Chapter 59: From RTL to Silicon — The Complete Design Flow](59-rtl-to-silicon.md)

A chip goes through many stages between a designer's first line of Verilog and a working chip in a customer's device. Front-end design: RTL coding, functional simulation, formal verification. Logic synthesis: mapping RTL to gates. Design for Test (DFT): adding scan chains and test logic so the chip can be tested after manufacture. Floorplanning: deciding the high-level layout — where does the CPU core go, where does the cache go, where are the I/O pads. Placement: positioning each of the millions of standard cells. Clock tree synthesis: designing the tree of buffers that distributes the clock signal evenly to every flip-flop. Routing: connecting every gate according to the netlist. Sign-off: final timing, power, and design rule checks. Tape-out: sending the final GDS-II file (the chip layout) to the foundry.

**Key topics:** front-end vs back-end design, logic synthesis, DFT, floorplanning, placement, clock tree synthesis, routing, GDSII, tape-out, NRE cost breakdown, time from tape-out to samples.

---

### [Chapter 60: Chiplets and Advanced Packaging — The New Paradigm](60-chiplets-and-advanced-packaging.md)

For decades, more performance meant putting more on a single chip (monolithic die). At advanced nodes, that becomes impossibly expensive — a large 5nm die has dramatically lower yield (more defects per wafer) than a small one. The solution: chiplets. Split the chip into multiple small dies, each optimized for a different process node, and connect them in one package. AMD Zen 2: 7nm compute chiplets + 14nm I/O die, connected by Infinity Fabric. Apple M1 Ultra: two M1 Max dies connected by UltraFusion (die-to-die bridge), appearing as one chip to software. Intel Foveros 3D: stacking chiplets vertically on top of each other for even shorter connections. TSMC CoWoS (Chip on Wafer on Substrate): the interposer-based packaging used to connect HBM memory to GPU dies in the H100.

**Key topics:** monolithic die vs chiplets, yield as motivation, AMD chiplet history, die-to-die interconnects (UCIe, Infinity Fabric, NVLink, UltraFusion), 2.5D and 3D stacking, HBM memory, TSMC CoWoS and SoIC, Intel Foveros.

---

### [Chapter 61: Testing and Quality — How We Know a Chip Works](61-testing-and-quality.md)

A fabricated wafer goes through an extensive testing process before any chip reaches a customer. Wafer probe testing: while still on the wafer, each die is probed electrically to identify defective dies (marked with ink) before the wafer is cut. Die cutting (singulation): the wafer is cut into individual dies using a diamond saw. Packaging: each good die is bonded into a protective package (BGA, LGA, QFN) with electrical connections exposed. Package testing: final comprehensive electrical testing at speed in an Automated Test Equipment (ATE) machine — this can last hours per chip. Binning: chips that fail full specification but pass lower specs are sold at lower clock speeds (this is why there are faster and slower variants of "the same" chip — Core i5 vs Core i7 are often the same die, just binned differently).

**Key topics:** wafer probe testing, die yield, singulation, chip packaging (BGA, LGA), ATE testing, burn-in testing, binning, quality and reliability standards, what "defects per million" means in practice.

---

---

# PART 7 — THE SHAKTI PROCESSOR

---

## Volume 8: India's Open Source CPU

> In 2011 a team at IIT Madras began a project that would become India's most significant contribution to processor architecture: SHAKTI, a complete family of open-source, RISC-V-based processors designed entirely in India. This volume tells the story of SHAKTI — what it is, how it works, why it matters, and what it means for India's technological future.

### [Chapter 62: The IIT Madras SHAKTI Program — Origins and Vision](62-shakti-origins-and-vision.md)

India imports nearly 100% of the semiconductors its economy depends on. Every phone, laptop, server, and defense system runs on chips designed abroad. In 2011, a group at IIT Madras led by Professor V. Kamakoti set out to change this: to design a family of processors entirely in India, using fully open-source tools, and release them freely for anyone to use. The SHAKTI program chose RISC-V as its ISA — a deliberate decision to avoid licensing fees and to contribute to the global open-source hardware movement. The name SHAKTI (शक्ति) means power or energy in Sanskrit. The Indian government, through MeitY (Ministry of Electronics and Information Technology), has funded SHAKTI as a strategic technology initiative.

**Key topics:** motivation for indigenous processor design, IIT Madras RISE Lab, Professor V. Kamakoti, RISC-V choice, MeitY funding, strategic importance of semiconductors, SHAKTI program goals.

---

### [Chapter 63: The SHAKTI Architecture — Design Philosophy](63-shakti-architecture.md)

SHAKTI uses the RISC-V ISA as its foundation — specifically the RV32 (32-bit) and RV64 (64-bit) base integer ISA with extensions. The design philosophy: keep it open, keep it verifiable, and prioritize security alongside performance. SHAKTI implements the full privilege architecture (machine, supervisor, user modes), enabling it to run Linux and real operating systems. The implementation is written in Bluespec SystemVerilog (BSV) — a high-level HDL developed at MIT that allows hardware design at a higher abstraction level than traditional Verilog. BSV's type system catches design errors early. The resulting RTL is then synthesized to gates and taped out. SHAKTI's open-source repository on GitHub contains the complete source code of every processor variant.

**Key topics:** RISC-V privilege architecture, Bluespec SystemVerilog (BSV), why BSV over Verilog, open-source RTL, privilege modes (M/S/U), SHAKTI repository structure, formal verification approach.

---

### [Chapter 64: SHAKTI Processor Families — From Embedded to Server](64-shakti-processor-families.md)

SHAKTI is not one processor but a family, covering the full range from tiny embedded microcontrollers to high-performance server chips:

**C-class (Copper):** The embedded-class processor. 32-bit, in-order, 3-stage pipeline. Comparable to ARM Cortex-M3. Target: IoT devices, microcontrollers, sensor nodes.

**E-class (Electron):** 32-bit, in-order, 5-stage pipeline. Higher performance than C-class. Comparable to ARM Cortex-M4. Target: embedded systems requiring more compute.

**I-class (Iron):** 64-bit, in-order, 5-stage pipeline. Boots Linux. Comparable to ARM Cortex-A5. Target: mid-range embedded Linux systems.

**M-class (Mercury):** 64-bit, out-of-order, superscalar. High-performance application processor. Comparable to ARM Cortex-A53. Target: smartphones, tablets.

**H-class (Helium):** 64-bit, highly-out-of-order, multi-issue. Server-class performance. Target: servers, HPC.

**S-class (Sulfur):** Multi-core server class with hardware security extensions.

---

### [Chapter 65: SHAKTI's Security Extensions — Security as a First-Class Citizen](65-shakti-security-extensions.md)

Traditional processor design adds security as an afterthought. SHAKTI takes a different approach: security features are designed in from the start. SHAKTI implements hardware support for Tagged Memory — every memory location carries a small tag indicating its security classification. The hardware enforces access rules automatically, preventing software vulnerabilities like buffer overflows and use-after-free bugs at the hardware level. SHAKTI also implements Physical Memory Protection (PMP) from the RISC-V specification, which restricts which memory regions each privilege level can access. The SHAKTI IOMMU (I/O Memory Management Unit) controls which devices can access which memory, preventing DMA-based attacks. These features make SHAKTI particularly interesting for defense, aerospace, and critical infrastructure applications.

**Key topics:** tagged memory, capability-based security, RISC-V PMP, SHAKTI IOMMU, hardware-enforced memory safety, comparison with ARM TrustZone, why hardware security matters for India's strategic applications.

---

### [Chapter 66: RISC-V Ecosystem and SHAKTI's Contributions](66-risc-v-ecosystem-shakti-contributions.md)

SHAKTI does not just use RISC-V — it actively contributes to the ecosystem. The IIT Madras team has developed open-source tooling, FPGA prototyping platforms, and documentation that benefits the entire global RISC-V community. SHAKTI on FPGAs: the complete SHAKTI C-class boots Linux on a Xilinx Artix-7 FPGA and is freely downloadable. The SHAKTI SDK: a complete software development environment for SHAKTI-based systems. SHAKTI's participation in the RISC-V International Foundation and contributions to ISA extensions. The Shakti platform board: a development board with SHAKTI SoC, RAM, and peripherals that Indian researchers and students can use to study real processor design.

**Key topics:** SHAKTI FPGA implementations, SHAKTI SDK, SHAKTI development boards, RISC-V International participation, open-source contributions, academic impact.

---

### [Chapter 67: SHAKTI Silicon Tapeouts — Going from Simulation to Real Chips](67-shakti-silicon-tapeouts.md)

SHAKTI has moved beyond simulation into actual silicon. The first SHAKTI tapeout: a 32nm chip fabricated at SCL (Semiconductor Laboratory, Chandigarh) — India's government-owned semiconductor fab. The chip successfully booted Linux. Subsequent tapeouts at TSMC's advanced nodes for higher-performance variants. India's Semiconductor Mission (ISM): the Indian government's plan to establish semiconductor fabrication in India, with $10 billion in incentives. Foxconn-HCL, Micron, and Tata Electronics have announced fab investments in India. The strategic context: India wants to move up the value chain from just designing chips to manufacturing them domestically.

**Key topics:** SCL Chandigarh fab, SHAKTI silicon tapeouts, India Semiconductor Mission, India Semiconductor Fab investments, Micron Gujarat ATMP, Tata Electronics fab Dholera, India's 2047 semiconductor goals.

---

### [Chapter 68: India's Semiconductor Ecosystem — The Road Ahead](68-india-semiconductor-ecosystem.md)

India's semiconductor story extends beyond SHAKTI. The design talent pipeline: India has over 20,000 chip design engineers — the largest pool outside the United States. Companies like Texas Instruments, Intel, Qualcomm, and AMD run major R&D centers in Bengaluru and Hyderabad. IIT and NIT graduates feed these centers. C-DAC (Centre for Development of Advanced Computing) develops India's own high-performance computers. The PARAM supercomputer series. DRDO develops processors for defense applications. ISRO designs chips for satellite and launch vehicle applications. The gap: India has design talent but lacks large-scale fabrication — the India Semiconductor Mission aims to close this gap by 2028. The opportunity: as geopolitical tensions push supply chain diversification, India is positioned to become a significant semiconductor hub.

**Key topics:** India's chip design ecosystem, TI/Intel/AMD India R&D centers, IIT/NIT talent pipeline, C-DAC, PARAM supercomputers, DRDO semiconductor programs, ISRO chip design, India Semiconductor Mission timeline, geopolitical opportunity.

---

---

# PART 8 — THE FUTURE

---

## Volume 9: What Comes Next

> The era of simply shrinking transistors to make chips faster is ending. The next decades of computing will be defined by new materials, new architectures, new ways of building chips, and entirely new computing paradigms. This final volume surveys the frontier.

### [Chapter 69: Multicore and Manycore — Parallel by Design](69-multicore-and-manycore.md)

When single-core clock speeds stopped scaling in 2004 (due to power and heat), the industry turned to parallelism: put multiple cores on one chip. Dual-core, quad-core, octa-core, and beyond. A 128-core AMD EPYC server CPU. A 4,096-core GPU. Parallelism is free in hardware; the challenge is in software. Amdahl's Law: if 5% of a program cannot be parallelized, the maximum speedup from any number of parallel cores is 20x — no matter how many cores you add. Thread-level parallelism (TLP): running multiple program threads simultaneously. NUMA (Non-Uniform Memory Access): in multi-socket servers, each CPU socket has local RAM that is faster to access than the other socket's RAM.

**Key topics:** multicore history, Amdahl's Law, Gustafson's Law, symmetric vs asymmetric multiprocessing, NUMA, cache coherence at scale, memory bandwidth wall, programming for parallelism.

---

### [Chapter 70: Heterogeneous Computing — The Right Tool for Each Job](70-heterogeneous-computing.md)

A modern laptop has a CPU (good at serial tasks), an integrated GPU (good at parallel graphics and AI), a neural engine (good at ML inference), a signal processor (good at audio), and maybe a dedicated encoder chip (good at video). This is heterogeneous computing: using multiple different processor types to handle different workloads efficiently. The challenge: moving data between these processors. Apple's unified memory architecture avoids this: CPU, GPU, and Neural Engine all share the same physical RAM. AMD's APUs (Accelerated Processing Units) combine CPU and GPU on the same die with shared cache. The future direction: every SoC will have many specialized accelerators, each consuming tiny power for its specific task.

**Key topics:** heterogeneous architectures, data movement bottleneck, unified memory, AMD APU, ARM big.LITTLE, Intel's P+E core hybrid, OpenCL/SYCL for heterogeneous programming, the "disaggregated" data center.

---

### [Chapter 71: Memory-Centric Computing — Moving Compute to Data](71-memory-centric-computing.md)

The biggest bottleneck in modern computing is not the speed of the CPU — it is moving data between where it is stored and where it is processed. For a GPU training a neural network, the data movement consumes more energy than the actual computation. Processing-In-Memory (PIM): place simple compute units directly inside the DRAM chips, so operations can happen where the data lives. Samsung's HBM-PIM: HBM memory modules with integrated MAC units, used for AI. UPMEM: a commercial PIM product with 2,500 small processors embedded in DRAM modules. Near-Data Processing (NDP): compute units very close to memory, even if not inside it. Storage-class memory: technologies like Intel Optane (3D XPoint) that bridge the gap between DRAM speed and SSD density.

**Key topics:** memory bandwidth wall, processing-in-memory (PIM), Samsung HBM-PIM, UPMEM, near-data processing, compute express link (CXL) memory expansion, storage class memory, implications for future architectures.

---

### [Chapter 72: The Post-Moore Era — What Replaces Shrinking?](72-the-post-moore-era.md)

Moore's Law gave the industry a free lunch for 50 years: every two years, the same chip got twice as fast for the same price. That free lunch is ending. The alternatives that will drive the next era of computing improvement: Specialization (domain-specific architectures — build the exact right hardware for AI, video, genomics). Advanced packaging (chiplets, 3D stacking — connect more compute in one package). New materials (gallium nitride GaN, silicon carbide SiC for power electronics; indium gallium arsenide InGaAs for RF; carbon nanotubes for logic). Architectural innovation (machine learning-guided microarchitecture, computing in memory). Quantum computing for specific hard problems. The end of Dennard scaling does not mean the end of progress — it means progress becomes harder and more differentiated.

**Key topics:** Dennard scaling end, the "free lunch is over," DSAs (Domain-Specific Architectures), architectural diversity, GaN/SiC power electronics, carbon nanotube transistors, IBM 2nm GAA, the next 20 years.

---

### [Chapter 73: 3D Integration and Advanced Packaging](73-3d-integration-and-advanced-packaging.md)

Chiplets spread compute horizontally; 3D stacking goes vertical. TSMC's SoIC (System on Integrated Chips): bond multiple dies face-to-face with copper-to-copper micro-bumps at 9µm pitch — far denser than any board-level connection. Intel's Foveros: 3D stacking with hybrid bonding, allowing the compute die on top and the base die (with I/O) below. The dream: stack DRAM directly on top of a CPU die — eliminating the entire off-chip memory bus. SK Hynix, Samsung, and Micron are all developing "cube" 3D DRAM. The challenge: heat removal. A stacked die cannot dissipate heat as easily as a single die. Thermal management innovations: microfluidic cooling channels inside the chip package itself.

**Key topics:** hybrid bonding vs micro-bumps, TSMC SoIC, Intel Foveros, 3D DRAM stacking, thermal limits of 3D integration, active cooling in package, wafer-on-wafer vs die-on-wafer bonding.

---

### [Chapter 74: Domain-Specific Architectures — The New Specialization](74-domain-specific-architectures.md)

When general-purpose processors hit scaling limits, the answer is to build processors that are really, really good at one specific domain. David Patterson (co-inventor of RISC) now argues that DSAs (Domain-Specific Architectures) are the most important architecture trend of the next decade. Examples: Google TPU (matrix multiply for neural networks), Cerebras CS-2 (entire neural network layer fits on chip — 850,000 AI cores), Graphcore IPU (graph-structured neural network computation), Groq LPU (sequential LLM token generation). Genomics accelerators (Oxford Nanopore), video encoding ASICs (Apple VideoToolbox), database query accelerators (Oracle SPARC M-class). The common pattern: the application domain has a specific mathematical structure; the DSA exploits that structure to achieve 10-1000x efficiency gain over a general-purpose processor.

**Key topics:** Patterson's DSA argument, TPU/Cerebras/Graphcore/Groq architectures, genomics ASICs, the cost of flexibility, DSA design principles, when to build a DSA vs use a GPU.

---

### [Chapter 75: The Next Frontier — What Are We Building Toward?](75-the-next-frontier.md)

Fifty years ago, a computer filled a room and was owned by universities and governments. Today you carry a trillion-transistor processor in your pocket. Where do the next fifty years lead? Near-term (5-10 years): 2nm and below GAA transistors, widespread chiplet adoption, AI integrated into every chip, the first practical quantum error correction. Medium-term (10-20 years): processing-in-memory mainstream, brain-computer interfaces moving from research to clinical use, 1 trillion transistors on a single chip, quantum computers solving industrially relevant optimization problems. Long-term speculation: biological computing using DNA as storage (a gram of DNA can store 215 petabytes), neuromorphic chips that approach the energy efficiency of biological brains (the human brain consumes ~20W and runs at 86 billion neurons), post-silicon substrates for computing elements. The consistent thread: human curiosity drives the next breakthrough, just as it always has.

**Key topics:** near-term roadmap (2025-2030), ITRS/IEEE roadmap, brain-computer interfaces, DNA data storage, trillion-transistor chips, the carbon nanotube CPU, post-silicon, and a call to the reader to join in building what comes next.

---

## A Note on the Journey

You began this course with a transistor — a switch made of sand. You end it at the frontier of what human ingenuity is reaching for: processors with a trillion switches, quantum computers that harness physics itself, and open-source chips designed in India that anyone on earth can build.

The transistor has not changed. It is still a switch. Everything you saw in between — the logic gates, the adders, the pipeline stages, the cache hierarchies, the billions of dollars of fabrication equipment — is the result of humans finding cleverer and cleverer ways to connect those switches.

Now you see how it works. You see why your laptop gets warm under heavy load (the switching energy becomes heat). You see why your phone battery drains when the GPU works hard. You see why Apple Silicon feels so fast (unified memory, tight hardware-software integration). You see why RISC-V matters (open hardware democratizes chip design the way Linux democratized software).

And if you are an engineer, a student, or a builder in India — you now see why SHAKTI matters, and what it means that a team at IIT Madras chose to spend ten years giving the world a free, open, capable processor family. That is the kind of contribution this field rewards and needs.

The next chapter is yours to write.

---

*Computer Architecture — From Sand and Electricity to Designing Your Own Processor*
