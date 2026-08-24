# Chapter 03: Binary and Number Systems

> "There are 10 kinds of people in the world: those who understand binary, and those who don't."
> — Classic programmer joke (the "10" is 2 in binary)

Every piece of data your computer has ever processed — every email, every video, every line of Astra code — exists as sequences of 1s and 0s deep in the hardware. Before we can build a compiler, we need to understand this fundamental reality. Not because we will work with raw bits every day, but because it explains *why* programming languages have the types they do, *why* there are integer overflow bugs, *why* floating point arithmetic sometimes gives surprising results, and *how* our Astra compiler will represent data in the compiled executable.

By the end of this chapter, you will understand exactly how Astra's types map to binary representations in memory — and that knowledge will directly inform the type system we design and implement.

---

## Table of Contents

1. Why Computers Use Binary
2. Counting in Binary
3. Converting Between Decimal and Binary
4. Binary Arithmetic
5. Hexadecimal: The Programmer's Shorthand
6. Bits, Nibbles, Bytes, and Beyond
7. ASCII and Character Encoding
8. Unicode and UTF-8
9. IEEE 754 Floating Point
10. Two's Complement: Storing Negative Integers
11. How Astra's Types Map to Binary
12. Exercises

---

## 1. Why Computers Use Binary

The answer is physical. A transistor — the fundamental building block of modern computer chips — is essentially a very fast electrical switch. It has two states:

- **On**: current flows through (high voltage, represents **1**)
- **Off**: current does not flow (low voltage, represents **0**)

```
┌─────────────────────────────────────────────────────────────────┐
│                    THE TRANSISTOR                               │
│                                                                 │
│  A transistor is like a light switch:                           │
│                                                                 │
│        Input                                                    │
│          │                                                      │
│          ▼                                                      │
│   ┌─────────────┐                                              │
│   │  Transistor  │  OFF (no current) → represents 0            │
│   │  (Switch)    │  ON  (current flows) → represents 1         │
│   └─────────────┘                                              │
│          │                                                      │
│          ▼                                                      │
│        Output                                                   │
│                                                                 │
│  Modern CPU: ~50 BILLION transistors on a chip the size        │
│  of your fingernail. Each one a tiny switch.                    │
│                                                                 │
│  1 transistor = 1 bit of information (0 or 1)                  │
│  8 transistors combined = 1 byte (256 possible values)          │
└─────────────────────────────────────────────────────────────────┘
```

Engineers could have tried to build circuits with 10 states (representing digits 0–9 for decimal), but it is physically very difficult to reliably distinguish 10 different voltage levels. Distinguishing just 2 (high vs low) is robust against noise, temperature variation, and manufacturing tolerances. Binary won because it maps perfectly to physical reality.

Everything else in this chapter — decimal-to-binary conversion, hexadecimal, floating point — is built on top of this single physical fact: hardware is easiest to build with two states.

---

## 2. Counting in Binary

In decimal (the number system you grew up with), we have 10 digits: 0, 1, 2, 3, 4, 5, 6, 7, 8, 9. When we run out of digits, we add another position to the left and start over. So after 9 comes 10, after 99 comes 100.

In binary, we have only 2 digits: 0 and 1. When we run out (after 1), we add another position and start over. The pattern is the same, just with base 2 instead of base 10.

**Decimal: base 10**
```
Position value:  1000   100    10     1
                  10³   10²   10¹   10⁰
Digit:             0      0     4     2
Value:         0×1000 + 0×100 + 4×10 + 2×1 = 42
```

**Binary: base 2**
```
Position value:   128    64    32    16     8     4     2     1
                   2⁷    2⁶    2⁵    2⁴    2³    2²    2¹    2⁰
Digit:              0     0     1     0     1     0     1     0
Value:           0+0+32+0+8+0+2+0 = 42
```

So **42 in decimal = 00101010 in binary**.

Here is the complete count from 0 to 15:

```
Decimal │ Binary  │ Explanation
────────┼─────────┼──────────────────────────────────
  0     │  0000   │ all zeros
  1     │  0001   │ 2⁰ = 1
  2     │  0010   │ 2¹ = 2
  3     │  0011   │ 2¹ + 2⁰ = 3
  4     │  0100   │ 2² = 4
  5     │  0101   │ 2² + 2⁰ = 5
  6     │  0110   │ 2² + 2¹ = 6
  7     │  0111   │ 2² + 2¹ + 2⁰ = 7
  8     │  1000   │ 2³ = 8
  9     │  1001   │ 2³ + 2⁰ = 9
 10     │  1010   │ 2³ + 2¹ = 10
 11     │  1011   │ 2³ + 2¹ + 2⁰ = 11
 12     │  1100   │ 2³ + 2² = 12
 13     │  1101   │ 2³ + 2² + 2⁰ = 13
 14     │  1110   │ 2³ + 2² + 2¹ = 14
 15     │  1111   │ 2³ + 2² + 2¹ + 2⁰ = 15
```

