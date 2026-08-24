# Chapter 44: Development Environment Setup

> **"Before you can write a single byte of your OS, you need to build the forge. A cross-compiler, a virtual machine, an assembler, a linker, a debugger — each one is a tool that exists specifically because OS development is different from application development. Getting this right once means smooth sailing for all 19 chapters that follow."**

---

## Table of Contents

1. [What We're Building](#1-what-were-building)
2. [Tools Overview](#2-tools-overview)
3. [Installing Prerequisites](#3-installing-prerequisites)
4. [Building the Cross-Compiler (GCC for i686-elf)](#4-building-the-cross-compiler-gcc-for-i686-elf)
5. [NASM — The Assembler](#5-nasm--the-assembler)
6. [QEMU — The Emulator](#6-qemu--the-emulator)
7. [GDB — The Debugger](#7-gdb--the-debugger)
8. [Project Structure and Makefile](#8-project-structure-and-makefile)
9. [Testing the Toolchain](#9-testing-the-toolchain)
10. [Summary](#summary)

---

## 1. What We're Building

Over Volumes 8 and 9, we will build a working x86 OS from scratch called **TinyOS**:

```
TinyOS features:
  - 32-bit protected mode (x86)
  - Bootloader (loads our kernel from disk)
  - GDT (segments for code/data)
  - IDT (exception and interrupt handlers)
  - Keyboard driver (scan codes → ASCII)
  - Physical memory manager (bitmap allocator)
  - Virtual memory (paging with page directory/tables)
  - Heap allocator (kmalloc/kfree)
  - Process model (PCB, context switching)
  - Round-robin scheduler
  - System calls (int 0x80)
  - User mode (Ring 3)
  - VFS + RAM-backed file system
  - Simple shell
```

**Target architecture:** x86 (32-bit protected mode, later 64-bit in Ch 64)
**Host environment:** Any Linux/macOS machine

---

## 2. Tools Overview

| Tool | Purpose |
|------|---------|
| NASM | Assembler — converts assembly to machine code |
| GCC (cross-compiler) | C compiler targeting i686-elf (not the host OS!) |
| ld (binutils) | Linker — combines object files |
| QEMU | x86 virtual machine — test OS without real hardware |
| GDB + QEMU | Debugger — step through kernel code |
| GNU Make | Build system — automates compilation |
| xorriso | Create bootable ISO images |
| grub-mkrescue | Create GRUB bootable ISOs |

**Why a CROSS-COMPILER?**
Your laptop's GCC targets your OS (Linux/macOS) and links against your OS's C library (glibc). We want GCC that targets a bare-metal i686 system with NO operating system.

Using your host's GCC is a common mistake — it will link against glibc by default, which requires a running OS to work.

---

## 3. Installing Prerequisites

**On Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y \
    build-essential \    # gcc, make, binutils
    nasm \               # assembler
    qemu-system-x86 \   # x86 emulator
    gdb \               # debugger
    xorriso \           # ISO creation
    grub-pc-bin \       # GRUB bootloader
    grub-common \       # GRUB utilities
    git \
    curl \
    bison \             # parser generator (needed for GCC build)
    flex \              # lexer generator (needed for GCC build)
    libgmp3-dev \       # math library (GCC dependency)
    libmpc-dev \        # math library (GCC dependency)
    libmpfr-dev \       # math library (GCC dependency)
    texinfo \           # documentation tools (GCC dependency)
    libisl-dev          # optimization library (GCC dependency)
```

**On macOS (with Homebrew):**
```bash
brew install nasm qemu gdb xorriso
# Cross-compiler via osxcross or direct build
# See cross-compiler section below
```

---

## 4. Building the Cross-Compiler (GCC for i686-elf)

We need GCC and binutils configured to produce code for `i686-elf` (bare metal x86-32).

```bash
# Set up environment:
export TARGET=i686-elf
export PREFIX="$HOME/opt/cross"
export PATH="$PREFIX/bin:$PATH"

mkdir -p $PREFIX

# Download sources:
cd /tmp
wget https://ftp.gnu.org/gnu/binutils/binutils-2.41.tar.gz
wget https://ftp.gnu.org/gnu/gcc/gcc-13.2.0/gcc-13.2.0.tar.gz

tar xf binutils-2.41.tar.gz
tar xf gcc-13.2.0.tar.gz

# Build binutils (as, ld, objcopy, nm for our target):
mkdir build-binutils && cd build-binutils
../binutils-2.41/configure \
    --target=$TARGET \
    --prefix=$PREFIX \
    --with-sysroot \
    --disable-nls \
    --disable-werror

make -j$(nproc)
make install
cd ..

# Build GCC (C compiler targeting i686-elf):
mkdir build-gcc && cd build-gcc
../gcc-13.2.0/configure \
    --target=$TARGET \
    --prefix=$PREFIX \
    --disable-nls \
    --enable-languages=c \
    --without-headers     # no C library headers (we're bare metal!)

make -j$(nproc) all-gcc
make -j$(nproc) all-target-libgcc
make install-gcc
make install-target-libgcc
cd ..

# Verify:
i686-elf-gcc --version
# i686-elf-gcc (GCC) 13.2.0

i686-elf-ld --version
# GNU ld (GNU Binutils) 2.41
```

**Add to your PATH permanently:**
```bash
echo 'export PATH="$HOME/opt/cross/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

---

## 5. NASM — The Assembler

**NASM (Netwide Assembler)** is the industry-standard x86 assembler for OS development.

```bash
# Verify NASM:
nasm --version
# NASM version 2.16.01

# Assemble a quick test:
cat > /tmp/test.asm << 'EOF'
bits 16
org 0x7C00

; Print 'A' using BIOS
mov ah, 0x0E
mov al, 'A'
int 0x10

; Halt
cli
hlt

times 510 - ($ - $$) db 0
dw 0xAA55
EOF

nasm -f bin /tmp/test.asm -o /tmp/test.bin

# Run in QEMU to see 'A' printed:
qemu-system-i386 -drive format=raw,file=/tmp/test.bin

# Should show a black screen with 'A' in the top-left
```

**NASM syntax reminders:**
```nasm
; Comments
mov ax, 5          ; move immediate 5 into ax
mov ax, [0x7E00]   ; move value AT address 0x7E00 into ax
mov [0x7E00], ax   ; move ax INTO address 0x7E00
db 0x41            ; define byte 0x41
dw 0xAA55          ; define word (16-bit)
dd 0xDEADBEEF      ; define dword (32-bit)
dq 0x123456789ABCDEF0  ; define qword (64-bit)
times 10 db 0      ; repeat 'db 0' 10 times
$ - $$             ; current address - start of section = bytes written so far
```

---

## 6. QEMU — The Emulator

**QEMU** lets you run your OS without rebooting a physical machine. Essential for fast development.

```bash
# Run a raw disk image (bootloader only):
qemu-system-i386 -drive format=raw,file=boot.bin

# Run with 128MB RAM, no display, serial output:
qemu-system-i386 \
    -m 128M \
    -drive format=raw,file=os.img \
    -serial stdio \      # redirect serial port to terminal
    -no-reboot \         # don't reboot on triple fault (helps debugging)
    -no-shutdown         # don't exit on shutdown (keeps QEMU open)

# Run with GDB debugging enabled:
qemu-system-i386 \
    -m 128M \
    -drive format=raw,file=os.img \
    -s \                 # -s = shorthand for -gdb tcp::1234 (wait for GDB on port 1234)
    -S \                 # -S = freeze CPU at start (wait for GDB to connect before running)
    -no-reboot

# Run from ISO (with GRUB):
qemu-system-i386 -cdrom os.iso -m 128M

# 64-bit VM:
qemu-system-x86_64 -m 256M -drive format=raw,file=os64.img
```

**QEMU monitor (interactive control):**
```bash
# Press Ctrl+Alt+2 to switch to QEMU monitor
# Press Ctrl+Alt+1 to switch back to VM display

# In monitor:
info registers    # show all CPU registers
x /10b 0x7C00    # hex dump 10 bytes at address 0x7C00
xp /10b 0x100000 # dump physical memory
stop             # pause VM
cont             # resume VM
quit             # quit QEMU
```

---

## 7. GDB — The Debugger

Combined QEMU + GDB allows **source-level debugging of kernel code**:

```bash
# Terminal 1: Start QEMU paused, waiting for GDB:
qemu-system-i386 -m 128M -drive format=raw,file=os.img -s -S

# Terminal 2: Connect GDB:
gdb

# In GDB:
(gdb) target remote localhost:1234   # connect to QEMU
(gdb) set architecture i8086          # for 16-bit boot code
(gdb) break *0x7C00                  # breakpoint at boot sector start
(gdb) continue                        # run until breakpoint
(gdb) info registers                  # show registers
(gdb) x/10i $pc                      # disassemble 10 instructions at PC
(gdb) stepi                          # execute one instruction
(gdb) set architecture i386          # switch to 32-bit after entering protected mode
(gdb) symbol-file kernel.elf         # load kernel symbols for source-level debug
(gdb) break kernel_main              # break on C function
(gdb) continue
(gdb) list                           # show source code at breakpoint
(gdb) print variable_name           # print variable value
(gdb) backtrace                     # show call stack
```

---

## 8. Project Structure and Makefile

Our TinyOS project structure:

```
tinyos/
├── Makefile
├── boot/
│   ├── boot.asm        # 512-byte bootloader
│   └── loader.asm      # second-stage loader (loads kernel)
├── kernel/
│   ├── kernel.c        # kernel entry point
│   ├── gdt.c / gdt.h   # Global Descriptor Table
│   ├── idt.c / idt.h   # Interrupt Descriptor Table
│   ├── isr.c / isr.h   # Interrupt Service Routines
│   ├── pic.c / pic.h   # PIC (8259) driver
│   ├── timer.c         # PIT timer
│   ├── keyboard.c      # keyboard driver
│   ├── pmm.c / pmm.h   # Physical Memory Manager
│   ├── vmm.c / vmm.h   # Virtual Memory Manager
│   ├── heap.c / heap.h # Heap allocator
│   ├── process.c       # Process management
│   ├── scheduler.c     # Round-robin scheduler
│   ├── syscall.c       # System call handler
│   ├── vfs.c           # Virtual File System
│   └── shell.c         # Simple shell
├── include/
│   ├── stdint.h        # uint8_t, uint32_t etc.
│   ├── string.h        # memcpy, memset, strlen
│   └── io.h            # inb, outb inline functions
├── linker.ld           # Linker script
└── grub.cfg            # GRUB configuration
```

**Makefile:**
```makefile
# TinyOS Makefile

# Cross-compiler:
CC      = i686-elf-gcc
LD      = i686-elf-ld
AS      = nasm

# Flags:
CFLAGS  = -std=c99 -ffreestanding -O2 -Wall -Wextra \
          -Iinclude -nostdlib -nostdinc -fno-builtin \
          -fno-stack-protector -fno-pic -m32

LDFLAGS = -T linker.ld -nostdlib

# Source files:
C_SRCS  = $(wildcard kernel/*.c)
ASM_SRCS = boot/boot.asm
C_OBJS  = $(C_SRCS:.c=.o)
ASM_OBJS = $(ASM_SRCS:.asm=.o)

# Default target:
all: os.iso

# Compile C files:
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# Assemble ASM files:
%.o: %.asm
	$(AS) -f elf32 $< -o $@

# Link kernel:
kernel.elf: $(C_OBJS) $(ASM_OBJS)
	$(LD) $(LDFLAGS) -o $@ $^

# Create bootable ISO with GRUB:
os.iso: kernel.elf
	mkdir -p isodir/boot/grub
	cp kernel.elf isodir/boot/
	cp grub.cfg isodir/boot/grub/
	grub-mkrescue -o os.iso isodir

# Run in QEMU:
run: os.iso
	qemu-system-i386 -m 128M -cdrom os.iso

# Debug with GDB:
debug: os.iso
	qemu-system-i386 -m 128M -cdrom os.iso -s -S &
	gdb kernel.elf -ex "target remote localhost:1234"

# Clean:
clean:
	rm -f kernel/*.o boot/*.o kernel.elf os.iso
	rm -rf isodir/

.PHONY: all run debug clean
```

**grub.cfg:**
```
set timeout=0
set default=0

menuentry "TinyOS" {
    multiboot /boot/kernel.elf
    boot
}
```

**linker.ld:**
```ld
ENTRY(kernel_main)

SECTIONS {
    . = 0x100000;   /* Load kernel at 1MB */
    
    .text : {
        *(.multiboot)   /* Multiboot header must be first */
        *(.text)
    }
    
    .rodata : {
        *(.rodata)
    }
    
    .data : {
        *(.data)
    }
    
    .bss : {
        *(COMMON)
        *(.bss)
    }
    
    kernel_end = .;  /* Symbol marking end of kernel */
}
```

---

## 9. Testing the Toolchain

A quick smoke test to verify everything works:

```c
// kernel/kernel.c — Hello World from our kernel
void kernel_main(void) {
    // VGA text mode: write directly to video memory at 0xB8000
    char *vga = (char*)0xB8000;
    
    const char *msg = "Hello from TinyOS!";
    for (int i = 0; msg[i]; i++) {
        vga[i * 2]     = msg[i];   // character
        vga[i * 2 + 1] = 0x0F;    // attribute: white on black
    }
    
    // Halt forever:
    while (1) {}
}
```

```bash
# Build and run:
make
make run
# QEMU window should show "Hello from TinyOS!"
```

If you see the message — your entire toolchain works! Cross-compiler, assembler, linker, GRUB, QEMU all functioning correctly.

---

## Summary

| Tool | Version | Purpose |
|------|---------|---------|
| i686-elf-gcc | 13.x | Cross-compiler: C → bare metal x86 machine code |
| i686-elf-ld | 2.41 | Linker: combine object files into kernel binary |
| NASM | 2.16 | Assembler: x86 assembly → machine code |
| QEMU | 8.x | Virtual machine: run and test OS without real hardware |
| GDB | 12+ | Debugger: breakpoints, register inspection in running kernel |
| GNU Make | 4.3 | Build automation |
| grub-mkrescue | 2.12 | Create bootable ISO with GRUB bootloader |
| xorriso | 1.5 | ISO image manipulation |
| `-ffreestanding` | GCC flag | No standard library, no OS assumptions |
| `-nostdlib` | GCC flag | Don't link standard libraries |
| linker.ld | Custom | Place kernel at 0x100000, define sections |
| `-s -S` | QEMU flags | Enable remote GDB debug, pause at startup |
