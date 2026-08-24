# Chapter 26: NTFS — The Windows File System

> **"NTFS was designed from scratch for Windows NT in the early 1990s, with the explicit goals that failed earlier: crash recovery, security, large files, and Unicode. Where FAT is a table, NTFS is a relational database — every piece of the file system, including its own metadata, is stored as a file."**

---

## Table of Contents

1. [NTFS History and Goals](#1-ntfs-history-and-goals)
2. [The Core Concept: Everything Is a File](#2-the-core-concept-everything-is-a-file)
3. [The Master File Table (MFT)](#3-the-master-file-table-mft)
4. [MFT Record Structure](#4-mft-record-structure)
5. [Attributes — How Data Is Stored](#5-attributes--how-data-is-stored)
6. [NTFS Journaling](#6-ntfs-journaling)
7. [NTFS Permissions (ACLs)](#7-ntfs-permissions-acls)
8. [NTFS Special Features](#8-ntfs-special-features)
9. [On-Disk Layout](#9-on-disk-layout)
10. [NTFS vs ext4 Comparison](#10-ntfs-vs-ext4-comparison)
11. [Summary](#summary)

---

## 1. NTFS History and Goals

**FAT's problems (by 1990):**
- No file permissions → any user can delete anything
- No crash recovery → power failure corrupts everything
- 2GB max file size → inadequate for databases and video
- 8.3 filenames → awkward for users
- ASCII only → can't handle international characters

**NTFS (New Technology File System)** was designed by David Cutler's team for Windows NT 3.1 (1993). Goals:
- Full crash recovery via journaling
- Per-file ACL permissions
- Large files (up to 16 EB theoretically)
- Long filenames with Unicode (UTF-16)
- Reliable, enterprise-grade storage

**Versions:**
```
NTFS 1.0: Windows NT 3.1 (1993)
NTFS 1.2: Windows NT 4.0 (1996) — compression
NTFS 3.0: Windows 2000 — encryption (EFS), disk quotas, sparse files
NTFS 3.1: Windows XP/Vista/7/8/10/11 — current version
```

---

## 2. The Core Concept: Everything Is a File

NTFS's most elegant design choice: **all file system metadata is stored in regular NTFS files**.

```
$MFT          (inode 0):  The Master File Table itself
$MFTMirr      (inode 1):  Partial backup of MFT (first 4 entries)
$LogFile      (inode 2):  Journal for crash recovery
$Volume       (inode 3):  Volume information (name, version, flags)
$AttrDef      (inode 4):  Attribute type definitions
.              (inode 5):  Root directory (the ".")
$Bitmap        (inode 6):  Allocation bitmap for all clusters
$Boot          (inode 7):  Boot sector and boot code
$BadClus       (inode 8):  List of bad clusters
$Secure        (inode 9):  Security descriptors (ACLs)
$UpCase        (inode 10): Uppercase translation table (for case-insensitive compare)
$Extend        (inode 11): Directory containing extended metadata files
  $Quota          Disk quotas
  $ObjId          Object IDs for files
  $Reparse        Reparse points (symlinks, junctions, mount points)
  $UsnJrnl        Update Sequence Number journal (change tracking)
```

You can actually read these files with the right tools! They're just files with special flags.

---

## 3. The Master File Table (MFT)

The **MFT** is the heart of NTFS. Think of it as the combined inode table + directory + data for all file system metadata.

**Structure:**
- An array of **MFT records**, each 1024 bytes (default)
- Every file and directory has an MFT record (like an inode in ext4)
- MFT record = collection of **attributes** describing the file

```
MFT layout:
┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬────────────┐
│  $MFT│MFTMir│LogFil│Volume│AttrDf│ Root │Bitmap│ Boot │ User files │
│ rec 0│ rec 1│ rec 2│ rec 3│ rec 4│ rec 5│ rec 6│ rec 7│ rec 16+    │
└──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┴────────────┘
```

**MFT reservation:**
When NTFS is formatted, it reserves 12.5% of the volume for MFT growth (to avoid fragmentation of the MFT itself). This is the "MFT zone."

**Finding a file:**
1. File's inode number = MFT record number
2. Calculate offset in MFT: offset = record_number × 1024
3. Read that MFT record
4. Parse its attributes to find the data

---

## 4. MFT Record Structure

Each 1024-byte MFT record:

```
MFT Record (1024 bytes):
┌─────────────────────────────────────────┐
│ FILE signature (4 bytes: "FILE")        │
│ Update sequence offset/size             │
│ Log file sequence number ($LogFile LSN) │
│ Sequence number (incremented on reuse)  │
│ Hard link count                         │
│ Offset to first attribute               │
│ Flags: 0x01=in use, 0x02=directory      │
│ Used size, allocated size               │
│ Base record reference (for extensions)  │
│ Next attribute ID                       │
│ Record number (matches MFT index)       │
├─────────────────────────────────────────┤
│ Attribute 1: $STANDARD_INFORMATION     │ ← always first
│ Attribute 2: $FILE_NAME                │ ← always present
│ Attribute 3: $DATA (or $INDEX_ROOT)    │ ← file content or directory
│ Attribute 4: $SECURITY_DESCRIPTOR      │ ← permissions
│ ...more attributes...                  │
│ 0xFFFFFFFF: end marker                 │
└─────────────────────────────────────────┘
```

**Multi-record files:**
If one MFT record can't hold all attributes (many alternate data streams, very long filename), NTFS uses **extension records** linked by the base record.

---

## 5. Attributes — How Data Is Stored

Everything in NTFS is an **attribute**. Every attribute has:
- Type (4-byte code)
- Length
- Resident/Non-resident flag
- Content

**Resident attribute:** Content stored DIRECTLY in the MFT record (small data fits):
```
$DATA attribute (resident, small file):
┌──────────────────────────────────────┐
│ Type: 0x80 ($DATA)                  │
│ Length: 72                           │
│ Resident flag: 1 (data is here)     │
│ Content length: 60                   │
│ [Actual file content — 60 bytes]    │
└──────────────────────────────────────┘
```

A file smaller than ~700 bytes can live ENTIRELY within the MFT record — zero extra I/O!

**Non-resident attribute:** Content stored in separate clusters, with a run list in the MFT record:
```
$DATA attribute (non-resident, large file):
┌──────────────────────────────────────────┐
│ Type: 0x80 ($DATA)                      │
│ Resident flag: 0 (data elsewhere)       │
│ Lowest VCN: 0                           │
│ Highest VCN: 131071                     │  ← 128MB file
│ Run list:                               │
│   Run 1: 65536 clusters @ LCN 40960    │  ← contiguous
│   Run 2: 65536 clusters @ LCN 200000   │  ← second contiguous run
└──────────────────────────────────────────┘

VCN = Virtual Cluster Number (position in file)
LCN = Logical Cluster Number (position on disk)
```

**Run list (data runs):**
Similar to ext4 extents, NTFS stores file data as runs (contiguous sequences of clusters):
```
Run encoding: length, start LCN (delta from previous)
[00] = end of runs
[11 10 00] = run of 0x10=16 clusters starting at LCN 0x00 (first run, absolute)
[11 20 FF F5] = run of 0x20=32 clusters starting at previous_LCN + 0xFFF5 (negative delta)
```

**Key attributes:**

| Attribute | Type | Description |
|-----------|------|-------------|
| $STANDARD_INFORMATION | 0x10 | Creation/modify/access times, flags, archive bit |
| $ATTRIBUTE_LIST | 0x20 | Points to extension MFT records |
| $FILE_NAME | 0x30 | Filename (UTF-16), parent directory reference |
| $OBJECT_ID | 0x40 | Globally unique identifier (GUID) |
| $SECURITY_DESCRIPTOR | 0x50 | ACL — who can access this file |
| $VOLUME_NAME | 0x60 | Volume label |
| $DATA | 0x80 | File content (the actual data!) |
| $INDEX_ROOT | 0x90 | Root of a B-tree index (used for directories) |
| $INDEX_ALLOCATION | 0xA0 | B-tree index for large directories |
| $BITMAP | 0xB0 | Free/used cluster bitmap |
| $REPARSE_POINT | 0xC0 | Symbolic link / junction / mount point data |
| $EA | 0x100 | Extended attributes (for POSIX compatibility) |

**Alternate Data Streams (ADS):**
A file can have multiple $DATA attributes with different names:
```
normalfile.txt:$DATA           ← default, unnamed data stream
normalfile.txt:hidden:$DATA    ← alternate stream named "hidden"
```
```cmd
:: Windows command to write to ADS:
echo "secret data" > normalfile.txt:hidden

:: Read ADS:
type normalfile.txt:hidden

:: Malware commonly hides in ADS
:: dir /r shows all streams
```

---

## 6. NTFS Journaling

NTFS uses a **transaction journal** stored in the `$LogFile` system file.

**How it works:**
```
Before any metadata change:
  1. Log intent: "I'm going to change X"
  2. Make the change
  3. Log completion: "Change X is done"

On recovery after crash:
  1. Read $LogFile
  2. Find incomplete transactions (intent logged, completion missing)
  3. For each incomplete: undo the change (roll back)
  4. File system returns to last consistent state
```

**NTFS journals metadata only** (like ext4 ordered mode by default). Data writes are not journaled unless you use specific APIs.

**$LogFile vs $UsnJrnl:**
- `$LogFile`: For crash recovery — short circular buffer, entries overwritten quickly
- `$UsnJrnl` (Update Sequence Number Journal): Change tracking for applications — records every file create/modify/delete. Applications like Windows Search index use it to track what's changed since they last ran.

---

## 7. NTFS Permissions (ACLs)

NTFS uses **Windows ACLs (Access Control Lists)** — much more granular than Unix rwxrwxrwx.

**Security Descriptor components:**
```
Security Descriptor:
  Owner SID:  S-1-5-21-....-1000  (owner's security ID)
  Group SID:  S-1-5-21-....-100   (primary group)
  
  DACL (Discretionary ACL) — who can access:
    ACE 1: Allow  S-1-5-21-...-1000  Read|Write|Execute  (Owner: full control)
    ACE 2: Allow  S-1-5-32-544       Read|Execute         (Administrators group)
    ACE 3: Allow  S-1-1-0            Read                  (Everyone: read only)
    
  SACL (System ACL) — what to audit:
    ACE: Audit  S-1-1-0  Write|Delete  (log write/delete attempts by anyone)
```

**NTFS file permissions:**
```
Permission         Windows Name     What it allows
Read               R                Read file content, view attributes
Write              W                Modify content, change attributes
Execute            X                Run as executable
Delete             D                Delete the file
Read Permissions   RC               See the ACL
Change Permissions WDAC             Modify the ACL
Take Ownership     WO               Become the owner
Full Control       F                All of the above
```

**Inheritance:**
Permissions can inherit from parent directories:
```
C:\Projects\  (DACL: Alice=Full, Bob=Read)
  └── docs\   (inherits: Alice=Full, Bob=Read) + can add own entries
       └── secret.docx  (additional: deny Bob=Read)  ← explicit deny overrides inherit
```

**SIDs (Security Identifiers):**
```
S-1-5-21-DOMAIN-USER_ID

Well-known SIDs:
  S-1-1-0       Everyone
  S-1-5-18      SYSTEM (local OS)
  S-1-5-32-544  BUILTIN\Administrators
  S-1-5-32-545  BUILTIN\Users
```

---

## 8. NTFS Special Features

**Compression:**
NTFS can compress individual files or entire directories. Uses LZ77 compression per 16-cluster unit. Transparent to applications.
```cmd
compact /c /s:C:\compressed_folder
```

**Encryption (EFS — Encrypting File System):**
Per-file encryption using a randomly generated file encryption key (FEK) encrypted with the user's RSA public key. Without the user's private key, the file is unreadable even if the disk is stolen.
```
File encryption:
  1. Generate random FEK (file encryption key)
  2. Encrypt file data with FEK using AES-256
  3. Encrypt FEK with user's RSA public key
  4. Store encrypted FEK in $EFS attribute
```

**Sparse files:**
A 1TB database file that's mostly zeros:
```
Without sparse: allocates 1TB of disk space filled with zeros
With sparse: stores only the non-zero regions; zero regions take no disk space

A process reads a "hole" in a sparse file → OS returns zeros without reading disk
```

**Reparse points (symlinks, junctions, mount points):**
```
NTFS Symlink:         C:\link → D:\target\dir  (like Unix symlink)
Junction:             C:\old → C:\new  (directory junction, same volume)
Mount point:          C:\data → \\Volume{GUID}\ (mount a volume at a path)
```

**Disk quotas:**
Limit how much disk space each user can use:
```
User Alice: quota 10GB / warning at 9GB
User Bob: no quota
```

---

## 9. On-Disk Layout

```
NTFS Volume:
┌─────────────────────────────────────────────────────────────────────┐
│ Boot Sector │ $MFT Zone       │ ... data ... │ $MFT Mirror (copy)  │
│ (sector 0)  │ (reserved 12.5%)│              │ (backup of 1st 4 MFT│
│ Offset 0    │                 │              │  records)            │
└─────────────────────────────────────────────────────────────────────┘

Boot Sector (BPB - BIOS Parameter Block for NTFS):
  Bytes per sector: 512
  Sectors per cluster: 8 (= 4KB clusters)
  MFT location: LCN (logical cluster number) of $MFT
  MFT Mirror location: LCN of $MFTMirr
  MFT record size: 1024 bytes
  Index block size: 4096 bytes
  Volume serial number: 8-byte random number
  NTFS signature: 0xAA55 (at offset 510)
```

**Cluster sizes:**
```
Volume size     Default cluster size
< 512MB         512 bytes
512MB – 1GB     1KB
1GB – 2GB       2KB
2GB – 4TB       4KB (most common)
4TB – 8TB       8KB
8TB – 16TB      16KB
16TB – 32TB     32KB
32TB – 256TB    64KB (Windows Server)
```

---

## 10. NTFS vs ext4 Comparison

| Feature | NTFS | ext4 |
|---------|------|------|
| Architecture | MFT + attributes | Inode table + separate data |
| Metadata storage | "Everything is a file" | Fixed inode + block bitmap |
| Journaling | Metadata (transaction log) | Ordered/writeback/data modes |
| Permissions | Windows ACLs (fine-grained) | Unix rwxrwxrwx + ACLs optional |
| Symlinks | Yes (reparse points) | Yes (dedicated type) |
| Hard links | Yes (up to ~4 billion) | Yes (up to 65,000 per file) |
| Max file size | 16 EB (theoretical) | 16TB (practical) |
| Max volume | 256 TB (Windows 10) | 1 EB |
| Compression | Built-in, per-file | Requires Btrfs or SquashFS |
| Encryption | Built-in EFS | Requires dm-crypt/fscrypt |
| Sparse files | Yes | Yes |
| Streams | Multiple $DATA per file | Extended attributes only |
| OS support | Native Windows; read-only Linux | Native Linux; limited Windows |
| Case sensitivity | Case-insensitive by default | Case-sensitive by default |

---

## Summary

| Concept | Description |
|---------|------------|
| MFT | Master File Table: array of 1KB records, one per file/directory |
| MFT record | Collection of attributes describing a file |
| Attribute | Named, typed data unit inside an MFT record |
| Resident attribute | Data stored directly inside MFT record (<~700B) |
| Non-resident attribute | Data in separate clusters, described by run list |
| Run list | List of (length, LCN) pairs describing contiguous disk runs |
| $DATA | Attribute holding file content (or directory index) |
| $STANDARD_INFORMATION | Timestamps, flags, archive bit |
| $FILE_NAME | UTF-16 filename and parent directory reference |
| $SECURITY_DESCRIPTOR | Windows ACL — who can access the file |
| $LogFile | Transaction journal for crash recovery |
| $UsnJrnl | Change tracking journal for applications |
| ADS | Alternate Data Streams: multiple named $DATA attributes per file |
| EFS | Encrypting File System: per-file AES encryption |
| Sparse files | Files with "holes" — zero ranges don't consume disk space |
| Reparse points | NTFS mechanism for symlinks, junctions, and mount points |
