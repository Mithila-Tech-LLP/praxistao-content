# Chapter 59: User Mode

> **"User mode is the last wall between your kernel and the chaos of untrusted code. When you jump to Ring 3, everything changes: the process can no longer access kernel memory, execute privileged instructions, or modify page tables. It is isolated, sandboxed, and can only communicate with the kernel through the system call gate. This is the boundary that makes operating systems secure."**

---

## Table of Contents

1. [Ring 3 vs Ring 0 — The Difference](#1-ring-3-vs-ring-0--the-difference)
2. [Setting Up a User Address Space](#2-setting-up-a-user-address-space)
3. [Loading a User Program](#3-loading-a-user-program)
4. [The Jump to Ring 3](#4-the-jump-to-ring-3)
5. [A Simple User Program](#5-a-simple-user-program)
6. [ELF Basics (Optional)](#6-elf-basics-optional)
7. [User Stack Setup](#7-user-stack-setup)
8. [Complete usermode.c / usermode.h](#8-complete-usermodec--usermodeh)
9. [Testing User Mode](#9-testing-user-mode)
10. [Summary](#summary)

---

## 1. Ring 3 vs Ring 0 — The Difference

```
Ring 0 (kernel):
  Can execute:   cli/sti, lgdt/lidt, mov cr0-cr4, in/out, hlt, wrmsr, etc.
  Can access:    any memory address (physical or virtual)
  Can load:      any segment selector including DPL=0 segments
  
Ring 3 (user):
  Cannot execute: ANY privileged instruction → GP fault (#GP)
  Cannot access:  pages with User bit (PTE_USER) = 0 → page fault (#PF)
  Cannot load:    DPL=0 segment selectors → GP fault
  CAN do:         normal math, logic, memory reads/writes (to mapped pages)
                  system calls via int $0x80 (DPL=3 in IDT)

What this means practically:
  A buggy user program that writes to address 0x100000 (kernel code!):
    → CPU checks page table: PTE_USER not set
    → Page fault exception
    → Kernel handles it: kill the process
    → Kernel continues unharmed
    
  Without user mode: the bug would corrupt the kernel directly!
```

---

## 2. Setting Up a User Address Space

User processes need their own isolated virtual address space:

```
User virtual address space layout:
  0x00000000 - 0x00000FFF: NULL page (unmapped — catches null pointer dereferences)
  0x00001000 - 0x00FFFFFF: User code + data (loaded here)
  0x01000000 - 0x0FFFFFFF: User heap (grows up via sys_brk)
  0x10000000 - 0xBFFFFFFF: Available for user stack, mmap, etc.
  0xC0000000 - 0xFFFFFFFF: Kernel space (mapped PTE_KERNEL, not user-accessible)
                            Same in ALL address spaces (kernel pages shared)

Kernel virtual addresses (0xC0000000+) are mapped in every user process's
page directory with PTE_USER=0, so user code can't access them.
The same physical RAM is mapped there; no copy needed.
```

```c
/* Create a fresh user address space: */
uint32_t *create_user_address_space(void) {
    uint32_t *pd = vmm_create_address_space();  /* copies kernel mappings */
    if (!pd) return NULL;
    
    /* NULL page is unmapped (PD entry 0 = 0) — catches null dereferences */
    /* User code starts at 0x1000 — we'll map it when loading the program */
    
    return pd;
}
```

---

## 3. Loading a User Program

For simplicity, we load a raw binary (not ELF) into the user address space.
The binary is embedded in the kernel as a byte array (compiled in):

```c
/* Load a raw binary into user address space at virtual address 'load_addr':
   The binary is treated as a flat executable — no headers.
   Entry point = load_addr. */

int load_user_binary(uint32_t *page_dir, const uint8_t *binary,
                     uint32_t binary_size, uint32_t load_addr) {
    /* How many pages do we need? */
    uint32_t pages = (binary_size + PAGE_SIZE - 1) / PAGE_SIZE;
    
    /* Map and copy each page: */
    for (uint32_t i = 0; i < pages; i++) {
        uint32_t vaddr = load_addr + i * PAGE_SIZE;
        
        /* Allocate a physical frame: */
        uint32_t phys = pmm_alloc_frame();
        if (!phys) return -1;
        
        /* Zero the frame: */
        memset((void *)phys, 0, PAGE_SIZE);
        
        /* Copy binary data into it: */
        uint32_t offset = i * PAGE_SIZE;
        uint32_t copy_size = binary_size - offset;
        if (copy_size > PAGE_SIZE) copy_size = PAGE_SIZE;
        memcpy((void *)phys, binary + offset, copy_size);
        
        /* Map into user page directory with USER flag: */
        vmm_map_page(page_dir, vaddr, phys, PTE_WRITABLE | PTE_USER);
    }
    
    return 0;
}
```

---

## 4. The Jump to Ring 3

The critical step: switching from Ring 0 to Ring 3 for the first time.
We do this using `iret` with carefully crafted values on the stack:

```
iret pops from stack (in this order):
  EIP     ← where to start in Ring 3
  CS      ← Ring 3 code selector (must have RPL=3)
  EFLAGS  ← must have IF=1 (interrupts enabled)
  ESP     ← user stack pointer
  SS      ← Ring 3 stack selector (must have RPL=3)

Since we're switching to a LOWER privilege (Ring 0 → Ring 3),
iret needs the full 5-item form (with ESP and SS).
```

```nasm
; boot/usermode.asm — Jump to Ring 3

[BITS 32]
[GLOBAL jump_to_usermode]

; void jump_to_usermode(uint32_t user_eip, uint32_t user_esp)
; Switches to Ring 3 and jumps to user_eip with user_esp as stack.
; This function DOES NOT RETURN (to the caller).

jump_to_usermode:
    ; Arguments on stack:
    ; [esp+4] = user_eip
    ; [esp+8] = user_esp
    
    mov eax, [esp + 4]   ; user EIP
    mov ecx, [esp + 8]   ; user ESP
    
    ; Load user data segment selector into DS, ES, FS, GS:
    ; GDT_USER_DATA | 3 = 0x20 | 3 = 0x23 (RPL=3)
    mov dx, 0x23
    mov ds, dx
    mov es, dx
    mov fs, dx
    mov gs, dx
    
    ; Build iret frame on the CURRENT (kernel) stack:
    push 0x23          ; SS  — user stack segment (GDT_USER_DATA | RPL=3)
    push ecx           ; ESP — user stack pointer
    pushf              ; EFLAGS — push current EFLAGS
    
    ; Set IF (interrupt enable) in EFLAGS on the stack:
    ; We MUST have interrupts enabled in user mode, otherwise no timer/keyboard
    pop edx
    or edx, 0x200      ; set bit 9 = IF
    push edx
    
    push 0x1B          ; CS  — user code segment (GDT_USER_CODE | RPL=3 = 0x18|3 = 0x1B)
    push eax           ; EIP — user entry point
    
    ; Jump to Ring 3:
    iret
    ; CPU pops: EIP, CS(Ring 3), EFLAGS, ESP, SS(Ring 3)
    ; We are now in Ring 3!
```

---

## 5. A Simple User Program

A flat binary that makes system calls. We compile it separately and embed it:

```c
/* user/hello.c — a minimal user program */
/* NOTE: compiled with -ffreestanding -nostdlib; no C library */

/* System call numbers: */
#define SYS_WRITE   2
#define SYS_EXIT    1
#define SYS_PUTCHAR 13

/* Inline system call: */
static int syscall(int num, int a, int b, int c) {
    int ret;
    __asm__ volatile(
        "int $0x80"
        : "=a"(ret)
        : "a"(num), "b"(a), "c"(b), "d"(c)
    );
    return ret;
}

static void puts_user(const char *s) {
    int len = 0;
    while (s[len]) len++;
    syscall(SYS_WRITE, 1, (int)s, len);
}

/* Entry point (must be at the very start of the binary): */
void _start(void) {
    puts_user("Hello from user mode!\n");
    puts_user("I am in Ring 3. The kernel cannot be touched.\n");
    
    /* Try writing to kernel address — will page fault: */
    /* volatile int *p = (int *)0xC0000000; *p = 1; */  /* DON'T: would page fault */
    
    /* Exit: */
    syscall(SYS_EXIT, 0, 0, 0);
    
    /* Should never reach here: */
    for (;;) {}
}
```

**Compiling the user program:**
```bash
# Compile to raw binary (no ELF headers, entry at byte 0):
i686-elf-gcc -ffreestanding -nostdlib -nostdinc -O2 -m32 \
    -T user_linker.ld -o hello.elf user/hello.c

# Extract raw binary:
i686-elf-objcopy -O binary hello.elf hello.bin

# Convert to C array for embedding:
xxd -i hello.bin > kernel/user_hello.h
```

**user_linker.ld** (flat binary, entry at 0x1000):
```ld
ENTRY(_start)
SECTIONS {
    . = 0x1000;
    .text : { *(.text) }
    .data : { *(.data) *(.rodata) }
    .bss  : { *(.bss) }
}
```

---

## 6. ELF Basics (Optional)

For a production OS, we'd load ELF binaries:

```
ELF header (52 bytes):
  Magic:        0x7F 'E' 'L' 'F'
  Class:        1 = 32-bit
  Data:         1 = little-endian
  e_type:       2 = ET_EXEC (executable)
  e_machine:    3 = EM_386 (x86)
  e_entry:      virtual address of entry point (e.g., 0x1000)
  e_phoff:      offset to program headers
  e_phnum:      number of program headers

Program header (segment):
  p_type:   1 = PT_LOAD (must be loaded into memory)
  p_offset: offset in ELF file
  p_vaddr:  virtual address to load at
  p_filesz: size in file
  p_memsz:  size in memory (>filesz for BSS: zero-fill the difference)
  p_flags:  R=4, W=2, X=1 (permissions)

To load ELF:
  1. Read ELF header → find e_entry (entry point), e_phoff, e_phnum
  2. For each PT_LOAD segment:
     a. Allocate pages from p_vaddr to p_vaddr + p_memsz
     b. Copy p_filesz bytes from file at p_offset
     c. Zero the remaining p_memsz - p_filesz bytes
  3. Jump to e_entry
```

---

## 7. User Stack Setup

```c
/* Set up a user stack for the process: */
uint32_t setup_user_stack(uint32_t *page_dir) {
    /* Allocate 4 pages (16KB) for the user stack: */
    #define USER_STACK_TOP   0xBFFFF000   /* Near top of user space */
    #define USER_STACK_PAGES 4
    
    for (int i = 0; i < USER_STACK_PAGES; i++) {
        uint32_t vaddr = USER_STACK_TOP - (i + 1) * PAGE_SIZE;
        uint32_t phys  = pmm_alloc_frame();
        if (!phys) return 0;
        memset((void *)phys, 0, PAGE_SIZE);
        vmm_map_page(page_dir, vaddr, phys, PTE_WRITABLE | PTE_USER);
    }
    
    /* Return top of stack (stack grows downward): */
    return USER_STACK_TOP;
}
```

---

## 8. Complete usermode.c / usermode.h

```c
/* kernel/usermode.c */

#include "usermode.h"
#include "process.h"
#include "vmm.h"
#include "pmm.h"
#include "heap.h"
#include "gdt.h"
#include "vga.h"
#include "string.h"

extern void jump_to_usermode(uint32_t eip, uint32_t esp);

/* Launch a user process from a flat binary: */
process_t *launch_user_process(const char *name,
                                const uint8_t *binary,
                                uint32_t binary_size) {
    /* 1. Create PCB: */
    process_t *proc = pcb_alloc();
    if (!proc) return NULL;
    
    proc->pid   = next_pid++;
    proc->ppid  = current_process ? current_process->pid : 0;
    proc->state = PROC_READY;
    proc->priority = 5;
    
    for (int i = 0; name[i] && i < PROC_NAME_MAX-1; i++) proc->name[i] = name[i];
    
    /* 2. Create user address space: */
    proc->page_dir = create_user_address_space();
    if (!proc->page_dir) { proc->state = PROC_DEAD; return NULL; }
    
    /* 3. Load binary at 0x1000: */
    if (load_user_binary(proc->page_dir, binary, binary_size, 0x1000) < 0) {
        proc->state = PROC_DEAD;
        return NULL;
    }
    
    /* 4. Set up user stack: */
    uint32_t user_esp = setup_user_stack(proc->page_dir);
    if (!user_esp) { proc->state = PROC_DEAD; return NULL; }
    proc->user_stack = user_esp;
    
    /* 5. Set up kernel stack for the initial context switch.
       
       When the scheduler first switches to this process, it will call
       switch_to() → switch_context() which pops EDI/ESI/EBX/EBP/EIP.
       
       We set up a "fake" kernel stack frame that will:
       - restore registers (all 0)
       - "return" to a trampoline function that calls jump_to_usermode() */
    proc->kernel_stack = (uint32_t)kmalloc(KERNEL_STACK_SIZE);
    
    uint32_t *kstack = (uint32_t *)(proc->kernel_stack + KERNEL_STACK_SIZE);
    
    /* Arguments for jump_to_usermode: */
    *(--kstack) = user_esp;          /* arg 2: user ESP */
    *(--kstack) = 0x1000;            /* arg 1: user EIP (entry point) */
    *(--kstack) = 0;                 /* fake return address (never used) */
    
    /* Callee-saved regs for switch_context (will be restored, then ret to trampoline): */
    *(--kstack) = (uint32_t)jump_to_usermode;  /* EIP = trampoline */
    *(--kstack) = 0;    /* EBP */
    *(--kstack) = 0;    /* EBX */
    *(--kstack) = 0;    /* ESI */
    *(--kstack) = 0;    /* EDI */
    
    proc->kernel_esp = (uint32_t)kstack;
    
    /* 6. Add to run queue: */
    proc->next = run_queue;
    if (run_queue) run_queue->prev = proc;
    run_queue  = proc;
    
    kprintf("Launched user process '%s' (PID=%u, entry=0x1000)\n",
            proc->name, proc->pid);
    return proc;
}
```

---

## 9. Testing User Mode

```c
/* Include the pre-compiled user program: */
#include "user_hello.h"   /* Generated: extern uint8_t hello_bin[]; extern uint32_t hello_bin_len; */

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    /* ... all init ... */
    
    /* Launch user-mode process: */
    launch_user_process("hello", hello_bin, hello_bin_len);
    
    __asm__ volatile("sti");
    schedule();  /* Start scheduling — user process will run */
    
    for (;;) __asm__ volatile("hlt");
}
```

Expected output:
```
Launched user process 'hello' (PID=1, entry=0x1000)
Starting scheduler...
Hello from user mode!
I am in Ring 3. The kernel cannot be touched.
[PID 1] exit(0)
```

If you see "Hello from user mode!" — congratulations! Your kernel has:
- A working address space isolation (page directory switch)
- Ring 3 user mode
- System calls across privilege levels
- Protected kernel memory

---

## Summary

| Concept | Description |
|---------|------------|
| Ring 3 | User mode: cannot execute privileged instructions; limited memory access |
| PTE_USER | Page table flag bit 2: must be set for user code to access the page |
| User code selector | 0x1B = GDT index 3, RPL=3; loaded as CS when jumping to Ring 3 |
| User data selector | 0x23 = GDT index 4, RPL=3; loaded as DS/ES/FS/GS/SS in user mode |
| iret for Ring switch | To go Ring 0→3: push SS, ESP, EFLAGS, CS, EIP then iret |
| EFLAGS.IF | Must be set (=1) in the pushed EFLAGS before iret to enable interrupts in Ring 3 |
| jump_to_usermode | Assembly function that builds the iret frame and jumps to Ring 3 |
| Flat binary | Simple executable: raw machine code, entry point at byte 0 or fixed address |
| User address space | Separate page directory; kernel mappings shared (DPL=0); user pages DPL=3 |
| NULL page | Virtual page 0 unmapped: crashes process on null pointer dereference (good!) |
| 0x1000 | User code entry point (above NULL page) |
| USER_STACK_TOP | Near 0xC0000000; grows downward into user address space |
| sys_brk | Syscall to grow user heap dynamically |
