# Mull: User Journeys and Requirements (V1)

Source: Savings Coach App Product Spec (Oct 4, 2026). This doc turns the spec into journeys, each with the requirements it needs. Items marked **[Suggested]** are improvements not in the spec; keep or drop them.

Ticket IDs (MULL-xx) point to `docs/backlog.md`.

## Who and what

- **User:** one person (you), self-hosted, BDT only, English verdicts.
- **Promise:** ask before you buy, get Buy / Wait / Skip in under 10 seconds with one sentence of reasoning.
- **Success after 3 months:** less spent on wants, DPS and emergency fund on target, most Wait/Skip verdicts followed.

## Gaps found in the spec

These need a decision before or during build. Each has a proposed answer.

| # | Gap | Why it matters | Proposed answer |
| --- | --- | --- | --- |
| G1 | Nothing records actual spending, but "Month state" needs "spent so far" and "wants spent" | Rules 3, 4 and 6 use remaining free money; without spend logging it never goes down | Add "I bought it" on a verdict plus a quick "Log a spend" action (MULL-14) |
| G2 | Rules need need vs want, but the buy check only asks for "category" | Rules 1, 2, 5 and 6 all branch on it | Explicit Need / Want toggle on the buy check; category stays for reporting (MULL-15) |
| G3 | Spec says the buy check has 4 fields, the screen lists 5 (name, price, category, what is wrong, why) | Affects the "keep it fast" goal | Required: name, price, need/want. Optional: what is wrong, why. Category defaults to "Other" |
| G4 | "Free money" (monthly, from the formula) and "remaining free money this month" are used interchangeably | Rule 2 uses the monthly figure, rule 3 the remaining figure | Name them `monthly_free` and `remaining_free` in code and UI |
| G5 | Can a purchase be paid from a goal balance (e.g. the wish item's own goal)? | A fully funded wish goal should give Buy, not Wait | Yes: if the item has a goal and its saved amount covers the price, verdict is Buy and the goal is marked spent (MULL-16) |
| G6 | No behaviour when the LLM call fails or is slow | Breaks the 10 second promise | Fall back to the rule verdict with a templated sentence after a 4 s timeout (MULL-20) |
| G7 | Goal contributions are planned but never confirmed as paid | Goal progress bars drift from reality | Monthly "mark contributions paid" step at pay day (MULL-17) |

## Journey 1: First-time setup

**Goal:** go from nothing to a usable dashboard in under 5 minutes.

1. Open the app, see a short "3 steps" setup.
2. Enter monthly income and pay day.
3. Enter fixed costs (rent, food, transport, bills, family support). Others can be added.
4. Enter goals: DPS (target, saved so far, monthly contribution), emergency fund (saved so far; target suggested as 3 or 6 months of fixed costs).
5. See the dashboard with free money and months of emergency cover already calculated.

**Requirements**

- R1.1 All money stored as integer BDT (no floats). MULL-05
- R1.2 Emergency fund target defaults to 3 x fixed costs, with a 6x option labelled "income varies". MULL-06
- R1.3 Setup can be finished in one sitting and edited later from Settings. MULL-11
- R1.4 **[Suggested]** Show a warning if free money is zero or negative after setup ("your plan already uses all your income"). MULL-11

## Journey 2: Quick glance at the dashboard

**Goal:** know in 2 seconds how much room is left this month.

1. Open the app.
2. See remaining free money this month, months of emergency cover, and a progress bar per goal.
3. Tap "Can I buy something?".

**Requirements**

- R2.1 Dashboard shows `remaining_free`, emergency cover in months (one decimal), goal bars with "on track / behind". MULL-12
- R2.2 A goal is "behind" if saved so far is less than the amount it should have by now given its contribution and start date. MULL-07
- R2.3 **[Suggested]** Show "wants spent this month" next to remaining free money, since that is the habit being tested. MULL-12

## Journey 3: Buy check gives Buy

**Goal:** a want that fits gets a quick yes.

1. Tap "Can I buy something?".
2. Enter name, price, Want, optional reason.
3. Verdict appears on the same screen: **Buy** plus one sentence.
4. Tap "I bought it". Remaining free money and wants spent update.

**Requirements**

- R3.1 Verdict returned in under 10 s end to end; rules alone under 100 ms. MULL-10, MULL-20
- R3.2 Every check is saved to history with inputs, rule verdict, final verdict, sentence, and the number shown. MULL-09
- R3.3 "I bought it" records a spend against the month. MULL-14

## Journey 4: Buy check gives Wait, then Save for this

**Goal:** turn a "not now" into a plan.

1. Buy check returns **Wait** with the month you could afford it (e.g. "January").
2. Tap "Save for this".
3. A wish-item goal is created, pre-filled with the price as target and a suggested monthly contribution that fits free money.
4. The goal appears on the dashboard with its own bar.
5. When the goal is funded, a new buy check on that item returns **Buy** (G5).

**Requirements**

- R4.1 Rule 3 computes the earliest month the price fits without cutting any goal contribution. MULL-03
- R4.2 "Save for this" never creates a goal automatically; user confirms. Several wish goals can run at once. MULL-16
- R4.3 New goal contribution cannot push `monthly_free` below zero; the form says so. MULL-16
- R4.4 **[Suggested]** On the affordable month, the item shows on the dashboard as "ready to re-check". MULL-24

## Journey 5: Buy check gives Skip

**Goal:** a clear no, with the reason, and honest logging if you buy anyway.

1. Buy check returns **Skip** with one sentence (e.g. DPS behind plan).
2. Either leave, or tap "I bought it anyway".
3. History records the purchase as not following the verdict.

**Requirements**

- R5.1 Skip and Wait results show both "OK, skip" and "I bought it anyway". MULL-14
- R5.2 "Bought anyway" counts against the monthly "verdicts followed" figure. MULL-18
- R5.3 **[Suggested]** For a want that got Wait or Skip, offer "Remind me in 48 hours" so the urge can pass (needs notifications, so Later). MULL-25

## Journey 6: Urgent need (broken phone)

**Goal:** a real need is not blocked by want rules, but the money still has to fit.

1. Enter item, price, mark as **Need**, describe what is wrong ("dead, needed for work").
2. Rules: if it fits `remaining_free`, **Buy** (rule 5). If not, **Wait** (rule 3).
3. AI may soften a Skip to Wait if the reason supports a genuine need, and may add a cheaper-option hint ("look under 18,000").
4. AI can never turn Skip or Wait into Buy when the money does not fit.

**Requirements**

- R6.1 Guard in code: final verdict can move at most one step from the rule verdict, and never to Buy unless the rule verdict allows it on money. MULL-19
- R6.2 The "one number" can be either the affordable month or a price ceiling. MULL-03, MULL-19
- R6.3 **[Suggested]** Allow an "emergency fund use" note when a need is paid from the emergency fund, and show cover dropping on the dashboard. MULL-26

## Journey 7: New month, pay day and income changes

**Goal:** the numbers stay correct month to month without redoing setup.

1. On pay day, a new month starts: spent and wants spent reset to zero.
2. Income uses the default; you can override this month's figure or add one-off extra income.
3. Confirm which goal contributions you actually paid; goal balances update.

**Requirements**

- R7.1 Month state is keyed by month; the engine always reads the current month's income. MULL-08
- R7.2 Overrides and extra income apply to one month only. MULL-08
- R7.3 Contribution confirmation updates "saved so far"; missed contributions mark the goal behind. MULL-17

## Journey 8: Manage goals

**Goal:** keep goals accurate as life changes.

1. Open Goals.
2. Add, edit, pause or close a goal (DPS, emergency, custom, wish item).
3. Change monthly contribution or target date and see the effect on free money.

**Requirements**

- R8.1 CRUD for goals with type, target, saved so far, monthly contribution, target date. MULL-07, MULL-13
- R8.2 Editing a contribution shows the new `monthly_free` before saving. MULL-13
- R8.3 **[Suggested]** Add money to a goal manually (bonus, gift) outside the monthly contribution. MULL-13

## Journey 9: Monthly review

**Goal:** see whether the app is changing your habits.

1. Open History.
2. See every check with verdict and outcome (bought, skipped, waited, bought anyway).
3. See the monthly line: "Followed X of Y Wait/Skip verdicts", wants spent vs last month.

**Requirements**

- R9.1 History list, filterable by month and verdict. MULL-18
- R9.2 Monthly summary: checks, verdict split, follow rate, wants spent. MULL-18
- R9.3 **[Suggested]** Export history as CSV so you can analyse it after 3 months. MULL-22

## Journey 10: Tune the rules **[Suggested]**

**Goal:** after a month of use, adjust cutoffs from real data, as the spec asks.

1. On any verdict, tap "Why?" to see which rule fired and the numbers it used.
2. In Settings, change the rule 1 cover threshold (default 1 month) and rule 2 multiplier (default 3 months).

**Requirements**

- R10.1 Rules return a trace: rule id, inputs, thresholds. MULL-02
- R10.2 Thresholds stored as config, not constants. MULL-21

## Non-functional requirements

- NFR1 Self-hosted, no bank linking, no third-party analytics.
- NFR2 Basic auth or a single password login, since the app holds financial data even with one user. MULL-23
- NFR3 Daily PostgreSQL backup. MULL-23
- NFR4 The item, reason and numbers sent to the LLM contain no identity data. MULL-19
- NFR5 UI copy presents verdicts as a spending-habit tool, not financial advice.
- NFR6 Works on a phone browser (most checks happen in a shop).

## Out of scope for V1 (from spec)

Bank or bKash/Nagad linking, multi-user, notifications (Later), auto-ranking the wish list (Later), languages other than English.
