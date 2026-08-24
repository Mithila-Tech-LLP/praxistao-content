# Electronics: A Complete Learning Guide
### From Atomic Physics to Modern Processors, Sensors, and IoT Systems

---

## Table of Contents

### Part I — Foundations of Matter and Electricity

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [01](01-atoms-and-matter.md) | Atoms and Matter | Atomic structure, Bohr model, quantum numbers, electron orbitals, crystal lattice, bonding |
| [02](02-electrical-properties-and-band-theory.md) | Electrical Properties and Band Theory | Energy bands, valence/conduction bands, conductors, semiconductors, insulators, Fermi level |
| [03](03-semiconductors-deep-dive.md) | Semiconductors Deep Dive | Silicon crystal, intrinsic/extrinsic, n-type/p-type doping, carrier concentration, Hall effect, mobility |
| [04](04-pn-junction-and-diodes.md) | PN Junction and Diodes | Depletion region, built-in potential, Shockley equation, forward/reverse bias, rectifier, Zener, LED, Schottky |

---

### Part II — Discrete Components and Transistors

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [05](05-basic-electronic-components.md) | Basic Electronic Components | Resistors, capacitors, inductors, transformers, RC/RL/LC filters, resonance, impedance |
| [06](06-bjt-transistors.md) | BJT Transistors | NPN/PNP operation, α/β parameters, biasing, CE/CB/CC amplifier configs, switch design, h-parameters |
| [07](07-mosfet-and-advanced-transistors.md) | MOSFET and Advanced Transistors | Enhancement/depletion MOSFET, CMOS, FinFET, GAA (Gate-All-Around), IGBT, JFET, body effect |

---

### Part III — Digital Logic and Sequential Circuits

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [08](08-logic-gates-and-boolean-algebra.md) | Logic Gates and Boolean Algebra | Number systems, Boolean algebra, De Morgan's, K-maps, NAND/NOR universality, logic families (CMOS/TTL) |
| [09](09-combinational-circuits.md) | Combinational Circuits | Half/full adders, ripple-carry, carry-lookahead, MUX/DEMUX, encoders, decoders, ALU, barrel shifter, parity |
| [10](10-sequential-circuits-and-memory.md) | Sequential Circuits and Memory | Latches, flip-flops (SR/D/JK/T), registers, counters, FSM (Mealy/Moore), SRAM cell, DRAM refresh |

---

### Part IV — Semiconductor Manufacturing and Computer Architecture

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [11](11-semiconductor-fabrication-and-moores-law.md) | Semiconductor Fabrication and Moore's Law | Czochralski process, photolithography, EUV, ion implantation, CMP, TSMC/Intel/Samsung nodes, Moore's Law |
| [12](12-computer-architecture-fundamentals.md) | Computer Architecture Fundamentals | Von Neumann vs Harvard, pipelining, OOO execution, branch prediction, Tomasulo, cache hierarchy, virtual memory, TLB |
| [13](13-instruction-set-architectures.md) | Instruction Set Architectures | RISC vs CISC, ARM (A32/T32/A64), x86-64, RISC-V, MIPS, 8/16/32/64-bit evolution, calling conventions |

---

### Part V — Memory Systems

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [14](14-memory-systems.md) | Memory Systems | SRAM 6T cell, DRAM 1T1C, DDR1–DDR5 evolution, LPDDR, HBM (3D stacking + TSV), GDDR, NAND Flash (SLC/MLC/TLC/QLC, 3D V-NAND), NOR Flash (XIP), eMMC, UFS, NVMe (PCIe), ECC, cache hierarchy |

---

### Part VI — Microcontroller Architecture

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [15](15-microcontroller-architecture.md) | Microcontroller Architecture | MCU vs MPU vs SoC, clock system (HSI/HSE/PLL), GPIO modes, NVIC, SysTick, ADC (SAR), DAC, timers/PWM, UART, SPI, I2C, CAN, DMA, USB, power sleep modes |
| [16](16-arduino-and-avr-architecture.md) | Arduino and AVR Architecture | AVR Harvard architecture, ATmega328P (registers, pipeline, SRAM map), PROGMEM, all Arduino boards (Uno/Mega/Nano/Due/MKR), bootloader (Optiboot), ISR |
| [17](17-esp32-and-xtensa-architecture.md) | ESP32 and Xtensa Architecture | ESP8266/NodeMCU, Xtensa LX6 dual-core (register windows, FPU), WiFi/BLE stack, all ESP32 variants (S2/S3/C3/C6/H2), FreeRTOS SMP, deep sleep, ULP coprocessor, OTA |
| [18](18-stm32-and-arm-cortex-m.md) | STM32 and ARM Cortex-M | Cortex-M0 through M85, Thumb-2 ISA, register file, NVIC, MPU, all STM32 series, Blue Pill (F103C8T6), SWD debugging, CubeMX/HAL/LL, FreeRTOS |
| [19](19-raspberry-pi-and-arm-cortex-a.md) | Raspberry Pi and ARM Cortex-A | Cortex-A53/A72/A76 microarchitecture, AArch64 registers, all RPi models (1→5 + Zero), BCM2711 SoC, boot sequence, GPIO, pigpio, RP2040 + PIO state machines |

