# Chapter 65: Standard Library — Core: string, math, time

> "A language is not complete until you can do something useful with it. The standard library is what transforms a compiler into a tool." — Rob Pike

---

## Overview

You have a compiler. It parses Astra source code, builds an AST, performs semantic analysis, generates IR, and emits machine code. Programs can declare variables, define functions, use control flow, and manage memory. But right now, an Astra program cannot even print the current time without writing C code by hand.

The standard library is the bridge between a working compiler and a useful language. Every professional language ships with one: Go's `fmt`, `os`, `net/http`; Python's `math`, `datetime`, `json`; Rust's `std::io`, `std::collections`. The standard library is so important that most programmers never think about it — it is simply always there, the invisible foundation beneath every real program.

This chapter builds the core of Astra's standard library: three foundational packages that nearly every program will use.

**The string package** — because every program works with text. Strings are the universal interface of computing: file paths, user input, network data, configuration, output. We will implement Astra's string type from the ground up, including full Unicode support, a rich method set, and an efficient string builder.

**The math package** — because numeric computation is fundamental. We will implement the full set of mathematical functions, constants, and a clean random number generation API.

**The time package** — because programs exist in time. Getting the current timestamp, sleeping, measuring elapsed duration, formatting and parsing dates — all essential.

Along the way, we will trace the complete pipeline: how a method call like `"hello".to_upper()` in Astra source code becomes a C function call in the generated machine code, passing through the AST, IR, and assembler along the way.

---

## What We're Building

Three Go files that implement Astra's runtime standard library, backed by C runtime functions for performance-critical operations:

```
stdlib/
  string/
    string.go         ← string package implementation (~350 lines)
  math/
    math.go           ← math package implementation (~200 lines)
  time/
    time.go           ← time package implementation (~180 lines)

runtime/
  astra_string.h      ← C struct definition and function declarations
  astra_string.c      ← C implementation of UTF-8 string operations
  astra_math.h        ← math runtime helpers
  astra_time.h        ← time runtime helpers
```

---

## Table of Contents

1. Astra String Internals — the `AstraString` C struct
2. Unicode and UTF-8 — why strings are hard
3. String methods — complete implementation
4. The String Builder — efficient concatenation
5. The Math Package — constants, functions, random numbers
6. The Time Package — timestamps, durations, formatting
7. How method calls compile — the full pipeline
8. Build Milestone
9. Exercises

---

## 1. Astra String Internals

### The AstraString Struct

Before we can implement string methods in Go, we need to understand how Astra stores strings in memory. Strings are not just `char*` pointers like in C. Astra strings are Unicode-aware, reference-counted value types stored in a heap-allocated struct.

```c
// runtime/astra_string.h

#ifndef ASTRA_STRING_H
#define ASTRA_STRING_H

#include <stdint.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    char*   data;       // UTF-8 bytes (null-terminated for C interop)
    int64_t len;        // number of Unicode code points (characters)
    int64_t byte_len;   // number of bytes (may differ from len for non-ASCII)
    int32_t ref_count;  // reference count for memory management
    int32_t _pad;       // padding to align to 8 bytes
} AstraString;

// Constructors
AstraString* astra_string_new(const char* data, int64_t byte_len);
AstraString* astra_string_from_cstr(const char* cstr);
AstraString* astra_string_empty(void);

// Reference counting
void astra_string_retain(AstraString* s);
void astra_string_release(AstraString* s);

// Basic info
int64_t astra_string_length(AstraString* s);      // Unicode char count
int64_t astra_string_byte_length(AstraString* s); // byte count
int     astra_string_is_empty(AstraString* s);

// Case conversion
AstraString* astra_string_to_upper(AstraString* s);
AstraString* astra_string_to_lower(AstraString* s);

// Search
int     astra_string_contains(AstraString* s, AstraString* sub);
int     astra_string_starts_with(AstraString* s, AstraString* prefix);
int     astra_string_ends_with(AstraString* s, AstraString* suffix);
int64_t astra_string_index_of(AstraString* s, AstraString* sub);

// Transformation
AstraString* astra_string_slice(AstraString* s, int64_t start, int64_t end);
AstraString* astra_string_replace(AstraString* s, AstraString* from, AstraString* to);
AstraString* astra_string_trim(AstraString* s);
AstraString* astra_string_trim_start(AstraString* s);
AstraString* astra_string_trim_end(AstraString* s);
AstraString* astra_string_repeat(AstraString* s, int64_t n);
AstraString* astra_string_reverse(AstraString* s);
AstraString* astra_string_concat(AstraString* a, AstraString* b);

// Conversion
int64_t astra_string_parse_int(AstraString* s, int* ok);
double  astra_string_parse_float(AstraString* s, int* ok);
AstraString* astra_int_to_string(int64_t n);
AstraString* astra_float_to_string(double f);

#endif // ASTRA_STRING_H
```

The key insight here is the separation of `len` (Unicode code points) and `byte_len` (raw bytes). Consider the string `"Hello, 世界! 🌍"`:

```
Character:  H  e  l  l  o  ,     世    界    !     🌍
Code point: U+0048 U+0065 U+006C U+006C U+006F U+002C U+0020 U+4E16 U+754C U+0021 U+0020 U+1F30D
UTF-8 bytes:
  H  = 0x48                        (1 byte)
  e  = 0x65                        (1 byte)
  l  = 0x6C                        (1 byte)
  l  = 0x6C                        (1 byte)
  o  = 0x6F                        (1 byte)
  ,  = 0x2C                        (1 byte)
  SP = 0x20                        (1 byte)
  世 = 0xE4 0xB8 0x96              (3 bytes)
  界 = 0xE7 0x95 0x8C              (3 bytes)
  !  = 0x21                        (1 byte)
  SP = 0x20                        (1 byte)
  🌍 = 0xF0 0x9F 0x8C 0x8D        (4 bytes)

len      = 12 (Unicode characters)
byte_len = 20 (UTF-8 bytes)
```

If we used only byte length, `s.slice(7, 9)` would give us garbage — we would be splitting in the middle of multi-byte UTF-8 sequences. The `len` field tracks character count, and all our operations work in terms of characters.

---

## 2. Unicode and UTF-8

### Why UTF-8?

UTF-8 is the dominant encoding of the modern web. Over 98% of web pages use it. It has three critical properties that make it ideal for Astra:

