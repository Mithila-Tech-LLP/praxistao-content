# Chapter 66: Standard Library — File I/O

> "Everything is a file." — Unix philosophy

---

## Overview

File I/O is the gateway between your program and the world beyond RAM. When a program starts, its entire state lives in memory — volatile, gone the moment the process ends. Files provide persistence. They are how programs read configuration, store results, load data, communicate with other programs, and survive a reboot.

Almost every real program touches the filesystem. A web server reads static assets and config files. A compiler reads source code and writes binaries. A data pipeline reads CSV files and writes JSON. A game loads level files and saves progress. File I/O is not optional — it is the fundamental mechanism by which programs become useful.

This chapter builds Astra's `file` package from scratch. We will cover every operation a real program needs: reading entire files, reading line by line, writing and appending, checking existence, listing directories, creating paths, copying, moving, deleting, and working with temporary files. We will go deep on the underlying system calls, buffering strategy, error handling, and the complete Go implementation.

We will also establish Astra's error handling idiom for I/O: the `Result<T, E>` pattern. File operations fail. Disks fill up, files get deleted, permissions are wrong, network filesystems disconnect. A good stdlib makes the error path visible and explicit, not something you can accidentally ignore.

---

## What We're Building

```
stdlib/
  file/
    file.go           ← main file operations (~400 lines)
    path.go           ← path manipulation helpers (~100 lines)
    temp.go           ← temporary file management (~60 lines)

runtime/
  astra_file.h        ← C file operation helpers (optional low-level layer)
```

---

## Table of Contents

1. The Filesystem — what programs interact with
2. System Calls — what happens beneath the Go layer
3. Reading Files — read, read_lines, read_bytes
4. Writing Files — write, append, write_bytes
5. File Metadata and Existence — exists, info
6. Directory Operations — list, mkdir, remove, rename, copy
7. Path Operations — dirname, basename, extension, abs, join
8. Temporary Files
9. The Result<T, E> Pattern for I/O
10. Buffered vs Unbuffered I/O
11. Complete Go Implementation
12. Build Milestone
13. Exercises

---

## 1. The Filesystem

Modern operating systems present programs with an abstraction called the **virtual filesystem** (VFS). From a program's perspective, the filesystem is a tree of directories and files rooted at `/` (on Unix) or a drive letter (on Windows). Every file has:

- A **path** — the unique address within the tree
- **Contents** — a sequence of bytes of arbitrary length
- **Metadata** — name, size, type, timestamps, permissions

The VFS abstracts over many different physical storage technologies: spinning HDDs, SSDs, NFS network shares, memory-mapped virtual filesystems like `/proc`, and more. From Astra's point of view, we do not care which physical storage backs a path — we just read and write bytes.

```
FILESYSTEM TREE

/
├── Users/
│   └── aditya/
│       ├── .zshrc          (file, 2.1 KB)
│       ├── code/
│       │   └── tasks/
│       │       ├── tasks.as  (file, 8.4 KB)
│       │       └── build/
│       │           └── tasks (file, 1.2 MB — compiled binary)
│       └── Documents/
│           └── notes.txt   (file, 512 B)
├── usr/
│   └── local/
│       └── bin/
│           └── astrac      (file, 4.8 MB — the compiler)
└── tmp/
    └── astra_abc123        (temp file, 0 B)
```

---

## 2. System Calls — What Happens Beneath Go

When Go (and by extension Astra) opens a file, it eventually calls into the operating system kernel via **system calls**. On Linux/macOS:

```
open(2)    — open/create a file, returns a file descriptor (int)
read(2)    — read bytes from a file descriptor into a buffer
write(2)   — write bytes from a buffer to a file descriptor
close(2)   — close a file descriptor and flush buffers
stat(2)    — get file metadata (size, timestamps, type)
lstat(2)   — like stat but does not follow symlinks
unlink(2)  — delete a file (remove its directory entry)
rename(2)  — atomically rename/move a file
mkdir(2)   — create a directory
rmdir(2)   — remove an empty directory
getdents(2)— list directory contents
```

Go's `os` package wraps all of these. We do not call them directly from our Go stdlib code — we let Go's runtime handle the OS interaction. But it is important to understand what is happening under the hood so we can reason about performance and error conditions.

### The File Descriptor

When you open a file, the OS returns an **integer file descriptor** (fd). Think of it as a handle into the OS's open file table:

