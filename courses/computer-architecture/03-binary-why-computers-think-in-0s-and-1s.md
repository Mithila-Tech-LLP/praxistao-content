# Chapter 03: Binary — Why Computers Think in 0s and 1s

You have probably heard that computers "use binary" or "speak in 0s and 1s." But why? Why not use the decimal system we use every day? And how exactly do 0s and 1s represent not just numbers, but text, images, audio, and every piece of information a computer works with?

## Table of Contents

1. [Why Binary?](#1-why-binary)
2. [Counting in Binary](#2-counting-in-binary)
3. [Converting Between Binary and Decimal](#3-converting-between-binary-and-decimal)
4. [Bits, Bytes, and Beyond](#4-bits-bytes-and-beyond)
5. [Hexadecimal — Binary's Shorthand](#5-hexadecimal--binarys-shorthand)
6. [Representing Negative Numbers](#6-representing-negative-numbers)
7. [Representing Text — ASCII and Unicode](#7-representing-text)
8. [Representing Images, Audio, and Everything Else](#8-representing-everything-else)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Binary?

Imagine you are designing a language for a machine that uses electricity. You decide to represent numbers using voltage levels. In decimal, you would need 10 distinct voltage levels: 0V for "0", 0.5V for "1", 1.0V for "2", ..., 4.5V for "9". 

The problem: noise. Real electrical circuits have noise — tiny random fluctuations in voltage. If your signal is supposed to be 1.0V (representing "2") but drifts to 0.95V due to interference, did the machine just read "1.9" instead of "2"? How do you tell the difference?

With only **two** levels — say, below 0.4V = "0" and above 0.7V = "1" — noise becomes much less dangerous. As long as the noise does not push a signal across the 0.4V–0.7V boundary, the machine reads the correct value. You can tolerate much larger noise margins.

Additionally, transistors are naturally binary devices. A MOSFET is either fully on (conducting → "1") or fully off (not conducting → "0"). Using a transistor in its "partially on" state (which you would need for multiple voltage levels) requires careful analog design that is harder to make reliable, fast, and miniaturizable.

**Binary is the natural language of transistors.** This is not a design choice — it is a consequence of how transistors work.

### Binary in Nature

Binary is also common in nature:
- A light switch: on or off
- A door: open or closed
- A coin flip: heads or tails
- A neuron: firing or not firing

The universe seems fond of binary.

### Quick Check

> 1. Why is binary more noise-tolerant than a 10-level (decimal) voltage system?
> 2. A transistor is naturally binary. What does it mean for a transistor to be in state "1"?
> 3. Name three everyday objects that have binary (two-state) behavior.

---

## 2. Counting in Binary

You already know how to count in decimal (base 10). You use 10 digits: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9. When you run out of digits, you put a 1 in the next column and start over: 9 → 10, 19 → 20, 99 → 100.

Binary (base 2) works exactly the same way, but you only have 2 digits: 0 and 1. When you run out, you carry into the next column immediately:

```
Decimal:  0   1   2   3   4   5   6   7   8   9   10  11  12  13  14  15
Binary:   0   1  10  11 100 101 110 111 1000 1001 1010 1011 1100 1101 1110 1111
```

### Place Values

In decimal, each column represents a power of 10:
```
  Thousands  Hundreds  Tens   Ones
  (10³=1000) (10²=100) (10¹=10) (10⁰=1)
       1        2        3       4     ←  This represents 1×1000 + 2×100 + 3×10 + 4×1 = 1234
```

In binary, each column represents a power of 2:
```
  128s  64s  32s  16s  8s  4s  2s  1s
  (2⁷)  (2⁶) (2⁵) (2⁴)(2³)(2²)(2¹)(2⁰)
   1     0    1    0   0   1   1   0  ←  This is the binary number 10100110
```

The rightmost bit is called the **LSB** (Least Significant Bit) because it has the smallest value (1). The leftmost bit is the **MSB** (Most Significant Bit) because it has the largest value.

### Counting Practice

```
0000 = 0     0001 = 1     0010 = 2     0011 = 3
0100 = 4     0101 = 5     0110 = 6     0111 = 7
1000 = 8     1001 = 9     1010 = 10    1011 = 11
1100 = 12    1101 = 13    1110 = 14    1111 = 15
```

Notice the pattern: the rightmost bit alternates every 1 number, the next bit every 2, the next every 4, the next every 8. This is how binary counting works.

### Quick Check

> 1. What is the decimal value of binary 1010?
> 2. What is the binary representation of decimal 13?
> 3. If you have 4 bits (a 4-bit number), what is the largest number you can represent?

---

## 3. Converting Between Binary and Decimal

### Binary to Decimal

To convert binary to decimal, multiply each bit by its place value and add:

**Example: Convert 11010101₂ to decimal**

```
Position:  7    6    5    4    3    2    1    0
Bit:       1    1    0    1    0    1    0    1
Value:   128   64    0   16    0    4    0    1

Sum: 128 + 64 + 0 + 16 + 0 + 4 + 0 + 1 = 213
```

So 11010101₂ = 213₁₀

### Decimal to Binary

Method: repeatedly divide by 2 and record the remainders. Read the remainders bottom-to-top.

**Example: Convert 213 to binary**

```
213 ÷ 2 = 106 remainder 1  ← LSB
106 ÷ 2 =  53 remainder 0
 53 ÷ 2 =  26 remainder 1
 26 ÷ 2 =  13 remainder 0
 13 ÷ 2 =   6 remainder 1
  6 ÷ 2 =   3 remainder 0
  3 ÷ 2 =   1 remainder 1
  1 ÷ 2 =   0 remainder 1  ← MSB

Read bottom-to-top: 11010101
```

So 213₁₀ = 11010101₂ ✓

### The Shortcut: Powers of 2

Memorizing powers of 2 makes conversion much faster:

```
2⁰ = 1       2⁴ = 16     2⁸ = 256
2¹ = 2       2⁵ = 32     2⁹ = 512
2² = 4       2⁶ = 64     2¹⁰ = 1024 (≈ 1 thousand)
2³ = 8       2⁷ = 128    2²⁰ = 1,048,576 (≈ 1 million)
```

These numbers — 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024... — appear constantly in computer science. If you see these numbers in a spec sheet or code, you know you are looking at something related to binary.

### Quick Check

> 1. Convert 10110011₂ to decimal. Show your work.
> 2. Convert 179 to binary using repeated division. Show your work.
> 3. Why do computer memory sizes come in values like 4GB, 8GB, 16GB, 32GB (all powers of 2) rather than 5GB, 10GB, 15GB?

---

## 4. Bits, Bytes, and Beyond

### The Bit

A **bit** (binary digit) is the smallest unit of information: a single 0 or 1. With 1 bit you can represent 2 states. With 2 bits you can represent 4 states (00, 01, 10, 11). With n bits you can represent **2ⁿ states**.

```
Bits    States         Examples
  1       2            on/off
  2       4            four colors
  4      16            hexadecimal digit, BCD digit
  8     256            one ASCII character
 16    65,536          65K different colors (16-bit color)
 32    4.3 billion     IPv4 address space, 32-bit integer
 64    18 quintillion  modern CPU word size
128    3.4 × 10³⁸      AES-128 encryption key space
256    1.2 × 10⁷⁷      enough to address every atom in the universe
```

### The Byte

A **byte** is 8 bits. This is the standard unit of computer data. Why 8? Historical reasons — early computers used various sizes (6-bit, 7-bit), but 8 settled as the standard because it could:
- Hold one ASCII character
- Be efficiently processed by hardware
- Divide evenly into two 4-bit nibbles (useful for hexadecimal)

The word "byte" was coined by Werner Buchholz in 1956 at IBM.

### Storage Sizes

```
Unit       Abbreviation    Exact value             Approximate
──────────────────────────────────────────────────────────────
Bit        b               1 bit                   —
Byte       B               8 bits                  —
Kilobyte   KB              1,024 bytes             10³ bytes
Megabyte   MB              1,048,576 bytes         10⁶ bytes
Gigabyte   GB              1,073,741,824 bytes     10⁹ bytes
Terabyte   TB              ~1.1 × 10¹² bytes       10¹² bytes
Petabyte   PB              ~1.1 × 10¹⁵ bytes       10¹⁵ bytes
Exabyte    EB              ~1.1 × 10¹⁸ bytes       10¹⁸ bytes
```

Note: Hard drive manufacturers use 1 KB = 1000 bytes (decimal). Memory manufacturers use 1 KB = 1024 bytes (binary). This is why a "1TB" drive shows up as ~931 GB in your operating system. To avoid confusion, the binary prefixes use "kibibyte" (KiB = 1024 bytes), "mebibyte" (MiB), "gibibyte" (GiB), etc. — though these are rarely used in everyday speech.

### Data Transfer Rates

Data transfer rates are measured in bits per second (bps):
- A USB 3.2 connection transfers at 20 Gbps (gigabits per second)
- A 5G connection can peak at ~4 Gbps
- LPDDR5X RAM in a modern phone transfers at ~88 GB/s (gigabytes per second, i.e., ~700 Gbps)

Note: lowercase b = bits, uppercase B = bytes. 1 Gbps = 125 MB/s.

### Quick Check

> 1. How many different values can a 16-bit number represent?
> 2. You download a 4GB file over a 100 Mbps internet connection. How many seconds will it take? (1 byte = 8 bits)
> 3. A RAM chip is 64-bit wide. How many bytes does it transfer per clock cycle?

---

## 5. Hexadecimal — Binary's Shorthand

Binary is great for computers but painful for humans. An 8-bit number like 10110101 is hard to read at a glance. Hexadecimal (base 16, "hex") provides a convenient shorthand.

Hex uses 16 symbols: 0-9 (decimal values 0-9) and A-F (decimal values 10-15):

```
Hex: 0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
Dec: 0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
Bin: 0  1 10 11 100 101 110 111 1000 1001 1010 1011 1100 1101 1110 1111
```

The key insight: **one hex digit = exactly 4 binary bits**. This makes conversion trivial:

```
Binary:   1011   0101
Hex:       B      5    →  0xB5

Binary:   1101   1110   1010   1101
Hex:       D      E      A      D    →  0xDEAD
```

To convert binary to hex, group the bits in fours from right to left, then convert each group.

### Hex Notation

Hex numbers are written with a "0x" prefix (0xDEAD), or with an "h" suffix (DEADh), or with a subscript ₁₆ (DEAD₁₆).

### Where You See Hex

Hex is everywhere in computing:
- **Memory addresses**: 0x7FFE45B30 — the hex address of a variable in RAM
- **Colors in web design**: #FF5733 (red=FF=255, green=57=87, blue=33=51)
- **Machine code**: the raw bytes of a program displayed in a hex editor
- **MAC addresses**: 00:1A:2B:3C:4D:5E — a network interface's identifier
- **SHA-256 hash**: a64f2fe4... — 64 hex digits (256 bits)

### Quick Check

> 1. Convert 0xAF to binary. Show your work.
> 2. Convert binary 11001010 11110000 to hex.
> 3. The web color #1E90FF is "dodger blue." What are the decimal values of the R, G, and B components?

---

## 6. Representing Negative Numbers

So far all our binary numbers are positive. But computers must handle negative numbers too. How?

### Sign-Magnitude

The obvious approach: use the MSB as a sign bit. 0 = positive, 1 = negative.
- 0101 = +5
- 1101 = -5

Problem: **two representations of zero** (0000 = +0 and 1000 = -0). Also, addition requires different logic for positive and negative numbers — inconvenient for hardware.

### Two's Complement

The standard solution used in virtually every modern computer. To negate a number:
1. Invert all the bits (flip every 0 to 1 and every 1 to 0) — this is called the one's complement
2. Add 1

**Example: represent -5 in 8-bit two's complement**

```
Start:    00000101  (5 in binary)
Invert:   11111010  (flip all bits)
Add 1:    11111011  (this is -5 in two's complement)
```

Verify: 11111011 + 00000101 = 100000000 → only 8 bits kept → 00000000 = 0 ✓

The beautiful property of two's complement: **ordinary binary addition works for both positive and negative numbers**. No special-case logic needed.

```
  5 + (-3) in 8-bit two's complement:
  
    00000101  (+5)
  + 11111101  (-3, which is 256-3=253 in unsigned)
  ──────────
   100000010  → discard carry bit → 00000010 = 2 ✓
```

### Range of Two's Complement

For an n-bit two's complement number:
- Minimum value: -2^(n-1)
- Maximum value: 2^(n-1) - 1
- Example for 8 bits: -128 to +127
- Example for 32 bits: -2,147,483,648 to +2,147,483,647

```
8-bit two's complement range:
10000000 = -128   (the only number where the sign bit means -128)
10000001 = -127
...
11111111 = -1
00000000 = 0
00000001 = 1
...
01111111 = 127
```

### Integer Overflow

What happens if you add 127 + 1 in a signed 8-bit system? 01111111 + 00000001 = 10000000 = -128. You have **overflowed** — the result wrapped around from the maximum positive to the minimum negative. This is a real programming bug that can cause serious problems. (The famous 1996 Ariane 5 rocket explosion was caused by an integer overflow.)

### Quick Check

> 1. Convert -42 to 8-bit two's complement. Show your steps.
> 2. Add the 8-bit two's complement numbers 11110000 and 00010000. What is the result in decimal?
> 3. In a 16-bit signed integer, what is the maximum positive value and the minimum negative value?

---

## 7. Representing Text

### ASCII

The American Standard Code for Information Interchange (ASCII) was developed in the 1960s to represent text. It assigns a 7-bit number (0-127) to each printable character and control code:

```
65 = 'A'    66 = 'B'    67 = 'C'    ...    90 = 'Z'
97 = 'a'    98 = 'b'    99 = 'c'    ...    122 = 'z'
48 = '0'    49 = '1'    50 = '2'    ...    57 = '9'
32 = ' '    33 = '!'    46 = '.'    ...
```

The word "Hello" in ASCII: 72 101 108 108 111 → in hex: 48 65 6C 6C 6F

Notice: uppercase 'A' is 65 and lowercase 'a' is 97. The difference is exactly 32. This means you can convert case by flipping bit 5 (the 32-place bit). Simple bit manipulation!

### The Problem with ASCII

ASCII only covers 128 characters — enough for English but not for Hindi (देवनागरी), Chinese (漢字), Arabic (عربي), or any of the thousands of scripts used around the world.

### Unicode and UTF-8

Unicode is the modern standard. It assigns a unique "code point" to every character in every language — over 143,000 characters today, including emoji 🎉.

**UTF-8** is the most common encoding of Unicode. It uses variable-length encoding:
- 1 byte for ASCII characters (backward compatible)
- 2-4 bytes for other Unicode characters

English "A" → 0x41 (1 byte)
Hindi "अ" → 0xE0 0xA4 0x85 (3 bytes)
Emoji "🎉" → 0xF0 0x9F 0x8E 0x89 (4 bytes)

Over 97% of all web pages use UTF-8. When you see garbled text (□□□ or mojibake), it is usually a UTF-8 text being interpreted as the wrong encoding.

### Quick Check

> 1. What is the ASCII code for the character '0' (zero)? And for '9'? What is the difference?
> 2. The text "Hi!" in ASCII. Convert each character to hex.
> 3. Why can't ASCII represent the character 'न' (the Hindi letter Na)?

---

## 8. Representing Everything Else

Binary can represent any information. The question is always: how do you encode it?

### Integers: Already covered (two's complement)

### Floating-Point Numbers

Real numbers (3.14, -0.001, 6.022 × 10²³) are stored using the IEEE 754 floating-point standard. A 32-bit float has three parts:
- 1 bit: sign (positive or negative)
- 8 bits: exponent (where the decimal point is)
- 23 bits: mantissa (the significant digits)

This represents numbers from ±1.4×10⁻⁴⁵ to ±3.4×10³⁸ — enormous range, at the cost of not being able to represent every decimal value exactly. (0.1 in float is not exactly 0.1 — this surprises every new programmer.)

### Images

A digital image is a grid of pixels. Each pixel has a color. In 24-bit color, each pixel uses 3 bytes: 1 byte each for Red, Green, and Blue intensity (0-255). A 1920×1080 image has 1,920 × 1,080 = 2,073,600 pixels. At 3 bytes per pixel, that is 6.2 MB uncompressed. JPEG compression reduces this by encoding the pattern of colors rather than each pixel individually.

### Audio

Audio is captured by sampling the air pressure (sound wave) at regular intervals. CD quality uses 44,100 samples per second (44.1 kHz), with each sample stored as a 16-bit signed integer. Stereo means two channels. One second of CD audio = 44,100 × 2 bytes × 2 channels = 176,400 bytes = 172 KB. MP3 compression reduces this by discarding sounds the human ear is least sensitive to.

### Video

Video is a sequence of images displayed rapidly (24-120 frames per second). Uncompressed 4K video (3840×2160) at 30fps at 24-bit color = 745 MB/s. H.264 and H.265 video codecs compress this dramatically by encoding only the differences between frames.

### Everything Else

- Boolean (true/false): 1 bit (or often 1 byte for alignment reasons)
- Date/time: a 64-bit integer counting seconds since January 1, 1970 (Unix timestamp)
- Geographic coordinates: two 64-bit floating-point numbers (latitude, longitude)
- A machine learning model weight: a 32-bit or 16-bit floating-point number

**Everything your computer knows about the world is bits. The encoding — the agreed-upon mapping between bit patterns and meanings — is what gives those bits their meaning.**

### Quick Check

> 1. A 640×480 image uses 24-bit color. What is the uncompressed file size in KB?
> 2. The number 0.1 cannot be represented exactly in IEEE 754 binary floating-point. What does this imply for computing financial amounts (like prices in dollars)?
> 3. Why does the JPEG format produce smaller files than the BMP format (which stores raw pixel values)?

---

## Summary

- Binary (base 2) uses only two digits — 0 and 1. This matches the two states of a transistor and is more noise-tolerant than multi-level systems.
- Binary place values are powers of 2: 1, 2, 4, 8, 16, 32, 64, 128...
- Converting binary↔decimal: multiply each bit by its place value and sum (binary→decimal); repeatedly divide by 2 and record remainders (decimal→binary).
- A **bit** is 0 or 1. A **byte** is 8 bits. Storage sizes: KB, MB, GB, TB...
- **Hexadecimal** (base 16) is binary's shorthand: one hex digit = 4 bits.
- **Two's complement** is the standard way to represent negative numbers. Invert all bits, add 1. Ordinary binary addition then works for both positive and negative numbers.
- **ASCII** encodes English characters in 7 bits. **Unicode/UTF-8** encodes every character in every language using 1-4 bytes.
- Every type of data — images, audio, video, text, coordinates — is ultimately encoded as a sequence of bits.

---

## Exercises

### Easy

1. Convert the following to decimal: (a) 1100₂, (b) 10111₂, (c) 11111111₂

2. Convert the following to binary: (a) 25, (b) 100, (c) 255

3. Convert the following hex to binary: (a) 0xA3, (b) 0xFF, (c) 0x1B

4. What is -7 in 4-bit two's complement? Show your steps.

### Medium

5. A university has 50,000 students. Each student record stores: student ID (32-bit integer), name (50 bytes UTF-8), GPA (32-bit float), graduation year (16-bit integer). Calculate the total storage in MB for all student records.

6. Write out the binary representation of the word "Code" in ASCII. Then write it in hex. (Hint: C=67, o=111, d=100, e=101)

7. A 32-bit unsigned integer has range 0 to 4,294,967,295. A 32-bit signed integer (two's complement) has range -2,147,483,648 to +2,147,483,647. Why do they have the same number of total values (4,294,967,296)? What changes between unsigned and signed?

### Hard

8. **Integer overflow bug**: The Ariane 5 rocket exploded in 1996 due to an integer overflow. A value that was known to be below 32,768 in the old Ariane 4 system was converted from a 64-bit float to a 16-bit signed integer. In Ariane 5, the value exceeded 32,767 — the maximum for a 16-bit signed integer. (a) What is 32,767 in binary? (b) What is 32,768 in binary? (c) If you store 32,768 in a 16-bit signed integer, what value do you get? (d) What general lesson does this teach about type safety in systems programming?

9. **Floating-point imprecision**: Open a Python interpreter (type `python3` in a terminal). Type `0.1 + 0.2`. You will see `0.30000000000000004` instead of `0.3`. Research why this happens (hint: 0.1 in binary is a repeating fraction, like 1/3 in decimal). Write an explanation of this phenomenon, and describe how financial software systems deal with it (hint: research "fixed-point arithmetic" and "decimal data types").
