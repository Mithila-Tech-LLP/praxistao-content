# Chapter 45: VaultDB — Wire Protocol and Network Layer

Right now VaultDB only works by calling Go functions directly. Real databases accept connections from clients over a network — your Python script, your Go app, your terminal. This chapter adds a TCP server with a binary protocol.

## Table of Contents

1. What a Wire Protocol Is
2. Message Format Design
3. The Server
4. The Request/Response Protocol
5. Client Driver (Go)
6. Exercises

---

## 1. What a Wire Protocol Is

A **wire protocol** defines how bytes flow between client and server over a network connection. PostgreSQL uses its own wire protocol (pgwire) — that's why any language (Python, Go, Rust, Java) can connect to PostgreSQL with a matching driver.

We'll design a simple binary protocol for VaultDB:

```
Client sends: [query message]
Server replies: [result message] or [error message]
```

**Why binary instead of text?** Binary is compact and fast to parse. Text (like JSON over HTTP) is human-readable but slower and larger.

---

## 2. Message Format Design

Every message starts with a 5-byte header:

```
Byte 0:    Message type (1 byte)
Bytes 1-4: Payload length in bytes (4 bytes, big-endian uint32)
Bytes 5+:  Payload (variable)
```

**Message types:**

| Type | Byte | Direction | Meaning |
|------|------|-----------|---------|
| Query | `0x51` ('Q') | Client → Server | Execute a SQL query |
| ResultSet | `0x52` ('R') | Server → Client | Query returned rows |
| CommandComplete | `0x43` ('C') | Server → Client | INSERT/UPDATE/DELETE done |
| Error | `0x45` ('E') | Server → Client | Query failed |
| Ready | `0x5A` ('Z') | Server → Client | Ready for next query |

**Query message payload:**
```
[length of SQL string (4 bytes)] [SQL string (UTF-8)]
```

**ResultSet payload:**
```
[num columns (2 bytes)]
  for each column:
    [name length (1 byte)] [name bytes]
[num rows (4 bytes)]
  for each row:
    for each column:
      [type (1 byte)] [value length (4 bytes)] [value bytes]
```

---

## 3. The Server

```go
// wire/server.go
package wire

import (
    "encoding/binary"
    "fmt"
    "io"
    "log"
    "net"

    "github.com/yourname/vaultdb/query"
    "github.com/yourname/vaultdb/storage"
    "github.com/yourname/vaultdb/wal"
)

type Server struct {
    addr    string
    db      *Database
}

type Database struct {
    dm      *storage.DiskManager
    bp      *storage.BufferPool
    wal     *wal.WAL
    catalog *storage.Catalog
}

func NewServer(addr string, db *Database) *Server {
    return &Server{addr: addr, db: db}
}

func (s *Server) ListenAndServe() error {
    ln, err := net.Listen("tcp", s.addr)
    if err != nil {
        return fmt.Errorf("listen: %w", err)
    }
    log.Printf("VaultDB listening on %s", s.addr)

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("accept error: %v", err)
            continue
        }
        go s.handleConn(conn)
    }
}

func (s *Server) handleConn(conn net.Conn) {
    defer conn.Close()
    log.Printf("new connection from %s", conn.RemoteAddr())

    exec := query.NewExecutor(s.db.dm, s.db.bp, s.db.wal, s.db.catalog)

    // Send "ready" to the client
    sendReady(conn)

    for {
        msgType, payload, err := readMessage(conn)
        if err != nil {
            if err != io.EOF {
                log.Printf("read message error: %v", err)
            }
            return
        }

        switch msgType {
        case 'Q': // Query
            sql, err := decodeQueryPayload(payload)
            if err != nil {
                sendError(conn, "bad query message: "+err.Error())
                continue
            }

            result, err := executeSQL(exec, sql)
            if err != nil {
                sendError(conn, err.Error())
            } else {
                sendResult(conn, result)
            }
            sendReady(conn)
        default:
            sendError(conn, fmt.Sprintf("unknown message type: %c", msgType))
        }
    }
}

func executeSQL(exec *query.Executor, sql string) (*query.Result, error) {
    stmt, err := query.Parse(sql)
    if err != nil {
        return nil, fmt.Errorf("parse error: %w", err)
    }
    return exec.Execute(stmt)
}
```

---

## 4. The Request/Response Protocol

