<!-- AI-CONTEXT
src: v0.2
phase: M1
direction: Remove apps/web, redesign user identity (line_user_id), implement LINE webhook + ticket submission
focus: [T-004, T-006, T-007]
done: [T-000, T-001, T-005]
blocked: none
next: T-004 > T-007 > T-006 > T-002 > T-003
risk: migration 000002 scope includes enum rename (lottery_type.N6→L6, prize_type.n6_*→l6_*); trunk/glo_result.json referenced in PRD but not yet committed to repo
adr: ADR-001
read_more:
  prd: doc/00-source/versions/v0.2/01-prd.md
  architecture: doc/07-decisions/README.md
  entities: doc/07-decisions/entity-register.md
  source_current: doc/00-source/versions/v0.2/
updated: 2026-04-30
-->

---

# Project Status — Lotto Journal

Last updated: 2026-04-30

## Source References

- `doc/00-source/versions/v0.2/01-prd.md` — current PRD (LINE-based)
- `doc/07-decisions/ADR-001-line-messaging-pivot.md` — Accepted

---

## Phase and Direction

**Current phase:** M1 — Design & Build

ADR-001 has been accepted (Option B). The project is now moving into M1:

1. Remove `apps/web` (Next.js) — no longer the user-facing product (T-006)
2. Redesign user identity: replace email/password with `line_user_id` via migration 000002 (T-004, T-007)
3. Design and implement the LINE webhook handler and ticket submission flow (T-002)

The cronjob (M2) and win notification (M3) follow after M1 is stable.

---

## Active Tasks

- `T-004` — Redesign user identity model for LINE (line_user_id) — design_validate
- `T-006` — Remove apps/web Next.js app
- `T-007` — Create migration 000002: users table + drop unused tables/enums
- `T-002` — Design LINE Messaging API webhook handler
- `T-003` — Design cronjob: lottery result fetch + comparison flow

---

## Completed Tasks

- `T-000` — Documentation setup (2026-04-30)
- `T-001` — Architecture pivot decided: Option B (LINE Messaging API) — ADR-001 Accepted (2026-04-30)
- `T-005` — Formal source docs written: PRD v0.2 created (2026-04-30)

---

## Blocked Tasks

- None currently

---

## Next Steps

1. **T-004:** Design the new `users` table schema (line_user_id) and migration plan
2. **T-007:** Write migration 000002 (users redesign + drop unused tables/enums)
3. **T-006:** Remove `apps/web` from the monorepo
4. **T-002:** Design LINE webhook handler (verify signature, parse events, handle follow/message)
5. **T-003:** Design cronjob — API endpoint now confirmed; commit `trunk/glo_result.json` first (T-008)

---

## Risks and Notes

- **⚠️ Enum rename — migration 000002 scope expanded:** PRD v0.2 uses `L6` / `l6_*` but
  migration 000001 created `N6` / `n6_*`. Migration 000002 must include
  `ALTER TYPE lottery_type RENAME VALUE 'N6' TO 'L6'` plus 13 rename statements for all
  `prize_type` values (`n6_*` → `l6_*`). See T-007 notes.
- **⚠️ `trunk/glo_result.json` missing from repo:** PRD §6.2 references this file as the
  sample GLO API response. Must be committed before T-003 can be fully designed. See T-008.
- **migration 000002:** No production data exists yet — migration is low-risk.
  Always write and test the `down` migration too.
