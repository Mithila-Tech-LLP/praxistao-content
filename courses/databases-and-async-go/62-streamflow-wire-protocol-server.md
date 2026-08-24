# Chapter 62: StreamFlow — Wire Protocol Server

The wire protocol server is the broker's front door. Every client — whether it's a Go service, a Python script, or a Rust program — connects here via TCP and speaks a shared binary language. This chapter builds that server.

## Table of Contents

1. The Protocol Framing
2. Encoding and Decoding Commands
3. The TCP Server
4. The Request Handler
5. Putting It Together in main.go
6. Testing with netcat
7. Exercises

---

## 1. The Protocol Framing

Every message in the StreamFlow protocol has this shape:

```
┌──────────────────────────────────────────────────────────────┐
│  2 bytes  │  4 bytes  │          N bytes                     │
│  command  │  length   │          payload                     │
└──────────────────────────────────────────────────────────────┘
```

This is the same idea as HTTP/2 frames, Redis's RESP protocol, or PostgreSQL's wire protocol: a fixed-size header tells you how many bytes to read next, so you never have to guess.

```go
// wire/framing.go
package wire

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
)

const headerSize = 6 // 2 (command) + 4 (length)

func WriteFrame(conn net.Conn, cmd Command, payload []byte) error {
    header := make([]byte, headerSize)
    binary.BigEndian.PutUint16(header[0:2], uint16(cmd))
    binary.BigEndian.PutUint32(header[2:6], uint32(len(payload)))

    // Write header + payload atomically using writev-style concat
    buf := append(header, payload...)
    _, err := conn.Write(buf)
    return err
}

func ReadFrame(conn net.Conn) (Command, []byte, error) {
    header := make([]byte, headerSize)
    if _, err := io.ReadFull(conn, header); err != nil {
        return 0, nil, fmt.Errorf("read header: %w", err)
    }

    cmd := Command(binary.BigEndian.Uint16(header[0:2]))
    length := binary.BigEndian.Uint32(header[2:6])

    if length == 0 {
        return cmd, nil, nil
    }

    if length > 10*1024*1024 { // 10 MB max message
        return 0, nil, fmt.Errorf("message too large: %d bytes", length)
    }

    payload := make([]byte, length)
    if _, err := io.ReadFull(conn, payload); err != nil {
        return 0, nil, fmt.Errorf("read payload: %w", err)
    }

    return cmd, payload, nil
}
```

---

## 2. Encoding and Decoding Commands

