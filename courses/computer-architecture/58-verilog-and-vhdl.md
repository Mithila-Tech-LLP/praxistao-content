# Chapter 58: Hardware Description Languages — Verilog and VHDL

Hardware Description Languages (HDLs) are the programming languages of chip design. Just as software engineers write Python or C++ to describe algorithms, hardware engineers write Verilog or VHDL to describe digital circuits. But the analogy ends quickly: an HDL describes **hardware structure** — registers, wires, combinational logic, and their connections — not a sequence of instructions for a processor to execute. Everything in Verilog runs conceptually **in parallel**, not sequentially. This chapter teaches the core concepts of Verilog (the more widely used HDL), contrasts it with VHDL, and introduces SystemVerilog — the modern superset used for both design and verification.

## Table of Contents

1. [HDLs vs Programming Languages](#1-hdls-vs-programming-languages)
2. [Verilog Fundamentals](#2-verilog-fundamentals)
3. [Sequential Logic in Verilog](#3-sequential-logic-in-verilog)
4. [Combinational Logic in Verilog](#4-combinational-logic-in-verilog)
5. [SystemVerilog — The Modern Standard](#5-systemverilog--the-modern-standard)
6. [VHDL — The Alternative](#6-vhdl--the-alternative)
7. [Common Design Patterns](#7-common-design-patterns)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. HDLs vs Programming Languages

The most important concept in HDL design is that it describes **hardware, not software**.

```
C++ (software):                  Verilog (hardware):
  x = a + b;                       assign sum = a + b;
  y = c + d;     // runs in         assign diff = a - b;  // BOTH exist
  z = x + y;     // sequence        assign result = sum + diff;  // simultaneously
  
  x is computed first,             sum and diff are both computed
  then y, then z.                  simultaneously as combinational logic.
  Sequential execution.            Parallel hardware.
```

**Key mental model shift**: In Verilog, you are not writing code that runs — you are describing circuits that **exist**. An `assign` statement describes a wire whose value is continuously driven by a combinational function. An `always @(posedge clk)` block describes a register that captures a value on the clock edge.

**Hardware is inherently parallel**: All combinational logic evaluates simultaneously. All registers update simultaneously at the clock edge. You cannot "sequence" two `assign` statements.

### Quick Check
> 1. What is the key difference between a software program and an HDL description?
> 2. In Verilog, what does `assign out = a & b` describe?
> 3. If you write two `assign` statements for different wires, do they execute sequentially or in parallel?

---

## 2. Verilog Fundamentals

**Module**: The fundamental building block in Verilog, equivalent to a function or class. A module has ports (inputs and outputs), internal wires and registers, and logic statements.

```verilog
// A simple full adder
module full_adder (
    input  a,         // 1-bit input
    input  b,         // 1-bit input
    input  cin,       // carry in
    output sum,       // sum output
    output cout       // carry out
);
    assign sum  = a ^ b ^ cin;       // XOR
    assign cout = (a & b) | (b & cin) | (a & cin);  // carry logic
endmodule
```

**Data types:**
- `wire`: A combinational connection. Continuously driven by a source. Default state = high-Z (disconnected).
- `reg`: A storage element (not always a flip-flop — depends on context). Used inside `always` blocks.
- `logic` (SystemVerilog): Unified type that can be either wire or reg — preferred in modern code.

**Bus declarations:**
```verilog
wire [7:0]  data_bus;      // 8-bit wire (bits 7 down to 0)
reg  [31:0] program_counter; // 32-bit register
input [3:0] opcode;        // 4-bit input port
```

**Number literals:**
```verilog
4'b1010   // 4-bit binary: 1010 = decimal 10
8'hFF     // 8-bit hex: 0xFF = 255
16'd1000  // 16-bit decimal: 1000
'1        // all-ones (any width)
```

**Operators:**
```verilog
// Bitwise (per-bit)
&  AND:    a & b
|  OR:     a | b
^  XOR:    a ^ b
~  NOT:    ~a

// Reduction (collapses all bits to one)
& a    // AND all bits of a
| a    // OR all bits of a (any bit set?)
^ a    // XOR all bits (odd parity?)

// Arithmetic
a + b   a - b   a * b   a / b   a % b

// Comparison
a == b  a != b  a > b  a < b  a >= b  a <= b

// Shift
a << 2  (shift left 2 bits)
a >> 2  (logical right shift)
a >>> 2 (arithmetic right shift, preserves sign)

// Concatenation
{a, b}        // concatenate bits: a[3:0], b[3:0] → 8 bits
{4{a}}        // replicate: 4 copies of a
{a[7:4], b[3:0]}  // pick bits and combine
```

### Quick Check
> 1. What is the difference between `wire` and `reg` in Verilog?
> 2. What does `{a[7:4], b[3:0]}` produce if a=0xAB and b=0xCD?
> 3. Write a Verilog expression to check if any bit in an 8-bit bus is set.

---

## 3. Sequential Logic in Verilog

Sequential logic (flip-flops, registers, state machines) is described with `always @(posedge clk)` blocks.

**D flip-flop with synchronous reset:**
```verilog
module d_flipflop (
    input  clk, reset, d,
    output reg q
);
    always @(posedge clk) begin
        if (reset)
            q <= 1'b0;    // non-blocking assignment
        else
            q <= d;
    end
endmodule
```

**Non-blocking assignment (`<=`)**: Used inside `always @(posedge clk)`. All right-hand sides are evaluated first (using current values), then all assignments happen simultaneously. This models how real flip-flops work — they all capture their input in the same clock edge.

**Blocking assignment (`=`)**: Used inside combinational `always @(*)` blocks. Executes like sequential software code. Use non-blocking in clocked blocks, blocking in combinational blocks — mixing them incorrectly causes subtle simulation bugs.

**4-bit counter:**
```verilog
module counter_4bit (
    input  clk, reset,
    output reg [3:0] count
);
    always @(posedge clk) begin
        if (reset)
            count <= 4'b0;
        else
            count <= count + 1;
    end
endmodule
```

**Finite State Machine (FSM):**
```verilog
module traffic_light (
    input  clk, reset,
    output reg [1:0] state  // 0=RED, 1=GREEN, 2=YELLOW
);
    // State encoding
    localparam RED    = 2'b00;
    localparam GREEN  = 2'b01;
    localparam YELLOW = 2'b10;
    
    always @(posedge clk) begin
        if (reset)
            state <= RED;
        else begin
            case (state)
                RED:    state <= GREEN;
                GREEN:  state <= YELLOW;
                YELLOW: state <= RED;
                default: state <= RED;
            endcase
        end
    end
endmodule
```

### Quick Check
> 1. What is the difference between blocking (`=`) and non-blocking (`<=`) assignment?
> 2. In an FSM, why is `default:` case important?
> 3. Write a Verilog register that captures input `d` on positive clock edge and resets to 0 on active-high reset.

---

## 4. Combinational Logic in Verilog

Combinational logic (circuits with no memory — output depends only on current input) is described with:
- `assign` statements
- `always @(*)` blocks (or `always_comb` in SystemVerilog)

**Multiplexer (2:1 mux):**
```verilog
// Using assign (ternary operator)
assign out = sel ? a : b;

// Using always (case statement)
always @(*) begin
    case (sel)
        1'b0: out = b;
        1'b1: out = a;
    endcase
end
```

**4-bit ripple carry adder (structural):**
```verilog
module adder_4bit (
    input  [3:0] a, b,
    input  cin,
    output [3:0] sum,
    output cout
);
    wire c1, c2, c3;
    
    // Instantiate four full adders
    full_adder FA0 (.a(a[0]), .b(b[0]), .cin(cin), .sum(sum[0]), .cout(c1));
    full_adder FA1 (.a(a[1]), .b(b[1]), .cin(c1),  .sum(sum[1]), .cout(c2));
    full_adder FA2 (.a(a[2]), .b(b[2]), .cin(c2),  .sum(sum[2]), .cout(c3));
    full_adder FA3 (.a(a[3]), .b(b[3]), .cin(c3),  .sum(sum[3]), .cout(cout));
endmodule
```

**Priority encoder:**
```verilog
module priority_enc (
    input  [3:0] in,
    output reg [1:0] out,
    output reg valid
);
    always @(*) begin
        valid = 1'b1;
        casez (in)  // casez: treat 'z' as don't care
            4'b1???:  out = 2'b11;
            4'b01??:  out = 2'b10;
            4'b001?:  out = 2'b01;
            4'b0001:  out = 2'b00;
            default: begin out = 2'b00; valid = 1'b0; end
        endcase
    end
endmodule
```

**Synthesis implication**: A `case` statement in an `always @(*)` block synthesizes to a mux tree. An `if-else` chain synthesizes to a priority encoder (first condition takes priority). Understand what hardware your code implies.

### Quick Check
> 1. What is the difference between `always @(*)` and `always @(posedge clk)`?
> 2. What hardware does a Verilog `case` statement typically synthesize to?
> 3. What is the `casez` statement and when is it useful?

---

## 5. SystemVerilog — The Modern Standard

**SystemVerilog** (IEEE 1800-2012) extends Verilog with:
- Better type system (`logic` type replaces wire/reg confusion)
- Interfaces (group related signals together)
- Classes (for object-oriented testbenches)
- Assertions (SVA — SystemVerilog Assertions)
- Enhanced verification features (UVM support)

**`logic` type** (replaces confusing wire/reg distinction):
```systemverilog
logic clk, reset;
logic [7:0] data;
logic [31:0] bus_data;
```

**Interface** (bundle of signals with direction):
```systemverilog
interface axi_bus (input logic clk);
    logic [31:0] awaddr, wdata, araddr, rdata;
    logic        awvalid, awready, wvalid, wready;
    // ... more AXI signals
    
    modport master (output awaddr, awvalid, input awready, ...);
    modport slave  (input awaddr, awvalid, output awready, ...);
endinterface
```

**Assertion (SVA):**
```systemverilog
// Assertion: if req is high, grant must be high within 4 cycles
property req_grant;
    @(posedge clk) req |=> ##[1:4] grant;
endproperty
assert property(req_grant) else $error("Grant timeout!");
```

**Enum for FSMs** (readable state names):
```systemverilog
typedef enum logic [1:0] {
    IDLE = 2'b00,
    FETCH = 2'b01,
    EXECUTE = 2'b10,
    WRITEBACK = 2'b11
} state_t;

state_t current_state, next_state;
```

### Quick Check
> 1. What does the `logic` type in SystemVerilog replace and why is it an improvement?
> 2. What is an interface in SystemVerilog and what problem does it solve?
> 3. Write a SystemVerilog assertion that checks `ack` is never high for more than 5 cycles.

---

## 6. VHDL — The Alternative

**VHDL** (VHSIC Hardware Description Language, IEEE 1076) was developed by the US Department of Defense in 1983. It is the other major HDL.

**Key differences from Verilog:**
- Strongly typed: you cannot assign an 8-bit value to a 4-bit signal without an explicit cast
- More verbose: explicit type declarations, library imports, architecture bodies
- Favored in Europe, aerospace/defense industries
- `std_logic` (9-value logic: 0, 1, X, Z, W, L, H, -, U) vs Verilog's 4-value (0, 1, X, Z)

```vhdl
-- VHDL D flip-flop
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

entity d_flipflop is
    Port (
        clk   : in  STD_LOGIC;
        reset : in  STD_LOGIC;
        d     : in  STD_LOGIC;
        q     : out STD_LOGIC
    );
end d_flipflop;

architecture Behavioral of d_flipflop is
begin
    process (clk)
    begin
        if rising_edge(clk) then
            if reset = '1' then
                q <= '0';
            else
                q <= d;
            end if;
        end if;
    end process;
end Behavioral;
```

**VHDL 9-value logic** is more accurate for simulation — `X` (unknown) and `U` (uninitialized) detect bugs that Verilog might miss because Verilog treats uninitialized as 0.

**Industry split**: Most ASIC design uses Verilog or SystemVerilog. FPGA vendors (Xilinx tools) support both. Military/aerospace and European companies often use VHDL.

### Quick Check
> 1. What is the main advantage of VHDL's strong typing vs Verilog?
> 2. What is `std_logic` in VHDL and how does its 9-value system differ from Verilog's 4-value?
> 3. Which industry sectors still prefer VHDL over Verilog?

---

## 7. Common Design Patterns

**Pipeline register:**
```verilog
// 2-stage pipeline: compute in stage 1, register, compute in stage 2
module pipeline_2stage (
    input  clk,
    input  [7:0] a, b, c,
    output [8:0] result  // a + b + c
);
    // Stage 1: compute a+b
    wire  [8:0] sum_ab  = {1'b0, a} + {1'b0, b};
    reg   [8:0] sum_ab_reg;
    reg   [7:0] c_reg;
    
    always @(posedge clk) begin
        sum_ab_reg <= sum_ab;
        c_reg      <= c;       // pipeline c to stay aligned
    end
    
    // Stage 2: add registered sum to c
    assign result = sum_ab_reg + {1'b0, c_reg};
endmodule
```

**FIFO (First In, First Out buffer) — structural:**
```verilog
module simple_fifo #(
    parameter DEPTH = 8,
    parameter WIDTH = 8
) (
    input  clk, reset, push, pop,
    input  [WIDTH-1:0] din,
    output [WIDTH-1:0] dout,
    output full, empty
);
    reg [WIDTH-1:0] mem [DEPTH-1:0];
    reg [$clog2(DEPTH)-1:0] wr_ptr, rd_ptr;
    reg [$clog2(DEPTH):0] count;
    
    assign full  = (count == DEPTH);
    assign empty = (count == 0);
    assign dout  = mem[rd_ptr];
    
    always @(posedge clk) begin
        if (reset) begin
            wr_ptr <= 0; rd_ptr <= 0; count <= 0;
        end else begin
            if (push && !full)  begin mem[wr_ptr] <= din; wr_ptr <= wr_ptr + 1; count <= count + 1; end
            if (pop  && !empty) begin rd_ptr <= rd_ptr + 1; count <= count - 1; end
        end
    end
endmodule
```

**Parameterized design** (`#(parameter)`) is critical for reuse — write an N-bit adder once, instantiate it at 8-bit, 32-bit, 64-bit.

### Quick Check
> 1. In the pipeline example, why is `c_reg` needed alongside `sum_ab_reg`?
> 2. What does `$clog2(DEPTH)` compute and why is it used for the pointer width?
> 3. Write a parameterized N-bit inverter module using `parameter`.

---

## Summary

- **HDLs** describe hardware structure, not algorithms. Everything runs in parallel.
- **Verilog**: `wire` for combinational, `reg` for sequential; `assign` for continuous logic; `always @(posedge clk)` for flip-flops; `always @(*)` for combinational.
- **Non-blocking (`<=`) vs blocking (`=`)**: non-blocking models simultaneous register updates; blocking models sequential combinational evaluation.
- **SystemVerilog**: superset of Verilog. Adds `logic` type, interfaces, enums, assertions (SVA), and OOP features for testbenches.
- **VHDL**: more verbose, strongly typed, 9-value logic. Preferred in defense/aerospace and Europe.
- **Common patterns**: pipeline registers, FIFOs, FSMs, parameterized modules.

---

## Exercises

### Easy
1. Write a Verilog module `and_gate` with 2-bit inputs `a`, `b` and 1-bit output `y` where `y = a[0] & a[1] & b[0] & b[1]`.
2. What is the difference between a non-blocking assignment `<=` and a blocking assignment `=`? When should each be used?
3. Write a Verilog 8-bit D register (all 8 bits captured on rising clock edge, synchronous active-high reset to 0).

### Medium
4. 4-bit ALU: Write a Verilog module for a 4-bit ALU supporting: ADD (op=000), SUB (op=001), AND (op=010), OR (op=011), XOR (op=100), NOT A (op=101), SHL (op=110), SHR (op=111). Inputs: `[3:0] a, b; [2:0] op`. Outputs: `[3:0] result; zero` (1 if result=0). Use a `case` statement. What hardware does the synthesis tool produce from this code?
5. 4-stage pipeline: A computation needs to compute `(a * b) + (c * d)` where a,b,c,d are 8-bit values. (a) Draw a pipeline diagram: which operations happen in which stage? (b) Write SystemVerilog for a 3-stage pipeline version: stage 1 = both multiplies, stage 2 = register products, stage 3 = add. (c) What is the latency (in cycles) of the pipeline? (d) What is the throughput (results per clock) once the pipeline is full?
6. UART transmitter FSM: A UART (serial communication) transmitter sends: 1 start bit (0), 8 data bits (LSB first), 1 stop bit (1). It asserts `tx_out` based on the current state. The `baud_clk` fires at the baud rate. (a) Define the states: IDLE, START, DATA[0..7], STOP. (b) Write a SystemVerilog FSM with a shift register holding the data byte. (c) Add an `enum` for state names. (d) Add an assertion: "When state=STOP, tx_out must be 1."

### Hard
7. Cache line array in Verilog: Implement a 4-way set-associative cache with 8 sets (index bits), 64-byte cache lines, and 32-bit addresses. (a) Calculate: offset bits (6), index bits (3), tag bits (remaining). (b) Declare Verilog memory arrays for: tag[8 sets][4 ways], valid[8 sets][4 ways], data[8 sets][4 ways][16 words]. (c) Write the lookup logic: given address, check all 4 ways for a tag match with valid bit — output hit/miss and data. (d) Write the LRU replacement logic using a 2-bit counter per set to track least-recently-used way.
8. AXI-Lite bus interface: AXI-Lite is a simplified version of the ARM AXI4 bus used in SoC design. Write a SystemVerilog module for an AXI-Lite slave register file with 4×32-bit registers (addr 0x00–0x0C). (a) Declare AXI-Lite interface signals: AWVALID/AWREADY/AWADDR (write address), WVALID/WREADY/WDATA (write data), BVALID/BREADY/BRESP (write response), ARVALID/ARREADY/ARADDR (read address), RVALID/RREADY/RDATA/RRESP (read data). (b) Implement write transaction FSM (IDLE → WADDR → WDATA → WRESP). (c) Implement read transaction FSM. (d) What is the key advantage of AXI's handshake mechanism (VALID/READY) over a simple bus with fixed timing?
