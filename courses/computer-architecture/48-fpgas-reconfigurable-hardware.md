# Chapter 48: FPGAs — Reconfigurable Hardware

An FPGA (Field-Programmable Gate Array) is a chip whose logic can be reconfigured after manufacturing. Unlike a CPU (fixed logic, flexible via software) or an ASIC (fixed logic, fixed function), an FPGA lets you define custom digital circuits — by programming the chip. You can implement a CPU, a cryptography engine, a neural network accelerator, or a custom image processing pipeline directly in hardware, and change it whenever your requirements change. FPGAs bridge the gap between software flexibility and hardware performance. They are used in high-frequency trading, aerospace, telecommunications, data center acceleration, and rapid hardware prototyping.

## Table of Contents

1. [What Is an FPGA?](#1-what-is-an-fpga)
2. [FPGA Architecture: LUTs, Flip-Flops, DSP Blocks](#2-fpga-architecture-luts-flip-flops-dsp-blocks)
3. [FPGA Programming: HDL and HLS](#3-fpga-programming-hdl-and-hls)
4. [FPGA in Practice](#4-fpga-in-practice)
5. [Xilinx vs Intel FPGA](#5-xilinx-vs-intel-fpga)
6. [When to Use an FPGA vs CPU vs ASIC](#6-when-to-use-an-fpga-vs-cpu-vs-asic)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What Is an FPGA?

An FPGA contains a large array of configurable logic blocks connected by a programmable routing network. By "programming" the FPGA (writing a bitstream to configuration memory), you configure what logic each block performs and how they connect.

Think of it as a blank silicon canvas: you tell each logic cell what function to implement and which other cells to connect to, and the FPGA becomes your custom digital circuit.

```
FPGA concept:
  ┌────────────────────────────────────────────────────────┐
  │                     FPGA                                │
  │  CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB      │
  │  CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB      │
  │  CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB CLB      │
  │    (CLB = Configurable Logic Block)                     │
  │                                                         │
  │  Between CLBs: programmable routing network            │
  │  (you connect any CLB output to any CLB input)          │
  │                                                         │
  │  On edges: I/O pads (GPIO, differential, LVDS, etc.)   │
  │                                                         │
  │  Special blocks: RAM (BRAM), DSP slices, PCIe, Ethernet │
  └────────────────────────────────────────────────────────┘
  
  Program the FPGA by loading a "bitstream" file that says:
  "CLB at row 5, col 3: implement XOR. Connect its output to CLB at row 5, col 4."
  (The bitstream encodes millions of such configuration bits)
```

### Quick Check
> 1. What is the key difference between an FPGA and an ASIC?
> 2. What is a "bitstream" in FPGA terminology?
> 3. What is a CLB and what does it implement?

---

## 2. FPGA Architecture: LUTs, Flip-Flops, DSP Blocks

**LUT (Look-Up Table)**: The basic unit of combinational logic in an FPGA. A 6-input LUT (6-LUT) can implement any Boolean function of 6 inputs. Internally, it is a 64-bit SRAM: the 6 inputs are the address, and the stored bit at that address is the output.

```
6-input LUT implementing: F = A AND B AND NOT(C)
  
  Truth table stored in 64-bit SRAM:
    A B C D E F (inputs) → output
    0 0 0 0 0 0 → 0
    0 0 0 0 0 1 → 0
    ...
    1 1 0 0 0 0 → 1   (A=1, B=1, C=0 → true)
    1 1 1 0 0 0 → 0   (A=1, B=1, C=1 → false, NOT(C) fails)
    ...
  
  Reconfigure: change the 64-bit SRAM contents → different logic function
```

**Flip-flop (FF)**: One D flip-flop per LUT (usually). The flip-flop stores one bit of state — essential for sequential logic (registers, counters, state machines).

**Slice/CLB**: A grouping of 4-8 LUTs + flip-flops + carry logic. The exact grouping varies by vendor.

**Block RAM (BRAM)**: On-chip RAM blocks of fixed size (18Kb or 36Kb typical). Used to implement local memories, FIFOs, frame buffers.

**DSP slices**: Hard-coded multiplier + accumulator blocks (18×18 or 27×18 multipliers). Much more efficient than implementing a multiplier from LUTs. Critical for signal processing and neural networks.

**PCIe, Ethernet, SerDes**: Modern FPGAs include hard-IP blocks for common interfaces, avoiding the need to implement these in LUTs.

```
Xilinx Ultrascale+ (XCVU9P) — a large FPGA:
  LUTs:     1,182,240 × 6-input LUTs
  FFs:      2,364,480 flip-flops
  BRAMs:    4,320 × 36Kb = 18 Mb total
  DSP slices: 6,840
  PCIe Gen4 blocks: 4
  100G Ethernet: 12
  
  This is a large FPGA, used in data centers and aerospace
```

### Quick Check
> 1. What is a 6-input LUT? How does it implement an arbitrary Boolean function?
> 2. Why are DSP slices important in FPGAs?
> 3. An FPGA needs to implement a 1024×1024 16-bit multiply. Would you use LUTs or DSP slices? Why?

---

## 3. FPGA Programming: HDL and HLS

FPGAs are programmed using **Hardware Description Languages (HDLs)** — languages that describe digital circuits, not algorithms. This is fundamentally different from software programming.

**Verilog** (industry standard):
```verilog
// 4-bit counter in Verilog
module counter(
    input  clk,
    input  reset,
    output reg [3:0] count  // 4-bit register
);
    always @(posedge clk) begin
        if (reset)
            count <= 4'b0000;
        else
            count <= count + 1;
    end
endmodule
```

This describes **hardware structure**: a register called `count`, a clock edge trigger, and logic that increments it. Not "run this code" — "build this circuit."

**VHDL** (alternative HDL, verbose but strict):
```vhdl
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;
entity counter is
    Port ( clk : in STD_LOGIC; reset : in STD_LOGIC; count : out STD_LOGIC_VECTOR(3 downto 0));
end counter;
architecture Behavioral of counter is
    signal count_int : STD_LOGIC_VECTOR(3 downto 0) := "0000";
begin
    count <= count_int;
    process(clk)
    begin
        if rising_edge(clk) then
            if reset = '1' then count_int <= "0000";
            else count_int <= count_int + 1;
            end if;
        end if;
    end process;
end Behavioral;
```

**HLS (High-Level Synthesis)**: Write C/C++ or OpenCL code; a tool (Xilinx Vitis HLS, Intel HLS Compiler) generates Verilog/VHDL automatically. Much more productive but less hardware control. Good for:
- Simple data processing pipelines
- Floating-point arithmetic
- Prototype implementation

```c
// HLS C code for FIR filter — tool generates hardware
#pragma HLS pipeline II=1   // pipeline with II=1 (one output per cycle)
void fir_filter(hls::stream<int16_t>& in, hls::stream<int16_t>& out, int16_t* coeff) {
    static int16_t shift_reg[TAPS];
    // ... (sliding window + MAC logic)
}
```

**FPGA design flow:**
1. Write HDL/HLS code
2. Synthesis: convert HDL → netlist (logical gates)
3. Implementation: place-and-route (map netlist to physical FPGA resources)
4. Bitstream generation: create programming file
5. Program FPGA: load bitstream via JTAG or flash memory

Implementation (place-and-route) is computationally intensive — a large FPGA can take 10–24 hours to compile.

### Quick Check
> 1. What is the fundamental difference between Verilog and C++ in terms of what they describe?
> 2. What is HLS and what does it sacrifice vs hand-written Verilog?
> 3. Why does FPGA implementation (place-and-route) take hours when CPU compilation takes seconds?

---

## 4. FPGA in Practice

**Prototyping**: FPGA is the standard tool for prototyping ASICs before tape-out. The ASIC design is first run on an FPGA to verify correctness. Apple, NVIDIA, Intel all use large FPGA farms for pre-silicon verification.

**High-frequency trading (HFT)**: Market data arrives as network packets; trading decisions must be made in microseconds. An FPGA receiving data directly from the network (FPGA NIC) can process and respond in 100–500 ns vs 1–10 µs for a CPU. FPGA HFT systems eliminate CPU OS overhead entirely.

**Telecommunications**: Base stations use FPGAs for the physical layer processing (FFT, OFDM, channel coding) — high volume, custom algorithms, needs to evolve with new standards.

**Data center acceleration**: Amazon AWS F1 instances: FPGA in EC2 accessible to developers. Microsoft Project Catapult: FPGAs in Azure data centers for Bing search ranking, network virtualization. These FPGAs are connected to CPUs via PCIe.

**Space and defense**: FPGAs are radiation-tolerant (special "rad-hard" FPGAs) and field-reprogrammable. NASA satellites use FPGAs. The Mars Curiosity rover used a Xilinx Virtex FPGA.

**Custom accelerators**: Implement a neural network inference pipeline directly in FPGA fabric — fixed latency, no OS jitter, tailored bit-widths (INT4, INT8).

### Quick Check
> 1. Why do HFT (high-frequency trading) firms use FPGAs instead of CPUs?
> 2. What is Amazon AWS F1 instance?
> 3. What makes FPGAs suitable for space applications?

---

## 5. Xilinx vs Intel FPGA

Two companies dominate the FPGA market:

**AMD-Xilinx** (AMD acquired Xilinx in 2022):
- Artix, Kintex, Virtex, Zynq, Versal product families
- **Zynq**: ARM Cortex-A9/A53 + FPGA fabric on one chip (PS + PL)
- **Versal AI Core**: FPGA + AI Engine (VLIW DSP array for ML) + ARM Cortex-A72
- **Alveo accelerator cards**: PCIe FPGA cards for data center
- Tools: Vivado (HDL), Vitis (HLS, AI), Vivado IP Integrator (block diagram)

**Intel (acquired Altera in 2015)**:
- Cyclone (low cost), Arria (mid-range), Stratix (high-end), Agilex (latest)
- **Agilex**: Intel 10nm process, HBM2e memory on-package option, PCIe Gen 5
- **PAC (Programmable Acceleration Card)**: Intel FPGA accelerator cards
- Tools: Quartus Prime (HDL), Intel HLS Compiler

**Market comparison (2024):**
- AMD-Xilinx: ~60% market share
- Intel FPGA: ~35% market share
- Other (Microchip Microsemi, Lattice): ~5%

**Lattice Semiconductor**: specializes in small, low-power FPGAs for edge/IoT/security applications. ECP5, CrossLink-NX, iCE40. The iCE40 is popular in hobbyist and small device markets.

**Open-source FPGA toolchain**: For Lattice iCE40/ECP5: IceStorm (iCE40) and Project Trellis (ECP5) reverse-engineered the bitstream format — enabling fully open-source Verilog → bitstream compilation without vendor tools. This is the only open-source FPGA toolchain that works on real silicon.

### Quick Check
> 1. What is the Zynq and why is it useful?
> 2. Which FPGA vendor has a fully open-source toolchain, and for which chip families?
> 3. What is the Versal AI Core and what makes it different from a standard FPGA?

---

## 6. When to Use an FPGA vs CPU vs ASIC

```
Decision matrix:

                FPGA            CPU             GPU             ASIC
Power           Medium          Medium-High     High            Low
Latency         Very Low        Medium          High (setup)    Very Low
Throughput      Medium-High     Medium          Very High       Very High
Flexibility     High            Very High       High            None
Time to market  Weeks-months    Days            Days            12-24 months
NRE cost        Low             None            None            $1M-$100M+
Unit cost       High ($50-$500) Low ($5-$500)   High ($300+)    Low at volume
```

**Use FPGA when:**
- Low-latency hard real-time (HFT, control systems, telecommunications)
- Custom interface protocols (you need hardware that doesn't exist)
- Prototyping before ASIC (verify design before committing to tape-out)
- Low-to-medium volume (not worth ASIC NRE cost)
- Algorithm will change (avionics: mission updates via software)

**Use ASIC when:**
- High volume (millions+): ASIC unit cost beats FPGA at scale
- Maximum efficiency required (mobile phone chip, implantable medical device)
- Fixed function (the algorithm is settled and stable)
- Power budget is extreme (battery-powered, implants)

**Use CPU when:**
- Software flexibility is paramount
- Algorithm changes frequently
- Single-threaded performance matters
- Small data volumes

**Use GPU when:**
- Batch data-parallel computation (ML training/inference, scientific simulation)
- Latency tolerance is acceptable (GPU launch overhead ~5µs)

### Quick Check
> 1. Why does an FPGA have lower latency than a GPU for many tasks?
> 2. When does an ASIC beat an FPGA on cost?
> 3. What is "NRE cost" for an ASIC and why does it matter?

---

## Summary

- An **FPGA** is reconfigurable hardware — its logic and routing can be changed by loading a bitstream. Unlike ASICs (fixed) or CPUs (software-programmed), FPGAs implement custom digital circuits.
- **Architecture**: LUTs (any Boolean function of 6 inputs), flip-flops (registers), DSP slices (hard multipliers), BRAM (on-chip memory), hard-IP (PCIe, Ethernet).
- **Programming**: Verilog or VHDL (hardware description languages), or HLS (C++ to hardware synthesis). Design flow includes synthesis → place-and-route → bitstream generation.
- **Applications**: HFT (sub-microsecond trading), telecom (base station physical layer), data center acceleration (AWS F1, Azure Catapult), ASIC prototyping, aerospace.
- **Vendors**: AMD-Xilinx (60% share, Zynq/Versal), Intel FPGA (35%, Agilex), Lattice (5%, iCE40 with open-source tools).
- **Trade-offs vs ASIC**: FPGA costs more per unit at high volume, uses more power, but can be reprogrammed and has no NRE cost.

---

## Exercises

### Easy
1. What is a LUT and how does it implement any Boolean function?
2. Give three applications where FPGAs are preferred over CPUs.
3. What is NRE cost for an ASIC and why does it favor FPGAs at low production volume?

### Medium
4. An HFT system receives market data packets at 10 Gbps. CPU latency: 5 µs (DPDK optimized). FPGA latency: 200 ns. The system makes one trade per second per instrument, with 1000 instruments. (a) What is the latency advantage per trade in ns? (b) In competitive HFT, being 4.8 µs faster equals "front running" by how many price updates at 10 Gbps? (c) What is the capital cost difference between FPGA (FPGA card ~$20K) + FPGA engineering ($500K) vs CPU server ($10K)? At what profit-per-trade advantage does FPGA break even in 1 year?
5. FPGA DSP slice utilization: implement a 64-tap FIR filter in FPGA. Each tap requires one 18×18 multiplier and one adder. (a) How many DSP slices are needed? (b) A Xilinx Artix-7 XC7A100T has 240 DSP slices. How many 64-tap FIR filters can it support in parallel? (c) If the FPGA runs at 200 MHz and each output sample requires 64 MACs executed in a pipelined fashion (one new output per cycle with II=1), what is the maximum sample rate?
6. FPGA vs ASIC for an IoT sensor: you need to produce 10 million chips for a smart home sensor (temperature, humidity, display, BLE). Compare: FPGA (iCE40LP: $3 each, program in-factory) vs ASIC (design cost $2M, fabrication per unit $0.30). At what production volume does ASIC break even? What other factors (time-to-market, power, upgradability) matter?

### Hard
7. Implement a simple CPU in FPGA: Design a minimal 8-bit CPU in Verilog with: 4 registers, 8-bit ALU (add, sub, AND, OR), 256-byte program memory, 256-byte data memory, 8 instructions (LOAD, STORE, ADD, SUB, AND, OR, JUMP, HALT). Sketch the Verilog module hierarchy: (a) what modules do you need? (b) how does the FSM (Fetch/Decode/Execute states) control the datapath? (c) how many LUTs does your design use? (estimate based on components)
8. Data center FPGA acceleration: Microsoft uses FPGAs for Bing search ranking. The ranking algorithm scores 100,000 documents per query, each document requiring 1000 floating-point operations. Query latency target: 2ms. (a) Required throughput: GFLOPS? (b) Would a CPU, GPU, or FPGA be fastest? (c) FPGA implementation uses fixed-point arithmetic (INT16 instead of FP32): what precision loss is acceptable and how do you validate it? (d) When the ranking algorithm changes (new ML model), how do you update 20,000 FPGAs deployed globally?
