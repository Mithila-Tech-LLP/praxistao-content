# Chapter 02: How the Internet Works — Packets, Protocols, and the Big Picture

*The internet is the primary attack surface for almost every hack. Before you can attack or defend networks, you need to understand how they work. This chapter covers everything from how your data physically travels to how protocols create the abstraction of a reliable connection.*

---

## The Big Picture

When you open Instagram on your phone in Mumbai:

1. Your phone sends a request to Instagram's servers in the US
2. That request travels through dozens of routers across the world
3. Instagram's server receives it, processes it, sends back your feed
4. Your phone receives the response and displays it

The data didn't travel as one piece. It was broken into small chunks called **packets**, each taking potentially different paths across the internet, and reassembled at the destination.

---

## The OSI Model — Seven Layers of Abstraction

Networking is complex, so engineers organized it into layers. Each layer handles a specific concern and talks only to the layers immediately above and below it.

```
┌─────────────────────────────────────────────────────┐
│ Layer 7: Application   │ HTTP, DNS, SMTP, FTP        │
├─────────────────────────────────────────────────────┤
│ Layer 6: Presentation  │ Encryption, encoding (TLS)  │
├─────────────────────────────────────────────────────┤
│ Layer 5: Session       │ Session management          │
├─────────────────────────────────────────────────────┤
│ Layer 4: Transport     │ TCP, UDP (ports)             │
├─────────────────────────────────────────────────────┤
│ Layer 3: Network       │ IP (routing between networks)│
├─────────────────────────────────────────────────────┤
│ Layer 2: Data Link     │ Ethernet, WiFi (MAC address) │
├─────────────────────────────────────────────────────┤
│ Layer 1: Physical      │ Cables, radio waves, light   │
└─────────────────────────────────────────────────────┘
```

**How data flows through layers:**
When you send data, it passes DOWN through layers, each adding its own header (metadata). When data is received, it passes UP, each layer removing its header.

This "wrapping" process is called **encapsulation**:
```
[Application Data]
[TCP Header][Application Data]
[IP Header][TCP Header][Application Data]
[Ethernet Header][IP Header][TCP Header][Application Data][Ethernet Trailer]
```

**Why this matters for security:** Attacks happen at specific layers:
- Layer 1-2: WiFi deauthentication, MAC spoofing
- Layer 3: IP spoofing, ICMP attacks
- Layer 4: TCP SYN flood, port scanning
- Layer 7: Web application attacks (SQL injection, XSS)

Security tools must understand which layer they're operating at.

---

## Layer 3: IP — Internet Protocol

IP is responsible for routing packets from source to destination across multiple networks.

**IP Address:** A 32-bit number (IPv4) that identifies a device on a network.
- Written as four decimal octets: `192.168.1.100`
- Two parts: network part and host part (determined by subnet mask)

**IPv4 vs IPv6:**
- IPv4: 32-bit addresses (~4 billion possible) — running out
- IPv6: 128-bit addresses (~340 undecillion possible) — the future

**Public vs Private IPs:**
Private IP ranges (used inside home/office networks):
- `10.0.0.0/8` (10.x.x.x)
- `172.16.0.0/12` (172.16-31.x.x)
- `192.168.0.0/16` (192.168.x.x)

Private IPs aren't routable on the public internet. NAT (Network Address Translation) converts private IPs to your router's public IP when accessing the internet.

**The IP Header** (simplified):
```
Bits:  0     4     8     16    24    31
       ┌─────┬─────┬─────────────────┐
       │ Ver │ IHL │  Type of Service │
       ├─────────────────────────────┤
       │ Total Length                 │
       ├─────────────────────────────┤
       │ ID          │ Flags │ Offset │
       ├─────────────────────────────┤
       │ TTL         │ Protocol       │
       ├─────────────────────────────┤
       │ Source IP Address            │
       ├─────────────────────────────┤
       │ Destination IP Address       │
       └─────────────────────────────┘
```

**TTL (Time to Live):** Every router that forwards a packet decrements TTL by 1. When TTL reaches 0, the packet is dropped. This prevents packets from circulating forever. Starting TTL: 64 (Linux), 128 (Windows), 255 (some devices). By observing TTL, you can guess what OS a device is running — this is called **OS fingerprinting**.

**Protocol field:** Tells you what's inside the IP packet:
- 6 = TCP
- 17 = UDP
- 1 = ICMP (ping)

---

