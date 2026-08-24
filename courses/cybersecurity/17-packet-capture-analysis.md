# Chapter 17: Packet Capture and Analysis — Seeing Raw Network Traffic

*Packet capture is the most powerful network diagnostic and offensive tool available. It lets you see exactly what's on the wire — credentials, tokens, data, and attack traffic.*

---

## Why Packet Analysis Matters

**Offensive uses:**
- Capture cleartext credentials (Telnet, FTP, HTTP Basic Auth)
- Session hijacking (steal authenticated cookies/tokens)
- Credential harvesting on compromised networks
- C2 traffic analysis

**Defensive uses:**
- Incident response — what data left the network?
- Intrusion detection — identify attack patterns
- Malware analysis — what does this malware phone home to?
- Troubleshooting — why is this connection failing?

---

## tcpdump — The Essential CLI Tool

```bash
# Capture on interface eth0
sudo tcpdump -i eth0

# Capture to file
sudo tcpdump -i eth0 -w capture.pcap

# Read from file
tcpdump -r capture.pcap

# Filter syntax (BPF - Berkeley Packet Filter)

# By host
tcpdump host 192.168.1.100

# By port
tcpdump port 80
tcpdump port 22 or port 443

# By protocol
tcpdump tcp
tcpdump udp
tcpdump icmp

# By direction
tcpdump src 192.168.1.100    # from this host
tcpdump dst 10.0.0.1         # to this host

# Complex filters
tcpdump 'tcp port 80 and host 192.168.1.1'
tcpdump 'tcp[tcpflags] & tcp-syn != 0'  # SYN packets only

# Show more detail
tcpdump -v   # verbose
tcpdump -vv  # more verbose
tcpdump -X   # show hex + ASCII
tcpdump -nn  # don't resolve names (faster)

# Useful combinations
sudo tcpdump -i eth0 -nn -X port 80 -w http_traffic.pcap
sudo tcpdump -i eth0 -nn 'port 443' | grep -v "length 0"
```

### Capturing Credentials

```bash
# HTTP Basic Auth (base64 encoded)
sudo tcpdump -i eth0 -A 'port 80' | grep -i "authorization: basic"

# FTP credentials
sudo tcpdump -i eth0 -A 'port 21' | grep -E "USER|PASS"

# Telnet (completely cleartext)
sudo tcpdump -i eth0 -A 'port 23'

# HTTP form data (POST bodies)
sudo tcpdump -i eth0 -A 'tcp port 80' | grep -E "password=|pass=|pwd="
```

---

## Wireshark — Visual Packet Analysis

Wireshark is the GUI equivalent of tcpdump, with far more analysis power.

**Key features:**
- Protocol dissectors for 3000+ protocols
- Follow TCP stream — reassemble entire conversation
- Statistics → Conversations — top talkers
- Filter expressions (same BPF syntax as tcpdump)
- Export objects — extract files transferred over HTTP/SMB

### Wireshark Display Filters

```
# IP address
ip.addr == 192.168.1.1
ip.src == 192.168.1.1 && ip.dst == 10.0.0.1

# Port
tcp.port == 80
tcp.dstport == 443

# Protocol
http
dns
tls

# HTTP specifics
http.request.method == "POST"
http.request.uri contains "login"
http contains "password"

# TCP flags
tcp.flags.syn == 1 && tcp.flags.ack == 0   # SYN packets (port scan)
tcp.flags.rst == 1                          # RST packets

# Find large transfers (potential exfiltration)
frame.len > 10000

# DNS anomalies
dns.qry.name contains ".onion"
dns.resp.ttl < 10               # fast flux DNS
```

---

## gopacket — Packet Capture in Go

