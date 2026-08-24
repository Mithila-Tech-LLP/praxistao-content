# Chapter 61: StreamFlow — Topic Manager and Consumer Group Manager

Topics and consumer groups are the two coordination layers above raw partition logs. The topic manager creates and routes to partitions. The consumer group manager tracks where each group has read up to.

## Table of Contents

1. The Topic Manager
2. Partition Assignment (Routing Messages)
3. The Consumer Group Manager
4. Offset Persistence
5. Group Rebalancing
6. Exercises

---

## 1. The Topic Manager

```go
// topic/manager.go
package topic

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    "github.com/yourname/streamflow/log"
)

type Config struct {
    Dir             string
    MaxSegmentBytes int64
}

type TopicInfo struct {
    Name       string `json:"name"`
    Partitions int    `json:"partitions"`
}

type Manager struct {
    mu     sync.RWMutex
    cfg    Config
    topics map[string][]*log.Log // topic → list of partition logs
    info   map[string]TopicInfo
}

func NewManager(cfg Config) (*Manager, error) {
    if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
        return nil, err
    }

    m := &Manager{
        cfg:    cfg,
        topics: make(map[string][]*log.Log),
        info:   make(map[string]TopicInfo),
    }

    // Load existing topics from disk
    if err := m.loadMetadata(); err != nil {
        return nil, err
    }

    return m, nil
}

func (m *Manager) CreateTopic(name string, numPartitions int) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if _, exists := m.topics[name]; exists {
        return fmt.Errorf("topic %q already exists", name)
    }

    partitions := make([]*log.Log, numPartitions)
    for i := 0; i < numPartitions; i++ {
        dir := filepath.Join(m.cfg.Dir, name, fmt.Sprintf("partition-%d", i))
        l, err := log.NewLog(dir, log.Config{
            MaxSegmentBytes: m.cfg.MaxSegmentBytes,
        })
        if err != nil {
            return fmt.Errorf("create partition %d: %w", i, err)
        }
        partitions[i] = l
    }

    m.topics[name] = partitions
    m.info[name] = TopicInfo{Name: name, Partitions: numPartitions}

    return m.saveMetadata()
}

// Produce appends a message to the given partition (or picks one automatically)
func (m *Manager) Produce(topicName string, partition int, key, value []byte) (int64, int, error) {
    m.mu.RLock()
    partitions, ok := m.topics[topicName]
    m.mu.RUnlock()

    if !ok {
        return 0, 0, fmt.Errorf("topic %q not found", topicName)
    }

    // Auto-assign partition
    if partition < 0 {
        if len(key) > 0 {
            // Hash-based: same key → same partition
            partition = int(fnv1a(key)) % len(partitions)
        } else {
            // Round-robin (simplified: use a counter)
            partition = int(rrCounter.Add(1)) % len(partitions)
        }
    }

    if partition >= len(partitions) {
        return 0, 0, fmt.Errorf("partition %d out of range", partition)
    }

    offset, err := partitions[partition].Append(key, value)
    return offset, partition, err
}

// Fetch reads messages from a topic partition
func (m *Manager) Fetch(topicName string, partition int, offset int64, maxBytes int32) ([]log.FetchMessage, error) {
    m.mu.RLock()
    partitions, ok := m.topics[topicName]
    m.mu.RUnlock()

    if !ok {
        return nil, fmt.Errorf("topic %q not found", topicName)
    }
    if partition >= len(partitions) {
        return nil, fmt.Errorf("partition %d out of range", partition)
    }

    return partitions[partition].Read(offset, maxBytes)
}

// HighestOffset returns the next-to-be-assigned offset for a partition
func (m *Manager) HighestOffset(topicName string, partition int) (int64, error) {
    m.mu.RLock()
    partitions, ok := m.topics[topicName]
    m.mu.RUnlock()

    if !ok {
        return 0, fmt.Errorf("topic %q not found", topicName)
    }
    if partition >= len(partitions) {
        return 0, fmt.Errorf("partition %d out of range", partition)
    }

    return partitions[partition].HighestOffset(), nil
}

func (m *Manager) TopicInfo(name string) (TopicInfo, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    info, ok := m.info[name]
    return info, ok
}

func (m *Manager) ListTopics() []TopicInfo {
    m.mu.RLock()
    defer m.mu.RUnlock()
    var list []TopicInfo
    for _, info := range m.info {
        list = append(list, info)
    }
    return list
}

// Persist topic metadata to disk
func (m *Manager) saveMetadata() error {
    path := filepath.Join(m.cfg.Dir, "metadata.json")
    data, err := json.MarshalIndent(m.info, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func (m *Manager) loadMetadata() error {
    path := filepath.Join(m.cfg.Dir, "metadata.json")
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return err
    }

    var info map[string]TopicInfo
    if err := json.Unmarshal(data, &info); err != nil {
        return err
    }

    for name, ti := range info {
        partitions := make([]*log.Log, ti.Partitions)
        for i := 0; i < ti.Partitions; i++ {
            dir := filepath.Join(m.cfg.Dir, name, fmt.Sprintf("partition-%d", i))
            l, err := log.NewLog(dir, log.Config{
                MaxSegmentBytes: m.cfg.MaxSegmentBytes,
            })
            if err != nil {
                return fmt.Errorf("reload partition %s/%d: %w", name, i, err)
            }
            partitions[i] = l
        }
        m.topics[name] = partitions
        m.info[name] = ti
    }
    return nil
}

func (m *Manager) Close() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    for _, partitions := range m.topics {
        for _, p := range partitions {
            p.Close()
        }
    }
    return nil
}

// FNV-1a hash for key-based partitioning
func fnv1a(data []byte) uint32 {
    h := uint32(2166136261)
    for _, b := range data {
        h ^= uint32(b)
        h *= 16777619
    }
    return h
}

// atomic round-robin counter
import "sync/atomic"
var rrCounter atomic.Uint64
```

