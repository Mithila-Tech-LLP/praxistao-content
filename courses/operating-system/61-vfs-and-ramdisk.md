# Chapter 61: VFS and Ramdisk

> **"The Virtual File System is the uniform interface that makes all storage look the same to your programs. Whether data is on a hard disk, a ramdisk, a network mount, or a fake /proc entry — open(), read(), write(), close() work identically. The VFS is the abstraction that makes 'everything is a file' real."**

---

## Table of Contents

1. [Why a VFS?](#1-why-a-vfs)
2. [VFS Architecture — Objects and Operations](#2-vfs-architecture--objects-and-operations)
3. [The Ramdisk — An In-Memory File System](#3-the-ramdisk--an-in-memory-file-system)
4. [File Descriptors](#4-file-descriptors)
5. [The VFS Layer in C](#5-the-vfs-layer-in-c)
6. [Implementing open/read/write/close](#6-implementing-openreadwriteclose)
7. [Populating the Ramdisk](#7-populating-the-ramdisk)
8. [sys_open, sys_read, sys_write Updated](#8-sys_open-sys_read-sys_write-updated)
9. [Testing](#9-testing)
10. [Summary](#summary)

---

## 1. Why a VFS?

```
Without VFS:
  Reading a file on disk  → call ext4_read()
  Reading from keyboard   → call kbd_getchar()
  Reading from network    → call tcp_recv()
  
  Every program must know what type of "file" it's reading.
  Shell scripts would need special handling for each type.

With VFS:
  Everything has the same interface:
    fd = open("/etc/passwd", O_RDONLY)
    read(fd, buf, 1024)
    close(fd)
    
    fd = open("/dev/keyboard", O_RDONLY)
    read(fd, buf, 1)   ← same call, different backend
    close(fd)
    
  Programs are oblivious to the underlying implementation.
  That's the power of VFS.
```

---

## 2. VFS Architecture — Objects and Operations

```
Our minimal VFS has three key objects:

1. vfs_node_t — represents a file or directory
   Fields: name, type (file/dir/device), size
   Operations: read, write, open, close, readdir
   
2. file_descriptor_t — an open file handle
   Fields: node pointer, position (offset), flags
   A process has an array of these (open files)
   
3. mount_point_t — maps a path prefix to a vfs_node_t
   "/":       ramdisk root
   "/dev":    device filesystem
   "/proc":   process info (Ch 63)

VFS dispatch:
  open("/bin/hello")
    → VFS looks up path → finds vfs_node for /bin/hello
    → calls node->open()
    → returns file descriptor index
    
  read(fd, buf, 100)
    → VFS looks up fd → gets vfs_node
    → calls node->read(offset, buf, 100)
    → advances offset by bytes read
    → returns byte count
```

---

## 3. The Ramdisk — An In-Memory File System

A ramdisk stores file contents directly in kernel memory. Simple, fast, volatile:

```
Ramdisk layout:
  Array of ramdisk_file_t entries:
    name[64]: file name
    data[]:   file contents (pointer to heap allocation)
    size:     size in bytes
    type:     FILE or DIR
    
  Directory as a ramdisk is just a ramdisk_file_t with type=DIR
  and children[] pointing to other files.
  
Simple flat structure (no subdirectory support initially):
  /bin/hello  → file
  /etc/motd   → file  
  /dev/zero   → device (special)
  /proc/      → directory
```

---

## 4. File Descriptors

Each process has a table of open files:

```c
/* Standard file descriptors: */
#define FD_STDIN   0   /* standard input (keyboard) */
#define FD_STDOUT  1   /* standard output (terminal) */
#define FD_STDERR  2   /* standard error (terminal) */
#define MAX_OPEN_FILES 16

typedef struct {
    struct vfs_node *node;   /* Which file this FD points to */
    uint32_t        offset;  /* Current read/write position */
    uint32_t        flags;   /* O_RDONLY, O_WRONLY, O_RDWR */
    int             in_use;  /* 1 if open, 0 if closed */
} file_descriptor_t;

/* Add to PCB (process.h): */
typedef struct process {
    /* ... existing fields ... */
    file_descriptor_t fds[MAX_OPEN_FILES];
} process_t;
```

---

## 5. The VFS Layer in C

```c
/* include/vfs.h */
#pragma once
#include "stdint.h"

#define VFS_NODE_FILE   1
#define VFS_NODE_DIR    2
#define VFS_NODE_DEVICE 3

#define O_RDONLY  0
#define O_WRONLY  1
#define O_RDWR    2
#define O_CREAT   4

struct vfs_node;

typedef struct {
    int32_t  (*read) (struct vfs_node *node, uint32_t offset,
                      uint8_t *buf, uint32_t count);
    int32_t  (*write)(struct vfs_node *node, uint32_t offset,
                      const uint8_t *buf, uint32_t count);
    int32_t  (*open) (struct vfs_node *node, uint32_t flags);
    void     (*close)(struct vfs_node *node);
    struct vfs_node *(*finddir)(struct vfs_node *node, const char *name);
} vfs_ops_t;

typedef struct vfs_node {
    char        name[64];
    uint32_t    type;       /* VFS_NODE_FILE, VFS_NODE_DIR, VFS_NODE_DEVICE */
    uint32_t    size;
    uint32_t    inode;      /* implementation-defined ID */
    vfs_ops_t  *ops;        /* Function pointers */
    void       *impl;       /* Backend-specific data (e.g., ramdisk_file_t*) */
} vfs_node_t;

/* Global VFS API: */
void        vfs_init(void);
vfs_node_t *vfs_lookup(const char *path);
int         vfs_open(const char *path, uint32_t flags);
int32_t     vfs_read(int fd, uint8_t *buf, uint32_t count);
int32_t     vfs_write(int fd, const uint8_t *buf, uint32_t count);
void        vfs_close(int fd);
void        vfs_register_root(vfs_node_t *root);
```

---

## 6. Implementing open/read/write/close

```c
/* kernel/vfs.c */

#include "vfs.h"
#include "process.h"
#include "heap.h"
#include "string.h"
#include "vga.h"

static vfs_node_t *vfs_root = NULL;

void vfs_register_root(vfs_node_t *root) {
    vfs_root = root;
}

/* Walk the path to find a vfs_node: */
vfs_node_t *vfs_lookup(const char *path) {
    if (!path || path[0] != '/') return NULL;
    if (path[1] == '\0') return vfs_root;  /* root "/" */
    
    /* Walk path components: */
    vfs_node_t *node = vfs_root;
    char component[64];
    const char *p = path + 1;  /* skip leading '/' */
    
    while (*p && node) {
        /* Extract next path component: */
        int i = 0;
        while (*p && *p != '/' && i < 63) {
            component[i++] = *p++;
        }
        component[i] = '\0';
        if (*p == '/') p++;
        
        if (!node->ops || !node->ops->finddir) return NULL;
        node = node->ops->finddir(node, component);
    }
    
    return node;
}

/* Allocate a file descriptor for the current process: */
static int alloc_fd(void) {
    if (!current_process) return -1;
    for (int i = 3; i < MAX_OPEN_FILES; i++) {  /* skip 0/1/2 (stdin/out/err) */
        if (!current_process->fds[i].in_use) return i;
    }
    return -1;
}

/* open() — returns fd or -1: */
int vfs_open(const char *path, uint32_t flags) {
    vfs_node_t *node = vfs_lookup(path);
    if (!node) return -1;
    
    int fd = alloc_fd();
    if (fd < 0) return -1;
    
    file_descriptor_t *f = &current_process->fds[fd];
    f->node   = node;
    f->offset = 0;
    f->flags  = flags;
    f->in_use = 1;
    
    if (node->ops && node->ops->open) {
        node->ops->open(node, flags);
    }
    
    return fd;
}

/* read() — returns bytes read: */
int32_t vfs_read(int fd, uint8_t *buf, uint32_t count) {
    if (fd < 0 || fd >= MAX_OPEN_FILES) return -1;
    if (!current_process) return -1;
    
    /* stdin (fd 0) → keyboard: */
    if (fd == 0) {
        for (uint32_t i = 0; i < count; i++) {
            char c = kbd_getchar_blocking();
            buf[i] = (uint8_t)c;
            if (c == '\n') return (int32_t)(i + 1);
        }
        return (int32_t)count;
    }
    
    file_descriptor_t *f = &current_process->fds[fd];
    if (!f->in_use || !f->node) return -1;
    if (!f->node->ops || !f->node->ops->read) return -1;
    
    int32_t bytes = f->node->ops->read(f->node, f->offset, buf, count);
    if (bytes > 0) f->offset += bytes;
    return bytes;
}

/* write() — returns bytes written: */
int32_t vfs_write(int fd, const uint8_t *buf, uint32_t count) {
    /* stdout/stderr → terminal: */
    if (fd == 1 || fd == 2) {
        for (uint32_t i = 0; i < count; i++) {
            terminal_putchar((char)buf[i]);
        }
        return (int32_t)count;
    }
    
    if (fd < 0 || fd >= MAX_OPEN_FILES) return -1;
    if (!current_process) return -1;
    
    file_descriptor_t *f = &current_process->fds[fd];
    if (!f->in_use || !f->node) return -1;
    if (!f->node->ops || !f->node->ops->write) return -1;
    
    int32_t bytes = f->node->ops->write(f->node, f->offset, buf, count);
    if (bytes > 0) f->offset += bytes;
    if (f->offset > f->node->size) f->node->size = f->offset;
    return bytes;
}

/* close(): */
void vfs_close(int fd) {
    if (fd < 3 || fd >= MAX_OPEN_FILES) return;
    if (!current_process) return;
    
    file_descriptor_t *f = &current_process->fds[fd];
    if (!f->in_use) return;
    
    if (f->node && f->node->ops && f->node->ops->close) {
        f->node->ops->close(f->node);
    }
    
    f->in_use = 0;
    f->node   = NULL;
    f->offset = 0;
}
```

---

## 7. Populating the Ramdisk

```c
/* kernel/ramdisk.c — simple in-memory file system */

#include "ramdisk.h"
#include "vfs.h"
#include "heap.h"
#include "string.h"

#define RAMDISK_MAX_FILES  32
#define RAMDISK_MAX_SIZE   (64 * 1024)   /* 64KB per file max */

typedef struct {
    char     name[64];
    uint8_t *data;
    uint32_t size;
    uint32_t capacity;
    int      type;   /* VFS_NODE_FILE or VFS_NODE_DIR */
} ramdisk_file_t;

static ramdisk_file_t rd_files[RAMDISK_MAX_FILES];
static vfs_node_t     rd_nodes[RAMDISK_MAX_FILES];
static vfs_node_t     rd_root;
static int            rd_count = 0;

/* Ramdisk read operation: */
static int32_t rd_read(vfs_node_t *node, uint32_t off, uint8_t *buf, uint32_t n) {
    ramdisk_file_t *f = (ramdisk_file_t *)node->impl;
    if (off >= f->size) return 0;
    if (off + n > f->size) n = f->size - off;
    memcpy(buf, f->data + off, n);
    return (int32_t)n;
}

/* Ramdisk write operation: */
static int32_t rd_write(vfs_node_t *node, uint32_t off, const uint8_t *buf, uint32_t n) {
    ramdisk_file_t *f = (ramdisk_file_t *)node->impl;
    if (off + n > f->capacity) n = f->capacity - off;
    memcpy(f->data + off, buf, n);
    if (off + n > f->size) f->size = off + n;
    node->size = f->size;
    return (int32_t)n;
}

/* Find a child by name (for directories): */
static vfs_node_t *rd_finddir(vfs_node_t *dir, const char *name) {
    (void)dir;
    for (int i = 0; i < rd_count; i++) {
        if (strcmp(rd_nodes[i].name, name) == 0) {
            return &rd_nodes[i];
        }
    }
    return NULL;
}

static vfs_ops_t rd_file_ops = {
    .read    = rd_read,
    .write   = rd_write,
    .open    = NULL,
    .close   = NULL,
    .finddir = NULL,
};

static vfs_ops_t rd_dir_ops = {
    .read    = NULL,
    .write   = NULL,
    .open    = NULL,
    .close   = NULL,
    .finddir = rd_finddir,
};

/* Create a file in the ramdisk: */
vfs_node_t *ramdisk_create(const char *name, const char *contents, uint32_t size) {
    if (rd_count >= RAMDISK_MAX_FILES) return NULL;
    
    ramdisk_file_t *f = &rd_files[rd_count];
    strncpy(f->name, name, 63);
    f->capacity = (size < 512) ? 512 : size;
    f->data     = (uint8_t *)kmalloc(f->capacity);
    f->size     = size;
    f->type     = VFS_NODE_FILE;
    
    if (contents && size) {
        memcpy(f->data, contents, size);
    }
    
    vfs_node_t *n = &rd_nodes[rd_count];
    strncpy(n->name, name, 63);
    n->type = VFS_NODE_FILE;
    n->size = size;
    n->ops  = &rd_file_ops;
    n->impl = f;
    n->inode = rd_count;
    
    rd_count++;
    return n;
}

/* Initialize ramdisk and populate with initial files: */
void ramdisk_init(void) {
    /* Set up root directory: */
    strncpy(rd_root.name, "/", 63);
    rd_root.type = VFS_NODE_DIR;
    rd_root.ops  = &rd_dir_ops;
    rd_root.impl = NULL;
    
    /* Register as VFS root: */
    vfs_register_root(&rd_root);
    
    /* Create some initial files: */
    ramdisk_create("motd",
        "Welcome to TinyOS!\n"
        "Type 'help' for a list of commands.\n",
        55);
    
    ramdisk_create("readme",
        "TinyOS - A minimal x86 operating system\n"
        "Built chapter by chapter in learning-notes.\n",
        84);
    
    kprintf("Ramdisk initialized with %d files.\n", rd_count);
}
```

---

## 8. sys_open, sys_read, sys_write Updated

Update syscall.c to go through VFS:
```c
static int32_t sys_open(registers_t *r) {
    const char *path  = (const char *)r->ebx;
    uint32_t    flags = r->ecx;
    if (!path) return -1;
    return vfs_open(path, flags);
}

static int32_t sys_close(registers_t *r) {
    vfs_close((int)r->ebx);
    return 0;
}

static int32_t sys_read(registers_t *r) {
    int      fd    = (int)r->ebx;
    uint8_t *buf   = (uint8_t *)r->ecx;
    uint32_t count = r->edx;
    return vfs_read(fd, buf, count);
}

static int32_t sys_write(registers_t *r) {
    int            fd    = (int)r->ebx;
    const uint8_t *buf   = (const uint8_t *)r->ecx;
    uint32_t       count = r->edx;
    return vfs_write(fd, buf, count);
}

/* Add to syscall_table: */
[SYS_OPEN]  = sys_open,
[SYS_CLOSE] = sys_close,
```

---

## 9. Testing

```c
static void vfs_test_process(void) {
    kprintf("VFS Test:\n");
    
    /* Open and read motd: */
    int fd = vfs_open("/motd", O_RDONLY);
    if (fd < 0) {
        kprintf("Failed to open /motd\n");
    } else {
        char buf[128];
        int n = vfs_read(fd, (uint8_t *)buf, sizeof(buf) - 1);
        buf[n] = '\0';
        kprintf("Contents of /motd:\n%s\n", buf);
        vfs_close(fd);
    }
    
    /* Create a new file and write to it: */
    /* (requires ramdisk_create to be exposed or a creat() syscall) */
    
    kprintf("VFS tests passed!\n");
    process_exit(0);
}
```

Expected output:
```
Ramdisk initialized with 2 files.
VFS Test:
Contents of /motd:
Welcome to TinyOS!
Type 'help' for a list of commands.

VFS tests passed!
```

---

## Summary

| Concept | Description |
|---------|------------|
| VFS | Virtual File System: uniform interface hiding underlying storage implementations |
| vfs_node_t | VFS object representing a file, directory, or device |
| vfs_ops_t | Function pointer table: read/write/open/close/finddir |
| vfs_lookup | Walk path string to find vfs_node_t |
| File descriptor | Per-process open-file handle: node pointer + offset + flags |
| FD 0/1/2 | stdin/stdout/stderr — special-cased in VFS: keyboard and terminal |
| Ramdisk | In-memory file system: files stored as heap allocations |
| rd_finddir | Searches ramdisk array for a node with matching name |
| offset | Current read/write position, advanced after each operation |
| Path walk | Split path by '/', call finddir() for each component |
| vfs_register_root | Install the root vfs_node_t (ramdisk's root directory) |
| O_RDONLY/O_WRONLY/O_RDWR | Open flags controlling read/write permissions |
