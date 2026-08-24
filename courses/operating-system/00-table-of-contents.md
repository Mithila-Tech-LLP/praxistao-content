# Operating Systems: From Zero to Building Your Own OS

## A Complete Course — Beginner to Advanced

> **Build your own OS step by step. Understand every principle behind Windows, Linux, macOS, iOS, Android, and every other OS running on every device around you.**

---

## How This Course Is Organized

This course is structured in **10 Volumes**, each building on the previous. You start with "what is an OS?" and finish by building a working operating system from scratch — bootloader, kernel, memory manager, scheduler, file system, and shell.

**Reading approach:** Every chapter explains the concept in plain English first, then goes into technical depth. Code chapters always explain what you're about to build BEFORE writing a single line. You should never be confused about why you're doing something.

---

## Volume 1: The Big Picture — What Is an Operating System?

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 01](01-what-is-an-operating-system.md) | What Is an Operating System? | OS as a manager between hardware and software |
| [Ch 02](02-history-of-operating-systems.md) | A Brief History of Operating Systems | From punch cards to Android |
| [Ch 03](03-types-of-operating-systems.md) | Types of Operating Systems | Desktop, Mobile, RTOS, Embedded, Distributed |
| [Ch 04](04-how-your-computer-really-works.md) | How Your Computer Really Works | CPU, RAM, I/O, buses — the hardware picture |
| [Ch 05](05-how-a-computer-boots.md) | How a Computer Boots | BIOS/UEFI, bootloader, kernel startup |
| [Ch 06](06-the-kernel.md) | The Kernel — Heart of the OS | Kernel space, user space, privilege rings |
| [Ch 07](07-system-calls.md) | System Calls — Talking to the Kernel | The API between programs and the OS |

---

## Volume 2: Processes and Threads — Managing Running Programs

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 08](08-what-is-a-process.md) | What Is a Process? | Programs vs. processes, PCB, address space |
| [Ch 09](09-process-states-and-lifecycle.md) | Process States and Lifecycle | New, ready, running, blocked, zombie |
| [Ch 10](10-cpu-scheduling.md) | CPU Scheduling — Who Runs Next? | FCFS, SJF, Round Robin, Priority, CFS |
| [Ch 11](11-threads.md) | Threads — Lightweight Processes | User threads, kernel threads, POSIX pthreads |
| [Ch 12](12-concurrency-problems.md) | Concurrency Problems | Race conditions, critical sections |
| [Ch 13](13-synchronization.md) | Synchronization — Locks and Semaphores | Mutex, semaphore, spinlock, monitor |
| [Ch 14](14-deadlocks.md) | Deadlocks — When Everyone Is Waiting | Conditions, prevention, detection, recovery |
| [Ch 15](15-inter-process-communication.md) | Inter-Process Communication (IPC) | Pipes, sockets, shared memory, message queues |

---

## Volume 3: Memory Management — The Art of Sharing RAM

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 16](16-physical-memory-basics.md) | Physical Memory Basics | Address space, RAM layout, byte addressing |
| [Ch 17](17-memory-allocation.md) | Memory Allocation — malloc and free | Heap, fragmentation, allocator algorithms |
| [Ch 18](18-virtual-memory-and-paging.md) | Virtual Memory and Paging | Page tables, virtual addresses, MMU |
| [Ch 19](19-how-page-tables-work.md) | How Page Tables Work | Multi-level paging, CR3, TLB |
| [Ch 20](20-page-faults-and-demand-paging.md) | Page Faults and Demand Paging | Page faults, swap space, working set |
| [Ch 21](21-page-replacement-algorithms.md) | Page Replacement Algorithms | FIFO, LRU, Clock algorithm |
| [Ch 22](22-segmentation-and-memory-protection.md) | Segmentation and Memory Protection | Segments, privilege checking, GDT |

---

## Volume 4: File Systems — Persistent Storage

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 23](23-file-system-concepts.md) | File System Concepts | Files, directories, inodes, paths |
| [Ch 24](24-fat-file-system.md) | FAT File System | FAT12/16/32 — used in USB drives everywhere |
| [Ch 25](25-ext4-linux-file-system.md) | ext4 — The Linux File System | Extents, journaling, directory trees |
| [Ch 26](26-ntfs-windows-file-system.md) | NTFS — Windows File System | MFT, permissions, journaling, compression |
| [Ch 27](27-modern-file-systems.md) | Modern File Systems (ZFS, APFS, Btrfs) | Copy-on-write, snapshots, checksums |
| [Ch 28](28-virtual-file-system.md) | Virtual File Systems (VFS) | Abstraction layer, mount points, /proc |

---

## Volume 5: I/O, Drivers, and Storage

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 29](29-interrupts.md) | Interrupts — Hardware Talking to the CPU | IRQs, ISRs, PIC, APIC |
| [Ch 30](30-device-drivers.md) | Device Drivers — How Hardware Gets Controlled | Driver model, character vs block devices |
| [Ch 31](31-disk-io-and-scheduling.md) | Disk I/O and Disk Scheduling | Cylinder scheduling, SSD differences |
| [Ch 32](32-networking-in-the-os.md) | Networking in the OS | Network stack, sockets, TCP/IP layers |

---

## Volume 6: Security and Protection

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 33](33-protection-rings.md) | Protection Rings and Privilege Levels | Ring 0–3, kernel/user mode |
| [Ch 34](34-access-control.md) | Access Control — Who Can Touch What | DAC, MAC, RBAC, ACLs |
| [Ch 35](35-os-security-attacks-and-defenses.md) | OS Security — Common Attacks and Defenses | Buffer overflow, ASLR, stack canaries |
| [Ch 36](36-virtualization-and-containers.md) | Virtualization and Containers | Hypervisors, VMs, Docker, namespaces |