```go
package main

import (
    "fmt"
    "log"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)

// Install: go get github.com/google/gopacket

func main() {
    // Open live capture
    handle, err := pcap.OpenLive("eth0", 65536, true, pcap.BlockForever)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()
    
    // Set BPF filter
    if err := handle.SetBPFFilter("tcp port 80"); err != nil {
        log.Fatal(err)
    }
    
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    
    for packet := range packetSource.Packets() {
        analyzePacket(packet)
    }
}

func analyzePacket(packet gopacket.Packet) {
    // IP layer
    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    if ipLayer == nil {
        return
    }
    ip := ipLayer.(*layers.IPv4)
    
    // TCP layer
    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer == nil {
        return
    }
    tcp := tcpLayer.(*layers.TCP)
    
    fmt.Printf("%s:%d → %s:%d ", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)
    
    // Flags
    flags := ""
    if tcp.SYN { flags += "SYN " }
    if tcp.ACK { flags += "ACK " }
    if tcp.FIN { flags += "FIN " }
    if tcp.RST { flags += "RST " }
    if tcp.PSH { flags += "PSH " }
    fmt.Printf("[%s] len=%d\n", flags, len(tcp.Payload))
    
    // HTTP detection
    if len(tcp.Payload) > 0 {
        payload := string(tcp.Payload)
        if len(payload) > 4 {
            if payload[:4] == "GET " || payload[:5] == "POST " ||
               payload[:4] == "HTTP" {
                fmt.Printf("  HTTP: %s\n", payload[:min(len(payload), 200)])
                // Look for credentials
                if containsCredentials(payload) {
                    fmt.Printf("  *** POTENTIAL CREDENTIALS DETECTED ***\n")
                }
            }
        }
    }
}

func containsCredentials(payload string) bool {
    keywords := []string{
        "password=", "passwd=", "pwd=", "pass=",
        "Authorization: Basic",
        "user=", "username=", "email=",
    }
    lower := strings.ToLower(payload)
    for _, kw := range keywords {
        if strings.Contains(lower, strings.ToLower(kw)) {
            return true
        }
    }
    return false
}

func min(a, b int) int { if a < b { return a }; return b }
```

---

## Reading PCAP Files

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)

type DNSQuery struct {
    Src    string
    Domain string
}

func analyzePCAP(filename string) {
    handle, err := pcap.OpenOffline(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer handle.Close()
    
    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    
    dnsQueries := make(map[string]int)
    tcpConns    := make(map[string]int)
    
    for packet := range packetSource.Packets() {
        // Count DNS queries
        if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
            dns := dnsLayer.(*layers.DNS)
            for _, q := range dns.Questions {
                domain := string(q.Name)
                dnsQueries[domain]++
            }
        }
        
        // Count TCP connections (SYN only)
        if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
            tcp := tcpLayer.(*layers.TCP)
            if tcp.SYN && !tcp.ACK {
                if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
                    ip := ipLayer.(*layers.IPv4)
                    key := fmt.Sprintf("%s → %s:%d", ip.SrcIP, ip.DstIP, tcp.DstPort)
                    tcpConns[key]++
                }
            }
        }
    }
    
    fmt.Println("=== DNS Queries (top 10) ===")
    // print top 10...
    
    fmt.Println("\n=== TCP Connections (SYN) ===")
    for conn, count := range tcpConns {
        if count > 1 {
            fmt.Printf("  %s (%d times)\n", conn, count)
        }
    }
}
```

---

## Summary

| Tool | Use case |
|------|---------|
| `tcpdump` | CLI capture, filtering, scripting |
| Wireshark | GUI analysis, protocol dissection |
| `gopacket` | Custom Go-based capture and analysis |
| `pcap` files | Store captures for later analysis |

---

## Exercises

1. Capture HTTP traffic to a DVWA instance and extract form POST data (username/password)
2. Write a Go program that reads a PCAP file and outputs all unique IP pairs that communicated
3. Use Wireshark to follow a TCP stream for an HTTP login. What does the full request look like?
4. Write a BPF filter for tcpdump that captures only DNS traffic and HTTP POST requests
