# Chapter 27: Modern File Systems — ZFS, APFS, Btrfs

> **"The classical file system design — inode tables, journaling, block bitmaps — was built when disks were mechanical and small. Modern file systems rethink storage from scratch: copy-on-write semantics, built-in checksums, snapshots, RAID, and deduplication. ZFS, APFS, and Btrfs represent where file systems are going."**

---

## Table of Contents

1. [Why Modern File Systems?](#1-why-modern-file-systems)
2. [Copy-On-Write (COW) — The Core Idea](#2-copy-on-write-cow--the-core-idea)
3. [ZFS — The "Last Word" in File Systems](#3-zfs--the-last-word-in-file-systems)
4. [Btrfs — Linux's Modern File System](#4-btrfs--linuxs-modern-file-system)
5. [APFS — Apple's File System](#5-apfs--apples-file-system)
6. [Snapshots — Time Travel for Files](#6-snapshots--time-travel-for-files)
7. [Data Integrity — Checksums and Scrubbing](#7-data-integrity--checksums-and-scrubbing)
8. [Comparing Modern File Systems](#8-comparing-modern-file-systems)
9. [Summary](#summary)

---

## 1. Why Modern File Systems?

**Problems with ext4/NTFS for modern workloads:**

| Problem | Classic FS | Modern FS |
|---------|-----------|-----------|
| Silent corruption | No checksums — bit flip = corrupted file | End-to-end checksums on all data |
| Snapshot cost | Expensive: copy entire volume | COW: free — just save the tree root |
| RAID complexity | Separate software/hardware RAID | Integrated into the file system |
| Deduplication | Not possible | Built-in block-level dedup |
| Large directories | Slow (B-tree helps, but still) | O(1) with hash trees |
| File system repair | fsck scans entire disk (hours) | COW trees always consistent |
| Thin provisioning | Requires LVM complexity | Built-in |

**The real driver: silent corruption.**

Studies by Microsoft, Google, and CERN found that:
- Cosmic rays cause bit flips in DRAM (detected by ECC RAM, but not all servers have ECC)
- Controller firmware bugs silently corrupt writes
- Torn writes (partially written blocks) aren't detected
- A file system without checksums will silently serve corrupted data

ext4 and NTFS trust the hardware to be correct. ZFS and Btrfs do not.

---

## 2. Copy-On-Write (COW) — The Core Idea

**Traditional update (in-place):**
```
Update block 500:
  1. Read block 500
  2. Modify it
  3. Write back to block 500 (SAME location)
  
If power fails during step 3: block 500 is half-written → corrupted!
→ Needs journaling to recover
```

**Copy-on-write (COW) update:**
```
Update block 500:
  1. Read block 500
  2. Write modified version to a NEW free block (block 9000)
  3. Update parent tree node to point to 9000 instead of 500
  4. Mark block 500 as free
  
If power fails during step 2: old block 500 still intact, new block 9000 abandoned
If power fails during step 3: same as above (old root still valid)
If power fails during step 4: both blocks 500 and 9000 exist; leak detected by GC
→ NO journal needed! Old tree is always valid.
```

**The tree is always consistent:**
A COW file system stores data in a tree. When any data changes, a new path from that leaf to the root is written. The old root is valid until the new root is atomically installed (single pointer update).

```
Before update:                After update:
    Root_A                        Root_B (new)
   /    \                        /    \
 Dir1   Dir2                  Dir1   Dir2' (new copy)
  |      |                     |      |
 File1  File2                 File1  File2 (new data)
                                      |
                                  [new block]
                                  
Old root_A and all old blocks: still valid until freed
```

---

## 3. ZFS — The "Last Word" in File Systems

**ZFS** was developed by Sun Microsystems (2001-2005), released open source with OpenSolaris. Available on Linux via OpenZFS project.

**ZFS philosophy: "The last word in file systems."**
ZFS refuses to compromise on data integrity. If it can't verify the data is correct, it won't serve it.

### Storage Pools (zpools)

ZFS replaces the traditional "partition → format → mount" workflow with **storage pools**.

```bash
# Create a pool from 3 disks (mirrored for redundancy):
zpool create mypool mirror /dev/sdb /dev/sdc /dev/sdd

# Create a pool with RAIDZ2 (similar to RAID6: 2-disk fault tolerance):
zpool create datapool raidz2 /dev/sdb /dev/sdc /dev/sdd /dev/sde /dev/sdf

# ZFS handles all RAID logic INTERNALLY — no mdadm or hardware RAID needed
```

**RAIDZ levels:**
```
RAIDZ1: 1 parity disk (tolerate 1 disk failure) — like RAID5
RAIDZ2: 2 parity disks (tolerate 2 disk failures) — like RAID6
RAIDZ3: 3 parity disks (tolerate 3 disk failures)

ZFS RAIDZ avoids the "write hole" problem of hardware RAID5/6:
Traditional RAID5: power failure during stripe write corrupts parity
ZFS RAIDZ: COW means new stripe written atomically — no write hole
```

### Datasets (Filesystems)

Inside a pool, you create **datasets** — independent file systems:
```bash
# Create datasets:
zfs create mypool/home
zfs create mypool/vm-images
zfs create mypool/logs

# Set compression:
zfs set compression=lz4 mypool/home
zfs set compression=zstd mypool/vm-images

# Set quotas:
zfs set quota=100G mypool/home

# View datasets:
zfs list
# NAME                USED  AVAIL     REFER  MOUNTPOINT
# mypool              142G   358G     1.34M  /mypool
# mypool/home         98G    358G     98G    /mypool/home
```

### The Merkle-Like Tree Structure

ZFS stores everything in a **tree of blocks with checksums**:

```
Uberblock (32 entries in a ring, located in the first 256KB of each vdev):
  txg: 12345  ← transaction group number
  root block pointer:
    dva[0]: vdev 0, offset 100GB    ← two copies of root object set
    dva[1]: vdev 1, offset 100GB
    checksum: SHA-256 or Fletcher4

    Root Object Set
    ├── Meta-Object Set (MOS): pool-level metadata
    │     ├── dnode for $DIRECTORY_OBJECT
    │     ├── dnode for Space Map Object
    │     └── dnode for Filesystem Label
    └── Object Set (per dataset)
          ├── dnode for each file (like inode)
          └── data blocks
```

**Every block has a 256-bit checksum stored in the PARENT block pointer**, not in the block itself (so a corrupted block doesn't also corrupt its own checksum).

### Checksums and Self-Healing

```bash
# ZFS verifies checksums on every read
# If a checksum fails on a mirrored pool:
#   1. Report an error
#   2. Automatically read the correct copy from the mirror
#   3. Repair the bad copy ("heal")
#   4. Serve the correct data to the application

# Run a full scrub (verify all data against checksums):
zpool scrub mypool

# View scrub results:
zpool status mypool
# scan: scrub repaired 0B in 00:15:22 with 0 errors on ...
```

### ARC — Adaptive Replacement Cache

ZFS has its own memory cache (ARC — Adaptive Replacement Cache) that bypasses the OS page cache:
```
ARC uses available RAM to cache frequently accessed data
Default: grows up to half of physical RAM
L2ARC: can extend ARC to an SSD (NVMe L2ARC for faster cold cache reads)
```

---

## 4. Btrfs — Linux's Modern File System

**Btrfs** (B-tree File System, pronounced "butter-FS") is the Linux answer to ZFS. Development began at Oracle in 2007, merged into the Linux kernel in 2009.

**Core features: COW, snapshots, built-in RAID, checksums.**

```bash
# Create Btrfs on a single drive:
mkfs.btrfs /dev/sdb

# Create Btrfs spanning multiple drives (RAID1):
mkfs.btrfs -d raid1 -m raid1 /dev/sdb /dev/sdc
# -d: data redundancy level
# -m: metadata redundancy level

# Mount:
mount -t btrfs /dev/sdb /mnt/data

# Create subvolume (like ZFS dataset):
btrfs subvolume create /mnt/data/home
btrfs subvolume create /mnt/data/var

# Snapshot a subvolume (instant, COW!):
btrfs subvolume snapshot /mnt/data/home /mnt/data/home-backup-$(date +%Y%m%d)
# Takes 0 bytes of extra space initially — shared blocks via COW

# Set compression:
btrfs property set /mnt/data/home compression zstd

# Check disk usage (COW makes `du` misleading):
btrfs filesystem df /mnt/data
btrfs filesystem usage /mnt/data
```

**Btrfs on-disk structures:**
```
Btrfs uses B-trees (hence "B-tree File System") for everything:

Filesystem Tree: maps (subvolume_id, objectid, offset) → data
Chunk Tree: maps logical addresses to physical locations (like LVM)
Device Tree: maps physical devices to chunks
Root Tree: maps subvolume IDs to root block pointers
Extent Tree: tracks which blocks are allocated
Checksum Tree: checksums for all data blocks
```

**Btrfs RAID:**
```
RAID0:  stripe, no redundancy
RAID1:  mirror (copies to 2 devices)
RAID10: stripe + mirror
RAID5/6: parity (EXPERIMENTAL — known issues, avoid for production)
```

**Send/receive (efficient backup):**
```bash
# Take snapshot of current state:
btrfs subvolume snapshot -r /mnt/data/home /mnt/data/home-snap1

# Send snapshot to another disk:
btrfs send /mnt/data/home-snap1 | btrfs receive /mnt/backup/

# Next day, send only the DIFFERENCE:
btrfs subvolume snapshot -r /mnt/data/home /mnt/data/home-snap2
btrfs send -p /mnt/data/home-snap1 /mnt/data/home-snap2 | btrfs receive /mnt/backup/
# Only sends changed blocks — extremely efficient!
```

**Btrfs stability note:**
As of 2026, Btrfs is stable for most workloads. RAID5/6 is still marked experimental. Facebook/Meta uses Btrfs at massive scale with their own patches. Fedora Linux made Btrfs the default in 2020.

---

## 5. APFS — Apple's File System

**APFS** (Apple File System) was announced at WWDC 2016, deployed in iOS 10.3 (2017) and macOS High Sierra 10.13 (2017). It replaced HFS+, which had been in use since 1998.

**Why Apple needed a new FS:**
- HFS+ was designed for HDDs; flash storage has different characteristics
- HFS+ had a single lock for the entire volume (bad for SSD parallelism)
- HFS+ timestamps had 1-second resolution (inadequate)
- HFS+ had no snapshots (Time Machine snapshots were awkward workarounds)
- HFS+ had no space sharing between volumes on the same partition

**APFS features:**

**Space sharing containers:**
```
APFS Container (one partition):
├── Volume: System    (used: 15GB, can grow as needed)
├── Volume: Data      (used: 100GB, can grow as needed)
├── Volume: Preboot   (small — boot files)
├── Volume: Recovery  (small — recovery OS)
└── Volume: VM        (swap, can grow)

All volumes share the same pool of free space.
No need to pre-partition space between System and Data!
```

**Copy-on-write:**
```
APFS is fully COW:
  - File writes go to new blocks; old blocks freed after all snapshots drop them
  - Crash safe without a traditional journal (COW provides consistency)
  - Apple calls their implementation "checkpoint-based consistency"
```

**Snapshots:**
```bash
# Create snapshot (macOS — Time Machine uses this):
tmutil snapshot

# List snapshots:
tmutil listlocalsnapshots /

# Mount snapshot read-only:
mount_apfs -o ro -s com.apple.TimeMachine.2026-07-01-120000 /dev/disk1s1 /mnt/snap
```

**Native encryption:**
APFS was designed with encryption as a first-class feature:
- Per-volume encryption
- Per-file encryption keys
- Hardware-backed key storage (Secure Enclave on Apple Silicon)
- FileVault on macOS uses APFS native encryption

**Clones:**
```bash
# Instant copy of a file (COW — shares blocks until modified):
cp --reflink /path/source /path/dest  # on Linux Btrfs
# macOS: any file copy within APFS is automatically a clone
# 1GB file copied "instantly" — uses zero extra space initially
```

**Fast directory sizing:**
APFS maintains a running count of directory sizes. `du` is instant (no tree walk needed).

**Space efficiency:**
APFS supports **compression** (transparent, used by macOS for system files) and **sparse files**.

---

## 6. Snapshots — Time Travel for Files

Snapshots are one of the most powerful features of COW file systems.

**How a COW snapshot works:**
```
Before snapshot:                    After taking snapshot:
                                    
Tree root → A → [data1]             Snapshot root = same tree root
            B → [data2]             
                                    Now modify file B:
                                    
                                    Tree root (new) → A → [data1]
                                                  → B' → [data3] (new data)
                                    
                                    Snapshot root (old) → A → [data1]
                                                       → B → [data2] (preserved!)
                                    
                                    Space used: data1 + data2 + data3
                                    (NOT data1 × 2 + data2 + data3)
                                    
                                    B's old block [data2] kept because snapshot references it.
                                    Only freed when snapshot is deleted.
```

**Use cases:**
```bash
# ZFS snapshot workflow:
zfs snapshot mypool/home@before-upgrade
# do OS upgrade...
# if upgrade broke things:
zfs rollback mypool/home@before-upgrade  # instant rollback!
# if upgrade succeeded:
zfs destroy mypool/home@before-upgrade   # free the space

# Btrfs snapshot for VM images:
btrfs subvolume snapshot /vms/ubuntu-base /vms/ubuntu-dev
# Launch VM from /vms/ubuntu-dev — starts from a clone of base, all changes stay separate
```

---

## 7. Data Integrity — Checksums and Scrubbing

**End-to-end checksums:**
```
Traditional: OS trusts disk controller to return correct data
ZFS/Btrfs: compute checksum of every block on write, verify on read

If checksum fails:
  Single copy: error reported to application (with details!)
  Mirrored: read from other copy, heal the bad copy
  RAIDZ: reconstruct from parity, heal the bad disk
```

**Checksums used:**
```
ZFS:   Fletcher2, Fletcher4 (fast), SHA-256 (secure, slower), SHA-512, Skein
Btrfs: CRC32C (fast, hardware-accelerated on modern CPUs), XXHASH, SHA-256, BLAKE2B
```

**Scrubbing:**
```bash
# ZFS — scrub all data (verify checksums):
zpool scrub mypool
# Reads every block, verifies checksum, repairs errors using redundancy
# Schedule monthly: zpool scrub is the file system equivalent of "fsck everything"

# Btrfs:
btrfs scrub start /mnt/data
btrfs scrub status /mnt/data
```

---

## 8. Comparing Modern File Systems

| Feature | ZFS | Btrfs | APFS |
|---------|-----|-------|------|
| Platform | Linux, FreeBSD, macOS (limited) | Linux | macOS, iOS |
| COW | Yes | Yes | Yes |
| Snapshots | Yes, excellent | Yes, good | Yes, good |
| RAID | RAIDZ1/2/3 + mirror | RAID0/1/10, RAID5/6 experimental | No (use mdadm) |
| Compression | LZ4, gzip, zstd | LZO, zlib, zstd | zlib (system files) |
| Deduplication | Yes (expensive, needs RAM) | Inline (experimental) | No |
| Encryption | Yes (native) | No (use dm-crypt) | Yes (native, HW-backed) |
| Checksums | Yes (Fletcher, SHA-256) | Yes (CRC32C, SHA-256) | Yes |
| Send/receive | zfs send/recv | btrfs send/recv | No |
| Space sharing | zpools | Pool with subvolumes | APFS containers |
| Stability | Very stable | Mostly stable, RAID5/6 avoid | Very stable |
| License | CDDL (not GPL — kernel module issues) | GPL (built into kernel) | Proprietary |

---

## Summary

| Concept | Description |
|---------|------------|
| COW (Copy-On-Write) | Never modify data in place; write to new blocks, update pointers atomically |
| Snapshot | Saved tree root; old blocks kept as long as snapshot exists |
| Checksum | Hash of block stored in parent pointer; detects corruption on read |
| Scrub | Periodic scan of all data; verify checksums, repair using redundancy |
| ZFS zpool | Storage pool spanning multiple devices; handles RAID internally |
| RAIDZ | ZFS's RAID implementation; avoids write hole with COW |
| Btrfs subvolume | Independently mountable section within a Btrfs pool |
| APFS container | Partition containing multiple volumes that share free space |
| APFS volume | Logical volume within a container; System/Data split on modern macOS |
| Deduplication | Store only one copy of identical blocks; others point to it |
| Clone/reflink | Instant file copy sharing blocks until one is modified (COW clone) |
| Send/receive | Incremental backup by sending only changed blocks between snapshots |
