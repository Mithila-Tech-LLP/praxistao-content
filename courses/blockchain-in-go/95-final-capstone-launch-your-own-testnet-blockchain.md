# Chapter 95: Final Capstone — Launch Your Own Testnet Blockchain

Ninety-four chapters ago, Chapter 01 asked you to imagine a shared notebook and a chain of wax-sealed envelopes. Every chapter since has replaced one more piece of that analogy with real, working code: real hashing, real signatures, real proof of work, real peer-to-peer networking, real storage, real smart contracts, real APIs, and — across this final volume — a real deployment. This chapter does not teach a new concept. It is the moment everything you have built stops being thirteen separate volumes of exercises and becomes one single, live thing: a multi-node GoChain testnet running on real servers, reachable by real strangers, mining real (worthless, by design) coins, monitored, backed up, and shipped through the exact CI/CD pipeline you built in Chapter 92. This is the finale.

## Table of Contents

1. [What "Capstone" Means Here](#1-what-capstone-means-here)
2. [The Complete System, Diagrammed](#2-the-complete-system-diagrammed)
3. [Walking the Whole Stack, Bottom to Top](#3-walking-the-whole-stack-bottom-to-top)
4. [A Full Live Trace: One Transaction's Entire Journey](#4-a-full-live-trace-one-transactions-entire-journey)
5. [Inviting Someone Else to Join, for Real](#5-inviting-someone-else-to-join-for-real)
6. [Watching a Stranger's Node Sync From Scratch](#6-watching-a-strangers-node-sync-from-scratch)
7. [What Could Still Go Wrong — an Honest Accounting](#7-what-could-still-go-wrong--an-honest-accounting)
8. [Capstone Project: Live GoChain Testnet](#8-capstone-project-live-gochain-testnet)
9. [Where to Go From Here](#9-where-to-go-from-here)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What "Capstone" Means Here

A capstone, literally, is the final stone placed at the top of an arch or a wall — not load-bearing on its own, but the piece whose placement proves every stone beneath it was cut and set correctly, because the whole structure only stands if everything underneath actually holds. That is exactly the spirit of this chapter. Nothing here is new machinery: every diagram, every command, and every piece of code below is something Chapters 01 through 94 already built, tested, and explained on its own. What this chapter adds is the thing no single earlier chapter could show you — what it looks like when all of it runs *at once*, on real infrastructure, with a real stranger's node syncing against it.

If you have been reading this course without actually running the code — a legitimate way to learn the concepts — this is the chapter where that stops being enough. The exercises at the end of this chapter, and the checklist in Section 8, are not thought experiments. They ask you to actually stand up a live testnet, invite someone else to it, and watch it work. That is the only way this capstone actually completes.

---

## 2. The Complete System, Diagrammed

Before the checklist, here is every single component this volume has built, and every connection between them, in one diagram. Nothing in this picture is hypothetical — every box is a container or a process you have already run in an earlier chapter; every arrow is a network connection you have already tested at least once.

```
                                     THE PUBLIC INTERNET
                                              |
                    +-------------------------+-------------------------+
                    |                         |                         |
              DNS lookups              a stranger's laptop        your own laptop
        explorer.gochain.example      (Priya, Section 6)          (gochain-wallet CLI)
        faucet.gochain.example                |                         |
        grafana.gochain.example                |                         |
                    |                         |                         |
                    v                         v                         v
   +-------------------------------------------------------------------------------+
   |  CLOUD VM  (Chapter 88)         203.0.113.42          ufw: 22, 80, 443, 9000  |
   |                                                                                 |
   |   +-------------------------------------------------------------------------+  |
   |   |                    Caddy reverse proxy (Chapter 93)                      |  |
   |   |     :80 -> redirect to :443        :443 -> TLS-terminated, routes by      |  |
   |   |     domain to the correct backend service below                          |  |
   |   +-------------------------------------------------------------------------+  |
   |       |                    |                     |                             |
   |       v                    v                     v                             |
   |  +----------+       +-------------+       +-------------+                      |
   |  | explorer |       |   faucet    |       |   grafana   |                      |
   |  | :8090    |       |   :8091     |       |   :3000     |                      |
   |  |(Ch72/75) |       |   (Ch90)    |       |   (Ch91)    |                      |
   |  +----------+       +-------------+       +-------------+                      |
   |       |                    |                     |                             |
   |       |  reads chain data  |  submits signed      |  queries                    |
   |       |  via node1's API   |  tx via node1's API  |  metrics                    |
   |       v                    v                     v                             |
   |  +--------------------------------------------------------------------------+  |
   |  |                      docker-compose.yml network (Ch87)                    |  |
   |  |                                                                            |  |
   |  |   +-----------+      +-----------+      +-----------+                     |  |
   |  |   |  node1    |<---->|  node2    |<---->|  node3    |   P2P gossip &       |  |
   |  |   |  (seed)   |<---------------------------->|  sync (Ch44-51)     |  |
   |  |   |  :8080 API|      |  :8080 API|      |  :8080 API|                     |  |
   |  |   |  :9000 P2P|      |  :9000 P2P|      |  :9000 P2P|                     |  |
   |  |   +-----------+      +-----------+      +-----------+                     |  |
   |  |       |  ^                |  ^                |  ^                        |  |
   |  |       |  |                |  |                |  |     each node runs:    |  |
   |  |       |  |                |  |                |  |     - core (Ch16-22)   |  |
   |  |       |  |                |  |                |  |     - consensus (Ch23-28)|  |
   |  |       |  |                |  |                |  |     - crypto (Ch08-15)  |  |
   |  |       v  |                v  |                v  |     - mempool (Ch34-35) |  |
   |  |    /metrics            /metrics            /metrics    - storage (Ch53-58) |  |
   |  |   (Ch91)               (Ch91)               (Ch91)     - vm (Ch59-68)      |  |
   |  |       |                    |                    |     - api (Ch70-71)      |  |
   |  |       +--------------------+--------------------+                          |  |
   |  |                            |                                                |  |
   |  |                            v                                                |  |
   |  |                     +-------------+                                         |  |
   |  |                     | prometheus  |  scrapes every node's /metrics          |  |
   |  |                     |   :9090     |  every 15s (Ch91)                       |  |
   |  |                     +-------------+                                         |  |
   |  |                            ^                                                |  |
   |  |                            | PromQL queries                                 |  |
   |  |                            +---------------- grafana (above)                |  |
   |  |                                                                            |  |
   |  |   named Docker volumes (Ch86): node1-data, node2-data, node3-data,          |  |
   |  |   caddy-data, grafana-data - each persists independently of container      |  |
   |  |   restarts, and is the exact input to Chapter 94's backup script            |  |
   |  +--------------------------------------------------------------------------+  |
   |                                                                                 |
   |   /opt/gochain/scripts/backup.sh  --(cron, nightly, Ch94)--> local + offsite   |
   |   object storage bucket, independent of this VM's own disk                     |
   +-------------------------------------------------------------------------------+
                    ^
                    |  docker compose pull && docker compose up -d
                    |  (SSH, triggered automatically on a tagged release)
                    |
   +-------------------------------------------------------------------------------+
   |                      GitHub  (your gochain repository)                        |
   |                                                                                 |
   |   .github/workflows/ci.yml       - runs on every push: go test -race -cover,  |
   |                                     go vet, gofmt check, docker build (Ch92)   |
   |                                                                                 |
   |   .github/workflows/release.yml  - runs on v*.*.* tags: test again, build,     |
   |                                     push to ghcr.io/you/gochain, then deploy    |
   |                                     over SSH to the VM above (Ch92)             |
   +-------------------------------------------------------------------------------+
```

Every arrow in that diagram is something you have personally typed a command to test at least once across the last nine chapters. The only thing new in this chapter is seeing all of it exist simultaneously.

---

## 3. Walking the Whole Stack, Bottom to Top

It helps to walk this diagram once, bottom to top, narrating what each layer is actually responsible for — partly as a review, and partly because this is exactly the mental model you will use for the rest of your career whenever you are handed an unfamiliar production system and need to understand it quickly.

**The GitHub layer** is where every change to GoChain starts and is verified before it is trusted. Nothing reaches the VM without first passing `go test ./... -race -cover` (Chapter 92) — the entire test suite you wrote across all thirteen volumes runs, automatically, on every single push.

**The VM layer** (Chapter 88) is the one piece of real, billed infrastructure this whole system runs on — a rented Linux server with a public IP address, a firewall that opens only the ports genuinely needed (22, 80, 443, 9000 — notice 8080 itself is no longer directly exposed to the public internet by this final chapter, since Caddy and the node's own P2P port are now the only public-facing surfaces that matter; more on this in Section 8).

**The Docker Compose layer** (Chapter 87) is where every application-level component actually runs: three GoChain nodes gossiping and syncing with each other over Chapter 44-51's P2P protocol, an explorer and faucet consuming the same public API any external client would use, and Prometheus/Grafana watching all of it.

**Caddy** (Chapter 93) is the single front door for every HTTP-speaking service in that layer — the only thing standing between the raw internet and your explorer, faucet, and Grafana dashboard, terminating HTTPS once so nothing downstream of it ever has to think about certificates.

**The backup layer** (Chapter 94), running quietly on a schedule underneath everything else, is the part nobody sees working until the one day something goes wrong — at which point it is the only thing that matters.

Every layer above is independently something you already built and understood in isolation. The capstone is confirming, out loud, that you can trace a single request or a single block all the way from one end of this diagram to the other, without getting lost — which is exactly what Section 4 does next.

---

## 4. A Full Live Trace: One Transaction's Entire Journey

To make Section 2's diagram concrete, here is one single transaction's entire life, traced through every layer it touches, with the exact chapter that built each step:

```
 1. A new user visits https://faucet.gochain.example and requests coins
    for their address.                                          (Ch90, Ch93)

 2. Caddy terminates TLS, forwards the plain HTTP request to the
    faucet container over Docker's internal network.                  (Ch93)

 3. The faucet's rate limiter checks the address and requester IP
    against their last payout time.                                   (Ch90)

 4. Allowed: the faucet signs a transaction with its own dedicated
    wallet's private key.                                       (Ch13, Ch33)

 5. The faucet POSTs the signed transaction to node1's /tx/send
    API endpoint, exactly like any external client would.             (Ch70)

 6. node1 validates the transaction - checks the signature, checks
    for double-spends against the mempool and UTXO set - and, if
    valid, adds it to its mempool.                                (Ch33, Ch34)

 7. node1 gossips the new transaction to its peers, node2 and
    node3, who each independently validate it again before
    forwarding it further.                                       (Ch48)

 8. A miner (any of the three nodes, whichever finds a valid nonce
    first) includes the transaction in its next block and solves
    the proof-of-work puzzle.                                    (Ch24-27)

 9. That node broadcasts the new block; the other two nodes
    validate it - correct hash, correct previous-block link,
    correct proof of work - and, finding it valid, append it to
    their own chains and remove the now-mined transaction from
    their mempools.                                          (Ch18-19, Ch49)

10. Every node's AddBlock call updates its gochain_block_height and
    gochain_mempool_size metrics at the exact moment the underlying
    values change.                                                    (Ch91)

11. Fifteen seconds later (at most), Prometheus scrapes all three
    nodes' /metrics endpoints and stores the new values.               (Ch91)

12. Grafana's dashboard, open in a browser at
    https://grafana.gochain.example, redraws its Block Height and
    Mempool Size panels live, without anyone refreshing the page. (Ch91, Ch93)

13. gochain-explorer, polling node1's API, shows the new block and
    the new user's incoming transaction on its recent-blocks page
    and on the user's own address page.                          (Ch72, Ch75)

14. The new user checks their balance, either through the explorer's
    UI or by querying node1's API directly, and sees their faucet
    payout confirmed.                                                 (Ch70)
```

Fourteen steps, fourteen different chapters' worth of code, and not one manual intervention anywhere in the middle. This is what "the system genuinely, independently works" means in concrete terms — not that each piece works in isolation, which you already knew, but that the full chain of custody from a stranger's browser click to a confirmed, monitored, explorable transaction holds together end to end.

---

## 5. Inviting Someone Else to Join, for Real

Chapter 90 walked through the mechanics of an external node joining your network. This chapter asks you to actually do it — not simulate it, not read about Priya doing it, but hand your testnet's real seed address to a real other person and watch, live, as their node becomes part of your network.

Send them exactly this, adapted to your own domain and IP:

```markdown
Join the GoChain testnet!

1. Install Go 1.22+ and clone https://github.com/you/gochain
2. Build the CLI:  go build -o gochain ./cmd/gochain
3. Start your node, pointed at our seed:

   ./gochain node start \
     --api-addr 0.0.0.0:8080 \
     --p2p-addr 0.0.0.0:9000 \
     --seed 203.0.113.42:9000 \
     --data ./data

4. Get some free test coins:
   curl -X POST https://faucet.gochain.example/faucet/request \
     -d '{"address": "<your address from `gochain wallet new`>"}'

5. Watch yourself show up in the explorer:
   https://explorer.gochain.example
```

Notice this invitation only needs five lines of instruction, and every single one of them is a direct, unmodified callback to a specific earlier chapter: cloning and building (Chapter 06), starting a node with a seed flag (Chapter 47), requesting from the faucet (Chapter 90), and checking the explorer (Chapter 75). Nothing about onboarding a stranger required inventing new mechanisms — the entire system was already designed, from Volume 7 onward, to make "a new, independent participant joins" an ordinary, expected event rather than a special case.

---

## 6. Watching a Stranger's Node Sync From Scratch

While your invited participant's node starts up, watch both sides simultaneously — their terminal and your `docker compose logs -f node1` — and narrate what you see against Chapter 49's synchronization design:

```
 their terminal:
 2026-08-01 14:02:11  node started, api=0.0.0.0:8080 p2p=0.0.0.0:9000
 2026-08-01 14:02:11  dialing seed 203.0.113.42:9000
 2026-08-01 14:02:12  handshake complete with 203.0.113.42 (their height: 4,213, our height: 0)
 2026-08-01 14:02:12  requesting blocks 1-4213
 2026-08-01 14:02:12  syncing... 500/4213
 2026-08-01 14:02:14  syncing... 2000/4213
 2026-08-01 14:02:16  syncing... 4213/4213 - fully synced
 2026-08-01 14:02:16  peer exchange: learned about 2 additional peers (node2, node3)

 your node1's logs, at the same time:
 2026-08-01 14:02:12  inbound connection from 198.51.100.7
 2026-08-01 14:02:12  handshake complete with 198.51.100.7 (their height: 0, our height: 4213)
 2026-08-01 14:02:12  peer 198.51.100.7 requested blocks 1-4213, streaming
 2026-08-01 14:02:16  peer 198.51.100.7 fully synced
```

The moment their node reports "fully synced" and their `/peers` endpoint lists your `node1`, `node2`, and `node3`, they are — genuinely, not symbolically — a full peer in your network, capable of mining their own blocks, gossiping transactions they hear about first, and serving chain data to the next stranger who points a node at *them*. Confirm it from both directions:

```bash
# On your server:
curl http://localhost:8080/peers
# {"peers": ["node2:9000", "node3:9000", "198.51.100.7:9000"]}

# On their machine:
curl http://localhost:8080/peers
# {"peers": ["203.0.113.42:9000"]}
```

If you want to go one step further, ask them to mine a block themselves, and watch it show up in your explorer within seconds — the single most convincing proof this course can offer that GoChain is a real, decentralized system rather than a demo that only ever worked because you controlled every machine involved.

---

## 7. What Could Still Go Wrong — an Honest Accounting

A capstone chapter that only celebrated success would be doing you a disservice. Worth being explicit about what this testnet is, and is not, before calling it done:

- **This is a testnet, not a production financial system.** `gochip`s are worthless by design, and nothing in this course's threat model has been reviewed with the rigor a system holding real value would need — recall Chapter 76-85's security and real-world case-study material as exactly the next layer of scrutiny a production chain would require, not a checklist this capstone claims to satisfy.
- **A handful of nodes is not meaningfully decentralized.** Chapter 51's Sybil and eclipse attack discussion still applies in full: a network with three or four nodes, several of which you personally operate, has nowhere near the resilience of a network with thousands of independent operators. This capstone proves the *mechanism* works, not that it is resistant to a well-resourced adversary at scale.
- **Single points of failure remain.** One VM, one Caddy instance, one Prometheus instance — Chapter 89's Kubernetes overview exists precisely because a system that needs to survive individual machine failures, not just individual container failures, needs a different operational model than this chapter's Docker Compose stack provides.
- **The faucet's economics do not scale to real abuse.** Chapter 90's rate limiter is a reasonable deterrent against casual misuse, not a defense against a genuinely motivated, well-resourced attacker with access to many distinct IP ranges.

None of this diminishes what you have built. It is the difference between "I built a working combustion engine and can explain exactly how every part of it works" and "I built a car that has passed every safety regulation in every country in the world" — the first is what this course delivers, honestly and completely; the second is a different, much larger undertaking that real production blockchain teams spend years on.

---

## 8. Capstone Project: Live GoChain Testnet

This is the actual assignment. Everything above was preparation and narration; this section is the checklist you follow, in order, to make it real.

### Full Deployment Checklist

1. **Confirm GoChain itself is complete and tested.** Run `go test ./... -race -cover` locally and confirm every package from every volume passes. *(Chapters 07-85, all of GoChain's own code.)*
2. **Confirm the Dockerfile builds a minimal image.** `docker build -t gochain:latest .` and check `docker images gochain` reports well under 50MB. *(Chapter 86.)*
3. **Confirm the multi-node Compose stack works locally.** `docker compose up --build`, then verify `/peers` on all three local nodes and a full send-and-mine cycle through the local explorer. *(Chapter 87.)*
4. **Provision a real cloud VM** with a public IP address, create a non-root `deploy` user, and install Docker. *(Chapter 88.)*
5. **(Optional) Review the Kubernetes overview** if you anticipate needing more than a handful of nodes or multiple physical machines — not required for this checklist to be considered complete. *(Chapter 89.)*
6. **Configure the VM's firewall**, initially opening only 22 (SSH), 8080 (API), and 9000 (P2P) via `ufw`. *(Chapter 88.)*
7. **Clone the repository to `/opt/gochain`** on the VM and bring the Compose stack up there with `docker compose up -d --build`. *(Chapter 88.)*
8. **Verify real internet reachability** from a network other than your own (mobile data, a friend's connection). *(Chapter 88.)*
9. **Write and publish your seed-node address** (`<your-vm-ip>:9000`) in a `TESTNET.md` file in your repository. *(Chapter 90.)*
10. **Generate a dedicated faucet wallet**, fund it via genesis allocation or mining, and deploy `gochain-faucet` as its own Compose service. *(Chapter 90.)*
11. **Test the faucet end to end**: request a payout, confirm the rate limiter rejects an immediate repeat request, and confirm a mined payout reflects correctly in a balance query. *(Chapter 90.)*
12. **Instrument GoChain with the five Prometheus metrics** (`gochain_block_height`, `gochain_mempool_size`, `gochain_peer_count`, `gochain_mining_hashrate`, `gochain_api_request_duration_seconds`) and confirm `/metrics` on port 8080 reports real values. *(Chapter 91.)*
13. **Deploy Prometheus and Grafana** as Compose services, confirm Prometheus is scraping every node, and build the five-panel dashboard. *(Chapter 91.)*
14. **Set up `.github/workflows/ci.yml`**, confirm it runs and passes on a push, and deliberately break something small to confirm it correctly fails. *(Chapter 92.)*
15. **Set up `.github/workflows/release.yml`**, add the `TESTNET_HOST`/`TESTNET_SSH_USER`/`TESTNET_SSH_KEY` repository secrets, and cut a real tagged release (`git tag v0.1.0 && git push origin v0.1.0`), confirming the image is built, pushed to GHCR, and automatically deployed to your VM. *(Chapter 92.)*
16. **Register a domain** (or subdomain) and add A records for `explorer.`, `faucet.`, and `grafana.` subdomains pointing at your VM's IP. *(Chapter 93.)*
17. **Deploy Caddy** with a `Caddyfile` reverse-proxying each of those three subdomains, open ports 80/443 in `ufw`, and confirm valid, browser-trusted HTTPS on all three URLs. *(Chapter 93.)*
18. **Implement and schedule the backup script** via cron, confirm a real backup file appears in `/opt/gochain/backups` and syncs to offsite object storage. *(Chapter 94.)*
19. **Run the chaos test** against a real node in your live testnet, and confirm it prints `PASS` — a genuinely restored, resynced node, not a simulation. *(Chapter 94.)*
20. **Invite a real other person to run a node against your published seed address**, and confirm, from both sides, that `/peers` lists each other and their node's height converges with the rest of the network. *(This chapter, Sections 5-6.)*
21. **Watch one full transaction's journey live** — faucet request, mempool, mining, gossip, explorer, and Grafana dashboard update — matching Section 4's trace against your own terminal output and browser windows, side by side.

### The Complete Architecture, One More Time

Section 2 already gave you the full diagram in detail. Here it is once more, collapsed to its essential shape — the picture worth keeping in your head as the one-sentence architecture of everything you built:

```
   strangers, your invited      DNS (domain names)      GitHub (source + CI/CD)
   participant, your own              |                          |
   wallet/browser                     v                          v
        |                    +----------------+          +--------------+
        +------------------->|  Caddy (TLS)    |          |  ci.yml       |
                              +----------------+          |  release.yml  |
                                /     |      \             +--------------+
                          explorer faucet  grafana                |
                                \     |      /                     | deploy over SSH
                                 v    v     v                      v
                        +--------------------------+      (pulls + restarts
                        |  3+ GoChain nodes          |       the stack below)
                        |  (P2P gossip + sync,        |
                        |   mining, mempool, VM)       |<---------------------+
                        +--------------------------+
                                     |
                                     v
                            Prometheus + Grafana
                            (metrics, dashboards)
                                     |
                                     v
                          nightly backups -> offsite storage
                          (chaos-tested restore procedure)
```

### What "Done" Looks Like

"Done," for this capstone, is not a green checkmark on a to-do list. It is a specific, observable state of the world, and you will know you have reached it because every one of the following is simultaneously true, at the same moment, without you doing anything to force it:

Your GoChain testnet is running on a real server you do not need to be logged into for it to keep working. A person you did not personally hand a laptop to has their own, independently-running GoChain node, synced to your chain, discovered through nothing but a published address and a piece of documentation. Anyone with the URL — not just you — can open `https://explorer.gochain.example` in a browser, see a valid certificate, and browse real, currently-mining blocks. That same anyone can request free test coins from your faucet and receive them within the time it takes to mine one block. You can open `https://grafana.gochain.example` and watch block height, mempool size, peer count, and hashrate update live, in real time, without refreshing anything. A `git push` of a passing change, followed by a `git tag`, results in a new version running on the live server a few minutes later, with no SSH session, no manual rebuild, and no manual restart on your part. And if the disk under any single node failed at this exact moment, you have already personally proven — not assumed — that you could bring it back and watch it resync.

That is what done looks like. Everything above this line is you, personally, having verified it with your own hands, not read about someone else doing it.

---

## 9. Where to Go From Here

Finishing this capstone does not mean there is nothing left to learn — it means you have crossed the specific threshold this course was built to get you across: from "I have heard of blockchain" to "I have built one, deployed it, broken it on purpose, and fixed it." A few honest directions worth naming for whatever comes next:

- **Harden the security model.** Revisit Chapter 76's attack lab and Chapter 85's real-world incidents with your own live system in mind — what would a real 51% attack against your specific three-or-four-node testnet actually require, concretely?
- **Grow the network for real.** Recruit more independent operators, on more diverse infrastructure and geography, and watch how Chapter 51's peer-diversity concerns stop being theoretical.
- **Take the Kubernetes path seriously**, if the testnet outgrows a handful of Compose-managed nodes — Chapter 89 was deliberately left as an on-ramp, not a dead end.
- **Read real production codebases** — `go-ethereum`, Bitcoin Core — with the vocabulary this course gave you. Chapters 81-83's comparisons were written for exactly this moment: you now have a genuine basis for comparison, not just a reading list.
- **Build the next thing.** Every skill this course exercised — careful data modeling, concurrent systems design, network protocol design, cryptographic reasoning, production deployment discipline — transfers directly to systems that have nothing to do with blockchains at all.

---

## Summary

- This chapter introduces no new mechanisms — it is the integration of all ninety-four preceding chapters into one live, observable system.
- The complete architecture spans a domain and TLS-terminating reverse proxy (Caddy), a multi-node Docker Compose stack (nodes, explorer, faucet), a monitoring stack (Prometheus/Grafana), a CI/CD pipeline (GitHub Actions), and a backup/restore process — each independently built in Chapters 86-94.
- Tracing one transaction end to end — faucet request, mempool validation, gossip, mining, sync, metrics, and explorer display — touches code from at least fourteen distinct earlier chapters, with no manual steps anywhere in the middle.
- The capstone project is a genuine deployment checklist, not a simulation: it asks you to actually provision infrastructure, publish a real seed address, and invite a real other person to join.
- A stranger's node joining your testnet runs through exactly Chapter 47's handshake and Chapter 49's synchronization logic, unmodified — proof that the networking layer never needed to know or care who, geographically, was on the other end.
- Being honest about scope matters: this is a testnet proving a mechanism works, not a production financial system that has been reviewed at the level real value at stake would demand.
- "Done" is a specific, verifiable state — a live server, an independent participant, a real HTTPS explorer and faucet, a live dashboard, an automated deploy pipeline, and a personally-tested disaster recovery process — not a checklist you can complete by reading about it.
- Finishing this capstone marks the transition from understanding blockchain concepts to having built, deployed, broken, and repaired a real one — a foundation that transfers far beyond blockchain specifically.

---

## Exercises

### Easy

1. Walk through the deployment checklist in Section 8 and, for each of the 21 items, write down which specific earlier chapter you would need to revisit if that step failed. This is a comprehension check, not a deployment step — you should be able to do this from memory or a quick skim, without redeploying anything.

2. Draw (on paper or in any tool you like) your own version of Section 2's architecture diagram, from memory, then compare it against the original and note anything you missed or misremembered.

3. Re-read Section 4's fourteen-step transaction trace and, for each step, name the specific Go file or package (per this course's structure) that code would live in.

### Medium

4. Actually complete items 1 through 13 of the Section 8 checklist: get GoChain tested, containerized, running as a local multi-node Compose stack, deployed to a real cloud VM, with a working faucet and a working Grafana dashboard. Do not proceed to CI/CD or TLS yet — this exercise is specifically about confirming the "core" deployment works before layering automation and public access on top.

5. Complete items 14 through 17: wire up CI/CD, cut a real tagged release, register a domain, and get valid HTTPS working on all three subdomains. Screenshot (or otherwise document) the green checkmarks on your GitHub Actions runs and the padlock icon in a browser visiting your explorer's HTTPS URL.

6. Complete items 18 and 19: implement scheduled backups and run the chaos test against your real, live deployment (not a local practice run). Paste the chaos test's final PASS/FAIL output and the before/after chain heights it reported.

### Hard

7. Complete the full checklist, items 20 and 21: recruit a genuinely independent person to run a node against your published seed address, and produce a side-by-side log capture (their terminal and your `node1`'s logs) showing the handshake and full sync happening in real time, similar to Section 6's example. Then have them submit a transaction through your faucet and mine it themselves, and confirm it appears in your explorer within one block's time.

8. Write a genuinely honest post-mortem-style retrospective (400-600 words) of your own deployment: what went wrong during the process that this course's chapters did not fully prepare you for, what you had to look up or improvise, and what you would do differently if you deployed a second, independent GoChain testnet from scratch tomorrow.

9. Pick one item from Section 7's "what could still go wrong" list and actually address it beyond this course's scope: for example, recruit at least five genuinely independent node operators across different hosting providers and geographies to meaningfully test Chapter 51's peer-diversity assumptions at a larger scale than this course's checklist requires, or perform a deliberate, controlled simulation of a 51% attack (Chapter 76) against your own live testnet and document exactly what real damage it could and could not do. Write up what you learned in enough detail that a future reader of this course could use it as a starting point for their own investigation.
