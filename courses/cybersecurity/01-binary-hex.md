# Chapter 01: Binary, Hex, and Why Computers Speak That Language

*When you look at a memory dump or network packet, you'll see numbers that look nothing like regular numbers. This chapter teaches you binary and hexadecimal — essential for reading raw computer data.*

---

## Why Binary?

Computers are built from transistors — microscopic switches that are either ON or OFF.

ON = 1  
OFF = 0

That's binary. Two states. Everything a computer does is built from these two states.

Why not use decimal (base 10)? Because it's much harder to build a physical switch that can reliably be in 10 distinct states. But making it just ON or OFF is easy and reliable. So computers use base 2 (binary).

---

## Binary — Base 2

In decimal (base 10), each digit position is a power of 10:

```
1,234 = 1×1000 + 2×100 + 3×10 + 4×1
      = 1×10³  + 2×10² + 3×10¹ + 4×10⁰
```

In binary (base 2), each digit position is a power of 2:

```
1011 = 1×8 + 0×4 + 1×2 + 1×1
     = 1×2³ + 0×2² + 1×2¹ + 1×2⁰
     = 8 + 0 + 2 + 1
     = 11 (decimal)
```

The binary digits are called "bits". 8 bits = 1 byte.

### Binary to Decimal Conversion

Position values (right to left):
```
Bit position: 7  6  5  4  3  2  1  0
Value:       128 64 32 16  8  4  2  1
```

Example: convert binary `10110101` to decimal:
```
1×128 + 0×64 + 1×32 + 1×16 + 0×8 + 1×4 + 0×2 + 1×1
= 128 + 0 + 32 + 16 + 0 + 4 + 0 + 1
= 181
```

### Decimal to Binary Conversion

Repeatedly divide by 2 and collect remainders:
```
181 ÷ 2 = 90 remainder 1
 90 ÷ 2 = 45 remainder 0
 45 ÷ 2 = 22 remainder 1
 22 ÷ 2 = 11 remainder 0
 11 ÷ 2 =  5 remainder 1
  5 ÷ 2 =  2 remainder 1
  2 ÷ 2 =  1 remainder 0
  1 ÷ 2 =  0 remainder 1
```
Read remainders bottom to top: `10110101`

### Practice

Convert these to decimal:
- `00001010` → ?
- `11111111` → ?
- `01000001` → ? (hint: this is the letter 'A' in ASCII)

Answers: 10, 255, 65

---

## Hexadecimal — Base 16

Binary is precise but tedious to write. `11111010000011001101011011001010` is hard to read.

Hexadecimal (hex) solves this. It's base 16, using digits 0-9 and letters A-F:

```
0=0, 1=1, 2=2, 3=3, 4=4, 5=5, 6=6, 7=7
8=8, 9=9, A=10, B=11, C=12, D=13, E=14, F=15
```

**The key insight:** Every 4 bits of binary converts to exactly 1 hex digit.

```
Binary: 1111 1010 0000 1100 1101 0110 1100 1010
Hex:      F    A    0    C    D    6    C    A
Result: 0xFA0CD6CA
```

This is why memory addresses, IP addresses (internal representation), file signatures, and cryptographic hashes are shown in hex.

### Hex to Decimal

```
0xFF = 15×16 + 15 = 240 + 15 = 255
0x41 = 4×16 + 1  = 64 + 1   = 65 (letter 'A' again)
0x7F = 7×16 + 15 = 112 + 15 = 127
```

### Binary to Hex (Easy Method)

Group binary into groups of 4 bits, convert each group:
```
Binary: 10101100 11001000
Groups: 1010 1100  1100 1000
Hex:      A    C     C    8
Result: 0xACC8
```

### Common Hex Values You'll See

| Hex | Decimal | Context |
|-----|---------|---------|
| `0xFF` | 255 | Max byte value, subnet mask components |
| `0x00` | 0 | Null byte, often marks end of strings |
| `0x41` | 65 | ASCII 'A' — found in buffer overflow shellcode |
| `0x90` | 144 | NOP instruction (no-operation) — used in exploits |
| `0x0A` | 10 | Newline character |
| `0x7F` | 127 | Loopback IP (127.0.0.1) |
| `0xDEADBEEF` | — | Debugging marker, famous in security |

### The `0x` Prefix

`0x` before a number means it's hexadecimal. So:
- `0x41` is hex 41 (decimal 65)
- `41` is decimal 41
- `0b01000001` is binary (the `b` notation)

---

## How Text Is Stored — ASCII and Unicode

Computers store text as numbers. The mapping between numbers and characters is called an encoding.

**ASCII (American Standard Code for Information Interchange):**
- 7-bit encoding, supports 128 characters
- Covers English letters, digits, punctuation, control characters

```
Decimal  Hex   Character
65       0x41  A
66       0x42  B
97       0x61  a
48       0x30  0 (digit zero)
32       0x20  Space
10       0x0A  Newline
0        0x00  Null
```

**Why null matters:** In C (the language most operating systems are written in), strings end with a null byte (`0x00`). Many buffer overflow attacks use null bytes or exploit assumptions about where strings end.

**Unicode / UTF-8:**
- Supports all world languages (1.1 million characters)
- ASCII is a subset of UTF-8
- Most modern systems use UTF-8

---

## How Numbers Are Stored — Integers

A signed 32-bit integer can hold values from -2,147,483,648 to 2,147,483,647.  
An unsigned 32-bit integer holds 0 to 4,294,967,295.

**Little-endian vs Big-endian:**

This is how multi-byte numbers are stored in memory:
- **Big-endian:** Most significant byte first (like writing numbers normally)
- **Little-endian:** Least significant byte first (x86/x64 CPUs use this)

