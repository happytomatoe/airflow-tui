---
name: describe-pr
description: Generate a complete pull request description for this repository. Use when the user wants a PR body drafted or updated from the actual diff, with either the repo template or a neutral fallback template.
---

# Describe PR

Use this skill when the user wants a PR description written from the actual branch
state.

## Workflow

1. Read `thoughts/shared/pr_description.md` if it exists and is in scope.
2. Find the target PR with `gh pr view` or `gh pr list`.
3. Check for an existing saved draft under `thoughts/shared/prs/` when that
   structure exists.
4. Gather:
   - PR metadata
   - diff
   - commits
   - base branch
5. Read referenced files needed for context.
6. Run practical verification commands and record what passed versus what still
   needs manual testing.
7. Write the completed description:
   - prefer `thoughts/shared/prs/{number}_description.md` when the repo uses that
     convention
   - otherwise use a neutral temporary draft path
8. Update the PR body with `gh pr edit --body-file ...` when the user wants the PR
   updated directly.

## Fallback Template

Use this template when no repo-local template is available:

```md
## What problem(s) was I solving?

## What user-facing changes did I ship?

## How I implemented it

## How to verify it

### Manual Testing

## Description for the changelog
```

## Writing Guidance

- Be specific about the problem solved.
- Separate user-facing changes from internal implementation details.
- Keep the result scannable.

