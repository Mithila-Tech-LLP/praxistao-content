# Chapter 65: StreamFlow — Major Project: Complete Working Broker

Everything we've built comes together here. This chapter assembles all the layers — segment files, partition logs, topic manager, consumer group manager, wire protocol server, and Go SDK — into a single runnable broker with a full test suite.

## Table of Contents

1. Final Directory Structure
2. Assembling the Broker
3. Integration Test Suite
4. Load Test
5. Running StreamFlow
6. Mini Project: Build a Chat App on StreamFlow
7. What We Built vs Kafka
8. Exercises

---

## 1. Final Directory Structure

```
streamflow/
├── main.go                   ← entry point, wires everything
├── go.mod
├── log/
│   ├── segment.go            ← one log file
│   └── log.go                ← multi-segment partition log
├── topic/
│   └── manager.go            ← topic lifecycle, routing
├── group/
│   └── manager.go            ← consumer group offset tracking
├── wire/
│   ├── protocol.go           ← data types, constants
│   ├── framing.go            ← ReadFrame / WriteFrame
│   ├── codec.go              ← encode/decode each command
│   ├── server.go             ← TCP accept loop
│   └── handler.go            ← request dispatch
├── client/
│   ├── conn.go               ← raw connection
│   ├── admin.go              ← AdminClient
│   ├── producer.go           ← batching producer
│   └── consumer.go           ← polling consumer
├── metrics/
│   └── metrics.go            ← Prometheus counters
└── tests/
    ├── integration_test.go
    └── load_test.go
```

---

## 2. The Complete main.go

```go
// main.go
package main

import (
    "context"
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/yourname/streamflow/group"
    "github.com/yourname/streamflow/metrics"
    "github.com/yourname/streamflow/topic"
    "github.com/yourname/streamflow/wire"
)

func main() {
    dataDir := flag.String("data", "./data", "data directory")
    addr := flag.String("addr", ":9999", "broker address")
    metricsAddr := flag.String("metrics", ":9998", "metrics address")
    maxSegMB := flag.Int64("max-segment-mb", 10, "max segment size in MB")
    flag.Parse()

    log.Printf("StreamFlow starting | data=%s addr=%s", *dataDir, *addr)

    // Start Prometheus metrics endpoint
    metrics.StartMetricsServer(*metricsAddr)

    // Initialize topic manager (loads existing topics from disk on startup)
    topics, err := topic.NewManager(topic.Config{
        Dir:             *dataDir + "/topics",
        MaxSegmentBytes: *maxSegMB * 1024 * 1024,
    })
    if err != nil {
        log.Fatalf("topic manager: %v", err)
    }
    defer topics.Close()

    // Initialize consumer group manager (loads committed offsets from disk)
    groups, err := group.NewManager(*dataDir + "/groups")
    if err != nil {
        log.Fatalf("group manager: %v", err)
    }

    // Periodically flush group offsets to disk
    // (they also flush on every Commit, but this ensures nothing is lost on crash)
    go func() {
        for {
            time.Sleep(5 * time.Second)
            if err := groups.Save(); err != nil {
                log.Printf("group save: %v", err)
            }
        }
    }()

    // Start TCP broker
    srv := wire.NewServer(*addr, topics, groups)

    ctx, cancel := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer cancel()

    if err := srv.ListenAndServe(ctx); err != nil {
        log.Printf("broker stopped: %v", err)
    }

    log.Println("StreamFlow shutdown complete")
}
```

---

## 3. Integration Test Suite

