package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
)

// Everything below this line, down to the "Chain" section, is provided
// for you: your own working versions of Tasks 01, 03, 04, 05, and 06,
// wired together with real digital signatures. Skim it to see how the
// pieces fit, then focus your work on the Chain type at the bottom,
// which is genuinely new.

// ─── crypto (Task 04) ───────────────────────────────────────────────────────

var curve = elliptic.P256()

func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func marshalPublicKey(pub *ecdsa.PublicKey) []byte {
	byteLen := (pub.Curve.Params().BitSize + 7) / 8
	buf := make([]byte, 2*byteLen)
	pub.X.FillBytes(buf[:byteLen])
	pub.Y.FillBytes(buf[byteLen:])
	return buf
}

func unmarshalPublicKey(pubKey []byte) (*ecdsa.PublicKey, error) {
	byteLen := (curve.Params().BitSize + 7) / 8
	if len(pubKey) != 2*byteLen {
		return nil, errors.New("crypto: invalid public key length")
	}
	x := new(big.Int).SetBytes(pubKey[:byteLen])
	y := new(big.Int).SetBytes(pubKey[byteLen:])
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	return &Wallet{PrivateKey: priv, PublicKey: marshalPublicKey(&priv.PublicKey)}
}

func (w *Wallet) Address() string {
	return hex.EncodeToString(Hash(w.PublicKey))
}

func Sign(priv *ecdsa.PrivateKey, data []byte) []byte {
	hash := Hash(data)
	sig, err := ecdsa.SignASN1(rand.Reader, priv, hash)
	if err != nil {
		panic(err)
	}
	return sig
}

func Verify(pubKey []byte, data, signature []byte) bool {
	pub, err := unmarshalPublicKey(pubKey)
	if err != nil {
		return false
	}
	return ecdsa.VerifyASN1(pub, Hash(data), signature)
}

func addressFromPubKey(pubKey []byte) string {
	return hex.EncodeToString(Hash(pubKey))
}

// ─── blocks & proof of work (Tasks 01-03) ──────────────────────────────────

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

func (b *Block) serializeForMining(nonce int64) []byte {
	var buf bytes.Buffer
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(b.Timestamp))
	buf.Write(ts)
	for _, tx := range b.Transactions {
		buf.Write(tx.ID)
	}
	buf.Write(b.PrevBlockHash)
	n := make([]byte, 8)
	binary.BigEndian.PutUint64(n, uint64(nonce))
	buf.Write(n)
	return buf.Bytes()
}

func mineBlock(transactions []*Transaction, prevBlockHash []byte, difficulty int) *Block {
	b := &Block{Timestamp: 1, Transactions: transactions, PrevBlockHash: prevBlockHash}

	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	hashInt := new(big.Int)

	var nonce int64
	var hash []byte
	for nonce = 0; ; nonce++ {
		hash = Hash(b.serializeForMining(nonce))
		hashInt.SetBytes(hash)
		if hashInt.Cmp(target) == -1 {
			break
		}
	}
	b.Nonce = nonce
	b.Hash = hash
	return b
}

// ─── transactions & UTXOs (Task 05), signed with Task 04's crypto ──────────

type TxOutput struct {
	Value      int64
	PubKeyHash []byte
}

type TxInput struct {
	TxID      []byte
	OutIndex  int
	Signature []byte
	PubKey    []byte
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []TxInput
	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{TxID: in.TxID, OutIndex: in.OutIndex})
	}
	return &Transaction{Inputs: inputs, Outputs: tx.Outputs}
}

func (tx *Transaction) serialize() []byte {
	var buf bytes.Buffer
	for _, in := range tx.Inputs {
		buf.Write(in.TxID)
		idx := make([]byte, 8)
		binary.BigEndian.PutUint64(idx, uint64(in.OutIndex))
		buf.Write(idx)
	}
	for _, out := range tx.Outputs {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, uint64(out.Value))
		buf.Write(val)
		buf.Write(out.PubKeyHash)
	}
	return buf.Bytes()
}

