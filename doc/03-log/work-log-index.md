<!-- AI-CONTEXT
last_session: 2026-04-30
tool: Human (direct PRD edit)
completed: []
checkpoint: none
next_from_last: T-008 T-004 T-007 T-006 T-002 T-003
notes: PRD v0.2 placeholders resolved by owner. Two issues flagged: enum rename N6→L6 expands T-007 scope; trunk/glo_result.json missing (T-008 added).
deep_context: none
-->

---

# Work Log Index — Lotto Journal

Last updated: 2026-04-30

---

## Milestone Summary

_(Updated when milestones close — never archived)_

- **M0 complete (2026-04-30):** ADR-001 accepted (Option B — LINE Messaging API).
  PRD v0.2 written. Entity register updated. doc/ structure established.

---

### 2026-04-30 — Session 3 — [Human]

- **Session summary:** Owner filled in all `<NEEDS_CLARIFICATION>` placeholders in PRD v0.2
  by direct file edit. PRD status effectively complete (only microcopy deferred).
- **Decisions resolved:**
  - Message input format: plain number(s), comma/space separated, optional `xN` quantity suffix
  - Ticket type terminology: `L6` (not `N6`) for 6-digit; `N3` unchanged
  - GLO API endpoint: `POST https://www.glo.or.th/api/lottery/getLatestLottery`
  - GLO API response format: documented in `trunk/glo_result.json` (to be committed)
  - Draw schedule: 1st & 16th monthly, **configurable** to handle holiday exceptions
  - Cronjob trigger time: 16:00 Bangkok time (GMT+0700)
  - Retry on API failure: 5 times
  - Non-win notification: YES — send push message even for non-winning tickets
  - `unfollow` event: mark user `status = inactive` (soft delete)
  - `user_profiles` table: DROP in migration 000002
  - `files` table: KEEP — reserved for future photo upload MVP
  - On-demand status query (§3.4): post-MVP
  - Scalability: ≤ 100 users for MVP
  - Monitoring: stdout logs only for MVP
- **New issues flagged by AI review:**
  1. **Enum rename** — PRD now uses `L6` / `l6_*` but migration 000001 created `N6` / `n6_*`.
     T-007 scope expanded to include `ALTER TYPE` rename statements.
  2. **`trunk/glo_result.json` missing** — T-008 added to track committing this file.
- **Tasks added:** `T-008` (commit glo_result.json)
- **Daily Log:** _(local only — not committed)_

---

### 2026-04-30 — Session 2 — [Claude (Sonnet 4.6)]

- **Session summary:** Accepted ADR-001 Option B (LINE Messaging API pivot). Updated ADR
  status to Accepted, filled in rationale and consequences. Updated entity register
  (LINE Messaging API active, Next.js deprecated). Wrote PRD v0.2
  (`doc/00-source/versions/v0.2/01-prd.md`). Updated all planning docs to reference v0.2.
  Advanced project phase from Pre-M0 → M1.
- **Tasks completed:** `T-001` (ADR-001 accepted), `T-005` (PRD v0.2 written)
- **Key decisions recorded in PRD v0.2:**
  - User identity: `line_user_id` replaces email/password
  - Tables to drop: `user_auth_methods`, `user_verifications`
  - Enums to drop: `provider_service`, `verification_type`
  - `users` table redesign → migration 000002 (T-007)
  - `apps/web` to be removed (T-006)
- **Open placeholders in PRD v0.2 that need resolution before implementation:**
  - Thai Gov Lottery API endpoint and response format (blocks T-003)
  - LINE message input format for ticket submission (blocks T-002 implementation)
  - Whether `user_profiles` and `files` tables should be kept or dropped
  - Whether to send "no win" notifications or only win notifications
  - Draw day schedule (1st & 16th assumed — needs confirmation)
- **Daily Log:** _(local only — not committed)_

---

### 2026-04-30 — Session 1 — [Claude (Sonnet 4.6)]

- **Session summary:** Initial documentation setup. Read all 19 core templates (00–18).
  Explored the existing codebase to understand current state before creating any docs.
- **Key findings from code exploration:**
  - `apps/api`: Go + Fiber setup with partial `SignUp` handler (hardcoded password — not production-ready)
  - `apps/web`: Next.js skeleton only (boilerplate pages, no business logic)
  - `trunk/db_diagram.dbml` + migration: Full DB schema exists — users, tickets, draws,
    draw_results, user_winnings, files, enums (account_status, lottery_type, prize_type, etc.)
  - Current `users` table is built for email/password + OAuth — does NOT have `line_user_id`
  - Lottery data model looks solid and likely survives the pivot
- **Tasks completed:** `T-000` (doc setup)
- **ADR created:** ADR-001 (Proposed) — architecture pivot web → LINE Messaging API
- **Validation:** Bootstrap checklist verified before closing session
- **Daily Log:** _(local only — not committed)_