## Layer 4: TCP — Transmission Control Protocol

TCP provides **reliable, ordered delivery** of data. It's used by HTTP, HTTPS, SSH, email, and most protocols that need all data to arrive correctly.

### The TCP Three-Way Handshake

Before any data is sent, TCP establishes a connection:

```
Client                    Server
  │                          │
  │──────── SYN ─────────────▶│   "I want to connect, my seq=1000"
  │                          │
  │◀────── SYN-ACK ──────────│   "OK, my seq=2000, ack=1001"
  │                          │
  │──────── ACK ─────────────▶│   "Acknowledged, ack=2001"
  │                          │
  │         [Connected]       │
```

SYN = Synchronize, ACK = Acknowledge

**Why this matters for security:**
- **Port scanning** works by sending SYN packets and seeing if you get SYN-ACK back (port is open) or RST (port is closed)
- **SYN flood attack:** Attacker sends many SYN packets but never completes the handshake, filling the server's connection table
- **Stateful firewalls** track connection state to allow only packets that are part of established connections

### TCP Flags

Each TCP packet has flags that indicate its purpose:
- **SYN:** Initiate connection
- **ACK:** Acknowledge received data
- **FIN:** Terminate connection gracefully
- **RST:** Reset/abort connection
- **PSH:** Push data to application immediately
- **URG:** Urgent data

An **ACK scan** (Nmap) sends packets with only ACK set. Firewalls that only block new connections will pass these through, revealing firewall rules.

### TCP Ports

Ports allow multiple services on one IP address. Think of an IP address as a building address and ports as apartment numbers.

Range: 0-65535

- **Ports 0-1023:** Well-known ports (require root/admin to bind)
- **Ports 1024-49151:** Registered ports
- **Ports 49152-65535:** Dynamic/ephemeral ports (used for client-side connections)

**Common ports to know:**

| Port | Protocol | Service |
|------|----------|---------|
| 21 | TCP | FTP (file transfer) |
| 22 | TCP | SSH (secure shell) |
| 23 | TCP | Telnet (insecure remote shell) |
| 25 | TCP | SMTP (email sending) |
| 53 | TCP/UDP | DNS (name resolution) |
| 80 | TCP | HTTP (web) |
| 443 | TCP | HTTPS (secure web) |
| 445 | TCP | SMB (Windows file sharing) |
| 3306 | TCP | MySQL database |
| 3389 | TCP | RDP (Windows remote desktop) |
| 8080 | TCP | HTTP alternate / proxy |

**Penetration testing focus:** Ports 22, 23, 80, 443, 445, 3389 are almost always the first things checked. An open port 23 (Telnet) on a production server is a critical finding.

---

## Layer 4: UDP — User Datagram Protocol

UDP is **fast but unreliable**. No handshake, no guaranteed delivery, no ordering.

Used for: DNS, DHCP, video streaming, VoIP, online gaming.

The tradeoff: UDP is faster (no handshake overhead, no retransmission), but you might lose packets.

**Security note:** UDP services are harder to scan and often overlooked. Many forgotten services run over UDP.

---

## Layer 7: DNS — The Internet's Phone Book

DNS (Domain Name System) converts human-readable names to IP addresses.

When you type `google.com`:
1. Your browser asks your OS: "What's the IP for google.com?"
2. OS checks its cache. If not there, asks your DNS resolver (usually your ISP or 8.8.8.8)
3. DNS resolver asks the root DNS servers: "Who handles .com?"
4. Root server: "Ask the .com nameservers"
5. .com nameserver: "Ask Google's nameservers"
6. Google's nameserver: "google.com is `142.250.195.78`"

**DNS record types:**
- **A record:** Domain → IPv4 address
- **AAAA record:** Domain → IPv6 address
- **MX record:** Mail servers for this domain
- **CNAME:** Alias (one domain points to another)
- **TXT:** Text records (used for SPF, DKIM email authentication)
- **NS:** Name servers for this domain
- **PTR:** Reverse lookup (IP → domain)

**Why DNS matters for security:**
- **DNS spoofing/poisoning:** Attacker injects false DNS records, redirecting traffic
- **DNS zone transfer:** If misconfigured, reveals all subdomains (AXFR request)
- **DNS tunneling:** Attackers exfiltrate data encoded in DNS queries
- **Subdomain enumeration:** Finding `admin.example.com`, `dev.example.com`, etc.

