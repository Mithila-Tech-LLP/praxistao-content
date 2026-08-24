# Chapter 97: Cloud Networking — VPCs, Subnets, Route Tables, and Security Groups

> **"A cloud provider's data center might have your competitor's servers in the very next rack, on the very same physical switches described in Chapter 94. A VPC is the reason that fact is never supposed to matter to either of you."**

---

## Table of Contents

1. [Where This Chapter Picks Up](#1-where-this-chapter-picks-up)
2. [The Problem: One Physical Data Center, Many Untrusting Tenants](#2-the-problem-one-physical-data-center-many-untrusting-tenants)
3. [Naive Fix #1: Dedicated Physical Hardware Per Customer](#3-naive-fix-1-dedicated-physical-hardware-per-customer)
4. [Naive Fix #2: Shared Hardware, No Isolation](#4-naive-fix-2-shared-hardware-no-isolation)
5. [The Real Solution: The Virtual Private Cloud](#5-the-real-solution-the-virtual-private-cloud)
6. [What a VPC Actually Is](#6-what-a-vpc-actually-is)
7. [Subnets Within a VPC](#7-subnets-within-a-vpc)
8. [Public Subnets vs. Private Subnets](#8-public-subnets-vs-private-subnets)
9. [Route Tables: Controlling Traffic Flow](#9-route-tables-controlling-traffic-flow)
10. [Security Groups: Stateful, Instance-Level Firewalls](#10-security-groups-stateful-instance-level-firewalls)
11. [Network ACLs: Stateless, Subnet-Level Firewalls](#11-network-acls-stateless-subnet-level-firewalls)
12. [Security Groups vs. NACLs, Directly Compared](#12-security-groups-vs-nacls-directly-compared)
13. [Full Worked Example: Designing a Small VPC](#13-full-worked-example-designing-a-small-vpc)
14. [Real-World Cloud Networking](#14-real-world-cloud-networking)
15. [Hands-On Experiment](#15-hands-on-experiment)
16. [Code: Simulating Stateful vs. Stateless Filtering in Go](#16-code-simulating-stateful-vs-stateless-filtering-in-go)
17. [Common Misconceptions](#17-common-misconceptions)
18. [Production Notes](#18-production-notes)
19. [What's Simplified Here](#19-whats-simplified-here)
20. [Interview Questions & Model Answers](#20-interview-questions--model-answers)
21. [Exercises](#21-exercises)
22. [Summary and Bridge to Chapter 98](#22-summary-and-bridge-to-chapter-98)

---

## 1. Where This Chapter Picks Up

Chapters 94-96 described data centers, load balancers, and CDNs largely as if you owned and operated the hardware yourself — racks, ToR switches, physical load-balancer appliances. Almost nobody actually does this anymore for a new service. Instead, they rent slices of somebody else's data center — AWS, Google Cloud, Microsoft Azure, and others — and get, from the outside, something that behaves exactly like their own private network, with their own IP address ranges, their own routers, their own firewalls, running on physical infrastructure they will never see, touch, or share knowledge of with anyone else using it. This chapter explains how that illusion is built, and it's genuinely useful to understand it as an illusion deliberately and carefully constructed on real, shared, physical hardware — not as some entirely separate kind of network.

---

## 2. The Problem: One Physical Data Center, Many Untrusting Tenants

A single cloud provider region might run tens of thousands of physical servers across multiple data centers, all built on the leaf-spine fabric Chapter 94 described. On that shared physical infrastructure, the provider needs to simultaneously host thousands of unrelated, mutually-untrusting customers — a bank, a startup, a government agency, and a hobbyist's side project, potentially with virtual machines running on the very same physical host.

The requirement this creates is blunt: **Customer A's traffic must never be visible to, reachable by, or able to interfere with Customer B's traffic or resources — even though both are running on hardware physically owned, wired, and operated by the same third party, possibly sharing the exact same physical server, ToR switch, or spine switch.** Whatever solves this also has to let each customer configure *their own* addressing scheme, routing behavior, and security policy — independently, without needing the provider's staff to hand-configure anything for them, and without one customer's configuration being visible to, or capable of affecting, any other's.

---

## 3. Naive Fix #1: Dedicated Physical Hardware Per Customer

The most obviously *safe* approach: give every customer their own physical servers, switches, and cabling, fully separate from every other customer's. This is real isolation, with no shared-infrastructure risk at all — and it is, in fact, still offered by cloud providers today as a premium option (AWS calls this a "Dedicated Host," for instance).

It fails as a *default*, general-purpose solution for the same reason Chapter 94, Section 7 rejected building one enormous, dedicated core switch per use case: **it destroys the entire economic argument for cloud computing.** The value of a cloud provider is precisely that thousands of customers *share* the same underlying physical fleet, so no individual customer has to buy, rack, cool, or maintain hardware sized for their own peak demand — a customer needing three servers ninety percent of the year and thirty for a two-week seasonal spike would need to own thirty servers year-round under dedicated-hardware isolation. Multi-tenancy — many customers sharing one physical fleet — is the entire point; dedicated hardware for everyone is not a general-purpose answer, it's opting out of cloud computing's central economic idea.

---

## 4. Naive Fix #2: Shared Hardware, No Isolation

The opposite naive approach: everyone shares the same physical network, addressed and routed exactly like Chapter 94's data-center fabric, with no separation at all — one giant, flat address space and network for every customer, at once.

This is obviously unacceptable the moment you say it plainly: any customer could send traffic to, or receive traffic addressed as if it came from, any other customer's servers; a compromised virtual machine belonging to one customer could scan and attack every other customer's resources on the same physical network as though they were all one organization's LAN. There is no privacy, no security boundary, and no way for two customers to even use overlapping private IP ranges (which Chapter 40 established as extremely common — nearly everyone's home router uses `192.168.1.0/24`) without immediate, catastrophic address collisions.

---

## 5. The Real Solution: The Virtual Private Cloud

The actual answer splits the difference cleanly: **share the physical hardware, but make the network itself virtual** — build a genuinely isolated *logical* network for each customer, entirely in software, running on top of the same shared physical fabric, such that from inside any one customer's logical network, it looks and behaves exactly like a private network they own outright, indistinguishable in practice from Chapter 41's private-addressing model, with no visibility into (or from) any other customer's equally-isolated logical network occupying the very same physical switches underneath.

This is a **Virtual Private Cloud (VPC)**: a customer's own isolated, logically-defined network, carved out of a cloud provider's shared physical infrastructure, with the customer given full control over its addressing, subnetting, routing, and security policy — exactly the set of concerns Chapters 36-52 spent entire volumes teaching you to reason about for a *physical* network, now applied to a network that has no dedicated physical wires of its own at all.

**Intuitive level:** think of a large shared office building, subdivided into many separate, soundproofed, individually keyed suites, each tenant free to arrange furniture and set their own house rules inside their own suite, with no tenant able to see, hear, or enter another tenant's space — even though the building's plumbing, electricity, and structural walls are all shared infrastructure underneath. The analogy breaks slightly in one important way: a VPC's isolation is enforced by network virtualization technology (a preview of Chapter 99's overlay networks and VXLAN, which explain the actual encapsulation mechanism making this practically real), not by physical walls at all — the "walls" here are entirely a property of how packets are tagged, encapsulated, and filtered by the underlying physical fabric's software, invisible to anyone actually looking at the wiring.

**Deep technical level:** each VPC's traffic is, in essentially every major cloud provider's implementation, carried over the shared physical fabric encapsulated inside an overlay network — the customer's own packets, complete with whatever private IP addressing they've chosen (even overlapping with another customer's, since the two never actually meet), are wrapped inside an outer packet the physical fabric actually forwards, and unwrapped again at the destination. Chapter 99 explains this encapsulation mechanism (VXLAN and its relatives) in full technical detail; for this chapter, the important fact is only that this is *how* the illusion of a dedicated private network is made physically real without dedicating any physical hardware at all.

---

## 6. What a VPC Actually Is

Concretely, when you create a VPC on a major cloud provider, you're defining:

- **An IP address range (CIDR block)** for the whole VPC — commonly a private range from Chapter 40's RFC 1918 space, like `10.0.0.0/16`, giving you roughly 65,000 usable addresses (Chapter 39's CIDR math) to allocate as you see fit within your own network.
- **A region** — the VPC exists within one cloud provider region (a specific geographic area, itself built from multiple physical data centers, tying back to Chapters 94-96's physical infrastructure), and typically cannot span regions directly (cross-region connectivity is its own separate feature, layered on top).
- **Complete isolation by default** — a brand-new VPC, by default, has no connectivity to the public Internet, to any other VPC, or to any other customer's resources at all; every path in or out has to be explicitly created (Sections 8-11, and Chapter 98's gateways).

This is the crucial design decision worth internalizing: a VPC starts **maximally closed**, and the customer opens exactly the specific paths they need — the same "default deny, explicitly allow" philosophy Chapter 84 introduced for firewalls generally, now applied as the VPC's starting posture rather than an optional configuration choice.

---

## 7. Subnets Within a VPC

A VPC's overall CIDR block is rarely used as one flat address space — exactly the subnetting motivation Chapter 38 introduced for physical networks applies here too, just inside a virtual one. A **subnet** carves out a smaller CIDR range from the VPC's block (e.g., `10.0.1.0/24` and `10.0.2.0/24` out of the VPC's `10.0.0.0/16`), and — critically, a cloud-specific detail with no direct physical-network equivalent — **each subnet exists in exactly one Availability Zone (AZ)**, a cloud provider's term for one physically distinct data center (or an isolated portion of one) within the region, engineered to fail independently of other AZs in the same region (separate power, separate cooling, separate physical building in most designs).

This means subnetting inside a VPC is doing double duty compared to Chapter 38's original motivation: it still organizes address space logically (by function — Section 8), but it *also* controls physical fault-tolerance, since placing resources in subnets across multiple AZs is how a cloud architecture survives one entire physical data center failing, directly extending the redundancy reasoning Chapter 95's load balancing already introduced at the server level, now applied at the level of whole physical facilities.

---

## 8. Public Subnets vs. Private Subnets

Within a VPC, subnets are conventionally split by their intended exposure to the outside world — a distinction created entirely by *routing configuration* (Section 9), not by any inherent property of the subnet's address range itself:

- **A public subnet** has a route (Section 9) to an **Internet Gateway** (Chapter 98's subject), meaning resources placed in it can be reached from, and can reach, the public Internet directly — typically where a load balancer (Chapter 95) or a bastion/jump host lives.
- **A private subnet** has no such route; resources here (application servers, databases) are unreachable directly from the Internet no matter what security policy is applied to them, because there is no *path* at all — a structural guarantee, not just a policy one, and one of the most important defense-in-depth patterns in real cloud architecture: sensitive resources are placed where a misconfigured security rule literally cannot expose them, because the routing layer beneath the security layer never offers a way out to begin with.

---

## 9. Route Tables: Controlling Traffic Flow

A **route table** is precisely what Chapter 44 already taught you a router's routing table is — a set of rules mapping destination prefixes to next hops — attached to each subnet, and it is the actual mechanism that makes Section 8's public/private distinction real rather than a mere label.

A typical private subnet's route table might look like this:

| Destination | Target |
|---|---|
| `10.0.0.0/16` (the VPC's own range) | local |
| `0.0.0.0/0` (everything else — the default route, Chapter 45) | NAT Gateway (Chapter 98) |

And a public subnet's route table, by contrast:

| Destination | Target |
|---|---|
| `10.0.0.0/16` (the VPC's own range) | local |
| `0.0.0.0/0` (default route) | Internet Gateway (Chapter 98) |

The only difference between "public" and "private" here is literally one line in a table: where the default route (Chapter 45's `0.0.0.0/0` catch-all) points. Point it at an Internet Gateway, and the subnet is directly Internet-reachable; point it at a NAT Gateway instead (Chapter 98 explains the crucial difference between the two), and outbound connections work but nothing on the Internet can *initiate* a connection inward — exactly the asymmetry Chapter 98's full worked example depends on.

---

## 10. Security Groups: Stateful, Instance-Level Firewalls

A **security group** is a firewall attached directly to an individual resource — most commonly a virtual machine instance, but also databases, load balancers, and other resources — controlling exactly what traffic is allowed in and out of *that specific resource*, regardless of what subnet it happens to sit in.

The single most important property of a security group, worth internalizing precisely because it's the source of the sharpest contrast with Section 11: **security groups are stateful.** If you allow inbound traffic on port 443 (Chapter 57's HTTPS port) from anywhere, you do *not* need a separate rule allowing the corresponding outbound *reply* traffic — the security group automatically recognizes that outbound packet as part of an already-permitted connection (the same connection-tracking concept Chapter 84 introduced for stateful firewalls generally) and lets it through without a matching rule of its own. This mirrors exactly the stateful-firewall behavior Chapter 84 described in general terms; a security group is that concept, productized as a specific cloud networking primitive.

A typical security group for a web server:

| Direction | Protocol | Port | Source/Destination | Purpose |
|---|---|---|---|---|
| Inbound | TCP | 443 | `0.0.0.0/0` | Accept HTTPS from anywhere |
| Inbound | TCP | 22 | `10.0.5.0/24` (bastion subnet only) | SSH, restricted to admin subnet |
| Outbound | (all) | (all) | `0.0.0.0/0` | Allow all outbound (common default) |

Notice there's no explicit rule needed for "allow the HTTPS *response* out" — statefulness handles it. Security groups, in every major cloud provider's model, also support **default deny**: with no matching allow rule, traffic is dropped, and (in most implementations) security groups only support allow rules at all — there's no explicit "deny" rule to write, because anything not explicitly allowed is denied by default.

---

## 11. Network ACLs: Stateless, Subnet-Level Firewalls

A **Network Access Control List (NACL)** is a firewall attached to an entire **subnet** (Section 7), not to an individual instance — every resource in that subnet is subject to the same NACL rules, regardless of its own security group configuration; both layers apply simultaneously, and traffic must pass both to succeed.

The defining, sharply contrasting property here: **NACLs are stateless.** Allowing inbound traffic on port 443 does **not** automatically allow the corresponding outbound reply — you must write a separate, explicit outbound rule permitting return traffic (conventionally, allowing outbound traffic from the high-numbered ephemeral port range, Chapter 57's ephemeral ports, back to the original client), or legitimate reply traffic will be silently dropped even though the original inbound request was explicitly allowed.

NACLs also differ from security groups in supporting **explicit deny rules**, evaluated in numbered order — a genuine, ordered rule list rather than only-allow, which enables patterns security groups structurally cannot express, like "allow everything from this whole address range, except this one specific smaller range within it," using an explicit deny rule for the exception evaluated before a broader allow.

A NACL example, showing the ordered-evaluation and explicit-outbound-return behavior a security group never needs:

| Rule # | Direction | Protocol | Port | Source/Destination | Allow/Deny |
|---|---|---|---|---|---|
| 100 | Inbound | TCP | 443 | `0.0.0.0/0` | ALLOW |
| 200 | Inbound | TCP | all | `203.0.113.0/24` (a known-bad range) | DENY |
| \* | Inbound | all | all | `0.0.0.0/0` | DENY (implicit default) |
| 100 | Outbound | TCP | 1024-65535 (ephemeral) | `0.0.0.0/0` | ALLOW |
| \* | Outbound | all | all | `0.0.0.0/0` | DENY (implicit default) |

Without that explicit outbound rule for the ephemeral port range, every inbound HTTPS connection this NACL otherwise allows would still fail — the server could receive the request, but its reply would be silently dropped on the way out, a genuinely common, genuinely confusing real-world misconfiguration that Section 12 and the exercises return to directly.

---

## 12. Security Groups vs. NACLs, Directly Compared

| Property | Security Group | NACL |
|---|---|---|
| Attached to | Individual instance/resource | Entire subnet |
| State tracking | Stateful — return traffic auto-allowed | Stateless — return traffic needs its own explicit rule |
| Rule types | Allow only (implicit deny for the rest) | Explicit allow and deny rules, evaluated in order |
| Scope of effect | Only the specific resource(s) it's attached to | Every resource in the subnet, regardless of their own security groups |
| Typical use | Fine-grained, per-service policy ("only my app server accepts 8080") | Coarse-grained, subnet-wide policy and explicit blocklisting |
| Evaluated relative to security groups | N/A (this is the SG layer) | Applied in addition to, not instead of, any security groups on resources inside the subnet |

In a real deployment, both apply to every packet at once: a NACL is the outer, coarser gate for the whole subnet, and a security group is the inner, finer gate for the specific instance — a packet has to pass both to reach its destination, directly mirroring the general "defense in depth" principle Chapter 84 introduced (a WAF plus a network firewall plus host-level filtering, layered rather than substituted).

---

## 13. Full Worked Example: Designing a Small VPC

```
VPC: 10.0.0.0/16   (region: us-east-1)

  ┌─────────────────────────── Availability Zone A ───────────────────────────┐
  │  Public Subnet: 10.0.1.0/24        Private Subnet: 10.0.11.0/24            │
  │  ┌──────────────────┐              ┌───────────────────────┐              │
  │  │ Load Balancer     │──(SG:443)──►│ App Server              │              │
  │  │ route: 0.0.0.0/0  │              │ route: 0.0.0.0/0        │              │
  │  │  -> Internet GW   │              │  -> NAT Gateway          │              │
  │  └──────────────────┘              └───────────────────────┘              │
  └──────────────────────────────────────────────────────────────────────────┘

  ┌─────────────────────────── Availability Zone B ───────────────────────────┐
  │  Public Subnet: 10.0.2.0/24        Private Subnet: 10.0.12.0/24            │
  │  ┌──────────────────┐              ┌───────────────────────┐              │
  │  │ Load Balancer     │──(SG:443)──►│ App Server              │              │
  │  │ (same LB, 2nd AZ) │              │ route: 0.0.0.0/0        │              │
  │  └──────────────────┘              │  -> NAT Gateway          │              │
  │                                     └───────────────────────┘              │
  └──────────────────────────────────────────────────────────────────────────┘

  NACL on both private subnets: deny inbound from 0.0.0.0/0 except from
  the public subnets' CIDR ranges, on the app port only.
```

This design combines nearly every idea in this chapter: one VPC, subnetted across two Availability Zones for fault tolerance (Section 7), a public/private split enforced entirely by route tables (Sections 8-9), security groups controlling exactly which instance-level traffic is allowed (Section 10), and a NACL adding a coarser, subnet-wide backstop (Section 11) — with Chapter 98 supplying the Internet Gateway and NAT Gateway that make the public subnet's inbound path and the private subnet's outbound path actually work.

---

## 14. Real-World Cloud Networking

- **AWS** calls this primitive a **VPC** directly, and its security-group/NACL split is exactly the terminology and stateful/stateless distinction this chapter has used throughout.
- **Google Cloud** calls the equivalent a **VPC network** too, though with some structural differences — a GCP VPC network can span multiple regions globally by default, unlike AWS's region-scoped VPC, and its firewall rules are closer to security groups' stateful model, without a separately-named NACL-equivalent as a default primitive (Hierarchical Firewall Policies serve a related, but not identical, coarser-grained role).
- **Microsoft Azure** calls the equivalent a **Virtual Network (VNet)**, with **Network Security Groups (NSGs)** playing a hybrid role that can attach to either a subnet or an individual network interface, and Azure's NSGs are themselves stateful — Azure's closest stateless, ordered-rule primitive is a separate feature (Azure Firewall) rather than a direct NACL equivalent.
- All three, underneath their differing terminology, solve Section 5's identical core problem the identical general way: software-defined, per-customer logical network isolation over shared physical infrastructure, secured by layered stateful and/or ordered rule-based filtering.

---

## 15. Hands-On Experiment

Even without a paid cloud account, most providers offer a generous enough free tier to try this directly:

1. Create a new VPC with CIDR `10.0.0.0/16`. Create two subnets, `10.0.1.0/24` and `10.0.2.0/24`, in two different Availability Zones.
2. Launch a small instance in each subnet. Attach a security group allowing inbound SSH (port 22) only from your own current public IP address (most cloud consoles offer a "my IP" auto-fill for exactly this).
3. Create an Internet Gateway, attach it to the VPC, and add a route in one subnet's route table sending `0.0.0.0/0` to it — this is now your "public" subnet.
4. Try to SSH into the instance in the subnet *without* that route. It should fail to connect at all, regardless of security group configuration — direct, hands-on proof of Section 8's claim that routing, not just security policy, is what makes a subnet reachable.
5. Add a NACL to the public subnet that explicitly denies inbound traffic from your own IP, while your security group still allows it. Confirm the connection now fails — direct proof that a NACL's subnet-wide deny overrides a permissive security group, since a packet must pass both layers.

---

## 16. Code: Simulating Stateful vs. Stateless Filtering in Go

A minimal simulation making Section 10 and Section 11's core distinction — statefulness — concrete and testable:

```go
package main

import "fmt"

type Rule struct {
	Direction string // "in" or "out"
	Port      int
	Allow     bool
}

// StatefulFirewall models a security group: only "in" rules are configured;
// matching "out" traffic for an already-permitted connection is automatic.
type StatefulFirewall struct {
	rules            []Rule
	establishedConns map[int]bool // tracks ports with an allowed inbound connection
}

func NewStatefulFirewall(rules []Rule) *StatefulFirewall {
	return &StatefulFirewall{rules: rules, establishedConns: make(map[int]bool)}
}

func (f *StatefulFirewall) AllowInbound(port int) bool {
	for _, r := range f.rules {
		if r.Direction == "in" && r.Port == port && r.Allow {
			f.establishedConns[port] = true // remember this connection's state
			return true
		}
	}
	return false
}

func (f *StatefulFirewall) AllowOutboundReply(port int) bool {
	// No matching "out" rule needed at all -- statefulness is the whole point.
	return f.establishedConns[port]
}

// StatelessFirewall models a NACL: "in" and "out" are evaluated completely
// independently, with no memory of prior connections.
type StatelessFirewall struct {
	rules []Rule
}

func (f *StatelessFirewall) Allow(direction string, port int) bool {
	for _, r := range f.rules {
		if r.Direction == direction && r.Port == port {
			return r.Allow
		}
	}
	return false // implicit deny, same as a real NACL's default
}

func main() {
	sg := NewStatefulFirewall([]Rule{{Direction: "in", Port: 443, Allow: true}})
	fmt.Println("Security group inbound 443:", sg.AllowInbound(443))       // true
	fmt.Println("Security group outbound reply on 443:", sg.AllowOutboundReply(443)) // true, no rule needed

	// A NACL configured with only the inbound rule, and no ephemeral outbound rule --
	// this is the classic real-world misconfiguration Section 11 warned about.
	nacl := &StatelessFirewall{rules: []Rule{{Direction: "in", Port: 443, Allow: true}}}
	fmt.Println("NACL inbound 443:", nacl.Allow("in", 443))   // true
	fmt.Println("NACL outbound 443:", nacl.Allow("out", 443)) // false -- the reply is silently dropped!
}
```

```
Security group inbound 443: true
Security group outbound reply on 443: true
NACL inbound 443: true
NACL outbound 443: false
```

That last `false` is exactly the real, common production bug Section 11 described: an inbound request an operator clearly intended to allow, whose reply silently never leaves, because a stateless NACL was configured with the mental model of a stateful security group.

---

## 17. Common Misconceptions

- **"A VPC is a physically separate network."** It's a logically isolated network running on shared physical infrastructure, made real through network virtualization/overlay encapsulation (Chapter 99), not through dedicated wiring.
- **"Security groups and NACLs do the same job, just at different scopes."** The scope difference (instance vs. subnet) is real, but the state-tracking difference is the more consequential one in practice — it's the reason NACL misconfigurations (forgetting the return-traffic rule) are a distinctly different, and distinctly common, class of bug from security group misconfigurations.
- **"A private subnet is private because of its security group rules."** A private subnet's isolation from inbound Internet traffic is a routing fact (no route to an Internet Gateway, Section 8) — it would remain unreachable from the Internet even with a maximally permissive security group, because there's no path there to begin with.
- **"Overlapping private IP ranges between two VPCs is automatically a conflict."** It isn't, precisely because VPCs are isolated logical networks — two completely separate customers (or even two VPCs owned by the same customer) can both use `10.0.0.0/16` with no conflict at all, right up until someone tries to *connect* the two (via VPC peering or a similar mechanism), at which point overlapping ranges genuinely do become a real, practical problem requiring redesign.
- **"More security layers (SG + NACL + host firewall) is always strictly better with no cost."** Layering is good defense-in-depth practice, but every additional stateless layer (NACLs especially) is another place a forgotten return-traffic rule can silently break otherwise-correct traffic — complexity itself is a real, recurring source of production incidents here.

---

## 18. Production Notes

- Most real production VPC designs use *more* than two subnets per Availability Zone — commonly a public subnet, a private "application" subnet, and an even more restricted private "data" subnet (for databases), each with progressively tighter NACLs and security groups, layering Section 8's public/private idea into finer gradations.
- Security group rules are commonly written to reference *other security groups* as their source/destination, not raw CIDR ranges — "allow inbound from any instance with the `app-server` security group attached" — so the rule automatically stays correct as instances are added or removed, without needing IP-address bookkeeping.
- NACLs are used sparingly in many real-world architectures specifically *because* of their statelessness pitfall (Section 16's bug) — many teams rely primarily on security groups for day-to-day policy and use NACLs mainly for coarse, rarely-changed guardrails (like an explicit deny against a known-malicious IP range) rather than fine-grained, frequently-updated rules.
- Terraform, CloudFormation, and similar infrastructure-as-code tools are the near-universal way real VPCs, subnets, route tables, and security groups are actually defined and changed in production — manually clicking through a cloud console to configure networking is common for learning, rare for production systems that need auditable, repeatable configuration.

---

## 19. What's Simplified Here

Real cloud VPC implementations include substantially more than this chapter covers: VPC peering and Transit Gateways for connecting multiple VPCs together; VPC endpoints for reaching a provider's other services without traversing the public Internet at all; flow logs for auditing exactly what traffic a security group or NACL actually allowed or blocked; and provider-specific quirks in exactly how many rules a security group or NACL can hold, and exactly how ordering and rule evaluation work (AWS NACL rules are evaluated in strict numeric order and stop at the first match; other providers' equivalents vary in these details). The core concepts — logical isolation over shared infrastructure, subnets tied to Availability Zones, routing-determined public/private status, and the stateful/stateless split between instance-level and subnet-level firewalls — are accurate and consistent across all major providers.

---

## 20. Interview Questions & Model Answers

**Beginner: What is a VPC, in one sentence?**
A Virtual Private Cloud is a customer's own logically isolated network, with its own address space, subnets, routing, and security policy, running on a cloud provider's shared physical infrastructure without needing dedicated physical hardware.

**Beginner: What determines whether a subnet is "public" or "private" in a cloud VPC?**
Its route table — specifically, whether its default route (`0.0.0.0/0`) points to an Internet Gateway (making it public) or to something else entirely, like a NAT Gateway or nothing at all (making it private). It is a routing fact, not a security-group setting.

**Intermediate: What is the core difference between a security group and a Network ACL?**
A security group is stateful and attaches to individual instances — allowing inbound traffic on a port automatically allows the matching outbound reply with no extra rule. A NACL is stateless and attaches to an entire subnet — inbound and outbound traffic are evaluated as completely independent rule sets, so allowing inbound traffic on a port requires a *separate*, explicit rule to allow the corresponding reply traffic out, or it will be silently dropped.

**Intermediate: Why can two different customers' VPCs both use the CIDR range 10.0.0.0/16 without any conflict?**
Because each VPC is a logically isolated network — its addressing exists only within that isolated context, carried over the shared physical fabric via network virtualization/overlay encapsulation. The two networks never actually interact unless something explicitly connects them (like VPC peering), so identical address ranges in two separate VPCs simply never meet.

**Advanced: Why does a VPC default to "closed" (no Internet access, no connectivity to anything else) rather than defaulting to open with explicit deny rules layered on?**
This mirrors the "default deny" security posture Chapter 84 introduced for firewalls generally — starting maximally closed means a misconfiguration or omission fails safe (nothing is reachable that wasn't explicitly enabled) rather than failing open (something is reachable that nobody intended to expose). Given that VPCs run on shared physical infrastructure specifically designed to isolate mutually-untrusting tenants, an accidentally-open default would be a far more severe failure mode than an accidentally-closed one, which just requires a follow-up fix rather than an active security incident.

**Advanced: A NACL on a private subnet allows inbound traffic on port 443 but has no outbound rule for the ephemeral port range. What happens, mechanically, to an inbound HTTPS request, and why?**
The initial inbound SYN packet on port 443 is allowed in and reaches the destination instance, which processes the request and attempts to send its reply back — sourced from port 443, destined to the client's ephemeral source port. Because the NACL is stateless, that outbound packet is evaluated against the NACL's outbound rules independently of the inbound rule that let the request in; with no outbound rule permitting traffic to the ephemeral port range, the reply is silently dropped by the NACL, and from the client's perspective, the connection times out despite the server having successfully received and processed the request.

---

## 21. Exercises

### Easy
1. In one sentence each, define VPC, subnet, and Availability Zone, and explain how the three relate to each other.
2. Explain why a private subnet is unreachable from the Internet even if its security group allows all inbound traffic from `0.0.0.0/0`.
3. List two differences between a security group and a NACL.

### Medium
4. Design a route table for a public subnet and a route table for a private subnet in the same VPC, using the format from Section 9, and explain the one-line difference between them that determines their public/private status.
5. Using the Go code in Section 16, add a second `StatelessFirewall` example that correctly includes the ephemeral-port outbound rule, and show that the reply is now allowed.
6. Explain why a security group is normally described as supporting only "allow" rules, while a NACL supports both "allow" and explicit "deny" rules, and give one real scenario where an explicit deny is genuinely useful.

### Hard
7. Two teams within the same company each build a VPC using `10.0.0.0/16`, unaware of each other. Six months later, the company wants to connect the two VPCs via peering so services in each can talk to services in the other. Explain exactly why this is now a problem, and describe two different ways to resolve it.
8. Design the full security posture (route tables, security groups, and NACLs) for a three-tier VPC architecture: a public load-balancer subnet, a private application-server subnet, and an even more restricted private database subnet, such that the database subnet can only ever receive traffic from the application-server subnet, on the database's specific port, from nowhere else, at any layer.
9. A production incident report says: "Users could successfully connect to our API on port 443, but every response timed out, even though CloudWatch/monitoring shows the backend successfully processed every request." Using this chapter's material, propose the single most likely root cause, and explain exactly how you would confirm it using the concepts from Sections 10-11.

---

## 22. Summary and Bridge to Chapter 98

| Term | Meaning |
|---|---|
| VPC (Virtual Private Cloud) | A customer's isolated logical network on a cloud provider's shared physical infrastructure |
| Region / Availability Zone | A geographic area (region) made of multiple independently-failing physical facilities (AZs) |
| Subnet | A smaller CIDR range within a VPC, tied to exactly one Availability Zone |
| Public / private subnet | Determined entirely by the subnet's route table — whether its default route reaches an Internet Gateway |
| Route table | Per-subnet rules mapping destination prefixes to next hops, controlling reachability |
| Security group | Stateful firewall attached to an individual instance; return traffic auto-allowed |
| NACL (Network ACL) | Stateless firewall attached to a whole subnet; inbound and outbound evaluated independently |

This chapter deliberately deferred two names that kept appearing in every route table example — Internet Gateway and NAT Gateway — because they deserve full treatment on their own. A VPC's public subnet can only actually reach, and be reached from, the Internet because of one of these; its private subnet can only reach *out* to the Internet (for software updates, external API calls) because of the other, and never be reached *from* it. Chapter 98 traces exactly how both work, and finishes with a complete, hop-by-hop request path through a realistic VPC — bringing this volume's cloud networking material to its full, concrete conclusion.
