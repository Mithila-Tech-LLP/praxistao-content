# Chapter 73: Certifications and Career Path

*Cybersecurity certifications validate your skills to employers. Choosing the right ones and preparing efficiently can transform your career trajectory.*

---

## The Certification Landscape

```
Entry Level
├── CompTIA Security+        — foundational, vendor-neutral
├── CompTIA Network+         — networking (prereq for Security+)
└── Google Cybersecurity     — beginner, Coursera-based

Intermediate
├── eJPT (INE)               — entry-level pentesting
├── CEH (EC-Council)         — ethical hacker certification
├── GIAC GSEC               — GIAC security essentials
└── CompTIA CySA+            — cybersecurity analyst

Advanced Pentesting
├── OSCP (Offensive Security) — gold standard, hands-on
├── CRTO (Red Team Ops)      — Cobalt Strike, red teaming
├── CRTE                     — Active Directory red team
└── OSED, OSEP, OSWE         — specialized OSCP family

Defensive
├── GCIA (GIAC intrusion analyst)
├── GCIH (GIAC incident handler)
└── GCFE (GIAC forensics examiner)

Cloud
├── AWS Security Specialty   — AWS security
├── GCP Professional Security Engineer
└── Microsoft SC-200         — Azure Defender

Management
├── CISSP                    — senior security management
├── CISM                     — information security manager
└── CISA                     — IS auditor
```

---

## OSCP — The Pentester's Standard

OSCP is widely regarded as the best practical certification:

```
Format:
- 23h 45min exam (24h with setup time)
- Hack 3 standalone machines + active directory set
- 70 points to pass (out of 100)
- Real machines (no multiple choice)
- Report due 24h after exam

Prerequisites:
- Complete PWK (Penetration Testing with Kali Linux) course
- Practice on HackTheBox / TryHackMe / VulnHub
- Master privilege escalation (Windows + Linux)
- Understand Active Directory attacks

Preparation Path (6-12 months):
Month 1-2: Linux basics, networking, scripting
Month 3-4: TryHackMe rooms (Jr Penetration Tester path)
Month 5-6: HackTheBox (Easy boxes consistently)
Month 7-8: PWK course content
Month 9-10: PWK lab machines (do all 60+)
Month 11+: TCM PNPT, practice exams, buffer overflow mastery
```

---

## Practical Career Paths

### Path 1: Penetration Tester / Red Teamer

```
Start:
1. Learn fundamentals (networking, Linux, scripting)
2. TryHackMe beginner rooms
3. CompTIA Security+ (if needed for job market)
4. HackTheBox Active (Easy → Medium consistently)
5. OSCP

Progress:
5. Mid-level: Junior pentester at consultancy
6. Specialize: web apps, AD, cloud, mobile
7. CRTO (red team) or OSEP (evasion) for senior roles
8. Start consulting independently or lead engagements
```

### Path 2: Defensive Security / SOC Analyst

```
Start:
1. Networking and OS fundamentals
2. Security+ or equivalent
3. Blue Team Labs Online / CyberDefenders
4. Build home SIEM (ELK Stack)
5. TryHackMe (SOC Level 1 path)

Progress:
5. SOC Analyst Tier 1 job
6. GIAC GCIH (incident handler)
7. Threat hunting and detection engineering
8. GIAC GREM (reverse engineering malware)
9. Senior threat hunter / detection engineer
```

### Path 3: Bug Bounty Hunter

```
Start:
1. Web security fundamentals (PortSwigger Web Academy)
2. Burp Suite proficiency
3. Complete all PortSwigger labs
4. Bug bounty on HackerOne (low-scope targets first)
5. Learn one platform deeply (WordPress, Shopify, etc.)

Progress:
5. Consistent small bounties ($100-500)
6. Specialize in cloud, APIs, or mobile
7. Build reputation on leaderboard
8. Private programs (invite-only, higher payouts)
9. Full-time independent researcher
```

---

## Study Resources

```
Free:
├── TryHackMe                — guided rooms, beginner-friendly
├── HackTheBox               — harder, real pentesting practice
├── PortSwigger Web Security Academy — best web security training
├── TCM Security (YouTube)   — free pentesting courses
├── Professor Messer          — CompTIA study videos
├── SANS Cyber Aces          — free fundamentals
└── PicoCTF                  — beginner CTFs

Paid (worth it):
├── TCM Security PNPT        — practical, OSCP prep
├── INE/eLearnSecurity       — eJPT, eCPPT paths
├── Hack The Box Academy     — structured learning paths
└── SANS courses             — expensive but premium

YouTube Channels:
├── IppSec                   — HackTheBox walkthroughs
├── John Hammond             — CTF and pentesting
├── TCM Security             — practical techniques
├── LiveOverflow             — binary exploitation, CTF
└── 0xdf                     — writeups with depth
```

