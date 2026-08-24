# Chapter 58: System Calls

> **"A system call is the only legal way for user-mode code to ask the kernel to do something privileged. Open a file, allocate memory, create a process — all of these require the CPU to switch from Ring 3 to Ring 0, execute kernel code, then safely return. The system call interface is the API boundary between your kernel and the programs it runs."**

---

## Table of Contents

1. [Why System Calls Exist](#1-why-system-calls-exist)
2. [How INT 0x80 Works](#2-how-int-0x80-works)
3. [System Call Numbering and Arguments](#3-system-call-numbering-and-arguments)
4. [The System Call Handler in Assembly](#4-the-system-call-handler-in-assembly)
5. [The System Call Dispatch Table in C](#5-the-system-call-dispatch-table-in-c)
6. [Implementing System Calls](#6-implementing-system-calls)
7. [System Call Return Values](#7-system-call-return-values)
8. [User-Side System Call Wrappers](#8-user-side-system-call-wrappers)
9. [Testing System Calls](#9-testing-system-calls)
10. [Summary](#summary)

---

## 1. Why System Calls Exist

```
Problem: user processes must NOT directly execute privileged operations.
  
  If user code could freely:
    - Write to hardware ports (I/O)
    - Modify page tables (CR3 writes)
    - Disable interrupts (CLI)
    - Read other processes' memory
    
  → No security, no isolation, no protection
  
Solution: User code runs in Ring 3 (restricted).
  Privileged operations are only available in Ring 0 (kernel).
  To request a privileged operation, user code must cross the ring boundary
  in a CONTROLLED way: the system call gate.
  
System call = a defined entry point into the kernel.
The kernel controls which operations are available, validates arguments,
performs the operation, and returns safely to Ring 3.

Linux: ~350 system calls (open, read, write, fork, exec, exit, ...)
Windows: ~500 NT native API calls (NtCreateFile, NtReadFile, ...)
Our TinyOS: ~15 system calls for learning purposes
```

---

## 2. How INT 0x80 Works

Linux historically used `int $0x80` for system calls. We do the same:

```
User code:
  mov eax, 4        ; syscall number: sys_write
  mov ebx, 1        ; arg 1: file descriptor (1 = stdout)
  mov ecx, message  ; arg 2: buffer address
  mov edx, 13       ; arg 3: length
  int 0x80          ; SYSCALL!
  
  ; After return: EAX = return value (bytes written, or -errno on error)

What happens inside the CPU when int $0x80 fires:
  1. CPU checks IDT[0x80].DPL — is the caller allowed to invoke this?
     (We set DPL=3 for vector 0x80, so Ring 3 can use it)
  2. CPU switches to Ring 0:
     - Loads SS0 and ESP0 from TSS (kernel stack)
     - Pushes old SS, old ESP, EFLAGS, old CS, old EIP onto kernel stack
  3. Clears IF (disables interrupts) — interrupt gate
  4. Jumps to IDT[0x80].handler
  5. Handler runs with Ring 0 privileges
  6. iret restores old CS (Ring 3), EIP, EFLAGS, ESP, SS
  7. User code continues after int $0x80 instruction
```

**Registering the syscall gate:**
```c
/* In idt_init(), after exception handlers: */
extern void syscall_stub(void);

/* Vector 0x80, DPL=3 (user can invoke), interrupt gate: */
/* flags = 0xEE = 1110 1110 = P=1, DPL=11, type=0xE (32-bit interrupt gate) */
idt_set_gate(0x80, syscall_stub, GDT_KERNEL_CODE, 0xEE);
```

---

## 3. System Call Numbering and Arguments

```c
/* include/syscall_numbers.h */
#pragma once

/* System call numbers: */
#define SYS_EXIT      1    /* exit(status) */
#define SYS_WRITE     2    /* write(fd, buf, count) */
#define SYS_READ      3    /* read(fd, buf, count) */
#define SYS_OPEN      4    /* open(path, flags) */
#define SYS_CLOSE     5    /* close(fd) */
#define SYS_GETPID    6    /* getpid() */
#define SYS_SLEEP     7    /* sleep(ms) */
#define SYS_FORK      8    /* fork() */
#define SYS_EXEC      9    /* exec(path, argv) */
#define SYS_WAITPID   10   /* waitpid(pid) */
#define SYS_BRK       11   /* brk(addr) — extend heap */
#define SYS_YIELD     12   /* yield() */
#define SYS_PUTCHAR   13   /* putchar(c) — simple debug output */
#define SYS_GETTIME   14   /* gettime() — return timer ticks */

#define SYSCALL_MAX   15

/* Calling convention:
   EAX = syscall number
   EBX = arg 1
   ECX = arg 2
   EDX = arg 3
   ESI = arg 4 (if needed)
   EDI = arg 5 (if needed)
   Return value in EAX
   Error: EAX = -1, set errno (we keep it simple: just return negative errno) */
```

---

## 4. The System Call Handler in Assembly

```nasm
; boot/syscall_stub.asm — System call entry point

[BITS 32]
[GLOBAL syscall_stub]
[EXTERN syscall_handler]

syscall_stub:
    ; When we get here (from int $0x80):
    ; CPU has pushed: old SS, old ESP, EFLAGS, old CS, old EIP
    ; We're now on the kernel stack (TSS.esp0)
    ; Interrupts are disabled (interrupt gate)
    ; Ring 0 is active (CS loaded from IDT gate)
    
    ; Save all registers:
    pusha               ; saves EAX, ECX, EDX, EBX, ESP, EBP, ESI, EDI
    push ds
    push es
    push fs
    push gs
    
    ; Load kernel data segments:
    mov ax, 0x10
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    
    ; Call C syscall handler with the registers struct pointer:
    push esp            ; argument: pointer to registers on stack
    call syscall_handler
    add esp, 4
    
    ; The return value is in EAX.
    ; Store it into the EAX slot in the saved registers
    ; (so when we popa, the user gets EAX = return value):
    mov [esp + 28], eax  ; offset to EAX in the pusha frame
    
    ; Restore segments and registers:
    pop gs
    pop fs
    pop es
    pop ds
    popa
    
    ; Return to user mode:
    iret
```

**The register layout on the kernel stack (after pusha+segs):**
```
esp+0:  gs
esp+4:  fs
esp+8:  es
esp+12: ds
esp+16: edi   (from pusha)
esp+20: esi
esp+24: ebp
esp+28: esp (old, pushed by pusha — not user ESP)
esp+32: ebx
esp+36: edx
esp+40: ecx
esp+44: eax   ← syscall number on entry; return value on exit
esp+48: eip   (saved by CPU)
esp+52: cs
esp+56: eflags
esp+60: esp   (old user-mode ESP)
esp+64: ss    (old user-mode SS)
```

So to write the return value into EAX: `mov [esp + 44], eax` (or use the registers_t struct offset).

---

## 5. The System Call Dispatch Table in C

```c
/* kernel/syscall.c */

#include "syscall.h"
#include "syscall_numbers.h"
#include "process.h"
#include "scheduler.h"
#include "vga.h"
#include "timer.h"
#include "isr.h"

/* System call function pointer type: */
typedef int32_t (*syscall_fn)(registers_t *);

/* Forward declarations: */
static int32_t sys_exit(registers_t *r);
static int32_t sys_write(registers_t *r);
static int32_t sys_read(registers_t *r);
static int32_t sys_getpid(registers_t *r);
static int32_t sys_sleep(registers_t *r);
static int32_t sys_yield(registers_t *r);
static int32_t sys_putchar(registers_t *r);
static int32_t sys_gettime(registers_t *r);
static int32_t sys_brk(registers_t *r);

/* Dispatch table indexed by syscall number: */
static syscall_fn syscall_table[SYSCALL_MAX] = {
    [0]            = NULL,           /* 0 = unused */
    [SYS_EXIT]     = sys_exit,
    [SYS_WRITE]    = sys_write,
    [SYS_READ]     = sys_read,
    [SYS_OPEN]     = NULL,           /* Implemented in Ch 61 */
    [SYS_CLOSE]    = NULL,
    [SYS_GETPID]   = sys_getpid,
    [SYS_SLEEP]    = sys_sleep,
    [SYS_FORK]     = NULL,           /* Advanced — omit for now */
    [SYS_EXEC]     = NULL,
    [SYS_WAITPID]  = NULL,
    [SYS_BRK]      = sys_brk,
    [SYS_YIELD]    = sys_yield,
    [SYS_PUTCHAR]  = sys_putchar,
    [SYS_GETTIME]  = sys_gettime,
};

/* Main syscall C handler — called from syscall_stub: */
int32_t syscall_handler(registers_t *regs) {
    uint32_t num = regs->eax;
    
    if (num >= SYSCALL_MAX || !syscall_table[num]) {
        kprintf("Unknown syscall %u from PID %u\n",
                num, current_process ? current_process->pid : 0);
        return -1;
    }
    
    return syscall_table[num](regs);
}
```

---

## 6. Implementing System Calls

```c
/* Exit: terminate current process */
static int32_t sys_exit(registers_t *r) {
    int exit_code = (int)r->ebx;
    kprintf("[PID %u] exit(%d)\n", current_process->pid, exit_code);
    process_exit(exit_code);
    /* Never returns */
    return 0;
}

/* Write: write to file descriptor */
static int32_t sys_write(registers_t *r) {
    int          fd    = (int)r->ebx;
    const char  *buf   = (const char *)r->ecx;
    uint32_t     count = r->edx;
    
    /* Basic validation: */
    if (!buf || count == 0) return -1;
    
    if (fd == 1 || fd == 2) {   /* stdout or stderr → screen */
        for (uint32_t i = 0; i < count; i++) {
            terminal_putchar(buf[i]);
        }
        return (int32_t)count;
    }
    
    return -1;  /* Other FDs not implemented yet */
}

/* Read: read from file descriptor */
static int32_t sys_read(registers_t *r) {
    int      fd    = (int)r->ebx;
    char    *buf   = (char *)r->ecx;
    uint32_t count = r->edx;
    
    if (fd == 0) {  /* stdin → keyboard (blocking) */
        /* Block until keyboard input available (implemented in Ch 60) */
        /* For now: return 0 (EOF) */
        (void)buf; (void)count;
        return 0;
    }
    
    return -1;
}

/* Get current process ID */
static int32_t sys_getpid(registers_t *r) {
    (void)r;
    return current_process ? (int32_t)current_process->pid : 0;
}

/* Sleep for N milliseconds */
static int32_t sys_sleep(registers_t *r) {
    uint32_t ms = r->ebx;
    sleep_ms(ms);
    return 0;
}

/* Yield CPU voluntarily */
static int32_t sys_yield(registers_t *r) {
    (void)r;
    yield();
    return 0;
}

/* Print a character to screen */
static int32_t sys_putchar(registers_t *r) {
    terminal_putchar((char)r->ebx);
    return 0;
}

/* Get timer ticks */
static int32_t sys_gettime(registers_t *r) {
    (void)r;
    return (int32_t)timer_get_ticks();
}

/* Extend the heap (change program break) */
static int32_t sys_brk(registers_t *r) {
    uint32_t new_brk = r->ebx;
    
    if (!current_process) return -1;
    
    if (new_brk == 0) {
        /* Query current break: */
        return (int32_t)current_process->brk;
    }
    
    /* Extend heap to new_brk by mapping more pages: */
    uint32_t old_brk = current_process->brk;
    if (new_brk <= old_brk) return (int32_t)old_brk;
    
    /* Map pages from old_brk to new_brk: */
    for (uint32_t addr = PAGE_ALIGN_UP(old_brk);
         addr < PAGE_ALIGN_UP(new_brk);
         addr += PAGE_SIZE) {
        vmm_alloc_page(current_process->page_dir, addr, PTE_WRITABLE | PTE_USER);
    }
    
    current_process->brk = new_brk;
    return (int32_t)new_brk;
}
```

---

## 7. System Call Return Values

```
Conventions:
  Success: EAX ≥ 0  (or specific positive value)
  Error:   EAX = -errno (negative error code)
  
Common error codes (we keep it simple):
  -1 = EPERM    (operation not permitted)
  -2 = ENOENT   (no such file or directory)
  -9 = EBADF    (bad file descriptor)
  -14 = EFAULT  (bad address)
  -22 = EINVAL  (invalid argument)
  -12 = ENOMEM  (out of memory)

The assembly stub writes EAX back into the saved register frame,
so when iret executes, the user code resumes with EAX = return value.
```

---

## 8. User-Side System Call Wrappers

For user programs, we provide inline wrappers:

```c
/* user/syscall.h — user-side system call wrappers */
#pragma once
#include "stdint.h"

/* Low-level syscall macro (up to 3 arguments): */
static inline int32_t syscall3(uint32_t num, uint32_t a, uint32_t b, uint32_t c) {
    int32_t ret;
    __asm__ volatile(
        "int $0x80"
        : "=a"(ret)
        : "a"(num), "b"(a), "c"(b), "d"(c)
        : "memory"
    );
    return ret;
}

static inline int32_t syscall1(uint32_t num, uint32_t a) {
    return syscall3(num, a, 0, 0);
}

static inline int32_t syscall0(uint32_t num) {
    return syscall3(num, 0, 0, 0);
}

/* High-level wrappers: */
static inline void exit(int code) {
    syscall1(1, (uint32_t)code);
}

static inline int write(int fd, const void *buf, uint32_t count) {
    return syscall3(2, (uint32_t)fd, (uint32_t)buf, count);
}

static inline int read(int fd, void *buf, uint32_t count) {
    return syscall3(3, (uint32_t)fd, (uint32_t)buf, count);
}

static inline int getpid(void) {
    return syscall0(6);
}

static inline void sleep(uint32_t ms) {
    syscall1(7, ms);
}

static inline void yield(void) {
    syscall0(12);
}

static inline void putchar(char c) {
    syscall1(13, (uint32_t)c);
}

static inline uint32_t gettime(void) {
    return (uint32_t)syscall0(14);
}
```

---

## 9. Testing System Calls

```c
/* kernel/kernel.c — test syscalls from within kernel mode
   (In Ring 0, we can call the handler directly without int $0x80) */

static void syscall_test_process(void) {
    /* Test sys_getpid: */
    int pid = getpid();  /* This calls int $0x80 */
    kprintf("[syscall test] PID = %d\n", pid);
    
    /* Test sys_write: */
    const char *msg = "Hello via syscall write!\n";
    int written = write(1, msg, 25);
    kprintf("[syscall test] write() returned %d\n", written);
    
    /* Test sys_gettime: */
    uint32_t t = gettime();
    kprintf("[syscall test] timer ticks = %u\n", t);
    
    /* Test sys_sleep: */
    kprintf("[syscall test] sleeping 200ms...\n");
    sleep(200);
    kprintf("[syscall test] awake!\n");
    
    /* Test sys_exit: */
    exit(0);
}

/* Register as a process and run: */
process_create("syscall-test", syscall_test_process, 5);
```

Expected output:
```
[syscall test] PID = 1
Hello via syscall write!
[syscall test] write() returned 25
[syscall test] timer ticks = 42
[syscall test] sleeping 200ms...
[syscall test] awake!
[PID 1] exit(0)
```

---

## Summary

| Concept | Description |
|---------|------------|
| System call | Controlled entry point from Ring 3 into Ring 0 |
| int $0x80 | x86 software interrupt for system calls; CPU privilege checks DPL in IDT gate |
| DPL=3 in IDT | Allows Ring 3 code to trigger this interrupt; otherwise GP fault |
| TSS.esp0 | Kernel stack used when int $0x80 causes Ring 3→0 transition |
| EAX | Syscall number on entry; return value on exit |
| EBX/ECX/EDX | Arguments 1/2/3 |
| syscall_table | Array of function pointers indexed by syscall number |
| ENOSYS | Error -38: function not implemented (when syscall_table[num] is NULL) |
| iret | Returns from syscall; restores Ring 3 CS/EIP/EFLAGS/ESP/SS |
| errno convention | Negative EAX = error; -1 = EPERM, -2 = ENOENT, etc. |
| User wrapper | Inline assembly around `int $0x80` for each syscall |
| Argument validation | Kernel must validate pointers from user space before dereferencing them |
