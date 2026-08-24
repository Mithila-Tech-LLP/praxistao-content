# Chapter 105: eBPF, XDP, TC, and Cilium

> **"Chapter 104 ended on a scaling problem: kube-proxy's iptables rules grow roughly linearly with the number of Services, and every packet pays that cost. This chapter is about the technology that made a genuinely different answer possible — not a faster set of static rules, but actual programs, written by operators, running safely inside the kernel's own packet path, changeable without ever rebooting or patching the kernel itself."**

---

## Table of Contents

1. [Where This Chapter Picks Up](#1-where-this-chapter-picks-up)
2. [The Problem: iptables Doesn't Scale, and the Kernel Can't Be Recompiled Per Cluster](#2-the-problem-iptables-doesnt-scale-and-the-kernel-cant-be-recompiled-per-cluster)
3. [Naive Fix: Write a Kernel Module](#3-naive-fix-write-a-kernel-module)
4. [Why Kernel Modules Are the Wrong Answer Here](#4-why-kernel-modules-are-the-wrong-answer-here)
5. [The Real Solution: A Sandboxed VM Inside the Kernel](#5-the-real-solution-a-sandboxed-vm-inside-the-kernel)
6. [What eBPF Actually Is](#6-what-ebpf-actually-is)
7. [The Verifier: How the Kernel Stays Safe](#7-the-verifier-how-the-kernel-stays-safe)
8. [Maps: How eBPF Programs Keep State and Talk to Userspace](#8-maps-how-ebpf-programs-keep-state-and-talk-to-userspace)
9. [Attachment Points: Where an eBPF Program Actually Runs](#9-attachment-points-where-an-ebpf-program-actually-runs)
10. [XDP: The Earliest Possible Hook](#10-xdp-the-earliest-possible-hook)
11. [TC (Traffic Control): A Later, More Flexible Hook](#11-tc-traffic-control-a-later-more-flexible-hook)
12. [XDP vs TC vs iptables, Compared Directly](#12-xdp-vs-tc-vs-iptables-compared-directly)
13. [Hands-On: Writing and Loading a Minimal XDP Program](#13-hands-on-writing-and-loading-a-minimal-xdp-program)
14. [Hands-On: Observing eBPF Programs on a Running System](#14-hands-on-observing-ebpf-programs-on-a-running-system)
15. [Code: A Userspace Go Loader for an eBPF Program](#15-code-a-userspace-go-loader-for-an-ebpf-program)
16. [Cilium: A Production Kubernetes CNI Built on eBPF](#16-cilium-a-production-kubernetes-cni-built-on-ebpf)
17. [How Cilium Replaces kube-proxy](#17-how-cilium-replaces-kube-proxy)
18. [How Cilium Enforces NetworkPolicy With eBPF](#18-how-cilium-enforces-networkpolicy-with-ebpf)
19. [Cilium's Observability: Hubble](#19-ciliums-observability-hubble)
20. [Common Misconceptions](#20-common-misconceptions)
21. [Production Notes](#21-production-notes)
22. [What's Simplified Here](#22-whats-simplified-here)
23. [Interview Questions & Model Answers](#23-interview-questions--model-answers)
24. [Exercises](#24-exercises)
25. [Summary and Bridge to Chapter 106](#25-summary-and-bridge-to-chapter-106)

---

## 1. Where This Chapter Picks Up

Chapter 104 closed on a specific, concrete scaling limitation: `kube-proxy`'s `iptables` mode evaluates a chain of rules whose length grows with the number of Services in a cluster, and every single packet pays that evaluation cost on its way through the kernel's networking stack (Chapter 102). This chapter explains the technology built specifically to let that packet-processing logic be replaced — not with faster static rules, but with actual custom programs running inside the kernel — and ends with Cilium, a real, widely deployed Kubernetes CNI (Chapter 103) built entirely on it.

---

## 2. The Problem: iptables Doesn't Scale, and the Kernel Can't Be Recompiled Per Cluster

Every mechanism this course has built through Chapter 104 — NAT rules (Chapter 41), bridge forwarding (Chapter 102), `kube-proxy`'s DNAT chains (Chapter 104) — works by configuring **existing, fixed kernel logic** with data: add a rule, add a route, add a forwarding table entry. The kernel code that *interprets* those rules was written once, by kernel developers, and ships as part of the Linux kernel itself.

This creates a hard ceiling. If an operator wants packet-processing logic the kernel doesn't already have — a custom load-balancing algorithm, a firewall rule keyed on application-layer content the kernel's built-in `netfilter` framework doesn't understand, a way to drop malicious packets in nanoseconds before they ever reach a socket — the only traditional options are: petition the Linux kernel maintainers to add that exact feature (a process that can take years, if it happens at all), or write and load a **kernel module**.

---

## 3. Naive Fix: Write a Kernel Module

Linux kernel modules are real, and they do let arbitrary code run inside the kernel, with full access to kernel memory and hardware. A custom kernel module could, in principle, implement exactly the packet-processing logic described in Section 2.

---

## 4. Why Kernel Modules Are the Wrong Answer Here

Kernel modules run with **zero safety net**. A kernel module has the same privilege level as the kernel itself — a null pointer dereference, an infinite loop, or a memory-safety bug in a kernel module doesn't crash one process, it crashes (or corrupts) the entire machine, taking down every container, every pod, and every other workload sharing that kernel (recall Chapter 102's central point: they all share exactly one kernel). Kernel modules also must be compiled against a specific kernel version's internal APIs, which change between kernel releases, making a module built for one kernel version potentially unloadable — or worse, silently broken — on another.

For infrastructure software meant to run across a large, heterogeneous fleet of machines, potentially updated by an operator who is not a kernel developer, "load arbitrary unverified code directly into the most privileged part of the operating system" is not an acceptable engineering trade-off. The naive fix technically works and is exactly why it's rarely chosen: the blast radius of a single bug is the entire machine.

---

## 5. The Real Solution: A Sandboxed VM Inside the Kernel

The actual solution, which grew out of a much older, narrower Linux feature (the classic BPF packet filter that `tcpdump`, previewed in Chapter 102 and detailed in Chapter 119, has used since the 1990s to filter which packets userspace sees), generalizes that idea into something far more powerful: let userspace load small, custom programs into the kernel, but run them inside a **restricted, verified virtual machine embedded in the kernel itself**, not as raw, unchecked native code.

This is **eBPF** (extended Berkeley Packet Filter — the name is now mostly historical; eBPF today is used for far more than packet filtering). The key property that makes it different from Section 3's kernel module: an eBPF program cannot crash the kernel, cannot read arbitrary kernel or user memory it wasn't explicitly given access to, and is guaranteed to terminate — all enforced *before* the program is ever allowed to run, by a component called the **verifier** (Section 7).

---

## 6. What eBPF Actually Is

At the intuitive level: eBPF is like a plugin system for the kernel, the same way a browser extension runs inside a browser's sandbox instead of as a separate native application with full system access — it can observe and influence specific things the host environment exposes, through a controlled interface, without being trusted with everything the host itself can do.

At the engineering level: eBPF programs are written in a restricted subset of C (or increasingly, higher-level frameworks that compile down to it), compiled to a small, RISC-like eBPF bytecode instruction set, loaded into the kernel via the `bpf()` system call, checked by the verifier (Section 7), and then either interpreted or — on virtually all modern production systems — **JIT-compiled to native machine code** for near-native execution speed once loaded.

At the deep-technical level: an eBPF program is not a general-purpose program that can run anywhere in the kernel — it is written for and attached to one specific **hook** (a defined point in the kernel where the kernel already calls out, if a program is attached there), it has a bounded, small stack (512 bytes), a limited, verifier-checked set of eBPF "helper functions" it's allowed to call (rather than open-ended access to arbitrary kernel functions), and communicates with userspace and with other eBPF programs exclusively through **maps** (Section 8), never through arbitrary shared memory.

```
    Userspace                          Kernel

  +---------------+   bpf() syscall  +------------------------+
  | Loader program | ---------------> | Verifier (Section 7)    |
  | (Section 15)   |                 | rejects unsafe programs |
  +---------------+                  +-----------|------------+
                                                  | passes
                                                  v
                                      +------------------------+
                                      | JIT compiler             |
                                      +-----------|------------+
                                                  v
                                      +------------------------+
                                      | Attached at a hook       |
                                      | (XDP, TC, etc. -- Sec 9)  |
                                      +------------------------+
                                                  |
                                      +-----------v------------+
                                      | Maps (Section 8) --       |
                                      | shared state, readable    |
                                      | from userspace too        |
                                      +------------------------+
```

---

## 7. The Verifier: How the Kernel Stays Safe

Before any eBPF program is allowed to run, the kernel's **verifier** performs a static analysis of every possible execution path through the program's bytecode, checking, among other things:

- **No unbounded loops.** Every loop must be provably bounded (or, in modern kernels, use a special bounded-loop helper) — the verifier must be convinced the program will always terminate, so it can never hang the kernel processing one packet forever.
- **No out-of-bounds memory access.** Every memory access is checked against known buffer sizes at verification time — an eBPF program cannot read or write memory it wasn't explicitly given a checked pointer to.
- **No arbitrary kernel function calls.** Only a specific, versioned set of "helper functions" the kernel exposes can be called — an eBPF program can't call just any internal kernel function, closing off most of the attack surface a native kernel module would have.
- **Bounded program size and complexity.** The verifier walks every possible path through the program and gives up (rejecting the program) if it becomes too complex to fully analyze in reasonable time — in practice this limits how large a single eBPF program can usefully be, which is part of why real systems like Cilium (Section 16) use many small, focused programs rather than one large one.

If a program fails any of these checks, the kernel refuses to load it at all — the equivalent of a strict, mandatory code review performed automatically, every single time, before code is ever allowed to execute at kernel privilege. This static, load-time guarantee — provable safety *before* execution, not runtime sandboxing *during* execution — is the specific engineering property that makes eBPF acceptable where Section 4's raw kernel module is not.

---

## 8. Maps: How eBPF Programs Keep State and Talk to Userspace

An eBPF program is stateless between individual invocations unless it explicitly stores state in an **eBPF map** — a kernel-managed key-value data structure (hash tables, arrays, and several more specialized types) that both eBPF programs and userspace processes can read and write, through the same `bpf()` system call interface.

Maps are how:

- A Cilium eBPF program attached to a network interface (Section 9) looks up "which policy applies to this connection" (Section 18) — the policy was written into the map by a userspace Cilium agent, and the in-kernel program just reads it, on every packet, without needing to be reloaded when policy changes.
- Statistics and events collected by an in-kernel eBPF program (packet counts, dropped-packet reasons) get exposed back to userspace tooling like Cilium's Hubble (Section 19) for observability.
- Multiple eBPF programs, potentially attached at different hooks (Section 9), share state — for instance, an XDP program (Section 10) recording a decision that a later TC program (Section 11) reads.

This map-based communication is precisely why eBPF is so effective for the Kubernetes networking problem Chapter 104 raised: policy and Service state can be updated in a map from userspace, instantly visible to already-running in-kernel programs on the very next packet, with no program reload and no `iptables`-style linear rule-chain growth (Section 12).

---

## 9. Attachment Points: Where an eBPF Program Actually Runs

eBPF isn't limited to networking — it can attach to system calls, function entry/exit points (for tracing and profiling), and security hooks. This chapter focuses on the **networking-relevant attachment points**, because they are the ones Cilium (Section 16) uses to replace `iptables`/`kube-proxy` from Chapter 104. The two most important for packet processing, in the order a packet would actually encounter them, are **XDP** and **TC**.

---

## 10. XDP: The Earliest Possible Hook

**XDP (eXpress Data Path)** attaches an eBPF program to the **network driver itself**, running the program on a raw packet **the instant it arrives from the NIC**, before the kernel has allocated its usual `sk_buff` packet-descriptor structure, before any of the standard IP/TCP stack processing Chapter 102 described, and before `iptables`/`netfilter` ever sees it at all.

This is the earliest point in the entire Linux networking stack a program can run, and it exists specifically for use cases where every nanosecond of per-packet overhead matters: an XDP program can inspect a packet and issue one of a small set of extremely cheap verdicts — `XDP_DROP` (discard immediately, at line rate, before the packet costs the kernel anything further — the standard mechanism for DDoS mitigation at the network edge), `XDP_PASS` (let it continue up the normal stack, unmodified), `XDP_TX` (bounce it back out the same interface, useful for load-balancing use cases), or `XDP_REDIRECT` (send it to a different interface or CPU entirely).

```
   NIC receives packet
          |
          v
   +--------------+
   | XDP hook       |  <- earliest possible point (Section 10)
   | (eBPF program) |
   +------|-------|--+
          |       |
    XDP_DROP   XDP_PASS
          |       |
       discarded  v
              normal kernel stack continues
              (sk_buff allocated, IP/TCP processing,
               TC hook (Section 11), netfilter/iptables...)
```

---

## 11. TC (Traffic Control): A Later, More Flexible Hook

**TC (Traffic Control)** is an older Linux subsystem, originally built for queuing and shaping outbound traffic (rate limiting, prioritization), that also supports attaching eBPF programs — at a point **later** than XDP, after the kernel has already allocated its standard `sk_buff` structure, and available on **both** ingress (incoming) and egress (outgoing) traffic, unlike XDP, which is ingress-only.

Because a TC eBPF program runs against a fully-formed `sk_buff`, it has access to more of the kernel's own packet metadata and can perform more complex operations — including some Section 8 map-driven decisions that need context XDP's earlier, more minimal view doesn't yet have. The trade-off is exactly what's implied by running later in the pipeline: somewhat higher per-packet overhead than XDP, in exchange for more flexibility and the ability to act on both directions of traffic.

---

## 12. XDP vs TC vs iptables, Compared Directly

| Mechanism | Where it runs | Typical use | Per-packet cost |
|---|---|---|---|
| `iptables`/`netfilter` (Chapters 41, 102, 104) | Deep in the standard IP stack, after significant kernel processing | General-purpose NAT, firewalling, `kube-proxy`'s DNAT rules | Grows with rule-chain length (Chapter 104, Section 11) |
| TC eBPF (Section 11) | After `sk_buff` allocation, ingress and egress | Cilium's main data path (Section 17), policy enforcement, load balancing | Lower than `iptables` at scale; independent of rule *count* the way linear `iptables` chains are |
| XDP eBPF (Section 10) | At the NIC driver, before `sk_buff` allocation, ingress only | DDoS mitigation, earliest-possible drop/redirect decisions | Lowest possible — the entire point is skipping standard stack overhead |

The critical structural difference from `iptables`: an eBPF program's own execution cost is roughly constant regardless of how much policy state it's enforcing, because policy lookups happen via Section 8's maps (hash-table lookups, effectively O(1)) rather than via sequential evaluation of a growing list of rules. This is precisely why Chapter 104's `iptables`-mode scaling problem — cost growing with Service count — doesn't recur in an eBPF-based data plane: the *number* of Services grows the map, not the number of instructions the eBPF program has to execute per packet.

---

## 13. Hands-On: Writing and Loading a Minimal XDP Program

A minimal XDP program that drops every incoming ICMP packet (i.e., blocks `ping` at the earliest possible point), written in restricted C and compiled with `clang`'s BPF target:

```c
// xdp_drop_icmp.c
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <bpf/bpf_helpers.h>

SEC("xdp")
int xdp_drop_icmp(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data     = (void *)(long)ctx->data;

    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS; // Section 7's bounds-checking discipline,
                          // required or the verifier rejects this program.

    if (eth->h_proto != __constant_htons(ETH_P_IP))
        return XDP_PASS;

    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;

    if (ip->protocol == IPPROTO_ICMP)
        return XDP_DROP; // Section 10's earliest-possible drop verdict.

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
```

```bash
# Compile to eBPF bytecode.
clang -O2 -target bpf -c xdp_drop_icmp.c -o xdp_drop_icmp.o

# Load and attach to a network interface using the ip command
# (iproute2 has built-in XDP support -- no separate loader needed
# for this simple case; Section 15 shows a Go loader for cases
# needing map interaction from userspace).
sudo ip link set dev eth0 xdp obj xdp_drop_icmp.o sec xdp

# Confirm it's attached:
ip link show eth0
# -> ... prog/xdp id 42 tag ...

# Test it -- pings should now silently fail from anywhere else on the network.
ping -c 3 <this-machine-ip>

# Detach it.
sudo ip link set dev eth0 xdp off
```

Every bounds check in the C source (the `if ((void *)(eth + 1) > data_end)` lines) exists specifically to satisfy Section 7's verifier — removing them causes the kernel to reject the program at load time with a verifier error, not a runtime crash, which is precisely the safety property Section 4 said a raw kernel module could never guarantee.

---

## 14. Hands-On: Observing eBPF Programs on a Running System

```bash
# List every eBPF program currently loaded on the system -- on a node
# running Cilium (Section 16), this will show dozens of programs.
sudo bpftool prog list

# Inspect a specific program's attachment point and basic stats.
sudo bpftool prog show id 42

# List eBPF maps (Section 8) currently in use.
sudo bpftool map list

# Dump the contents of a specific map -- e.g., a Cilium policy map,
# showing exactly the key-value state an in-kernel program is
# consulting on every packet.
sudo bpftool map dump id 7
```

`bpftool` (part of the standard Linux kernel tooling) is the direct analogue of Chapter 56's `ip`/`ss` toolbox, purpose-built for inspecting eBPF's own kernel-resident state — a genuinely new category of "what's running in my kernel right now" that didn't exist before eBPF, because before eBPF the kernel's packet-processing logic was fixed at compile time, not something loaded and swapped at runtime.

---

## 15. Code: A Userspace Go Loader for an eBPF Program

Real production tools (including Cilium itself) commonly use the `cilium/ebpf` Go library to load, attach, and interact with eBPF programs from userspace, rather than shelling out to `ip link` as Section 13 did for simplicity. A minimal loader reading map contents from Section 8's mechanism:

```go
package main

import (
	"fmt"
	"log"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

func main() {
	// Load the compiled eBPF object file from Section 13.
	spec, err := ebpf.LoadCollectionSpec("xdp_drop_icmp.o")
	if err != nil {
		log.Fatalf("loading spec: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		// This is where Section 7's verifier rejection would surface,
		// as a Go error, if the program failed a safety check.
		log.Fatalf("verifier rejected program: %v", err)
	}
	defer coll.Close()

	prog := coll.Programs["xdp_drop_icmp"]

	// Attach to a real interface -- the userspace-driven equivalent
	// of Section 13's `ip link set dev eth0 xdp obj ...` command.
	iface := "eth0"
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   prog,
		Interface: ifaceIndex(iface),
	})
	if err != nil {
		log.Fatalf("attaching XDP program: %v", err)
	}
	defer l.Close()

	fmt.Println("XDP program attached to", iface, "- press Ctrl+C to detach and exit")
	select {} // Keep the process alive; detaching happens on exit via defer.
}

func ifaceIndex(name string) int {
	// In real code this resolves via net.InterfaceByName(name).Index;
	// omitted here for brevity.
	return 2
}
```

This is precisely the pattern Cilium's own agent uses at a much larger scale: a long-running userspace Go process that compiles or loads eBPF programs, attaches them to the hooks described in Sections 10–11, and continuously updates their maps (Section 8) as Kubernetes state (Chapter 104's Services, Endpoints, and NetworkPolicies) changes.

---

## 16. Cilium: A Production Kubernetes CNI Built on eBPF

**Cilium** is a CNCF-graduated Kubernetes CNI plugin (satisfying exactly Chapter 103's CNI contract and Chapter 104's Section 3 networking-model requirements) whose entire data plane — pod-to-pod forwarding, Service load balancing, and NetworkPolicy enforcement — is implemented as a coordinated set of eBPF programs attached to TC and, optionally, XDP hooks on every node, instead of `iptables` rules.

Architecturally, Cilium runs a userspace **agent** (`cilium-agent`) as a daemon on every node — the long-running companion process Chapter 103, Section 19's misconception note previewed — which:

- Watches the Kubernetes API for pods, Services, Endpoints, and NetworkPolicy objects (exactly what Chapter 104's Section 19 Go code did, read-only, for one Service).
- Compiles and loads the eBPF programs needed for that node's current pod set.
- Continuously updates eBPF maps (Section 8) with the live state those programs need — pod IP-to-node mappings, Service-to-backend mappings, and compiled policy rules — so the in-kernel programs themselves never need to be reloaded when state changes, only their maps updated.

```
   cilium-agent (userspace, one per node)
        |
        | watches K8s API: pods, Services, Endpoints, NetworkPolicy
        | (Chapter 104's objects)
        v
   Updates eBPF maps (Section 8) -----------------+
                                                   |
   +-----------------------------------------------|--+
   | Node's kernel                                    |
   |  TC eBPF programs (Section 11) attached per pod   |
   |  interface, consulting the maps on every packet   |
   |  for: Service load balancing, policy enforcement, |
   |  pod-to-pod forwarding (replacing iptables)        |
   +---------------------------------------------------+
```

---

## 17. How Cilium Replaces kube-proxy

In its "kube-proxy replacement" mode, Cilium disables `kube-proxy` entirely and implements Chapter 104 Section 11's exact job — translating a Service's ClusterIP into a real backend pod IP — using an eBPF map keyed by ClusterIP, populated with the live list of backend pod IPs (from watching Endpoints, exactly as Chapter 104's Section 19 code did), consulted by a TC eBPF program attached at each pod's veth (Chapter 102, Section 9) interface.

Because that map lookup is a hash-table operation with cost independent of how many Services exist, this is the concrete resolution of Chapter 104's Section 21 production note: Cilium's eBPF-based Service load balancing scales to a very large number of Services without the roughly-linear-in-rule-count cost `iptables` mode pays, and without `ipvs` mode's own additional operational complexity — the direct answer to the scaling problem this chapter opened with in Section 2.

---

## 18. How Cilium Enforces NetworkPolicy With eBPF

Chapter 104, Section 18 left NetworkPolicy enforcement as "up to the CNI plugin." Cilium's answer: compile each NetworkPolicy object's rules into compact, verifier-checked eBPF logic and identity-keyed map entries, then consult that state in the TC eBPF program attached to every pod's network interface, on every single packet, before it's allowed to leave or enter that pod.

Cilium goes one step further than raw IP-based policy by assigning every pod a numeric **security identity** derived from its Kubernetes labels (so pods with `app: backend` all share one identity, regardless of their constantly-changing IP addresses from Chapter 104's Section 8), and keys its policy maps by identity rather than IP — meaning a policy rule doesn't need to be rewritten every time a pod is rescheduled and gets a new IP, only when its labels genuinely change. This directly compounds the ephemeral-IP problem Chapter 104 Section 8 raised for application code with the same underlying fix Section 9 of that chapter used for Services: stop keying anything durable on an IP address that's expected to change.

---

## 19. Cilium's Observability: Hubble

Because every packet decision already flows through Cilium's own eBPF programs, those programs are naturally positioned to also record what they observed — which Cilium exposes via **Hubble**, its built-in observability component. Hubble surfaces, per-flow, exactly which pods (by identity, per Section 18) talked to which other pods, on which ports, and whether a NetworkPolicy allowed or denied the connection — all sourced from eBPF map data (Section 8) and event data emitted directly from the kernel-resident programs, without needing a separate packet-capture pass the way `tcpdump` (Chapter 119) would.

```bash
# A representative Hubble CLI query (assuming Cilium + Hubble installed):
hubble observe --pod frontend --to-pod backend
# TIMESTAMP   SOURCE            DESTINATION       VERDICT   SUMMARY
# 10:03:12    frontend-7f9c8    backend-4d2a1     FORWARDED TCP Flags: SYN
# 10:03:12    backend-4d2a1     frontend-7f9c8    FORWARDED TCP Flags: SYN, ACK
```

This is Chapter 102's `tcpdump -i br0` observability idea (Section 14 of that chapter), generalized cluster-wide and made policy-aware — instead of raw packet bytes on one bridge, Hubble reports Kubernetes-identity-aware flow decisions sourced directly from the same eBPF programs doing the actual forwarding and policy enforcement.

---

## 20. Common Misconceptions

- **"eBPF programs can crash the kernel like a bad kernel module can."** Section 7 was explicit: the verifier statically proves safety properties (bounded execution, memory safety within checked bounds) before a program is ever allowed to load — this is precisely the property that makes eBPF acceptable in a way Section 4's raw kernel module never was.
- **"XDP is strictly better than TC because it's faster."** XDP is faster specifically because it runs earlier and skips more of the kernel's own stack, but that earliness also means less context is available (no `sk_buff`, ingress-only) — TC eBPF exists precisely for logic that genuinely needs the fuller context XDP's minimal view doesn't yet have (Section 11).
- **"Cilium is just Calico with extra branding."** Both are real, CNCF-ecosystem CNI plugins, but Calico's default data plane is `iptables`/routing-based (with an optional eBPF mode added later), while Cilium was built eBPF-first from the start — the architectural difference described in Sections 16–18 is real, not naming.
- **"eBPF is only for networking."** This chapter focused entirely on networking hooks (XDP, TC) because that's the throughline from Chapters 102–104, but eBPF is also heavily used for security monitoring, system call tracing, and performance profiling — networking is eBPF's most famous production use case, not its only one.
- **"Because eBPF programs are 'in the kernel,' they can see and touch anything the kernel can."** Section 6 and Section 7 both stressed the opposite: an eBPF program's access is deliberately limited to what specific helper functions and maps it's been granted, verifier-checked bounds, and the specific data made available at its attachment point — meaningfully less access than a raw kernel module has.

---

## 21. Production Notes

- Migrating an existing cluster from `iptables`-mode `kube-proxy` to Cilium's kube-proxy-replacement mode (Section 17) is a well-trodden but non-trivial operational change — it typically requires kernel version minimums (specific eBPF features were added across kernel releases) and careful rollout, since it changes the node's entire Service data path.
- eBPF program verification time (Section 7) can itself become a real constraint for very large, complex programs — this is part of why production eBPF systems like Cilium favor composing many small, focused programs (tail-called between each other via a mechanism this chapter didn't detail) rather than one large monolithic program.
- Hubble's (Section 19) per-flow visibility is frequently cited as a concrete operational win over `iptables`-based setups, where reconstructing "which policy denied this connection and why" from raw `iptables` counters is materially harder than reading a Hubble flow log with an explicit verdict and policy name attached.
- Not every environment can use XDP's fastest, driver-level mode ("native" XDP) — it requires NIC driver support; environments without it fall back to "generic" XDP (running slightly later, in the kernel's standard receive path) or "offloaded" XDP (running directly on supported NIC hardware, where the eBPF program executes on the network card itself, ahead of even native XDP's driver-level hook).

---

## 22. What's Simplified Here

The verifier's core safety guarantees (Section 7), the XDP/TC attachment-point distinction (Sections 10–11), and Cilium's high-level architecture (Sections 16–18) are accurate and reflect current, real production behavior. Left out for focus: the full eBPF instruction set and its calling conventions; tail calls and how large real-world eBPF programs like Cilium's are actually composed from many smaller ones to work around Section 7's verifier complexity limits; `sockmap`/`sk_msg` eBPF hooks that operate at the socket layer rather than the packet layer, used for some service-mesh acceleration use cases (a further evolution beyond even Chapter 101's sidecar model); and the full depth of Cilium's additional features (its own service-mesh capabilities, multi-cluster networking, and encryption modes) beyond the CNI and Service-replacement core described here. The core idea — the kernel now supports safely loading and swapping custom packet-processing programs at runtime, verified for safety before execution, communicating with userspace through maps — is accurate and is the real foundation of Cilium and of eBPF's broader, fast-growing role in modern networking infrastructure.

---

## 23. Interview Questions & Model Answers

**Beginner: What problem does eBPF solve that a traditional kernel module does not?**
It lets custom code run inside the kernel's privileged execution context without the safety risk a traditional kernel module carries — a verifier statically checks every eBPF program before it's allowed to load, proving properties like bounded execution and memory safety within checked limits, so a buggy eBPF program is rejected at load time rather than crashing or corrupting the whole machine.

**Beginner: What is the difference between XDP and TC as eBPF attachment points?**
XDP attaches at the network driver, running on a raw packet before the kernel allocates its standard packet descriptor and before any normal stack processing — the earliest and cheapest point, ingress-only. TC attaches later, after that descriptor exists, supporting both ingress and egress, with more context available at a somewhat higher per-packet cost.

**Intermediate: Why does an eBPF-based data plane avoid the scaling problem that iptables-mode kube-proxy has as Service count grows?**
`iptables` evaluates rules sequentially, so its cost grows roughly with the number of rules (and thus Services); an eBPF program instead performs a map lookup (Section 8), a hash-table operation whose cost is essentially independent of how many entries the map holds — so adding more Services grows the map's size, not the number of instructions the eBPF program executes per packet.

**Intermediate: What role does the verifier play in making eBPF safe, and what specific properties does it check?**
The verifier statically analyzes an eBPF program's bytecode before allowing it to load, checking that all loops are provably bounded (guaranteeing termination), that all memory accesses stay within checked bounds, and that the program only calls a specific, restricted set of kernel-provided helper functions — rejecting the program at load time if any check fails, rather than allowing potentially unsafe code to run and fail at runtime.

**Advanced: Explain how Cilium's use of security identities, rather than raw pod IPs, for NetworkPolicy enforcement solves a problem introduced back in Chapter 104.**
Chapter 104 established that pod IPs change on every reschedule; if NetworkPolicy enforcement were keyed directly on IP addresses, every pod restart would require rewriting policy state. Cilium instead derives a numeric identity from a pod's Kubernetes labels and keys its eBPF policy maps by that identity — since labels typically stay stable across reschedules even when IPs don't, policy state doesn't need to be touched when a pod is merely recreated with a new IP, only when its actual labels change.

**Advanced: A team wants sub-microsecond packet drop decisions for DDoS mitigation versus a team that wants to enforce complex, stateful, application-aware policy on both inbound and outbound pod traffic. Which eBPF attachment point is the better fit for each, and why?**
The DDoS mitigation case fits XDP (Section 10): it needs the earliest possible, cheapest verdict, and `XDP_DROP` discards malicious packets before the kernel has spent any further processing on them, which is exactly XDP's design point. The stateful, bidirectional policy case fits TC (Section 11): it needs the fuller packet context TC's later attachment point provides, and it needs to run on both ingress and egress, which XDP (ingress-only) cannot do — matching exactly how Cilium itself splits responsibilities between the two hooks in production.

---

## 24. Exercises

### Easy
1. In your own words, explain why a kernel module and an eBPF program both "run inside the kernel" but carry very different risk profiles.
2. List the four example XDP verdicts from Section 10 and give a one-sentence use case for each.
3. Run `sudo bpftool prog list` (or `sudo bpftool prog list` inside a VM if not on Linux) and identify whether any XDP or TC programs are already loaded on the system.

### Medium
4. Modify Section 13's `xdp_drop_icmp.c` to instead drop only UDP packets destined for port 53 (DNS, Chapter 66), and explain what additional header parsing and bounds checks are needed.
5. Using Section 15's Go loader as a model, write a program that reads and prints the contents of an eBPF map by name using the `cilium/ebpf` library's `ebpf.LoadPinnedMap`, explaining what "pinning" a map means and why it's needed for a separate process to access it later.
6. Referencing Section 12's comparison table, explain in your own words why an XDP program dropping malicious traffic is cheaper for the system overall than the same traffic reaching an `iptables` DROP rule.

### Hard
7. Design, at a conceptual level (no code required), the map structure Cilium would need to implement Chapter 104's ClusterIP-to-backend-pod translation (Section 17) — what should the map's key be, what should its value be, and what has to update it when a pod is rescheduled?
8. Explain, referencing Section 7, why an eBPF program is not allowed to contain an unbounded `while(1)` loop even if the programmer is confident it will always terminate quickly in practice — what would the verifier need to prove, and why can it not simply run the program and see what happens?
9. A cluster running Cilium in kube-proxy-replacement mode (Section 17) reports much lower p99 latency for Service-to-Service calls than the same cluster previously did under `iptables`-mode kube-proxy, especially as the number of Services grew past several thousand. Referencing Sections 8, 11, and 12, explain precisely which mechanism accounts for this improvement, and why the improvement should be expected to grow more pronounced, not less, as Service count keeps increasing.

---

## 25. Summary and Bridge to Chapter 106

| Term | Meaning |
|---|---|
| eBPF | Sandboxed, verified programs loadable into the running Linux kernel at runtime |
| Verifier | The kernel component that statically proves an eBPF program's safety before allowing it to load |
| eBPF map | A kernel-managed key-value store shared between eBPF programs and userspace |
| XDP | The earliest network hook, at the NIC driver, ingress-only, for extreme-performance decisions |
| TC (Traffic Control) | A later network hook, ingress and egress, with fuller packet context |
| Cilium | A CNCF Kubernetes CNI whose entire data plane is built on eBPF instead of iptables |
| Security identity | Cilium's label-derived, IP-independent identifier used for stable policy enforcement |
| Hubble | Cilium's eBPF-sourced, per-flow network observability component |

This closes Part 16's arc: Chapter 102 built the raw kernel primitives (namespaces, veth, bridges) by hand; Chapter 103 showed Docker and CNI automating them into a pluggable standard; Chapter 104 showed Kubernetes building a stricter networking model and the Service/`kube-proxy` abstraction on top of it; and this chapter showed eBPF, XDP, and TC replacing the `iptables`-based mechanics underneath all of it with programs that run inside the kernel itself, verified safe before they ever execute a single packet. Every chapter in this arc has been about mechanisms working *underneath* an application that never has to know any of this exists — a service just opens a socket (Chapter 57) and sends bytes. Part 17 changes altitude one final time and stops describing those mechanisms from the outside: Chapter 106 opens with the socket API itself, and builds, in real Go code, the actual client and server programs that everything from Chapter 57 onward has been describing in prose.
