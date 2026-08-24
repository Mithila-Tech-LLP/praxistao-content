# Chapter 90: Setting Up a Public Testnet and a Faucet

Chapter 88 got GoChain running on a real, internet-reachable server. Chapter 89 took a detour through Kubernetes for anyone curious about running this at larger scale. But a server that only *you* know the address of is not yet a public testnet — it is just your private network with a fancier address. This chapter does two things: publishes the information a stranger needs to actually join your network, and builds a **faucet**, a small web service that hands out free test `gochip`s so a brand-new participant can start sending transactions immediately, without first having to mine a block themselves.

## Table of Contents

1. [What Makes a Testnet "Public"](#1-what-makes-a-testnet-public)
2. [Publishing Seed-Node Addresses](#2-publishing-seed-node-addresses)
3. [How an External Node Joins and Syncs](#3-how-an-external-node-joins-and-syncs)
4. [Why New Participants Need a Faucet](#4-why-new-participants-need-a-faucet)
5. [Designing the Faucet's Rules](#5-designing-the-faucets-rules)
6. [Building the Faucet Wallet](#6-building-the-faucet-wallet)
7. [The Rate Limiter](#7-the-rate-limiter)
8. [The Faucet HTTP Handler](#8-the-faucet-http-handler)
9. [Wiring Up the Faucet Server](#9-wiring-up-the-faucet-server)
10. [Adding the Faucet to Docker Compose](#10-adding-the-faucet-to-docker-compose)
11. [Testing the Faucet End to End](#11-testing-the-faucet-end-to-end)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What Makes a Testnet "Public"

A **testnet** is a blockchain network that behaves exactly like a real production network — real mining, real transactions, real peer-to-peer gossip — except the coins on it are worthless by design, meant purely for testing and experimentation rather than holding real value. Nearly every major blockchain project (Bitcoin's testnet, Ethereum's Sepolia and Holesky) maintains one for exactly this reason: developers need somewhere to try things that might break, without risking real money.

What actually makes a testnet **public**, as opposed to merely "running on a server," is two things working together: a stranger's GoChain node needs to be able to *find* it, and once found, that stranger needs a reasonable way to *get started* using it. Chapter 88 solved the first half only partially — the server is reachable, but nobody except you knows its address yet. This chapter solves the rest.

---

## 2. Publishing Seed-Node Addresses

Recall Chapter 47's seed-node pattern: a new node needs at least one already-known peer address to bootstrap from, after which peer-exchange messages let it discover the rest of the network organically. Up to now, every seed address used in this course (`node1:9000` in Chapter 87's Compose network, `127.0.0.1:53001` in Chapter 52's manual multi-terminal setup) has only ever been meaningful to processes that already knew about each other. A public testnet needs its seed address written down somewhere a total stranger can find it.

Concretely, this means documenting, in your project's README or a dedicated `TESTNET.md` file, the exact information anyone needs to connect:

```markdown
# GoChain Public Testnet

Seed node address: 203.0.113.42:9000
Explorer:           http://203.0.113.42:8090   (HTTPS version: Chapter 93)
Faucet:             http://203.0.113.42:8090/faucet  (this chapter)

To join with your own GoChain node:

    gochain node start \
      --api-addr 0.0.0.0:8080 \
      --p2p-addr 0.0.0.0:9000 \
      --seed 203.0.113.42:9000 \
      --data ./data
```

This is genuinely the entire "publishing" step — no special registration process, no central authority to notify. The seed address is just a fact about the world (a reachable IP and port) written down somewhere discoverable, exactly the same way Bitcoin's own client ships with a small hardcoded list of long-lived, well-known node addresses to try first. Anyone who reads it and runs the command above has, from GoChain's own code's perspective, done exactly what `node2` and `node3` did automatically inside Chapter 87's Compose network — dialed a known address and started the Chapter 47 handshake.

---

## 3. How an External Node Joins and Syncs

Walking through exactly what happens, end to end, the moment a stranger — call her Priya, running GoChain on her own laptop for the first time — points her node at your published seed address:

```
  Priya's laptop                              Your public testnet (from the internet)
 +------------------+                        +----------------------------------------+
 | gochain node      |                        |  node1 (seed)   node2      node3       |
 | start             |                        |  :8080 :9000    :8081/9001 :8082/9002  |
 | --seed            |   1. TCP dial          |                                        |
 |  203.0.113.42:9000|  --------------------> |  node1 accepts the connection            |
 |                    |                        |                                        |
 |                    |  2. handshake          |                                        |
 |                    |  <-------------------> |  (Chapter 47's version/verack exchange) |
 |                    |                        |                                        |
 |                    |  3. "what's your       |                                        |
 |                    |     height?"           |                                        |
 |                    |  --------------------> |  node1: "height 4,213"                  |
 |                    |                        |                                        |
 |                    |  4. request missing    |                                        |
 |                    |     blocks in order    |                                        |
 |                    |  <-------------------- |  node1 streams blocks 1..4213           |
 |                    |                        |                                        |
 |                    |  5. validate each       |                                        |
 |                    |     block as it        |                                        |
 |                    |     arrives            |                                        |
 |                    |     (Chapter 19/49)    |                                        |
 |                    |                        |                                        |
 |                    |  6. peer-exchange:      |                                        |
 |                    |     node1 shares        |                                        |
 |                    |     node2/node3's       |                                        |
 |                    |     addresses too       |                                        |
 |                    |  <-------------------- |                                        |
 |                    |                        |                                        |
 |  now connected to  |                        |                                        |
 |  node1, node2,     |                        |                                        |
 |  node3 - a full     |                        |                                        |
 |  peer, not a         |                        |                                        |
 |  second-class client|                        |                                        |
 +--------------------+                        +----------------------------------------+
```

Nothing here is new code — this is Chapter 47's peer discovery, Chapter 49's synchronization, and Chapter 19's block validation, running exactly as they always have. What is new is that the "other side" of the connection is a genuine stranger's machine on a different continent instead of another process on your own laptop. If GoChain's networking and validation code is correctly written, it cannot tell the difference — and that indifference is precisely the point of building a peer-to-peer system this way from Chapter 43 onward: correctness never depended on trusting *who* you were talking to, only on independently verifying *what* they sent you.

Once synced, Priya's node is not a passive observer — it can mine its own blocks, gossip transactions it hears about, and serve chain data to the *next* stranger who points a node at *it*, exactly the organically-growing network topology Chapter 47 described.

---

## 4. Why New Participants Need a Faucet

Priya's node is now fully synced and participating in consensus. But she cannot do anything useful with it yet — sending a transaction requires spending existing UTXOs, and her brand-new wallet address owns none. She could mine a block herself and collect the coinbase reward from Chapter 37, but proof-of-work mining against an established network's current difficulty (Chapter 26) could take a genuinely long time on a single laptop competing against your testnet's existing miners.

A **faucet** is the standard blockchain-testnet solution to exactly this bootstrapping problem: a service, typically reachable over plain HTTP, that holds a supply of test coins and gives out a small, fixed amount to any address that asks — no mining, no proof of work, no waiting. Every well-known public testnet (Ethereum's Sepolia faucet, Bitcoin's testnet faucets) works this way, and it exists purely to remove friction for new participants and for your own future demos, never as a security-relevant part of the network itself.

---

## 5. Designing the Faucet's Rules

Before writing code, the faucet needs a small, explicit policy — otherwise nothing stops one address (or one person, using many addresses) from repeatedly draining it:

- **A fixed payout per request** — every request gets exactly the same amount, say `10 gochips`, removing any incentive to game the amount requested.
- **A per-address cooldown** — the same address cannot request again until a fixed amount of time has passed (this chapter uses 24 hours), which stops a single address from being drained continuously.
- **A per-IP-address cooldown**, as a second, independent check — since an attacker can generate unlimited new addresses for free (Chapter 13's key generation costs nothing), address-based limiting alone is not enough; limiting by the *requester's* network origin closes that gap, at the cost of occasionally being too strict for two legitimate users behind the same shared IP (a fair, well-understood trade-off every public faucet makes).
- **A maximum total supply the faucet will ever give out**, as a final backstop, so a bug elsewhere in the rate limiter cannot silently drain funds forever.

---

## 6. Building the Faucet Wallet

The faucet needs its own dedicated wallet — never reuse a wallet that also does anything else, so its blast radius if compromised is limited to "an attacker can give out free testnet coins faster than intended," not anything worse:

```bash
gochain-wallet new --output faucet-wallet.json
# Wallet created.
# Address: gcT1qFaucetAddr9k2m...
```

The faucet's own address needs to actually hold coins before it can give any away — either mined directly (run the faucet's address through a few rounds of Chapter 25's mining) or, more simply, granted a starting balance directly in your testnet's genesis block (Chapter 18), the same way a real testnet's genesis allocation typically pre-funds a small number of well-known addresses.

---

## 7. The Rate Limiter

The rate limiter is the faucet's core safety mechanism — a small, self-contained package that tracks the last time each address and each IP successfully received funds, and rejects a request if either cooldown has not yet elapsed:

```go
// gochain/faucet/limiter.go
//
// Package faucet implements a rate-limited coin dispenser for GoChain's
// public testnet. Its only job is deciding who is allowed a payout right
// now - it knows nothing about transactions or wallets itself.
package faucet

import (
	"sync"
	"time"
)

// Limiter tracks the last successful payout time per address and per
// requester IP, independently. Both checks must pass for a new request
// to be allowed - this is what stops both "one address asking
// repeatedly" and "one person generating unlimited free addresses."
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	byAddr   map[string]time.Time
	byIP     map[string]time.Time
}

// NewLimiter creates a rate limiter with a fixed cooldown - this chapter
// uses 24 hours, a reasonable balance between "usable for a legitimate
// developer testing repeatedly over a few days" and "not drainable in a
// tight loop."
func NewLimiter(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		byAddr:   make(map[string]time.Time),
		byIP:     make(map[string]time.Time),
	}
}

// Allow reports whether a request from the given address and IP should
// be granted right now, and if so, records both as having just been
// paid out - the same lock protects the check and the record so two
// concurrent requests can never both slip through a half-updated map.
func (l *Limiter) Allow(address, ip string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if last, ok := l.byAddr[address]; ok {
		if remaining := l.cooldown - now.Sub(last); remaining > 0 {
			return false, remaining
		}
	}
	if last, ok := l.byIP[ip]; ok {
		if remaining := l.cooldown - now.Sub(last); remaining > 0 {
			return false, remaining
		}
	}

	l.byAddr[address] = now
	l.byIP[ip] = now
	return true, 0
}
```

Recording the timestamp for *both* keys happens only once a request is actually allowed, and happens atomically under the same lock as the check itself — a classic check-then-act pattern that, done separately (check, then later record), would let two nearly-simultaneous requests both slip past the check before either recorded anything, defeating the whole limiter. This is the same category of bug Chapter 34's mempool double-spend check has to guard against, solved the same way: hold one lock across the entire read-then-write sequence.

---

## 8. The Faucet HTTP Handler

With the rate limiter in place, the HTTP handler itself is a thin layer that validates the request, checks the limiter, and — if allowed — asks the faucet's own wallet to build, sign, and submit a transaction to the local GoChain node's API, reusing the exact `POST /tx/send`-style endpoint from Chapter 70 that any other client would use:

```go
// gochain/faucet/handler.go
package faucet

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/you/gochain/wallet"
)

// payoutAmount is fixed and identical for every request, by design -
// see Section 5's reasoning against a variable or requester-chosen
// amount.
const payoutAmount = 10 // gochips

type requestBody struct {
	Address string `json:"address"`
}

type responseBody struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	TxID    string `json:"txId,omitempty"`
}

// Server ties together the rate limiter, the faucet's own funded
// wallet, and the node API it submits transactions through.
type Server struct {
	limiter *Limiter
	wallet  *wallet.Wallet
	nodeURL string // e.g. "http://localhost:8080" - the faucet's own
	               // node, reached the same way any external client
	               // would reach it
}

func NewServer(limiter *Limiter, w *wallet.Wallet, nodeURL string) *Server {
	return &Server{limiter: limiter, wallet: w, nodeURL: nodeURL}
}

// HandleRequest is registered at POST /faucet/request. It expects a
// JSON body of {"address": "gcT1q..."} and returns a JSON result
// explaining whether the payout succeeded, and if not, why - including
// exactly how long the requester must wait before trying again.
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Address == "" {
		writeJSON(w, http.StatusBadRequest, responseBody{
			OK:      false,
			Message: "missing or invalid \"address\" field",
		})
		return
	}

	// RemoteAddr includes a port ("203.0.113.7:54321"); SplitHostPort
	// strips it, since limiting by IP should ignore the ephemeral
	// source port a client's OS happened to pick for this connection.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	allowed, remaining := s.limiter.Allow(body.Address, ip)
	if !allowed {
		writeJSON(w, http.StatusTooManyRequests, responseBody{
			OK:      false,
			Message: "rate limited, try again in " + remaining.Round(1e9).String(),
		})
		return
	}

	txID, err := s.sendPayout(body.Address)
	if err != nil {
		log.Printf("faucet: payout to %s failed: %v", body.Address, err)
		writeJSON(w, http.StatusInternalServerError, responseBody{
			OK:      false,
			Message: "internal error sending payout",
		})
		return
	}

	writeJSON(w, http.StatusOK, responseBody{
		OK:      true,
		Message: "sent " + itoa(payoutAmount) + " gochips",
		TxID:    txID,
	})
}

func writeJSON(w http.ResponseWriter, status int, body responseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
```

`sendPayout` (a small helper, using the same `wallet.Send`-style flow Chapter 36's CLI wallet already implements) builds a transaction from the faucet's wallet to the requested address, signs it with the faucet's private key, and posts it to the node's `/tx/send` endpoint — nothing about that inner mechanism is new; the faucet is simply an automated, rate-limited caller of the exact same wallet-and-API machinery a human already uses through `gochain-wallet send`.

Notice, too, what the handler deliberately does *not* do: it never lets the caller specify an amount, never trusts any client-supplied signature (the faucet signs its own transaction with its own key — the requester only ever supplies a destination address), and never leaks internal error detail to the response body, logging it server-side instead.

---

## 9. Wiring Up the Faucet Server

```go
// cmd/gochain-faucet/main.go
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/you/gochain/faucet"
	"github.com/you/gochain/wallet"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0:8091", "address for the faucet HTTP server")
	nodeURL := flag.String("node", "http://localhost:8080", "GoChain node API to submit payouts through")
	walletFile := flag.String("wallet", "faucet-wallet.json", "path to the faucet's own wallet file")
	flag.Parse()

	w, err := wallet.Load(*walletFile)
	if err != nil {
		log.Fatalf("faucet: failed to load wallet: %v", err)
	}

	limiter := faucet.NewLimiter(24 * time.Hour)
	server := faucet.NewServer(limiter, w, *nodeURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/faucet/request", server.HandleRequest)

	log.Printf("faucet listening on %s, funding from %s, paying out via node %s", *listenAddr, w.Address(), *nodeURL)
	log.Fatal(http.ListenAndServe(*listenAddr, mux))
}
```

The faucet runs as its own small, independent binary — a fifth participant in the testnet alongside the three nodes and the explorer, exactly the kind of standalone tool the `cmd/` layout from Chapter 06 was designed to keep separate from GoChain's core library code.

---

## 10. Adding the Faucet to Docker Compose

Following the exact multi-stage pattern from Chapter 86 and the `Dockerfile.explorer` precedent from Chapter 87, add a small `Dockerfile.faucet` and a new service to `docker-compose.yml`:

```dockerfile
# Dockerfile.faucet
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/gochain-faucet ./cmd/gochain-faucet

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/gochain-faucet /usr/local/bin/gochain-faucet
COPY faucet-wallet.json /app/faucet-wallet.json
EXPOSE 8091
ENTRYPOINT ["gochain-faucet"]
CMD ["--listen", "0.0.0.0:8091", "--node", "http://node1:8080", "--wallet", "/app/faucet-wallet.json"]
```

```yaml
# docker-compose.yml — addition for Chapter 90
services:
  # ... existing node1, node2, node3, explorer from Chapter 87 ...

  faucet:
    build:
      context: .
      dockerfile: Dockerfile.faucet
    container_name: gochain-faucet
    # gochain-faucet talks to node1's API by service name, exactly the
    # same DNS mechanism the explorer already relies on.
    ports:
      - "8091:8091"
    networks:
      - gochain-net
    depends_on:
      - node1
```

Note that port 8091 is not opened in Chapter 88's firewall alongside 8080/9000 — like the explorer's 8090, it stays reachable only from inside the Compose network and via SSH tunnel until Chapter 93 puts it behind a proper domain and reverse proxy, rather than exposing a bare, unauthenticated port straight to the internet.

---

## 11. Testing the Faucet End to End

```bash
curl -X POST http://localhost:8091/faucet/request \
  -H "Content-Type: application/json" \
  -d '{"address": "gcT1qNewUserAddr..."}'

# {"ok":true,"message":"sent 10 gochips","txId":"a1b2c3d4..."}
```

Mine the pending transaction (or wait for a testnet miner to pick it up) and confirm the new address's balance:

```bash
curl http://localhost:8080/mine
curl "http://localhost:8080/balance?address=gcT1qNewUserAddr..."
# {"balance": 10}
```

Then confirm the rate limiter actually works, by immediately requesting again from the same address:

```bash
curl -X POST http://localhost:8091/faucet/request \
  -H "Content-Type: application/json" \
  -d '{"address": "gcT1qNewUserAddr..."}'

# {"ok":false,"message":"rate limited, try again in 23h59m58s"}
```

That second response is the entire faucet policy from Section 5 working exactly as designed: the same address, requesting again seconds later, is turned away with a clear, honest explanation rather than silently failing or being drained repeatedly.

---

## Summary

- A public testnet requires both a *findable* seed address and a low-friction way for newcomers to get started — Chapter 88 alone only solved reachability.
- Publishing a seed-node address is just documentation: a written-down IP:port that anyone can pass to `gochain node --seed`, exactly mirroring Chapter 47's peer-discovery design.
- An external node joining runs through the exact same handshake (Chapter 47), sync (Chapter 49), and validation (Chapter 19) code as any other GoChain node — the network layer was built from Chapter 43 onward to never care who, geographically, is on the other end of a connection.
- A faucet solves the cold-start problem: a new participant cannot spend UTXOs they do not have, and mining against an established difficulty is impractical for onboarding.
- The faucet's policy is a fixed payout amount, a per-address cooldown, and a per-IP cooldown — the second check exists because new addresses are free to generate, so address-only limiting is not enough.
- The rate limiter's check-and-record step must happen under a single lock, exactly like the mempool's double-spend check, or concurrent requests can both slip through.
- The faucet is just an automated, rate-limited caller of the same wallet-signing and `/tx/send` API flow a human already uses through `gochain-wallet` — no new transaction machinery, only a policy layer in front of it.
- The faucet runs as its own container in `docker-compose.yml`, reachable only internally (and via SSH tunnel) until Chapter 93 puts it behind a real domain and HTTPS.

---

## Exercises

### Easy

1. Generate a dedicated faucet wallet with `gochain-wallet new`, fund it via your testnet's genesis allocation or a few mined blocks, and start `gochain-faucet` pointed at your local node. Confirm `curl -X POST .../faucet/request` with a valid address returns `{"ok":true,...}`.

2. Immediately repeat the same request from Exercise 1 with the same address, and confirm you get a `429`-style rate-limited response with an accurate remaining-cooldown duration.

3. Send a request with a missing `address` field (`{}`) and confirm the handler returns a `400 Bad Request` with a clear error message, rather than crashing or hanging.

### Medium

4. Add a maximum-total-supply backstop to the `Limiter` (or `Server`): track cumulative gochips paid out, and reject all further requests once a configured ceiling is reached, returning a distinct message ("faucet is empty, please contact the operator") rather than the ordinary rate-limit message.

5. Write a small load-testing script (a shell loop with `curl`, or a short Go program) that fires 20 concurrent faucet requests from 20 different randomly generated addresses but all from the same machine (same IP). Confirm the per-IP cooldown correctly rejects all but the first, even under real concurrency, and explain what would go wrong if the `Allow` method's check and record were not protected by the same lock.

6. Extend the faucet's rate limiter to persist its cooldown records to disk (a simple JSON file, or a BoltDB bucket reusing Chapter 54's storage engine) so a faucet restart does not silently reset everyone's cooldown to zero. Test this by requesting, restarting the faucet process, and confirming the same address is still rate-limited afterward.

### Hard

7. Recruit a friend (or use a second machine/VPN connection you control) to actually join your public testnet as an external node, following Section 2's published instructions. Capture logs from both sides showing the handshake, sync, and peer-exchange steps from Section 3 actually happening, and confirm their node's chain height converges with yours.

8. Add a CAPTCHA-style or proof-of-work challenge (a small, cheap "solve this before your request is accepted" step) to the faucet's HTTP handler, as an additional defense beyond IP/address rate limiting against a scripted attacker rotating through many IP addresses (for example, via a botnet or proxy pool). Explain what class of abuse this defends against that per-IP limiting alone does not.

9. Design and implement a `/faucet/stats` endpoint that reports, without exposing individual requesters' addresses or IPs, aggregate faucet health: total gochips paid out, number of unique addresses served in the last 24 hours, and remaining balance. Explain what information you deliberately excluded from this endpoint and why, given that it will eventually be publicly reachable per Chapter 93.
