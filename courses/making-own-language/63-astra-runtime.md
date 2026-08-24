# Chapter 63: The Astra Runtime — The Foundation That Runs Every Astra Program

> "The runtime is the invisible stage crew that sets up the theater before the actors arrive, and cleans up after they leave. Nobody sees it, but without it, there is no show."
> — paraphrased from systems programming lore

---

## Overview

You have built a lexer, a parser, a type checker, and a code generator. When you run `astrac build main.as`, your compiler reads the Astra source and produces a binary. But what happens when that binary actually *runs*?

Before your `main()` function gets its first instruction, something else happens: the **runtime** starts up. It prepares the stack, the heap, the panic handler, and the built-in functions. Only then does it hand control to your code. When `main()` returns — or when a panic occurs — the runtime cleans up and exits.

This chapter builds `runtime/runtime.c`, the C file that every Astra binary is linked against. It is the bedrock of the language.

**What you will understand after this chapter:**
- What a "runtime" actually is and what it must do
- How the Astra binary starts (hint: it is not `main()`)
- How `print("Hello!")` in Astra becomes a `write()` system call
- How panics propagate and terminate programs cleanly
- How to compile and link the runtime against your programs

---

## What We Are Building

```
┌─────────────────────────────────────────────────────────────┐
│                    ASTRA BINARY (todo_server)                │
├─────────────────────────────────────────────────────────────┤
│  Your Astra Code (compiled to machine code)                  │
│    fn main() { print("hello") }                              │
├─────────────────────────────────────────────────────────────┤
│  Astra Runtime (runtime.o, linked in)                        │
│    _astra_start()   — real entry point                       │
│    astra_panic()    — error handling                         │
│    astra_print()    — I/O                                    │
│    astra_alloc()    — memory                                 │
│    astra_gc_*()     — garbage collection                     │
├─────────────────────────────────────────────────────────────┤
│  C Standard Library (libc)                                   │
│    malloc, free, printf, write, exit                         │
├─────────────────────────────────────────────────────────────┤
│  Operating System (macOS / Linux)                            │
│    System calls: write(2), mmap(2), exit(2)                  │
└─────────────────────────────────────────────────────────────┘
```

Every Astra program links against `runtime.o`. The runtime provides the C glue between Astra's compiled code and the OS.

---

## Table of Contents

1. What Is a Runtime? (And What It Is Not)
2. Runtime Responsibilities — The Full List
3. The Entry Point: `_astra_start`
4. The Panic System
5. Built-in String Functions
6. Built-in I/O Functions
7. Runtime Assertions
8. Stack Unwinding with setjmp/longjmp
9. Exit Codes and Process Cleanup
10. The Complete runtime.c File
11. Compiling and Linking the Runtime
12. Tracing "Hello!" from Astra to the OS
13. Astra Build Milestone
14. Exercises

---

## 1. What Is a Runtime? (And What It Is Not)

The word "runtime" is genuinely confusing because it has two meanings in programming:

**Meaning 1 — A period of time:** "Runtime error" means an error that happens *while the program is running* (as opposed to a compile-time error). This is the meaning you already know.

**Meaning 2 — A support library:** A runtime is a library of C (or assembly) code that is automatically linked into every program written in a given language. This is what we are building.

Think of it this way: when you write `fn main()` in Astra, you are the playwright. But the theater needs stage crew. The runtime is the stage crew.

Languages you already know all have runtimes:
- **Python**: CPython is almost entirely "runtime" — it provides the object model, GC, module system, and interpreter loop.
- **Java**: The JVM is the runtime. It starts up before your `main()`, verifies bytecode, compiles it to machine code, and manages memory.
- **Go**: The Go runtime (runtime package) manages goroutines, GC, channel communication, and the scheduler.
- **C**: Even C has a runtime! It is called `crt0` ("C Runtime 0") and is provided by your OS. It calls your `main()`.
- **Astra**: We build a small custom runtime in C. It handles startup, panic, I/O, and memory.

