# Chapter 71: Real-Time Updates with WebSockets

The REST and JSON-RPC endpoints from Chapter 70 answer one question at a time: "what's the balance right now," "what's in this block right now." But a wallet showing "transaction pending..." or a block explorer's homepage both really want to know the instant something changes, not a second (or a poll interval) later. This chapter adds a WebSocket endpoint that pushes new-block and new-transaction events to every subscribed client the moment they happen, plus a small terminal program that proves it works by printing each event live as it streams in.

## Table of Contents

1. [The Problem With Polling](#1-the-problem-with-polling)
2. [WebSockets in One Paragraph](#2-websockets-in-one-paragraph)
3. [Designing an Event Hub](#3-designing-an-event-hub)
4. [Implementing `HandleWS`](#4-implementing-handlews)
5. [Broadcasting Block and Transaction Events](#5-broadcasting-block-and-transaction-events)
6. [A Terminal Client That Watches Live](#6-a-terminal-client-that-watches-live)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Problem With Polling

Imagine you are waiting for a package. **Polling** is calling the delivery company's hotline every thirty seconds to ask "has it arrived yet?" It works, eventually, but it wastes your time, wastes their operator's time, and either annoys everyone if you call too often or leaves you waiting too long if you call too rarely. **Push** is the delivery company texting you the moment the package is scanned as delivered — no wasted calls, no waiting around, and the fastest possible notice.

A wallet or block explorer polling GoChain's REST API for "is there a new block yet?" has exactly the same trade-off:

```
POLLING (what we have after Chapter 70)

  Client                          Server
    |--- GET /blocks/latest ------->|   "nothing new"
    |<------------------------------|
    |         (wait 2 seconds)      |
    |--- GET /blocks/latest ------->|   "nothing new"
    |<------------------------------|
    |         (wait 2 seconds)      |
    |--- GET /blocks/latest ------->|   "here's a new block!" (up to 2s late)
    |<------------------------------|
```

Poll too often and you flood the server with requests that almost always answer "nothing changed." Poll too rarely and users notice the lag — a wallet that takes ten seconds to notice a payment arrived feels broken, even though the payment landed on-chain instantly. Every real production blockchain node (Ethereum's `go-ethereum`, Bitcoin Core) solves this the same way: a persistent connection the server can write to whenever it wants, instead of only when asked.

---

## 2. WebSockets in One Paragraph

An HTTP request-response is like a letter: you send a request, the server sends exactly one reply, and the "conversation" ends. A **WebSocket** is like an open phone line: after a single handshake (which starts as a normal HTTP request and then "upgrades" to a different protocol), both sides can send messages to each other, at any time, for as long as the connection stays open — no new request needed for every message. GoChain uses this so that once a client subscribes, the server can shout "new block!" or "new transaction!" down the same open connection the instant it happens, with no polling loop on either end.

```
   Client                                    Server
     |                                          |
     | 1. HTTP GET /ws  (Upgrade: websocket)     |
     |----------------------------------------->|
     |                                          |
     | 2. 101 Switching Protocols                |
     |<-----------------------------------------|
     |                                          |
     |========= connection now "upgraded" ======|
     |                                          |
     |          <---  {"type":"newBlock", ...}  |  (pushed whenever mined)
     |          <---  {"type":"newTx", ...}      |  (pushed whenever received)
     |          <---  {"type":"newBlock", ...}  |
     |                                          |
```

We reach for `gorilla/websocket`, the de facto standard third-party library for WebSockets in Go, rather than hand-rolling the WebSocket wire protocol (a binary framing format defined in RFC 6455) ourselves — it is exactly the kind of well-trodden, security-sensitive parsing code not worth reimplementing.

```bash
go get github.com/gorilla/websocket
```

---

## 3. Designing an Event Hub

A single incoming connection is easy to think about, but GoChain's node will have *many* WebSocket clients connected at once — an explorer website, several wallets, a monitoring dashboard — and a new block needs to reach every single one of them, not just whichever client happened to make the most recent request. We need something that tracks every currently-connected client and can fan a single event out to all of them at once. This is the job of a **hub**: a small, mutex-protected registry of connections.

```go
// gochain/api/hub.go
package api

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// wsEvent is the JSON envelope every pushed message uses. Type tells the
// client what kind of thing happened ("newBlock" or "newTx"); Data
// carries the payload, shaped the same way the REST/JSON-RPC responses
// from Chapter 70 already are, so a client parsing one already knows
// how to parse the other.
type wsEvent struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// Hub tracks every currently-connected WebSocket client and knows how
// to fan a single event out to all of them. Each client gets its own
// buffered channel rather than being written to directly, so one slow
// or stuck client can never block delivery to every other client.
type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]chan []byte
}

func newHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]chan []byte)}
}

// register adds a newly-upgraded connection to the hub and returns the
// channel HandleWS should read from and forward onto the socket.
func (h *Hub) register(conn *websocket.Conn) chan []byte {
	ch := make(chan []byte, 16) // small buffer absorbs bursts (e.g. several events firing close together)
	h.mu.Lock()
	h.clients[conn] = ch
	h.mu.Unlock()
	return ch
}

// unregister removes a connection and closes its channel. It's safe to
// call more than once for the same connection (e.g. once when the read
// loop detects a disconnect, once from a deferred cleanup) because it
// checks the map before acting.
func (h *Hub) unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.clients[conn]; ok {
		close(ch)
		delete(h.clients, conn)
	}
}

// broadcast marshals an event once and pushes it onto every connected
// client's channel. A full channel (a client that isn't reading fast
// enough) is skipped rather than blocked on -- one unresponsive client
// must never stall event delivery to everyone else.
func (h *Hub) broadcast(event wsEvent) {
	payload, err := json.Marshal(event)
	if err != nil {
		return // an event we can't even marshal isn't worth crashing over
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.clients {
		select {
		case ch <- payload:
		default:
			// Client's buffer is full -- drop this message for them
			// rather than block the whole broadcast on one slow reader.
		}
	}
}
```

`Hub` is deliberately small: a map from an open connection to a private outbound channel, guarded by a `sync.Mutex` because multiple goroutines (one per connected client, plus whichever goroutine calls `broadcast`) touch it concurrently. The `select`/`default` pattern in `broadcast` is the important design decision — without the `default` case, sending to one full channel would block the entire loop, meaning a single slow or frozen client could silently stop *every other client* from receiving new-block events. Dropping a message for a slow client is the right trade-off for a live event feed: a client that fell behind can always re-fetch the latest state over REST.

The contract's `Server` type keeps exactly its three public fields (`Chain`, `Mempool`, `UTXOSet`); the hub is added as a small **unexported** field, purely an implementation detail of the WebSocket layer, invisible to anything outside this package:

```go
// gochain/api/server.go (extended from Chapter 70)
package api

import (
	"github.com/you/gochain/core"
	"github.com/you/gochain/storage"
)

type Server struct {
	Chain   *core.Blockchain
	Mempool *core.Mempool
	UTXOSet *storage.UTXOSet

	hub *Hub // internal: tracks connected WebSocket clients (this chapter)
}

func NewServer(chain *core.Blockchain, mempool *core.Mempool, utxoSet *storage.UTXOSet) *Server {
	return &Server{
		Chain:   chain,
		Mempool: mempool,
		UTXOSet: utxoSet,
		hub:     newHub(),
	}
}
```

---

## 4. Implementing `HandleWS`

`HandleWS` has one job beyond a normal HTTP handler: **upgrade** the incoming HTTP connection into a WebSocket connection, then keep that connection alive, forwarding whatever the hub sends it, until the client disconnects.

```go
// gochain/api/handlers_ws.go
package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader configures how an HTTP connection is promoted to a
// WebSocket. CheckOrigin is wide open here for local development --
// a production deployment (Volume 13) should restrict it to the
// explorer's own domain to prevent other websites from silently
// opening WebSocket connections to your node on a visitor's behalf.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// HandleWS serves GET /ws. Once upgraded, this handler doesn't return
// until the client disconnects -- unlike every other handler in this
// package, which returns almost immediately after writing one response.
func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade already wrote an error response to w itself on failure.
		return
	}
	defer conn.Close()

	outbound := s.hub.register(conn)
	defer s.hub.unregister(conn)

	// gorilla/websocket requires *something* to be reading incoming
	// frames at all times, even if we don't expect the client to send
	// us anything -- otherwise we'd never notice the client closing the
	// connection (a "close" frame arrives as a read error). This
	// goroutine's only job is detecting that and cleaning up.
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				s.hub.unregister(conn) // safe to call twice; see Hub.unregister
				return
			}
		}
	}()

	// The main loop: forward every event the hub sends us onto the
	// actual socket. When unregister() closes `outbound`, this range
	// exits cleanly and the handler returns, closing the connection.
	for payload := range outbound {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return
		}
	}
}
```

Three things are worth calling out by name. `upgrader.Upgrade` performs the actual HTTP-to-WebSocket handshake (the `101 Switching Protocols` response from Section 2's diagram) and hands back a `*websocket.Conn` we can read from and write to like a long-lived pipe. The background goroutine reading incoming messages exists purely to detect disconnection — GoChain's clients never need to *send* anything over this socket, but the underlying TCP connection's "the other end hung up" signal only surfaces through a failed read, so something has to be reading. The `for payload := range outbound` loop is what actually delivers events: it blocks until the hub sends something (via `broadcast`, next section) or the channel is closed (because `unregister` ran), at which point the loop — and the handler — exits.

Wiring `/ws` into the router from Chapter 70 is a one-line addition:

```go
mux.HandleFunc("GET /ws", s.HandleWS)
```

---

## 5. Broadcasting Block and Transaction Events

The hub can fan events out, and `HandleWS` can deliver them to a socket — but nothing calls `hub.broadcast` yet. Two small public methods on `Server` do that, and they're called from exactly the two places in the codebase where "something new happened" is already known: right after a transaction is accepted into the mempool, and right after a block is mined or received.

```go
// gochain/api/broadcast.go
package api

import "github.com/you/gochain/core"

// BroadcastNewBlock notifies every connected WebSocket client that a
// new block was added to the chain. Call this once, right after the
// block is accepted -- whether it was mined locally (Chapter 74's
// `gochain node start` loop) or received from a peer (Volume 7).
func (s *Server) BroadcastNewBlock(block *core.Block) {
	s.hub.broadcast(wsEvent{Type: "newBlock", Data: toBlockResponse(block)})
}

// BroadcastNewTx notifies every connected client that a new transaction
// entered the mempool. It reuses the same shape a REST client would
// see, so one JSON parser on the client side handles both.
func (s *Server) BroadcastNewTx(tx *core.Transaction) {
	s.hub.broadcast(wsEvent{Type: "newTx", Data: map[string]string{"txId": tx.IDHex()}})
}
```

`HandleSendTx` from Chapter 70 gets exactly one new line, right after a transaction clears the mempool:

```go
func (s *Server) HandleSendTx(w http.ResponseWriter, r *http.Request) {
	// ... decode req, as in Chapter 70 ...

	if err := s.Mempool.Add(&req.Transaction); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.BroadcastNewTx(&req.Transaction) // NEW: tell every subscribed client immediately

	writeJSON(w, http.StatusAccepted, map[string]string{
		"txId":   req.Transaction.IDHex(),
		"status": "pending",
	})
}
```

And the node's mining loop (built fully as part of `gochain node start` in Chapter 74) calls `BroadcastNewBlock` the moment a block is successfully appended:

```go
// inside the node's mining loop
block := chain.MineBlock()
if err := chain.AddBlock(block); err == nil {
	apiServer.BroadcastNewBlock(block) // every subscribed wallet/explorer finds out instantly
}
```

The complete path from "a block gets mined" to "a browser tab updates" now looks like this:

```
  Miner            core.Blockchain      api.Server         Hub          WebSocket clients
    |                     |                  |               |                 |
    |-- MineBlock() ----->|                  |               |                 |
    |<-- new block -------|                  |               |                 |
    |-- AddBlock(block) ->|                  |               |                 |
    |<-- ok --------------|                  |               |                 |
    |-- BroadcastNewBlock(block) ----------->|               |                 |
    |                     |                  |-- broadcast ->|                 |
    |                     |                  |               |-- push -------->| (wallet)
    |                     |                  |               |-- push -------->| (explorer)
    |                     |                  |               |-- push -------->| (gochain-watch)
```

---

## 6. A Terminal Client That Watches Live

To prove all of this actually works end to end, here's a tiny standalone program — not part of `gochain` itself, just a demonstration client — that connects to a running node's `/ws` endpoint and prints every event as it arrives.

```go
// cmd/gochain-watch/main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "gochain node HTTP address")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	fmt.Println("connecting to", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial failed:", err)
	}
	defer conn.Close()

	fmt.Println("connected -- watching for live events (Ctrl-C to quit)")

	// event mirrors api.wsEvent's JSON shape. We keep Data as raw JSON
	// here since this client doesn't need to fully decode it -- it just
	// needs to print it, so re-parsing it into a Go struct would be
	// wasted work.
	var event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("connection closed:", err)
		}
		if err := json.Unmarshal(message, &event); err != nil {
			continue // ignore anything we can't parse rather than crash the watcher
		}
		switch event.Type {
		case "newBlock":
			fmt.Println("[new block]", string(event.Data))
		case "newTx":
			fmt.Println("[new tx]   ", string(event.Data))
		default:
			fmt.Println("[unknown event]", event.Type, string(event.Data))
		}
	}
}
```

Running it against a live node produces exactly what you'd hope for — a live feed, no polling:

```
$ go run ./cmd/gochain-watch -addr localhost:8080
connecting to ws://localhost:8080/ws
connected -- watching for live events (Ctrl-C to quit)
[new tx]    {"txId":"a91f...","status":"pending"}
[new block] {"height":42,"hash":"0000ab..","txCount":3,...}
[new tx]    {"txId":"c204...","status":"pending"}
```

`gochain-watch` is deliberately minimal — no reconnection logic, no filtering — precisely so it's obvious, by reading it top to bottom, that everything interesting is happening on the *server* side. This is also exactly the pattern Chapter 73's browser-based explorer frontend reuses, just with the browser's built-in `WebSocket` object instead of `gorilla/websocket`'s Go client.

---

## Summary

- Polling wastes requests and adds latency; a WebSocket keeps one connection open so the server can push events the instant they happen, with no client asking first.
- `gorilla/websocket` handles the WebSocket wire protocol (RFC 6455 framing) so GoChain doesn't need to implement it by hand.
- A `Hub` tracks every connected client in a `map[*websocket.Conn]chan []byte`, guarded by a mutex, and fans a single event out to all of them via `broadcast`.
- `broadcast` uses `select`/`default` to skip a full channel rather than block on it — one slow client must never stall delivery to everyone else.
- `HandleWS` upgrades the HTTP connection, registers it with the hub, runs a background read loop purely to detect disconnection, and forwards hub events onto the socket until the client goes away.
- `Server.BroadcastNewBlock` and `Server.BroadcastNewTx` are the two call sites — one in the mining loop, one in `HandleSendTx` — where "something new happened" turns into a pushed event.
- The hub is added to `Server` as an unexported field, keeping the contract's three public fields (`Chain`, `Mempool`, `UTXOSet`) unchanged.
- `gochain-watch`, a tiny standalone terminal client, proves the whole pipeline works by printing every pushed event live, and previews exactly the pattern Chapter 73's browser frontend will reuse.

---

## Exercises

### Easy

1. Add a third event type, `"newPeer"`, broadcast whenever a network peer connects (Volume 7's `network.Node`). You don't need to wire it into the real networking code — just add the `Server` method, the `wsEvent` case, and update `gochain-watch` to print it distinctly.
2. `gochain-watch` currently exits the whole program if the connection drops. Change it to print an error and exit with a non-zero status code via `os.Exit(1)`, and explain in a comment why silently retrying forever might hide a real problem from someone running it interactively.
3. Add a `ClientCount() int` method to `Hub` that returns the current number of connected clients (safely, under the mutex), and expose it through the existing `/health` REST endpoint from Chapter 70 as `"wsClients": N`.

### Medium

4. Right now, `broadcast`'s dropped-message behavior is silent. Add a per-client dropped-message counter to `Hub`, incremented inside the `default` case, and expose it via a new debug endpoint `GET /debug/ws` that lists each connection's buffer size and drop count — useful for noticing a client that's chronically falling behind.
5. Implement **topic-based subscriptions**: let a client connect to `/ws?topics=newBlock` (only) or `/ws?topics=newBlock,newTx` (both), and have `HandleWS` only forward events whose `Type` is in the requested set. Default to all topics if the query parameter is missing.
6. Add a periodic **ping/pong** using `gorilla/websocket`'s `SetPingHandler` and a `time.Ticker` that sends a ping frame every 30 seconds, closing (and unregistering) any connection that fails to respond in time — this is how real WebSocket servers detect "half-open" connections where the network dropped without either side sending a proper close frame.

### Hard

7. `Hub.broadcast` currently `json.Marshal`s the event once per call, which is efficient, but every client receives every event regardless of load. Rework the hub so a burst of many rapid `BroadcastNewTx` calls (say, 50 transactions arriving in one second) gets **coalesced** into a single "here are the last N transaction IDs" message per client every 250ms, while `newBlock` events remain immediate — and explain the trade-off you're making for transaction-heavy periods.
8. Write an integration test that starts a real `httptest.Server` wrapping `NewRouter`, opens a WebSocket connection to it using `gorilla/websocket`'s dialer, calls `HandleSendTx` (via a normal HTTP POST to the same test server) with a valid transaction, and asserts the WebSocket client receives a `"newTx"` event within a short timeout — proving the full REST-to-WebSocket pipeline works without a real blockchain node running.
9. Design and implement graceful shutdown: when the node process receives `SIGINT` (Ctrl-C), it should stop accepting new WebSocket upgrades, send a proper WebSocket close frame to every currently-connected client (rather than just dropping the TCP connection), wait up to a configurable timeout for clients to disconnect cleanly, and only then let `http.ListenAndServe` return.