```go
const (
    MsgQuery           byte = 'Q'
    MsgResultSet       byte = 'R'
    MsgCommandComplete byte = 'C'
    MsgError           byte = 'E'
    MsgReady           byte = 'Z'
)

// readMessage reads one message from the connection
func readMessage(r io.Reader) (byte, []byte, error) {
    var header [5]byte
    if _, err := io.ReadFull(r, header[:]); err != nil {
        return 0, nil, err
    }

    msgType := header[0]
    payloadLen := binary.BigEndian.Uint32(header[1:5])

    if payloadLen > 64*1024*1024 { // 64 MB safety limit
        return 0, nil, fmt.Errorf("message too large: %d bytes", payloadLen)
    }

    payload := make([]byte, payloadLen)
    if _, err := io.ReadFull(r, payload); err != nil {
        return 0, nil, err
    }
    return msgType, payload, nil
}

// writeMessage writes one message to the connection
func writeMessage(w io.Writer, msgType byte, payload []byte) error {
    header := make([]byte, 5)
    header[0] = msgType
    binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
    if _, err := w.Write(header); err != nil {
        return err
    }
    if len(payload) > 0 {
        _, err := w.Write(payload)
        return err
    }
    return nil
}

func decodeQueryPayload(data []byte) (string, error) {
    if len(data) < 4 {
        return "", fmt.Errorf("query payload too short")
    }
    sqlLen := binary.BigEndian.Uint32(data[:4])
    if int(sqlLen) > len(data)-4 {
        return "", fmt.Errorf("truncated SQL string")
    }
    return string(data[4 : 4+sqlLen]), nil
}

func sendReady(w io.Writer) {
    writeMessage(w, MsgReady, nil)
}

func sendError(w io.Writer, msg string) {
    payload := make([]byte, 4+len(msg))
    binary.BigEndian.PutUint32(payload[:4], uint32(len(msg)))
    copy(payload[4:], msg)
    writeMessage(w, MsgError, payload)
}

func sendResult(w io.Writer, result *query.Result) {
    if result == nil {
        writeMessage(w, MsgCommandComplete, []byte{0, 0, 0, 0}) // affected=0
        return
    }

    if len(result.Rows) == 0 && result.Affected > 0 {
        // Command complete (INSERT/UPDATE/DELETE)
        payload := make([]byte, 8)
        binary.BigEndian.PutUint64(payload, uint64(result.Affected))
        writeMessage(w, MsgCommandComplete, payload)
        return
    }

    // Build ResultSet payload
    var payload []byte

    // Number of columns (2 bytes)
    numCols := len(result.Columns)
    payload = append(payload, byte(numCols>>8), byte(numCols))

    // Column names
    for _, col := range result.Columns {
        payload = append(payload, byte(len(col)))
        payload = append(payload, []byte(col)...)
    }

    // Number of rows (4 bytes)
    numRows := len(result.Rows)
    payload = append(payload,
        byte(numRows>>24), byte(numRows>>16), byte(numRows>>8), byte(numRows))

    // Rows
    for _, row := range result.Rows {
        for _, val := range row {
            payload = append(payload, byte(val.Type))
            payload = append(payload,
                byte(len(val.Data)>>24), byte(len(val.Data)>>16),
                byte(len(val.Data)>>8), byte(len(val.Data)))
            payload = append(payload, val.Data...)
        }
    }

    writeMessage(w, MsgResultSet, payload)
}
```

---

## 5. Client Driver (Go)

Now anyone can build a client in any language. Here's the Go driver:

