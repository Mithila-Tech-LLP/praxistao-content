# Chapter 13: What Is the Internet?

> **"The internet is not a place. It's not a building or a computer. It's an agreement — a set of rules that lets any two computers talk to each other, anywhere in the world, as long as they both follow the rules."**

---

## Table of Contents

1. [The Simple Answer](#1-the-simple-answer)
2. [How Computers Connect to Each Other](#2-how-computers-connect-to-each-other)
3. [IP Addresses — The Internet's Phone Book](#3-ip-addresses--the-internets-phone-book)
4. [Packets — How Data Travels](#4-packets--how-data-travels)
5. [How Data Physically Travels Across the World](#5-how-data-physically-travels-across-the-world)
6. [Wi-Fi and Cellular — Connecting Without Wires](#6-wi-fi-and-cellular--connecting-without-wires)
7. [What Happens When You Load Google.com](#7-what-happens-when-you-load-googlecom)
8. [The Internet vs. The Web](#8-the-internet-vs-the-web)
9. [Summary](#summary)

---

## 1. The Simple Answer

```
The internet is:
  Billions of computers (and phones, and servers, and TVs, and fridges)
  connected together by cables and wireless signals
  following the same rules (protocols) to send and receive data.
  
That's it. There's no single "internet machine" somewhere.
The internet IS the connection between all those devices.

How many devices?
  ~15 billion connected devices today.
  By 2030: ~50 billion (Internet of Things).
```

---

## 2. How Computers Connect to Each Other

Computers are connected through a hierarchy of networks:

```
Your phone / laptop
    ↓
Your home router (Wi-Fi or wired)
    ↓
Your ISP (Internet Service Provider)
(AT&T, Comcast, Jio, BT, Vodafone — whoever you pay for internet)
    ↓
Regional network
    ↓
Internet backbone
(massive fiber optic cables between cities and countries)
    ↓
Another ISP on the other side
    ↓
A server (Google's, Netflix's, your friend's computer)
```

**ISP (Internet Service Provider):**
The company you pay to give you internet access. They connect your home to the broader network.

**Router:**
The device in your home that:
- Connects to your ISP via a cable
- Creates your home Wi-Fi network
- Routes traffic between your devices and the internet
- Acts as a firewall (basic security)

---

## 3. IP Addresses — The Internet's Phone Book

Every device on the internet has an **IP address** — a unique number identifying it.

```
IPv4 address: four numbers 0–255, separated by dots
  Example: 142.250.80.46 ← this is one of Google's IP addresses
  
  142.250.80.46
  ↑           ↑
  Network     Device on that network

IPv6 address: newer, much more addresses available
  Example: 2607:f8b0:4004:0c08:0000:0000:0000:200e
  IPv4 can handle 4 billion addresses.
  IPv6 can handle 340 undecillion (3.4 × 10^38) addresses.
```

**Domain names:**
IP addresses are hard to remember. So we use domain names.

```
google.com   → 142.250.80.46
amazon.com   → 205.251.242.103
wikipedia.org → 208.80.154.224
```

**DNS (Domain Name System):**
Like a phone book for the internet. You type `google.com`, your computer asks a DNS server "what's the IP address for google.com?" and gets back the number. Then your computer connects to that number.

```
You type:  google.com
               ↓
         DNS lookup
               ↓
         142.250.80.46
               ↓
         Connect to Google
```

---

## 4. Packets — How Data Travels

The internet doesn't send files as one big blob. It breaks everything into small pieces called **packets**.

```
You send a photo (3MB):
  
  Photo is split into ~2,000 packets of ~1.5KB each
  Each packet contains:
    Where it came from (your IP)
    Where it's going (server IP)
    Packet number (1 of 2000, 2 of 2000, etc.)
    Checksum (verify data not corrupted)
    The actual data chunk
  
  Each packet may take a DIFFERENT route through the internet
  (whatever path is fastest at that moment)
  
  Destination reassembles packets in order
  If any packet is lost → destination asks for it again
  
This is TCP/IP: the fundamental protocol of the internet.
```

**Why packets instead of whole files?**
- Multiple people can share the same cable simultaneously (each packet interleaves)
- If a packet is dropped, only that packet is resent (not the whole file)
- Automatic load balancing (packets take whatever path is free)

---

## 5. How Data Physically Travels Across the World

The internet runs on physical infrastructure:

```
Within your city:
  Fiber optic cables    → glass cables that carry light pulses
  (Light travels at 200,000 km/second in glass)
  Coaxial cable         → old cable TV infrastructure
  Telephone copper wire → DSL internet
  
Between cities/countries:
  More fiber optic cables, buried underground or laid under the ocean
  
Under the ocean:
  Submarine cables — actual physical cables on the ocean floor
  
  Atlantic cable: New York ↔ London = 6,900 km
  Trans-Pacific: Los Angeles ↔ Tokyo = 9,000 km
  
  When you load a website from another country,
  your data physically travels under the ocean on these cables.
  
Satellite internet:
  SpaceX Starlink: 5,000+ satellites in low Earth orbit
  Sends data up to satellite → satellite sends to another → back down
  Covers areas without cable (remote areas, ships, planes)
  
Cell towers:
  Your phone connects to the nearest cell tower
  Cell tower connects via fiber to the internet backbone
```

---

## 6. Wi-Fi and Cellular — Connecting Without Wires

**Wi-Fi:**
```
Your router broadcasts a Wi-Fi signal using radio waves (2.4 GHz or 5 GHz)
Your device has a Wi-Fi chip that receives this signal
They exchange data wirelessly — up to ~100m range
The router then connects via cable to your ISP

Wi-Fi generations:
  Wi-Fi 4 (802.11n)    → ~150 Mbps
  Wi-Fi 5 (802.11ac)   → ~1,000 Mbps (most common currently)
  Wi-Fi 6/6E           → ~9,600 Mbps (newest routers/devices)
```

**Cellular (3G/4G/5G):**
```
Your phone has a cellular radio that talks to cell towers
Cell towers are every few kilometers in cities

4G LTE: typical 20–100 Mbps (current standard globally)
5G:     typical 100–1,000 Mbps (newer phones and areas)

The "G" = Generation of mobile network standards
5G promises very low latency (important for self-driving cars, VR)
```

---

## 7. What Happens When You Load Google.com

Let's trace exactly what happens in the ~200 milliseconds it takes to load google.com:

```
0ms     You type google.com and press Enter

1ms     Your browser checks its cache:
        "Have I been to google.com recently? If yes, skip DNS."
        
2ms     DNS lookup:
        Browser asks DNS server: "What's the IP for google.com?"
        DNS server responds: "142.250.80.46"

5ms     TCP connection:
        Your computer sends a "hello" (SYN packet) to 142.250.80.46
        Google's server responds "hello back" (SYN-ACK)
        Your computer confirms (ACK)
        (This is called the "TCP handshake")

10ms    TLS handshake (for HTTPS security):
        Your computer and Google's server agree on encryption
        
20ms    HTTP Request:
        Your browser sends: "GET /  HTTP/1.1"
        (please send me your homepage)

100ms   Network transit:
        Request travels through cables/routers to Google's server
        
120ms   Server processing:
        Google's server prepares the response HTML

150ms   Response arrives:
        HTML (the page structure) arrives at your browser
        Browser starts reading it, discovers it needs CSS/JS/images
        
170ms   More requests:
        Browser fetches CSS, JavaScript, fonts, images in parallel
        
200ms   Page rendered:
        Browser draws the page on screen
        You see Google's homepage
```

200 milliseconds. That's how long it takes for your click to travel potentially halfway around the world and back.

---

## 8. The Internet vs. The Web

Many people use "internet" and "web" interchangeably. They're different things.

```
The Internet:
  The physical + logical network connecting all computers
  The infrastructure: cables, routers, protocols (TCP/IP)
  Has existed since 1969 (ARPANET)
  
The Web (World Wide Web):
  ONE service that runs ON the internet
  Uses HTTP/HTTPS protocol to transfer web pages
  Invented by Tim Berners-Lee in 1991
  What you see in a browser

Other services that use the internet (not the web):
  Email (SMTP, IMAP) — not the web
  WhatsApp messages — not the web
  FaceTime / Zoom — not the web
  Online gaming — not the web
  
The Web is to the Internet as
Roads are to the Country —
one way to travel through it,
not the country itself.
```

---

## Summary

| Concept | What It Means |
|---------|--------------|
| Internet | Global network of all connected computers and devices |
| IP address | Unique number identifying each device on the internet |
| DNS | Translates domain names (google.com) to IP addresses |
| Packet | Small chunk of data (files are split into many packets) |
| TCP/IP | The fundamental rules (protocols) of the internet |
| Router | Device that connects your home network to the internet |
| ISP | Company you pay for internet access |
| Wi-Fi | Short-range wireless connection to your router |
| 4G/5G | Cellular network for phone internet access |

**The internet connects computers. The web is what you do with that connection. Next: what exactly IS the web?**
