# Chapter 33: Process Injection — Hiding Code in Other Processes

*Process injection is one of the most powerful evasion techniques. By injecting malicious code into a legitimate process (explorer.exe, svchost.exe), attackers make their code look like normal system activity.*

---

## Why Process Injection?

- **Evasion:** Code runs under a trusted process name
- **Privilege:** Inject into a higher-privilege process to inherit its token
- **Persistence:** Code survives even if the original loader exits
- **Network evasion:** Outbound connections from explorer.exe look normal; from malware.exe don't

---

## Classic DLL Injection (Windows)

The most fundamental injection technique:

```
1. Open target process (OpenProcess)
2. Allocate memory in target (VirtualAllocEx)
3. Write DLL path into that memory (WriteProcessMemory)
4. Create remote thread that calls LoadLibrary (CreateRemoteThread)
5. Windows loads your DLL into the target process
```

```c
// C code (Windows API)
HANDLE hProc = OpenProcess(PROCESS_ALL_ACCESS, FALSE, targetPID);

LPVOID pMem = VirtualAllocEx(hProc, NULL, strlen(dllPath)+1, 
                              MEM_COMMIT, PAGE_READWRITE);

WriteProcessMemory(hProc, pMem, dllPath, strlen(dllPath)+1, NULL);

HANDLE hThread = CreateRemoteThread(hProc, NULL, 0,
    (LPTHREAD_START_ROUTINE)GetProcAddress(GetModuleHandle("kernel32.dll"), "LoadLibraryA"),
    pMem, 0, NULL);
```

---

## Shellcode Injection

Instead of a DLL, inject raw shellcode:

```
1. OpenProcess
2. VirtualAllocEx (with PAGE_EXECUTE_READWRITE)
3. WriteProcessMemory (write shellcode bytes)
4. CreateRemoteThread (execute at shellcode address)
```

**Detection by GoShield:**
- `VirtualAllocEx` with PAGE_EXECUTE_READWRITE = memory injection
- `WriteProcessMemory` from non-parent process = injection
- `CreateRemoteThread` across processes = injection
- Anonymous executable memory regions in `/proc/PID/maps`

---

## Process Hollowing

Replace a legitimate process's code with malicious code:

```
1. Create a suspended legitimate process (e.g., svchost.exe)
2. Unmap its memory (ZwUnmapViewOfSection)
3. Allocate new memory at same base address
4. Write malicious PE into the memory
5. Fix up import tables, relocations
6. Resume thread
```

The process still shows as "svchost.exe" in task manager but runs your code.

---

## Linux Process Injection

```c
// Linux: inject via ptrace
ptrace(PTRACE_ATTACH, pid, NULL, NULL);  // attach to target
waitpid(pid, NULL, 0);                    // wait for stop

// Read/write process memory
ptrace(PTRACE_PEEKDATA, pid, addr, NULL);  // read word
ptrace(PTRACE_POKEDATA, pid, addr, data);  // write word

// Inject shellcode by overwriting code at current instruction pointer
// Then resume execution
ptrace(PTRACE_CONT, pid, NULL, NULL);
```

---

## Detecting Injection in Go (GoShield)

```go
// In GoShield process monitor — detect injection indicators

func (pw *ProcessWatcher) checkInjection(pid int32) []string {
    var indicators []string
    
    // Check 1: anonymous executable memory regions
    mapsPath := fmt.Sprintf("/proc/%d/maps", pid)
    data, err := os.ReadFile(mapsPath)
    if err != nil {
        return nil
    }
    
    for _, line := range strings.Split(string(data), "\n") {
        // Format: addr perms offset dev ino path
        // e.g.: 7f1234560000-7f1234570000 r-xp 00000000 00:00 0
        if len(line) < 20 {
            continue
        }
        // Anonymous executable page (no file backing, executable)
        if strings.Contains(line, "r-xp") || strings.Contains(line, "rwxp") {
            parts := strings.Fields(line)
            if len(parts) < 5 || parts[4] == "0" {
                // inode is 0 = anonymous mapping
                if strings.Contains(line, "rwxp") {
                    indicators = append(indicators, 
                        fmt.Sprintf("rwx anonymous memory at %s", parts[0]))
                }
            }
        }
    }
    
    // Check 2: process opened suspicious file descriptors
    fdPath := fmt.Sprintf("/proc/%d/fd", pid)
    fds, err := os.ReadDir(fdPath)
    if err == nil {
        for _, fd := range fds {
            target, _ := os.Readlink(filepath.Join(fdPath, fd.Name()))
            if strings.Contains(target, "memfd:") {
                // memfd_create is used for fileless execution
                indicators = append(indicators, "memfd file descriptor (fileless exec)")
            }
        }
    }
    
    return indicators
}
```

---

## Summary

| Technique | How | Detection |
|-----------|-----|-----------|
| DLL injection | LoadLibrary in remote thread | CreateRemoteThread across processes |
| Shellcode injection | VirtualAllocEx + WriteProcessMemory | RWX memory, remote thread |
| Process hollowing | Unmap legit process, write malcode | Image mismatch, unbacked executable memory |
| Linux ptrace inject | ptrace POKEDATA | ptrace calls on non-child processes |

---

## Exercises

1. Study the GoShield process monitor code (Chapter 61). How does it detect parent-child process anomalies?
2. Research memfd_create — what is it used for legitimately, and how is it abused?
3. Write detection rules for GoShield to catch process injection (Chapter 65's rule format)
4. Research "Reflective DLL Injection" — how does it avoid the WriteProcessMemory pattern?