---

## 2. The Consumer Group Manager

The consumer group manager tracks: for each `(groupID, topic, partition)`, what is the committed offset?

```go
// group/manager.go
package group

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"
)

type CommitKey struct {
    GroupID   string
    Topic     string
    Partition int
}

type Manager struct {
    mu      sync.RWMutex
    dir     string
    offsets map[CommitKey]int64
    dirty   bool
}

func NewManager(dir string) (*Manager, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    m := &Manager{
        dir:     dir,
        offsets: make(map[CommitKey]int64),
    }
    return m, m.load()
}

// GetOffset returns the committed offset for a consumer group+topic+partition.
// Returns 0 if the group has never committed for this partition (start from beginning).
func (m *Manager) GetOffset(groupID, topic string, partition int) int64 {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.offsets[CommitKey{groupID, topic, partition}]
}

// CommitOffset records that a consumer group has processed up to (and including) offset.
func (m *Manager) CommitOffset(groupID, topic string, partition int, offset int64) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    key := CommitKey{groupID, topic, partition}

    // Only move forward (never backward)
    if current := m.offsets[key]; offset <= current {
        return nil
    }

    m.offsets[key] = offset
    m.dirty = true
    return nil
}

// GetAllOffsets returns all offsets for a consumer group
func (m *Manager) GetAllOffsets(groupID string) map[string]map[int]int64 {
    m.mu.RLock()
    defer m.mu.RUnlock()

    result := make(map[string]map[int]int64)
    for key, offset := range m.offsets {
        if key.GroupID != groupID {
            continue
        }
        if result[key.Topic] == nil {
            result[key.Topic] = make(map[int]int64)
        }
        result[key.Topic][key.Partition] = offset
    }
    return result
}

// Save persists the offset table to disk
func (m *Manager) Save() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if !m.dirty {
        return nil
    }

    // Serialize as a flat list for simplicity
    type entry struct {
        GroupID   string `json:"group_id"`
        Topic     string `json:"topic"`
        Partition int    `json:"partition"`
        Offset    int64  `json:"offset"`
    }

    var entries []entry
    for key, offset := range m.offsets {
        entries = append(entries, entry{key.GroupID, key.Topic, key.Partition, offset})
    }

    data, err := json.MarshalIndent(entries, "", "  ")
    if err != nil {
        return err
    }

    path := filepath.Join(m.dir, "offsets.json")
    if err := os.WriteFile(path+".tmp", data, 0644); err != nil {
        return err
    }
    // Atomic rename
    if err := os.Rename(path+".tmp", path); err != nil {
        return err
    }

    m.dirty = false
    return nil
}

func (m *Manager) load() error {
    path := filepath.Join(m.dir, "offsets.json")
    data, err := os.ReadFile(path)
    if os.IsNotExist(err) {
        return nil
    }
    if err != nil {
        return err
    }

    type entry struct {
        GroupID   string `json:"group_id"`
        Topic     string `json:"topic"`
        Partition int    `json:"partition"`
        Offset    int64  `json:"offset"`
    }

    var entries []entry
    if err := json.Unmarshal(data, &entries); err != nil {
        return err
    }

    for _, e := range entries {
        m.offsets[CommitKey{e.GroupID, e.Topic, e.Partition}] = e.Offset
    }
    return nil
}
```

