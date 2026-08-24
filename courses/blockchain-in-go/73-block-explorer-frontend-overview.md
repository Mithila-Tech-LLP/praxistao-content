# Chapter 73: Block Explorer Frontend — Overview

Chapter 72 gave GoChain a backend that can answer every question a curious visitor might have: what are the recent blocks, what's in this one, what did this transaction move, what has this address ever done. None of that is useful to an actual human yet, though — a JSON blob is not something you hand your non-technical friend and say "here, check whether your payment went through." This chapter builds the thin layer on top: a handful of HTML pages, wired together with plain JavaScript `fetch()` calls, that turn Chapter 72's four endpoints into something you can open in a browser and click through. The goal here is not to teach you a frontend framework — it's to show, as directly as possible, that a block explorer's frontend is *nothing more* than "call an endpoint, render the JSON, add a link to the next page."

## Table of Contents

1. [What This Chapter Builds](#1-what-this-chapter-builds)
2. [Project Layout for a Minimal Frontend](#2-project-layout-for-a-minimal-frontend)
3. [A Shared Fetch Helper](#3-a-shared-fetch-helper)
4. [The Homepage: Recent Blocks](#4-the-homepage-recent-blocks)
5. [The Block Detail Page](#5-the-block-detail-page)
6. [The Transaction Detail Page](#6-the-transaction-detail-page)
7. [The Address Page](#7-the-address-page)
8. [Tying the Pages Together with Links](#8-tying-the-pages-together-with-links)
9. [What a Production Explorer Adds](#9-what-a-production-explorer-adds)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What This Chapter Builds

Think of Chapter 72's API as a reference library's card catalog — every card is precise and complete, but nobody browses a library by reading catalog cards one at a time. A block explorer frontend is the library's reading room: the same information, arranged into a handful of well-lit pages a visitor can wander through by simply clicking whatever catches their eye. Every page this chapter builds maps onto exactly one endpoint from Chapter 72; there is no new blockchain logic here at all, only presentation.

```
   Chapter 72 endpoint                              This chapter's page
   -----------------------------------------        --------------------------
   GET /explorer/blocks                     ------>  Homepage (recent blocks)
   GET /explorer/blocks/{hash}/transactions ------>  Block detail page
   GET /explorer/transactions/{txid}        ------>  Transaction detail page
   GET /explorer/address/{address}/txs      ------>  Address page
```

The whole frontend is four static HTML files and one small shared JavaScript file — no build step, no `npm install`, no framework. That is a deliberate choice for this chapter: once you understand that "render this JSON as a table with links" is the entire job, reaching for React, Vue, or Svelte later (which Section 9 discusses) is purely a matter of developer convenience at scale, not a fundamentally different idea.

```
                     Visitor's browser
   +------------------------------------------------------+
   |  index.html   block.html   tx.html   address.html      |
   |       |            |          |            |           |
   |       +------------+----------+------------+           |
   |                    |                                    |
   |              shared explorer.js                         |
   |        (fetchJSON, formatting helpers)                  |
   +--------------------------|-----------------------------+
                              |  fetch() -- plain HTTP requests
                              v
   +--------------------------------------------------------+
   |            gochain/api  (Chapters 70-72)                 |
   |   /explorer/blocks, /explorer/blocks/{hash}/transactions, |
   |   /explorer/transactions/{txid},                          |
   |   /explorer/address/{address}/transactions                |
   +--------------------------------------------------------+
```

---

## 2. Project Layout for a Minimal Frontend

A small, static-file directory is all this needs. There is no server-side templating here — the API returns JSON, and the browser assembles the page after the fact, entirely in JavaScript. This keeps the frontend genuinely decoupled from the backend: it could just as easily be served from a CDN, a different machine, or (as Chapter 75 does) bundled straight into the Go binary itself.

```
gochain/explorer-frontend/
├── index.html       -- homepage: recent blocks
├── block.html        -- one block's transactions
├── tx.html            -- one transaction's full detail
├── address.html       -- one address's history + balance
├── explorer.js         -- shared fetch + formatting helpers
└── style.css            -- minimal shared styling
```

Every HTML file follows the same three-part shape: a skeleton of empty containers, a `<script>` tag pulling in `explorer.js`, and a small inline script that calls one function to populate the page once the DOM is ready. Nothing here needs a bundler, because browsers have supported `fetch()`, `<script type="module">`, and template literals natively for years.

---

## 3. A Shared Fetch Helper

Every page in this chapter needs the same three things: call an endpoint, handle the "it failed" case, and turn a raw address or hash into a clickable link. Centralizing these avoids four slightly-different copies of the same fetch logic.

```js
// explorer.js -- shared by every page in this chapter

// fetchJSON wraps fetch() with the two things every page here needs:
// throwing on a non-2xx response (so callers can use a single catch
// block) and parsing the body as JSON automatically.
async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `request failed: ${res.status}`);
  }
  return res.json();
}

// shortHash trims a long hex hash down to a readable "0000ab..1234"
// form for table cells, while linkHash keeps the full value in the
// href so clicking still goes to the right place.
function shortHash(hex) {
  if (!hex || hex.length <= 14) return hex;
  return `${hex.slice(0, 8)}..${hex.slice(-6)}`;
}

// linkBlock / linkTx / linkAddress build the anchor tags every page
// uses to cross-link into the other three pages -- centralizing the
// URL shape here means changing a page's filename only requires an
// edit in one place, not four.
function linkBlock(hash) {
  return `<a href="block.html?hash=${hash}">${shortHash(hash)}</a>`;
}
function linkTx(txId) {
  return `<a href="tx.html?txid=${txId}">${shortHash(txId)}</a>`;
}
function linkAddress(address) {
  return `<a href="address.html?address=${address}">${shortHash(address)}</a>`;
}

// showError renders a friendly error message into a page's #error
// container -- every page below calls this from the same catch block.
function showError(message) {
  const el = document.getElementById('error');
  if (el) {
    el.textContent = message;
    el.style.display = 'block';
  }
}
```

Every function above does exactly one small job, and every page from here on leans on all four of them instead of reinventing link-building or error handling per page. This is the same "small shared helper, reused everywhere" instinct GoChain's own `writeJSON`/`writeError` pair followed back in Chapter 70 — it just happens to be JavaScript this time, not Go.

---

## 4. The Homepage: Recent Blocks

The homepage's entire job is to call `GET /explorer/blocks` and render the result as a table, with a "next page" link that just increments `?page=`.

```html
<!-- index.html -->
<!DOCTYPE html>
<html>
<head>
  <title>GoChain Explorer</title>
  <link rel="stylesheet" href="style.css">
  <script src="explorer.js"></script>
</head>
<body>
  <h1>GoChain Explorer</h1>
  <div id="error" style="display:none; color:red;"></div>
  <table id="blocks-table">
    <thead>
      <tr><th>Height</th><th>Hash</th><th>Txs</th><th>Timestamp</th></tr>
    </thead>
    <tbody id="blocks-body"></tbody>
  </table>
  <div id="pagination"></div>

  <script>
    // Read the page number from the URL (?page=2), defaulting to 1 --
    // exactly the same query parameter Chapter 72's API expects, so
    // this page can just forward it straight through.
    const params = new URLSearchParams(window.location.search);
    const page = parseInt(params.get('page') || '1', 10);

    async function loadBlocks() {
      try {
        const resp = await fetchJSON(`/explorer/blocks?page=${page}&pageSize=20`);
        const body = document.getElementById('blocks-body');
        body.innerHTML = resp.data.map(b => `
          <tr>
            <td>${b.height}</td>
            <td>${linkBlock(b.hash)}</td>
            <td>${b.txCount}</td>
            <td>${new Date(b.timestamp * 1000).toLocaleString()}</td>
          </tr>
        `).join('');

        // total lets us know whether a "next" link even makes sense --
        // no point linking past the last page of real data.
        const totalPages = Math.ceil(resp.total / resp.pageSize);
        const pag = document.getElementById('pagination');
        pag.innerHTML = '';
        if (page > 1) pag.innerHTML += `<a href="?page=${page - 1}">Previous</a> `;
        if (page < totalPages) pag.innerHTML += `<a href="?page=${page + 1}">Next</a>`;
      } catch (e) {
        showError(e.message);
      }
    }

    loadBlocks();
  </script>
</body>
</html>
```

Notice how little there is here beyond plumbing: read `page` from the URL, ask the API for that page, turn each block into a table row, and build "Previous"/"Next" links from `resp.total` and `resp.pageSize` — the exact same `pagedResponse` envelope Chapter 72's `writePaged` produces for every listing endpoint. Because every listing endpoint shares that one envelope shape, this pagination logic could be copy-pasted into `block.html` and `address.html` almost verbatim, which Section 8 makes explicit.

---

## 5. The Block Detail Page

Clicking a block's hash on the homepage should land on a page listing every transaction inside that block — Chapter 72's `GET /explorer/blocks/{hash}/transactions`.

```html
<!-- block.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Block Detail — GoChain Explorer</title>
  <link rel="stylesheet" href="style.css">
  <script src="explorer.js"></script>
</head>
<body>
  <a href="index.html">&laquo; Back to recent blocks</a>
  <h1>Block <span id="block-hash"></span></h1>
  <div id="error" style="display:none; color:red;"></div>
  <table>
    <thead><tr><th>Tx ID</th><th>Inputs</th><th>Outputs</th><th>Fee</th></tr></thead>
    <tbody id="tx-body"></tbody>
  </table>

  <script>
    const params = new URLSearchParams(window.location.search);
    const hash = params.get('hash');
    document.getElementById('block-hash').textContent = shortHash(hash);

    async function loadBlockTxs() {
      try {
        const resp = await fetchJSON(
          `/explorer/blocks/${hash}/transactions?page=1&pageSize=50`
        );
        const body = document.getElementById('tx-body');
        body.innerHTML = resp.data.map(tx => `
          <tr>
            <td>${linkTx(tx.txId)}</td>
            <td>${tx.inputs.length}</td>
            <td>${tx.outputs.length}</td>
            <td>${tx.fee}</td>
          </tr>
        `).join('');
      } catch (e) {
        showError(e.message);
      }
    }

    loadBlockTxs();
  </script>
</body>
</html>
```

`hash` arrives as a query parameter (`block.html?hash=0000ab...`) rather than a path segment, purely because this is a static file with no server-side router of its own — a static HTML file can't match `/block/{hash}` the way `net/http.ServeMux` can. That's a small, honest limitation of "no build step, no server-side routing," and it's exactly the kind of thing a real framework (Section 9) would clean up with client-side routing, in exchange for needing a build step.

---

## 6. The Transaction Detail Page

A transaction's detail page is the one place a visitor sees the full picture Chapter 72's `toTxResponse` assembled: every input's source address and amount, every output's destination, the fee, and the confirmation count.

```html
<!-- tx.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Transaction Detail — GoChain Explorer</title>
  <link rel="stylesheet" href="style.css">
  <script src="explorer.js"></script>
</head>
<body>
  <h1>Transaction <span id="txid"></span></h1>
  <div id="error" style="display:none; color:red;"></div>
  <p id="status"></p>
  <h3>Inputs</h3>
  <table><tbody id="inputs-body"></tbody></table>
  <h3>Outputs</h3>
  <table><tbody id="outputs-body"></tbody></table>

  <script>
    const params = new URLSearchParams(window.location.search);
    const txid = params.get('txid');
    document.getElementById('txid').textContent = shortHash(txid);

    async function loadTx() {
      try {
        const tx = await fetchJSON(`/explorer/transactions/${txid}`);

        // A transaction with no blockHash is still sitting in the
        // mempool (Chapter 72, Section 5) -- render that distinctly
        // rather than showing a confusing "0 confirmations" as if it
        // were a mined-but-brand-new transaction.
        document.getElementById('status').textContent = tx.blockHash
          ? `Confirmed in block ${linkBlock(tx.blockHash)} (${tx.confirmations} confirmations)`
          : 'Pending in mempool';

        document.getElementById('inputs-body').innerHTML = tx.inputs.map(i => `
          <tr><td>${linkAddress(i.fromAddress)}</td><td>${i.amount}</td></tr>
        `).join('');

        document.getElementById('outputs-body').innerHTML = tx.outputs.map(o => `
          <tr><td>${linkAddress(o.toAddress)}</td><td>${o.amount}</td></tr>
        `).join('');
      } catch (e) {
        showError(e.message);
      }
    }

    loadTx();
  </script>
</body>
</html>
```

The `tx.blockHash ? ... : 'Pending in mempool'` check mirrors, on the frontend, the exact same distinction Chapter 72's `HandleGetTransaction` makes on the backend: a transaction can be real and viewable before it is ever mined. Getting this one conditional right is the difference between an explorer that looks broken the instant someone submits a transaction, and one that correctly shows "pending" the way Etherscan or a block explorer for Bitcoin does.

---

## 7. The Address Page

The address page is the frontend's single busiest page: it shows a balance (a fast, precise number) alongside a scrollable history (a browsable list) — exactly the two things `HandleAddressHistory` bundles into one response in Chapter 72.

```html
<!-- address.html -->
<!DOCTYPE html>
<html>
<head>
  <title>Address — GoChain Explorer</title>
  <link rel="stylesheet" href="style.css">
  <script src="explorer.js"></script>
</head>
<body>
  <h1>Address <span id="address"></span></h1>
  <div id="error" style="display:none; color:red;"></div>
  <p>Balance: <strong id="balance"></strong> gochips</p>
  <h3>Transaction History</h3>
  <table><tbody id="history-body"></tbody></table>
  <div id="pagination"></div>

  <script>
    const params = new URLSearchParams(window.location.search);
    const address = params.get('address');
    const page = parseInt(params.get('page') || '1', 10);
    document.getElementById('address').textContent = shortHash(address);

    async function loadAddress() {
      try {
        const resp = await fetchJSON(
          `/explorer/address/${address}/transactions?page=${page}&pageSize=20`
        );
        document.getElementById('balance').textContent = resp.balance;

        document.getElementById('history-body').innerHTML = resp.data.map(tx => `
          <tr>
            <td>${linkTx(tx.txId)}</td>
            <td>${tx.blockHeight ?? 'pending'}</td>
            <td>${tx.fee}</td>
          </tr>
        `).join('');

        const totalPages = Math.ceil(resp.total / resp.pageSize);
        const pag = document.getElementById('pagination');
        pag.innerHTML = '';
        if (page > 1) pag.innerHTML += `<a href="?address=${address}&page=${page - 1}">Previous</a> `;
        if (page < totalPages) pag.innerHTML += `<a href="?address=${address}&page=${page + 1}">Next</a>`;
      } catch (e) {
        showError(e.message);
      }
    }

    loadAddress();
  </script>
</body>
</html>
```

Notice this page's `fetchJSON` call hits `resp.balance` directly on the top-level object rather than inside `resp.data` — a small but deliberate reminder that `HandleAddressHistory`'s response shape (Chapter 72, Section 6) is *not* a plain `pagedResponse`; it adds `address` and `balance` fields alongside the usual `data`/`page`/`pageSize`/`total` envelope. A frontend that assumed every listing endpoint had an identical shape would silently render `undefined` here — which is exactly the kind of small mismatch worth testing for, as this chapter's exercises ask you to do.

---

## 8. Tying the Pages Together with Links

Every page above already calls `linkBlock`, `linkTx`, or `linkAddress` wherever a hash or address appears in a table cell — that's the entire "navigation" story for this explorer. There is no menu, no search box yet (Chapter 72's `search` endpoint from its own exercises would back one), and no client-side router: every click is a plain `<a href="...">` to a different static HTML file, and every page independently re-fetches whatever it needs from the API. This is intentionally the simplest possible version of "an explorer" — four pages, one shared helper file, and a strict one-to-one mapping to backend endpoints.

```
   index.html  --(click block hash)-->  block.html
       ^                                     |
       |                              (click tx id)
       |                                     v
   address.html <--(click address)--     tx.html
       ^                                     |
       +----------(click input/output address)+
```

Every arrow in that diagram is a plain hyperlink built by one of Section 3's `link*` helpers — there is no hidden machinery. If you can trace this diagram by reading the four HTML files above, you understand the entire frontend.

---

## 9. What a Production Explorer Adds

A real explorer — Etherscan, Blockchain.com's explorer, or a Bitcoin block explorer — starts from exactly this same shape and adds layers of polish this chapter deliberately skipped, because none of them change the *core idea*, only its presentation and robustness:

- **Client-side routing and a build step** (React Router, Vue Router, or similar) so `/block/{hash}` is a real URL instead of `block.html?hash=...`, and so JavaScript is bundled and minified rather than loaded as loose `<script>` tags.
- **A search box** hitting a `search` endpoint (sketched as an exercise back in Chapter 72) that guesses whether the input is a block hash, a transaction ID, or an address.
- **Live updates** — Chapter 71's WebSocket endpoint pushing new blocks onto the homepage the instant they're mined, instead of requiring a manual refresh.
- **Rich formatting** — human-friendly relative timestamps ("3 minutes ago"), QR codes for addresses, charts of network activity over time.
- **Caching and CDN delivery** for pages and assets that rarely change, and the same "don't cache the current tip" caveat Chapter 72's exercises raised for the backend.

None of these change what you already understand: a page calls an endpoint, renders JSON, and links to the next page. Chapter 75's mini project builds a leaner, deployable version of exactly this frontend, bundled directly into a single Go binary via `embed` — proof that "minimal" and "actually shippable" are not in tension.

---

## Summary

- A block explorer frontend is a thin presentation layer over Chapter 72's API — every page maps to exactly one endpoint, with no new blockchain logic anywhere in this chapter.
- Four static HTML pages (homepage, block detail, transaction detail, address) and one shared `explorer.js` file are enough for a fully clickable explorer, with no build step or framework required.
- A shared `fetchJSON` helper centralizes error handling, and `linkBlock`/`linkTx`/`linkAddress` centralize how one page links to another, so the URL shape for each page only needs to be defined once.
- The homepage and address page both consume Chapter 72's `pagedResponse` envelope (`data`, `page`, `pageSize`, `total`) identically, but the address page's response also carries a top-level `balance` field the plain envelope does not — a shape difference worth testing for explicitly.
- The transaction detail page must distinguish a mined transaction (has a `blockHash`, a real confirmation count) from a pending one sitting in the mempool (no `blockHash`, shown distinctly) — mirroring the exact same distinction Chapter 72's backend makes.
- Because a static HTML file has no server-side router, hash and address values are passed as query parameters (`?hash=...`) rather than path segments — a small, honest limitation a real client-side router removes.
- Production explorers add client-side routing, search, live updates via WebSockets, and caching on top of this exact same foundation — none of it changes the core "fetch, render, link" idea this chapter teaches.

---

## Exercises

### Easy

1. Add a fifth static page, `mempool.html`, that lists currently pending transactions using the `GET /explorer/mempool` endpoint from Chapter 72's exercises, reusing `linkTx` for each row.
2. `shortHash` currently always trims to `first 8 + last 6` characters. Add a second parameter controlling how many characters to keep, and use a shorter form for narrow table columns and a longer form for page headings.
3. The homepage's `loadBlocks` function silently does nothing useful if the API is unreachable beyond calling `showError`. Add a "Retry" button that appears alongside the error message and re-calls `loadBlocks()` when clicked.

### Medium

4. Add a small in-browser test (a plain `<script>` you can open as its own HTML file, no test framework required) that calls `fetchJSON` against a mock `fetch` returning a `HandleAddressHistory`-shaped response, and asserts that `resp.balance` is read correctly — the exact shape mismatch Section 7 warned about.
5. Implement a basic client-side "recently viewed" list: every time `block.html`, `tx.html`, or `address.html` loads successfully, store the hash/address and page type into `localStorage`, and render a small "Recently viewed" sidebar on every page from that stored list.
6. Currently, navigating from the address page's "Next" link loses the current scroll position and forces a full page reload. Rewrite `loadAddress` so pagination is handled by re-fetching and re-rendering `#history-body` in place (updating `history.pushState` for the URL) instead of a full navigation, and explain what real routing libraries do to solve this same problem more generally.

### Hard

7. Implement the search box described in Section 9: a single input box on `index.html` that, on submit, calls `GET /explorer/search?q=...` (built as an exercise in Chapter 72) and redirects the browser to `block.html`, `tx.html`, or `address.html` based on the returned `type` field.
8. Wire `index.html` to Chapter 71's WebSocket endpoint (`GET /ws`) so a newly mined block is prepended to the blocks table live, without a page refresh, while leaving the existing `fetch`-based initial load in place as a fallback for browsers or proxies that don't support WebSockets.
9. Rebuild this same four-page explorer using a client-side router of your choice (plain `history.pushState`-based routing is enough; a full framework is not required) so that block, transaction, and address pages use real path segments (`/block/{hash}`) instead of query parameters, and a single `index.html` shell loads all four "pages" as JavaScript-rendered views. Discuss, in a short comment, what you had to add (a router, a not-found view) that a multi-file static site got for free.
