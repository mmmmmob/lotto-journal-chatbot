# Project Plan — Lotto Journal

Status: Draft baseline
Date started: 2026-04-30

## Source References

- `doc/00-source/versions/v0.2/01-prd.md` — current PRD (LINE-based)
- `doc/07-decisions/ADR-001-line-messaging-pivot.md` — architecture decision (Accepted)
- `trunk/db_diagram.dbml` — data model (post-migration 000002)

---

## Project Objective

**Lotto Journal** is a service that lets users record their lottery ticket numbers and
be automatically notified if any of their tickets win on official draw dates.

Key flows:

1. **Ticket submission** — users record 6-digit or 3-digit ticket numbers
2. **Result fetch** — a scheduled job fetches official results from the Thai Government
   Lottery Office API on every draw day
3. **Result check** — the system compares every user's tickets against the draw results
4. **Win notification** — users whose tickets match are notified automatically

**Architecture:** LINE Messaging API (ADR-001 Accepted 2026-04-30 — Option B chosen).

---

## Scope

### In Scope (MVP)

- User identity and registration (LINE Messaging API per ADR-001)
- Ticket submission: 6-digit and 3-digit lottery numbers
- Draw management: tracking draw dates and status
- Cronjob: scheduled fetch of official lottery results from external API
- Result comparison: checking all user tickets against draw results
- Win notification: alerting users when a ticket wins

### Out of Scope (MVP)

- Ticket image OCR / photo upload for bulk entry
- Ticket resale or marketplace functionality
- Multiple notification channels simultaneously
- Admin dashboard / back-office UI
- Payout or cash settlement handling
- Historical statistics or analytics

---

## Deliverables

- Go backend (Fiber) with business logic for ticket management, draw tracking, and result checking
- Cronjob service (scheduled result fetch + comparison)
- User notification delivery (LINE push message)
- PostgreSQL database with finalized schema (existing schema is a strong foundation)

---

## Milestones

| Milestone | Description                                               | Source Reference              | Status      |
| --------- | --------------------------------------------------------- | ----------------------------- | ----------- |
| M0        | Architecture decided + formal source docs written         | ADR-001, v0.2/01-prd.md       | Done        |
| M1        | User identity redesign + LINE webhook + ticket submission | v0.2/01-prd.md §3.1–3.2, §5.1 | In Progress |
| M2        | Cronjob: result fetch + comparison + win detection        | v0.2/01-prd.md §3.3, §6.2     | Next        |
| M3        | Win notification via LINE push message                    | v0.2/01-prd.md §3.3, §6.1     | Pending     |
| M4        | Hardening: idempotency, error handling, testing, launch   | v0.2/01-prd.md §7             | Pending     |

---

## Risks and Assumptions

### Risks

- **[HIGH — RESOLVED] Architecture pivot:** ADR-001 Accepted. Option B (LINE Messaging API) chosen.
  Entity register and PRD v0.2 updated accordingly.

- **[HIGH — RESOLVED] User identity migration:** Migration 000002 complete. `users` table redesigned
  around `line_user_id`; auth tables and enums dropped.

- **[MEDIUM] External API reliability:** The Thai Government Lottery Office API availability,
  response format, and uptime are not yet confirmed. A fallback strategy may be needed.

- **[MEDIUM] LINE Messaging API limits:** LINE has rate limits and webhook verification
  requirements. These need to be understood before designing the integration.

- **[LOW] apps/web removal:** If the LINE pivot is chosen, the Next.js app may be removed.
  This is low risk as it's currently a skeleton with no production traffic.

### Assumptions

- The Thai Government Lottery Office provides a machine-readable data source for draw results
- LINE Messaging API is the chosen user channel (pending ADR-001)
- Lottery draw days are the 1st and 16th of each month (Thai Government Lottery)
- The existing PostgreSQL schema for tickets, draws, draw_results, and user_winnings
  is fundamentally sound and will be retained (with possible user identity changes)

---

## Change Control

1. If scope expands beyond source docs, create an extension doc (`doc/06-extensions/`)
   or bump to a new source version
2. Update `work-status.md` and `task-board.md`
3. Log decisions in `work-log-index.md`
4. Architecture changes require an ADR in `doc/07-decisions/`

---

## Quality Gates

**Verdict values:** `PASS` | `CONCERNS` | `FAIL`

### Entry Gate (before starting a milestone)

- [ ] Source docs version used is clearly stated in project-plan
- [ ] Every task in the milestone has a source reference
- [ ] No `blocked` tasks without a resolution plan
- [ ] ADR index has been read; no conflicting decisions found

**Verdict Entry Gate:** `PASS` / `CONCERNS` / `FAIL`
If FAIL → do not start milestone until resolved.

### Exit Gate (before closing a milestone)

- [ ] All tasks in the milestone are `done` or have a documented reason
- [ ] Each `done` task passed review and has validation evidence
- [ ] No tasks stuck at `design_validate` or `in_progress`
- [ ] `work-status.md` reflects post-milestone state
- [ ] `work-log-index.md` has an entry for this milestone
- [ ] Architecture decisions have ADRs
- [ ] No scope added without source reference or extension doc

**Verdict Exit Gate:** `PASS` / `CONCERNS` / `FAIL`
If FAIL → milestone is not closed; resolve issues first.