func isCoinbase(tx *Transaction) bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0
}

// newCoinbaseTX builds a reward transaction. height must be unique per
// block -- without it, two coinbase transactions paying the same
// address the same reward would serialize to identical bytes and
// collide on the same transaction ID.
func newCoinbaseTX(address string, reward int64, height int64) *Transaction {
	pubKeyHash, _ := hex.DecodeString(address)
	tx := &Transaction{
		Inputs:  []TxInput{{TxID: []byte{}, OutIndex: -1 - int(height)}},
		Outputs: []TxOutput{{Value: reward, PubKeyHash: pubKeyHash}},
	}
	tx.ID = Hash(tx.serialize())
	return tx
}

func findUTXOsWithRefs(transactions []*Transaction, pubKeyHash []byte) []struct {
	TxID  []byte
	Index int
	Value int64
} {
	spent := make(map[string]map[int]bool)
	for _, tx := range transactions {
		if isCoinbase(tx) {
			continue
		}
		for _, in := range tx.Inputs {
			key := fmt.Sprintf("%x", in.TxID)
			if spent[key] == nil {
				spent[key] = make(map[int]bool)
			}
			spent[key][in.OutIndex] = true
		}
	}

	var result []struct {
		TxID  []byte
		Index int
		Value int64
	}
	for _, tx := range transactions {
		key := fmt.Sprintf("%x", tx.ID)
		for i, out := range tx.Outputs {
			if spent[key][i] {
				continue
			}
			if bytes.Equal(out.PubKeyHash, pubKeyHash) {
				result = append(result, struct {
					TxID  []byte
					Index int
					Value int64
				}{tx.ID, i, out.Value})
			}
		}
	}
	return result
}

func balanceOf(transactions []*Transaction, address string) int64 {
	pubKeyHash, err := hex.DecodeString(address)
	if err != nil {
		return 0
	}
	var total int64
	for _, u := range findUTXOsWithRefs(transactions, pubKeyHash) {
		total += u.Value
	}
	return total
}

// newTransaction builds and signs a transaction sending amount from the
// wallet behind priv/pubKey to the given recipient address.
func newTransaction(transactions []*Transaction, priv *ecdsa.PrivateKey, pubKey []byte, to string, amount int64) (*Transaction, error) {
	fromHash := Hash(pubKey)
	toHash, err := hex.DecodeString(to)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient address: %w", err)
	}

	candidates := findUTXOsWithRefs(transactions, fromHash)

	var inputs []TxInput
	var accumulated int64
	for _, c := range candidates {
		if accumulated >= amount {
			break
		}
		inputs = append(inputs, TxInput{TxID: c.TxID, OutIndex: c.Index})
		accumulated += c.Value
	}
	if accumulated < amount {
		return nil, fmt.Errorf("insufficient funds: have %d, need %d", accumulated, amount)
	}

	outputs := []TxOutput{{Value: amount, PubKeyHash: toHash}}
	if accumulated > amount {
		outputs = append(outputs, TxOutput{Value: accumulated - amount, PubKeyHash: fromHash})
	}

	tx := &Transaction{Inputs: inputs, Outputs: outputs}
	tx.ID = Hash(tx.serialize())

	sigData := tx.trimmedCopy().serialize()
	sig := Sign(priv, sigData)
	for i := range tx.Inputs {
		tx.Inputs[i].Signature = sig
		tx.Inputs[i].PubKey = pubKey
	}

	return tx, nil
}

func verifyTransaction(tx *Transaction) bool {
	if isCoinbase(tx) {
		return true
	}
	sigData := tx.trimmedCopy().serialize()
	for _, in := range tx.Inputs {
		if !Verify(in.PubKey, sigData, in.Signature) {
			return false
		}
	}
	return true
}

// ─── mempool (Task 06) ─────────────────────────────────────────────────────

type Mempool struct {
	pending map[string]*Transaction
	claimed map[string]bool
	order   []string
}

