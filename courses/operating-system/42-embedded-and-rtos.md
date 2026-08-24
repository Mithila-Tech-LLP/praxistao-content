# Chapter 42: Embedded and Real-Time Operating Systems

> **"An airbag controller must detect a crash and fire in 15 milliseconds — every time, without exception. A missed deadline doesn't cause slow performance; it causes death. This is why real-time operating systems exist: not for speed, but for predictability. Guaranteed response times are non-negotiable."**

---

## Table of Contents

1. [What Is Embedded Systems?](#1-what-is-embedded-systems)
2. [Real-Time Requirements](#2-real-time-requirements)
3. [RTOS Architecture](#3-rtos-architecture)
4. [Scheduling in RTOS](#4-scheduling-in-rtos)
5. [Famous RTOS Examples](#5-famous-rtos-examples)
6. [FreeRTOS — The Most Popular Embedded RTOS](#6-freertos--the-most-popular-embedded-rtos)
7. [QNX — Microkernel RTOS](#7-qnx--microkernel-rtos)
8. [Embedded Linux](#8-embedded-linux)
9. [OS for Microcontrollers vs Microprocessors](#9-os-for-microcontrollers-vs-microprocessors)
10. [Summary](#summary)

---

## 1. What Is Embedded Systems?

An **embedded system** is a computer built into a larger device to perform a specific, dedicated function:

```
Examples:
  Automotive: Engine control unit (ECU), airbag controller, ABS, transmission
  Medical: Pacemaker, insulin pump, MRI controller, ventilator
  Industrial: PLC (Programmable Logic Controller), robot controller
  Consumer: Microwave oven, washing machine, TV remote
  Infrastructure: Traffic lights, power grid controllers
  Communications: Router firmware, modem, switch ASIC
  IoT: Smart thermostat, sensor node, smart speaker
  Aerospace: Flight controller, satellite attitude control
```

**Embedded system characteristics:**
```
Resource-constrained:
  CPU: 8-bit to 32-bit (ARM Cortex-M0 = 256KB flash, 32KB RAM)
  Memory: kilobytes to megabytes (not gigabytes)
  Power: milliwatts to watts (battery-powered for years)

Single-purpose: runs ONE application (not a general-purpose OS)
Real-time: must respond within guaranteed time bounds
Reliability: must work for years without reset (no blue screen allowed)
Often no display: no user interface at all
```

---

## 2. Real-Time Requirements

**Hard Real-Time:**
A missed deadline is a system FAILURE — not just reduced performance:
```
Airbag controller: must fire within 15ms of crash detection
  Miss deadline: airbag doesn't deploy → injury/death
  
Fly-by-wire aircraft: must update control surfaces every 20ms
  Miss deadline: aircraft becomes uncontrollable
  
ABS (Anti-lock Braking): must sample wheel speed every 5ms
  Miss deadline: incorrect brake pressure → loss of control
```

**Soft Real-Time:**
A missed deadline degrades quality but doesn't cause failure:
```
Video playback: should decode frame every 16.7ms (60 FPS)
  Miss deadline: frame drop → visible stutter (annoying, not dangerous)
  
Audio streaming: should output buffer every 5ms
  Miss deadline: audio glitch (audible pop)
  
Web server: should respond within 200ms
  Miss deadline: user sees "slow" (bad UX, not dangerous)
```

**Key metric: Worst-Case Execution Time (WCET)**
```
Hard RT system guarantee:
  WCET(task) ≤ deadline

General-purpose OS (Linux, Windows):
  Average case is good
  WORST case is NOT bounded (scheduler can preempt you for arbitrary time)
  → NOT suitable for hard real-time
  
RTOS guarantee:
  Every task has a known worst-case execution time
  Scheduler guarantees tasks meet their deadlines
  Interrupt latency is bounded (typically <10 microseconds)
```

---

## 3. RTOS Architecture

**RTOS design principles:**

**1. Minimal and deterministic:**
```
General OS: thousand of system calls, complex scheduler, virtual memory, page faults
RTOS: ~20-100 API functions, simple fixed-priority scheduler, NO virtual memory
→ Every operation takes a bounded, predictable time
```

**2. No dynamic memory (often):**
```
malloc() in general OS: variable time (depends on heap state, fragmentation)
RTOS typical approach:
  - Fixed-size memory pools (alloc = O(1), free = O(1))
  - Static allocation (define all data structures at compile time)
  - Stack checking (stack overflow detection)
```

**3. Predictable context switch:**
```
Linux context switch: ~1-5 microseconds (variable)
RTOS context switch: fixed time, typically ~1-3 microseconds
  AND: guaranteed to happen within a bounded time after trigger
```

**4. Priority inversion prevention:**
A common RTOS bug where a low-priority task holds a resource needed by a high-priority task:
```
Priority inversion scenario:
  HP (High Priority) task needs mutex M
  LP (Low Priority) task holds mutex M
  MP (Medium Priority) task preempts LP
  
  HP waits → LP can't run (preempted by MP) → deadlock-like situation
  
Solution: Priority inheritance
  When LP holds a mutex that HP needs:
    LP temporarily gets HP's priority
    LP completes quickly, releases mutex
    LP's priority drops back to normal
    HP can now acquire mutex
```

---

## 4. Scheduling in RTOS

**Fixed-Priority Preemptive Scheduling:**
Most RTOSes use this. Each task has a fixed priority. Highest-priority runnable task always runs.

```
Tasks:
  Task A: priority 1 (highest), period 10ms, execution 2ms
  Task B: priority 2, period 20ms, execution 5ms
  Task C: priority 3 (lowest), period 50ms, execution 10ms

Timeline:
0ms:   A runs (2ms)
2ms:   B runs
5ms:   A preempts B (new period started)
7ms:   B resumes
10ms:  A preempts B again
12ms:  B finishes, C starts
15ms:  A preempts C
17ms:  C resumes
20ms:  A preempts C, B also due
17ms:  A (2ms), then B (5ms), then C resumes
...
```

**Rate Monotonic Scheduling (RMS):**
The optimal priority assignment for periodic tasks:
```
Rule: shorter period → higher priority

Theoretical utilization bound:
  n tasks can be scheduled if: Σ(Ci/Ti) ≤ n(2^(1/n) - 1)
  For large n: approaches ln(2) ≈ 69%

  Task A: 2ms execution / 10ms period = 0.2
  Task B: 5ms execution / 20ms period = 0.25
  Task C: 10ms execution / 50ms period = 0.2
  Total: 0.65 ≤ 0.69 ✓ → schedulable!
```

**EDF (Earliest Deadline First):**
Dynamic priority based on which task's deadline is soonest:
```
Optimal for single CPU: can achieve 100% utilization
More complex to implement: priorities change dynamically
Used in: SCHED_DEADLINE in Linux, multimedia systems
```

---

## 5. Famous RTOS Examples

| RTOS | Users | Key feature |
|------|-------|-------------|
| FreeRTOS | IoT, hobbyist, medical | Open source, runs on tiny microcontrollers |
| Zephyr | IoT, wearables (Amazon Echo) | Linux Foundation, modern, well-maintained |
| QNX | Automotive (BMW, Ford), medical | POSIX-compatible microkernel, very reliable |
| VxWorks | Aerospace (Mars rovers), military | Certified for DO-178C (aviation safety) |
| INTEGRITY | Aerospace, automotive | ARINC 653 certified, partitioned |
| ThreadX | IoT chips, printers, cameras | Very lightweight, Azure RTOS (Microsoft) |
| RTEMS | NASA spacecraft (Curiosity, Perseverance Mars rover) | Open source, space-proven |
| Tizen | Samsung TVs, wearables | Linux-based, Samsung |
| AUTOSAR | Modern cars (most OEMs) | Standard for automotive software |

**VxWorks on Mars:**
```
NASA Mars rovers (Spirit, Opportunity, Curiosity, Perseverance):
  → Run VxWorks RTOS
  
Mars Science Laboratory (Curiosity) computer:
  CPU: 200MHz RAD750 (radiation-hardened PowerPC)
  RAM: 256MB DRAM
  Flash: 2GB
  OS: VxWorks

The 2004 Spirit rover had a memory management bug:
  → VxWorks process table filled up
  → Rover rebooted ~128 times
  → Mission controllers uploaded a software fix via deep space network
  → Rover survived and operated for 6+ years
```

---

## 6. FreeRTOS — The Most Popular Embedded RTOS

**FreeRTOS** runs on over 50 microcontroller architectures with just 6-10KB of ROM:

**Core API:**
```c
#include "FreeRTOS.h"
#include "task.h"
#include "queue.h"

// Create a task:
void temperature_task(void *pvParameters) {
    while (1) {
        int temp = read_sensor();
        printf("Temperature: %d°C\n", temp);
        vTaskDelay(pdMS_TO_TICKS(1000));  // delay 1 second
    }
}

// Create an LED blink task:
void led_task(void *pvParameters) {
    while (1) {
        toggle_led();
        vTaskDelay(pdMS_TO_TICKS(500));   // blink every 500ms
    }
}

int main(void) {
    // Create tasks:
    xTaskCreate(temperature_task, "TempTask", 256, NULL, 2, NULL);
    //           function          name        stack  param  priority  handle
    
    xTaskCreate(led_task, "LEDTask", 128, NULL, 1, NULL);
    
    // Start the scheduler (never returns):
    vTaskStartScheduler();
    
    // Should never reach here:
    for (;;);
}
```

**FreeRTOS Inter-Task Communication:**
```c
// Queue: pass data between tasks
QueueHandle_t xQueue = xQueueCreate(10, sizeof(uint32_t));

// Sender task:
uint32_t data = 42;
xQueueSend(xQueue, &data, portMAX_DELAY);  // blocks if queue full

// Receiver task:
uint32_t received;
xQueueReceive(xQueue, &received, portMAX_DELAY);  // blocks if queue empty

// Semaphore: synchronization
SemaphoreHandle_t xSemaphore = xSemaphoreCreateBinary();
xSemaphoreGive(xSemaphore);    // signal
xSemaphoreTake(xSemaphore, portMAX_DELAY);  // wait

// Mutex: mutual exclusion
SemaphoreHandle_t xMutex = xSemaphoreCreateMutex();
xSemaphoreTake(xMutex, portMAX_DELAY);  // lock
// ... critical section ...
xSemaphoreGive(xMutex);                 // unlock
```

**FreeRTOS memory usage:**
```
Kernel: ~6KB ROM, ~4KB RAM
Per task stack: defined at creation (256 words = 1KB for ARM32)
Total for simple app: 15-20KB ROM, 8-16KB RAM
→ Fits comfortably on STM32F103 (128KB flash, 20KB RAM)
```

---

## 7. QNX — Microkernel RTOS

**QNX** is the most commercially successful microkernel RTOS. Used in:
- Automotive infotainment systems (BMW iDrive, Audi, Ford Sync)
- Medical devices (ultrasound, patient monitors)
- Industrial automation
- Network infrastructure

**QNX architecture:**
```
QNX Microkernel: only ~12KB of code in Ring 0!
  - Message passing (synchronous)
  - Process scheduling
  - Interrupt handling
  
Everything else runs as user-space processes:
  devb-ram:   RAM disk driver
  devb-eide:  IDE/ATA disk driver
  io-blk:     Block I/O manager
  io-fs:      File system manager
  io-net:     Network manager
  devc-con:   Console driver
  
If a driver crashes: only that process dies, kernel continues
  → kernel restarts the driver automatically
  → system never needs a reboot for driver failures
```

**QNX message passing:**
```c
// QNX IPC: synchronous message passing
// Server creates channel:
int chid = ChannelCreate(0);

// Client sends message:
int coid = ConnectAttach(0, server_pid, chid, _NTO_SIDE_CHANNEL, 0);
MsgSend(coid, &send_buf, sizeof(send_buf), &recv_buf, sizeof(recv_buf));
// Blocks until server replies

// Server receives and replies:
int rcvid = MsgReceive(chid, &send_buf, sizeof(send_buf), NULL);
// ... process request ...
MsgReply(rcvid, EOK, &reply_buf, sizeof(reply_buf));
```

**Why QNX for safety-critical systems:**
- Certified IEC 61508 (functional safety), ISO 26262 (automotive)
- Mandatory access control
- Process isolation (driver bug doesn't crash system)
- Deterministic message passing latency

---

## 8. Embedded Linux

Many embedded systems use Linux when:
- More capability is needed (networking, filesystems, broad driver support)
- Real-time is "soft" (consumer electronics, IoT gateways, routers)
- Budget allows more capable hardware (ARM Cortex-A series, not just Cortex-M)

**Real-Time Linux — PREEMPT_RT:**
```
Standard Linux: not real-time (unbounded latencies)
Linux with PREEMPT_RT patch:
  - Converts most spinlocks to mutexes (preemptible)
  - Interrupt handlers run in kernel threads (preemptible by higher-priority RT task)
  - Achieves 10-50 microsecond worst-case latency

Use:
  Audio production (professional DAW, audio servers like JACK)
  Industrial automation (Beckhoff TwinCAT on Linux)
  CNC controllers (LinuxCNC)
  
PREEMPT_RT merged into mainline Linux kernel 6.12 (2024) — no separate patch needed!
```

**Buildroot and Yocto:**
```
Building minimal Linux for embedded:
  Buildroot: simple, fast, generates minimal root FS
    make menuconfig → configure → make → get bootable image
    
  Yocto/OpenEmbedded: complex but very flexible
    Layers: BSP layer (board support) + distro layer + app layer
    Recipe-based: precise control over every package version
    Used by: automotive (Automotive Grade Linux), many IoT platforms

Cross-compilation:
  x86 laptop cannot run ARM binary directly
  → cross-compiler: gcc for ARM target, runs on x86 host
  arm-linux-gnueabihf-gcc -o myapp myapp.c
  → copy to ARM device → runs on device
```

---

## 9. OS for Microcontrollers vs Microprocessors

```
Microcontroller (MCU):          Microprocessor (MPU / SoC):
  ARM Cortex-M0/M3/M4             ARM Cortex-A53/A72/A76
  Flash: 16KB - 2MB               Storage: 4GB+ (eMMC, SD)
  RAM: 2KB - 512KB                 RAM: 256MB - 8GB
  No MMU (no virtual memory)      Full MMU (virtual memory)
  Single application              Multiple processes
  
  OS: FreeRTOS, Zephyr, Mbed      OS: Embedded Linux, Android, QNX
  
MCU (no MMU) challenges:
  - All code runs in same address space
  - Bug in one task can corrupt any memory
  - No demand paging (all code must fit in RAM)
  - Fixed stack size (overflow corrupts adjacent memory)
  
MPU (with MMU) benefits:
  - Process isolation (virtual memory)
  - Can run Linux with full POSIX API
  - Support filesystems, networking, etc.
```

---

## Summary

| Concept | Description |
|---------|------------|
| Hard real-time | Missed deadline = system failure (airbag, flight controller) |
| Soft real-time | Missed deadline = degraded quality (video, audio) |
| WCET | Worst-Case Execution Time; must be ≤ deadline for hard RT |
| Fixed-priority preemption | Highest-priority runnable task always runs; others wait |
| Rate Monotonic | Optimal static priority assignment: shorter period → higher priority |
| EDF | Earliest Deadline First; optimal dynamic scheduling |
| Priority inversion | Low-priority task blocks high-priority; fix with priority inheritance |
| FreeRTOS | Open-source RTOS; 6-10KB; runs on MCUs with 32KB RAM |
| QNX | Microkernel RTOS; drivers in user space; used in automotive/medical |
| VxWorks | Aerospace-certified RTOS; DO-178C; used in Mars rovers |
| PREEMPT_RT | Linux kernel patch (now mainline) for soft real-time; <50µs latency |
| Buildroot | Minimal Linux builder for embedded targets |
| MCU | Microcontroller: no MMU, small RAM, runs RTOS or bare metal |
| MPU/SoC | Application processor: full MMU, runs Linux/Android |
