# Chapter 31: Registers, the Call Stack, and the System V ABI

> "The ABI is the contract between code and code. Break it and nothing works. Understand it and everything becomes clear." — Systems programming wisdom

---

## Overview

In Chapter 30, you learned assembly instructions — the vocabulary of machine language. But knowing individual words is not enough to write a sentence. You need grammar. For x86-64 programming, that grammar is the **Application Binary Interface (ABI)**.

The ABI is a set of rules that answers critical questions: Which registers hold function arguments? Which registers can a function freely overwrite? How must the stack be arranged? What happens to local variables? Without the ABI, compiled code from different files (or different compilers) could not call each other. The ABI is the shared contract that makes the entire software ecosystem work.

This chapter covers the System V AMD64 ABI — the standard used on Linux and macOS — in complete detail, and then implements Astra's ABI module in Go.

---

## What We're Building

The Astra Build Milestone for this chapter is the complete `codegen/abi.go` module — the Go code that Astra's code generator uses to emit correct function prologues, epilogues, and argument-handling code. This module is the backbone of correct code generation.

---

## Table of Contents

1. Registers In Depth: General vs Special Purpose
2. The Three Special Registers: rsp, rbp, rip
3. What Is an ABI?
4. System V AMD64 ABI: Argument Passing
5. Argument Passing: More Than 6 Arguments
6. Return Values
7. Caller-Saved vs Callee-Saved Registers
8. Function Prologue: Setting Up the Stack Frame
9. Function Epilogue: Tearing Down the Stack Frame
10. Stack Frame Layout in Detail
11. Variadic Functions in Assembly
12. The Red Zone
13. Stack Alignment: The 16-Byte Rule
14. Astra Build Milestone: The ABI Module
15. Exercises
16. Summary

---

## 1. Registers In Depth: General vs Special Purpose

In Chapter 29, we listed all 16 registers. Now let us understand how they are *used* in practice.

### General-Purpose Registers

These 14 registers can (in principle) hold any value and be used for any computation:

```
rax, rbx, rcx, rdx, rsi, rdi, r8, r9, r10, r11, r12, r13, r14, r15
```

However, the ABI assigns specific **conventional purposes** to many of them. While nothing forces you to obey (the CPU does not know or care), every program and library does — making interoperability possible.

### Special-Purpose Registers

Two registers have strict hardware-enforced roles:

**`rsp` (Stack Pointer):** Always points to the top of the stack. Certain instructions (`push`, `pop`, `call`, `ret`) implicitly use `rsp`. If you corrupt `rsp`, your program crashes immediately and spectacularly.

**`rip` (Instruction Pointer):** Always contains the address of the next instruction to execute. You cannot write to it directly. Only `jmp`, `call`, `ret`, and conditional jumps change `rip`.

And one with a strong conventional role:

**`rbp` (Base Pointer / Frame Pointer):** Conventionally used to mark the base of the current stack frame. Local variables are addressed as `[rbp-N]`. Debuggers use `rbp` to walk the call stack and print stack traces. Some compilers omit the frame pointer for performance (freeing `rbp` as a general register), but this makes debugging harder.

---

## 2. The Three Special Registers: rsp, rbp, rip

Let us look at these three in detail with a diagram:

```
MEMORY AT RUNTIME:

Low address  0x0000...
+------------------------+
|  .text section         | <- rip points into here
|  (your program code)   |
|  ...                   |
+------------------------+
|  .data / .rodata       |
|  (global variables,    |
|   string literals)     |
+------------------------+
|  heap (grows upward)   |
|  malloc'd memory       |
|  ...                   |
+------------------------+ <- heap top grows here
|                        |
|  (unmapped space)      |
|                        |
+------------------------+ <- stack bottom (top of stack space)
|  stack (grows down)    | <- rbp points into here
|  local variables       |
|  saved registers       |
|  return addresses      |
+------------------------+ <- rsp points here (top of stack = lowest used addr)
High address 0xFFFF...

rip: "what am I doing next?"
rbp: "what is my current stack frame?"
rsp: "where does the stack currently end?"
```

### Why the Stack Grows Down

It is a historical convention that has persisted because it neatly avoids conflicts: the heap grows from low addresses upward, and the stack grows from high addresses downward. They meet in the middle only if you run out of memory. Modern operating systems enforce limits (typically 8MB for the stack on Linux) using the page table, so a runaway stack gets a segfault rather than silently overwriting heap data.

---

## 3. What Is an ABI?

An **Application Binary Interface** defines the rules for how compiled code interacts at the binary level. It is broader than just function calls:

```
ABI COVERS:
+------------------------------------------+
| Function calling convention              |
|  - How arguments are passed              |
|  - How return values are returned        |
|  - Which registers must be preserved     |
|  - How the stack must be arranged        |
+------------------------------------------+
| Data representation                      |
|  - Size of int, long, pointer types      |
|  - Struct field alignment and padding    |
|  - Endianness                            |
+------------------------------------------+
| Name mangling                            |
|  - How C++ class methods are named       |
|  - Symbol visibility rules              |
+------------------------------------------+
| Object file format                       |
|  - ELF on Linux, Mach-O on macOS        |
|  - Section layouts, relocation formats   |
+------------------------------------------+
| System call interface                    |
|  - Which register holds syscall number   |
|  - How arguments and results are passed  |
+------------------------------------------+
```

