# Chapter 13: Go System Interaction — OS, Exec, File I/O, Syscalls

*Security tools need to reach deep into the operating system — run commands, read files, inspect processes, make syscalls. Go's standard library makes this straightforward while keeping the code safe and portable.*

---

## Running System Commands

```go
package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"
    "context"
)

// Run a command, capture output
func run(name string, args ...string) (string, string, int) {
    cmd := exec.Command(name, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    code := 0
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            code = exitErr.ExitCode()
        }
    }
    return stdout.String(), stderr.String(), code
}

// Run with timeout
func runWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    out, err := exec.CommandContext(ctx, name, args...).Output()
    return string(out), err
}

// DANGER: Never do this with user input!
func badCommand(userInput string) string {
    // COMMAND INJECTION — userInput could be "foo; rm -rf /"
    out, _ := exec.Command("sh", "-c", "grep "+userInput+" /etc/hosts").Output()
    return string(out)
}

// SAFE: Pass args separately, never through shell
func safeCommand(pattern string) string {
    out, _ := exec.Command("grep", pattern, "/etc/hosts").Output()
    return string(out)
}

func main() {
    // Get system info
    out, _, _ := run("uname", "-a")
    fmt.Println("System:", strings.TrimSpace(out))
    
    // List processes
    out, _, _ = run("ps", "aux")
    lines := strings.Split(out, "\n")
    fmt.Printf("Running processes: %d\n", len(lines))
    
    // With timeout
    out, err := runWithTimeout(5*time.Second, "nmap", "-sV", "127.0.0.1")
    if err != nil {
        fmt.Println("Timed out or error:", err)
    } else {
        fmt.Println(out)
    }
}
```

---

## File I/O for Security Tools

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// Read file line by line (memory-efficient for large files)
func readLines(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    var lines []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if line != "" && !strings.HasPrefix(line, "#") {
            lines = append(lines, line)
        }
    }
    return lines, scanner.Err()
}

// Walk directory tree — find interesting files
func walkForSensitiveFiles(root string) {
    sensitiveExts := map[string]bool{
        ".key": true, ".pem": true, ".p12": true,
        ".env": true, ".conf": true, ".cfg": true,
        ".bak": true, ".sql": true, ".db":  true,
    }
    sensitiveNames := map[string]bool{
        "id_rsa": true, "id_ed25519": true,
        ".htpasswd": true, "shadow": true,
        "docker-compose.yml": true, "secrets.yaml": true,
    }
    
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        
        name := filepath.Base(path)
        ext  := filepath.Ext(path)
        
        if sensitiveExts[ext] || sensitiveNames[name] {
            fmt.Printf("[SENSITIVE] %s (%d bytes)\n", path, info.Size())
        }
        return nil
    })
}

// Write results atomically (temp file + rename)
func writeResultsAtomic(path string, data []byte) error {
    tmpFile := path + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0600); err != nil {
        return err
    }
    return os.Rename(tmpFile, path)  // atomic on most filesystems
}

// Check if path is readable/writable by current user
func checkPermissions(path string) {
    // Read
    f, err := os.Open(path)
    if err == nil {
        f.Close()
        fmt.Printf("[READABLE] %s\n", path)
    }
    
    // Write
    f, err = os.OpenFile(path, os.O_WRONLY, 0)
    if err == nil {
        f.Close()
        fmt.Printf("[WRITABLE] %s ← potential security risk\n", path)
    }
}