### DNS Lookup in Go

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Forward lookup: domain → IPs
    ips, err := net.LookupHost("google.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("google.com resolves to:")
    for _, ip := range ips {
        fmt.Println(" ", ip)
    }
    
    // MX record lookup
    mxRecords, _ := net.LookupMX("google.com")
    fmt.Println("\nMail servers:")
    for _, mx := range mxRecords {
        fmt.Printf("  %s (priority: %d)\n", mx.Host, mx.Pref)
    }
    
    // Reverse lookup: IP → domain
    names, _ := net.LookupAddr("8.8.8.8")
    fmt.Println("\n8.8.8.8 reverse lookup:", names)
}
```

---

## Layer 7: HTTP — How the Web Works

HTTP (HyperText Transfer Protocol) is the language browsers and servers use.

**Request format:**
```
GET /search?q=hello HTTP/1.1
Host: google.com
User-Agent: Mozilla/5.0
Accept: text/html
Cookie: session=abc123
                        ← blank line = end of headers
```

**Response format:**
```
HTTP/1.1 200 OK
Content-Type: text/html
Content-Length: 12345
Set-Cookie: session=xyz789

<html>...</html>         ← response body
```

**HTTP methods:**
- `GET`: Retrieve data
- `POST`: Submit data (create/update)
- `PUT`: Replace data
- `DELETE`: Remove data
- `OPTIONS`: What methods are allowed?
- `HEAD`: Like GET but no body (check if resource exists)

**HTTP status codes:**
- `200 OK`: Success
- `301/302`: Redirect
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Authenticated but not allowed
- `404 Not Found`: Resource doesn't exist
- `500 Server Error`: Server-side bug

**Why HTTP matters for security:** Most web attacks manipulate HTTP. SQL injection is in HTTP parameters. XSS is in HTTP responses. CSRF exploits HTTP cookies. Understanding HTTP is the foundation of web application security.

---

## HTTPS — Encrypting HTTP

HTTPS = HTTP over TLS (Transport Layer Security).

Before any HTTP data is exchanged, a TLS handshake establishes an encrypted channel:
1. Client says "Hello" with list of supported cipher suites
2. Server responds with its certificate (proof of identity)
3. Client verifies the certificate against trusted Certificate Authorities (CAs)
4. Keys are exchanged (using asymmetric crypto)
5. All subsequent communication is encrypted with symmetric crypto

**What HTTPS protects against:** Eavesdropping (man-in-the-middle seeing your data) and tampering (someone modifying data in transit).

**What HTTPS does NOT protect:** Once the data reaches the server, it's decrypted. SQL injection, XSS, and other application-layer attacks still work over HTTPS.

---

## Packets in Practice — What Network Traffic Looks Like

A single HTTP GET request to `http://example.com/index.html` generates:

1. DNS query/response (UDP, port 53)
2. TCP SYN, SYN-ACK, ACK (three-way handshake)
3. HTTP GET request (TCP payload)
4. HTTP 200 OK response (TCP payload, possibly multiple packets)
5. TCP FIN/ACK (connection teardown)

**Wireshark** can capture all of this. You'll build a Go-based packet sniffer in Chapter 20.

---

## Summary

```
Your Browser              Your Router            Web Server
     │                        │                      │
     │──DNS query(UDP 53)─────▶│──forward──────────▶│
     │◀─DNS response──────────│◀──response──────────│
     │                        │                      │
     │──TCP SYN (port 443)────▶│──route──────────────▶│
     │◀─TCP SYN-ACK───────────│◀────────────────────│
     │──TCP ACK───────────────▶│──route──────────────▶│
     │                        │                      │
     │──TLS Handshake─────────▶│──────────────────────▶│
     │◀─TLS Established───────│◀────────────────────│
     │                        │                      │
     │──HTTPS GET /──────────▶│──────────────────────▶│
     │◀─HTTPS 200 OK──────────│◀────────────────────│
```

---

## Exercises

1. Open a terminal and run `ping google.com`. What does TTL tell you about the number of routers between you and Google?
2. Run `nslookup flipkart.com`. How many IPs does it have? What does this tell you?
3. Open your browser's Developer Tools (F12) → Network tab. Load a website. Find a request and look at its headers. What interesting headers do you see?
4. Write a Go program that does a DNS lookup for 5 Indian websites and prints their IP addresses.
5. What's the difference between `192.168.1.0/24` and `10.0.0.0/8` in terms of number of hosts?
