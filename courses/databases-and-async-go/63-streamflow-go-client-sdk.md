# Chapter 63: StreamFlow — Go Client SDK

The broker speaks TCP and binary. But nobody wants to hand-craft binary frames in their application code. We build a clean Go SDK that hides all that complexity behind three simple types: `Producer`, `Consumer`, and `AdminClient`.

## Table of Contents

1. The Client Connection
2. AdminClient — Managing Topics
3. Producer — Sending Messages
4. Consumer — Receiving Messages with Auto-Commit
5. A Complete Example: Order Pipeline
6. Exercises

---

## 1. The Client Connection

```go
// client/conn.go
package client

import (
    "fmt"
    "net"
    "sync"
    "time"

    "github.com/yourname/streamflow/wire"
)

// Conn is a single TCP connection to the broker.
// All operations are serialized through a mutex.
// In production you'd use a connection pool — this is educational.

type Conn struct {
    mu   sync.Mutex
    conn net.Conn
    addr string
}

func Dial(addr string) (*Conn, error) {
    nc, err := net.DialTimeout("tcp", addr, 5*time.Second)
    if err != nil {
        return nil, fmt.Errorf("dial %s: %w", addr, err)
    }
    return &Conn{conn: nc, addr: addr}, nil
}

func (c *Conn) Close() error {
    return c.conn.Close()
}

// sendRecv sends a command and reads the response frame
func (c *Conn) sendRecv(cmd wire.Command, payload []byte) (wire.Command, []byte, error) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if err := wire.WriteFrame(c.conn, cmd, payload); err != nil {
        return 0, nil, fmt.Errorf("write: %w", err)
    }

    respCmd, respPayload, err := wire.ReadFrame(c.conn)
    if err != nil {
        return 0, nil, fmt.Errorf("read: %w", err)
    }

    if respCmd == wire.CmdError {
        return 0, nil, fmt.Errorf("broker error: %s", string(respPayload))
    }

    return respCmd, respPayload, nil
}
```

---

## 2. AdminClient

```go
// client/admin.go
package client

import (
    "github.com/yourname/streamflow/wire"
)

type AdminClient struct {
    conn *Conn
}

func NewAdminClient(addr string) (*AdminClient, error) {
    c, err := Dial(addr)
    if err != nil {
        return nil, err
    }
    return &AdminClient{conn: c}, nil
}

func (a *AdminClient) Close() error { return a.conn.Close() }

func (a *AdminClient) CreateTopic(name string, partitions int) error {
    payload := wire.EncodeCreateTopicRequest(wire.CreateTopicRequest{
        Topic:      name,
        Partitions: partitions,
    })

    _, _, err := a.conn.sendRecv(wire.CmdCreate, payload)
    return err
}

func (a *AdminClient) GetOffsets(groupID string) (map[string]map[int]int64, error) {
    groupBytes := wire.EncodeString(groupID)
    _, respPayload, err := a.conn.sendRecv(wire.CmdOffsets, groupBytes)
    if err != nil {
        return nil, err
    }

    var offsets map[string]map[int]int64
    if err := json.Unmarshal(respPayload, &offsets); err != nil {
        return nil, err
    }
    return offsets, nil
}
```

---

## 3. Producer

