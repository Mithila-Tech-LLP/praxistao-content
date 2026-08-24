# Chapter 75: Mini Project — Block Explorer API

Chapters 72 and 73 built two halves of the same thing: a backend that answers browsing questions, and a frontend of static files that renders the answers. Right now those two halves live as separate concerns — a Go server on one side, loose HTML files you'd have to serve from *somewhere* on the other. This chapter closes that gap by building `gochain-explorer`: a single, self-contained Go binary that serves both the JSON API and the bundled frontend from the exact same process, using Go's `embed` package to compile the frontend's HTML, JS, and CSS directly into the executable. The result is a tool you can hand someone as one file — point it at any running GoChain node's RPC address, and it just works, with nothing else to install or configure. This is also, deliberately, the shape Part 6 needs: a single static binary is the easiest possible thing to put inside a Docker container.

## Table of Contents

1. [What We're Building](#1-what-were-building)
2. [Project Layout](#2-project-layout)
3. [Embedding the Frontend with `embed`](#3-embedding-the-frontend-with-embed)
4. [Talking to a Remote GoChain Node](#4-talking-to-a-remote-gochain-node)
5. [The Explorer Server](#5-the-explorer-server)
6. [Configuration](#6-configuration)
7. [`main.go` and the CLI Flags](#7-maingo-and-the-cli-flags)
8. [Mini Project: `gochain-explorer`](#mini-project-gochain-explorer)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What We're Building

Think about how you'd hand a friend a slideshow presentation. You could send them the individual slide files and a note saying "you'll also need this particular viewer app installed" — technically complete, practically annoying. Or you could export the whole thing as one self-contained PDF that opens anywhere, no separate app required. `gochain-explorer` is the PDF version of the block explorer: instead of "here's a Go server, and separately, here are some HTML files you need to serve from somewhere," it's one binary that does both.

```
   BEFORE (Chapters 72 + 73, two separate concerns)

   +------------------------+       +---------------------------+
   |  gochain node process   |       |  explorer-frontend/ files  |
   |  (Chapters 70-72 API)    |       |  (served by... what,       |
   |  /explorer/...            |       |   exactly? a second        |
   +------------------------+       |   process? a CDN?)          |
                                      +---------------------------+

   AFTER (this chapter)

   +----------------------------------------------------------+
   |                  gochain-explorer  (one binary)             |
   |                                                              |
   |   embedded frontend files            explorer API client     |
   |   (index.html, explorer.js, ...)      (calls a REMOTE node)  |
   |            |                                    |            |
   |            +---------------+--------------------+           |
   |                            |                                 |
   |                     one http.Server                          |
   +----------------------------|---------------------------------+
                                 |
                     HTTP requests to any GoChain node
                     (--node-addr flag, Section 6)
```

The crucial design decision here is that `gochain-explorer` does **not** open a `core.Blockchain` directly, and does **not** import `gochain/core` at all. It is a pure HTTP client of a *separately running* GoChain node's Chapter 70-72 API — the same relationship any external wallet or exchange backend has with a node. This matters for deployment: you can run one GoChain node and point several `gochain-explorer` instances at it (say, one per region, behind a load balancer), or point a single explorer at any node on the network, public or private, without recompiling anything.

---

## 2. Project Layout

```
gochain/cmd/gochain-explorer/
├── main.go              -- entry point, flag parsing, starts the server
├── config.go             -- Config struct and flag/env loading
├── client.go               -- thin HTTP client for a remote GoChain node
├── server.go                -- http.Handler wiring: API proxy + embedded frontend
└── frontend/                  -- the exact static files from Chapter 73
    ├── index.html
    ├── block.html
    ├── tx.html
    ├── address.html
    ├── explorer.js
    └── style.css
```

Nothing in `frontend/` changes from Chapter 73 — this chapter's whole job is packaging, not rewriting the explorer's pages. The one adjustment: Chapter 73's pages called `fetch('/explorer/blocks?...')` assuming the API and the frontend shared an origin. That assumption becomes literally true here, since `gochain-explorer` serves both from the same `http.Server` on the same port — no CORS configuration needed at all, which is one of the nicer side effects of bundling both halves into one process.

---

## 3. Embedding the Frontend with `embed`

Go's standard library `embed` package, available since Go 1.16, lets a program include arbitrary files from disk directly inside the compiled binary — no separate deployment step for static assets, no risk of the binary and its frontend files drifting apart on some server because someone forgot to copy one of them.

```go
// gochain/cmd/gochain-explorer/server.go
package main

import (
	"embed"
	"io/fs"
	"net/http"
)

// frontendFS embeds every file under frontend/ directly into the
// compiled gochain-explorer binary at build time. //go:embed is a
// compiler directive, not a runtime call -- the files listed become
// part of the binary's own data section, exactly like a string constant
// would, just much larger.
//
//go:embed frontend
var frontendFS embed.FS

// frontendHandler strips the "frontend/" prefix embed.FS always adds
// (since the directive embeds the whole frontend/ directory, not just
// its contents) so that a request for "/index.html" correctly resolves
// to the embedded "frontend/index.html" file.
func frontendHandler() (http.Handler, error) {
	sub, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}
```

`fs.Sub` is the one subtlety worth understanding here: `//go:embed frontend` embeds the *directory itself*, so every path inside `frontendFS` starts with `frontend/` (e.g., `frontend/index.html`). `fs.Sub(frontendFS, "frontend")` re-roots the filesystem so that `index.html` (no prefix) is what `http.FileServer` sees — otherwise every URL a visitor typed would need an unwanted `/frontend/` prefix, which is not what Chapter 73's pages (or a visitor's browser) expect.

---

## 4. Talking to a Remote GoChain Node

`gochain-explorer` needs to answer its own `/explorer/...` requests by forwarding them to whatever real GoChain node it's configured to point at — it holds no blockchain data of its own.

```go
// gochain/cmd/gochain-explorer/client.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// nodeClient is a thin wrapper around http.Client that knows one thing:
// the base URL of the GoChain node to forward requests to. It has no
// knowledge of core.Blockchain, core.Transaction, or any other internal
// GoChain type -- as far as this package is concerned, a GoChain node is
// just an HTTP server that happens to speak the Chapter 70-72 API.
type nodeClient struct {
	baseURL string
	http    *http.Client
}

func newNodeClient(baseURL string) *nodeClient {
	return &nodeClient{baseURL: baseURL, http: &http.Client{}}
}

// proxy forwards an incoming request's path and query string to the
// configured node, and copies the node's response back verbatim --
// status code, headers, and body -- so gochain-explorer's own handlers
// never need to know the shape of any specific endpoint's JSON. It is,
// deliberately, a dumb pipe: Chapter 72 already validated and shaped
// every response once; this package's only job is relaying it.
func (c *nodeClient) proxy(w http.ResponseWriter, r *http.Request) {
	target := c.baseURL + r.URL.Path
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad upstream request: %v", err), http.StatusBadGateway)
		return
	}
	req.Header = r.Header.Clone()

	resp, err := c.http.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("upstream node unreachable: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// validateBaseURL is a small startup-time sanity check -- catching a
// malformed --node-addr immediately, at startup, is far friendlier than
// discovering it only once the first visitor's request fails.
func validateBaseURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid --node-addr %q: must be a full URL like http://localhost:8080", raw)
	}
	return nil
}
```

`proxy` is intentionally "dumb": it does not decode the JSON, does not reshape it, does not add caching — it just relays bytes. Every real decision about what an endpoint returns was already made once, correctly, in Chapter 72's handlers; duplicating that logic here would be exactly the kind of "two places doing the same job, destined to quietly drift apart" problem Chapter 70 warned against for its own REST/JSON-RPC split.

---

## 5. The Explorer Server

```go
// gochain/cmd/gochain-explorer/server.go (continued)

// newExplorerServer wires together the embedded frontend (Section 3) and
// the proxying node client (Section 4) into one http.Handler. Every
// request either matches an /explorer/... (or /blocks, /balance, /rpc)
// API path -- forwarded to the real node -- or falls through to serving
// a static file from the embedded frontend.
func newExplorerServer(nc *nodeClient) (http.Handler, error) {
	frontend, err := frontendHandler()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	// Every API path this explorer's frontend calls gets proxied
	// straight through to the configured GoChain node -- listed
	// explicitly (rather than "proxy everything") so a request for an
	// unrecognized path falls through to the frontend's 404 page
	// instead of silently forwarding arbitrary traffic to the node.
	apiPrefixes := []string{"/explorer/", "/blocks/", "/balance/", "/rpc", "/health"}
	for _, prefix := range apiPrefixes {
		p := prefix // capture for the closure
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			nc.proxy(w, r)
		})
	}

	// Anything else -- index.html, block.html, explorer.js, style.css --
	// is served directly from the embedded filesystem.
	mux.Handle("/", frontend)

	return mux, nil
}
```

Registering explicit prefixes (`/explorer/`, `/blocks/`, and so on) rather than a single catch-all proxy is a deliberate safety choice: it means `gochain-explorer` only ever forwards the exact set of paths its own frontend is written to call, and cannot be tricked into relaying arbitrary traffic to the upstream node — a small but real difference between "a proxy for this specific frontend" and "an open proxy to whatever node it points at."

---

## 6. Configuration

```go
// gochain/cmd/gochain-explorer/config.go
package main

import (
	"fmt"
	"os"
)

// Config holds everything gochain-explorer needs to start, gathered
// from CLI flags with environment-variable fallbacks -- a common
// pattern for tools that need to run both interactively (a developer
// typing flags) and inside a container (Part 6), where flags are
// awkward but environment variables are the norm.
type Config struct {
	ListenAddr string // e.g. ":9090" -- where THIS explorer serves HTTP
	NodeAddr   string // e.g. "http://localhost:8080" -- the GoChain node to query
}

// loadConfig reads flags (already parsed by main.go into these two
// variables) and falls back to environment variables, then hardcoded
// defaults, in that order -- flags win because they are the most
// explicit and immediate way to configure a single run.
func loadConfig(listenFlag, nodeFlag string) (Config, error) {
	cfg := Config{
		ListenAddr: firstNonEmpty(listenFlag, os.Getenv("GOCHAIN_EXPLORER_LISTEN"), ":9090"),
		NodeAddr:   firstNonEmpty(nodeFlag, os.Getenv("GOCHAIN_NODE_ADDR"), "http://localhost:8080"),
	}
	if err := validateBaseURL(cfg.NodeAddr); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func (c Config) String() string {
	return fmt.Sprintf("listen=%s node=%s", c.ListenAddr, c.NodeAddr)
}
```

The flag-then-env-then-default precedence here previews exactly the configuration convention Part 6's Docker chapters lean on: a container typically has no convenient way to pass `--flags` at `docker run` time without extra scripting, but setting environment variables (`-e GOCHAIN_NODE_ADDR=http://gochain-node:8080`) is the normal, idiomatic way containers get configured — building that fallback in now means Chapter 86's Dockerfile needs zero changes to this binary's code.

---

## 7. `main.go` and the CLI Flags

```go
// gochain/cmd/gochain-explorer/main.go
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	listenFlag := flag.String("listen", "", "address to serve the explorer on (e.g. :9090)")
	nodeFlag := flag.String("node-addr", "", "GoChain node to query (e.g. http://localhost:8080)")
	flag.Parse()

	cfg, err := loadConfig(*listenFlag, *nodeFlag)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	nc := newNodeClient(cfg.NodeAddr)
	handler, err := newExplorerServer(nc)
	if err != nil {
		log.Fatalf("failed to build explorer server: %v", err)
	}

	log.Printf("gochain-explorer serving on %s, querying node at %s", cfg.ListenAddr, cfg.NodeAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, handler))
}
```

This standalone `cmd/gochain-explorer` binary deliberately uses plain `flag`, not Chapter 74's Cobra tree — it is its own small, single-purpose tool (one job: serve the explorer), not a growing family of related subcommands, which is exactly the situation Chapter 74 itself argued `flag` is still perfectly fine for.

---

## Mini Project: `gochain-explorer`

Putting the whole thing together: build the binary, point it at a running GoChain node, and open it in a browser.

```bash
# Terminal 1: run an actual GoChain node (Chapter 74's node start command,
# or any earlier chapter's equivalent) so there's real data to explore.
$ gochain node start --miner-address gochain1alice... --addr :8080

# Terminal 2: build and run the explorer, pointed at that node.
$ go build -o gochain-explorer ./cmd/gochain-explorer
$ ./gochain-explorer --listen :9090 --node-addr http://localhost:8080
2026/08/01 09:30:11 gochain-explorer serving on :9090, querying node at http://localhost:8080
```

Opening `http://localhost:9090` in a browser now shows Chapter 73's homepage, listing recent blocks fetched live from the node on port 8080 — but everything, frontend and API-proxying alike, is served from the single `gochain-explorer` binary on port 9090. A quick sanity check confirms the embedding actually worked, with no separate file dependency at runtime:

```bash
$ ls
gochain-explorer          # <- one file

$ mv gochain-explorer /tmp/ && cd /tmp
$ ./gochain-explorer --node-addr http://localhost:8080
2026/08/01 09:31:02 gochain-explorer serving on :9090, querying node at http://localhost:8080
# still works -- no frontend/ directory needed alongside the binary,
# because the files are compiled INTO it, not read from disk at runtime.
```

Pointing the very same binary at a different node — a friend's, or a public testnet — is just a different flag value, with no rebuild required:

```bash
$ ./gochain-explorer --node-addr http://198.51.100.20:8080
2026/08/01 09:32:40 gochain-explorer serving on :9090, querying node at http://198.51.100.20:8080
```

This is exactly the property Part 6 needs: a `Dockerfile` for `gochain-explorer` only ever needs to `COPY` a single compiled binary and set an environment variable for `GOCHAIN_NODE_ADDR` — no volume mount for frontend assets, no risk of the container's frontend files silently going stale relative to a rebuilt image.

---

## Summary

- `gochain-explorer` packages Chapter 72's backend and Chapter 73's frontend into a single deployable Go binary, using `embed` to compile the frontend's static files directly into the executable.
- The explorer holds no blockchain data of its own — it is a pure HTTP client of a separately running GoChain node, proxying `/explorer/...`, `/blocks/...`, `/balance/...`, and `/rpc` requests straight through and serving everything else from the embedded frontend.
- `//go:embed frontend` plus `fs.Sub` re-roots the embedded filesystem so requests like `/index.html` resolve correctly without an unwanted `frontend/` prefix.
- The proxy forwards requests and responses as raw bytes rather than decoding and re-encoding JSON, so there is exactly one place (Chapter 72's own handlers) that decides what any endpoint's response looks like.
- Registering explicit API path prefixes, rather than proxying every request indiscriminately, keeps `gochain-explorer` from being usable as an open relay to whatever node it's pointed at.
- Configuration follows a flag-then-environment-variable-then-default precedence, so the exact same binary configures naturally both for a developer running it locally and for a container in Part 6, with zero code changes.
- Serving both the frontend and the API-proxy from one `http.Server` and one origin eliminates any need for CORS configuration between them.
- Because both the API-forwarding logic and the frontend files are compiled into one binary, `gochain-explorer` can be moved to any machine and pointed at any GoChain node with a single flag — no separate deployment step for static assets.

---

## Exercises

### Easy

1. Add a `/healthz` endpoint to `gochain-explorer` itself (distinct from the proxied `/health`, which reports the upstream *node's* health) that returns `{"status": "ok", "node": "<configured node-addr>"}, confirming the explorer process itself is alive even if the upstream node is unreachable.
2. `loadConfig` currently exits with `log.Fatalf` on an invalid `--node-addr`. Add a startup self-check that also performs a real HTTP request to the configured node's `/health` endpoint before serving any traffic, logging a clear warning (not a fatal error) if the node doesn't respond, since the node might simply not be running *yet*.
3. Add a `--version` flag to `gochain-explorer` (plain `flag`-based, matching this chapter's style) that prints a hardcoded version string and exits, without starting the server.

### Medium

4. `proxy`'s explicit `apiPrefixes` list currently must be kept in sync by hand with whatever paths Chapter 73's `explorer.js` actually calls. Write a small test that parses `frontend/explorer.js` for string literals matching an API-looking pattern (e.g., starting with `/explorer/`, `/blocks/`, etc.) and asserts every one of them is covered by an entry in `apiPrefixes`, catching drift between the two automatically.
5. Add response caching in front of the proxy specifically for `GET /blocks/{hash}` and `GET /explorer/blocks/{hash}/transactions` requests for blocks that are clearly not the current tip (you'll need one extra call to the node to learn the current height) — since a mined block's contents are immutable, these are safe to cache aggressively in `gochain-explorer` itself, taking real load off the upstream node.
6. Add a `--node-addrs` flag (plural) accepting a comma-separated list of node addresses, and round-robin `nodeClient.proxy` across them, so `gochain-explorer` can spread read traffic across several GoChain nodes instead of depending on exactly one.

### Hard

7. Right now, if the upstream node is unreachable, every proxied request fails individually with a `502`. Implement a small circuit breaker: after some number of consecutive proxy failures, `gochain-explorer` should stop attempting to reach the node for a cooldown period and immediately return a clear "upstream temporarily unavailable" response instead, then automatically try again after the cooldown — explain, in a comment, why this protects a struggling upstream node from being hammered by an explorer that keeps retrying every single visitor's request.
8. Extend `gochain-explorer` to also proxy Chapter 71's WebSocket endpoint (`GET /ws`), so a browser connecting to `ws://localhost:9090/ws` transparently gets live block/transaction events relayed from the upstream node's own WebSocket connection — research Go's approach to proxying a `net/http` WebSocket upgrade and implement it.
9. Containerize `gochain-explorer` in a minimal multi-stage `Dockerfile` (a `golang:1.22` build stage producing the static binary, copied into a `scratch` or `distroless` final image) and confirm the resulting image runs with nothing but `docker run -p 9090:9090 -e GOCHAIN_NODE_ADDR=http://host.docker.internal:8080 gochain-explorer`, with no volume mounts required for the frontend — this is a preview of Chapter 86's full containerization treatment, done here for just this one binary.
