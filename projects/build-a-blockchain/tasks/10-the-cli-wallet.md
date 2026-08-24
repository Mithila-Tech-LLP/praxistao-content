# Task 10: The CLI Wallet

## What you will build

The capstone: wire the blockchain, wallet, transactions, mempool, and mining from every earlier task into one cohesive `Chain` type, then a small CLI on top of it — the same fundamental experience as a real cryptocurrency wallet, entirely backed by code you wrote yourself.

## Concepts

### Composition, not new machinery

Nothing in this task is a genuinely new concept — it's every previous task's pieces, wired together in the order a real user would actually exercise them: generate a wallet, mine a block to get some starting balance (via a coinbase reward), send some of it, mine again to confirm the transaction, check the new balance.

```
                          Chain
                            |
        +-------------------+-------------------+
        |                   |                   |
   Blockchain            Mempool             (miner address)
   (Task 02, 03)         (Task 06)
        |
   []*Transaction  (Task 05)
```

### Why mining is what confirms a transaction

A transaction sitting in the mempool hasn't affected anyone's balance yet — `BalanceOf` (Task 05) only sees transactions that are actually part of a mined block. This is worth confirming for yourself hands-on: send a transaction, check the recipient's balance (still zero), mine a block, check again (now correct).

## Interface to implement

```go
type Chain struct {
	// unexported fields: a *Blockchain (Task 02), a *Mempool (Task 06),
	// and something tracking all mined transactions for BalanceOf/UTXO lookups
}

// NewChain creates an empty chain with a genesis block and an empty mempool.
func NewChain() *Chain

// Send builds a transaction sending amount from from's balance to to,
// signs it (using priv, from Task 04), and adds it to the mempool.
// Returns an error if from's confirmed balance can't cover amount, or if
// the mempool rejects it as a conflicting spend.
func (c *Chain) Send(priv *ecdsa.PrivateKey, from, to string, amount int64) (*Transaction, error)

// Mine takes every pending transaction out of the mempool, adds a
// coinbase reward transaction paying minerAddress, mines a new block
// containing all of it with the given proof-of-work difficulty, and
// appends it to the chain.
func (c *Chain) Mine(minerAddress string, reward int64, difficulty int) *Block

// Balance returns minerAddress's -- or any address's -- current
// confirmed balance, based only on transactions inside mined blocks.
func (c *Chain) Balance(address string) int64
```

Then build a `main()` with subcommands `new`, `mine`, `balance`, and `send`, following the example session in `project.md`.

## Hints

- `Send` needs *some* way to look up `from`'s current UTXOs to build the transaction (Task 05's `NewTransaction`) — gather every transaction from every mined block into one slice and hand it to `FindUTXOs`/`NewTransaction`. It's fine if this is a full rescan for now; a real index (like the full course covers) is a later optimization, not a correctness requirement here.
- Keep the CLI itself simple: Go's standard `flag` package, one `flag.FlagSet` per subcommand, dispatching on `os.Args[1]` — you don't need a CLI framework for four subcommands.
- Since a real session spans multiple separate program runs, persisting wallet keys and chain state between invocations matters for a genuinely usable CLI — but for this task's tests, focus on exercising `Chain` directly within a single test process; persistence across process runs is a great extension exercise once the core logic is solid (see below).
- Write an end-to-end test: create a `Chain`, mine a block rewarding wallet A, `Send` from A to B, mine again, and assert both `Balance(A)` and `Balance(B)` land exactly where the arithmetic predicts (`A`'s original reward minus the amount sent, `B`'s balance equal to the amount received).

## Stretch goals

Once your tests pass, you've built a genuinely working (if simplified) blockchain. A few natural next steps, if you want to keep pushing on it:

- Persist the chain and wallet keys to disk between CLI invocations (revisit Task 01's `Serialize` idea, applied to the whole chain).
- Wire Task 08's P2P networking in: run two `chain-wallet` processes, and have a mined block or a sent transaction on one automatically reach the other.
- Add a `deploy`/`call` subcommand that runs a Task 09 VM program against the chain, gated by the same signature-based authorization your transactions already use.
- For the full depth version of everything in this project — HD wallets, a real P2P gossip network, an embedded database, a complete smart-contract system, a block explorer, and deploying a live multi-node testnet to the cloud — see the companion [**Blockchain in Go**](/course/blockchain-in-go) course.

## Run the tests

```bash
cd starter/task-10-the-cli-wallet
go test ./...
```