1. **ASCII compatible**: the first 128 code points (U+0000–U+007F) encode as a single byte identical to ASCII. Pure ASCII strings have `len == byte_len`.

2. **Self-synchronizing**: you can always tell whether a byte is the start of a character or a continuation byte. Bytes 0x80–0xBF are continuation bytes; all others start a new character.

3. **No null bytes in multi-byte sequences**: this means UTF-8 strings can be null-terminated for C interop without ambiguity.

### UTF-8 Encoding Rules

```
Code point range       Byte sequence
U+0000   – U+007F      0xxxxxxx
U+0080   – U+07FF      110xxxxx 10xxxxxx
U+0800   – U+FFFF      1110xxxx 10xxxxxx 10xxxxxx
U+10000  – U+10FFFF    11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```

The C runtime needs to walk UTF-8 sequences to implement character-level operations. Here is the core utility:

```c
// runtime/astra_string.c

#include "astra_string.h"
#include <stdio.h>
#include <ctype.h>
#include <assert.h>

// Returns number of bytes in the UTF-8 sequence starting at `b`.
// Returns 1 for invalid bytes (graceful degradation).
static int utf8_char_len(unsigned char b) {
    if      ((b & 0x80) == 0x00) return 1;  // 0xxxxxxx
    else if ((b & 0xE0) == 0xC0) return 2;  // 110xxxxx
    else if ((b & 0xF0) == 0xE0) return 3;  // 1110xxxx
    else if ((b & 0xF8) == 0xF0) return 4;  // 11110xxx
    return 1; // invalid: treat as single byte
}

// Count the number of Unicode code points in `byte_len` bytes of `data`.
static int64_t utf8_count_chars(const char* data, int64_t byte_len) {
    int64_t count = 0;
    int64_t i = 0;
    while (i < byte_len) {
        i += utf8_char_len((unsigned char)data[i]);
        count++;
    }
    return count;
}

// Get the byte offset of the Nth character (0-indexed).
// Returns -1 if n >= char_count.
static int64_t utf8_char_offset(const char* data, int64_t byte_len, int64_t n) {
    int64_t i = 0;
    int64_t char_idx = 0;
    while (i < byte_len && char_idx < n) {
        i += utf8_char_len((unsigned char)data[i]);
        char_idx++;
    }
    return (char_idx == n) ? i : -1;
}

// Constructor: create from raw bytes
AstraString* astra_string_new(const char* data, int64_t byte_len) {
    AstraString* s = (AstraString*)malloc(sizeof(AstraString));
    if (!s) return NULL;
    s->data = (char*)malloc(byte_len + 1);
    if (!s->data) { free(s); return NULL; }
    memcpy(s->data, data, byte_len);
    s->data[byte_len] = '\0';
    s->byte_len = byte_len;
    s->len = utf8_count_chars(data, byte_len);
    s->ref_count = 1;
    s->_pad = 0;
    return s;
}

AstraString* astra_string_from_cstr(const char* cstr) {
    return astra_string_new(cstr, (int64_t)strlen(cstr));
}

AstraString* astra_string_empty(void) {
    return astra_string_new("", 0);
}

void astra_string_retain(AstraString* s) {
    if (s) s->ref_count++;
}

void astra_string_release(AstraString* s) {
    if (!s) return;
    s->ref_count--;
    if (s->ref_count <= 0) {
        free(s->data);
        free(s);
    }
}

int64_t astra_string_length(AstraString* s)      { return s ? s->len : 0; }
int64_t astra_string_byte_length(AstraString* s) { return s ? s->byte_len : 0; }
int     astra_string_is_empty(AstraString* s)    { return !s || s->len == 0; }
```

---

## 3. String Methods — Complete Implementation

### The Go Layer

The Go stdlib is what the Astra compiler links against. It wraps the C runtime functions and provides the Astra-facing API. Go calls C via `cgo`.

