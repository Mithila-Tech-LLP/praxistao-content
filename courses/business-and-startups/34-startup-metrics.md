# Chapter 34: Startup Metrics — What to Measure and Why It Matters

*Not everything that counts can be counted, and not everything that can be counted counts. But in business, the right numbers tell the truth about your company.*

---

## 1. Vanity Metrics vs Real Metrics

A vanity metric is a number that looks impressive but doesn't tell you whether your business is working.

| Vanity Metric | Why It's Misleading | Better Metric |
|---------------|---------------------|---------------|
| Total registered users | Anyone can sign up; most never use the product | DAU / MAU, retention |
| App downloads | Downloads ≠ active users | Day 7 retention |
| Total GMV | Gross value doesn't reflect what you actually earn | Net revenue, take rate |
| Social media followers | Followers don't pay bills | Engagement rate, conversions |
| Funding raised | You can burn money efficiently or not | Burn multiple (revenue per $ burned) |
| Press mentions | PR doesn't mean customers are buying | Revenue, user growth |

Vanity metrics are dangerous because they feel like progress. You can have 10 million app downloads and be 6 months from bankruptcy.

**Rule:** Every metric you track should lead to a decision. If the metric went up and you'd do nothing differently, it's not a useful metric.

---

## 2. The Metrics That Matter at Each Stage

### Pre-PMF (Stages 1-2)

You're not focused on growth yet. You're trying to validate.

- **Active users per day/week:** Are people using the product?
- **Retention (D7, D30):** Do users come back?
- **Sean Ellis score:** What % say they'd be "very disappointed" if the product disappeared?
- **Qualitative NPS + user interviews:** Why do users love it? Why do they leave?

Don't worry about revenue, CAC, or LTV yet. If nobody is using the product, these metrics are meaningless.

### Post-PMF, Pre-Growth (Stage 3)

You've found the signal. Now validate it can be replicated.

- **MoM (Month-over-Month) growth rate:** Growing 10-30% monthly is the target
- **Retention cohorts:** Does the D30 retention curve flatten?
- **First acquisition channel CAC:** What does it cost to acquire one customer through your primary channel?
- **Early unit economics:** Revenue per user vs. cost per user

### Growth Stage (Stages 4-5)

You know the model works. Now scale and optimize.

- **Revenue growth rate** (MoM, YoY)
- **CAC by channel**
- **LTV**
- **LTV:CAC ratio** (must be >3x for sustainable economics)
- **Churn rate**
- **Net Revenue Retention (NRR)** for B2B
- **Burn multiple** (cash burned per $1 of net new ARR added)

---

## 3. Consumer App Metrics

### DAU, WAU, MAU

- **DAU (Daily Active Users):** Users who open the app at least once in a day
- **WAU (Weekly Active Users):** At least once in a week
- **MAU (Monthly Active Users):** At least once in a month

### Stickiness

**Stickiness = DAU / MAU**

If your DAU is 1 million and MAU is 10 million, stickiness = 10%.

- WhatsApp: ~70% (people use it almost every day)
- Most social apps: 20-40%
- News apps: 10-20%
- Transactional apps (travel, insurance): 5-15%

Higher stickiness means higher engagement and better retention.

### Retention Curves

Plot the % of users who return on each day after signup. A healthy retention curve looks like a swoosh — drops sharply at first, then flattens.

```
100%|
 80%|  \
 60%|    \
 40%|      \_______   ← Flat tail = retained core users
 20%|               (say 30% come back every month)
  0%+————————————————————→ Days
    1  7  14  30  60  90
```

A curve that goes to zero means no one finds enduring value. This is a product problem, not a marketing problem.

### Session Metrics

- **Session length:** How long is a typical session?
- **Sessions per user per day:** How often do users open the app?
- **Actions per session:** What do they do when they're in it?

For a game: longer sessions and more sessions per day = good. For a utility (task manager, navigation): shorter sessions that accomplish the task = good. It depends on the product's job.