---

## 3. Group Rebalancing

In real Kafka, when a consumer joins or leaves a group, partitions are redistributed (rebalanced). StreamFlow's simplified version:

```go
// Assign partitions to consumers in a group
// Simple strategy: distribute partitions round-robin across consumers

type Assignment struct {
    ConsumerID string
    Partitions []int
}

func AssignPartitions(consumers []string, numPartitions int) []Assignment {
    if len(consumers) == 0 {
        return nil
    }

    assignments := make([]Assignment, len(consumers))
    for i, id := range consumers {
        assignments[i] = Assignment{ConsumerID: id}
    }

    for p := 0; p < numPartitions; p++ {
        consumerIdx := p % len(consumers)
        assignments[consumerIdx].Partitions = append(
            assignments[consumerIdx].Partitions, p)
    }

    return assignments
}

// ConsumerRegistry tracks active consumers per group
type ConsumerRegistry struct {
    mu        sync.Mutex
    consumers map[string][]string // groupID → list of consumer IDs
    sessions  map[string]time.Time // consumerID → last heartbeat
}

func (cr *ConsumerRegistry) Join(groupID, consumerID string) {
    cr.mu.Lock()
    defer cr.mu.Unlock()
    for _, id := range cr.consumers[groupID] {
        if id == consumerID {
            return // already registered
        }
    }
    cr.consumers[groupID] = append(cr.consumers[groupID], consumerID)
    cr.sessions[consumerID] = time.Now()
}

func (cr *ConsumerRegistry) Heartbeat(consumerID string) {
    cr.mu.Lock()
    cr.sessions[consumerID] = time.Now()
    cr.mu.Unlock()
}

func (cr *ConsumerRegistry) evictStale() {
    cr.mu.Lock()
    defer cr.mu.Unlock()
    timeout := 30 * time.Second
    for id, lastSeen := range cr.sessions {
        if time.Since(lastSeen) > timeout {
            delete(cr.sessions, id)
            for groupID, consumers := range cr.consumers {
                for i, cid := range consumers {
                    if cid == id {
                        cr.consumers[groupID] = append(consumers[:i], consumers[i+1:]...)
                        break
                    }
                }
            }
        }
    }
}
```

---

## Summary

- The topic manager routes `Produce` calls to the correct partition log. Key-based routing: `fnv1a(key) % numPartitions`.
- Topic metadata (name, partition count) is persisted to `metadata.json`.
- The consumer group manager tracks `(groupID, topic, partition) → committedOffset`.
- Offsets are only moved forward — committing a lower offset is a no-op.
- Simple rebalancing: round-robin assignment of partitions to consumers.

### Exercises

**Easy:** Create a topic "logs" with 3 partitions. Produce 30 messages with keys "key-0" through "key-29". Print which partition each key was routed to. Verify keys with the same hash always go to the same partition.

**Medium:** Write a test for the consumer group manager: 3 consumers in a group each commit offsets for different partitions. After restarting the manager (reload from disk), verify all offsets are preserved.

**Hard:** Implement consumer lag calculation: `GetLag(groupID, topic, partition)` returns `highestOffset - committedOffset`. Expose this as a `GET /lag/{groupID}/{topic}` HTTP endpoint that returns JSON.