```
PROCESS                    OS KERNEL
┌──────────────────┐       ┌─────────────────────────────────┐
│ fd table         │       │ Open File Table                  │
│ fd 0 ──────────►─┼──────►│ stdin  (terminal read end)       │
│ fd 1 ──────────►─┼──────►│ stdout (terminal write end)      │
│ fd 2 ──────────►─┼──────►│ stderr (terminal error write)    │
│ fd 3 ──────────►─┼──────►│ notes.txt (offset: 0, mode: r)  │
│ fd 4 ──────────►─┼──────►│ output.txt (offset: 0, mode: w) │
└──────────────────┘       └─────────────────────────────────┘
                                          │
                                          ▼ inode lookup
                           ┌─────────────────────────────────┐
                           │ Inode Table (disk)               │
                           │ inode 48271: notes.txt           │
                           │   size: 512 bytes                │
                           │   blocks: [0x1a4f, 0x0000, ...]  │
                           │   mtime: 2024-01-15 14:30:00     │
                           └─────────────────────────────────┘
```

Go abstracts file descriptors as `*os.File`. When the file is closed, the fd is returned to the pool and can be reused by the next `open` call.

---

## 3. Reading Files

### 3.1 Read Entire File

The most common operation: read the entire contents of a file into a string.

```astra
import file

let content = file.read("data.txt")
// content is a string containing all bytes of data.txt, decoded as UTF-8
```

Internally, Go uses `os.ReadFile` which reads the file in a single call (after getting the file size from `stat`), allocating exactly the right buffer. This is O(n) time and O(n) memory.

### 3.2 Read Lines

For text processing, reading line by line is often more ergonomic:

```astra
let lines = file.read_lines("data.csv")
for line in lines {
    let parts = line.split(",")
    print(parts[0])
}
```

`read_lines` returns a `List<string>` where each element is one line, with the trailing newline stripped. It handles both `\n` (Unix) and `\r\n` (Windows) line endings correctly.

### 3.3 Read Raw Bytes

For binary files (images, audio, compiled binaries), we need raw bytes:

```astra
let bytes = file.read_bytes("image.png")
// bytes is List<byte>, i.e. List<int> where each value is 0–255
let magic = bytes.slice(0, 4)  // PNG magic: [137, 80, 78, 71]
```

---

## 4. Writing Files

### 4.1 Write (Create or Overwrite)

```astra
file.write("output.txt", "Hello, World!\n")
// Creates the file if it does not exist, truncates it if it does.
// Permissions: 0644 (owner rw, group r, other r)
```

### 4.2 Append

```astra
file.append("server.log", "[2024-01-15 14:30:00] Request received\n")
// Opens the file in O_APPEND mode. Writes are atomic at the OS level
// for writes smaller than PIPE_BUF (typically 4096 bytes).
```

The difference between write and append at the syscall level:

```c
// write mode: O_WRONLY | O_CREATE | O_TRUNC
// append mode: O_WRONLY | O_CREATE | O_APPEND
```

### 4.3 Write Raw Bytes

```astra
let bytes = [137, 80, 78, 71, 13, 10, 26, 10]  // PNG magic bytes
file.write_bytes("header.bin", bytes)
```

---

## 5. File Metadata and Existence

```astra
// Check existence before operating
if file.exists("config.toml") {
    let info = file.info("config.toml")
    print(info.name)            // "config.toml"
    print(info.size.to_string() + " bytes")  // "4096 bytes"
    print(info.is_dir.to_string())            // "false"
    print(info.modified)        // "2024-01-15 14:30:00"
}
```

The `FileInfo` struct:

```astra
struct FileInfo {
    name:     string   // file name (not full path)
    size:     int      // size in bytes
    is_dir:   bool     // true if this is a directory
    mode:     int      // Unix permission bits (e.g., 0644)
    modified: string   // last modification time, formatted as ISO 8601
}
```

**Important**: checking `exists` then acting on the result is a **TOCTOU** (time-of-check-time-of-use) race condition in concurrent programs. Between the `exists` check and the `read`, another process could delete the file. The correct pattern for concurrent safety is to attempt the operation and handle the error:

```astra
match file.read("config.toml") {
    Ok(content) => // use content
    Err(e)      => if e.contains("not found") { use_defaults() } else { panic(e) }
}
```

---

## 6. Directory Operations

```astra
// List directory contents
let entries = file.list("./src")
for entry in entries {
    let prefix = if entry.is_dir { "dir:  " } else { "file: " }
    print(prefix + entry.name + " (" + entry.size.to_string() + " bytes)")
}

// Create directory (recursive — creates parent dirs as needed)
file.mkdir("output/logs/2024")

// Delete file
file.remove("temp.txt")

// Delete directory (must be empty, or use remove_dir_all for recursive delete)
file.remove_dir("empty_dir")
file.remove_dir_all("old_build")   // recursive delete, USE WITH CARE

// Rename or move
file.rename("old_name.txt", "new_name.txt")
file.rename("./file.txt", "./subdir/file.txt")  // move to subdirectory

// Copy a file (does not copy directory trees)
file.copy("src/template.html", "dist/index.html")
```

