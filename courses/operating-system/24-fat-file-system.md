# Chapter 24: FAT — The File Allocation Table

> **"FAT was designed in 1977 for floppy disks with 160KB of storage. It has no permissions, no journaling, and limited filename support. Yet FAT is still used in billions of USB drives, SD cards, and cameras today. Its simplicity is its strength: every device in the world can read it."**

---

## Table of Contents

1. [What Is FAT?](#1-what-is-fat)
2. [FAT Variants — FAT12, FAT16, FAT32](#2-fat-variants--fat12-fat16-fat32)
3. [On-Disk Layout](#3-on-disk-layout)
4. [The File Allocation Table Explained](#4-the-file-allocation-table-explained)
5. [Directory Entries](#5-directory-entries)
6. [How FAT Reads a File](#6-how-fat-reads-a-file)
7. [Long File Names (VFAT)](#7-long-file-names-vfat)
8. [FAT32 vs exFAT](#8-fat32-vs-exfat)
9. [Limitations of FAT](#9-limitations-of-fat)
10. [Why FAT Still Exists](#10-why-fat-still-exists)
11. [Summary](#summary)

---

## 1. What Is FAT?

**FAT (File Allocation Table)** is a file system designed by Bill Gates and Marc McDonald in 1977 for Microsoft's early disk operating systems.

**The core idea:**
Instead of storing block pointers in a separate structure (like inodes), FAT stores a linked list of cluster numbers in a single table at the start of the disk. The "File Allocation Table" IS the map from file data to disk locations.

**FAT is:**
- Simple — no inodes, no bitmaps, minimal metadata
- Universally compatible — supported by every OS, phone, camera, car
- Fast for sequential reads — but slow for fragmented files
- Not crash-safe — no journaling, no transactions

---

## 2. FAT Variants — FAT12, FAT16, FAT32

The number after "FAT" refers to the **bit-width of each entry** in the File Allocation Table:

| Variant | Entry size | Max clusters | Max volume size | Typical use |
|---------|-----------|--------------|-----------------|-------------|
| FAT12 | 12 bits | 4,084 | ~32MB | Floppy disks |
| FAT16 | 16 bits | 65,536 | ~2GB (4GB with 64KB clusters) | Old HDDs |
| FAT32 | 28 bits (!) | 268,435,456 | ~2TB (8TB with 32KB clusters) | USB drives, SD cards |
| exFAT | 32 bits | ~1 billion | ~128 PB | Large flash drives, SDXC |

Note: FAT32 uses only 28 of the 32 bits — the upper 4 bits are reserved.

**How volume size is calculated:**
Max size = max clusters × cluster size
FAT32: 268M clusters × 4KB clusters = ~1TB practical max

---

## 3. On-Disk Layout

A FAT-formatted volume is divided into these regions:

```
┌──────────────────────────────────────────────────────────────────┐
│ Boot Sector │  Reserved  │ FAT 1  │ FAT 2  │ Root Dir │ Data    │
│ (sector 0)  │  Sectors   │ (copy) │ (copy) │ (FAT16)  │ Clusters│
└──────────────────────────────────────────────────────────────────┘
  512 bytes    BPB defines   FAT1=FAT2 (redundancy)
```

**Boot Sector / BPB (BIOS Parameter Block):**
Sector 0 contains the volume parameters:
```
Offset  Size  Description
0       3     Jump instruction to boot code
3       8     OEM Name ("MSDOS5.0", "mkdosfs ", etc.)
11      2     Bytes per sector (almost always 512)
13      1     Sectors per cluster (1,2,4,8,...128)
14      2     Number of reserved sectors (includes boot sector; FAT32: usually 32)
16      1     Number of FATs (always 2 for redundancy)
17      2     Max root directory entries (FAT16: 512, FAT32: 0)
19      2     Total sectors (FAT16; if > 65535, use offset 32)
21      1     Media type (0xF8 = fixed disk, 0xF0 = removable)
22      2     Sectors per FAT (FAT16)
24      2     Sectors per track (for CHS geometry)
26      2     Number of heads
28      4     Hidden sectors
32      4     Total sectors 32-bit (FAT32)
36      4     FAT size in sectors (FAT32 only)
44      4     Root directory start cluster (FAT32 only, usually 2)
```

**Reserved sectors:** Space after boot sector before the first FAT. FAT32 uses this for the FS Info sector (tracks free cluster count).

**Two copies of FAT:** FAT1 and FAT2 are identical. If FAT1 gets corrupted, FAT2 is the backup.

**Root directory:** FAT16 stores root directory as fixed-size array after FAT2. FAT32 stores root directory as a regular cluster chain (starting at cluster 2).

**Data area:** Starts at cluster 2 (by convention). Clusters are numbered starting at 2.

---

## 4. The File Allocation Table Explained

The FAT is an array. **Each element N stores the "next cluster" for cluster N**.

```
FAT array (FAT32, each entry is 28-bit):

Index: 0    1    2    3    4    5    6    7    8    9   10
Value: res  res  3    4    5    EOC  7    8    EOC  FREE FREE

Special values:
  0x0000000  = free cluster
  0x0FFFFFF8 to 0x0FFFFFFF = End Of Chain (EOC) — last cluster of a file
  0x0FFFFFF7 = bad cluster (don't use)
  0x0000001  = reserved (cluster 1)
```

**Reading cluster chain for a file starting at cluster 2:**
```
file starts at cluster 2
FAT[2] = 3  → next cluster is 3
FAT[3] = 4  → next cluster is 4
FAT[4] = 5  → next cluster is 5
FAT[5] = EOC → this is the last cluster
```
File spans clusters: 2 → 3 → 4 → 5 (4 clusters × cluster_size = file content)

**Free space:** Clusters with FAT entry = 0 are free. Finding free space = scanning the FAT for zeros.

**Fragmentation:**
```
Initial state (file A in clusters 2,3,4,5; file B in clusters 6,7):
FAT: [2→3, 3→4, 4→5, 5→EOC, 6→7, 7→EOC]

Delete file A:
FAT: [2→0, 3→0, 4→0, 5→0, 6→7, 7→EOC]

Write new file C (3 clusters):
OS finds free: 2, 3, 4 (first fit)
FAT: [2→3, 3→4, 4→EOC, 5→0, 6→7, 7→EOC]

Write another file D (2 clusters):
OS finds free: 5, then... must search further
→ fragmented file
```

**Fragmentation** means related clusters are scattered, causing extra seek time.

---

## 5. Directory Entries

Each FAT directory entry is exactly **32 bytes**:

```
Offset  Size  Description
0       8     Short name (padded with spaces)
8       3     Extension (padded with spaces) — SFN is "NAME    EXT"
11      1     Attributes byte
               0x01 = Read-Only
               0x02 = Hidden
               0x04 = System
               0x08 = Volume Label
               0x10 = Directory
               0x20 = Archive
               0x0F = Long File Name entry (special)
12      1     Windows NT reserved
13      1     Creation time milliseconds (0-199)
14      2     Creation time (5 bits hour, 6 bits minute, 5 bits 2-second)
16      2     Creation date (7 bits year since 1980, 4 bits month, 5 bits day)
18      2     Last access date
20      2     High word of first cluster (FAT32 only!)
22      2     Modified time
24      2     Modified date
26      2     Low word of first cluster
28      4     File size in bytes
```

**Short File Name (SFN) example:**
```
"readme.txt" → stored as "README  TXT" (8+3 format)
"my file.doc" → CANNOT be stored directly in 8.3!
```

---

## 6. How FAT Reads a File

**Step-by-step: read `/readme.txt`**

1. **Find root directory** (FAT32: start at cluster from BPB; FAT16: fixed location)
2. **Search directory** for entry with name "README  TXT"
3. **Get starting cluster** from directory entry (offset 20+26 = cluster high + cluster low)
4. **Follow FAT chain** to find all clusters
5. **Read cluster data** from data area

```c
// Cluster number to disk sector:
uint32_t cluster_to_sector(uint32_t cluster, BPB *bpb) {
    return bpb->data_start_sector + (cluster - 2) * bpb->sectors_per_cluster;
}
// Note: cluster 2 is the FIRST data cluster (0 and 1 are reserved)
```

**Reading is sequential through the FAT chain:**
```
cluster 2 → read data → get FAT[2]=3 → seek to cluster 3 → read data → ...
```

On a fragmented file, this requires many seeks (slow on HDDs).

---

## 7. Long File Names (VFAT)

Microsoft added Long File Name (LFN) support in Windows 95 via **VFAT** — without breaking compatibility.

**Trick:** Use multiple consecutive 32-byte directory entries with attribute = 0x0F (LFN marker) before the regular 8.3 entry:

```
LFN entries (each holds 13 UCS-2 characters = 26 bytes of filename):
  LFN entry 3 (last): "ng name.txt\0" (end of name + null)
  LFN entry 2:        "is a very lo"
  LFN entry 1:        "this is a ve"
Regular 8.3 entry:    "THISIS~1TXT" (truncated for old systems)
```

**Checksum:** Each LFN entry includes a checksum of the 8.3 name to verify it belongs to that entry.

**Compatibility:**
- Old DOS: ignores LFN entries (attribute 0x0F was never used before), uses 8.3 name
- New Windows: reads LFN entries, presents full name

FAT32 supports filenames up to **255 UTF-16 characters**.

---

## 8. FAT32 vs exFAT

**exFAT** (Extended FAT) was introduced by Microsoft in 2006 for large flash storage (SDXC cards, large USB drives):

| Feature | FAT32 | exFAT |
|---------|-------|-------|
| Max file size | 4GB - 1 byte | 16 EB (exabytes) |
| Max volume | ~2TB (with 32KB clusters ~8TB) | 128 PB |
| Cluster size | 512B – 32KB | 512B – 32MB |
| Directory entries | 512 entries max per dir (effectively unlimited) | Unlimited |
| Timestamps | Seconds resolution | 10ms resolution |
| Permissions | None | None |
| Journaling | No | No |
| OS support | Universal | Most modern OSes |

**exFAT is the default format for SDXC cards** (≥32GB) and used for external drives shared between Windows and macOS.

---

## 9. Limitations of FAT

FAT's simplicity comes at a cost:

| Limitation | Details |
|-----------|---------|
| No permissions | Any user can read/write any file |
| No ownership | No uid/gid — not suitable for multi-user systems |
| No journaling | Power failure can corrupt FAT entries → lost files |
| No hard links | Each file has exactly one directory entry |
| Max file size | 4GB for FAT32 (32-bit size field) |
| No sparse files | Can't create 1TB file with mostly zeros efficiently |
| Fragmentation | No built-in defragmentation; degrades over time |
| No encryption | No native encryption support |
| Poor performance | Following FAT chain is sequential; bad for large files |
| No compression | No native file compression |
| Timestamps | Year 1980–2107 only (FAT32) |

**FAT is unsuitable as a system partition for any serious OS** because it lacks permissions and journaling.

---

## 10. Why FAT Still Exists

Despite its limitations, FAT survives because:

1. **Universal compatibility:** Every OS, phone, camera, car stereo, router reads FAT32.
2. **No driver needed:** The spec is simple enough to implement in 500 lines of C.
3. **Embedded systems:** Microcontrollers (STM32, Arduino with SD shield) use FAT because it's tiny.
4. **EFI System Partition:** UEFI requires FAT32 for the boot partition (all modern computers have FAT32 on the EFI partition).
5. **Interchange format:** "Format as FAT32" is the universal "make this readable everywhere" command.
6. **SD card standard:** The SD Association mandates FAT for SD cards up to 32GB, exFAT for SDXC.

---

## Summary

| Concept | Description |
|---------|------------|
| FAT | File system using linked list of cluster numbers in a table |
| Cluster | Allocation unit (1–128 sectors); fixed size per volume |
| FAT entry | 12/16/28-bit value: 0=free, number=next cluster, EOC=last |
| Directory entry | 32-byte record: name, attributes, start cluster, file size |
| 8.3 format | Short file name: 8 chars + 3 char extension |
| VFAT/LFN | Multiple 32-byte entries per file for long names |
| FAT12 | 12-bit entries; floppy disks |
| FAT16 | 16-bit entries; small HDDs |
| FAT32 | 28-bit entries; USB drives, SD cards |
| exFAT | 32-bit entries; large flash storage; no 4GB file limit |
| BPB | BIOS Parameter Block: volume parameters in boot sector |
| Two FAT copies | Redundancy: FAT1 = FAT2 for corruption recovery |