The number `0x12345678` in memory:
- Big-endian: `12 34 56 78`
- Little-endian: `78 56 34 12`

**Why this matters:** When analyzing network packets and memory dumps, you need to know byte order to correctly read multi-byte values. TCP/IP uses big-endian ("network byte order"). x86 CPUs use little-endian.

---

## IP Addresses in Binary

IPv4 addresses are 32-bit numbers written in "dotted decimal notation":

```
192.168.1.100

192 = 11000000
168 = 10101000
  1 = 00000001
100 = 01100100

Binary: 11000000.10101000.00000001.01100100
Hex:       C0       A8      01      64
```

So `192.168.1.100` is the same as `0xC0A80164` in hex.

**Subnet masks:**
`255.255.255.0` = `0xFFFFFF00` = `11111111.11111111.11111111.00000000`

The 1s indicate "network part", 0s indicate "host part". CIDR notation `/24` means 24 bits are the network part.

---

## Bitwise Operations

Security tools frequently use bitwise operations — operations performed on individual bits.

**AND (`&`):** Both bits must be 1
```
01001010  (74)
& 11110000  (240)
= 01000000  (64)
```

**OR (`|`):** At least one bit must be 1  
**XOR (`^`):** Bits must be different (1 if different, 0 if same)  
**NOT (`~`):** Flip all bits  
**Left shift (`<<`):** Shift bits left (multiplies by power of 2)  
**Right shift (`>>`):** Shift bits right (divides by power of 2)

**Security use cases:**
- Subnet masking: `IP & subnet_mask` extracts the network address
- Flag checking: `flags & 0x02` checks if bit 1 is set
- Simple XOR encryption: `data ^ key` — used in many malware samples
- Permission bits: Linux file permissions are bit flags

---

## Reading a Hex Dump

When you analyze malware, memory, or network traffic, you'll see hex dumps like this:

```
00000000  4d 5a 90 00 03 00 00 00  04 00 00 00 ff ff 00 00  |MZ..............|
00000010  b8 00 00 00 00 00 00 00  40 00 00 00 00 00 00 00  |........@.......|
00000020  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
```

Format:
- **Left column:** Memory address (hex)
- **Middle:** 16 bytes per row in hex
- **Right:** ASCII representation (`.` for non-printable)

**The `4d 5a` / `MZ` header:** Every Windows executable file starts with bytes `0x4D 0x5A` — "MZ" in ASCII (Magic Ziggy, after Mark Zbikowski who designed the DOS executable format). When you see this in a memory dump, you've found a Windows PE executable.

**File signatures (magic bytes):** Every file format starts with specific bytes:
- `4D 5A` — Windows executable (.exe, .dll)
- `7F 45 4C 46` — Linux ELF executable
- `FF D8 FF` — JPEG image
- `89 50 4E 47` — PNG image
- `25 50 44 46` — PDF document
- `50 4B 03 04` — ZIP file (also Office documents!)

Malware often disguises itself by faking file extensions. But you can always check the magic bytes to find the real type.

---

## Go: Working with Binary and Hex

Here's how you work with hex and binary in Go:

```go
package main

import "fmt"

func main() {
    // Hex literals
    var address uint32 = 0xC0A80164  // 192.168.1.100
    fmt.Printf("IP as decimal: %d\n", address)
    fmt.Printf("IP as hex: 0x%X\n", address)
    
    // Binary literals (Go 1.13+)
    var perms uint8 = 0b00110110  // Unix permission bits
    fmt.Printf("Permissions: %d (octal: %o)\n", perms, perms)
    
    // Bitwise AND — subnet masking
    ip      := uint32(0xC0A80164)  // 192.168.1.100
    mask    := uint32(0xFFFFFF00)  // 255.255.255.0
    network := ip & mask
    fmt.Printf("Network: 0x%X\n", network)  // 0xC0A80100 = 192.168.1.0
    
    // XOR — simple cipher
    data    := byte(0x41)  // 'A'
    key     := byte(0xFF)
    encoded := data ^ key
    decoded := encoded ^ key
    fmt.Printf("Original: 0x%X, Encoded: 0x%X, Decoded: 0x%X\n", 
               data, encoded, decoded)
    
    // Converting between bytes and strings
    hexString := fmt.Sprintf("%X", []byte("Hello"))
    fmt.Println("Hello in hex:", hexString)
}
```

Output:
```
IP as decimal: 3232235876
IP as hex: 0xC0A80164
Permissions: 54 (octal: 66)
Network: 0xC0A80100
Original: 0x41, Encoded: 0xBE, Decoded: 0x41
Hello in hex: 48656C6C6F
```

---

## Summary

| Concept | Example | Why it matters |
|---------|---------|----------------|
| Binary | `01000001` | How data is physically stored |
| Hexadecimal | `0x41` = 65 = 'A' | Reading memory, packets, files |
| ASCII | 65 = 'A' | Interpreting text in raw data |
| Endianness | `0x1234` stored as `34 12` in little-endian | Reading multi-byte values correctly |
| Magic bytes | `4D 5A` = Windows executable | Identifying file types regardless of extension |
| XOR | `A ^ B ^ B = A` | Simple cipher, parity checks, checksums |

---

## Exercises

1. Convert `10011010` (binary) to decimal, then to hex
2. What character does `0x61` represent in ASCII?
3. Convert IP address `10.0.0.1` to hex
4. What's `0xFF & 0x0F`? What operation would you use this for?
5. Run the Go code above and modify it to XOR the string "password" with key `0xAA`. What do you get?
