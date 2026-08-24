# Chapter 62: A Basic Shell

> **"The shell is the face of your OS. It reads a line of text, parses it into a command and arguments, finds the program to run, and executes it. When you have a working shell, you have a working OS — something you can actually interact with, test, and demonstrate. This is the finish line of Volume 9."**

---

## Table of Contents

1. [What a Shell Does](#1-what-a-shell-does)
2. [Shell Loop — Read, Parse, Execute](#2-shell-loop--read-parse-execute)
3. [Parsing the Command Line](#3-parsing-the-command-line)
4. [Built-in Commands](#4-built-in-commands)
5. [Executing Programs](#5-executing-programs)
6. [The Complete Shell](#6-the-complete-shell)
7. [Built-in Command Implementations](#7-built-in-command-implementations)
8. [Putting It All Together](#8-putting-it-all-together)
9. [Summary](#summary)

---

## 1. What a Shell Does

```
A shell is a Read-Eval-Print Loop (REPL):

  while (running) {
    print prompt
    read a line of text from keyboard
    parse the line into: command + arguments
    
    if command is built-in (echo, ls, ps, help, exit):
      execute it directly in the shell
    else:
      look up the command in the file system (/bin/command)
      create a new process to run it
      wait for it to finish
    
    go back to start
  }

Built-in commands (run inside the shell process, no child):
  help, exit, echo, clear, ps, ls, cat, time, meminfo

External commands (run as a child process):
  /bin/hello, /bin/calc, any user-mode binary
  (In our OS these are ramdisk programs)
```

---

## 2. Shell Loop — Read, Parse, Execute

```c
/* kernel/shell.c */

#include "shell.h"
#include "keyboard.h"
#include "vga.h"
#include "process.h"
#include "scheduler.h"
#include "timer.h"
#include "pmm.h"
#include "vfs.h"
#include "string.h"
#include "heap.h"

#define MAX_LINE    256
#define MAX_ARGS    16
#define MAX_ARG_LEN 64

static int shell_running = 1;

void shell_main(void) {
    terminal_set_color(VGA_COLOR_LIGHT_GREEN, VGA_COLOR_BLACK);
    kprintf("\n");
    kprintf("  ████████╗██╗███╗   ██╗██╗   ██╗ ██████╗ ███████╗\n");
    kprintf("     ██╔══╝██║████╗  ██║╚██╗ ██╔╝██╔═══██╗██╔════╝\n");
    kprintf("     ██║   ██║██╔██╗ ██║ ╚████╔╝ ██║   ██║███████╗\n");
    kprintf("     ██║   ██║██║╚██╗██║  ╚██╔╝  ██║   ██║╚════██║\n");
    kprintf("     ██║   ██║██║ ╚████║   ██║   ╚██████╔╝███████║\n");
    kprintf("     ╚═╝   ╚═╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ ╚══════╝\n");
    terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
    kprintf("\nTinyOS Shell v1.0  (type 'help' for commands)\n\n");
    
    char line[MAX_LINE];
    
    while (shell_running) {
        /* Print prompt: */
        terminal_set_color(VGA_COLOR_LIGHT_CYAN, VGA_COLOR_BLACK);
        kprintf("tinyos");
        terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
        kprintf(":/ $ ");
        terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
        
        /* Read a line: */
        int len = kbd_readline(line, MAX_LINE);
        if (len == 0) continue;
        
        /* Parse and execute: */
        shell_execute(line);
    }
    
    kprintf("Shell exiting.\n");
    process_exit(0);
}
```

---

## 3. Parsing the Command Line

```c
/* Parse a command line into argv[]: */
static int parse_args(char *line, char *argv[], int max_args) {
    int argc = 0;
    char *p  = line;
    
    while (*p && argc < max_args) {
        /* Skip whitespace: */
        while (*p == ' ' || *p == '\t') p++;
        if (!*p) break;
        
        /* Start of argument: */
        argv[argc++] = p;
        
        /* Find end of argument: */
        while (*p && *p != ' ' && *p != '\t') p++;
        
        /* Null-terminate it: */
        if (*p) {
            *p = '\0';
            p++;
        }
    }
    
    return argc;
}

void shell_execute(char *line) {
    char *argv[MAX_ARGS];
    int   argc = parse_args(line, argv, MAX_ARGS);
    
    if (argc == 0) return;
    
    /* Try built-in commands first: */
    if (shell_builtin(argc, argv)) return;
    
    /* Try to execute as a file from ramdisk: */
    /* In our OS: look up in ramdisk by name, load+run as user process */
    shell_run_file(argv[0], argc, argv);
}
```

---

## 4. Built-in Commands

```c
/* Built-in command handler:
   Returns 1 if the command was a built-in (handled), 0 if not found. */

typedef int (*builtin_fn)(int argc, char *argv[]);

typedef struct {
    const char *name;
    builtin_fn   fn;
    const char  *help;
} builtin_t;

/* Forward declarations: */
static int cmd_help(int argc, char *argv[]);
static int cmd_exit(int argc, char *argv[]);
static int cmd_echo(int argc, char *argv[]);
static int cmd_clear(int argc, char *argv[]);
static int cmd_ps(int argc, char *argv[]);
static int cmd_ls(int argc, char *argv[]);
static int cmd_cat(int argc, char *argv[]);
static int cmd_time(int argc, char *argv[]);
static int cmd_meminfo(int argc, char *argv[]);
static int cmd_reboot(int argc, char *argv[]);

static const builtin_t builtins[] = {
    { "help",    cmd_help,    "Show this help message" },
    { "exit",    cmd_exit,    "Exit the shell" },
    { "echo",    cmd_echo,    "Print arguments to screen" },
    { "clear",   cmd_clear,   "Clear the screen" },
    { "ps",      cmd_ps,      "Show process table" },
    { "ls",      cmd_ls,      "List files in ramdisk" },
    { "cat",     cmd_cat,     "Print file contents" },
    { "time",    cmd_time,    "Show system time (timer ticks)" },
    { "meminfo", cmd_meminfo, "Show memory usage" },
    { "reboot",  cmd_reboot,  "Reboot the system" },
    { NULL, NULL, NULL }
};

int shell_builtin(int argc, char *argv[]) {
    for (int i = 0; builtins[i].name; i++) {
        if (strcmp(argv[0], builtins[i].name) == 0) {
            builtins[i].fn(argc, argv);
            return 1;
        }
    }
    return 0;  /* Not a built-in */
}
```

---

## 5. Executing Programs

```c
/* Run a file from the ramdisk as a new process: */
void shell_run_file(const char *name, int argc, char *argv[]) {
    /* Look up the file: */
    vfs_node_t *node = vfs_lookup(name);
    
    if (!node) {
        /* Try with leading /bin/ prefix: */
        char path[80];
        path[0] = '/';
        for (int i = 0; name[i] && i < 74; i++) path[i+1] = name[i];
        path[strlen(name)+1] = '\0';
        node = vfs_lookup(path);
    }
    
    if (!node || node->type != VFS_NODE_FILE) {
        kprintf("Command not found: %s\n", name);
        return;
    }
    
    /* Read the file into a buffer: */
    uint8_t *binary = (uint8_t *)kmalloc(node->size);
    if (!binary) {
        kprintf("Out of memory loading '%s'\n", name);
        return;
    }
    
    if (node->ops && node->ops->read) {
        node->ops->read(node, 0, binary, node->size);
    }
    
    /* Launch as a user process: */
    process_t *child = launch_user_process(name, binary, node->size);
    kfree(binary);
    
    if (!child) {
        kprintf("Failed to launch '%s'\n", name);
        return;
    }
    
    /* Wait for child to finish (simplified blocking wait): */
    while (child->state != PROC_ZOMBIE && child->state != PROC_DEAD) {
        yield();  /* Give up CPU while waiting */
    }
    
    if (child->state == PROC_ZOMBIE) {
        kprintf("[Process exited with code %d]\n", child->exit_code);
        child->state = PROC_DEAD;  /* Reap the zombie */
    }
}
```

---

## 6. The Complete Shell

```c
/* shell.h */
#pragma once

void shell_main(void);
void shell_execute(char *line);
int  shell_builtin(int argc, char *argv[]);
void shell_run_file(const char *name, int argc, char *argv[]);
```

---

## 7. Built-in Command Implementations

```c
static int cmd_help(int argc, char *argv[]) {
    (void)argc; (void)argv;
    terminal_set_color(VGA_COLOR_YELLOW, VGA_COLOR_BLACK);
    kprintf("\nTinyOS Shell Commands:\n\n");
    terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
    for (int i = 0; builtins[i].name; i++) {
        terminal_set_color(VGA_COLOR_LIGHT_CYAN, VGA_COLOR_BLACK);
        kprintf("  %-10s", builtins[i].name);
        terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLACK);
        kprintf("  %s\n", builtins[i].help);
    }
    kprintf("\n");
    return 1;
}

static int cmd_exit(int argc, char *argv[]) {
    (void)argc; (void)argv;
    shell_running = 0;
    return 1;
}

static int cmd_echo(int argc, char *argv[]) {
    for (int i = 1; i < argc; i++) {
        kprintf("%s", argv[i]);
        if (i < argc - 1) kprintf(" ");
    }
    kprintf("\n");
    return 1;
}

static int cmd_clear(int argc, char *argv[]) {
    (void)argc; (void)argv;
    terminal_init();  /* Reinitialize = clear screen + reset cursor */
    return 1;
}

static int cmd_ps(int argc, char *argv[]) {
    (void)argc; (void)argv;
    process_print_all();
    return 1;
}

static int cmd_ls(int argc, char *argv[]) {
    (void)argc; (void)argv;
    kprintf("\nFiles in ramdisk:\n");
    /* List all ramdisk files: */
    extern int rd_count;
    extern vfs_node_t rd_nodes[];
    for (int i = 0; i < rd_count; i++) {
        kprintf("  %-20s  %u bytes\n", rd_nodes[i].name, rd_nodes[i].size);
    }
    kprintf("\n");
    return 1;
}

static int cmd_cat(int argc, char *argv[]) {
    if (argc < 2) {
        kprintf("Usage: cat <filename>\n");
        return 1;
    }
    
    int fd = vfs_open(argv[1], O_RDONLY);
    if (fd < 0) {
        kprintf("cat: cannot open '%s'\n", argv[1]);
        return 1;
    }
    
    uint8_t buf[256];
    int32_t n;
    while ((n = vfs_read(fd, buf, sizeof(buf))) > 0) {
        for (int i = 0; i < n; i++) {
            terminal_putchar((char)buf[i]);
        }
    }
    vfs_close(fd);
    return 1;
}

static int cmd_time(int argc, char *argv[]) {
    (void)argc; (void)argv;
    uint32_t ticks = timer_get_ticks();
    kprintf("Timer ticks: %u  (uptime: %u seconds at 100Hz)\n",
            ticks, ticks / 100);
    return 1;
}

static int cmd_meminfo(int argc, char *argv[]) {
    (void)argc; (void)argv;
    uint32_t total = pmm_total_frame_count() * PAGE_SIZE / 1024;
    uint32_t free  = pmm_free_frame_count()  * PAGE_SIZE / 1024;
    uint32_t used  = total - free;
    kprintf("\nMemory Information:\n");
    kprintf("  Total:  %u KB (%u MB)\n", total, total / 1024);
    kprintf("  Used:   %u KB (%u MB)\n", used,  used  / 1024);
    kprintf("  Free:   %u KB (%u MB)\n", free,  free  / 1024);
    kprintf("  Usage:  %u%%\n\n", used * 100 / total);
    return 1;
}

static int cmd_reboot(int argc, char *argv[]) {
    (void)argc; (void)argv;
    kprintf("Rebooting...\n");
    /* Trigger triple fault (the cleanest reboot in a bare-metal OS): */
    __asm__ volatile(
        "cli\n"           /* disable interrupts */
        "lidt [0]\n"      /* load invalid IDT */
        "int $0\n"        /* trigger exception → triple fault → CPU reset */
    );
    for (;;) {}
    return 1;
}
```

---

## 8. Putting It All Together

Final `kernel_main`:

```c
void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    
    kprintf("TinyOS Booting...\n");
    
    struct multiboot_info *mbi = (struct multiboot_info *)mbi_ptr;
    if (magic != 0x2BADB002) {
        kprintf("FATAL: Bad magic 0x%x\n", magic);
        for (;;) {}
    }
    
    kprintf("[1/9] GDT...      "); gdt_init();                        kprintf("OK\n");
    kprintf("[2/9] IDT...      "); idt_init();                        kprintf("OK\n");
    kprintf("[3/9] PIC...      "); pic_init(); pic_disable();         kprintf("OK\n");
    kprintf("[4/9] PMM...      "); pmm_init(mbi->mmap_addr, mbi->mmap_length, mbi->mem_upper);
    kprintf("[5/9] VMM...      "); vmm_init();                        kprintf("OK\n");
    kprintf("[6/9] Heap...     "); heap_init();                       kprintf("OK\n");
    kprintf("[7/9] Processes..."); process_init();                    kprintf("OK\n");
    kprintf("[8/9] VFS...      "); ramdisk_init();                    kprintf("OK\n");
    kprintf("[9/9] Timer+Kbd..."); timer_init(100); keyboard_init();  kprintf("OK\n");
    
    kprintf("\nAll systems initialized. Starting shell...\n\n");
    
    /* Create the shell process: */
    process_create("shell", shell_main, 8);
    
    /* Start scheduling: */
    __asm__ volatile("sti");
    schedule();
    
    /* Should never reach: */
    for (;;) __asm__ volatile("hlt");
}
```

**Running TinyOS:**
```bash
make clean && make && make run
```

You should see:
```
TinyOS Booting...
[1/9] GDT...      OK
[2/9] IDT...      OK
[3/9] PIC...      OK
[4/9] PMM...      128 MB total, 121 MB free
[5/9] VMM...      OK
[6/9] Heap...     OK
[7/9] Processes...OK
[8/9] VFS...      Ramdisk initialized with 2 files. OK
[9/9] Timer+Kbd...OK

All systems initialized. Starting shell...

  ████████╗██╗███╗   ██╗██╗   ██╗ ██████╗ ███████╗
     ██╔══╝ ...
     
TinyOS Shell v1.0  (type 'help' for commands)

tinyos:/ $ help
tinyos:/ $ ls
tinyos:/ $ cat motd
tinyos:/ $ meminfo
tinyos:/ $ ps
tinyos:/ $ time
tinyos:/ $ echo Hello World!
tinyos:/ $ exit
```

---

## Summary

| Concept | Description |
|---------|------------|
| Shell | REPL: read line → parse → execute → repeat |
| Parse | Split input by spaces; argv[0]=command, argv[1..]=args |
| Built-in | Command handled directly inside the shell process (no child created) |
| External command | Loaded from VFS, executed as a new user-mode process |
| shell_run_file | Look up file in VFS → read binary → launch_user_process → wait |
| Waiting for child | Shell calls yield() in a loop until child state = PROC_ZOMBIE |
| help | List all built-in commands with descriptions |
| ps | Call process_print_all() to show all PCBs |
| ls | List VFS (ramdisk) nodes |
| cat | Open file via VFS, read chunks, print to terminal |
| meminfo | Show PMM statistics: total/used/free frames |
| time | Show timer_get_ticks() and uptime in seconds |
| reboot | Trigger triple fault by loading an invalid IDT |
| shell_running | Flag: set to 0 by 'exit' command to end the shell loop |
