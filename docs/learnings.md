# Learnings

One line per surprise: what the AI got wrong, the prompt that fixed it, the rule added.

- T01: AI drafted the PR description without reading the template in `docs/ai-workflow.md`; fixed by asking it to reshape to the template. Rule: read `docs/ai-workflow.md` before drafting a PR.
- T01: I gave a wrong API (e.Shutdown, not in Echo v5). Claude checked the installed library instead of guessing. Always verify APIs against the installed version.
