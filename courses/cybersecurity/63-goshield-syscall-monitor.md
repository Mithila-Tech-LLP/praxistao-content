# Chapter 63: GoShield — Syscall Monitor

*The syscall monitor is GoShield's kernel-level visibility layer. By observing system calls, we can detect attacks that would be invisible to file-based or network-based monitoring alone.*

---

## Why Syscall Monitoring?

Every attack eventually makes syscalls:
- Shellcode → `execve`, `mmap`, `mprotect`
- Process injection → `ptrace`, anonymous RWX `mmap`
- Credential theft → `openat` on `/etc/shadow`, `/proc/*/mem`
- Reverse shell → `socket`, `connect`, `dup2`, `execve`
- Rootkit → `write` to `/proc/sys/kernel/`, kernel module load

Syscalls are the single chokepoint between user-space attacks and the OS.

---

## Linux Syscall Monitoring Approaches

```
1. auditd — Linux kernel audit framework
   + Production-grade, stable
   + Rule-based filtering
   - Overhead on busy systems

2. eBPF (Extended Berkeley Packet Filter)
   + Modern, low overhead
   + Rich context (arguments, return values)
   + Production safe
   - Requires Linux 4.1+, complex to write

3. strace / ptrace
   + Easy to use
   - Only traces specific process
   - Not suitable for always-on monitoring

4. Seccomp profiles
   + Block syscalls before they execute
   + No overhead when not violated
   - Limits what process can do (used in containers)
```

---

## auditd-Based Monitoring

```bash
# Install auditd
apt install auditd

# Start and enable
systemctl enable --now auditd

# Add rules to watch critical operations:
# Watch execve (all program executions)
auditctl -a always,exit -F arch=b64 -S execve -k exec_watch

# Watch /etc/passwd and /etc/shadow
auditctl -w /etc/passwd -p wa -k passwd_changes
auditctl -w /etc/shadow -p wa -k shadow_access

# Watch loading kernel modules
auditctl -a always,exit -F arch=b64 -S init_module -S finit_module -k module_load

# Watch ptrace (injection attempts)
auditctl -a always,exit -F arch=b64 -S ptrace -k ptrace_detect

# Make rules permanent
cat >> /etc/audit/rules.d/goshield.rules << 'EOF'
-a always,exit -F arch=b64 -S execve -k exec_watch
-w /etc/passwd -p wa -k passwd_changes
-w /etc/shadow -p wa -k shadow_access
-a always,exit -F arch=b64 -S init_module -S finit_module -k module_load
-a always,exit -F arch=b64 -S ptrace -k ptrace_detect
EOF

# View audit logs
ausearch -k exec_watch
ausearch -k ptrace_detect
aureport --summary
```

---

## Go: Parsing Audit Logs

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "time"
)

type AuditEvent struct {
    Time     time.Time
    Type     string
    Syscall  string
    Comm     string  // process name
    Exe      string  // executable path
    Pid      string
    UID      string
    Args     map[string]string
    Key      string  // rule key that triggered
}

func parseAuditLine(line string) *AuditEvent {
    event := &AuditEvent{Args: make(map[string]string)}
    
    // Extract type: type=SYSCALL
    typeRe := regexp.MustCompile(`type=(\w+)`)
    if m := typeRe.FindStringSubmatch(line); len(m) > 1 {
        event.Type = m[1]
    }
    
    // Extract key=value pairs
    kvRe := regexp.MustCompile(`(\w+)=([^\s]+)`)
    for _, m := range kvRe.FindAllStringSubmatch(line, -1) {
        key, val := m[1], m[2]
        // Strip quotes
        val = strings.Trim(val, `"`)
        event.Args[key] = val
        
        switch key {
        case "comm":
            event.Comm = val
        case "exe":
            event.Exe = val
        case "pid":
            event.Pid = val
        case "uid":
            event.UID = val
        case "key":
            event.Key = val
        case "syscall":
            event.Syscall = val
        }
    }
    
    return event
}

type SyscallDetector struct {
    suspiciousPatterns []DetectionPattern
}

type DetectionPattern struct {
    Name     string
    Severity string
    Check    func(*AuditEvent) bool
}

