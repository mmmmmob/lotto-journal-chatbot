# Lotto Journal — Documentation

This directory is the single source of truth for all project documentation.

## Reading Order (Every AI Session)

Read these files **in order** at the start of every session — the first 3 each have an
`AI-CONTEXT` block at the top; read that block first for fast orientation:

1. `01-plan/work-status.md` — current phase, active tasks, blockers
2. `03-log/work-log-index.md` — what the last session accomplished
3. `02-task/task-board.md` — full task state
4. `01-plan/project-plan.md` — overall plan and milestones
5. `04-way-of-work/coding-standards.md` — coding rules
6. `04-way-of-work/versioning-policy.md` — semver rules (when preparing a release)
7. `04-way-of-work/release-checklist.md` — release steps
8. `00-source/README.md` — source docs index
9. `07-decisions/README.md` — ADR index (read before any architecture decision)

## Folder Structure

| Folder            | Purpose                                                    |
| ----------------- | ---------------------------------------------------------- |
| `00-source/`      | Source docs (requirements, PRD, specs) — versioned         |
| `01-plan/`        | Project plan and work status                               |
| `02-task/`        | Task board                                                 |
| `03-log/`         | Work logs and session history                              |
| `04-way-of-work/` | Working guidelines, coding standards, AI decision protocol, versioning/release rules |
| `05-summary/`     | Monthly and milestone summaries                            |
| `06-extensions/`  | Extension docs for scope additions                         |
| `07-decisions/`   | Architecture Decision Records (ADR) and entity register    |

## Source Version in Use

Current: **`v0.2`** — LINE-based PRD. See `00-source/versions/v0.2/01-prd.md`.

## Current Status

**Phase:** M1 — Design & Build

- ADR-001 accepted (Option B — LINE Messaging API pivot)
- Migration 000002 and 000003 complete (LINE identity redesign + webhook idempotency)
- Active tasks: T-003 (cronjob design) — T-002 complete
- See `01-plan/work-status.md` for full detail
