# Chapter 35: OSINT — Open Source Intelligence Gathering

*OSINT is intelligence from public sources. Before touching a single system, a skilled recon operator can map an organization's employees, technologies, offices, partners, and vulnerabilities — all without leaving a trace.*

---

## OSINT for What?

- **Employee enumeration** → phishing targets, credential stuffing
- **Technology fingerprinting** → what CMS, frameworks, cloud providers
- **Infrastructure mapping** → IP ranges, domains, DNS records
- **Breach data** → leaked credentials, past breaches
- **Social engineering prep** → company structure, relationships
- **Job postings** → reveal technologies in use ("must know Kubernetes, Postgres, AWS")

---

## Passive vs Active Recon

**Passive (OSINT):** No direct contact with target
- Search engines, LinkedIn, GitHub, breach databases
- Target can't detect you (unless they monitor breach DB queries)

**Active recon:** Direct contact
- DNS queries, port scans, HTTP requests
- Appears in target's logs

---

## People and Organization Intelligence

```bash
# LinkedIn — employees, org structure
# Shodan, Censys — internet-exposed infrastructure
# Hunter.io — email addresses: firstname.lastname@company.com
# theHarvester — email + domain discovery
theHarvester -d company.com -b google,linkedin,bing

# Email pattern discovery
curl -s "https://api.hunter.io/v2/domain-search?domain=company.com&api_key=KEY"
# Output: firstname.lastname@company.com → likely pattern

# Verifying emails without sending
curl -s "https://api.hunter.io/v2/email-verifier?email=john@company.com&api_key=KEY"
```

---

## GitHub Reconnaissance

GitHub is a goldmine of secrets accidentally committed:

```bash
# Search GitHub for secrets
# On github.com, search:
# "company.com" password
# "company.com" api_key
# "company.com" secret
# org:companyname password
# org:companyname AWS_ACCESS_KEY

# GitLeaks — scan repos for secrets
gitleaks detect --source=. --report-format json --report-path leaks.json

# TruffleHog — finds secrets in git history
trufflehog github --org=companyname

# git log search
git log -p | grep -i "password\|api_key\|secret\|token"

# Search all commits (even deleted ones)
git log --all -p | grep "password ="
```

### What to Look For

```
AWS keys:     AKIA[0-9A-Z]{16}
GitHub token: ghp_[a-zA-Z0-9]{36}
Slack token:  xox[baprs]-[0-9a-zA-Z]+
JWT:          eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+
Private key:  -----BEGIN RSA PRIVATE KEY-----
Database URL: postgres://user:password@host/dbname
```

---

## Shodan — The Search Engine for Internet-Connected Devices

```bash
# Shodan CLI
pip install shodan
shodan init YOUR_API_KEY

# Search for company's exposed services
shodan search "company.com"
shodan search "ssl:company.com"

# Find specific services
shodan search "hostname:company.com http"
shodan search "org:'Company Name'"

# Find all their IP ranges
shodan search "org:'Company Name'" --fields ip_str,port,hostnames

# Find exposed admin panels
shodan search "http.title:'admin' org:'Company Name'"

# Find outdated software
shodan search "Apache/2.2 org:'Company Name'"
shodan search "openssh 5. org:'Company Name'"
```

### Shodan Dorks

```
ssl.cert.subject.cn:*.company.com    # SSL certs for subdomains
http.favicon.hash:-335242539          # WordPress default favicon
port:8080 http.title:"Tomcat"         # Exposed Tomcat servers
net:203.0.113.0/24 port:22            # SSH on specific subnet
country:IN port:3389                  # Indian RDP servers
product:Redis                         # Exposed Redis (no auth!)
```

---

## Maltego — Visual OSINT (Link Analysis)

Maltego visualizes relationships between entities:
- Domain → subdomains → IPs → organizations → people → emails
- Transforms query online data sources automatically
- Community edition is free

---

## OSINT Framework

```
Target Organization
├── People
│   ├── LinkedIn (employees, org structure)
│   ├── Twitter/X (executives, personalities)
│   ├── Email (theHarvester, Hunter.io)
│   └── Facebook/Instagram (personal info)
├── Technical Infrastructure
│   ├── Domains (subdomain enum, DNS)
│   ├── IPs (Shodan, Censys, ARIN/RIPE)
│   ├── Technologies (BuiltWith, Wappalyzer)
│   └── Certificates (crt.sh)
├── Code
│   ├── GitHub (repos, leaked secrets)
│   └── Pastebin (leaked data)
├── Breach Data
│   ├── Have I Been Pwned
│   └── DeHashed, IntelligenceX
└── Documents
    ├── Google dorks (filetype:pdf)
    └── LinkedIn job postings
```

---

## Google Dorks

