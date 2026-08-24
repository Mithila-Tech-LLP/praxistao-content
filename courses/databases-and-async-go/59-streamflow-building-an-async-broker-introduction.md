# Chapter 59: StreamFlow — Building an Async Broker from Scratch

We've used Kafka and RabbitMQ. Now we build our own message broker. StreamFlow is a simplified but real message broker written in Go — it handles topics, partitions, consumer groups, and persistence. Building it teaches you everything Kafka's source code took millions of dollars to learn.

## Table of Contents

1. What We're Building
2. Architecture Overview
3. Project Setup
4. Core Data Structures
5. The Storage Layer
6. What's Next

---

## 1. What We're Building

StreamFlow will support:
- **Topics with partitions** — data is distributed across N partitions
- **Producers** — publish messages to a topic
- **Consumer groups** — multiple services each get all messages
- **Offset tracking** — consumers resume from where they left off
- **Persistence** — messages survive restarts
- **TCP wire protocol** — Go, Python, or any language can connect

It won't support (to keep it learnable):
- Replication (one broker)
- Schema registry
- Message compression

By the end, you'll have a working broker that can handle millions of messages per day and teach you why every design decision in Kafka exists.

---

## 2. Architecture Overview

```
                    ┌──────────────────────────────────────────────────────────┐
                    │                    StreamFlow Broker                      │
                    │                                                            │
    Producer ───────►  Wire Protocol (TCP)  ─── Message Router                 │
                    │                               │                           │
    Consumer ───────►                               ▼                           │
                    │                         Topic Manager                     │
                    │                         ┌──────────────────────┐          │
                    │                         │ Topic: "orders"      │          │
                    │                         │  Partition 0 (log)   │          │
                    │                         │  Partition 1 (log)   │          │
                    │                         │  Partition 2 (log)   │          │
                    │                         └──────────────────────┘          │
                    │                                                            │
                    │                         Consumer Group Manager            │
                    │                         ┌──────────────────────┐          │
                    │                         │ Group: "email-svc"   │          │
                    │                         │  Partition 0: offset=42│        │
                    │                         │  Partition 1: offset=38│        │
                    │                         └──────────────────────┘          │
                    └──────────────────────────────────────────────────────────┘
```

**Each layer:**

- **Wire Protocol:** TCP server. Clients connect, send commands, receive responses.
- **Message Router:** Parses commands, routes to topic manager or consumer group manager.
- **Topic Manager:** Manages topic and partition lifecycle.
- **Partition Log:** The append-only log for one partition. Same concept as Kafka's segment files.
- **Consumer Group Manager:** Tracks offset per consumer group per partition.

---

## 3. Project Setup

```bash
mkdir streamflow && cd streamflow
go mod init github.com/yourname/streamflow
mkdir -p log topic group wire
touch main.go
```

Directory structure:
```
streamflow/
├── main.go
├── log/
│   ├── segment.go    # One log segment file
│   └── log.go        # Manages multiple segments
├── topic/
│   ├── partition.go  # One partition = one log
│   └── manager.go    # Topic manager
├── group/
│   └── manager.go    # Consumer group offset tracking
└── wire/
    ├── protocol.go   # Message format
    └── server.go     # TCP server
```

---

## 4. Core Data Structures

