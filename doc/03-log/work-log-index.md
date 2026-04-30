<!-- AI-CONTEXT
last_session: 2026-04-30
tool: Claude (Sonnet 4.6)
completed: [T-000]
checkpoint: none
next_from_last: T-001 T-005
notes: Architecture pivot (web app vs LINE Messaging API) is the critical open question. ADR-001 created as Proposed. No user-facing implementation until ADR-001 is accepted. Source docs (v0.2) to be written after ADR-001 resolves.
deep_context: none
-->

---

# Work Log Index — Lotto Journal

Last updated: 2026-04-30

---

## Milestone Summary

_(Updated when milestones close — never archived)_

- **Setup (2026-04-30):** Initial `doc/` structure created. Project state assessed.
  Architecture pivot identified as the critical open decision (ADR-001 created as Proposed).

---

## Recent Sessions

### 2026-04-30 — [Claude (Sonnet 4.6)]

- **Session summary:** Initial documentation setup. Read all 19 core templates (00–18).
  Explored the existing codebase to understand current state before creating any docs.
- **Key findings from code exploration:**
  - `apps/api`: Go + Fiber setup with partial `SignUp` handler (hardcoded password — not production-ready)
  - `apps/web`: Next.js skeleton only (boilerplate pages, no business logic)
  - `trunk/db_diagram.dbml` + migration: Full DB schema exists — users, tickets, draws,
    draw_results, user_winnings, files, enums (account_status, lottery_type, prize_type, etc.)
  - Current `users` table is built for email/password + OAuth (provider_service enum: google,
    facebook, apple, local) — does NOT have `line_user_id`
  - Lottery data model looks solid and likely survives the pivot
- **Tasks completed:** `T-000` (doc setup)
- **ADR created:** ADR-001 (Proposed) — architecture pivot web → LINE Messaging API
- **Validation:** Bootstrap checklist verified before closing session
- **Daily Log:** _(local only — not committed)_
