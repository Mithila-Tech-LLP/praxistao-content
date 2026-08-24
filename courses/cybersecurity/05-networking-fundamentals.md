# Chapter 05: Networking Fundamentals — TCP/IP, UDP, and How Data Travels

*Networks are the battlefield for most cyberattacks. Understanding how packets travel is the foundation of everything — attack and defense.*

---

## The Big Picture: How Data Gets from A to B

When you type `google.com` in your browser:

1. **DNS** — converts `google.com` → `142.250.185.46` (an IP address)
2. **TCP handshake** — your computer connects to Google's server
3. **TLS** — encrypts the connection (HTTPS)
4. **HTTP request** — browser asks for the webpage
5. **HTTP response** — Google sends HTML, CSS, JavaScript
6. **Browser renders** — you see the page

Every attack targets some part of this chain.

---

## IP Addresses

Every device on a network has an IP address — a unique identifier.

**IPv4:** 32-bit number, written as 4 groups of 0-255
```
192.168.1.100
```

In binary: each group is 8 bits (an octet)
```
192      168      1        100
11000000 10101000 00000001 01100100
```

**IPv6:** 128-bit, handles the address shortage
```
2001:0db8:85a3:0000:0000:8a2e:0370:7334
```

### Private vs Public IP

| Range | Type | Used for |
|-------|------|----------|
| `10.0.0.0/8` | Private | Large networks |
| `172.16.0.0/12` | Private | Medium networks |
| `192.168.0.0/16` | Private | Home routers |
| Everything else | Public | Internet |

Private addresses never appear on the public internet. Your home router gets one public IP; inside your home all devices share it (via NAT — Network Address Translation).

### CIDR Notation

`192.168.1.0/24` means:
- Network portion: first 24 bits (`192.168.1`)
- Host portion: last 8 bits (0-255)
- **254 usable addresses**: `192.168.1.1` to `192.168.1.254`

As a pentester, `/24` tells you "there are 254 targets to scan in this network."

---

## TCP — The Reliable Protocol

TCP (Transmission Control Protocol) guarantees delivery. Every packet is acknowledged. If a packet is lost, it's retransmitted.

### TCP Three-Way Handshake

```
Client                          Server
  |                               |
  |------- SYN ------------------>|    "I want to connect"
  |                               |
  |<------ SYN-ACK ---------------|    "OK, I'm ready"
  |                               |
  |------- ACK ------------------>|    "Great, connection open"
  |                               |
  |===== DATA FLOWS BOTH WAYS =====|
```

