# Chapter 38: Linux Architecture Deep Dive

> **"Linux is not just an operating system kernel. It is an engineering marvel — a 30+ million line codebase maintained by thousands of contributors, running on everything from smart watches to supercomputers, updated continuously while maintaining backward compatibility for decades. Understanding Linux's architecture explains why it has won the server, the cloud, and the phone."**

---

## Table of Contents

1. [Linux Overview](#1-linux-overview)
2. [Kernel Subsystems](#2-kernel-subsystems)
3. [The Process Scheduler — CFS](#3-the-process-scheduler--cfs)
4. [Memory Management in Linux](#4-memory-management-in-linux)
5. [Linux File System — VFS in Practice](#5-linux-file-system--vfs-in-practice)
6. [The Networking Stack](#6-the-networking-stack)
7. [Linux Device Model](#7-linux-device-model)
8. [Kernel Module System](#8-kernel-module-system)
9. [eBPF — Programmable Kernel](#9-ebpf--programmable-kernel)
10. [Linux Boot Process Revisited](#10-linux-boot-process-revisited)
11. [Summary](#summary)

---

## 1. Linux Overview

**Linux statistics (as of 2026):**
- ~30 million lines of code
- ~20,000 contributors historically
- ~4,000 patches per release cycle (~2 months)
- Runs on: phones (Android), servers, cloud (98%+ of top 1M websites), supercomputers (100% of top 500), embedded systems, TVs, cars

**Linux is a monolithic kernel:**
All kernel code runs in Ring 0 (kernel space). Subsystems (scheduler, networking, filesystems) are compiled into one big binary. However, **loadable kernel modules** (LKMs) can be loaded/unloaded at runtime.

**Kernel source tree organization:**
```
linux/
├── arch/          Architecture-specific code (x86, arm64, riscv, ...)
│   └── x86/
│       ├── kernel/ CPU init, syscall entry, interrupts, SMP
│       ├── mm/     x86-specific memory management
│       └── boot/   Boot code (compressed/decompressed kernel)
├── drivers/       ~60% of kernel size! All device drivers
│   ├── char/
│   ├── block/
│   ├── net/
│   ├── gpu/       (drm, i915 for Intel, amdgpu, etc.)
│   └── usb/
├── fs/            File system implementations
│   ├── ext4/
│   ├── btrfs/
│   ├── proc/
│   └── ...
├── kernel/        Core kernel: scheduler, timer, kmod, ...
├── mm/            Memory management
├── net/           Networking stack
├── ipc/           Inter-process communication
├── security/      LSM (Linux Security Module) framework, SELinux, AppArmor
├── include/       Header files
└── init/          Kernel startup code (init/main.c)
```

---

## 2. Kernel Subsystems

```
Linux Kernel Subsystems:
┌─────────────────────────────────────────────────────────┐
│                System Call Interface                     │
├────────────┬────────────┬─────────────┬─────────────────┤
│  Process   │  Memory    │    File     │    Network      │
│ Management │ Management │   System    │    Stack        │
│            │            │   (VFS)     │                 │
├────────────┴────────────┴─────────────┴─────────────────┤
│              Device Drivers / Hardware Interface         │
│  Block Layer │ Network Drivers │ Input/USB │ GPU/Display │
└─────────────────────────────────────────────────────────┘
     │                │               │
  Hardware         Hardware       Hardware
  (Disk, SSD)    (NIC, WiFi)   (Keyboard, GPU)
```

**LSM (Linux Security Module) framework:**
A hook-based framework allowing multiple security modules to be loaded:
```
Every security-sensitive operation:
  → VFS calls security_file_open() (LSM hook)
  → LSM framework calls ALL registered modules (SELinux, AppArmor, etc.)
  → If any module denies: operation is denied
  → If all allow: operation proceeds
```

---

## 3. The Process Scheduler — CFS

The **CFS (Completely Fair Scheduler)** is Linux's default process scheduler.

**Core idea: vruntime (virtual runtime)**
Every runnable process has a `vruntime` — how much CPU time it has consumed (weighted by priority):

```c
// Updating vruntime when process uses CPU:
vruntime += actual_runtime_ns × (NICE_0_WEIGHT / process_weight)
// Lower-priority (nice 19) process has lower weight → faster vruntime growth
// Higher-priority (nice -20) process has higher weight → slower vruntime growth

// Result: lower-priority processes catch up to higher-priority ones' vruntime faster
// → They appear to need more CPU → CFS gives them more
```

**Scheduling decision:**
Always pick the process with the lowest vruntime next.

**Data structure: red-black tree (rbtree)**
All runnable processes stored in an rbtree, ordered by vruntime:
```
Left-most node = lowest vruntime = next process to run
Adding a process: insert into rbtree in O(log n)
Selecting next: pop leftmost node in O(1)
```

**Preemption:**
When the running process's vruntime catches up to the next process's vruntime, it's preempted. This ensures proportional CPU sharing.

**Scheduling classes (priority order):**
```
1. SCHED_DEADLINE: EDF (earliest deadline first) for hard real-time
2. SCHED_FIFO / SCHED_RR: real-time processes (priority 1-99)
3. SCHED_NORMAL: regular processes (the CFS domain)
4. SCHED_BATCH: background tasks (like SCHED_NORMAL but with longer timeslices)
5. SCHED_IDLE: run only when nothing else wants to run
```

**SMP load balancing:**
```
Per-CPU run queues:
  CPU 0: [nginx, python, bash]
  CPU 1: [empty]
  CPU 2: [mysqld, mysqld, mysqld]
  CPU 3: [apache, apache]

Load balancer fires periodically:
  Detects imbalance (CPU 1 idle, CPU 2 overloaded)
  Migrates mysqld from CPU 2's queue to CPU 1's queue
  Balance: each CPU gets roughly equal vruntime load
```

---

## 4. Memory Management in Linux

**Key Linux memory management features:**

**Buddy allocator:**
Manages physical page frames in power-of-2 sizes:
```
Free lists per order (0=4KB, 1=8KB, 2=16KB, ..., 10=4MB):
  Order 0: [frame 0, frame 3, frame 7, ...]
  Order 1: [frames 4-5, frames 10-11, ...]
  ...
  Order 10: [frames 0-1023, ...]

Allocate 16KB (order 2):
  If order-2 free: return it
  Else: split an order-3 (32KB) block into two order-2 → return one, keep other in order-2 list

Free:
  Add to order-2 free list
  Check if "buddy" (the other half of the split) is also free → merge to order-3
  Recursively merge upward → reduces fragmentation
```

**SLAB/SLUB allocator:**
For kernel objects smaller than a full page:
```
Caches:
  task_struct cache: maintains ~100 ready-to-use task_structs (avoids alloc+init cost)
  inode cache: maintains ready-to-use inodes
  kmalloc-64 cache: 64-byte generic allocations
  kmalloc-128 cache: 128-byte generic allocations
  ...

SLUB (the default): minimalist, good for NUMA, less memory overhead than SLAB
```

**Memory zones:**
```
ZONE_DMA:      0-16MB   (ISA DMA-capable range — legacy)
ZONE_DMA32:    0-4GB    (32-bit DMA-capable range for some devices)
ZONE_NORMAL:   rest of RAM (normal allocations)
ZONE_HIGHMEM:  (32-bit only) above 896MB — needs kmap to access
ZONE_MOVABLE:  Pages that can be migrated (for hot-plug memory removal)
```

**Transparent Huge Pages (THP):**
```bash
# Linux automatically promotes 512 consecutive 4KB pages to one 2MB huge page
# Benefits: fewer TLB misses, less page table overhead

# Control:
cat /sys/kernel/mm/transparent_hugepage/enabled
# [always] madvise never

# Check THP stats:
cat /proc/vmstat | grep thp
# thp_fault_alloc 12345
# thp_collapse_alloc 678
```

---

## 5. Linux File System — VFS in Practice

**ext4 default + many others:**
```bash
# Linux can mount many FS types simultaneously:
findmnt
# /       ext4    /dev/sda1
# /home   btrfs   /dev/sda2
# /boot   ext4    /dev/sda3
# /tmp    tmpfs   tmpfs
# /proc   proc    proc
# /sys    sysfs   sysfs
# /run    tmpfs   tmpfs
# /dev    devtmpfs devtmpfs
# /media/usb0 vfat /dev/sdb1
```

**Page cache as the center:**
```
Read a file:
  1. Check page cache: is page in memory?
  2. If yes: return from cache (no disk I/O!)
  3. If no: disk read → load into page cache → return

Write a file:
  1. Write to page cache (dirty page)
  2. Return to application (very fast)
  3. Background: writeback thread flushes dirty pages to disk

The page cache is shared between processes:
  Two processes reading the same file → same physical pages
  mmap() of the same file → same physical pages
  Saves RAM, improves performance
```

---

## 6. The Networking Stack

**Linux is a reference implementation for TCP/IP:**

```
Key networking subsystems:
  Socket layer: sys_socket, sys_bind, sys_connect, sys_send, sys_recv
  Protocol layer: inet_family_ops → tcp_prot, udp_prot
  IP layer: ip_rcv, ip_route_input, ip_forward
  Netfilter: iptables/nftables hooks (packet filtering, NAT)
  Traffic control (tc): packet scheduling, shaping, policing
  NIC driver: NAPI polling, DMA ring buffers
```

**Netfilter and iptables/nftables:**
```bash
# iptables hook points:
PREROUTING:  before routing decision (NAT, DNAT)
INPUT:       packets destined for local processes
FORWARD:     packets being routed through (router)
OUTPUT:      packets generated by local processes
POSTROUTING: after routing (NAT, SNAT/MASQUERADE)

# Example: block all SSH from the internet except one IP:
iptables -A INPUT -p tcp --dport 22 -s 192.168.1.100 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -j DROP

# NAT for containers:
iptables -t nat -A POSTROUTING -s 172.17.0.0/16 -j MASQUERADE
```

**eBPF in networking (XDP):**
```c
// XDP (eXpress Data Path): eBPF program runs at NIC driver level
// Can drop, modify, or redirect packets with near-line-rate performance

SEC("xdp")
int xdp_drop_icmp(struct xdp_md *ctx) {
    struct ethhdr *eth = data;
    struct iphdr  *ip  = data + sizeof(*eth);
    
    if (ip->protocol == IPPROTO_ICMP)
        return XDP_DROP;    // drop ping packets
    return XDP_PASS;        // allow everything else
}
// Used by: Cloudflare (DDoS mitigation), Facebook (load balancing), Cilium (k8s networking)
```

---

## 7. Linux Device Model

**sysfs reflects the device hierarchy:**
```
/sys/bus/pci/devices/:
  0000:00:1f.2  ← PCIe device (domain:bus:device.function notation)
    driver →  ahci  (symlink to driver managing this device)
    vendor:   0x8086 (Intel)
    device:   0x8c02 (8 Series SATA Controller)

/sys/class/block/:
  sda → ../../devices/pci0000:00/0000:00:1f.2/ata1/host0/target0:0:0/0:0:0:0/block/sda

/sys/class/net/:
  eth0 → ../../devices/pci.../net/eth0
  lo   → ../../devices/virtual/net/lo
```

**udev:**
The userspace device manager. When a new device is plugged in:
1. Kernel creates a kobject and signals udev via netlink socket
2. udev reads device attributes from sysfs
3. udev runs rules from `/etc/udev/rules.d/`:
   ```
   # Create /dev/disk/by-id symlinks for USB drives:
   SUBSYSTEM=="block", KERNEL=="sd*", ATTR{removable}=="1",
     SYMLINK+="disk/by-id/$env{ID_SERIAL}"
   
   # Assign specific device name to a USB NIC by MAC address:
   SUBSYSTEM=="net", ACTION=="add", ATTR{address}=="aa:bb:cc:dd:ee:ff",
     NAME="eth_office"
   ```
4. udev creates device nodes in /dev, applies permissions, runs scripts

---

## 8. Kernel Module System

**How modules are loaded:**

```
modprobe ext4
  → reads /lib/modules/$(uname -r)/modules.dep to find dependencies
  → loads dependencies first (e.g., mbcache.ko, jbd2.ko)
  → calls init_module() syscall with .ko file content
    → kernel maps the .ko into kernel memory
    → resolves symbol references (ext4 calls jbd2 functions)
    → calls the module's init function (ext4_init_fs)
    → module is now part of the kernel
```

**Module autoloading:**
The kernel automatically loads modules when needed:
```bash
# User types: mount -t ext4 /dev/sdb1 /mnt/data
# Kernel: "I don't have ext4 handler"
# Kernel calls: request_module("fs-ext4")
# Which calls: modprobe ext4
# ext4 module loads, mount proceeds

# Module aliases enable this:
cat /lib/modules/$(uname -r)/modules.alias | grep ext4
# alias fs-ext4 ext4

# USB device plugged in:
# Kernel reads VID:PID → request_module("usb:v046DpC05Ad4101...")
# udev matches and calls modprobe with the alias
```

---

## 9. eBPF — Programmable Kernel

**eBPF (extended Berkeley Packet Filter)** is one of the most transformative additions to Linux in the last decade. It allows running sandboxed programs inside the kernel without modifying kernel source.

**How eBPF works:**
```
1. User writes eBPF program in C (restricted subset)
2. Compiler (clang) compiles to eBPF bytecode
3. User loads bytecode into kernel via bpf() syscall
4. Kernel's eBPF verifier checks the program:
   - No infinite loops (bounded loops only)
   - No invalid memory accesses
   - Terminates in bounded time
5. JIT compiler converts eBPF bytecode → native machine code
6. eBPF program hooks into kernel events:
   - kprobe: any kernel function entry/exit
   - tracepoint: static kernel trace points
   - XDP: network packet receive (before network stack)
   - socket: per-socket filtering
   - perf event: performance counters
   - cgroup: per-cgroup hooks
```

**Use cases:**
```bash
# Performance profiling (bcc tools):
bpftrace -e 'kprobe:do_sys_open { printf("Open: %s\n", str(arg1)); }'
# Trace every file open in the kernel — live output, no restart needed

# Network packet filter (XDP for DDoS mitigation):
# Write 100G-rate packet filter that runs before the network stack

# System call auditing:
# Trace exact syscalls with arguments for any process

# Latency measurement:
bpftrace -e 'kprobe:ext4_file_read_iter { @start[tid] = nsecs; }
             kretprobe:ext4_file_read_iter { @ms = hist((nsecs - @start[tid])/1000); }'
# Histogram of ext4 read latencies in microseconds
```

**eBPF-based tools:**
```
bcc:      BPF Compiler Collection — Python/Lua wrapping eBPF programs
bpftrace: High-level eBPF tracing language (like dtrace for Linux)
Cilium:   Kubernetes networking + security using eBPF
Falco:    Security monitoring using eBPF
XDP:      High-performance packet processing
```

---

## 10. Linux Boot Process Revisited

**Detailed boot trace:**
```
1. BIOS/UEFI POST → finds bootable device
2. Bootloader (GRUB) loads:
   - Linux kernel image (vmlinuz or bzImage)
   - initramfs (compressed cpio archive containing early rootfs)
3. Kernel decompresses itself
4. arch/x86/boot/compressed/head_64.S: early setup (enable long mode, etc.)
5. arch/x86/kernel/head_64.S: initialize BSS, call start_kernel()
6. init/main.c: start_kernel()
   a. setup_arch() — CPU detection, NUMA topology, parse kernel cmdline
   b. mm_init() — memory allocator initialization (buddy allocator ready)
   c. sched_init() — scheduler initialization (run queues, CFS ready)
   d. rcu_init() — Read-Copy-Update (RCU) synchronization init
   e. softirq_init() — softirq vectors
   f. timekeeping_init() — timers
   g. init_IRQ() — interrupt controllers (PIC/APIC)
   h. tick_init() — periodic timer tick
   i. radix_tree_init(), idr_init() — data structure init
   j. early_irq_init(), init_IRQ() — set up IRQ subsystem
   k. trap_init() — set up IDT (exception handlers)
   l. console_init() — early console (tty0, serial port)
   m. mem_init() — finalize memory map
   n. kmem_cache_init() — SLAB/SLUB allocator
   o. vfs_caches_init() — VFS caches (dcache, icache)
   p. signals_init() — signal queues
   q. proc_root_init() — mount /proc
   r. rest_init()
     → kernel_thread(kernel_init) → spawns PID 1 (init/systemd)
     → kernel_thread(kthreadd)    → spawns PID 2 (kernel threads)
     → schedule() → idle loop
7. kernel_init (PID 1):
   a. do_initcalls() — init all subsystems (drivers, filesystems) via __initcall macros
   b. wait for initramfs setup
   c. prepare_namespace() — mount real root FS
   d. execve("/sbin/init") or "/lib/systemd/systemd"
8. systemd (or SysV init):
   a. Mounts filesystems from /etc/fstab
   b. Starts services (network, logging, ssh, etc.)
   c. Reaches login prompt or display manager
```

---

## Summary

| Component | Description |
|-----------|------------|
| Monolithic kernel | All code in kernel space; fast function calls between subsystems |
| CFS | Completely Fair Scheduler; red-black tree sorted by vruntime |
| vruntime | Virtual runtime; normalized CPU time to ensure fair sharing |
| Buddy allocator | Physical page allocator; power-of-2 sizes; coalescing reduces fragmentation |
| SLUB | Kernel object allocator; per-type caches with pooling |
| VFS | Virtual File System; uniform file interface for all FS types |
| Netfilter | Kernel packet filtering framework; iptables/nftables hooks into it |
| udev | Userspace device manager; creates /dev nodes, applies udev rules on hotplug |
| eBPF | Safe, sandboxed kernel programs for observability, networking, security |
| LKM | Loadable Kernel Module; driver loaded at runtime without recompile |
| initramfs | In-RAM root filesystem loaded at boot; initializes hardware before mounting real FS |
| systemd | PID 1; parallel service startup, dependency management, cgroup enforcement |