func NewMempool() *Mempool {
	return &Mempool{pending: make(map[string]*Transaction), claimed: make(map[string]bool)}
}

func claimKey(txID []byte, outIndex int) string {
	return fmt.Sprintf("%x:%d", txID, outIndex)
}

func (mp *Mempool) Add(tx *Transaction) error {
	if !verifyTransaction(tx) {
		return errors.New("mempool: invalid signature")
	}
	for _, in := range tx.Inputs {
		if mp.claimed[claimKey(in.TxID, in.OutIndex)] {
			return fmt.Errorf("mempool: conflicting spend of %x:%d", in.TxID, in.OutIndex)
		}
	}
	key := fmt.Sprintf("%x", tx.ID)
	mp.pending[key] = tx
	mp.order = append(mp.order, key)
	for _, in := range tx.Inputs {
		mp.claimed[claimKey(in.TxID, in.OutIndex)] = true
	}
	return nil
}

func (mp *Mempool) Pending() []*Transaction {
	result := make([]*Transaction, 0, len(mp.order))
	for _, key := range mp.order {
		result = append(result, mp.pending[key])
	}
	return result
}

func (mp *Mempool) Clear() {
	mp.pending = make(map[string]*Transaction)
	mp.claimed = make(map[string]bool)
	mp.order = nil
}

// ─── Chain: wiring it all together (this is the part you implement) ───────

type Chain struct {
	blocks       []*Block
	mempool      *Mempool
	allConfirmed []*Transaction // every transaction from every mined block, flattened
}

// NewChain creates an empty chain with a genesis block and an empty mempool.
func NewChain() *Chain {
	panic("TODO: implement NewChain")
}

// Send builds a transaction sending amount from the wallet behind priv,
// signs it, and adds it to the mempool. from must be the address
// derived from priv's matching public key. Returns an error if from's
// confirmed balance can't cover amount, or if the mempool rejects it as
// a conflicting spend.
func (c *Chain) Send(priv *ecdsa.PrivateKey, from, to string, amount int64) (*Transaction, error) {
	panic("TODO: implement Send using newTransaction and c.mempool.Add")
}

// Mine takes every pending transaction out of the mempool, adds a
// coinbase reward transaction paying minerAddress, mines a new block
// containing all of it with the given proof-of-work difficulty, and
// appends it to the chain.
func (c *Chain) Mine(minerAddress string, reward int64, difficulty int) *Block {
	panic("TODO: implement Mine using newCoinbaseTX and mineBlock")
}

// Balance returns address's current confirmed balance, based only on
// transactions inside mined blocks.
func (c *Chain) Balance(address string) int64 {
	panic("TODO: implement Balance using balanceOf")
}

// main runs a single-process demo of the whole system end to end.
//
// A CLI whose subcommands persist state across SEPARATE process
// invocations (so "chain-wallet new" today and "chain-wallet send"
// tomorrow share the same chain) needs real disk persistence -- see
// this task's "Stretch goals". This demo shows the same flow within
// one run, which is exactly what the test suite exercises directly
// against Chain.
func main() {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	fs.Parse(os.Args[1:])

	chain := NewChain()

	alice := NewWallet()
	bob := NewWallet()
	fmt.Println("Alice's address:", alice.Address())
	fmt.Println("Bob's address:  ", bob.Address())

	chain.Mine(alice.Address(), 50, 14)
	fmt.Printf("Mined block 1 (reward: 50 coins)\n")
	fmt.Println("Alice's balance:", chain.Balance(alice.Address()))

	tx, err := chain.Send(alice.PrivateKey, alice.Address(), bob.Address(), 20)
	if err != nil {
		fmt.Println("send failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Transaction %x submitted to the mempool\n", tx.ID)

	chain.Mine(alice.Address(), 50, 14)
	fmt.Println("Mined a second block (includes the pending transaction)")

	fmt.Println("Alice's balance:", chain.Balance(alice.Address()))
	fmt.Println("Bob's balance:  ", chain.Balance(bob.Address()))
}
