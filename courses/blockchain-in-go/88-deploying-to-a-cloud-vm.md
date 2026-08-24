# Chapter 88: Deploying to a Cloud VM

Every GoChain testnet you have run so far — Chapter 86's single container, Chapter 87's whole `docker compose up` network — has lived entirely on your own laptop. That is a fine place to develop, but it is not a network anyone else can reach: the moment you close your laptop's lid, the "testnet" disappears, and nobody outside your Wi-Fi network could ever have dialed in to it anyway. This chapter moves GoChain onto a real, internet-reachable computer — a **cloud VM** — that stays running whether or not your laptop is even turned on, and configures a firewall so that only the ports GoChain actually needs are open to the world.

## Table of Contents

1. [What a Cloud VM Actually Is](#1-what-a-cloud-vm-actually-is)
2. [Choosing a Provider — the Concepts Transfer Everywhere](#2-choosing-a-provider--the-concepts-transfer-everywhere)
3. [Provisioning Your First VM](#3-provisioning-your-first-vm)
4. [First Login and Basic Hardening](#4-first-login-and-basic-hardening)
5. [Installing Docker on the VM](#5-installing-docker-on-the-vm)
6. [Adding Swap Space for a Small VM](#6-adding-swap-space-for-a-small-vm)
7. [Getting Your Code Onto the Server](#7-getting-your-code-onto-the-server)
8. [Configuring the Firewall](#8-configuring-the-firewall)
9. [Bringing the Testnet Up on the VM](#9-bringing-the-testnet-up-on-the-vm)
10. [Making the Stack Survive a Reboot](#10-making-the-stack-survive-a-reboot)
11. [Verifying Real Internet Reachability](#11-verifying-real-internet-reachability)
12. [What to Do Before You Walk Away](#12-what-to-do-before-you-walk-away)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. What a Cloud VM Actually Is

A **virtual machine (VM)** is a computer that exists as software, running on a much larger physical machine owned by a cloud provider, but that behaves — from the inside — exactly like a dedicated computer with its own CPU allocation, its own memory, its own disk, and its own IP address. A **cloud VM** is simply one of these you rent, by the hour or by the month, from a company that owns the physical hardware in a data center somewhere. You never touch the physical machine; you interact with your slice of it entirely over the network, usually starting with SSH (Secure Shell, a protocol for securely logging into and running commands on a remote computer).

Think of it like renting a serviced apartment instead of building a house. You do not own the building, the plumbing, or the electrical wiring — the provider maintains all of that. What you get is a private, locked space with a working address (an IP address, in this case) that is genuinely yours to configure, furnish, and lock behind your own door. Just as a serviced apartment can be rented for an afternoon or a year, a cloud VM can be spun up for a quick experiment and destroyed an hour later, billed only for the time it existed.

This is the single property that matters most for this chapter: unlike your laptop, a cloud VM has a **public IP address** — a network address reachable from any other computer on the internet, not just other devices on your home Wi-Fi. That public reachability is precisely what turns "a testnet running on my laptop" into "a testnet a stranger on the other side of the world can join," which is the entire point of Chapter 90's public testnet.

---

## 2. Choosing a Provider — the Concepts Transfer Everywhere

There are many cloud providers — DigitalOcean, AWS (Amazon Web Services), Google Cloud, Linode, Hetzner, and others — and this chapter deliberately does not lock you into one. Every provider's dashboard looks slightly different, but every single one asks you to make the same handful of decisions, using different names for the same ideas:

| Concept | DigitalOcean's name | AWS's name | What it means |
|---|---|---|---|
| The VM itself | Droplet | EC2 instance | The virtual computer you are renting |
| The template it boots from | Image / distribution | AMI (Amazon Machine Image) | Which operating system it starts with |
| Its size | Droplet size (e.g. Basic, 1 vCPU/1GB) | Instance type (e.g. `t3.micro`) | How much CPU, RAM, and disk you get |
| Network access rules | Firewall | Security Group | Which ports are open to the internet |
| Your login credential | SSH key | Key pair | The private key that lets you log in without a password |

If you have never used any cloud provider before, DigitalOcean is the friendliest on-ramp — its dashboard is small, its pricing is a flat monthly rate with no surprise line items, and a Droplet capable of running this chapter's testnet costs a few dollars a month. AWS EC2 is more powerful and more commonly used professionally, but its dashboard has considerably more surface area to learn before you find the three or four settings that actually matter for this chapter. Whichever you pick, every step below maps directly onto the equivalent screen or command for your provider — the underlying Linux server you end up with is the same either way.

---

## 3. Provisioning Your First VM

**Provisioning** just means "requesting that the provider create the VM for you, according to the choices you made." Concretely, for essentially any provider, this means:

1. **Create an account** and add a payment method (most providers require this even if your usage stays within a free tier).
2. **Generate an SSH key pair** on your own laptop, if you do not already have one:

   ```bash
   ssh-keygen -t ed25519 -C "gochain-testnet"
   # Creates ~/.ssh/id_ed25519 (private key — never share this)
   # and ~/.ssh/id_ed25519.pub (public key — safe to upload anywhere)
   ```

   This is the exact same key-pair idea from Chapter 11 — a private half you keep secret, a public half you can hand to anyone. Here, instead of signing transactions, the private key proves to the VM that you are allowed to log in, and the public key is what you upload to the provider so it can be installed on the VM automatically at creation time.

3. **Choose an image**: pick a plain Linux distribution — Ubuntu 22.04 LTS is a safe, extremely well-documented default that every provider offers.
4. **Choose a size**: a testnet of three or four lightweight GoChain containers (recall Chapter 86's ~20MB images) runs comfortably on the smallest paid tier most providers offer — typically 1 vCPU and 1-2GB of RAM. You can always resize later if the testnet grows.
5. **Attach your SSH key** to the new VM during creation — this is the step that makes password-free login work the moment it boots.
6. **Create the VM**, and note the **public IP address** the provider assigns it. This one number is the address of your entire testnet from this point forward, in every chapter that follows.

```
   Your laptop                         Cloud provider's data center
  +-------------+                     +------------------------------+
  | ~/.ssh/     |   "create a VM,     |  A new Linux server is        |
  |  id_ed25519 |    install my       |  provisioned, given a public  |
  |  (private)  |    public key"      |  IP (e.g. 203.0.113.42), and  |
  | id_ed25519. |  ---------------->  |  your public key is installed |
  |  pub (public)|                    |  in its authorized_keys file  |
  +-------------+                     +------------------------------+
```

---

## 4. First Login and Basic Hardening

Log in over SSH using the public IP address from Section 3:

```bash
ssh root@203.0.113.42
```

Most providers hand you a fresh VM logged in as `root` — the all-powerful administrator account. Running everything as `root` forever is a bad habit: any mistake, or any bug in software you later run, has unrestricted power over the entire machine. The very first thing to do is create a dedicated, less-privileged user for day-to-day work:

```bash
# Create a new user with a home directory and add it to the sudo group,
# which lets it run individual commands as root via `sudo` when actually
# needed, without staying logged in as root all the time.
adduser deploy
usermod -aG sudo deploy

# Copy your SSH key to the new user so you can log in as `deploy` the
# same password-free way you just logged in as root.
rsync --archive --chown=deploy:deploy ~/.ssh /home/deploy
```

From this point on, log in as `deploy`, not `root`:

```bash
ssh deploy@203.0.113.42
```

The name `deploy` is deliberate — this is the exact account Chapter 92's CI/CD pipeline will SSH into automatically (as `secrets.TESTNET_SSH_USER`) to push new releases, so it is worth setting it up correctly here rather than improvising a different account later.

While you are here, run the two commands that every fresh Linux server should get immediately:

```bash
sudo apt update && sudo apt upgrade -y
```

This pulls the latest security patches for the base operating system — a freshly provisioned VM's image is often a few weeks or months old by the time you boot it, and this closes that gap.

---

## 5. Installing Docker on the VM

Everything from Chapter 86 onward assumes Docker exists. Ubuntu's official convenience script is the fastest reliable way to install it on a fresh server:

```bash
# Docker's official install script detects your distribution and
# installs the correct packages automatically - simpler than manually
# adding Docker's package repository yourself.
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Let the `deploy` user run `docker` commands without needing `sudo`
# every single time, by adding it to the `docker` group.
sudo usermod -aG docker deploy

# Log out and back in for the group change to take effect.
exit
ssh deploy@203.0.113.42

# Confirm Docker actually works.
docker run hello-world
```

`docker run hello-world` pulls a tiny test image, runs it, and prints a confirmation message — if you see that message, Docker is correctly installed and the `deploy` user has permission to use it, exactly the two things every later step in this chapter depends on.

Also install the Docker Compose plugin, if the convenience script did not already include it (it usually does on recent Ubuntu versions):

```bash
docker compose version
# Docker Compose version v2.27.0
```

---

## 6. Adding Swap Space for a Small VM

The cheapest VM tier from Section 3 — 1 vCPU, 1-2GB of RAM — is genuinely enough to run Chapter 87's whole testnet, but it leaves little room for error. Building four Docker images in a row (three nodes plus the explorer, per Chapter 87's `docker-compose.yml`) can briefly spike memory usage during compilation, and a VM with no **swap space** simply kills the offending process the instant physical RAM runs out, via Linux's out-of-memory ("OOM") killer — often taking down something unrelated and important along with it, since the OOM killer picks a victim process using its own heuristics, not necessarily the one that actually caused the spike.

**Swap space** is disk space set aside to act as an overflow for RAM: when physical memory is nearly full, the kernel moves rarely-used pages of memory out to disk, freeing up real RAM for whatever needs it right now. It is much slower than real RAM, so it is not a substitute for having enough memory in the first place — but as a shock absorber against short, occasional spikes on a small VM, it is a cheap and standard precaution:

```bash
# Create a 2GB file to use as swap space. `fallocate` reserves the
# space instantly, without slowly zeroing every byte the way `dd`
# would.
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile   # only root should be able to read this file

# Format it as swap and turn it on.
sudo mkswap /swapfile
sudo swapon /swapfile

# Confirm it's active.
free -h
#                total        used        free      shared  buff/cache   available
# Mem:           985Mi       210Mi       520Mi       1.0Mi       255Mi       650Mi
# Swap:          2.0Gi          0B       2.0Gi
```

Make it permanent, so it survives a reboot rather than needing to be recreated by hand every time:

```bash
echo "/swapfile none swap sw 0 0" | sudo tee -a /etc/fstab
```

This single step has saved more small-VM deployments from a mysterious, hard-to-diagnose crash than almost any other piece of advice in this chapter — worth doing before, not after, the first time a `docker compose up --build` unexpectedly hangs or a container disappears with no clear error in its own logs.

---

## 7. Getting Your Code Onto the Server

The VM needs the same two files Chapter 86 and Chapter 87 already produced: the `Dockerfile` (and `Dockerfile.explorer`) and `docker-compose.yml`, plus the rest of the `gochain` source tree the multi-stage build compiles from. The cleanest way to get them there — and the one Chapter 92's CI/CD pipeline continues to assume — is a Git clone directly onto the server, into the fixed path `/opt/gochain`:

```bash
# /opt is the conventional Linux location for "optional," self-contained
# application software - a sensible, discoverable home for GoChain's
# deployed copy, distinct from the deploy user's own home directory.
sudo mkdir -p /opt/gochain
sudo chown deploy:deploy /opt/gochain

git clone https://github.com/you/gochain.git /opt/gochain
cd /opt/gochain
```

If your repository is private, either use a GitHub personal access token in the clone URL, or generate a fresh SSH key *on the VM itself* and add its public half as a GitHub deploy key — never copy your laptop's private key onto a server, since a compromised server should never be able to compromise your GitHub account too.

---

## 8. Configuring the Firewall

A firewall is a set of rules that decides which network connections are allowed to reach the machine at all, checked before the connection ever gets to Docker, GoChain, or anything else running on the VM. Ubuntu ships with **ufw** (Uncomplicated Firewall), a simple front end over Linux's underlying `iptables` rules; AWS instead calls the equivalent concept a **security group**, configured from its dashboard rather than a command line — the rules you are about to set are identical in spirit either way.

The rule that matters most: **open only the ports GoChain actually needs, plus SSH, and deny everything else.** Recall from Chapter 86 that GoChain conventionally uses port 8080 for its API and port 9000 for P2P — those, plus port 22 for your own SSH access, are the entire list.

```bash
# Deny everything by default - the safest possible starting posture.
# Every rule after this is an explicit exception to that default.
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow SSH so you do not lock yourself out of your own server.
sudo ufw allow 22/tcp

# GoChain's two conventional ports: 8080 for the API, 9000 for P2P.
sudo ufw allow 8080/tcp
sudo ufw allow 9000/tcp

# Turn the firewall on.
sudo ufw enable
# Command may disrupt existing ssh connections. Proceed with operation (y|n)? y

# Confirm exactly what is open.
sudo ufw status verbose
```

Expected output from that last command:

```
Status: active
Logging: on (low)
Default: deny (incoming), allow (outgoing), disabled (routed)

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW IN    Anywhere
8080/tcp                   ALLOW IN    Anywhere
9000/tcp                   ALLOW IN    Anywhere
22/tcp (v6)                ALLOW IN    Anywhere (v6)
8080/tcp (v6)              ALLOW IN    Anywhere (v6)
9000/tcp (v6)              ALLOW IN    Anywhere (v6)
```

Notice what is deliberately **not** open: Chapter 87's `docker-compose.yml` also publishes `node2` on 8081/9001, `node3` on 8082/9002, and the explorer on 8090. Those stay unreachable from the public internet for now — only `node1`, the testnet's seed node on the conventional 8080/9000 pair, is exposed. This is intentional: a stranger joining the network in Chapter 90 only ever needs one reachable address to bootstrap from, and every extra open port is one more thing an attacker could probe. Later chapters open more access deliberately and narrowly: Chapter 91's Grafana and Prometheus stay behind an SSH tunnel rather than a public port, and Chapter 93 puts the explorer behind a proper domain and HTTPS reverse proxy instead of a bare, unauthenticated port.

On AWS EC2, the equivalent security-group configuration is three inbound rules — TCP 22, TCP 8080, TCP 9000, each from source `0.0.0.0/0` (meaning "any IP address") — set from the EC2 console's "Security Groups" screen rather than a shell command, with everything else denied by the security group's own default.

---

## 9. Bringing the Testnet Up on the VM

With Docker installed, the code cloned, and the firewall configured, starting the testnet is the exact same command as Chapter 87, just run over an SSH session instead of on your laptop:

```bash
cd /opt/gochain
docker compose up -d --build
```

`-d` runs it detached, which matters far more here than it did locally — you want the testnet to keep running after you close your SSH session and disconnect, not die the moment your terminal does. Confirm everything started:

```bash
docker compose ps
```

```
NAME                IMAGE            STATUS         PORTS
gochain-node1       gochain-node1    Up 8 seconds   0.0.0.0:8080->8080/tcp, 0.0.0.0:9000->9000/tcp
gochain-node2       gochain-node2    Up 7 seconds   0.0.0.0:8081->8080/tcp, 0.0.0.0:9001->9000/tcp
gochain-node3       gochain-node3    Up 7 seconds   0.0.0.0:8082->8080/tcp, 0.0.0.0:9002->9000/tcp
gochain-explorer    gochain-explorer Up 6 seconds   0.0.0.0:8090->8090/tcp
```

---

## 10. Making the Stack Survive a Reboot

A VM restarts more often than you might expect — a provider performing host-level maintenance, a kernel security update requiring a reboot, or simply you resizing the VM later, as one of this chapter's exercises suggests. `docker compose up -d` alone does not guarantee the testnet comes back automatically after a reboot; that depends on two separate things being configured correctly.

First, Docker's own daemon needs to be enabled to start automatically on boot — the install script from Section 5 typically does this already, but it is worth confirming rather than assuming:

```bash
sudo systemctl is-enabled docker
# enabled
```

Second, each container needs a **restart policy** telling Docker what to do if the container stops — including "the whole VM rebooted, so every container stopped." Add `restart: unless-stopped` to every service in `docker-compose.yml`:

```yaml
# docker-compose.yml — add this to every service (node1, node2, node3, explorer)
services:
  node1:
    # ... existing config from Chapter 87 ...
    restart: unless-stopped
```

`unless-stopped` restarts a container automatically after a crash or a host reboot, but respects a *deliberate* `docker compose stop` — the distinction matters: you want a crash to self-heal, but you do not want a container you intentionally stopped for maintenance to spring back to life the next time the VM restarts.

Test this without waiting for an actual, unplanned reboot to find out whether it works:

```bash
sudo reboot
# wait a minute or two for the VM to come back up, then reconnect:
ssh deploy@203.0.113.42

docker compose ps
# all four services should already be Up, with no manual `docker
# compose up` needed on your part at all
```

If any service is missing after a reboot, the most common cause is a missing `restart:` line on that specific service — this is exactly the kind of gap that is far better to discover now, deliberately, than to discover the next time your cloud provider reboots the host for maintenance while you are asleep.

---

## 11. Verifying Real Internet Reachability

Everything so far could, in principle, still just be "working on the VM the same way it worked on your laptop." The step that actually proves something new: reach it from a **different** network entirely — your phone's cellular connection (with Wi-Fi turned off), a friend's computer, or an online tool like a web-based port checker.

```bash
# From your phone (cellular data, not the same Wi-Fi as your laptop),
# or from any machine that is not on your home network:
curl http://203.0.113.42:8080/chain/height
# {"height": 0}
```

If that returns a real response, you have just reached a GoChain node from genuinely across the public internet — not localhost, not a Docker network, not your home router's internal addressing, but an actual round trip out to your cloud provider's data center and back. This is the moment GoChain stops being a program that only runs where you are sitting.

Also confirm the port that should **not** be reachable is in fact not:

```bash
curl --max-time 5 http://203.0.113.42:8090
# curl: (28) Connection timed out after 5000 milliseconds
```

A timeout (rather than an immediate refusal) is exactly what `ufw`'s default-deny posture produces — the connection request is silently dropped rather than rejected with an error, which is deliberately slightly harder for a port-scanning attacker to distinguish from "nothing is listening here at all."

---

## 12. What to Do Before You Walk Away

A cloud VM, unlike your laptop, keeps running (and keeps costing money, and keeps being a target for automated internet-wide scanning) whether or not you are paying attention to it. Before moving on to later chapters, a short checklist worth internalizing as a habit for every VM you ever provision:

- **Disable root login over SSH** entirely, now that the `deploy` user works, by editing `/etc/ssh/sshd_config` and setting `PermitRootLogin no`, then `sudo systemctl restart sshd`.
- **Confirm `ufw status`** shows exactly the ports you intend, no more.
- **Note the monthly cost** of the VM size you chose, and set a calendar reminder to revisit it — an idle testnet left running for a year is a real, avoidable expense.
- **Write down the public IP address** somewhere durable; Chapter 90's seed-node address, Chapter 92's `TESTNET_HOST` secret, and Chapter 93's DNS record all point at this exact number.

---

## Summary

- A cloud VM is a rented virtual computer with its own public IP address, reachable from anywhere on the internet — the property that turns a local testnet into a real one.
- Every provider (DigitalOcean, AWS, or otherwise) asks the same handful of questions — image, size, network rules, SSH key — under different names; the resulting Linux server is the same either way.
- SSH key pairs work exactly like Chapter 11's signing keys: a private half you never share, a public half you hand to the provider so it can be installed automatically.
- A fresh VM should immediately get a non-root `deploy` user, patched packages, and Docker installed via `curl -fsSL https://get.docker.com | sh`.
- GoChain's code lives at `/opt/gochain` on the server — the exact path Chapter 92's automated deploy step later assumes.
- `ufw default deny incoming`, followed by explicit `allow` rules for port 22 (SSH), 8080 (API), and 9000 (P2P), is the entire firewall policy this chapter needs — everything else stays closed until a later chapter deliberately opens it.
- A small VM benefits from a swap file as a shock absorber against short memory spikes during image builds, and `restart: unless-stopped` on every Compose service is what lets the testnet come back automatically after a reboot.
- `docker compose up -d --build`, run over SSH from `/opt/gochain`, brings the exact same testnet from Chapter 87 up on a real server instead of your laptop.
- The real test of success is reaching the node from a genuinely different network — your phone's cellular connection, not your home Wi-Fi — and getting back a real response.

---

## Exercises

### Easy

1. Provision a VM on any provider of your choice, generate an SSH key pair if you do not already have one, and log in successfully as `root`. Then create the `deploy` user exactly as shown in Section 4, and confirm you can log back in as `deploy` without a password prompt.

2. Install Docker on your VM using the official convenience script, and run `docker run hello-world` to confirm it works. Then run `docker compose version` and record the exact version string it prints.

3. Configure `ufw` exactly as shown in Section 8, then run `sudo ufw status verbose` and paste the output. Confirm it shows exactly three allowed ports (22, 8080, 9000) and a default-deny posture for everything else.

### Medium

4. Clone your GoChain repository to `/opt/gochain` on the VM and run `docker compose up -d --build`. From a network that is *not* your home Wi-Fi (mobile data, a friend's connection, or an online port-checking tool), confirm `curl http://<your-vm-ip>:8080/chain/height` succeeds and `curl http://<your-vm-ip>:8090` times out.

5. Disable root SSH login (`PermitRootLogin no` in `/etc/ssh/sshd_config`, then restart `sshd`), and confirm, by attempting `ssh root@<your-vm-ip>`, that the connection is now refused while `ssh deploy@<your-vm-ip>` still works.

6. Resize your VM to a smaller tier than the one you started with (most providers support this without recreating the VM), and confirm the running testnet survives the resize and reboot, resuming automatically once the VM comes back up. If it does not restart automatically, investigate and fix that using Docker's restart policies (`restart: unless-stopped` in `docker-compose.yml`).

7. Add the 2GB swap file from Section 6, then run `free -h` before and after to confirm it is active. Deliberately trigger a heavy rebuild (`docker compose build --no-cache`) while watching `free -h` in a second terminal, and note whether swap usage rises during the build.

### Hard

8. Simulate an internet-wide port scan against your own VM using a tool like `nmap` (run from a different machine, with permission — never scan infrastructure you do not own or control) and confirm only ports 22, 8080, and 9000 report as open, with everything else showing filtered or closed. Write 100-150 words on what an attacker scanning your VM would and would not be able to learn from this result.

9. Set up unattended security updates on the VM (`unattended-upgrades` on Ubuntu) so operating-system security patches install automatically without you having to SSH in and run `apt upgrade` manually. Explain, in your own words, the trade-off between automatic patching (which can occasionally restart services unexpectedly) and manual patching (which is safer per-update but easy to forget).

10. Provision a *second* VM from a different cloud provider than your first, install Docker, and bring up the exact same `docker-compose.yml` testnet on it. Configure `node1` on this second VM with `--seed <first-vm-ip>:9000`, and confirm, using `/peers` on both machines, that you now have a real two-server, cross-provider GoChain network — the direct predecessor to Chapter 90's public testnet.