func NewSyscallDetector() *SyscallDetector {
    d := &SyscallDetector{}
    
    d.suspiciousPatterns = []DetectionPattern{
        {
            Name:     "Ptrace on non-child process",
            Severity: "HIGH",
            Check: func(e *AuditEvent) bool {
                return e.Key == "ptrace_detect" && e.Syscall == "ptrace"
            },
        },
        {
            Name:     "Shadow file accessed",
            Severity: "CRITICAL",
            Check: func(e *AuditEvent) bool {
                return e.Key == "shadow_access"
            },
        },
        {
            Name:     "Kernel module loaded",
            Severity: "CRITICAL",
            Check: func(e *AuditEvent) bool {
                return e.Key == "module_load"
            },
        },
        {
            Name:     "Suspicious process spawned from web server",
            Severity: "HIGH",
            Check: func(e *AuditEvent) bool {
                webServers := []string{"apache2", "nginx", "httpd", "php-fpm", "uwsgi"}
                suspicious := []string{"bash", "sh", "python3", "perl", "ruby"}
                
                commLow := strings.ToLower(e.Comm)
                for _, ws := range webServers {
                    if strings.Contains(commLow, ws) {
                        for _, susp := range suspicious {
                            if strings.Contains(e.Exe, susp) {
                                return true
                            }
                        }
                    }
                }
                return false
            },
        },
    }
    
    return d
}

func (d *SyscallDetector) Analyze(event *AuditEvent) []string {
    var alerts []string
    for _, pattern := range d.suspiciousPatterns {
        if pattern.Check(event) {
            alerts = append(alerts, fmt.Sprintf("[%s] %s: pid=%s exe=%s",
                pattern.Severity, pattern.Name, event.Pid, event.Exe))
        }
    }
    return alerts
}

func MonitorAuditLog(path string) {
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("Cannot open audit log: %v\n", err)
        return
    }
    defer f.Close()
    
    detector := NewSyscallDetector()
    
    // Tail file (simplified — real implementation would seek to end)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        event := parseAuditLine(line)
        
        if alerts := detector.Analyze(event); len(alerts) > 0 {
            for _, alert := range alerts {
                fmt.Printf("[ALERT] %s\n", alert)
            }
        }
    }
}

func main() {
    fmt.Println("GoShield Syscall Monitor starting...")
    MonitorAuditLog("/var/log/audit/audit.log")
}
```

---

## eBPF-Based Monitoring (Overview)

For production GoShield, eBPF gives lower overhead and richer data:

```go
// Using cilium/ebpf library
// go get github.com/cilium/ebpf

// eBPF program (C, compiled to BPF bytecode)
// bpf_program.c:
/*
SEC("tracepoint/syscalls/sys_enter_execve")
int trace_execve(struct trace_event_raw_sys_enter* ctx) {
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(comm, sizeof(comm));
    
    // Send event to user space via ring buffer
    bpf_ringbuf_output(...);
    return 0;
}
*/

// Go side reads from ring buffer
// reader, _ := ringbuf.NewReader(objs.Events)
// for {
//     record, _ := reader.Read()
//     // process event
// }
```

---

## Seccomp Profiles (Blocking Dangerous Syscalls)

```go
import "golang.org/x/sys/unix"

// Create a seccomp filter for a sandboxed process
func restrictSyscalls() error {
    // Install seccomp filter that blocks dangerous syscalls
    // Uses BPF filter language
    
    // In real GoShield: apply to the agent process to reduce attack surface
    // If agent is compromised, it can't ptrace, load modules, etc.
    
    // Using libseccomp Go binding:
    // filter, _ := seccomp.NewFilter(seccomp.ActKill)
    // filter.AddRule(seccomp.ScmpSyscall("ptrace"), seccomp.ActAllow)  // selectively allow
    
    return nil
}
```

---

## Summary

| Mechanism | Overhead | Use case |
|-----------|----------|---------|
| auditd | Medium | Production logging, rule-based alerts |
| eBPF | Low | Real-time, high-frequency monitoring |
| strace/ptrace | High | Debugging specific processes |
| Seccomp | Minimal | Preventing dangerous syscalls |

---

## Exercises

1. Install auditd on a Linux VM and add rules for execve and ptrace monitoring
2. Trigger each rule by running programs — verify the audit log captures them
3. Write the Go audit log parser and test it detects a mock "shadow file access" event
4. Research eBPF tracing tools (BCC project, Falco) — how does Falco detect container escapes?
