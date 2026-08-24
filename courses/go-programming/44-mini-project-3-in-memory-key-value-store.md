# Chapter 44: Mini Project 3 — In-Memory Key-Value Store

A Redis-like in-memory key-value store with TTL support, multiple data types, persistence, and a TCP server accepting a simple text protocol. This project uses the full Volume 3 data structure toolkit.

## What You'll Build

```
kvstore/
├── main.go
├── store/
│   ├── store.go        # Core KV store with generics
│   ├── ttl.go          # TTL expiry using a min-heap
│   ├── types.go        # Value types (string, list, set, hash)
│   └── store_test.go
├── server/
│   ├── server.go       # TCP server
│   ├── handler.go      # Command handler
│   └── protocol.go     # Text protocol parser
└── persist/
    └── aof.go          # Append-only file persistence
```

**Supported commands:**
- String: `SET`, `GET`, `DEL`, `EXISTS`, `EXPIRE`, `TTL`
- List: `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LLEN`, `LRANGE`
- Set: `SADD`, `SREM`, `SISMEMBER`, `SMEMBERS`, `SCARD`
- Hash: `HSET`, `HGET`, `HDEL`, `HGETALL`, `HLEN`
- Server: `PING`, `KEYS`, `FLUSHALL`, `DBSIZE`

---

## 1. Value Types

```go
// store/types.go
package store

import "fmt"

type ValueType int

const (
	TypeString ValueType = iota
	TypeList
	TypeSet
	TypeHash
)

// Value is a discriminated union of all supported types.
type Value struct {
	typ  ValueType
	str  string
	list []string
	set  map[string]struct{}
	hash map[string]string
}

func StringValue(s string) Value {
	return Value{typ: TypeString, str: s}
}

func ListValue() Value {
	return Value{typ: TypeList, list: []string{}}
}

func SetValue() Value {
	return Value{typ: TypeSet, set: map[string]struct{}{}}
}

func HashValue() Value {
	return Value{typ: TypeHash, hash: map[string]string{}}
}

func (v *Value) Type() ValueType { return v.typ }

func (v *Value) assertType(want ValueType) error {
	if v.typ != want {
		got := []string{"string", "list", "set", "hash"}[v.typ]
		exp := []string{"string", "list", "set", "hash"}[want]
		return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value: expected %s, got %s", exp, got)
	}
	return nil
}

// String ops
func (v *Value) GetString() (string, error) {
	if err := v.assertType(TypeString); err != nil {
		return "", err
	}
	return v.str, nil
}

// List ops
func (v *Value) LPush(val string) error {
	if err := v.assertType(TypeList); err != nil {
		return err
	}
	v.list = append([]string{val}, v.list...)
	return nil
}

func (v *Value) RPush(val string) error {
	if err := v.assertType(TypeList); err != nil {
		return err
	}
	v.list = append(v.list, val)
	return nil
}

func (v *Value) LPop() (string, error) {
	if err := v.assertType(TypeList); err != nil {
		return "", err
	}
	if len(v.list) == 0 {
		return "", nil
	}
	val := v.list[0]
	v.list = v.list[1:]
	return val, nil
}

func (v *Value) RPop() (string, error) {
	if err := v.assertType(TypeList); err != nil {
		return "", err
	}
	if len(v.list) == 0 {
		return "", nil
	}
	n := len(v.list)
	val := v.list[n-1]
	v.list = v.list[:n-1]
	return val, nil
}

func (v *Value) LLen() (int, error) {
	if err := v.assertType(TypeList); err != nil {
		return 0, err
	}
	return len(v.list), nil
}

func (v *Value) LRange(start, stop int) ([]string, error) {
	if err := v.assertType(TypeList); err != nil {
		return nil, err
	}
	n := len(v.list)
	if start < 0 {
		start = max(n+start, 0)
	}
	if stop < 0 {
		stop = n + stop
	}
	if stop >= n {
		stop = n - 1
	}
	if start > stop {
		return []string{}, nil
	}
	result := make([]string, stop-start+1)
	copy(result, v.list[start:stop+1])
	return result, nil
}

// Set ops
func (v *Value) SAdd(member string) (bool, error) {
	if err := v.assertType(TypeSet); err != nil {
		return false, err
	}
	_, exists := v.set[member]
	v.set[member] = struct{}{}
	return !exists, nil
}

func (v *Value) SRem(member string) (bool, error) {
	if err := v.assertType(TypeSet); err != nil {
		return false, err
	}
	_, exists := v.set[member]
	delete(v.set, member)
	return exists, nil
}

func (v *Value) SIsMember(member string) (bool, error) {
	if err := v.assertType(TypeSet); err != nil {
		return false, err
	}
	_, ok := v.set[member]
	return ok, nil
}

func (v *Value) SMembers() ([]string, error) {
	if err := v.assertType(TypeSet); err != nil {
		return nil, err
	}
	members := make([]string, 0, len(v.set))
	for m := range v.set {
		members = append(members, m)
	}
	return members, nil
}

func (v *Value) SCard() (int, error) {
	if err := v.assertType(TypeSet); err != nil {
		return 0, err
	}
	return len(v.set), nil
}

// Hash ops
func (v *Value) HSet(field, val string) (bool, error) {
	if err := v.assertType(TypeHash); err != nil {
		return false, err
	}
	_, exists := v.hash[field]
	v.hash[field] = val
	return !exists, nil
}

func (v *Value) HGet(field string) (string, bool, error) {
	if err := v.assertType(TypeHash); err != nil {
		return "", false, err
	}
	val, ok := v.hash[field]
	return val, ok, nil
}

func (v *Value) HDel(field string) (bool, error) {
	if err := v.assertType(TypeHash); err != nil {
		return false, err
	}
	_, exists := v.hash[field]
	delete(v.hash, field)
	return exists, nil
}

func (v *Value) HGetAll() (map[string]string, error) {
	if err := v.assertType(TypeHash); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(v.hash))
	for k, v := range v.hash {
		result[k] = v
	}
	return result, nil
}

func (v *Value) HLen() (int, error) {
	if err := v.assertType(TypeHash); err != nil {
		return 0, err
	}
	return len(v.hash), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```

