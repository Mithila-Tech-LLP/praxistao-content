# Chapter 117: Building a Simple VPN

> **"Chapter 85 drew a box inside a box inside a box and called it a tunnel. This chapter makes that drawing literal: a real file descriptor the kernel hands you raw IP packets through, a real AES key encrypting them, and a real UDP socket carrying the result across an untrusted network to a peer that undoes all three steps in reverse."**

---

## Table of Contents

1. [Recap: Chapter 85's Tunneling, Chapters 78-79's Crypto](#1-recap-chapter-85s-tunneling-chapters-78-79s-crypto)
2. [The Problem: Turning "Wrap One Packet In Another" Into Actual Code](#2-the-problem-turning-wrap-one-packet-in-another-into-actual-code)
3. [A Naive Attempt: A Plaintext UDP Relay — And Why That Isn't a VPN](#3-a-naive-attempt-a-plaintext-udp-relay--and-why-that-isnt-a-vpn)
4. [The Real Solution: TUN Device + AEAD Envelope + UDP Transport](#4-the-real-solution-tun-device--aead-envelope--udp-transport)
5. [Setting Up a Linux TUN Interface (Requires root)](#5-setting-up-a-linux-tun-interface-requires-root)
6. [Code: Opening the TUN Device](#6-code-opening-the-tun-device)
7. [Code: The Encrypted Envelope (AES-256-GCM)](#7-code-the-encrypted-envelope-aes-256-gcm)
8. [Code: The VPN Peer — TUN → Encrypt → UDP](#8-code-the-vpn-peer--tun--encrypt--udp)
9. [Code: The VPN Peer — UDP → Decrypt → TUN](#9-code-the-vpn-peer--udp--decrypt--tun)
10. [Code: main() — Wiring Both Directions Together](#10-code-main--wiring-both-directions-together)
11. [The Fallback: A Simulated TUN for Readers Without root](#11-the-fallback-a-simulated-tun-for-readers-without-root)
12. [Hands-On Experiment: Two Peers, ping Across the Tunnel](#12-hands-on-experiment-two-peers-ping-across-the-tunnel)
13. [Packet-Level Walkthrough: One ICMP Echo, Start to Finish](#13-packet-level-walkthrough-one-icmp-echo-start-to-finish)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Production Notes: What WireGuard/IPsec/OpenVPN Add On Top](#15-production-notes-what-wireguardipsecopenvpn-add-on-top)
16. [What's Simplified Here](#16-whats-simplified-here)
17. [Interview Questions & Model Answers](#17-interview-questions--model-answers)
18. [Exercises](#18-exercises)
19. [Summary](#19-summary)
20. [What's Next in This Volume](#20-whats-next-in-this-volume)

---

## 1. Recap: Chapter 85's Tunneling, Chapters 78-79's Crypto

Chapter 85 established tunneling as the general idea underneath every VPN: take a whole packet, private headers and all, and use it as the payload of another packet that a public network already knows how to route. It compared IPsec, OpenVPN, and WireGuard as three real implementations of that idea, at three different points on the complexity/speed trade-off.

Chapter 78 gave you AES and the AEAD construction (`GCM`) that turns encryption into encryption-plus-tamper-detection in one pass. Chapter 79 explained the problem AES alone doesn't solve — how two strangers agree on a shared key over a network an eavesdropper is watching — and Chapter 115 showed that a UDP socket can stand in for a "virtual wire" between two programs on a real network.

This chapter fuses all three. You will open a real Linux TUN virtual network interface, read whatever raw IP packets the kernel hands you through it, seal each one in an AES-256-GCM envelope, send that envelope as one UDP datagram to a peer, and have the peer reverse every step — decrypt, verify, and hand the original packet back to *its* kernel through *its own* TUN device. That is a minimal, working VPN: not a simulation of one, an actual one, small enough to read in one sitting.

---

## 2. The Problem: Turning "Wrap One Packet In Another" Into Actual Code

Chapter 85's diagram is easy to draw and easy to agree with. Making it real raises three concrete engineering questions that prose glossed over:

1. **Where do the "inner" packets even come from?** Chapter 114's packet sniffer read packets that were *already* traveling across a real interface. A VPN needs the opposite: it needs the *operating system itself* to hand it packets that some application generated for a destination on a private network that doesn't physically exist yet — and to accept packets back for delivery to that same application. That's exactly what a **TUN virtual interface** is for.
2. **How does "encrypt the payload" actually work, byte for byte?** Chapter 78 explained AES-GCM conceptually. This chapter has to pick an exact envelope layout: where does the nonce go, where does the authentication tag go, and how does the receiver know where one ends and the next begins.
3. **What carries the encrypted result across the real network?** Chapter 115 already answered this for a different problem (a toy router) — an ordinary UDP socket, exactly like Chapter 107 built, works perfectly as the "outer" transport, because a VPN envelope is just an opaque blob of bytes as far as the public Internet is concerned.

---

## 3. A Naive Attempt: A Plaintext UDP Relay — And Why That Isn't a VPN

The simplest thing that could possibly "tunnel" packets is: read a packet from the TUN device, send its raw bytes over UDP to a peer, have the peer write those exact bytes into its own TUN device. No encryption at all.

This actually works, mechanically — packets do arrive on the other side. But it fails the one requirement that makes a VPN a *security* tool rather than just a relay: Chapter 85 Section 2 was explicit that only the VPN server, holding a shared key, should be able to recover the inner packet. A plaintext relay hands the inner packet — private IP addresses, TCP ports, and unencrypted application payloads — to every router and every passive eavesdropper on the path between the two peers, in the clear. It is functionally a tunnel and not at all a *virtual private* network. The gap between "packets arrive" and "packets arrive safely" is exactly Chapter 78's AEAD encryption, and it's the one piece the naive version is missing.

---

## 4. The Real Solution: TUN Device + AEAD Envelope + UDP Transport

```
 Host A                                                        Host B
 ┌─────────────┐                                          ┌─────────────┐
 │ application │  writes/reads IP packets to 10.0.0.2      │ application │
 │ (e.g. ping) │◄──────────────┐                ┌─────────►│ (e.g. ping) │
 └──────┬──────┘               │                │          └──────┬──────┘
        │ kernel routes to tun0│                │tun1 delivers to │
        ▼                      │                │      kernel     ▼
 ┌─────────────┐         ┌─────┴─────┐    ┌──────┴────┐    ┌─────────────┐
 │  tun0 (fd)  │────────►│ this Go   │    │  this Go  │◄───│  tun1 (fd)  │
 │ 10.0.0.1/32 │  read   │ program A │    │ program B │write│ 10.0.0.2/32│
 └─────────────┘         │  seal()   │    │  open()   │    └─────────────┘
                          └─────┬─────┘    └─────▲─────┘
                                │ UDP, encrypted  │ UDP, encrypted
                                └──────────►──────┘
                              public Internet (only this
                              blob of bytes is visible)
```

The pipeline in one sentence: **TUN gives you the packet, AES-GCM seals it, UDP carries the seal, the peer's AES-GCM opens it, the peer's TUN delivers it.** Every step is something an earlier chapter already built or explained — this chapter's only genuinely new piece is the TUN device itself.

---

## 5. Setting Up a Linux TUN Interface (Requires root)

TUN devices are a Linux (and BSD/macOS-with-a-different-API) kernel feature that only root, or a process with the `CAP_NET_ADMIN` capability, can create. If you don't have access to a Linux root shell right now, skip ahead to Section 11 — everything else in this chapter still runs.

**Two Linux machines (or two network namespaces on one machine) are needed** — one per VPN peer. On each machine, after starting this chapter's Go program (which creates the interface as a side effect of the `TUNSETIFF` call in Section 6), configure it as a point-to-point link in a second terminal:

```
# Host A ("server" side, private address 10.0.0.1)
$ sudo ip addr add 10.0.0.1/32 peer 10.0.0.2/32 dev tun0
$ sudo ip link set dev tun0 mtu 1400   # see Section 14 for why this is 1400, not 1500
$ sudo ip link set tun0 up

# Host B ("client" side, private address 10.0.0.2)
$ sudo ip addr add 10.0.0.2/32 peer 10.0.0.1/32 dev tun1
$ sudo ip link set dev tun1 mtu 1400
$ sudo ip link set tun1 up
```

`peer` here declares a point-to-point route rather than an ordinary subnet — exactly right for a two-node tunnel, and avoids two interfaces on the same host or LAN both believing they own the same `/24`. If you only have one Linux machine available, `ip netns add` gives you two isolated network namespaces to play both roles safely without a second physical host.

---

## 6. Code: Opening the TUN Device

```go
//go:build linux

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	tunsetiff  = 0x400454ca // TUNSETIFF ioctl request number (linux/if_tun.h)
	iffTun     = 0x0001     // IFF_TUN: layer-3 IP tunnel (vs. IFF_TAP's layer-2 Ethernet)
	iffNoPI    = 0x1000     // IFF_NO_PI: don't prefix packets with a 4-byte protocol-info header
	ifNameSize = 16         // IFNAMSIZ
)

// ifReq mirrors just enough of the kernel's struct ifreq (a fixed 40 bytes on
// Linux) for the TUNSETIFF ioctl: an interface name, a flags field, and
// padding to match the kernel's real struct size.
type ifReq struct {
	Name  [ifNameSize]byte
	Flags uint16
	pad   [40 - ifNameSize - 2]byte
}

// openTUN opens /dev/net/tun and asks the kernel to attach it to an
// interface named `name` (creating it if it doesn't exist), configured as a
// pure IP tunnel — Chapter 85 Section 2's picture made real: from this point
// on, the kernel hands this program whole, raw IP packets to do something
// with, instead of this program having to construct or parse Ethernet frames.
func openTUN(name string) (*os.File, error) {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/net/tun: %w (are you root?)", err)
	}
	var req ifReq
	copy(req.Name[:], name)
	req.Flags = iffTun | iffNoPI

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(tunsetiff), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		syscall.Close(fd)
		return nil, fmt.Errorf("ioctl TUNSETIFF: %w", errno)
	}
	return os.NewFile(uintptr(fd), "/dev/net/tun"), nil
}
```

`IFF_NO_PI` matters: without it, the kernel prepends a 4-byte "protocol information" header (2 flag bytes + 2 bytes of `EtherType`) to every packet read from the device, which this chapter's code would otherwise have to strip before treating the rest as a clean IPv4 packet. Asking for `IFF_NO_PI` means every `Read` from the resulting `*os.File` returns exactly one raw IP packet, and nothing else — the cleanest possible interface for Section 8's code to build on.

---

## 7. Code: The Encrypted Envelope (AES-256-GCM)

```go
// deriveKey turns a human-typed shared passphrase into a fixed 32-byte
// AES-256 key. AES itself doesn't care where a key comes from (Chapter 78),
// but a raw passphrase is the wrong length and lower-entropy per byte than a
// real key — SHA-256 here is a stand-in for what a real system uses a proper
// KDF for (Section 16 explains the gap).
func deriveKey(passphrase string) []byte {
	sum := sha256.Sum256([]byte(passphrase))
	return sum[:] // exactly 32 bytes -> selects AES-256, not AES-128
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	return cipher.NewGCM(block) // Chapter 78 Section 6's AEAD mode: confidentiality + integrity in one pass
}
```

The envelope this chapter sends over UDP has a fixed, simple layout:

| Field | Size | Contents |
|---|---|---|
| Nonce | 12 bytes | Random, unique per packet — GCM's required IV, sent in the clear (it must be, so the receiver can decrypt) |
| Ciphertext | = plaintext length | The encrypted raw IP packet read from TUN |
| Auth tag | 16 bytes | GCM's authentication tag — detects any bit flipped in transit, whether by noise or by an attacker |

```go
// sealPacket encrypts one raw IP packet into a self-contained envelope:
// [12-byte nonce][ciphertext][16-byte GCM tag], per the table above.
func sealPacket(aead cipher.AEAD, packet []byte) ([]byte, error) {
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Seal appends ciphertext+tag after dst; passing nonce as dst makes the
	// nonce the envelope's own prefix — exactly the layout in the table above.
	return aead.Seal(nonce, nonce, packet, nil), nil
}

// openPacket reverses sealPacket, or returns an error if the envelope was
// corrupted or tampered with, or the wrong key was used — an AEAD mode
// (Chapter 78 Section 6) cannot tell you which; it can only refuse to
// produce plaintext it can't vouch for (Chapter 80's integrity guarantee).
func openPacket(aead cipher.AEAD, envelope []byte) ([]byte, error) {
	ns := aead.NonceSize()
	if len(envelope) < ns {
		return nil, fmt.Errorf("envelope too short to contain a nonce")
	}
	nonce, ciphertext := envelope[:ns], envelope[ns:]
	return aead.Open(nil, nonce, ciphertext, nil)
}
```

---

## 8. Code: The VPN Peer — TUN → Encrypt → UDP

```go
// VPNPeer holds one side's complete state: the TUN device it reads/writes
// raw IP packets through, the UDP socket it uses as the outer transport
// (Chapter 107's connectionless socket, doing exactly what it was built
// for), and the shared AEAD used to protect every envelope.
type VPNPeer struct {
	tun      io.ReadWriteCloser // a real TUN device, or Section 11's simulated fallback
	conn     *net.UDPConn
	peerAddr *net.UDPAddr // may start nil and be learned dynamically — see udpToTun
	mu       sync.Mutex
	aead     cipher.AEAD
}

// tunToUDP is one direction of the tunnel: whatever the kernel hands us on
// the TUN device is, by construction, a whole raw IP packet addressed to
// something on the private network on the other side (Chapter 85 Section
// 2). Seal it whole, in one UDP datagram, to the peer.
func (p *VPNPeer) tunToUDP() {
	buf := make([]byte, 65535)
	for {
		n, err := p.tun.Read(buf)
		if err != nil {
			log.Printf("tun read error: %v", err)
			return
		}
		p.mu.Lock()
		peerAddr := p.peerAddr
		p.mu.Unlock()
		if peerAddr == nil {
			log.Printf("no peer known yet, dropping %d-byte outbound packet (waiting for -peer or an inbound datagram)", n)
			continue
		}
		envelope, err := sealPacket(p.aead, buf[:n])
		if err != nil {
			log.Printf("seal error: %v", err)
			continue
		}
		if _, err := p.conn.WriteToUDP(envelope, peerAddr); err != nil {
			log.Printf("udp write error: %v", err)
			continue
		}
		log.Printf("tun -> udp: sealed %d bytes into a %d-byte envelope, sent to %s", n, len(envelope), peerAddr)
	}
}
```

---

## 9. Code: The VPN Peer — UDP → Decrypt → TUN

```go
// udpToTun is the reverse direction: an envelope arrives from the public
// network, gets authenticated and decrypted, and — if and only if that
// succeeds — the recovered raw IP packet is handed to the kernel through
// the TUN device, exactly as if it had arrived on a real interface.
func (p *VPNPeer) udpToTun() {
	buf := make([]byte, 65535)
	for {
		n, addr, err := p.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("udp read error: %v", err)
			return
		}

		p.mu.Lock()
		known := p.peerAddr
		p.mu.Unlock()
		if known != nil && !addr.IP.Equal(known.IP) {
			log.Printf("dropping envelope from unexpected sender %s (expected %s)", addr, known)
			continue
		}

		packet, err := openPacket(p.aead, buf[:n])
		if err != nil {
			// Could be a corrupted datagram, a stale/wrong key, or a forged
			// one — GCM correctly refuses to guess which (Chapter 78/80).
			log.Printf("decrypt/auth failed for %d-byte envelope from %s, dropping: %v", n, addr, err)
			continue
		}

		p.mu.Lock()
		if p.peerAddr == nil {
			p.peerAddr = addr // "roaming": learn the peer's endpoint from its first valid packet
			log.Printf("learned peer address dynamically: %s", addr)
		}
		p.mu.Unlock()

		if _, err := p.tun.Write(packet); err != nil {
			log.Printf("tun write error: %v", err)
			continue
		}
		log.Printf("udp -> tun: opened %d bytes from %s, wrote %d-byte packet to tun", n, addr, len(packet))
	}
}
```

Learning `peerAddr` from the first successfully-authenticated packet — rather than requiring it to be configured on both sides — mirrors how real WireGuard handles a roaming client whose public IP changes: it trusts whichever source address a validly-decrypted packet most recently arrived from, and updates its notion of "the peer" accordingly.

---

## 10. Code: main() — Wiring Both Directions Together

```go
func main() {
	tunName := flag.String("tun", "tun0", "TUN interface name (Linux, requires root)")
	listenAddr := flag.String("listen", ":55555", "UDP address to listen on")
	peerAddr := flag.String("peer", "", "peer's host:port to send encrypted packets to (optional — see Section 9)")
	passphrase := flag.String("key", "correct horse battery staple", "shared passphrase (Chapter 79 would replace this with a real key exchange)")
	simulate := flag.Bool("simulate", false, "use the simulated TUN fallback instead of a real device (Section 11)")
	flag.Parse()

	var tun io.ReadWriteCloser
	var err error
	if *simulate {
		tun = newSimulatedTUN()
	} else {
		tun, err = openTUN(*tunName)
		if err != nil {
			log.Fatalf("openTUN: %v (try -simulate if you're not root, or not on Linux)", err)
		}
	}
	defer tun.Close()

	udpAddr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	aead, err := newAEAD(deriveKey(*passphrase))
	if err != nil {
		log.Fatal(err)
	}

	peer := &VPNPeer{tun: tun, conn: conn, aead: aead}
	if *peerAddr != "" {
		resolved, err := net.ResolveUDPAddr("udp", *peerAddr)
		if err != nil {
			log.Fatal(err)
		}
		peer.peerAddr = resolved
	}

	log.Printf("VPN peer up: tun=%s listen=%s peer=%q simulate=%v", *tunName, *listenAddr, *peerAddr, *simulate)
	go peer.tunToUDP()
	peer.udpToTun() // block the main goroutine here
}
```

---

## 11. The Fallback: A Simulated TUN for Readers Without root

Not every reader has a spare Linux root shell handy. This fallback swaps the real TUN device for an in-memory `io.ReadWriteCloser` that manufactures structurally real (correctly checksummed) IPv4/ICMP packets on a timer instead of reading them from a kernel — every other line of code in this chapter, from Section 7's `sealPacket` through Section 9's `udpToTun`, runs completely unmodified against it.

```go
// simulatedTUN behaves like a TUN device that periodically originates a
// synthetic ICMP echo request "from the private network," and prints
// whatever decrypted packet it's handed back — letting you observe the
// full seal/send/receive/open pipeline without root or a Linux kernel at all.
type simulatedTUN struct {
	outbound chan []byte
}

func newSimulatedTUN() *simulatedTUN {
	s := &simulatedTUN{outbound: make(chan []byte, 8)}
	go s.generate()
	return s
}

func (s *simulatedTUN) generate() {
	var seq uint16
	for {
		time.Sleep(3 * time.Second)
		seq++
		s.outbound <- buildFakeICMPEcho("10.0.0.1", "10.0.0.2", seq)
	}
}

func (s *simulatedTUN) Read(buf []byte) (int, error) {
	packet := <-s.outbound
	return copy(buf, packet), nil
}

func (s *simulatedTUN) Write(buf []byte) (int, error) {
	// A real TUN device would hand this decrypted packet to the kernel for
	// routing to a real application. The simulation just parses and prints
	// it, proving decryption worked without a kernel in the loop.
	fmt.Printf("[simulated TUN] received a decrypted packet, %d bytes: %s\n", len(buf), describeIPPacket(buf))
	return len(buf), nil
}

func (s *simulatedTUN) Close() error { return nil }

// buildFakeICMPEcho constructs a minimal, correctly checksummed IPv4 packet
// carrying an ICMP echo request — reusing Chapter 115 Section 7's
// checksum/TTL approach, so the bytes traveling through this simulation are
// real IP packets, not arbitrary placeholder bytes.
func buildFakeICMPEcho(srcIP, dstIP string, seq uint16) []byte {
	icmp := make([]byte, 8)
	icmp[0] = 8 // Type 8: Echo Request
	icmp[1] = 0 // Code 0
	binary.BigEndian.PutUint16(icmp[4:6], 1)   // identifier
	binary.BigEndian.PutUint16(icmp[6:8], seq) // sequence number
	binary.BigEndian.PutUint16(icmp[2:4], internetChecksum(icmp))

	header := make([]byte, 20)
	header[0] = 0x45 // version 4, IHL 5 (20 bytes, no options)
	totalLen := 20 + len(icmp)
	binary.BigEndian.PutUint16(header[2:4], uint16(totalLen))
	binary.BigEndian.PutUint16(header[4:6], seq) // identification
	header[8] = 64                               // TTL (Chapter 54)
	header[9] = 1                                // protocol 1 = ICMP
	copy(header[12:16], net.ParseIP(srcIP).To4())
	copy(header[16:20], net.ParseIP(dstIP).To4())
	binary.BigEndian.PutUint16(header[10:12], internetChecksum(header))

	return append(header, icmp...)
}

// internetChecksum is the one's-complement checksum from Chapter 19, the
// same algorithm Chapter 115 Section 7 used.
func internetChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(data); i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func describeIPPacket(b []byte) string {
	if len(b) < 20 {
		return "too short to be a valid IPv4 header"
	}
	return fmt.Sprintf("%s -> %s, protocol %d", net.IP(b[12:16]), net.IP(b[16:20]), b[9])
}
```

---

## 12. Hands-On Experiment: Two Peers, ping Across the Tunnel

**With real TUN devices** (two Linux hosts, or two network namespaces, configured per Section 5):

```
# Host A
$ sudo go run . -tun=tun0 -listen=:55555 -peer=<hostB-ip>:55555 -key=demo-shared-secret
VPN peer up: tun=tun0 listen=:55555 peer="<hostB-ip>:55555" simulate=false

# Host B
$ sudo go run . -tun=tun1 -listen=:55555 -peer=<hostA-ip>:55555 -key=demo-shared-secret
VPN peer up: tun=tun1 listen=:55555 peer="<hostA-ip>:55555" simulate=false

# From Host A, ping the private address that only exists inside the tunnel:
$ ping -c 3 10.0.0.2
64 bytes from 10.0.0.2: icmp_seq=1 ttl=64 time=1.8 ms
64 bytes from 10.0.0.2: icmp_seq=2 ttl=64 time=1.2 ms
64 bytes from 10.0.0.2: icmp_seq=3 ttl=64 time=1.1 ms
```

The two programs' own logs show every packet crossing:

```
[hostA] tun -> udp: sealed 84 bytes into a 112-byte envelope, sent to hostB:55555
[hostB] udp -> tun: opened 112 bytes from hostA:55555, wrote 84-byte packet to tun
[hostB] tun -> udp: sealed 84 bytes into a 112-byte envelope, sent to hostA:55555
[hostA] udp -> tun: opened 112 bytes from hostB:55555, wrote 84-byte packet to tun
```

And `tcpdump` on either host's *real* interface confirms Chapter 85 Section 2's claim directly — the outer network sees nothing but opaque UDP datagrams:

```
$ sudo tcpdump -ni eth0 udp port 55555
14:02:11.442113 IP hostA.55555 > hostB.55555: UDP, length 112
14:02:11.443560 IP hostB.55555 > hostA.55555: UDP, length 112
```

**Without root, using the simulated fallback:**

```
$ go run . -simulate -listen=:55555 -peer=127.0.0.1:55556 -key=demo &
$ go run . -simulate -listen=:55556 -peer=127.0.0.1:55555 -key=demo &

[simulated TUN] received a decrypted packet, 28 bytes: 10.0.0.1 -> 10.0.0.2, protocol 1
[simulated TUN] received a decrypted packet, 28 bytes: 10.0.0.1 -> 10.0.0.2, protocol 1
```

Every step this chapter built — sealing, sending, receiving, opening — runs identically in both cases; only where the packets originate differs.

---

## 13. Packet-Level Walkthrough: One ICMP Echo, Start to Finish

1. `ping` on Host A constructs an ICMP echo request addressed to `10.0.0.2` (Chapter 54).
2. The kernel's routing table (Chapter 45) has a route for `10.0.0.2` pointing at `tun0` (installed by Section 5's `ip addr ... peer ...` command) — so it writes the raw IP packet to the TUN file descriptor instead of a physical NIC.
3. Section 8's `tunToUDP` goroutine is blocked in `p.tun.Read`; it wakes up with those exact 84 bytes.
4. `sealPacket` generates a fresh random 12-byte nonce, runs AES-256-GCM over the 84-byte packet, and returns `nonce || ciphertext || tag` — 112 bytes total (84 + 12 + 16).
5. `conn.WriteToUDP` sends those 112 bytes as one UDP datagram to Host B's `55555`. Every router between the two hosts forwards it purely on the outer IP/UDP headers — none of them can see, and don't need to see, that it's carrying ICMP at all.
6. Host B's `udpToTun` goroutine wakes up from `ReadFromUDP` with the 112-byte envelope and Host A's address.
7. `openPacket` splits off the 12-byte nonce, decrypts and verifies the remaining 100 bytes, and returns the original 84-byte ICMP packet — or an error, if anything was altered.
8. `p.tun.Write` hands those 84 bytes to Host B's kernel through `tun1`. The kernel sees an ordinary-looking IP packet addressed to its own `10.0.0.2` and delivers it up to its ICMP stack exactly as if it had arrived on a physical wire.
9. The kernel's ICMP implementation replies with an echo reply, and the entire six-step process repeats in reverse, symmetric down to the byte counts.

---

## 14. Common Misconceptions

- **"The VPN encrypts the outer UDP/IP headers too."** It cannot — the outer headers must stay in plaintext, or no router on the public Internet could route the envelope at all. Only the *inner* packet, the part `sealPacket` operates on, is protected. This is exactly Chapter 85 Section 2's diagram: the outer box is visible to everyone; only what's inside the inner boxes is hidden.
- **"A SHA-256-hashed passphrase is a secure key exchange."** It's a *shared secret*, not an *exchange* — both operators must already agree on the same passphrase out of band before running the program. Chapter 79's entire subject was solving key agreement over a network an eavesdropper is watching; this chapter deliberately sidesteps that problem (Section 16) rather than solving it.
- **"TUN and TAP are interchangeable names for the same thing."** `IFF_TUN` delivers layer-3 IP packets (what this chapter reads and writes); `IFF_TAP` delivers whole layer-2 Ethernet frames, MAC headers included, and is what you'd reach for to bridge two remote Ethernet segments rather than route between two IP subnets.
- **"A decryption failure means the network corrupted the datagram."** It might — or it might mean an attacker forged garbage, or the two sides' keys don't actually match. GCM's authentication tag can only tell you *that* the datagram isn't trustworthy, never *why* (Chapter 80's integrity guarantee, not a diagnostic tool).
- **"The 28 bytes of crypto overhead is negligible."** 12 bytes of nonce plus 16 bytes of tag, on top of an 8-byte UDP header and a 20-byte outer IP header, is 56 bytes of total overhead per packet. A full 1500-byte inner packet would push the outer UDP datagram past the typical 1500-byte link MTU, forcing IP fragmentation — which is exactly why Section 5 set the TUN interface's own MTU to 1400, leaving headroom for every layer wrapped around it.

---

## 15. Production Notes: What WireGuard/IPsec/OpenVPN Add On Top

- **A real handshake, not a static shared key.** WireGuard runs a Noise-protocol-based handshake using ephemeral Diffie-Hellman key exchanges (Chapter 79) for every session, giving **forward secrecy** — compromising today's key doesn't expose yesterday's traffic. This chapter's `deriveKey` produces the same key forever, which has no such property.
- **Replay protection.** Real tunnel protocols track a monotonically increasing counter (or a sliding window of recently-seen values) per peer and reject any envelope whose counter has already been seen — defeating an attacker who simply records and re-sends ("replays") a captured, still-validly-encrypted envelope. This chapter's `openPacket` has no such check: a captured envelope can be resent successfully as long as the key hasn't rotated.
- **Multiple peers and routing, not just one.** WireGuard's `AllowedIPs` config lets one interface multiplex many peers, routing each inner destination IP to the correct peer's tunnel — this chapter hardcodes exactly one peer per running instance.
- **IKE's two-phase negotiation (IPsec, Chapter 85 Section 4)** and **TLS-based session setup (OpenVPN)** both solve the same key-agreement problem this chapter skips, at the cost of the implementation complexity Chapter 85 discussed directly.
- **Kernel-space performance.** Modern WireGuard runs as a Linux kernel module, avoiding the per-packet user-space/kernel-space copy this chapter's `os.File`-based TUN I/O pays on every single packet — a meaningful throughput difference at multi-gigabit line rates, though irrelevant for this chapter's teaching goals.

---

## 16. What's Simplified Here

- No key exchange — a static, pre-shared passphrase stands in for a real Diffie-Hellman/Noise handshake (Chapter 79).
- No replay protection — a captured envelope can be resent and will decrypt successfully again.
- Exactly one peer per instance, not a routing table of many peers and allowed subnets.
- IPv4 only, and the code assumes Linux's specific `/dev/net/tun` ioctl interface — macOS's `utun` and Windows' `wintun` use different mechanisms entirely, which is exactly why Section 11's fallback exists.
- No automatic interface configuration — `ip addr`/`ip link` are run by hand rather than by the program itself, unlike a production VPN client.
- No key rotation or session renegotiation — the same AES key is used for the entire lifetime of the process.

---

## 17. Interview Questions & Model Answers

**Beginner: What's the difference between what a TUN device delivers and what a raw Ethernet capture (Chapter 114) delivers?**
A TUN device, opened with `IFF_TUN`, delivers whole IP packets — no Ethernet header, no MAC addresses, just the IP layer and everything above it (Chapter 26's IP-and-up). A raw packet capture off a real interface (Chapter 114) sees full Ethernet frames including MAC headers, because it's observing an actual layer-2 medium. `IFF_TAP` would give TUN-style layer-2 frames instead, if that were needed.

**Intermediate: Why must the nonce in Section 7's envelope be sent in plaintext alongside the ciphertext, and why does that not weaken the encryption?**
GCM decryption is mathematically defined in terms of the same nonce that was used to encrypt — without it, the receiver has no way to reverse the keystream generation, so the nonce has to be transmitted somehow. It doesn't weaken security because a nonce isn't secret by design; what actually matters is that it's *never reused with the same key* (Chapter 78's AES-GCM caveat) — reuse, not disclosure, is what breaks GCM's security guarantees.

**Advanced: This chapter's implementation has no replay protection. Describe exactly how you'd add it, and what state each peer would need to keep.**
Each peer needs to remember, per sender, either a strictly increasing counter it expects the next packet to exceed, or a bitmap ("sliding window") of recently-accepted sequence numbers to tolerate limited reordering without rejecting legitimately delayed packets. The sender embeds a monotonically increasing counter as GCM's "associated data" (or folds it into the nonce construction itself, as WireGuard does) so the value is authenticated and can't be forged; on receipt, `openPacket` succeeding is necessary but no longer sufficient — the receiver must additionally check the counter against its window and reject (without even attempting decryption, ideally, to save CPU) anything already seen or too far outside the window, exactly the guard this chapter's `openPacket` currently lacks.

---

## 18. Exercises

### Easy
1. Add a byte counter that logs total bytes sealed and opened per minute, giving this VPN a crude throughput readout.
2. Change `-key` handling to read the passphrase from an environment variable instead of a command-line flag, and explain in a comment why that's a meaningful (if partial) security improvement.
3. Trace, by hand, the exact byte length of the UDP envelope for a 64-byte inner ICMP packet, showing your arithmetic against Section 7's table.

### Medium
4. Add the counter-based replay protection described in Section 17's advanced answer, rejecting any envelope whose embedded counter has already been seen.
5. Extend `VPNPeer` to support multiple peers, each identified by its own `net.UDPAddr`, routing outbound packets to the correct peer based on the inner packet's destination IP (a minimal version of WireGuard's `AllowedIPs`).
6. Add a `-mtu` flag that logs a warning whenever an inbound TUN packet's size, plus the 28-byte crypto overhead and 28-byte outer UDP/IP headers, would exceed a configurable physical link MTU.

### Hard
7. Replace `deriveKey`'s static passphrase with a minimal ephemeral X25519 (Chapter 79's ECC) key exchange run once at startup between the two peers over the same UDP socket, deriving a fresh session key each run — closing this chapter's forward-secrecy gap.
8. Implement key rotation: renegotiate a new AES key every N minutes without dropping in-flight packets, requiring both the old and new key to be briefly valid for decryption simultaneously.
9. Port Section 6's `openTUN` to also support macOS's `utun` interface (a different, BSD-style ioctl interface), behind the same `io.ReadWriteCloser` abstraction, so Section 8-10's code runs unmodified on both platforms.

---

## 19. Summary

| Term | Meaning |
|---|---|
| TUN device | A virtual network interface that hands a user-space program raw IP packets instead of routing them physically |
| `IFF_TUN` vs `IFF_TAP` | Layer-3 IP packets vs. layer-2 Ethernet frames delivered by a virtual interface |
| `IFF_NO_PI` | Flag requesting packets with no extra 4-byte protocol-info prefix |
| Envelope | This chapter's wire format: `nonce (12B) \|\| ciphertext \|\| tag (16B)` |
| AES-256-GCM | The AEAD cipher used to seal/open each packet — confidentiality and integrity in one operation |
| Peer roaming | Learning a peer's UDP address dynamically from its first validly-decrypted packet |
| Forward secrecy | The property (missing here, present in WireGuard) that a leaked key doesn't expose past traffic |
| Replay protection | The property (missing here) that a captured, valid envelope can't be successfully resent later |

Chapter 117 turned Chapter 85's tunneling diagram into a program you can actually run: a real TUN device, real AES-GCM encryption, and a real UDP socket carrying the result, with an honest accounting of exactly what a production tunnel protocol like WireGuard adds on top.

---

## 20. What's Next in This Volume

Chapter 118 closes out Volume 17 with its capstone: a small distributed key-value service, with multiple nodes discovering each other over the network and replicating writes between themselves through a custom wire protocol — bringing together the sockets, serialization, and concurrency patterns built across every chapter in this volume, from Chapter 106's first accept loop through this chapter's encrypted tunnel.
