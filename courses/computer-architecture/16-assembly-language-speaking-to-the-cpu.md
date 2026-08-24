# Chapter 16: Assembly Language — Speaking to the CPU

Assembly language is the lowest level of human-readable programming. Each assembly instruction corresponds to exactly one machine instruction — there is virtually no abstraction. Learning to read assembly gives you a direct window into what the CPU is doing and why high-level programs behave the way they do. This chapter uses RISC-V assembly as the primary example, with x86-64 shown for comparison.

## Table of Contents

1. [What is Assembly Language?](#1-what-is-assembly-language)
2. [RISC-V Assembly Basics](#2-risc-v-assembly-basics)
3. [Arithmetic and Logic Operations](#3-arithmetic-and-logic-operations)
4. [Memory: Loads and Stores](#4-memory-loads-and-stores)
5. [Control Flow: Branches and Jumps](#5-control-flow-branches-and-jumps)
6. [A Complete Program: Sum 1 to N](#6-a-complete-program-sum-1-to-n)
7. [Functions: Call and Return](#7-functions-call-and-return)
8. [The Stack Frame](#8-the-stack-frame)
9. [Assembly for C Patterns](#9-assembly-for-c-patterns)
10. [x86-64 Assembly: The Comparison](#10-x86-64-assembly-the-comparison)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What is Assembly Language?

### Machine Code vs Assembly

The CPU executes **machine code**: binary numbers (0s and 1s). For example, the RISC-V instruction `ADD x5, x1, x2` is encoded as:

```
Binary:  0000000 00010 00001 000 00101 0110011
Hex:     0x00208293

Field breakdown:
funct7 = 0000000   (ADD, not SUB)
rs2    = 00010     (x2)
rs1    = 00001     (x1)
funct3 = 000
rd     = 00101     (x5)
opcode = 0110011   (R-type)
```

Nobody writes binary by hand. **Assembly language** is a textual representation where:
- Each binary instruction has a mnemonic: `ADD`, `LW`, `BEQ`
- Registers have names: `x0`, `x1`, ... or symbolic aliases like `zero`, `ra`, `t0`
- Memory addresses can be symbolic labels: `loop:`, `main:`, `data:`

An **assembler** (a program like `as` or `riscv-gnu-assembler`) converts assembly text to binary machine code.

```
Assembly text                   Machine code (binary)
─────────────────────────────   ─────────────────────────────────
add  t0, a0, a1              →  0x00a50233
lw   t1, 0(a0)               →  0x0005002b  (approximately)
beq  t0, zero, done          →  0x00028463  (approximately)
```

### Why Learn Assembly?

1. **Debugging**: assembly shows exactly what the program does — no hidden magic
2. **Performance**: inner loops, SIMD, crypto code is often written in assembly for maximum performance
3. **Security**: understanding exploits requires reading assembly (buffer overflow, return-oriented programming)
4. **Embedded systems**: microcontrollers with kilobytes of memory may require hand-written assembly
5. **Compilers**: if you write a compiler, you generate assembly

### Quick Check

> 1. What is the difference between assembly language and machine code?
> 2. What is an assembler and what does it produce?
> 3. Name three reasons a programmer would write assembly directly instead of C.

---

## 2. RISC-V Assembly Basics

### Registers

RISC-V has 32 integer registers (x0-x31). Each is 64 bits wide in RV64:

```
Register  ABI Name   Role                        Saved by?
────────  ─────────  ─────────────────────────── ─────────
x0        zero       Always 0 (hardwired)         N/A
x1        ra         Return address               Caller
x2        sp         Stack pointer                Callee
x3        gp         Global pointer               N/A
x4        tp         Thread pointer               N/A
x5-x7     t0-t2      Temporaries                  Caller
x8        s0/fp      Saved register / frame ptr   Callee
x9        s1         Saved register               Callee
x10-x11   a0-a1      Fn args / return values      Caller
x12-x17   a2-a7      Function arguments           Caller
x18-x27   s2-s11     Saved registers              Callee
x28-x31   t3-t6      Temporaries                  Caller
```

**The golden rule**: `x0` always reads as 0. Any write to `x0` is discarded. This is incredibly useful — you never need a "clear register" instruction:
```
addi x5, x0, 42     # x5 = 0 + 42 = 42   (load immediate)
add  x5, x0, x6     # x5 = 0 + x6 = x6   (register move!)
beq  x5, x0, done   # branch if x5 == 0   (compare to zero)
```

### Basic Instruction Syntax

```
mnemonic  rd, rs1, rs2      # R-type: rd = rs1 OP rs2
mnemonic  rd, rs1, imm      # I-type: rd = rs1 OP imm
mnemonic  rs2, imm(rs1)     # Store:  Memory[rs1+imm] = rs2
mnemonic  rs1, rs2, label   # Branch: if rs1 OP rs2 goto label
```

### Pseudo-Instructions

RISC-V assemblers support pseudo-instructions that expand to one or more real instructions:

```
Pseudo           Expands to
──────────────── ──────────────────────────────────────────
li  t0, 100      addi t0, zero, 100      (small immediate)
li  t0, 0x12345  lui t0, 0x12; addi t0, t0, 0x345 (large)
mv  t0, t1       addi t0, t1, 0          (register move)
neg t0, t1       sub  t0, zero, t1       (negate)
nop              addi zero, zero, 0      (no operation)
j   label        jal  zero, label        (unconditional jump)
ret              jalr zero, 0(ra)        (return from function)
call func        auipc ra, ...; jalr ra, ra, ... (function call)
```

### Quick Check

> 1. What does `addi x5, x0, 42` do? Why is `x0` useful here?
> 2. What does the `mv t0, t1` pseudo-instruction expand to?
> 3. What is the difference between `t` registers and `s` registers in the calling convention?

---

## 3. Arithmetic and Logic Operations

### Integer Arithmetic

```riscv
# Addition
add  t0, t1, t2      # t0 = t1 + t2
addi t0, t1, 10      # t0 = t1 + 10   (I-type, imm is signed 12-bit: -2048 to +2047)

# Subtraction
sub  t0, t1, t2      # t0 = t1 - t2   (no SUBI: use addi with negative immediate)
addi t0, t1, -5      # t0 = t1 - 5    (equivalent to SUBI t0, t1, 5)

# Multiplication (requires M extension)
mul  t0, t1, t2      # t0 = (t1 × t2) lower 64 bits
mulh t0, t1, t2      # t0 = (t1 × t2) upper 64 bits (signed)

# Division (requires M extension)
div  t0, t1, t2      # t0 = t1 ÷ t2   (signed integer division)
rem  t0, t1, t2      # t0 = t1 mod t2 (remainder)
```

### Bitwise Logic

```riscv
and  t0, t1, t2      # t0 = t1 AND t2
andi t0, t1, 0xFF    # t0 = t1 AND 0xFF   (extract lowest byte)
or   t0, t1, t2      # t0 = t1 OR t2
ori  t0, t1, 0x01    # t0 = t1 OR 1      (set bit 0)
xor  t0, t1, t2      # t0 = t1 XOR t2
xori t0, t1, -1      # t0 = NOT t1       (-1 in 12-bit imm = all 1s → XOR flips all bits)
```

### Shifts

```riscv
sll  t0, t1, t2      # t0 = t1 << t2    (shift left logical, fills with 0)
slli t0, t1, 3       # t0 = t1 << 3     (multiply by 8)
srl  t0, t1, t2      # t0 = t1 >> t2    (logical shift right, fills with 0)
srli t0, t1, 2       # t0 = t1 >> 2     (divide by 4, unsigned)
sra  t0, t1, t2      # t0 = t1 >> t2    (arithmetic shift right, fills with sign bit)
srai t0, t1, 2       # t0 = t1 >> 2     (divide by 4, signed)
```

### Compare (Set Less Than)

RISC-V doesn't have a flags register like x86. Instead, it has compare instructions that write a 0 or 1 to a register:

```riscv
slt  t0, t1, t2      # t0 = (t1 < t2) ? 1 : 0   (signed)
sltu t0, t1, t2      # t0 = (t1 < t2) ? 1 : 0   (unsigned)
slti t0, t1, 5       # t0 = (t1 < 5) ? 1 : 0
```

### Quick Check

> 1. How do you implement `t0 = t1 - 5` in RISC-V (there is no SUBI)?
> 2. How do you logically NOT a register in RISC-V (there is no NOT instruction)?
> 3. What is the difference between SRL and SRA? When would you use each?

---

## 4. Memory: Loads and Stores

### Load Instructions

```riscv
# Load into rd from Memory[rs1 + imm]
lb   t0, 0(a0)       # load byte (8-bit), sign-extended to 64 bits
lbu  t0, 0(a0)       # load byte, zero-extended (unsigned)
lh   t0, 0(a0)       # load halfword (16-bit), sign-extended
lhu  t0, 0(a0)       # load halfword, zero-extended
lw   t0, 0(a0)       # load word (32-bit), sign-extended
lwu  t0, 0(a0)       # load word, zero-extended
ld   t0, 0(a0)       # load doubleword (64-bit)
```

### Store Instructions

```riscv
# Store rs2 to Memory[rs1 + imm]
sb   t0, 0(a0)       # store lowest byte of t0
sh   t0, 0(a0)       # store lowest halfword of t0
sw   t0, 0(a0)       # store lowest word of t0
sd   t0, 0(a0)       # store full doubleword of t0
```

### Memory Access Examples

```riscv
# int array[10] accessed through a0 (pointer to array)
# Access array[0]:
lw   t0, 0(a0)       # t0 = array[0]

# Access array[3] (offset = 3 * 4 bytes = 12):
lw   t0, 12(a0)      # t0 = array[3]

# Access array[i] where i is in t1:
slli t1, t1, 2       # t1 = i * 4  (each int is 4 bytes)
add  t2, a0, t1      # t2 = &array[i]
lw   t0, 0(t2)       # t0 = array[i]

# struct access:
# struct Point { int x; int y; }  —  x at offset 0, y at offset 4
# Pointer to struct in a0
lw   t0, 0(a0)       # t0 = point.x
lw   t1, 4(a0)       # t1 = point.y
sw   t0, 0(a1)       # another_point.x = t0
```

### Quick Check

> 1. What is the difference between `lb` and `lbu`?
> 2. How do you access `array[5]` if the array base is in register `a0` and each element is 4 bytes?
> 3. Why does RISC-V have separate instructions for 8-bit, 16-bit, 32-bit, and 64-bit loads?

---

## 5. Control Flow: Branches and Jumps

### Branch Instructions (Conditional)

All RISC-V branches compare two registers and jump to a PC-relative offset if the condition is true:

```riscv
beq  rs1, rs2, label   # branch if rs1 == rs2
bne  rs1, rs2, label   # branch if rs1 != rs2
blt  rs1, rs2, label   # branch if rs1 < rs2   (signed)
bge  rs1, rs2, label   # branch if rs1 >= rs2  (signed)
bltu rs1, rs2, label   # branch if rs1 < rs2   (unsigned)
bgeu rs1, rs2, label   # branch if rs1 >= rs2  (unsigned)

# Pseudo-instructions using zero register:
beqz rs1, label        # branch if rs1 == 0    → beq rs1, zero, label
bnez rs1, label        # branch if rs1 != 0    → bne rs1, zero, label
bltz rs1, label        # branch if rs1 < 0     → blt rs1, zero, label
bgtz rs1, label        # branch if rs1 > 0     → blt zero, rs1, label
```

### Jump Instructions (Unconditional)

```riscv
jal  rd, label         # Jump And Link: rd = PC+4; PC = PC + offset
                       # Used for function calls when rd = ra
                       # Used for unconditional jump when rd = zero (j pseudo-insn)

jalr rd, imm(rs1)      # Jump And Link Register: rd = PC+4; PC = rs1 + imm
                       # Used for function return: jalr zero, 0(ra)  → ret
                       # Used for indirect function calls (function pointers)
```

The `jal` instruction saves `PC+4` (the address of the next instruction) into `rd`. For function calls, `rd = ra` so the callee can return. For unconditional jumps, `rd = zero` (discards the return address).

### Quick Check

> 1. How does `jal` support both function calls and unconditional jumps?
> 2. What does `beq t0, zero, end` do?
> 3. Why does RISC-V have both `jal` (PC-relative) and `jalr` (register-based) jump instructions?

---

## 6. A Complete Program: Sum 1 to N

Let's write a program that computes sum = 1 + 2 + 3 + ... + N, where N is passed in register a0.

```riscv
# int sum_to_n(int n)
# Input:  a0 = n
# Output: a0 = sum (1 + 2 + ... + n)
# Uses:   t0 (current i), t1 (sum)

sum_to_n:
    # Initialize: sum = 0, i = 1
    li   t1, 0          # t1 = 0  (sum)
    li   t0, 1          # t0 = 1  (i = 1)

loop:
    # while (i <= n) { sum += i; i++; }
    bgt  t0, a0, done   # if i > n, exit loop
    add  t1, t1, t0     # sum += i
    addi t0, t0, 1      # i++
    j    loop           # goto loop

done:
    mv   a0, t1         # return sum (in a0)
    ret                 # return to caller (jalr zero, 0(ra))
```

### Trace for N = 4

```
Cycle  t0 (i)  t1 (sum)  Instruction                 Action
─────  ──────  ─────────  ────────────────────────── ──────────────────
 1      1        0        li t1, 0                   sum = 0
 2      1        0        li t0, 1                   i = 1
 3      1        0        bgt t0, a0, done  (1>4? N) no branch
 4      1        1        add t1, t1, t0             sum = 0+1 = 1
 5      2        1        addi t0, t0, 1             i = 2
 6      2        1        j loop                     goto loop
 7      2        1        bgt t0, a0, done  (2>4? N) no branch
 8      2        3        add t1, t1, t0             sum = 1+2 = 3
 9      3        3        addi t0, t0, 1             i = 3
10      3        3        j loop                     goto loop
11      3        3        bgt t0, a0, done  (3>4? N) no branch
12      3        6        add t1, t1, t0             sum = 3+3 = 6
13      4        6        addi t0, t0, 1             i = 4
14      4        6        j loop                     goto loop
15      4        6        bgt t0, a0, done  (4>4? N) no branch
16      4       10        add t1, t1, t0             sum = 6+4 = 10
17      5       10        addi t0, t0, 1             i = 5
18      5       10        j loop                     goto loop
19      5       10        bgt t0, a0, done  (5>4? Y) BRANCH TAKEN
20      -       10        mv a0, t1                  a0 = 10
21      -       10        ret                        return 10
```

Result: a0 = 10 (= 1+2+3+4). Correct!

### Quick Check

> 1. In the sum_to_n program, what registers hold the function result at return?
> 2. Trace the program for N=2. How many iterations does the loop run?
> 3. If you changed `bgt` to `bge`, would the result be different? Why?

---

## 7. Functions: Call and Return

### The Function Call Mechanism

When calling a function in RISC-V:

1. **Caller** puts arguments in a0-a7 (up to 8 arguments)
2. **Caller** saves its t-registers (if needed — they're caller-saved)
3. **Caller** executes `jal ra, function_name` — this saves PC+4 in ra (x1) and jumps to function
4. **Callee** executes the function body
5. **Callee** puts return value in a0 (or a0+a1 for 128-bit results)
6. **Callee** executes `ret` (= `jalr zero, 0(ra)`) — PC = ra (return to caller)
7. **Caller** uses return value from a0

### What if the Callee Calls Another Function?

Problem: if `sum_to_n` calls another function, `jal ra, helper` will overwrite `ra` — destroying the return address from the caller!

Solution: the callee must **save ra on the stack** before calling any other function.

```riscv
# A function that calls another function
fib:                             # int fib(int n)
    addi sp, sp, -16             # allocate 16 bytes on stack
    sd   ra, 8(sp)               # save return address
    sd   s0, 0(sp)               # save s0 (callee-saved)
    mv   s0, a0                  # save n in s0 (a0 will be clobbered by recursive calls)

    li   t0, 1
    ble  s0, t0, base_case       # if n <= 1, return n

    addi a0, s0, -1              # argument = n-1
    call fib                     # ra = fib(n-1)
    mv   s1, a0                  # save fib(n-1) result in s1

    addi a0, s0, -2              # argument = n-2
    call fib                     # a0 = fib(n-2)
    add  a0, a0, s1              # a0 = fib(n-1) + fib(n-2)
    j    done

base_case:
    mv   a0, s0                  # return n (which is 0 or 1)

done:
    ld   ra, 8(sp)               # restore return address
    ld   s0, 0(sp)               # restore s0
    addi sp, sp, 16              # deallocate stack frame
    ret
```

### Quick Check

> 1. Where does `jal ra, func` save the return address?
> 2. Why must a function that calls other functions save `ra` on the stack?
> 3. Why are a0-a7 "caller-saved" registers? Who is responsible for preserving their values?

---

## 8. The Stack Frame

The stack is a LIFO (Last In, First Out) region of memory that grows downward in memory. The stack pointer `sp` (x2) always points to the top of the stack.

```
High addresses
┌─────────────────────────────────────────────────┐
│  main's stack frame                             │
│  ...                                            │
├─────────────────────────────────────────────────┤  ← sp (before calling fib)
│  fib's stack frame:                             │
│  [sp+8]  = saved ra                             │
│  [sp+0]  = saved s0                             │
├─────────────────────────────────────────────────┤  ← sp (inside fib)
│  fib(n-1)'s stack frame (during recursive call) │
├─────────────────────────────────────────────────┤  ← sp (in recursive call)
│  ...                                            │
Low addresses (stack grows down)
```

### Stack Frame Layout

A typical RISC-V stack frame contains:
1. Saved callee-saved registers (s0-s11) that the function modifies
2. Saved `ra` (if the function calls other functions)
3. Local variables that don't fit in registers
4. Space for spilled registers

The frame size must be a multiple of 16 bytes (RISC-V ABI requirement for alignment).

### Quick Check

> 1. Which direction does the stack grow in memory (toward higher or lower addresses)?
> 2. List three things stored in a typical stack frame.
> 3. What happens if the stack pointer is not aligned to 16 bytes on AArch64 (ARM)?

---

## 9. Assembly for C Patterns

### if-else

```c
// C
if (x > 0) { result = x; } else { result = -x; }
```

```riscv
# a0 = x
    blez a0, else_branch     # if x <= 0, go to else
    mv   a1, a0              # result = x (then-branch)
    j    endif
else_branch:
    neg  a1, a0              # result = -x (else-branch)
endif:
    # a1 = result
```

### for Loop

```c
// C
int sum = 0;
for (int i = 0; i < 10; i++) sum += i;
```

```riscv
    li   t0, 0               # sum = 0
    li   t1, 0               # i = 0
    li   t2, 10              # loop limit
loop:
    bge  t1, t2, done        # if i >= 10, exit
    add  t0, t0, t1          # sum += i
    addi t1, t1, 1           # i++
    j    loop
done:
    # t0 = sum (= 45)
```

### while with Array

```c
// C: find first element > threshold
int* find_first(int* arr, int n, int threshold) {
    for (int i = 0; i < n; i++)
        if (arr[i] > threshold) return &arr[i];
    return NULL;
}
```

```riscv
# a0 = arr, a1 = n, a2 = threshold
find_first:
    li   t0, 0               # i = 0
loop:
    bge  t0, a1, not_found   # if i >= n, return NULL
    slli t1, t0, 2           # t1 = i * 4 (offset in bytes)
    add  t2, a0, t1          # t2 = &arr[i]
    lw   t3, 0(t2)           # t3 = arr[i]
    bgt  t3, a2, found       # if arr[i] > threshold, found it
    addi t0, t0, 1           # i++
    j    loop
found:
    mv   a0, t2              # return &arr[i]
    ret
not_found:
    li   a0, 0               # return NULL
    ret
```

---

## 10. x86-64 Assembly: The Comparison

The same sum-to-N function in x86-64 AT&T syntax (used by GAS assembler):

```x86asm
# int sum_to_n(int n)
# Input:  %edi = n  (first argument in System V AMD64 ABI)
# Output: %eax = sum

sum_to_n:
    xor  %eax, %eax          # eax = 0 (sum)
    test %edi, %edi          # set flags based on n
    jle  .done               # if n <= 0, done
    mov  $1, %ecx            # ecx = 1 (i)
.loop:
    add  %ecx, %eax          # sum += i
    inc  %ecx                # i++
    cmp  %edi, %ecx          # compare i to n
    jle  .loop               # if i <= n, continue
.done:
    ret
```

Key differences from RISC-V:
- **Implicit flags register (EFLAGS)**: `test`, `cmp`, `add`, `inc` all modify flags; `jle`, `jg` read them. In RISC-V, branches compare registers directly.
- **Variable instruction length**: `xor %eax, %eax` is 2 bytes; `mov $1, %ecx` is 5 bytes.
- **Fewer registers**: x86-64 has 16 GPRs (rax, rbx, rcx, rdx, rsi, rdi, rsp, rbp, r8-r15) vs RISC-V's 32.
- **Memory-register operations**: x86 can `add` directly from/to memory without explicit load/store.
- **`ret` uses the stack**: x86 `ret` pops the return address from the stack (no ra register).

---

## Summary

- **Assembly language** is a human-readable representation of machine code, where each line corresponds to one CPU instruction.
- **RISC-V registers**: 32 general-purpose registers (x0-x31), with ABI names (zero, ra, sp, t0-t6, a0-a7, s0-s11). x0 is hardwired to 0.
- **Load-store architecture**: only `lw`/`ld`/`lh`/`lb` and `sw`/`sd`/`sh`/`sb` access memory. ALU ops (add, sub, and...) use only registers.
- **Control flow**: conditional branches (beq, bne, blt...) and jumps (jal, jalr). `jal ra, label` is a function call; `ret` (= `jalr zero, 0(ra)`) is a return.
- **Calling convention**: arguments in a0-a7, return value in a0, return address in ra, stack pointer in sp. Callee saves s-registers and ra (if it calls further); caller saves t-registers.
- **Stack frame**: allocated by subtracting from sp; stores saved registers, local variables, return address.
- **x86-64 differences**: implicit flags register, variable-length encoding, fewer GPRs, memory operands in ALU instructions.

---

## Exercises

### Easy

1. Translate to RISC-V assembly:
   - `c = a + b;` (a in t0, b in t1, c in t2)
   - `x = x * 8;` (x in t0, use shift)
   - `if (a == 0) b = 1;` (a in a0, b in a1)

2. Trace the execution of `sum_to_n` for N=3. What are the values of t0 and t1 after each loop iteration?

3. What does `xori t0, t1, -1` compute? Why is this equivalent to bitwise NOT?

### Medium

4. Write RISC-V assembly for this C code:
   ```c
   int max(int* arr, int n) {
       int m = arr[0];
       for (int i = 1; i < n; i++)
           if (arr[i] > m) m = arr[i];
       return m;
   }
   ```
   Assume a0 = arr (pointer), a1 = n. Return result in a0. Show your register allocation.

5. The RISC-V calling convention says a0-a7 are "caller-saved" and s0-s11 are "callee-saved." In your `fib` implementation, why must you save `s0` to the stack? What would go wrong if you used `t0` instead of `s0` to save n?

6. **Stack overflow**: recursion uses the stack. If each `fib` call uses 16 bytes of stack, and the system has 8 MB of stack, what is the maximum recursion depth? What happens if you call `fib(1000000)` (which requires ~10^6 stack frames)?

### Hard

7. **Compiler output analysis**: compile this C code with `riscv64-unknown-elf-gcc -O0 -S test.c` (or use the Godbolt compiler explorer with RISC-V target, `-O0`):
   ```c
   int factorial(int n) {
       if (n <= 1) return 1;
       return n * factorial(n-1);
   }
   ```
   Then compile with `-O2`. Compare the two assembly outputs:
   (a) What extra code does `-O0` generate (e.g., stack frame setup)?
   (b) Does `-O2` convert the recursion to iteration (tail call optimization)? Why is the compiler free to do this?
   (c) Count the instructions at `-O0` vs `-O2`. What is the ratio?

8. **Position-independent code**: RISC-V uses PC-relative addressing for all branches and the `auipc` + `addi` sequence for loading addresses. Why must shared libraries use position-independent code? What would happen if a shared library used absolute addresses? How does Linux's dynamic linker (ld.so) resolve function calls between a main executable and a shared library at runtime?
