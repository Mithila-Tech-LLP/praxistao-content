# Chapter 12: Strings and Character Encoding — Text in the Digital World

> *"There Ain't No Such Thing As Plain Text. If you have a string, you absolutely must know what encoding it uses. 'Plain text' using 7-bit ASCII is a historical accident."*
> — Joel Spolsky, "The Absolute Minimum Every Software Developer Absolutely, Positively Must Know About Unicode and Character Sets"

Text looks simple on the surface. You type letters; they appear on screen. You store a name in a variable; you print it back. What could be hard about that? As it turns out, almost everything — if you want text to work correctly for all seven billion humans on Earth, across every writing system ever devised.

This chapter builds your understanding from first principles: what a character actually is, why computers had so much trouble representing text historically, how the Unicode standard solved the problem, and how UTF-8 elegantly encodes all of Unicode in a backward-compatible way. We then look at how Go handles strings (with important surprises), and finally design the internal string representation Astra uses — including the C struct that the compiled runtime will use at execution time. By the end, you will understand exactly what happens to every byte when you write `let s = "Hello, 世界"` in an Astra program.

## What We're Building

In Astra, strings are first-class values:

```astra
let greeting = "Hello, World"
let japanese = "こんにちは"
let emoji    = "Astra says 😊"
let mixed    = "Café résumé naïve"
let name     = "Hello, " + "World"
```

The Astra runtime must store these correctly, count characters correctly (not bytes!), and provide efficient operations. Everything here feeds into Chapter 63 (Astra Runtime) and Chapter 64 (Memory Management).

## Table of Contents

1. Why Strings Are Harder Than They Look
2. ASCII — The 1963 Foundation
3. Extended ASCII and Code Pages — The Failed Solution
4. Unicode — One Standard for All Human Writing
5. UTF-8 — The Clever Encoding
6. UTF-16 and UTF-32 — When They're Used
7. Strings in Go — Surprises and Best Practices
8. Strings in Astra — Internal Representation
9. The Astra String Standard Library
10. Astra Build Milestone — AstraString Implementation
11. Exercises
12. Summary

---

## 1. Why Strings Are Harder Than They Look

A string is a sequence of characters. But what is a **character**?

In English, a character is a letter, digit, punctuation mark, or space. Easy. But what about:

- `é` — Is it one character? (It looks like one.) In some encodings it is one byte; in others it is two.
- `😊` — This emoji is a single character to a human, but it occupies **4 bytes** in UTF-8.
- `가` — Korean Hangul syllable. One character, 3 bytes in UTF-8.
- `𝒜` — Mathematical script letter A. One character, 4 bytes in UTF-8.
- `ñ` — Spanish letter eñe. One code point, but can be represented as EITHER `U+00F1` (one code point) or `n` + `U+0303` (combining tilde) — two code points! Both look identical on screen.

If your code assumes `1 character = 1 byte`, it breaks on any of these. This is not a theoretical problem — it is a real bug that has caused data corruption, security vulnerabilities, and crashes in real systems.

The root problem is that computers were designed in America in the 1960s, when English was the only language that mattered. The fix took 30 years of international effort and resulted in Unicode.

---

## 2. ASCII — The 1963 Foundation

**ASCII** (American Standard Code for Information Interchange) was created in 1963. It uses **7 bits**, giving 128 possible values (0 through 127).

```
┌────────────────────────────────────────────────────────────┐
│                    ASCII CHARACTER TABLE                    │
│                   (Printable Characters)                   │
├────────┬──────┬───────┬────────┬──────┬───────┬───────────┤
│ Dec    │ Hex  │ Char  │ Dec    │ Hex  │ Char  │           │
├────────┼──────┼───────┼────────┼──────┼───────┤           │
│  32    │ 0x20 │ SPACE │  64    │ 0x40 │  @    │           │
│  33    │ 0x21 │  !    │  65    │ 0x41 │  A    │           │
│  34    │ 0x22 │  "    │  66    │ 0x42 │  B    │           │
│  35    │ 0x23 │  #    │ ...    │ ...  │ ...   │           │
│  48    │ 0x30 │  0    │  90    │ 0x5A │  Z    │           │
│  49    │ 0x31 │  1    │  97    │ 0x61 │  a    │           │
│  ...   │ ...  │ ...   │ ...    │ ...  │ ...   │           │
│  57    │ 0x39 │  9    │ 122    │ 0x7A │  z    │           │
└────────┴──────┴───────┴────────┴──────┴───────┴───────────┘

Control characters: 0-31 (tab=9, newline=10, carriage return=13)
Maximum printable: 127 (0x7F = DEL)
```

**The binary representation of "Hi":**