func main() {
    // Load wordlist
    words, err := readLines("/usr/share/wordlists/rockyou.txt")
    if err != nil {
        fmt.Println("No wordlist:", err)
    } else {
        fmt.Printf("Loaded %d words\n", len(words))
    }
    
    // Hunt for sensitive files
    fmt.Println("\nSearching /home for sensitive files:")
    walkForSensitiveFiles("/home")
    
    // Check critical file permissions
    for _, path := range []string{"/etc/passwd", "/etc/shadow", "/etc/sudoers"} {
        checkPermissions(path)
    }
    
    _ = io.Discard
}
```

---

## Environment Variables

```go
package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    // Get a single env var
    apiKey := os.Getenv("GOSHIELD_API_KEY")
    if apiKey == "" {
        fmt.Fprintln(os.Stderr, "GOSHIELD_API_KEY not set")
        os.Exit(1)
    }
    
    // Get with default
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
        logLevel = "info"
    }
    
    // Dump all environment variables (useful for recon)
    fmt.Println("Environment variables:")
    for _, env := range os.Environ() {
        parts := strings.SplitN(env, "=", 2)
        key := parts[0]
        val := parts[1]
        
        // Flag sensitive-looking variables
        sensitive := false
        for _, word := range []string{"PASSWORD", "SECRET", "KEY", "TOKEN", "CREDENTIAL"} {
            if strings.Contains(strings.ToUpper(key), word) {
                sensitive = true
                break
            }
        }
        
        if sensitive {
            fmt.Printf("  [SENSITIVE] %s = %s...\n", key, val[:min(len(val), 4)])
        }
    }
    
    // Set/unset
    os.Setenv("MY_TOOL_DEBUG", "true")
    os.Unsetenv("MY_TOOL_DEBUG")
    
    // Read /proc/PID/environ for another process (if readable)
    pid := os.Getpid()
    data, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
    if err == nil {
        // Env vars are NUL-separated
        for _, entry := range strings.Split(string(data), "\x00") {
            if strings.Contains(strings.ToUpper(entry), "PATH") {
                fmt.Println("PATH env:", entry)
            }
        }
    }
}

func min(a, b int) int { if a < b { return a }; return b }
```

---

## OS-Level Information

```go
package main

import (
    "fmt"
    "os"
    "os/user"
    "runtime"
)

func systemInfo() {
    fmt.Println("=== System Info ===")
    fmt.Println("OS:", runtime.GOOS)
    fmt.Println("Arch:", runtime.GOARCH)
    fmt.Println("CPUs:", runtime.NumCPU())
    
    hostname, _ := os.Hostname()
    fmt.Println("Hostname:", hostname)
    
    fmt.Println("PID:", os.Getpid())
    fmt.Println("PPID:", os.Getppid())
    fmt.Println("UID:", os.Getuid())
    fmt.Println("GID:", os.Getgid())
    fmt.Println("EUID:", os.Geteuid())
    
    // Current user
    u, err := user.Current()
    if err == nil {
        fmt.Printf("User: %s (home: %s)\n", u.Username, u.HomeDir)
    }
    
    // Groups
    groups, _ := os.Getgroups()
    fmt.Println("Groups:", groups)
    
    // Check if running as root
    if os.Geteuid() == 0 {
        fmt.Println("[!] Running as root")
    }
    
    // Working directory
    wd, _ := os.Getwd()
    fmt.Println("Working dir:", wd)
    
    // Temp directory
    fmt.Println("Temp dir:", os.TempDir())
}

func main() {
    systemInfo()
}
```

---

## Signals — Graceful Shutdown

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Set up signal handler
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh,
        os.Interrupt,     // Ctrl+C (SIGINT)
        syscall.SIGTERM,  // kill (default)
        syscall.SIGHUP,   // terminal closed
    )
    
    // Do work
    fmt.Println("Tool running... Ctrl+C to stop")
    
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case sig := <-sigCh:
            fmt.Printf("\nReceived signal: %v\n", sig)
            fmt.Println("Cleaning up...")
            // close connections, flush buffers, save state
            fmt.Println("Done.")
            return
        case t := <-ticker.C:
            fmt.Printf("Working... %v\n", t.Unix())
        }
    }
}
```

---

## Summary

| Task | Go approach |
|------|-----------|
| Run command | `exec.Command(name, args...).Output()` |
| Run with timeout | `exec.CommandContext(ctx, ...)` |
| Read file lines | `bufio.Scanner` over `os.Open` |
| Walk directories | `filepath.Walk` |
| Get env var | `os.Getenv` |
| Get UID/GID | `os.Getuid()`, `os.Getgid()` |
| Graceful shutdown | `signal.Notify` + channel select |
| Atomic write | Write temp file + `os.Rename` |

---

## Exercises

1. Write a tool that reads `/proc/PID/maps` for a given PID and flags any anonymous executable memory regions (no file backing, executable = suspicious)
2. Build a credential scanner: walk a directory tree and search all `.env`, `.cfg`, `.yaml` files for patterns matching `password=`, `secret=`, `api_key=`
3. Write a tool that reads `/etc/passwd` and `/etc/shadow` and reports which users have passwordless logins or old password hashes
4. Implement a file watcher without `fsnotify`: poll a directory every second, detect new files, report them