---

## 7. Path Operations

Path manipulation is platform-sensitive: Unix uses `/`, Windows uses `\`. Astra's path functions always use `/` internally and handle the OS translation:

```astra
// Decompose paths
let full = "/usr/local/bin/astrac"
print(file.dirname(full))        // "/usr/local/bin"
print(file.basename(full))       // "astrac"
print(file.extension("main.as")) // "as"
print(file.extension("Makefile"))// "" (no extension)

// Combine paths safely (no double slashes)
let path = file.join("/usr/local", "bin", "astrac")
// → "/usr/local/bin/astrac"

// Resolve relative paths
let abs = file.abs("./src/../src/main.as")
// → "/Users/aditya/code/myproject/src/main.as"

// Special directories
let home = file.home_dir()         // "/Users/aditya"
let cwd  = file.working_dir()      // current working directory
let tmp  = file.temp_dir()         // "/tmp" on Unix, "C:\Temp" on Windows
```

---

## 8. Temporary Files

Temporary files are created, used briefly, and then should be cleaned up. Astra's `file.temp_file` creates a file with a unique name in the system's temp directory:

```astra
// Create a temp file with a name matching the pattern
// * is replaced with random characters to ensure uniqueness
let tmp = file.temp_file("astra_*")
defer file.remove(tmp.path)    // always clean up, even on error

tmp.write("processing data...")
// ... do work with tmp.path ...
// When the function returns, defer fires and tmp is deleted
```

The `defer` keyword (like Go's `defer`) ensures cleanup even if a panic occurs. It registers a function call to execute when the current scope exits.

---

## 9. The Result<T, E> Pattern for File I/O

Every file operation can fail. The Go stdlib signals failure via multiple return values `(T, error)`. Astra uses `Result<T, E>`:

```astra
// Result type definition (built into the language)
enum Result<T, E> {
    Ok(T)
    Err(E)
}
```

All file operations return `Result<T, string>` where the string is the error message:

```astra
fn read_config(path: string) -> Result<Config, string> {
    // file.read returns Result<string, string>
    let content = file.read(path)?     // ? propagates Err upward
    let cfg = json.unmarshal<Config>(content)?
    return Ok(cfg)
}

fn main() {
    match read_config("./config.json") {
        Ok(cfg)  => start_server(cfg)
        Err(msg) => {
            print("Failed to load config: " + msg)
            print("Using defaults.")
            start_server(Config.default())
        }
    }
}
```

The `?` operator (question mark) is syntactic sugar for:

```astra
let content = match file.read(path) {
    Ok(v)  => v
    Err(e) => return Err(e)
}
```

It propagates the error to the caller. Any function that uses `?` must return a `Result`.

### Error Categories

```
FileNotFound     — the path does not exist
PermissionDenied — the process lacks read/write permission
IsADirectory     — tried to read a dir as a file (use list instead)
NotADirectory    — tried to list a file as a dir
DiskFull         — no space left on device
TooManyOpenFiles — OS file descriptor limit reached
InvalidPath      — path contains illegal characters
```

In Astra, these are all represented as strings in the `Err` variant. A future version might use a typed error enum.

---

## 10. Buffered vs Unbuffered I/O

Understanding when I/O is buffered vs unbuffered is critical for performance.

### Unbuffered I/O

Every read or write goes directly to the OS via a system call. Fast for large sequential transfers, but each `write("x")` call incurs the overhead of a kernel context switch:

```
userspace write("x") → syscall write(fd, "x", 1) → kernel → disk
                                                    ↑
                                             ~1-10 μs per call
```

### Buffered I/O

Go's `bufio.Writer` accumulates writes in a memory buffer (default 4 KB). The system call only happens when the buffer is full or when `Flush()` is explicitly called:

```
write("x") → buffer [x......]   (no syscall yet)
write("y") → buffer [xy.....]   (no syscall yet)
write("z") → buffer [xyz....]   (no syscall yet)
...4096 chars later...
flush()    → syscall write(fd, buffer, 4096) → kernel → disk
             ↑
             one syscall for 4096 bytes — much more efficient
```

Astra's `file.write` and `file.append` use buffered I/O internally and flush automatically when the write completes. For high-throughput writing (e.g., writing millions of log lines), use the `file.Writer` streaming API:

```astra
import file

fn write_million_lines(path: string) {
    let w = file.Writer.new(path)
    defer w.close()   // flushes and closes when done

    for i in 0..1_000_000 {
        w.write_line("Line " + i.to_string())
        // Internal buffer accumulates lines; syscalls happen every ~4KB
    }
    // w.close() flushes remaining buffer and closes the fd
}
```

Performance comparison for writing 1 million lines:
```
file.write (build string first, then write): ~800ms, 80MB memory for string
file.Writer (streaming buffered):             ~120ms, ~4KB buffer overhead
```

---

## 11. Complete Go Implementation

```go
// stdlib/file/file.go