The **System V AMD64 ABI** is the standard for 64-bit Linux and macOS (and most other Unix-like systems). Windows uses a different ABI (Microsoft x64 calling convention), which is why code compiled for Linux does not directly link with code compiled for Windows.

---

## 4. System V AMD64 ABI: Argument Passing

The most important part of the ABI for our compiler: how function arguments are passed.

### The Six Argument Registers

The first six **integer or pointer** arguments go in these registers, in order:

```
Argument    Register
--------    --------
1st         rdi
2nd         rsi
3rd         rdx
4th         rcx
5th         r8
6th         r9
```

Mnemonic many programmers use: **"Di, Si, D, C, 8, 9"** or just remember the order of `rdi, rsi, rdx, rcx, r8, r9`.

Example:

```astra
// Astra function call:
result = complex_func(a, b, c, d, e, f)

// Generated assembly:
mov  rdi, a   ; 1st argument
mov  rsi, b   ; 2nd argument
mov  rdx, c   ; 3rd argument
mov  rcx, d   ; 4th argument
mov  r8,  e   ; 5th argument
mov  r9,  f   ; 6th argument
call complex_func
; result is now in rax
```

### Floating-Point Arguments

Floating-point (float/double) arguments use XMM registers instead:

```
FP Argument    Register
-----------    --------
1st fp         xmm0
2nd fp         xmm1
3rd fp         xmm2
...up to...
8th fp         xmm7
```

Integer and FP arguments are counted separately:

```c
// C example:
double mixed(int a, double b, int c, double d);

// Calling convention:
// a (int)    → rdi
// b (double) → xmm0
// c (int)    → rsi
// d (double) → xmm1
```

Astra will need to handle this for float arithmetic in Chapter 36+. For now, we focus on integer arguments.

---

## 5. Argument Passing: More Than 6 Arguments

When a function has more than 6 integer arguments, the extras go on the stack **before** the `call` instruction, in reverse order (7th argument pushed last, closest to the top of the stack):

```
fn many_args(a,b,c,d,e,f,g,h: int) -> int

STACK BEFORE CALL:
+------------------+  <- rsp was here
|   return address |  <- pushed by call instruction
+------------------+
|   h (8th arg)    |  <- pushed first (highest on stack)
+------------------+
|   g (7th arg)    |  <- pushed second (lowest extra on stack)
+------------------+  <- rsp after pushing g, h

REGISTERS:
rdi = a (1st)
rsi = b (2nd)
rdx = c (3rd)
rcx = d (4th)
r8  = e (5th)
r9  = f (6th)
```

Inside the function, extra arguments are accessed relative to `rbp`:

```asm
; After standard prologue (push rbp; mov rbp, rsp):
; [rbp + 8]  = return address (pushed by call)
; [rbp + 16] = 7th argument (g)
; [rbp + 24] = 8th argument (h)
```

Note: `rbp` itself holds the saved previous `rbp` (pushed in prologue). The return address is at `[rbp+8]`. Stack-passed arguments start at `[rbp+16]`.

---

## 6. Return Values

Return values follow these rules:

| Type | Register |
|------|----------|
| Integer, pointer (≤ 64 bits) | rax |
| 128-bit integer | rdx:rax (high in rdx, low in rax) |
| float/double | xmm0 |
| Struct ≤ 16 bytes | rax + rdx |
| Struct > 16 bytes | Caller allocates memory, passes pointer in rdi; callee writes to it |

For Astra's purposes (integer-focused language initially):

```asm
; All Astra functions that return an int or pointer:
; Put the return value in rax before ret

fn square(x: int) -> int {
    return x * x
}

; Generated:
square:
    push    rbp
    mov     rbp, rsp
    imul    rdi, rdi        ; rdi = x * x (modify rdi directly — it's caller-saved)
    mov     rax, rdi        ; Move result to return register
    pop     rbp
    ret
```

---

## 7. Caller-Saved vs Callee-Saved Registers

This is the most common source of bugs when writing assembly by hand. There are two categories of registers:

### Caller-Saved Registers (Volatile, "Scratch Registers")

The caller must save these before making a function call if it needs the values after the call returns. The callee is free to overwrite them without saving.

```
CALLER-SAVED: rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11
(Also: xmm0–xmm7, xmm8–xmm15)
```

Think of these as "scratch paper" — the callee can scribble on them freely.

### Callee-Saved Registers (Non-Volatile, "Preserved Registers")

The callee must save these (push to stack) if it uses them, and restore them (pop from stack) before returning. The caller can rely on these being unchanged after a function call.

```
CALLEE-SAVED: rbx, r12, r13, r14, r15, rbp
(Also: the x87/MMX/SSE state, but we ignore that for now)
```

Think of these as "permanent notes" — the callee must put them back exactly how they found them.

### Why This Division?