```go
// client/producer.go
package client

import (
    "context"
    "encoding/binary"
    "fmt"
    "sync"
    "time"

    "github.com/yourname/streamflow/wire"
)

type ProducerConfig struct {
    Addr string
    // BatchSize: collect up to this many messages before flushing
    BatchSize int
    // BatchTimeout: flush at most every this duration
    BatchTimeout time.Duration
}

type ProducerRecord struct {
    Topic     string
    Partition int // -1 = auto
    Key       []byte
    Value     []byte
}

type ProduceResult struct {
    Offset    int64
    Partition int
    Err       error
}

type pendingRecord struct {
    record ProducerRecord
    result chan ProduceResult
}

type Producer struct {
    conn    *Conn
    cfg     ProducerConfig
    pending chan pendingRecord
    wg      sync.WaitGroup
    cancel  context.CancelFunc
}

func NewProducer(cfg ProducerConfig) (*Producer, error) {
    if cfg.BatchSize == 0 {
        cfg.BatchSize = 100
    }
    if cfg.BatchTimeout == 0 {
        cfg.BatchTimeout = 5 * time.Millisecond
    }

    conn, err := Dial(cfg.Addr)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())
    p := &Producer{
        conn:    conn,
        cfg:     cfg,
        pending: make(chan pendingRecord, cfg.BatchSize*2),
        cancel:  cancel,
    }

    p.wg.Add(1)
    go p.runFlusher(ctx)

    return p, nil
}

// Send produces a record asynchronously. The returned channel gets exactly one result.
func (p *Producer) Send(rec ProducerRecord) <-chan ProduceResult {
    ch := make(chan ProduceResult, 1)
    p.pending <- pendingRecord{record: rec, result: ch}
    return ch
}

// SendSync produces a record and waits for acknowledgement.
func (p *Producer) SendSync(rec ProducerRecord) (int64, int, error) {
    ch := p.Send(rec)
    r := <-ch
    return r.Offset, r.Partition, r.Err
}

func (p *Producer) Close() error {
    p.cancel()
    p.wg.Wait()
    return p.conn.Close()
}

func (p *Producer) runFlusher(ctx context.Context) {
    defer p.wg.Done()
    ticker := time.NewTicker(p.cfg.BatchTimeout)
    defer ticker.Stop()

    var batch []pendingRecord

    flush := func() {
        if len(batch) == 0 {
            return
        }
        p.flushBatch(batch)
        batch = batch[:0]
    }

    for {
        select {
        case rec := <-p.pending:
            batch = append(batch, rec)
            if len(batch) >= p.cfg.BatchSize {
                flush()
            }
        case <-ticker.C:
            flush()
        case <-ctx.Done():
            // Drain remaining
            for len(p.pending) > 0 {
                rec := <-p.pending
                batch = append(batch, rec)
            }
            flush()
            return
        }
    }
}

func (p *Producer) flushBatch(batch []pendingRecord) {
    // Group by topic+partition for efficient batching
    type key struct {
        topic     string
        partition int
    }

    grouped := make(map[key][]pendingRecord)
    for _, pr := range batch {
        k := key{pr.record.Topic, pr.record.Partition}
        grouped[k] = append(grouped[k], pr)
    }

    for k, recs := range grouped {
        msgs := make([]wire.Message, len(recs))
        for i, r := range recs {
            msgs[i] = wire.Message{Key: r.record.Key, Value: r.record.Value}
        }

        req := wire.ProduceRequest{
            Topic:     k.topic,
            Partition: k.partition,
            Messages:  msgs,
        }

        _, respPayload, err := p.conn.sendRecv(wire.CmdProduce, wire.EncodeProduceRequest(req))

        var offset int64
        var partition int
        if err == nil && len(respPayload) >= 12 {
            offset = int64(binary.BigEndian.Uint64(respPayload[0:8]))
            partition = int(int32(binary.BigEndian.Uint32(respPayload[8:12])))
        }

        // Reply to all callers in this batch
        // Offsets are assigned sequentially: last record gets the returned offset,
        // each prior one gets offset - (n - i - 1)
        for i, rec := range recs {
            result := ProduceResult{Err: err}
            if err == nil {
                result.Offset = offset - int64(len(recs)-1-i)
                result.Partition = partition
            }
            rec.result <- result
        }
    }
}
```

---

## 4. Consumer

