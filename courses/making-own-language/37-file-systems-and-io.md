# Chapter 37: File Systems, I/O, and System Calls

> "Everything is a file."
> — The Unix Philosophy

---

## Overview

When an Astra program reads a source file, writes compiled output, or opens a network socket, it is interacting with one of the most fundamental abstractions in an operating system: the **file**. Unix-like systems take this abstraction to an extreme — not only disk files, but network connections, hardware devices, inter-process pipes, and even the terminal are all treated as files, opened with `open()`, read with `read()`, and closed with `close()`. This uniform interface is one of the great design decisions in computer science.

In this chapter we explore the full stack of file and I/O operations: from the physical disk all the way up to the user-facing file API in Astra. We will study how file systems organize data with inodes and directory entries, how the kernel's Virtual File System (VFS) layer provides a uniform interface across completely different on-disk formats, and how Go's I/O primitives map to the underlying system calls. We will also study non-blocking I/O and event loops — the technology that lets a web server handle ten thousand simultaneous connections with a single OS thread. By the end of this chapter, you will have implemented Astra's complete file standard library.

---

## What We're Building

We will implement Astra's **file standard library** — the `file` module that Astra programs use for all file operations. This includes reading, writing, appending, listing directories, checking existence, and path manipulation. We will show exactly which OS system call each Astra function triggers.

---

## Table of Contents

1. What Is a File System?
2. Inodes — Metadata Without Names
3. Directory Entries — Names Without Data
4. Hard Links vs Soft Links
5. The Virtual File System (VFS)
6. File Descriptors — The OS Interface to Files
7. System Calls: open, read, write, close
8. Buffered vs Unbuffered I/O
9. Memory-Mapped Files
10. Pipes and FIFOs
11. Non-Blocking I/O and I/O Multiplexing
12. Event Loops — How Web Servers Handle Thousands of Connections
13. File I/O in Go
14. Astra Build Milestone: The Complete File Standard Library
15. Exercises
16. Summary

---

## 1. What Is a File System?

A file system is the layer of software that organizes data on a storage device (SSD, HDD, USB drive) into a hierarchy of named files and directories that humans and programs can work with.

Without a file system, a 1TB SSD is just 1 trillion bytes of raw storage — there is no concept of "a file called main.astra that contains this text." The file system imposes structure: it divides the disk into blocks, keeps a catalog of which blocks belong to which file, tracks permissions, and provides the directory tree that lets you navigate from `/` to `/home/adityapathak/code/astra/main.astra`.

```
Physical Disk (raw blocks):
┌──────────────────────────────────────────────────────┐
│ Block 0 │ Block 1 │ Block 2 │ Block 3 │ ... │ Block N│
└──────────────────────────────────────────────────────┘

File System Layer:
┌──────────────────────────────────────────────────────┐
│  Superblock │ Inode table │ Data blocks │ Bitmap     │
│  (FS info)  │  (metadata) │  (content)  │ (free map) │
└──────────────────────────────────────────────────────┘

User's View (directory hierarchy):
/
├── home/
│   └── adityapathak/
│       └── code/
│           └── astra/
│               ├── main.astra
│               └── lexer.astra
├── etc/
│   └── hosts
└── usr/
    └── bin/
        └── ls
```

Common file systems:
- **ext4**: the default Linux file system. Journaling (keeps a log of changes to recover from crashes).
- **APFS**: Apple File System, used on macOS and iOS. Copy-on-write, snapshots, native encryption.
- **NTFS**: Windows file system. Journaling, ACLs, alternate data streams.
- **FAT32/exFAT**: legacy/cross-platform (USB drives). Simple but limited (4GB max file size for FAT32).
- **tmpfs**: an in-memory file system (Linux). Files live in RAM; they vanish on reboot. Used for `/tmp`.
- **NFS**: Network File System. Files live on a remote server; accessed over the network using the same API.

---

## 2. Inodes — Metadata Without Names

Every file on a Unix file system has an **inode** (index node) — a fixed-size data structure that stores everything about a file EXCEPT its name.

