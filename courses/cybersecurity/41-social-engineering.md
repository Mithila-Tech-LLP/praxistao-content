# Chapter 41: Social Engineering — Hacking Humans

*Technical defenses only protect against technical attacks. Social engineering bypasses all technology by targeting the weakest link: people. Understanding these attacks is essential for both executing and defending against them.*

---

## What is Social Engineering?

Social engineering manipulates people into:
- Revealing credentials or sensitive information
- Installing malware
- Granting unauthorized access
- Transferring money (BEC attacks)

The attack surface is human psychology, not software.

---

## Core Principles of Manipulation

Robert Cialdini's influence principles exploited in attacks:

| Principle | How attackers use it |
|-----------|---------------------|
| **Authority** | Impersonate CEO, IT department, auditor |
| **Urgency** | "Your account will be locked in 24 hours" |
| **Scarcity** | "Only 2 spots left for the training" |
| **Social proof** | "All other employees have already updated..." |
| **Liking** | Build rapport, seem friendly/helpful |
| **Reciprocity** | Offer help first, ask for favor later |
| **Commitment** | Get small yes first, then larger requests |

---

## Phishing

The most common initial access vector:

```
Types of phishing:
├── Phishing          — mass, generic (same email to millions)
├── Spear Phishing    — targeted (specific person/organization)
├── Whaling           — targeting executives (CEO, CFO)
├── Vishing           — voice phishing (phone calls)
├── Smishing          — SMS phishing
└── Business Email Compromise (BEC) — impersonate CEO for wire transfer
```

### Phishing Email Anatomy

```
From: IT-Security <it.security@company-corp.com>   ← lookalike domain
Subject: URGENT: Password Expiration Notice

Dear [Employee Name],                               ← personalized (spear)

Your network password expires in 24 hours.         ← urgency
Please update immediately at the link below:

https://company-secure-portal.com/update           ← lookalike domain
                                                    ← actually attacker's site

Failure to update will result in account lockout.  ← threat

Best regards,
IT Security Team                                    ← authority
```

### Technical Phishing Setup

```bash
# GoPhish — open-source phishing framework
# Set up phishing campaigns with tracking

# Create lookalike domain
# company.com → c0mpany.com, compàny.com (IDN homograph), company-secure.com

# Clone legitimate login page
httrack https://company.com/login -O cloned/

# Gophish configuration
# 1. SMTP relay server
# 2. Cloned landing page
# 3. Target list
# 4. Email template

# Evilginx2 — proxy phishing (captures session cookies, bypasses MFA)
# Sits between user and real site, captures all traffic including session tokens
```

---

## Pretexting

Creating a false scenario (pretext) to manipulate the target:

```
Scenario 1: IT Support
"Hi, this is Mike from the IT helpdesk. We're seeing unusual activity from your 
account and need to verify your identity. Can you confirm your username and 
current password so I can check the system?"

Scenario 2: New Employee
Walk into a building dressed as an IT contractor carrying equipment.
"Hi, I'm installing new hardware for the network upgrade. Where's the server room?"

Scenario 3: Vendor
"I'm calling from Microsoft. Our monitoring detected a problem with your server.
I need to walk you through a remote fix..."

Scenario 4: CEO Fraud (BEC)
Email appearing from CEO to CFO:
"I need you to process an urgent wire transfer of $85,000 to our new 
acquisition partner. This is confidential. Do it now before end of day."
```

---

## Vishing (Voice Phishing)

```
Phone call attack flow:
1. Research target (LinkedIn, company website)
2. Build pretext (IT support, vendor, executive)
3. Call victim — use spoofed caller ID (looks like internal number)
4. Establish rapport (reference real info found via OSINT)
5. Create urgency
6. Extract credentials or install "security software" (RAT)

Common pretexts:
- IT help desk password reset
- Benefits enrollment deadline
- Payroll/direct deposit update
- "Your email is about to be shut down"
```

---

## Physical Security Testing

```
Physical attack vectors:
├── Tailgating        — follow an employee through secure door
├── Impersonation     — delivery person, IT contractor, visitor
├── Dumpster diving   — find sensitive docs in trash
├── USB drops         — leave malicious USB in parking lot
└── Shoulder surfing  — watch someone type password

USB Drop Attack:
1. Create malicious USB (auto-run payload via HID attack)
2. Label it "Q4 Salary Data.xlsx" or "Layoff List 2024"
3. Drop in company parking lot
4. 45% of people plug in random USBs (research shows)

Rubber Ducky / O.MG Cable / Hak5 tools:
- Appear as keyboard (HID device) — types commands at 1000 WPM
- Bypasses all endpoint security (it's "just a keyboard")
```

---

## Go: Phishing Simulation Tool (Authorized Testing)

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type PhishClick struct {
    Target    string
    IP        string
    UserAgent string
    Time      time.Time
    Opened    bool
    Clicked   bool
}

var clicks []PhishClick

func trackHandler(w http.ResponseWriter, r *http.Request) {
    target := r.URL.Query().Get("t")
    
    click := PhishClick{
        Target:    target,
        IP:        r.RemoteAddr,
        UserAgent: r.UserAgent(),
        Time:      time.Now(),
        Clicked:   true,
    }
    clicks = append(clicks, click)
    
    log.Printf("[PHISH CLICK] Target: %s | IP: %s | UA: %s",
        target, r.RemoteAddr, r.UserAgent())
    
    // Redirect to real site (so it's not obvious)
    http.Redirect(w, r, "https://company.com/login", http.StatusFound)
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "total_targets": 100,
        "clicked":       len(clicks),
        "click_rate":    fmt.Sprintf("%.1f%%", float64(len(clicks))/100*100),
        "clicks":        clicks,
    })
}

func main() {
    http.HandleFunc("/t", trackHandler)
    http.HandleFunc("/report", reportHandler)
    
    fmt.Println("Phishing tracker running on :8080")
    fmt.Println("Track URL: http://attacker.com/t?t=TARGET_ID")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Defending Against Social Engineering

```
Technical Controls:
- Email filtering (SPF, DKIM, DMARC)
- Anti-phishing email gateway
- MFA everywhere (preferably hardware keys / FIDO2)
- URL filtering / sandboxing
- Disable USB auto-run
- Visitor badging with escort policies

Human Controls:
- Security awareness training (regular, not annual)
- Simulated phishing campaigns (measure click rates)
- Verify phone requests via callback to known number
- "If in doubt, don't" culture (no punishment for reporting)
- Clear procedures for wire transfers, password resets

Process Controls:
- Dual approval for large transfers
- Out-of-band verification for sensitive requests
- Clear escalation path for suspicious requests
```

---

## Summary

| Attack | Target | Defense |
|--------|--------|---------|
| Phishing | Mass email credential theft | Anti-phishing gateway, DMARC, MFA |
| Spear phishing | Specific employee | Security awareness training |
| Vishing | Phone — credentials/access | Callback verification policy |
| BEC | CFO wire transfer | Dual approval, out-of-band confirm |
| USB drops | Physical access | Disable USB, security culture |
| Tailgating | Physical building access | Access control, anti-tailgate doors |

---

## Exercises

1. Run a GoPhish campaign against yourself or your team (with permission) — what click rate do you get?
2. Research the Twitter 2020 hack — it was a social engineering attack. What exactly happened?
3. Test your organization's physical security: can you tailgate into a restricted area?
4. Study DMARC — implement it for a domain you control and verify it blocks spoofed emails