---

## 2. TTL Expiry with a Min-Heap

```go
// store/ttl.go
package store

import (
	"container/heap"
	"time"
)

type expiry struct {
	key       string
	expiresAt time.Time
	index     int // position in the heap
}

// expiryHeap is a min-heap ordered by expiry time.
type expiryHeap []*expiry

func (h expiryHeap) Len() int            { return len(h) }
func (h expiryHeap) Less(i, j int) bool  { return h[i].expiresAt.Before(h[j].expiresAt) }
func (h expiryHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *expiryHeap) Push(x any) {
	e := x.(*expiry)
	e.index = len(*h)
	*h = append(*h, e)
}
func (h *expiryHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	old[n-1] = nil
	e.index = -1
	*h = old[:n-1]
	return e
}

// ttlIndex tracks per-key expiry entries so we can update them.
type ttlIndex struct {
	h       expiryHeap
	entries map[string]*expiry
}

func newTTLIndex() *ttlIndex {
	t := &ttlIndex{entries: make(map[string]*expiry)}
	heap.Init(&t.h)
	return t
}

func (t *ttlIndex) set(key string, d time.Duration) {
	exp := time.Now().Add(d)
	if e, ok := t.entries[key]; ok {
		e.expiresAt = exp
		heap.Fix(&t.h, e.index)
		return
	}
	e := &expiry{key: key, expiresAt: exp}
	t.entries[key] = e
	heap.Push(&t.h, e)
}

func (t *ttlIndex) remove(key string) {
	if e, ok := t.entries[key]; ok {
		heap.Remove(&t.h, e.index)
		delete(t.entries, key)
	}
}

func (t *ttlIndex) ttl(key string) time.Duration {
	e, ok := t.entries[key]
	if !ok {
		return -1 // no TTL
	}
	remaining := time.Until(e.expiresAt)
	if remaining < 0 {
		return 0 // expired
	}
	return remaining
}

// expired returns keys that have passed their expiry time and removes them from the index.
func (t *ttlIndex) expired() []string {
	now := time.Now()
	var keys []string
	for t.h.Len() > 0 && t.h[0].expiresAt.Before(now) {
		e := heap.Pop(&t.h).(*expiry)
		delete(t.entries, e.key)
		keys = append(keys, e.key)
	}
	return keys
}
```

