# Chapter 60: GoShield Agent — File Integrity Monitoring

*File Integrity Monitoring (FIM) watches for unauthorized changes to critical files. This is how EDRs detect ransomware, webshell drops, and configuration tampering. We'll build a production-quality FIM module.*

---

## What File Events to Watch

**Why files matter for detection:**
- **Ransomware:** Mass writes to files with encrypted extensions
- **Webshells:** New `.php/.jsp/.asp` files in web directories
- **Persistence:** New cron jobs, startup scripts, SSH keys
- **Credential theft:** Reading `/etc/shadow`, `id_rsa` files
- **Data exfiltration:** Large file reads, data staging in `/tmp`
- **Config tampering:** Writes to `/etc/passwd`, `/etc/sudoers`
- **Malware drop:** New executables in unusual locations

---

## The `fsnotify` Library

Go's `fsnotify` library provides cross-platform file system notifications (uses inotify on Linux, FSEvents on Mac, ReadDirectoryChangesW on Windows):

```bash
go get github.com/fsnotify/fsnotify
```

Events it provides:
- `Create` — file or directory created
- `Write` — file written to
- `Remove` — file deleted
- `Rename` — file renamed
- `Chmod` — permissions changed

---

## Building the File Watcher

```go
// pkg/agent/filewatcher.go
package agent

import (
    "crypto/sha256"
    "encoding/hex"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
    
    "github.com/fsnotify/fsnotify"
    "github.com/yourname/goshield/pkg/events"
)

// FileWatcher monitors file system for changes
type FileWatcher struct {
    agentID  string
    hostname string
    watcher  *fsnotify.Watcher
    config   FileWatchConfig
    eventCh  chan<- interface{}  // Output channel — events go here
    stopCh   chan struct{}
    wg       sync.WaitGroup
    
    // Track recent events to deduplicate
    recentEvents map[string]time.Time
    mu           sync.Mutex
}

// FileWatchConfig configures what to watch
type FileWatchConfig struct {
    Paths          []string `yaml:"paths"`
    ExcludePatterns []string `yaml:"exclude_patterns"`
    TrackContent   bool     `yaml:"track_content"`    // Compute SHA256?
    MaxFileSize    int64    `yaml:"max_file_size"`    // Don't hash files larger than this
}

// DefaultFileWatchConfig returns sensible defaults
func DefaultFileWatchConfig() FileWatchConfig {
    return FileWatchConfig{
        Paths: []string{
            "/etc",
            "/bin",
            "/sbin",
            "/usr/bin",
            "/usr/sbin",
            "/usr/local/bin",
            "/home",
            "/root",
            "/tmp",
            "/var/tmp",
            "/var/www",   // Web files
            "/var/spool/cron",  // Cron jobs
        },
        ExcludePatterns: []string{
            "*.log",
            "*.pid",
            "*.sock",
            "*.swp",
            "*~",
        },
        TrackContent: true,
        MaxFileSize:  10 * 1024 * 1024, // 10MB max for hashing
    }
}

// NewFileWatcher creates and starts a file watcher
func NewFileWatcher(agentID, hostname string, config FileWatchConfig, eventCh chan<- interface{}) (*FileWatcher, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }
    
    fw := &FileWatcher{
        agentID:      agentID,
        hostname:     hostname,
        watcher:      watcher,
        config:       config,
        eventCh:      eventCh,
        stopCh:       make(chan struct{}),
        recentEvents: make(map[string]time.Time),
    }
    
    // Add all configured paths
    for _, path := range config.Paths {
        if err := fw.addPathRecursive(path); err != nil {
            // Log but don't fail — path might not exist on all systems
            // log.Printf("Warning: cannot watch %s: %v", path, err)
        }
    }
    
    return fw, nil
}

// addPathRecursive adds a directory and all subdirectories to watch
func (fw *FileWatcher) addPathRecursive(root string) error {
    return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil // Skip inaccessible paths
        }
        if info.IsDir() {
            return fw.watcher.Add(path)
        }
        return nil
    })
}

// Start begins watching for file system events
func (fw *FileWatcher) Start() {
    fw.wg.Add(1)
    go fw.processEvents()
}

// Stop gracefully stops the watcher
func (fw *FileWatcher) Stop() {
    close(fw.stopCh)
    fw.watcher.Close()
    fw.wg.Wait()
}

// processEvents is the main event loop
func (fw *FileWatcher) processEvents() {
    defer fw.wg.Done()
    
    // Cleanup deduplication cache every minute
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-fw.stopCh:
            return
            
        case event, ok := <-fw.watcher.Events:
            if !ok {
                return
            }
            fw.handleEvent(event)
            
        case err, ok := <-fw.watcher.Errors:
            if !ok {
                return
            }
            _ = err // log in production
            
        case <-ticker.C:
            fw.cleanupRecentCache()
        }
    }
}

// handleEvent processes a single fsnotify event
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
    // Skip excluded patterns
    if fw.isExcluded(event.Name) {
        return
    }
    
    // Deduplicate rapid writes (e.g., app writing every 100ms)
    key := event.Name + event.Op.String()
    fw.mu.Lock()
    if last, seen := fw.recentEvents[key]; seen && time.Since(last) < 500*time.Millisecond {
        fw.mu.Unlock()
        return // Skip duplicate within 500ms
    }
    fw.recentEvents[key] = time.Now()
    fw.mu.Unlock()
    
    // Determine action
    action := ""
    switch {
    case event.Op&fsnotify.Create != 0:
        action = "create"
        // If new directory, add it to watcher
        if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
            fw.watcher.Add(event.Name)
        }
    case event.Op&fsnotify.Write != 0:
        action = "write"
    case event.Op&fsnotify.Remove != 0:
        action = "delete"
    case event.Op&fsnotify.Rename != 0:
        action = "rename"
    case event.Op&fsnotify.Chmod != 0:
        action = "chmod"
    default:
        return
    }
    
    // Build the file event
    fileEvent := &events.FileEvent{
        Base: events.Base{
            ID:        events.NewID(),
            AgentID:   fw.agentID,
            Hostname:  fw.hostname,
            Timestamp: time.Now(),
            Type:      events.EventTypeFile,
        },
        Action:    action,
        Path:      event.Name,
        Extension: strings.ToLower(filepath.Ext(event.Name)),
    }
    
    // Get file metadata (may fail if file was deleted)
    if action != "delete" {
        if info, err := os.Stat(event.Name); err == nil {
            fileEvent.Size = info.Size()
        }
        
        // Compute SHA256 hash if file is small enough
        if fw.config.TrackContent && fileEvent.Size > 0 && fileEvent.Size < fw.config.MaxFileSize {
            if hash, err := hashFile(event.Name); err == nil {
                fileEvent.SHA256 = hash
            }
        }
        
        // Check if file is hidden (starts with dot on Unix)
        base := filepath.Base(event.Name)
        fileEvent.IsHidden = len(base) > 1 && base[0] == '.'
    }
    
    // Send event downstream
    select {
    case fw.eventCh <- fileEvent:
    default:
        // Channel full — in production, log dropped event
    }
}

// isExcluded checks if a path matches any exclude pattern
func (fw *FileWatcher) isExcluded(path string) bool {
    base := filepath.Base(path)
    for _, pattern := range fw.config.ExcludePatterns {
        if matched, _ := filepath.Match(pattern, base); matched {
            return true
        }
    }
    return false
}

// cleanupRecentCache removes old entries from deduplication cache
func (fw *FileWatcher) cleanupRecentCache() {
    fw.mu.Lock()
    defer fw.mu.Unlock()
    cutoff := time.Now().Add(-1 * time.Minute)
    for key, t := range fw.recentEvents {
        if t.Before(cutoff) {
            delete(fw.recentEvents, key)
        }
    }
}

// hashFile computes the SHA256 hash of a file
func hashFile(path string) (string, error) {
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

## Testing the File Watcher

```go
// Simple test harness
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
    
    "github.com/yourname/goshield/pkg/agent"
)