```
┌─────────────────────────────────────────────────┐
│                    INODE                        │
├─────────────────────────────────────────────────┤
│ Inode number:  1234                             │
│ File type:     regular file                     │
│ Permissions:   rw-r--r-- (0644)                 │
│ Owner (UID):   1000 (adityapathak)              │
│ Group (GID):   1000                             │
│ Size:          4096 bytes                       │
│ Hard link count: 1                              │
│ Access time (atime): 2026-06-07 10:30:00        │
│ Modify time (mtime): 2026-06-07 09:15:00        │
│ Change time (ctime): 2026-06-07 09:15:00        │
│                                                 │
│ Block pointers:                                 │
│   Direct[0]:   block 5001  ─────────────────┐  │
│   Direct[1]:   block 5002  ──────────────┐  │  │
│   Direct[2]:   0 (file smaller than 3 blocks)  │
│   ...                                       │  │
│   Indirect:    block 7000 (points to 512    │  │
│                 more block pointers)        │  │
│   Double indirect: ...                      │  │
└─────────────────────────────────────────────────┘
          │                       │
          ▼                       ▼
    ┌───────────┐           ┌───────────┐
    │  Block    │           │  Block    │
    │  5001     │           │  5002     │
    │  (data)   │           │  (data)   │
    └───────────┘           └───────────┘
```

The inode does NOT contain the file name. The name is stored in the directory. This design allows one inode to have multiple names (hard links) and allows directories to be searched independently of file data.

**ls -li** shows inode numbers:
```
$ ls -li /home/adityapathak/code/astra/
1234 -rw-r--r-- 1 adityapathak 4096 Jun 7 09:15 main.astra
1235 -rw-r--r-- 1 adityapathak 2048 Jun 7 09:10 lexer.astra
```

---

## 3. Directory Entries — Names Without Data

A directory is itself a file, whose content is a list of **directory entries** (dirents): name → inode number mappings.

```
Directory: /home/adityapathak/code/astra/
┌──────────────────────────────────────┐
│ dirent: "."        → inode 1230      │ (this directory itself)
│ dirent: ".."       → inode 1225      │ (parent directory)
│ dirent: "main.astra" → inode 1234    │
│ dirent: "lexer.astra" → inode 1235   │
│ dirent: "parser.astra" → inode 1236  │
└──────────────────────────────────────┘
```

To open `/home/adityapathak/code/astra/main.astra`, the kernel:
1. Starts at the root inode (always inode 2)
2. Looks up "home" in root's directory → inode 1100
3. Looks up "adityapathak" in inode 1100's directory → inode 1200
4. Looks up "code" in inode 1200's directory → inode 1225
5. Looks up "astra" in inode 1225's directory → inode 1230
6. Looks up "main.astra" in inode 1230's directory → inode 1234
7. Opens inode 1234 → allocates a file descriptor → returns it

Path resolution is O(d) where d is the depth of the path — this is why deeply nested paths are slightly slower than shallow ones.

---

## 4. Hard Links vs Soft Links

**Hard link**: a directory entry that points to an existing inode. Multiple directory entries can point to the SAME inode.

```
/home/user/docs/report.astra    ─┐
                                  ├──► inode 1234 (refcount = 2)
/home/user/backup/report.astra  ─┘

The file data has two names. Deleting one name decrements refcount.
The file is only deleted when refcount reaches 0 (all names removed).
```

**Soft link (symbolic link)**: a special file whose content is a path string. Following the link causes the kernel to re-traverse the path.

```
/home/user/docs/current.astra → "/home/user/docs/report_v3.astra"
                                  ↑ stored as data in the symlink's inode

Dangling symlink: if the target is deleted, the symlink still exists
but following it gives ENOENT (no such file or directory).
```

Hard links are fast and transparent; symlinks are flexible (can point to directories, can cross file system boundaries) but add one extra path resolution.

---

## 5. The Virtual File System (VFS)

Your system may have an ext4 disk, an APFS USB drive, an NFS network share, and a tmpfs all mounted simultaneously. Yet `open("/mnt/usb/file.txt")` works the same as `open("/home/user/file.txt")`. How?

The **Virtual File System (VFS)** is a kernel abstraction layer. It defines a standard interface (a set of C function pointers) that every concrete file system must implement:

```mermaid
flowchart TD
    UP["User program\nopen(\"/mnt/usb/file.txt\", O_RDONLY)"]
    SC["System call interface"]
    VFS["VFS Layer\npath_lookup() → find inode\ninode→i_op→open() → FS-specific open\nfile→f_op→read() → FS-specific read\nfile→f_op→write() → FS-specific write"]
    EXT4["ext4\n(local disk)"]
    APFS["APFS\n(USB drive)"]
    NFS["NFS\n(network)"]
    UP --> SC --> VFS
    VFS --> EXT4
    VFS --> APFS
    VFS --> NFS
```