```go
// client/client.go
package client

import (
    "encoding/binary"
    "fmt"
    "io"
    "net"
)

type Client struct {
    conn net.Conn
}

type Row []string

type ResultSet struct {
    Columns []string
    Rows    []Row
}

func Connect(addr string) (*Client, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }
    c := &Client{conn: conn}

    // Wait for Ready message
    msgType, _, err := c.readMessage()
    if err != nil {
        return nil, err
    }
    if msgType != 'Z' {
        return nil, fmt.Errorf("expected Ready, got %c", msgType)
    }
    return c, nil
}

func (c *Client) Query(sql string) (*ResultSet, error) {
    // Send query
    payload := make([]byte, 4+len(sql))
    binary.BigEndian.PutUint32(payload[:4], uint32(len(sql)))
    copy(payload[4:], sql)

    if err := c.writeMessage('Q', payload); err != nil {
        return nil, err
    }

    // Read response
    msgType, data, err := c.readMessage()
    if err != nil {
        return nil, err
    }

    var result *ResultSet
    switch msgType {
    case 'E':
        msgLen := binary.BigEndian.Uint32(data[:4])
        return nil, fmt.Errorf("server error: %s", string(data[4:4+msgLen]))
    case 'C':
        affected := binary.BigEndian.Uint64(data[:8])
        result = &ResultSet{Columns: []string{"affected"}, Rows: []Row{{fmt.Sprintf("%d", affected)}}}
    case 'R':
        result, err = decodeResultSet(data)
        if err != nil {
            return nil, err
        }
    default:
        return nil, fmt.Errorf("unknown response type: %c", msgType)
    }

    // Wait for Ready
    c.readMessage()

    return result, nil
}

func (c *Client) Exec(sql string) (int64, error) {
    rs, err := c.Query(sql)
    if err != nil {
        return 0, err
    }
    if len(rs.Rows) > 0 && len(rs.Rows[0]) > 0 {
        var n int64
        fmt.Sscanf(rs.Rows[0][0], "%d", &n)
        return n, nil
    }
    return 0, nil
}

func (c *Client) Close() error {
    return c.conn.Close()
}

func (c *Client) readMessage() (byte, []byte, error) {
    var header [5]byte
    if _, err := io.ReadFull(c.conn, header[:]); err != nil {
        return 0, nil, err
    }
    msgType := header[0]
    payloadLen := binary.BigEndian.Uint32(header[1:5])
    payload := make([]byte, payloadLen)
    if payloadLen > 0 {
        if _, err := io.ReadFull(c.conn, payload); err != nil {
            return 0, nil, err
        }
    }
    return msgType, payload, nil
}

func (c *Client) writeMessage(msgType byte, payload []byte) error {
    header := make([]byte, 5)
    header[0] = msgType
    binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))
    if _, err := c.conn.Write(header); err != nil {
        return err
    }
    if len(payload) > 0 {
        _, err := c.conn.Write(payload)
        return err
    }
    return nil
}

func decodeResultSet(data []byte) (*ResultSet, error) {
    off := 0
    numCols := int(binary.BigEndian.Uint16(data[off:]))
    off += 2

    rs := &ResultSet{}
    for i := 0; i < numCols; i++ {
        nameLen := int(data[off])
        off++
        rs.Columns = append(rs.Columns, string(data[off:off+nameLen]))
        off += nameLen
    }

    numRows := int(binary.BigEndian.Uint32(data[off:]))
    off += 4

    for i := 0; i < numRows; i++ {
        row := make(Row, numCols)
        for j := 0; j < numCols; j++ {
            _ = data[off] // type byte
            off++
            valLen := int(binary.BigEndian.Uint32(data[off:]))
            off += 4
            row[j] = string(data[off : off+valLen])
            off += valLen
        }
        rs.Rows = append(rs.Rows, row)
    }
    return rs, nil
}

// Usage example
func Example() {
    c, err := Connect("localhost:5555")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    c.Exec("CREATE TABLE users (id INT, name VARCHAR, age INT)")
    c.Exec("INSERT INTO users VALUES (1, 'Alice', 25)")
    c.Exec("INSERT INTO users VALUES (2, 'Bob', 30)")

    rs, _ := c.Query("SELECT * FROM users WHERE age > 24")
    for _, col := range rs.Columns {
        fmt.Printf("%-15s", col)
    }
    fmt.Println()
    for _, row := range rs.Rows {
        for _, val := range row {
            fmt.Printf("%-15s", val)
        }
        fmt.Println()
    }
}
```

---

## Summary

- A wire protocol defines the binary message format between client and server. We use a simple 5-byte header (type + length) followed by a typed payload.
- The server accepts TCP connections, reads messages in a loop, parses SQL, executes it, and sends back results.
- The client driver builds query messages, sends them, and decodes result messages.
- Any language that implements the same protocol can talk to VaultDB — Go, Python, Rust, etc.

### Exercises

**Easy:** Test the server and client: start the server in one goroutine, connect from the client, run `CREATE TABLE` + `INSERT` + `SELECT`, and verify the results.

**Medium:** Add authentication: require a password on connect. Add a handshake phase: client sends `AUTH password`, server checks against a configured password and replies with `AUTH_OK` or `AUTH_FAIL`. Only start serving queries after successful auth.

**Hard:** Add support for prepared statements: `PREPARE stmt_name FROM 'SELECT * FROM users WHERE id = ?'`, then `EXECUTE stmt_name (42)`. The server parses the SQL once and stores the AST, then on execute just substitutes the parameters and runs the executor.