package astra_file

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// ─── Types ────────────────────────────────────────────────────────────────────

// FileInfo holds metadata about a filesystem entry.
type FileInfo struct {
    Name     string // base name (not full path)
    Size     int64  // size in bytes (0 for directories)
    IsDir    bool   // true if this is a directory
    Mode     int    // Unix permission bits (e.g. 0644)
    Modified string // last modification time, ISO 8601
}

// AstraFile wraps an os.File with a buffered writer for efficient writes.
type AstraFile struct {
    path   string
    f      *os.File
    writer *bufio.Writer
    reader *bufio.Reader
}

// ─── Reading ──────────────────────────────────────────────────────────────────

// Read reads the entire contents of the file at path as a UTF-8 string.
// Returns Err if the file does not exist, cannot be opened, or is not valid UTF-8.
func Read(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", formatError("read", path, err)
    }
    return string(data), nil
}

// ReadLines reads the file and returns each line as a separate string.
// Trailing newlines (\n and \r\n) are stripped from each line.
// Empty lines are preserved as empty strings.
func ReadLines(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, formatError("read_lines", path, err)
    }
    defer f.Close()

    var lines []string
    scanner := bufio.NewScanner(f)
    // Increase scanner buffer for long lines (default 64KB may be too small)
    scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // up to 10MB per line
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, formatError("read_lines", path, err)
    }
    return lines, nil
}

// ReadBytes reads the file as raw bytes.
// Returns a []byte slice (represented as List<byte> in Astra).
func ReadBytes(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, formatError("read_bytes", path, err)
    }
    return data, nil
}

// ─── Writing ──────────────────────────────────────────────────────────────────

// Write creates or overwrites the file at path with the given content.
// Creates parent directories if they do not exist.
// Permissions: 0644 (owner rw, group r, other r).
func Write(path, content string) error {
    if err := ensureParentDir(path); err != nil {
        return formatError("write", path, err)
    }
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        return formatError("write", path, err)
    }
    return nil
}

// Append opens the file for appending (creates if not exists) and writes content.
// The write is atomic for content smaller than PIPE_BUF (typically 4096 bytes).
func Append(path, content string) error {
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return formatError("append", path, err)
    }
    defer f.Close()

    _, err = f.WriteString(content)
    if err != nil {
        return formatError("append", path, err)
    }
    return nil
}

// WriteBytes writes raw bytes to the file at path (create or overwrite).
func WriteBytes(path string, data []byte) error {
    if err := ensureParentDir(path); err != nil {
        return formatError("write_bytes", path, err)
    }
    if err := os.WriteFile(path, data, 0644); err != nil {
        return formatError("write_bytes", path, err)
    }
    return nil
}

// ─── Existence and Metadata ───────────────────────────────────────────────────

// Exists returns true if the path exists (file or directory).
// Does not follow symlinks.
func Exists(path string) bool {
    _, err := os.Lstat(path)
    return err == nil
}

// Info returns metadata for the file or directory at path.
func Info(path string) (FileInfo, error) {
    stat, err := os.Stat(path) // follows symlinks
    if err != nil {
        return FileInfo{}, formatError("info", path, err)
    }
    return FileInfo{
        Name:     stat.Name(),
        Size:     stat.Size(),
        IsDir:    stat.IsDir(),
        Mode:     int(stat.Mode().Perm()),
        Modified: stat.ModTime().Format(time.RFC3339),
    }, nil
}

// IsFile returns true if path exists and is a regular file (not a directory).
func IsFile(path string) bool {
    stat, err := os.Stat(path)
    return err == nil && !stat.IsDir()
}

// IsDir returns true if path exists and is a directory.
func IsDir(path string) bool {
    stat, err := os.Stat(path)
    return err == nil && stat.IsDir()
}

// Size returns the size of the file in bytes.
func Size(path string) (int64, error) {
    stat, err := os.Stat(path)
    if err != nil {
        return 0, formatError("size", path, err)
    }
    return stat.Size(), nil
}

// ─── Directory Operations ─────────────────────────────────────────────────────

// List returns the entries in the directory at path.
// Does not recurse into subdirectories.
// Entries are sorted by name.
func List(path string) ([]FileInfo, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, formatError("list", path, err)
    }

    result := make([]FileInfo, 0, len(entries))
    for _, entry := range entries {
        info, err := entry.Info()
        if err != nil {
            continue // skip unreadable entries
        }
        result = append(result, FileInfo{
            Name:     entry.Name(),
            Size:     info.Size(),
            IsDir:    entry.IsDir(),
            Mode:     int(info.Mode().Perm()),
            Modified: info.ModTime().Format(time.RFC3339),
        })
    }
    return result, nil
}

