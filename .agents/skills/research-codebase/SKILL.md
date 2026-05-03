---
name: research-codebase
description: Research and document how the codebase works today. Use when the user wants a technical map, implementation explanation, or historical context, with either thoughts-aware, generic, or no-thoughts research modes.
---

# Research Codebase

Use this skill for documentation and code understanding work, not for critique by
default.

## Workflow

1. Read any directly mentioned files first.
2. Break the question into smaller research tracks.
3. Decide the research mode:
   - default: use both live code and relevant `thoughts/` documents
   - generic: do not assume a full `thoughts/` layout exists
   - no-thoughts: rely only on the live codebase
4. Use `update_plan` to track the investigation.
5. Spawn parallel explorers guided by:
   - `.codex/agents/codebase-locator.md`
   - `.codex/agents/codebase-analyzer.md`
   - `.codex/agents/codebase-pattern-finder.md`
   - `.codex/agents/thoughts-locator.md` when `thoughts/` is in scope
   - `.codex/agents/thoughts-analyzer.md` when deeper document analysis is needed
6. Use `.codex/agents/web-search-researcher.md` only when the user explicitly asks
   for external research or current information.
7. Wait for all research to complete, then synthesize findings.
8. Save the result in the repo's research-notes location when the user wants a saved
   artifact:
   - prefer `thoughts/shared/research/` when it exists and is in scope
   - otherwise use the repo's existing research-notes location

## Rules

- Document what is true now.
- Do not critique or recommend unless the user asks.
- Prefer live code as the primary source of truth.
- Include concrete `file:line` references for important claims.

