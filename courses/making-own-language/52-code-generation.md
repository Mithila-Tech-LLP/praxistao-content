# Chapter 52: Code Generation — From IR to Machine Code

> "Good code is its own best documentation. As you're about to add a comment, ask yourself, 'How can I improve the code so that this comment isn't needed?'"
> — Steve McConnell
>
> (For compiler writers, replace "comment" with "manual instruction selection": good IR makes the translation to machine code obvious.)

---

This is the moment everything has been building toward. We have a lexer that tokenizes source code. A parser that builds an AST. A semantic analyzer that verifies names and scopes. A type checker that validates every operation. An IR builder that flattens structured code into simple three-address instructions. Now we need to translate that IR into real machine code — actual bytes that a CPU can execute.

Code generation is the back-end of the compiler: the part that knows about the target machine. While everything before this chapter was architecture-independent, code generation is inherently platform-specific. We must know where function arguments go (registers or stack?), how memory is laid out, which instructions perform which operations, and how many registers the CPU has. In this chapter we target **x86-64** — the 64-bit Intel architecture that runs virtually every laptop, desktop, and server in the world.

By the end of this chapter you will have a complete `codegen/x86_64.go` and `codegen/regalloc.go` that can take Astra's IR and emit valid x86-64 assembly, ready to be assembled with `as` and linked with `ld` or `gcc`. We will implement register allocation using the linear scan algorithm, emit function prologues and epilogues, handle the System V AMD64 calling convention, and generate a real assembly file with proper sections.

---

## What We're Building

A complete x86-64 code generator for Astra: an instruction selector, a linear-scan register allocator with spill handling, a function call code generator following the System V AMD64 ABI, and an assembly file emitter with `.text`, `.data`, and `.rodata` sections. Plus peephole optimizations to clean up obvious redundancies.

## Table of Contents

1. The Code Generator's Job
2. Target Architecture: x86-64
3. Instruction Selection
4. Register Allocation: The Core Problem
5. The Linear Scan Algorithm
6. Spilling: When Registers Run Out
7. Function Call Code Generation
8. Generating a Complete Assembly File
9. Calling External Functions
10. Peephole Optimization
11. Implementation: codegen/x86_64.go and codegen/regalloc.go
12. Astra Build Milestone

---

## 1. The Code Generator's Job

The code generator takes the compiler's IR — clean, architecture-independent, using unlimited virtual registers — and produces assembly code for a specific CPU. The IR is designed for the compiler's convenience. Assembly is designed for the CPU's convenience. The code generator bridges this gap.

The three major tasks of code generation are:

**Instruction Selection:** For each IR instruction, choose the corresponding machine instruction(s). `t1 = a + b` in IR becomes `mov rax, [a_slot]; add rax, [b_slot]; mov [t1_slot], rax` in assembly. There is usually more than one way to express an operation in machine code; instruction selection chooses the best one.