The VFS defines key kernel objects:
- **superblock**: represents a mounted file system
- **inode**: represents a file or directory
- **dentry** (directory entry): represents a path component (cached in the dentry cache for performance)
- **file**: represents an open file in a process (has position, flags, etc.)

Because of VFS, your code that calls `open()` works identically across all file systems.

---

## 6. File Descriptors — The OS Interface to Files

When a process opens a file, the kernel creates an **open file description** (kernel-side state) and returns a small integer to the process: the **file descriptor** (fd).

```
Process (PID 4217)
File Descriptor Table:
┌────┬──────────────────────────────────────────┐
│ 0  │  stdin (keyboard/pipe)                   │
│ 1  │  stdout (terminal/pipe)                  │
│ 2  │  stderr (terminal)                       │
│ 3  │  open file: /home/user/main.astra        │
│ 4  │  open socket: 192.168.1.1:80             │
│ 5  │  pipe read end                           │
└────┴──────────────────────────────────────────┘

Kernel: Open File Table (shared across processes)
┌────────────────────────────────────────────────┐
│  File Position: 1024                           │
│  Flags: O_RDONLY                               │
│  Ref count: 1                                  │
│  Pointer → inode 1234                          │
└────────────────────────────────────────────────┘
```

File descriptors 0, 1, 2 (stdin, stdout, stderr) are opened automatically when a process starts. This is why `printf` works without opening any file — it writes to fd 1.

File descriptors are just integers. This means:
- You can `dup2(fd, 1)` to redirect stdout to a file
- `fork()` copies the fd table — parent and child share the same open file description (same position!)
- Passing fd 3 to another process via a Unix socket transfers access to that open file

---

## 7. System Calls: open, read, write, close

These four system calls are the heart of all file I/O on Unix:

```c
#include <fcntl.h>
#include <unistd.h>

// open(): returns fd or -1 on error
// flags: O_RDONLY, O_WRONLY, O_RDWR, O_CREAT, O_TRUNC, O_APPEND
// mode: permissions if creating (e.g., 0644)
int fd = open("/home/user/main.astra", O_RDONLY);

// read(): reads up to n bytes into buf
// Returns: bytes actually read (0 = EOF, -1 = error)
char buf[4096];
ssize_t bytes_read = read(fd, buf, sizeof(buf));

// write(): writes n bytes from buf to fd
// Returns: bytes actually written (-1 = error)
ssize_t bytes_written = write(fd, "hello\n", 6);

// close(): releases the file descriptor
close(fd);
```

Each of these is a **system call** — a transition from user mode to kernel mode. The CPU switches privilege levels, the kernel handles the request, then switches back. This is expensive (~100-1000ns per call). Minimizing system call count is a key I/O optimization.