**What the Astra runtime is NOT:**
- It is not the compiler (`astrac`). The compiler runs before the program, to translate source code.
- It is not an interpreter. Astra programs are fully compiled to machine code; the runtime is just a C library they link against.
- It is not the standard library. The standard library (string, math, file, http) is separate and built on top of the runtime.

---

## 2. Runtime Responsibilities — The Full List

Here is everything the Astra runtime must handle:

| Responsibility | Function | Description |
|---|---|---|
| Program startup | `_astra_start` | Real entry point; sets up heap, calls `main` |
| Panic handling | `astra_panic` | Print error, stack trace, exit(1) |
| Standard output | `astra_print`, `astra_println` | Write strings to stdout |
| Integer to string | `astra_int_to_string` | Convert int64 → heap string |
| Float to string | `astra_float_to_string` | Convert float64 → heap string |
| String concat | `astra_string_concat` | Allocate and join two strings |
| String equality | `astra_string_eq` | Compare two strings byte-by-byte |
| Runtime assertion | `astra_assert` | Panic if condition is false |
| Memory allocation | `astra_alloc` | Thin wrapper over GC allocator |
| Bounds checking | `astra_bounds_check` | Panic on out-of-bounds array access |
| Division by zero | (implicit) | Caught by OS signal SIGFPE → panic |
| Null dereference | (implicit) | Caught by OS signal SIGSEGV → panic |

The runtime deliberately stays small. Most things belong in the standard library (stdlib) or in compiled Astra code. The runtime only contains what *every* Astra program needs, with no exceptions.

---

## 3. The Entry Point: `_astra_start`

Here is something that might surprise you: your program does not start at `main()`. On most operating systems, the OS hands control to a function called `_start` (or `_astra_start` in our case). This function then calls `main()`.

Why? Because several things must happen before your code runs:
1. The C standard library must be initialized (`__libc_start_main` on Linux).
2. Command-line arguments (`argc`, `argv`) must be set up.
3. Environment variables must be available.
4. Any language-level setup (GC heap, panic handlers, etc.) must be done.

In Astra, `_astra_start` handles all of this:

```c
// runtime/runtime.c — excerpt: startup

// Forward declaration of the user's compiled main function
extern void astra_user_main(void);

// _astra_start is the true entry point of every Astra binary.
// The linker is told to use this instead of the default `main`.
void _astra_start(void) {
    // 1. Initialize the GC heap (allocate the initial heap block)
    gc_init(ASTRA_HEAP_SIZE_INITIAL);

    // 2. Set up the panic jump buffer (for setjmp/longjmp unwinding)
    if (setjmp(astra_panic_buf) != 0) {
        // We land here when a panic unwinds the stack
        write(STDERR_FILENO, astra_panic_message,
              strlen(astra_panic_message));
        write(STDERR_FILENO, "\n", 1);
        _exit(1);
    }

    // 3. Call the user's compiled main function
    astra_user_main();

    // 4. Clean up and exit with success
    gc_shutdown();
    _exit(0);
}
```

When the Astra compiler compiles `fn main() { ... }`, it emits a C function called `astra_user_main`. The runtime calls that. When it returns, the runtime exits with code 0.

**Why `_exit` instead of `exit`?**

`exit()` runs C atexit handlers and flushes stdio buffers. `_exit()` goes directly to the OS `exit` syscall. Astra manages its own cleanup, so we bypass C's cleanup chain.

---

## 4. The Panic System

When something goes wrong in an Astra program — index out of bounds, null dereference, explicit `panic("message")` — the runtime must:
1. Print the error message.
2. Print a stack trace (simplified in v1.0).
3. Exit with code 1.

```c
// runtime/runtime.c — panic handling

// The jump buffer for panic unwinding
jmp_buf astra_panic_buf;

// The current panic message (set before longjmp)
static char astra_panic_message[4096];

// Called by compiled Astra code when a panic occurs.
// msg: the panic message string (null-terminated)
void astra_panic(const char* msg) {
    // Format the panic message into our buffer
    snprintf(astra_panic_message, sizeof(astra_panic_message),
             "\n[PANIC] %s\n"
             "Stack trace: (not available in Astra v1.0)\n"
             "Program aborted.\n",
             msg);

    // longjmp back to _astra_start's setjmp, which will print and exit
    longjmp(astra_panic_buf, 1);
}
```