```
┌─────────────────────────────────────────────────────────┐
│  String "Hi" in ASCII / UTF-8                           │
│                                                         │
│  'H' = ASCII 72 = 0x48                                  │
│  Binary: 0 1 0 0 1 0 0 0                               │
│          ↑ bit 7 always 0 in ASCII (7-bit standard)    │
│                                                         │
│  'i' = ASCII 105 = 0x69                                 │
│  Binary: 0 1 1 0 1 0 0 1                               │
│                                                         │
│  Memory layout:                                         │
│  ┌────────┬────────┐                                    │
│  │ 0x48   │ 0x69   │                                    │
│  │  'H'   │  'i'   │                                    │
│  └────────┴────────┘                                    │
│   addr 0    addr 1                                      │
└─────────────────────────────────────────────────────────┘
```

ASCII was elegant for its time. Every character fits in one byte (with bit 7 unused). String length in characters equals string length in bytes. Indexing is O(1).

**The problem:** ASCII covers only 128 characters. No accented characters (é, ü, ñ). No Chinese, Japanese, Korean, Arabic, Hebrew, Thai, Hindi, or any of thousands of other writing systems. No emoji (invented 35 years later). No mathematical symbols. For purely English text, ASCII is fine. For the rest of humanity, it fails completely.

---

## 3. Extended ASCII and Code Pages — The Failed Solution

The obvious fix: use the 8th bit. With 8 bits, you get 256 values (0-255). The first 128 are standard ASCII; the upper 128 can be used for extra characters.

The problem: **which** extra characters? In the 1980s, different organizations invented different **code pages**:

- **ISO 8859-1 (Latin-1):** Western European languages (é, ü, ñ, ø)
- **ISO 8859-5:** Cyrillic (А, Б, В, ...)
- **Windows-1252:** Windows Western European (similar to Latin-1 with minor differences)
- **Shift-JIS:** Japanese
- **GB2312:** Simplified Chinese
- **Big5:** Traditional Chinese

This created the **mojibake** problem (文字化け — Japanese for "transformed characters"). If you write a document in Shift-JIS and send it to someone who opens it expecting ISO 8859-1, every character above 127 is interpreted as the wrong character. The document appears as meaningless garbage:

```
Original (Shift-JIS):   日本語のテキスト
Displayed as Latin-1:   æ—¥æœ¬èªžã®ãƒ†ã‚­ã‚¹ãƒˆ
```

By the early 1990s, the internet was connecting computers from every country, and the Babel of incompatible encodings was a disaster. Something had to change.

---

## 4. Unicode — One Standard for All Human Writing

The **Unicode Consortium** was founded in 1991 with one mission: create a single encoding that covers every character in every writing system, past and present.

The core idea: every character is assigned a unique **code point** — a non-negative integer with the prefix `U+` followed by hexadecimal digits.

```
┌────────────────────────────────────────────────────────┐
│              SELECTED UNICODE CODE POINTS               │
├─────────────┬────────┬──────────────────────────────────┤
│  Character  │  U+    │         Description              │
├─────────────┼────────┼──────────────────────────────────┤
│     A       │ U+0041 │ Latin capital letter A           │
│     a       │ U+0061 │ Latin small letter A             │
│     é       │ U+00E9 │ Latin small letter E with acute  │
│     ñ       │ U+00F1 │ Latin small letter N with tilde  │
│     ü       │ U+00FC │ Latin small letter U with diaer. │
│     中      │ U+4E2D │ CJK Unified Ideograph (middle)  │
│     世      │ U+4E16 │ CJK Unified Ideograph (world)   │
│     α       │ U+03B1 │ Greek small letter alpha         │
│     ←       │ U+2190 │ Leftwards arrow                  │
│     ™       │ U+2122 │ Trade mark sign                  │
│     😀      │ U+1F600 │ Grinning face emoji             │
│     😊      │ U+1F60A │ Smiling face with smiling eyes  │
│     🔨      │ U+1F528 │ Hammer                          │
└─────────────┴────────┴──────────────────────────────────┘
```

**Total range:** U+0000 to U+10FFFF — over 1.1 million possible code points.

**Currently assigned:** about 149,000+ characters (as of Unicode 15).

### Unicode Planes

The full Unicode space is divided into **17 planes**, each containing 65,536 code points:

```
┌────────────────────────────────────────────────────────┐
│                   UNICODE PLANES                        │
├────────────┬───────────────────────────────────────────┤
│ Plane 0    │ U+0000 - U+FFFF   (BMP — Basic Multilingual│
│            │                    Plane) — most everyday  │
│            │                    characters here         │
├────────────┼───────────────────────────────────────────┤
│ Plane 1    │ U+10000 - U+1FFFF (Supplementary Multilin- │
│            │                    gual — emoji, historic  │
│            │                    scripts, math symbols)  │
├────────────┼───────────────────────────────────────────┤
│ Planes 2-3 │ Rare CJK ideographs, historic scripts      │
├────────────┼───────────────────────────────────────────┤
│ Planes 4-13│ Unassigned (future use)                    │
├────────────┼───────────────────────────────────────────┤
│ Planes 14  │ Tags, variation selectors                  │
├────────────┼───────────────────────────────────────────┤
│ Planes 15-16│ Private use areas                         │
└────────────┴───────────────────────────────────────────┘
```

