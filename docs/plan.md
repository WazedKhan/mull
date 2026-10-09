# Mull: Build Plan

Each task is one branch and one PR, sized for a review you can finish in one sitting (target under ~300 changed lines, excluding generated code). Tasks run top to bottom; within a phase, later tasks depend on earlier ones unless noted. Ticket IDs (MULL-xx) are in `backlog.md`.

Tick a task in its own PR. "Done when" is the check the AI must run and show.

## Phase 0: Foundation (MULL-01)

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [x] T00 | Docs and rules | Add `CLAUDE.md`, `docs/*`, `.claude/agents/code-reviewer.md`, `.claude/skills/task/SKILL.md`, empty `docs/learnings.md` | Files render on GitHub |
| [x] T01 | Go API skeleton | `go mod init`, `cmd/api`, Echo server, `GET /healthz`, slog, config from env | Handler test for `/healthz` passes |
| [ ] T02 | Make, lint, CI | Makefile (`run`, `test`, `lint`, `check`), golangci-lint config, GitHub Actions running `make check` | CI green on a push |
| [ ] T03 | Postgres and migrations | docker compose Postgres, goose wired, first empty migration, `make migrate` | `make migrate` runs clean twice |
| [ ] T04 | Web skeleton | Vite React TS in `web/`, calls `/healthz`, ESLint, web build added to `make check` | Page shows "API ok" |

## Phase 1: Rules engine (MULL-02, 03, 04)

No DB, no HTTP. This is the heart of the app, so read every line.

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [ ] T05 | Money type | `internal/money`: int64 BDT, parse, format with commas, add/sub/compare | Table tests incl. negative and zero |
| [ ] T06 | Engine types and rule chain | `Input`, `Result`, `Verdict`, `Trace`, `Thresholds`; first-match chain with no rules yet | Test: empty chain returns an error, not Buy |
| [ ] T07 | Buy rules (5, 6) | Need fits `remaining_free`; want fits and goals on track | One firing and one non-firing case each |
| [ ] T08 | Wait rules (3, 4) | Price above `remaining_free`; would cut a contribution | Boundary test: price equal to `remaining_free` |
| [ ] T09 | Skip rules (1, 2) | Cover below threshold for wants; price above N x `monthly_free` for wants | Rule order test: Skip wins over Wait |
| [ ] T10 | Affordable month and price ceiling | Earliest month price fits; ceiling for needs; "not within 12 months" | Tests with a fixed "now" |
| [ ] T11 | Spec cases and edge cases | The three phone cases from the spec, zero fixed costs, zero income | Spec cases give the spec's verdicts |
| [ ] T12 | Funded wish goal (G5) | Wish goal balance covering the price gives Buy | Test with a fully and partly funded goal |

## Phase 2: Data and API (MULL-05 to 10, 14, 16, 17)

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [ ] T13 | sqlc and settings tables | sqlc config, migrations for `income_defaults` and `fixed_costs`, queries | Store test against real Postgres |
| [ ] T14 | Settings API | GET/PUT income and pay day; CRUD fixed costs; validation | Handler tests incl. negative amount rejected |
| [ ] T15 | Goals store | Migration and queries for goals | Store tests |
| [ ] T16 | Goals API and on-track | CRUD, pause/close, computed on-track flag | Test: behind vs on-track |
| [ ] T17 | Derived numbers | Service for `monthly_free`, `remaining_free`, emergency cover, suggested target | Unit tests per formula |
| [ ] T18 | Months | Migration; lazy creation of current month from pay day | Test: day before and on pay day |
| [ ] T19 | Income override and extra income | One-month override and one-off income endpoints | Test: default unchanged next month |
| [ ] T20 | Checks store | Migration and queries for checks with outcome | Store tests |
| [ ] T21 | POST /checks | Service loads state, runs engine, uses an `Explainer` interface (stub returns template), saves | E2E handler test, rules path under 100 ms |
| [ ] T22 | Spends | Migration, log/edit/delete spend, need or want | `remaining_free` drops after a spend |
| [ ] T23 | Check outcomes | Bought / bought anyway / skipped; bought creates a linked spend | Test both outcome paths |
| [ ] T24 | Save for this | Create wish goal from a check; block contribution that makes `monthly_free` negative | Test the block message |
| [ ] T25 | Contribution confirmation | List planned contributions for the month; mark paid updates goal | Test unpaid leaves goal behind |

## Phase 3: Core screens (MULL-11, 12, 13, 15)

Check every screen at phone width.

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [ ] T26 | Web foundation | API client, TanStack Query, routes, mobile layout, `lib/money.ts` | Unit test for money formatting |
| [ ] T27 | Setup: income and fixed costs | Steps 1 and 2 | Screenshot at 375 px |
| [ ] T28 | Setup: goals and warning | Step 3, emergency target suggestion, zero free money warning | Timed run under 5 minutes |
| [ ] T29 | Dashboard | Remaining free, wants spent, cover, goal bars, check button | Screenshot with sample data |
| [ ] T30 | Buy check form | Name, price, Need/Want required; reason fields optional | Numeric keyboard on mobile |
| [ ] T31 | Buy check result | Verdict, sentence, number, outcome buttons, Save for this | Screenshots for Buy, Wait, Skip |
| [ ] T32 | Goals screen | List, edit, pause, add money, shows new `monthly_free` before save | Screenshot |
| [ ] T33 | Log a spend and contributions banner | Quick spend form; banner until contributions confirmed | Manual run-through |

## Phase 4: AI layer (MULL-19, 20)

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [ ] T34 | Guard | `ai.Guard(rule, proposed, moneyFits)`: max one step, never to Buy when money does not fit | Test: AI tries Skip to Buy, gets rule verdict |
| [ ] T35 | Prompt and parsing | Prompt template, strict JSON output, one-sentence trim, no identity data | Tests with recorded sample responses, no network |
| [ ] T36 | LLM client and wiring | Real provider behind `Explainer`, model and key from env | Manual check on 3 real items |
| [ ] T37 | Timeout and fallback | 4 s timeout, templated sentence per rule, log failure | Test with a slow fake client |

## Phase 5: History and real use (MULL-18, 21, 22, 23)

| ID | Task | Scope | Done when |
| --- | --- | --- | --- |
| [ ] T38 | History API and summary | List by month/verdict; follow rate, verdict split, wants spent vs last month | Tests on seeded data |
| [ ] T39 | History screen | List, filters, monthly line | Screenshot |
| [ ] T40 | Why panel and thresholds | Show rule trace; settings for rule 1 and 2 cutoffs | Changing a cutoff changes a verdict in a test |
| [ ] T41 | CSV export | Checks and spends as CSV | File opens in a spreadsheet |
| [ ] T42 | Login | Single password, session cookie | Unauthenticated request gets 401 |
| [ ] T43 | Deploy and backups | VPS, Nginx, HTTPS, daily `pg_dump`, 14-day retention | Restore tested once |
| [ ] T44 | Month of use | Use it daily for a month; log issues in `learnings.md`; tune cutoffs | Retro written |

## Later

T45 ready to re-check (MULL-24), T46 48-hour cool-off (MULL-25), T47 emergency fund use for needs (MULL-26). Plan these after T44.