```go
// client/consumer.go
package client

import (
    "context"
    "fmt"
    "time"

    "github.com/yourname/streamflow/wire"
)

type ConsumerConfig struct {
    Addr       string
    GroupID    string
    Topic      string
    Partition  int
    MaxBytes   int32
    PollInterval time.Duration
    AutoCommit bool
}

type Record struct {
    Offset    int64
    Timestamp time.Time
    Key       []byte
    Value     []byte
    Partition int
}

type Consumer struct {
    conn          *Conn
    cfg           ConsumerConfig
    currentOffset int64
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
    if cfg.MaxBytes == 0 {
        cfg.MaxBytes = 1024 * 1024 // 1 MB
    }
    if cfg.PollInterval == 0 {
        cfg.PollInterval = 100 * time.Millisecond
    }

    conn, err := Dial(cfg.Addr)
    if err != nil {
        return nil, err
    }

    c := &Consumer{conn: conn, cfg: cfg, currentOffset: -1}

    // Fetch committed offset for this group
    committedOffset, err := c.fetchCommittedOffset()
    if err != nil {
        conn.Close()
        return nil, err
    }
    c.currentOffset = committedOffset

    return c, nil
}

func (c *Consumer) fetchCommittedOffset() (int64, error) {
    req := wire.FetchRequest{
        Topic:     c.cfg.Topic,
        Partition: c.cfg.Partition,
        Offset:    -1, // signal: use committed offset
        MaxBytes:  0,  // peek only
        GroupID:   c.cfg.GroupID,
    }
    _, respPayload, err := c.conn.sendRecv(wire.CmdFetch, wire.EncodeFetchRequest(req))
    if err != nil {
        return 0, err
    }

    resp, err := wire.DecodeFetchResponse(respPayload)
    if err != nil {
        return 0, err
    }
    return resp.Offset, nil
}

// Poll fetches the next batch of records. Blocks until records are available
// or ctx is cancelled.
func (c *Consumer) Poll(ctx context.Context) ([]Record, error) {
    for {
        records, err := c.fetch()
        if err != nil {
            return nil, err
        }

        if len(records) > 0 {
            if c.cfg.AutoCommit && len(records) > 0 {
                last := records[len(records)-1]
                c.commitOffset(last.Offset + 1) // commit next-to-process
            }
            return records, nil
        }

        // No messages yet — wait before polling again
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(c.cfg.PollInterval):
        }
    }
}

func (c *Consumer) fetch() ([]Record, error) {
    req := wire.FetchRequest{
        Topic:     c.cfg.Topic,
        Partition: c.cfg.Partition,
        Offset:    c.currentOffset,
        MaxBytes:  c.cfg.MaxBytes,
        GroupID:   c.cfg.GroupID,
    }

    _, respPayload, err := c.conn.sendRecv(wire.CmdFetch, wire.EncodeFetchRequest(req))
    if err != nil {
        return nil, err
    }

    resp, err := wire.DecodeFetchResponse(respPayload)
    if err != nil {
        return nil, err
    }

    records := make([]Record, len(resp.Messages))
    for i, msg := range resp.Messages {
        records[i] = Record{
            Offset:    msg.Offset,
            Timestamp: time.Unix(0, msg.Timestamp),
            Key:       msg.Key,
            Value:     msg.Value,
            Partition: resp.Partition,
        }
        if msg.Offset+1 > c.currentOffset {
            c.currentOffset = msg.Offset + 1
        }
    }

    return records, nil
}

// CommitOffset manually commits the offset for this consumer group
func (c *Consumer) CommitOffset(offset int64) error {
    return c.commitOffset(offset)
}

func (c *Consumer) commitOffset(offset int64) error {
    req := wire.CommitRequest{
        GroupID:   c.cfg.GroupID,
        Topic:     c.cfg.Topic,
        Partition: c.cfg.Partition,
        Offset:    offset,
    }

    _, _, err := c.conn.sendRecv(wire.CmdCommit, wire.EncodeCommitRequest(req))
    return err
}

func (c *Consumer) Close() error { return c.conn.Close() }
```

---

## 5. A Complete Example: Order Pipeline

