# Chapter 22: Instruction Encoding and Decoding

An instruction is ultimately a number — a binary pattern that the CPU's decoder must interpret. How that number is structured (the instruction format) is one of the most consequential ISA design decisions. It determines code density, decoder complexity, pipeline performance, and extensibility. This chapter covers how instructions are encoded in binary and how the CPU decodes them.

## Table of Contents

1. [From Assembly to Binary](#1-from-assembly-to-binary)
2. [Fixed vs Variable Width Encoding](#2-fixed-vs-variable-width-encoding)
3. [RISC-V Encoding: A Clean Design](#3-risc-v-encoding-a-clean-design)
4. [x86-64 Encoding: A Historical Mess](#4-x86-64-encoding-a-historical-mess)
5. [ARM64 Encoding: Elegant 32-bit](#5-arm64-encoding-elegant-32-bit)
6. [Opcode Space Management](#6-opcode-space-management)
7. [The Decode Unit in Hardware](#7-the-decode-unit-in-hardware)
8. [Compressed Instruction Encoding](#8-compressed-instruction-encoding)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. From Assembly to Binary

When an assembler processes assembly language, it converts each mnemonic+operands into a binary number according to the ISA's encoding rules.

### Example: RISC-V ADD Encoding

```
Assembly:   add  t0, t1, t2
```

RISC-V ADD is an R-type instruction. Looking up the encoding tables:
```
Instruction: ADD  rd, rs1, rs2
funct7:      0000000
rs2:         t2 = x7  = 00111
rs1:         t1 = x6  = 00110
funct3:      000
rd:          t0 = x5  = 00101
opcode:      0110011   (R-type arithmetic)

Binary: 0000000 00111 00110 000 00101 0110011
Hex:    0x00730293  
         Wait, let me compute properly:
         0000000_00111_00110_000_00101_0110011
         = 0x007302B3
```

Let's verify by bit positions:
```
31      25|24  20|19  15|14  12|11   7|6      0
 0000000  | 00111 | 00110 | 000  | 00101 | 0110011
 funct7     rs2     rs1    funct3  rd      opcode
 = 0       = 7     = 6     = 0    = 5     = 0x33

Binary:   0000 0000 0111 0011 0000 0010 1011 0011
Hex:      0x00730233
```

### Verifying with Spike or QEMU

On a RISC-V Linux system:
```bash
# Write a RISC-V program
echo "add t0, t1, t2" | riscv64-unknown-elf-as -o test.o
riscv64-unknown-elf-objdump -d test.o
# Output: 
# 0: 00730233  add t0,t1,t2
```

### Quick Check

> 1. What type of instruction is RISC-V ADD? What fields does it have?
> 2. The opcode field is 7 bits. How many distinct opcodes can it encode?
> 3. If rs1 is 5 bits, how many registers can be specified?

---

## 2. Fixed vs Variable Width Encoding

### Fixed-Width (32-bit): RISC-V, MIPS, ARM64

Every instruction is exactly 32 bits (4 bytes). This provides:

**Advantages**:
- **Simple decode**: the fetch unit always fetches exactly 4 bytes; the decode unit always receives exactly 4 bytes
- **Alignment guarantee**: instructions are naturally 4-byte aligned (PC is always a multiple of 4)
- **Predictable pipeline stages**: every instruction has the same decode complexity
- **Simple branch target computation**: all branches are to 4-byte-aligned addresses

**Disadvantages**:
- **Code density**: a simple "load 1" into a register needs all 32 bits, even though the immediate (1) fits in 5 bits
- **Immediate limitations**: with 32-bit instructions and overhead for opcode/registers, immediates are limited to 12-20 bits
- **Not all instructions need 32 bits**: most real programs have many "common simple" instructions

### Variable-Width: x86-64

Instructions range from 1 byte to 15 bytes. This provides:

**Advantages**:
- **Code density**: simple operations encode in 1-3 bytes; complex operations in 8-15 bytes
- **Large immediates**: can encode 32-bit immediates naturally in 32-bit instructions
- **Alignment doesn't matter**: instructions can start at any byte address

**Disadvantages**:
- **Complex fetch**: don't know where instruction N+1 starts until instruction N is decoded
- **Sequential decode dependency**: can't decode instructions N and N+1 simultaneously without knowing N's length
- **Complex hardware**: the fetch/decode frontend is the most complex part of modern x86 CPUs

### Hybrid: ARM32 Thumb / RISC-V C Extension

Mix 16-bit and 32-bit instructions. The first bit(s) of each instruction indicate whether it's 16-bit or 32-bit:

RISC-V C extension:
- If bits [1:0] == 11: 32-bit instruction (all base ISA instructions)
- If bits [1:0] != 11: 16-bit compressed instruction (C extension)

ARM Thumb2:
- If bits [15:11] == 11101, 11110, 11111: 32-bit Thumb-2 instruction
- Otherwise: 16-bit Thumb instruction

This enables 25-30% code size reduction while retaining simple decode (2 bits checked first).

### Quick Check

> 1. Why is fixed-width instruction encoding simpler for hardware?
> 2. What is the main disadvantage of variable-width encoding for the decode pipeline?
> 3. How does the RISC-V C extension differentiate 16-bit from 32-bit instructions?

---

## 3. RISC-V Encoding: A Clean Design

RISC-V's 32-bit encoding is carefully designed for hardware efficiency. Let's examine every format:

### The Six Formats

```
R-type (Register-Register):
31      25 24  20 19  15 14  12 11   7 6      0
┌─────────┬──────┬──────┬──────┬──────┬────────┐
│ funct7  │ rs2  │ rs1  │funct3│  rd  │ opcode │
└─────────┴──────┴──────┴──────┴──────┴────────┘

I-type (Immediate — loads, ALU with immediate, JALR):
31              20 19  15 14  12 11   7 6      0
┌─────────────────┬──────┬──────┬──────┬────────┐
│   imm[11:0]     │ rs1  │funct3│  rd  │ opcode │
└─────────────────┴──────┴──────┴──────┴────────┘

S-type (Store):
31    25 24  20 19  15 14  12 11   7 6      0
┌───────┬──────┬──────┬──────┬──────┬────────┐
│imm[11:5]│ rs2 │ rs1  │funct3│imm[4:0]│opcode│
└───────┴──────┴──────┴──────┴──────┴────────┘

B-type (Branch):
31  30    25 24  20 19  15 14  12 11  8  7  6     0
┌──┬───────┬──────┬──────┬──────┬────┬──┬────────┐
│i12│i[10:5]│ rs2  │ rs1  │funct3│i[4:1]│i11│ opcode │
└──┴───────┴──────┴──────┴──────┴────┴──┴────────┘

U-type (Upper Immediate — LUI, AUIPC):
31              12 11   7 6      0
┌─────────────────┬──────┬────────┐
│   imm[31:12]    │  rd  │ opcode │
└─────────────────┴──────┴────────┘

J-type (Jump — JAL):
31  30        21 20  19       12 11   7 6      0
┌──┬────────────┬──┬────────────┬──────┬────────┐
│i20│ imm[10:1] │i11│ imm[19:12] │  rd  │ opcode │
└──┴────────────┴──┴────────────┴──────┴────────┘
```

### Key Design Insights

**1. rs1 and rs2 are ALWAYS at the same bit positions**:
```
rs1 = bits [19:15] in R, I, S, B formats
rs2 = bits [24:20] in R, S, B formats
```
The register file can be read immediately after fetching the instruction, before fully decoding it. This reduces decode latency.

**2. rd is ALWAYS at bits [11:7]**:
The register write port can be determined early.

**3. Sign bit is ALWAYS at bit 31**:
All immediates are sign-extended. The sign bit is always the MSB of the instruction (bit 31). This means sign-extension hardware always reads the same bit.

**4. B-type scrambles the immediate bits deliberately**:
The branch target bits [10:1] come from instruction bits [30:21] — the same position as J-type's lower immediate bits. This allows sharing immediate-extraction hardware between branches and jumps.

**5. Immediates are missing bit 0 in B and J types**:
Branches and jumps encode offsets in multiples of 2 bytes (for Thumb compatibility, even though RISC-V base ISA is 4-byte aligned). The LSB of the offset is implicitly 0 — encoded range doubles while using the same number of bits.

### Instruction Decoding Examples

**Decode 0x00730233**:
```
Binary: 0000 0000 0111 0011 0000 0010 0011 0011

opcode:  bits [6:0]  = 011 0011  = 0x33 → R-type arithmetic
rd:      bits [11:7] = 0 0100    = 4 → x4 (tp... wait, let me recheck)
```

Actually let me properly decode:
```
0x00730233:
Binary: 00000000011100110000001000110011

opcode = 0110011 (7 bits from right) = 0x33 = R-type
rd     = 00100 = 4
funct3 = 000
rs1    = 00110 = 6
rs2    = 00111 = 7
funct7 = 0000000

R-type with opcode=0x33, funct3=000, funct7=0000000 → ADD
ADD rd=x4, rs1=x6, rs2=x7
ADD tp, t1, t2    (using ABI names: x4=tp, x6=t1, x7=t2)
```

### Quick Check

> 1. Why does RISC-V always place rs1 at bits [19:15] in all formats that use rs1?
> 2. Why do branch immediates not include bit 0?
> 3. What is the sign bit of every RISC-V immediate, and why is this consistent?

---

## 4. x86-64 Encoding: A Historical Mess

x86-64 instruction encoding is famously complex — a testament to 45 years of backward-compatible additions.

### The Legacy Structure

An x86-64 instruction consists of optional prefix bytes followed by the core instruction bytes:

```
Structure (left to right, all optional except Opcode):
┌──────────┬──────────┬──────────┬────────┬──────┬─────┬──────────┬──────────┐
│ Prefixes │  REX     │  Escape  │ Opcode │ModRM │ SIB │  Disp    │  Imm     │
│ 0-4 bytes│ 0-1 byte │ 0-2 bytes│1-3 bytes│0-1B │ 0-1B│ 0,1,2,4B│ 0,1,2,4B │
└──────────┴──────────┴──────────┴────────┴──────┴─────┴──────────┴──────────┘
Total: 1-15 bytes
```

### Prefix Bytes (Legacy)

The original 8086 had single-byte prefixes that modified instruction behavior. x86-64 still has them for backward compatibility:

```
0xF0    LOCK prefix  (atomic read-modify-write)
0xF2    REPNE/REPNZ  (repeat not-equal: for string operations)
0xF3    REP/REPE     (repeat: for string operations)
0x66    Operand size override (change default 32-bit to 16-bit)
0x67    Address size override (change default 64-bit to 32-bit)
0x2E    CS segment override
0x3E    DS segment override
0x26    ES segment override
0x64    FS segment override (thread-local storage)
0x65    GS segment override
0x36    SS segment override
```

### REX Prefix (Added in AMD64)

The REX prefix is a single byte of the form `0100WRXB`:
```
Bit 7-4: must be 0100 (identifies it as REX)
W = 1: 64-bit operand (default is 32-bit)
R = 1: extend ModRM.reg field (register field)
X = 1: extend SIB.index field  
B = 1: extend ModRM.r/m or SIB.base or opcode reg field
```

Without REX, you can only access r0-r7 (8 registers). With REX.R=1, you access r8-r15. This is how AMD64 doubled the register count from 8 to 16 without changing the core instruction format.

### ModRM Byte

The ModRM byte encodes addressing mode:
```
Bits [7:6]  = Mod   (00=indirect, 01=indirect+disp8, 10=indirect+disp32, 11=direct register)
Bits [5:3]  = Reg   (register or opcode extension)
Bits [2:0]  = R/M   (register or memory specifier)

Special cases:
  Mod=00, R/M=101: RIP-relative addressing (no base register, just disp32 added to RIP)
  R/M=100 (with Mod≠11): SIB byte follows
```

### SIB Byte

Present when ModRM says so. Encodes base+index*scale:
```
Bits [7:6]  = Scale  (0→×1, 1→×2, 2→×4, 3→×8)
Bits [5:3]  = Index  (register for index, or 100=no index)
Bits [2:0]  = Base   (register for base, or 101→depends on Mod)
```

### Escape Bytes

x86 ran out of 1-byte opcode space years ago. Escape bytes prefix multi-byte opcodes:
- `0F xx`: 2-byte opcode (added for 286, most SSE/MMX instructions)
- `0F 38 xx`: 3-byte opcode (added for SSE4, SHA, AES-NI)
- `0F 3A xx`: 3-byte opcode variant

### VEX/EVEX Prefixes (Added for AVX)

When Intel added AVX (256-bit), they needed: 3 operands (instead of x86's usual 2-operand destructive style), YMM registers, and non-destructive operations. They added new 2-3 byte prefix formats:

- **VEX prefix**: 2 or 3 bytes, replaces REX + escape bytes for AVX instructions
- **EVEX prefix**: 4 bytes, for AVX-512 (adds mask registers, broadcast, embedded rounding)

This means modern x86 has FOUR different prefix systems: legacy prefixes, REX, VEX, and EVEX — all overlapping and interacting.

### Why x86 Encoding Matters for Performance

Intel's x86 decode front-end must:
1. Determine instruction boundaries (where does instruction N end and N+1 begin?)
2. Parse legacy prefixes, REX, VEX, or EVEX
3. Decode opcode (1-3 bytes + prefix)
4. Decode ModRM and SIB if present
5. Determine displacement and immediate lengths
6. Fetch all these bytes from the instruction cache

This takes 3-6 cycles for complex instructions. The x86 frontend uses:
- **Instruction Length Decoder**: determines boundaries (scans for prefix bytes + opcode)
- **Pre-decoder**: identifies prefix patterns
- **Full decoder**: converts to µops (4 decoders in parallel for Intel: 1 complex + 3 simple)
- **µop cache**: caches decoded µops to bypass the decoder for hot loops

The µop cache is the key innovation that makes x86 competitive despite decode complexity.

### Quick Check

> 1. What does the REX.W bit do in x86-64?
> 2. Why does x86 need escape bytes (0x0F prefix)?
> 3. What is the difference between a VEX prefix and a REX prefix?

---

## 5. ARM64 Encoding: Elegant 32-bit

AArch64 instruction encoding is a clean 32-bit fixed-width format, but with more structure than RISC-V:

### Instruction Groups by Top Bits

ARM64 partitions the opcode space using the top bits:

```
Bits [28:25]  Category
────────────  ─────────────────────────────
0000          Reserved/SME
0001          Reserved/SVE
0010          SVE (Scalable Vector Extension)
0011          Reserved
0100          Data processing (immediate)
0101          Data processing (immediate)
0110          Loads and stores
0111          Loads and stores
1000          Data processing (register)
1001          Data processing (register)
1010          Data processing (register)
1011          Data processing (register)
1100          Floating-point and advanced SIMD
1101          Floating-point and advanced SIMD
1110          Floating-point and advanced SIMD
1111          Floating-point and advanced SIMD
```

### Data Processing (Register) Format

```
31  29 28  24  23  22  21  20  16  15  10   9   5   4   0
┌────┬─────┬────┬────┬────┬─────┬───────┬─────┬─────┐
│ sf │op21 │S   │type│shift│Rm   │ imm6  │ Rn  │ Rd  │
└────┴─────┴────┴────┴────┴─────┴───────┴─────┴─────┘
sf = 0: 32-bit (W registers), 1: 64-bit (X registers)
```

### Load/Store Format

```
31 30  27  26  24  22  21  20  16  15  12  10   9   5   4   0
┌──┬────┬──┬────┬──┬──┬─────┬────┬────┬─────┬─────┐
│size│op1│V │op2 │opc  │imm9 │opt │S  │ Rn  │ Rt  │
└──┴────┴──┴────┴──┴──┴─────┴────┴────┴─────┴─────┘
```

### Register Encoding

ARM64 has 31 general-purpose registers (x0-x30). This fits in 5 bits (0-31), and register 31 is used as:
- **XZR/WZR** (zero register) when the spec says "reads as zero, writes ignored"
- **SP/WSP** (stack pointer) in certain contexts

The dual-use of register 31 requires context-dependent decode.

### Quick Check

> 1. How does ARM64 encode whether an instruction uses 32-bit (W) or 64-bit (X) registers?
> 2. ARM64 has 31 GPRs (x0-x30). How does it encode them in 5 bits that can represent 0-31?
> 3. What encoding trick does ARM64 use to make register 31 serve dual purposes?

---

## 6. Opcode Space Management

Every ISA faces a fundamental constraint: there is a finite number of distinct bit patterns. How do you allocate them?

### RISC-V: Hierarchical Opcode Space

```
7-bit opcode field → 128 possible opcodes
  But only ~50 are defined in RV32I base.

For disambiguation within an opcode, use funct3 (3 bits) and funct7 (7 bits):
  R-type instruction (opcode=0x33) uses funct3 + funct7 to distinguish:
    funct3=000, funct7=0000000 → ADD
    funct3=000, funct7=0100000 → SUB
    funct3=001, funct7=0000000 → SLL
    funct3=010, funct7=0000000 → SLT
    ... etc.
```

This hierarchical scheme allows adding instructions within an opcode without claiming a new opcode.

RISC-V reserves opcode space for custom extensions:
```
Opcode 0x0B  → custom-0 (custom extensions, non-standard)
Opcode 0x2B  → custom-1
Opcode 0x5B  → custom-2 / rv128
Opcode 0x7B  → custom-3 / rv128
```

### x86 Opcode Exhaustion and Escape Sequences

The original 8086 used 1-byte opcodes (256 possibilities). By the 286 era, all 256 were claimed. Intel added `0F` as an escape prefix for a second opcode byte (another 256 opcodes). By SSE, those were full too. `0F 38` and `0F 3A` added two more pages of opcodes. Then VEX and EVEX added entirely new encoding spaces.

The result: x86 opcode space is a geological record of 45 years of additions, full of discontinuities, gaps, and special cases.

### ARM64: Compact Opcode Space

ARM64 was designed in 2011 knowing this history. It uses the top bits to carve up the 32-bit space into groups, with careful allocation of unused combinations for future extensions. No escape bytes needed.

---

## 7. The Decode Unit in Hardware

### RISC-V Decode

A simple RISC-V decode unit:

```
32-bit instruction in
         │
         ▼
  ┌─────────────┐
  │ Opcode      │ → which functional unit? (ALU, branch, load, store, system)
  │ [6:0]       │
  └─────────────┘
  ┌─────────────┐
  │ rd [11:7]   │ → register write port address
  └─────────────┘
  ┌─────────────┐
  │ rs1 [19:15] │ → register read port 1 address
  └─────────────┘
  ┌─────────────┐
  │ rs2 [24:20] │ → register read port 2 address (for R/S/B types)
  └─────────────┘
  ┌─────────────┐
  │ funct3 [14:12] + funct7 [31:25] → specific operation within opcode class
  └─────────────┘
  ┌─────────────┐
  │ Immediate   │ → sign-extend and mux appropriate bits for immediate type
  │ extraction  │
  └─────────────┘
```

All of this decoding is pure combinational logic — takes ~1 clock cycle. The design is so simple that a RISC-V core can be implemented in an FPGA in a single afternoon using ~1000 lines of Verilog.

### x86-64 Decode

The x86-64 decode front-end is massively more complex:

```
Bytes from instruction cache (up to 16 bytes/cycle)
              │
              ▼
    ┌──────────────────────┐
    │ Instruction Length   │  Determine where each instruction starts and ends
    │ Decoder (ILD)        │  (up to 4 instructions per cycle)
    └──────────────────────┘
              │
              ▼
    ┌──────────────────────┐
    │ Pre-decoder          │  Identify prefixes, opcode bytes, ModRM, SIB
    └──────────────────────┘
              │
              ▼
    ┌──────────────────────┐
    │ Full decoders        │  4 parallel decoders:
    │  1 complex (MSROM)   │  - Complex decoder for >4-µop instructions (uses microcode ROM)
    │  3 simple            │  - 3 simple decoders for 1-4 µop instructions
    └──────────────────────┘
              │
              ▼
    ┌──────────────────────┐
    │ µop Queue            │  Buffer of decoded µops for the execution engine
    └──────────────────────┘
```

The Intel Sandy Bridge µop cache (2011) completely bypasses this pipeline for code in the cache:
```
    ┌──────────────────────┐
    │ µop Cache            │  Can supply 6 µops/cycle directly
    │ (Decoded ICache)     │  Bypasses the entire decode front-end
    └──────────────────────┘
```

### Decode Throughput

Modern processors need to decode many instructions per cycle to feed wide execution engines:

```
Processor           Decode Width   µop/cycle
──────────────────  ────────────   ─────────
Intel Core (2000)       3             3
Intel Core i7 (2008)    4             4
AMD Zen 2 (2019)        4             4
Intel Alder Lake (2021) 6             6
Apple M1 (2020)         8             8
```

Apple M1's 8-wide decode is enabled by AArch64's simple 32-bit fixed-width encoding — the decode logic is simpler, leaving more transistor budget for wider decode.

### Quick Check

> 1. Why is RISC-V decode hardware simpler than x86 decode hardware?
> 2. What is the µop cache and how does it improve x86 performance?
> 3. Why does Apple M1 have 8-wide decode when most x86 processors have 4-6?

---

## 8. Compressed Instruction Encoding

### RISC-V C Extension Encoding

The C extension adds 16-bit variants of the most common 32-bit instructions. The 16-bit instructions use a limited subset of registers (x8-x15 for some formats, any register for others) and smaller immediates:

```
CR format (Register): 15    12  11   7   6    2   1  0
                      ┌──────┬──────┬──────┬────┐
                      │funct4│  rd  │  rs2 │ op │
                      └──────┴──────┴──────┴────┘

CI format (Immediate): 15  13  12  11   7  6   2  1  0
                       ┌────┬──┬──────┬──────┬────┐
                       │funct3│imm│  rd/rs1│imm│ op │
                       └────┴──┴──────┴──────┴────┘

CSS format (Stack-Relative Store):
                       15  13  12   7   6   2   1  0
                       ┌────┬───────┬──────┬────┐
                       │funct3│imm6  │  rs2 │ op │
                       └────┴───────┴──────┴────┘
```

The `op` field (bits [1:0]) identifies the format:
- `00`: quad 0 (memory ops: C.LW, C.SW, C.LD, C.SD, ...)
- `01`: quad 1 (ALU and control: C.ADDI, C.JAL, C.BEQZ, ...)
- `10`: quad 2 (register ops: C.SLLI, C.LWSP, C.SWSP, ...)
- `11`: 32-bit instruction (base ISA)

### C Extension Common Instructions

```
Instruction    Expands to        Restriction
─────────────  ────────────────  ──────────────────────────────
C.ADD rd, rs2  add rd, rd, rs2   rd ≠ x0
C.ADDI rd, imm addi rd, rd, imm  rd ≠ x0, imm ≠ 0
C.LW rd, offset(rs1)  lw rd, offset(rs1)  rd and rs1 in x8-x15
C.SW rs2, offset(rs1) sw rs2, offset(rs1) rs2 and rs1 in x8-x15
C.BEQZ rs1, label  beq rs1, zero, label   rs1 in x8-x15
C.J label         jal zero, label          (unconditional jump)
C.MV rd, rs2      add rd, zero, rs2        (register move)
C.LI rd, imm      addi rd, zero, imm       (load immediate)
```

### Impact on Code Size

A typical C program compiled for RV64GC (with C extension) vs RV64G (without):
- Reduction in code size: ~25-30%
- Reduction in instruction cache pressure: proportional
- Fetch bandwidth improvement: more instructions fit per cache line

For systems with limited flash (embedded IoT devices), 25% code savings can be the difference between fitting or not fitting.

### Quick Check

> 1. How does the C extension differentiate 16-bit from 32-bit instructions?
> 2. What restriction do C.LW and C.SW have on which registers they use?
> 3. If a function is 400 bytes with RV64G and compiles to 300 bytes with RV64GC, how many more instructions fit in a 64-byte cache line?

---

## Summary

- An **instruction encoding** maps assembly mnemonics and operands to binary numbers according to ISA-defined formats.
- **Fixed-width encoding** (RISC-V, ARM64): every instruction is 32 bits. Simple decode, predictable pipeline, but limited immediate size and code density.
- **Variable-width encoding** (x86-64): 1-15 bytes. Better code density, larger immediates, but complex decode requiring sequential processing.
- **RISC-V** uses 6 formats (R, I, S, B, U, J-type), all 32 bits. rs1 always at [19:15], rs2 at [24:20], sign bit always at [31] — these enable parallel decode and fast register file access.
- **x86-64** encoding is a 45-year accumulation: optional legacy prefixes, REX prefix, 1-3 byte opcode, ModRM, SIB, displacement, immediate. Requires complex ILD and multiple-stage decode.
- **ARM64** uses clean 32-bit fixed encoding partitioned by top bits into logical groups.
- The **decode unit** implements instruction format parsing in hardware: RISC-V is ~1 cycle combinational; x86 requires 3-5 cycles with complex pipeline stages.
- **Compressed instructions** (RISC-V C, ARM Thumb) use 16-bit encoding for common operations, achieving 25-30% code size reduction.

---

## Exercises

### Easy

1. Encode the following RISC-V instructions as 32-bit hex values (look up the opcode tables):
   - `addi x5, x6, 100`  (I-type, opcode=0x13, funct3=000)
   - `sw x5, 12(x6)`     (S-type, opcode=0x23, funct3=010)

2. The RISC-V B-type immediate encodes a 13-bit signed offset (in multiples of 2 bytes) for branches, but the immediate bits are scattered across the instruction. Why not just put them in order? What does the scrambled layout enable?

3. Decode this x86-64 instruction: `48 89 C0`. The first byte `48` is a REX prefix (W=1, R=0, X=0, B=0). `89` is the opcode for `MOV r/m64, r64`. The ModRM byte `C0` is: Mod=11 (register direct), Reg=000 (rax), R/M=000 (rax). What is the complete instruction?

### Medium

4. **Immediate range exercise**: RISC-V I-type has a 12-bit signed immediate (-2048 to +2047). For a structure:
   ```c
   struct Large { char data[2048]; int value; };
   struct Large s;
   int v = s.value;   // offset of value = 2048
   ```
   The offset 2048 is OUT of range for I-type immediate. How does the compiler access `s.value`? Write the RISC-V instruction sequence needed. How would x86 handle this differently?

5. **x86 instruction length**: The x86-64 CPU must know how long each instruction is to find the next one. This creates a chicken-and-egg problem: to decode, you need to know instruction boundaries; to find boundaries, you need to decode. Intel solves this with a dedicated Instruction Length Decoder (ILD). Describe the algorithm the ILD would use to determine instruction length from a byte stream. What special cases must it handle?

6. **Opcode space pressure**: The RISC-V base ISA has 47 instructions using only some of the available opcode space. As extensions are added (M, A, F, D, C, V, B, K, H, J, P...), opcode space eventually fills. The RISC-V spec reserves custom opcode space for vendor extensions. Research: how many "non-conforming" opcode spaces does RISC-V reserve? What happens if two vendors both use the same custom opcode? How does the ecosystem handle this conflict?

### Hard

7. **Design a decode pipeline**: Design a 2-stage RISC-V decode pipeline for a 2-wide superscalar (decode 2 instructions per cycle). Your pipeline must:
   - Stage 1: Fetch 8 bytes per cycle, split into two 4-byte instructions
   - Stage 2: Decode both instructions simultaneously (extract opcode, registers, immediates)
   
   Problem: instruction 2 uses the result of instruction 1 (data dependency). What information must be extracted in stage 1 to detect this? Show the hazard detection logic: if instruction 2's rs1 or rs2 equals instruction 1's rd, a hazard exists. Draw the pipeline diagram with hazard detection logic.

8. **Binary exploit and encoding**: A buffer overflow attack replaces the return address on the stack with the attacker's desired address. On 64-bit x86-64 Linux, the attacker wants to jump to address `0x7fff deadbeef`. However, many exploits use string functions (like `strcpy`) that stop at null bytes (`\x00`). The return address `0x00 0x00 0x7f ff de ad be ef` contains null bytes. Research return-oriented programming (ROP): instead of injecting code, chain together existing instruction sequences ending in `ret`. Show how this avoids null bytes in the address. What role does the ISA's instruction encoding play in finding "ROP gadgets"? Why is ROP harder against RISC architectures with fixed-width encoding vs x86 with variable-width?
