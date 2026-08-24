# Chapter 67: Product Management — Building Things People Actually Want

*Product management is one of the most important and misunderstood roles in tech companies. The PM sits at the intersection of customer needs, business goals, and technical constraints — deciding what to build and why.*

---

## What Is Product Management?

**The classic definition:** A product manager is the CEO of the product. They own the "what" (what to build) while engineering owns the "how" (how to build it).

**The honest definition:** A PM is responsible for:
1. Understanding what customers need
2. Deciding what to build that serves both customer needs and business goals
3. Communicating priorities to engineering, design, and other teams
4. Measuring whether what was built is actually working

What PMs don't do: write code (usually), design (usually), or manage engineers directly.

What makes it hard: PMs have significant responsibility but usually no direct authority over the people who build the product.

---

## The Product Triangle — Customer, Business, Technology

Every product decision sits at the intersection of three constraints:

**Customer (Desirability):** Does the customer want this? Will they use it? Will it solve their problem?

**Business (Viability):** Does this make business sense? Does it generate revenue or reduce costs? Does it fit the company's strategy?

**Technology (Feasibility):** Can we build this with our current team and resources? How long will it take?

A good product decision is desirable (customers want it), viable (the business benefits), and feasible (the team can build it in reasonable time).

**The tension:** These three often conflict.
- Customers want features that are expensive to build (customer vs tech)
- Business wants features customers don't want (business vs customer)
- The most technically elegant solution isn't what customers need (tech vs customer)

The PM's job is to find decisions where all three align — and make difficult tradeoffs when they don't.

---

## Product Discovery — Finding What to Build

Most PMs spend too much time in delivery (building things) and not enough in discovery (figuring out what to build).

**The discovery process:**

**1. Define the problem (not the solution)**
Weak: "We should build a search feature."
Strong: "Users can't find items they're looking for when our catalog exceeds 5,000 products — 23% of users who use search abandon within 2 minutes."

Starting with the problem prevents building solutions to the wrong problem.

**2. Talk to customers**
Not just surveys. Real conversations. The goal: understand their mental model, their workflow, their pain points.

The Mom Test rules (from Rob Fitzpatrick):
- Ask about their life, not your product
- Ask about past behavior, not future intentions ("Have you ever tried to do X?" not "Would you use X?")
- Ask for specifics, not generalities

**3. Define success metrics before building**
What number will change if this feature succeeds? If you can't define this in advance, you won't know if you succeeded.

**4. Explore multiple solutions**
Before building, explore 3-5 different ways to solve the problem. The first idea is rarely the best idea.

**5. Validate with minimal effort**
Build the cheapest possible test of your hypothesis. A mockup, a prototype, a manual experiment — before engineering writes a single line of code.

---

## Prioritization — Deciding What to Build First

The hardest PM skill. You'll always have 10x more ideas than capacity to build.

**Prioritization frameworks:**

**RICE Scoring:**
- **R**each: How many users are affected per time period?
- **I**mpact: How much does it affect users? (1=minimal, 3=massive)
- **C**onfidence: How sure are you about reach and impact? (as %)
- **E**ffort: Engineering weeks required

Score = (Reach × Impact × Confidence) ÷ Effort

Higher RICE score = higher priority.

**ICE Scoring (simpler):**
- **I**mpact: How big is the potential impact?
- **C**onfidence: How confident are we?
- **E**ase: How easy to implement?

ICE = Impact × Confidence × Ease

**Jobs-To-Be-Done (JTBD) Framework:**
Instead of building features, focus on "jobs" — things customers are trying to accomplish.

The famous example: "People don't want a drill. They want a hole in the wall." Even more accurately: "They want a picture hung." Even more accurately: "They want to feel proud of their home."

When you understand the underlying job, you might find a completely different solution.

---

## Writing Good Product Requirements

The PM writes the requirements (PRD — Product Requirements Document) that communicates what needs to be built to engineering and design.

**A good PRD includes:**
1. **Context:** Why are we building this? What problem does it solve?
2. **Goals:** What outcome are we optimizing for? What metrics will we measure?
3. **Non-goals:** What are we explicitly NOT doing? (This prevents scope creep)
4. **User stories:** As a [user], I want to [action] so that [benefit]
5. **Functional requirements:** Specific behaviors the product must have
6. **Edge cases:** What happens in unusual situations?
7. **Success metrics:** How do we know this worked?

**The most important part:** The "why." If engineers understand why they're building something, they make better micro-decisions. If they only know the "what," every ambiguous edge case requires a PM decision.

---