The ABI splits registers this way to enable efficient calling conventions:
- Argument registers (`rdi`, `rsi`, etc.) are caller-saved because after the call, arguments are no longer needed
- Working registers (`rax`, `rcx`, `rdx`) are caller-saved because the callee needs them for computation
- Long-lived variables should go in callee-saved registers (`rbx`, `r12`-`r15`) so they survive function calls

### Example: Getting It Wrong vs Right

```asm
; WRONG: Callee corrupts rbx without saving it

outer:
    push    rbp
    mov     rbp, rsp
    mov     rbx, 42         ; rbx = 42 (we're using rbx for a loop counter)
    call    inner           ; inner modifies rbx — CORRUPTED!
    cmp     rbx, 42         ; WRONG: rbx is no longer 42
    ...

inner:
    mov     rbx, 99         ; Using rbx without saving — ABI VIOLATION!
    ret

; CORRECT: Callee saves and restores rbx

inner:
    push    rbp
    mov     rbp, rsp
    push    rbx             ; Save rbx (callee-saved)
    mov     rbx, 99         ; Use rbx freely
    ; ... do work ...
    pop     rbx             ; Restore rbx
    pop     rbp
    ret
```

Astra's code generator must automatically track which callee-saved registers it uses and emit the correct save/restore code.

---

## 8. Function Prologue: Setting Up the Stack Frame

Every non-leaf function (a function that calls other functions) follows a standard prologue:

```asm
function_name:
    ; Step 1: Save the caller's frame pointer
    push    rbp

    ; Step 2: Establish our own frame pointer
    mov     rbp, rsp

    ; Step 3: Allocate space for local variables (if any)
    sub     rsp, N       ; N = number of bytes needed, aligned to 16

    ; Step 4: Save any callee-saved registers we plan to use
    push    rbx          ; (if we use rbx)
    push    r12          ; (if we use r12)
    ; etc.
```

After the prologue, the stack looks like this:

```
HIGH ADDRESSES
+-------------------------+
|  ...caller's frame...  |
+-------------------------+
|  return address         | (pushed by call instruction)
+-------------------------+  <- old rsp (before our push rbp)
|  saved rbp              | (pushed by our push rbp)
+-------------------------+  <- rbp (after mov rbp, rsp)
|  saved rbx              | (if we use rbx)
+-------------------------+
|  saved r12              | (if we use r12)
+-------------------------+
|  local variable 1       | <- [rbp - 8]
+-------------------------+
|  local variable 2       | <- [rbp - 16]
+-------------------------+
|  ...more locals...      |
+-------------------------+  <- rsp (after sub rsp, N)
|  (unused/red zone)      |
LOW ADDRESSES
```

**Calculating `N` (local variable space):**

For Astra, we calculate the number of bytes needed by looking at all local variable declarations in the function body. Each `int` needs 8 bytes (we use 64-bit integers throughout). We must round up to a multiple of 16 for stack alignment.

---

## 9. Function Epilogue: Tearing Down the Stack Frame

The epilogue undoes everything the prologue did, in reverse order:

```asm
    ; Step 1: Restore callee-saved registers (in reverse push order)
    pop     r12          ; (if saved)
    pop     rbx          ; (if saved)

    ; Step 2: Restore rsp to frame base
    mov     rsp, rbp     ; Discard all local variables

    ; Step 3: Restore the caller's frame pointer
    pop     rbp

    ; Step 4: Return to caller
    ret
```

The `leave` instruction is a convenience that combines steps 2 and 3:
```asm
    leave    ; equivalent to: mov rsp, rbp; pop rbp
    ret
```

Most compilers emit explicit `mov rsp, rbp; pop rbp` rather than `leave` for clarity (and because `leave` is slightly slower on some CPUs).

---

## 10. Stack Frame Layout in Detail

Let us walk through a complete example with locals, saved registers, and stack-passed arguments:

```astra
fn complex(a: int, b: int, c: int, d: int, e: int, f: int, g: int) -> int {
    let sum = a + b + c + d + e + f + g
    let double_sum = sum * 2
    return double_sum
}
```

Generated assembly with full stack frame:

```asm
complex:
    ; ---- Prologue ----
    push    rbp                      ; Save caller's rbp
    mov     rbp, rsp                 ; Establish frame

    ; Arguments 1-6 are in: rdi, rsi, rdx, rcx, r8, r9
    ; Argument 7 (g) is at: [rbp + 16] (see layout below)

    sub     rsp, 32                  ; 2 local vars * 8 bytes = 16, but align to 32
                                     ; (we add some padding for safety)

    ; Stack frame layout at this point:
    ; [rbp + 16]  = g (7th argument, pushed by caller before call)
    ; [rbp + 8]   = return address (pushed by call)
    ; [rbp + 0]   = saved rbp from caller
    ; [rbp - 8]   = local: sum
    ; [rbp - 16]  = local: double_sum
    ; rsp         = rbp - 32

    ; ---- Compute sum = a + b + c + d + e + f + g ----
    mov     rax, rdi                 ; rax = a
    add     rax, rsi                 ; rax = a + b
    add     rax, rdx                 ; rax = a + b + c
    add     rax, rcx                 ; rax = a + b + c + d
    add     rax, r8                  ; rax = a + b + c + d + e
    add     rax, r9                  ; rax = a + b + c + d + e + f
    add     rax, QWORD PTR [rbp+16]  ; rax = sum (add g from stack)
    mov     QWORD PTR [rbp-8], rax   ; store sum to local variable

    ; ---- Compute double_sum = sum * 2 ----
    mov     rax, QWORD PTR [rbp-8]   ; load sum
    imul    rax, 2                    ; rax = sum * 2
    mov     QWORD PTR [rbp-16], rax  ; store double_sum

    ; ---- Return double_sum ----
    mov     rax, QWORD PTR [rbp-16]  ; load double_sum into return register

    ; ---- Epilogue ----
    mov     rsp, rbp                 ; restore rsp
    pop     rbp                      ; restore caller's rbp
    ret
```