// ListAll returns all files recursively under path.
// Returns full paths relative to path.
func ListAll(rootPath string) ([]FileInfo, error) {
    var result []FileInfo

    err := filepath.Walk(rootPath, func(p string, info os.FileInfo, err error) error {
        if err != nil {
            return nil // skip unreadable entries
        }
        rel, _ := filepath.Rel(rootPath, p)
        if rel == "." {
            return nil // skip the root itself
        }
        result = append(result, FileInfo{
            Name:     rel,
            Size:     info.Size(),
            IsDir:    info.IsDir(),
            Mode:     int(info.Mode().Perm()),
            Modified: info.ModTime().Format(time.RFC3339),
        })
        return nil
    })
    if err != nil {
        return nil, formatError("list_all", rootPath, err)
    }
    return result, nil
}

// Mkdir creates the directory at path, including all parent directories.
// Equivalent to `mkdir -p`. No error if directory already exists.
func Mkdir(path string) error {
    if err := os.MkdirAll(path, 0755); err != nil {
        return formatError("mkdir", path, err)
    }
    return nil
}

// Remove deletes the file at path. Returns an error if it is a directory.
func Remove(path string) error {
    if err := os.Remove(path); err != nil {
        return formatError("remove", path, err)
    }
    return nil
}

// RemoveDir removes the directory at path. Fails if not empty.
func RemoveDir(path string) error {
    if err := os.Remove(path); err != nil {
        return formatError("remove_dir", path, err)
    }
    return nil
}

// RemoveDirAll recursively deletes path and everything inside it.
// USE WITH EXTREME CARE. There is no undo.
func RemoveDirAll(path string) error {
    if err := os.RemoveAll(path); err != nil {
        return formatError("remove_dir_all", path, err)
    }
    return nil
}

// Rename renames (or moves) oldPath to newPath.
// This is an atomic operation on the same filesystem.
// If newPath already exists, it is replaced (on Unix).
func Rename(oldPath, newPath string) error {
    if err := ensureParentDir(newPath); err != nil {
        return formatError("rename", newPath, err)
    }
    if err := os.Rename(oldPath, newPath); err != nil {
        return formatError("rename", oldPath, err)
    }
    return nil
}

// Copy copies the file at src to dst (create or overwrite).
// Copies file contents and permissions. Does not copy directories.
func Copy(src, dst string) error {
    // Open source
    in, err := os.Open(src)
    if err != nil {
        return formatError("copy", src, err)
    }
    defer in.Close()

    // Get source permissions
    srcInfo, err := in.Stat()
    if err != nil {
        return formatError("copy", src, err)
    }

    // Ensure destination parent exists
    if err := ensureParentDir(dst); err != nil {
        return formatError("copy", dst, err)
    }

    // Create destination
    out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
    if err != nil {
        return formatError("copy", dst, err)
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()

    // Copy with buffered I/O (io.Copy uses a 32KB internal buffer)
    _, err = io.Copy(out, in)
    if err != nil {
        return formatError("copy", src, err)
    }
    return nil
}

// ─── Path Operations ─────────────────────────────────────────────────────────

// Dirname returns the directory component of path.
func Dirname(path string) string {
    return filepath.Dir(path)
}

// Basename returns the base name component of path (last element).
func Basename(path string) string {
    return filepath.Base(path)
}

// Extension returns the file extension (without the dot).
// Returns "" for files without an extension.
func Extension(path string) string {
    ext := filepath.Ext(path)
    if len(ext) > 0 && ext[0] == '.' {
        return ext[1:]
    }
    return ext
}

// Join joins path elements with the OS path separator, cleaning the result.
func Join(parts ...string) string {
    return filepath.Join(parts...)
}

// Abs returns the absolute path for path. Resolves . and .. components.
// Returns an error if the path cannot be resolved.
func Abs(path string) (string, error) {
    abs, err := filepath.Abs(path)
    if err != nil {
        return "", fmt.Errorf("file.abs: %v", err)
    }
    return abs, nil
}

// HomeDir returns the current user's home directory.
func HomeDir() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("file.home_dir: %v", err)
    }
    return home, nil
}

// WorkingDir returns the current working directory.
func WorkingDir() (string, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("file.working_dir: %v", err)
    }
    return cwd, nil
}

// TempDir returns the directory for temporary files.
func TempDir() string {
    return os.TempDir()
}

// ─── Temporary Files ─────────────────────────────────────────────────────────

// TempFileHandle wraps a temporary file with its path.
type TempFileHandle struct {
    Path string
    f    *os.File
}

// Write writes content to the temporary file.
func (t *TempFileHandle) Write(content string) error {
    _, err := t.f.WriteString(content)
    return err
}

