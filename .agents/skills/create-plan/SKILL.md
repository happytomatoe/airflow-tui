---
name: create-plan
description: Create implementation plans for this repository from tickets, files, or task descriptions. Use when the user wants a grounded technical plan, including repo-specific thoughts-aware planning, generic planning, or no-thoughts planning.
---

# Create Plan

Use this skill when the user wants a detailed implementation plan rather than code.

## Workflow

1. Read all directly mentioned files fully before delegating.
2. Choose the planning mode that matches the repo and the request:
   - default: use `thoughts/` documents when they exist and are in scope
   - generic: do not assume a `thoughts/` layout beyond what the repo actually has
   - no-thoughts: ignore historical `thoughts/` research and ground the plan in the
     live codebase
3. Research in parallel with `spawn_agent`:
   - use an `explorer` guided by `.codex/agents/codebase-locator.md`
   - use an `explorer` guided by `.codex/agents/codebase-analyzer.md`
   - use an `explorer` guided by `.codex/agents/codebase-pattern-finder.md` when
     similar implementations matter
   - use an `explorer` guided by `.codex/agents/thoughts-locator.md` when
     `thoughts/` is in scope
   - use a default agent guided by `.codex/agents/jira-agent.md` when a Jira ticket
     is involved
4. Read the most relevant files the research identifies.
5. Present only questions that cannot be resolved from code or documents.
6. Use `update_plan` to track research and drafting work.
7. Draft the plan in the repo's established planning location:
   - prefer `thoughts/shared/plans/` when it exists and is in scope
   - otherwise use the repo's existing planning location
   - if no planning location exists, ask the user where the plan should live

## Plan Requirements

- Ground every recommendation in the current codebase.
- Include concrete file references.
- Separate automated verification from manual verification.
- Include out-of-scope notes when they prevent likely confusion.
- Do not leave unresolved questions in the saved plan.

