<!-- AI-CONTEXT
src: v0.2
phase: M1
direction: Implement LINE webhook handler + ticket submission flow
focus: [T-002, T-003]
done: [T-000, T-001, T-005, T-008, T-004, T-007, T-006]
future: [T-009 LIFF — post-MVP]
blocked: none
next: T-002 > T-003
risk: none active
adr: ADR-001
read_more:
  prd: doc/00-source/versions/v0.2/01-prd.md
  migration_design: doc/06-extensions/T-004-migration-002-design.md
  architecture: doc/07-decisions/README.md
  entities: doc/07-decisions/entity-register.md
  source_current: doc/00-source/versions/v0.2/
updated: 2026-04-30
-->

---

# Project Status — Lotto Journal

Last updated: 2026-04-30 (session 4)

## Source References

- `doc/00-source/versions/v0.2/01-prd.md` — current PRD (LINE-based)
- `doc/07-decisions/ADR-001-line-messaging-pivot.md` — Accepted

---

## Phase and Direction

**Current phase:** M1 — Design & Build

ADR-001 has been accepted (Option B). M1 work remaining:

1. Design and implement the LINE webhook handler and ticket submission flow (T-002)
2. Design the cronjob for lottery result fetch + comparison (T-003)

The cronjob (M2) and win notification (M3) follow after M1 is stable.

**Post-MVP direction:** A LIFF (LINE Front-end Framework) web app is planned to complement
the chatbot. The monorepo structure is intentionally preserved for this. See T-009.

---

## Active Tasks

- `T-002` — Design LINE Messaging API webhook handler — todo
- `T-003` — Design cronjob: lottery result fetch + comparison flow — todo

---

## Completed Tasks

- `T-000` — Documentation setup (2026-04-30)
- `T-001` — Architecture pivot decided: Option B (LINE Messaging API) — ADR-001 Accepted (2026-04-30)
- `T-005` — Formal source docs written: PRD v0.2 created (2026-04-30)
- `T-008` — `trunk/glo_result.json` committed by owner (2026-04-30)
- `T-004` — User identity schema designed; DBML updated; owner approved (2026-04-30)
- `T-007` — Migration 000002 written (up + down); Go model + code updated; build passes (2026-04-30)
- `T-006` — `apps/web` deleted; `turbo.json` + `pnpm-workspace.yaml` cleaned up (2026-04-30)

---

## Blocked Tasks

- None currently

---

## Next Steps

1. **T-002:** Design LINE webhook handler (verify signature, parse events, follow/unfollow/message)
2. **T-003:** Design cronjob — `trunk/glo_result.json` now committed; ready to design

---

## Future Direction

- **T-009 — LIFF app (post-MVP):** A LIFF (LINE Front-end Framework) web app is planned to
  complement the chatbot. Lives in `apps/liff`. Monorepo structure (Turbo, pnpm workspaces)
  intentionally kept for this purpose. Task: design when post-MVP phase begins.

---

## Risks and Notes

- No active risks.
- **Migration 000002 notes (for reference):**
  - `account_status` was drop+recreated (PostgreSQL has no `DROP VALUE`)
  - `user_winnings.user_id` was missing from 000001 SQL — added in 000002 [FOUND-IN-PASSING]
  - Enum renames: `N6`→`L6`, `n6_*`→`l6_*` (9 prize_type values)
  - `AUTH` code removed: `auth_handler.go`, `auth_service.go`, `auth_dto.go` deleted