```go
// example/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/yourname/streamflow/client"
)

type Order struct {
    ID     string  `json:"id"`
    UserID string  `json:"user_id"`
    Amount float64 `json:"amount"`
}

const brokerAddr = "localhost:9999"
const topic = "orders"

func main() {
    // Setup: create topic
    admin, err := client.NewAdminClient(brokerAddr)
    if err != nil {
        log.Fatal(err)
    }
    admin.CreateTopic(topic, 3) // ignore error if already exists
    admin.Close()

    // Start producer goroutine
    go runProducer()

    // Start two consumers in the same group
    go runConsumer("consumer-1")
    go runConsumer("consumer-2")

    select {} // run forever
}

func runProducer() {
    p, err := client.NewProducer(client.ProducerConfig{
        Addr:         brokerAddr,
        BatchSize:    50,
        BatchTimeout: 10 * time.Millisecond,
    })
    if err != nil {
        log.Fatalf("producer: %v", err)
    }
    defer p.Close()

    for i := 0; ; i++ {
        order := Order{
            ID:     fmt.Sprintf("order-%d", i),
            UserID: fmt.Sprintf("user-%d", i%100),
            Amount: float64(i%100) * 9.99,
        }

        data, _ := json.Marshal(order)
        offset, partition, err := p.SendSync(client.ProducerRecord{
            Topic: topic,
            Key:   []byte(order.UserID), // same user → same partition
            Value: data,
        })

        if err != nil {
            log.Printf("produce error: %v", err)
        } else {
            log.Printf("[Producer] order-%d → partition=%d offset=%d", i, partition, offset)
        }

        time.Sleep(200 * time.Millisecond)
    }
}

func runConsumer(consumerID string) {
    c, err := client.NewConsumer(client.ConsumerConfig{
        Addr:       brokerAddr,
        GroupID:    "order-processors",
        Topic:      topic,
        Partition:  0, // simplified: both consumers read partition 0
        AutoCommit: false,
        MaxBytes:   512 * 1024,
    })
    if err != nil {
        log.Fatalf("%s: %v", consumerID, err)
    }
    defer c.Close()

    ctx := context.Background()
    log.Printf("[%s] started, waiting for messages", consumerID)

    for {
        records, err := c.Poll(ctx)
        if err != nil {
            log.Printf("[%s] poll error: %v", consumerID, err)
            return
        }

        for _, r := range records {
            var order Order
            json.Unmarshal(r.Value, &order)
            log.Printf("[%s] Processing order %s, amount=$%.2f (offset=%d)",
                consumerID, order.ID, order.Amount, r.Offset)

            // Simulate work
            time.Sleep(50 * time.Millisecond)
        }

        // Manual commit after processing the batch
        if len(records) > 0 {
            last := records[len(records)-1]
            c.CommitOffset(last.Offset + 1)
        }
    }
}
```

---

## Summary

- The SDK has three objects: `Conn` (raw TCP), `AdminClient` (topic management), `Producer` (async batching), `Consumer` (poll + commit).
- The producer batches records by topic+partition and flushes on batch size or timeout — exactly like Kafka's producer.
- The consumer starts from the committed offset. `AutoCommit: false` gives you manual control for at-least-once delivery.
- Manual commit after successful processing: if the process crashes after processing but before committing, the message is replayed (at-least-once, not exactly-once).

### Exercises

**Easy:** Write a program that produces 1,000 messages to a topic with 5 partitions. After all messages are produced, print how many messages landed on each partition.

**Medium:** Add a `Consumer.Seek(offset int64)` method that resets the consumer to a specific offset and sends a corresponding `CommitOffset` to the broker. Use it to "replay" messages from offset 0 after processing 100 messages.

**Hard:** Implement a connection pool in the SDK: `Pool{min, max}` that maintains min to max connections, handing out connections round-robin to producers and consumers. Benchmark single-connection vs pool-of-10 with 100 concurrent goroutines producing messages.
