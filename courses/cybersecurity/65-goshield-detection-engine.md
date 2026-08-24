# Chapter 65: GoShield Detection Engine — Rule-Based Threat Detection

*The detection engine is where security logic lives. It evaluates every event against rules and decides: threat or benign? This chapter builds a complete YAML-rule-based detection engine in Go.*

---

## Detection Philosophy

Two major approaches in EDR products:

**1. Signature-based (what we're building):**
Human-written rules describing known-bad patterns. Like Sigma rules.
- Pros: precise, explainable, easy to tune
- Cons: only catches known threats; attackers adapt

**2. ML-based behavioral (SentinelOne's approach):**
Train models on normal behavior, flag anomalies.
- Pros: catches novel threats
- Cons: false positives, black box, needs huge training data

We'll build the rule-based engine (approach 1). Adding ML is Chapter 69.

---

## The Rule Structure

Our YAML rules support:

```yaml
name: string               # Human-readable name
id: string                 # Unique identifier (e.g., PROC-001)
severity: info|low|medium|high|critical
description: string
mitre: string              # MITRE ATT&CK reference
references: [string]       # URLs with more context

condition:
  event_type: process|file|network
  
  # Simple field match
  all_of:
    - field: name
      equals: "bash"
  
  any_of:
    - field: command_line
      contains: "-i"
    - field: command_line
      contains: ">& /dev/tcp"
  
  # Threshold-based
  threshold:
    count: 50
    window: 30s
    group_by: pid    # Count per PID

response:
  alert: true
  kill_process: false
  quarantine: false
```

---

## Rule Types and Operations

```go
// pkg/detector/rule.go
package detector

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
    
    "gopkg.in/yaml.v3"
    "github.com/yourname/goshield/pkg/events"
)

// Severity levels
type Severity string

const (
    SeverityInfo     Severity = "info"
    SeverityLow      Severity = "low"
    SeverityMedium   Severity = "medium"
    SeverityHigh     Severity = "high"
    SeverityCritical Severity = "critical"
)

// Rule is a parsed detection rule
type Rule struct {
    Name        string        `yaml:"name"`
    ID          string        `yaml:"id"`
    Severity    Severity      `yaml:"severity"`
    Description string        `yaml:"description"`
    MITRE       string        `yaml:"mitre"`
    References  []string      `yaml:"references"`
    Condition   RuleCondition `yaml:"condition"`
    Response    RuleResponse  `yaml:"response"`
    
    // Compiled state (not from YAML)
    compiledMatchers []FieldMatcher
}

// RuleCondition defines when a rule fires
type RuleCondition struct {
    EventType string         `yaml:"event_type"`
    AllOf     []FieldCondition `yaml:"all_of"`
    AnyOf     []FieldCondition `yaml:"any_of"`
    Threshold *ThresholdCond  `yaml:"threshold"`
}

// FieldCondition is one condition on a single field
type FieldCondition struct {
    Field          string   `yaml:"field"`
    Equals         string   `yaml:"equals"`
    EqualsAny      []string `yaml:"equals_any"`
    Contains       string   `yaml:"contains"`
    ContainsAny    []string `yaml:"contains_any"`
    ContainsAll    []string `yaml:"contains_all"`
    StartsWith     string   `yaml:"starts_with"`
    StartsWithAny  []string `yaml:"starts_with_any"`
    EndsWith       string   `yaml:"ends_with"`
    EndsWithAny    []string `yaml:"ends_with_any"`
    Regex          string   `yaml:"regex"`
    Not            bool     `yaml:"not"`
    CaseInsensitive bool    `yaml:"case_insensitive"`
    
    // Compiled regex
    compiledRegex *regexp.Regexp
}

// ThresholdCond detects volume anomalies
type ThresholdCond struct {
    Count   int    `yaml:"count"`
    Window  string `yaml:"window"`       // e.g. "30s", "1m"
    GroupBy string `yaml:"group_by"`     // e.g. "pid", "hostname"
    
    windowDuration time.Duration
}

// RuleResponse defines what to do when rule fires
type RuleResponse struct {
    Alert        bool `yaml:"alert"`
    KillProcess  bool `yaml:"kill_process"`
    Quarantine   bool `yaml:"quarantine"`
}

// LoadRules loads all YAML rules from a directory
func LoadRules(dir string) ([]*Rule, error) {
    var rules []*Rule
    
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
            return nil
        }
        
        rule, err := loadRuleFile(path)
        if err != nil {
            return fmt.Errorf("loading %s: %w", path, err)
        }
        
        if err := rule.Compile(); err != nil {
            return fmt.Errorf("compiling %s: %w", path, err)
        }
        
        rules = append(rules, rule)
        return nil
    })
    
    return rules, err
}

func loadRuleFile(path string) (*Rule, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var rule Rule
    if err := yaml.Unmarshal(data, &rule); err != nil {
        return nil, err
    }
    
    return &rule, nil
}

// Compile pre-compiles regexes and parses durations
func (r *Rule) Compile() error {
    // Compile regexes in conditions
    for i := range r.Condition.AllOf {
        if err := r.Condition.AllOf[i].compile(); err != nil {
            return fmt.Errorf("all_of[%d]: %w", i, err)
        }
    }
    for i := range r.Condition.AnyOf {
        if err := r.Condition.AnyOf[i].compile(); err != nil {
            return fmt.Errorf("any_of[%d]: %w", i, err)
        }
    }
    
    // Parse threshold window duration
    if r.Condition.Threshold != nil {
        d, err := time.ParseDuration(r.Condition.Threshold.Window)
        if err != nil {
            return fmt.Errorf("invalid threshold window: %w", err)
        }
        r.Condition.Threshold.windowDuration = d
    }
    
    return nil
}

func (fc *FieldCondition) compile() error {
    if fc.Regex != "" {
        re, err := regexp.Compile(fc.Regex)
        if err != nil {
            return err
        }
        fc.compiledRegex = re
    }
    return nil
}

// Matches checks if a field value satisfies a condition
func (fc *FieldCondition) Matches(value string) bool {
    if fc.CaseInsensitive {
        value = strings.ToLower(value)
    }
    
    normalize := func(s string) string {
        if fc.CaseInsensitive {
            return strings.ToLower(s)
        }
        return s
    }
    
    result := false
    
    switch {
    case fc.Equals != "":
        result = value == normalize(fc.Equals)
    case len(fc.EqualsAny) > 0:
        for _, v := range fc.EqualsAny {
            if value == normalize(v) {
                result = true
                break
            }
        }
    case fc.Contains != "":
        result = strings.Contains(value, normalize(fc.Contains))
    case len(fc.ContainsAny) > 0:
        for _, v := range fc.ContainsAny {
            if strings.Contains(value, normalize(v)) {
                result = true
                break
            }
        }
    case len(fc.ContainsAll) > 0:
        result = true
        for _, v := range fc.ContainsAll {
            if !strings.Contains(value, normalize(v)) {
                result = false
                break
            }
        }
    case fc.StartsWith != "":
        result = strings.HasPrefix(value, normalize(fc.StartsWith))
    case len(fc.StartsWithAny) > 0:
        for _, v := range fc.StartsWithAny {
            if strings.HasPrefix(value, normalize(v)) {
                result = true
                break
            }
        }
    case fc.EndsWith != "":
        result = strings.HasSuffix(value, normalize(fc.EndsWith))
    case len(fc.EndsWithAny) > 0:
        for _, v := range fc.EndsWithAny {
            if strings.HasSuffix(value, normalize(v)) {
                result = true
                break
            }
        }
    case fc.compiledRegex != nil:
        result = fc.compiledRegex.MatchString(value)
    }
    
    if fc.Not {
        return !result
    }
    return result
}
```

---

## The Detection Engine

```go
// pkg/detector/engine.go
package detector

import (
    "fmt"
    "sync"
    "time"
    
    "github.com/yourname/goshield/pkg/events"
)

// Detector evaluates events against rules
type Detector struct {
    rules []*Rule
    
    // Threshold tracking: ruleID -> groupKey -> []timestamp
    thresholdMu     sync.Mutex
    thresholdCounts map[string]map[string][]time.Time
}

// NewDetector creates a detector with rules from dir
func NewDetector(rulesDir string) (*Detector, error) {
    rules, err := LoadRules(rulesDir)
    if err != nil {
        return nil, err
    }
    
    return &Detector{
        rules:           rules,
        thresholdCounts: make(map[string]map[string][]time.Time),
    }, nil
}

// RuleCount returns the number of loaded rules
func (d *Detector) RuleCount() int {
    return len(d.rules)
}

// Evaluate tests an event against all rules and returns any alerts
func (d *Detector) Evaluate(event interface{}) []*events.Alert {
    var alerts []*events.Alert
    
    for _, rule := range d.rules {
        if alert := d.evaluateRule(rule, event); alert != nil {
            alerts = append(alerts, alert)
        }
    }
    
    return alerts
}

// evaluateRule checks one rule against one event
func (d *Detector) evaluateRule(rule *Rule, event interface{}) *events.Alert {
    // Extract common fields based on event type
    eventType, fields := extractFields(event)
    if eventType == "" {
        return nil
    }
    
    // Check event type matches
    if rule.Condition.EventType != "" && rule.Condition.EventType != eventType {
        return nil
    }
    
    // Evaluate all_of conditions
    for _, cond := range rule.Condition.AllOf {
        value, ok := fields[cond.Field]
        if !ok || !cond.Matches(value) {
            return nil // All conditions must match
        }
    }
    
    // Evaluate any_of conditions
    if len(rule.Condition.AnyOf) > 0 {
        matched := false
        for _, cond := range rule.Condition.AnyOf {
            value, ok := fields[cond.Field]
            if ok && cond.Matches(value) {
                matched = true
                break
            }
        }
        if !matched {
            return nil
        }
    }
    
    // Check threshold if defined
    if rule.Condition.Threshold != nil {
        if !d.checkThreshold(rule, fields) {
            return nil
        }
    }
    
    // Rule matched — create alert
    base := extractBase(event)
    return &events.Alert{
        Base: events.Base{
            ID:        events.NewID(),
            AgentID:   base.AgentID,
            Hostname:  base.Hostname,
            Timestamp: time.Now(),
            Type:      events.EventTypeAlert,
        },
        RuleID:      rule.ID,
        RuleName:    rule.Name,
        Severity:    events.Severity(rule.Severity),
        Description: rule.Description,
        EventID:     base.ID,
        EventType:   base.Type,
        MITRE:       rule.MITRE,
    }
}

// checkThreshold implements sliding window counting
func (d *Detector) checkThreshold(rule *Rule, fields map[string]string) bool {
    thr := rule.Condition.Threshold
    
    // Determine the group key
    groupKey := fields[thr.GroupBy]
    if groupKey == "" {
        groupKey = "*"
    }
    
    d.thresholdMu.Lock()
    defer d.thresholdMu.Unlock()
    
    // Initialize if needed
    if d.thresholdCounts[rule.ID] == nil {
        d.thresholdCounts[rule.ID] = make(map[string][]time.Time)
    }
    
    now := time.Now()
    windowStart := now.Add(-thr.windowDuration)
    
    // Add current event timestamp
    times := d.thresholdCounts[rule.ID][groupKey]
    times = append(times, now)
    
    // Remove events outside the window
    start := 0
    for start < len(times) && times[start].Before(windowStart) {
        start++
    }
    times = times[start:]
    d.thresholdCounts[rule.ID][groupKey] = times
    
    // Check if count exceeds threshold
    return len(times) >= thr.Count
}

// extractFields pulls field values from an event into a string map
func extractFields(event interface{}) (string, map[string]string) {
    fields := make(map[string]string)
    
    switch e := event.(type) {
    case *events.ProcessEvent:
        fields["action"] = e.Action
        fields["pid"] = fmt.Sprintf("%d", e.PID)
        fields["ppid"] = fmt.Sprintf("%d", e.PPID)
        fields["name"] = e.Name
        fields["command_line"] = e.CommandLine
        fields["username"] = e.Username
        fields["exe_path"] = e.ExePath
        fields["sha256"] = e.SHA256
        if e.IsElevated {
            fields["is_elevated"] = "true"
        }
        return "process", fields
        
    case *events.FileEvent:
        fields["action"] = e.Action
        fields["path"] = e.Path
        fields["new_path"] = e.NewPath
        fields["extension"] = e.Extension
        fields["sha256"] = e.SHA256
        fields["pid"] = fmt.Sprintf("%d", e.PID)
        fields["process"] = e.Process
        if e.IsHidden {
            fields["is_hidden"] = "true"
        }
        return "file", fields
        
    case *events.NetworkEvent:
        fields["action"] = e.Action
        fields["protocol"] = e.Protocol
        fields["src_ip"] = e.SrcIP
        fields["src_port"] = fmt.Sprintf("%d", e.SrcPort)
        fields["dst_ip"] = e.DstIP
        fields["dst_port"] = fmt.Sprintf("%d", e.DstPort)
        fields["pid"] = fmt.Sprintf("%d", e.PID)
        fields["process"] = e.Process
        fields["domain"] = e.Domain
        return "network", fields
    }
    
    return "", nil
}

func extractBase(event interface{}) events.Base {
    switch e := event.(type) {
    case *events.ProcessEvent:
        return e.Base
    case *events.FileEvent:
        return e.Base
    case *events.NetworkEvent:
        return e.Base
    }
    return events.Base{}
}
```

---

## Testing Detection Rules

```go
// Simple test to verify rules fire correctly
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
    
    "github.com/yourname/goshield/pkg/detector"
    "github.com/yourname/goshield/pkg/events"
)

func main() {
    // Write a test rule file
    ruleContent := `
name: PowerShell Encoded Command
id: TEST-001
severity: high
description: Test detection
condition:
  event_type: process
  all_of:
    - field: name
      contains: "powershell"
      case_insensitive: true
    - field: command_line
      contains: "-EncodedCommand"
      case_insensitive: true
response:
  alert: true
`
    os.MkdirAll("/tmp/test-rules", 0755)
    os.WriteFile("/tmp/test-rules/test.yaml", []byte(ruleContent), 0644)
    
    d, err := detector.NewDetector("/tmp/test-rules")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Loaded %d rules\n", d.RuleCount())
    
    // Test event that SHOULD fire
    badEvent := &events.ProcessEvent{
        Base: events.Base{
            ID: "test-event-1", Hostname: "testhost",
            Timestamp: time.Now(), Type: events.EventTypeProcess,
        },
        Action:      "create",
        Name:        "powershell.exe",
        CommandLine: "powershell.exe -EncodedCommand aGVsbG8gd29ybGQ=",
        PID:         1234,
    }
    
    alerts := d.Evaluate(badEvent)
    if len(alerts) > 0 {
        fmt.Println("[✓] Detection worked!")
        data, _ := json.MarshalIndent(alerts[0], "", "  ")
        fmt.Println(string(data))
    } else {
        fmt.Println("[✗] No detection — check rule")
    }
    
    // Test event that should NOT fire
    goodEvent := &events.ProcessEvent{
        Base: events.Base{ID: "test-event-2", Hostname: "testhost", Timestamp: time.Now()},
        Action: "create",
        Name:   "chrome.exe",
        CommandLine: "chrome --profile-directory=Default",
    }
    
    alerts2 := d.Evaluate(goodEvent)
    if len(alerts2) == 0 {
        fmt.Println("[✓] No false positive — correct!")
    } else {
        fmt.Println("[✗] False positive — check rule")
    }
}
```

---

## Built-in Rule Library

```
rules/
├── process/
│   ├── psh_encoded.yaml          # PowerShell encoded command
│   ├── psh_download.yaml         # PowerShell downloading from internet  
│   ├── reverse_shell.yaml        # Bash/nc reverse shells
│   ├── suspicious_parent.yaml    # Office spawning shells
│   ├── scheduled_task.yaml       # New scheduled tasks
│   ├── crontab_modification.yaml # cron changes
│   └── privilege_escalation.yaml # sudo -i, su -
│
├── file/
│   ├── sensitive_write.yaml      # /etc/passwd, /etc/shadow changes
│   ├── webshell_drop.yaml        # .php in web root
│   ├── ransomware_ext.yaml       # .encrypted extension renames
│   ├── suid_binary_new.yaml      # New SUID binary
│   └── ssh_key_add.yaml          # New authorized_keys entry
│
└── network/
    ├── outbound_rare_port.yaml   # Connections to unusual ports
    ├── dns_tunneling.yaml        # Unusually long DNS queries
    ├── beaconing.yaml            # Regular interval connections (C2)
    └── tor_exit_node.yaml        # Connection to known Tor exit nodes
```

---

## Summary

The detection engine:
1. Loads YAML rules from disk at startup
2. Compiles regex patterns and parses durations
3. For each incoming event: checks event type, evaluates `all_of` + `any_of` conditions
4. For threshold rules: maintains sliding window per group key
5. Returns `Alert` objects when rules match

This is architecturally equivalent to how **Sigma rules** work — the same format used by commercial EDR vendors for sharing detection logic.

---

## Exercises

1. Write a rule that detects `nc -l -p XXXX` (netcat listening on a port)
2. Write a threshold rule: alert if any single PID creates more than 10 new processes in 60 seconds
3. Add a `not_any_of` condition type to the rule schema
4. Implement rule tags (e.g., tag: ["ransomware", "windows"]) and add filtering
5. Benchmark the detection engine: how many events per second can it evaluate?
