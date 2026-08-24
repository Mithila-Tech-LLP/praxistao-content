# Chapter 23: File System Concepts

> **"A file system is the OS's answer to the question: how do you store data permanently and find it again? The answer — files organized in directories, identified by names, described by metadata — seems simple. But building it to be fast, reliable, and correct across power failures is one of the most complex engineering challenges in operating systems."**

---

## Table of Contents

1. [What Is a File System?](#1-what-is-a-file-system)
2. [Files — What They Really Are](#2-files--what-they-really-are)
3. [File Attributes and Metadata](#3-file-attributes-and-metadata)
4. [Directories — Organizing Files](#4-directories--organizing-files)
5. [Paths — Finding Files](#5-paths--finding-files)
6. [The Inode — The Heart of Unix File Systems](#6-the-inode--the-heart-of-unix-file-systems)
7. [Hard Links and Soft Links](#7-hard-links-and-soft-links)
8. [File System Operations](#8-file-system-operations)
9. [How Files Are Stored on Disk](#9-how-files-are-stored-on-disk)
10. [Journaling — Surviving Power Failures](#10-journaling--surviving-power-failures)
11. [File Descriptors and the Open File Table](#11-file-descriptors-and-the-open-file-table)
12. [Special Files — Not Just Regular Data](#12-special-files--not-just-regular-data)
13. [Summary](#summary)

---

## 1. What Is a File System?

A **file system** is the OS component that manages persistent storage. It answers:
- How to organize data so it can be found
- How to name data so it can be referenced
- How to store data so it survives power cycles
- How to protect data (permissions)
- How to share data between users and processes

**Without a file system:**
A 1TB disk is just 2 trillion random bytes with addresses. To find your document, you'd need to remember its exact byte offset. To add to it, you'd need to know exactly where it ends. Utterly impractical.

**With a file system:**
The same disk appears as a tree of named directories containing named files. `open("Documents/resume.pdf", O_RDONLY)` finds and opens the file. The file system handles all the low-level details.

**On-disk data structure:**
A file system is a data structure stored on a block device (disk, SSD, USB drive). It maps:
- File names → file metadata (inode)
- Inode → list of data blocks
- Free space tracking

---

## 2. Files — What They Really Are

**A file is a named, typed, persistent sequence of bytes.**

Properties:
- **Named:** Has a human-readable name
- **Typed:** Has a type (regular file, directory, symlink, device, pipe)
- **Persistent:** Survives power cycles and process exits
- **Sequence:** Ordered bytes from offset 0 to size-1
- **Extensible:** Can grow (write past end) or shrink (truncate)

**File types (Unix):**
```
- (dash)  Regular file: text, binary, executables, images, etc.
d         Directory: contains other files
l         Symbolic link: points to another file/directory
b         Block device: disk drive, USB (read/write in blocks)
c         Character device: keyboard, serial port (read/write stream)
p         Named pipe (FIFO): IPC mechanism
s         Unix socket: IPC mechanism
```

```bash
ls -la /dev/sda /etc/passwd /tmp/mypipe
brw-rw----  /dev/sda       # b = block device
-rw-r--r--  /etc/passwd    # - = regular file
prw-r--r--  /tmp/mypipe    # p = named pipe
```

---

## 3. File Attributes and Metadata

Files have more than just content — they have **metadata**:

```
File: /home/user/documents/report.pdf

Metadata:
  Name:         report.pdf
  Type:         regular file
  Size:         1,048,576 bytes (1MB)
  Inode:        94372
  Owner:        user (UID 1000)
  Group:        users (GID 100)
  Permissions:  -rw-r--r-- (owner: rw, group: r, others: r)
  Access time:  2026-07-01 09:30:00 (last read)
  Modify time:  2026-06-28 14:15:30 (last write)
  Change time:  2026-06-28 14:15:30 (last metadata change)
  Link count:   1 (number of hard links to this inode)
  Blocks:       2048 (number of 512-byte blocks allocated)
  Block size:   4096 (optimal I/O block size)
```

**Getting metadata in C:**
```c
#include <sys/stat.h>

struct stat info;
stat("/etc/passwd", &info);

printf("Size: %ld bytes\n", info.st_size);
printf("Inode: %lu\n", info.st_ino);
printf("UID: %u\n", info.st_uid);
printf("Permissions: %o\n", info.st_mode & 0777);
printf("Links: %lu\n", info.st_nlink);
```

**Unix permission bits:**
```
-rw-r--r--
│ │  │  │
│ │  │  └── Others: r (read), no write, no exec
│ │  └───── Group:  r (read), no write, no exec
│ └──────── Owner:  rw (read+write), no exec
└────────── Type: - (regular file)

Encoded as octal: 0644
  0600 = owner read+write
  0040 = group read
  0004 = others read
  Total: 0644
```

**Execute permission:**
For files: can be executed as a program.
For directories: can "enter" (cd into) the directory.

---

## 4. Directories — Organizing Files

A **directory** is a special file that maps file names to inode numbers.

**Directory structure:**
```
Directory: /home/user/documents/
  Entry: "."        → inode 94369  (this directory itself)
  Entry: ".."       → inode 94200  (parent directory)
  Entry: "resume.pdf" → inode 94372
  Entry: "notes.txt"  → inode 94380
  Entry: "photos"     → inode 94400  (subdirectory)
```

A directory is JUST a list of (name, inode) pairs. That's all. The name "resume.pdf" only exists in the directory. The inode 94372 doesn't know its own name.

**This has consequences:**
- Multiple directory entries can point to the same inode (hard links)
- Moving a file within the same file system: just rename the directory entry (no data copying)
- Moving a file across file systems: must copy data (different inodes on different FS)

---

## 5. Paths — Finding Files

**Absolute path:** Starts from the root `/`:
```
/home/user/documents/report.pdf
```

**Relative path:** Relative to the current working directory:
```
documents/report.pdf        (if CWD is /home/user)
./documents/report.pdf      (same as above)
../user2/data.txt           (go up one level, then into user2)
```

**Path resolution:**
```
Resolve: /home/user/documents/report.pdf

1. Start at root inode (inode 2, always)
2. Find "home" in root directory entries → inode 500
3. Check inode 500: is it a directory? yes
4. Find "user" in inode 500's directory → inode 94200
5. Check inode 94200: is it a directory? yes
6. Find "documents" in inode 94200's directory → inode 94369
7. Check inode 94369: is it a directory? yes
8. Find "report.pdf" in inode 94369's directory → inode 94372
9. Return inode 94372
```

Each `.` = current directory, `..` = parent directory.

---

## 6. The Inode — The Heart of Unix File Systems

The **inode (index node)** is the fundamental data structure of Unix file systems. Every file/directory has exactly one inode.

**What an inode contains:**
```
inode #94372:
  type:        regular file
  size:        1,048,576 bytes
  uid:         1000  (owner)
  gid:         100   (group)
  mode:        100644 (regular, rw-r--r--)
  atime:       2026-07-01 09:30:00
  mtime:       2026-06-28 14:15:30
  ctime:       2026-06-28 14:15:30
  link count:  1
  data blocks: [block 4096, block 4097, block 4098, ..., block 5119]
               (256 blocks × 4096 bytes = 1MB)
```

**What an inode does NOT contain:**
- The file's NAME (the name is in the directory entry)
- Any file content (content is in the data blocks the inode points to)

**inode allocation:**
File systems pre-allocate a fixed number of inodes at format time:
```bash
# Check inode usage:
df -i /
# Filesystem  Inodes  IUsed  IFree  IUse%  Mounted on
# /dev/sda1   6553600  45000  6508600  1%   /

# Can run out of inodes even with free disk space!
# (Creating millions of tiny files can exhaust inodes)
```

**Block pointers in inode:**
For large files, inodes use indirect blocks:
```
inode data blocks:
  0..11:    Direct pointers (12 × 4KB = 48KB directly)
  12:       Indirect block pointer → points to a block containing 1024 more pointers
              (1024 × 4KB = 4MB)
  13:       Double indirect → block of pointers → each → block of pointers
              (1024 × 1024 × 4KB = 4GB)
  14:       Triple indirect → 1024 × 1024 × 1024 × 4KB = 4TB
```

ext2/ext3 used this scheme. ext4 uses **extents** instead (see Chapter 25).

---

## 7. Hard Links and Soft Links

**Hard link:**
Multiple directory entries pointing to the SAME inode.

```bash
ln original.txt hardlink.txt   # create hard link
```

```
Directory:
  "original.txt" → inode 1234
  "hardlink.txt" → inode 1234  ← same inode!

inode 1234: link count = 2

ls -li original.txt hardlink.txt
94372 -rw-r--r-- 2 user user 1000 Jul 1 ... hardlink.txt
94372 -rw-r--r-- 2 user user 1000 Jul 1 ... original.txt
# Both have the same inode number (94372) and link count (2)
```

**Properties of hard links:**
- Same permissions, same size, same content (same inode)
- Deleting one doesn't delete the data — decrements link count
- Data deleted only when link count = 0 AND no open file descriptors
- Cannot cross file system boundaries (inode numbers are local to a FS)
- Cannot link to directories (prevents loops in the directory tree — special case for . and ..)

**Soft link (symbolic link):**
A file whose content is a PATH to another file/directory.

```bash
ln -s /home/user/documents symlink_docs   # create symlink
```

```
Directory:
  "symlink_docs" → inode 9999 (contains: "/home/user/documents")
  
inode 9999: type = symlink, content = "/home/user/documents"
```

**Properties of soft links:**
- Its own inode, can cross file systems
- Can point to a directory
- Can be broken (dangling symlink) if target is deleted
- Transparent to most operations (open on a symlink → opens the target)

```bash
ls -la /lib    # most system links are symlinks
# /lib -> usr/lib     (symlink to /usr/lib)
```

---

## 8. File System Operations

Key operations that file systems must support:

**Creating a file:**
1. Allocate a new inode
2. Initialize inode metadata (type, permissions, timestamps, uid, gid)
3. Add a directory entry `(name → inode)` to the parent directory

**Writing to a file:**
1. Find the inode for the file
2. Determine which data blocks to write to (allocate new blocks if needed)
3. Write data to the block(s)
4. Update inode: file size, mtime, block list

**Reading a file:**
1. Find the inode
2. Determine which data blocks contain the requested offset range
3. Read data from those blocks
4. Update inode: atime (access time)

**Deleting a file (unlink):**
1. Remove the directory entry
2. Decrement link count in inode
3. If link count == 0 AND no open fds: free data blocks, free inode

**Rename/move (same FS):**
1. Add directory entry in destination directory
2. Remove directory entry from source directory
3. No data copying — inode stays the same!

**Rename/move (different FS):**
1. Copy all data to new location
2. Create new inode in destination FS
3. Delete source

---

## 9. How Files Are Stored on Disk

A disk (or partition) formatted with a file system is divided into regions:

**ext4 on-disk layout:**
```
Partition:
┌─────────────────────────────────────────────────────────────────┐
│ Superblock │ Block Group 0 │ Block Group 1 │ ... │ Block Group N│
└─────────────────────────────────────────────────────────────────┘

Block Group layout:
┌─────────────┬────────────┬────────────┬────────────┬───────────┐
│ Superblock  │ GDT        │ Block      │ Inode      │ Data      │
│ (copy)      │ (copy)     │ Bitmap     │ Table      │ Blocks    │
└─────────────┴────────────┴────────────┴────────────┴───────────┘
  (optional)    (optional)   (1 block)   (many inodes)  (most space)
```

**Superblock:**
The first block after the boot area. Contains:
- Magic number (identifies file system type)
- Block size (1KB, 2KB, 4KB)
- Total blocks, free blocks
- Total inodes, free inodes
- Time last mounted, last checked
- File system state (clean/dirty)

```bash
# View superblock:
tune2fs -l /dev/sda1 | head -30
```

**Block bitmap:**
One bit per data block: 0 = free, 1 = used.
Fast allocation: `find_first_zero_bit()` → O(1) with BITMAPFFS instruction.

**Inode table:**
Array of inodes. Each inode is 128 bytes (ext2/3) or 256 bytes (ext4).
`inode_number → table_base + (inode_number - 1) × inode_size`

---

## 10. Journaling — Surviving Power Failures

**The crash problem:**
Writing a file requires multiple disk writes (inode update, data blocks, directory entry). If power fails halfway through, the disk is in an inconsistent state.

**Example:**
Appending to a file:
1. Write new data to block 5000 ← power fails here!
2. Update inode: size += 4096, blocks += 5000
3. Update block bitmap: block 5000 is used

Power fails after step 1:
- Data is in block 5000 on disk
- Inode still shows old size (doesn't know about block 5000)
- Block 5000 not in bitmap (appears free)
- Block 5000 might be allocated to another file later → data corruption!

**Old solution: fsck (file system check)**
On boot after crash, run fsck — scans entire file system for inconsistencies and fixes them. On a 1TB disk: takes 10–30 minutes. Unacceptable for servers!

**Journaling:**
Before making changes to the file system, write a log entry (journal entry) describing what you're about to do. Then make the changes. If power fails:
- Before log entry written: nothing happened, nothing to recover
- After log entry, before changes: replay the log on next boot (fast!)
- After changes AND log committed: mark log as done (clean)

**Journal location:** A reserved area of the file system (ext4: `/dev/sda1` journal at the beginning of the partition).

**Journaling modes (ext4):**
```
writeback:  Only journal metadata (not data). Fastest. Data may be stale after crash.
ordered:    Journal metadata; data written to disk BEFORE metadata committed. Default.
data:       Journal everything (metadata AND data). Slowest. Safest.
```

```bash
# Check ext4 journal mode:
tune2fs -l /dev/sda1 | grep "Default mount options"
# Default mount options:    user_xattr acl
# data=ordered (default mode)
```

---

## 11. File Descriptors and the Open File Table

When a process opens a file, the OS creates a **file descriptor** (fd) — an integer handle.

**Three-level structure:**
```
Process 1:                                           
  fd[0] = 0 → [Open File Table entry 5]  → inode 12345 (stdin)
  fd[1] = 1 → [Open File Table entry 6]  → inode 12346 (stdout)
  fd[2] = 2 → [Open File Table entry 7]  → inode 12347 (stderr)
  fd[3] = 3 → [Open File Table entry 10] → inode 94372 (report.pdf)

Process 2:
  fd[0] = 0 → [Open File Table entry 5]  → inode 12345 (stdin - shared!)
  fd[3] = 3 → [Open File Table entry 15] → inode 94372 (same file, different offset!)

Open File Table (kernel):
  Entry 5:  offset=0,      flags=O_RDONLY, inode=12345, refcount=2
  Entry 6:  offset=0,      flags=O_WRONLY, inode=12346, refcount=1
  Entry 10: offset=524288, flags=O_RDONLY, inode=94372, refcount=1
  Entry 15: offset=0,      flags=O_RDONLY, inode=94372, refcount=1
```

**Per-file-descriptor info (open file table entry):**
- **Current offset:** Where the next read/write starts
- **Flags:** O_RDONLY, O_WRONLY, O_RDWR, O_NONBLOCK, O_APPEND, etc.
- **Reference count:** How many fds point to this entry

**File descriptor inheritance:**
Child processes (after fork) inherit parent's file descriptors. Both process's fd[3] points to the same open file table entry → shared offset! (That's why `fork()` + `dup2()` is used for shell pipes.)

**Predefined fds:**
- 0: stdin (usually terminal input)
- 1: stdout (usually terminal output)
- 2: stderr (usually terminal error output)

---

## 12. Special Files — Not Just Regular Data

Unix "everything is a file" philosophy extends to non-regular-file objects:

**/proc pseudo-filesystem:**
```
/proc/1234/       — directory for process 1234
/proc/1234/maps   — memory map
/proc/1234/status — process status
/proc/cpuinfo     — CPU information
/proc/meminfo     — memory stats
/proc/sys/        — kernel parameters (sysctl)
```
Reading these "files" runs kernel code to generate data on the fly.

**/sys (sysfs):**
```
/sys/block/sda/queue/read_ahead_kb — disk readahead setting
/sys/class/net/eth0/speed          — network interface speed
/sys/devices/                      — device tree
```
Exposes hardware devices and kernel state as files.

**/dev (devices):**
```
/dev/sda          — first SATA disk
/dev/sda1         — first partition
/dev/null         — discards all writes; reads return EOF
/dev/zero         — reads return infinite zeros
/dev/urandom      — reads return random bytes
/dev/tty          — current terminal
/dev/loop0        — loop device (mount a file as a disk)
```

**Pipes as files:**
Named pipes (FIFOs) appear in the file system and can be accessed like files.

---

## Summary

| Concept | Definition |
|---------|-----------|
| File | Named, typed, persistent sequence of bytes |
| Directory | Special file mapping names to inodes |
| Inode | Metadata + block list for a file/directory |
| Path resolution | Walk directory tree using inode numbers |
| Hard link | Multiple names pointing to same inode |
| Soft link | File containing a path to another file |
| Journaling | Log of pending changes; ensures consistency after crashes |
| File descriptor | Integer handle to an open file; per-process |
| Open file table | Kernel table of open files with offset and flags |
| Superblock | File system metadata (block size, inode count, etc.) |
| Block bitmap | Tracks which blocks are free/used |
| Inode table | Array of inodes; one per file/directory |
| /proc /sys /dev | Special pseudo-filesystems exposing kernel state as files |