Notice: 4 bits can represent 2⁴ = 16 values (0–15). 8 bits can represent 2⁸ = 256 values (0–255). 64 bits can represent 2⁶⁴ ≈ 18.4 quintillion values.

---

## 3. Converting Between Decimal and Binary

### Decimal to Binary

The standard method is **repeated division by 2**, reading remainders bottom-to-top.

Example: Convert **156** to binary.

```
156 ÷ 2 = 78  remainder  0   ← Least significant bit (rightmost)
 78 ÷ 2 = 39  remainder  0
 39 ÷ 2 = 19  remainder  1
 19 ÷ 2 =  9  remainder  1
  9 ÷ 2 =  4  remainder  1
  4 ÷ 2 =  2  remainder  0
  2 ÷ 2 =  1  remainder  0
  1 ÷ 2 =  0  remainder  1   ← Most significant bit (leftmost)

Reading remainders bottom-to-top: 10011100

Verification: 128 + 0 + 0 + 16 + 8 + 4 + 0 + 0 = 156  ✓
```

Another example: Convert **255** to binary.

```
255 ÷ 2 = 127 remainder 1
127 ÷ 2 =  63 remainder 1
 63 ÷ 2 =  31 remainder 1
 31 ÷ 2 =  15 remainder 1
 15 ÷ 2 =   7 remainder 1
  7 ÷ 2 =   3 remainder 1
  3 ÷ 2 =   1 remainder 1
  1 ÷ 2 =   0 remainder 1

Result: 11111111 (all eight 1s)
This is the maximum value for an unsigned 8-bit number!
```

### Binary to Decimal

Multiply each bit by its position value (power of 2) and sum.

Example: Convert **10110101** to decimal.

```
Position:  7    6    5    4    3    2    1    0
Bit:       1    0    1    1    0    1    0    1
Value:   128    0   32   16    0    4    0    1

128 + 32 + 16 + 4 + 1 = 181
```

A quick way to verify: if the leftmost bit is 1, the number is at least 2^(n-1). For an 8-bit number with leftmost bit 1, the value is between 128 and 255.

### Powers of 2 Worth Memorizing

```
2⁰  = 1
2¹  = 2
2²  = 4
2³  = 8
2⁴  = 16
2⁵  = 32
2⁶  = 64
2⁷  = 128
2⁸  = 256
2¹⁰ = 1,024              (approximately 1 thousand, hence "kilo")
2²⁰ = 1,048,576          (approximately 1 million, hence "mega")
2³⁰ = 1,073,741,824      (approximately 1 billion, hence "giga")
2³² = 4,294,967,296      (max value + 1 for 32-bit unsigned int)
2⁶⁴ = 18,446,744,073,709,551,616 (max value + 1 for 64-bit)
```

---

## 4. Binary Arithmetic

### Binary Addition

Binary addition follows the same rules as decimal addition, but with carries happening at 2 instead of 10.

Rules:
```
0 + 0 = 0
0 + 1 = 1
1 + 0 = 1
1 + 1 = 0, carry 1     (because 1+1=2, and 2 in binary is "10")
1 + 1 + 1(carry) = 1, carry 1
```

Example: **42 + 27 = 69**

```
   42 =  0010 1010
+  27 =  0001 1011
────────────────────
Carries:  0011 0000
Result:   0100 0101 = 69

Let's verify:
  0010 1010
+ 0001 1011
──────────
  0100 0101

Column by column (right to left):
  0+1=1
  1+1=0, carry 1
  0+0+1=1
  1+1=0, carry 1
  0+1+1=0, carry 1
  1+0+1=0, carry 1
  0+0+1=1
  0+0=0
Result: 01000101 = 64+4+1 = 69  ✓
```

### Integer Overflow

What happens when you add two numbers and the result is too large to fit in the number of bits available?

```
8-bit unsigned integer: max value = 11111111 = 255

  255 =  1111 1111
+   1 =  0000 0001
──────────────────
         0000 0000  with a carry out of the leftmost bit
         
The carry is DISCARDED! Result: 0

255 + 1 = 0 in 8-bit arithmetic!
```

