# Chapter 41: Android Architecture

> **"Android is the most deployed OS on Earth — over 3 billion active devices. It runs on the Linux kernel but is fundamentally different from a Linux distribution: no glibc, no systemd, no X11, no traditional Unix filesystem layout. It's a clean-room mobile OS built on Linux primitives, Java/Kotlin on a purpose-built runtime, with security as a foundational constraint."**

---

## Table of Contents

1. [Android Stack Overview](#1-android-stack-overview)
2. [The Linux Kernel Layer](#2-the-linux-kernel-layer)
3. [HAL — Hardware Abstraction Layer](#3-hal--hardware-abstraction-layer)
4. [Android Runtime (ART)](#4-android-runtime-art)
5. [The Application Framework](#5-the-application-framework)
6. [Binder IPC](#6-binder-ipc)
7. [Android Security Model](#7-android-security-model)
8. [Android Memory Management](#8-android-memory-management)
9. [Android vs Standard Linux](#9-android-vs-standard-linux)
10. [Summary](#summary)

---

## 1. Android Stack Overview

```
Android Stack (bottom to top):
┌────────────────────────────────────────────────────────────────┐
│  Applications (Calculator, Maps, 3rd-party apps)               │
├────────────────────────────────────────────────────────────────┤
│  Application Framework (Java/Kotlin APIs)                       │
│  ActivityManager, PackageManager, WindowManager, LocationManager│
│  NotificationManager, ContentResolver, InputManager...         │
├────────────────────────────────────────────────────────────────┤
│  Android Runtime (ART) + Core Libraries                        │
│  ART: AOT+JIT compiler for Java/Kotlin bytecode (DEX)          │
│  Core libs: java.lang, java.io, java.util (not full OpenJDK)   │
├──────────────────┬─────────────────────────────────────────────┤
│  C/C++ Libraries │  Android HAL                               │
│  libc (bionic)   │  Camera HAL, Audio HAL, Sensors HAL...     │
│  libssl, zlib    │  (C++ or Java HIDL/AIDL interfaces)         │
│  SQLite, OpenGL  │                                             │
├──────────────────┴─────────────────────────────────────────────┤
│  Linux Kernel (modified for Android)                           │
│  Binder IPC, ION allocator, Wakelocks, PMEM, Low Memory Killer │
└────────────────────────────────────────────────────────────────┘
```

---

## 2. The Linux Kernel Layer

Android uses the Linux kernel but with Android-specific modifications:

**Android-specific kernel additions:**

**Binder:**
A fast IPC mechanism (not in upstream Linux):
```
/dev/binder — device file for Binder communication
Binder IPC is used for EVERYTHING in Android (unlike Unix sockets in Linux)
  App → ActivityManagerService: "start this Activity"
  App → SurfaceFlinger: "draw this frame"
  App → CameraService: "open camera"
  
Performance: ~3× faster than Unix sockets for same message
Why: mapped shared memory, object passing, no extra copies
```

**Wakelocks:**
Mechanism for preventing the CPU from sleeping:
```c
// User space:
wake_lock_init(&my_wakelock, WAKE_LOCK_SUSPEND, "myapp_wakelock");
wake_lock(&my_wakelock);       // prevent sleep
// ... do work that must complete ...
wake_unlock(&my_wakelock);     // allow sleep

// Android's challenge: sleeping saves battery but wakelock abuse drains it
// "Wakelocks" are the #1 cause of battery drain in Android apps
// ADB command: adb shell dumpsys power | grep PARTIAL_WAKE_LOCK
```

**ION memory allocator:**
Manages physically contiguous memory for GPU, camera, and display:
```
GPU, Camera, Display all need memory that:
  - Is physically contiguous (DMA from hardware)
  - Can be shared between CPU and GPU
  
ION provides a unified heap for zero-copy sharing between components
(Being replaced by DMA-BUF heap in newer kernels)
```

**Android Low Memory Killer (LMK):**
Proactively kills background processes before the OOM killer fires:
```
Process groups by oom_adj_score:
  -17:  core system services (never killed)
  -12:  persistent processes (kill only in extreme pressure)
  0:    foreground application (rarely killed)
  100:  visible application
  200:  perceptible (audio in background)
  300:  A services (bound services)
  500:  B services (recently used services)
  900:  cached background activities (kill first)
  1000: completely empty (kill immediately)

LMK kills processes in reverse order when free memory falls below thresholds
```

**SELinux on Android:**
```
Android 5.0 (Lollipop) enforced SELinux everywhere (earlier versions were permissive)
Every app has its own SELinux domain: untrusted_app_N (N = API target level)
System services have specific domains: camera_service, surfaceflinger, etc.
Each domain: tightly scoped policy (only access what it needs)
```

---

## 3. HAL — Hardware Abstraction Layer

Android's HAL allows hardware vendors to provide drivers without exposing kernel interfaces directly to apps:

```
Traditional Linux:            Android with HAL:
  App                           App (Java/Kotlin)
  → kernel syscall               → Framework API
  → kernel driver                → HAL interface (HIDL/AIDL)
                                 → Vendor HAL implementation (C++/Java)
                                 → Kernel driver
```

**HIDL/AIDL (HAL Interface Definition Language):**
```java
// Camera HAL interface (simplified):
interface ICameraDevice {
    open(ICameraDeviceCallback callback) generates (CameraStatus status);
    close();
    getCameraCharacteristics() generates (CameraStatus status, CameraMetadata characteristics);
    configureStreams(StreamConfiguration requestedConfiguration) generates (HalStreamConfiguration halConfiguration);
    processCaptureRequest(CaptureRequest request) generates (CaptureError captureError);
}
```

Hardware vendors implement these interfaces. Android framework calls them. If a vendor changes their kernel driver, only the HAL layer changes (not the framework).

**This is why Android updates are hard:**
```
Google releases Android 15
  ↓
OEM (Samsung/Xiaomi) needs to:
  1. Port new Android code
  2. Verify their HAL implementations still work
  3. Test with their hardware
  4. Pass certification (CTS - Compatibility Test Suite)
  5. Ship update
  
Takes 6-18 months typically → many devices run outdated Android
```

**Project Treble (Android 8+):**
Separated the vendor partition from the system partition:
```
Before Treble: vendor code and Android framework were tightly coupled
After Treble:  /system and /vendor are separate; HIDL defines the interface

Result: Google can update /system independently of vendor's /vendor
Speeds up updates but doesn't eliminate the process entirely
```

---

## 4. Android Runtime (ART)

**History:**
- Original: Dalvik VM — JIT-compiled bytecode (similar to JVM but not the same)
- Android 5.0 (2014): Replaced by ART (Android Runtime)

**ART vs Dalvik:**
```
Dalvik: interpret + JIT (just-in-time compilation at runtime)
  Pro: small on-disk (DEX files are compact)
  Con: slow startup (compilation on first run), consumes CPU/battery

ART: AOT (ahead-of-time compilation) + JIT + profile-guided optimization
  At install time: DEX → OAT (native machine code) via dex2oat
  At runtime: JIT collects hot methods, writes profile
  In background (when idle/charging): AOT recompiles hot paths using profile
  Result: fast startup, fast execution, reasonable disk usage
```

**DEX format:**
Android uses DEX (Dalvik Executable) bytecode, not standard JVM .class files:
```
DEX features vs JVM class files:
  Shared string pool (cross-class deduplication)
  Compressed format (fewer bytes on device storage)
  Single file for multiple classes (vs one .class per class)
  
ART processes:
  .java/.kt source → kotlinc/javac → .class files → d8/r8 → .dex → dex2oat → .oat
```

**Garbage collector:**
```
ART's GC:
  Concurrent Mark-Sweep (CMS): runs mostly concurrent with app
  Semi-space GC: for app startup (fast allocation, accepts pauses)
  Generational GC (newer): young/old generation splitting

Impact on app development:
  "GC pauses" cause jank (frame drops)
  Android 7.0+ GC improved to minimize pauses
  Go (Golang apps via cgo) avoid GC for embedded native code
```

---

## 5. The Application Framework

The Android application framework provides Java/Kotlin APIs that apps use:

**Core system services (running as system processes):**
```
ActivityManagerService (AMS): manages Activities, their lifecycle, tasks, back stack
PackageManagerService (PMS):  manages installed apps, permissions, package metadata
WindowManagerService (WMS):   manages windows, z-order, input dispatch
SurfaceFlinger:               composites all app surfaces into final screen (GPU-driven)
AudioFlinger:                 audio routing, mixing, HAL calls
CameraService:                manages camera hardware, HAL calls
LocationManagerService:       GPS, network location, geofencing
InputManagerService:          touchscreen, key events, routing to correct window
NotificationManagerService:   notification management, alerts
ContentResolver:              databases (contacts, calendar, media store)
```

**The Activity lifecycle (most critical API):**
```
Created → Started → Resumed (foreground, visible, interactive)
                 ↓ home button / new activity
              Paused (visible but not focused)
                 ↓ covered completely
              Stopped (not visible, but alive)
                 ↓ memory pressure
           Destroyed (gone, may be recreated)

App developer MUST save state in onPause/onStop:
Android WILL kill background activities to free memory
If killed, onSaveInstanceState() lets you restore UI state
```

---

## 6. Binder IPC

**Binder** is Android's most important kernel addition — a fast IPC mechanism:

**Why not Unix pipes/sockets?**
```
Unix socket call path:
  App → write() → kernel → read() by service
  Data copied: app buffer → kernel socket buffer → service buffer
  TWO copies + TWO context switches

Binder call path:
  1. App mmap()s a shared memory region
  2. Binder kernel driver:
     - Copies data ONCE (app → shared region)
     - Maps that region read-only into service's address space
     - Sends signal to service thread
  
  Only ONE copy; service reads directly from shared region
```

**The Binder protocol:**
```java
// AIDL: Android Interface Definition Language
// Define service interface:
interface IMyService {
    int add(int a, int b);
    String greet(String name);
}

// Generated stub (in service process):
class MyService extends IMyService.Stub {
    @Override
    public int add(int a, int b) { return a + b; }
    @Override
    public String greet(String name) { return "Hello, " + name; }
}

// Client side (in app process):
IMyService service = IMyService.Stub.asInterface(binder);
int result = service.add(3, 4);  // This goes through Binder!
// Binder serializes (3, 4), sends to service process, deserializes, calls add(), returns
```

**Binder transaction overhead:**
- Same process: direct call (no copy)
- Different process: ~1µs round-trip (very fast for IPC)
- Kernel mode: ~0.5µs for context switch

---

## 7. Android Security Model

**Android sandbox: each app = separate Linux user**
```
App A: uid 10050 (assigned at install time)
App B: uid 10051
App C: uid 10052
System: uid 1000

Each app:
  - Runs in its own process (separate vm_space)
  - Has its own Linux UID (so DAC prevents app-to-app file access)
  - Has its own SELinux domain
  - Has its own storage namespace
  - Communicates with other apps only via:
    Intent (explicit), ContentProvider (shared data), Binder (system services)
```

**Permissions model:**
```
Install-time permissions (old, still used for some):
  App declares in AndroidManifest.xml
  Granted when app is installed (user sees list on older Android)

Runtime permissions (Android 6+):
  Dangerous permissions require explicit user grant at runtime:
  CAMERA, MICROPHONE, LOCATION, CONTACTS, CALENDAR, STORAGE, PHONE

Permission groups:
  Location: ACCESS_FINE_LOCATION, ACCESS_COARSE_LOCATION
  Camera: CAMERA
  Storage: READ_EXTERNAL_STORAGE, WRITE_EXTERNAL_STORAGE (pre-Android 10)

Scoped storage (Android 10+):
  Apps can no longer access arbitrary files via path
  Must use MediaStore API for photos/media
  Must use SAF (Storage Access Framework) for documents
```

**Verified Boot:**
```
Boot ROM → Bootloader (Android Verified Boot 2.0) → Kernel → System
Each stage verifies the signature of the next before executing it

DM-Verity: kernel verifies /system partition in real-time (hash tree)
If any block in /system is modified → boot halts
Makes persistent rooting very difficult
```

---

## 8. Android Memory Management

**Limited RAM, aggressive reclamation:**
```
Phone has 4GB RAM. Typical pressure:
  Android OS + system services: ~1.5GB
  Running app: 200-500MB
  Cached apps (background): 200MB each × many apps
  
LMK constantly trimming cache to free memory for foreground app
```

**ZRAM (compressed swap):**
```
Android uses ZRAM — compressed RAM as virtual swap:
  Page evicted from physical RAM → compressed 2-3× → stored in ZRAM
  No disk swap (flash wears out fast with swap)
  ZRAM lives in RAM but takes less space
  
Typical: 2GB ZRAM on 4GB phone → ~1.5GB effective extra "memory"
```

**Memory trim callbacks:**
```java
// App can respond to memory pressure:
@Override
public void onTrimMemory(int level) {
    if (level >= TRIM_MEMORY_MODERATE) {
        // Release non-essential caches
        imageCache.evictAll();
    }
    if (level >= TRIM_MEMORY_COMPLETE) {
        // Release everything — we're about to be killed
        clearAllData();
    }
}
```

---

## 9. Android vs Standard Linux

| Feature | Standard Linux | Android |
|---------|---------------|---------|
| libc | glibc | bionic (smaller, BSD-licensed, Android-tuned) |
| Init system | systemd | init (Android-specific) |
| IPC | D-Bus, Unix sockets | Binder (primary), Unix sockets (secondary) |
| GUI | X11 / Wayland | SurfaceFlinger + SurfaceView (direct GPU) |
| Service manager | systemd units | ServiceManager (Binder registry) |
| Package manager | apt/dnf/pacman | PackageManager + APK format |
| App runtime | Native processes | ART (Java/Kotlin) |
| Security | SELinux (optional), capabilities | SELinux (mandatory), per-app UID, runtime permissions |
| Shell | /bin/bash or zsh | /system/bin/sh (toybox-based) |
| Init process | /sbin/init → systemd | /init (Android init, reads .rc files) |
| Filesystems | ext4 default, btrfs optional | ext4, F2FS (flash-optimized) |
| Swap | Disk-based | ZRAM only (no disk swap) |

---

## Summary

| Concept | Description |
|---------|------------|
| Android stack | Linux kernel → HAL → bionic/ART → Framework → Apps |
| Binder | Android's primary IPC; single-copy via shared memory; 1µs round-trip |
| ART | Android Runtime; AOT + JIT compilation of Dalvik DEX bytecode |
| DEX | Dalvik Executable: Android's bytecode format (not JVM class files) |
| AIDL | Android Interface Definition Language; generates Binder stubs/proxies |
| SELinux | Mandatory on Android 5+; each app domain, system service domain |
| Per-app UID | Each app gets unique Linux UID; DAC enforces app isolation |
| HAL | Hardware Abstraction Layer; vendor-provided HIDL/AIDL implementations |
| HIDL | HAL Interface Definition Language; separates /system from /vendor |
| LMK | Android Low Memory Killer; kills background apps before OOM |
| ZRAM | Compressed RAM acting as swap; no disk swap on Android |
| Wakelocks | Prevent CPU sleep; abuse causes battery drain |
| Verified Boot | Chain of trust from ROM through kernel to system partition |
| Scoped storage | Android 10+; apps can't access arbitrary file paths |
| AMS | ActivityManagerService: manages app lifecycle, task stack |