**SYN scan (Nmap's -sS):** Send SYN, don't complete handshake
- SYN-ACK response = port is open
- RST response = port is closed
- No response = firewall dropped the packet

### TCP Header Fields

```
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|   Source Port     |   Destination Port           |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|             Sequence Number                      |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|           Acknowledgment Number                  |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|Data | | |U|A|P|R|S|F|     Window Size            |
|Off  |Res|Res|R|C|S|S|Y|I|                         |
|set  |   |   |G|K|H|T|N|N|                         |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|         Checksum        |    Urgent Pointer       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

**Flags attackers care about:**
- **SYN (S):** Initiate connection
- **ACK (A):** Acknowledge received data
- **FIN (F):** Gracefully close connection
- **RST (R):** Abruptly close connection (port closed, or firewall)
- **PSH (P):** Push data to application immediately
- **URG (U):** Urgent data pointer is valid

**ACK scan:** Send ACK to a port. RST = port exists (firewall is stateless). No response = stateful firewall blocks it. Used for firewall fingerprinting.

---

## UDP — The Fast, Unreliable Protocol

UDP (User Datagram Protocol) fires and forgets. No handshake. No acknowledgment.

```
Client                          Server
  |                               |
  |------- UDP packet ----------->|    "Here's your data"
  |                               |    (no response required)
```

**Used for:** DNS, video streaming, VoIP, gaming, NTP

**Why it matters for security:**
- UDP services often get less attention from admins → more likely to have vulnerabilities
- DNS amplification attacks use UDP — send small request, get giant response, flood victim
- Slower to scan (no handshake to time out on)

---

## ICMP — The Network Health Protocol

ICMP (Internet Control Message Protocol) is used for network diagnostics.

```bash
ping 8.8.8.8          # sends ICMP echo request, expects echo reply
traceroute 8.8.8.8    # uses ICMP TTL exceeded messages
```

**ICMP attacks:**
- **Ping of Death:** Oversized ICMP packet to crash old systems (patched)
- **ICMP flood (ping flood):** DDoS attack — overwhelming pings
- **Smurf attack:** Spoof victim's IP, ping broadcast address → entire network pings victim

Many firewalls block ICMP. `ping` not working doesn't mean the host is down.

---

## Ports — Application Addressing

IP gets you to the right machine. The **port number** tells the machine which application should handle the packet.

Ports 0-65535:
- **0-1023:** Well-known ports (root required to listen)
- **1024-49151:** Registered ports
- **49152-65535:** Ephemeral/dynamic (used by clients)

**Critical ports to know:**

| Port | Protocol | Notes |
|------|----------|-------|
| 21 | FTP | File transfer, often unencrypted |
| 22 | SSH | Encrypted remote shell |
| 23 | Telnet | Unencrypted remote shell (dangerous!) |
| 25 | SMTP | Email sending |
| 53 | DNS | TCP and UDP |
| 80 | HTTP | Unencrypted web |
| 110 | POP3 | Email retrieval |
| 143 | IMAP | Email |
| 443 | HTTPS | Encrypted web |
| 445 | SMB | Windows file sharing (EternalBlue!) |
| 1433 | MSSQL | Microsoft SQL Server |
| 3306 | MySQL | MySQL database |
| 3389 | RDP | Windows Remote Desktop |
| 5432 | PostgreSQL | PostgreSQL database |
| 5900 | VNC | Remote desktop |
| 6379 | Redis | Cache database (often unauthenticated!) |
| 8080 | HTTP-alt | Web servers, proxies |
| 27017 | MongoDB | MongoDB (often unauthenticated!) |

---

## DNS — The Internet's Phone Book

When you type `google.com`, your computer doesn't know the IP. DNS translates names to addresses.

### DNS Query Process

```
Browser
  |
  +-- Cache miss
  |
  +-> Recursive Resolver (your ISP or 8.8.8.8)
        |
        +-> Root nameserver (.com)
              |
              +-> TLD nameserver (google.com NS)
                    |
                    +-> Authoritative nameserver
                          |
                          +-> Returns: 142.250.185.46
```

### DNS Record Types

| Record | Purpose | Example |
|--------|---------|---------|
| **A** | IPv4 address | `google.com` → `142.250.185.46` |
| **AAAA** | IPv6 address | `google.com` → `2607:f8b0:4004:c08::8b` |
| **CNAME** | Alias | `www.google.com` → `google.com` |
| **MX** | Mail server | `google.com` → `smtp.google.com` |
| **TXT** | Text (SPF, DKIM) | `"v=spf1 include:_spf.google.com ~all"` |
| **NS** | Nameserver | `google.com` NS → `ns1.google.com` |
| **PTR** | Reverse DNS | `142.250.185.46` → `lga25s57-in-f14.1e100.net` |
| **SOA** | Zone authority | Zone metadata |

### DNS for Attackers and Defenders

**Attackers:**
```bash
# Zone transfer — dump all DNS records (often misconfigured)
dig axfr @ns1.target.com target.com

# Subdomain enumeration
for sub in www mail dev api staging admin; do
    dig +short $sub.target.com
done

# Find mail servers (target for phishing)
dig MX target.com

# SPF record — understand email setup
dig TXT target.com | grep spf
```

**Defenders:**
- Zone transfers should be restricted to authorized nameservers only
- Monitor DNS for fast-flux domains (malware C2)
- DNS over HTTPS (DoH) hides DNS queries from network monitors

---

## TCP/IP Stack — The Layers

```
Application Layer    HTTP, HTTPS, SSH, DNS, SMTP
                     ↕ (data / payload)
Transport Layer      TCP, UDP
                     ↕ (segments)
Internet Layer       IP, ICMP
                     ↕ (packets)
Network Access       Ethernet, WiFi, ARP
                     ↕ (frames)
Physical             Cables, radio waves, light
```

**Encapsulation:** Each layer wraps the data from the layer above:

```
[Ethernet header][IP header][TCP header][HTTP data][Ethernet trailer]
```

**Why this matters:**
- **ARP poisoning** attacks the Network Access layer
- **IP spoofing** attacks the Internet layer
- **TCP session hijacking** attacks the Transport layer
- **HTTP injection** attacks the Application layer

---

## ARP — Address Resolution Protocol

Within a LAN, communication uses MAC addresses (hardware addresses), not IPs.

ARP maps IP addresses to MAC addresses.

```
Host A wonders: "Who has 192.168.1.1?"
Broadcasts on LAN:  "ARP Request: Who has 192.168.1.1?"
Router responds:    "ARP Reply: 192.168.1.1 is at aa:bb:cc:dd:ee:ff"
```

**ARP Poisoning / ARP Spoofing:**
```
Attacker sends:     "192.168.1.1 is at my MAC address (lie!)"
All hosts update their ARP cache with the lie
All traffic for 192.168.1.1 now goes to attacker first
= Man-in-the-Middle attack
```

Tool: `arpspoof` (dsniff), `ettercap`, or custom Go tool

---

## Firewalls and NAT

**Stateful firewall:** Tracks connection state. Allows established connections; blocks unexpected inbound.

**Stateless firewall (ACL):** Simple rule matching. No state tracking. Easier to bypass.

**NAT (Network Address Translation):**
```
Your computer: 192.168.1.100:54321 → 8.8.8.8:443
Router translates: 203.0.113.5:54321 → 8.8.8.8:443
Response returns to router, translated back to your IP
```

**Why NAT matters:** Millions of private IPs hide behind one public IP. Attacker can't reach your machine directly — they need you to initiate contact (phishing, malicious websites, etc.)

---

## Go: Building a Network Sniffer Concept

```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
)

