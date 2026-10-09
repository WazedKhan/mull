# Mull

Personal savings coach. Before a purchase, the user gets a Buy / Wait / Skip verdict with one sentence of reasoning. Single user, BDT only, self-hosted.

Read before working: `docs/user-journeys.md` (what and why), `docs/plan.md` (task list), `docs/ai-workflow.md` (how we work).

## Commands

- `make run` starts the API, web and Postgres (docker compose)
- `make test` runs Go tests; `make test-engine` runs only the rules engine
- `make lint` runs golangci-lint and the web linter
- `make check` runs lint, tests and the web build. IMPORTANT: run it before saying a task is done, and paste the result
- `make migrate-new name=<x>` creates a migration; `make migrate` applies them
- `make sqlc` regenerates the query code after editing `internal/store/queries/*.sql`

## Layout and dependency rules

```
cmd/api/            main, wiring only
internal/money/     Money type (int64 BDT). No other imports from this repo
internal/engine/    rules engine. Pure functions, imports only money. No DB, HTTP, time.Now or AI
internal/ai/        LLM explanation + guard. Depends on engine types
internal/store/     Postgres via sqlc. Migrations in internal/store/migrations
internal/service/   use cases; owns interfaces for store and ai
internal/http/      Echo handlers: parse, call service, render. No business logic
web/                React + TypeScript (Vite)
```

- Dependencies point inward: http -> service -> engine/store/ai interfaces. Never import `http` or `store` from `engine`.
- Define interfaces in the package that uses them, not the one that implements them.
- Pass `time.Time` ("now") and thresholds into the engine as inputs so tests are deterministic.

## Domain rules (never break these)

- Money is always `money.Money` (int64 BDT). No float64 for money in Go, SQL or TypeScript.
- The rules engine owns the verdict. The AI may move it at most one step and never to Buy unless the money fits. Enforce this in `ai.Guard`, with tests.
- Use the names `monthly_free` and `remaining_free`; never just "free money" in code.
- Never edit a migration that is already on `main`; add a new one.
- No bank linking, no analytics, no data sent anywhere except the LLM call, which must contain no identity data.

## Code style

- Go: wrap errors with `fmt.Errorf("doing x: %w", err)`; use `log/slog`; pass `context.Context` first.
- Tests are table-driven. Engine changes need a test per rule branch. Store tests run against the real Postgres from docker compose, not mocks.
- Web: TanStack Query for server data; format money only through `web/src/lib/money.ts`.
- No new dependency without asking first and saying why the standard library is not enough.

## Workflow

- One task from `docs/plan.md` per branch and PR. Branch: `t07-engine-buy-rules`.
- Stay inside the task's scope. If you find something else, add it to "Follow-ups" in the PR description instead of fixing it.
- Target under ~300 changed lines per PR, excluding generated code.
- Commits: Conventional Commits (`feat(engine): add rule 3 wait`).
- Before coding a non-trivial task, write a short plan and wait for approval.
- When done: `make check`, then a fresh-context review (`code-reviewer` agent), then the PR description from the template in `docs/ai-workflow.md`.
- Do not push or open PRs. The human reviews and pushes.