```go
// stdlib/string/string.go

package astra_string

/*
#cgo CFLAGS: -I../../runtime
#cgo LDFLAGS: -L../../runtime -lastra_runtime
#include "astra_string.h"
#include <stdlib.h>
*/
import "C"
import (
    "fmt"
    "strings"
    "unsafe"
)

// AstraString wraps the C AstraString struct for Go.
type AstraString struct {
    ptr *C.AstraString
}

// NewString creates an AstraString from a Go string.
func NewString(s string) AstraString {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    ptr := C.astra_string_new(cs, C.int64_t(len(s)))
    return AstraString{ptr: ptr}
}

// GoString returns the Go string value.
func (s AstraString) GoString() string {
    if s.ptr == nil {
        return ""
    }
    return C.GoStringN(s.ptr.data, C.int(s.ptr.byte_len))
}

// --- Basic Info ---

// Length returns the number of Unicode code points.
func (s AstraString) Length() int64 {
    return int64(C.astra_string_length(s.ptr))
}

// ByteLength returns the number of UTF-8 bytes.
func (s AstraString) ByteLength() int64 {
    return int64(C.astra_string_byte_length(s.ptr))
}

// IsEmpty returns true if the string has no characters.
func (s AstraString) IsEmpty() bool {
    return C.astra_string_is_empty(s.ptr) != 0
}

// --- Case Conversion ---

// ToUpper returns a new string with all characters in upper case.
// For ASCII characters, uses the standard ASCII rule.
// For Unicode letters, we delegate to Go's unicode package through
// a pure-Go implementation for correctness, then create a new AstraString.
func (s AstraString) ToUpper() AstraString {
    // Pure Go: handles full Unicode case folding correctly
    upper := strings.ToUpper(s.GoString())
    return NewString(upper)
}

// ToLower returns a new string with all characters in lower case.
func (s AstraString) ToLower() AstraString {
    lower := strings.ToLower(s.GoString())
    return NewString(lower)
}

// --- Search ---

// Contains returns true if the string contains the substring.
func (s AstraString) Contains(sub AstraString) bool {
    return strings.Contains(s.GoString(), sub.GoString())
}

// StartsWith returns true if the string begins with prefix.
func (s AstraString) StartsWith(prefix AstraString) bool {
    return strings.HasPrefix(s.GoString(), prefix.GoString())
}

// EndsWith returns true if the string ends with suffix.
func (s AstraString) EndsWith(suffix AstraString) bool {
    return strings.HasSuffix(s.GoString(), suffix.GoString())
}

// IndexOf returns the character index (not byte index) of the first
// occurrence of sub, or -1 if not found.
func (s AstraString) IndexOf(sub AstraString) int64 {
    str := s.GoString()
    subStr := sub.GoString()
    byteIdx := strings.Index(str, subStr)
    if byteIdx < 0 {
        return -1
    }
    // Convert byte index to character index
    charIdx := int64(len([]rune(str[:byteIdx])))
    return charIdx
}

// --- Transformation ---

// Slice returns the substring from character index `start` (inclusive)
// to `end` (exclusive). Supports negative indices (Python-style).
func (s AstraString) Slice(start, end int64) AstraString {
    runes := []rune(s.GoString())
    n := int64(len(runes))

    // Clamp to valid range
    if start < 0 { start = 0 }
    if end > n   { end = n }
    if start >= end {
        return NewString("")
    }

    return NewString(string(runes[start:end]))
}

// Replace replaces the first occurrence of `from` with `to`.
// Use ReplaceAll for all occurrences.
func (s AstraString) Replace(from, to AstraString) AstraString {
    result := strings.Replace(s.GoString(), from.GoString(), to.GoString(), 1)
    return NewString(result)
}

// ReplaceAll replaces all occurrences of `from` with `to`.
func (s AstraString) ReplaceAll(from, to AstraString) AstraString {
    result := strings.ReplaceAll(s.GoString(), from.GoString(), to.GoString())
    return NewString(result)
}

// Trim removes leading and trailing whitespace (spaces, tabs, newlines).
func (s AstraString) Trim() AstraString {
    return NewString(strings.TrimSpace(s.GoString()))
}

// TrimStart removes only leading whitespace.
func (s AstraString) TrimStart() AstraString {
    return NewString(strings.TrimLeft(s.GoString(), " \t\n\r\v\f"))
}

// TrimEnd removes only trailing whitespace.
func (s AstraString) TrimEnd() AstraString {
    return NewString(strings.TrimRight(s.GoString(), " \t\n\r\v\f"))
}

// Split splits the string by `sep` and returns a slice of AstraStrings.
func (s AstraString) Split(sep AstraString) []AstraString {
    parts := strings.Split(s.GoString(), sep.GoString())
    result := make([]AstraString, len(parts))
    for i, p := range parts {
        result[i] = NewString(p)
    }
    return result
}

// Repeat returns the string repeated n times.
func (s AstraString) Repeat(n int64) AstraString {
    return NewString(strings.Repeat(s.GoString(), int(n)))
}

// Reverse returns the string with Unicode characters in reverse order.
// Note: this reverses code points, not grapheme clusters. Composed
// characters (e.g. é = e + combining accent) may not reverse intuitively.
func (s AstraString) Reverse() AstraString {
    runes := []rune(s.GoString())
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return NewString(string(runes))
}

// Concat concatenates two strings.
func (s AstraString) Concat(other AstraString) AstraString {
    return NewString(s.GoString() + other.GoString())
}

// --- Conversion ---

// ParseInt parses the string as a base-10 integer.
// Returns (0, false) if parsing fails.
func (s AstraString) ParseInt() (int64, bool) {
    var n int64
    str := strings.TrimSpace(s.GoString())
    _, err := fmt.Sscanf(str, "%d", &n)
    return n, err == nil
}

// ParseFloat parses the string as a 64-bit floating-point number.
func (s AstraString) ParseFloat() (float64, bool) {
    var f float64
    str := strings.TrimSpace(s.GoString())
    _, err := fmt.Sscanf(str, "%f", &f)
    return f, err == nil
}

// IntToString converts an int64 to an AstraString.
func IntToString(n int64) AstraString {
    return NewString(fmt.Sprintf("%d", n))
}

// FloatToString converts a float64 to an AstraString.
func FloatToString(f float64) AstraString {
    // Format with enough precision to round-trip, strip trailing zeros
    s := fmt.Sprintf("%g", f)
    return NewString(s)
}

// --- Format ---

// Format replaces {} placeholders with the given arguments in order.
// Example: "Hello, {}! You are {} years old.".format("Aditya", 25)
func (s AstraString) Format(args ...interface{}) AstraString {
    str := s.GoString()
    var result strings.Builder
    argIdx := 0

    i := 0
    for i < len(str) {
        if str[i] == '{' && i+1 < len(str) && str[i+1] == '}' {
            if argIdx < len(args) {
                result.WriteString(fmt.Sprintf("%v", args[argIdx]))
                argIdx++
            } else {
                result.WriteString("{}")
            }
            i += 2
        } else {
            result.WriteByte(str[i])
            i++
        }
    }
    return NewString(result.String())
}
```

### The String Builder

Naive string concatenation inside a loop is O(n²) — each `+` allocates a new string and copies all the old bytes. The `StringBuilder` accumulates bytes and performs a single allocation at the end.

```go
// stdlib/string/builder.go

package astra_string

import "strings"

// Builder efficiently builds strings through multiple Append calls.
// Uses Go's strings.Builder internally, which grows exponentially
// to amortize allocation cost.
type Builder struct {
    inner strings.Builder
    charCount int64
}

// NewBuilder creates a fresh string builder.
func NewBuilder() *Builder {
    return &Builder{}
}

// Append adds a string to the builder.
func (b *Builder) Append(s AstraString) {
    str := s.GoString()
    b.inner.WriteString(str)
    // Count runes for accurate char length tracking
    for _, _ = range str {
        b.charCount++
    }
}

// AppendByte adds a single byte (ASCII) to the builder.
func (b *Builder) AppendByte(c byte) {
    b.inner.WriteByte(c)
    b.charCount++
}

// AppendRune adds a single Unicode code point to the builder.
func (b *Builder) AppendRune(r rune) {
    b.inner.WriteRune(r)
    b.charCount++
}

// Len returns the number of Unicode characters accumulated so far.
func (b *Builder) Len() int64 {
    return b.charCount
}

// ByteLen returns the number of bytes accumulated so far.
func (b *Builder) ByteLen() int {
    return b.inner.Len()
}

// Build returns the final AstraString and resets the builder.
func (b *Builder) Build() AstraString {
    result := NewString(b.inner.String())
    b.inner.Reset()
    b.charCount = 0
    return result
}

// Reset clears the builder without building.
func (b *Builder) Reset() {
    b.inner.Reset()
    b.charCount = 0
}
```

Usage in Astra:

```astra
import string

fn build_report(items: List<string>) -> string {
    let sb = string.Builder.new()
    sb.append("=== Report ===\n")
    for i in 0..items.len() {
        sb.append(i.to_string())
        sb.append(". ")
        sb.append(items[i])
        sb.append("\n")
    }
    sb.append("=== End ===\n")
    return sb.build()
}

fn main() {
    let items = ["Build lexer", "Build parser", "Build codegen", "Build stdlib"]
    let report = build_report(items)
    print(report)
}
```

Output:
```
=== Report ===
0. Build lexer
1. Build parser
2. Build codegen
3. Build stdlib
=== End ===
```

The string builder is 100–1000x faster than `+=` concatenation for large strings because it avoids repeated allocation and copying.

---

## 4. String Memory Layout Diagram

```
STACK
┌────────────────────────────────────────┐
│  AstraString* s                        │  (8 bytes — pointer)
│  ──────────────────────────────────    │
│  value: 0x00007f8a4b200040 ──────────►─┼──┐
└────────────────────────────────────────┘  │
                                            │
HEAP                                        ▼
┌────────────────────────────────────────────────────────────────┐
│  AstraString struct @ 0x00007f8a4b200040                       │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ data:      0x00007f8a4b200070 ──────────────────────────►│──┼──┐
│  │ len:       12                                             │  │  │
│  │ byte_len:  20                                             │  │  │
│  │ ref_count: 1                                             │  │  │
│  │ _pad:      0                                             │  │  │
│  └──────────────────────────────────────────────────────────┘  │  │
└────────────────────────────────────────────────────────────────┘  │
                                                                     ▼
HEAP (string data)                                                   │
┌────────────────────────────────────────────────────────────────┐  │
│  char[] @ 0x00007f8a4b200070                                    │◄─┘
│  48 65 6C 6C 6F 2C 20 E4 B8 96 E7 95 8C 21 20 F0 9F 8C 8D 00  │
│  H  e  l  l  o  ,  SP 世         界         !  SP 🌍           \0│
└────────────────────────────────────────────────────────────────┘
```

---

## 5. The Math Package

The math package wraps Go's `math` standard library. Since Go's math functions are already IEEE 754-compliant and highly optimized, there is no need to write C wrappers for most of them.

```go
// stdlib/math/math.go

package astra_math

import (
    "math"
    "math/rand"
    "time"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
    PI    = math.Pi              // 3.141592653589793
    E     = math.E               // 2.718281828459045
    TAU   = math.Pi * 2          // 6.283185307179586 (full circle in radians)
    PHI   = 1.618033988749895    // golden ratio (1 + sqrt(5)) / 2
    SQRT2 = math.Sqrt2           // 1.4142135623730951
    LN2   = math.Ln2             // 0.6931471805599453
    LN10  = math.Log(10)         // 2.302585092994046
    INF   = math.Inf(1)          // positive infinity
    NaN   = math.NaN()           // not-a-number
)

// ─── Rounding ────────────────────────────────────────────────────────────────

// Floor returns the greatest integer value ≤ x.
// Floor(3.7) → 3.0,  Floor(-3.7) → -4.0
func Floor(x float64) float64 { return math.Floor(x) }

// Ceil returns the smallest integer value ≥ x.
// Ceil(3.2) → 4.0,  Ceil(-3.2) → -3.0
func Ceil(x float64) float64 { return math.Ceil(x) }

// Round returns the nearest integer, rounding half away from zero.
// Round(3.5) → 4.0,  Round(-3.5) → -4.0
func Round(x float64) float64 { return math.Round(x) }

// Trunc returns the integer part by removing the fractional part.
// Trunc(3.9) → 3.0,  Trunc(-3.9) → -3.0
func Trunc(x float64) float64 { return math.Trunc(x) }

// ─── Bounds ─────────────────────────────────────────────────────────────────

// Abs returns the absolute value of x.
func Abs(x float64) float64 { return math.Abs(x) }

// AbsInt returns the absolute value of an integer.
func AbsInt(x int64) int64 {
    if x < 0 { return -x }
    return x
}

// Min returns the smaller of a and b.
func Min(a, b float64) float64 { return math.Min(a, b) }

// MinInt returns the smaller of two integers.
func MinInt(a, b int64) int64 {
    if a < b { return a }
    return b
}

// Max returns the larger of a and b.
func Max(a, b float64) float64 { return math.Max(a, b) }

// MaxInt returns the larger of two integers.
func MaxInt(a, b int64) int64 {
    if a > b { return a }
    return b
}

// Clamp returns x clamped to the range [min, max].
// If x < min, returns min. If x > max, returns max. Otherwise returns x.
func Clamp(x, minVal, maxVal float64) float64 {
    if x < minVal { return minVal }
    if x > maxVal { return maxVal }
    return x
}

// ClampInt returns an integer clamped to [min, max].
func ClampInt(x, minVal, maxVal int64) int64 {
    if x < minVal { return minVal }
    if x > maxVal { return maxVal }
    return x
}

// ─── Powers and Roots ────────────────────────────────────────────────────────

// Sqrt returns the square root of x. Panics if x < 0.
func Sqrt(x float64) float64 { return math.Sqrt(x) }

// Cbrt returns the cube root of x (handles negative x correctly).
func Cbrt(x float64) float64 { return math.Cbrt(x) }

// Pow returns x raised to the power y.
func Pow(x, y float64) float64 { return math.Pow(x, y) }

// Exp returns e^x (Euler's number raised to x).
func Exp(x float64) float64 { return math.Exp(x) }

// Exp2 returns 2^x.
func Exp2(x float64) float64 { return math.Exp2(x) }

// ─── Logarithms ──────────────────────────────────────────────────────────────

// Log returns the natural logarithm (base e) of x. x must be > 0.
func Log(x float64) float64 { return math.Log(x) }

// Log2 returns the base-2 logarithm of x.
func Log2(x float64) float64 { return math.Log2(x) }

// Log10 returns the base-10 logarithm of x.
func Log10(x float64) float64 { return math.Log10(x) }

// LogBase returns log of x in the given base.
func LogBase(x, base float64) float64 {
    return math.Log(x) / math.Log(base)
}

// ─── Trigonometry ────────────────────────────────────────────────────────────
// All trig functions take/return values in radians.

func Sin(x float64) float64 { return math.Sin(x) }
func Cos(x float64) float64 { return math.Cos(x) }
func Tan(x float64) float64 { return math.Tan(x) }

// Asin returns the arcsine of x in radians. x must be in [-1, 1].
func Asin(x float64) float64 { return math.Asin(x) }

// Acos returns the arccosine of x in radians. x must be in [-1, 1].
func Acos(x float64) float64 { return math.Acos(x) }

// Atan returns the arctangent of x in radians.
func Atan(x float64) float64 { return math.Atan(x) }

// Atan2 returns the angle in radians between the positive X axis and
// the point (x, y). More numerically stable than Atan(y/x).
func Atan2(y, x float64) float64 { return math.Atan2(y, x) }

// Degrees converts radians to degrees.
func Degrees(rad float64) float64 { return rad * (180.0 / math.Pi) }

// Radians converts degrees to radians.
func Radians(deg float64) float64 { return deg * (math.Pi / 180.0) }

// ─── Hyperbolic ──────────────────────────────────────────────────────────────

func Sinh(x float64) float64 { return math.Sinh(x) }
func Cosh(x float64) float64 { return math.Cosh(x) }
func Tanh(x float64) float64 { return math.Tanh(x) }

// ─── Modular Arithmetic ──────────────────────────────────────────────────────

// Fmod returns the floating-point remainder of x/y.
// The result has the same sign as x.
func Fmod(x, y float64) float64 { return math.Mod(x, y) }

// IsNaN returns true if x is not-a-number.
func IsNaN(x float64) bool { return math.IsNaN(x) }

// IsInf returns true if x is positive or negative infinity.
func IsInf(x float64) bool { return math.IsInf(x, 0) }

// IsFinite returns true if x is a normal finite number.
func IsFinite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

// ─── Random Numbers ──────────────────────────────────────────────────────────

// Global random source. Seeded with current time at startup for
// non-deterministic behavior by default.
var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Seed sets the random seed for reproducible sequences.
// Use Seed(42) in tests to get deterministic behavior.
func Seed(s int64) {
    globalRand = rand.New(rand.NewSource(s))
}

// RandInt returns a random integer in the range [min, max] (inclusive).
// Panics if min > max.
func RandInt(min, max int64) int64 {
    if min > max {
        panic(fmt.Sprintf("math.rand_int: min (%d) > max (%d)", min, max))
    }
    return min + globalRand.Int63n(max-min+1)
}

// RandFloat returns a random float64 in [0.0, 1.0).
func RandFloat() float64 {
    return globalRand.Float64()
}

// RandFloat64 returns a random float64 in [min, max).
func RandFloat64(min, max float64) float64 {
    return min + globalRand.Float64()*(max-min)
}

// Shuffle randomly reorders a slice of int64s in-place.
// Uses the Fisher-Yates shuffle algorithm: O(n).
func ShuffleInts(arr []int64) {
    globalRand.Shuffle(len(arr), func(i, j int) {
        arr[i], arr[j] = arr[j], arr[i]
    })
}
```

