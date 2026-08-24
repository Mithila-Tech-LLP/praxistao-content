# Chapter 36: Virtualization and Containers

> **"Virtualization is how one physical computer pretends to be many. Containers are how one OS instance appears as many isolated environments. Both solve the same business problem — efficient, isolated, reproducible environments — but at completely different layers. Understanding both is essential for modern systems."**

---

## Table of Contents

1. [What Is Virtualization?](#1-what-is-virtualization)
2. [Type 1 Hypervisors (Bare Metal)](#2-type-1-hypervisors-bare-metal)
3. [Type 2 Hypervisors (Hosted)](#3-type-2-hypervisors-hosted)
4. [Hardware-Assisted Virtualization (VT-x/AMD-V)](#4-hardware-assisted-virtualization-vt-xamd-v)
5. [Memory Virtualization — EPT and Shadow Page Tables](#5-memory-virtualization--ept-and-shadow-page-tables)
6. [I/O Virtualization](#6-io-virtualization)
7. [Containers vs VMs](#7-containers-vs-vms)
8. [Linux Containers — Namespaces and cgroups](#8-linux-containers--namespaces-and-cgroups)
9. [Docker Architecture](#9-docker-architecture)
10. [Security Isolation Comparison](#10-security-isolation-comparison)
11. [Summary](#summary)

---

## 1. What Is Virtualization?

**Virtualization** creates a software abstraction of hardware — a **virtual machine (VM)** that runs a guest OS as if it had its own dedicated hardware.

**Why virtualization?**
```
Without virtualization:
  1 physical server → runs 1 OS → runs a few applications
  Server CPU utilization: 5-15% (mostly idle!)
  → Expensive, wasteful
  
With virtualization:
  1 physical server → hypervisor → runs 10-50 VMs
  Each VM: 1 OS + its applications
  Server CPU utilization: 60-80%
  → Server consolidation: fewer physical machines, lower cost

Additional benefits:
  Isolation: one VM crash doesn't affect others
  Snapshot: save/restore entire VM state
  Migration: move running VM between physical servers (live migration)
  Testing: safely test malware in a VM sandbox
```

**Hypervisor (Virtual Machine Monitor, VMM):**
Software that creates and manages virtual machines. Provides each VM with:
- Virtual CPUs (vCPUs)
- Virtual RAM
- Virtual network interfaces
- Virtual disks (files on host, appear as real disks to guest)

---

## 2. Type 1 Hypervisors (Bare Metal)

**Type 1:** Hypervisor runs DIRECTLY on the hardware, no host OS:

```
┌─────────────────────────────────────────────────┐
│ VM 1: Linux    │ VM 2: Windows │ VM 3: FreeBSD   │
│ OS + Apps      │ OS + Apps     │ OS + Apps        │
├─────────────────────────────────────────────────┤
│                Hypervisor (Type 1)               │
│          (bare metal, no host OS)               │
├─────────────────────────────────────────────────┤
│              Physical Hardware                   │
└─────────────────────────────────────────────────┘
```

**Examples:**
- **VMware ESXi**: enterprise server virtualization; most popular data center hypervisor
- **Microsoft Hyper-V**: built into Windows Server; also used in Azure
- **Xen**: open-source; used by AWS EC2
- **KVM (Kernel-based Virtual Machine)**: built into Linux kernel; used by GCP, many cloud providers

**KVM (Linux's Type 1):**
Technically a hybrid: KVM is a kernel module that turns the Linux kernel into a Type 1 hypervisor:
```bash
# Check if KVM is available:
ls /dev/kvm   # exists if CPU supports VT-x/AMD-V

# Create a VM with QEMU+KVM:
qemu-system-x86_64 \
    -enable-kvm \                      # use hardware virtualization
    -cpu host \                        # expose host CPU features to guest
    -m 2048 \                          # 2GB RAM
    -smp 2 \                           # 2 vCPUs
    -drive file=ubuntu.qcow2 \         # disk image
    -net user,hostfwd=tcp::2222-:22   # NAT network, forward port 2222 to guest SSH
```

**KVM architecture:**
```
Host Linux kernel with KVM module:
  /dev/kvm → userspace API for managing VMs
  QEMU process: uses /dev/kvm ioctls to:
    - Create VM: KVM_CREATE_VM
    - Create vCPU: KVM_CREATE_VCPU
    - Set memory regions: KVM_SET_USER_MEMORY_REGION
    - Run vCPU: KVM_RUN (puts vCPU into guest mode)
    - Handle VM exits (emulate I/O, handle MMIO)
```

---

## 3. Type 2 Hypervisors (Hosted)

**Type 2:** Hypervisor runs as an application on top of an existing host OS:

```
┌─────────────────────────────────────────────────┐
│ VM 1: Ubuntu        │ VM 2: Windows XP          │
│ Guest OS + Apps     │ Guest OS + Apps             │
├─────────────────────────────────────────────────┤
│        VMware Workstation / VirtualBox           │
├─────────────────────────────────────────────────┤
│          Host OS (Windows / macOS / Linux)       │
├─────────────────────────────────────────────────┤
│              Physical Hardware                   │
└─────────────────────────────────────────────────┘
```

**Examples:**
- **VMware Workstation/Fusion**: popular for development
- **VirtualBox**: open-source, cross-platform
- **Parallels Desktop**: macOS-focused, good Apple Silicon support
- **QEMU** (without KVM): pure software emulation

**Type 1 vs Type 2:**
| Feature | Type 1 | Type 2 |
|---------|--------|--------|
| Performance | Better (direct hardware access) | Overhead from host OS |
| Use case | Production data center | Development, testing, desktop |
| Security | Better (smaller attack surface) | Host OS vulnerabilities affect all VMs |
| Ease of setup | Harder | Easy (install like an app) |

---

## 4. Hardware-Assisted Virtualization (VT-x/AMD-V)

Without hardware assistance, hypervisors had to use **binary translation** — rewriting all privileged guest instructions on the fly. This was slow.

**Intel VT-x and AMD-V** add hardware support directly to the CPU:

```
CPU operating modes:
  VMX Root mode (host mode):     Hypervisor runs here
  VMX Non-Root mode (guest mode): Guest OS runs here

VMX instructions:
  VMXON:    Enable VMX operation
  VMLAUNCH: Enter guest for first time
  VMRESUME: Return to guest after VM exit
  VMPTRLD:  Load VMCS for this vCPU
  VMREAD:   Read from VMCS
  VMWRITE:  Write to VMCS
  VMXOFF:   Disable VMX operation

VMCS (Virtual Machine Control Structure):
  Per-vCPU data structure controlling what causes VM exits and what state is saved
```

**VM Exit and VM Entry:**
```
VM Entry (hypervisor → guest):
  1. CPU loads guest state from VMCS (registers, RFLAGS, segment registers, CR3)
  2. CPU begins executing guest code in VMX Non-Root mode
  3. Guest OS runs as if on real hardware
  
VM Exit (guest → hypervisor):
  Triggered by: privileged instruction in guest, I/O access, interrupt
  1. CPU saves guest state to VMCS
  2. CPU loads host state from VMCS
  3. CPU begins executing hypervisor exit handler
  4. Hypervisor handles: emulate I/O, inject interrupt, etc.
  5. VM Entry to return to guest
```

---

## 5. Memory Virtualization — EPT and Shadow Page Tables

**The challenge:**
Guest OS manages its own page tables (Guest Virtual Address → Guest Physical Address).
But Guest Physical Address is not really physical — it's a virtual address in the hypervisor.

We need two-level translation:
```
Guest Virtual → [Guest page tables] → Guest Physical → [Hypervisor] → Host Physical
```

**Two approaches:**

**Shadow page tables (pre-hardware support):**
Hypervisor maintains "shadow" page tables that map GVA directly to HPA:
```
Guest writes to its page table → trap to hypervisor → hypervisor updates shadow PT
→ Expensive: every guest page table modification = VM exit + shadow PT update
```

**EPT/NPT — Extended Page Tables / Nested Page Tables (hardware):**
Hardware walks TWO levels of page tables:
1. Guest page table: GVA → GPA (walked by hardware automatically)
2. EPT (maintained by hypervisor): GPA → HPA

```
MOV RAX, [GVA]
  → hardware walks guest page tables: GVA → GPA
  → hardware walks EPT: GPA → HPA
  → actual memory access at HPA
  
Only 1 VM exit if EPT entry is missing (EPT violation)
Otherwise: transparent hardware translation, no hypervisor involvement!
```

---

## 6. I/O Virtualization

**I/O is the hardest part to virtualize:**

**Emulation (slowest):**
```
Guest writes to I/O port → VM exit → hypervisor emulates the device in software
(QEMU emulates thousands of devices: IDE disk, E1000 NIC, VGA, etc.)
Very compatible but slow: every I/O operation = VM exit
```

**Para-virtualization (virtio):**
Guest OS knows it's in a VM and uses special drivers optimized for virtual environments:
```
virtio: standardized device interface for VMs

virtio-blk:  Virtual block device (disk) — much faster than emulated IDE
virtio-net:  Virtual network card — much faster than emulated E1000
virtio-gpu:  Virtual GPU for display
virtio-fs:   Shared filesystem between host and guest

Guest: uses virtio driver → sends commands to hypervisor via shared ring buffers
Hypervisor: processes batch of I/O requests without VM exits per request
```

**SR-IOV — PCIe Single Root I/O Virtualization:**
A physical NIC can present itself as multiple independent virtual NICs:
```
Physical NIC → Physical Function (PF, managed by hypervisor)
            → Virtual Function 1 (VF1, assigned to VM 1)
            → Virtual Function 2 (VF2, assigned to VM 2)
            → ...

VM gets DIRECT hardware access to its VF — zero hypervisor in the data path!
Achieves near-native performance for network/storage I/O
```

---

## 7. Containers vs VMs

```
Virtual Machine:                    Container:
┌────────────────────┐              ┌────────────────────┐
│ Application        │              │ Application        │
│ Libraries/Runtime  │              │ Libraries/Runtime  │
│ Guest OS kernel    │              │ ─────────────────  │
│ Hypervisor emulated│              │ (No OS kernel!     │
│ hardware           │              │  shares host kernel│
├────────────────────┤              ├────────────────────┤
│ Hypervisor         │              │ Container Runtime  │
│ (KVM, VMware, etc) │              │ (Docker, containerd│
├────────────────────┤              ├────────────────────┤
│ Host OS Kernel     │              │ Host OS Kernel     │
├────────────────────┤              ├────────────────────┤
│ Hardware           │              │ Hardware           │
└────────────────────┘              └────────────────────┘
```

| Feature | VM | Container |
|---------|-----|-----------|
| Isolation | Strong (separate OS) | Weaker (shared kernel) |
| Boot time | Minutes | Milliseconds |
| Size | GBs (full OS) | MBs (just app + libs) |
| Performance | Near-native (with VT-x) | Native (no hypervisor overhead) |
| Security | Strong | Weaker (kernel vuln affects all containers) |
| Density | ~10-50 VMs per server | ~100-1000 containers per server |
| Kernel bugs | Guest kernel can be different | All share host kernel vulnerability |
| Use case | Strong isolation, different OS | Microservices, dev/prod parity |

---

## 8. Linux Containers — Namespaces and cgroups

Linux containers are built from two kernel features:

**Namespaces — isolation:**
Each namespace type provides isolation for one aspect of the system:

```
pid namespace:     Container sees its own process tree (container's init = PID 1)
mount namespace:   Container has its own file system view (separate root FS)
net namespace:     Container has its own network stack (separate IP, routes, sockets)
user namespace:    Container has its own UID/GID mapping (uid 0 inside = non-root outside)
ipc namespace:     Separate System V IPC and POSIX message queues
uts namespace:     Separate hostname and domainname
cgroup namespace:  Separate view of cgroup hierarchy
time namespace:    Different clock offsets (new, kernel 5.6+)
```

**cgroups (control groups) — resource limits:**
```
Limit how much CPU, RAM, I/O, network a group of processes can use:

cpu:     CPU time (shares, quotas, periods)
memory:  RAM + swap limits (hard limits, soft limits, OOM killer behavior)
blkio:   Block I/O bandwidth and IOPS limits
net_cls: Tag network packets with cgroup ID (for traffic shaping)
devices: Allow/deny access to device files

Example: limit nginx to 50% CPU and 512MB RAM:
  /sys/fs/cgroup/cpu/nginx/cpu.cfs_period_us = 100000  (100ms)
  /sys/fs/cgroup/cpu/nginx/cpu.cfs_quota_us  = 50000   (50ms/period = 50%)
  /sys/fs/cgroup/memory/nginx/memory.limit_in_bytes = 536870912  (512MB)
```

---

## 9. Docker Architecture

**Docker** packages applications with their dependencies into **images** (immutable snapshots) and runs them as containers.

```
Docker architecture:

User types: docker run ubuntu bash
       ↓
Docker CLI → REST API → Docker daemon (dockerd)
       ↓
dockerd → containerd (container runtime)
       ↓
containerd → runc (low-level OCI runtime)
       ↓
runc:
  1. Create new namespaces (pid, mount, net, uts, ipc, user)
  2. Set up cgroups (resource limits)
  3. Mount overlay filesystem (container's root FS)
  4. Drop privileges (capabilities)
  5. Apply seccomp filter
  6. exec() the process inside the container
```

**Container file system (overlay):**
Docker uses **OverlayFS** — a union file system that layers read-only image layers with a writable container layer:
```
Container overlay:
  Upper (writable): [ container's writes ]
  Lower (read-only): [ layer 3: app ]
                     [ layer 2: python ]
                     [ layer 1: ubuntu base ]

When container writes to a file in lower layers:
  Copy-on-Write: copy file to upper layer, then write to upper copy
  Lower layers unchanged

Image layers are shared between ALL containers using that image → saves disk space
```

```bash
# Docker commands:
docker pull ubuntu:22.04          # download image
docker run -it ubuntu:22.04 bash  # start container with interactive shell
docker ps                         # list running containers
docker exec -it <id> bash        # exec into running container
docker images                     # list images
docker inspect <id>               # detailed container info (namespaces, mounts, etc.)
```

---

## 10. Security Isolation Comparison

```
VMs (strongest isolation):
  Attack: kernel exploit in a guest VM
  Blast radius: only that guest VM — hypervisor is not affected
  Why: guest has a DIFFERENT kernel; hypervisor enforces hardware boundaries
  
Standard containers (weaker):
  Attack: kernel exploit via a container
  Blast radius: ALL containers on the same host
  Why: all containers share the host kernel — kernel exploit = root on host

Sandboxed containers (middle ground):
  gVisor (Google): User-space kernel interception
    - Every syscall from container → goes to gVisor's user-space kernel
    - gVisor reimplements Linux kernel subset in Go
    - Real host kernel sees only gVisor's safe subset of syscalls
    
  Kata Containers:
    - Each container runs in a lightweight VM (using KVM)
    - Combines container speed with VM isolation
    - VM overhead: ~150ms startup (vs 1ms for runc)
    
  Firecracker (AWS Lambda):
    - Ultra-lightweight VMM (virtual machine monitor) in Rust
    - Starts in < 125ms
    - Powers AWS Lambda and AWS Fargate
```

---

## Summary

| Concept | Description |
|---------|------------|
| Hypervisor | Software that creates and manages virtual machines |
| Type 1 | Bare-metal hypervisor; runs directly on hardware (KVM, ESXi, Xen, Hyper-V) |
| Type 2 | Hosted hypervisor; runs on a host OS (VirtualBox, VMware Workstation) |
| VT-x/AMD-V | CPU hardware for efficient virtualization; VMX root/non-root modes |
| VMCS | Per-vCPU structure controlling what triggers VM exits and what state is saved |
| VM Exit | Trap from guest to hypervisor (on privileged instruction, I/O, etc.) |
| EPT/NPT | Hardware two-level page tables for memory virtualization; no shadow PT needed |
| virtio | Paravirtual device interface; much faster than emulated legacy devices |
| SR-IOV | Physical device with multiple virtual functions; VMs get direct hardware access |
| Container | Isolated process group sharing host kernel; fast, lightweight |
| Namespace | Linux kernel isolation: pid, mount, net, user, uts, ipc, cgroup |
| cgroup | Linux kernel resource limits: CPU, memory, I/O per process group |
| OverlayFS | Union file system: layered read-only image + writable container layer |
| runc | Low-level OCI container runtime that does the actual namespace/cgroup setup |
| gVisor | User-space kernel interception for stronger container isolation |
| Kata Containers | Containers in lightweight VMs; combines speed and isolation |
