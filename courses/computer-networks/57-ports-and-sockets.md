# Chapter 57: Ports and Sockets — How Many Programs Share One Network Connection

> **"An IP address gets a packet to the right building. A port gets it to the right room."**

---

## Table of Contents

1. [The Problem: One Address, Many Programs](#1-the-problem-one-address-many-programs)
2. [A Naive First Attempt](#2-a-naive-first-attempt)
3. [The Real Solution: Ports](#3-the-real-solution-ports)
4. [Well-Known, Registered, and Ephemeral Ports](#4-well-known-registered-and-ephemeral-ports)
5. [The Socket — Identifying One Conversation](#5-the-socket--identifying-one-conversation)
6. [The Four-Tuple, Worked](#6-the-four-tuple-worked)
7. [How the OS Actually Demultiplexes Packets](#7-how-the-os-actually-demultiplexes-packets)
8. [The Socket API, Briefly](#8-the-socket-api-briefly)
9. [Seeing It Live: `ss` and `netstat`](#9-seeing-it-live-ss-and-netstat)
10. [Hands-On Experiment: Watch Ports Get Claimed and Freed](#10-hands-on-experiment-watch-ports-get-claimed-and-freed)
11. [Common Misconceptions](#11-common-misconceptions)
12. [Production Notes](#12-production-notes)
13. [What's Simplified Here](#13-whats-simplified-here)
14. [Interview Questions & Model Answers](#interview-questions--model-answers)
15. [Exercises](#exercises)
16. [Summary](#summary)

---

## 1. The Problem: One Address, Many Programs

By Chapter 45, you know how a packet finds its way across the Internet: routers look at the destination IP address, do a longest-prefix match, and forward the packet one hop closer to its destination. By the time a packet arrives at your laptop, that job is done. The packet is *here*.

But "here" is not good enough. Your laptop, right now, probably has:

- A browser talking to three or four websites at once
- An email client polling a mail server
- A music app streaming audio
- A messaging app holding a connection open in the background
- An SSH session to a remote server

All of these programs share **one network interface** and **one IP address**. When a packet lands on that interface carrying data for a web page, how does the operating system know it's for the browser and not the email client? IP headers (Chapter 36) only carry source and destination *addresses* — they say nothing about which *program* on the destination machine should receive the data. IP's job stops at the machine. Something else has to finish the delivery.

This is the same kind of problem you've already seen solved twice in this course:

- Chapter 29 solved "which machine on this Ethernet segment?" with MAC addresses.
- Chapter 36 solved "which machine on this planet?" with IP addresses.

Now we need: "which **program** on this machine?" That is the problem this chapter solves, and it's the first problem the *transport layer* — the layer TCP and UDP live at — exists to answer.

---

## 2. A Naive First Attempt

Before reaching for the real answer, it's worth trying the obvious first ideas and watching them fail — the same derivation habit this whole course uses.

**Naive idea #1: give every program its own IP address.**

If your browser had IP address 192.168.1.10 and your email client had 192.168.1.11, a packet's destination IP alone would tell the OS which program to hand it to. This actually sort of works — and cloud engineers do something like it deliberately (one pod, one IP, in Kubernetes, as you'll see in Chapter 104). But as a *general* solution it collapses immediately:

- Your machine would need as many IP addresses as programs that might ever want to talk on the network — dozens, sometimes.
- IPv4 addresses are already scarce (Chapter 42 exists because of this).
- Every time you opened a new program, something would need to assign it a fresh address and tell the rest of the Internet how to route to it. That's absurd for a program that lives for ten seconds.

**Naive idea #2: let the operating system just guess, or hand every packet to every program.**

Also broken: it means every program has to inspect every packet arriving on the machine to see if it's "for them," wasting CPU, and worse, leaking data — your email client would see your browser's traffic.

**What we actually need** is a second, smaller address that rides *inside* the same IP packet, chosen freely by each program, requiring no coordination with any router or ISP, needing no global registration — because it only has to be meaningful on the one machine that receives it. That's a port.

---

## 3. The Real Solution: Ports

A **port** is a 16-bit number — a value from 0 to 65535 — that identifies a specific program (technically, a specific communication endpoint of a process) on a given machine.

Every TCP or UDP packet carries two of these numbers in its header:

- A **source port** — which program on the sending machine this data came from
- A **destination port** — which program on the receiving machine this data is for

```
                 IP header                    Transport header (TCP/UDP)
        ┌───────────────────────────┐   ┌──────────────────────────────┐
        │ src IP: 203.0.113.5       │   │ src port: 51342               │
        │ dst IP: 142.250.80.46     │   │ dst port: 443                 │
        └───────────────────────────┘   └──────────────────────────────┘
                                                        │
                                                        ▼
                                    142.250.80.46's kernel:
                                    "port 443 → hand this to the
                                     process listening there (nginx)"
```

Sixteen bits gives 2^16 = 65,536 possible port numbers per IP address, per transport protocol. That's why a single web server, listening on port 443, can be reached by tens of thousands of different browsers simultaneously without a single collision on the server's *destination* port — the server side of every connection uses the *same* port 443, and it's the combination with each client's unique address and port that tells connections apart (Section 6 makes this precise).

Note the phrase "per transport protocol." TCP port 53 and UDP port 53 are two entirely separate numbers that happen to share a decimal value — the OS keeps two independent port tables, one for TCP and one for UDP, because the two header formats and the two socket types are unrelated. DNS actually uses both: UDP/53 for normal queries, and TCP/53 for large or zone-transfer responses.

It's also worth being precise about what ports have nothing to do with: IPv4 versus IPv6 (Chapter 42). The port field itself is identical in size and meaning regardless of which IP version carries it — a server listening on TCP port 443 over IPv4 and one listening on TCP port 443 over IPv6 are, from the port's point of view, doing exactly the same thing with exactly the same 16-bit number. What changes is only the address half of the eventual 4-tuple (Section 6): an IPv6 socket's tuple looks like `(2001:db8::5, 51342, 2606:4700::1, 443)` instead of using dotted-decimal addresses, but the port arithmetic, the well-known/ephemeral ranges, and the demultiplexing logic in Section 7 are completely unchanged. Many servers today run **dual-stack** (Chapter 43), binding the same port number on both an IPv4 and an IPv6 socket simultaneously — two entirely separate sockets, in fact, since the address families are different, even though they share a port number and often even the same underlying application logic.

---

## 4. Well-Known, Registered, and Ephemeral Ports

If ports were assigned totally at random, you couldn't reliably connect to "the web server" on a machine you'd never talked to before — you wouldn't know which port to try. So IANA (the Internet Assigned Numbers Authority) maintains a registry that divides the 65,536 possible ports into three ranges:

| Range | Name | Purpose |
|---|---|---|
| 0 – 1023 | **Well-known ports** | Reserved for standard, widely-used services. Binding to one usually requires administrator/root privileges on Unix-like systems. |
| 1024 – 49151 | **Registered ports** | Vendors and application authors can register a specific port for their software (e.g., 3306 for MySQL, 5432 for PostgreSQL, 8080 as a common "alternate HTTP" convention). No special privilege needed to use these. |
| 49152 – 65535 | **Dynamic / private / ephemeral ports** | Never permanently assigned to anything. The operating system hands these out temporarily to client-side connections. |

Some well-known ports you have already met in this course, or will meet soon:

| Port | Protocol | Service |
|---|---|---|
| 20/21 | TCP | FTP (data/control) |
| 22 | TCP | SSH |
| 23 | TCP | Telnet |
| 25 | TCP | SMTP (mail sending) |
| 53 | UDP/TCP | DNS (Chapter 66) |
| 67/68 | UDP | DHCP (Chapter 55) |
| 80 | TCP | HTTP (Chapter 71) |
| 123 | UDP | NTP (time sync) |
| 443 | TCP | HTTPS / TLS (Chapter 82) |
| 445 | TCP | SMB (Windows file sharing) |

**Ephemeral ports, worked through:** when your browser connects *out* to google.com's port 443, it doesn't use port 443 on your side too. Instead, your operating system's networking stack picks a temporary, unused port number from the ephemeral range for the *source* side of that connection — say, 51342 — and that's the number stamped into the source port field of every packet your browser sends for that conversation. When the browser closes the connection, that port is freed for reuse.

Different operating systems use slightly different ephemeral ranges in practice:

```
IANA-recommended:  49152 – 65535   (~16,384 ports)
Linux default:      32768 – 60999   (~28,000 ports; see /proc/sys/net/ipv4/ip_local_port_range)
Windows (modern):   49152 – 65535
```

This is why the number 49152 shows up so often: it's 0xC000, chosen because it's exactly three-quarters of the way through the 16-bit space, leaving a clean quarter for dynamic use — though as you can see, Linux ignores the recommendation and claims a much larger slice for itself.

---

## 5. The Socket — Identifying One Conversation

Ports solve "which program," but there's a subtlety hiding just underneath: a *port number alone* does not identify one specific conversation.

Think about a busy web server. It listens on TCP port 443. At any given moment it might be serving five thousand different browsers, all of whom are sending packets with **destination port 443**. If the server's networking stack only looked at the destination port, it would have no way to tell "this packet belongs to the connection with the laptop in Mumbai" from "this packet belongs to the connection with the phone in São Paulo" — they'd collide into the same bucket.

The piece of information that actually, uniquely identifies one ongoing conversation is called a **socket**, and it is defined by four numbers together, not one:

```
(source IP, source port, destination IP, destination port)
```

This is often just called **the 4-tuple**. Two TCP segments belong to the same connection if and only if all four values match. Change any one of the four — a different client IP, a different client port, a different server, a different server port — and it's a different conversation, tracked separately by both ends' operating systems.

> **Intuitive analogy:** think of an apartment building's mail room. The building's street address is the IP address. Each apartment number is a port. But if two different couriers are *both* delivering a package to apartment 443, the only way to keep their deliveries straight is to also record who each courier is (the source) — "package for 443, from courier A" versus "package for 443, from courier B." The apartment number alone isn't a big enough label once more than one sender is involved.
>
> **Where the analogy breaks:** a real apartment doesn't get a *different* mailroom slot for every courier — it's still "apartment 443" either way. A TCP socket really does behave as if it's a separate, independent mailbox per (courier, apartment) pair, each with its own buffered, ordered stream of bytes, isolated from every other pair using that same apartment number.

---

## 6. The Four-Tuple, Worked

Let's make this concrete with real numbers. Suppose `example-server.com` (142.250.80.46) runs a web server on port 443, and two different laptops connect to it at the same time.

```
Laptop A: 203.0.113.5,  ephemeral port 51342
Laptop B: 198.51.100.9, ephemeral port 60110

Connection 1 (A → server):
  (203.0.113.5, 51342, 142.250.80.46, 443)

Connection 2 (B → server):
  (198.51.100.9, 60110, 142.250.80.46, 443)
```

Both connections share the same destination IP and the same destination port — 142.250.80.46:443 — yet the server's kernel keeps them completely separate, because the *source* half of the tuple differs in both cases. The server can even be serving Laptop A twice at once (two browser tabs open to the same site), and as long as the OS assigned each tab a different ephemeral source port, those are two different sockets too:

```
Tab 1: (203.0.113.5, 51342, 142.250.80.46, 443)
Tab 2: (203.0.113.5, 51343, 142.250.80.46, 443)
```

Notice something important here: **the server's listening port never changes.** Port 443 is not "used up" or "consumed" by each new connection — a single `listen()`ing socket on port 443 spawns a brand-new, independent *connected* socket (with its own 4-tuple) for each accepted TCP connection. This is exactly how one Nginx or Apache process, bound to one port, serves tens of thousands of simultaneous clients: the OS distinguishes them by tuple, not by port alone.

A theoretical upper bound follows directly from this: since the tuple includes a 16-bit source port, a single client machine can hold at most ~65,536 simultaneous connections *to one specific server IP:port pair* before it runs out of distinct source ports to use (this is the real mechanism behind "port exhaustion," covered in Section 11). A server, by contrast, is limited only by memory and file descriptors — its side of the tuple (its own IP:port) is fixed, but the *other* three fields vary per client, so it can track millions of simultaneous sockets across different client IPs and ports.

---

## 7. How the OS Actually Demultiplexes Packets

Putting Chapters 27, 36, and this chapter together, here is the exact path a packet takes once it physically arrives at a machine:

```
1. NIC receives an Ethernet frame → checks destination MAC address (Ch 29)
   matches this machine? keep it, else discard.
2. Kernel strips the Ethernet header → looks at EtherType → hands
   payload to the IP layer.
3. IP layer checks destination IP address (Ch 36) → matches this
   machine? keep it, else (if forwarding is enabled) route it onward.
4. IP layer reads the Protocol field in the IP header (6 = TCP, 17 = UDP)
   → hands payload to the matching transport-layer handler.
5. Transport layer (TCP or UDP) reads source port + destination port
   from its own header, combines them with the source/destination IP
   already known from step 3, and looks up the full 4-tuple in its
   connection table.
6. Match found → deliver the payload to that socket's receive buffer,
   which the owning process reads from with a system call (read/recv).
   No match (for TCP, no matching connection and no listening socket
   on that port) → kernel replies with a TCP RST, or for UDP, an ICMP
   "port unreachable" message.
```

Every layer strips its own header and asks one question — "who is this for, at my layer?" — before handing the remainder up to the next. Ports are simply the question the transport layer asks.

---

## 8. The Socket API, Briefly

"Socket" is also the name of the actual programming abstraction operating systems expose for this. You will build real TCP and UDP servers with this API in Chapters 106–107, but here is the shape of it, since the vocabulary (`bind`, `listen`, `accept`, `connect`) maps directly onto everything above:

```go
// Server side — claims a port and waits for connections
listener, _ := net.Listen("tcp", ":443")   // bind() + listen() under the hood
for {
    conn, _ := listener.Accept()           // blocks until a client connects;
                                            // returns a NEW socket, one per client,
                                            // each with its own 4-tuple
    go handleConnection(conn)
}

// Client side — picks an ephemeral source port automatically
conn, _ := net.Dial("tcp", "142.250.80.46:443") // connect()
```

`listener.Accept()` is the crucial detail: it does not hand you the *listening* socket again and again. Each call returns a brand-new `conn` object, bound to a specific 4-tuple, the moment a new client's handshake (Chapter 59) completes. The listening socket's only job is to sit on the well-known port and mint new connected sockets.

Underneath this Go code, on a Unix-like OS, a socket is ultimately just a **file descriptor** — the same small integer handle used for regular files, pipes, and devices. This is a deliberate design decision, not an implementation accident: it means the same `read()`/`write()` system calls, the same `select`/`epoll`/`kqueue` event-notification machinery, and the same permission model that already existed for files could be reused for network communication, rather than inventing a parallel universe of network-only primitives. When a program hits its file descriptor limit (`ulimit -n` on Linux), it is, among other things, hitting a ceiling on how many simultaneous sockets it can hold open — a very real production constraint for servers handling large numbers of concurrent connections (Chapter 106 builds on this directly when constructing a real TCP server).

One more API detail worth knowing by name because it resurfaces in Section 11's misconceptions: `SO_REUSEPORT`. Ordinarily, only one socket may be bound to a given IP:port combination. `SO_REUSEPORT` relaxes this, letting several independent sockets — often one per CPU core, in separate goroutines, threads, or even separate processes — all bind the identical address and port simultaneously, with the kernel itself load-balancing each new incoming connection across them (usually by hashing the 4-tuple):

```go
lc := net.ListenConfig{
    Control: func(network, address string, c syscall.RawConn) error {
        var opErr error
        c.Control(func(fd uintptr) {
            opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
        })
        return opErr
    },
}
listener, _ := lc.Listen(context.Background(), "tcp", ":443")
// Four separate processes can each run this and all successfully bind :443 —
// the kernel hands each new connection's SYN to exactly one of them.
```

This is a common way to scale a single-threaded-per-core server design (common in Go, Nginx, and many high-performance servers) across all of a machine's CPU cores without needing a separate load balancer process in front just to spread work between workers on the *same* machine.

---

## 9. Seeing It Live: `ss` and `netstat`

You can watch this 4-tuple bookkeeping directly on your own machine. `ss` (socket statistics, the modern Linux tool — `netstat` is its older, slower ancestor, still common on other systems) prints exactly the structure described above:

```
$ ss -tn
State      Recv-Q  Send-Q   Local Address:Port      Peer Address:Port
ESTAB      0       0        192.168.1.42:51342      142.250.80.46:443
ESTAB      0       0        192.168.1.42:51343      142.250.80.46:443
ESTAB      0       0        192.168.1.42:60110      13.107.42.14:443
LISTEN     0       128      0.0.0.0:22               0.0.0.0:*
```

Read the last line first: `LISTEN` on `0.0.0.0:22` means "sshd is bound to port 22 on every interface, and has no peer yet — it's a listening socket, not a connected one." The three `ESTAB` (established) lines above it are three separate 4-tuples: two connections to the same server (142.250.80.46:443, different local ports — two browser tabs) and one connection to a different server entirely. Every field you've learned in this chapter is right there in the output.

On the server side of a busy web server, the same command would show thousands of `ESTAB` lines all sharing `Local Address:Port` of `0.0.0.0:443` (or the server's real IP), differentiated only by `Peer Address:Port` — a direct, visual confirmation that the destination port is shared but the tuple is not.

---

## 10. Hands-On Experiment: Watch Ports Get Claimed and Freed

Everything above is easy to verify directly, with tools nearly every machine already has. This is a genuine, runnable experiment — try each step.

**Step 1 — see your machine's ephemeral range.**

```
$ cat /proc/sys/net/ipv4/ip_local_port_range
32768   60999
```

(On macOS: `sysctl net.inet.ip.portrange.first net.inet.ip.portrange.last`.) This is the exact range Section 4 described — the pool the OS draws from every time a program calls `connect()` without first binding a specific source port itself.

**Step 2 — open a listener and watch it appear.**

```
$ nc -l 9000 &
$ ss -tln | grep 9000
LISTEN  0  1  0.0.0.0:9000  0.0.0.0:*
```

A single `LISTEN` line appears — no peer address yet, exactly as Section 9 described, because nothing has connected.

**Step 3 — connect twice, and watch two distinct tuples appear.**

```
$ nc 127.0.0.1 9000 &     # first client
$ nc 127.0.0.1 9000 &     # second client
$ ss -tn | grep 9000
ESTAB  0  0  127.0.0.1:52011  127.0.0.1:9000
ESTAB  0  0  127.0.0.1:52013  127.0.0.1:9000
ESTAB  0  0  127.0.0.1:9000   127.0.0.1:52011
ESTAB  0  0  127.0.0.1:9000   127.0.0.1:52013
```

Notice each connection produces *two* lines — one for each end of the socket, since both the client and server processes have their own kernel-tracked socket for it — and the two clients received two different ephemeral source ports (`52011` and `52013`) from the pool in Step 1, which is exactly what lets the listening `nc` process (and the kernel underneath it) tell the two connections apart despite sharing the same destination port.

**Step 4 — close a connection and watch the tuple disappear (eventually).**

```
$ kill %2   # kill the second nc client
$ ss -tn | grep 9000
```

The pair of lines for port `52013` should vanish almost immediately on the client side; on some systems, the *server*-side socket may briefly persist depending on how the peer's close was detected — a small preview of the `TIME_WAIT` behavior Chapter 64 covers in full.

**Step 5 — try to steal a bound port.**

```
$ python3 -m http.server 9000
$ python3 -m http.server 9000   # in a second terminal, same port
OSError: [Errno 98] Address already in use
```

This is Section 11's "a port can only have one process bound and listening on it at a time" rule, observed directly — and it's the exact error `SO_REUSEPORT` (Section 8) is designed to bypass deliberately, when that's actually what you want.

---

## 11. Common Misconceptions

- **"A port is like a physical door that only one thing can use."** Only half right. A port *can* only have one process **bound and listening** on it at a time (per IP, per protocol) — two web servers can't both `bind()` to TCP port 443 on the same address. But a listening port fans out into thousands of simultaneously connected sockets, each with a different tuple, which is precisely what Section 6 walked through. It's less "one door" and more "one lobby with an unlimited number of private rooms behind it, one room assigned to each visitor."
- **"Port 80 is HTTP and port 443 is HTTPS, always."** These are strong conventions enforced by nobody. Nothing stops a server operator from running HTTP on port 8080 or SSH on port 443 (people do this deliberately to get past restrictive firewalls that only allow "web" traffic outward). The port number is a *convention* IANA registers, not a technical requirement enforced by the protocol itself.
- **"UDP doesn't have ports."** It does — Chapter 58 shows the UDP header has source and destination port fields in exactly the same positions as TCP's. Ports are a transport-layer concept shared by both protocols, not a TCP-only feature.
- **"Two processes can never share a port."** Modern kernels support `SO_REUSEPORT`, which lets multiple processes (or threads) bind the *same* IP:port combination simultaneously, with the kernel load-balancing incoming connections across them — a common trick for scaling a server across CPU cores without a separate load balancer in front.
- **"Well-known ports are somehow more secure or reserved by the network."** They're not blessed by any router or ISP — the "reserved" enforcement (needing root/administrator privileges to bind below 1024) is purely a local operating-system convention, meant to stop an unprivileged user process from impersonating a system service like SSH or HTTP on a shared machine.

---

## 12. Production Notes

- **Ephemeral port exhaustion** is a real production incident, not a theoretical curiosity. A server or client making very high rates of outbound connections to the *same* remote IP:port (e.g., a proxy calling one downstream API very heavily) can exhaust its local ephemeral port range, because each outbound connection needs a distinct 4-tuple and the destination side of that tuple never changes. Symptoms: `connect: cannot assign requested address`. Fixes include widening the ephemeral port range, enabling `SO_REUSEADDR`/connection pooling, or spreading outbound traffic across more source or destination addresses.
- **`TIME_WAIT` accumulation** (previewed here, covered fully in Chapter 64) ties directly into this: a closed connection's tuple is held in a lingering state for a couple of minutes to avoid confusing it with a future connection reusing the same numbers — under high connection churn, this can also exhaust available ephemeral ports.
- **Firewalls and security groups** (Chapters 84, 97) are, mechanically, mostly port-matching rules: "allow inbound TCP to port 443, deny everything else." Understanding that a port identifies a *service*, not a *machine*, is why "open port 22 only from these IP ranges" is a meaningful, common rule.
- **NAT** (Chapter 41) has to rewrite not just IP addresses but also port numbers, precisely because many internal machines behind one public IP need to be told apart by an external server — NAT effectively multiplexes many internal 4-tuples onto a shared set of public source ports.
- **Container and Kubernetes port mapping** (Chapters 103–104) is, mechanically, another layer of exactly this same idea: a container might believe it's bound to port 8080 "on its own machine," while the host actually exposes that service to the outside world on a completely different port, with the container runtime rewriting addresses and ports in between — a modern, software-defined echo of what NAT boxes have done in hardware and firmware for decades.
- **Load balancer health checks** (Chapter 95) are frequently implemented as nothing more than "can I open a TCP connection to this backend's port and complete the handshake" — which is why a backend process crashing (and its listening socket disappearing) is detected almost immediately: the very next health-check connection attempt gets an immediate `RST` or connection-refused error instead of completing, since there is no listener left for the OS to hand the SYN to.

---

## 13. What's Simplified Here

This chapter treats "process" and "socket owner" as interchangeable for clarity; in reality a single process can own many sockets, threads within a process can share sockets, and container/namespace boundaries (Chapter 103) can make "which program" a more layered question than presented here. The port ranges given (well-known/registered/ephemeral) are IANA's official convention (RFC 6335); real operating systems are free to (and do) deviate, as shown with Linux's non-standard ephemeral range. Socket-level details like backlog queues, `SO_REUSEPORT` load-balancing behavior, and multi-homed hosts (multiple IP addresses on one machine, each with its own set of independent port spaces) are real and touched on only briefly above.

---

## Interview Questions & Model Answers

**Beginner: What is a port, and why does it need to be 16 bits?**

A port is a number from 0–65535 that identifies which program on a machine a piece of network data is meant for, since an IP address alone only identifies the machine. Sixteen bits is simply the size TCP and UDP's designers chose for the field in their headers — it caps the space at 65,536 values, which is why the well-known/registered/ephemeral ranges all fit inside that ceiling.

**Intermediate: A server listens on port 443. Explain how it can handle 10,000 simultaneous client connections through that single port.**

The listening port isn't consumed per-connection. Each accepted connection is identified by the full 4-tuple — source IP, source port, destination IP, destination port — not by the destination port alone. Since every client has a different source IP and/or source port, the server's kernel can distinguish 10,000 simultaneous conversations that all share the same destination IP:port, by tracking 10,000 distinct tuples in its connection table, and delivering each incoming packet to the matching socket's buffer.

**Advanced: Why can a busy reverse proxy making many outbound connections to the same backend run out of ports, and what are two ways to fix it?**

Every outbound connection from the proxy to one fixed backend IP:port needs a unique 4-tuple, and since the destination half is fixed, uniqueness has to come entirely from the proxy's own (source IP, source port) pair. With one source IP, that caps concurrent connections to that one backend at roughly 65,536 minus the well-known/registered range, and in practice often much lower once `TIME_WAIT` sockets are counted, since a recently closed connection's tuple is unusable for a couple of minutes. Fixes: (1) connection pooling/keep-alive so fewer new connections are opened per unit time, reusing existing sockets instead of exhausting the tuple space; (2) binding outbound connections across multiple source IP addresses, multiplying the available tuple space by the number of source IPs used.

**Advanced: How does `SO_REUSEPORT` allow four worker processes to all bind the same port without violating the 4-tuple uniqueness that made sockets useful in the first place?**

It doesn't violate 4-tuple uniqueness at all — it changes what happens at the moment a new connection is *accepted*, not what identifies a connection once it exists. Each of the four processes calls `bind()` on the identical `IP:port`, and the kernel keeps all four listening sockets registered against that one address. When a new SYN (Chapter 59) arrives, the kernel picks exactly one of the four listening sockets (commonly by hashing the incoming 4-tuple) to hand that specific connection to, and from that point on the resulting connected socket has its own fully unique 4-tuple exactly as described in Section 6 — no two processes ever end up owning the same established connection. `SO_REUSEPORT` is purely about which of several *candidate listeners* gets first refusal on a brand-new connection; it changes nothing about how already-established connections are told apart.

---

## Exercises

### Easy
1. Name the well-known port for HTTP, HTTPS, SSH, and DNS (UDP).
2. Explain, in one sentence, why a client's source port is usually not 80 or 443, but the server's destination port often is.
3. Run `ss -tn` (Linux/macOS with `ss` installed) or `netstat -an` on your own machine and identify one `LISTEN` line and one `ESTABLISHED` line.

### Medium
4. Two browser tabs on the same laptop are both connected to `wikipedia.org:443`. Write out both full 4-tuples using made-up but plausible IP and port numbers, and explain what must differ between them.
5. A server binds to `0.0.0.0:8080`. What does `0.0.0.0` mean here, and how is it different from binding to `127.0.0.1:8080`?
6. Explain why DNS needs both a UDP port 53 and a TCP port 53, rather than just one or the other.

### Hard
7. A load-testing tool opens 70,000 outbound TCP connections per minute, all to the same target IP:port, from a single source IP, and starts failing with "cannot assign requested address." Diagnose why, using the concepts in this chapter, and propose two independent fixes.
8. Explain how `SO_REUSEPORT` lets four separate worker processes all successfully call `bind()` on the same `0.0.0.0:443`, when Section 11 states a listening port can normally only be bound by one process. What has to be true about how the kernel routes incoming connections for this not to break the 4-tuple uniqueness guarantee?
9. A dual-stack server binds TCP port 8443 on both an IPv4 socket and an IPv6 socket. A client connects over IPv6. Write out the full 4-tuple the server's kernel uses to track that connection, and explain why an IPv4 client connecting "at the same time" to the same port number can never collide with it, even though both use port 8443.
10. A container orchestrator maps a container's internal port 3000 to port 31921 on the host machine. Using the demultiplexing steps from Section 7, explain at which step in that pipeline the port number would actually get rewritten, and why the container's own application code never needs to know that 31921 is involved at all.

---

## Summary

| Term | Meaning |
|---|---|
| Port | A 16-bit number (0–65535) identifying which program on a machine a packet is for |
| Well-known ports | 0–1023, reserved for standard services, usually require admin/root to bind |
| Registered ports | 1024–49151, registerable by vendors for specific applications |
| Ephemeral ports | 49152–65535 (or OS-specific range), assigned temporarily to outgoing client connections |
| Socket | One endpoint of a network conversation; identified in practice by the full 4-tuple |
| 4-tuple | (source IP, source port, destination IP, destination port) — uniquely identifies one connection |
| Listening socket | Bound to a port, waiting to spawn new connected sockets — never itself "consumed" by a connection |
| Port exhaustion | Running out of distinct source ports for outbound connections to one fixed destination |

Ports and sockets tell us *who* on each machine is talking. The next question is *what* they're actually allowed to say to each other — and the simplest possible answer, adding almost nothing to raw IP, is UDP. Chapter 58 covers it.
