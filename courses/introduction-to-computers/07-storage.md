# Chapter 07: Storage — Where Your Files Live Permanently

> **"RAM forgets everything when you turn off the computer. Storage remembers everything. Without storage, you'd have to re-download every app and retype every document every single time you turned on your computer."**

---

## Table of Contents

1. [Why We Need Permanent Storage](#1-why-we-need-permanent-storage)
2. [Hard Disk Drives (HDD)](#2-hard-disk-drives-hdd)
3. [Solid State Drives (SSD)](#3-solid-state-drives-ssd)
4. [Flash Storage — Phones and USB Drives](#4-flash-storage--phones-and-usb-drives)
5. [Cloud Storage — Storing on the Internet](#5-cloud-storage--storing-on-the-internet)
6. [How Data Is Actually Stored](#6-how-data-is-actually-stored)
7. [File Formats — What's a .jpg? What's a .mp3?](#7-file-formats--whats-a-jpg-whats-an-mp3)
8. [Summary](#summary)

---

## 1. Why We Need Permanent Storage

```
Without storage:
  Turn off computer → lose everything
  Every time you turn it on: empty, like a newborn computer
  No photos, no apps, no saved games, no documents
  
With storage:
  Turn off computer → data is safe on disk
  Turn it back on → everything where you left it
  Storage is what makes a computer actually useful long-term
```

---

## 2. Hard Disk Drives (HDD)

The hard disk drive was the standard for 50+ years (1950s–2010s).

```
Inside an HDD:
  
  ┌──────────────────────────────────────────────┐
  │                                              │
  │     ┌────────────────────────────────┐       │
  │     │          Platter               │       │
  │     │       (magnetic disk)          │       │
  │     │    ┌────────────────┐          │       │
  │     │    │    Platter 2   │          │       │
  │     │    └────────────────┘          │       │
  │     └────────────────────────────────┘       │
  │                   ↕                          │
  │              Spindle motor                   │
  │              (spins at 5,400–7,200 RPM)      │
  │                                              │
  │     ═══ Read/write arm (like a record        │
  │         player needle)                       │
  └──────────────────────────────────────────────┘
```

**How it stores data:**
- The disk surface is coated with a magnetic material
- Tiny magnetic domains can point one of two ways: North or South
- North = 1, South = 0
- The read/write head changes these magnetic orientations to store data
- The same head detects them to read data

**Characteristics:**
- Cheap: $20–$30 for 1TB
- Slow: 100–150 MB/s read speed
- Fragile: has moving parts — drop it and lose data
- Loud: you can hear it clicking and spinning
- Still used for: cheap bulk storage (backups, NAS servers)

---

## 3. Solid State Drives (SSD)

SSDs replaced HDDs in most laptops and phones. No moving parts.

```
Inside an SSD:
  
  ┌─────────────────────────────────────────────┐
  │                                             │
  │   ┌─────────┐  ┌─────────┐  ┌─────────┐   │
  │   │  NAND   │  │  NAND   │  │  NAND   │   │
  │   │  Flash  │  │  Flash  │  │  Flash  │   │
  │   │  chip   │  │  chip   │  │  chip   │   │
  │   └─────────┘  └─────────┘  └─────────┘   │
  │                                             │
  │   ┌──────────────────────────────────────┐  │
  │   │         Controller chip              │  │
  │   │  (manages reading/writing to chips)  │  │
  │   └──────────────────────────────────────┘  │
  └─────────────────────────────────────────────┘
```

**How it stores data:**
- NAND Flash memory cells trap electrons in a floating gate
- Electrons present = 1, Electrons absent = 0
- No mechanical movement needed — electrons are set/read electronically

**Characteristics:**
- Faster: 500 MB/s (SATA SSD) to 7,000 MB/s (NVMe SSD)
- Silent: no moving parts
- Durable: survives drops
- More expensive: ~$80 for 1TB (vs $25 for HDD)
- SSDs have a limited number of write cycles (but modern ones last 5–10+ years of normal use)

**SSD types:**
```
2.5" SATA SSD:    Same size as laptop HDD, connects via SATA cable
M.2 SATA:         Small stick, slots directly into motherboard
M.2 NVMe:         Same shape but much faster (direct PCIe connection)
U.2:              Enterprise server SSDs
```

---

## 4. Flash Storage — Phones and USB Drives

```
USB Flash Drive ("USB stick" / "Thumb drive"):
  Portable storage you plug into USB port
  Same technology as SSD but slower
  Sizes: 4GB to 2TB
  Cost: ~$10–$30
  Use: transferring files between computers
  
SD Card / MicroSD:
  Used in cameras, some Android phones, Nintendo Switch
  Tiny: about 1cm × 1.5cm
  Sizes: 16GB to 1TB
  
Phone Storage:
  Built-in flash memory (not removable usually)
  iPhone 16: 128GB to 1TB options
  Same NAND flash technology as SSD
  
eMMC (Embedded MultiMediaCard):
  Cheaper, slower flash used in budget phones and devices
  Built directly onto the circuit board
```

---

## 5. Cloud Storage — Storing on the Internet

"The Cloud" is really just someone else's hard drive — huge data centers full of drives storing your files.

```
How cloud storage works:
  
  Your photo is taken on your phone
       ↓
  Uploaded to Apple/Google/Dropbox servers via Internet
       ↓
  Stored on drives in a data center (maybe across the world)
       ↓
  You open your laptop → download the photo from the cloud
  
  The photo exists in:
    Your phone (local copy)
    The data center (cloud copy)
    Your laptop (synced copy)
```

**Popular cloud storage services:**
- **iCloud** (Apple) — 5GB free, $1/month for 50GB
- **Google Drive** — 15GB free, $3/month for 100GB
- **OneDrive** (Microsoft) — 5GB free, included with Office 365
- **Dropbox** — 2GB free, paid plans for more
- **Amazon Photos** — unlimited photo storage for Prime members

**Benefits:**
- Access your files from any device
- Files safe even if your laptop breaks
- Share files with others easily

**Considerations:**
- Requires internet connection
- Monthly cost for more than a few GBs
- Your data is on someone else's servers

---

## 6. How Data Is Actually Stored

At the lowest level, all data is just 0s and 1s (binary). But how do we represent things like photos and text as 0s and 1s?

```
Text:
  Every letter has a number (ASCII/Unicode encoding)
  Letter 'A' = 65 = 01000001 in binary
  "Hello" = 72, 101, 108, 108, 111
  5 bytes to store "Hello"
  
Numbers:
  The number 255 = 11111111 in binary (8 bits = 1 byte)
  
Images:
  Each pixel has a color
  Each color = 3 numbers (Red, Green, Blue) from 0–255
  A 1,920 × 1,080 photo = 2,073,600 pixels
  × 3 bytes per pixel = ~6MB uncompressed
  JPEG compression reduces it to ~300KB (1/20th the size)
  
Audio:
  Sound is captured 44,100 times per second
  Each sample = 2 bytes
  3 minutes of audio = 44,100 × 2 × 180 seconds = ~16MB uncompressed
  MP3 compression reduces it to ~4MB
  
Video:
  Just many images per second, plus audio
  24–60 images per second × hours = massive files
  H.264/H.265 compression makes it manageable
```

---

## 7. File Formats — What's a .jpg? What's an .mp3?

A **file format** defines how data is organized inside a file. The file extension (.jpg, .pdf, .mp3) tells apps what format to expect.

```
Common file formats:

Documents:
  .txt      → plain text (no formatting)
  .docx     → Microsoft Word document
  .pdf      → Portable Document Format (looks same on all devices)
  
Images:
  .jpg/.jpeg → compressed photo (small file, slight quality loss)
  .png       → compressed, no quality loss, good for graphics
  .gif       → animated images (limited colors)
  .heic      → Apple's modern format (smaller than jpg, same quality)
  .raw       → uncompressed camera data (huge files, for professionals)
  
Audio:
  .mp3       → compressed audio (good quality, small size)
  .flac      → lossless audio (perfect quality, large size)
  .wav       → uncompressed audio (perfect quality, huge size)
  .aac       → Apple's format (better than mp3 at same size)
  
Video:
  .mp4       → most common video format
  .mov       → Apple's format
  .avi       → older Windows format
  .mkv       → container that holds video + audio + subtitles
  
Programs:
  .exe       → Windows program
  .app       → macOS program (actually a folder)
  .apk       → Android app package
  .ipa       → iOS app package
```

---

## Summary

| Type | Speed | Cost/TB | Permanent? | Best For |
|------|-------|---------|-----------|---------|
| HDD | Slow | ~$20 | Yes | Cheap bulk storage, backups |
| SSD | Fast | ~$80 | Yes | OS, apps, everyday use |
| NVMe SSD | Very fast | ~$100 | Yes | High-performance computers |
| USB Drive | Medium | ~$15 | Yes | Transferring files |
| SD Card | Medium | ~$15 | Yes | Cameras, phones |
| Cloud | Internet speed | $1-3/month | Yes | Backups, sharing, any-device access |

**Now you know about storage. Let's look at the physical devices you actually use to interact with a computer — keyboards, mice, monitors, and more.**
