# Chapter 61: Testing and Quality

You cannot test quality into a chip — you must build it in. But manufacturing is imperfect: crystal defects, process variations, particle contamination, and lithography imperfections create chips that don't work or work incorrectly. Testing finds these failures before they reach customers. This chapter covers how chips are tested (at wafer level and package level), how manufacturing defects are modeled (fault models), how Design for Test (DFT) techniques make chips testable, and how reliability and quality are maintained over a product's lifetime.

## Table of Contents

1. [Why Test? Manufacturing Defects](#1-why-test-manufacturing-defects)
2. [Fault Models](#2-fault-models)
3. [Design for Test (DFT)](#3-design-for-test-dft)
4. [Wafer Test and Final Test](#4-wafer-test-and-final-test)
5. [JTAG — The Universal Debug Interface](#5-jtag--the-universal-debug-interface)
6. [Reliability: HTOL, ESD, Electromigration](#6-reliability-htol-esd-electromigration)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Why Test? Manufacturing Defects

A modern chip fabrication process has hundreds of steps. Each step can introduce defects:

**Sources of manufacturing defects:**
- **Particle contamination**: a dust particle on the wafer during lithography creates a short or open
- **Process variation**: resist thickness variation → printed feature too wide/too narrow
- **Oxidation non-uniformity**: gate oxide too thin → gate leakage; too thick → transistor won't turn on
- **Etch non-uniformity**: over-etch → wire breaks; under-etch → wire shorts to neighbor
- **Ion implantation dose error**: too few dopants → wrong threshold voltage

```
Defect type examples:
  
  Stuck-at-0:   Wire is stuck at logic 0 regardless of driver
                (caused by short to GND, oxide breakdown)
  
  Stuck-at-1:   Wire is stuck at logic 1
                (caused by short to VDD)
  
  Bridging:     Two adjacent wires are shorted together
                (particle in gap, under-etched metal)
  
  Open:         Wire is broken, no connection
                (over-etch, missing via)
  
  Timing fault: Wire has extra capacitance → too slow to meet timing
                (extra metal deposit, neighboring wire coupling)
```

**DPPM (Defective Parts Per Million)**: The industry quality metric. Consumer electronics targets <10 DPPM. Automotive (AEC-Q100) targets <1 DPPM. A 1 DPPM failure rate means 1 in every million shipped chips has a field failure. Testing is what achieves this.

### Quick Check
> 1. Name three sources of manufacturing defects in chip fabrication.
> 2. What is a "stuck-at fault" and what manufacturing defect causes it?
> 3. What does DPPM mean and what is the automotive quality target?

---

## 2. Fault Models

To test for defects, we need models of how defects affect circuit behavior:

**Stuck-at fault model (most common)**: Assumes any wire can be stuck at 0 or stuck at 1. A circuit with N wires has 2N possible stuck-at faults. Testing must apply inputs that:
1. **Sensitize** the fault: make the expected value on the faulty wire different from the stuck value
2. **Propagate** the effect: route the fault effect to an observable output

```
Fault coverage calculation:
  
  Simple 2-input AND gate: a & b → z
  
  Faults to test:
  - a stuck-at-0: apply a=1, b=1 → expect z=1 (faulty: z=0) ✓
  - a stuck-at-1: apply a=0, b=1 → expect z=0 (faulty: z=1) ✓
  - b stuck-at-0: apply a=1, b=1 → expect z=1 ✓
  - b stuck-at-1: apply a=1, b=0 → expect z=0 ✓
  - z stuck-at-0: apply a=1, b=1 → observe z output ✓
  - z stuck-at-1: apply a=0, b=1 → observe z output ✓
  
  Fault coverage = (faults tested) / (total faults) × 100%
  Goal: >99% fault coverage
```

**Transition fault model**: A wire transitions 0→1 or 1→0 too slowly (timing fault from extra capacitance). Requires tests that exercise the timing of each transition, not just the logic level.

**Bridging fault model**: Two adjacent wires short together. More complex — requires testing that detects when wire A forces wire B to an incorrect state.

**Cell-Aware Test (CAT)**: Model defects at the transistor level inside standard cells (intra-cell bridging, gate oxide shorts). Used for automotive and high-reliability parts.

### Quick Check
> 1. What are the two steps required to test for a stuck-at fault?
> 2. What is "fault coverage" and what percentage is typically targeted?
> 3. What is a transition fault and what manufacturing defect causes it?

---

## 3. Design for Test (DFT)

**DFT** inserts extra hardware into the chip design specifically to enable testing. Without DFT, many internal nodes are unreachable from the chip's output pins — you cannot observe whether they are working correctly.

**Scan chains**: The most important DFT technique. Every flip-flop in the design is made into a **scan flip-flop** — it has two inputs: normal D input (functional mode) and SI (scan input, test mode). Scan flip-flops are connected in a long chain:

```
Scan chain concept:
  
  Functional mode:       Scan test mode:
  
  D1 → FF1 → Q1         SI → FF1 → FF2 → FF3 → ... → FFn → SO
  D2 → FF2 → Q2         (scan chain: shift register)
  D3 → FF3 → Q3
  
  In test mode:
  1. Shift in test pattern through SI (n clock cycles for n FFs)
  2. Apply one clock cycle in functional mode
  3. Shift out captured state through SO (n clock cycles)
  4. Compare captured output to expected output → fault detected?
```

A chip might have 10 million flip-flops divided into 1000 scan chains of 10,000 FFs each. All scan chains run in parallel → reduces test time from 10M cycles to 10K cycles.

**ATPG (Automatic Test Pattern Generation)**: Software tool (Mentor TestKompress, Synopsys TetraMAX) that automatically generates scan test patterns to achieve target fault coverage. ATPG produces a set of vectors (test patterns) that, when shifted through the scan chains, test for stuck-at, transition, and bridging faults.

**BIST (Built-In Self Test)**: Add hardware to the chip that generates test patterns and checks responses internally — no external test equipment needed. Used for SRAMs (memory BIST), PLLs, and SerDes. The chip tests itself.

```
Memory BIST (MBIST) concept:
  
  ┌────────────────────────────────┐
  │        SRAM array              │
  │  ┌─────────────────────────┐  │
  │  │   MBIST controller      │  │
  │  │  1. March algorithm:    │  │
  │  │     Write 0 to all cells│  │
  │  │     Read all (expect 0) │  │
  │  │     Write 1 to all cells│  │
  │  │     Read all (expect 1) │  │
  │  │  2. Pass/Fail output    │  │
  │  └─────────────────────────┘  │
  └────────────────────────────────┘
  MBIST runs from ROM at power-on, requires no external tester
```

**Boundary scan (JTAG IEEE 1149.1)**: Special shift registers at every I/O pin allow testing the chip's I/O connections without probing individual pins. Also enables board-level test (test solder joints on a PCB).

### Quick Check
> 1. What is a scan chain and how does it enable testing of internal flip-flops?
> 2. What does ATPG do?
> 3. What is MBIST and when does it typically run?

---

## 4. Wafer Test and Final Test

**Wafer probe test:**
Before dicing, every die on the wafer is tested using a **probe card** — a device with thousands of tiny pins (probes) that contact the die's test pads.

The **ATE (Automatic Test Equipment)** applies test vectors and compares responses. Failing dies are marked with ink or recorded in a map.

```
Wafer probe setup:
  
  Prober: mechanical stage moves wafer under probe card
  Probe card: 5000–15000 probe needles, custom for each design
  ATE: Teradyne J750, Advantest T2000 (~$5M each)
  
  Test sequence:
  1. SCAN test: shift test vectors through scan chains (90% of test time)
  2. BIST: run MBIST for SRAMs
  3. At-speed test: apply test at functional clock frequency
  4. Burn-in: elevate voltage/temperature briefly to screen infant mortality
  5. Record pass/fail for each die
  
  Duration: 0.5–5 seconds per die (for a consumer chip)
  Throughput: 500-1000 wafers/day per ATE machine
```

**Final test (post-package):**
After packaging, chips are tested again in sockets on load boards:
- **OS (Outgoing Specification) test**: verify chip meets datasheet spec over voltage/temperature
- **Parametric test**: measure leakage current, supply voltage vs frequency, I/O characteristics
- **Functional test**: run full software boot / application test

**Binning**: Not all chips are identical. Process variation causes some chips to run faster (at lower voltage) than spec, others to run slower (or only at higher voltage). **Binning** sorts chips into "bins":
- Bin 1: meets high-performance spec (sold as i9)
- Bin 2: meets mainstream spec (sold as i7)
- Bin 3: meets budget spec (sold as i5)
- Bin 4: only some cores working (sold as i3 with disabled cores)
- Fail: scrapped or sold for non-critical apps

This is why the same die is sold at different price points — Intel, AMD, and NVIDIA all use binning.

### Quick Check
> 1. What is a probe card and how does wafer probe testing work?
> 2. What is "binning" and why is it economically important?
> 3. Why do chips pass wafer test but sometimes fail final test?

---

## 5. JTAG — The Universal Debug Interface

**JTAG (Joint Test Action Group, IEEE 1149.1)** is the standard chip debug and test interface, present on virtually every serious chip. It uses 4 signals:

- **TCK**: Test Clock
- **TMS**: Test Mode Select (state machine control)
- **TDI**: Test Data In
- **TDO**: Test Data Out

JTAG implements a state machine (TAP — Test Access Port) that allows:
1. **Boundary scan**: control/observe every I/O pin via shift registers
2. **Internal scan**: access internal scan chains for manufacturing test
3. **Debug access**: connect to an on-chip debugger (see CPU registers, halt execution, set breakpoints)

```
JTAG TAP state machine:
  
  Test-Logic-Reset
       │
  Run-Test/Idle
       │
  ┌────┴────┐
  DR-Scan  IR-Scan   (DR=data register, IR=instruction register)
  │         │
  Capture  Capture
  │         │
  Shift    Shift
  │         │
  Exit1    Exit1
  │         │
  Update   Update
  
  TCK drives all state transitions
  TMS controls which branch to take
```

**ARM CoreSight** (ARM's debug architecture): Provides JTAG + SWD (Serial Wire Debug, 2-pin) access to ARM cores. Allows debugger (GDB + OpenOCD or LLDB) to halt the CPU, read registers, access memory, set hardware breakpoints.

**RISC-V External Debug Spec**: RISC-V equivalent — defines debug registers accessible via JTAG. Used in SiFive and all RISC-V implementations including SHAKTI.

**Security implications**: JTAG is powerful — it can read memory (including secrets), modify execution, and bypass security measures. Production chips often implement JTAG fusing: after manufacturing test, the JTAG interface is disabled or locked with a password stored in OTP (One-Time Programmable) memory.

### Quick Check
> 1. What are the four JTAG signals and what does each one do?
> 2. What is boundary scan and what does it enable?
> 3. Why do production chips often disable or lock JTAG?

---

## 6. Reliability: HTOL, ESD, Electromigration

A chip must not only work initially — it must work for the product's lifetime (consumer: 3–5 years; automotive: 15–20 years; aerospace: 30+ years).

**HTOL (High Temperature Operating Life) / Burn-in:**
Accelerated aging by running chips at elevated temperature and voltage for 168–1000 hours. Failures from infant mortality (early defects) are screened out. Uses Arrhenius equation to predict real-world lifetime from accelerated conditions.

```
Arrhenius acceleration:
  AF = e^(Ea/k × (1/T_use - 1/T_test))
  
  Where:
    Ea = activation energy (~0.7 eV for typical failure mechanisms)
    k = Boltzmann constant = 8.62×10⁻⁵ eV/K
    T_use = use temperature (e.g., 85°C = 358K)
    T_test = test temperature (e.g., 125°C = 398K)
  
  AF ≈ 30× (testing 168h at 125°C ≈ 5000h at 85°C)
```

**ESD (Electrostatic Discharge):**
A person walking across carpet can accumulate 2,000–5,000V of static charge. Touching a chip pin discharges this energy in nanoseconds — enough to destroy transistor gate oxide (1nm SiO₂ breaks down at <10V). ESD protection circuits (diodes to VDD/GND, RC circuits) at every I/O pin clamp the discharge. ESD is tested with HBM (Human Body Model) and CDM (Charged Device Model) tests.

**Electromigration (EM):**
At high current density, metal atoms in copper wires are pushed by electron flow. Over time, voids (wire thinning → open) and hillocks (wire thickening → short) develop. EM is the primary long-term reliability failure mechanism in chips.

```
Electromigration:
  
  Electrons →→→→→→→→ current flow →→→→→→→→
  Cu atoms ←←← pushed by electron wind ←←←
  
  After years:
  Void forms:       Hillock forms:
  ──── ────         ────╔╗────
       gap               bump
  (open circuit)    (may short to adjacent wire)
```

Black's Law: Mean Time to Failure (MTTF) = A × j⁻ⁿ × e^(Ea/kT), where j = current density, n ≈ 2. Keeping current density below limits (typically 1–10 MA/cm²) prevents EM failures.

**Hot electrons and NBTI:**
- **Hot carrier injection (HCI)**: High-energy electrons in the channel get injected into the gate oxide, permanently changing V_th over time → transistor degrades
- **NBTI (Negative Bias Temperature Instability)**: pMOS transistors degrade when held at negative gate voltage over time → V_th shifts

Both cause gradual performance degradation over chip lifetime. Circuit designers must include timing margin to account for this aging.

### Quick Check
> 1. What is HTOL and how does it screen for early-life failures?
> 2. What is electromigration and what is its long-term effect on chips?
> 3. What is ESD and how do ESD protection circuits work?

---

## Summary

- **Manufacturing defects** (particle contamination, process variation, etch/implant errors) create stuck-at, open, bridge, and timing faults.
- **Fault models**: stuck-at (most common), transition (timing faults), bridging, cell-aware.
- **DFT**: Scan chains (flip-flop shift registers for parallel test), ATPG (automatic test vector generation), MBIST (self-test for memories), boundary scan (I/O test via JTAG).
- **Test flow**: ATPG vector generation → wafer probe (ATE + probe card) → binning (performance sorting) → final test.
- **JTAG**: standard 4-pin debug/test interface (TCK, TMS, TDI, TDO). Boundary scan + internal scan + debug access. Disabled in production for security.
- **Reliability**: HTOL/burn-in (screen early failures), ESD protection (guard every pin), EM analysis (current density limits), aging compensation (HCI, NBTI timing margins).

---

## Exercises

### Easy
1. What is a stuck-at fault and what manufacturing defect causes it?
2. What is a scan chain and how does it reduce test time?
3. What does "binning" mean in chip manufacturing?

### Medium
4. Fault coverage math: A combinational circuit has 500 total stuck-at faults. After ATPG generates 200 test patterns, 480 faults are detected. (a) Fault coverage? (b) Remaining undetected faults? (c) If manufacturing defect rate is 0.01 defects per die, what fraction of defects escape detection? (d) For automotive 1 DPPM target with 1 million chips shipped: how many total defects at 0.01 defect rate? How many escapes at your coverage? Is this acceptable?
5. Scan chain test time: A chip has 2 million flip-flops divided into 200 scan chains of 10,000 FFs each. The test has 10,000 test patterns. Test clock = 100 MHz. (a) Time to shift in one pattern (10,000 shift cycles)? (b) Time to apply functional capture (1 cycle)? (c) Time to shift out (10,000 cycles)? (d) Total time per pattern? (e) Total test time for 10,000 patterns? (f) If test clock is raised to 200 MHz: new total time? (g) Compare to testing without scan (need 10,000 patterns applied functionally at 1 GHz functional clock): how much faster is scan?
6. HTOL acceleration: A chip needs a 10-year reliability guarantee at 70°C ambient (Tj = 90°C). HTOL test runs at 125°C for 168 hours. Activation energy Ea = 0.7 eV. (a) Calculate acceleration factor AF using Arrhenius equation. (b) Equivalent real-world time represented by 168h HTOL? (c) 168h × AF > 10 years? If not: how long must HTOL run to represent 10 years? (d) If 1000 chips are tested and 5 fail in HTOL: what is the estimated failure rate per million chip-hours at use conditions?

### Hard
7. DFT area overhead: A CPU core has 1 million flip-flops. Implementing scan: each regular DFF is replaced by a scan DFF (add a 2:1 mux before the D input). In TSMC 28nm: regular DFF = 6µm², scan DFF = 8µm². (a) Area overhead from scan flip-flops alone? (b) 200 scan chains need scan-in (SI) and scan-out (SO) pins, adding 400 extra I/O connections routed to the chain endpoints. Estimating 50µm² routing overhead per connection: routing overhead? (c) MBIST controller for 4 SRAM arrays (32KB each): ~5000 gates × 8µm² per gate. MBIST area? (d) Total DFT overhead as % of core area (assume core area = 4mm²). (e) Typical DFT overhead is 10–15%. Is your estimate in this range? What optimizations can reduce it?
8. Post-silicon debug scenario: Your chip boots but crashes randomly under load. The crash appears to be a random bit flip in SRAM during high-frequency operation. (a) Design a MBIST sequence to determine if the SRAM has systematic (always-failing) bits vs random (intermittent) failures. (b) The crash correlates with temperature > 80°C. Describe how you would use JTAG to: set a watchpoint on the SRAM address, halt the CPU on the write, read back the written data, and verify correctness. (c) Failure analysis suspects NBTI degradation on a pMOS transistor in the SRAM cell. How would you verify this with lab equipment? (hint: FIB-SEM, parametric measurement). (d) A metal ECO can add extra decoupling capacitors near the SRAM to stabilize supply voltage during high-frequency switching. How would this help if the root cause is supply noise rather than NBTI?