Unicode defines **what** characters exist and their properties. It does not define **how** to store them in memory — that is the job of an encoding like UTF-8.

---

## 5. UTF-8 — The Clever Encoding

**UTF-8** (Unicode Transformation Format, 8-bit) was designed by Ken Thompson and Rob Pike (the same people who created Go and Unix). It uses between 1 and 4 bytes to encode each code point.

### The Encoding Rules

```
┌──────────────────────────────────────────────────────────────────┐
│                    UTF-8 ENCODING TABLE                           │
├──────────────────┬────────────────────────────────────────────────┤
│  Code Point Range │  Byte Sequence (x = payload bits)             │
├──────────────────┼────────────────────────────────────────────────┤
│ U+0000 - U+007F  │  0xxxxxxx                      (1 byte)        │
│  (0 - 127)       │  7 payload bits                                │
├──────────────────┼────────────────────────────────────────────────┤
│ U+0080 - U+07FF  │  110xxxxx  10xxxxxx             (2 bytes)      │
│  (128 - 2047)    │  5 + 6 = 11 payload bits                       │
├──────────────────┼────────────────────────────────────────────────┤
│ U+0800 - U+FFFF  │  1110xxxx  10xxxxxx  10xxxxxx   (3 bytes)      │
│  (2048 - 65535)  │  4 + 6 + 6 = 16 payload bits                  │
├──────────────────┼────────────────────────────────────────────────┤
│ U+10000-U+10FFFF │  11110xxx  10xxxxxx  10xxxxxx  10xxxxxx (4B)  │
│ (65536 - 1114111)│  3 + 6 + 6 + 6 = 21 payload bits              │
└──────────────────┴────────────────────────────────────────────────┘

Key design:
- 1-byte: starts with 0 (pure ASCII compatibility)
- 2-byte: starts with 110
- 3-byte: starts with 1110
- 4-byte: starts with 11110
- Continuation bytes always start with 10
```

### Four Detailed Examples

**Example 1: 'A' (U+0041 = decimal 65)**

```
U+0041 = 65 = 0100 0001 (binary)
Falls in range U+0000-U+007F → 1-byte encoding
UTF-8: 0 1000001  =  0x41  = 65
       ↑ 0 indicates 1-byte character (pure ASCII!)

Memory: [0x41]
```

**Example 2: 'é' (U+00E9 = decimal 233)**

```
U+00E9 = 233 = 1110 1001 (binary, 8 bits)
Falls in range U+0080-U+07FF → 2-byte encoding

Template: 110xxxxx 10xxxxxx
          233 = 000 1110 1001 (11 bits to fill template)
               ↑ leading zeros added to get 11 bits

Split: 000 1110 | 10 1001
       byte1: 110 + 000 1110 = 1100 0011 = 0xC3
       byte2:  10 + 10 1001  = 1010 1001 = 0xA9

Memory: [0xC3] [0xA9]  (2 bytes for one character)
```

**Example 3: '中' (U+4E2D = decimal 19,501)**

```
U+4E2D = 19501 = 0100 1110 0010 1101 (binary, 16 bits)
Falls in range U+0800-U+FFFF → 3-byte encoding

Template: 1110xxxx 10xxxxxx 10xxxxxx
          19501 = 0100 111000 101101 (16 bits split into 4+6+6)
          4-bit group:  0100
          6-bit group:  111000
          6-bit group:  101101

byte1: 1110 + 0100   = 1110 0100 = 0xE4
byte2:   10 + 111000 = 1011 1000 = 0xB8
byte3:   10 + 101101 = 1010 1101 = 0xAD

Memory: [0xE4] [0xB8] [0xAD]  (3 bytes for one character)
```

**Example 4: '😀' (U+1F600 = decimal 128,512)**

```
U+1F600 = 128512 = 0001 1111 0110 0000 0000 (binary, 21 bits needed)
Falls in range U+10000-U+10FFFF → 4-byte encoding

Template: 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
          128512 = 000 011111 011000 000000 (21 bits split 3+6+6+6)
          3-bit:   000
          6-bit:   011111
          6-bit:   011000
          6-bit:   000000

byte1: 11110 + 000   = 1111 0000 = 0xF0
byte2:    10 + 011111 = 1001 1111 = 0x9F
byte3:    10 + 011000 = 1001 1000 = 0x98
byte4:    10 + 000000 = 1000 0000 = 0x80

Memory: [0xF0] [0x9F] [0x98] [0x80]  (4 bytes for one emoji)
```

### Why UTF-8 Is Brilliant

1. **Backward-compatible with ASCII.** Any ASCII file is a valid UTF-8 file. The bit 7 of single-byte characters is always 0, matching ASCII exactly.

2. **Self-synchronizing.** You can find the start of any character by scanning backwards until you find a byte that does NOT start with `10`. This means if data is corrupted in the middle, the rest of the stream is still readable.

