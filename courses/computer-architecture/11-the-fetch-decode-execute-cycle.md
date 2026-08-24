# Chapter 11: The Fetch-Decode-Execute Cycle

This is the heartbeat of every computer ever built. Every program that has ever run — from the first ENIAC calculation in 1945 to the AI inference running on your GPU today — passes through these three steps for every single instruction. No exceptions. No shortcuts. The fetch-decode-execute cycle is not a software abstraction or a design pattern: it is what a processor physically *is*. Understand this cycle in full, and you understand the atom from which all computing is built.

---

## Table of Contents

1. [The Big Picture: What the Cycle Is](#1-the-big-picture)
2. [The Internal Registers: MAR, MDR, IR, and PC](#2-the-internal-registers-mar-mdr-ir-and-pc)
3. [The Three Stages in Detail](#3-the-three-stages-in-detail)
4. [The System Bus: How Data Actually Travels](#4-the-system-bus-how-data-actually-travels)
5. [A Complete Concrete Trace: Adding Two Numbers in Memory](#5-a-complete-concrete-trace-adding-two-numbers-in-memory)
6. [Instruction-by-Instruction Traces: LOAD, STORE, ADD, BRANCH](#6-instruction-by-instruction-traces-load-store-add-branch)
7. [Timing Diagrams: Visualizing Time Across Multiple Cycles](#7-timing-diagrams-visualizing-time-across-multiple-cycles)
8. [Cycles Per Instruction (CPI): Measuring CPU Performance](#8-cycles-per-instruction-cpi-measuring-cpu-performance)
9. [The Bigger Picture: Why This All Matters](#9-the-bigger-picture-why-this-all-matters)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Big Picture

### The Assembly Line Analogy

Imagine a factory assembly line. A car chassis enters at one end. At Station 1, workers bolt on the engine. At Station 2, workers install the seats. At Station 3, workers attach the wheels. The same three stations, in the same order, for every car. The workers do not need to know why a particular car is being built — they just follow the same process every time.

A CPU works exactly like this. The "car" is an instruction. The three stations are **Fetch**, **Decode**, and **Execute**. Every instruction, regardless of what it does — add two numbers, read from memory, jump to a different location in the program — passes through the same three stations in the same order.

### The Cycle in One Sentence

**Fetch** the instruction from memory, **decode** what it means, **execute** what it says to do — then immediately start over with the next instruction.

```
  ┌──────────────────────────────────────────────────────────┐
  │                                                          │
  ▼                                                          │
┌──────┐    ┌──────────┐    ┌──────────┐                    │
│FETCH │───►│  DECODE  │───►│ EXECUTE  │────────────────────┘
└──────┘    └──────────┘    └──────────┘
   │
   │  (PC increments here — pointing to NEXT instruction)
   ▼
```

The cycle loops forever. There is no outer code that says "run the cycle again." The CPU simply *is* the cycle. When power flows through the chip, the cycle runs. The only way to stop it is to halt the CPU explicitly (with a HALT instruction) or cut the power.

### What Makes It Elegant

The genius of this design is its uniformity. The CPU does not need special logic for "what type of instruction is this?" at the top level. It always fetches, always decodes, always executes. The *content* of those stages differs per instruction, but the rhythm — the 1-2-3 beat — never changes.

This uniformity is what makes pipelining possible (covered in a later chapter). If every instruction follows the same sequence of steps, you can overlap them: while instruction N is being decoded, instruction N+1 can be fetched. Like an assembly line running multiple cars simultaneously.

---

### Quick Check 1.1

> 1. What are the three stages of the fetch-decode-execute cycle?
> 2. After executing an instruction, does the CPU need special logic to decide whether to fetch the next one?
> 3. Using the factory analogy, what is the "car" and what are the "stations"?

---

## 2. The Internal Registers: MAR, MDR, IR, and PC

Before tracing through the cycle, you need to meet the four internal registers that make it work. Think of these as the CPU's scratchpad for managing the cycle itself — not for storing your program's data, but for the CPU's own bookkeeping.

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                        CPU INTERNALS                            │
  │                                                                 │
  │  ┌──────────────────┐     ┌──────────────────────────────────┐  │
  │  │   PC             │     │   IR (Instruction Register)      │  │
  │  │ Program Counter  │     │ Holds current instruction bits   │  │
  │  │ Points to NEXT   │     │ 32 bits (for 32-bit instructions)│  │
  │  │ instruction addr │     └──────────────────────────────────┘  │
  │  └──────────────────┘                                           │
  │                                                                 │
  │  ┌──────────────────┐     ┌──────────────────────────────────┐  │
  │  │   MAR            │     │   MDR                            │  │
  │  │ Memory Address   │     │ Memory Data Register             │  │
  │  │ Register         │     │ Holds data being read from       │  │
  │  │ Holds the address│     │ or written to memory             │  │
  │  │ we want to access│     └──────────────────────────────────┘  │
  │  └──────────────────┘                                           │
  └─────────────────────────────────────────────────────────────────┘
```

Let's understand each one deeply.

---

### The Program Counter (PC)

The **Program Counter** (sometimes called the **Instruction Pointer** or **IP** in x86) is the CPU's bookmark. It holds the memory address of the *next* instruction to be fetched.

When the CPU starts up, the PC is set to a specific address (often called the **reset vector** — perhaps address 0x00000000 or 0xFFFF0000 depending on the architecture). That address contains the first instruction of the boot code.

After fetching each instruction, the PC is incremented:
- By 4 bytes for a 32-bit fixed-length instruction set (like RISC-V or MIPS)
- By 2 bytes for a 16-bit instruction
- By a variable amount for variable-length instruction sets (like x86, where instructions can be 1-15 bytes long)

Branch and jump instructions work by *overwriting* the PC with a new address — effectively saying "the next instruction is not the one right after me; it's over there."

### The Instruction Register (IR)

The **Instruction Register** holds the bits of the instruction currently being processed. When the CPU fetches an instruction from memory, those bits land in the IR. The decode stage reads the IR and figures out what operation to perform.

Think of the IR like a recipe card pulled from a recipe box (memory). The PC is your index finger pointing to which card to pull next. The IR is the countertop where you lay out the card to read it.

### The Memory Address Register (MAR)

The **Memory Address Register** is the staging area for memory addresses. Whenever the CPU needs to access memory — whether to fetch an instruction or to load/store data — it first loads the desired address into the MAR. Then it signals the memory unit to perform the read or write.

The MAR is like dialing a telephone number before the call connects. You can't just "think" the address at the memory — you have to explicitly place it somewhere the memory bus can read.

### The Memory Data Register (MDR)

The **Memory Data Register** (also called the **Memory Buffer Register** or **MBR** in some texts) is the staging area for data being transferred to or from memory.

- On a **read** (LOAD): the data that memory returns lands in the MDR first, then the CPU copies it to where it's needed.
- On a **write** (STORE): the CPU places data in the MDR first, then signals memory to store MDR's contents at the address in MAR.

The MDR is like the loading dock at a warehouse. Incoming shipments are staged there before being brought into the building; outgoing shipments are staged there before being loaded onto the truck.

```
  Memory Access Pattern:

  READ from memory:                    WRITE to memory:
  ┌──────────────────────────────────┐  ┌──────────────────────────────────┐
  │ 1. MAR ← address                 │  │ 1. MAR ← address                 │
  │ 2. Signal: Memory READ           │  │ 2. MDR ← data to write           │
  │ 3. MDR ← Memory[MAR]             │  │ 3. Signal: Memory WRITE          │
  │ 4. Use MDR's contents            │  │ 4. Memory[MAR] ← MDR             │
  └──────────────────────────────────┘  └──────────────────────────────────┘
```

---

### Quick Check 2.1

> 1. What does the PC hold, and what happens to it after each fetch?
> 2. What is the role of the MAR? Why can the CPU not send an address directly to memory without the MAR?
> 3. During a LOAD instruction, which register does the data from memory land in first?

---

## 3. The Three Stages in Detail

Now let's trace through each stage step by step, showing every register transfer.

### Stage 1: FETCH

**Goal:** Get the next instruction from memory and place it in the IR. Update the PC.

The fetch stage is the same for *every* instruction — the CPU does not yet know what the instruction is. It just goes and gets it.

```
  FETCH Stage — Step-by-Step Register Transfers:
  ═══════════════════════════════════════════════

  Step F1:  MAR ← PC
            (Copy the PC address into the Memory Address Register)

  Step F2:  Memory READ: MDR ← Memory[MAR]
            (Signal memory to read from address in MAR)
            (Wait for memory to respond — the instruction bits arrive)

  Step F3:  IR ← MDR
            (Copy the fetched instruction bits from MDR into the IR)

  Step F4:  PC ← PC + 4
            (Increment the PC to point to the NEXT instruction)
            (This happens in parallel with step F3 in most designs)
```

Visually, here is what flows where:

```
                       MEMORY
                    ┌──────────┐
                    │ addr 100 │◄──── MAR receives PC = 100
                    │ addr 104 │      Memory outputs instruction bits
                    │ addr 108 │
                    └──────────┘
                         │
                         │  instruction bits flow down the memory bus
                         ▼
                    ┌──────────┐
                    │   MDR    │  ← instruction bits land here first
                    └──────────┘
                         │
                         ▼
                    ┌──────────┐
                    │    IR    │  ← instruction bits copied here for decode
                    └──────────┘

                    ┌──────────┐
                    │   PC     │  100 → 104 (incremented simultaneously)
                    └──────────┘
```

**Why increment the PC during fetch?** Because by the time execute is done and the CPU is ready to fetch again, the PC must already point to the next instruction. Also, some architectures compute branch targets relative to the *current* PC (after increment). Getting the increment done early simplifies the branch calculation.

---

### Stage 2: DECODE

**Goal:** Read the IR bits and figure out: what operation? which registers? any immediate value?

The decode stage is handled by the **Control Unit** (Chapter 10). It examines the opcode field of the IR and generates control signals — binary wires that tell every other component what to do.

```
  DECODE Stage — Step-by-Step:
  ═══════════════════════════

  Step D1:  Control Unit reads IR[opcode field]
            (Determines instruction type: R-type, Load, Store, Branch, etc.)

  Step D2:  Control Unit generates control signals
            RegWrite ← 1 or 0
            ALUSrc   ← 1 or 0
            MemRead  ← 1 or 0
            MemWrite ← 1 or 0
            Branch   ← 1 or 0
            ALUOp    ← ADD, SUB, AND, OR, etc.

  Step D3:  Register addresses extracted from IR
            rs1 ← IR[19:15]  (source register 1)
            rs2 ← IR[24:20]  (source register 2)
            rd  ← IR[11:7]   (destination register)

  Step D4:  Register file is read
            A ← Register[rs1]  (first operand)
            B ← Register[rs2]  (second operand)

  Step D5:  Immediate value sign-extended (if instruction uses one)
            imm ← SignExtend(IR[31:20])   (for I-type instructions)
```

Steps D3 through D5 happen simultaneously — the control unit, register file, and immediate extractor all read the IR in parallel.

```
  ┌──────────────────────────────────────────────────────────────────┐
  │                     IR (32 bits)                                 │
  │  [31:25][24:20][19:15][14:12][11:7][6:0]                        │
  │   funct7  rs2   rs1  funct3   rd  opcode                        │
  └──────────┬──────┬─────┬──────────┬──────┬───────────────────────┘
             │      │     │          │      │
             │      │     │          │      ▼
             │      │     │          │  ┌──────────────┐
             │      │     │          │  │ Control Unit │
             │      │     │          │  │ generates    │
             │      │     │          │  │ control sigs │
             │      │     │          │  └──────────────┘
             │      │     │          │
             ▼      ▼     ▼          ▼
         ┌──────────────────────┐  ┌──────────────────┐
         │    Register File     │  │ Immediate        │
         │  Read[rs1] → A       │  │ Sign Extender    │
         │  Read[rs2] → B       │  │ imm → Imm_val    │
         └──────────────────────┘  └──────────────────┘
```

**What is sign extension?** Many instructions encode a small constant — called an **immediate value** — directly in the instruction bits. These constants might be 12 bits wide (for a RISC-V I-type instruction), but the ALU works on 32-bit or 64-bit values. Sign extension copies the *sign bit* (the most significant bit) of the 12-bit value into the upper 20 bits. This preserves the numeric value, including negative numbers. For example, the 12-bit value `0b111111111110` represents -2. Sign-extended to 32 bits: `0b11111111111111111111111111111110` — still -2.

---

### Stage 3: EXECUTE

**Goal:** Carry out the instruction. This may involve the ALU, memory, and/or updating the PC.

Execute is the most varied stage — what happens here depends entirely on the instruction type. But the control signals generated during decode ensure that every component does exactly the right thing.

**For an ALU instruction (ADD, SUB, AND, etc.):**

```
  EXECUTE: ALU Instruction
  ════════════════════════
  Step E1:  ALU receives operands A and B (from register file)
            ALU performs operation indicated by ALUOp control signal
            Result ← ALU(A, B)

  Step E2:  Result written to destination register
            Register[rd] ← Result

  Step E3:  PC is already updated (done in Fetch)
            No memory access needed
```

**For a LOAD instruction:**

```
  EXECUTE: LOAD (read from memory into register)
  ═════════════════════════════════════════════
  Step E1:  ALU computes memory address
            MAR ← A + imm    (base register + offset)

  Step E2:  Memory READ
            MDR ← Memory[MAR]

  Step E3:  Write data to destination register
            Register[rd] ← MDR
```

**For a STORE instruction:**

```
  EXECUTE: STORE (write register value to memory)
  ════════════════════════════════════════════════
  Step E1:  ALU computes memory address
            MAR ← A + imm    (base register + offset)

  Step E2:  Prepare data to write
            MDR ← B          (value from rs2 register)

  Step E3:  Memory WRITE
            Memory[MAR] ← MDR

  Step E4:  No register writeback (RegWrite = 0)
```

**For a BRANCH instruction:**

```
  EXECUTE: BRANCH (conditional jump)
  ═══════════════════════════════════
  Step E1:  ALU subtracts or compares A and B
            Sets Zero flag if A == B (for BEQ)
            Sets Negative flag if A < B (for BLT)

  Step E2:  Compute branch target address
            BranchTarget ← PC_at_fetch + SignExtend(imm)
            (PC was already incremented to PC+4 during fetch,
             but branch target is calculated from the fetch PC)

  Step E3:  Decide which PC to use
            IF (condition is TRUE):
                PC ← BranchTarget   (branch taken)
            ELSE:
                PC ← PC + 4         (already incremented, keep it)
                (branch not taken)
```

**The key insight about branches:** The PC was incremented during fetch to PC+4. If the branch is NOT taken, that incremented value is correct — the next instruction is at PC+4. If the branch IS taken, execute overwrites the PC with the branch target. So execute either *corrects* the PC (taken) or *confirms* it (not taken).

---

### Quick Check 3.1

> 1. List every register that changes during the Fetch stage.
> 2. During Decode, what is happening to the register file while the control unit is reading the opcode?
> 3. For a Store instruction, does the destination register (rd) get written? What gets written instead?

---

## 4. The System Bus: How Data Actually Travels

Between the CPU and memory runs the **system bus** — a set of parallel wires carrying signals between components. Understanding the bus helps you understand why MAR and MDR exist at all.

### Three Sub-Buses

```
  ┌──────────────────────────────────────────────────────────────────┐
  │                         SYSTEM BUS                               │
  │                                                                  │
  │  ┌──────────────────────────────────────────────────────────┐   │
  │  │ ADDRESS BUS  (32 bits wide)                               │   │
  │  │ Carries: memory address from MAR                         │   │
  │  │ Direction: CPU → Memory (one-way)                        │   │
  │  └──────────────────────────────────────────────────────────┘   │
  │                                                                  │
  │  ┌──────────────────────────────────────────────────────────┐   │
  │  │ DATA BUS  (32 or 64 bits wide)                           │   │
  │  │ Carries: data between MDR and memory                     │   │
  │  │ Direction: bidirectional (read: mem→CPU, write: CPU→mem) │   │
  │  └──────────────────────────────────────────────────────────┘   │
  │                                                                  │
  │  ┌──────────────────────────────────────────────────────────┐   │
  │  │ CONTROL BUS  (several lines)                             │   │
  │  │ Carries: READ/WRITE signals, READY signal, etc.          │   │
  │  │ Direction: mostly CPU → Memory; READY goes mem → CPU     │   │
  │  └──────────────────────────────────────────────────────────┘   │
  └──────────────────────────────────────────────────────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
  ┌──────────┐         ┌──────────┐         ┌──────────┐
  │   CPU    │         │   RAM    │         │  I/O     │
  │  (MAR    │         │  (stores │         │ devices  │
  │   MDR)   │         │   data)  │         │          │
  └──────────┘         └──────────┘         └──────────┘
```

### Why MAR and MDR Exist

You might wonder: why not just let the CPU send the address directly down the address bus, without buffering it in the MAR first? The answer is timing and arbitration.

The bus is *shared* — multiple components can request to use it. The MAR acts as the CPU's output buffer: it holds the address stable on the address bus while memory completes its read operation (which can take many nanoseconds). The CPU does not have to keep holding the address manually; the MAR holds it while the CPU can proceed with other operations.

Similarly, the MDR acts as the buffer at the "receiving dock." Memory puts the data on the data bus; the MDR captures it so the bus can be freed while the CPU processes the data.

Think of it like postal mail. You don't *hold* your letter in the air until the recipient reads it. You drop it in a mailbox (MDR/MAR) and the postal system handles the rest.

### Bus Transactions During Fetch

```
  Clock Cycle 1:
    CPU places PC value onto address bus (via MAR)
    CPU asserts READ on control bus

  Clock Cycle 2 (memory latency):
    Memory decodes address, reads the requested location
    Memory places instruction bits on the data bus
    Memory asserts READY on control bus

  Clock Cycle 3:
    CPU reads data bus into MDR
    CPU copies MDR into IR
    CPU releases address bus
    CPU increments PC
```

In fast cache-based systems, steps 1-3 may complete in a single cycle. In slow RAM-only systems, many cycles might pass between steps 1 and 3.

---

### Quick Check 4.1

> 1. What are the three sub-buses in a system bus? What does each carry?
> 2. Why does the address bus go one-way (CPU to memory) while the data bus is bidirectional?
> 3. If memory needs 5 clock cycles to respond to a read, how many cycles does the CPU spend idle waiting? What might it do instead?

---

## 5. A Complete Concrete Trace: Adding Two Numbers in Memory

Now let us trace a real program all the way through the cycle, showing every register value change, every bus transaction, and every control signal. We will use a simple assembly program that:

1. Loads the number 15 from memory address 200
2. Loads the number 27 from memory address 204
3. Adds them together
4. Stores the result (42) back to memory address 208

### The Program

```assembly
# Assume this program lives at addresses 100-115
# Registers: R1=scratch, R2=first number, R3=second number, R4=result

LOAD  R2, 200      # R2 ← Memory[200]   (loads the value 15)
LOAD  R3, 204      # R3 ← Memory[204]   (loads the value 27)
ADD   R4, R2, R3   # R4 ← R2 + R3       (computes 15 + 27 = 42)
STORE R4, 208      # Memory[208] ← R4   (stores 42)
HALT               # Stop
```

### Memory Contents Before Execution

```
  Address │ Contents
  ────────┼──────────────────────────────────────────────────
  100     │  LOAD R2, 200    (instruction bits for "load R2 from address 200")
  104     │  LOAD R3, 204    (instruction bits for "load R3 from address 204")
  108     │  ADD  R4, R2, R3 (instruction bits for "add R2+R3 into R4")
  112     │  STORE R4, 208   (instruction bits for "store R4 to address 208")
  116     │  HALT            (instruction bits for "halt")
  ...     │  ...
  200     │  15              (first number to add)
  204     │  27              (second number to add)
  208     │  0               (will store result here)
```

### Initial CPU Register State

```
  PC  = 100   (start of program)
  IR  = 0     (undefined/empty)
  MAR = 0     (undefined)
  MDR = 0     (undefined)
  R1  = 0, R2 = 0, R3 = 0, R4 = 0
```

---

### Instruction 1: LOAD R2, 200

**FETCH:**

```
  Step F1:  MAR ← PC
            MAR = 100

  Step F2:  MDR ← Memory[MAR]
            (Bus: address bus = 100, control bus = READ)
            (Memory returns instruction bits for LOAD R2, 200)
            MDR = [LOAD R2, 200 encoding]

  Step F3:  IR ← MDR
            IR = [LOAD R2, 200 encoding]

  Step F4:  PC ← PC + 4
            PC = 104

  Register state after Fetch:
    PC=104  MAR=100  MDR=[LOAD encoding]  IR=[LOAD encoding]
```

**DECODE:**

```
  Control Unit reads IR opcode → LOAD instruction
  Control signals generated:
    RegWrite  = 1   (yes, we will write to a register)
    MemRead   = 1   (yes, we need to read memory)
    MemWrite  = 0   (no memory write)
    ALUSrc    = 1   (second ALU input is immediate, not register)
    MemToReg  = 1   (write memory data to register, not ALU result)

  From IR:
    rd  = R2                  (destination register)
    imm = 200                 (memory address to load from)

  Register file read:
    A ← R0 = 0                (base register, often R0 for absolute addresses)
```

**EXECUTE:**

```
  Step E1:  Compute memory address
            MAR ← A + imm = 0 + 200 = 200

  Step E2:  Memory READ
            (Bus: address bus = 200, control bus = READ)
            MDR ← Memory[200]
            MDR = 15

  Step E3:  Write to destination register
            R2 ← MDR = 15

  Register state after Execute:
    PC=104  MAR=200  MDR=15  IR=[LOAD R2,200]  R2=15
```

---

### Instruction 2: LOAD R3, 204

**FETCH:**

```
  MAR ← PC (= 104)
  MDR ← Memory[104]  → [LOAD R3, 204 encoding]
  IR  ← MDR
  PC  ← 108

  After Fetch:  PC=108  MAR=104  IR=[LOAD R3,204]
```

**DECODE:**

```
  Same control signals as Instruction 1 (another LOAD)
  rd = R3,  imm = 204
```

**EXECUTE:**

```
  MAR ← 0 + 204 = 204
  MDR ← Memory[204] = 27
  R3  ← MDR = 27

  After Execute:  PC=108  MAR=204  MDR=27  R2=15  R3=27
```

---

### Instruction 3: ADD R4, R2, R3

**FETCH:**

```
  MAR ← PC (= 108)
  MDR ← Memory[108]  → [ADD R4, R2, R3 encoding]
  IR  ← MDR
  PC  ← 112

  After Fetch:  PC=112  MAR=108  IR=[ADD R4,R2,R3]
```

**DECODE:**

```
  Control Unit reads IR opcode → ADD instruction
  Control signals:
    RegWrite  = 1   (write result to register)
    MemRead   = 0   (no memory read)
    MemWrite  = 0   (no memory write)
    ALUSrc    = 0   (second input is register, not immediate)
    MemToReg  = 0   (write ALU result to register, not memory data)
    ALUOp     = ADD

  From IR:
    rs1 = R2  (first source register)
    rs2 = R3  (second source register)
    rd  = R4  (destination register)

  Register file read:
    A ← R2 = 15
    B ← R3 = 27
```

**EXECUTE:**

```
  Step E1:  ALU computes result
            Result ← A + B = 15 + 27 = 42

  Step E2:  Write result to destination register
            R4 ← 42

  Step E3:  No memory access (MemRead=0, MemWrite=0)

  After Execute:  PC=112  R2=15  R3=27  R4=42
```

---

### Instruction 4: STORE R4, 208

**FETCH:**

```
  MAR ← PC (= 112)
  MDR ← Memory[112]  → [STORE R4, 208 encoding]
  IR  ← MDR
  PC  ← 116

  After Fetch:  PC=116  MAR=112  IR=[STORE R4,208]
```

**DECODE:**

```
  Control Unit reads IR → STORE instruction
  Control signals:
    RegWrite  = 0   (no register write for STORE)
    MemRead   = 0   (no memory read)
    MemWrite  = 1   (write to memory)
    ALUSrc    = 1   (address from immediate)
    MemToReg  = X   (don't care, no register write)

  From IR:
    rs2 = R4  (register whose value will be stored)
    imm = 208 (destination address)

  Register file read:
    A ← R0 = 0    (base register)
    B ← R4 = 42   (value to store)
```

**EXECUTE:**

```
  Step E1:  Compute memory address
            MAR ← A + imm = 0 + 208 = 208

  Step E2:  Stage data to write
            MDR ← B = 42

  Step E3:  Memory WRITE
            (Bus: address bus = 208, data bus = 42, control bus = WRITE)
            Memory[208] ← MDR = 42

  Step E4:  No register writeback (RegWrite = 0)

  After Execute:  PC=116  MAR=208  MDR=42  Memory[208]=42
```

---

### Final State

```
  Memory[200] = 15    (unchanged)
  Memory[204] = 27    (unchanged)
  Memory[208] = 42    (result stored successfully)

  R2 = 15, R3 = 27, R4 = 42

  PC = 116            (pointing to HALT instruction)
```

The program correctly computed 15 + 27 = 42 and stored the result. Every register transfer happened exactly as the cycle demands — no magic, no shortcuts.

---

### Quick Check 5.1

> 1. During Instruction 3 (ADD), what value does the MAR hold after the Fetch stage completes?
> 2. For the STORE instruction, why is RegWrite set to 0?
> 3. How many times did the data bus carry a value from memory to the MDR in this entire program?

---

## 6. Instruction-by-Instruction Traces: LOAD, STORE, ADD, BRANCH

Let us now trace one of each instruction type with full detail, including control signals.

### 6.1 LOAD Instruction Trace

```
  Instruction:  LOAD R5, 8(R1)    (load from address R1 + 8 into R5)
  Assume:       R1 = 1000

  ═══ FETCH ═══════════════════════════════════════════════
  MAR  ← PC           MAR = current PC
  MDR  ← Mem[MAR]     MDR = LOAD encoding
  IR   ← MDR          IR  = LOAD encoding
  PC   ← PC + 4       PC increments

  ═══ DECODE ══════════════════════════════════════════════
  opcode → LOAD
  Control signals:  RegWrite=1, MemRead=1, MemWrite=0,
                    ALUSrc=1, MemToReg=1
  rs1 = R1,  rd = R5,  imm = 8
  A ← R1 = 1000

  ═══ EXECUTE ══════════════════════════════════════════════
  MAR ← A + imm     MAR = 1000 + 8 = 1008
  MDR ← Mem[1008]   (memory read: address bus=1008, READ signal)
  MDR = [value at address 1008]
  R5  ← MDR         (MemToReg=1 routes MDR to register, not ALU result)
```

### 6.2 STORE Instruction Trace

```
  Instruction:  STORE R7, 12(R2)   (store R7 to address R2 + 12)
  Assume:       R2 = 2000, R7 = 99

  ═══ FETCH ═══════════════════════════════════════════════
  (Standard fetch as above)
  IR = STORE encoding, PC increments

  ═══ DECODE ══════════════════════════════════════════════
  opcode → STORE
  Control signals:  RegWrite=0, MemRead=0, MemWrite=1,
                    ALUSrc=1
  rs1 = R2,  rs2 = R7,  imm = 12
  A ← R2 = 2000
  B ← R7 = 99

  ═══ EXECUTE ══════════════════════════════════════════════
  MAR ← A + imm     MAR = 2000 + 12 = 2012
  MDR ← B           MDR = 99
  Mem[2012] ← MDR   (memory write: addr bus=2012, data bus=99, WRITE signal)
  (no register writeback)
```

### 6.3 ADD Instruction Trace

```
  Instruction:  ADD R8, R5, R6    (R8 ← R5 + R6)
  Assume:       R5 = 100, R6 = 200

  ═══ FETCH ═══════════════════════════════════════════════
  (Standard fetch as above)
  IR = ADD encoding, PC increments

  ═══ DECODE ══════════════════════════════════════════════
  opcode → ADD
  Control signals:  RegWrite=1, MemRead=0, MemWrite=0,
                    ALUSrc=0, MemToReg=0, ALUOp=ADD
  rs1 = R5, rs2 = R6, rd = R8
  A ← R5 = 100
  B ← R6 = 200

  ═══ EXECUTE ══════════════════════════════════════════════
  ALU: Result ← 100 + 200 = 300
  R8 ← 300          (RegWrite=1, MemToReg=0 routes ALU result to register)
  (no memory access: MAR and MDR unchanged)
```

### 6.4 BRANCH Instruction Trace

```
  Instruction:  BEQ R1, R2, -8    (if R1==R2, jump back 8 bytes)
  Assume:       R1 = 5, R2 = 5    (condition is TRUE — branch taken)
  Fetch PC:     0x1020

  ═══ FETCH ═══════════════════════════════════════════════
  MAR ← 0x1020
  MDR ← Mem[0x1020]   (BEQ encoding)
  IR  ← MDR
  PC  ← 0x1024        (PC incremented to 0x1020 + 4)

  ═══ DECODE ══════════════════════════════════════════════
  opcode → BEQ
  Control signals:  RegWrite=0, MemRead=0, MemWrite=0,
                    Branch=1, ALUOp=SUB (or compare)
  rs1 = R1, rs2 = R2, imm = -8
  A ← R1 = 5
  B ← R2 = 5

  Branch Target Address = Fetch_PC + imm = 0x1020 + (-8) = 0x1018

  ═══ EXECUTE ══════════════════════════════════════════════
  ALU: Result ← A - B = 5 - 5 = 0   → Zero flag SET

  Branch condition: BEQ checks Zero flag
    Zero = 1 (they are equal) → BRANCH TAKEN

  PC ← BranchTarget = 0x1018
  (overwrites the 0x1024 that was set during fetch)

  Next fetch will be at PC = 0x1018

  ═════════════════════════════════════════════════════════
  Had R1 ≠ R2 (e.g., R1=5, R2=6):
    ALU Result = 5 - 6 = -1  → Zero flag CLEAR
    Branch NOT taken: PC stays at 0x1024
```

---

### Quick Check 6.1

> 1. For a STORE instruction, both R2 and R7 are read from the register file during Decode. Why are two source registers needed?
> 2. In the BRANCH trace, the PC was set to 0x1024 during Fetch. What overwrites it, and when?
> 3. If a BEQ instruction has the condition FALSE (registers not equal), how many clock cycles does it "waste" compared to not being there at all?

---

## 7. Timing Diagrams: Visualizing Time Across Multiple Cycles

A **timing diagram** shows the values of signals and registers plotted against time. Each row is a signal; each column is a clock edge. Reading a timing diagram tells you exactly what is happening at every moment.

### Single-Cycle CPU: Four Instructions in Sequence

The program: LOAD, ADD, STORE, HALT — four instructions at addresses 100, 104, 108, 112.

```
  Clock:    ___     ___     ___     ___     ___
         __|   |___|   |___|   |___|   |___|   |__
           1       2       3       4       5

  PC:   100     104     108     112     116
        ├───────┤───────┤───────┤───────┤

  MAR:  100  200 104  204 108  208 112     116
        (fetch)(load)(fetch)(load)(fetch)(store)(fetch)

  MDR: [LOAD][data][ADD ][    ][STORE][R7  ][HALT]
           15           data         42

  IR:  [LOAD ][ADD   ][STORE ][HALT  ]
        ├──────┤───────┤───────┤───────┤

  Note: In single-cycle design, everything happens within one clock cycle.
        The timeline above compresses sub-operations within each cycle.
```

### A More Detailed View: Sub-Cycle Events

Within a single clock cycle, there are internal phases (not clock-edge-separated, but time-separated by propagation delays):

```
  Within ONE clock cycle (e.g., cycle 1 = LOAD instruction):

  Time ──────────────────────────────────────────────────►
       0ps     100ps   200ps   300ps   400ps   500ps   600ps
       │       │       │       │       │       │       │
       │←Fetch→│←Decode+RdReg─►│←──────Execute──────►│
       │       │       │       │       │       │       │
  PC   │100────────────────────────────────────────────│
  MAR  │100────┤ then →│200────────────────────────────│
  MDR  │  ←mem read→   │  ←data mem read──────────────│
  IR   │       │LOAD───────────────────────────────────│
  Ctrl │       │sig gen─────────────────────────────── │
  ALU  │       │       │       │addr=1008──────────────│
  Reg  │       │       │R1─────────────────────────────│→R5=val
       │       │       │       │       │       │       │
```

### Multi-Cycle Timing: Fetch + Execute Separated

In a **multi-cycle CPU** (a stepping stone between single-cycle and pipelined designs), each stage of the cycle takes its own clock cycle:

```
  Clock:  1       2       3       4       5       6       7
          ┌───┐   ┌───┐   ┌───┐   ┌───┐   ┌───┐   ┌───┐   ┌───┐
          │   │   │   │   │   │   │   │   │   │   │   │   │   │
  ────────┘   └───┘   └───┘   └───┘   └───┘   └───┘   └───┘   └──

  Instruction 1 (LOAD):
  Clock 1: IF  — Fetch instruction from PC=100, PC→104
  Clock 2: ID  — Decode LOAD, read R1 from register file
  Clock 3: EX  — ALU computes 1000+8=1008, MAR=1008
  Clock 4: MEM — Read Memory[1008] → MDR
  Clock 5: WB  — Write MDR to R5

  Instruction 2 (ADD):
  Clock 6: IF  — Fetch from PC=104
  Clock 7: ID  — Decode ADD, read R5, R6
  ...

  Total for 2 instructions: 5 + 5 = 10 cycles
  (In a pipelined design, this would be approximately 5 + 1 = 6 cycles)
```

### The CPI (Cycles Per Instruction) Visualization

```
  Single-cycle CPU (CPI = 1.0):

  Cycle:  1       2       3       4       5
          LOAD    ADD     STORE   ADD     BRANCH
          (done)  (done)  (done)  (done)  (done)

  Multi-cycle CPU (variable CPI):

  Cycle:  1  2  3  4  5  6  7  8  9  10 11
          IF ID EX ME WB IF ID EX WB IF ID ...
          ├──── LOAD ────┤├── ADD ──┤
          (5 cycles)      (4 cycles, no MEM)
```

---

### Quick Check 7.1

> 1. In a single-cycle CPU, does one clock cycle correspond to one stage (Fetch, Decode, or Execute) or to the entire instruction?
> 2. In the multi-cycle CPU timing diagram, why does ADD take 4 cycles instead of 5?
> 3. What would a timing diagram look like for two consecutive HALT instructions?

---

## 8. Cycles Per Instruction (CPI): Measuring CPU Performance

**CPI** — Cycles Per Instruction — is one of the most fundamental metrics in computer architecture. It tells you how many clock cycles, on average, each instruction requires.

### The Performance Equation

```
  CPU Execution Time = Instruction Count × CPI × Clock Period

  Or equivalently:
  CPU Execution Time = Instruction Count × CPI / Clock Frequency
```

This simple equation has three independent knobs:
1. **Instruction Count** — reduced by better compilers, better algorithms, or higher-level ISAs
2. **CPI** — reduced by pipelining, caching, out-of-order execution
3. **Clock Frequency** — increased by better transistors, better circuit design, deeper pipelines

Improving one knob often worsens another. This tension defines computer architecture.

### CPI in a Single-Cycle CPU

In a single-cycle CPU, every instruction takes exactly **one** clock cycle. CPI = 1.0 by definition.

But the clock must be slow enough for the *slowest* instruction to complete all its work. If a Load instruction needs 650ps and an ADD needs 450ps:

```
  Single-cycle clock period = 650ps  (must accommodate Load)
  Clock frequency = 1 / 650ps ≈ 1.54 GHz

  Time for 1000 instructions (mix of LOADs and ADDs):
  = 1000 × 1.0 × 650ps = 650,000ps = 650ns
```

But the ADDs only *needed* 450ps — they wasted 200ps each.

### CPI in a Multi-Cycle CPU

A multi-cycle CPU allows different instructions to take different numbers of cycles. The clock can be fast (one stage per cycle), and simple instructions finish sooner:

```
  Instruction Type    Cycles Needed
  ─────────────────   ─────────────
  R-type (ADD, SUB)   4 cycles  (IF, ID, EX, WB)
  I-type arithmetic   4 cycles
  Load                5 cycles  (IF, ID, EX, MEM, WB)
  Store               4 cycles  (IF, ID, EX, MEM — no WB)
  Branch              3 cycles  (IF, ID, EX — no MEM, no WB)
```

If a typical program has this instruction mix:
```
  25% Loads       → 5 cycles each
  10% Stores      → 4 cycles each
  52% ALU ops     → 4 cycles each
  13% Branches    → 3 cycles each
```

Average CPI = (0.25 × 5) + (0.10 × 4) + (0.52 × 4) + (0.13 × 3)
             = 1.25 + 0.40 + 2.08 + 0.39
             = 4.12 cycles per instruction

If the clock period is 200ps (one stage = 200ps):
- Multi-cycle clock frequency = 1/200ps = 5 GHz
- Single-cycle clock frequency = 1/1000ps = 1 GHz (5 stages × 200ps each)

```
  Performance comparison:
  
  Single-cycle:  1000 instructions × 1.0 CPI × 1000ps = 1,000,000ps = 1μs
  
  Multi-cycle:   1000 instructions × 4.12 CPI × 200ps = 824,000ps = 0.824μs
  
  Multi-cycle is ~21% faster, despite higher CPI, because it runs a faster clock.
```

### CPI in a Pipelined CPU

Pipelining aims for CPI = 1.0 while running a *fast* clock (one stage per cycle). In an ideal 5-stage pipeline:

```
  Pipelined clock period  = 200ps  (one stage)
  Pipelined CPI           = 1.0    (ideal, ignoring hazards)
  
  Performance: 1000 × 1.0 × 200ps = 200,000ps = 200ns
  
  That is 5× faster than single-cycle!
```

Real pipelines have **hazards** (data, control, structural) that increase effective CPI above 1.0, but modern CPUs use sophisticated techniques to keep CPI close to 1.0 or even below 1.0 (via superscalar execution — multiple instructions per cycle).

### CPI Comparison Table

```
  ┌─────────────────────┬─────────┬───────────┬──────────────────┐
  │ CPU Design          │ CPI     │ Clock     │ Relative Speed   │
  ├─────────────────────┼─────────┼───────────┼──────────────────┤
  │ Single-cycle        │ 1.0     │ Slow      │ Baseline         │
  │ Multi-cycle         │ 3-5     │ Fast      │ ~20% faster      │
  │ Pipelined (ideal)   │ 1.0     │ Fast      │ 5× faster        │
  │ Pipelined (real)    │ 1.1-1.4 │ Fast      │ 3-4× faster      │
  │ Superscalar (modern)│ 0.2-0.5 │ Very fast │ 10-30× faster    │
  └─────────────────────┴─────────┴───────────┴──────────────────┘

  Note: "faster" means relative to a simple single-cycle baseline.
  Modern CPUs at 3-4 GHz with CPI ~0.3 execute ~10 billion instructions/sec.
```

### What Raises CPI in Real CPUs

Even a pipelined CPU sees CPI > 1.0 in practice:

| Hazard Type | Cause | Effect on CPI |
|---|---|---|
| Data hazard | Instruction needs result from previous instruction | Insert 1-3 stall cycles |
| Control hazard | Branch direction not known until Execute | Insert 1-2 stall cycles (misprediction penalty) |
| Structural hazard | Two instructions need same hardware simultaneously | Insert 1 stall cycle |
| Cache miss (I-cache) | Instruction not in cache | Insert 50-200 stall cycles |
| Cache miss (D-cache) | Data not in cache | Insert 50-200 stall cycles |

Branch prediction (used by all modern CPUs) dramatically reduces control hazard stalls. Cache hierarchies reduce the stall impact of cache misses. These are among the most important architectural innovations ever developed.

---

### Quick Check 8.1

> 1. Write the CPU performance equation. What are its three terms?
> 2. Why does a multi-cycle CPU often outperform a single-cycle CPU even though its CPI is higher?
> 3. A program has 40% Load instructions (5 cycles each) and 60% ADD instructions (4 cycles each) in a multi-cycle CPU. What is the average CPI?

---

## 9. The Bigger Picture: Why This All Matters

### Every Computer Ever Built

It is worth pausing to appreciate what we have just traced. The fetch-decode-execute cycle was not invented by Intel or ARM. It was conceived by John von Neumann and his colleagues in the 1940s, and every general-purpose computer ever built — from the ENIAC to the Apollo Guidance Computer to the smartphone in your pocket — operates on this same principle.

The specific registers have different names. The instruction formats differ. The number of pipeline stages varies. But the fundamental operation — fetch a pattern of bits from memory, interpret those bits as a command, execute that command, repeat — is universal.

### Why Memory Is the Heart of the Cycle

Notice that memory appears in *every* stage of the cycle:
- **Fetch**: memory read (instruction fetch)
- **Execute**: optional memory read (load) or write (store)

The CPU's speed is, in a deep sense, constrained by how fast it can access memory. This is the **memory bottleneck** — also called the **von Neumann bottleneck**. No matter how fast the ALU, if memory is slow, the CPU waits. This drove the invention of cache memory, register files, and complex memory hierarchies — all of which we will explore later in this course.

### The Illusion of Programs

Here is perhaps the most mind-bending insight about the cycle: **a program is just data in memory**. The CPU fetches bits from memory and interprets them as instructions. But those same bits could be interpreted as numbers, or as text, or as an image. There is nothing intrinsically "executable" about instruction bits — the CPU just treats whatever it fetches as instructions.

This is exactly how **viruses** work: they place code (data) somewhere the CPU will eventually execute. This is why **buffer overflows** are dangerous: they let attackers put their own data into memory in a way the CPU later treats as instructions. The uniform nature of the fetch-decode-execute cycle — "just execute whatever is at PC" — is both its greatest strength and a source of profound security implications.

### The Heartbeat Metaphor

The fetch-decode-execute cycle is often called the CPU's "heartbeat." This is more than poetic. Just as a heartbeat is the repeated, rhythmic action that keeps a biological system alive, the cycle is the repeated, rhythmic action that keeps a computational system operating. Stop the cycle (cut power, execute HALT) and the computer ceases to compute. Resume it (restart, resume from halt) and computation continues exactly where it left off.

Every animation you've watched, every web page you've loaded, every game you've played — these were the result of billions of fetch-decode-execute cycles happening per second, each one a tiny, perfectly executed three-step dance.

---

## Summary

- The **fetch-decode-execute cycle** is the fundamental, universal operating mode of all general-purpose CPUs. It repeats endlessly until explicitly stopped.

- **Four internal registers** manage the cycle:
  - **PC** (Program Counter): holds address of next instruction
  - **MAR** (Memory Address Register): buffers addresses for memory access
  - **MDR** (Memory Data Register): buffers data to/from memory
  - **IR** (Instruction Register): holds the current instruction being processed

- **Fetch stage** (same for every instruction):
  - MAR ← PC
  - MDR ← Memory[MAR] (bus transaction)
  - IR ← MDR
  - PC ← PC + 4

- **Decode stage**:
  - Control unit reads opcode from IR, generates control signals
  - Register file is read (rs1, rs2) in parallel
  - Immediate value is sign-extended if needed

- **Execute stage** (varies by instruction type):
  - ADD/ALU: ALU computes, result written to register
  - LOAD: ALU computes address, memory read, data written to register
  - STORE: ALU computes address, register value written to memory
  - BRANCH: ALU compares, PC overwritten with branch target if condition true

- **The system bus** has three sub-buses: address (CPU→memory), data (bidirectional), control (signals like READ/WRITE). MAR and MDR interface the CPU to these buses.

- **CPI** (Cycles Per Instruction) combined with clock frequency and instruction count gives CPU performance. Single-cycle CPU: CPI=1.0 but slow clock. Multi-cycle: variable CPI but faster clock. Pipelined: CPI≈1.0 with fast clock — best of both worlds.

- The cycle's universality means a program is just data in memory — the CPU executes whatever bits it finds at the PC address, which has profound implications for security and system design.

---

## Exercises

### Easy

1. Trace through the complete fetch-decode-execute cycle for the single instruction `ADD R5, R3, R4`, assuming R3=12 and R4=8. List every register (PC, MAR, MDR, IR, R5) and its value at the end of each stage. The instruction is located at address 0x2000.

2. For a LOAD instruction `LOAD R7, 500`, why must the execute stage read memory *twice* across the whole instruction's lifetime? Which memory read is different from which?

3. A CPU running at 2 GHz executes a program with 6 billion instructions. The average CPI is 1.5. How long does the program take to run? (Give your answer in seconds.)

### Medium

4. Consider the following program:

   ```assembly
   LOAD  R1, 300        # Load value from address 300
   LOAD  R2, 304        # Load value from address 304
   SUB   R3, R1, R2     # R3 = R1 - R2
   BEQ   R3, R0, done   # If R3 == 0, jump to "done"
   STORE R3, 308        # Store result to address 308
   done: HALT
   ```

   Memory[300] = 50, Memory[304] = 50. Trace the program completely. What happens at the branch instruction? Does the STORE execute? What is Memory[308] at the end?

5. In the BRANCH instruction trace (Section 6.4), the branch target was calculated as `Fetch_PC + imm` where Fetch_PC is the address the instruction was fetched from (0x1020), NOT the already-incremented PC (0x1024). Why does it use the un-incremented PC? What would go wrong if it used PC+4 instead?

6. A multi-cycle CPU has the following stage times: IF=200ps, ID=150ps, EX=300ps, MEM=250ps, WB=100ps. The instruction mix is 20% Load (5 stages), 10% Store (4 stages: IF/ID/EX/MEM), 60% ALU (4 stages: IF/ID/EX/WB), 10% Branch (3 stages: IF/ID/EX). Calculate the average CPI and the performance (instructions per second). How does this compare to a single-cycle CPU with clock period = sum of all 5 stage times?

### Hard

7. **The Control Signal Difference Between LOAD and STORE**: Both LOAD and STORE compute a memory address using the ALU (base + offset). But LOAD reads from that address and STORE writes to it. Design a complete control signal table for both instructions. Include: RegWrite, MemRead, MemWrite, ALUSrc, MemToReg, and ALUOp. Then explain: for STORE, two registers are read during decode (the base register for the address, and the register whose value will be stored). But for LOAD, only one register is read (the base). How does the control unit know which register values to use?

8. **The Von Neumann Bottleneck**: In a simple CPU, the same memory bus is used for instruction fetch (during Fetch stage) and data access (during Execute for LOAD/STORE). This creates a **structural hazard** — both operations cannot happen simultaneously on the same bus.
   - (a) If a single-cycle CPU must do instruction fetch AND data memory access in the same cycle (for a LOAD instruction), and both use the same memory, what problem arises?
   - (b) How does the Harvard Architecture solve this? Draw a simple diagram showing the two separate memories.
   - (c) Modern CPUs use a modified Harvard architecture with a unified memory but separate L1 instruction cache and L1 data cache. How does this solve the problem in practice?

9. **Self-Modifying Code**: The fetch-decode-execute cycle always fetches from the address in the PC. In principle, a program could write new instruction bits to a memory address, then branch to that address and execute the newly written instructions. This is called **self-modifying code**.
   - (a) Trace through how this would work using the instructions we have studied: STORE (to write new instruction bits), then a BRANCH to the stored location.
   - (b) Why do modern CPUs and operating systems generally prohibit or discourage self-modifying code? (Hint: consider the instruction cache and security.)
   - (c) Despite the risks, can you think of one legitimate use case for self-modifying code?
