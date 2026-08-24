# Chapter 60: Keyboard Driver

> **"A keyboard driver turns electrical signals into characters. The hardware sends a scan code — a number representing which key was pressed or released. The driver translates that number to ASCII, handles shift/caps lock/special keys, and puts the character into a buffer where the shell can read it. It's a complete I/O pipeline from physics to text."**

---

## Table of Contents

1. [How PS/2 Keyboards Work](#1-how-ps2-keyboards-work)
2. [Scan Codes — Set 1](#2-scan-codes--set-1)
3. [The Keyboard Input Buffer](#3-the-keyboard-input-buffer)
4. [The Complete Keyboard Driver](#4-the-complete-keyboard-driver)
5. [Shift, Caps Lock, Special Keys](#5-shift-caps-lock-special-keys)
6. [Blocking Read — Waiting for Input](#6-blocking-read--waiting-for-input)
7. [Connecting to sys_read](#7-connecting-to-sys_read)
8. [Complete keyboard.c / keyboard.h](#8-complete-keyboardh)
9. [Testing](#9-testing)
10. [Summary](#summary)

---

## 1. How PS/2 Keyboards Work

```
Physical key press → keyboard controller → IRQ1 → OS

Hardware flow:
  1. Key pressed on keyboard
  2. Keyboard microcontroller encodes it as a scan code
  3. Keyboard sends scan code to PS/2 controller (i8042 chip)
  4. i8042 raises IRQ1 (hardware interrupt)
  5. CPU pauses execution → calls IRQ1 handler
  6. Handler reads scan code from port 0x60
  7. Translates scan code to key meaning
  8. Puts character into keyboard buffer
  9. Shell's read() call unblocks and returns the character

I/O Ports:
  0x60: Data port   — read scan code from here
  0x64: Status/Cmd  — status byte: bit 0 = output buffer full (read when 1)
                                  bit 1 = input buffer full  (write when 0)

Reading a scan code:
  uint8_t scancode = inb(0x60);
  
That's it — one I/O port read gives you the scan code.
```

---

## 2. Scan Codes — Set 1

Scan code set 1 (default for compatibility):

```
Key press:   scan code 0x01-0x58 (press code)
Key release: scan code = press code + 0x80 (bit 7 set)

Example:
  'A' press:   0x1E
  'A' release: 0x9E (0x1E | 0x80)

Extended keys (two bytes):
  E0 prefix → Fn keys, arrows, Insert, Delete, Home, End, etc.
  0xE0 0x48 = Up arrow press
  0xE0 0x50 = Down arrow press

Scan code map (partial):
  0x01: Escape      0x0E: Backspace   0x0F: Tab
  0x1C: Enter       0x2A: Left Shift  0x36: Right Shift
  0x38: Left Alt    0x1D: Left Ctrl   0x3A: Caps Lock

Key positions 1-14 (top row):
  0x02='1'  0x03='2'  0x04='3'  0x05='4'  0x06='5'
  0x07='6'  0x08='7'  0x09='8'  0x0A='9'  0x0B='0'
  0x0C='-'  0x0D='='  0x0E=Backspace

Key positions (QWERTY row):
  0x10='q'  0x11='w'  0x12='e'  0x13='r'  0x14='t'
  0x15='y'  0x16='u'  0x17='i'  0x18='o'  0x19='p'
  0x1A='['  0x1B=']'  0x1C=Enter
  
Key positions (ASDF row):
  0x1E='a'  0x1F='s'  0x20='d'  0x21='f'  0x22='g'
  0x23='h'  0x24='j'  0x25='k'  0x26='l'  0x27=';'
  0x28='\''
  
Key positions (ZXCV row):
  0x2C='z'  0x2D='x'  0x2E='c'  0x2F='v'  0x30='b'
  0x31='n'  0x32='m'  0x33=','  0x34='.'  0x35='/'
  
Space bar: 0x39
```

---

## 3. The Keyboard Input Buffer

A circular (ring) buffer decouples the IRQ handler from the consumer:

```
Ring buffer:
  head: next read position
  tail: next write position
  buf[]:  storage

Write (from IRQ handler):
  buf[tail] = char
  tail = (tail + 1) % BUFFER_SIZE

Read (from sys_read / shell):
  if head == tail: buffer empty, block (wait)
  char = buf[head]
  head = (head + 1) % BUFFER_SIZE
  return char

Full: (tail + 1) % SIZE == head → drop character
Empty: head == tail

Benefits:
  IRQ handler is fast (no blocking, just write + advance pointer)
  Reader can take its time
  Up to BUFFER_SIZE - 1 characters can be buffered
```

---

## 4. The Complete Keyboard Driver

```c
/* kernel/keyboard.c */

#include "keyboard.h"
#include "irq.h"
#include "pic.h"
#include "io.h"
#include "vga.h"
#include "process.h"
#include "scheduler.h"

#define KBD_DATA_PORT  0x60
#define KBD_BUF_SIZE   256

/* Scan code → ASCII mapping (US QWERTY, unshifted): */
static const char scancode_lower[128] = {
    0,   27,  '1','2','3','4','5','6','7','8','9','0','-','=', '\b',
    '\t','q', 'w','e','r','t','y','u','i','o','p','[',']', '\n',
    0,   'a', 's','d','f','g','h','j','k','l',';','\'','`',
    0,   '\\','z','x','c','v','b','n','m',',','.','/', 0,
    '*', 0,   ' ', 0,
    /* F1-F10: */
    0,0,0,0,0,0,0,0,0,0,
    0,0,                     /* NumLock, ScrollLock */
    '7','8','9','-',         /* numpad */
    '4','5','6','+',
    '1','2','3','0','.', 0,0,0,
    0,0                      /* F11, F12 */
};

/* Shifted version: */
static const char scancode_upper[128] = {
    0,   27,  '!','@','#','$','%','^','&','*','(',')','_','+', '\b',
    '\t','Q', 'W','E','R','T','Y','U','I','O','P','{','}', '\n',
    0,   'A', 'S','D','F','G','H','J','K','L',':','"','~',
    0,   '|', 'Z','X','C','V','B','N','M','<','>','?', 0,
    '*', 0,   ' ', 0,
    0,0,0,0,0,0,0,0,0,0,
    0,0,
    '7','8','9','-',
    '4','5','6','+',
    '1','2','3','0','.', 0,0,0,
    0,0
};

/* Input ring buffer: */
static char     kbd_buf[KBD_BUF_SIZE];
static uint32_t kbd_head = 0;
static uint32_t kbd_tail = 0;

/* Modifier key state: */
static int shift_pressed   = 0;
static int ctrl_pressed    = 0;
static int alt_pressed     = 0;
static int caps_lock       = 0;

/* Process waiting for keyboard input (for blocking read): */
static process_t *kbd_waiter = NULL;

/* Put a character into the ring buffer: */
static void kbd_buf_put(char c) {
    uint32_t next = (kbd_tail + 1) % KBD_BUF_SIZE;
    if (next == kbd_head) return;  /* Buffer full — drop character */
    kbd_buf[kbd_tail] = c;
    kbd_tail = next;
    
    /* Wake up any process waiting for keyboard input: */
    if (kbd_waiter) {
        unblock_process(kbd_waiter);
        kbd_waiter = NULL;
    }
}

/* Read one character from the ring buffer (returns 0 if empty): */
char kbd_buf_get(void) {
    if (kbd_head == kbd_tail) return 0;  /* Empty */
    char c = kbd_buf[kbd_head];
    kbd_head = (kbd_head + 1) % KBD_BUF_SIZE;
    return c;
}

/* Check if buffer has data: */
int kbd_buf_available(void) {
    return kbd_head != kbd_tail;
}

/* IRQ1 handler: */
static void keyboard_irq_handler(registers_t *regs) {
    (void)regs;
    
    uint8_t scancode = inb(KBD_DATA_PORT);
    
    /* Handle extended key prefix (0xE0): */
    if (scancode == 0xE0) {
        /* Read next byte (the actual key): */
        /* For simplicity, we'll just read and ignore for now: */
        /* Full implementation would handle arrows, Page Up/Down, etc. */
        /* inb(KBD_DATA_PORT); */
        return;
    }
    
    /* Key release (bit 7 set): */
    if (scancode & 0x80) {
        uint8_t release = scancode & 0x7F;
        if (release == 0x2A || release == 0x36) shift_pressed = 0;  /* Shift up */
        if (release == 0x1D) ctrl_pressed = 0;                       /* Ctrl up */
        if (release == 0x38) alt_pressed  = 0;                       /* Alt up */
        return;
    }
    
    /* Key press: */
    switch (scancode) {
        case 0x2A: case 0x36: shift_pressed = 1; return;  /* Left/Right Shift */
        case 0x1D:             ctrl_pressed  = 1; return;  /* Left Ctrl */
        case 0x38:             alt_pressed   = 1; return;  /* Left Alt */
        case 0x3A:             caps_lock ^= 1;    return;  /* Caps Lock toggle */
        case 0x01:             /* Escape */
            kbd_buf_put(27);
            return;
        default: break;
    }
    
    if (scancode >= 128) return;  /* Unknown key */
    
    /* Translate to ASCII: */
    char c;
    int use_upper = shift_pressed ^ caps_lock;
    c = use_upper ? scancode_upper[scancode] : scancode_lower[scancode];
    
    if (!c) return;  /* Non-printable key */
    
    /* Ctrl+key combinations: */
    if (ctrl_pressed) {
        if (c >= 'a' && c <= 'z') c -= 96;   /* Ctrl+a = 1, Ctrl+b = 2, etc. */
        if (c >= 'A' && c <= 'Z') c -= 64;   /* Ctrl+A = 1 */
        /* Ctrl+C = 3, Ctrl+D = 4 */
    }
    
    /* Echo to screen: */
    terminal_putchar(c);
    
    /* Put in buffer: */
    kbd_buf_put(c);
}

/* Initialize the keyboard driver: */
void keyboard_init(void) {
    irq_register(1, keyboard_irq_handler);
    pic_unmask_irq(1);
}
```

---

## 5. Shift, Caps Lock, Special Keys

```
Modifier logic:
  Shift: held down → use scancode_upper[]
  Caps Lock: toggle → affects letters only (A-Z become uppercase)
  Combined: shift + caps lock → back to lowercase letters
  
  use_upper = shift_pressed XOR caps_lock
  
  For letters: use_upper → uppercase; else lowercase
  For numbers/symbols: only shift matters (caps lock doesn't affect '1'→'!')
  
Special keys we handle:
  Backspace (0x0E): send '\b'; terminal_putchar handles destructive backspace
  Tab (0x0F):       send '\t'
  Enter (0x1C):     send '\n'
  Escape (0x01):    send 0x1B (ASCII ESC)
  
Ctrl combinations:
  Ctrl+C → ASCII 3  (traditionally: interrupt/kill)
  Ctrl+D → ASCII 4  (traditionally: EOF)
  Ctrl+Z → ASCII 26 (traditionally: suspend)
```

---

## 6. Blocking Read — Waiting for Input

```c
/* Read one character, blocking until one is available: */
char kbd_getchar_blocking(void) {
    while (!kbd_buf_available()) {
        /* Nothing in buffer — block this process: */
        kbd_waiter = current_process;
        block_current();   /* State → BLOCKED, calls schedule() */
        /* When unblock_process() is called (from IRQ handler),
           we'll be rescheduled and loop again to check the buffer. */
    }
    return kbd_buf_get();
}

/* Read a line (until '\n'), storing in buf[0..max-1], null-terminated: */
int kbd_readline(char *buf, int max) {
    int len = 0;
    while (len < max - 1) {
        char c = kbd_getchar_blocking();
        
        if (c == '\b') {
            /* Backspace: */
            if (len > 0) {
                len--;
                /* Erase character on screen: */
                terminal_putchar('\b');
                terminal_putchar(' ');
                terminal_putchar('\b');
            }
            continue;
        }
        
        if (c == '\n') {
            buf[len] = '\0';
            terminal_putchar('\n');
            return len;
        }
        
        buf[len++] = c;
    }
    buf[len] = '\0';
    return len;
}
```

---

## 7. Connecting to sys_read

Update the `sys_read` syscall (from Chapter 58) to use the keyboard:

```c
/* In syscall.c, update sys_read: */
static int32_t sys_read(registers_t *r) {
    int      fd    = (int)r->ebx;
    char    *buf   = (char *)r->ecx;
    uint32_t count = r->edx;
    
    if (!buf || count == 0) return -1;
    
    if (fd == 0) {  /* stdin — keyboard */
        uint32_t i = 0;
        while (i < count) {
            char c = kbd_getchar_blocking();
            buf[i++] = c;
            if (c == '\n') break;  /* Stop at newline (line-buffered) */
        }
        return (int32_t)i;
    }
    
    return -1;
}
```

---

## 8. Complete keyboard.h

```c
/* include/keyboard.h */
#pragma once
#include "stdint.h"

void keyboard_init(void);
char kbd_buf_get(void);
int  kbd_buf_available(void);
char kbd_getchar_blocking(void);
int  kbd_readline(char *buf, int max);
```

---

## 9. Testing

```c
/* In kernel_main, after initializing everything: */

void keyboard_test_process(void) {
    kprintf("\nKeyboard test: type something and press Enter\n");
    kprintf("> ");
    
    char line[80];
    while (1) {
        int len = kbd_readline(line, sizeof(line));
        kprintf("You typed (%d chars): '%s'\n", len, line);
        kprintf("> ");
        
        /* Handle some commands: */
        if (line[0] == 'q' && line[1] == '\0') {
            kprintf("Goodbye!\n");
            process_exit(0);
        }
        if (line[0] == 'p' && line[1] == '\0') {
            process_print_all();
        }
    }
}

process_create("kbd-test", keyboard_test_process, 5);
```

You should be able to:
1. Type characters and see them appear on screen
2. Use Backspace to correct mistakes
3. Press Enter to submit the line
4. See the echo back: `You typed (5 chars): 'hello'`

---

## Summary

| Concept | Description |
|---------|------------|
| PS/2 keyboard | Sends scan codes to port 0x60; raises IRQ1 for each key event |
| Scan code | Number representing which physical key was pressed (1-127 for press, +128 for release) |
| Scan code set 1 | Original IBM PC format; press = code, release = code | 0x80 |
| Release event | Scancode with bit 7 set: `scancode & 0x80` is non-zero |
| Port 0x60 | Read from here to get the scan code |
| Ring buffer | Fixed-size circular queue; IRQ handler writes, shell reads; decouples producer/consumer |
| kbd_head / kbd_tail | Read/write pointers into the circular buffer; head==tail means empty |
| Modifier state | Track shift/ctrl/alt with booleans; set on press, clear on release |
| caps_lock | Toggle on press (XOR with shift for letter case) |
| kbd_waiter | Process waiting for input; IRQ handler calls unblock_process(kbd_waiter) when data arrives |
| kbd_getchar_blocking | Blocks current process if buffer empty; unblocked by IRQ handler |
| kbd_readline | Read until '\n'; handle backspace; return null-terminated string |
| Echo | IRQ handler calls terminal_putchar(c) to show typed characters on screen |
