# Chapter 70: Building a JSON-RPC and REST API

Every package GoChain needed to actually *be* a blockchain is done: cryptography, blocks, consensus, transactions, wallets, networking, storage, and a virtual machine. But right now, the only way to ask "what's my balance?" or "send this transaction" is to write Go code that imports `gochain/core` directly and calls its functions. A wallet app, a block explorer website, or an exchange's backend cannot do that — they are usually written in a different language, running on a different machine, and should never need to link against GoChain's internals just to ask it a question. This chapter gives GoChain a front door: an HTTP API, speaking both JSON-RPC and REST, that any program on the internet can talk to using nothing more exotic than an HTTP client.

## Table of Contents

1. [Why a Node Needs an API](#1-why-a-node-needs-an-api)
2. [JSON-RPC and REST, Compared](#2-json-rpc-and-rest-compared)
3. [Designing the `api.Server` Type](#3-designing-the-apiserver-type)
4. [Wiring Up the HTTP Router](#4-wiring-up-the-http-router)
5. [Building the REST Endpoints](#5-building-the-rest-endpoints)
6. [Building the JSON-RPC Endpoint](#6-building-the-json-rpc-endpoint)
7. [Trying It Out with curl](#7-trying-it-out-with-curl)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why a Node Needs an API

Think about how you actually use a bank. You never walk into the vault and count your own cash — you use a teller window, an ATM, or a banking app, and all three of them speak to the same underlying ledger through a small, well-defined set of requests: "what's my balance," "move money from A to B," "show me this statement." You don't need to know how the bank's internal ledger software works, and the bank does not let you touch it directly. The teller window is the **interface**: a narrow, controlled point of contact between the outside world and a system that would be dangerous (and impractical) to expose directly.

A GoChain node needs exactly the same kind of teller window. Right now, `gochain/core`, `gochain/storage`, and `gochain/network` are Go packages — useful only to another Go program running in the same process. The moment you want:

- a web-based wallet that lets someone check their balance from a browser,
- a block explorer website (Chapter 72) that displays recent blocks to visitors,
- an exchange's backend server, written in Python or Node.js, that needs to submit withdrawal transactions,

...you need a way to reach into a running GoChain node from **outside its process**, over the network, using a protocol every programming language already knows how to speak: HTTP. That's what this chapter builds — the `gochain/api` package, a small HTTP server that sits in front of the blockchain, mempool, and UTXO set you have already built, and translates plain HTTP requests into calls against them.

```
                     Outside world (any language, any machine)
   +-----------+   +-----------+   +-----------+   +-----------+
   |  Wallet   |   |  Explorer |   |  Exchange |   |   curl /   |
   |    app    |   |  website  |   |  backend  |   |   Postman  |
   +-----+-----+   +-----+-----+   +-----+-----+   +-----+-----+
         |               |               |               |
         +-------+-------+-------+-------+-------+-------+
                 |             HTTP requests
                 v
        +-----------------------------------------+
        |          gochain/api  (this volume)      |
        |   REST endpoints  +  JSON-RPC endpoint    |
        +-------------------+-----------------------+
                             |
                 direct Go function calls
                             v
        +-----------------------------------------+
        |   gochain/core (Blockchain, Mempool)      |
        |   gochain/storage (UTXOSet)                |
        +-----------------------------------------+
```

The `gochain/api` package is deliberately thin. It does not reimplement any blockchain logic — it never recomputes a balance by hand or re-derives a hash itself. It only translates an HTTP request into a call against the `Blockchain`, `Mempool`, and `UTXOSet` types you already built, and translates their answer back into JSON. This separation matters: the same validation and consensus rules apply whether a transaction arrives from a peer over the P2P network (Volume 7) or from a wallet hitting this API — there is exactly one place those rules live, and it isn't here.

---

## 2. JSON-RPC and REST, Compared

There are two popular conventions for designing an HTTP API, and real blockchains typically expose both. It's worth understanding each on its own before writing a line of code, because they solve the same problem with a genuinely different shape.

**REST** ("Representational State Transfer") treats everything as a **resource**, identified by a URL, and uses HTTP methods as verbs on that resource:

```
GET  /blocks/000000abc...      -> fetch a block by hash
GET  /balance/1Bv11k7c...       -> fetch a balance
POST /transactions              -> create (submit) a transaction
```

**JSON-RPC** treats everything as a **remote procedure call** — you name a method and pass parameters, much like calling a local function, except the "function" runs on a different machine. There is exactly one URL; the method name lives inside the request body:

```
POST /rpc
{"jsonrpc": "2.0", "method": "getBlock", "params": ["000000abc..."], "id": 1}
```

Ethereum's `go-ethereum` node is the most famous real-world example: essentially everything it exposes (`eth_getBlockByHash`, `eth_sendRawTransaction`, `eth_getBalance`) is JSON-RPC, because a blockchain node's "API surface" is naturally a list of named operations — get this, send that — rather than a tree of nested resources the way a typical web app's REST API is. Bitcoin Core's RPC interface works the same way.

Why build both for GoChain, then? Because they serve different audiences well:

| | REST | JSON-RPC |
|---|---|---|
| **Best for** | browsers, simple `fetch()` calls, tools expecting URLs per resource | wallets and libraries that already speak "Ethereum-style" JSON-RPC |
| **Discoverable via URL?** | yes — `/blocks/{hash}` reads naturally | no — everything goes through one `/rpc` endpoint |
| **Caching-friendly** | yes (GET requests can be cached by URL) | no (everything is a POST) |
| **Matches real blockchain convention** | partially | closely (this is exactly what `go-ethereum` and Bitcoin Core do) |

GoChain exposes both from the *same* `Server` type, backed by the *same* underlying logic, so a developer can pick whichever style fits their tool.

---

## 3. Designing the `api.Server` Type

Every handler this chapter writes needs the same three things: the blockchain (to look up blocks), the mempool (to accept new transactions), and the UTXO set (to answer balance queries fast, using the index built in Volume 8 instead of rescanning the whole chain). Rather than pass all three into every single handler function, we bundle them once into a `Server` type and hang every handler off it as a method — this is the same "bundle your dependencies into a struct, then attach methods" pattern GoChain has used since `core.Blockchain` itself.

```go
// gochain/api/server.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/you/gochain/core"
	"github.com/you/gochain/storage"
)

// Server bundles everything an HTTP handler needs to answer a request
// about GoChain's state. It holds no state of its own beyond these three
// references — every handler reads from (or writes into) them, but the
// Server itself never duplicates or caches blockchain data.
type Server struct {
	Chain   *core.Blockchain
	Mempool *core.Mempool
	UTXOSet *storage.UTXOSet
}

// NewServer wires up a Server around an already-running node's chain,
// mempool, and UTXO set. Nothing here starts a listener yet -- that
// happens separately, in main() or the CLI's `node start` command
// (Chapter 74), by handing the Server's router to http.ListenAndServe.
func NewServer(chain *core.Blockchain, mempool *core.Mempool, utxoSet *storage.UTXOSet) *Server {
	return &Server{
		Chain:   chain,
		Mempool: mempool,
		UTXOSet: utxoSet,
	}
}

// writeJSON is a small shared helper every handler in this package uses
// to send a JSON response with the right status code and content type,
// so that behavior doesn't get copy-pasted (and subtly drift) across
// a dozen handler functions.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Encoding errors here are rare (they mean we tried to serialize
	// something JSON can't represent) and, by this point, the status
	// code is already written, so there's nothing more useful to do
	// than log it -- a real node would use slog here.
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a REST-style error body: {"error": "message"}.
// Every REST handler in this chapter uses this exact shape, so any
// client can rely on `error` always being present and always a string
// when a request fails.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
```

`Server` is exactly the shape agreed on across this whole volume: three fields, `Chain`, `Mempool`, and `UTXOSet`, and a `NewServer` constructor that takes them in that order. Every method from here through Chapter 71 (`HandleRPC`, `HandleGetBlock`, `HandleGetBalance`, `HandleSendTx`, and later `HandleWS`) is defined as a method on this one type, so all of them share access to the same live chain, mempool, and UTXO set without any global variables. `writeJSON` and `writeError` are small conveniences: every REST handler in this chapter ends by calling one of them, guaranteeing every response — success or failure — has a consistent shape.

---

## 4. Wiring Up the HTTP Router

Go's standard `net/http` package has, since Go 1.22, grown a router capable of matching both the HTTP method and named path segments directly in `http.ServeMux` — the same job a third-party router used to be required for. That means GoChain's API needs no extra dependency at all for basic routing.

```go
// gochain/api/router.go
package api

import "net/http"

// NewRouter builds the HTTP routing table for a Server. Each pattern
// below combines an HTTP method and a URL pattern -- "GET /blocks/{hash}"
// only matches GET requests whose path looks like /blocks/<anything>,
// and the {hash} segment is later readable via r.PathValue("hash").
func NewRouter(s *Server) http.Handler {
	mux := http.NewServeMux()

	// -- JSON-RPC: a single endpoint, method name lives in the body --
	mux.HandleFunc("POST /rpc", s.HandleRPC)

	// -- REST: one URL per resource, HTTP method as the verb --
	mux.HandleFunc("GET /blocks/{hash}", s.HandleGetBlock)
	mux.HandleFunc("GET /balance/{address}", s.HandleGetBalance)
	mux.HandleFunc("POST /transactions", s.HandleSendTx)

	// A tiny liveness check -- useful for load balancers and for your
	// own sanity when a node refuses to start.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}
```

`NewRouter` takes a `*Server` and returns an `http.Handler` — the standard interface every Go HTTP server understands, whether it comes from `net/http`, a third-party router, or (as here) directly from `http.ServeMux`. Each `mux.HandleFunc` call registers one **pattern**: the method (`GET`, `POST`) followed by a path, with `{name}` segments acting as wildcards a handler can read back out. `/health` is a small bonus endpoint with no corresponding contract method — a plain closure is fine for something this simple, since it doesn't need any of `Server`'s fields.

Starting the server, once this router exists, is a single standard-library call:

```go
// gochain/cmd/gochain/serve.go (simplified; the real version lives
// behind the `gochain node start` command built in Chapter 74)
func startAPIServer(s *api.Server, addr string) error {
	router := api.NewRouter(s)
	// http.ListenAndServe blocks forever, serving requests until the
	// process exits or an unrecoverable error occurs (like the port
	// already being in use).
	return http.ListenAndServe(addr, router)
}
```

---

## 5. Building the REST Endpoints

Each REST handler follows the same three-step shape: pull whatever the client provided (a path segment, a query parameter, a JSON body) out of the `*http.Request`, call into `Chain`, `Mempool`, or `UTXOSet` to do the real work, and write a JSON response. None of these handlers contain blockchain logic themselves — they are translators, nothing more.

### `HandleGetBlock` — `GET /blocks/{hash}`

```go
// gochain/api/handlers_rest.go
package api

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/you/gochain/core"
)

// blockResponse is the JSON shape we send back for a block -- a
// deliberately flattened, hex-friendly view of core.Block, rather than
// serializing core.Block directly. This keeps the wire format stable
// even if core.Block's internal field types change later.
type blockResponse struct {
	Height        int64    `json:"height"`
	Hash          string   `json:"hash"`
	PrevBlockHash string   `json:"prevBlockHash"`
	MerkleRoot    string   `json:"merkleRoot"`
	Timestamp     int64    `json:"timestamp"`
	Nonce         int64    `json:"nonce"`
	TxCount       int      `json:"txCount"`
	TxIDs         []string `json:"txIds"`
}

func toBlockResponse(b *core.Block) blockResponse {
	ids := make([]string, len(b.Transactions))
	for i, tx := range b.Transactions {
		ids[i] = tx.IDHex()
	}
	return blockResponse{
		Height:        b.Height,
		Hash:          hex.EncodeToString(b.Hash),
		PrevBlockHash: hex.EncodeToString(b.PrevBlockHash),
		MerkleRoot:    hex.EncodeToString(b.MerkleRoot),
		Timestamp:     b.Timestamp,
		Nonce:         b.Nonce,
		TxCount:       len(b.Transactions),
		TxIDs:         ids,
	}
}

// HandleGetBlock serves GET /blocks/{hash}. The hash arrives as a hex
// string in the URL (hashes are raw bytes internally, but hex is what
// every block explorer, wallet, and human expects to see and type).
func (s *Server) HandleGetBlock(w http.ResponseWriter, r *http.Request) {
	hashHex := r.PathValue("hash")
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		writeError(w, http.StatusBadRequest, "hash must be valid hex")
		return
	}

	block, err := s.Chain.GetBlockByHash(hashBytes)
	if err != nil {
		if errors.Is(err, core.ErrBlockNotFound) {
			writeError(w, http.StatusNotFound, "block not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load block")
		return
	}

	writeJSON(w, http.StatusOK, toBlockResponse(block))
}
```

`blockResponse` is a small but important decision: we do **not** serialize `core.Block` directly. `core.Block`'s fields are raw `[]byte` hashes, meant for hashing and storage, not for JSON — encoding raw bytes as JSON produces an unreadable base64 blob by default. `toBlockResponse` converts every hash to a familiar hex string and flattens the transaction list down to just their IDs (a client that wants full transaction detail calls the transaction endpoint from Chapter 72). `HandleGetBlock` itself reads the `{hash}` wildcard via `r.PathValue("hash")`, decodes it from hex back into the raw bytes `Chain.GetBlockByHash` expects, and maps a "not found" chain error onto an HTTP `404` — a `500` is reserved for genuinely unexpected failures, not "the block simply doesn't exist."

### `HandleGetBalance` — `GET /balance/{address}`

```go
// HandleGetBalance serves GET /balance/{address}. It never scans the
// chain itself -- it defers entirely to storage.UTXOSet, the fast index
// built in Volume 8 specifically so this kind of lookup doesn't require
// replaying every transaction ever mined.
func (s *Server) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	if address == "" {
		writeError(w, http.StatusBadRequest, "address is required")
		return
	}

	balance, err := s.UTXOSet.Balance(address)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to compute balance")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"address": address,
		"balance": balance,
	})
}
```

There's no validation here beyond "is the address non-empty" — an address that has simply never received any gochips is not an error, it's a real, valid address with a balance of zero, and `UTXOSet.Balance` returns `0, nil` for it rather than an error. This mirrors how real chains behave: querying an unused Ethereum address's balance returns `0`, not a `404`.

### `HandleSendTx` — `POST /transactions`

```go
// sendTxRequest is the JSON body a client POSTs to submit a transaction.
// The transaction itself is built and signed entirely client-side (by a
// wallet, using the crypto.Sign machinery from Volume 2) -- the API
// never sees a private key, only the already-signed result.
type sendTxRequest struct {
	Transaction core.Transaction `json:"transaction"`
}

// HandleSendTx serves POST /transactions. It accepts a fully-formed,
// already-signed transaction and hands it to the mempool, which is
// where Volume 5's double-spend and signature checks actually happen --
// this handler adds no validation of its own, so there is exactly one
// place in the whole codebase those rules live.
func (s *Server) HandleSendTx(w http.ResponseWriter, r *http.Request) {
	var req sendTxRequest
	// MaxBytesReader caps the request body so a malformed or malicious
	// client can't force us to buffer an unbounded amount of memory.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid transaction JSON: "+err.Error())
		return
	}

	if err := s.Mempool.Add(&req.Transaction); err != nil {
		// Mempool.Add returns a descriptive error for a bad signature,
		// a double-spend attempt, or a reference to a UTXO that
		// doesn't exist -- all of those are the client's fault, so
		// they map onto 400, not 500.
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"txId": req.Transaction.IDHex(),
		"status": "pending",
	})
}
```

`HandleSendTx` is intentionally the thinnest handler in the file. It decodes JSON, hands the resulting `core.Transaction` straight to `s.Mempool.Add`, and reports whatever comes back. Every check that matters — is the signature valid, does the referenced UTXO exist, is it already spent — was built in Volume 5 and lives inside `Mempool.Add` itself. The handler's only real job is turning a mempool rejection into an honest `400 Bad Request` rather than a vague `500`, since a rejected transaction is almost always the client's fault (bad signature, insufficient funds, stale UTXO reference), not the server's. The response status is `202 Accepted`, not `200 OK` — the transaction has been *accepted into the mempool*, not yet mined into a block, and `202` communicates exactly that "we got it, work is pending" state.

---

## 6. Building the JSON-RPC Endpoint

JSON-RPC needs one route, `POST /rpc`, but a bit more internal plumbing than REST, because a single handler now has to dispatch to *different* logic based on a `method` field buried inside the request body, rather than relying on the URL and HTTP method to do that job.

```go
// gochain/api/handlers_rpc.go
package api

import (
	"encoding/json"
	"net/http"
)

// rpcRequest matches the JSON-RPC 2.0 request envelope. Params is left
// as json.RawMessage (undecoded bytes) because its actual shape depends
// entirely on which method was requested -- getBlock expects a hash
// string, sendTransaction expects a whole transaction object.
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      any             `json:"id"`
}

// rpcError matches JSON-RPC 2.0's error object shape.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// rpcResponse matches the JSON-RPC 2.0 response envelope. Result and
// Error are both pointers/interfaces so that exactly one of them is
// present in the encoded JSON -- omitempty drops whichever is unused.
type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
	ID      any       `json:"id"`
}

// Standard JSON-RPC 2.0 error codes -- these exact numbers are part of
// the spec, not something we invented, so any JSON-RPC-aware client
// library already knows what -32601 means without reading GoChain docs.
const (
	rpcParseError     = -32700
	rpcInvalidRequest = -32600
	rpcMethodNotFound = -32601
	rpcInvalidParams  = -32602
	rpcInternalError  = -32603
)

func rpcErrorResponse(id any, code int, message string) rpcResponse {
	return rpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: code, Message: message}, ID: id}
}

// HandleRPC serves POST /rpc -- the single entry point for every
// JSON-RPC method GoChain exposes. It decodes the envelope, dispatches
// on the method name, and always writes back a well-formed JSON-RPC
// response, even on failure -- JSON-RPC callers expect an envelope with
// an "error" field, not a bare HTTP 400.
func (s *Server) HandleRPC(w http.ResponseWriter, r *http.Request) {
	var req rpcRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusOK, rpcErrorResponse(nil, rpcParseError, "invalid JSON"))
		return
	}

	var (
		result any
		rpcErr *rpcError
	)

	switch req.Method {
	case "getBlock":
		result, rpcErr = s.rpcGetBlock(req.Params)
	case "getBalance":
		result, rpcErr = s.rpcGetBalance(req.Params)
	case "sendTransaction":
		result, rpcErr = s.rpcSendTransaction(req.Params)
	default:
		rpcErr = &rpcError{Code: rpcMethodNotFound, Message: "unknown method: " + req.Method}
	}

	if rpcErr != nil {
		writeJSON(w, http.StatusOK, rpcResponse{JSONRPC: "2.0", Error: rpcErr, ID: req.ID})
		return
	}
	writeJSON(w, http.StatusOK, rpcResponse{JSONRPC: "2.0", Result: result, ID: req.ID})
}
```

Notice every JSON-RPC response is written with HTTP status `200 OK` — even error responses. That surprises people coming from REST, but it's a deliberate part of the JSON-RPC convention: the HTTP layer only ever reports "did the transport work," and the `error` field inside the JSON body reports "did the *call* succeed." A client can't tell the difference between "method not found" and "insufficient balance" by looking at the HTTP status alone — it has to look inside the envelope either way, so JSON-RPC keeps the transport status boring and uniform.

Now the three per-method handlers, each responsible for decoding its own `params` shape and calling exactly the same underlying `Chain`/`Mempool`/`UTXOSet` methods the REST handlers used above — proof that both API styles really are two faces on the same logic:

```go
// gochain/api/rpc_methods.go
package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/you/gochain/core"
)

// rpcGetBlock implements the "getBlock" method. Params is expected to
// be a single-element array containing a hex-encoded block hash, e.g.
// "params": ["00000abc..."] -- matching the array-of-positional-
// arguments convention most JSON-RPC methods (including Ethereum's) use.
func (s *Server) rpcGetBlock(params json.RawMessage) (any, *rpcError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		return nil, &rpcError{Code: rpcInvalidParams, Message: "expected params: [hashHex]"}
	}
	hashBytes, err := hex.DecodeString(args[0])
	if err != nil {
		return nil, &rpcError{Code: rpcInvalidParams, Message: "hash must be valid hex"}
	}
	block, err := s.Chain.GetBlockByHash(hashBytes)
	if err != nil {
		if errors.Is(err, core.ErrBlockNotFound) {
			return nil, &rpcError{Code: rpcInvalidParams, Message: "block not found"}
		}
		return nil, &rpcError{Code: rpcInternalError, Message: "failed to load block"}
	}
	return toBlockResponse(block), nil
}

// rpcGetBalance implements the "getBalance" method: params: [address].
func (s *Server) rpcGetBalance(params json.RawMessage) (any, *rpcError) {
	var args []string
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		return nil, &rpcError{Code: rpcInvalidParams, Message: "expected params: [address]"}
	}
	balance, err := s.UTXOSet.Balance(args[0])
	if err != nil {
		return nil, &rpcError{Code: rpcInternalError, Message: "failed to compute balance"}
	}
	return map[string]any{"address": args[0], "balance": balance}, nil
}

// rpcSendTransaction implements the "sendTransaction" method:
// params: [transactionObject] -- a single already-signed transaction,
// exactly like the REST /transactions body, just wrapped for JSON-RPC.
func (s *Server) rpcSendTransaction(params json.RawMessage) (any, *rpcError) {
	var args []core.Transaction
	if err := json.Unmarshal(params, &args); err != nil || len(args) != 1 {
		return nil, &rpcError{Code: rpcInvalidParams, Message: "expected params: [transaction]"}
	}
	tx := args[0]
	if err := s.Mempool.Add(&tx); err != nil {
		return nil, &rpcError{Code: rpcInvalidParams, Message: err.Error()}
	}
	return map[string]string{"txId": tx.IDHex(), "status": "pending"}, nil
}
```

Each `rpcXxx` method takes raw, undecoded `params` bytes and returns either a result (anything JSON-serializable) or an `*rpcError` — never both. `HandleRPC`'s `switch` statement is the only place that decides *which* of these to call; the methods themselves know nothing about HTTP, JSON-RPC envelopes, or status codes, which keeps them easy to unit test directly, without spinning up an HTTP server at all.

---

## 7. Trying It Out with curl

With `NewRouter(s)` wired up behind `http.ListenAndServe`, both API styles are reachable with nothing more than `curl`:

```bash
# REST: fetch a block by hash
curl http://localhost:8080/blocks/00000f3a9b...

# REST: check a balance
curl http://localhost:8080/balance/1Bv11k7cCG...

# REST: submit a transaction
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"transaction": {"id": "...", "inputs": [...], "outputs": [...]}}'

# JSON-RPC: the same three operations, one endpoint
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"getBalance","params":["1Bv11k7cCG..."],"id":1}'
```

A typical REST round trip for the balance check looks like this over the wire:

```
 Client                      Server (api.Server)              storage.UTXOSet
   |                                |                                |
   |  GET /balance/1Bv11k7c...      |                                |
   |------------------------------->|                                |
   |                                |  Balance("1Bv11k7c...")        |
   |                                |------------------------------->|
   |                                |                                | index lookup
   |                                |         balance, nil           |
   |                                |<-------------------------------|
   |     200 OK                     |                                |
   |     {"address":..,"balance":..}|                                |
   |<-------------------------------|                                |
```

And here's exactly how the two API styles map onto the same underlying calls, side by side:

```
   REST                              JSON-RPC
   ---------------------------       ---------------------------
   GET  /blocks/{hash}          ==>  {"method": "getBlock",
                                       "params": [hash]}

   GET  /balance/{address}      ==>  {"method": "getBalance",
                                       "params": [address]}

   POST /transactions           ==>  {"method": "sendTransaction",
        body: {transaction}            "params": [transaction]}

   Both sides call the same:
     s.Chain.GetBlockByHash()
     s.UTXOSet.Balance()
     s.Mempool.Add()
```

---

## Summary

- A blockchain node needs an API so programs outside its own process — wallets, explorers, exchanges — can query and interact with it over plain HTTP, without linking against `gochain/core` directly.
- REST models the world as **resources** (`/blocks/{hash}`, `/balance/{address}`) with HTTP methods as verbs; JSON-RPC models it as **named remote calls** (`getBlock`, `getBalance`, `sendTransaction`) through one endpoint. Real chains like Ethereum and Bitcoin lean heavily on JSON-RPC; GoChain offers both.
- The `api.Server` struct bundles exactly three dependencies — `Chain`, `Mempool`, `UTXOSet` — and every handler in this volume is a method on it, sharing the same live state.
- Go's standard `net/http.ServeMux`, since Go 1.22, supports method-and-path routing (`"GET /blocks/{hash}"`) and named wildcards (`r.PathValue("hash")`) with no third-party router required.
- REST handlers translate URL/query/body input into calls against `Chain`, `Mempool`, and `UTXOSet`, and map domain errors (not found, invalid transaction) onto sensible HTTP status codes — `404`, `400`, `202`, never a vague catch-all.
- JSON-RPC's `HandleRPC` always answers with HTTP `200`; success or failure is reported *inside* the JSON envelope's `result` or `error` field, per the JSON-RPC 2.0 convention.
- Neither API layer contains blockchain logic of its own — signature checks, double-spend detection, and UTXO lookups all live exactly where earlier volumes put them, keeping this package a thin, honest translator.
- `blockResponse` (and similar DTOs you'll build for transactions in Chapter 72) decouple the *wire format* clients depend on from `core.Block`'s internal representation, so internal refactors don't silently break every client.

---

## Exercises

### Easy

1. Add a REST endpoint `GET /height` that returns `{"height": N}`, the current chain height. Wire it into `NewRouter` and add a matching JSON-RPC method `getHeight` that takes no parameters.
2. `HandleGetBalance` currently treats an empty address as a `400`. Add a check that rejects addresses shorter than some minimum length you choose, with a clear error message, and write down why validating input shape (even loosely) at the API boundary is worth doing even though `UTXOSet.Balance` would also just return zero for a nonsense address.
3. The `/health` endpoint currently always returns `{"status": "ok"}`. Extend it to also report the current block height and the number of pending transactions in the mempool, so it doubles as a lightweight status page.

### Medium

4. JSON-RPC supports **batch requests**: a client can POST a JSON array of request objects instead of a single object, and expects back a JSON array of responses in the same order. Extend `HandleRPC` to detect a JSON array body and process each request in the batch, returning an array of `rpcResponse` values.
5. Add request logging middleware that wraps the `http.Handler` returned by `NewRouter`: log the method, path, status code, and duration of every request. Use `http.ResponseWriter` wrapping to capture the status code, since the standard interface doesn't expose it directly.
6. Currently, a malformed transaction JSON body in `HandleSendTx` and a mempool rejection both return `400`, with no way for a client to tell them apart programmatically. Design (and implement) a small structured REST error format — `{"error": "...", "code": "INVALID_JSON" | "MEMPOOL_REJECTED"}` — and update both REST and JSON-RPC error paths to use codes consistently.

### Hard

7. Add a simple API key mechanism: `HandleSendTx` should require a valid `X-API-Key` header for any address beyond a configurable rate limit, returning `401` otherwise, while `HandleGetBlock` and `HandleGetBalance` remain open to anyone. Discuss, in a short comment, why write operations (`sendTransaction`) are a more natural place to add access control than read operations on a public chain.
8. Implement a shared `withRecover` middleware that wraps every handler so that a panic inside one request (say, a nil-pointer bug in a future handler) is recovered, logged, and turned into a clean `500` response instead of crashing the entire node process — then write a test handler that panics on purpose and confirm the server keeps serving other requests afterward.
9. Real JSON-RPC nodes commonly namespace methods, e.g. Ethereum's `eth_getBalance`, `net_version`, `web3_clientVersion`. Refactor `HandleRPC`'s dispatch table into a registration-based design — a `map[string]func(json.RawMessage) (any, *rpcError)` built once at `NewServer` time — and add a `chain_` namespace prefix to GoChain's three methods (`chain_getBlock`, and so on) without changing any handler's internal logic.
