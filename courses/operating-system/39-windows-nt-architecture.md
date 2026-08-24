# Chapter 39: Windows NT Architecture

> **"Windows NT was designed to be a professional OS from day one — portable, reliable, secure, and compatible with existing applications. While the world debated whether Microsoft could compete with Unix, NT quietly built a multi-layer architecture that now runs on hundreds of millions of devices."**

---

## Table of Contents

1. [Windows NT History](#1-windows-nt-history)
2. [NT Architecture Overview](#2-nt-architecture-overview)
3. [The NT Kernel (Ntoskrnl.exe)](#3-the-nt-kernel-ntoskrnlexe)
4. [HAL — Hardware Abstraction Layer](#4-hal--hardware-abstraction-layer)
5. [Executive Services](#5-executive-services)
6. [Win32 Subsystem](#6-win32-subsystem)
7. [The Windows Registry](#7-the-windows-registry)
8. [Windows Memory Management](#8-windows-memory-management)
9. [Windows Security Model](#9-windows-security-model)
10. [Windows vs Linux Architecture Comparison](#10-windows-vs-linux-architecture-comparison)
11. [Summary](#summary)

---

## 1. Windows NT History

**Windows NT (New Technology)** was developed by David Cutler, who previously designed the VAX/VMS operating system at Digital Equipment Corporation (DEC). Key design goals:

- **Portability:** Write in C; abstraction layer (HAL) hides hardware differences
- **Security:** C2 (later Common Criteria) security certification
- **Compatibility:** Support for POSIX, OS/2, and Win32 applications
- **Reliability:** Structured exception handling, protected kernel/user separation
- **Scalability:** SMP support from day one

**Windows NT family timeline:**
```
Windows NT 3.1 (1993) — first release, desktop + server
Windows NT 4.0 (1996) — moved GDI/graphics into kernel (performance)
Windows 2000 (NT 5.0) — Active Directory, NTFS improvements
Windows XP (NT 5.1) — merged consumer Windows with NT
Windows Vista (NT 6.0) — UAC, new security features
Windows 7 (NT 6.1) — performance improvements
Windows 8 (NT 6.2) — tablet-first design
Windows 10 (NT 10.0, 2015) — current, Windows-as-a-service
Windows 11 (NT 10.0.22000+, 2021) — new UI, security requirements
```

All are the SAME NT kernel, incrementally improved. Windows 10/11 share a kernel with Windows Server 2019/2022.

---

## 2. NT Architecture Overview

Windows NT uses a **hybrid kernel** architecture:

```
User Mode:
┌──────────────────────────────────────────────────────────────────────┐
│  Win32 Apps  │  Win32 Service  │  NT Native Apps  │  POSIX Apps     │
│  (Notepad,   │  (svchost.exe)  │  (smss.exe,      │  (Cygwin,       │
│   Chrome)    │                 │   csrss.exe)     │   WSL2)         │
│              │                 │                  │                 │
│         Subsystem DLLs (kernel32.dll, ntdll.dll)                    │
└──────────────────────────────────────────────────────────────────────┘
              │ Win32 API calls → NT native API calls
┌──────────────────────────────────────────────────────────────────────┐
│  Executive Services (kernel32.dll calls → ntdll.dll → syscall)       │
├──────────────┬────────────────────────────────────────────────────────┤
│  NT Kernel   │  HAL (Hardware Abstraction Layer)                     │
│  (ntos.exe)  │                                                       │
└──────────────┴────────────────────────────────────────────────────────┘
                    Physical Hardware
```

**Key separation:**
- **User mode:** Applications, services, subsystem DLLs
- **Kernel mode:** Executive, NT Kernel, HAL, drivers
- **ntdll.dll:** The boundary — translates Win32 calls to NT native syscalls

---

## 3. The NT Kernel (Ntoskrnl.exe)

The NT kernel has two components:

**Microkernel portion (low-level):**
- Thread scheduling and dispatching
- Interrupt and exception handling
- Multiprocessor synchronization (spinlocks)
- Kernel objects (events, mutexes, semaphores)
- Trap handling

**Executive portion (higher-level):**
Everything else in kernel mode (see section 5).

**Object Manager:**
Windows uses an **object system** to represent kernel resources:
```
Object types (like file types in Unix, but for kernel resources):
  Process      → task_struct equivalent
  Thread       → thread control block
  File         → open file instance
  Section      → memory-mapped file
  Event        → signaling mechanism
  Mutex        → mutual exclusion
  Semaphore    → counting semaphore
  Timer        → kernel timer
  Directory    → object namespace directory (not file system!)
  SymbolicLink → alias to another object
  Key          → registry key
  Token        → security access token
```

**Object namespace (separate from file system!):**
```
\                 ← root of object namespace
├── Device\       ← device objects
│   ├── HarddiskVolume1  (C: drive)
│   ├── HarddiskVolume2  (D: drive)
│   └── Null
├── GLOBAL??      ← global namespace
├── Sessions\     ← per-session namespaces
│   └── 1\       ← session 1
│       └── DosDevices\  ← C:, D:, etc. (drive letter aliases)
├── KnownDlls\
└── Windows\      ← named pipes, semaphores, shared sections
```

**Handles:**
Like file descriptors but for any object type:
```c
// Every NT resource is accessed via a handle (similar to Unix fds)
HANDLE hProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, pid);
HANDLE hFile = CreateFile("C:\\file.txt", GENERIC_READ, 0, NULL, OPEN_EXISTING, 0, NULL);
HANDLE hEvent = CreateEvent(NULL, FALSE, FALSE, NULL);

// Close handle when done (decrements reference count):
CloseHandle(hProcess);
```

---

## 4. HAL — Hardware Abstraction Layer

**HAL (hal.dll)** abstracts hardware-specific details:
```
HAL provides:
  HalGetInterruptVector()   — map hardware IRQ to software vector
  HalEnableSystemInterrupt() — enable hardware interrupt
  HalTranslateBusAddress()  — translate bus addresses
  KeStallExecutionProcessor() — precise delays
  HalMakeBeep()              — PC speaker (legacy)
  
HAL hides:
  PIC vs APIC interrupt controller differences
  UniProcessor vs SMP initialization
  Different chipset timer configurations
```

**Why HAL?**
The same NT kernel binary runs on different hardware by swapping the HAL. Originally allowed NT to run on x86, Alpha, MIPS, PowerPC, Itanium, ARM. Today, the HAL layer is thinner (ARM and x86-64 HALs, mainly), but the abstraction remains.

---

## 5. Executive Services

The **Windows Executive** is a set of kernel-mode services:

```
Executive Components:
  Object Manager (ObM):     Creates/manages kernel objects; reference counting; handles
  Process Manager:          CreateProcess, CreateThread, process/thread lifecycle
  Virtual Memory Manager (VMM): Page tables, page faults, working sets, pagefile
  I/O Manager:              IRP (I/O Request Packets) subsystem, driver framework
  Cache Manager (CcM):      File cache (similar to Linux page cache)
  Security Reference Monitor (SRM): Access checks, auditing
  Plug and Play Manager:    Device detection, driver loading, power management
  Power Manager:            Sleep/hibernate, device power states
  Win32k.sys:               GDI, USER (windows, messages) — in kernel since NT 4.0
  Configuration Manager:    Windows Registry implementation
  Local Procedure Call (LPC/ALPC): Fast IPC between local processes
```

**I/O Request Packets (IRPs):**
Windows I/O is based on IRPs — data structures passed through a driver stack:
```
Application: ReadFile(hFile, buf, size, &bytesRead, NULL)
  → I/O Manager creates an IRP_MJ_READ
  → IRP travels through driver stack:
    File system driver (NTFS) → volume manager → disk driver → controller driver
  → Each driver processes IRP and passes it to next (or completes it)
  → Completion routine called back up the stack
  → ReadFile returns
```

---

## 6. Win32 Subsystem

**Win32** is the main application environment on Windows. Most applications use it.

**Key DLLs:**
```
ntdll.dll:    NT native API; user/kernel boundary (like glibc in Linux)
kernel32.dll: Win32 kernel services (files, processes, memory, threads)
user32.dll:   Win32 user interface (windows, messages, input)
gdi32.dll:    Win32 graphics (device contexts, drawing)
advapi32.dll: Advanced Win32 (registry, security, services, events)
shell32.dll:  Windows Shell (Explorer integration, file associations)
ws2_32.dll:   Winsock (Windows networking API)
```

**Calling chain:**
```
Win32 application calls CreateFile():
  → kernel32.dll!CreateFileW()
     → ntdll.dll!NtCreateFile()
        → syscall instruction  (KiSystemCall64 in NT kernel)
           → kernel32 → ntdll are user-mode, everything below is kernel
           → Io Manager creates IRP
           → Driver stack handles the request
```

**CSRSS (Client/Server Runtime Subsystem):**
The Win32 subsystem process (csrss.exe) handles:
- Console window management
- Critical Windows bookkeeping (process tracking)
- Win32 server functions that don't need kernel mode

**Windows Subsystem for Linux (WSL2):**
Microsoft's Linux-on-Windows solution. WSL2 runs an actual Linux kernel in a lightweight Hyper-V VM:
```
WSL2 architecture:
  Windows user: runs wsl.exe
  Hyper-V lightweight VM: Linux kernel (real, unmodified Linux)
  Inside VM: Ubuntu/Debian/etc. userspace
  
  File access between Windows and Linux via 9P filesystem protocol
  Network access: shared virtual NIC
  Display: WSLg → Wayland + X11 forwarding via RDP
```

---

## 7. The Windows Registry

**The Registry** is Windows' centralized configuration database — a hierarchical key-value store:

```
Registry structure (hives):
HKEY_LOCAL_MACHINE (HKLM):  Machine-wide settings
  \SYSTEM\CurrentControlSet\Services\  → driver and service configuration
  \SOFTWARE\Microsoft\Windows NT\CurrentVersion  → Windows version info
  \SECURITY                            → security policy
  \SAM                                 → local user accounts

HKEY_USERS (HKU):           Per-user settings (loaded when user logs in)
  \S-1-5-21-...-1000        → user's settings (hive from NTUSER.DAT)

HKEY_CURRENT_USER (HKCU):  Current user's settings (alias into HKU)
  \Software\                → application settings
  \Control Panel\           → user preferences

HKEY_CLASSES_ROOT (HKCR):  File associations, COM classes
HKEY_CURRENT_CONFIG:        Current hardware profile
```

**Registry vs Linux:**
Linux stores config in text files (`/etc/`, `~/.config/`). Windows centralizes it:
```
Linux: /etc/nginx/nginx.conf     → Windows: HKLM\SOFTWARE\nginx\config
Linux: ~/.bashrc                 → Windows: HKCU\Software\profile
Linux: /etc/hosts                → Windows: HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Hosts
```

**Transactional registry:**
Windows Vista+ supports transactional registry operations:
```c
HANDLE txn = CreateTransaction(NULL, NULL, 0, 0, 0, INFINITE, NULL);
RegOpenKeyTransacted(HKEY_LOCAL_MACHINE, "SOFTWARE\\MyApp", 0, KEY_WRITE, &hKey, txn, NULL);
RegSetValueEx(hKey, "Setting", 0, REG_SZ, "value", strlen("value")+1);
CommitTransaction(txn);  // atomically applies all changes
// or RollbackTransaction(txn) to undo
```

---

## 8. Windows Memory Management

**Windows virtual address space (64-bit):**
```
0x0000000000001000  - 0x00007FFFFFFFFFFF: User space (128TB)
  0x00007FF000000000: highest user-mode address
  Stack: grows down from near top
  Heap: grows up from low addresses
  Loaded modules: various locations (ASLR randomized)
  
0xFFFF800000000000  - 0xFFFFFFFFFFFFFFFF: Kernel space (128TB)
  Kernel code, Executive data, driver space
```

**Working set:**
Each process has a **working set** — the set of pages currently in physical RAM. The **working set manager** trims process working sets when the system is under memory pressure.

**Pagefile (Windows swap):**
```
pagefile.sys: Located on the system drive (usually C:\)
Virtual memory = physical RAM + pagefile
Windows manages pagefile size automatically, or you can set manually

System-managed pagefile:
  Minimum: 1× RAM
  Maximum: 3× RAM
  
On crash: Windows writes memory dump to pagefile for analysis
```

**Section objects (memory-mapped files):**
Windows uses Section objects instead of mmap:
```c
// Create a section from a file:
HANDLE hFile = CreateFile("data.bin", GENERIC_READ|GENERIC_WRITE, ...);
HANDLE hSection = CreateFileMapping(hFile, NULL, PAGE_READWRITE, 0, 0, NULL);

// Map into process address space:
LPVOID pView = MapViewOfFile(hSection, FILE_MAP_ALL_ACCESS, 0, 0, 0);

// Now pView points to the file's content, like mmap() on Linux
// Two processes mapping the same section → shared memory!

UnmapViewOfFile(pView);
CloseHandle(hSection);
```

---

## 9. Windows Security Model

**Access tokens:**
Every process and thread has an **access token** containing:
```
Token contents:
  User SID:          S-1-5-21-domain-1000 (alice)
  Group SIDs:        S-1-5-32-544 (Administrators), S-1-5-32-545 (Users), ...
  Privilege set:     SeShutdownPrivilege, SeCreateSymbolicLinkPrivilege, ...
  Token type:        Primary (process) or Impersonation (thread)
  Integrity level:   Low, Medium, High, System (Windows Vista+)
```

**Integrity levels (UAC mechanism):**
```
System:  OS processes (Services, System process)
High:    Elevated administrator processes (after UAC prompt)
Medium:  Normal user processes (default)
Low:     Protected processes (IE sandbox, sandboxed browsers)
Untrusted: Downloaded content, sandboxed apps

Rule: Process can only write to objects at SAME OR LOWER integrity level
  Medium process cannot write to High-integrity registry keys
  Low-integrity browser cannot write to %APPDATA% (Medium) without UAC
```

**UAC (User Account Control):**
```
Application requires admin → UAC prompt
User clicks Yes → process gets a High integrity token
User clicks No → process gets Medium integrity token and might fail

Administrator account in Windows Vista+:
  Two tokens: filtered (Medium) and full (High)
  Normal operations: use filtered token
  Elevated operations: request full token via UAC
```

**Access Check:**
```
Process opens C:\Windows\System32\config\SAM (sensitive registry backup)

Security Reference Monitor checks:
  1. Get caller's token: groups, privileges, integrity level
  2. Get SAM file's DACL (Access Control List)
  3. Check ACEs:
     SYSTEM: Full Control ← but we're not SYSTEM
     Administrators: Full Control ← we're in Administrators
  4. Check integrity: SAM is System integrity, caller is High → write denied
  5. Check UAC: deny write-down across integrity levels
  
Result: Access denied!
```

---

## 10. Windows vs Linux Architecture Comparison

| Aspect | Windows NT | Linux |
|--------|-----------|-------|
| Kernel type | Hybrid (monolithic + microkernel ideas) | Monolithic |
| Kernel entry point | ntoskrnl.exe | vmlinuz |
| Hardware abstraction | HAL (hal.dll) | architecture-specific code in arch/ |
| Object model | NT Object Manager (typed objects) | file descriptors (everything is a file) |
| IPC | LPC/ALPC, COM, WCF | pipes, Unix sockets, shared memory |
| Configuration | Registry (hierarchical database) | text files in /etc, ~/.config |
| Driver framework | WDM/WDF (IRP-based) | Driver model (probe/remove, file_ops) |
| Service management | Service Control Manager | systemd |
| GUI | Integrated (Win32k.sys in kernel) | X11/Wayland (user space) |
| Security | ACLs + integrity levels + UAC | DAC + capabilities + SELinux/AppArmor |
| Source availability | Closed source | Open source |
| Modularity | DLLs (delay-load), COM | shared libraries (.so), kernel modules |

---

## Summary

| Concept | Description |
|---------|------------|
| NT kernel | Core Windows kernel; handles scheduling, interrupts, synchronization |
| Executive | Higher-level kernel services: process, I/O, memory, security managers |
| HAL | Hardware Abstraction Layer; hides chipset differences from kernel |
| Object Manager | Manages typed kernel objects; reference counting; handle table |
| IRP | I/O Request Packet; data structure passed through driver stack for I/O |
| ntdll.dll | User/kernel boundary DLL; translates Win32 API to NT native syscalls |
| Win32k.sys | Kernel-mode graphics and windowing system (GDI, USER) |
| Registry | Centralized configuration database; hierarchical key-value store |
| Pagefile | Windows swap file (pagefile.sys); extends virtual memory |
| Access token | Per-process security credential: SID, groups, privileges, integrity |
| Integrity level | Low/Medium/High/System; enforces write-down restriction; UAC mechanism |
| ALPC | Advanced Local Procedure Call; fast local IPC in Windows |
| WSL2 | Windows Subsystem for Linux 2; real Linux kernel in Hyper-V micro-VM |
