# Chapter 32: Networking in the OS

> **"The network stack is one of the most complex parts of the operating system. It must handle millions of connections, validate every packet, buffer data efficiently, and provide a simple API (sockets) that hasn't fundamentally changed since BSD Unix in 1983. Getting it right means the difference between a 1Gbps and a 100Gbps server."**

---

## Table of Contents

1. [The Network Stack — Layers](#1-the-network-stack--layers)
2. [Sockets — The OS Network API](#2-sockets--the-os-network-api)
3. [The TCP/IP Stack in the Kernel](#3-the-tcpip-stack-in-the-kernel)
4. [Packet Journey — Receive Path](#4-packet-journey--receive-path)
5. [Packet Journey — Send Path](#5-packet-journey--send-path)
6. [Socket Buffers (sk_buff)](#6-socket-buffers-sk_buff)
7. [TCP Flow Control and Congestion](#7-tcp-flow-control-and-congestion)
8. [Zero-Copy Networking](#8-zero-copy-networking)
9. [Network Namespaces — OS-Level Isolation](#9-network-namespaces--os-level-isolation)
10. [Summary](#summary)

---

## 1. The Network Stack — Layers

The OS implements the **TCP/IP model** (a simplified version of OSI):

```
┌──────────────────────────────────────────────────┐
│ Application Layer         HTTP, DNS, SSH, SMTP   │  User space
├──────────────────────────────────────────────────┤
│ Transport Layer           TCP, UDP               │  Kernel
├──────────────────────────────────────────────────┤
│ Network Layer             IP, ICMP, ARP          │  Kernel
├──────────────────────────────────────────────────┤
│ Link Layer                Ethernet, Wi-Fi        │  Kernel + Driver
├──────────────────────────────────────────────────┤
│ Physical Layer            Cables, radio waves    │  Hardware
└──────────────────────────────────────────────────┘

Data encapsulation (sending):
  Application data
  → + TCP header (ports, seq#, flags)   → TCP segment
  → + IP header (src/dst IP, TTL, proto) → IP packet
  → + Ethernet header (src/dst MAC)      → Ethernet frame
  → sent as electrical/optical/radio signal
```

**Each layer adds a header:**
```
[Ethernet HDR | IP HDR | TCP HDR | HTTP Data]
              │         │         │
              └── 14B   └── 20B   └── variable
```

**On receive, each layer strips its header** and passes the payload up.

---

## 2. Sockets — The OS Network API

A **socket** is the OS abstraction for a network connection. It works like a file — identified by a file descriptor, read/written with `read()`/`write()`.

**Creating a TCP server:**
```c
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

// 1. Create socket:
int sock = socket(AF_INET, SOCK_STREAM, 0);
// AF_INET = IPv4, SOCK_STREAM = TCP (reliable, ordered, bidirectional)
// AF_INET6 = IPv6, SOCK_DGRAM = UDP (unreliable, unordered)
// AF_UNIX = Unix domain sockets (local IPC)

// 2. Bind to address/port:
struct sockaddr_in addr = {
    .sin_family = AF_INET,
    .sin_port   = htons(8080),      // big-endian port number
    .sin_addr.s_addr = INADDR_ANY,  // listen on all interfaces
};
bind(sock, (struct sockaddr*)&addr, sizeof(addr));

// 3. Listen for connections:
listen(sock, 128);  // 128 = backlog (pending connections queue depth)

// 4. Accept a connection (blocks until client connects):
int client_fd = accept(sock, NULL, NULL);
// client_fd is a NEW socket for communicating with this client
// sock continues listening for more connections

// 5. Read/write:
char buf[1024];
ssize_t n = read(client_fd, buf, sizeof(buf));
write(client_fd, "HTTP/1.1 200 OK\r\n...", ...);

// 6. Close:
close(client_fd);
close(sock);
```

**Socket states (TCP):**
```
State transition (server side):
  CLOSED → LISTEN (after bind+listen)
  LISTEN → SYN_RECV (received SYN from client)
  SYN_RECV → ESTABLISHED (3-way handshake complete)
  ESTABLISHED → CLOSE_WAIT (received FIN from client)
  CLOSE_WAIT → LAST_ACK (sent FIN to client)
  LAST_ACK → CLOSED (received ACK of our FIN)
  
  ESTABLISHED → FIN_WAIT1 (we initiated close)
  FIN_WAIT1 → FIN_WAIT2 (received ACK of our FIN)
  FIN_WAIT2 → TIME_WAIT (received FIN from peer)
  TIME_WAIT → CLOSED (after 2×MSL = 60-120 seconds)
```

**TIME_WAIT** exists to handle delayed packets from the old connection. After closing, the OS keeps the socket in TIME_WAIT for 2 minutes to absorb any stray packets.

```bash
# View socket states:
ss -tan
# State      Recv-Q  Send-Q  Local Address:Port  Peer Address:Port
# LISTEN     0       128     0.0.0.0:8080        0.0.0.0:*
# ESTABLISHED 0      0       192.168.1.5:8080    10.0.0.1:54321
# TIME-WAIT  0       0       192.168.1.5:8080    10.0.0.1:54320
```

---

## 3. The TCP/IP Stack in the Kernel

**Key kernel data structures for a TCP socket:**

```c
struct socket {             // VFS-visible abstraction
    socket_state state;
    struct sock *sk;        // protocol-specific data
    const struct proto_ops *ops;  // send, recv, accept, etc.
    struct file *file;      // back-pointer to the open file
};

struct sock {               // generic socket for all protocols
    struct sk_buff_head sk_receive_queue;  // received data waiting for recv()
    struct sk_buff_head sk_write_queue;    // data queued for sending
    int sk_rcvbuf;          // receive buffer size (bytes)
    int sk_sndbuf;          // send buffer size (bytes)
    // ...
};

struct tcp_sock {           // TCP-specific (inherits from sock)
    // Sequence numbers:
    u32 snd_nxt;            // next seq# to send
    u32 rcv_nxt;            // next seq# expected to receive
    u32 snd_una;            // oldest unacknowledged seq#
    
    // Congestion control:
    u32 snd_cwnd;           // congestion window (how many bytes we can send)
    u32 ssthresh;           // slow start threshold
    
    // RTT estimation:
    u32 srtt_us;            // smoothed round-trip time
    u32 mdev_us;            // RTT deviation
    
    // Flow control:
    u32 rcv_wnd;            // receive window advertised to peer
    // ...
};
```

---

## 4. Packet Journey — Receive Path

**A TCP packet arrives at the network card:**

```
1. NIC hardware:
   - Packet arrives on wire (electrical/optical signal)
   - NIC receives Ethernet frame
   - NIC computes checksum (offloaded from CPU)
   - NIC performs DMA: writes frame to kernel ring buffer (RX ring)
   - NIC raises hardware interrupt (or signals via polling in NAPI mode)

2. Hard IRQ handler (top half):
   - Acknowledge interrupt to NIC
   - Schedule NET_RX_SOFTIRQ (deferred processing)

3. NET_RX_SOFTIRQ (bottom half):
   - Call driver's napi_poll() to drain the RX ring
   - For each packet in ring:
     a. Allocate sk_buff (socket buffer)
     b. Copy packet data into sk_buff
     c. Pass to netif_receive_skb()

4. Network layer dispatch (netif_receive_skb):
   - Deliver to registered packet type handlers
   - For Ethernet frame with ethertype 0x0800 (IPv4): → ip_rcv()
   
5. IP layer (ip_rcv):
   - Validate IP header (version, header length, checksum)
   - Handle fragmentation (reassemble if fragmented)
   - Check destination IP → for us? or forward?
   - Dispatch by protocol:
     - TCP (proto 6): → tcp_v4_rcv()
     - UDP (proto 17): → udp_rcv()
     - ICMP (proto 1): → icmp_rcv()

6. TCP layer (tcp_v4_rcv):
   - Find matching socket: (src_ip, src_port, dst_ip, dst_port) → sock lookup
   - Validate TCP header (checksum, sequence numbers)
   - Handle TCP state machine (SYN/SYN-ACK/ACK for handshake)
   - Process data:
     a. In-order: add to sk_receive_queue
     b. Out-of-order: store in reorder queue
   - Send ACK (may be delayed for piggybacking)
   - Wake up process blocked in recv() if data is now available

7. Application:
   - recv(fd, buf, len) → copies from sk_receive_queue to user buffer
   - Returns number of bytes read
```

---

## 5. Packet Journey — Send Path

**Application calls send():**

```
1. Application:
   send(fd, data, len, 0)  // or write(fd, data, len)
   
2. TCP layer:
   - Copy data from user space to sk_buff
   - Add TCP header (port numbers, sequence numbers, checksum)
   - Add to send queue
   - Check send window (can we send now?)
   - If window allows → pass to IP layer
   
3. IP layer:
   - Add IP header (src/dst IP, TTL=64, proto=TCP, checksum)
   - Route lookup: which interface/gateway to use?
     (consult routing table: ip route show)
   - For remote destination: look up next-hop MAC via ARP
   
4. ARP (if needed):
   - Check ARP cache: (192.168.1.1) → (aa:bb:cc:dd:ee:ff)?
   - If not cached: send ARP broadcast "who has 192.168.1.1?"
   - Wait for ARP reply, cache the MAC
   
5. Link layer:
   - Add Ethernet header (src MAC, dst MAC, ethertype 0x0800)
   
6. NIC driver:
   - Place sk_buff into TX ring buffer
   - Signal NIC to transmit
   - NIC performs DMA: reads frame from TX ring
   - NIC transmits frame on wire
   - NIC raises interrupt when transmission complete
   
7. TX completion interrupt:
   - Free sk_buff
   - Wake up sender if send buffer was full
```

---

## 6. Socket Buffers (sk_buff)

The **sk_buff (socket buffer, skb)** is the central data structure of Linux networking. It represents a single packet plus its metadata.

```c
struct sk_buff {
    // Packet data:
    unsigned char   *head;   // start of allocated buffer
    unsigned char   *data;   // start of current data (packet content)
    unsigned char   *tail;   // end of current data
    unsigned char   *end;    // end of allocated buffer
    
    // Header areas (before data):
    // As packet travels DOWN the stack, headers are PREPENDED (pushed before data):
    //   head → [headroom][data: [TCP][IP][Ethernet][payload]][tailroom] ← end
    
    // Metadata:
    struct sock     *sk;     // socket this packet belongs to
    struct net_device *dev;  // network device
    unsigned short   protocol; // Ethernet type
    char             ip_summed; // checksum state (NONE, PARTIAL, COMPLETE)
    
    // TCP/IP info:
    struct skb_shared_info *shinfo; // fragmentation info, timestamp
    
    // Navigation helpers:
    // skb->data points to the current layer's header
    // skb_pull(skb, sizeof(struct ethhdr)) → advances data past Ethernet header
    // skb_push(skb, sizeof(struct iphdr))  → prepends IP header before data
};
```

**Zero-copy architecture:**
sk_buff has a headroom before `data` specifically so that each protocol layer can **prepend** its header without copying the payload:
```
Initial (transport layer):   [    headroom    ][TCP header][payload]
After IP pushes header:      [ headroom][IP header][TCP header][payload]
After Ethernet pushes header:[Ethernet][IP header][TCP header][payload]
                              ^data points here
```

No copies of the payload at all — just pointer manipulation!

---

## 7. TCP Flow Control and Congestion

**Flow control (receiver-based):**
Prevents the sender from overwhelming the receiver's buffer.

```
Receiver: "My receive buffer has 65535 bytes available"
→ Sends this as window size in every TCP ACK header

Sender: only send up to window_size bytes without receiving ACK
→ If receiver is slow, window shrinks → sender slows down
→ If receiver is overwhelmed: window = 0 → sender pauses completely

Modern: TCP window scaling (RFC 1323) — window up to 1GB
  (default 65535 bytes is too small for high-latency links)
```

**Congestion control (network-based):**
Prevents overwhelming the network (not just the receiver).

**Slow Start:**
```
Start: cwnd (congestion window) = 1 MSS (1 MTU-worth, ~1460 bytes)
For each ACK received: cwnd += 1 MSS  (doubles each RTT!)
1 RTT → cwnd = 2 MSS
2 RTT → cwnd = 4 MSS
3 RTT → cwnd = 8 MSS
...until cwnd reaches ssthresh (slow start threshold, initially 64KB)
```

**Congestion Avoidance:**
After reaching ssthresh: cwnd += 1 MSS per RTT (linear growth, not exponential)

**Packet loss detected:**
```
Triple duplicate ACK (fast retransmit):
  ssthresh = cwnd / 2
  cwnd = ssthresh
  Enter fast recovery

Timeout (severe congestion):
  ssthresh = cwnd / 2
  cwnd = 1 MSS
  Restart slow start (full backoff)
```

**Modern congestion control algorithms:**
```bash
# View available congestion control algorithms:
cat /proc/sys/net/ipv4/tcp_available_congestion_control
# cubic reno bbr2

# Change to BBR (Google's algorithm, better for high BDP links):
echo bbr > /proc/sys/net/ipv4/tcp_congestion_control

# CUBIC: default in Linux, good for local networks
# BBR: model-based, better for links with bufferbloat
# QUIC: Google's UDP-based transport, used by HTTP/3
```

---

## 8. Zero-Copy Networking

**Normal send path:**
```c
// Application: 4 copies of data occur!
// 1. Application buffer → kernel space (copy_from_user in write())
// 2. Kernel buffer → NIC DMA (NIC reads from kernel buffer)
// (Plus 2 mode switches: user→kernel for write, kernel→user for nothing here)
```

**sendfile() — zero-copy for file serving:**
```c
// Send a file to a socket WITHOUT copying to user space:
int file_fd = open("large_file.mp4", O_RDONLY);
int sock_fd = ...; // connected TCP socket

// Sends file_fd contents directly to sock_fd:
sendfile(sock_fd, file_fd, NULL, file_size);

// Data path:
// Page cache → NIC DMA (no CPU copying, no user space involvement)
// Used by nginx, Apache for static file serving — massive throughput improvement
```

**splice() — pipe-based zero-copy:**
```c
// Move data between two file descriptors through a pipe (no copy to user space):
splice(src_fd, NULL, pipe_fd[1], NULL, len, SPLICE_F_MOVE);
splice(pipe_fd[0], NULL, dst_fd, NULL, len, SPLICE_F_MOVE);
```

**io_uring — async I/O with zero syscall overhead:**
```c
// Submit batch of operations without context switches:
struct io_uring ring;
io_uring_queue_init(256, &ring, 0);

struct io_uring_sqe *sqe = io_uring_get_sqe(&ring);
io_uring_prep_read(sqe, fd, buf, len, 0);  // prepare read
io_uring_submit(&ring);  // submit all prepared operations in one syscall

// ... do other work ...

struct io_uring_cqe *cqe;
io_uring_wait_cqe(&ring, &cqe);  // wait for completion
int result = cqe->res;
io_uring_cqe_seen(&ring, cqe);
```

---

## 9. Network Namespaces — OS-Level Isolation

Linux **network namespaces** give each process group its own isolated network stack:
- Separate network interfaces
- Separate routing tables
- Separate firewall rules (iptables/nftables)
- Separate socket state

**This is how Docker containers and Kubernetes pods get network isolation.**

```bash
# Create a network namespace:
ip netns add mycontainer

# Create a virtual Ethernet pair (veth = virtual wire):
ip link add veth0 type veth peer name veth1

# Move one end into the namespace:
ip link set veth1 netns mycontainer

# Configure outside end:
ip addr add 10.0.0.1/24 dev veth0
ip link set veth0 up

# Configure inside end:
ip netns exec mycontainer ip addr add 10.0.0.2/24 dev veth1
ip netns exec mycontainer ip link set veth1 up
ip netns exec mycontainer ip link set lo up

# Run a command in the namespace:
ip netns exec mycontainer ping 10.0.0.1
# (reaches veth0 — can communicate with the outside)

# Enable forwarding between namespaces and the internet:
echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE
```

**Other namespaces used by containers:**
```
Network namespace:  separate network stack
PID namespace:      separate process tree (container sees PID 1 as init)
Mount namespace:    separate file system view (container's own root FS)
UTS namespace:      separate hostname/domainname
IPC namespace:      separate System V IPC, POSIX MQs
User namespace:     separate UID/GID mapping (root inside = non-root outside)
Cgroup namespace:   separate cgroup hierarchy view
```

---

## Summary

| Concept | Description |
|---------|------------|
| Socket | File descriptor for a network connection; unified API for TCP/UDP/Unix domain |
| TCP three-way handshake | SYN → SYN-ACK → ACK; establishes connection state |
| TIME_WAIT | 2-minute post-close state; absorbs stray packets |
| sk_buff | Kernel packet structure; supports header prepend without data copy |
| netif_receive_skb | Kernel function dispatching received packets to protocol handlers |
| ip_rcv | IP layer handler; validates, reassembles, routes packets |
| tcp_v4_rcv | TCP layer; finds socket, validates sequence numbers, delivers data |
| sk_receive_queue | Per-socket buffer of received data waiting for recv() |
| Slow start | TCP ramps up from 1 MSS, doubling cwnd per RTT until ssthresh |
| Congestion avoidance | Linear cwnd growth after slow start; backs off on loss |
| BBR | Google's model-based TCP congestion control; better for high-BDP links |
| sendfile() | Zero-copy: sends file to socket directly from page cache |
| io_uring | Async I/O interface; batches operations to minimize syscall overhead |
| Network namespace | Isolated network stack per process group; used by containers |
| VETH pair | Virtual Ethernet cable connecting two network namespaces |
| ARP | Address Resolution Protocol; maps IP address to MAC address |
| Window scaling | TCP extension allowing windows up to 1GB (needed for fast long-distance links) |
