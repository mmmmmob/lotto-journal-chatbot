# Extensions — Lotto Journal

This folder stores extension documents: scope additions, design decisions for a single
component, or feature ideas that don't yet warrant a new source doc version.

## When to Create an Extension Doc

- New scope that doesn't change the core "product truth"
- Additional details to a specific area (e.g., LINE message format spec)
- Feature ideas still being explored
- Reverse-documentation for undocumented code found during a session

## When to Create a New Source Version Instead

See `doc/00-source/README.md` → Source Version Policy.

## Extension Doc Template

```
# <EXTENSION_TITLE>

Date: YYYY-MM-DD
Status: Draft

## Reason

Why this document exists, and how it differs from existing source docs.

## Related Source References

- `doc/00-source/versions/<VERSION>/<file>.md`

## Description

What is being added or described.

## Impact

- Effect on the plan
- Tasks affected
- Data model / API / UX impact (if any)

## Required Updates

- [ ] doc/01-plan/project-plan.md
- [ ] doc/01-plan/work-status.md
- [ ] doc/02-task/task-board.md
```

## Reverse-Document Protocol

If undocumented code is found:

1. Do not modify the code before documenting it
2. Create `doc/06-extensions/reverse-<name>.md`
3. Create a task `[REVERSE-DOC]` in the task-board
4. Wait for human review before modifying the code

---

## Index

| File                                    | Title                                                                   | Date       | Status |
| --------------------------------------- | ----------------------------------------------------------------------- | ---------- | ------ |
| `T-004-migration-002-design.md`         | Migration 000002 — Design Document                                      | 2026-04-30 | done   |
| `T-020-photo-ocr-openai-r2-proposal.md` | T-020 — Photo Ticket Capture (LINE Image → OpenAI OCR → Confirm → Save) | 2026-05-12 | draft  |
