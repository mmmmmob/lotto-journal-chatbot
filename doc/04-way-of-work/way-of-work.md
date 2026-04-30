# Way of Work — Lotto Journal

## Purpose

This document defines how we work on the **Lotto Journal** project so that every session
knows: where to start reading, how far the work has progressed, what is done vs pending,
how to update status and logs, and which source docs version to reference.

---

## Files That Must Be Read First

Read in this order every session — the first 3 have an `AI-CONTEXT` block at the top;
read that block first for quick orientation, then read the body only if more detail is needed:

1. `doc/01-plan/work-status.md` — read AI-CONTEXT block first
2. `doc/03-log/work-log-index.md` — read AI-CONTEXT block first
3. `doc/02-task/task-board.md` — read AI-CONTEXT block first
4. `doc/01-plan/project-plan.md`
5. `doc/04-way-of-work/coding-standards.md`
6. `doc/00-source/README.md` + source docs for the current version
7. `doc/07-decisions/README.md` — before any architecture decision

---

## Language Policy

| Context                                                | Language | Reason                               |
| ------------------------------------------------------ | -------- | ------------------------------------ |
| Internal reasoning / thinking                          | English  | Token-efficient; AI processes faster |
| Output to the user                                     | English  | Project language                     |
| Document content (work-status, task-board, logs, etc.) | English  | Humans must read and maintain        |
| AI-CONTEXT blocks                                      | English  | Read by AI, not humans               |
| Code, variable names, technical identifiers            | English  | International convention             |

**Project language: English**

Rule: AI must not output in a different language without explicit permission.

---

## Source of Truth by Purpose

| Purpose                         | Location                                       |
| ------------------------------- | ---------------------------------------------- |
| Business / product requirements | `doc/00-source/versions/<version>/`            |
| Primary plan                    | `doc/01-plan/project-plan.md`                  |
| Current status                  | `doc/01-plan/work-status.md`                   |
| Tasks and tickets               | `doc/02-task/task-board.md`                    |
| Session handoff summary         | `doc/03-log/work-log-index.md`                 |
| Architecture decisions          | `doc/07-decisions/README.md` + individual ADRs |
| Entity/tech status              | `doc/07-decisions/entity-register.md`          |

---

## Compliance Status

Status: **active**
Reason: —

Compliance scan runs automatically every session. To pause: type `pause compliance`.

---

## Core Working Rules

- Do not start implementation without knowing which source doc version is being referenced
- Do not overwrite source doc revisions — create a new version if requirements change
- If new scope appears beyond the current source docs, create an extension doc or a new source version
- Update status and log index at the end of every session, no exceptions

---

## Logging Rules

- Daily log (`doc/03-log/YYYY/MM/YYYY-MM-DD-log.md`) — detailed record for the day
- `work-log-index.md` — session summary for handoff; AI reads this every session
- Daily logs: local only (not committed to git) — commit only monthly summaries and above

---

## Start-of-Session Checklist

- [ ] Read AI-CONTEXT block of `work-status.md` — know current phase/focus/blockers
- [ ] Read AI-CONTEXT block of `work-log-index.md` — know what the last session did
- [ ] Read AI-CONTEXT block of `task-board.md` — know which tasks are active/blocked
- [ ] Check for gaps: do `in_progress` tasks on the board have source references?
- [ ] If gap or inconsistency found → follow Scenario H in `ai-decision-protocol.md`
- [ ] Run compliance scan (unless paused)

## End-of-Session Checklist

Before closing any session, complete all of the following:

- [ ] `work-status.md` — update body **and** AI-CONTEXT block
- [ ] `work-log-index.md` — add new entry **and** update AI-CONTEXT block
- [ ] `task-board.md` — update task status **and** AI-CONTEXT block if any task changed
- [ ] Daily log / daily summary in local workspace
- [ ] Any task left in-progress must have `[IN_PROGRESS: checkpoint saved]` plus a summary
      of what was done so far

**Sync rule:** The AI-CONTEXT block and the body must reflect the same information.
If they differ, trust the body and update the block immediately.

---

## Context Window Management

**Minimal Context Set** — 4 files that give the most information in the least time:

1. `doc/01-plan/work-status.md`
2. `doc/02-task/task-board.md`
3. `doc/03-log/work-log-index.md` (latest entry)
4. `doc/07-decisions/README.md`

**Pre-compact protocol** (before context gets compressed):

1. Save checkpoint in work-log: summarize work done and where it stopped
2. Update work-status to reflect current state
3. Mark task as `[IN_PROGRESS: checkpoint saved — <short summary>]`
4. Note the clear next action for the next session

**Post-compact protocol** (after context is compressed):

1. Re-read the minimal context set before continuing
2. Confirm that the in-progress task matches the saved checkpoint
3. If inconsistency found → follow Scenario B in `ai-decision-protocol.md`

---

## Memory Scope Protocol

When new information is discovered during a session:

| Type of information                           | Where to store                                            |
| --------------------------------------------- | --------------------------------------------------------- |
| Architectural decision                        | `doc/07-decisions/ADR-NNN-*.md`                           |
| New entity or entity status change            | `doc/07-decisions/entity-register.md`                     |
| Pattern / lesson applicable to other projects | `~/ai-workspace/cross-project-memory.md` (ask user first) |
| Progress / detail / in-session decisions      | `doc/03-log/work-log-index.md`                            |
| New task or task status change                | `doc/02-task/task-board.md`                               |

Rule: one piece of information may belong in multiple places — not mutually exclusive.
If unsure: save to work-log first and note "location uncertain".

---

## Multi-AI / Multi-Tool Coordination

This project currently uses a single AI tool. If multiple tools are introduced:

- Create `doc/03-log/agents/<tool-name>.md` per tool
- Each tool writes only its own diary
- `work-log-index.md` remains the master index
- Identify the tool in every work-log entry: `[Claude Code]`, `[Claude.ai]`, etc.