This is **integer overflow** — the bane of many bugs. In C and C++, signed integer overflow is undefined behavior. In many languages, it silently wraps around.

Astra's `int` type is 64-bit, so overflow is extremely rare in practice (it would require numbers larger than 9 quintillion), but we need to be aware of it. Our compiler will generate code that uses 64-bit integer arithmetic by default.

---

## 5. Hexadecimal: The Programmer's Shorthand

Binary is great for hardware, but long binary numbers are tedious to read. Consider:

```
Binary:       1101 1110 1010 1101 1011 1110 1110 1111
Hexadecimal:  D    E    A    D    B    E    E    F
Decimal:      3,735,928,559
```

**Hexadecimal** (base 16) is a compact notation for binary. Each hex digit represents exactly 4 binary bits. Since 2⁴ = 16, we need 16 symbols: 0–9 and A–F.

```
Hex Digit │ Decimal │ Binary
──────────┼─────────┼────────
    0     │    0    │  0000
    1     │    1    │  0001
    2     │    2    │  0010
    3     │    3    │  0011
    4     │    4    │  0100
    5     │    5    │  0101
    6     │    6    │  0110
    7     │    7    │  0111
    8     │    8    │  1000
    9     │    9    │  1001
    A     │   10    │  1010
    B     │   11    │  1011
    C     │   12    │  1100
    D     │   13    │  1101
    E     │   14    │  1110
    F     │   15    │  1111
```

### Converting Binary to Hex

Group bits into 4s from the right, then convert each group:

```
Binary: 1011 0110 1111 0001
Group:  B    6    F    1
Hex: 0xB6F1
```

### Hex Notation

In programming, hexadecimal numbers are prefixed with `0x`:
- `0xFF` = 255
- `0x100` = 256
- `0x1000` = 4096
- `0xDEADBEEF` = 3,735,928,559

You will see hex everywhere in compiler output, debuggers, memory dumps, and executable files. Memory addresses are typically shown in hex (e.g., `0x00400000`).

In Go, you can write hex literals directly:
```go
address := 0x400000  // memory address in hex
mask := 0xFF         // bitmask: 11111111
color := 0xRRGGBB   // 24-bit RGB color
```

### Why Hex Matters for Compilers

When we write our code generator in Go, we will emit bytes that form machine code. Those bytes are much easier to work with in hex. For example, the x86-64 instruction to move a value into register `rax` starts with the byte `0x48 0xB8`. If we see `0x48 0xB8` in a disassembled binary, we immediately know what it is. If we saw `72 184`, we would have to convert first.

---

## 6. Bits, Nibbles, Bytes, and Beyond

These terms come up constantly in computer science and compiler design:

```
┌────────────────────────────────────────────────────────────────┐
│                    DATA SIZE TERMINOLOGY                       │
│                                                                │
│  1 bit      = single binary digit (0 or 1)                    │
│               Values: 0 or 1                                   │
│               Used for: boolean flags, individual switches     │
│                                                                │
│  4 bits     = 1 nibble                                         │
│               Values: 0 to 15 (0x0 to 0xF)                    │
│               Used for: hex notation, BCD, color channels      │
│                                                                │
│  8 bits     = 1 byte                                           │
│               Values: 0 to 255 (0x00 to 0xFF)                 │
│               Used for: characters (ASCII), smallest addressable unit │
│                                                                │
│  16 bits    = 2 bytes = 1 word (on 16-bit systems)            │
│               Values: 0 to 65,535                              │
│               Used for: Unicode code points (BMP), port numbers│
│                                                                │
│  32 bits    = 4 bytes = 1 double word (dword)                 │
│               Values: 0 to ~4.3 billion (unsigned)            │
│               Used for: 32-bit int, IPv4 addresses, Unicode    │
│                                                                │
│  64 bits    = 8 bytes = 1 quad word (qword)                   │
│               Values: 0 to ~18.4 quintillion (unsigned)       │
│               Used for: 64-bit int, memory addresses, float64  │
│                                                                │
│  128 bits   = 16 bytes = 1 double quad word (SIMD)            │
│               Used for: SIMD vector operations                 │
└────────────────────────────────────────────────────────────────┘
```

### Storage Size Prefixes

```
1 Byte     =           1 byte
1 Kilobyte = 1,024     bytes   (2¹⁰)
1 Megabyte = 1,048,576 bytes   (2²⁰)
1 Gigabyte ≈ 1 billion bytes   (2³⁰)
1 Terabyte ≈ 1 trillion bytes  (2⁴⁰)

Note: Storage manufacturers use decimal (1 KB = 1000 bytes)
      Operating systems use binary (1 KB = 1024 bytes)
      This causes the "my 1TB drive only shows 931 GB" confusion
```