Usage in Astra:

```astra
import math

fn main() {
    // Constants
    print(math.PI)    // 3.141592653589793
    print(math.TAU)   // 6.283185307179586
    print(math.PHI)   // 1.618033988749895

    // Area of a circle
    let radius = 5.0
    let area = math.PI * math.pow(radius, 2.0)
    print("Area: " + area.to_string())  // Area: 78.53981633974483

    // Degrees and radians
    let angle_deg = 45.0
    let angle_rad = math.radians(angle_deg)
    print(math.sin(angle_rad))   // 0.7071067811865476 (√2/2)
    print(math.cos(angle_rad))   // 0.7071067811865476

    // Distance between two points using Pythagorean theorem
    let dx = 3.0
    let dy = 4.0
    let dist = math.sqrt(math.pow(dx, 2.0) + math.pow(dy, 2.0))
    print("Distance: " + dist.to_string())  // Distance: 5

    // Clamping user input
    let user_volume = 150
    let safe_volume = math.clamp(user_volume, 0, 100)
    print(safe_volume)  // 100

    // Random numbers
    math.seed(42)
    for _ in 0..5 {
        print(math.rand_int(1, 6).to_string() + " ")  // dice roll
    }
    // Reproducible: 6 1 4 4 2 (always the same with seed 42)

    // Logarithms: how many digits does n have?
    let n = 1_000_000
    let digits = math.floor(math.log10(n.to_float()) + 1.0)
    print("Digits: " + digits.to_string())  // Digits: 7
}
```

---

## 6. The Time Package

Time is surprisingly complex: timezones, leap seconds, daylight saving, calendar systems. Astra's time package takes a pragmatic approach — it wraps Go's `time` package, which handles all of this correctly.

