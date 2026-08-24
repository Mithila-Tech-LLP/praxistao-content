# Chapter 19: DNS Reconnaissance — Mapping a Target Through Their DNS

*DNS is the internet's phone book — and it's often poorly secured. A thorough DNS recon can reveal all subdomains, mail servers, third-party services, and internal infrastructure before you touch a single target system.*

---

## Why DNS Recon is the First Step

DNS leaks a target's entire infrastructure if not properly locked down:
- Subdomains = additional attack surfaces
- MX records = email provider (useful for phishing)
- TXT records = SPF/DKIM policies, cloud services, domain verification
- Zone transfer = dump ALL records at once
- Reverse DNS = find hostnames for IPs

A thorough recon often reveals: dev, staging, admin, api, vpn, mail, ftp subdomains — each potentially vulnerable.

---

## Basic DNS Queries

```bash
# A record (IPv4)
dig example.com A
dig +short example.com

# All record types
dig example.com ANY

# Mail servers
dig example.com MX

# Name servers
dig example.com NS

# TXT records (SPF, DKIM, domain verification)
dig example.com TXT

# Reverse DNS (PTR)
dig -x 93.184.216.34

# Specific DNS server
dig @8.8.8.8 example.com

# Using host command
host example.com
host -t mx example.com
```

---

## Zone Transfer — The Jackpot

If DNS is misconfigured, you can download ALL their DNS records.

```bash
# Find nameservers first
dig NS example.com

# Try zone transfer from each nameserver
dig axfr @ns1.example.com example.com
dig axfr @ns2.example.com example.com

# Alternative tool
host -l example.com ns1.example.com

# Real example that always works for practice:
dig axfr @nsztm1.digi.ninja zonetransfer.me
```

A successful zone transfer dumps every subdomain, internal IP, mail server, and more.

---

## Subdomain Enumeration

Zone transfers rarely work (most orgs block them). Use brute force or passive methods.

```bash
# Brute force with wordlist
gobuster dns -d example.com -w /usr/share/wordlists/dns/subdomains-top1million-5000.txt
ffuf -w /usr/share/wordlists/dns/subdomains-top1million-5000.txt -u http://FUZZ.example.com

# Subfinder (passive — uses OSINT sources)
subfinder -d example.com

# Amass (comprehensive)
amass enum -d example.com
amass enum -active -d example.com  # active probing

# Certificate transparency logs (crt.sh)
curl -s "https://crt.sh/?q=%.example.com&output=json" | \
    python3 -c "import sys,json; [print(x['name_value']) for x in json.load(sys.stdin)]" | \
    sort -u
```

---

## Go: DNS Brute Forcer

```go
package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "os"
    "strings"
    "sync"
    "time"
)

type DNSResult struct {
    Subdomain string
    IPs       []string
}

func trySubdomain(subdomain, domain string, timeout time.Duration) *DNSResult {
    fqdn := subdomain + "." + domain
    
    resolver := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx interface{}, network, address string) (net.Conn, error) {
            d := net.Dialer{Timeout: timeout}
            return d.DialContext(ctx.(interface{ Deadline() (interface{}, bool) }), "udp", "8.8.8.8:53")
        },
    }
    
    // Simple approach: just use net.LookupHost
    ips, err := net.LookupHost(fqdn)
    if err != nil {
        return nil
    }
    
    _ = resolver
    return &DNSResult{
        Subdomain: fqdn,
        IPs:       ips,
    }
}

func bruteForce(domain, wordlistPath string, workers int) {
    f, err := os.Open(wordlistPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Cannot open wordlist:", err)
        return
    }
    defer f.Close()
    
    jobs    := make(chan string, workers)
    results := make(chan *DNSResult, workers)
    var wg sync.WaitGroup
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for sub := range jobs {
                if r := trySubdomain(sub, domain, 2*time.Second); r != nil {
                    results <- r
                }
            }
        }()
    }
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    go func() {
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            sub := strings.TrimSpace(scanner.Text())
            if sub != "" && !strings.HasPrefix(sub, "#") {
                jobs <- sub
            }
        }
        close(jobs)
    }()
    
    found := 0
    for r := range results {
        fmt.Printf("[+] %s → %s\n", r.Subdomain, strings.Join(r.IPs, ", "))
        found++
    }
    fmt.Printf("\nFound %d subdomains\n", found)
}

func main() {
    domain   := flag.String("d", "", "Target domain")
    wordlist := flag.String("w", "subdomains.txt", "Wordlist path")
    workers  := flag.Int("t", 50, "Concurrent workers")
    flag.Parse()
    
    if *domain == "" {
        fmt.Println("Usage: dnsbf -d example.com -w wordlist.txt")
        return
    }
    
    fmt.Printf("Brute-forcing subdomains of %s\n", *domain)
    bruteForce(*domain, *wordlist, *workers)
}
```

---

## Passive DNS — Finding Subdomains Without Querying

```bash
# Certificate Transparency (no direct queries to target)
curl -s "https://crt.sh/?q=%.example.com&output=json" | \
    python3 -m json.tool | grep "name_value" | sort -u

# Wayback Machine CDX API
curl -s "http://web.archive.org/cdx/search/cdx?url=*.example.com&output=json&fl=original&collapse=urlkey" | \
    python3 -c "import sys,json; [print(x[0]) for x in json.load(sys.stdin)]"

# VirusTotal (need API key)
curl -s "https://www.virustotal.com/vtapi/v2/domain/report?apikey=KEY&domain=example.com"

# SecurityTrails, Shodan, Censys (commercial APIs)
```

---

## DNS for Detecting Attack Infrastructure

Defenders also use DNS to detect malware:

```bash
# Fast-flux domains (changing IPs every seconds = botnet C2)
# Run dig repeatedly:
for i in $(seq 1 5); do dig +short evil-domain.com; sleep 2; done
# If IPs change every query = fast-flux!

# Domain Generation Algorithms (DGA) — malware creates random domains
# Look for domains with high entropy (random-looking)
echo "xkj38dkas9f.com" | python3 -c "
import sys, math, collections
s = sys.stdin.read().strip()
freq = collections.Counter(s)
entropy = -sum(f/len(s)*math.log2(f/len(s)) for f in freq.values())
print(f'Entropy: {entropy:.2f} (>3.5 = likely random/DGA)')"
```

---

## Summary

| Technique | Tool | What it finds |
|-----------|------|--------------|
| Zone transfer | `dig axfr` | All DNS records at once |
| Subdomain brute force | `gobuster dns`, custom Go | Subdomains from wordlist |
| Passive enum | `subfinder`, `crt.sh` | Subdomains without querying target |
| Reverse DNS | `dig -x` | Hostname for IP |
| MX lookup | `dig MX` | Email provider |
| TXT lookup | `dig TXT` | Cloud services, SPF policy |

---

## Exercises

1. Perform a zone transfer against `zonetransfer.me` (intentionally vulnerable). How many subdomains does it have?
2. Use `subfinder` and `crt.sh` against a company you're authorized to test. Compare results.
3. Extend the Go DNS brute forcer to also check for wildcard DNS (if `*.domain.com` resolves, all subdomains will — filter them out)
4. Write a Go tool that checks a domain's SPF record and identifies if it's configured correctly to prevent spoofing
