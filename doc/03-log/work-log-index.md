<!-- AI-CONTEXT
last_session: 2026-05-08 (session 6)
tool: Claude (Sonnet 4.6)
completed: [T-010, T-011]
in_progress: []
checkpoint: none
next_from_last: T-003
notes: T-010 done. Middleware (recover+requestid+Logging+timeout) implemented then immediately upgraded: fiber v2→v3 (v3.2.0) after deprecated timeout.New warning. T-011 done. GET /health with DB ping. Build passes. README + webhook-flow.md updated.
deep_context: doc/06-extensions/T-004-migration-002-design.md
-->

---

# Work Log Index — Lotto Journal

Last updated: 2026-05-08 (session 6)

---

## Milestone Summary

_(Updated when milestones close — never archived)_

- **M0 complete (2026-04-30):** ADR-001 accepted (Option B — LINE Messaging API).
  PRD v0.2 written. Entity register updated. doc/ structure established.

---

### 2026-05-08 — Session 6 — [Claude (Sonnet 4.6)]

- **Session summary:** T-010 (middleware hardening) implemented, then immediately upgraded to Fiber v3 after the deprecated `timeout.New` warning surfaced. Fetched official Fiber v3 docs, migrated the full codebase from v2 → v3 (v3.2.0). T-011 (GET /health) implemented. Updated `apps/api/README.md` and `trunk/webhook-flow.md` to reflect new middleware stack and corrected a false LINE retry-interval claim. Build passes.
- **Work done:**
  - `internal/handler/health_handler.go`: created; `GET /health`; calls `db.DB().Ping()`; returns `{"status":"ok","db":"ok"}` (200) or `{"status":"degraded","db":"<err>"}` (503)
  - `app/main.go`: wired `healthHandler`; registered `GET /health` (no timeout wrapper)
  - `middlewares/log.go`: upgraded to log status code + request ID; `c.Locals("requestid")` → `requestid.FromContext(c)` (v3 API); handler sig `*fiber.Ctx` → `fiber.Ctx`
  - `app/main.go`: `recoverer.New(EnableStackTrace: true)` globally; `requestid.New()` globally; `/webhook` wrapped with `timeout.New(handler, timeout.Config{Timeout: 25s})` (v3 race-free timeout); `log.Fatal(app.Listen(...))` (v3 always returns error); `recover` import aliased as `recoverer` (v3 idiom)
  - `internal/handler/line_handler.go`: import `v2` → `v3`; handler sig `*fiber.Ctx` → `fiber.Ctx`
  - `go.mod`: added `github.com/gofiber/fiber/v3 v3.2.0`; `go mod tidy` removed v2 entirely
  - `apps/api/README.md`: added **Middleware stack** section (table + log format + timeout rationale with LINE redelivery reference)
  - `trunk/webhook-flow.md`: updated Big Picture diagram (middleware chain shown); Step 1 (middleware registrations + timeout-wrapped route); Step 2 (handler sig); Step 5 (removed false "within 30 seconds" claim, linked LINE docs); Files Involved table
- **Decisions resolved this session:**
  - Fiber v3 replaces v2 — `timeout.New` in v3 is race-free via Abandon mechanism; no `NewWithContext` needed
  - `requestid.FromContext(c)` is the v3 accessor — v3 drops the `ContextKey` config field
  - `recover` middleware import aliased as `recoverer` to avoid shadowing the Go built-in
  - LINE's webhook retry count/interval is not publicly disclosed — removed the false "30 seconds" claim from webhook-flow.md
- **Tasks changed:**
  - T-010: done (build passes, all middleware wired, Fiber v3)
  - T-011: done (build passes)
- **Awaiting owner action:** None
- **Daily Log:** _(local only — not committed)_

---

### 2026-05-07 — Session 5 — [Claude (Sonnet 4.6)]

- **Session summary:** T-002 (LINE webhook handler) fully designed and implemented. LINE Bot SDK v8 added. All layers built from scratch (models, repos, services, handler). Migration 000003 added for idempotency. Build passes.
- **Work done:**
  - Added `github.com/line/line-bot-sdk-go/v8` (v8.20.0) to go.mod
  - Created `migrations/000003_webhook_events.up/down.sql` — idempotency table
  - Created `models/draw.go`, `models/ticket.go`, `models/webhook_event.go`
  - Updated `repository/user_repository.go` — added `FindByLineUserID`, `FindOrCreate`, `UpdateStatus`
  - Created `repository/draw_repository.go` — `FindByDate`, `FindOrCreate`
  - Created `repository/ticket_repository.go` — `Create`
  - Created `repository/webhook_event_repository.go` — atomic `MarkProcessed` (ON CONFLICT DO NOTHING)
  - Created `service/user_service.go` — `FindOrCreate`, `Deactivate`
  - Created `service/draw_service.go` — `NextDrawDate` (Bangkok time, 1st/16th logic), `FindOrCreateUpcoming`
  - Created `service/ticket_service.go` — `ParseTicketInput` (commas+spaces, xN quantity, 3/6-digit validation), `SubmitTickets`
  - Created `handler/line_handler.go` — Fiber→SDK bridge (synthetic `*http.Request`), `follow`/`unfollow`/`message` routing, reply builder
  - Updated `config/config.go` — added `LINE_CHANNEL_SECRET`, `LINE_CHANNEL_ACCESS_TOKEN`
  - Updated `app/main.go` — wired all layers; registered `POST /webhook`
- **Decisions resolved this session:**
  - LINE SDK bridge pattern: build synthetic `*http.Request` from Fiber context to pass to `webhook.ParseRequest()`
  - Idempotency: `webhook_events` table with `ON CONFLICT DO NOTHING` insert + `RowsAffected` check
  - `NextDrawDate`: Bangkok time, candidates = [1st, 16th of current month, 1st of next month], first >= today
  - `ParseTicketInput`: normalise commas→spaces, merge `\d+ x\d+` pattern, validate 3/6-digit only
  - `UserSource` type assertion is value type (`webhook.UserSource`, not pointer)
  - Microcopy: Thai language placeholder (PRD marks it TBD)
- **Tasks changed:**
  - T-002: done (build passes, all events handled)
- **Awaiting owner action:** Run `make migrate-up` against the DB to apply migration 000003 before testing
- **Daily Log:** _(local only — not committed)_

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