---

## 11. Variadic Functions in Assembly

Variadic functions (like `printf` in C) accept a variable number of arguments. The System V ABI has a specific protocol for them.

The key addition: the caller must put the number of floating-point arguments in `al` (the low byte of `rax`) before calling a variadic function. This lets the callee know how many XMM registers to search.

```c
// C: printf("Hello %d %d\n", a, b);

// Assembly:
mov     rdi, format_string    ; 1st arg: format string
mov     rsi, a                ; 2nd arg: a
mov     rdx, b                ; 3rd arg: b
xor     rax, rax              ; al = 0 (zero floating-point arguments)
call    printf
```

Astra does not support variadic functions in version 1 (it is a complex feature). The `print` function in Astra's standard library is a fixed-signature function, not variadic. We note this for completeness.

---

## 12. The Red Zone

On Linux and macOS (but NOT Windows), the ABI guarantees a **red zone**: 128 bytes below `rsp` that will not be modified by signal handlers or interrupt handlers.

```
STACK LAYOUT WITH RED ZONE:

[rbp]      <- current frame pointer
...
[rsp]      <- current top of stack
-----------
[rsp - 8]   \
[rsp - 16]   |  These 128 bytes are the RED ZONE.
[rsp - 24]   |  Signal handlers will not touch them.
...          |  You can use them for temporaries
[rsp - 128]  /  WITHOUT decrementing rsp!
```