```go
// wire/codec.go
package wire

import (
    "encoding/binary"
    "fmt"
)

// ProduceRequest encoding:
//   2: topic len
//   N: topic
//   4: partition (-1 = auto)
//   2: message count
//   for each: 2 key len, N key, 4 val len, N val

func EncodeProduceRequest(r ProduceRequest) []byte {
    var buf []byte

    topicBytes := []byte(r.Topic)
    tmp := make([]byte, 2)
    binary.BigEndian.PutUint16(tmp, uint16(len(topicBytes)))
    buf = append(buf, tmp...)
    buf = append(buf, topicBytes...)

    p := make([]byte, 4)
    binary.BigEndian.PutUint32(p, uint32(int32(r.Partition)))
    buf = append(buf, p...)

    cnt := make([]byte, 2)
    binary.BigEndian.PutUint16(cnt, uint16(len(r.Messages)))
    buf = append(buf, cnt...)

    for _, msg := range r.Messages {
        kl := make([]byte, 2)
        binary.BigEndian.PutUint16(kl, uint16(len(msg.Key)))
        buf = append(buf, kl...)
        buf = append(buf, msg.Key...)

        vl := make([]byte, 4)
        binary.BigEndian.PutUint32(vl, uint32(len(msg.Value)))
        buf = append(buf, vl...)
        buf = append(buf, msg.Value...)
    }

    return buf
}

func DecodeProduceRequest(data []byte) (ProduceRequest, error) {
    var r ProduceRequest
    off := 0

    if len(data) < 2 {
        return r, fmt.Errorf("short data")
    }
    topicLen := int(binary.BigEndian.Uint16(data[off:]))
    off += 2
    r.Topic = string(data[off : off+topicLen])
    off += topicLen

    r.Partition = int(int32(binary.BigEndian.Uint32(data[off:])))
    off += 4

    msgCount := int(binary.BigEndian.Uint16(data[off:]))
    off += 2

    r.Messages = make([]Message, msgCount)
    for i := 0; i < msgCount; i++ {
        kl := int(binary.BigEndian.Uint16(data[off:]))
        off += 2
        r.Messages[i].Key = data[off : off+kl]
        off += kl

        vl := int(binary.BigEndian.Uint32(data[off:]))
        off += 4
        r.Messages[i].Value = data[off : off+vl]
        off += vl
    }

    return r, nil
}

// FetchRequest encoding:
//   2: topic len, N: topic, 1: partition, 8: offset, 4: maxBytes, 2: groupID len, N: groupID

func EncodeFetchRequest(r FetchRequest) []byte {
    var buf []byte
    buf = append(buf, EncodeString(r.Topic)...)

    tmp := make([]byte, 13) // 1+8+4
    tmp[0] = byte(r.Partition)
    binary.BigEndian.PutUint64(tmp[1:9], uint64(r.Offset))
    binary.BigEndian.PutUint32(tmp[9:13], uint32(r.MaxBytes))
    buf = append(buf, tmp...)

    buf = append(buf, EncodeString(r.GroupID)...)
    return buf
}

func DecodeFetchRequest(data []byte) (FetchRequest, error) {
    var r FetchRequest
    off := 0

    topic, n := DecodeString(data, off)
    r.Topic = topic
    off = n

    r.Partition = int(data[off])
    off++
    r.Offset = int64(binary.BigEndian.Uint64(data[off:]))
    off += 8
    r.MaxBytes = int32(binary.BigEndian.Uint32(data[off:]))
    off += 4

    groupID, _ := DecodeString(data, off)
    r.GroupID = groupID

    return r, nil
}

// CommitRequest encoding: 2 groupID, 2 topic, 1 partition, 8 offset

func EncodeCommitRequest(r CommitRequest) []byte {
    var buf []byte
    buf = append(buf, EncodeString(r.GroupID)...)
    buf = append(buf, EncodeString(r.Topic)...)

    tmp := make([]byte, 9) // 1+8
    tmp[0] = byte(r.Partition)
    binary.BigEndian.PutUint64(tmp[1:], uint64(r.Offset))
    buf = append(buf, tmp...)
    return buf
}

func DecodeCommitRequest(data []byte) (CommitRequest, error) {
    var r CommitRequest
    off := 0

    groupID, n := DecodeString(data, off)
    r.GroupID = groupID
    off = n

    topic, n := DecodeString(data, off)
    r.Topic = topic
    off = n

    r.Partition = int(data[off])
    off++
    r.Offset = int64(binary.BigEndian.Uint64(data[off:]))

    return r, nil
}

// CreateTopicRequest: 2 topic, 1 partitions

func EncodeCreateTopicRequest(r CreateTopicRequest) []byte {
    buf := EncodeString(r.Topic)
    buf = append(buf, byte(r.Partitions))
    return buf
}

func DecodeCreateTopicRequest(data []byte) (CreateTopicRequest, error) {
    topic, n := DecodeString(data, 0)
    return CreateTopicRequest{
        Topic:      topic,
        Partitions: int(data[n]),
    }, nil
}

// FetchResponse encoding:
//   2: topic, 1: partition, 8: base offset, 2: message count
//   for each: 8 offset, 8 ts, 2 keyLen, N key, 4 valLen, N val

func EncodeFetchResponse(r FetchResponse) []byte {
    var buf []byte
    buf = append(buf, EncodeString(r.Topic)...)

    hdr := make([]byte, 11) // 1 partition + 8 offset + 2 count
    hdr[0] = byte(r.Partition)
    binary.BigEndian.PutUint64(hdr[1:9], uint64(r.Offset))
    binary.BigEndian.PutUint16(hdr[9:11], uint16(len(r.Messages)))
    buf = append(buf, hdr...)

    for _, msg := range r.Messages {
        entry := make([]byte, 8+8+2+4+len(msg.Key)+len(msg.Value))
        binary.BigEndian.PutUint64(entry[0:8], uint64(msg.Offset))
        binary.BigEndian.PutUint64(entry[8:16], uint64(msg.Timestamp))
        binary.BigEndian.PutUint16(entry[16:18], uint16(len(msg.Key)))
        binary.BigEndian.PutUint32(entry[18:22], uint32(len(msg.Value)))
        copy(entry[22:], msg.Key)
        copy(entry[22+len(msg.Key):], msg.Value)
        buf = append(buf, entry...)
    }

    return buf
}

func DecodeFetchResponse(data []byte) (FetchResponse, error) {
    var r FetchResponse
    off := 0

    topic, n := DecodeString(data, off)
    r.Topic = topic
    off = n

    r.Partition = int(data[off])
    off++
    r.Offset = int64(binary.BigEndian.Uint64(data[off:]))
    off += 8

    count := int(binary.BigEndian.Uint16(data[off:]))
    off += 2

    r.Messages = make([]FetchMessage, 0, count)
    for i := 0; i < count; i++ {
        msg := FetchMessage{}
        msg.Offset = int64(binary.BigEndian.Uint64(data[off:]))
        off += 8
        msg.Timestamp = int64(binary.BigEndian.Uint64(data[off:]))
        off += 8
        kl := int(binary.BigEndian.Uint16(data[off:]))
        off += 2
        vl := int(binary.BigEndian.Uint32(data[off:]))
        off += 4
        msg.Key = data[off : off+kl]
        off += kl
        msg.Value = data[off : off+vl]
        off += vl
        r.Messages = append(r.Messages, msg)
    }

    return r, nil
}
```