3. **No byte-order issues.** Unlike UTF-16, UTF-8 is the same on all hardware regardless of endianness.

4. **Compact for common text.** English text is 1 byte/char. Western European is 1-2 bytes/char. CJK is 3 bytes/char. Only rare supplementary characters need 4 bytes.

```
┌─────────────────────────────────────────────────────────────┐
│  "Hello, 世界" — byte layout in UTF-8                       │
│                                                             │
│  H    e    l    l    o    ,    SP   世        界            │
│  0x48 0x65 0x6C 0x6C 0x6F 0x2C 0x20 E4 B8 96 E7 95 8C     │
│  ←──────────────────────────────────→ ←──────→ ←──────→   │
│  7 bytes (ASCII range)                 3 bytes  3 bytes     │
│                                                             │
│  Total: 13 bytes for 9 characters                          │
│  len() in Astra returns 9 (characters)                     │
│  byte_len() in Astra returns 13 (bytes)                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. UTF-16 and UTF-32 — When They're Used

### UTF-16

UTF-16 uses 2 bytes for most characters and 4 bytes (called **surrogate pairs**) for characters above U+FFFF.

- Characters U+0000-U+FFFF: stored as 2 bytes (the code point value directly)
- Characters U+10000-U+10FFFF: encoded using two 2-byte "surrogates":
  - High surrogate: 0xD800-0xDBFF
  - Low surrogate: 0xDC00-0xDFFF

**Used by:** Windows API, Java (its `char` type is UTF-16), JavaScript (internally), .NET.

**Problem:** Even though most characters are 2 bytes, characters on supplementary planes (emoji, historic scripts) require 4 bytes and are NOT the same as two normal characters. Java/JavaScript code that does `str.length()` or `str[i]` gets broken results for emoji.

### UTF-32

UTF-32 uses exactly **4 bytes** for every code point. Simple and predictable — any character is at offset `index * 4`.

**Used by:** Internally in some Python implementations, some Unix systems, theoretical cases where O(1) character indexing matters more than memory.

**Problem:** Wastes memory. An English string uses 4x more memory than UTF-8.

### The Byte Order Mark (BOM)

A **BOM** is the character U+FEFF placed at the very start of a file to indicate the byte order:

```
UTF-8 BOM:    EF BB BF  (redundant but sometimes used by Windows)
UTF-16 LE:    FF FE
UTF-16 BE:    FE FF
UTF-32 LE:    FF FE 00 00
UTF-32 BE:    00 00 FE FF
```

Go and Astra use UTF-8 without BOM, which is the standard on Unix/Linux/macOS.

---

## 7. Strings in Go — Surprises and Best Practices

Go's `string` type looks simple but has important subtleties.

### A Go String Is a Byte Slice

In Go, a `string` is literally a pointer to bytes plus a length — it is an **immutable sequence of bytes**. The bytes are interpreted as UTF-8, but Go does not enforce this. You could store invalid UTF-8 in a Go string; it would compile and run, but string operations would produce garbage.

```go
s := "Hello, 世界"

// len() counts BYTES, not characters!
fmt.Println(len(s))   // Output: 13  (7 ASCII + 3 + 3 for two Chinese chars)

// Indexing gives a BYTE, not a character:
fmt.Println(s[7])     // Output: 228  (0xE4, first byte of '世')
fmt.Println(s[8])     // Output: 184  (0xB8, second byte of '世')
// s[7] does NOT give you '世'!
```

### The `rune` Type

Go's solution: the **rune** type. A `rune` is an alias for `int32` and represents a single Unicode code point (character).

```go
// Convert string to rune slice to get characters:
runes := []rune("Hello, 世界")
fmt.Println(len(runes))  // Output: 9  (nine characters!)
fmt.Println(string(runes[7]))  // Output: 世
fmt.Println(string(runes[8]))  // Output: 界
```

### Range Over String Gives Runes

When you use `for range` over a string in Go, it automatically decodes UTF-8 and gives you runes:

```go
s := "Hello, 世界"
for i, r := range s {
    // i is the BYTE index of the start of the rune
    // r is the rune (decoded code point)
    fmt.Printf("byte[%d] = U+%04X ('%c')\n", i, r, r)
}
// Output:
// byte[0] = U+0048 ('H')
// byte[1] = U+0065 ('e')
// ...
// byte[7] = U+4E16 ('世')   ← note: next index will be 10, not 8!
// byte[10] = U+754C ('界')
```

Notice: after the character at byte index 7 (which is 3 bytes), the next character starts at byte index 10.

### Important Go String Functions

```go
import (
    "strings"
    "fmt"
    "unicode/utf8"
)

s := "Hello, World"

