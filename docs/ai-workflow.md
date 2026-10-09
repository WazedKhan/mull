# How We Build Mull with AI

The goal is two things at once: ship Mull, and get good at AI-assisted development. The AI writes most of the code; you own the decisions, the tests that define "done", and the review. Fast generation is not a reason to review less.

## The five rules

1. **Small tasks.** One task from `plan.md` = one branch = one PR, one axis of change, under ~300 lines. Splitting a plan takes minutes; splitting a finished PR takes hours.
2. **Verification first.** Every task has a "done when" check the AI can run itself (tests, `make check`, a screenshot). If it cannot be verified, it is not ready to build.
3. **Plan before code** for anything touching more than one file. Skip the plan only if you can describe the diff in one sentence.
4. **Fresh context.** One session per task. Clear between tasks. After two failed corrections, stop, clear, and write a better prompt.
5. **Rules live in files.** When the AI makes the same mistake twice, add one line to `CLAUDE.md`. When a rule must never be broken, make it a test, a linter rule or a hook instead.

## The loop for every task

| Step | Who | What |
| --- | --- | --- |
| 1. Pick | You | Take the next task in `plan.md`. Read its scope and "done when". |
| 2. Start clean | You | New session (or `/clear`). Create the branch. |
| 3. Explore and plan | AI, in plan mode | `/task T07`, or: "Read CLAUDE.md, plan.md task T07 and the files it touches. Propose a plan: files, functions, tests. Don't write code yet." |
| 4. Review the plan | You | Is it in scope? Are the tests the right ones? Edit the plan. **This is the highest-value review.** |
| 5. Tests first (engine and services) | You + AI | You list the cases in plain words; AI turns them into a table test that fails. |
| 6. Implement | AI | Implement until `make check` passes. Show the output as evidence. |
| 7. Fresh review | AI subagent | "Use the code-reviewer agent on this diff against task T07." Only fix findings about correctness or scope. |
| 8. Explain back | AI, then you | Ask: "Explain this diff to me as if I'm reviewing it. What would break if X?" If you can't explain it, don't merge it. |
| 9. Your review | You | Use the checklist below. Read every line of the engine, money and guard code. |
| 10. Commit and push | You | Fill in the PR template. Tick the task in `plan.md`. Note anything learned in `docs/learnings.md`. |

## Review checklist for AI-written PRs

- [ ] Diff matches the task's scope; nothing extra snuck in
- [ ] "Done when" check passed and the output is in the PR
- [ ] Tests would fail if the logic were wrong (try breaking one line and rerun)
- [ ] No float money, no `time.Now()` inside the engine
- [ ] Dependency direction respected (no `store` or `http` import in `engine`)
- [ ] Errors are wrapped and handled, not swallowed or `_`-ignored
- [ ] No new dependency without a reason
- [ ] Names match the domain (`monthly_free`, `remaining_free`, verdict names)
- [ ] You can explain every function in the diff

## PR description template

```
## Task
T07: engine Buy rules (MULL-02)

## What changed
- ...

## How it was verified
make check output / test names / screenshot

## Decisions
- ...

## Follow-ups (out of scope)
- ...
```

## Prompt patterns that work

- **Scope + check:** "Implement rule 3 in internal/engine. Cases: price above remaining_free gives Wait; equal gives no match. Run make test-engine and fix failures."
- **Point to a pattern:** "Add the goals handler following the same structure as internal/httpapi/settings.go."
- **Symptom + location + fixed state:** "POST /checks returns 500 when there are no goals. Look at service/check.go. Write a failing test first, then fix."
- **Interview me:** for a fuzzy feature, "Interview me about X, ask about edge cases and tradeoffs, then write the plan."
- **Root cause:** "Fix the cause, don't silence the error or skip the test."

## Anti-patterns

- The kitchen-sink session: several unrelated tasks in one context.
- Approving a plan you skimmed.
- Merging tests you didn't read. AI-written tests often assert what the code does, not what it should do.
- Letting the reviewer agent over-engineer: ignore style nits and "what if" findings outside the requirements.
- Growing `CLAUDE.md` with obvious rules. If the AI already does it right, delete the line.

## Tooling in this repo

- `CLAUDE.md`: loaded every session. Keep it short; prune monthly.
- `.claude/agents/code-reviewer.md`: fresh-context reviewer for step 7.
- `.claude/skills/task/SKILL.md`: `/task T07` runs steps 3 to 6 for a task.
- Later, once a rule keeps being broken: a hook that runs `make lint` after edits.

## Learning log

Keep `docs/learnings.md` with one line per surprise: what the AI got wrong, which prompt fixed it, and which rule you added. After a month, this log is your personal playbook, and a good story for interviews.
