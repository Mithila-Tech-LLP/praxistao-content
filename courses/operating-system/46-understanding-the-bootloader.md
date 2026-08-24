# Chapter 46: Understanding the Bootloader

> **"The bootloader is the first code that runs when your OS boots — before any C, before any standard library, before any file system. It operates with nothing but a CPU in real mode, 512 bytes of storage, and BIOS. Its only job: find and load the kernel. Every line must count."**

---

## Table of Contents

1. [What the BIOS Gives Us](#1-what-the-bios-gives-us)
2. [The Boot Sector — 512 Bytes](#2-the-boot-sector--512-bytes)
3. [BIOS Interrupts — Our Only I/O](#3-bios-interrupts--our-only-io)
4. [Reading from Disk with BIOS](#4-reading-from-disk-with-bios)
5. [Entering Protected Mode](#5-entering-protected-mode)
6. [A Two-Stage Bootloader](#6-a-two-stage-bootloader)
7. [Full Bootloader Code — Stage 1](#7-full-bootloader-code--stage-1)
8. [Full Bootloader Code — Stage 2](#8-full-bootloader-code--stage-2)
9. [Testing the Bootloader](#9-testing-the-bootloader)
10. [Summary](#summary)

---

## 1. What the BIOS Gives Us

When the CPU starts, the PC is:
```
Mode:        Real mode (16-bit, no protection, no paging)
CS:IP:       0xF000:0xFFF0 (BIOS ROM entry point)
Memory:      1MB addressable
BIOS:        Code in ROM initializes hardware, tests RAM (POST)
Then BIOS:   Searches for bootable device
             Reads first 512 bytes of bootable device into 0x7C00
             Checks last 2 bytes for 0x55AA magic
             If found: jumps to 0x7C00 with CS=0x0000, IP=0x7C00
```

**What we have at boot:**
```
AX: 0x0000
BX: 0x0000
DL: boot device number (0x00=floppy, 0x80=first HDD) ← VERY USEFUL!
CS: 0x0000
DS: 0x0000 (often — don't assume)
SS: 0x0000 (often)
SP: 0xFFFE (often — stack at end of conventional memory)

Memory layout:
0x00000 - 0x003FF: IVT (Interrupt Vector Table) — BIOS interrupt handlers
0x00400 - 0x004FF: BDA (BIOS Data Area)
0x00500 - 0x07BFF: free (our stack can go here)
0x07C00 - 0x07DFF: Boot sector loaded here (512 bytes)
0x07E00 - 0x7FFFF: free (load our second stage / kernel here!)
0x80000 - 0x9FFFF: Extended BIOS data area
0xA0000 - 0xBFFFF: VGA memory (text: 0xB8000)
0xC0000 - 0xFFFFF: BIOS ROM, option ROMs
```

---

## 2. The Boot Sector — 512 Bytes

Rules:
- Exactly 512 bytes
- Last 2 bytes must be `0x55, 0xAA` (little-endian 0xAA55)
- Runs in 16-bit real mode at 0x7C00
- Must fit setup + load kernel/stage2 in 510 bytes (512 - 2 for magic)
- If too small: load a second stage to do more work

---

## 3. BIOS Interrupts — Our Only I/O

Before we have our own drivers, BIOS provides basic I/O via `int` instructions:

**Video (INT 0x10):**
```nasm
; Print a character:
mov ah, 0x0E    ; BIOS teletype output function
mov al, 'H'     ; character to print
mov bh, 0       ; page number (0)
int 0x10        ; call BIOS

; Set video mode:
mov ah, 0x00
mov al, 0x03    ; 80×25 text mode (VGA)
int 0x10

; Get cursor position:
mov ah, 0x03
mov bh, 0       ; page
int 0x10        ; returns: DH=row, DL=column
```

**Disk (INT 0x13):**
```nasm
; Read sectors from disk (CHS addressing):
mov ah, 0x02    ; function: read sectors
mov al, 4       ; number of sectors to read
mov ch, 0       ; cylinder 0
mov cl, 2       ; sector 2 (sectors start at 1; sector 1 = boot sector)
mov dh, 0       ; head 0
mov dl, 0x80    ; drive (0x80 = first HDD)
mov bx, 0x7E00  ; ES:BX = destination buffer (ES=0, BX=0x7E00)
int 0x13
; CF=0: success, AL=sectors read; CF=1: error, AH=error code

; Reset disk (call if error):
mov ah, 0x00
mov dl, 0x80
int 0x13
```

**Memory detection (INT 0x15, AX=0xE820):**
```nasm
; Get memory map (crucial for OS — need to know which RAM is usable):
mov eax, 0xE820
mov ebx, 0           ; continuation value (0 for first call)
mov ecx, 24          ; bytes to write per entry
mov edx, 0x534D4150  ; magic signature "SMAP"
mov di, memory_map   ; ES:DI = destination buffer
int 0x15             ; returns: carry set = error; EBX=0 = last entry
; Each entry: 8 bytes base address, 8 bytes length, 4 bytes type
; Type 1 = usable RAM
; Repeat with returned EBX until EBX=0 (last entry)
```

---

## 4. Reading from Disk with BIOS

**CHS (Cylinder-Head-Sector) addressing:**
Old but used by BIOS. Our kernel is stored after the boot sector on the same disk.

```nasm
; Read 'count' sectors starting from LBA 'start' into memory at 'buffer':
load_disk:
    ; Using CHS addressing — for simplicity, assumes head 0, cylinder 0
    ; This works for loading the first ~63 sectors (enough for our kernel)
    
    mov ah, 0x02
    mov al, [count]         ; sectors to read
    mov ch, 0               ; cylinder 0
    mov cl, [start]         ; starting sector (1-based)
    mov dh, 0               ; head 0
    mov dl, 0x80            ; first hard disk
    mov bx, [buffer]        ; buffer in ES:BX
    int 0x13
    
    jc disk_error           ; carry flag = error
    ret

disk_error:
    mov si, disk_err_msg
    call print_string
    hlt
```

---

## 5. Entering Protected Mode

After loading the kernel, we switch from 16-bit real mode to 32-bit protected mode:

```nasm
enter_protected_mode:
    cli                      ; disable interrupts (we don't have an IDT yet)
    
    ; Load the GDT:
    lgdt [gdt_descriptor]
    
    ; Set PE bit in CR0:
    mov eax, cr0
    or eax, 1
    mov cr0, eax
    
    ; Far jump: flush pipeline, reload CS with kernel code selector (0x08):
    jmp 0x08:protected_mode_entry
    
; GDT:
gdt_start:
    ; Null descriptor (required):
    dd 0, 0
    
    ; Kernel code: base=0, limit=4GB, DPL=0, executable, readable
    dw 0xFFFF           ; limit[15:0]
    dw 0x0000           ; base[15:0]
    db 0x00             ; base[23:16]
    db 0b10011010       ; access byte: P=1, DPL=0, S=1, type=0b1010 (code, read)
    db 0b11001111       ; flags+limit[19:16]: G=1, DB=1, L=0, limit=0xF
    db 0x00             ; base[31:24]
    
    ; Kernel data: base=0, limit=4GB, DPL=0, writable
    dw 0xFFFF
    dw 0x0000
    db 0x00
    db 0b10010010       ; access byte: P=1, DPL=0, S=1, type=0b0010 (data, write)
    db 0b11001111
    db 0x00
gdt_end:

gdt_descriptor:
    dw gdt_end - gdt_start - 1   ; GDT limit (size - 1)
    dd gdt_start                   ; GDT base address

[BITS 32]
protected_mode_entry:
    ; Now in 32-bit protected mode!
    ; Reload all segment registers with data selector:
    mov ax, 0x10    ; kernel data segment selector
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    mov ss, ax
    mov esp, 0x90000  ; set up a stack
    
    ; Jump to kernel:
    call kernel_main
    
    ; Shouldn't return, but halt if it does:
    cli
    hlt
```

---

## 6. A Two-Stage Bootloader

**Why two stages?**
Stage 1 (512 bytes = boot sector): too small to do much. Just enough to:
- Set up a basic stack
- Load stage 2 from disk into memory
- Jump to stage 2

Stage 2 (can be larger — a few KB): does the real work:
- Use INT 0x15 E820 to map available memory
- Load the kernel from disk
- Parse kernel ELF header (or just raw binary)
- Switch to protected mode
- Enable paging
- Jump to kernel

For our TinyOS we'll skip the ELF parsing and just use GRUB (Chapter 47) which handles all of this.

But for learning, here's a minimal two-stage setup:

```
Disk layout:
  Sector 1 (LBA 0): Stage 1 bootloader (512 bytes)
  Sector 2 (LBA 1): Stage 2 bootloader (loaded to 0x7E00)
  Sectors 3-50:     Kernel binary (loaded to 0x100000)
```

---

## 7. Full Bootloader Code — Stage 1

```nasm
; boot/stage1.asm — Stage 1 bootloader (must fit in 512 bytes)

[BITS 16]
[ORG 0x7C00]

start:
    ; Clear direction flag and disable interrupts during setup:
    cld
    cli
    
    ; Set up segments and stack:
    xor ax, ax
    mov ds, ax
    mov es, ax
    mov ss, ax
    mov sp, 0x7C00  ; stack grows down from 0x7C00 (below boot sector)
    
    ; Re-enable interrupts:
    sti
    
    ; Save boot drive number:
    mov [boot_drive], dl
    
    ; Print boot message:
    mov si, msg_boot
    call print_string
    
    ; Load Stage 2 from disk:
    ; Read 4 sectors (2KB) starting at sector 2 into 0x7E00
    mov ah, 0x02        ; read sectors
    mov al, 4           ; 4 sectors (2KB)
    mov ch, 0           ; cylinder 0
    mov cl, 2           ; sector 2 (sector after boot sector)
    mov dh, 0           ; head 0
    mov dl, [boot_drive]
    mov bx, 0x7E00      ; load at 0x7E00
    int 0x13
    
    jc disk_error
    
    ; Jump to stage 2:
    jmp 0x0000:0x7E00

disk_error:
    mov si, msg_disk_error
    call print_string
    jmp $               ; hang

; Print string (DS:SI = null-terminated string):
print_string:
    push ax
    push bx
.loop:
    lodsb               ; load byte from [SI] into AL, increment SI
    test al, al
    jz .done            ; if null terminator, done
    mov ah, 0x0E
    mov bh, 0
    int 0x10            ; print character
    jmp .loop
.done:
    pop bx
    pop ax
    ret

; Data:
boot_drive   db 0
msg_boot     db "Stage1 OK", 13, 10, 0
msg_disk_error db "Disk Error!", 13, 10, 0

; Boot sector padding and magic:
times 510 - ($ - $$) db 0
dw 0xAA55
```

---

## 8. Full Bootloader Code — Stage 2

```nasm
; boot/stage2.asm — Stage 2 bootloader

[BITS 16]
[ORG 0x7E00]

start_stage2:
    ; Print stage 2 message:
    mov si, msg_stage2
    call print_string
    
    ; Load kernel into memory at 0x10000 (64KB mark)
    ; Kernel is at sectors 6..50 (45 sectors = 22.5KB)
    mov ah, 0x02
    mov al, 45          ; 45 sectors = 22.5KB
    mov ch, 0
    mov cl, 6           ; starting sector 6
    mov dh, 0
    mov dl, 0x80
    ; We'll load to 0x1000:0x0000 = physical 0x10000:
    push es
    mov ax, 0x1000
    mov es, ax
    xor bx, bx          ; ES:BX = 0x1000:0x0000 = 0x10000
    int 0x13
    pop es
    jc disk_error_s2
    
    ; Detect memory (INT 0x15, E820):
    ; (simplified: just use 0x88 for systems <= 64MB)
    ; Full E820 implementation would go here
    
    ; Enter protected mode:
    cli
    lgdt [gdt_descriptor]
    
    mov eax, cr0
    or eax, 1
    mov cr0, eax
    
    jmp 0x08:protected_entry

disk_error_s2:
    mov si, msg_disk_err2
    call print_string
    jmp $

print_string:
    push ax
    push bx
.loop:
    lodsb
    test al, al
    jz .done
    mov ah, 0x0E
    mov bh, 0
    int 0x10
    jmp .loop
.done:
    pop bx
    pop ax
    ret

; --- GDT (Global Descriptor Table) ---
gdt_start:
    ; Null descriptor:
    dq 0
    
    ; Kernel code segment: base=0, limit=4GB, DPL=0, code
    dw 0xFFFF, 0x0000
    db 0x00, 0b10011010, 0b11001111, 0x00
    
    ; Kernel data segment: base=0, limit=4GB, DPL=0, data
    dw 0xFFFF, 0x0000
    db 0x00, 0b10010010, 0b11001111, 0x00
gdt_end:

gdt_descriptor:
    dw gdt_end - gdt_start - 1
    dd gdt_start

msg_stage2   db "Stage2 OK, loading kernel...", 13, 10, 0
msg_disk_err2 db "Stage2 disk error!", 13, 10, 0

; --- 32-bit protected mode entry ---
[BITS 32]
protected_entry:
    ; Set all data segments to kernel data selector (0x10):
    mov ax, 0x10
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    mov ss, ax
    mov esp, 0x9FC00    ; top of low memory
    
    ; Jump to kernel entry point at 0x10000:
    call 0x10000
    
    ; Should not return:
    cli
    hlt
```

---

## 9. Testing the Bootloader

```bash
# Assemble stage 1:
nasm -f bin boot/stage1.asm -o stage1.bin

# Assemble stage 2:
nasm -f bin boot/stage2.asm -o stage2.bin

# Create a disk image (1.44MB floppy-sized):
dd if=/dev/zero of=disk.img bs=512 count=2880

# Write stage 1 to sector 1:
dd if=stage1.bin of=disk.img bs=512 count=1 conv=notrunc

# Write stage 2 to sectors 2-5:
dd if=stage2.bin of=disk.img bs=512 seek=1 count=4 conv=notrunc

# Write kernel to sectors 6+:
dd if=kernel.bin of=disk.img bs=512 seek=5 conv=notrunc

# Run in QEMU:
qemu-system-i386 -drive format=raw,file=disk.img

# Debug with GDB:
qemu-system-i386 -drive format=raw,file=disk.img -s -S &
gdb -ex "target remote localhost:1234" \
    -ex "set architecture i8086" \
    -ex "break *0x7C00" \
    -ex "continue"
```

---

## Summary

| Concept | Description |
|---------|------------|
| Boot sector | First 512 bytes of bootable disk; must end with 0xAA55; loaded at 0x7C00 |
| Real mode | 16-bit startup mode; BIOS interrupts available; 1MB addressable |
| INT 0x10 | BIOS video services (print character, set mode, get cursor) |
| INT 0x13 | BIOS disk services (read sectors by CHS or LBA) |
| INT 0x15 E820 | Memory map detection; returns list of usable/reserved RAM ranges |
| DL | Boot drive number passed by BIOS (0x80 = first hard disk) |
| Two-stage bootloader | Stage 1: load stage 2; Stage 2: load kernel, switch to protected mode |
| Protected mode entry | cli → lgdt → set CR0.PE → far jump to flush pipeline |
| GDT | Global Descriptor Table; defines code/data segments; loaded with lgdt |
| 0x7E00 | Standard load address for stage 2 bootloader (right after boot sector) |
| 0x100000 | 1MB — standard kernel load address (above real-mode memory hole) |
| Far jump | `jmp 0x08:label` — jumps to CS=0x08, reloads segment register, flushes prefetch |