// Basic operations:
fmt.Println(strings.Contains(s, "World"))       // true
fmt.Println(strings.HasPrefix(s, "Hello"))      // true
fmt.Println(strings.HasSuffix(s, "World"))      // true
fmt.Println(strings.Index(s, "World"))          // 7
fmt.Println(strings.ToUpper(s))                 // HELLO, WORLD
fmt.Println(strings.ToLower(s))                 // hello, world
fmt.Println(strings.TrimSpace("  hello  "))     // "hello"
fmt.Println(strings.Split("a,b,c", ","))        // ["a", "b", "c"]
fmt.Println(strings.Join([]string{"a","b"}, "-")) // "a-b"
fmt.Println(strings.Replace(s, "World", "Astra", 1)) // "Hello, Astra"
fmt.Println(strings.Count(s, "l"))              // 3

// Unicode-aware length:
fmt.Println(utf8.RuneCountInString("Hello, 世界"))  // 9 (characters)
fmt.Println(len("Hello, 世界"))                     // 13 (bytes)
```

### String Builder for Efficient Concatenation

In Go (and in Astra's runtime), repeated string concatenation with `+` is slow because each `+` allocates a new string. Use `strings.Builder`:

```go
// SLOW: creates a new string for each +:
result := ""
for i := 0; i < 1000; i++ {
    result = result + "ha"   // 1000 allocations!
}

// FAST: strings.Builder accumulates into a single buffer:
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("ha")
}
result := sb.String()   // one final allocation
```

**Why is repeated `+` slow?** Each concatenation allocates a new string of length `len(a) + len(b)`, copies all bytes from `a`, then copies all bytes from `b`. For N concatenations each adding M bytes, this is O(N²·M) total work. `strings.Builder` uses an exponentially-growing buffer, amortizing allocations to O(N·M).

### Format Strings with `fmt.Sprintf`

```go
name := "Astra"
version := 1
greeting := fmt.Sprintf("Welcome to %s version %d!", name, version)
// "Welcome to Astra version 1!"

pi := 3.14159
formatted := fmt.Sprintf("Pi is approximately %.2f", pi)
// "Pi is approximately 3.14"
```

### Common Go String Mistakes

```go
// MISTAKE 1: using len() to count characters
s := "café"
fmt.Println(len(s))         // 5, not 4! ('é' is 2 bytes in UTF-8)
fmt.Println(len([]rune(s))) // 4 — correct character count

// MISTAKE 2: indexing into a multi-byte string
s := "世界"
fmt.Println(s[0])           // 228 — a byte value, NOT '世'
// Correct: use []rune(s)[0] or range loop

// MISTAKE 3: string comparison is fine for content, but:
a := "café"   // might be 'e' + combining accent (2 code points)
b := "café"   // might be é as single code point (U+00E9)
fmt.Println(a == b)  // might be FALSE! Same visual, different bytes
// For robust comparison, use unicode normalization (golang.org/x/text)
```

---

## 8. Strings in Astra — Internal Representation

Astra's compiler generates C code (or LLVM IR) that is compiled to native machine code. At the native level, Astra strings are represented as a C struct:

```c
/* runtime/astra_string.h */

typedef struct {
    uint8_t* data;       /* pointer to UTF-8 encoded bytes (NOT null-terminated) */
    int64_t  byte_len;   /* number of bytes in data */
    int64_t  len;        /* number of Unicode code points (characters) */
    uint64_t hash;       /* cached hash value (0 = not computed yet) */
} AstraString;
```

**Why not null-terminated (like C strings)?**

C uses `\0` (byte value 0) to mark the end of a string. This fails for binary data (which may contain `\0`) and requires O(n) `strlen()` to find the length. Astra stores the length explicitly, so `len()` is O(1).

**Why store both `byte_len` and `len`?**

- `byte_len`: needed for memory operations (how many bytes to copy, allocate, or hash)
- `len`: needed for user-facing operations (how many characters does the user see?)

If you only stored one, every call to the other would require O(n) scanning.

### Memory Layout

```
┌───────────────────────────────────────────────────────────────┐
│  AstraString for "Hello, 世界"                                │
│                                                               │
│  ┌──────────┬──────────┬──────────┬──────────┐               │
│  │  data    │ byte_len │   len    │   hash   │               │
│  │  ptr     │   13     │    9     │    0     │               │
│  └────┬─────┴──────────┴──────────┴──────────┘               │
│       │                                                       │
│       ▼  (points to UTF-8 bytes in .rodata or heap)          │
│  ┌────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┐ │
│  │ H  │ e  │ l  │ l  │ o  │ ,  │ SP │E4  │B8  │96  │E7  │95  │8C │ │
│  │ 48 │ 65 │ 6C │ 6C │ 6F │ 2C │ 20 │(世=3 bytes)│(界=3 bytes)│ │
│  └────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┴────┘ │
│   0    1    2    3    4    5    6    7    8    9    10   11   12      │
└───────────────────────────────────────────────────────────────┘
```

### String Literals in the Compiled Binary

String literals like `"Hello, World"` are placed in the `.rodata` (read-only data) section of the compiled binary. They are never freed or modified — they live for the lifetime of the program.

```
┌────────────────────────────────────────┐
│  Compiled Astra Binary                 │
├────────────────────────────────────────┤
│  .text     (machine code)              │
├────────────────────────────────────────┤
│  .rodata   (read-only data)            │
│  ┌──────────────────────────────────┐  │
│  │ 0x10a0: "Hello, World\0"         │  │
│  │ 0x10ae: "Hello, 世界\0"          │  │
│  │ 0x10bc: "Error: division by zero"│  │
│  └──────────────────────────────────┘  │
├────────────────────────────────────────┤
│  .data     (mutable global variables)  │
├────────────────────────────────────────┤
│  .bss      (zero-initialized globals)  │
└────────────────────────────────────────┘
```

### Copy-on-Write Semantics

Astra strings are **immutable**. When you do:

```astra
let a = "hello"
let b = a           // b and a share the same bytes! No copy.
let c = a + "!"     // NEW string allocated; a is unchanged
```

This is safe because strings never change after creation. Two variables can point to the same bytes without conflict. Only concatenation (or other "mutating" operations) allocates new memory.

---

## 9. The Astra String Standard Library

Astra provides a rich string API. Here is the complete set of operations for v1:

```astra
// All string methods — available as s.method() syntax in Astra

