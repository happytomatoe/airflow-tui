---
name: validate-plan
description: Validate whether an implementation matches its saved plan. Use when the user wants a gap analysis between a plan and the current code, including verification results, deviations, and remaining manual checks.
---

# Validate Plan

Use this skill when the user wants an implementation audited against its plan rather
than continued.

## Workflow

1. Read the plan fully.
2. Identify expected file changes, success criteria, and verification commands.
3. Inspect recent commits and diffs when needed.
4. Run the plan's automated checks when practical.
5. Compare claimed completion against the actual code.
6. Produce a validation report with:
   - implementation status by phase
   - automated verification results
   - deviations from the plan
   - missing coverage
   - manual verification still required

## Rules

- Be explicit about what was verified versus inferred.
- Prioritize real gaps, regressions, and incomplete work.
- Keep manual verification steps actionable.

