# Task 03: Proof of Work

## What you will build

A real proof-of-work mining loop: given a block, search for a `Nonce` value that makes the block's hash start with a required number of zero bits, exactly the mechanism that makes Bitcoin-style chains expensive to rewrite.

## Concepts

### The puzzle

Proof of work asks: find a number (the nonce) that, combined with the block's other data, produces a hash meeting some target — in this task, a hash whose first `difficulty` bits are all zero. There's no shortcut: you try nonces one at a time until you get lucky. Trying is cheap per-attempt but the odds of any single guess succeeding are low, so finding a valid nonce takes real, measurable work — and anyone else can instantly verify your answer by hashing it once and checking the target.

```
nonce=0  -> hash 7f3a91...   (doesn't start with enough zero bits)
nonce=1  -> hash c81b22...   (nope)
nonce=2  -> hash 003fd0...   (nope, only 1 leading zero nibble found so far)
   ...
nonce=8452 -> hash 0000e7...  <- meets a 4-hex-digit-zero target! this is our answer.
```

### Difficulty as a target

Rather than counting bits by hand, it's simplest to compare the hash (read as a big number) against a target number: `target = 1 << (256 - difficulty)`. A valid hash, interpreted as a big integer, must be less than `target`. Higher `difficulty` means a smaller target, which means fewer valid hashes exist, which means more attempts are needed on average.

## Interface to implement

```go
type ProofOfWork struct {
	Block      *Block
	Target     *big.Int
}

// NewProofOfWork creates a ProofOfWork for b with the given difficulty
// (number of leading zero bits the hash must have).
func NewProofOfWork(b *Block, difficulty int) *ProofOfWork

// Run searches for a nonce that satisfies the target, starting from 0.
// It returns the winning nonce and the resulting hash.
func (pow *ProofOfWork) Run() (int64, []byte)

// Validate reports whether pow.Block's current Nonce and Hash actually
// satisfy pow.Target -- used to check a solved block without redoing
// the search.
func (pow *ProofOfWork) Validate() bool
```

Extend your `Block` type from Task 02 with a `Nonce int64` field.

## Hints

- Hash together `PrevBlockHash`, `Data`, and the candidate `Nonce` (converted to bytes) on each attempt — the exact same fields your `ComputeHash` already serializes, plus the nonce.
- Use `math/big` to build the target: `target := big.NewInt(1); target.Lsh(target, uint(256-difficulty))`.
- To compare a hash against the target: `hashInt := new(big.Int).SetBytes(hash); hashInt.Cmp(target) == -1` means the hash is below target (valid).
- Keep `difficulty` small for tests (e.g. 12-16 bits) so `go test` finishes in well under a second — a difficulty tuned for a real network would make your test suite take minutes.
- `Validate()` should NOT search for anything — it just recomputes the hash once, using the block's *existing* `Nonce`, and checks it against the target.

## Run the tests

```bash
cd starter/task-03-proof-of-work
go test ./...
```

All tests must pass before moving to Task 04.
