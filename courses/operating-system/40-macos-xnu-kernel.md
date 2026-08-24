# Chapter 40: macOS and the XNU Kernel

> **"XNU is the unusual hybrid at the heart of every Mac, iPhone, iPad, and Apple Watch. It fuses the Mach microkernel (for message-passing and virtual memory) with the BSD Unix subsystem (for file systems, networking, POSIX) into a single kernel — keeping the theoretical clarity of a microkernel while running with the performance of a monolithic design."**

---

## Table of Contents

1. [macOS History — From NeXT to Apple Silicon](#1-macos-history--from-next-to-apple-silicon)
2. [XNU Kernel Architecture](#2-xnu-kernel-architecture)
3. [Mach — The Microkernel Foundation](#3-mach--the-microkernel-foundation)
4. [BSD Layer — Unix Compatibility](#4-bsd-layer--unix-compatibility)
5. [IOKit — The Driver Framework](#5-iokit--the-driver-framework)
6. [macOS File System — APFS and HFS+](#6-macos-file-system--apfs-and-hfs)
7. [Darwin — The Open Source Foundation](#7-darwin--the-open-source-foundation)
8. [iOS — The Mobile Descendant](#8-ios--the-mobile-descendant)
9. [Apple Silicon — From Intel to ARM](#9-apple-silicon--from-intel-to-arm)
10. [macOS Security Features](#10-macos-security-features)
11. [Summary](#summary)

---

## 1. macOS History — From NeXT to Apple Silicon

**1985:** Steve Jobs leaves Apple, founds NeXT.

**1988:** NeXTSTEP OS released — built on Mach microkernel + BSD Unix + Objective-C runtime. Visually stunning, technically advanced. First implementation of the Dock, real-time spell checking, display PostScript.

**1997:** Apple acquires NeXT for $429 million. Steve Jobs returns. NeXTSTEP becomes the foundation for the next Mac OS.

**2001:** Mac OS X 10.0 ships. Based on Darwin (Mach + BSD) with Aqua GUI on top. A complete break from Classic Mac OS.

**2007:** iPhone OS (later renamed iOS) ships. Uses the same XNU kernel as Mac OS X with a different userspace (no Aqua, no X11, no full POSIX).

**2020:** Apple Silicon M1. Apple transitions Macs from Intel x86-64 to its own ARM-based chips. The kernel is the same XNU — just compiled for ARM64.

```
XNU runs on:
  Intel x86-64:  Macs (2006-2020)
  Apple Silicon (ARM64): Macs (2020+), iPhones, iPads, Apple Watch, Apple TV
```

---

## 2. XNU Kernel Architecture

**XNU** stands for "X is Not Unix." It's the hybrid kernel in Darwin:

```
User Space:
┌────────────────────────────────────────────────────────────────┐
│  macOS Applications (Cocoa, UIKit, SwiftUI)                    │
│  POSIX Applications (Terminal, Homebrew)                       │
│  System Daemons (launchd, configd, powerd, etc.)               │
└────────────────────────────────────────────────────────────────┘
             │ System calls + Mach traps
Kernel Space:
┌─────────────────────────────────────────────────────────────────┐
│  BSD Layer (4.4BSD-derived):                                     │
│  Networking (TCP/IP, sockets), POSIX API, VFS, file systems     │
│  Process model, signals, Unix permissions                        │
├─────────────────────────────────────────────────────────────────┤
│  Mach Layer:                                                     │
│  Virtual Memory Manager, Mach IPC (ports, messages), scheduling │
│  Tasks (processes), threads, IPC rights                         │
├─────────────────────────────────────────────────────────────────┤
│  IOKit (C++ object system for drivers)                          │
├─────────────────────────────────────────────────────────────────┤
│  Platform Support Layer (architecture-specific)                 │
└─────────────────────────────────────────────────────────────────┘
```

**The key architectural choice:**
Mach and BSD are in the SAME address space (same kernel). This is NOT a true microkernel (where they'd be in separate processes). XNU chose this for performance — Mach was too slow when IPC crossed process boundaries.

---

## 3. Mach — The Microkernel Foundation

**Mach** was developed at Carnegie Mellon University in the 1980s.

**Mach abstractions:**

**Tasks:**
Like processes — an address space + set of resources:
```
Mach task:
  virtual address space (map of memory regions)
  set of threads
  port set (communication rights)
  exception ports (debugging)
```

**Threads:**
Execution contexts within tasks. Mach threads have a BSD peer (`struct thread` has a `uu_proc` pointer to the BSD proc).

**Ports:**
Mach's fundamental IPC mechanism. A port is a message queue protected by capability:
```
Port types:
  SEND:      Can send messages to this port
  RECEIVE:   Can receive messages from this port (only one holder!)
  SEND_ONCE: Can send ONE message (like a reply-to address)
  DEAD_NAME: Port was destroyed (notification mechanism)

Communication:
  Sender: knows a port SEND right → mach_msg(SEND, port, message)
  Receiver: holds SEND right → mach_msg(RECEIVE, port, &message)
  
  Kernel buffers up to N messages; receiver blocks if empty
```

**Mach IPC is used for:**
```
System calls → kernel replies on reply port
Task creation, thread management
XPC services (macOS app-to-app IPC built on Mach ports)
Exception handling (debugger: gdb/lldb attaches via exception ports)
Launchd: services registered as Mach ports; clients look up by name
```

**Mach VM:**
Mach's virtual memory is one of its strongest features:
```
vm_allocate(task, &address, size, TRUE)   — allocate virtual region
vm_map(task, &address, size, ...)         — map memory
vm_copy(src_task, src_addr, size, dst_task, dst_addr)  — copy between tasks
vm_remap(dst_task, dst_addr, size, src_task, src_addr) — remap (COW)
vm_protect(task, address, size, FALSE, VM_PROT_READ)   — change protections
```

The Mach VM implements **Copy-On-Write** sharing between tasks (fork uses this).

---

## 4. BSD Layer — Unix Compatibility

The BSD layer provides the Unix interface that most programs use:

**BSD system calls (POSIX):**
```
File I/O: open, read, write, close, stat, lseek, mmap
Processes: fork, execve, waitpid, exit, getpid, kill
Signals: signal, sigaction, sigprocmask
Network: socket, bind, connect, send, recv
Filesystem: mkdir, rmdir, link, unlink, rename, readdir
IPC: pipe, socketpair, shm_open, mq_open
```

**BSD proc/thread overlay:**
Every Mach task has a corresponding BSD `proc`. Every Mach thread in that task has a BSD `uthread`:
```
Mach task → struct proc (BSD process)
  Mach thread 1 → struct uthread (BSD thread state for this thread)
  Mach thread 2 → struct uthread
```

**The syscall path:**
```
System call from user space:
  POSIX call: goes to BSD layer first → may call Mach VM for memory operations
  Mach trap:  goes directly to Mach layer

Example: mmap()
  1. User calls mmap() → BSD mmap handler
  2. BSD calls vm_map() (Mach VM) to create memory mapping
  3. BSD wraps result in POSIX return values
```

---

## 5. IOKit — The Driver Framework

**IOKit** is Apple's driver framework, written in a **restricted subset of C++**:

```
C++ class hierarchy for drivers:
  OSObject (base object class)
  └── IORegistryEntry (has a place in the IO registry tree)
      └── IOService (can be matched, probed, started, stopped)
          ├── IOPlatformDevice (platform devices)
          ├── IOPCIDevice (PCIe devices)
          ├── IOUSBDevice (USB devices)
          ├── IOBlockStorageDevice (disks)
          └── IONetworkController (NICs)
```

**IO Registry:**
A tree of active objects representing all connected hardware:
```bash
# View the IO Registry from Terminal:
ioreg -l | head -100    # raw output
ioreg -p IOService      # service plane (how devices are connected)
ioreg -p IOUSBPlane     # USB tree

# Example output:
#   +-o Root  <class IORegistryEntry>
#     +-o AppleACPIPlatformExpert  <class AppleACPIPlatformExpert>
#       +-o PCI0@0  <class IOPCIBridge>
#         +-o NVME@1C  <class IOPCIDevice>
#           +-o IONVMeController <class IONVMeController>
```

**Driver lifecycle:**
```
1. New device appears (USB plugged in, PCIe enumeration)
2. IOKit creates an IOService nub (stub) for the device
3. IOKit runs matching (finds driver with matching personality dict)
4. IOKit calls driver's probe() → check if it really handles this device
5. IOKit calls start() → driver initializes hardware
6. Driver publishes itself as a provider for higher-level drivers
7. Higher-level drivers (e.g., storage driver for USB drive) attach
8. Device removed → stop() → terminate() called up the stack
```

---

## 6. macOS File System — APFS and HFS+

**APFS (Apple File System, 2017):**
(Covered in depth in Chapter 27, APFS section)
- Copy-on-write, snapshots, clones
- Space-sharing APFS containers
- Native encryption (AFPS volumes can be encrypted independently)
- Nanosecond timestamps
- Fast directory sizing

**HFS+ (legacy):**
Still used for some Time Machine backups and non-APFS volumes. Will eventually be deprecated.

**Special macOS VFS features:**

**Extended attributes (xattr):**
```bash
# macOS uses xattrs for resource forks, quarantine flags, Spotlight metadata:
xattr file.txt                     # list xattrs
xattr -p com.apple.quarantine file.txt  # read quarantine flag
xattr -d com.apple.quarantine file.txt  # remove quarantine
# Files downloaded from internet get com.apple.quarantine xattr → "open anyway?" dialog
```

**Spotlight (metadata search):**
```bash
# Spotlight indexes file content and metadata via fsevents (kernel notifications)
mdfind "kMDItemDisplayName == '*.py'"   # find Python files
mdfind -onlyin ~/Documents "budget"    # search within a directory
```

**FSEvents:**
```c
// FSEvents: kernel API to watch directory trees for changes
// Used by: Spotlight, Time Machine, file sync tools (Dropbox, iCloud)
CFStringRef path = CFSTR("/Users/alice/Documents");
FSEventStreamRef stream = FSEventStreamCreate(NULL,
    callback, NULL, 1, &path, 0, 0);
FSEventStreamScheduleWithRunLoop(stream, CFRunLoopGetCurrent(), kCFRunLoopDefaultMode);
FSEventStreamStart(stream);
```

---

## 7. Darwin — The Open Source Foundation

**Darwin** is Apple's open-source Unix foundation:
- XNU kernel (source at github.com/apple/darwin-xnu)
- BSD userspace tools (ls, cp, grep, etc.)
- launchd (PID 1 replacement for init/systemd)
- dyld (dynamic linker)
- Shell (/bin/zsh since macOS Catalina, was /bin/bash)

**Not open source:**
- Aqua GUI (Quartz Compositor, WindowServer)
- CoreGraphics, CoreAnimation, Metal
- Most Apple-specific frameworks (Cocoa, UIKit)

**launchd:**
```xml
<!-- /Library/LaunchDaemons/com.example.myservice.plist -->
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" ...>
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.example.myservice</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/myservice</string>
        <string>--config</string>
        <string>/etc/myservice.conf</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

```bash
launchctl load /Library/LaunchDaemons/com.example.myservice.plist
launchctl unload /Library/LaunchDaemons/com.example.myservice.plist
launchctl list | grep myservice
```

---

## 8. iOS — The Mobile Descendant

**iOS uses the exact same XNU kernel as macOS.** The kernel binary is nearly identical. What differs:

**Missing in iOS (security by design):**
- No user-accessible shell (no /bin/sh via terminal)
- No app-to-app file system access (each app sandboxed in ~/Library/Mobile Documents/)
- No sideloading (apps must go through App Store — except developer mode)
- No DTrace/eBPF for debugging production devices

**iOS-specific additions:**
- IOKit drivers for touch screen, GPU, secure enclave, Face ID sensor
- Springboard (home screen) instead of Finder/Dock
- Jetsam: iOS memory pressure manager (aggressively kills background apps)
- App sandbox enforced by macOS sandbox + TCC (Transparency, Consent, Control)

**Jetsam (iOS memory killer):**
Unlike Linux's OOM killer (kills most expensive process), iOS Jetsam preemptively kills background apps based on:
- Memory pressure level
- App priority (frontmost app gets most resources)
- App type (audio app gets priority to keep playing music)
- Time since last use

---

## 9. Apple Silicon — From Intel to ARM

**Apple M1 (2020)** was a landmark — Apple's first Mac chip:
- ARM64 (8 high-performance + efficiency cores)
- Unified memory architecture (CPU, GPU, Neural Engine share LPDDR5)
- 16GB memory feels like 32GB Intel (no separate GPU VRAM)

**How XNU runs on both x86-64 and ARM64:**
```
XNU source is platform-independent (Mach and BSD layers)
Architecture-specific code in arch/x86_64/ and arch/arm64/

Rosetta 2: binary translator for x86-64 apps on Apple Silicon
  - AOT (ahead-of-time) translation at first launch
  - Translated binary cached; subsequent launches are fast
  - Runtime: dynamic translation of x86 JIT code
```

**Unified Memory Architecture impact on OS:**
```
Traditional: CPU RAM + discrete GPU VRAM (separate)
Apple Silicon: one pool of memory, accessible by CPU, GPU, and Neural Engine

OS implication:
  - No copy between CPU and GPU memory (Metal API directly addresses same DRAM)
  - Memory pressure applies to GPU too (not just CPU allocations)
  - IOKit maps GPU resources differently (no PCI DMA separate from system RAM)
```

---

## 10. macOS Security Features

**System Integrity Protection (SIP, rootless):**
```bash
csrutil status
# System Integrity Protection status: enabled.

# SIP prevents modifications to:
# /System, /usr, /bin, /sbin, /private/var/db/receipts
# Even root (uid=0) cannot modify these!

# Protected by: kernel checks, not just permissions
# To disable (requires Recovery Mode): csrutil disable
```

**TCC (Transparency, Consent, Control):**
User must grant permission for apps to access:
```
Camera, Microphone, Location, Contacts, Calendar, Reminders, Photos, Files
Network access (enterprise MDM-managed)
Accessibility, Screen Recording, Full Disk Access

Stored in: ~/Library/Application Support/com.apple.TCC/TCC.db
Enforced by: TCC daemon + kernel sandbox
```

**Gatekeeper:**
```bash
# Applications must be:
# 1. Downloaded from App Store, OR
# 2. Signed with a Developer ID certificate from Apple
# 3. Notarized (scanned by Apple's servers for malware)

# To bypass for trusted software (user-approved):
# System Settings → Privacy & Security → "Open Anyway"

# Command line:
spctl --assess /Applications/MyApp.app
codesign --verify --deep --strict /Applications/MyApp.app
```

**Secure Boot:**
```
Apple Silicon:
  Boot ROM (on chip) → verifies signed macOS kernel
  Chain of trust: Apple ROM → iBoot → XNU
  
Recovery Mode: separate firmware partition, separate signed kernel
```

---

## Summary

| Concept | Description |
|---------|------------|
| XNU | Hybrid kernel: Mach + BSD in same address space |
| Mach | Microkernel providing VM, scheduling, ports (IPC) |
| Mach port | Capability-based IPC endpoint; send/receive rights |
| BSD layer | POSIX API, networking, VFS, signals on top of Mach |
| IOKit | C++ object-oriented driver framework; IO Registry |
| APFS | Apple File System: COW, snapshots, encryption, space sharing |
| launchd | PID 1 on macOS; service manager + socket activation |
| Darwin | Open-source Unix foundation of macOS/iOS |
| Rosetta 2 | x86-64 → ARM64 binary translator; enables Intel app compatibility |
| Jetsam | iOS memory killer; proactively kills background apps |
| SIP | System Integrity Protection; kernel-enforced write protection on system dirs |
| TCC | Transparency, Consent, Control; user-approved access to sensitive resources |
| Gatekeeper | Requires App Store or Developer ID + notarization |
| Secure Boot | Chain of trust from Apple ROM through signed boot stages to kernel |