**Why setjmp/longjmp?**

When a panic occurs deep inside a call stack:

```
_astra_start
  └── astra_user_main        // fn main()
        └── process_request  // fn process_request()
              └── find_user   // fn find_user()
                    └── astra_panic("user not found")  ← PANIC HERE
```

We need to unwind all those frames and return to `_astra_start` to print the error and exit. `longjmp` does exactly this — it jumps back to wherever `setjmp` was called, regardless of how deep the call stack is.

This is a simplified approach. A production runtime (like Go's) maintains explicit stack metadata and can print a real stack trace. We will add that in a later volume.

---

## 5. Built-in String Functions

Astra's string operations are implemented as C functions in the runtime. The Astra compiler emits calls to these functions whenever string operations appear in your code.

**String representation:** An Astra string is not just a `char*`. It is a struct:

```c
// runtime/runtime.h

typedef struct {
    char*  data;    // pointer to UTF-8 bytes (not null-terminated internally)
    size_t len;     // number of bytes
    size_t cap;     // allocated capacity (for mutable strings)
} AstraString;
```

Most Astra strings are immutable. When you write `let s = "hello"`, the compiler embeds the bytes in the binary as a constant, and creates an `AstraString` pointing to them with `len = 5`, `cap = 0` (immutable marker).

**String concatenation:**

```c
// Concatenate two Astra strings, return a new heap-allocated string.
AstraString* astra_string_concat(const AstraString* a, const AstraString* b) {
    size_t new_len = a->len + b->len;
    char* buf = (char*)astra_alloc(new_len + 1);  // +1 for safety null byte
    
    memcpy(buf, a->data, a->len);
    memcpy(buf + a->len, b->data, b->len);
    buf[new_len] = '\0';  // null-terminate for C interop safety
    
    AstraString* result = (AstraString*)astra_alloc(sizeof(AstraString));
    result->data = buf;
    result->len  = new_len;
    result->cap  = new_len;
    return result;
}
```

**String equality:**

```c
// Compare two Astra strings. Returns 1 if equal, 0 if not.
int astra_string_eq(const AstraString* a, const AstraString* b) {
    if (a->len != b->len) return 0;
    return memcmp(a->data, b->data, a->len) == 0 ? 1 : 0;
}
```

**Integer to string:**

```c
// Convert an int64 to a heap-allocated AstraString.
AstraString* astra_int_to_string(int64_t n) {
    // snprintf to a local buffer first to measure length
    char buf[32];
    int  len = snprintf(buf, sizeof(buf), "%lld", (long long)n);
    
    char* data = (char*)astra_alloc(len + 1);
    memcpy(data, buf, len + 1);
    
    AstraString* s = (AstraString*)astra_alloc(sizeof(AstraString));
    s->data = data;
    s->len  = (size_t)len;
    s->cap  = (size_t)len;
    return s;
}
```

**Float to string:**

```c
// Convert a float64 (double) to a heap-allocated AstraString.
AstraString* astra_float_to_string(double f) {
    char buf[64];
    int len;
    
    // Handle special cases
    if (isinf(f)) {
        len = snprintf(buf, sizeof(buf), f > 0 ? "Inf" : "-Inf");
    } else if (isnan(f)) {
        len = snprintf(buf, sizeof(buf), "NaN");
    } else {
        // Use %g for compact representation: "3.14" not "3.140000"
        len = snprintf(buf, sizeof(buf), "%g", f);
    }
    
    char* data = (char*)astra_alloc(len + 1);
    memcpy(data, buf, len + 1);
    
    AstraString* s = (AstraString*)astra_alloc(sizeof(AstraString));
    s->data = data;
    s->len  = (size_t)len;
    s->cap  = (size_t)len;
    return s;
}
```

---

## 6. Built-in I/O Functions

Astra's `print()` and `println()` are the most fundamental I/O operations. They map directly to the runtime:

```c
// Write a string to stdout without a newline.
void astra_print(const AstraString* s) {
    if (s == NULL || s->data == NULL || s->len == 0) return;
    
    // Use write() directly for reliability and simplicity.
    // write(fd, buffer, count) is a POSIX system call.
    ssize_t written = 0;
    while ((size_t)written < s->len) {
        ssize_t n = write(STDOUT_FILENO,
                          s->data + written,
                          s->len - (size_t)written);
        if (n < 0) {
            if (errno == EINTR) continue;  // interrupted by signal, retry
            break;                          // real error, give up
        }
        written += n;
    }
}

// Write a string to stdout, then write a newline.
void astra_println(const AstraString* s) {
    astra_print(s);
    write(STDOUT_FILENO, "\n", 1);
}

// Write a string to stderr (for internal runtime messages).
void astra_eprint(const char* msg) {
    write(STDERR_FILENO, msg, strlen(msg));
}
```

---

## 7. Runtime Assertions

Astra inserts automatic runtime checks for safety. The compiler emits calls to `astra_bounds_check` before every array access, and `astra_assert` wherever user code calls `assert()`.

```c
// Bounds check: called before every array index operation.
// If index >= length, panic.
void astra_bounds_check(int64_t index, int64_t length,
                         const char* file, int line) {
    if (index < 0 || index >= length) {
        char msg[256];
        snprintf(msg, sizeof(msg),
                 "index out of bounds: index %lld, length %lld "
                 "(at %s:%d)",
                 (long long)index, (long long)length, file, line);
        astra_panic(msg);
    }
}

// User-callable assertion.
void astra_assert(int condition, const AstraString* message) {
    if (!condition) {
        char msg[512];
        snprintf(msg, sizeof(msg), "assertion failed: %.*s",
                 (int)message->len, message->data);
        astra_panic(msg);
    }
}
```

When the Astra compiler sees:
```astra
let x = arr[5]
```

It compiles it to something like:
```c
astra_bounds_check(5, arr->length, "main.as", 42);
int64_t x = arr->data[5];
```

This means Astra programs can *never* access memory out of bounds and cause a buffer overflow — the runtime catches it first and panics gracefully.

---

## 8. Stack Unwinding with setjmp/longjmp

Let us dig deeper into how panics unwind the stack. Here is what `setjmp` and `longjmp` actually do:

```
                   CALL STACK (grows downward)
                   
 _astra_start ──► [setjmp saves CPU state here]
                    ↓
 astra_user_main  [frame A]
                    ↓
 process_items    [frame B]
                    ↓
 parse_item       [frame C]
                    ↓
 astra_panic()   ──► longjmp(astra_panic_buf, 1)
                         │
                         │ STACK UNWINDS — all frames B, C, D are
                         │ popped instantly. CPU state restored.
                         │
                         ▼
 _astra_start ◄── setjmp returns 1 (the longjmp value)
                  [prints panic message, calls _exit(1)]
```

`setjmp` saves the current CPU state (stack pointer, instruction pointer, register values) into a `jmp_buf` structure. Later, `longjmp` restores that saved state, making the CPU act as if `setjmp` just returned — but with a non-zero return value.

**What does NOT happen during longjmp unwinding:**
- C++ destructors do not run (Astra is not C++).
- `defer` statements in Astra do NOT run in v1.0. (A production runtime would run defers during unwind — this is future work.)
- No memory is freed. The GC will handle it on the next collection cycle.

**Important:** The `jmp_buf` must remain valid for the entire life of the program. That is why it is a global variable in the runtime, not a local variable in a function.

---

## 9. Exit Codes and Process Cleanup

POSIX programs signal success or failure through an **exit code** — an integer 0-255 that the OS gives back to whoever ran the program (a shell, another program, a CI system):

| Exit Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Generic error (Astra panic) |
| 2 | Misuse of shell command |
| 126 | Permission denied |
| 127 | Command not found |
| 128+N | Killed by signal N |

Astra uses:
- Exit 0: `main()` returned normally.
- Exit 1: Any panic occurred.

```c
// Called at the end of _astra_start after panic or success.
static void astra_cleanup(int exit_code) {
    // Flush any buffered output (we use write() directly so this is a no-op)
    // Run GC finalization
    gc_shutdown();
    // Exit with the given code
    _exit(exit_code);
}
```

---

## 10. The Complete runtime.c File

Here is the full `runtime/runtime.c` — about 220 lines of carefully written C that every Astra program depends on:

```c
// runtime/runtime.c
// The Astra Language Runtime
// Every Astra binary is linked against this file.
// Build: cc -c -O2 runtime.c -o runtime.o

#include <stdint.h>
#include <stddef.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <math.h>
#include <setjmp.h>
#include "runtime.h"
#include "gc.h"

// ── Panic system ─────────────────────────────────────────────

jmp_buf astra_panic_buf;
static char astra_panic_message[4096];
static int  astra_in_panic = 0;

void astra_panic(const char* msg) {
    if (astra_in_panic) {
        // Nested panic: go straight to _exit to avoid infinite loop
        const char* nested = "[PANIC] nested panic, aborting\n";
        write(STDERR_FILENO, nested, strlen(nested));
        _exit(1);
    }
    astra_in_panic = 1;

    snprintf(astra_panic_message, sizeof(astra_panic_message),
             "\n[PANIC] %s\n"
             "Stack trace: (run with ASTRA_DEBUG=1 for details)\n"
             "Program exited with code 1.\n",
             msg ? msg : "(no message)");

    longjmp(astra_panic_buf, 1);
}

// ── Program entry point ──────────────────────────────────────

// Forward declaration: the user's compiled main function.
extern void astra_user_main(void);

// The linker sets this as the entry point instead of the default main().
void _astra_start(void) {
    // Initialize the garbage collector's heap
    gc_init(8 * 1024 * 1024);  // 8 MB initial heap

    // Set up the panic landing zone
    int panic_code = setjmp(astra_panic_buf);
    if (panic_code != 0) {
        // We arrive here after a longjmp from astra_panic()
        write(STDERR_FILENO,
              astra_panic_message,
              strlen(astra_panic_message));
        gc_shutdown();
        _exit(1);
    }

    // Call the user's Astra main() function
    astra_user_main();

    // Normal exit
    gc_shutdown();
    _exit(0);
}

// ── I/O functions ─────────────────────────────────────────────

void astra_print(const AstraString* s) {
    if (!s || !s->data || s->len == 0) return;
    ssize_t written = 0;
    while ((size_t)written < s->len) {
        ssize_t n = write(STDOUT_FILENO,
                          s->data + written,
                          s->len - (size_t)written);
        if (n < 0) {
            if (errno == EINTR) continue;
            return;
        }
        written += n;
    }
}

void astra_println(const AstraString* s) {
    astra_print(s);
    write(STDOUT_FILENO, "\n", 1);
}

void astra_eprint(const char* msg) {
    if (!msg) return;
    write(STDERR_FILENO, msg, strlen(msg));
}

// ── String functions ──────────────────────────────────────────

AstraString* astra_string_new(const char* data, size_t len) {
    AstraString* s = (AstraString*)gc_alloc(sizeof(AstraString),
                                             ASTRA_TYPE_STRING);
    s->data = (char*)gc_alloc(len + 1, ASTRA_TYPE_RAW_BYTES);
    memcpy(s->data, data, len);
    s->data[len] = '\0';
    s->len = len;
    s->cap = len;
    return s;
}

AstraString* astra_string_concat(const AstraString* a,
                                  const AstraString* b) {
    if (!a) return (AstraString*)b;
    if (!b) return (AstraString*)a;

    size_t new_len = a->len + b->len;
    char*  buf     = (char*)gc_alloc(new_len + 1, ASTRA_TYPE_RAW_BYTES);

    memcpy(buf, a->data, a->len);
    memcpy(buf + a->len, b->data, b->len);
    buf[new_len] = '\0';

    AstraString* result = (AstraString*)gc_alloc(sizeof(AstraString),
                                                   ASTRA_TYPE_STRING);
    result->data = buf;
    result->len  = new_len;
    result->cap  = new_len;
    return result;
}

int astra_string_eq(const AstraString* a, const AstraString* b) {
    if (!a && !b) return 1;
    if (!a || !b) return 0;
    if (a->len != b->len) return 0;
    return memcmp(a->data, b->data, a->len) == 0 ? 1 : 0;
}

AstraString* astra_int_to_string(int64_t n) {
    char buf[32];
    int  len = snprintf(buf, sizeof(buf), "%lld", (long long)n);
    return astra_string_new(buf, (size_t)len);
}

AstraString* astra_float_to_string(double f) {
    char buf[64];
    int  len;
    if (isinf(f))      len = snprintf(buf, sizeof(buf),
                                       f > 0 ? "Inf" : "-Inf");
    else if (isnan(f)) len = snprintf(buf, sizeof(buf), "NaN");
    else               len = snprintf(buf, sizeof(buf), "%g", f);
    return astra_string_new(buf, (size_t)len);
}

AstraString* astra_bool_to_string(int b) {
    return b ? astra_string_new("true", 4)
             : astra_string_new("false", 5);
}

// ── Safety checks ─────────────────────────────────────────────

void astra_bounds_check(int64_t index, int64_t length,
                         const char* file, int line) {
    if (index < 0 || index >= length) {
        char msg[256];
        snprintf(msg, sizeof(msg),
                 "index out of bounds: index=%lld, length=%lld (at %s:%d)",
                 (long long)index, (long long)length, file, line);
        astra_panic(msg);
    }
}

void astra_null_check(const void* ptr, const char* name,
                       const char* file, int line) {
    if (!ptr) {
        char msg[256];
        snprintf(msg, sizeof(msg),
                 "null pointer dereference: '%s' is null (at %s:%d)",
                 name ? name : "value", file, line);
        astra_panic(msg);
    }
}

void astra_assert(int condition, const AstraString* message) {
    if (!condition) {
        char msg[512];
        if (message && message->data) {
            snprintf(msg, sizeof(msg), "assertion failed: %.*s",
                     (int)message->len, message->data);
        } else {
            snprintf(msg, sizeof(msg), "assertion failed");
        }
        astra_panic(msg);
    }
}

// ── Division guard ────────────────────────────────────────────

int64_t astra_int_div(int64_t a, int64_t b,
                       const char* file, int line) {
    if (b == 0) {
        char msg[128];
        snprintf(msg, sizeof(msg),
                 "division by zero (at %s:%d)", file, line);
        astra_panic(msg);
    }
    return a / b;
}

double astra_float_div(double a, double b,
                        const char* file, int line) {
    if (b == 0.0) {
        char msg[128];
        snprintf(msg, sizeof(msg),
                 "float division by zero (at %s:%d)", file, line);
        astra_panic(msg);
    }
    return a / b;
}
```

---

## 11. Compiling and Linking the Runtime

The runtime is compiled to an object file and linked into every Astra binary. Here is how the Astra compiler's build pipeline works:

```bash
# Step 1: Compile the runtime (done once, or when runtime.c changes)
cc -c -O2 -Wall \
   runtime/runtime.c \
   -o runtime/runtime.o

# Step 2: Compile the GC (Chapter 64)
cc -c -O2 -Wall \
   runtime/gc.c \
   -o runtime/gc.o

# Step 3: Compile the user's Astra code
#   (astrac generates a .c file, then compiles it)
astrac codegen main.as -o main_generated.c
cc -c -O2 main_generated.c -o main.o

# Step 4: Link everything together
cc main.o \
   runtime/runtime.o \
   runtime/gc.o \
   -lm \          # math library
   -o main        # final binary

# Step 5: Run!
./main
```

The Astra compiler (`astrac`) handles steps 2-4 automatically when you run `astrac build main.as`. But this is what it does internally.

**Header file (runtime.h):**

```c
// runtime/runtime.h
// Declarations for the Astra runtime.
// Included by generated C code from the Astra compiler.

#ifndef ASTRA_RUNTIME_H
#define ASTRA_RUNTIME_H

#include <stdint.h>
#include <stddef.h>
#include <setjmp.h>

// String type
typedef struct {
    char*  data;
    size_t len;
    size_t cap;
} AstraString;

// Type IDs for GC
#define ASTRA_TYPE_STRING     1
#define ASTRA_TYPE_RAW_BYTES  2
#define ASTRA_TYPE_STRUCT     3
#define ASTRA_TYPE_LIST       4
#define ASTRA_TYPE_CLOSURE    5

// Panic
extern jmp_buf astra_panic_buf;
void astra_panic(const char* msg);

// I/O
void astra_print(const AstraString* s);
void astra_println(const AstraString* s);
void astra_eprint(const char* msg);

// Strings
AstraString* astra_string_new(const char* data, size_t len);
AstraString* astra_string_concat(const AstraString* a,
                                  const AstraString* b);
int          astra_string_eq(const AstraString* a, const AstraString* b);
AstraString* astra_int_to_string(int64_t n);
AstraString* astra_float_to_string(double f);
AstraString* astra_bool_to_string(int b);

// Safety checks
void    astra_bounds_check(int64_t index, int64_t length,
                            const char* file, int line);
void    astra_null_check(const void* ptr, const char* name,
                          const char* file, int line);
void    astra_assert(int condition, const AstraString* message);
int64_t astra_int_div(int64_t a, int64_t b, const char* file, int line);

#endif // ASTRA_RUNTIME_H
```

---

## 12. Tracing "Hello!" from Astra to the OS

Let us follow the complete call chain when an Astra program calls `print("Hello!")`.

**Astra source code:**
```astra
fn main() {
    print("Hello!")
}
```

**Step 1: Compiler emits C code**

The Astra compiler (`astrac`) sees the `print()` call and emits:
```c
// main_generated.c (generated by astrac)
#include "runtime.h"

void astra_user_main(void) {
    // String literal "Hello!" is stored in the binary's read-only data section
    static const char _str0_data[] = "Hello!";
    static AstraString _str0 = { (char*)_str0_data, 6, 0 };

    // print("Hello!") compiles to astra_println(&_str0)
    astra_println(&_str0);
}
```

**Step 2: astra_println (runtime.c)**
```c
void astra_println(const AstraString* s) {
    astra_print(s);         // write the string
    write(STDOUT_FILENO, "\n", 1);  // write newline
}
```

**Step 3: astra_print (runtime.c)**
```c
void astra_print(const AstraString* s) {
    // s->data = "Hello!", s->len = 6
    write(STDOUT_FILENO, s->data, s->len);
    // This calls the write(2) system call
}
```

**Step 4: write(2) — OS system call**

`write(STDOUT_FILENO, "Hello!", 6)` is a Linux/macOS system call. The CPU executes a special instruction (`syscall` on x86-64) that transfers control to the kernel. The kernel:
1. Looks up file descriptor 1 (stdout).
2. Copies 6 bytes from the process's memory to the terminal's buffer.
3. Returns to the process.

The complete call chain:

```
Astra source:  print("Hello!")
                      │
Generated C:   astra_println(&_str0)
                      │
runtime.c:     astra_print(s)
                      │
runtime.c:     write(1, "Hello!", 6)
                      │
OS kernel:     copy bytes to terminal buffer
                      │
Terminal:      displays "Hello!"
```

Six layers from your Astra source to the screen. Every one of them is something you have built or understand.

---

## 13. Astra Build Milestone

At this point, the `runtime/` directory looks like:

```
runtime/
├── runtime.h        ← declarations, AstraString type
├── runtime.c        ← complete runtime (220 lines)
├── gc.h             ← GC declarations (Chapter 64)
└── gc.c             ← GC implementation (Chapter 64)
```

Let us compile the runtime and test it directly from C:

```c
// test_runtime.c — tests the runtime functions in isolation
#include "runtime.h"
#include "gc.h"
#include <stdio.h>

void astra_user_main(void) {
    // Test astra_println
    AstraString* hello = astra_string_new("Hello from the Astra runtime!", 29);
    astra_println(hello);

    // Test astra_int_to_string
    AstraString* n = astra_int_to_string(42);
    astra_println(n);

    // Test astra_string_concat
    AstraString* a = astra_string_new("foo", 3);
    AstraString* b = astra_string_new("bar", 3);
    AstraString* ab = astra_string_concat(a, b);
    astra_println(ab);  // → "foobar"

    // Test astra_bounds_check (should pass)
    astra_bounds_check(2, 5, "test_runtime.c", 30);
    astra_println(astra_string_new("bounds check passed", 19));

    // Uncomment to test panic:
    // astra_panic("something went wrong");
}
```

Build and run:
```bash
cc -c runtime/runtime.c -o runtime/runtime.o
cc -c runtime/gc.c      -o runtime/gc.o
cc -c test_runtime.c    -o test_runtime.o
cc test_runtime.o runtime/runtime.o runtime/gc.o -lm -o test_runtime
./test_runtime
```

Expected output:
```
Hello from the Astra runtime!
42
foobar
bounds check passed
```

The runtime works. Every Astra program that compiles successfully links against this and gets all these functions for free.

---

## 14. Exercises

**Exercise 63.1 — Add `astra_string_length`**

The `len` field of `AstraString` counts bytes, not Unicode characters. Write a function `size_t astra_string_char_length(const AstraString* s)` that counts UTF-8 code points (hint: a byte that starts with `10xxxxxx` is a continuation byte, not a new character).

**Exercise 63.2 — Improve panic messages**

Modify `astra_panic()` to also print the current time (using `time(NULL)` from `<time.h>`) in the panic message. Format: `[PANIC at 2024-01-15 10:30:00] message`.

**Exercise 63.3 — Add `astra_string_repeat`**

Implement `AstraString* astra_string_repeat(const AstraString* s, int64_t n)` that returns a new string containing `s` repeated `n` times. Handle `n <= 0` by returning an empty string.

**Exercise 63.4 — Test the bounds check**

Write a C test that calls `astra_bounds_check(-1, 5, "test.c", 1)` and `astra_bounds_check(5, 5, "test.c", 1)`. What happens? Verify the panic message is correct. Then test `astra_bounds_check(4, 5, "test.c", 1)` — this should NOT panic.

**Exercise 63.5 — Memory accounting**

Add a global variable `size_t astra_total_bytes_allocated = 0` to the runtime. Every time `gc_alloc` is called, add the size to this counter. Add a function `void astra_print_memory_usage(void)` that prints the total bytes allocated. Call it at the end of `_astra_start` just before exit.

**Exercise 63.6 — The C runtime's crt0**

Research what `crt0.o` is in a standard C program. How does it differ from `_astra_start`? What does `__libc_start_main` do on Linux? Write a 200-word explanation in a comment at the top of `runtime.c`.

**Exercise 63.7 — Signal handling**

When a C program dereferences a null pointer, the OS sends signal `SIGSEGV`. Add a signal handler using `sigaction()` that catches `SIGSEGV` and calls `astra_panic("null pointer dereference")` instead of letting the OS print "Segmentation fault". (Hint: look up `sigaction`, `SA_SIGACTION`, and `siginfo_t`.)

---

## Chapter Summary

| Concept | What it is | Where in code |
|---|---|---|
| Runtime | C library linked into every Astra binary | `runtime/runtime.c` |
| `_astra_start` | True entry point of every Astra program | `runtime.c:_astra_start` |
| `astra_panic` | Print error and exit(1) | `runtime.c:astra_panic` |
| `setjmp/longjmp` | Stack unwinding mechanism | `astra_panic_buf`, `setjmp` in `_astra_start` |
| `AstraString` | Internal string representation: data + len + cap | `runtime.h:AstraString` |
| `astra_println` | Calls `write(1, ...)` system call | `runtime.c:astra_print` |
| `astra_bounds_check` | Panics on out-of-bounds array access | `runtime.c:astra_bounds_check` |
| Exit code 0 | Success | `_exit(0)` in `_astra_start` |
| Exit code 1 | Any panic | `_exit(1)` in panic handler |
| Build command | `cc -c runtime.c -o runtime.o` | Linked into every Astra binary |

The runtime is done. In Chapter 64, we build the garbage collector — the system that makes memory management automatic so Astra programmers never write `free()`.
