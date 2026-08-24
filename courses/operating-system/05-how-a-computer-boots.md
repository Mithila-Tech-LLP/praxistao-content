# Chapter 05: How a Computer Boots

> **"Booting a computer is a bootstrapping problem — the OS isn't loaded yet, so what loads the OS? Something simpler loads something slightly more complex, which loads something more complex still, until the OS is fully running. It's turtles all the way up."**

---

## Table of Contents

1. [The Boot Problem](#1-the-boot-problem)
2. [Stage 0 — Power On](#2-stage-0--power-on)
3. [Stage 1 — BIOS/UEFI Firmware](#3-stage-1--biosuefi-firmware)
4. [Stage 2 — The Bootloader (MBR)](#4-stage-2--the-bootloader-mbr)
5. [Stage 3 — The Second-Stage Bootloader (GRUB)](#5-stage-3--the-second-stage-bootloader-grub)
6. [Stage 4 — Kernel Loading and Initialization](#6-stage-4--kernel-loading-and-initialization)
7. [Stage 5 — The Init System](#7-stage-5--the-init-system)
8. [Full Boot Sequence Summary](#8-full-boot-sequence-summary)
9. [BIOS Boot vs. UEFI Boot](#9-bios-boot-vs-uefi-boot)
10. [The Boot Sector in Detail (x86 Assembly)](#10-the-boot-sector-in-detail-x86-assembly)
11. [What Happens When the OS Kernel Starts](#11-what-happens-when-the-os-kernel-starts)
12. [Summary](#summary)

---

## 1. The Boot Problem

Here's the fundamental paradox of booting a computer:

**The OS is a program stored on disk. But to run a program from disk, you need an OS. So who loads the OS?**

This is the **bootstrapping problem** (from the expression "pulling yourself up by your own bootstraps").

The solution is a chain of increasingly capable programs:

```
FIRMWARE (in ROM, runs first, no OS needed)
    ↓
BOOTLOADER STAGE 1 (tiny, 512 bytes, in MBR on disk)
    ↓
BOOTLOADER STAGE 2 (larger, can read file systems, find kernel)
    ↓
OS KERNEL (the full operating system, in RAM)
    ↓
INIT SYSTEM (first user-space process)
    ↓
USER SESSION (login screen, desktop, shell)
```

Each stage is simple enough to be loaded without needing the stage above it. Let's trace this chain step by step.

---

## 2. Stage 0 — Power On

You press the power button. Here's what happens in the first fractions of a second:

**0. Power supply:** The power supply converts AC wall power to DC voltages (+12V, +5V, +3.3V). It outputs a "Power Good" signal when voltages stabilize (takes ~100–500 milliseconds).

**1. CPU reset:** When Power Good arrives, the CPU resets. The instruction pointer (RIP/EIP) is set to a hardwired address: `0xFFFFFFF0` (top of 4GB address space, 16 bytes before 4GB boundary).

**2. Jump to firmware:** At physical address `0xFFFFFFF0`, the motherboard maps ROM containing the firmware. The CPU immediately starts executing firmware code.

```
CPU state at power-on:
  CS:IP = 0xF000:0xFFF0  (real mode; 0xF000 << 4 + 0xFFF0 = 0xFFFF0)
  Mode: 16-bit real mode
  All other registers: undefined
```

---

## 3. Stage 1 — BIOS/UEFI Firmware

The firmware is the first code that runs. It has several jobs:

**A. POST (Power-On Self Test):**
The firmware tests hardware:
- Is RAM present and working? (It writes patterns and reads them back)
- Does the CPU work?
- Is there a graphics controller? (Initialize it so errors can be displayed)
- Are required devices present?

If something fails, you get beep codes (before video is available) or on-screen error messages.

**B. Hardware enumeration:**
The firmware discovers all connected devices:
- PCI/PCIe: scans the PCI bus, finds all cards
- USB: enumerates connected USB devices
- SATA/NVMe: finds connected drives
- Builds a table of all discovered hardware (ACPI tables)

**C. Build memory map:**
The firmware determines exactly how much RAM is installed and what physical addresses are usable vs. reserved (for devices, ROM, etc.). It stores this information for the OS.

**D. Find boot device:**
The firmware has a configured boot order (e.g., USB first, then NVMe, then network). It checks each:
1. Is this device present and bootable?
2. Does it have a valid boot sector? (For BIOS: magic bytes `0x55 0xAA` at bytes 510-511)

**E. Load boot sector:**
For legacy BIOS:
```
BIOS reads exactly 512 bytes from the first sector of the boot device
BIOS loads those 512 bytes into RAM at physical address 0x7C00
BIOS jumps to 0x7C00
```

From this point, the code at `0x7C00` takes over. The BIOS is done. Control is transferred.

**What state the BIOS leaves for the bootloader:**
```
CPU mode: 16-bit real mode
DL register: boot drive number (0x80 for first hard drive)
CS:IP: 0x0000:0x7C00 (i.e., the bootloader starts executing)
A20 line: may or may not be enabled (more on this in Volume 9)
```

---

## 4. Stage 2 — The Bootloader (MBR)

The **MBR (Master Boot Record)** is the first 512 bytes of a disk. Its layout:

```
Offset  Size     Content
──────────────────────────────────────────────────────────────────
0x000   446 bytes  Bootstrap code (the actual first-stage bootloader)
0x1BE   64 bytes   Partition table (4 entries × 16 bytes each)
0x1FE   2 bytes    Magic: 0x55 0xAA (marks this as a valid boot sector)
```

**The first-stage bootloader (446 bytes) has ONE job:**
Load the second-stage bootloader from disk and jump to it.

446 bytes is not enough to do much. No file system reading. No complex logic. Just:
1. Find the active partition in the partition table
2. Read the first sector of that partition into memory
3. Jump to it

**GRUB's first stage** is exactly this — a tiny MBR boostrap code that loads GRUB's second stage from a specific location on disk.

**Important limitation:**
In BIOS mode, the bootloader can only use BIOS disk services to read from disk. BIOS disk services use the old Cylinder-Head-Sector (CHS) addressing, which can't address disks beyond 8GB. Modern bootloaders use LBA (Logical Block Addressing) via INT 0x13 extensions to work around this.

---

## 5. Stage 3 — The Second-Stage Bootloader (GRUB)

**GRUB (Grand Unified Bootloader)** is the bootloader used by most Linux distributions. (Windows uses its own: Windows Boot Manager.)

GRUB's job is sophisticated:
1. **Switch to a more capable mode** (32-bit protected mode) to access more RAM
2. **Read the file system** (EXT4, FAT32, Btrfs, etc.) on the boot partition
3. **Read its configuration file** (`/boot/grub/grub.cfg`)
4. **Display a menu** (if configured) letting you choose which OS or kernel to boot
5. **Load the kernel image** into memory
6. **Load the initial RAM disk** (initrd/initramfs) into memory
7. **Set up the machine state** the kernel expects (Multiboot spec)
8. **Jump to the kernel entry point**

**GRUB's configuration example:**
```
# /boot/grub/grub.cfg
menuentry "Ubuntu 24.04" {
    linux /boot/vmlinuz-6.8.0 root=/dev/sda2 quiet splash
    initrd /boot/initrd.img-6.8.0
}
menuentry "Ubuntu recovery" {
    linux /boot/vmlinuz-6.8.0 root=/dev/sda2 recovery nomodeset
    initrd /boot/initrd.img-6.8.0
}
```

**What is `initramfs` (initial RAM filesystem)?**
The kernel needs some drivers just to read the real root file system from disk. But those drivers are ON the disk. Chicken-and-egg problem.

Solution: A tiny compressed file system (initramfs) is loaded into RAM alongside the kernel. It contains just enough drivers and scripts to:
- Load disk driver (e.g., NVMe driver)
- Read the real root partition
- Mount it as the real root
- Hand off to the real init system

---

## 6. Stage 4 — Kernel Loading and Initialization

GRUB (or another bootloader) decompresses and loads the kernel into memory, then jumps to the kernel's entry point.

**The Linux kernel image is a compressed executable.** At its start is a small decompressor that uncompresses the real kernel, then jumps to `startup_32` (or `startup_64` for 64-bit).

**What the kernel does when it first starts executing:**

```
1. Verify CPU capabilities (is this CPU supported?)
2. Set up early paging (turn on virtual memory — must happen very early)
3. Set up the early exception handlers
4. Detect physical memory via BIOS/UEFI memory map
5. Copy kernel to its final memory location (if needed)
6. Set up the final GDT (Global Descriptor Table)
7. Set up the IDT (Interrupt Descriptor Table)
8. Initialize the CPU scheduler
9. Initialize the memory subsystem (buddy allocator, slab allocator)
10. Initialize the VFS (Virtual File System)
11. Mount the initramfs as root
12. Initialize each device driver
13. Start kernel threads (kthreadd, ksoftirqd, etc.)
14. Run /sbin/init (or whatever init= kernel parameter says)
```

By step 14, the kernel is fully initialized. User space begins.

---

## 7. Stage 5 — The Init System

The first user-space program ever started is **init** (PID 1). Every other process on the system is a descendant of PID 1.

Init's job: Start all the system services needed for a working system.

**Modern Linux: systemd**
```
systemd starts:
├── systemd-journald (logging)
├── udevd (hot-plug device management)
├── networkd (networking)
├── sshd (SSH server, if enabled)
├── dbus (message bus between processes)
├── display manager (login screen)
│   └── desktop session
│       ├── Xorg or Wayland (display server)
│       ├── window manager
│       └── terminal, browser, etc. (user apps)
└── ... (hundreds of services for a full desktop)
```

**Older Linux: SysV init**
SysV init ran shell scripts in numbered runlevels. Systemd replaced it because it's faster (starts services in parallel) and more capable (dependency management, socket activation, etc.).

**macOS: launchd**
macOS uses launchd as PID 1. It's similar to systemd — launches services on demand.

**Android: init**
Android's init is a custom binary that reads `.rc` configuration files to start Android system services (SurfaceFlinger for graphics, ActivityManager for apps, etc.).

---

## 8. Full Boot Sequence Summary

```
Power On
    ↓
CPU jumps to 0xFFFFFFF0
    ↓
BIOS/UEFI firmware runs
  - POST: test hardware
  - Enumerate devices
  - Build memory map
  - Find boot device
    ↓
MBR (512 bytes) at 0x7C00
  - Tiny bootstrap: finds active partition
    ↓
GRUB Stage 2 (several hundred KB)
  - Switch to protected mode
  - Read file system
  - Load kernel + initramfs
  - Set up kernel parameters
    ↓
Kernel decompresses and starts
  - Set up virtual memory
  - Initialize subsystems
  - Start kernel threads
    ↓
/sbin/init (PID 1) starts
  - Start system services
    ↓
Login screen / shell prompt
    ↓
User logs in
    ↓
Desktop or shell session
```

**Typical boot time breakdown:**
- BIOS/UEFI firmware: 1–3 seconds
- GRUB loading + menu timeout: 0–3 seconds
- Kernel initialization: 0.5–2 seconds
- systemd starting services: 3–10 seconds
- Total: 5–20 seconds (fast) to 30–60 seconds (slow HDD, many services)

---

## 9. BIOS Boot vs. UEFI Boot

**BIOS boot (legacy, still widely used for learning):**
```
Disk layout:
  ┌──────────────────────────────────────────────────────┐
  │ MBR (512 bytes) │ Partition 1 │ Partition 2 │ ...   │
  │ 0x7C00          │ (boot part) │             │       │
  └──────────────────────────────────────────────────────┘
  MBR contains: bootstrap code + partition table
  Partition table: max 4 primary partitions, max 2TB disk
```

**UEFI boot (modern):**
```
Disk layout (GPT — GUID Partition Table):
  ┌───────────────────────────────────────────────────────────┐
  │ Protective MBR │ GPT Header │ GPT Partitions... │        │
  │                │            │                    │        │
  │ Partition 1: EFI System Partition (FAT32)        │        │
  │   /EFI/BOOT/BOOTX64.EFI  ← UEFI loads this     │        │
  │   /EFI/ubuntu/grubx64.efi                        │        │
  │ Partition 2: /boot (ext4)                        │        │
  │ Partition 3: / (root ext4)                       │        │
  └───────────────────────────────────────────────────────────┘
```

**UEFI differences:**
- Reads FAT32 partitions directly (no 512-byte limitation)
- Loads `.efi` executables (PE32+ format, like Windows executables)
- Has its own boot manager (can store multiple boot entries in NVRAM)
- Secure Boot: verifies EFI executables are signed by trusted keys
- Supports disks > 2TB (GPT vs. MBR)
- Can boot in 32 or 64-bit mode directly

For our OS course, we'll use BIOS/Multiboot2 — simpler and better supported by QEMU out of the box.

---

## 10. The Boot Sector in Detail (x86 Assembly)

Let's look at the absolute minimum x86 assembly boot sector — the foundation of everything:

```nasm
; boot.asm — Minimal x86 boot sector
; Assemble with: nasm -f bin -o boot.bin boot.asm
; Run with: qemu-system-x86_64 -drive format=raw,file=boot.bin

[BITS 16]           ; We start in 16-bit real mode
[ORG 0x7C00]        ; BIOS loads us at physical address 0x7C00

start:
    ; Clear segment registers
    xor ax, ax      ; ax = 0
    mov ds, ax      ; data segment = 0
    mov es, ax      ; extra segment = 0
    mov ss, ax      ; stack segment = 0
    mov sp, 0x7C00  ; stack grows downward from 0x7C00

    ; Print a character (BIOS interrupt 0x10, function 0x0E)
    mov ah, 0x0E    ; function: teletype output
    mov al, 'H'     ; character to print
    int 0x10        ; call BIOS video service
    
    mov al, 'i'
    int 0x10
    
    ; Hang forever
hang:
    hlt             ; halt CPU until next interrupt
    jmp hang        ; if interrupt woke us, halt again

; Padding and magic number
times 510 - ($ - $$) db 0   ; pad to 510 bytes with zeros
dw 0xAA55                    ; magic number: bytes 510-511 = 0x55, 0xAA
                              ; (little-endian: stored as AA 55 on disk)
```

**Every detail matters:**
- `[ORG 0x7C00]` — tells NASM that addresses start at 0x7C00 (where BIOS loaded us)
- `[BITS 16]` — generate 16-bit code (real mode)
- `int 0x10` — BIOS video interrupt; the OS doesn't exist yet so we use BIOS services
- `times 510 - ($ - $$) db 0` — pad to exactly 510 bytes
- `dw 0xAA55` — the 2-byte magic signature BIOS checks before jumping to us

This is the absolute foundation. In Volume 9, we'll build on this to create a real OS.

---

## 11. What Happens When the OS Kernel Starts

When GRUB jumps to the kernel entry point, the kernel starts in a half-initialized state:

**What's set up (by GRUB/Multiboot):**
- CPU is in 32-bit protected mode (not real mode anymore)
- Basic memory map provided by BIOS is available
- Kernel code is in RAM
- Stack is minimal (GRUB set up a small temporary stack)

**What's NOT set up:**
- Interrupts are disabled
- Paging (virtual memory) is not enabled
- No console output (no driver yet)
- No heap (no malloc)
- No processes (no scheduler)
- Nothing

**First things the kernel must do (in order):**
1. Set up its own stack (allocate kernel stack area)
2. Clear BSS segment (uninitialized globals, required by C standard)
3. Set up early console (so we can print error messages)
4. Parse the Multiboot info structure (find how much RAM we have)
5. Set up GDT (Global Descriptor Table — segment descriptors)
6. Set up IDT (Interrupt Descriptor Table — interrupt handlers)
7. Set up paging (turn on virtual memory)
8. Initialize the physical memory allocator
9. Initialize the heap
10. Initialize the scheduler
11. Start the first process

At the end of this, the kernel is fully operational and can run user programs.

---

## Summary

| Stage | What runs | Loaded by | Location |
|-------|----------|-----------|---------|
| Firmware | BIOS/UEFI | CPU hardwired | ROM chip on motherboard |
| Stage 1 | MBR bootloader | BIOS (loads to 0x7C00) | First 512 bytes of disk |
| Stage 2 | GRUB/Windows Boot Manager | MBR code | Disk partition / EFI partition |
| Kernel | Linux kernel / NT kernel | GRUB | /boot/vmlinuz or FAT32 partition |
| Init | systemd / launchd / Android init | Kernel | /sbin/init |
| Services | daemons, drivers, desktop | Init | Everywhere in user space |

**Key addresses (BIOS boot):**

| Address | Meaning |
|---------|---------|
| `0xFFFFFFF0` | First instruction CPU executes (jumps to BIOS ROM) |
| `0x7C00` | Where BIOS loads the 512-byte boot sector |
| `0xB8000` | VGA text mode buffer (write here → text appears on screen) |
| `0x100000` | 1MB mark; kernel typically loaded above here |
