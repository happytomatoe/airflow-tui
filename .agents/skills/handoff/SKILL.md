---
name: handoff
description: Create or resume Codex handoff documents for this repository. Use when the user wants to capture current progress for a later session or restart from an existing handoff with validation against current code.
---

# Handoff

Use this skill for either creating a handoff or resuming from one.

## Create Handoff

1. Determine the target ticket or use `general` when none exists.
2. Create the handoff in the repo's established handoff location:
   - prefer `thoughts/shared/handoffs/EMS-XXXX/` when that structure exists
3. Gather:
   - current date and timezone
   - current branch
   - current commit
   - repository name
   - task status
   - recent changes
   - key learnings
   - next steps
4. Write YAML frontmatter and concise sections for:
   - task summary
   - critical references
   - recent changes
   - learnings
   - artifacts
   - action items and next steps
   - other notes
5. Respond with the handoff path and how to resume from it later.

## Resume Handoff

1. Read the handoff fully, or locate the most recent relevant handoff if the user
   supplied only a ticket ID.
2. Read any linked plan or research documents that still matter.
3. Verify the handoff against the current codebase with focused research.
4. Present:
   - original tasks
   - current status
   - validated learnings
   - recent changes still present or missing
   - recommended next actions
5. Get confirmation before starting new implementation work.
6. Use `update_plan` to track resumed work when implementation begins.

## Rules

- Be thorough, but prefer `file:line` references over large code snippets.
- Do not assume the codebase still matches an old handoff. Verify first.

