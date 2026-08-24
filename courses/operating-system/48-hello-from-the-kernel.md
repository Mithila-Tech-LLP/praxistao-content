# Chapter 48: Hello from the Kernel

> **"Writing to VGA memory is the 'Hello, World!' of OS development — it proves your toolchain works, your bootloader works, your entry point works, and your code runs on bare metal. Those few characters on screen represent the entire stack from transistors to text. Celebrate them."**

---

## Table of Contents

1. [VGA Text Mode Basics](#1-vga-text-mode-basics)
2. [The VGA Buffer Layout](#2-the-vga-buffer-layout)
3. [Colors and Attributes](#3-colors-and-attributes)
4. [Writing Characters — Raw Approach](#4-writing-characters--raw-approach)
5. [Building a Terminal Driver](#5-building-a-terminal-driver)
6. [Scrolling](#6-scrolling)
7. [Printf-Style Formatting](#7-printf-style-formatting)
8. [The Complete vga.c / vga.h](#8-the-complete-vgah)
9. [Testing It All](#9-testing-it-all)
10. [Summary](#summary)

---

## 1. VGA Text Mode Basics

When the machine boots in VGA text mode (80 columns × 25 rows), the hardware maps a special region of memory at physical address **0xB8000**. Writing bytes there directly controls what appears on screen — no system calls, no drivers, no file system. Just memory writes.

```
How VGA text mode works:
  Hardware: VGA controller continuously reads from 0xB8000
            converts bytes to pixels and drives the display
            
  Software: we write bytes to 0xB8000
            hardware picks them up automatically (no "flush" needed)
            
  Result: instantaneous, direct screen output
  
Dimensions: 80 columns × 25 rows = 2000 character cells
Memory size: 2000 × 2 bytes per cell = 4000 bytes total
```

Why 2 bytes per cell? Each cell has:
- Byte 0: ASCII character code
- Byte 1: Attribute (color: foreground + background)

---

## 2. The VGA Buffer Layout

```
VGA buffer at 0xB8000:

Cell (col=0, row=0):  address = 0xB8000 + (0  * 80 + 0)  * 2 = 0xB8000
Cell (col=1, row=0):  address = 0xB8000 + (0  * 80 + 1)  * 2 = 0xB8002
Cell (col=0, row=1):  address = 0xB8000 + (1  * 80 + 0)  * 2 = 0xB80A0
Cell (col=79, row=24): address = 0xB8000 + (24 * 80 + 79) * 2 = 0xB8F9E

Formula:
  address = 0xB8000 + (row * 80 + col) * 2

Memory layout (first few bytes):
  0xB8000: char[0]   attr[0]   char[1]   attr[1]   ...
           'H'       0x0F      'e'       0x0F      ...
           ├── col 0 ────────┤ ├── col 1 ────────┤

As a 16-bit value per cell:
  cell = (attribute << 8) | ascii_char
  'H' white on black = (0x0F << 8) | 'H' = 0x0F48
```

---

## 3. Colors and Attributes

The attribute byte encodes foreground and background color:

```
Attribute byte layout:
  Bit 7:   Blink (on older hardware; modern: high background intensity)
  Bits 6-4: Background color (3 bits = 8 colors, no bright background natively)
  Bit 3:   Bright/intense foreground
  Bits 2-0: Foreground color (3 bits)

Color values:
  0 = Black          8 = Dark gray (bright black)
  1 = Blue           9 = Bright blue
  2 = Green          A = Bright green
  3 = Cyan           B = Bright cyan
  4 = Red            C = Bright red
  5 = Magenta        D = Bright magenta
  6 = Brown/Yellow   E = Yellow (bright)
  7 = Light gray     F = White (bright)

Examples:
  White text on black background:  0x0F = 0000 1111
  Green text on black background:  0x02 = 0000 0010
  White text on blue background:   0x1F = 0001 1111
  Red text on white background:    0x74 = 0111 0100
  Black text on red background:    0x40 = 0100 0000  ← error screen style
```

```c
/* Color definitions (include/vga.h): */

typedef enum {
    VGA_COLOR_BLACK         = 0,
    VGA_COLOR_BLUE          = 1,
    VGA_COLOR_GREEN         = 2,
    VGA_COLOR_CYAN          = 3,
    VGA_COLOR_RED           = 4,
    VGA_COLOR_MAGENTA       = 5,
    VGA_COLOR_BROWN         = 6,
    VGA_COLOR_LIGHT_GREY    = 7,
    VGA_COLOR_DARK_GREY     = 8,
    VGA_COLOR_LIGHT_BLUE    = 9,
    VGA_COLOR_LIGHT_GREEN   = 10,
    VGA_COLOR_LIGHT_CYAN    = 11,
    VGA_COLOR_LIGHT_RED     = 12,
    VGA_COLOR_LIGHT_MAGENTA = 13,
    VGA_COLOR_YELLOW        = 14,
    VGA_COLOR_WHITE         = 15,
} vga_color_t;

static inline uint8_t vga_make_attr(vga_color_t fg, vga_color_t bg) {
    return fg | (bg << 4);
}

static inline uint16_t vga_make_cell(char c, uint8_t attr) {
    return (uint16_t)c | ((uint16_t)attr << 8);
}
```

---

## 4. Writing Characters — Raw Approach

The simplest possible way to write to the screen:

```c
/* Simplest possible "Hello, World!": */
void kernel_main(void) {
    volatile uint16_t *vga = (volatile uint16_t *)0xB8000;
    
    const char *msg = "Hello from TinyOS!";
    uint8_t attr = 0x0F; /* white on black */
    
    for (int i = 0; msg[i] != '\0'; i++) {
        vga[i] = (uint16_t)msg[i] | ((uint16_t)attr << 8);
    }
    
    for (;;) {} /* halt */
}
```

This works but has no cursor tracking, no newlines, no scrolling. We need a proper terminal driver.

---

## 5. Building a Terminal Driver

```c
/* kernel/vga.c — VGA terminal driver */

#include "vga.h"
#include "stdint.h"

#define VGA_COLS    80
#define VGA_ROWS    25
#define VGA_BASE    ((volatile uint16_t *)0xB8000)

/* Terminal state: */
static int term_row;
static int term_col;
static uint8_t term_attr;
static volatile uint16_t *vga_buf;

/* Hardware cursor I/O ports: */
#define VGA_CTRL_REG  0x3D4
#define VGA_DATA_REG  0x3D5

static inline void outb(uint16_t port, uint8_t val) {
    __asm__ volatile ("outb %0, %1" : : "a"(val), "Nd"(port));
}

/* Move the hardware text cursor (the blinking underscore): */
static void update_cursor(void) {
    uint16_t pos = term_row * VGA_COLS + term_col;
    outb(VGA_CTRL_REG, 0x0F);          /* cursor location low byte register */
    outb(VGA_DATA_REG, (uint8_t)(pos & 0xFF));
    outb(VGA_CTRL_REG, 0x0E);          /* cursor location high byte register */
    outb(VGA_DATA_REG, (uint8_t)((pos >> 8) & 0xFF));
}

/* Initialize the terminal (call once at boot): */
void terminal_init(void) {
    vga_buf  = VGA_BASE;
    term_row = 0;
    term_col = 0;
    term_attr = vga_make_attr(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
    
    /* Clear the screen: */
    for (int y = 0; y < VGA_ROWS; y++) {
        for (int x = 0; x < VGA_COLS; x++) {
            vga_buf[y * VGA_COLS + x] = vga_make_cell(' ', term_attr);
        }
    }
    
    update_cursor();
}

/* Set the text color: */
void terminal_set_color(vga_color_t fg, vga_color_t bg) {
    term_attr = vga_make_attr(fg, bg);
}

/* Write a single character: */
void terminal_putchar(char c) {
    if (c == '\n') {
        term_col = 0;
        term_row++;
    } else if (c == '\r') {
        term_col = 0;
    } else if (c == '\t') {
        /* Tab: advance to next 8-column tab stop: */
        term_col = (term_col + 8) & ~7;
    } else if (c == '\b') {
        /* Backspace: */
        if (term_col > 0) {
            term_col--;
            vga_buf[term_row * VGA_COLS + term_col] = vga_make_cell(' ', term_attr);
        }
    } else {
        vga_buf[term_row * VGA_COLS + term_col] = vga_make_cell(c, term_attr);
        term_col++;
    }
    
    /* Wrap at end of line: */
    if (term_col >= VGA_COLS) {
        term_col = 0;
        term_row++;
    }
    
    /* Scroll if we've gone past the last row: */
    if (term_row >= VGA_ROWS) {
        terminal_scroll();
    }
    
    update_cursor();
}

/* Write a null-terminated string: */
void terminal_write(const char *s) {
    for (int i = 0; s[i]; i++) {
        terminal_putchar(s[i]);
    }
}

/* Write len bytes from a buffer (may contain null bytes): */
void terminal_write_n(const char *s, int len) {
    for (int i = 0; i < len; i++) {
        terminal_putchar(s[i]);
    }
}
```

---

## 6. Scrolling

When text reaches the bottom, we need to scroll up:

```c
/* Scroll the terminal up by one line: */
void terminal_scroll(void) {
    /* Move all rows up by one: */
    for (int y = 0; y < VGA_ROWS - 1; y++) {
        for (int x = 0; x < VGA_COLS; x++) {
            vga_buf[y * VGA_COLS + x] = vga_buf[(y + 1) * VGA_COLS + x];
        }
    }
    
    /* Clear the last row: */
    for (int x = 0; x < VGA_COLS; x++) {
        vga_buf[(VGA_ROWS - 1) * VGA_COLS + x] = vga_make_cell(' ', term_attr);
    }
    
    term_row = VGA_ROWS - 1;
}
```

---

## 7. Printf-Style Formatting

We need `printf` but can't use libc. Here's a simple implementation:

```c
/* kernel/printf.c — minimal printf for the kernel */

#include "vga.h"
#include "stdarg.h"  /* va_list — provided by GCC even in freestanding mode */

/* Convert integer to string in given base: */
static void itoa(int64_t value, char *buf, int base, int is_unsigned) {
    static const char digits[] = "0123456789ABCDEF";
    char tmp[32];
    int i = 0;
    int negative = 0;
    uint64_t uval;
    
    if (!is_unsigned && value < 0) {
        negative = 1;
        uval = (uint64_t)(-value);
    } else {
        uval = (uint64_t)value;
    }
    
    if (uval == 0) {
        buf[0] = '0';
        buf[1] = '\0';
        return;
    }
    
    while (uval > 0) {
        tmp[i++] = digits[uval % base];
        uval /= base;
    }
    
    int j = 0;
    if (negative) buf[j++] = '-';
    while (i > 0) buf[j++] = tmp[--i];
    buf[j] = '\0';
}

void kprintf(const char *fmt, ...) {
    va_list args;
    va_start(args, fmt);
    
    char buf[32];
    
    for (int i = 0; fmt[i]; i++) {
        if (fmt[i] != '%') {
            terminal_putchar(fmt[i]);
            continue;
        }
        
        i++; /* move past '%' */
        switch (fmt[i]) {
            case 'd': {
                int val = va_arg(args, int);
                itoa(val, buf, 10, 0);
                terminal_write(buf);
                break;
            }
            case 'u': {
                unsigned val = va_arg(args, unsigned);
                itoa(val, buf, 10, 1);
                terminal_write(buf);
                break;
            }
            case 'x': {
                unsigned val = va_arg(args, unsigned);
                itoa(val, buf, 16, 1);
                terminal_write(buf);
                break;
            }
            case 'p': {
                /* Print pointer as 0xADDRESS */
                terminal_write("0x");
                unsigned val = (unsigned)va_arg(args, void *);
                itoa(val, buf, 16, 1);
                terminal_write(buf);
                break;
            }
            case 's': {
                char *s = va_arg(args, char *);
                terminal_write(s ? s : "(null)");
                break;
            }
            case 'c': {
                char c = (char)va_arg(args, int);
                terminal_putchar(c);
                break;
            }
            case '%': {
                terminal_putchar('%');
                break;
            }
        }
    }
    
    va_end(args);
}
```

---

## 8. The Complete vga.h

```c
/* include/vga.h */
#pragma once
#include "stdint.h"

typedef enum {
    VGA_COLOR_BLACK = 0, VGA_COLOR_BLUE = 1, VGA_COLOR_GREEN = 2,
    VGA_COLOR_CYAN = 3, VGA_COLOR_RED = 4, VGA_COLOR_MAGENTA = 5,
    VGA_COLOR_BROWN = 6, VGA_COLOR_LIGHT_GREY = 7, VGA_COLOR_DARK_GREY = 8,
    VGA_COLOR_LIGHT_BLUE = 9, VGA_COLOR_LIGHT_GREEN = 10,
    VGA_COLOR_LIGHT_CYAN = 11, VGA_COLOR_LIGHT_RED = 12,
    VGA_COLOR_LIGHT_MAGENTA = 13, VGA_COLOR_YELLOW = 14, VGA_COLOR_WHITE = 15,
} vga_color_t;

static inline uint8_t vga_make_attr(vga_color_t fg, vga_color_t bg) {
    return (uint8_t)(fg | (bg << 4));
}
static inline uint16_t vga_make_cell(char c, uint8_t attr) {
    return (uint16_t)c | ((uint16_t)attr << 8);
}

void terminal_init(void);
void terminal_set_color(vga_color_t fg, vga_color_t bg);
void terminal_putchar(char c);
void terminal_write(const char *s);
void terminal_scroll(void);
void kprintf(const char *fmt, ...);
```

---

## 9. Testing It All

Update `kernel_main` to use our new terminal:

```c
/* kernel/kernel.c */
#include "multiboot.h"
#include "vga.h"
#include "stdint.h"

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    
    if (magic != 0x2BADB002) {
        terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_RED);
        terminal_write("FATAL: Not loaded by Multiboot bootloader!\n");
        for (;;) {}
    }
    
    /* Print banner: */
    terminal_set_color(VGA_COLOR_YELLOW, VGA_COLOR_BLACK);
    terminal_write("  _____ _             ___  ____  \n");
    terminal_write(" |_   _(_)_ __  _   _/ _ \\/ ___| \n");
    terminal_write("   | | | | '_ \\| | | | | | \\___ \\ \n");
    terminal_write("   | | | | | | | |_| | |_| |___) |\n");
    terminal_write("   |_| |_|_| |_|\\__, |\\___/|____/ \n");
    terminal_write("                |___/              \n");
    
    terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
    terminal_write("\nTinyOS v0.1 — Chapter 48\n");
    terminal_write("=========================\n\n");
    
    /* Show memory info from Multiboot: */
    struct multiboot_info *mbi = (struct multiboot_info *)mbi_ptr;
    if (mbi->flags & (1 << 0)) {
        kprintf("Lower memory: %u KB\n", mbi->mem_lower);
        kprintf("Upper memory: %u KB (%u MB)\n",
                mbi->mem_upper, mbi->mem_upper / 1024);
    }
    
    kprintf("\nHex: 0x%x  Dec: %d  Str: %s\n", 0xDEAD, 42, "hello!");
    
    terminal_set_color(VGA_COLOR_LIGHT_GREEN, VGA_COLOR_BLACK);
    terminal_write("\nSystem halted. All good!\n");
    
    for (;;) {}
}
```

**Build and run:**
```bash
make && make run
```

Expected output (on a QEMU black screen):
```
  _____ _             ___  ____
 |_   _(_)_ __  _   _/ _ \/ ___|
   | | | | '_ \| | | | | | \___ \
   | | | | | | | |_| | |_| |___) |
   |_| |_|_| |_|\__, |\___/|____/
                |___/

TinyOS v0.1 — Chapter 48
=========================

Lower memory: 639 KB
Upper memory: 130048 KB (127 MB)

Hex: 0xDEAD  Dec: 42  Str: hello!

System halted. All good!
```

---

## Summary

| Concept | Description |
|---------|------------|
| VGA buffer | Physical memory at 0xB8000; 80×25 cells; 2 bytes per cell |
| Cell format | Low byte = ASCII char; high byte = attribute (fg color + bg color) |
| Attribute byte | Bits [3:0] = foreground color; bits [6:4] = background color; bit 7 = blink |
| VGA colors | 16 colors (0-15): black, blue, green, cyan, red, magenta, brown, light grey, 8 bright variants |
| volatile | Required for VGA buffer pointer — prevents compiler from optimizing away writes |
| Hardware cursor | Controlled via I/O ports 0x3D4/0x3D5 (CRT Controller registers 0x0E and 0x0F) |
| Scrolling | Copy rows 1..24 to rows 0..23; clear row 24; set term_row = 24 |
| kprintf | Kernel-side printf using va_list (from GCC freestanding headers) |
| terminal_init | Clear screen, set default color, reset cursor to (0,0) |