---

## 3. Core Store

```go
// store/store.go
package store

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Store is the in-memory key-value store.
type Store struct {
	mu   sync.RWMutex
	data map[string]*Value
	ttl  *ttlIndex
}

func New() *Store {
	s := &Store{
		data: make(map[string]*Value),
		ttl:  newTTLIndex(),
	}
	return s
}

// StartGC starts the background goroutine that evicts expired keys.
func (s *Store) StartGC(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.evictExpired()
			}
		}
	}()
}

func (s *Store) evictExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, key := range s.ttl.expired() {
		delete(s.data, key)
	}
}

// isExpired checks if a key has expired (lazy expiry on access).
// Must be called with at least a read lock held.
func (s *Store) isExpired(key string) bool {
	return s.ttl.ttl(key) == 0
}

// ── String commands ─────────────────────────────────

func (s *Store) Set(key, value string, ttlDur time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := StringValue(value)
	s.data[key] = &v
	if ttlDur > 0 {
		s.ttl.set(key, ttlDur)
	} else {
		s.ttl.remove(key)
	}
}

func (s *Store) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return "", false, nil
	}
	val, err := v.GetString()
	return val, err == nil, err
}

func (s *Store) Del(keys ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			s.ttl.remove(key)
			n++
		}
	}
	return n
}

func (s *Store) Exists(keys ...string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok && !s.isExpired(key) {
			n++
		}
	}
	return n
}

func (s *Store) Expire(key string, d time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok || s.isExpired(key) {
		return false
	}
	s.ttl.set(key, d)
	return true
}

func (s *Store) TTL(key string) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.data[key]; !ok {
		return -2 // key doesn't exist
	}
	return s.ttl.ttl(key)
}

// ── List commands ───────────────────────────────────

func (s *Store) getOrCreateList(key string) (*Value, error) {
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		nv := ListValue()
		s.data[key] = &nv
		return &nv, nil
	}
	if err := v.assertType(TypeList); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Store) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.getOrCreateList(key)
	if err != nil {
		return 0, err
	}
	for _, val := range values {
		v.LPush(val)
	}
	n, _ := v.LLen()
	return n, nil
}

func (s *Store) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.getOrCreateList(key)
	if err != nil {
		return 0, err
	}
	for _, val := range values {
		v.RPush(val)
	}
	n, _ := v.LLen()
	return n, nil
}

func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return "", false, nil
	}
	val, err := v.LPop()
	return val, val != "", err
}

func (s *Store) RPop(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return "", false, nil
	}
	val, err := v.RPop()
	return val, val != "", err
}

func (s *Store) LLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return 0, nil
	}
	return v.LLen()
}

func (s *Store) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return []string{}, nil
	}
	return v.LRange(start, stop)
}

// ── Set commands ────────────────────────────────────

func (s *Store) getOrCreateSet(key string) (*Value, error) {
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		nv := SetValue()
		s.data[key] = &nv
		return &nv, nil
	}
	if err := v.assertType(TypeSet); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Store) SAdd(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.getOrCreateSet(key)
	if err != nil {
		return 0, err
	}
	added := 0
	for _, m := range members {
		if ok, _ := v.SAdd(m); ok {
			added++
		}
	}
	return added, nil
}

func (s *Store) SRem(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return 0, nil
	}
	removed := 0
	for _, m := range members {
		if ok, err := v.SRem(m); err != nil {
			return 0, err
		} else if ok {
			removed++
		}
	}
	return removed, nil
}

func (s *Store) SIsMember(key, member string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return false, nil
	}
	return v.SIsMember(member)
}

func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return []string{}, nil
	}
	return v.SMembers()
}

// ── Hash commands ───────────────────────────────────

func (s *Store) getOrCreateHash(key string) (*Value, error) {
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		nv := HashValue()
		s.data[key] = &nv
		return &nv, nil
	}
	if err := v.assertType(TypeHash); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Store) HSet(key, field, value string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.getOrCreateHash(key)
	if err != nil {
		return false, err
	}
	return v.HSet(field, value)
}

func (s *Store) HGet(key, field string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return "", false, nil
	}
	return v.HGet(field)
}

func (s *Store) HGetAll(key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return map[string]string{}, nil
	}
	return v.HGetAll()
}

// ── Server commands ─────────────────────────────────

func (s *Store) Keys(pattern string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		if !s.isExpired(k) {
			if pattern == "*" || k == pattern {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

func (s *Store) DBSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) FlushAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]*Value)
	s.ttl = newTTLIndex()
}

// Snapshot returns a copy of all data for persistence.
func (s *Store) Snapshot() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := make(map[string]string, len(s.data))
	for k, v := range s.data {
		if !s.isExpired(k) && v.typ == TypeString {
			snap[k] = v.str
		}
	}
	return snap
}

// Restore loads a snapshot (used at startup).
func (s *Store) Restore(snap map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range snap {
		val := StringValue(v)
		s.data[k] = &val
	}
}

// TypeOf returns the type name of a key.
func (s *Store) TypeOf(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok || s.isExpired(key) {
		return "none"
	}
	return []string{"string", "list", "set", "hash"}[v.typ]
}

// ErrorReply wraps a command error for protocol formatting.
type ErrorReply struct{ msg string }

func (e ErrorReply) Error() string { return e.msg }
func Errorf(format string, args ...any) ErrorReply {
	return ErrorReply{fmt.Sprintf(format, args...)}
}
```