**Register Allocation:** IR uses unlimited temporaries (`t1`, `t2`, ..., `t1000`). Real CPUs have a small, fixed number of registers (x86-64 has 16 general-purpose 64-bit registers). The register allocator assigns each temporary to either a register or a stack slot (a memory location in the function's stack frame).

**Code Emission:** Produce the actual assembly text (or binary machine code) in the correct file format, with the right sections, symbol tables, and relocation information.

A secondary task — **peephole optimization** — examines the generated assembly a few instructions at a time and replaces obvious inefficiencies with better code. This is distinct from the IR-level optimizations in Chapter 70.

---

## 2. Target Architecture: x86-64

**x86-64** (also called AMD64, x86_64, or Intel 64) is the 64-bit extension of the classic x86 architecture. It is the instruction set architecture (ISA) of virtually all modern desktop processors (Intel Core, AMD Ryzen) and server processors (Intel Xeon, AMD EPYC).

### Why x86-64?

We choose x86-64 for Astra's initial backend because:
1. It is the most common server/desktop architecture — every developer's laptop can run the output directly
2. Its assembly language has extensive documentation and tooling
3. Understanding x86-64 assembly is a widely transferable skill
4. macOS (x86-64), Linux (x86-64), and Windows all use compatible calling conventions (with minor differences)

Alternative backends we can add later (Chapter 71 via LLVM):
- **ARM64 (AArch64):** Used by Apple Silicon (M1, M2, M3), AWS Graviton, Android phones
- **RISC-V:** Open ISA used in embedded systems and research
- **WebAssembly:** Stack-based VM ISA, runs in browsers and via WASI

### Registers

x86-64 has 16 general-purpose 64-bit registers:

```
General Purpose Registers (64-bit):
┌────────┬────────────────────────────────────────────────────────────┐
│ Name   │ Purpose                                                    │
├────────┼────────────────────────────────────────────────────────────┤
│ rax    │ Return value, accumulator, caller-saved                    │
│ rbx    │ Base register, callee-saved                                │
│ rcx    │ Counter, 4th argument, caller-saved                        │
│ rdx    │ Data, 3rd argument, caller-saved                           │
│ rsi    │ Source index, 2nd argument, caller-saved                   │
│ rdi    │ Destination index, 1st argument, caller-saved              │
│ rsp    │ Stack pointer (do not use for data!)                       │
│ rbp    │ Frame pointer, callee-saved                                │
│ r8     │ 5th argument, caller-saved                                 │
│ r9     │ 6th argument, caller-saved                                 │
│ r10    │ Caller-saved scratch                                       │
│ r11    │ Caller-saved scratch                                       │
│ r12    │ Callee-saved                                               │
│ r13    │ Callee-saved                                               │
│ r14    │ Callee-saved                                               │
│ r15    │ Callee-saved                                               │
└────────┴────────────────────────────────────────────────────────────┘

Each 64-bit register has 32-bit (eXX), 16-bit (XX), 8-bit (Xh/Xl) sub-registers:
  rax → eax (low 32) → ax (low 16) → ah (bits 8-15), al (bits 0-7)
  r8  → r8d (low 32) → r8w (low 16) → r8b (low 8)
```

### Key Instructions

```nasm
; Data movement
mov  rax, rbx          ; rax = rbx
mov  rax, [rbp-8]      ; rax = memory at address (rbp - 8)
mov  [rbp-8], rax      ; memory at (rbp - 8) = rax
push rax               ; decrement rsp, then store rax at [rsp]
pop  rax               ; load rax from [rsp], then increment rsp
lea  rax, [rip+.str0]  ; load address of .str0 into rax (position-independent)

; Arithmetic
add  rax, rbx          ; rax = rax + rbx
sub  rax, rbx          ; rax = rax - rbx
imul rax, rbx          ; rax = rax * rbx (signed multiply)
idiv rbx               ; rax = rax / rbx, rdx = rax % rbx (signed divide)
                       ; (must zero-extend rax into rdx first with cqo)
neg  rax               ; rax = -rax

; Comparison and jumps
cmp  rax, rbx          ; set flags based on rax - rbx (does not store result)
je   label             ; jump if equal (ZF=1)
jne  label             ; jump if not equal
jl   label             ; jump if less (signed)
jg   label             ; jump if greater (signed)
jle  label             ; jump if less or equal
jge  label             ; jump if greater or equal

; Logic
and  rax, rbx          ; rax = rax & rbx
or   rax, rbx          ; rax = rax | rbx
xor  rax, rbx          ; rax = rax ^ rbx (also used as: xor rax, rax → rax = 0)
shl  rax, cl           ; rax <<= cl (logical left shift)
shr  rax, cl           ; rax >>= cl (logical right shift)

; Function calls
call label             ; push return address, jump to label
ret                    ; pop return address, jump to it
```

---

## 3. Instruction Selection

Instruction selection maps each IR instruction to one or more assembly instructions. Astra uses a straightforward **template-based** approach: for each IR instruction pattern, we have a fixed template of assembly instructions.

```
IR Pattern          → Assembly Template
──────────────────────────────────────────────────────────────────────
t = a + b           → mov rax, <a>
                      add rax, <b>
                      mov <t>, rax

t = a - b           → mov rax, <a>
                      sub rax, <b>
                      mov <t>, rax

t = a * b           → mov rax, <a>
                      imul rax, <b>
                      mov <t>, rax

t = a / b           → mov rax, <a>
                      cqo            ; sign-extend rax into rdx:rax
                      idiv <b>
                      mov <t>, rax   ; quotient in rax

t = a % b           → mov rax, <a>
                      cqo
                      idiv <b>
                      mov <t>, rdx   ; remainder in rdx

t = a == b          → xor rcx, rcx  ; rcx = 0
                      cmp <a>, <b>
                      sete cl        ; cl = 1 if equal, 0 otherwise
                      mov <t>, rcx

t = a < b           → xor rcx, rcx
                      cmp <a>, <b>
                      setl cl
                      mov <t>, rcx

if t goto L1        → cmp <t>, 0
 else goto L2         jne L1
                      jmp L2

call f(a,b,...) → t → mov rdi, <a>  ; arg1
                      mov rsi, <b>  ; arg2
                      call f
                      mov <t>, rax  ; capture return value

return v            → mov rax, <v>
                      <function epilogue>
                      ret
```

In the templates, `<a>` means "the register or memory location holding variable `a`" — determined by the register allocator. After instruction selection and register allocation are combined, `<a>` becomes something like `rdi` (if allocated to a register) or `qword [rbp-16]` (if spilled to a stack slot).

### BURS and Optimal Instruction Selection

The template approach works but is not optimal. A more sophisticated approach is **BURS (Bottom-Up Rewriting Systems)**, which treats the IR instruction tree as a pattern-matching problem:

```
# Instead of:
t1 = a * 4
t2 = b + t1

# BURS might match the combined tree:
t2 = b + (a * 4)
# And emit a single instruction:
lea t2, [b + a*4]    ; x86 LEA handles base + index*scale!
```

The `lea` instruction in x86-64 can compute `base + index * scale + displacement` in a single instruction, which is much more efficient than separate multiply and add. BURS discovers these opportunities by matching patterns in the IR tree.

Astra uses the simpler template approach in this chapter. Chapter 70 adds tree-pattern matching.

---

## 4. Register Allocation: The Core Problem

IR has unlimited temporaries. x86-64 has 15 usable general-purpose registers (all except `rsp`). The **register allocation problem** is: assign each temporary either a register or a stack slot (in the function's activation record).

The goal is to maximize the number of temporaries that live in registers, because:
- Register access: ~1 clock cycle
- Cache L1 access (closest memory): ~4 clock cycles
- Cache L2 access: ~12 clock cycles
- Main memory access: ~200 clock cycles

Getting the allocation right can be the difference between a function running in 5 nanoseconds vs 50 nanoseconds.

### The Interference Graph

The classical approach to register allocation is **graph coloring**. Two temporaries **interfere** if they are both live at the same point in the program (i.e., they both hold values that will be needed in the future). Two interfering temporaries cannot share a register — they need separate locations.

We build an **interference graph**: one node per temporary, one edge between any two temporaries that interfere. Register allocation is then a graph coloring problem: assign colors (registers) to nodes such that no two adjacent nodes share a color, using at most k colors (where k is the number of available registers).

```
Example: temporaries t1, t2, t3 are live simultaneously at some point.
They all interfere with each other:

   t1 ─── t2
    │  ╲  │
    │    ╲│
   t3 ─── (t2 and t3 also interfere)

If k = 3 registers, this is colorable (one register per node).
If k = 2 registers, it is NOT colorable (K3 requires 3 colors).
→ One of t1, t2, t3 must be spilled to the stack.
```

Graph k-coloring is NP-complete in general. For register allocation, we use efficient heuristics. The most important for production compilers is the **Chaitin-Briggs** algorithm, which is what GCC and LLVM use.

For Astra's first version, we use a simpler algorithm that is still O(n) and produces good results for most code: **Linear Scan**.

---

## 5. The Linear Scan Algorithm

**Linear Scan** (Poletto and Sarkar, 1999) is a fast, practical register allocation algorithm used in the JVM's JIT compilers, Dalvik/ART (Android), and many other just-in-time and ahead-of-time compilers.

The key insight: instead of building an interference graph, compute **live intervals** — for each temporary, a range of instruction indices [start, end] during which the temporary is live (i.e., has been defined and will be used again). Then allocate registers to intervals in order of start point, greedily.

### Step 1: Compute Live Intervals

```
Instructions:              Intervals:
──────────────────────────  ────────────────────────────────
[0]  t1 = base             t1: [0, 4]  (defined at 0, last used at 4)
[1]  t2 = exp              t2: [1, 7]  (defined at 1, last used at 7)
[2]  t3 = 0                t3: [2, 3]  (defined at 2, last used at 3)
[3]  t4 = t2 > t3          t4: [3, 4]  (defined at 3, last used at 4)
[4]  if t4 goto body       
[5]  ...
[6]  t5 = result * t1      t5: [6, 7]  (defined at 6, last used at 7)
[7]  result = t5           
[8]  t6 = t2 - 1           t6: [8, 9]  
[9]  exp = t6              
```

A live interval starts at the instruction that **defines** the temporary and ends at the last instruction that **uses** it.

### Step 2: Allocate Registers in Order

Sort intervals by start point. Process each one:

```go
type Interval struct {
    Temp  string
    Start int
    End   int
}
```

For each interval in start-point order:
1. **Expire old intervals:** Remove from the "active" set any intervals whose end is before our start point. Return their registers to the free list.
2. **Allocate:** If a free register is available, assign it to this interval and add to active.
3. **Spill:** If no register is free, spill either this interval or the active interval with the farthest end point (whichever has the larger End — it's "in the way" for longer). The spilled interval gets a stack slot.

### Step 3: Linear Scan in Action

Available registers (simplified example with 3): `rax`, `rbx`, `rcx`

```
Process intervals in order of Start:

t1 [0,4]: active=[], free=[rax,rbx,rcx]
  → assign rax to t1. active=[t1(rax)], free=[rbx,rcx]

t2 [1,7]: active=[t1(rax)], free=[rbx,rcx]
  → assign rbx to t2. active=[t1(rax),t2(rbx)], free=[rcx]

t3 [2,3]: active=[t1(rax),t2(rbx)], free=[rcx]
  → assign rcx to t3. active=[t1(rax),t2(rbx),t3(rcx)], free=[]

t4 [3,4]: active=[t1(rax),t2(rbx),t3(rcx)], free=[]
  expire: t3.End=3 < t4.Start=3? No (equal). No expiry yet.
  NO FREE REGISTERS → must spill.
  Farthest-end active: t2 has End=7, farthest.
  Spill t2 to stack slot [rbp-8].
  Assign t2's register rbx to t4.
  active=[t1(rax),t4(rbx),t3(rcx)], free=[]
  (t2 → stack slot [rbp-8])

t5 [6,7]: 
  expire: t4.End=4 < 6 → expire t4, free rbx. t3.End=3 < 6 → expire t3, free rcx.
  active=[t1(rax)], free=[rbx,rcx]
  → assign rbx to t5. active=[t1(rax),t5(rbx)]
```

**Final allocation:**
```
t1 → rax      (register)
t2 → [rbp-8]  (spilled to stack)
t3 → rcx      (register, freed before t4 arrives)
t4 → rbx      (register, reused after t3 expires)
t5 → rbx      (register, reused after t4 expires)
```

Linear scan is O(n log n) (dominated by the sort of intervals by start point, plus O(n) active-set operations). It does not produce optimal register assignments like graph coloring can, but it is fast enough for JIT compilation and produces good results for ahead-of-time compilation.

---

## 6. Spilling: When Registers Run Out

When the register allocator decides a temporary must be **spilled**, it means that temporary's value must be saved to a stack slot (a memory location in the function's stack frame).

```
Stack Frame Layout (typical Astra function):
───────────────────────────────────────────────
rsp + 8  →  (return address, placed by call)
rbp      →  saved frame pointer (our rbp save)
rbp - 8  →  first local / spill slot (8 bytes)
rbp - 16 →  second local / spill slot
rbp - 24 →  third local / spill slot
...
rsp      →  (top of stack, aligned to 16)
───────────────────────────────────────────────
```

For each spilled temporary, we assign a stack slot offset from `rbp`. We also insert two kinds of code around uses:

**Spill (store) before a definition:** When the value is computed, store it to the stack slot.
```nasm
; t2 was spilled. After computing t2 = a + b:
add  rax, rsi
mov  [rbp-8], rax   ; spill store: save t2 to stack
```

**Reload (load) before a use:** When t2 is needed as an operand, load it from the stack slot.
```nasm
; Before using t2 in t5 = result * t2:
mov  rbx, [rbp-8]   ; reload: load spilled t2 into rbx
imul rax, rbx       ; now use rbx instead of t2
```

Spilling adds memory traffic, which hurts performance. Good register allocators minimize spills by choosing which temporary to spill wisely (usually the one with the fewest remaining uses, or the longest remaining live interval).

---

## 7. Function Call Code Generation

Function calls are the most complex part of code generation because they must follow the **calling convention** — the ABI (Application Binary Interface) that specifies exactly how arguments are passed, where return values go, and which registers must be preserved across calls.

Astra on macOS and Linux uses the **System V AMD64 ABI**:

```
Argument Passing (integers and pointers):
  1st argument → rdi
  2nd argument → rsi
  3rd argument → rdx
  4th argument → rcx
  5th argument → r8
  6th argument → r9
  7th+ arguments → pushed on stack, right to left

Return Value:
  Integer/pointer → rax
  64-bit float    → xmm0

Caller-Saved Registers (caller must save if needed across a call):
  rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11

Callee-Saved Registers (called function must restore before returning):
  rbx, rbp, r12, r13, r14, r15

Stack Alignment:
  rsp must be 16-byte aligned immediately before a call instruction
  (the call instruction pushes an 8-byte return address, making rsp
  temporarily misaligned, but the called function's prologue restores it)
```

### Function Prologue and Epilogue

Every Astra function begins with a **prologue** that sets up the stack frame and ends with an **epilogue** that tears it down:

```nasm
; PROLOGUE (at the start of every function)
_function_name:
    push  rbp           ; save caller's frame pointer
    mov   rbp, rsp      ; establish our frame pointer
    sub   rsp, <N>      ; allocate N bytes for locals/spills
                        ; N must keep rsp 16-byte aligned

; BODY
; ... function instructions ...

; EPILOGUE (before each return instruction)
    mov   rsp, rbp      ; restore stack pointer (frees locals)
    pop   rbp           ; restore caller's frame pointer
    ret                 ; return to caller
```

### Generating a Call Site

```nasm
; Calling: result = add(x, y)  where x→rdi, y→rsi, result→rax

; 1. Save caller-saved registers that are live across this call
;    (if we have live values in rdi, rsi, etc.)
push rdi    ; save if rdi holds something we need after the call
push rsi    ; save if rsi holds something we need after the call

; 2. Load arguments into the right registers
mov  rdi, [rbp-8]    ; x → 1st arg register (if x is in a stack slot)
mov  rsi, [rbp-16]   ; y → 2nd arg register (if y is in a stack slot)

; 3. Ensure 16-byte stack alignment before call
;    (rsp is 8-byte aligned here due to pushed return addr and rbp)
;    The call pushes 8 more bytes (return address), so rsp becomes
;    16-byte aligned. If not naturally aligned, add a sub rsp, 8.

; 4. Call the function
call _add

; 5. Retrieve return value
mov  [rbp-24], rax   ; store result from rax to result's location

; 6. Restore saved caller-saved registers
pop  rsi
pop  rdi
```

---

## 8. Generating a Complete Assembly File

A complete assembly file for Astra has four sections:

```
Assembly File Structure:
───────────────────────────────────────────────────────────────────
.rodata section:
  Read-only data. String literals live here.
  .str0: .asciz "Hello, World!\n"
  .str1: .asciz "Error: division by zero\n"

.data section:
  Initialized mutable global variables.
  .global_x: .quad 0     ; int x = 0 (64-bit integer, initialized to 0)

.bss section:
  Uninitialized global variables (just reserves space, no initial bytes).
  .global_y: .zero 8     ; int y (uninitialized, 8 bytes)

.text section:
  Executable code. All function definitions.
  .globl _main
  _main:
    push rbp
    mov  rbp, rsp
    ...
    pop  rbp
    ret
───────────────────────────────────────────────────────────────────
```

### Symbol Visibility

- **Global symbols** (`.globl _main`): visible to the linker, used for functions that must be callable from other object files or from the OS.
- **Local symbols** (no `.globl`): visible only within this object file.

On macOS, C symbols are prefixed with underscore: `main` in C becomes `_main` in assembly. On Linux, there is no underscore prefix. Astra's code generator handles this difference:

```go
func symbolName(name string) string {
    if runtime.GOOS == "darwin" {
        return "_" + name
    }
    return name
}
```

### String Literal Emission

Strings in Astra are UTF-8. Each unique string literal gets an entry in `.rodata`:

```nasm
; .rodata
.str0:
    .ascii "Hello, Astra!\0"   ; null-terminated string

; Using it in code:
lea  rdi, [rip + .str0]        ; RIP-relative addressing for PIC
call _println
```

The `[rip + .str0]` addressing mode is **position-independent code (PIC)**: the address is relative to the current instruction pointer, not an absolute address. This is required on macOS and modern Linux for position-independent executables.

---

## 9. Calling External Functions

Astra programs call into a small runtime library for operations that are hard to implement in pure assembly: printing, memory allocation, string operations, and so on.

```nasm
; Calling the Astra runtime function astra_print_int(n: int):
mov  rdi, [rbp-8]   ; load the integer value
call _astra_print_int

; Calling astra_alloc(size: int) -> *void:
mov  rdi, 64         ; allocate 64 bytes
call _astra_alloc
mov  [rbp-16], rax   ; store the returned pointer
```

These external symbols are declared at the top of each assembly file and resolved by the linker when it combines the Astra object file with the runtime library object file:

```nasm
; At the top of every generated .s file:
.extern _astra_print_int
.extern _astra_print_string
.extern _astra_print_float
.extern _astra_println
.extern _astra_alloc
.extern _astra_free
.extern _astra_string_concat
.extern _astra_string_length
```

---

## 10. Peephole Optimization

After full code generation, we scan the assembly output a few instructions at a time (a "peephole") and replace inefficient patterns with better ones. This is fast to implement and can produce significant speedups.

Common peephole optimizations:

**Remove redundant moves:**
```nasm
; Before:
mov rax, rax        ; move register to itself → useless
; After:
(removed)
```

**Collapse push/pop pairs:**
```nasm
; Before:
push rax
pop  rax            ; immediate push then pop of same register
; After:
(removed)
```

**Replace multiply by power of 2 with shift:**
```nasm
; Before:
imul rax, 4         ; multiply by 4
; After:
shl  rax, 2         ; shift left by 2 (same result, 1 cycle vs 3 cycles)
```

**Replace divide by power of 2 with shift:**
```nasm
; Before:
cqo
mov rbx, 8
idiv rbx            ; divide by 8 (expensive: ~20-90 cycles on modern CPUs)
; After:
sar  rax, 3         ; arithmetic right shift by 3 (1 cycle)
```

**Replace compare with zero after subtraction:**
```nasm
; Before:
sub  rax, rbx
cmp  rax, 0         ; sub already sets zero flag
; After:
sub  rax, rbx
; (the cmp is redundant — use jz/jnz directly after sub)
```

**Use xor to zero a register (faster than mov ... 0):**
```nasm
; Before:
mov  rax, 0         ; 5 bytes (REX + opcode + modrm + 4-byte imm)
; After:
xor  rax, rax       ; 3 bytes, and sets flags in 0 cycles (zeroing idiom)
```

---

## 11. Implementation: codegen/regalloc.go and codegen/x86_64.go

```go
// codegen/regalloc.go
package codegen

import "sort"

// The x86-64 general-purpose registers available for allocation.
// We exclude rsp (stack pointer) and rbp (frame pointer).
var allRegisters = []string{
    "rdi", "rsi", "rdx", "rcx", "r8", "r9",  // argument registers (caller-saved)
    "rax",                                      // return value / accumulator
    "r10", "r11",                               // caller-saved scratch
    "rbx", "r12", "r13", "r14", "r15",          // callee-saved
}

// Caller-saved: the function making a call must save these if needed across the call.
var callerSaved = map[string]bool{
    "rax": true, "rcx": true, "rdx": true, "rsi": true, "rdi": true,
    "r8": true, "r9": true, "r10": true, "r11": true,
}

// Callee-saved: the called function must restore these before returning.
var calleeSaved = map[string]bool{
    "rbx": true, "r12": true, "r13": true, "r14": true, "r15": true,
}

// Interval records the live range of a temporary.
type Interval struct {
    Temp  string
    Start int // instruction index where temp is first defined
    End   int // instruction index where temp is last used
}

// Allocation records the final register or stack location for a temp.
type Allocation struct {
    InReg    bool
    Register string // if InReg == true
    Offset   int    // if InReg == false: byte offset from rbp (negative)
}

// LinearScanAllocator assigns registers to temporaries using the linear scan algorithm.
type LinearScanAllocator struct {
    freeRegs  []string               // currently available registers
    active    []Interval             // intervals currently in registers, sorted by End
    location  map[string]*Allocation // result: temp → register or stack slot
    stackSize int                    // bytes used for spill slots (grows by 8 each spill)
}

func NewLinearScanAllocator() *LinearScanAllocator {
    regs := make([]string, len(allRegisters))
    copy(regs, allRegisters)
    return &LinearScanAllocator{
        freeRegs: regs,
        location: make(map[string]*Allocation),
    }
}

// Allocate performs register allocation for a list of live intervals.
// intervals must be sorted by Start.
func (a *LinearScanAllocator) Allocate(intervals []Interval) {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].Start < intervals[j].Start
    })

    for _, interval := range intervals {
        // Step 1: expire old intervals (those that ended before our start)
        a.expireOldIntervals(interval)

        if len(a.freeRegs) > 0 {
            // Step 2a: allocate a register
            reg := a.freeRegs[0]
            a.freeRegs = a.freeRegs[1:]
            a.location[interval.Temp] = &Allocation{InReg: true, Register: reg}
            // Insert into active list, sorted by End
            a.active = append(a.active, interval)
            sort.Slice(a.active, func(i, j int) bool {
                return a.active[i].End < a.active[j].End
            })
        } else {
            // Step 2b: spill — no free register available
            a.spillAtInterval(interval)
        }
    }
}

// expireOldIntervals removes intervals from active that ended before i.Start.
// Their registers are returned to freeRegs.
func (a *LinearScanAllocator) expireOldIntervals(i Interval) {
    remaining := a.active[:0]
    for _, ai := range a.active {
        if ai.End < i.Start {
            // This interval has ended; reclaim its register
            reg := a.location[ai.Temp].Register
            a.freeRegs = append(a.freeRegs, reg)
        } else {
            remaining = append(remaining, ai)
        }
    }
    a.active = remaining
}

// spillAtInterval spills either the current interval or the one with the
// farthest end point (whichever ends later).
func (a *LinearScanAllocator) spillAtInterval(i Interval) {
    if len(a.active) == 0 {
        a.spill(i)
        return
    }
    // Find the active interval with the farthest End
    spill := a.active[len(a.active)-1] // active is sorted by End ascending
    if spill.End > i.End {
        // Spill the active interval; give its register to i
        reg := a.location[spill.Temp].Register
        a.location[spill.Temp] = &Allocation{InReg: false, Offset: a.newStackSlot()}
        a.location[i.Temp] = &Allocation{InReg: true, Register: reg}
        // Remove spill from active, add i
        a.active = a.active[:len(a.active)-1]
        a.active = append(a.active, i)
        sort.Slice(a.active, func(x, y int) bool {
            return a.active[x].End < a.active[y].End
        })
    } else {
        // Spill i itself
        a.spill(i)
    }
}

func (a *LinearScanAllocator) spill(i Interval) {
    a.location[i.Temp] = &Allocation{InReg: false, Offset: a.newStackSlot()}
}

// newStackSlot allocates 8 bytes on the stack and returns the rbp-relative offset.
// Stack slots are laid out below rbp: first slot at rbp-8, second at rbp-16, etc.
func (a *LinearScanAllocator) newStackSlot() int {
    a.stackSize += 8
    return -a.stackSize // negative offset from rbp
}

// Location returns the allocation for a temporary.
func (a *LinearScanAllocator) Location(temp string) *Allocation {
    if loc, ok := a.location[temp]; ok {
        return loc
    }
    // Unknown temp: assign a new stack slot (error recovery)
    return &Allocation{InReg: false, Offset: a.newStackSlot()}
}

// StackFrameSize returns the total bytes needed for the stack frame.
// Rounded up to the nearest 16 to maintain 16-byte stack alignment.
func (a *LinearScanAllocator) StackFrameSize() int {
    size := a.stackSize
    if size%16 != 0 {
        size += 16 - (size % 16)
    }
    return size
}
```

```go
// codegen/x86_64.go
package codegen

import (
    "fmt"
    "runtime"
    "strings"
    "astra/ir"
)

// CodeGen translates Astra IR to x86-64 assembly text.
type CodeGen struct {
    out    strings.Builder
    alloc  *LinearScanAllocator
    fnName string // current function being generated
}

func NewCodeGen() *CodeGen {
    return &CodeGen{}
}

// Generate produces a complete assembly file from an IR program.
func (g *CodeGen) Generate(prog *ir.Program) string {
    g.emitLine("; Generated by Astra Compiler")
    g.emitLine("; Target: x86-64 (%s)", runtime.GOOS)
    g.emitLine("")

    // Declare external runtime symbols
    g.emitExterns()
    g.emitLine("")

    // .rodata section: string literals
    if len(prog.Strings) > 0 {
        g.emitLine(".section __TEXT,__cstring,cstring_literals" )
        for _, s := range prog.Strings {
            g.emitLine("%s:", s.Label)
            g.emitLine("    .asciz %q", s.Content)
        }
        g.emitLine("")
    }

    // .data section: initialized globals
    if len(prog.Globals) > 0 {
        g.emitLine(".data")
        for _, gv := range prog.Globals {
            g.emitLine("    .globl %s", g.sym(gv.Name))
            g.emitLine("%s:", g.sym(gv.Name))
            g.emitLine("    .quad %s", gv.Value)
        }
        g.emitLine("")
    }

    // .text section: function code
    g.emitLine(".text")
    for _, fn := range prog.Functions {
        g.generateFunction(fn)
    }

    return g.out.String()
}

func (g *CodeGen) emitExterns() {
    externs := []string{
        "astra_print_int", "astra_print_string", "astra_print_bool",
        "astra_println", "astra_alloc", "astra_free",
        "astra_string_concat", "astra_string_length",
    }
    for _, e := range externs {
        g.emitLine(".extern %s", g.sym(e))
    }
}

// ─────────────────────────────────────────────────────────────────────
// Function code generation
// ─────────────────────────────────────────────────────────────────────

func (g *CodeGen) generateFunction(fn *ir.Function) {
    g.fnName = fn.Name
    g.alloc = NewLinearScanAllocator()

    // Step 1: compute live intervals for all temporaries
    intervals := g.computeIntervals(fn)

    // Step 2: allocate registers
    g.alloc.Allocate(intervals)
    frameSize := g.alloc.StackFrameSize()

    // Step 3: emit prologue
    g.emitLine("")
    g.emitLine("    .globl %s", g.sym(fn.Name))
    g.emitLine("%s:", g.sym(fn.Name))
    g.emitLine("    pushq   %%rbp")
    g.emitLine("    movq    %%rsp, %%rbp")
    if frameSize > 0 {
        g.emitLine("    subq    $%d, %%rsp", frameSize)
    }

    // Save callee-saved registers that we use
    g.emitCalleeSaves()

    // Step 4: map function parameters from argument registers to their allocations
    argRegs := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
    for i, param := range fn.Params {
        if i < len(argRegs) {
            // Parameter arrived in argRegs[i]; copy to its allocated location
            loc := g.alloc.Location(param)
            if loc.InReg {
                if loc.Register != argRegs[i] {
                    g.emitLine("    movq    %%%s, %%%s", argRegs[i], loc.Register)
                }
            } else {
                g.emitLine("    movq    %%%s, %d(%%rbp)", argRegs[i], loc.Offset)
            }
        }
        // params 7+ would need to be loaded from the stack (not implemented here)
    }

    // Step 5: emit each basic block
    for _, block := range fn.Blocks {
        g.generateBlock(block)
    }

    // Step 6: emit epilogue (also emitted before each return instruction)
    g.emitLine(".%s_epilogue:", fn.Name)
    g.emitCalleeRestores()
    g.emitLine("    movq    %%rbp, %%rsp")
    g.emitLine("    popq    %%rbp")
    g.emitLine("    ret")
}

func (g *CodeGen) generateBlock(block *ir.BasicBlock) {
    g.emitLine("%s_%s:", g.fnName, block.Label)
    for _, instr := range block.Instructions {
        g.generateInstr(instr)
    }
}

// ─────────────────────────────────────────────────────────────────────
// Instruction code generation
// ─────────────────────────────────────────────────────────────────────

func (g *CodeGen) generateInstr(instr ir.Instruction) {
    switch i := instr.(type) {
    case *ir.LoadImm:
        dest := g.loc(i.Dest)
        switch i.Kind {
        case "int":
            g.emitLine("    movq    $%d, %s", i.IntVal, dest)
        case "bool":
            v := int64(0)
            if i.BoolVal { v = 1 }
            g.emitLine("    movq    $%d, %s", v, dest)
        case "string":
            label := fmt.Sprintf(".str_inline_%s", i.Dest)
            g.emitLine("    leaq    %s(%%rip), %%rax", label)
            g.emitLine("    movq    %%rax, %s", dest)
        }

    case *ir.Copy:
        src  := g.loc(i.Src)
        dest := g.loc(i.Dest)
        // If both are memory, use rax as intermediary
        if !g.isReg(i.Src) && !g.isReg(i.Dest) {
            g.emitLine("    movq    %s, %%rax", src)
            g.emitLine("    movq    %%rax, %s", dest)
        } else if src != dest {
            g.emitLine("    movq    %s, %s", src, dest)
        }

    case *ir.BinOp:
        g.generateBinOp(i)

    case *ir.UnOp:
        src  := g.loc(i.Src)
        dest := g.loc(i.Dest)
        g.emitLine("    movq    %s, %%rax", src)
        switch i.Op {
        case "-":
            g.emitLine("    negq    %%rax")
        case "!":
            g.emitLine("    xorq    $1, %%rax")
        }
        g.emitLine("    movq    %%rax, %s", dest)

    case *ir.Jump:
        g.emitLine("    jmp     %s_%s", g.fnName, i.Target)

    case *ir.CondJump:
        cond := g.loc(i.Cond)
        g.emitLine("    cmpq    $0, %s", cond)
        g.emitLine("    jne     %s_%s", g.fnName, i.TrueTarget)
        g.emitLine("    jmp     %s_%s", g.fnName, i.FalseTarget)

    case *ir.Call:
        g.generateCall(i)

    case *ir.Return:
        if i.Value != "" {
            src := g.loc(i.Value)
            if src != "%rax" {
                g.emitLine("    movq    %s, %%rax", src)
            }
        }
        g.emitLine("    jmp     .%s_epilogue", g.fnName)

    case *ir.IndexStore:
        arr   := g.loc(i.Arr)
        index := g.loc(i.Index)
        val   := g.loc(i.Val)
        g.emitLine("    movq    %s, %%rax",  arr)
        g.emitLine("    movq    %s, %%rcx",  index)
        g.emitLine("    movq    %s, %%rdx",  val)
        g.emitLine("    movq    %%rdx, (%%rax,%%rcx,8)") // arr[index*8] = val

    case *ir.IndexLoad:
        arr   := g.loc(i.Arr)
        index := g.loc(i.Index)
        dest  := g.loc(i.Dest)
        g.emitLine("    movq    %s, %%rax", arr)
        g.emitLine("    movq    %s, %%rcx", index)
        g.emitLine("    movq    (%%rax,%%rcx,8), %%rdx") // rdx = arr[index*8]
        g.emitLine("    movq    %%rdx, %s", dest)
    }
}

func (g *CodeGen) generateBinOp(i *ir.BinOp) {
    left  := g.loc(i.Left)
    right := g.loc(i.Right)
    dest  := g.loc(i.Dest)

    switch i.Op {
    case "+":
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    addq    %s, %%rax", right)
        g.emitLine("    movq    %%rax, %s", dest)
    case "-":
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    subq    %s, %%rax", right)
        g.emitLine("    movq    %%rax, %s", dest)
    case "*":
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    imulq   %s", right)
        g.emitLine("    movq    %%rax, %s", dest)
    case "/":
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    cqo")                      // sign-extend rax into rdx:rax
        g.emitLine("    idivq   %s", right)
        g.emitLine("    movq    %%rax, %s", dest)
    case "%":
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    cqo")
        g.emitLine("    idivq   %s", right)
        g.emitLine("    movq    %%rdx, %s", dest)  // remainder in rdx
    case "==", "!=", "<", ">", "<=", ">=":
        g.emitLine("    xorq    %%rcx, %%rcx")     // rcx = 0
        g.emitLine("    movq    %s, %%rax", left)
        g.emitLine("    cmpq    %s, %%rax", right)
        switch i.Op {
        case "==":
            g.emitLine("    sete    %%cl")
        case "!=":
            g.emitLine("    setne   %%cl")
        case "<":
            g.emitLine("    setl    %%cl")
        case ">":
            g.emitLine("    setg    %%cl")
        case "<=":
            g.emitLine("    setle   %%cl")
        case ">=":
            g.emitLine("    setge   %%cl")
        }
        g.emitLine("    movq    %%rcx, %s", dest)
    }
}

func (g *CodeGen) generateCall(i *ir.Call) {
    argRegs := []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

    // Save caller-saved registers that are currently in use
    // (In a real implementation, we track which are live across the call)
    // For simplicity, we save/restore rdi..r11 except those used for args.

    // Move arguments into argument registers
    for j, arg := range i.Args {
        if j >= len(argRegs) {
            // Push extra args onto stack (right to left)
            // (simplified: only handles up to 6 args here)
            break
        }
        src := g.loc(arg)
        argReg := "%" + argRegs[j]
        if src != argReg {
            g.emitLine("    movq    %s, %s", src, argReg)
        }
    }

    // Align stack to 16 bytes before call (if frame size is odd multiple of 8)
    // In production we track this precisely; here we use a simple heuristic.
    g.emitLine("    call    %s", g.sym(i.FuncName))

    if i.Dest != "" {
        dest := g.loc(i.Dest)
        if dest != "%rax" {
            g.emitLine("    movq    %%rax, %s", dest)
        }
    }
}

// ─────────────────────────────────────────────────────────────────────
// Helper methods
// ─────────────────────────────────────────────────────────────────────

// loc returns the assembly operand for a temporary (register or memory).
func (g *CodeGen) loc(temp string) string {
    alloc := g.alloc.Location(temp)
    if alloc.InReg {
        return "%" + alloc.Register
    }
    return fmt.Sprintf("%d(%%rbp)", alloc.Offset)
}

// isReg returns true if the temp is allocated to a register.
func (g *CodeGen) isReg(temp string) bool {
    return g.alloc.Location(temp).InReg
}

// sym applies the platform-specific symbol name prefix.
func (g *CodeGen) sym(name string) string {
    if runtime.GOOS == "darwin" {
        return "_" + name
    }
    return name
}

func (g *CodeGen) emitCalleeSaves() {
    // Save callee-saved registers that our allocation uses
    for reg := range calleeSaved {
        // In a real implementation, check if we actually allocate anything to this reg.
        // For simplicity, save all callee-saved registers.
        g.emitLine("    pushq   %%%s", reg)
    }
}

func (g *CodeGen) emitCalleeRestores() {
    // Restore in reverse order
    saved := []string{"r15", "r14", "r13", "r12", "rbx"}
    for _, reg := range saved {
        g.emitLine("    popq    %%%s", reg)
    }
}

func (g *CodeGen) emitLine(format string, args ...interface{}) {
    if len(args) > 0 {
        g.out.WriteString(fmt.Sprintf(format, args...) + "\n")
    } else {
        g.out.WriteString(format + "\n")
    }
}

// computeIntervals computes live intervals for all temporaries in a function.
// This is a simplified liveness analysis: we scan instructions linearly.
func (g *CodeGen) computeIntervals(fn *ir.Function) []Interval {
    // Map from temp name → [first_def, last_use]
    defs := make(map[string]int)
    uses := make(map[string]int)

    instrIdx := 0
    for _, block := range fn.Blocks {
        for _, instr := range block.Instructions {
            // Record definitions and uses based on instruction type
            switch i := instr.(type) {
            case *ir.BinOp:
                if _, ok := defs[i.Dest]; !ok { defs[i.Dest] = instrIdx }
                uses[i.Left] = instrIdx; uses[i.Right] = instrIdx
            case *ir.UnOp:
                if _, ok := defs[i.Dest]; !ok { defs[i.Dest] = instrIdx }
                uses[i.Src] = instrIdx
            case *ir.Copy:
                if _, ok := defs[i.Dest]; !ok { defs[i.Dest] = instrIdx }
                uses[i.Src] = instrIdx
            case *ir.LoadImm:
                if _, ok := defs[i.Dest]; !ok { defs[i.Dest] = instrIdx }
            case *ir.Call:
                for _, arg := range i.Args { uses[arg] = instrIdx }
                if i.Dest != "" {
                    if _, ok := defs[i.Dest]; !ok { defs[i.Dest] = instrIdx }
                }
            case *ir.Return:
                if i.Value != "" { uses[i.Value] = instrIdx }
            case *ir.CondJump:
                uses[i.Cond] = instrIdx
            }
            instrIdx++
        }
    }

    // Also add parameters as defined at instruction 0
    for _, param := range fn.Params {
        if _, ok := defs[param]; !ok {
            defs[param] = 0
        }
        if _, ok := uses[param]; !ok {
            uses[param] = 0
        }
    }

    var intervals []Interval
    for temp, start := range defs {
        end := start
        if u, ok := uses[temp]; ok && u > end {
            end = u
        }
        intervals = append(intervals, Interval{Temp: temp, Start: start, End: end})
    }
    return intervals
}
```

---

## 🔨 Astra Build Milestone

### Complete Assembly for Factorial

```astra
fn factorial(n: int) -> int {
    if n <= 1 { return 1 }
    return n * factorial(n - 1)
}
```

**IR (from the lowerer):**
```
fn factorial(n):
.entry_1:
    t1 = 1
    t2 = n <= t1
    if t2 goto .then_2 else goto .else_3
.then_2:
    t3 = 1
    return t3
.else_3:
    t4 = 1
    t5 = n - t4
    t6 = call factorial(t5)
    t7 = n * t6
    return t7
```

**Register Allocation (with 5 temps, 3 available registers: rdi, rsi, rax):**

| Temp | Live Range | Assignment |
|---|---|---|
| n (param) | [0, 8] | rdi (parameter register, then spill) |
| t1 | [1, 2] | rax |
| t2 | [2, 3] | rsi |
| t3 | [4, 5] | rax |
| t4 | [6, 7] | rax |
| t5 | [7, 8] | rsi |
| t6 | [8, 9] | rdi |
| t7 | [9, 10] | rax |

**Generated x86-64 Assembly (macOS format):**

```nasm
    ; Generated by Astra Compiler
    ; Target: x86-64 darwin

    .extern _astra_print_int

    .text
    .globl _factorial
_factorial:
    pushq   %rbp                    ; save caller's frame pointer
    movq    %rsp, %rbp              ; establish our frame pointer
    subq    $48, %rsp               ; allocate 48 bytes for locals/spills
                                    ; (must be 16-byte aligned: ✓)
    pushq   %rbx                    ; save callee-saved registers
    pushq   %r12
    pushq   %r13
    pushq   %r14
    pushq   %r15

    ; Parameter: n arrives in %rdi, save it to stack for preservation across call
    movq    %rdi, -8(%rbp)          ; n → [rbp-8]

_factorial_.entry_1:
    movq    $1, %rax                ; t1 = 1       (LoadImm: literal 1)
    movq    -8(%rbp), %rcx         ; load n into rcx for comparison
    cmpq    %rax, %rcx             ; compare n with t1 (value 1)
    xorq    %rsi, %rsi             ; t2 = 0 initially
    setle   %sil                    ; t2 = (n <= 1) ? 1 : 0
    movq    %rsi, -16(%rbp)        ; save t2 to stack

    cmpq    $0, -16(%rbp)          ; test condition t2
    jne     _factorial_.then_2      ; if t2 != 0 (i.e., n <= 1), jump to then
    jmp     _factorial_.else_3      ; otherwise fall to else

_factorial_.then_2:
    movq    $1, %rax               ; t3 = 1       (the return value 1)
    jmp     .factorial_epilogue    ; return 1

_factorial_.else_3:
    movq    $1, %rax               ; t4 = 1       (for n - 1)
    movq    -8(%rbp), %rcx        ; reload n
    subq    %rax, %rcx            ; t5 = n - 1
    movq    %rcx, -24(%rbp)       ; save t5

    ; Prepare for recursive call: factorial(n - 1)
    movq    -24(%rbp), %rdi       ; arg1 = t5 = n - 1
    call    _factorial             ; call factorial(n-1)
    movq    %rax, -32(%rbp)       ; t6 = return value, save to stack

    ; Compute n * factorial(n-1)
    movq    -8(%rbp), %rax        ; reload n
    imulq   -32(%rbp)             ; t7 = n * t6 (result in rax)
                                   ; (idiv/imul use rax implicitly)
    ; rax now holds n * factorial(n-1), which is our return value
    jmp     .factorial_epilogue

.factorial_epilogue:
    ; rax already holds the return value
    popq    %r15                    ; restore callee-saved registers (reverse order)
    popq    %r14
    popq    %r13
    popq    %r12
    popq    %rbx
    movq    %rbp, %rsp             ; restore stack pointer
    popq    %rbp                   ; restore frame pointer
    ret                            ; return (pops saved return address, jumps to it)
```

### Assembly Instruction Annotations

| Assembly Line | What It Does | Why |
|---|---|---|
| `pushq %rbp` | Save caller's frame pointer | Required by ABI: we use rbp as frame pointer |
| `movq %rsp, %rbp` | Set our frame pointer | rbp now marks the base of our frame |
| `subq $48, %rsp` | Allocate 48 bytes on stack | Space for local variables and spill slots |
| `pushq %rbx` etc. | Save callee-saved registers | ABI requires: restore before returning |
| `movq %rdi, -8(%rbp)` | Save parameter n to stack | rdi will be overwritten by recursive call |
| `movq $1, %rax` | Load constant 1 into rax | For the comparison n <= 1 |
| `cmpq %rax, %rcx` | Compare n and 1 | Sets CPU flags for setle |
| `setle %sil` | Set sil = 1 if n <= 1 | Converts condition to 0/1 boolean |
| `jne _factorial_.then_2` | Jump if t2 != 0 | Branches to base case |
| `subq %rax, %rcx` | rcx = n - 1 | Computes recursive argument |
| `movq %rcx, %rdi` | Load arg into rdi | Prepares for recursive call |
| `call _factorial` | Recursive call | Return value will be in rax |
| `movq %rax, -32(%rbp)` | Save return value to stack | rax may be overwritten; save it |
| `movq -8(%rbp), %rax` | Reload n | n was saved before the call |
| `imulq -32(%rbp)` | rax = rax * [rbp-32] | n * factorial(n-1) |
| `popq %r15` etc. | Restore callee-saved regs | ABI requirement, must match pushq |
| `movq %rbp, %rsp` | Restore stack pointer | Frees all local variables |
| `popq %rbp` | Restore caller's rbp | Restores caller's frame |
| `ret` | Return to caller | Pops saved return address, jumps |

---

## Exercises

1. **Peephole Pass Implementation:** Implement a `PeepholeOptimize(lines []string) []string` function that applies at least three of the peephole optimizations described in section 10. Test it on the factorial assembly output and measure how many instructions are removed.

2. **Stack Frame Calculator:** Given a function with 8 temporaries and 3 available registers, manually run the linear scan algorithm to determine which temporaries get registers and which get spilled. Then calculate the minimum stack frame size (in bytes, 16-byte aligned) needed for the spill slots and saved callee-saved registers.

3. **6-Argument Function:** Write the complete assembly code generation (prologue, argument marshalling, body, epilogue) for this function:
   ```astra
   fn clamp(val: int, lo: int, hi: int) -> int {
       if val < lo { return lo }
       if val > hi { return hi }
       return val
   }
   ```
   Pay attention to which argument registers are used and whether any need to be saved to the stack before the comparison.

4. **Assembly File Driver:** Write a complete Go program `cmd/compile/main.go` that reads an Astra source file from the command line, runs the full compilation pipeline (lex → parse → sema → type check → lower to IR → codegen), and writes the resulting assembly to a `.s` file. Test it by feeding the output to the system assembler (`as`) and linker.

5. **Calling Convention Compliance Test:** Write a small C program that calls an Astra-compiled function. This tests that the generated assembly follows the System V AMD64 ABI correctly. The C program declares the Astra function as an `extern` C function and calls it with several arguments. Trace through the calling convention rules manually and verify the assembly is correct.

6. **Tail Call Optimization:** The factorial function as written is not tail-recursive. Rewrite it to be tail-recursive (using an accumulator parameter), and then explain what assembly-level optimization **tail call optimization (TCO)** would produce: instead of a `call` and `ret`, TCO replaces the recursive call with a `jmp` back to the function's entry point, eliminating stack growth. Implement this optimization for self-recursive tail calls.

7. **Float Support:** Extend the code generator to handle `float` type (64-bit IEEE 754 doubles). On x86-64, floating-point operations use `xmm0`–`xmm15` registers and SSE2 instructions (`movsd`, `addsd`, `mulsd`, `comisd`). What changes need to be made to the register allocator, instruction selector, and calling convention handling?

8. **Live Interval Refinement:** The simplified `computeIntervals` function in this chapter uses a linear scan through instructions without regard for control flow (loops). A live range in a loop might be longer than computed because the value survives across loop iterations. Explain how you would fix this by processing the CFG using a reverse post-order traversal and iterating to a fixed point.

---

## Summary

| Concept | Key Idea |
|---|---|
| Code generation | Translates architecture-independent IR into target-specific assembly |
| Instruction selection | Maps each IR instruction to one or more machine instructions |
| Register allocation | Assigns unlimited IR temps to a finite set of CPU registers |
| Linear scan | O(n log n) allocation: sort intervals by start, greedily assign registers |
| Live interval | Range [start, end] in instruction stream where a temp holds a live value |
| Spilling | When registers are exhausted, store temp's value in a stack memory slot |
| Stack frame | Memory region for a function's locals and spills; accessed via rbp offsets |
| System V AMD64 ABI | Calling convention: args in rdi/rsi/rdx/rcx/r8/r9, return in rax |
| Caller-saved | rax, rcx, rdx, rdi, rsi, r8-r11: save before calls if needed |
| Callee-saved | rbx, r12-r15: function must restore before returning |
| Function prologue | push rbp; mov rbp, rsp; sub rsp, N — sets up stack frame |
| Function epilogue | mov rsp, rbp; pop rbp; ret — tears down stack frame |
| Peephole optimization | Local 2-3 instruction improvements; removes obvious redundancies |
| RIP-relative addressing | Position-independent addressing for string/data literals |
| `.text` / `.data` / `.rodata` | Assembly file sections for code, mutable data, and read-only data |
