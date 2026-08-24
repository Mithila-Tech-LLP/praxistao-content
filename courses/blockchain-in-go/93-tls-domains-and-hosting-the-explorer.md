# Chapter 93: TLS, Domains, and Hosting the Explorer

Right now, reaching your GoChain explorer means telling someone "go to `http://203.0.113.42:8090`" — a bare IP address, an unfamiliar port, and a browser that may well show a "Not Secure" warning the moment they get there. That is a fine way to demo something to yourself, and a bad way to hand a link to anyone else. This chapter fixes both problems at once: a real domain name in place of the IP address, and a **Caddy** reverse proxy in front of the explorer that gets you free, automatically-renewing HTTPS with almost no configuration. This is the chapter that makes the testnet feel like a real, shareable product instead of a personal experiment.

## Table of Contents

1. [Why a Bare IP and Port Is Not Good Enough](#1-why-a-bare-ip-and-port-is-not-good-enough)
2. [Domains and DNS in Plain Language](#2-domains-and-dns-in-plain-language)
3. [Pointing a Domain at Your VM](#3-pointing-a-domain-at-your-vm)
4. [TLS and HTTPS, Practically](#4-tls-and-https-practically)
5. [What a Reverse Proxy Does](#5-what-a-reverse-proxy-does)
6. [Why Caddy Specifically](#6-why-caddy-specifically)
7. [Writing the Caddyfile](#7-writing-the-caddyfile)
8. [Adding Caddy to Docker Compose](#8-adding-caddy-to-docker-compose)
9. [Opening the Firewall for HTTPS](#9-opening-the-firewall-for-https)
10. [Bringing It Up and Verifying the Certificate](#10-bringing-it-up-and-verifying-the-certificate)
11. [Avoiding Let's Encrypt Rate Limits While Testing](#11-avoiding-lets-encrypt-rate-limits-while-testing)
12. [What This Changes About the Faucet and Grafana](#12-what-this-changes-about-the-faucet-and-grafana)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. Why a Bare IP and Port Is Not Good Enough

Three concrete problems with `http://203.0.113.42:8090`, beyond it simply looking unpolished:

- **It is not memorable or trustworthy.** A number like `203.0.113.42` tells a visitor nothing about who runs it or whether it is safe to click — compare that to `explorer.gochain.example`, which at least states a claim (right or wrong) about ownership.
- **It is plain HTTP, not HTTPS.** Every request and response between a visitor's browser and your explorer travels as unencrypted plain text, readable by anything sitting on the network path between them — a coffee-shop Wi-Fi network, an ISP, anyone. Modern browsers actively flag this with a "Not Secure" warning in the address bar.
- **It will not survive you changing servers.** If you ever move the testnet to a new VM (a bigger one, a different provider), every link anyone has ever shared or bookmarked breaks, because it was pointing at a specific machine's specific address rather than a name you control independently of any one server.

A domain name plus HTTPS fixes all three at once, and — as this chapter shows — costs very little effort to set up correctly.

---

## 2. Domains and DNS in Plain Language

A **domain name** (like `gochain.example`) is a human-readable name that stands in for one or more actual IP addresses. **DNS (Domain Name System)** is the distributed lookup system, briefly introduced back in Chapter 87 for Docker's *internal* version of the same idea, that translates a domain name into the IP address a computer actually needs to open a connection. Think of DNS like a phone book: you look up a business by name, and the phone book hands you back the actual phone number to dial — the name is just a convenient, memorable alias for the number that does the real work.

A **DNS record** is one entry in that phone book. The one this chapter cares about is an **A record** — the simplest kind, mapping a name directly to an IPv4 address:

```
   DNS record type:  A
   Name:             explorer.gochain.example
   Value:            203.0.113.42     <- your VM's public IP from Chapter 88
   TTL:              3600 seconds     <- how long resolvers may cache this answer
```

`TTL` (Time To Live) controls how long other computers' DNS resolvers are allowed to remember this answer before checking again — a short TTL (like the 3600 seconds/one hour shown above) means a future change to the record propagates to the rest of the internet reasonably quickly, at the cost of every resolver having to re-check slightly more often.

---

## 3. Pointing a Domain at Your VM

You need to actually own a domain name first — registered through any domain registrar (Namecheap, Google Domains, Cloudflare, and dozens of others), typically costing somewhere around ten to fifteen dollars a year for a common top-level domain like `.com` or `.dev`. Once registered, every registrar (or a separate DNS provider like Cloudflare, which many people point their domain's nameservers at for its free extra features) gives you a dashboard for adding DNS records exactly like the one in Section 2.

For this chapter's setup, add one A record for the explorer's subdomain, pointing at your Chapter 88 VM's public IP:

```
Type    Name                          Value              TTL
A       explorer.gochain.example      203.0.113.42       3600
```

A **subdomain** (`explorer.` prefixed onto `gochain.example`) is just a further-qualified name under your domain — free to create as many of as you like, each pointing wherever you choose. Add a second A record now, for the Grafana dashboard this chapter also puts behind HTTPS in Section 7:

```
Type    Name                          Value              TTL
A       grafana.gochain.example       203.0.113.42       3600
```

Both records point at the exact same IP — your one VM — because DNS only answers "what address does this name mean," never "what happens once a request arrives there." Routing each name to the correct backend *service* once it reaches that shared IP is entirely Caddy's job, covered starting in Section 5.

A brief note on the **apex domain** (`gochain.example` with no subdomain prefix at all, sometimes called the "root" or "naked" domain): some registrars restrict what record types an apex domain can use (a technical limitation of the DNS root itself, not a Caddy or GoChain concern), which is exactly why this chapter routes everything through subdomains rather than the bare domain — a simpler, more portable choice that every registrar supports without exception.

Confirm the record has propagated before moving on — DNS changes are not instant, and can take anywhere from a few minutes to the full TTL to become visible everywhere:

```bash
dig +short explorer.gochain.example
# 203.0.113.42
```

If `dig` (a standard DNS lookup tool available on macOS and Linux; `nslookup` is the closest Windows equivalent) prints your VM's IP address back, the domain is correctly wired up and ready for the next step.

---

## 4. TLS and HTTPS, Practically

**TLS (Transport Layer Security)** is the protocol that encrypts a connection between a client and a server, so anything traveling between them is unreadable to anyone else on the network path. **HTTPS** is simply "HTTP running over a TLS connection" — the same requests and responses your explorer already handles, just wrapped in an encrypted tunnel first. A padlock icon in a browser's address bar means the browser successfully negotiated TLS with the server and verified its identity.

That identity verification is the part that requires a **TLS certificate** — a small, digitally signed file that says "the operator of this certificate has proven they control `explorer.gochain.example`," issued by a **Certificate Authority (CA)**, an organization browsers already trust to only issue certificates to someone who actually controls the domain in question. **Let's Encrypt** is a free, automated, widely-trusted CA that issues exactly this kind of certificate, and is what makes free, automatic HTTPS practical for a small testnet project rather than something only well-funded companies could previously afford.

The traditional way to get one of these certificates — running `certbot` by hand, copying the resulting certificate and key files into your web server's configuration, and remembering to renew before the certificate expires (Let's Encrypt certificates last 90 days) — works, but is one more manual, easy-to-forget operational chore. This chapter uses a tool that removes that chore entirely.

---

## 5. What a Reverse Proxy Does

A **reverse proxy** is a server that sits in front of one or more backend services, receiving all incoming traffic itself and forwarding each request to the correct backend based on rules you configure — the mirror image of a normal ("forward") proxy, which sits in front of *clients* rather than servers. From a visitor's perspective, they are only ever talking to the reverse proxy; it decides, invisibly, which real service actually handles each request.

```
                    the internet
                         |
                         v
              +----------------------+
              |   Reverse proxy       |     listens on 443 (HTTPS)
              |   (Caddy)              |     terminates TLS here
              +----------------------+
               /          |          \
    explorer.*     grafana.*      (future subdomains)
         |              |
         v              v
  +-------------+  +-------------+
  | gochain-    |  |  grafana    |
  | explorer    |  |  :3000      |
  | :8090       |  |             |
  +-------------+  +-------------+
      (plain HTTP, only reachable
       from inside the Compose
       network - never directly
       from the internet)
```

This is exactly the right place for TLS to live: rather than teaching `gochain-explorer` itself (or Grafana, or the faucet) how to manage certificates, one single reverse proxy terminates HTTPS once, for every backend behind it, and talks plain, unencrypted HTTP to each backend over Docker's own internal, already-private network from Chapter 87 — a service GoChain's own Go code never needs to know or care exists.

---

## 6. Why Caddy Specifically

Several reverse proxies could do this job — **Nginx** and **Traefik** are both common, capable choices. This chapter uses **Caddy** specifically because it makes automatic HTTPS the *default* behavior rather than an add-on you configure: point Caddy at a domain name, and it automatically requests a Let's Encrypt certificate, serves HTTPS immediately, and silently renews the certificate well before its 90-day expiry, with zero certificate-specific configuration written by you at all. Where Nginx or Traefik can absolutely be configured to do the same thing, Caddy's central design idea is that this should be the easy, obvious path rather than an advanced feature to discover later.

---

## 7. Writing the Caddyfile

Caddy's configuration file is conventionally named `Caddyfile`, and its whole design philosophy shows immediately in how little it takes to configure a working HTTPS reverse proxy:

```
# Caddyfile — reverse proxy with automatic HTTPS for the GoChain
# explorer and Grafana dashboard.

explorer.gochain.example {
    # Everything under this block applies only to requests whose Host
    # header matches this domain. Caddy automatically requests and
    # renews a Let's Encrypt certificate for exactly this name the
    # first time it starts serving it - no certbot, no manual renewal
    # cron job, nothing else to configure.
    reverse_proxy gochain-explorer:8090

    # A small, free extra: gzip/zstd compress responses where the
    # client supports it, shrinking the explorer's HTML/JS payload.
    encode gzip
}

grafana.gochain.example {
    reverse_proxy grafana:3000
}
```

That is the entire file. `explorer.gochain.example { reverse_proxy gochain-explorer:8090 }` reads almost like plain English: "for requests to this domain, forward them to `gochain-explorer` on port 8090" — and `gochain-explorer` here is, once again, Docker Compose's internal service-name DNS from Chapter 87, resolved from inside the same `gochain-net` network Caddy will join. Caddy listens on the standard HTTPS port 443 (and plain port 80, which it uses briefly during the certificate-issuance process and then automatically redirects to HTTPS) — nothing in this file mentions certificates, keys, or renewal timers, because Caddy's whole point is handling all of that invisibly.

---

## 8. Adding Caddy to Docker Compose

```yaml
# docker-compose.yml — addition for Chapter 93
services:
  # ... existing node1, node2, node3, explorer, faucet, prometheus,
  # grafana from earlier chapters ...

  caddy:
    image: caddy:2.8-alpine
    container_name: gochain-caddy
    ports:
      - "80:80"   # required briefly for Let's Encrypt's HTTP-01 challenge
      - "443:443" # the real, ongoing HTTPS traffic
    volumes:
      - ./Caddyfile:/etc/caddyfile/Caddyfile:ro
      - caddy-data:/data       # persists issued certificates across restarts
      - caddy-config:/config
    command: ["caddy", "run", "--config", "/etc/caddyfile/Caddyfile"]
    networks:
      - gochain-net
    depends_on:
      - explorer
      - grafana

volumes:
  caddy-data:
  caddy-config:
```

The `caddy-data` volume matters more than it looks: it is where Caddy stores the certificates it obtains, so a container restart (or an image update) does not throw away a valid certificate and force Caddy to request a brand-new one from Let's Encrypt unnecessarily — Let's Encrypt does enforce rate limits on how often the same domain can request new certificates, so preserving this volume is not just an optimization but a real reliability concern.

---

## 9. Opening the Firewall for HTTPS

Recall Chapter 88's `ufw` rules opened only 22 (SSH), 8080 (API), and 9000 (P2P). Caddy needs two more, standard web ports opened before any of this works from the outside:

```bash
sudo ufw allow 80/tcp   # required for Let's Encrypt's certificate challenge
sudo ufw allow 443/tcp  # the actual HTTPS traffic visitors will use
sudo ufw status verbose
```

This is exactly the kind of deliberate, narrow port-opening Chapter 88 anticipated — the explorer's raw port 8090 and Grafana's raw port 3000 stay closed to the public internet entirely; only Caddy's 80 and 443 are opened, and Caddy alone decides, based on the `Caddyfile`, what plain-HTTP backend each encrypted request actually reaches.

---

## 10. Bringing It Up and Verifying the Certificate

```bash
docker compose up -d caddy
docker compose logs -f caddy
```

Expected log output on first start, abbreviated:

```
gochain-caddy | {"level":"info","msg":"serving initial configuration"}
gochain-caddy | {"level":"info","msg":"obtaining certificate","identifier":"explorer.gochain.example"}
gochain-caddy | {"level":"info","msg":"certificate obtained successfully","identifier":"explorer.gochain.example"}
gochain-caddy | {"level":"info","msg":"obtaining certificate","identifier":"grafana.gochain.example"}
gochain-caddy | {"level":"info","msg":"certificate obtained successfully","identifier":"grafana.gochain.example"}
```

From any browser, visit `https://explorer.gochain.example` — the padlock icon should appear immediately, with no warning. From the command line, confirm both the redirect and the certificate details:

```bash
curl -I http://explorer.gochain.example
# HTTP/1.1 308 Permanent Redirect
# Location: https://explorer.gochain.example/

curl -vI https://explorer.gochain.example 2>&1 | grep -E "subject:|issuer:"
# subject: CN=explorer.gochain.example
# issuer: C=US, O=Let's Encrypt, CN=R11
```

The plain-HTTP request being redirected (rather than served directly) confirms Caddy is enforcing HTTPS rather than merely offering it, and the certificate's `issuer` line confirms it really was issued by Let's Encrypt, not a self-signed placeholder a browser would reject.

---

## 11. Avoiding Let's Encrypt Rate Limits While Testing

Let's Encrypt enforces a real limit on how many certificates it will issue for the same domain in a rolling week — generous for normal use, but easy to accidentally hit while you are still experimenting with your `Caddyfile`, restarting the `caddy` container repeatedly, tweaking a typo, and triggering a fresh certificate request each time. Once you hit that limit, every further request for that domain is refused for the rest of the week, no exceptions — an entirely avoidable, frustrating way to lose testing time.

Let's Encrypt's solution is a separate **staging environment**: a parallel certificate authority, using the exact same protocol, that issues certificates browsers do *not* trust (so you cannot use them for a real demo) but that has dramatically more generous rate limits, specifically meant for exactly this kind of iterative configuration testing. Caddy makes switching to it a one-line addition near the top of the `Caddyfile`:

```
# Caddyfile — add this block while iterating on your configuration,
# then delete it once you are confident the config is correct and are
# ready for a real, browser-trusted certificate.
{
    acme_ca https://acme-staging-v02.api.letsencrypt.org/directory
}

explorer.gochain.example {
    reverse_proxy gochain-explorer:8090
    encode gzip
}
```

With that block present, Caddy requests every certificate from the staging CA instead — your browser will show a certificate warning (expected and correct, since staging certificates are deliberately untrusted), but you can restart, break, and fix your `Caddyfile` as many times as you need without any risk of hitting the production rate limit. Once the configuration is confirmed working end to end, delete the `acme_ca` override entirely and restart Caddy once, to request the real, trusted certificate exactly one time.

---

## 12. What This Changes About the Faucet and Grafana

Chapter 90 left the faucet reachable only inside the Compose network (port 8091, never opened in `ufw`), and Chapter 91 left Grafana and Prometheus reachable only via an SSH tunnel to ports 3000 and 9090. With Caddy now in place, extending real, HTTPS-secured public access to any of them is a matter of adding one more block to the `Caddyfile` — no new firewall rule, since Caddy already owns the only two ports (80/443) that matter for this purpose:

```
faucet.gochain.example {
    reverse_proxy gochain-faucet:8091
}
```

Grafana was already given a subdomain in Section 7; Prometheus itself deliberately is not, since its raw query interface exposes more operational detail than is worth handing to the public internet — Grafana, which queries Prometheus on your behalf and presents curated dashboards, remains the intended public-facing surface for that data, exactly the distinction Chapter 91 drew when it left Prometheus's own port 9090 tunnel-only.

This is the general pattern this chapter leaves you with: every new HTTP-speaking service this course adds from here on gets a subdomain and a `reverse_proxy` block, never a newly opened raw port — Caddy is now the one and only front door.

---

## Summary

- A bare IP address and port is unmemorable, untrusted-looking, unencrypted, and tied to one specific server — a domain plus HTTPS fixes all four.
- DNS translates a human-readable domain name into an IP address via records; an **A record** is the simple case, mapping a name directly to your VM's IPv4 address.
- TLS encrypts a connection; HTTPS is HTTP running over TLS; a certificate, issued by a trusted Certificate Authority like Let's Encrypt, proves you actually control the domain the browser is connecting to.
- A **reverse proxy** sits in front of your backend services, terminating HTTPS once and forwarding plain HTTP to each backend over Docker's already-private internal network.
- **Caddy** is used specifically because automatic HTTPS is its default behavior — a `Caddyfile` block naming a domain and a backend is the entire configuration needed, with certificate issuance and renewal handled invisibly.
- Caddy's data volume must persist across restarts, both so certificates survive container recreation and to avoid unnecessarily hitting Let's Encrypt's rate limits.
- Let's Encrypt's **staging environment** (`acme_ca` pointed at its staging URL) issues untrusted-but-rate-limit-friendly certificates for iterating on a `Caddyfile` safely, before switching to production for the one real certificate you actually need.
- Only ports 80 and 443 need to be newly opened in `ufw` — every backend behind Caddy (explorer, faucet, Grafana) keeps its own raw port closed to the public internet entirely.
- From this chapter forward, every new public-facing service gets a subdomain and a `Caddyfile` block, never a new raw port opened directly.

---

## Exercises

### Easy

1. Register a domain (or use a free subdomain service if you do not want to purchase one) and add an A record pointing a subdomain at your Chapter 88 VM's IP address. Confirm it resolves correctly with `dig +short <your-subdomain>`.

2. Write a `Caddyfile` with a single block reverse-proxying your subdomain to `gochain-explorer:8090`, add the `caddy` service to `docker-compose.yml`, and bring it up. Visit your domain in a browser and confirm the padlock icon appears with no warnings.

3. Run `curl -I http://<your-subdomain>` and confirm you get a redirect response to the `https://` version, rather than a plain 200 response served over HTTP.

### Medium

4. Add a second `Caddyfile` block for a `faucet.` subdomain, following Section 12's pattern, and confirm you can submit a faucet request through the new HTTPS URL instead of the old internal-only port.

5. Add the staging `acme_ca` override from Section 11 to your `Caddyfile`, restart Caddy, and confirm your browser now shows an untrusted-certificate warning when visiting your domain. Then remove the override, restart once more, and confirm the warning disappears and a real, trusted certificate is issued.

6. Stop and remove the `caddy` container (`docker compose rm -f caddy`), then bring it back up, and confirm — by checking the logs — that it reuses the certificate already stored in the `caddy-data` volume rather than requesting a brand-new one from Let's Encrypt. Explain, in 100 words, why this matters given Let's Encrypt's rate limits.

7. Add HTTP Basic Authentication to the Grafana subdomain's `Caddyfile` block (Caddy supports this natively with a `basicauth` directive) so Grafana's dashboards require a username and password even though they are now technically internet-reachable. Test that an incorrect password is rejected and a correct one succeeds.

### Hard

8. Configure Caddy to serve a custom error page for the explorer subdomain when the backend `gochain-explorer` container is down (simulate this by stopping it), instead of Caddy's default error response. Explain what a real production deployment would want this error page to say to a visitor versus what it would want logged internally.

9. Investigate Caddy's on-demand TLS feature, which can issue certificates for domains not listed in the `Caddyfile` at all, decided dynamically at request time. Explain, in 200-300 words, why this feature would be dangerous to enable for a public-facing GoChain deployment without an additional authorization check, and what that check should verify before Caddy is allowed to request a certificate for an arbitrary incoming domain name.

10. Set up a second, independent Caddy instance (or a second `Caddyfile` block) that reverse-proxies to a *second* GoChain node's explorer instance on a different subdomain, and add HTTP-to-HTTPS redirect testing plus a TLS certificate expiry monitoring check (using a tool like `openssl s_client` scripted on a cron job) that alerts you if a certificate's remaining validity ever drops below 14 days — a real safety net in case Caddy's automatic renewal ever silently fails.