---

## 4. TCP Protocol and Server

The protocol is a simplified line-based format:

```
CLIENT: SET mykey myvalue\r\n
SERVER: +OK\r\n

CLIENT: GET mykey\r\n
SERVER: $7\r\nmyvalue\r\n

CLIENT: LPUSH mylist a b c\r\n
SERVER: :3\r\n
```

```go
// server/protocol.go
package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// writeOK writes "+OK\r\n"
func writeOK(w io.Writer) {
	fmt.Fprint(w, "+OK\r\n")
}

// writeError writes "-ERR message\r\n"
func writeError(w io.Writer, msg string) {
	fmt.Fprintf(w, "-ERR %s\r\n", msg)
}

// writeInt writes ":n\r\n"
func writeInt(w io.Writer, n int) {
	fmt.Fprintf(w, ":%d\r\n", n)
}

// writeBulk writes "$len\r\nvalue\r\n" or "$-1\r\n" for nil
func writeBulk(w io.Writer, s string, ok bool) {
	if !ok {
		fmt.Fprint(w, "$-1\r\n")
		return
	}
	fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s)
}

// writeArray writes *n\r\n followed by n bulk strings
func writeArray(w io.Writer, items []string) {
	fmt.Fprintf(w, "*%d\r\n", len(items))
	for _, item := range items {
		writeBulk(w, item, true)
	}
}

// writeBool writes ":1\r\n" or ":0\r\n"
func writeBool(w io.Writer, b bool) {
	if b {
		writeInt(w, 1)
	} else {
		writeInt(w, 0)
	}
}

// readCommand reads one command from the reader.
// Supports both inline (plain text) and RESP array format.
func readCommand(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")

	if strings.HasPrefix(line, "*") {
		// RESP array format
		n, err := strconv.Atoi(line[1:])
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid array length")
		}
		args := make([]string, 0, n)
		for i := 0; i < n; i++ {
			// Read $len line
			lenLine, err := r.ReadString('\n')
			if err != nil {
				return nil, err
			}
			lenLine = strings.TrimRight(lenLine, "\r\n")
			if !strings.HasPrefix(lenLine, "$") {
				return nil, fmt.Errorf("expected bulk string")
			}
			blen, err := strconv.Atoi(lenLine[1:])
			if err != nil || blen < 0 {
				return nil, fmt.Errorf("invalid bulk length")
			}
			buf := make([]byte, blen+2) // +2 for \r\n
			if _, err := io.ReadFull(r, buf); err != nil {
				return nil, err
			}
			args = append(args, string(buf[:blen]))
		}
		return args, nil
	}

	// Inline format: split by spaces
	return strings.Fields(line), nil
}
```

