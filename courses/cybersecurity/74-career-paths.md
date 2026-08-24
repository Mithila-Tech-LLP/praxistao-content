# Chapter 74: Career Paths in Cybersecurity

*The cybersecurity industry has a persistent global talent shortage of 4 million professionals. Opportunities are enormous — the question is which path fits you.*

---

## The Two Sides: Red Team vs Blue Team

Cybersecurity broadly divides into:

**Red Team (Offensive):** You think and act like an attacker.
- Penetration tester
- Bug bounty hunter
- Exploit developer
- Red team operator
- Security researcher

**Blue Team (Defensive):** You detect and respond to attacks.
- SOC Analyst
- Incident Responder
- Threat Hunter
- Security Engineer
- Detection Engineer

**Purple Team:** Bridges both — uses offensive techniques to improve defenses.

Most professionals start in one, then learn the other. The best defenders understand attack; the best attackers understand defense.

---

## Career Paths in Detail

### 1. Penetration Tester (Ethical Hacker)

**What you do:** Companies pay you to try to hack their systems, then report what you found.

**Types:**
- **Network pentesting:** Scanning, exploitation, privilege escalation on internal networks
- **Web app pentesting:** SQL injection, XSS, logic flaws in web applications
- **Mobile pentesting:** Android/iOS app security
- **Red teaming:** Simulating advanced persistent threats (APT) — full attack simulation
- **Physical pentesting:** Social engineering, physical access to buildings

**Tools:** Burp Suite, Metasploit, Nmap, BloodHound, Cobalt Strike, custom scripts

**Salary:** 
- India: ₹6-25 lakh (junior to senior)
- Global/remote: $80,000-$180,000+

**Certifications:** OSCP (Offensive Security Certified Professional) — the most respected, hands-on 24-hour exam. CEH (Certified Ethical Hacker) — more theoretical.

---

### 2. Security Operations Center (SOC) Analyst

**What you do:** Monitor security dashboards (SIEM), investigate alerts, escalate true positives.

**Tiers:**
- **Tier 1:** Alert triage — real or false positive?
- **Tier 2:** Deeper investigation, incident analysis
- **Tier 3:** Threat hunting, advanced investigation

**Tools:** Splunk, Microsoft Sentinel, IBM QRadar, CrowdStrike, SentinelOne

**Salary:**
- India: ₹4-15 lakh
- Global: $50,000-$110,000

**Reality:** Tier 1 SOC is often tedious and high-volume. Most alerts are false positives. The skill is building pattern recognition to find real threats fast.

---

### 3. Incident Responder

**What you do:** When a company is breached, you investigate, contain, eradicate, and recover.

**Skills needed:**
- Memory forensics (Volatility)
- Disk forensics (Autopsy, FTK)
- Log analysis (SIEM, grep/python)
- Network forensics (Wireshark, Zeek)
- Malware analysis (reverse engineering, sandboxing)

**Salary:** $80,000-$150,000 globally (high demand, high stress)

**Reality:** Incident response is extremely exciting and extremely stressful. You're called during crises, often working 24-hour shifts during major breaches. The work is consequential.

---

### 4. Security Engineer / Detection Engineer

**What you do:** Build and maintain security systems. Write detection rules, configure EDR/SIEM, build security tooling.

This is what we've been building throughout this course.

**Skills needed:**
- Programming (Go, Python, Bash)
- Cloud security (AWS/GCP/Azure)
- SIEM administration
- EDR administration
- Detection rule writing (Sigma, YARA, Snort/Suricata rules)

**Salary:**
- India: ₹12-40 lakh
- Global: $100,000-$180,000

---

### 5. Malware Analyst / Reverse Engineer

**What you do:** Analyze malicious software to understand what it does and how to detect/remove it.

**Skills needed:**
- Assembly language (x86/x64)
- Reverse engineering tools (Ghidra — free, IDA Pro — expensive)
- Dynamic analysis (running malware in sandbox, monitoring behavior)
- Static analysis (reading disassembly without running)

**This is rare and very well compensated.**

**Tools:** Ghidra, IDA Pro, x64dbg, Wireshark, Cuckoo Sandbox

**Salary:** $100,000-$200,000+ (specialist role)

---

### 6. Cloud Security Engineer

**What you do:** Secure cloud infrastructure. AWS/GCP/Azure misconfigurations are the #1 source of data breaches today.

**Skills needed:**
- AWS/GCP/Azure fundamentals
- Infrastructure as Code (Terraform)
- Kubernetes security
- Cloud-native security tools (AWS GuardDuty, Security Hub)
- Container security (Docker, K8s RBAC)

**Fastest growing area in security.** Every company is migrating to cloud; almost none have adequate cloud security.

**Salary:**
- India: ₹15-50 lakh
- Global: $120,000-$220,000

---

### 7. Bug Bounty Hunter

**What you do:** Find vulnerabilities in companies' systems, report them, receive payment.

**Platforms:** HackerOne, Bugcrowd, Intigriti

**Reality:**
- Top earners make $500,000+/year
- Most participants earn very little
- Extremely competitive — thousands of hunters targeting the same programs
- Best entry point: learn web application security, start with small private programs

**Success factors:** Depth in specific vulnerability classes (not breadth across all of them), persistence, creativity in chaining bugs.

---

## Certifications