func main() {
    eventCh := make(chan interface{}, 1000)
    
    config := agent.DefaultFileWatchConfig()
    config.Paths = []string{"/tmp/test-watch"} // watch just this dir
    
    // Create test directory
    os.MkdirAll("/tmp/test-watch", 0755)
    
    watcher, err := agent.NewFileWatcher("test-agent", "localhost", config, eventCh)
    if err != nil {
        panic(err)
    }
    watcher.Start()
    
    fmt.Println("Watching /tmp/test-watch — try creating/modifying files there!")
    
    // Print events
    go func() {
        for event := range eventCh {
            data, _ := json.MarshalIndent(event, "", "  ")
            fmt.Println(string(data))
        }
    }()
    
    // Run for 30 seconds
    time.Sleep(30 * time.Second)
    watcher.Stop()
}
```

Run this and in another terminal:
```bash
echo "test" > /tmp/test-watch/hello.txt
chmod 777 /tmp/test-watch/hello.txt
rm /tmp/test-watch/hello.txt
```

You'll see file events in real time.

---

## FIM Detection: What to Alert On

```yaml
# rules/file/sensitive_write.yaml
name: Sensitive File Modified
id: FILE-001
severity: high
condition:
  event_type: file
  any_of:
    - field: path
      equals_any: ["/etc/passwd", "/etc/shadow", "/etc/sudoers", "/root/.ssh/authorized_keys"]
    - field: path
      starts_with_any: ["/etc/cron", "/etc/init.d/", "/etc/systemd/system/"]