---

## 7. ASCII and Character Encoding

Computers store everything as numbers. Text is no exception. We need a standard mapping from characters to numbers. The original standard was **ASCII** (American Standard Code for Information Interchange), created in 1963.

ASCII uses 7 bits (values 0–127) to encode:
- Control characters (0–31): newline, tab, backspace, etc.
- Printable characters (32–126): space, letters, digits, punctuation
- Delete (127)

```
ASCII Table (selected):

Dec │ Hex │ Char │    Dec │ Hex │ Char │    Dec │ Hex │ Char
────┼─────┼──────┤    ────┼─────┼──────┤    ────┼─────┼──────
 32 │ 20  │ (sp) │     65 │ 41  │  A   │     97 │ 61  │  a
 33 │ 21  │  !   │     66 │ 42  │  B   │     98 │ 62  │  b
 48 │ 30  │  0   │     67 │ 43  │  C   │     99 │ 63  │  c
 49 │ 31  │  1   │     ...       ...   │    ...       ...
 57 │ 39  │  9   │     90 │ 5A  │  Z   │    122 │ 7A  │  z
 10 │ 0A  │ \n   │
  9 │ 09  │ \t   │
```

Important patterns to note:
- Digits '0'–'9' are ASCII 48–57 (0x30–0x39). So `'5' - '0' = 5` converts character to number.
- Uppercase 'A'–'Z' are 65–90 (0x41–0x5A)
- Lowercase 'a'–'z' are 97–122 (0x61–0x7A)
- The difference between uppercase and lowercase is exactly 32 (0x20). So `'a' - 'A' = 32`.

When our Astra lexer reads source code, it is reading a sequence of bytes from disk. When it sees the byte `0x66` (`f`), `0x6E` (`n`), `0x20` (space), it knows it is looking at the keyword `fn` followed by a space.

In Go, we will write code like:

```go
// In our Astra lexer
func isLetter(ch byte) bool {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
    return ch >= '0' && ch <= '9'
}

// This works because Go character literals use ASCII values
// 'a' is 97, 'z' is 122, '0' is 48, '9' is 57
```

---

## 8. Unicode and UTF-8

ASCII only handles 128 characters — fine for English, terrible for the rest of the world. **Unicode** is the universal standard that assigns a unique number (called a **code point**) to every character in every writing system: Chinese, Arabic, emoji, ancient Egyptian hieroglyphics — all of it.

Unicode defines over 140,000 characters (and growing). Code points are written as `U+XXXX` where X is a hex digit.

```
Some Unicode Code Points:

U+0041  'A'         (Latin A — same as ASCII)
U+00E9  'é'         (Latin e with acute accent)
U+4E2D  '中'        (Chinese character for "middle/China")
U+0041  'α'         (Actually U+03B1, Greek alpha)
U+1F600 '😀'       (Grinning face emoji)
U+1F4BB '💻'       (Laptop computer emoji)
```

### UTF-8: The Encoding

Having code points is great, but how do we actually store them as bytes? **UTF-8** is the dominant answer and it is clever: it is variable-length.

```
┌───────────────────────────────────────────────────────────────┐
│                    UTF-8 ENCODING RULES                       │
│                                                               │
│  Code Point Range        │ UTF-8 Byte Pattern                 │
│  ────────────────────────┼─────────────────────────────────── │
│  U+0000   to U+007F      │ 0xxxxxxx                           │
│  (0–127, ASCII)          │ 1 byte                             │
│                          │                                    │
│  U+0080   to U+07FF      │ 110xxxxx 10xxxxxx                  │
│  (128–2047)              │ 2 bytes                            │
│                          │                                    │
│  U+0800   to U+FFFF      │ 1110xxxx 10xxxxxx 10xxxxxx        │
│  (2048–65535, most langs)│ 3 bytes                            │
│                          │                                    │
│  U+10000  to U+10FFFF    │ 11110xxx 10xxxxxx 10xxxxxx 10xx   │
│  (65536+, emoji etc)     │ 4 bytes                            │
└───────────────────────────────────────────────────────────────┘
```

This is brilliant because:
1. **ASCII compatibility**: all ASCII characters encode as a single byte (their ASCII value). An ASCII file is already valid UTF-8.
2. **Self-synchronizing**: if you are in the middle of a multi-byte sequence, you can tell — continuation bytes always start with `10`.
3. **Variable length**: common characters (ASCII, Western European) use 1-2 bytes. Rare characters use 3-4. This saves space compared to fixed 4-byte encoding (UTF-32).

