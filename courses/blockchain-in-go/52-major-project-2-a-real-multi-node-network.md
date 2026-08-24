# Chapter 52: Major Project 2 — A Real Multi-Node Network

Every chapter in this volume has been building one piece of machinery at a time — handshakes (Chapter 46), peer discovery (Chapter 47), gossip (Chapter 48), synchronization (Chapter 49), fork resolution (Chapter 50), and an honest look at what still goes wrong (Chapter 51). This chapter assembles every one of those pieces into something you run for real: three independent GoChain processes, three independent wallets, one transaction that starts on a node that didn't mine it, and a block that reaches every node without you manually copying a single byte between them.

## Table of Contents

1. [What "Real" Means for This Project](#1-what-real-means-for-this-project)
2. [Extending main.go for a Scriptable Demo Node](#2-extending-maingo-for-a-scriptable-demo-node)
3. [Sync-on-Connect: Wiring In Chapter 49](#3-sync-on-connect-wiring-in-chapter-49)
4. [A Background Mining Loop](#4-a-background-mining-loop)
5. [Tracing One Transaction's Full Journey](#5-tracing-one-transactions-full-journey)
6. [Major Project: Multi-Node Testnet](#6-major-project-multi-node-testnet)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What "Real" Means for This Project

Up to now, every demo in this volume has been small on purpose: two nodes in Chapter 46, five in Chapter 47's peer-discovery trace — enough to see one mechanism work, not enough to feel like an actual network. This project insists on a specific, deliberately concrete bar for "real": **three separate operating-system processes**, each with its own listen address, its own on-disk chain (via `core.OpenBlockchain`, from Chapter 20, rather than the in-memory `core.NewBlockchain` earlier chapters used for brevity), and its own independently generated wallet — no shared memory, no cheating by having one process peek at another's variables. The only channel any node has to any other is the TCP-based, gob-encoded protocol this entire volume has been building. If a transaction submitted on Node A ends up mined and validated on Node B and Node C, it is because the actual network code did the work, not because we simulated it.

```
   Node A                    Node B                    Node C
   (own process,             (own process,             (own process,
    own data dir,             own data dir,             own data dir,
    own wallet)                own wallet)               own wallet)
       |                         |                         |
       +---------- TCP connections, gob-encoded ------------+
                    Envelope-framed P2P messages
                    (the ONLY way these processes
                     ever exchange information)
```

---

## 2. Extending main.go for a Scriptable Demo Node

Chapter 46 introduced `cmd/gochain-node/main.go`; Chapter 47 added `-seeds` and background peer exchange. This project adds the last few flags needed to script an entire demo without touching any code once it starts: a persisted wallet, an optional up-front "premine" so a node has coins to send in the first place, a scheduled send, and a background mining loop.

```go
package main

import (
	"encoding/hex"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/you/gochain/core"
	"github.com/you/gochain/network"
	"github.com/you/gochain/wallet"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3000", "this node's listen address")
	dataDir := flag.String("datadir", "./data", "directory for this node's on-disk chain")
	walletFile := flag.String("walletfile", "", "path to this node's wallet file (default: <datadir>/wallet.dat)")
	seeds := flag.String("seeds", "", "comma-separated seed addresses to bootstrap from")
	premine := flag.Int("premine", 0, "mine this many blocks immediately at startup, rewarding this node's own wallet")
	mineEvery := flag.Duration("mine-every", 0, "if > 0, mine a block on this interval whenever the mempool is non-empty")
	sendTo := flag.String("send-to", "", "if set, send send-amount gochips to this address after send-after")
	sendAmount := flag.Int("send-amount", 0, "amount to send to send-to")
	sendAfter := flag.Duration("send-after", 5*time.Second, "how long to wait before sending, to give the network time to form")
	statusEvery := flag.Duration("status-every", 5*time.Second, "how often to log this node's height, tip hash, and balance")
	flag.Parse()

	if *walletFile == "" {
		*walletFile = *dataDir + "/wallet.dat"
	}

	// Load an existing wallet if this node has run before; otherwise
	// generate a fresh one and save it, exactly like the CLI wallet from
	// Chapter 36. A wallet is nothing more than a key pair and the address
	// derived from it, so a brand-new gochip balance always starts at zero.
	w, err := wallet.Load(*walletFile)
	if err != nil {
		w = wallet.New()
		if err := w.Save(*walletFile); err != nil {
			log.Fatalf("could not save new wallet: %v", err)
		}
		log.Printf("[%s] generated new wallet, address: %s", *addr, w.Address())
	} else {
		log.Printf("[%s] loaded existing wallet, address: %s", *addr, w.Address())
	}

	// The chain rewards THIS node's own wallet address for every block it
	// mines locally -- exactly the pattern the very first code sample in
	// this course's table of contents used.
	chain, err := core.OpenBlockchain(*dataDir, w.Address())
	if err != nil {
		log.Fatalf("could not open chain: %v", err)
	}
	defer chain.Close()

	if *premine > 0 {
		for i := 0; i < *premine; i++ {
			block := chain.MineBlock()
			log.Printf("[%s] premined block %d (%s)", *addr, block.Height, block.HashHex())
		}
	}

	mempool := core.NewMempool()
	node := network.NewNode(*addr, chain, mempool)
	if err := node.Listen(); err != nil {
		log.Fatal(err)
	}
	node.StartPeerExchange(10 * time.Second)

	if *seeds != "" {
		node.Bootstrap(strings.Split(*seeds, ","))
	}

	if *mineEvery > 0 {
		startMiningLoop(node, *addr, *mineEvery)
	}
	if *sendTo != "" && *sendAmount > 0 {
		go func() {
			time.Sleep(*sendAfter)
			sendCoins(node, w, *addr, *sendTo, *sendAmount)
		}()
	}
	startStatusLoop(chain, w, *addr, *statusEvery)

	select {} // block forever
}
```

`wallet.Load`/`wallet.Save` are the small persistence helpers Chapter 36's CLI wallet already relies on; using them here means restarting a node (say, after a crash) reuses the same address and the same coin ownership rather than silently starting over with a brand-new, empty wallet. `premine` exists purely so this demo doesn't need an external faucet: mining a block or two locally, before ever talking to a peer, is enough for Node A to have real gochips to send.

---

## 3. Sync-on-Connect: Wiring In Chapter 49

Chapter 49 built `SyncWithPeer` but deliberately left the *policy* of when to call it to whoever owns a `Node` — and named "sync-on-connect" as this exact project's job. `Bootstrap` (Chapter 47) already knows which seeds it successfully dialed; this project's `main.go` extends that success path with one call:

```go
import "log"

// bootstrapAndSync wraps Chapter 47's Bootstrap with sync-on-connect: the
// moment a seed answers, immediately ask it to catch us up, rather than
// waiting for gossip to eventually mention blocks we may have already
// missed before we ever connected.
func bootstrapAndSync(n *network.Node, seeds []string, myAddr string) {
	for _, seed := range seeds {
		if seed == myAddr {
			continue
		}
		if err := n.Dial(seed); err != nil {
			log.Printf("[%s] could not reach seed %s: %v", myAddr, seed, err)
			continue
		}
		log.Printf("[%s] connected to seed %s, requesting sync", myAddr, seed)
		if err := n.SyncWithPeer(seed); err != nil {
			log.Printf("[%s] sync with %s failed (may just mean nothing to sync yet): %v", myAddr, seed, err)
		}
	}
}
```

This replaces the plain `node.Bootstrap(strings.Split(*seeds, ","))` call in Section 2's `main` with `bootstrapAndSync(node, strings.Split(*seeds, ","), *addr)`. Notice the log wording on failure: an error from `SyncWithPeer` immediately after connecting to a brand-new seed is extremely common and not alarming — Chapter 49, Section 8 already established that "the peer has nothing new" is a perfectly normal, non-error outcome, and a truly empty `InvPayload` returns `nil` from `SyncWithPeer` rather than an error at all. A logged error here more often means a timeout on a slow-starting peer than anything actually wrong.

---

## 4. A Background Mining Loop

Node B, in this project, is the one designated to mine — meaning it periodically checks whether anything is waiting in its mempool and, if so, mines it into a real block, exactly the way Chapter 48's `HandleBlock` already anticipates ("either mined locally or received from a peer"):

```go
import "time"

// startMiningLoop mines a new block on a fixed interval whenever the
// mempool has at least one pending transaction, then feeds the freshly
// mined block through HandleBlock -- the exact same path a block received
// from a peer would take, so gossip (Chapter 48) picks it up and forwards
// it to every connected peer without any special-casing for "blocks we
// mined ourselves" versus "blocks a peer sent us."
func startMiningLoop(n *network.Node, myAddr string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if len(n.Mempool.GetAll()) == 0 {
				continue // nothing waiting -- no point mining an empty block
			}
			block := n.Chain.MineBlock()
			log.Printf("[%s] mined block %d (%s) with %d transaction(s)",
				myAddr, block.Height, block.HashHex(), len(block.Transactions))
			n.HandleBlock(block)
		}
	}()
}
```

Routing a locally mined block through `n.HandleBlock` rather than calling `n.Broadcast` directly is deliberate, and mirrors exactly how `sendCoins` below routes a locally created transaction through `n.HandleTx` rather than broadcasting it by hand: `HandleBlock` already does the seen-set bookkeeping (Chapter 48) that prevents this exact block from being needlessly reprocessed if it somehow loops back to us, and it is the one place fork-resolution logic (Chapter 50) lives, in case this node was — however briefly — behind on a competing chain when it started mining.

```go
// sendCoins builds and signs a transaction (core.Blockchain.Send, from
// Volume 5) and feeds it through HandleTx, the exact same entry point a
// transaction arriving from a peer would use. This is what actually gets
// a transaction created on THIS node gossiped out to every peer.
func sendCoins(n *network.Node, w *wallet.Wallet, myAddr, to string, amount int) {
	tx, err := n.Chain.Send(w, to, amount)
	if err != nil {
		log.Printf("[%s] could not send %d gochips to %s: %v", myAddr, amount, to, err)
		return
	}
	log.Printf("[%s] submitted transaction %s: sending %d gochips to %s", myAddr, tx.IDHex(), amount, to)
	n.HandleTx(tx)
}
```

---

## 5. Tracing One Transaction's Full Journey

Before running the real script, it helps to have the entire expected sequence in one picture — every step below is code this course has already built, in earlier chapters, being exercised together for the first time:

```
  Node A                    Node B (miner)              Node C
  ------                    --------------               ------
  sendCoins() (Sec. 4)
    chain.Send(...)
    HandleTx(tx)   ---------> gossiped via
       |                      Broadcast (Ch. 46)
       |                          |
       |                          v
       |                    A's peer connection
       |                    delivers MsgTx
       |                          |
       |                          v
       |                    HandleTx(tx) (Ch. 48):
       |                    dedup, Mempool.Add,
       |                    re-broadcast -----------------> HandleTx(tx):
       |                                                     dedup, Mempool.Add
       |                          |
       |                    mining loop (Sec. 4)
       |                    fires: MineBlock(),
       |                    HandleBlock(block)
       |                          |
       |                          v
       |                    Broadcast(MsgBlock)
       |                          |
       v                          v                            v
  HandleBlock(block):       (already applied         HandleBlock(block):
  ValidateBlock, AddBlock    locally, block             ValidateBlock, AddBlock
  (Ch. 19, 25, 48)           already reflects            (Ch. 19, 25, 48)
                              its own new balance)
       |                                                       |
       v                                                       v
  A's balance for                                        C's balance for
  the recipient address                                  the recipient address
  now matches B's                                         now matches B's
```

Node A never talks to Node C directly in this trace (only A-B and B-C connections exist, in the shell script's topology below) — yet the transaction and the block both reach C anyway, purely through gossip forwarding, exactly the mechanism Chapter 48 proved converges the whole network from any single origin point.

---

## 6. Major Project: Multi-Node Testnet

This is the whole thing, assembled and runnable: a shell script that starts three GoChain nodes, submits a transaction, waits for it to be mined and propagated, and then automatically checks that all three nodes agree on the resulting chain.

### The Script

```bash
#!/usr/bin/env bash
# run-testnet.sh -- starts a 3-node GoChain testnet, submits and mines one
# transaction, and asserts that every node converges on the same chain.
set -uo pipefail

ROOT="./testnet-run"
rm -rf "$ROOT"
mkdir -p "$ROOT/a" "$ROOT/b" "$ROOT/c"

PIDS=()
cleanup() {
	echo "--- shutting down all nodes ---"
	for pid in "${PIDS[@]}"; do
		kill "$pid" 2>/dev/null
	done
}
trap cleanup EXIT

# Step 1: start Node B first, alone, with no peers yet. It generates its
# own wallet on first run -- we need its address before Node A can be told
# to send coins to it.
go run ./cmd/gochain-node \
	-addr 127.0.0.1:3001 -datadir "$ROOT/b" \
	-mine-every 3s -status-every 4s \
	> "$ROOT/b/log.txt" 2>&1 &
PIDS+=($!)

echo "waiting for Node B's wallet to be generated..."
until grep -q "wallet, address:" "$ROOT/b/log.txt" 2>/dev/null; do sleep 0.5; done
B_ADDR=$(grep "wallet, address:" "$ROOT/b/log.txt" | head -n1 | sed -E 's/.*address: //')
echo "Node B's address: $B_ADDR"

# Step 2: start Node C, bootstrapping from B.
go run ./cmd/gochain-node \
	-addr 127.0.0.1:3002 -datadir "$ROOT/c" \
	-seeds 127.0.0.1:3001 -status-every 4s \
	> "$ROOT/c/log.txt" 2>&1 &
PIDS+=($!)

# Step 3: start Node A, bootstrapping from B, premining one block so it
# has coins to send, and scheduled to send 10 gochips to Node B's address
# five seconds after startup -- enough time for the handshake and initial
# sync to complete first.
go run ./cmd/gochain-node \
	-addr 127.0.0.1:3000 -datadir "$ROOT/a" \
	-seeds 127.0.0.1:3001 -premine 1 \
	-send-to "$B_ADDR" -send-amount 10 -send-after 5s \
	-status-every 4s \
	> "$ROOT/a/log.txt" 2>&1 &
PIDS+=($!)

echo "network running -- letting it settle for 25 seconds..."
sleep 25

# Step 4: convergence check. Every node has been logging its own status
# line (Section 2's startStatusLoop) on a timer; the LAST such line in
# each log file reflects each node's most up-to-date view of the chain.
TIP_A=$(grep "status:" "$ROOT/a/log.txt" | tail -n1 | grep -oE 'tip=[0-9a-f]+' )
TIP_B=$(grep "status:" "$ROOT/b/log.txt" | tail -n1 | grep -oE 'tip=[0-9a-f]+' )
TIP_C=$(grep "status:" "$ROOT/c/log.txt" | tail -n1 | grep -oE 'tip=[0-9a-f]+' )

echo "Node A tip: $TIP_A"
echo "Node B tip: $TIP_B"
echo "Node C tip: $TIP_C"

if [[ -n "$TIP_A" && "$TIP_A" == "$TIP_B" && "$TIP_B" == "$TIP_C" ]]; then
	echo "CONVERGED: all three nodes agree on the chain tip."
	exit 0
else
	echo "NOT CONVERGED: nodes disagree -- see $ROOT/{a,b,c}/log.txt"
	exit 1
fi
```

The one function this script assumes but Section 2 didn't yet show is `startStatusLoop`, a small companion to `startMiningLoop` that prints exactly the line the script greps for:

```go
// startStatusLoop periodically logs this node's height, short tip hash,
// mempool size, and this node's own wallet balance -- the single line
// the demo script's convergence check greps for.
func startStatusLoop(chain *core.Blockchain, w *wallet.Wallet, myAddr string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			tip, err := chain.Tip()
			if err != nil {
				continue
			}
			shortHash := hex.EncodeToString(tip.Hash)[:8]
			balance := chain.BalanceOf(w.Address())
			log.Printf("[%s] status: height=%d tip=%s mempool=%d balance=%d",
				myAddr, chain.Height(), shortHash, len(chain.Mempool.GetAll()), balance)
		}
	}()
}
```

### Running It

```bash
chmod +x run-testnet.sh
./run-testnet.sh
```

### Expected Output

Node B's log (`testnet-run/b/log.txt`) — the miner:

```
2026/07/31 12:00:00 [127.0.0.1:3001] generated new wallet, address: 1Bv1Zk...B7q
2026/07/31 12:00:00 [127.0.0.1:3001] listening for peers
2026/07/31 12:00:05 [127.0.0.1:3001] status: height=0 tip=00000000 mempool=0 balance=0
2026/07/31 12:00:09 [127.0.0.1:3001] peer 127.0.0.1:3002 is caught up (height 0)
2026/07/31 12:00:10 [127.0.0.1:3001] peer 127.0.0.1:3000 is caught up (height 1)
2026/07/31 12:00:12 [127.0.0.1:3001] mined block 1 (7f3a9c2e) with 1 transaction(s)
2026/07/31 12:00:13 [127.0.0.1:3001] status: height=1 tip=7f3a9c2e mempool=0 balance=10
```

Node A's log (`testnet-run/a/log.txt`) — the sender:

```
2026/07/31 12:00:00 [127.0.0.1:3000] generated new wallet, address: 1FgQ8m...X2d
2026/07/31 12:00:00 [127.0.0.1:3000] premined block 1 (a1b2c3d4)
2026/07/31 12:00:00 [127.0.0.1:3000] listening for peers
2026/07/31 12:00:01 [127.0.0.1:3000] connected to seed 127.0.0.1:3001, requesting sync
2026/07/31 12:00:05 [127.0.0.1:3000] submitted transaction 9e8d7c6b: sending 10 gochips to 1Bv1Zk...B7q
2026/07/31 12:00:13 [127.0.0.1:3000] status: height=1 tip=7f3a9c2e mempool=0 balance=90
```

Node C's log (`testnet-run/c/log.txt`) — never talked to directly by A, yet fully caught up:

```
2026/07/31 12:00:00 [127.0.0.1:3002] generated new wallet, address: 1Ck4Wn...T9z
2026/07/31 12:00:00 [127.0.0.1:3002] listening for peers
2026/07/31 12:00:00 [127.0.0.1:3002] connected to seed 127.0.0.1:3001, requesting sync
2026/07/31 12:00:13 [127.0.0.1:3002] status: height=1 tip=7f3a9c2e mempool=0 balance=0
```

And the script's own final output:

```
Node A tip: tip=7f3a9c2e
Node B tip: tip=7f3a9c2e
Node C tip: tip=7f3a9c2e
CONVERGED: all three nodes agree on the chain tip.
```

Three things are worth confirming for yourself in a real run, not just reading here: Node A's `premine` block (height 1, hash `a1b2c3d4`) never appears in any node's *final* status line — it was A's own private chain state before it ever connected to anyone, and the moment A's `SyncWithPeer` (or later fork resolution) reconciled with the network's real height-1 block (`7f3a9c2e`, mined by B), `ReplaceChain` from Chapter 50 is exactly the mechanism that resolves which of the two competing height-1 blocks the whole network actually keeps. Second, Node A's final balance (90, assuming a starting coinbase reward of 100) reflects that it *spent* 10 gochips, even though the block that confirmed that spend was mined by a different process entirely. Third, Node C's balance (0) is correct and expected — C never received any gochips in this script; what converged for C is the *chain*, not any particular balance.

---

## Summary

- This project's bar for "real" is strict on purpose: three separate OS processes, three separate on-disk chains, three separate wallets, communicating only through the TCP/gob protocol built across this entire volume.
- `main.go` gained persisted wallets (`wallet.Load`/`wallet.Save`), an up-front `-premine` flag to seed a node with spendable coins, and scheduled, flag-driven actions (`-send-to`/`-send-amount`/`-send-after`) so an entire demo scenario can be scripted without touching Go code at runtime.
- `bootstrapAndSync` wires Chapter 49's `SyncWithPeer` into the moment a seed connection succeeds — the concrete "sync-on-connect" policy Chapter 49, Section 9 explicitly deferred to this chapter.
- A background mining loop calls `MineBlock` whenever the mempool is non-empty, then routes the result through `HandleBlock` — the same entry point a block received from a peer uses, so gossip and fork resolution apply uniformly regardless of a block's origin.
- `sendCoins` builds a transaction with `Chain.Send` and routes it through `HandleTx` for the same reason — local and remote transactions flow through identical logic.
- The major project's shell script starts the network, waits for a scheduled cross-node transaction to be mined and propagate, and then automatically checks a "status" line each node already logs, asserting all three converge on an identical tip hash.
- The expected run demonstrates every mechanism this volume built working together: peer discovery finding Node C without it ever dialing Node A directly, sync-on-connect catching every node up, gossip carrying both the transaction and the resulting block, and (implicitly) fork resolution reconciling Node A's own premined block against the network's real history.

---

## Exercises

### Easy

1. **Run `run-testnet.sh` yourself** at least twice. Compare the two runs' tip hashes (they should differ between runs, since each run generates fresh random wallets and thus a different genesis-adjacent transaction) but confirm all three nodes agree *within* each individual run.

2. **Change `-send-amount` to more than Node A's premined balance** (for instance, 1000 when the coinbase reward is 100) and observe what happens. Which layer of validation from earlier chapters catches this, and where in the logs does the rejection show up?

3. **Add a fourth node, Node D**, bootstrapping only from Node C (not from B directly), and update the script's convergence check to include it. Confirm D still converges despite never connecting to B or A directly.

### Medium

4. **Modify the mining loop (Section 4) to mine even when the mempool is empty**, producing regular empty blocks (a real, if wasteful, thing some blockchains do to keep block times steady). Rerun the testnet and observe how this changes the final height and the timing of convergence.

5. **Introduce an artificial network delay**: before starting Node C, insert a `sleep 15` in the script, so it joins the network well after the transaction has already been mined. Confirm Node C still ends up converged, entirely through Chapter 49's synchronization rather than gossip (since gossip alone cannot deliver something that happened before a node connected).

6. **Add a `-crash-after duration` flag** to `main.go` that calls `os.Exit(1)` after the given duration, and use it to simulate Node B (the miner) crashing moments after mining the block but before every peer has necessarily received it via gossip. Confirm the remaining two nodes still converge (assuming at least one of them received the block before the crash) and explain what would happen if the crash occurred before *any* peer received it.

### Hard

7. **Rewrite the convergence check as a Go test** instead of a shell script, using `os/exec` to launch the three node processes, polling each one's log file (or, better, adding a minimal Unix-domain-socket or loopback-TCP status endpoint to `main.go` that a test can query directly instead of parsing logs) until all three report matching tip hashes or a timeout expires. Discuss why parsing log output, as this chapter's script does, is fragile compared to a real status endpoint, and what changes once Chapter 70's JSON-RPC API exists.

8. **Deliberately create a fork in the script**: start Node A and Node B connected to each other but NOT to Node C, have A and B each mine a competing block containing different transactions at roughly the same time (you'll need to briefly pause one node's mining relative to the other to make the timing work), then connect C to both and confirm — via logs — that `ReplaceChain` (Chapter 50) picks a winner and all three nodes converge on it, with the losing block's transaction returned to some node's mempool as pending.

9. **Extend the testnet to five nodes with a deliberately sparse topology** (each node bootstraps from only one other, forming a chain: A to B, B to C, C to D, D to E, with no other connections at all), submit a transaction on Node A, and measure (by timestamping status lines) how many seconds it takes to reach Node E purely through gossip and peer-exchange-driven connection growth from Chapter 47 and Chapter 48. Compare this propagation time against a fully-meshed topology of the same five nodes and discuss the trade-off in your own words.