```go
// stdlib/time/time.go

package astra_time

import (
    "fmt"
    "time"
)

// ─── Types ────────────────────────────────────────────────────────────────────

// AstraTime wraps Go's time.Time.
type AstraTime struct {
    inner time.Time
}

// Duration represents a time interval in nanoseconds (same as Go's time.Duration).
type Duration struct {
    ns int64 // nanoseconds
}

// ─── Constructors ─────────────────────────────────────────────────────────────

// Now returns the current local time.
func Now() AstraTime {
    return AstraTime{inner: time.Now()}
}

// Unix returns the current Unix timestamp (seconds since Jan 1 1970 UTC).
func Unix() int64 {
    return time.Now().Unix()
}

// UnixMilli returns current time as milliseconds since Unix epoch.
func UnixMilli() int64 {
    return time.Now().UnixMilli()
}

// UnixNano returns current time as nanoseconds since Unix epoch.
func UnixNano() int64 {
    return time.Now().UnixNano()
}

// FromUnix creates an AstraTime from a Unix timestamp.
func FromUnix(ts int64) AstraTime {
    return AstraTime{inner: time.Unix(ts, 0)}
}

// ─── Time Fields ──────────────────────────────────────────────────────────────

func (t AstraTime) Year()       int { return t.inner.Year() }
func (t AstraTime) Month()      int { return int(t.inner.Month()) }
func (t AstraTime) Day()        int { return t.inner.Day() }
func (t AstraTime) Hour()       int { return t.inner.Hour() }
func (t AstraTime) Minute()     int { return t.inner.Minute() }
func (t AstraTime) Second()     int { return t.inner.Second() }
func (t AstraTime) Nanosecond() int { return t.inner.Nanosecond() }

// Weekday returns the day of the week: 0 = Sunday, 6 = Saturday.
func (t AstraTime) Weekday() int { return int(t.inner.Weekday()) }

// YearDay returns the day of the year, in the range [1, 365] (or 366).
func (t AstraTime) YearDay() int { return t.inner.YearDay() }

// UnixTimestamp returns the underlying Unix timestamp for this time.
func (t AstraTime) UnixTimestamp() int64 { return t.inner.Unix() }

// ─── Time Arithmetic ─────────────────────────────────────────────────────────

// Add returns a new AstraTime with the duration added.
func (t AstraTime) Add(d Duration) AstraTime {
    return AstraTime{inner: t.inner.Add(time.Duration(d.ns))}
}

// Sub returns the Duration between t and other (t - other).
func (t AstraTime) Sub(other AstraTime) Duration {
    return Duration{ns: int64(t.inner.Sub(other.inner))}
}

// Before returns true if t is before other.
func (t AstraTime) Before(other AstraTime) bool {
    return t.inner.Before(other.inner)
}

// After returns true if t is after other.
func (t AstraTime) After(other AstraTime) bool {
    return t.inner.After(other.inner)
}

// Equal returns true if t and other represent the same instant.
func (t AstraTime) Equal(other AstraTime) bool {
    return t.inner.Equal(other.inner)
}

// Since returns the Duration elapsed since t (equivalent to Now().Sub(t)).
func Since(t AstraTime) Duration {
    return Duration{ns: int64(time.Since(t.inner))}
}

// Until returns the Duration until t (equivalent to t.Sub(Now())).
func Until(t AstraTime) Duration {
    return Duration{ns: int64(time.Until(t.inner))}
}

// ─── Duration Constructors ───────────────────────────────────────────────────

func Nanoseconds(n int64)  Duration { return Duration{ns: n} }
func Microseconds(n int64) Duration { return Duration{ns: n * 1_000} }
func Milliseconds(n int64) Duration { return Duration{ns: n * 1_000_000} }
func Seconds(n int64)      Duration { return Duration{ns: n * 1_000_000_000} }
func Minutes(n int64)      Duration { return Duration{ns: n * 60 * 1_000_000_000} }
func Hours(n int64)        Duration { return Duration{ns: n * 3600 * 1_000_000_000} }
func Days(n int64)         Duration { return Duration{ns: n * 86400 * 1_000_000_000} }

// ─── Duration Fields ─────────────────────────────────────────────────────────

func (d Duration) AsNanoseconds()  int64   { return d.ns }
func (d Duration) AsMicroseconds() float64 { return float64(d.ns) / 1_000 }
func (d Duration) AsMilliseconds() float64 { return float64(d.ns) / 1_000_000 }
func (d Duration) AsSeconds()      float64 { return float64(d.ns) / 1_000_000_000 }
func (d Duration) AsMinutes()      float64 { return float64(d.ns) / 60_000_000_000 }
func (d Duration) AsHours()        float64 { return float64(d.ns) / 3_600_000_000_000 }

func (d Duration) String() string {
    return time.Duration(d.ns).String()
}

// ─── Sleep ───────────────────────────────────────────────────────────────────

// Sleep pauses the current goroutine for `ms` milliseconds.
func Sleep(ms int64) {
    time.Sleep(time.Duration(ms) * time.Millisecond)
}

// SleepDuration pauses for the given Duration.
func SleepDuration(d Duration) {
    time.Sleep(time.Duration(d.ns))
}

// ─── Formatting ──────────────────────────────────────────────────────────────

// Format formats a time using Go's reference time format.
// The format string uses Go's "magic reference time": Mon Jan 2 15:04:05 2006
// Common formats:
//   "2006-01-02"              → "2024-01-15"
//   "2006-01-02 15:04:05"     → "2024-01-15 14:30:00"
//   "2006-01-02T15:04:05Z07:00" → ISO 8601 / RFC 3339
//   "January 2, 2006"         → "January 15, 2024"
//   "Mon, 02 Jan 2006"        → "Mon, 15 Jan 2024"
//   "3:04 PM"                 → "2:30 PM"
//   "15:04:05.000"            → "14:30:00.000" (with milliseconds)
func Format(t AstraTime, layout string) string {
    return t.inner.Format(layout)
}

// FormatISO returns the time in ISO 8601 / RFC 3339 format.
func FormatISO(t AstraTime) string {
    return t.inner.Format(time.RFC3339)
}

// FormatDate returns "YYYY-MM-DD".
func FormatDate(t AstraTime) string {
    return t.inner.Format("2006-01-02")
}

// FormatDateTime returns "YYYY-MM-DD HH:MM:SS".
func FormatDateTime(t AstraTime) string {
    return t.inner.Format("2006-01-02 15:04:05")
}

// ─── Parsing ─────────────────────────────────────────────────────────────────

// Parse parses a time string using the given layout.
// Returns (AstraTime, true) on success or (zero, false) on failure.
func Parse(layout, value string) (AstraTime, bool) {
    t, err := time.Parse(layout, value)
    if err != nil {
        return AstraTime{}, false
    }
    return AstraTime{inner: t}, true
}

// ParseISO parses an ISO 8601 / RFC 3339 time string.
func ParseISO(s string) (AstraTime, bool) {
    t, err := time.Parse(time.RFC3339, s)
    if err != nil {
        return AstraTime{}, false
    }
    return AstraTime{inner: t}, true
}

// ParseDate parses a "YYYY-MM-DD" string.
func ParseDate(s string) (AstraTime, bool) {
    return Parse("2006-01-02", s)
}
```

Usage in Astra:

```astra
import time

fn main() {
    // Current time
    let now = time.now()
    print("Year:    " + now.year().to_string())
    print("Month:   " + now.month().to_string())
    print("Day:     " + now.day().to_string())
    print("Hour:    " + now.hour().to_string())
    print("Minute:  " + now.minute().to_string())

    // Formatting
    let formatted = time.format(now, "2006-01-02 15:04:05")
    print("Now: " + formatted)  // e.g., "2024-01-15 14:30:00"

    // Unix timestamp
    let ts = time.unix()
    print("Unix: " + ts.to_string())  // e.g., "1705329000"

    // Arithmetic
    let deadline = now.add(time.days(7))
    let remaining = deadline.sub(now)
    print("Days remaining: " + remaining.as_hours().to_string())

    // Parsing
    let birth = time.parse_date("1995-06-15")
    let age_duration = now.sub(birth)
    let age_years = (age_duration.as_hours() / 8760.0).trunc()
    print("Age: " + age_years.to_string() + " years")

    // Benchmarking
    let start = time.now()
    let sum = 0
    for i in 0..1_000_000 {
        sum = sum + i
    }
    let elapsed = time.since(start)
    print("Sum: " + sum.to_string())
    print("Time: " + elapsed.as_milliseconds().to_string() + "ms")

    // Sleep
    print("Sleeping 500ms...")
    time.sleep(500)
    print("Done!")
}
```

