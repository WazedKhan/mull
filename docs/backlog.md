# Mull: Backlog (V1)

Tickets follow the spec's build order. Journeys and requirements (R-numbers) are in `docs/user-journeys.md`. **[Suggested]** marks tickets that go beyond the spec.

Labels: `backend`, `frontend`, `engine`, `ai`, `infra`, `ux`, `suggested`.
Size: S (under half a day), M (1 to 2 days), L (3+ days).

## Milestones

| Milestone | Goal | Tickets |
| --- | --- | --- |
| M0 Foundation | Repo runs locally with one command | 01 |
| M1 Rules engine | Verdicts correct and tested, no UI | 02, 03, 04 |
| M2 Data and API | All money data stored and served | 05, 06, 07, 08, 09, 10, 14, 16, 17 |
| M3 Core screens | Setup, Dashboard, Buy check usable on a phone | 11, 12, 13, 15 |
| M4 AI layer | One-sentence explanations with hard guardrails | 19, 20 |
| M5 History and month of use | History, summary, then a full month of real use | 18, 21, 22, 23 |
| Later | After the month of use | 24, 25, 26 |

## M0 Foundation

### MULL-01 Project scaffolding
`infra` · S
Go (Echo) API, React (Vite) frontend, PostgreSQL, all in `docker-compose`. Migrations tool (golang-migrate or goose). Makefile targets: `run`, `test`, `migrate`. CI runs `go test ./...` and frontend build on push.
- [ ] `make run` brings up API, web and DB
- [ ] CI is green on an empty test
- [ ] README updated with setup steps

## M1 Rules engine

### MULL-02 Rules engine core with trace
`engine` `backend` · M · depends on 01
Pure Go package, no DB or HTTP. Input: item (price, need/want), month state, goals, fixed costs, thresholds. Output: verdict, rule id, the one number, and a trace of inputs used. First matching rule wins, in spec order (Skip 1, Skip 2, Wait 3, Wait 4, Buy 5, Buy 6).
- [ ] All six rules implemented in order
- [ ] Thresholds passed in (rule 1 cover months, rule 2 multiplier), not hard-coded
- [ ] Returns a trace struct (R10.1)
- [ ] Handles zero fixed costs without dividing by zero

### MULL-03 Affordable month and price ceiling
`engine` · S · depends on 02
For Wait verdicts, compute the earliest month the price fits without cutting any goal contribution (R4.1). For a need over budget, also return a price ceiling equal to `remaining_free` (R6.2).
- [ ] Month counted forward from current month using `monthly_free`
- [ ] Returns "not within 12 months" instead of a far-off date

### MULL-04 Engine tests from the spec's sample cases
`engine` · S · depends on 02, 03
Table-driven tests: the three phone cases from the spec, plus edge cases for each rule boundary (price equal to free money, cover exactly 1 month, funded wish goal).
- [ ] Every rule has at least one firing case and one non-firing case
- [ ] Spec sample cases produce the spec's verdicts

## M2 Data and API

### MULL-05 Schema and money type
`backend` · M · depends on 01
Tables: `income_defaults`, `fixed_costs`, `goals`, `months`, `month_income`, `spends`, `checks`. All money as `BIGINT` BDT (R1.1). Use a `Money` type in Go.
- [ ] Migrations up and down work
- [ ] No float money anywhere in Go structs or SQL

### MULL-06 Derived numbers service
`backend` · S · depends on 05
`monthly_free = income - fixed costs - monthly contributions`, `remaining_free = monthly_free - spent this month`, emergency cover in months, suggested emergency target (3x or 6x) (R1.2, gap G4).
- [ ] Unit tests for each formula
- [ ] Names `monthly_free` and `remaining_free` used consistently in API responses

### MULL-07 Goals model and on-track status
`backend` · M · depends on 05
CRUD for goals: name, type (DPS, emergency, custom, wish), target, saved so far, monthly contribution, target date, status (active, paused, done). On-track = saved so far at least the expected amount by now (R2.2).
- [ ] REST endpoints with validation
- [ ] On-track flag computed, not stored

### MULL-08 Month state and income overrides
`backend` · M · depends on 05
A month starts on pay day. Income defaults each month and can be overridden for one month; one-off extra income adds to that month only (R7.1, R7.2).
- [ ] Current month created lazily on first request after pay day
- [ ] Override and extra income do not change the default

### MULL-09 Check history persistence
`backend` · S · depends on 05
Store every check: inputs, rule verdict, rule id, final verdict, sentence, number, outcome (pending, bought, bought anyway, skipped, waited) (R3.2).
- [ ] Outcome can be updated after the check

### MULL-10 Buy check endpoint
`backend` · M · depends on 02, 06, 07, 08, 09
`POST /checks`: loads current state, runs the engine, calls the AI layer (stubbed until M4), saves and returns the result.
- [ ] Rules-only path responds in under 100 ms locally (R3.1)
- [ ] Returns the trace for the "Why?" panel

