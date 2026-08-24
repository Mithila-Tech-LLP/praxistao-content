# Chapter 18: Numbers and Binary — How Computers Count

> **"Computers only understand two things: on and off, true and false, 1 and 0. Yet from these two things, we build video games, music, movies, AI, and the entire internet. The magic is not the 0s and 1s — the magic is what you can build by combining billions of them."**

---

## Table of Contents

1. [Why Computers Use Binary](#1-why-computers-use-binary)
2. [What Is Binary?](#2-what-is-binary)
3. [Counting in Binary](#3-counting-in-binary)
4. [Bits and Bytes](#4-bits-and-bytes)
5. [How Binary Represents Text](#5-how-binary-represents-text)
6. [How Binary Represents Colors and Images](#6-how-binary-represents-colors-and-images)
7. [Hexadecimal — A Shorthand for Binary](#7-hexadecimal--a-shorthand-for-binary)
8. [Summary](#summary)

---

## 1. Why Computers Use Binary

Why do computers use 0s and 1s instead of 0–9 like humans?

```
The answer is electricity.
  
  A wire can carry electricity or not.
  ON = electricity flowing = 1
  OFF = no electricity = 0
  
  These two states are easy to distinguish reliably.
  A wire at 3.3V is ON. At 0V is OFF. No confusion.
  
  What if we tried 10 states (0–9)?
    0V = 0, 0.3V = 1, 0.6V = 2, ..., 2.7V = 9
    Voltage fluctuates. Hard to reliably distinguish 0.6V from 0.7V.
    Errors would be constant.
  
  Binary is reliable because the gap between 0 and 1 is big.
  A transistor is either fully ON or fully OFF — no ambiguity.
  
  This is why every digital device in existence uses binary.
  It's not a choice — it's physics.
```

---

## 2. What Is Binary?

We count in **decimal** — base 10, using digits 0–9.

Computers count in **binary** — base 2, using only 0 and 1.

```
Decimal (base 10):
  Uses 10 symbols: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
  Each position is worth 10× the previous position
  
  Example: 342
    3 × 100  = 300
    4 × 10   = 40
    2 × 1    = 2
    Total    = 342

Binary (base 2):
  Uses 2 symbols: 0, 1
  Each position is worth 2× the previous position
  
  Example: 1011
    1 × 8  = 8
    0 × 4  = 0
    1 × 2  = 2
    1 × 1  = 1
    Total  = 11 (in decimal)
```

---

## 3. Counting in Binary

```
Decimal  Binary     Think of it as...
0        0000       All lights off
1        0001       Last light on
2        0010       Second light on
3        0011       Last two lights on
4        0100       Third light on
5        0101
6        0110
7        0111       Last three lights on
8        1000       Just the first light
9        1001
10       1010
11       1011
12       1100
13       1101
14       1110
15       1111       All four lights on

Pattern: each time you run out of symbols, you carry to the next position.
Same rule as decimal (9 + 1 = 10 — carry the 1), but with 1 + 1 = 10 in binary.
```

---

## 4. Bits and Bytes

```
Bit:  One binary digit. Either 0 or 1.
      The smallest unit of information.
      
Byte: 8 bits grouped together.
      Why 8? It's a convenient power of 2 that can represent 256 values (0–255).
      
  1 bit    → 2 possible values   (0 or 1)
  2 bits   → 4 possible values   (00, 01, 10, 11)
  4 bits   → 16 possible values  (0–15)
  8 bits   → 256 possible values (0–255)
  16 bits  → 65,536 values
  32 bits  → 4,294,967,296 values
  64 bits  → 18,446,744,073,709,551,616 values (18 quintillion)

This is why:
  RGB color has values 0–255 (fits in 1 byte each)
  Old games had 256 colors (1 byte per pixel)
  Old PCs were "32-bit" — addresses up to 4GB RAM (2^32 bytes)
  Modern PCs are "64-bit" — addresses up to 16 exabytes
```

---

## 5. How Binary Represents Text

Every letter, number, and symbol has a number. That number is stored in binary.

```
ASCII (American Standard Code for Information Interchange):
  Assigns numbers 0–127 to characters.
  
  A = 65 = 01000001
  B = 66 = 01000010
  C = 67 = 01000011
  a = 97 = 01100001
  b = 98 = 01100010
  0 = 48 = 00110000
  1 = 49 = 00110001
  ! = 33 = 00100001
  Space = 32 = 00100000
  
"Hello" in binary:
  H = 72  = 01001000
  e = 101 = 01100101
  l = 108 = 01101100
  l = 108 = 01101100
  o = 111 = 01101111
  
  That's 5 bytes of data (5 × 8 bits = 40 bits).
  
Unicode:
  ASCII only covers English.
  Unicode covers ALL writing systems: 144,000+ characters.
  Chinese, Arabic, Hindi, emoji 😊 — all have Unicode values.
  
  😊 = U+1F60A = 11110000 10011111 10011000 10001010 (4 bytes in UTF-8)
```

---

## 6. How Binary Represents Colors and Images

```
Every color is a combination of Red, Green, and Blue light.
  Each channel: 0 (no color) to 255 (full color).
  Stored as 1 byte each = 3 bytes per pixel.
  
Color examples:
  Red:     R=255, G=0,   B=0    (11111111 00000000 00000000)
  Green:   R=0,   G=255, B=0    (00000000 11111111 00000000)
  Blue:    R=0,   G=0,   B=255  (00000000 00000000 11111111)
  White:   R=255, G=255, B=255  (11111111 11111111 11111111)
  Black:   R=0,   G=0,   B=0    (00000000 00000000 00000000)
  Yellow:  R=255, G=255, B=0    (11111111 11111111 00000000)

Image = a grid of pixels, each with 3 bytes.
  
  1920 × 1080 image:
    = 2,073,600 pixels
    × 3 bytes per pixel
    = 6,220,800 bytes
    = ~6 MB uncompressed
    
    JPEG compression: reduces to ~300KB by finding patterns and approximating.
    You lose some detail, but the file is 20× smaller.
```

---

## 7. Hexadecimal — A Shorthand for Binary

Binary is powerful but verbose. Hexadecimal (hex) is a shorthand that programmers use.

```
Hexadecimal = base 16.
Uses: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, A, B, C, D, E, F
(A=10, B=11, C=12, D=13, E=14, F=15)

The nice property: each hex digit = exactly 4 bits.
So one byte = exactly 2 hex digits.

Byte  Binary      Hex
0   = 00000000  = 00
15  = 00001111  = 0F
16  = 00010000  = 10
255 = 11111111  = FF

Where you see hex:
  Colors in CSS: #FF0000 = Red (FF=255 red, 00=0 green, 00=0 blue)
                 #FFFFFF = White
                 #000000 = Black
                 #FF5733 = Orange-ish
  
  Web addresses: %20 = space (hex 20 = decimal 32 = space character)
  
  Memory addresses: 0x7FFF800A
  
  Hex numbers are often prefixed with "0x": 0x1A = 26 in decimal
```

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Binary | Base-2 number system using only 0 and 1 |
| Bit | Single binary digit (0 or 1) |
| Byte | 8 bits grouped together (can hold values 0–255) |
| ASCII | Standard mapping of characters to numbers |
| Unicode | Universal character encoding covering all languages |
| RGB | Red-Green-Blue color model; each 0–255, stored as 3 bytes |
| Hexadecimal | Base-16 shorthand for binary (0–9, A–F) |

**Binary is the foundation of all computing. Now let's build on it: what does it mean to "program" a computer?**