```go
// wire/protocol.go
package wire

import "encoding/binary"

// StreamFlow Wire Protocol
//
// Every message:
//   2 bytes: command type
//   4 bytes: payload length
//   N bytes: payload
//
// Commands:
//   PRODUCE  (0x01): publish messages
//   FETCH    (0x02): read messages
//   COMMIT   (0x03): commit consumer offset
//   CREATE   (0x04): create a topic
//   OFFSETS  (0x05): get current offsets for a consumer group
//   ERROR    (0xFF): server error response

type Command uint16

const (
    CmdProduce Command = 0x01
    CmdFetch   Command = 0x02
    CmdCommit  Command = 0x03
    CmdCreate  Command = 0x04
    CmdOffsets Command = 0x05
    CmdError   Command = 0xFF
    CmdOK      Command = 0x00
)

// ProduceRequest: publish messages to a topic
// Payload:
//   2 bytes: topic name length
//   N bytes: topic name
//   1 byte:  partition count (0 = auto-assign)
//   2 bytes: number of messages
//   For each message:
//     2 bytes: key length
//     N bytes: key
//     4 bytes: value length
//     N bytes: value

type Message struct {
    Key   []byte
    Value []byte
}

type ProduceRequest struct {
    Topic     string
    Partition int // -1 = auto-assign
    Messages  []Message
}

// FetchRequest: read messages from a topic
// Payload:
//   2 bytes: topic name length
//   N bytes: topic name
//   1 byte:  partition
//   8 bytes: offset to start from
//   4 bytes: max bytes to return
//   2 bytes: consumer group name length
//   N bytes: consumer group name

type FetchRequest struct {
    Topic     string
    Partition int
    Offset    int64
    MaxBytes  int32
    GroupID   string
}

// FetchResponse: the returned messages
type FetchResponse struct {
    Topic     string
    Partition int
    Offset    int64 // offset of the first message
    Messages  []FetchMessage
}

type FetchMessage struct {
    Offset    int64
    Timestamp int64
    Key       []byte
    Value     []byte
}

// CommitRequest: commit consumer offset
type CommitRequest struct {
    GroupID   string
    Topic     string
    Partition int
    Offset    int64
}

// CreateTopicRequest: create a new topic
type CreateTopicRequest struct {
    Topic      string
    Partitions int
}

// Encode/decode helpers

func EncodeString(s string) []byte {
    b := make([]byte, 2+len(s))
    binary.BigEndian.PutUint16(b[:2], uint16(len(s)))
    copy(b[2:], s)
    return b
}

func DecodeString(data []byte, off int) (string, int) {
    length := int(binary.BigEndian.Uint16(data[off:]))
    return string(data[off+2 : off+2+length]), off + 2 + length
}

func EncodeInt64(n int64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(n))
    return b
}

func DecodeInt64(data []byte, off int) (int64, int) {
    return int64(binary.BigEndian.Uint64(data[off:])), off + 8
}
```

---

## 5. The Storage Layer

The core of any broker is its storage. We use a simple append-only binary log:

```go
// log/segment.go
package log

import (
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "sync"
)

// Segment is one log file (like Kafka's .log files)
// A partition consists of multiple segments. When a segment reaches maxBytes,
// a new one is created. Old segments can be deleted for retention.

// Message format in segment file:
//   8 bytes: offset
//   8 bytes: timestamp (unix nanoseconds)
//   2 bytes: key length
//   N bytes: key
//   4 bytes: value length
//   N bytes: value

const (
    msgHeaderSize = 8 + 8 + 2 + 4 // offset + ts + keyLen + valueLen
)

type Segment struct {
    mu          sync.RWMutex
    file        *os.File
    path        string
    baseOffset  int64   // offset of the first message in this segment
    nextOffset  int64   // offset to assign to the next message
    size        int64   // bytes written
}

func OpenSegment(path string, baseOffset int64) (*Segment, error) {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return nil, fmt.Errorf("open segment %s: %w", path, err)
    }

    seg := &Segment{
        file:       f,
        path:       path,
        baseOffset: baseOffset,
        nextOffset: baseOffset,
    }

    // Scan existing messages to find nextOffset and size
    if err := seg.scan(); err != nil {
        return nil, err
    }

    return seg, nil
}

func (s *Segment) scan() error {
    info, err := s.file.Stat()
    if err != nil {
        return err
    }

    s.size = info.Size()
    if s.size == 0 {
        return nil
    }

    // Read through to find the last offset
    _, _ = s.file.Seek(0, io.SeekStart)
    for {
        var header [msgHeaderSize]byte
        _, err := io.ReadFull(s.file, header[:])
        if err == io.EOF || err == io.ErrUnexpectedEOF {
            break
        }
        if err != nil {
            return err
        }

        offset := int64(binary.BigEndian.Uint64(header[0:8]))
        keyLen := int(binary.BigEndian.Uint16(header[16:18]))
        valLen := int(binary.BigEndian.Uint32(header[18:22]))

        // Skip key and value
        if _, err := s.file.Seek(int64(keyLen+valLen), io.SeekCurrent); err != nil {
            break
        }
        s.nextOffset = offset + 1
    }
    _, _ = s.file.Seek(0, io.SeekEnd)
    return nil
}

// Append writes a message to the segment and returns its offset
func (s *Segment) Append(key, value []byte, timestamp int64) (int64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    offset := s.nextOffset

    // Build message bytes
    buf := make([]byte, msgHeaderSize+len(key)+len(value))
    binary.BigEndian.PutUint64(buf[0:8], uint64(offset))
    binary.BigEndian.PutUint64(buf[8:16], uint64(timestamp))
    binary.BigEndian.PutUint16(buf[16:18], uint16(len(key)))
    binary.BigEndian.PutUint32(buf[18:22], uint32(len(value)))
    copy(buf[msgHeaderSize:], key)
    copy(buf[msgHeaderSize+len(key):], value)

    if _, err := s.file.Write(buf); err != nil {
        return 0, fmt.Errorf("segment write: %w", err)
    }

    s.nextOffset++
    s.size += int64(len(buf))
    return offset, nil
}

// Read reads messages starting from 'fromOffset', up to maxBytes
func (s *Segment) Read(fromOffset int64, maxBytes int32) ([]FetchMessage, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    f, err := os.Open(s.path) // open for reading without affecting write position
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var messages []FetchMessage
    var bytesRead int32

    for {
        var header [msgHeaderSize]byte
        if _, err := io.ReadFull(f, header[:]); err != nil {
            break
        }

        offset := int64(binary.BigEndian.Uint64(header[0:8]))
        ts := int64(binary.BigEndian.Uint64(header[8:16]))
        keyLen := int(binary.BigEndian.Uint16(header[16:18]))
        valLen := int(binary.BigEndian.Uint32(header[18:22]))

        if offset < fromOffset {
            // Skip this message
            if _, err := f.Seek(int64(keyLen+valLen), io.SeekCurrent); err != nil {
                break
            }
            continue
        }

        totalSize := int32(msgHeaderSize + keyLen + valLen)
        if bytesRead+totalSize > maxBytes && len(messages) > 0 {
            break
        }

        key := make([]byte, keyLen)
        val := make([]byte, valLen)
        io.ReadFull(f, key)
        io.ReadFull(f, val)

        messages = append(messages, FetchMessage{
            Offset:    offset,
            Timestamp: ts,
            Key:       key,
            Value:     val,
        })
        bytesRead += totalSize
    }

    return messages, nil
}

func (s *Segment) Sync() error {
    return s.file.Sync()
}

func (s *Segment) Close() error {
    s.file.Sync()
    return s.file.Close()
}

func (s *Segment) BaseOffset() int64 { return s.baseOffset }
func (s *Segment) NextOffset() int64 {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.nextOffset
}
func (s *Segment) Size() int64 {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.size
}
```

---

## 6. What's Next

Over the next 8 chapters:

| Chapter | What We Build |
|---------|--------------|
| 60 | The Partition Log (managing multiple segments) |
| 61 | The Topic Manager (creating topics, routing to partitions) |
| 62 | Consumer Group Manager (offset tracking, group rebalancing) |
| 63 | The Wire Protocol Server (TCP, parsing commands) |
| 64 | The Go Client SDK |
| 65 | Performance and Tuning |
| 66 | Clustering (multi-node basics) |
| 67 | Major Project: Complete StreamFlow |

---

## Summary

- StreamFlow follows the same design as Kafka: append-only log per partition, consumer groups with offsets, TCP wire protocol.
- A segment is a single file storing messages as length-prefixed binary records.
- Each message has: offset (assigned by broker), timestamp, key, value.
- The scan step on startup reads through the segment to find the current `nextOffset`.
- New messages are always appended to the end — no seeks needed for writes.

### Exercises

**Easy:** Write a test for the `Segment`: append 100 messages with increasing keys and values. Then read them back and verify all 100 are present with correct offsets (0-99).

**Medium:** Implement `Segment.ReadAt(offset int64) (*FetchMessage, error)` that reads a single message at a specific offset efficiently (without scanning from the beginning).

**Hard:** Implement a sparse index for fast offset lookup: every 100 messages, record `(offset, bytePosition)` in a small index file. Implement `Segment.SeekToOffset(offset int64)` that uses the index to seek to the right position in O(log n) time.