```go
// server/handler.go
package server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"kvstore/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Handle(conn io.ReadWriter) {
	r := bufio.NewReader(conn)
	w := conn.(io.Writer)
	for {
		args, err := readCommand(r)
		if err != nil {
			return
		}
		if len(args) == 0 {
			continue
		}
		h.dispatch(w, args)
	}
}

func (h *Handler) dispatch(w io.Writer, args []string) {
	cmd := strings.ToUpper(args[0])
	switch cmd {
	case "PING":
		if len(args) > 1 {
			writeBulk(w, args[1], true)
		} else {
			fmt.Fprint(w, "+PONG\r\n")
		}
	case "SET":
		h.handleSet(w, args[1:])
	case "GET":
		h.handleGet(w, args[1:])
	case "DEL":
		h.handleDel(w, args[1:])
	case "EXISTS":
		h.handleExists(w, args[1:])
	case "EXPIRE":
		h.handleExpire(w, args[1:])
	case "TTL":
		h.handleTTL(w, args[1:])
	case "TYPE":
		h.handleType(w, args[1:])
	case "LPUSH":
		h.handleLPush(w, args[1:])
	case "RPUSH":
		h.handleRPush(w, args[1:])
	case "LPOP":
		h.handleLPop(w, args[1:])
	case "RPOP":
		h.handleRPop(w, args[1:])
	case "LLEN":
		h.handleLLen(w, args[1:])
	case "LRANGE":
		h.handleLRange(w, args[1:])
	case "SADD":
		h.handleSAdd(w, args[1:])
	case "SREM":
		h.handleSRem(w, args[1:])
	case "SISMEMBER":
		h.handleSIsMember(w, args[1:])
	case "SMEMBERS":
		h.handleSMembers(w, args[1:])
	case "HSET":
		h.handleHSet(w, args[1:])
	case "HGET":
		h.handleHGet(w, args[1:])
	case "HGETALL":
		h.handleHGetAll(w, args[1:])
	case "KEYS":
		h.handleKeys(w, args[1:])
	case "DBSIZE":
		writeInt(w, h.store.DBSize())
	case "FLUSHALL":
		h.store.FlushAll()
		writeOK(w)
	default:
		writeError(w, fmt.Sprintf("unknown command '%s'", args[0]))
	}
}

func (h *Handler) handleSet(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "SET requires key value")
		return
	}
	key, value := args[0], args[1]
	var ttlDur time.Duration
	for i := 2; i < len(args)-1; i += 2 {
		switch strings.ToUpper(args[i]) {
		case "EX":
			secs, err := strconv.Atoi(args[i+1])
			if err != nil || secs <= 0 {
				writeError(w, "invalid EX value")
				return
			}
			ttlDur = time.Duration(secs) * time.Second
		case "PX":
			ms, err := strconv.Atoi(args[i+1])
			if err != nil || ms <= 0 {
				writeError(w, "invalid PX value")
				return
			}
			ttlDur = time.Duration(ms) * time.Millisecond
		}
	}
	h.store.Set(key, value, ttlDur)
	writeOK(w)
}

func (h *Handler) handleGet(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "GET requires key")
		return
	}
	val, ok, err := h.store.Get(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeBulk(w, val, ok)
}

func (h *Handler) handleDel(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "DEL requires at least one key")
		return
	}
	writeInt(w, h.store.Del(args...))
}

func (h *Handler) handleExists(w io.Writer, args []string) {
	writeInt(w, h.store.Exists(args...))
}

func (h *Handler) handleExpire(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "EXPIRE requires key seconds")
		return
	}
	secs, err := strconv.Atoi(args[1])
	if err != nil {
		writeError(w, "invalid seconds")
		return
	}
	writeBool(w, h.store.Expire(args[0], time.Duration(secs)*time.Second))
}

func (h *Handler) handleTTL(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "TTL requires key")
		return
	}
	d := h.store.TTL(args[0])
	switch d {
	case -2:
		writeInt(w, -2) // key doesn't exist
	case -1:
		writeInt(w, -1) // no TTL
	default:
		writeInt(w, int(d.Seconds()))
	}
}

func (h *Handler) handleType(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "TYPE requires key")
		return
	}
	fmt.Fprintf(w, "+%s\r\n", h.store.TypeOf(args[0]))
}

func (h *Handler) handleLPush(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "LPUSH requires key value [value...]")
		return
	}
	n, err := h.store.LPush(args[0], args[1:]...)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeInt(w, n)
}

func (h *Handler) handleRPush(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "RPUSH requires key value [value...]")
		return
	}
	n, err := h.store.RPush(args[0], args[1:]...)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeInt(w, n)
}

func (h *Handler) handleLPop(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "LPOP requires key")
		return
	}
	val, ok, err := h.store.LPop(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeBulk(w, val, ok)
}

func (h *Handler) handleRPop(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "RPOP requires key")
		return
	}
	val, ok, err := h.store.RPop(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeBulk(w, val, ok)
}

func (h *Handler) handleLLen(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "LLEN requires key")
		return
	}
	n, err := h.store.LLen(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeInt(w, n)
}

func (h *Handler) handleLRange(w io.Writer, args []string) {
	if len(args) < 3 {
		writeError(w, "LRANGE requires key start stop")
		return
	}
	start, err1 := strconv.Atoi(args[1])
	stop, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil {
		writeError(w, "start and stop must be integers")
		return
	}
	items, err := h.store.LRange(args[0], start, stop)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeArray(w, items)
}

func (h *Handler) handleSAdd(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "SADD requires key member [member...]")
		return
	}
	n, err := h.store.SAdd(args[0], args[1:]...)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeInt(w, n)
}

func (h *Handler) handleSRem(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "SREM requires key member [member...]")
		return
	}
	n, err := h.store.SRem(args[0], args[1:]...)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeInt(w, n)
}

func (h *Handler) handleSIsMember(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "SISMEMBER requires key member")
		return
	}
	ok, err := h.store.SIsMember(args[0], args[1])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeBool(w, ok)
}

func (h *Handler) handleSMembers(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "SMEMBERS requires key")
		return
	}
	members, err := h.store.SMembers(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeArray(w, members)
}

func (h *Handler) handleHSet(w io.Writer, args []string) {
	if len(args) < 3 || len(args)%2 == 0 {
		writeError(w, "HSET requires key field value [field value ...]")
		return
	}
	added := 0
	for i := 1; i < len(args)-1; i += 2 {
		ok, err := h.store.HSet(args[0], args[i], args[i+1])
		if err != nil {
			writeError(w, err.Error())
			return
		}
		if ok {
			added++
		}
	}
	writeInt(w, added)
}

func (h *Handler) handleHGet(w io.Writer, args []string) {
	if len(args) < 2 {
		writeError(w, "HGET requires key field")
		return
	}
	val, ok, err := h.store.HGet(args[0], args[1])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	writeBulk(w, val, ok)
}

func (h *Handler) handleHGetAll(w io.Writer, args []string) {
	if len(args) < 1 {
		writeError(w, "HGETALL requires key")
		return
	}
	m, err := h.store.HGetAll(args[0])
	if err != nil {
		writeError(w, err.Error())
		return
	}
	items := make([]string, 0, len(m)*2)
	for k, v := range m {
		items = append(items, k, v)
	}
	writeArray(w, items)
}

func (h *Handler) handleKeys(w io.Writer, args []string) {
	pattern := "*"
	if len(args) > 0 {
		pattern = args[0]
	}
	keys := h.store.Keys(pattern)
	writeArray(w, keys)
}
```