// LENGTH AND NAVIGATION
s.len() -> int              // number of Unicode characters (code points)
s.byte_len() -> int         // number of UTF-8 bytes
s.is_empty() -> bool        // true if len() == 0
s.char_at(i: int) -> string // character at position i (returns 1-char string)
s.bytes() -> List<int>      // UTF-8 bytes as list of integers

// SEARCH AND TEST
s.contains(sub: string) -> bool      // does s contain sub?
s.starts_with(prefix: string) -> bool
s.ends_with(suffix: string) -> bool
s.index_of(sub: string) -> int       // -1 if not found
s.count(sub: string) -> int          // count non-overlapping occurrences

// TRANSFORMATION
s.to_upper() -> string       // "hello" → "HELLO"
s.to_lower() -> string       // "HELLO" → "hello"
s.trim() -> string           // remove leading/trailing whitespace
s.trim_start() -> string     // remove leading whitespace only
s.trim_end() -> string       // remove trailing whitespace only
s.replace(old: string, new: string) -> string  // replace first occurrence
s.replace_all(old: string, new: string) -> string
s.reverse() -> string        // reverse characters (not bytes!)

// SPLITTING AND JOINING
s.split(sep: string) -> List<string>  // split by separator
s.split_lines() -> List<string>       // split by \n or \r\n
"sep".join(parts: List<string>) -> string  // join with separator

// CONVERSION
s.to_int() -> int            // parse decimal integer, panics on error
s.to_float() -> float        // parse float
s.to_int_or(default: int) -> int  // parse or return default
```

---

## 10. 🔨 Astra Build Milestone — AstraString Implementation

Here is the complete Go implementation of Astra's string handling — the code the compiler uses to represent, compare, and intern strings during compilation.

```go
// runtime/string.go
// This represents how the Astra runtime handles strings.
// The compiler generates code that calls these (or equivalent C functions).
package runtime

import (
    "fmt"
    "hash/fnv"
    "sync"
    "unicode/utf8"
)

// ============================================================
// AstraString — the compiler's representation of a string value
// ============================================================

// AstraString is the internal representation of an Astra string.
// It mirrors the C struct in the compiled runtime.
type AstraString struct {
    data    []byte  // UTF-8 encoded bytes
    byteLen int64   // len(data)
    charLen int64   // number of Unicode code points
    hash    uint64  // cached FNV hash (0 means not computed)
}

// NewAstraString creates an AstraString from a Go string.
// This is called by the compiler when it processes string literals.
func NewAstraString(s string) *AstraString {
    data := []byte(s)

    // Count Unicode code points (characters, not bytes)
    charLen := int64(utf8.RuneCountInString(s))

    return &AstraString{
        data:    data,
        byteLen: int64(len(data)),
        charLen: charLen,
        hash:    0, // computed lazily
    }
}

// Len returns the number of Unicode characters (code points).
// This is what Astra's s.len() exposes to the user.
func (s *AstraString) Len() int64 { return s.charLen }

// ByteLen returns the number of UTF-8 bytes.
func (s *AstraString) ByteLen() int64 { return s.byteLen }

// GoString converts back to a Go string (for use in compiler internals).
func (s *AstraString) GoString() string { return string(s.data) }

// Hash returns the FNV-64 hash of the string content.
// Used for string interning and hash map keys.
func (s *AstraString) Hash() uint64 {
    if s.hash == 0 {
        h := fnv.New64a()
        h.Write(s.data)
        s.hash = h.Sum64()
        if s.hash == 0 {
            s.hash = 1 // 0 is reserved for "not computed"
        }
    }
    return s.hash
}