---

## 4. SaaS Metrics — The Complete Set

SaaS is unique because revenue is recurring. This creates a specific set of metrics that capture the health of the business.

### MRR and ARR

**MRR (Monthly Recurring Revenue):** The predictable revenue your SaaS earns each month. A customer paying ₹10,000/month contributes ₹10,000 MRR.

**ARR (Annual Recurring Revenue):** MRR × 12. Investors use ARR to evaluate SaaS companies.

**MRR Movements:**
- **New MRR:** Revenue from brand new customers
- **Expansion MRR:** Revenue from existing customers who upgraded
- **Churned MRR:** Revenue lost from customers who cancelled
- **Net New MRR** = New MRR + Expansion MRR - Churned MRR

A healthy SaaS has Expansion MRR > Churned MRR. This means existing customers are growing with you, which offsets cancellations.

### Churn Rate

**Logo churn:** % of customers who cancel in a given period
**Revenue churn:** % of revenue lost due to cancellations

Monthly churn benchmarks:
- <1%: Excellent
- 1-2%: Good
- 2-3%: Needs work
- >3%: Urgent problem

Why churn kills SaaS: at 3% monthly churn, you lose 30% of your customer base per year. You must acquire 30% new customers just to stay flat.

### NRR (Net Revenue Retention)

**NRR = (Starting MRR + Expansion MRR - Churn MRR) / Starting MRR × 100%**

NRR tells you: if you stopped acquiring new customers today, would your revenue grow or shrink?

- NRR > 100%: Revenue grows from existing customers alone (expansion > churn). This is the holy grail.
- NRR = 100%: Flat
- NRR < 100%: You're shrinking from existing customers — a death spiral

**World-class NRR benchmarks:**
- Snowflake: ~130-140%
- Twilio: ~120-130%
- Strong Indian SaaS: 110-120%

NRR > 100% means the company could theoretically stop acquiring new customers and still grow. That's an incredibly powerful economic property.

### CAC (Customer Acquisition Cost)

**CAC = Total sales & marketing spend / Number of new customers acquired**

If you spent ₹50 lakh on marketing and sales this quarter and acquired 100 new customers, your CAC = ₹50,000.

**Payback period** = CAC / (Monthly Revenue per Customer × Gross Margin %)

If CAC is ₹50,000 and a customer pays ₹5,000/month at 80% gross margin, payback = ₹50,000 / (₹5,000 × 0.80) = 12.5 months.

Best-in-class SaaS: 12-18 month payback period. Startups in growth mode: 18-24 months is acceptable.

### LTV (Lifetime Value)

**LTV = (Average Revenue per Customer per Month × Gross Margin %) / Monthly Churn Rate**

If customers pay ₹5,000/month, gross margin is 80%, and monthly churn is 2%:

LTV = (₹5,000 × 0.80) / 0.02 = ₹2,00,000

### The Magic Number: LTV:CAC Ratio

If LTV is ₹2,00,000 and CAC is ₹50,000:

LTV:CAC = 4:1