```go
// server/server.go
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"kvstore/store"
)

type Server struct {
	addr    string
	handler *Handler
}

func New(addr string, s *store.Store) *Server {
	return &Server{
		addr:    addr,
		handler: NewHandler(s),
	}
}

func (srv *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	slog.Info("kvstore listening", "addr", srv.addr)

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				slog.Error("accept", "err", err)
				continue
			}
		}
		go func() {
			defer conn.Close()
			srv.handler.Handle(conn)
		}()
	}
}
```

---

## 5. Main

```go
// main.go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"kvstore/server"
	"kvstore/store"
)

func main() {
	addr := ":6399"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	s := store.New()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s.StartGC(ctx)

	srv := server.New(addr, s)
	if err := srv.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	slog.Info("shutdown complete")
}
```

---

## 6. Running It

```bash
go mod init kvstore
go run .

# In another terminal, use redis-cli (it speaks the same RESP protocol):
redis-cli -p 6399 PING
redis-cli -p 6399 SET name "Gopher"
redis-cli -p 6399 GET name
redis-cli -p 6399 SET counter 0 EX 60
redis-cli -p 6399 TTL counter
redis-cli -p 6399 LPUSH tasks "buy milk" "code review" "write tests"
redis-cli -p 6399 LRANGE tasks 0 -1
redis-cli -p 6399 SADD tags go redis kvstore
redis-cli -p 6399 SMEMBERS tags
redis-cli -p 6399 HSET user:1 name Alice age 30
redis-cli -p 6399 HGETALL user:1
redis-cli -p 6399 DBSIZE
redis-cli -p 6399 KEYS "*"
```