// Close closes the file without deleting it. Use file.remove(tmp.path) to delete.
func (t *TempFileHandle) Close() error {
    return t.f.Close()
}

// TempFile creates a temporary file with a name matching the pattern.
// The '*' in the pattern is replaced with a random string.
// Example: TempFile("astra_*") might create "/tmp/astra_382719234"
func TempFile(pattern string) (*TempFileHandle, error) {
    f, err := os.CreateTemp("", pattern)
    if err != nil {
        return nil, fmt.Errorf("file.temp_file: %v", err)
    }
    return &TempFileHandle{Path: f.Name(), f: f}, nil
}

// ─── Streaming Writer ─────────────────────────────────────────────────────────

// Writer is a buffered file writer for high-throughput output.
type Writer struct {
    f   *os.File
    buf *bufio.Writer
}

// NewWriter opens path for writing (creates or truncates) with a 4KB buffer.
func NewWriter(path string) (*Writer, error) {
    if err := ensureParentDir(path); err != nil {
        return nil, formatError("Writer.new", path, err)
    }
    f, err := os.Create(path)
    if err != nil {
        return nil, formatError("Writer.new", path, err)
    }
    return &Writer{
        f:   f,
        buf: bufio.NewWriterSize(f, 4*1024), // 4KB buffer
    }, nil
}

// Write writes a string to the buffer (flushed to disk when buffer fills or Close is called).
func (w *Writer) Write(s string) error {
    _, err := w.buf.WriteString(s)
    return err
}

// WriteLine writes a string followed by a newline.
func (w *Writer) WriteLine(s string) error {
    if err := w.Write(s); err != nil {
        return err
    }
    return w.buf.WriteByte('\n')
}

// Flush flushes the buffer to disk without closing the file.
func (w *Writer) Flush() error {
    return w.buf.Flush()
}

// Close flushes the buffer and closes the file.
func (w *Writer) Close() error {
    if err := w.buf.Flush(); err != nil {
        w.f.Close()
        return err
    }
    return w.f.Close()
}

// ─── Streaming Reader ─────────────────────────────────────────────────────────

// LineReader reads a file line by line without loading it all into memory.
// Suitable for very large files (gigabytes).
type LineReader struct {
    f       *os.File
    scanner *bufio.Scanner
    err     error
}

// NewLineReader opens path for line-by-line reading.
func NewLineReader(path string) (*LineReader, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, formatError("LineReader.new", path, err)
    }
    scanner := bufio.NewScanner(f)
    scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)
    return &LineReader{f: f, scanner: scanner}, nil
}

// Next advances to the next line. Returns false when done or on error.
func (r *LineReader) Next() bool {
    return r.scanner.Scan()
}

// Line returns the current line text (without newline).
func (r *LineReader) Line() string {
    return r.scanner.Text()
}

// Err returns any scanner error (not io.EOF).
func (r *LineReader) Err() error {
    return r.scanner.Err()
}

// Close closes the underlying file.
func (r *LineReader) Close() error {
    return r.f.Close()
}

// ─── Internal Helpers ─────────────────────────────────────────────────────────

// formatError creates a user-friendly error message.
func formatError(op, path string, err error) error {
    if os.IsNotExist(err) {
        return fmt.Errorf("file.%s: file not found: %s", op, path)
    }
    if os.IsPermission(err) {
        return fmt.Errorf("file.%s: permission denied: %s", op, path)
    }
    if os.IsExist(err) {
        return fmt.Errorf("file.%s: already exists: %s", op, path)
    }
    // Check for disk full
    if strings.Contains(err.Error(), "no space left") {
        return fmt.Errorf("file.%s: disk full: %s", op, path)
    }
    return fmt.Errorf("file.%s: %s: %v", op, path, err)
}

// ensureParentDir creates the parent directory of path if it does not exist.
func ensureParentDir(path string) error {
    dir := filepath.Dir(path)
    if dir == "." || dir == "/" {
        return nil
    }
    return os.MkdirAll(dir, 0755)
}
```

---

## 12. Real-World Example: JSON Config Processing Pipeline

Here is a complete Astra program that demonstrates all the file operations working together. It reads a JSON configuration file, processes data files listed in the config, and writes a summary report.

```astra
import file
import json
import time
import string
import math

// Config file: config.json
// {
//   "input_dir": "./data",
//   "output_file": "./report.txt",
//   "max_lines": 1000
// }

struct Config {
    input_dir:   string
    output_file: string
    max_lines:   int
}

struct LineStats {
    file_name:    string
    total_lines:  int
    empty_lines:  int
    longest_line: int
    total_chars:  int
}

fn load_config(path: string) -> Result<Config, string> {
    let content = file.read(path)?
    let cfg = json.unmarshal<Config>(content)?
    return Ok(cfg)
}

