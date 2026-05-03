---
name: iterate-plan
description: Revise an existing implementation plan based on feedback or new constraints. Use when the user wants a saved plan updated without starting planning from scratch, with fresh research when needed.
---

# Iterate Plan

Use this skill when a plan exists and needs to be revised rather than replaced.

## Workflow

1. Read the plan fully.
2. Parse the requested changes and determine whether `thoughts/` documents are in
   scope or intentionally excluded.
3. If fresh technical context is needed, research with explorers guided by:
   - `.codex/agents/codebase-locator.md`
   - `.codex/agents/codebase-analyzer.md`
   - `.codex/agents/codebase-pattern-finder.md`
   - `.codex/agents/thoughts-locator.md` when `thoughts/` is in scope
   - `.codex/agents/thoughts-analyzer.md` when `thoughts/` is in scope
4. Read any newly identified source files fully.
5. Present your understanding before editing when the requested change is
   substantial.
6. Edit the plan surgically with `apply_patch`.

## Rules

- Preserve the good parts of the plan.
- Keep file references and success criteria concrete.
- Keep automated and manual verification distinct.
- Do not leave unresolved questions in the plan text.

