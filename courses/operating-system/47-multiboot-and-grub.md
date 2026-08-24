# Chapter 47: Multiboot and GRUB

> **"Writing a custom bootloader is educational — but GRUB already solves every hard problem: loading from FAT/ext4/ISO, detecting memory, switching CPU modes, passing boot information to your kernel. Why reinvent it? Let GRUB do the boot; you write the kernel. That's the professional approach used by every major open-source OS."**

---

## Table of Contents

1. [The Problem GRUB Solves](#1-the-problem-grub-solves)
2. [The Multiboot Specification](#2-the-multiboot-specification)
3. [Multiboot Header — What Your Kernel Must Provide](#3-multiboot-header--what-your-kernel-must-provide)
4. [What GRUB Does for Us](#4-what-grub-does-for-us)
5. [The Multiboot Information Structure](#5-the-multiboot-information-structure)
6. [Our Kernel Entry Point](#6-our-kernel-entry-point)
7. [The Linker Script](#7-the-linker-script)
8. [grub.cfg](#8-grubcfg)
9. [Building and Booting the ISO](#9-building-and-booting-the-iso)
10. [Summary](#summary)

---

## 1. The Problem GRUB Solves

Writing a bare bootloader (Chapter 46) taught us the fundamentals. But a real bootloader must handle:

```
Hard problems in a bootloader:
  ✗ Reading kernels from ext4, FAT32, ISO 9660 file systems
  ✗ Detecting all available RAM (E820 memory map)
  ✗ Switching from real mode → protected mode
  ✗ Loading a kernel ELF binary at the right address
  ✗ Handling BIOS quirks across hundreds of hardware vendors
  ✗ Supporting boot over network (PXE)
  ✗ Showing a boot menu
  ✗ Loading initial ramdisk (initrd)
```

All of this is already solved by GRUB (GRand Unified Bootloader). Rather than solving it again, we can use a standard — **Multiboot** — that tells GRUB how to load our kernel.

---

## 2. The Multiboot Specification

**Multiboot** is an open standard that defines the contract between a bootloader and a kernel:

```
The deal:
  Kernel (your side):
    - Put a magic header in the first 8KB of the kernel binary
    - Header identifies itself as Multiboot-compliant
    - Header tells GRUB what the kernel needs (flags, load address, etc.)
    
  Bootloader (GRUB's side):
    - Finds the Multiboot header
    - Loads the kernel into memory at the right address
    - Switches CPU to 32-bit protected mode
    - Collects system info (memory map, video mode, boot device)
    - Jumps to kernel entry point with:
        EAX = 0x2BADB002 (Multiboot magic number)
        EBX = pointer to Multiboot Information structure
```

**Two versions:**
- **Multiboot 1**: Original. Simpler. We'll use this.
- **Multiboot 2**: Supports UEFI, 64-bit kernels. More complex.

---

## 3. Multiboot Header — What Your Kernel Must Provide

The kernel binary must contain a Multiboot header in its first 8KB:

```nasm
; boot/multiboot.asm — Multiboot header

section .multiboot
align 4

; Multiboot 1 magic values:
MULTIBOOT_MAGIC    equ 0x1BADB002
MULTIBOOT_FLAGS    equ 0x00000003  ; bit 0: align modules on page boundary
                                    ; bit 1: provide memory map
                                    ; bit 2 (NOT SET): use video mode info
MULTIBOOT_CHECKSUM equ -(MULTIBOOT_MAGIC + MULTIBOOT_FLAGS)

; The header (3 required dwords):
dd MULTIBOOT_MAGIC      ; magic: identifies this as a Multiboot header
dd MULTIBOOT_FLAGS      ; flags: what we need GRUB to do
dd MULTIBOOT_CHECKSUM   ; checksum: magic + flags + checksum must equal 0
```

**Flags breakdown:**
```
Bit 0 (0x01) = ALIGN_MODULES
  Ask GRUB to align modules (initrd) on page (4KB) boundaries
  Good practice — easier for us to work with

Bit 1 (0x02) = MEMORY_INFO
  Ask GRUB to give us the full memory map via the MBI structure
  ESSENTIAL — we need to know which RAM is usable for our allocator

Bit 2 (0x04) = VIDEO_MODE (we don't set this)
  Would ask GRUB to set a graphics mode before booting
  We'll use VGA text mode (already set by BIOS)

Bit 16 (0x10000) = LOAD_ADDRESS (we don't need this)
  If set, header contains load address info
  We don't need it because we use an ELF kernel — GRUB reads addresses from ELF headers
```

**The checksum rule:**
`magic + flags + checksum = 0 (mod 2^32)`

So: `checksum = -(magic + flags) = -(0x1BADB002 + 0x00000003) = -0x1BADB005`

In NASM: `dd -(MULTIBOOT_MAGIC + MULTIBOOT_FLAGS)` — NASM handles the two's complement automatically.

---

## 4. What GRUB Does for Us

Before jumping to our kernel, GRUB has already:

```
1. Loaded the kernel ELF binary from disk (reading our kernel.elf file from the ISO)
   - GRUB understands ELF format, reads program headers, maps sections to correct addresses

2. Switched CPU to 32-bit protected mode with a flat GDT
   - Segment registers all loaded with flat (base=0, limit=4GB) selectors

3. Disabled paging (paging is OFF — we get physical addresses)

4. Loaded a basic GDT (we'll replace this with our own in Chapter 49)

5. Set up EAX = 0x2BADB002 (we check this to confirm Multiboot)
   EBX = physical address of Multiboot Information structure

6. Left interrupts DISABLED (EFLAGS.IF = 0)

7. The FPU state is undefined — we should initialize it if we use FPU/SSE

What GRUB has NOT done:
  ✗ Set up our IDT (we do that — Chapter 50)
  ✗ Enabled paging (we do that — Chapter 53)
  ✗ Set up a proper stack (we do that right after entry)
```

---

## 5. The Multiboot Information Structure

GRUB fills in a structure and passes its address in EBX. Here's the layout:

```c
/* include/multiboot.h */

#define MULTIBOOT_MAGIC 0x2BADB002

/* Flags in mbi->flags — which fields are valid: */
#define MBI_FLAG_MEM        (1 << 0)  /* mem_lower/mem_upper valid */
#define MBI_FLAG_BOOTDEV    (1 << 1)  /* boot_device valid */
#define MBI_FLAG_CMDLINE    (1 << 2)  /* cmdline valid */
#define MBI_FLAG_MODS       (1 << 3)  /* modules list valid */
#define MBI_FLAG_MMAP       (1 << 6)  /* mmap_length/mmap_addr valid */
#define MBI_FLAG_DRIVES     (1 << 7)  /* drives_length/drives_addr valid */

struct multiboot_info {
    uint32_t flags;         /* Flags: which fields below are valid */
    
    /* Memory (if flags & MBI_FLAG_MEM): */
    uint32_t mem_lower;     /* KB of lower memory (below 1MB, typically ~640) */
    uint32_t mem_upper;     /* KB of upper memory (above 1MB, e.g. 130048 = 127MB) */
    
    /* Boot device (if flags & MBI_FLAG_BOOTDEV): */
    uint32_t boot_device;   /* BIOS drive number of boot device */
    
    /* Command line (if flags & MBI_FLAG_CMDLINE): */
    uint32_t cmdline;       /* Physical address of null-terminated command line string */
    
    /* Modules (if flags & MBI_FLAG_MODS): */
    uint32_t mods_count;    /* Number of loaded modules (initrd etc.) */
    uint32_t mods_addr;     /* Physical address of first module entry */
    
    /* ELF section info (used by debuggers): */
    uint32_t syms[4];
    
    /* Memory map (if flags & MBI_FLAG_MMAP): */
    uint32_t mmap_length;   /* Size in bytes of the memory map array */
    uint32_t mmap_addr;     /* Physical address of memory map entries */
    
    /* Drives (if flags & MBI_FLAG_DRIVES): */
    uint32_t drives_length;
    uint32_t drives_addr;
    
    /* Boot loader name: */
    uint32_t boot_loader_name; /* Physical address of name string, e.g. "GRUB 2.12" */
    
    /* ... more fields we don't need right now ... */
} __attribute__((packed));

/* Each memory map entry (pointed to by mmap_addr): */
struct multiboot_mmap_entry {
    uint32_t size;          /* Size of this entry (usually 20) */
    uint64_t base_addr;     /* Physical base address */
    uint64_t length;        /* Length of region in bytes */
    uint32_t type;          /* 1 = usable RAM, 2 = reserved, 3 = ACPI reclaimable, etc. */
} __attribute__((packed));
```

**Reading the memory map:**
```c
void parse_memory_map(struct multiboot_info *mbi) {
    if (!(mbi->flags & MBI_FLAG_MMAP)) {
        /* No memory map — fall back to mem_upper */
        return;
    }
    
    struct multiboot_mmap_entry *entry = (void *)mbi->mmap_addr;
    uint8_t *end = (uint8_t *)mbi->mmap_addr + mbi->mmap_length;
    
    while ((uint8_t *)entry < end) {
        if (entry->type == 1) {
            /* Usable RAM region */
            uint32_t base = (uint32_t)entry->base_addr;
            uint32_t len  = (uint32_t)entry->length;
            /* Mark these pages as free in our physical memory manager */
            pmm_mark_free(base, len);
        }
        /* Move to next entry (size field + 4 bytes for size itself): */
        entry = (void *)((uint8_t *)entry + entry->size + sizeof(uint32_t));
    }
}
```

---

## 6. Our Kernel Entry Point

The kernel's assembly entry point — the first code that runs after GRUB hands off control:

```nasm
; boot/entry.asm — Kernel entry point (called by GRUB)

[BITS 32]                   ; We're in 32-bit protected mode (GRUB set this up)
[GLOBAL kernel_entry]       ; Make this symbol visible to the linker
[EXTERN kernel_main]        ; kernel_main is defined in kernel/kernel.c

section .multiboot
align 4
    dd 0x1BADB002           ; Multiboot magic
    dd 0x00000003           ; Flags: align + memory map
    dd -(0x1BADB002 + 0x00000003)  ; Checksum

section .bss
align 16
stack_bottom:
    resb 16384              ; Reserve 16KB for the initial kernel stack
stack_top:

section .text
kernel_entry:
    ; GRUB hands us control here:
    ;   EAX = 0x2BADB002 (Multiboot magic)
    ;   EBX = pointer to Multiboot Information structure
    ;   Segments are flat (base=0, limit=4GB)
    ;   Paging is OFF
    ;   Interrupts are DISABLED
    
    ; Set up the kernel stack first — we can't call C without a stack:
    mov esp, stack_top
    
    ; Reset the stack frame base pointer:
    xor ebp, ebp
    
    ; Push Multiboot arguments for kernel_main(uint32_t magic, uint32_t mbi_ptr):
    push ebx                ; arg 2: Multiboot Information structure address
    push eax                ; arg 1: Multiboot magic (0x2BADB002)
    
    ; Call the C kernel:
    call kernel_main
    
    ; If kernel_main ever returns (it shouldn't), halt:
.hang:
    cli                     ; Disable interrupts
    hlt                     ; Halt the CPU
    jmp .hang               ; Loop in case of NMI
```

---

## 7. The Linker Script

The linker script tells the linker where to place each section in memory:

```ld
/* linker.ld */

ENTRY(kernel_entry)         /* Entry point symbol name */

SECTIONS {
    /* Kernel loads at 1MB (0x100000) — standard for Multiboot kernels */
    . = 1M;
    
    /* .multiboot must be first (within first 8KB of the binary): */
    .multiboot : {
        *(.multiboot)
    }
    
    /* Executable code: */
    .text ALIGN(4K) : {
        *(.text)
    }
    
    /* Read-only data (string literals, const globals): */
    .rodata ALIGN(4K) : {
        *(.rodata)
    }
    
    /* Read-write initialized data: */
    .data ALIGN(4K) : {
        *(.data)
    }
    
    /* Uninitialized data (zero-initialized at boot): */
    .bss ALIGN(4K) : {
        *(COMMON)
        *(.bss)
        
        /* Stack is also in BSS (reserved but not initialized): */
    }
    
    /* Symbol we can use from C to find where kernel ends: */
    kernel_end = .;
}
```

**Why 1MB (0x100000)?**
```
Memory layout after GRUB loads us:
  0x00000 - 0x00FFF: IVT + BDA (BIOS uses this)
  0x01000 - 0x7FFFF: GRUB stack, modules, etc.
  0x80000 - 0x9FFFF: EBDA (Extended BIOS Data Area)
  0xA0000 - 0xFFFFF: VGA + BIOS ROM (DO NOT WRITE HERE)
  0x100000+: our kernel starts here (above the "640KB memory hole")
  
The region 0x0009FC00 - 0x000FFFFF is the "upper memory area" —
mostly reserved. Starting at exactly 1MB is the conventional choice.
```

---

## 8. grub.cfg

GRUB reads this configuration file from the ISO to know what to boot:

```bash
# boot/grub/grub.cfg

# Don't wait at the menu (0 second timeout):
set timeout=0
set default=0

# Boot menu entry:
menuentry "TinyOS" {
    # Load our kernel as a Multiboot module:
    multiboot /boot/kernel.elf
    
    # (Optional) pass a command line to the kernel:
    # multiboot /boot/kernel.elf root=/dev/sda1 debug
    
    # (Optional) load an initial ramdisk:
    # module /boot/initrd.img
    
    boot
}
```

The path `/boot/kernel.elf` is relative to the ISO root. When GRUB boots from the ISO, it reads this file, finds the Multiboot header, and loads the ELF segments to their correct physical addresses.

---

## 9. Building and Booting the ISO

Putting it all together:

**Project structure for this chapter:**
```
tinyos/
├── boot/
│   ├── entry.asm           ← kernel entry + multiboot header
│   └── grub/
│       └── grub.cfg        ← GRUB configuration
├── kernel/
│   └── kernel.c            ← kernel_main
├── include/
│   ├── multiboot.h         ← MBI structure definitions
│   └── stdint.h            ← uint32_t etc.
├── linker.ld               ← linker script
└── Makefile
```

**kernel/kernel.c:**
```c
#include "multiboot.h"
#include "stdint.h"

/* VGA text mode buffer at physical 0xB8000: */
static volatile uint16_t *vga = (uint16_t *)0xB8000;

static void vga_put(int x, int y, char c, uint8_t color) {
    vga[y * 80 + x] = (uint16_t)c | ((uint16_t)color << 8);
}

static void print(const char *s, int row) {
    for (int i = 0; s[i]; i++) {
        vga_put(i, row, s[i], 0x0F); /* white on black */
    }
}

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    /* Clear screen: */
    for (int i = 0; i < 80 * 25; i++) {
        vga[i] = 0x0F20; /* space, white on black */
    }
    
    /* Verify Multiboot magic: */
    if (magic != 0x2BADB002) {
        print("ERROR: Not booted by Multiboot-compliant bootloader!", 0);
        while (1) {}
    }
    
    print("TinyOS is alive!", 0);
    
    /* Parse the Multiboot Info structure: */
    struct multiboot_info *mbi = (struct multiboot_info *)mbi_ptr;
    
    if (mbi->flags & MBI_FLAG_MEM) {
        /* mem_upper is KB of RAM above 1MB */
        /* mem_lower is KB of RAM below 1MB (usually 640) */
        print("Memory detected via Multiboot!", 1);
    }
    
    if (mbi->flags & MBI_FLAG_MMAP) {
        print("Memory map available — ready for PMM!", 2);
    }
    
    /* Halt — more functionality in upcoming chapters: */
    while (1) {}
}
```

**Makefile:**
```makefile
CC      = i686-elf-gcc
LD      = i686-elf-ld
AS      = nasm
CFLAGS  = -std=c99 -ffreestanding -O2 -Wall -Iinclude -nostdlib -nostdinc \
          -fno-builtin -fno-stack-protector -m32
LDFLAGS = -T linker.ld -nostdlib

OBJS = boot/entry.o kernel/kernel.o

all: os.iso

boot/entry.o: boot/entry.asm
	$(AS) -f elf32 $< -o $@

kernel/kernel.o: kernel/kernel.c
	$(CC) $(CFLAGS) -c $< -o $@

kernel.elf: $(OBJS)
	$(LD) $(LDFLAGS) -o $@ $^

os.iso: kernel.elf
	mkdir -p isodir/boot/grub
	cp kernel.elf isodir/boot/
	cp boot/grub/grub.cfg isodir/boot/grub/
	grub-mkrescue -o os.iso isodir

run: os.iso
	qemu-system-i386 -m 128M -cdrom os.iso -serial stdio

debug: kernel.elf os.iso
	qemu-system-i386 -m 128M -cdrom os.iso -s -S &
	gdb kernel.elf \
	    -ex "target remote localhost:1234" \
	    -ex "set architecture i386" \
	    -ex "break kernel_main" \
	    -ex "continue"

clean:
	rm -f boot/entry.o kernel/kernel.o kernel.elf os.iso
	rm -rf isodir/

.PHONY: all run debug clean
```

**Build and run:**
```bash
make
make run
```

You should see a black QEMU window with white text:
```
TinyOS is alive!
Memory detected via Multiboot!
Memory map available — ready for PMM!
```

If you see those three lines, GRUB successfully:
- Loaded your kernel ELF
- Passed CPU control to `kernel_entry`
- Your kernel verified the Multiboot magic
- Your kernel read the MBI structure

You now have a working foundation for everything in Volume 9.

---

## Summary

| Concept | Description |
|---------|------------|
| Multiboot | Standard contract between bootloader and kernel |
| Multiboot header | 3 dwords in first 8KB of kernel: magic (0x1BADB002) + flags + checksum |
| GRUB | GNU bootloader; reads Multiboot header, loads ELF kernel, switches to protected mode |
| EAX at entry | 0x2BADB002 — Multiboot magic (verifies GRUB booted us) |
| EBX at entry | Physical address of Multiboot Information (MBI) structure |
| MBI | Struct from GRUB: memory map, command line, boot device, module list |
| MBI flags | Bitfield: which MBI fields are valid (check before reading) |
| Memory map entry | 20-byte struct: base address, length, type (type 1 = usable RAM) |
| linker.ld | Places kernel at 0x100000; .multiboot section first |
| grub.cfg | `multiboot /boot/kernel.elf` — tells GRUB what to load |
| grub-mkrescue | Creates bootable ISO with GRUB + our kernel |
| kernel.elf | ELF binary — GRUB reads ELF headers to load sections at correct addresses |
| Stack | We allocate 16KB in .bss, set ESP = stack_top before calling kernel_main |