---

## 7. How Method Calls Compile — The Full Pipeline

Let's trace `"hello".to_upper()` through the entire compiler.

### Step 1 — Source Code

```astra
let result = "hello".to_upper()
```

### Step 2 — AST

The parser recognizes `"hello"` as a `StringLiteralExpr` and `.to_upper()` as a `MethodCallExpr`.

```
LetStatement
└── name: "result"
└── value: MethodCallExpr
    ├── receiver: StringLiteralExpr { value: "hello" }
    ├── method:   "to_upper"
    └── args:     []
```

### Step 3 — Semantic Analysis

The type checker resolves the receiver type as `string`, looks up `to_upper` in the string method table, and annotates the node with its return type `string`.

```go
// In semantic/type_checker.go

case *MethodCallExpr:
    receiverType := tc.inferType(node.Receiver)
    if receiverType == TypeString {
        method, ok := stringMethods[node.Method]
        if !ok {
            tc.error("string has no method '%s'", node.Method)
        }
        node.ResolvedType = method.ReturnType
        node.RuntimeFunc  = method.RuntimeFunc  // "astra_string_to_upper"
    }
```

The method table entry for `to_upper`:

```go
var stringMethods = map[string]MethodInfo{
    "to_upper": {
        Params:      []Type{},
        ReturnType:  TypeString,
        RuntimeFunc: "astra_string_to_upper",
    },
    "to_lower": {
        Params:      []Type{},
        ReturnType:  TypeString,
        RuntimeFunc: "astra_string_to_lower",
    },
    // ... all other methods
}
```

### Step 4 — IR Generation

The IR generator emits a call instruction with the resolved runtime function name:

```
; Astra IR for: let result = "hello".to_upper()

%0 = astra_string_const "hello"        ; load string constant
%1 = call astra_string_to_upper(%0)    ; call method
store result, %1                       ; store in variable
```

The IR call instruction looks like:

```go
// In ir/instructions.go

type CallInstr struct {
    Dest     string    // "%1"
    FuncName string    // "astra_string_to_upper"
    Args     []Value   // ["%0"]
}
```

### Step 5 — Code Generation (x86-64 Assembly)

The backend translates the IR call into x86-64 assembly following the System V AMD64 ABI. The first argument goes in `rdi`.

```asm
; Generated x86-64 for "hello".to_upper()

    ; Load the string constant (pointer to AstraString on heap)
    lea     rdi, [rip + .str_hello]     ; rdi = pointer to "hello" AstraString
    call    astra_string_to_upper       ; call C function
    ; result (AstraString*) is now in rax
    mov     [rbp - 8], rax             ; store into local variable "result"

.str_hello:
    ; AstraString struct for "hello" (BSS or .rodata section)
    .quad   .str_hello_data            ; data pointer
    .quad   5                          ; len (5 Unicode chars)
    .quad   5                          ; byte_len (5 bytes, pure ASCII)
    .long   -1                         ; ref_count = -1 means "don't free" (static)
    .long   0                          ; _pad

.str_hello_data:
    .byte   0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x00   ; "hello\0"
```

### Step 6 — C Runtime Function

The actual `astra_string_to_upper` function in the C runtime:

```c
// runtime/astra_string.c

#include "astra_string.h"
#include <ctype.h>
#include <wctype.h>
#include <locale.h>

// For a production implementation we would use a full Unicode case table.
// This simplified version handles ASCII and common Latin characters.
AstraString* astra_string_to_upper(AstraString* s) {
    if (!s || s->byte_len == 0) {
        return astra_string_empty();
    }

    // Allocate output buffer — upper-case conversion cannot increase UTF-8
    // byte length for most scripts (rare exceptions exist in Unicode).
    // We allocate 4x for safety (worst-case Unicode expansion).
    char* buf = (char*)malloc(s->byte_len * 4 + 1);
    if (!buf) return NULL;

    int64_t in  = 0;
    int64_t out = 0;

    while (in < s->byte_len) {
        unsigned char b = (unsigned char)s->data[in];

        if (b < 0x80) {
            // ASCII: simple lookup table
            buf[out++] = (char)toupper(b);
            in += 1;
        } else {
            // Non-ASCII: use the wide character functions
            // Decode one UTF-8 code point
            wchar_t wc = 0;
            int char_bytes = utf8_char_len(b);

            // Simplified: copy non-ASCII unchanged for now.
            // A production implementation would use a full Unicode case table
            // (ICU library or similar).
            for (int i = 0; i < char_bytes && in + i < s->byte_len; i++) {
                buf[out++] = s->data[in + i];
            }
            in += char_bytes;
        }
    }

    buf[out] = '\0';
    AstraString* result = astra_string_new(buf, out);
    free(buf);
    return result;
}
```

### The Complete Pipeline Diagram

```mermaid
flowchart TD
    SRC["Source Code: \"hello\".to_upper()"]
    LEX["Lexer\nToken(STRING), Token(DOT),\nToken(IDENT to_upper), Token(LPAREN), Token(RPAREN)"]
    PAR["Parser\nMethodCallExpr { receiver: StringLiteral, method: to_upper, args: [] }"]
    SEM["Semantic Analysis\nResolve receiver → string\nResolve method → astra_string_to_upper\nInfer return type → string"]
    IRGEN["IR Generation\n%0 = string_const hello\n%1 = call astra_string_to_upper(%0)"]
    CODEGEN["Code Generation x86-64\nlea rdi, [rip + .str_hello]\ncall astra_string_to_upper\nmov [rbp-8], rax"]
    RT["Runtime C\nastra_string_to_upper(AstraString* s)"]
    RES["Result: AstraString* pointing to HELLO"]
    SRC --> LEX --> PAR --> SEM --> IRGEN --> CODEGEN --> RT --> RES
```

---

## 8. Build Milestone — Complete stdlib/string, math, time

### File Structure