## Working with Engineering — The Most Critical Relationship

The PM-engineering relationship determines whether great products get built.

**What engineers want from PMs:**
- Clear problem statements (not prescribed solutions)
- Understanding of technical constraints (don't ask for what's architecturally impossible in the timeline)
- Respect for engineering judgment on implementation
- Clear priorities (not "everything is P0")
- Stability — not changing requirements every day

**What good PMs give engineers:**
- Context: why are we building this?
- Clear acceptance criteria: what does "done" look like?
- Quick decisions: unblocking engineers fast when they have questions
- Protection from stakeholder noise: shielding engineering from constant direction changes

**The PM trap:** Designing the solution in the PRD rather than the problem. Engineers often have better ideas for implementation than PMs — but PMs sometimes get attached to their vision. "Falling in love with the solution" is a common PM failure mode.

---

## Measuring Product Success

**The North Star Metric:**
One metric that captures the core value you deliver to customers. All team activity should ultimately move this.

- Spotify: time spent listening
- Airbnb: nights booked
- Slack: messages sent per user per day
- Zomato: orders per week per user

The North Star metric is a lagging indicator (it measures output). PMs also track leading indicators (things that predict the North Star).

**The AARRR metrics (Pirate Metrics):**
- Acquisition: how are users finding you?
- Activation: are new users having a "first value" experience?
- Retention: are users coming back?
- Revenue: are users paying?
- Referral: are users telling others?

**What to measure by product stage:**

| Stage | Primary Metric |
|-------|---------------|
| Pre-launch | Waitlist/interest signals |
| Early stage | Activation rate (% of signups who get value) |
| Growing | Retention (% of users active 30 days after signup) |
| Scaling | Revenue per user, LTV:CAC |
| Mature | Net Revenue Retention, NPS |

---

## The PM Career Path

**Associate PM (APM):** Junior PM, often recent graduate. Typically works on one product area with mentorship.

**PM:** Responsible for a specific product or feature set. Owns PRDs, discovery, measurement.

**Senior PM:** Leads a product area with multiple features. Influences strategy. Mentors junior PMs.

**Principal PM / Group PM:** Strategy, cross-team product thinking. Often drives major product bets.

**VP of Product / Chief Product Officer (CPO):** Sets product direction for the company. Works with CEO on long-term roadmap.

**The builder alternative:** Many experienced PMs become founders — they've spent years learning what customers want and how to build products. This is the natural entrepreneurial path.

---

## Product Management in India

**The Indian PM market:**
- Significant demand at Indian tech companies (Zomato, Flipkart, Razorpay, Zepto)
- Even stronger demand for PMs who can work on global products
- Indian PMs are increasingly hired by US tech companies remotely

**What makes a great Indian PM:**
- Deep understanding of the Indian user (non-English, price-sensitive, mobile-first)
- Data fluency (SQL, analytics tools)
- Communication clarity in writing (PRDs, strategy docs)
- Ability to work with large, diverse engineering teams

**Salary ranges (2024):**
- APM at startup: ₹15-25 lakh
- PM at Flipkart/Zomato/Razorpay: ₹35-60 lakh
- Senior PM: ₹60-1.5 crore
- PM at US tech (Google, Meta, Amazon, remote): $150,000-$300,000+

---

## Key Lessons

1. **Start with the problem, not the solution** — falling in love with your solution before understanding the problem is the most common PM mistake
2. **Talk to users constantly** — not surveys, real conversations
3. **Define success metrics before building** — if you can't measure it, you can't improve it
4. **Prioritization is the hardest skill** — every team thinks their request is most important
5. **The PM-engineer relationship determines product quality** — invest in it

---

## Summary

- PMs own the "what" — deciding what to build — while engineering owns the "how"
- Discovery (understanding what customers need) is as important as delivery (building it)
- RICE/ICE scoring provides systematic prioritization frameworks
- North Star Metric + AARRR funnel is the core measurement framework
- India's PM market is strong; experienced PMs increasingly work globally

---

## Exercises

1. **Write a product requirements document** for one feature you want in an app you use. Include problem statement, goals, user stories, and success metrics.
2. **Do 3 customer interviews** using Mom Test principles. What did you learn that surprised you?
3. **Prioritize 5 features** using RICE scoring for a hypothetical product. Which scores highest? Does it match your intuition?
4. **Find the North Star Metric** for three apps you use. How do you know that's the right metric?
5. **Analyze a product failure.** Find a feature that was shipped but removed or ignored (Google+, Facebook Poke, Snapchat's redesign). What went wrong in the discovery or prioritization?
