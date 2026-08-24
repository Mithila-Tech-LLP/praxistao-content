# Chapter 46: Microcontrollers — Processors in Everything

A microcontroller is a complete computer on a single chip: CPU, RAM, Flash storage, and peripherals like UART, SPI, I2C, ADC — all integrated. Unlike a microprocessor (which needs external memory and peripherals), a microcontroller is self-contained. This makes microcontrollers cheap, small, and power-efficient. They're in your car (dozens of ECUs), your kitchen appliances, your keyboard, your USB charger, your smoke detector, your medical devices, and millions of other places you don't think about. The global microcontroller market ships over 30 billion units per year.

## Table of Contents

1. [Microcontroller vs Microprocessor](#1-microcontroller-vs-microprocessor)
2. [AVR — Arduino's Heart](#2-avr--arduinos-heart)
3. [PIC — The Industrial Workhorse](#3-pic--the-industrial-workhorse)
4. [ARM Cortex-M — The Modern Standard](#4-arm-cortex-m--the-modern-standard)
5. [STM32 — The Most Popular MCU Family](#5-stm32--the-most-popular-mcu-family)
6. [Peripherals That Make MCUs Useful](#6-peripherals-that-make-mcus-useful)
7. [Real-Time Operating Systems (RTOS)](#7-real-time-operating-systems-rtos)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Microcontroller vs Microprocessor

```
Microprocessor (e.g., Intel Core i7):
  CPU on chip
  Needs: external RAM, external storage, external USB/PCI controller
  Runs: operating system, multiple complex applications
  Cost: $100–$1000
  Power: 15–300W
  
Microcontroller (e.g., STM32F103):
  CPU + RAM + Flash + peripherals, all on one chip
  Runs: bare-metal firmware or RTOS
  Cost: $0.20–$20
  Power: mW to µW range
  
Examples of what's inside a microcontroller:
  ┌─────────────────────────────────────────────────────┐
  │                 Microcontroller                      │
  │  ┌──────────┐  ┌───────────┐  ┌────────────────┐  │
  │  │  CPU     │  │  Flash    │  │  SRAM          │  │
  │  │ (Cortex-M│  │  (Program)│  │  (Working RAM) │  │
  │  │  or AVR) │  │  64KB     │  │  20KB          │  │
  │  └──────────┘  └───────────┘  └────────────────┘  │
  │                                                      │
  │  ┌──────────┐  ┌───────────┐  ┌────────────────┐  │
  │  │  Timers  │  │  UART/SPI │  │  ADC/DAC       │  │
  │  │          │  │  I2C      │  │  (analog I/O)  │  │
  │  └──────────┘  └───────────┘  └────────────────┘  │
  │                                                      │
  │  GPIO pins (general-purpose digital I/O)            │
  │  Clock management (PLL, prescalers)                  │
  │  Watchdog timer, reset logic                         │
  └─────────────────────────────────────────────────────┘
```

**Use cases**: Microcontrollers run deterministic, dedicated tasks:
- Read a temperature sensor every second
- Control a motor's speed based on feedback
- Implement a USB HID device (keyboard, mouse)
- Manage battery charging
- Run a simple communication protocol

They are **not** designed for running browsers, operating systems, or general-purpose computing.

### Quick Check
> 1. Name four things a microcontroller has on-chip that a microprocessor does not.
> 2. Why would you choose a $0.50 microcontroller over a $50 Raspberry Pi for a smoke detector?
> 3. What is "bare-metal" programming and when is it appropriate?

---

## 2. AVR — Arduino's Heart

**AVR** is Atmel's (now Microchip's) family of 8-bit RISC microcontrollers. The ATmega328P is the chip in the original Arduino Uno, and it made microcontroller programming accessible to millions of students and hobbyists.

**ATmega328P specs:**
```
Architecture: Modified Harvard, 8-bit AVR RISC
Clock: 16 MHz (with external crystal)
Flash: 32KB (program memory)
SRAM: 2KB (data memory)
EEPROM: 1KB (non-volatile data)
I/O: 23 digital GPIO pins, 6 analog inputs
Peripherals: UART, SPI, I2C, 6 PWM channels, ADC (10-bit)
Power: 3.3mA active at 5V (16.5mW), 0.1µA power-down mode
Cost: ~$2 each
```

**Why Harvard architecture**: AVR uses strict Harvard — program (Flash) and data (SRAM) are separate address spaces. The program counter is 16 bits (addressing 32K of 16-bit instructions). Data pointers are 16 bits (addressing 64KB of SRAM, though physical SRAM is only 2KB). A pointer to a data variable means nothing in the instruction space.

**AVR instruction set**:
- Most instructions execute in 1 clock cycle
- 32 × 8-bit general-purpose registers (R0–R31)
- Special pairs: X (R26:R27), Y (R28:R29), Z (R30:R31) for 16-bit address pointers
- Simple: no pipelining complexity, no caches, execution is predictable to the cycle

**Arduino framework** abstracts away the hardware — `digitalWrite(13, HIGH)` turns on the LED, `analogRead(A0)` reads a sensor. This hides the register-level complexity, though real embedded engineers still write directly to hardware registers for performance.

### Quick Check
> 1. What are the Flash, SRAM, and EEPROM sizes of the ATmega328P?
> 2. Why is AVR's Harvard architecture a limitation for C programs with string constants?
> 3. What are the X, Y, and Z register pairs in AVR used for?

---

## 3. PIC — The Industrial Workhorse

**PIC** (Peripheral Interface Controller, now just "PIC") is Microchip Technology's microcontroller family — one of the most widely used MCU families in the world, especially in industrial and consumer electronics.

**PIC architecture** (8-bit PIC10/12/16/18):
- Modified Harvard architecture (like AVR)
- Reduced instruction set: PIC10/12/16 have only 33–35 instructions
- Accumulator-based (not multi-register like AVR): most operations go through W (working register)
- Program memory: 12–22 bit instruction words (wider than data bus)
- 8-bit data bus

**PIC popularity reasons:**
- Extremely cheap: PIC10F220 under $0.25 in volume
- Simple instruction set = easy to learn
- Large portfolio from 6-pin tiny chips to complex 100-pin devices
- Strong industrial pedigree and reputation for reliability
- Microchip provides free development tools (MPLAB X IDE, XC compilers)

**Modern PIC32 (32-bit)**:
- Based on MIPS architecture (Microchip acquired MIPS)
- 32-bit, 80–200 MHz
- USB, Ethernet, CAN
- Targets industrial control, medical devices

**Where PICs appear**: Garage door openers, appliances, toy controllers, car key fobs, power supply management. Anywhere you need simple, reliable, cheap processing.

### Quick Check
> 1. What makes PIC microcontrollers so popular for industrial applications?
> 2. PIC uses an "accumulator-based" architecture. What does this mean compared to AVR's register file?
> 3. What is the cheapest PIC microcontroller and what would you use it for?

---

## 4. ARM Cortex-M — The Modern Standard

**ARM Cortex-M** is the dominant microcontroller architecture in new designs. Every major MCU vendor (ST, NXP, Nordic, Renesas, Silicon Labs, Cypress/Infineon) uses Cortex-M.

**Cortex-M family:**

| Core | Architecture | Pipeline | FPU | Use Case |
|------|-------------|----------|-----|----------|
| Cortex-M0 | ARMv6-M | 3-stage | No | Ultra-small, ultra-cheap |
| Cortex-M0+ | ARMv6-M | 2-stage | No | Lowest power |
| Cortex-M3 | ARMv7-M | 3-stage | No | General embedded |
| Cortex-M4 | ARMv7E-M | 3-stage | Optional | DSP, motor control, audio |
| Cortex-M7 | ARMv7E-M | 6-stage, OOO | Yes | High-performance embedded |
| Cortex-M33 | ARMv8-M | 3-stage | Optional | IoT with security (TrustZone) |
| Cortex-M55 | ARMv8.1-M | 3-stage | Yes + Helium SIMD | Edge AI |
| Cortex-M85 | ARMv8.1-M | 8-stage OOO | Yes + Helium | High-perf embedded |

**Key Cortex-M features:**
- **Thumb-2 ISA**: 16-bit and 32-bit instruction mix for code density
- **NVIC (Nested Vectored Interrupt Controller)**: Hardware-managed interrupt prioritization and tail-chaining
- **CMSIS (Cortex Microcontroller Software Interface Standard)**: Vendor-neutral hardware abstraction layer
- **MPU (Memory Protection Unit)**: Basic memory protection (no MMU/virtual memory)

**NVIC tail-chaining**: When one interrupt finishes and another is waiting, instead of restoring all registers (12 cycles) then saving them again, the CPU goes directly to the next interrupt handler — saving 12 cycles of overhead. Important for high-interrupt-rate applications.

**Thumb-2 encoding**: ARM32 instructions are 32 bits; ARM Thumb is 16-bit encodings of common instructions. Thumb-2 mixes both. Most C code compiles to Thumb-2, getting 30–40% smaller code than ARM32 while retaining nearly the same performance.

### Quick Check
> 1. What Cortex-M core would you use for an ML inference task?
> 2. What is "NVIC tail-chaining" and why does it matter for interrupt-heavy applications?
> 3. Why is Thumb-2 encoding important for microcontroller applications?

---

## 5. STM32 — The Most Popular MCU Family

**STMicroelectronics STM32** is the most popular professional microcontroller family, used in everything from consumer electronics to industrial robots.

**STM32 naming scheme:**
```
STM32 F 103 C 8 T 6
      │ │   │ │ │ │
      │ │   │ │ │ └─ Temp range (6 = industrial -40 to 85°C)
      │ │   │ │ └─── Package (T = LQFP48)
      │ │   │ └───── Flash density (8 = 64KB)
      │ │   └─────── Pin count (C = 48 pins)
      │ └─────────── Specific number
      └───────────── Family (F = Mainstream, H = High Performance, L = Low Power, etc.)
```

**Key STM32 families:**
- **STM32F0/F1**: Entry level, Cortex-M0/M3, 48–72 MHz
- **STM32F4**: Cortex-M4F, 168 MHz, DSP — popular for audio, motor control
- **STM32H7**: Cortex-M7, 480 MHz, dual-core (M7+M4) — high-performance
- **STM32L4/L5**: Cortex-M4/M33, ultra-low power (<10 µA/MHz)
- **STM32U5**: Cortex-M33, ultra-low power with TrustZone security
- **STM32WB**: Cortex-M4 + M0 (for BLE), built-in BT LE radio

**STM32CubeMX**: GUI tool to configure STM32 peripherals visually and generate initialization code. Drag-drop to assign GPIO pins, configure clock tree, enable DMA channels — reduces bring-up time from weeks to days.

**HAL (Hardware Abstraction Layer)**: STM32's vendor HAL allows writing code like:
```c
HAL_GPIO_WritePin(GPIOA, GPIO_PIN_5, GPIO_PIN_SET);  // set PA5 high
HAL_UART_Transmit(&huart2, (uint8_t*)"Hello\n", 6, 100);  // send over UART
```

Without HAL (direct register access):
```c
GPIOA->BSRR = GPIO_BSRR_BS5;  // set bit 5 of BSRR in GPIOA
```

Both work; HAL is more readable; direct register access is faster and more predictable.

### Quick Check
> 1. What does "STM32F103C8T6" mean?
> 2. What STM32 family would you choose for a battery-powered sensor?
> 3. What is the difference between STM32's HAL and direct register access?

---

## 6. Peripherals That Make MCUs Useful

MCUs are useful because of their integrated peripherals — hardware modules that handle common communication and I/O tasks.

**UART (Universal Asynchronous Receiver/Transmitter)**: Serial communication with no clock signal. Just TX and RX wires. Used for: debug output, GPS modules, Bluetooth modules (HC-05), AT command communication.

**SPI (Serial Peripheral Interface)**: 4-wire synchronous serial: MOSI, MISO, SCK, CS. Fast (up to 50 MHz). Used for: SD cards, SPI Flash, TFT displays, IMU sensors.

**I2C (Inter-Integrated Circuit)**: 2-wire synchronous: SCL (clock), SDA (data). Multiple devices share the bus, each with a 7-bit address. Slower (100kHz–3.4MHz) but simpler wiring. Used for: OLED displays, temperature sensors, accelerometers, RTCs.

**PWM (Pulse Width Modulation)**: Digital output that switches on/off rapidly. Average voltage ∝ duty cycle. Used for: motor speed control, LED brightness, servo position, audio output.

**ADC (Analog to Digital Converter)**: Converts analog voltage to digital value. Resolution: 12-bit (0–4095 range) typically. Sample rate: 1–10 MHz. Used for: potentiometers, temperature sensors, microphones, light sensors.

**DMA (Direct Memory Access)**: Hardware that copies data between peripherals and RAM without CPU involvement. Example: UART receives bytes into a DMA buffer; CPU is interrupted only when the buffer is full — not per-byte. Critical for high-throughput peripherals.

**Timers**: Hardware counters used for: generating precise delays, PWM, input capture (measure pulse width), real-time clocking.

```
Example: DMA + UART for efficient serial receiving:
  Without DMA: 
    Byte arrives → interrupt → CPU saves context → reads byte → stores → resumes
    At 115200 baud: 11,520 interrupts/second, each ~20 cycles overhead = 230K cycles/sec wasted
  
  With DMA:
    Configure DMA: UART data register → RAM buffer
    CPU does other work
    When buffer full → single interrupt → CPU processes entire buffer
    At 64-byte buffer: 1,125 interrupts/sec instead of 11,520 — 10× less overhead
```

### Quick Check
> 1. What is the difference between I2C and SPI in terms of wiring and speed?
> 2. How does PWM control motor speed?
> 3. Why is DMA important for high-baud-rate UART communication?

---

## 7. Real-Time Operating Systems (RTOS)

For more complex MCU applications, a **Real-Time Operating System (RTOS)** provides:
- **Task scheduling**: Multiple tasks running cooperatively or preemptively
- **Timing guarantees**: Tasks execute within specified deadlines
- **Inter-task communication**: Queues, semaphores, mutexes
- **Small footprint**: FreeRTOS fits in 6-10KB of Flash

**FreeRTOS**: Most popular open-source RTOS, used in Amazon products, AWS IoT. Amazon acquired FreeRTOS in 2017.

```c
// FreeRTOS example: two tasks running concurrently
void sensorTask(void* param) {
    for (;;) {
        int reading = readTemperature();
        xQueueSend(sensorQueue, &reading, portMAX_DELAY);
        vTaskDelay(pdMS_TO_TICKS(1000));  // wait 1 second
    }
}

void displayTask(void* param) {
    int reading;
    for (;;) {
        if (xQueueReceive(sensorQueue, &reading, portMAX_DELAY)) {
            updateDisplay(reading);
        }
    }
}
// Both tasks run "simultaneously" via preemptive scheduling
```

**Scheduling**: FreeRTOS uses priority-based preemptive scheduling. High-priority tasks preempt low-priority ones. The scheduler runs in the SysTick timer interrupt (every 1ms by default).

**RTOS vs bare-metal**: For simple, single-task firmware (read sensor → send data → sleep), bare-metal is simpler. For devices with multiple concurrent operations (handle BLE connection + log data + update display + process commands), RTOS simplifies the code significantly.

**Zephyr**: Linux Foundation's RTOS, used in many IoT devices. More features than FreeRTOS, actively developed with industry support (Intel, Nordic, Google).

### Quick Check
> 1. What problem does an RTOS solve that bare-metal programming cannot easily handle?
> 2. What is FreeRTOS and who owns it?
> 3. In the FreeRTOS example above, how does `vTaskDelay()` allow the display task to run while the sensor task is waiting?

---

## Summary

- **Microcontrollers** integrate CPU, Flash, SRAM, and peripherals on one chip. They are cheap, small, and power-efficient — used in 30+ billion units per year.
- **AVR** (ATmega328P): 8-bit Harvard RISC, used in Arduino. 32KB Flash, 2KB SRAM, simple and educational.
- **PIC**: Microchip's ultra-cheap industrial MCU. 33–35 instruction RISC. Used in appliances, consumer electronics.
- **ARM Cortex-M** (M0 to M85): The modern MCU standard, 32-bit, ranging from ultra-cheap M0 to high-performance M7/M85.
- **STM32**: Most popular professional MCU family. Wide range from low-power L-series to high-performance H7.
- **Peripherals**: UART, SPI, I2C, PWM, ADC, DMA, timers — hardware that handles communication and I/O.
- **RTOS** (FreeRTOS, Zephyr): Allows multiple concurrent tasks on a single-core MCU with scheduling, queues, and timing guarantees.

---

## Exercises

### Easy
1. A smoke detector needs to read a sensor every second and beep if smoke is detected. Would you use a Raspberry Pi 4 or an STM32L010 (Cortex-M0+)? Justify your answer.
2. What is the difference between Flash memory and SRAM in a microcontroller? What is stored in each?
3. Your MCU needs to communicate with both an OLED display and a microSD card simultaneously. Which peripherals would you use for each?

### Medium
4. Calculate power consumption for a home temperature logger: STM32L010 Cortex-M0+. Active current: 4.5 mA at 3.3V, 32MHz. Stop mode: 0.9µA. Wake up every 60 seconds, take reading (5ms active), store to flash (10ms active), go back to sleep. Battery: 2xAA = 3000mAh. (a) Average current draw? (b) Battery life?
5. DMA vs interrupt-driven UART: An STM32 receives GPS data at 9600 baud (960 bytes/second). Each byte triggers an interrupt with 50-cycle overhead. The CPU runs at 72 MHz. (a) What fraction of CPU time is spent on UART interrupts? (b) With DMA and 32-byte buffer: how often is the CPU interrupted? (c) What fraction of CPU time with DMA?
6. Design an embedded system for a smart plant watering device: reads soil moisture sensor (ADC), controls water pump (GPIO), displays moisture level on OLED (I2C), connects to Wi-Fi (ESP32 module via UART), runs on 5V/2A adapter. Specify: MCU family, Flash/SRAM requirements, peripheral connections, firmware architecture (bare-metal or RTOS?).

### Hard
7. Bootloader design for firmware update: An STM32 device in the field needs to receive firmware updates over UART. Design a bootloader: (a) bootloader lives in the first 16KB of Flash, application in remaining 112KB — how does the bootloader know to jump to the application? (b) how does the bootloader receive the new firmware (UART XModem/YMODEM protocol)? (c) how does it verify integrity (CRC32) and prevent bricking if power fails mid-flash? (d) how does the application trigger a "goto bootloader" event?
8. Real-time control system timing analysis: A brushless DC motor controller runs on STM32F4 at 168 MHz. The current control loop must run every 50µs (20kHz). Each iteration: ADC sampling (1µs), current calculation (2µs), PWM update (0.5µs). (a) What fraction of CPU time does the control loop use? (b) A position control loop runs at 1kHz (every 1ms). Can both loops run safely if the current loop has higher priority? (c) A UART debug interface also needs to run. At 115200 baud, what is the maximum safe buffer size to avoid missing control loop deadlines?