**lseek()**: moves the file position (like `Seek` in Go's `io.Seeker`):

```c
// Seek to byte 1024 from the start
lseek(fd, 1024, SEEK_SET);

// Seek 100 bytes forward from current position
lseek(fd, 100, SEEK_CUR);

// Seek to 10 bytes before end-of-file
lseek(fd, -10, SEEK_END);
```

**stat()**: get file metadata without opening:

```c
struct stat st;
stat("/home/user/main.astra", &st);
printf("size: %lld\n", st.st_size);
printf("permissions: %o\n", st.st_mode & 0777);
```

---

## 8. Buffered vs Unbuffered I/O

System calls are expensive (kernel crossing). If your program writes one character at a time with 1000 `write()` calls, that's 1000 kernel crossings. Much better to buffer writes in user space and flush in one big `write()` call.

**C stdio (buffered I/O)**: The standard C library wraps raw `read()`/`write()` with a buffer:

```c
FILE* f = fopen("output.txt", "w"); // buffered
fprintf(f, "hello");    // writes to 4KB buffer in user space
fprintf(f, " world\n"); // still just in the buffer
fflush(f);              // NOW calls write() once with "hello world\n"
fclose(f);              // flushes and close()s
```

Three buffering modes:
- **Fully buffered**: writes accumulate until the buffer is full (or fflush). Used for regular files.
- **Line buffered**: writes flush at newlines. Used for terminals (this is why `printf("hi")` without `\n` sometimes doesn't appear immediately).
- **Unbuffered**: every write goes directly to the OS. Used for `stderr` (so error messages appear immediately even if the program crashes).

**Go's bufio**: Go's standard library provides explicit buffering via `bufio.NewReader` and `bufio.NewWriter`:

```go
// Without buffering: many small read() syscalls
f, _ := os.Open("big_file.txt")
buf := make([]byte, 1)
for {
    n, _ := f.Read(buf) // ONE syscall per byte — terrible!
    if n == 0 { break }
}

// With buffering: reads 4KB at a time, serves small reads from buffer
f, _ := os.Open("big_file.txt")
br := bufio.NewReader(f) // 4KB buffer by default
for {
    line, err := br.ReadString('\n') // O(1) for lines in buffer
    if err != nil { break }
    process(line)
}
```

**O_DIRECT flag**: bypass the kernel's page cache entirely — data goes directly between user buffer and disk. Used by databases (they manage their own caching) to avoid double buffering.

---

## 9. Memory-Mapped Files

`mmap()` maps a file (or a range of a file) directly into the process's virtual address space. Instead of calling `read()`, you just access memory addresses — the kernel brings in the data as page faults.

```c
#include <sys/mman.h>

int fd = open("source.astra", O_RDONLY);
struct stat st;
fstat(fd, &st);

// Map the entire file into virtual memory
char* data = mmap(NULL, st.st_size, PROT_READ, MAP_PRIVATE, fd, 0);

// Now access the file's content as a C string!
printf("First byte: %c\n", data[0]);
printf("100th byte: %c\n", data[99]);

munmap(data, st.st_size);
close(fd);
```

**How it works**: `mmap()` sets up page table entries for the file range, but doesn't read any data yet. On first access to each page, a page fault triggers the kernel to read that 4KB chunk from disk into a physical frame. Subsequently, the data is in the kernel's page cache — subsequent accesses are as fast as RAM.

**Benefits**:
- Zero-copy: no explicit `read()` call, no user-space buffer allocation
- Random access is O(1) — just index into the mapped pointer
- The OS automatically manages what to keep in RAM vs evict to disk

**The Astra compiler uses mmap for lexing**: reading the entire source file via `mmap` and scanning it byte-by-byte is the fastest possible approach, especially for large files.

**Shared memory via mmap**: `MAP_SHARED` allows multiple processes to map the same file and see each other's writes — a fast inter-process communication mechanism used by databases and message queues.

---

## 10. Pipes and FIFOs

A **pipe** is a one-directional byte stream connecting two ends: a write end and a read end. Data written to the write end can be read from the read end.

```c
int pipefd[2]; // pipefd[0] = read end, pipefd[1] = write end
pipe(pipefd);

pid_t child = fork();
if (child == 0) {
    // Child: read from pipe
    close(pipefd[1]); // close write end (child doesn't write)
    char buf[256];
    read(pipefd[0], buf, sizeof(buf));
    printf("Child received: %s\n", buf);
    close(pipefd[0]);
    exit(0);
} else {
    // Parent: write to pipe
    close(pipefd[0]); // close read end (parent doesn't read)
    write(pipefd[1], "hello from parent", 17);
    close(pipefd[1]); // closing write end sends EOF to reader
    wait(NULL);
}
```

When you run `cat file.txt | grep error | wc -l` in your shell, the shell creates two pipes and connects the standard output of each command to the standard input of the next. This is how Unix pipelines work.

**FIFO (named pipe)**: `mkfifo("/tmp/myfifo", 0644)` creates a pipe that has a name in the file system. Any process can open it by name — unrelated processes can communicate this way.

---

## 11. Non-Blocking I/O and I/O Multiplexing

The standard `read()` call **blocks** — the calling thread sleeps until data arrives. For a server handling 10,000 clients, you cannot afford a thread per client (10,000 threads × 8MB stack = 80GB RAM). You need to handle many connections in ONE thread.

The key: **non-blocking I/O + I/O multiplexing**.

With `O_NONBLOCK`, `read()` returns immediately with `EAGAIN` if no data is available (instead of blocking). The program can then try other file descriptors.

**select()** and **poll()**: wait for one of N file descriptors to become ready (readable/writable):

```c
fd_set readfds;
FD_ZERO(&readfds);
for (int i = 0; i < num_clients; i++) {
    FD_SET(client_fds[i], &readfds);
}
struct timeval timeout = { .tv_sec = 1 };

// Blocks until at least one fd is ready, or timeout expires
int ready = select(maxfd + 1, &readfds, NULL, NULL, &timeout);

for (int i = 0; i < num_clients; i++) {
    if (FD_ISSET(client_fds[i], &readfds)) {
        // This client has data — read it
        handle_client(client_fds[i]);
    }
}
```

`select()` has a problem: it requires copying the entire fd_set to the kernel on every call, and scanning all N fds even if only 1 is ready. For N=10,000, this is O(N) per event — terrible.

**epoll() (Linux)**: the modern solution. Register file descriptors once; get notified only of active ones. O(1) per active fd, regardless of total connections.

```c
// Create epoll instance
int epfd = epoll_create1(0);

// Register a file descriptor for monitoring
struct epoll_event ev;
ev.events = EPOLLIN;     // notify when readable
ev.data.fd = client_fd;
epoll_ctl(epfd, EPOLL_CTL_ADD, client_fd, &ev);

// Wait for events (blocks until at least 1 fd is ready)
struct epoll_event events[64];
int n = epoll_wait(epfd, events, 64, -1); // -1 = no timeout

for (int i = 0; i < n; i++) {
    int fd = events[i].data.fd;
    handle_client(fd); // only active clients
}
```

**kqueue** (macOS/BSD): equivalent to epoll, used on macOS. Go's `net` package uses kqueue on macOS and epoll on Linux internally.

---

## 12. Event Loops — How Web Servers Handle Thousands of Connections

epoll enables the **event loop** pattern: a single thread manages thousands of connections by reacting to I/O events.

```
┌─────────────────────────────────────────────────────────┐
│                  EVENT LOOP                             │
│                                                         │
│  while true:                                            │
│    events = epoll_wait(epfd)   // wait for any I/O      │
│                                                         │
│    for each event:                                      │
│      if event is "new connection":                      │
│        accept(); register new fd with epoll             │
│      if event is "data readable on client fd":          │
│        read request from fd                             │
│        process request (may do async DB call)           │
│        write response to fd                             │
│      if event is "fd writable":                         │
│        flush pending write buffer to fd                 │
└─────────────────────────────────────────────────────────┘

10,000 clients connected, but only ~10 active at once
→ 1 thread handles all of them efficiently
```

Node.js, Nginx, Redis, and Go's `net/http` all use event loops (or variants like goroutines-per-connection with Go's netpoller, which uses epoll/kqueue under the hood). Go's runtime has a special goroutine scheduler integration with epoll: when a goroutine blocks on network I/O, the runtime registers the fd with epoll and parks the goroutine; when the fd becomes ready, the goroutine is resumed. This gives you the simplicity of blocking reads in goroutines with the efficiency of epoll.

---

## 13. File I/O in Go

Go provides a layered I/O system:

```go
// Layer 1: raw syscall wrappers (os package)
f, err := os.Open("file.txt")         // open for reading
f, err := os.Create("output.txt")     // create/truncate for writing
f, err := os.OpenFile("file.txt",     // full control
    os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)

// Reading
data, err := os.ReadFile("file.txt")  // read entire file
buf := make([]byte, 4096)
n, err := f.Read(buf)                 // read up to 4096 bytes

// Writing
os.WriteFile("out.txt", data, 0644)   // write entire file
f.Write([]byte("hello"))              // write bytes
f.WriteString("hello")                // write string

// Layer 2: buffered I/O (bufio package)
br := bufio.NewReader(f)
line, err := br.ReadString('\n')
scanner := bufio.NewScanner(f)        // line-by-line scanning
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

bw := bufio.NewWriter(f)
bw.WriteString("buffered write")
bw.Flush()                            // must flush!

// Layer 3: path manipulation (filepath package)
abs, _ := filepath.Abs("../file.txt")   // absolute path
dir   := filepath.Dir("/a/b/c.txt")     // "/a/b"
base  := filepath.Base("/a/b/c.txt")    // "c.txt"
ext   := filepath.Ext("file.astra")     // ".astra"
joined := filepath.Join("dir", "sub", "file.txt") // "dir/sub/file.txt"

// Directory operations
entries, err := os.ReadDir(".")
for _, entry := range entries {
    fmt.Printf("%s (dir: %v)\n", entry.Name(), entry.IsDir())
}
os.Mkdir("newdir", 0755)
os.MkdirAll("a/b/c/d", 0755) // create all intermediate dirs
os.Remove("file.txt")
os.RemoveAll("dir/")
os.Rename("old.txt", "new.txt")

// File info
info, err := os.Stat("file.txt")
info.Size()       // bytes
info.ModTime()    // time.Time
info.IsDir()      // bool
info.Mode()       // permissions
```

---

## 14. Astra Build Milestone: The Complete File Standard Library

Here is Astra's complete `file` standard library, implemented in Go:

```go
// stdlib/file/file.go
// The Astra file standard library.
// Astra programs import this as: import "std/file"
//
// Each function here corresponds to a system call sequence:
//   ReadFile  → open() + read() × N + close()
//   WriteFile → open(O_WRONLY|O_CREATE|O_TRUNC) + write() + close()
//   AppendFile → open(O_WRONLY|O_APPEND|O_CREATE) + write() + close()
//   Exists    → stat() (no open needed)
//   ...

package file

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// ─── Reading ────────────────────────────────────────────────────────────────

// ReadFile reads the entire contents of the file at path.
// Syscalls: open(O_RDONLY) + read() × N + close()
func ReadFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("file.ReadFile(%q): %w", path, err)
    }
    return string(data), nil
}

// ReadLines reads a file and returns its lines as a slice.
// Empty lines are preserved. The final newline is NOT included as an empty line.
func ReadLines(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("file.ReadLines(%q): %w", path, err)
    }
    defer f.Close()

    var lines []string
    scanner := bufio.NewScanner(f)
    scanner.Buffer(make([]byte, 1<<20), 1<<20) // 1MB max line length

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("file.ReadLines(%q): %w", path, err)
    }
    return lines, nil
}

// ReadBytes reads the entire contents of the file as a byte slice.
// Use this for binary files (images, compiled objects, etc.).
func ReadBytes(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("file.ReadBytes(%q): %w", path, err)
    }
    return data, nil
}

// ─── Writing ─────────────────────────────────────────────────────────────────

// WriteFile creates or truncates the file at path and writes content.
// Syscalls: open(O_WRONLY|O_CREATE|O_TRUNC, 0644) + write() + close()
func WriteFile(path, content string) error {
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        return fmt.Errorf("file.WriteFile(%q): %w", path, err)
    }
    return nil
}

// WriteBytes creates or truncates the file at path and writes raw bytes.
// Use for binary output (compiled object files, images, etc.).
func WriteBytes(path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("file.WriteBytes(%q): %w", path, err)
    }
    return nil
}

// AppendFile appends content to the file at path, creating it if needed.
// Syscalls: open(O_WRONLY|O_APPEND|O_CREATE, 0644) + write() + close()
func AppendFile(path, content string) error {
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("file.AppendFile(%q): %w", path, err)
    }
    defer f.Close()

    _, err = f.WriteString(content)
    if err != nil {
        return fmt.Errorf("file.AppendFile(%q): write: %w", path, err)
    }
    return nil
}

// WriteLines writes a slice of strings to a file, one per line.
func WriteLines(path string, lines []string) error {
    content := strings.Join(lines, "\n")
    if len(lines) > 0 {
        content += "\n"
    }
    return WriteFile(path, content)
}

// ─── Querying ─────────────────────────────────────────────────────────────────

// Exists returns true if the file or directory at path exists.
// Syscalls: stat()
func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

// IsFile returns true if path exists and is a regular file.
func IsFile(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.Mode().IsRegular()
}

// IsDir returns true if path exists and is a directory.
func IsDir(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}

// FileInfo holds metadata about a file.
type FileInfo struct {
    Path     string
    Name     string
    Size     int64
    IsDir    bool
    Modified time.Time
    Mode     os.FileMode
}

// Stat returns metadata about the file at path.
// Syscalls: stat()
func Stat(path string) (FileInfo, error) {
    info, err := os.Stat(path)
    if err != nil {
        return FileInfo{}, fmt.Errorf("file.Stat(%q): %w", path, err)
    }
    return FileInfo{
        Path:     path,
        Name:     info.Name(),
        Size:     info.Size(),
        IsDir:    info.IsDir(),
        Modified: info.ModTime(),
        Mode:     info.Mode(),
    }, nil
}

// Size returns the size in bytes of the file at path.
func Size(path string) (int64, error) {
    info, err := os.Stat(path)
    if err != nil {
        return 0, fmt.Errorf("file.Size(%q): %w", path, err)
    }
    return info.Size(), nil
}

// ─── Manipulation ─────────────────────────────────────────────────────────────

// Copy copies the file at src to dst, creating dst if needed.
func Copy(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("file.Copy: open src %q: %w", src, err)
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("file.Copy: create dst %q: %w", dst, err)
    }
    defer out.Close()

    if _, err := io.Copy(out, in); err != nil {
        return fmt.Errorf("file.Copy(%q → %q): %w", src, dst, err)
    }
    return nil
}

// Move renames/moves the file at src to dst.
// Syscalls: rename() (atomic on same filesystem)
func Move(src, dst string) error {
    if err := os.Rename(src, dst); err != nil {
        // rename() fails across filesystems — fall back to copy+delete
        if err := Copy(src, dst); err != nil {
            return fmt.Errorf("file.Move(%q → %q): %w", src, dst, err)
        }
        return os.Remove(src)
    }
    return nil
}

// Delete removes the file at path.
// Syscalls: unlink()
func Delete(path string) error {
    if err := os.Remove(path); err != nil {
        return fmt.Errorf("file.Delete(%q): %w", path, err)
    }
    return nil
}

// ─── Directories ─────────────────────────────────────────────────────────────

// MakeDir creates the directory at path.
// Syscalls: mkdir()
func MakeDir(path string) error {
    if err := os.Mkdir(path, 0755); err != nil {
        return fmt.Errorf("file.MakeDir(%q): %w", path, err)
    }
    return nil
}

// MakeDirAll creates the directory at path and all intermediate parents.
// Syscalls: mkdir() × depth
func MakeDirAll(path string) error {
    if err := os.MkdirAll(path, 0755); err != nil {
        return fmt.Errorf("file.MakeDirAll(%q): %w", path, err)
    }
    return nil
}

// ListDir returns the names of all entries in the directory at path.
// Syscalls: open(O_RDONLY) + getdents() + close()
func ListDir(path string) ([]string, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, fmt.Errorf("file.ListDir(%q): %w", path, err)
    }
    names := make([]string, len(entries))
    for i, entry := range entries {
        names[i] = entry.Name()
    }
    return names, nil
}

// ListDirInfo returns FileInfo for all entries in the directory at path.
func ListDirInfo(path string) ([]FileInfo, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, fmt.Errorf("file.ListDirInfo(%q): %w", path, err)
    }
    infos := make([]FileInfo, 0, len(entries))
    for _, entry := range entries {
        info, err := entry.Info()
        if err != nil { continue }
        infos = append(infos, FileInfo{
            Path:     filepath.Join(path, entry.Name()),
            Name:     entry.Name(),
            Size:     info.Size(),
            IsDir:    entry.IsDir(),
            Modified: info.ModTime(),
            Mode:     info.Mode(),
        })
    }
    return infos, nil
}

// Walk recursively visits all files and directories under root.
// callback is called for each file and directory.
func Walk(root string, callback func(path string, info FileInfo) error) error {
    return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil { return err }
        return callback(path, FileInfo{
            Path:     path,
            Name:     info.Name(),
            Size:     info.Size(),
            IsDir:    info.IsDir(),
            Modified: info.ModTime(),
            Mode:     info.Mode(),
        })
    })
}

// ─── Paths ────────────────────────────────────────────────────────────────────

// Join joins path elements with the OS separator.
func Join(parts ...string) string {
    return filepath.Join(parts...)
}

// Abs returns the absolute path of the given path.
func Abs(path string) (string, error) {
    return filepath.Abs(path)
}

// Dir returns the directory component of a path.
// Dir("/a/b/c.txt") == "/a/b"
func Dir(path string) string {
    return filepath.Dir(path)
}

// Base returns the final component of a path.
// Base("/a/b/c.txt") == "c.txt"
func Base(path string) string {
    return filepath.Base(path)
}

// Ext returns the file extension (with dot).
// Ext("main.astra") == ".astra"
func Ext(path string) string {
    return filepath.Ext(path)
}

// Stem returns the filename without extension.
// Stem("main.astra") == "main"
func Stem(path string) string {
    base := filepath.Base(path)
    ext  := filepath.Ext(base)
    return strings.TrimSuffix(base, ext)
}

// TempFile creates a temporary file and returns its path.
func TempFile(prefix string) (string, error) {
    f, err := os.CreateTemp("", prefix)
    if err != nil {
        return "", fmt.Errorf("file.TempFile: %w", err)
    }
    name := f.Name()
    f.Close()
    return name, nil
}
```

Now the corresponding Astra source declarations (the standard library interface Astra programs see):

```astra
// std/file.astra — Astra declarations for the file stdlib
// These map directly to the Go implementations above.

module file

// Read the entire file as a string
extern fn read_file(path: string) -> Result<string, Error>

// Read file as lines
extern fn read_lines(path: string) -> Result<List<string>, Error>

// Write content to file (creates or truncates)
extern fn write_file(path: string, content: string) -> Result<(), Error>

// Append to file
extern fn append_file(path: string, content: string) -> Result<(), Error>

// Check if path exists
extern fn exists(path: string) -> bool

// Get file size in bytes
extern fn size(path: string) -> Result<int, Error>

// List directory contents
extern fn list_dir(path: string) -> Result<List<string>, Error>

// Create directory (and parents)
extern fn make_dir_all(path: string) -> Result<(), Error>

// Path utilities
extern fn join(parts: List<string>) -> string
extern fn base(path: string) -> string
extern fn dir(path: string) -> string
extern fn ext(path: string) -> string
extern fn stem(path: string) -> string

// Usage example in Astra:
fn compile_all_files(dir: string) -> Result<(), Error> {
    let files = file.list_dir(dir)?
    for f in files {
        if file.ext(f) == ".astra" {
            let source = file.read_file(file.join([dir, f]))?
            compile(source)?
        }
    }
    return ok(())
}
```

The `?` operator is Astra's error propagation syntax (similar to Rust's `?`) — it unwraps `Result::Ok` or returns `Result::Err` early.

