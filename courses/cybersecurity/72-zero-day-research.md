# Chapter 72: Zero-Day Research — Finding Unknown Vulnerabilities

*A zero-day is a vulnerability unknown to the vendor — no patch exists. Zero-day research is the apex of offensive security: you're finding bugs that nobody else knows about.*

---

## What Makes a Zero-Day?

```
Day 0     → Researcher discovers vulnerability
Day 0     → (if malicious) Attacker exploits in the wild
Day N     → Vendor notified OR vulnerability disclosed publicly
Day N     → Vendor aware: now it's just a "vulnerability"
Day N+X   → Patch released: now it's fully "patched"

"N-day" vulnerabilities have a known CVE but no patch yet (or organizations haven't patched yet)
"Zero-day" — the original undisclosed discovery
```

---

## Vulnerability Research Methodology

```
1. Target Selection
   - Open source software (can read source code)
   - Widely deployed software (higher impact)
   - Software with trust boundaries (parser inputs, network protocols)
   
2. Attack Surface Mapping
   - All entry points that accept untrusted input
   - Parsers (JSON, XML, image formats, protocols)
   - IPC mechanisms, sockets, shared memory
   
3. Code Audit / Fuzzing
   - Manual code review (for open source)
   - Static analysis tools
   - Dynamic fuzzing
   
4. Root Cause Analysis
   - Is this exploitable?
   - What's the worst-case impact?
   
5. Proof of Concept
   - Demonstrate exploitability
   - Minimal reliable reproducer
   
6. Responsible Disclosure
   - Report to vendor first
   - Give time to patch (90 days standard)
   - Publish after patch released
```

---

## Fuzzing — Automated Bug Finding

Fuzzing: generate random/mutated inputs, look for crashes:

```bash
# AFL++ (American Fuzzy Lop) — coverage-guided fuzzer
# Instruments binary to track code coverage
# Evolves inputs that reach new code paths

# Compile target with AFL instrumentation
CC=afl-clang-fast ./configure
make

# Create initial corpus (starting inputs)
mkdir corpus
echo "hello" > corpus/test1.txt

# Fuzz!
afl-fuzz -i corpus -o output ./target_binary @@
# @@ = placeholder for the mutated input file

# Libfuzzer (compile-into approach)
# Better for library functions

# Go fuzzing (built-in since Go 1.18!)
# Example fuzz test:
```

```go
// go test -fuzz=FuzzParseRequest -fuzztime=60s

package mypackage

import "testing"

func FuzzParseRequest(f *testing.F) {
    // Seed corpus
    f.Add([]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"))
    f.Add([]byte(""))
    f.Add([]byte("AAAAAAAAAAAAAAAAAAAAAAAAA"))
    
    f.Fuzz(func(t *testing.T, data []byte) {
        // ParseRequest should never panic or crash
        // regardless of input
        defer func() {
            if r := recover(); r != nil {
                t.Errorf("ParseRequest panicked: %v", r)
            }
        }()
        
        _ = ParseRequest(data)  // function under test
    })
}
```

---

## Common Vulnerability Classes to Hunt

### Memory Safety

```c
// Buffer overflow — check:
// - strcpy, strcat, sprintf, gets (unbounded)
// - memcpy without length check
// - Array access without bounds check

// Use-after-free — check:
// - free() then access same pointer
// - Dangling pointers in complex state machines

// Integer overflow:
// - size_t arithmetic → negative result used as allocation size
uint32_t size = user_input;
char *buf = malloc(size + 1);  // if size=0xFFFFFFFF → size+1=0 → malloc(0)
```

### Logic Bugs

```
Authentication bypasses:
- State machine flaws (skip authentication step)
- Type confusion (0 == false == null → passes check)
- Unicode normalization bypasses (admin → ádmin passes filter)

Authorization bypasses:
- IDOR (predictable IDs)
- Privilege checks in wrong place
- Race conditions in permission checks

Race conditions:
- TOCTOU (Time of Check, Time of Use)
  Check: is file safe? → Use: process file
  Between check and use: attacker swaps file
  /tmp file races are classic
```

---

## Static Analysis for Bug Finding

```bash
# Semgrep — pattern-based code scanner
# Free, open source, fast

# Scan for dangerous C functions
semgrep --config p/c.dangerous-functions ./src/

# Scan Go for security issues
semgrep --config p/golang-security ./src/

# Custom rule: find SQL concatenation
cat > sqli.yaml << 'EOF'
rules:
- id: go-sql-injection
  pattern: |
    $DB.Query("..." + $USERINPUT)
  message: Possible SQL injection via string concatenation
  severity: ERROR
  languages: [go]
EOF
semgrep --config sqli.yaml ./src/

# CodeQL (GitHub's semantic analysis)
# More powerful, requires setup
codeql database create mydb --language=go
codeql database analyze mydb go-security-and-quality.qls
```

---

## Go: Simple Fuzzer for Protocol Parsers