// Equals compares two AstraStrings by content (NOT by pointer).
// This is what Astra's == operator does for strings.
func (s *AstraString) Equals(other *AstraString) bool {
    if s.byteLen != other.byteLen {
        return false  // fast path: different length → definitely different
    }
    if s.hash != 0 && other.hash != 0 && s.hash != other.hash {
        return false  // fast path: different hash → definitely different
    }
    // Full byte-by-byte comparison:
    for i := range s.data {
        if s.data[i] != other.data[i] {
            return false
        }
    }
    return true
}

// Concat returns a new AstraString that is s + other.
// Allocates a new byte slice — strings are immutable.
func (s *AstraString) Concat(other *AstraString) *AstraString {
    newData := make([]byte, s.byteLen+other.byteLen)
    copy(newData, s.data)
    copy(newData[s.byteLen:], other.data)
    return &AstraString{
        data:    newData,
        byteLen: s.byteLen + other.byteLen,
        charLen: s.charLen + other.charLen,
        hash:    0,
    }
}

// CharAt returns the Unicode character at the given position (0-indexed).
// Returns (char, true) on success, ("", false) if out of bounds.
// This is O(n) because UTF-8 is variable-width.
func (s *AstraString) CharAt(pos int64) (string, bool) {
    if pos < 0 || pos >= s.charLen {
        return "", false
    }
    // Scan through runes to find the one at position pos
    var i int64
    for _, r := range string(s.data) {
        if i == pos {
            return string(r), true
        }
        i++
    }
    return "", false // should not reach here
}

// String returns a human-readable representation (for debugging).
func (s *AstraString) String() string {
    return fmt.Sprintf("AstraString{%q, bytes=%d, chars=%d}",
        string(s.data), s.byteLen, s.charLen)
}

// ============================================================
// String Interning Pool
// ============================================================
// String interning: store each unique string ONCE and reuse the pointer.
// Benefits:
//   1. == comparison becomes pointer equality (O(1) instead of O(n))
//   2. Saves memory when many variables hold the same string value
//   3. String literals in the binary are automatically interned

// StringPool is a thread-safe pool of interned strings.
type StringPool struct {
    mu      sync.RWMutex
    strings map[string]*AstraString  // content → canonical pointer
}

// Global string pool used by the compiler and runtime.
var globalPool = &StringPool{
    strings: make(map[string]*AstraString),
}

// Intern returns the canonical AstraString for the given content.
// If this content has been seen before, returns the existing pointer.
// Otherwise, creates and stores a new AstraString.
func (p *StringPool) Intern(s string) *AstraString {
    // Fast path: check with read lock
    p.mu.RLock()
    if existing, ok := p.strings[s]; ok {
        p.mu.RUnlock()
        return existing
    }
    p.mu.RUnlock()

    // Slow path: insert with write lock
    p.mu.Lock()
    defer p.mu.Unlock()
    // Check again (another goroutine might have inserted between our two locks)
    if existing, ok := p.strings[s]; ok {
        return existing
    }
    as := NewAstraString(s)
    p.strings[s] = as
    return as
}

// InternString is a convenience function using the global pool.
func InternString(s string) *AstraString {
    return globalPool.Intern(s)
}

// ============================================================
// Demonstration: "Hello, 世界" memory layout
// ============================================================

func DemoHelloWorld() {
    s := NewAstraString("Hello, 世界")

    fmt.Printf("String: %q\n", s.GoString())
    fmt.Printf("Byte length (byte_len): %d\n", s.ByteLen())  // 13
    fmt.Printf("Char length (len):      %d\n", s.Len())      // 9
    fmt.Printf("Hash:                   %d\n", s.Hash())

    fmt.Println("\nByte-by-byte layout:")
    for i, b := range s.data {
        fmt.Printf("  byte[%2d] = 0x%02X (%3d)\n", i, b, b)
    }

    fmt.Println("\nCharacter-by-character (Unicode code points):")
    charIdx := int64(0)
    for _, r := range s.GoString() {
        char, _ := s.CharAt(charIdx)
        fmt.Printf("  char[%d] = U+%04X = %q\n", charIdx, r, char)
        charIdx++
    }
}