---

### Part VII — Modern Processor Architectures

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [20](20-apple-silicon-architecture.md) | Apple Silicon Architecture | M1–M4 chips, unified memory (PoP LPDDR), Firestorm P-cores (8-wide, 630-entry ROB), Icestorm E-cores, Neural Engine (11–38 TOPS), M1 Pro/Max/Ultra, Rosetta 2, Metal, Core ML |
| [21](21-intel-amd-x86-architecture.md) | Intel and AMD x86 Architecture | x86-64 registers, Intel Golden Cove (12 ports, 6-wide decode), Gracemont E-cores, Arrow Lake tiles, AMD Zen 1–5 (CCX→CCD), 3D V-Cache, Infinity Fabric, AVX-512, DDR5 timings |

---

### Part VIII — Sensors

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [22](22-sensors-fundamentals.md) | Sensor Fundamentals | Sensitivity, resolution, accuracy, precision, linearity, noise (SNR), drift, signal conditioning (TIA, Wheatstone bridge, 4–20 mA), MEMS fabrication, calibration techniques |
| [23](23-sensors-ir-ultrasonic-lidar.md) | IR, Ultrasonic, and LiDAR Sensors | TCRT5000, Sharp GP2Y (IR triangulation), HC-SR501 PIR, HC-SR04 ultrasonic, VL53L0X/VL53L1X (VCSEL+SPAD ToF), RPLiDAR, Velodyne, mmWave radar (HLK-LD2410) |
| [24](24-sensors-temperature-humidity-pressure.md) | Temperature, Humidity, and Pressure Sensors | Seebeck effect, thermocouple types (K/J/T/E), MAX31855, RTD (Callendar-Van Dusen), MAX31865, NTC/PTC thermistors, DHT11/DHT22, DS18B20 (1-Wire), LM35, BME280, SHT31 |
| [25](25-sensors-imu-accelerometer-gyroscope.md) | IMU: Accelerometer and Gyroscope | MEMS proof mass, Coriolis effect, MPU-6050 (DMP, FIFO, register map), ICM-42688-P, LSM6DSO, HMC5883L magnetometer, hard/soft iron calibration, complementary/Mahony/Madgwick/EKF filters, quaternions |
| [26](26-sensors-optical-chemical-others.md) | Optical, Chemical, and Other Sensors | LDR, photodiode (TIA), TCS34725 (RGBC color), optical encoder (quadrature), MQ gas sensors (SnO₂), CCS811/SCD40 (CO₂), pH electrode (Nernst equation), ACS712/INA219 (current), FSR, YF-S201 (flow), I2S MEMS mic |

---

### Part IX — Systems Integration

| Chapter | Title | Key Topics |
|---------|-------|------------|
| [27](27-putting-it-all-together.md) | Putting It All Together | Wired/wireless protocol comparison, MQTT (QoS, topics), PCB design (4-layer stack, trace width, ESD), LDO vs DC-DC, battery types, Li-Ion CC-CV charging, IoT architecture, Node-RED, Home Assistant, ESPHome, debugging (SWD, logic analyzer), complete knowledge map |

---

## Reading Paths

### Beginner — Start Here
```
01 → 02 → 03 → 04 → 05 → 06 → 08 → 09 → 10 → 16 (Arduino)
```

### Embedded Systems / Firmware Engineer
```
08 → 09 → 10 → 14 → 15 → 16 → 17 → 18 → 22 → 23–26 → 27
```

### Hardware / PCB Designer
```
04 → 05 → 06 → 07 → 11 → 14 → 15 → 22 → 27
```

### Computer Architecture / CPU Enthusiast
```
08 → 09 → 10 → 11 → 12 → 13 → 14 → 19 → 20 → 21
```

### IoT / Maker
```
16 → 17 → 22 → 23 → 24 → 25 → 26 → 27
```

### Complete Linear Read (Recommended)
```
01 → 02 → ... → 27  (follow chapter order)
```

---

## Quick Reference

| Topic | Chapter |
|-------|---------|
| How transistors work | 06, 07 |
| How CPUs work internally | 12, 13 |
| DDR5 / HBM memory | 14 |
| Arduino pinout and internals | 16 |
| ESP32 WiFi + BLE | 17 |
| STM32 + FreeRTOS | 18 |
| Raspberry Pi boot sequence | 19 |
| Apple M-series chips | 20 |
| AMD Zen architecture | 21 |
| Reading sensor datasheets | 22 |
| HC-SR04 ultrasonic wiring | 23 |
| DS18B20 / DHT22 | 24 |
| MPU-6050 + sensor fusion | 25 |
| Gas sensors / pH / current | 26 |
| PCB design basics | 27 |
| MQTT / IoT architecture | 27 |

---

*27 chapters · ~700 KB of content · No external sources needed*