---

## Go: Certification Exam Tracker

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type Certification struct {
    Name        string    `json:"name"`
    Provider    string    `json:"provider"`
    Status      string    `json:"status"`     // planned, studying, passed, failed
    TargetDate  time.Time `json:"target_date"`
    PassedDate  time.Time `json:"passed_date,omitempty"`
    Score       int       `json:"score,omitempty"`
    Notes       string    `json:"notes"`
}

type StudyPlan struct {
    Certifications []Certification `json:"certifications"`
    UpdatedAt      time.Time       `json:"updated_at"`
}

func loadPlan(path string) *StudyPlan {
    data, err := os.ReadFile(path)
    if err != nil {
        return &StudyPlan{}
    }
    var plan StudyPlan
    json.Unmarshal(data, &plan)
    return &plan
}

func (p *StudyPlan) Save(path string) error {
    p.UpdatedAt = time.Now()
    data, err := json.MarshalIndent(p, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func (p *StudyPlan) Summary() {
    fmt.Printf("=== CERTIFICATION STUDY PLAN ===\n\n")
    
    statusGroups := map[string][]Certification{}
    for _, c := range p.Certifications {
        statusGroups[c.Status] = append(statusGroups[c.Status], c)
    }
    
    order := []string{"studying", "planned", "passed", "failed"}
    for _, status := range order {
        certs := statusGroups[status]
        if len(certs) == 0 {
            continue
        }
        
        fmt.Printf("[%s]\n", status)
        for _, c := range certs {
            if status == "studying" {
                daysLeft := int(time.Until(c.TargetDate).Hours() / 24)
                fmt.Printf("  %-20s (%s) — %d days to exam\n", c.Name, c.Provider, daysLeft)
            } else if status == "passed" {
                fmt.Printf("  %-20s (%s) — passed %s\n", c.Name, c.Provider,
                    c.PassedDate.Format("2006-01-02"))
            } else {
                fmt.Printf("  %-20s (%s)\n", c.Name, c.Provider)
            }
        }
        fmt.Println()
    }
}

func main() {
    plan := &StudyPlan{
        Certifications: []Certification{
            {
                Name:       "CompTIA Security+",
                Provider:   "CompTIA",
                Status:     "studying",
                TargetDate: time.Now().AddDate(0, 2, 0),
                Notes:      "Professor Messer videos + Darril Gibson book",
            },
            {
                Name:       "OSCP",
                Provider:   "Offensive Security",
                Status:     "planned",
                TargetDate: time.Now().AddDate(0, 8, 0),
                Notes:      "After completing HTB consistently",
            },
        },
    }
    
    plan.Summary()
    plan.Save("study-plan.json")
}
```

---

## Realistic Timeline

```
Complete Beginner → OSCP Level
Year 1:
  Q1: Linux + Networking fundamentals
      TryHackMe: Pre-Security path
  Q2: Web security basics
      PortSwigger: All Apprentice labs
  Q3: Pentesting basics
      TryHackMe: Jr Pentesting path
  Q4: HTB Easy boxes consistently
      CompTIA Security+ (if needed)

Year 2:
  Q1-Q2: OSCP course + labs
  Q3: More practice, review weak areas
  Q4: OSCP exam

OSCP → Senior Level (Years 3-5):
  Specialize in one area
  Active Directory expertise (CRTO/CRTE)
  Bug bounty side income
  Lead engagements, mentor juniors
```

---

## Summary

| Certification | Best For | Difficulty | Cost |
|--------------|----------|-----------|------|
| Security+ | First job, DoD requirement | Beginner | $370 |
| eJPT | Pentesting entry | Beginner | $200/year |
| PNPT (TCM) | Practical pentest | Intermediate | $400 |
| OSCP | Gold standard pentest | Advanced | $1,499 |
| CRTO | Red teaming | Advanced | £399 |
| CISSP | Management | Expert | $699 |
| GCIH | Incident response | Intermediate | $949 |

---

## Exercises

1. Make your own study plan using the Go tracker above
2. Sign up for TryHackMe and complete the "Jr Penetration Tester" path
3. Create a profile on HackTheBox and solve your first "Easy" box (start with retired machines with walkthroughs)
4. Read 3 OSCP exam reports shared publicly — understand what the exam structure requires
