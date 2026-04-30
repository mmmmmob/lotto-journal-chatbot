<!-- AI-CONTEXT
src: v0.1
phase: Pre-M0
direction: Resolve architecture pivot (web app vs LINE Messaging API) before any feature implementation
focus: [T-001]
done: [T-000]
blocked: none
next: T-001 > T-005 > T-002 > T-003 > T-004
risk: Architecture pivot unresolved — no user-facing implementation until ADR-001 is accepted
adr: ADR-001
read_more:
  architecture: doc/07-decisions/README.md
  entities: doc/07-decisions/entity-register.md
  source_current: doc/00-source/versions/v0.1/
  adr_pivot: doc/07-decisions/ADR-001-line-messaging-pivot.md
updated: 2026-04-30
-->

---

# Project Status — Lotto Journal

Last updated: 2026-04-30

## Source References

- `doc/00-source/versions/v0.1/00-setup-placeholder.md` — no formal PRD yet
- `doc/07-decisions/ADR-001-line-messaging-pivot.md` — architecture pivot (Proposed)

---

## Phase and Direction

**Current phase:** Pre-M0

The project has existing code for a web app architecture (Next.js frontend + Go/Fiber API
with email/password auth). The team is considering a significant pivot to **LINE Messaging
API** as the primary user interaction channel.

Before any further implementation on user-facing features, **ADR-001 must be decided**.
Once decided, formal source docs (v0.2) should be written.

---

## Active Tasks

- `T-001` — Decide architecture pivot: web app vs LINE Messaging API (create/finalize ADR-001)

---

## Completed Tasks

- `T-000` — Documentation setup: `doc/` structure created (2026-04-30)

---

## Blocked Tasks

- None currently

---

## Next Steps

1. **T-001:** Resolve ADR-001 — decide whether to keep web app or pivot to LINE Messaging API
2. **T-005:** Write formal source docs (PRD v0.2) after T-001 is resolved
3. **T-002:** Design LINE webhook integration (if LINE pivot is chosen)
4. **T-003:** Design cronjob flow for fetching lottery results
5. **T-004:** Redesign user identity model (especially if LINE pivot changes the users table)

---

## Risks and Notes

- **Architecture pivot unresolved:** Do not start M1 tasks until ADR-001 is accepted.
- **Source docs missing:** All plans are based on the setup brief in `doc/00-source/README.md`.
  Formal PRD must follow after M0.
- **User identity:** The current `users` table uses email/password. A LINE pivot means
  users are identified by `line_user_id` — this changes auth, registration, and profile flows.
- **External API:** Thai Government Lottery Office API format and availability not yet confirmed.
  Confirm before designing the cronjob (T-003).
