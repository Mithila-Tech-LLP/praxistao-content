# Chapter 17: Online Safety — Staying Safe on the Internet

> **"The internet is the most powerful tool humanity has ever created. It's also the largest crime scene in history. You don't have to be afraid — you just have to be smart. A few good habits protect you from 99% of threats."**

---

## Table of Contents

1. [The Real Threats Online](#1-the-real-threats-online)
2. [Passwords — Your First Line of Defense](#2-passwords--your-first-line-of-defense)
3. [Two-Factor Authentication (2FA)](#3-two-factor-authentication-2fa)
4. [Phishing — Fake Websites and Fake Emails](#4-phishing--fake-websites-and-fake-emails)
5. [Malware and Viruses](#5-malware-and-viruses)
6. [Privacy — Who's Watching You?](#6-privacy--whos-watching-you)
7. [Safe Online Habits — The Short List](#7-safe-online-habits--the-short-list)
8. [For Parents — Keeping Kids Safe](#8-for-parents--keeping-kids-safe)
9. [Summary](#summary)

---

## 1. The Real Threats Online

```
Realistic threats (most common):
  
  Account takeover:
    Someone guesses or steals your password → logs into your email/bank/social
    Prevention: strong passwords + 2FA
    
  Phishing:
    Fake email/website tricks you into entering your password
    Prevention: check URLs, be suspicious of urgency
    
  Malware:
    Malicious software on your computer / phone
    Usually from downloading pirated software or clicking bad links
    Prevention: don't download random stuff, keep OS updated
    
  Data breaches:
    A company you use (LinkedIn, Facebook, Target, Adobe) gets hacked
    Your email + password are now in criminal databases
    Prevention: use unique passwords for every site (password manager)
    
Less common (but scary-sounding) threats:
  
  Man-in-the-middle:
    Someone intercepts your traffic on public Wi-Fi
    Prevention: use HTTPS (🔒), use VPN on public Wi-Fi
    
  Ransomware:
    Malware encrypts all your files and demands Bitcoin to unlock
    Mainly targets businesses and organizations
    Prevention: keep data backed up, keep software updated
```

---

## 2. Passwords — Your First Line of Defense

```
Bad password habits (extremely common):
  
  Using simple passwords: "password", "123456", "qwerty", "iloveyou"
  These are cracked in seconds.
  
  Using the same password everywhere:
  If ANY website gets hacked → your password is now known.
  Hackers try stolen passwords on Gmail, banks, everything.
  This is called "credential stuffing."
  
  Using personal info: your birthday, name, pet's name, city.
  A determined attacker who knows you will try these.

What makes a good password:
  
  Long: at least 12 characters (length matters most)
  Random: no real words, no patterns
  Unique: different for every website
  
  Good examples:
    Xq7#mPk2@Yw9
    purple!horse!battery!staple
    (a passphrase — 4 random words — is very strong and memorable)
  
  Bad examples:
    December2024
    Mike1985!
    P@ssw0rd
    
The solution: a PASSWORD MANAGER
  Apps like 1Password, Bitwarden, Apple Keychain, Google Passwords
  Generate and store random unique passwords for every site
  You only need to remember ONE master password
  They fill passwords automatically when you log in
  
  THIS IS THE SINGLE BEST THING YOU CAN DO FOR ONLINE SECURITY.
```

---

## 3. Two-Factor Authentication (2FA)

Passwords can be stolen. 2FA means even if someone has your password, they still can't log in without a second factor.

```
How 2FA works:

  You log in with your password (factor 1: something you KNOW)
  
  The service asks for a second verification (factor 2):
    SMS code: "Your verification code is 847291"
    App code: Authenticator app shows a 6-digit code that changes every 30 seconds
    Physical key: a USB security key (YubiKey) you plug in
    Biometric: your fingerprint or face
    Email code
  
  Even if a hacker has your password, they don't have your phone.
  So they can't complete login.

Strength of different 2FA methods (weakest to strongest):
  
  SMS code            → decent (can be SIM-swapped in sophisticated attacks)
  Email code          → decent (if email itself is secured)
  Authenticator app   → good (Google Authenticator, Authy, 1Password)
  Physical key        → best (virtually unphishable)

Enable 2FA on at minimum:
  ✅ Your email (most important — used to reset all other passwords)
  ✅ Your bank and financial accounts
  ✅ Your Apple ID / Google Account
  ✅ Social media accounts
```

---

## 4. Phishing — Fake Websites and Fake Emails

Phishing is when criminals pretend to be a trusted company to steal your login.

```
Classic phishing email:
  
  FROM: security@app1e.com  ← note the "1" instead of "i"
  SUBJECT: Your Apple ID has been locked
  
  "Dear Customer,
  Your Apple ID has been locked due to suspicious activity.
  You must verify your account within 24 hours or it will be permanently suspended.
  
  Click here to verify: http://apple-security-verify.xyz/login"
  
  ↑ This link goes to a fake Apple website that looks identical to real Apple.
  ↑ You enter your password → criminals capture it.

Red flags:
  Urgency: "Act now!" "Within 24 hours!" "Account suspended!"
  Suspicious sender address: look carefully at the full email address
  Suspicious link: hover over it — the real URL is often obviously wrong
  Unexpected: you didn't initiate this contact
  Bad spelling/grammar (though AI is making fake emails more convincing)
  
Golden rules:
  Never click links in emails about your account.
  Instead, type the website address directly in your browser.
  
  Your bank will NEVER ask for your password by email.
  Apple will NEVER email you asking to verify your Apple ID.
  
How to check a URL:
  Real Apple:   apple.com
  Fake Apple:   apple-support-verify.com / app1e.com / secure-apple.net
  
  The REAL domain is the part just before .com/.org/.net
  In "secure.apple-login.net" → domain is "apple-login.net" — NOT apple!
  In "apple.com/support" → domain is "apple.com" — REAL ✓
```

---

## 5. Malware and Viruses

**Malware** = malicious software. Viruses, trojans, ransomware, spyware — all are types of malware.

```
How malware gets onto your computer:
  
  Downloading pirated software:
    "Free Photoshop" → contains malware
    "Free Netflix account generator" → malware
    
  Malicious email attachments:
    "Invoice_March.pdf.exe" → executable disguised as PDF
    
  Malicious websites:
    Some websites try to run code to exploit browser vulnerabilities
    Keeping your browser updated prevents most of these
    
  USB drives:
    Malware can spread via physical USB drives
    Never plug in a USB drive you found or don't trust
    
How to protect yourself:
  
  Keep OS and apps updated:
    Most malware exploits known vulnerabilities.
    Updates patch those vulnerabilities.
    Enable automatic updates.
    
  Windows: Windows Defender is now very good. Free, built-in.
  Mac: macOS has Gatekeeper + XProtect. Also good.
  Don't download pirated software.
  Don't click "You have a virus! Click here to remove it!" pop-ups (these ARE the virus)
  
Signs you might have malware:
  Computer suddenly very slow
  Strange pop-ups
  Browser redirects to weird sites
  Unknown programs in your taskbar
  Files encrypted with ransom note
```

---

## 6. Privacy — Who's Watching You?

```
Who knows what you do online:
  
  Your ISP (Internet Service Provider):
    Can see every website you visit (even with incognito mode)
    Required by law to log this in some countries
    Solution: VPN encrypts your traffic so ISP only sees "VPN server"
    
  Websites you visit:
    Your IP address (approximate location)
    What you clicked, how long you stayed
    Your browser type and version
    
  Google/Meta (if logged in):
    Search history, location, videos watched, ads clicked
    Used for targeted advertising
    
  Apps on your phone:
    Many apps request access to location, contacts, camera, microphone
    Some use these legitimately, others collect data unnecessarily
    
Protecting your privacy:
  
  VPN (Virtual Private Network):
    Encrypts your internet traffic
    Hides your real IP address from websites
    Good for: public Wi-Fi, avoiding ISP tracking
    Not a magic shield: VPN provider now sees your traffic instead
    
  Browser: Firefox or Brave block more trackers than Chrome
  Search: DuckDuckGo doesn't track searches
  Email: ProtonMail is encrypted
  
  Phone settings:
    Go to Settings → Privacy
    Review what permissions each app has (location, camera, microphone)
    Revoke anything that seems unnecessary
```

---

## 7. Safe Online Habits — The Short List

```
✅ DO:
  Use a password manager with unique strong passwords
  Enable 2FA on important accounts (especially email)
  Keep OS and apps updated
  Look at the full sender email and URL before clicking
  Use HTTPS websites (look for 🔒)
  Back up your data (photos, documents) to an external drive or cloud
  Use your phone's biometric login (Face ID, fingerprint) for apps
  
❌ DON'T:
  Use the same password on multiple sites
  Click links in unexpected emails
  Download pirated software
  Share your password with anyone (including "tech support" who calls you)
  Click "You won a prize! Click here!" pop-ups
  Use public Wi-Fi for banking without a VPN
  Post personal details publicly (birthday, phone number, home address)
  
🤔 BE SKEPTICAL OF:
  Anyone calling claiming to be Microsoft/Apple/Your Bank needing access
  (Real tech companies never cold-call you asking for access)
  Emails asking you to verify account details urgently
  "Free" software that seems too good to be true
  Offers that require your credit card for a "free trial"
```

---

## 8. For Parents — Keeping Kids Safe

```
Age-appropriate use:
  Under 13: No social media (most platforms require 13+ by law, TikTok 13+, etc.)
  13-16: Supervised use, privacy settings locked down
  
Tools:
  Screen time limits: iPhone Screen Time, Android Digital Wellbeing, Google Family Link
  Content filters: parental controls on your router and devices
  Location sharing: Find My (Apple), Google Family Sharing
  
Conversations > Restrictions:
  Teach kids HOW to be safe rather than only blocking access.
  They will find ways around restrictions — understanding why matters more.
  
Key messages for kids:
  Never share your full name, school, address, or phone online
  Never meet someone in person you only know online
  If something makes you uncomfortable, tell a trusted adult
  Screenshots are permanent — don't send anything you wouldn't want everyone to see
```

---

## Summary

| Threat | What It Is | How to Prevent It |
|--------|-----------|------------------|
| Account takeover | Stolen password used to log into your accounts | Unique passwords + password manager + 2FA |
| Phishing | Fake email/website stealing your password | Check URLs carefully, don't click email links |
| Malware | Malicious software on your device | Don't download pirated software, keep updated |
| Data breach | Company you use gets hacked, exposing your password | Unique passwords so breach of one site doesn't cascade |
| Privacy tracking | Companies/ISPs tracking your behavior | VPN, privacy-focused browser, review app permissions |

**You're now equipped to use the internet safely. Next: let's go deeper — how do computers count, and what IS programming?**
