---
name: implement-plan
description: Execute an approved implementation plan from this repository. Use when the user wants work resumed from a saved plan and expects progress tracking, phased execution, and verification after each phase.
---

# Implement Plan

Use this skill when the user already has an approved plan and wants implementation.

## Workflow

1. Read the plan fully.
2. Read the ticket and all plan-linked references fully.
3. Identify completed checkboxes or phase markers and resume from the first
   incomplete item.
4. Use `update_plan` to track execution progress.
5. Implement one phase at a time.
6. After each phase:
   - run the automated verification required by the plan
   - update completed checkboxes in the plan when the repo tracks progress there
   - summarize what is ready for manual verification
7. If the codebase no longer matches the plan, stop and present:
   - expected state
   - actual state
   - why the mismatch matters
   - a concrete proposed next move

## Rules

- Follow the plan's intent, but verify against the current code.
- Do not mark manual verification complete until the user confirms it.
- Use delegation only for bounded side work that does not block the immediate next
  step.