```

```yaml
# rules/file/webshell_drop.yaml
name: Potential Webshell Dropped
id: FILE-002
severity: critical
condition:
  event_type: file
  all_of:
    - field: action
      equals: "create"
    - field: path
      starts_with_any: ["/var/www/", "/srv/http/", "/usr/share/nginx/"]
    - field: extension
      equals_any: [".php", ".jsp", ".asp", ".aspx", ".phtml"]
```

```yaml
# rules/file/ransomware_extension.yaml
name: File Renamed to Ransomware Extension
id: FILE-003
severity: critical
description: Files being renamed to common ransomware extensions
condition:
  event_type: file
  all_of:
    - field: action
      equals: "rename"
    - field: new_path
      ends_with_any: [".locked", ".encrypted", ".enc", ".crypto", ".crypt", ".WNCRYPT", ".locky"]
```

---

## Performance Considerations

A busy server might generate thousands of file events per second. Real performance tuning:

**1. Batch events before sending**
```go
// Collect events for 5 seconds then send as batch
ticker := time.NewTicker(5 * time.Second)
batch := make([]interface{}, 0, 100)
for {
    select {
    case event := <-eventCh:
        batch = append(batch, event)
        if len(batch) >= 100 {
            sendBatch(batch)
            batch = batch[:0]
        }
    case <-ticker.C:
        if len(batch) > 0 {
            sendBatch(batch)
            batch = batch[:0]
        }
    }
}
```

**2. Only hash new/written files, not all events**

**3. Exclude high-frequency paths** (log files, temp files) in config

**4. Rate limit per-directory** — if one dir generates >1000 events/sec, throttle

---

## Summary

| Capability | Implementation |
|------------|---------------|
| Watch directories | `fsnotify.Watcher` |
| File metadata | `os.Stat()` |
| File content hash | `sha256.New()` + `io.Copy()` |
| Deduplication | In-memory cache with TTL |
| Excluded patterns | `filepath.Match()` |
| New dir watching | Recursive `filepath.Walk()` |

---

## Exercises

1. Build and run the file watcher. Create files in `/tmp/test-watch` and observe events.
2. Add a check: if a created file is executable (permissions include `+x`), flag it at higher severity.
3. Implement a "baseline" feature: hash all existing files on startup, then alert only when files change from their baseline hash.
4. Add Windows support: on Windows, the paths would be `C:\Windows\System32`, `C:\Users`. How would you detect that?
5. Measure the agent's CPU usage with `htop` while file events are happening at high rate. Optimize if needed.