```go
package main

import (
    "bytes"
    "crypto/rand"
    "encoding/binary"
    "fmt"
    "math/big"
    "time"
)

type FuzzResult struct {
    Input     []byte
    Panicked  bool
    Duration  time.Duration
    ErrorMsg  string
}

func fuzzParser(parser func([]byte) error, corpus [][]byte, iterations int) []FuzzResult {
    var crashes []FuzzResult
    
    for i := 0; i < iterations; i++ {
        // Pick a corpus entry to mutate
        base := corpus[randInt(len(corpus))]
        input := mutate(base)
        
        var result FuzzResult
        result.Input = input
        
        start := time.Now()
        
        func() {
            defer func() {
                if r := recover(); r != nil {
                    result.Panicked = true
                    result.ErrorMsg = fmt.Sprintf("%v", r)
                }
            }()
            
            if err := parser(input); err != nil {
                // Errors are expected — panics are bugs
            }
        }()
        
        result.Duration = time.Since(start)
        
        if result.Panicked {
            crashes = append(crashes, result)
            fmt.Printf("[CRASH] Input: %x\n  Panic: %s\n", input, result.ErrorMsg)
        }
    }
    
    return crashes
}

func mutate(input []byte) []byte {
    if len(input) == 0 {
        return []byte{byte(randInt(256))}
    }
    
    result := make([]byte, len(input))
    copy(result, input)
    
    strategy := randInt(5)
    switch strategy {
    case 0: // Bit flip
        pos := randInt(len(result))
        result[pos] ^= byte(1 << randInt(8))
    case 1: // Byte change
        pos := randInt(len(result))
        result[pos] = byte(randInt(256))
    case 2: // Insert bytes
        pos := randInt(len(result))
        insert := make([]byte, randInt(16)+1)
        rand.Read(insert)
        result = append(result[:pos], append(insert, result[pos:]...)...)
    case 3: // Delete bytes
        if len(result) > 1 {
            pos := randInt(len(result))
            result = append(result[:pos], result[pos+1:]...)
        }
    case 4: // Magic values
        magics := [][]byte{
            {0xFF, 0xFF, 0xFF, 0xFF},
            {0x00, 0x00, 0x00, 0x00},
            {0x80, 0x00, 0x00, 0x00},
            {0x7F, 0xFF, 0xFF, 0xFF},
        }
        magic := magics[randInt(len(magics))]
        if len(result) >= 4 {
            pos := randInt(len(result) - 4)
            copy(result[pos:], magic)
        }
    }
    
    return result
}

func randInt(n int) int {
    if n <= 0 {
        return 0
    }
    max := big.NewInt(int64(n))
    v, _ := rand.Int(rand.Reader, max)
    return int(v.Int64())
}

// Example: fuzz a binary protocol parser
func exampleParser(data []byte) error {
    if len(data) < 4 {
        return fmt.Errorf("too short")
    }
    
    length := binary.BigEndian.Uint32(data[:4])
    
    // BUG: if length is huge, this allocates too much
    // VULNERABILITY: integer not checked against data length
    _ = make([]byte, length)  // This would panic/OOM on 0xFFFFFFFF
    
    return nil
}

func main() {
    corpus := [][]byte{
        {0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'},
        {0x00, 0x00, 0x00, 0x00},
    }
    
    fmt.Println("Fuzzing exampleParser...")
    crashes := fuzzParser(exampleParser, corpus, 1000)
    fmt.Printf("Found %d crashes in 1000 iterations\n", len(crashes))
    
    // Minimize crash case
    if len(crashes) > 0 {
        fmt.Printf("First crash input: %x\n", crashes[0].Input)
        _ = bytes.NewBuffer(crashes[0].Input)
    }
}
```

---

## Responsible Disclosure

```
90-day standard (Google Project Zero policy):
Day 0   → Report to vendor via security@vendor.com
Day 0   → Vendor acknowledges within 7 days
Day 0-90 → Vendor develops and releases patch
Day 90  → Researcher publishes full details

If no response after 7 days → escalate
If no patch after 90 days → publish (with or without patch)

Bug Bounty Programs:
- HackerOne, Bugcrowd — platform-mediated
- Direct vendor programs (Google, Microsoft, Apple, Meta)
- CVSS score determines bounty amount
- Critical RCE can pay $20,000-$1,000,000+
```

---

## Summary

| Activity | Tools | Goal |
|----------|-------|------|
| Fuzzing | AFL++, libfuzzer, Go fuzz | Find crashes/panics |
| Static analysis | Semgrep, CodeQL, Ghidra | Find unsafe code patterns |
| Dynamic analysis | GDB, valgrind, AddressSanitizer | Confirm exploitability |
| Disclosure | Email, HackerOne, CVE database | Responsibly notify |

---

## Exercises

1. Write a Go fuzz test for a JSON parser function — find an input that causes a panic
2. Run Semgrep with the security rules on a Go project you have access to
3. Read a real bug bounty write-up from HackerOne ($10k+ bounties) — understand the root cause
4. Research AddressSanitizer and how it helps find memory bugs that fuzzers trigger
