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
6. `00-source/README.md` — source docs index
7. `07-decisions/README.md` — ADR index (read before any architecture decision)

## Folder Structure

| Folder            | Purpose                                                    |
| ----------------- | ---------------------------------------------------------- |
| `00-source/`      | Source docs (requirements, PRD, specs) — versioned         |
| `01-plan/`        | Project plan and work status                               |
| `02-task/`        | Task board                                                 |
| `03-log/`         | Work logs and session history                              |
| `04-way-of-work/` | Working guidelines, coding standards, AI decision protocol |
| `05-summary/`     | Monthly and milestone summaries                            |
| `06-extensions/`  | Extension docs for scope additions                         |
| `07-decisions/`   | Architecture Decision Records (ADR) and entity register    |

## Source Version in Use

Current: **`v0.1`** — ⚠️ No formal source docs yet. See `00-source/README.md`.

## Key Pending Decision

> ⚠️ **Architecture pivot under active consideration:** migrating user interaction
> from the current web app to LINE Messaging API.
>
> **No new implementation should start on user-facing features until
> `07-decisions/ADR-001-line-messaging-pivot.md` is resolved.**