```
# Find specific file types on a domain
site:company.com filetype:pdf
site:company.com filetype:xlsx
site:company.com filetype:doc

# Find login pages
site:company.com inurl:login
site:company.com inurl:admin
site:company.com intitle:"admin"

# Find exposed configuration
site:company.com "db_password"
site:company.com "api_key" filetype:txt

# Subdomain enumeration via Google
site:*.company.com -site:www.company.com

# Find exposed panels
site:company.com "PhpMyAdmin"
site:company.com "Kibana"
site:company.com "Grafana"
```

---

## Go: OSINT Domain Profiler

```go
package main

import (
    "fmt"
    "net"
    "strings"
    "time"
)

type DomainProfile struct {
    Domain      string
    IPs         []string
    MXRecords   []string
    NSRecords   []string
    TXTRecords  []string
    Subdomains  []string
    Technologies []string
}

func profileDomain(domain string) *DomainProfile {
    profile := &DomainProfile{Domain: domain}
    
    // A records
    ips, err := net.LookupHost(domain)
    if err == nil {
        profile.IPs = ips
    }
    
    // MX records (email provider)
    mxs, err := net.LookupMX(domain)
    if err == nil {
        for _, mx := range mxs {
            host := strings.TrimSuffix(mx.Host, ".")
            provider := identifyEmailProvider(host)
            profile.MXRecords = append(profile.MXRecords, 
                fmt.Sprintf("%s (pref %d) [%s]", host, mx.Pref, provider))
        }
    }
    
    // TXT records (SPF, DKIM, cloud services)
    txts, err := net.LookupTXT(domain)
    if err == nil {
        for _, txt := range txts {
            profile.TXTRecords = append(profile.TXTRecords, txt)
            // Detect cloud services
            if tech := identifyTechFromTXT(txt); tech != "" {
                profile.Technologies = append(profile.Technologies, tech)
            }
        }
    }
    
    // NS records
    nss, err := net.LookupNS(domain)
    if err == nil {
        for _, ns := range nss {
            profile.NSRecords = append(profile.NSRecords, 
                strings.TrimSuffix(ns.Host, "."))
        }
    }
    
    return profile
}

func identifyEmailProvider(mx string) string {
    providers := map[string]string{
        "google":   "Google Workspace",
        "outlook":  "Microsoft 365",
        "yahoodns": "Yahoo Mail",
        "mxroute":  "MXRoute",
        "zoho":     "Zoho Mail",
        "protonmail": "ProtonMail",
        "mailgun":  "Mailgun",
        "sendgrid": "SendGrid",
    }
    lower := strings.ToLower(mx)
    for keyword, provider := range providers {
        if strings.Contains(lower, keyword) {
            return provider
        }
    }
    return "Unknown"
}

func identifyTechFromTXT(txt string) string {
    techPatterns := map[string]string{
        "v=spf1":              "Email: SPF configured",
        "MS=ms":               "Cloud: Microsoft 365",
        "google-site-verification": "Cloud: Google Workspace",
        "atlassian-domain-verification": "Tools: Atlassian",
        "stripe-verification": "Payment: Stripe",
        "amazonses":           "AWS SES",
        "docusign":            "DocuSign",
        "hubspot":             "HubSpot",
        "segment":             "Segment Analytics",
    }
    for pattern, tech := range techPatterns {
        if strings.Contains(txt, pattern) {
            return tech
        }
    }
    return ""
}

func main() {
    domain := "github.com"
    fmt.Printf("=== OSINT Profile: %s ===\n\n", domain)
    
    profile := profileDomain(domain)
    
    fmt.Printf("IPs: %s\n", strings.Join(profile.IPs, ", "))
    
    fmt.Println("\nMX (Email):")
    for _, mx := range profile.MXRecords {
        fmt.Printf("  %s\n", mx)
    }
    
    fmt.Println("\nNameservers:")
    for _, ns := range profile.NSRecords {
        fmt.Printf("  %s\n", ns)
    }
    
    fmt.Println("\nDetected Technologies:")
    for _, tech := range profile.Technologies {
        fmt.Printf("  %s\n", tech)
    }
    
    _ = time.Second
}
```

---

## Summary

| OSINT Source | What it reveals |
|-------------|----------------|
| LinkedIn | Employees, titles, org structure, tech stack |
| GitHub | Source code, leaked secrets, tech choices |
| Shodan/Censys | Internet-exposed services, versions |
| Hunter.io | Employee emails and patterns |
| crt.sh | All SSL certificates (→ subdomains) |
| Google dorks | Exposed files, admin panels, credentials |
| Have I Been Pwned | Past breaches, leaked credentials |
| Job postings | Technologies in use, planned migrations |

---

## Exercises

1. Profile a company you're authorized to research using theHarvester and Shodan
2. Use GitHub dorks to find any accidentally committed AWS keys or passwords in public repos
3. Use `crt.sh` to enumerate subdomains of a target domain — how many did you find?
4. Build the Go domain profiler above and add technology detection from HTTP headers (Server, X-Powered-By, cookies)