```
stdlib/
├── string/
│   ├── string.go         (NewString, ToUpper, ToLower, Contains, Slice, ...)
│   └── builder.go        (Builder, Append, Build, ...)
├── math/
│   └── math.go           (PI, E, TAU, Sqrt, Sin, Cos, RandInt, ...)
└── time/
    └── time.go           (Now, Unix, Sleep, Format, Parse, ...)

runtime/
├── astra_string.h        (AstraString struct, function declarations)
├── astra_string.c        (UTF-8 implementation)
├── astra_math.h          (thin wrappers around libm)
└── astra_time.h          (thin wrappers around <time.h>)
```

### Verification Test

```astra
// tests/stdlib_core_test.as

import string
import math
import time

fn test_strings() {
    let s = "Hello, 世界! 🌍"
    assert(s.length() == 12, "length should be 12")
    assert(s.byte_length() == 20, "byte_length should be 20")
    assert(s.to_upper() == "HELLO, 世界! 🌍", "to_upper failed")
    assert(s.to_lower() == "hello, 世界! 🌍", "to_lower failed")
    assert(s.contains("世界") == true, "contains failed")
    assert(s.starts_with("Hello") == true, "starts_with failed")
    assert(s.ends_with("🌍") == true, "ends_with failed")
    assert(s.index_of("世") == 7, "index_of failed")
    assert(s.slice(0, 5) == "Hello", "slice failed")
    assert(s.replace("Hello", "Hi") == "Hi, 世界! 🌍", "replace failed")
    assert("  hello  ".trim() == "hello", "trim failed")
    assert("ha".repeat(3) == "hahaha", "repeat failed")
    assert("42".parse_int() == 42, "parse_int failed")
    print("All string tests passed!")
}

fn test_math() {
    assert(math.abs(-5.0) == 5.0, "abs failed")
    assert(math.floor(3.7) == 3.0, "floor failed")
    assert(math.ceil(3.2) == 4.0, "ceil failed")
    assert(math.round(3.5) == 4.0, "round failed")
    assert(math.sqrt(9.0) == 3.0, "sqrt failed")
    assert(math.pow(2.0, 10.0) == 1024.0, "pow failed")
    assert(math.clamp(150, 0, 100) == 100, "clamp failed")
    assert(math.min(3.0, 5.0) == 3.0, "min failed")
    assert(math.max(3.0, 5.0) == 5.0, "max failed")
    math.seed(42)
    let r = math.rand_int(1, 10)
    assert(r >= 1 && r <= 10, "rand_int out of range")
    print("All math tests passed!")
}

fn test_time() {
    let now = time.now()
    assert(now.year() >= 2024, "year too small")
    assert(now.month() >= 1 && now.month() <= 12, "month out of range")
    assert(now.day() >= 1 && now.day() <= 31, "day out of range")
    let ts = time.unix()
    assert(ts > 1_000_000_000, "unix ts too small")
    let formatted = time.format(now, "2006-01-02")
    assert(formatted.length() == 10, "formatted date length wrong")
    let start = time.now()
    time.sleep(10)
    let elapsed = time.since(start)
    assert(elapsed.as_milliseconds() >= 9.0, "sleep too short")
    print("All time tests passed!")
}

fn main() {
    test_strings()
    test_math()
    test_time()
    print("All stdlib core tests passed!")
}
```

---

## 🔨 Astra Build Milestone

At this point, your Astra installation should have working implementations of the three core stdlib packages. Here is the expected project state:

| File | Lines | Status |
|------|-------|--------|
| `stdlib/string/string.go` | ~350 | Complete |
| `stdlib/string/builder.go` | ~80 | Complete |
| `stdlib/math/math.go` | ~200 | Complete |
| `stdlib/time/time.go` | ~180 | Complete |
| `runtime/astra_string.h` | ~60 | Complete |
| `runtime/astra_string.c` | ~300 | Complete |

Build and run the verification test:

```bash
$ cd making-own-language/
$ astrac build tests/stdlib_core_test.as -o test_core
$ ./test_core
All string tests passed!
All math tests passed!
All time tests passed!
All stdlib core tests passed!
```

---

## 9. Exercises

**Exercise 1 — String Padding**
Implement `s.pad_left(width, char)` and `s.pad_right(width, char)`. For example, `"42".pad_left(5, '0')` should return `"00042"`. Handle the case where `s.length() >= width` (return `s` unchanged). Add this to `stdlib/string/string.go`.

**Exercise 2 — String Split with Limit**
Implement `s.split_n(sep, n)` that splits into at most `n` parts. `"a,b,c,d".split_n(",", 2)` should return `["a", "b,c,d"]`. The last element gets all remaining content.

**Exercise 3 — Math: Is Prime?**
Implement `math.is_prime(n: int) -> bool` using the trial division algorithm. Then implement the Sieve of Eratosthenes as `math.primes_up_to(n: int) -> List<int>` to generate all primes below n. Benchmark both for n = 1,000,000.

**Exercise 4 — Time Formatting Table**
Write an Astra program that prints the current time in 10 different formats (ISO 8601, Unix timestamp, RFC 1123, "January 2 2006", "Mon Jan 2", etc.). Use the format strings from the time package documentation.

**Exercise 5 — Stopwatch**
Implement a `Stopwatch` type in Astra using the time package:
```astra
let sw = Stopwatch.new()
sw.start()
// ... some work ...
sw.stop()
print(sw.elapsed().as_milliseconds().to_string() + "ms")
sw.reset()
```
Support `lap()` that records intermediate times without stopping.

**Exercise 6 — String Format Extension**
Extend the `format` method to support named placeholders: `"Hello, {name}! You are {age} years old.".format_named({"name": "Aditya", "age": "25"})`. The implementation should take a `Map<string, string>` and replace `{key}` with the corresponding value.

---

## Summary

| Topic | Key Takeaway |
|-------|-------------|
| AstraString struct | `len` = Unicode chars, `byte_len` = UTF-8 bytes. These differ for non-ASCII text. |
| UTF-8 | Self-synchronizing encoding. ASCII-compatible. `byte_len >= len` always. |
| String methods | All return new strings (immutable). Operate in character space, not byte space. |
| String Builder | Use `string.Builder` for any loop that concatenates strings. O(n) vs O(n²). |
| Math package | Wraps Go's `math` package. Constants: PI, E, TAU, PHI, SQRT2. |
| Random numbers | `math.seed()` for reproducibility. `rand_int(min, max)` inclusive on both ends. |
| Time package | Wraps Go's `time` package. Go's "magic reference time" format strings. |
| Method call pipeline | Source → AST (MethodCallExpr) → IR (call) → Assembly (call) → C runtime function. |
