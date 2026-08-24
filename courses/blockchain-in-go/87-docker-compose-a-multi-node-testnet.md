# Chapter 87: Docker Compose — A Multi-Node Testnet

Chapter 52's major project ran a real multi-node GoChain network by opening several terminal tabs, each starting a node with its own set of flags, and manually keeping track of which port belonged to which node. This chapter replaces every one of those terminal tabs with a single file and a single command — `docker compose up` — that starts an entire, pre-wired, multi-node testnet, including a block explorer, in one shot.

## Table of Contents

1. [From Several Terminals to One Command](#1-from-several-terminals-to-one-command)
2. [What Docker Compose Is](#2-what-docker-compose-is)
3. [Docker's Internal DNS and Service Discovery](#3-dockers-internal-dns-and-service-discovery)
4. [Designing the Testnet Topology](#4-designing-the-testnet-topology)
5. [The docker-compose.yml File](#5-the-docker-composeyml-file)
6. [Explaining Each Service](#6-explaining-each-service)
7. [Seed Nodes, Revisited](#7-seed-nodes-revisited)
8. [Bringing the Testnet Up](#8-bringing-the-testnet-up)
9. [Verifying the Nodes Found Each Other](#9-verifying-the-nodes-found-each-other)
10. [Adding the Block Explorer](#10-adding-the-block-explorer)
11. [Tearing Down and Data Persistence](#11-tearing-down-and-data-persistence)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. From Several Terminals to One Command

Recall Chapter 52's setup: three separate `gochain node start` processes, each in its own terminal, each with a different `--api-addr`, `--p2p-addr`, and `--seed` flag pointing at the others by `127.0.0.1:<port>`. That works, but it does not scale past a handful of manual terminal windows, it is easy to get a flag wrong, and it gives you no easy way to tear the whole thing down and bring it back up identically later.

Chapter 86 solved "how do I package one node." This chapter solves "how do I start several of them, pre-configured to find each other, with one command" — which is precisely what running a realistic testnet, and every later chapter in this volume, requires.

---

## 2. What Docker Compose Is

**Docker Compose** is a tool that reads a single YAML file describing a set of containers — which images to run, what ports to publish, what environment variables or command-line flags to pass, how they should be networked together — and starts (or stops) the *entire group* as one unit, with one command. Where Chapter 86's `docker run` started one container by hand, Compose lets you describe "three GoChain nodes and an explorer, wired together like this" once, in a file, and reproduce that exact setup forever after with `docker compose up`.

```
   docker-compose.yml                    docker compose up
  +--------------------+                +------------------------+
  | services:          |                |                        |
  |   node1:  ...      |   -------->    |  [node1]  [node2]      |
  |   node2:  ...      |                |  [node3]  [explorer]   |
  |   node3:  ...      |                |                        |
  |   explorer: ...     |                |  all started, networked|
  +--------------------+                |  together, in one shot|
                                         +------------------------+
```

The file itself is just YAML (a human-readable configuration format, already familiar from any earlier `.yaml` config in this course) — no new language to learn, only a small, fixed vocabulary of keys Compose understands: `services`, `build`, `image`, `ports`, `volumes`, `networks`, `command`, `depends_on`, and a few others you will see below.

---

## 3. Docker's Internal DNS and Service Discovery

When Compose starts a group of services, it automatically creates a private virtual network for them and gives every container an entry in that network's internal **DNS** (Domain Name System — the system that translates human-readable names into machine-reachable addresses). Concretely: a container named `node2` can reach a container named `node1` simply by connecting to the hostname `node1`, on whatever port `node1` is listening on *inside* the network — no IP addresses, no manual configuration, no port-mapping tricks required between containers on the same Compose network.

```
                 Docker's internal DNS
                 (built into the bridge network)

   node2  --- "what's node1's address?" --->  DNS
   node2  <--- "10.0.5.3" -------------------  DNS
   node2  ---------------- connects to 10.0.5.3:9000 -------> node1
```

This is the single mechanism that makes Chapter 52's manual `--seed 127.0.0.1:3001` flags unnecessary here: instead of a loopback address and a specific host port, every node's `--seed` flag simply says `node1:9000` — the *service name* from the compose file, on the *container-internal* port — and Docker resolves it automatically, regardless of which host machine or which random internal IP the container actually ends up with.

---

## 4. Designing the Testnet Topology

Before writing YAML, decide on the shape of the network this chapter builds: three GoChain nodes, one of them (`node1`) acting as the network's seed node (the first-known peer everyone else connects through, exactly the role described in Chapter 47), and a `gochain-explorer` instance pointed at `node1`'s API.

```
                         gochain-net (Docker bridge network)

  +-----------+        +-----------+        +-----------+
  |   node1   |<------>|   node2   |        |   node3   |
  | (seed)    |<-----------------------------> |         |
  | :8080 API |        | :8080 API |        | :8080 API |
  | :9000 P2P |        | :9000 P2P |        | :9000 P2P |
  +-----------+        +-----------+        +-----------+
        ^
        | HTTP API calls
        |
  +--------------+
  |   explorer   |
  |  :8090 HTTP  |
  +--------------+

  Host machine ports published outward:
   localhost:8080 -> node1:8080     localhost:9000 -> node1:9000
   localhost:8081 -> node2:8080     localhost:9001 -> node2:9000
   localhost:8082 -> node3:8080     localhost:9002 -> node3:9000
   localhost:8090 -> explorer:8090
```

Each node still gets its own host-side port mapping (so you can `curl` any of them individually from your laptop, exactly as in Chapter 86), but *inside* the `gochain-net` network, every service always listens on the same conventional ports — 8080 and 9000 — since container-to-container traffic never touches the host's port numbers at all.

---

## 5. The docker-compose.yml File

Here is the full file, placed at the repository root next to the `Dockerfile` from Chapter 86:

```yaml
# docker-compose.yml — a 3-node GoChain testnet plus a block explorer,
# all built from the same image, wired together with Docker's built-in
# service-name DNS instead of manual IP/port bookkeeping.

services:
  node1:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gochain-node1
    # node1 has no --seed flag: it is this testnet's seed node, the first
    # peer everyone else connects to (Chapter 47's seed-node role).
    command: >
      node start
      --api-addr 0.0.0.0:8080
      --p2p-addr 0.0.0.0:9000
      --data /app/data
    ports:
      - "8080:8080"
      - "9000:9000"
    volumes:
      - node1-data:/app/data
    networks:
      - gochain-net

  node2:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gochain-node2
    # node1:9000 resolves via Docker's internal DNS to node1's
    # container-internal address - no IP addresses hardcoded anywhere.
    command: >
      node start
      --api-addr 0.0.0.0:8080
      --p2p-addr 0.0.0.0:9000
      --data /app/data
      --seed node1:9000
    ports:
      - "8081:8080"
      - "9001:9000"
    volumes:
      - node2-data:/app/data
    networks:
      - gochain-net
    depends_on:
      - node1

  node3:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gochain-node3
    command: >
      node start
      --api-addr 0.0.0.0:8080
      --p2p-addr 0.0.0.0:9000
      --data /app/data
      --seed node1:9000
    ports:
      - "8082:8080"
      - "9002:9000"
    volumes:
      - node3-data:/app/data
    networks:
      - gochain-net
    depends_on:
      - node1

  explorer:
    build:
      context: .
      dockerfile: Dockerfile.explorer
    container_name: gochain-explorer
    # gochain-explorer (Chapter 75) is pointed at node1's API by service
    # name - it never needs to know node1's real IP address.
    command: ["--node", "http://node1:8080", "--listen", "0.0.0.0:8090"]
    ports:
      - "8090:8090"
    networks:
      - gochain-net
    depends_on:
      - node1

networks:
  gochain-net:
    driver: bridge

volumes:
  node1-data:
  node2-data:
  node3-data:
```

This assumes a second, small Dockerfile, `Dockerfile.explorer`, built the same multi-stage way as Chapter 86's `Dockerfile` but compiling `./cmd/gochain-explorer` (the self-contained explorer binary with its embedded frontend from Chapter 75) instead of `./cmd/gochain`:

```dockerfile
# Dockerfile.explorer — same multi-stage shape as Chapter 86's Dockerfile,
# just pointed at the explorer's own main package.
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/gochain-explorer ./cmd/gochain-explorer

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/gochain-explorer /usr/local/bin/gochain-explorer
EXPOSE 8090
ENTRYPOINT ["gochain-explorer"]
```

---

## 6. Explaining Each Service

- **`build: { context: ., dockerfile: Dockerfile }`** — rather than pulling a pre-built image from a registry, Compose builds the image itself from the local `Dockerfile`, the same one Chapter 86 wrote. All three node services reuse this exact same build, since they are all just GoChain nodes with different flags.
- **`command: >`** — the `>` YAML syntax folds a multi-line block into a single space-joined string, purely for readability in the file; it produces the same command Chapter 86's `CMD` would, just with an explicit `--seed` flag appended for `node2` and `node3`.
- **`ports: - "8081:8080"`** — host-port-to-container-port mapping, exactly as in Chapter 86's `docker run -p`, just expressed in YAML. Notice each node maps to a *different* host port (8080/8081/8082) even though every container listens on the same internal `8080` — this is required because your host machine cannot have three different processes all bound to host port 8080 simultaneously, even though three different *containers* can each privately use 8080 internally without conflict.
- **`volumes: - node1-data:/app/data`** — a distinct named volume per node, so each node's blockchain data, wallet files, and UTXO index (Volume 8) are kept completely separate, exactly modeling three independently-operated machines rather than three processes secretly sharing one disk.
- **`networks: - gochain-net`** — places every service on the same custom bridge network, which is what makes the DNS-by-service-name behavior from Section 3 available between them. Compose would create a default network automatically even without this, but naming it explicitly makes the topology diagram in Section 4 match the file exactly.
- **`depends_on: - node1`** — tells Compose to start `node1` before `node2`, `node3`, and `explorer`. This only controls *start order*, not "wait until node1's API is actually ready" — Section 9 covers why nodes and the explorer still need their own retry logic for that gap.

---

## 7. Seed Nodes, Revisited

Chapter 47 introduced the seed-node pattern: a new node needs at least one already-known peer address to bootstrap from, after which peer-exchange messages let it discover the rest of the network organically. In this compose file, `node1` plays exactly that seed role, and `node2`/`node3` are handed `--seed node1:9000` as their one bootstrap address.

The only thing that has changed from Chapter 52's version is *what the seed address looks like*: a bare service name and a fixed internal port (`node1:9000`) instead of a loopback IP and an arbitrary host port (`127.0.0.1:53001`). The underlying peer-discovery and handshake code inside `gochain/network` is completely unaware of the difference — it just dials whatever address string it was given, and Docker's DNS makes `node1` resolve correctly without the node's own code needing to know anything about containers at all. This is a useful thing to notice: good separation of concerns (Volume 7's `network.Node` never hardcoding how addresses are resolved) is exactly what makes deployment changes like this one require zero changes to GoChain's actual Go code.

---

## 8. Bringing the Testnet Up

From the repository root:

```bash
docker compose up --build
```

`--build` forces Compose to (re)build the images from the Dockerfiles rather than using any previously cached versions — worth doing the first time, and any time you've changed GoChain's source since the last build. Expected output (abbreviated, with logs from all four containers interleaved by default):

```
[+] Running 4/4
 ✔ Container gochain-node1     Started
 ✔ Container gochain-node2     Started
 ✔ Container gochain-node3     Started
 ✔ Container gochain-explorer  Started

gochain-node1     | 2024/03/01 10:00:01 node started, api=0.0.0.0:8080 p2p=0.0.0.0:9000
gochain-node2     | 2024/03/01 10:00:02 node started, api=0.0.0.0:8080 p2p=0.0.0.0:9000
gochain-node2     | 2024/03/01 10:00:02 dialing seed node1:9000
gochain-node2     | 2024/03/01 10:00:02 handshake complete with node1 (height 0)
gochain-node3     | 2024/03/01 10:00:02 node started, api=0.0.0.0:8080 p2p=0.0.0.0:9000
gochain-node3     | 2024/03/01 10:00:02 dialing seed node1:9000
gochain-node3     | 2024/03/01 10:00:02 handshake complete with node1 (height 0)
gochain-explorer  | 2024/03/01 10:00:03 explorer listening on 0.0.0.0:8090, node=http://node1:8080
```

To run it detached (in the background, freeing your terminal) instead, use `docker compose up -d --build`, and follow logs on demand with `docker compose logs -f`.

---

## 9. Verifying the Nodes Found Each Other

With the testnet running, confirm from your host machine that all three nodes are up and have connected to each other:

```bash
curl http://localhost:8080/peers   # node1's view of its peers
# {"peers": ["node2:9000", "node3:9000"]}

curl http://localhost:8081/peers   # node2's view
# {"peers": ["node1:9000"]}
```

Then exercise the actual point of a network: submit a transaction against `node2`'s API, mine it on `node3`, and confirm the resulting block shows up on `node1` too, exactly mirroring Chapter 52's convergence test but across containers instead of terminal tabs:

```bash
curl -X POST http://localhost:8081/tx/send -d '{"to": "...", "amount": 10}'
curl -X POST http://localhost:8082/mine
curl http://localhost:8080/chain/height
# {"height": 1}   <- node1 received and validated the block node3 mined
```

If `node2` or `node3` never show a peer, the most common cause is starting them before `node1`'s listener is actually accepting connections yet, even though `depends_on` started the container itself — `depends_on` only orders container *start*, not application *readiness*. GoChain's own peer-connection code should already retry a failed dial a few times with a short backoff (a good habit from Volume 7); if it does not yet, that is worth revisiting before relying on this compose file for anything beyond local experimentation.

---

## 10. Adding the Block Explorer

With the testnet's peers confirmed, point your browser at `http://localhost:8090` — the `gochain-explorer` container, built in Section 5, pointed at `node1`'s API via the internal `http://node1:8080` address. You should see the same block you just mined in Section 9 listed as the chain's most recent block, its transaction detail page showing the transfer you submitted, and the recipient address's page showing its updated balance — the exact explorer views built in Chapter 72, now running against a real, multi-container network with a single `docker compose up`, instead of a node running bare on your laptop.

This is the direct payoff of this chapter: what used to be "open four terminals, remember four sets of flags, and hope you typed the seed address correctly" is now one file, checked into version control, that anyone on your team (or, from Chapter 88 onward, a real cloud server) can bring up identically.

---

## 11. Tearing Down and Data Persistence

To stop everything:

```bash
docker compose down
```

This stops and removes all four containers and the `gochain-net` network, but — because the named volumes (`node1-data`, `node2-data`, `node3-data`) are declared separately under the top-level `volumes:` key — it does **not** delete their contents. Running `docker compose up` again later restarts the same testnet with the same blockchain history intact on each node.

To wipe everything, including the chain data itself (useful when you want a genuinely fresh testnet, for instance before Chapter 90's public testnet setup), add the `-v` flag:

```bash
docker compose down -v
```

---

## Summary

- Docker Compose reads one YAML file describing a group of services and starts, networks, and stops them all together, replacing Chapter 52's multi-terminal manual setup with a single `docker compose up`.
- Compose automatically creates a private network for a project's services and gives each one a DNS entry matching its service name, so containers reach each other by name (`node1`) instead of by IP address or host port.
- This chapter's `docker-compose.yml` builds three GoChain node services (from Chapter 86's `Dockerfile`) and one explorer service (from a parallel `Dockerfile.explorer`), each with its own published host ports and its own named data volume.
- `node1` plays the seed-node role from Chapter 47; `node2` and `node3` are configured with `--seed node1:9000`, using the service name and internal port rather than a loopback address.
- `depends_on` only controls container start order, not application readiness — GoChain's own peer-dial retry logic (from Volume 7) is what actually makes nodes wait for each other successfully.
- `docker compose up --build` builds and starts the whole testnet; `docker compose logs -f` follows all services' logs; `docker compose down` stops it while preserving data; `docker compose down -v` also deletes the underlying volumes.
- Verifying the testnet means checking `/peers` on each node and running an actual send-and-mine transaction end to end, then confirming it propagates to every node and shows up correctly in the explorer.
- This exact compose file, unmodified in its core shape, is what Chapter 88 deploys onto a real cloud VM.

---

## Exercises

### Easy

1. **Bring the testnet up** with `docker compose up --build`, then run `docker compose ps` to list the four running services and their port mappings. Confirm all four show a healthy/running state.

2. **Curl the `/peers` endpoint** on all three nodes (ports 8080, 8081, 8082) and write down what each one reports. Confirm every node's peer list contains at least `node1`, and that `node1`'s own list contains both `node2` and `node3`.

3. **Submit a transaction and mine it**, following Section 9's example, then open `http://localhost:8090` in a browser and confirm the transaction and resulting block appear correctly in the explorer.

### Medium

4. **Add a fourth node, `node4`**, to the compose file, following the exact pattern of `node2`/`node3` (its own host ports, its own named volume, `--seed node1:9000`). Bring the testnet up and confirm, via `/peers`, that `node4` successfully joins and that the other three nodes eventually learn about it too (recall Chapter 48's gossip and Chapter 47's peer-exchange messages — `node4` may reach the others indirectly, not only through `node1`).

5. **Run `docker compose down` (without `-v`) and then `docker compose up` again**, without rebuilding, and confirm each node's chain height and balances are exactly as they were before teardown, proving the named volumes preserved state correctly. Then run `docker compose down -v` and bring the testnet back up, confirming every node starts from a fresh genesis block.

6. **Change `node2` and `node3`'s seed to each other instead of `node1`** (`--seed node3:9000` for node2, `--seed node2:9000` for node3), leaving `node1` with no seed and no `depends_on` from the others. Bring the testnet up and investigate, using `/peers` on each node, whether `node1` ever gets discovered by `node2`/`node3` at all — and explain, in 100-150 words, what this reveals about the difference between "a seed node" and "the only entry point into a network."

### Hard

7. **Simulate a node crash and recovery**: with the testnet running and several blocks mined, run `docker compose kill node3` (a hard kill, no graceful shutdown), wait, then `docker compose start node3`. Confirm `node3` catches back up to the current chain height using the synchronization logic from Chapter 49, and record how long it takes to fully resync relative to how many blocks it missed.

8. **Convert the three fixed node services into a single scalable service** using Compose's `deploy.replicas` (or the older `--scale` flag with a template-style service definition), where each replica gets a different container name and port mapping automatically. Explain, in writing, what breaks about the current `--seed node1:9000` approach once nodes are dynamically scaled up and down, and propose (without necessarily implementing) a fix.

9. **Add a `healthcheck:` block** to each node service that calls its own `/chain/height` endpoint, and a `condition: service_healthy` variant of `depends_on` for `node2`, `node3`, and `explorer` so they only start dialing `node1` once it is confirmed ready, rather than relying purely on GoChain's own retry logic. Test that this removes or reduces the "no peers found" failure mode discussed in Section 9, and write a short comparison of this approach versus in-application retry logic — including when you would still want both.