**Benchmarks:**
- <1: Losing money on every customer (don't scale this)
- 1-3: Marginal; needs improvement
- >3: Good; unit economics support scaling
- >5: Excellent; scale aggressively

---

## 5. E-commerce Metrics

### GMV vs Revenue vs Net Revenue

**GMV (Gross Merchandise Value):** Total value of all goods sold on the platform. Amazon India's GMV includes every product sold, including by third-party sellers.

**Revenue (Take Rate):** GMV × Take Rate. If GMV is ₹1,000 crore and take rate is 5%, Revenue = ₹50 crore.

**Net Revenue:** Revenue minus direct costs (logistics subsidies, payment fees, returns).

Investors in e-commerce look at GMV for market position, but ultimately care about revenue and net revenue for business viability.

### AOV (Average Order Value)

AOV = Total Revenue / Number of Orders

Higher AOV = more revenue per delivery. E-commerce companies work hard to increase AOV through:
- Minimum order for free delivery
- Bundle recommendations
- Upsells at checkout

### Repeat Purchase Rate

For a grocery app, what % of customers who ordered in month 1 ordered again in month 2? Month 3?

High repeat rate = the product is genuinely useful. Low repeat rate = either the product disappointed or the need isn't frequent.

---

## 6. Marketplace Metrics

Marketplaces have two-sided metrics: supply (sellers/providers) and demand (buyers).

**Supply side:**
- Number of active sellers/providers
- Provider utilization rate (what % of their capacity is being used?)
- Provider NPS (are they happy with the platform?)

**Demand side:**
- Number of active buyers
- Purchase frequency
- Buyer NPS

**Core marketplace metrics:**
- **Take Rate:** Revenue as a % of GMV. Higher is better for economics; too high and sellers leave.
- **Liquidity:** % of searches/requests that result in a match. Low liquidity = marketplace isn't working.
- **Repeat rate:** % of buyers who return within 90 days.

---

## 7. Burn Rate and Runway

**Monthly Burn Rate:** Cash spent per month (net outflow). If you spend ₹2 crore/month and earn ₹50 lakh/month in revenue, your net burn is ₹1.5 crore/month.

**Runway:** How many months until you run out of cash.

Runway = Cash in Bank / Monthly Net Burn

If you have ₹18 crore in the bank and burn ₹1.5 crore/month: Runway = 12 months.

**Why runway matters:** You should always know your runway. When it drops below 6 months, you should already be in fundraising conversations. When it drops below 3 months, you're in survival mode.

**Rule of thumb:** Start your next fundraise when you have 9-12 months of runway. It takes 3-6 months to close a round.

**Burn multiple:** Net burn / Net new ARR. How much money do you burn for every $1 of new recurring revenue?
- <0.5: Efficient growth
- 0.5-1.5: Acceptable
- >2: Expensive growth; needs improvement

---

## 8. Building a Dashboard

For investors, your dashboard shows: the metrics that prove the business is working.

For internal use, your dashboard shows: the metrics that help you run the business.

These are often different.

**Investor dashboard (monthly):**
- Revenue / ARR
- MoM / YoY growth rate
- Gross margin
- Key unit economics (CAC, LTV, NRR)
- Burn and runway
- Headcount

**Operational dashboard (weekly/daily):**
- North Star Metric
- Daily/weekly active users
- Acquisition by channel
- Funnel conversion rates
- Support tickets/escalations
- Revenue

Don't build a 50-metric dashboard. If you track everything, you focus on nothing. Pick the 5-10 metrics most relevant to your stage and focus on them relentlessly.

---

## Summary

- **Vanity metrics** look good, don't lead to decisions. Track metrics that tell you whether the business is working.
- Consumer apps: focus on retention curves, stickiness, North Star Metric
- SaaS: MRR, churn, NRR (>100% is magic), LTV:CAC >3x
- E-commerce: GMV vs Revenue, AOV, repeat rate
- Always know your **burn rate and runway** — this is not optional
- Build two dashboards: investor-facing (monthly) and operational (weekly) — keep both lean

---

## Exercises

1. **Pick a public Indian tech company** (Zomato, Paytm, Info Edge). Read their quarterly earnings report. What metrics do they report? What's not in the report that you wish was?

2. **Build a unit economics model** for a hypothetical SaaS with: ₹3,000/month price, 75% gross margin, 3% monthly churn, ₹30,000 CAC. Calculate LTV, payback period, LTV:CAC.

3. **Define the North Star Metric** for 5 different startups. Defend each choice.

4. **Calculate the runway** for a startup with ₹5 crore in the bank, ₹80 lakh/month in revenue, and ₹2 crore/month in expenses.

5. **Find a startup that reported impressive metrics before failing.** Were the metrics vanity metrics? What real metrics would have predicted the failure?
