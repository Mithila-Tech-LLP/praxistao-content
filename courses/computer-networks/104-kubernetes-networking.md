# Chapter 104: Kubernetes Networking

> **"Chapter 103 gave every container a plug-in standard for getting a network stack. Kubernetes takes that standard and builds a hard requirement on top of it: every one of potentially thousands of constantly-appearing-and-disappearing pods must be individually addressable, from anywhere in the cluster, with no NAT in the way — and then has to solve the much harder problem of letting anything else reliably find a pod that might not exist by the time you finish reading its IP address."**

---

## Table of Contents

1. [Where This Chapter Picks Up](#1-where-this-chapter-picks-up)
2. [The Problem: Orchestrating Chapter 103's Containers at Scale](#2-the-problem-orchestrating-chapter-103s-containers-at-scale)
3. [Kubernetes's Networking Model: Three Non-Negotiable Rules](#3-kubernetess-networking-model-three-non-negotiable-rules)
4. [Naive Fix: Just Use Docker's Default Bridge Network](#4-naive-fix-just-use-dockers-default-bridge-network)
5. [Why the Naive Fix Fails at Cluster Scale](#5-why-the-naive-fix-fails-at-cluster-scale)
6. [The Real Solution, Part 1: A Pod-Per-IP CNI Setup](#6-the-real-solution-part-1-a-pod-per-ip-cni-setup)
7. [Hands-On: Inspecting Pod IPs and the Pod Network](#7-hands-on-inspecting-pod-ips-and-the-pod-network)
8. [The Next Problem: Pods Are Ephemeral, Everything Else Isn't](#8-the-next-problem-pods-are-ephemeral-everything-else-isnt)
9. [The Real Solution, Part 2: The Service Abstraction](#9-the-real-solution-part-2-the-service-abstraction)
10. [ClusterIP: The Default Service Type](#10-clusterip-the-default-service-type)
11. [How kube-proxy Actually Implements a ClusterIP](#11-how-kube-proxy-actually-implements-a-clusterip)
12. [Hands-On: Watching kube-proxy's iptables Rules](#12-hands-on-watching-kube-proxys-iptables-rules)
13. [NodePort: Reaching a Service From Outside the Cluster](#13-nodeport-reaching-a-service-from-outside-the-cluster)
14. [LoadBalancer: Handing Off to Chapter 95's Load Balancers](#14-loadbalancer-handing-off-to-chapter-95s-load-balancers)
15. [Service Discovery: DNS Names for Services](#15-service-discovery-dns-names-for-services)
16. [Headless Services: When You Want the Pod IPs Directly](#16-headless-services-when-you-want-the-pod-ips-directly)
17. [Ingress: HTTP-Aware Routing on Top of Services](#17-ingress-http-aware-routing-on-top-of-services)
18. [Network Policies: Firewalling Pod-to-Pod Traffic](#18-network-policies-firewalling-pod-to-pod-traffic)
19. [Code: A Minimal ClusterIP-to-Endpoint Resolver in Go](#19-code-a-minimal-clusterip-to-endpoint-resolver-in-go)
20. [Common Misconceptions](#20-common-misconceptions)
21. [Production Notes](#21-production-notes)
22. [What's Simplified Here](#22-whats-simplified-here)
23. [Interview Questions & Model Answers](#23-interview-questions--model-answers)
24. [Exercises](#24-exercises)
25. [Summary and Bridge to Chapter 105](#25-summary-and-bridge-to-chapter-105)

---

## 1. Where This Chapter Picks Up

Chapter 103 closed by naming exactly what Kubernetes needs from a CNI plugin: pods, not just containers. This chapter answers what a "pod" networking-wise even is, why Kubernetes imposes networking rules stricter than Docker's own default (Chapter 103, Section 4), and how the Service abstraction — combined with `kube-proxy` and, at the edge, the load balancers Chapter 95 already built up in detail — turns a constantly-churning set of pod IPs into something stable enough to actually depend on.

---

## 2. The Problem: Orchestrating Chapter 103's Containers at Scale

A single host running Docker (Chapter 103) is one failure domain: if that machine dies, every container on it dies. Real production workloads need containers spread across many machines, automatically rescheduled onto a healthy machine when one fails, scaled up and down based on load, and — critically for this chapter — able to find and talk to each other correctly regardless of which specific machine they end up running on at any given moment.

Kubernetes is the orchestrator that manages this scheduling and lifecycle. Its fundamental unit isn't a single container but a **pod** — one or more containers that are always scheduled together, onto the same machine, sharing the same network namespace (Chapter 102, Section 7). Two containers in the same pod share one IP address and can reach each other over `localhost`, exactly like two processes in the same Chapter 102 network namespace, because that is literally what they are.

The problem this chapter solves: given a cluster of many machines (called **nodes**), each running many pods that are created, destroyed, and rescheduled constantly, how does any pod reliably reach any other pod, and how does anything reach a *group* of interchangeable pods (like "the current set of web server replicas") without needing to track which specific ones currently exist?

---

## 3. Kubernetes's Networking Model: Three Non-Negotiable Rules

Before describing any implementation, Kubernetes's own specification states three requirements every networking implementation (every CNI plugin, Chapter 103) must satisfy, no matter how it's implemented underneath:

1. **Every pod gets its own real, cluster-wide IP address.** Not a per-host private address that needs translation to be seen elsewhere (as Chapter 103, Section 4's `172.17.0.0/16` addresses were, scoped to one host) — an IP that means the same pod no matter which node in the cluster is asking.
2. **Pods can reach each other's IPs directly, with no NAT.** A pod on Node A must be able to send a packet straight to a pod's real IP on Node B and have it arrive unmodified — the packet's source and destination addresses must be exactly the pod IPs on both ends, matching what each pod itself believes its own address to be.
3. **A node can reach any pod on any node**, and (with a Kubernetes-specific carve-out) **a pod can reach itself via the address other pods would use to reach it.**

Rule 2 is the sharpest departure from Chapter 103's Docker default: Section 4 of that chapter had every container behind host-level NAT, needing an explicit port-publish to be reached from anywhere else. Kubernetes explicitly forbids that model for pod-to-pod traffic — it requires what Chapter 99 called a genuinely flat address space across the whole cluster, whatever the underlying physical topology actually looks like.

---

## 4. Naive Fix: Just Use Docker's Default Bridge Network

The obvious first attempt, given that Kubernetes pods are ultimately containers: reuse exactly Chapter 103's Section 4 default bridge setup on every node — each node gets its own `docker0`-equivalent bridge, its own private subnet, pods get IPs from that subnet.

---

## 5. Why the Naive Fix Fails at Cluster Scale

This fails Rule 1 and Rule 2 immediately. If every node independently runs its own `172.17.0.0/16`, two different nodes will hand out the *exact same* pod IP to two different pods — there is no cluster-wide uniqueness. And even if the subnets were made unique per node, Chapter 103, Section 4's model puts every pod behind host-level NAT: a pod on Node A trying to reach a pod on Node B by its real IP has no route there at all, because that IP only exists inside Node B's private, host-scoped bridge network, invisible from outside that one host.

This is structurally the exact problem Chapter 103, Section 18 introduced for multi-host container traffic in general, now stated as a hard requirement instead of an optional nice-to-have. Kubernetes's answer is the same class of fix — non-overlapping address allocation plus cluster-wide routability — just enforced as policy rather than left optional.

---

## 6. The Real Solution, Part 1: A Pod-Per-IP CNI Setup

A conformant Kubernetes CNI setup (any of Calico, Cilium — Chapter 105 — Flannel, or several others) satisfies Section 3's rules with a pattern that generalizes Chapter 103's IPAM discussion (Section 14 of that chapter):

- The cluster is assigned one large address block up front — a **cluster CIDR**, for example `10.244.0.0/16`.
- Each node is assigned a distinct, non-overlapping slice of it — a **pod CIDR**, for example `10.244.1.0/24` for Node A and `10.244.2.0/24` for Node B — guaranteeing Rule 1's uniqueness with zero cross-node coordination needed at pod-creation time, because each node already owns its own exclusive range.
- Each node runs a CNI plugin (Chapter 103, Sections 10–14) that, on every pod creation, does exactly Chapter 102's namespace/veth/bridge dance, assigning the new pod an IP from that node's own pod CIDR.
- Cluster-wide routability (Rule 2) is provided one of two ways: either an **overlay network** (Chapter 103, Section 18; Chapter 99's VXLAN) encapsulating pod traffic between nodes so it rides over the real underlay network transparently, or, in "native routing" mode (Calico's default, and Cilium's when the underlying network supports it), each node's real router/switch fabric is configured with routes for every other node's pod CIDR, letting pod traffic be forwarded as ordinary, unencapsulated IP packets.

```
   Node A (pod CIDR 10.244.1.0/24)      Node B (pod CIDR 10.244.2.0/24)

   +----------------------------+      +----------------------------+
   | Pod 1: 10.244.1.5           |      | Pod 3: 10.244.2.7           |
   | Pod 2: 10.244.1.6           |      | Pod 4: 10.244.2.8           |
   |         |                  |      |         |                  |
   |   CNI-managed bridge/veths  |      |   CNI-managed bridge/veths  |
   +------------|---------------+      +------------|---------------+
                |                                    |
                +---- cluster network (overlay ------+
                        or native routing) ----
```

This is Chapter 103's Section 4 diagram, generalized from one host's bridge to a whole cluster's worth of nodes, with the previously-hardcoded per-host subnet now sliced out of one coordinated cluster CIDR instead of every node reusing the same default.

---

## 7. Hands-On: Inspecting Pod IPs and the Pod Network

Assuming access to a running cluster (a local one via `kind` or `minikube` works identically for this purpose):

```bash
# List every pod's IP, and which node it landed on -- Rule 1 in action.
kubectl get pods -o wide

# Example output:
# NAME              READY   STATUS    IP            NODE
# web-7d4b9-abc12   1/1     Running   10.244.1.5    node-a
# web-7d4b9-def34   1/1     Running   10.244.2.7    node-b

# Confirm Rule 2 directly: exec into one pod, curl another pod's IP
# straight across nodes, with no port-forwarding or NAT involved.
kubectl exec -it web-7d4b9-abc12 -- curl -s 10.244.2.7:80

# Inspect a node's own pod CIDR allocation (Section 6):
kubectl get node node-a -o jsonpath='{.spec.podCIDR}'
# -> 10.244.1.0/24
```

The `curl` in the second command reaching a pod on a different node, by its raw pod IP, with a real response coming back, is Rule 2 made concrete — the exact guarantee Chapter 103's Docker default (Section 4 of that chapter) explicitly did not provide across hosts.

---

## 8. The Next Problem: Pods Are Ephemeral, Everything Else Isn't

Section 7's `curl 10.244.2.7:80` works, but it depends on a fact that Kubernetes deliberately does not guarantee: that pod's IP staying the same. Kubernetes reschedules pods constantly — a pod crashes and gets recreated, a deployment is scaled from three replicas to five, a node is drained for maintenance — and every single time a pod is recreated, it typically gets a **brand-new IP address**.

Hardcoding `10.244.2.7` anywhere in application configuration is therefore a latent bug: it will work until the next reschedule, then silently point at nothing (or worse, at whatever different pod later happens to receive that now-recycled address). This is the same class of problem Chapter 101 described service meshes solving with service discovery, and Chapter 66 described DNS solving for human-readable names over changing IPs — except here the thing that "changes" is not a server being migrated occasionally, but a routine, constant, and expected event.

---

## 9. The Real Solution, Part 2: The Service Abstraction

A Kubernetes **Service** is a stable, virtual IP address and DNS name that sits in front of a *set* of pods selected by a label (for instance, "every pod labeled `app=web`"), and stays constant even as the specific pods behind it come and go.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  selector:
    app: web
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
```

This declares: give me a stable virtual IP, reachable at port 80, that forwards to port 8080 on whichever pods currently have the label `app: web`. Kubernetes continuously tracks which pods match that selector in an object called an **Endpoints** (or the newer **EndpointSlice**) resource — updated automatically the instant a matching pod is created, destroyed, or fails a readiness check — so the Service's membership list is always current without anyone manually maintaining it.

This is, conceptually, exactly a load balancer's job as Chapter 95 defined it: a single stable address in front of a changing pool of backends, with health-aware membership. Kubernetes's Service is that concept, implemented as a first-class, always-on cluster primitive rather than a separately deployed appliance.

---

## 10. ClusterIP: The Default Service Type

`ClusterIP` (Section 9's example, and the default `type` if none is specified) allocates a virtual IP address from a separate, cluster-internal range (distinct from the pod CIDR of Section 6) — reachable only from *within* the cluster, not from the outside internet. It exists purely to solve Section 8's problem for internal, service-to-service traffic: one stable address and DNS name (Section 15) that every other pod in the cluster can depend on indefinitely, no matter how many times the actual backing pods are rescheduled.

Critically, a ClusterIP is **not** a real interface anywhere, and no single process is "listening" on it the way a physical server listens on a bound socket (Chapter 57). It is a purely virtual address that only becomes meaningful because of the mechanism Section 11 describes.

---

## 11. How kube-proxy Actually Implements a ClusterIP

**`kube-proxy`** is the component, running on every node, responsible for making Section 10's virtual IP actually do something. Despite the name, in its most common modes it is not a proxy that traffic passes *through* in the traditional sense (Chapter 76's reverse proxy) — it is a program that continuously watches the Kubernetes API for Service and Endpoints changes and, on every change, rewrites the local node's packet-forwarding rules so the kernel itself does the redirection, with zero extra hops through a userspace process.

`kube-proxy` has run in three different modes over Kubernetes's history, each an evolution of the same idea:

| Mode | Mechanism | Trade-off |
|---|---|---|
| `userspace` (legacy, removed) | `kube-proxy` itself accepted connections and proxied them in userspace | Simple, but every packet paid a kernel-to-userspace-and-back cost |
| `iptables` (long-time default) | `kube-proxy` programs `iptables` `DNAT` rules (Chapter 103, Section 7's mechanism) so the kernel rewrites the ClusterIP to a real pod IP in-kernel | No userspace hop, but rule evaluation is roughly linear in the number of Services, which gets slow on very large clusters |
| `ipvs` (newer default option) | Uses the kernel's IP Virtual Server module — a purpose-built in-kernel load balancer — with hash-table-based lookups | Scales far better with Service count; supports more load-balancing algorithms than `iptables`'s round-robin-via-random-probability trick |

In `iptables` mode, a packet destined for a ClusterIP never actually gets forwarded to that IP at all — the moment it enters the sending pod's own network namespace, the kernel's `iptables` rules match the ClusterIP as a destination and rewrite (DNAT) it, in-flight, to one of the real backing pod IPs from the current Endpoints list, chosen essentially at random for basic load balancing. This is exactly Chapter 103 Section 7's DNAT mechanism, generalized from "one host port to one container" to "one cluster-wide virtual IP to a pseudo-random member of a dynamically updated pod set."

---

## 12. Hands-On: Watching kube-proxy's iptables Rules

```bash
# Find the ClusterIP kube-proxy assigned to the Service from Section 9.
kubectl get svc web
# NAME   TYPE        CLUSTER-IP     PORT(S)
# web    ClusterIP   10.96.34.201   80/TCP

# On any node, inspect the iptables rules kube-proxy generated for it
# (run this on a node, e.g. via `kubectl debug node/<name>` or SSH):
sudo iptables -t nat -L KUBE-SERVICES -n | grep 10.96.34.201
# -> KUBE-SVC-XXXX  tcp -- 0.0.0.0/0  10.96.34.201  tcp dpt:80

sudo iptables -t nat -L KUBE-SVC-XXXX -n
# -> a chain with one rule per backing pod, each with a statistical
#    probability (e.g. --probability 0.5000000000 for two pods) of
#    matching, DNAT-ing to that specific pod's real IP:targetPort
```

That `--probability` field is `iptables`'s literal mechanism for load balancing across an arbitrary number of backends purely with static, sequentially-evaluated rules: pod 1 matches with probability 1/N, and if it doesn't match, pod 2 matches with probability 1/(N-1) of what remains, and so on — a chained-coin-flip implementation of what a real load balancer (Chapter 95) would call round-robin or random selection.

---

## 13. NodePort: Reaching a Service From Outside the Cluster

A `ClusterIP` (Section 10) is unreachable from outside the cluster by design. `NodePort` builds directly on top of it, adding one more `kube-proxy`-managed rule:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  selector:
    app: web
  ports:
    - port: 80
      targetPort: 8080
      nodePort: 30080
  type: NodePort
```

`kube-proxy` opens port `30080` on **every node's** real network interface (not just the node actually running a matching pod), and any traffic arriving there gets DNAT'd through the exact same ClusterIP mechanism from Section 11, to a real pod, wherever in the cluster it happens to be running. Reaching the service from outside now means `curl http://<any-node-ip>:30080`, whether or not that specific node is running a `web` pod at all — `kube-proxy`'s rules forward it internally, across the pod network (Section 6), to whichever node actually is.

This is functionally similar to Chapter 103, Section 7's port publishing, generalized: instead of one host mapping one port to one container, every node in the cluster maps one port to the whole Service's current pod set.

---

## 14. LoadBalancer: Handing Off to Chapter 95's Load Balancers

`NodePort` (Section 13) technically makes a Service reachable from outside, but it's a poor fit for real external traffic: clients would need to know a specific node's real IP address, and if that node goes down, so does their access — exactly the single-point-of-failure problem Chapter 95 opened with.

`type: LoadBalancer` fixes this by asking the cluster's underlying infrastructure (a cloud provider like AWS, GCP, or Azure, or an on-premise solution like MetalLB) to provision a **real, external load balancer** — the actual hardware or managed-service load balancer Chapter 95 described in full — and configure it to forward traffic to the Service's `NodePort` on a healthy subset of the cluster's nodes.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  selector:
    app: web
  ports:
    - port: 80
      targetPort: 8080
  type: LoadBalancer
```

```bash
kubectl get svc web
# NAME   TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)
# web    LoadBalancer   10.96.34.201   203.0.113.50     80:30080/TCP
```

The `EXTERNAL-IP` here is a real, internet-routable address, owned and health-checked by the cloud provider's load balancer — which is itself just doing Chapter 95's job (distributing traffic across healthy backends, in this case the cluster's nodes) — layered on top of `NodePort` (Section 13), which is itself layered on top of `ClusterIP` (Section 10). Each Service type in this chapter is strictly additive: `LoadBalancer` implies `NodePort` implies `ClusterIP`, never a separate mechanism from scratch.

---

## 15. Service Discovery: DNS Names for Services

Every Service also gets a cluster-internal DNS name, resolved by **CoreDNS** (Kubernetes's built-in DNS server, itself running as ordinary pods), typically in the form `<service-name>.<namespace>.svc.cluster.local`:

```bash
kubectl exec -it some-other-pod -- nslookup web.default.svc.cluster.local
# Server:    10.96.0.10
# Address:   10.96.0.10#53
# Name:      web.default.svc.cluster.local
# Address:   10.96.34.201
```

This resolves straight to the Service's ClusterIP (Section 10) — the same Chapter 66–69 DNS mechanics this course already built in full depth, just running as an in-cluster resolver whose records are generated automatically from live Service objects instead of hand-edited zone files. Application code inside the cluster should always address other services by this DNS name, never by ClusterIP directly, for exactly the reason Chapter 66 gave for DNS existing at all: names are stable abstractions over addresses that change.

---

## 16. Headless Services: When You Want the Pod IPs Directly

Some workloads — most notably stateful systems like a database cluster where a client needs to talk to a *specific* replica, not a load-balanced random one — don't want Section 10's virtual-IP indirection at all. Setting `clusterIP: None` creates a **headless Service**: no virtual IP is allocated, and a DNS lookup on its name returns the **individual IPs of every matching pod directly**, as multiple `A`/`AAAA` records, rather than one virtual IP.

```bash
kubectl exec -it some-pod -- nslookup web-headless.default.svc.cluster.local
# Address: 10.244.1.5
# Address: 10.244.2.7
```

This is a deliberate opt-out of Section 11's `kube-proxy` DNAT mechanism entirely, useful whenever the application itself needs to implement its own logic for choosing among a known, stable set of named endpoints (each StatefulSet pod also gets its own stable per-pod DNS name in this mode) rather than accepting the ClusterIP's essentially random selection.

---

## 17. Ingress: HTTP-Aware Routing on Top of Services

Everything in Sections 10–14 operates at the transport layer (Chapter 57's ports, TCP/UDP) — a Service has no idea whether the traffic hitting it is HTTP, and can't route based on a URL path or `Host` header. **Ingress** is a separate Kubernetes API object that adds exactly that: HTTP/HTTPS-aware routing, typically implemented by an Ingress controller (commonly NGINX, Traefik, or a cloud provider's own) that itself runs as pods behind a `LoadBalancer` Service (Section 14).

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /users
            pathType: Prefix
            backend:
              service:
                name: users-svc
                port: { number: 80 }
          - path: /orders
            pathType: Prefix
            backend:
              service:
                name: orders-svc
                port: { number: 80 }
```

This is Chapter 76's reverse proxy pattern, expressed declaratively as a Kubernetes object: one external entry point, routing by `Host` header and path to whichever internal ClusterIP Service actually owns that piece of functionality — TLS termination (Chapter 82) is also commonly configured here, at the Ingress layer, rather than inside every individual pod.

---

## 18. Network Policies: Firewalling Pod-to-Pod Traffic

Section 3's Rule 2 guarantees every pod *can* reach every other pod by default — which is a connectivity guarantee, not a security posture. Left alone, a compromised pod anywhere in the cluster can attempt to reach any other pod anywhere else in the cluster. A **NetworkPolicy** object lets an operator restrict this, declaratively, to a firewall-like allow-list:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: only-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - port: 8080
```

This declares: pods labeled `app: backend` accept incoming traffic only from pods labeled `app: frontend`, on port 8080 — everything else is dropped. NetworkPolicy is only an API object; it takes real enforcement effort from the CNI plugin (Chapter 103) actually running the cluster's pod network. Not every CNI plugin implements it (the simplest ones, like basic Flannel, don't); Calico and Cilium (Chapter 105) both do, typically by compiling these rules down into `iptables`/`nftables` rules or, in Cilium's case, eBPF programs attached directly to each pod's network path.

---

## 19. Code: A Minimal ClusterIP-to-Endpoint Resolver in Go

To make Section 11's mechanism concrete in code (rather than just `iptables` output), a minimal Go program using `client-go` that watches a Service's Endpoints and prints the current backing pod IPs — conceptually the exact information `kube-proxy` consumes to generate its DNAT rules:

```go
package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// printServiceEndpoints mirrors, at a read-only level, exactly what
// kube-proxy watches continuously (Section 11) in order to know which
// real pod IPs currently sit behind a Service's stable ClusterIP.
func printServiceEndpoints(ctx context.Context, clientset *kubernetes.Clientset, namespace, service string) error {
	eps, err := clientset.CoreV1().Endpoints(namespace).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("fetching endpoints for %s/%s: %w", namespace, service, err)
	}

	for _, subset := range eps.Subsets {
		for _, addr := range subset.Addresses {
			for _, port := range subset.Ports {
				// This IP:port pair is exactly what an iptables DNAT
				// rule (Section 12) or an IPVS entry would send traffic
				// to, chosen essentially at random per new connection.
				fmt.Printf("backend: %s:%d\n", addr.IP, port.Port)
			}
		}
	}
	return nil
}

func main() {
	// In-cluster config: reads the ServiceAccount token Kubernetes
	// automatically mounts into every pod, the same credential kube-proxy
	// itself uses to watch the API server.
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	if err := printServiceEndpoints(context.Background(), clientset, "default", "web"); err != nil {
		panic(err)
	}
}
```

Run as a pod inside the cluster (with RBAC permission to read Endpoints), this prints the live backing pod IPs for the `web` Service from Section 9 — the same list `kube-proxy` uses, `CoreDNS` uses (indirectly, via the Service's ClusterIP), and any headless-Service DNS lookup (Section 16) would return directly.

---

## 20. Common Misconceptions

- **"A Kubernetes Service is a load balancer running somewhere."** For `ClusterIP` and `NodePort`, there is no separate proxy process traffic passes through — Section 11 was explicit that `kube-proxy`'s common modes program the kernel's own packet-rewriting rules; only `LoadBalancer` (Section 14) actually provisions a real, separate load-balancing appliance.
- **"Pods keep the same IP across restarts."** They generally do not (Section 8) — a rescheduled pod almost always gets a new IP, which is the entire reason Services and DNS names (Section 15), not raw pod IPs, are the correct thing for application code to depend on.
- **"NAT-free pod-to-pod communication means there's no NAT anywhere in a Kubernetes cluster."** Rule 2 (Section 3) only forbids NAT for pod-to-pod traffic specifically — `NodePort` and often outbound internet traffic from pods still involve NAT/masquerading at the node boundary, exactly as Chapter 103's container networking did.
- **"Ingress is just another name for a Service."** A Service (Sections 9–14) operates at the transport layer with no knowledge of HTTP; Ingress (Section 17) is a separate object providing HTTP-aware Layer 7 routing, typically implemented by a controller running behind its own `LoadBalancer` Service.
- **"NetworkPolicy is enforced by Kubernetes itself."** Kubernetes only stores and exposes the NetworkPolicy object; actual enforcement is entirely the responsibility of whichever CNI plugin the cluster runs, and some plugins simply ignore NetworkPolicy objects if they don't implement support for them.

---

## 21. Production Notes

- Large clusters (many thousands of Services) commonly move from `iptables` mode to `ipvs` mode (Section 11) specifically because `iptables` rule evaluation scales roughly linearly with Service count, becoming a measurable source of latency and CPU overhead at scale — a concrete, production-observed motivation for eBPF-based data planes (Chapter 105) that avoid this scaling pattern entirely.
- `LoadBalancer` Services on public clouds each typically provision one real cloud load balancer, which usually carries a direct dollar cost per Service — a common cost-optimization pattern is fronting many Services with one shared Ingress controller (Section 17) behind a single `LoadBalancer`, instead of one `LoadBalancer` Service per application.
- Pod restarts changing IPs (Section 8) is a frequent source of confusion in incident debugging — logs or monitoring dashboards captured "the IP of the failing pod" are often useless minutes later once Kubernetes has already rescheduled it; correlating by pod *name* or labels, not IP, is the production-correct habit.
- NetworkPolicy's default-allow-everything posture (Section 18) if no policies are defined at all is a common security gap in real clusters — many production security baselines require an explicit default-deny NetworkPolicy in every namespace as a starting point.

---

## 22. What's Simplified Here

The Service types (Sections 10, 13, 14), `kube-proxy`'s DNAT mechanism (Section 11), and DNS-based service discovery (Section 15) are described accurately and match current Kubernetes behavior. Left out for focus: `EndpointSlice`'s specific scalability improvements over the older, single-object `Endpoints` API; `Topology Aware Routing`, which lets `kube-proxy` prefer same-zone backends to reduce cross-zone traffic costs; `ExternalName` Services, a fourth Service type that's purely a DNS CNAME record with no proxying involved at all; and the full Ingress API's annotation-based extension mechanisms, which vary significantly between Ingress controller implementations despite sharing one core spec. The core idea — pod IPs are real and NAT-free, Services are a stable indirection layer implemented via kernel-level packet rewriting, and external access is layered additively through NodePort and LoadBalancer — is accurate and is exactly the substrate Chapter 105's eBPF-based CNIs optimize without changing.

---

## 23. Interview Questions & Model Answers

**Beginner: What are the three requirements Kubernetes places on any networking implementation?**
Every pod gets its own real, cluster-wide IP address; pods can reach each other's real IPs directly with no NAT in between; and any node can reach any pod on any node — collectively meaning the pod network must behave like one large, flat Layer 3 network regardless of the underlying physical topology.

**Beginner: Why can't application code safely hardcode a pod's IP address?**
Because pods are ephemeral — Kubernetes reschedules them constantly (crashes, scaling, node maintenance), and a recreated pod almost always receives a brand-new IP address, so any hardcoded IP will eventually point at nothing or, worse, a different pod entirely.

**Intermediate: What is the actual mechanism kube-proxy uses to make a ClusterIP work, in its common `iptables` mode?**
`kube-proxy` watches the Kubernetes API for Service and Endpoints changes and programs `iptables` DNAT rules on every node, so that any packet whose destination matches a Service's virtual ClusterIP gets its destination address rewritten, in-kernel, to one of the currently-live backing pod IPs — no separate proxy process handles the actual packets.

**Intermediate: How do ClusterIP, NodePort, and LoadBalancer relate to each other?**
They're additive, not alternative: `NodePort` opens a port on every node that forwards into the same ClusterIP mechanism; `LoadBalancer` provisions a real external load balancer (Chapter 95) that forwards traffic to that Service's NodePort — so a `LoadBalancer` Service always has a working ClusterIP and NodePort underneath it too.

**Advanced: Explain why Ingress exists as a separate object from Service, rather than Service simply supporting HTTP routing directly.**
A Service operates purely at the transport layer (Chapter 57) and has no concept of HTTP semantics like `Host` headers or URL paths; Ingress adds Layer 7-aware routing on top, letting one external entry point (and often one shared `LoadBalancer`, saving the per-Service cloud LB cost noted in Section 21) fan traffic out to many different backend Services based on HTTP request content, mirroring the general Layer 4 vs Layer 7 load balancing distinction Chapter 95 already drew.

**Advanced: A cluster runs a CNI plugin that satisfies Kubernetes's three networking-model requirements (Section 3) but does not implement NetworkPolicy enforcement. What is the practical security consequence, and why does this not violate Kubernetes's own conformance requirements?**
The practical consequence is that any pod can reach any other pod on any port, regardless of NetworkPolicy objects defined in the cluster, because NetworkPolicy is enforced entirely by the CNI plugin's data plane, not by Kubernetes core — this doesn't violate conformance because NetworkPolicy support has historically been optional for CNI plugins; Kubernetes's hard requirements are only the three connectivity rules from Section 3, and policy enforcement is treated as an additional, separately-adopted capability (which plugins like Calico and Cilium provide and simpler plugins like basic Flannel do not).

---

## 24. Exercises

### Easy
1. Referencing Section 3, explain in one sentence why Docker's default bridge networking (Chapter 103, Section 4) fails Kubernetes's pod networking requirements.
2. Create a `ClusterIP` Service in a local cluster (`kind`/`minikube`) and use `kubectl get endpoints` to see its current backing pod IPs.
3. What DNS name would a Service named `orders` in namespace `shop` resolve to, following Section 15's pattern?

### Medium
4. Using Section 12's approach, compare the `iptables` rules generated for a Service with one backing pod versus three backing pods, and explain how the `--probability` values change.
5. Convert the ClusterIP Service from Section 9 into a headless Service (Section 16) and show, with `nslookup`, how the DNS response changes.
6. Extend Section 19's Go program to also watch (not just fetch once) for Endpoints changes using `client-go`'s informer/watch API, printing a line every time the backing pod set changes.

### Hard
7. Design a NetworkPolicy (Section 18) that allows a `frontend`-labeled pod to reach a `backend`-labeled pod on port 8080, and a `backend`-labeled pod to reach a `database`-labeled pod on port 5432, but denies `frontend` from reaching `database` directly. Explain, referencing Section 18, what CNI-level mechanism would actually enforce this.
8. A `LoadBalancer` Service's `EXTERNAL-IP` stays `<pending>` indefinitely on a bare-metal cluster with no cloud provider integration. Referencing Section 14, explain why, and what a tool like MetalLB does to make `LoadBalancer` Services work in that environment.
9. Compare `iptables` mode and `ipvs` mode `kube-proxy` (Section 11) on a cluster with 5,000 Services, referencing Chapter 105's eBPF motivation: what specific per-packet cost does each mode incur, and why does an eBPF-based data plane avoid the scaling problem that motivated `ipvs`'s creation in the first place?

---

## 25. Summary and Bridge to Chapter 105

| Term | Meaning |
|---|---|
| Pod | One or more containers sharing a network namespace, Kubernetes's smallest deployable unit |
| Pod CIDR | The unique per-node address range a node's pods are assigned from |
| Service | A stable virtual IP/DNS name in front of a dynamic set of pods |
| ClusterIP | The default Service type: internal-only virtual IP |
| kube-proxy | The node-level component programming kernel rules to implement Service virtual IPs |
| NodePort | A Service type opening the same port on every node, forwarding into the ClusterIP mechanism |
| LoadBalancer | A Service type provisioning a real external load balancer (Chapter 95) in front of NodePort |
| Ingress | HTTP/HTTPS-aware Layer 7 routing on top of Services |
| NetworkPolicy | A declarative firewall-like object restricting pod-to-pod traffic, enforced by the CNI plugin |

This chapter showed Kubernetes demanding a stricter networking guarantee than Chapter 103's Docker default — real, NAT-free, cluster-wide pod IPs — and then layering Services, `kube-proxy`, and Chapter 95's load balancers on top to make that raw connectivity usable and stable despite constant pod churn. Every mechanism described here still assumes the same underlying tools Chapter 103 introduced: `iptables` rules, kernel routing, and CNI plugins doing the actual packet plumbing. Chapter 105 asks what happens when those `iptables`-based mechanisms stop scaling, and shows the kernel-level answer reshaping all of it: eBPF, the sandboxed-program technology that lets a CNI plugin like Cilium replace `kube-proxy`'s rule-based forwarding, and much of Chapters 102–104's packet path, with programs running directly inside the kernel's own packet-processing pipeline.
