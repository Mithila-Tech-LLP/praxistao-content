# Chapter 31: The Memory Bus and I/O

The CPU is not an island. It needs to communicate with RAM, with the GPU, with NVMe SSDs, with USB devices, and a hundred other things. All this communication flows through buses — high-speed electrical pathways that connect components on and around the motherboard. The design of these buses — their width, frequency, protocol, and topology — determines much of a system's overall performance. A CPU with a perfect cache hierarchy is still bottlenecked if the memory bus can't deliver data fast enough.

## Table of Contents

1. [What Is a Bus?](#1-what-is-a-bus)
2. [The Memory Bus — CPU to DRAM](#2-the-memory-bus--cpu-to-dram)
3. [DDR SDRAM — Double Data Rate Memory](#3-ddr-sdram--double-data-rate-memory)
4. [PCIe — The Universal I/O Interconnect](#4-pcie--the-universal-io-interconnect)
5. [DMA — Direct Memory Access](#5-dma--direct-memory-access)
6. [NVMe and Storage I/O](#6-nvme-and-storage-io)
7. [System Interconnects — QPI, UPI, Infinity Fabric](#7-system-interconnects--qpi-upi-infinity-fabric)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is a Bus?

A **bus** is a shared communication pathway — a set of electrical wires plus a protocol for who talks when. Think of it like a highway: multiple lanes (wires), cars (data packets) following traffic rules (protocol).

Bus characteristics:
- **Width**: How many bits are transferred in parallel (e.g., 64-bit bus = 8 bytes per transfer)
- **Frequency**: How fast the clock runs (e.g., 3200 MHz)
- **Bandwidth**: Width × Frequency = bytes per second
- **Latency**: How long a single request takes (independent of bandwidth)
- **Protocol**: Rules for who initiates a transfer, addressing, error handling

Early computers used a **shared bus** (a single set of wires connecting CPU, RAM, and all peripherals). Simple, but only one device can transfer at a time — a bottleneck. Modern systems use a **point-to-point** topology: dedicated links between each pair of components. More wires, but no contention.

```
Old shared bus architecture:
  CPU ──────────── Memory Bus (shared) ──────── RAM
                         │
                      PCI Bus ──── GPU, NIC, disk
                      (all sharing the same bus!)

Modern point-to-point topology:
  CPU ── DDR5 channel 0 ── RAM DIMM 0
  CPU ── DDR5 channel 1 ── RAM DIMM 1
  CPU ── PCIe 5.0 x16 ──── GPU
  CPU ── PCIe 5.0 x4 ───── NVMe SSD
  CPU ── USB controller
```

### Quick Check
> 1. Define bus bandwidth and latency. Can you have high bandwidth but high latency? Give an example.
> 2. Why did the industry move from shared buses to point-to-point connections?
> 3. Calculate the bandwidth of a 64-bit bus running at 3200 MHz (in GB/s).

---

## 2. The Memory Bus — CPU to DRAM

The **memory bus** connects the CPU's memory controller to DRAM. In modern CPUs, the memory controller is **integrated on the CPU die itself** (on Intel since 2008 Nehalem, on AMD since 2003 Athlon 64), eliminating the need for a separate Northbridge chip and dramatically reducing latency.

Memory is organized into **channels** — independent memory buses that can be accessed simultaneously. A dual-channel configuration doubles bandwidth over a single channel.

```
DDR5 memory bus example (Intel Alder Lake):
  Channels: 2 (dual-channel)
  Bus width per channel: 64 bits (8 bytes)
  DDR5-4800: 2400 MHz × 2 (DDR) = 4800 MT/s (megatransfers per second)
  Bandwidth per channel: 4800 × 10^6 × 8 bytes = 38.4 GB/s
  Total bandwidth (dual-channel): 76.8 GB/s
```

**Memory ranks and banks**: Each DIMM has multiple ranks (independent sets of chips) and each chip has multiple banks (sub-arrays). These internal structures allow pipelining: while one bank is precharging, another can be accessed. The memory controller orchestrates all of this.

**Memory latency**: The key metrics are:
- **CL (CAS Latency)**: Number of clock cycles from request to first data word
- **tRCD**: Time to activate a row before column can be accessed
- **tRP**: Row precharge time
- Typical DDR5: CL40-40-40-76 at 6000 MHz ≈ 13 ns total latency

Despite high bandwidth, DRAM latency remains ~60–100 ns. Bandwidth improvements from DDR3 → DDR4 → DDR5 are significant; latency improvements are modest.

### Quick Check
> 1. What is a "memory channel" and how does dual-channel increase bandwidth?
> 2. Modern CPUs have the memory controller on the CPU die. What was the old arrangement and why was it slower?
> 3. DDR5-6000 dual-channel: 6000 MT/s × 8 bytes × 2 channels. Calculate total bandwidth in GB/s.

---

## 3. DDR SDRAM — Double Data Rate Memory

**SDRAM (Synchronous Dynamic RAM)**: Memory that operates on a clock signal, synchronizing with the CPU's memory bus. **DDR (Double Data Rate)** transfers data on both the rising AND falling edge of the clock — doubling bandwidth for the same clock frequency.

```
DDR evolution:
  DDR1 (2000): 200–400 MT/s,  3.3V
  DDR2 (2003): 400–1066 MT/s, 1.8V
  DDR3 (2007): 800–2133 MT/s, 1.5V
  DDR4 (2014): 1600–3200 MT/s, 1.2V  ← still common in most desktops/servers
  DDR5 (2020): 3200–8400 MT/s, 1.1V  ← mainstream in new builds (2022+)
  DDR6 (expected ~2025): up to 17600 MT/s
```

Each generation roughly doubles bandwidth and reduces voltage. DDR5 also moves the power management IC (PMIC) onto the DIMM module itself, improving power delivery.

**HBM (High Bandwidth Memory)**: A completely different approach used in GPUs and high-end CPUs. HBM stacks DRAM dies vertically, connected by a wide interface through the package. The interface is thousands of bits wide (4096 bits for HBM2E vs 64 bits for a DDR5 channel).

```
HBM2E: 1024-bit interface × 2 Gbps = 256 GB/s per stack
AMD Instinct MI250X GPU: 8 HBM2E stacks = 3.2 TB/s total memory bandwidth
vs.
DDR5-6000 dual-channel: ~96 GB/s total
```

HBM is expensive and has limited capacity (typically 16–96GB vs 1–4TB for server DDR), but the bandwidth is extraordinary. It's used in AI accelerators, HPC GPUs, and Intel's high-bandwidth Xeon with HBM (Sapphire Rapids HBM).

**LPDDR (Low Power DDR)**: Used in mobile and laptops (Apple M-series uses LPDDR5X, Samsung Galaxy phones use LPDDR5). Lower power through lower voltages and optimized idle states. The memory is often **soldered directly onto the SoC package** (unlike desktop DIMMs that are removable), reducing signal path length and latency.

### Quick Check
> 1. What does "double data rate" mean in DDR memory?
> 2. Why is HBM used in GPUs rather than DDR5?
> 3. Apple's M-series uses LPDDR5X soldered to the SoC. What advantages does this give over DDR5 DIMMs?

---

## 4. PCIe — The Universal I/O Interconnect

**PCIe (Peripheral Component Interconnect Express)** is the dominant standard for connecting high-speed peripherals (GPUs, NVMe SSDs, NICs, FPGAs) to the CPU. It replaced the old PCI, AGP, and ISA buses, which were parallel shared buses. PCIe is a serial, point-to-point, full-duplex protocol.

PCIe uses **lanes** — each lane is a pair of differential signal wires in each direction (4 wires per lane: Tx+, Tx-, Rx+, Rx-). More lanes = more bandwidth.

```
PCIe bandwidth per lane (each direction):
  PCIe 3.0: 985 MB/s per lane  (8 GT/s raw, ~80% efficiency)
  PCIe 4.0: 1.97 GB/s per lane (16 GT/s)
  PCIe 5.0: 3.94 GB/s per lane (32 GT/s)
  PCIe 6.0: 7.88 GB/s per lane (64 GT/s, PAM4 signaling)

Common configurations:
  x1:  1 lane   (cheap NICs, sound cards)
  x4:  4 lanes  (NVMe SSDs)
  x16: 16 lanes (GPUs — the primary GPU slot)

GPU at PCIe 4.0 x16:
  Bandwidth = 16 × 1.97 GB/s × 2 (bidirectional) = 63 GB/s
```

**PCIe topology**: A CPU has a root complex (the PCIe controller). Devices connect directly to the root complex (best latency) or through PCIe switches (fan-out for more devices). The GPU is almost always connected directly to the CPU's root complex.

**PCIe vs. proprietary**: AMD's Infinity Fabric and Intel's UPI are used for CPU-to-CPU links (Section 7). NVLink is NVIDIA's proprietary GPU-to-GPU interconnect (much faster than PCIe for multi-GPU AI training).

### Quick Check
> 1. What is a PCIe lane? What is the bandwidth of a single PCIe 5.0 lane in each direction?
> 2. Why does a GPU use a x16 PCIe slot?
> 3. An NVMe SSD at PCIe 4.0 x4: what is its maximum theoretical bandwidth?

---

## 5. DMA — Direct Memory Access

Without DMA, transferring data from a hard drive to memory works like this: the CPU reads data from the drive's I/O port one byte at a time, then writes each byte to memory — **polling** the drive. During this time, the CPU is 100% occupied with the transfer and cannot do anything else. This is called **programmed I/O (PIO)**.

**DMA (Direct Memory Access)** offloads this to dedicated hardware — a **DMA controller**. The CPU programs the DMA controller with:
1. Source address (e.g., a disk sector)
2. Destination address (e.g., a buffer in RAM)
3. Transfer size

The DMA controller then performs the transfer independently, while the CPU is free to do other work. When done, the DMA controller signals the CPU via an **interrupt** (Chapter 32).

```
DMA transfer sequence:
  CPU → DMA controller: "Copy 512 bytes from disk sector 42 to memory 0x1F000"
  CPU → does other work
  DMA controller → requests memory bus → transfers data
  DMA controller → asserts interrupt line
  CPU: "Interrupt! DMA done" → processes the transferred data
```

**Bus mastering**: Modern DMA is implemented as PCIe **bus mastering** — the PCIe device (GPU, NIC, NVMe SSD) can initiate memory reads/writes directly. The device contains its own DMA engine. This is how a GPU can copy a model's weights to VRAM without involving the CPU at all.

**IOMMU (I/O Memory Management Unit)**: Just as the MMU protects process memory from each other, the IOMMU protects system memory from rogue devices. It translates device-physical addresses to host-physical addresses and can prevent a malicious device from DMA'ing into arbitrary memory regions. Essential for passing PCIe devices through to virtual machines (VT-d on Intel, AMD-Vi on AMD).

### Quick Check
> 1. What is the problem with Programmed I/O (PIO)?
> 2. Describe DMA in one paragraph. What are the three things the CPU tells the DMA controller?
> 3. What is the IOMMU and why is it important for virtualization?

---

## 6. NVMe and Storage I/O

Traditional hard drives (HDDs) use spinning magnetic platters — access times of 5–10 ms. SSDs use NAND flash — access times of 50–100 µs. But early SSDs were connected via SATA, an interface designed for slow HDDs with a theoretical maximum of 600 MB/s — a bottleneck.

**NVMe (Non-Volatile Memory Express)** is a protocol designed specifically for SSDs, using PCIe directly. It removes the SATA bottleneck and takes advantage of SSD's internal parallelism.

```
Storage interface comparison:
  SATA III: up to 600 MB/s, HDD-era protocol
  SATA SSD: typically 500–550 MB/s sequential read
  NVMe PCIe 3.0 x4: up to 3.5 GB/s sequential read
  NVMe PCIe 4.0 x4: up to 7 GB/s sequential read
  NVMe PCIe 5.0 x4: up to 14 GB/s sequential read

  Samsung 990 Pro (PCIe 4.0): 7.45 GB/s read, 6.9 GB/s write
```

**NVMe queue depth**: SATA can handle 1 command at a time with a queue depth of 32. NVMe supports 64K commands per queue and 64K queues — designed for the parallelism of NAND flash and for server workloads with many concurrent I/O requests.

**Computational storage**: Some cutting-edge NVMe SSDs (e.g., Samsung SmartSSD, ScaleFlux) include an FPGA or ARM processor onboard. Instead of reading raw data to the host CPU for filtering/compression, the computation runs inside the drive. This reduces the amount of data transferred over PCIe — critical for analytics workloads with huge datasets.

### Quick Check
> 1. Why was SATA a bottleneck for SSDs even though SSDs are much faster than HDDs?
> 2. NVMe PCIe 4.0 x4 theoretical bandwidth: 4 × 1.97 GB/s. Is this consistent with the 7 GB/s quoted above? Why might they differ?
> 3. What is "computational storage" and what problem does it solve?

---

## 7. System Interconnects — QPI, UPI, Infinity Fabric

In multi-socket servers (two or more CPUs sharing one system), the CPUs must communicate with each other. This requires a high-speed CPU-to-CPU interconnect.

**Intel UPI (Ultra Path Interconnect)**: Previously called QPI (QuickPath Interconnect). Used in Intel Xeon server processors. UPI is a high-speed point-to-point serial link between CPUs.

```
Intel Xeon Scalable (Sapphire Rapids): 3 UPI links at 16 GT/s each
  Bandwidth: 3 × 16 GT/s × 2B/transfer ≈ 96 GB/s CPU-to-CPU bandwidth
```

**AMD Infinity Fabric**: AMD's proprietary interconnect used within a single CPU die (connecting CCDs — core complex dies — to the I/O die) and between CPUs in multi-socket servers. In Ryzen desktop CPUs, Infinity Fabric also connects the CCDs to the central I/O die that contains the memory controller.

```
AMD EPYC Genoa (Zen 4): 3 Infinity Fabric links to other CPUs (xGMI)
  Also connects up to 12 CCDs internally via the I/O die
```

**NUMA (Non-Uniform Memory Access)**: In a multi-socket system, each CPU has "local" DRAM (connected directly to that CPU's memory controller) and "remote" DRAM (connected to the other CPU's memory controller, accessible via the interconnect). Accessing local memory is fast; accessing remote memory crosses the interconnect and takes 2–3× longer.

```
NUMA topology (2-socket server):
  CPU 0 ── local memory (64GB, 76 GB/s bandwidth)
  CPU 0 ── UPI → CPU 1 → remote memory (64GB, ~35 GB/s effective)

  OS must place memory near the CPU that uses it (NUMA-aware allocation)
  to avoid bottlenecking on the UPI interconnect.
```

Modern operating systems (Linux, Windows Server) are NUMA-aware and try to allocate memory on the same NUMA node as the thread accessing it. Databases and HPC applications often do explicit NUMA placement (`numactl`).

### Quick Check
> 1. What is the purpose of Intel UPI / AMD Infinity Fabric in a multi-socket server?
> 2. What is NUMA and why does accessing "remote" memory take longer?
> 3. Why do Linux and databases care about NUMA topology?

---

## Summary

- Buses have **bandwidth** (data per second) and **latency** (time per request). Modern systems use point-to-point links instead of shared buses.
- The **memory bus** (DDR5, HBM) connects the on-chip memory controller to DRAM. Bandwidth scales with DDR generation and channel count. HBM offers extreme bandwidth for GPUs.
- **PCIe** is the universal high-speed I/O interconnect for GPUs, NVMe SSDs, and NICs. Higher gen numbers and more lanes = more bandwidth.
- **DMA** offloads data transfers from the CPU to dedicated hardware (DMA controllers, PCIe bus masters). The CPU is notified via interrupt when done. The IOMMU protects memory from rogue devices.
- **NVMe over PCIe** is the modern SSD interface — up to 14 GB/s vs 600 MB/s for SATA.
- **CPU-to-CPU interconnects** (Intel UPI, AMD Infinity Fabric) enable multi-socket servers and create NUMA topologies where local memory is faster than remote memory.

---

## Exercises

### Easy
1. A DDR5-6400 dual-channel system: each channel is 6400 MT/s and 64 bits wide. Calculate total system memory bandwidth in GB/s.
2. What is DMA? Why is it better than having the CPU do memory transfers directly (Programmed I/O)?
3. An NVMe SSD at PCIe 4.0 x4: calculate maximum theoretical bandwidth. What would it be at PCIe 5.0 x4?

### Medium
4. A GPU has 16GB of GDDR6X memory at 1008 GB/s and is connected to the CPU via PCIe 4.0 x16 (64 GB/s). The GPU needs to process a 10GB dataset. (a) What is the bottleneck: GPU memory bandwidth or PCIe bandwidth? (b) How long does it take to transfer the dataset from CPU RAM to GPU VRAM? (c) How long would PCIe 5.0 x16 take?
5. A 2-socket server uses Intel Xeon with UPI at 96 GB/s CPU-to-CPU bandwidth. NUMA node 0 has 256GB RAM at 200 GB/s, NUMA node 1 has 256GB RAM at 200 GB/s. A workload running on CPU 0 randomly accesses all 512GB of memory. What fraction of accesses are remote? What is the effective memory bandwidth available to the workload?
6. Explain bus mastering in PCIe. When a GPU's DMA engine wants to copy a 1GB buffer from system RAM to GPU VRAM: (a) Who initiates the transfer? (b) How does the CPU know when it's done? (c) What role does the IOMMU play?

### Hard
7. Memory bandwidth vs latency: An HPC application alternates between two access patterns: (a) sequential array scan (bandwidth-limited), (b) random pointer chasing (latency-limited). Compare how DDR5-6400 dual-channel vs HBM3 would perform for each pattern. HBM3: 1.2 TB/s bandwidth, 100ns latency. DDR5-6400 dual-channel: 100 GB/s bandwidth, 70ns latency. Which system wins for which workload and why?
8. Design a PCIe bandwidth measurement experiment. You want to measure: (a) peak PCIe read bandwidth (CPU reading from GPU buffer), (b) peak PCIe write bandwidth (CPU writing to GPU buffer), (c) PCIe latency (round-trip time for a 1-byte transfer). Describe the specific operations in CUDA or OpenCL you would use, and explain why the peak bandwidth might be less than the theoretical maximum.