---

## 3. The TCP Server

```go
// wire/server.go
package wire

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    "github.com/yourname/streamflow/group"
    "github.com/yourname/streamflow/topic"
)

type Server struct {
    addr    string
    topics  *topic.Manager
    groups  *group.Manager
    wg      sync.WaitGroup
}

func NewServer(addr string, topics *topic.Manager, groups *group.Manager) *Server {
    return &Server{
        addr:   addr,
        topics: topics,
        groups: groups,
    }
}

func (s *Server) ListenAndServe(ctx context.Context) error {
    ln, err := net.Listen("tcp", s.addr)
    if err != nil {
        return fmt.Errorf("listen: %w", err)
    }

    log.Printf("StreamFlow broker listening on %s", s.addr)

    go func() {
        <-ctx.Done()
        ln.Close()
    }()

    for {
        conn, err := ln.Accept()
        if err != nil {
            if ctx.Err() != nil {
                break
            }
            log.Printf("accept error: %v", err)
            continue
        }

        s.wg.Add(1)
        go func() {
            defer s.wg.Done()
            s.handleConn(conn)
        }()
    }

    s.wg.Wait()
    return nil
}
```

---

## 4. The Request Handler

```go
// wire/handler.go
package wire

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
)

func (s *Server) handleConn(conn net.Conn) {
    defer conn.Close()
    remote := conn.RemoteAddr().String()
    log.Printf("[broker] client connected: %s", remote)

    for {
        cmd, payload, err := ReadFrame(conn)
        if err != nil {
            log.Printf("[broker] client %s disconnected: %v", remote, err)
            return
        }

        switch cmd {
        case CmdCreate:
            s.handleCreate(conn, payload)
        case CmdProduce:
            s.handleProduce(conn, payload)
        case CmdFetch:
            s.handleFetch(conn, payload)
        case CmdCommit:
            s.handleCommit(conn, payload)
        case CmdOffsets:
            s.handleOffsets(conn, payload)
        default:
            sendError(conn, fmt.Sprintf("unknown command: %d", cmd))
        }
    }
}

func (s *Server) handleCreate(conn net.Conn, payload []byte) {
    req, err := DecodeCreateTopicRequest(payload)
    if err != nil {
        sendError(conn, "bad CreateTopic request: "+err.Error())
        return
    }

    if err := s.topics.CreateTopic(req.Topic, req.Partitions); err != nil {
        sendError(conn, err.Error())
        return
    }

    log.Printf("[broker] created topic %q with %d partitions", req.Topic, req.Partitions)
    WriteFrame(conn, CmdOK, []byte("ok"))
}

func (s *Server) handleProduce(conn net.Conn, payload []byte) {
    req, err := DecodeProduceRequest(payload)
    if err != nil {
        sendError(conn, "bad Produce request: "+err.Error())
        return
    }

    // For each message, produce it
    var lastOffset int64
    var lastPartition int
    for _, msg := range req.Messages {
        offset, partition, err := s.topics.Produce(req.Topic, req.Partition, msg.Key, msg.Value)
        if err != nil {
            sendError(conn, err.Error())
            return
        }
        lastOffset = offset
        lastPartition = partition
    }

    // Return: offset of last message, partition used
    resp := make([]byte, 12)
    encodeInt64Into(resp[0:8], lastOffset)
    encodeInt32Into(resp[8:12], int32(lastPartition))
    WriteFrame(conn, CmdOK, resp)
}

func (s *Server) handleFetch(conn net.Conn, payload []byte) {
    req, err := DecodeFetchRequest(payload)
    if err != nil {
        sendError(conn, "bad Fetch request: "+err.Error())
        return
    }

    // If client has a group, use committed offset as starting point when offset == -1
    fetchOffset := req.Offset
    if fetchOffset < 0 && req.GroupID != "" {
        fetchOffset = s.groups.GetOffset(req.GroupID, req.Topic, req.Partition)
    }

    messages, err := s.topics.Fetch(req.Topic, req.Partition, fetchOffset, req.MaxBytes)
    if err != nil {
        sendError(conn, err.Error())
        return
    }

    resp := FetchResponse{
        Topic:     req.Topic,
        Partition: req.Partition,
        Offset:    fetchOffset,
        Messages:  messages,
    }

    WriteFrame(conn, CmdOK, EncodeFetchResponse(resp))
}

func (s *Server) handleCommit(conn net.Conn, payload []byte) {
    req, err := DecodeCommitRequest(payload)
    if err != nil {
        sendError(conn, "bad Commit request: "+err.Error())
        return
    }

    if err := s.groups.CommitOffset(req.GroupID, req.Topic, req.Partition, req.Offset); err != nil {
        sendError(conn, err.Error())
        return
    }

    // Persist asynchronously (best effort)
    go s.groups.Save()

    WriteFrame(conn, CmdOK, nil)
}

func (s *Server) handleOffsets(conn net.Conn, payload []byte) {
    groupID, _ := DecodeString(payload, 0)
    offsets := s.groups.GetAllOffsets(groupID)

    data, _ := json.Marshal(offsets)
    WriteFrame(conn, CmdOK, data)
}

func sendError(conn net.Conn, msg string) {
    WriteFrame(conn, CmdError, []byte(msg))
    log.Printf("[broker] error sent: %s", msg)
}

func encodeInt64Into(dst []byte, n int64) {
    dst[0] = byte(n >> 56)
    dst[1] = byte(n >> 48)
    dst[2] = byte(n >> 40)
    dst[3] = byte(n >> 32)
    dst[4] = byte(n >> 24)
    dst[5] = byte(n >> 16)
    dst[6] = byte(n >> 8)
    dst[7] = byte(n)
}

func encodeInt32Into(dst []byte, n int32) {
    dst[0] = byte(n >> 24)
    dst[1] = byte(n >> 16)
    dst[2] = byte(n >> 8)
    dst[3] = byte(n)
}
```

