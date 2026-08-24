# Chapter 21: Calling Conventions and the ABI

Programs are never monolithic. They consist of thousands of functions, compiled separately, possibly in different languages, possibly by different compilers. The **Application Binary Interface (ABI)** is the contract that makes all these pieces work together. It specifies: how are arguments passed? Where is the return value? What happens to registers when a function is called? How is the stack managed? This chapter covers calling conventions in depth, from the reasoning behind them to real-world consequences.

## Table of Contents

1. [What is an ABI?](#1-what-is-an-abi)
2. [Why Calling Conventions Matter](#2-why-calling-conventions-matter)
3. [RISC-V Calling Convention (RISC-V psABI)](#3-risc-v-calling-convention-risc-v-psabi)
4. [System V AMD64 ABI (x86-64 Linux/macOS)](#4-system-v-amd64-abi-x86-64-linuxmacos)
5. [Windows x64 ABI](#5-windows-x64-abi)
6. [ARM64 ABI (AAPCS64)](#6-arm64-abi-aapcs64)
7. [Stack Frame Layout](#7-stack-frame-layout)
8. [Passing Large Arguments](#8-passing-large-arguments)
9. [Varargs and Printf-style Functions](#9-varargs-and-printf-style-functions)
10. [Tail Call Optimization](#10-tail-call-optimization)
11. [Cross-Language ABI](#11-cross-language-abi)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What is an ABI?

An ISA defines what instructions exist. An ABI defines how those instructions are used when functions interact.

More precisely, an ABI specifies:

1. **Calling convention**: which registers hold function arguments, return values, and which registers must be preserved across calls
2. **Stack management**: how the stack is organized, the direction it grows, alignment requirements
3. **Data layout**: how structs, unions, and arrays are laid out in memory (including padding and alignment)
4. **Object file format**: how compiled code is stored (ELF on Linux, Mach-O on macOS, PE on Windows)
5. **Name mangling**: how function names are encoded in symbol tables (C++ mangles names to encode overloads; C does not)
6. **Dynamic linking**: how shared libraries are loaded and how function calls to them are resolved at runtime
7. **Exception handling**: how exceptions propagate across stack frames (C++ unwind tables)
8. **System call interface**: the mechanism to call OS services

### ABI vs API

- **API (Application Programming Interface)**: what functions exist and what arguments they take — a source-code level concept
- **ABI (Application Binary Interface)**: how those functions are called at the binary level — a machine code concept

If you compile two libraries with different compilers that follow the same ABI, they can call each other. If they use different ABIs, linking fails or produces crashes.

### Quick Check

> 1. What is the difference between API and ABI?
> 2. List four things that an ABI specifies beyond what the ISA specifies.
> 3. Why would two separately compiled libraries need to follow the same ABI?

---

## 2. Why Calling Conventions Matter

### The Fundamental Problem

When function A calls function B, who owns each register? If both A and B modify `t0` without coordination, one of them will lose their value.

There are three possible approaches:

**Caller-saves all**: The caller saves every register it cares about before calling any function. The callee can freely modify all registers.
- Pro: callee code is simpler (no saves needed)
- Con: caller must save many registers even if callee doesn't use them

**Callee-saves all**: The callee saves all registers it modifies at entry and restores at exit.
- Pro: caller doesn't need to know what callee uses
- Con: callee must save registers it doesn't need to modify

**Split convention**: Some registers are caller-saved, others are callee-saved. This is what all real ABIs use.

### The Modern Split Convention

Real calling conventions divide registers:

**Caller-saved registers** (also called "volatile" or "call-clobbered"):
- The caller must save these before calling a function (if it needs them after the call)
- The callee may freely modify them
- Used for: function arguments, temporary values

**Callee-saved registers** (also called "non-volatile" or "call-preserved"):  
- The callee must save these before using them and restore before returning
- The caller can rely on these being unchanged after a function call
- Used for: long-lived variables that span multiple function calls

### Why This Split?

Consider a loop calling a function on each iteration:
```c
int total = 0;
for (int i = 0; i < n; i++) {
    total += compute(i);   // total and i must survive this call
}
```

The compiler wants to keep `total` and `i` in registers. It puts them in callee-saved registers (e.g., `s0`, `s1` in RISC-V). The `compute` function promises to preserve those registers — so `total` and `i` survive across each call without the loop needing to reload them from memory.

If everything were caller-saved, the loop would have to spill `total` and `i` to the stack on every iteration — enormous overhead.

### Quick Check

> 1. What is the difference between caller-saved and callee-saved registers?
> 2. Why are callee-saved registers useful for long-lived variables?
> 3. If `t0` is a caller-saved register and you need its value after a function call, what must you do?

---

## 3. RISC-V Calling Convention (RISC-V psABI)

The RISC-V Processor-Specific ABI (psABI) defines the standard calling convention for RISC-V:

### Integer Argument and Return Registers

```
Register  ABI Name  Role
────────  ────────  ──────────────────────────────────────────────────────────
x10       a0        Argument 1 / Return value 1
x11       a1        Argument 2 / Return value 2 (for 128-bit returns)
x12       a2        Argument 3
x13       a3        Argument 4
x14       a4        Argument 5
x15       a5        Argument 6
x16       a6        Argument 7
x17       a7        Argument 8 / Syscall number
```

If there are more than 8 integer arguments, additional arguments are pushed onto the stack (in left-to-right order, pushed right-to-left so the first extra argument is at the lowest stack address).

### Callee-Saved and Caller-Saved Split

```
Callee-Saved (preserved across calls):    s0-s11, sp, fp (=s0)
Caller-Saved (may be clobbered by calls): a0-a7, t0-t6, ra, gp, tp
```

The callee saves callee-saved registers it uses at function entry (PUSH) and restores them at exit (POP).

### Float Arguments

The RISC-V F and D extensions add 32 floating-point registers. The ABI defines:

```
Floating-point registers:
  fa0-fa7 (f10-f17)   Float arguments and return values (caller-saved)
  ft0-ft11 (f0-f7, f28-f31) Float temporaries (caller-saved)
  fs0-fs11 (f8-f9, f18-f27) Float saved registers (callee-saved)
```

Float arguments go in fa0-fa7; overflow goes on the stack. A function with mixed int+float args uses both a-registers and fa-registers.

### Stack Pointer Rules

- Stack grows downward (toward lower addresses)
- Stack pointer (sp) must be 16-byte aligned before a function call (`jal`/`call`)
- Functions allocate stack space by subtracting from sp at entry; restore by adding back at exit

### Return Address Convention

The return address is in `ra` (x1). The caller's `jal ra, func` loads `ra` with `PC+4`. The callee returns with `ret` (= `jalr zero, 0(ra)`).

If the callee itself calls any function, it must save `ra` to the stack first (otherwise `jal ra, inner_func` overwrites `ra`).

### Concrete Example

```c
// int bar(int x, int y, int z);
// int foo(int a, int b) {
//     return bar(a+b, b*2, a-b);
// }
```

```riscv
foo:
    # Arguments: a0=a, a1=b
    # We need a and b after computing all three arguments to bar.
    # Since bar is called with a0/a1/a2, we must compute all three args first.
    # a0 and a1 will be overwritten, so save b:
    
    addi  sp, sp, -16          # allocate stack frame (16 bytes, aligned)
    sd    ra, 8(sp)            # save return address
    sd    s0, 0(sp)            # save s0 (callee-saved)
    
    mv    s0, a1               # save b in s0 (callee-saved → survives bar call)
    
    add   t0, a0, a1           # t0 = a+b (first arg for bar)
    add   t1, s0, s0           # t1 = b*2 (second arg for bar)
    sub   t2, a0, s0           # t2 = a-b (third arg for bar)
    
    mv    a0, t0               # arg1 = a+b
    mv    a1, t1               # arg2 = b*2
    mv    a2, t2               # arg3 = a-b
    
    call  bar                  # bar(a+b, b*2, a-b)
    # a0 = bar's return value (leave it in a0, that's our return value too)
    
    ld    ra, 8(sp)            # restore return address
    ld    s0, 0(sp)            # restore s0
    addi  sp, sp, 16           # deallocate stack frame
    ret
```

### Quick Check

> 1. Where does the 9th integer argument go in RISC-V calling convention?
> 2. Why must `ra` be saved to the stack when calling another function?
> 3. What size must the stack frame be aligned to in RISC-V?

---

## 4. System V AMD64 ABI (x86-64 Linux/macOS)

### Integer Arguments and Return

```
Arguments:  rdi, rsi, rdx, rcx, r8, r9   (first 6 integer/pointer arguments)
Return:     rax (primary), rdx (secondary for 128-bit return)
Additional: pushed right-to-left on stack (arg7 first pushed, last to push is arg7,
            so arg7 is at lowest stack address above old rsp)
```

### Register Classification

```
Caller-saved:   rax, rcx, rdx, rsi, rdi, r8, r9, r10, r11
Callee-saved:   rbx, rbp, r12, r13, r14, r15
Stack pointer:  rsp (by convention, always callee saves sp by restoring it)
```

### Red Zone

In the System V AMD64 ABI, user-space functions can use the **128 bytes below rsp** (the "red zone") as a scratch area WITHOUT adjusting rsp. Signal handlers and interrupt handlers must not touch this area.

This allows leaf functions (those that don't call other functions) to avoid adjusting rsp entirely:

```x86asm
# Leaf function using red zone — no stack frame needed!
int_abs:
    mov  eax, edi              # eax = x
    neg  eax                   # eax = -x
    cmovl eax, edi             # if x was positive, use original
    ret
    # Never touched rsp — used red zone for nothing because no locals needed
```

The red zone is an optimization that avoids the overhead of `sub rsp, N` / `add rsp, N` in leaf functions.

### Float and SIMD Arguments (SSE)

```
Float/SIMD:    xmm0-xmm7 (first 8 float/SSE arguments)
Caller-saved:  xmm0-xmm15 (all XMM registers are caller-saved!)
Return:        xmm0 (primary), xmm1 (secondary)
```

A function like `double pow(double x, double y)` takes both arguments in `xmm0` and `xmm1`, returns result in `xmm0`.

### Quick Check

> 1. What registers hold the first 6 arguments in the System V AMD64 ABI?
> 2. What is the "red zone" and what is its purpose?
> 3. Are XMM registers caller-saved or callee-saved in the System V ABI?

---

## 5. Windows x64 ABI

Windows uses a different calling convention from Linux/macOS, even though both run on x86-64!

### Windows x64 Register Usage

```
Arguments:  rcx, rdx, r8, r9   (first 4 integer args — ONLY 4, not 6!)
Float args: xmm0, xmm1, xmm2, xmm3   (first 4 float args)

Return:     rax or xmm0

Callee-saved: rbx, rbp, rdi, rsi, r12, r13, r14, r15, xmm6-xmm15
Caller-saved: rax, rcx, rdx, r8, r9, r10, r11, xmm0-xmm5
```

Key Windows differences from System V:
- Only **4** register arguments (System V: 6)
- Different register set for callee-saved (Windows saves more!)
- No red zone (Windows has no red zone concept)
- Shadow space (Windows requires 32 bytes of stack space reserved before each call)
- XMM6-XMM15 are callee-saved (System V: none are callee-saved)

### Shadow Space

Windows requires 32 bytes of "shadow space" (also called "home space") reserved above the called function's stack frame. This space is for the callee to dump its register arguments for debugging or variadic functions.

```x86asm
# Windows caller:
sub rsp, 40          # 8-byte alignment adjustment + 32 bytes shadow space
mov rcx, first_arg   # arg1
mov rdx, second_arg  # arg2
mov r8, third_arg    # arg3
mov r9, fourth_arg   # arg4
call function
add rsp, 40          # restore stack
```

This shadow space is ALWAYS required, even if the called function doesn't use it.

### Why Windows Differs

Windows chose a more conservative convention (fewer register args) to make debugging and exception handling easier. The shadow space guarantees that every function's arguments are always accessible at known stack locations, even without frame pointer chains.

Cross-platform code that calls Windows DLLs from Linux (via Wine or cross-compilation) must carefully handle ABI differences.

### Quick Check

> 1. How many register arguments does the Windows x64 ABI support?
> 2. What is "shadow space" and why does Windows require it?
> 3. Name two callee-saved registers that differ between System V AMD64 and Windows x64 ABIs.

---

## 6. ARM64 ABI (AAPCS64)

The Procedure Call Standard for AArch64 (AAPCS64) is ARM's official calling convention for 64-bit ARM:

### Register Assignment

```
Arguments:   x0-x7   (8 integer/pointer args)
             v0-v7   (8 float/SIMD args)
Return:      x0 (primary), x1 (secondary), or v0

Caller-saved (temporary): x0-x17, v0-v7, v16-v31
Callee-saved:             x18-x28, x29 (frame pointer), x30 (link register = ra)
                          v8-v15 (lower 8 bytes of each; upper 8 are caller-saved)
```

### Special Registers

```
x29  = frame pointer (FP) — points to the saved FP of the previous frame
x30  = link register (LR) — holds return address (like RISC-V's ra)
sp   = stack pointer (x31 when used as SP, not GPR)
xzr  = zero register (x31 when used as source)
```

### Stack Requirements

- Stack must be 16-byte aligned at all times in AArch64 (stricter than RISC-V's "aligned at call time")
- Violation causes a stack alignment fault on ARM64 hardware

### AAPCS64 vs RISC-V psABI Comparison

| Feature | RISC-V psABI | AAPCS64 |
|---------|-------------|---------|
| Integer arg registers | 8 (a0-a7) | 8 (x0-x7) |
| Float arg registers | 8 (fa0-fa7) | 8 (v0-v7) |
| Return value | a0 (+ a1) | x0 (+ x1) |
| Return address | ra (x1) | x30 (LR) |
| Stack alignment | 16-byte at call | 16-byte always |
| Frame pointer | optional (s0) | conventional (x29) |
| Callee-saved GPRs | s0-s11 | x18-x28, x29, x30 |

---

## 7. Stack Frame Layout

A typical stack frame for a function in RISC-V:

```
Higher addresses
┌───────────────────────────────────────┐
│   Previous frame                      │
├───────────────────────────────────────┤  ← sp before function entry
│   Saved return address (ra/x1)        │  -8(previous sp) or fp+8
│   Saved frame pointer (s0/fp)         │  -16(previous sp) or fp
│   ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─   │
│   Saved callee-saved registers        │
│   (s1-s11, fs0-fs11 as needed)        │
│   ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─   │
│   Local variables                     │
│   (variables that don't fit in regs)  │
│   ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─   │
│   Outgoing arguments (if > 8 args)    │
└───────────────────────────────────────┘  ← sp after function entry (sp = sp - frame_size)
Lower addresses
```

### Frame Pointer (FP) Chain

The frame pointer (FP / s0 in RISC-V, x29 in AArch64) creates a linked list of stack frames that debuggers traverse:

```
FP → [saved FP of caller][saved RA of caller][...frame...]
      ↓
     FP of caller → [saved FP of caller's caller][saved RA of caller's caller]
                      ↓
                     ...
```

Following the FP chain produces a **stack trace** — the sequence of function calls currently active. When a program crashes, the debugger walks this chain to show you `main → function_a → function_b → crash_location`.

Modern compilers with `-O2` often omit the frame pointer (using the `-fomit-frame-pointer` option), saving one register (giving it to the compiler's allocator). Stack traces still work via separate "unwind tables" embedded in the binary (`.eh_frame` section in ELF).

### Quick Check

> 1. Draw the stack frame layout for a function that has 2 local int variables and calls one other function.
> 2. What is the frame pointer used for?
> 3. What happens to the frame pointer chain when a compiler compiles with `-fomit-frame-pointer`?

---

## 8. Passing Large Arguments

### Integers and Pointers

Integer and pointer arguments (8 bytes or smaller) go in registers directly.

### Structs and Unions

The rules for passing structs depend on the struct size:

**Small structs (≤ 16 bytes in RISC-V psABI)**:
- Packed into registers if they fit
- A 8-byte struct goes in one argument register
- A 16-byte struct goes in two argument registers

```c
struct { int x; int y; } point = {1, 2};
foo(point);  // passed as a single 64-bit register: [y=2][x=1] (little-endian packed)
```

**Large structs (> 16 bytes)**:
- Passed via pointer
- Caller allocates space on stack, copies the struct there, passes pointer in argument register

```c
struct BigStruct { int data[10]; };    // 40 bytes > 16
void process(struct BigStruct s);
// Becomes effectively: process(struct BigStruct* s_copy_on_stack);
```

**Return of large structs**:
- Caller allocates space for the return value
- Caller passes a hidden first argument: pointer to the return value space
- Callee writes result through that pointer

```c
struct BigStruct make_struct();
struct BigStruct result = make_struct();
// Compiles to:
// struct BigStruct __result_space;
// make_struct(&__result_space);  // &__result_space passed as hidden arg
// result = __result_space;       (or result IS __result_space if compiler is clever)
```

### Quick Check

> 1. How is a 12-byte struct passed as a function argument in RISC-V?
> 2. How is a 100-byte struct passed?
> 3. What is the "hidden parameter" for returning large structs?

---

## 9. Varargs and Printf-style Functions

Functions with variable argument lists (`printf`, `scanf`) have special ABI treatment:

```c
int printf(const char* format, ...);   // variable arguments
```

### RISC-V Varargs

In RISC-V, `printf("Hello %d %f\n", 42, 3.14)`:
- `a0` = pointer to format string
- `a1` = 42 (first `%d` argument)
- `fa0` = 3.14 (first `%f` argument in float register)

But wait — `printf` doesn't know the types when it's called. It must deduce them from the format string. This means the ABI must ensure all arguments are accessible in a predictable order.

The challenge: if some args go to integer registers (a1-a7) and some to float registers (fa0-fa7), printf must know which register holds which argument.

### The `va_list` Mechanism

```c
#include <stdarg.h>

void my_printf(const char* fmt, ...) {
    va_list args;
    va_start(args, fmt);       // initialize args to point at first vararg
    
    while (*fmt) {
        if (*fmt == '%' && *(fmt+1) == 'd') {
            int val = va_arg(args, int);   // read next int argument
            print_integer(val);
            fmt++;
        }
        // ...
        fmt++;
    }
    
    va_end(args);
}
```

The compiler expands `va_start` and `va_arg` into code that reads arguments from registers in the order they were passed, following the platform's ABI rules.

### Quick Check

> 1. What does `va_start(args, fmt)` do?
> 2. In `printf("x=%d y=%f", 5, 3.14)`, which register holds `5` and which holds `3.14` in RISC-V?
> 3. Why can't varargs functions be inlined easily?

---

## 10. Tail Call Optimization

### What is a Tail Call?

A **tail call** is a function call that is the very last operation in a function — the caller doesn't do anything with the callee's return value except return it.

```c
int foo(int n) {
    return bar(n - 1);   // tail call: bar's return value IS foo's return value
}

int baz(int n) {
    int result = bar(n - 1);
    return result + 1;   // NOT a tail call: adds 1 to bar's result
}
```

### Tail Call Optimization (TCO)

For a tail call, the compiler can optimize by **reusing the current stack frame**:

Without TCO:
```
foo(5) → foo(4) → foo(3) → foo(2) → foo(1) → base_case
Stack: main/foo(5)/foo(4)/foo(3)/foo(2)/foo(1)/base_case
       (n+1 stack frames for depth n)
```

With TCO:
```
foo(5) → foo(4) → foo(3) → ... (reuses same frame!)
Stack: main/foo  (only 2 stack frames regardless of depth!)
```

The optimization works because:
1. The callee doesn't need the caller's frame anymore (return value goes directly to caller's caller)
2. The return address can be the SAME (caller's caller's address)
3. We can jump to the callee instead of calling it, reusing the current stack frame

In RISC-V, a tail call becomes:
```riscv
# Instead of:
call  bar
ret

# Tail call optimization:
j     bar   # or jal zero, bar — just JUMP, don't save return address
            # bar's "ret" returns directly to foo's caller!
```

### When TCO is Crucial: Recursive Algorithms

Functional languages (Haskell, Erlang, Scheme) rely on TCO for recursion — without it, even simple recursive programs would overflow the stack. With TCO, a tail-recursive function is as efficient as a loop.

```c
// Tail-recursive factorial (with accumulator):
long factorial(int n, long acc) {
    if (n <= 1) return acc;
    return factorial(n - 1, n * acc);  // tail call!
}
long result = factorial(1000000, 1);   // safe with TCO (never deep stack)
```

The C standard does NOT require TCO. GCC and Clang implement it opportunistically under `-O2`.

### Quick Check

> 1. What is a "tail call"?
> 2. How does TCO transform the stack usage of a tail-recursive function?
> 3. Why is TCO important for functional programming languages?

---

## 11. Cross-Language ABI

### C as the Lingua Franca

Every major programming language can call C code, and C follows the platform ABI. This means any language that can call C can indirectly interoperate with any other language.

The typical pattern:
```
Rust ──FFI─→ C API ──→ C library
Go   ──CGO─→ C API ──→ C library  
Python ─ctypes/CFFI─→ C API ──→ C library
Java ──JNI─→ C API ──→ C library
```

### Name Mangling in C++

C++ supports function overloading (multiple functions with the same name but different argument types). Since object files identify functions by name, C++ must "mangle" function names to encode type information:

```cpp
void foo(int x);         // C++ mangled: _Z3fooi
void foo(double x);      // C++ mangled: _Z3food
void foo(int x, float y); // C++ mangled: _Z3fooif
```

This is why you must `extern "C"` in C++ headers to expose functions callable from C:

```cpp
extern "C" {
    void my_function(int x);  // exported with C linkage: name = "my_function"
}
```

Without `extern "C"`, C can't find C++ functions because the name is mangled differently by every compiler.

### Quick Check

> 1. Why is C called the "lingua franca" of programming language interoperability?
> 2. What is "name mangling" in C++ and why is it needed?
> 3. What does `extern "C"` do in a C++ header?

---

## Summary

- **ABI** is the binary interface contract: calling convention, data layout, object format, name mangling, exception handling, system call interface.
- **Calling conventions** define which registers are used for arguments, return values, and which must be preserved across calls.
- **Caller-saved** registers: the caller saves if needed. **Callee-saved**: the callee saves if it uses them.
- **RISC-V psABI**: a0-a7 for integer args, fa0-fa7 for float, s0-s11 callee-saved, t0-t6 caller-saved, sp must be 16-byte aligned at call.
- **System V AMD64** (Linux/macOS): rdi/rsi/rdx/rcx/r8/r9 for args, rbx/rbp/r12-r15 callee-saved, 128-byte red zone below rsp.
- **Windows x64**: rcx/rdx/r8/r9 (only 4!), no red zone, 32-byte shadow space required.
- **AAPCS64** (ARM64): x0-x7 for args, x19-x28/x29/x30 callee-saved.
- Large structs are passed by reference (caller allocates, passes pointer). Return of large structs uses hidden pointer argument.
- **Tail call optimization**: reuse current stack frame for tail calls, enabling O(1) stack recursion.
- **Cross-language ABI**: C is the universal FFI layer. C++ uses name mangling; `extern "C"` removes it.

---

## Exercises

### Easy

1. Write the RISC-V calling convention for a function `int64_t add_three(int64_t a, int64_t b, int64_t c)`. Which registers receive a, b, c? Where does the return value go?

2. What is the difference between callee-saved and caller-saved registers? Give one example of each from RISC-V.

3. What is "tail call optimization" and how does it change stack frame management?

### Medium

4. **Stack frame construction**: Write the RISC-V assembly prologue and epilogue for a function `foo(int a, int b, int c)` that calls two other functions and uses one local variable stored in `s0`. Show: stack allocation, saving ra and s0, restoring, and returning.

5. **ABI portability bug**: A developer writes a shared library in C++ and exports a function without `extern "C"`. The library is compiled with GCC. Another developer tries to call this function from MSVC-compiled code on Windows. Explain all the ways this can go wrong: (a) name mangling difference, (b) calling convention difference (System V vs Windows x64), (c) C++ runtime ABI differences. What must the developer do to make the library portable?

6. **Varargs internals**: The `vprintf` function takes a format string and a `va_list`. The `va_list` on x86-64 (System V ABI) is defined in `<stdarg.h>` as a struct:
   ```c
   typedef struct {
       unsigned int gp_offset;     // offset to next integer register arg
       unsigned int fp_offset;     // offset to next float register arg
       void* overflow_arg_area;    // pointer to stack args overflow area
       void* reg_save_area;        // pointer to saved registers
   } va_list[1];
   ```
   Explain how `va_arg(list, int)` uses this structure to get the next integer argument. What happens when integer register arguments are exhausted?

### Hard

7. **ABI stability and shared libraries**: Shared libraries (.so / .dylib) have their ABI baked into the binary. If you add a parameter to a function in a shared library and recompile it, all programs that linked against the old version will crash (or worse, silently compute wrong results). Linux distributions manage this with "soname" versioning. Research: (a) what is the "soname" in an ELF shared library and how does it relate to ABI versioning? (b) How does `SOVERSION` in CMake control this? (c) What techniques can you use to add a new function to a library without breaking existing binaries? (d) How does the "symbol versioning" feature (GNU ld `--version-script`) allow maintaining multiple ABI versions simultaneously in one .so file?

8. **Implement a calling convention from scratch**: Design a minimal calling convention for a hypothetical 8-register RISC processor where registers are R0-R7. Your design must handle: (a) up to 4 function arguments, (b) one return value, (c) nested function calls without stack corruption, (d) saving and restoring registers without using 3 or more callee-saved registers for simple functions. Write out the complete rules (which registers for what purpose) and justify each choice. Then trace through a two-level function call (main calls foo, foo calls bar) to verify your convention works.