fn analyze_file(path: string, max_lines: int) -> Result<LineStats, string> {
    let lines = file.read_lines(path)?

    let stats = LineStats {
        file_name:    file.basename(path),
        total_lines:  0,
        empty_lines:  0,
        longest_line: 0,
        total_chars:  0
    }

    let limit = math.min_int(lines.len(), max_lines)

    for i in 0..limit {
        let line = lines[i]
        stats.total_lines = stats.total_lines + 1
        stats.total_chars = stats.total_chars + line.length()
        if line.is_empty() {
            stats.empty_lines = stats.empty_lines + 1
        }
        if line.length() > stats.longest_line {
            stats.longest_line = line.length()
        }
    }

    return Ok(stats)
}

fn write_report(cfg: Config, all_stats: List<LineStats>) -> Result<void, string> {
    let sb = string.Builder.new()

    sb.append("=== File Analysis Report ===\n")
    sb.append("Generated: ")
    sb.append(time.format(time.now(), "2006-01-02 15:04:05"))
    sb.append("\n\n")

    for stats in all_stats {
        sb.append("File: ")
        sb.append(stats.file_name)
        sb.append("\n")
        sb.append("  Total lines:  ")
        sb.append(stats.total_lines.to_string())
        sb.append("\n")
        sb.append("  Empty lines:  ")
        sb.append(stats.empty_lines.to_string())
        sb.append("\n")
        sb.append("  Longest line: ")
        sb.append(stats.longest_line.to_string())
        sb.append(" chars\n")
        sb.append("  Total chars:  ")
        sb.append(stats.total_chars.to_string())
        sb.append("\n\n")
    }

    let report = sb.build()
    file.write(cfg.output_file, report)?
    return Ok(void)
}

fn main() {
    let start = time.now()

    // Load configuration
    let cfg = match load_config("config.json") {
        Ok(c)  => c
        Err(e) => {
            print("Error loading config: " + e)
            return
        }
    }

    print("Input directory: " + cfg.input_dir)
    print("Output file:     " + cfg.output_file)
    print("Max lines:       " + cfg.max_lines.to_string())

    // Check that input directory exists
    if !file.exists(cfg.input_dir) {
        print("Input directory does not exist: " + cfg.input_dir)
        return
    }

    // List all .txt files in the input directory
    let entries = match file.list(cfg.input_dir) {
        Ok(e)  => e
        Err(e) => {
            print("Error listing directory: " + e)
            return
        }
    }

    let txt_files: List<string> = List.new()
    for entry in entries {
        if !entry.is_dir && file.extension(entry.name) == "txt" {
            txt_files.push(file.join(cfg.input_dir, entry.name))
        }
    }

    print("Found " + txt_files.len().to_string() + " .txt files")

    // Analyze each file
    let all_stats: List<LineStats> = List.new()
    for path in txt_files {
        print("Analyzing: " + path)
        match analyze_file(path, cfg.max_lines) {
            Ok(stats) => all_stats.push(stats)
            Err(e)    => print("  Warning: " + e)
        }
    }

    // Write report
    match write_report(cfg, all_stats) {
        Ok(_)  => print("\nReport written to: " + cfg.output_file)
        Err(e) => print("Error writing report: " + e)
    }

    let elapsed = time.since(start)
    print("Total time: " + elapsed.as_milliseconds().to_string() + "ms")
}
```

Running it:

```bash
$ astrac build process.as -o process
$ echo '{"input_dir":"./data","output_file":"./report.txt","max_lines":1000}' > config.json
$ mkdir -p data
$ echo -e "Hello world\nLine 2\n\nLine 4" > data/sample.txt
$ echo -e "One\nTwo\nThree" > data/numbers.txt
$ ./process
Input directory: ./data
Output file:     ./report.txt
Max lines:       1000
Found 2 .txt files
Analyzing: ./data/sample.txt
Analyzing: ./data/numbers.txt

Report written to: ./report.txt
Total time: 3ms

$ cat report.txt
=== File Analysis Report ===
Generated: 2024-01-15 14:30:00

File: sample.txt
  Total lines:  4
  Empty lines:  1
  Longest line: 11 chars
  Total chars:  23

File: numbers.txt
  Total lines:  3
  Empty lines:  0
  Longest line: 5 chars
  Total chars:  11