```go
// tests/integration_test.go
package tests

import (
    "context"
    "encoding/json"
    "fmt"
    "testing"
    "time"

    "github.com/yourname/streamflow/client"
)

const brokerAddr = "localhost:9999"

func TestProduceAndConsume(t *testing.T) {
    topic := fmt.Sprintf("test-%d", time.Now().UnixNano())

    admin, err := client.NewAdminClient(brokerAddr)
    if err != nil {
        t.Skip("broker not running:", err)
    }
    defer admin.Close()

    if err := admin.CreateTopic(topic, 1); err != nil {
        t.Fatalf("create topic: %v", err)
    }

    // Produce 100 messages
    p, _ := client.NewProducer(client.ProducerConfig{Addr: brokerAddr, BatchSize: 10})
    defer p.Close()

    for i := 0; i < 100; i++ {
        _, _, err := p.SendSync(client.ProducerRecord{
            Topic: topic,
            Key:   []byte(fmt.Sprintf("key-%d", i)),
            Value: []byte(fmt.Sprintf("value-%d", i)),
        })
        if err != nil {
            t.Fatalf("produce %d: %v", i, err)
        }
    }

    // Consume all 100
    c, _ := client.NewConsumer(client.ConsumerConfig{
        Addr:      brokerAddr,
        GroupID:   "test-group",
        Topic:     topic,
        Partition: 0,
        MaxBytes:  1024 * 1024,
    })
    defer c.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var received []client.Record
    for len(received) < 100 {
        records, err := c.Poll(ctx)
        if err != nil {
            t.Fatalf("poll: %v", err)
        }
        received = append(received, records...)
    }

    if len(received) != 100 {
        t.Fatalf("want 100 records, got %d", len(received))
    }

    // Verify ordering and content
    for i, r := range received {
        wantKey := fmt.Sprintf("key-%d", i)
        if string(r.Key) != wantKey {
            t.Errorf("record %d: want key=%q got %q", i, wantKey, r.Key)
        }
        if r.Offset != int64(i) {
            t.Errorf("record %d: want offset=%d got %d", i, i, r.Offset)
        }
    }
}

func TestConsumerResumesAfterRestart(t *testing.T) {
    topic := fmt.Sprintf("resume-%d", time.Now().UnixNano())
    group := "resume-group"

    admin, _ := client.NewAdminClient(brokerAddr)
    admin.CreateTopic(topic, 1)
    admin.Close()

    // Produce 50 messages
    p, _ := client.NewProducer(client.ProducerConfig{Addr: brokerAddr})
    for i := 0; i < 50; i++ {
        p.SendSync(client.ProducerRecord{Topic: topic, Value: []byte(fmt.Sprintf("%d", i))})
    }
    p.Close()

    // Consumer 1: reads and commits first 25
    c1, _ := client.NewConsumer(client.ConsumerConfig{
        Addr: brokerAddr, GroupID: group, Topic: topic, Partition: 0,
    })
    ctx := context.Background()
    var records []client.Record
    for len(records) < 25 {
        batch, _ := c1.Poll(ctx)
        records = append(records, batch...)
    }
    c1.CommitOffset(25)
    c1.Close()

    // Consumer 2 (new connection, same group): should start from 25
    c2, _ := client.NewConsumer(client.ConsumerConfig{
        Addr: brokerAddr, GroupID: group, Topic: topic, Partition: 0,
    })
    defer c2.Close()

    batch, _ := c2.Poll(ctx)
    if len(batch) == 0 {
        t.Fatal("expected messages after offset 25")
    }

    if batch[0].Offset != 25 {
        t.Fatalf("want first offset=25, got %d", batch[0].Offset)
    }
}

func TestKeyBasedPartitioning(t *testing.T) {
    topic := fmt.Sprintf("keys-%d", time.Now().UnixNano())

    admin, _ := client.NewAdminClient(brokerAddr)
    admin.CreateTopic(topic, 4)
    admin.Close()

    p, _ := client.NewProducer(client.ProducerConfig{Addr: brokerAddr})
    defer p.Close()

    // Same key must always land on the same partition
    key := []byte("user-42")
    var seenPartitions []int

    for i := 0; i < 20; i++ {
        _, partition, err := p.SendSync(client.ProducerRecord{
            Topic: topic,
            Key:   key,
            Value: []byte(fmt.Sprintf("msg-%d", i)),
        })
        if err != nil {
            t.Fatal(err)
        }
        seenPartitions = append(seenPartitions, partition)
    }

    // All 20 messages should be on the same partition
    first := seenPartitions[0]
    for i, p := range seenPartitions {
        if p != first {
            t.Fatalf("message %d went to partition %d, want %d", i, p, first)
        }
    }
    t.Logf("key %q consistently routed to partition %d", key, first)
}

func TestLargeMessages(t *testing.T) {
    topic := fmt.Sprintf("large-%d", time.Now().UnixNano())
    admin, _ := client.NewAdminClient(brokerAddr)
    admin.CreateTopic(topic, 1)
    admin.Close()

    p, _ := client.NewProducer(client.ProducerConfig{Addr: brokerAddr})
    defer p.Close()

    // 500 KB message
    bigValue := make([]byte, 500*1024)
    for i := range bigValue {
        bigValue[i] = byte(i % 256)
    }

    offset, _, err := p.SendSync(client.ProducerRecord{Topic: topic, Value: bigValue})
    if err != nil {
        t.Fatalf("produce large message: %v", err)
    }

    c, _ := client.NewConsumer(client.ConsumerConfig{
        Addr: brokerAddr, Topic: topic, Partition: 0,
        MaxBytes: 1024 * 1024,
    })
    defer c.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    records, err := c.Poll(ctx)
    if err != nil {
        t.Fatal(err)
    }

    if len(records) == 0 || records[0].Offset != offset {
        t.Fatal("failed to retrieve large message")
    }
    if len(records[0].Value) != len(bigValue) {
        t.Fatalf("value length mismatch: want %d got %d", len(bigValue), len(records[0].Value))
    }
    t.Logf("Large message round-trip OK (%d bytes)", len(bigValue))
}
```