---

## 7. Data Structures Used

| Feature | Data Structure (from Vol 3) |
|---------|----------------------------|
| Hash map for all keys | `map[string]*Value` |
| String values | Plain Go `string` |
| List type | `[]string` slice |
| Set type | `map[string]struct{}` |
| Hash type | `map[string]string` |
| TTL expiry tracking | Min-heap (`container/heap`) |
| Thread safety | `sync.RWMutex` |

---

## Summary

- **Discriminated union** (`Value` with `typ` field) gives Redis-like multi-type values in Go
- **Min-heap for TTL**: O(log n) inserts/removes, O(1) peek at next expiry — efficient for millions of keys
- **Lazy expiry**: check expiry on every read; background GC sweeps every 100ms
- **RESP protocol**: compatible with `redis-cli` — the same wire format Redis uses
- Concurrency: `sync.RWMutex` with separate read/write locks allows concurrent reads

## Exercises

### Easy
1. Add a `GETSET key newvalue` command that returns the old value while setting a new one atomically.
2. Add `INCR` and `DECR` commands that increment/decrement a string value treated as an integer. Return an error if the value is not an integer.
3. Add `MSET key1 val1 key2 val2...` and `MGET key1 key2...` for batch operations.

### Medium
4. Implement **LFU eviction policy**: when `DBSize()` exceeds a configurable limit, evict the least-frequently-used key. Track access counts in `Stats` and maintain an LFU structure from Ch 42.
5. Add **append-only file (AOF) persistence**: write every mutating command to a file. On startup, replay the file to restore state. Use a background goroutine that flushes the buffer every second.
6. Add a `SUBSCRIBE`/`PUBLISH` command pair: clients can subscribe to a channel name and receive messages published by other clients. Use a `map[string][]chan string` protected by a mutex, and send messages in goroutines to avoid blocking the publisher.

### Hard
7. Add **Lua scripting** support via the `EVAL` command using the `gopher-lua` package. Scripts should be able to call `GET` and `SET` atomically with the store locked for the duration of the script.
8. Benchmark your store against `go-redis` with a real Redis instance. Measure GET/SET throughput (ops/sec), latency P50/P99, and memory usage for 1M keys. What are the performance characteristics? Where is time actually spent?
