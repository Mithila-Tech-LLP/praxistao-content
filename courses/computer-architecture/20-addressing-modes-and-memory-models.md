# Chapter 20: Addressing Modes and Memory Models

How does a CPU compute the address of a memory location? The answer — "addressing modes" — reveals a deep design choice in every ISA. And once we have an address, how does the system ensure that multiple programs and multiple processors see a consistent view of memory? The "memory model" answers this. This chapter covers both: the mechanics of address computation and the subtle rules of memory consistency.

## Table of Contents

1. [Addressing Modes in Depth](#1-addressing-modes-in-depth)
2. [Addressing in RISC-V vs x86](#2-addressing-in-risc-v-vs-x86)
3. [Memory Alignment](#3-memory-alignment)
4. [Endianness](#4-endianness)
5. [Memory Models: Coherence and Consistency](#5-memory-models-coherence-and-consistency)
6. [Cache Coherence Protocols](#6-cache-coherence-protocols)
7. [Memory Ordering: Sequential Consistency](#7-memory-ordering-sequential-consistency)
8. [Relaxed Memory Models and Fences](#8-relaxed-memory-models-and-fences)
9. [Practical Implications for Programmers](#9-practical-implications-for-programmers)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Addressing Modes in Depth

An **addressing mode** is the method by which an instruction computes the memory address of an operand. The address computed from the instruction fields is called the **effective address** (EA).

### Register Direct

No memory access — the operand is the register value itself.

```
EA = none (the value IS in the register)
Example: ADD R1, R2, R3    (RISC-V: add t0, t1, t2)
```

### Immediate (Literal)

The operand is a constant embedded in the instruction.

```
EA = none
Example: ADDI R1, R2, 5   (add 5 to R2, store in R1)
```

### Absolute (Direct)

The instruction contains the full memory address.

```
EA = address_field
Example: LOAD R1, [0x12345678]    (rare in RISC — addresses don't fit in instruction)
```

This mode is impractical in RISC architectures where instructions are 32 bits — a 32-bit address leaves no room for opcode and register fields. x86 supports it for some instructions.

### Register Indirect

The register holds the memory address.

```
EA = reg
Example: LOAD R1, [R2]   (RISC-V: lw t0, 0(t1) with offset=0)
```

### Base + Displacement (Base + Offset)

The most important and most used addressing mode. The effective address is the sum of a base register and a signed constant (displacement).

```
EA = base_register + displacement
Example: LW t0, 12(a0)   (effective address = a0 + 12)
```

Real-world uses:
- **Struct field access**: base = pointer to struct, displacement = field offset
- **Array access**: base = array pointer, displacement = index × element_size (if element_size = immediate-friendly)
- **Stack access**: base = sp (stack pointer), displacement = local variable offset

```c
struct Point { float x; float y; float z; };
void zero_z(struct Point* p) { p->z = 0.0f; }
```

```riscv
# p in a0, z is at offset 8 (x at 0, y at 4, z at 8)
sw zero, 8(a0)      # store 0.0f (integer 0 = float 0.0 bit pattern) to p->z
```

### Base + Index (Scaled)

EA = base_register + index_register × scale_factor

This mode is particularly important for array access where the index is variable:
```c
int arr[100]; int i = 42;
int val = arr[i];    // effective address = &arr + i*4
```

x86 supports 4 scale factors: 1, 2, 4, 8 (byte, halfword, word, doubleword).

```x86asm
mov eax, [rbx + rcx*4]    ; eax = *(rbx + rcx*4)
```

RISC-V does NOT have scaled index. The programmer/compiler must compute the address:
```riscv
slli t1, t1, 2        # t1 = i * 4 (multiply by 4 = scale)
add  t2, a0, t1       # t2 = &arr[i]
lw   t0, 0(t2)        # t0 = arr[i]
```

### PC-Relative

EA = PC + offset. The effective address is relative to the current program counter.

Why this matters: code that uses PC-relative addressing can be loaded at any memory address (position-independent code). Shared libraries (.so files, .dylib files) MUST use PC-relative addressing because they can be loaded at different addresses in different processes.

```riscv
auipc t0, %pcrel_hi(data_label)     # t0 = PC + upper 20 bits of offset
addi  t0, t0, %pcrel_lo(data_label) # t0 = address of data_label
```

In x86-64, RIP-relative addressing was added specifically for this:
```x86asm
lea rax, [rip + label]    ; rax = address of label (PC-relative)
```

### Quick Check

> 1. What is an "effective address"?
> 2. Why can't RISC-V use absolute addressing for 64-bit addresses?
> 3. Why must shared libraries use PC-relative addressing?

---

## 2. Addressing in RISC-V vs x86

The contrast between RISC-V and x86 addressing modes illustrates the RISC vs CISC philosophy:

### RISC-V Addressing

RISC-V has exactly ONE memory addressing mode for load/store: **base + 12-bit signed offset**.

```
EA = rs1 + imm12   (where imm12 is a 12-bit signed immediate: -2048 to +2047)
```

That's it. If you need a different addressing mode, you compute the address in registers first, then use the computed address as the base with offset=0.

```riscv
# Array with variable index:
slli t1, t1, 2          # t1 = index * 4
add  t2, a0, t1         # t2 = &arr[index]
lw   t0, 0(t2)          # load arr[index]

# Struct nested inside another struct:
# outer.inner.field at offset outer.inner=20, field=8 → total offset=28
lw   t0, 28(a0)         # one instruction if offset < 2048

# Very large offset:
lui  t1, %hi(BIGOFFSET)
addi t1, t1, %lo(BIGOFFSET)
add  t1, a0, t1         # compute full address
lw   t0, 0(t1)          # load
```

### x86 Addressing

x86 supports an extremely rich set of addressing modes. The full formula:

```
EA = base_register + index_register × scale + displacement
     where scale ∈ {1, 2, 4, 8}
           displacement is 0, 8-bit, or 32-bit signed
           base_register: any GPR
           index_register: any GPR except RSP
```

Examples:
```x86asm
[rbx]                    # base only
[rbx + 8]                # base + 8-bit displacement
[rbx + 0x12345678]       # base + 32-bit displacement
[rbx + rcx]              # base + index
[rbx + rcx*4]            # base + scaled index
[rbx + rcx*4 + 8]        # base + scaled index + displacement
[rip + offset]           # RIP-relative (for PIC code)
```

The x86 encoded addressing mode uses a complex **ModRM** byte + optional **SIB** (Scale-Index-Base) byte to encode this variety. This complexity requires a dedicated address generation unit (AGU) and adds decode overhead — but eliminates explicit address computation instructions for the common cases.

### Tradeoff Summary

| Mode | RISC-V | x86 |
|------|--------|-----|
| base only | `lw t0, 0(a0)` | `mov eax, [rbx]` |
| base+offset | `lw t0, 8(a0)` | `mov eax, [rbx+8]` |
| base+index | add+lw (2 instr) | `mov eax, [rbx+rcx]` |
| base+scaled index | slli+add+lw (3 instr) | `mov eax, [rbx+rcx*4]` |
| Decode complexity | Simple | Complex (ModRM+SIB) |
| Hardware area | Small | Larger AGU |

### Quick Check

> 1. How many addressing modes does RISC-V have for loads and stores?
> 2. Write the RISC-V sequence to load from address `base + index * 8`.
> 3. What do the ModRM and SIB bytes encode in x86?

---

## 3. Memory Alignment

**Alignment** is the requirement that memory accesses start at addresses that are multiples of the access size.

### Why Alignment Matters

Consider loading a 32-bit (4-byte) integer from address 0x1001:

```
Memory:  byte at 0x1001  byte at 0x1002  byte at 0x1003  byte at 0x1004
                ↑ start of unaligned 4-byte read

The 4 bytes span two naturally-aligned 4-byte blocks:
Block 0x1000: [0x1000, 0x1001, 0x1002, 0x1003]   (reads 3 bytes from here)
Block 0x1004: [0x1004, 0x1005, 0x1006, 0x1007]   (reads 1 byte from here)
```

Hardware that naturally fetches data in aligned chunks must:
1. Read the first aligned block
2. Read the second aligned block  
3. Combine the relevant bytes

This is called an **unaligned access**. It requires either:
- Hardware support (Intel x86 handles unaligned accesses transparently but slower)
- A hardware exception + software fixup (MIPS, early ARM would trap on unaligned access)
- Undefined behavior (some ISAs say: don't do it)

### RISC-V Alignment

The base RISC-V spec says: aligned accesses are guaranteed to work. Unaligned accesses may work (hardware supports them) or may trap (hardware raises an exception, software handles it). The Zicclsm extension adds fast unaligned hardware support.

In practice, Linux on RISC-V handles unaligned access traps in software — this is correct but slow (100x slower than aligned access).

### Alignment Rules for C Structs

```c
struct Example {
    char  a;      // 1 byte at offset 0
    // 3 bytes PADDING here (to align int to 4-byte boundary)
    int   b;      // 4 bytes at offset 4
    char  c;      // 1 byte at offset 8
    // 1 byte PADDING (to align short to 2-byte boundary)
    short d;      // 2 bytes at offset 10
    // 0 bytes padding (total size = 12, multiple of largest member = 4)
};
// sizeof(struct Example) = 12 (not 8!)
```

Padding is added silently by the compiler. This is why `sizeof` of a struct often surprises programmers.

**Optimization**: reorder fields from largest to smallest to minimize padding:
```c
struct Optimized {
    int   b;      // 4 bytes at offset 0
    short d;      // 2 bytes at offset 4
    char  a;      // 1 byte at offset 6
    char  c;      // 1 byte at offset 7
    // 0 bytes padding!
};
// sizeof(struct Optimized) = 8  (vs 12 above — 33% smaller!)
```

### Quick Check

> 1. What is an "unaligned memory access" and why is it problematic?
> 2. Why does `struct { char a; int b; }` not have size 5 bytes?
> 3. How would you reorder the fields of `struct { char a; int b; char c; double d; }` to minimize struct size?

---

## 4. Endianness

**Endianness** determines how multi-byte values are stored in memory — which byte goes first?

### Little-Endian

The **least significant byte** (LSB) is stored at the **lowest address**.

Value `0x12345678` stored at address `0x1000`:
```
Address: 0x1000  0x1001  0x1002  0x1003
Byte:      0x78    0x56    0x34    0x12
           LSB                     MSB
```

x86, x86-64, ARM (default mode), RISC-V, and most modern ISAs are little-endian.

### Big-Endian

The **most significant byte** (MSB) is stored at the **lowest address**.

Value `0x12345678` stored at address `0x1000`:
```
Address: 0x1000  0x1001  0x1002  0x1003
Byte:      0x12    0x34    0x56    0x78
           MSB                     LSB
```

Network byte order, SPARC, PowerPC (optional), MIPS (optional), and early ARM (optional) are big-endian.

### Why Does Endianness Matter?

1. **Network communication**: TCP/IP uses network byte order (big-endian). Code that sends integers over a network must convert with `htons()`/`htonl()` (host-to-network short/long) and convert back on the other end.

2. **File formats**: Binary file formats specify endianness. Reading a PowerPC-format file on x86 requires byte-swapping.

3. **Cross-architecture debugging**: When examining memory dumps from different architectures, bytes appear in different orders.

4. **Exploits**: Buffer overflow exploits must use the target's endianness when encoding addresses.

### Swapping Endianness

```riscv
# RISC-V: swap endianness of a 32-bit value in t0
# (Straightforward without a hardware bswap instruction)
# t0 = 0xAABBCCDD → want 0xDDCCBBAA

andi  t1, t0, 0xFF            # t1 = 0x000000DD (byte 0)
srli  t2, t0, 8
andi  t2, t2, 0xFF            # t2 = 0x000000CC (byte 1)
srli  t3, t0, 16
andi  t3, t3, 0xFF            # t3 = 0x000000BB (byte 2)
srli  t4, t0, 24              # t4 = 0x000000AA (byte 3)

slli  t1, t1, 24              # t1 = 0xDD000000
slli  t2, t2, 16              # t2 = 0x00CC0000
slli  t3, t3, 8               # t3 = 0x0000BB00
or    t0, t1, t2
or    t0, t0, t3
or    t0, t0, t4              # t0 = 0xDDCCBBAA
```

x86 has a one-instruction `bswap eax` that does the same in 1 cycle.

### Quick Check

> 1. What does "little-endian" mean? Where is the least significant byte stored?
> 2. Why must network code convert integer byte order?
> 3. If you read bytes at address 0x1000-0x1003 and get {0x01, 0x02, 0x03, 0x04}, what 32-bit integer do they represent in little-endian? In big-endian?

---

## 5. Memory Models: Coherence and Consistency

When multiple processors (or multiple cores) all access the same memory, we need rules about what they see. There are two related concepts:

### Cache Coherence

**Problem**: if cores A and B both have a copy of memory location X in their L1 caches, and core A writes X, what does core B see?

**Cache coherence** ensures that all cores see the same value for a memory location. Any write by one core eventually becomes visible to all other cores.

### Memory Consistency

**Problem**: even with coherence, in what ORDER do writes from one core become visible to other cores?

**Memory consistency** defines the ordering guarantees of memory operations.

### Quick Check

> 1. What is the difference between cache coherence and memory consistency?
> 2. If core A writes value 42 to address X while core B reads X from its cache, what does coherence require?

---

## 6. Cache Coherence Protocols

The most widely used coherence protocol family is **MESI** (Modified-Exclusive-Shared-Invalid):

### MESI States

Each cache line can be in one of four states:

```
State     Meaning
─────     ─────────────────────────────────────────────────────
M (Modified)   This cache has the only copy, and it's been modified.
               Main memory is OUT OF DATE. This core must write back
               before another core can access this line.

E (Exclusive)  This cache has the only copy, and it matches main memory.
               Can upgrade to M without bus transaction.

S (Shared)     This line may be in other caches too. Read-only.
               If this core wants to write, it must invalidate others.

I (Invalid)    This cache line doesn't have valid data.
               Must fetch from memory or another cache on access.
```

### MESI State Transitions

```
          Core reads, no other cache has it      
    I ─────────────────────────────────────→ E
                                              |
    E ──── Core writes ─────────────────────→ M
                                              |
    M ──── Another core reads ──────────────→ S (after writeback)
    E ──── Another core reads ──────────────→ S
    S ──── Core writes ──────────────────────→ M (after invalidating others)
    S ──── Another core invalidates ─────────→ I
    M ──── Bus transaction ───────────────────→ I (write back first)
```

### How Invalidation Works

```
Core 0 has X=5 in Shared state.
Core 1 has X=5 in Shared state.

Core 0 writes X=42:
  Core 0 broadcasts "INVALIDATE X" on the coherence bus
  Core 1 receives the message, transitions X to Invalid state
  Core 0 transitions X to Modified state, writes 42 to its L1 cache

Core 1 now reads X:
  Cache miss (X is Invalid in Core 1's cache)
  Core 1 requests X
  Core 0's cache controller detects it has Modified copy of X
  Core 0 writes back X=42 to memory
  Core 1 gets X=42 from memory (or directly from Core 0's cache)
  Both go to Shared state
```

### MOESI, MESIF Variants

- **MOESI**: Adds Owned state — allows dirty (modified) cache lines to be read by other cores without writing back to memory first (faster).
- **MESIF**: Adds Forward state — one cache designated to respond to read requests, reducing bus traffic.

Intel uses MESIF. AMD uses MOESI.

### Quick Check

> 1. What does each letter in MESI stand for?
> 2. When a cache line is in Modified state, what must happen before another core can read it?
> 3. What is "invalidation" and why is it necessary?

---

## 7. Memory Ordering: Sequential Consistency

### The Ideal: Sequential Consistency (SC)

**Lamport's Sequential Consistency (1979)**: The result of any execution is the same as if the operations of all processors were executed in some sequential order, and the operations of each individual processor appear in this sequence in the order specified by its program.

In simpler terms: the memory appears as a single shared memory where reads and writes happen in some total order that respects each thread's program order.

Sequential consistency is the intuitive model most programmers assume:

```
Thread 1:          Thread 2:
x = 1              y = 1
r1 = y             r2 = x
```

Under SC, if `r1 = 0` (Thread 1 read y before Thread 2 wrote it), then `r2 = 1` must hold (Thread 2 read x after Thread 1 wrote it). It's impossible for both `r1 = 0` AND `r2 = 0` under SC.

### The Problem: SC is Expensive

Implementing SC requires:
1. Completing all previous memory operations before starting the next one
2. No reordering of memory operations by the processor

Both of these kill performance. Modern out-of-order processors reorder memory operations all the time — a store might complete before an earlier load, and loads from the same address might be forwarded from the store buffer without waiting for it to commit to cache.

Enforcing SC would require flushing the store buffer after every store, serializing all memory operations. Performance would drop by 40-80%.

### Quick Check

> 1. What does "sequential consistency" mean for a multi-processor system?
> 2. Why is sequential consistency expensive to implement?

---

## 8. Relaxed Memory Models and Fences

### Relaxed Memory Models

Real processors implement **relaxed** memory models that allow certain reorderings:

**Total Store Order (TSO)** — used by x86:
- Stores may be visible to other processors AFTER younger loads from the same processor
- A processor can read its own stores before they are globally visible
- But stores are globally visible in program order

Under TSO, the `r1=0, r2=0` outcome IS possible:
- Thread 1 buffers `x=1` in its store buffer
- Thread 1 reads `y` → gets 0 (Thread 2's write hasn't been visible yet)
- Thread 2 buffers `y=1` in its store buffer
- Thread 2 reads `x` → gets 0 (Thread 1's write hasn't been visible yet)
- Both store buffers drain
- Result: r1=0, r2=0

**ARM's and RISC-V's Relaxed Model** — even weaker than TSO:
- Loads and stores can be reordered with respect to each other
- An earlier store might become visible to other processors after a later load from the same thread

### Memory Fences (Barriers)

To prevent unwanted reordering, ISAs provide **fence** instructions that enforce ordering:

```riscv
# RISC-V fence:
fence  iorw, iorw     # full fence: predecessor and successor cannot be reordered
fence.i               # instruction fence: ensures instruction cache coherence
# Components: i=instruction fetch, o=device output, r=read (load), w=write (store)
# fence rw, rw is most common for data ordering
```

```x86asm
# x86 memory fences:
mfence        # full memory fence (sequentially consistent for all memory operations)
sfence        # store fence (all stores before this are visible before stores after)
lfence        # load fence (all loads before this complete before loads after)
```

```arm
# AArch64 barriers:
DMB ISH       # Data Memory Barrier, Inner Shareable domain — full data ordering
DSB ISH       # Data Synchronization Barrier (stronger than DMB)
ISB           # Instruction Synchronization Barrier — flushes pipeline
```

### When You Need Fences: The Classic Example

```c
// Thread 1                    // Thread 2
data = 42;                     while (ready == 0) {} // spin
                               use(data);  // safe?
// fence: STORE fence
ready = 1;
```

Without the fence, the CPU might reorder `ready = 1` before `data = 42` (even though they're in separate addresses). Thread 2 would see `ready = 1` but read uninitialized `data`.

With the fence, `data = 42` is guaranteed to be globally visible before `ready = 1` is visible.

### Higher-Level Synchronization

In practice, you don't write fences directly. You use:
- **Mutexes/locks**: contain the appropriate fences internally
- **Atomic operations**: `std::atomic<int>` with `memory_order_acquire`/`memory_order_release`
- **Message passing**: sending data to another thread automatically orders the memory

The C++11 memory model and POSIX threads define exactly which operations create which memory orderings, so you don't need to insert bare fences in most application code.

---

## 9. Practical Implications for Programmers

### Data Races are Undefined Behavior in C/C++

If two threads access the same memory location and at least one access is a write, and there is no synchronization between them — this is a **data race**. In C++, data races are undefined behavior: the compiler can generate code that does anything, including incorrect results, crashes, or security vulnerabilities.

Use `std::atomic<T>` or a mutex for any shared mutable data.

### The "Happens-Before" Relationship

Correct concurrent programs establish "happens-before" relationships:
- Operation A "happens before" operation B if they're in the same thread (program order)
- A lock release "happens before" the subsequent lock acquire
- A fence "happens before" operations after it on other threads

If A happens-before B, any write in A is visible in B.

### False Sharing

Two threads modifying different variables that happen to be on the same cache line causes **false sharing** — cache coherence traffic even though the threads are logically independent:

```c
struct __attribute__((aligned(64))) {
    // Put these on SEPARATE cache lines to avoid false sharing!
    int counter_thread0;
    char padding[60];     // pad to fill 64-byte cache line
    int counter_thread1;
} shared;
```

When Thread 0 writes `counter_thread0`, the cache coherence protocol invalidates Thread 1's copy of the cache line — including `counter_thread1`. Thread 1 gets a cache miss even though it only accesses `counter_thread1`. Adding padding eliminates this.

### Quick Check

> 1. What is a "data race" and why is it undefined behavior in C++?
> 2. What is "false sharing" and how does it hurt performance?
> 3. When should a programmer use `std::atomic<T>` vs a mutex?

---

## Summary

- **Addressing modes** specify how to compute the effective address: register, immediate, absolute, register indirect, base+displacement, base+index (scaled), PC-relative.
- RISC-V has ONE addressing mode: base + 12-bit offset. x86 has the full range including scaled index.
- **Alignment**: data at address that is a multiple of its size (4-byte int → 4-byte-aligned). Unaligned access is slow or traps. Compiler adds padding to structs.
- **Endianness**: little-endian (LSB at lowest address, used by x86/ARM/RISC-V) vs big-endian (MSB first, used by network byte order). Matters for cross-platform binary data.
- **Cache coherence** (MESI/MOESI/MESIF) ensures all cores see a consistent value for each memory location, using invalidation and write-back protocols.
- **Memory consistency models**: SC (sequential consistency) is ideal but expensive. TSO (x86) allows store buffering. RISC-V and ARM use fully relaxed models.
- **Fences** enforce ordering: `fence rw, rw` (RISC-V), `mfence` (x86), `DMB` (ARM).
- **Data races** are undefined behavior. Use `std::atomic` or mutexes for shared mutable state. False sharing occurs when independent variables share a cache line.

---

## Exercises

### Easy

1. Compute the effective address for each:
   - RISC-V: `lw t0, 16(a0)` where a0 = 0x1000
   - x86: `mov eax, [rbx + rcx*4]` where rbx = 0x2000, rcx = 5
   - PC-relative: `auipc t0, 0` where PC = 0x10000

2. What is the size of `struct { char a; int b; char c; }` in C on a 32-bit system? Draw the memory layout including padding.

3. Value `0xDEADBEEF` stored at address 0x1000 in little-endian. What byte is at 0x1000? At 0x1003?

### Medium

4. **Cache coherence trace**: Two cores with 2-entry direct-mapped MESI caches. Each cache line holds 1 word (4 bytes). Initial state: both caches empty.

   Trace these operations and show each cache line's state after each step:
   - Core 0 reads address 0x100 (value = 5)
   - Core 1 reads address 0x100
   - Core 0 writes 42 to address 0x100
   - Core 1 reads address 0x100
   
   How many cache coherence bus transactions occur total?

5. **Memory model puzzle**: Under the x86 TSO memory model, can the following outcome occur? Explain:
   ```
   Thread 1: store x=1, load y into r1
   Thread 2: store y=1, load x into r2
   
   Result: r1=0, r2=0 (both read 0)
   ```
   Under what memory ordering would this be possible? Under what would it be impossible?

6. **False sharing exercise**: A parallel counter increment:
   ```c
   int counters[4];  // shared array
   // Thread i increments counters[i] 1,000,000 times
   ```
   On a 64-byte cache line, `counters[0]` through `counters[3]` (4×4=16 bytes) all fit in ONE cache line. Estimate the performance impact of false sharing vs a version where each counter is on its own cache line. How would you fix it in C?

### Hard

7. **Lock-free programming**: Implement a thread-safe single-producer single-consumer ring buffer in C using RISC-V atomic operations (lr.w/sc.w). The producer writes to the buffer; the consumer reads from it. What memory ordering is required for the write to `data` to be visible before the write to `head`? Write out the fence instructions needed. Then implement the same using C11 `stdatomic.h` with appropriate `memory_order` arguments.

8. **Spectre via cache timing**: The Spectre variant 1 attack uses the cache as a side-channel. The attacker code:
   ```c
   if (index < array1_size) {
       temp = array2[array1[index] * 64];
   }
   ```
   The CPU speculatively executes the if-body even when `index >= array1_size` (out-of-bounds). The speculative access to `array2` leaves a cache timing trace. Explain: (a) how an attacker chooses `index` to read an arbitrary kernel address via `array1[index]`, (b) how timing `array2` access reveals what was read, (c) what does the RISC-V memory fence specification say about speculation past a fence? Does `fence` prevent Spectre? (d) What software mitigation does the GCC/LLVM compiler insert (look up `array_index_nospec`)?