---

## Exercises

1. **Inode explorer**: Write a Go program that, given a file path as a command-line argument, uses `os.Stat()` and `syscall.Stat_t` to print: the inode number, number of hard links, file size, permissions (in octal), user ID, last modification time. Verify your output matches `ls -li` for the same file.

2. **Buffering benchmark**: Write a Go benchmark comparing three approaches to reading a 10MB text file line by line: (a) raw `os.File.Read` with a 1-byte buffer, (b) `bufio.Scanner`, (c) `os.ReadFile` then `strings.Split`. Measure throughput (lines/second) for each. Explain why the results differ.

3. **Pipe-based compiler pipeline**: Implement a Unix pipe in Go using `os.Pipe()`. Write a "producer" goroutine that writes Astra source code line by line to the pipe's write end. Write a "consumer" goroutine that reads from the pipe's read end and counts the total tokens. This simulates a streaming lexer pipeline.

4. **Directory tree printer**: Implement a `tree`-like command in Go using `filepath.Walk`. Print a directory tree with ASCII box-drawing characters (├── for non-last items, └── for last items, │ for indent), showing file sizes for regular files. Handle deeply nested directories correctly.

5. **mmap lexer**: Implement a fast Astra lexer that uses `syscall.Mmap` to map the source file into memory instead of calling `os.ReadFile`. Compare its performance to the `os.ReadFile`-based lexer on a 1MB source file. Measure the time for 100 repetitions.