| Certification | Provider | Type | Difficulty | Value |
|--------------|----------|------|------------|-------|
| **OSCP** | Offensive Security | Hands-on, 24hr exam | Hard | Very high |
| **CRTO** | Zero-Point Security | Red team operations | Medium | High |
| **CEH** | EC-Council | Theory + lab | Medium | Medium |
| **CompTIA Security+** | CompTIA | Theory | Easy | Low-medium |
| **CISSP** | ISC² | Management/theory | Hard | High (management) |
| **CISM** | ISACA | Management | Hard | High (management) |
| **AWS Security** | Amazon | Cloud security | Medium | High for cloud |
| **GREM** | GIAC | Malware analysis | Hard | Very high |

**Recommendation for India:**
1. **CEH** — good for getting first job, recognized by Indian companies
2. **OSCP** — premium certification, dramatically better job prospects globally
3. **AWS Security** — if going cloud path

---

## Building Your Skills

### The Right Order

1. **Linux, networking, programming basics** — Chapters 1-8 of this course
2. **Learn Go/Python** for security tools
3. **Set up a home lab** (Chapter 75)
4. **Web application security** — PortSwigger Web Security Academy (free, excellent)
5. **CTF competitions** — picoCTF, Hack The Box, TryHackMe
6. **Specialize** — pick web, network, cloud, or malware
7. **Build a portfolio** — tools you built, CTF writeups, bug bounty reports

### Practice Platforms

**Beginner:**
- **TryHackMe** (tryhackme.com) — guided rooms, great for complete beginners
- **PortSwigger Web Security Academy** — best free web security education
- **PicoCTF** — beginner CTF competitions

**Intermediate:**
- **Hack The Box** (hackthebox.com) — machines to root, great community
- **VulnHub** — downloadable vulnerable VMs

**Advanced:**
- **CTF competitions** — global competitions, practice with a team
- **Bug bounty** — HackerOne, Bugcrowd
- **DVWA, WebGoat, Metasploitable** — locally hosted practice targets

---

## The Indian Cybersecurity Ecosystem

**Companies hiring security professionals in India:**
- **Global tech:** Google, Microsoft, Amazon, Meta (Security Engineering)
- **Indian enterprises:** Infosys, TCS, Wipro, HCL, Tech Mahindra (MSSP services)
- **Indian fintech:** Razorpay, Paytm, Zerodha, CRED (internal security)
- **Security vendors:** Palo Alto, CrowdStrike, SentinelOne, Symantec (India offices)
- **MSSPs:** Multiple Managed Security Service Providers run large SOC operations in India
- **Indian security companies:** Sequretek, Lucideus (Safe Security), CyberArk India

**Indian certifications/training:**
- **CERT-In** (Indian Computer Emergency Response Team) — government cybersecurity body
- **DSCI** (Data Security Council of India) — industry body
- **C|EH** training available from many Indian providers

**Salary ranges in India (2024):**
- Junior SOC Analyst: ₹4-8 lakh
- Mid-level Pentester: ₹10-20 lakh
- Senior Security Engineer: ₹20-45 lakh
- Security Architect: ₹35-70 lakh
- CISO (large company): ₹70 lakh - ₹2 crore

---

## Building a Portfolio

**What stands out to employers:**

1. **GitHub with real security tools** — this course gave you a port scanner, SQL injection tester, GoShield EDR, and more. Document them well.

2. **CTF writeups** — solve a Hack The Box machine, write a detailed blog post explaining your methodology. Employers read these.

3. **Bug bounty finds** — even a single valid finding (P4/P3) on a public bug bounty program shows you can do real security work.

4. **Open source contributions** — contribute to security tools like Nuclei, subfinder, or similar.

5. **Research blog** — write about a vulnerability you studied, a technique you learned, or a security concept you understood deeply.

---

## Key Lessons

1. **The field is wide — specialize early.** Jack of all security trades is a master of none. Pick web, cloud, malware, or red team and go deep.

2. **Certifications matter for the first job, skills matter for the next one.** CEH gets your resume past HR; OSCP signals real skill.

3. **Practice matters more than coursework.** A student who has rooted 50 Hack The Box machines is more hireable than one who read 10 security books.

4. **Build in public.** Every tool you publish, every CTF writeup you write, every vulnerability you report builds your reputation.

5. **The hacker mindset is a lifestyle.** The best security professionals are permanently curious — they want to understand how everything works and where it breaks.

---

## Summary

| Path | Entry point | Best for | Salary range (India) |
|------|------------|----------|---------------------|
| Penetration tester | Web/network hacking labs | Attack-minded, problem solver | ₹8-25L |
| SOC Analyst | TryHackMe, SIEM training | Alert investigation, patterns | ₹4-15L |
| Security Engineer | Programming + cloud | Builders, coders | ₹12-40L |
| Incident Responder | Forensics training | High-pressure, crisis | ₹15-35L |
| Bug Bounty | PortSwigger Academy | Independent, creative | Variable |
| Cloud Security | AWS/GCP training | Cloud-native era | ₹15-50L |

---

## Exercises

1. Create a Hack The Box account and root your first machine. Write a writeup.
2. Complete PortSwigger's SQL injection lab series (free)
3. Sign up for a bug bounty program on HackerOne (many have free/public programs)
4. Research OSCP — what topics does it cover? Is it the right cert for your path?
5. Build a GitHub portfolio: push all the tools from this course with READMEs
