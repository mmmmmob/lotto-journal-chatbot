# Compliance — Lotto Journal

This document defines the compliance scan protocol for this project.
The scan runs automatically every session unless paused.

Source: adapted from `core/15-compliance-check-template.md`

---

## Control Commands

| Command             | Effect                                   |
| ------------------- | ---------------------------------------- |
| _(no command)_      | Scan runs automatically every session    |
| `pause compliance`  | Pause scan for this session              |
| `resume compliance` | Re-enable scan                           |
| `scan`              | Run scan immediately; output full report |
| `scan refactor`     | Scan REFACTOR-PENDING items only         |

---

## What Gets Scanned

### Level 1 — Fix Now or Defer (must decide)

| Code | What is checked                                | Threshold                                       |
| ---- | ---------------------------------------------- | ----------------------------------------------- |
| C-01 | File size                                      | > 500 lines                                     |
| C-02 | Task has no source reference                   | Any task without `doc/00-source/...` or ADR ref |
| C-03 | Task marked `done` without validation evidence | No record of how it was verified                |
| C-04 | Placeholder still present in a file            | `<PROJECT_NAME>`, `<NEEDS_CLARIFICATION>`, etc. |

### Level 2 — Always Defer (tag in code)

| Code | What is checked                        | Threshold                                 |
| ---- | -------------------------------------- | ----------------------------------------- |
| C-05 | Function/method too long               | > 50 lines                                |
| C-06 | New dependency without ADR             | Import not previously used in the project |
| C-07 | AI-CONTEXT block out of sync with body | Value in block ≠ value in body            |
| C-08 | TODO / FIXME without task reference    | Comment missing T-XXX reference           |

### Level 3 — Notice Only

| Code | What is checked                                | Threshold                                             |
| ---- | ---------------------------------------------- | ----------------------------------------------------- |
| C-09 | work-status not updated while task in progress | `in_progress` task, status not updated > 3 days       |
| C-10 | No log entry for last session                  | `work-log-index` missing the latest session           |
| C-11 | Security baseline failed                       | See security rules below                              |
| C-12 | work-log-index too large                       | > 300 lines — recommend archive                       |
| C-13 | task-board done section too large              | > 15 items — recommend archive                        |
| C-14 | entity-register not updated when tech changed  | Tech added/deprecated but entity-register not updated |

---

## Security Baseline (C-11)

Check every session when code is modified — flag immediately if found:

| Issue                            | What to check                                                      |
| -------------------------------- | ------------------------------------------------------------------ |
| Hardcoded secrets                | API keys, passwords, tokens in code or unencrypted config          |
| SQL injection                    | Query built by string concatenation instead of parameterized query |
| XSS                              | HTML rendered from user input without sanitization                 |
| Insecure direct object reference | User-supplied ID used without ownership check                      |
| Sensitive data in logs           | Logs containing passwords, tokens, or PII                          |
| LINE channel secret exposed      | Channel secret hardcoded or logged                                 |
| Known vulnerable dependency      | Package version with a known CVE                                   |

If C-11 is found → **Fix immediately before continuing any other work. No deferral.**

---

## Violation Tag Format

Use this format in code comments when deferring:

```
// REFACTOR-PENDING[C-01]: file too long (620 lines), needs splitting — T-XXX
// REFACTOR-PENDING[C-05]: function too long (65 lines) — T-XXX
// REFACTOR-PENDING[C-06]: new lib introduced, ADR needed — T-XXX
// REFACTOR-PENDING[C-08]: TODO without task reference — T-XXX
```

Format: `// REFACTOR-PENDING[C-XX]: <description> — <T-XXX>`

---

## Fix vs Defer Decision

**Fix now when:**

- Within the scope of the current task (no extra cost)
- Fix takes less than ~30 minutes
- Affects correctness, not just cleanliness

**Defer when:**

- Outside the current task's scope
- Would take more than ~30 minutes
- Pure refactor with no behavior change

When deferring → create a task T-XXX in task-board immediately, tagged `[COMPLIANCE-DEFER]`

---

## Report Format

```
=== Compliance Report — YYYY-MM-DD ===

LEVEL 1 (Fix or Defer):
  [C-01] apps/api/internal/handler/xxx.go — NNN lines (limit: 500)
  [C-03] T-XXX marked done — no validation evidence

LEVEL 2 (Defer):
  [C-05] apps/api/internal/service/xxx.go:FunctionName() — NN lines
  [C-06] import 'xxx' in apps/api/... — no ADR found

LEVEL 3 (Notice):
  [C-09] T-XXX in_progress — work-status not updated for N days

Action required:
  - Fix or create REFACTOR-PENDING tasks for Level 1 items
  - Create REFACTOR-PENDING tasks for Level 2 items
  - Note Level 3 items in work-log
```
