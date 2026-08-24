# Chapter 86: Containerizing a Node with Docker

Every GoChain node you have run so far has depended on whatever happened to be installed on your laptop: your specific Go version, your specific operating system, whatever files happened to already exist in `./data`. That is fine for development, but it is exactly the kind of hidden dependency that breaks the moment you try to run GoChain on a cloud server, hand it to a teammate, or run three nodes side by side without them tripping over each other's files. This chapter fixes that permanently by packaging GoChain into a Docker container — a small, self-contained, portable unit that runs identically anywhere Docker itself runs.

## Table of Contents

1. [Why Containerize a Blockchain Node](#1-why-containerize-a-blockchain-node)
2. [Docker Concepts in Plain Language](#2-docker-concepts-in-plain-language)
3. [The Problem With a Single-Stage Build](#3-the-problem-with-a-single-stage-build)
4. [Multi-Stage Builds, Explained](#4-multi-stage-builds-explained)
5. [Writing the GoChain Dockerfile](#5-writing-the-gochain-dockerfile)
6. [Explaining the Dockerfile, Line by Line](#6-explaining-the-dockerfile-line-by-line)
7. [Building the Image](#7-building-the-image)
8. [Running Your First Containerized Node](#8-running-your-first-containerized-node)
9. [Connecting From Your Host Machine](#9-connecting-from-your-host-machine)
10. [Inspecting, Logs, and Cleaning Up](#10-inspecting-logs-and-cleaning-up)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Containerize a Blockchain Node

Think about everything a working GoChain node currently depends on: a specific version of Go installed to compile it, a `go.mod` with the exact right dependency versions, a `./data` directory that has to exist and be writable, and specific ports free on your machine. Every one of those is a way a node can work perfectly on your laptop and fail mysteriously on someone else's — the classic "it works on my machine" problem.

**Docker** is a tool that packages an application together with everything it needs to run — the compiled binary, any files it depends on, the exact environment it expects — into a single unit called a **container**, which behaves the same way no matter what computer it runs on. Rather than telling a teammate "install Go 1.22, clone this repo, run `go build`, then run the binary with these three flags," you tell them "run this one Docker command." Rather than hoping a cloud server has the right Go version installed, you ship it a container that already contains a working binary.

For GoChain specifically, this matters for three concrete reasons that the rest of this volume all builds on:

1. **Volume 7's multi-node testnet** required opening several terminal tabs, each running its own `gochain node start` with different flags. Chapter 87 replaces that with a single file that starts several *containers* at once — each one an isolated, identical copy of the same image.
2. **Volume 13's later chapters** — deploying to a cloud VM (Chapter 88), Kubernetes (Chapter 89), and CI/CD (Chapter 92) — all assume GoChain ships as a container. A cloud provider, a Kubernetes cluster, and a GitHub Actions pipeline all speak the same language: container images, not "clone and `go build`."
3. **Reproducibility.** A container built today and a container built in six months from the same source produce the same behavior, because the exact runtime environment is baked into the image, not assembled fresh from whatever happens to be on the host machine at the time.

---

## 2. Docker Concepts in Plain Language

Before writing anything, four terms need concrete definitions, because every later chapter in this volume uses all four constantly.

- **Image** — a read-only template that contains an application and everything it needs to run: the compiled binary, system libraries, configuration files, and the instructions for how to start it. Think of an image like a frozen, shrink-wrapped snapshot of a fully set-up computer, ready to be duplicated. You do not run an image directly — you run *containers from* it.
- **Container** — a running instance of an image. Just like you can start a Go program from a single compiled binary multiple times (each becoming its own separate process with its own memory), you can start multiple containers from a single image, and each one runs independently, with its own filesystem view, network address, and process, even though they all came from the exact same frozen template.
- **Dockerfile** — a plain-text script of instructions (`FROM`, `RUN`, `COPY`, `CMD`, and so on) that tells Docker how to *build* an image, step by step, starting from some base image and layering your own application on top.
- **Registry** — a server that stores and serves images by name and version tag, the way GitHub stores and serves source code. Docker Hub is the most common public registry; Chapter 92 pushes GoChain's own image to one as part of CI/CD.

```
   Dockerfile                Image                  Containers
  (instructions)     -->   (frozen         -->    (running copies,
                             template)               each independent)

  FROM golang...                                   +-----------+
  RUN go build...    -->   gochain:latest   -->     | container |
  ENTRYPOINT ...                              -->   +-----------+
                                                     +-----------+
                                              -->    | container |
                                                     +-----------+
```

A container is *not* a tiny virtual machine, even though it can feel like one. A virtual machine (VM) emulates an entire computer, including its own operating system kernel — heavy, slow to start, often gigabytes in size. A container shares the host machine's kernel and only isolates the application's view of the filesystem, processes, and network — lightweight, starts in a fraction of a second, and typically measured in megabytes. This is exactly why running three GoChain "nodes" as three containers in Chapter 87 is practical, where running three full virtual machines on your laptop would be uncomfortably heavy.

---

## 3. The Problem With a Single-Stage Build

The naive way to containerize GoChain would be a single `FROM golang:1.22` image that installs the Go toolchain, copies in the source code, and runs `go build` — with the *same* image then used to run the compiled binary. This works, but it drags three problems into every deployment:

- **Size.** The official `golang` image (with the full compiler toolchain, standard library source, and build caches) is several hundred megabytes before your code is even added. GoChain's compiled binary itself is a few tens of megabytes at most. Shipping the whole toolchain to run a binary that does not need it anywhere close to doubles or triples the image size for no runtime benefit.
- **Pull time.** Every server, every CI run, and every Kubernetes pod that starts a new container has to first *download* the image if it does not already have it cached. A 900MB image takes meaningfully longer to pull than a 20MB one — this adds up across dozens of deployments.
- **Attack surface.** A full Go toolchain, plus whatever package manager and shell utilities `golang:1.22` includes, gives an attacker who somehow gets code execution inside the container far more to work with than a container that contains nothing but the compiled binary and its bare runtime dependencies.

The fix is to use two separate images for two separate jobs: one *only* for compiling the code, and a second, much smaller one *only* for running the already-compiled result. That is exactly what a multi-stage build gives you.

---

## 4. Multi-Stage Builds, Explained

A **multi-stage build** is a Dockerfile with more than one `FROM` instruction, where each `FROM` starts a fresh, independent build stage, and later stages can selectively copy specific files out of earlier stages — while everything else from those earlier stages (the compiler, source code, intermediate build artifacts) is discarded and never becomes part of the final image.

```
STAGE 1: "builder"                         STAGE 2: "runtime" (final image)
+----------------------------+             +----------------------------+
| FROM golang:1.22-alpine    |             | FROM alpine:3.19           |
|                            |             |                            |
|  full Go compiler          |   copy      |  (no compiler)             |
|  go.mod / go.sum           |   ONLY      |  (no source code)         |
|  entire source tree        |  ------->   |                            |
|  go build -> /out/gochain  |   the       |  /usr/local/bin/gochain   |
|                            |   binary    |  (~20MB total image)      |
| (discarded after build)    |             |                            |
+----------------------------+             +----------------------------+
        ~700MB+ intermediate                   ships to production
```

Everything in the "builder" stage exists purely as scaffolding — Docker builds it, runs the compile step inside it, and then, because the final stage only ever says "copy this one file from the builder stage," none of the toolchain, source, or intermediate layers make it into the image you actually deploy. The result: a final image that contains exactly one thing GoChain actually needs at runtime — the compiled `gochain` binary — plus the tiny base operating system layer required to run it at all.

---

## 5. Writing the GoChain Dockerfile

Here is the full, working `Dockerfile`, placed at the root of the `gochain` module, next to `go.mod`:

```dockerfile
# ---------- Stage 1: build ----------
FROM golang:1.22-alpine AS builder

# Alpine's base image is intentionally minimal; git and ca-certificates
# aren't included by default, but `go mod download` needs them to fetch
# modules over HTTPS from version-control hosts.
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy only the dependency manifests first. Docker caches each instruction
# as a "layer," and reuses a cached layer if its inputs haven't changed.
# Since go.mod/go.sum change far less often than application source, this
# ordering means `go mod download` (slow) is skipped on most rebuilds,
# and only reruns when a dependency actually changes.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source tree and build.
COPY . .

# CGO_ENABLED=0 disables cgo, producing a fully static binary with no
# dynamic C library dependencies - this is what lets the final image use
# a minimal base like alpine (or even `scratch`) without missing shared
# libraries at runtime. -ldflags="-s -w" strips debug symbols and the
# symbol table, shrinking the binary further since we don't need to
# attach a debugger to a production container.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/gochain ./cmd/gochain

# ---------- Stage 2: runtime ----------
FROM alpine:3.19

# The final image still needs CA certificates so any outbound HTTPS
# connections GoChain makes (e.g. Chapter 92's registry pulls, or future
# integrations) can verify TLS certificates correctly.
RUN apk add --no-cache ca-certificates

# Create a dedicated, unprivileged user. If a bug in gochain were ever
# exploited, the attacker lands in this restricted account, not root -
# a basic but important container-hardening habit.
RUN addgroup -S gochain && adduser -S gochain -G gochain

WORKDIR /app

# Copy ONLY the compiled binary from the builder stage. None of the Go
# toolchain, module cache, or source code exists in this image at all.
COPY --from=builder /out/gochain /usr/local/bin/gochain

# The directory GoChain's storage layer (Volume 8) uses for its BoltDB
# file and wallet files. Owned by the unprivileged user so it can write.
RUN mkdir -p /app/data && chown -R gochain:gochain /app/data

USER gochain

# 8080 = the REST/JSON-RPC/WebSocket API from Volume 10.
# 9000 = the P2P port other GoChain nodes dial into, from Volume 7.
# EXPOSE is documentation for humans and tooling - it does not, by
# itself, publish the port to the host; `docker run -p` does that.
EXPOSE 8080 9000

# Marks /app/data as a volume - a hint to Docker (and to anyone reading
# this file) that this directory holds state that should persist across
# container restarts and rebuilds, not be treated as disposable.
VOLUME ["/app/data"]

# ENTRYPOINT fixes the program that always runs; CMD supplies its
# default arguments, which anyone running the image can still override.
ENTRYPOINT ["gochain"]
CMD ["node", "start", "--api-addr", "0.0.0.0:8080", "--p2p-addr", "0.0.0.0:9000", "--data", "/app/data"]
```

---

## 6. Explaining the Dockerfile, Line by Line

Walking through the choices made above, by name:

- **`FROM golang:1.22-alpine AS builder`** — starts the first stage from the official Go image, using its `alpine` variant (built on the minimal Alpine Linux distribution) rather than the full Debian-based variant, since we only need this stage to compile code, not to be small — but starting small anyway keeps the build itself fast. `AS builder` names this stage so the second stage can refer back to it.
- **`RUN apk add --no-cache git ca-certificates`** — `apk` is Alpine's package manager (the same role `apt` plays on Debian/Ubuntu). `--no-cache` avoids leaving a local package index cached inside the image layer, keeping this intermediate stage (which gets discarded anyway) a little leaner.
- **`COPY go.mod go.sum ./` then `RUN go mod download`, *before* `COPY . .`** — this ordering is a deliberate Docker layer-caching trick: as long as your dependencies haven't changed, Docker reuses the cached "download modules" layer on every subsequent build, even if you've edited a hundred `.go` files, because that `COPY . .` (and everything after it) is a separate layer from the dependency download.
- **`CGO_ENABLED=0 GOOS=linux go build ...`** — explained above: this produces a static binary with no external C library dependencies, which is exactly what a minimal runtime image like `alpine` (or `scratch`, an image with *nothing* in it at all) requires to run the binary without missing-library errors.
- **`FROM alpine:3.19`** — starts the second, completely separate stage. Nothing from the builder stage exists here unless explicitly copied in with `COPY --from=builder`.
- **`RUN addgroup ... && adduser ...` and `USER gochain`** — running as a non-root user inside the container is a basic security hardening step recommended for any container that will ever run on a shared host or a cloud VM (Chapter 88 uses exactly this image on a real, internet-reachable server).
- **`COPY --from=builder /out/gochain /usr/local/bin/gochain`** — the one line that actually pulls a file across the stage boundary. This is the entire point of the multi-stage pattern: everything else the builder stage touched is left behind.
- **`EXPOSE 8080 9000`** — documents, for both humans reading the Dockerfile and tools that introspect images, that this container expects to serve traffic on GoChain's two conventional ports: 8080 for the API and 9000 for P2P.
- **`ENTRYPOINT ["gochain"]` / `CMD [...]`** — `ENTRYPOINT` is the fixed command that always runs when the container starts; `CMD` supplies default arguments to it. Anyone running this image can override the `CMD` portion (for example, to run `gochain wallet new` instead of starting a node) without needing a different Dockerfile.

---

## 7. Building the Image

From the root of the `gochain` module (next to `go.mod` and the `Dockerfile` you just wrote), build the image:

```bash
docker build -t gochain:latest .
```

`-t gochain:latest` tags the resulting image with a human-readable name and version (`latest` here just means "the most recent build" — a real release, in Chapter 92, tags images with version numbers instead). The trailing `.` tells Docker to use the current directory as the **build context** — everything Docker is allowed to `COPY` from during the build.

Expected output (abbreviated):

```
[+] Building 24.3s (16/16) FINISHED
 => [builder 1/6] FROM docker.io/library/golang:1.22-alpine
 => [builder 2/6] RUN apk add --no-cache git ca-certificates
 => [builder 3/6] WORKDIR /src
 => [builder 4/6] COPY go.mod go.sum ./
 => [builder 5/6] RUN go mod download
 => [builder 6/6] COPY . .
 => [builder 7/6] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/gochain ./cmd/gochain
 => [stage-1 1/5] FROM docker.io/library/alpine:3.19
 => [stage-1 2/5] RUN apk add --no-cache ca-certificates
 => [stage-1 3/5] RUN addgroup -S gochain && adduser -S gochain -G gochain
 => [stage-1 4/5] COPY --from=builder /out/gochain /usr/local/bin/gochain
 => [stage-1 5/5] RUN mkdir -p /app/data && chown -R gochain:gochain /app/data
 => exporting to image
 => => naming to docker.io/library/gochain:latest
```

Check the final image's size to see the multi-stage build's payoff directly:

```bash
docker images gochain
# REPOSITORY   TAG      IMAGE ID       CREATED          SIZE
# gochain      latest   a1b2c3d4e5f6   10 seconds ago   19.4MB
```

Compare that to the `golang:1.22-alpine` builder image alone, which you can check with `docker images golang` — typically 300-400MB. The final `gochain` image ships none of that.

---

## 8. Running Your First Containerized Node

Start a container from the image you just built:

```bash
docker run -d \
  --name gochain-node1 \
  -p 8080:8080 \
  -p 9000:9000 \
  -v gochain-data:/app/data \
  gochain:latest
```

- **`-d`** runs the container detached, in the background, printing its container ID and returning control of your terminal immediately.
- **`--name gochain-node1`** gives the container a fixed, memorable name instead of Docker's randomly generated one (useful the moment you are running more than one, exactly as Chapter 87 will).
- **`-p 8080:8080 -p 9000:9000`** publishes ports: the left side of each pair is the port on your host machine, the right side is the port inside the container. Without `-p`, the container's ports exist only inside Docker's internal network and are unreachable from your host at all.
- **`-v gochain-data:/app/data`** mounts a **named Docker volume** called `gochain-data` at `/app/data` inside the container. A named volume is storage that Docker manages and keeps alive independently of the container's lifecycle — if you remove and recreate the container, the blockchain data in `/app/data` survives, exactly the durability property a real node needs.

Expected output is just the new container's full ID (a long hex string), confirming it started. Verify it is actually running:

```bash
docker ps
# CONTAINER ID   IMAGE            COMMAND                  STATUS         PORTS
# 7f2a9c1d4b3e   gochain:latest   "gochain node start…"   Up 3 seconds   0.0.0.0:8080->8080/tcp, 0.0.0.0:9000->9000/tcp
```

---

## 9. Connecting From Your Host Machine

With ports published, GoChain's API is reachable from your host exactly as if the node were running natively — the container boundary is invisible to a client making an HTTP request:

```bash
curl http://localhost:8080/chain/height
# {"height": 0}
```

You can also connect a `gochain-wallet` (from Chapter 42) or `gochain-explorer` (from Chapter 75) running directly on your host machine at the same address, since as far as either tool is concerned, `localhost:8080` is just a normal HTTP server:

```bash
gochain-wallet balance --node http://localhost:8080 --address <your-address>
```

The P2P port behaves the same way: another GoChain process running natively on your laptop (or, once you reach Chapter 87, another container) can dial `localhost:9000` and complete a handshake exactly as it would with a native process, because Docker's port publishing is transparent to both sides of the TCP connection.

---

## 10. Inspecting, Logs, and Cleaning Up

A few commands you will use constantly from here on, both in this chapter and every later one in this volume:

```bash
# Follow the node's logs live, exactly like tailing a log file.
docker logs -f gochain-node1

# Open an interactive shell inside the running container, useful for
# poking around the filesystem or checking that /app/data is populated.
docker exec -it gochain-node1 sh

# Stop the container (the process inside receives SIGTERM, then SIGKILL
# after a grace period if it hasn't exited).
docker stop gochain-node1

# Remove the stopped container. The named volume (and its data) survives
# this - only `docker volume rm gochain-data` would delete the data too.
docker rm gochain-node1
```

Running `docker exec -it gochain-node1 sh` and then `ls -la /app/data` inside the container is a good sanity check the first time through this chapter — you should see GoChain's BoltDB file and any wallet files it created, owned by the `gochain` user rather than root, confirming the non-root setup from Section 6 took effect.

---

## Summary

- Docker packages an application and everything it needs into a portable **image**, from which you can start one or more independent **containers**.
- A container is lighter than a virtual machine because it shares the host's kernel rather than emulating an entire separate computer.
- A single-stage build that includes the full Go toolchain in the final image is unnecessarily large, slow to pull, and has more attack surface than the running binary actually needs.
- A **multi-stage build** uses one `FROM` stage to compile the code and a second, separate `FROM` stage to run it, copying across only the compiled binary — shrinking GoChain's final image from hundreds of megabytes to under 20MB.
- `CGO_ENABLED=0` produces a fully static binary, which is what allows the final stage to use a minimal base image like `alpine` without missing shared libraries.
- GoChain's `Dockerfile` exposes ports 8080 (API) and 9000 (P2P) by convention, runs as a non-root user, and stores durable state under `/app/data`, backed by a named Docker volume.
- `docker build -t gochain:latest .` builds the image; `docker run -d -p 8080:8080 -p 9000:9000 -v gochain-data:/app/data gochain:latest` runs it, with ports published so your host machine (and other tools like `gochain-wallet` or `gochain-explorer`) can reach it exactly as if it were a native process.
- This single containerized node is the exact building block Chapter 87 uses to run a whole multi-node testnet with one command.

---

## Exercises

### Easy

1. **Build the image** using the `docker build` command from Section 7, then run `docker images gochain` and record the exact size you get. Compare it to `docker images golang` (pull `golang:1.22-alpine` first if needed) and write one sentence stating the ratio between the two.

2. **Start a container** as shown in Section 8, then run `docker logs -f gochain-node1` and, in a separate terminal, hit `curl http://localhost:8080/chain/height` a few times while mining a block or two through any means available to you from earlier volumes. Confirm the height reported by `curl` changes.

3. **Stop and remove the container**, then start a *new* one with the same `-v gochain-data:/app/data` flag, and confirm (via `curl` or `docker exec`) that your blockchain data survived the container's removal. Then run `docker volume rm gochain-data` (after stopping and removing any container using it) and start a brand-new node, confirming the chain is now empty again.

### Medium

4. **Remove the `CGO_ENABLED=0` environment variable** from the `RUN go build` line, rebuild the image, and observe what happens when you try `docker run` it (it will likely fail at startup, or fail to even build depending on your host). Write 100-150 words explaining, in your own terms, why a dynamically-linked binary built against `golang:1.22-alpine`'s C libraries fails to run correctly against `alpine:3.19`'s runtime environment.

5. **Change the final stage's base image from `alpine:3.19` to `scratch`** (an image with literally nothing in it — no shell, no package manager, not even `/bin/sh`). You will need to remove the `RUN apk add` and `RUN addgroup/adduser` lines (since there's no package manager to run them with) and instead copy `/etc/passwd` and CA certificates from the builder stage directly. Get the image to build and run successfully, and write a short comparison of the resulting image size versus the `alpine` version.

6. **Deliberately break layer caching**: reorder the Dockerfile so `COPY . .` happens *before* `COPY go.mod go.sum ./` and `RUN go mod download`. Make a small, unrelated change to a `.go` file (not `go.mod`/`go.sum`) and rebuild twice, timing each build. Explain, using Docker's own build output, exactly which layers got rebuilt unnecessarily and why.

### Hard

7. **Add a `HEALTHCHECK` instruction** to the Dockerfile that periodically calls `/chain/height` on the running node's own API and marks the container unhealthy if it fails to respond within a reasonable time. Test it by starting a container, confirming `docker ps` shows a healthy status, then simulating a hang (e.g., by sending `SIGSTOP` to the process inside via `docker exec`) and watching the status change.

8. **Write a second Dockerfile, `Dockerfile.wallet`**, that builds and packages `gochain-wallet` (Chapter 42) using the same multi-stage pattern, and run it as a one-off container (`docker run --rm gochain-wallet:latest new`) that generates a wallet, prints the address, and exits. Explain what changes (if any) were needed to the non-root user and volume setup given that a wallet CLI, unlike a long-running node, does not need an open port.

9. **Investigate and apply Docker's build cache mount feature** (`RUN --mount=type=cache,target=/root/go/pkg/mod go mod download`, requiring BuildKit) to speed up repeated builds across entirely separate `docker build` invocations (not just layer caching within one build history). Measure and report the difference in `go mod download` time on a clean build cache versus a warm one, and explain in your own words why this is a different caching mechanism from the layer caching discussed in Exercise 6.