// Expected output:
// String: "Hello, 世界"
// Byte length (byte_len): 13
// Char length (len):      9
// Hash:                   <some number>
//
// Byte-by-byte layout:
//   byte[ 0] = 0x48 ( 72)    'H'
//   byte[ 1] = 0x65 (101)    'e'
//   byte[ 2] = 0x6C (108)    'l'
//   byte[ 3] = 0x6C (108)    'l'
//   byte[ 4] = 0x6F (111)    'o'
//   byte[ 5] = 0x2C ( 44)    ','
//   byte[ 6] = 0x20 ( 32)    ' '
//   byte[ 7] = 0xE4 (228)    first byte of '世'
//   byte[ 8] = 0xB8 (184)    second byte of '世'
//   byte[ 9] = 0x96 (150)    third byte of '世'
//   byte[10] = 0xE7 (231)    first byte of '界'
//   byte[11] = 0x95 (149)    second byte of '界'
//   byte[12] = 0x8C (140)    third byte of '界'
//
// Character-by-character:
//   char[0] = U+0048 = "H"
//   char[1] = U+0065 = "e"
//   ...
//   char[7] = U+4E16 = "世"
//   char[8] = U+754C = "界"
```

### How `let s = "Hello, 世界"` Is Compiled

When the Astra compiler sees this source code:

```astra
let s = "Hello, 世界"
```

It performs these steps:

1. **Lexer** tokenizes `"Hello, 世界"` as a `STRING_LITERAL` token with value `Hello, 世界`
2. **Parser** creates a `StringLiteral` AST node with `Value = "Hello, 世界"`
3. **Code generator** calls `InternString("Hello, 世界")` to get the canonical `AstraString` pointer
4. **Backend** places the UTF-8 bytes `[0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0xE4, 0xB8, 0x96, 0xE7, 0x95, 0x8C]` in the `.rodata` section
5. **Generated code** at runtime: creates an `AstraString` struct on the stack with `data = &rodata[offset]`, `byte_len = 13`, `len = 9`

---

## Exercises

1. **UTF-8 encoding by hand.** Encode the following characters to UTF-8 bytes (show your work using the encoding table):
   - `ñ` (U+00F1)
   - `α` (U+03B1 = decimal 945)
   - `★` (U+2605 = decimal 9733)
   *Hint: determine which row of the encoding table applies, then fill in the bit template.*

2. **The length bug.** In Go, explain why `len("café")` returns 5 instead of 4. Write a Go function `runeLen(s string) int` that correctly counts characters. Test it on: `"hello"`, `"café"`, `"こんにちは"`, `"😊😊"`.
   *Hint: use `utf8.RuneCountInString` or convert to `[]rune`.*

3. **Mojibake detective.** The following string was stored in Windows-1252 encoding but read as UTF-8: `â€œHello worldâ€`. What was the original string? What kind of text uses the bytes that Windows-1252 uses for those characters?
   *Hint: the bytes 0x93 and 0x94 in Windows-1252 are "smart quotes" `"` and `"`.*

4. **String interning performance.** Add a method `Size() int` to `StringPool` that returns the number of interned strings. Write a Go program that interns 10,000 strings where 5,000 are duplicates. Verify that the pool contains exactly 5,000 entries. Measure the time for the first intern vs the time for an already-interned string.
   *Hint: use `time.Now()` and `time.Since()` to measure.*

5. **Implement `CharAt` efficiently.** The current `CharAt` implementation in the milestone is O(n). Under what circumstances could you make it O(1)? Design a data structure (`OffsetTable`) that pre-computes the byte offset of each character. When is this worth doing, and when is the O(n) scan faster?
   *Hint: An offset table is an array where `offsets[i]` is the byte position of character `i`. It costs O(n) to build and O(1) to query. For strings that are indexed many times, this pays off.*

6. **Reverse a Unicode string correctly.** Write a Go function `reverseString(s string) string` that reverses the characters (not the bytes) of a UTF-8 string. Verify: `reverseString("café")` = `"éfac"`, `reverseString("Hello, 世界")` = `"界世 ,olleH"`.
   *Hint: convert to `[]rune`, reverse the slice, convert back to string.*

---

## Summary

| Concept | Definition | Key Detail |
|---|---|---|
| ASCII | 7-bit encoding, 128 characters | Only covers English/basic Latin |
| Code page | 8-bit encoding for one language | Incompatible between languages; caused mojibake |
| Unicode | Standard assigning numbers to all characters | Over 1.1M possible code points (U+0000-U+10FFFF) |
| Code point | A unique number for one character | `U+1F600` = 😀 |
| UTF-8 | Variable-width Unicode encoding | 1-4 bytes; backward-compatible with ASCII |
| UTF-16 | 2 or 4 bytes per code point | Used by Java, Windows API |
| UTF-32 | 4 bytes per code point | Simple but wasteful |
| BOM | Byte Order Mark, U+FEFF at file start | Identifies encoding and byte order |
| Go `string` | Immutable byte slice | `len()` counts bytes, not characters! |
| Go `rune` | Alias for `int32`; one Unicode code point | Use `[]rune(s)` for character indexing |
| String interning | Store each unique string once | Enables O(1) `==` comparison |
| `.rodata` | Read-only data section of binary | String literals live here |
| AstraString | `{data, byte_len, len, hash}` | `len` = chars; `byte_len` = bytes |
| Mojibake | Garbled text from encoding mismatch | Open Shift-JIS as Latin-1 |
| Self-synchronizing | Can find char boundary by scanning back | Bytes starting with `10` are continuation bytes |

Text encoding is one of those topics where ignorance causes real bugs. Every Astra programmer who understands this chapter will avoid the classic mistakes: using `len()` for character counting, comparing floats for equality, and assuming one byte equals one character. The `AstraString` struct you built here is the foundation that makes Astra's string operations correct for all of humanity's writing systems.
