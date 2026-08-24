# Chapter 89: Kubernetes for Blockchain Nodes, Overview

Chapter 87's `docker-compose.yml` comfortably runs three or four GoChain nodes on one server. But what happens when a testnet needs to run twenty nodes, spread across several machines, where individual containers occasionally crash and need to be restarted automatically, without anyone SSHing in at 3 a.m. to run `docker compose up` again by hand? This is the exact gap **Kubernetes** fills. This chapter is a conceptual overview, not a hard requirement — every later chapter in this volume keeps using Docker Compose, and you can read this chapter, understand the shape of the idea, and move on without ever running a Kubernetes cluster yourself.

## Table of Contents

1. [Where Docker Compose Stops Being Enough](#1-where-docker-compose-stops-being-enough)
2. [What Kubernetes Actually Is](#2-what-kubernetes-actually-is)
3. [Pods, Explained](#3-pods-explained)
4. [Deployments, Explained](#4-deployments-explained)
5. [Services, Explained](#5-services-explained)
6. [Why a Blockchain Node Is Not a Stateless Web App](#6-why-a-blockchain-node-is-not-a-stateless-web-app)
7. [The Whole Picture, Diagrammed](#7-the-whole-picture-diagrammed)
8. [A Sample GoChain Deployment Manifest](#8-a-sample-gochain-deployment-manifest)
9. [A Sample GoChain Service Manifest](#9-a-sample-gochain-service-manifest)
10. [Applying the Manifests](#10-applying-the-manifests)
11. [Namespaces, Briefly](#11-namespaces-briefly)
12. [What Kubernetes Buys You That Compose Does Not](#12-what-kubernetes-buys-you-that-compose-does-not)
13. [When to Actually Reach for Kubernetes](#13-when-to-actually-reach-for-kubernetes)
14. [Summary](#summary)
15. [Exercises](#exercises)

---

## 1. Where Docker Compose Stops Being Enough

Docker Compose, from Chapter 87, is genuinely excellent at one specific job: starting a fixed, known set of containers on **one machine**, wired together, with one command. Every service in `docker-compose.yml` lives on the same server, restarts are handled with a simple `restart:` policy, and scaling means editing the YAML file by hand to add another `nodeN:` block.

That model strains once a testnet grows past what a single server comfortably runs, or once you want nodes distributed across multiple physical machines for real geographic or organizational diversity — recall Chapter 51's point that concentrating every peer on one server undermines exactly the decentralization a blockchain network is supposed to have. Compose has no built-in concept of "a machine went down, reschedule its containers somewhere else" — if the one server running your `docker-compose.yml` crashes, every container on it is simply gone until a human intervenes.

**Kubernetes** (often abbreviated **k8s** — "k," eight letters, "s") is a system for running containers across a *cluster* of many machines at once, automatically deciding which machine runs which container, restarting containers that crash, and rescheduling them elsewhere if an entire machine disappears. Where Compose answers "how do I start several containers on this one server," Kubernetes answers "how do I keep several containers running correctly across a whole fleet of servers, indefinitely, without a human watching."

---

## 2. What Kubernetes Actually Is

Think of Docker Compose as a single restaurant kitchen: one head chef (you, running `docker compose up`) directly tells each cook (container) what to do, on one physical countertop (one server). Kubernetes is closer to running an entire restaurant *chain*: you do not personally walk into every branch and tell the staff what to cook. Instead, you write down a standing policy — "every branch should always have three pizza cooks on shift" — and a management layer (Kubernetes itself) continuously checks every branch, hires a replacement the moment a cook calls in sick, and can even open a new branch in a different city if demand grows, all without you personally supervising each one.

That "management layer" in Kubernetes is called the **control plane** — software that constantly compares "what you asked for" (declared in YAML files, much like `docker-compose.yml`) against "what is actually running," and takes action to close any gap. This is the single most important idea in this chapter: with Kubernetes, you do not issue one-time commands like `docker run`; you declare a **desired state** ("I want 3 copies of this GoChain node running, always"), and the control plane's job, forever, is making reality match that declaration.

```
   You declare:                     Kubernetes continuously enforces:
  "I want 3 GoChain node            +--------------------------------+
   pods running, always"    -->     | check every few seconds:       |
                                     |   actual state == desired?     |
                                     |     yes -> do nothing          |
                                     |     no  -> start/stop pods     |
                                     |            until it matches    |
                                     +--------------------------------+
```

---

## 3. Pods, Explained

A **pod** is Kubernetes' smallest deployable unit — one or more containers that are always scheduled together, on the same machine, sharing the same network address. For GoChain, a pod is almost always just one container: a single running GoChain node, the same image Chapter 86 built. (Pods can hold multiple tightly-coupled containers — a common pattern is a main application container plus a small "sidecar" container that does logging or monitoring for it — but GoChain does not need that complexity here.)

The important mental shift from Chapter 86's `docker run`: you do not create pods directly in normal use. Pods are disposable and expected to die — a crashed pod is simply deleted and a fresh one is created to replace it, with a new internal IP address and a clean filesystem (unless backed by persistent storage, the Kubernetes equivalent of Chapter 86's named Docker volumes). Nothing about a specific pod is meant to be precious; what is precious is the *declaration* that says how many should exist and what they should run.

---

## 4. Deployments, Explained

A **Deployment** is the object that actually expresses "keep N copies of this pod running, forever, and here is exactly what each one should look like." You write a Deployment manifest once, declaring the container image, the number of replicas (copies), resource limits, and so on — and the control plane takes it from there: if a pod crashes, the Deployment notices the replica count dropped below what was declared and starts a replacement immediately, with no human involved.

This maps directly onto Chapter 87's `node2`/`node3` services, just generalized: instead of hand-writing three nearly-identical `nodeN:` blocks in `docker-compose.yml`, a single Deployment with `replicas: 3` produces three identical, independently-scheduled GoChain node pods, and asking for a fourth is a one-line change (`replicas: 4`) rather than copy-pasting an entire new service block.

---

## 5. Services, Explained

Pods come and go, and each one gets a fresh internal IP address every time it is recreated — which is a problem, because Chapter 87's whole seed-node approach (`--seed node1:9000`) depends on a stable address to dial. A **Service** in Kubernetes solves exactly this: it is a stable, unchanging network name and address that sits in front of a group of pods (selected by a label, a simple key-value tag attached to each pod) and automatically routes traffic to whichever pods currently exist and are healthy, even as individual pods are replaced underneath it.

This is the direct Kubernetes equivalent of Chapter 87's Docker Compose service-name DNS (`node1` resolving automatically to whatever container is currently playing that role) — the mechanism is different under the hood, but the effect for GoChain's own code is identical: it dials a stable name, and the platform handles making sure something real answers on the other end.

---

## 6. Why a Blockchain Node Is Not a Stateless Web App

Sections 3-5 quietly borrowed the standard Kubernetes mental model, which was designed first and foremost for **stateless** web applications — pods where any replica can answer any request, none of them holding data the others do not also have, and where losing one and replacing it with a fresh, empty one is a complete non-event. A plain Deployment, exactly as shown in Section 7 below, is genuinely the right tool for that shape of workload.

A GoChain node does not quite fit that shape, and it is worth being honest about the mismatch rather than gliding past it. Each node holds its own BoltDB file (Chapter 54), its own view of the mempool (Chapter 34), and — critically — its own identity as a specific peer other nodes have specifically connected to and are tracking in their peer tables (Chapter 46). Delete one of three ordinary Deployment pods and Kubernetes creates a replacement with a brand-new internal identity and, depending on your storage configuration, a completely empty data directory — indistinguishable, from the outside, from a node that just lost its entire chain history and needs to resync from scratch every single time it is rescheduled.

Kubernetes has a purpose-built object for exactly this situation: a **StatefulSet**. Where a Deployment's pods are interchangeable, a StatefulSet gives each pod a stable, predictable name (`gochain-node-0`, `gochain-node-1`, `gochain-node-2`, numbered and never reshuffled) and, paired with a `volumeClaimTemplate`, its own dedicated persistent storage that follows that specific numbered pod even if it is rescheduled onto a completely different physical machine. `gochain-node-1` deleted and recreated comes back as `gochain-node-1` again, with the same disk attached, not a fresh, anonymous replacement.

```
   Deployment (Section 4)                  StatefulSet
  +---------------------------+           +---------------------------+
  | pod: random-suffix-a       |           | pod: gochain-node-0        |
  | pod: random-suffix-b       |           | pod: gochain-node-1        |
  | pod: random-suffix-c       |           | pod: gochain-node-2        |
  |                             |           |                             |
  | delete one -> replacement   |           | delete one -> replacement   |
  | gets a NEW random name and  |           | gets the SAME name and      |
  | (usually) fresh storage     |           | (via volumeClaimTemplate)   |
  |                             |           | the SAME storage reattached |
  +---------------------------+           +---------------------------+
```

This chapter's own sample manifests in Sections 8-9 deliberately use the simpler Deployment shape, both because it is the more broadly useful concept to understand first and because Section 7's `emptyDir` storage was already flagged as a simplification. A production-minded GoChain-on-Kubernetes setup would reach for a StatefulSet the moment node identity and data locality actually matter — which, for a blockchain node, is essentially always. This is exactly the kind of nuance worth knowing exists, even in an overview chapter that stops short of implementing it: recognizing "this workload has state and identity, therefore Deployment alone is the wrong tool" is a transferable piece of judgment, useful for any stateful system you deploy on Kubernetes in the future, not just GoChain.

---

## 7. The Whole Picture, Diagrammed

```
                         Kubernetes Cluster
     +-----------------------------------------------------------+
     |                                                             |
     |   Deployment: gochain-node (declares: replicas: 3)          |
     |     |                                                       |
     |     |  control plane creates/replaces pods to match          |
     |     v                                                       |
     |  +--------+   +--------+   +--------+                       |
     |  |  Pod   |   |  Pod   |   |  Pod   |   <- each pod runs     |
     |  | node   |   | node   |   | node   |      one gochain       |
     |  | :8080  |   | :8080  |   | :8080  |      container         |
     |  | :9000  |   | :9000  |   | :9000  |                       |
     |  +--------+   +--------+   +--------+                       |
     |       ^            ^            ^                           |
     |       |            |            |                           |
     |       +------------+------------+                           |
     |                    |                                        |
     |         Service: gochain-node-svc                            |
     |         (stable name + address, routes to whichever          |
     |          pods currently exist and are healthy)                |
     |                    |                                        |
     +--------------------|----------------------------------------+
                           |
                     external traffic
                    (a wallet, a peer,
                     another node)
```

A pod can crash and be replaced three times in an hour, on three different underlying machines in the cluster, and from outside the cluster — from a `gochain-wallet` or another peer trying to connect — none of that churn is ever visible. They only ever see the Service's one stable address.

---

## 8. A Sample GoChain Deployment Manifest

Kubernetes manifests, like `docker-compose.yml`, are YAML files — the vocabulary is different, but the format itself holds no surprises. Here is a Deployment that runs three replicas of the exact same GoChain image built in Chapter 86:

```yaml
# gochain-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gochain-node
  labels:
    app: gochain-node
spec:
  replicas: 3 # the desired state: always keep 3 pods running
  selector:
    matchLabels:
      app: gochain-node # tells the Deployment which pods it owns
  template:
    # everything under `template` describes what each pod should
    # look like - this is essentially one container spec, repeated
    # `replicas` times.
    metadata:
      labels:
        app: gochain-node # must match `selector.matchLabels` above
    spec:
      containers:
        - name: gochain
          image: ghcr.io/you/gochain:latest # the image from Chapter 92's CI/CD
          ports:
            - containerPort: 8080 # API
            - containerPort: 9000 # P2P
          args:
            - "node"
            - "start"
            - "--api-addr=0.0.0.0:8080"
            - "--p2p-addr=0.0.0.0:9000"
            - "--seed=gochain-node-svc:9000"
            - "--data=/app/data"
          volumeMounts:
            - name: gochain-data
              mountPath: /app/data
          resources:
            # requests/limits tell Kubernetes how much CPU and memory
            # to reserve, so it can schedule pods onto machines that
            # actually have room, and so one misbehaving pod cannot
            # starve every other pod on the same physical node.
            requests:
              cpu: "100m" # 100 millicores = one-tenth of one CPU core
              memory: "128Mi"
            limits:
              cpu: "500m"
              memory: "256Mi"
      volumes:
        - name: gochain-data
          emptyDir: {} # simplest option; a real cluster would use a
                       # PersistentVolumeClaim so data survives a pod
                       # being rescheduled onto a different machine
```

Notice `--seed=gochain-node-svc:9000`: every pod points at the *Service's* stable name from Section 9 below, not at any individual pod's own address, for exactly the reason Section 5 explained — an individual pod's address is not something to depend on.

The `args:` list above hardcodes every flag directly into the manifest, which is fine for a small overview but becomes awkward the moment several Deployments need to share most of the same flags with only one or two differing. Kubernetes' answer to that is a **ConfigMap** — a small, separately-managed object holding configuration values, which a pod can then reference instead of repeating the same literal strings in every manifest:

```yaml
# gochain-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gochain-config
data:
  API_ADDR: "0.0.0.0:8080"
  P2P_ADDR: "0.0.0.0:9000"
  SEED_ADDR: "gochain-node-svc:9000"
```

A pod spec then reads these as environment variables (`envFrom: - configMapRef: { name: gochain-config }`) instead of a fixed `args:` list, the same underlying idea as keeping configuration in a `.env` file rather than hardcoding it into source code — a pattern this course has not needed until now, since Chapter 87's Compose file was small enough that repeating flags per-service stayed perfectly readable.

---

## 9. A Sample GoChain Service Manifest

```yaml
# gochain-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: gochain-node-svc
spec:
  selector:
    app: gochain-node # routes to every pod carrying this label -
                       # exactly the pods the Deployment above creates
  ports:
    - name: api
      port: 8080
      targetPort: 8080
    - name: p2p
      port: 9000
      targetPort: 9000
  type: ClusterIP # reachable only from inside the cluster by default;
                   # a real public-facing deployment would use
                   # `type: LoadBalancer` instead, to get a real
                   # external, internet-reachable IP address
```

`type: ClusterIP` (the default) gives the Service a stable address usable only by other things running inside the same cluster — fine for pods talking to each other, but not yet reachable from the public internet the way Chapter 88's VM was. `type: LoadBalancer` is the setting that would ask the cloud provider running the cluster to provision a real, public IP address in front of the Service, the Kubernetes equivalent of Chapter 88's `ufw allow` rules — a decision deliberately left for later, since this chapter's goal is understanding the shape of the system, not standing up a production cluster.

---

## 10. Applying the Manifests

If you do want to try this hands-on, the lowest-friction way to get a real (if tiny, single-machine) Kubernetes cluster is a local tool like **minikube** or **kind** (Kubernetes IN Docker), both of which run an entire small cluster inside your own laptop for learning purposes, with no cloud provider or billing involved at all.

```bash
# Using minikube as an example - installs a single-node local cluster.
minikube start

# kubectl is Kubernetes' command-line tool, the same conceptual role
# `docker` plays for plain containers.
kubectl apply -f gochain-deployment.yaml
kubectl apply -f gochain-service.yaml

# Confirm three pods came up.
kubectl get pods
# NAME                            READY   STATUS    RESTARTS   AGE
# gochain-node-7d9f4c8b6d-4k2pl   1/1     Running   0          12s
# gochain-node-7d9f4c8b6d-8xqwr   1/1     Running   0          12s
# gochain-node-7d9f4c8b6d-jm3nz   1/1     Running   0          12s

# Confirm the Service exists and has an internal address.
kubectl get service gochain-node-svc
```

Deliberately kill one pod and watch the Deployment replace it automatically — the single clearest hands-on demonstration of "declared state, continuously enforced" this chapter can offer:

```bash
kubectl delete pod gochain-node-7d9f4c8b6d-4k2pl
# pod "gochain-node-7d9f4c8b6d-4k2pl" deleted

kubectl get pods
# a brand-new pod, with a new name, is already Running - the Deployment
# noticed the replica count dropped to 2 and immediately created a
# replacement to bring it back to the declared 3.
```

---

## 11. Namespaces, Briefly

One more organizing concept worth a short mention before the final wrap-up: a **namespace** is a way of partitioning a single Kubernetes cluster into several isolated virtual clusters, each with its own set of named objects. Two Deployments both named `gochain-node` can coexist peacefully in the same physical cluster as long as one lives in namespace `testnet-staging` and the other in `testnet-production` — Kubernetes treats them as entirely separate objects, the same way two files named `main.go` can coexist in two different directories without conflict.

```bash
kubectl create namespace gochain-staging
kubectl apply -f gochain-deployment.yaml --namespace gochain-staging
kubectl get pods --namespace gochain-staging
```

For a single small testnet, namespaces are not strictly necessary — everything can live in Kubernetes' `default` namespace with no ill effect. They become genuinely useful the moment you want to run more than one environment on the same cluster at once (echoing Chapter 92's Exercise 9, which sketched exactly this idea — a staging deployment alongside a production one), giving each one a clean, separate name space to avoid accidental collisions rather than relying on careful naming discipline alone.

---

## 12. What Kubernetes Buys You That Compose Does Not

Concretely, compared to Chapter 87's `docker-compose.yml`:

- **Multi-machine clusters.** A Kubernetes cluster can span dozens or hundreds of physical machines; Compose is fundamentally scoped to one.
- **Self-healing by default.** A pod that crashes is replaced automatically, with no `restart:` policy to remember to set (though Compose's `restart: unless-stopped` covers the single-machine version of this reasonably well).
- **Declarative scaling.** `kubectl scale deployment gochain-node --replicas=10` is the entire command to go from 3 nodes to 10 — no hand-editing YAML to add seven new near-duplicate service blocks.
- **Rolling updates with zero downtime.** A Deployment can replace pods with a new image version a few at a time, keeping the rest of the fleet serving traffic throughout — directly relevant to Chapter 92's Exercise 8 about zero-downtime deploys, which Compose's `pull && up -d` approach cannot fully achieve on its own.
- **Scheduling across failure domains.** A properly configured cluster can spread pods across different physical machines (and even different data centers) so a single hardware failure does not take down every node at once — a real, infrastructure-level version of the peer-diversity idea from Chapter 51.

---

## 13. When to Actually Reach for Kubernetes

None of this makes Kubernetes automatically the right choice. It adds a meaningful new layer of concepts (pods, Deployments, Services, StatefulSets, namespaces, and several more this overview does not cover at all — ConfigMaps, Secrets, Ingress controllers) that a small testnet simply does not need in order to work. Every remaining chapter in this volume — the faucet, monitoring, CI/CD, TLS, backups, and the final capstone — is written entirely on top of Chapter 87's Docker Compose stack, and none of it requires Kubernetes at all.

A reasonable rule of thumb, echoing Chapter 80's consensus/architecture decision framework: reach for Kubernetes when you are actually operating enough machines, and enough churn across them, that manually running `docker compose up` on each one has become the bottleneck — not because a real production blockchain project "should" use Kubernetes on principle. A three-to-ten-node learning testnet, which is exactly what this volume builds, is comfortably inside Docker Compose's sweet spot.

---

## Summary

- Kubernetes automates running containers across a *cluster* of many machines, where Docker Compose is scoped to just one.
- The core idea is **declarative desired state**: you declare what should be running, and a control plane continuously reconciles reality to match it, rather than you issuing one-time `docker run` commands.
- A **pod** is the smallest deployable unit — usually one GoChain node container — and pods are treated as disposable, expected to be replaced freely.
- A **Deployment** declares how many replicas of a pod should exist and what each one should look like; it recreates crashed pods automatically to restore that count.
- A **Service** gives a stable, unchanging name and address in front of a group of pods, solving the same "how do peers find each other reliably" problem Chapter 87's Docker DNS solved for Compose.
- A sample `Deployment` + `Service` pair for GoChain looks almost exactly like `docker-compose.yml`'s services, just in Kubernetes' own YAML vocabulary, with `--seed` pointed at the Service's stable name instead of an individual node.
- A blockchain node is not a stateless web app — a **StatefulSet** gives each pod a stable, numbered identity and its own persistent storage that follows it across rescheduling, unlike a plain Deployment's interchangeable pods.
- **Namespaces** partition one physical cluster into isolated virtual clusters, useful once you want more than one environment (staging, production) coexisting on the same hardware.
- Kubernetes buys you multi-machine scheduling, self-healing, one-line scaling, and zero-downtime rolling updates — real advantages, but ones that matter most at a scale this course's testnet does not require.
- This chapter is optional background: every later chapter in this volume continues to build directly on Chapter 87's Docker Compose stack.

---

## Exercises

### Easy

1. Install `minikube` (or `kind`) locally, run `minikube start`, and confirm `kubectl get nodes` shows one node in a `Ready` state.

2. Apply the `gochain-deployment.yaml` and `gochain-service.yaml` manifests from Sections 8-9 to your local cluster, and confirm `kubectl get pods` shows three `gochain-node` pods in a `Running` state.

3. Delete one pod with `kubectl delete pod <name>` and, within a few seconds, run `kubectl get pods` again. Confirm a replacement pod with a new name has already appeared, and write one sentence describing what component of Kubernetes caused that to happen.

### Medium

4. Scale the Deployment from 3 replicas to 6 using `kubectl scale deployment gochain-node --replicas=6`, and confirm `kubectl get pods` reflects the new count within a few seconds. Then scale back down to 2 and confirm Kubernetes terminates the extra pods rather than leaving them running.

5. Change the Service's `type` from `ClusterIP` to `NodePort` in `gochain-service.yaml`, reapply it, and use `kubectl get service gochain-node-svc` to find the automatically assigned external port. Confirm you can `curl` a GoChain node's `/chain/height` endpoint from your host machine through that port.

6. Add a `livenessProbe` to the Deployment manifest that periodically checks each pod's `/chain/height` endpoint over HTTP, and confirm (by breaking a pod's networking with `kubectl exec` in a controlled way) that Kubernetes detects an unhealthy pod and restarts it, distinguishing this from the ordinary crash-and-replace behavior you observed in Exercise 3.

### Hard

7. Convert `gochain-deployment.yaml`'s storage from `emptyDir: {}` to a `PersistentVolumeClaim`, so a pod's blockchain data survives it being deleted and recreated rather than starting from genesis every time. Test this explicitly: mine a few blocks, delete the pod, and confirm the replacement pod's chain height is unchanged.

8. Research and explain, in 200-300 words, why running three GoChain nodes as three replicas of the *same* Deployment (as this chapter does) is actually a subtly wrong model compared to Chapter 87's Compose setup, where `node1` plays a distinct seed-node role from `node2`/`node3`. Propose a manifest structure (hint: a separate, single-replica Deployment for the seed node, plus a scalable Deployment for ordinary peers) that fixes this.

9. Install a real Kubernetes Ingress controller in your local cluster, and configure an Ingress resource that routes `gochain.local` to the `gochain-node-svc` Service's API port, using your machine's `/etc/hosts` file to make `gochain.local` resolve locally. Compare, in writing, what an Ingress controller is doing here to what Chapter 93's Caddy reverse proxy does for a plain Docker Compose deployment — are they solving the same problem at a different layer, or genuinely different problems?

10. Rewrite `gochain-deployment.yaml` as a `StatefulSet` instead of a `Deployment`, using a `volumeClaimTemplate` for each pod's `/app/data` directory, following Section 6's reasoning. Confirm, via `kubectl get pods`, that the resulting pods are named predictably (`gochain-node-0`, `gochain-node-1`, `gochain-node-2`), then delete `gochain-node-1` specifically and confirm its replacement comes back with the same name and the same chain data intact, unlike the fresh, differently-named replacement pods you observed with the plain Deployment in Exercise 3.