```

---

## 🔨 Astra Build Milestone

Complete `stdlib/file/file.go` (~400 lines). The file stdlib is done when:

| Capability | Function | Status |
|-----------|----------|--------|
| Read entire file | `file.read(path)` | Complete |
| Read lines | `file.read_lines(path)` | Complete |
| Read bytes | `file.read_bytes(path)` | Complete |
| Write file | `file.write(path, content)` | Complete |
| Append to file | `file.append(path, content)` | Complete |
| Write bytes | `file.write_bytes(path, bytes)` | Complete |
| Check existence | `file.exists(path)` | Complete |
| File metadata | `file.info(path)` | Complete |
| List directory | `file.list(path)` | Complete |
| Create directory | `file.mkdir(path)` | Complete |
| Delete file | `file.remove(path)` | Complete |
| Delete directory | `file.remove_dir(path)` | Complete |
| Rename/move | `file.rename(old, new)` | Complete |
| Copy file | `file.copy(src, dst)` | Complete |
| Path dirname | `file.dirname(path)` | Complete |
| Path basename | `file.basename(path)` | Complete |
| Path extension | `file.extension(path)` | Complete |
| Path join | `file.join(parts...)` | Complete |
| Absolute path | `file.abs(path)` | Complete |
| Home directory | `file.home_dir()` | Complete |
| Working directory | `file.working_dir()` | Complete |
| Temp file | `file.temp_file(pattern)` | Complete |
| Streaming writer | `file.Writer` | Complete |
| Streaming reader | `file.LineReader` | Complete |

Test run:
```bash
$ astrac build tests/file_test.as -o file_test
$ ./file_test
file.read: PASSED
file.write: PASSED
file.append: PASSED
file.read_lines: PASSED
file.exists: PASSED
file.info: PASSED
file.list: PASSED
file.mkdir: PASSED
file.copy: PASSED
file.rename: PASSED
file.remove: PASSED
path.dirname: PASSED
path.basename: PASSED
path.extension: PASSED
All file I/O tests passed!
```

---

## 13. Exercises

**Exercise 1 — File Search**
Implement `file.find(root: string, pattern: string) -> List<string>` that recursively searches for files matching a glob pattern (e.g., `"*.as"`, `"test_*.txt"`). Use `filepath.Match` from Go's standard library for the pattern matching.

**Exercise 2 — File Watcher (Polling)**
Implement a simple file watcher that calls a callback when a file changes:
```astra
file.watch("config.json", fn(path: string) {
    print("Config changed, reloading...")
    reload_config(path)
})
```
Poll every 500ms using `time.sleep` and compare the last modification time from `file.info`. This is not as efficient as OS-level watching (inotify/kqueue) but is simpler to implement.

**Exercise 3 — CSV Reader**
Using `file.read_lines`, implement a CSV parser in Astra:
```astra
fn read_csv(path: string) -> List<List<string>> {
    let lines = file.read_lines(path)
    // split each line by comma, handling quoted fields
    // "Alice, \"the great\"", 25 → ["Alice, \"the great\"", "25"]
}
```
Handle quoted fields correctly (a comma inside quotes is not a delimiter).

**Exercise 4 — Atomic Write**
Implement `file.write_atomic(path, content)` that writes to a temporary file first, then renames it over the destination. This prevents data corruption if the process crashes mid-write:
```
1. Write to path + ".tmp"
2. Rename path + ".tmp" to path (atomic on same filesystem)
```

**Exercise 5 — File Size Formatting**
Implement a function `format_size(bytes: int) -> string` that formats a byte count in human-readable form:
- < 1024 → "512 B"
- < 1024 * 1024 → "4.2 KB"
- < 1024^3 → "1.3 MB"
- else → "2.1 GB"
Then print a directory tree with formatted sizes using `file.list_all`.

**Exercise 6 — Safe Delete (Trash)**
Implement `file.trash(path)` that moves a file to a `~/.astra_trash/` directory instead of deleting it permanently. Include a `file.trash_restore(name)` function and a `file.trash_empty()` to permanently delete all trashed files.

---

## Summary

| Topic | Key Takeaway |
|-------|-------------|
| Everything is a file | Unix's virtual filesystem presents one unified tree for all storage |
| System calls | `open`, `read`, `write`, `close` are the 4 fundamental I/O syscalls |
| Read variants | `file.read` (full file), `file.read_lines` (line array), `file.read_bytes` (raw) |
| Write modes | `file.write` (overwrite), `file.append` (add to end) — determined by `O_TRUNC` vs `O_APPEND` |
| Buffered I/O | `bufio.Writer` batches writes into 4KB chunks — 10x faster for many small writes |
| Result pattern | All file ops return `Result<T, string>`. Use `?` to propagate errors, `match` to handle them |
| TOCTOU race | Never check `exists` and then act — attempt the operation and handle the error |
| Streaming | Use `file.Writer` and `file.LineReader` for files too large to load into memory |
| Path operations | `dirname`, `basename`, `extension`, `join`, `abs` are your toolkit for path manipulation |
| Atomic write | Write to temp file, then rename — the only safe way to update a file without data loss |
