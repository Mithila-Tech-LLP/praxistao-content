# Chapter 72: Building a Block Explorer Backend

Anyone who has ever pasted a transaction hash into Etherscan to see if a payment went through has used a block explorer — a website that lets a total stranger browse and search a blockchain without running a node themselves. This chapter builds the backend half of GoChain's own explorer: endpoints for browsing recent blocks, drilling into a single block's transactions, inspecting a single transaction's inputs and outputs, and searching an address's entire history, every one of them paginated so they hold up against a chain with tens of thousands of blocks.

## Table of Contents

1. [What an Explorer Backend Actually Needs](#1-what-an-explorer-backend-actually-needs)
2. [A Reusable Pagination Pattern](#2-a-reusable-pagination-pattern)
3. [Listing Recent Blocks](#3-listing-recent-blocks)
4. [Viewing a Block's Transactions](#4-viewing-a-blocks-transactions)
5. [Viewing a Single Transaction](#5-viewing-a-single-transaction)
6. [Searching an Address's History](#6-searching-an-addresss-history)
7. [Wiring the Explorer Routes](#7-wiring-the-explorer-routes)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What an Explorer Backend Actually Needs

Chapter 70's API answers narrow questions well: "what's this one block," "what's this one balance." A block explorer needs to answer broader, more human questions: "show me what's happened recently," "show me everything about this specific transaction," "show me everything that ever happened to this address." None of these are new blockchain concepts — they're new *views* over data GoChain has had all along, assembled specifically for browsing rather than for a single targeted lookup.

Think of the difference between a bank's ATM (Chapter 70's API — you ask one precise question, like "what's my balance," and get one precise answer) and a bank's paper statement (this chapter — a scrollable, chronological, browsable history designed for a human's eyes, with more context per item and support for flipping through pages). Both sit in front of the same underlying ledger; they're just different lenses onto it.

```
                     gochain/api  (Chapter 70's Server, extended)
   +----------------------------------------------------------------+
   |  Narrow, single-answer endpoints        Explorer endpoints      |
   |  (Chapter 70)                            (this chapter)          |
   |  --------------------------------        ------------------      |
   |  GET /blocks/{hash}                      GET /explorer/blocks    |
   |  GET /balance/{address}                  GET /explorer/blocks/{hash}/transactions
   |  POST /transactions                      GET /explorer/transactions/{txid}
   |                                           GET /explorer/address/{address}/transactions
   +----------------------------------------------------------------+
                              |
                    both read from the same
                    Chain, Mempool, UTXOSet
```

Every handler in this chapter is still a method on the exact same `api.Server` from Chapter 70 — no new struct, no new dependencies. We're adding capability, not architecture.

---

## 2. A Reusable Pagination Pattern

A chain with 50,000 blocks and millions of transactions cannot ever return "all of them" in one response — the response would be enormous, slow to generate, and slow for a browser to render. Every listing endpoint in this chapter accepts a `page` and `pageSize` query parameter and returns a consistent envelope so a frontend can build one reusable "next page" control instead of learning a different shape per endpoint.

```go
// gochain/api/pagination.go
package api

import (
	"net/http"
	"strconv"
)

// pageParams holds a parsed, clamped page and pageSize -- clamped so a
// client (accidentally or otherwise) can't request pageSize=1000000 and
// force the server to build a huge response.
type pageParams struct {
	Page     int
	PageSize int
}

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// parsePageParams reads ?page= and ?pageSize= from the request's query
// string, defaulting and clamping both to sane values. Invalid or
// missing values silently fall back to the default rather than
// erroring -- a malformed page number isn't worth failing a whole
// request over when "just show me page 1" is a reasonable fallback.
func parsePageParams(r *http.Request) pageParams {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return pageParams{Page: page, PageSize: pageSize}
}

// offset converts a 1-indexed page number into a 0-indexed skip count
// for the underlying storage query -- page 1 skips nothing, page 2
// skips one full page, and so on.
func (p pageParams) offset() int {
	return (p.Page - 1) * p.PageSize
}

// pagedResponse is the envelope every listing endpoint in this chapter
// returns. total is the full count of matching items (not just what's
// on this page), which is what lets a frontend compute how many pages
// exist and render page-number links.
type pagedResponse struct {
	Data     any `json:"data"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func writePaged(w http.ResponseWriter, p pageParams, total int, data any) {
	writeJSON(w, http.StatusOK, pagedResponse{
		Data:     data,
		Page:     p.Page,
		PageSize: p.PageSize,
		Total:    total,
	})
}
```

`parsePageParams` centralizes every rule about page numbers this chapter needs — defaulting, clamping to `maxPageSize` so nobody can request an enormous page in one go, and never erroring out over a bad query string. `pagedResponse` is the one envelope shape every listing endpoint below returns; a frontend developer only has to learn `{data, page, pageSize, total}` once and reuses that knowledge for blocks, block transactions, and address history alike.

---

## 3. Listing Recent Blocks

An explorer's homepage almost always starts the same way: the most recently mined blocks, newest first, with a "load more" or page-number control underneath.

```go
// gochain/api/handlers_explorer.go
package api

import "net/http"

// HandleListBlocks serves GET /explorer/blocks?page=&pageSize=. It
// returns blocks newest-first (highest height first) -- the order
// every real block explorer uses, since the newest blocks are almost
// always what a visitor cares about first.
func (s *Server) HandleListBlocks(w http.ResponseWriter, r *http.Request) {
	p := parsePageParams(r)

	// RecentBlocks walks the chain backward from the tip, using the
	// height index storage.Store already maintains (Volume 8) --
	// no need to load every block just to skip most of them.
	blocks, total, err := s.Chain.RecentBlocks(p.offset(), p.PageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list blocks")
		return
	}

	responses := make([]blockResponse, len(blocks))
	for i, b := range blocks {
		responses[i] = toBlockResponse(b)
	}
	writePaged(w, p, total, responses)
}
```

`HandleListBlocks` reuses `toBlockResponse` from Chapter 70 — the same hex-friendly view of a block, no duplicated conversion logic — and hands the actual "which blocks, in what order, how many total" work to `Chain.RecentBlocks`, which walks backward from the current tip using the height index built in Volume 8 rather than scanning every block on disk. A request for `GET /explorer/blocks?page=2&pageSize=10` on a 42-block chain returns something like:

```json
{
  "data": [
    { "height": 32, "hash": "0000ab..", "txCount": 2, "...": "..." },
    { "height": 31, "hash": "0000cd..", "txCount": 0, "...": "..." }
  ],
  "page": 2,
  "pageSize": 10,
  "total": 42
}
```

---

## 4. Viewing a Block's Transactions

Chapter 70's `HandleGetBlock` deliberately returned only transaction *IDs* inside a block, to keep that response small and predictable. An explorer's block detail page wants the opposite: full detail on every transaction in that block, ready to render directly.

```go
// HandleBlockTransactions serves GET /explorer/blocks/{hash}/transactions.
// A block's full transaction list is already sitting in memory once
// the block itself is loaded -- unlike address history (Section 6),
// this endpoint doesn't need a dedicated index, just the block.
func (s *Server) HandleBlockTransactions(w http.ResponseWriter, r *http.Request) {
	hashHex := r.PathValue("hash")
	hashBytes, err := hexDecodeOrBadRequest(w, hashHex)
	if err != nil {
		return
	}

	block, err := s.Chain.GetBlockByHash(hashBytes)
	if err != nil {
		writeError(w, http.StatusNotFound, "block not found")
		return
	}

	p := parsePageParams(r)
	txs := block.Transactions
	total := len(txs)

	start := p.offset()
	if start > total {
		start = total
	}
	end := start + p.PageSize
	if end > total {
		end = total
	}

	responses := make([]txResponse, 0, end-start)
	for _, tx := range txs[start:end] {
		responses = append(responses, toTxResponse(tx, block))
	}
	writePaged(w, p, total, responses)
}
```

Because a block's transactions already live in memory as a plain Go slice the moment the block is loaded, pagination here is ordinary slice-slicing (`txs[start:end]`) rather than a database-level query — a useful reminder that "paginate everything the same way at the API layer" doesn't mean "fetch everything the same way underneath." Most GoChain blocks will hold far fewer than a full page's worth of transactions anyway; this endpoint mostly exists for consistency and to handle the rare unusually large block gracefully. `hexDecodeOrBadRequest` is a tiny shared helper (shown in full below) that both this handler and Section 5 use to avoid repeating the same hex-decode-or-400 logic.

```go
// gochain/api/handlers_explorer.go (continued)
import (
	"encoding/hex"
	"net/http"
)

// hexDecodeOrBadRequest decodes a hex string from a path parameter,
// writing a 400 response and returning a non-nil error itself if the
// string isn't valid hex -- callers just check the returned error and
// return immediately, keeping every handler's early-exit style uniform.
func hexDecodeOrBadRequest(w http.ResponseWriter, s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		writeError(w, http.StatusBadRequest, "expected valid hex")
		return nil, err
	}
	return b, nil
}
```

---

## 5. Viewing a Single Transaction

A transaction detail page needs more than Chapter 70 ever exposed: not just "this transaction exists," but every input it spends (and, ideally, which address and amount each input came from) and every output it creates.

```go
// txResponse is the JSON shape for a single transaction's full detail --
// unlike blockResponse's txIds summary, this expands every input and
// output so an explorer's transaction page can render them directly.
type txInputResponse struct {
	PrevTxID  string `json:"prevTxId"`
	OutIndex  int    `json:"outIndex"`
	FromAddr  string `json:"fromAddress"`
	Amount    int64  `json:"amount"`
}

type txOutputResponse struct {
	ToAddress string `json:"toAddress"`
	Amount    int64  `json:"amount"`
}

type txResponse struct {
	TxID          string             `json:"txId"`
	BlockHash     string             `json:"blockHash,omitempty"`
	BlockHeight   int64              `json:"blockHeight,omitempty"`
	Confirmations int64              `json:"confirmations"`
	Inputs        []txInputResponse  `json:"inputs"`
	Outputs       []txOutputResponse `json:"outputs"`
	Fee           int64              `json:"fee"`
}

// toTxResponse builds a full explorer view of a transaction. block may
// be nil (a still-pending transaction sitting in the mempool has no
// block yet), in which case BlockHash/BlockHeight are left empty and
// Confirmations is zero.
func toTxResponse(tx *core.Transaction, block *core.Block) txResponse {
	inputs := make([]txInputResponse, len(tx.Inputs))
	var inputTotal int64
	for i, in := range tx.Inputs {
		// ResolveInput looks up the earlier output an input spends, so
		// we can show *which address* the funds came from -- an input
		// on its own only stores a reference (prev tx id + index), not
		// an address, mirroring how Bitcoin's inputs work.
		addr, amount := resolveInputSource(in)
		inputs[i] = txInputResponse{
			PrevTxID: hex.EncodeToString(in.PrevTxID),
			OutIndex: in.OutIndex,
			FromAddr: addr,
			Amount:   amount,
		}
		inputTotal += amount
	}

	outputs := make([]txOutputResponse, len(tx.Outputs))
	var outputTotal int64
	for i, out := range tx.Outputs {
		outputs[i] = txOutputResponse{ToAddress: out.Address, Amount: out.Amount}
		outputTotal += out.Amount
	}

	resp := txResponse{
		TxID:    tx.IDHex(),
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     inputTotal - outputTotal, // fee mechanics from Volume 5, Chapter 35
	}
	if block != nil {
		resp.BlockHash = hex.EncodeToString(block.Hash)
		resp.BlockHeight = block.Height
	}
	return resp
}

// HandleGetTransaction serves GET /explorer/transactions/{txid}. Unlike
// block lookups, a transaction ID alone doesn't tell us which block (if
// any) contains it, so Chain.GetTransaction does the work of searching
// the transaction index built in Volume 8 and returning both the
// transaction and its containing block in one call.
func (s *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {
	txIDHex := r.PathValue("txid")
	txIDBytes, err := hexDecodeOrBadRequest(w, txIDHex)
	if err != nil {
		return
	}

	tx, block, err := s.Chain.GetTransaction(txIDBytes)
	if err != nil {
		// Not found in a block -- check the mempool before giving up;
		// a transaction that was just submitted (Chapter 70) is real
		// and should be viewable immediately, even before it's mined.
		if pending := s.Mempool.Get(txIDBytes); pending != nil {
			writeJSON(w, http.StatusOK, toTxResponse(pending, nil))
			return
		}
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	resp := toTxResponse(tx, block)
	resp.Confirmations = s.Chain.Height() - block.Height + 1
	writeJSON(w, http.StatusOK, resp)
}
```

`toTxResponse` is the interesting part: it turns raw `Inputs`/`Outputs` (which, following the UTXO model from Volume 5, only store *references* to earlier outputs, not addresses directly) into something a human can actually read, by resolving each input back to the address and amount it originally came from via `resolveInputSource` (a thin wrapper the explorer package adds around `UTXOSet`'s lookup logic). **Confirmations** — how many blocks have been mined on top of this transaction's block, a rough proxy for "how safe is it to treat this as final" — is computed as `currentHeight - blockHeight + 1`, so a transaction in the very latest block shows `1` confirmation, matching the convention Bitcoin and Ethereum explorers both use. `HandleGetTransaction` also checks the mempool before returning `404`, so a transaction someone just submitted through Chapter 70's `POST /transactions` is visible on the explorer immediately, marked as unconfirmed (no `blockHash`, zero confirmations), rather than appearing to vanish until it's mined.

---

## 6. Searching an Address's History

The most "explorer-like" feature of all: paste in an address, see everything that ever happened to it. Unlike a balance lookup (Chapter 70), which only needs to know about *currently unspent* outputs, a full history needs every transaction that ever sent to, or spent from, that address — including ones long since spent. This is why GoChain's storage layer (Volume 8) maintains a dedicated **address index** alongside the UTXO set: a mapping from address to every transaction ID that ever touched it, updated incrementally as new blocks and transactions arrive, precisely so this endpoint never has to rescan the entire chain from block zero.

```go
// HandleAddressHistory serves GET /explorer/address/{address}/transactions.
// It returns the address's current balance alongside a paginated list
// of every transaction that ever involved it, newest first -- the two
// things anyone pasting an address into an explorer wants to see
// together on one page.
func (s *Server) HandleAddressHistory(w http.ResponseWriter, r *http.Request) {
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

	p := parsePageParams(r)
	// TransactionsByAddress delegates to the address index maintained
	// in storage.Store since Volume 8 -- an O(1) index lookup followed
	// by a bounded page of results, never a full-chain scan.
	txs, total, err := s.Chain.TransactionsByAddress(address, p.offset(), p.PageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load address history")
		return
	}

	responses := make([]txResponse, len(txs))
	for i, tx := range txs {
		block, _ := s.Chain.BlockContaining(tx)
		responses[i] = toTxResponse(tx, block)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"address":  address,
		"balance":  balance,
		"data":     responses,
		"page":     p.Page,
		"pageSize": p.PageSize,
		"total":    total,
	})
}
```

This handler is the clearest illustration in the chapter of why Volume 8's indexing work mattered beyond just speeding up balance checks: without a dedicated address index, answering "show me everything that ever happened to this address" would mean scanning every transaction in every block, every single time — fine on a toy chain with a dozen blocks, unusable on anything resembling a real network. `TransactionsByAddress` returns already-paginated results directly from the index, and the handler only has to look up each transaction's containing block (for confirmations and block linkage) before shaping the final response.

---

## 7. Wiring the Explorer Routes

Four new patterns join the router built in Chapter 70, all under an `/explorer` prefix to keep them visually and organizationally separate from the narrower Chapter 70 API, even though every handler is still just a method on the same `Server`.

```go
// gochain/api/router.go (extended)
func NewRouter(s *Server) http.Handler {
	mux := http.NewServeMux()

	// Chapter 70: JSON-RPC + core REST
	mux.HandleFunc("POST /rpc", s.HandleRPC)
	mux.HandleFunc("GET /blocks/{hash}", s.HandleGetBlock)
	mux.HandleFunc("GET /balance/{address}", s.HandleGetBalance)
	mux.HandleFunc("POST /transactions", s.HandleSendTx)

	// Chapter 71: live events
	mux.HandleFunc("GET /ws", s.HandleWS)

	// Chapter 72: block explorer backend
	mux.HandleFunc("GET /explorer/blocks", s.HandleListBlocks)
	mux.HandleFunc("GET /explorer/blocks/{hash}/transactions", s.HandleBlockTransactions)
	mux.HandleFunc("GET /explorer/transactions/{txid}", s.HandleGetTransaction)
	mux.HandleFunc("GET /explorer/address/{address}/transactions", s.HandleAddressHistory)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}
```

Here is the full picture of every endpoint this chapter (and the two before it) adds, and roughly what a browsing session through them looks like:

```
  /rpc  (POST)                     -- single JSON-RPC entry point (Ch 70)
  /blocks/{hash}         (GET)     -- one block, summary view (Ch 70)
  /balance/{address}     (GET)     -- one balance (Ch 70)
  /transactions          (POST)    -- submit a signed transaction (Ch 70)
  /ws                    (GET)     -- live event stream (Ch 71)

  /explorer/blocks                          (GET, paginated) -- recent blocks
       |
       v  (visitor clicks a block)
  /explorer/blocks/{hash}/transactions       (GET, paginated) -- that block's txs
       |
       v  (visitor clicks a transaction)
  /explorer/transactions/{txid}              (GET)            -- full tx detail
       |
       v  (visitor clicks an input/output address)
  /explorer/address/{address}/transactions   (GET, paginated) -- full address history
```

This click-path — blocks, to a block's transactions, to one transaction, to an address, and back around to that address's own transaction list — is exactly the navigation Chapter 73's frontend builds pages for, one endpoint per page.

---

## Summary

- A block explorer backend answers broader, browsing-oriented questions than Chapter 70's narrow API, but reuses the exact same `Server`, `Chain`, `Mempool`, and `UTXOSet` — it's new capability, not new architecture.
- `pageParams` and `pagedResponse` give every listing endpoint (`page`, `pageSize`, `total`) one consistent shape, clamped to a sane maximum page size so no request can force an unbounded response.
- `HandleListBlocks` returns the most recent blocks first, backed by the height index from Volume 8 rather than a full scan.
- `HandleBlockTransactions` expands a block's full transaction detail (not just IDs), paginated the same way as every other listing endpoint, even though most blocks fit on one page.
- `HandleGetTransaction` resolves a transaction's inputs back to their source addresses and amounts, computes a Bitcoin/Ethereum-style confirmation count, and checks the mempool before reporting `404` so freshly-submitted transactions are visible immediately.
- `HandleAddressHistory` depends on a dedicated address index (built in Volume 8) mapping an address to every transaction that ever touched it — without it, this endpoint would require scanning the entire chain on every request.
- All four new endpoints live under an `/explorer/...` prefix, kept visually distinct from Chapter 70's narrower API even though they share one router and one `Server`.
- The endpoint-to-endpoint click path (blocks → a block's transactions → one transaction → an address's history) mirrors exactly the page structure Chapter 73's frontend builds.

---

## Exercises

### Easy

1. Add a `GET /explorer/mempool` endpoint returning the current list of pending (not yet mined) transactions, paginated the same way as every other listing endpoint in this chapter.
2. `HandleListBlocks` currently always sorts newest-first. Add an optional `?order=asc` query parameter that returns oldest-first instead, useful for someone paging through history from genesis forward.
3. Add a `blockCount` field to `HandleAddressHistory`'s response: the number of distinct blocks (not transactions) an address has appeared in. Compute it from the already-fetched `txs` slice rather than issuing a new query.

### Medium

4. Real explorers commonly support **search-by-anything**: a single search box where a user can paste a block hash, a transaction ID, or an address, and get routed to the right detail page. Implement `GET /explorer/search?q=...` that tries, in order, "is this a valid block hash that exists," then "is this a valid transaction ID that exists," then "treat it as an address," and returns `{"type": "block" | "transaction" | "address", "result": {...}}` accordingly, or `404` if nothing matches.
5. `HandleAddressHistory` currently fetches each transaction's containing block one at a time inside a loop (`s.Chain.BlockContaining(tx)`), which means one storage lookup per transaction on the page. Profile (in a comment, reasoning from first principles is fine) why this could be slow for a busy address, then redesign it to batch those lookups — for example, by first collecting all needed block hashes and issuing one multi-get against the store instead of `pageSize` separate calls.
6. Add response caching for `HandleGetBlock` and `HandleBlockTransactions` specifically for blocks that are *not* the current chain tip — since a mined block's contents never change once several confirmations have passed, these responses are safe to cache aggressively (say, with an in-memory `map[string][]byte` guarded by a mutex, or an `ETag` header per Chapter 70's REST conventions). Explain why caching the *current tip* the same way would be a bug.

### Hard

7. Design and implement **rich pagination cursors** instead of page numbers: replace `?page=&pageSize=` with `?cursor=&limit=` where the cursor is an opaque, base64-encoded value encoding "the last height/txID seen." Update `HandleListBlocks` to use it, and explain in a comment why cursor-based pagination is more correct than page numbers for a fast-moving list where new blocks are constantly being added at the front (hint: think about what happens to page 2's contents between two requests if you use page numbers).
8. Add a `GET /explorer/stats` endpoint reporting chain-wide statistics useful for a homepage dashboard: total blocks, total transactions ever recorded, average transactions per block over the last 100 blocks, and average block time over the last 100 blocks. Discuss which of these are cheap to compute on every request versus which should be precomputed and cached, and implement whichever approach you choose for each.
9. Write a benchmark (using Go's `testing.B`) that populates an in-memory test chain with 10,000 blocks and roughly 50,000 transactions spread across 500 distinct addresses, then benchmarks `HandleAddressHistory` for one of the busier addresses. Compare the benchmark's results with and without the address index from Volume 8 wired in (i.e., falling back to a full linear scan), and report the actual speedup you measure.
