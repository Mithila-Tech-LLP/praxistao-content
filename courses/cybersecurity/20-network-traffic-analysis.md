# Chapter 20: Network Traffic Analysis — Finding Attacks in Packets

*Traffic analysis bridges offense and defense. Attackers analyze traffic to find credentials and patterns; defenders analyze it to detect attacks, investigate incidents, and understand attacker behavior.*

---

## Flow Analysis vs Deep Packet Inspection

**NetFlow/IPFIX (flow data):**
- Summary: src IP, dst IP, port, protocol, bytes, packets
- No payload — low storage, scalable to millions of flows
- Good for: DDoS detection, port scans, beaconing, volume anomalies

**Full packet capture (DPI):**
- Complete network traffic
- High storage requirements
- Good for: credential capture, data exfiltration investigation, malware C2

For most enterprise security operations: flow data for detection, full capture for investigation.

---

## Detecting Attack Patterns in Traffic

### Port Scan Detection

```
Pattern: One source IP → many destination IPs/ports, rapid succession
```

```bash
# tcpdump: capture SYN packets from one source
sudo tcpdump -i eth0 'tcp[tcpflags] & tcp-syn != 0 and tcp[tcpflags] & tcp-ack == 0' -nn

# In Zeek (network monitoring framework):
# Look for port scan signatures in conn.log
cat conn.log | awk '{print $3, $5, $6}' | sort | uniq -c | sort -rn | head -20
# $3=src_ip, $5=dst_ip, $6=dst_port
```

### Beaconing Detection (C2 Communication)

C2 (command and control) traffic has regular intervals — malware "phones home" on a schedule.

```
Pattern: Fixed interval connections to same destination
         e.g., every 60 seconds to 185.220.101.5:443
```

```bash
# Extract connection times to an IP
tcpdump -r capture.pcap -nn 'host 185.220.101.5' | \
    awk '{print $1}' | head -20

# Calculate intervals
tcpdump -r capture.pcap -nn 'host 185.220.101.5' | \
    awk '{print $1}' | \
    awk 'NR>1{diff=$1-prev; printf "%.2f\n", diff} {prev=$1}' | \
    sort -n
# Consistent intervals (e.g., all ~60.0) = beaconing
```

### Data Exfiltration Detection

```
Pattern: Unusually large outbound transfer to unknown destination
         Especially: slow and steady (avoids rate-based detection)
```

```bash
# Find top outbound flows by byte count
# In Wireshark: Statistics → Conversations → TCP → Sort by bytes
# Filter: ip.dst not in known_good_list and frame.len > 5000

# tcpdump: capture large outbound packets
sudo tcpdump -i eth0 'dst not 192.168.1.0/24 and greater 1000' -nn
```

---

## Zeek (Bro) — Network Security Monitor

Zeek processes network traffic and generates structured logs — far more useful than raw pcap for monitoring.

```bash
# Run Zeek on a pcap
zeek -r capture.pcap

# Generates:
# conn.log       — all connections
# dns.log        — DNS queries
# http.log       — HTTP requests (URLs, headers, bodies)
# ssl.log        — TLS connections
# files.log      — transferred files
# weird.log      — protocol anomalies
# notice.log     — security notices

# Query http.log for POSTs
cat http.log | zeek-cut ts id.orig_h method uri | grep POST

# Find large transfers
cat files.log | zeek-cut ts id.orig_h total_bytes | awk '$3 > 100000'

# DNS anomalies
cat dns.log | zeek-cut ts id.orig_h query answers | \
    awk '{if(length($4) > 200) print}'  # unusually long DNS responses (tunneling?)
```

---

## Go: Traffic Analyzer (from PCAP)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strings"
)

// Analyze a tcpdump text output (tcpdump -nn -r file.pcap)
type Connection struct {
    SrcIP   string
    DstIP   string
    DstPort string
    Proto   string
    Count   int
}

type TrafficSummary struct {
    TopTalkers    map[string]int // src IP → byte count
    TopDestPorts  map[string]int // dst port → connection count
    DNSQueries    map[string]int // domain → query count
    SuspiciousIPs []string
}

func analyzeConnLog(filename string) *TrafficSummary {
    f, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer f.Close()
    
    summary := &TrafficSummary{
        TopTalkers:   make(map[string]int),
        TopDestPorts: make(map[string]int),
        DNSQueries:   make(map[string]int),
    }
    
    // Example: parse simplified conn.log format
    // timestamp src_ip src_port dst_ip dst_port proto bytes
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") {
            continue
        }
        
        fields := strings.Fields(line)
        if len(fields) < 7 {
            continue
        }
        
        srcIP  := fields[1]
        dstPort := fields[4]
        
        summary.TopTalkers[srcIP]++
        summary.TopDestPorts[dstPort]++
    }
    
    // Flag suspicious: IPs connecting to >100 unique ports (port scan)
    for ip, count := range summary.TopTalkers {
        if count > 100 {
            summary.SuspiciousIPs = append(summary.SuspiciousIPs, 
                fmt.Sprintf("%s (%d connections)", ip, count))
        }
    }
    
    return summary
}

type kv struct{ k string; v int }

func topN(m map[string]int, n int) []kv {
    var pairs []kv
    for k, v := range m {
        pairs = append(pairs, kv{k, v})
    }
    sort.Slice(pairs, func(i, j int) bool { return pairs[i].v > pairs[j].v })
    if len(pairs) > n { return pairs[:n] }
    return pairs
}

func main() {
    summary := analyzeConnLog("conn.log")
    if summary == nil {
        fmt.Println("Could not read log file")
        return
    }
    
    fmt.Println("=== Top Talkers ===")
    for _, p := range topN(summary.TopTalkers, 10) {
        fmt.Printf("  %-20s %d connections\n", p.k, p.v)
    }
    
    fmt.Println("\n=== Most Accessed Ports ===")
    for _, p := range topN(summary.TopDestPorts, 10) {
        fmt.Printf("  Port %-10s %d times\n", p.k, p.v)
    }
    
    if len(summary.SuspiciousIPs) > 0 {
        fmt.Println("\n=== Suspicious IPs (possible port scan) ===")
        for _, ip := range summary.SuspiciousIPs {
            fmt.Println(" ", ip)
        }
    }
}
```

---

## Wireless Traffic Analysis

```bash
# Put WiFi card in monitor mode
sudo airmon-ng start wlan0   # creates wlan0mon

# Capture all WiFi traffic
sudo airodump-ng wlan0mon

# Capture specific AP
sudo airodump-ng -c 6 --bssid AA:BB:CC:DD:EE:FF -w capture wlan0mon

# View capture in Wireshark
# Filter: wlan.fc.type_subtype == 0x08  (beacon frames)
# Filter: eapol  (WPA handshakes)
```

---

## Summary

| Technique | Detection target |
|-----------|-----------------|
| Flow analysis | DDoS, port scans, beaconing, large transfers |
| Packet capture | Credentials, session tokens, data content |
| Zeek logs | Structured analysis, protocol-level visibility |
| DNS analysis | C2 domain generation, DNS tunneling |
| Beacon detection | C2 communication intervals |

---

## Exercises

1. Download a malware PCAP from malware-traffic-analysis.net. Identify the C2 server, malware family, and what data was exfiltrated.
2. Write a beaconing detector in Go: takes a pcap file, groups connections by dst IP, calculates the standard deviation of connection intervals, flags IPs with std dev < 5 seconds.
3. Use Zeek on a pcap and extract all HTTP URLs accessed. Are any suspicious?
4. Write a Go program that reads tcpdump output and calculates the top 10 "chattiest" pairs (src→dst with most packets).