**Example: Encoding '中' (U+4E2D)**

```
U+4E2D in binary: 100 111000 101101
                   (break into groups: 4 + 6 + 6 bits)

UTF-8 3-byte template: 1110xxxx 10xxxxxx 10xxxxxx
Fill in:              1110 0100 10 111000 10 101101
                       E4        B8        AD

So '中' = bytes 0xE4 0xB8 0xAD
Verify: $ python3 -c "print('中'.encode('utf-8'))"
b'\xe4\xb8\xad'  ✓
```

### Astra and Unicode

Astra source files are UTF-8 encoded. Astra string literals contain UTF-8 bytes. Our lexer must handle multi-byte characters correctly.

In Go, `string` is already a sequence of UTF-8 bytes. The `rune` type is a Unicode code point (32-bit integer). We can iterate over a Go string character by character using `range`:

```go
// Correctly iterating over Unicode in Go
s := "Hello, 世界"
for i, r := range s {
    fmt.Printf("index %d: %c (U+%04X)\n", i, r, r)
}
// Output:
// index 0: H (U+0048)
// index 1: e (U+0065)
// ...
// index 7: 世 (U+4E16)
// index 10: 界 (U+754C)
// Note: indices skip by byte count, not character count
```

---

## 9. IEEE 754 Floating Point

How does a computer store `3.14159`? We cannot represent all decimal fractions exactly in binary — just as 1/3 cannot be written exactly in decimal (0.333...), 0.1 cannot be written exactly in binary.

**IEEE 754** is the international standard for floating-point arithmetic. It is used by virtually all modern CPUs and programming languages including Astra.

A 64-bit (double precision) IEEE 754 float has three parts:

```
┌─────────────────────────────────────────────────────────────────┐
│              64-BIT IEEE 754 FLOAT LAYOUT                       │
│                                                                 │
│  Bit 63  │  Bits 62–52  │  Bits 51–0                           │
│  ────────┼──────────────┼──────────────────────────────────    │
│  S (sign)│  E (exponent)│  M (mantissa / fraction)             │
│  1 bit   │  11 bits     │  52 bits                             │
│                                                                 │
│  Value = (-1)^S × 1.M × 2^(E - 1023)                          │
│                                                                 │
│  Example: 3.14159265358979...                                   │
│  S = 0 (positive)                                              │
│  E = 1024 (biased, represents 2^1)                             │
│  M = 1001001000011111101101010... (52 bits of 1.5707...)        │
│                                                                 │
│  Bit layout:                                                    │
│  0 10000000000 1001001000011111101101010100010001000010110100   │
│  │ └─────────┘ └──────────────────────────────────────────┘   │
│  S exponent=1024  mantissa (52 bits of fraction part)           │
└─────────────────────────────────────────────────────────────────┘
```

### The Floating Point Trap: Imprecision

Because not all decimal fractions have exact binary representations, floating point arithmetic can produce surprising results:

```go
// In Go (same happens in any IEEE 754 language)
package main
import "fmt"

func main() {
    a := 0.1
    b := 0.2
    fmt.Println(a + b)         // prints: 0.30000000000000004 !
    fmt.Println(a + b == 0.3)  // prints: false !
}
```

This is not a bug — it is the fundamental limitation of representing real numbers in finite binary digits. For Astra, this means:

1. `float` values in Astra should never be compared with `==` for exact equality (use a small epsilon margin)
2. Financial calculations should use integer arithmetic (store cents, not dollars)
3. Our compiler should generate IEEE 754 instructions (x86-64 uses `ADDSD`, `MULSD`, `DIVSD` for 64-bit float operations)

### Special Float Values

IEEE 754 reserves certain bit patterns for special values:

| Value | Meaning |
|-------|---------|
| `+Infinity` | Division by zero: `1.0 / 0.0` |
| `-Infinity` | Negative division by zero: `-1.0 / 0.0` |
| `NaN` | Not a Number: `0.0 / 0.0`, `sqrt(-1)` |
| `+0.0` and `-0.0` | Positive and negative zero (they compare equal) |

Astra's standard library will need to handle these cases when implementing math functions.

---

## 10. Two's Complement: Storing Negative Integers

Unsigned integers are simple — they represent values 0 to 2^n - 1. But how do we represent negative numbers? We need a scheme where the hardware can perform addition and subtraction using the same circuits for both positive and negative numbers.