---

## Volume 7: Real Operating Systems Explained

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 37](37-unix-and-unix-philosophy.md) | Unix and the Unix Philosophy | Everything is a file, pipes, POSIX |
| [Ch 38](38-linux-architecture.md) | Linux Architecture Deep Dive | Monolithic kernel, modules, Linux internals |
| [Ch 39](39-windows-nt-architecture.md) | Windows NT Architecture | HAL, kernel, subsystems, registry |
| [Ch 40](40-macos-xnu-kernel.md) | macOS and XNU Kernel | Darwin, Mach + BSD hybrid, iOS shared base |
| [Ch 41](41-android-architecture.md) | Android Architecture | Linux kernel + Android runtime + Binder IPC |
| [Ch 42](42-embedded-and-rtos.md) | Embedded and Real-Time OS | FreeRTOS, VxWorks, Zephyr, RTOS principles |
| [Ch 43](43-os-landscape-comparison.md) | The OS Landscape — Comparison | Every major OS, where it runs, why it exists |

---

## Volume 8: Setting Up to Build an OS

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 44](44-development-environment-setup.md) | Development Environment Setup | Cross-compiler, QEMU, NASM, GCC, Makefile |
| [Ch 45](45-x86-for-os-developers.md) | x86 Architecture for OS Developers | Registers, modes, instructions, calling convention |
| [Ch 46](46-understanding-the-bootloader.md) | Understanding the Bootloader | Real mode, BIOS calls, loading the kernel |
| [Ch 47](47-multiboot-and-grub.md) | Multiboot and GRUB | Multiboot spec, kernel entry point |

---

## Volume 9: Building Your OS — Step by Step

| Chapter | Title | What You Build |
|---------|-------|----------------|
| [Ch 48](48-hello-from-the-kernel.md) | Hello from the Kernel | Minimal kernel that prints to screen (VGA) |
| [Ch 49](49-gdt-setup.md) | Setting Up Protected Mode (GDT) | GDT, segment descriptors, entering 32-bit mode |
| [Ch 50](50-idt-exception-handlers.md) | Interrupt Descriptor Table (IDT) | IDT setup, exception handlers, IRQ handling |
| [Ch 51](51-pic-and-hardware-interrupts.md) | PIC and Hardware Interrupts | 8259 PIC, enabling IRQs, keyboard/timer |
| [Ch 52](52-physical-memory-manager.md) | Physical Memory Manager | Bitmap allocator, tracking free/used pages |
| [Ch 53](53-paging-virtual-memory.md) | Paging — Virtual Memory | Setting up page directory and page tables |
| [Ch 54](54-heap-allocator.md) | Heap Allocator (kmalloc/kfree) | Kernel heap, first-fit allocator |
| [Ch 55](55-processes-and-pcb.md) | Processes and the PCB | Creating processes, process control block |
| [Ch 56](56-context-switching.md) | Context Switching | Saving/restoring registers, TSS |
| [Ch 57](57-process-scheduler.md) | Process Scheduler | Round-robin scheduler, timer interrupt |
| [Ch 58](58-system-calls.md) | System Calls | int 0x80 / syscall, system call table |
| [Ch 59](59-user-mode.md) | User Mode — Running User Programs | Ring 3 transition, ELF loader basics |
| [Ch 60](60-keyboard-driver.md) | Keyboard Driver | Scan code to ASCII, input buffer |
| [Ch 61](61-vfs-and-ramdisk.md) | VFS and a RAM Disk File System | initrd, basic file operations |
| [Ch 62](62-a-basic-shell.md) | A Basic Shell | Read input, parse commands, exec |

---

## Volume 10: Advanced Topics

| Chapter | Title | Key Concept |
|---------|-------|-------------|
| [Ch 63](63-smp-multi-core-os.md) | SMP — Multi-Core Operating Systems | APIC, CPU affinity, SMP-safe locking |
| [Ch 64](64-64-bit-os.md) | 64-bit OS (x86-64 Long Mode) | Entering long mode, 4-level paging |
| [Ch 65](65-future-of-operating-systems.md) | The Future of Operating Systems | Unikernels, WASM OS, seL4, Fuchsia |

---

## What You Will Be Able to Do After This Course

1. **Explain** exactly how any operating system works — from process scheduling to virtual memory to file systems
2. **Read** Linux kernel source code and understand what it does
3. **Debug** OS-level problems (memory leaks, deadlocks, performance bottlenecks)
4. **Build** a working x86 OS kernel that boots, runs processes, and has a file system
5. **Design** an OS for a specific use case (embedded, real-time, secure, distributed)
6. **Understand** why Windows, Linux, macOS, Android, and iOS make the design choices they do

---

## Prerequisites

- Basic C programming (you know what a pointer is)
- Basic command-line usage (Linux/macOS terminal or WSL)
- No prior OS knowledge needed — we start from "what even is an OS?"

## Tools Used

- **Language:** C (kernel), x86 Assembly (boot code, low-level glue)
- **Assembler:** NASM
- **Compiler:** GCC (cross-compiler targeting i686-elf)
- **Emulator:** QEMU (runs your OS without real hardware)
- **Debugger:** GDB + QEMU's built-in debugging
- **Build system:** GNU Make
- **OS to run tools on:** Linux or macOS (WSL2 works on Windows)

---

*"An operating system is not magic — it is just software with extra responsibilities. Once you understand those responsibilities, the magic disappears and clarity takes its place."*
