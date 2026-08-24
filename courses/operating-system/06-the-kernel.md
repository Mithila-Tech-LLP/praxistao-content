# Chapter 06: The Kernel — Heart of the OS

> **"The kernel is the most privileged piece of software on your computer. It can do anything — and that's exactly why it must be written with extreme care. Every bug is a potential security hole or system crash."**

---

## Table of Contents

1. [What Is the Kernel?](#1-what-is-the-kernel)
2. [Kernel Space vs. User Space](#2-kernel-space-vs-user-space)
3. [What the Kernel Actually Contains](#3-what-the-kernel-actually-contains)
4. [Kernel Subsystems](#4-kernel-subsystems)
5. [How the Kernel Gets Loaded](#5-how-the-kernel-gets-loaded)
6. [The Kernel's Memory Layout](#6-the-kernels-memory-layout)
7. [Kernel Data Structures](#7-kernel-data-structures)
8. [Kernel Modules — Extending the Kernel](#8-kernel-modules--extending-the-kernel)
9. [Preemption and Non-Preemptive Kernels](#9-preemption-and-non-preemptive-kernels)
10. [The Linux Kernel — A Real Monolithic Kernel](#10-the-linux-kernel--a-real-monolithic-kernel)
11. [The Windows NT Kernel — A Hybrid Kernel](#11-the-windows-nt-kernel--a-hybrid-kernel)
12. [Summary](#summary)

---

## 1. What Is the Kernel?

The **kernel** is the core of the operating system. It's one program (one binary file) that:
- Loads at boot time
- Stays in memory forever (never swapped out)
- Runs with full hardware privilege (Ring 0 on x86)
- Manages ALL other software

Everything that runs on your computer — browsers, games, compilers — runs "on top of" the kernel. They can only do things the kernel allows.

**Mental model:**
Think of the kernel as the **immune system** of the computer. It's invisible when working correctly. It protects the body (other programs) from threats. When it malfunctions, everything goes wrong.

Or think of it as the **operating system's operating system** — it provides services to all software the same way an OS provides services to users.

---

## 2. Kernel Space vs. User Space

This is the most fundamental division in any modern OS:

```
┌────────────────────────────────────────────────────────┐
│                    USER SPACE                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │  Chrome  │  │  Python  │  │   vim    │ ...       │
│  └──────────┘  └──────────┘  └──────────┘           │
│                                                        │
│  Programs run here. Cannot:                           │
│  - Access arbitrary physical memory                   │
│  - Execute privileged CPU instructions                │
│  - Talk to hardware directly                          │
│  - Modify kernel memory                               │
├────────────────────────────────────────────────────────┤
│  ═══════ SYSTEM CALL BOUNDARY (the only gate) ════════ │
├────────────────────────────────────────────────────────┤
│                   KERNEL SPACE                         │
│  ┌────────────────────────────────────────────────┐   │
│  │  Scheduler │ Memory Manager │ File System │    │   │
│  │  Device Drivers │ Network Stack │ IPC │ ...   │   │
│  └────────────────────────────────────────────────┘   │
│                                                        │
│  Kernel runs here. Can:                               │
│  - Access all physical memory                         │
│  - Execute all CPU instructions (including privileged) │
│  - Configure hardware                                 │
│  - Set up/modify page tables                          │
└────────────────────────────────────────────────────────┘
         HARDWARE (CPU, RAM, disk, network, ...)
```

**Why this division is non-negotiable:**
- Without it: a buggy program could corrupt the kernel's data → immediate crash
- Without it: a malicious program could read any process's memory → steal passwords
- Without it: any program could talk to hardware → chaos and security disaster

The x86 CPU enforces this in hardware. When a user-mode program tries to execute a privileged instruction (like `cli` — clear interrupts), the CPU raises a General Protection Fault. The kernel catches it and kills the program.

---

## 3. What the Kernel Actually Contains

Depending on kernel architecture (monolithic vs microkernel), the kernel contains different things:

**Monolithic kernel (Linux, BSD):**
Everything is inside one big kernel binary:
- Process scheduler
- Virtual memory manager
- File systems (ext4, FAT32, NFS, ...)
- Device drivers (network cards, storage, USB, ...)
- Network stack (TCP/IP, sockets)
- IPC (pipes, Unix domain sockets, shared memory)
- System call interface
- Security subsystem

**Microkernel (seL4, QNX):**
Only the absolute minimum in kernel space:
- Message passing (IPC) between processes
- Memory address space management
- Basic scheduling (thread abstraction)
Everything else (file systems, drivers, network) runs as privileged user-space processes.

**For our course:** We'll build a monolithic kernel (simpler to understand and build).

---

## 4. Kernel Subsystems

A real kernel is organized into subsystems. Here's how Linux organizes its kernel:

**1. Process Management:**
- Creating and destroying processes
- Loading program binaries (ELF executables)
- CPU scheduling (who runs next)
- Context switching (save state of one process, load another)
- Signals (kill, segfault notifications)
- `fork()`, `exec()`, `wait()`

**2. Memory Management:**
- Physical memory allocator (buddy system for pages)
- Virtual memory (page tables per process)
- Demand paging (load pages from disk on demand)
- Swap management (move pages to disk when RAM is full)
- Kernel memory allocator (slab/slub for kernel objects)
- `mmap()`, `brk()`, `munmap()`

**3. File System:**
- VFS (Virtual File System) — abstraction over all file systems
- Specific file systems: ext4, FAT32, NTFS, NFS, proc, sysfs, ...
- Directory caches, inode caches
- Page cache (file contents cached in RAM)
- `open()`, `read()`, `write()`, `close()`, `mkdir()`, `stat()`

**4. I/O and Device Drivers:**
- Character devices (keyboard, serial port, /dev/null)
- Block devices (disk drives, loop devices)
- Network devices (Ethernet, WiFi)
- Driver framework (how drivers register with the kernel)
- `ioctl()`, direct I/O, DMA

**5. Network Stack:**
- Socket interface (TCP, UDP, Unix sockets)
- TCP/IP implementation
- Routing, filtering (iptables/nftables)
- Network device interface
- `socket()`, `bind()`, `connect()`, `send()`, `recv()`

**6. IPC (Inter-Process Communication):**
- Pipes and FIFOs
- Unix domain sockets
- Shared memory (`shmget`, `mmap`)
- Message queues
- Semaphores and mutexes

**7. Security:**
- DAC (Discretionary Access Control) — standard Unix permissions
- Capabilities — fine-grained privileges
- SELinux / AppArmor — Mandatory Access Control
- Seccomp — restrict system calls per process
- Namespaces — isolation (used by containers)

---

## 5. How the Kernel Gets Loaded

Recall from Chapter 5:

GRUB loads the kernel binary (e.g., `vmlinuz`) into physical memory and jumps to it. The kernel file is an ELF binary (Executable and Linkable Format) — same format as compiled C programs.

**The kernel's ELF sections:**
```
.text    — kernel code (instructions)
.data    — initialized kernel data (global variables with values)
.bss     — uninitialized global variables (kernel zeroes this at startup)
.rodata  — read-only data (string constants, lookup tables)
```

**The kernel is linked to run at a specific virtual address.** On 32-bit x86, the Linux kernel is traditionally linked to run at virtual address `0xC0000000` (3GB mark). The lower 3GB of virtual address space (0x00000000–0xBFFFFFFF) is user space. The upper 1GB (0xC0000000–0xFFFFFFFF) is kernel space.

On 64-bit x86, the kernel occupies the top half of the 48-bit virtual address space: `0xFFFF800000000000` and above.

---

## 6. The Kernel's Memory Layout

Here's how memory looks from the kernel's perspective (32-bit Linux, classic layout):

```
Virtual Address Space (32-bit Linux)
┌─────────────────────────────────────┐ 0xFFFFFFFF (4GB)
│  Kernel space (1GB)                 │
│  ─────────────────────────────────  │
│  0xFFFF0000  Exception vectors      │
│  0xF8000000  Kernel modules         │
│  0xF0000000  vmalloc area           │
│  0xC0000000  Direct-mapped RAM (1GB)│
│              Kernel code+data here  │
├─────────────────────────────────────┤ 0xC0000000 (3GB)
│  User space (3GB)                   │
│  ─────────────────────────────────  │
│  0xBFFFF000  Stack (grows down)     │
│  0x40000000  Shared libraries       │
│  0x08048000  Program text (.text)   │
│  0x00000000  (unmapped — null ptr)  │
└─────────────────────────────────────┘ 0x00000000
```

**Important:**
- Every process sees the SAME kernel virtual address range (0xC0000000+)
- The kernel is mapped into every process's page table — this is how system calls can switch from user to kernel space without changing CR3 (page table register)
- Physical RAM is "direct-mapped" into kernel space: virtual address `0xC0000000 + offset` = physical address `offset`

---

## 7. Kernel Data Structures

The kernel manages complex state about every process, file, device, and network connection. Here are the most important data structures:

**`task_struct` (Linux) — Process Descriptor:**
Every running or sleeping process is represented by a `task_struct`. It contains:
```c
struct task_struct {
    volatile long   state;      /* -1 unrunnable, 0 runnable, >0 stopped */
    void           *stack;      /* kernel stack pointer */
    pid_t           pid;        /* process ID */
    pid_t           tgid;       /* thread group ID */
    struct task_struct *parent; /* pointer to parent process */
    struct list_head children;  /* list of child processes */
    struct mm_struct *mm;       /* memory descriptor (page tables, etc.) */
    struct files_struct *files; /* open file descriptors */
    struct signal_struct *signal; /* signal handlers */
    // ... hundreds more fields
};
```
On a running Linux system, there's one `task_struct` for every process and thread.

**`mm_struct` — Memory Descriptor:**
Describes a process's virtual address space: where the stack is, where the heap is, what files are mapped in, the page table base address.

**`inode` — File System Node:**
Every file, directory, symlink, or device in a file system has an inode:
```c
struct inode {
    umode_t         i_mode;      /* permissions: rwxrwxrwx + type */
    uid_t           i_uid;       /* owner user ID */
    gid_t           i_gid;       /* owner group ID */
    loff_t          i_size;      /* file size in bytes */
    struct timespec i_atime;     /* last access time */
    struct timespec i_mtime;     /* last modification time */
    unsigned long   i_ino;       /* inode number (unique per file system) */
    // ...
};
```

**`file` structure — Open File:**
When a process opens a file, the kernel creates a `file` structure:
```c
struct file {
    struct path         f_path;  /* dentry + mount point */
    struct inode       *f_inode; /* underlying inode */
    const struct file_operations *f_op; /* function pointers: read, write, ... */
    loff_t              f_pos;   /* current read/write position */
    unsigned int        f_flags; /* O_RDONLY, O_WRONLY, O_NONBLOCK, etc. */
    // ...
};
```

**`socket` structure — Network Connection:**
Every TCP/UDP socket has a `socket` → `sock` structure with state, send/receive buffers, timers.

---

## 8. Kernel Modules — Extending the Kernel

The Linux kernel supports **loadable kernel modules (LKM)** — pieces of kernel code that can be loaded or unloaded at runtime without rebooting.

```bash
# List loaded modules
lsmod

# Load a module
sudo modprobe nvidia    # load NVIDIA GPU driver

# Unload a module
sudo rmmod nvidia

# Check if a module is for your kernel
modinfo nvidia
```

**Why modules exist:**
- Drivers for every possible hardware would make the kernel enormous
- Device drivers can be loaded only when the device is connected
- Kernel can be extended without recompiling
- Third-party drivers (NVIDIA, VirtualBox) can be distributed without being in the main kernel

**Module structure (simplified):**
```c
#include <linux/module.h>
#include <linux/init.h>

static int __init my_driver_init(void) {
    printk(KERN_INFO "My driver loaded!\n");
    // register with the kernel subsystem
    return 0;
}

static void __exit my_driver_exit(void) {
    printk(KERN_INFO "My driver unloaded!\n");
    // unregister
}

module_init(my_driver_init);
module_exit(my_driver_exit);
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Your name");
```

**Security concern:**
Kernel modules run in Ring 0 with full kernel privilege. A malicious or buggy module can crash or compromise the entire system. This is why:
- Linux Secure Boot + module signing: only signed modules can be loaded
- Some security-focused OSes (seLinux strict, grsecurity) restrict module loading

---

## 9. Preemption and Non-Preemptive Kernels

A subtle but important question: **Can the kernel be interrupted while it's running?**

**Non-preemptive kernel (old Linux, ≤ 2.4):**
- Once a thread enters kernel space (via system call), it cannot be forcibly taken off the CPU
- It runs until it voluntarily yields (e.g., sleeps waiting for disk I/O)
- Simpler to implement (no locks needed for many kernel data structures)
- Problem: a slow system call holds up all other processes

**Preemptive kernel (Linux ≥ 2.6, Windows NT always):**
- The scheduler can forcibly switch away from a process even while it's in the middle of a system call
- Better responsiveness (real-time and interactive tasks get CPU quickly)
- More complex: all kernel data structures must be protected with locks
- Required for real-time Linux (PREEMPT_RT patches)

**Fully preemptible kernel (PREEMPT_RT Linux):**
- Even interrupt handlers can be preempted
- Maximizes worst-case response time
- Used for audio production, industrial control on Linux

For our OS build, we'll start with a non-preemptible kernel (simpler), then explain how to make it preemptible.

---

## 10. The Linux Kernel — A Real Monolithic Kernel

Linux is the most widely deployed kernel in history. Let's look at its architecture:

**Key facts:**
- ~30 million lines of code (kernel 6.x)
- 2,000+ developers contribute per release
- New version every ~3 months
- Runs on: x86, ARM, RISC-V, MIPS, PowerPC, s390, and more
- Monolithic but with loadable modules

**Linux kernel source tree structure:**
```
linux/
├── arch/         Architecture-specific code (x86, arm, riscv, ...)
├── kernel/       Core: scheduler, signals, timers, IRQs
├── mm/           Memory management
├── fs/           Virtual file system + all file system drivers
├── drivers/      ALL device drivers (huge — 70% of the kernel)
├── net/          Network stack (TCP/IP, socket layer)
├── ipc/          IPC: pipes, message queues, semaphores
├── security/     SELinux, AppArmor, seccomp
├── crypto/       Cryptographic algorithms
├── lib/          Kernel utility library
├── include/      Header files
└── init/         Kernel initialization (main.c — kernel "main")
```

**The Linux boot flow (simplified):**
```
arch/x86/boot/compressed/head_64.S  → decompresses kernel
arch/x86/kernel/head_64.S           → sets up 64-bit mode, paging
init/main.c: start_kernel()         → main kernel init function
  → setup_arch()                    → CPU/memory architecture setup
  → mm_init()                       → memory management init
  → sched_init()                    → scheduler init
  → vfs_caches_init()               → file system init
  → rest_init()                     → creates init thread
    → kernel_thread(kernel_init)    → starts PID 1
      → run_init_process("/sbin/init")  → user space begins
```

---

## 11. The Windows NT Kernel — A Hybrid Kernel

The Windows NT kernel was designed from scratch in the early 1990s, separate from DOS.

**Windows NT kernel architecture:**

```
┌──────────────────────────────────────────────────────────────────┐
│                       USER SPACE                                  │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │           NTDLL.DLL (NT Native API)                         │ │
│  │  Win32 API ← user32.dll, kernel32.dll, gdi32.dll           │ │
│  └─────────────────────────────────────────────────────────────┘ │
├──────────────────────────────────────────────────────────────────┤
│                      KERNEL SPACE                                 │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │                  Win32k.sys (graphics)                    │   │
│  ├───────────────────────────────────────────────────────────┤   │
│  │  Kernel Executive:                                        │   │
│  │  Object Manager │ Process/Thread │ Memory Manager         │   │
│  │  I/O Manager    │ Security Ref   │ Cache Manager          │   │
│  ├───────────────────────────────────────────────────────────┤   │
│  │  NT Kernel: Scheduler, Interrupt handling, Synchronization│   │
│  ├───────────────────────────────────────────────────────────┤   │
│  │  HAL (Hardware Abstraction Layer)                         │   │
│  └───────────────────────────────────────────────────────────┘   │
│                        HARDWARE                                   │
└──────────────────────────────────────────────────────────────────┘
```

**Key Windows NT concepts:**

**HAL (Hardware Abstraction Layer):**
A thin layer that abstracts hardware differences. Different HAL versions for different hardware (SMP vs. single CPU, APIC vs. PIC). Higher kernel layers never deal with hardware directly — they call HAL.

**Executive:**
Upper layer of kernel — high-level policies (process management, file I/O, networking). Calls the NT Kernel below it for low-level operations.

**NT Kernel:**
Lowest software layer (above HAL). Handles interrupts, thread scheduling, low-level synchronization. Very small and focused.

**Object Manager:**
Windows manages almost everything as "objects" (files, processes, threads, events, mutexes). The Object Manager provides reference counting, security checking, and a namespace.

**Registry:**
Unlike Unix (which uses files for configuration), Windows uses the Registry — a hierarchical database of key-value pairs for system and application configuration.

---

## Summary

| Concept | Definition |
|---------|------------|
| Kernel | Core OS software; runs Ring 0; always in memory; manages everything |
| Kernel space | Memory region where kernel runs; fully privileged |
| User space | Memory where programs run; restricted; can't touch hardware |
| System call boundary | Only controlled way to go from user space to kernel space |
| Kernel subsystems | Scheduler, memory manager, VFS, drivers, network, IPC, security |
| Kernel modules | Loadable kernel extensions; run in Ring 0 |
| Preemptive kernel | Scheduler can forcibly switch even during system calls |
| Monolithic kernel | Everything in one binary (Linux, BSD) — fast but one bug can crash all |
| Microkernel | Only IPC+memory in kernel; everything else as user processes (seL4, QNX) |
| Hybrid kernel | Mix of both (Windows NT, macOS XNU) |

**The kernel is software, but it behaves like infrastructure.** Just as you trust the bridge you drive over, you trust the kernel running your programs. Building your own is the best way to truly understand this trust.