The answer is **two's complement**. It is the standard for signed integer representation in virtually all modern hardware.

### The Idea

For an n-bit signed integer:
- The leftmost bit is the **sign bit**: 0 = positive, 1 = negative
- Positive numbers are the same as their binary representation
- Negative numbers are represented as 2^n minus the absolute value

For 8-bit signed integers (range: -128 to +127):

```
Decimal │ Binary   │ Explanation
────────┼──────────┼──────────────────────────────────────
   127  │ 01111111 │ Max positive (sign bit 0)
     1  │ 00000001 │
     0  │ 00000000 │
    -1  │ 11111111 │ 256 - 1 = 255 = 11111111
    -2  │ 11111110 │ 256 - 2 = 254 = 11111110
  -127  │ 10000001 │ 256 - 127 = 129 = 10000001
  -128  │ 10000000 │ Min negative (sign bit 1, no positive twin)
```

### How to Compute Two's Complement

To negate a number in two's complement:
1. Flip all bits (bitwise NOT)
2. Add 1

Example: Negate 5 (00000101) to get -5:
```
Step 1: Flip bits:  00000101 → 11111010
Step 2: Add 1:      11111010 + 1 = 11111011

-5 in 8-bit two's complement = 11111011

Verify: 11111011 = 128+64+32+16+8+0+2+1 = 251
As unsigned: 251. As signed: 251 - 256 = -5  ✓
```

### Why Two's Complement Is Brilliant

The hardware can add any two 8-bit values using the same adder circuit, regardless of whether they are positive or negative:

```
  5 + (-3):
  
  5  = 00000101
 -3  = 11111101   (256-3 = 253)
 ──────────────
Sum:  100000010  → discard the overflow bit → 00000010 = 2  ✓

5 + (-3) = 2. Correct!
```

No special "negative number adder" needed — the same hardware works for everything.

### Signed vs Unsigned in Astra

Astra's `int` type is a 64-bit signed integer (two's complement):
- Range: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- Minimum: `0x8000000000000000` (1 followed by 63 zeros)
- Maximum: `0x7FFFFFFFFFFFFFFF` (0 followed by 63 ones)

This is the same as Go's `int64` type. For most programs, this range is more than sufficient.

---

## 11. How Astra's Types Map to Binary

Now we can define precisely how each Astra primitive type is stored in memory:

```
┌─────────────────────────────────────────────────────────────────┐
│                ASTRA PRIMITIVE TYPES — BINARY LAYOUT            │
│                                                                 │
│  Type    │ Bits │ Encoding            │ Range / Notes           │
│  ────────┼──────┼─────────────────────┼──────────────────────── │
│  int     │  64  │ Two's complement    │ -9.2×10¹⁸ to 9.2×10¹⁸  │
│          │      │ (signed integer)    │ Same as int64 in Go     │
│  ────────┼──────┼─────────────────────┼──────────────────────── │
│  float   │  64  │ IEEE 754 double     │ ±1.7×10³⁰⁸, ~15 digits │
│          │      │ precision           │ Same as float64 in Go   │
│  ────────┼──────┼─────────────────────┼──────────────────────── │
│  bool    │   1  │ 0 = false, 1 = true │ (stored as 1 byte in    │
│          │ (1B) │                     │  practice for alignment)│
│  ────────┼──────┼─────────────────────┼──────────────────────── │
│  char    │  32  │ Unicode code point  │ U+0000 to U+10FFFF     │
│          │      │ (UTF-32)            │ Same as int32/rune in Go│
│  ────────┼──────┼─────────────────────┼──────────────────────── │
│  string  │ 128  │ Pointer + Length    │ ptr: 64-bit address     │
│          │ (16B)│                     │ len: 64-bit byte count  │
│          │      │                     │ Data: UTF-8 in heap     │
└─────────────────────────────────────────────────────────────────┘
```

### The String Type in Detail

Astra strings are not a fixed-size primitive — they can be arbitrarily long. Instead, a string *value* is a small fixed-size structure containing:
1. A **pointer** (64 bits) to the actual bytes of the string in memory (on the heap)
2. A **length** (64 bits) counting the number of bytes

