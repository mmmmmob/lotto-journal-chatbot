<!-- AI-CONTEXT
last_session: 2026-04-30 (session 4)
tool: Claude (Sonnet 4.6)
completed: [T-008 (by owner), T-004, T-007, T-006]
in_progress: []
checkpoint: none
next_from_last: T-002 > T-003
notes: Full M1 setup complete. Migration 000002 done. apps/web removed. LIFF noted as T-009 post-MVP. Monorepo kept for LIFF. Only T-002 and T-003 remain active.
deep_context: doc/06-extensions/T-004-migration-002-design.md
-->

---

# Work Log Index — Lotto Journal

Last updated: 2026-04-30 (session 4)

---

## Milestone Summary

_(Updated when milestones close — never archived)_

- **M0 complete (2026-04-30):** ADR-001 accepted (Option B — LINE Messaging API).
  PRD v0.2 written. Entity register updated. doc/ structure established.

---

### 2026-04-30 — Session 4 — [Claude (Sonnet 4.6)]

- **Session summary:** T-004 design completed. T-007 (migration 000002) written. T-006 (remove apps/web) done. LIFF planned as T-009 post-MVP. T-008 closed (glo_result.json committed by owner).
- **Work done:**
  - Analysed migration 000002 scope; produced design doc `doc/06-extensions/T-004-migration-002-design.md`
  - Updated `trunk/db_diagram.dbml` to post-000002 target state
  - Written `000002_line_identity.up/down.sql`; updated `models/user.go`; removed auth code; build passes
  - Deleted `apps/web`; cleaned `turbo.json` (removed generate/prisma tasks); cleaned `pnpm-workspace.yaml`
  - Extended `Makefile` and `package.json` scripts (db:start, db:stop, migrate:\* targets)
  - Updated README with pnpm-first setup guide
  - Added T-009 (LIFF planning) to task board as post-MVP future task
  - Added LIFF note to PRD v0.2 §8 (Out of Scope) — intentionally deferred, not abandoned
- **Decisions resolved this session:**
  - `account_status` enum: DROP `pending`, ADD `inactive`, KEEP `active`+`suspended`
  - `user_winnings.user_id` [FOUND-IN-PASSING]: added in migration 000002
  - Monorepo structure kept intentionally — LIFF app will use `apps/liff` in future
- **Tasks changed:**
  - T-008: done (committed by owner)
  - T-004: done (design approved)
  - T-007: done (migration written, build passes)
  - T-006: done (apps/web deleted, configs cleaned)
  - T-009: added (future / post-MVP)
- **Awaiting owner action:** None — ready to start T-002 or T-003
- **Daily Log:** _(local only — not committed)_

---

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
