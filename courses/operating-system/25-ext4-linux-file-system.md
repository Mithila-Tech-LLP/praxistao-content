# Chapter 25: ext4 — The Linux File System

> **"ext4 is the evolution of 30 years of Unix file system thinking. It combines the classical inode-based design with extents for large files, journaling for crash recovery, delayed allocation for performance, and features borrowed from the best file systems of the 1990s. It's not the newest or fanciest, but it's the rock-solid default that billions of Linux systems trust."**

---

## Table of Contents

1. [The ext Family — A Brief History](#1-the-ext-family--a-brief-history)
2. [ext4 Features Overview](#2-ext4-features-overview)
3. [On-Disk Layout](#3-on-disk-layout)
4. [Inodes in ext4](#4-inodes-in-ext4)
5. [Extents — Efficient Large File Storage](#5-extents--efficient-large-file-storage)
6. [The Journal — Crash Recovery](#6-the-journal--crash-recovery)
7. [Directory Structure in ext4](#7-directory-structure-in-ext4)
8. [Delayed Allocation — Performance Optimization](#8-delayed-allocation--performance-optimization)
9. [Block Groups and Locality](#9-block-groups-and-locality)
10. [Creating and Inspecting ext4](#10-creating-and-inspecting-ext4)
11. [Summary](#summary)

---

## 1. The ext Family — A Brief History

```
1992: ext (Extended File System)
      First real Linux FS; replaced Minix FS
      No timestamps, 2GB max, no journaling

1993: ext2
      Inodes with proper metadata, 2TB max, 255-char filenames
      Still used today for boot partitions (no journaling = fewer writes to flash)

2001: ext3
      Added journaling to ext2 (could upgrade in-place!)
      Same on-disk format + journal area
      
2008: ext4
      Extents (replace indirect blocks), 1EB max file size
      Nanosecond timestamps, 64-bit block numbers
      Delayed allocation, faster fsck, persistent preallocation
      Backward compatible: ext3 can be mounted as ext4 with limited features
```

**Design philosophy:** ext4 is a production-ready, stable, backward-compatible upgrade. It is NOT experimental. The goal was to handle modern workloads (large files, millions of files, SSDs) while remaining reliable.

---

## 2. ext4 Features Overview

| Feature | Description |
|---------|------------|
| Max file size | 16TB (with 4KB blocks) |
| Max volume size | 1 EB (exabyte) |
| Max files | 4 billion (32-bit inode numbers) |
| Block size | 1KB, 2KB, 4KB (default), 64KB |
| Filename length | 255 bytes (UTF-8) |
| Journaling | writeback, ordered (default), data modes |
| Extents | Contiguous block ranges (replaces indirect blocks) |
| Delayed allocation | Batch block allocations for better locality |
| Multiblock allocation | Allocate multiple blocks at once |
| dir_index | HTree (hash tree) for large directories |
| 64-bit block numbers | Volumes > 8TB (requires ext4 mode) |
| Nanosecond timestamps | Creation time, atime, mtime, ctime all at ns resolution |
| Persistent preallocation | Reserve disk space without writing data (fallocate) |
| Online defragmentation | e4defrag tool (limited) |

---

## 3. On-Disk Layout

```
ext4 partition:
┌────────────────────────────────────────────────────────────────────┐
│  Boot  │ Block Group 0     │ Block Group 1 │ ... │ Block Group N  │
│  Block │ (may have copies  │               │     │                │
│(1KB)   │  of superblock)   │               │     │                │
└────────────────────────────────────────────────────────────────────┘

Each Block Group (e.g., 32768 blocks × 4KB = 128MB per group):
┌──────────┬───────┬────────┬────────┬───────────┬──────────────────┐
│Superblock│ Group │ Block  │ Inode  │ Inode     │ Data Blocks      │
│ (copy)   │ Desc  │ Bitmap │ Bitmap │ Table     │                  │
│ 1 block  │ Table │ 1 blk  │ 1 blk  │ N blocks  │ (most of group)  │
└──────────┴───────┴────────┴────────┴───────────┴──────────────────┘
 (not in all groups — sparse superblock feature)
```

**Superblock (block group 0, offset 1024 bytes):**
```c
struct ext4_super_block {
    uint32_t s_inodes_count;       // total inodes
    uint32_t s_blocks_count_lo;    // total blocks (low 32 bits)
    uint32_t s_r_blocks_count_lo;  // reserved blocks for root
    uint32_t s_free_blocks_count_lo;
    uint32_t s_free_inodes_count;
    uint32_t s_first_data_block;   // 0 for 4KB blocks, 1 for 1KB blocks
    uint32_t s_log_block_size;     // block size = 1024 << s_log_block_size
    uint32_t s_blocks_per_group;   // blocks per block group
    uint32_t s_inodes_per_group;   // inodes per block group
    uint16_t s_magic;              // 0xEF53 (must be this value)
    uint16_t s_state;              // 1=clean, 2=errors
    uint32_t s_rev_level;          // 0=original, 1=dynamic
    uint32_t s_feature_compat;     // compatible features
    uint32_t s_feature_incompat;   // incompatible features (must support)
    // ... more fields
};
```

**Block Group Descriptor Table:**
After the superblock, a table of group descriptors:
```c
struct ext4_group_desc {
    uint32_t bg_block_bitmap_lo;    // block number of block bitmap
    uint32_t bg_inode_bitmap_lo;    // block number of inode bitmap
    uint32_t bg_inode_table_lo;     // block number of inode table
    uint16_t bg_free_blocks_count_lo;
    uint16_t bg_free_inodes_count_lo;
    uint16_t bg_used_dirs_count_lo;
    // ...
};
```

---

## 4. Inodes in ext4

Each inode in ext4 is **256 bytes** (up from 128 bytes in ext2/3):

```c
struct ext4_inode {
    uint16_t i_mode;          // file type + permissions (0o100644 = regular rw-r--r--)
    uint16_t i_uid;           // owner user id (lower 16 bits)
    uint32_t i_size_lo;       // file size in bytes (low 32 bits)
    uint32_t i_atime;         // access time (seconds since epoch)
    uint32_t i_ctime;         // change time
    uint32_t i_mtime;         // modify time
    uint32_t i_dtime;         // deletion time (when inode was deleted)
    uint16_t i_gid;           // group id
    uint16_t i_links_count;   // hard link count
    uint32_t i_blocks_lo;     // block count (in 512-byte units)
    uint32_t i_flags;         // flags (EXT4_EXTENTS_FL, EXT4_INLINE_DATA_FL, etc.)
    
    union {
        uint32_t i_block[15]; // old: 12 direct + 1 indirect + 1 double + 1 triple
                               // new: extent tree (if EXT4_EXTENTS_FL flag set)
    } i_data;
    
    uint32_t i_generation;    // file version (NFS use)
    uint32_t i_file_acl_lo;   // extended attributes block
    uint32_t i_size_high;     // high 32 bits of file size (for large files)
    // ...
    uint16_t i_extra_isize;   // size of extra inode fields
    uint32_t i_ctime_extra;   // change time (nanoseconds)
    uint32_t i_mtime_extra;   // modify time (nanoseconds)
    uint32_t i_atime_extra;   // access time (nanoseconds)
    uint32_t i_crtime;        // creation time (not in ext3!)
    uint32_t i_crtime_extra;
};
```

**Inode addressing:**
Given inode number N:
1. Block group = (N - 1) / inodes_per_group
2. Inode index within group = (N - 1) % inodes_per_group
3. Inode offset in inode table = inode_index × inode_size
4. Block containing inode = bg_inode_table + offset / block_size

---

## 5. Extents — Efficient Large File Storage

**The old way (indirect blocks) had problems:**
```
Inode → [12 direct blocks, 1 indirect pointer, 1 double indirect, 1 triple indirect]

A 1GB file needs:
  256KB of indirect blocks
  Requires 3 extra disk reads just to find where data is
  
A 100MB write might touch 25,600 separate indirect block entries
```

**The ext4 way: Extents**

An **extent** is a contiguous range of blocks described by just three numbers:
- Starting logical block (position in file)
- Starting physical block (location on disk)
- Length in blocks

```c
struct ext4_extent {
    uint32_t ee_block;    // first logical block this extent covers
    uint16_t ee_len;      // number of blocks (max 32768 for uninitialized)
    uint16_t ee_start_hi; // high 16 bits of physical block
    uint32_t ee_start_lo; // low 32 bits of physical block
};
```

**Extent tree:**
The 60 bytes in the inode's i_block field (where block pointers used to live) now hold the root of an **extent tree**:

```
Inode's i_block (60 bytes):
┌──────────────────────────────────────────────────────┐
│ ext4_extent_header │ extent/index entries (4 max)    │
└──────────────────────────────────────────────────────┘

If file fits in 4 extents (common case):
  [header, depth=0] [extent1: block 0-8191] [extent2: block 8192-16383] ...
  → only 1 access to find any block!

If file needs more extents:
  [header, depth=1] [index1 → block X] [index2 → block Y]
                              ↓
               [header, depth=0] [extent1] [extent2] ... (up to 340 per block)
```

**Huge win for large files:**
A 1GB file stored contiguously (SSDs favor contiguous allocation):
- Old: hundreds of indirect block entries, multiple extra reads
- New: ONE extent `{start_block, 262144_blocks}`, found in the inode itself

**Fragmented file:**
If a file is fragmented, it uses more extents, and the tree grows deeper.

---

## 6. The Journal — Crash Recovery

ext4 uses a **journal** (separate file `.journal` or `.journal` inode) to ensure consistency.

**Journal modes:**

**writeback mode:**
```
Journal records: only metadata operations (inode updates, directory changes)
Data: written directly to disk (no journal)
After crash: metadata consistent, but file data might be stale
→ You might see an inode claiming a file has X bytes but the data blocks aren't updated
```

**ordered mode (default):**
```
1. Write file DATA to disk first (before updating metadata)
2. Journal metadata: start transaction, record metadata changes, commit
3. Apply metadata changes to disk
After crash:
  If step 1 didn't finish: inode still has old size, old data → consistent (old state)
  If step 2 didn't commit: journal incomplete, ignore → consistent (old state)
  If committed: replay journal → consistent (new state)
→ Never shows a file with new size but old/garbage data
```

**data mode:**
```
Both metadata AND data go through journal
After crash: both metadata and data are fully consistent
Slowest: everything written twice (journal + real location)
```

**Journal format:**
The journal contains **transactions** made of **journal blocks**:
```
Journal Block Types:
  JBD2_DESCRIPTOR_BLOCK: list of blocks that will be modified
  JBD2_COMMIT_BLOCK: marks transaction as complete
  JBD2_REVOKE_BLOCK: marks blocks as no longer valid

Transaction flow:
  BEGIN → write descriptor (list of dirty blocks) → write data copies → COMMIT
  On crash replay: find complete transactions (BEGIN to COMMIT) → re-apply
```

---

## 7. Directory Structure in ext4

ext4 supports two directory storage formats:

**Linear (classic):**
Directory data blocks contain a linked list of variable-length entries:
```c
struct ext4_dir_entry_2 {
    uint32_t inode;       // inode number (0 = empty entry)
    uint16_t rec_len;     // length of this directory entry (to find next)
    uint8_t  name_len;    // length of filename
    uint8_t  file_type;   // 1=regular, 2=directory, 7=symlink, etc.
    char     name[];      // filename (not null-terminated)
};
```

Finding a file in a large directory with linear search = O(n) where n = number of files.

**HTree (Hash Tree) — dir_index feature:**
For directories with many files (hundreds or thousands), ext4 uses an **HTree** — a B-tree of filename hashes.

```
Root block:
  [H-Tree header] [hash range → leaf block] [hash range → leaf block] ...

Leaf block:
  [directory entries for filenames with hashes in this range]

Lookup:
  hash("filename") → find range → read one leaf block → find entry
  → O(log n) instead of O(n)
```

Creating a file in a large directory with HTree is fast even with 100,000 entries.

---

## 8. Delayed Allocation — Performance Optimization

**The problem with eager allocation:**
When an application writes data, the OS might allocate blocks immediately before knowing the final file size. Result: fragmented files (blocks allocated piecemeal as writes happen).

**Delayed allocation (delalloc):**
```
Application writes to a file:
  1. Data goes into page cache (RAM only — no disk allocation yet)
  2. Inode is dirtied (marked as needing write)
  3. NO block allocation yet

When the page is finally flushed to disk:
  1. OS knows exactly how much data needs to be written
  2. Allocates a contiguous run of blocks for the entire dirty data
  3. Writes all data to the contiguous run
```

**Benefits:**
- Fewer, larger extents → less fragmentation
- Multiblock allocator (mballoc) can allocate 100MB of contiguous blocks at once
- Better for SSDs (larger, sequential writes are more efficient)

**Risk:**
If power fails after write() returns but before flush: data is lost! This is why databases use `fsync()` to force data to disk before considering a write committed.

---

## 9. Block Groups and Locality

**Block groups** are a key design principle for performance: **keep a file's inode, directory entry, and data blocks physically close together**.

```
Block group = 128MB of disk space

Rules:
1. Create a file in directory D → allocate inode in SAME group as D's inode
2. Allocate data blocks → try to allocate in SAME group as the inode
3. If group is full → use a nearby group

Effect: reading a file requires disk head to move very little
  → much better performance on HDDs
  → better cache utilization on SSDs
```

**Inode allocation strategy:**
- Directories: spread across block groups (to distribute load)
- Files in same directory: try to put in same block group

---

## 10. Creating and Inspecting ext4

```bash
# Format a partition as ext4:
mkfs.ext4 /dev/sdb1

# Format with specific options:
mkfs.ext4 -L "mydata" -m 1 -E lazy_itable_init=0 /dev/sdb1
#          -L: volume label
#          -m 1: 1% reserved blocks (default 5%)
#          lazy_itable_init=0: initialize all inodes now (slower format, faster first mount)

# Mount:
mount -t ext4 /dev/sdb1 /mnt/data

# View superblock info:
tune2fs -l /dev/sdb1
# Shows: block count, inode count, journal size, mount count, last check, features

# Dump inode info:
debugfs /dev/sdb1
debugfs: stat /etc/passwd   # shows full inode structure
debugfs: dump <2> /tmp/root_dir_raw  # dump raw block 2 (root inode)

# Check filesystem (must be unmounted):
fsck.ext4 -n /dev/sdb1   # -n = dry run (no changes)

# View extent tree of a file:
filefrag -v /path/to/bigfile
# Shows: physical/logical block mapping, number of extents

# Get detailed filesystem stats:
df -i /   # inode usage
dumpe2fs /dev/sda1 | grep -E "(Block count|Inode count|Journal)"
```

---

## Summary

| Concept | Details |
|---------|---------|
| ext family | ext → ext2 → ext3 (add journal) → ext4 (extents, nanosec timestamps) |
| Magic number | 0xEF53 in superblock identifies ext2/3/4 |
| Block group | 128MB region; inode+data kept together for locality |
| Superblock | Volume metadata; copied in some block groups |
| ext4 inode | 256 bytes; includes extent tree root, nanosecond timestamps |
| Extents | Contiguous block range: {logical_start, physical_start, length}; replaces indirect blocks |
| Extent tree | B-tree of extents; root fits in 60 bytes of inode |
| Journal modes | writeback (metadata only) / ordered (data first, default) / data (everything) |
| HTree | Hash tree for large directories; O(log n) lookup |
| Delayed allocation | Defer block allocation until flush; better contiguity and fewer extents |
| Block bitmaps | Per-group; 1 bit per block; fast free-space search |
| Inode bitmaps | Per-group; 1 bit per inode; fast free-inode search |
