# Chapter 31: Disk I/O and Disk Scheduling

> **"A hard drive head moves at 5,400 or 7,200 RPM. Getting the seek order right means the difference between 50 IOPS and 150 IOPS. A naive sequential request queue is left on the table. But SSDs and NVMe have no head — and understanding that changes everything about how the OS should schedule I/O."**

---

## Table of Contents

1. [How Hard Disks Work](#1-how-hard-disks-work)
2. [HDD Performance Math](#2-hdd-performance-math)
3. [The Block I/O Stack](#3-the-block-io-stack)
4. [Disk Scheduling Algorithms](#4-disk-scheduling-algorithms)
5. [NVMe and SSDs — A Different World](#5-nvme-and-ssds--a-different-world)
6. [Linux I/O Schedulers](#6-linux-io-schedulers)
7. [RAID — Combining Multiple Disks](#7-raid--combining-multiple-disks)
8. [Disk Caching and Write-Back](#8-disk-caching-and-write-back)
9. [Summary](#summary)

---

## 1. How Hard Disks Work

**Mechanical structure:**
```
HDD internals:
┌──────────────────────────────────────┐
│  Platter (magnetic disk, coated)     │
│  ┌────────────────────────────────┐  │
│  │  Track 0 (outermost)          │  │
│  │    ┌──────────────────────┐   │  │
│  │    │   Track 100          │   │  │
│  │    │   ...                │   │  │
│  │    └──────────────────────┘   │  │
│  │  Track N (innermost)          │  │
│  └────────────────────────────────┘  │
│                                      │
│  Read/Write head (on actuator arm)   │
│  → moves radially across platters    │
└──────────────────────────────────────┘

Multiple platters stacked vertically
Multiple heads (one per platter surface)
All heads move together (same cylinder)
```

**Addressing data:**
```
CHS (Cylinder-Head-Sector) — legacy:
  Cylinder: which track (0 to 4999 for a typical drive)
  Head:     which platter surface (0-7)
  Sector:   which 512-byte unit on the track (1-63)

LBA (Logical Block Addressing) — modern:
  Single sequential number: LBA 0, 1, 2, 3, ...
  Drive firmware maps LBA → physical location
  OS only sees LBA numbers
```

---

## 2. HDD Performance Math

Three components of I/O latency:

**1. Seek time:** Head moves to the correct cylinder.
```
Average seek time: 5-12ms for 7200 RPM drives
Full stroke (inner to outer track): ~15-20ms
Adjacent track: ~0.3ms
Typical average: ~8ms for random access patterns
```

**2. Rotational latency:** Disk rotates until the target sector is under the head.
```
7200 RPM = 120 rotations/second
One full rotation = 1/120 = 8.33ms
Average rotational latency = half rotation = 4.17ms
```

**3. Transfer time:** Time to read/write the actual data.
```
7200 RPM drive outer tracks: ~200 MB/s sequential
For 4KB block: 4096 / (200 × 10^6) = 0.02ms (negligible)
```

**Total I/O time (random 4KB read on HDD):**
```
Seek:     8ms
Rotation: 4ms
Transfer: 0.02ms
Total:    ~12ms
```

**IOPS (I/O Operations Per Second):**
```
Random IOPS = 1000ms / 12ms = ~83 IOPS (for 7200 RPM HDD)
Sequential throughput = 200 MB/s (bypasses seek/rotation overhead)
```

**Key insight:** For HDDs, **seek order matters enormously**. Reading 1,000 random blocks scattered across the disk takes 1000 × 12ms = 12 seconds. Reading them in order minimizes seeks.

---

## 3. The Block I/O Stack

Linux's block I/O architecture:

```
Application
    │ read() / write()
    ▼
File System (ext4, etc.)
    │ Translates file offsets → block numbers
    │ Creates bio (block I/O descriptor)
    ▼
Block Layer — General
    │ Merging: adjacent bios combined into one larger request
    │ Plugging: queue requests briefly to allow merging
    ▼
I/O Scheduler
    │ Reorders requests to minimize seek time
    ▼
Block Device Driver (SCSI, NVMe, etc.)
    │ Converts block layer requests to hardware commands
    ▼
Hardware (HDD, SSD, NVMe drive)
    │ DMA transfer to/from RAM
    ▼
Interrupt: I/O complete
    │ Block layer processes completion
    ▼
File System: data available
    │
    ▼
Application: read() returns
```

**`struct bio` (block I/O descriptor):**
```c
struct bio {
    sector_t    bi_sector;    // starting sector (LBA)
    struct block_device *bi_bdev;  // target block device
    unsigned short bi_vcnt;   // number of scatter-gather segments
    struct bio_vec *bi_io_vec;  // array of {page, offset, len} — where data goes in RAM
    bio_end_io_t *bi_end_io;  // completion callback
    void *bi_private;
    // ...
};
```

---

## 4. Disk Scheduling Algorithms

**Goal:** Reorder the queue of pending I/O requests to minimize total seek time.

### FCFS (First Come, First Served)

No reordering — process requests in order received:
```
Head position: track 50
Request queue: 98, 183, 37, 122, 14, 124, 65, 67

Seek sequence: 50→98→183→37→122→14→124→65→67
Total distance: 48+85+146+85+108+110+59+2 = 643 tracks

Poor performance — lots of back-and-forth (thrashing)
```

### SSTF (Shortest Seek Time First)

Always service the request closest to the current head position:
```
Head position: 50
Queue: 98, 183, 37, 122, 14, 124, 65, 67

Step 1: Closest to 50 → 37 (distance 13)  → seek to 37
Step 2: Closest to 37 → 14 (distance 23)  → seek to 14
Step 3: Closest to 14 → 65 (distance 51)  → seek to 65
Step 4: Closest to 65 → 67 (distance 2)   → seek to 67
Step 5: Closest to 67 → 98 (distance 31)  → seek to 98
Step 6: Closest to 98 → 122 (distance 24) → seek to 122
Step 7: Closest to 122 → 124 (distance 2) → seek to 124
Step 8: 183 (only one left)               → seek to 183

Total distance: 13+23+51+2+31+24+2+59 = 205 tracks (vs 643 for FCFS)
```

**Problem with SSTF:** **Starvation**. A request far from the current head might wait forever if new nearby requests keep arriving.

### SCAN (Elevator Algorithm)

Head sweeps back and forth, servicing requests in the direction of movement:
```
Head position: 50, moving toward higher tracks
Queue: 98, 183, 37, 122, 14, 124, 65, 67

Moving up (toward 199):
  50 → 65 → 67 → 98 → 122 → 124 → 183 → 199 (end)
Moving down (toward 0):
  199 → 37 → 14 → 0 (end)

Total distance: 15+2+31+24+2+59+16+162+23+14 = 348 tracks
```

No starvation: every request is serviced within one full sweep.

### C-SCAN (Circular SCAN)

Like SCAN but only services requests in one direction; on reaching the end, jumps back to beginning:
```
Head moves 50 → 199 (services requests going up)
Then: jumps from 199 → 0 (no servicing)
Then moves 0 → 199 again (services remaining requests)

More uniform wait times than SCAN
```

### LOOK and C-LOOK

Like SCAN/C-SCAN but doesn't go all the way to disk end — only as far as the last request in each direction. More efficient than SCAN in practice.

**Comparison:**
| Algorithm | Total distance | Starvation risk | Fairness |
|-----------|---------------|-----------------|---------|
| FCFS | Highest | None | Perfect |
| SSTF | Low | High | Poor |
| SCAN | Medium | None | Good |
| C-SCAN | Medium | None | Better (more uniform) |
| C-LOOK | Low | None | Good |

---

## 5. NVMe and SSDs — A Different World

**SSD (Solid State Drive):** Stores data in NAND flash memory cells, not on magnetic platters.

**Key differences from HDD:**
```
                HDD          SSD (SATA)    NVMe (PCIe)
Seek time:      5-12ms       0.1ms         0.02ms
Sequential R:   100-200MB/s  500-600MB/s   3,500-7,000MB/s
Random IOPS:    80-100       50,000-90,000 500,000-1,000,000
Latency:        ~10ms        ~0.1ms        ~0.02ms
```

**SSD internal structure:**
```
NVMe SSD:
  Controller chip (CPU + DRAM cache)
  ↓ (manages flash translation layer, wear leveling, garbage collection)
  Die 0: Plane 0, Plane 1 (each with thousands of blocks)
  Die 1: ...
  (Each block: ~512 pages of 4KB = 2MB per block)
```

**Flash cell types:**
```
SLC (Single Level Cell): 1 bit per cell — fastest, most durable, expensive
MLC (Multi Level Cell):  2 bits per cell — balanced
TLC (Triple Level Cell): 3 bits per cell — slower, less durable, cheaper (consumer drives)
QLC (Quad Level Cell):   4 bits per cell — densest, slowest, least durable
```

**Flash limitations:**
1. **Erase before write:** Flash must be erased in large blocks (128-256 pages) before pages can be rewritten. Can't overwrite a page in place.
2. **Write amplification:** A small update can force erasing a large block and rewriting everything.
3. **Wear:** Each cell can only be erased ~100-10,000 times (QLC: ~100, SLC: ~100,000).

**Flash Translation Layer (FTL):**
The SSD controller hides all this complexity with a **FTL**:
- **Logical-to-physical mapping:** Like a page table for flash blocks
- **Wear leveling:** Distributes writes evenly across all blocks to prevent early wear
- **Garbage collection:** Collects "dirty" blocks (data overwritten) and reclaims space
- **TRIM command:** OS tells SSD which LBA ranges are deleted → SSD can erase them proactively

**Why seek scheduling is mostly irrelevant for SSDs:**
- No mechanical head movement = no seek time
- Random IOPS ≈ Sequential IOPS (within ~2×)
- The FTL already parallelizes I/O across multiple flash dies
- For NVMe: the hardware has 64,000+ queues with 65,535 entries each — can handle massive parallelism

---

## 6. Linux I/O Schedulers

**For HDDs:**

**BFQ (Budget Fair Queueing) — default on HDDs:**
- Assigns I/O bandwidth budgets to processes (similar to CFS for CPU)
- Prioritizes interactive processes (low-latency for desktop use)
- Reduces latency at the cost of some throughput
```bash
echo bfq > /sys/block/sda/queue/scheduler
# cat /sys/block/sda/queue/scheduler
# [bfq] mq-deadline none
```

**mq-deadline:**
- Deadline scheduling: each request has a deadline (read: 500ms, write: 5s default)
- Maintains sorted queues (like SSTF) but promotes old requests near deadline
- Good for databases and latency-sensitive workloads
```bash
echo mq-deadline > /sys/block/sda/queue/scheduler
```

**For SSDs/NVMe:**

**none (no-op):**
- No scheduling — just submit requests in order received
- For NVMe SSDs: the hardware handles its own scheduling efficiently
```bash
echo none > /sys/block/nvme0n1/queue/scheduler
```

**Tuning:**
```bash
# Read-ahead for sequential workloads:
echo 2048 > /sys/block/sda/queue/read_ahead_kb  # 2MB readahead

# Queue depth (number of concurrent requests to the device):
cat /sys/block/nvme0n1/queue/nr_requests  # NVMe: 1023 (deep queue)
cat /sys/block/sda/queue/nr_requests      # HDD: typically 128

# Enable TRIM for SSDs:
# Add "discard" mount option in /etc/fstab:
# /dev/nvme0n1p1 / ext4 defaults,discard 0 1
# Or run fstrim periodically:
fstrim -v /
```

---

## 7. RAID — Combining Multiple Disks

**RAID (Redundant Array of Independent Disks)** combines multiple disks for:
- **Performance:** read/write data in parallel across drives
- **Redundancy:** survive drive failures without data loss

**RAID levels:**

**RAID 0 — Striping (no redundancy):**
```
Data:          [A][B][C][D][E][F]
Disk 0:        [A][C][E]         ← even chunks
Disk 1:        [B][D][F]         ← odd chunks

Read speed: 2×  Write speed: 2×  Fault tolerance: NONE
If either disk fails → all data lost!
Use case: video editing where speed matters more than safety
```

**RAID 1 — Mirroring:**
```
Data:          [A][B][C]
Disk 0:        [A][B][C]  (identical copy)
Disk 1:        [A][B][C]  (identical copy)

Read speed: 2× (can read from either)  Write speed: 1× (must write both)
Fault tolerance: survive 1 disk failure
Use case: OS drive, database transaction logs
```

**RAID 5 — Striping with distributed parity:**
```
Data:     [A][B][C][D]  [E][F][G][H]
Disk 0:   [A]   [C]  P  [E]   [G]
Disk 1:   [B]   [C]     [F]   [G]   P
Disk 2:   [P]   [D]     [P]   [H]

(P = parity block = XOR of the data blocks in that stripe)

Read speed: (n-1)× where n = drives
Write speed: slower than RAID 0 (parity must be updated)
Fault tolerance: survive 1 disk failure (rebuild from parity)
Space efficiency: (n-1)/n  — e.g., 4 disks: 75% usable
Problem: write hole (parity inconsistency after power failure during write)
```

**RAID 6 — Double parity (survives 2 failures):**
```
Like RAID 5 but two parity blocks per stripe (P and Q)
Survive 2 simultaneous disk failures
Space efficiency: (n-2)/n
```

**RAID 10 — Mirror + Stripe:**
```
4 disks: 2 mirrored pairs, then stripe across them
Survive 1 failure per mirror pair (up to 2 total with luck)
Speed: fast reads and writes
Space: 50% usable
Best for: high-performance databases
```

```bash
# Linux software RAID with mdadm:
mdadm --create /dev/md0 --level=1 --raid-devices=2 /dev/sdb /dev/sdc

# Check RAID status:
cat /proc/mdstat
# Personalities : [raid1]
# md0 : active raid1 sdc[1] sdb[0]
#       976760832 blocks super 1.2 [2/2] [UU]

# Monitor rebuild after replacing a failed disk:
watch -n 1 cat /proc/mdstat
```

---

## 8. Disk Caching and Write-Back

**Page cache (read cache):**
Disk reads are cached in the page cache. Re-reading the same file is served from RAM.

**Write-back (write cache):**
```
Application: write(fd, data, 4096)
  → data goes to page cache (dirty page) — RETURNS IMMEDIATELY
  → application continues, thinks write is done

Later (a few seconds): kernel writes dirty page to disk
  → pdflush / flusher threads
  → triggered by: dirty ratio exceeded, sync() call, umount, or timeout

Risk: system crash between write() return and disk flush → data lost!
```

**Write-through:**
```
write(fd, ...) → kernel writes to disk BEFORE returning to application
Slower but safer
Used when: data must survive crash
Application uses: fsync() to force write-through for specific fd
```

**fsync() and fdatasync():**
```c
// Wait until all pending writes for fd are on disk:
fsync(fd);       // flushes data + metadata (timestamps, file size)
fdatasync(fd);   // flushes data only (faster — skips timestamp update)

// Databases call fsync() after every committed transaction
// Without fsync(): committed transaction can be lost on power failure
```

**O_SYNC flag:**
```c
int fd = open("/var/log/app.log", O_WRONLY | O_APPEND | O_SYNC);
// Every write() call on this fd blocks until data is on disk
```

**Dirty ratio tuning:**
```bash
# Maximum percentage of RAM that can be dirty before sync is forced:
cat /proc/sys/vm/dirty_ratio         # default 20 (%)
echo 10 > /proc/sys/vm/dirty_ratio   # reduce for less data loss on crash

# Background writeback starts at:
cat /proc/sys/vm/dirty_background_ratio  # default 10 (%)

# Maximum time a dirty page can sit before being written:
cat /proc/sys/vm/dirty_expire_centisecs  # default 3000 (30 seconds)
```

---

## Summary

| Concept | Description |
|---------|------------|
| Seek time | Time for HDD head to reach correct track (~5-12ms average) |
| Rotational latency | Wait for sector to rotate under head (~4ms average) |
| IOPS | I/O Operations Per Second: ~80 (HDD), ~90,000 (SATA SSD), ~500,000+ (NVMe) |
| FCFS | First Come First Served; no reordering; poor for random workloads |
| SSTF | Shortest Seek Time First; greedy; risk of starvation |
| SCAN | Elevator: sweep back and forth; no starvation; good average latency |
| C-SCAN | Circular SCAN: one-way sweep; more uniform wait time |
| NVMe | PCIe-attached flash; ~0.02ms latency; 64K+ queues; no seek scheduling needed |
| FTL | Flash Translation Layer: hides flash erase/write limitations from OS |
| TRIM | OS tells SSD which blocks are deleted so SSD can erase proactively |
| BFQ | Linux scheduler for HDDs; fair bandwidth + low latency for interactive I/O |
| RAID 0 | Striping; 2× speed; zero redundancy |
| RAID 1 | Mirroring; survive 1 failure; 50% capacity |
| RAID 5 | Striping + single parity; survive 1 failure; (n-1)/n capacity |
| RAID 10 | Mirror + stripe; high performance + redundancy |
| Page cache | OS caches disk reads/writes in RAM for performance |
| Write-back | Writes cached in RAM, flushed to disk later (fast but loses on crash) |
| fsync() | Flush a file's cached writes to disk before returning |
