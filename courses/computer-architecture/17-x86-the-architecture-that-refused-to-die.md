# Chapter 17: x86 — The Architecture That Refused to Die

No ISA in history has defied more predictions of its demise than x86. Declared "too complex" by academics in 1980, "dying" when RISC arrived in 1985, "obsolete" when 64-bit was needed in 1999, "doomed" when ARM arrived in mobile in 2010 — x86 survived every challenge and still runs the world's most powerful servers and desktops. This chapter covers the full x86 family: its origin, evolution, extensions, and the extraordinary engineering that keeps it competitive.

## Table of Contents

1. [The IBM PC Accident](#1-the-ibm-pc-accident)
2. [x86 Evolution: 8086 to x86-64](#2-x86-evolution-8086-to-x86-64)
3. [x86-64 Registers and Encoding](#3-x86-64-registers-and-encoding)
4. [x86-64 Instruction Set](#4-x86-64-instruction-set)
5. [Extensions That Changed Computing](#5-extensions-that-changed-computing)
6. [x86's Ugly Parts](#6-x86s-ugly-parts)
7. [How Modern Intel/AMD Make x86 Fast](#7-how-modern-intelamd-make-x86-fast)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The IBM PC Accident

### Intel 8086 (1978)

Intel designed the 8086 as an upgrade from the 8080 — a 16-bit processor for embedded systems and industrial control. It was not designed to be the foundation of personal computing.

The 8086:
- 16-bit data bus, 20-bit address bus (could address 1 MB)
- 8 general-purpose registers (AX, BX, CX, DX, SI, DI, BP, SP)
- Many instructions operating on specific registers (AX for multiply, CX for loops, DX for I/O)
- Variable-length instructions (1-6 bytes)
- 29,000 transistors at 5 MHz

### The IBM PC Choice (1980)

IBM was developing a personal computer and needed an operating system and processor quickly. IBM chose Intel 8088 (a crippled 8086 with an 8-bit external bus to save cost) and licensed DOS from a tiny company called Microsoft. The choice was made partly because Intel had an established sales relationship with IBM.

Had IBM chosen the Motorola 68000 (a far cleaner 32-bit architecture popular with engineers), the history of computing would be unrecognizable. But they didn't.

The IBM PC became the dominant personal computer standard by 1983. Every clone maker (Compaq, Dell, HP, Acer) had to be x86-compatible to run IBM PC software. x86 lock-in was born.

### Intel 80286 (1982) and 80386 (1985)

The 286 added **protected mode** — memory protection between programs, essential for multitasking operating systems. The 386 added:
- Full 32-bit processing (IA-32, also called i386)
- 32-bit address bus (4 GB address space)
- Paging and virtual memory
- Backward compatibility with all 16-bit 8086 software

Every new processor maintained backward compatibility. The 386 could run code from the 8086 (1978) unchanged. This commitment to compatibility became a defining characteristic — and a burden — of x86.

### Quick Check

> 1. What accident of history made x86 the dominant PC architecture?
> 2. What processor was the main alternative to x86 for personal computers in 1980?
> 3. What major feature did the 80386 add that made it suitable for modern operating systems?

---

## 2. x86 Evolution: 8086 to x86-64

The x86 family grew by accretion — each generation added features without breaking old ones:

```
Year  Processor   Bits  Key Addition
────  ──────────  ────  ────────────────────────────────────────────────
1978  Intel 8086   16   Original x86 ISA
1982  Intel 80286  16   Protected mode (segmentation, memory protection)
1985  Intel 80386  32   Full 32-bit (IA-32), paging, 4 GB address space
1989  Intel 80486  32   Built-in FPU, 8KB L1 cache on chip, pipeline
1993  Pentium      32   Dual issue (2 instructions/cycle), 64-bit data bus
1995  Pentium Pro  32   P6 µop architecture (RISC core inside CISC)
1997  Pentium MMX  32   MMX extension (SIMD for multimedia)
1999  Pentium III  32   SSE (Streaming SIMD: 128-bit XMM registers, 4×float)
2000  Pentium 4    32   NetBurst: very deep pipeline (20-31 stages), SSE2
2003  AMD Opteron  64   x86-64 / AMD64 (64-bit extension to x86)
2006  Core 2       64   Intel adopts AMD64, macro-fusion, SSSE3
2011  Sandy Bridge 64   AVX (256-bit SIMD), ring bus, integrated GPU
2013  Haswell      64   AVX2, FMA (fused multiply-add), Transactional Memory
2017  Skylake-X    64   AVX-512 (512-bit SIMD)
2019  Ice Lake     64   Sunny Cove: better IPC, faster AVX-512
2021  Alder Lake   64   Hybrid: P-cores + E-cores (Intel Thread Director)
2023  Meteor Lake  64   Tile-based chiplet design, NPU for AI
```

### The 64-bit Transition (2003)

By 2003, the 32-bit address space (4 GB) was becoming a bottleneck for servers and memory-hungry applications. Two approaches existed:

**Intel's approach: IA-64 / Itanium** — a completely new, incompatible 64-bit ISA using VLIW (Very Long Instruction Word). Programs had to be recompiled for IA-64. Itanium launched in 2001. Despite technical merits, it was a market disaster — nobody wanted to abandon their x86 software.

**AMD's approach: AMD64 (x86-64)** — add a 64-bit extension to x86 that runs existing 32-bit code unchanged. AMD Opteron launched in 2003. Intel quickly adopted AMD's extension (calling it Intel 64 or EM64T). x86-64 became the standard.

The lesson: binary compatibility always beats technical purity in the market.

### Quick Check

> 1. Which company created the 64-bit extension to x86? Why did this approach win over Intel's IA-64?
> 2. What were the two main reasons the 32-bit address space became a bottleneck by 2003?
> 3. In what year did Intel add the floating-point unit (FPU) directly onto the processor chip?

---

## 3. x86-64 Registers and Encoding

### General-Purpose Registers

x86-64 has 16 general-purpose registers, each 64 bits wide. Historical naming is chaotic because each was extended over generations:

```
64-bit  32-bit  16-bit  8-bit high  8-bit low   Role
──────  ──────  ──────  ──────────  ─────────   ─────────────────────────
rax     eax     ax      ah          al           Accumulator / return value
rbx     ebx     bx      bh          bl           Base (general / callee-saved)
rcx     ecx     cx      ch          cl           Counter (loop / shift amount)
rdx     edx     dx      dh          dl           Data (multiply / I/O)
rsi     esi     si                  sil          Source Index
rdi     edi     di                  dil          Destination Index
rsp     esp     sp                  spl          Stack Pointer
rbp     ebp     bp                  bpl          Base/Frame Pointer
r8      r8d     r8w                 r8b          (AMD64 addition)
r9      r9d     r9w                 r9b
r10     r10d    r10w                r10b
r11     r11d    r11w                r11b
r12     r12d    r12w                r12b
r13     r13d    r13w                r13b
r14     r14d    r14w                r14b
r15     r15d    r15w                r15b
```

Writing to a 32-bit register (e.g., `eax`) zeroes the upper 32 bits of the 64-bit register. Writing to 8-bit or 16-bit registers does NOT zero the upper bits (historical quirk for backward compatibility).

### Special Registers

```
rip        Instruction pointer (PC), 64-bit
rflags     Flags register: CF (carry), ZF (zero), SF (sign), OF (overflow), PF (parity)
```

### SIMD Register Files

x86-64 has accumulated multiple overlapping SIMD register sets:
```
MMX:  mm0-mm7    64-bit (aliased to FPU x87 registers — terrible design)
SSE:  xmm0-xmm15  128-bit
AVX:  ymm0-ymm15  256-bit (xmm is the lower 128 bits of ymm)
AVX-512: zmm0-zmm31  512-bit (ymm is the lower 256 bits of zmm)
```

### x86-64 Instruction Encoding

x86-64 instructions are 1-15 bytes. The encoding has many optional fields:

```
Optional   Optional  Optional  Optional  Optional  Required  Optional  Optional
REX Prefix  Prefix    Prefix   Opcode   ModRM     SIB       Disp      Imm
1 byte      1 byte    1 byte  1-3 bytes  1 byte   1 byte   0-4 bytes 0-4 bytes
```

The `REX` prefix (added in x86-64) allows access to the 8 new registers (r8-r15) and 64-bit operation. Without it, x86-64 code looks like 32-bit code.

This complexity (compared to RISC-V's clean 4-byte encoding) requires the x86 decode frontend to be significantly more complex.

### Quick Check

> 1. How many general-purpose 64-bit registers does x86-64 have?
> 2. What happens to the upper 32 bits of `rax` when you write to `eax`?
> 3. What is the maximum size (in bytes) of a single x86-64 instruction?

---

## 4. x86-64 Instruction Set

### System V AMD64 ABI (Linux calling convention)

```
Arguments:  rdi, rsi, rdx, rcx, r8, r9  (first 6 integer args)
            xmm0-xmm7                     (first 8 float args)
Return:     rax (integer), xmm0 (float)
Caller-saved: rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11
Callee-saved: rbx, rbp, r12, r13, r14, r15
Stack: 16-byte aligned before call instruction
```

### Key Instructions

```x86asm
# Data movement
mov  rax, rbx          # rax = rbx
mov  rax, [rbx]        # rax = Memory[rbx]  (load)
mov  [rbx], rax        # Memory[rbx] = rax  (store)
mov  rax, [rbx+rcx*4]  # indexed: rax = Memory[rbx + rcx*4]
movzx rax, byte [rbx]  # load 1 byte, zero-extend
movsx rax, dword [rbx] # load 4 bytes, sign-extend
lea  rax, [rbx+rcx*4]  # load effective address (no memory access!)
push rbp               # rsp -= 8; Memory[rsp] = rbp
pop  rbp               # rbp = Memory[rsp]; rsp += 8

# Arithmetic
add  rax, rbx          # rax += rbx
sub  rax, rbx          # rax -= rbx
imul rax, rbx          # rax = rax * rbx (signed)
idiv rbx               # rax = rdx:rax / rbx; rdx = rdx:rax % rbx
inc  rax               # rax++
dec  rax               # rax--
neg  rax               # rax = -rax

# Logic
and  rax, rbx
or   rax, rbx
xor  rax, rbx          # also: xor rax, rax zeros rax efficiently
not  rax               # bitwise NOT
shl  rax, 3            # logical left shift by 3
shr  rax, 3            # logical right shift
sar  rax, 3            # arithmetic right shift (preserves sign)

# Comparison and flags
cmp  rax, rbx          # set flags based on rax - rbx (discards result)
test rax, rbx          # set flags based on rax AND rbx (discards result)

# Conditional jumps (read flags set by cmp/test)
je   label   # jump if equal (ZF=1)
jne  label   # jump if not equal
jl   label   # jump if less (signed)
jle  label   # jump if less or equal
jg   label   # jump if greater
jge  label   # jump if greater or equal
jb   label   # jump if below (unsigned)
ja   label   # jump if above (unsigned)
jz   label   # jump if zero (= je)
jnz  label   # jump if not zero (= jne)

# Function call/return
call func    # push rip+size; jmp func
ret          # pop rip; jmp rip
```

### Sum 1 to N in x86-64

```x86asm
# int64_t sum_to_n(int64_t n)
# rdi = n; return value in rax

sum_to_n:
    xor  eax, eax          # rax = 0 (sum); using eax zeros upper 32 bits
    test rdi, rdi
    jle  .done             # if n <= 0, return 0
    mov  ecx, 1            # ecx = 1 (i)
.loop:
    add  rax, rcx          # sum += i
    inc  ecx               # i++
    cmp  ecx, edi          # compare i to n
    jle  .loop             # if i <= n, continue
.done:
    ret
```

### Quick Check

> 1. In the System V AMD64 ABI, which register holds the first function argument?
> 2. What does `lea rax, [rbx+rcx*4]` do? How is it different from `mov rax, [rbx+rcx*4]`?
> 3. What is the purpose of `xor rax, rax` and why is it preferred over `mov rax, 0`?

---

## 5. Extensions That Changed Computing

### MMX (1996)

The first x86 SIMD extension. 8 × 64-bit MMX registers (mm0-mm7), aliased to x87 FPU registers (a catastrophic design decision — you couldn't use MMX and floating-point simultaneously without saving/restoring context). Used for integer SIMD on video data. Superseded by SSE.

### SSE/SSE2 (1999/2001)

Streaming SIMD Extensions. Introduced separate 128-bit XMM registers (xmm0-xmm15). SSE added:
- 4 × 32-bit float packed operations (4 floats processed simultaneously)
- New register file (finally separate from x87)

SSE2 added:
- 2 × 64-bit float packed
- Integer SIMD (16 × 8-bit, 8 × 16-bit, 4 × 32-bit, 2 × 64-bit)

SSE2 became the baseline for 64-bit x86 — every x86-64 CPU must support SSE2.

### AVX/AVX2 (2011/2013)

Advanced Vector Extensions. Extended XMM to 256-bit YMM registers. AVX2 added:
- Integer SIMD at 256-bit width
- Gather instructions (load from non-contiguous memory locations)
- FMA (Fused Multiply-Add): `result = a*b + c` in a single instruction (critical for matrix math)

A key compiler trick: compile the same function twice, once for SSE2 (baseline) and once for AVX2. Use CPUID at runtime to detect support and call the faster version. This is "function multi-versioning."

### AVX-512 (2017)

Extended to 512-bit ZMM registers. 32 registers (zmm0-zmm31) vs 16 in AVX. New mask registers for conditional (predicated) SIMD. Extremely powerful for HPC and ML workloads.

Controversy: AVX-512 instructions generate so much heat that Intel's Alder Lake (2021) and later processors DOWN-CLOCK the CPU when AVX-512 executes — the core temperature spikes. This is why some Intel chips dropped AVX-512 support entirely (Apple M1 has none).

### AES-NI (2010)

Hardware acceleration for AES encryption/decryption. What takes 60-100 cycles in software takes 1-3 cycles with AES-NI. All modern x86 and ARM chips include AES-NI.

```x86asm
# One round of AES encryption in one instruction:
aesenc  xmm0, xmm1     # one AES round
aesenclast xmm0, xmm2  # final AES round
```

### Quick Check

> 1. Why was the MMX register aliasing to x87 FPU registers a bad design?
> 2. How many single-precision floats can AVX2 process simultaneously?
> 3. Why does AVX-512 cause clock frequency to drop on some processors?

---

## 6. x86's Ugly Parts

### Implicit Register Usage

Many instructions write implicit registers without the programmer specifying:
- `MUL rbx` — multiplies rax by rbx, puts result in rdx:rax (128-bit result)
- `DIV rbx` — divides rdx:rax by rbx, quotient in rax, remainder in rdx
- `LOOP label` — decrements rcx and branches if rcx ≠ 0 (ecx dependency)
- `PUSH rax` — implicitly modifies rsp
- `CPUID` — writes to eax, ebx, ecx, edx simultaneously

These implicit dependencies cause pipeline hazards because the CPU must track hidden register writes.

### Segment Registers (Legacy Baggage)

x86 inherited 16-bit segment registers: CS (code), DS (data), ES (extra), FS (thread-local), GS (another thread-local), SS (stack). In 64-bit mode, most segments are ignored (base = 0, limit = full address space), but FS and GS are still used (FS:0 = thread-local storage in Linux; GS:0 in Windows). The segment register mechanism adds complexity to every memory access.

### FLAGS Register Dependencies

The flags register creates a serial dependency chain. Consider:

```x86asm
cmp  rax, rbx    # writes ZF, CF, SF, OF, AF, PF
je   .label      # reads ZF
sete al          # also reads ZF — all three are linked!
```

Between any two flag-using instructions, the CPU must preserve the flag state. This prevents many optimizations that would be trivial in RISC-V (which has no flags register).

### CISC Memory Operands and Decode Complexity

The x86 decode front-end must handle instructions like:
```x86asm
add qword [rax+rbx*4+16], rcx    # read: 64-bit value at address (rax+rbx*4+16), add rcx, write back
```

Parsing this single instruction requires:
1. Decode the opcode (1-3 bytes)
2. Parse the ModRM byte
3. Parse the SIB byte (for scaled index)
4. Read the 32-bit displacement
5. Compute the effective address: rax + rbx×4 + 16
6. Issue a load
7. Issue an add
8. Issue a store
9. Track that the flags are modified

This decodes to 3-5 µops and takes multiple decode cycles. The x86 frontend is one of the most complex pieces of hardware on the chip, consuming ~15% of die area and significant power.

### Quick Check

> 1. What implicit registers does the `DIV` instruction use?
> 2. Why do segment registers still exist in 64-bit mode even though most are unused?
> 3. Why does the FLAGS register create pipeline problems?

---

## 7. How Modern Intel/AMD Make x86 Fast

Despite all its legacy complexity, modern x86 processors are extraordinarily fast. How?

### µop Translation and Fusion

The decode frontend translates x86 instructions to µops. Modern Intel processors additionally perform **macro-fusion**: combining pairs of instructions into a single µop:

```x86asm
cmp  rax, rbx       →  combined into single CMP+JE µop
je   .target
```

And **micro-fusion**: combining a memory operand with an operation:
```x86asm
add  rax, [rbx]     →  single µop (load + add fused)
```

These fusions reduce the µop count, allowing more instructions to flow through the pipeline per cycle.

### Decoded ICache (µop Cache)

Intel since Sandy Bridge (2011) has a **µop cache** (also called "Decoded ICache" or L0 cache). After decoding x86 instructions to µops, the µops are cached. On a cache hit, the frontend fetches µops directly, bypassing the x86 decoder entirely. This saves the decode complexity cost for repeated code (loops).

### Out-of-Order Execution

Modern CPUs don't execute instructions in the order they appear in the program. The out-of-order engine (Chapter 22) finds independent instructions and executes them simultaneously. An Intel Skylake core can have 224 µops "in flight" simultaneously — only 5 of them were being "executed" in the original program order at any given moment.

### Deep Cache Hierarchy

L1: 32-64KB, 4-5 cycles latency
L2: 256-512KB, 10-15 cycles
L3: 4-64MB (shared), 35-50 cycles
DRAM: 16-64 GB, 100-300 cycles

The L1 and L2 caches can service most memory accesses without hitting slow DRAM.

### Quick Check

> 1. What is macro-fusion? Give an example.
> 2. What is the µop cache and how does it help performance?
> 3. What does "out-of-order execution" mean? How many µops can an Intel Skylake have in flight?

---

## Summary

- x86 was chosen for the IBM PC by accident (1981); ISA lock-in made it dominant.
- x86 grew from 8086 (16-bit, 1978) through i386 (32-bit, 1985) to x86-64 (64-bit, 2003). AMD created x86-64; Intel adopted it after Itanium failed.
- x86-64 has 16 × 64-bit GPRs, multiple overlapping SIMD register files (XMM/YMM/ZMM), and the FLAGS register.
- Key extensions: SSE/SSE2 (128-bit SIMD), AVX/AVX2 (256-bit), AVX-512 (512-bit), AES-NI, FMA.
- x86's ugly parts: implicit register use, segment registers, FLAGS dependencies, complex variable-length encoding.
- Modern Intel/AMD hide complexity: µop translation, macro/micro-fusion, µop cache, out-of-order execution, deep cache hierarchy.

---

## Exercises

### Easy

1. What year was the Intel 8086 introduced, and what was its design purpose?

2. List three x86-64 SIMD register types and their widths (in bits).

3. What does the `LEA` instruction do differently from `MOV` with a memory operand?

### Medium

4. **Register naming exercise**: x86-64 registers are awkwardly named for historical reasons. Given that `rax` is 64-bit:
   - Name the 32-bit, 16-bit, and 8-bit views of `rax`
   - If you write `0xABCDEF1234567890` to `rax`, then execute `mov bl, al`, what value is in `bl`?
   - What value is in `eax` after executing `mov eax, 1`? What is in `rax`?

5. The System V AMD64 ABI allows only 6 integer arguments in registers (rdi, rsi, rdx, rcx, r8, r9). What happens if a function needs 10 arguments? Where are arguments 7-10 passed? Write a small x86-64 example showing how a caller sets up a 7-argument function call.

6. **Explain the Itanium failure**: Itanium (IA-64) used VLIW (Very Long Instruction Word) where the compiler explicitly schedules instructions for parallel execution, encoding instruction-level parallelism statically. Compare this to Intel's x86-64 out-of-order execution engine that dynamically finds instruction-level parallelism at runtime. What are the trade-offs? Why did VLIW fail commercially despite being technically elegant?

### Hard

7. **SPECTRE and ISA design**: The Spectre vulnerability (2018) exploits speculative execution in out-of-order processors. When the CPU speculatively executes past a branch (before knowing the branch outcome), it may access memory it shouldn't — and although the speculative results are discarded, the cache state is not fully reset, leaking information via timing side channels. This is a fundamental conflict: the ISA says "execute instruction B only after verifying condition A," but the microarchitecture executes B speculatively before A is confirmed. Research and answer: (a) Can the x86 ISA itself be patched to prevent Spectre? (b) What hardware mitigations exist (IBRS, STIBP, SSBD)? (c) What performance penalty do these mitigations impose? (d) Does RISC-V's simpler ISA make it immune to Spectre? Why or why not?

8. **Compiler explorer deep-dive**: Use the Godbolt compiler explorer (godbolt.org) with x86-64 GCC. Compile this code with `-O3 -mavx2`:
   ```c
   void add_arrays(float* a, float* b, float* c, int n) {
       for (int i = 0; i < n; i++) c[i] = a[i] + b[i];
   }
   ```
   (a) What AVX2 instructions does the compiler generate?
   (b) How many floats are processed per loop iteration with vectorization?
   (c) What does the loop prologue (before the main loop) do? Why?
   (d) What does the loop epilogue handle?
   (e) Compile again with `-O3` but without `-mavx2`. How many floats per iteration now? What is the performance ratio?
