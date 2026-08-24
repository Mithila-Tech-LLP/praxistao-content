# Chapter 62: The IIT Madras SHAKTI Program — Origins and Vision

India's semiconductor industry imports nearly $25 billion worth of chips annually. Despite having world-class software engineers and a growing electronics market, India has historically been absent from chip design and fabrication at the leading edge. The **SHAKTI program** at IIT Madras is India's most ambitious attempt to change this — by developing a complete, open-source processor ecosystem based on RISC-V that could seed indigenous chip design capability. This chapter covers the origins of the program, the people behind it, the funding and institutional context, and the vision that drives it.

## Table of Contents

1. [Why India Needs Indigenous Processors](#1-why-india-needs-indigenous-processors)
2. [SHAKTI's Origins at IIT Madras](#2-shaktis-origins-at-iit-madras)
3. [The RISC-V Connection](#3-the-risc-v-connection)
4. [Funding and Institutional Support](#4-funding-and-institutional-support)
5. [The SHAKTI Vision — Not Just Chips](#5-the-shakti-vision--not-just-chips)
6. [Global Context: Open-Source Hardware](#6-global-context-open-source-hardware)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Why India Needs Indigenous Processors

**Strategic dependency**: Every military system, every critical infrastructure, every smartphone in India runs on processors designed abroad — primarily by US companies (Intel, AMD, Qualcomm, ARM) or Taiwanese foundries. This creates:
- **National security risk**: chips can contain backdoors or kill-switches; supply can be cut in geopolitical conflict
- **Economic dependency**: India exports software services but imports hardware — a structural imbalance
- **Technology transfer limitations**: US export controls (EAR) restrict access to cutting-edge design tools and processes

**The import gap**: India imported $25B of semiconductors in FY2023 and exported only $6B (mostly packaging, not design). The government's target: reduce import dependency and grow domestic semiconductor from $25B import to $300B domestic production by 2030.

**Software strength, hardware gap**: India produces world-class software engineers (IITs, IIITs), many of whom join AMD, Intel, Qualcomm, and ARM's India design centers. But the IP created stays with US companies. A domestic processor program retains that IP in India.

**The China example**: China recognized this dependency in the 2010s and began massive investment in semiconductors. RISC-V gives India (and China) a way to develop processors without ARM/x86 licensing fees. Alibaba's Xuantai T-Head, HiSilicon Kirin — these are the models SHAKTI can emulate.

### Quick Check
> 1. Why is India's chip import dependency a national security concern?
> 2. What is the target for India's domestic semiconductor production by 2030?
> 3. Why does RISC-V specifically enable countries to develop independent processor ecosystems?

---

## 2. SHAKTI's Origins at IIT Madras

**IIT Madras** (Indian Institute of Technology Madras, Chennai) is one of India's premier technology institutions and home to the RISE (RISC-V based Systems and Embedded Research) lab, which houses the SHAKTI project.

**The founding team**: Prof. Kamakoti Veezhinathan (now Director of IIT Madras) was the primary driving force behind SHAKTI. He recognized in 2014–2015 that RISC-V's open ISA provided a unique opportunity for India to develop indigenous processors without licensing barriers.

**Timeline:**
- 2014: RISC-V ISA becomes public (RISC-V Foundation formed)
- 2015: IIT Madras begins SHAKTI development under Prof. Kamakoti
- 2017: First SHAKTI chip (Riscy): tape-out at Intel's IITM-Intel Research Center
- 2018: SHAKTI C-class 64-bit processor silicon demonstrated
- 2019: SHAKTI deployed in Indian government applications (ISRO evaluation, strategic applications)
- 2021: SHAKTI processors in ISRO satellite ground station systems
- 2023: SHAKTI-T (security-enhanced) processor, multiple tape-outs

**Why IIT Madras?**
- Existing strong VLSI (Very Large Scale Integration) research program
- Government support through MHRD (Ministry of Human Resource Development)
- Intel partnership for early tape-outs
- Proximity to Chennai's semiconductor industry (Texas Instruments, Samsung have major design centers there)

**The name SHAKTI**: Named after the Sanskrit/Hindu concept of power/energy (Shakti = power). The name reflects the program's aspiration to build India's technological power through indigenous chip design.

### Quick Check
> 1. Who is the primary driving force behind the SHAKTI program?
> 2. What was significant about the 2017 Riscy tape-out?
> 3. Why is IIT Madras the right institution to lead this program?

---

## 3. The RISC-V Connection

SHAKTI chose RISC-V as its ISA foundation — not by accident, but by deliberate strategic reasoning.

**Why not x86 (Intel/AMD)?**
- Proprietary ISA: cannot implement without a license from Intel
- License terms are restrictive and expensive
- No path to full independence

**Why not ARM?**
- ARM architecture license is expensive ($millions+)
- ARM Holdings (owned by SoftBank, with significant stake now by NVIDIA's investment attempts) can revoke licenses in geopolitical disputes
- India cannot guarantee long-term license availability
- The 2020 Arm-NVIDIA acquisition attempt raised concerns about future access

**Why RISC-V:**
- Completely free and open ISA: no license fees, no royalties, no restrictions
- Extensible: custom instruction extensions are permitted without license
- RISC-V International foundation is Switzerland-based (neutral), governed by members (including IIT Madras, which is a member)
- Growing ecosystem: Western Digital, Google, NVIDIA, Alibaba, Qualcomm are all RISC-V members
- No US export control on the ISA itself (the ISA is just a specification — open knowledge)

**RISC-V India ecosystem**: Beyond IIT Madras SHAKTI, other Indian entities working with RISC-V:
- IIT Bombay: IITB-RISC (RISC-V processor for education)
- CDAC (Centre for Development of Advanced Computing): working on RISC-V for government applications
- Indian space research: ISRO evaluating RISC-V for space-grade processors

**SHAKTI as RISC-V ambassador**: IIT Madras has been a vocal advocate for RISC-V in India and contributed to RISC-V International standards discussions.

### Quick Check
> 1. Why did SHAKTI choose RISC-V over ARM?
> 2. What geopolitical concern made ARM licensing risky for India?
> 3. What is RISC-V International and why is its Swiss location relevant?

---

## 4. Funding and Institutional Support

**Government funding:**
- **MeitY (Ministry of Electronics and IT)**: Primary funder; total SHAKTI funding ~₹100 crore ($12M USD) over 2015–2023
- **MHRD (now MoE, Ministry of Education)**: IITM funding through IMPRINT (Impacting Research Innovation and Technology)
- **DST (Department of Science and Technology)**: Research grants
- **India Semiconductor Mission (ISM)**: launched 2021 — dedicated $10B for semiconductor incentives; SHAKTI is a beneficiary

**Industry partnerships:**
- **Intel**: Early tape-outs of SHAKTI processors at Intel IITM Research Center; Intel provided EDA tools and fabrication services
- **IIT Madras incubated startup (Mindgrove Technologies)**: Spin-off from SHAKTI project to commercialize SHAKTI IP in products; raised ~$3M in 2022
- **CDAC**: Collaboration on post-quantum cryptography extensions for SHAKTI-T

**Comparative context:**
- US: DARPA invested $500M+ in open-source hardware research (OpenROAD, RISC-V development)
- China: $150B government investment in semiconductor industry
- India $12M for SHAKTI: tiny by comparison, but meaningful as seed funding for an open ecosystem

**The IndiaChips program**: India announced a $10B incentive scheme (India Semiconductor Mission) in 2021 for chip manufacturing and design, with SHAKTI as a centerpiece for domestic design IP.

### Quick Check
> 1. What is the approximate total government funding for SHAKTI?
> 2. What is Mindgrove Technologies and what is its relationship to SHAKTI?
> 3. What is the India Semiconductor Mission and what is its scale?

---

## 5. The SHAKTI Vision — Not Just Chips

Prof. Kamakoti and the SHAKTI team have articulated a vision that goes beyond building individual chips:

**Five pillars of SHAKTI vision:**

1. **Open hardware ecosystem**: Publish all RTL code, documentation, and toolchains under open licenses — enabling Indian and global engineers to use, study, and improve SHAKTI. Not a proprietary system but a public good.

2. **Security by design**: Every SHAKTI processor includes hardware security features — Physical Memory Protection (PMP), secure enclaves, post-quantum cryptographic accelerators. India's strategic applications (military, government) require cryptographic independence.

3. **Education and workforce development**: SHAKTI source code is used in university curricula across India. Students design modifications, write drivers, compile operating systems for SHAKTI. Building the next generation of Indian chip designers.

4. **Indigenous manufacturing path**: Partnering with India's nascent semiconductor ecosystem (ISMC in Karnataka, Tata Semiconductor, HCL Semiconductors) to eventually tape out SHAKTI on Indian soil — not just design in India but fabricate in India.

5. **Commercial viability**: Through Mindgrove Technologies and other spinoffs, demonstrate that Indian-designed RISC-V chips can compete in the global market — IoT chips, automotive MCUs, edge AI processors.

```
SHAKTI ecosystem layers:
  
  Application layer:    IoT, edge AI, secure communications, aerospace
  Software layer:       Linux port, FreeRTOS, SHAKTI SDK, GCC toolchain
  IP layer:             SHAKTI processor cores (E-class to S-class)
  Tool layer:           Open-source EDA (OpenROAD), simulation
  Process layer:        TSMC 22nm, IIT-Intel tape-outs, India fab (target)
  
  Goal: India-designed and India-fabricated, end to end
```

### Quick Check
> 1. What are the five pillars of SHAKTI's vision?
> 2. How does SHAKTI contribute to education in India?
> 3. What is the long-term goal for SHAKTI regarding chip fabrication?

---

## 6. Global Context: Open-Source Hardware

SHAKTI is part of a global open-source hardware movement:

**Key players:**
- **RISC-V International**: the ISA foundation. Members include Google, Western Digital, NVIDIA, Qualcomm, Intel, Samsung. IIT Madras is a member.
- **OpenROAD**: DARPA-funded open-source place-and-route tool. Used by SHAKTI for its OpenLane flow.
- **SiFive**: US company commercializing RISC-V cores. Different from SHAKTI — proprietary implementation of an open ISA.
- **ESP32 (Espressif)**: Chinese RISC-V MCU, 100M+ units/year. Proves commercial viability.
- **T-Head (Alibaba)**: Xuantai 710 — 128-core server-grade RISC-V chip with 3.4 GHz frequency. Shows RISC-V can compete with ARM at server level.

**Open source hardware vs open source software**: Open source software (Linux, LLVM) has been transformative. Open source hardware (RISC-V, SHAKTI) is earlier in its adoption curve. The key difference: hardware requires expensive fabrication — unlike software, you can't just download and run it; you need a foundry.

**CHIPS and Science Act (US, 2022)**: $52.7 billion for US semiconductor R&D and manufacturing. Some of this flows to open-source hardware research (through DARPA, NSF). The US recognizes that open hardware ecosystems accelerate domestic chip capability.

**SHAKTI in the global picture**: IIT Madras SHAKTI is recognized internationally as one of the most serious open-source processor projects. Prof. Kamakoti presents regularly at RISC-V summits. SHAKTI's security extensions (post-quantum, tagged memory) have influenced broader RISC-V discussions.

### Quick Check
> 1. What is the difference between SiFive and SHAKTI?
> 2. What is the CHIPS and Science Act and what does it fund?
> 3. Why is open-source hardware adoption slower than open-source software?

---

## Summary

- **SHAKTI** is IIT Madras's open-source RISC-V processor program, initiated ~2015 under Prof. Kamakoti Veezhinathan.
- **Strategic motivation**: India imports $25B of chips annually; indigenous processor capability reduces strategic dependency.
- **RISC-V chosen** because it is free, open, extensible — no license fees, no foreign control of the ISA.
- **Funding**: ~₹100 crore ($12M) from MeitY/MHRD; Intel partnership for early tape-outs; Mindgrove Technologies spinoff for commercialization.
- **Vision**: Open hardware ecosystem, security by design, education platform, path to India-fabricated chips.
- **Global context**: Part of the RISC-V open-source hardware movement, alongside SiFive (commercial), T-Head/Alibaba (server-class), ESP32 (IoT).

---

## Exercises

### Easy
1. Why did SHAKTI choose RISC-V over ARM or x86 as its ISA?
2. What is the India Semiconductor Mission and why was it launched?
3. What is Mindgrove Technologies and what is its role?

### Medium
4. Strategic analysis: India imports $25B of chips and exports $6B (packaging services). The US imposes an ARM licensing embargo on India (hypothetical). (a) Which Indian products would be immediately affected? (b) How long would it take to replace ARM-based chips with RISC-V (SHAKTI) chips in: (i) a government missile guidance system, (ii) a consumer smartphone, (iii) an industrial IoT sensor? (c) What investments would India need to make in the next 5 years to ensure semiconductor sovereignty for strategic applications?
5. RISC-V geopolitics: RISC-V International is based in Switzerland (neutral country). In 2022, some US politicians proposed restricting Chinese access to RISC-V. (a) Is the ISA itself export-controlled (it's a specification published openly)? (b) What IS export-controlled: EDA tools (Synopsys, Cadence), leading-edge fab (ASML machines), or RISC-V processor IP blocks? (c) What is the difference between restricting the ISA specification vs restricting commercial implementations? (d) How would SHAKTI be affected by US restrictions on EDA tools?
6. Open-source hardware economics: The Linux kernel is developed by thousands of contributors at near-zero marginal cost. Compare: (a) What is the equivalent "development cost" for a processor like SHAKTI-C? (b) Why can't an open-source hardware project attract the same number of contributors as Linux? (hint: you can't run the hardware on your laptop for free) (c) RISC-V International has 3,500+ member organizations. What incentive does each member have to contribute to the common ISA rather than proprietary extensions?

### Hard
7. Technology gap analysis: SHAKTI's most advanced processor (S-class) targets ~1GHz at TSMC 22nm. Compare to Apple M3 (3nm, 4+ GHz). (a) What is the performance gap in raw compute (MIPS/GFLOPS)? (b) What is the power efficiency gap? (c) What process node and architectural improvements would SHAKTI need to be competitive with embedded ARM cores (Cortex-A55, ~3GHz at 5nm) in a 5-year roadmap? (d) Is "competitive" the right goal, or should SHAKTI target specific niches (secure government applications, space, defense) where strategic value outweighs performance-per-dollar?
8. Ecosystem bootstrapping: The ARM ecosystem (compilers, OS, libraries, toolchains) was built over 30 years. SHAKTI needs to bootstrap its software ecosystem from scratch. (a) What is the "minimum viable software stack" needed for SHAKTI to be used in a real product? List components. (b) RISC-V already has GCC/LLVM/Linux support. How much of this is directly usable for SHAKTI vs. how much needs SHAKTI-specific work? (c) The SHAKTI SDK must include BSP (Board Support Package), bootloader, RTOS port. Estimate person-months to develop each. (d) If India had a domestically-designed SHAKTI-based SoC in 10 million smartphones by 2030: what market share would this represent and what economic impact would be achieved (semiconductor import substitution)?
