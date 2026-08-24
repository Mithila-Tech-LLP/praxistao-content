# Chapter 43: CTF Methodology — Competing in Capture the Flag

*CTF (Capture the Flag) competitions are how security professionals sharpen skills in a legal, controlled environment. They're puzzles that teach real exploitation techniques.*

---

## What is a CTF?

CTF competitions give you challenges with a hidden "flag" (a string like `flag{abc123def456}`). You find the flag by exploiting vulnerabilities, solving crypto puzzles, reversing code, or analyzing forensics.

```
CTF Types:
├── Jeopardy         — categories of independent challenges, pick any
│   ├── Web          — SQL injection, XSS, SSRF, auth bypasses
│   ├── Pwn          — binary exploitation (buffer overflow, ROP)
│   ├── Reversing    — reverse engineer compiled code
│   ├── Crypto       — break cryptographic implementations
│   ├── Forensics    — analyze disk images, memory dumps, network captures
│   ├── OSINT        — open source intelligence
│   └── Misc         — misc puzzles, steg, trivia
│
└── Attack/Defense   — teams attack each other's services and defend their own
```

---

## Setting Up Your CTF Environment

```bash
# Kali Linux (recommended)
# - Full-featured pentest distro
# - All tools pre-installed

# Essential tools
apt install -y pwntools gdb python3-pwntools
pip3 install pycryptodome angr z3-solver

# For reversing
apt install ghidra radare2 gdb-peda

# For forensics
apt install foremost binwalk volatility3 autopsy

# For web
apt install burpsuite nikto sqlmap gobuster

# Docker for isolated environments
docker pull ctfd/ctfd    # self-hosted CTF platform
```

---

## Web CTF Methodology

```
Step 1: Reconnaissance
- Read the challenge description carefully
- View page source (Ctrl+U)
- Check HTML comments: <!-- debug mode: true -->
- Check JavaScript for hints
- Look for hidden form fields
- Check robots.txt, sitemap.xml, .git/

Step 2: Identify the vulnerability type
- SQL injection? Try ' or "
- XSS? Try <script>alert(1)</script>
- SSRF? Try URL parameters
- LFI? Try ?page=../../../../etc/passwd
- Command injection? Try ; id or | ls

Step 3: Exploit
- Manual first, tools second
- Burp Suite for intercepting/modifying requests
```

```bash
# Quick web recon
curl -sv "http://challenge.com/" 2>&1 | grep -i "server\|x-powered\|cookie"
curl "http://challenge.com/robots.txt"
gobuster dir -u http://challenge.com -w /usr/share/wordlists/dirbuster/medium.txt

# SQLi test
sqlmap -u "http://challenge.com/item?id=1" --batch --dump

# LFI
curl "http://challenge.com/?file=../../../../etc/passwd"
curl "http://challenge.com/?file=php://filter/convert.base64-encode/resource=index.php"
```

---

## Binary Exploitation (Pwn) Methodology

```bash
# 1. Understand the binary
file ./challenge        # what type of binary?
checksec --file=./challenge   # what protections?
strings ./challenge     # any interesting strings?

# 2. Run it and understand behavior
./challenge
python3 -c "print('A'*100)" | ./challenge

# 3. Open in GDB
gdb ./challenge
gdb-peda $ run          # run with peda
gdb-peda $ pattern create 200  # cyclic pattern
gdb-peda $ run <<< $(python3 -c "from pwn import *; print(cyclic(200))")
gdb-peda $ pattern offset $rip  # find offset after crash

# 4. Disassemble
gdb $ disas main
gdb $ disas vulnerable_function

# 5. Exploit
from pwn import *
p = process('./challenge')
# or remote: p = remote('challenge.ctf.com', 1337)
```

---

## Reverse Engineering Methodology

```bash
# 1. Static analysis
file ./challenge
strings ./challenge | grep -i flag
objdump -d ./challenge | head -100

# 2. Ghidra (GUI decompiler)
# Open binary → analyze → view decompiled C code
# Look for: strcmp, flag[], password checks

# 3. Dynamic analysis
ltrace ./challenge    # library calls (strcmp, strcmp comparisons!)
strace ./challenge    # system calls
gdb + breakpoints

# 4. Common patterns
# - strcmp with flag: run with the right argument
# - XOR encoding: find key
# - Base64 decoded flag: decode it
# - Custom encoding: reverse the algorithm

# Example: find flag comparison
strings ./challenge | grep -i "correct\|wrong\|flag\|CTF"
# If it says "Enter password:" then "Correct!" — it's comparing your input
```

