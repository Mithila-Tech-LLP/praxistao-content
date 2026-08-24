# Chapter 93: 6G and the Future of Mobile Networks

> **"The most honest sentence to open a chapter about 6G with is this: as of this writing, 6G does not exist. There is no finalized standard, no deployed network, and no phone in anyone's pocket that speaks it. Everything in this chapter is either an early standards effort, a university or industry research paper, or a marketing department's guess about what 2030 will look like — and the single most useful skill this chapter can teach you is how to tell those three things apart when someone hands you a headline."**

---

## Table of Contents

1. [Why This Chapter Starts With a Warning](#1-why-this-chapter-starts-with-a-warning)
2. [What "6G" Even Refers To Right Now](#2-what-6g-even-refers-to-right-now)
3. [Where 5G Falls Short: The Real Problems Driving 6G Research](#3-where-5g-falls-short-the-real-problems-driving-6g-research)
4. [The Standards Timeline, Honestly](#4-the-standards-timeline-honestly)
5. [Terahertz Spectrum](#5-terahertz-spectrum)
6. [AI-Native RAN](#6-ai-native-ran)
7. [Integrated Sensing and Communication (ISAC)](#7-integrated-sensing-and-communication-isac)
8. [Non-Terrestrial Network Integration](#8-non-terrestrial-network-integration)
9. [Further-Out Ideas: Holographic Communication, Extreme XR](#9-further-out-ideas-holographic-communication-extreme-xr)
10. [The Labeled Summary Table: What's Actually True Right Now](#10-the-labeled-summary-table-whats-actually-true-right-now)
11. [A History Lesson in Hype: What 5G's Marketing Taught Us](#11-a-history-lesson-in-hype-what-5gs-marketing-taught-us)
12. [How to Evaluate Any Future-Tech Claim Yourself](#12-how-to-evaluate-any-future-tech-claim-yourself)
13. [Common Misconceptions](#13-common-misconceptions)
14. [What's Simplified Here](#14-whats-simplified-here)
15. [Interview Questions & Model Answers](#15-interview-questions--model-answers)
16. [Exercises](#16-exercises)
17. [Summary](#17-summary)

---

## 1. Why This Chapter Starts With a Warning

Every other chapter in this course has described something real: a protocol you can capture in Wireshark, a header you can decode byte by byte, a piece of hardware you can point to. This chapter is structurally different, and it would be dishonest to pretend otherwise. **6G is a research and early-standardization topic, not a deployed technology.** Chapter 92 closed by promising this chapter would label every claim honestly — deployed, commercially emerging, standardized, active research, or speculative — and that labeling discipline matters *more* here than anywhere else in this course, because the entire subject is more exposed to hype, marketing pre-announcement, and genuine scientific uncertainty than any prior chapter's material.

This isn't a reason to skip the topic. Understanding what serious researchers and standards bodies are actually working on — and being able to separate that from vendor press releases — is a genuinely useful, practical skill, and the specific research directions covered here (Sections 5-9) are real, funded, actively-published lines of work, not invented for this chapter. The goal is to give you an accurate picture of *where things stand*, not a confident prediction of what 6G will be.

---

## 2. What "6G" Even Refers To Right Now

As of this writing, "6G" refers to **a research agenda and an early standardization process**, not a finished specification. The relevant standards bodies are the same ones that defined every prior generation:

- **The ITU (International Telecommunication Union)**, which defines high-level global requirements for each "IMT" generation (recall Chapter 91, Section 7's IMT-Advanced requirement for 4G), has been developing its next framework, referred to as **IMT-2030**, describing the target capabilities a 6G-era system should aim for.
- **3GPP**, the same body that standardized every 5G release (Chapter 92), has begun early 6G-related study items, with formal 6G specification work expected to ramp up through the late 2020s.

**Status label: standards work in progress (early stage), not standardized.** No finalized 6G radio interface, core network architecture, or interoperability specification exists yet. Anything you read describing "6G's architecture" in concrete technical detail, as of this writing, is describing a research proposal or an early standards discussion input — not a ratified specification the way this course could describe LTE's EPC or 5G's 5GC.

---

## 3. Where 5G Falls Short: The Real Problems Driving 6G Research

6G research isn't happening in a vacuum — it's motivated by specific, real gaps between 5G's ambitions (Chapter 92, Section 2's three pillars) and 5G's actual, honestly-assessed real-world delivery, several of which Chapter 92 already flagged directly:

- **URLLC's most demanding promises (guaranteed sub-millisecond, five-nines-reliable latency) remain largely unrealized at public-network scale**, genuinely deployed today mostly in narrow, private, pilot contexts (Chapter 92, Sections 2 and 9) rather than broadly available.
- **Sub-6GHz vs. mmWave's coverage/capacity trade-off (Chapter 92, Section 3) is still an unsolved tension**, not a problem 5G fully closed — researchers are actively looking for spectrum and antenna approaches that narrow that gap further.
- **Network slicing and the fully software-defined 5GC (Chapter 92, Sections 6-7) remain unevenly deployed**, with a large share of real-world "5G" still running on NSA/EPC-derived infrastructure (Chapter 92, Section 10) rather than realizing 5G's own architectural vision fully, let alone extending it further.
- **New application demands are emerging that 5G's original 2015-2018-era requirements didn't fully anticipate**: dramatically more capable AI/ML workloads wanting to run partly on-device and partly at the network edge, immersive XR (extended reality) applications wanting far higher sustained throughput and lower motion-to-photon latency than even 5G's eMBB comfortably delivers, and a growing interest in networks that don't just carry communication but also directly sense their physical environment (Section 7).

6G research, in other words, is substantially about **finishing and extending 5G's own unfinished ambitions**, not starting from a blank slate — which is worth keeping in mind whenever "6G" is presented as an entirely novel leap rather than, in large part, a continuation of problems this course has already introduced you to.

---

## 4. The Standards Timeline, Honestly

| Milestone | Expected Timing | Status Label |
|---|---|---|
| ITU IMT-2030 framework (high-level requirements) | Largely established mid-2020s | **Standardized** (framework/requirements level only) |
| 3GPP formal 6G study and work items | Ramping up mid-to-late 2020s | **Active standards work in progress** |
| First finalized 3GPP 6G radio/core specification | Widely expected around 2028-2029 (based on the roughly decade-long cadence of prior generations: 4G ~2008, 5G ~2018) | **Speculative timeline** — a reasonable extrapolation from historical cadence, not a guaranteed date |
| First commercial 6G network launches | Widely expected around 2030 | **Speculative timeline** |

**Why the honest label matters here specifically:** every prior generation in this course (1G through 5G, Chapters 90-92) took roughly eight to ten years from serious standardization work beginning to commercial launch. That historical pattern is a reasonable basis for an *educated guess* about 6G's timeline — but it is still a guess, not a commitment from any standards body, and it's worth noticing that this entire section is the most speculative part of an already speculative chapter, purely about *timing*, before even getting to *what* 6G will technically include.

---

## 5. Terahertz Spectrum

**What it is:** research into using frequencies in the terahertz range (roughly 100 GHz to several THz) — even higher than 5G's mmWave band (24-100 GHz, Chapter 92, Section 3) — to access an enormous amount of additional, largely unused bandwidth, in principle enabling data rates far beyond anything 5G mmWave can achieve.

**Why it's hard, using physics this course already covered:** Chapter 92, Section 3 already showed mmWave paying a steep propagation-distance and penetration cost for its bandwidth; terahertz frequencies push that same trade-off to a genuinely extreme point. At terahertz frequencies, signals are absorbed even more severely by atmospheric moisture, blocked almost completely by any solid obstacle (including, in some studies, heavy fog or even a person's body passing between transmitter and receiver), and require extremely precise, tightly-focused beamforming (a further extension of Chapter 92, Section 4's Massive MIMO concept) just to have any usable range at all — likely measured in tens of meters, not the hundreds of meters mmWave already struggles to exceed.

**Status label: active research.** Terahertz communication is a genuinely active academic and industry research field, with real published experimental results (including short-range, high-speed link demonstrations in laboratory settings), but no standardized terahertz cellular interface exists, and the practical engineering challenges (semiconductor components efficient enough at these frequencies, workable beamforming at this scale, coping with severe atmospheric absorption) remain substantially unsolved for wide-area, real-world deployment as of this writing.

---

## 6. AI-Native RAN

**What it is:** the idea of building the Radio Access Network's core functions — resource scheduling, beamforming decisions, handover timing, interference management — around machine learning models trained on real network conditions, rather than the largely rule-based, human-engineered algorithms that have driven every prior generation's RAN (including 5G's, Chapter 92).

**An important nuance worth being precise about, split into two different maturity levels:**

- **Using AI/ML as an optimization tool bolted onto an otherwise conventional RAN — tuning parameters, predicting congestion, optimizing handover timing using trained models — is already happening today**, including in real 5G network operations by major carriers and vendors. This narrower use of AI in networking is **commercially emerging**, not speculative.
- **"AI-native RAN" as a *foundational design philosophy*** — where the radio access network's core control logic is designed from the ground up to be a learned, adaptive system rather than a fixed, human-specified protocol with AI as an add-on — **is a genuinely active research direction**, discussed in early 6G standards study items and academic literature, but nowhere near a finalized architectural approach. It raises real, currently-unresolved engineering questions this course's earlier chapters should make you appropriately skeptical of by default: how do you guarantee interoperability between an AI-native radio component built by one vendor and another built by a different vendor, when their internal decision logic isn't a fixed, auditable protocol the way 5G NR's specification is? How do you debug, in the field, a scheduling decision made by a trained model rather than a documented algorithm? Chapter 122's debugging playbook chapter should make clear just how much operational value comes from a system's behavior being predictable and specifiable — properties a learned system doesn't automatically retain.

**Status label: commercially emerging (as an optimization layer) / active research (as a foundational architecture).** Both halves of this claim are true simultaneously, and conflating them — treating today's real, deployed AI-assisted network optimization as proof that "AI-native 6G" is already a solved architectural approach — is exactly the kind of imprecision this chapter is trying to train you out of.

**Where this connects to a standards effort you already know about:** the O-RAN Alliance (Chapter 92, Section 15) has been standardizing an interface, the **RIC (RAN Intelligent Controller)**, specifically designed to let third-party applications — including ML models — plug into a RAN's control loop through an open, vendor-neutral interface, rather than each vendor building AI optimization as a closed, proprietary add-on. This is a genuinely useful, concrete anchor point for evaluating "AI-native RAN" claims: a specific, named, standardized *interface* for AI-based RAN control already exists and is being deployed today, which is meaningfully different evidence than a vague claim about "6G will use AI everywhere" — and precisely the kind of distinction Section 12's evaluation framework asks you to draw.

---

## 7. Integrated Sensing and Communication (ISAC)

**What it is:** the idea of using the same radio signals a cellular network already transmits for communication to *also* sense the physical environment — detecting the position, movement, and even some physical characteristics of objects (vehicles, people, obstacles) near a base station or device, similar in underlying physical principle to radar, but piggybacked onto ordinary communication signals rather than requiring separate, dedicated sensing hardware.

**Intuitive picture:** a 5G/6G base station's radio signal, in addition to carrying data to a phone, also reflects off nearby objects and returns to the transmitter (or a nearby receiver) — and if you can extract useful information from those reflections (how far away is that reflecting object, is it moving, how fast), you've gotten environmental sensing essentially "for free," reusing spectrum and hardware that's already there for communication. **Where the idea's promise is real but genuinely unproven at scale:** doing this well enough to be useful (distinguishing a pedestrian from a parked car reliably, at a range and accuracy useful for something like traffic management or industrial safety) requires solving real signal-processing problems that remain open research questions — extracting a clean, reliable sensing signal from a communication waveform that wasn't originally designed with sensing accuracy as a goal.

**Real, concrete proposed use cases:** traffic monitoring and vehicle detection integrated directly into roadside cellular infrastructure, industrial safety systems (detecting a person entering a hazardous zone using the same radios already providing factory connectivity, connecting directly to Chapter 92 Section 9's private-5G industrial context), and environmental/weather sensing using existing cellular infrastructure as an incidental sensor network.

**Status label: active research.** ISAC is one of the most frequently cited 6G research themes in both academic literature and early 3GPP/ITU study discussions, with real experimental demonstrations, but no standardized ISAC capability exists in any deployed cellular generation as of this writing, and it remains meaningfully unproven at the scale and reliability real safety-relevant use cases would require.

**Why ISAC connects naturally to the rest of this volume's radio material:** ISAC's basic technique — inferring an object's range and velocity from a reflected radio signal's timing and frequency shift — is the same underlying physics radar has used since the 1930s-40s, and it's a direct extension of ideas this course already introduced: Chapter 92, Section 4's beamforming already showed a base station computing precisely where a device is, using signal timing and phase, in order to aim a beam at it; ISAC proposes using that same kind of directional, timing-sensitive signal processing to characterize *any* reflecting object nearby, not just an intentionally-communicating device carrying its own radio. Seen this way, ISAC isn't an entirely separate technology bolted onto communication — it's an extension of capabilities Massive MIMO systems already need to have, pointed at a new problem.

---

## 8. Non-Terrestrial Network Integration

**What it is:** deliberately integrating satellite connectivity (Chapter 129's later, dedicated treatment of satellite internet and LEO constellations goes much further into this) directly into the standard cellular architecture, so that a device can seamlessly use a satellite link when no terrestrial tower is available, ideally without needing separate hardware or a separate subscription from its ordinary cellular service.

**Where this is genuinely further along than most other items in this chapter:** unlike terahertz spectrum or AI-native RAN, **NTN (Non-Terrestrial Network) support is already part of 5G-Advanced 3GPP specifications** (later 5G releases, beyond the initial Release 15 baseline Chapter 92 described), and real commercial services already exist that let an ordinary, unmodified smartphone send emergency messages or basic connectivity via satellite when out of terrestrial coverage — a capability several major smartphone manufacturers and carriers had already launched, in limited form, before this chapter's writing. Broader integration — full-speed cellular data seamlessly handed off to and from satellite links as a completely ordinary part of everyday connectivity, anywhere on Earth — remains a work in progress, expected to deepen substantially as part of 6G's broader architecture.

**Status label: commercially emerging** (basic satellite messaging/SOS capability, already shipping) **and standardized** (as part of 5G-Advanced) **for the basic case; active research** for full-featured, seamless broadband-speed non-terrestrial integration as a mainstream 6G capability.

---

## 9. Further-Out Ideas: Holographic Communication, Extreme XR

A handful of ideas appear regularly in 6G research papers, vendor whitepapers, and technology conference keynotes that deserve honest mention specifically *because* they're further from reality than Sections 5-8's items, and deserve a more skeptical label:

- **"Holographic" or true volumetric telepresence communication** — transmitting enough real-time visual data to render a realistic, full 3D representation of a remote person or object, viewable from any angle, rather than a flat 2D video call. This would require data rates and processing well beyond even 5G's most aggressive eMBB targets, and no clear, standardized path to delivering it exists yet.
- **Extreme, fully immersive XR with imperceptible motion-to-photon latency at population scale** — going meaningfully further than today's best VR/AR headsets, at a scale and cost making it broadly consumer-accessible rather than a specialized/enterprise product.
- **Fully autonomous, self-optimizing "zero-touch" networks**, extending Section 6's AI-native RAN idea to the entire network's operations (not just the radio access portion), managing themselves with minimal human network-engineering intervention.

**Status label: speculative.** These ideas appear regularly in forward-looking research and marketing material, and elements of each connect to real, ongoing research threads (Sections 5-7), but none of them describe a specific, technically-scoped engineering target the way Sections 5-8's items do — they are best understood as illustrative *visions* for what sufficiently advanced future networks might eventually enable, not near-term engineering roadmaps.

---

## 10. The Labeled Summary Table: What's Actually True Right Now

Consolidating every claim across Sections 4-9 into one table, using this course's full five-level honesty scale:

| Item | Status |
|---|---|
| ITU IMT-2030 high-level framework | Standardized (requirements-level) |
| 3GPP formal 6G radio/core specification | Active standards work in progress; not yet finalized |
| 6G commercial launch timeline (~2030) | Speculative (reasonable historical extrapolation, not a commitment) |
| Terahertz spectrum for cellular use | Active research |
| AI used to optimize existing 5G RAN operations | Commercially emerging (already deployed in limited form today) |
| "AI-native" RAN as a foundational architecture | Active research |
| Integrated Sensing and Communication (ISAC) | Active research |
| Basic satellite SOS/messaging via NTN on ordinary phones | Commercially emerging / Standardized (5G-Advanced) |
| Full seamless broadband-speed non-terrestrial integration | Active research |
| Holographic/volumetric telepresence communication | Speculative |
| Extreme population-scale immersive XR | Speculative |
| Fully autonomous "zero-touch" network operations | Speculative (extending an active research thread) |

---

## 11. A History Lesson in Hype: What 5G's Marketing Taught Us

Chapter 92, Section 12 already showed a concrete, recent, verifiable example worth remembering here: 5G was marketed for years around "sub-1ms latency" and "10 Gbps everywhere," while the honest, typical real-world experience — genuinely good, genuinely a real improvement over 4G, but categorically more modest than the headline figures — took years longer to materialize even partially, and several of the most ambitious capabilities (true network slicing at consumer scale, ubiquitous URLLC) remain narrowly deployed rather than mainstream even as this chapter is being written, roughly half a decade after 5G's initial commercial launch.

This isn't cited to be cynical about 6G, or to suggest none of Sections 5-9's research will pan out — some of it certainly will, in some form, on some timeline. It's cited because **the same institutional incentives that produced 5G's gap between marketing and delivered reality are already visibly present in 6G's earliest public discussion**, years before any specification even exists. Treating 6G's most ambitious claims with the same healthy, evidence-seeking skepticism this course has tried to model throughout — and that Chapter 92 explicitly applied to 5G's own marketing — is simply applying a lesson this course already taught you, one generation earlier, to a case where the stakes for getting it right are, if anything, higher, precisely because there's so much less concrete, checkable reality to compare the claims against yet.

---

## 12. How to Evaluate Any Future-Tech Claim Yourself

A practical, reusable checklist, distilled from this entire chapter's approach, for evaluating any "next-generation networking" claim you encounter after finishing this course:

1. **Is there a finalized, ratified standard, or a specific numbered specification/RFC/3GPP release behind this claim?** If not, it is not yet a technology you can build a dependent product on with confidence — it's a proposal.
2. **Is the number being quoted a typical, real-world, measured figure, or a theoretical peak under specifically favorable, disclosed conditions?** Chapter 92, Section 12 already showed you exactly this gap for 5G; assume the same gap exists until you've seen independently verified, typical-case numbers.
3. **Who is making the claim, and what are they selling?** A peer-reviewed research paper with disclosed experimental conditions, a standards body's published specification, and a vendor's product marketing page are three very different levels of evidentiary weight — treat them accordingly, and don't let a marketing page borrow the credibility of a genuine research result it's loosely citing.
4. **Does the claim describe something physically demonstrated at small scale, or extrapolated to a scale nobody has actually tested?** Section 5's terahertz research includes real laboratory demonstrations; that is meaningfully different evidence than an extrapolated claim about nationwide terahertz coverage, even if both get called "6G research."
5. **What specifically would have to be additionally solved to go from "demonstrated in a lab/pilot" to "deployed at scale"?** If you can't answer this, you likely don't yet understand the claim well enough to evaluate it — and that's a completely reasonable, honest place to stop.

**A worked example, applying all five questions to a representative headline:** suppose you read the claim "Researchers achieve 100 Gbps over a terahertz link, paving the way for 6G." Walking it through the checklist:

1. *Is there a finalized standard behind this?* No — this is Section 5's territory, a research demonstration, not a 3GPP or ITU specification.
2. *Typical or theoretical-peak-under-favorable-conditions?* Almost certainly the latter — headline lab results like this are reported under carefully controlled conditions (precise antenna alignment, short distance, controlled atmosphere), not typical real-world deployment conditions.
3. *Who's making the claim?* If it's a peer-reviewed conference paper (common venues: IEEE journals, academic 6G workshops) disclosing its exact test setup, that's meaningfully stronger evidence than a vendor press release citing the same number without the underlying methodology.
4. *Small-scale demonstration or tested-at-scale claim?* A single link between two fixed, precisely-aligned antennas in a lab is a small-scale demonstration — it says essentially nothing yet about whether the same result holds for a moving device, at a realistic distance, through realistic obstacles.
5. *What's the gap to real deployment?* Everything Section 5 already flagged: efficient terahertz-capable semiconductor hardware at consumer-device cost and power levels, workable beamforming robust to a moving target, and a solution to severe atmospheric/obstacle absorption at any meaningful range.

Running through these five questions doesn't tell you the research is unimportant — a genuine 100 Gbps link, even under narrow lab conditions, is real scientific progress and a legitimate input to future standards work. It tells you specifically what the claim does and doesn't establish, which is exactly the skill this chapter is trying to build, and exactly the skill Section 11's 5G hindsight shows was worth having all along.

---

## 13. Common Misconceptions

- **"6G already exists somewhere, just not where I live yet."** As Section 2 stated plainly, no finalized 6G standard or genuine commercial 6G network exists anywhere as of this writing — any product marketed as "6G" today is either misusing the term or referring to something else (enhanced 5G, private branding) entirely.
- **"6G will definitely launch in 2030 because that's what's been announced."** Section 4 showed this is a reasonable extrapolation from historical generational cadence, not a guaranteed commitment from any standards body — treat it as an educated estimate, not a scheduled fact.
- **"AI-native RAN means today's 5G networks are already run entirely by AI."** Section 6 drew a specific, important distinction: AI already assists with optimization in real deployed networks today, which is different from a foundational, ground-up AI-native architecture, which remains active research.
- **"6G health/safety concerns are a new, unstudied risk."** The same non-ionizing-radiation physics Chapter 92, Section 14 already explained for 5G applies to any terrestrial cellular frequency under serious research consideration for 6G as well; higher frequency does not mean ionizing radiation, and this remains true regardless of which generation is being discussed.
- **"Every 6G research direction mentioned in this chapter will make it into the final standard."** Standards processes routinely explore, and then narrow or drop, many candidate technologies between early research and a finalized specification — Section 9's "further-out" ideas, in particular, are exactly the kind of proposals that may or may not survive that narrowing process in any recognizable form.

---

## 14. What's Simplified Here

This chapter's coverage of 6G research is necessarily a snapshot of a fast-moving, still-early field — new research results, standards study-item outcomes, and vendor announcements will continue to appear after this chapter is written, and some specific claims and timelines here may already look dated or be superseded by the time you're reading this. This is an inherent property of writing about pre-standardization technology, not a correctable flaw — which is precisely why Section 12's evaluation framework, meant to outlast any specific fact in this chapter, is arguably more durable and important than any individual research direction described in Sections 5-9.

---

## 15. Interview Questions & Model Answers

**Beginner: "Does 6G exist yet?"**

*Model answer:* "No. As of now, there's no finalized 6G standard from 3GPP or the ITU, and no genuine commercial 6G network anywhere. What exists is early standards study work (like the ITU's IMT-2030 framework) and active academic and industry research into candidate technologies — terahertz spectrum, AI-native radio access, integrated sensing and communication, and better integration with satellite networks. Commercial 6G is widely expected around 2030, based on the roughly decade-long cadence of prior generations, but that's an estimate, not a commitment."

**Intermediate: "What problem is Integrated Sensing and Communication trying to solve, and how mature is it?"**

*Model answer:* "ISAC is about using the same radio signals a cellular network already sends for communication to also sense the physical environment — detecting nearby objects' position and movement using reflections of the communication signal itself, similar in principle to radar, without needing separate dedicated sensing hardware. Potential uses include traffic monitoring and industrial safety systems. It's currently an active research area, with real experimental demonstrations in academic and industry literature, but it isn't standardized in any deployed cellular generation yet, and real-world reliability at the accuracy safety-relevant use cases would need remains an open problem."

**Advanced: "Why should an engineer be skeptical of 6G marketing claims specifically, given that they weren't necessarily as skeptical of past cellular generation claims?"**

*Model answer:* "Because we have a very recent, well-documented precedent: 5G was marketed for years with headline figures like sub-1ms latency and multi-gigabit speeds 'everywhere,' but the typical real-world experience for most users turned out to be meaningfully more modest, and some of 5G's most ambitious architectural promises — genuine network slicing at consumer scale, ubiquitous URLLC — remain narrowly deployed years after commercial launch rather than mainstream. 6G's earliest public discussion is already showing the same pattern of ambitious headline claims preceding any finalized specification, let alone real deployment data. That doesn't mean none of it will happen — some of it likely will, in some form — but an engineer evaluating 6G claims today should apply the same standard applied to 5G in hindsight: ask whether a claim is backed by a finalized standard and independently measured, typical-case data, or whether it's a theoretical peak figure or a research proposal being presented with more confidence than the underlying evidence currently supports."

**Advanced: "A press release claims a 'real terahertz 6G link at 100 Gbps.' Walk through how you'd evaluate it before repeating the claim."**

*Model answer:* "I'd apply the same five checks I'd apply to any pre-standardization technology claim: first, whether there's a finalized standard behind it — for terahertz 6G, there isn't one, so this is research, not deployed technology. Second, whether the number is a typical or a theoretical-peak, favorable-conditions figure — lab terahertz results are almost always the latter, achieved over short distances with precisely aligned antennas, not typical usage conditions. Third, who's making the claim and what they're selling — a peer-reviewed paper disclosing its methodology carries more weight than a press release citing the same number without it. Fourth, whether it's a small-scale demonstration or something actually tested at real deployment scale — a two-antenna lab link says very little yet about a moving device in a real environment. And fifth, what's still missing to go from this demonstration to real deployment — for terahertz specifically, that's efficient, low-cost consumer-grade hardware, workable beamforming for moving targets, and a solution to severe atmospheric and obstacle absorption. None of that means the research isn't valuable — it's a legitimate, real scientific result — but it means I would describe it precisely as 'a promising early research result,' not as evidence that '6G will deliver 100 Gbps.'"

---

## 16. Exercises

### Easy

1. In one sentence, state the current status of 6G as of this chapter's writing.
2. Name two of this chapter's five status labels, and give one 6G-related example of each from Section 10's table.
3. What does ISAC stand for, and what is its basic idea?

### Medium

4. Using Section 3, explain in your own words why 6G research is substantially motivated by unfinished 5G ambitions rather than being a completely fresh set of goals.
5. Explain the distinction Section 6 draws between "AI assisting an existing RAN" and "AI-native RAN as a foundational architecture." Why does conflating the two lead to an inaccurate picture of how mature AI-driven networking actually is today?
6. Using Section 5's physics-based explanation, describe why terahertz spectrum's enormous available bandwidth comes with a correspondingly severe practical cost, and name at least two specific real-world conditions that would degrade a terahertz link's usable range.

### Hard

7. Section 11 draws a direct parallel between 5G's marketing-versus-reality gap and early 6G discussion. Pick one specific claim from Sections 5-9 of this chapter, and apply Section 12's five-point evaluation checklist to it explicitly, point by point, explaining what additional evidence would change your confidence in it.
8. Section 8 noted that basic satellite messaging via NTN is already commercially shipping, while full broadband-speed non-terrestrial integration remains active research. Explain, using what you know from this chapter and Chapter 92's spectrum/latency material, what specific technical gaps likely separate "occasional emergency text via satellite" from "seamless, full-speed cellular data with a satellite handing off to and from terrestrial towers."
9. Section 7 draws a connection between ISAC and Chapter 92's Massive MIMO beamforming. Explain, in your own words, what capability a base station already needs for beamforming that ISAC proposes reusing for a different purpose — and what new problem (beyond simply locating an intentional radio transmitter) ISAC has to solve that beamforming alone does not.

---

## 17. Summary

| Term | Meaning |
|---|---|
| 6G | The next cellular generation, currently in early research and standardization; no finalized standard or commercial deployment exists as of this writing |
| IMT-2030 | The ITU's high-level requirements framework for 6G-era systems |
| Terahertz spectrum | Frequencies above mmWave (100 GHz+) offering enormous bandwidth at the cost of extreme propagation loss — active research |
| AI-native RAN | The vision of a radio access network built around learned models as foundational architecture, distinct from today's real but narrower AI-assisted RAN optimization |
| ISAC | Integrated Sensing and Communication — using communication signals to also sense the physical environment, radar-like — active research |
| NTN | Non-Terrestrial Networks — integrating satellite connectivity into standard cellular architecture; basic messaging already commercially emerging, full broadband integration still active research |
| Five-level honesty scale | Deployed / commercially emerging / standardized / active research / speculative — the labeling discipline this chapter (and Chapter 92) applied to every claim |
| Hype-vs-reality gap | The documented pattern, seen clearly in 5G, of marketing headline figures outrunning typical real-world delivered performance for years after a generation's launch |

This chapter, and this volume, close with an honest picture of the frontier rather than a confident prediction — which is itself the right note to end a course volume about mobile networks on. Every generation this volume covered, from 1G's analog voice to 5G's software-defined, sliceable core, was built to move data across the hardest possible physical environment: open air, moving devices, no fixed wire at all. Chapter 94 now turns to the opposite extreme of that same problem — networking inside a single building, under one roof, with as much fixed, controlled infrastructure as an engineering team could possibly want — and opens Part 15 with a look inside a modern data center: the servers, NICs, top-of-rack switches, and leaf-spine architecture that give every machine in a hyperscale facility roughly equal bandwidth to every other machine.
