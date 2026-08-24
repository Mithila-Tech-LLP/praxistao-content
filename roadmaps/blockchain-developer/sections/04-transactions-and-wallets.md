---
title: Transactions & Wallets
---
Blocks are only interesting because of what's inside them: transactions moving value between wallets. This section covers the UTXO model, building and signing transactions, and the mempool that holds unconfirmed ones.

### Transactions & the UTXO Model
Bitcoin-style chains don't track "account balances" directly — they track unspent transaction outputs (UTXOs), and your balance is just the sum of the ones you can prove you own.

**Resources:**
- [What Is a Transaction](course:blockchain-in-go#29-what-is-a-transaction)
- [The UTXO Model Explained](course:blockchain-in-go#30-the-utxo-model-explained)

### Building & Signing Transactions
A transaction is data plus a signature proving whoever created it actually owns the funds being spent. See both halves implemented in Go.

**Resources:**
- [Building Transactions in Go](course:blockchain-in-go#32-building-transactions-in-go)
- [Signing and Verifying Transactions](course:blockchain-in-go#33-signing-and-verifying-transactions)

### The Mempool & Double-Spend Prevention
Before a transaction makes it into a block, it waits in the mempool — and that's exactly where you have to stop someone from spending the same coins twice.

**Resources:**
- [The Mempool and Preventing Double-Spending](course:blockchain-in-go#34-the-mempool-and-preventing-double-spending)

### CLI Wallets
A wallet is really just a key pair plus a friendly interface for creating and sending transactions. Build the interface here.

**Resources:**
- [Building a CLI Wallet](course:blockchain-in-go#36-building-a-cli-wallet)

### Practice: The Full CLI Wallet
> branches-from: CLI Wallets

Build the wallet task of the standalone project: generate keys, check balances, and send transactions from the command line.

**Resources:**
- [Build a Blockchain project](project:build-a-blockchain)

### Major Project: Send & Receive Coins
Bring transactions, signatures, the mempool, and the wallet together into a working send-and-receive flow — the first point where GoChain feels like real money moving around.

**Resources:**
- [Major Project 1: Send and Receive Coins](course:blockchain-in-go#37-major-project-1-send-and-receive-coins)

### HD Wallets & Mnemonics
> optional

Real wallets don't make you back up one key at a time — a single seed phrase can deterministically regenerate every key you'll ever need.

**Resources:**
- [HD Wallets Explained](course:blockchain-in-go#38-hd-wallets-explained)
- [Implementing Mnemonic Seed Phrases](course:blockchain-in-go#39-implementing-mnemonic-seed-phrases)