### MULL-14 Spend logging (gap G1)
`backend` `frontend` · M · depends on 08, 09
"I bought it" and "I bought it anyway" on a result record a spend linked to the check. A quick "Log a spend" action records spends that did not go through a check, marked need or want (R3.3, R5.1).
- [ ] `remaining_free` and wants spent update immediately
- [ ] Spends can be edited or deleted for typos

### MULL-16 Save for this and funded wish goals (gap G5)
`backend` `frontend` · M · depends on 07, 10
On Wait, "Save for this" opens a pre-filled goal form (price as target, suggested contribution that fits). Multiple wish goals allowed. When a wish goal covers the price, the check returns Buy and buying marks the goal done (R4.2, R4.3).
- [ ] Goal is only created after user confirms
- [ ] Contribution that would push `monthly_free` below zero is blocked with a message

### MULL-17 Monthly contribution confirmation (gap G7)
`backend` `frontend` · S · depends on 07, 08
At month start, list planned contributions with "paid" checkboxes. Paid adds to saved so far; unpaid leaves the goal behind (R7.3).
- [ ] Dashboard shows a banner until confirmed

## M3 Core screens

### MULL-11 Setup screen
`frontend` `ux` · M · depends on 06, 07, 08
Three steps: income and pay day, fixed costs, goals. Under 5 minutes. Editable later from Settings. Warn when `monthly_free` is zero or negative (R1.3, R1.4).
- [ ] Timed run-through under 5 minutes
- [ ] Works on a phone screen

### MULL-12 Dashboard
`frontend` · M · depends on 06, 07
`remaining_free`, wants spent this month, emergency cover, goal bars with on-track status, big "Can I buy something?" button (R2.1, R2.3).
- [ ] Loads in under 1 s on a phone

### MULL-13 Goals screen
`frontend` · S · depends on 07
Add, edit, pause, close goals. Shows new `monthly_free` before saving a contribution change. Manual "add money" to a goal (R8.1 to R8.3).

### MULL-15 Buy check screen (gaps G2, G3)
`frontend` `ux` · M · depends on 10
Required: name, price, Need/Want toggle. Optional: what is wrong with what you have, why you need it, category (default Other). Result appears on the same screen with verdict, sentence, number, and the outcome buttons.
- [ ] Three taps plus typing to a verdict
- [ ] Price input uses a numeric keyboard on mobile

## M4 AI layer

### MULL-19 AI explanation with guardrail
`ai` `backend` · M · depends on 10
Send item, reason, condition, rule verdict and key numbers (no identity data, NFR4). The model returns a verdict adjustment and one sentence as JSON. Code enforces: at most one step from the rule verdict, never to Buy unless money fits; otherwise the rule verdict stands (R6.1).
- [ ] Guard unit-tested with an "AI tries to upgrade to Buy" case
- [ ] Sentence trimmed to one sentence, English only
- [ ] Prompt and model name in config

### MULL-20 LLM timeout and fallback (gap G6)
`ai` `backend` · S · depends on 19
4 s timeout. On error or timeout, return the rule verdict with a templated sentence per rule. Log failures.
- [ ] End-to-end check stays under 10 s with the LLM down

## M5 History and a month of use

### MULL-18 History and monthly summary
`frontend` `backend` · M · depends on 09, 14
List of checks filterable by month and verdict. Summary: total checks, verdict split, "followed X of Y Wait/Skip", wants spent vs last month (R9.1, R9.2).

### MULL-21 Why panel and configurable thresholds [Suggested]
`frontend` `backend` `suggested` · S · depends on 02, 10
"Why?" on a result shows the rule that fired and its numbers. Settings let you change the rule 1 and rule 2 cutoffs (R10.1, R10.2).

### MULL-22 CSV export of history [Suggested]
`backend` `suggested` · S · depends on 09
Download all checks and spends as CSV (R9.3).

### MULL-23 Auth, backups and deploy
`infra` · M · depends on 01
Single-user password login with a session cookie, HTTPS behind Nginx, daily `pg_dump` with 14-day retention, deploy to a small VPS (NFR2, NFR3).
- [ ] Restore from backup tested once

**Then: use it yourself for a full month before starting Later tickets.**

## Later

### MULL-24 Ready to re-check [Suggested]
`suggested` · S
When a Wait item's affordable month arrives, show it on the dashboard as "ready to re-check" (R4.4).

### MULL-25 48-hour cool-off reminder [Suggested]
`suggested` · M
After Wait or Skip on a want, optional reminder in 48 hours. Depends on notifications, which the spec puts in Later (R5.3).

### MULL-26 Emergency fund use for needs [Suggested]
`suggested` · S
Allow a need to be paid from the emergency fund with a note; dashboard cover drops and the fund goal shows a refill plan (R6.3).