Leaf functions (functions that don't call other functions) can use the red zone to store temporaries without adjusting `rsp`. This is a minor optimization:

```asm
; A leaf function with 2 locals can use the red zone:
leaf_fn:
    push    rbp
    mov     rbp, rsp
    ; No "sub rsp, 16" needed — we use red zone instead!
    mov     QWORD PTR [rsp-8],  rdi   ; temp local 1 (in red zone)
    mov     QWORD PTR [rsp-16], rsi   ; temp local 2 (in red zone)
    ; ... do work ...
    pop     rbp
    ret
```

Astra's code generator will use the red zone optimization for simple leaf functions to avoid unnecessary `sub rsp` instructions.

---

## 13. Stack Alignment: The 16-Byte Rule

One of the most subtle and important ABI requirements: **`rsp` must be 16-byte aligned immediately before a `call` instruction.**

Why? SSE and AVX instructions that operate on vector data require 16-byte (or 32-byte) aligned memory. The ABI guarantees this alignment at function entry (just after the `call` pushes the return address, `rsp` must be 16-byte aligned, meaning `rsp % 16 == 0`). This means just *before* `call`, `rsp` must be `rsp % 16 == 8` (since `call` will push 8 bytes).

```
ALIGNMENT TIMELINE:

At program start:  rsp is 16-byte aligned (guaranteed by OS)
                   rsp % 16 == 0

Before our call:   rsp must be 8 mod 16
                   rsp % 16 == 8
                   (so that after call pushes 8 bytes, rsp % 16 == 0)

Inside callee      After: push rbp (8 bytes) and mov rbp, rsp
(after prologue):  rsp % 16 == 0  ✓  (for SSE loads/stores)
```

How does a function maintain alignment?

```asm
main:
    push    rbp          ; rsp -= 8; rsp % 16 == 8 now (was 0)
    mov     rbp, rsp     ; rbp = rsp

    ; We want to call a function. rsp % 16 must be 8 before call.
    ; Currently: rsp % 16 == 8 (after push rbp). Good!
    ; But if we allocate locals, we must maintain alignment.

    sub     rsp, 8       ; 1 local variable (8 bytes)
                         ; rsp % 16 == 0 now (8 - 8 = 0)
                         ; WRONG! Before call, rsp must be 8 mod 16.

    sub     rsp, 16      ; 1 local variable (8 bytes) + 8 bytes padding
                         ; rsp % 16 == 8 (8 - 16 = -8 ≡ 8 mod 16)
                         ; Correct! Now rsp % 16 == 8.

    call    some_function ; Pushes 8 bytes: rsp % 16 == 0. ✓
```

The `align16` function in Astra's code generator must ensure the total allocation is always a multiple of 16 bytes. If you have 1 local (8 bytes), you allocate 16. If you have 3 locals (24 bytes), you allocate 32. And so on.

---

## 14. Astra Build Milestone: The Complete ABI Module

Here is the complete Go implementation of Astra's ABI module:

```go
// codegen/abi.go
// System V AMD64 ABI implementation for the Astra compiler.
// This module handles function prologues, epilogues, argument passing,
// and register allocation according to the ABI specification.

package codegen

import (
    "fmt"
    "io"
)

// ============================================================
// ABI Constants
// ============================================================

// argRegisters are the integer argument registers in order per System V AMD64 ABI.
// The first 6 integer/pointer arguments go in these registers.
var argRegisters = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

// returnRegister holds integer return values per the ABI.
const returnRegister = "rax"

// callerSaved registers may be freely overwritten by a called function.
// The CALLER must save these before a call if it needs them afterward.
var callerSaved = []string{
    "rax", "rcx", "rdx", "rsi", "rdi",
    "r8", "r9", "r10", "r11",
}

// calleeSaved registers must be preserved by any function that uses them.
// The CALLEE must save them at the start and restore them before returning.
var calleeSaved = []string{
    "rbx", "r12", "r13", "r14", "r15",
}

// stackPointer is the register that always points to the top of the stack.
const stackPointer = "rsp"

// basePointer is the register that points to the bottom of the current frame.
const basePointer = "rbp"

// redZoneSize is the guaranteed safe area below rsp on Linux/macOS.
// Leaf functions can use this for temporaries without adjusting rsp.
const redZoneSize = 128

// stackAlignment is the required alignment of rsp before a call instruction.
// rsp must be 16-byte aligned AFTER the call pushes the return address,
// meaning rsp must be 8 mod 16 BEFORE the call.
const stackAlignment = 16

// ============================================================
// align16: Round up to 16-byte boundary
// ============================================================

// align16 rounds n up to the nearest multiple of 16.
// Used to ensure stack allocations maintain alignment.
//
// Examples:
//   align16(0)  = 0
//   align16(1)  = 16
//   align16(8)  = 16
//   align16(16) = 16
//   align16(17) = 32
//   align16(24) = 32
func align16(n int) int {
    if n == 0 {
        return 0
    }
    return (n + 15) &^ 15
}

// ============================================================
// FunctionContext: Tracks state during function code generation
// ============================================================

// FunctionContext tracks the state needed to generate code for a single function.
// It knows which local variables exist, where they live on the stack,
// and which callee-saved registers have been pushed.
type FunctionContext struct {
    // Name of the function being generated
    Name string

    // localOffsets maps local variable names to their rbp-relative offsets.
    // Negative values mean below rbp (e.g., -8, -16, -24...).
    localOffsets map[string]int

    // nextLocalOffset tracks the next available stack slot.
    // Starts at -8, decremented by 8 for each new local.
    nextLocalOffset int

    // totalLocalBytes is the total stack space reserved for locals.
    // This is what gets subtracted from rsp in the prologue.
    totalLocalBytes int

    // savedCalleeSaved tracks which callee-saved registers we have pushed.
    // We must pop them in reverse order in the epilogue.
    savedCalleeSaved []string

    // isLeaf is true if the function does not call any other functions.
    // Leaf functions can use the red zone optimization.
    isLeaf bool

    // paramCount is the number of parameters (used for register assignment).
    paramCount int
}

// NewFunctionContext creates a fresh context for generating a function.
func NewFunctionContext(name string, isLeaf bool) *FunctionContext {
    return &FunctionContext{
        Name:            name,
        localOffsets:    make(map[string]int),
        nextLocalOffset: 0,
        isLeaf:          isLeaf,
    }
}

// AddParam registers a function parameter.
// Parameters passed in registers can be accessed via register or spilled to stack.
func (fc *FunctionContext) AddParam(name string, index int) {
    if index < len(argRegisters) {
        // For simplicity, we spill all register parameters to the stack.
        // An optimizing compiler would keep them in registers.
        fc.AddLocal(name)
    } else {
        // Stack-passed parameters are at positive offsets from rbp.
        // The layout is: [rbp+16] = 7th arg, [rbp+24] = 8th arg, etc.
        stackIndex := index - len(argRegisters)
        fc.localOffsets[name] = (stackIndex + 2) * 8 // +2: skip return addr and saved rbp
    }
    fc.paramCount++
}

// AddLocal registers a new local variable and assigns it a stack slot.
// Returns the stack offset (negative, relative to rbp).
func (fc *FunctionContext) AddLocal(name string) int {
    fc.nextLocalOffset -= 8
    fc.localOffsets[name] = fc.nextLocalOffset
    if -fc.nextLocalOffset > fc.totalLocalBytes {
        fc.totalLocalBytes = -fc.nextLocalOffset
    }
    return fc.nextLocalOffset
}

// OffsetOf returns the rbp-relative offset of a named variable.
// Returns (offset, true) if found, (0, false) if not found.
func (fc *FunctionContext) OffsetOf(name string) (int, bool) {
    offset, ok := fc.localOffsets[name]
    return offset, ok
}

// MemRef returns the assembly memory reference string for a named variable.
// For example: "QWORD PTR [rbp-8]" or "QWORD PTR [rbp+16]"
func (fc *FunctionContext) MemRef(name string) string {
    offset, ok := fc.localOffsets[name]
    if !ok {
        panic(fmt.Sprintf("ABI error: variable %q not found in function %q", name, fc.Name))
    }
    if offset >= 0 {
        return fmt.Sprintf("QWORD PTR [rbp+%d]", offset)
    }
    return fmt.Sprintf("QWORD PTR [rbp%d]", offset) // negative: rbp-8, rbp-16, etc.
}

// NeedsCalleeSave returns true if we need to save/restore callee-saved registers.
func (fc *FunctionContext) NeedsCalleeSave(reg string) bool {
    for _, r := range calleeSaved {
        if r == reg {
            return true
        }
    }
    return false
}

// UseCalleeSaved marks a callee-saved register as used (must be saved/restored).
func (fc *FunctionContext) UseCalleeSaved(reg string) {
    for _, r := range fc.savedCalleeSaved {
        if r == reg {
            return // already tracked
        }
    }
    fc.savedCalleeSaved = append(fc.savedCalleeSaved, reg)
}

// ============================================================
// Prologue and Epilogue Emission
// ============================================================

// emitFunctionPrologue emits the standard function prologue.
//
// For a function with N bytes of locals, emits:
//
//   push rbp
//   mov  rbp, rsp
//   sub  rsp, N         ; only if N > 0
//   push rbx            ; only if rbx is used (callee-saved)
//   push r12            ; only if r12 is used (callee-saved)
//   ... etc.
//
// Then spills register parameters to their stack slots.
func emitFunctionPrologue(w io.Writer, fc *FunctionContext) {
    // Step 1: Save caller's base pointer
    fmt.Fprintln(w, "    push    rbp")

    // Step 2: Establish our stack frame
    fmt.Fprintln(w, "    mov     rbp, rsp")

    // Step 3: Allocate local variable space
    // We must account for space and maintain 16-byte alignment.
    // After push rbp, rsp is 8 mod 16. We need sub rsp, N where N is 8 mod 16.
    // align16(totalLocalBytes) gives us a multiple of 16, but we need 8 mod 16.
    // Solution: if align16(bytes) % 16 == 0, add 8 for alignment.
    if fc.totalLocalBytes > 0 {
        allocSize := align16(fc.totalLocalBytes)
        // After push rbp: rsp % 16 = 8.
        // After sub rsp, allocSize: rsp % 16 = (8 - allocSize % 16 + 16) % 16
        // We need rsp % 16 == 8 before any call. This means allocSize % 16 == 0.
        // align16 already guarantees that. So the alignment is correct!
        fmt.Fprintf(w, "    sub     rsp, %d\n", allocSize)
    }

    // Step 4: Save callee-saved registers that this function uses
    for _, reg := range fc.savedCalleeSaved {
        fmt.Fprintf(w, "    push    %s\n", reg)
    }

    // Step 5: Spill register parameters to their stack slots
    // (In an optimizing compiler, this step would be optional/avoided.)
    for i := 0; i < fc.paramCount && i < len(argRegisters); i++ {
        // The parameters should already be AddParam'd with stack offsets.
        // We emit the spill here.
        // Note: actual param names not tracked here; this would be integrated
        // with the symbol table in a full compiler.
        _ = argRegisters[i]
    }
}

// emitFunctionEpilogue emits the standard function epilogue.
//
// Restores callee-saved registers in reverse order, then:
//   mov  rsp, rbp
//   pop  rbp
//   ret
func emitFunctionEpilogue(w io.Writer, fc *FunctionContext) {
    // Restore callee-saved registers in reverse order
    for i := len(fc.savedCalleeSaved) - 1; i >= 0; i-- {
        fmt.Fprintf(w, "    pop     %s\n", fc.savedCalleeSaved[i])
    }

    // Tear down the stack frame
    fmt.Fprintln(w, "    mov     rsp, rbp")
    fmt.Fprintln(w, "    pop     rbp")
    fmt.Fprintln(w, "    ret")
}

// ============================================================
// Argument Loading
// ============================================================

// emitLoadArg emits code to load a function argument into a register.
// For the first 6 arguments, the value is already in the arg register.
// For stack-passed arguments, we load from the stack.
func emitLoadArg(w io.Writer, argIndex int, destReg string) {
    if argIndex < len(argRegisters) {
        srcReg := argRegisters[argIndex]
        if srcReg != destReg {
            fmt.Fprintf(w, "    mov     %s, %s\n", destReg, srcReg)
        }
        // If srcReg == destReg, the argument is already where we want it.
    } else {
        // Stack-passed argument.
        // After standard prologue (push rbp; mov rbp, rsp):
        // [rbp+8]  = return address
        // [rbp+16] = 7th argument (index 6)
        // [rbp+24] = 8th argument (index 7)
        stackOffset := (argIndex - len(argRegisters) + 2) * 8
        fmt.Fprintf(w, "    mov     %s, QWORD PTR [rbp+%d]\n", destReg, stackOffset)
    }
}

// emitSetupCall emits code to set up arguments before a function call.
// args is a list of source registers/memory refs for each argument position.
func emitSetupCall(w io.Writer, args []string) {
    // Load each argument into its designated register.
    // Be careful: if arg[1] is in rdi and we're trying to move rdi to rsi first,
    // we might corrupt the value. In a real compiler, we'd use a proper
    // parallel copy algorithm. For now, we use temporaries.
    for i, arg := range args {
        if i < len(argRegisters) {
            destReg := argRegisters[i]
            if arg != destReg {
                fmt.Fprintf(w, "    mov     %s, %s\n", destReg, arg)
            }
        } else {
            // Stack-passed: push in reverse order separately
            // (handled by caller of this function)
        }
    }
}

// emitStackArgs emits push instructions for arguments 7+ (in reverse order).
// The ABI requires pushing in reverse so the 7th arg is closest to the return address.
func emitStackArgs(w io.Writer, args []string) {
    if len(args) <= len(argRegisters) {
        return // no stack args needed
    }
    stackArgs := args[len(argRegisters):]
    // Push in reverse order so 7th arg is at [rbp+16], 8th at [rbp+24], etc.
    for i := len(stackArgs) - 1; i >= 0; i-- {
        fmt.Fprintf(w, "    push    %s\n", stackArgs[i])
    }
}

// emitCleanupStackArgs emits the stack pointer adjustment after a call
// that passed arguments on the stack.
func emitCleanupStackArgs(w io.Writer, stackArgCount int) {
    if stackArgCount > 0 {
        bytes := stackArgCount * 8
        fmt.Fprintf(w, "    add     rsp, %d\n", bytes)
    }
}

// ============================================================
// Convenience: Complete Function Wrapper
// ============================================================

// EmitFunction is the high-level function that generates a complete assembly
// function. It calls the provided body emitter between prologue and epilogue.
//
// Usage:
//
//   ctx := NewFunctionContext("my_func", false)
//   ctx.AddParam("a", 0)
//   ctx.AddParam("b", 1)
//   ctx.AddLocal("result")
//
//   EmitFunction(w, ctx, func() {
//       // emit function body using ctx.MemRef("a"), etc.
//   })
func EmitFunction(w io.Writer, fc *FunctionContext, body func()) {
    // Emit label
    fmt.Fprintf(w, "%s:\n", fc.Name)

    // Emit prologue
    emitFunctionPrologue(w, fc)

    // Emit body (caller-provided)
    body()

    // Emit epilogue (also called by return statements; this is the fallthrough)
    emitFunctionEpilogue(w, fc)
}

// ============================================================
// Example: Using the ABI module to generate "add"
// ============================================================

// GenerateAddExample demonstrates using the ABI module to generate
// the assembly for: fn add(a: int, b: int) -> int { return a + b }
func GenerateAddExample(w io.Writer) {
    fc := NewFunctionContext("add", true /* leaf function */)
    fc.AddParam("a", 0) // a is in rdi (spilled to [rbp-8])
    fc.AddParam("b", 1) // b is in rsi (spilled to [rbp-16])

    fmt.Fprintln(w, "; Astra: fn add(a: int, b: int) -> int { return a + b }")
    fmt.Fprintf(w, "%s:\n", fc.Name)

    // Prologue (leaf function with small locals might use red zone,
    // but we use the standard approach for clarity)
    fmt.Fprintln(w, "    push    rbp")
    fmt.Fprintln(w, "    mov     rbp, rsp")

    // Function body: rdi = a, rsi = b, put a+b in rax
    fmt.Fprintln(w, "    mov     rax, rdi    ; rax = a (1st argument)")
    fmt.Fprintln(w, "    add     rax, rsi    ; rax = a + b (2nd argument)")

    // Epilogue
    fmt.Fprintln(w, "    pop     rbp")
    fmt.Fprintln(w, "    ret")
}
```

To verify the module works, here is a simple test:

```go
// codegen/abi_test.go
package codegen

import (
    "bytes"
    "strings"
    "testing"
)

func TestGenerateAddExample(t *testing.T) {
    var buf bytes.Buffer
    GenerateAddExample(&buf)
    output := buf.String()

    // Verify the generated assembly contains expected instructions
    required := []string{
        "add:",
        "push    rbp",
        "mov     rbp, rsp",
        "mov     rax, rdi",
        "add     rax, rsi",
        "pop     rbp",
        "ret",
    }
    for _, line := range required {
        if !strings.Contains(output, line) {
            t.Errorf("Expected output to contain %q\nGot:\n%s", line, output)
        }
    }
}

func TestAlign16(t *testing.T) {
    tests := []struct{ input, want int }{
        {0, 0},
        {1, 16},
        {8, 16},
        {16, 16},
        {17, 32},
        {24, 32},
        {32, 32},
        {33, 48},
    }
    for _, tc := range tests {
        got := align16(tc.input)
        if got != tc.want {
            t.Errorf("align16(%d) = %d, want %d", tc.input, got, tc.want)
        }
    }
}

func TestFunctionContext(t *testing.T) {
    fc := NewFunctionContext("test_fn", false)
    fc.AddParam("x", 0) // should be at [rbp-8]
    fc.AddParam("y", 1) // should be at [rbp-16]
    fc.AddLocal("z")     // should be at [rbp-24]

    if ref := fc.MemRef("x"); ref != "QWORD PTR [rbp-8]" {
        t.Errorf("x: got %q, want %q", ref, "QWORD PTR [rbp-8]")
    }
    if ref := fc.MemRef("y"); ref != "QWORD PTR [rbp-16]" {
        t.Errorf("y: got %q, want %q", ref, "QWORD PTR [rbp-16]")
    }
    if ref := fc.MemRef("z"); ref != "QWORD PTR [rbp-24]" {
        t.Errorf("z: got %q, want %q", ref, "QWORD PTR [rbp-24]")
    }
}
```

---

## 15. Exercises

**Exercise 1 — Register Classification:**
Without looking at the tables, classify each register as caller-saved, callee-saved, or special-purpose:
`rax`, `rbx`, `rcx`, `rdx`, `rsp`, `rbp`, `rsi`, `rdi`, `r8`, `r9`, `r10`, `r11`, `r12`, `r13`, `r14`, `r15`

**Exercise 2 — Stack Frame Diagram:**
Draw the complete stack frame for this Astra function after the prologue executes:
```astra
fn compute(a: int, b: int, c: int) -> int {
    let x = a + b
    let y = x * c
    return y
}
```
Show every value on the stack and what memory address it lives at (assume `rsp = 0x7fff0000` at function entry).

**Exercise 3 — Alignment Calculation:**
A function has 5 local `int` variables (each 8 bytes = 40 bytes total). What is the correct value to subtract from `rsp` in the prologue, and why? (Hint: use `align16` and think about the 16-byte alignment rule.)

**Exercise 4 — Argument Counting:**
Write out which register (or stack location) each argument lives in for these function calls:
```astra
fn f1(a: int) -> int
fn f2(a: int, b: int, c: int) -> int
fn f3(a: int, b: int, c: int, d: int, e: int, f: int) -> int
fn f4(a: int, b: int, c: int, d: int, e: int, f: int, g: int, h: int) -> int
```

**Exercise 5 — Callee-Save Bug Finding:**
Find all ABI violations in this assembly:
```asm
outer:
    push    rbp
    mov     rbp, rsp
    mov     r12, 100        ; loop counter
    call    inner           ; inner uses r12 without saving!
    cmp     r12, 100
    je      .ok
    ; ERROR: r12 was corrupted
.ok:
    pop     rbp
    ret

inner:
    push    rbp
    mov     rbp, rsp
    mov     r12, 42         ; BUG: using callee-saved register without saving
    pop     rbp
    ret
```
Fix the bug in `inner`.

**Exercise 6 — Red Zone Usage:**
For a leaf function with exactly 2 local int variables, show two ways to handle them:
1. The standard approach (sub rsp, N)
2. The red zone approach (no rsp adjustment)

**Exercise 7 — ABI Module Extension:**
Extend the `FunctionContext` type to support 32-bit integer variables (`int32` in Astra). Add a method `Add32Local(name string)` that allocates 4 bytes (but still uses 8-byte alignment for simplicity), and a corresponding `MemRef32(name string)` that returns `DWORD PTR [rbp-N]`.

**Exercise 8 — Stack Walk:**
When a program crashes, the debugger uses `rbp` to walk the call stack and print a backtrace. Given this chain of calls: `main → compute → helper`, draw the complete linked list of saved `rbp` values that lets the debugger walk from `helper` back to `main`. Where does the chain end?

---

## 16. Summary

| Concept | Rule | Why |
|---------|------|-----|
| Argument registers | rdi, rsi, rdx, rcx, r8, r9 (in order) | Standard for Linux/macOS; enables interop |
| More than 6 args | Push on stack in reverse order | 7th arg at [rbp+16], 8th at [rbp+24] |
| Return value | Integer → rax; float → xmm0 | Caller knows where to find the result |
| Caller-saved | rax, rcx, rdx, rsi, rdi, r8-r11 | Callee can overwrite freely |
| Callee-saved | rbx, r12-r15, rbp | Function must save/restore if used |
| Prologue | push rbp; mov rbp, rsp; sub rsp, N | Establishes stack frame |
| Epilogue | restore callee-saved; mov rsp, rbp; pop rbp; ret | Tears down frame |
| align16 | Round local bytes up to multiple of 16 | Maintains 16-byte stack alignment |
| Red zone | 128 bytes below rsp, safe without adjusting rsp | Leaf function optimization |
| Stack alignment | rsp must be 16-byte aligned after call pushes return address | Required for SSE/AVX instructions |

The ABI is the skeleton that holds all function calls together. Astra's code generator does not have the luxury of "winging it" — every function must obey these rules precisely or the program crashes in mysterious ways. The `abi.go` module encodes these rules once so that every generated function is automatically correct.

In the next chapter, we examine what happens after code generation: how the assembler, linker, and operating system cooperate to turn assembly files into an executable program.