```
┌─────────────────────────────────────────────────────────────────┐
│                 STRING IN MEMORY                                 │
│                                                                  │
│  Variable: name = "Aditya"                                       │
│                                                                  │
│  Stack:                       Heap:                              │
│  ┌────────────────────┐       ┌─────────────────────────────┐  │
│  │ name               │       │ Address: 0x00C0001A8000      │  │
│  │ ┌────────────────┐ │       │ 0x41 0x64 0x69 0x74 0x79 0x61 │ │
│  │ │ ptr: 0x00C0001A│─┼──────►│  A    d    i    t    y    a  │  │
│  │ │ len: 6         │ │       └─────────────────────────────┘  │
│  │ └────────────────┘ │                                         │
│  └────────────────────┘                                         │
│                                                                  │
│  The string header (ptr + len) lives on the stack or as a field │
│  The actual bytes live on the heap                              │
└─────────────────────────────────────────────────────────────────┘
```

This is exactly how Go represents strings internally, and it is a great design:
- Passing a string to a function only copies 16 bytes (the header), not the entire string data
- `len(s)` is O(1) — we just read the `len` field, no counting needed
- Strings can be arbitrarily large without changing the size of the header

### Type Sizes Matter for Alignment

CPUs are fastest when data is **aligned** — meaning its address is a multiple of its size. An 8-byte (64-bit) integer should be at an address divisible by 8. A 4-byte (32-bit) value should be at an address divisible by 4.

When our code generator lays out struct fields in memory, it must insert **padding** to maintain alignment:

```astra
struct Point {
    x: float  // 8 bytes at offset 0
    y: float  // 8 bytes at offset 8
}
// Total size: 16 bytes, no padding needed

struct Mixed {
    flag: bool  // 1 byte at offset 0
    // 7 bytes PADDING to align 'value' to 8-byte boundary
    value: int  // 8 bytes at offset 8
}
// Total size: 16 bytes (7 bytes wasted on padding!)
```

This is a real concern in language design. Astra's compiler will handle alignment automatically, but language designers should understand why it exists.

---

## Astra Build Milestone: Defining Astra's Type Representations

At this point in the guide, we have defined Astra's primitive types in terms of their binary representations. This is a critical design document that all future chapters will reference. Here it is formally:

### Astra Primitive Type Reference

```
╔══════════════════════════════════════════════════════════════════╗
║           ASTRA LANGUAGE: TYPE → BINARY SPECIFICATION           ║
╠══════════════╦═══════╦════════════════╦════════════════════════╣
║  Astra Type  ║ Width ║ Encoding       ║ Example Values          ║
╠══════════════╬═══════╬════════════════╬════════════════════════╣
║ int          ║ 64b   ║ Two's compl.   ║ 0, 42, -17, 2^63-1    ║
╠══════════════╬═══════╬════════════════╬════════════════════════╣
║ float        ║ 64b   ║ IEEE 754 dbl   ║ 0.0, 3.14, -2.7, Inf  ║
╠══════════════╬═══════╬════════════════╬════════════════════════╣
║ bool         ║ 8b*   ║ 0=false,1=true ║ true, false            ║
╠══════════════╬═══════╬════════════════╬════════════════════════╣
║ char         ║ 32b   ║ Unicode point  ║ 'A', 'é', '中', '😀'  ║
╠══════════════╬═══════╬════════════════╬════════════════════════╣
║ string       ║ 128b  ║ ptr(64)+len(64)║ "hello", "", "Aditya"  ║
║              ║       ║ + heap data    ║                         ║
╚══════════════╩═══════╩════════════════╩════════════════════════╝
* bool is stored as 1 byte (8 bits) for memory alignment purposes,
  but only uses 1 bit of information.
```

### Integer Overflow Behavior in Astra

```
Astra int (64-bit signed two's complement):
  Max positive: 9,223,372,036,854,775,807  (0x7FFFFFFFFFFFFFFF)
  Min negative: -9,223,372,036,854,775,808 (0x8000000000000000)
  Overflow wraps around (same behavior as Go, C, Rust in release mode)
```

### Float Precision in Astra

```
Astra float (64-bit IEEE 754 double):
  Significant decimal digits: ~15-16
  Approximate range: ±1.7 × 10^308
  
  WARNING: Floating point is not exact!
  0.1 + 0.2 == 0.3  → false in Astra (and all IEEE 754 languages)
  Use: (0.1 + 0.2 - 0.3).abs() < 0.0001  for approximate equality
```

### String Encoding

```
Astra strings: UTF-8 encoded byte sequences
  len(s): returns byte count, not character count
  s[i]:   returns i-th BYTE, not i-th character
  
  Example:
  let s = "Hello, 世界"
  len(s) == 13   (7 ASCII bytes + 3 bytes for 世 + 3 bytes for 界)
  s[7]   == 0xE4 (first byte of 世's UTF-8 encoding)
```

