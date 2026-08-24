# Chapter 15: Inter-Process Communication (IPC)

> **"Processes are isolated by design — they can't read each other's memory. But they often need to work together. IPC is the set of controlled, OS-mediated channels that allow processes to collaborate without compromising their isolation."**

---

## Table of Contents

1. [Why IPC Exists](#1-why-ipc-exists)
2. [Overview of IPC Mechanisms](#2-overview-of-ipc-mechanisms)
3. [Pipes — The Unix Way](#3-pipes--the-unix-way)
4. [Named Pipes (FIFOs)](#4-named-pipes-fifos)
5. [Unix Domain Sockets](#5-unix-domain-sockets)
6. [Shared Memory — Fastest IPC](#6-shared-memory--fastest-ipc)
7. [Message Queues](#7-message-queues)
8. [Signals — Asynchronous Notifications](#8-signals--asynchronous-notifications)
9. [Network Sockets — Cross-Machine IPC](#9-network-sockets--cross-machine-ipc)
10. [Memory-Mapped Files](#10-memory-mapped-files)
11. [IPC in Practice: Designing a Pipeline](#11-ipc-in-practice-designing-a-pipeline)
12. [Windows IPC Mechanisms](#12-windows-ipc-mechanisms)
13. [Summary](#summary)

---

## 1. Why IPC Exists

Processes have isolated memory spaces by design. Process A cannot directly access process B's variables. This isolation is:
- **Good for security:** One buggy/malicious process can't corrupt others
- **Good for stability:** One crashing process doesn't affect others
- **Bad for communication:** Processes often need to share data or coordinate

IPC (Inter-Process Communication) is the set of OS-provided mechanisms that allow processes to communicate in a controlled way, without breaking isolation.

**Real uses of IPC:**
- Shell pipes: `ls | grep .c | wc -l` — three processes passing data through pipes
- Web server: parent process receives connections, forks child processes to handle them, uses shared memory for caching
- Chrome: browser UI process, renderer process, GPU process, network process — all communicate via IPC
- DBus: all GUI applications communicate with the desktop environment via a message bus
- Microservices: services communicate via network sockets (HTTP, gRPC)

---

## 2. Overview of IPC Mechanisms

| Mechanism | Speed | Persistence | Same machine? | Notes |
|-----------|-------|-------------|---------------|-------|
| Pipes | Medium | Temporary | Same machine only | Sequential; one direction |
| Named pipes | Medium | Until deleted | Same machine | Both directions (two pipes) |
| Unix sockets | Fast | Temporary | Same machine only | Full duplex, stream or datagram |
| Shared memory | Very fast | Until deleted | Same machine | Requires synchronization |
| Message queues | Medium | Until deleted | Same machine | Structured messages |
| Signals | Very fast | No data | Same machine | Only notification, no data |
| Network sockets | Slowest | N/A | Cross-machine | Full duplex, network overhead |
| Memory-mapped files | Fast | Persistent | Same machine | File-backed shared memory |

---

## 3. Pipes — The Unix Way

A **pipe** is a unidirectional, in-kernel buffer connecting two processes. Write to one end, read from the other.

```
        PIPE BUFFER (kernel)
        ┌──────────────────┐
write ─►│  data flows ──► │─► read
   (fd[1]) └──────────────────┘ (fd[0])
```

**Creating a pipe:**
```c
#include <unistd.h>

int fd[2];
pipe(fd);   // fd[0] = read end, fd[1] = write end
```

**Classic fork+pipe pattern:**
```c
#include <unistd.h>
#include <stdio.h>
#include <string.h>

int main() {
    int fd[2];
    pipe(fd);
    
    pid_t pid = fork();
    
    if (pid == 0) {
        // CHILD: writer
        close(fd[0]);                  // close unused read end
        char *msg = "Hello from child!\n";
        write(fd[1], msg, strlen(msg));
        close(fd[1]);
        return 0;
    } else {
        // PARENT: reader
        close(fd[1]);                  // close unused write end
        char buf[100];
        int n = read(fd[0], buf, sizeof(buf) - 1);
        buf[n] = '\0';
        printf("Parent received: %s", buf);
        close(fd[0]);
        return 0;
    }
}
```

**The shell pipe `ls | grep .c`:**
```c
// Simplified shell implementation of: ls | grep .c
int fd[2];
pipe(fd);

if (fork() == 0) {
    // ls process:
    dup2(fd[1], STDOUT_FILENO);  // stdout → pipe write end
    close(fd[0]);
    close(fd[1]);
    execve("/bin/ls", ...);
} else if (fork() == 0) {
    // grep process:
    dup2(fd[0], STDIN_FILENO);   // stdin → pipe read end
    close(fd[0]);
    close(fd[1]);
    execve("/bin/grep", (char*[]){"/bin/grep", ".c", NULL}, environ);
} else {
    // Parent (shell):
    close(fd[0]);
    close(fd[1]);
    wait(NULL); wait(NULL);
}
```

**Pipe properties:**
- **Atomic writes:** Writes ≤ PIPE_BUF (usually 4096 bytes) are atomic — won't be interleaved with other writers
- **Blocking reads:** If no data and write end open: read blocks
- **Broken pipe:** Writing to a pipe with no readers: SIGPIPE signal (or EPIPE error)
- **Capacity:** Typically 64KB on Linux (configurable)

---

## 4. Named Pipes (FIFOs)

**Anonymous pipes** only work between related processes (parent-child sharing the fd). 

**Named pipes (FIFOs)** appear as a file in the filesystem and can be used by ANY two processes.

```bash
# Create a named pipe:
mkfifo /tmp/mypipe

# Process 1 (reader):
cat /tmp/mypipe

# Process 2 (writer — separate terminal):
echo "Hello!" > /tmp/mypipe
```

**In C:**
```c
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

// Create the FIFO:
mkfifo("/tmp/mypipe", 0666);

// Writer:
int wfd = open("/tmp/mypipe", O_WRONLY);
write(wfd, "Hello", 5);
close(wfd);

// Reader (in another process):
int rfd = open("/tmp/mypipe", O_RDONLY);  // blocks until writer opens
char buf[100];
read(rfd, buf, sizeof(buf));
close(rfd);

// Cleanup:
unlink("/tmp/mypipe");
```

**Note:** Opening a FIFO for reading blocks until a writer opens it, and vice versa (by default). This natural synchronization can be a feature or a hazard.

---

## 5. Unix Domain Sockets

Unix domain sockets look like network sockets but communicate in-kernel (no network). They support:
- **SOCK_STREAM:** Bidirectional, reliable, ordered (like TCP but local)
- **SOCK_DGRAM:** Message-based (like UDP but reliable because it's local)

**Advantages over pipes:**
- Bidirectional (one socket, two directions)
- Can pass file descriptors between processes (magical!)
- Higher bandwidth than pipes on some systems

```c
// Server (listens for connections):
int server_fd = socket(AF_UNIX, SOCK_STREAM, 0);

struct sockaddr_un addr;
addr.sun_family = AF_UNIX;
strcpy(addr.sun_path, "/tmp/myapp.sock");

bind(server_fd, (struct sockaddr*)&addr, sizeof(addr));
listen(server_fd, 5);

int client_fd = accept(server_fd, NULL, NULL);
// Now communicate with client via client_fd

// Client:
int sock = socket(AF_UNIX, SOCK_STREAM, 0);
connect(sock, (struct sockaddr*)&addr, sizeof(addr));
// Now communicate via sock
```

**Passing file descriptors:**
A unique Unix feature: you can send an open file descriptor from one process to another through a Unix socket, using `sendmsg()` with `SCM_RIGHTS`.

```c
// Process A: send an open file descriptor to Process B:
// This transfers the actual open file (with its position, flags) to B
// B can then read/write the file even if it has no path access!
```

This is how privilege separation works in modern software: a privileged process opens a sensitive file, then passes the fd to an unprivileged process to use.

---

## 6. Shared Memory — Fastest IPC

**Shared memory** maps the same physical memory pages into multiple processes' virtual address spaces.

```
Process A's virtual space:    [code][heap][...][SHARED PAGE]
                                                     ↑
                                              same physical page
                                                     ↓
Process B's virtual space:    [code][heap][...][SHARED PAGE]
```

When A writes to the shared page, B immediately sees it — no copying, no kernel involvement. This is the **fastest possible IPC**.

**POSIX shared memory:**
```c
#include <sys/mman.h>
#include <fcntl.h>

// Process A (creator):
int shm_fd = shm_open("/myshm", O_CREAT | O_RDWR, 0666);
ftruncate(shm_fd, 4096);  // set size to 4KB

void *ptr = mmap(NULL, 4096, PROT_READ | PROT_WRITE, MAP_SHARED, shm_fd, 0);
// ptr now points to 4KB of shared memory

strcpy((char*)ptr, "Hello from A!");

// Process B (reader):
int shm_fd = shm_open("/myshm", O_RDONLY, 0666);
void *ptr = mmap(NULL, 4096, PROT_READ, MAP_SHARED, shm_fd, 0);
printf("B reads: %s\n", (char*)ptr);  // "Hello from A!"

// Cleanup:
munmap(ptr, 4096);
shm_unlink("/myshm");  // removes the shared memory object
```

**Critical: shared memory requires explicit synchronization!**
If both processes read and write simultaneously, you have a race condition.

```c
// Use a semaphore to protect the shared memory:
sem_t *sem = sem_open("/mysem", O_CREAT, 0666, 1);

// Before reading/writing shared memory:
sem_wait(sem);
// ... use shared memory ...
sem_post(sem);
```

---

## 7. Message Queues

Message queues allow processes to send discrete messages (typed packets).

**Advantages over pipes:**
- Messages have types — a reader can receive only messages of a specific type
- Non-blocking is easy
- Messages are preserved across process restarts (until explicitly deleted)
- Multiple readers can receive different message types from the same queue

**POSIX message queues:**
```c
#include <mqueue.h>

// Writer:
mqd_t mq = mq_open("/myqueue", O_CREAT | O_WRONLY, 0666, NULL);
char *msg = "Hello!";
mq_send(mq, msg, strlen(msg), 1);  // priority = 1
mq_close(mq);

// Reader:
mqd_t mq = mq_open("/myqueue", O_RDONLY);
char buf[1024];
unsigned int priority;
mq_receive(mq, buf, sizeof(buf), &priority);
printf("Received: %s (priority %u)\n", buf, priority);
mq_close(mq);
mq_unlink("/myqueue");
```

**Message queues vs. pipes:**
- Pipes: byte stream, no message boundaries
- Message queues: discrete messages, preserved boundaries, typed priorities

---

## 8. Signals — Asynchronous Notifications

**Signals** are asynchronous notifications sent to a process. They don't carry arbitrary data — just a signal number.

Think of signals as software interrupts.

**Common signals:**

| Signal | Number | Default action | Meaning |
|--------|--------|----------------|---------|
| SIGTERM | 15 | Terminate | Request graceful shutdown |
| SIGKILL | 9 | Kill | Force kill (cannot be caught) |
| SIGINT | 2 | Terminate | Ctrl+C |
| SIGQUIT | 3 | Core dump | Ctrl+\\ |
| SIGSEGV | 11 | Core dump | Segmentation fault |
| SIGFPE | 8 | Core dump | Floating point error |
| SIGCHLD | 17 | Ignore | Child exited |
| SIGHUP | 1 | Terminate | Hangup (terminal closed) |
| SIGSTOP | 19 | Stop | Pause process (cannot be caught) |
| SIGCONT | 18 | Continue | Resume paused process |
| SIGUSR1/2 | 10/12 | Terminate | User-defined |

**Sending signals:**
```bash
kill -SIGTERM 1234   # request graceful shutdown of PID 1234
kill -9 1234         # force kill (SIGKILL)
kill -SIGUSR1 1234   # send custom signal
```

```c
#include <signal.h>
kill(1234, SIGTERM);   // send SIGTERM to PID 1234
kill(0, SIGTERM);      // send to all processes in my process group
```

**Handling signals:**
```c
#include <signal.h>

void sigterm_handler(int sig) {
    printf("Received SIGTERM, shutting down gracefully...\n");
    cleanup();
    exit(0);
}

void sigusr1_handler(int sig) {
    printf("Received SIGUSR1 — reloading config\n");
    reload_config();
}

int main() {
    signal(SIGTERM, sigterm_handler);   // register handler
    signal(SIGUSR1, sigusr1_handler);
    signal(SIGCHLD, SIG_IGN);          // ignore child exit signals
    
    // Main program loop:
    while (1) {
        // ... do work ...
    }
}
```

**Signals are limited for data passing:**
- Cannot carry arbitrary data (just the signal number)
- Small queue (RT signals can be queued, but classic signals are not)
- Signal handlers must be async-signal-safe (only certain functions allowed)
- Better for notifications ("reload config", "dump stats", "shut down") than data

---

## 9. Network Sockets — Cross-Machine IPC

For communication between different machines (or even different processes across a network):

**TCP socket (reliable, ordered, stream):**
```c
// Server:
int server_fd = socket(AF_INET, SOCK_STREAM, 0);

struct sockaddr_in addr = {
    .sin_family = AF_INET,
    .sin_port = htons(8080),
    .sin_addr.s_addr = INADDR_ANY
};
bind(server_fd, (struct sockaddr*)&addr, sizeof(addr));
listen(server_fd, 10);

int client_fd = accept(server_fd, NULL, NULL);
char buf[1024];
int n = recv(client_fd, buf, sizeof(buf), 0);
// process request...
send(client_fd, "HTTP/1.1 200 OK\r\n...", ..., 0);

// Client:
int sock = socket(AF_INET, SOCK_STREAM, 0);
connect(sock, (struct sockaddr*)&server_addr, sizeof(server_addr));
send(sock, "GET / HTTP/1.1\r\n...", ..., 0);
recv(sock, response, sizeof(response), 0);
```

Network sockets are the foundation of all internet communication, microservices, and distributed systems. Every HTTP request, database query, and API call uses sockets under the hood.

---

## 10. Memory-Mapped Files

`mmap()` can map a FILE into memory. Multiple processes mapping the same file share the same physical pages.

```c
#include <sys/mman.h>
#include <fcntl.h>

int fd = open("datafile.dat", O_RDWR);
lseek(fd, 1024*1024 - 1, SEEK_SET);
write(fd, "", 1);  // create 1MB file

// Map the file into memory:
uint8_t *data = mmap(NULL, 1024*1024, 
                      PROT_READ | PROT_WRITE, 
                      MAP_SHARED,           // MAP_SHARED: writes go to file
                      fd, 0);

// Now access the file like an array:
data[0] = 0xFF;    // writes to the file!
uint32_t x = *(uint32_t*)(data + 100);  // reads from file

// Changes to MAP_SHARED are visible to all processes mapping the same file!
// This is inter-process shared memory backed by a persistent file.

msync(data, 1024*1024, MS_SYNC);  // ensure changes written to disk
munmap(data, 1024*1024);
close(fd);
```

**mmap vs. read/write:**
- mmap: OS handles paging; reads only what's accessed; great for large files
- read/write: explicit, portable, works over network sockets (mmap doesn't)
- mmap is faster for random access to large files (no user-kernel copying)

---

## 11. IPC in Practice: Designing a Pipeline

Real applications combine multiple IPC mechanisms. Let's design a image processing pipeline:

```
[Camera Process] ─── shared memory ──► [Processor Process] ─── pipe ──► [Writer Process]
                         ↑                      │                               │
                     semaphore              signals                         file I/O
                   (synchronize)          (control)                     (persistence)
```

**Camera Process:**
```c
// Grabs frames from camera hardware
// Writes each frame to shared memory buffer
// Signals Processor via semaphore
sem_post(&frame_ready);
```

**Processor Process:**
```c
// Waits for frame in shared memory
sem_wait(&frame_ready);
// Processes frame (filter, encode, etc.)
// Sends result via pipe to Writer
write(pipe_fd[1], &result_metadata, sizeof(result_metadata));
```

**Writer Process:**
```c
// Reads from pipe
read(pipe_fd[0], &metadata, sizeof(metadata));
// Reads image data from shared memory
// Writes to disk file
```

**Control (signal):**
```bash
# Admin sends SIGUSR1 → all processes rotate their log files
kill -SIGUSR1 -$(pgrep -g camera_pipeline)
```

This is how real media processing pipelines work (GStreamer, FFmpeg pipelines, etc.).

---

## 12. Windows IPC Mechanisms

Windows has its own IPC mechanisms:

**Anonymous pipes:**
Similar to Unix pipes. `CreatePipe()` creates read/write handles. Used for standard I/O redirection.

**Named pipes:**
More powerful than Unix FIFOs:
- Bidirectional (one handle for both read and write)
- Multiple clients can connect
- Works over network (via `\\server\pipe\name`)

**Mailslots:**
One-to-many broadcasting of messages. One writer can send to multiple readers.

**Shared memory (file mapping):**
```c
// Windows:
HANDLE hMapping = CreateFileMapping(INVALID_HANDLE_VALUE, NULL,
                                    PAGE_READWRITE, 0, 4096, L"MySharedMemory");
void *ptr = MapViewOfFile(hMapping, FILE_MAP_ALL_ACCESS, 0, 0, 4096);
```

**COM/DCOM (Component Object Model):**
Windows's rich IPC framework for structured communication between Windows applications (used internally by Windows components, Office, etc.).

**Windows messages (WM_*):**
GUI applications communicate via Windows messages posted to message queues. Used for all GUI events (mouse clicks, keyboard, redraw, etc.). Inter-process messaging is possible via `PostMessage` to another window's handle.

---

## Summary

| Mechanism | Use When | Speed | Notes |
|-----------|---------|-------|-------|
| Pipes | Parent-child, shell pipelines | Medium | Unidirectional, byte stream |
| Named pipes (FIFOs) | Unrelated processes, simple | Medium | Filesystem-visible |
| Unix sockets | Fast local IPC, fd passing | Fast | Full-duplex, stream or datagram |
| Shared memory | High-speed data sharing | Fastest | Must add synchronization |
| Message queues | Typed messages, priority | Medium | Preserved messages |
| Signals | Asynchronous notifications | Instant | No data, just signal number |
| Network sockets | Cross-machine, services | Slowest | Standard for distributed systems |
| mmap files | Large file sharing | Fast | File-backed, persistent |

**Choosing the right IPC:**

```
Need to pass data between parent and child?
    → pipes (simplest)
    
Need bidirectional communication, same machine?
    → Unix domain sockets
    
Need maximum throughput, same machine?
    → shared memory + semaphore for sync
    
Need to send to unrelated process, any machine?
    → network sockets
    
Need to send a notification (not data)?
    → signals

Need typed messages that survive process restart?
    → message queues

Need persistent shared data?
    → mmap files
```