// Parse a raw IP header from bytes
type IPHeader struct {
    Version     uint8
    IHL         uint8
    TOS         uint8
    TotalLength uint16
    TTL         uint8
    Protocol    uint8
    SrcIP       net.IP
    DstIP       net.IP
}

func parseIPHeader(data []byte) (*IPHeader, error) {
    if len(data) < 20 {
        return nil, fmt.Errorf("too short")
    }
    h := &IPHeader{
        Version:     data[0] >> 4,
        IHL:         (data[0] & 0x0F) * 4,
        TOS:         data[1],
        TotalLength: binary.BigEndian.Uint16(data[2:4]),
        TTL:         data[8],
        Protocol:    data[9],
        SrcIP:       net.IP(data[12:16]),
        DstIP:       net.IP(data[16:20]),
    }
    return h, nil
}

func protocolName(proto uint8) string {
    switch proto {
    case 1:  return "ICMP"
    case 6:  return "TCP"
    case 17: return "UDP"
    default: return fmt.Sprintf("Unknown(%d)", proto)
    }
}

// Parse TCP header (starts after IP header)
type TCPHeader struct {
    SrcPort  uint16
    DstPort  uint16
    Seq      uint32
    Ack      uint32
    Flags    uint8
    Window   uint16
}

func parseTCPHeader(data []byte) (*TCPHeader, error) {
    if len(data) < 20 {
        return nil, fmt.Errorf("too short")
    }
    return &TCPHeader{
        SrcPort: binary.BigEndian.Uint16(data[0:2]),
        DstPort: binary.BigEndian.Uint16(data[2:4]),
        Seq:     binary.BigEndian.Uint32(data[4:8]),
        Ack:     binary.BigEndian.Uint32(data[8:12]),
        Flags:   data[13],
        Window:  binary.BigEndian.Uint16(data[14:16]),
    }, nil
}

func (t *TCPHeader) FlagsString() string {
    flags := ""
    if t.Flags&0x02 != 0 { flags += "SYN " }
    if t.Flags&0x10 != 0 { flags += "ACK " }
    if t.Flags&0x01 != 0 { flags += "FIN " }
    if t.Flags&0x04 != 0 { flags += "RST " }
    if t.Flags&0x08 != 0 { flags += "PSH " }
    return flags
}

func main() {
    // To actually capture packets, you'd use gopacket with libpcap
    // This shows the parsing logic
    
    // Simulate an IP+TCP packet
    packet := []byte{
        0x45, 0x00, 0x00, 0x3c,  // version/IHL, TOS, total length
        0x00, 0x00, 0x40, 0x00,  // ID, flags/frag offset
        0x40, 0x06, 0x00, 0x00,  // TTL=64, protocol=TCP, checksum
        0xc0, 0xa8, 0x01, 0x64,  // src IP: 192.168.1.100
        0x08, 0x08, 0x08, 0x08,  // dst IP: 8.8.8.8
        // TCP header
        0xd3, 0x4a, 0x01, 0xbb,  // src port: 54090, dst port: 443
        0x00, 0x00, 0x00, 0x01,  // sequence number
        0x00, 0x00, 0x00, 0x00,  // ack number
        0x50, 0x02, 0x20, 0x00,  // data offset, flags=SYN, window
        0x00, 0x00, 0x00, 0x00,  // checksum, urgent
    }
    
    ip, _ := parseIPHeader(packet)
    fmt.Printf("IP: %s → %s (TTL=%d, Proto=%s)\n",
        ip.SrcIP, ip.DstIP, ip.TTL, protocolName(ip.Protocol))
    
    tcp, _ := parseTCPHeader(packet[ip.IHL:])
    fmt.Printf("TCP: %d → %d [%s] Window=%d\n",
        tcp.SrcPort, tcp.DstPort, tcp.FlagsString(), tcp.Window)
}
```

Real packet capture requires `gopacket` library + libpcap. The above shows how raw bytes become structured data.

---

## Summary

| Concept | Why it matters |
|---------|---------------|
| IP addresses | Every device has one; CIDR tells you the network scope |
| TCP handshake | Port scanning exploits this; session hijacking targets it |
| UDP | DNS/VoIP/gaming; less scrutinized, harder to scan |
| DNS | First step in recon; zone transfers can dump entire records |
| Ports | Services = attack surface; knowing common ports is essential |
| ARP | LAN attacks; MitM attacks through poisoning |
| Firewalls | Understand what blocks your scan vs what means "closed" |

---

## Exercises

1. Use `tcpdump -n -i eth0 tcp` to capture TCP traffic. Can you identify a full three-way handshake?
2. Perform a DNS zone transfer against a deliberately vulnerable DNS server (try `digi.ninja/projects/zonetransferme.htm`)
3. Run `arp -a` to see your ARP cache. Which entries are there? Can you identify your router's MAC address?
4. Use Wireshark to capture and analyze an HTTP request. Find the TCP handshake, the HTTP GET, and the response.
5. Write a Go program that sends a DNS query for `A` record of a domain using the `net` package.
