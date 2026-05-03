---
name: create-worktree
description: Prepare a worktree for follow-on implementation work in this repository. Use when the user wants a new worktree plus the exact branch, path, plan reference, and launch prompt needed for the next Codex session.
---

# Create Worktree

Use this skill when a task should move into a dedicated worktree before
implementation.

## Workflow

1. Read `hack/create_worktree.sh` if it exists.
2. Determine:
   - branch name
   - worktree path
   - plan path
   - launch prompt
3. Show the user the exact derived values before executing.
4. Create the worktree only after the derived values are confirmed.
5. Prefer relative `thoughts/shared/...` paths in launch prompts when the repo syncs
   those paths across worktrees.

## Output

Report:
- worktree path
- branch name
- plan path
- launch prompt or command to use next

## Rules

- Adapt the launch command to the actual local Codex workflow.
- Do not hard-code legacy launcher tooling.

