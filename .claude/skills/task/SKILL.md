---
name: task
description: Start a task from docs/plan.md by ID, e.g. /task T07. Plans first, waits for approval, then implements and verifies.
disable-model-invocation: true
---
Work on task $ARGUMENTS from docs/plan.md.

1. Read CLAUDE.md, the task's row in docs/plan.md, and the linked ticket in docs/backlog.md.
2. Read only the files the task touches. Do not explore the whole repo.
3. Write a short plan: files to create or change, functions and types, test cases in plain words, and anything out of scope you noticed. Stop and wait for approval.
4. After approval, write the tests first and show them failing.
5. Implement until `make check` passes. Paste the output.
6. Run the code-reviewer agent on the diff. Fix "must fix" findings only.
7. Draft the PR description using the template in docs/ai-workflow.md.
8. Do not commit, push, or tick the task. The human does that.