These definitions are not arbitrary — they were chosen to match exactly how the Go `int64`, `float64`, `uint32` (for char), and `string` types work internally. This means our compiler's code generator can leverage Go's own runtime for many operations.

---

## 12. Exercises

1. **Binary Conversion Practice**: Convert the following decimal numbers to binary: 17, 63, 100, 200, 255. Then convert back from binary to verify. Show all steps of the division method.
   *Hint: Use repeated division by 2 and read remainders bottom-to-top.*

2. **The Hexadecimal Challenge**: Convert the memory address `0x00400000` to decimal. Then convert `255` to hex. Finally, write out what memory address comes right after `0x7FFFFFFE` (the answer requires careful thought about carry).
   *Hint: 0x00400000 = 4 × 16^5. For the last question: 0x7FFFFFFF + 1.*

3. **Float Imprecision**: In Go (or any language with IEEE 754 floats), calculate `0.1 + 0.2`. You will get `0.30000000000000004`. Why does this happen? Write an Astra-style function that correctly checks if two floats are "approximately equal" (within 0.0001). What parameter would you need?
   *Hint: No float value exactly represents 0.1 in binary. Think of it like 1/3 in decimal.*

4. **Two's Complement**: Manually compute the two's complement representation of the following numbers in 8-bit binary: -1, -42, -128. Verify by flipping bits and adding 1. Also determine: what is the 8-bit two's complement of 128? (It cannot be represented — explain why.)
   *Hint: Start from the positive binary representation and apply the negation algorithm.*

5. **String Memory Layout**: An Astra string `s = "programming"` (11 characters, all ASCII). Draw the memory layout showing the stack portion (pointer + length) and the heap portion (the actual bytes). What is the total memory used?
   *Hint: The header is 16 bytes (8 for ptr + 8 for len). The heap data is 11 bytes.*

6. **UTF-8 Encoding**: The emoji 😀 has Unicode code point U+1F600. Using the UTF-8 encoding table from this chapter, determine how many bytes it takes to encode and what those bytes are. Verify by checking: does U+1F600 fall in the 1-byte, 2-byte, 3-byte, or 4-byte range?
   *Hint: U+1F600 = 128512 decimal. It's in the U+10000–U+10FFFF range (4 bytes).*

7. **Struct Alignment Puzzle**: Design an Astra struct that wastes the most padding bytes possible with 4 fields of types: `bool`, `int`, `bool`, `float`. Draw the memory layout, showing offset of each field and any padding. Then rearrange the fields to minimize wasted space.
   *Hint: Arrange large types first, small types last. A bool padded to align an int wastes 7 bytes.*

8. **Overflow Detection**: Write the Astra code that would detect integer overflow before it happens (checking if `a + b` would overflow a 64-bit signed integer). Think about what conditions cause overflow. 
   *Hint: Overflow when adding two positives gives a negative, or adding two negatives gives a positive. You can check: if b > 0 and a > MAX_INT - b, overflow will occur.*

---

## Summary: Key Concepts

| Concept | Definition | Relevance to Astra |
|---------|-----------|-------------------|
| Binary | Base-2 number system; hardware uses it | All Astra data stored as binary |
| Bit | Single binary digit (0 or 1) | Smallest unit of data |
| Byte | 8 bits; smallest addressable memory unit | Memory addresses are byte addresses |
| Hexadecimal | Base-16; compact binary notation | Used in memory addresses, machine code |
| Binary addition | Same as decimal, carry at 2 | CPU's ADD instruction does this |
| Integer overflow | Value wraps when too large for bit width | Astra int can overflow for very large numbers |
| ASCII | 7-bit encoding for English characters | Astra source files, basic strings |
| Unicode | Universal character standard (140k+ chars) | Astra strings support all languages |
| UTF-8 | Variable-length Unicode encoding | Astra's string encoding on disk and in memory |
| IEEE 754 | Standard for floating point numbers | Astra `float` type uses this |
| Float imprecision | 0.1 + 0.2 ≠ 0.3 exactly | Never compare floats with == |
| NaN | Not a Number (0.0/0.0) | Astra math functions can produce this |
| Two's complement | How negative integers are stored | Astra `int` uses this for negatives |
| Sign bit | Leftmost bit indicates sign | 1 = negative in two's complement |
| Memory alignment | Data at addresses divisible by its size | Compiler must pad struct fields |
| String as ptr+len | String header (16 bytes) + heap data | Astra string representation |
| char as Unicode | 32-bit Unicode code point | Astra `char` type holds any character |