---

## 5. main.go

```go
// main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/yourname/streamflow/group"
    "github.com/yourname/streamflow/topic"
    "github.com/yourname/streamflow/wire"
)

func main() {
    dataDir := "./data"
    addr := ":9999"

    // Initialize topic manager
    topics, err := topic.NewManager(topic.Config{
        Dir:             dataDir + "/topics",
        MaxSegmentBytes: 10 * 1024 * 1024, // 10 MB per segment
    })
    if err != nil {
        log.Fatalf("topic manager: %v", err)
    }
    defer topics.Close()

    // Initialize consumer group manager
    groups, err := group.NewManager(dataDir + "/groups")
    if err != nil {
        log.Fatalf("group manager: %v", err)
    }

    // Start TCP server
    srv := wire.NewServer(addr, topics, groups)

    ctx, cancel := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer cancel()

    log.Printf("StreamFlow broker starting on %s, data dir: %s", addr, dataDir)
    if err := srv.ListenAndServe(ctx); err != nil {
        log.Printf("server stopped: %v", err)
    }

    log.Println("StreamFlow broker shutdown complete")
}
```

---

## 6. Testing with a Go Smoke Test

```go
// smoke_test.go — run: go test -v -run TestSmoke
package main_test

import (
    "net"
    "testing"
    "time"

    "github.com/yourname/streamflow/wire"
)

func TestSmoke(t *testing.T) {
    conn, err := net.DialTimeout("tcp", "localhost:9999", 2*time.Second)
    if err != nil {
        t.Fatalf("connect: %v", err)
    }
    defer conn.Close()

    // Create topic
    wire.WriteFrame(conn, wire.CmdCreate,
        wire.EncodeCreateTopicRequest(wire.CreateTopicRequest{
            Topic: "smoke-test", Partitions: 2,
        }))

    cmd, resp, err := wire.ReadFrame(conn)
    if err != nil || cmd != wire.CmdOK {
        t.Fatalf("create topic: cmd=%d err=%v resp=%s", cmd, err, resp)
    }

    // Produce a message
    wire.WriteFrame(conn, wire.CmdProduce,
        wire.EncodeProduceRequest(wire.ProduceRequest{
            Topic:     "smoke-test",
            Partition: -1, // auto
            Messages: []wire.Message{
                {Key: []byte("k1"), Value: []byte("hello world")},
            },
        }))

    cmd, _, err = wire.ReadFrame(conn)
    if err != nil || cmd != wire.CmdOK {
        t.Fatalf("produce: %v", err)
    }

    // Fetch it back
    wire.WriteFrame(conn, wire.CmdFetch,
        wire.EncodeFetchRequest(wire.FetchRequest{
            Topic:     "smoke-test",
            Partition: 0,
            Offset:    0,
            MaxBytes:  1024 * 1024,
            GroupID:   "test-group",
        }))

    cmd, data, err := wire.ReadFrame(conn)
    if err != nil || cmd != wire.CmdOK {
        t.Fatalf("fetch: %v", err)
    }

    fetchResp, err := wire.DecodeFetchResponse(data)
    if err != nil {
        t.Fatalf("decode: %v", err)
    }

    t.Logf("Fetched %d messages from partition %d", len(fetchResp.Messages), fetchResp.Partition)
}
```

---

## Summary

- The wire protocol uses a 6-byte header (2B command + 4B length) before every payload.
- `ReadFrame` and `WriteFrame` are the only two transport primitives — everything else is codec.
- The server dispatches on command type, delegating to the topic manager and group manager.
- Consumer groups use offset -1 as "give me from where I last left off."
- Commit is asynchronous (we call `groups.Save()` in a goroutine for performance).

### Exercises

**Easy:** Start the broker (`go run main.go`), then write a Go program that creates a topic, produces 5 messages, fetches them back, and prints each one.

**Medium:** Add a `CmdMetadata` command that returns, for a given topic: number of partitions, highest offset per partition. Encode the response as JSON.

**Hard:** Add connection-level authentication: the first frame from any client must be `CmdAuth` with a preshared secret (configured via env var). If the client sends a wrong secret or skips the auth frame, close the connection. All subsequent commands require authentication to have succeeded.