---

## 4. Mini Project: Chat App on StreamFlow

Build a simple terminal chat app where users publish messages to a "chat" topic and all connected clients receive them in real time.

```go
// chat/main.go
package main

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/yourname/streamflow/client"
)

type ChatMessage struct {
    User    string    `json:"user"`
    Text    string    `json:"text"`
    SentAt  time.Time `json:"sent_at"`
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: chat <username>")
        os.Exit(1)
    }
    username := os.Args[1]

    // Admin: create topic if it doesn't exist
    admin, _ := client.NewAdminClient("localhost:9999")
    admin.CreateTopic("chat", 1) // single partition keeps ordering
    admin.Close()

    // Producer for sending
    p, _ := client.NewProducer(client.ProducerConfig{
        Addr:         "localhost:9999",
        BatchTimeout: 5 * time.Millisecond,
    })
    defer p.Close()

    // Consumer for receiving (different group per user so everyone sees all messages)
    c, _ := client.NewConsumer(client.ConsumerConfig{
        Addr:       "localhost:9999",
        GroupID:    "chat-" + username,
        Topic:      "chat",
        Partition:  0,
        AutoCommit: true,
        MaxBytes:   512 * 1024,
    })
    defer c.Close()

    // Goroutine: poll for incoming messages
    go func() {
        ctx := context.Background()
        for {
            records, err := c.Poll(ctx)
            if err != nil {
                return
            }
            for _, r := range records {
                var msg ChatMessage
                if err := json.Unmarshal(r.Value, &msg); err != nil {
                    continue
                }
                if msg.User != username { // don't echo own messages
                    fmt.Printf("\r[%s] %s: %s\n> ", msg.SentAt.Format("15:04:05"), msg.User, msg.Text)
                }
            }
        }
    }()

    // Main loop: read from stdin and publish
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Printf("Connected as %s. Type and press Enter.\n> ", username)
    for scanner.Scan() {
        text := scanner.Text()
        if text == "" {
            continue
        }

        msg := ChatMessage{User: username, Text: text, SentAt: time.Now()}
        data, _ := json.Marshal(msg)
        p.SendSync(client.ProducerRecord{Topic: "chat", Value: data})
        fmt.Print("> ")
    }
}
```

Run multiple terminals:
```bash
go run chat/main.go alice
go run chat/main.go bob
```

---

## 5. What We Built vs Kafka

| Feature | StreamFlow (ours) | Apache Kafka |
|---------|------------------|--------------|
| Append-only log | ✅ | ✅ |
| Multiple partitions | ✅ | ✅ |
| Consumer groups + offsets | ✅ | ✅ |
| Key-based partitioning | ✅ | ✅ |
| Group commit (batched fsync) | ✅ | ✅ |
| TCP wire protocol | ✅ | ✅ |
| Per-segment files with rotation | ✅ | ✅ |
| Replication | ❌ | ✅ (ISR, acks) |
| Schema registry | ❌ | ✅ (Confluent) |
| Exactly-once semantics | ❌ | ✅ (idempotent producer) |
| Compression | ❌ | ✅ (gzip, snappy, lz4, zstd) |
| Rack-aware assignment | ❌ | ✅ |
| Admin UI | ❌ | ✅ (multiple) |

StreamFlow is ~2,500 lines of Go. Kafka's Java source is ~600,000 lines. We have the core 80% of the value in 0.4% of the code.

---

## Summary

You've built a complete, working message broker from scratch:

- Segment files with group commit for crash-safe writes
- Multi-segment partition logs with retention
- Topic manager with FNV-1a key-based routing
- Consumer group offset tracking with atomic rename persistence
- Binary TCP wire protocol
- Batching Go producer and polling consumer SDK
- Prometheus metrics
- A complete integration test suite

This is the same architecture Kafka, Pulsar, and Redpanda use. Now when you use Kafka in production, you understand every design decision from the ground up.

### Exercises

**Easy:** Add a `streamflow-cli` binary with subcommands: `topics list`, `topics create <name> --partitions N`, `produce <topic> <message>`, `consume <topic> --group <group>`.

**Medium:** Add message TTL at the topic level: topics can be created with `--retention 24h`. Implement a background goroutine that calls `ApplyRetention()` on all partition logs every 5 minutes.

**Hard — Major Project:** Add a replication protocol. Run two StreamFlow instances. When a message is produced to the leader, it replicates to the follower before acknowledging the client. If the leader crashes, the follower takes over. This is the core of Kafka's ISR (in-sync replicas) model.