6. **epoll simulation**: Since epoll is Linux-specific, use Go's `net` package to write an echo server that handles 100 simultaneous clients using a single goroutine with `net.Listener` and non-blocking reads (set a very short deadline). Measure how many messages per second the server handles compared to a goroutine-per-client server.

---

## Summary

| Concept | Key Point |
|---|---|
| File system | Organizes disk data into hierarchical namespace of files and directories |
| Inode | Fixed-size metadata record: permissions, size, timestamps, block pointers — no name |
| Directory entry | Name → inode number mapping; the "file name" lives in the directory, not the inode |
| Hard link | Multiple directory entries pointing to the same inode |
| Soft link | File whose content is a path string; one extra path resolution per dereference |
| VFS | Kernel abstraction layer; uniform open/read/write API across all file system types |
| File descriptor | Small integer; per-process index into open-file table; 0=stdin, 1=stdout, 2=stderr |
| open/read/write/close | The four system calls underlying all file I/O |
| Buffered I/O | User-space buffer accumulates writes; flushes in one large syscall |
| mmap | Maps file into virtual address space; access via pointer; zero-copy |
| Pipe | One-way byte stream; used for shell pipelines and inter-process communication |
| Non-blocking I/O | `O_NONBLOCK` flag; `read()` returns EAGAIN instead of blocking |
| epoll/kqueue | Register many fds once; get only active ones; O(1) per event |
| Event loop | Single thread handles thousands of connections using epoll |
| Go bufio | Buffered reader/writer wrapping os.File; use for line-by-line processing |
| Astra file stdlib | Complete file API mapping Astra functions to OS syscalls via Go |
