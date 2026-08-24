# Chapter 38: HD Wallets — BIP-32, BIP-39, and BIP-44 Explained

The wallet you built in Chapter 36 holds exactly one key pair. Lose the private key, lose the coins — and if you want a second address for a second purpose, you have to generate a whole new, separate key pair and back it up separately too. Real wallets solve this with **hierarchical deterministic (HD) wallets**: an entire, unlimited tree of key pairs, all regenerable from a single starting point you only ever have to write down once. This chapter builds the conceptual picture — with no code yet — so that Chapters 39 and 40's implementations feel like an obvious continuation rather than a pile of new cryptography.

## Table of Contents

1. [The Problem: One Key Pair Isn't Enough](#1-the-problem-one-key-pair-isnt-enough)
2. [The Family Tree Idea](#2-the-family-tree-idea)
3. [BIP-39: Turning Randomness Into Words You Can Write Down](#3-bip-39-turning-randomness-into-words-you-can-write-down)
4. [BIP-32: Deriving a Tree of Keys From One Seed](#4-bip-32-deriving-a-tree-of-keys-from-one-seed)
5. [Hardened vs. Normal Derivation](#5-hardened-vs-normal-derivation)
6. [BIP-44: A Standard Map So Every Wallet Agrees](#6-bip-44-a-standard-map-so-every-wallet-agrees)
7. [Putting the Three BIPs Together](#7-putting-the-three-bips-together)
8. [Where This Lands in gochain/wallet](#8-where-this-lands-in-gochainwallet)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Problem: One Key Pair Isn't Enough

Think about how you actually use money day to day. You might want one account for savings, one for a joint household budget, and one you hand out publicly on a business card without worrying it links back to your savings. A single key pair forces you to either reuse the same address everywhere (bad for privacy — anyone can see every payment you've ever received) or generate a pile of unrelated key pairs, each of which needs its own separate backup. Lose one backup, lose that slice of your funds forever, with no relationship to the others that might help you recover it.

Real wallet software — MetaMask, Ledger, Trezor, every mobile Bitcoin or Ethereum wallet — solves this with a completely different design: generate **one** piece of secret randomness, called a **seed**, and derive every key pair the wallet will ever use from it, following a fixed, repeatable recipe. Back up the seed once (usually written down as a list of words), and you can regenerate every single derived address again from scratch, on any device, at any time, even years later.

This is what **hierarchical deterministic (HD)** means, word by word:

- **Hierarchical** — the keys are organized in a tree, with parent keys able to produce child keys, which can themselves produce grandchild keys, and so on.
- **Deterministic** — given the same seed and the same position in the tree, you always get the exact same key pair back, every time, on every device. Nothing about derivation involves fresh randomness.

Three standards work together to make this real, and this chapter takes them one at a time: **BIP-39** (turning randomness into a human-writable phrase), **BIP-32** (deriving the actual tree of keys from that phrase), and **BIP-44** (agreeing on *where* in that tree different kinds of keys should live). "BIP" stands for **Bitcoin Improvement Proposal** — a design document the Bitcoin community uses to standardize behavior — but all three of these particular BIPs are followed far beyond Bitcoin itself, including by Ethereum wallets and, in this volume, by GoChain.

---

## 2. The Family Tree Idea

Before any acronyms, picture an actual family tree. One ancestor sits at the root. Every descendant can be reached by following a specific path down from that root — child 3, then their child 1, then *their* child 0 — and that path uniquely and permanently identifies one specific person. Nobody needs to separately remember who each descendant is; if you know the root ancestor and the rule for how children are produced, you can reconstruct the entire tree, on demand, indefinitely.

An HD wallet's key tree works exactly the same way:

```
                         master seed
                        (the one thing
                         you back up)
                              |
                       derive master key
                              |
              +---------------+---------------+
              |               |               |
        account 0         account 1         account 2
              |               |               |
      +-------+-------+       |               |
      |               |       |               |
  address 0       address 1  address 0    address 0
  (child key)     (child key) (child key) (child key)
      |               |
   Address:        Address:
   1A2b3C...        1X9y8Z...
```

Every box in that tree — the master key, each account key, each address key — is a full, valid `crypto.KeyPair`, capable of signing transactions on its own. None of them are stored anywhere except the master seed at the very top. Every other key is *recomputed on the fly*, on demand, by walking down from the seed using a fixed mathematical rule. If your hard drive with your address-0 private key dies, that private key was never actually "the backup" in the first place — the seed at the top was. Regenerate it from the seed and you get the exact same key pair back, byte for byte.

This single idea — one backup, an entire regenerable tree — is what the rest of this chapter (and Chapters 39-40's code) is building toward.

---

## 3. BIP-39: Turning Randomness Into Words You Can Write Down

The master seed itself is nothing more than a large block of random bytes — typically 128 to 256 bits, generated the same way you'd generate any cryptographic key. But random bytes are a terrible thing to ask a human being to write down and later retype correctly: `4f 3a 91 c2 88 0d ...` invites transcription errors, and a single wrong hex digit silently produces a *different, completely unrelated* seed with no warning.

**BIP-39** fixes this by mapping those random bytes onto a fixed, publicly known list of exactly 2048 English words (other languages have their own official lists), chosen specifically because they are common, unambiguous to spell, and hard to confuse with each other. The random bytes get sliced into small chunks, and each chunk becomes an index into that word list. The result — your **mnemonic** or **seed phrase** — looks like this:

```
witch  collapse  practice  feed  shame  open
maple  north     rescue    fee   coyote crumble
```

Twelve (or twenty-four) ordinary words, easy to write on paper, easy to read aloud over the phone to yourself while double-checking, and — crucially — self-checking: BIP-39 also folds a small **checksum** into the last word, so a single mistyped word is caught immediately as invalid, instead of quietly generating a completely different, wrong wallet that silently controls no funds. Chapter 39 implements the exact algorithm behind this: how many bits become how many words, and how the checksum is computed.

A **passphrase** is an optional extra word or sentence you can add on top of the mnemonic (BIP-39 calls this the "25th word" convention, even though it isn't literally added as a 2048-list word). Two different passphrases combined with the *same* twelve words produce two entirely different, unrelated seeds — a feature some wallets offer as a form of plausible deniability or an extra secret layer, which we will also support in `SeedFromMnemonic()`.

---

## 4. BIP-32: Deriving a Tree of Keys From One Seed

Having a seed and a human-readable phrase for it solves the backup problem, but doesn't yet explain how a *tree* of keys comes out of it. That's **BIP-32**'s job: hierarchical deterministic key derivation.

The process starts by combining the seed with a fixed derivation formula (Chapter 40 implements the actual math) to produce a **master key** — which is really two things bundled together: a master *private key*, and a **chain code**, an extra 32 bytes of entropy whose only job is to make child derivation unpredictable to anyone who doesn't have it. Think of the chain code as a special ingredient mixed into every derivation step; without it, watching one child key's derivation might let an attacker guess the rule and derive siblings on their own.

From the master key, BIP-32 defines a repeatable step: **given a parent key, a parent chain code, and an index number, produce a child key and a new chain code.** Apply that step once, and you get one child. Apply it again to that child, and you get a grandchild. The whole tree in Section 2's diagram is just this one step, applied repeatedly along different paths.

```
   master key + chain code
            |
            |  derive(index = 0)
            v
   child(0) key + chain code
            |
            |  derive(index = 5)
            v
   grandchild(0/5) key + chain code
```

Two properties make this genuinely useful, not just a cute recursive trick:

- **One-way, per level.** Knowing a child key does not let you work backward to recover its parent or siblings — derivation only runs forward, down the tree.
- **Fully deterministic.** Running `derive(masterKey, chainCode, index=0)` today and running it again next year, on a completely different machine, produces the exact same child key and chain code, provided you start from the same master key. This is what lets "restore my wallet" work at all: no state needs to be saved except the original seed.

---

## 5. Hardened vs. Normal Derivation

BIP-32 actually defines two flavors of the derivation step, and the difference matters for security, not just organization.

**Normal derivation** computes a child key using the parent's *public* key as part of its input. This has a useful side effect: someone holding only a parent public key (not the private key) can derive all the corresponding child *public* keys — handy for, say, a web store that wants to generate a fresh receiving address per customer without holding any private key on the server at all.

**Hardened derivation** instead uses the parent's *private* key as input, which means it is impossible to derive a hardened child from a public key alone — you must have the private key. The trade-off buys real security: without hardened derivation, if an attacker ever obtained *one* child private key plus the shared parent public key, they could mathematically reconstruct the parent private key itself, exposing every other child in the entire tree. Hardened derivation closes that specific attack for the levels of the tree where it matters most.

By convention, hardened indices are written with an apostrophe — `0'`, `44'` — and are internally just the index number plus 2³¹ (`0x80000000`), which is how software tells the two flavors apart using a single 32-bit index field.

```
Index 0            -> normal derivation  (public-key math only)
Index 0' (=2^31)   -> hardened derivation (requires the private key)
```

BIP-44, covered next, specifies exactly which levels of a wallet's tree should use hardened derivation and which can stay normal.

---

## 6. BIP-44: A Standard Map So Every Wallet Agrees

BIP-32 tells you *how* to derive a child key from a parent. It says nothing about *where in the tree* a given address should live — and without an agreed convention, every wallet vendor would invent its own layout, and importing your seed phrase into a different wallet app would produce a completely different, useless set of addresses.

**BIP-44** fixes this by defining a standard **derivation path**: a fixed sequence of five levels every compliant wallet uses, in the same order, for the same purpose:

```
m / purpose' / coin_type' / account' / change / address_index

m            — the master key (the root of the tree)
purpose'     — always 44' for BIP-44-style wallets (hardened)
coin_type'   — which cryptocurrency (0' = Bitcoin, 60' = Ethereum, ...) (hardened)
account'     — which "account" within the wallet, e.g. 0', 1', 2'... (hardened)
change       — 0 for a normal receiving address, 1 for internal "change" addresses
address_index — 0, 1, 2, ... — which specific address within the account
```

```
                              m
                              |
                            44'                (purpose: BIP-44)
                              |
                          coin_type'            (which chain)
                              |
                          account'              (which "wallet" within the app)
                          /        \
                    change=0      change=1      (receiving vs. change addresses)
                    /     \
              index=0   index=1  ...            (individual addresses)
```

Notice the first three levels are hardened (the apostrophe), while `change` and `address_index` are normal. This exactly follows Section 5's reasoning: `purpose`, `coin_type`, and `account` sit high in the tree and protect the most value if compromised, so they use hardened derivation; `change` and `address_index` are derived far more often (potentially generating a fresh address per transaction) and benefit from normal derivation's ability to be watched by public-key-only software, like a receipt-printing web server.

The practical payoff: a recovery phrase generated by one BIP-44-compliant wallet, typed into a *completely different* BIP-44-compliant wallet, reconstructs the exact same addresses in the exact same order — because both wallets are walking the exact same fixed path down the exact same tree.

---

## 7. Putting the Three BIPs Together

Here is the full pipeline, start to finish, with each BIP's job labeled:

```
random bytes (128-256 bits)
        |
        |  BIP-39: map bytes onto the 2048-word list, add a checksum word
        v
mnemonic ("witch collapse practice feed shame open maple ...")
        |
        |  BIP-39: PBKDF2 stretch (mnemonic + optional passphrase) -> 64-byte seed
        v
seed (64 bytes)
        |
        |  BIP-32: HMAC-derive master private key + master chain code
        v
master key + chain code
        |
        |  BIP-32 derivation step, walked along a
        |  BIP-44 path: m/44'/coin_type'/account'/change/index
        v
one specific address's key pair
        |
        |  Volume 2, Chapter 14: hash public key, add version byte + checksum, Base58-encode
        v
a real GoChain address, ready to receive gochips
```

Every arrow in that pipeline is implemented as real Go code before this volume ends: BIP-39's word mapping and seed stretching in Chapter 39, BIP-32's key tree and BIP-44's path in Chapter 40.

---

## 8. Where This Lands in gochain/wallet

To ground all of this in the actual code you are about to write, here is the exact shape Chapters 39 and 40 fill in, so you know where each concept from this chapter is headed:

```go
package wallet

// NewMnemonic generates fresh randomness and encodes it as a BIP-39
// seed phrase. Chapter 39.
func NewMnemonic() (string, error)

// SeedFromMnemonic turns a mnemonic (plus an optional passphrase) back
// into the 64-byte seed used for derivation. Chapter 39.
func SeedFromMnemonic(mnemonic, passphrase string) []byte

// HDWallet holds the one thing that must ever be backed up (the seed)
// plus the derived master key used to walk the BIP-32 tree. Chapter 40.
type HDWallet struct {
	Seed      []byte
	MasterKey []byte
}

// NewHDWallet derives the master key from a seed. Chapter 40.
func NewHDWallet(seed []byte) *HDWallet

// DeriveWallet walks a fixed BIP-44 path down to a single address,
// returning a normal *Wallet indistinguishable from one made with
// wallet.New(). Chapter 40.
func (hd *HDWallet) DeriveWallet(index uint32) (*Wallet, error)
```

Notice that `DeriveWallet` returns a plain `*Wallet` — the exact same type Chapter 36's simple, single-key-pair wallet used. This is deliberate: everything downstream of a wallet (signing transactions, checking balances, the CLI from Chapter 36) keeps working unchanged whether that `*Wallet` came from `wallet.New()` or from walking a whole derivation tree. HD wallets change how a key pair is *produced*, not what a key pair *is*.

---

## Summary

- A single key pair forces an all-or-nothing backup and hurts privacy through address reuse; **hierarchical deterministic (HD) wallets** solve both by deriving an entire tree of key pairs from one seed.
- **Hierarchical** means the keys form a tree; **deterministic** means the same seed and the same tree position always reproduce the exact same key pair.
- **BIP-39** turns random bytes into a human-writable, self-checking mnemonic (seed phrase), and stretches it (with an optional passphrase) into a 64-byte seed.
- **BIP-32** defines the actual derivation step — parent key + chain code + index → child key + new chain code — that builds the tree, in **hardened** (private-key-based, more secure) or **normal** (public-key-based, more flexible) flavors.
- **BIP-44** standardizes *where* in that tree different kinds of keys live via the path `m/purpose'/coin_type'/account'/change/address_index`, so different wallet software agrees on the same layout.
- The **chain code** is the extra ingredient in BIP-32 derivation that keeps sibling keys unpredictable to anyone without it.
- `gochain/wallet`'s `DeriveWallet()` will return a plain `*Wallet` — HD derivation changes how keys are produced, not what they are once produced.
- Nothing in this chapter was implemented yet — Chapters 39 and 40 turn every arrow in Section 7's pipeline into real, tested Go code.

---

## Exercises

### Easy

1. **Draw your own three-level family tree** (on paper or in a text file) representing a master key deriving two accounts, each of which derives two addresses. Label every node with a plausible index number, and mark which arrows in your drawing would be hardened versus normal, based on Section 6's convention.

2. **In your own words**, explain why writing down twelve BIP-39 words is safer against transcription error than writing down the equivalent raw hex bytes of a seed. Give one concrete example of a mistake a human could make with hex that BIP-39's checksum word would catch.

3. **List the five levels of a BIP-44 path** in order and, for each one, write one sentence explaining what real-world thing it lets a wallet distinguish (e.g., which cryptocurrency, which sub-account).

### Medium

4. **Explain the security argument for hardened derivation** in your own words: specifically, describe what an attacker could do if they obtained one *normal* (non-hardened) child private key together with its parent's public key, and why hardened derivation prevents that specific attack.

5. **A user wants two completely separate "wallets" inside one app** — one for everyday spending, one for long-term savings — without generating two separate seed phrases. Using only concepts from Section 6, explain which BIP-44 path level they should use to separate the two, and why that level (rather than, say, `address_index`) is the appropriate one.

6. **Trace the full pipeline from Section 7 by hand** for a hypothetical 128-bit entropy value, writing down, in your own words, what data exists at each of the six stages (random bytes, mnemonic, seed, master key, derived key, address) and which BIP is responsible for producing it from the stage before it.

### Hard

7. **Research (from general cryptography knowledge) what "hardened" derivation loses** compared to normal derivation in terms of functionality, and describe a real product feature (like the receipt-printing web server mentioned in Section 5) that becomes impossible once every level of a path is hardened. Argue why BIP-44 nonetheless hardens the first three levels rather than none.

8. **Suppose two different, unrelated companies both build BIP-44-compliant wallets** but each hardcodes a different `coin_type'` value for the same coin due to a documentation mistake. Explain, step by step, what would go wrong if a user generated a seed phrase in Wallet A and imported it into Wallet B, even though both wallets correctly implement BIP-32 and BIP-39.

9. **Design (on paper) an extension to the BIP-44 path** that would support a hypothetical "shared family wallet" where three people each have a numbered sub-tree of addresses, but any of them can independently prove ownership of the funds in their own sub-tree using only their own seed. Identify which existing path level(s) you would repurpose, and explain any new hardened boundary you would introduce and why.