---

## Cryptography CTF Methodology

```python
# Common crypto challenges:

# 1. Caesar cipher / ROT13
import codecs
codecs.encode('uryyb jbeyq', 'rot_13')

# 2. Base64 decode
import base64
base64.b64decode('aGVsbG8gd29ybGQ=')

# 3. XOR cipher (key reuse / short key)
from pwn import xor
ciphertext = bytes.fromhex('1a3f2b4c')
xor(ciphertext, b'key')

# 4. RSA with small N (factor it)
from sympy import factorint
n = 24892460398...
factors = factorint(n)  # if n is small enough

# 5. Hash cracking
import hashlib
for word in open('rockyou.txt'):
    if hashlib.md5(word.strip().encode()).hexdigest() == target_hash:
        print(word.strip())
```

---

## Forensics Methodology

```bash
# 1. File analysis
file mystery.dat         # what kind of file is this really?
binwalk mystery.dat      # embedded files inside?
hexdump -C mystery.dat | head   # look at raw bytes

# 2. Image forensics (steganography)
exiftool image.jpg           # metadata in EXIF
strings image.jpg            # strings embedded in image
steghide extract -sf image.jpg -p ""   # try empty password
zsteg image.png              # PNG steganography

# 3. PCAP analysis
wireshark capture.pcap
tshark -r capture.pcap -Y "http" -T fields -e http.file_data | xxd
# Look for: POST data, file transfers, credentials

# 4. Memory forensics
volatility3 -f memory.raw imageinfo
volatility3 -f memory.raw windows.pslist
volatility3 -f memory.raw windows.cmdline
volatility3 -f memory.raw windows.filescan | grep flag
```

---

## CTF Platforms

```
PicoCTF           — beginner-friendly, permanent challenges
HackTheBox        — harder, machine-based hacking
TryHackMe         — guided rooms, beginner to advanced
CTFtime.org       — list of all upcoming CTF competitions
pwn.college       — binary exploitation focus
CryptoHack        — crypto-only challenges
PortSwigger Labs  — web security labs
```

---

## Go: CTF Challenge Solver Framework

```go
package main

import (
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "strings"
)

// Common CTF decoding operations
func tryDecodings(data string) {
    fmt.Printf("Input: %s\n\n", data)
    
    // Base64
    if decoded, err := base64.StdEncoding.DecodeString(data); err == nil {
        fmt.Printf("Base64: %s\n", decoded)
    }
    
    // Hex
    if decoded, err := hex.DecodeString(data); err == nil {
        fmt.Printf("Hex: %s\n", decoded)
    }
    
    // ROT13
    rot13 := strings.Map(func(r rune) rune {
        switch {
        case r >= 'A' && r <= 'Z':
            return 'A' + (r-'A'+13)%26
        case r >= 'a' && r <= 'z':
            return 'a' + (r-'a'+13)%26
        }
        return r
    }, data)
    fmt.Printf("ROT13: %s\n", rot13)
    
    // XOR with single byte (brute force key)
    fmt.Println("\nXOR brute force (single byte key):")
    raw, _ := hex.DecodeString(data)
    for key := byte(0); key < 255; key++ {
        result := make([]byte, len(raw))
        for i, b := range raw {
            result[i] = b ^ key
        }
        s := string(result)
        if strings.Contains(s, "flag{") || strings.Contains(s, "CTF{") {
            fmt.Printf("  key=0x%02x: %s\n", key, s)
        }
    }
}

func main() {
    // Example: XOR encoded flag
    tryDecodings("666c61677b6578616d706c655f666c61677d")  // hex
}
```

---

## Summary

| CTF Category | Key Skill | Essential Tools |
|-------------|-----------|----------------|
| Web | Bug identification, HTTP | Burp Suite, sqlmap, gobuster |
| Pwn | Memory exploitation, GDB | pwntools, GDB+peda, checksec |
| Reversing | Disassembly, decompilation | Ghidra, ltrace, strings |
| Crypto | Math, encoding | Python, pycryptodome |
| Forensics | File analysis, steg | Wireshark, volatility, binwalk |

---

## Exercises

1. Complete 5 "easy" web challenges on PicoCTF — document your methodology for each
2. Solve a buffer overflow challenge on pwn.college — what's the offset and what payload works?
3. Use Ghidra to reverse a simple "crackme" binary and find the correct password
4. Analyze a PCAP file in Wireshark — find credentials or flag-containing HTTP traffic
