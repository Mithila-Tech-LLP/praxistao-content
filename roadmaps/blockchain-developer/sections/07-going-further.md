---
title: Going Further
---
With the core chain, network, and VM built, the last stretch is about opening the chain up to the outside world and understanding how real production blockchains make different tradeoffs than GoChain does.

### APIs & Block Explorers
Nobody interacts with a node over raw TCP. A JSON-RPC/REST API and a block explorer are how humans and other software actually see what's happening on the chain.

**Resources:**
- [Building a JSON-RPC and REST API](course:blockchain-in-go#70-building-a-json-rpc-and-rest-api)
- [Building a Block Explorer Backend](course:blockchain-in-go#72-building-a-block-explorer-backend)

### Proof of Stake & BFT Consensus
> optional

Proof of Work isn't the only way to reach consensus — Proof of Stake and BFT-style protocols (used by many modern chains) make very different energy and finality tradeoffs.

**Resources:**
- [Proof of Stake: Explained and Implemented](course:blockchain-in-go#77-proof-of-stake-explained-and-implemented)
- [BFT Consensus: PBFT and Tendermint Overview](course:blockchain-in-go#78-bft-consensus-pbft-and-tendermint-overview)

### Real-World Architectures: Bitcoin & Ethereum
> optional

Compare GoChain's deliberately simplified design against the real architectures of Bitcoin and Ethereum to see which of your simplifications were reasonable — and which real chains had to solve differently at scale.

**Resources:**
- [Bitcoin Architecture Deep Dive](course:blockchain-in-go#81-bitcoin-architecture-deep-dive)
- [Ethereum Architecture Deep Dive](course:blockchain-in-go#82-ethereum-architecture-deep-dive)

### Deploying a Testnet: Docker, Kubernetes, Monitoring
Take a multi-node network from your laptop to something resembling production: containerized, orchestrated, and observable.

**Resources:**
- [Docker Compose: A Multi-Node Testnet](course:blockchain-in-go#87-docker-compose-a-multi-node-testnet)
- [Monitoring with Prometheus and Grafana](course:blockchain-in-go#91-monitoring-with-prometheus-and-grafana)

### Final Capstone: Launch Your Own Testnet
The finish line: take everything built across this roadmap — chain, consensus, wallets, networking, storage, and the VM — and launch a real public testnet.

**Resources:**
- [Final Capstone: Launch Your Own Testnet Blockchain](course:blockchain-in-go#95-final-capstone-launch-your-own-testnet-blockchain)
